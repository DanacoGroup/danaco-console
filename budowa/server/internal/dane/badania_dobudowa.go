// Dobudowa obszaru Research po stronie danych: kontrakt RepozytoriumBadanDobudowa
// w całości oraz katalogowanie, lektura, treść, załączniki, streszczenia i tabele źródła.
package dane

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type KatalogZrodlaBadania struct {
	Etykiety []string
	Kolekcje []string
	Pytania  []string
}

type LekturaZrodlaBadania struct {
	StanLektury string
	EtapIndeks  *int64
	Identyfikat *string
	CslJson     *string
}

// TrescZrodlaBadania leży przy źródle, nie w tabeli własnej (migracja 150).
type TrescZrodlaBadania struct {
	Tekst         *string
	Stron         *int64
	WarstwaTekstu bool
}

type WycofanieZrodlaBadania struct {
	Stan       string
	Adres      *string
	Sprawdzono string
}

type ZalacznikZrodlaBadania struct {
	Kod              string
	ZrodloKod        string
	Rodzaj           string
	PlikBibliotekiID *string
	Sciezka          *string
	RozmiarBajtow    *int64
	Utworzono        string
}

type StreszczenieZrodlaBadania struct {
	Abstrakt    *string
	Tezy        []string
	Metodologia *string
	Wnioski     []string
}

type TabelaZrodlaBadania struct {
	Kod       string
	ZrodloKod string
	Podpis    *string
	Naglowki  []string
	Wiersze   string
	Strona    *int64
}

