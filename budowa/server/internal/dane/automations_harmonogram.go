// Odpowiedzialność pliku: harmonogram automatyki i jego wyzwalacze (tabele
// `harmonogram_automatyki`, `wyzwalacz_automatyki`) — trwałość okna Scheduler.
//
// Harmonogram i wyzwalacze zapisują się razem. Wyzwalacz bez harmonogramu nie
// ma czego wyzwalać, a harmonogram zapisany bez wyzwalaczy zostawiłby w bazie
// wyzwalacze poprzedniej wersji. Jedna transakcja zamyka obie możliwości.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// Harmonogram to wiersz tabeli `harmonogram_automatyki`.
type Harmonogram struct {
	ID                   int64
	Kod                  string
	AutomatykaID         int64
	Cron                 *string
	StrefaCzasowa        *string
	Czynny               bool
	NastepneUruchomienie *string
	Utworzono            string
	Zaktualizowano       string
	// Cztery pola nadzoru i wejścia zdalnego (migracja 275). Trzymane w bazie,
	// bo mają obowiązywać po ponownym złożeniu rdzenia: okno tolerancji
	// w pamięci procesu przestałoby nadzorować cokolwiek po pierwszym restarcie.
	// OdwolaniePodpisu niesie REFERENCJĘ klucza HMAC w sejfie, nigdy wartość.
	TolerancjaSekundy   int
	RegulaNadzoru       *string
	OdwolaniePodpisu    *string
	OknoDeduplikacjiSek int
}

// WyzwalaczAutomatyki to wiersz tabeli `wyzwalacz_automatyki`.
type WyzwalaczAutomatyki struct {
	ID        int64
	Kod       string
	Rodzaj    string
	Wyrazenie string
	Czynny    bool
	Kolejnosc int
}

const (
	kolumnyHarmonogramu = `id, identyfikator_zewnetrzny, automatyka_id, cron, strefa_czasowa,
	                       czynny, nastepne_uruchomienie, utworzono, zaktualizowano,
	                       tolerancja_sekundy, regula_nadzoru, odwolanie_podpisu,
	                       okno_deduplikacji_sekundy`

	zapiszHarmonogram = `INSERT INTO harmonogram_automatyki
	                     (identyfikator_zewnetrzny, automatyka_id, cron, strefa_czasowa,
	                      czynny, nastepne_uruchomienie)
	                     VALUES (?, ?, ?, ?, ?, ?)
	                     ON CONFLICT(automatyka_id) DO UPDATE SET
	                         cron = excluded.cron,
	                         strefa_czasowa = excluded.strefa_czasowa,
	                         czynny = excluded.czynny,
	                         nastepne_uruchomienie = excluded.nastepne_uruchomienie,
	                         zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')`

	pobierzHarmonogram = `SELECT ` + kolumnyHarmonogramu + ` FROM harmonogram_automatyki
	                      WHERE automatyka_id = ?`

	usunWyzwalacze = `DELETE FROM wyzwalacz_automatyki WHERE harmonogram_id = ?`

	wstawWyzwalacz = `INSERT INTO wyzwalacz_automatyki
	                  (harmonogram_id, identyfikator_zewnetrzny, rodzaj, wyrazenie, czynny, kolejnosc)
	                  VALUES (?, ?, ?, ?, ?, ?)`

	listaWyzwalaczy = `SELECT id, identyfikator_zewnetrzny, rodzaj, wyrazenie, czynny, kolejnosc
	                   FROM wyzwalacz_automatyki WHERE harmonogram_id = ?
	                   ORDER BY kolejnosc, id`

	// Budzik pyta o harmonogramy czynne, których wyliczony termin już minął.
	// Termin nadany (`nastepne_uruchomienie IS NOT NULL`), bo harmonogram bez
	// zrozumiałej cykliczności terminu nie ma i odpalić się nie może.
	naleznHarmonogramy = `SELECT ` + kolumnyHarmonogramu + ` FROM harmonogram_automatyki
	                      WHERE czynny = 1 AND nastepne_uruchomienie IS NOT NULL
	                        AND nastepne_uruchomienie <= ?
	                      ORDER BY nastepne_uruchomienie, id`

	ustawNastepneUruchomienie = `UPDATE harmonogram_automatyki
	                             SET nastepne_uruchomienie = ?,
	                                 zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                             WHERE automatyka_id = ?`
)

