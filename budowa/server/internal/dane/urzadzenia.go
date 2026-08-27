// Plik prowadzi katalog urządzeń w tabeli urzadzenie wraz z kontraktem
// repozytorium i odczytem, podczas gdy zapis leży w pliku urzadzenia_zapis.go.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// Urzadzenie to wiersz tabeli urzadzenie: nazwa hosta i identyfikator
// sprzętowy są faktami maszyny ustalanymi przy rozpoznaniu, a pole biezace
// oznacza maszynę, na której działa rdzeń.
type Urzadzenie struct {
	ID                     int64
	Nazwa                  string
	IdentyfikatorSprzetowy string
	NazwaHosta             string
	SystemOperacyjny       string
	WersjaKlienta          string
	Zaufane                bool
	Biezace                bool
	OstatnioWidziane       *string
	Utworzono              string
}

// RepozytoriumUrzadzen jest kontraktem katalogu urządzeń, wskazującym operacje
// repozytorium dostępne warstwie wyższej rdzenia.
type RepozytoriumUrzadzen interface {
	Lista(ctx context.Context) ([]Urzadzenie, error)
	Pobierz(ctx context.Context, id int64) (Urzadzenie, error)
	PoIdentyfikatorze(ctx context.Context, identyfikator string) (Urzadzenie, error)
	Biezace(ctx context.Context) (Urzadzenie, error)
	Dodaj(ctx context.Context, urzadzenie Urzadzenie) (int64, error)
	Aktualizuj(ctx context.Context, urzadzenie Urzadzenie) error
	// ZapewnijBiezace zakłada albo odświeża wiersz maszyny bieżącej rdzenia.
	ZapewnijBiezace(ctx context.Context, urzadzenie Urzadzenie) (Urzadzenie, error)
}

const (
	kolumnyUrzadzenia = `id, nazwa, identyfikator_sprzetowy, nazwa_hosta, system_operacyjny,
	                     wersja_klienta, zaufane, biezace, ostatnio_widziane, utworzono`

	listaUrzadzen = `SELECT ` + kolumnyUrzadzenia + ` FROM urzadzenie ORDER BY nazwa, id`

	pobierzUrzadzenie = `SELECT ` + kolumnyUrzadzenia + ` FROM urzadzenie WHERE id = ?`

	urzadzeniePoIdentyfikatorze = `SELECT ` + kolumnyUrzadzenia + ` FROM urzadzenie
	                               WHERE identyfikator_sprzetowy = ?`

	urzadzenieBiezace = `SELECT ` + kolumnyUrzadzenia + ` FROM urzadzenie WHERE biezace = 1`
)

type repozytoriumUrzadzen struct {
	zapytania *zapytania
	db        *sql.DB
}

// Zgodność implementacji repozytoriumUrzadzen z kontraktem RepozytoriumUrzadzen
// jest sprawdzana wskazaniem typu przy kompilacji pakietu.
var _ RepozytoriumUrzadzen = (*repozytoriumUrzadzen)(nil)

// noweRepozytoriumUrzadzen zakłada repozytorium katalogu urządzeń nad pamięcią
// przygotowanych zapytań i połączeniem bazy danych.
func noweRepozytoriumUrzadzen(z *zapytania, db *sql.DB) *repozytoriumUrzadzen {
	return &repozytoriumUrzadzen{zapytania: z, db: db}
}

// Lista zwraca komplet urządzeń w kolejności nazw. Katalog pusty nie jest błędem:
// platforma bez zarejestrowanej maszyny pracuje dalej, tylko nie ma na czym
// osadzić katalogu lokalnego.
func (r *repozytoriumUrzadzen) Lista(ctx context.Context) ([]Urzadzenie, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaUrzadzen)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać katalogu urządzeń: %w", err)
	}
	defer wiersze.Close()

	lista := []Urzadzenie{}
	for wiersze.Next() {
		urzadzenie, err := odczytajUrzadzenie(wiersze)
		if err != nil {
			return nil, err
		}
		lista = append(lista, urzadzenie)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt katalogu urządzeń: %w", err)
	}
	return lista, nil
}

// Pobierz zwraca urządzenie wskazane kluczem wiersza — tym samym, który niesie
// kolumna `punkt_dostepu.urzadzenie_id`.
func (r *repozytoriumUrzadzen) Pobierz(ctx context.Context, id int64) (Urzadzenie, error) {
	return r.jedno(ctx, pobierzUrzadzenie, fmt.Sprintf("%d", id), id)
}

// PoIdentyfikatorze zwraca urządzenie o wskazanym identyfikatorze sprzętowym.
// Identyfikator jest trwałym rozpoznaniem maszyny, klucz wiersza — jej numerem
// w tej bazie.
func (r *repozytoriumUrzadzen) PoIdentyfikatorze(ctx context.Context, identyfikator string) (Urzadzenie, error) {
	return r.jedno(ctx, urzadzeniePoIdentyfikatorze, fmt.Sprintf("%q", identyfikator), identyfikator)
}

// Biezace zwraca maszynę, na której działa rdzeń. Brak wiersza znaczy, że
// rozpoznanie startowe jeszcze nie przebiegło — warstwa wyższa odróżnia ten
// przypadek przez ErrBrakWiersza, nie przez treść komunikatu.
func (r *repozytoriumUrzadzen) Biezace(ctx context.Context) (Urzadzenie, error) {
	return r.jedno(ctx, urzadzenieBiezace, "bieżące")
}

// jedno wykonuje zapytanie zwracające najwyżej jeden wiersz urządzenia
// i tłumaczy brak wiersza na błąd ErrBrakWiersza.
func (r *repozytoriumUrzadzen) jedno(ctx context.Context, zapytanie, opis string,
	argumenty ...any) (Urzadzenie, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return Urzadzenie{}, err
	}
	urzadzenie, err := odczytajUrzadzenie(polecenie.QueryRowContext(ctx, argumenty...))
	if errors.Is(err, sql.ErrNoRows) {
		return Urzadzenie{}, fmt.Errorf("dane: urządzenie %s nie istnieje: %w", opis, ErrBrakWiersza)
	}
	if err != nil {
		return Urzadzenie{}, fmt.Errorf("dane: nieczytelny wiersz urządzenia %s: %w", opis, err)
	}
	return urzadzenie, nil
}

// odczytajUrzadzenie składa strukturę z jednego wiersza wyniku. Kolumny
// opcjonalne (`system_operacyjny`, `wersja_klienta`) czytamy jako napis pusty —
// dla warstw wyższych brak rozpoznania i rozpoznanie puste znaczą to samo.
func odczytajUrzadzenie(wiersz skaner) (Urzadzenie, error) {
	var urzadzenie Urzadzenie
	var system, wersja, widziane sql.NullString
	var zaufane, biezace int
	err := wiersz.Scan(&urzadzenie.ID, &urzadzenie.Nazwa, &urzadzenie.IdentyfikatorSprzetowy,
		&urzadzenie.NazwaHosta, &system, &wersja, &zaufane, &biezace, &widziane,
		&urzadzenie.Utworzono)
	if err != nil {
		return Urzadzenie{}, err
	}
	urzadzenie.SystemOperacyjny = system.String
	urzadzenie.WersjaKlienta = wersja.String
	urzadzenie.Zaufane = zaufane != 0
	urzadzenie.Biezace = biezace != 0
	urzadzenie.OstatnioWidziane = tekstZKolumny(widziane)
	return urzadzenie, nil
}
