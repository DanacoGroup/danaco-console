// Plik dobudowuje obszar badań po stronie danych: odkrywanie źródeł, monitory
// tematów, przestrzeń badania oraz style cytowania (kontrakt: badania_dobudowa.go).
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type WynikOdkryciaBadania struct {
	Klucz           string
	Okno            string
	Tytul           string
	Adres           *string
	Autorzy         []string
	Rok             *int64
	Dostawca        string
	Identyfikator   *string
	Fragment        *string
	OtwartyDostep   *int64
	Duplikat        bool
	Odrzucony       bool
	PowodOdrzucenia *string
	ZrodloKod       *string
}

type MonitorBadania struct {
	Kod           string
	Okno          string
	Rodzaj        string
	Zapytanie     *string
	Adres         *string
	InterwalMinut *int64
	Wlaczony      bool
	Oczekujace    int
	OdswiezonoO   *string
}

type PytanieBadania struct {
	Kod       string
	Tekst     string
	Kolejnosc int
}

type SzczegolyPrzestrzeniBadania struct {
	Odbiorca *string
	Protokol *string
	Granice  *string
	Notatka  *string
}

type StylCytowaniaBadania struct {
	Kod    string
	Nazwa  string
	Wlasny bool
}

// Więzy UNIQUE tych tabel obejmują całą tabelę, nie konto: zero zmienionych
// wierszy przy DO UPDATE jest odmową, nie powodzeniem.
func (r *repozytoriumBadan) zapiszWGranicyBadania(ctx context.Context, sqlTekst, byt,
	wskazanie string, argumenty ...any) error {

	polecenie, err := r.zapytania.przygotuj(ctx, sqlTekst)
	if err != nil {
		return err
	}
	wynik, err := polecenie.ExecContext(ctx, argumenty...)
	if err != nil {
		return fmt.Errorf("dane: badania — polecenie nie powiodło się: %w", err)
	}
	return sprawdzTrafienieZapisu(wynik, byt, wskazanie)
}

// Pozycja już znana zachowuje swój przesiew przy ponownym zapytaniu.
func (r *repozytoriumBadan) ZapiszWynikiOdkrycia(ctx context.Context, wyniki []WynikOdkryciaBadania) error {
	for _, wynik := range wyniki {
		if strings.TrimSpace(wynik.Klucz) == "" {
			continue
		}
		err := r.zapiszWGranicyBadania(ctx, `INSERT INTO wynik_odkrycia_badania
		    (klucz, okno, tytul, adres, autorzy, rok, dostawca, identyfikator, fragment,
		     otwarty_dostep, duplikat, zrodlo_kod, konto_id)
		    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, `+WskazanieKonta+`)
		    ON CONFLICT(klucz, COALESCE(konto_id, 0)) DO UPDATE SET
		        tytul = excluded.tytul, adres = excluded.adres, autorzy = excluded.autorzy,
		        rok = excluded.rok, dostawca = excluded.dostawca,
		        identyfikator = excluded.identyfikator, fragment = excluded.fragment,
		        otwarty_dostep = excluded.otwarty_dostep, duplikat = excluded.duplikat,
		        zrodlo_kod = COALESCE(excluded.zrodlo_kod, wynik_odkrycia_badania.zrodlo_kod)
		    WHERE `+WarunekKonta,
			"wynik odkrycia", wynik.Klucz,
			wynik.Klucz, wynik.Okno, wynik.Tytul, tekstDoKolumny(wynik.Adres),
			listaJakoBadania(wynik.Autorzy), liczbaDoKolumny(wynik.Rok), wynik.Dostawca,
			tekstDoKolumny(wynik.Identyfikator), tekstDoKolumny(wynik.Fragment),
			liczbaDoKolumny(wynik.OtwartyDostep), wartoscCalkowitaBadania(wynik.Duplikat),
			tekstDoKolumny(wynik.ZrodloKod), KontoOperatora(ctx), KontoOperatora(ctx))
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *repozytoriumBadan) WynikiOdkrycia(ctx context.Context, okno string) ([]WynikOdkryciaBadania, error) {
	wiersze, err := r.pytajBadania(ctx, `SELECT klucz, okno, tytul, adres, autorzy, rok, dostawca,
	        identyfikator, fragment, otwarty_dostep, duplikat, odrzucony, powod_odrzucenia, zrodlo_kod
	    FROM wynik_odkrycia_badania WHERE okno = ? AND `+WarunekKonta+`
	    ORDER BY utworzono DESC, id DESC`, okno, KontoOperatora(ctx))
	if err != nil {
		return nil, err
	}
	defer wiersze.Close()

	lista := []WynikOdkryciaBadania{}
	for wiersze.Next() {
		var wynik WynikOdkryciaBadania
		var adres, identyfikator, fragment, powod, zrodlo sql.NullString
		var autorzy string
		var rok, otwarty sql.NullInt64
		var duplikat, odrzucony int64
		if err := wiersze.Scan(&wynik.Klucz, &wynik.Okno, &wynik.Tytul, &adres, &autorzy, &rok,
			&wynik.Dostawca, &identyfikator, &fragment, &otwarty, &duplikat, &odrzucony,
			&powod, &zrodlo); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz wyniku odkrycia: %w", err)
		}
		wynik.Adres = tekstZKolumny(adres)
		wynik.Autorzy = listaZBadania(autorzy)
		wynik.Rok = liczbaZKolumny(rok)
		wynik.Identyfikator = tekstZKolumny(identyfikator)
		wynik.Fragment = tekstZKolumny(fragment)
		wynik.OtwartyDostep = liczbaZKolumny(otwarty)
		wynik.Duplikat = duplikat != 0
		wynik.Odrzucony = odrzucony != 0
		wynik.PowodOdrzucenia = tekstZKolumny(powod)
		wynik.ZrodloKod = tekstZKolumny(zrodlo)
		lista = append(lista, wynik)
	}
	return lista, wiersze.Err()
}

