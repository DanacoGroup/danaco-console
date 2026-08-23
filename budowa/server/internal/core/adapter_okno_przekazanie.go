// Odpowiedzialność pliku: obszar window.* — komenda `window.handoff`, która
// utrwala przekazanie zlecenia od okna koordynatora do okna wykonawcy
// Port `PrzekazanieOkna` i komenda `window.action` leżą w
// `adapter_okno_akcja.go` i `adapter_okno_przekazanie_uchwyty.go`; ten plik
// deklaruje wyłącznie typ adaptera, jego konstruktor i wiązania zależności.
//
// Więź koordynator–wykonawca (kolumna `okno_komunikacji.okno_koordynatora_id`)
// bez utrwalenia tutaj żyłaby wyłącznie w pamięci przeglądarki i ginęła z jej
// zamknięciem, a Mission Control nie miałby czego pokazać po ponownym
// uruchomieniu rdzenia.
//
// Komenda robi trzy rzeczy naraz — wszystkie albo żadną:
//  1. utrwala więź koordynator–wykonawca,
//  2. zapisuje zlecenie wraz z kompletem kontekstu,
//  3. zakłada pozycję kolejki, której identyfikator wraca w `QueueItemId`.
//
// Warstwa danych (`dane/przekazanie_okna*.go`) przyjmuje `PozycjaKolejkiID` jako
// identyfikator gotowy i sama pozycji nie zakłada. Zakładanie należy więc do
// tego adaptera i idzie przez jedyny silnik kolejek — ten sam `adapterKolejek`,
// którym pracuje pętla sesyjna i moduł Automations. Wzorem
// `adapterAutomatyk.ZKolejkami` (`adapter_modul_automations.go`) konstruktor
// bierze wyłącznie repozytorium obszaru, a silnik kolejek dochodzi osobnym
// wiązaniem po złożeniu grafu zależności w `montaz_moduly.go`. Bez podpiętego
// adaptera kolejek `Przekaz` odmawia wprost, zamiast meldować wykonanie.
//
// Repozytoria okien, sesji, modułów i kanałów stoją tu, bo `window.handoff`
// przyjmuje identyfikatory zewnętrzne (tekstowe) okien i sesji, a warstwa danych
// obszaru window.* oraz silnik kolejek pracują na kluczach wewnętrznych
// (`int64`) tabel `okno_komunikacji` i `sesja`. Rozwiązanie identyfikatora na
// wiersz — i odmowa `not_found`, gdy okna nie ma — jest obowiązkiem tego
// adaptera, bo warstwa danych przyjmuje klucze już rozwiązane. Moduły i kanały
// modelu służą wyłącznie złożeniu odpowiedzi: kontrakt `Window` niesie kody
// modułu i kanału, a wiersz okna niesie klucze obce do nich. Wszystkie te
// zależności powstają w `montaz_moduly.go`, dlatego adapter przyjmuje je gotowe
// przez wiązania tego pliku.
package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// Przedrostki i stałe wartości bytów zakładanych przez tę komendę. Zlecenie
// przekazania samo nie ma identyfikatora tekstowego kontraktu (kontrakt oddaje
// wyłącznie `Window` i `QueueItemId`), więc jedyny przedrostek dotyczy rodzaju
// kolejki, którą komenda zakłada dla każdego przekazania z osobna.
const (
	// rodzajKolejkiPrzekazania to `multitasking`, a nie własna wartość.
	//
	// Schemat dopuszcza dokładnie dwa rodzaje kolejki — `sesyjna`
	// i `multitasking` — bo jeden silnik kolejek obsługuje pętlę sesyjną
	// i MultitaskingAI. Przekazanie zlecenia z okna koordynatora do okna
	// wykonawcy jest pętlą MultitaskingAI, więc mieści się w rodzaju już
	// istniejącym; własny rodzaj byłby trzecim znaczeniem tego samego pojęcia
	// i padłby na warunku CHECK kolumny.
	rodzajKolejkiPrzekazania = "multitasking"
)

// adapterPrzekazaniaOkna wypełnia port PrzekazanieOkna.
// Zależności poza repozytorium własnego obszaru są bytami montażu: silnik
// kolejek (jedyny wykonawca pozycji), repozytoria okien i sesji (rozwiązanie
// identyfikatorów zewnętrznych na wewnętrzne) oraz słowniki modułów i kanałów
// modelu (złożenie okna wyniku).
type adapterPrzekazaniaOkna struct {
	repozytorium dane.RepozytoriumPrzekazan
	kolejki      *adapterKolejek
	okna         dane.RepozytoriumOkien
	sesje        dane.RepozytoriumSesji
	moduly       dane.RepozytoriumModulow
	kanaly       dane.RepozytoriumKanalow
	akcje        dane.RepozytoriumAkcji
}

