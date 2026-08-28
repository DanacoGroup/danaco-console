// Plik jest magazynem śladu wywołań kanału modelu: zapis wywołania wraz z odcinkami oraz
// odczyt zawężony do pytania operatora. Uzasadnienie modelu danych niesie rozdział
// prowenancja.go dokumentacji architektury.
package dane

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// WywolanieModelu jest jednym wierszem śladu wywołania kanału modelu. Wskaźniki niosą
// brak wartości jako stan różny od zera: wywołanie biegnące nie ma jeszcze końca,
// a wywołanie bez cennika nie ma kosztu.
type WywolanieModelu struct {
	ID               int64
	Kod              string
	SladKod          *string
	RodzicKod        *string
	ProcesKod        *string
	SesjaKod         *string
	OknoKod          *string
	WiadomoscKod     *string
	KanalKod         *string
	Dostawca         *string
	Model            *string
	KontoKod         *string
	ProjektKod       *string
	Srodowisko       *string
	Stan             string
	KodBledu         *string
	Poczatek         int64
	Koniec           *int64
	OpoznienieMs     *int
	TokenyPromptu    *int
	TokenyOdpowiedzi *int
	TokenyCache      *int
	TokenyRazem      *int
	Koszt            *float64
	Waluta           *string
	LiczbaNarzedzi   int
	LiczbaPodagentow int
	Ocena            string
	OcenaNotatka     *string
	TrescZapisana    bool
	Zredagowane      bool
	Prompt           *string
	Odpowiedz        *string
}

// OdcinekWywolania jest jedną gałęzią drzewa wywołania: pojedynczym krokiem
// przetwarzania wraz z jego czasem, kosztem i stanem.
type OdcinekWywolania struct {
	ID             int64
	Kod            string
	WywolanieKod   string
	RodzicKod      *string
	Nazwa          string
	Rodzaj         string
	Poczatek       int64
	Koniec         *int64
	CzasMs         *int
	Tokeny         *int
	Koszt          *float64
	Stan           string
	KodBledu       *string
	TrafienieCache bool
	Atrybuty       *string
}

// SitoWywolan zawęża wykaz. Pole puste znaczy „bez zawężenia po tej osi" —
// nie „wartość pusta", bo o wywołanie bez sesji też można zapytać wprost.
type SitoWywolan struct {
	Od        *int64
	Do        *int64
	SesjaKod  string
	OknoKod   string
	ProcesKod string
	SladKod   string
	KanalKod  string
	Model     string
	KontoKod  string
	Stan      string
	Fraza     string
	MinKoszt  *float64
	Granica   int
}

// SumaZuzycia jest jednym wierszem rozliczenia zużycia po wymiarze: liczbą zadań,
// tokenami i kosztem zsumowanymi dla jednego klucza.
type SumaZuzycia struct {
	Wymiar            string
	Klucz             string
	Zadania           int
	ZadaniaBledne     int
	TokenyPromptu     int
	TokenyOdpowiedzi  int
	TokenyCache       int
	TokenyRazem       int
	Koszt             float64
	BezCeny           int
	SrednieOpoznienie int
}

