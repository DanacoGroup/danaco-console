// Warstwa danych obsługuje odczyt warsztatu modułu Developer: dziennik
// budowań, punkty przerwania, kolekcje zapytań, połączenia bazodanowe,
// znaleziska skanów oraz wyniki testów i pokrycie.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

const (
	pobierzPrzebiegiBudowania = `SELECT ` + kolumnyPrzebieguBudowania + `
	                             FROM developer_budowanie
	                             WHERE okno_kod = ?
	                               AND (? = '' OR stan = ?)
	                             ORDER BY uruchomiono DESC, id DESC
	                             LIMIT CASE WHEN ? > 0 THEN ? ELSE -1 END`

	pobierzPrzebiegBudowania = `SELECT ` + kolumnyPrzebieguBudowania + `
	                            FROM developer_budowanie WHERE kod = ?`

	kolumnyPunktuPrzerwaniaDevelopera = `kod, okno_kod, sciezka, wiersz, rodzaj, warunek,
	                                     warunek_trafien, wpis, zweryfikowany`

	pobierzPunktyPrzerwaniaDevelopera = `SELECT ` + kolumnyPunktuPrzerwaniaDevelopera + `
	                                     FROM developer_punkt_przerwania
	                                     WHERE okno_kod = ? AND (? = '' OR sciezka = ?)
	                                     ORDER BY sciezka, wiersz`

	kolumnyKolekcjiApiDevelopera = `kod, okno_kod, nazwa, zapytania, srodowiska, zmieniono`

	pobierzKolekcjeApiDevelopera = `SELECT ` + kolumnyKolekcjiApiDevelopera + `
	                                FROM developer_kolekcja_api
	                                WHERE okno_kod = ? AND (? = '' OR kod = ?)
	                                  AND ` + WarunekKonta + `
	                                ORDER BY zmieniono DESC, id DESC`

	kolumnyPolaczeniaDanychDevelopera = `kod, okno_kod, nazwa, silnik, host, port, baza,
	                                     uzytkownik, poswiadczenie, tylko_odczyt`

	pobierzPolaczeniaDanychDevelopera = `SELECT ` + kolumnyPolaczeniaDanychDevelopera + `
	                                     FROM developer_polaczenie_danych
	                                     WHERE okno_kod = ? AND ` + WarunekKonta + `
	                                     ORDER BY nazwa`

	pobierzPolaczenieDanychDevelopera = `SELECT ` + kolumnyPolaczeniaDanychDevelopera + `
	                                     FROM developer_polaczenie_danych
	                                     WHERE kod = ? AND ` + WarunekKonta

	kolumnyZnaleziskaDevelopera = `z.kod, z.skan_kod, z.rodzaj, z.waga, z.tytul, z.opis,
	                               z.sciezka, z.wiersz, z.regula, z.cve, z.pakiet, z.wersja_naprawy`

	// Kolejność wagi jest jawna, bo alfabet ustawiłby ostrzeżenie przed błędem:
	// „error” < „info” < „warning” tekstowo, a Operator ma zobaczyć najcięższe
	// spostrzeżenia u góry wykazu.
	pobierzZnaleziskaDevelopera = `SELECT ` + kolumnyZnaleziskaDevelopera + `
	                               FROM developer_znalezisko z
	                               JOIN developer_skan s ON s.kod = z.skan_kod
	                               WHERE (? = '' OR z.skan_kod = ?)
	                                 AND (? = '' OR s.okno_kod = ?)
	                                 AND (? = '' OR z.rodzaj = ?)
	                                 AND (? = '' OR z.waga = ?)
	                               ORDER BY CASE z.waga WHEN 'error' THEN 0 WHEN 'warning' THEN 1
	                                                    ELSE 2 END,
	                                        s.uruchomiono DESC, z.id
	                               LIMIT CASE WHEN ? > 0 THEN ? ELSE -1 END`

	pobierzWynikiTestowDevelopera = `SELECT budowanie_kod, zestaw, nazwa, stan, czas_ms, tresc,
	                                        sciezka, wiersz
	                                 FROM developer_wynik_testu
	                                 WHERE budowanie_kod = ? AND (? = '' OR stan = ?)
	                                 ORDER BY id`

	pobierzPokrycieDevelopera = `SELECT budowanie_kod, sciezka, instrukcje, pokryte, procent,
	                                    wiersze_bez_pokrycia
	                             FROM developer_pokrycie
	                             WHERE budowanie_kod = ? AND (? = '' OR sciezka = ?)
	                             ORDER BY sciezka`
)

// Przebiegi zwraca dziennik budowań okna od najnowszego; stan pusty i limit
// niedodatni nie zawężają wykazu.
func (r *repozytoriumDevelopera) Przebiegi(ctx context.Context, oknoKod, stan string,
	limit int) ([]PrzebiegBudowania, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzPrzebiegiBudowania)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, oknoKod, stan, stan, limit, limit)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać dziennika budowań okna %q: %w", oknoKod, err)
	}
	defer wiersze.Close()

	przebiegi := make([]PrzebiegBudowania, 0, 16)
	for wiersze.Next() {
		przebieg, err := odczytajPrzebiegBudowania(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz dziennika budowań: %w", err)
		}
		przebiegi = append(przebiegi, przebieg)
	}
	return przebiegi, wiersze.Err()
}

