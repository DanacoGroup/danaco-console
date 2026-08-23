// Odpowiedzialność pliku: nadania dostępu (tabela `nadanie_dostepu`) —
// struktura, kontrakt repozytorium i odczyt. Zapis leży w `dostep_nadania_zapis.go`.
//
// Nadanie wiąże okno komunikacji z punktem dostępu. Żyje per okno rozmowy, nie
// per sesja i nie per platforma. Okno ma zbiór nadań — kolejność i oznaczenie
// głównego niosą znaczenie. Okno bez nadań pracuje dalej, tylko niczego nie
// widzi.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"danacoconsole/shared"
)

// Nadanie to wiersz tabeli `nadanie_dostepu` wraz z listą korzeni. Lista pusta
// znaczy „komplet korzeni punktu", zgodnie z opisem pola `roots` struktury
// AccessGrant kontraktu.
type Nadanie struct {
	ID                int64
	OknoKomunikacjiID int64
	PunktDostepuID    int64
	Tryb              shared.AccessMode
	Korzenie          []string
	Kolejnosc         int
	Glowne            bool
	Aktywne           bool
	// IdentyfikatorZewnetrzny wiąże wiersz z nadaniem rdzenia, które żyje pod
	// identyfikatorem tekstowym. NULL oznacza wiersz założony wprost w bazie, bez
	// odpowiednika w pamięci rdzenia.
	IdentyfikatorZewnetrzny *string
	Utworzono               string
}

// RepozytoriumNadan jest kontraktem obszaru nadań dostępu.
type RepozytoriumNadan interface {
	ListaOkna(ctx context.Context, oknoID int64, tylkoAktywne bool) ([]Nadanie, error)
	Pobierz(ctx context.Context, id int64) (Nadanie, error)
	PoIdentyfikatorze(ctx context.Context, identyfikator string) (Nadanie, error)
	Dodaj(ctx context.Context, nadanie Nadanie) (int64, error)
	Aktualizuj(ctx context.Context, nadanie Nadanie) error
	OznaczGlowne(ctx context.Context, id int64) error
	Usun(ctx context.Context, id int64) error
}

const (
	kolumnyNadania = `id, okno_komunikacji_id, punkt_dostepu_id, tryb, kolejnosc, glowne,
	                  aktywne, identyfikator_zewnetrzny, utworzono`

	listaNadanOkna = `SELECT ` + kolumnyNadania + ` FROM nadanie_dostepu
	                  WHERE okno_komunikacji_id = ? AND (? = 0 OR aktywne = 1)
	                  ORDER BY kolejnosc, id`

	pobierzNadanie = `SELECT ` + kolumnyNadania + ` FROM nadanie_dostepu WHERE id = ?`

	nadaniePoIdentyfikatorze = `SELECT ` + kolumnyNadania + ` FROM nadanie_dostepu
	                            WHERE identyfikator_zewnetrzny = ?`
)

type repozytoriumNadan struct {
	zapytania *zapytania
	db        *sql.DB
}

// Zgodność implementacji z kontraktem sprawdzana jest przy kompilacji, a nie
// dopiero przy złożeniu zestawu repozytoriów.
var _ RepozytoriumNadan = (*repozytoriumNadan)(nil)

// noweRepozytoriumNadan zakłada repozytorium nadań dostępu.
func noweRepozytoriumNadan(z *zapytania, db *sql.DB) *repozytoriumNadan {
	return &repozytoriumNadan{zapytania: z, db: db}
}

// ListaOkna zwraca zbiór nadań jednego okna w kolejności zapisanej przez
// Operatora. Zbiór pusty nie jest błędem.
func (r *repozytoriumNadan) ListaOkna(ctx context.Context, oknoID int64,
	tylkoAktywne bool) ([]Nadanie, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaNadanOkna)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, oknoID, liczbaLogiczna(tylkoAktywne))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać nadań okna %d: %w", oknoID, err)
	}
	defer wiersze.Close()

	lista := []Nadanie{}
	for wiersze.Next() {
		nadanie, err := odczytajNadanie(wiersze)
		if err != nil {
			return nil, err
		}
		lista = append(lista, nadanie)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt nadań okna %d: %w", oknoID, err)
	}
	for i, nadanie := range lista {
		if lista[i].Korzenie, err = r.korzenieNadania(ctx, nadanie.ID); err != nil {
			return nil, err
		}
	}
	return lista, nil
}

// Pobierz zwraca nadanie wskazane kluczem wiersza.
func (r *repozytoriumNadan) Pobierz(ctx context.Context, id int64) (Nadanie, error) {
	return r.jedno(ctx, pobierzNadanie, fmt.Sprintf("%d", id), id)
}

// PoIdentyfikatorze zwraca nadanie wskazane identyfikatorem rdzenia.
func (r *repozytoriumNadan) PoIdentyfikatorze(ctx context.Context, identyfikator string) (Nadanie, error) {
	return r.jedno(ctx, nadaniePoIdentyfikatorze, fmt.Sprintf("%q", identyfikator), identyfikator)
}

// jedno odczytuje pojedyncze nadanie wraz z jego korzeniami.
func (r *repozytoriumNadan) jedno(ctx context.Context, zapytanie, opis string,
	argument any) (Nadanie, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return Nadanie{}, err
	}
	nadanie, err := odczytajNadanie(polecenie.QueryRowContext(ctx, argument))
	if errors.Is(err, sql.ErrNoRows) {
		return Nadanie{}, fmt.Errorf("dane: nadanie dostępu %s nie istnieje: %w", opis, ErrBrakWiersza)
	}
	if err != nil {
		return Nadanie{}, err
	}
	if nadanie.Korzenie, err = r.korzenieNadania(ctx, nadanie.ID); err != nil {
		return Nadanie{}, err
	}
	return nadanie, nil
}

// korzenieNadania zwraca zawężenie korzeni punktu zapisane przy nadaniu.
func (r *repozytoriumNadan) korzenieNadania(ctx context.Context, nadanieID int64) ([]string, error) {
	return wczytajKorzenie(ctx, r.zapytania, listaKorzeniNadania, nadanieID, "nadania")
}

// odczytajNadanie składa strukturę z jednego wiersza wyniku. Korzenie dokłada
// repozytorium — leżą w tabeli podrzędnej.
func odczytajNadanie(wiersz skaner) (Nadanie, error) {
	var nadanie Nadanie
	var identyfikator sql.NullString
	var tryb string
	var glowne, aktywne int
	err := wiersz.Scan(&nadanie.ID, &nadanie.OknoKomunikacjiID, &nadanie.PunktDostepuID, &tryb,
		&nadanie.Kolejnosc, &glowne, &aktywne, &identyfikator, &nadanie.Utworzono)
	if err != nil {
		return Nadanie{}, err
	}
	if nadanie.Tryb, err = trybDostepuZBazy(tryb, "nadanie_dostepu.tryb"); err != nil {
		return Nadanie{}, err
	}
	nadanie.IdentyfikatorZewnetrzny = tekstZKolumny(identyfikator)
	nadanie.Glowne = glowne != 0
	nadanie.Aktywne = aktywne != 0
	nadanie.Korzenie = []string{}
	return nadanie, nil
}
