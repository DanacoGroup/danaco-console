package store

import (
	"database/sql"
	"path/filepath"
	"strings"
	"testing"
)

// Kontrola spójności to jedyny przyrząd, którym rdzeń mówi, że plik bazy przestał trzymać się kupy.

// TestNaruszonyKluczObcyJestWykrywany zakłada wiersz-sierotę i sprawdza, czy kontrola spójności go widzi.
func TestNaruszonyKluczObcyJestWykrywany(t *testing.T) {
	sciezka := filepath.Join(t.TempDir(), "dane.sqlite")
	baza, err := Otworz(sciezka)
	if err != nil {
		t.Fatalf("nie można otworzyć bazy: %v", err)
	}
	t.Cleanup(func() { _ = baza.Zamknij() })

	if err := baza.SprawdzSpojnosc(); err != nil {
		t.Fatalf("baza świeżo zmigrowana jest już niespójna: %v", err)
	}

	zalozSierote(t, sciezka)

	err = baza.SprawdzSpojnosc()
	if err == nil {
		t.Fatal("kontrola spójności przeszła mimo wiersza-sieroty")
	}
	if !strings.Contains(err.Error(), "sesja") {
		t.Errorf("odmowa nie nazywa tabeli, w której leży sierota: %v", err)
	}
}

// Funkcja zalozSierote wstawia w tej bazie wiersz sesji wskazujący na nieistniejącą kartę tej samej sesji.
func zalozSierote(t *testing.T, sciezka string) {
	t.Helper()

	// DSN bez pragmy kluczy obcych to połączenie oboczne, wyłącznie do złożenia stanu spoza drogi zwykłej.
	oboczne, err := sql.Open(nazwaSterownika, filepath.ToSlash(sciezka)+"?_pragma=busy_timeout(5000)")
	if err != nil {
		t.Fatalf("nie można otworzyć połączenia obocznego: %v", err)
	}
	defer oboczne.Close()

	if _, err := oboczne.Exec(
		"INSERT INTO sesja (karta_sesji_id, tytul) VALUES (?, ?)",
		999999, "sesja bez karty"); err != nil {
		t.Fatalf("nie można założyć wiersza-sieroty: %v", err)
	}
}

// TestOpisNaruszeniaNazywaObieStrony sprawdza tekst, który Operator zobaczy
// w diagnostyce. Sam komunikat „niespójna" nie prowadzi do naprawy — prowadzi
// do niej wskazanie tabeli, wiersza i rodzica.
func TestOpisNaruszeniaNazywaObieStrony(t *testing.T) {
	opis := naruszenieKluczaObcego{
		Tabela:        "sesja",
		Wiersz:        7,
		TabelaRodzica: "karta_sesji",
		NumerKlucza:   0,
	}.String()

	for _, czesc := range []string{"sesja", "karta_sesji", "7"} {
		if !strings.Contains(opis, czesc) {
			t.Errorf("opis naruszenia nie zawiera %q: %s", czesc, opis)
		}
	}
}

// TestZamknijNaBazieZerowejJestBezpieczne pilnuje obietnicy z opisu metody.
// Ścieżka zamknięcia biegnie po nieudanym starcie, więc odbiornik zerowy jest
// tam stanem realnym, nie teoretycznym.
func TestZamknijNaBazieZerowejJestBezpieczne(t *testing.T) {
	var zerowa *Baza
	if err := zerowa.Zamknij(); err != nil {
		t.Errorf("zamknięcie bazy zerowej zgłosiło błąd: %v", err)
	}
	if err := (&Baza{}).Zamknij(); err != nil {
		t.Errorf("zamknięcie bazy bez puli zgłosiło błąd: %v", err)
	}
}

// TestPragmaNieznanaJestOdmowa sprawdza drogę błędu odczytu pragmy — tą samą,
// którą idzie kontrola integralności na pliku uszkodzonym.
func TestPragmaNieznanaJestOdmowa(t *testing.T) {
	baza := swiezaBaza(t)
	if _, err := baza.PragmaTekstowa("pragma_ktorej_nie_ma"); err == nil {
		t.Error("odczyt pragmy nieistniejącej zakończył się powodzeniem")
	}
}

// TestDsnNiesieKompletPragm pilnuje, by pragmy obowiązywały każde połączenie
// z puli. Pragma wykonana raz po otwarciu bazy objęłaby jedno połączenie
// z czterech — trzy pozostałe pracowałyby bez kluczy obcych i bez trybu WAL.
func TestDsnNiesieKompletPragm(t *testing.T) {
	dsn := zbudujDSN("/katalog/dane.sqlite")
	for _, pragma := range pragmyPolaczenia {
		if !strings.Contains(dsn, "_pragma=") {
			t.Fatalf("DSN nie niesie pragm wcale: %s", dsn)
		}
		// Wartości są w DSN zakodowane procentowo, więc porównywana jest
		// nazwa pragmy, nie cały zapis.
		nazwa, _, _ := strings.Cut(pragma, "(")
		if !strings.Contains(dsn, nazwa) {
			t.Errorf("DSN nie niesie pragmy %q: %s", nazwa, dsn)
		}
	}
}