type RepozytoriumBadanDobudowa interface {
	UstawKatalogZrodla(ctx context.Context, kodZrodla string, katalog KatalogZrodlaBadania) error
	KatalogZrodlaBadania(ctx context.Context, kodZrodla string) (KatalogZrodlaBadania, error)
	UstawLektureZrodla(ctx context.Context, kodZrodla string, lektura LekturaZrodlaBadania) error
	LekturaZrodlaBadania(ctx context.Context, kodZrodla string) (LekturaZrodlaBadania, error)
	UstawTrescZrodla(ctx context.Context, kodZrodla string, tresc TrescZrodlaBadania) error
	TrescZrodlaBadania(ctx context.Context, kodZrodla string) (TrescZrodlaBadania, error)
	UstawWycofanieZrodla(ctx context.Context, kodZrodla string, flaga WycofanieZrodlaBadania) error
	WycofanieZrodlaBadania(ctx context.Context, kodZrodla string) (WycofanieZrodlaBadania, error)
	UsunZrodlo(ctx context.Context, kodZrodla string) (int, error)
	PrzeniesUstaleniaZrodel(ctx context.Context, kodyZrodel []string, kodDocelowy string) (int, error)
	ZapiszZalacznik(ctx context.Context, zalacznik ZalacznikZrodlaBadania) (ZalacznikZrodlaBadania, error)
	Zalaczniki(ctx context.Context, kodZrodla string) ([]ZalacznikZrodlaBadania, error)
	ZapiszStreszczenieZrodla(ctx context.Context, kodZrodla string, s StreszczenieZrodlaBadania) error
	StreszczenieZrodlaBadania(ctx context.Context, kodZrodla string) (StreszczenieZrodlaBadania, error)
	ZapiszTabeleZrodla(ctx context.Context, tabele []TabelaZrodlaBadania) error
	TabeleZrodla(ctx context.Context, kodZrodla string) ([]TabelaZrodlaBadania, error)

	ZapiszAdnotacje(ctx context.Context, a AdnotacjaBadania) (AdnotacjaBadania, error)
	Adnotacje(ctx context.Context, kodZrodla, okno, rodzaj string, limit int) ([]AdnotacjaBadania, error)
	Adnotacja(ctx context.Context, kod string) (AdnotacjaBadania, error)
	UsunAdnotacje(ctx context.Context, kod string) error
	UstawSzczegolyUstalenia(ctx context.Context, kod string, s SzczegolyUstaleniaBadania) error
	SzczegolyUstaleniaBadania(ctx context.Context, kod string) (SzczegolyUstaleniaBadania, error)
	UsunUstalenie(ctx context.Context, kod string) (int, error)
	PrzeniesZrodlaUstalen(ctx context.Context, kodyUstalen []string, kodDocelowy string) (int, error)

	ZapiszKsiazkeKodow(ctx context.Context, okno string, kody []KodBadania) ([]KodBadania, error)
	KsiazkaKodow(ctx context.Context, okno string) ([]KodBadania, error)
	UstawKodyUstalenia(ctx context.Context, kodUstalenia string, kodyKodow []string) error
	KodyUstalenia(ctx context.Context, kodUstalenia string) ([]KodBadania, error)
	MacierzKodow(ctx context.Context, okno string) ([]KomorkaMacierzyBadania, error)

	ZapiszSprzecznosc(ctx context.Context, s SprzecznoscBadania) (SprzecznoscBadania, error)
	Sprzecznosci(ctx context.Context, okno string) ([]SprzecznoscBadania, error)
	Sprzecznosc(ctx context.Context, kod string) (SprzecznoscBadania, error)
	ZapiszWeryfikacje(ctx context.Context, kodUstalenia, werdykt, uzasadnienie string) (string, error)
	ZapiszWatki(ctx context.Context, okno string, watki []WatekBadania) ([]WatekBadania, error)
	ZapiszProwenancje(ctx context.Context, w WpisProwenancjiBadania) error
	Prowenancja(ctx context.Context, kodUstalenia string) ([]WpisProwenancjiBadania, error)

	UstawPytaniaBadania(ctx context.Context, pytania []PytanieBadania) ([]PytanieBadania, error)
	PytaniaBadania(ctx context.Context) ([]PytanieBadania, error)
	UstawSzczegolyPrzestrzeni(ctx context.Context, s SzczegolyPrzestrzeniBadania) error
	SzczegolyPrzestrzeniBadania(ctx context.Context) (SzczegolyPrzestrzeniBadania, error)
	UstawNotatkePrzestrzeni(ctx context.Context, tresc string) (string, error)

	ZapiszWynikiOdkrycia(ctx context.Context, wyniki []WynikOdkryciaBadania) error
	WynikiOdkrycia(ctx context.Context, okno string) ([]WynikOdkryciaBadania, error)
	OdrzucWynikiOdkrycia(ctx context.Context, okno string, klucze []string, powod string) (int, error)
	ZapiszMonitor(ctx context.Context, m MonitorBadania) (MonitorBadania, error)
	Monitory(ctx context.Context, okno string, tylkoWlaczone bool) ([]MonitorBadania, error)
	Monitor(ctx context.Context, kod string) (MonitorBadania, error)
	ZapiszMonitorZeStanem(ctx context.Context, m MonitorBadania) (MonitorBadania, error)

	ZapiszStylCytowania(ctx context.Context, kod, nazwa string) error
	StyleCytowania(ctx context.Context) ([]StylCytowaniaBadania, error)

	Raporty(ctx context.Context, okno string) ([]RaportBadania, error)
	ZapiszSzablonRaportu(ctx context.Context, s SzablonRaportuBadania) error
	SzablonyRaportu(ctx context.Context) ([]SzablonRaportuBadania, error)
	ZapiszWersjeRaportu(ctx context.Context, w WersjaRaportuBadania) (WersjaRaportuBadania, error)
	WersjeRaportu(ctx context.Context, kodRaportu string, limit int) ([]WersjaRaportuBadania, error)
	WersjaRaportuBadania(ctx context.Context, kod string) (WersjaRaportuBadania, error)
	ZapiszKomentarzRaportu(ctx context.Context, k KomentarzRaportuBadania) (KomentarzRaportuBadania, error)
	KomentarzeRaportu(ctx context.Context, kodRaportu string, tylkoOtwarte bool) ([]KomentarzRaportuBadania, error)
	ZapiszBlokRaportu(ctx context.Context, b BlokRaportuBadania) (BlokRaportuBadania, error)
	UstawPrzypisyRaportu(ctx context.Context, kodRaportu, umiejscowienie string, skrocone bool) error

	Eksporty(ctx context.Context, kodRaportu, okno string, limit int) ([]EksportRaportu, error)
	ZapiszSzablonEksportu(ctx context.Context, s SzablonEksportuBadania) (SzablonEksportuBadania, error)
	SzablonyEksportu(ctx context.Context) ([]SzablonEksportuBadania, error)
	ZapiszUdostepnienie(ctx context.Context, kodRaportu, kod, adres string, wygasaO *string) error
}

