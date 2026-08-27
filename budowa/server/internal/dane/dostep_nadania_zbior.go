// Odpowiedzialność pliku: reguły zbioru nadań jednego okna — miejsce nowego
// nadania w kolejności, oznaczenie głównego oraz sprawdzenie, czy zawężenie
// korzeni mieści się w obszarze punktu.
package dane

import (
	"context"
	"database/sql"
	"fmt"
)

// miejsceWZbiorze ustala kolejność nowego nadania w zbiorze oraz rozstrzyga, czy zostaje ono nadaniem głównym okna.
func (r *repozytoriumNadan) miejsceWZbiorze(ctx context.Context, transakcja *sql.Tx,
	nadanie Nadanie) (int, bool, error) {

	kolejnosc := nadanie.Kolejnosc
	if kolejnosc <= 0 {
		polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, nastepnaKolejnoscNadania)
		if err != nil {
			return 0, false, err
		}
		if err := polecenie.QueryRowContext(ctx, nadanie.OknoKomunikacjiID).Scan(&kolejnosc); err != nil {
			return 0, false, fmt.Errorf("dane: nie można ustalić kolejności nadania okna %d: %w",
				nadanie.OknoKomunikacjiID, err)
		}
	}
	if nadanie.Glowne {
		return kolejnosc, true, nil
	}
	polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, liczbaNadanOkna)
	if err != nil {
		return 0, false, err
	}
	var liczba int
	if err := polecenie.QueryRowContext(ctx, nadanie.OknoKomunikacjiID).Scan(&liczba); err != nil {
		return 0, false, fmt.Errorf("dane: nie można policzyć nadań okna %d: %w",
			nadanie.OknoKomunikacjiID, err)
	}
	return kolejnosc, liczba == 0, nil
}

// zdejmijGlowne kasuje oznaczenie głównego z pozostałych nadań okna. Wywoływane
// przed nadaniem oznaczenia nowemu wierszowi — indeks częściowy bazy dopuszcza
// najwyżej jedno główne nadanie okna.
func (r *repozytoriumNadan) zdejmijGlowne(ctx context.Context, transakcja *sql.Tx,
	oknoID, pomijaneID int64) error {

	polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, zdejmijGlowneOkna)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, oknoID, pomijaneID); err != nil {
		return fmt.Errorf("dane: nie można zdjąć oznaczenia głównego nadania okna %d: %w", oknoID, err)
	}
	return nil
}

// oknoNadania odczytuje okno, do którego należy nadanie. Brak wiersza jest
// sygnałem ErrBrakWiersza, nie awarią odczytu.
func (r *repozytoriumNadan) oknoNadania(ctx context.Context, transakcja *sql.Tx,
	id int64) (int64, error) {

	polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, oknoNadania)
	if err != nil {
		return 0, err
	}
	var oknoID int64
	err = polecenie.QueryRowContext(ctx, id).Scan(&oknoID)
	if err == sql.ErrNoRows {
		return 0, fmt.Errorf("dane: nadanie dostępu %d nie istnieje: %w", id, ErrBrakWiersza)
	}
	if err != nil {
		return 0, fmt.Errorf("dane: nie można odczytać okna nadania %d: %w", id, err)
	}
	return oknoID, nil
}

// sprawdzKorzenie pilnuje, żeby zawężenie korzeni nadania mieściło się w obszarze
// punktu. Lista pusta znaczy „komplet korzeni punktu" i nie wymaga sprawdzenia.
func (r *repozytoriumNadan) sprawdzKorzenie(ctx context.Context, nadanie Nadanie) error {
	if len(uporzadkujKorzenie(nadanie.Korzenie)) == 0 {
		return nil
	}
	korzeniePunktu, err := wczytajKorzenie(ctx, r.zapytania, listaKorzeniPunktu,
		nadanie.PunktDostepuID, "punktu dostępu")
	if err != nil {
		return err
	}
	return sprawdzZawezenieKorzeni(korzeniePunktu, nadanie.Korzenie)
}
