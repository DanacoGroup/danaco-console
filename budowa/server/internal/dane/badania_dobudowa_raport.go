// Plik dobudowuje obszar badań po stronie danych: raport i jego otoczenie —
// szablony struktury, wersje z migawką sekcji, komentarze recenzji, bloki
// wstawek, ustawienia przypisów oraz eksport. Kontrakt obszaru deklaruje plik
// badania_dobudowa.go.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// SzablonRaportuBadania odwzorowuje wiersz tabeli szablon_raportu_badania: nazwany zestaw
// tytułów sekcji, który operator może zastosować przy zakładaniu nowego raportu badania.
type SzablonRaportuBadania struct {
	Kod          string
	Nazwa        string
	TytulySekcji []string
	Wlasny       bool
}

// WersjaRaportuBadania odwzorowuje wiersz tabeli wersja_raportu_badania. Pole Migawka niesie
// pełny stan sekcji raportu zapisany w chwili utworzenia wersji, niezależnie od
// późniejszych zmian samego raportu.
type WersjaRaportuBadania struct {
	Kod          string
	RaportKod    string
	Etykieta     *string
	Migawka      string
	LiczbaSekcji int
	Utworzono    string
}

// KomentarzRaportuBadania odwzorowuje wiersz tabeli komentarz_raportu_badania: wpis recenzji
// przypisany do sekcji albo wątku raportu, z opcjonalnym cytatem fragmentu tekstu i stanem rozstrzygnięcia.
type KomentarzRaportuBadania struct {
	Kod            string
	RaportKod      string
	SekcjaKod      *string
	WatekKod       *string
	Tresc          string
	Cytat          *string
	Rozstrzygniety bool
	Utworzono      string
}

// BlokRaportuBadania to wiersz `blok_raportu_badania` — wstawka osadzona
// w sekcji: macierz, oś czasu, wykres albo tabela dowodów.
type BlokRaportuBadania struct {
	Kod           string
	RaportKod     string
	SekcjaKod     string
	Rodzaj        string
	Podpis        *string
	Naglowki      []string
	Wiersze       string
	UstalenieKody []string
}

// SzablonEksportuBadania odwzorowuje wiersz tabeli szablon_eksportu_badania: zapisaną
// kombinację formatu i zawartości, którą operator może ponownie zastosować przy kolejnym eksporcie raportu.
type SzablonEksportuBadania struct {
	Kod       string
	Nazwa     string
	Format    string
	Cel       *string
	Zawartosc string
}

// ── Raporty okna ───────────────────────────────────────────────────────────

// Raporty oddaje raporty okna badania od najświeżej zmienionego. Potrzebne
// `research.report.get` bez wskazania raportu: okno pyta o raport bieżący,
// a bieżącym jest ten zmieniony ostatnio.
func (r *repozytoriumBadan) Raporty(ctx context.Context, okno string) ([]RaportBadania, error) {
	wiersze, err := r.pytajBadania(ctx, `SELECT `+kolumnyRaportuBadania+`
	    FROM raport_badania WHERE okno = ? ORDER BY zaktualizowano DESC, id DESC`, okno)
	if err != nil {
		return nil, err
	}
	defer wiersze.Close()

	lista := []RaportBadania{}
	for wiersze.Next() {
		var raport RaportBadania
		if err := wiersze.Scan(&raport.ID, &raport.Kod, &raport.Okno, &raport.Tytul,
			&raport.Utworzono, &raport.Zaktualizowano); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz raportu badania: %w", err)
		}
		lista = append(lista, raport)
	}
	return lista, wiersze.Err()
}

// ── Szablony struktury raportu ─────────────────────────────────────────────

// ZapiszSzablonRaportu utrwala szablon własny operatora w tabeli szablon_raportu_badania,
// nadpisując nazwę i tytuły sekcji istniejącego wiersza o tym samym identyfikatorze zewnętrznym.
func (r *repozytoriumBadan) ZapiszSzablonRaportu(ctx context.Context, s SzablonRaportuBadania) error {
	return r.wykonajBadania(ctx, `INSERT INTO szablon_raportu_badania
	    (identyfikator_zewnetrzny, nazwa, tytuly_sekcji, wlasny) VALUES (?, ?, ?, ?)
	    ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	        nazwa = excluded.nazwa, tytuly_sekcji = excluded.tytuly_sekcji`,
		s.Kod, s.Nazwa, listaJakoBadania(s.TytulySekcji), wartoscCalkowitaBadania(s.Wlasny))
}

