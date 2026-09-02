// Odpowiedzialność pliku: zapis nadań dostępu — dwie tabele (`nadanie_dostepu`,
// `korzen_nadania`) i porządek zbioru nadań okna w jednej transakcji.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var (
	wstawNadanie = `INSERT INTO nadanie_dostepu
	                (okno_komunikacji_id, punkt_dostepu_id, tryb, kolejnosc, glowne, aktywne,
	                 identyfikator_zewnetrzny)
	                SELECT ?, ?, ?, ?, ?, ?, ?
	                 WHERE ` + kontoOknaNadan + `
	                   AND EXISTS (SELECT 1 FROM punkt_dostepu WHERE id = ? AND ` + WarunekKonta + `)`

	aktualizujNadanie = `UPDATE nadanie_dostepu
	                     SET punkt_dostepu_id = ?, tryb = ?, kolejnosc = ?, aktywne = ?,
	                         identyfikator_zewnetrzny = ?
	                     WHERE id = ? AND ` + kontoNadania

	usunNadanie = `DELETE FROM nadanie_dostepu WHERE id = ? AND ` + kontoNadania

	zdejmijGlowneOkna = `UPDATE nadanie_dostepu SET glowne = 0
	                     WHERE okno_komunikacji_id = ? AND glowne = 1 AND id <> ? AND ` + kontoNadania

	ustawGlowneNadanie = `UPDATE nadanie_dostepu SET glowne = 1 WHERE id = ? AND ` + kontoNadania

	nastepnaKolejnoscNadania = `SELECT COALESCE(MAX(kolejnosc), 0) + 1 FROM nadanie_dostepu
	                            WHERE okno_komunikacji_id = ? AND ` + kontoNadania

	liczbaNadanOkna = `SELECT COUNT(*) FROM nadanie_dostepu
	                   WHERE okno_komunikacji_id = ? AND ` + kontoNadania

	oknoNadania = `SELECT okno_komunikacji_id FROM nadanie_dostepu WHERE id = ? AND ` + kontoNadania

	istnienieNadania = `SELECT 1 FROM nadanie_dostepu WHERE id = ?`
)

func (r *repozytoriumNadan) Dodaj(ctx context.Context, nadanie Nadanie) (int64, error) {
	tryb, err := trybDostepuNaBaze(nadanie.Tryb, "nadanie_dostepu.tryb")
	if err != nil {
		return 0, err
	}
	if err := r.sprawdzKorzenie(ctx, nadanie); err != nil {
		return 0, err
	}
	var id int64
	err = wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		kolejnosc, glowne, err := r.miejsceWZbiorze(ctx, transakcja, nadanie)
		if err != nil {
			return err
		}
		polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, wstawNadanie)
		if err != nil {
			return err
		}
		if glowne {
			if err := r.zdejmijGlowne(ctx, transakcja, nadanie.OknoKomunikacjiID, 0); err != nil {
				return err
			}
		}
		konto := KontoOperatora(ctx)
		wynik, err := polecenie.ExecContext(ctx, nadanie.OknoKomunikacjiID, nadanie.PunktDostepuID,
			tryb, kolejnosc, liczbaLogiczna(glowne), liczbaLogiczna(nadanie.Aktywne),
			tekstDoKolumny(nadanie.IdentyfikatorZewnetrzny),
			nadanie.OknoKomunikacjiID, konto, nadanie.PunktDostepuID, konto)
		if err != nil {
			return fmt.Errorf("dane: nie można nadać dostępu oknu %d: %w",
				nadanie.OknoKomunikacjiID, err)
		}
		if err := sprawdzTrafienieZapisu(wynik, "okno komunikacji nadania",
			fmt.Sprintf("%d", nadanie.OknoKomunikacjiID)); err != nil {
			return err
		}
		if id, err = wynik.LastInsertId(); err != nil {
			return fmt.Errorf("dane: nieznany identyfikator zapisanego nadania: %w", err)
		}
		return zapiszKorzenie(ctx, r.zapytania, transakcja, usunKorzenieNadania, wstawKorzenNadania,
			id, nadanie.Korzenie, "nadania")
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}

