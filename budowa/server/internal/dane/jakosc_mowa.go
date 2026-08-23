// Odpowiedzialność pliku: ślad syntezy mowy panelu tłumaczenia (odsłuch),
// tabela `panel_tlumaczenia_synteza_mowy` (`store/migracja_055_jakosc_i_mowa.sql`).
// Zasila komendę `translate.speech.synthesize`; plik osobny od `jakosc.go`
// i `jakosc_eksport.go`.
//
// RDZEŃ NIE SYNTEZUJE MOWY — TEN SAM BRAK CO W MODULE ASSISTANT
// (`wpis_dziennika_asystenta.nagranie_odnosnik`). `NagranieOdnosnik` niesie
// odwołanie do pliku dostarczonego z zewnątrz, jeśli kiedykolwiek powstanie;
// metoda zapisu nie dorabia mu wartości domyślnej — brak zostaje NULL, bo
// rdzeń nie syntezuje i zmyślona ścieżka byłaby obietnicą bez pokrycia.
package dane

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// SyntezaMowy to wiersz tabeli `panel_tlumaczenia_synteza_mowy` — jeden ślad
// żądania odsłuchu panelu. `NagranieOdnosnik` bywa NULL, dopóki nagranie
// realnie nie istnieje (rdzeń go nie wytwarza).
type SyntezaMowy struct {
	ID               int64
	PanelID          int64
	NagranieOdnosnik *string
	Utworzono        int64
}

const (
	kolumnySyntezyMowy = `id, panel_id, nagranie_odnosnik, utworzono`

	wstawSyntezeMowy = `INSERT INTO panel_tlumaczenia_synteza_mowy
	                    (panel_id, nagranie_odnosnik, utworzono)
	                    VALUES (?, ?, ?)`

	pobierzSyntezyPanelu = `SELECT ` + kolumnySyntezyMowy + `
	                        FROM panel_tlumaczenia_synteza_mowy
	                        WHERE panel_id = ?
	                        ORDER BY utworzono DESC, id DESC`
)

// ZapiszSyntezeMowy dopisuje ślad żądania syntezy mowy dla panelu. Panel bez
// identyfikatora nie ma do czego przypiąć śladu — odrzucamy go od razu,
// zamiast pozwolić kluczowi obcemu bazy zgłosić to dopiero przy zapisie.
func (r *repozytoriumTlumaczen) ZapiszSyntezeMowy(ctx context.Context, synteza SyntezaMowy) (SyntezaMowy, error) {
	if synteza.PanelID == 0 {
		return SyntezaMowy{}, fmt.Errorf("dane: synteza mowy bez panelu")
	}
	utworzono := synteza.Utworzono
	if utworzono == 0 {
		utworzono = time.Now().UnixMilli()
	}
	polecenie, err := r.zapytania.przygotuj(ctx, wstawSyntezeMowy)
	if err != nil {
		return SyntezaMowy{}, err
	}
	wynik, err := polecenie.ExecContext(ctx, synteza.PanelID, tekstDoKolumny(synteza.NagranieOdnosnik), utworzono)
	if err != nil {
		return SyntezaMowy{}, fmt.Errorf("dane: nie można zapisać syntezy mowy panelu %d: %w", synteza.PanelID, err)
	}
	id, err := wynik.LastInsertId()
	if err != nil {
		return SyntezaMowy{}, fmt.Errorf("dane: nie można odczytać identyfikatora syntezy mowy panelu %d: %w", synteza.PanelID, err)
	}
	synteza.ID = id
	synteza.Utworzono = utworzono
	return synteza, nil
}

// Syntezy zwraca ślady syntezy mowy panelu od najnowszego. Brak panelu nie
// jest tu błędem — pusty wykaz znaczy po prostu „jeszcze żadnego odsłuchu”,
// zgodnie z ErrBrakWiersza zarezerwowanym dla odczytu pojedynczego bytu.
func (r *repozytoriumTlumaczen) Syntezy(ctx context.Context, panelID int64) ([]SyntezaMowy, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzSyntezyPanelu)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, panelID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać syntez mowy panelu %d: %w", panelID, err)
	}
	defer wiersze.Close()
	syntezy := make([]SyntezaMowy, 0, 8)
	for wiersze.Next() {
		synteza, err := odczytajSyntezeMowy(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz syntezy mowy panelu %d: %w", panelID, err)
		}
		syntezy = append(syntezy, synteza)
	}
	return syntezy, wiersze.Err()
}

// odczytajSyntezeMowy przenosi jeden wiersz zapytania do bytu obszaru.
func odczytajSyntezeMowy(s skaner) (SyntezaMowy, error) {
	var synteza SyntezaMowy
	var nagranieOdnosnik sql.NullString
	err := s.Scan(&synteza.ID, &synteza.PanelID, &nagranieOdnosnik, &synteza.Utworzono)
	if err != nil {
		return SyntezaMowy{}, err
	}
	synteza.NagranieOdnosnik = tekstZKolumny(nagranieOdnosnik)
	return synteza, nil
}
