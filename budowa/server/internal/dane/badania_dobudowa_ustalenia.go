// Odpowiedzialność pliku: dobudowa obszaru Research po stronie danych — byty
// krążące wokół ustalenia: adnotacje, kodowanie, sprzeczności, wątki, prowenancja.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type KotwicaBadania struct {
	Rodzaj   string
	Strona   *int64
	Od       *int64
	Do       *int64
	Selektor *string
	CzasMs   *int64
}

type AdnotacjaBadania struct {
	Kod          string
	ZrodloKod    string
	ZrodloTytul  string
	Rodzaj       string
	Cytat        *string
	Komentarz    *string
	Kolor        *string
	Kotwica      KotwicaBadania
	UstalenieKod *string
	Utworzono    string
}

type SzczegolyUstaleniaBadania struct {
	Rodzaj             string
	Waga               string
	Notatka            *string
	WymagaPotwierdzeni bool
	AdnotacjaKod       *string
	Kotwica            KotwicaBadania
}

type KodBadania struct {
	Kod         string
	Okno        string
	Nazwa       string
	Opis        *string
	Nadrzedny   *string
	Wystapienia int
}

type KomorkaMacierzyBadania struct {
	KodKodu       string
	ZrodloKod     string
	Liczba        int
	UstalenieKody []string
}

type SprzecznoscBadania struct {
	Kod                 string
	Okno                string
	Streszczenie        string
	RoznicaLiczbowa     *string
	Rozstrzygnieta      bool
	UstalenieRozstrzyga *string
	Uzasadnienie        *string
	UstalenieKody       []string
}

type WatekBadania struct {
	Kod           string
	Nazwa         string
	Automatyczny  bool
	UstalenieKody []string
}

type WpisProwenancjiBadania struct {
	UstalenieKod string
	OCzasie      string
	Aktor        string
	Czynnosc     string
	ZrodloKod    *string
	Strona       *int64
}

// ── Adnotacje lektury ──────────────────────────────────────────────────────

func (r *repozytoriumBadan) ZapiszAdnotacje(ctx context.Context,
	a AdnotacjaBadania) (AdnotacjaBadania, error) {

	id, err := r.idZrodla(ctx, a.ZrodloKod)
	if err != nil {
		return AdnotacjaBadania{}, err
	}
	rodzajKotwicy := a.Kotwica.Rodzaj
	if rodzajKotwicy == "" {
		rodzajKotwicy = "page"
	}
	err = r.zapiszWGranicyBadania(ctx, `INSERT INTO adnotacja_badania
	    (identyfikator_zewnetrzny, zrodlo_id, rodzaj, cytat, komentarz, kolor,
	     kotwica_rodzaj, kotwica_strona, kotwica_od, kotwica_do, kotwica_selektor,
	     kotwica_czas_ms, ustalenie_kod)
	    SELECT ?, id, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
	      FROM zrodlo_badania WHERE id = ? AND `+WarunekKonta+`
	    ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	        rodzaj = excluded.rodzaj, cytat = excluded.cytat,
	        komentarz = excluded.komentarz, kolor = excluded.kolor,
	        kotwica_rodzaj = excluded.kotwica_rodzaj,
	        kotwica_strona = excluded.kotwica_strona,
	        kotwica_od = excluded.kotwica_od, kotwica_do = excluded.kotwica_do,
	        kotwica_selektor = excluded.kotwica_selektor,
	        kotwica_czas_ms = excluded.kotwica_czas_ms,
	        ustalenie_kod = excluded.ustalenie_kod
	    WHERE `+warunekZrodlaAdnotacji, "adnotacja badania", a.Kod,
		a.Kod, a.Rodzaj, tekstDoKolumny(a.Cytat), tekstDoKolumny(a.Komentarz),
		tekstDoKolumny(a.Kolor), rodzajKotwicy, liczbaDoKolumny(a.Kotwica.Strona),
		liczbaDoKolumny(a.Kotwica.Od), liczbaDoKolumny(a.Kotwica.Do),
		tekstDoKolumny(a.Kotwica.Selektor), liczbaDoKolumny(a.Kotwica.CzasMs),
		tekstDoKolumny(a.UstalenieKod), id, KontoOperatora(ctx), KontoOperatora(ctx))
	if err != nil {
		return AdnotacjaBadania{}, err
	}
	return r.Adnotacja(ctx, a.Kod)
}

// Wypis kontraktu `ResearchExcerpt` niesie tytuł źródła, więc odczyt adnotacji
// bierze go jednym złączeniem.
const zapytanieAdnotacjiBadania = `SELECT a.identyfikator_zewnetrzny, z.identyfikator_zewnetrzny, z.tytul,
	        a.rodzaj, a.cytat, a.komentarz, a.kolor, a.kotwica_rodzaj, a.kotwica_strona,
	        a.kotwica_od, a.kotwica_do, a.kotwica_selektor, a.kotwica_czas_ms,
	        a.ustalenie_kod, a.utworzono
	   FROM adnotacja_badania a JOIN zrodlo_badania z ON z.id = a.zrodlo_id `

