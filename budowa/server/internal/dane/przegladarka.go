// Plik niesie kontrakt obszaru Browser oraz jedyną metodę implementowaną tu
// w całości: migawkę strony. Migawka jest historią nawigacji, nie osobnym
// bytem — kolejne wiersze uporządkowane po utworzeniu są tą historią.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// MigawkaStrony to wiersz tabeli `migawka_strony` — zrzut stanu strony w danej
// chwili, zwracany przez `browser.navigate` i `browser.snapshot.get` pod
// wspólnym kształtem `BrowserSnapshot`.
type MigawkaStrony struct {
	ID              int64
	Kod             string
	Okno            string
	Url             string
	Tytul           *string
	TekstOdwolanie  *string
	ZrodloOdwolanie *string
	ZrzutOdwolanie  *string
	Utworzono       string
}

// RepozytoriumPrzegladania jest kontraktem obszaru Browser: migawki, źródła,
// notatki, karty, monitory, kanały, zakładki, wytwory i zestawy przeglądania.
type RepozytoriumPrzegladania interface {
	// --- migawki ---
	ZapiszMigawke(ctx context.Context, migawka MigawkaStrony) (MigawkaStrony, error)
	Migawka(ctx context.Context, kod string) (MigawkaStrony, error)
	OstatniaMigawka(ctx context.Context, okno string) (MigawkaStrony, error)
	Historia(ctx context.Context, okno string, limit int) ([]MigawkaStrony, error)

	// --- źródła ---
	ZapiszZrodlo(ctx context.Context, zrodlo ZrodloPrzegladania) (ZrodloPrzegladania, error)
	Zrodla(ctx context.Context, filtr FiltrZrodelPrzegladania) ([]ZrodloPrzegladania, error)

	// --- notatki ---
	ZapiszNotatke(ctx context.Context, notatka NotatkaPrzegladania) (NotatkaPrzegladania, error)
	Notatki(ctx context.Context, filtr FiltrNotatekPrzegladania) ([]NotatkaPrzegladania, error)

	// --- notatki: zmiana i odczyt pojedynczej ---
	AktualizujNotatke(ctx context.Context, zmiana ZmianaNotatki) (NotatkaPrzegladania, error)
	Notatka(ctx context.Context, kod string) (NotatkaPrzegladania, error)

	// --- karty, grupy kart i przestrzenie robocze (migracja 170) ---
	ZapiszKarte(ctx context.Context, karta KartaPrzegladania) (KartaPrzegladania, error)
	Karta(ctx context.Context, kod string) (KartaPrzegladania, error)
	Karty(ctx context.Context, filtr FiltrKartPrzegladania) ([]KartaPrzegladania, error)
	NastepnaKolejnoscKarty(ctx context.Context, okno string) (int64, error)
	ZamknijKarte(ctx context.Context, kod string) (bool, error)
	OdznaczPozostaleKarty(ctx context.Context, okno, kod string) error
	PrzypiszKarteDoGrupy(ctx context.Context, okno, karta, grupa string) error
	PrzypiszKarteDoPrzestrzeni(ctx context.Context, okno, karta, przestrzen string) error
	ZapiszGrupeKart(ctx context.Context, grupa GrupaKart) (GrupaKart, error)
	GrupaKart(ctx context.Context, kod string) (GrupaKart, error)
	GrupyKart(ctx context.Context, okno string, limit int) ([]GrupaKart, error)
	UsunGrupeKart(ctx context.Context, kod string) (bool, error)
	ZapiszPrzestrzen(ctx context.Context, przestrzen PrzestrzenPrzegladania) (PrzestrzenPrzegladania, error)
	Przestrzen(ctx context.Context, kod string) (PrzestrzenPrzegladania, error)
	Przestrzenie(ctx context.Context, okno string, limit int) ([]PrzestrzenPrzegladania, error)
	UsunPrzestrzen(ctx context.Context, kod string) (bool, error)
	LiczbaKartPrzestrzeni(ctx context.Context, kod string) (int64, error)

	// --- monitory, kanały, kolejka czytania i zakładki (migracje 171-174) ---
	ZapiszMonitor(ctx context.Context, monitor MonitorPrzegladania) (MonitorPrzegladania, error)
	Monitor(ctx context.Context, kod string) (MonitorPrzegladania, error)
	Monitory(ctx context.Context, okno string, tylkoWlaczone bool, limit int) ([]MonitorPrzegladania, error)
	UsunMonitor(ctx context.Context, kod string) (bool, error)
	ZapiszKanal(ctx context.Context, kanal KanalPrzegladania) (KanalPrzegladania, error)
	Kanal(ctx context.Context, kod string) (KanalPrzegladania, error)
	KanalPoAdresie(ctx context.Context, okno, url string) (KanalPrzegladania, error)
	Kanaly(ctx context.Context, okno, kod string, limit int) ([]KanalPrzegladania, error)
	UsunKanal(ctx context.Context, kod string) (bool, error)
	ZapiszWpisKanalu(ctx context.Context, wpis WpisKanalu) error
	WpisyKanalu(ctx context.Context, kanal string, tylkoNieprzeczytane bool, limit int) ([]WpisKanalu, error)
	NieprzeczytaneKanalu(ctx context.Context, kanal string) (int64, error)
	ZapiszPozycjeCzytania(ctx context.Context, pozycja PozycjaCzytania) (PozycjaCzytania, error)
	PozycjaCzytania(ctx context.Context, kod string) (PozycjaCzytania, error)
	KolejkaCzytania(ctx context.Context, okno string, tylkoNieprzeczytane bool, limit int) ([]PozycjaCzytania, error)
	UsunPozycjeCzytania(ctx context.Context, kod string) (bool, error)
	ZapiszZakladke(ctx context.Context, zakladka ZakladkaPrzegladania) (ZakladkaPrzegladania, error)
	Zakladka(ctx context.Context, kod string) (ZakladkaPrzegladania, error)
	Zakladki(ctx context.Context, filtr FiltrZakladek) ([]ZakladkaPrzegladania, error)
	UsunZakladke(ctx context.Context, kod string) (bool, error)

	// --- wytwory, zrzuty, pobrania, makra i granice (migracje 175-178) ---
	ZapiszWytwor(ctx context.Context, wytwor WytworPrzegladania) (WytworPrzegladania, error)
	Wytwor(ctx context.Context, kod string) (WytworPrzegladania, error)
	Wytwory(ctx context.Context, okno, rodzaj string, limit int) ([]WytworPrzegladania, error)
	ZapiszZrzut(ctx context.Context, zrzut ZrzutPrzegladania) (ZrzutPrzegladania, error)
	Zrzut(ctx context.Context, kod string) (ZrzutPrzegladania, error)
	ZrzutPoOdwolaniu(ctx context.Context, odwolanie string) (ZrzutPrzegladania, error)
	ZrzutMigawki(ctx context.Context, migawka string) (ZrzutPrzegladania, error)
	ZapiszPobranie(ctx context.Context, pobranie PobraniePrzegladania) (PobraniePrzegladania, error)
	Pobranie(ctx context.Context, kod string) (PobraniePrzegladania, error)
	Pobrania(ctx context.Context, okno, stan string, limit int) ([]PobraniePrzegladania, error)
	UsunPobranie(ctx context.Context, kod string) (bool, error)
	ZapiszMakro(ctx context.Context, makro MakroPrzegladania) (MakroPrzegladania, error)
	Makro(ctx context.Context, kod string) (MakroPrzegladania, error)
	Makra(ctx context.Context, okno string, limit int) ([]MakroPrzegladania, error)
	ZapiszGranice(ctx context.Context, granica GranicaWykonawcy) (GranicaWykonawcy, error)
	Granice(ctx context.Context, zasieg, zasiegID string) (GranicaWykonawcy, error)

	// --- zestawy źródeł i wątki notatek (migracja 179) ---
	ZapiszZestawZrodel(ctx context.Context, zestaw ZestawZrodel) (ZestawZrodel, error)
	ZestawZrodel(ctx context.Context, kod string) (ZestawZrodel, error)
	ZestawyZrodel(ctx context.Context, okno string, limit int) ([]ZestawZrodel, error)
	UsunZestawZrodel(ctx context.Context, kod string) (bool, error)
	PrzypiszZrodloDoZestawu(ctx context.Context, okno, zrodlo, zestaw string) error
	KodyZrodelZestawu(ctx context.Context, zestaw string) ([]string, error)
	UsunZrodlo(ctx context.Context, kod string) (bool, error)
	ZapiszWatekNotatek(ctx context.Context, watek WatekNotatek) (WatekNotatek, error)
	WatekNotatek(ctx context.Context, kod string) (WatekNotatek, error)
	WatkiNotatek(ctx context.Context, okno string, limit int) ([]WatekNotatek, error)
	UsunWatekNotatek(ctx context.Context, kod string) (bool, error)
	PrzypiszNotatkeDoWatku(ctx context.Context, okno, notatka, watek string) error
	KodyNotatekWatku(ctx context.Context, watek string) ([]string, error)

	// --- wspólne: rozróżnienie okna nieznanego od okna pustego ---
	OknoZnane(ctx context.Context, okno string) (bool, error)
}

