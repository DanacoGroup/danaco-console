package store

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// Zapora katalogu akcji sprawdza, czy każdy wiersz katalogu wskazuje komendę istniejącą w kontrakcie.

// Funkcja komendyKontraktu czyta nazwy wszystkich komend zdefiniowanych w kontrakcie na potrzeby sprawdzianu.
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

// pozycjaKatalogu jest wierszem tabeli akcja w zakresie kolumn, o który dokładnie pyta zapora katalogu.
type pozycjaKatalogu struct {
	Kod     string
	Komenda string
	Ikona   string
	Poziom  string
}

// Funkcja pozycjeKataloguAkcji czyta katalog akcji z bazy po pełnym przejeździe wszystkich jej migracji.
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

// TestKatalogAkcjiWskazujeKomendyKontraktu jest zaporą główną: każda pozycja katalogu wskazuje komendę, która istnieje w kontrakcie.
func TestKatalogAkcjiWskazujeKomendyKontraktu(t *testing.T) {
	baza := swiezaBaza(t)

	for _, rozjazd := range pozycjeBezPokrycia(pozycjeKataloguAkcji(t, baza), komendyKontraktu(t)) {
		t.Errorf("pozycja katalogu akcji wskazuje komendę spoza kontraktu: %s — "+
			"u Operatora jest to kontrolka oddająca odmowę *.unknown; naprawa: wnieść "+
			"komendę do shared/contract.json wraz z uchwytem rdzenia albo zdjąć pozycję "+
			"z katalogu", rozjazd)
	}
}

// Funkcja pozycjeBezPokrycia wylicza pozycje katalogu, których komendy nie ma w kontrakcie, licząc też pozycje bez żadnej komendy.
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

// TestZaporaKataloguWykrywaAtrape dowodzi, że zapora katalogu akcji naprawdę wykrywa wiersz udający działającą kontrolkę.
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

// TestKatalogAkcjiNieZmyslaIkon pilnuje, aby każda ikona pozycji katalogu istniała w zestawie ikon klienta.
func TestKatalogAkcjiNieZmyslaIkon(t *testing.T) {
	if !zrodlaIkonKlientaStoja() {
		t.Skipf("klient nie ma jeszcze zestawu ikon — brak katalogu %s; sprawdzian "+
			"wraca sam, gdy zestaw wejdzie wraz z ramą aplikacji", sciezkaZrodelIkon)
	}

	baza := swiezaBaza(t)
	ikony := nazwyIkonKlienta(t)

	for _, pozycja := range pozycjeKataloguAkcji(t, baza) {
		if strings.TrimSpace(pozycja.Ikona) == "" {
			// Ikona pusta jest świadomym brakiem; panel stawia wtedy kontrolkę samym napisem, co nie jest usterką.
			continue
		}
		if !ikony[pozycja.Ikona] {
			t.Errorf("pozycja katalogu %s wskazuje ikonę %q, której nie ma w zestawie "+
				"klienta (%s) — w oknie wyjdzie kontrolka bez znaku",
				pozycja.Kod, pozycja.Ikona, sciezkaZrodelIkon)
		}
	}
}

// sciezkaZrodelIkon wskazuje katalog źródeł ikon klienta, liczony od katalogu tego pakietu store aplikacji.
const sciezkaZrodelIkon = "../../../klient/src/ikony/zrodla"

// Funkcja zrodlaIkonKlientaStoja orzeka, czy klient ma już zestaw ikon, będący warunkiem powrotu sprawdzianu wyżej.
func zrodlaIkonKlientaStoja() bool {
	opis, err := os.Stat(filepath.Clean(sciezkaZrodelIkon))
	return err == nil && opis.IsDir()
}

// Funkcja nazwyIkonKlienta czyta nazwy wszystkich ikon zapisane w plikach źródłowych zestawu ikon klienta.
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

// Funkcja nazwaIkonyZWiersza wyjmuje nazwę ikony z jednego wiersza tekstu pliku zapisu zestawu ikon klienta.
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
