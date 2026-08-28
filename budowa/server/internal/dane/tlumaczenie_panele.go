// Repozytorium panelu tłumaczenia obsługuje zapis, odczyt pojedynczego panelu
// i wykaz paneli okna w tabeli `panel_tlumaczenia`; zmiany treści leżą
// w pliku `tlumaczenie_tresc.go`.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// PanelTlumaczenia to wiersz tabeli `panel_tlumaczenia`, w którym pole Stan
// bierze wartości TranslationStatus z kontraktu wprost, bez przekładu.
type PanelTlumaczenia struct {
	ID     int64
	Kod    string
	OknoID int64
	// OknoKod to identyfikator zewnętrzny okna, doczytywany złączeniem, nie
	// zapisywany drugi raz.
	OknoKod        string
	Jezyk          string
	Tresc          *string
	TrescOdwolanie *string
	Stan           string
	Ton            *string
	TrescZwrotna   *string
	// Migawka obiegu zatwierdzeń z migracji 162; historię niesie tabela
	// `zatwierdzenie_panelu`.
	EtapZatwierdzenia *string
	Zatwierdzil       *string
	Zatwierdzono      *int64
	Zaktualizowano    int64
}

const (
	kolumnyPaneluTlumaczenia = `p.id, p.identyfikator_zewnetrzny, p.okno_id,
	                            o.identyfikator_zewnetrzny, p.jezyk, p.tresc,
	                            p.tresc_odwolanie, p.stan, p.ton, p.tresc_zwrotna,
	                            p.etap_zatwierdzenia, p.zatwierdzil, p.zatwierdzono,
	                            p.zaktualizowano`

	// Złączenie z oknem doczytuje kod zewnętrzny okna, ponieważ panel nie
	// zapisuje go drugi raz; panel bez okna nie istnieje, więc złączenie
	// wewnętrzne nie gubi wierszy.
	zrodloPaneluTlumaczenia = ` FROM panel_tlumaczenia p
	                            JOIN okno_tlumaczenia o ON o.id = p.okno_id`

	wstawPanelTlumaczenia = `INSERT INTO panel_tlumaczenia
	                         (identyfikator_zewnetrzny, okno_id, jezyk, tresc,
	                          tresc_odwolanie, stan, ton, tresc_zwrotna, zaktualizowano)
	                         VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	pobierzPanelTlumaczenia = `SELECT ` + kolumnyPaneluTlumaczenia + zrodloPaneluTlumaczenia +
		` WHERE p.identyfikator_zewnetrzny = ?`

	pobierzPaneleOkna = `SELECT ` + kolumnyPaneluTlumaczenia + zrodloPaneluTlumaczenia +
		` WHERE p.okno_id = ? ORDER BY p.jezyk`

	pobierzWszystkiePanele = `SELECT ` + kolumnyPaneluTlumaczenia + zrodloPaneluTlumaczenia +
		` ORDER BY p.okno_id, p.jezyk`
)

// ZapiszPanel zakłada panel tłumaczenia dla wskazanego okna i nie nadpisuje
// panelu istniejącego po kodzie zewnętrznym.
func (r *repozytoriumTlumaczen) ZapiszPanel(ctx context.Context, oknoID int64, panel PanelTlumaczenia) (PanelTlumaczenia, error) {
	if panel.Kod == "" {
		return PanelTlumaczenia{}, fmt.Errorf("dane: panel tłumaczenia bez identyfikatora")
	}
	if oknoID == 0 {
		return PanelTlumaczenia{}, fmt.Errorf("dane: panel tłumaczenia %q bez okna", panel.Kod)
	}
	if panel.Jezyk == "" {
		return PanelTlumaczenia{}, fmt.Errorf("dane: panel tłumaczenia %q bez języka docelowego", panel.Kod)
	}
	stan := panel.Stan
	if stan == "" {
		stan = "pending"
	}
	teraz := time.Now().UnixMilli()
	polecenie, err := r.zapytania.przygotuj(ctx, wstawPanelTlumaczenia)
	if err != nil {
		return PanelTlumaczenia{}, err
	}
	_, err = polecenie.ExecContext(ctx, panel.Kod, oknoID, panel.Jezyk,
		tekstDoKolumny(panel.Tresc), tekstDoKolumny(panel.TrescOdwolanie), stan,
		tekstDoKolumny(panel.Ton), tekstDoKolumny(panel.TrescZwrotna), teraz)
	if err != nil {
		return PanelTlumaczenia{}, fmt.Errorf("dane: nie można zapisać panelu tłumaczenia %q: %w", panel.Kod, err)
	}
	return r.Panel(ctx, panel.Kod)
}

// Panel zwraca panel tłumaczenia o wskazanym kodzie zewnętrznym. Brak wiersza
// wraca jako ErrBrakWiersza — warstwa wyższa odróżnia „nie ma” od „odczyt się
// nie powiódł”.
func (r *repozytoriumTlumaczen) Panel(ctx context.Context, kod string) (PanelTlumaczenia, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzPanelTlumaczenia)
	if err != nil {
		return PanelTlumaczenia{}, err
	}
	panel, err := odczytajPanelTlumaczenia(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return PanelTlumaczenia{}, ErrBrakWiersza
	}
	if err != nil {
		return PanelTlumaczenia{}, fmt.Errorf("dane: nieczytelny wiersz panelu tłumaczenia %q: %w", kod, err)
	}
	return panel, nil
}

// Panele zwraca wszystkie panele okna, po języku docelowym — `source.set`
// zasila kontraktowe `Panels []TranslationPanel` tym odczytem po uruchomieniu
// aktualizacji wszystkich paneli okna naraz.
func (r *repozytoriumTlumaczen) Panele(ctx context.Context, oknoID int64) ([]PanelTlumaczenia, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzPaneleOkna)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, oknoID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać paneli okna %d: %w", oknoID, err)
	}
	defer wiersze.Close()

	lista := []PanelTlumaczenia{}
	for wiersze.Next() {
		panel, err := odczytajPanelTlumaczenia(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz panelu tłumaczenia okna %d: %w", oknoID, err)
		}
		lista = append(lista, panel)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt paneli okna %d: %w", oknoID, err)
	}
	return lista, nil
}

// WszystkiePanele oddaje panele całej instalacji, uporządkowane po oknie
// i języku, bez zawężenia do jednego okna.
func (r *repozytoriumTlumaczen) WszystkiePanele(ctx context.Context) ([]PanelTlumaczenia, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzWszystkiePanele)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać paneli tłumaczenia: %w", err)
	}
	defer wiersze.Close()

	lista := []PanelTlumaczenia{}
	for wiersze.Next() {
		panel, err := odczytajPanelTlumaczenia(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz panelu tłumaczenia: %w", err)
		}
		lista = append(lista, panel)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt paneli tłumaczenia: %w", err)
	}
	return lista, nil
}

// odczytajPanelTlumaczenia składa strukturę PanelTlumaczenia z jednego
// wiersza wyniku zapytania do bazy.
func odczytajPanelTlumaczenia(wiersz skaner) (PanelTlumaczenia, error) {
	var panel PanelTlumaczenia
	var tresc, trescOdwolanie, ton, trescZwrotna, etap, zatwierdzil sql.NullString
	var zatwierdzono sql.NullInt64
	err := wiersz.Scan(&panel.ID, &panel.Kod, &panel.OknoID, &panel.OknoKod, &panel.Jezyk, &tresc,
		&trescOdwolanie, &panel.Stan, &ton, &trescZwrotna, &etap, &zatwierdzil, &zatwierdzono,
		&panel.Zaktualizowano)
	if err != nil {
		return PanelTlumaczenia{}, err
	}
	panel.Tresc = tekstZKolumny(tresc)
	panel.TrescOdwolanie = tekstZKolumny(trescOdwolanie)
	panel.Ton = tekstZKolumny(ton)
	panel.TrescZwrotna = tekstZKolumny(trescZwrotna)
	panel.EtapZatwierdzenia = tekstZKolumny(etap)
	panel.Zatwierdzil = tekstZKolumny(zatwierdzil)
	panel.Zatwierdzono = liczbaZKolumny(zatwierdzono)
	return panel, nil
}
