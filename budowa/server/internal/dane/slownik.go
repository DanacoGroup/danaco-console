// Odpowiedzialność pliku: terminy słownika modułu Translate (tabela
// `termin_slownika`). Ślady importu i eksportu leżą w `slownik_wymiana.go`,
// pamięć tłumaczeń w `slownik_pamiec.go` — jedno repozytorium rozdzielone na
// pliki wedle odpowiedzialności. Typ, interfejs i konstruktor deklaruje
// wyłącznie `tlumaczenie.go`; ten plik implementuje na `*repozytoriumTlumaczen`
// wyłącznie metody terminu.
//
// Wystąpienia terminu nie mają tu tabeli: `translate.glossary.occurrences` liczy
// się w locie z treści okna albo panelu przeszukanej względem `Zrodlo` terminu,
// żeby nie unieważniać zapisu przy każdej korekcie panelu.
//
// Zapis ma jedną drogę. `translate.glossary.set` nadsyła zawsze komplet zmian
// naraz — `ZapiszTerminy` przyjmuje wykaz i zapisuje go w jednej transakcji
// (`dane/transakcja.go`) przez UPSERT po `identyfikator_zewnetrzny`: termin ze
// wskazanym kodem aktualizuje się, termin bez zastanego wiersza o tym kodzie
// zakłada się — ten sam SQL obsługuje obie ścieżki, nie dwie osobne metody.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// TerminSlownika to wiersz tabeli `termin_slownika` — odpowiednik terminu
// (Glossary Term) stosowany przy generowaniu tłumaczeń. `NieTlumaczyc` niesie
// wartość logiczną wprost (kolumna INTEGER z warunkiem CHECK na 0 albo 1) —
// warstwa wyższa nie widzi liczby, tylko `bool`.
type TerminSlownika struct {
	ID           int64
	Kod          string
	Zrodlo       string
	Jezyk        string
	Cel          *string
	NieTlumaczyc bool
	Uwaga        *string
	// Stan i Dziedzina doszły z migracją 162 pod `translate.glossary.list`:
	// kontrakt zawęża wykaz terminów stanem (`GlossaryTermStatus`) i dziedziną.
	// Oba pola bywają puste — termin zastany nikogo o stan nie pytał, a nadanie
	// mu stanu domyślnego byłoby wydaniem zgody, której Operator nie wydał.
	Stan           *string
	Dziedzina      *string
	Zaktualizowano int64
}

const (
	kolumnyTerminuSlownika = `id, identyfikator_zewnetrzny, zrodlo, jezyk, cel,
	                          nie_tlumaczyc, uwaga, stan, dziedzina, zaktualizowano`

	zapiszTerminSlownika = `INSERT INTO termin_slownika
	                        (identyfikator_zewnetrzny, zrodlo, jezyk, cel,
	                         nie_tlumaczyc, uwaga, stan, dziedzina, zaktualizowano)
	                        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	                        ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                            zrodlo = excluded.zrodlo,
	                            jezyk = excluded.jezyk,
	                            cel = excluded.cel,
	                            nie_tlumaczyc = excluded.nie_tlumaczyc,
	                            uwaga = excluded.uwaga,
	                            stan = excluded.stan,
	                            dziedzina = excluded.dziedzina,
	                            zaktualizowano = excluded.zaktualizowano`

	pobierzTerminSlownika = `SELECT ` + kolumnyTerminuSlownika + ` FROM termin_slownika
	                         WHERE identyfikator_zewnetrzny = ?`

	listaTerminowSlownika = `SELECT ` + kolumnyTerminuSlownika + ` FROM termin_slownika
	                         ORDER BY jezyk, zrodlo`
)