// warunekZrodlaAdnotacji sięga konta przez `zrodlo_badania` (konto_id od migracji 484);
// nazwy stoją niekwalifikowane, bo ten sam warunek stoi w DO UPDATE bez aliasu.
const warunekZrodlaAdnotacji = `EXISTS (SELECT 1 FROM zrodlo_badania w
	        WHERE w.id = zrodlo_id AND ` + WarunekKonta + `)`

func odczytajAdnotacjeBadania(wiersz skaner) (AdnotacjaBadania, error) {
	var a AdnotacjaBadania
	var cytat, komentarz, kolor, selektor, ustalenie sql.NullString
	var strona, od, doZnaku, czas sql.NullInt64
	err := wiersz.Scan(&a.Kod, &a.ZrodloKod, &a.ZrodloTytul, &a.Rodzaj, &cytat, &komentarz,
		&kolor, &a.Kotwica.Rodzaj, &strona, &od, &doZnaku, &selektor, &czas, &ustalenie, &a.Utworzono)
	if err != nil {
		return AdnotacjaBadania{}, err
	}
	a.Cytat = tekstZKolumny(cytat)
	a.Komentarz = tekstZKolumny(komentarz)
	a.Kolor = tekstZKolumny(kolor)
	a.Kotwica.Strona = liczbaZKolumny(strona)
	a.Kotwica.Od = liczbaZKolumny(od)
	a.Kotwica.Do = liczbaZKolumny(doZnaku)
	a.Kotwica.Selektor = tekstZKolumny(selektor)
	a.Kotwica.CzasMs = liczbaZKolumny(czas)
	a.UstalenieKod = tekstZKolumny(ustalenie)
	return a, nil
}

func (r *repozytoriumBadan) Adnotacja(ctx context.Context, kod string) (AdnotacjaBadania, error) {
	polecenie, err := r.zapytania.przygotuj(ctx,
		zapytanieAdnotacjiBadania+`WHERE a.identyfikator_zewnetrzny = ? AND `+warunekZrodlaAdnotacji)
	if err != nil {
		return AdnotacjaBadania{}, err
	}
	a, err := odczytajAdnotacjeBadania(polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return AdnotacjaBadania{}, ErrBrakWiersza
	}
	if err != nil {
		return AdnotacjaBadania{}, fmt.Errorf("dane: nieczytelna adnotacja %q: %w", kod, err)
	}
	return a, nil
}

// Adnotacje oddaje adnotacje źródła albo całego okna badania. Puste wskazanie
// jest zawężeniem pominiętym, nie zawężeniem do wartości pustej.
func (r *repozytoriumBadan) Adnotacje(ctx context.Context, kodZrodla, okno, rodzaj string,
	limit int) ([]AdnotacjaBadania, error) {

	warunki := []string{warunekZrodlaAdnotacji}
	argumenty := []any{KontoOperatora(ctx)}
	if strings.TrimSpace(kodZrodla) != "" {
		warunki = append(warunki, "z.identyfikator_zewnetrzny = ?")
		argumenty = append(argumenty, kodZrodla)
	}
	if strings.TrimSpace(okno) != "" {
		warunki = append(warunki, "z.okno = ?")
		argumenty = append(argumenty, okno)
	}
	if strings.TrimSpace(rodzaj) != "" {
		warunki = append(warunki, "a.rodzaj = ?")
		argumenty = append(argumenty, rodzaj)
	}
	if limit <= 0 {
		limit = 500
	}
	argumenty = append(argumenty, limit)

	wiersze, err := r.pytajBadania(ctx, zapytanieAdnotacjiBadania+"WHERE "+strings.Join(warunki, " AND ")+
		" ORDER BY a.utworzono DESC, a.id DESC LIMIT ?", argumenty...)
	if err != nil {
		return nil, err
	}
	defer wiersze.Close()

	lista := []AdnotacjaBadania{}
	for wiersze.Next() {
		a, err := odczytajAdnotacjeBadania(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz adnotacji: %w", err)
		}
		lista = append(lista, a)
	}
	return lista, wiersze.Err()
}

// UsunAdnotacje zdejmuje adnotację. Brak wiersza jest ErrBrakWiersza, bo okno
// ma odróżnić „nie było czego zdjąć" od „zdjęto".
func (r *repozytoriumBadan) UsunAdnotacje(ctx context.Context, kod string) error {
	polecenie, err := r.zapytania.przygotuj(ctx,
		`DELETE FROM adnotacja_badania WHERE identyfikator_zewnetrzny = ? AND `+warunekZrodlaAdnotacji)
	if err != nil {
		return err
	}
	wynik, err := polecenie.ExecContext(ctx, kod, KontoOperatora(ctx))
	if err != nil {
		return fmt.Errorf("dane: nie można usunąć adnotacji %q: %w", kod, err)
	}
	ile, _ := wynik.RowsAffected()
	if ile == 0 {
		return ErrBrakWiersza
	}
	return nil
}

