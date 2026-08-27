// Odpowiedzialność pliku: dostęp do rejestru kanałów modelu (tabela
// `kanal_modelu`). Rejestr jest sterowany danymi — nowy kanał to nowy wiersz, nie nowy typ w kodzie.
package dane

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// Kanal to wiersz rejestru kanałów modelu, niosący jego konfigurację połączenia i parametry działania.
type Kanal struct {
	ID                     int64
	Kod                    string
	Nazwa                  string
	Dostawca               string
	IdentyfikatorModelu    string
	RodzajKanalu           string
	KontoID                *int64
	PoswiadczenieOdwolanie *string
	ParametryJSON          string
	Multimodalny           bool
	Aktywny                bool
	Kolejnosc              int
	Utworzono              string
}

// RepozytoriumKanalow jest kontraktem rejestru kanałów, określającym operacje dostępne na wykazie kanałów.
type RepozytoriumKanalow interface {
	Dodaj(ctx context.Context, kanal Kanal) (int64, error)
	Aktualizuj(ctx context.Context, kanal Kanal) error
	Usun(ctx context.Context, id int64) error
	Lista(ctx context.Context, tylkoAktywne bool) ([]Kanal, error)
	PobierzPoKodzie(ctx context.Context, kod string) (Kanal, error)
}