// nowyAdapterPrzekazaniaOkna wiąże port z repozytorium obszaru window.*.
// Pozostałe zależności dochodzą osobnymi wiązaniami po złożeniu grafu w
// `montaz_moduly.go`: konstruktor nie odmawia z powodu bytów, których jeszcze
// nie ma w kolejności montażu.
func nowyAdapterPrzekazaniaOkna(repozytorium dane.RepozytoriumPrzekazan) *adapterPrzekazaniaOkna {
	return &adapterPrzekazaniaOkna{repozytorium: repozytorium}
}

// ZKolejkami podpina jedyny silnik kolejek. Bez niego `Przekaz` odmawia wprost —
// wzór: komentarz przy `adapterAutomatyk.ZKolejkami` w
// `adapter_modul_automations.go`.
func (a *adapterPrzekazaniaOkna) ZKolejkami(kolejki *adapterKolejek) *adapterPrzekazaniaOkna {
	a.kolejki = kolejki
	return a
}

// ZBytamiOkien podpina repozytoria potrzebne do rozwiązania identyfikatorów
// zewnętrznych żądania (okna, sesja) i do złożenia okna wyniku (moduły,
// kanały modelu). Cztery byty naraz, bo wszystkie cztery powstają w tym samym
// miejscu montażu i żaden z osobna nie czyni komendy wykonalną.
func (a *adapterPrzekazaniaOkna) ZBytamiOkien(okna dane.RepozytoriumOkien, sesje dane.RepozytoriumSesji,
	moduly dane.RepozytoriumModulow, kanaly dane.RepozytoriumKanalow) *adapterPrzekazaniaOkna {

	a.okna = okna
	a.sesje = sesje
	a.moduly = moduly
	a.kanaly = kanaly
	return a
}

// ZKatalogiemAkcji podpina katalog akcji (tabela `akcja`). Bez niego
// `window.action` nie odróżniłaby akcji istniejącej od zmyślonej i kwitowałaby
// każdą powodzeniem.
// Katalog jest odczytywany wprost z repozytorium, a nie z `RejestrAkcji`,
// bo rozstrzygnięcie „czy ta akcja istnieje" ma widzieć stan bieżący wiersza,
// nie kopię z chwili startu rdzenia.
func (a *adapterPrzekazaniaOkna) ZKatalogiemAkcji(akcje dane.RepozytoriumAkcji) *adapterPrzekazaniaOkna {
	a.akcje = akcje
	return a
}

// Przekaz wykonuje `window.handoff`. Trzy zapisy naraz: pozycja kolejki, więź
// koordynator–wykonawca, zlecenie z kontekstem. Kolejność jest celowa: pozycja
// kolejki idzie pierwsza, bo jest jedynym z trzech zapisów, który silnik kolejek
// umie cofnąć samodzielnie (kolejka porzucona bez zlecenia jest stanem
// nieszkodliwym); więź i zlecenie idą po niej, gdy wiadomo już, że jest czym
// wykonać.
//
// Wspólnej transakcji SQL między repozytoriami nie ma — warstwa danych obszaru
// window.* i silnik kolejek to dwa oddzielne repozytoria. Odmowa wczesna, przed
// pierwszym zapisem, pokrywa najczęstszy przypadek: okno albo sesja, których nie
// ma. Usterka po pierwszym zapisie zostawia założoną pozycję kolejki bez
// zlecenia, co Queue Manager pokazuje jako pozycję do ręcznego domknięcia, a nie
// jako ciche zaginięcie zlecenia.
func (a *adapterPrzekazaniaOkna) Przekaz(ctx context.Context,
	z shared.WindowHandoffRequest) (shared.WindowHandoffResponse, error) {

	if a.kolejki == nil {
		return shared.WindowHandoffResponse{}, errBrakSilnikaKolejekPrzekazania
	}
	zrodlowe, docelowe, sesja, err := a.rozwiazBytyZadania(ctx, z)
	if err != nil {
		return shared.WindowHandoffResponse{}, err
	}
	pozycjaID, err := a.zalozPozycjeKolejki(ctx, sesja, zrodlowe, docelowe, z)
	if err != nil {
		return shared.WindowHandoffResponse{}, bladPrzekazania(err)
	}
	if err := a.repozytorium.UstawKoordynatora(ctx, z.ToWindowId, z.FromWindowId); err != nil {
		return shared.WindowHandoffResponse{}, bladPrzekazania(err)
	}
	if err := a.zapiszZlecenie(ctx, sesja, zrodlowe, docelowe, pozycjaID, z); err != nil {
		return shared.WindowHandoffResponse{}, bladPrzekazania(err)
	}
	oknoWyniku, err := a.oknoWynikuKontraktu(ctx, docelowe, z)
	if err != nil {
		return shared.WindowHandoffResponse{}, bladPrzekazania(err)
	}
	return shared.WindowHandoffResponse{Window: oknoWyniku, QueueItemId: strconv.FormatInt(pozycjaID, 10)}, nil
}

