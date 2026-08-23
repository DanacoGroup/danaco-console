package store

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// Zapora katalogu akcji.
//
// Panel akcji i siatka szybkich akcji nie mają własnej listy pozycji — biorą ją
// komendą `action.list` z tabeli `akcja`. Wiersz katalogu wskazuje komendę
// kolumną `komenda`, a kolumna ta nie jest kluczem obcym i być nim nie może:
// kontrakt mieszka w `shared/contract.json`, nie w bazie. Baza przyjmie więc
// każdą nazwę, także nazwę komendy, której nie ma.
//
// Skutek takiego wiersza u Operatora: kontrolka jest, daje się nacisnąć i wraca
// odmową `*.unknown`. To jest ATRAPA — element, który obiecuje czynność, a nie
// ma za sobą ani jednego wykonawcy. Wykaz braków w tym produkcie już raz mówił
// o brakach, których nie było; atrapa jest tym samym kłamstwem w drugą stronę.
//
// Dlatego kontrola stoi tutaj, a nie w kodzie rdzenia: pyta o stan po PEŁNYM
// przejeździe migracji, więc obejmuje każdy wiersz katalogu — także wiersz
// wniesiony migracją, która jeszcze nie istnieje.
//
// Sprawdzian wypada niepomyślnie także wtedy, gdy katalog zostanie opróżniony.
// Zaczyn akcji jest treścią produktu, nie danymi przykładowymi: pusty katalog
// znaczy panel akcji bez ani jednej pozycji w każdym oknie platformy.

// komendyKontraktu czyta nazwy wszystkich komend kontraktu.
func komendyKontraktu(t *testing.T) map[string]bool {
	t.Helper()

	tresc, err := os.ReadFile(filepath.Clean(sciezkaKontraktu))
	if err != nil {
		t.Fatalf("nie można odczytać kontraktu: %v", err)
	}
	var kontrakt struct {
		Komendy []struct {
			Typ string `json:"typ"`
		} `json:"komendy"`
	}
	if err := json.Unmarshal(tresc, &kontrakt); err != nil {
		t.Fatalf("kontrakt nie jest poprawnym JSON: %v", err)
	}
	nazwy := make(map[string]bool, len(kontrakt.Komendy))
	for _, komenda := range kontrakt.Komendy {
		if nazwa := strings.TrimSpace(komenda.Typ); nazwa != "" {
			nazwy[nazwa] = true
		}
	}
	if len(nazwy) == 0 {
		t.Fatal("kontrakt nie niesie ani jednej komendy — albo zmienił kształt, " +
			"albo sprawdzian czyta nie ten plik")
	}
	return nazwy
}

// pozycjaKatalogu jest wierszem tabeli `akcja` w zakresie, o który pyta zapora.
type pozycjaKatalogu struct {
	Kod     string
	Komenda string
	Ikona   string
	Poziom  string
}

// pozycjeKataloguAkcji czyta katalog akcji z bazy po pełnym przejeździe migracji.
func pozycjeKataloguAkcji(t *testing.T, baza *Baza) []pozycjaKatalogu {
	t.Helper()

	wiersze, err := baza.DB.Query(`
        SELECT a.kod, a.komenda, a.ikona, p.kod
          FROM akcja a
          JOIN poziom_zasiegu p ON p.id = a.poziom_zasiegu_id
         ORDER BY a.kod`)
	if err != nil {
		t.Fatalf("nie można odczytać katalogu akcji: %v", err)
	}
	defer wiersze.Close()

	pozycje := []pozycjaKatalogu{}
	for wiersze.Next() {
		var pozycja pozycjaKatalogu
		if err := wiersze.Scan(&pozycja.Kod, &pozycja.Komenda, &pozycja.Ikona, &pozycja.Poziom); err != nil {
			t.Fatalf("nie można odczytać wiersza katalogu akcji: %v", err)
		}
		pozycje = append(pozycje, pozycja)
	}
	if err := wiersze.Err(); err != nil {
		t.Fatalf("odczyt katalogu akcji przerwany: %v", err)
	}
	if len(pozycje) == 0 {
		t.Fatal("katalog akcji jest pusty — `action.list` oddawałby wykaz bez ani jednej " +
			"pozycji, a panele akcji nazywałyby to brakiem pozycji w każdym oknie platformy")
	}
	return pozycje
}

