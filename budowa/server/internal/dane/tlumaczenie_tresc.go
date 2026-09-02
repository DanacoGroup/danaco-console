// Plik zmienia treść panelu tłumaczenia: tłumaczenie, tłumaczenie zwrotne, ton
// i stan, jako metody typu repozytoriumTlumaczen z pliku tlumaczenie.go.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

// Panel nie niesie konta (migracja 484); granica idzie przez okno_tlumaczenia.
var granicaOknaPanelu = ` AND EXISTS (SELECT 1 FROM okno_tlumaczenia k
                            WHERE k.id = panel_tlumaczenia.okno_id AND ` +
	strings.ReplaceAll(WarunekKonta, "konto_id", "k.konto_id") + `)`

var (
	ustawTrescPanelu = `UPDATE panel_tlumaczenia
	                    SET tresc = ?, tresc_odwolanie = ?, zaktualizowano = ?
	                    WHERE identyfikator_zewnetrzny = ?` + granicaOknaPanelu

	ustawTrescZwrotnaPanelu = `UPDATE panel_tlumaczenia
	                           SET tresc_zwrotna = ?, zaktualizowano = ?
	                           WHERE identyfikator_zewnetrzny = ?` + granicaOknaPanelu

	ustawTonPanelu = `UPDATE panel_tlumaczenia
	                  SET ton = ?, zaktualizowano = ?
	                  WHERE identyfikator_zewnetrzny = ?` + granicaOknaPanelu

	ustawStanPanelu = `UPDATE panel_tlumaczenia
	                   SET stan = ?, zaktualizowano = ?
	                   WHERE identyfikator_zewnetrzny = ?` + granicaOknaPanelu
)

const istniejePanelTlumaczenia = `SELECT 1 FROM panel_tlumaczenia WHERE identyfikator_zewnetrzny = ?`

func (r *repozytoriumTlumaczen) UstawTlumaczenie(ctx context.Context, kodPanelu string,
	tresc, odwolanie *string) (PanelTlumaczenia, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, ustawTrescPanelu)
	if err != nil {
		return PanelTlumaczenia{}, err
	}
	wynik, err := polecenie.ExecContext(ctx, tekstDoKolumny(tresc), tekstDoKolumny(odwolanie),
		time.Now().UnixMilli(), kodPanelu, KontoOperatora(ctx))
	if err != nil {
		return PanelTlumaczenia{}, fmt.Errorf("dane: nie można ustawić tłumaczenia panelu %q: %w", kodPanelu, err)
	}
	if err := r.trafieniePanelu(ctx, wynik, kodPanelu); err != nil {
		return PanelTlumaczenia{}, err
	}
	return r.Panel(ctx, kodPanelu)
}

func (r *repozytoriumTlumaczen) UstawTlumaczenieZwrotne(ctx context.Context,
	kodPanelu, tresc string) (PanelTlumaczenia, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, ustawTrescZwrotnaPanelu)
	if err != nil {
		return PanelTlumaczenia{}, err
	}
	wynik, err := polecenie.ExecContext(ctx, tresc, time.Now().UnixMilli(), kodPanelu, KontoOperatora(ctx))
	if err != nil {
		return PanelTlumaczenia{}, fmt.Errorf("dane: nie można ustawić tłumaczenia zwrotnego panelu %q: %w", kodPanelu, err)
	}
	if err := r.trafieniePanelu(ctx, wynik, kodPanelu); err != nil {
		return PanelTlumaczenia{}, err
	}
	return r.Panel(ctx, kodPanelu)
}

func (r *repozytoriumTlumaczen) UstawTon(ctx context.Context, kodPanelu, ton string) (PanelTlumaczenia, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, ustawTonPanelu)
	if err != nil {
		return PanelTlumaczenia{}, err
	}
	wynik, err := polecenie.ExecContext(ctx, ton, time.Now().UnixMilli(), kodPanelu, KontoOperatora(ctx))
	if err != nil {
		return PanelTlumaczenia{}, fmt.Errorf("dane: nie można ustawić tonu panelu %q: %w", kodPanelu, err)
	}
	if err := r.trafieniePanelu(ctx, wynik, kodPanelu); err != nil {
		return PanelTlumaczenia{}, err
	}
	return r.Panel(ctx, kodPanelu)
}

// Stan panelu: CHECK w migracja_053_tlumaczenie.sql, wartości kontraktu wprost.
func (r *repozytoriumTlumaczen) UstawStanPanelu(ctx context.Context, kodPanelu, stan string) (PanelTlumaczenia, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, ustawStanPanelu)
	if err != nil {
		return PanelTlumaczenia{}, err
	}
	wynik, err := polecenie.ExecContext(ctx, stan, time.Now().UnixMilli(), kodPanelu, KontoOperatora(ctx))
	if err != nil {
		return PanelTlumaczenia{}, fmt.Errorf("dane: nie można ustawić stanu panelu %q: %w", kodPanelu, err)
	}
	if err := r.trafieniePanelu(ctx, wynik, kodPanelu); err != nil {
		return PanelTlumaczenia{}, err
	}
	return r.Panel(ctx, kodPanelu)
}

// Silnik oddaje zero wierszy tak dla panelu cudzego, jak dla panelu, którego nie ma.
func (r *repozytoriumTlumaczen) trafieniePanelu(ctx context.Context, wynik sql.Result, kodPanelu string) error {
	zmienione, err := wynik.RowsAffected()
	if err != nil {
		return fmt.Errorf("dane: nieczytelny wynik zmiany panelu %q: %w", kodPanelu, err)
	}
	if zmienione > 0 {
		return nil
	}
	polecenie, err := r.zapytania.przygotuj(ctx, istniejePanelTlumaczenia)
	if err != nil {
		return err
	}
	var jeden int
	err = polecenie.QueryRowContext(ctx, kodPanelu).Scan(&jeden)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrBrakWiersza
	}
	if err != nil {
		return fmt.Errorf("dane: nie można sprawdzić panelu %q: %w", kodPanelu, err)
	}
	return fmt.Errorf("dane: panel tłumaczenia %q należy do innego konta: %w", kodPanelu, ErrKolizjaWiersza)
}
