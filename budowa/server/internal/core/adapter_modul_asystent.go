// Moduł Assistant: typ adaptera, konstruktor i obsługa
// `assistant.voice.command`. Stan zleceń (`assistant.action.status`) i dziennik
// (`assistant.activity.list`) leżą w `adapter_modul_asystent_czynnosci.go`, a
// port `Asystent` i `zarejestrujAsystenta` — w
// `adapter_modul_asystent_uchwyty.go`. Wszystkie te pliki piszą metody na
// jednym typie `*adapterAsystenta`.
//
// Rdzeń rozpoznaje mowę. Pakiet `server/internal/mowa` wnosi silnik, a port
// `Mowa` (`handlers_mowa.go`) wystawia go w kształcie kontraktu. Adapter bierze
// go zależnością opcjonalną (`ZMowa`) i gdy przychodzi samo nagranie —
// przepisuje je, zanim założy zlecenie. Dlatego `assistant.voice.command` ma
// pole `Transcript` osobno od `AudioRef`: tekst poprawiony ręcznie ma
// pierwszeństwo przed rozpoznaniem nagrania.
//
// Bez wpiętego silnika (`ZMowa`) zlecenie powstaje ze śladem nagrania, a
// transkrypcja zostaje pusta; odmowy nie ma — moduł pracuje w zakresie, w
// którym może.
//
// Syntezy mowy to nie dotyczy: `Speak` i `SpeechRef` idą w drugą stronę i rdzeń
// ich nie spełnia — patrz komentarz przy `PolecenieGlosowe`.
package core