// TestKatalogAkcjiWskazujeKomendyKontraktu jest zaporą główną: każda pozycja
// katalogu wskazuje komendę, która w kontrakcie JEST.
//
// Zapory nie wolno osłabić wykazem wyjątków. Pozycja bez pokrycia w kontrakcie
// nie ma stanu przejściowego „jeszcze nie" — dopóki komendy nie ma, kontrolki
// też nie ma być, a wiersz dochodzi migracją razem z komendą.
func TestKatalogAkcjiWskazujeKomendyKontraktu(t *testing.T) {
	baza := swiezaBaza(t)

	for _, rozjazd := range pozycjeBezPokrycia(pozycjeKataloguAkcji(t, baza), komendyKontraktu(t)) {
		t.Errorf("pozycja katalogu akcji wskazuje komendę spoza kontraktu: %s — "+
			"u Operatora jest to kontrolka oddająca odmowę *.unknown; naprawa: wnieść "+
			"komendę do shared/contract.json wraz z uchwytem rdzenia albo zdjąć pozycję "+
			"z katalogu", rozjazd)
	}
}

// pozycjeBezPokrycia wylicza pozycje katalogu, których komendy nie ma
// w kontrakcie. Pozycja bez ani jednej komendy liczy się do tego samego wykazu:
// kontrolka bez komendy i kontrolka z komendą nieistniejącą kończą się u
// Operatora tym samym — naciśnięciem bez skutku.
//
// Wyliczenie stoi osobno od sprawdzianu po to, żeby dało się je nakarmić
// wierszem, którego w katalogu nie ma — patrz TestZaporaKataloguWykrywaAtrape.
func pozycjeBezPokrycia(pozycje []pozycjaKatalogu, komendy map[string]bool) []string {
	bezPokrycia := []string{}
	for _, pozycja := range pozycje {
		komenda := strings.TrimSpace(pozycja.Komenda)
		if komenda == "" {
			bezPokrycia = append(bezPokrycia, pozycja.Kod+" → (bez komendy)")
			continue
		}
		if !komendy[komenda] {
			bezPokrycia = append(bezPokrycia, pozycja.Kod+" → "+komenda)
		}
	}
	sort.Strings(bezPokrycia)
	return bezPokrycia
}

// TestZaporaKataloguWykrywaAtrape dowodzi, że zapora wyżej naprawdę zapiera.
//
// Sprawdzian, który przechodzi na katalogu zdrowym, ale przeszedłby też na
// katalogu z atrapą, jest sprawdzianem pozornym — a pozorny sprawdzian jest
// gorszy od jego braku, bo świeci zielono i nikt nie patrzy dalej. Dlatego droga
// niepomyślna jest tu mierzona wprost: wiersz wskazujący komendę, której nie ma,
// oraz wiersz bez komendy muszą wyjść z wyliczenia oba.
func TestZaporaKataloguWykrywaAtrape(t *testing.T) {
	komendy := komendyKontraktu(t)
	if komendy["komenda.ktorej.nie.ma"] {
		t.Fatal("kontrakt niesie komendę użytą w tym sprawdzianie jako nieistniejącą — " +
			"dobierz inną nazwę, bo sprawdzian przestał mierzyć drogę niepomyślną")
	}

	wykryte := pozycjeBezPokrycia([]pozycjaKatalogu{
		{Kod: "home.enter", Komenda: "home.enter"},
		{Kod: "atrapa.z.komenda.nieistniejaca", Komenda: "komenda.ktorej.nie.ma"},
		{Kod: "atrapa.bez.komendy", Komenda: "  "},
	}, komendy)

	if len(wykryte) != 2 {
		t.Fatalf("zapora wykryła %d atrap z dwóch wniesionych: %v", len(wykryte), wykryte)
	}
	if !strings.Contains(strings.Join(wykryte, "\n"), "komenda.ktorej.nie.ma") {
		t.Errorf("zapora nie nazwała komendy spoza kontraktu: %v", wykryte)
	}
}

