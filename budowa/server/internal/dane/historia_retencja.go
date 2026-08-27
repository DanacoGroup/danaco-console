// Odpowiedzialność pliku: zasada przechowywania historii (tabela
// `zasada_przechowywania`) — jej zapis, rozstrzygnięcie dla okna i egzekucja na wykazie pozycji.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// ZasadaPrzechowywania to wiersz tabeli `zasada_przechowywania`. Oba progi puste
// znaczą zasadę wyłączoną, nie zasadę zerową.
type ZasadaPrzechowywania struct {
	// Zakres niesie wartość okna, sesji albo globalną; ZakresKod bywa pusty tylko przy zakresie globalnym.
	Zakres          string
	ZakresKod       string
	DniTrzymania    *int
	PozycjeTrzymane *int
}

const (
	// Granicę czasu wylicza rdzeń, nie `now` bazy: dwa zegary to dwie prawdy.
	// Retencja liczbą pozycji zostawia N najnowszych.
	usunHistorieStarsza    = usunHistorieOkna + ` AND utworzono < ?`
	usunHistorieNadmiarowa = usunHistorieOkna + ` AND id NOT IN (
	                             SELECT w.id FROM wiadomosc w
	                             JOIN okno_komunikacji o ON o.id = w.okno_komunikacji_id
	                             WHERE o.identyfikator_zewnetrzny = ?
	                             ORDER BY w.utworzono DESC, w.id DESC
	                             LIMIT ?)`

	// Zasada bez progów nie jest wierszem, tylko jego brakiem: nastawa wyłączona usuwa wiersz zamiast zapisywać
	// go pustym. Kod pusty (zakres global) porównuje się przez funkcję COALESCE, bo kolumna trzyma wtedy wartość NULL.
	usunZasadePrzechowywania = `DELETE FROM zasada_przechowywania
	                            WHERE zakres = ? AND COALESCE(zakres_kod, '') = ?`

	zapiszZasadePrzechowywania = `INSERT INTO zasada_przechowywania
	                              (zakres, zakres_kod, dni_trzymania, pozycje_trzymane)
	                              VALUES (?, ?, ?, ?)
	                              ON CONFLICT (zakres, COALESCE(zakres_kod, '')) DO UPDATE SET
	                                  dni_trzymania = excluded.dni_trzymania,
	                                  pozycje_trzymane = excluded.pozycje_trzymane,
	                                  zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')`

	// Zakres rozstrzyga się od najwęższego do najszerszego: najpierw okno, potem sesja, na końcu zakres globalny.
	zasadaOkna = `SELECT zakres, COALESCE(zakres_kod, ''), dni_trzymania, pozycje_trzymane
	              FROM zasada_przechowywania
	              WHERE (zakres = 'window' AND zakres_kod = ?)
	                 OR (zakres = 'session' AND zakres_kod = (
	                         SELECT s.identyfikator_zewnetrzny
	                         FROM okno_komunikacji o JOIN sesja s ON s.id = o.sesja_id
	                         WHERE o.identyfikator_zewnetrzny = ?))
	                 OR zakres = 'global'
	              ORDER BY CASE zakres WHEN 'window' THEN 0 WHEN 'session' THEN 1 ELSE 2 END
	              LIMIT 1`
	oknaZakresuSesji = `SELECT o.identyfikator_zewnetrzny
	                    FROM okno_komunikacji o JOIN sesja s ON s.id = o.sesja_id
	                    WHERE s.identyfikator_zewnetrzny = ?
	                      AND o.identyfikator_zewnetrzny IS NOT NULL`
	oknaZakresuGlobalnego = `SELECT identyfikator_zewnetrzny FROM okno_komunikacji
	                         WHERE identyfikator_zewnetrzny IS NOT NULL`
	// Byt zakresu — okno albo sesja — sprawdzany jest przed zapisaniem zasady, żeby nastawa nie odnosiła się
	// do bytu, którego nie ma.
	istnienieOkna  = `SELECT 1 FROM okno_komunikacji WHERE identyfikator_zewnetrzny = ? LIMIT 1`
	istnienieSesji = `SELECT 1 FROM sesja WHERE identyfikator_zewnetrzny = ? LIMIT 1`
)

// ZapiszZasade zakłada zasadę zakresu albo nadpisuje istniejącą; zasada bez obu progów zdejmuje wiersz zakresu,
// zamiast zapisywać wiersz pusty.
func (r *repozytoriumHistorii) ZapiszZasade(ctx context.Context,
	zasada ZasadaPrzechowywania) (ZasadaPrzechowywania, error) {

	if zasada.DniTrzymania == nil && zasada.PozycjeTrzymane == nil {
		return r.zdejmijZasade(ctx, zasada)
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszZasadePrzechowywania)
	if err != nil {
		return ZasadaPrzechowywania{}, err
	}
	_, err = polecenie.ExecContext(ctx, zasada.Zakres, kodZakresuDoKolumny(zasada),
		progDoKolumny(zasada.DniTrzymania), progDoKolumny(zasada.PozycjeTrzymane))
	if err != nil {
		return ZasadaPrzechowywania{}, fmt.Errorf(
			"dane: nie można zapisać zasady przechowywania zakresu %q: %w", zasada.Zakres, err)
	}
	return zasada, nil
}