// Przebieg zwraca jeden przebieg budowania po jego identyfikatorze,
// odczytany z tabeli developer_budowanie.
func (r *repozytoriumDevelopera) Przebieg(ctx context.Context, kod string) (PrzebiegBudowania, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzPrzebiegBudowania)
	if err != nil {
		return PrzebiegBudowania{}, err
	}
	przebieg, err := odczytajPrzebiegBudowania(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return PrzebiegBudowania{}, ErrBrakWiersza
	}
	if err != nil {
		return PrzebiegBudowania{}, fmt.Errorf("dane: nieczytelny przebieg budowania %q: %w", kod, err)
	}
	return przebieg, nil
}

// PunktyPrzerwania zwraca punkty przerwania postawione w oknie; pusta
// ścieżka pliku znaczy wszystkie pliki.
func (r *repozytoriumDevelopera) PunktyPrzerwania(ctx context.Context,
	oknoKod, sciezka string) ([]PunktPrzerwania, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzPunktyPrzerwaniaDevelopera)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, oknoKod, sciezka, sciezka)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać punktów przerwania okna %q: %w", oknoKod, err)
	}
	defer wiersze.Close()

	punkty := make([]PunktPrzerwania, 0, 8)
	for wiersze.Next() {
		var punkt PunktPrzerwania
		var zweryfikowany int64
		if err := wiersze.Scan(&punkt.Kod, &punkt.OknoKod, &punkt.Sciezka, &punkt.Wiersz,
			&punkt.Rodzaj, &punkt.Warunek, &punkt.WarunekTrafien, &punkt.Wpis,
			&zweryfikowany); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny punkt przerwania: %w", err)
		}
		punkt.Zweryfikowany = zweryfikowany != 0
		punkty = append(punkty, punkt)
	}
	return punkty, wiersze.Err()
}

// KolekcjeApi zwraca kolekcje zapytań HTTP okna; niepusty kod zawęża wykaz
// do jednej kolekcji zapytań.
func (r *repozytoriumDevelopera) KolekcjeApi(ctx context.Context,
	oknoKod, kod string) ([]KolekcjaApi, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzKolekcjeApiDevelopera)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, oknoKod, kod, kod, KontoOperatora(ctx))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać kolekcji zapytań okna %q: %w", oknoKod, err)
	}
	defer wiersze.Close()

	kolekcje := make([]KolekcjaApi, 0, 8)
	for wiersze.Next() {
		var kolekcja KolekcjaApi
		if err := wiersze.Scan(&kolekcja.Kod, &kolekcja.OknoKod, &kolekcja.Nazwa,
			&kolekcja.Zapytania, &kolekcja.Srodowiska, &kolekcja.Zmieniono); err != nil {
			return nil, fmt.Errorf("dane: nieczytelna kolekcja zapytań: %w", err)
		}
		kolekcje = append(kolekcje, kolekcja)
	}
	return kolekcje, wiersze.Err()
}

// PolaczeniaDanych zwraca połączenia bazodanowe okna w kolejności nazwy,
// bez treści hasła dostępowego.
func (r *repozytoriumDevelopera) PolaczeniaDanych(ctx context.Context,
	oknoKod string) ([]PolaczenieDanych, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzPolaczeniaDanychDevelopera)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, oknoKod, KontoOperatora(ctx))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać połączeń okna %q: %w", oknoKod, err)
	}
	defer wiersze.Close()

	polaczenia := make([]PolaczenieDanych, 0, 8)
	for wiersze.Next() {
		polaczenie, err := odczytajPolaczenieDanychDevelopera(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelne połączenie bazodanowe: %w", err)
		}
		polaczenia = append(polaczenia, polaczenie)
	}
	return polaczenia, wiersze.Err()
}

// PolaczenieDanychPoKodzie zwraca jedno połączenie bazodanowe po jego
// identyfikatorze wiersza w bazie.
func (r *repozytoriumDevelopera) PolaczenieDanychPoKodzie(ctx context.Context,
	kod string) (PolaczenieDanych, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzPolaczenieDanychDevelopera)
	if err != nil {
		return PolaczenieDanych{}, err
	}
	polaczenie, err := odczytajPolaczenieDanychDevelopera(
		polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return PolaczenieDanych{}, ErrBrakWiersza
	}
	if err != nil {
		return PolaczenieDanych{}, fmt.Errorf("dane: nieczytelne połączenie %q: %w", kod, err)
	}
	return polaczenie, nil
}

