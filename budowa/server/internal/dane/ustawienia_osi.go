// Plik odczytuje i kasuje ustawienia pod adresem złożonym: poziom zasięgu
// razem z osią rozstrzygania, dla platformy, modelu albo konta.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"danacoconsole/shared"
)

// RepozytoriumKonfiguracjiOsi rozszerza obszar konfiguracji o oś rozstrzygania,
// nie zastępując wersji bez osi — każda implementacja obsługuje obie postacie
// adresu.
type RepozytoriumKonfiguracjiOsi interface {
	RepozytoriumKonfiguracji
	OdczytajOsi(ctx context.Context, poziom shared.ConfigScope, kluczZasiegu string,
		os shared.ConfigAxis, kluczOsi, klucz string) (Ustawienie, bool, error)
	ListaOsi(ctx context.Context, poziom shared.ConfigScope, kluczZasiegu string,
		os shared.ConfigAxis, kluczOsi string) ([]Ustawienie, error)
	UsunOsi(ctx context.Context, poziom shared.ConfigScope, kluczZasiegu string,
		os shared.ConfigAxis, kluczOsi, klucz string) error
}

// OdczytajOsi zwraca ustawienie spod adresu złożonego. Drugi wynik mówi, czy
// wartość ustawiono — brak wiersza nie jest błędem.
func (r *repozytoriumKonfiguracji) OdczytajOsi(ctx context.Context, poziom shared.ConfigScope,
	kluczZasiegu string, os shared.ConfigAxis, kluczOsi, klucz string) (Ustawienie, bool, error) {

	kod, kodOsi, err := adresZlozony(poziom, os)
	if err != nil {
		return Ustawienie{}, false, err
	}
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzUstawienie)
	if err != nil {
		return Ustawienie{}, false, err
	}
	wiersz := polecenie.QueryRowContext(ctx, kod, kluczZasiegu, kodOsi, kluczOsi, klucz,
		KontoOperatora(ctx))
	ustawienie, err := odczytajUstawienie(wiersz)
	if errors.Is(err, sql.ErrNoRows) {
		return Ustawienie{}, false, nil
	}
	if err != nil {
		return Ustawienie{}, false, err
	}
	return ustawienie, true, nil
}

// ListaOsi zwraca ustawienia jednego bytu poziomu na jednej osi rozstrzygania,
// w kolejności zapisu w tabeli ustawienie.
func (r *repozytoriumKonfiguracji) ListaOsi(ctx context.Context, poziom shared.ConfigScope,
	kluczZasiegu string, os shared.ConfigAxis, kluczOsi string) ([]Ustawienie, error) {

	kod, kodOsi, err := adresZlozony(poziom, os)
	if err != nil {
		return nil, err
	}
	polecenie, err := r.zapytania.przygotuj(ctx, listaUstawienPoziomu)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, kod, kluczZasiegu, kodOsi, kluczOsi,
		KontoOperatora(ctx))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać ustawień poziomu %q osi %q: %w",
			kod, kodOsi, err)
	}
	defer wiersze.Close()

	lista := []Ustawienie{}
	for wiersze.Next() {
		ustawienie, err := odczytajUstawienie(wiersze)
		if err != nil {
			return nil, err
		}
		lista = append(lista, ustawienie)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt ustawień poziomu %q osi %q: %w",
			kod, kodOsi, err)
	}
	return lista, nil
}

// UsunOsi kasuje ustawienie spod adresu złożonego. Skasowanie wiersza przywraca
// wartość poziomu szerszego, a w ostateczności wartość domyślną.
func (r *repozytoriumKonfiguracji) UsunOsi(ctx context.Context, poziom shared.ConfigScope,
	kluczZasiegu string, os shared.ConfigAxis, kluczOsi, klucz string) error {

	kod, kodOsi, err := adresZlozony(poziom, os)
	if err != nil {
		return err
	}
	polecenie, err := r.zapytania.przygotuj(ctx, usunUstawienie)
	if err != nil {
		return err
	}
	_, err = polecenie.ExecContext(ctx, kod, kluczZasiegu, kodOsi, kluczOsi, klucz,
		KontoOperatora(ctx))
	if err != nil {
		return fmt.Errorf("dane: nie można usunąć ustawienia %q poziomu %q osi %q: %w",
			klucz, kod, kodOsi, err)
	}
	return nil
}

// adresZlozony przekłada parę poziom i oś rozstrzygania na wartości kolumn
// tabeli ustawienie, gotowe do zapytania.
func adresZlozony(poziom shared.ConfigScope, os shared.ConfigAxis) (string, string, error) {
	kod, err := poziomZasieguNaBaze(poziom)
	if err != nil {
		return "", "", err
	}
	kodOsi, err := osZasieguNaBaze(os)
	if err != nil {
		return "", "", err
	}
	return kod, kodOsi, nil
}
