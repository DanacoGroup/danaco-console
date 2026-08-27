// Odpowiedzialność pliku: zestawy tematyczne źródeł i wątki tematyczne notatek, wraz z usunięciem źródła z wykazu okna przeglądania.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// ZestawZrodel to wiersz tabeli zestaw_zrodel_przegladania, niosący nazwę i kod zestawu tematycznego źródeł okna.
type ZestawZrodel struct {
	ID             int64
	Kod            string
	Okno           string
	Nazwa          string
	Utworzono      string
	Zaktualizowano string
}

// WatekNotatek to wiersz tabeli watek_notatek_przegladania, niosący nazwę i kod wątku tematycznego notatek okna.
type WatekNotatek struct {
	ID             int64
	Kod            string
	Okno           string
	Nazwa          string
	Utworzono      string
	Zaktualizowano string
}

const (
	kolumnyZestawuZrodel = `id, identyfikator_zewnetrzny, okno, nazwa, utworzono, zaktualizowano`

	zapiszZestawZrodel = `INSERT INTO zestaw_zrodel_przegladania
	                      (identyfikator_zewnetrzny, okno, nazwa)
	                      VALUES (?, ?, ?)
	                      ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                          nazwa = excluded.nazwa,
	                          zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')`

	pobierzZestawZrodel = `SELECT ` + kolumnyZestawuZrodel + `
	                       FROM zestaw_zrodel_przegladania WHERE identyfikator_zewnetrzny = ?`

	listaZestawowZrodel = `SELECT ` + kolumnyZestawuZrodel + `
	                       FROM zestaw_zrodel_przegladania WHERE okno = ?
	                       ORDER BY utworzono DESC, id DESC LIMIT ?`

	usunZestawZrodel = `DELETE FROM zestaw_zrodel_przegladania WHERE identyfikator_zewnetrzny = ?`

	odepnijZrodlaZestawu = `UPDATE zrodlo_przegladania SET grupa = NULL WHERE grupa = ?`

	przypiszZrodloDoZestawu = `UPDATE zrodlo_przegladania SET grupa = ?
	                           WHERE identyfikator_zewnetrzny = ? AND okno = ?`

	kodyZrodelZestawu = `SELECT identyfikator_zewnetrzny FROM zrodlo_przegladania
	                     WHERE grupa = ? ORDER BY utworzono DESC, id DESC`

	usunZrodloPrzegladania = `DELETE FROM zrodlo_przegladania WHERE identyfikator_zewnetrzny = ?`

	kolumnyWatkuNotatek = `id, identyfikator_zewnetrzny, okno, nazwa, utworzono, zaktualizowano`

	zapiszWatekNotatek = `INSERT INTO watek_notatek_przegladania
	                      (identyfikator_zewnetrzny, okno, nazwa)
	                      VALUES (?, ?, ?)
	                      ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                          nazwa = excluded.nazwa,
	                          zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')`

	pobierzWatekNotatek = `SELECT ` + kolumnyWatkuNotatek + `
	                       FROM watek_notatek_przegladania WHERE identyfikator_zewnetrzny = ?`

	listaWatkowNotatek = `SELECT ` + kolumnyWatkuNotatek + `
	                      FROM watek_notatek_przegladania WHERE okno = ?
	                      ORDER BY utworzono DESC, id DESC LIMIT ?`

	usunWatekNotatek = `DELETE FROM watek_notatek_przegladania WHERE identyfikator_zewnetrzny = ?`

	odepnijNotatkiWatku = `UPDATE notatka_przegladania SET watek = NULL WHERE watek = ?`

	przypiszNotatkeDoWatku = `UPDATE notatka_przegladania SET watek = ?
	                          WHERE identyfikator_zewnetrzny = ? AND okno = ?`

	kodyNotatekWatku = `SELECT identyfikator_zewnetrzny FROM notatka_przegladania
	                    WHERE watek = ? ORDER BY utworzono DESC, id DESC`
)

// ZapiszZestawZrodel zakłada zestaw tematyczny źródeł albo zmienia nazwę zestawu zapisanego wcześniej.
func (r *repozytoriumPrzegladania) ZapiszZestawZrodel(ctx context.Context, zestaw ZestawZrodel) (ZestawZrodel, error) {
	if zestaw.Kod == "" || zestaw.Okno == "" || zestaw.Nazwa == "" {
		return ZestawZrodel{}, fmt.Errorf("dane: zestaw źródeł bez identyfikatora, okna albo nazwy")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszZestawZrodel)
	if err != nil {
		return ZestawZrodel{}, err
	}
	if _, err := polecenie.ExecContext(ctx, zestaw.Kod, zestaw.Okno, zestaw.Nazwa); err != nil {
		return ZestawZrodel{}, fmt.Errorf("dane: nie można zapisać zestawu źródeł %q: %w", zestaw.Kod, err)
	}
	return r.ZestawZrodel(ctx, zestaw.Kod)
}

