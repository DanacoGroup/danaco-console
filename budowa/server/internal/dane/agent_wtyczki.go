// Plik prowadzi wtyczki eksperta — katalog rozszerzeń powłoki, osobny od konektorów wskazujących usługi przez most
// punktów dostępu; warstwy promptu i tożsamość własna eksperta leżą w agent_warstwy.go jako ta sama implementacja repozytorium.
package dane

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"
)

// przedrostekWtyczki znakuje kod wtyczki eksperta. Konektor ma „ak-”
// (`core/adapter_modul_agents_przeklad.go`); wtyczka jest bytem odrębnym, więc
// dostaje własny przedrostek, po którym widać ją w dzienniku bez zaglądania
// do tabeli.
const przedrostekWtyczki = "aw-"

// licznikWtyczek rozróżnia kody wytworzone w tej samej milisekundzie przez
// równoległe wywołania — wzorzec `core/identyfikator.go`.
var licznikWtyczek atomic.Uint64

const (
	wstawWtyczkeAgenta = `INSERT INTO agent_wtyczka (kod, agent_id, nazwa, zrodlo, wersja, aktywna)
	                      VALUES (?, ?, ?, ?, ?, 1)`

	usunWtyczkeAgenta = `DELETE FROM agent_wtyczka WHERE kod = ? AND agent_id = ?`

	wtyczkiAgenta = `SELECT w.id, w.kod, a.kod, w.nazwa, w.zrodlo, w.wersja, w.aktywna, w.utworzono
	                   FROM agent_wtyczka w JOIN agent a ON a.id = w.agent_id
	                  WHERE w.agent_id = ?
	                  ORDER BY w.nazwa, w.kod`
)

// Wtyczka nie ma wskazania konta; granicę niesie korzeń `agent` po stronie złączenia.
var (
	wtyczkaPoKodzie = `SELECT w.id, w.kod, a.kod, w.nazwa, w.zrodlo, w.wersja, w.aktywna, w.utworzono
	                     FROM agent_wtyczka w JOIN agent a ON a.id = w.agent_id
	                    WHERE w.kod = ? AND ` + warunekKontaEksperta

	wtyczkiWszystkich = `SELECT w.id, w.kod, a.kod, w.nazwa, w.zrodlo, w.wersja, w.aktywna, w.utworzono
	                       FROM agent_wtyczka w JOIN agent a ON a.id = w.agent_id
	                      WHERE ` + warunekKontaEksperta + `
	                      ORDER BY a.kod, w.nazwa, w.kod`
)

// DodajWtyczke przypisuje ekspertowi wtyczkę i oddaje zapisany wiersz. Wstawienie
// i odczyt zwracanego wiersza idą w jednej transakcji, żeby oddany wpis był
// dokładnie tym, który właśnie powstał, a nie stanem zastanym chwilę później.
func (r *repozytoriumWarstwAgenta) DodajWtyczke(ctx context.Context, kodAgenta, nazwa string,
	zrodlo, wersja *string) (WtyczkaAgenta, error) {

	numer, err := r.numerAgenta(ctx, kodAgenta)
	if err != nil {
		return WtyczkaAgenta{}, err
	}
	nazwaWtyczki := strings.TrimSpace(nazwa)
	if nazwaWtyczki == "" {
		return WtyczkaAgenta{}, fmt.Errorf("dane: wtyczka eksperta %q wymaga nazwy", kodAgenta)
	}
	kod := nowyKodWtyczki()
	var zapisana WtyczkaAgenta
	err = wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		wstaw, err := r.zapytania.wTransakcji(ctx, transakcja, wstawWtyczkeAgenta)
		if err != nil {
			return err
		}
		_, err = wstaw.ExecContext(ctx, kod, numer, nazwaWtyczki,
			tekstDoKolumny(zrodlo), tekstDoKolumny(wersja))
		if err != nil {
			return fmt.Errorf("dane: nie można przypisać wtyczki %q ekspertowi %q: %w",
				nazwaWtyczki, kodAgenta, err)
		}
		odczyt, err := r.zapytania.wTransakcji(ctx, transakcja, wtyczkaPoKodzie)
		if err != nil {
			return err
		}
		zapisana, err = odczytajWtyczke(odczyt.QueryRowContext(ctx, kod, KontoOperatora(ctx)))
		if err != nil {
			return fmt.Errorf("dane: nie można odczytać wtyczki %q: %w", kod, err)
		}
		return nil
	})
	if err != nil {
		return WtyczkaAgenta{}, err
	}
	return zapisana, nil
}

