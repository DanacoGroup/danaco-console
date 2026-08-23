// Odpowiedzialność pliku: panel tłumaczenia (tabela `panel_tlumaczenia`) —
// zapis, odczyt pojedynczego panelu i wykaz paneli okna. Zmiany treści
// (tłumaczenie, tłumaczenie zwrotne, ton, stan) leżą w `tlumaczenie_tresc.go` —
// ten plik trzyma wyłącznie założenie panelu i odczyty, żeby jedna droga zmiany
// pola nie rozjechała się z drugą.
//
// Niezgodności nie są polem tego typu. Kontraktowe `TranslationPanel.Issues`
// warstwa wyższa składa z osobnego odczytu `Niezgodnosci` (`jakosc.go`) po
// `PanelID`. Panel i jego niezgodności to dwa byty w dwóch tabelach; trzymanie
// ich razem w jednej strukturze Go byłoby fałszywym obrazem schematu.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// PanelTlumaczenia to wiersz tabeli `panel_tlumaczenia`. `Stan` bierze wartości
// `TranslationStatus` z kontraktu wprost, bez przekładu. `Ton` ustawia komenda
// `translate.panel.tone.set`, a `TrescZwrotna` — `backtranslation.run`; obie
// kolumny mieszkają tu, bo to pola tego panelu, ale ten plik ich nie modyfikuje.
type PanelTlumaczenia struct {
	ID     int64
	Kod    string
	OknoID int64
	// OknoKod to identyfikator zewnętrzny okna — ten, którym okno wychodzi
	// kontraktem jako `TranslationPanel.windowId`. Bez niego z panelu nie da się
	// dojść do jego okna, bo odczyt okna przyjmuje kod zewnętrzny, a wiersz
	// panelu niesie sam klucz wewnętrzny. Kolumna jest doczytywana złączeniem,
	// nie zapisywana drugi raz — prawda o kodzie okna zostaje
	// w `okno_tlumaczenia`.
	OknoKod        string
	Jezyk          string
	Tresc          *string
	TrescOdwolanie *string
	Stan           string
	Ton            *string
	TrescZwrotna   *string
	// Migawka obiegu zatwierdzeń (migracja 162). Historię obiegu niesie tabela
	// `zatwierdzenie_panelu` (`tlumaczenie_kontrola.go`); te trzy pola są
	// wyłącznie stanem bieżącym, żeby odczyt panelu nie musiał dociągać
	// ostatniego wiersza obiegu przy każdym wykazie paneli okna.
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

	// Złączenie z oknem doczytuje kod zewnętrzny — patrz komentarz przy polu
	// `OknoKod`. Panel bez okna nie istnieje (klucz obcy NOT NULL), więc
	// złączenie wewnętrzne nie gubi wierszy.
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

// ZapiszPanel zakłada panel tłumaczenia dla wskazanego okna. Panel jest
// dodawany raz przez `translate.target.add` — kontrakt nie przewiduje
// nadpisania panelu po kodzie zewnętrznym (to robią osobne komendy zmiany
// treści, tonu i stanu), więc tu nie ma `ON CONFLICT`, jak przy zleceniu
// asystenta.
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
// i języku.
//
// Służy dwóm komendom obejmującym cały słownik. Kontrakt `glossary.apply`
// mówi, że puste `panelId` znaczy „wszystkie panele", a `glossary.occurrences`
// niesie sam termin, bez wskazania okna. Zakres bez zawężenia jest tu
// poprawny: słownik jest jeden na instalację, więc jego zastosowanie też.
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

// odczytajPanelTlumaczenia składa strukturę z jednego wiersza wyniku.
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
