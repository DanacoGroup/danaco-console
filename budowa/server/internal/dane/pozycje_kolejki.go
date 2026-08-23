// Odpowiedzialność pliku: zlecenia kolejki (tabela `pozycja_kolejki`).
// Licznik obiegów naprawczych rośnie bez limitu — przerwanie należy do
// użytkownika, nie do warstwy trwałości.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"danacoconsole/shared"
)

// Pozycja to wiersz tabeli `pozycja_kolejki`. Stan i werdykt są słownikami
// schematu; kontrakt nie ma dla nich wyliczenia.
type Pozycja struct {
	ID                 int64
	KolejkaID          int64
	OknoWykonawcyID    *int64
	Kolejnosc          int
	Tytul              string
	TrescZlecenia      *string
	TrescOdwolanie     *string
	Stan               string
	WerdyktWeryfikacji *string
	LicznikObiegow     int
	Utworzono          string
	Zaktualizowano     string
}

// domyslnyStanPozycji odpowiada wartości domyślnej kolumny w schemacie.
const domyslnyStanPozycji = "oczekuje"

const (
	kolumnyPozycji = `id, kolejka_id, okno_wykonawcy_id, kolejnosc, tytul, tresc_zlecenia,
	                  tresc_odwolanie, stan, werdykt_weryfikacji, licznik_obiegow,
	                  utworzono, zaktualizowano`

	wstawPozycje = `INSERT INTO pozycja_kolejki
	                (kolejka_id, okno_wykonawcy_id, kolejnosc, tytul, tresc_zlecenia,
	                 tresc_odwolanie, stan, werdykt_weryfikacji, licznik_obiegow)
	                VALUES (?, ?, (SELECT COALESCE(MAX(kolejnosc) + 1, 0) FROM pozycja_kolejki
	                               WHERE kolejka_id = ?), ?, ?, ?, ?, ?, ?)`

	zmienStanPozycji = `UPDATE pozycja_kolejki
	                    SET stan = ?, werdykt_weryfikacji = ?,
	                        zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                    WHERE id = ?`

	zwiekszObiegPozycji = `UPDATE pozycja_kolejki
	                       SET licznik_obiegow = licznik_obiegow + 1,
	                           zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                       WHERE id = ?
	                       RETURNING kolejka_id, licznik_obiegow`

	listaPozycjiKolejki = `SELECT ` + kolumnyPozycji + ` FROM pozycja_kolejki
	                       WHERE kolejka_id = ? ORDER BY kolejnosc, id`
)

// DodajPozycje dokłada zlecenie na koniec kolejki i odnotowuje je w dzienniku.
func (r *repozytoriumKolejek) DodajPozycje(ctx context.Context, pozycja Pozycja) (int64, error) {
	stan := pozycja.Stan
	if stan == "" {
		stan = domyslnyStanPozycji
	}
	var id int64
	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, wstawPozycje)
		if err != nil {
			return err
		}
		wynik, err := polecenie.ExecContext(ctx, pozycja.KolejkaID,
			liczbaDoKolumny(pozycja.OknoWykonawcyID), pozycja.KolejkaID, pozycja.Tytul,
			tekstDoKolumny(pozycja.TrescZlecenia), tekstDoKolumny(pozycja.TrescOdwolanie),
			stan, tekstDoKolumny(pozycja.WerdyktWeryfikacji), pozycja.LicznikObiegow)
		if err != nil {
			return fmt.Errorf("dane: nie można dodać pozycji do kolejki %d: %w", pozycja.KolejkaID, err)
		}
		if id, err = wynik.LastInsertId(); err != nil {
			return fmt.Errorf("dane: nieznany identyfikator dodanej pozycji: %w", err)
		}
		return dopiszWpisDziennika(ctx, r.zapytania, transakcja, WpisDziennika{
			KolejkaID: pozycja.KolejkaID, PozycjaID: &id, Akcja: "dodanie_pozycji", StanPo: &stan,
		})
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}

