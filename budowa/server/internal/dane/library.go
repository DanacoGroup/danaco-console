// Odpowiedzialność pliku: obszar plików modułu Library — definicja pliku
// repozytorium wiedzy (tabela `plik_biblioteki`) wraz z
// kontraktem całego obszaru. Wersje leżą w `library_wersje.go`, kolekcje
// i etykiety w `library_kolekcje.go` — jedno repozytorium, trzy pliki wedle
// odpowiedzialności, tak jak `dane/automations*.go`.
//
// Interfejs deklaruje wyłącznie ten plik, w całości — wraz z metodami, które
// implementują pozostałe pliki obszaru. Interfejs rozdzielony na trzy pliki
// byłby trzema prawdami o jednym kontrakcie.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// PlikBiblioteki to wiersz tabeli `plik_biblioteki`. Kod jest identyfikatorem,
// którym plik wychodzi kontraktem (`LibraryFile.id`).
type PlikBiblioteki struct {
	ID    int64
	Kod   string
	Nazwa string
	// Sciezka to ścieżka źródłowa podana przy wgraniu w `sourcePath` — zapis
	// tego, skąd plik przyszedł na maszynie Operatora. Bajty leżą gdzie indziej:
	// wciąga je magazyn treści rdzenia, a wskazuje `TrescOdwolanie`
	// (`adapter_modul_library_tresc.go`). Ta kolumna jest więc prowenancją, nie
	// drogą do treści.
	//
	// Nie wychodzi z rdzenia. Pole `path` kontraktu — `LibraryFile.path`,
	// „Sciezka w strukturze repozytorium" — to droga wewnątrz biblioteki, nie
	// ścieżka systemowa: wypełnienie go tą kolumną wynosiłoby do klienta i do
	// modelu katalog roboczy Operatora, a klient czyta pierwszy człon po `/`
	// jako katalog nawigacji i zawęża wykaz `startsWith`
	// (`client/src/moduly/library/wykaz-plikow.ts`, `stan-biblioteki.ts`).
	// Porządek biblioteki niosą kolekcje, powiązanie wiele-do-wielu, z którego
	// jednej ścieżki nie da się wyprowadzić. Pole kontraktu zostaje więc puste;
	// nadanie mu treści wymaga pojęcia ścieżki w kontrakcie — zgłoszone, nie
	// zrobione.
	Sciezka *string
	// SciezkaRepozytorium to droga WEWNĄTRZ biblioteki — pole `LibraryFile.path`
	// kontraktu, którym rozporządza `library.file.move`. Rozłączna z `Sciezka`:
	// tamta mówi, skąd plik przyszedł na maszynie Operatora, ta — gdzie leży
	// w porządku repozytorium (migracja 180).
	SciezkaRepozytorium *string
	// Stan rozdziela wykaz czynny od archiwum: `aktywny` albo `zarchiwizowany`
	// (kontrakt: `LibraryFileStatus`). Archiwizacja jest przeniesieniem między
	// stanami, nie usunięciem wiersza.
	Stan            string
	MimeType        *string
	RozmiarBajtow   *int64
	ProjektID       *string
	ModulZrodlowyID *string
	WersjaBiezacaID *int64
	SumaKontrolna   *string
	TrescOdwolanie  *string
	Utworzono       string
	Zaktualizowano  string
}

// FiltrPlikow niesie zawężenia wspólne dla `library.file.list` i
// `library.file.search`. Fraza filtru działa na nazwie pliku — `library.file.list`
// zawęża wykaz, a nie szuka w dokumentach. Przeszukanie treści należy do
// `Szukaj`, które dokłada do tego samego zawężenia klauzulę indeksu treści
// (`biblioteka_indeks_tresci.go`).
type FiltrPlikow struct {
	Fraza       *string
	Etykiety    []string
	KolekcjaKod *string
	ProjektID   *string
	// Stan zawęża do zasobów czynnych albo archiwalnych. Puste znaczy „bez
	// zawężenia" — wykaz domyślny modułu podaje tu `aktywny`, żeby kosz
	// repozytorium naprawdę zdejmował zasób z widoku, a czynności higieny
	// obchodziły całość.
	Stan   *string
	Limit  int
	Offset int
}