// SzablonyRaportu oddaje szablony zapisane w instalacji. Zestaw wbudowany
// dokłada adapter — repozytorium mówi wyłącznie o tym, co leży w bazie.
func (r *repozytoriumBadan) SzablonyRaportu(ctx context.Context) ([]SzablonRaportuBadania, error) {
	wiersze, err := r.pytajBadania(ctx,
		`SELECT identyfikator_zewnetrzny, nazwa, tytuly_sekcji, wlasny
		 FROM szablon_raportu_badania ORDER BY nazwa`)
	if err != nil {
		return nil, err
	}
	defer wiersze.Close()

	lista := []SzablonRaportuBadania{}
	for wiersze.Next() {
		var szablon SzablonRaportuBadania
		var tytuly string
		var wlasny int64
		if err := wiersze.Scan(&szablon.Kod, &szablon.Nazwa, &tytuly, &wlasny); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz szablonu raportu: %w", err)
		}
		szablon.TytulySekcji = listaZBadania(tytuly)
		szablon.Wlasny = wlasny != 0
		lista = append(lista, szablon)
	}
	return lista, wiersze.Err()
}

// ── Wersje raportu ─────────────────────────────────────────────────────────

// ZapiszWersjeRaportu utrwala migawkę aktualnego stanu sekcji raportu jako kolejną,
// niemodyfikowalną wersję powiązaną z raportem wskazanym kodem.
func (r *repozytoriumBadan) ZapiszWersjeRaportu(ctx context.Context,
	w WersjaRaportuBadania) (WersjaRaportuBadania, error) {

	id, err := r.idRaportu(ctx, w.RaportKod)
	if err != nil {
		return WersjaRaportuBadania{}, err
	}
	err = r.wykonajBadania(ctx, `INSERT INTO wersja_raportu_badania
	    (identyfikator_zewnetrzny, raport_id, etykieta, migawka, liczba_sekcji)
	    VALUES (?, ?, ?, ?, ?)`,
		w.Kod, id, tekstDoKolumny(w.Etykieta), w.Migawka, w.LiczbaSekcji)
	if err != nil {
		return WersjaRaportuBadania{}, err
	}
	return r.WersjaRaportuBadania(ctx, w.Kod)
}

// zapytanieWersjiBadania jest wspólnym tekstem zapytania SQL, z którego korzystają wszystkie
// funkcje odczytujące pojedynczą albo wiele wersji raportu.
const zapytanieWersjiBadania = `SELECT w.identyfikator_zewnetrzny, r.identyfikator_zewnetrzny,
	        w.etykieta, w.migawka, w.liczba_sekcji, w.utworzono
	   FROM wersja_raportu_badania w JOIN raport_badania r ON r.id = w.raport_id `

// WersjaRaportuBadania oddaje jedną wersję raportu wskazaną kodem zewnętrznym albo błąd
// ErrBrakWiersza, gdy taka wersja nie istnieje w bazie.
func (r *repozytoriumBadan) WersjaRaportuBadania(ctx context.Context, kod string) (WersjaRaportuBadania, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, zapytanieWersjiBadania+`WHERE w.identyfikator_zewnetrzny = ?`)
	if err != nil {
		return WersjaRaportuBadania{}, err
	}
	w, err := odczytajWersjeRaportuBadania(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return WersjaRaportuBadania{}, ErrBrakWiersza
	}
	if err != nil {
		return WersjaRaportuBadania{}, fmt.Errorf("dane: nieczytelna wersja raportu %q: %w", kod, err)
	}
	return w, nil
}

