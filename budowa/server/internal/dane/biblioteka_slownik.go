// Odpowiedzialność pliku: słownik etykiet biblioteki (`etykieta_slownika_biblioteki`) i tezaurus
// relacji między etykietami (`relacja_tezaurusa_biblioteki`).
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type EtykietaSlownikaBiblioteki struct {
	Nazwa        string
	Barwa        *string
	LiczbaPlikow int
	Utworzono    string
}

type RelacjaTezaurusaBiblioteki struct {
	Zrodlo    string
	Cel       string
	Rodzaj    string
	Utworzono string
}

const (
	// Etykieta przy zasobie własnej kolumny konta nie ma: granica dochodzi przez `plik_id`.
	warunekKontaEtykietyPliku = `EXISTS (SELECT 1 FROM plik_biblioteki
	                                      WHERE plik_biblioteki.id = etykieta_pliku_biblioteki.plik_id
	                                        AND ` + WarunekKonta + `)`

	zapiszEtykieteSlownika = `INSERT INTO etykieta_slownika_biblioteki (nazwa, barwa, konto_id)
	                          VALUES (?, ?, ` + WskazanieKonta + `)
	                          ON CONFLICT(nazwa) DO UPDATE SET
	                              barwa = COALESCE(excluded.barwa, etykieta_slownika_biblioteki.barwa)
	                          WHERE ` + WarunekKonta

	pobierzEtykieteSlownika = `SELECT s.nazwa, s.barwa, s.utworzono,
	                                  (SELECT COUNT(*) FROM etykieta_pliku_biblioteki
	                                    WHERE etykieta_pliku_biblioteki.etykieta = s.nazwa
	                                      AND ` + warunekKontaEtykietyPliku + `) AS uzycie
	                           FROM etykieta_slownika_biblioteki s
	                           WHERE s.nazwa = ? AND ` + WarunekKonta

	// Wykaz łączy wpisy słownika i etykiety nadane przy zasobach; UNION zdejmuje powtórzenia.
	wykazEtykietSlownika = `WITH nazwy AS (
	                            SELECT nazwa FROM etykieta_slownika_biblioteki
	                             WHERE ` + WarunekKonta + `
	                            UNION
	                            SELECT etykieta AS nazwa FROM etykieta_pliku_biblioteki
	                             WHERE ` + warunekKontaEtykietyPliku + `
	                        )
	                        SELECT n.nazwa,
	                               (SELECT barwa FROM etykieta_slownika_biblioteki s
	                                 WHERE s.nazwa = n.nazwa AND ` + WarunekKonta + `) AS barwa,
	                               COALESCE((SELECT utworzono FROM etykieta_slownika_biblioteki s
	                                          WHERE s.nazwa = n.nazwa AND ` + WarunekKonta + `), '') AS utworzono,
	                               (SELECT COUNT(*) FROM etykieta_pliku_biblioteki
	                                 WHERE etykieta_pliku_biblioteki.etykieta = n.nazwa
	                                   AND ` + warunekKontaEtykietyPliku + `) AS uzycie
	                        FROM nazwy n`

	przemianujEtykietePliku = `UPDATE OR REPLACE etykieta_pliku_biblioteki
	                           SET etykieta = ? WHERE etykieta = ? AND ` + warunekKontaEtykietyPliku

	// Nazwa jest kluczem głównym całej tabeli: OR REPLACE kasowałby wpis konta cudzego, więc kolizja wraca błędem.
	przemianujEtykieteSlownika = `UPDATE etykieta_slownika_biblioteki
	                              SET nazwa = ? WHERE nazwa = ? AND ` + WarunekKonta

	wpisSlownikaWKoncie = `SELECT 1 FROM etykieta_slownika_biblioteki
	                       WHERE nazwa = ? AND ` + WarunekKonta

	usunEtykietePlikow = `DELETE FROM etykieta_pliku_biblioteki
	                      WHERE etykieta = ? AND ` + warunekKontaEtykietyPliku

	usunEtykieteSlownikaZapis = `DELETE FROM etykieta_slownika_biblioteki
	                             WHERE nazwa = ? AND ` + WarunekKonta

	wstawRelacjeTezaurusa = `INSERT INTO relacja_tezaurusa_biblioteki
	                         (etykieta_zrodlowa, etykieta_docelowa, rodzaj)
	                         VALUES (?, ?, ?)
	                         ON CONFLICT(etykieta_zrodlowa, etykieta_docelowa, rodzaj) DO NOTHING`

	usunRelacjeTezaurusa = `DELETE FROM relacja_tezaurusa_biblioteki
	                        WHERE etykieta_zrodlowa = ? AND etykieta_docelowa = ? AND rodzaj = ?`

	wykazRelacjiTezaurusa = `SELECT etykieta_zrodlowa, etykieta_docelowa, rodzaj, utworzono
	                         FROM relacja_tezaurusa_biblioteki
	                         ORDER BY etykieta_zrodlowa, rodzaj, etykieta_docelowa`
)