const (
	kolumnyMigawki = `id, identyfikator_zewnetrzny, okno, url, tytul,
	                  tekst_odwolanie, zrodlo_odwolanie, zrzut_odwolanie, utworzono`

	// Migawka nigdy nie nadpisuje poprzedniej — każde wywołanie navigate/snapshot
	// jest nowym wierszem historii, więc wstawienie jest zwykłym INSERT-em, bez
	// ON CONFLICT (w przeciwieństwie do `automatyka`).
	zapiszMigawke = `INSERT INTO migawka_strony
	                 (identyfikator_zewnetrzny, okno, url, tytul,
	                  tekst_odwolanie, zrodlo_odwolanie, zrzut_odwolanie)
	                 VALUES (?, ?, ?, ?, ?, ?, ?)`

	pobierzMigawke = `SELECT ` + kolumnyMigawki + ` FROM migawka_strony
	                  WHERE identyfikator_zewnetrzny = ?`

	// Kolejność zgodna z indeksem idx_migawka_strony_okno (okno, utworzono DESC, id) —
	// najświeższy wiersz jest zarówno ostatnią migawką, jak i głową historii.
	ostatniaMigawkaOkna = `SELECT ` + kolumnyMigawki + ` FROM migawka_strony
	                       WHERE okno = ? ORDER BY utworzono DESC, id DESC LIMIT 1`

	historiaMigawekOkna = `SELECT ` + kolumnyMigawki + ` FROM migawka_strony
	                       WHERE okno = ? ORDER BY utworzono DESC, id DESC LIMIT ?`
)

