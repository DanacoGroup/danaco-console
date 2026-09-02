// Plik dobudowuje obszar badań po stronie danych: odkrywanie źródeł, monitory tematów,
// przestrzeń badania oraz style cytowania. Kontrakt obszaru deklaruje plik
// badania_dobudowa.go.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// WynikOdkryciaBadania odwzorowuje wiersz tabeli wynik_odkrycia_badania — pozycja zwrócona
// przez dostawcę wyszukiwania, zapamiętana po to, żeby dało się ją odrzucić i policzyć
// w przesiewie PRISMA.
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

// MonitorBadania odwzorowuje wiersz tabeli monitor_badania: subskrypcję zapytania albo
// adresu, którą instalacja odświeża w ustalonym interwale i zgłasza jako nowe pozycje.
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

// PytanieBadania odwzorowuje wiersz tabeli pytanie_badania: jedno pytanie badawcze
// przestrzeni badania wraz z jego kolejnością w wykazie.
type PytanieBadania struct {
	Kod       string
	Tekst     string
	Kolejnosc int
}

// SzczegolyPrzestrzeniBadania niesie odbiorcę, protokół, granice oraz notatkę roboczą
// jedynej przestrzeni badania instalacji.
type SzczegolyPrzestrzeniBadania struct {
	Odbiorca *string
	Protokol *string
	Granice  *string
	Notatka  *string
}

// StylCytowaniaBadania odwzorowuje wiersz tabeli styl_cytowania_badania: styl cytowania
// dodany własnoręcznie przez operatora, poza zestawem stylów wbudowanych.
type StylCytowaniaBadania struct {
	Kod    string
	Nazwa  string
	Wlasny bool
}

// zapiszWGranicyBadania wykonuje polecenie zawężone kontem i odróżnia zapis
// wykonany od odciętego granicą. Więzy UNIQUE tych tabel obejmują je w całości,
// więc konflikt sięga wiersza cudzego konta, a zero zmienionych wierszy przy
// DO UPDATE jest odmową, nie powodzeniem.
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

// ── Wyniki odkrycia ────────────────────────────────────────────────────────

