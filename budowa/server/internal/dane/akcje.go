// Plik odczytuje katalog akcji sterowany danymi: nowa akcja to nowy wiersz, nie nowa gałąź w kodzie, wzorem rejestru
// kanałów modelu; akcja należy do jednego z ośmiu poziomów zasięgu, a pusty klucz zasięgu znaczy każdy byt tego poziomu.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"danacoconsole/shared"
)

// Akcja to wiersz katalogu akcji wraz z poziomem zasięgu i kluczem wskazującym konkretny byt tego poziomu.
type Akcja struct {
	ID                 int64
	Kod                string
	Nazwa              string
	Opis               string
	Ikona              string
	PoziomZasiegu      shared.ConfigScope
	KluczZasiegu       string
	Komenda            shared.MessageType
	WarunekDostepnosci string
	Kolejnosc          int
	Aktywna            bool
}

// RepozytoriumAkcji jest kontraktem katalogu akcji: zapisu tu nie ma, wiersze wnosi zaczyn migracji, a wybór akcji
// jednego bytu poziomu rozstrzyga rejestr akcji w rdzeniu na wierszach już odczytanych.
type RepozytoriumAkcji interface {
	Lista(ctx context.Context, tylkoAktywne bool) ([]Akcja, error)
	PoKodzie(ctx context.Context, kod string) (Akcja, error)
}

const (
	kolumnyAkcji = `a.id, a.kod, a.nazwa, a.opis, a.ikona, p.kod, a.klucz_zasiegu,
	                a.komenda, a.warunek_dostepnosci, a.kolejnosc, a.aktywna`

	zrodloAkcji = ` FROM akcja a JOIN poziom_zasiegu p ON p.id = a.poziom_zasiegu_id`

	listaAkcji = `SELECT ` + kolumnyAkcji + zrodloAkcji + `
	              WHERE (? = 0 OR a.aktywna = 1)
	              ORDER BY p.pierwszenstwo, a.klucz_zasiegu, a.kolejnosc, a.kod`

	pobierzAkcjePoKodzie = `SELECT ` + kolumnyAkcji + zrodloAkcji + ` WHERE a.kod = ?`
)

type repozytoriumAkcji struct {
	zapytania *zapytania
}

// noweRepozytoriumAkcji zakłada repozytorium katalogu akcji nad pamięcią
// przygotowanych zapytań zestawu.
func noweRepozytoriumAkcji(z *zapytania) *repozytoriumAkcji {
	return &repozytoriumAkcji{zapytania: z}
}

// Lista zwraca cały katalog akcji uporządkowany od poziomu najszerszego.
// Katalog pusty nie jest błędem — rejestr zbuduje się bez akcji, a platforma
// pracuje dalej.
func (r *repozytoriumAkcji) Lista(ctx context.Context, tylkoAktywne bool) ([]Akcja, error) {
	return r.wykaz(ctx, listaAkcji, "katalogu akcji", liczbaLogiczna(tylkoAktywne))
}

// PoKodzie zwraca jedną pozycję katalogu. Brak wiersza jest sygnałem
// ErrBrakWiersza, nie awarią odczytu.
func (r *repozytoriumAkcji) PoKodzie(ctx context.Context, kod string) (Akcja, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzAkcjePoKodzie)
	if err != nil {
		return Akcja{}, err
	}
	akcja, err := odczytajAkcje(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return Akcja{}, fmt.Errorf("dane: akcja %q nie istnieje w katalogu: %w", kod, ErrBrakWiersza)
	}
	return akcja, err
}

// wykaz wykonuje zapytanie zwracające wiele wierszy katalogu akcji dla wskazanego poziomu zasięgu i klucza.
func (r *repozytoriumAkcji) wykaz(ctx context.Context, zapytanie, opis string,
	argumenty ...any) ([]Akcja, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, argumenty...)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać %s: %w", opis, err)
	}
	defer wiersze.Close()

	lista := []Akcja{}
	for wiersze.Next() {
		akcja, err := odczytajAkcje(wiersze)
		if err != nil {
			return nil, err
		}
		lista = append(lista, akcja)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt %s: %w", opis, err)
	}
	return lista, nil
}

// odczytajAkcje składa strukturę pojedynczej akcji katalogu wprost z jednego wiersza wyniku zapytania.
func odczytajAkcje(wiersz skaner) (Akcja, error) {
	var akcja Akcja
	var poziom string
	var aktywna int
	err := wiersz.Scan(&akcja.ID, &akcja.Kod, &akcja.Nazwa, &akcja.Opis, &akcja.Ikona,
		&poziom, &akcja.KluczZasiegu, &akcja.Komenda, &akcja.WarunekDostepnosci,
		&akcja.Kolejnosc, &aktywna)
	if err != nil {
		return Akcja{}, err
	}
	if akcja.PoziomZasiegu, err = poziomZasieguZBazy(poziom); err != nil {
		return Akcja{}, err
	}
	akcja.Aktywna = aktywna != 0
	return akcja, nil
}
