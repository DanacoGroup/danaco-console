// Plik przechowuje dane sesji przeglądania prowadzonej przez rdzeń: wytwory i
// zrzuty ekranu trwałe w magazynie, pobrania plików, zapisane makra oraz
// granice działania roli Wykonawcy dla kolejnych przebiegów.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// WytworPrzegladania to wiersz tabeli wytwor_przegladania: zapisany rezultat
// sesji przeglądania wraz z odwołaniem do treści przechowywanej w magazynie
// plików.
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

// ZrzutPrzegladania to wiersz tabeli zrzut_przegladania: migawka ekranu
// utrwalona w magazynie wraz z trybem, formatem oraz wymiarami obrazu w
// pikselach.
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

// PobraniePrzegladania to wiersz tabeli pobranie_przegladania: stan
// pobierania pliku spod adresu url wraz z postępem w bajtach i ewentualnym
// komunikatem błędu.
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

// MakroPrzegladania to wiersz tabeli makro_przegladania: nagrany ciąg kroków
// możliwy do odtworzenia w kolejnej sesji przeglądania.
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

// GranicaWykonawcy to wiersz tabeli granica_wykonawcy_przegladania: limity
// kroków i czasu działania oraz wykaz domen dozwolonych i zablokowanych dla
// roli Wykonawcy.
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
	                 rozmiar_bajtow, url_zrodla, migawka_zewnetrzna_id, konto_id)
	                VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ` + WskazanieKonta + `)`

	pobierzWytwor = `SELECT ` + kolumnyWytworu + ` FROM wytwor_przegladania
	                 WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	listaWytworow = `SELECT ` + kolumnyWytworu + ` FROM wytwor_przegladania
	                 WHERE okno = ? AND (? = '' OR rodzaj = ?) AND ` + WarunekKonta + `
	                 ORDER BY utworzono DESC, id DESC LIMIT ?`

	kolumnyZrzutu = `id, identyfikator_zewnetrzny, okno, migawka_zewnetrzna_id, tresc_odwolanie,
	                 tryb, format, szerokosc, wysokosc, rozmiar_bajtow, utworzono`

	zapiszZrzut = `INSERT INTO zrzut_przegladania
	               (identyfikator_zewnetrzny, okno, migawka_zewnetrzna_id, tresc_odwolanie,
	                tryb, format, szerokosc, wysokosc, rozmiar_bajtow, konto_id)
	               VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ` + WskazanieKonta + `)`

	pobierzZrzut = `SELECT ` + kolumnyZrzutu + ` FROM zrzut_przegladania
	                WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	pobierzZrzutPoOdwolaniu = `SELECT ` + kolumnyZrzutu + ` FROM zrzut_przegladania
	                           WHERE tresc_odwolanie = ? AND ` + WarunekKonta + `
	                           ORDER BY id DESC LIMIT 1`

	pobierzZrzutMigawki = `SELECT ` + kolumnyZrzutu + ` FROM zrzut_przegladania
	                       WHERE migawka_zewnetrzna_id = ? AND ` + WarunekKonta + `
	                       ORDER BY id DESC LIMIT 1`

	kolumnyPobrania = `id, identyfikator_zewnetrzny, okno, url, nazwa_pliku, sciezka_docelowa,
	                   typ_mime, stan, odebrano_bajtow, razem_bajtow, komunikat_bledu,
	                   rozpoczeto, zakonczono`

	zapiszPobranie = `INSERT INTO pobranie_przegladania
	                  (identyfikator_zewnetrzny, okno, url, nazwa_pliku, sciezka_docelowa, typ_mime,
	                   stan, odebrano_bajtow, razem_bajtow, komunikat_bledu, zakonczono, konto_id)
	                  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ` + WskazanieKonta + `)
	                  ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                      url = excluded.url,
	                      nazwa_pliku = excluded.nazwa_pliku,
	                      sciezka_docelowa = excluded.sciezka_docelowa,
	                      typ_mime = excluded.typ_mime,
	                      stan = excluded.stan,
	                      odebrano_bajtow = excluded.odebrano_bajtow,
	                      razem_bajtow = excluded.razem_bajtow,
	                      komunikat_bledu = excluded.komunikat_bledu,
	                      zakonczono = excluded.zakonczono
	                  WHERE ` + WarunekKonta

	pobierzPobranie = `SELECT ` + kolumnyPobrania + ` FROM pobranie_przegladania
	                   WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	listaPobran = `SELECT ` + kolumnyPobrania + ` FROM pobranie_przegladania
	               WHERE (? = '' OR okno = ?) AND (? = '' OR stan = ?) AND ` + WarunekKonta + `
	               ORDER BY rozpoczeto DESC, id DESC LIMIT ?`

	usunPobranie = `DELETE FROM pobranie_przegladania
	                WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	kolumnyMakra = `id, identyfikator_zewnetrzny, okno, nazwa, kroki_json, nagrywanie,
	                automatyka_zewnetrzna_id, utworzono, zaktualizowano`

	zapiszMakro = `INSERT INTO makro_przegladania
	               (identyfikator_zewnetrzny, okno, nazwa, kroki_json, nagrywanie,
	                automatyka_zewnetrzna_id, konto_id)
	               VALUES (?, ?, ?, ?, ?, ?, ` + WskazanieKonta + `)
	               ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                   nazwa = excluded.nazwa,
	                   kroki_json = excluded.kroki_json,
	                   nagrywanie = excluded.nagrywanie,
	                   automatyka_zewnetrzna_id = excluded.automatyka_zewnetrzna_id,
	                   zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	               WHERE ` + WarunekKonta

	pobierzMakro = `SELECT ` + kolumnyMakra + ` FROM makro_przegladania
	                WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	listaMakr = `SELECT ` + kolumnyMakra + ` FROM makro_przegladania
	             WHERE okno = ? AND ` + WarunekKonta + `
	             ORDER BY utworzono DESC, id DESC LIMIT ?`

	kolumnyGranicy = `id, zasieg, zasieg_id, max_krokow, max_czas_sekund,
	                  domeny_dozwolone_json, domeny_zablokowane_json,
	                  potwierdzaj_wyslanie, zaktualizowano`

	zapiszGranice = `INSERT INTO granica_wykonawcy_przegladania
	                 (zasieg, zasieg_id, max_krokow, max_czas_sekund, domeny_dozwolone_json,
	                  domeny_zablokowane_json, potwierdzaj_wyslanie, konto_id)
	                 VALUES (?, ?, ?, ?, ?, ?, ?, ` + WskazanieKonta + `)
	                 ON CONFLICT(zasieg, zasieg_id, COALESCE(konto_id, 0)) DO UPDATE SET
	                     max_krokow = excluded.max_krokow,
	                     max_czas_sekund = excluded.max_czas_sekund,
	                     domeny_dozwolone_json = excluded.domeny_dozwolone_json,
	                     domeny_zablokowane_json = excluded.domeny_zablokowane_json,
	                     potwierdzaj_wyslanie = excluded.potwierdzaj_wyslanie,
	                     zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                 WHERE ` + WarunekKonta

	pobierzGranice = `SELECT ` + kolumnyGranicy + ` FROM granica_wykonawcy_przegladania
	                  WHERE zasieg = ? AND zasieg_id = ? AND ` + WarunekKonta
)

