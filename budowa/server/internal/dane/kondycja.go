// Magazyn sond kondycji i serii ich wyników (sonda_kondycji, wynik_sondy_kondycji) dla rodziny poleceń health.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// SondaKondycji to definicja jednego pomiaru wraz z odbiciem ostatniego przebiegu.
type SondaKondycji struct {
	ID              int64
	Kod             string
	Nazwa           string
	Rodzaj          string
	Cel             string
	KomponentKod    *string
	OdstepMs        int64
	LimitCzasuMs    *int64
	OczekiwanyKod   *int64
	TrescWysylana   *string
	CelDostepnosci  *float64
	Czynna          bool
	OstatniStan     *string
	OstatniPrzebieg *int64
	Utworzono       int64
	Zaktualizowano  *int64
}

// WynikSondyKondycji to jeden niezmienny pomiar wykonany o znanej godzinie.
type WynikSondyKondycji struct {
	Kod              string
	SondaKod         string
	Stan             string
	Wykonano         int64
	CzasOdpowiedziMs *int64
	StatusHttp       *int64
	Szczegol         *string
	BladKod          *string
}

// SitoWynikowKondycji zawęża odczyt serii: sondą, stanem, przedziałem czasu i granicą wierszy.
type SitoWynikowKondycji struct {
	SondaKod string
	Stan     string
	OdCzasu  int64
	DoCzasu  int64
	Granica  int
}

// RepozytoriumKondycji jest kontraktem magazynu sond, serii wyników i pulsu bazy rdzenia.
type RepozytoriumKondycji interface {
	ZapiszSonde(ctx context.Context, sonda SondaKondycji) (SondaKondycji, bool, error)
	Sonda(ctx context.Context, kod string) (SondaKondycji, error)
	Sondy(ctx context.Context, rodzaj, komponent string, tylkoCzynne bool, granica int) ([]SondaKondycji, error)
	UsunSonde(ctx context.Context, kod string) (int, error)
	ZapiszWynik(ctx context.Context, wynik WynikSondyKondycji) (WynikSondyKondycji, error)
	Wyniki(ctx context.Context, sito SitoWynikowKondycji) ([]WynikSondyKondycji, int, error)
	Puls(ctx context.Context) error
}