func (r *repozytoriumBiblioteki) EtykietySlownika(ctx context.Context, fraza *string,
	tylkoNieuzywane bool, limit int) ([]EtykietaSlownikaBiblioteki, int, error) {

	warunki := []string{"1 = 1"}
	argumenty := []any{}
	if fraza != nil && strings.TrimSpace(*fraza) != "" {
		warunki = append(warunki, "n.nazwa LIKE ?")
		argumenty = append(argumenty, "%"+strings.TrimSpace(*fraza)+"%")
	}
	if tylkoNieuzywane {
		warunki = append(warunki, "uzycie = 0")
	}
	warunek := strings.Join(warunki, " AND ")

	wskazania := wskazaniaKonta(wykazEtykietSlownika, KontoOperatora(ctx))
	zapytanie := wykazEtykietSlownika + ` WHERE ` + warunek + ` ORDER BY n.nazwa LIMIT ?`
	wiersze, err := r.db.QueryContext(ctx, zapytanie,
		append(append(append([]any{}, wskazania...), argumenty...), granicaWykazu(limit))...)
	if err != nil {
		return nil, 0, fmt.Errorf("dane: nie można odczytać słownika etykiet: %w", err)
	}
	defer wiersze.Close()

	lista := []EtykietaSlownikaBiblioteki{}
	for wiersze.Next() {
		etykieta, err := odczytajEtykieteSlownika(wiersze)
		if err != nil {
			return nil, 0, fmt.Errorf("dane: nieczytelny wiersz słownika etykiet: %w", err)
		}
		lista = append(lista, etykieta)
	}
	if err := wiersze.Err(); err != nil {
		return nil, 0, fmt.Errorf("dane: przerwany odczyt słownika etykiet: %w", err)
	}

	var lacznie int
	zapytanieLiczby := `SELECT COUNT(*) FROM (` + wykazEtykietSlownika + ` WHERE ` + warunek + `)`
	if err := r.db.QueryRowContext(ctx, zapytanieLiczby,
		append(append([]any{}, wskazania...), argumenty...)...).Scan(&lacznie); err != nil {
		return nil, 0, fmt.Errorf("dane: nie można policzyć etykiet słownika: %w", err)
	}
	return lista, lacznie, nil
}

// Etykieta nosząca zasoby bez wpisu słownikowego wraca z licznikiem i pustym czasem założenia.
func (r *repozytoriumBiblioteki) EtykietaSlownika(ctx context.Context,
	nazwa string) (EtykietaSlownikaBiblioteki, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzEtykieteSlownika)
	if err != nil {
		return EtykietaSlownikaBiblioteki{}, err
	}
	var etykieta EtykietaSlownikaBiblioteki
	var barwa sql.NullString
	err = polecenie.QueryRowContext(ctx, KontoOperatora(ctx), nazwa, KontoOperatora(ctx)).
		Scan(&etykieta.Nazwa, &barwa, &etykieta.Utworzono, &etykieta.LiczbaPlikow)
	if errors.Is(err, sql.ErrNoRows) {
		uzycie, err := r.uzycieEtykiety(ctx, nazwa)
		if err != nil {
			return EtykietaSlownikaBiblioteki{}, err
		}
		if uzycie == 0 {
			return EtykietaSlownikaBiblioteki{}, ErrBrakWiersza
		}
		return EtykietaSlownikaBiblioteki{Nazwa: nazwa, LiczbaPlikow: uzycie}, nil
	}
	if err != nil {
		return EtykietaSlownikaBiblioteki{}, fmt.Errorf("dane: nieczytelna etykieta słownika %q: %w",
			nazwa, err)
	}
	etykieta.Barwa = tekstZKolumny(barwa)
	return etykieta, nil
}