func (r *repozytoriumBadan) wykonajBadania(ctx context.Context, sqlTekst string, argumenty ...any) error {
	polecenie, err := r.zapytania.przygotuj(ctx, sqlTekst)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, argumenty...); err != nil {
		return fmt.Errorf("dane: badania — polecenie nie powiodło się: %w", err)
	}
	return nil
}

func (r *repozytoriumBadan) pytajBadania(ctx context.Context, sqlTekst string, argumenty ...any) (*sql.Rows, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, sqlTekst)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, argumenty...)
	if err != nil {
		return nil, fmt.Errorf("dane: badania — odczyt nie powiódł się: %w", err)
	}
	return wiersze, nil
}

func (r *repozytoriumBadan) idZrodla(ctx context.Context, kod string) (int64, error) {
	polecenie, err := r.zapytania.przygotuj(ctx,
		`SELECT id FROM zrodlo_badania WHERE identyfikator_zewnetrzny = ? AND `+WarunekKonta)
	if err != nil {
		return 0, err
	}
	var id int64
	err = polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrBrakWiersza
	}
	if err != nil {
		return 0, fmt.Errorf("dane: nie można odnaleźć źródła badania %q: %w", kod, err)
	}
	return id, nil
}

func (r *repozytoriumBadan) idUstalenia(ctx context.Context, kod string) (int64, error) {
	polecenie, err := r.zapytania.przygotuj(ctx,
		`SELECT id FROM ustalenie_badania WHERE identyfikator_zewnetrzny = ? AND `+WarunekKonta)
	if err != nil {
		return 0, err
	}
	var id int64
	err = polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrBrakWiersza
	}
	if err != nil {
		return 0, fmt.Errorf("dane: nie można odnaleźć ustalenia badania %q: %w", kod, err)
	}
	return id, nil
}

func (r *repozytoriumBadan) idRaportu(ctx context.Context, kod string) (int64, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, `SELECT id FROM raport_badania WHERE identyfikator_zewnetrzny = ? AND `+WarunekKonta)
	if err != nil {
		return 0, err
	}
	var id int64
	err = polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrBrakWiersza
	}
	if err != nil {
		return 0, fmt.Errorf("dane: nie można odnaleźć raportu badania %q: %w", kod, err)
	}
	return id, nil
}

// SQLite nie ma typu logicznego; kolumna niesie 0/1.
func wartoscCalkowitaBadania(wartosc bool) int64 {
	if wartosc {
		return 1
	}
	return 0
}

// Lista bez tożsamości własnej idzie kolumną tekstową JSON, nie tabelą podrzędną (migracja 150).
func listaJakoBadania(wartosci []string) string {
	if wartosci == nil {
		wartosci = []string{}
	}
	bajty, err := json.Marshal(wartosci)
	if err != nil {
		return "[]"
	}
	return string(bajty)
}

// Treść nieczytelna wraca listą pustą: kolumna zapisana innym kształtem nie wywraca odczytu wiersza.
func listaZBadania(tekst string) []string {
	lista := []string{}
	if strings.TrimSpace(tekst) == "" {
		return lista
	}
	if err := json.Unmarshal([]byte(tekst), &lista); err != nil {
		return []string{}
	}
	return lista
}

func napisyZWierszyBadania(wiersze *sql.Rows) ([]string, error) {
	defer wiersze.Close()
	lista := []string{}
	for wiersze.Next() {
		var wartosc string
		if err := wiersze.Scan(&wartosc); err != nil {
			return nil, fmt.Errorf("dane: nieczytelna kolumna tekstowa: %w", err)
		}
		lista = append(lista, wartosc)
	}
	return lista, wiersze.Err()
}