// RepozytoriumBiblioteki jest kontraktem obszaru Library.
type RepozytoriumBiblioteki interface {
	// --- agent A: plik biblioteki ---
	ZapiszPlik(ctx context.Context, plik PlikBiblioteki) (PlikBiblioteki, error)
	Plik(ctx context.Context, kod string) (PlikBiblioteki, error)
	Pliki(ctx context.Context, filtr FiltrPlikow) ([]PlikBiblioteki, int, error)
	Szukaj(ctx context.Context, fraza string, filtr FiltrPlikow) ([]PlikBiblioteki, int, error)
	// ZapiszIndeksTresci zasila indeks treści (FTS5) wyciągiem
	// tekstowym pliku — bez niego `Szukaj` widzi wyłącznie nazwy
	// (`biblioteka_indeks_tresci.go`).
	ZapiszIndeksTresci(ctx context.Context, plikID int64, wyciag string) error
	// OdwolaniaTresci wymienia bloby trzymane jeszcze przy życiu przez plik
	// albo wersję — dla sprzątania magazynu (`biblioteka_odwolania.go`).
	OdwolaniaTresci(ctx context.Context) ([]string, error)

	// --- agent B: wersje ---
	ZapiszWersje(ctx context.Context, plikID int64, wersja WersjaPlikuBiblioteki) (WersjaPlikuBiblioteki, error)
	// DolozWersje dokłada wersję ORAZ przestawia na nią plik macierzysty w
	// jednej transakcji (`library_wersje_zapis.go`) — inaczej niż `ZapiszWersje`,
	// które samo wskaźnika bieżącej wersji nie rusza.
	DolozWersje(ctx context.Context, plikID int64, wersja WersjaPlikuBiblioteki) (WersjaPlikuBiblioteki, PlikBiblioteki, error)
	Wersje(ctx context.Context, plikID int64) ([]WersjaPlikuBiblioteki, error)
	PrzywrocWersje(ctx context.Context, plikID int64, kodWersji string) (PlikBiblioteki, error)

	// --- agent C: kolekcje i etykiety ---
	UtworzKolekcje(ctx context.Context, kolekcja KolekcjaBiblioteki) (KolekcjaBiblioteki, error)
	PrzypiszDoKolekcji(ctx context.Context, kodKolekcji string, kodyPlikow []string) ([]string, error)
	Kolekcje(ctx context.Context) ([]KolekcjaBiblioteki, error)
	UstawEtykiety(ctx context.Context, kodPliku string, etykiety []string) ([]string, error)
	Etykiety(ctx context.Context, plikID int64) ([]string, error)
	// KolekcjePliku odpowiada na pytanie „do których kolekcji należy ten plik"
	// (`LibraryFile.collectionIds`), a UstawKolekcjePliku czyni wykaz kolekcji
	// pliku dokładnie takim, jaki podano — wraz ze zdjęciem z kolekcji spoza
	// wykazu, czego `PrzypiszDoKolekcji` nie umie (`biblioteka_kolekcje_pliku.go`).
	KolekcjePliku(ctx context.Context, plikID int64) ([]string, error)
	UstawKolekcjePliku(ctx context.Context, kodPliku string, kodyKolekcji []string) ([]string, error)

	// --- opis zasobu i schemat metadanych (`biblioteka_opis.go`) ---
	Opis(ctx context.Context, plikID int64) (OpisZasobuBiblioteki, error)
	ZapiszOpis(ctx context.Context, plikID int64, opis OpisZasobuBiblioteki) error
	PolaSchematu(ctx context.Context, mimeType, kolekcjaKod *string) ([]PoleSchematuBiblioteki, error)
	ZapiszPoleSchematu(ctx context.Context, pole PoleSchematuBiblioteki) (PoleSchematuBiblioteki, error)
	UsunPoleSchematu(ctx context.Context, kod string) (bool, error)
	// ZasobyZPolem liczy zasoby, przy których pole niestandardowe ma już
	// wartość — liczba wchodzi do odpowiedzi `library.schema.set`, żeby zdjęcie
	// definicji nie wyglądało na czynność bez skutku ubocznego.
	ZasobyZPolem(ctx context.Context, kodPola string) (int, error)

	// --- słownik etykiet i tezaurus (`biblioteka_slownik.go`) ---
	EtykietySlownika(ctx context.Context, fraza *string, tylkoNieuzywane bool,
		limit int) ([]EtykietaSlownikaBiblioteki, int, error)
	ZapiszEtykieteSlownika(ctx context.Context, nazwa string, barwa *string) (EtykietaSlownikaBiblioteki, error)
	PrzemianujEtykiete(ctx context.Context, stara, nowa string) (int, error)
	UsunEtykieteZeSlownika(ctx context.Context, nazwa string) (int, error)
	EtykietaSlownika(ctx context.Context, nazwa string) (EtykietaSlownikaBiblioteki, error)
	UstawRelacjeTezaurusa(ctx context.Context, zrodlo, cel, rodzaj string, zdejmij bool) (bool, error)
	RelacjeTezaurusa(ctx context.Context) ([]RelacjaTezaurusaBiblioteki, error)

	// --- kolekcje w postaci pełnej (`biblioteka_reguly.go`) ---
	KolekcjeWykaz(ctx context.Context, rodzicKod, fraza *string, limit int) ([]KolekcjaBiblioteki, int, error)

	// --- reguły repozytorium (`biblioteka_reguly.go`) ---
	ZapiszRegule(ctx context.Context, regula RegulaBiblioteki) (RegulaBiblioteki, error)
	Regula(ctx context.Context, kod string) (RegulaBiblioteki, error)
	Reguly(ctx context.Context, rodzaj *string, tylkoCzynne bool) ([]RegulaBiblioteki, error)
	UsunRegule(ctx context.Context, kod string) (bool, error)
	// PrzypiszRegula wpisuje zasoby do kolekcji ze znacznikiem pochodzenia
	// `regula`, a OdepnijPrzypisaniaReguly zdejmuje wyłącznie te wpisy —
	// przypisanie ręczne Operatora zostaje.
	PrzypiszRegula(ctx context.Context, kodKolekcji string, kodyPlikow []string) (int, error)
	OdepnijPrzypisaniaReguly(ctx context.Context, kodKolekcji string) (int, error)

	// --- dziennik audytu (`biblioteka_audyt.go`) ---
	ZapiszWpisAudytu(ctx context.Context, wpis WpisAudytuBiblioteki) (WpisAudytuBiblioteki, error)
	WpisyAudytu(ctx context.Context, filtr FiltrAudytuBiblioteki) ([]WpisAudytuBiblioteki, int, error)

	// --- retencja i utrwalenie (`biblioteka_archiwum.go`) ---
	ZapiszPolitykeRetencji(ctx context.Context, polityka PolitykaRetencjiBiblioteki) (PolitykaRetencjiBiblioteki, error)
	PolitykiRetencji(ctx context.Context, zasieg *string) ([]PolitykaRetencjiBiblioteki, error)
	UsunPolitykeRetencji(ctx context.Context, kod string) (bool, error)
	ZapiszUtrwalenie(ctx context.Context, zadanie ZadanieUtrwaleniaBiblioteki) (ZadanieUtrwaleniaBiblioteki, error)

	// --- udostępnienia i nasłuchy (`biblioteka_udostepnienia.go`) ---
	ZapiszUdostepnienie(ctx context.Context, udostepnienie UdostepnienieBiblioteki) (UdostepnienieBiblioteki, error)
	Udostepnienia(ctx context.Context, celKod *string, tylkoCzynne bool) ([]UdostepnienieBiblioteki, error)
	OdwolajUdostepnienie(ctx context.Context, kod string) (bool, error)
	ZapiszWebhook(ctx context.Context, webhook WebhookBiblioteki) (WebhookBiblioteki, error)
	Webhooki(ctx context.Context, tylkoCzynne bool) ([]WebhookBiblioteki, error)
	UsunWebhook(ctx context.Context, kod string) (bool, error)

	// --- sugestie porządkujące (`biblioteka_sugestie.go`) ---
	ZapiszSugestie(ctx context.Context, sugestia SugestiaBiblioteki) (SugestiaBiblioteki, error)
	Sugestie(ctx context.Context, plikKod *string, rodzaje []string, limit int) ([]SugestiaBiblioteki, int, error)
	Sugestia(ctx context.Context, kod string) (SugestiaBiblioteki, error)
	RozstrzygnijSugestie(ctx context.Context, kody []string, przyjeto bool) (int, error)

	// --- cykl życia zasobu i pulpit stanu (`biblioteka_cykl.go`) ---
	UstawStanPlikow(ctx context.Context, kody []string, stan string) ([]PlikBiblioteki, error)
	UstawSciezkeRepozytorium(ctx context.Context, kody []string, sciezka string) ([]PlikBiblioteki, error)
	PrzemianujPlik(ctx context.Context, kod, nazwa string) (PlikBiblioteki, error)
	UsunPliki(ctx context.Context, kody []string) (int, int, error)
	Statystyki(ctx context.Context, kolekcjaKod, projektID *string, topN int) (StatystykiBiblioteki, error)
}

