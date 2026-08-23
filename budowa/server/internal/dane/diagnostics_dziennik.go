// Odpowiedzialność pliku: dziennik modułu Diagnostics — dopisanie wpisów
// partią, odczyt zawężony filtrem, wykaz źródeł i rozkład poziomów.
//
// Zawężenie idzie parametrem, nie sklejaniem tekstu. Jedno przygotowane
// zapytanie obsługuje pięć zawężeń naraz, bo pusty parametr znaczy „nie
// zawężaj". Wartości nigdy nie wchodzą do treści SQL, a pamięć podręczna
// zapytań ma jedną pozycję zamiast trzydziestu dwóch.
package dane

import (
	"context"
	"database/sql"
	"fmt"

	"danacoconsole/shared"
)

// granicaDziennika obowiązuje, gdy Operator nie poda własnej. Bez niej pierwsze
// otwarcie Logs Viewer po tygodniu pracy ściągnęłoby cały dziennik do
// przeglądarki.
const granicaDziennika = 500

const (
	kolumnyWpisu = `kod, chwila, poziom, zrodlo, tresc, sesja_kod, okno_kod, proces_kod, odcisk`

	warunkiWpisu = ` WHERE (? = 0 OR chwila >= ?)
	                   AND (? = 0 OR chwila <= ?)
	                   AND (? = '' OR poziom = ?)
	                   AND (? = '' OR zrodlo = ?)
	                   AND (? = '' OR instr(lower(tresc), lower(?)) > 0)`

	wstawWpisDiagnostyki = `INSERT INTO diagnostyka_wpis
	                        (kod, chwila, poziom, zrodlo, tresc, sesja_kod, okno_kod, proces_kod, odcisk)
	                        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	pobierzWpisy = `SELECT ` + kolumnyWpisu + `, 1 FROM diagnostyka_wpis` + warunkiWpisu +
		` ORDER BY chwila DESC, id DESC LIMIT CASE WHEN ? > 0 THEN ? ELSE -1 END`

	// Wariant scalający: SQLite w grupowaniu z MAX() oddaje kolumny gołe
	// z wiersza, na którym maksimum wypadło — wpis reprezentujący grupę jest
	// więc wpisem najnowszym, a nie przypadkowym.
	pobierzWpisyScalone = `SELECT kod, MAX(chwila), poziom, zrodlo, tresc, sesja_kod, okno_kod,
	                              proces_kod, odcisk, COUNT(*)
	                       FROM diagnostyka_wpis` + warunkiWpisu +
		` GROUP BY odcisk ORDER BY 2 DESC LIMIT CASE WHEN ? > 0 THEN ? ELSE -1 END`

	policzWpisy = `SELECT COUNT(*) FROM diagnostyka_wpis` + warunkiWpisu

	policzWpisyScalone = `SELECT COUNT(DISTINCT odcisk) FROM diagnostyka_wpis` + warunkiWpisu

	pobierzZrodlaWpisow = `SELECT DISTINCT zrodlo FROM diagnostyka_wpis
	                       WHERE zrodlo IS NOT NULL AND zrodlo <> '' ORDER BY zrodlo`

	policzPoziomyWpisow = `SELECT poziom, COUNT(*) FROM diagnostyka_wpis
	                       WHERE (? = 0 OR chwila >= ?) AND (? = 0 OR chwila <= ?)
	                       GROUP BY poziom`
)

type repozytoriumDiagnostyki struct {
	zapytania *zapytania
	db        *sql.DB
}

func noweRepozytoriumDiagnostyki(z *zapytania, db *sql.DB) *repozytoriumDiagnostyki {
	return &repozytoriumDiagnostyki{zapytania: z, db: db}
}

// DopiszWpisy zapisuje partię wpisów w jednej transakcji. Partia, a nie wpis
// pojedynczy: dziennik rdzenia bywa gęsty, a osobna transakcja na każdą linię
// zamieniłaby zapis dziennika w wąskie gardło pracy rdzenia.
func (r *repozytoriumDiagnostyki) DopiszWpisy(ctx context.Context, wpisy []WpisDiagnostyki) error {
	if len(wpisy) == 0 {
		return nil
	}
	return wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, wstawWpisDiagnostyki)
		if err != nil {
			return err
		}
		for _, wpis := range wpisy {
			_, err := polecenie.ExecContext(ctx, wpis.Kod, wpis.Chwila, string(wpis.Poziom),
				tekstDoKolumny(wpis.Zrodlo), wpis.Tresc, tekstDoKolumny(wpis.SesjaKod),
				tekstDoKolumny(wpis.OknoKod), tekstDoKolumny(wpis.ProcesKod), wpis.Odcisk)
			if err != nil {
				return fmt.Errorf("dane: nie można dopisać wpisu dziennika %q: %w", wpis.Kod, err)
			}
		}
		return nil
	})
}

// Wpisy zwraca dziennik zawężony filtrem oraz liczbę wpisów spełniających
// warunki przed ucięciem granicą — klient odróżnia „tyle jest" od „tyle widać".
func (r *repozytoriumDiagnostyki) Wpisy(ctx context.Context,
	filtr FiltrDziennika) ([]WpisDiagnostyki, int, error) {

	zapytanie, zliczanie := pobierzWpisy, policzWpisy
	if filtr.Scal {
		zapytanie, zliczanie = pobierzWpisyScalone, policzWpisyScalone
	}
	granica := filtr.Granica
	if granica <= 0 {
		granica = granicaDziennika
	}
	warunki := warunkiDziennika(filtr)

	razem, err := r.policz(ctx, zliczanie, warunki)
	if err != nil {
		return nil, 0, err
	}

	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return nil, 0, err
	}
	wiersze, err := polecenie.QueryContext(ctx, append(warunki, granica, granica)...)
	if err != nil {
		return nil, 0, fmt.Errorf("dane: nie można odczytać dziennika diagnostyki: %w", err)
	}
	defer wiersze.Close()

	wpisy := make([]WpisDiagnostyki, 0, 64)
	for wiersze.Next() {
		wpis, err := odczytajWpisDiagnostyki(wiersze)
		if err != nil {
			return nil, 0, fmt.Errorf("dane: nieczytelny wpis dziennika diagnostyki: %w", err)
		}
		wpisy = append(wpisy, wpis)
	}
	return wpisy, razem, wiersze.Err()
}

// ZrodlaWpisow zwraca źródła obecne w dzienniku, uporządkowane alfabetycznie.
func (r *repozytoriumDiagnostyki) ZrodlaWpisow(ctx context.Context) ([]string, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzZrodlaWpisow)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać źródeł dziennika: %w", err)
	}
	defer wiersze.Close()

	zrodla := make([]string, 0, 8)
	for wiersze.Next() {
		var zrodlo string
		if err := wiersze.Scan(&zrodlo); err != nil {
			return nil, fmt.Errorf("dane: nieczytelne źródło dziennika: %w", err)
		}
		zrodla = append(zrodla, zrodlo)
	}
	return zrodla, wiersze.Err()
}

// PoziomyWpisow zwraca rozkład wpisów po poziomach w zakresie czasu.
func (r *repozytoriumDiagnostyki) PoziomyWpisow(ctx context.Context, od, do int64) (LicznikPoziomow, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, policzPoziomyWpisow)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, od, od, do, do)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można policzyć poziomów dziennika: %w", err)
	}
	defer wiersze.Close()

	licznik := LicznikPoziomow{}
	for wiersze.Next() {
		var poziom shared.LogLevel
		var liczba int
		if err := wiersze.Scan(&poziom, &liczba); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny rozkład poziomów: %w", err)
		}
		licznik[poziom] = liczba
	}
	return licznik, wiersze.Err()
}

// policz wykonuje zapytanie zliczające na tych samych warunkach co odczyt.
func (r *repozytoriumDiagnostyki) policz(ctx context.Context, zapytanie string, warunki []any) (int, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return 0, err
	}
	var razem int
	if err := polecenie.QueryRowContext(ctx, warunki...).Scan(&razem); err != nil {
		return 0, fmt.Errorf("dane: nie można policzyć wpisów dziennika: %w", err)
	}
	return razem, nil
}

// warunkiDziennika składa argumenty zawężeń w kolejności znaków zapytania.
func warunkiDziennika(filtr FiltrDziennika) []any {
	poziom, zrodlo, wzorzec := string(filtr.Poziom), filtr.Zrodlo, filtr.Wzorzec
	return []any{
		filtr.Od, filtr.Od,
		filtr.Do, filtr.Do,
		poziom, poziom,
		zrodlo, zrodlo,
		wzorzec, wzorzec,
	}
}

// odczytajWpisDiagnostyki przenosi wiersz dziennika do bytu obszaru.
func odczytajWpisDiagnostyki(s skaner) (WpisDiagnostyki, error) {
	var wpis WpisDiagnostyki
	err := s.Scan(&wpis.Kod, &wpis.Chwila, &wpis.Poziom, &wpis.Zrodlo, &wpis.Tresc,
		&wpis.SesjaKod, &wpis.OknoKod, &wpis.ProcesKod, &wpis.Odcisk, &wpis.Powtorzenia)
	return wpis, err
}