// rozwiazBytyZadania odnajduje okno źródłowe, okno docelowe i sesję po ich
// identyfikatorach zewnętrznych. Okno, którego nie ma — po którejkolwiek
// stronie przekazania — kończy się odmową `not_found`, nie cichą zgodą: cicha
// zgoda potwierdzałaby więź, która w rzeczywistości nie powstała.
func (a *adapterPrzekazaniaOkna) rozwiazBytyZadania(ctx context.Context,
	z shared.WindowHandoffRequest) (dane.Okno, dane.Okno, dane.Sesja, error) {

	zrodlowe, err := a.okna.PoIdentyfikatorze(ctx, z.FromWindowId)
	if err != nil {
		return dane.Okno{}, dane.Okno{}, dane.Sesja{}, bladNieznanegoOknaPrzekazania(z.FromWindowId, err)
	}
	docelowe, err := a.okna.PoIdentyfikatorze(ctx, z.ToWindowId)
	if err != nil {
		return dane.Okno{}, dane.Okno{}, dane.Sesja{}, bladNieznanegoOknaPrzekazania(z.ToWindowId, err)
	}
	sesja, err := a.sesje.PoIdentyfikatorze(ctx, z.SessionId)
	if err != nil {
		return dane.Okno{}, dane.Okno{}, dane.Sesja{}, bladNieznanejSesjiPrzekazania(z.SessionId, err)
	}
	return zrodlowe, docelowe, sesja, nil
}

// zalozPozycjeKolejki zakłada kolejkę przekazania i jej jedyną pozycję przez
// silnik kolejek (`a.kolejki.repozytorium` — ten sam obiekt, którym jedzie
// `adapter_modul_automations_kolejka.go`). Kolejka nowa dla każdego przekazania,
// bo handoff nie ma pojęcia „kolejki tej pary okien" do ponownego użycia —
// wzorzec `queue.create` w `adapter_kolejki.go` też zakłada kolejkę za każdym
// wywołaniem.
func (a *adapterPrzekazaniaOkna) zalozPozycjeKolejki(ctx context.Context, sesja dane.Sesja,
	zrodlowe, docelowe dane.Okno, z shared.WindowHandoffRequest) (int64, error) {

	idSesji, idZrodlowe, idDocelowe := sesja.ID, zrodlowe.ID, docelowe.ID
	kolejkaID, err := a.kolejki.repozytorium.UtworzKolejke(ctx, dane.Kolejka{
		Nazwa: "przekazanie: " + z.Instruction, Rodzaj: rodzajKolejkiPrzekazania,
		SesjaID: &idSesji, OknoKoordynatoraID: &idZrodlowe, Stan: shared.QueueStatusIdle,
	})
	if err != nil {
		return 0, err
	}
	return a.kolejki.repozytorium.DodajPozycje(ctx, dane.Pozycja{
		KolejkaID: kolejkaID, OknoWykonawcyID: &idDocelowe, Tytul: z.Instruction,
		TrescZlecenia: &z.Instruction, Stan: stanPozycjiOczekuje,
	})
}

// zapiszZlecenie utrwala treść zlecenia wraz z kompletem kontekstu.
// `Bundle` idzie jako surowy zapis JSON — ten adapter go nie rozbiera, tak jak
// nakazuje komentarz `dane/przekazanie_okna.go`.
func (a *adapterPrzekazaniaOkna) zapiszZlecenie(ctx context.Context, sesja dane.Sesja,
	zrodlowe, docelowe dane.Okno, pozycjaID int64, z shared.WindowHandoffRequest) error {

	kontekst, err := kompletKontekstuJSON(z.Bundle)
	if err != nil {
		return bladWskazaniaPrzekazania("komplet kontekstu jest nieczytelny: " + err.Error())
	}
	_, err = a.repozytorium.ZapiszZlecenie(ctx, dane.ZleceniePrzekazania{
		SesjaID: sesja.ID, OknoZrodloweID: zrodlowe.ID, OknoDoceloweID: docelowe.ID,
		PozycjaKolejkiID: pozycjaID, Polecenie: z.Instruction, KompletKontekstu: kontekst,
	})
	return err
}

// kompletKontekstuJSON przekłada `ContextBundle` opcjonalny na surowy zapis
// JSON opcjonalny. Żądanie bez kompletu — samo polecenie — zapisuje się bez
// tej kolumny; to stan poprawny, nie brak danych.
func kompletKontekstuJSON(pakiet *shared.ContextBundle) (*string, error) {
	if pakiet == nil {
		return nil, nil
	}
	surowy, err := json.Marshal(pakiet)
	if err != nil {
		return nil, err
	}
	tekst := string(surowy)
	return &tekst, nil
}

