// Odpowiedzialność pliku: obsługa `browser.source.add` i `browser.note.add` — dwie
// szuflady modułu Browser zasilane z toku przeglądania, każda swoim wierszem repozytorium.
package core

import (
	"context"
	"fmt"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

const (
	przedrostekZrodlaPrzegladania  = "src-"
	przedrostekNotatkiPrzegladania = "note-"
)

// Zapis jest idempotentny wobec identyfikatora zewnętrznego: powtórzenie nadpisuje wiersz, nie dubluje.
func (a *adapterPrzegladarki) DodajZrodlo(ctx context.Context, z shared.BrowserSourceAddRequest) (shared.BrowserSourceAddResponse, error) {
	if brak := brakiZrodlaPrzegladania(z); brak != "" {
		return shared.BrowserSourceAddResponse{}, bladWskazaniaPrzegladarki(brak)
	}

	zrodlo := dane.ZrodloPrzegladania{
		Kod:                 nowyIdentyfikator(przedrostekZrodlaPrzegladania),
		Okno:                z.WindowId,
		Url:                 z.Url,
		Tytul:               z.Title,
		MigawkaZewnetrznaID: z.SnapshotId,
		Kluczowe:            wartoscLogiczna(z.Key),
		Grupa:               z.GroupId,
	}
	zapisane, err := a.repozytorium.ZapiszZrodlo(ctx, zrodlo)
	if err != nil {
		return shared.BrowserSourceAddResponse{}, bladPrzegladarki(err)
	}
	return shared.BrowserSourceAddResponse{Source: zrodloKontraktu(zapisane)}, nil
}

func (a *adapterPrzegladarki) DodajNotatke(ctx context.Context, z shared.BrowserNoteAddRequest) (shared.BrowserNoteAddResponse, error) {
	if brak := brakiNotatkiPrzegladania(z); brak != "" {
		return shared.BrowserNoteAddResponse{}, bladWskazaniaPrzegladarki(brak)
	}

	notatka := dane.NotatkaPrzegladania{
		Kod:       nowyIdentyfikator(przedrostekNotatkiPrzegladania),
		Okno:      z.WindowId,
		ZrodloID:  z.SourceId,
		Tresc:     z.Content,
		Cytat:     z.Quote,
		Watek:     z.ThreadId,
		Przypieta: wartoscLogiczna(z.Pinned),
	}
	// Klasyfikacja przychodzi wyliczeniem kontraktu, a kolumna trzyma tekst; przekład jest tu.
	if z.Classification != nil {
		klasyfikacja := string(*z.Classification)
		notatka.Klasyfikacja = &klasyfikacja
	}
	zapisana, err := a.repozytorium.ZapiszNotatke(ctx, notatka)
	if err != nil {
		return shared.BrowserNoteAddResponse{}, fmt.Errorf("core: nie można dodać notatki przeglądania: %w", err)
	}
	return shared.BrowserNoteAddResponse{Note: notatkaKontraktu(zapisana)}, nil
}

func brakiZrodlaPrzegladania(z shared.BrowserSourceAddRequest) string {
	var puste []string
	if strings.TrimSpace(z.WindowId) == "" {
		puste = append(puste, "windowId")
	}
	if strings.TrimSpace(z.Url) == "" {
		puste = append(puste, "url")
	}
	if len(puste) == 0 {
		return ""
	}
	return "źródło przeglądania z pustymi polami: " + strings.Join(puste, ", ")
}

func brakiNotatkiPrzegladania(z shared.BrowserNoteAddRequest) string {
	var puste []string
	if strings.TrimSpace(z.WindowId) == "" {
		puste = append(puste, "windowId")
	}
	if strings.TrimSpace(z.Content) == "" {
		puste = append(puste, "content")
	}
	if len(puste) == 0 {
		return ""
	}
	return "notatka przeglądania z pustymi polami: " + strings.Join(puste, ", ")
}

func zrodloKontraktu(z dane.ZrodloPrzegladania) shared.BrowserSource {
	zrodlo := shared.BrowserSource{
		Id:         z.Kod,
		WindowId:   z.Okno,
		Url:        z.Url,
		Title:      z.Tytul,
		SnapshotId: z.MigawkaZewnetrznaID,
		CreatedAt:  chwilaBazy(z.Utworzono),
		GroupId:    z.Grupa,
	}
	if z.Kluczowe {
		kluczowe := true
		zrodlo.Key = &kluczowe
	}
	return zrodlo
}

func notatkaKontraktu(n dane.NotatkaPrzegladania) shared.BrowserNote {
	notatka := shared.BrowserNote{
		Id:        n.Kod,
		WindowId:  n.Okno,
		SourceId:  n.ZrodloID,
		Content:   n.Tresc,
		Quote:     n.Cytat,
		CreatedAt: chwilaBazy(n.Utworzono),
		UpdatedAt: chwilaBazy(n.Zaktualizowano),
		ThreadId:  n.Watek,
	}
	if n.Klasyfikacja != nil {
		klasyfikacja := shared.BrowserNoteClassification(*n.Klasyfikacja)
		notatka.Classification = &klasyfikacja
	}
	if n.Przypieta {
		przypieta := true
		notatka.Pinned = &przypieta
	}
	return notatka
}

// Brak pola i fałsz jawny niosą tę samą wartość w kolumnie NOT NULL tabeli.
func wartoscLogiczna(w *bool) bool {
	return w != nil && *w
}