// Kontrakt wymienia zbiór w całości: lista pusta znaczy „bez pozycji”, nie „bez zmiany”.
func (r *repozytoriumBadan) UstawKatalogZrodla(ctx context.Context, kodZrodla string,
	katalog KatalogZrodlaBadania) error {

	id, err := r.idZrodla(ctx, kodZrodla)
	if err != nil {
		return err
	}
	return wTransakcji(ctx, r.db, func(t *sql.Tx) error {
		pary := []struct {
			usun, wstaw string
			wartosci    []string
		}{
			{`DELETE FROM etykieta_zrodla_badania WHERE zrodlo_id = ?`,
				`INSERT OR IGNORE INTO etykieta_zrodla_badania (zrodlo_id, etykieta) VALUES (?, ?)`,
				katalog.Etykiety},
			{`DELETE FROM kolekcja_zrodla_badania WHERE zrodlo_id = ?`,
				`INSERT OR IGNORE INTO kolekcja_zrodla_badania (zrodlo_id, kolekcja) VALUES (?, ?)`,
				katalog.Kolekcje},
			{`DELETE FROM pytanie_zrodla_badania WHERE zrodlo_id = ?`,
				`INSERT OR IGNORE INTO pytanie_zrodla_badania (zrodlo_id, pytanie_kod) VALUES (?, ?)`,
				katalog.Pytania},
		}
		for _, para := range pary {
			if para.wartosci == nil {
				continue
			}
			czyszczenie, err := r.zapytania.wTransakcji(ctx, t, para.usun)
			if err != nil {
				return err
			}
			if _, err := czyszczenie.ExecContext(ctx, id); err != nil {
				return fmt.Errorf("dane: nie można wyczyścić katalogu źródła %q: %w", kodZrodla, err)
			}
			wstawienie, err := r.zapytania.wTransakcji(ctx, t, para.wstaw)
			if err != nil {
				return err
			}
			for _, wartosc := range para.wartosci {
				if strings.TrimSpace(wartosc) == "" {
					continue
				}
				if _, err := wstawienie.ExecContext(ctx, id, wartosc); err != nil {
					return fmt.Errorf("dane: nie można zapisać katalogu źródła %q: %w", kodZrodla, err)
				}
			}
		}
		return nil
	})
}

func (r *repozytoriumBadan) KatalogZrodlaBadania(ctx context.Context, kodZrodla string) (KatalogZrodlaBadania, error) {
	id, err := r.idZrodla(ctx, kodZrodla)
	if err != nil {
		return KatalogZrodlaBadania{}, err
	}
	katalog := KatalogZrodlaBadania{}
	for _, para := range []struct {
		sqlTekst string
		cel      *[]string
	}{
		{`SELECT etykieta FROM etykieta_zrodla_badania WHERE zrodlo_id = ? ORDER BY etykieta`, &katalog.Etykiety},
		{`SELECT kolekcja FROM kolekcja_zrodla_badania WHERE zrodlo_id = ? ORDER BY kolekcja`, &katalog.Kolekcje},
		{`SELECT pytanie_kod FROM pytanie_zrodla_badania WHERE zrodlo_id = ? ORDER BY pytanie_kod`, &katalog.Pytania},
	} {
		wiersze, err := r.pytajBadania(ctx, para.sqlTekst, id)
		if err != nil {
			return KatalogZrodlaBadania{}, err
		}
		lista, err := napisyZWierszyBadania(wiersze)
		if err != nil {
			return KatalogZrodlaBadania{}, err
		}
		*para.cel = lista
	}
	return katalog, nil
}

func (r *repozytoriumBadan) UstawLektureZrodla(ctx context.Context, kodZrodla string,
	lektura LekturaZrodlaBadania) error {

	if _, err := r.idZrodla(ctx, kodZrodla); err != nil {
		return err
	}
	return r.zapiszWGranicyBadania(ctx, `UPDATE zrodlo_badania SET
	        stan_lektury  = COALESCE(NULLIF(?, ''), stan_lektury),
	        etap_indeks   = COALESCE(?, etap_indeks),
	        identyfikator = COALESCE(?, identyfikator),
	        csl_json      = COALESCE(?, csl_json)
	    WHERE identyfikator_zewnetrzny = ? AND `+WarunekKonta, "źródło badania", kodZrodla,
		lektura.StanLektury, liczbaDoKolumny(lektura.EtapIndeks),
		tekstDoKolumny(lektura.Identyfikat), tekstDoKolumny(lektura.CslJson), kodZrodla,
		KontoOperatora(ctx))
}

