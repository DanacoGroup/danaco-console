// Plik prowadzi warstwy promptu eksperta oraz jego tożsamość własną: imię widoczne i favikon; wtyczki eksperta leżą
// w agent_wtyczki.go jako ta sama implementacja repozytorium, rozdzielona wyłącznie na dwa pliki.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// WarstwaAgenta to wiersz `agent_warstwa`. `Warstwa` i `Tryb` niosą wartości
// kontraktu wprost — 'constitution' | 'profile' | 'expertise' oraz
// 'ZASTAP' | 'DOLACZ'.
type WarstwaAgenta struct {
	Warstwa        string
	Tresc          string
	Tryb           string
	Aktywna        bool
	Zaktualizowano string
}

// WtyczkaAgenta to wiersz `agent_wtyczka`. `Kod` jest identyfikatorem trwałym
// i odpowiada polu `AgentPlugin.id` kontraktu; `AgentKod` niesie kod eksperta,
// bo warstwy wyższe wskazują eksperta kodem, nie numerem wiersza.
type WtyczkaAgenta struct {
	ID        int64
	Kod       string
	AgentKod  string
	Nazwa     string
	Zrodlo    *string
	Wersja    *string
	Aktywna   bool
	Utworzono string
}

// RepozytoriumWarstwAgenta jest kontraktem tożsamości własnej eksperta:
// warstw jego promptu i jego wtyczek.
type RepozytoriumWarstwAgenta interface {
	UstawWarstwe(ctx context.Context, kodAgenta, warstwa, tresc, tryb string, aktywna bool) error
	UsunWarstwe(ctx context.Context, kodAgenta, warstwa string) (bool, error)
	Warstwy(ctx context.Context, kodAgenta string) ([]WarstwaAgenta, error)
	// WarstwyWszystkich oddaje warstwy wszystkich ekspertów jednym zapytaniem, po kodzie eksperta.
	WarstwyWszystkich(ctx context.Context) (map[string][]WarstwaAgenta, error)
	DodajWtyczke(ctx context.Context, kodAgenta, nazwa string, zrodlo, wersja *string) (WtyczkaAgenta, error)
	UsunWtyczke(ctx context.Context, kodAgenta, kodWtyczki string) (bool, error)
	Wtyczki(ctx context.Context, kodAgenta string) ([]WtyczkaAgenta, error)
	// WtyczkiWszystkich oddaje wtyczki wszystkich ekspertów — z tego samego powodu.
	WtyczkiWszystkich(ctx context.Context) (map[string][]WtyczkaAgenta, error)
	UstawTozsamosc(ctx context.Context, kodAgenta, imieWlasne, favikon string) error
	// UstawTrybNakladki zapisuje tryb nałożenia instrukcji; wartość spoza katalogu odrzuca warunek CHECK.
	UstawTrybNakladki(ctx context.Context, kodAgenta, tryb string) error
}

const (
	numerAgentaPoKodzie = `SELECT id FROM agent WHERE kod = ?`

	zapiszWarstweAgenta = `INSERT INTO agent_warstwa (agent_id, warstwa, tresc, tryb, aktywna)
	                       VALUES (?, ?, ?, ?, ?)
	                       ON CONFLICT(agent_id, warstwa) DO UPDATE SET
	                           tresc          = excluded.tresc,
	                           tryb           = excluded.tryb,
	                           aktywna        = excluded.aktywna,
	                           zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ', 'now')`

	usunWarstweAgenta = `DELETE FROM agent_warstwa WHERE agent_id = ? AND warstwa = ?`

	// Porządek warstw idzie wg krytyczności, tak samo jak w warstwie nakładki: konstytucja stoi najwyżej, ekspertyza zadaniowa najniżej.
	warstwyWszystkich = `SELECT a.kod, w.warstwa, w.tresc, w.tryb, w.aktywna, w.zaktualizowano
	                       FROM agent_warstwa w JOIN agent a ON a.id = w.agent_id
	                      ORDER BY a.kod, CASE w.warstwa
	                                          WHEN 'constitution' THEN 1
	                                          WHEN 'profile'      THEN 2
	                                          WHEN 'expertise'    THEN 3
	                                          ELSE 4
	                                      END, w.warstwa`

	warstwyAgenta = `SELECT warstwa, tresc, tryb, aktywna, zaktualizowano
	                   FROM agent_warstwa WHERE agent_id = ?
	                  ORDER BY CASE warstwa
	                               WHEN 'constitution' THEN 1
	                               WHEN 'profile'      THEN 2
	                               WHEN 'expertise'    THEN 3
	                               ELSE 4
	                           END, warstwa`

	zapiszTozsamoscAgenta = `UPDATE agent
	                            SET imie_wlasne = ?, favikon = ?,
	                                zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ', 'now')
	                          WHERE kod = ?`

	zapiszTrybNakladki = `UPDATE agent
	                         SET tryb_nakladki = ?,
	                             zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ', 'now')
	                       WHERE kod = ?`
)

