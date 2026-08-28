// Odpowiedzialność pliku: zapis nadań dostępu. Zapis dotyka dwóch tabel
// (`nadanie_dostepu`, `korzen_nadania`) i porządku całego zbioru nadań okna,
// więc idzie w jednej transakcji.
package dane

import (
	"context"
	"database/sql"
	"fmt"
)

const (
	wstawNadanie = `INSERT INTO nadanie_dostepu
	                (okno_komunikacji_id, punkt_dostepu_id, tryb, kolejnosc, glowne, aktywne,
	                 identyfikator_zewnetrzny)
	                VALUES (?, ?, ?, ?, ?, ?, ?)`

	aktualizujNadanie = `UPDATE nadanie_dostepu
	                     SET punkt_dostepu_id = ?, tryb = ?, kolejnosc = ?, aktywne = ?,
	                         identyfikator_zewnetrzny = ?
	                     WHERE id = ?`

	usunNadanie = `DELETE FROM nadanie_dostepu WHERE id = ?`

	zdejmijGlowneOkna = `UPDATE nadanie_dostepu SET glowne = 0
	                     WHERE okno_komunikacji_id = ? AND glowne = 1 AND id <> ?`

	ustawGlowneNadanie = `UPDATE nadanie_dostepu SET glowne = 1 WHERE id = ?`

	nastepnaKolejnoscNadania = `SELECT COALESCE(MAX(kolejnosc), 0) + 1 FROM nadanie_dostepu
	                            WHERE okno_komunikacji_id = ?`

	liczbaNadanOkna = `SELECT COUNT(*) FROM nadanie_dostepu WHERE okno_komunikacji_id = ?`

	oknoNadania = `SELECT okno_komunikacji_id FROM nadanie_dostepu WHERE id = ?`
)

// Dodaj wpisuje nadanie do zbioru okna wraz z zawężeniem korzeni punktu dostępu, wyznaczając miejsce w kolejności.
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
		wynik, err := polecenie.ExecContext(ctx, nadanie.OknoKomunikacjiID, nadanie.PunktDostepuID,
			tryb, kolejnosc, liczbaLogiczna(glowne), liczbaLogiczna(nadanie.Aktywne),
			tekstDoKolumny(nadanie.IdentyfikatorZewnetrzny))
		if err != nil {
			return fmt.Errorf("dane: nie można nadać dostępu oknu %d: %w",
				nadanie.OknoKomunikacjiID, err)
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

// Aktualizuj zapisuje zmienione nadanie. Oznaczenia głównego nie zmienia —
// należy ono do zbioru okna, nie do pojedynczego wiersza, więc ma osobną metodę.
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
			nadanie.ID)
		if err != nil {
			return fmt.Errorf("dane: nie można zapisać nadania %d: %w", nadanie.ID, err)
		}
		if err := sprawdzTrafienie(wynik, "nadanie_dostepu", nadanie.ID); err != nil {
			return err
		}
		return zapiszKorzenie(ctx, r.zapytania, transakcja, usunKorzenieNadania, wstawKorzenNadania,
			nadanie.ID, nadanie.Korzenie, "nadania")
	})
}

// OznaczGlowne przenosi oznaczenie głównego na wskazane nadanie. Zdjęcie
// poprzedniego i nadanie nowego idą w jednej transakcji, bo indeks częściowy
// bazy dopuszcza najwyżej jedno główne nadanie okna.
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
		wynik, err := polecenie.ExecContext(ctx, id)
		if err != nil {
			return fmt.Errorf("dane: nie można oznaczyć nadania %d jako głównego: %w", id, err)
		}
		return sprawdzTrafienie(wynik, "nadanie_dostepu", id)
	})
}

// Usun cofa nadanie dostępu; korzenie nadania kasuje kaskada więzu klucza obcego zdefiniowana w schemacie bazy.
func (r *repozytoriumNadan) Usun(ctx context.Context, id int64) error {
	polecenie, err := r.zapytania.przygotuj(ctx, usunNadanie)
	if err != nil {
		return err
	}
	wynik, err := polecenie.ExecContext(ctx, id)
	if err != nil {
		return fmt.Errorf("dane: nie można cofnąć nadania %d: %w", id, err)
	}
	return sprawdzTrafienie(wynik, "nadanie_dostepu", id)
}