func (r *repozytoriumBadan) OdrzucWynikiOdkrycia(ctx context.Context, okno string,
	klucze []string, powod string) (int, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, `UPDATE wynik_odkrycia_badania
	    SET odrzucony = 1, powod_odrzucenia = ?
	    WHERE okno = ? AND klucz = ? AND odrzucony = 0 AND `+WarunekKonta)
	if err != nil {
		return 0, err
	}
	odrzucone := 0
	for _, klucz := range klucze {
		wynik, err := polecenie.ExecContext(ctx, powod, okno, klucz, KontoOperatora(ctx))
		if err != nil {
			return odrzucone, fmt.Errorf("dane: nie można odrzucić wyniku %q: %w", klucz, err)
		}
		ile, _ := wynik.RowsAffected()
		odrzucone += int(ile)
	}
	return odrzucone, nil
}

func (r *repozytoriumBadan) ZapiszMonitor(ctx context.Context, m MonitorBadania) (MonitorBadania, error) {
	err := r.zapiszWGranicyBadania(ctx, `INSERT INTO monitor_badania
	    (identyfikator_zewnetrzny, okno, rodzaj, zapytanie, adres, interwal_minut, wlaczony,
	     oczekujace, odswiezono_o, konto_id)
	    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, `+WskazanieKonta+`)
	    ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	        rodzaj = excluded.rodzaj, zapytanie = excluded.zapytanie, adres = excluded.adres,
	        interwal_minut = excluded.interwal_minut, wlaczony = excluded.wlaczony
	    WHERE `+WarunekKonta,
		"monitor badania", m.Kod,
		m.Kod, m.Okno, m.Rodzaj, tekstDoKolumny(m.Zapytanie), tekstDoKolumny(m.Adres),
		liczbaDoKolumny(m.InterwalMinut), wartoscCalkowitaBadania(m.Wlaczony), m.Oczekujace,
		tekstDoKolumny(m.OdswiezonoO), KontoOperatora(ctx), KontoOperatora(ctx))
	if err != nil {
		return MonitorBadania{}, err
	}
	return r.Monitor(ctx, m.Kod)
}

// Osobno od ZapiszMonitor: ustawienie monitora przez operatora nie kasuje licznika.
func (r *repozytoriumBadan) ZapiszMonitorZeStanem(ctx context.Context,
	m MonitorBadania) (MonitorBadania, error) {

	// Zero zmienionych wierszy odmawia odczyt poniżej brakiem wiersza.
	err := r.wykonajBadania(ctx, `UPDATE monitor_badania
	    SET oczekujace = ?, odswiezono_o = ?
	    WHERE identyfikator_zewnetrzny = ? AND `+WarunekKonta,
		m.Oczekujace, tekstDoKolumny(m.OdswiezonoO), m.Kod, KontoOperatora(ctx))
	if err != nil {
		return MonitorBadania{}, err
	}
	return r.Monitor(ctx, m.Kod)
}

const zapytanieMonitoraBadania = `SELECT identyfikator_zewnetrzny, okno, rodzaj, zapytanie, adres,
	        interwal_minut, wlaczony, oczekujace, odswiezono_o FROM monitor_badania `

