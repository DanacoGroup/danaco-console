// Odpowiedzialność pliku: odczyt pozycji katalogu ustawień wraz z ich
// dopuszczalnymi wartościami, poziomami zasięgu i osiami. Pozycja katalogu
// odpowiada strukturze SettingDefinition kontraktu — klient buduje z niej pole
// formularza i nie zna ani jednego klucza z osobna.
//
// Zbiory poboczne (opcje, zasięgi, osie) czytane są trzema zapytaniami zbiorczo
// i dokładane do pozycji w pamięci. Zapytania per pozycja nie ma.
package dane

import (
	"context"
	"database/sql"
	"fmt"

	"danacoconsole/shared"
)

const (
	listaDefinicjiUstawien = `SELECT d.klucz, k.kod, d.nazwa, d.opis, d.rodzaj_wartosci,
	                                 d.wartosc_domyslna, d.minimum, d.maksimum, d.skok,
	                                 d.wzorzec, d.jednostka, d.podpowiedz, d.wymagane,
	                                 d.wymaga_restartu, d.widoczne_gdy_klucz,
	                                 d.widoczne_gdy_wartosc, d.kolejnosc, d.aktywna
	                            FROM definicja_ustawienia d
	                            JOIN kategoria_ustawien k ON k.id = d.kategoria_id
	                           WHERE (? = 0 OR d.aktywna = 1)
	                           ORDER BY k.kolejnosc, d.kolejnosc, d.klucz`

	listaOpcjiUstawien = `SELECT d.klucz, o.wartosc, o.etykieta, o.opis, o.kolejnosc
	                        FROM opcja_ustawienia o
	                        JOIN definicja_ustawienia d ON d.id = o.definicja_id
	                       ORDER BY d.klucz, o.kolejnosc, o.wartosc`

	listaZasiegowDefinicji = `SELECT d.klucz, p.kod
	                            FROM definicja_ustawienia_zasieg z
	                            JOIN definicja_ustawienia d ON d.id = z.definicja_id
	                            JOIN poziom_zasiegu p ON p.id = z.poziom_zasiegu_id
	                           ORDER BY d.klucz, p.pierwszenstwo`

	listaOsiDefinicji = `SELECT d.klucz, o.kod
	                       FROM definicja_ustawienia_os w
	                       JOIN definicja_ustawienia d ON d.id = w.definicja_id
	                       JOIN os_zasiegu o ON o.kod = w.os
	                      ORDER BY d.klucz, o.pierwszenstwo`
)

// Definicje zwraca pozycje katalogu wraz ze zbiorami pobocznymi.
func (r *repozytoriumKatalogUstawien) Definicje(ctx context.Context,
	tylkoAktywne bool) ([]shared.SettingDefinition, error) {

	pozycje, err := r.pozycje(ctx, tylkoAktywne)
	if err != nil {
		return nil, err
	}
	opcje, err := r.opcje(ctx)
	if err != nil {
		return nil, err
	}
	zasiegi, err := r.pary(ctx, listaZasiegowDefinicji, "zasięgów pozycji katalogu")
	if err != nil {
		return nil, err
	}
	osie, err := r.pary(ctx, listaOsiDefinicji, "osi pozycji katalogu")
	if err != nil {
		return nil, err
	}
	for indeks := range pozycje {
		klucz := pozycje[indeks].Key
		pozycje[indeks].Options = opcje[klucz]
		pozycje[indeks].AllowedScopes = zasiegiKontraktu(zasiegi[klucz])
		pozycje[indeks].AllowedAxes = osieKontraktu(osie[klucz])
	}
	return pozycje, nil
}