type repozytoriumWarstwAgenta struct {
	zapytania *zapytania
	db        *sql.DB
}

var _ RepozytoriumWarstwAgenta = (*repozytoriumWarstwAgenta)(nil)

func noweRepozytoriumWarstwAgenta(z *zapytania, db *sql.DB) *repozytoriumWarstwAgenta {
	return &repozytoriumWarstwAgenta{zapytania: z, db: db}
}

// UstawWarstwe zapisuje treść jednej warstwy promptu eksperta. Zapis powtórzony
// nadpisuje warstwę zamiast dublować wiersz. Treść pusta nie jest błędem —
// znaczy „warstwa istnieje, ale nic nie wnosi", a to inny stan niż jej brak.
func (r *repozytoriumWarstwAgenta) UstawWarstwe(ctx context.Context,
	kodAgenta, warstwa, tresc, tryb string, aktywna bool) error {

	numer, err := r.numerAgenta(ctx, kodAgenta)
	if err != nil {
		return err
	}
	nazwaWarstwy := strings.TrimSpace(warstwa)
	if nazwaWarstwy == "" {
		return fmt.Errorf("dane: warstwa eksperta %q wymaga nazwy", kodAgenta)
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszWarstweAgenta)
	if err != nil {
		return err
	}
	_, err = polecenie.ExecContext(ctx, numer, nazwaWarstwy, tresc,
		strings.TrimSpace(tryb), liczbaLogiczna(aktywna))
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać warstwy %q eksperta %q: %w",
			nazwaWarstwy, kodAgenta, err)
	}
	return nil
}

// UsunWarstwe zdejmuje warstwę z promptu eksperta. Zwraca informację, czy wiersz
// istniał — brak warstwy znaczy „prompt bez niej", nie błąd.
func (r *repozytoriumWarstwAgenta) UsunWarstwe(ctx context.Context,
	kodAgenta, warstwa string) (bool, error) {

	numer, err := r.numerAgenta(ctx, kodAgenta)
	if err != nil {
		return false, err
	}
	polecenie, err := r.zapytania.przygotuj(ctx, usunWarstweAgenta)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, numer, strings.TrimSpace(warstwa))
	if err != nil {
		return false, fmt.Errorf("dane: nie można usunąć warstwy %q eksperta %q: %w",
			warstwa, kodAgenta, err)
	}
	zdjete, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nie można ustalić skutku usunięcia warstwy %q: %w", warstwa, err)
	}
	return zdjete > 0, nil
}

// Warstwy zwraca warstwy eksperta w porządku krytyczności. Ekspert bez ani
// jednej warstwy oddaje wykaz pusty, nie błąd.
func (r *repozytoriumWarstwAgenta) Warstwy(ctx context.Context, kodAgenta string) ([]WarstwaAgenta, error) {
	numer, err := r.numerAgenta(ctx, kodAgenta)
	if err != nil {
		return nil, err
	}
	polecenie, err := r.zapytania.przygotuj(ctx, warstwyAgenta)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, numer)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać warstw eksperta %q: %w", kodAgenta, err)
	}
	defer wiersze.Close()

	zebrane := []WarstwaAgenta{}
	for wiersze.Next() {
		var wpis WarstwaAgenta
		var aktywna int
		if err := wiersze.Scan(&wpis.Warstwa, &wpis.Tresc, &wpis.Tryb, &aktywna, &wpis.Zaktualizowano); err != nil {
			return nil, fmt.Errorf("dane: przerwany odczyt warstw eksperta %q: %w", kodAgenta, err)
		}
		wpis.Aktywna = aktywna == 1
		zebrane = append(zebrane, wpis)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt warstw eksperta %q: %w", kodAgenta, err)
	}
	return zebrane, nil
}

