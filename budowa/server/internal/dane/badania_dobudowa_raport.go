// Dobudowa obszaru badań: raport i jego otoczenie — szablony, wersje, komentarze, bloki, przypisy, eksport.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// SzablonRaportuBadania to nazwany zestaw tytułów sekcji do zastosowania przy zakładaniu raportu.
type SzablonRaportuBadania struct {
	Kod          string
	Nazwa        string
	TytulySekcji []string
	Wlasny       bool
}

// WersjaRaportuBadania niesie w Migawce pełny stan sekcji raportu z chwili utworzenia wersji.
type WersjaRaportuBadania struct {
	Kod          string
	RaportKod    string
	Etykieta     *string
	Migawka      string
	LiczbaSekcji int
	Utworzono    string
}

// KomentarzRaportuBadania to wpis recenzji przypisany do sekcji albo wątku, z cytatem i stanem rozstrzygnięcia.
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

// BlokRaportuBadania to wstawka osadzona w sekcji: macierz, oś czasu, wykres albo tabela dowodów.
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

// SzablonEksportuBadania to zapisana kombinacja formatu i zawartości do ponownego zastosowania przy eksporcie.
type SzablonEksportuBadania struct {
	Kod       string
	Nazwa     string
	Format    string
	Cel       *string
	Zawartosc string
}

// Raporty oddaje raporty okna badania od najświeżej zmienionego (raport bieżący dla `research.report.get`).
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

// ZapiszSzablonRaportu utrwala szablon własny Operatora, nadpisując nazwę i tytuły sekcji wiersza o tym kodzie.
func (r *repozytoriumBadan) ZapiszSzablonRaportu(ctx context.Context, s SzablonRaportuBadania) error {
	return r.wykonajBadania(ctx, `INSERT INTO szablon_raportu_badania
	    (identyfikator_zewnetrzny, nazwa, tytuly_sekcji, wlasny, konto_id) VALUES (?, ?, ?, ?, `+WskazanieKonta+`)
	    ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	        nazwa = excluded.nazwa, tytuly_sekcji = excluded.tytuly_sekcji
	    WHERE `+WarunekKonta,
		s.Kod, s.Nazwa, listaJakoBadania(s.TytulySekcji), wartoscCalkowitaBadania(s.Wlasny),
		KontoOperatora(ctx), KontoOperatora(ctx))
}

func (r *repozytoriumBadan) SzablonyRaportu(ctx context.Context) ([]SzablonRaportuBadania, error) {
	wiersze, err := r.pytajBadania(ctx,
		`SELECT identyfikator_zewnetrzny, nazwa, tytuly_sekcji, wlasny
		 FROM szablon_raportu_badania WHERE `+WarunekKonta+` ORDER BY nazwa`,
		KontoOperatora(ctx))
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

// ZapiszWersjeRaportu utrwala migawkę stanu sekcji raportu jako kolejną, niemodyfikowalną wersję.
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

const zapytanieWersjiBadania = `SELECT w.identyfikator_zewnetrzny, r.identyfikator_zewnetrzny,
	        w.etykieta, w.migawka, w.liczba_sekcji, w.utworzono
	   FROM wersja_raportu_badania w JOIN raport_badania r ON r.id = w.raport_id `

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

// WersjeRaportu oddaje wersje raportu od najświeższej, ograniczone limitem (niedodatni = domyślny).
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

// ZapiszKomentarzRaportu zakłada komentarz recenzji albo nadpisuje zastany po kodzie, zachowując datę utworzenia.
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

// KomentarzeRaportu oddaje komentarze raportu; zawężenie do otwartych jest pytaniem trybu recenzji.
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

// ZapiszBlokRaportu utrwala wstawkę osadzoną w sekcji (macierz, oś czasu, wykres, tabela) po kodzie.
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

// UstawPrzypisyRaportu zapisuje umiejscowienie i tryb skrócony przypisów raportu i odświeża znacznik czasu.
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

// Eksporty oddaje ślad eksportów raportu po kodzie albo, gdy kod pusty, wszystkich eksportów okna.
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

// ZapiszSzablonEksportu utrwala kombinację formatu i zawartości eksportu jako szablon, nadpisując po kodzie.
func (r *repozytoriumBadan) ZapiszSzablonEksportu(ctx context.Context,
	s SzablonEksportuBadania) (SzablonEksportuBadania, error) {

	err := r.wykonajBadania(ctx, `INSERT INTO szablon_eksportu_badania
	    (identyfikator_zewnetrzny, nazwa, format, cel, zawartosc, konto_id) VALUES (?, ?, ?, ?, ?, `+WskazanieKonta+`)
	    ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	        nazwa = excluded.nazwa, format = excluded.format,
	        cel = excluded.cel, zawartosc = excluded.zawartosc
	    WHERE `+WarunekKonta,
		s.Kod, s.Nazwa, s.Format, tekstDoKolumny(s.Cel), s.Zawartosc,
		KontoOperatora(ctx), KontoOperatora(ctx))
	if err != nil {
		return SzablonEksportuBadania{}, err
	}
	return s, nil
}

func (r *repozytoriumBadan) SzablonyEksportu(ctx context.Context) ([]SzablonEksportuBadania, error) {
	wiersze, err := r.pytajBadania(ctx,
		`SELECT identyfikator_zewnetrzny, nazwa, format, cel, zawartosc
		 FROM szablon_eksportu_badania WHERE `+WarunekKonta+` ORDER BY nazwa`,
		KontoOperatora(ctx))
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

// ZapiszUdostepnienie utrwala odnośnik udostępnienia raportu z adresem i czasem wygaśnięcia, po kodzie.
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
