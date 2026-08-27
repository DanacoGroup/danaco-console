// Plik odczytuje harmonogramy automatyk po ich identyfikatorze zewnętrznym
// oraz odczytuje pełny wykaz harmonogramów z tabeli harmonogram_automatyki.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

const (
	pobierzHarmonogramPoKodzie = `SELECT ` + kolumnyHarmonogramu + ` FROM harmonogram_automatyki
	                              WHERE identyfikator_zewnetrzny = ?`

	// Warunek ograniczający wynik do harmonogramów czynnych działa w zapytaniu
	// SQL, a wynik jest uporządkowany według terminu najbliższego uruchomienia.
	listaHarmonogramow = `SELECT ` + kolumnyHarmonogramu + ` FROM harmonogram_automatyki
	                      WHERE (? = 0 OR czynny = 1)
	                      ORDER BY nastepne_uruchomienie IS NULL, nastepne_uruchomienie, id`
)

// HarmonogramPoKodzie zwraca harmonogram o wskazanym identyfikatorze zewnętrznym.
// Brak wiersza wraca jako ErrBrakWiersza — warstwa wyższa odróżnia „harmonogramu
// nie ma" od „odczyt się nie powiódł" i nazywa byt, którego brakuje.
func (r *repozytoriumAutomatyk) HarmonogramPoKodzie(ctx context.Context, kod string) (Harmonogram, error) {
	if kod == "" {
		return Harmonogram{}, fmt.Errorf("dane: odczyt harmonogramu bez identyfikatora")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzHarmonogramPoKodzie)
	if err != nil {
		return Harmonogram{}, err
	}
	harmonogram, err := odczytajHarmonogram(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return Harmonogram{}, ErrBrakWiersza
	}
	if err != nil {
		return Harmonogram{}, fmt.Errorf("dane: nieczytelny wiersz harmonogramu %q: %w", kod, err)
	}
	return harmonogram, nil
}

// Harmonogramy zwraca komplet harmonogramów, opcjonalnie zawężony do czynnych.
// Pusty wykaz jest stanem poprawnym: platforma bez ani jednej automatyki
// z harmonogramem nie ma czego oddać.
func (r *repozytoriumAutomatyk) Harmonogramy(ctx context.Context, tylkoCzynne bool) ([]Harmonogram, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaHarmonogramow)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, liczbaLogiczna(tylkoCzynne))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać harmonogramów: %w", err)
	}
	defer wiersze.Close()

	lista := []Harmonogram{}
	for wiersze.Next() {
		harmonogram, err := odczytajHarmonogram(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz harmonogramu: %w", err)
		}
		lista = append(lista, harmonogram)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt harmonogramów: %w", err)
	}
	return lista, nil
}
