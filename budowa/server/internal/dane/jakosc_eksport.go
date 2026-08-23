// Odpowiedzialność pliku: ślady eksportu panelu tłumaczenia (tabela
// `panel_tlumaczenia_eksport`) — część `RepozytoriumTlumaczen`, deklarowanego
// w całości w `tlumaczenie.go`. Tu wyłącznie zapis i odczyt śladu eksportu; typ,
// interfejs i konstruktor leżą tam.
//
// Rdzeń nie ma magazynu plików binarnych, więc `translate.panel.export` nie
// wytwarza pliku i `plik_odnosnik` zostaje NULL, dopóki plik realnie nie powstał
// poza rdzeniem. Zapisywany jest wyłącznie ślad zlecenia eksportu: w jakim
// formacie i kiedy.
//
// Czas jest liczbą milisekund epoki, wzorem `dane/asystent.go`.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// EksportPanelu to wiersz `panel_tlumaczenia_eksport` — ślad zlecenia
// eksportu panelu. `Format` niesie wartości ExportFormat kontraktu wprost
// (pdf/docx/markdown/html/txt), bez przekładu na polski. `PlikOdnosnik` jest
// odwołaniem do pliku wyniku, jeśli realnie powstał — rdzeń go nie wytwarza,
// więc zwykle zostaje NULL.
type EksportPanelu struct {
	ID           int64
	PanelID      int64
	Format       string
	PlikOdnosnik *string
	Utworzono    int64
}

const (
	kolumnyEksportuPanelu = `id, panel_id, format, plik_odnosnik, utworzono`

	wstawEksportPanelu = `INSERT INTO panel_tlumaczenia_eksport
	                      (panel_id, format, plik_odnosnik, utworzono)
	                      VALUES (?, ?, ?, ?)`

	pobierzEksportPanelu = `SELECT ` + kolumnyEksportuPanelu + `
	                        FROM panel_tlumaczenia_eksport WHERE id = ?`

	listaEksportowPanelu = `SELECT ` + kolumnyEksportuPanelu + `
	                        FROM panel_tlumaczenia_eksport
	                        WHERE panel_id = ? ORDER BY utworzono DESC, id DESC`

	istniejePanelDlaEksportu = `SELECT id FROM panel_tlumaczenia WHERE id = ?`
)

// ZapiszEksportPanelu dokłada ślad zlecenia eksportu panelu: każde
// wywołanie `translate.panel.export` dopisuje nowy wiersz historii, nic nie
// nadpisuje — ślad ma pokazywać wszystkie dotychczasowe eksporty panelu, nie
// tylko ostatni. Panel, którego nie ma, wraca jako ErrBrakWiersza — cicha
// zgoda na eksport bytu, którego nie ma, byłaby potwierdzeniem czynności,
// która się nie odbyła.
func (r *repozytoriumTlumaczen) ZapiszEksportPanelu(ctx context.Context, eksport EksportPanelu) (EksportPanelu, error) {
	if eksport.PanelID == 0 {
		return EksportPanelu{}, fmt.Errorf("dane: eksport panelu bez identyfikatora panelu")
	}
	if eksport.Format == "" {
		return EksportPanelu{}, fmt.Errorf("dane: eksport panelu %d bez formatu", eksport.PanelID)
	}
	sprawdzenie, err := r.zapytania.przygotuj(ctx, istniejePanelDlaEksportu)
	if err != nil {
		return EksportPanelu{}, err
	}
	var istniejeID int64
	err = sprawdzenie.QueryRowContext(ctx, eksport.PanelID).Scan(&istniejeID)
	if errors.Is(err, sql.ErrNoRows) {
		return EksportPanelu{}, fmt.Errorf("dane: panel tłumaczenia %d nie istnieje: %w", eksport.PanelID, ErrBrakWiersza)
	}
	if err != nil {
		return EksportPanelu{}, fmt.Errorf("dane: nie można znaleźć panelu tłumaczenia %d: %w", eksport.PanelID, err)
	}

	utworzono := eksport.Utworzono
	if utworzono == 0 {
		utworzono = time.Now().UnixMilli()
	}
	wstawienie, err := r.zapytania.przygotuj(ctx, wstawEksportPanelu)
	if err != nil {
		return EksportPanelu{}, err
	}
	wynik, err := wstawienie.ExecContext(ctx, eksport.PanelID, eksport.Format, tekstDoKolumny(eksport.PlikOdnosnik), utworzono)
	if err != nil {
		return EksportPanelu{}, fmt.Errorf("dane: nie można zapisać eksportu panelu %d: %w", eksport.PanelID, err)
	}
	nowyID, err := wynik.LastInsertId()
	if err != nil {
		return EksportPanelu{}, fmt.Errorf("dane: nie można ustalić identyfikatora eksportu panelu %d: %w", eksport.PanelID, err)
	}

	odczyt, err := r.zapytania.przygotuj(ctx, pobierzEksportPanelu)
	if err != nil {
		return EksportPanelu{}, err
	}
	zapisany, err := odczytajEksportPanelu(odczyt.QueryRowContext(ctx, nowyID))
	if err != nil {
		return EksportPanelu{}, fmt.Errorf("dane: nie można odczytać zapisanego eksportu panelu %d: %w", eksport.PanelID, err)
	}
	return zapisany, nil
}

// EksportyPanelu zwraca ślady eksportu panelu, od najświeższego. Panel bez
// żadnego eksportu wraca jako wykaz pusty, nie błąd — brak eksportów jest
// stanem startowym panelu, nie usterką.
func (r *repozytoriumTlumaczen) EksportyPanelu(ctx context.Context, panelID int64) ([]EksportPanelu, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaEksportowPanelu)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, panelID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać eksportów panelu %d: %w", panelID, err)
	}
	defer wiersze.Close()

	lista := []EksportPanelu{}
	for wiersze.Next() {
		eksport, err := odczytajEksportPanelu(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz eksportu panelu %d: %w", panelID, err)
		}
		lista = append(lista, eksport)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt eksportów panelu %d: %w", panelID, err)
	}
	return lista, nil
}

// odczytajEksportPanelu składa strukturę z jednego wiersza.
func odczytajEksportPanelu(wiersz skaner) (EksportPanelu, error) {
	var eksport EksportPanelu
	var plikOdnosnik sql.NullString
	err := wiersz.Scan(&eksport.ID, &eksport.PanelID, &eksport.Format, &plikOdnosnik, &eksport.Utworzono)
	if err != nil {
		return EksportPanelu{}, err
	}
	eksport.PlikOdnosnik = tekstZKolumny(plikOdnosnik)
	return eksport, nil
}