// WersjeRaportu oddaje wersje wskazanego raportu od najświeższej, ograniczone parametrem
// limit; wartość niedodatnia ustala limit domyślny.
func (r *repozytoriumBadan) WersjeRaportu(ctx context.Context, kodRaportu string,
	limit int) ([]WersjaRaportuBadania, error) {

	if limit <= 0 {
		limit = 50
	}
	wiersze, err := r.pytajBadania(ctx, zapytanieWersjiBadania+
		`WHERE r.identyfikator_zewnetrzny = ? ORDER BY w.utworzono DESC, w.id DESC LIMIT ?`,
		kodRaportu, limit)
	if err != nil {
		return nil, err
	}
	defer wiersze.Close()

	lista := []WersjaRaportuBadania{}
	for wiersze.Next() {
		w, err := odczytajWersjeRaportuBadania(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz wersji raportu: %w", err)
		}
		lista = append(lista, w)
	}
	return lista, wiersze.Err()
}

// odczytajWersjeRaportuBadania składa strukturę WersjaRaportuBadania z jednego wiersza wyniku
// zapytania, niezależnie od tego, czy wiersz pochodzi z pojedynczego odczytu, czy z iteracji po wielu wierszach.
func odczytajWersjeRaportuBadania(wiersz skaner) (WersjaRaportuBadania, error) {
	var w WersjaRaportuBadania
	var etykieta sql.NullString
	err := wiersz.Scan(&w.Kod, &w.RaportKod, &etykieta, &w.Migawka, &w.LiczbaSekcji, &w.Utworzono)
	if err != nil {
		return WersjaRaportuBadania{}, err
	}
	w.Etykieta = tekstZKolumny(etykieta)
	return w, nil
}

// ── Komentarze recenzji ────────────────────────────────────────────────────

// ZapiszKomentarzRaportu zakłada nowy komentarz recenzji albo nadpisuje zastany wiersz o tym
// samym identyfikatorze zewnętrznym, zachowując datę jego utworzenia.
func (r *repozytoriumBadan) ZapiszKomentarzRaportu(ctx context.Context,
	k KomentarzRaportuBadania) (KomentarzRaportuBadania, error) {

	id, err := r.idRaportu(ctx, k.RaportKod)
	if err != nil {
		return KomentarzRaportuBadania{}, err
	}
	err = r.wykonajBadania(ctx, `INSERT INTO komentarz_raportu_badania
	    (identyfikator_zewnetrzny, raport_id, sekcja_kod, watek_kod, tresc, cytat, rozstrzygniety)
	    VALUES (?, ?, ?, ?, ?, ?, ?)
	    ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	        sekcja_kod = excluded.sekcja_kod, watek_kod = excluded.watek_kod,
	        tresc = excluded.tresc, cytat = excluded.cytat,
	        rozstrzygniety = excluded.rozstrzygniety`,
		k.Kod, id, tekstDoKolumny(k.SekcjaKod), tekstDoKolumny(k.WatekKod), k.Tresc,
		tekstDoKolumny(k.Cytat), wartoscCalkowitaBadania(k.Rozstrzygniety))
	if err != nil {
		return KomentarzRaportuBadania{}, err
	}
	lista, err := r.KomentarzeRaportu(ctx, k.RaportKod, false)
	if err != nil {
		return KomentarzRaportuBadania{}, err
	}
	for _, wiersz := range lista {
		if wiersz.Kod == k.Kod {
			return wiersz, nil
		}
	}
	return KomentarzRaportuBadania{}, ErrBrakWiersza
}

// KomentarzeRaportu oddaje komentarze raportu; zawężenie do otwartych jest
// pytaniem trybu recenzji, nie osobnym bytem.
func (r *repozytoriumBadan) KomentarzeRaportu(ctx context.Context, kodRaportu string,
	tylkoOtwarte bool) ([]KomentarzRaportuBadania, error) {

	sqlTekst := `SELECT k.identyfikator_zewnetrzny, r.identyfikator_zewnetrzny, k.sekcja_kod,
	        k.watek_kod, k.tresc, k.cytat, k.rozstrzygniety, k.utworzono
	   FROM komentarz_raportu_badania k JOIN raport_badania r ON r.id = k.raport_id
	  WHERE r.identyfikator_zewnetrzny = ?`
	if tylkoOtwarte {
		sqlTekst += ` AND k.rozstrzygniety = 0`
	}
	sqlTekst += ` ORDER BY k.utworzono, k.id`

	wiersze, err := r.pytajBadania(ctx, sqlTekst, kodRaportu)
	if err != nil {
		return nil, err
	}
	defer wiersze.Close()

	lista := []KomentarzRaportuBadania{}
	for wiersze.Next() {
		var k KomentarzRaportuBadania
		var sekcja, watek, cytat sql.NullString
		var rozstrzygniety int64
		if err := wiersze.Scan(&k.Kod, &k.RaportKod, &sekcja, &watek, &k.Tresc, &cytat,
			&rozstrzygniety, &k.Utworzono); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz komentarza raportu: %w", err)
		}
		k.SekcjaKod = tekstZKolumny(sekcja)
		k.WatekKod = tekstZKolumny(watek)
		k.Cytat = tekstZKolumny(cytat)
		k.Rozstrzygniety = rozstrzygniety != 0
		lista = append(lista, k)
	}
	return lista, wiersze.Err()
}