// TestKatalogAkcjiNieZmyslaIkon pilnuje drugiej połowy tej samej obietnicy.
// Pozycja z ikoną, której nie ma w zestawie klienta, wychodzi w oknie kontrolką
// bez znaku — pustym prostokątem, którego Operator nie umie odczytać. Zestaw
// czyta się z plików źródłowych ikon, bo one są jedyną prawdą o tym, co klient
// umie narysować.
func TestKatalogAkcjiNieZmyslaIkon(t *testing.T) {
	baza := swiezaBaza(t)
	ikony := nazwyIkonKlienta(t)

	for _, pozycja := range pozycjeKataloguAkcji(t, baza) {
		if strings.TrimSpace(pozycja.Ikona) == "" {
			// Ikona pusta jest świadomym brakiem — panel stawia wtedy kontrolkę
			// samym napisem. To jest czytelne, więc nie jest usterką.
			continue
		}
		if !ikony[pozycja.Ikona] {
			t.Errorf("pozycja katalogu %s wskazuje ikonę %q, której nie ma w zestawie "+
				"klienta (client/src/ikony/zrodla) — w oknie wyjdzie kontrolka bez znaku",
				pozycja.Kod, pozycja.Ikona)
		}
	}
}

// sciezkaZrodelIkon wskazuje katalog źródeł ikon klienta, licząc od katalogu
// pakietu store.
const sciezkaZrodelIkon = "../../../client/src/ikony/zrodla"

// nazwyIkonKlienta czyta nazwy ikon z plików źródłowych zestawu.
//
// Odczyt idzie po kluczach zapisu obiektu — nazwa ikony stoi w tych plikach jako
// klucz wcięty dwoma znakami odstępu. Rozbiór składni TypeScriptu byłby tu
// kodem, który sam może się mylić; wzorzec klucza wystarcza, bo pliki źródeł mają
// jeden kształt i pilnuje go formater klienta.
func nazwyIkonKlienta(t *testing.T) map[string]bool {
	t.Helper()

	wpisy, err := os.ReadDir(filepath.Clean(sciezkaZrodelIkon))
	if err != nil {
		t.Fatalf("nie można przejrzeć źródeł ikon klienta: %v", err)
	}
	nazwy := map[string]bool{}
	for _, wpis := range wpisy {
		if wpis.IsDir() || !strings.HasSuffix(wpis.Name(), ".ts") {
			continue
		}
		tresc, err := os.ReadFile(filepath.Join(sciezkaZrodelIkon, wpis.Name()))
		if err != nil {
			t.Fatalf("nie można odczytać %s: %v", wpis.Name(), err)
		}
		for _, wiersz := range strings.Split(string(tresc), "\n") {
			if nazwa, jest := nazwaIkonyZWiersza(wiersz); jest {
				nazwy[nazwa] = true
			}
		}
	}
	if len(nazwy) == 0 {
		t.Fatal("źródła ikon klienta nie dały ani jednej nazwy — albo zmieniły kształt, " +
			"albo sprawdzian czyta nie ten katalog")
	}
	return nazwy
}

// nazwaIkonyZWiersza wyjmuje nazwę ikony z wiersza zapisu zestawu.
func nazwaIkonyZWiersza(wiersz string) (string, bool) {
	if !strings.HasPrefix(wiersz, "  ") || strings.HasPrefix(strings.TrimSpace(wiersz), "//") {
		return "", false
	}
	klucz, _, rozdzielone := strings.Cut(strings.TrimSpace(wiersz), ":")
	if !rozdzielone {
		return "", false
	}
	klucz = strings.Trim(strings.TrimSpace(klucz), "'\"")
	if klucz == "" || strings.ContainsAny(klucz, " (){}[],") {
		return "", false
	}
	return klucz, true
}