// ── Szczegóły ustalenia ────────────────────────────────────────────────────

// UstawSzczegolyUstalenia nadpisuje cechy dobudowane ustalenia; pole puste
// albo zerowe zostaje bez zmiany, poza wymogiem potwierdzenia.
func (r *repozytoriumBadan) UstawSzczegolyUstalenia(ctx context.Context, kod string,
	s SzczegolyUstaleniaBadania) error {

	if _, err := r.idUstalenia(ctx, kod); err != nil {
		return err
	}
	return r.wykonajBadania(ctx, `UPDATE ustalenie_badania SET
	        rodzaj = COALESCE(NULLIF(?, ''), rodzaj),
	        waga = COALESCE(NULLIF(?, ''), waga),
	        notatka = COALESCE(?, notatka),
	        wymaga_potwierdzenia = ?,
	        adnotacja_kod = COALESCE(?, adnotacja_kod),
	        kotwica_rodzaj = COALESCE(NULLIF(?, ''), kotwica_rodzaj),
	        kotwica_strona = COALESCE(?, kotwica_strona),
	        kotwica_od = COALESCE(?, kotwica_od),
	        kotwica_do = COALESCE(?, kotwica_do),
	        kotwica_selektor = COALESCE(?, kotwica_selektor),
	        kotwica_czas_ms = COALESCE(?, kotwica_czas_ms),
	        zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	    WHERE identyfikator_zewnetrzny = ? AND `+WarunekKonta,
		s.Rodzaj, s.Waga, tekstDoKolumny(s.Notatka), wartoscCalkowitaBadania(s.WymagaPotwierdzeni),
		tekstDoKolumny(s.AdnotacjaKod), s.Kotwica.Rodzaj, liczbaDoKolumny(s.Kotwica.Strona),
		liczbaDoKolumny(s.Kotwica.Od), liczbaDoKolumny(s.Kotwica.Do),
		tekstDoKolumny(s.Kotwica.Selektor), liczbaDoKolumny(s.Kotwica.CzasMs), kod,
		KontoOperatora(ctx))
}

func (r *repozytoriumBadan) SzczegolyUstaleniaBadania(ctx context.Context, kod string) (SzczegolyUstaleniaBadania, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, `SELECT rodzaj, waga, notatka, wymaga_potwierdzenia,
	        adnotacja_kod, kotwica_rodzaj, kotwica_strona, kotwica_od, kotwica_do,
	        kotwica_selektor, kotwica_czas_ms
	    FROM ustalenie_badania WHERE identyfikator_zewnetrzny = ? AND `+WarunekKonta)
	if err != nil {
		return SzczegolyUstaleniaBadania{}, err
	}
	var s SzczegolyUstaleniaBadania
	var rodzaj, waga, notatka, adnotacja, kotwicaRodzaj, selektor sql.NullString
	var wymaga int64
	var strona, od, doZnaku, czas sql.NullInt64
	err = polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)).Scan(&rodzaj, &waga, &notatka, &wymaga,
		&adnotacja, &kotwicaRodzaj, &strona, &od, &doZnaku, &selektor, &czas)
	if errors.Is(err, sql.ErrNoRows) {
		return SzczegolyUstaleniaBadania{}, ErrBrakWiersza
	}
	if err != nil {
		return SzczegolyUstaleniaBadania{}, fmt.Errorf("dane: nieczytelne szczegóły ustalenia %q: %w", kod, err)
	}
	s.Rodzaj = rodzaj.String
	s.Waga = waga.String
	s.Notatka = tekstZKolumny(notatka)
	s.WymagaPotwierdzeni = wymaga != 0
	s.AdnotacjaKod = tekstZKolumny(adnotacja)
	s.Kotwica = KotwicaBadania{
		Rodzaj: kotwicaRodzaj.String, Strona: liczbaZKolumny(strona), Od: liczbaZKolumny(od),
		Do: liczbaZKolumny(doZnaku), Selektor: tekstZKolumny(selektor), CzasMs: liczbaZKolumny(czas),
	}
	return s, nil
}

