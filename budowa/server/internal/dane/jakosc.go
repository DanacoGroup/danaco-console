// Odpowiedzialność pliku: niezgodności kontroli jakości panelu tłumaczenia
// (tabela `panel_tlumaczenia_niezgodnosc`) — faseta „jakość” modułu Translate.
package dane

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// NiezgodnoscTlumaczenia to wiersz tabeli `panel_tlumaczenia_niezgodnosc`.
// `Rodzaj` bierze wartości kontraktu (TranslationIssueKind) wprost, bez
// tłumaczenia na polski.
type NiezgodnoscTlumaczenia struct {
	ID        int64
	PanelID   int64
	Rodzaj    string
	Segment   *string
	Szczegol  *string
	Utworzono int64
}

const (
	usunNiezgodnosciPanelu = `DELETE FROM panel_tlumaczenia_niezgodnosc WHERE panel_id = ?`

	wstawNiezgodnoscPanelu = `INSERT INTO panel_tlumaczenia_niezgodnosc
	                          (panel_id, rodzaj, segment, szczegol, utworzono)
	                          VALUES (?, ?, ?, ?, ?)`

	listaNiezgodnosciPanelu = `SELECT id, panel_id, rodzaj, segment, szczegol, utworzono
	                           FROM panel_tlumaczenia_niezgodnosc WHERE panel_id = ?
	                           ORDER BY utworzono DESC, id DESC`
)

// ZapiszNiezgodnosci podmienia komplet niezgodności panelu wynikiem bieżącej
// kontroli jakości. Wykaz pusty zostawia panel bez zastrzeżeń — kontrola
// czysta jest stanem poprawnym, nie brakiem zapisu.
func (r *repozytoriumTlumaczen) ZapiszNiezgodnosci(ctx context.Context,
	panelID int64, niezgodnosci []NiezgodnoscTlumaczenia) error {

	teraz := time.Now().UnixMilli()
	return wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		czyszczenie, err := r.zapytania.wTransakcji(ctx, transakcja, usunNiezgodnosciPanelu)
		if err != nil {
			return err
		}
		if _, err := czyszczenie.ExecContext(ctx, panelID); err != nil {
			return fmt.Errorf("dane: nie można wyczyścić niezgodności panelu %d: %w", panelID, err)
		}
		wstawienie, err := r.zapytania.wTransakcji(ctx, transakcja, wstawNiezgodnoscPanelu)
		if err != nil {
			return err
		}
		for _, niezgodnosc := range niezgodnosci {
			utworzono := niezgodnosc.Utworzono
			if utworzono == 0 {
				utworzono = teraz
			}
			_, err := wstawienie.ExecContext(ctx, panelID, niezgodnosc.Rodzaj,
				tekstDoKolumny(niezgodnosc.Segment), tekstDoKolumny(niezgodnosc.Szczegol), utworzono)
			if err != nil {
				return fmt.Errorf("dane: nie można zapisać niezgodności %q panelu %d: %w",
					niezgodnosc.Rodzaj, panelID, err)
			}
		}
		return nil
	})
}

// Niezgodnosci zwraca niezgodności panelu z ostatniej kontroli jakości, uporządkowane od najnowszej do najstarszej.
func (r *repozytoriumTlumaczen) Niezgodnosci(ctx context.Context, panelID int64) ([]NiezgodnoscTlumaczenia, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaNiezgodnosciPanelu)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, panelID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać niezgodności panelu %d: %w", panelID, err)
	}
	defer wiersze.Close()

	lista := []NiezgodnoscTlumaczenia{}
	for wiersze.Next() {
		var niezgodnosc NiezgodnoscTlumaczenia
		var segment, szczegol sql.NullString
		err := wiersze.Scan(&niezgodnosc.ID, &niezgodnosc.PanelID, &niezgodnosc.Rodzaj,
			&segment, &szczegol, &niezgodnosc.Utworzono)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz niezgodności panelu %d: %w", panelID, err)
		}
		niezgodnosc.Segment = tekstZKolumny(segment)
		niezgodnosc.Szczegol = tekstZKolumny(szczegol)
		lista = append(lista, niezgodnosc)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt niezgodności panelu %d: %w", panelID, err)
	}
	return lista, nil
}