const (
	kolumnyPlikuBiblioteki = `id, identyfikator_zewnetrzny, nazwa, sciezka, sciezka_repozytorium,
	                          stan, mime_type,
	                          rozmiar_bajtow, projekt_id, modul_zrodlowy_id, wersja_biezaca_id,
	                          suma_kontrolna, tresc_odwolanie, utworzono, zaktualizowano`

	// Stan wchodzi do zapisu z wartością domyślną nadaną przez wołającego:
	// pusty łańcuch nie przejdzie warunku CHECK kolumny, więc `ZapiszPlik`
	// podstawia `aktywny` (`stanZapisu`) zamiast pozwolić bazie odmówić przy
	// każdym wgraniu, które o stanie nie myśli.
	zapiszPlikBiblioteki = `INSERT INTO plik_biblioteki
	                        (identyfikator_zewnetrzny, nazwa, sciezka, sciezka_repozytorium, stan,
	                         mime_type, rozmiar_bajtow,
	                         projekt_id, modul_zrodlowy_id, wersja_biezaca_id, suma_kontrolna,
	                         tresc_odwolanie)
	                        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	                        ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                            nazwa = excluded.nazwa,
	                            sciezka = excluded.sciezka,
	                            sciezka_repozytorium = excluded.sciezka_repozytorium,
	                            stan = excluded.stan,
	                            mime_type = excluded.mime_type,
	                            rozmiar_bajtow = excluded.rozmiar_bajtow,
	                            projekt_id = excluded.projekt_id,
	                            modul_zrodlowy_id = excluded.modul_zrodlowy_id,
	                            wersja_biezaca_id = excluded.wersja_biezaca_id,
	                            suma_kontrolna = excluded.suma_kontrolna,
	                            tresc_odwolanie = excluded.tresc_odwolanie,
	                            zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')`

	pobierzPlikBiblioteki = `SELECT ` + kolumnyPlikuBiblioteki + ` FROM plik_biblioteki
	                         WHERE identyfikator_zewnetrzny = ?`
)

