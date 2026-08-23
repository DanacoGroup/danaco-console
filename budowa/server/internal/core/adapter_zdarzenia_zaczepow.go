// Odpowiedzialność pliku: odbiór zdarzeń wykonawczych kanału głównego —
// zdarzeń zaczepów i zamknięć tur — oraz rozstrzyganie stanu odpowiedzi
// ze zdarzenia, a nie z treści wypowiedzi modelu.
//
// Dziennik zdarzeń dostaje każde zdarzenie zaczepu wraz z surową kopertą jako
// dowodem, tabela zamknięć dostaje zdarzenie `result` wraz ze stanem, który
// z niego wyprowadzono. Zapis następuje po zdarzeniu i nie steruje przebiegiem
// tury.
package core

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"sync"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/injection"
	"danacoconsole/shared"
)

// zdarzeniaWykonawcze przyjmuje zdarzenia wykonawcze tury i prowadzi ich ślad.
type zdarzeniaWykonawcze struct {
	zycie    context.Context
	repo     dane.RepozytoriumZdarzenWykonawczych
	dziennik *log.Logger

	// diagnostyka dopina się po montażu rejestru kanałów — rejestr powstaje
	// przed modułami — stąd zamek zamiast pola ustalanego w konstruktorze.
	mu          sync.RWMutex
	diagnostyka *adapterDiagnostyki
}

// nowyOdbiorZdarzenWykonawczych składa odbiornik nad repozytorium zdarzeń.
// Zestaw pusty daje odbiornik pracujący bez trwałości — tura biegnie, znika
// wyłącznie ślad, a dziennik procesu mówi o tym raz przy montażu.
func nowyOdbiorZdarzenWykonawczych(zycie context.Context, repozytoria *dane.Zestaw, dziennik *log.Logger) *zdarzeniaWykonawcze {
	if zycie == nil {
		zycie = context.Background()
	}
	odbior := &zdarzeniaWykonawcze{zycie: zycie, dziennik: dziennik}
	if repozytoria != nil {
		odbior.repo = repozytoria.ZdarzeniaWykonawcze
	}
	if odbior.repo == nil && dziennik != nil {
		dziennik.Printf("zdarzenia wykonawcze: brak trwałości — ślad zaczepów i zamknięć nie będzie zapisywany")
	}
	return odbior
}

// PodepnijDiagnostyke wpina odbiorcę odmów, dzięki czemu zdarzenia zaczepów
// trafiają do diagnostyki zamiast ginąć bez śladu.
func (o *zdarzeniaWykonawcze) PodepnijDiagnostyke(d *adapterDiagnostyki) {
	if o == nil {
		return
	}
	o.mu.Lock()
	o.diagnostyka = d
	o.mu.Unlock()
}

// HaczykZaczepow oddaje haczyk dla jednego okna i jednej wiadomości — kanał
// woła go przy każdym zdarzeniu zaczepu na strumieniu. Odbiornik pusty daje
// haczyk pusty, a kanał pusty haczyk pomija.
func (o *zdarzeniaWykonawcze) HaczykZaczepow(okno, wiadomosc string) func(injection.ZdarzenieZaczepu) {
	if o == nil {
		return nil
	}
	return func(e injection.ZdarzenieZaczepu) { o.zanotujZaczep(okno, wiadomosc, e) }
}

// rodzaje wierszy dziennika zdarzeń — słownik kolumny `dziennik_zdarzen.rodzaj`.
const (
	rodzajZaczepStart     = "zaczep_start"
	rodzajZaczepOdpowiedz = "zaczep_odpowiedz"
)

// zanotujZaczep zapisuje zdarzenie zaczepu w dzienniku zdarzeń i przekazuje
// je diagnostyce. Niepowodzenie zapisu nie przerywa tury — zostaje
// wpis w dzienniku procesu, żeby strata śladu nie była cicha.
func (o *zdarzeniaWykonawcze) zanotujZaczep(okno, wiadomosc string, e injection.ZdarzenieZaczepu) {
	if o == nil {
		return
	}
	if o.repo != nil {
		rodzaj := rodzajZaczepStart
		if e.Odpowiedz() {
			rodzaj = rodzajZaczepOdpowiedz
		}
		var kodWyjscia *int64
		if e.KodWyjscia != nil {
			kod := int64(*e.KodWyjscia)
			kodWyjscia = &kod
		}
		wiersz := dane.ZdarzenieZaczepuWiersz{
			Chwila: time.Now().UnixMilli(), OknoKod: okno, WiadomoscKod: wiadomosc,
			Rodzaj: rodzaj, ZaczepID: e.IdZaczepu, Zaczep: e.Nazwa,
			Zdarzenie: e.Zdarzenie, Wynik: e.Wynik, KodWyjscia: kodWyjscia,
			Tresc: trescZaczepu(e), Ladunek: string(e.Surowe), SesjaCLI: e.IdSesjiCLI,
		}
		if err := o.repo.ZapiszZaczep(o.zycie, wiersz); err != nil && o.dziennik != nil {
			o.dziennik.Printf("dziennik zdarzeń: zapis zaczepu %q nie powiódł się: %v", e.Nazwa, err)
		}
	}
	o.przekazDiagnostyce(okno, e)
}