// ZapiszTerminy zapisuje cały nadesłany wykaz terminów w jednej transakcji —
// `translate.glossary.set` przychodzi termin po terminie z warstwy wyżej, ale
// jeden termin nie ma prawa zostać zapisany, gdy kolejny w tym samym wywołaniu
// zawiedzie. Zwraca terminy po zapisie, odczytane z bazy, nie przepisane
// żądanie — `Zaktualizowano` i pola ustalone przez UPSERT mają wyjść ze stanu
// faktycznego.
func (r *repozytoriumTlumaczen) ZapiszTerminy(ctx context.Context,
	terminy []TerminSlownika) ([]TerminSlownika, error) {

	kody := make([]string, 0, len(terminy))
	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		zapis, err := r.zapytania.wTransakcji(ctx, transakcja, zapiszTerminSlownika)
		if err != nil {
			return err
		}
		teraz := time.Now().UnixMilli()
		for _, termin := range terminy {
			if termin.Kod == "" {
				return fmt.Errorf("dane: termin słownika bez identyfikatora")
			}
			if termin.Zrodlo == "" {
				return fmt.Errorf("dane: termin słownika %q bez treści źródłowej", termin.Kod)
			}
			if termin.Jezyk == "" {
				return fmt.Errorf("dane: termin słownika %q bez języka odpowiednika", termin.Kod)
			}
			zaktualizowano := termin.Zaktualizowano
			if zaktualizowano == 0 {
				zaktualizowano = teraz
			}
			_, err := zapis.ExecContext(ctx, termin.Kod, termin.Zrodlo, termin.Jezyk,
				tekstDoKolumny(termin.Cel), liczbaLogiczna(termin.NieTlumaczyc),
				tekstDoKolumny(termin.Uwaga), tekstDoKolumny(termin.Stan),
				tekstDoKolumny(termin.Dziedzina), zaktualizowano)
			if err != nil {
				return fmt.Errorf("dane: nie można zapisać terminu słownika %q: %w", termin.Kod, err)
			}
			kody = append(kody, termin.Kod)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	zapisane := make([]TerminSlownika, 0, len(kody))
	for _, kod := range kody {
		termin, err := r.Termin(ctx, kod)
		if err != nil {
			return nil, err
		}
		zapisane = append(zapisane, termin)
	}
	return zapisane, nil
}

// Terminy zwraca cały słownik uporządkowany po języku i treści źródłowej —
// obsługuje odczyt zasilający `glossary.apply` i `glossary.occurrences`,
// którym trzeba przejrzeć wszystkie terminy naraz.
func (r *repozytoriumTlumaczen) Terminy(ctx context.Context) ([]TerminSlownika, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaTerminowSlownika)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać terminów słownika: %w", err)
	}
	defer wiersze.Close()

	lista := []TerminSlownika{}
	for wiersze.Next() {
		termin, err := odczytajTerminSlownika(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz terminu słownika: %w", err)
		}
		lista = append(lista, termin)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt terminów słownika: %w", err)
	}
	return lista, nil
}

// Termin zwraca termin słownika o wskazanym kodzie. Brak wiersza wraca jako
// ErrBrakWiersza — warstwa wyższa odróżnia „nie ma” od „odczyt się nie
// powiódł”.
func (r *repozytoriumTlumaczen) Termin(ctx context.Context, kod string) (TerminSlownika, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzTerminSlownika)
	if err != nil {
		return TerminSlownika{}, err
	}
	termin, err := odczytajTerminSlownika(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return TerminSlownika{}, ErrBrakWiersza
	}
	if err != nil {
		return TerminSlownika{}, fmt.Errorf("dane: nieczytelny wiersz terminu słownika %q: %w", kod, err)
	}
	return termin, nil
}

// odczytajTerminSlownika składa strukturę z jednego wiersza wyniku.
func odczytajTerminSlownika(wiersz skaner) (TerminSlownika, error) {
	var termin TerminSlownika
	var cel, uwaga, stan, dziedzina sql.NullString
	var nieTlumaczyc int
	err := wiersz.Scan(&termin.ID, &termin.Kod, &termin.Zrodlo, &termin.Jezyk, &cel,
		&nieTlumaczyc, &uwaga, &stan, &dziedzina, &termin.Zaktualizowano)
	if err != nil {
		return TerminSlownika{}, err
	}
	termin.Cel = tekstZKolumny(cel)
	termin.NieTlumaczyc = nieTlumaczyc != 0
	termin.Uwaga = tekstZKolumny(uwaga)
	termin.Stan = tekstZKolumny(stan)
	termin.Dziedzina = tekstZKolumny(dziedzina)
	return termin, nil
}
