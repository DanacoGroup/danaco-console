package store

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Kontrakt deklaruje, że wyliczenie z odpowiednikiem w schemacie niesie pole kolumnaBazy i pole baza.

// sciezkaKontraktu wskazuje jedyne źródło prawdy nazw kontraktu, liczone od katalogu tego pakietu store.
const sciezkaKontraktu = "../../../shared/contract.json"

// wyliczenieKontraktu jest wycinkiem kontraktu potrzebnym temu sprawdzianowi.
// Wycinek, nie pełny model: sprawdzian pyta wyłącznie o odwzorowanie bazy,
// a pełna struktura kontraktu należy do generatora.
type wyliczenieKontraktu struct {
	Nazwa       string `json:"nazwa"`
	KolumnaBazy string `json:"kolumnaBazy"`
	Wartosci    []struct {
		Wartosc string `json:"wartosc"`
		Baza    string `json:"baza"`
		// Przelotowa oznacza wartość, która przechodzi przez strumień i nigdy nie trafia do kolumny bazy.
		Przelotowa bool `json:"przelotowa"`
	} `json:"wartosci"`
}

// odwzorowaniaRozeszlyeSieZeSchematem wylicza wyliczenia, których pole kolumnaBazy wskazuje dziś na nieistniejącą tabelę albo kolumnę.
var odwzorowaniaRozeszlyeSieZeSchematem = map[string]string{
	"ProgressStatus": "tabela proces_sesji zdjęta krokiem sieroty_transportu; rejestr procesów żyje w pamięci",

	// Automations — schemat modułu wszedł migracjami 260-275, więc wierszy rozjazdu tu nie ma.

	// Diagnostics — prowenancja wywołań ma już swoje tabele, więc jej wiersze zeszły stąd razem z powodem.
	"AlertRuleKind":      "tabela alert_regula powstaje z migracją modułu Diagnostics",
	"AlertMetric":        "tabela alert_regula powstaje z migracją modułu Diagnostics",
	"AlertComparison":    "tabela alert_regula powstaje z migracją modułu Diagnostics",
	"AlertChannel":       "tabela alert_regula_kanal powstaje z migracją modułu Diagnostics",
	"AlertTriggerStatus": "tabela alert_wyzwolenie powstaje z migracją modułu Diagnostics",
	"HealthProbeKind":    "tabela kondycja_sonda powstaje z migracją modułu Diagnostics",
	"HealthProbeStatus":  "tabela kondycja_wynik powstaje z migracją modułu Diagnostics",

	// Library — schemat modułu wszedł migracjami 180-186, więc wierszy rozjazdu tu nie ma.
}

// TestKazdaKolumnaOdwzorowaniaIstniejeWSchemacie sprawdza pierwszą połowę
// obietnicy: wskazana tabela i kolumna są w bazie po pełnym przejeździe
// migracji.
func TestKazdaKolumnaOdwzorowaniaIstniejeWSchemacie(t *testing.T) {
	baza := swiezaBaza(t)
	wyliczenia := wyliczeniaOdwzorowane(t)

	for _, wyliczenie := range wyliczenia {
		t.Run(wyliczenie.Nazwa, func(t *testing.T) {
			tabela, kolumna := rozlozOdwolanie(t, wyliczenie.KolumnaBazy)
			obecne := tabelaIstnieje(t, baza, tabela) && kolumnaIstnieje(t, baza, tabela, kolumna)
			powod, uznany := odwzorowaniaRozeszlyeSieZeSchematem[wyliczenie.Nazwa]

			if uznany && obecne {
				t.Fatalf("odwzorowanie %s → %s znów zgadza się ze schematem — "+
					"zdejmij je z wykazu odwzorowaniaRozeszlyeSieZeSchematem",
					wyliczenie.Nazwa, wyliczenie.KolumnaBazy)
			}
			if uznany {
				t.Skipf("odwzorowanie %s → %s jest rozjazdem uznanym: %s",
					wyliczenie.Nazwa, wyliczenie.KolumnaBazy, powod)
			}
			if !tabelaIstnieje(t, baza, tabela) {
				t.Fatalf("kontrakt odwzorowuje %s na tabelę %q, której w schemacie nie ma",
					wyliczenie.Nazwa, tabela)
			}
			if !kolumnaIstnieje(t, baza, tabela, kolumna) {
				t.Fatalf("tabela %q nie ma kolumny %q wskazanej przez kontrakt (%s)",
					tabela, kolumna, wyliczenie.Nazwa)
			}
		})
	}
}