// UsunUstalenie zdejmuje ustalenie i oddaje liczbę sekcji raportu, które przez
// to straciły odwołanie — odpowiedź kontraktu mówi o skutku, nie o zamiarze.
func (r *repozytoriumBadan) UsunUstalenie(ctx context.Context, kod string) (int, error) {
	id, err := r.idUstalenia(ctx, kod)
	if err != nil {
		return 0, err
	}
	var odwiazane int
	err = wTransakcji(ctx, r.db, func(t *sql.Tx) error {
		liczenie, err := r.zapytania.wTransakcji(ctx, t,
			`SELECT COUNT(*) FROM ustalenie_sekcji_raportu_badania WHERE ustalenie_id = ?`)
		if err != nil {
			return err
		}
		if err := liczenie.QueryRowContext(ctx, id).Scan(&odwiazane); err != nil {
			return fmt.Errorf("dane: nie można policzyć sekcji ustalenia %q: %w", kod, err)
		}
		usuniecie, err := r.zapytania.wTransakcji(ctx, t,
			`DELETE FROM ustalenie_badania WHERE id = ? AND `+WarunekKonta)
		if err != nil {
			return err
		}
		if _, err := usuniecie.ExecContext(ctx, id, KontoOperatora(ctx)); err != nil {
			return fmt.Errorf("dane: nie można usunąć ustalenia %q: %w", kod, err)
		}
		return nil
	})
	return odwiazane, err
}

func (r *repozytoriumBadan) PrzeniesZrodlaUstalen(ctx context.Context, kodyUstalen []string,
	kodDocelowy string) (int, error) {

	docelowe, err := r.idUstalenia(ctx, kodDocelowy)
	if err != nil {
		return 0, err
	}
	przeniesione := 0
	err = wTransakcji(ctx, r.db, func(t *sql.Tx) error {
		for _, kod := range kodyUstalen {
			if kod == kodDocelowy {
				continue
			}
			wskazanie, err := r.zapytania.wTransakcji(ctx, t,
				`SELECT id FROM ustalenie_badania WHERE identyfikator_zewnetrzny = ? AND `+WarunekKonta)
			if err != nil {
				return err
			}
			var scalane int64
			err = wskazanie.QueryRowContext(ctx, kod, KontoOperatora(ctx)).Scan(&scalane)
			if errors.Is(err, sql.ErrNoRows) {
				continue
			}
			if err != nil {
				return fmt.Errorf("dane: nie można odnaleźć ustalenia scalanego %q: %w", kod, err)
			}
			przepiecie, err := r.zapytania.wTransakcji(ctx, t,
				`INSERT OR IGNORE INTO zrodlo_ustalenia_badania (ustalenie_id, zrodlo_id)
				 SELECT ?, zrodlo_id FROM zrodlo_ustalenia_badania WHERE ustalenie_id = ?`)
			if err != nil {
				return err
			}
			wynik, err := przepiecie.ExecContext(ctx, docelowe, scalane)
			if err != nil {
				return fmt.Errorf("dane: nie można przepiąć źródeł ustalenia %q: %w", kod, err)
			}
			ile, _ := wynik.RowsAffected()
			przeniesione += int(ile)

			usuniecie, err := r.zapytania.wTransakcji(ctx, t,
				`DELETE FROM ustalenie_badania WHERE id = ? AND `+WarunekKonta)
			if err != nil {
				return err
			}
			if _, err := usuniecie.ExecContext(ctx, scalane, KontoOperatora(ctx)); err != nil {
				return fmt.Errorf("dane: nie można usunąć ustalenia scalonego %q: %w", kod, err)
			}
		}
		return nil
	})
	return przeniesione, err
}

// ── Książka kodów i macierz ────────────────────────────────────────────────

func (r *repozytoriumBadan) ZapiszKsiazkeKodow(ctx context.Context, okno string,
	kody []KodBadania) ([]KodBadania, error) {

	for _, kod := range kody {
		if strings.TrimSpace(kod.Kod) == "" || strings.TrimSpace(kod.Nazwa) == "" {
			continue
		}
		err := r.zapiszWGranicyBadania(ctx, `INSERT INTO kod_badania
		    (identyfikator_zewnetrzny, okno, nazwa, opis, nadrzedny_kod, konto_id)
		    VALUES (?, ?, ?, ?, ?, `+WskazanieKonta+`)
		    ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
		        nazwa = excluded.nazwa, opis = excluded.opis,
		        nadrzedny_kod = excluded.nadrzedny_kod
		    WHERE `+WarunekKonta, "kod badania", kod.Kod,
			kod.Kod, okno, kod.Nazwa, tekstDoKolumny(kod.Opis), tekstDoKolumny(kod.Nadrzedny),
			KontoOperatora(ctx), KontoOperatora(ctx))
		if err != nil {
			return nil, err
		}
	}
	return r.KsiazkaKodow(ctx, okno)
}