// ZanotujZamkniecie utrwala zamknięcie tury razem ze stanem wiadomości
// wyprowadzonym ze zdarzenia.
func (o *zdarzeniaWykonawcze) ZanotujZamkniecie(okno, wiadomosc string, z *zamkniecieTury, stan shared.MessageStatus) {
	if o == nil || o.repo == nil || z == nil {
		return
	}
	wiersz := dane.ZamkniecieTuryWiersz{
		Chwila: time.Now().UnixMilli(), OknoKod: okno, WiadomoscKod: wiadomosc,
		Podtyp: z.Podtyp, Blad: z.Blad, StanNadany: shared.WartosciBazyMessageStatus[stan],
		KosztUSD: z.Koszt, Tury: z.Tury, CzasMs: z.CzasMs, Konto: z.Konto,
		SesjaCLI: z.IdRozmowyCLI, TypyZdarzen: strings.Join(z.TypyZdarzen, ","),
		LinieNierozpoznane: z.LinieNierozpoznane,
	}
	if err := o.repo.ZapiszZamkniecie(o.zycie, wiersz); err != nil && o.dziennik != nil {
		o.dziennik.Printf("zamknięcie tury: zapis dla wiadomości %q nie powiódł się: %v", wiadomosc, err)
	}
}

// ZanotujPrzelaczenie utrwala wykonane przełączenie kanału. Zapis następuje
// po przełączeniu; w samym strumieniu przełączenie widać po fragmencie
// metadanych konta i po prowenancji drugiego wywołania.
func (o *zdarzeniaWykonawcze) ZanotujPrzelaczenie(okno, wiadomosc, zKanalu, naKanal, powod string) {
	if o == nil || o.repo == nil {
		return
	}
	wiersz := dane.PrzelaczenieKanaluWiersz{
		Chwila: time.Now().UnixMilli(), OknoKod: okno, WiadomoscKod: wiadomosc,
		ZKanalu: zKanalu, NaKanal: naKanal, Powod: powod,
	}
	if err := o.repo.ZapiszPrzelaczenie(o.zycie, wiersz); err != nil && o.dziennik != nil {
		o.dziennik.Printf("przełączenie kanału: zapis %s→%s nie powiódł się: %v", zKanalu, naKanal, err)
	}
}

// trescZaczepu skleja wyjście polecenia zaczepu do kolumny `tresc`. Surowa
// koperta i tak jedzie w `ladunek` — to jest skrót do czytania, nie dowód.
func trescZaczepu(e injection.ZdarzenieZaczepu) string {
	tresc := strings.TrimSpace(e.Wyjscie)
	if blad := strings.TrimSpace(e.BladWyjscia); blad != "" {
		if tresc != "" {
			tresc += "; "
		}
		tresc += "stderr: " + blad
	}
	const limit = 2000
	if len(tresc) > limit {
		return tresc[:limit] + "…"
	}
	return tresc
}

// zamkniecieTury jest podsumowaniem zdarzenia `result` odczytanym z ostatniego
// fragmentu strumienia. Pola wiąże z injection.ZakonczenieTury zgodność tagów
// JSON, bez importu tamtego typu.
type zamkniecieTury struct {
	IdRozmowyCLI       string   `json:"cliSessionId"`
	Podtyp             string   `json:"subtype"`
	Blad               bool     `json:"isError"`
	Koszt              float64  `json:"costUsd"`
	Tury               int      `json:"turns"`
	CzasMs             int64    `json:"durationMs"`
	Konto              string   `json:"account"`
	TypyZdarzen        []string `json:"eventTypes"`
	LinieNierozpoznane int      `json:"unparsedLines"`
}

// zamkniecieZFragmentu wyjmuje podsumowanie tury z danych fragmentu. Fragment
// bez podtypu nie jest zamknięciem — prowenancja i błędy również niosą dane,
// a zdarzenie `result` zawsze ma pole `subtype`.
func zamkniecieZFragmentu(dane json.RawMessage) (zamkniecieTury, bool) {
	if len(dane) == 0 {
		return zamkniecieTury{}, false
	}
	var z zamkniecieTury
	if err := json.Unmarshal(dane, &z); err != nil || z.Podtyp == "" {
		return zamkniecieTury{}, false
	}
	return z, true
}

// stanOdpowiedziZeZdarzen rozstrzyga stan wiadomości po turze na podstawie
// zdarzeń wykonawczych.
//
// Kolejność warunków odpowiada hierarchii wiarygodności: przerwanie tury,
// potem błąd toru (kanał nie dowiózł strumienia), potem zdarzenie `result` —
// jego pole `is_error` rozstrzyga niezależnie od podtypu, ponieważ podtyp
// `success` występuje także przy `is_error: true`. Tura bez zamknięcia — na
// kanale, który zamknięcia nie nadaje (echo, api) — kończy się stanem
// ukończonym.
func stanOdpowiedziZeZdarzen(kontekst context.Context, err error, z *zamkniecieTury) shared.MessageStatus {
	if kontekst.Err() != nil {
		return shared.MessageStatusStopped
	}
	if err != nil {
		return shared.MessageStatusError
	}
	if z != nil && z.Blad {
		return shared.MessageStatusError
	}
	return shared.MessageStatusComplete
}