// oknoWynikuKontraktu składa okno docelowe kontraktu po przyjęciu zlecenia.
// Korzysta z tego samego przekładu co odczyt utrwalony (`oknoWierszaKontraktu`
// w `przeklad_nawigacja.go`); dokłada wyłącznie identyfikator sesji zewnętrzny
// (już znany z żądania — drugi odczyt sesji byłby zapytaniem po to samo) i
// świeżą więź koordynatora, którą ta komenda właśnie zapisała.
func (a *adapterPrzekazaniaOkna) oknoWynikuKontraktu(ctx context.Context, docelowe dane.Okno,
	z shared.WindowHandoffRequest) (shared.Window, error) {

	moduly, err := a.moduly.Lista(ctx)
	if err != nil {
		return shared.Window{}, err
	}
	kodModulu, err := kodSlownika(docelowe.ModulID, moduly, func(m dane.Modul) (int64, string) { return m.ID, m.Kod })
	if err != nil {
		return shared.Window{}, err
	}
	kanaly, err := a.kanaly.Lista(ctx, false)
	if err != nil {
		return shared.Window{}, err
	}
	kodKanalu, err := kodSlownika(docelowe.KanalModeluID, kanaly, func(k dane.Kanal) (int64, string) { return k.ID, k.Kod })
	if err != nil {
		return shared.Window{}, err
	}
	okno := oknoWierszaKontraktu(docelowe, kodModulu, kodKanalu)
	okno.SessionId = z.SessionId
	koordynator := z.FromWindowId
	okno.CoordinatorWindowId = &koordynator
	return okno, nil
}

// kodSlownika odnajduje kod bytu o wskazanym kluczu wewnętrznym w już
// odczytanym wykazie słownika (moduły albo kanały modelu). Typ ogólny, bo oba
// słowniki mają ten sam kształt wyszukania — jeden przebieg wykazu zamiast
// dwóch niemal identycznych funkcji; wykaz kanałów czyta się bez filtra
// aktywności, bo okno mogło zostać założone na kanale, który od tamtej pory
// wyłączono.
func kodSlownika[T any](id int64, wykaz []T, klucz func(T) (int64, string)) (string, error) {
	for _, element := range wykaz {
		elementID, elementKod := klucz(element)
		if elementID == id {
			return elementKod, nil
		}
	}
	return "", fmt.Errorf("dane: byt słownika %d nie istnieje", id)
}

// bladPrzekazania znakuje usterkę kodem kontraktu. Błąd, któremu kod już
// nadano — odmowa wskazania, brak bytu, brak wykonawcy — przechodzi tędy bez
// zmiany kodu; dopiero usterka bez kodu staje się usterką wewnętrzną rdzenia.
func bladPrzekazania(err error) error {
	if err == nil {
		return nil
	}
	return protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeInternalError, err))
}

// bladWskazaniaPrzekazania nazywa brak albo uszkodzenie danych w żądaniu — to
// błąd żądania, nie rdzenia.
func bladWskazaniaPrzekazania(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"window.handoff: "+powod))
}

// bladNieznanegoOknaPrzekazania odróżnia „okna nie ma" od „odczyt się nie
// powiódł". Okno źródłowe albo docelowe, którego nie ma, kończy się odmową
// `not_found` — nigdy cichą zgodą na więź, która nie powstała.
func bladNieznanegoOknaPrzekazania(idOkna string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"window.handoff: okno nie istnieje: "+idOkna))
	}
	return bladPrzekazania(err)
}

// bladNieznanejSesjiPrzekazania odróżnia „sesji nie ma" od „odczyt się nie
// powiódł" — tak samo jak przy oknie, sesja nieznana jest błędem żądania.
func bladNieznanejSesjiPrzekazania(idSesji string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"window.handoff: sesja nie istnieje: "+idSesji))
	}
	return bladPrzekazania(err)
}

// errBrakSilnikaKolejekPrzekazania nazywa brak wykonawcy pozycji kolejki.
// Kontrakt nie ma kodu „domena niewpięta" (tak samo jak przy Automations —
// `adapter_modul_automations_pozycje.go`), więc odmowa idzie tym samym kodem
// `channel_unavailable`: ponawialna, wpięcie silnika czyni żądanie wykonalnym
// bez zmiany treści.
var errBrakSilnikaKolejekPrzekazania = protocol.JakoError(protocol.NowyBlad(
	shared.ErrorCodeChannelUnavailable,
	"window.handoff: silnik kolejek nie jest wpięty — przekazania nie ma kto wykonać"))
