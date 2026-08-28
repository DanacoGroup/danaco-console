// Odpowiedzialność pliku: obsługa browser.source.list i browser.note.list — odczytowa strona szuflad źródeł i notatek, bez dokładania nowej wiedzy.
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

// WykazNotatek oddaje notatki okna, od najświeższej, w całości albo zawężone do jednego źródła; źródło spoza okna daje wykaz pusty, nie odmowę.
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

// upewnijSieOknoZnane przepuszcza dalej wyłącznie okno naprawdę widziane przez moduł Browser, rozdzielając trzy różne odmowy według winy.
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
