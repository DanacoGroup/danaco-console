// Plik prowadzi obszar Apps — dobudowę Architecture Designera: historię wersji układu oraz adnotacje projektowe kanwy;
// wiersz historii dopisuje ZapiszArchitekture w tej samej transakcji, w której wymienia komponenty.
package dane

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// WersjaArchitekturyApp to wiersz tabeli `wersja_architektury_apps` niosący migawkę układu komponentów.
type WersjaArchitekturyApp struct {
	Wersja            int
	KodArchitektury   string
	LiczbaKomponentow int
	Roznica           *string
	Utworzono         string
}

// AdnotacjaArchitekturyApp to wiersz tabeli `adnotacja_architektury_apps` niosący notatkę projektową kanwy.
type AdnotacjaArchitekturyApp struct {
	ID             int64
	Kod            string
	ArchitekturaID int64
	KomponentKod   *string
	ZaleznoscZ     *string
	ZaleznoscDo    *string
	Tresc          string
	Utworzono      string
	Zaktualizowano string
}

const (
	wstawWersjeArchitekturyApp = `INSERT INTO wersja_architektury_apps
	                              (architektura_id, wersja, liczba_komponentow, roznica)
	                              VALUES (?, ?, ?, ?)
	                              ON CONFLICT(architektura_id, wersja) DO UPDATE SET
	                                  liczba_komponentow = excluded.liczba_komponentow,
	                                  roznica = excluded.roznica`

	kolumnyAdnotacjiApp = `id, identyfikator_zewnetrzny, architektura_id, komponent_kod,
	                       zaleznosc_z, zaleznosc_do, tresc, utworzono, zaktualizowano`

	listaAdnotacjiApp = `SELECT ` + kolumnyAdnotacjiApp + `
	                     FROM adnotacja_architektury_apps
	                     WHERE architektura_id = ? ORDER BY id`
)

// Wersja i adnotacja własnej kolumny konta nie mają; granica idzie drogą po
// `architektura_id` do `architektura_apps` (kolumna `konto_id` z migracji 484).
var (
	warunekKontaArchitekturyApp = strings.ReplaceAll(WarunekKonta, "konto_id", "a.konto_id")

	listaWersjiArchitekturyApp = `SELECT w.wersja, a.identyfikator_zewnetrzny,
	                                     w.liczba_komponentow, w.roznica, w.utworzono
	                              FROM wersja_architektury_apps w
	                              JOIN architektura_apps a ON a.id = w.architektura_id
	                              WHERE w.architektura_id = ? AND ` + warunekKontaArchitekturyApp + `
	                              ORDER BY w.wersja DESC`

	adnotacjaWKoncie = ` EXISTS (SELECT 1 FROM architektura_apps a
	                             WHERE a.id = adnotacja_architektury_apps.architektura_id
	                               AND ` + warunekKontaArchitekturyApp + `)`

	// Identyfikator adnotacji jest jednoznaczny w całej tabeli — gałąź konfliktu bez
	// warunku konta sięgnęłaby adnotacji konta cudzego.
	zapiszAdnotacjeApp = `INSERT INTO adnotacja_architektury_apps
	                      (identyfikator_zewnetrzny, architektura_id, komponent_kod,
	                       zaleznosc_z, zaleznosc_do, tresc)
	                      VALUES (?, ?, ?, ?, ?, ?)
	                      ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                          komponent_kod = excluded.komponent_kod,
	                          zaleznosc_z = excluded.zaleznosc_z,
	                          zaleznosc_do = excluded.zaleznosc_do,
	                          tresc = excluded.tresc,
	                          zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                      WHERE ` + adnotacjaWKoncie

	pobierzAdnotacjeApp = `SELECT ` + kolumnyAdnotacjiApp + `
	                       FROM adnotacja_architektury_apps
	                       WHERE identyfikator_zewnetrzny = ? AND ` + adnotacjaWKoncie
)

