// Odpowiedzialność pliku: materiał wytworzony w sesji przeglądania i czynności
// prowadzone przez rdzeń — wytwory i zrzuty (migracja 175), pobrania (176),
// makra (177) oraz granice działania Wykonawcy (178).
//
// Wspólne im jest to, że każde z nich zostawia ślad poza bazą albo poza chwilą:
// wytwór i zrzut mają bajty w magazynie, pobranie ma plik na dysku, makro ma
// kroki do odtworzenia, granica obowiązuje kolejne przebiegi. Wiersz jest tu
// wskazaniem na coś, co istnieje naprawdę — dlatego kolumna odwołania jest
// NOT NULL: wytwór bez bajtów byłby meldunkiem bez skutku.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// WytworPrzegladania to wiersz tabeli `wytwor_przegladania`.
type WytworPrzegladania struct {
	ID                  int64
	Kod                 string
	Okno                string
	Rodzaj              string
	Tytul               *string
	TrescOdwolanie      string
	TypMime             *string
	RozmiarBajtow       *int64
	UrlZrodla           *string
	MigawkaZewnetrznaID *string
	Utworzono           string
}

// ZrzutPrzegladania to wiersz tabeli `zrzut_przegladania`.
type ZrzutPrzegladania struct {
	ID                  int64
	Kod                 string
	Okno                string
	MigawkaZewnetrznaID *string
	TrescOdwolanie      string
	Tryb                string
	Format              string
	Szerokosc           int64
	Wysokosc            int64
	RozmiarBajtow       *int64
	Utworzono           string
}

// PobraniePrzegladania to wiersz tabeli `pobranie_przegladania`.
type PobraniePrzegladania struct {
	ID              int64
	Kod             string
	Okno            string
	Url             string
	NazwaPliku      *string
	SciezkaDocelowa *string
	TypMime         *string
	Stan            string
	OdebranoBajtow  *int64
	RazemBajtow     *int64
	KomunikatBledu  *string
	Rozpoczeto      string
	Zakonczono      *string
}

// MakroPrzegladania to wiersz tabeli `makro_przegladania`.
type MakroPrzegladania struct {
	ID             int64
	Kod            string
	Okno           string
	Nazwa          string
	KrokiJson      *string
	Nagrywanie     bool
	AutomatykaKod  *string
	Utworzono      string
	Zaktualizowano string
}

// GranicaWykonawcy to wiersz tabeli `granica_wykonawcy_przegladania`.
type GranicaWykonawcy struct {
	ID                    int64
	Zasieg                string
	ZasiegID              string
	MaxKrokow             int64
	MaxCzasSekund         int64
	DomenyDozwoloneJson   *string
	DomenyZablokowaneJson *string
	PotwierdzajWyslanie   bool
	Zaktualizowano        string
}

