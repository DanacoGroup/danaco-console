// Obszar gradientów wypełnienia modułu Design (tabele `gradient_design`
// i `stopien_gradientu_design`, migracja 326) — część `RepozytoriumDesignu`
// zadeklarowanego w `design.go`.
//
// Gradient jest jeden na cel. Celem jest kompozycja, a wskazanie ścieżki albo
// warstwy zawęża go do jednego bytu na niej — stąd klucz jedyny na trójce
// (kompozycja, ścieżka, warstwa). Wskazanie pominięte zapisuje się pustym
// napisem, nie NULL-em: SQLite liczy dwa NULL-e za różne w indeksie UNIQUE,
// więc gradient całej kompozycji zakładany dwa razy powstałby dwa razy zamiast
// nadpisać się raz.
//
// Stopnie podmieniają się kompletem, tak jak warstwy kompozycji: kontrakt
// (`DesignGradient.Stops`) nadsyła cały gradient, a stopień usunięty w oknie
// musi zniknąć także w bazie.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// GradientDesignu to wiersz tabeli `gradient_design` wraz z jego stopniami.
type GradientDesignu struct {
	ID             int64
	Kod            string
	KompozycjaID   int64
	SciezkaID      string
	WarstwaID      string
	Rodzaj         string
	Kat            *float64
	Stopnie        []StopienGradientuDesignu
	Zaktualizowano string
}

// StopienGradientuDesignu to wiersz tabeli `stopien_gradientu_design` — jeden
// stopień w kolejności, w jakiej wszedł.
type StopienGradientuDesignu struct {
	Polozenie float64
	Barwa     string
	Krycie    *float64
}

const (
	kolumnyGradientuDesignu = `g.id, g.identyfikator_zewnetrzny, g.kompozycja_id, g.sciezka_id,
	                           g.warstwa_id, g.rodzaj, g.kat, g.zaktualizowano`

	zapiszGradientDesignuSQL = `INSERT INTO gradient_design
	                            (identyfikator_zewnetrzny, kompozycja_id, sciezka_id, warstwa_id,
	                             rodzaj, kat, zaktualizowano)
	                            VALUES (?, ?, ?, ?, ?, ?, strftime('%Y-%m-%dT%H:%M:%fZ','now'))
	                            ON CONFLICT(kompozycja_id, sciezka_id, warstwa_id) DO UPDATE SET
	                                rodzaj = excluded.rodzaj,
	                                kat = excluded.kat,
	                                zaktualizowano = excluded.zaktualizowano`

	pobierzGradientDesignuPoCelu = `SELECT ` + kolumnyGradientuDesignu + ` FROM gradient_design g
	                                WHERE g.kompozycja_id = ? AND g.sciezka_id = ? AND g.warstwa_id = ?`

	usunStopnieGradientuDesignu = `DELETE FROM stopien_gradientu_design WHERE gradient_id = ?`

	wstawStopienGradientuDesignu = `INSERT INTO stopien_gradientu_design
	                                (gradient_id, kolejnosc, polozenie, barwa, krycie)
	                                VALUES (?, ?, ?, ?, ?)`

	listaStopniGradientuDesignu = `SELECT polozenie, barwa, krycie FROM stopien_gradientu_design
	                               WHERE gradient_id = ? ORDER BY kolejnosc`
)