type repozytoriumPrzegladania struct {
	zapytania *zapytania
	db        *sql.DB
}

func noweRepozytoriumPrzegladania(z *zapytania, db *sql.DB) *repozytoriumPrzegladania {
	return &repozytoriumPrzegladania{zapytania: z, db: db}
}

// ZapiszMigawke wstawia nowy zrzut stanu strony i zwraca go odczytany z bazy,
// z nadanym identyfikatorem i czasem utworzenia.
func (r *repozytoriumPrzegladania) ZapiszMigawke(ctx context.Context,
	migawka MigawkaStrony) (MigawkaStrony, error) {

	if migawka.Kod == "" {
		return MigawkaStrony{}, fmt.Errorf("dane: migawka strony bez identyfikatora")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszMigawke)
	if err != nil {
		return MigawkaStrony{}, err
	}
	_, err = polecenie.ExecContext(ctx, migawka.Kod, migawka.Okno, migawka.Url,
		tekstDoKolumny(migawka.Tytul), tekstDoKolumny(migawka.TekstOdwolanie),
		tekstDoKolumny(migawka.ZrodloOdwolanie), tekstDoKolumny(migawka.ZrzutOdwolanie))
	if err != nil {
		return MigawkaStrony{}, fmt.Errorf("dane: nie można zapisać migawki %q: %w", migawka.Kod, err)
	}
	return r.Migawka(ctx, migawka.Kod)
}