func (r *repozytoriumBadan) Monitor(ctx context.Context, kod string) (MonitorBadania, error) {
	polecenie, err := r.zapytania.przygotuj(ctx,
		zapytanieMonitoraBadania+`WHERE identyfikator_zewnetrzny = ? AND `+WarunekKonta)
	if err != nil {
		return MonitorBadania{}, err
	}
	m, err := odczytajMonitorBadania(polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return MonitorBadania{}, ErrBrakWiersza
	}
	if err != nil {
		return MonitorBadania{}, fmt.Errorf("dane: nieczytelny monitor %q: %w", kod, err)
	}
	return m, nil
}

func (r *repozytoriumBadan) Monitory(ctx context.Context, okno string,
	tylkoWlaczone bool) ([]MonitorBadania, error) {

	sqlTekst := zapytanieMonitoraBadania + `WHERE okno = ? AND ` + WarunekKonta
	if tylkoWlaczone {
		sqlTekst += ` AND wlaczony = 1`
	}
	sqlTekst += ` ORDER BY identyfikator_zewnetrzny`
	wiersze, err := r.pytajBadania(ctx, sqlTekst, okno, KontoOperatora(ctx))
	if err != nil {
		return nil, err
	}
	defer wiersze.Close()

	lista := []MonitorBadania{}
	for wiersze.Next() {
		m, err := odczytajMonitorBadania(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz monitora: %w", err)
		}
		lista = append(lista, m)
	}
	return lista, wiersze.Err()
}

func odczytajMonitorBadania(wiersz skaner) (MonitorBadania, error) {
	var m MonitorBadania
	var zapytanie, adres, odswiezono sql.NullString
	var interwal sql.NullInt64
	var wlaczony int64
	err := wiersz.Scan(&m.Kod, &m.Okno, &m.Rodzaj, &zapytanie, &adres, &interwal,
		&wlaczony, &m.Oczekujace, &odswiezono)
	if err != nil {
		return MonitorBadania{}, err
	}
	m.Zapytanie = tekstZKolumny(zapytanie)
	m.Adres = tekstZKolumny(adres)
	m.InterwalMinut = liczbaZKolumny(interwal)
	m.Wlaczony = wlaczony != 0
	m.OdswiezonoO = tekstZKolumny(odswiezono)
	return m, nil
}