const (
	kolumnyWytworu = `id, identyfikator_zewnetrzny, okno, rodzaj, tytul, tresc_odwolanie,
	                  typ_mime, rozmiar_bajtow, url_zrodla, migawka_zewnetrzna_id, utworzono`

	zapiszWytwor = `INSERT INTO wytwor_przegladania
	                (identyfikator_zewnetrzny, okno, rodzaj, tytul, tresc_odwolanie, typ_mime,
	                 rozmiar_bajtow, url_zrodla, migawka_zewnetrzna_id)
	                VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	pobierzWytwor = `SELECT ` + kolumnyWytworu + ` FROM wytwor_przegladania
	                 WHERE identyfikator_zewnetrzny = ?`

	listaWytworow = `SELECT ` + kolumnyWytworu + ` FROM wytwor_przegladania
	                 WHERE okno = ? AND (? = '' OR rodzaj = ?)
	                 ORDER BY utworzono DESC, id DESC LIMIT ?`

	kolumnyZrzutu = `id, identyfikator_zewnetrzny, okno, migawka_zewnetrzna_id, tresc_odwolanie,
	                 tryb, format, szerokosc, wysokosc, rozmiar_bajtow, utworzono`

	zapiszZrzut = `INSERT INTO zrzut_przegladania
	               (identyfikator_zewnetrzny, okno, migawka_zewnetrzna_id, tresc_odwolanie,
	                tryb, format, szerokosc, wysokosc, rozmiar_bajtow)
	               VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	pobierzZrzut = `SELECT ` + kolumnyZrzutu + ` FROM zrzut_przegladania
	                WHERE identyfikator_zewnetrzny = ?`

	pobierzZrzutPoOdwolaniu = `SELECT ` + kolumnyZrzutu + ` FROM zrzut_przegladania
	                           WHERE tresc_odwolanie = ? ORDER BY id DESC LIMIT 1`

	pobierzZrzutMigawki = `SELECT ` + kolumnyZrzutu + ` FROM zrzut_przegladania
	                       WHERE migawka_zewnetrzna_id = ? ORDER BY id DESC LIMIT 1`

	kolumnyPobrania = `id, identyfikator_zewnetrzny, okno, url, nazwa_pliku, sciezka_docelowa,
	                   typ_mime, stan, odebrano_bajtow, razem_bajtow, komunikat_bledu,
	                   rozpoczeto, zakonczono`

	zapiszPobranie = `INSERT INTO pobranie_przegladania
	                  (identyfikator_zewnetrzny, okno, url, nazwa_pliku, sciezka_docelowa, typ_mime,
	                   stan, odebrano_bajtow, razem_bajtow, komunikat_bledu, zakonczono)
	                  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	                  ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                      url = excluded.url,
	                      nazwa_pliku = excluded.nazwa_pliku,
	                      sciezka_docelowa = excluded.sciezka_docelowa,
	                      typ_mime = excluded.typ_mime,
	                      stan = excluded.stan,
	                      odebrano_bajtow = excluded.odebrano_bajtow,
	                      razem_bajtow = excluded.razem_bajtow,
	                      komunikat_bledu = excluded.komunikat_bledu,
	                      zakonczono = excluded.zakonczono`

	pobierzPobranie = `SELECT ` + kolumnyPobrania + ` FROM pobranie_przegladania
	                   WHERE identyfikator_zewnetrzny = ?`

	listaPobran = `SELECT ` + kolumnyPobrania + ` FROM pobranie_przegladania
	               WHERE (? = '' OR okno = ?) AND (? = '' OR stan = ?)
	               ORDER BY rozpoczeto DESC, id DESC LIMIT ?`

	usunPobranie = `DELETE FROM pobranie_przegladania WHERE identyfikator_zewnetrzny = ?`

	kolumnyMakra = `id, identyfikator_zewnetrzny, okno, nazwa, kroki_json, nagrywanie,
	                automatyka_zewnetrzna_id, utworzono, zaktualizowano`

	zapiszMakro = `INSERT INTO makro_przegladania
	               (identyfikator_zewnetrzny, okno, nazwa, kroki_json, nagrywanie, automatyka_zewnetrzna_id)
	               VALUES (?, ?, ?, ?, ?, ?)
	               ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                   nazwa = excluded.nazwa,
	                   kroki_json = excluded.kroki_json,
	                   nagrywanie = excluded.nagrywanie,
	                   automatyka_zewnetrzna_id = excluded.automatyka_zewnetrzna_id,
	                   zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')`

	pobierzMakro = `SELECT ` + kolumnyMakra + ` FROM makro_przegladania
	                WHERE identyfikator_zewnetrzny = ?`

	listaMakr = `SELECT ` + kolumnyMakra + ` FROM makro_przegladania
	             WHERE okno = ? ORDER BY utworzono DESC, id DESC LIMIT ?`

	kolumnyGranicy = `id, zasieg, zasieg_id, max_krokow, max_czas_sekund,
	                  domeny_dozwolone_json, domeny_zablokowane_json,
	                  potwierdzaj_wyslanie, zaktualizowano`

	zapiszGranice = `INSERT INTO granica_wykonawcy_przegladania
	                 (zasieg, zasieg_id, max_krokow, max_czas_sekund, domeny_dozwolone_json,
	                  domeny_zablokowane_json, potwierdzaj_wyslanie)
	                 VALUES (?, ?, ?, ?, ?, ?, ?)
	                 ON CONFLICT(zasieg, zasieg_id) DO UPDATE SET
	                     max_krokow = excluded.max_krokow,
	                     max_czas_sekund = excluded.max_czas_sekund,
	                     domeny_dozwolone_json = excluded.domeny_dozwolone_json,
	                     domeny_zablokowane_json = excluded.domeny_zablokowane_json,
	                     potwierdzaj_wyslanie = excluded.potwierdzaj_wyslanie,
	                     zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')`

	pobierzGranice = `SELECT ` + kolumnyGranicy + ` FROM granica_wykonawcy_przegladania
	                  WHERE zasieg = ? AND zasieg_id = ?`
)