// TestKazdaWartoscOdwzorowaniaPrzechodziPrzezWarunek sprawdza, że wartość zapisana w kontrakcie jako baza jest wartością, którą kolumna przyjmuje.
func TestKazdaWartoscOdwzorowaniaPrzechodziPrzezWarunek(t *testing.T) {
	baza := swiezaBaza(t)
	wyliczenia := wyliczeniaOdwzorowane(t)

	for _, wyliczenie := range wyliczenia {
		t.Run(wyliczenie.Nazwa, func(t *testing.T) {
			if powod, uznany := odwzorowaniaRozeszlyeSieZeSchematem[wyliczenie.Nazwa]; uznany {
				t.Skipf("odwzorowanie %s → %s jest rozjazdem uznanym: %s",
					wyliczenie.Nazwa, wyliczenie.KolumnaBazy, powod)
			}
			tabela, kolumna := rozlozOdwolanie(t, wyliczenie.KolumnaBazy)
			schemat := schematTabeli(t, baza, tabela)
			if !strings.Contains(schemat, "CHECK("+kolumna+" ") {
				t.Skipf("kolumna %s.%s nie ma warunku CHECK — nie ma czego porównać", tabela, kolumna)
			}
			for _, wartosc := range wyliczenie.Wartosci {
				// Wartość przelotowa nigdy nie trafia do kolumny; kontrakt mówi to wprost, więc jej brak jest zgodny.
				if wartosc.Przelotowa {
					continue
				}
				if wartosc.Baza == "" {
					t.Errorf("wartość %q wyliczenia %s nie ma odpowiednika w bazie "+
						"ani nie jest oznaczona jako przelotowa",
						wartosc.Wartosc, wyliczenie.Nazwa)
					continue
				}
				if !strings.Contains(schemat, "'"+wartosc.Baza+"'") {
					t.Errorf("warunek kolumny %s.%s nie dopuszcza wartości %q, którą kontrakt "+
						"przypisuje wartości %q wyliczenia %s",
						tabela, kolumna, wartosc.Baza, wartosc.Wartosc, wyliczenie.Nazwa)
				}
			}
		})
	}
}

// Funkcja wyliczeniaOdwzorowane czyta z kontraktu wyliczenia niosące odwzorowanie na schemat bazy danych.
func wyliczeniaOdwzorowane(t *testing.T) []wyliczenieKontraktu {
	t.Helper()

	tresc, err := os.ReadFile(filepath.Clean(sciezkaKontraktu))
	if err != nil {
		t.Fatalf("nie można odczytać kontraktu: %v", err)
	}
	var kontrakt struct {
		Wyliczenia []wyliczenieKontraktu `json:"wyliczenia"`
	}
	if err := json.Unmarshal(tresc, &kontrakt); err != nil {
		t.Fatalf("kontrakt nie jest poprawnym JSON: %v", err)
	}

	odwzorowane := make([]wyliczenieKontraktu, 0, len(kontrakt.Wyliczenia))
	for _, wyliczenie := range kontrakt.Wyliczenia {
		if wyliczenie.KolumnaBazy != "" {
			odwzorowane = append(odwzorowane, wyliczenie)
		}
	}
	if len(odwzorowane) == 0 {
		t.Fatal("kontrakt nie niesie ani jednego wyliczenia z odwzorowaniem bazy — " +
			"albo zmienił kształt, albo sprawdzian czyta nie ten plik")
	}
	return odwzorowane
}

// Funkcja rozlozOdwolanie rozdziela zapis odwołania w postaci tabela.kolumna na dwie osobne jego części.
func rozlozOdwolanie(t *testing.T, odwolanie string) (string, string) {
	t.Helper()
	tabela, kolumna, rozdzielone := strings.Cut(odwolanie, ".")
	if !rozdzielone || tabela == "" || kolumna == "" {
		t.Fatalf("odwołanie %q nie ma postaci tabela.kolumna", odwolanie)
	}
	return tabela, kolumna
}

// Funkcja tabelaIstnieje pyta schemat bazy o istnienie tabeli o podanej nazwie po przejeździe migracji.
func tabelaIstnieje(t *testing.T, baza *Baza, tabela string) bool {
	t.Helper()
	var liczba int
	if err := baza.DB.QueryRow(
		"SELECT count(*) FROM sqlite_master WHERE type='table' AND name=?", tabela).
		Scan(&liczba); err != nil {
		t.Fatalf("nie można odczytać schematu: %v", err)
	}
	return liczba > 0
}

// Funkcja kolumnaIstnieje pyta schemat bazy o istnienie kolumny wskazanej tabeli po przejeździe migracji.
func kolumnaIstnieje(t *testing.T, baza *Baza, tabela, kolumna string) bool {
	t.Helper()
	wiersze, err := baza.DB.Query("SELECT name FROM pragma_table_info(?)", tabela)
	if err != nil {
		t.Fatalf("nie można odczytać kolumn tabeli %q: %v", tabela, err)
	}
	defer wiersze.Close()

	for wiersze.Next() {
		var nazwa string
		if err := wiersze.Scan(&nazwa); err != nil {
			t.Fatalf("nieczytelny wynik pragma_table_info: %v", err)
		}
		if nazwa == kolumna {
			return true
		}
	}
	if err := wiersze.Err(); err != nil {
		t.Fatalf("przerwany odczyt kolumn tabeli %q: %v", tabela, err)
	}
	return false
}

// schematTabeli zwraca zapis CREATE TABLE po pełnym przejeździe migracji —
// czyli po wszystkich przebudowach tabeli, jakie kroki na niej wykonały.
func schematTabeli(t *testing.T, baza *Baza, tabela string) string {
	t.Helper()
	var schemat string
	if err := baza.DB.QueryRow(
		"SELECT sql FROM sqlite_master WHERE type='table' AND name=?", tabela).
		Scan(&schemat); err != nil {
		t.Fatalf("nie można odczytać schematu tabeli %q: %v", tabela, err)
	}
	return schemat
}