// ZapiszWytwor odkłada wytwór sesji przeglądania jako kolejny wiersz
// historii, nie nadpisując wcześniej zapisanych wytworów tego samego okna.
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
		tekstDoKolumny(wytwor.MigawkaZewnetrznaID), KontoOperatora(ctx))
	if err != nil {
		return WytworPrzegladania{}, fmt.Errorf("dane: nie można zapisać wytworu %q: %w", wytwor.Kod, err)
	}
	return r.Wytwor(ctx, wytwor.Kod)
}

// Wytwor oddaje wytwór przeglądania o wskazanym kodzie zewnętrznym,
// zwracając błąd ErrBrakWiersza, gdy taki wytwór nie istnieje w magazynie.
func (r *repozytoriumPrzegladania) Wytwor(ctx context.Context, kod string) (WytworPrzegladania, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzWytwor)
	if err != nil {
		return WytworPrzegladania{}, err
	}
	wytwor, err := odczytajWytwor(polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return WytworPrzegladania{}, ErrBrakWiersza
	}
	if err != nil {
		return WytworPrzegladania{}, fmt.Errorf("dane: nieczytelny wytwór %q: %w", kod, err)
	}
	return wytwor, nil
}

// Wytwory oddaje wytwory zapisane dla wskazanego okna, opcjonalnie zawężone
// do jednego rodzaju, uporządkowane od najświeższego według czasu utworzenia.
func (r *repozytoriumPrzegladania) Wytwory(ctx context.Context, okno, rodzaj string, limit int) ([]WytworPrzegladania, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaWytworow)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno, rodzaj, rodzaj,
		KontoOperatora(ctx), granicaWykazu(limit))
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