// ── Bloki wstawek ──────────────────────────────────────────────────────────

// ZapiszBlokRaportu utrwala wstawkę osadzoną w sekcji raportu — macierz, oś czasu, wykres
// albo tabelę dowodów — nadpisując zastany wiersz o tym samym kodzie.
func (r *repozytoriumBadan) ZapiszBlokRaportu(ctx context.Context,
	b BlokRaportuBadania) (BlokRaportuBadania, error) {

	id, err := r.idRaportu(ctx, b.RaportKod)
	if err != nil {
		return BlokRaportuBadania{}, err
	}
	err = r.wykonajBadania(ctx, `INSERT INTO blok_raportu_badania
	    (identyfikator_zewnetrzny, raport_id, sekcja_kod, rodzaj, podpis, naglowki, wiersze, ustalenia)
	    VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	    ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	        sekcja_kod = excluded.sekcja_kod, rodzaj = excluded.rodzaj,
	        podpis = excluded.podpis, naglowki = excluded.naglowki,
	        wiersze = excluded.wiersze, ustalenia = excluded.ustalenia`,
		b.Kod, id, b.SekcjaKod, b.Rodzaj, tekstDoKolumny(b.Podpis),
		listaJakoBadania(b.Naglowki), b.Wiersze, listaJakoBadania(b.UstalenieKody))
	if err != nil {
		return BlokRaportuBadania{}, err
	}
	return b, nil
}

// UstawPrzypisyRaportu zapisuje ustawienia menedżera przypisów raportu — umiejscowienie oraz
// tryb skrócony — i odświeża znacznik czasu aktualizacji raportu.
func (r *repozytoriumBadan) UstawPrzypisyRaportu(ctx context.Context, kodRaportu,
	umiejscowienie string, skrocone bool) error {

	if _, err := r.idRaportu(ctx, kodRaportu); err != nil {
		return err
	}
	return r.wykonajBadania(ctx, `UPDATE raport_badania
	    SET przypisy_umiejscowienie = ?, przypisy_skrocone = ?,
	        zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	    WHERE identyfikator_zewnetrzny = ?`,
		umiejscowienie, wartoscCalkowitaBadania(skrocone), kodRaportu)
}

// ── Eksport ────────────────────────────────────────────────────────────────

