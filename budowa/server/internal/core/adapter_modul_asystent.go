// Moduł Assistant: typ adaptera, konstruktor i obsługa assistant.voice.command.
// Rdzeń rozpoznaje mowę, gdy przychodzi samo nagranie, i przepisuje je, zanim
// założy zlecenie.
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
// jedna: wspólne repozytorium modułu, rozłożone po stronie danych na trzy
// pliki wedle odpowiedzialności, ale niosące jeden typ.
type adapterAsystenta struct {
	repozytorium dane.RepozytoriumAsystenta

	kanaly   *models.Rejestr
	nadzorca *session.Nadzorca
	nadajnik Nadajnik
	// mowa rozpoznaje nagranie, gdy Operator nie dosłał tekstu, zależność opcjonalna widziana portem.
	mowa Mowa
	// mosty składają konfigurację MCP okna, którym model zlecenia dostaje sterowanie platformą.
	mosty *mostyOkna
	// zycie jest kontekstem rdzenia, nie połączenia — tura biegnie dłużej niż komenda i jej nie przerywa.
	zycie context.Context
	// biegi wiążą kod zlecenia z przerwaniem tury, sterowanie zmienia biegnącą turę, nie sam wiersz.
	biegi map[string]context.CancelFunc
	// odwolane niesie decyzje Operatora zapadłe, zanim wykonawca zdążył przestawić zlecenie na running.
	odwolane map[string]bool
	zamki    sync.Mutex
}

// nowyAdapterAsystenta wiąże adapter z repozytorium modułu. Wykonawca zleceń
// dochodzi metodami budującymi (`ZKanalami`, `ZSesjami`, `ZWyjsciem`).
func nowyAdapterAsystenta(repozytorium dane.RepozytoriumAsystenta) *adapterAsystenta {
	return &adapterAsystenta{repozytorium: repozytorium}
}

// ZMowa podpina silnik rozpoznawania mowy. Bez niego moduł przyjmuje nagranie bez transkrypcji tekstowej.
func (a *adapterAsystenta) ZMowa(m Mowa) *adapterAsystenta {
	// Montaż podaje tu wskaźnik, a wskaźnik nil schowany w interfejsie nie jest nil dla porównania.
	if m == nil {
		return a
	}
	if adapter, jest := m.(*adapterMowy); jest && adapter == nil {
		return a
	}
	a.mowa = m
	return a
}

// rozpoznajNagranie zamienia nagranie na tekst silnikiem mowy, wedle trzech wyników speech.transcribe.
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

// PolecenieGlosowe obsługuje assistant.voice.command, zakładając zlecenie i pierwszy wpis dziennika razem.
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

	// Przez silnik idzie samo nagranie — tekst Operatora ma pierwszeństwo i nie idzie przez nie ponownie.
	if transkrypcja == "" && a.mowa != nil {
		rozpoznana, err := a.rozpoznajNagranie(ctx, *nagranie)
		if err != nil {
			return shared.AssistantVoiceCommandResponse{}, err
		}
		transkrypcja = rozpoznana
	}

	// Droga zlecenia wynika z tego, co Operator przysłał: nagranie znaczy voice, sama transkrypcja — text.
	droga := shared.AssistantOriginText
	if !brakNagrania {
		droga = shared.AssistantOriginVoice
	}

	// Profil odkłada się przy zleceniu, bo wykonawca biegnie gorutyną długo po odpowiedzi na komendę.
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

	// Podejmij zlecenie z kolejki: świeżo założone queued ma zostać wykonane od razu, bez ręcznego wpisu.
	a.podejmij(ctx, zapisaneZlecenie.Kod)

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

// bladWskazaniaAsystenta nazywa brak danych w żądaniu — błąd Operatora, nie usterkę rdzenia platformy.
func bladWskazaniaAsystenta(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed, "moduł Assistant: "+powod))
}
