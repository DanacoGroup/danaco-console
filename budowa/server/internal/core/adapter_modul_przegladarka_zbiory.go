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

// Przedrostki identyfikatorów nadawanych przez rdzeń — źródło i notatka nie przychodzą od
// klienta z własnym identyfikatorem, więc rdzeń nadaje je tak samo jak innym bytom.
const (
	przedrostekZrodlaPrzegladania  = "src-"
	przedrostekNotatkiPrzegladania = "note-"
)

// DodajZrodlo zapisuje stronę w wykazie źródeł okna; zapis jest idempotentny wobec
// identyfikatora zewnętrznego — powtórne wywołanie z tym samym kodem nadpisuje wiersz, nie dubluje go.
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
		return shared.BrowserSourceAddResponse{}, fmt.Errorf("core: nie można dodać źródła przeglądania: %w", err)
	}
	return shared.BrowserSourceAddResponse{Source: zrodloKontraktu(zapisane)}, nil
}

// DodajNotatke zapisuje notatkę Operatora, opcjonalnie powiązaną ze źródłem
// i cytatem fragmentu strony.
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

// brakiZrodlaPrzegladania i brakiNotatkiPrzegladania nazywają pole, które przyszło puste,
// zamiast ogólnikowej odmowy — Operator ma się dowiedzieć, co dopisać do żądania.
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

// zrodloKontraktu przekłada wiersz źródła bazy danych na byt kontraktu `BrowserSource` odpowiedzi żądania.
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

// notatkaKontraktu przekłada wiersz notatki bazy danych na byt kontraktu `BrowserNote` odpowiedzi żądania.
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

// wartoscLogiczna odczytuje wskaźnik logiczny żądania — brak pola znaczy
// „nie zaznaczono", nie „fałsz jawnie wybrany", ale obie sytuacje niosą tę
// samą wartość w kolumnie NOT NULL tabeli.
func wartoscLogiczna(w *bool) bool {
	return w != nil && *w
}