// ZapiszZrzut odkłada wiersz zrzutu ekranu wraz z trybem, formatem i
// wymiarami obrazu, wiążąc go z magazynem przez odwołanie do treści.
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
		zrzut.Format, zrzut.Szerokosc, zrzut.Wysokosc, liczbaDoKolumny(zrzut.RozmiarBajtow),
		KontoOperatora(ctx))
	if err != nil {
		return ZrzutPrzegladania{}, fmt.Errorf("dane: nie można zapisać zrzutu %q: %w", zrzut.Kod, err)
	}
	return r.Zrzut(ctx, zrzut.Kod)
}

// Zrzut oddaje zrzut ekranu o wskazanym kodzie zewnętrznym, zwracając błąd
// ErrBrakWiersza, gdy taki zrzut nie istnieje w magazynie.
func (r *repozytoriumPrzegladania) Zrzut(ctx context.Context, kod string) (ZrzutPrzegladania, error) {
	return r.jedenZrzut(ctx, pobierzZrzut, kod)
}

// ZrzutPoOdwolaniu oddaje zrzut leżący pod wskazanym odwołaniem magazynu,
// czyli ten sam zrzut, który zwraca punkt browser.snapshot.screenshot.get w
// polu screenshotRef.
func (r *repozytoriumPrzegladania) ZrzutPoOdwolaniu(ctx context.Context, odwolanie string) (ZrzutPrzegladania, error) {
	return r.jedenZrzut(ctx, pobierzZrzutPoOdwolaniu, odwolanie)
}

// ZrzutMigawki oddaje najświeższy zrzut ekranu wykonany przy wskazanej
// migawce zewnętrznej, uporządkowany malejąco według identyfikatora wiersza.
func (r *repozytoriumPrzegladania) ZrzutMigawki(ctx context.Context, migawka string) (ZrzutPrzegladania, error) {
	return r.jedenZrzut(ctx, pobierzZrzutMigawki, migawka)
}

// jedenZrzut wykonuje odczyt jednego zrzutu wskazanym zapytaniem SQL,
// mapując brak wiersza na błąd ErrBrakWiersza wspólny dla całego repozytorium.
func (r *repozytoriumPrzegladania) jedenZrzut(ctx context.Context, zapytanie, wskazanie string) (ZrzutPrzegladania, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return ZrzutPrzegladania{}, err
	}
	zrzut, err := odczytajZrzut(polecenie.QueryRowContext(ctx, wskazanie, KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return ZrzutPrzegladania{}, ErrBrakWiersza
	}
	if err != nil {
		return ZrzutPrzegladania{}, fmt.Errorf("dane: nieczytelny zrzut %q: %w", wskazanie, err)
	}
	return zrzut, nil
}

// ZapiszPobranie zakłada nowe pobranie albo nadpisuje jego stan i postęp,
// gdy pobranie o tym samym kodzie zewnętrznym już istnieje w magazynie.
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
		tekstDoKolumny(pobranie.Zakonczono), KontoOperatora(ctx), KontoOperatora(ctx))
	if err != nil {
		return PobraniePrzegladania{}, fmt.Errorf("dane: nie można zapisać pobrania %q: %w", pobranie.Kod, err)
	}
	return r.Pobranie(ctx, pobranie.Kod)
}

// Pobranie oddaje pobranie o wskazanym kodzie zewnętrznym, zwracając błąd
// ErrBrakWiersza, gdy takie pobranie nie istnieje w magazynie danych.
func (r *repozytoriumPrzegladania) Pobranie(ctx context.Context, kod string) (PobraniePrzegladania, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzPobranie)
	if err != nil {
		return PobraniePrzegladania{}, err
	}
	pobranie, err := odczytajPobranie(polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return PobraniePrzegladania{}, ErrBrakWiersza
	}
	if err != nil {
		return PobraniePrzegladania{}, fmt.Errorf("dane: nieczytelne pobranie %q: %w", kod, err)
	}
	return pobranie, nil
}

// Pobrania oddaje pobrania zapisane w magazynie, opcjonalnie zawężone do
// jednego okna i do jednego stanu, uporządkowane od najświeższego.
func (r *repozytoriumPrzegladania) Pobrania(ctx context.Context, okno, stan string, limit int) ([]PobraniePrzegladania, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaPobran)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno, okno, stan, stan,
		KontoOperatora(ctx), granicaWykazu(limit))
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

