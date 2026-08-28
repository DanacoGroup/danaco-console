// Odpowiedzialność pliku: repozytorium poczty prowadzi trwałość konfiguracji skrzynek pocztowych, dziennika listów odebranych i śladu listów wysłanych.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// SkrzynkaPocztowa to wiersz tabeli skrzynka_pocztowa, niosący konfigurację połączenia odbiorczego i wysyłkowego jednej skrzynki.
type SkrzynkaPocztowa struct {
	ID             int64
	Kod            string
	Adres          string
	HostOdbioru    string
	PortOdbioru    int
	HostWysylki    string
	PortWysylki    int
	Uzytkownik     string
	HasloOdwolanie *string
	TLSWeryfikacja bool
	TaktSekundy    int
	Aktywna        bool
}

// ListOdebrany to wiersz tabeli list_odebrany, niosący treść i stan przetworzenia jednego listu odebranego przez skrzynkę.
type ListOdebrany struct {
	ID                 int64
	SkrzynkaID         int64
	IdentyfikatorListu string
	Nadawca            string
	Temat              string
	Chwila             *string
	Tresc              string
	Odebrano           string
	Przetworzony       bool
}

// ListWyslany to wiersz tabeli `list_wyslany`. Wysyłka nieudana też jest
// wierszem — z `Powodzenie=false` i treścią błędu.
type ListWyslany struct {
	ID            int64
	SkrzynkaID    int64
	Adresat       string
	Temat         string
	Tresc         string
	WOdpowiedziNa *string
	Powodzenie    bool
	Blad          *string
}

// RepozytoriumPoczty jest kontraktem odczytu i zapisu skrzynek pocztowych, listów odebranych oraz śladu listów wysłanych.
type RepozytoriumPoczty interface {
	SkrzynkiCzynne(ctx context.Context) ([]SkrzynkaPocztowa, error)
	SkrzynkaPoKodzie(ctx context.Context, kod string) (SkrzynkaPocztowa, error)

	// DodajList zapisuje list odebrany; powtórny odbiór jest stanem taktowania, nie błędem.
	DodajList(ctx context.Context, list ListOdebrany) (bool, error)
	ListyNieprzetworzone(ctx context.Context, skrzynkaID int64) ([]ListOdebrany, error)
	OznaczPrzetworzony(ctx context.Context, id int64) error
	// OstatniList zwraca najświeższy list odebrany skrzynki, z ErrBrakWiersza, gdy dziennik jest pusty.
	OstatniList(ctx context.Context, skrzynkaID int64) (ListOdebrany, error)
	ListPoIdentyfikatorze(ctx context.Context, skrzynkaID int64, identyfikator string) (ListOdebrany, error)

	ZapiszWyslany(ctx context.Context, list ListWyslany) error
}

const (
	kolumnySkrzynki = `id, kod, adres, host_odbioru, port_odbioru, host_wysylki,
	                   port_wysylki, uzytkownik, haslo_odwolanie, tls_weryfikacja,
	                   takt_sekundy, aktywna`

	skrzynkiCzynneSQL = `SELECT ` + kolumnySkrzynki + ` FROM skrzynka_pocztowa
	                     WHERE aktywna = 1 ORDER BY kod`

	skrzynkaPoKodzieSQL = `SELECT ` + kolumnySkrzynki + ` FROM skrzynka_pocztowa
	                       WHERE kod = ?`

	kolumnyListu = `id, skrzynka_id, identyfikator_listu, nadawca, temat, chwila,
	                tresc, odebrano, przetworzony`

	// Zero zmienionych wierszy przy tym zapisie znaczy, że list o podanym identyfikatorze już jest zapisany w dzienniku listów.
	dodajListSQL = `INSERT OR IGNORE INTO list_odebrany
	                (skrzynka_id, identyfikator_listu, nadawca, temat, chwila, tresc)
	                VALUES (?, ?, ?, ?, ?, ?)`

	listyNieprzetworzoneSQL = `SELECT ` + kolumnyListu + ` FROM list_odebrany
	                           WHERE skrzynka_id = ? AND przetworzony = 0
	                           ORDER BY id`

	oznaczPrzetworzonySQL = `UPDATE list_odebrany SET przetworzony = 1 WHERE id = ?`

	ostatniListSQL = `SELECT ` + kolumnyListu + ` FROM list_odebrany
	                  WHERE skrzynka_id = ? ORDER BY id DESC LIMIT 1`

	listPoIdentyfikatorzeSQL = `SELECT ` + kolumnyListu + ` FROM list_odebrany
	                            WHERE skrzynka_id = ? AND identyfikator_listu = ?`

	zapiszWyslanySQL = `INSERT INTO list_wyslany
	                    (skrzynka_id, adresat, temat, tresc, w_odpowiedzi_na,
	                     powodzenie, blad)
	                    VALUES (?, ?, ?, ?, ?, ?, ?)`
)

type repozytoriumPoczty struct {
	zapytania *zapytania
}

func noweRepozytoriumPoczty(zapytania *zapytania) RepozytoriumPoczty {
	return &repozytoriumPoczty{zapytania: zapytania}
}

func (r *repozytoriumPoczty) SkrzynkiCzynne(ctx context.Context) ([]SkrzynkaPocztowa, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, skrzynkiCzynneSQL)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać skrzynek pocztowych: %w", err)
	}
	defer wiersze.Close()
	lista := []SkrzynkaPocztowa{}
	for wiersze.Next() {
		skrzynka, err := odczytajSkrzynke(wiersze)
		if err != nil {
			return nil, err
		}
		lista = append(lista, skrzynka)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt skrzynek pocztowych: %w", err)
	}
	return lista, nil
}