type repozytoriumBiblioteki struct {
	zapytania *zapytania
	db        *sql.DB
}

func noweRepozytoriumBiblioteki(z *zapytania, db *sql.DB) *repozytoriumBiblioteki {
	return &repozytoriumBiblioteki{zapytania: z, db: db}
}

// ZapiszPlik zakłada wiersz pliku albo nadpisuje zastany i zwraca stan po
// zapisie. Etykiety i przypisania do kolekcji nie są tu ruszane — to obszar
// agenta C (`UstawEtykiety`, `PrzypiszDoKolekcji`), wołany osobno przez adapter.
func (r *repozytoriumBiblioteki) ZapiszPlik(ctx context.Context, plik PlikBiblioteki) (PlikBiblioteki, error) {
	if plik.Kod == "" {
		return PlikBiblioteki{}, fmt.Errorf("dane: plik biblioteki bez identyfikatora")
	}
	if plik.Nazwa == "" {
		return PlikBiblioteki{}, fmt.Errorf("dane: plik biblioteki %q bez nazwy", plik.Kod)
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszPlikBiblioteki)
	if err != nil {
		return PlikBiblioteki{}, err
	}
	_, err = polecenie.ExecContext(ctx, plik.Kod, plik.Nazwa, tekstDoKolumny(plik.Sciezka),
		tekstDoKolumny(plik.SciezkaRepozytorium), stanZapisu(plik.Stan),
		tekstDoKolumny(plik.MimeType), liczbaDoKolumny(plik.RozmiarBajtow),
		tekstDoKolumny(plik.ProjektID), tekstDoKolumny(plik.ModulZrodlowyID),
		liczbaDoKolumny(plik.WersjaBiezacaID), tekstDoKolumny(plik.SumaKontrolna),
		tekstDoKolumny(plik.TrescOdwolanie))
	if err != nil {
		return PlikBiblioteki{}, fmt.Errorf("dane: nie można zapisać pliku biblioteki %q: %w", plik.Kod, err)
	}
	return r.Plik(ctx, plik.Kod)
}