const (
	kolumnySondyKondycji = `s.id, s.identyfikator_zewnetrzny, s.nazwa, s.rodzaj, s.cel,
	                        s.komponent_kod, s.odstep_ms, s.limit_czasu_ms, s.oczekiwany_status,
	                        s.tresc_wysylana, s.cel_dostepnosci, s.czynna, s.ostatni_stan,
	                        s.ostatni_przebieg, s.utworzono, s.zaktualizowano`

	wstawSondeKondycji = `INSERT INTO sonda_kondycji
	                      (identyfikator_zewnetrzny, nazwa, rodzaj, cel, komponent_kod, odstep_ms,
	                       limit_czasu_ms, oczekiwany_status, tresc_wysylana, cel_dostepnosci,
	                       czynna, utworzono, zaktualizowano, konto_id)
	                      VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ` + WskazanieKonta + `)`

	aktualizujSondeKondycji = `UPDATE sonda_kondycji SET
	                               nazwa = ?, rodzaj = ?, cel = ?, komponent_kod = ?, odstep_ms = ?,
	                               limit_czasu_ms = ?, oczekiwany_status = ?, tresc_wysylana = ?,
	                               cel_dostepnosci = ?, czynna = ?, zaktualizowano = ?
	                           WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	pobierzSondeKondycji = `SELECT ` + kolumnySondyKondycji +
		` FROM sonda_kondycji s WHERE s.identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	usunSondeKondycji = `DELETE FROM sonda_kondycji WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	policzWynikiSondyKondycji = `SELECT COUNT(*) FROM wynik_sondy_kondycji WHERE sonda_kod = ?`

	kolumnyWynikuKondycji = `identyfikator_zewnetrzny, sonda_kod, stan, wykonano,
	                         czas_odpowiedzi_ms, status_http, szczegol, blad_kod`

	wstawWynikKondycji = `INSERT INTO wynik_sondy_kondycji (` + kolumnyWynikuKondycji + `)
	                      VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	odbijPrzebiegSondyKondycji = `UPDATE sonda_kondycji
	                              SET ostatni_stan = ?, ostatni_przebieg = ?
	                              WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta
)

type repozytoriumKondycji struct {
	zapytania *zapytania
	db        *sql.DB
}

func noweRepozytoriumKondycji(z *zapytania, db *sql.DB) *repozytoriumKondycji {
	return &repozytoriumKondycji{zapytania: z, db: db}
}

// ZapiszSonde zakłada definicję albo nadpisuje zastaną po kodzie; drugi wynik mówi, czy sonda powstała teraz.
func (r *repozytoriumKondycji) ZapiszSonde(ctx context.Context,
	sonda SondaKondycji) (SondaKondycji, bool, error) {

	if strings.TrimSpace(sonda.Kod) == "" {
		return SondaKondycji{}, false, fmt.Errorf("dane: sonda kondycji bez identyfikatora")
	}
	powstala := false
	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		odczyt, err := r.zapytania.wTransakcji(ctx, transakcja, pobierzSondeKondycji)
		if err != nil {
			return err
		}
		zastana, err := odczytajSondeKondycji(odczyt.QueryRowContext(ctx, sonda.Kod, KontoOperatora(ctx)))
		switch {
		case errors.Is(err, sql.ErrNoRows):
			powstala = true
		case err != nil:
			return fmt.Errorf("dane: nie można odczytać sondy kondycji %q: %w", sonda.Kod, err)
		default:
			sonda.Utworzono = zastana.Utworzono
			sonda.OstatniStan = zastana.OstatniStan
			sonda.OstatniPrzebieg = zastana.OstatniPrzebieg
		}

		tekst := aktualizujSondeKondycji
		if powstala {
			tekst = wstawSondeKondycji
		}
		zapis, err := r.zapytania.wTransakcji(ctx, transakcja, tekst)
		if err != nil {
			return err
		}
		wspolne := []any{sonda.Nazwa, sonda.Rodzaj, sonda.Cel, tekstDoKolumny(sonda.KomponentKod),
			sonda.OdstepMs, liczbaDoKolumny(sonda.LimitCzasuMs), liczbaDoKolumny(sonda.OczekiwanyKod),
			tekstDoKolumny(sonda.TrescWysylana), liczbaRzeczywistaDoKolumny(sonda.CelDostepnosci),
			liczbaLogiczna(sonda.Czynna)}
		var argumenty []any
		if powstala {
			argumenty = append([]any{sonda.Kod}, wspolne...)
			argumenty = append(argumenty, sonda.Utworzono, liczbaDoKolumny(sonda.Zaktualizowano),
				KontoOperatora(ctx))
		} else {
			argumenty = append(wspolne, liczbaDoKolumny(sonda.Zaktualizowano), sonda.Kod,
				KontoOperatora(ctx))
		}
		if _, err := zapis.ExecContext(ctx, argumenty...); err != nil {
			return fmt.Errorf("dane: nie można zapisać sondy kondycji %q: %w", sonda.Kod, err)
		}
		return nil
	})
	if err != nil {
		return SondaKondycji{}, false, err
	}
	zapisana, err := r.Sonda(ctx, sonda.Kod)
	return zapisana, powstala, err
}

// Sonda zwraca jedną definicję sondy po kodzie; brak wiersza sygnalizuje ErrBrakWiersza.
func (r *repozytoriumKondycji) Sonda(ctx context.Context, kod string) (SondaKondycji, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzSondeKondycji)
	if err != nil {
		return SondaKondycji{}, err
	}
	sonda, err := odczytajSondeKondycji(polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return SondaKondycji{}, fmt.Errorf("dane: sonda kondycji %q nie istnieje: %w", kod, ErrBrakWiersza)
	}
	return sonda, err
}

// Sondy zwraca definicje zawężone rodzajem, komponentem i stanem czynności, od najstarszej.
func (r *repozytoriumKondycji) Sondy(ctx context.Context, rodzaj, komponent string,
	tylkoCzynne bool, granica int) ([]SondaKondycji, error) {

	warunki := []string{"1 = 1", WarunekKonta}
	argumenty := []any{KontoOperatora(ctx)}
	if strings.TrimSpace(rodzaj) != "" {
		warunki = append(warunki, "s.rodzaj = ?")
		argumenty = append(argumenty, rodzaj)
	}
	if strings.TrimSpace(komponent) != "" {
		warunki = append(warunki, "s.komponent_kod = ?")
		argumenty = append(argumenty, komponent)
	}
	if tylkoCzynne {
		warunki = append(warunki, "s.czynna = 1")
	}
	tekst := `SELECT ` + kolumnySondyKondycji + ` FROM sonda_kondycji s WHERE ` +
		strings.Join(warunki, " AND ") + ` ORDER BY s.id`
	if granica > 0 {
		tekst += fmt.Sprintf(" LIMIT %d", granica)
	}

	wiersze, err := r.db.QueryContext(ctx, tekst, argumenty...)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać sond kondycji: %w", err)
	}
	defer wiersze.Close()

	lista := []SondaKondycji{}
	for wiersze.Next() {
		sonda, err := odczytajSondeKondycji(wiersze)
		if err != nil {
			return nil, err
		}
		lista = append(lista, sonda)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt sond kondycji: %w", err)
	}
	return lista, nil
}

// UsunSonde wykreśla definicję z serią wyników i oddaje liczbę pomiarów, liczoną przed kaskadą.
func (r *repozytoriumKondycji) UsunSonde(ctx context.Context, kod string) (int, error) {
	usunietych := 0
	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		licznik, err := r.zapytania.wTransakcji(ctx, transakcja, policzWynikiSondyKondycji)
		if err != nil {
			return err
		}
		if err := licznik.QueryRowContext(ctx, kod).Scan(&usunietych); err != nil {
			return fmt.Errorf("dane: nie można policzyć wyników sondy kondycji %q: %w", kod, err)
		}
		kasowanie, err := r.zapytania.wTransakcji(ctx, transakcja, usunSondeKondycji)
		if err != nil {
			return err
		}
		wynik, err := kasowanie.ExecContext(ctx, kod)
		if err != nil {
			return fmt.Errorf("dane: nie można usunąć sondy kondycji %q: %w", kod, err)
		}
		zmienione, err := wynik.RowsAffected()
		if err == nil && zmienione == 0 {
			return fmt.Errorf("dane: sonda kondycji %q nie istnieje: %w", kod, ErrBrakWiersza)
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return usunietych, nil
}

// ZapiszWynik dopisuje pomiar i podnosi odbicie w definicji sondy w jednej transakcji.
func (r *repozytoriumKondycji) ZapiszWynik(ctx context.Context,
	wynik WynikSondyKondycji) (WynikSondyKondycji, error) {

	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		zapis, err := r.zapytania.wTransakcji(ctx, transakcja, wstawWynikKondycji)
		if err != nil {
			return err
		}
		if _, err := zapis.ExecContext(ctx, wynik.Kod, wynik.SondaKod, wynik.Stan, wynik.Wykonano,
			liczbaDoKolumny(wynik.CzasOdpowiedziMs), liczbaDoKolumny(wynik.StatusHttp),
			tekstDoKolumny(wynik.Szczegol), tekstDoKolumny(wynik.BladKod)); err != nil {
			return fmt.Errorf("dane: nie można zapisać wyniku sondy kondycji %q: %w", wynik.Kod, err)
		}
		odbicie, err := r.zapytania.wTransakcji(ctx, transakcja, odbijPrzebiegSondyKondycji)
		if err != nil {
			return err
		}
		if _, err := odbicie.ExecContext(ctx, wynik.Stan, wynik.Wykonano, wynik.SondaKod, KontoOperatora(ctx)); err != nil {
			return fmt.Errorf("dane: nie można odnotować przebiegu sondy kondycji %q: %w",
				wynik.SondaKod, err)
		}
		return nil
	})
	if err != nil {
		return WynikSondyKondycji{}, err
	}
	return wynik, nil
}

// Wyniki zwraca serię pomiarów od najnowszego wraz z liczbą wszystkich bez granicy.
func (r *repozytoriumKondycji) Wyniki(ctx context.Context,
	sito SitoWynikowKondycji) ([]WynikSondyKondycji, int, error) {

	warunki := []string{"1 = 1", `EXISTS (SELECT 1 FROM sonda_kondycji s
		WHERE s.identyfikator_zewnetrzny = wynik_sondy_kondycji.sonda_kod AND ` + WarunekKonta + `)`}
	argumenty := []any{KontoOperatora(ctx)}
	if strings.TrimSpace(sito.SondaKod) != "" {
		warunki = append(warunki, "sonda_kod = ?")
		argumenty = append(argumenty, sito.SondaKod)
	}
	if strings.TrimSpace(sito.Stan) != "" {
		warunki = append(warunki, "stan = ?")
		argumenty = append(argumenty, sito.Stan)
	}
	if sito.OdCzasu > 0 {
		warunki = append(warunki, "wykonano >= ?")
		argumenty = append(argumenty, sito.OdCzasu)
	}
	if sito.DoCzasu > 0 {
		warunki = append(warunki, "wykonano <= ?")
		argumenty = append(argumenty, sito.DoCzasu)
	}
	gdzie := strings.Join(warunki, " AND ")

	wszystkich := 0
	if err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM wynik_sondy_kondycji WHERE `+gdzie, argumenty...).
		Scan(&wszystkich); err != nil {
		return nil, 0, fmt.Errorf("dane: nie można policzyć wyników sond kondycji: %w", err)
	}

	tekst := `SELECT ` + kolumnyWynikuKondycji + ` FROM wynik_sondy_kondycji WHERE ` + gdzie +
		` ORDER BY wykonano DESC, id DESC`
	if sito.Granica > 0 {
		tekst += fmt.Sprintf(" LIMIT %d", sito.Granica)
	}
	wiersze, err := r.db.QueryContext(ctx, tekst, argumenty...)
	if err != nil {
		return nil, 0, fmt.Errorf("dane: nie można odczytać wyników sond kondycji: %w", err)
	}
	defer wiersze.Close()

	lista := []WynikSondyKondycji{}
	for wiersze.Next() {
		wynik, err := odczytajWynikKondycji(wiersze)
		if err != nil {
			return nil, 0, err
		}
		lista = append(lista, wynik)
	}
	if err := wiersze.Err(); err != nil {
		return nil, 0, fmt.Errorf("dane: przerwany odczyt wyników sond kondycji: %w", err)
	}
	return lista, wszystkich, nil
}

