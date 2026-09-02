// Odpowiedzialność pliku: dostęp do obszaru sesji; jednostką wykonania jest okno komunikacji, nie sama sesja, choć sesja pozostaje wspólna.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"danacoconsole/shared"
)

// Sesja to wiersz tabeli sesja; moduł nie jest tu kolumną, bo należy do okna komunikacji, nie do sesji.
type Sesja struct {
	ID           int64
	KartaSesjiID int64
	Tytul        string
	Projekt      *string
	Stan         shared.SessionStatus
	// IdentyfikatorZewnetrzny wiąże wiersz z sesją rdzenia; pusty oznacza wiersz założony wprost w bazie.
	IdentyfikatorZewnetrzny *string
	Utworzono               string
	Zaktualizowano          string
	Zakonczono              *string
}

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

	sesjaIstnieje = `SELECT 1 FROM sesja WHERE id = ?`
)

// Tabela sesja nie ma kolumny konto_id: konto sięga się przez karta_sesji (migracja 407).
var (
	warunekKontaKartySesji = strings.ReplaceAll(WarunekKonta, "konto_id", "k.konto_id")

	sesjaKonta = ` EXISTS (SELECT 1 FROM karta_sesji k
	                       WHERE k.id = sesja.karta_sesji_id AND ` + warunekKontaKartySesji + `)`

	pobierzSesje = `SELECT ` + kolumnySesji + ` FROM sesja WHERE id = ? AND` + sesjaKonta

	sesjaPoIdentyfikatorze = `SELECT ` + kolumnySesji + ` FROM sesja
	                          WHERE identyfikator_zewnetrzny = ? AND` + sesjaKonta

	// Sesje w koszu nie wchodzą do wykazu sesji żywych; ten jeden warunek zdejmuje je z wykazu, archiwum i odtworzenia rejestru.
	listaSesji = `SELECT ` + kolumnySesji + ` FROM sesja
	              WHERE (? = 0 OR karta_sesji_id = ?) AND usunieto_o IS NULL AND` + sesjaKonta + `
	              ORDER BY id`

	zmienStanSesji = `UPDATE sesja
	                  SET stan = ?,
	                      zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now'),
	                      zakonczono = CASE WHEN ? IN ('zakonczona','archiwalna')
	                                        THEN strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                                        ELSE zakonczono END
	                  WHERE id = ? AND` + sesjaKonta

	usunSesje = `DELETE FROM sesja WHERE id = ? AND` + sesjaKonta
)

type repozytoriumSesji struct {
	zapytania *zapytania
}

func noweRepozytoriumSesji(z *zapytania) *repozytoriumSesji {
	return &repozytoriumSesji{zapytania: z}
}

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

func (r *repozytoriumSesji) Pobierz(ctx context.Context, id int64) (Sesja, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzSesje)
	if err != nil {
		return Sesja{}, err
	}
	sesja, err := odczytajSesje(polecenie.QueryRowContext(ctx, id, KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return Sesja{}, fmt.Errorf("dane: sesja %d nie istnieje", id)
	}
	return sesja, err
}

// PoIdentyfikatorze jest jedyną drogą odnalezienia wiersza po restarcie procesu: rdzeń zna wyłącznie identyfikator tekstowy, nie klucz główny.
func (r *repozytoriumSesji) PoIdentyfikatorze(ctx context.Context, identyfikator string) (Sesja, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, sesjaPoIdentyfikatorze)
	if err != nil {
		return Sesja{}, err
	}
	sesja, err := odczytajSesje(polecenie.QueryRowContext(ctx, identyfikator, KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return Sesja{}, fmt.Errorf("%w: sesja %q", ErrBrakWiersza, identyfikator)
	}
	return sesja, err
}

func (r *repozytoriumSesji) Lista(ctx context.Context, kartaSesjiID int64) ([]Sesja, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaSesji)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, kartaSesjiID, kartaSesjiID, KontoOperatora(ctx))
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

func (r *repozytoriumSesji) ZmienStan(ctx context.Context, id int64, stan shared.SessionStatus) error {
	kolumna, err := stanSesjiNaBaze(stan)
	if err != nil {
		return err
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zmienStanSesji)
	if err != nil {
		return err
	}
	wynik, err := polecenie.ExecContext(ctx, kolumna, kolumna, id, KontoOperatora(ctx))
	if err != nil {
		return fmt.Errorf("dane: nie można zmienić stanu sesji %d: %w", id, err)
	}
	return r.trafienieSesji(ctx, wynik, id)
}

// Usun kasuje sesję wraz z oknami i wiadomościami kaskadą schematu; torem produktu jest jednak kosz, nie usunięcie wprost.
func (r *repozytoriumSesji) Usun(ctx context.Context, id int64) error {
	polecenie, err := r.zapytania.przygotuj(ctx, usunSesje)
	if err != nil {
		return err
	}
	wynik, err := polecenie.ExecContext(ctx, id, KontoOperatora(ctx))
	if err != nil {
		return fmt.Errorf("dane: nie można usunąć sesji %d: %w", id, err)
	}
	return r.trafienieSesji(ctx, wynik, id)
}

// trafienieSesji przy zerze zmienionych wierszy odróżnia sesję innego konta (ErrKolizjaWiersza) od sesji nieistniejącej (ErrBrakWiersza).
func (r *repozytoriumSesji) trafienieSesji(ctx context.Context, wynik sql.Result, id int64) error {
	zmienione, err := wynik.RowsAffected()
	if err != nil {
		return fmt.Errorf("dane: nieznana liczba zmienionych wierszy sesji %d: %w", id, err)
	}
	if zmienione > 0 {
		return nil
	}
	polecenie, err := r.zapytania.przygotuj(ctx, sesjaIstnieje)
	if err != nil {
		return err
	}
	var jest int
	switch err := polecenie.QueryRowContext(ctx, id).Scan(&jest); {
	case errors.Is(err, sql.ErrNoRows):
		return fmt.Errorf("%w: sesja %d", ErrBrakWiersza, id)
	case err != nil:
		return fmt.Errorf("dane: nie można sprawdzić sesji %d: %w", id, err)
	}
	return fmt.Errorf("dane: sesja %d należy do innego konta: %w", id, ErrKolizjaWiersza)
}

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

// sprawdzTrafienie odróżnia zapis rzeczywiście wykonany od zapisu, który nie trafił w żaden wiersz tej tabeli.
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