// Oznaczenie głównego należy do zbioru okna, nie do wiersza, więc Aktualizuj go nie zmienia.
func (r *repozytoriumNadan) Aktualizuj(ctx context.Context, nadanie Nadanie) error {
	tryb, err := trybDostepuNaBaze(nadanie.Tryb, "nadanie_dostepu.tryb")
	if err != nil {
		return err
	}
	if err := r.sprawdzKorzenie(ctx, nadanie); err != nil {
		return err
	}
	return wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, aktualizujNadanie)
		if err != nil {
			return err
		}
		wynik, err := polecenie.ExecContext(ctx, nadanie.PunktDostepuID, tryb, nadanie.Kolejnosc,
			liczbaLogiczna(nadanie.Aktywne), tekstDoKolumny(nadanie.IdentyfikatorZewnetrzny),
			nadanie.ID, KontoOperatora(ctx))
		if err != nil {
			return fmt.Errorf("dane: nie można zapisać nadania %d: %w", nadanie.ID, err)
		}
		if err := r.trafienieNadania(ctx, transakcja, wynik, nadanie.ID); err != nil {
			return err
		}
		return zapiszKorzenie(ctx, r.zapytania, transakcja, usunKorzenieNadania, wstawKorzenNadania,
			nadanie.ID, nadanie.Korzenie, "nadania")
	})
}

// Zdjęcie poprzedniego i nadanie nowego oznaczenia idą w jednej transakcji: indeks
// częściowy bazy dopuszcza najwyżej jedno główne nadanie okna.
func (r *repozytoriumNadan) OznaczGlowne(ctx context.Context, id int64) error {
	return wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		oknoID, err := r.oknoNadania(ctx, transakcja, id)
		if err != nil {
			return err
		}
		if err := r.zdejmijGlowne(ctx, transakcja, oknoID, id); err != nil {
			return err
		}
		polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, ustawGlowneNadanie)
		if err != nil {
			return err
		}
		wynik, err := polecenie.ExecContext(ctx, id, KontoOperatora(ctx))
		if err != nil {
			return fmt.Errorf("dane: nie można oznaczyć nadania %d jako głównego: %w", id, err)
		}
		return r.trafienieNadania(ctx, transakcja, wynik, id)
	})
}

// Korzenie nadania kasuje kaskada więzu klucza obcego ze schematu bazy.
func (r *repozytoriumNadan) Usun(ctx context.Context, id int64) error {
	return wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, usunNadanie)
		if err != nil {
			return err
		}
		wynik, err := polecenie.ExecContext(ctx, id, KontoOperatora(ctx))
		if err != nil {
			return fmt.Errorf("dane: nie można cofnąć nadania %d: %w", id, err)
		}
		return r.trafienieNadania(ctx, transakcja, wynik, id)
	})
}

// trafienieNadania odróżnia zapis zatrzymany przez granicę konta (ErrKolizjaWiersza)
// od zapisu w wiersz, którego nie ma (ErrBrakWiersza); silnik obu przypadków nie rozróżnia.
func (r *repozytoriumNadan) trafienieNadania(ctx context.Context, transakcja *sql.Tx,
	wynik sql.Result, id int64) error {

	zmienione, err := wynik.RowsAffected()
	if err != nil {
		return fmt.Errorf("dane: nieznana liczba zmienionych wierszy tabeli nadanie_dostepu: %w", err)
	}
	if zmienione > 0 {
		return nil
	}
	polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, istnienieNadania)
	if err != nil {
		return err
	}
	var obecny int
	err = polecenie.QueryRowContext(ctx, id).Scan(&obecny)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("dane: nadanie dostępu %d nie istnieje: %w", id, ErrBrakWiersza)
	}
	if err != nil {
		return fmt.Errorf("dane: nie można odczytać nadania %d: %w", id, err)
	}
	return fmt.Errorf("dane: nadanie dostępu %d należy do innego konta: %w", id, ErrKolizjaWiersza)
}
