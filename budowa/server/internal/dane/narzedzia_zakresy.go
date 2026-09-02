// Odpowiedzialność pliku: zakresy uprawnień i limity wywołań pozycji katalogu
// narzędzi dla profilu asystenta (tabele `narzedzie_zakres_profilu` i `narzedzie_wywolanie` z migracji 278).
package dane

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type ZakresNarzedzia struct {
	ProfilKod       string
	NazwaPelna      string
	Dostepne        bool
	Potwierdzenie   bool
	LimitWywolan    int
	OknoSekund      int
	ZuzyteWywolania int
	Dopuszczone     []string
	Uzasadnienie    string
	Zaktualizowano  string
}

type RepozytoriumZakresowNarzedzi interface {
	ZakresyProfilu(ctx context.Context, kodProfilu, nazwaPelna string) ([]ZakresNarzedzia, error)
	ZuzycieProfilu(ctx context.Context, kodProfilu, kodSesji string) (map[string]int, error)
	ZapiszZakresNarzedzia(ctx context.Context, zakres ZakresNarzedzia) (ZakresNarzedzia, error)
	OdnotujWywolanieNarzedzia(ctx context.Context, kodProfilu, nazwaPelna, kodSesji string) error
}

const (
	kolumnyZakresuNarzedzia = `z.nazwa_pelna, z.dostepne, z.potwierdzenie, z.limit_wywolan,
	                           z.okno_sekund, z.dopuszczone, z.uzasadnienie, z.zaktualizowano`

	// Zakres i wywołanie nie niosą konta (migracja 484); granica idzie przez profil_asystenta.
	zakresyNarzedziProfilu = `SELECT ` + kolumnyZakresuNarzedzia + `
	                            FROM narzedzie_zakres_profilu z
	                            JOIN profil_asystenta p ON p.id = z.profil_id
	                           WHERE p.identyfikator_zewnetrzny = ?
	                             AND (? = '' OR z.nazwa_pelna = ?)
	                             AND ` + WarunekKonta + `
	                           ORDER BY z.nazwa_pelna`

	zapiszZakresNarzedziaProfilu = `INSERT INTO narzedzie_zakres_profilu
	    (profil_id, nazwa_pelna, dostepne, potwierdzenie, limit_wywolan, okno_sekund,
	     dopuszczone, uzasadnienie, zaktualizowano)
	    VALUES (?, ?, ?, ?, ?, ?, ?, ?, strftime('%Y-%m-%dT%H:%M:%fZ','now'))
	    ON CONFLICT(profil_id, nazwa_pelna) DO UPDATE SET
	        dostepne = excluded.dostepne,
	        potwierdzenie = excluded.potwierdzenie,
	        limit_wywolan = excluded.limit_wywolan,
	        okno_sekund = excluded.okno_sekund,
	        dopuszczone = excluded.dopuszczone,
	        uzasadnienie = excluded.uzasadnienie,
	        zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')`

	zuzycieZakresowNarzedzi = `SELECT z.nazwa_pelna, COUNT(w.id)
	                             FROM narzedzie_zakres_profilu z
	                             JOIN profil_asystenta p ON p.id = z.profil_id
	                             LEFT JOIN narzedzie_wywolanie w
	                                    ON w.profil_id = z.profil_id
	                                   AND w.nazwa_pelna = z.nazwa_pelna
	                                   AND w.wykonano >= strftime('%Y-%m-%dT%H:%M:%fZ','now',
	                                                              '-' || z.okno_sekund || ' seconds')
	                                   AND (? = '' OR w.sesja_kod = ?)
	                            WHERE p.identyfikator_zewnetrzny = ?
	                              AND ` + WarunekKonta + `
	                            GROUP BY z.nazwa_pelna`

	wstawWywolanieNarzedzia = `INSERT INTO narzedzie_wywolanie (profil_id, nazwa_pelna, sesja_kod)
	                           VALUES (?, ?, ?)`

	sprzatnijWywolaniaNarzedzia = `DELETE FROM narzedzie_wywolanie
	                                WHERE profil_id = ? AND nazwa_pelna = ?
	                                  AND wykonano < strftime('%Y-%m-%dT%H:%M:%fZ','now', '-' ||
	                                      (SELECT COALESCE(MAX(okno_sekund), 3600)
	                                         FROM narzedzie_zakres_profilu
	                                        WHERE profil_id = ? AND nazwa_pelna = ?) || ' seconds')`

	numerProfiluAsystenta = `SELECT id FROM profil_asystenta WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta
)

type repozytoriumZakresowNarzedzi struct {
	zapytania *zapytania
	db        *sql.DB
}

var _ RepozytoriumZakresowNarzedzi = (*repozytoriumZakresowNarzedzi)(nil)

func noweRepozytoriumZakresowNarzedzi(z *zapytania, db *sql.DB) *repozytoriumZakresowNarzedzi {
	return &repozytoriumZakresowNarzedzi{zapytania: z, db: db}
}

func (r *repozytoriumZakresowNarzedzi) numerProfiluZakresu(ctx context.Context, kod string) (int64, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, numerProfiluAsystenta)
	if err != nil {
		return 0, err
	}
	var id int64
	switch err := polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)).Scan(&id); {
	case err == sql.ErrNoRows:
		return 0, ErrBrakWiersza
	case err != nil:
		return 0, fmt.Errorf("dane: nie można odczytać profilu asystenta %q: %w", kod, err)
	}
	return id, nil
}

// ZakresyProfilu oddaje zakresy zapisane dla profilu; nazwa pusta zwraca komplet zakresów.
func (r *repozytoriumZakresowNarzedzi) ZakresyProfilu(ctx context.Context,
	kodProfilu, nazwaPelna string) ([]ZakresNarzedzia, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, zakresyNarzedziProfilu)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, kodProfilu, nazwaPelna, nazwaPelna, KontoOperatora(ctx))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać zakresów narzędzi profilu %q: %w", kodProfilu, err)
	}
	defer wiersze.Close()
	zakresy := make([]ZakresNarzedzia, 0, 16)
	for wiersze.Next() {
		zakres := ZakresNarzedzia{ProfilKod: kodProfilu}
		var dostepne, potwierdzenie int
		var dopuszczone string
		if err := wiersze.Scan(&zakres.NazwaPelna, &dostepne, &potwierdzenie, &zakres.LimitWywolan,
			&zakres.OknoSekund, &dopuszczone, &zakres.Uzasadnienie, &zakres.Zaktualizowano); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny zakres narzędzia profilu %q: %w", kodProfilu, err)
		}
		zakres.Dostepne = dostepne == 1
		zakres.Potwierdzenie = potwierdzenie == 1
		zakres.Dopuszczone = wartosciDopuszczoneZakresu(dopuszczone)
		zakresy = append(zakresy, zakres)
	}
	if err := wiersze.Err(); err != nil {
		return nil, err
	}
	return zakresy, nil
}

// ZuzycieProfilu oddaje liczbę wywołań w bieżącym oknie czasu; kod sesji pusty liczy wszystkie karty.
func (r *repozytoriumZakresowNarzedzi) ZuzycieProfilu(ctx context.Context,
	kodProfilu, kodSesji string) (map[string]int, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, zuzycieZakresowNarzedzi)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, kodSesji, kodSesji, kodProfilu, KontoOperatora(ctx))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można policzyć zużycia limitów profilu %q: %w", kodProfilu, err)
	}
	defer wiersze.Close()
	zuzycie := make(map[string]int, 16)
	for wiersze.Next() {
		var nazwa string
		var liczba int
		if err := wiersze.Scan(&nazwa, &liczba); err != nil {
			return nil, fmt.Errorf("dane: nieczytelne zużycie limitu profilu %q: %w", kodProfilu, err)
		}
		zuzycie[nazwa] = liczba
	}
	return zuzycie, wiersze.Err()
}

// ZapiszZakresNarzedzia zakłada albo zmienia zakres i oddaje go po zapisie.
func (r *repozytoriumZakresowNarzedzi) ZapiszZakresNarzedzia(ctx context.Context,
	zakres ZakresNarzedzia) (ZakresNarzedzia, error) {

	id, err := r.numerProfiluZakresu(ctx, zakres.ProfilKod)
	if err != nil {
		return ZakresNarzedzia{}, err
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszZakresNarzedziaProfilu)
	if err != nil {
		return ZakresNarzedzia{}, err
	}
	_, err = polecenie.ExecContext(ctx, id, zakres.NazwaPelna, liczbaLogiczna(zakres.Dostepne),
		liczbaLogiczna(zakres.Potwierdzenie), zakres.LimitWywolan, zakres.OknoSekund,
		strings.Join(zakres.Dopuszczone, "\n"), zakres.Uzasadnienie)
	if err != nil {
		return ZakresNarzedzia{}, fmt.Errorf("dane: nie można zapisać zakresu %q profilu %q: %w",
			zakres.NazwaPelna, zakres.ProfilKod, err)
	}
	zapisane, err := r.ZakresyProfilu(ctx, zakres.ProfilKod, zakres.NazwaPelna)
	if err != nil {
		return ZakresNarzedzia{}, err
	}
	if len(zapisane) == 0 {
		return ZakresNarzedzia{}, ErrBrakWiersza
	}
	return zapisane[0], nil
}

// OdnotujWywolanieNarzedzia dopisuje jedno wywołanie do rachunku limitu i sprząta wiersze spoza okna.
func (r *repozytoriumZakresowNarzedzi) OdnotujWywolanieNarzedzia(ctx context.Context,
	kodProfilu, nazwaPelna, kodSesji string) error {

	id, err := r.numerProfiluZakresu(ctx, kodProfilu)
	if err != nil {
		return err
	}
	polecenie, err := r.zapytania.przygotuj(ctx, wstawWywolanieNarzedzia)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, id, nazwaPelna, tekstDoKolumny(wskaznikTekstuZakresu(kodSesji))); err != nil {
		return fmt.Errorf("dane: nie można odnotować wywołania %q: %w", nazwaPelna, err)
	}
	sprzataczka, err := r.zapytania.przygotuj(ctx, sprzatnijWywolaniaNarzedzia)
	if err != nil {
		return err
	}
	if _, err := sprzataczka.ExecContext(ctx, id, nazwaPelna, id, nazwaPelna); err != nil {
		return fmt.Errorf("dane: nie można sprzątnąć rachunku wywołań %q: %w", nazwaPelna, err)
	}
	return nil
}

func wartosciDopuszczoneZakresu(tresc string) []string {
	przyciety := strings.TrimSpace(tresc)
	if przyciety == "" {
		return nil
	}
	pozycje := strings.Split(przyciety, "\n")
	wynik := make([]string, 0, len(pozycje))
	for _, pozycja := range pozycje {
		if wartosc := strings.TrimSpace(pozycja); wartosc != "" {
			wynik = append(wynik, wartosc)
		}
	}
	return wynik
}

// Kolumna `sesja_kod` niesie NULL dla wywołania spoza karty sesji (migracja 278).
func wskaznikTekstuZakresu(wartosc string) *string {
	if strings.TrimSpace(wartosc) == "" {
		return nil
	}
	return &wartosc
}
