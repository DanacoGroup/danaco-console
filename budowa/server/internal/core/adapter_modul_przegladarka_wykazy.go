// Odpowiedzialność pliku: obsługa `browser.source.list` i `browser.note.list` —
// odczytowa strona szuflady źródeł i szuflady notatek. Obie komendy nie
// dokładają nowej wiedzy — oddają to, co `browser.source.add`
// i `browser.note.add` już zapisały w tabelach modułu.
//
// Okno nieznane to odmowa, okno puste to wynik. Wykaz pusty jest prawidłową
// odpowiedzią: okno przeglądania istnieje, tylko nic w nim jeszcze nie
// zebrano. Okno, którego moduł nigdy nie widział, dostaje `not_found` — gdyby
// oddać na nie pustą tablicę, literówka w identyfikatorze okna wyglądałaby
// dokładnie tak samo jak uczciwie pusta szuflada, a Operator szukałby braku
// danych zamiast braku okna. Znaczenie „znane" rozstrzyga `dane.OknoZnane`
// i jego nagłówek.
//
// Limit zerowy znaczy wykaz pełny, nie pusty. Kontrakt daje `limit` jako pole
// nieobowiązkowe; jego brak to „nie ograniczaj", nie „oddaj nic" — przekład
// robi `granicaWykazu` warstwy danych.
package core

import (
	"context"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// WykazZrodel oddaje źródła zebrane w oknie, od najnowszego, ewentualnie
// zawężone do oznaczonych jako kluczowe.
func (a *adapterPrzegladarki) WykazZrodel(ctx context.Context,
	z shared.BrowserSourceListRequest) (shared.BrowserSourceListResponse, error) {

	if err := a.upewnijSieOknoZnane(ctx, z.WindowId, "source.list"); err != nil {
		return shared.BrowserSourceListResponse{}, err
	}
	wiersze, err := a.repozytorium.Zrodla(ctx, dane.FiltrZrodelPrzegladania{
		Okno:          z.WindowId,
		TylkoKluczowe: wartoscLogiczna(z.KeyOnly),
		Limit:         wartoscLiczby(z.Limit),
		Zestaw:        wartoscTekstu(z.GroupId),
		Szukaj:        wartoscTekstu(z.Query),
	})
	if err != nil {
		return shared.BrowserSourceListResponse{}, bladPrzegladarki(err)
	}

	zrodla := make([]shared.BrowserSource, 0, len(wiersze))
	for _, wiersz := range wiersze {
		zrodla = append(zrodla, zrodloKontraktu(wiersz))
	}
	return shared.BrowserSourceListResponse{Sources: zrodla}, nil
}

// WykazNotatek oddaje notatki okna, od najświeższej, w całości albo zawężone
// do jednego źródła. Zawężenie do źródła, którego w oknie nie ma, daje wykaz
// pusty, nie odmowę: notatka bez źródła jest dozwolona, a kolumna źródła nie
// niesie więzu obcego, więc moduł nie ma czego sprawdzać poza samym oknem.
func (a *adapterPrzegladarki) WykazNotatek(ctx context.Context,
	z shared.BrowserNoteListRequest) (shared.BrowserNoteListResponse, error) {

	if err := a.upewnijSieOknoZnane(ctx, z.WindowId, "note.list"); err != nil {
		return shared.BrowserNoteListResponse{}, err
	}
	wiersze, err := a.repozytorium.Notatki(ctx, dane.FiltrNotatekPrzegladania{
		Okno:                z.WindowId,
		ZrodloID:            wartoscTekstu(z.SourceId),
		Limit:               wartoscLiczby(z.Limit),
		Watek:               wartoscTekstu(z.ThreadId),
		Klasyfikacja:        klasyfikacjaNotatki(z.Classification),
		Szukaj:              wartoscTekstu(z.Query),
		PrzypieteNaPoczatku: wartoscLogiczna(z.PinnedFirst),
	})
	if err != nil {
		return shared.BrowserNoteListResponse{}, bladPrzegladarki(err)
	}

	notatki := make([]shared.BrowserNote, 0, len(wiersze))
	for _, wiersz := range wiersze {
		notatki = append(notatki, notatkaKontraktu(wiersz))
	}
	return shared.BrowserNoteListResponse{Notes: notatki}, nil
}

// upewnijSieOknoZnane przepuszcza dalej wyłącznie okno, które moduł Browser
// naprawdę widział. Rozdziela trzy odmowy, bo to trzy różne winy: brak
// wskazania okna jest usterką żądania (`validation_failed`), okno nieznane —
// pomyłką Operatora (`not_found`), a niepowodzenie odczytu — usterką rdzenia.
func (a *adapterPrzegladarki) upewnijSieOknoZnane(ctx context.Context, okno, komenda string) error {
	if okno == "" {
		return bladWskazaniaPrzegladarki("komenda " + komenda + " bez okna")
	}
	znane, err := a.repozytorium.OknoZnane(ctx, okno)
	if err != nil {
		return bladPrzegladarki(err)
	}
	if !znane {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Browser: nie ma okna przeglądania o wskazaniu: "+okno))
	}
	return nil
}

// klasyfikacjaNotatki przekłada wyliczenie kontraktu na tekst kolumny. Brak
// wskazania znaczy „wszystkie rodzaje", nie „rodzaj pusty".
func klasyfikacjaNotatki(wskazanie *shared.BrowserNoteClassification) string {
	if wskazanie == nil {
		return ""
	}
	return string(*wskazanie)
}
