// Odpowiedzialność pliku: dostęp do obszaru sesji (tabela `sesja`).
// Sesja pozostaje wspólna dla plików, pamięci, projektu i agentów; jednostką
// wykonania jest okno komunikacji, nie sesja.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"danacoconsole/shared"
)

// Sesja to wiersz tabeli `sesja`. Moduł nie jest tu kolumną — należy do okna
// komunikacji.
type Sesja struct {
	ID           int64
	KartaSesjiID int64
	Tytul        string
	Projekt      *string
	Stan         shared.SessionStatus
	// IdentyfikatorZewnetrzny wiąże wiersz z sesją rdzenia, która żyje pod
	// identyfikatorem tekstowym i pod nim wychodzi kontraktem do klienta.
	// Pusty oznacza wiersz założony wprost w bazie.
	IdentyfikatorZewnetrzny *string
	Utworzono               string
	Zaktualizowano          string
	Zakonczono              *string
}

// RepozytoriumSesji jest kontraktem obszaru sesji dla warstw wyższych.
type RepozytoriumSesji interface {
	Utworz(ctx context.Context, sesja Sesja) (int64, error)
	Pobierz(ctx context.Context, id int64) (Sesja, error)
	PoIdentyfikatorze(ctx context.Context, identyfikator string) (Sesja, error)
	Lista(ctx context.Context, kartaSesjiID int64) ([]Sesja, error)
	ZmienStan(ctx context.Context, id int64, stan shared.SessionStatus) error
	ZmienTytul(ctx context.Context, id int64, tytul string) error
	ZmienProjekt(ctx context.Context, id int64, projekt string) error
	Usun(ctx context.Context, id int64) error
}

const (
	kolumnySesji = `id, karta_sesji_id, tytul, projekt, stan, identyfikator_zewnetrzny,
	                utworzono, zaktualizowano, zakonczono`

	wstawSesje = `INSERT INTO sesja (karta_sesji_id, tytul, projekt, stan, identyfikator_zewnetrzny)
	              VALUES (?, ?, ?, ?, ?)`

	pobierzSesje = `SELECT ` + kolumnySesji + ` FROM sesja WHERE id = ?`

	sesjaPoIdentyfikatorze = `SELECT ` + kolumnySesji + ` FROM sesja
	                          WHERE identyfikator_zewnetrzny = ?`

	// Sesje w koszu (kolumna `usunieto_o`) nie wchodzą do wykazu sesji żywych —
	// przez ten jeden warunek znikają z session.list, z archiwum i z odtworzenia
	// rejestru po restarcie. Widzi je wyłącznie repozytorium kosza
	// (budowa/server/internal/dane/sesje_kosz.go); odczyt po identyfikatorze ich
	// nie kryje, bo po nim odbywa się przywrócenie i czyszczenie.
	listaSesji = `SELECT ` + kolumnySesji + ` FROM sesja
	              WHERE (? = 0 OR karta_sesji_id = ?) AND usunieto_o IS NULL
	              ORDER BY id`

	zmienStanSesji = `UPDATE sesja
	                  SET stan = ?,
	                      zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now'),
	                      zakonczono = CASE WHEN ? IN ('zakonczona','archiwalna')
	                                        THEN strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                                        ELSE zakonczono END
	                  WHERE id = ?`

	usunSesje = `DELETE FROM sesja WHERE id = ?`
)

type repozytoriumSesji struct {
	zapytania *zapytania
}

func noweRepozytoriumSesji(z *zapytania) *repozytoriumSesji {
	return &repozytoriumSesji{zapytania: z}
}

// Utworz zakłada sesję i zwraca jej identyfikator.
func (r *repozytoriumSesji) Utworz(ctx context.Context, sesja Sesja) (int64, error) {
	stan, err := stanSesjiNaBaze(sesja.Stan)
	if err != nil {
		return 0, err
	}
	polecenie, err := r.zapytania.przygotuj(ctx, wstawSesje)
	if err != nil {
		return 0, err
	}
	wynik, err := polecenie.ExecContext(ctx, sesja.KartaSesjiID, sesja.Tytul,
		tekstDoKolumny(sesja.Projekt), stan, tekstDoKolumny(sesja.IdentyfikatorZewnetrzny))
	if err != nil {
		return 0, fmt.Errorf("dane: nie można zapisać sesji %q: %w", sesja.Tytul, err)
	}
	return wynik.LastInsertId()
}