func (r *repozytoriumBiblioteki) ZapiszEtykieteSlownika(ctx context.Context, nazwa string,
	barwa *string) (EtykietaSlownikaBiblioteki, error) {

	if strings.TrimSpace(nazwa) == "" {
		return EtykietaSlownikaBiblioteki{}, fmt.Errorf("dane: etykieta słownika bez nazwy")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszEtykieteSlownika)
	if err != nil {
		return EtykietaSlownikaBiblioteki{}, err
	}
	wynik, err := polecenie.ExecContext(ctx, nazwa, tekstDoKolumny(barwa),
		KontoOperatora(ctx), KontoOperatora(ctx))
	if err != nil {
		return EtykietaSlownikaBiblioteki{}, fmt.Errorf("dane: nie można zapisać etykiety %q: %w",
			nazwa, err)
	}
	if err := sprawdzTrafienieZapisu(wynik, "etykieta słownika biblioteki", nazwa); err != nil {
		return EtykietaSlownikaBiblioteki{}, err
	}
	return r.EtykietaSlownika(ctx, nazwa)
}

// Wpis docelowy stojący już w koncie zostaje, a wpis źródłowy schodzi: `library.tag.merge` celuje w nazwę istniejącą.
func (r *repozytoriumBiblioteki) PrzemianujEtykiete(ctx context.Context, stara, nowa string) (int, error) {
	if strings.TrimSpace(stara) == "" || strings.TrimSpace(nowa) == "" {
		return 0, fmt.Errorf("dane: zmiana nazwy etykiety bez wskazania nazw")
	}
	dotkniete := 0
	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		przyZasobach, err := r.zapytania.wTransakcji(ctx, transakcja, przemianujEtykietePliku)
		if err != nil {
			return err
		}
		wynik, err := przyZasobach.ExecContext(ctx, nowa, stara, KontoOperatora(ctx))
		if err != nil {
			return fmt.Errorf("dane: nie można zmienić nazwy etykiety %q przy zasobach: %w", stara, err)
		}
		liczba, err := wynik.RowsAffected()
		if err != nil {
			return fmt.Errorf("dane: nie można policzyć zasobów zmienionej etykiety: %w", err)
		}
		dotkniete = int(liczba)

		return r.przemianujWpisSlownika(ctx, transakcja, stara, nowa)
	})
	if err != nil {
		return 0, err
	}
	return dotkniete, nil
}

func (r *repozytoriumBiblioteki) przemianujWpisSlownika(ctx context.Context, transakcja *sql.Tx,
	stara, nowa string) error {

	docelowy, err := r.zapytania.wTransakcji(ctx, transakcja, wpisSlownikaWKoncie)
	if err != nil {
		return err
	}
	var jeden int
	err = docelowy.QueryRowContext(ctx, nowa, KontoOperatora(ctx)).Scan(&jeden)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("dane: nieczytelny wpis słownika %q: %w", nowa, err)
	}
	if err == nil {
		zrodlowy, err := r.zapytania.wTransakcji(ctx, transakcja, usunEtykieteSlownikaZapis)
		if err != nil {
			return err
		}
		if _, err := zrodlowy.ExecContext(ctx, stara, KontoOperatora(ctx)); err != nil {
			return fmt.Errorf("dane: nie można usunąć etykiety %q ze słownika: %w", stara, err)
		}
		return nil
	}
	wSlowniku, err := r.zapytania.wTransakcji(ctx, transakcja, przemianujEtykieteSlownika)
	if err != nil {
		return err
	}
	if _, err := wSlowniku.ExecContext(ctx, nowa, stara, KontoOperatora(ctx)); err != nil {
		if czyKolizja(err) {
			return fmt.Errorf("dane: etykieta %q należy do słownika innego konta: %w", nowa, ErrKolizjaWiersza)
		}
		return fmt.Errorf("dane: nie można zmienić nazwy etykiety %q w słowniku: %w", stara, err)
	}
	return nil
}