func (r *repozytoriumBadan) LekturaZrodlaBadania(ctx context.Context, kodZrodla string) (LekturaZrodlaBadania, error) {
	polecenie, err := r.zapytania.przygotuj(ctx,
		`SELECT stan_lektury, etap_indeks, identyfikator, csl_json
		 FROM zrodlo_badania WHERE identyfikator_zewnetrzny = ? AND `+WarunekKonta)
	if err != nil {
		return LekturaZrodlaBadania{}, err
	}
	var lektura LekturaZrodlaBadania
	var etap sql.NullInt64
	var identyfikator, csl sql.NullString
	err = polecenie.QueryRowContext(ctx, kodZrodla, KontoOperatora(ctx)).Scan(&lektura.StanLektury, &etap, &identyfikator, &csl)
	if errors.Is(err, sql.ErrNoRows) {
		return LekturaZrodlaBadania{}, ErrBrakWiersza
	}
	if err != nil {
		return LekturaZrodlaBadania{}, fmt.Errorf("dane: nieczytelna lektura źródła %q: %w", kodZrodla, err)
	}
	lektura.EtapIndeks = liczbaZKolumny(etap)
	lektura.Identyfikat = tekstZKolumny(identyfikator)
	lektura.CslJson = tekstZKolumny(csl)
	return lektura, nil
}

func (r *repozytoriumBadan) UstawTrescZrodla(ctx context.Context, kodZrodla string,
	tresc TrescZrodlaBadania) error {

	if _, err := r.idZrodla(ctx, kodZrodla); err != nil {
		return err
	}
	return r.zapiszWGranicyBadania(ctx, `UPDATE zrodlo_badania
	    SET tresc = ?, stron = ?, warstwa_tekstu = ?
	    WHERE identyfikator_zewnetrzny = ? AND `+WarunekKonta, "źródło badania", kodZrodla,
		tekstDoKolumny(tresc.Tekst), liczbaDoKolumny(tresc.Stron),
		wartoscCalkowitaBadania(tresc.WarstwaTekstu), kodZrodla, KontoOperatora(ctx))
}

func (r *repozytoriumBadan) TrescZrodlaBadania(ctx context.Context, kodZrodla string) (TrescZrodlaBadania, error) {
	polecenie, err := r.zapytania.przygotuj(ctx,
		`SELECT tresc, stron, warstwa_tekstu FROM zrodlo_badania
		 WHERE identyfikator_zewnetrzny = ? AND `+WarunekKonta)
	if err != nil {
		return TrescZrodlaBadania{}, err
	}
	var tresc TrescZrodlaBadania
	var tekst sql.NullString
	var stron sql.NullInt64
	var warstwa int64
	err = polecenie.QueryRowContext(ctx, kodZrodla, KontoOperatora(ctx)).Scan(&tekst, &stron, &warstwa)
	if errors.Is(err, sql.ErrNoRows) {
		return TrescZrodlaBadania{}, ErrBrakWiersza
	}
	if err != nil {
		return TrescZrodlaBadania{}, fmt.Errorf("dane: nieczytelna treść źródła %q: %w", kodZrodla, err)
	}
	tresc.Tekst = tekstZKolumny(tekst)
	tresc.Stron = liczbaZKolumny(stron)
	tresc.WarstwaTekstu = warstwa != 0
	return tresc, nil
}

func (r *repozytoriumBadan) UstawWycofanieZrodla(ctx context.Context, kodZrodla string,
	flaga WycofanieZrodlaBadania) error {

	if _, err := r.idZrodla(ctx, kodZrodla); err != nil {
		return err
	}
	return r.zapiszWGranicyBadania(ctx, `UPDATE zrodlo_badania
	    SET wycofanie_stan = ?, wycofanie_adres = ?,
	        wycofanie_sprawdzono = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	    WHERE identyfikator_zewnetrzny = ? AND `+WarunekKonta, "źródło badania", kodZrodla,
		flaga.Stan, tekstDoKolumny(flaga.Adres), kodZrodla, KontoOperatora(ctx))
}

