// Odpowiedzialność pliku: trwałość dziennika transkrypcji — tabela
// `transkrypcja` z migracji `store/migracja_075_mowa.sql`.
//
// Plik powiela kształt repozytoriów pakietu `dane` (blok stałych ze składanym
// SQL, typ wiersza, filtr, interfejs repozytorium, asercja `var _`, odczyt
// wiersza przez interfejs skanera), ale stoi w pakiecie `mowa` i bierze
// `*sql.DB` wprost. Stąd własny interfejs skanera: odpowiednik `dane.skaner`
// jest nieeksportowany i nie da się go użyć spoza pakietu `dane`.
//
// Dziennik nie ma własnego zegara — czas przychodzi z góry w polu `Utworzono`
// wpisu, a baza go nie wstawia (kolumna nie ma wyrażenia domyślnego). Jeden
// zegar na jedno zdarzenie: drugi rozjeżdżałby ślad z chwilą zapisu.
//
// Dziennik zapisuje także odmowy: `Zapisz` przyjmuje wpis o stanie `odmowa`
// i jest to jego zwykłe użycie, bo powód odmowy jest pierwszą informacją
// potrzebną, gdy nic się nie przepisało. Więz spójności stanu z powodem
// pilnuje CHECK tabeli `transkrypcja`, a sprawdzenia niżej odrzucają niespójny
// wpis wcześniej, żeby wołający dostał zdanie po polsku zamiast komunikatu
// sterownika SQLite.
package mowa

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

// WpisTranskrypcji to wiersz tabeli `transkrypcja`.
type WpisTranskrypcji struct {
	// ID jest kluczem sztucznym bazy; na zewnątrz idzie Identyfikator.
	ID int64
	// Identyfikator odpowiada kolumnie `identyfikator_zewnetrzny` — tożsamość
	// wpisu widoczna poza bazą.
	Identyfikator string
	// OknoId bywa pusty: transkrypcję wołają także ścieżki spoza okna
	// komunikacji. Pustka jest tu prawdą o zdarzeniu, nie brakiem danych, więc
	// w bazie odpowiada jej NULL.
	OknoId string
	// NagranieOdnosnik to ścieżka pliku na dysku Operatora, nie bajty dźwięku
	// — rdzeń otwiera plik w miejscu i niczego nie kopiuje.
	NagranieOdnosnik string
	// Model i Jezyk zapisane tak, jak obowiązywały w chwili zlecenia, a nie
	// odczytane z ustawień przy odczycie dziennika.
	Model string
	Jezyk string
	// Znakow to długość rozpoznanego tekstu. Samego tekstu dziennik nie trzyma:
	// idzie on do wołającego, a jego kopia obok byłaby drugą prawdą.
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

// FiltrTranskrypcji zawęża wykaz dziennika.
type FiltrTranskrypcji struct {
	// OknoId pusty znaczy „wszystkie okna, także wpisy spoza okna".
	OknoId string
	// Limit niedodatni znaczy LimitWykazuDomyslny — brak wskazania jest
	// wskazaniem wartości domyślnej, a nie prośbą o całość dziennika.
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
		` WHERE t.identyfikator_zewnetrzny = ?`

	// Jedno zapytanie na oba warianty żądania: puste zawężenie okna wyłącza
	// pierwszy warunek. Ze wskazanym oknem porządek biegnie indeksem
	// idx_transkrypcja_wykaz (okno_id, utworzono, id) czytanym wstecz, więc
	// „najnowsze najpierw" nie kosztuje sortowania wyniku (EXPLAIN QUERY PLAN:
	// SEARCH ... USING COVERING INDEX, bez kroku ORDER BY). Bez wskazania okna
	// sortowanie zostaje — i tak ma być: to zapytanie diagnostyczne przez cały
	// dziennik, a nie widok otwierany przy każdym zleceniu. Kolumna `id`
	// w porządku rozstrzyga wpisy z tej samej milisekundy, żeby kolejność była
	// stała między wywołaniami.
	listaTranskrypcji = `SELECT ` + kolumnyTranskrypcji + zrodloTranskrypcji +
		` WHERE (? = '' OR t.okno_id = ?)
		  ORDER BY t.utworzono DESC, t.id DESC
		  LIMIT ?`

	wstawTranskrypcje = `INSERT INTO transkrypcja
	                     (identyfikator_zewnetrzny, okno_id, nagranie_odnosnik, model,
	                      jezyk, znakow, trwanie_ms, stan, powod, utworzono)
	                     VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
)

// dziennikTranskrypcji jest jedyną implementacją Dziennika.
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

// Zapisz wstawia wpis dziennika i oddaje go odczytany z bazy.
func (d *dziennikTranskrypcji) Zapisz(ctx context.Context,
	wpis WpisTranskrypcji) (WpisTranskrypcji, error) {

	if err := sprawdzWpis(wpis); err != nil {
		return WpisTranskrypcji{}, err
	}
	if _, err := d.db.ExecContext(ctx, wstawTranskrypcje, wpis.Identyfikator,
		napisDoKolumny(wpis.OknoId), wpis.NagranieOdnosnik, wpis.Model, wpis.Jezyk,
		wpis.Znakow, wpis.TrwanieMs, wpis.Stan, napisDoKolumny(wpis.Powod),
		wpis.Utworzono); err != nil {

		return WpisTranskrypcji{}, fmt.Errorf("mowa: nie można zapisać transkrypcji %q: %w",
			wpis.Identyfikator, err)
	}
	return d.wpisPoIdentyfikatorze(ctx, wpis.Identyfikator)
}

// Wykaz oddaje wpisy dziennika od najnowszego.
func (d *dziennikTranskrypcji) Wykaz(ctx context.Context,
	filtr FiltrTranskrypcji) ([]WpisTranskrypcji, error) {

	limit := filtr.Limit
	if limit <= 0 {
		limit = LimitWykazuDomyslny
	}
	wiersze, err := d.db.QueryContext(ctx, listaTranskrypcji, filtr.OknoId, filtr.OknoId, limit)
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

	wpis, err := odczytajTranskrypcje(d.db.QueryRowContext(ctx, pobierzTranskrypcje, identyfikator))
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

// sprawdzWpis odrzuca wpisy, których baza i tak by nie przyjęła — po to, żeby
// wołający dostał zdanie po polsku zamiast komunikatu o naruszonym CHECK-u.
// Zdublowaniem więzu to nie jest: baza pozostaje ostatecznym strażnikiem,
// bo pisać do niej może też przyszła ścieżka, która tej funkcji nie wywoła.
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

// odczytajTranskrypcje składa wpis z jednego wiersza wyniku.
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

// napisDoKolumny przekłada napis pusty na NULL. Kolumny `okno_id` i `powod`
// dopuszczają pustkę, a pustka ta coś znaczy — „zlecenie spoza okna" oraz „nie
// było odmowy". Napis pusty zapisany wprost udawałby wartość podaną, a przy
// `powod` naruszałby CHECK tabeli `transkrypcja`.
func napisDoKolumny(wartosc string) any {
	if wartosc == "" {
		return nil
	}
	return wartosc
}