func (r *repozytoriumBadan) KsiazkaKodow(ctx context.Context, okno string) ([]KodBadania, error) {
	wiersze, err := r.pytajBadania(ctx, `SELECT k.identyfikator_zewnetrzny, k.nazwa, k.opis, k.nadrzedny_kod,
	        (SELECT COUNT(*) FROM kod_ustalenia_badania ku
	          WHERE ku.kod_id = k.id
	            AND EXISTS (SELECT 1 FROM ustalenie_badania w
	                         WHERE w.id = ku.ustalenie_id AND `+WarunekKonta+`))
	   FROM kod_badania k WHERE k.okno = ? AND `+WarunekKonta+` ORDER BY k.nazwa`,
		KontoOperatora(ctx), okno, KontoOperatora(ctx))
	if err != nil {
		return nil, err
	}
	defer wiersze.Close()

	lista := []KodBadania{}
	for wiersze.Next() {
		kod := KodBadania{Okno: okno}
		var opis, nadrzedny sql.NullString
		if err := wiersze.Scan(&kod.Kod, &kod.Nazwa, &opis, &nadrzedny, &kod.Wystapienia); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz książki kodów: %w", err)
		}
		kod.Opis = tekstZKolumny(opis)
		kod.Nadrzedny = tekstZKolumny(nadrzedny)
		lista = append(lista, kod)
	}
	return lista, wiersze.Err()
}

func (r *repozytoriumBadan) UstawKodyUstalenia(ctx context.Context, kodUstalenia string,
	kodyKodow []string) error {

	id, err := r.idUstalenia(ctx, kodUstalenia)
	if err != nil {
		return err
	}
	return wTransakcji(ctx, r.db, func(t *sql.Tx) error {
		czyszczenie, err := r.zapytania.wTransakcji(ctx, t,
			`DELETE FROM kod_ustalenia_badania WHERE ustalenie_id = ?`)
		if err != nil {
			return err
		}
		if _, err := czyszczenie.ExecContext(ctx, id); err != nil {
			return fmt.Errorf("dane: nie można wyczyścić kodów ustalenia %q: %w", kodUstalenia, err)
		}
		wstawienie, err := r.zapytania.wTransakcji(ctx, t,
			`INSERT OR IGNORE INTO kod_ustalenia_badania (ustalenie_id, kod_id)
			 SELECT ?, id FROM kod_badania WHERE identyfikator_zewnetrzny = ? AND `+WarunekKonta)
		if err != nil {
			return err
		}
		for _, kod := range kodyKodow {
			if _, err := wstawienie.ExecContext(ctx, id, kod, KontoOperatora(ctx)); err != nil {
				return fmt.Errorf("dane: nie można przypisać kodu %q: %w", kod, err)
			}
		}
		return nil
	})
}

func (r *repozytoriumBadan) KodyUstalenia(ctx context.Context, kodUstalenia string) ([]KodBadania, error) {
	// Warunek konta idzie podzapytaniem, bo `kod_badania` niesie własną kolumnę
	// `konto_id` i nazwa niekwalifikowana w złączeniu byłaby dwuznaczna.
	wiersze, err := r.pytajBadania(ctx, `SELECT k.identyfikator_zewnetrzny, k.okno, k.nazwa, k.opis, k.nadrzedny_kod
	   FROM kod_badania k
	   JOIN kod_ustalenia_badania ku ON ku.kod_id = k.id
	   JOIN ustalenie_badania u ON u.id = ku.ustalenie_id
	  WHERE u.identyfikator_zewnetrzny = ?
	    AND EXISTS (SELECT 1 FROM ustalenie_badania w WHERE w.id = u.id AND `+WarunekKonta+`)
	  ORDER BY k.nazwa`, kodUstalenia, KontoOperatora(ctx))
	if err != nil {
		return nil, err
	}
	defer wiersze.Close()

	lista := []KodBadania{}
	for wiersze.Next() {
		var kod KodBadania
		var opis, nadrzedny sql.NullString
		if err := wiersze.Scan(&kod.Kod, &kod.Okno, &kod.Nazwa, &opis, &nadrzedny); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz kodu ustalenia: %w", err)
		}
		kod.Opis = tekstZKolumny(opis)
		kod.Nadrzedny = tekstZKolumny(nadrzedny)
		lista = append(lista, kod)
	}
	return lista, wiersze.Err()
}

