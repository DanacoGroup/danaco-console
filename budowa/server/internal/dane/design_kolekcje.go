// Kolekcje zasobów modułu Design: przypisanie jest dokładką albo odjęciem, nigdy podmianą kompletu.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// KolekcjaDesignu to wiersz `kolekcja_design`; `Zasoby` niesie identyfikatory pozycji w kolejności przypisania.
type KolekcjaDesignu struct {
	ID             int64
	Kod            string
	Okno           string
	Nazwa          string
	Opis           *string
	Zasoby         []string
	Liczba         int
	Zaktualizowano string
}

const (
	kolumnyKolekcjiDesignu = `k.id, k.identyfikator_zewnetrzny, k.okno, k.nazwa, k.opis,
	                          (SELECT COUNT(*) FROM pozycja_kolekcji_design p
	                            WHERE p.kolekcja_id = k.id),
	                          k.zaktualizowano`

	zapiszKolekcjeDesignu = `INSERT INTO kolekcja_design
	                         (identyfikator_zewnetrzny, okno, nazwa, opis, zaktualizowano, konto_id)
	                         VALUES (?, ?, ?, ?, strftime('%Y-%m-%dT%H:%M:%fZ','now'), ` + WskazanieKonta + `)
	                         ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                             nazwa = excluded.nazwa,
	                             opis = excluded.opis,
	                             zaktualizowano = excluded.zaktualizowano
	                         WHERE ` + WarunekKonta

	pobierzKolekcjeDesignu = `SELECT ` + kolumnyKolekcjiDesignu + ` FROM kolekcja_design k
	                          WHERE k.identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	listaKolekcjiDesignu = `SELECT ` + kolumnyKolekcjiDesignu + ` FROM kolekcja_design k
	                        WHERE k.okno = ? AND ` + WarunekKonta + `
	                        ORDER BY k.zaktualizowano DESC, k.id DESC`

	// Zawężenie do kolekcji zawierających wskazany zasób, po jego identyfikatorze zewnętrznym w żądaniu okna.
	listaKolekcjiDesignuZZasobem = `SELECT ` + kolumnyKolekcjiDesignu + ` FROM kolekcja_design k
	                                WHERE k.okno = ? AND ` + WarunekKonta + `
	                                  AND EXISTS (SELECT 1 FROM pozycja_kolekcji_design p
	                                               WHERE p.kolekcja_id = k.id AND p.zasob_id = ?)
	                                ORDER BY k.zaktualizowano DESC, k.id DESC`

	// Kolejność nowej pozycji jest o jeden dalsza od największej zastanej, żeby dokładanie zachowywało porządek dołączania.
	wstawPozycjeKolekcjiDesignu = `INSERT INTO pozycja_kolekcji_design (kolekcja_id, zasob_id, kolejnosc)
	                               VALUES (?, ?, (SELECT COALESCE(MAX(kolejnosc), 0) + 1
	                                                FROM pozycja_kolekcji_design WHERE kolekcja_id = ?))
	                               ON CONFLICT(kolekcja_id, zasob_id) DO NOTHING`

	usunPozycjeKolekcjiDesignu = `DELETE FROM pozycja_kolekcji_design
	                              WHERE kolekcja_id = ? AND zasob_id = ?`

	dotknijKolekcjeDesignu = `UPDATE kolekcja_design
	                          SET zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                          WHERE id = ?`

	listaZasobowKolekcjiDesignu = `SELECT zasob_id FROM pozycja_kolekcji_design
	                               WHERE kolekcja_id = ? ORDER BY kolejnosc, rowid`
)

// ZapiszKolekcjeDesignu zakłada kolekcję albo nadpisuje zastaną po kodzie i oddaje stan po zapisie.
func (r *repozytoriumDesignu) ZapiszKolekcjeDesignu(ctx context.Context,
	kolekcja KolekcjaDesignu) (KolekcjaDesignu, error) {

	if kolekcja.Kod == "" {
		return KolekcjaDesignu{}, fmt.Errorf("dane: kolekcja design bez identyfikatora")
	}
	if kolekcja.Okno == "" {
		return KolekcjaDesignu{}, fmt.Errorf("dane: kolekcja design %q bez okna", kolekcja.Kod)
	}
	if kolekcja.Nazwa == "" {
		return KolekcjaDesignu{}, fmt.Errorf("dane: kolekcja design %q bez nazwy", kolekcja.Kod)
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszKolekcjeDesignu)
	if err != nil {
		return KolekcjaDesignu{}, err
	}
	if _, err := polecenie.ExecContext(ctx, kolekcja.Kod, kolekcja.Okno, kolekcja.Nazwa,
		tekstDoKolumny(kolekcja.Opis), KontoOperatora(ctx), KontoOperatora(ctx)); err != nil {
		return KolekcjaDesignu{}, fmt.Errorf("dane: nie można zapisać kolekcji design %q: %w",
			kolekcja.Kod, err)
	}
	return r.KolekcjaDesignuPoKodzie(ctx, kolekcja.Kod)
}

// KolekcjaDesignuPoKodzie zwraca kolekcję po kodzie wraz z zasobami; brak wiersza wraca jako ErrBrakWiersza.
func (r *repozytoriumDesignu) KolekcjaDesignuPoKodzie(ctx context.Context,
	kod string) (KolekcjaDesignu, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzKolekcjeDesignu)
	if err != nil {
		return KolekcjaDesignu{}, err
	}
	kolekcja, err := odczytajKolekcjeDesignu(polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return KolekcjaDesignu{}, ErrBrakWiersza
	}
	if err != nil {
		return KolekcjaDesignu{}, fmt.Errorf("dane: nieczytelny wiersz kolekcji design %q: %w", kod, err)
	}
	zasoby, err := r.ZasobyKolekcjiDesignu(ctx, kolekcja.ID)
	if err != nil {
		return KolekcjaDesignu{}, err
	}
	kolekcja.Zasoby = zasoby
	return kolekcja, nil
}

// KolekcjeDesignu zwraca kolekcje okna wraz z zasobami; wskazanie zasobu zawęża wykaz do kolekcji z tym zasobem.
func (r *repozytoriumDesignu) KolekcjeDesignu(ctx context.Context,
	okno string, zasob *string) ([]KolekcjaDesignu, error) {

	zapytanie := listaKolekcjiDesignu
	argumenty := []any{okno, KontoOperatora(ctx)}
	if zasob != nil && *zasob != "" {
		zapytanie = listaKolekcjiDesignuZZasobem
		argumenty = append(argumenty, *zasob)
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, argumenty...)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać kolekcji design okna %q: %w", okno, err)
	}
	defer wiersze.Close()

	lista := []KolekcjaDesignu{}
	for wiersze.Next() {
		kolekcja, err := odczytajKolekcjeDesignu(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz kolekcji design okna %q: %w", okno, err)
		}
		lista = append(lista, kolekcja)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt kolekcji design okna %q: %w", okno, err)
	}
	// Zasoby idą osobnym odczytem po zamknięciu kursora: silnik bazy nie lubi zapytań zagnieżdżonych.
	for numer := range lista {
		zasoby, err := r.ZasobyKolekcjiDesignu(ctx, lista[numer].ID)
		if err != nil {
			return nil, err
		}
		lista[numer].Zasoby = zasoby
	}
	return lista, nil
}

// ZmienPrzypisaniaKolekcjiDesignu dokłada zasoby do kolekcji albo je z niej zdejmuje i oddaje liczbę przypisań, które naprawdę się zmieniły; całość idzie w jednej transakcji.
func (r *repozytoriumDesignu) ZmienPrzypisaniaKolekcjiDesignu(ctx context.Context,
	kolekcjaID int64, zasoby []string, zdejmij bool) (int, error) {

	zmienione := 0
	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		zapytanie := wstawPozycjeKolekcjiDesignu
		if zdejmij {
			zapytanie = usunPozycjeKolekcjiDesignu
		}
		polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, zapytanie)
		if err != nil {
			return err
		}
		for _, zasob := range zasoby {
			if zasob == "" {
				continue
			}
			var wynik sql.Result
			if zdejmij {
				wynik, err = polecenie.ExecContext(ctx, kolekcjaID, zasob)
			} else {
				wynik, err = polecenie.ExecContext(ctx, kolekcjaID, zasob, kolekcjaID)
			}
			if err != nil {
				return fmt.Errorf("dane: nie można zmienić przypisania zasobu %q w kolekcji design %d: %w",
					zasob, kolekcjaID, err)
			}
			ile, err := wynik.RowsAffected()
			if err != nil {
				return fmt.Errorf("dane: nie można policzyć zmian przypisań kolekcji design %d: %w",
					kolekcjaID, err)
			}
			zmienione += int(ile)
		}
		if zmienione == 0 {
			return nil
		}
		dotkniecie, err := r.zapytania.wTransakcji(ctx, transakcja, dotknijKolekcjeDesignu)
		if err != nil {
			return err
		}
		if _, err := dotkniecie.ExecContext(ctx, kolekcjaID); err != nil {
			return fmt.Errorf("dane: nie można odnotować zmiany kolekcji design %d: %w", kolekcjaID, err)
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return zmienione, nil
}

// ZasobyKolekcjiDesignu zwraca identyfikatory zewnętrzne zasobów kolekcji w kolejności przypisania z bazy.
func (r *repozytoriumDesignu) ZasobyKolekcjiDesignu(ctx context.Context,
	kolekcjaID int64) ([]string, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaZasobowKolekcjiDesignu)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, kolekcjaID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać zasobów kolekcji design %d: %w", kolekcjaID, err)
	}
	defer wiersze.Close()

	lista := []string{}
	for wiersze.Next() {
		var zasob string
		if err := wiersze.Scan(&zasob); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz zasobu kolekcji design %d: %w", kolekcjaID, err)
		}
		lista = append(lista, zasob)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt zasobów kolekcji design %d: %w", kolekcjaID, err)
	}
	return lista, nil
}

// odczytajKolekcjeDesignu składa strukturę kolekcji wprost z jednego wiersza wyniku zapytania SQL do bazy.
func odczytajKolekcjeDesignu(wiersz skaner) (KolekcjaDesignu, error) {
	var kolekcja KolekcjaDesignu
	var opis sql.NullString
	err := wiersz.Scan(&kolekcja.ID, &kolekcja.Kod, &kolekcja.Okno, &kolekcja.Nazwa,
		&opis, &kolekcja.Liczba, &kolekcja.Zaktualizowano)
	if err != nil {
		return KolekcjaDesignu{}, err
	}
	kolekcja.Opis = tekstZKolumny(opis)
	return kolekcja, nil
}
