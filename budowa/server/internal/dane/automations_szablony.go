// Plik prowadzi bibliotekę szablonów przepływów: szablon jest bytem odrębnym od automatyki, z której powstał,
// i ma dalej zakładać automatyki o kształcie z chwili zapisu, dlatego kroki leżą tu jako migawka, bez klucza obcego do automatyki.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// SzablonAutomatyki to wiersz tabeli `szablon_automatyki` niosący migawkę kroków wraz z parametrami szablonu.
type SzablonAutomatyki struct {
	ID             int64
	Kod            string
	Nazwa          string
	Opis           *string
	Kroki          string
	Zaktualizowano string
}

// ParametrSzablonu to wiersz tabeli `parametr_szablonu_automatyki` — pole
// formularza uzupełnianego przy zastosowaniu szablonu.
type ParametrSzablonu struct {
	Nazwa           string
	Etykieta        *string
	Wymagany        bool
	WartoscDomyslna *string
	Kolejnosc       int
}

const (
	kolumnySzablonuAutomatyki = `id, identyfikator_zewnetrzny, nazwa, opis, kroki, zaktualizowano`

	zapiszSzablonAutomatyki = `INSERT INTO szablon_automatyki
	                           (identyfikator_zewnetrzny, nazwa, opis, kroki)
	                           VALUES (?, ?, ?, ?)
	                           ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                               nazwa = excluded.nazwa,
	                               opis = excluded.opis,
	                               kroki = excluded.kroki,
	                               zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')`

	pobierzSzablonAutomatyki = `SELECT ` + kolumnySzablonuAutomatyki + ` FROM szablon_automatyki
	                            WHERE identyfikator_zewnetrzny = ?`

	listaSzablonowAutomatyki = `SELECT ` + kolumnySzablonuAutomatyki + ` FROM szablon_automatyki
	                            ORDER BY nazwa, id LIMIT ?`

	usunParametrySzablonu = `DELETE FROM parametr_szablonu_automatyki WHERE szablon_id = ?`

	wstawParametrSzablonu = `INSERT INTO parametr_szablonu_automatyki
	                         (szablon_id, nazwa, etykieta, wymagany, wartosc_domyslna, kolejnosc)
	                         VALUES (?, ?, ?, ?, ?, ?)`

	listaParametrowSzablonu = `SELECT nazwa, etykieta, wymagany, wartosc_domyslna, kolejnosc
	                           FROM parametr_szablonu_automatyki WHERE szablon_id = ?
	                           ORDER BY kolejnosc, nazwa`
)

// ZapiszSzablonAutomatyki zakłada szablon albo nadpisuje zastany wraz z jego
// parametrami. Jedna transakcja, bo szablon zapisany z parametrami poprzedniej
// wersji byłby formularzem pytającym o pola, których szablon już nie zna.
func (r *repozytoriumAutomatyk) ZapiszSzablonAutomatyki(ctx context.Context,
	szablon SzablonAutomatyki, parametry []ParametrSzablonu) (SzablonAutomatyki, error) {

	if szablon.Kod == "" {
		return SzablonAutomatyki{}, fmt.Errorf("dane: szablon automatyki bez identyfikatora")
	}
	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		if err := wykonajWTransakcjiAutomatyzacji(ctx, r.zapytania, transakcja, zapiszSzablonAutomatyki,
			szablon.Kod, szablon.Nazwa, tekstDoKolumny(szablon.Opis), szablon.Kroki); err != nil {
			return err
		}
		zapisany, err := szablonWTransakcji(ctx, r.zapytania, transakcja, szablon.Kod)
		if err != nil {
			return err
		}
		if err := wykonajWTransakcjiAutomatyzacji(ctx, r.zapytania, transakcja,
			usunParametrySzablonu, zapisany.ID); err != nil {
			return err
		}
		for numer, parametr := range parametry {
			if parametr.Nazwa == "" {
				continue
			}
			err := wykonajWTransakcjiAutomatyzacji(ctx, r.zapytania, transakcja, wstawParametrSzablonu,
				zapisany.ID, parametr.Nazwa, tekstDoKolumny(parametr.Etykieta),
				liczbaLogiczna(parametr.Wymagany), tekstDoKolumny(parametr.WartoscDomyslna),
				numer+1)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return SzablonAutomatyki{}, err
	}
	return r.SzablonAutomatyki(ctx, szablon.Kod)
}