// Plik zwraca plik o wskazanym kodzie. Brak wiersza wraca jako ErrBrakWiersza —
// warstwa wyższa odróżnia „nie ma” od „odczyt się nie powiódł”.
func (r *repozytoriumBiblioteki) Plik(ctx context.Context, kod string) (PlikBiblioteki, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzPlikBiblioteki)
	if err != nil {
		return PlikBiblioteki{}, err
	}
	plik, err := odczytajPlikBiblioteki(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return PlikBiblioteki{}, ErrBrakWiersza
	}
	if err != nil {
		return PlikBiblioteki{}, fmt.Errorf("dane: nieczytelny wiersz pliku biblioteki %q: %w", kod, err)
	}
	return plik, nil
}

// Pliki zwraca pliki spełniające filtr, posortowane od najnowszych, wraz
// z całkowitą liczbą trafień sprzed obcięcia limitem (`LibraryFileListResponse.total`).
func (r *repozytoriumBiblioteki) Pliki(ctx context.Context, filtr FiltrPlikow) ([]PlikBiblioteki, int, error) {
	warunek, argumenty := warunkiFiltruPlikow(filtr)
	return r.pliki(ctx, warunek, argumenty, filtr.Limit, filtr.Offset)
}

// Szukaj realizuje `library.file.search` — fraza trafia w nazwę albo w treść.
//
// Dwa dopasowania, jedna fraza. Żądanie kontraktu niesie jedną frazę bez
// rozróżnienia „po nazwie" / „po treści", więc obie drogi muszą być pytane
// naraz, a trafienie którąkolwiek jest trafieniem. Nazwę dopasowuje `LIKE`
// (fragment dowolnego miejsca nazwy, także pół słowa), treść — indeks FTS5
// (`warunekTresciBiblioteki`), zasilany przy wgraniu pliku i przy każdej nowej
// wersji.
//
// Plik niezaindeksowany nie znika z wyszukiwania. Wiersze niezaindeksowane
// oraz treści nietekstowe nie mają wiersza w indeksie — dopasowanie nazwy
// działa dla nich dalej, bo jest osobnym członem alternatywy, a nie warunkiem
// dodatkowym. Zawężenia pozostałe (projekt, kolekcja, etykiety) obowiązują
// obie drogi tak samo.
func (r *repozytoriumBiblioteki) Szukaj(ctx context.Context, fraza string,
	filtr FiltrPlikow) ([]PlikBiblioteki, int, error) {

	// Fraza nie idzie przez `FiltrPlikow.Fraza`, bo tamten człon zawęża do
	// nazwy; tu nazwa i treść stoją w alternatywie i muszą być jedną klauzulą.
	filtr.Fraza = nil
	warunek, argumenty := warunkiFiltruPlikow(filtr)
	warunek += ` AND (nazwa LIKE ? OR ` + warunekTresciBiblioteki + `)`
	argumenty = append(argumenty, "%"+fraza+"%", zapytanieTresci(fraza))
	return r.pliki(ctx, warunek, argumenty, filtr.Limit, filtr.Offset)
}