// Eksporty oddaje ślad eksportów raportu wskazanego kodem albo, gdy kod raportu jest pusty,
// wszystkich eksportów całego okna badania wskazanego nazwą.
func (r *repozytoriumBadan) Eksporty(ctx context.Context, kodRaportu, okno string,
	limit int) ([]EksportRaportu, error) {

	warunki := []string{"1 = 1"}
	argumenty := []any{}
	if strings.TrimSpace(kodRaportu) != "" {
		warunki = append(warunki, "r.identyfikator_zewnetrzny = ?")
		argumenty = append(argumenty, kodRaportu)
	}
	if strings.TrimSpace(okno) != "" {
		warunki = append(warunki, "r.okno = ?")
		argumenty = append(argumenty, okno)
	}
	if limit <= 0 {
		limit = 50
	}
	argumenty = append(argumenty, limit)

	wiersze, err := r.pytajBadania(ctx, `SELECT e.identyfikator_zewnetrzny, r.identyfikator_zewnetrzny,
	        e.format, e.cel, e.sciezka_docelowa, e.plik_biblioteki_id, e.sciezka_wyniku,
	        e.rozmiar_bajtow, e.utworzono
	   FROM eksport_raportu_badania e JOIN raport_badania r ON r.id = e.raport_id
	  WHERE `+strings.Join(warunki, " AND ")+
		` ORDER BY e.utworzono DESC, e.id DESC LIMIT ?`, argumenty...)
	if err != nil {
		return nil, err
	}
	defer wiersze.Close()

	lista := []EksportRaportu{}
	for wiersze.Next() {
		var eksport EksportRaportu
		var cel string
		var docelowa, plik, wynik sql.NullString
		var rozmiar sql.NullInt64
		if err := wiersze.Scan(&eksport.Kod, &eksport.RaportKod, &eksport.Format, &cel,
			&docelowa, &plik, &wynik, &rozmiar, &eksport.Utworzono); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz eksportu raportu: %w", err)
		}
		eksport.SciezkaDocelowa = tekstZKolumny(docelowa)
		eksport.PlikBibliotekiID = tekstZKolumny(plik)
		eksport.SciezkaWyniku = tekstZKolumny(wynik)
		eksport.RozmiarBajtow = liczbaZKolumny(rozmiar)
		eksport.Cel = cel
		lista = append(lista, eksport)
	}
	return lista, wiersze.Err()
}

// ZapiszSzablonEksportu utrwala kombinację formatu i zawartości eksportu jako szablon
// wielokrotnego użytku, nadpisując zastany wiersz o tym samym kodzie.
func (r *repozytoriumBadan) ZapiszSzablonEksportu(ctx context.Context,
	s SzablonEksportuBadania) (SzablonEksportuBadania, error) {

	err := r.wykonajBadania(ctx, `INSERT INTO szablon_eksportu_badania
	    (identyfikator_zewnetrzny, nazwa, format, cel, zawartosc) VALUES (?, ?, ?, ?, ?)
	    ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	        nazwa = excluded.nazwa, format = excluded.format,
	        cel = excluded.cel, zawartosc = excluded.zawartosc`,
		s.Kod, s.Nazwa, s.Format, tekstDoKolumny(s.Cel), s.Zawartosc)
	if err != nil {
		return SzablonEksportuBadania{}, err
	}
	return s, nil
}

// SzablonyEksportu oddaje wszystkie szablony eksportu zapisane w instalacji, uporządkowane
// alfabetycznie według nazwy szablonu.
func (r *repozytoriumBadan) SzablonyEksportu(ctx context.Context) ([]SzablonEksportuBadania, error) {
	wiersze, err := r.pytajBadania(ctx,
		`SELECT identyfikator_zewnetrzny, nazwa, format, cel, zawartosc
		 FROM szablon_eksportu_badania ORDER BY nazwa`)
	if err != nil {
		return nil, err
	}
	defer wiersze.Close()

	lista := []SzablonEksportuBadania{}
	for wiersze.Next() {
		var szablon SzablonEksportuBadania
		var cel sql.NullString
		if err := wiersze.Scan(&szablon.Kod, &szablon.Nazwa, &szablon.Format, &cel,
			&szablon.Zawartosc); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz szablonu eksportu: %w", err)
		}
		szablon.Cel = tekstZKolumny(cel)
		lista = append(lista, szablon)
	}
	return lista, wiersze.Err()
}

// ZapiszUdostepnienie utrwala odnośnik udostępnienia raportu wraz z adresem oraz opcjonalnym
// czasem wygaśnięcia, nadpisując zastany wiersz o tym samym kodzie.
func (r *repozytoriumBadan) ZapiszUdostepnienie(ctx context.Context, kodRaportu, kod,
	adres string, wygasaO *string) error {

	id, err := r.idRaportu(ctx, kodRaportu)
	if err != nil {
		return err
	}
	return r.wykonajBadania(ctx, `INSERT INTO udostepnienie_raportu_badania
	    (identyfikator_zewnetrzny, raport_id, adres, wygasa_o) VALUES (?, ?, ?, ?)
	    ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	        adres = excluded.adres, wygasa_o = excluded.wygasa_o`,
		kod, id, adres, tekstDoKolumny(wygasaO))
}