func (r *repozytoriumBadan) WycofanieZrodlaBadania(ctx context.Context, kodZrodla string) (WycofanieZrodlaBadania, error) {
	polecenie, err := r.zapytania.przygotuj(ctx,
		`SELECT wycofanie_stan, wycofanie_adres, wycofanie_sprawdzono
		 FROM zrodlo_badania WHERE identyfikator_zewnetrzny = ? AND `+WarunekKonta)
	if err != nil {
		return WycofanieZrodlaBadania{}, err
	}
	var stan, adres, sprawdzono sql.NullString
	err = polecenie.QueryRowContext(ctx, kodZrodla, KontoOperatora(ctx)).Scan(&stan, &adres, &sprawdzono)
	if errors.Is(err, sql.ErrNoRows) {
		return WycofanieZrodlaBadania{}, ErrBrakWiersza
	}
	if err != nil {
		return WycofanieZrodlaBadania{}, fmt.Errorf("dane: nieczytelna flaga wycofania %q: %w", kodZrodla, err)
	}
	return WycofanieZrodlaBadania{Stan: stan.String, Adres: tekstZKolumny(adres), Sprawdzono: sprawdzono.String}, nil
}

func (r *repozytoriumBadan) UsunZrodlo(ctx context.Context, kodZrodla string) (int, error) {
	id, err := r.idZrodla(ctx, kodZrodla)
	if err != nil {
		return 0, err
	}
	var odwiazane int
	err = wTransakcji(ctx, r.db, func(t *sql.Tx) error {
		// Wiązanie wisi na ustaleniu, a konto niesie ustalenie, więc liczba
		// bez sięgnięcia do korzenia mówiłaby o ustaleniach konta obcego.
		liczenie, err := r.zapytania.wTransakcji(ctx, t,
			`SELECT COUNT(*) FROM zrodlo_ustalenia_badania zu
			  WHERE zu.zrodlo_id = ?
			    AND EXISTS (SELECT 1 FROM ustalenie_badania w
			                 WHERE w.id = zu.ustalenie_id AND `+WarunekKonta+`)`)
		if err != nil {
			return err
		}
		if err := liczenie.QueryRowContext(ctx, id, KontoOperatora(ctx)).Scan(&odwiazane); err != nil {
			return fmt.Errorf("dane: nie można policzyć ustaleń źródła %q: %w", kodZrodla, err)
		}
		usuniecie, err := r.zapytania.wTransakcji(ctx, t,
			`DELETE FROM zrodlo_badania WHERE id = ? AND `+WarunekKonta)
		if err != nil {
			return err
		}
		wynik, err := usuniecie.ExecContext(ctx, id, KontoOperatora(ctx))
		if err != nil {
			return fmt.Errorf("dane: nie można usunąć źródła %q: %w", kodZrodla, err)
		}
		return sprawdzTrafienieZapisu(wynik, "źródło badania", kodZrodla)
	})
	return odwiazane, err
}

func (r *repozytoriumBadan) PrzeniesUstaleniaZrodel(ctx context.Context, kodyZrodel []string,
	kodDocelowy string) (int, error) {

	docelowe, err := r.idZrodla(ctx, kodDocelowy)
	if err != nil {
		return 0, err
	}
	przeniesione := 0
	err = wTransakcji(ctx, r.db, func(t *sql.Tx) error {
		for _, kod := range kodyZrodel {
			if kod == kodDocelowy {
				continue
			}
			wskazanie, err := r.zapytania.wTransakcji(ctx, t,
				`SELECT id FROM zrodlo_badania WHERE identyfikator_zewnetrzny = ? AND `+WarunekKonta)
			if err != nil {
				return err
			}
			var scalane int64
			err = wskazanie.QueryRowContext(ctx, kod, KontoOperatora(ctx)).Scan(&scalane)
			if errors.Is(err, sql.ErrNoRows) {
				continue
			}
			if err != nil {
				return fmt.Errorf("dane: nie można odnaleźć źródła scalanego %q: %w", kod, err)
			}
			// Przepięcie zmienia wiersze dziecka, więc sięga korzenia: bez tego
			// scalenie źródeł przestawiałoby wiązania ustaleń konta obcego.
			przepiecie, err := r.zapytania.wTransakcji(ctx, t,
				`INSERT OR IGNORE INTO zrodlo_ustalenia_badania (ustalenie_id, zrodlo_id)
				 SELECT zu.ustalenie_id, ? FROM zrodlo_ustalenia_badania zu
				  WHERE zu.zrodlo_id = ?
				    AND EXISTS (SELECT 1 FROM ustalenie_badania w
				                 WHERE w.id = zu.ustalenie_id AND `+WarunekKonta+`)`)
			if err != nil {
				return err
			}
			wynik, err := przepiecie.ExecContext(ctx, docelowe, scalane, KontoOperatora(ctx))
			if err != nil {
				return fmt.Errorf("dane: nie można przepiąć ustaleń źródła %q: %w", kod, err)
			}
			ile, _ := wynik.RowsAffected()
			przeniesione += int(ile)

			usuniecie, err := r.zapytania.wTransakcji(ctx, t,
				`DELETE FROM zrodlo_badania WHERE id = ? AND `+WarunekKonta)
			if err != nil {
				return err
			}
			usuniete, err := usuniecie.ExecContext(ctx, scalane, KontoOperatora(ctx))
			if err != nil {
				return fmt.Errorf("dane: nie można usunąć źródła scalonego %q: %w", kod, err)
			}
			if err := sprawdzTrafienieZapisu(usuniete, "źródło badania", kod); err != nil {
				return err
			}
		}
		return nil
	})
	return przeniesione, err
}

