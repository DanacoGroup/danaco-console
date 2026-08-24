// Odpowiedzialność pliku: obsługa `browser.source.add` i `browser.note.add` —
// dwie szuflady modułu Browser zasilane z toku przeglądania, każda swoim
// wierszem w `dane.RepozytoriumPrzegladania`.
//
// Warstwa danych oddaje znacznik czasu jako surowy tekst bazy (ISO 8601 ze
// strftime), a kontrakt chce milisekund epoki. Przekład robi `chwilaBazy`
// z `przeklad_nawigacja.go` — jedno miejsce przekładu czasu bazy dla całego
// rdzenia, nie własna kopia w każdym module.
//
// Notatka bez źródła jest dozwolona: `SourceId` w żądaniu jest opcjonalny, bo
// notatka Operatora może dotyczyć całej strony, nie tylko zebranego źródła
// (`migracja_047_przegladarka.sql`, `dane/przegladarka_notatki.go`). Adapter nie
// dogląda istnienia źródła przy zapisie notatki — kolumna nie niesie więzu
// obcego, więc i port go nie udaje.
package core

import (
	"context"
	"fmt"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// Przedrostki identyfikatorów nadawanych przez rdzeń — źródło i notatka nie
// przychodzą od klienta z własnym `Id` (kontrakt ich nie niesie w żądaniu),
// więc rdzeń nadaje je tak samo jak zdarzeniom i wiadomościom (`identyfikator.go`).
const (
	przedrostekZrodlaPrzegladania  = "src-"
	przedrostekNotatkiPrzegladania = "note-"
)

// DodajZrodlo zapisuje stronę w wykazie źródeł okna. Zapis jest idempotentny
// wobec identyfikatora zewnętrznego: powtórne wywołanie z tym samym `Id`
// nadpisuje wiersz, nie dubluje go (patrz `ZapiszZrodlo`).
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
	// Klasyfikacja przychodzi wyliczeniem kontraktu, a kolumna trzyma tekst:
	// przekład jest tu, a nie w warstwie danych, bo to kontrakt zna dopuszczalne
	// wartości i on je rozstrzyga przy wejściu żądania.
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

// brakiZrodlaPrzegladania i brakiNotatkiPrzegladania nazywają pole, które
// przyszło puste, zamiast mówić „okna albo adresu": odmowa idzie do Operatora
// i ma powiedzieć, co dopisać. Brak samego pola w treści żądania odsiewa brama
// kontraktu (`brama_kontraktu.go`) — tutaj rozstrzyga się wyłącznie wartość
// pusta, bo tylko dziedzina wie, że pusty adres źródłem nie jest.
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

// zrodloKontraktu przekłada wiersz źródła na byt kontraktu `BrowserSource`.
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

// notatkaKontraktu przekłada wiersz notatki na byt kontraktu `BrowserNote`.
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