// ZapiszWytwor odkłada wytwór sesji przeglądania. Wytwór jest wpisem
// historii — powstaje nowym wierszem, nie nadpisaniem poprzedniego.
func (r *repozytoriumPrzegladania) ZapiszWytwor(ctx context.Context,
	wytwor WytworPrzegladania) (WytworPrzegladania, error) {

	if wytwor.Kod == "" || wytwor.Okno == "" || wytwor.TrescOdwolanie == "" {
		return WytworPrzegladania{}, fmt.Errorf("dane: wytwór przeglądania bez identyfikatora, okna albo odwołania do treści")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszWytwor)
	if err != nil {
		return WytworPrzegladania{}, err
	}
	_, err = polecenie.ExecContext(ctx, wytwor.Kod, wytwor.Okno, wytwor.Rodzaj,
		tekstDoKolumny(wytwor.Tytul), wytwor.TrescOdwolanie, tekstDoKolumny(wytwor.TypMime),
		liczbaDoKolumny(wytwor.RozmiarBajtow), tekstDoKolumny(wytwor.UrlZrodla),
		tekstDoKolumny(wytwor.MigawkaZewnetrznaID))
	if err != nil {
		return WytworPrzegladania{}, fmt.Errorf("dane: nie można zapisać wytworu %q: %w", wytwor.Kod, err)
	}
	return r.Wytwor(ctx, wytwor.Kod)
}

// Wytwor oddaje wytwór o wskazanym kodzie.
func (r *repozytoriumPrzegladania) Wytwor(ctx context.Context, kod string) (WytworPrzegladania, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzWytwor)
	if err != nil {
		return WytworPrzegladania{}, err
	}
	wytwor, err := odczytajWytwor(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return WytworPrzegladania{}, ErrBrakWiersza
	}
	if err != nil {
		return WytworPrzegladania{}, fmt.Errorf("dane: nieczytelny wytwór %q: %w", kod, err)
	}
	return wytwor, nil
}

// Wytwory oddaje wytwory okna, opcjonalnie zawężone do jednego rodzaju.
func (r *repozytoriumPrzegladania) Wytwory(ctx context.Context, okno, rodzaj string, limit int) ([]WytworPrzegladania, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaWytworow)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno, rodzaj, rodzaj, granicaWykazu(limit))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać wytworów okna %q: %w", okno, err)
	}
	defer wiersze.Close()

	lista := []WytworPrzegladania{}
	for wiersze.Next() {
		wytwor, err := odczytajWytwor(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz wytworów okna %q: %w", okno, err)
		}
		lista = append(lista, wytwor)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt wytworów okna %q: %w", okno, err)
	}
	return lista, nil
}

// ZapiszZrzut odkłada wiersz zrzutu strony.
func (r *repozytoriumPrzegladania) ZapiszZrzut(ctx context.Context, zrzut ZrzutPrzegladania) (ZrzutPrzegladania, error) {
	if zrzut.Kod == "" || zrzut.Okno == "" || zrzut.TrescOdwolanie == "" {
		return ZrzutPrzegladania{}, fmt.Errorf("dane: zrzut bez identyfikatora, okna albo odwołania do treści")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszZrzut)
	if err != nil {
		return ZrzutPrzegladania{}, err
	}
	_, err = polecenie.ExecContext(ctx, zrzut.Kod, zrzut.Okno,
		tekstDoKolumny(zrzut.MigawkaZewnetrznaID), zrzut.TrescOdwolanie, zrzut.Tryb,
		zrzut.Format, zrzut.Szerokosc, zrzut.Wysokosc, liczbaDoKolumny(zrzut.RozmiarBajtow))
	if err != nil {
		return ZrzutPrzegladania{}, fmt.Errorf("dane: nie można zapisać zrzutu %q: %w", zrzut.Kod, err)
	}
	return r.Zrzut(ctx, zrzut.Kod)
}

// Zrzut oddaje zrzut o wskazanym kodzie.
func (r *repozytoriumPrzegladania) Zrzut(ctx context.Context, kod string) (ZrzutPrzegladania, error) {
	return r.jedenZrzut(ctx, pobierzZrzut, kod)
}