// warunkiFiltruPlikow składa klauzulę WHERE i argumenty wspólne dla `Pliki`
// i `Szukaj` — oba budują to samo zawężenie, różni je tylko obowiązkowość frazy.
func warunkiFiltruPlikow(filtr FiltrPlikow) (string, []any) {
	warunki := []string{"1 = 1"}
	argumenty := []any{}

	if filtr.Fraza != nil && *filtr.Fraza != "" {
		warunki = append(warunki, "nazwa LIKE ?")
		argumenty = append(argumenty, "%"+*filtr.Fraza+"%")
	}
	if filtr.Stan != nil && *filtr.Stan != "" {
		warunki = append(warunki, "stan = ?")
		argumenty = append(argumenty, *filtr.Stan)
	}
	if filtr.ProjektID != nil && *filtr.ProjektID != "" {
		warunki = append(warunki, "projekt_id = ?")
		argumenty = append(argumenty, *filtr.ProjektID)
	}
	if filtr.KolekcjaKod != nil && *filtr.KolekcjaKod != "" {
		warunki = append(warunki, `EXISTS (SELECT 1 FROM przypisanie_kolekcji_biblioteki pk
		                          JOIN kolekcja_biblioteki k ON k.id = pk.kolekcja_id
		                          WHERE pk.plik_id = plik_biblioteki.id
		                            AND k.identyfikator_zewnetrzny = ?)`)
		argumenty = append(argumenty, *filtr.KolekcjaKod)
	}
	// Plik musi nieść każdą podaną etykietę (zawężenie, nie dopasowanie
	// dowolnej) — jedna klauzula EXISTS na etykietę.
	for _, etykieta := range filtr.Etykiety {
		if etykieta == "" {
			continue
		}
		warunki = append(warunki, `EXISTS (SELECT 1 FROM etykieta_pliku_biblioteki e
		                          WHERE e.plik_id = plik_biblioteki.id AND e.etykieta = ?)`)
		argumenty = append(argumenty, etykieta)
	}
	return strings.Join(warunki, " AND "), argumenty
}

