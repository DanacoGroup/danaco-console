// Obszar ikon własnych modułu Design (tabele `ikona_design`
// i `etykieta_ikony_design`, migracja 326) — część `RepozytoriumDesignu`
// zadeklarowanego w `design.go`.
//
// W bazie leżą wyłącznie ikony WŁASNE: te narysowane na siatce
// (`design.icon.set`) i te wytworzone kanałem modelu (`design.icon.generate`).
// Katalog ikon otwartoźródłowych bazy nie dotyka — jest wkompilowany w binarium
// (`core/ikony_katalogu`), więc rdzeń zna go bez kroku zasiewu, a stanowisko
// bez ani jednej ikony własnej i tak ma czego szukać.
//
// Etykiety podmieniają się kompletem, wzorem `UstawEtykietyZasobu`: kontrakt
// (`DesignIconSetRequest.Tags`) nadsyła stan docelowy, więc etykieta zdjęta
// w oknie musi zniknąć także w bazie.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// IkonaDesignu to wiersz tabeli `ikona_design` wraz z etykietami wyszukiwania.
// Etykiety leżą w osobnej tabeli i wchodzą tym samym odczytem, bo `DesignIcon`
// kontraktu niesie je przy każdej ikonie.
type IkonaDesignu struct {
	ID             int64
	Kod            string
	Okno           string
	Nazwa          string
	Zestaw         *string
	SVG            string
	Siatka         *int64
	GruboscObrysu  *float64
	Etykiety       []string
	Zaktualizowano string
}

const (
	kolumnyIkonyDesignu = `i.id, i.identyfikator_zewnetrzny, i.okno, i.nazwa, i.zestaw,
	                       i.svg, i.siatka, i.grubosc_obrysu, i.zaktualizowano`

	zapiszIkoneDesignuSQL = `INSERT INTO ikona_design
	                         (identyfikator_zewnetrzny, okno, nazwa, zestaw, svg, siatka,
	                          grubosc_obrysu, zaktualizowano)
	                         VALUES (?, ?, ?, ?, ?, ?, ?, strftime('%Y-%m-%dT%H:%M:%fZ','now'))
	                         ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                             nazwa = excluded.nazwa,
	                             zestaw = excluded.zestaw,
	                             svg = excluded.svg,
	                             siatka = excluded.siatka,
	                             grubosc_obrysu = excluded.grubosc_obrysu,
	                             zaktualizowano = excluded.zaktualizowano`

	pobierzIkoneDesignu = `SELECT ` + kolumnyIkonyDesignu + ` FROM ikona_design i
	                       WHERE i.identyfikator_zewnetrzny = ?`

	listaIkonDesignuOkna = `SELECT ` + kolumnyIkonyDesignu + ` FROM ikona_design i
	                        WHERE i.okno = ? ORDER BY i.zaktualizowano DESC, i.id DESC`

	usunEtykietyIkonyDesignu = `DELETE FROM etykieta_ikony_design WHERE ikona_id = ?`

	wstawEtykieteIkonyDesignu = `INSERT INTO etykieta_ikony_design (ikona_id, etykieta)
	                             VALUES (?, ?) ON CONFLICT(ikona_id, etykieta) DO NOTHING`

	listaEtykietIkonyDesignu = `SELECT etykieta FROM etykieta_ikony_design
	                            WHERE ikona_id = ? ORDER BY etykieta`
)

// ZapiszIkoneDesignu zakłada ikonę albo nadpisuje zastaną po identyfikatorze
// zewnętrznym i podmienia komplet jej etykiet w jednej transakcji.
//
// Nazwa ikony jest jedyna w oknie (indeks UNIQUE migracji 326). Zderzenie
// wychodzi stąd jako błąd bazy, nie jako cicha podmiana cudzej ikony: dwie
// ikony o jednej nazwie w jednym oknie dałyby sprite z dwoma symbolami o tym
// samym `id`, czyli plik, którego przeglądarka nie złoży.
func (r *repozytoriumDesignu) ZapiszIkoneDesignu(ctx context.Context,
	ikona IkonaDesignu) (IkonaDesignu, error) {

	if ikona.Kod == "" {
		return IkonaDesignu{}, fmt.Errorf("dane: ikona design bez identyfikatora")
	}
	if ikona.Okno == "" {
		return IkonaDesignu{}, fmt.Errorf("dane: ikona design %q bez okna", ikona.Kod)
	}
	if ikona.Nazwa == "" {
		return IkonaDesignu{}, fmt.Errorf("dane: ikona design %q bez nazwy", ikona.Kod)
	}
	if ikona.SVG == "" {
		return IkonaDesignu{}, fmt.Errorf("dane: ikona design %q bez treści SVG", ikona.Kod)
	}

	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		zapis, err := r.zapytania.wTransakcji(ctx, transakcja, zapiszIkoneDesignuSQL)
		if err != nil {
			return err
		}
		if _, err := zapis.ExecContext(ctx, ikona.Kod, ikona.Okno, ikona.Nazwa,
			tekstDoKolumny(ikona.Zestaw), ikona.SVG, liczbaDoKolumny(ikona.Siatka),
			liczbaRzeczywistaDoKolumny(ikona.GruboscObrysu)); err != nil {
			return fmt.Errorf("dane: nie można zapisać ikony design %q: %w", ikona.Kod, err)
		}

		// Klucz wiersza wchodzi dopiero po zapisie — ikona mogła powstać w tej
		// transakcji, a `etykieta_ikony_design.ikona_id` wymaga klucza.
		odczyt, err := r.zapytania.wTransakcji(ctx, transakcja, pobierzIkoneDesignu)
		if err != nil {
			return err
		}
		zapisana, err := odczytajIkoneDesignu(odczyt.QueryRowContext(ctx, ikona.Kod))
		if err != nil {
			return fmt.Errorf("dane: nie można odczytać zapisanej ikony design %q: %w", ikona.Kod, err)
		}

		czyszczenie, err := r.zapytania.wTransakcji(ctx, transakcja, usunEtykietyIkonyDesignu)
		if err != nil {
			return err
		}
		if _, err := czyszczenie.ExecContext(ctx, zapisana.ID); err != nil {
			return fmt.Errorf("dane: nie można wyczyścić etykiet ikony design %q: %w", ikona.Kod, err)
		}
		if len(ikona.Etykiety) == 0 {
			return nil
		}
		wstawienie, err := r.zapytania.wTransakcji(ctx, transakcja, wstawEtykieteIkonyDesignu)
		if err != nil {
			return err
		}
		for _, etykieta := range ikona.Etykiety {
			if etykieta == "" {
				continue
			}
			if _, err := wstawienie.ExecContext(ctx, zapisana.ID, etykieta); err != nil {
				return fmt.Errorf("dane: nie można zapisać etykiety ikony design %q: %w", ikona.Kod, err)
			}
		}
		return nil
	})
	if err != nil {
		return IkonaDesignu{}, err
	}
	return r.IkonaDesignuPoKodzie(ctx, ikona.Kod)
}

