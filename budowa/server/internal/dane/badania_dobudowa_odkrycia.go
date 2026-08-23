// Odpowiedzialność pliku: dobudowa obszaru Research po stronie danych —
// odkrywanie źródeł (wyniki wyszukiwania i ich przesiew), monitory tematów,
// przestrzeń badania (pytania badawcze, odbiorca, protokół, notatka robocza)
// oraz style cytowania. Kontrakt obszaru deklaruje `badania_dobudowa.go`.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// WynikOdkryciaBadania to wiersz `wynik_odkrycia_badania` — pozycja zwrócona przez
// dostawcę wyszukiwania, zapamiętana po to, żeby dało się ją odrzucić i policzyć
// w przesiewie PRISMA (patrz nagłówek migracji 150).
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

// MonitorBadania to wiersz `monitor_badania`.
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

// PytanieBadania to wiersz `pytanie_badania`.
type PytanieBadania struct {
	Kod       string
	Tekst     string
	Kolejnosc int
}

// SzczegolyPrzestrzeniBadania to pola przestrzeni badania dobudowane migracją 150.
type SzczegolyPrzestrzeniBadania struct {
	Odbiorca *string
	Protokol *string
	Granice  *string
	Notatka  *string
}

// StylCytowaniaBadania to wiersz `styl_cytowania_badania`.
type StylCytowaniaBadania struct {
	Kod    string
	Nazwa  string
	Wlasny bool
}

// ── Wyniki odkrycia ────────────────────────────────────────────────────────

// ZapiszWynikiOdkrycia utrwala pozycje zwrócone przez dostawców. Pozycja już
// znana zachowuje swój przesiew: odrzucenie nie ma prawa zniknąć dlatego, że
// to samo zapytanie puszczono po raz drugi.
func (r *repozytoriumBadan) ZapiszWynikiOdkrycia(ctx context.Context, wyniki []WynikOdkryciaBadania) error {
	for _, wynik := range wyniki {
		if strings.TrimSpace(wynik.Klucz) == "" {
			continue
		}
		err := r.wykonajBadania(ctx, `INSERT INTO wynik_odkrycia_badania
		    (klucz, okno, tytul, adres, autorzy, rok, dostawca, identyfikator, fragment,
		     otwarty_dostep, duplikat, zrodlo_kod)
		    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		    ON CONFLICT(klucz) DO UPDATE SET
		        tytul = excluded.tytul, adres = excluded.adres, autorzy = excluded.autorzy,
		        rok = excluded.rok, dostawca = excluded.dostawca,
		        identyfikator = excluded.identyfikator, fragment = excluded.fragment,
		        otwarty_dostep = excluded.otwarty_dostep, duplikat = excluded.duplikat,
		        zrodlo_kod = COALESCE(excluded.zrodlo_kod, wynik_odkrycia_badania.zrodlo_kod)`,
			wynik.Klucz, wynik.Okno, wynik.Tytul, tekstDoKolumny(wynik.Adres),
			listaJakoBadania(wynik.Autorzy), liczbaDoKolumny(wynik.Rok), wynik.Dostawca,
			tekstDoKolumny(wynik.Identyfikator), tekstDoKolumny(wynik.Fragment),
			liczbaDoKolumny(wynik.OtwartyDostep), wartoscCalkowitaBadania(wynik.Duplikat),
			tekstDoKolumny(wynik.ZrodloKod))
		if err != nil {
			return err
		}
	}
	return nil
}

// WynikiOdkrycia oddaje pozycje zapamiętane dla okna badania.
func (r *repozytoriumBadan) WynikiOdkrycia(ctx context.Context, okno string) ([]WynikOdkryciaBadania, error) {
	wiersze, err := r.pytajBadania(ctx, `SELECT klucz, okno, tytul, adres, autorzy, rok, dostawca,
	        identyfikator, fragment, otwarty_dostep, duplikat, odrzucony, powod_odrzucenia, zrodlo_kod
	    FROM wynik_odkrycia_badania WHERE okno = ? ORDER BY utworzono DESC, id DESC`, okno)
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

// OdrzucWynikiOdkrycia znakuje wskazane pozycje jako odrzucone z uzasadnieniem
// i oddaje liczbę pozycji rzeczywiście odrzuconych — klucz nieznany nie liczy
// się jako odrzucenie.
func (r *repozytoriumBadan) OdrzucWynikiOdkrycia(ctx context.Context, okno string,
	klucze []string, powod string) (int, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, `UPDATE wynik_odkrycia_badania
	    SET odrzucony = 1, powod_odrzucenia = ?
	    WHERE okno = ? AND klucz = ? AND odrzucony = 0`)
	if err != nil {
		return 0, err
	}
	odrzucone := 0
	for _, klucz := range klucze {
		wynik, err := polecenie.ExecContext(ctx, powod, okno, klucz)
		if err != nil {
			return odrzucone, fmt.Errorf("dane: nie można odrzucić wyniku %q: %w", klucz, err)
		}
		ile, _ := wynik.RowsAffected()
		odrzucone += int(ile)
	}
	return odrzucone, nil
}

