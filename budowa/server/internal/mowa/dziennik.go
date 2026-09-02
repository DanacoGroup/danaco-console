// Plik utrwala dziennik transkrypcji mowy: wpisy tabeli `transkrypcja` w stanach
// gotowa, bez_mowy i odmowa, każdy z czasem zapisu podanym przez wołającego.
package mowa

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"danacoconsole/server/internal/dane"
)

// Wartości kolumny `stan`, przepisane z CHECK-u tabeli `transkrypcja`. Napis
// w wywołaniu byłby kolejnym miejscem, w którym ta sama wartość musiałaby się
// zgadzać co do znaku.
const (
	StanGotowa = "gotowa"
	// StanZapisuBezMowy — przetworzono, mowy w nagraniu nie ma. To wynik,
	// nie odmowa: pomiar się odbył i właśnie tyle zmierzył. Bez osobnej
	// wartości czytający musiałby wnioskować z `znakow = 0`, a to zgadywanie.
	StanZapisuBezMowy = "bez_mowy"
	StanOdmowa        = "odmowa"
)

// LimitWykazuDomyslny obowiązuje, gdy filtr nie niesie własnego. Wykaz bez
// żadnej granicy oddawałby dziennik od początku instalacji — a pytanie brzmi
// „co ostatnio przepisano", nie „co kiedykolwiek".
const LimitWykazuDomyslny = 50

// WpisTranskrypcji to wiersz tabeli `transkrypcja`: jedno zlecenie rozpoznania
// mowy, jego wynik albo powód odmowy, wraz z czasem zapisu.
type WpisTranskrypcji struct {
	// ID jest kluczem sztucznym bazy; na zewnątrz idzie Identyfikator.
	ID int64
	// Identyfikator odpowiada kolumnie `identyfikator_zewnetrzny` — tożsamość
	// wpisu widoczna poza bazą.
	Identyfikator string
	// OknoId bywa pusty dla transkrypcji spoza okna komunikacji; w bazie
	// odpowiada mu NULL.
	OknoId string
	// NagranieOdnosnik to ścieżka pliku na dysku, nie bajty dźwięku.
	NagranieOdnosnik string
	// Model i Jezyk niosą wartości obowiązujące w chwili zlecenia.
	Model string
	Jezyk string
	// Znakow to długość rozpoznanego tekstu; samego tekstu dziennik nie trzyma.
	Znakow int64
	// TrwanieMs to długość nagrania w milisekundach, nie czas przetwarzania.
	TrwanieMs int64
	// Stan niesie StanGotowa, StanZapisuBezMowy albo StanOdmowa.
	Stan string
	// Powod wypełniony jest wyłącznie przy odmowie; przy wpisie udanym musi
	// być pusty.
	Powod string
	// Utworzono to milisekundy epoki podane przez wołającego.
	Utworzono int64
}

// FiltrTranskrypcji zawęża wykaz dziennika transkrypcji do jednego okna
// komunikacji i do zadanej liczby najnowszych wpisów.
type FiltrTranskrypcji struct {
	// OknoId pusty znaczy „wszystkie okna, także wpisy spoza okna".
	OknoId string
	// Limit niedodatni znaczy LimitWykazuDomyslny, nie całość dziennika.
	Limit int
}

// Dziennik jest kontraktem trwałości transkrypcji — wołający zna
// interfejs, nie strukturę, i nie zna SQL-a.
type Dziennik interface {
	// Zapisz wstawia wpis i oddaje go po zapisie, razem z nadanym ID.
	Zapisz(ctx context.Context, wpis WpisTranskrypcji) (WpisTranskrypcji, error)
	// Wykaz oddaje wpisy od najnowszego.
	Wykaz(ctx context.Context, filtr FiltrTranskrypcji) ([]WpisTranskrypcji, error)
}