func (r *repozytoriumBiblioteki) UsunEtykieteZeSlownika(ctx context.Context, nazwa string) (int, error) {
	zdjete := 0
	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		przyZasobach, err := r.zapytania.wTransakcji(ctx, transakcja, usunEtykietePlikow)
		if err != nil {
			return err
		}
		wynik, err := przyZasobach.ExecContext(ctx, nazwa, KontoOperatora(ctx))
		if err != nil {
			return fmt.Errorf("dane: nie można zdjąć etykiety %q z zasobów: %w", nazwa, err)
		}
		liczba, err := wynik.RowsAffected()
		if err != nil {
			return fmt.Errorf("dane: nie można policzyć zasobów zdjętej etykiety: %w", err)
		}
		zdjete = int(liczba)

		wSlowniku, err := r.zapytania.wTransakcji(ctx, transakcja, usunEtykieteSlownikaZapis)
		if err != nil {
			return err
		}
		if _, err := wSlowniku.ExecContext(ctx, nazwa, KontoOperatora(ctx)); err != nil {
			return fmt.Errorf("dane: nie można usunąć etykiety %q ze słownika: %w", nazwa, err)
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return zdjete, nil
}

func (r *repozytoriumBiblioteki) UstawRelacjeTezaurusa(ctx context.Context,
	zrodlo, cel, rodzaj string, zdejmij bool) (bool, error) {

	zapytanie := wstawRelacjeTezaurusa
	if zdejmij {
		zapytanie = usunRelacjeTezaurusa
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return false, err
	}
	if _, err := polecenie.ExecContext(ctx, zrodlo, cel, rodzaj); err != nil {
		return false, fmt.Errorf("dane: nie można zapisać relacji tezaurusa %q → %q: %w",
			zrodlo, cel, err)
	}
	return !zdejmij, nil
}

func (r *repozytoriumBiblioteki) RelacjeTezaurusa(ctx context.Context) ([]RelacjaTezaurusaBiblioteki, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, wykazRelacjiTezaurusa)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać tezaurusa: %w", err)
	}
	defer wiersze.Close()

	lista := []RelacjaTezaurusaBiblioteki{}
	for wiersze.Next() {
		var relacja RelacjaTezaurusaBiblioteki
		if err := wiersze.Scan(&relacja.Zrodlo, &relacja.Cel, &relacja.Rodzaj,
			&relacja.Utworzono); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz tezaurusa: %w", err)
		}
		lista = append(lista, relacja)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt tezaurusa: %w", err)
	}
	return lista, nil
}

func (r *repozytoriumBiblioteki) uzycieEtykiety(ctx context.Context, nazwa string) (int, error) {
	var liczba int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM etykieta_pliku_biblioteki
		  WHERE etykieta_pliku_biblioteki.etykieta = ? AND `+warunekKontaEtykietyPliku,
		nazwa, KontoOperatora(ctx)).Scan(&liczba)
	if err != nil {
		return 0, fmt.Errorf("dane: nie można policzyć użycia etykiety %q: %w", nazwa, err)
	}
	return liczba, nil
}

func wskazaniaKonta(zapytanie string, kontoID int64) []any {
	lista := []any{}
	for i := strings.Count(zapytanie, WarunekKonta); i > 0; i-- {
		lista = append(lista, kontoID)
	}
	return lista
}

func odczytajEtykieteSlownika(wiersz skaner) (EtykietaSlownikaBiblioteki, error) {
	var etykieta EtykietaSlownikaBiblioteki
	var barwa sql.NullString
	if err := wiersz.Scan(&etykieta.Nazwa, &barwa, &etykieta.Utworzono,
		&etykieta.LiczbaPlikow); err != nil {
		return EtykietaSlownikaBiblioteki{}, err
	}
	etykieta.Barwa = tekstZKolumny(barwa)
	return etykieta, nil
}