// RepozytoriumProwenancji jest kontraktem magazynu śladu wywołań modelu wraz z odcinkami
// i rozliczeniem zużycia.
type RepozytoriumProwenancji interface {
	// ZapiszWywolanie wnosi wiersz albo nadpisuje istniejący po kodzie.
	ZapiszWywolanie(ctx context.Context, wywolanie WywolanieModelu) (WywolanieModelu, error)
	// Wywolanie oddaje jeden ślad po kodzie.
	Wywolanie(ctx context.Context, kod string) (WywolanieModelu, error)
	// Wywolania oddaje wykaz zawężony sitem wraz z liczbą wszystkich pasujących.
	Wywolania(ctx context.Context, sito SitoWywolan) ([]WywolanieModelu, int, error)
	// ZapiszOdcinek wnosi gałąź drzewa wywołania.
	ZapiszOdcinek(ctx context.Context, odcinek OdcinekWywolania) (OdcinekWywolania, error)
	// Odcinki oddaje drzewo jednego wywołania w kolejności rozpoczęcia.
	Odcinki(ctx context.Context, wywolanieKod string) ([]OdcinekWywolania, error)
	// OcenWywolanie zapisuje zdanie Operatora o wywołaniu.
	OcenWywolanie(ctx context.Context, kod, ocena string, notatka *string) (WywolanieModelu, error)
	// Zuzycie liczy sumy po wskazanym wymiarze w zadanym oknie czasu.
	Zuzycie(ctx context.Context, wymiar string, od, do *int64, granica int) ([]SumaZuzycia, error)
}

// kolumnaWymiaru przekłada wymiar kontraktu na kolumnę tabeli. Wymiar spoza wykazu jest
// odmową, nie cichym zejściem do wymiaru domyślnego, żeby suma policzona po innej osi
// nie wyglądała identycznie jak zamówiona.
var kolumnaWymiaru = map[string]string{
	"channel":     "kanal_kod",
	"provider":    "dostawca",
	"account":     "konto_kod",
	"session":     "sesja_kod",
	"project":     "projekt_kod",
	"environment": "srodowisko",
	"window":      "okno_kod",
}

type repozytoriumProwenancji struct {
	db *sql.DB
	z  *zapytania
}

func noweRepozytoriumProwenancji(z *zapytania, db *sql.DB) *repozytoriumProwenancji {
	return &repozytoriumProwenancji{db: db, z: z}
}

const kolumnyWywolania = `kod, slad_kod, rodzic_kod, proces_kod, sesja_kod, okno_kod,
	wiadomosc_kod, kanal_kod, dostawca, model, konto_kod, projekt_kod, srodowisko,
	stan, kod_bledu, poczatek, koniec, opoznienie_ms, tokeny_promptu,
	tokeny_odpowiedzi, tokeny_cache, tokeny_razem, koszt, waluta,
	liczba_narzedzi, liczba_podagentow, ocena, ocena_notatka,
	tresc_zapisana, zredagowane, prompt, odpowiedz`

func (r *repozytoriumProwenancji) odczytajWywolanie(s skaner) (WywolanieModelu, error) {
	var w WywolanieModelu
	err := s.Scan(&w.Kod, &w.SladKod, &w.RodzicKod, &w.ProcesKod, &w.SesjaKod, &w.OknoKod,
		&w.WiadomoscKod, &w.KanalKod, &w.Dostawca, &w.Model, &w.KontoKod, &w.ProjektKod,
		&w.Srodowisko, &w.Stan, &w.KodBledu, &w.Poczatek, &w.Koniec, &w.OpoznienieMs,
		&w.TokenyPromptu, &w.TokenyOdpowiedzi, &w.TokenyCache, &w.TokenyRazem, &w.Koszt,
		&w.Waluta, &w.LiczbaNarzedzi, &w.LiczbaPodagentow, &w.Ocena, &w.OcenaNotatka,
		&w.TrescZapisana, &w.Zredagowane, &w.Prompt, &w.Odpowiedz)
	return w, err
}

