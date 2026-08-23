// Odpowiedzialność pliku: ślady wymiany słownika (Glossary Panel) modułu
// Translate — import i eksport terminów (TBX/CSV), tabele
// `slad_importu_slownika` i `slad_eksportu_slownika`. Typ, interfejs
// i konstruktor deklaruje `tlumaczenie.go`; tu wyłącznie implementacja czterech
// metod wymiany na `*repozytoriumTlumaczen`.
//
// Tabele są dwie, bo import i eksport to różne kierunki z różną kolumną wyniku
// (`liczba_zaimportowanych` kontra `liczba_wyeksportowanych`) — wspólna tabela
// byłaby dwiema prawdami o jednym bycie.
//
// Ślad eksportu nie niesie dowodu powstania pliku: rdzeń nie ma magazynu blobów,
// a kontrakt `TranslateGlossaryExportResponse` oddaje wyłącznie `ExportedCount`,
// bez identyfikatora pliku ani rozmiaru — inaczej niż `research.report.export`
// (`dane/badania_raport.go`, `EksportRaportu.PlikBibliotekiID`/`RozmiarBajtow`).
// `Sciezka` niesie więc ścieżkę żądaną, nie ścieżkę wyniku.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// SladImportuSlownika to wiersz `slad_importu_slownika` — ślad wyniku importu
// glosariusza (translate.glossary.import): ile terminów weszło i z jakiej
// ścieżki.
type SladImportuSlownika struct {
	ID                    int64
	Kod                   string
	Sciezka               string
	LiczbaZaimportowanych int64
	Utworzono             int64
}

// SladEksportuSlownika to wiersz `slad_eksportu_slownika` — ślad wyniku
// eksportu glosariusza (translate.glossary.export). `Sciezka` to ścieżka
// docelowa zgłoszona w żądaniu, nie dowód powstania pliku (rdzeń nie ma
// magazynu blobów).
type SladEksportuSlownika struct {
	ID                     int64
	Kod                    string
	Sciezka                string
	LiczbaWyeksportowanych int64
	Utworzono              int64
}

const (
	kolumnySladuImportuSlownika = `id, identyfikator_zewnetrzny, sciezka, liczba_zaimportowanych, utworzono`
	zapiszSladImportuSlownika   = `INSERT INTO slad_importu_slownika
	                               (identyfikator_zewnetrzny, sciezka, liczba_zaimportowanych, utworzono)
	                               VALUES (?, ?, ?, ?)`
	pobierzSladImportuSlownika = `SELECT ` + kolumnySladuImportuSlownika + ` FROM slad_importu_slownika
	                              WHERE identyfikator_zewnetrzny = ?`
	listaSladowImportuSlownika = `SELECT ` + kolumnySladuImportuSlownika + ` FROM slad_importu_slownika
	                              ORDER BY utworzono DESC LIMIT ?`

	kolumnySladuEksportuSlownika = `id, identyfikator_zewnetrzny, sciezka, liczba_wyeksportowanych, utworzono`
	zapiszSladEksportuSlownika   = `INSERT INTO slad_eksportu_slownika
	                               (identyfikator_zewnetrzny, sciezka, liczba_wyeksportowanych, utworzono)
	                               VALUES (?, ?, ?, ?)`
	pobierzSladEksportuSlownika = `SELECT ` + kolumnySladuEksportuSlownika + ` FROM slad_eksportu_slownika
	                               WHERE identyfikator_zewnetrzny = ?`
	listaSladowEksportuSlownika = `SELECT ` + kolumnySladuEksportuSlownika + ` FROM slad_eksportu_slownika
	                               ORDER BY utworzono DESC LIMIT ?`
)

// ZapiszImportSlownika dokłada ślad importu (historia, nic nie
// nadpisuje) — każde wywołanie `glossary.import` to osobny fakt o wyniku.
func (r *repozytoriumTlumaczen) ZapiszImportSlownika(ctx context.Context, slad SladImportuSlownika) (SladImportuSlownika, error) {
	if slad.Kod == "" {
		return SladImportuSlownika{}, fmt.Errorf("dane: ślad importu słownika bez identyfikatora")
	}
	if slad.Sciezka == "" {
		return SladImportuSlownika{}, fmt.Errorf("dane: ślad importu słownika %q bez ścieżki", slad.Kod)
	}
	teraz := slad.Utworzono
	if teraz == 0 {
		teraz = time.Now().UnixMilli()
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszSladImportuSlownika)
	if err != nil {
		return SladImportuSlownika{}, err
	}
	if _, err := polecenie.ExecContext(ctx, slad.Kod, slad.Sciezka, slad.LiczbaZaimportowanych, teraz); err != nil {
		return SladImportuSlownika{}, fmt.Errorf("dane: nie można zapisać śladu importu słownika %q: %w", slad.Kod, err)
	}
	odczyt, err := r.zapytania.przygotuj(ctx, pobierzSladImportuSlownika)
	if err != nil {
		return SladImportuSlownika{}, err
	}
	zapisany, err := odczytajSladImportuSlownika(odczyt.QueryRowContext(ctx, slad.Kod))
	if err != nil {
		return SladImportuSlownika{}, fmt.Errorf("dane: nie można odczytać zapisanego śladu importu słownika %q: %w", slad.Kod, err)
	}
	return zapisany, nil
}