func (r *repozytoriumBadan) ZapiszZalacznik(ctx context.Context,
	zalacznik ZalacznikZrodlaBadania) (ZalacznikZrodlaBadania, error) {

	id, err := r.idZrodla(ctx, zalacznik.ZrodloKod)
	if err != nil {
		return ZalacznikZrodlaBadania{}, err
	}
	err = r.wykonajBadania(ctx, `INSERT INTO zalacznik_zrodla_badania
	    (identyfikator_zewnetrzny, zrodlo_id, rodzaj, plik_biblioteki_id, sciezka, rozmiar_bajtow)
	    VALUES (?, ?, ?, ?, ?, ?)
	    ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	        rodzaj = excluded.rodzaj,
	        plik_biblioteki_id = excluded.plik_biblioteki_id,
	        sciezka = excluded.sciezka,
	        rozmiar_bajtow = excluded.rozmiar_bajtow`,
		zalacznik.Kod, id, zalacznik.Rodzaj, tekstDoKolumny(zalacznik.PlikBibliotekiID),
		tekstDoKolumny(zalacznik.Sciezka), liczbaDoKolumny(zalacznik.RozmiarBajtow))
	if err != nil {
		return ZalacznikZrodlaBadania{}, err
	}
	lista, err := r.Zalaczniki(ctx, zalacznik.ZrodloKod)
	if err != nil {
		return ZalacznikZrodlaBadania{}, err
	}
	for _, wiersz := range lista {
		if wiersz.Kod == zalacznik.Kod {
			return wiersz, nil
		}
	}
	return ZalacznikZrodlaBadania{}, ErrBrakWiersza
}

func (r *repozytoriumBadan) Zalaczniki(ctx context.Context, kodZrodla string) ([]ZalacznikZrodlaBadania, error) {
	id, err := r.idZrodla(ctx, kodZrodla)
	if err != nil {
		return nil, err
	}
	wiersze, err := r.pytajBadania(ctx, `SELECT identyfikator_zewnetrzny, rodzaj, plik_biblioteki_id,
	        sciezka, rozmiar_bajtow, utworzono
	    FROM zalacznik_zrodla_badania WHERE zrodlo_id = ? ORDER BY utworzono DESC, id DESC`, id)
	if err != nil {
		return nil, err
	}
	defer wiersze.Close()

	lista := []ZalacznikZrodlaBadania{}
	for wiersze.Next() {
		wiersz := ZalacznikZrodlaBadania{ZrodloKod: kodZrodla}
		var plik, sciezka sql.NullString
		var rozmiar sql.NullInt64
		if err := wiersze.Scan(&wiersz.Kod, &wiersz.Rodzaj, &plik, &sciezka, &rozmiar, &wiersz.Utworzono); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz załącznika źródła: %w", err)
		}
		wiersz.PlikBibliotekiID = tekstZKolumny(plik)
		wiersz.Sciezka = tekstZKolumny(sciezka)
		wiersz.RozmiarBajtow = liczbaZKolumny(rozmiar)
		lista = append(lista, wiersz)
	}
	return lista, wiersze.Err()
}