const (
	kolumnyTranskrypcji = `t.id, t.identyfikator_zewnetrzny, t.okno_id, t.nagranie_odnosnik,
	                       t.model, t.jezyk, t.znakow, t.trwanie_ms, t.stan, t.powod,
	                       t.utworzono`

	zrodloTranskrypcji = ` FROM transkrypcja t`

	pobierzTranskrypcje = `SELECT ` + kolumnyTranskrypcji + zrodloTranskrypcji +
		` WHERE t.identyfikator_zewnetrzny = ? AND ` + dane.WarunekKonta

	// Zapytanie obsługuje oba warianty żądania: puste zawężenie okna wyłącza
	// pierwszy warunek, a wskazane oddaje wpisy tylko z tego okna, obie
	// gałęzie od najnowszego wpisu.
	listaTranskrypcji = `SELECT ` + kolumnyTranskrypcji + zrodloTranskrypcji +
		` WHERE (? = '' OR t.okno_id = ?) AND ` + dane.WarunekKonta + `
		  ORDER BY t.utworzono DESC, t.id DESC
		  LIMIT ?`

	wstawTranskrypcje = `INSERT INTO transkrypcja
	                     (identyfikator_zewnetrzny, okno_id, nagranie_odnosnik, model,
	                      jezyk, znakow, trwanie_ms, stan, powod, utworzono, konto_id)
	                     VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ` + dane.WskazanieKonta + `)`
)

// dziennikTranskrypcji jest jedyną implementacją Dziennika, opartą wprost
// na puli połączeń bazy danych.
type dziennikTranskrypcji struct {
	db *sql.DB
}

// Zgodność implementacji z kontraktem sprawdza kompilator, a nie dopiero
// pierwsze wywołanie w czasie pracy.
var _ Dziennik = (*dziennikTranskrypcji)(nil)

// NowyDziennik zakłada dziennik nad otwartą pulą połączeń. Bierze `*sql.DB`,
// a nie pamięć zapytań pakietu `dane`, bo ta jest nieeksportowana.
func NowyDziennik(db *sql.DB) Dziennik {
	if db == nil {
		return nil
	}
	return &dziennikTranskrypcji{db: db}
}

// Zapisz wstawia wpis dziennika i oddaje go odczytany z bazy, razem
// z nadanym identyfikatorem sztucznym.
func (d *dziennikTranskrypcji) Zapisz(ctx context.Context,
	wpis WpisTranskrypcji) (WpisTranskrypcji, error) {

	if err := sprawdzWpis(wpis); err != nil {
		return WpisTranskrypcji{}, err
	}
	if _, err := d.db.ExecContext(ctx, wstawTranskrypcje, wpis.Identyfikator,
		napisDoKolumny(wpis.OknoId), wpis.NagranieOdnosnik, wpis.Model, wpis.Jezyk,
		wpis.Znakow, wpis.TrwanieMs, wpis.Stan, napisDoKolumny(wpis.Powod),
		wpis.Utworzono, dane.KontoOperatora(ctx)); err != nil {

		return WpisTranskrypcji{}, fmt.Errorf("mowa: nie można zapisać transkrypcji %q: %w",
			wpis.Identyfikator, err)
	}
	return d.wpisPoIdentyfikatorze(ctx, wpis.Identyfikator)
}

// Wykaz oddaje wpisy dziennika od najnowszego, zawężone filtrem okna
// komunikacji i liczbą zwracanych wpisów.
func (d *dziennikTranskrypcji) Wykaz(ctx context.Context,
	filtr FiltrTranskrypcji) ([]WpisTranskrypcji, error) {

	limit := filtr.Limit
	if limit <= 0 {
		limit = LimitWykazuDomyslny
	}
	wiersze, err := d.db.QueryContext(ctx, listaTranskrypcji, filtr.OknoId, filtr.OknoId,
		dane.KontoOperatora(ctx), limit)
	if err != nil {
		return nil, fmt.Errorf("mowa: nie można odczytać dziennika transkrypcji: %w", err)
	}
	defer wiersze.Close()

	dziennik := make([]WpisTranskrypcji, 0, 16)
	for wiersze.Next() {
		wpis, err := odczytajTranskrypcje(wiersze)
		if err != nil {
			return nil, fmt.Errorf("mowa: uszkodzony wiersz transkrypcji: %w", err)
		}
		dziennik = append(dziennik, wpis)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("mowa: nie można odczytać dziennika transkrypcji: %w", err)
	}
	return dziennik, nil
}