// ZapiszGradientDesignu zakłada gradient celu albo nadpisuje zastany
// i podmienia komplet jego stopni w jednej transakcji.
//
// Identyfikator zewnętrzny nadaje wołający i służy wyłącznie temu, żeby
// gradient nowy dostał własną nazwę; rozstrzyga CEL, nie identyfikator —
// gradient wysłany drugi raz na tę samą warstwę nadpisuje ten, który tam stoi,
// choćby przyszedł z nowym identyfikatorem.
func (r *repozytoriumDesignu) ZapiszGradientDesignu(ctx context.Context,
	gradient GradientDesignu) (GradientDesignu, error) {

	if gradient.Kod == "" {
		return GradientDesignu{}, fmt.Errorf("dane: gradient design bez identyfikatora")
	}
	if gradient.KompozycjaID == 0 {
		return GradientDesignu{}, fmt.Errorf("dane: gradient design %q bez kompozycji", gradient.Kod)
	}
	if gradient.Rodzaj == "" {
		return GradientDesignu{}, fmt.Errorf("dane: gradient design %q bez rodzaju", gradient.Kod)
	}
	if len(gradient.Stopnie) == 0 {
		return GradientDesignu{}, fmt.Errorf("dane: gradient design %q bez stopni", gradient.Kod)
	}

	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		zapis, err := r.zapytania.wTransakcji(ctx, transakcja, zapiszGradientDesignuSQL)
		if err != nil {
			return err
		}
		if _, err := zapis.ExecContext(ctx, gradient.Kod, gradient.KompozycjaID,
			gradient.SciezkaID, gradient.WarstwaID, gradient.Rodzaj,
			liczbaRzeczywistaDoKolumny(gradient.Kat)); err != nil {
			return fmt.Errorf("dane: nie można zapisać gradientu design %q: %w", gradient.Kod, err)
		}

		odczyt, err := r.zapytania.wTransakcji(ctx, transakcja, pobierzGradientDesignuPoCelu)
		if err != nil {
			return err
		}
		zapisany, err := odczytajGradientDesignu(odczyt.QueryRowContext(ctx,
			gradient.KompozycjaID, gradient.SciezkaID, gradient.WarstwaID))
		if err != nil {
			return fmt.Errorf("dane: nie można odczytać zapisanego gradientu design %q: %w",
				gradient.Kod, err)
		}

		czyszczenie, err := r.zapytania.wTransakcji(ctx, transakcja, usunStopnieGradientuDesignu)
		if err != nil {
			return err
		}
		if _, err := czyszczenie.ExecContext(ctx, zapisany.ID); err != nil {
			return fmt.Errorf("dane: nie można wyczyścić stopni gradientu design %q: %w",
				gradient.Kod, err)
		}
		wstawienie, err := r.zapytania.wTransakcji(ctx, transakcja, wstawStopienGradientuDesignu)
		if err != nil {
			return err
		}
		for kolejnosc, stopien := range gradient.Stopnie {
			if _, err := wstawienie.ExecContext(ctx, zapisany.ID, kolejnosc,
				stopien.Polozenie, stopien.Barwa,
				liczbaRzeczywistaDoKolumny(stopien.Krycie)); err != nil {
				return fmt.Errorf("dane: nie można zapisać stopnia gradientu design %q: %w",
					gradient.Kod, err)
			}
		}
		return nil
	})
	if err != nil {
		return GradientDesignu{}, err
	}
	return r.GradientDesignuPoCelu(ctx, gradient.KompozycjaID, gradient.SciezkaID, gradient.WarstwaID)
}

// GradientDesignuPoCelu zwraca gradient stojący na wskazanym celu wraz ze
// stopniami. Brak wiersza wraca jako ErrBrakWiersza.
func (r *repozytoriumDesignu) GradientDesignuPoCelu(ctx context.Context,
	kompozycjaID int64, sciezka, warstwa string) (GradientDesignu, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzGradientDesignuPoCelu)
	if err != nil {
		return GradientDesignu{}, err
	}
	gradient, err := odczytajGradientDesignu(polecenie.QueryRowContext(ctx,
		kompozycjaID, sciezka, warstwa))
	if errors.Is(err, sql.ErrNoRows) {
		return GradientDesignu{}, ErrBrakWiersza
	}
	if err != nil {
		return GradientDesignu{}, fmt.Errorf("dane: nieczytelny wiersz gradientu design: %w", err)
	}
	stopnie, err := r.stopnieGradientuDesignu(ctx, gradient.ID)
	if err != nil {
		return GradientDesignu{}, err
	}
	gradient.Stopnie = stopnie
	return gradient, nil
}

// stopnieGradientuDesignu czyta stopnie jednego gradientu w ich kolejności.
func (r *repozytoriumDesignu) stopnieGradientuDesignu(ctx context.Context,
	gradientID int64) ([]StopienGradientuDesignu, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaStopniGradientuDesignu)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, gradientID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać stopni gradientu design %d: %w", gradientID, err)
	}
	defer wiersze.Close()

	stopnie := []StopienGradientuDesignu{}
	for wiersze.Next() {
		var stopien StopienGradientuDesignu
		var krycie sql.NullFloat64
		if err := wiersze.Scan(&stopien.Polozenie, &stopien.Barwa, &krycie); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny stopień gradientu design %d: %w", gradientID, err)
		}
		stopien.Krycie = liczbaRzeczywistaZKolumny(krycie)
		stopnie = append(stopnie, stopien)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt stopni gradientu design %d: %w", gradientID, err)
	}
	return stopnie, nil
}

// odczytajGradientDesignu składa wiersz gradientu ze skanera.
func odczytajGradientDesignu(wiersz skaner) (GradientDesignu, error) {
	var gradient GradientDesignu
	var kat sql.NullFloat64

	if err := wiersz.Scan(&gradient.ID, &gradient.Kod, &gradient.KompozycjaID,
		&gradient.SciezkaID, &gradient.WarstwaID, &gradient.Rodzaj, &kat,
		&gradient.Zaktualizowano); err != nil {
		return GradientDesignu{}, err
	}
	gradient.Kat = liczbaRzeczywistaZKolumny(kat)
	gradient.Stopnie = []StopienGradientuDesignu{}
	return gradient, nil
}