// ZapiszWynikiOdkrycia utrwala pozycje zwrócone przez dostawców. Pozycja już
// znana zachowuje swój przesiew: odrzucenie nie ma prawa zniknąć dlatego, że
// to samo zapytanie puszczono po raz drugi.
func (r *repozytoriumBadan) ZapiszWynikiOdkrycia(ctx context.Context, wyniki []WynikOdkryciaBadania) error {
	for _, wynik := range wyniki {
		if strings.TrimSpace(wynik.Klucz) == "" {
			continue
		}
		err := r.zapiszWGranicyBadania(ctx, `INSERT INTO wynik_odkrycia_badania
		    (klucz, okno, tytul, adres, autorzy, rok, dostawca, identyfikator, fragment,
		     otwarty_dostep, duplikat, zrodlo_kod, konto_id)
		    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, `+WskazanieKonta+`)
		    ON CONFLICT(klucz) DO UPDATE SET
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

// WynikiOdkrycia oddaje pozycje zapamiętane dla okna badania, od najświeższej zapisanej,
// wraz ze stanem ich przesiewu.
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

// OdrzucWynikiOdkrycia znakuje wskazane pozycje jako odrzucone z uzasadnieniem
// i oddaje liczbę pozycji rzeczywiście odrzuconych — klucz nieznany nie liczy
// się jako odrzucenie.
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

// ── Monitory tematów ───────────────────────────────────────────────────────

// ZapiszMonitor zakłada monitor tematu albo nadpisuje zastany wiersz o tym samym
// identyfikatorze zewnętrznym, nie naruszając jego stanu odświeżenia.
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

// ZapiszMonitorZeStanem zapisuje monitor razem ze stanem odświeżenia: licznikiem
// oczekujących i chwilą ostatniego przebiegu. Osobno od ZapiszMonitor, bo tam stan jest
// świadomie pomijany — ustawienie monitora przez operatora nie ma skasować licznika.
func (r *repozytoriumBadan) ZapiszMonitorZeStanem(ctx context.Context,
	m MonitorBadania) (MonitorBadania, error) {

	// Zero zmienionych wierszy nie idzie tu przez sprawdzTrafienieZapisu: przy
	// zwykłym UPDATE znaczy tak samo „monitora nie ma" jak „monitor jest cudzy",
	// a odmowę wypowiada odczyt poniżej brakiem wiersza.
	err := r.wykonajBadania(ctx, `UPDATE monitor_badania
	    SET oczekujace = ?, odswiezono_o = ?
	    WHERE identyfikator_zewnetrzny = ? AND `+WarunekKonta,
		m.Oczekujace, tekstDoKolumny(m.OdswiezonoO), m.Kod, KontoOperatora(ctx))
	if err != nil {
		return MonitorBadania{}, err
	}
	return r.Monitor(ctx, m.Kod)
}

// zapytanieMonitoraBadania jest wspólnym tekstem zapytania SQL, z którego korzystają
// funkcje odczytujące pojedynczy monitor oraz listy monitorów.
const zapytanieMonitoraBadania = `SELECT identyfikator_zewnetrzny, okno, rodzaj, zapytanie, adres,
	        interwal_minut, wlaczony, oczekujace, odswiezono_o FROM monitor_badania `

// Monitor oddaje jeden monitor tematu wskazany kodem zewnętrznym albo błąd
// ErrBrakWiersza, gdy taki monitor nie istnieje w bazie.
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

// Monitory oddaje monitory okna badania, opcjonalnie zawężone wyłącznie do włączonych,
// uporządkowane według identyfikatora zewnętrznego.
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

// odczytajMonitorBadania składa strukturę MonitorBadania z jednego wiersza wyniku zapytania,
// niezależnie od tego, czy pochodzi z pojedynczego odczytu, czy z iteracji po wielu wierszach.
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

// UstawPytaniaBadania wymienia w całości pytania badawcze przestrzeni: usuwa zastany
// zestaw i zapisuje podany w jednej transakcji, zachowując kolejność wejściową.
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

// PytaniaBadania oddaje pytania badawcze przestrzeni badania w kolejności ustalonej
// polem kolejnosc każdego wiersza.
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

// SzczegolyPrzestrzeniBadania oddaje odbiorcę, protokół, granice i notatkę roboczą jedynej
// przestrzeni badania instalacji.
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

// UstawNotatkePrzestrzeni nadpisuje notatkę roboczą przestrzeni badania i oddaje chwilę,
// w której zapis rzeczywiście nastąpił.
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

// zapewnijPrzestrzen zakłada jedyny wiersz przestrzeni, gdy jeszcze nie stoi. Przestrzeń
// jest bytem jednowierszowym instalacji, więc jej brak jest stanem początkowym, a nie usterką.
func (r *repozytoriumBadan) zapewnijPrzestrzen(ctx context.Context) error {
	return r.wykonajBadania(ctx, `INSERT OR IGNORE INTO przestrzen_badania (id, zakres) VALUES (1, '')`)
}

// ── Style cytowania ────────────────────────────────────────────────────────

// ZapiszStylCytowania utrwala styl cytowania własny operatora albo nadpisuje nazwę
// zastanego stylu o tym samym identyfikatorze zewnętrznym.
func (r *repozytoriumBadan) ZapiszStylCytowania(ctx context.Context, kod, nazwa string) error {
	return r.zapiszWGranicyBadania(ctx, `INSERT INTO styl_cytowania_badania
	    (identyfikator_zewnetrzny, nazwa, wlasny, konto_id) VALUES (?, ?, 1, `+WskazanieKonta+`)
	    ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET nazwa = excluded.nazwa
	    WHERE `+WarunekKonta,
		"styl cytowania", kod, kod, nazwa, KontoOperatora(ctx), KontoOperatora(ctx))
}

// StyleCytowania oddaje style własne zapisane w instalacji. Style wbudowane
// dokłada adapter — repozytorium mówi wyłącznie o tym, co leży w bazie.
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