import (
	"context"
	"sync"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/models"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// Przedrostki identyfikatorów bytów modułu. Zadeklarowane tu w całości —
// łącznie z przedrostkiem wpisu dziennika, którego ten plik nie nadaje sam —
// żeby nie deklarować ich po raz drugi w innym pliku modułu.
const (
	przedrostekZleceniaAsystenta = "zlec-asyst-"
	przedrostekWpisuAsystenta    = "wpis-asyst-"
)

// adapterAsystenta wypełnia część portu Asystent. Zależność obowiązkowa jest
// jedna: wspólne repozytorium modułu, rozłożone po stronie danych na trzy pliki
// wedle odpowiedzialności, ale niosące jeden typ.
//
// Zależności wykonawcy zleceń są opcjonalne i wpinane w montażu metodami
// budującymi (jak w module Roundtable): rejestr kanałów modelu, nadzorca sesji
// (źródło kanału i katalogów okna) oraz nadajnik zdarzeń. Bez nich moduł
// przyjmuje polecenia i prowadzi ich stan ręcznie, ale nie podejmuje zleceń
// z kolejki — zlecenie zostaje `queued`, zamiast udawać wykonanie, tak jak
// debata bez rejestru kanałów odmawia uruchomienia tury.
type adapterAsystenta struct {
	repozytorium dane.RepozytoriumAsystenta

	kanaly   *models.Rejestr
	nadzorca *session.Nadzorca
	nadajnik Nadajnik
	// mowa rozpoznaje nagranie, gdy Operator nie dosłał tekstu. Zależność
	// opcjonalna i widziana portem kontraktu, nie typem pakietu `mowa`: moduł
	// potrzebuje tu wyłącznie przekładu „nagranie → tekst", a przez port dostaje
	// tę samą instancję silnika, którą obsługuje rodzina `speech.*`.
	mowa Mowa
	// mosty składają konfigurację MCP okna wraz z wpisem serwera narzędzi
	// modelu — tędy model zlecenia dostaje sterowanie platformą
	// (`most_narzedzi.go`). Zależność opcjonalna: bez niej tura zlecenia idzie
	// samym tekstem, bez ani jednego narzędzia
	// (`adapter_modul_asystent_sterowanie.go`).
	mosty *mostyOkna
	// zycie jest kontekstem rdzenia, nie połączenia: tura zlecenia biegnie dłużej
	// niż komenda, a rozłączenie klienta nie ma jej przerywać.
	zycie context.Context
	// biegi wiążą kod zlecenia z przerwaniem jego tury. Bez tego wykazu
	// sterowanie Operatora (`cancel`, `pause`) zmieniałoby sam wiersz, a tura
	// pracowałaby dalej i domykała zlecenie mimo odwołania. Zamki strzegą
	// wykazu, bo tury biegną gorutynami, a sterowanie przychodzi z pętli komend.
	biegi map[string]context.CancelFunc
	// odwolane niesie decyzje Operatora, które zapadły, zanim wykonawca zdążył
	// przestawić zlecenie na `running`. Bez tego znacznika `cancel` zapisuje
	// `cancelled`, a wykonawca — który stan odczytał chwilę wcześniej —
	// nadpisuje go z powrotem na `running` i dowozi zlecenie do `done`.
	// Odwołanie ma być mocniejsze od tury, nie odwrotnie.
	odwolane map[string]bool
	zamki    sync.Mutex
}

// nowyAdapterAsystenta wiąże adapter z repozytorium modułu. Wykonawca zleceń
// dochodzi metodami budującymi (`ZKanalami`, `ZSesjami`, `ZWyjsciem`).
func nowyAdapterAsystenta(repozytorium dane.RepozytoriumAsystenta) *adapterAsystenta {
	return &adapterAsystenta{repozytorium: repozytorium}
}

// ZMowa podpina silnik rozpoznawania mowy. Bez niego moduł przyjmuje nagranie
// bez transkrypcji.
func (a *adapterAsystenta) ZMowa(m Mowa) *adapterAsystenta {
	// Montaż podaje tu wskaźnik, a wskaźnik nil schowany w interfejsie nie jest
	// nil dla porównania `a.mowa == nil`. Bez rozpakowania go brak silnika
	// zamieniłby się w panikę przy pierwszym podyktowanym poleceniu, zamiast
	// w pracę w zakresie niepełnym.
	if m == nil {
		return a
	}
	if adapter, jest := m.(*adapterMowy); jest && adapter == nil {
		return a
	}
	a.mowa = m
	return a
}

// rozpoznajNagranie zamienia nagranie na tekst silnikiem mowy.
//
// Trzy wyniki, trzy zachowania — te same, które rozdziela pole `processed`
// kontraktu `speech.transcribe`:
//
//   - rozpoznano tekst → tekst wraca i staje się treścią polecenia;
//   - przetworzono, mowy brak → odmowa nazywająca ten fakt.
//     `assistant.voice.command` ma wytworzyć polecenie, a puste polecenie nie
//     jest poleceniem: zlecenie z pustą treścią zostawiłoby w dzienniku modułu
//     wpis, po którym nic się nie dzieje;
//   - nie przetworzono → odmowa silnika wraca nietknięta, z jej własnym kodem
//     kontraktu i pełnym komunikatem (`adapter_modul_mowa.go`).
func (a *adapterAsystenta) rozpoznajNagranie(ctx context.Context, odnosnik string) (string, error) {
	wynik, err := a.mowa.Przepisz(ctx, shared.SpeechTranscribeRequest{AudioRef: odnosnik})
	if err != nil {
		return "", err
	}
	if wynik.Transcript == "" {
		return "", bladWskazaniaAsystenta("nagranie " + odnosnik +
			" zostało przetworzone i nie zawiera mowy, więc nie ma z czego złożyć" +
			" polecenia; naprawa: nagrać polecenie ponownie albo przysłać je tekstem")
	}
	return wynik.Transcript, nil
}

// PolecenieGlosowe obsługuje `assistant.voice.command`. Zakłada zlecenie
// i pierwszy wpis dziennika w jednej transakcji — woła `PrzyjmijPolecenie`
// warstwy danych, nie dwa osobne zapisy, żeby nigdy nie powstał wpis bez
// zlecenia ani zlecenie bez śladu w rozmowie.
//
// Profil ma byt trwały: `ProfileId` z żądania odkłada się w kolumnie
// `zlecenie_asystenta.profil_kod`, bo niesie warstwę promptu rozstrzygającą,
// czy model sięgnie po narzędzia platformy, czy odpisze samym tekstem. Wejście
// tej warstwy do tury opisuje `adapter_modul_asystent_profil.go`.
//
// Syntezy mowy to nie dotyczy: `Speak` z żądania i `SpeechRef` w odpowiedzi
// zostają niespełnione, a `MessageId` puste — rdzeń nie ma silnika syntezy mowy
// ani szyny wiadomości okna w tym module. Profil niesie nastawę głosu
// (`glos_syntezy`), lecz nastawa nie jest wykonaniem.
func (a *adapterAsystenta) PolecenieGlosowe(ctx context.Context,
	z shared.AssistantVoiceCommandRequest) (shared.AssistantVoiceCommandResponse, error) {

	if z.WindowId == "" {
		return shared.AssistantVoiceCommandResponse{}, bladWskazaniaAsystenta("żądanie bez okna asystenta")
	}

	transkrypcja := ""
	if z.Transcript != nil {
		transkrypcja = *z.Transcript
	}
	nagranie := z.AudioRef
	brakNagrania := nagranie == nil || *nagranie == ""
	if transkrypcja == "" && brakNagrania {
		return shared.AssistantVoiceCommandResponse{}, bladWskazaniaAsystenta(
			"żądanie bez transkrypcji i bez odnośnika do nagrania")
	}

	// Przez silnik idzie samo nagranie. Tekst przysłany przez Operatora ma
	// pierwszeństwo i nie jest przepuszczany ponownie: `Transcript` niesie
	// wypowiedź już przez niego poprawioną, więc drugie rozpoznanie cofnęłoby
	// tę poprawkę.
	if transkrypcja == "" && a.mowa != nil {
		rozpoznana, err := a.rozpoznajNagranie(ctx, *nagranie)
		if err != nil {
			return shared.AssistantVoiceCommandResponse{}, err
		}
		transkrypcja = rozpoznana
	}

	// Droga zlecenia wynika z tego, co Operator naprawdę przysłał: nagranie
	// znaczy `voice`, sama poprawiona transkrypcja bez nagrania — `text`.
	droga := shared.AssistantOriginText
	if !brakNagrania {
		droga = shared.AssistantOriginVoice
	}

	// Profil odkłada się przy zleceniu. Wykonawca biegnie gorutyną długo po
	// odpowiedzi na tę komendę i warunki tury czyta z bazy, więc wskazanie
	// trzymane tylko w pamięci nie doczekałoby tury. Kod nieznany kończy się
	// nazwaną odmową jeszcze przed zapisem, bo klucz obcy kolumny zamieniłby go
	// w usterkę bez powodu.
	profil, err := a.wskazanyProfil(ctx, z.ProfileId)
	if err != nil {
		return shared.AssistantVoiceCommandResponse{}, err
	}

	zlecenie := dane.ZlecenieAsystenta{
		Kod:       nowyIdentyfikator(przedrostekZleceniaAsystenta),
		OknoKod:   z.WindowId,
		Stan:      string(shared.AssistantActionStatusQueued),
		Droga:     droga,
		ProfilKod: profil,
	}
	wpis := dane.WpisDziennikaAsystenta{
		Kod:              nowyIdentyfikator(przedrostekWpisuAsystenta),
		Rodzaj:           string(shared.AssistantActivityKindCommand),
		Tresc:            transkrypcja,
		NagranieOdnosnik: nagranie,
	}

	zapisaneZlecenie, _, err := a.repozytorium.PrzyjmijPolecenie(ctx, zlecenie, wpis)
	if err != nil {
		return shared.AssistantVoiceCommandResponse{}, bladAsystenta(err)
	}

	// Podejmij zlecenie z kolejki: świeżo założone `queued` ma zostać wykonane,
	// a nie czekać na ręczny `assistant.action.status` z okna. Bieg jest
	// asynchroniczny — komenda potwierdza przyjęcie od razu, a tura modelu trwa
	// dłużej niż wykonanie komendy.
	//
	// Z wpiętym silnikiem mowy treść zlecenia jest złożona, zanim zlecenie
	// powstanie. Bez silnika zlecenie z samym nagraniem zostaje w kolejce i
	// czeka na tekst przysłany wprost.
	a.podejmij(zapisaneZlecenie.Kod)

	return shared.AssistantVoiceCommandResponse{
		Transcript: transkrypcja,
		Action:     zlozZlecenieAsystenta(zapisaneZlecenie),
	}, nil
}

// zlozZlecenieAsystenta składa zlecenie kontraktu z wiersza repozytorium.
// Wspólne dla `PolecenieGlosowe` i `StanCzynnosci` drugiego agenta — obie
// komendy oddają ten sam kształt `AssistantAction` z tego samego wiersza.
func zlozZlecenieAsystenta(z dane.ZlecenieAsystenta) shared.AssistantAction {
	return shared.AssistantAction{
		Id:          z.Kod,
		WindowId:    z.OknoKod,
		Title:       z.Tytul,
		Status:      shared.AssistantActionStatus(z.Stan),
		Origin:      shared.AssistantOrigin(z.Droga),
		CurrentStep: wskaznikMalej(z.EtapBiezacy),
		TotalSteps:  wskaznikMalej(z.LiczbaEtapow),
		Priority:    wskaznikMalej(z.Priorytet),
		Result:      z.Wynik,
		CreatedAt:   z.Utworzono,
		UpdatedAt:   z.Zaktualizowano,
	}
}

// bladAsystenta znakuje usterkę kodem kontraktu, żeby okno modułu pokazało
// powód, a nie samo „nie udało się". Błąd, któremu kod już nadano, przechodzi
// bez zmiany; dopiero usterka bez kodu staje się usterką wewnętrzną rdzenia.
func bladAsystenta(err error) error {
	if err == nil {
		return nil
	}
	return protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeInternalError, err))
}

// bladWskazaniaAsystenta nazywa brak danych w żądaniu — błąd Operatora, nie rdzenia.
func bladWskazaniaAsystenta(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed, "moduł Assistant: "+powod))
}