// ZapiszEksportSlownika dokłada ślad eksportu — `Sciezka` to
// ścieżka żądana z `TranslateGlossaryExportRequest.Path`, nie dowód powstania
// pliku (rdzeń nie ma magazynu blobów).
func (r *repozytoriumTlumaczen) ZapiszEksportSlownika(ctx context.Context, slad SladEksportuSlownika) (SladEksportuSlownika, error) {
	if slad.Kod == "" {
		return SladEksportuSlownika{}, fmt.Errorf("dane: ślad eksportu słownika bez identyfikatora")
	}
	if slad.Sciezka == "" {
		return SladEksportuSlownika{}, fmt.Errorf("dane: ślad eksportu słownika %q bez ścieżki", slad.Kod)
	}
	teraz := slad.Utworzono
	if teraz == 0 {
		teraz = time.Now().UnixMilli()
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszSladEksportuSlownika)
	if err != nil {
		return SladEksportuSlownika{}, err
	}
	if _, err := polecenie.ExecContext(ctx, slad.Kod, slad.Sciezka, slad.LiczbaWyeksportowanych, teraz); err != nil {
		return SladEksportuSlownika{}, fmt.Errorf("dane: nie można zapisać śladu eksportu słownika %q: %w", slad.Kod, err)
	}
	odczyt, err := r.zapytania.przygotuj(ctx, pobierzSladEksportuSlownika)
	if err != nil {
		return SladEksportuSlownika{}, err
	}
	zapisany, err := odczytajSladEksportuSlownika(odczyt.QueryRowContext(ctx, slad.Kod))
	if err != nil {
		return SladEksportuSlownika{}, fmt.Errorf("dane: nie można odczytać zapisanego śladu eksportu słownika %q: %w", slad.Kod, err)
	}
	return zapisany, nil
}

// ImportySlownika zwraca najświeższe ślady importu, najnowsze pierwsze.
func (r *repozytoriumTlumaczen) ImportySlownika(ctx context.Context, limit int) ([]SladImportuSlownika, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaSladowImportuSlownika)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać śladów importu słownika: %w", err)
	}
	defer wiersze.Close()

	lista := []SladImportuSlownika{}
	for wiersze.Next() {
		slad, err := odczytajSladImportuSlownika(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz śladu importu słownika: %w", err)
		}
		lista = append(lista, slad)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt śladów importu słownika: %w", err)
	}
	return lista, nil
}

// EksportySlownika zwraca najświeższe ślady eksportu, najnowsze pierwsze.
func (r *repozytoriumTlumaczen) EksportySlownika(ctx context.Context, limit int) ([]SladEksportuSlownika, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaSladowEksportuSlownika)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać śladów eksportu słownika: %w", err)
	}
	defer wiersze.Close()

	lista := []SladEksportuSlownika{}
	for wiersze.Next() {
		slad, err := odczytajSladEksportuSlownika(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz śladu eksportu słownika: %w", err)
		}
		lista = append(lista, slad)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt śladów eksportu słownika: %w", err)
	}
	return lista, nil
}

// odczytajSladImportuSlownika składa strukturę z jednego wiersza wyniku. Brak
// wiersza wraca jako ErrBrakWiersza.
func odczytajSladImportuSlownika(wiersz skaner) (SladImportuSlownika, error) {
	var slad SladImportuSlownika
	err := wiersz.Scan(&slad.ID, &slad.Kod, &slad.Sciezka, &slad.LiczbaZaimportowanych, &slad.Utworzono)
	if errors.Is(err, sql.ErrNoRows) {
		return SladImportuSlownika{}, ErrBrakWiersza
	}
	return slad, err
}

// odczytajSladEksportuSlownika składa strukturę z jednego wiersza wyniku.
func odczytajSladEksportuSlownika(wiersz skaner) (SladEksportuSlownika, error) {
	var slad SladEksportuSlownika
	err := wiersz.Scan(&slad.ID, &slad.Kod, &slad.Sciezka, &slad.LiczbaWyeksportowanych, &slad.Utworzono)
	if errors.Is(err, sql.ErrNoRows) {
		return SladEksportuSlownika{}, ErrBrakWiersza
	}
	return slad, err
}