// MacierzKodow liczy zapytaniem, nie pętlą po ustaleniach: baza ma indeksy,
// a pętla miałaby tyle zapytań, ile par kod × źródło.
func (r *repozytoriumBadan) MacierzKodow(ctx context.Context, okno string) ([]KomorkaMacierzyBadania, error) {
	wiersze, err := r.pytajBadania(ctx, `SELECT k.identyfikator_zewnetrzny, z.identyfikator_zewnetrzny,
	        COUNT(DISTINCT u.id), GROUP_CONCAT(DISTINCT u.identyfikator_zewnetrzny)
	   FROM kod_badania k
	   JOIN kod_ustalenia_badania ku ON ku.kod_id = k.id
	   JOIN ustalenie_badania u ON u.id = ku.ustalenie_id
	   JOIN zrodlo_ustalenia_badania zu ON zu.ustalenie_id = u.id
	   JOIN zrodlo_badania z ON z.id = zu.zrodlo_id
	  WHERE k.okno = ? AND u.okno = ?
	    AND EXISTS (SELECT 1 FROM ustalenie_badania w WHERE w.id = u.id AND `+WarunekKonta+`)
	  GROUP BY k.identyfikator_zewnetrzny, z.identyfikator_zewnetrzny
	  ORDER BY k.identyfikator_zewnetrzny, z.identyfikator_zewnetrzny`, okno, okno, KontoOperatora(ctx))
	if err != nil {
		return nil, err
	}
	defer wiersze.Close()

	lista := []KomorkaMacierzyBadania{}
	for wiersze.Next() {
		var komorka KomorkaMacierzyBadania
		var kody sql.NullString
		if err := wiersze.Scan(&komorka.KodKodu, &komorka.ZrodloKod, &komorka.Liczba, &kody); err != nil {
			return nil, fmt.Errorf("dane: nieczytelna komórka macierzy kodów: %w", err)
		}
		if kody.Valid && kody.String != "" {
			komorka.UstalenieKody = strings.Split(kody.String, ",")
		} else {
			komorka.UstalenieKody = []string{}
		}
		lista = append(lista, komorka)
	}
	return lista, wiersze.Err()
}

// ── Sprzeczności ───────────────────────────────────────────────────────────

func (r *repozytoriumBadan) ZapiszSprzecznosc(ctx context.Context,
	s SprzecznoscBadania) (SprzecznoscBadania, error) {

	err := wTransakcji(ctx, r.db, func(t *sql.Tx) error {
		zapis, err := r.zapytania.wTransakcji(ctx, t, `INSERT INTO sprzecznosc_badania
		    (identyfikator_zewnetrzny, okno, streszczenie, roznica_liczbowa,
		     rozstrzygnieta, ustalenie_rozstrzygajace, uzasadnienie, konto_id)
		    VALUES (?, ?, ?, ?, ?, ?, ?, `+WskazanieKonta+`)
		    ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
		        streszczenie = excluded.streszczenie,
		        roznica_liczbowa = excluded.roznica_liczbowa,
		        rozstrzygnieta = excluded.rozstrzygnieta,
		        ustalenie_rozstrzygajace = excluded.ustalenie_rozstrzygajace,
		        uzasadnienie = excluded.uzasadnienie
		    WHERE `+WarunekKonta)
		if err != nil {
			return err
		}
		wynik, err := zapis.ExecContext(ctx, s.Kod, s.Okno, s.Streszczenie,
			tekstDoKolumny(s.RoznicaLiczbowa), wartoscCalkowitaBadania(s.Rozstrzygnieta),
			tekstDoKolumny(s.UstalenieRozstrzyga), tekstDoKolumny(s.Uzasadnienie),
			KontoOperatora(ctx), KontoOperatora(ctx))
		if err != nil {
			return fmt.Errorf("dane: nie można zapisać sprzeczności %q: %w", s.Kod, err)
		}
		if err := sprawdzTrafienieZapisu(wynik, "sprzeczność badania", s.Kod); err != nil {
			return err
		}
		if s.UstalenieKody == nil {
			return nil
		}
		wskazanie, err := r.zapytania.wTransakcji(ctx, t,
			`SELECT id FROM sprzecznosc_badania WHERE identyfikator_zewnetrzny = ? AND `+WarunekKonta)
		if err != nil {
			return err
		}
		var id int64
		if err := wskazanie.QueryRowContext(ctx, s.Kod, KontoOperatora(ctx)).Scan(&id); err != nil {
			return fmt.Errorf("dane: nie można odnaleźć sprzeczności %q: %w", s.Kod, err)
		}
		czyszczenie, err := r.zapytania.wTransakcji(ctx, t,
			`DELETE FROM ustalenie_sprzecznosci_badania WHERE sprzecznosc_id = ?`)
		if err != nil {
			return err
		}
		if _, err := czyszczenie.ExecContext(ctx, id); err != nil {
			return fmt.Errorf("dane: nie można wyczyścić ustaleń sprzeczności %q: %w", s.Kod, err)
		}
		wstawienie, err := r.zapytania.wTransakcji(ctx, t,
			`INSERT OR IGNORE INTO ustalenie_sprzecznosci_badania (sprzecznosc_id, ustalenie_kod)
			 VALUES (?, ?)`)
		if err != nil {
			return err
		}
		for _, kod := range s.UstalenieKody {
			if _, err := wstawienie.ExecContext(ctx, id, kod); err != nil {
				return fmt.Errorf("dane: nie można powiązać ustalenia %q ze sprzecznością: %w", kod, err)
			}
		}
		return nil
	})
	if err != nil {
		return SprzecznoscBadania{}, err
	}
	return r.Sprzecznosc(ctx, s.Kod)
}