// UstawTozsamosc zapisuje imię własne i favikon eksperta; obie wartości są napisami wolnymi, a pusta znaczy „nie nadano”.
func (r *repozytoriumWarstwAgenta) UstawTozsamosc(ctx context.Context,
	kodAgenta, imieWlasne, favikon string) error {

	polecenie, err := r.zapytania.przygotuj(ctx, zapiszTozsamoscAgenta)
	if err != nil {
		return err
	}
	wynik, err := polecenie.ExecContext(ctx, strings.TrimSpace(imieWlasne),
		strings.TrimSpace(favikon), kodAgenta)
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać tożsamości eksperta %q: %w", kodAgenta, err)
	}
	zmienione, err := wynik.RowsAffected()
	if err != nil {
		return fmt.Errorf("dane: nie można ustalić skutku zapisu tożsamości eksperta %q: %w", kodAgenta, err)
	}
	if zmienione == 0 {
		return fmt.Errorf("dane: ekspert %q nie istnieje: %w", kodAgenta, ErrBrakWiersza)
	}
	return nil
}

// numerAgenta przekłada kod trwały eksperta na numer wiersza. Ekspert
// nieistniejący wraca jako ErrBrakWiersza, żeby warstwa wyższa odróżniła
// „nie ma takiego eksperta" od „odczyt padł".
func (r *repozytoriumWarstwAgenta) numerAgenta(ctx context.Context, kodAgenta string) (int64, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, numerAgentaPoKodzie)
	if err != nil {
		return 0, err
	}
	var numer int64
	err = polecenie.QueryRowContext(ctx, kodAgenta).Scan(&numer)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("dane: ekspert %q nie istnieje: %w", kodAgenta, ErrBrakWiersza)
	}
	if err != nil {
		return 0, fmt.Errorf("dane: nie można odczytać eksperta %q: %w", kodAgenta, err)
	}
	return numer, nil
}

// WarstwyWszystkich oddaje warstwy wszystkich ekspertów jednym zapytaniem na całe wywołanie; ekspert bez ani jednej warstwy nie dostaje wpisu w mapie.
func (r *repozytoriumWarstwAgenta) WarstwyWszystkich(
	ctx context.Context,
) (map[string][]WarstwaAgenta, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, warstwyWszystkich)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać warstw ekspertów: %w", err)
	}
	defer wiersze.Close()

	zebrane := map[string][]WarstwaAgenta{}
	for wiersze.Next() {
		var kod string
		var wpis WarstwaAgenta
		var aktywna int
		if err := wiersze.Scan(
			&kod, &wpis.Warstwa, &wpis.Tresc, &wpis.Tryb, &aktywna, &wpis.Zaktualizowano,
		); err != nil {
			return nil, fmt.Errorf("dane: przerwany odczyt warstw ekspertów: %w", err)
		}
		wpis.Aktywna = aktywna == 1
		zebrane[kod] = append(zebrane[kod], wpis)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt warstw ekspertów: %w", err)
	}
	return zebrane, nil
}

// UstawTrybNakladki zapisuje tryb nałożenia instrukcji eksperta jako pole eksperta, nie warstwy; wartość spoza katalogu odrzuca warunek CHECK bazy.
func (r *repozytoriumWarstwAgenta) UstawTrybNakladki(ctx context.Context,
	kodAgenta, tryb string) error {

	polecenie, err := r.zapytania.przygotuj(ctx, zapiszTrybNakladki)
	if err != nil {
		return err
	}
	wynik, err := polecenie.ExecContext(ctx, strings.TrimSpace(tryb), kodAgenta)
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać trybu nakładki eksperta %q: %w", kodAgenta, err)
	}
	zmienione, err := wynik.RowsAffected()
	if err != nil {
		return fmt.Errorf("dane: nie można ustalić skutku zapisu trybu eksperta %q: %w", kodAgenta, err)
	}
	if zmienione == 0 {
		return fmt.Errorf("dane: ekspert %q nie istnieje: %w", kodAgenta, ErrBrakWiersza)
	}
	return nil
}
