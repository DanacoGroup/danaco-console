// Wypełnia obszar window.* komendą `window.handoff`, która utrwala
// przekazanie zlecenia od okna koordynatora do okna wykonawcy wraz z pozycją
// kolejki i zapisanym kontekstem; deklaruje typ adaptera, konstruktor
// i wiązania zależności.
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

// Stałe wartości bytów zakładanych przez tę komendę. Zlecenie przekazania nie
// ma własnego identyfikatora tekstowego kontraktu, więc jedyna stała dotyczy
// rodzaju kolejki zakładanej dla każdego przekazania z osobna.
const (
	// rodzajKolejkiPrzekazania to `multitasking`, a nie własna wartość: schemat
	// dopuszcza tylko `sesyjna` i `multitasking`, a przekazanie zlecenia jest
	// pętlą MultitaskingAI.
	rodzajKolejkiPrzekazania = "multitasking"
)

// adapterPrzekazaniaOkna wypełnia port PrzekazanieOkna. Zależności poza
// repozytorium własnego obszaru — silnik kolejek, repozytoria okien i sesji,
// słowniki modułów i kanałów — są bytami montażu.
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

// ZBytamiOkien podpina repozytoria do rozwiązania identyfikatorów żądania
// (okna, sesja) i do złożenia okna wyniku (moduły, kanały). Cztery naraz,
// bo powstają w tym samym miejscu montażu.
func (a *adapterPrzekazaniaOkna) ZBytamiOkien(okna dane.RepozytoriumOkien, sesje dane.RepozytoriumSesji,
	moduly dane.RepozytoriumModulow, kanaly dane.RepozytoriumKanalow) *adapterPrzekazaniaOkna {

	a.okna = okna
	a.sesje = sesje
	a.moduly = moduly
	a.kanaly = kanaly
	return a
}

// ZKatalogiemAkcji podpina katalog akcji. Bez niego `window.action` nie
// odróżniłaby akcji istniejącej od zmyślonej. Katalog czyta się wprost
// z repozytorium, żeby widzieć stan bieżący wiersza.
func (a *adapterPrzekazaniaOkna) ZKatalogiemAkcji(akcje dane.RepozytoriumAkcji) *adapterPrzekazaniaOkna {
	a.akcje = akcje
	return a
}

// Przekaz wykonuje `window.handoff`: zakłada pozycję kolejki, utrwala więź
// koordynator–wykonawca i zapisuje zlecenie z kontekstem. Kolejność zapisów
// jest celowa, a wspólnej transakcji SQL między repozytoriami nie ma.
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

// rozwiazBytyZadania odnajduje okno źródłowe, okno docelowe i sesję po
// identyfikatorach zewnętrznych. Brak okna po którejkolwiek stronie kończy
// się odmową `not_found`, nie cichą zgodą.
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
// silnik kolejek. Kolejka jest nowa dla każdego przekazania, bo handoff nie
// ma pojęcia kolejki tej pary okien do ponownego użycia.
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
// Korzysta z przekładu `oknoWierszaKontraktu` i dokłada identyfikator sesji
// oraz świeżą więź koordynatora zapisaną przez tę komendę.
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

// kodSlownika odnajduje kod bytu o wskazanym kluczu wewnętrznym w odczytanym
// wykazie słownika, modułów albo kanałów. Typ ogólny zastępuje dwie niemal
// identyczne funkcje jednym przebiegiem wykazu.
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
// Kontrakt nie ma kodu domena niewpięta, więc odmowa idzie kodem
// `channel_unavailable`, ponawialnym po wpięciu silnika.
var errBrakSilnikaKolejekPrzekazania = protocol.JakoError(protocol.NowyBlad(
	shared.ErrorCodeChannelUnavailable,
	"window.handoff: silnik kolejek nie jest wpięty — przekazania nie ma kto wykonać"))