const (
	kolumnyKanalu = `id, kod, nazwa, dostawca, identyfikator_modelu, rodzaj_kanalu, konto_id,
	                 poswiadczenie_odwolanie, parametry_json, multimodalny, aktywny, kolejnosc, utworzono`

	wstawKanal = `INSERT INTO kanal_modelu
	              (kod, nazwa, dostawca, identyfikator_modelu, rodzaj_kanalu, konto_id,
	               poswiadczenie_odwolanie, parametry_json, multimodalny, aktywny, kolejnosc)
	              VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	aktualizujKanal = `UPDATE kanal_modelu
	                   SET nazwa = ?, dostawca = ?, identyfikator_modelu = ?, rodzaj_kanalu = ?,
	                       konto_id = ?, poswiadczenie_odwolanie = ?, parametry_json = ?,
	                       multimodalny = ?, aktywny = ?, kolejnosc = ?
	                   WHERE id = ?`

	usunKanal = `DELETE FROM kanal_modelu WHERE id = ?`

	listaKanalow = `SELECT ` + kolumnyKanalu + ` FROM kanal_modelu
	                WHERE (? = 0 OR aktywny = 1) ORDER BY kolejnosc, id`

	pobierzKanalPoKodzie = `SELECT ` + kolumnyKanalu + ` FROM kanal_modelu WHERE kod = ?`
)

type repozytoriumKanalow struct {
	zapytania *zapytania
}

func noweRepozytoriumKanalow(z *zapytania) *repozytoriumKanalow {
	return &repozytoriumKanalow{zapytania: z}
}

// Dodaj wpisuje nowy kanał do rejestru wraz z jego parametrami połączenia zapisanymi jako dokument JSON.
func (r *repozytoriumKanalow) Dodaj(ctx context.Context, kanal Kanal) (int64, error) {
	parametry, err := parametryKanalu(kanal)
	if err != nil {
		return 0, err
	}
	polecenie, err := r.zapytania.przygotuj(ctx, wstawKanal)
	if err != nil {
		return 0, err
	}
	wynik, err := polecenie.ExecContext(ctx, kanal.Kod, kanal.Nazwa, kanal.Dostawca,
		kanal.IdentyfikatorModelu, kanal.RodzajKanalu, liczbaDoKolumny(kanal.KontoID),
		tekstDoKolumny(kanal.PoswiadczenieOdwolanie), parametry,
		liczbaLogiczna(kanal.Multimodalny), liczbaLogiczna(kanal.Aktywny), kanal.Kolejnosc)
	if err != nil {
		return 0, fmt.Errorf("dane: nie można dodać kanału %q: %w", kanal.Kod, err)
	}
	return wynik.LastInsertId()
}

// Aktualizuj zapisuje zmieniony wiersz rejestru. Kod kanału pozostaje stały —
// jest jego identyfikatorem w konfiguracji okien.
func (r *repozytoriumKanalow) Aktualizuj(ctx context.Context, kanal Kanal) error {
	parametry, err := parametryKanalu(kanal)
	if err != nil {
		return err
	}
	polecenie, err := r.zapytania.przygotuj(ctx, aktualizujKanal)
	if err != nil {
		return err
	}
	wynik, err := polecenie.ExecContext(ctx, kanal.Nazwa, kanal.Dostawca, kanal.IdentyfikatorModelu,
		kanal.RodzajKanalu, liczbaDoKolumny(kanal.KontoID),
		tekstDoKolumny(kanal.PoswiadczenieOdwolanie), parametry,
		liczbaLogiczna(kanal.Multimodalny), liczbaLogiczna(kanal.Aktywny), kanal.Kolejnosc, kanal.ID)
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać kanału %d: %w", kanal.ID, err)
	}
	return sprawdzTrafienie(wynik, "kanal_modelu", kanal.ID)
}

// Usun kasuje wiersz rejestru kanałów wskazany identyfikatorem, trwale usuwając kanał z konfiguracji modelu.
func (r *repozytoriumKanalow) Usun(ctx context.Context, id int64) error {
	polecenie, err := r.zapytania.przygotuj(ctx, usunKanal)
	if err != nil {
		return err
	}
	wynik, err := polecenie.ExecContext(ctx, id)
	if err != nil {
		return fmt.Errorf("dane: nie można usunąć kanału %d: %w", id, err)
	}
	return sprawdzTrafienie(wynik, "kanal_modelu", id)
}

// Lista zwraca rejestr kanałów modelu — komplet wpisów albo wyłącznie kanały czynne, zależnie od parametru.
func (r *repozytoriumKanalow) Lista(ctx context.Context, tylkoAktywne bool) ([]Kanal, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaKanalow)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, liczbaLogiczna(tylkoAktywne))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać rejestru kanałów: %w", err)
	}
	defer wiersze.Close()

	lista := []Kanal{}
	for wiersze.Next() {
		kanal, err := odczytajKanal(wiersze)
		if err != nil {
			return nil, err
		}
		lista = append(lista, kanal)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt rejestru kanałów: %w", err)
	}
	return lista, nil
}

// PobierzPoKodzie zwraca kanał modelu wskazany kodem, jaki jest użyty w konfiguracji danego okna komunikacji.
func (r *repozytoriumKanalow) PobierzPoKodzie(ctx context.Context, kod string) (Kanal, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzKanalPoKodzie)
	if err != nil {
		return Kanal{}, err
	}
	kanal, err := odczytajKanal(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return Kanal{}, fmt.Errorf("dane: kanał %q nie istnieje w rejestrze", kod)
	}
	return kanal, err
}

// parametryKanalu pilnuje, żeby kolumna parametrów była poprawnym obiektem JSON.
// Pusta wartość oznacza brak parametrów, nie błąd.
func parametryKanalu(kanal Kanal) (string, error) {
	parametry := strings.TrimSpace(kanal.ParametryJSON)
	if parametry == "" {
		return "{}", nil
	}
	if !json.Valid([]byte(parametry)) {
		return "", fmt.Errorf("dane: parametry kanału %q nie są poprawnym JSON", kanal.Kod)
	}
	return parametry, nil
}

// odczytajKanal składa pełną strukturę kanału modelu z jednego wiersza wyniku zapytania do bazy danych.
func odczytajKanal(wiersz skaner) (Kanal, error) {
	var kanal Kanal
	var kontoID sql.NullInt64
	var odwolanie sql.NullString
	var multimodalny, aktywny int
	err := wiersz.Scan(&kanal.ID, &kanal.Kod, &kanal.Nazwa, &kanal.Dostawca,
		&kanal.IdentyfikatorModelu, &kanal.RodzajKanalu, &kontoID, &odwolanie,
		&kanal.ParametryJSON, &multimodalny, &aktywny, &kanal.Kolejnosc, &kanal.Utworzono)
	if err != nil {
		return Kanal{}, err
	}
	kanal.KontoID = liczbaZKolumny(kontoID)
	kanal.PoswiadczenieOdwolanie = tekstZKolumny(odwolanie)
	kanal.Multimodalny = multimodalny == 1
	kanal.Aktywny = aktywny == 1
	return kanal, nil
}