// pozycje odczytuje same wiersze pozycji, bez zbiorów pobocznych.
func (r *repozytoriumKatalogUstawien) pozycje(ctx context.Context,
	tylkoAktywne bool) ([]shared.SettingDefinition, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaDefinicjiUstawien)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, liczbaLogiczna(tylkoAktywne))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać pozycji katalogu ustawień: %w", err)
	}
	defer wiersze.Close()

	lista := []shared.SettingDefinition{}
	for wiersze.Next() {
		pozycja, err := odczytajPozycjeKatalogu(wiersze)
		if err != nil {
			return nil, err
		}
		lista = append(lista, pozycja)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt pozycji katalogu ustawień: %w", err)
	}
	return lista, nil
}

// opcje zwraca dopuszczalne wartości pozycji, zgrupowane kluczem ustawienia.
func (r *repozytoriumKatalogUstawien) opcje(ctx context.Context) (map[string][]shared.SettingOption, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaOpcjiUstawien)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać opcji katalogu ustawień: %w", err)
	}
	defer wiersze.Close()

	zebrane := map[string][]shared.SettingOption{}
	for wiersze.Next() {
		var klucz, opis string
		var opcja shared.SettingOption
		if err := wiersze.Scan(&klucz, &opcja.Value, &opcja.Label, &opis, &opcja.Order); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz opcji katalogu ustawień: %w", err)
		}
		opcja.Description = tekstNiepusty(opis)
		zebrane[klucz] = append(zebrane[klucz], opcja)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt opcji katalogu ustawień: %w", err)
	}
	return zebrane, nil
}

// pary odczytuje relację klucz ustawienia → kod słownika (poziom zasięgu, oś).
func (r *repozytoriumKatalogUstawien) pary(ctx context.Context, zapytanie,
	obszar string) (map[string][]string, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać %s: %w", obszar, err)
	}
	defer wiersze.Close()

	zebrane := map[string][]string{}
	for wiersze.Next() {
		var klucz, kod string
		if err := wiersze.Scan(&klucz, &kod); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz %s: %w", obszar, err)
		}
		zebrane[klucz] = append(zebrane[klucz], kod)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt %s: %w", obszar, err)
	}
	return zebrane, nil
}

// odczytajPozycjeKatalogu składa pozycję kontraktu z jednego wiersza.
func odczytajPozycjeKatalogu(wiersz skaner) (shared.SettingDefinition, error) {
	var pozycja shared.SettingDefinition
	var opis, domyslna, wzorzec, jednostka, podpowiedz, gdyKlucz, gdyWartosc string
	var minimum, maksimum, skok sql.NullFloat64
	var wymagane, wymagaRestartu, aktywna int
	err := wiersz.Scan(&pozycja.Key, &pozycja.CategoryId, &pozycja.Name, &opis,
		&pozycja.ValueType, &domyslna, &minimum, &maksimum, &skok, &wzorzec, &jednostka,
		&podpowiedz, &wymagane, &wymagaRestartu, &gdyKlucz, &gdyWartosc,
		&pozycja.Order, &aktywna)
	if err != nil {
		return shared.SettingDefinition{}, fmt.Errorf("dane: nieczytelny wiersz pozycji katalogu ustawień: %w", err)
	}
	pozycja.Description = tekstNiepusty(opis)
	pozycja.DefaultValue = wartoscDomyslnaJSON(domyslna, pozycja.ValueType)
	pozycja.Minimum = liczbaRzeczywistaZKolumny(minimum)
	pozycja.Maximum = liczbaRzeczywistaZKolumny(maksimum)
	pozycja.Step = liczbaRzeczywistaZKolumny(skok)
	pozycja.Pattern = tekstNiepusty(wzorzec)
	pozycja.Unit = tekstNiepusty(jednostka)
	pozycja.Placeholder = tekstNiepusty(podpowiedz)
	pozycja.Required = wymagane == 1
	pozycja.RestartRequired = wartoscLogiczna(wymagaRestartu == 1)
	pozycja.VisibleWhenKey = tekstNiepusty(gdyKlucz)
	pozycja.VisibleWhenValue = tekstNiepusty(gdyWartosc)
	pozycja.Enabled = aktywna == 1
	return pozycja, nil
}