// UsunPobranie zdejmuje pobranie o wskazanym kodzie zewnętrznym z wykazu,
// oddając informację, czy wiersz istniał przed usunięciem.
func (r *repozytoriumPrzegladania) UsunPobranie(ctx context.Context, kod string) (bool, error) {
	return r.usunWiersz(ctx, usunPobranie, kod, "pobranie", KontoOperatora(ctx))
}

// ZapiszMakro zakłada nowe makro albo nadpisuje zastane wraz z krokami, gdy
// makro o tym samym kodzie zewnętrznym już istnieje w magazynie.
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
		tekstDoKolumny(makro.AutomatykaKod), KontoOperatora(ctx), KontoOperatora(ctx))
	if err != nil {
		return MakroPrzegladania{}, fmt.Errorf("dane: nie można zapisać makra %q: %w", makro.Kod, err)
	}
	return r.Makro(ctx, makro.Kod)
}

// Makro oddaje makro przeglądania o wskazanym kodzie zewnętrznym, zwracając
// błąd ErrBrakWiersza, gdy takie makro nie istnieje w magazynie.
func (r *repozytoriumPrzegladania) Makro(ctx context.Context, kod string) (MakroPrzegladania, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzMakro)
	if err != nil {
		return MakroPrzegladania{}, err
	}
	makro, err := odczytajMakro(polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return MakroPrzegladania{}, ErrBrakWiersza
	}
	if err != nil {
		return MakroPrzegladania{}, fmt.Errorf("dane: nieczytelne makro %q: %w", kod, err)
	}
	return makro, nil
}

// Makra oddaje makra zapisane dla wskazanego okna, uporządkowane od
// najświeższego według czasu utworzenia i identyfikatora wiersza.
func (r *repozytoriumPrzegladania) Makra(ctx context.Context, okno string, limit int) ([]MakroPrzegladania, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaMakr)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno, KontoOperatora(ctx), granicaWykazu(limit))
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

// ZapiszGranice ustala granice działania Wykonawcy dla pary zasięg i
// wskazanie, nadpisując wartości zastane dla tej samej pary kolumn.
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
		tekstDoKolumny(granica.DomenyZablokowaneJson), liczbaLogiczna(granica.PotwierdzajWyslanie),
		KontoOperatora(ctx), KontoOperatora(ctx))
	if err != nil {
		return GranicaWykonawcy{}, fmt.Errorf("dane: nie można zapisać granic Wykonawcy: %w", err)
	}
	return r.Granice(ctx, granica.Zasieg, granica.ZasiegID)
}

// Granice oddaje granice obowiązujące dla pary zasięg i wskazanie. Brak
// wiersza wraca jako ErrBrakWiersza, ponieważ wartość domyślna należy do
// warstwy wyższej, która zna liczbę kroków domyślnych.
func (r *repozytoriumPrzegladania) Granice(ctx context.Context, zasieg, zasiegID string) (GranicaWykonawcy, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzGranice)
	if err != nil {
		return GranicaWykonawcy{}, err
	}
	granica, err := odczytajGranice(polecenie.QueryRowContext(ctx, zasieg, zasiegID, KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return GranicaWykonawcy{}, ErrBrakWiersza
	}
	if err != nil {
		return GranicaWykonawcy{}, fmt.Errorf("dane: nieczytelne granice Wykonawcy (%s/%s): %w", zasieg, zasiegID, err)
	}
	return granica, nil
}

// odczytajWytwor składa wytwór przeglądania z jednego wiersza wyniku
// zapytania, zamieniając kolumny nullowalne na wskaźniki opcjonalne.
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

// odczytajZrzut składa zrzut ekranu z jednego wiersza wyniku zapytania,
// zamieniając kolumny nullowalne na wskaźniki opcjonalne struktury.
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

// odczytajPobranie składa pobranie z jednego wiersza wyniku zapytania,
// zamieniając kolumny nullowalne na wskaźniki opcjonalne struktury.
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

// odczytajMakro składa makro przeglądania z jednego wiersza wyniku
// zapytania, zamieniając znacznik liczbowy nagrywania na wartość logiczną.
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

// odczytajGranice składa granice Wykonawcy z jednego wiersza wyniku
// zapytania, zamieniając znacznik liczbowy potwierdzenia na wartość logiczną.
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