// ZestawZrodel oddaje zestaw tematyczny źródeł o wskazanym kodzie zewnętrznym wraz z jego pełną nazwą.
func (r *repozytoriumPrzegladania) ZestawZrodel(ctx context.Context, kod string) (ZestawZrodel, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzZestawZrodel)
	if err != nil {
		return ZestawZrodel{}, err
	}
	var zestaw ZestawZrodel
	err = polecenie.QueryRowContext(ctx, kod).Scan(&zestaw.ID, &zestaw.Kod, &zestaw.Okno,
		&zestaw.Nazwa, &zestaw.Utworzono, &zestaw.Zaktualizowano)
	if errors.Is(err, sql.ErrNoRows) {
		return ZestawZrodel{}, ErrBrakWiersza
	}
	if err != nil {
		return ZestawZrodel{}, fmt.Errorf("dane: nieczytelny zestaw źródeł %q: %w", kod, err)
	}
	return zestaw, nil
}

// ZestawyZrodel oddaje wszystkie zestawy tematyczne źródeł danego okna, od najnowszego do najstarszego.
func (r *repozytoriumPrzegladania) ZestawyZrodel(ctx context.Context, okno string, limit int) ([]ZestawZrodel, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaZestawowZrodel)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno, granicaWykazu(limit))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać zestawów źródeł okna %q: %w", okno, err)
	}
	defer wiersze.Close()

	lista := []ZestawZrodel{}
	for wiersze.Next() {
		var zestaw ZestawZrodel
		if err := wiersze.Scan(&zestaw.ID, &zestaw.Kod, &zestaw.Okno, &zestaw.Nazwa,
			&zestaw.Utworzono, &zestaw.Zaktualizowano); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz zestawów źródeł: %w", err)
		}
		lista = append(lista, zestaw)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt zestawów źródeł: %w", err)
	}
	return lista, nil
}

// UsunZestawZrodel zdejmuje zestaw i odpina od niego źródła. Źródła zostają
// w wykazie okna — zestaw jest ich grupowaniem, nie właścicielem.
func (r *repozytoriumPrzegladania) UsunZestawZrodel(ctx context.Context, kod string) (bool, error) {
	odpiecie, err := r.zapytania.przygotuj(ctx, odepnijZrodlaZestawu)
	if err != nil {
		return false, err
	}
	if _, err := odpiecie.ExecContext(ctx, kod); err != nil {
		return false, fmt.Errorf("dane: nie można odpiąć źródeł zestawu %q: %w", kod, err)
	}
	return r.usunWiersz(ctx, usunZestawZrodel, kod, "zestaw źródeł")
}

// PrzypiszZrodloDoZestawu wiąże wskazane źródło z zestawem tematycznym albo zdejmuje istniejące wiązanie.
func (r *repozytoriumPrzegladania) PrzypiszZrodloDoZestawu(ctx context.Context, okno, zrodlo, zestaw string) error {
	polecenie, err := r.zapytania.przygotuj(ctx, przypiszZrodloDoZestawu)
	if err != nil {
		return err
	}
	var wartosc any
	if zestaw != "" {
		wartosc = zestaw
	}
	if _, err := polecenie.ExecContext(ctx, wartosc, zrodlo, okno); err != nil {
		return fmt.Errorf("dane: nie można przypisać źródła %q do zestawu: %w", zrodlo, err)
	}
	return nil
}

// KodyZrodelZestawu oddaje skład zestawu — identyfikatory źródeł do niego
// należących, liczone z kolumny, a nie z drugiej listy.
func (r *repozytoriumPrzegladania) KodyZrodelZestawu(ctx context.Context, zestaw string) ([]string, error) {
	return r.kody(ctx, kodyZrodelZestawu, zestaw, "zestawu źródeł")
}

// UsunZrodlo zdejmuje wskazane źródło z wykazu okna, obsługując komendę browser.source.remove w całości.
func (r *repozytoriumPrzegladania) UsunZrodlo(ctx context.Context, kod string) (bool, error) {
	return r.usunWiersz(ctx, usunZrodloPrzegladania, kod, "źródło przeglądania")
}