func odczytajSondeKondycji(s skaner) (SondaKondycji, error) {
	var sonda SondaKondycji
	var komponent, ostatniStan, tresc sql.NullString
	var limit, oczekiwany, przebieg, zaktualizowano sql.NullInt64
	var cel sql.NullFloat64
	var czynna int
	err := s.Scan(&sonda.ID, &sonda.Kod, &sonda.Nazwa, &sonda.Rodzaj, &sonda.Cel, &komponent,
		&sonda.OdstepMs, &limit, &oczekiwany, &tresc, &cel, &czynna, &ostatniStan, &przebieg,
		&sonda.Utworzono, &zaktualizowano)
	if err != nil {
		return SondaKondycji{}, err
	}
	sonda.KomponentKod = tekstZKolumny(komponent)
	sonda.LimitCzasuMs = liczbaZKolumny(limit)
	sonda.OczekiwanyKod = liczbaZKolumny(oczekiwany)
	sonda.TrescWysylana = tekstZKolumny(tresc)
	sonda.CelDostepnosci = liczbaRzeczywistaZKolumny(cel)
	sonda.Czynna = czynna == 1
	sonda.OstatniStan = tekstZKolumny(ostatniStan)
	sonda.OstatniPrzebieg = liczbaZKolumny(przebieg)
	sonda.Zaktualizowano = liczbaZKolumny(zaktualizowano)
	return sonda, nil
}

func odczytajWynikKondycji(s skaner) (WynikSondyKondycji, error) {
	var wynik WynikSondyKondycji
	var czas, status sql.NullInt64
	var szczegol, blad sql.NullString
	err := s.Scan(&wynik.Kod, &wynik.SondaKod, &wynik.Stan, &wynik.Wykonano,
		&czas, &status, &szczegol, &blad)
	if err != nil {
		return WynikSondyKondycji{}, err
	}
	wynik.CzasOdpowiedziMs = liczbaZKolumny(czas)
	wynik.StatusHttp = liczbaZKolumny(status)
	wynik.Szczegol = tekstZKolumny(szczegol)
	wynik.BladKod = tekstZKolumny(blad)
	return wynik, nil
}

// Puls wykonuje najtańsze możliwe zapytanie do bazy, by zmierzyć czas jej obiegu; liczy się sama odpowiedź.
func (r *repozytoriumKondycji) Puls(ctx context.Context) error {
	var jeden int
	if err := r.db.QueryRowContext(ctx, "SELECT 1").Scan(&jeden); err != nil {
		return fmt.Errorf("dane: baza serwera nie odpowiedziała na puls: %w", err)
	}
	return nil
}