// Migawka zwraca zrzut o wskazanym kodzie. Brak wiersza wraca jako
// ErrBrakWiersza — warstwa wyższa odróżnia „nie ma" od „odczyt się nie powiódł".
func (r *repozytoriumPrzegladania) Migawka(ctx context.Context, kod string) (MigawkaStrony, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzMigawke)
	if err != nil {
		return MigawkaStrony{}, err
	}
	migawka, err := odczytajMigawke(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return MigawkaStrony{}, ErrBrakWiersza
	}
	if err != nil {
		return MigawkaStrony{}, fmt.Errorf("dane: nieczytelny wiersz migawki %q: %w", kod, err)
	}
	return migawka, nil
}

// OstatniaMigawka zwraca najświeższy zrzut okna — służy `browser.snapshot.get`
// wywołanemu bez podanego identyfikatora migawki. Brak wiersza wraca jako
// ErrBrakWiersza.
func (r *repozytoriumPrzegladania) OstatniaMigawka(ctx context.Context, okno string) (MigawkaStrony, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, ostatniaMigawkaOkna)
	if err != nil {
		return MigawkaStrony{}, err
	}
	migawka, err := odczytajMigawke(polecenie.QueryRowContext(ctx, okno))
	if errors.Is(err, sql.ErrNoRows) {
		return MigawkaStrony{}, ErrBrakWiersza
	}
	if err != nil {
		return MigawkaStrony{}, fmt.Errorf("dane: nieczytelny wiersz ostatniej migawki okna %q: %w", okno, err)
	}
	return migawka, nil
}

// Historia zwraca migawki okna od najświeższej — kolejne wiersze uporządkowane
// po utworzeniu są historią nawigacji. Limit 0 lub ujemny znaczy wykaz pełny.
func (r *repozytoriumPrzegladania) Historia(ctx context.Context, okno string, limit int) ([]MigawkaStrony, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, historiaMigawekOkna)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno, granicaWykazu(limit))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać historii okna %q: %w", okno, err)
	}
	defer wiersze.Close()

	lista := []MigawkaStrony{}
	for wiersze.Next() {
		migawka, err := odczytajMigawke(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz historii okna %q: %w", okno, err)
		}
		lista = append(lista, migawka)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt historii okna %q: %w", okno, err)
	}
	return lista, nil
}

// odczytajMigawke składa strukturę MigawkaStrony z jednego wiersza wyniku,
// w kolejności kolumn kolumnyMigawki.
func odczytajMigawke(wiersz skaner) (MigawkaStrony, error) {
	var migawka MigawkaStrony
	var tytul, tekstOdwolanie, zrodloOdwolanie, zrzutOdwolanie sql.NullString
	err := wiersz.Scan(&migawka.ID, &migawka.Kod, &migawka.Okno, &migawka.Url, &tytul,
		&tekstOdwolanie, &zrodloOdwolanie, &zrzutOdwolanie, &migawka.Utworzono)
	if err != nil {
		return MigawkaStrony{}, err
	}
	migawka.Tytul = tekstZKolumny(tytul)
	migawka.TekstOdwolanie = tekstZKolumny(tekstOdwolanie)
	migawka.ZrodloOdwolanie = tekstZKolumny(zrodloOdwolanie)
	migawka.ZrzutOdwolanie = tekstZKolumny(zrzutOdwolanie)
	return migawka, nil
}