// ZrzutPoOdwolaniu oddaje zrzut leżący pod wskazanym odwołaniem magazynu —
// tą drogą pyta `browser.snapshot.screenshot.get` z polem `screenshotRef`.
func (r *repozytoriumPrzegladania) ZrzutPoOdwolaniu(ctx context.Context, odwolanie string) (ZrzutPrzegladania, error) {
	return r.jedenZrzut(ctx, pobierzZrzutPoOdwolaniu, odwolanie)
}

// ZrzutMigawki oddaje najświeższy zrzut wykonany przy wskazanej migawce.
func (r *repozytoriumPrzegladania) ZrzutMigawki(ctx context.Context, migawka string) (ZrzutPrzegladania, error) {
	return r.jedenZrzut(ctx, pobierzZrzutMigawki, migawka)
}

// jedenZrzut wykonuje odczyt jednego zrzutu wskazanym zapytaniem.
func (r *repozytoriumPrzegladania) jedenZrzut(ctx context.Context, zapytanie, wskazanie string) (ZrzutPrzegladania, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return ZrzutPrzegladania{}, err
	}
	zrzut, err := odczytajZrzut(polecenie.QueryRowContext(ctx, wskazanie))
	if errors.Is(err, sql.ErrNoRows) {
		return ZrzutPrzegladania{}, ErrBrakWiersza
	}
	if err != nil {
		return ZrzutPrzegladania{}, fmt.Errorf("dane: nieczytelny zrzut %q: %w", wskazanie, err)
	}
	return zrzut, nil
}

// ZapiszPobranie zakłada pobranie albo nadpisuje jego stan i postęp.
func (r *repozytoriumPrzegladania) ZapiszPobranie(ctx context.Context,
	pobranie PobraniePrzegladania) (PobraniePrzegladania, error) {

	if pobranie.Kod == "" || pobranie.Okno == "" || pobranie.Url == "" {
		return PobraniePrzegladania{}, fmt.Errorf("dane: pobranie bez identyfikatora, okna albo adresu")
	}
	if pobranie.Stan == "" {
		pobranie.Stan = "queued"
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszPobranie)
	if err != nil {
		return PobraniePrzegladania{}, err
	}
	_, err = polecenie.ExecContext(ctx, pobranie.Kod, pobranie.Okno, pobranie.Url,
		tekstDoKolumny(pobranie.NazwaPliku), tekstDoKolumny(pobranie.SciezkaDocelowa),
		tekstDoKolumny(pobranie.TypMime), pobranie.Stan, liczbaDoKolumny(pobranie.OdebranoBajtow),
		liczbaDoKolumny(pobranie.RazemBajtow), tekstDoKolumny(pobranie.KomunikatBledu),
		tekstDoKolumny(pobranie.Zakonczono))
	if err != nil {
		return PobraniePrzegladania{}, fmt.Errorf("dane: nie można zapisać pobrania %q: %w", pobranie.Kod, err)
	}
	return r.Pobranie(ctx, pobranie.Kod)
}

// Pobranie oddaje pobranie o wskazanym kodzie.
func (r *repozytoriumPrzegladania) Pobranie(ctx context.Context, kod string) (PobraniePrzegladania, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzPobranie)
	if err != nil {
		return PobraniePrzegladania{}, err
	}
	pobranie, err := odczytajPobranie(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return PobraniePrzegladania{}, ErrBrakWiersza
	}
	if err != nil {
		return PobraniePrzegladania{}, fmt.Errorf("dane: nieczytelne pobranie %q: %w", kod, err)
	}
	return pobranie, nil
}

// Pobrania oddaje pobrania, opcjonalnie zawężone do okna i do jednego stanu.
func (r *repozytoriumPrzegladania) Pobrania(ctx context.Context, okno, stan string, limit int) ([]PobraniePrzegladania, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaPobran)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno, okno, stan, stan, granicaWykazu(limit))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać pobrań: %w", err)
	}
	defer wiersze.Close()

	lista := []PobraniePrzegladania{}
	for wiersze.Next() {
		pobranie, err := odczytajPobranie(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz pobrań: %w", err)
		}
		lista = append(lista, pobranie)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt pobrań: %w", err)
	}
	return lista, nil
}

// UsunPobranie zdejmuje pobranie z wykazu.
func (r *repozytoriumPrzegladania) UsunPobranie(ctx context.Context, kod string) (bool, error) {
	return r.usunWiersz(ctx, usunPobranie, kod, "pobranie")
}