// Pobierz zwraca sesję o wskazanym identyfikatorze.
func (r *repozytoriumSesji) Pobierz(ctx context.Context, id int64) (Sesja, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzSesje)
	if err != nil {
		return Sesja{}, err
	}
	sesja, err := odczytajSesje(polecenie.QueryRowContext(ctx, id))
	if errors.Is(err, sql.ErrNoRows) {
		return Sesja{}, fmt.Errorf("dane: sesja %d nie istnieje", id)
	}
	return sesja, err
}

// PoIdentyfikatorze zwraca sesję po identyfikatorze nadanym przez rdzeń. Jest to
// jedyna droga odnalezienia wiersza sesji po restarcie procesu — rdzeń zna
// wyłącznie identyfikator tekstowy, nie klucz główny wiersza.
func (r *repozytoriumSesji) PoIdentyfikatorze(ctx context.Context, identyfikator string) (Sesja, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, sesjaPoIdentyfikatorze)
	if err != nil {
		return Sesja{}, err
	}
	sesja, err := odczytajSesje(polecenie.QueryRowContext(ctx, identyfikator))
	if errors.Is(err, sql.ErrNoRows) {
		return Sesja{}, fmt.Errorf("%w: sesja %q", ErrBrakWiersza, identyfikator)
	}
	return sesja, err
}

// Lista zwraca sesje karty. Wartość 0 oznacza wszystkie karty.
func (r *repozytoriumSesji) Lista(ctx context.Context, kartaSesjiID int64) ([]Sesja, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaSesji)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, kartaSesjiID, kartaSesjiID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać listy sesji: %w", err)
	}
	defer wiersze.Close()

	lista := []Sesja{}
	for wiersze.Next() {
		sesja, err := odczytajSesje(wiersze)
		if err != nil {
			return nil, err
		}
		lista = append(lista, sesja)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt listy sesji: %w", err)
	}
	return lista, nil
}

// ZmienStan zapisuje nowy stan sesji; stan końcowy odnotowuje czas zakończenia.
func (r *repozytoriumSesji) ZmienStan(ctx context.Context, id int64, stan shared.SessionStatus) error {
	kolumna, err := stanSesjiNaBaze(stan)
	if err != nil {
		return err
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zmienStanSesji)
	if err != nil {
		return err
	}
	wynik, err := polecenie.ExecContext(ctx, kolumna, kolumna, id)
	if err != nil {
		return fmt.Errorf("dane: nie można zmienić stanu sesji %d: %w", id, err)
	}
	return sprawdzTrafienie(wynik, "sesja", id)
}

// Usun kasuje sesję wraz z oknami i wiadomościami (kaskada schematu).
// Prymityw warstwy danych — torem produktu jest kosz: session.delete stawia
// znacznik, a fizyczny DELETE wykonuje czyszczenie po terminie
// w budowa/server/internal/dane/sesje_kosz.go, bo tylko ono sprząta też bloki
// wiadomości.
func (r *repozytoriumSesji) Usun(ctx context.Context, id int64) error {
	polecenie, err := r.zapytania.przygotuj(ctx, usunSesje)
	if err != nil {
		return err
	}
	wynik, err := polecenie.ExecContext(ctx, id)
	if err != nil {
		return fmt.Errorf("dane: nie można usunąć sesji %d: %w", id, err)
	}
	return sprawdzTrafienie(wynik, "sesja", id)
}

// odczytajSesje składa strukturę z jednego wiersza wyniku.
func odczytajSesje(wiersz skaner) (Sesja, error) {
	var sesja Sesja
	var projekt, zakonczono, identyfikator sql.NullString
	var stan string
	err := wiersz.Scan(&sesja.ID, &sesja.KartaSesjiID, &sesja.Tytul, &projekt, &stan,
		&identyfikator, &sesja.Utworzono, &sesja.Zaktualizowano, &zakonczono)
	if err != nil {
		return Sesja{}, err
	}
	sesja.Projekt = tekstZKolumny(projekt)
	sesja.IdentyfikatorZewnetrzny = tekstZKolumny(identyfikator)
	sesja.Zakonczono = tekstZKolumny(zakonczono)
	sesja.Stan, err = stanSesjiZBazy(stan)
	if err != nil {
		return Sesja{}, err
	}
	return sesja, nil
}

// sprawdzTrafienie odróżnia zapis wykonany od zapisu, który nie trafił w wiersz.
func sprawdzTrafienie(wynik sql.Result, tabela string, id int64) error {
	liczba, err := wynik.RowsAffected()
	if err != nil {
		return fmt.Errorf("dane: nieznana liczba zmienionych wierszy tabeli %s: %w", tabela, err)
	}
	if liczba == 0 {
		return fmt.Errorf("dane: wiersz %d tabeli %s nie istnieje", id, tabela)
	}
	return nil
}
