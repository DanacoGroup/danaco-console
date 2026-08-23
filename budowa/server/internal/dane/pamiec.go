// Odpowiedzialność pliku: dostęp do zasobów pamięci wielopoziomowej (tabela
// `zasob_pamieci`). Pamięć zostaje na poziomie sesji — okno komunikacji nie jest
// jej poziomem. Treść obszerna trafia do pliku, baza trzyma odwołanie
// . Konfigurację pamięci sesji obsługuje `pamiec_sesji.go`.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// PoziomPamieci jest słownikiem schematu obszaru pamięci. Kontrakt nie ma dla
// niego wyliczenia — pamięć nie ma jeszcze komend.
type PoziomPamieci string

// Katalog poziomów pamięci niesie kolumna `zasob_pamieci.poziom`; wartość składa
// się z napisu odczytanego z bazy — zob. rozlozPoziomy w pamiec_sesji.go.

// Zasob to wiersz tabeli `zasob_pamieci`.
type Zasob struct {
	ID             int64
	Poziom         PoziomPamieci
	KluczZasiegu   string
	Klucz          string
	Tresc          *string
	TrescOdwolanie *string
	Waga           int
	Utworzono      string
	Zaktualizowano string
}

// RepozytoriumPamieci jest kontraktem obszaru pamięci.
type RepozytoriumPamieci interface {
	Zapisz(ctx context.Context, zasob Zasob) error
	Pobierz(ctx context.Context, poziom PoziomPamieci, kluczZasiegu, klucz string) (Zasob, bool, error)
	ListaPoziomu(ctx context.Context, poziom PoziomPamieci, kluczZasiegu string) ([]Zasob, error)
	Usun(ctx context.Context, poziom PoziomPamieci, kluczZasiegu, klucz string) error
	UstawKonfiguracjeSesji(ctx context.Context, konfiguracja KonfiguracjaPamieci) error
	KonfiguracjaSesji(ctx context.Context, sesjaID int64) (KonfiguracjaPamieci, bool, error)
}

const (
	kolumnyZasobu = `id, poziom, klucz_zasiegu, klucz, tresc, tresc_odwolanie, waga,
	                 utworzono, zaktualizowano`

	zapiszZasobPamieci = `INSERT INTO zasob_pamieci
	                      (poziom, klucz_zasiegu, klucz, tresc, tresc_odwolanie, waga)
	                      VALUES (?, ?, ?, ?, ?, ?)
	                      ON CONFLICT(poziom, klucz_zasiegu, klucz) DO UPDATE SET
	                          tresc = excluded.tresc,
	                          tresc_odwolanie = excluded.tresc_odwolanie,
	                          waga = excluded.waga,
	                          zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')`

	pobierzZasobPamieci = `SELECT ` + kolumnyZasobu + ` FROM zasob_pamieci
	                       WHERE poziom = ? AND klucz_zasiegu = ? AND klucz = ?`

	listaZasobowPoziomu = `SELECT ` + kolumnyZasobu + ` FROM zasob_pamieci
	                       WHERE poziom = ? AND klucz_zasiegu = ?
	                       ORDER BY waga DESC, klucz`

	usunZasobPamieci = `DELETE FROM zasob_pamieci
	                    WHERE poziom = ? AND klucz_zasiegu = ? AND klucz = ?`
)

type repozytoriumPamieci struct {
	zapytania *zapytania
}

func noweRepozytoriumPamieci(z *zapytania) *repozytoriumPamieci {
	return &repozytoriumPamieci{zapytania: z}
}

// Zapisz utrwala zasób pamięci; zasób o tym samym kluczu nadpisuje.
func (r *repozytoriumPamieci) Zapisz(ctx context.Context, zasob Zasob) error {
	if zasob.Klucz == "" {
		return fmt.Errorf("dane: zasób pamięci bez klucza")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszZasobPamieci)
	if err != nil {
		return err
	}
	_, err = polecenie.ExecContext(ctx, string(zasob.Poziom), zasob.KluczZasiegu, zasob.Klucz,
		tekstDoKolumny(zasob.Tresc), tekstDoKolumny(zasob.TrescOdwolanie), zasob.Waga)
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać zasobu pamięci %q poziomu %q: %w",
			zasob.Klucz, zasob.Poziom, err)
	}
	return nil
}

// Pobierz zwraca zasób pamięci. Drugi wynik mówi, czy zasób istnieje — brak
// wpisu nie jest błędem.
func (r *repozytoriumPamieci) Pobierz(ctx context.Context, poziom PoziomPamieci,
	kluczZasiegu, klucz string) (Zasob, bool, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzZasobPamieci)
	if err != nil {
		return Zasob{}, false, err
	}
	zasob, err := odczytajZasob(polecenie.QueryRowContext(ctx, string(poziom), kluczZasiegu, klucz))
	if errors.Is(err, sql.ErrNoRows) {
		return Zasob{}, false, nil
	}
	if err != nil {
		return Zasob{}, false, err
	}
	return zasob, true, nil
}

// ListaPoziomu zwraca zasoby jednego bytu poziomu, od najwyższej wagi.
func (r *repozytoriumPamieci) ListaPoziomu(ctx context.Context, poziom PoziomPamieci,
	kluczZasiegu string) ([]Zasob, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaZasobowPoziomu)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, string(poziom), kluczZasiegu)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać pamięci poziomu %q: %w", poziom, err)
	}
	defer wiersze.Close()

	lista := []Zasob{}
	for wiersze.Next() {
		zasob, err := odczytajZasob(wiersze)
		if err != nil {
			return nil, err
		}
		lista = append(lista, zasob)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt pamięci poziomu %q: %w", poziom, err)
	}
	return lista, nil
}

// Usun kasuje zasób pamięci.
func (r *repozytoriumPamieci) Usun(ctx context.Context, poziom PoziomPamieci,
	kluczZasiegu, klucz string) error {

	polecenie, err := r.zapytania.przygotuj(ctx, usunZasobPamieci)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, string(poziom), kluczZasiegu, klucz); err != nil {
		return fmt.Errorf("dane: nie można usunąć zasobu pamięci %q: %w", klucz, err)
	}
	return nil
}

// odczytajZasob składa strukturę z jednego wiersza wyniku.
func odczytajZasob(wiersz skaner) (Zasob, error) {
	var zasob Zasob
	var poziom string
	var tresc, odwolanie sql.NullString
	err := wiersz.Scan(&zasob.ID, &poziom, &zasob.KluczZasiegu, &zasob.Klucz, &tresc, &odwolanie,
		&zasob.Waga, &zasob.Utworzono, &zasob.Zaktualizowano)
	if err != nil {
		return Zasob{}, err
	}
	zasob.Poziom = PoziomPamieci(poziom)
	zasob.Tresc = tekstZKolumny(tresc)
	zasob.TrescOdwolanie = tekstZKolumny(odwolanie)
	return zasob, nil
}