// ZapiszMakro zakłada makro albo nadpisuje zastane wraz z krokami.
func (r *repozytoriumPrzegladania) ZapiszMakro(ctx context.Context, makro MakroPrzegladania) (MakroPrzegladania, error) {
	if makro.Kod == "" || makro.Okno == "" || makro.Nazwa == "" {
		return MakroPrzegladania{}, fmt.Errorf("dane: makro przeglądania bez identyfikatora, okna albo nazwy")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszMakro)
	if err != nil {
		return MakroPrzegladania{}, err
	}
	_, err = polecenie.ExecContext(ctx, makro.Kod, makro.Okno, makro.Nazwa,
		tekstDoKolumny(makro.KrokiJson), liczbaLogiczna(makro.Nagrywanie),
		tekstDoKolumny(makro.AutomatykaKod))
	if err != nil {
		return MakroPrzegladania{}, fmt.Errorf("dane: nie można zapisać makra %q: %w", makro.Kod, err)
	}
	return r.Makro(ctx, makro.Kod)
}

// Makro oddaje makro o wskazanym kodzie.
func (r *repozytoriumPrzegladania) Makro(ctx context.Context, kod string) (MakroPrzegladania, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzMakro)
	if err != nil {
		return MakroPrzegladania{}, err
	}
	makro, err := odczytajMakro(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return MakroPrzegladania{}, ErrBrakWiersza
	}
	if err != nil {
		return MakroPrzegladania{}, fmt.Errorf("dane: nieczytelne makro %q: %w", kod, err)
	}
	return makro, nil
}

// Makra oddaje makra okna od najnowszego.
func (r *repozytoriumPrzegladania) Makra(ctx context.Context, okno string, limit int) ([]MakroPrzegladania, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaMakr)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno, granicaWykazu(limit))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać makr okna %q: %w", okno, err)
	}
	defer wiersze.Close()

	lista := []MakroPrzegladania{}
	for wiersze.Next() {
		makro, err := odczytajMakro(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz makr okna %q: %w", okno, err)
		}
		lista = append(lista, makro)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt makr okna %q: %w", okno, err)
	}
	return lista, nil
}

// ZapiszGranice ustala granice działania Wykonawcy dla pary zasięg + wskazanie.
func (r *repozytoriumPrzegladania) ZapiszGranice(ctx context.Context, granica GranicaWykonawcy) (GranicaWykonawcy, error) {
	if granica.Zasieg == "" {
		return GranicaWykonawcy{}, fmt.Errorf("dane: granice Wykonawcy bez zasięgu")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszGranice)
	if err != nil {
		return GranicaWykonawcy{}, err
	}
	_, err = polecenie.ExecContext(ctx, granica.Zasieg, granica.ZasiegID, granica.MaxKrokow,
		granica.MaxCzasSekund, tekstDoKolumny(granica.DomenyDozwoloneJson),
		tekstDoKolumny(granica.DomenyZablokowaneJson), liczbaLogiczna(granica.PotwierdzajWyslanie))
	if err != nil {
		return GranicaWykonawcy{}, fmt.Errorf("dane: nie można zapisać granic Wykonawcy: %w", err)
	}
	return r.Granice(ctx, granica.Zasieg, granica.ZasiegID)
}

// Granice oddaje granice obowiązujące dla pary zasięg + wskazanie. Brak wiersza
// wraca jako ErrBrakWiersza — wartość domyślną wstawia warstwa wyższa, bo to
// ona wie, ile kroków znaczy „domyślnie".
func (r *repozytoriumPrzegladania) Granice(ctx context.Context, zasieg, zasiegID string) (GranicaWykonawcy, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzGranice)
	if err != nil {
		return GranicaWykonawcy{}, err
	}
	granica, err := odczytajGranice(polecenie.QueryRowContext(ctx, zasieg, zasiegID))
	if errors.Is(err, sql.ErrNoRows) {
		return GranicaWykonawcy{}, ErrBrakWiersza
	}
	if err != nil {
		return GranicaWykonawcy{}, fmt.Errorf("dane: nieczytelne granice Wykonawcy (%s/%s): %w", zasieg, zasiegID, err)
	}
	return granica, nil
}

