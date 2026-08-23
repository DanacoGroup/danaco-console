// Odpowiedzialność pliku: dostęp do kart sesji (tabela `karta_sesji`). Karta jest
// kontenerem sesji w obrębie środowiska — sesja bez karty nie istnieje, bo więz
// klucza obcego jest obowiązkowy. Dlatego repozytorium ma metodę `Zapewnij`:
// warstwa wyższa zapisuje sesję, a karta ma powstać po drodze, nie zablokować
// zapisu.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// KartaSesji to wiersz tabeli `karta_sesji`.
type KartaSesji struct {
	ID             int64
	SrodowiskoID   int64
	Nazwa          string
	Opis           *string
	Przypieta      bool
	Kolejnosc      int
	Utworzono      string
	Zaktualizowano string
}

// RepozytoriumKartSesji jest kontraktem obszaru kart sesji.
type RepozytoriumKartSesji interface {
	Zapewnij(ctx context.Context, srodowiskoID int64, nazwa string) (int64, error)
	Lista(ctx context.Context, srodowiskoID int64) ([]KartaSesji, error)
}

const (
	kolumnyKartySesji = `id, srodowisko_id, nazwa, opis, przypieta, kolejnosc, utworzono, zaktualizowano`

	kartaSesjiPoNazwie = `SELECT ` + kolumnyKartySesji + ` FROM karta_sesji
	                      WHERE srodowisko_id = ? AND nazwa = ? ORDER BY id LIMIT 1`

	wstawKarteSesji = `INSERT INTO karta_sesji (srodowisko_id, nazwa, kolejnosc)
	                   VALUES (?, ?, (SELECT COALESCE(MAX(kolejnosc) + 1, 0) FROM karta_sesji
	                                  WHERE srodowisko_id = ?))`

	listaKartSesji = `SELECT ` + kolumnyKartySesji + ` FROM karta_sesji
	                  WHERE srodowisko_id = ? ORDER BY kolejnosc, id`
)

type repozytoriumKartSesji struct {
	zapytania *zapytania
}

func noweRepozytoriumKartSesji(z *zapytania) *repozytoriumKartSesji {
	return &repozytoriumKartSesji{zapytania: z}
}

// Zapewnij zwraca identyfikator karty o wskazanej nazwie, zakładając ją, gdy
// jeszcze nie istnieje. Wywołanie powtórne nie mnoży wierszy.
func (r *repozytoriumKartSesji) Zapewnij(ctx context.Context, srodowiskoID int64, nazwa string) (int64, error) {
	odczyt, err := r.zapytania.przygotuj(ctx, kartaSesjiPoNazwie)
	if err != nil {
		return 0, err
	}
	karta, err := odczytajKarteSesji(odczyt.QueryRowContext(ctx, srodowiskoID, nazwa))
	if err == nil {
		return karta.ID, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}
	zapis, err := r.zapytania.przygotuj(ctx, wstawKarteSesji)
	if err != nil {
		return 0, err
	}
	wynik, err := zapis.ExecContext(ctx, srodowiskoID, nazwa, srodowiskoID)
	if err != nil {
		return 0, fmt.Errorf("dane: nie można zapisać karty sesji %q: %w", nazwa, err)
	}
	return wynik.LastInsertId()
}

// Lista zwraca karty środowiska w kolejności wyświetlania.
func (r *repozytoriumKartSesji) Lista(ctx context.Context, srodowiskoID int64) ([]KartaSesji, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaKartSesji)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, srodowiskoID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać kart środowiska %d: %w", srodowiskoID, err)
	}
	defer wiersze.Close()

	lista := []KartaSesji{}
	for wiersze.Next() {
		karta, err := odczytajKarteSesji(wiersze)
		if err != nil {
			return nil, err
		}
		lista = append(lista, karta)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt kart środowiska %d: %w", srodowiskoID, err)
	}
	return lista, nil
}

// odczytajKarteSesji składa strukturę z jednego wiersza wyniku.
func odczytajKarteSesji(wiersz skaner) (KartaSesji, error) {
	var karta KartaSesji
	var opis sql.NullString
	var przypieta int
	err := wiersz.Scan(&karta.ID, &karta.SrodowiskoID, &karta.Nazwa, &opis, &przypieta,
		&karta.Kolejnosc, &karta.Utworzono, &karta.Zaktualizowano)
	if err != nil {
		return KartaSesji{}, err
	}
	karta.Opis = tekstZKolumny(opis)
	karta.Przypieta = przypieta != 0
	return karta, nil
}