func (r *repozytoriumBadan) UstawPytaniaBadania(ctx context.Context,
	pytania []PytanieBadania) ([]PytanieBadania, error) {

	err := wTransakcji(ctx, r.db, func(t *sql.Tx) error {
		czyszczenie, err := r.zapytania.wTransakcji(ctx, t,
			`DELETE FROM pytanie_badania WHERE `+WarunekKonta)
		if err != nil {
			return err
		}
		if _, err := czyszczenie.ExecContext(ctx, KontoOperatora(ctx)); err != nil {
			return fmt.Errorf("dane: nie można wyczyścić pytań badania: %w", err)
		}
		zapis, err := r.zapytania.wTransakcji(ctx, t,
			`INSERT INTO pytanie_badania (identyfikator_zewnetrzny, tekst, kolejnosc, konto_id)
			 VALUES (?, ?, ?, `+WskazanieKonta+`)`)
		if err != nil {
			return err
		}
		for _, pytanie := range pytania {
			if _, err := zapis.ExecContext(ctx, pytanie.Kod, pytanie.Tekst, pytanie.Kolejnosc,
				KontoOperatora(ctx)); err != nil {
				return fmt.Errorf("dane: nie można zapisać pytania %q: %w", pytanie.Kod, err)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return r.PytaniaBadania(ctx)
}

func (r *repozytoriumBadan) PytaniaBadania(ctx context.Context) ([]PytanieBadania, error) {
	wiersze, err := r.pytajBadania(ctx,
		`SELECT identyfikator_zewnetrzny, tekst, kolejnosc FROM pytanie_badania
		 WHERE `+WarunekKonta+` ORDER BY kolejnosc, id`, KontoOperatora(ctx))
	if err != nil {
		return nil, err
	}
	defer wiersze.Close()

	lista := []PytanieBadania{}
	for wiersze.Next() {
		var pytanie PytanieBadania
		if err := wiersze.Scan(&pytanie.Kod, &pytanie.Tekst, &pytanie.Kolejnosc); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz pytania badania: %w", err)
		}
		lista = append(lista, pytanie)
	}
	return lista, wiersze.Err()
}

// Pole puste zostaje bez zmiany: wywołujący podaje to, co zmienia.
func (r *repozytoriumBadan) UstawSzczegolyPrzestrzeni(ctx context.Context, s SzczegolyPrzestrzeniBadania) error {
	if err := r.zapewnijPrzestrzen(ctx); err != nil {
		return err
	}
	return r.zapiszWGranicyBadania(ctx, `UPDATE przestrzen_badania SET
	        odbiorca = COALESCE(?, odbiorca),
	        protokol = COALESCE(?, protokol),
	        granice  = COALESCE(?, granice),
	        zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	    WHERE `+WarunekKonta,
		"przestrzeń badania", "1",
		tekstDoKolumny(s.Odbiorca), tekstDoKolumny(s.Protokol), tekstDoKolumny(s.Granice),
		KontoOperatora(ctx))
}

func (r *repozytoriumBadan) SzczegolyPrzestrzeniBadania(ctx context.Context) (SzczegolyPrzestrzeniBadania, error) {
	if err := r.zapewnijPrzestrzen(ctx); err != nil {
		return SzczegolyPrzestrzeniBadania{}, err
	}
	polecenie, err := r.zapytania.przygotuj(ctx,
		`SELECT odbiorca, protokol, granice, notatka FROM przestrzen_badania WHERE `+WarunekKonta)
	if err != nil {
		return SzczegolyPrzestrzeniBadania{}, err
	}
	var odbiorca, protokol, granice, notatka sql.NullString
	err = polecenie.QueryRowContext(ctx, KontoOperatora(ctx)).Scan(&odbiorca, &protokol, &granice, &notatka)
	if errors.Is(err, sql.ErrNoRows) {
		return SzczegolyPrzestrzeniBadania{}, nil
	}
	if err != nil {
		return SzczegolyPrzestrzeniBadania{}, fmt.Errorf("dane: nieczytelna przestrzeń badania: %w", err)
	}
	return SzczegolyPrzestrzeniBadania{
		Odbiorca: tekstZKolumny(odbiorca), Protokol: tekstZKolumny(protokol),
		Granice: tekstZKolumny(granice), Notatka: tekstZKolumny(notatka),
	}, nil
}

func (r *repozytoriumBadan) UstawNotatkePrzestrzeni(ctx context.Context, tresc string) (string, error) {
	if err := r.zapewnijPrzestrzen(ctx); err != nil {
		return "", err
	}
	err := r.zapiszWGranicyBadania(ctx, `UPDATE przestrzen_badania
	    SET notatka = ?, zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	    WHERE `+WarunekKonta, "przestrzeń badania", "1", tresc, KontoOperatora(ctx))
	if err != nil {
		return "", err
	}
	polecenie, err := r.zapytania.przygotuj(ctx,
		`SELECT zaktualizowano FROM przestrzen_badania WHERE `+WarunekKonta)
	if err != nil {
		return "", err
	}
	var chwila string
	if err := polecenie.QueryRowContext(ctx, KontoOperatora(ctx)).Scan(&chwila); err != nil {
		return "", fmt.Errorf("dane: nieczytelna chwila przestrzeni badania: %w", err)
	}
	return chwila, nil
}

// Od kroku 495 każde konto ma własny wiersz przestrzeni; wskaźnik po wyrażeniu
// pilnuje, żeby drugie wywołanie go nie zdublowało.
func (r *repozytoriumBadan) zapewnijPrzestrzen(ctx context.Context) error {
	return r.wykonajBadania(ctx, `INSERT INTO przestrzen_badania (zakres, konto_id)
	    VALUES ('', `+WskazanieKonta+`)
	    ON CONFLICT(COALESCE(konto_id, 0)) DO NOTHING`, KontoOperatora(ctx))
}

func (r *repozytoriumBadan) ZapiszStylCytowania(ctx context.Context, kod, nazwa string) error {
	return r.zapiszWGranicyBadania(ctx, `INSERT INTO styl_cytowania_badania
	    (identyfikator_zewnetrzny, nazwa, wlasny, konto_id) VALUES (?, ?, 1, `+WskazanieKonta+`)
	    ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET nazwa = excluded.nazwa
	    WHERE `+WarunekKonta,
		"styl cytowania", kod, kod, nazwa, KontoOperatora(ctx), KontoOperatora(ctx))
}

func (r *repozytoriumBadan) StyleCytowania(ctx context.Context) ([]StylCytowaniaBadania, error) {
	wiersze, err := r.pytajBadania(ctx,
		`SELECT identyfikator_zewnetrzny, nazwa, wlasny FROM styl_cytowania_badania
		 WHERE `+WarunekKonta+` ORDER BY nazwa`, KontoOperatora(ctx))
	if err != nil {
		return nil, err
	}
	defer wiersze.Close()

	lista := []StylCytowaniaBadania{}
	for wiersze.Next() {
		var styl StylCytowaniaBadania
		var wlasny int64
		if err := wiersze.Scan(&styl.Kod, &styl.Nazwa, &wlasny); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz stylu cytowania: %w", err)
		}
		styl.Wlasny = wlasny != 0
		lista = append(lista, styl)
	}
	return lista, wiersze.Err()
}