// ── Monitory tematów ───────────────────────────────────────────────────────

// ZapiszMonitor zakłada monitor albo nadpisuje zastany.
func (r *repozytoriumBadan) ZapiszMonitor(ctx context.Context, m MonitorBadania) (MonitorBadania, error) {
	err := r.wykonajBadania(ctx, `INSERT INTO monitor_badania
	    (identyfikator_zewnetrzny, okno, rodzaj, zapytanie, adres, interwal_minut, wlaczony,
	     oczekujace, odswiezono_o)
	    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	    ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	        rodzaj = excluded.rodzaj, zapytanie = excluded.zapytanie, adres = excluded.adres,
	        interwal_minut = excluded.interwal_minut, wlaczony = excluded.wlaczony`,
		m.Kod, m.Okno, m.Rodzaj, tekstDoKolumny(m.Zapytanie), tekstDoKolumny(m.Adres),
		liczbaDoKolumny(m.InterwalMinut), wartoscCalkowitaBadania(m.Wlaczony), m.Oczekujace,
		tekstDoKolumny(m.OdswiezonoO))
	if err != nil {
		return MonitorBadania{}, err
	}
	return r.Monitor(ctx, m.Kod)
}

// ZapiszMonitorZeStanem zapisuje monitor razem ze stanem odświeżenia: licznikiem
// oczekujących i chwilą ostatniego przebiegu. Osobno od `ZapiszMonitor`, bo tam
// stan jest świadomie pomijany — ustawienie monitora przez Operatora nie ma
// prawa skasować licznika nowych pozycji, którego Operator nie dotykał.
func (r *repozytoriumBadan) ZapiszMonitorZeStanem(ctx context.Context,
	m MonitorBadania) (MonitorBadania, error) {

	err := r.wykonajBadania(ctx, `UPDATE monitor_badania
	    SET oczekujace = ?, odswiezono_o = ? WHERE identyfikator_zewnetrzny = ?`,
		m.Oczekujace, tekstDoKolumny(m.OdswiezonoO), m.Kod)
	if err != nil {
		return MonitorBadania{}, err
	}
	return r.Monitor(ctx, m.Kod)
}

// zapytanieMonitoraBadania jest wspólnym odczytem monitorów.
const zapytanieMonitoraBadania = `SELECT identyfikator_zewnetrzny, okno, rodzaj, zapytanie, adres,
	        interwal_minut, wlaczony, oczekujace, odswiezono_o FROM monitor_badania `

// Monitor oddaje jeden monitor po kodzie.
func (r *repozytoriumBadan) Monitor(ctx context.Context, kod string) (MonitorBadania, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, zapytanieMonitoraBadania+`WHERE identyfikator_zewnetrzny = ?`)
	if err != nil {
		return MonitorBadania{}, err
	}
	m, err := odczytajMonitorBadania(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return MonitorBadania{}, ErrBrakWiersza
	}
	if err != nil {
		return MonitorBadania{}, fmt.Errorf("dane: nieczytelny monitor %q: %w", kod, err)
	}
	return m, nil
}

