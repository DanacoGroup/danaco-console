// Plik dokłada do modułu automatyk dwie czynności ułożenia pozycji kolejki:
// zmianę priorytetu oraz skierowanie pozycji do innej kolejki, bez zmiany
// stanu ani biegu naprawczego zlecenia.
package dane

import (
	"context"
	"database/sql"
	"fmt"
)

const (
	ustawKolejnoscPozycji = `UPDATE pozycja_kolejki
	                         SET kolejnosc = ?,
	                             zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                         WHERE id = ?`

	przeniesPozycjeDoKolejki = `UPDATE pozycja_kolejki
	                            SET kolejka_id = ?,
	                                kolejnosc = (SELECT COALESCE(MAX(kolejnosc) + 1, 0)
	                                             FROM pozycja_kolejki WHERE kolejka_id = ?),
	                                zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                            WHERE id = ?`
)

// UstawKolejnoscPozycji zmienia priorytet zlecenia w kolejce. Niższa liczba
// znaczy wcześniejsze wykonanie — tak samo jak przy zasilaniu kolejki, gdzie
// kolejność rośnie wraz z dokładaniem zleceń.
func (r *repozytoriumAutomatyk) UstawKolejnoscPozycji(ctx context.Context,
	pozycjaID int64, kolejnosc int) error {

	return wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		kolejkaID, obieg, err := polozeniePozycji(ctx, r.zapytania, transakcja, pozycjaID)
		if err != nil {
			return err
		}
		polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, ustawKolejnoscPozycji)
		if err != nil {
			return err
		}
		wynik, err := polecenie.ExecContext(ctx, kolejnosc, pozycjaID)
		if err != nil {
			return fmt.Errorf("dane: nie można zmienić priorytetu pozycji %d: %w", pozycjaID, err)
		}
		if err := sprawdzTrafienie(wynik, "pozycja_kolejki", pozycjaID); err != nil {
			return err
		}
		szczegoly := fmt.Sprintf("priorytet=%d", kolejnosc)
		return dopiszWpisDziennika(ctx, r.zapytania, transakcja, WpisDziennika{
			KolejkaID: kolejkaID, PozycjaID: &pozycjaID, Akcja: "zmiana_priorytetu",
			NumerObiegu: obieg, Szczegoly: &szczegoly,
		})
	})
}

// PrzeniesPozycje kieruje zlecenie do innej kolejki, na jej koniec. Ślad
// zostaje w dzienniku kolejki źródłowej, bo to z niej pozycja znika — inaczej
// przegląd kolejki pokazywałby zniknięcie bez przyczyny.
func (r *repozytoriumAutomatyk) PrzeniesPozycje(ctx context.Context,
	pozycjaID, kolejkaDocelowaID int64) error {

	return wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		zrodlowa, obieg, err := polozeniePozycji(ctx, r.zapytania, transakcja, pozycjaID)
		if err != nil {
			return err
		}
		if zrodlowa == kolejkaDocelowaID {
			return nil
		}
		polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, przeniesPozycjeDoKolejki)
		if err != nil {
			return err
		}
		wynik, err := polecenie.ExecContext(ctx, kolejkaDocelowaID, kolejkaDocelowaID, pozycjaID)
		if err != nil {
			return fmt.Errorf("dane: nie można skierować pozycji %d do kolejki %d: %w",
				pozycjaID, kolejkaDocelowaID, err)
		}
		if err := sprawdzTrafienie(wynik, "pozycja_kolejki", pozycjaID); err != nil {
			return err
		}
		szczegoly := fmt.Sprintf("kolejka_docelowa=%d", kolejkaDocelowaID)
		return dopiszWpisDziennika(ctx, r.zapytania, transakcja, WpisDziennika{
			KolejkaID: zrodlowa, PozycjaID: &pozycjaID, Akcja: "skierowanie_pozycji",
			NumerObiegu: obieg, Szczegoly: &szczegoly,
		})
	})
}
