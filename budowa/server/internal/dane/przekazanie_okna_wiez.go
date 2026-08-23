// Odpowiedzialność pliku: więź koordynator–wykonawca w obszarze window.*. Więź nie
// ma własnej tabeli — mieszka w kolumnie `okno_komunikacji.okno_koordynatora_id`,
// do której kontrakt odwołuje się wprost jako `Window.coordinatorWindowId`. Ten
// plik dokłada do repozytorium `RepozytoriumPrzekazan` (deklarowanego przez
// `przekazanie_okna.go`) drogę zapisu i odczytu tej kolumny z poziomu
// identyfikatora zewnętrznego okna: `dane/okna.go` czyta i pisze tę kolumnę
// wyłącznie jako część pełnego wiersza okna (`Pobierz`/`Aktualizuj`, po
// identyfikatorze wewnętrznym `int64`), a Mission Control operuje na
// identyfikatorach zewnętrznych pojedynczej więzi, nie całego okna.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

const (
	// idOknaPoIdentyfikatorze odnajduje wewnętrzny identyfikator okna po jego
	// identyfikatorze zewnętrznym — obie metody zapisu i odczytu więzi
	// przyjmują identyfikatory zewnętrzne, tak jak kontrakt Mission Control.
	idOknaPoIdentyfikatorze = `SELECT id FROM okno_komunikacji WHERE identyfikator_zewnetrzny = ?`

	ustawKoordynatoraOkna = `UPDATE okno_komunikacji
	                         SET okno_koordynatora_id = ?,
	                             zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                         WHERE id = ?`

	koordynatorOkna = `SELECT koordynator.identyfikator_zewnetrzny
	                   FROM okno_komunikacji wykonawca
	                   LEFT JOIN okno_komunikacji koordynator
	                       ON koordynator.id = wykonawca.okno_koordynatora_id
	                   WHERE wykonawca.identyfikator_zewnetrzny = ?`

	wykonawcyKoordynatora = `SELECT wykonawca.identyfikator_zewnetrzny
	                         FROM okno_komunikacji wykonawca
	                         JOIN okno_komunikacji koordynator
	                             ON koordynator.id = wykonawca.okno_koordynatora_id
	                         WHERE koordynator.identyfikator_zewnetrzny = ?
	                         ORDER BY wykonawca.kolejnosc, wykonawca.id`
)

// idOknaZewnetrzne odnajduje wewnętrzny identyfikator okna po zewnętrznym.
// Brak wiersza wraca jako ErrBrakWiersza — wspólne dla wszystkich trzech metod
// tego pliku, bo każda z nich operuje na oknie wskazanym z zewnątrz.
func (r *repozytoriumPrzekazan) idOknaZewnetrzne(ctx context.Context, identyfikator string) (int64, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, idOknaPoIdentyfikatorze)
	if err != nil {
		return 0, err
	}
	var id int64
	err = polecenie.QueryRowContext(ctx, identyfikator).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("%w: okno %q", ErrBrakWiersza, identyfikator)
	}
	if err != nil {
		return 0, fmt.Errorf("dane: nie można odnaleźć okna %q: %w", identyfikator, err)
	}
	return id, nil
}

// UstawKoordynatora zapisuje więź koordynator–wykonawca: oknoWykonawcy zaczyna
// wskazywać oknoKoordynatora jako swojego koordynatora.
//
// Oba okna muszą istnieć — cicha zgoda na więź z oknem, którego nie ma, dałaby
// potwierdzenie relacji, która w rzeczywistości nie powstała. Obie strony więzi
// rozwiązujemy na identyfikatory wewnętrzne przed zapisem, a nie podzapytaniem
// w UPDATE, które ciche niedopasowanie zamieniłoby w NULL.
func (r *repozytoriumPrzekazan) UstawKoordynatora(ctx context.Context, oknoWykonawcy, oknoKoordynatora string) error {
	idWykonawcy, err := r.idOknaZewnetrzne(ctx, oknoWykonawcy)
	if err != nil {
		return err
	}
	idKoordynatora, err := r.idOknaZewnetrzne(ctx, oknoKoordynatora)
	if err != nil {
		return err
	}
	polecenie, err := r.zapytania.przygotuj(ctx, ustawKoordynatoraOkna)
	if err != nil {
		return err
	}
	wynik, err := polecenie.ExecContext(ctx, idKoordynatora, idWykonawcy)
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać koordynatora okna %q: %w", oknoWykonawcy, err)
	}
	return sprawdzTrafienie(wynik, "okno_komunikacji", idWykonawcy)
}

// Koordynator zwraca identyfikator zewnętrzny koordynatora okna wykonawcy.
// Okno samodzielne, bez koordynatora, oddaje pusty napis — to poprawny stan, nie
// usterka. Brak samego okna wykonawcy wraca jako ErrBrakWiersza.
func (r *repozytoriumPrzekazan) Koordynator(ctx context.Context, oknoWykonawcy string) (string, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, koordynatorOkna)
	if err != nil {
		return "", err
	}
	var koordynator sql.NullString
	err = polecenie.QueryRowContext(ctx, oknoWykonawcy).Scan(&koordynator)
	if errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("%w: okno %q", ErrBrakWiersza, oknoWykonawcy)
	}
	if err != nil {
		return "", fmt.Errorf("dane: nie można odczytać koordynatora okna %q: %w", oknoWykonawcy, err)
	}
	return koordynator.String, nil
}

// Wykonawcy zwraca identyfikatory zewnętrzne okien podległych koordynatorowi —
// zapytanie widoku Mission Control. Kolejność ustala `kolejnosc, id` okna
// wykonawcy, tak jak każda inna lista okien w tym module (`dane/okna.go`),
// żeby widok nie tasował wierszy między odświeżeniami.
func (r *repozytoriumPrzekazan) Wykonawcy(ctx context.Context, oknoKoordynatora string) ([]string, error) {
	if _, err := r.idOknaZewnetrzne(ctx, oknoKoordynatora); err != nil {
		return nil, err
	}
	polecenie, err := r.zapytania.przygotuj(ctx, wykonawcyKoordynatora)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, oknoKoordynatora)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać wykonawców okna %q: %w", oknoKoordynatora, err)
	}
	defer wiersze.Close()

	lista := []string{}
	for wiersze.Next() {
		var identyfikator string
		if err := wiersze.Scan(&identyfikator); err != nil {
			return nil, fmt.Errorf("dane: przerwany odczyt wykonawców okna %q: %w", oknoKoordynatora, err)
		}
		lista = append(lista, identyfikator)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt wykonawców okna %q: %w", oknoKoordynatora, err)
	}
	return lista, nil
}