func (r *repozytoriumBadan) Sprzecznosc(ctx context.Context, kod string) (SprzecznoscBadania, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, `SELECT identyfikator_zewnetrzny, okno, streszczenie,
	        roznica_liczbowa, rozstrzygnieta, ustalenie_rozstrzygajace, uzasadnienie
	    FROM sprzecznosc_badania WHERE identyfikator_zewnetrzny = ? AND `+WarunekKonta)
	if err != nil {
		return SprzecznoscBadania{}, err
	}
	s, err := odczytajSprzecznoscBadania(polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return SprzecznoscBadania{}, ErrBrakWiersza
	}
	if err != nil {
		return SprzecznoscBadania{}, fmt.Errorf("dane: nieczytelna sprzeczność %q: %w", kod, err)
	}
	kody, err := r.ustaleniaSprzecznosci(ctx, kod)
	if err != nil {
		return SprzecznoscBadania{}, err
	}
	s.UstalenieKody = kody
	return s, nil
}

func (r *repozytoriumBadan) Sprzecznosci(ctx context.Context, okno string) ([]SprzecznoscBadania, error) {
	wiersze, err := r.pytajBadania(ctx, `SELECT identyfikator_zewnetrzny, okno, streszczenie,
	        roznica_liczbowa, rozstrzygnieta, ustalenie_rozstrzygajace, uzasadnienie
	    FROM sprzecznosc_badania WHERE okno = ? AND `+WarunekKonta+`
	    ORDER BY utworzono DESC, id DESC`, okno, KontoOperatora(ctx))
	if err != nil {
		return nil, err
	}
	defer wiersze.Close()

	lista := []SprzecznoscBadania{}
	for wiersze.Next() {
		s, err := odczytajSprzecznoscBadania(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz sprzeczności: %w", err)
		}
		lista = append(lista, s)
	}
	if err := wiersze.Err(); err != nil {
		return nil, err
	}
	for numer := range lista {
		kody, err := r.ustaleniaSprzecznosci(ctx, lista[numer].Kod)
		if err != nil {
			return nil, err
		}
		lista[numer].UstalenieKody = kody
	}
	return lista, nil
}

func (r *repozytoriumBadan) ustaleniaSprzecznosci(ctx context.Context, kod string) ([]string, error) {
	wiersze, err := r.pytajBadania(ctx, `SELECT us.ustalenie_kod
	   FROM ustalenie_sprzecznosci_badania us
	   JOIN sprzecznosc_badania s ON s.id = us.sprzecznosc_id
	  WHERE s.identyfikator_zewnetrzny = ?
	    AND EXISTS (SELECT 1 FROM sprzecznosc_badania w WHERE w.id = s.id AND `+WarunekKonta+`)
	  ORDER BY us.ustalenie_kod`, kod, KontoOperatora(ctx))
	if err != nil {
		return nil, err
	}
	return napisyZWierszyBadania(wiersze)
}

func odczytajSprzecznoscBadania(wiersz skaner) (SprzecznoscBadania, error) {
	var s SprzecznoscBadania
	var roznica, rozstrzygajace, uzasadnienie sql.NullString
	var rozstrzygnieta int64
	err := wiersz.Scan(&s.Kod, &s.Okno, &s.Streszczenie, &roznica, &rozstrzygnieta,
		&rozstrzygajace, &uzasadnienie)
	if err != nil {
		return SprzecznoscBadania{}, err
	}
	s.RoznicaLiczbowa = tekstZKolumny(roznica)
	s.Rozstrzygnieta = rozstrzygnieta != 0
	s.UstalenieRozstrzyga = tekstZKolumny(rozstrzygajace)
	s.Uzasadnienie = tekstZKolumny(uzasadnienie)
	s.UstalenieKody = []string{}
	return s, nil
}

// ── Weryfikacja twierdzenia ────────────────────────────────────────────────

func (r *repozytoriumBadan) ZapiszWeryfikacje(ctx context.Context, kodUstalenia,
	werdykt, uzasadnienie string) (string, error) {

	id, err := r.idUstalenia(ctx, kodUstalenia)
	if err != nil {
		return "", err
	}
	err = r.wykonajBadania(ctx, `INSERT INTO weryfikacja_ustalenia_badania
	    (ustalenie_id, werdykt, uzasadnienie, sprawdzono)
	    VALUES (?, ?, ?, strftime('%Y-%m-%dT%H:%M:%fZ','now'))
	    ON CONFLICT(ustalenie_id) DO UPDATE SET
	        werdykt = excluded.werdykt, uzasadnienie = excluded.uzasadnienie,
	        sprawdzono = excluded.sprawdzono`, id, werdykt, uzasadnienie)
	if err != nil {
		return "", err
	}
	polecenie, err := r.zapytania.przygotuj(ctx,
		`SELECT sprawdzono FROM weryfikacja_ustalenia_badania WHERE ustalenie_id = ?`)
	if err != nil {
		return "", err
	}
	var sprawdzono string
	if err := polecenie.QueryRowContext(ctx, id).Scan(&sprawdzono); err != nil {
		return "", fmt.Errorf("dane: nieczytelna weryfikacja ustalenia %q: %w", kodUstalenia, err)
	}
	return sprawdzono, nil
}