// IkonaDesignuPoKodzie zwraca ikonę o wskazanym identyfikatorze zewnętrznym
// wraz z etykietami. Brak wiersza wraca jako ErrBrakWiersza.
func (r *repozytoriumDesignu) IkonaDesignuPoKodzie(ctx context.Context,
	kod string) (IkonaDesignu, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzIkoneDesignu)
	if err != nil {
		return IkonaDesignu{}, err
	}
	ikona, err := odczytajIkoneDesignu(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return IkonaDesignu{}, ErrBrakWiersza
	}
	if err != nil {
		return IkonaDesignu{}, fmt.Errorf("dane: nieczytelny wiersz ikony design %q: %w", kod, err)
	}
	etykiety, err := r.etykietyIkonyDesignu(ctx, ikona.ID)
	if err != nil {
		return IkonaDesignu{}, err
	}
	ikona.Etykiety = etykiety
	return ikona, nil
}

// IkonyDesignuOkna zwraca ikony własne okna, od ostatnio zmienianej, wraz
// z etykietami.
func (r *repozytoriumDesignu) IkonyDesignuOkna(ctx context.Context,
	okno string) ([]IkonaDesignu, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaIkonDesignuOkna)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać ikon design okna %q: %w", okno, err)
	}
	defer wiersze.Close()

	lista := []IkonaDesignu{}
	for wiersze.Next() {
		ikona, err := odczytajIkoneDesignu(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz ikony design okna %q: %w", okno, err)
		}
		lista = append(lista, ikona)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt ikon design okna %q: %w", okno, err)
	}
	// Etykiety idą po zamknięciu kursora — powód ten sam, co przy kolekcjach.
	for numer := range lista {
		etykiety, err := r.etykietyIkonyDesignu(ctx, lista[numer].ID)
		if err != nil {
			return nil, err
		}
		lista[numer].Etykiety = etykiety
	}
	return lista, nil
}

// etykietyIkonyDesignu czyta etykiety jednej ikony.
func (r *repozytoriumDesignu) etykietyIkonyDesignu(ctx context.Context,
	ikonaID int64) ([]string, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaEtykietIkonyDesignu)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, ikonaID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać etykiet ikony design %d: %w", ikonaID, err)
	}
	defer wiersze.Close()

	etykiety := []string{}
	for wiersze.Next() {
		var etykieta string
		if err := wiersze.Scan(&etykieta); err != nil {
			return nil, fmt.Errorf("dane: nieczytelna etykieta ikony design %d: %w", ikonaID, err)
		}
		etykiety = append(etykiety, etykieta)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt etykiet ikony design %d: %w", ikonaID, err)
	}
	return etykiety, nil
}

// odczytajIkoneDesignu składa wiersz ikony ze skanera — jedno miejsce
// odwzorowania kolumn na pola dla odczytu pojedynczego i dla wykazu.
func odczytajIkoneDesignu(wiersz skaner) (IkonaDesignu, error) {
	var ikona IkonaDesignu
	var zestaw sql.NullString
	var siatka sql.NullInt64
	var grubosc sql.NullFloat64

	if err := wiersz.Scan(&ikona.ID, &ikona.Kod, &ikona.Okno, &ikona.Nazwa, &zestaw,
		&ikona.SVG, &siatka, &grubosc, &ikona.Zaktualizowano); err != nil {
		return IkonaDesignu{}, err
	}
	ikona.Zestaw = tekstZKolumny(zestaw)
	ikona.Siatka = liczbaZKolumny(siatka)
	ikona.GruboscObrysu = liczbaRzeczywistaZKolumny(grubosc)
	ikona.Etykiety = []string{}
	return ikona, nil
}