// ZapiszWatekNotatek zakłada nowy wątek tematyczny notatek albo zmienia nazwę wątku już zapisanego wcześniej.
func (r *repozytoriumPrzegladania) ZapiszWatekNotatek(ctx context.Context, watek WatekNotatek) (WatekNotatek, error) {
	if watek.Kod == "" || watek.Okno == "" || watek.Nazwa == "" {
		return WatekNotatek{}, fmt.Errorf("dane: wątek notatek bez identyfikatora, okna albo nazwy")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszWatekNotatek)
	if err != nil {
		return WatekNotatek{}, err
	}
	if _, err := polecenie.ExecContext(ctx, watek.Kod, watek.Okno, watek.Nazwa); err != nil {
		return WatekNotatek{}, fmt.Errorf("dane: nie można zapisać wątku notatek %q: %w", watek.Kod, err)
	}
	return r.WatekNotatek(ctx, watek.Kod)
}

// WatekNotatek oddaje wątek tematyczny notatek o wskazanym kodzie zewnętrznym, razem z jego pełną nazwą.
func (r *repozytoriumPrzegladania) WatekNotatek(ctx context.Context, kod string) (WatekNotatek, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzWatekNotatek)
	if err != nil {
		return WatekNotatek{}, err
	}
	var watek WatekNotatek
	err = polecenie.QueryRowContext(ctx, kod).Scan(&watek.ID, &watek.Kod, &watek.Okno,
		&watek.Nazwa, &watek.Utworzono, &watek.Zaktualizowano)
	if errors.Is(err, sql.ErrNoRows) {
		return WatekNotatek{}, ErrBrakWiersza
	}
	if err != nil {
		return WatekNotatek{}, fmt.Errorf("dane: nieczytelny wątek notatek %q: %w", kod, err)
	}
	return watek, nil
}

// WatkiNotatek oddaje wszystkie wątki tematyczne notatek okna operacyjnego, od najnowszego do najstarszego.
func (r *repozytoriumPrzegladania) WatkiNotatek(ctx context.Context, okno string, limit int) ([]WatekNotatek, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaWatkowNotatek)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno, granicaWykazu(limit))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać wątków notatek okna %q: %w", okno, err)
	}
	defer wiersze.Close()

	lista := []WatekNotatek{}
	for wiersze.Next() {
		var watek WatekNotatek
		if err := wiersze.Scan(&watek.ID, &watek.Kod, &watek.Okno, &watek.Nazwa,
			&watek.Utworzono, &watek.Zaktualizowano); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz wątków notatek: %w", err)
		}
		lista = append(lista, watek)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt wątków notatek: %w", err)
	}
	return lista, nil
}

// UsunWatekNotatek zdejmuje wskazany wątek tematyczny notatek i odpina od niego wszystkie powiązane notatki.
func (r *repozytoriumPrzegladania) UsunWatekNotatek(ctx context.Context, kod string) (bool, error) {
	odpiecie, err := r.zapytania.przygotuj(ctx, odepnijNotatkiWatku)
	if err != nil {
		return false, err
	}
	if _, err := odpiecie.ExecContext(ctx, kod); err != nil {
		return false, fmt.Errorf("dane: nie można odpiąć notatek wątku %q: %w", kod, err)
	}
	return r.usunWiersz(ctx, usunWatekNotatek, kod, "wątek notatek")
}

// PrzypiszNotatkeDoWatku wiąże wskazaną notatkę z wątkiem tematycznym albo zdejmuje istniejące wiązanie.
func (r *repozytoriumPrzegladania) PrzypiszNotatkeDoWatku(ctx context.Context, okno, notatka, watek string) error {
	polecenie, err := r.zapytania.przygotuj(ctx, przypiszNotatkeDoWatku)
	if err != nil {
		return err
	}
	var wartosc any
	if watek != "" {
		wartosc = watek
	}
	if _, err := polecenie.ExecContext(ctx, wartosc, notatka, okno); err != nil {
		return fmt.Errorf("dane: nie można przypisać notatki %q do wątku: %w", notatka, err)
	}
	return nil
}

// KodyNotatekWatku oddaje skład wątku, czyli identyfikatory notatek do niego należących, liczone z kolumny.
func (r *repozytoriumPrzegladania) KodyNotatekWatku(ctx context.Context, watek string) ([]string, error) {
	return r.kody(ctx, kodyNotatekWatku, watek, "wątku notatek")
}

// kody wykonuje odczyt jednej kolumny identyfikatorów. Pusty wykaz jest
// prawidłowym składem: zestaw założony i jeszcze niezapełniony istnieje.
func (r *repozytoriumPrzegladania) kody(ctx context.Context, zapytanie, wskazanie, nazwa string) ([]string, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, wskazanie)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać składu %s %q: %w", nazwa, wskazanie, err)
	}
	defer wiersze.Close()

	lista := []string{}
	for wiersze.Next() {
		var kod string
		if err := wiersze.Scan(&kod); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz składu %s %q: %w", nazwa, wskazanie, err)
		}
		lista = append(lista, kod)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt składu %s %q: %w", nazwa, wskazanie, err)
	}
	return lista, nil
}