// ZmienStanPozycji zapisuje stan zlecenia i werdykt weryfikacji.
func (r *repozytoriumKolejek) ZmienStanPozycji(ctx context.Context, id int64,
	stan string, werdykt *string) error {

	return wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		kolejkaID, obieg, err := polozeniePozycji(ctx, r.zapytania, transakcja, id)
		if err != nil {
			return err
		}
		polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, zmienStanPozycji)
		if err != nil {
			return err
		}
		wynik, err := polecenie.ExecContext(ctx, stan, tekstDoKolumny(werdykt), id)
		if err != nil {
			return fmt.Errorf("dane: nie można zmienić stanu pozycji %d: %w", id, err)
		}
		if err := sprawdzTrafienie(wynik, "pozycja_kolejki", id); err != nil {
			return err
		}
		return dopiszWpisDziennika(ctx, r.zapytania, transakcja, WpisDziennika{
			KolejkaID: kolejkaID, PozycjaID: &id, Akcja: "zmiana_stanu_pozycji",
			StanPo: &stan, NumerObiegu: obieg, Szczegoly: werdykt,
		})
	})
}

// ZwiekszObieg podnosi licznik biegu naprawczego. Bez limitu i bez warunku —
// krok z błędami powtarza się dowolną liczbę razy.
func (r *repozytoriumKolejek) ZwiekszObieg(ctx context.Context, pozycjaID int64) (int, error) {
	var obieg int
	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, zwiekszObiegPozycji)
		if err != nil {
			return err
		}
		var kolejkaID int64
		err = polecenie.QueryRowContext(ctx, pozycjaID).Scan(&kolejkaID, &obieg)
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("dane: pozycja kolejki %d nie istnieje", pozycjaID)
		}
		if err != nil {
			return fmt.Errorf("dane: nie można podnieść licznika obiegów pozycji %d: %w", pozycjaID, err)
		}
		return dopiszWpisDziennika(ctx, r.zapytania, transakcja, WpisDziennika{
			KolejkaID: kolejkaID, PozycjaID: &pozycjaID,
			Akcja: string(shared.QueueActionRetry), NumerObiegu: obieg,
		})
	})
	if err != nil {
		return 0, err
	}
	return obieg, nil
}

// ListaPozycji zwraca zlecenia kolejki w zapisanej kolejności.
func (r *repozytoriumKolejek) ListaPozycji(ctx context.Context, kolejkaID int64) ([]Pozycja, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaPozycjiKolejki)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, kolejkaID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać pozycji kolejki %d: %w", kolejkaID, err)
	}
	defer wiersze.Close()

	lista := []Pozycja{}
	for wiersze.Next() {
		pozycja, err := odczytajPozycje(wiersze)
		if err != nil {
			return nil, err
		}
		lista = append(lista, pozycja)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt pozycji kolejki %d: %w", kolejkaID, err)
	}
	return lista, nil
}

// odczytajPozycje składa strukturę z jednego wiersza wyniku.
func odczytajPozycje(wiersz skaner) (Pozycja, error) {
	var pozycja Pozycja
	var okno sql.NullInt64
	var zlecenie, odwolanie, werdykt sql.NullString
	err := wiersz.Scan(&pozycja.ID, &pozycja.KolejkaID, &okno, &pozycja.Kolejnosc, &pozycja.Tytul,
		&zlecenie, &odwolanie, &pozycja.Stan, &werdykt, &pozycja.LicznikObiegow,
		&pozycja.Utworzono, &pozycja.Zaktualizowano)
	if err != nil {
		return Pozycja{}, fmt.Errorf("dane: nieczytelny wiersz pozycji kolejki: %w", err)
	}
	pozycja.OknoWykonawcyID = liczbaZKolumny(okno)
	pozycja.TrescZlecenia = tekstZKolumny(zlecenie)
	pozycja.TrescOdwolanie = tekstZKolumny(odwolanie)
	pozycja.WerdyktWeryfikacji = tekstZKolumny(werdykt)
	return pozycja, nil
}
