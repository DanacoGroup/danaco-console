// Repozytorium przechowuje zależności między zadaniami projektu w tabeli
// `zaleznosc_zadan_projektu`, pilnując jednoznaczności każdej krawędzi.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"danacoconsole/shared"
)

// ZaleznoscWorkspace to wiersz tabeli `zaleznosc_zadan_projektu`, reprezentujący
// zależność między dwoma zadaniami.
type ZaleznoscWorkspace struct {
	ProjektID     int64
	Identyfikator string
	Poprzednik    string
	Nastepnik     string
	Rodzaj        shared.WorkspaceDependencyKind
	OdstepMinut   int
}

const (
	kolumnyZaleznosciWorkspace = `projekt_id, identyfikator_zewnetrzny, poprzednik, nastepnik,
	                              rodzaj, odstep_minut`

	zapiszZaleznoscWorkspace = `INSERT INTO zaleznosc_zadan_projektu
	    (identyfikator_zewnetrzny, projekt_id, poprzednik, nastepnik, rodzaj, odstep_minut)
	    VALUES (?, ?, ?, ?, ?, ?)
	    ON CONFLICT(projekt_id, poprzednik, nastepnik) DO UPDATE SET
	        rodzaj = excluded.rodzaj, odstep_minut = excluded.odstep_minut`

	pobierzZaleznoscWorkspace = `SELECT ` + kolumnyZaleznosciWorkspace +
		` FROM zaleznosc_zadan_projektu WHERE identyfikator_zewnetrzny = ?`

	pobierzZaleznoscParaWorkspace = `SELECT ` + kolumnyZaleznosciWorkspace +
		` FROM zaleznosc_zadan_projektu WHERE projekt_id = ? AND poprzednik = ? AND nastepnik = ?`

	listaZaleznosciWorkspace = `SELECT ` + kolumnyZaleznosciWorkspace +
		` FROM zaleznosc_zadan_projektu WHERE projekt_id = ? ORDER BY id`

	usunZaleznoscWorkspace = `DELETE FROM zaleznosc_zadan_projektu WHERE identyfikator_zewnetrzny = ?`
)

// ZapiszZaleznoscWorkspace zakłada zależność i oddaje jej stan po zapisie.
// Powtórne założenie tej samej krawędzi zmienia jej rodzaj i odstęp zamiast
// zakładać drugą.
func (r *repozytoriumPrzestrzeniRoboczej) ZapiszZaleznoscWorkspace(ctx context.Context,
	zaleznosc ZaleznoscWorkspace) (ZaleznoscWorkspace, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, zapiszZaleznoscWorkspace)
	if err != nil {
		return ZaleznoscWorkspace{}, err
	}
	_, err = polecenie.ExecContext(ctx, zaleznosc.Identyfikator, zaleznosc.ProjektID,
		zaleznosc.Poprzednik, zaleznosc.Nastepnik, string(zaleznosc.Rodzaj), zaleznosc.OdstepMinut)
	if err != nil {
		return ZaleznoscWorkspace{}, fmt.Errorf("dane: nie można zapisać zależności zadań %q→%q: %w",
			zaleznosc.Poprzednik, zaleznosc.Nastepnik, err)
	}
	polecenieOdczytu, err := r.zapytania.przygotuj(ctx, pobierzZaleznoscParaWorkspace)
	if err != nil {
		return ZaleznoscWorkspace{}, err
	}
	zapisana, err := odczytajZaleznoscWorkspace(polecenieOdczytu.QueryRowContext(ctx,
		zaleznosc.ProjektID, zaleznosc.Poprzednik, zaleznosc.Nastepnik))
	if err != nil {
		return ZaleznoscWorkspace{}, fmt.Errorf("dane: nieczytelna zależność po zapisie: %w", err)
	}
	return zapisana, nil
}

// ZaleznoscWorkspace zwraca z tabeli `zaleznosc_zadan_projektu` jedną
// zależność po jej identyfikatorze.
func (r *repozytoriumPrzestrzeniRoboczej) ZaleznoscWorkspace(ctx context.Context,
	identyfikator string) (ZaleznoscWorkspace, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzZaleznoscWorkspace)
	if err != nil {
		return ZaleznoscWorkspace{}, err
	}
	zaleznosc, err := odczytajZaleznoscWorkspace(polecenie.QueryRowContext(ctx, identyfikator))
	if errors.Is(err, sql.ErrNoRows) {
		return ZaleznoscWorkspace{}, ErrBrakWiersza
	}
	if err != nil {
		return ZaleznoscWorkspace{}, fmt.Errorf("dane: nieczytelna zależność %q: %w",
			identyfikator, err)
	}
	return zaleznosc, nil
}

// ZaleznosciWorkspace zwraca z tabeli `zaleznosc_zadan_projektu` komplet
// zależności całego projektu naraz.
func (r *repozytoriumPrzestrzeniRoboczej) ZaleznosciWorkspace(ctx context.Context,
	projektID int64) ([]ZaleznoscWorkspace, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaZaleznosciWorkspace)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, projektID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać zależności projektu %d: %w", projektID, err)
	}
	defer wiersze.Close()

	lista := []ZaleznoscWorkspace{}
	for wiersze.Next() {
		zaleznosc, err := odczytajZaleznoscWorkspace(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz zależności zadań: %w", err)
		}
		lista = append(lista, zaleznosc)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt zależności projektu %d: %w", projektID, err)
	}
	return lista, nil
}

// UsunZaleznoscWorkspace znosi zależność. Fałsz znaczy, że takiej krawędzi nie
// było — to nie jest błąd wywołania.
func (r *repozytoriumPrzestrzeniRoboczej) UsunZaleznoscWorkspace(ctx context.Context,
	identyfikator string) (bool, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, usunZaleznoscWorkspace)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, identyfikator)
	if err != nil {
		return false, fmt.Errorf("dane: nie można znieść zależności %q: %w", identyfikator, err)
	}
	zmienione, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nieznany wynik zniesienia zależności %q: %w",
			identyfikator, err)
	}
	return zmienione > 0, nil
}

// odczytajZaleznoscWorkspace składa strukturę ZaleznoscWorkspace z jednego
// wiersza wyniku zapytania SQL.
func odczytajZaleznoscWorkspace(wiersz skaner) (ZaleznoscWorkspace, error) {
	var zaleznosc ZaleznoscWorkspace
	var rodzaj string
	err := wiersz.Scan(&zaleznosc.ProjektID, &zaleznosc.Identyfikator, &zaleznosc.Poprzednik,
		&zaleznosc.Nastepnik, &rodzaj, &zaleznosc.OdstepMinut)
	if err != nil {
		return ZaleznoscWorkspace{}, err
	}
	zaleznosc.Rodzaj = shared.WorkspaceDependencyKind(rodzaj)
	return zaleznosc, nil
}