// ZapiszHarmonogram zapisuje harmonogram wraz z kompletem wyzwalaczy i zwraca
// stan po zapisie.
func (r *repozytoriumAutomatyk) ZapiszHarmonogram(ctx context.Context, harmonogram Harmonogram,
	wyzwalacze []WyzwalaczAutomatyki) (Harmonogram, error) {

	if harmonogram.Kod == "" || harmonogram.AutomatykaID == 0 {
		return Harmonogram{}, fmt.Errorf("dane: harmonogram bez identyfikatora albo bez automatyki")
	}
	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		zapis, err := r.zapytania.wTransakcji(ctx, transakcja, zapiszHarmonogram)
		if err != nil {
			return err
		}
		_, err = zapis.ExecContext(ctx, harmonogram.Kod, harmonogram.AutomatykaID,
			tekstDoKolumny(harmonogram.Cron), tekstDoKolumny(harmonogram.StrefaCzasowa),
			liczbaLogiczna(harmonogram.Czynny), tekstDoKolumny(harmonogram.NastepneUruchomienie))
		if err != nil {
			return fmt.Errorf("dane: nie można zapisać harmonogramu automatyki %d: %w",
				harmonogram.AutomatykaID, err)
		}
		zapisany, err := harmonogramWTransakcji(ctx, r.zapytania, transakcja, harmonogram.AutomatykaID)
		if err != nil {
			return err
		}
		return zapiszWyzwalacze(ctx, r.zapytania, transakcja, zapisany.ID, wyzwalacze)
	})
	if err != nil {
		return Harmonogram{}, err
	}
	return r.Harmonogram(ctx, harmonogram.AutomatykaID)
}

// Harmonogram zwraca harmonogram automatyki. Brak wiersza wraca jako
// ErrBrakWiersza — automatyka bez harmonogramu jest stanem poprawnym.
func (r *repozytoriumAutomatyk) Harmonogram(ctx context.Context, automatykaID int64) (Harmonogram, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzHarmonogram)
	if err != nil {
		return Harmonogram{}, err
	}
	harmonogram, err := odczytajHarmonogram(polecenie.QueryRowContext(ctx, automatykaID))
	if errors.Is(err, sql.ErrNoRows) {
		return Harmonogram{}, ErrBrakWiersza
	}
	if err != nil {
		return Harmonogram{}, fmt.Errorf("dane: nieczytelny wiersz harmonogramu automatyki %d: %w",
			automatykaID, err)
	}
	return harmonogram, nil
}

// Wyzwalacze zwraca wyzwalacze harmonogramu w zapisanej kolejności.
func (r *repozytoriumAutomatyk) Wyzwalacze(ctx context.Context,
	harmonogramID int64) ([]WyzwalaczAutomatyki, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaWyzwalaczy)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, harmonogramID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać wyzwalaczy harmonogramu %d: %w", harmonogramID, err)
	}
	defer wiersze.Close()

	lista := []WyzwalaczAutomatyki{}
	for wiersze.Next() {
		var wyzwalacz WyzwalaczAutomatyki
		var czynny int
		err := wiersze.Scan(&wyzwalacz.ID, &wyzwalacz.Kod, &wyzwalacz.Rodzaj,
			&wyzwalacz.Wyrazenie, &czynny, &wyzwalacz.Kolejnosc)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz wyzwalacza: %w", err)
		}
		wyzwalacz.Czynny = czynny == 1
		lista = append(lista, wyzwalacz)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt wyzwalaczy harmonogramu %d: %w", harmonogramID, err)
	}
	return lista, nil
}

// HarmonogramyNalezne zwraca harmonogramy czynne, których termin najbliższego
// uruchomienia nie jest późniejszy niż `teraz` — to zbiór, który budzik ma
// odpalić. `teraz` jest znacznikiem bazy (ten sam format co kolumna terminu).
func (r *repozytoriumAutomatyk) HarmonogramyNalezne(ctx context.Context, teraz string) ([]Harmonogram, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, naleznHarmonogramy)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, teraz)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać należnych harmonogramów: %w", err)
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
		return nil, fmt.Errorf("dane: przerwany odczyt należnych harmonogramów: %w", err)
	}
	return lista, nil
}