// pliki wykonuje zapytanie o listę i zapytanie o liczbę spełniających
// warunek — jedno miejsce dla `Pliki` i `Szukaj`, bo różnią się tylko
// klauzulą WHERE, nie sposobem odczytu.
func (r *repozytoriumBiblioteki) pliki(ctx context.Context, warunek string, argumenty []any,
	limit, offset int) ([]PlikBiblioteki, int, error) {

	zapytanieListy := `SELECT ` + kolumnyPlikuBiblioteki + ` FROM plik_biblioteki
	                    WHERE ` + warunek + `
	                    ORDER BY zaktualizowano DESC, id DESC
	                    LIMIT ? OFFSET ?`
	wiersze, err := r.db.QueryContext(ctx, zapytanieListy,
		append(append([]any{}, argumenty...), granicaWykazu(limit), przesuniecieWykazu(offset))...)
	if err != nil {
		return nil, 0, fmt.Errorf("dane: nie można odczytać plików biblioteki: %w", err)
	}
	defer wiersze.Close()

	lista := []PlikBiblioteki{}
	for wiersze.Next() {
		plik, err := odczytajPlikBiblioteki(wiersze)
		if err != nil {
			return nil, 0, fmt.Errorf("dane: nieczytelny wiersz pliku biblioteki: %w", err)
		}
		lista = append(lista, plik)
	}
	if err := wiersze.Err(); err != nil {
		return nil, 0, fmt.Errorf("dane: przerwany odczyt plików biblioteki: %w", err)
	}

	zapytanieLiczby := `SELECT COUNT(*) FROM plik_biblioteki WHERE ` + warunek
	var lacznie int
	if err := r.db.QueryRowContext(ctx, zapytanieLiczby, argumenty...).Scan(&lacznie); err != nil {
		return nil, 0, fmt.Errorf("dane: nie można policzyć plików biblioteki: %w", err)
	}
	return lista, lacznie, nil
}

// przesuniecieWykazu przekłada przesunięcie żądania na argument zapytania —
// wartość ujemna nie ma znaczenia w SQL, więc zjeżdża do zera.
func przesuniecieWykazu(offset int) int {
	if offset <= 0 {
		return 0
	}
	return offset
}

// stanZapisu podstawia stan domyślny za wartość niepodaną. Kolumna ma warunek
// CHECK, więc pusty łańcuch byłby odmową bazy przy każdym wgraniu — a wgranie
// o stanie nie rozstrzyga: zasób wchodzi do wykazu czynnego.
func stanZapisu(stan string) string {
	if stan == "" {
		return StanZasobuCzynny
	}
	return stan
}

// Stany zasobu repozytorium — odwzorowanie `LibraryFileStatus` kontraktu na
// wartości kolumny `plik_biblioteki.stan`.
const (
	StanZasobuCzynny         = "aktywny"
	StanZasobuZarchiwizowany = "zarchiwizowany"
)

// odczytajPlikBiblioteki składa strukturę z jednego wiersza wyniku.
func odczytajPlikBiblioteki(wiersz skaner) (PlikBiblioteki, error) {
	var plik PlikBiblioteki
	var sciezka, sciezkaRepozytorium, mimeType sql.NullString
	var projektID, modulZrodlowyID, sumaKontrolna, trescOdwolanie sql.NullString
	var rozmiarBajtow, wersjaBiezacaID sql.NullInt64
	err := wiersz.Scan(&plik.ID, &plik.Kod, &plik.Nazwa, &sciezka, &sciezkaRepozytorium,
		&plik.Stan, &mimeType, &rozmiarBajtow,
		&projektID, &modulZrodlowyID, &wersjaBiezacaID, &sumaKontrolna, &trescOdwolanie,
		&plik.Utworzono, &plik.Zaktualizowano)
	if err != nil {
		return PlikBiblioteki{}, err
	}
	plik.Sciezka = tekstZKolumny(sciezka)
	plik.SciezkaRepozytorium = tekstZKolumny(sciezkaRepozytorium)
	plik.MimeType = tekstZKolumny(mimeType)
	plik.RozmiarBajtow = liczbaZKolumny(rozmiarBajtow)
	plik.ProjektID = tekstZKolumny(projektID)
	plik.ModulZrodlowyID = tekstZKolumny(modulZrodlowyID)
	plik.WersjaBiezacaID = liczbaZKolumny(wersjaBiezacaID)
	plik.SumaKontrolna = tekstZKolumny(sumaKontrolna)
	plik.TrescOdwolanie = tekstZKolumny(trescOdwolanie)
	return plik, nil
}