// wpisPoIdentyfikatorze odczytuje jeden wpis. Brak wiersza wraca jako
// ErrBrakWpisu — warstwa wyżej odróżnia „nie ma" od „odczyt padł", tak samo jak
// przy `dane.ErrBrakWiersza`.
func (d *dziennikTranskrypcji) wpisPoIdentyfikatorze(ctx context.Context,
	identyfikator string) (WpisTranskrypcji, error) {

	wpis, err := odczytajTranskrypcje(d.db.QueryRowContext(ctx, pobierzTranskrypcje,
		identyfikator, dane.KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return WpisTranskrypcji{}, fmt.Errorf("mowa: transkrypcja %q nie istnieje: %w",
			identyfikator, ErrBrakWpisu)
	}
	if err != nil {
		return WpisTranskrypcji{}, fmt.Errorf("mowa: nie można odczytać transkrypcji %q: %w",
			identyfikator, err)
	}
	return wpis, nil
}

// ErrBrakWpisu odróżnia brak wiersza od awarii odczytu. Wartownik, a nie
// pusta struktura: wołający sprawdza go przez `errors.Is` i nie musi zgadywać,
// czy zerowy wpis znaczy „nie ma", czy „coś się popsuło".
var ErrBrakWpisu = errors.New("mowa: brak wpisu dziennika")

// sprawdzWpis odrzuca wpisy, których baza i tak by nie przyjęła, żeby
// wołający dostał zdanie po polsku zamiast komunikatu o naruszonym więzie.
func sprawdzWpis(wpis WpisTranskrypcji) error {
	if wpis.Identyfikator == "" {
		return fmt.Errorf("mowa: transkrypcja bez identyfikatora")
	}
	if wpis.NagranieOdnosnik == "" {
		return fmt.Errorf("mowa: transkrypcja %q bez odnośnika nagrania", wpis.Identyfikator)
	}
	switch wpis.Stan {
	case StanGotowa, StanZapisuBezMowy:
		if wpis.Powod != "" {
			return fmt.Errorf("mowa: transkrypcja %q udana, a niesie powód odmowy",
				wpis.Identyfikator)
		}
	case StanOdmowa:
		if wpis.Powod == "" {
			return fmt.Errorf("mowa: transkrypcja %q odmówiona bez nazwanego powodu",
				wpis.Identyfikator)
		}
	default:
		return fmt.Errorf("mowa: transkrypcja %q o nieznanym stanie %q",
			wpis.Identyfikator, wpis.Stan)
	}
	if wpis.Utworzono <= 0 {
		return fmt.Errorf("mowa: transkrypcja %q bez czasu zapisu — czas podaje wołający",
			wpis.Identyfikator)
	}
	return nil
}

// skaner obejmuje `*sql.Row` i `*sql.Rows`, żeby odczyt wiersza wyglądał tak
// samo niezależnie od tego, czy zapytanie zwraca jeden wiersz, czy wiele.
type skaner interface {
	Scan(cele ...any) error
}

// odczytajTranskrypcje składa wpis dziennika z jednego wiersza wyniku
// zapytania o transkrypcję, niezależnie od liczby zwróconych wierszy.
func odczytajTranskrypcje(wiersz skaner) (WpisTranskrypcji, error) {
	var wpis WpisTranskrypcji
	var oknoId, powod sql.NullString
	if err := wiersz.Scan(&wpis.ID, &wpis.Identyfikator, &oknoId, &wpis.NagranieOdnosnik,
		&wpis.Model, &wpis.Jezyk, &wpis.Znakow, &wpis.TrwanieMs, &wpis.Stan, &powod,
		&wpis.Utworzono); err != nil {

		return WpisTranskrypcji{}, err
	}
	wpis.OknoId = oknoId.String
	wpis.Powod = powod.String
	return wpis, nil
}

// napisDoKolumny przekłada napis pusty na NULL, bo w kolumnach `okno_id`
// i `powod` pustka niesie odrębne znaczenie od wartości podanej wprost.
func napisDoKolumny(wartosc string) any {
	if wartosc == "" {
		return nil
	}
	return wartosc
}