// SzablonAutomatyki zwraca szablon po jego kodzie zewnętrznym; brak wiersza wraca jako ErrBrakWiersza.
func (r *repozytoriumAutomatyk) SzablonAutomatyki(ctx context.Context,
	kod string) (SzablonAutomatyki, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzSzablonAutomatyki)
	if err != nil {
		return SzablonAutomatyki{}, err
	}
	szablon, err := odczytajSzablonAutomatyki(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return SzablonAutomatyki{}, ErrBrakWiersza
	}
	if err != nil {
		return SzablonAutomatyki{}, fmt.Errorf("dane: nieczytelny szablon automatyki %q: %w", kod, err)
	}
	return szablon, nil
}

// SzablonyAutomatyki zwraca całą bibliotekę szablonów w kolejności nazw wprost z bazy danych repozytorium.
func (r *repozytoriumAutomatyk) SzablonyAutomatyki(ctx context.Context,
	limit int) ([]SzablonAutomatyki, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaSzablonowAutomatyki)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, granicaWykazu(limit))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać szablonów automatyk: %w", err)
	}
	defer wiersze.Close()

	lista := []SzablonAutomatyki{}
	for wiersze.Next() {
		szablon, err := odczytajSzablonAutomatyki(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz szablonu automatyki: %w", err)
		}
		lista = append(lista, szablon)
	}
	return lista, wiersze.Err()
}

// ParametrySzablonuAutomatyki zwraca parametry szablonu w kolejności ich zapisu wprost z bazy danych repozytorium.
func (r *repozytoriumAutomatyk) ParametrySzablonuAutomatyki(ctx context.Context,
	szablonID int64) ([]ParametrSzablonu, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaParametrowSzablonu)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, szablonID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać parametrów szablonu %d: %w", szablonID, err)
	}
	defer wiersze.Close()

	parametry := []ParametrSzablonu{}
	for wiersze.Next() {
		var parametr ParametrSzablonu
		var etykieta, wartosc sql.NullString
		var wymagany int
		if err := wiersze.Scan(&parametr.Nazwa, &etykieta, &wymagany, &wartosc,
			&parametr.Kolejnosc); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz parametru szablonu: %w", err)
		}
		parametr.Etykieta = tekstZKolumny(etykieta)
		parametr.Wymagany = wymagany == 1
		parametr.WartoscDomyslna = tekstZKolumny(wartosc)
		parametry = append(parametry, parametr)
	}
	return parametry, wiersze.Err()
}

// szablonWTransakcji dobiera świeżo zapisany szablon wewnątrz tej samej
// transakcji — potrzebny jest jego klucz wiersza pod parametry.
func szablonWTransakcji(ctx context.Context, z *zapytania, transakcja *sql.Tx,
	kod string) (SzablonAutomatyki, error) {

	polecenie, err := z.wTransakcji(ctx, transakcja, pobierzSzablonAutomatyki)
	if err != nil {
		return SzablonAutomatyki{}, err
	}
	szablon, err := odczytajSzablonAutomatyki(polecenie.QueryRowContext(ctx, kod))
	if err != nil {
		return SzablonAutomatyki{}, fmt.Errorf("dane: nieczytelny szablon automatyki %q: %w", kod, err)
	}
	return szablon, nil
}

// odczytajSzablonAutomatyki składa szablon automatyki wprost z jednego wiersza wyniku zapytania do bazy danych.
func odczytajSzablonAutomatyki(wiersz skaner) (SzablonAutomatyki, error) {
	var szablon SzablonAutomatyki
	var opis sql.NullString
	err := wiersz.Scan(&szablon.ID, &szablon.Kod, &szablon.Nazwa, &opis,
		&szablon.Kroki, &szablon.Zaktualizowano)
	if err != nil {
		return SzablonAutomatyki{}, err
	}
	szablon.Opis = tekstZKolumny(opis)
	return szablon, nil
}