// zdejmijZasade kasuje wiersz zasady zakresu. Brak wiersza nie jest błędem:
// zdjęcie zasady, której nie było, zostawia stan taki, jakiego Operator chciał,
// więc czynność jest idempotentna.
func (r *repozytoriumHistorii) zdejmijZasade(ctx context.Context,
	zasada ZasadaPrzechowywania) (ZasadaPrzechowywania, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, usunZasadePrzechowywania)
	if err != nil {
		return ZasadaPrzechowywania{}, err
	}
	if _, err := polecenie.ExecContext(ctx, zasada.Zakres, zasada.ZakresKod); err != nil {
		return ZasadaPrzechowywania{}, fmt.Errorf(
			"dane: nie można zdjąć zasady przechowywania zakresu %q: %w", zasada.Zakres, err)
	}
	return zasada, nil
}

// ZasadaOkna rozstrzyga zasadę przechowywania obowiązującą wskazane okno komunikacji. Brak zasady nie jest błędem.
func (r *repozytoriumHistorii) ZasadaOkna(ctx context.Context,
	oknoKod string) (ZasadaPrzechowywania, bool, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, zasadaOkna)
	if err != nil {
		return ZasadaPrzechowywania{}, false, err
	}
	var zasada ZasadaPrzechowywania
	var dni, pozycje *int64
	err = polecenie.QueryRowContext(ctx, oknoKod, oknoKod).
		Scan(&zasada.Zakres, &zasada.ZakresKod, &dni, &pozycje)
	if errors.Is(err, sql.ErrNoRows) {
		return ZasadaPrzechowywania{}, false, nil
	}
	if err != nil {
		return ZasadaPrzechowywania{}, false,
			fmt.Errorf("dane: nie można rozstrzygnąć zasady przechowywania okna %q: %w", oknoKod, err)
	}
	zasada.DniTrzymania = progZKolumny(dni)
	zasada.PozycjeTrzymane = progZKolumny(pozycje)
	return zasada, true, nil
}

// IstniejeByt rozstrzyga, czy byt wskazany przez zasadę zakresu w ogóle istnieje w bazie danych rdzenia.
func (r *repozytoriumHistorii) IstniejeByt(ctx context.Context, zakres, zakresKod string) (bool, error) {
	zapytanie := istnienieOkna
	switch zakres {
	case "global":
		return true, nil
	case "session":
		zapytanie = istnienieSesji
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return false, err
	}
	var jeden int
	err = polecenie.QueryRowContext(ctx, zakresKod).Scan(&jeden)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("dane: nie można sprawdzić bytu zakresu %q %q: %w", zakres, zakresKod, err)
	}
	return true, nil
}

// OknaZakresu wylicza okna objęte zasadą zakresu. Zakres okna nie pyta bazy —
// oknem zakresu jest samo wskazanie; egzekucja na oknie nieistniejącym niczego
// nie usuwa i nie jest błędem.
func (r *repozytoriumHistorii) OknaZakresu(ctx context.Context, zakres, zakresKod string) ([]string, error) {
	if zakres == "window" {
		return []string{zakresKod}, nil
	}
	zapytanie, argumenty := oknaZakresuGlobalnego, []any{}
	if zakres == "session" {
		zapytanie, argumenty = oknaZakresuSesji, []any{zakresKod}
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, argumenty...)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można wyliczyć okien zakresu %q: %w", zakres, err)
	}
	defer wiersze.Close()

	okna := []string{}
	for wiersze.Next() {
		var kod string
		if err := wiersze.Scan(&kod); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz okna zakresu %q: %w", zakres, err)
		}
		okna = append(okna, kod)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwane wyliczanie okien zakresu %q: %w", zakres, err)
	}
	return okna, nil
}

// Egzekwuj stosuje zasadę przechowywania okna; zasada bez progów niczego nie usuwa — jest wyłączona, nie zerowa.
func (r *repozytoriumHistorii) Egzekwuj(ctx context.Context, oknoKod string) (int, error) {
	zasada, jest, err := r.ZasadaOkna(ctx, oknoKod)
	if err != nil || !jest {
		return 0, err
	}
	if zasada.DniTrzymania == nil && zasada.PozycjeTrzymane == nil {
		return 0, nil
	}
	usuniete := 0
	err = wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		if zasada.DniTrzymania != nil {
			granica := time.Now().UTC().AddDate(0, 0, -*zasada.DniTrzymania).Format(znacznikCzasuHistorii)
			liczba, err := wykonajUsuniecie(ctx, transakcja, usunHistorieStarsza, oknoKod, granica)
			if err != nil {
				return err
			}
			usuniete += liczba
		}
		if zasada.PozycjeTrzymane != nil {
			liczba, err := wykonajUsuniecie(ctx, transakcja, usunHistorieNadmiarowa,
				oknoKod, oknoKod, *zasada.PozycjeTrzymane)
			if err != nil {
				return err
			}
			usuniete += liczba
		}
		if usuniete == 0 {
			return nil
		}
		return sprzatnijBloki(ctx, transakcja, oknoKod)
	})
	if err != nil {
		return 0, err
	}
	return usuniete, nil
}

// kodZakresuDoKolumny znosi kod pusty do NULL — zakres global bytu nie
// wskazuje, a warunek CHECK pustego napisu tam nie wpuści.
func kodZakresuDoKolumny(zasada ZasadaPrzechowywania) any {
	if zasada.ZakresKod == "" {
		return nil
	}
	return zasada.ZakresKod
}

// progDoKolumny znosi próg nieustawiony do wartości pustej kolumny: brak znaczy „bez ograniczenia", nie zero.
func progDoKolumny(prog *int) any {
	if prog == nil {
		return nil
	}
	return *prog
}

// progZKolumny podnosi wartość kolumny bazy do progu kontraktu; wartość pusta kolumny zostaje brakiem progu.
func progZKolumny(kolumna *int64) *int {
	if kolumna == nil {
		return nil
	}
	prog := int(*kolumna)
	return &prog
}