func (r *repozytoriumAplikacji) WersjeArchitekturyApp(ctx context.Context,
	architekturaID int64) ([]WersjaArchitekturyApp, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaWersjiArchitekturyApp)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, architekturaID, KontoOperatora(ctx))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać wersji architektury %d: %w", architekturaID, err)
	}
	defer wiersze.Close()

	lista := []WersjaArchitekturyApp{}
	for wiersze.Next() {
		var wersja WersjaArchitekturyApp
		var roznica sql.NullString
		err := wiersze.Scan(&wersja.Wersja, &wersja.KodArchitektury,
			&wersja.LiczbaKomponentow, &roznica, &wersja.Utworzono)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz wersji architektury: %w", err)
		}
		wersja.Roznica = tekstZKolumny(roznica)
		lista = append(lista, wersja)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt wersji architektury %d: %w", architekturaID, err)
	}
	return lista, nil
}

func (r *repozytoriumAplikacji) ZapiszAdnotacjeApp(ctx context.Context,
	adnotacja AdnotacjaArchitekturyApp) (AdnotacjaArchitekturyApp, error) {

	if adnotacja.Kod == "" {
		return AdnotacjaArchitekturyApp{}, fmt.Errorf("dane: adnotacja architektury bez identyfikatora")
	}
	if adnotacja.ArchitekturaID == 0 {
		return AdnotacjaArchitekturyApp{}, fmt.Errorf("dane: adnotacja %q bez architektury", adnotacja.Kod)
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszAdnotacjeApp)
	if err != nil {
		return AdnotacjaArchitekturyApp{}, err
	}
	wynik, err := polecenie.ExecContext(ctx, adnotacja.Kod, adnotacja.ArchitekturaID,
		tekstDoKolumny(adnotacja.KomponentKod), tekstDoKolumny(adnotacja.ZaleznoscZ),
		tekstDoKolumny(adnotacja.ZaleznoscDo), adnotacja.Tresc, KontoOperatora(ctx))
	if err != nil {
		return AdnotacjaArchitekturyApp{}, fmt.Errorf("dane: nie można zapisać adnotacji %q: %w",
			adnotacja.Kod, err)
	}
	if err := sprawdzTrafienieZapisu(wynik, "adnotacja architektury", adnotacja.Kod); err != nil {
		return AdnotacjaArchitekturyApp{}, err
	}
	return r.AdnotacjaApp(ctx, adnotacja.Kod)
}

func (r *repozytoriumAplikacji) AdnotacjaApp(ctx context.Context, kod string) (AdnotacjaArchitekturyApp, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzAdnotacjeApp)
	if err != nil {
		return AdnotacjaArchitekturyApp{}, err
	}
	adnotacja, err := odczytajAdnotacjeApp(polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)))
	if err == sql.ErrNoRows {
		return AdnotacjaArchitekturyApp{}, ErrBrakWiersza
	}
	if err != nil {
		return AdnotacjaArchitekturyApp{}, fmt.Errorf("dane: nieczytelna adnotacja %q: %w", kod, err)
	}
	return adnotacja, nil
}

func (r *repozytoriumAplikacji) AdnotacjeApp(ctx context.Context,
	architekturaID int64) ([]AdnotacjaArchitekturyApp, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaAdnotacjiApp)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, architekturaID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać adnotacji architektury %d: %w", architekturaID, err)
	}
	defer wiersze.Close()

	lista := []AdnotacjaArchitekturyApp{}
	for wiersze.Next() {
		adnotacja, err := odczytajAdnotacjeApp(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz adnotacji architektury: %w", err)
		}
		lista = append(lista, adnotacja)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt adnotacji architektury %d: %w", architekturaID, err)
	}
	return lista, nil
}

func odczytajAdnotacjeApp(wiersz skaner) (AdnotacjaArchitekturyApp, error) {
	var adnotacja AdnotacjaArchitekturyApp
	var komponent, zaleznoscZ, zaleznoscDo sql.NullString
	err := wiersz.Scan(&adnotacja.ID, &adnotacja.Kod, &adnotacja.ArchitekturaID,
		&komponent, &zaleznoscZ, &zaleznoscDo, &adnotacja.Tresc,
		&adnotacja.Utworzono, &adnotacja.Zaktualizowano)
	if err != nil {
		return AdnotacjaArchitekturyApp{}, err
	}
	adnotacja.KomponentKod = tekstZKolumny(komponent)
	adnotacja.ZaleznoscZ = tekstZKolumny(zaleznoscZ)
	adnotacja.ZaleznoscDo = tekstZKolumny(zaleznoscDo)
	return adnotacja, nil
}