func (r *repozytoriumPoczty) SkrzynkaPoKodzie(ctx context.Context, kod string) (SkrzynkaPocztowa, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, skrzynkaPoKodzieSQL)
	if err != nil {
		return SkrzynkaPocztowa{}, err
	}
	skrzynka, err := odczytajSkrzynke(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return SkrzynkaPocztowa{}, ErrBrakWiersza
	}
	return skrzynka, err
}

func (r *repozytoriumPoczty) DodajList(ctx context.Context, list ListOdebrany) (bool, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, dodajListSQL)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, list.SkrzynkaID, list.IdentyfikatorListu,
		list.Nadawca, list.Temat, tekstDoKolumny(list.Chwila), list.Tresc)
	if err != nil {
		return false, fmt.Errorf("dane: nie można zapisać listu %q: %w", list.IdentyfikatorListu, err)
	}
	zmienione, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nieznany wynik zapisu listu %q: %w", list.IdentyfikatorListu, err)
	}
	return zmienione > 0, nil
}

func (r *repozytoriumPoczty) ListyNieprzetworzone(ctx context.Context, skrzynkaID int64) ([]ListOdebrany, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listyNieprzetworzoneSQL)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, skrzynkaID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać listów skrzynki %d: %w", skrzynkaID, err)
	}
	defer wiersze.Close()
	lista := []ListOdebrany{}
	for wiersze.Next() {
		list, err := odczytajList(wiersze)
		if err != nil {
			return nil, err
		}
		lista = append(lista, list)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt listów skrzynki %d: %w", skrzynkaID, err)
	}
	return lista, nil
}

func (r *repozytoriumPoczty) OznaczPrzetworzony(ctx context.Context, id int64) error {
	polecenie, err := r.zapytania.przygotuj(ctx, oznaczPrzetworzonySQL)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, id); err != nil {
		return fmt.Errorf("dane: nie można oznaczyć listu %d jako przetworzonego: %w", id, err)
	}
	return nil
}

func (r *repozytoriumPoczty) OstatniList(ctx context.Context, skrzynkaID int64) (ListOdebrany, error) {
	return r.jedenList(ctx, ostatniListSQL, skrzynkaID)
}

func (r *repozytoriumPoczty) ListPoIdentyfikatorze(ctx context.Context,
	skrzynkaID int64, identyfikator string) (ListOdebrany, error) {
	return r.jedenList(ctx, listPoIdentyfikatorzeSQL, skrzynkaID, identyfikator)
}

func (r *repozytoriumPoczty) ZapiszWyslany(ctx context.Context, list ListWyslany) error {
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszWyslanySQL)
	if err != nil {
		return err
	}
	_, err = polecenie.ExecContext(ctx, list.SkrzynkaID, list.Adresat, list.Temat,
		list.Tresc, tekstDoKolumny(list.WOdpowiedziNa),
		liczbaLogiczna(list.Powodzenie), tekstDoKolumny(list.Blad))
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać śladu wysyłki do %q: %w", list.Adresat, err)
	}
	return nil
}

// jedenList wykonuje zapytanie o pojedynczy wiersz listu i oddaje ErrBrakWiersza, gdy zapytanie nie trafia w żaden wiersz.
func (r *repozytoriumPoczty) jedenList(ctx context.Context, sqlText string, argumenty ...any) (ListOdebrany, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, sqlText)
	if err != nil {
		return ListOdebrany{}, err
	}
	list, err := odczytajList(polecenie.QueryRowContext(ctx, argumenty...))
	if errors.Is(err, sql.ErrNoRows) {
		return ListOdebrany{}, ErrBrakWiersza
	}
	return list, err
}

// odczytajSkrzynke składa strukturę skrzynki pocztowej z jednego wiersza wyniku zapytania, odczytując kolumny w ustalonej kolejności.
func odczytajSkrzynke(wiersz skaner) (SkrzynkaPocztowa, error) {
	var skrzynka SkrzynkaPocztowa
	var odwolanie sql.NullString
	var tls, aktywna int
	err := wiersz.Scan(&skrzynka.ID, &skrzynka.Kod, &skrzynka.Adres,
		&skrzynka.HostOdbioru, &skrzynka.PortOdbioru, &skrzynka.HostWysylki,
		&skrzynka.PortWysylki, &skrzynka.Uzytkownik, &odwolanie, &tls,
		&skrzynka.TaktSekundy, &aktywna)
	if errors.Is(err, sql.ErrNoRows) {
		return SkrzynkaPocztowa{}, err
	}
	if err != nil {
		return SkrzynkaPocztowa{}, fmt.Errorf("dane: nieczytelny wiersz skrzynki pocztowej: %w", err)
	}
	skrzynka.HasloOdwolanie = tekstZKolumny(odwolanie)
	skrzynka.TLSWeryfikacja, skrzynka.Aktywna = tls == 1, aktywna == 1
	return skrzynka, nil
}

// odczytajList składa strukturę listu odebranego z jednego wiersza wyniku zapytania, odczytując kolumny w ustalonej kolejności.
func odczytajList(wiersz skaner) (ListOdebrany, error) {
	var list ListOdebrany
	var chwila sql.NullString
	var przetworzony int
	err := wiersz.Scan(&list.ID, &list.SkrzynkaID, &list.IdentyfikatorListu,
		&list.Nadawca, &list.Temat, &chwila, &list.Tresc, &list.Odebrano, &przetworzony)
	if errors.Is(err, sql.ErrNoRows) {
		return ListOdebrany{}, err
	}
	if err != nil {
		return ListOdebrany{}, fmt.Errorf("dane: nieczytelny wiersz listu odebranego: %w", err)
	}
	list.Chwila = tekstZKolumny(chwila)
	list.Przetworzony = przetworzony == 1
	return list, nil
}
