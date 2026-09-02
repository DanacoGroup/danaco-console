// Odpowiedzialność pliku: zestawy tematyczne źródeł i wątki tematyczne notatek, wraz z usunięciem źródła z wykazu okna przeglądania.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type ZestawZrodel struct {
	ID             int64
	Kod            string
	Okno           string
	Nazwa          string
	Utworzono      string
	Zaktualizowano string
}

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

	// Warunek przy DO UPDATE zostawia wiersz cudzego konta nietknięty;
	// zapis kończy się wtedy ErrKolizjaWiersza (migracja 484).
	zapiszZestawZrodel = `INSERT INTO zestaw_zrodel_przegladania
	                      (identyfikator_zewnetrzny, okno, nazwa, konto_id)
	                      VALUES (?, ?, ?, ` + WskazanieKonta + `)
	                      ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                          nazwa = excluded.nazwa,
	                          zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                      WHERE ` + WarunekKonta

	pobierzZestawZrodel = `SELECT ` + kolumnyZestawuZrodel + `
	                       FROM zestaw_zrodel_przegladania
	                       WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	listaZestawowZrodel = `SELECT ` + kolumnyZestawuZrodel + `
	                       FROM zestaw_zrodel_przegladania
	                       WHERE okno = ? AND ` + WarunekKonta + `
	                       ORDER BY utworzono DESC, id DESC LIMIT ?`

	usunZestawZrodel = `DELETE FROM zestaw_zrodel_przegladania
	                    WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	odepnijZrodlaZestawu = `UPDATE zrodlo_przegladania SET grupa = NULL
	                        WHERE grupa = ? AND ` + WarunekKonta

	przypiszZrodloDoZestawu = `UPDATE zrodlo_przegladania SET grupa = ?
	                           WHERE identyfikator_zewnetrzny = ? AND okno = ? AND ` + WarunekKonta

	kodyZrodelZestawu = `SELECT identyfikator_zewnetrzny FROM zrodlo_przegladania
	                     WHERE grupa = ? AND ` + WarunekKonta + `
	                     ORDER BY utworzono DESC, id DESC`

	usunZrodloPrzegladania = `DELETE FROM zrodlo_przegladania
	                          WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	kolumnyWatkuNotatek = `id, identyfikator_zewnetrzny, okno, nazwa, utworzono, zaktualizowano`

	zapiszWatekNotatek = `INSERT INTO watek_notatek_przegladania
	                      (identyfikator_zewnetrzny, okno, nazwa, konto_id)
	                      VALUES (?, ?, ?, ` + WskazanieKonta + `)
	                      ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                          nazwa = excluded.nazwa,
	                          zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                      WHERE ` + WarunekKonta

	pobierzWatekNotatek = `SELECT ` + kolumnyWatkuNotatek + `
	                       FROM watek_notatek_przegladania
	                       WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	listaWatkowNotatek = `SELECT ` + kolumnyWatkuNotatek + `
	                      FROM watek_notatek_przegladania
	                      WHERE okno = ? AND ` + WarunekKonta + `
	                      ORDER BY utworzono DESC, id DESC LIMIT ?`

	usunWatekNotatek = `DELETE FROM watek_notatek_przegladania
	                    WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	odepnijNotatkiWatku = `UPDATE notatka_przegladania SET watek = NULL
	                       WHERE watek = ? AND ` + WarunekKonta

	przypiszNotatkeDoWatku = `UPDATE notatka_przegladania SET watek = ?
	                          WHERE identyfikator_zewnetrzny = ? AND okno = ? AND ` + WarunekKonta

	kodyNotatekWatku = `SELECT identyfikator_zewnetrzny FROM notatka_przegladania
	                    WHERE watek = ? AND ` + WarunekKonta + `
	                    ORDER BY utworzono DESC, id DESC`
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
	wynik, err := polecenie.ExecContext(ctx, zestaw.Kod, zestaw.Okno, zestaw.Nazwa,
		KontoOperatora(ctx), KontoOperatora(ctx))
	if err != nil {
		return ZestawZrodel{}, fmt.Errorf("dane: nie można zapisać zestawu źródeł %q: %w", zestaw.Kod, err)
	}
	if err := sprawdzTrafienieZapisu(wynik, "zestaw źródeł", zestaw.Kod); err != nil {
		return ZestawZrodel{}, err
	}
	return r.ZestawZrodel(ctx, zestaw.Kod)
}

func (r *repozytoriumPrzegladania) ZestawZrodel(ctx context.Context, kod string) (ZestawZrodel, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzZestawZrodel)
	if err != nil {
		return ZestawZrodel{}, err
	}
	var zestaw ZestawZrodel
	err = polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)).Scan(&zestaw.ID, &zestaw.Kod,
		&zestaw.Okno, &zestaw.Nazwa, &zestaw.Utworzono, &zestaw.Zaktualizowano)
	if errors.Is(err, sql.ErrNoRows) {
		return ZestawZrodel{}, ErrBrakWiersza
	}
	if err != nil {
		return ZestawZrodel{}, fmt.Errorf("dane: nieczytelny zestaw źródeł %q: %w", kod, err)
	}
	return zestaw, nil
}

func (r *repozytoriumPrzegladania) ZestawyZrodel(ctx context.Context, okno string, limit int) ([]ZestawZrodel, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaZestawowZrodel)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno, KontoOperatora(ctx), granicaWykazu(limit))
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

// UsunZestawZrodel zdejmuje zestaw i odpina od niego źródła; źródła zostają w wykazie okna.
func (r *repozytoriumPrzegladania) UsunZestawZrodel(ctx context.Context, kod string) (bool, error) {
	odpiecie, err := r.zapytania.przygotuj(ctx, odepnijZrodlaZestawu)
	if err != nil {
		return false, err
	}
	if _, err := odpiecie.ExecContext(ctx, kod, KontoOperatora(ctx)); err != nil {
		return false, fmt.Errorf("dane: nie można odpiąć źródeł zestawu %q: %w", kod, err)
	}
	return r.usunWiersz(ctx, usunZestawZrodel, kod, "zestaw źródeł", KontoOperatora(ctx))
}

// PrzypiszZrodloDoZestawu wiąże źródło z zestawem tematycznym albo zdejmuje wiązanie.
func (r *repozytoriumPrzegladania) PrzypiszZrodloDoZestawu(ctx context.Context, okno, zrodlo, zestaw string) error {
	polecenie, err := r.zapytania.przygotuj(ctx, przypiszZrodloDoZestawu)
	if err != nil {
		return err
	}
	var wartosc any
	if zestaw != "" {
		wartosc = zestaw
	}
	wynik, err := polecenie.ExecContext(ctx, wartosc, zrodlo, okno, KontoOperatora(ctx))
	if err != nil {
		return fmt.Errorf("dane: nie można przypisać źródła %q do zestawu: %w", zrodlo, err)
	}
	zmienione, err := wynik.RowsAffected()
	if err != nil {
		return fmt.Errorf("dane: nieznany skutek przypisania źródła %q: %w", zrodlo, err)
	}
	if zmienione == 0 {
		return fmt.Errorf("dane: źródło %q nie należy do okna %q: %w", zrodlo, okno, ErrBrakWiersza)
	}
	return nil
}

func (r *repozytoriumPrzegladania) KodyZrodelZestawu(ctx context.Context, zestaw string) ([]string, error) {
	return r.kody(ctx, kodyZrodelZestawu, zestaw, "zestawu źródeł", KontoOperatora(ctx))
}

// UsunZrodlo zdejmuje źródło z wykazu okna (komenda browser.source.remove).
func (r *repozytoriumPrzegladania) UsunZrodlo(ctx context.Context, kod string) (bool, error) {
	return r.usunWiersz(ctx, usunZrodloPrzegladania, kod, "źródło przeglądania", KontoOperatora(ctx))
}

// ZapiszWatekNotatek zakłada wątek tematyczny notatek albo zmienia nazwę wątku zapisanego wcześniej.
func (r *repozytoriumPrzegladania) ZapiszWatekNotatek(ctx context.Context, watek WatekNotatek) (WatekNotatek, error) {
	if watek.Kod == "" || watek.Okno == "" || watek.Nazwa == "" {
		return WatekNotatek{}, fmt.Errorf("dane: wątek notatek bez identyfikatora, okna albo nazwy")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszWatekNotatek)
	if err != nil {
		return WatekNotatek{}, err
	}
	wynik, err := polecenie.ExecContext(ctx, watek.Kod, watek.Okno, watek.Nazwa,
		KontoOperatora(ctx), KontoOperatora(ctx))
	if err != nil {
		return WatekNotatek{}, fmt.Errorf("dane: nie można zapisać wątku notatek %q: %w", watek.Kod, err)
	}
	if err := sprawdzTrafienieZapisu(wynik, "wątek notatek", watek.Kod); err != nil {
		return WatekNotatek{}, err
	}
	return r.WatekNotatek(ctx, watek.Kod)
}

func (r *repozytoriumPrzegladania) WatekNotatek(ctx context.Context, kod string) (WatekNotatek, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzWatekNotatek)
	if err != nil {
		return WatekNotatek{}, err
	}
	var watek WatekNotatek
	err = polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)).Scan(&watek.ID, &watek.Kod,
		&watek.Okno, &watek.Nazwa, &watek.Utworzono, &watek.Zaktualizowano)
	if errors.Is(err, sql.ErrNoRows) {
		return WatekNotatek{}, ErrBrakWiersza
	}
	if err != nil {
		return WatekNotatek{}, fmt.Errorf("dane: nieczytelny wątek notatek %q: %w", kod, err)
	}
	return watek, nil
}

func (r *repozytoriumPrzegladania) WatkiNotatek(ctx context.Context, okno string, limit int) ([]WatekNotatek, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaWatkowNotatek)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno, KontoOperatora(ctx), granicaWykazu(limit))
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

// UsunWatekNotatek zdejmuje wątek tematyczny notatek i odpina od niego notatki.
func (r *repozytoriumPrzegladania) UsunWatekNotatek(ctx context.Context, kod string) (bool, error) {
	odpiecie, err := r.zapytania.przygotuj(ctx, odepnijNotatkiWatku)
	if err != nil {
		return false, err
	}
	if _, err := odpiecie.ExecContext(ctx, kod, KontoOperatora(ctx)); err != nil {
		return false, fmt.Errorf("dane: nie można odpiąć notatek wątku %q: %w", kod, err)
	}
	return r.usunWiersz(ctx, usunWatekNotatek, kod, "wątek notatek", KontoOperatora(ctx))
}

func (r *repozytoriumPrzegladania) PrzypiszNotatkeDoWatku(ctx context.Context, okno, notatka, watek string) error {
	polecenie, err := r.zapytania.przygotuj(ctx, przypiszNotatkeDoWatku)
	if err != nil {
		return err
	}
	var wartosc any
	if watek != "" {
		wartosc = watek
	}
	if _, err := polecenie.ExecContext(ctx, wartosc, notatka, okno, KontoOperatora(ctx)); err != nil {
		return fmt.Errorf("dane: nie można przypisać notatki %q do wątku: %w", notatka, err)
	}
	return nil
}

func (r *repozytoriumPrzegladania) KodyNotatekWatku(ctx context.Context, watek string) ([]string, error) {
	return r.kody(ctx, kodyNotatekWatku, watek, "wątku notatek", KontoOperatora(ctx))
}

func (r *repozytoriumPrzegladania) kody(ctx context.Context, zapytanie, wskazanie, nazwa string,
	dalsze ...any) ([]string, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, append([]any{wskazanie}, dalsze...)...)
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