// UsunWtyczke odłącza wtyczkę od eksperta i zwraca informację, czy wiersz istniał; warunek na numer eksperta
// pilnuje, żeby kod wtyczki cudzej nie skasował wiersza innego eksperta.
func (r *repozytoriumWarstwAgenta) UsunWtyczke(ctx context.Context,
	kodAgenta, kodWtyczki string) (bool, error) {

	numer, err := r.numerAgenta(ctx, kodAgenta)
	if err != nil {
		return false, err
	}
	polecenie, err := r.zapytania.przygotuj(ctx, usunWtyczkeAgenta)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, strings.TrimSpace(kodWtyczki), numer)
	if err != nil {
		return false, fmt.Errorf("dane: nie można odłączyć wtyczki %q eksperta %q: %w",
			kodWtyczki, kodAgenta, err)
	}
	zdjete, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nie można ustalić skutku odłączenia wtyczki %q: %w", kodWtyczki, err)
	}
	return zdjete > 0, nil
}

// Wtyczki zwraca wtyczki eksperta w porządku nazw — tak samo jak wykaz
// konektorów. Ekspert bez wtyczek oddaje wykaz pusty, nie błąd.
func (r *repozytoriumWarstwAgenta) Wtyczki(ctx context.Context, kodAgenta string) ([]WtyczkaAgenta, error) {
	numer, err := r.numerAgenta(ctx, kodAgenta)
	if err != nil {
		return nil, err
	}
	polecenie, err := r.zapytania.przygotuj(ctx, wtyczkiAgenta)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, numer)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać wtyczek eksperta %q: %w", kodAgenta, err)
	}
	defer wiersze.Close()

	zebrane := []WtyczkaAgenta{}
	for wiersze.Next() {
		wpis, err := odczytajWtyczke(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: przerwany odczyt wtyczek eksperta %q: %w", kodAgenta, err)
		}
		zebrane = append(zebrane, wpis)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt wtyczek eksperta %q: %w", kodAgenta, err)
	}
	return zebrane, nil
}

// odczytajWtyczke składa strukturę wtyczki eksperta wprost z jednego wiersza wyniku zapytania do bazy.
func odczytajWtyczke(wiersz skaner) (WtyczkaAgenta, error) {
	var wpis WtyczkaAgenta
	var zrodlo, wersja sql.NullString
	var aktywna int
	err := wiersz.Scan(&wpis.ID, &wpis.Kod, &wpis.AgentKod, &wpis.Nazwa,
		&zrodlo, &wersja, &aktywna, &wpis.Utworzono)
	if err != nil {
		return WtyczkaAgenta{}, err
	}
	wpis.Zrodlo = tekstZKolumny(zrodlo)
	wpis.Wersja = tekstZKolumny(wersja)
	wpis.Aktywna = aktywna == 1
	return wpis, nil
}

// nowyKodWtyczki składa kod trwały wtyczki z licznika i części losowej. Gdy
// źródło losowe zawiedzie, wystarcza część licznikowa — brak losowości nie ma
// prawa wstrzymać zapisu. Kształt jest ten sam co w
// `core/identyfikator.go:nowyIdentyfikator`.
func nowyKodWtyczki() string {
	kolejny := licznikWtyczek.Add(1)
	losowe := make([]byte, 8)
	if _, err := rand.Read(losowe); err != nil {
		return przedrostekWtyczki + strconv.FormatUint(kolejny, 36)
	}
	return przedrostekWtyczki + strconv.FormatUint(kolejny, 36) + "-" + hex.EncodeToString(losowe)
}

// WtyczkiWszystkich oddaje wtyczki wszystkich ekspertów jednym zapytaniem, tą samą drogą co WarstwyWszystkich; ekspert bez wtyczek nie dostaje wpisu.
func (r *repozytoriumWarstwAgenta) WtyczkiWszystkich(
	ctx context.Context,
) (map[string][]WtyczkaAgenta, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, wtyczkiWszystkich)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, KontoOperatora(ctx))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać wtyczek ekspertów: %w", err)
	}
	defer wiersze.Close()

	zebrane := map[string][]WtyczkaAgenta{}
	for wiersze.Next() {
		wpis, err := odczytajWtyczke(wiersze)
		if err != nil {
			return nil, err
		}
		zebrane[wpis.AgentKod] = append(zebrane[wpis.AgentKod], wpis)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt wtyczek ekspertów: %w", err)
	}
	return zebrane, nil
}