func (r *repozytoriumProwenancji) ZapiszWywolanie(ctx context.Context,
	w WywolanieModelu) (WywolanieModelu, error) {

	if strings.TrimSpace(w.Kod) == "" {
		return WywolanieModelu{}, fmt.Errorf("dane: wywołanie bez kodu nie ma czym się zapisać")
	}
	if w.Ocena == "" {
		w.Ocena = "bez_oceny"
	}
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO prowenancja_wywolanie (`+kolumnyWywolania+`)
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
		ON CONFLICT(kod) DO UPDATE SET
			stan = excluded.stan, kod_bledu = excluded.kod_bledu,
			koniec = excluded.koniec, opoznienie_ms = excluded.opoznienie_ms,
			tokeny_promptu = excluded.tokeny_promptu,
			tokeny_odpowiedzi = excluded.tokeny_odpowiedzi,
			tokeny_cache = excluded.tokeny_cache, tokeny_razem = excluded.tokeny_razem,
			koszt = excluded.koszt, waluta = excluded.waluta,
			liczba_narzedzi = excluded.liczba_narzedzi,
			liczba_podagentow = excluded.liczba_podagentow,
			tresc_zapisana = excluded.tresc_zapisana, zredagowane = excluded.zredagowane,
			prompt = excluded.prompt, odpowiedz = excluded.odpowiedz`,
		w.Kod, w.SladKod, w.RodzicKod, w.ProcesKod, w.SesjaKod, w.OknoKod, w.WiadomoscKod,
		w.KanalKod, w.Dostawca, w.Model, w.KontoKod, w.ProjektKod, w.Srodowisko, w.Stan,
		w.KodBledu, w.Poczatek, w.Koniec, w.OpoznienieMs, w.TokenyPromptu, w.TokenyOdpowiedzi,
		w.TokenyCache, w.TokenyRazem, w.Koszt, w.Waluta, w.LiczbaNarzedzi, w.LiczbaPodagentow,
		w.Ocena, w.OcenaNotatka, w.TrescZapisana, w.Zredagowane, w.Prompt, w.Odpowiedz)
	if err != nil {
		return WywolanieModelu{}, fmt.Errorf("dane: zapis wywołania %s: %w", w.Kod, err)
	}
	return r.Wywolanie(ctx, w.Kod)
}

func (r *repozytoriumProwenancji) Wywolanie(ctx context.Context, kod string) (WywolanieModelu, error) {
	wiersz := r.db.QueryRowContext(ctx,
		`SELECT `+kolumnyWywolania+` FROM prowenancja_wywolanie WHERE kod = ?`, kod)
	w, err := r.odczytajWywolanie(wiersz)
	if err == sql.ErrNoRows {
		return WywolanieModelu{}, ErrBrakWiersza
	}
	if err != nil {
		return WywolanieModelu{}, fmt.Errorf("dane: odczyt wywołania %s: %w", kod, err)
	}
	return w, nil
}

// warunkiSita składa listę warunków wraz z argumentami, wspólną dla wykazu i licznika
// wszystkich pasujących, żeby licznik nie liczył po innym zawężeniu niż wykaz.
func warunkiSita(sito SitoWywolan) ([]string, []any) {
	var warunki []string
	var argumenty []any
	dodaj := func(warunek string, wartosc any) {
		warunki = append(warunki, warunek)
		argumenty = append(argumenty, wartosc)
	}
	if sito.Od != nil {
		dodaj("poczatek >= ?", *sito.Od)
	}
	if sito.Do != nil {
		dodaj("poczatek <= ?", *sito.Do)
	}
	for _, para := range []struct {
		kolumna string
		wartosc string
	}{
		{"sesja_kod", sito.SesjaKod}, {"okno_kod", sito.OknoKod},
		{"proces_kod", sito.ProcesKod}, {"slad_kod", sito.SladKod},
		{"kanal_kod", sito.KanalKod}, {"model", sito.Model},
		{"konto_kod", sito.KontoKod}, {"stan", sito.Stan},
	} {
		if strings.TrimSpace(para.wartosc) != "" {
			dodaj(para.kolumna+" = ?", para.wartosc)
		}
	}
	if sito.MinKoszt != nil {
		dodaj("koszt >= ?", *sito.MinKoszt)
	}
	if fraza := strings.TrimSpace(sito.Fraza); fraza != "" {
		warunki = append(warunki, "(model LIKE ? OR kanal_kod LIKE ? OR kod_bledu LIKE ?)")
		wzor := "%" + fraza + "%"
		argumenty = append(argumenty, wzor, wzor, wzor)
	}
	return warunki, argumenty
}

func (r *repozytoriumProwenancji) Wywolania(ctx context.Context,
	sito SitoWywolan) ([]WywolanieModelu, int, error) {

	warunki, argumenty := warunkiSita(sito)
	gdzie := ""
	if len(warunki) > 0 {
		gdzie = " WHERE " + strings.Join(warunki, " AND ")
	}

	var wszystkie int
	if err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM prowenancja_wywolanie`+gdzie, argumenty...).Scan(&wszystkie); err != nil {
		return nil, 0, fmt.Errorf("dane: licznik wywołań: %w", err)
	}

	granica := sito.Granica
	if granica <= 0 {
		granica = 100
	}
	wiersze, err := r.db.QueryContext(ctx,
		`SELECT `+kolumnyWywolania+` FROM prowenancja_wywolanie`+gdzie+
			` ORDER BY poczatek DESC LIMIT ?`, append(argumenty, granica)...)
	if err != nil {
		return nil, 0, fmt.Errorf("dane: wykaz wywołań: %w", err)
	}
	defer wiersze.Close()

	var wykaz []WywolanieModelu
	for wiersze.Next() {
		w, err := r.odczytajWywolanie(wiersze)
		if err != nil {
			return nil, 0, fmt.Errorf("dane: nieczytelny wiersz wywołania: %w", err)
		}
		wykaz = append(wykaz, w)
	}
	return wykaz, wszystkie, wiersze.Err()
}