func (r *repozytoriumBadan) ZapiszStreszczenieZrodla(ctx context.Context, kodZrodla string,
	s StreszczenieZrodlaBadania) error {

	id, err := r.idZrodla(ctx, kodZrodla)
	if err != nil {
		return err
	}
	return r.wykonajBadania(ctx, `INSERT INTO streszczenie_zrodla_badania
	    (zrodlo_id, abstrakt, tezy, metodologia, wnioski)
	    VALUES (?, ?, ?, ?, ?)
	    ON CONFLICT(zrodlo_id) DO UPDATE SET
	        abstrakt = excluded.abstrakt, tezy = excluded.tezy,
	        metodologia = excluded.metodologia, wnioski = excluded.wnioski`,
		id, tekstDoKolumny(s.Abstrakt), listaJakoBadania(s.Tezy),
		tekstDoKolumny(s.Metodologia), listaJakoBadania(s.Wnioski))
}

func (r *repozytoriumBadan) StreszczenieZrodlaBadania(ctx context.Context, kodZrodla string) (StreszczenieZrodlaBadania, error) {
	id, err := r.idZrodla(ctx, kodZrodla)
	if err != nil {
		return StreszczenieZrodlaBadania{}, err
	}
	polecenie, err := r.zapytania.przygotuj(ctx,
		`SELECT abstrakt, tezy, metodologia, wnioski FROM streszczenie_zrodla_badania WHERE zrodlo_id = ?`)
	if err != nil {
		return StreszczenieZrodlaBadania{}, err
	}
	var abstrakt, metodologia sql.NullString
	var tezy, wnioski string
	err = polecenie.QueryRowContext(ctx, id).Scan(&abstrakt, &tezy, &metodologia, &wnioski)
	if errors.Is(err, sql.ErrNoRows) {
		return StreszczenieZrodlaBadania{}, ErrBrakWiersza
	}
	if err != nil {
		return StreszczenieZrodlaBadania{}, fmt.Errorf("dane: nieczytelne streszczenie źródła %q: %w", kodZrodla, err)
	}
	return StreszczenieZrodlaBadania{
		Abstrakt: tekstZKolumny(abstrakt), Tezy: listaZBadania(tezy),
		Metodologia: tekstZKolumny(metodologia), Wnioski: listaZBadania(wnioski),
	}, nil
}

func (r *repozytoriumBadan) ZapiszTabeleZrodla(ctx context.Context, tabele []TabelaZrodlaBadania) error {
	for _, tabela := range tabele {
		id, err := r.idZrodla(ctx, tabela.ZrodloKod)
		if err != nil {
			return err
		}
		err = r.wykonajBadania(ctx, `INSERT INTO tabela_zrodla_badania
		    (identyfikator_zewnetrzny, zrodlo_id, podpis, naglowki, wiersze, kotwica_strona)
		    VALUES (?, ?, ?, ?, ?, ?)
		    ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
		        podpis = excluded.podpis, naglowki = excluded.naglowki,
		        wiersze = excluded.wiersze, kotwica_strona = excluded.kotwica_strona`,
			tabela.Kod, id, tekstDoKolumny(tabela.Podpis), listaJakoBadania(tabela.Naglowki),
			tabela.Wiersze, liczbaDoKolumny(tabela.Strona))
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *repozytoriumBadan) TabeleZrodla(ctx context.Context, kodZrodla string) ([]TabelaZrodlaBadania, error) {
	id, err := r.idZrodla(ctx, kodZrodla)
	if err != nil {
		return nil, err
	}
	wiersze, err := r.pytajBadania(ctx, `SELECT identyfikator_zewnetrzny, podpis, naglowki, wiersze, kotwica_strona
	    FROM tabela_zrodla_badania WHERE zrodlo_id = ? ORDER BY id`, id)
	if err != nil {
		return nil, err
	}
	defer wiersze.Close()

	lista := []TabelaZrodlaBadania{}
	for wiersze.Next() {
		tabela := TabelaZrodlaBadania{ZrodloKod: kodZrodla}
		var podpis sql.NullString
		var naglowki string
		var strona sql.NullInt64
		if err := wiersze.Scan(&tabela.Kod, &podpis, &naglowki, &tabela.Wiersze, &strona); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz tabeli źródła: %w", err)
		}
		tabela.Podpis = tekstZKolumny(podpis)
		tabela.Naglowki = listaZBadania(naglowki)
		tabela.Strona = liczbaZKolumny(strona)
		lista = append(lista, tabela)
	}
	return lista, wiersze.Err()
}
