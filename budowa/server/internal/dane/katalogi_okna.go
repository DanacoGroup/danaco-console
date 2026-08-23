// Odpowiedzialność pliku: lista katalogów roboczych okna (tabela `katalog_okna`).
// Katalog roboczy jest listą, nie pojedynczym polem, więc zapis okna
// wymienia całą listę w tej samej transakcji, w której zapisuje samo okno.
package dane

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

const (
	usunKatalogiOkna = `DELETE FROM katalog_okna WHERE okno_komunikacji_id = ?`

	wstawKatalogOkna = `INSERT INTO katalog_okna (okno_komunikacji_id, sciezka, kolejnosc)
	                    VALUES (?, ?, ?)`

	listaKatalogowOkna = `SELECT sciezka FROM katalog_okna
	                      WHERE okno_komunikacji_id = ? ORDER BY kolejnosc, id`

	listaKatalogowSesji = `SELECT k.okno_komunikacji_id, k.sciezka
	                       FROM katalog_okna k
	                       JOIN okno_komunikacji o ON o.id = k.okno_komunikacji_id
	                       WHERE o.sesja_id = ?
	                       ORDER BY k.okno_komunikacji_id, k.kolejnosc, k.id`
)

// zapiszKatalogiOkna wymienia listę katalogów okna. Wywoływane wyłącznie
// wewnątrz transakcji zapisu okna — inaczej okno i jego katalogi mogłyby się
// rozjechać.
func zapiszKatalogiOkna(ctx context.Context, z *zapytania, transakcja *sql.Tx,
	oknoID int64, katalogi []string) error {

	kasowanie, err := z.wTransakcji(ctx, transakcja, usunKatalogiOkna)
	if err != nil {
		return err
	}
	if _, err := kasowanie.ExecContext(ctx, oknoID); err != nil {
		return fmt.Errorf("dane: nie można wyczyścić katalogów okna %d: %w", oknoID, err)
	}
	wstawianie, err := z.wTransakcji(ctx, transakcja, wstawKatalogOkna)
	if err != nil {
		return err
	}
	kolejnosc := 0
	for _, sciezka := range katalogi {
		sciezka = strings.TrimSpace(sciezka)
		if sciezka == "" {
			continue
		}
		if _, err := wstawianie.ExecContext(ctx, oknoID, sciezka, kolejnosc); err != nil {
			return fmt.Errorf("dane: nie można zapisać katalogu %q okna %d: %w", sciezka, oknoID, err)
		}
		kolejnosc++
	}
	return nil
}

// katalogiOkna zwraca listę katalogów jednego okna w zapisanej kolejności.
func katalogiOkna(ctx context.Context, z *zapytania, oknoID int64) ([]string, error) {
	polecenie, err := z.przygotuj(ctx, listaKatalogowOkna)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, oknoID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać katalogów okna %d: %w", oknoID, err)
	}
	defer wiersze.Close()

	katalogi := []string{}
	for wiersze.Next() {
		var sciezka string
		if err := wiersze.Scan(&sciezka); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny katalog okna %d: %w", oknoID, err)
		}
		katalogi = append(katalogi, sciezka)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt katalogów okna %d: %w", oknoID, err)
	}
	return katalogi, nil
}

// katalogiOkienSesji zwraca katalogi wszystkich okien sesji jednym zapytaniem —
// lista okien sesji nie mnoży zapytań przez liczbę okien.
func katalogiOkienSesji(ctx context.Context, z *zapytania, sesjaID int64) (map[int64][]string, error) {
	polecenie, err := z.przygotuj(ctx, listaKatalogowSesji)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, sesjaID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać katalogów okien sesji %d: %w", sesjaID, err)
	}
	defer wiersze.Close()

	wynik := map[int64][]string{}
	for wiersze.Next() {
		var oknoID int64
		var sciezka string
		if err := wiersze.Scan(&oknoID, &sciezka); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny katalog okna sesji %d: %w", sesjaID, err)
		}
		wynik[oknoID] = append(wynik[oknoID], sciezka)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt katalogów okien sesji %d: %w", sesjaID, err)
	}
	return wynik, nil
}