func (r *repozytoriumProwenancji) ZapiszOdcinek(ctx context.Context,
	o OdcinekWywolania) (OdcinekWywolania, error) {

	if strings.TrimSpace(o.Kod) == "" || strings.TrimSpace(o.WywolanieKod) == "" {
		return OdcinekWywolania{}, fmt.Errorf("dane: odcinek bez kodu albo bez wywołania")
	}
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO prowenancja_odcinek
			(kod, wywolanie_kod, rodzic_kod, nazwa, rodzaj, poczatek, koniec, czas_ms,
			 tokeny, koszt, stan, kod_bledu, trafienie_cache, atrybuty)
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)
		ON CONFLICT(kod) DO UPDATE SET
			koniec = excluded.koniec, czas_ms = excluded.czas_ms,
			tokeny = excluded.tokeny, koszt = excluded.koszt,
			stan = excluded.stan, kod_bledu = excluded.kod_bledu,
			trafienie_cache = excluded.trafienie_cache, atrybuty = excluded.atrybuty`,
		o.Kod, o.WywolanieKod, o.RodzicKod, o.Nazwa, o.Rodzaj, o.Poczatek, o.Koniec,
		o.CzasMs, o.Tokeny, o.Koszt, o.Stan, o.KodBledu, o.TrafienieCache, o.Atrybuty)
	if err != nil {
		return OdcinekWywolania{}, fmt.Errorf("dane: zapis odcinka %s: %w", o.Kod, err)
	}
	return o, nil
}

func (r *repozytoriumProwenancji) Odcinki(ctx context.Context,
	wywolanieKod string) ([]OdcinekWywolania, error) {

	wiersze, err := r.db.QueryContext(ctx, `
		SELECT kod, wywolanie_kod, rodzic_kod, nazwa, rodzaj, poczatek, koniec, czas_ms,
		       tokeny, koszt, stan, kod_bledu, trafienie_cache, atrybuty
		  FROM prowenancja_odcinek WHERE wywolanie_kod = ? ORDER BY poczatek`, wywolanieKod)
	if err != nil {
		return nil, fmt.Errorf("dane: odcinki wywołania %s: %w", wywolanieKod, err)
	}
	defer wiersze.Close()

	var wykaz []OdcinekWywolania
	for wiersze.Next() {
		var o OdcinekWywolania
		if err := wiersze.Scan(&o.Kod, &o.WywolanieKod, &o.RodzicKod, &o.Nazwa, &o.Rodzaj,
			&o.Poczatek, &o.Koniec, &o.CzasMs, &o.Tokeny, &o.Koszt, &o.Stan, &o.KodBledu,
			&o.TrafienieCache, &o.Atrybuty); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny odcinek: %w", err)
		}
		wykaz = append(wykaz, o)
	}
	return wykaz, wiersze.Err()
}

func (r *repozytoriumProwenancji) OcenWywolanie(ctx context.Context,
	kod, ocena string, notatka *string) (WywolanieModelu, error) {

	wynik, err := r.db.ExecContext(ctx,
		`UPDATE prowenancja_wywolanie SET ocena = ?, ocena_notatka = ? WHERE kod = ?`,
		ocena, notatka, kod)
	if err != nil {
		return WywolanieModelu{}, fmt.Errorf("dane: ocena wywołania %s: %w", kod, err)
	}
	if zmienione, err := wynik.RowsAffected(); err == nil && zmienione == 0 {
		return WywolanieModelu{}, ErrBrakWiersza
	}
	return r.Wywolanie(ctx, kod)
}

func (r *repozytoriumProwenancji) Zuzycie(ctx context.Context, wymiar string,
	od, do *int64, granica int) ([]SumaZuzycia, error) {

	kolumna, znany := kolumnaWymiaru[wymiar]
	if !znany {
		return nil, fmt.Errorf("dane: wymiar rozliczenia %q nie ma odpowiednika w schemacie", wymiar)
	}

	var warunki []string
	var argumenty []any
	if od != nil {
		warunki = append(warunki, "poczatek >= ?")
		argumenty = append(argumenty, *od)
	}
	if do != nil {
		warunki = append(warunki, "poczatek <= ?")
		argumenty = append(argumenty, *do)
	}
	gdzie := ""
	if len(warunki) > 0 {
		gdzie = " WHERE " + strings.Join(warunki, " AND ")
	}
	if granica <= 0 {
		granica = 50
	}

	// Wywołanie bez wymiaru wpada do klucza pustego zamiast wypadać z sumy rozliczenia.
	wiersze, err := r.db.QueryContext(ctx, `
		SELECT COALESCE(`+kolumna+`, '') AS klucz,
		       COUNT(*),
		       SUM(CASE WHEN stan = 'bledne' THEN 1 ELSE 0 END),
		       COALESCE(SUM(tokeny_promptu), 0), COALESCE(SUM(tokeny_odpowiedzi), 0),
		       COALESCE(SUM(tokeny_cache), 0), COALESCE(SUM(tokeny_razem), 0),
		       COALESCE(SUM(koszt), 0),
		       SUM(CASE WHEN koszt IS NULL THEN 1 ELSE 0 END),
		       COALESCE(CAST(AVG(opoznienie_ms) AS INTEGER), 0)
		  FROM prowenancja_wywolanie`+gdzie+`
		 GROUP BY klucz ORDER BY SUM(COALESCE(koszt, 0)) DESC, COUNT(*) DESC LIMIT ?`,
		append(argumenty, granica)...)
	if err != nil {
		return nil, fmt.Errorf("dane: rozliczenie zużycia po %s: %w", wymiar, err)
	}
	defer wiersze.Close()

	var wykaz []SumaZuzycia
	for wiersze.Next() {
		suma := SumaZuzycia{Wymiar: wymiar}
		if err := wiersze.Scan(&suma.Klucz, &suma.Zadania, &suma.ZadaniaBledne,
			&suma.TokenyPromptu, &suma.TokenyOdpowiedzi, &suma.TokenyCache,
			&suma.TokenyRazem, &suma.Koszt, &suma.BezCeny,
			&suma.SrednieOpoznienie); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz rozliczenia: %w", err)
		}
		wykaz = append(wykaz, suma)
	}
	return wykaz, wiersze.Err()
}