// odczytajWytwor składa wytwór z jednego wiersza wyniku.
func odczytajWytwor(wiersz skaner) (WytworPrzegladania, error) {
	var wytwor WytworPrzegladania
	var tytul, mime, url, migawka sql.NullString
	var rozmiar sql.NullInt64
	err := wiersz.Scan(&wytwor.ID, &wytwor.Kod, &wytwor.Okno, &wytwor.Rodzaj, &tytul,
		&wytwor.TrescOdwolanie, &mime, &rozmiar, &url, &migawka, &wytwor.Utworzono)
	if err != nil {
		return WytworPrzegladania{}, err
	}
	wytwor.Tytul = tekstZKolumny(tytul)
	wytwor.TypMime = tekstZKolumny(mime)
	wytwor.RozmiarBajtow = liczbaZKolumny(rozmiar)
	wytwor.UrlZrodla = tekstZKolumny(url)
	wytwor.MigawkaZewnetrznaID = tekstZKolumny(migawka)
	return wytwor, nil
}

// odczytajZrzut składa zrzut z jednego wiersza wyniku.
func odczytajZrzut(wiersz skaner) (ZrzutPrzegladania, error) {
	var zrzut ZrzutPrzegladania
	var migawka sql.NullString
	var rozmiar sql.NullInt64
	err := wiersz.Scan(&zrzut.ID, &zrzut.Kod, &zrzut.Okno, &migawka, &zrzut.TrescOdwolanie,
		&zrzut.Tryb, &zrzut.Format, &zrzut.Szerokosc, &zrzut.Wysokosc, &rozmiar, &zrzut.Utworzono)
	if err != nil {
		return ZrzutPrzegladania{}, err
	}
	zrzut.MigawkaZewnetrznaID = tekstZKolumny(migawka)
	zrzut.RozmiarBajtow = liczbaZKolumny(rozmiar)
	return zrzut, nil
}

// odczytajPobranie składa pobranie z jednego wiersza wyniku.
func odczytajPobranie(wiersz skaner) (PobraniePrzegladania, error) {
	var pobranie PobraniePrzegladania
	var nazwa, sciezka, mime, blad, zakonczono sql.NullString
	var odebrano, razem sql.NullInt64
	err := wiersz.Scan(&pobranie.ID, &pobranie.Kod, &pobranie.Okno, &pobranie.Url, &nazwa,
		&sciezka, &mime, &pobranie.Stan, &odebrano, &razem, &blad,
		&pobranie.Rozpoczeto, &zakonczono)
	if err != nil {
		return PobraniePrzegladania{}, err
	}
	pobranie.NazwaPliku = tekstZKolumny(nazwa)
	pobranie.SciezkaDocelowa = tekstZKolumny(sciezka)
	pobranie.TypMime = tekstZKolumny(mime)
	pobranie.OdebranoBajtow = liczbaZKolumny(odebrano)
	pobranie.RazemBajtow = liczbaZKolumny(razem)
	pobranie.KomunikatBledu = tekstZKolumny(blad)
	pobranie.Zakonczono = tekstZKolumny(zakonczono)
	return pobranie, nil
}

// odczytajMakro składa makro z jednego wiersza wyniku.
func odczytajMakro(wiersz skaner) (MakroPrzegladania, error) {
	var makro MakroPrzegladania
	var kroki, automatyka sql.NullString
	var nagrywanie int64
	err := wiersz.Scan(&makro.ID, &makro.Kod, &makro.Okno, &makro.Nazwa, &kroki,
		&nagrywanie, &automatyka, &makro.Utworzono, &makro.Zaktualizowano)
	if err != nil {
		return MakroPrzegladania{}, err
	}
	makro.KrokiJson = tekstZKolumny(kroki)
	makro.AutomatykaKod = tekstZKolumny(automatyka)
	makro.Nagrywanie = nagrywanie == 1
	return makro, nil
}

// odczytajGranice składa granice z jednego wiersza wyniku.
func odczytajGranice(wiersz skaner) (GranicaWykonawcy, error) {
	var granica GranicaWykonawcy
	var dozwolone, zablokowane sql.NullString
	var potwierdzaj int64
	err := wiersz.Scan(&granica.ID, &granica.Zasieg, &granica.ZasiegID, &granica.MaxKrokow,
		&granica.MaxCzasSekund, &dozwolone, &zablokowane, &potwierdzaj, &granica.Zaktualizowano)
	if err != nil {
		return GranicaWykonawcy{}, err
	}
	granica.DomenyDozwoloneJson = tekstZKolumny(dozwolone)
	granica.DomenyZablokowaneJson = tekstZKolumny(zablokowane)
	granica.PotwierdzajWyslanie = potwierdzaj == 1
	return granica, nil
}