// Monitory oddaje monitory okna badania.
func (r *repozytoriumBadan) Monitory(ctx context.Context, okno string,
	tylkoWlaczone bool) ([]MonitorBadania, error) {

	sqlTekst := zapytanieMonitoraBadania + `WHERE okno = ?`
	if tylkoWlaczone {
		sqlTekst += ` AND wlaczony = 1`
	}
	sqlTekst += ` ORDER BY identyfikator_zewnetrzny`
	wiersze, err := r.pytajBadania(ctx, sqlTekst, okno)
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

// odczytajMonitorBadania składa strukturę z jednego wiersza wyniku.
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

// ── Przestrzeń badania ─────────────────────────────────────────────────────

// UstawPytaniaBadania wymienia w całości pytania badawcze przestrzeni.
func (r *repozytoriumBadan) UstawPytaniaBadania(ctx context.Context,
	pytania []PytanieBadania) ([]PytanieBadania, error) {

	err := wTransakcji(ctx, r.db, func(t *sql.Tx) error {
		czyszczenie, err := r.zapytania.wTransakcji(ctx, t, `DELETE FROM pytanie_badania`)
		if err != nil {
			return err
		}
		if _, err := czyszczenie.ExecContext(ctx); err != nil {
			return fmt.Errorf("dane: nie można wyczyścić pytań badania: %w", err)
		}
		zapis, err := r.zapytania.wTransakcji(ctx, t,
			`INSERT INTO pytanie_badania (identyfikator_zewnetrzny, tekst, kolejnosc) VALUES (?, ?, ?)`)
		if err != nil {
			return err
		}
		for _, pytanie := range pytania {
			if _, err := zapis.ExecContext(ctx, pytanie.Kod, pytanie.Tekst, pytanie.Kolejnosc); err != nil {
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

// PytaniaBadania oddaje pytania badawcze w kolejności zapisu.
func (r *repozytoriumBadan) PytaniaBadania(ctx context.Context) ([]PytanieBadania, error) {
	wiersze, err := r.pytajBadania(ctx,
		`SELECT identyfikator_zewnetrzny, tekst, kolejnosc FROM pytanie_badania ORDER BY kolejnosc, id`)
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

// UstawSzczegolyPrzestrzeni nadpisuje odbiorcę, protokół i granice badania.
// Pole puste zostaje bez zmiany — wywołujący podaje to, co zmienia.
func (r *repozytoriumBadan) UstawSzczegolyPrzestrzeni(ctx context.Context, s SzczegolyPrzestrzeniBadania) error {
	if err := r.zapewnijPrzestrzen(ctx); err != nil {
		return err
	}
	return r.wykonajBadania(ctx, `UPDATE przestrzen_badania SET
	        odbiorca = COALESCE(?, odbiorca),
	        protokol = COALESCE(?, protokol),
	        granice  = COALESCE(?, granice),
	        zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	    WHERE id = 1`,
		tekstDoKolumny(s.Odbiorca), tekstDoKolumny(s.Protokol), tekstDoKolumny(s.Granice))
}

// SzczegolyPrzestrzeniBadania oddaje odbiorcę, protokół, granice i notatkę roboczą.
func (r *repozytoriumBadan) SzczegolyPrzestrzeniBadania(ctx context.Context) (SzczegolyPrzestrzeniBadania, error) {
	if err := r.zapewnijPrzestrzen(ctx); err != nil {
		return SzczegolyPrzestrzeniBadania{}, err
	}
	polecenie, err := r.zapytania.przygotuj(ctx,
		`SELECT odbiorca, protokol, granice, notatka FROM przestrzen_badania WHERE id = 1`)
	if err != nil {
		return SzczegolyPrzestrzeniBadania{}, err
	}
	var odbiorca, protokol, granice, notatka sql.NullString
	if err := polecenie.QueryRowContext(ctx).Scan(&odbiorca, &protokol, &granice, &notatka); err != nil {
		return SzczegolyPrzestrzeniBadania{}, fmt.Errorf("dane: nieczytelna przestrzeń badania: %w", err)
	}
	return SzczegolyPrzestrzeniBadania{
		Odbiorca: tekstZKolumny(odbiorca), Protokol: tekstZKolumny(protokol),
		Granice: tekstZKolumny(granice), Notatka: tekstZKolumny(notatka),
	}, nil
}

// UstawNotatkePrzestrzeni nadpisuje notatkę roboczą i oddaje chwilę zapisu.
func (r *repozytoriumBadan) UstawNotatkePrzestrzeni(ctx context.Context, tresc string) (string, error) {
	if err := r.zapewnijPrzestrzen(ctx); err != nil {
		return "", err
	}
	err := r.wykonajBadania(ctx, `UPDATE przestrzen_badania
	    SET notatka = ?, zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now') WHERE id = 1`, tresc)
	if err != nil {
		return "", err
	}
	polecenie, err := r.zapytania.przygotuj(ctx,
		`SELECT zaktualizowano FROM przestrzen_badania WHERE id = 1`)
	if err != nil {
		return "", err
	}
	var chwila string
	if err := polecenie.QueryRowContext(ctx).Scan(&chwila); err != nil {
		return "", fmt.Errorf("dane: nieczytelna chwila przestrzeni badania: %w", err)
	}
	return chwila, nil
}

// zapewnijPrzestrzen zakłada jedyny wiersz przestrzeni, gdy jeszcze nie stoi.
// Przestrzeń jest bytem jednowierszowym instalacji (patrz `badania_raport.go`),
// więc jej brak jest stanem początkowym, a nie usterką.
func (r *repozytoriumBadan) zapewnijPrzestrzen(ctx context.Context) error {
	return r.wykonajBadania(ctx, `INSERT OR IGNORE INTO przestrzen_badania (id, zakres) VALUES (1, '')`)
}

// ── Style cytowania ────────────────────────────────────────────────────────

// ZapiszStylCytowania utrwala styl własny Operatora.
func (r *repozytoriumBadan) ZapiszStylCytowania(ctx context.Context, kod, nazwa string) error {
	return r.wykonajBadania(ctx, `INSERT INTO styl_cytowania_badania
	    (identyfikator_zewnetrzny, nazwa, wlasny) VALUES (?, ?, 1)
	    ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET nazwa = excluded.nazwa`, kod, nazwa)
}

// StyleCytowania oddaje style własne zapisane w instalacji. Style wbudowane
// dokłada adapter — repozytorium mówi wyłącznie o tym, co leży w bazie.
func (r *repozytoriumBadan) StyleCytowania(ctx context.Context) ([]StylCytowaniaBadania, error) {
	wiersze, err := r.pytajBadania(ctx,
		`SELECT identyfikator_zewnetrzny, nazwa, wlasny FROM styl_cytowania_badania ORDER BY nazwa`)
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