// ── Wątki tematyczne ───────────────────────────────────────────────────────

// ZapiszWatki wymienia w całości wątki okna badania. Grupowanie jest wynikiem
// jednego przebiegu klastrowania, więc dopisywanie zostawiłoby wątki poprzednie
// obok nowych — dwa podziały tego samego zbioru naraz.
func (r *repozytoriumBadan) ZapiszWatki(ctx context.Context, okno string,
	watki []WatekBadania) ([]WatekBadania, error) {

	err := wTransakcji(ctx, r.db, func(t *sql.Tx) error {
		czyszczenie, err := r.zapytania.wTransakcji(ctx, t,
			`DELETE FROM watek_ustalen_badania WHERE okno = ? AND `+WarunekKonta)
		if err != nil {
			return err
		}
		if _, err := czyszczenie.ExecContext(ctx, okno, KontoOperatora(ctx)); err != nil {
			return fmt.Errorf("dane: nie można wyczyścić wątków okna %q: %w", okno, err)
		}
		zapis, err := r.zapytania.wTransakcji(ctx, t, `INSERT INTO watek_ustalen_badania
		    (identyfikator_zewnetrzny, okno, nazwa, automatyczny, konto_id)
		    VALUES (?, ?, ?, ?, `+WskazanieKonta+`)`)
		if err != nil {
			return err
		}
		wiazanie, err := r.zapytania.wTransakcji(ctx, t,
			`INSERT OR IGNORE INTO ustalenie_watku_badania (watek_id, ustalenie_kod) VALUES (?, ?)`)
		if err != nil {
			return err
		}
		for _, watek := range watki {
			wynik, err := zapis.ExecContext(ctx, watek.Kod, okno, watek.Nazwa,
				wartoscCalkowitaBadania(watek.Automatyczny), KontoOperatora(ctx))
			if err != nil {
				return fmt.Errorf("dane: nie można zapisać wątku %q: %w", watek.Kod, err)
			}
			id, err := wynik.LastInsertId()
			if err != nil {
				return fmt.Errorf("dane: nie można odczytać klucza wątku %q: %w", watek.Kod, err)
			}
			for _, kod := range watek.UstalenieKody {
				if _, err := wiazanie.ExecContext(ctx, id, kod); err != nil {
					return fmt.Errorf("dane: nie można powiązać ustalenia %q z wątkiem: %w", kod, err)
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return watki, nil
}

// ── Ślad prowenancji ───────────────────────────────────────────────────────

// ZapiszProwenancje dopisuje wpis do śladu ustalenia. Ślad narasta i nie jest
// nadpisywany — historia zmian, którą da się nadpisać, historią nie jest.
func (r *repozytoriumBadan) ZapiszProwenancje(ctx context.Context, w WpisProwenancjiBadania) error {
	aktor := w.Aktor
	if aktor == "" {
		aktor = "operator"
	}
	return r.wykonajBadania(ctx, `INSERT INTO prowenancja_badania
	    (ustalenie_kod, aktor, czynnosc, zrodlo_kod, kotwica_strona, konto_id)
	    VALUES (?, ?, ?, ?, ?, `+WskazanieKonta+`)`,
		w.UstalenieKod, aktor, w.Czynnosc, tekstDoKolumny(w.ZrodloKod), liczbaDoKolumny(w.Strona),
		KontoOperatora(ctx))
}

func (r *repozytoriumBadan) Prowenancja(ctx context.Context,
	kodUstalenia string) ([]WpisProwenancjiBadania, error) {

	wiersze, err := r.pytajBadania(ctx, `SELECT ustalenie_kod, o_czasie, aktor, czynnosc, zrodlo_kod, kotwica_strona
	    FROM prowenancja_badania WHERE ustalenie_kod = ? AND `+WarunekKonta+`
	    ORDER BY o_czasie, id`, kodUstalenia, KontoOperatora(ctx))
	if err != nil {
		return nil, err
	}
	defer wiersze.Close()

	lista := []WpisProwenancjiBadania{}
	for wiersze.Next() {
		var wpis WpisProwenancjiBadania
		var zrodlo sql.NullString
		var strona sql.NullInt64
		if err := wiersze.Scan(&wpis.UstalenieKod, &wpis.OCzasie, &wpis.Aktor, &wpis.Czynnosc,
			&zrodlo, &strona); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wpis prowenancji: %w", err)
		}
		wpis.ZrodloKod = tekstZKolumny(zrodlo)
		wpis.Strona = liczbaZKolumny(strona)
		lista = append(lista, wpis)
	}
	return lista, wiersze.Err()
}