// Znaleziska zwraca spostrzeżenia skanów wedle filtru, od najcięższej wagi
// i najnowszego przebiegu skanu.
func (r *repozytoriumDevelopera) Znaleziska(ctx context.Context,
	filtr FiltrZnalezisk) ([]ZnaleziskoSkanu, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzZnaleziskaDevelopera)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx,
		filtr.SkanKod, filtr.SkanKod, filtr.OknoKod, filtr.OknoKod,
		filtr.Rodzaj, filtr.Rodzaj, filtr.Waga, filtr.Waga, filtr.Limit, filtr.Limit)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać znalezisk skanowania: %w", err)
	}
	defer wiersze.Close()

	znaleziska := make([]ZnaleziskoSkanu, 0, 16)
	for wiersze.Next() {
		var z ZnaleziskoSkanu
		if err := wiersze.Scan(&z.Kod, &z.SkanKod, &z.Rodzaj, &z.Waga, &z.Tytul, &z.Opis,
			&z.Sciezka, &z.Wiersz, &z.Regula, &z.Cve, &z.Pakiet, &z.WersjaNaprawy); err != nil {
			return nil, fmt.Errorf("dane: nieczytelne znalezisko skanowania: %w", err)
		}
		znaleziska = append(znaleziska, z)
	}
	return znaleziska, wiersze.Err()
}

// WynikiTestow zwraca wyniki testów przebiegu budowania okna; pusty stan
// znaczy wszystkie wyniki testu.
func (r *repozytoriumDevelopera) WynikiTestow(ctx context.Context,
	budowanieKod, stan string) ([]WynikTestu, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzWynikiTestowDevelopera)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, budowanieKod, stan, stan)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać wyników testów przebiegu %q: %w",
			budowanieKod, err)
	}
	defer wiersze.Close()

	wyniki := make([]WynikTestu, 0, 32)
	for wiersze.Next() {
		var w WynikTestu
		if err := wiersze.Scan(&w.BudowanieKod, &w.Zestaw, &w.Nazwa, &w.Stan, &w.CzasMs,
			&w.Tresc, &w.Sciezka, &w.Wiersz); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wynik testu: %w", err)
		}
		wyniki = append(wyniki, w)
	}
	return wyniki, wiersze.Err()
}

// Pokrycie zwraca pomiar pokrycia kodu przebiegu budowania; pusta ścieżka
// znaczy wszystkie pliki pokrycia.
func (r *repozytoriumDevelopera) Pokrycie(ctx context.Context,
	budowanieKod, sciezka string) ([]PokryciePliku, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzPokrycieDevelopera)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, budowanieKod, sciezka, sciezka)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać pokrycia przebiegu %q: %w",
			budowanieKod, err)
	}
	defer wiersze.Close()

	pokrycie := make([]PokryciePliku, 0, 16)
	for wiersze.Next() {
		var p PokryciePliku
		if err := wiersze.Scan(&p.BudowanieKod, &p.Sciezka, &p.Instrukcje, &p.Pokryte,
			&p.Procent, &p.WierszeBezPokrycia); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz pokrycia: %w", err)
		}
		pokrycie = append(pokrycie, p)
	}
	return pokrycie, wiersze.Err()
}

// odczytajPolaczenieDanychDevelopera składa połączenie PolaczenieDanych
// z jednego wiersza wyniku zapytania.
func odczytajPolaczenieDanychDevelopera(
	wiersz interface{ Scan(...any) error }) (PolaczenieDanych, error) {

	var polaczenie PolaczenieDanych
	var tylkoOdczyt int64
	err := wiersz.Scan(&polaczenie.Kod, &polaczenie.OknoKod, &polaczenie.Nazwa,
		&polaczenie.Silnik, &polaczenie.Host, &polaczenie.Port, &polaczenie.Baza,
		&polaczenie.Uzytkownik, &polaczenie.Poswiadczenie, &tylkoOdczyt)
	polaczenie.TylkoOdczyt = tylkoOdczyt != 0
	return polaczenie, err
}

// pobierzWersjePoKodzie czyta jedną migawkę pliku po jej identyfikatorze
// wiersza tabeli developer_wersja_pliku.
const pobierzWersjePoKodzie = `SELECT ` + kolumnyWersjiPliku + `
                               FROM developer_wersja_pliku
                               WHERE kod = ? AND ` + WarunekKonta

// WersjaPoKodzie zwraca jedną migawkę treści pliku po jej identyfikatorze
// wiersza w tabeli bazy danych.
func (r *repozytoriumDevelopera) WersjaPoKodzie(ctx context.Context, kod string) (WersjaPliku, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzWersjePoKodzie)
	if err != nil {
		return WersjaPliku{}, err
	}
	wersja, err := odczytajWersjePliku(polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return WersjaPliku{}, ErrBrakWiersza
	}
	if err != nil {
		return WersjaPliku{}, fmt.Errorf("dane: nieczytelna wersja pliku %q: %w", kod, err)
	}
	return wersja, nil
}
