// Odpowiedzialność pliku: dostęp do słownika środowisk (tabela `srodowisko`).
// Środowisko jest profilem widoczności modułów, nie pojemnikiem, do którego
// moduł należy na wyłączność. Wiersze wnosi zaczyn schematu — repozytorium ich
// nie zakłada, wyłącznie czyta.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// Srodowisko to wiersz tabeli `srodowisko`.
//
// Opis i Motto to dwa różne zdania o środowisku, drukowane w dwóch różnych
// miejscach: opis objaśnia tryb pracy na karcie wejścia, motto jest podtytułem
// w nagłówku kolumny matrycy. Oba bywają puste — brak motta jest stanem
// normalnym, nie usterką wiersza.
type Srodowisko struct {
	ID        int64
	Kod       string
	Nazwa     string
	Opis      *string
	Motto     *string
	Kolejnosc int
	Aktywne   bool
}

// RepozytoriumSrodowisk jest kontraktem słownika środowisk.
type RepozytoriumSrodowisk interface {
	Lista(ctx context.Context) ([]Srodowisko, error)
	PoKodzie(ctx context.Context, kod string) (Srodowisko, error)
	Pierwsze(ctx context.Context) (Srodowisko, error)
}

const (
	kolumnySrodowiska = `id, kod, nazwa, opis, motto, kolejnosc, aktywne`

	listaSrodowisk = `SELECT ` + kolumnySrodowiska + ` FROM srodowisko ORDER BY kolejnosc, kod`

	srodowiskoPoKodzie = `SELECT ` + kolumnySrodowiska + ` FROM srodowisko WHERE kod = ?`

	pierwszeSrodowisko = `SELECT ` + kolumnySrodowiska + ` FROM srodowisko
	                      WHERE aktywne = 1 ORDER BY kolejnosc, kod LIMIT 1`
)

type repozytoriumSrodowisk struct {
	zapytania *zapytania
}

func noweRepozytoriumSrodowisk(z *zapytania) *repozytoriumSrodowisk {
	return &repozytoriumSrodowisk{zapytania: z}
}

// Lista zwraca wszystkie środowiska w kolejności nawigacji.
func (r *repozytoriumSrodowisk) Lista(ctx context.Context) ([]Srodowisko, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaSrodowisk)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać środowisk: %w", err)
	}
	defer wiersze.Close()

	lista := []Srodowisko{}
	for wiersze.Next() {
		srodowisko, err := odczytajSrodowisko(wiersze)
		if err != nil {
			return nil, err
		}
		lista = append(lista, srodowisko)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt środowisk: %w", err)
	}
	return lista, nil
}

// PoKodzie zwraca środowisko o wskazanym kodzie.
func (r *repozytoriumSrodowisk) PoKodzie(ctx context.Context, kod string) (Srodowisko, error) {
	return r.jedno(ctx, srodowiskoPoKodzie, fmt.Sprintf("o kodzie %q", kod), kod)
}

// Pierwsze zwraca środowisko czynne o najniższej kolejności. Służy tam, gdzie
// środowiska nie wskazano wprost, a wiersz podrzędny musi je mieć.
func (r *repozytoriumSrodowisk) Pierwsze(ctx context.Context) (Srodowisko, error) {
	return r.jedno(ctx, pierwszeSrodowisko, "czynne")
}

// jedno wykonuje zapytanie zwracające najwyżej jeden wiersz środowiska.
func (r *repozytoriumSrodowisk) jedno(ctx context.Context, zapytanie, opis string, argumenty ...any) (Srodowisko, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return Srodowisko{}, err
	}
	srodowisko, err := odczytajSrodowisko(polecenie.QueryRowContext(ctx, argumenty...))
	if errors.Is(err, sql.ErrNoRows) {
		return Srodowisko{}, fmt.Errorf("%w: środowisko %s", ErrBrakWiersza, opis)
	}
	return srodowisko, err
}

// odczytajSrodowisko składa strukturę z jednego wiersza wyniku.
func odczytajSrodowisko(wiersz skaner) (Srodowisko, error) {
	var srodowisko Srodowisko
	var opis, motto sql.NullString
	var aktywne int
	err := wiersz.Scan(&srodowisko.ID, &srodowisko.Kod, &srodowisko.Nazwa, &opis, &motto,
		&srodowisko.Kolejnosc, &aktywne)
	if err != nil {
		return Srodowisko{}, err
	}
	srodowisko.Opis = tekstZKolumny(opis)
	srodowisko.Motto = tekstZKolumny(motto)
	srodowisko.Aktywne = aktywne != 0
	return srodowisko, nil
}