// UstawNastepneUruchomienie przesuwa termin najbliższego uruchomienia. Budzik
// robi to przed odpaleniem, żeby ten sam harmonogram nie ruszył ponownie, gdyby
// odpalenie trwało dłużej niż takt zegara. Termin `nil` znaczy brak następnego
// terminu (harmonogram wyłączony albo cykliczność niezrozumiała).
func (r *repozytoriumAutomatyk) UstawNastepneUruchomienie(ctx context.Context,
	automatykaID int64, nastepne *string) error {

	polecenie, err := r.zapytania.przygotuj(ctx, ustawNastepneUruchomienie)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, tekstDoKolumny(nastepne), automatykaID); err != nil {
		return fmt.Errorf("dane: nie można przesunąć terminu harmonogramu automatyki %d: %w",
			automatykaID, err)
	}
	return nil
}

// zapiszWyzwalacze podmienia komplet wyzwalaczy harmonogramu w transakcji.
func zapiszWyzwalacze(ctx context.Context, z *zapytania, transakcja *sql.Tx,
	harmonogramID int64, wyzwalacze []WyzwalaczAutomatyki) error {

	czyszczenie, err := z.wTransakcji(ctx, transakcja, usunWyzwalacze)
	if err != nil {
		return err
	}
	if _, err := czyszczenie.ExecContext(ctx, harmonogramID); err != nil {
		return fmt.Errorf("dane: nie można wyczyścić wyzwalaczy harmonogramu %d: %w", harmonogramID, err)
	}
	wstawienie, err := z.wTransakcji(ctx, transakcja, wstawWyzwalacz)
	if err != nil {
		return err
	}
	for numer, wyzwalacz := range wyzwalacze {
		kolejnosc := wyzwalacz.Kolejnosc
		if kolejnosc == 0 {
			kolejnosc = numer + 1
		}
		_, err := wstawienie.ExecContext(ctx, harmonogramID, wyzwalacz.Kod, wyzwalacz.Rodzaj,
			wyzwalacz.Wyrazenie, liczbaLogiczna(wyzwalacz.Czynny), kolejnosc)
		if err != nil {
			return fmt.Errorf("dane: nie można zapisać wyzwalacza %q: %w", wyzwalacz.Kod, err)
		}
	}
	return nil
}

// harmonogramWTransakcji odczytuje wiersz harmonogramu wewnątrz transakcji
// zapisu — identyfikator jest potrzebny wyzwalaczom, zanim transakcja się
// zamknie.
func harmonogramWTransakcji(ctx context.Context, z *zapytania, transakcja *sql.Tx,
	automatykaID int64) (Harmonogram, error) {

	polecenie, err := z.wTransakcji(ctx, transakcja, pobierzHarmonogram)
	if err != nil {
		return Harmonogram{}, err
	}
	harmonogram, err := odczytajHarmonogram(polecenie.QueryRowContext(ctx, automatykaID))
	if err != nil {
		return Harmonogram{}, fmt.Errorf("dane: nieczytelny harmonogram automatyki %d: %w", automatykaID, err)
	}
	return harmonogram, nil
}

// odczytajHarmonogram składa strukturę z jednego wiersza wyniku.
func odczytajHarmonogram(wiersz skaner) (Harmonogram, error) {
	var harmonogram Harmonogram
	var cron, strefa, nastepne, regula, podpis sql.NullString
	var czynny int
	err := wiersz.Scan(&harmonogram.ID, &harmonogram.Kod, &harmonogram.AutomatykaID,
		&cron, &strefa, &czynny, &nastepne, &harmonogram.Utworzono, &harmonogram.Zaktualizowano,
		&harmonogram.TolerancjaSekundy, &regula, &podpis, &harmonogram.OknoDeduplikacjiSek)
	if err != nil {
		return Harmonogram{}, err
	}
	harmonogram.Cron, harmonogram.StrefaCzasowa = tekstZKolumny(cron), tekstZKolumny(strefa)
	harmonogram.NastepneUruchomienie = tekstZKolumny(nastepne)
	harmonogram.RegulaNadzoru, harmonogram.OdwolaniePodpisu = tekstZKolumny(regula), tekstZKolumny(podpis)
	harmonogram.Czynny = czynny == 1
	return harmonogram, nil
}
