package store

import (
	"context"
	"path/filepath"
	"strings"
	"testing"
)

// Funkcja swiezaBaza zakłada bazę w katalogu sprawdzianu i doprowadza ją do bieżącej wersji schematu migracjami.
func swiezaBaza(t *testing.T) *Baza {
	t.Helper()
	baza, err := Otworz(filepath.Join(t.TempDir(), "dane.sqlite"))
	if err != nil {
		t.Fatalf("nie można otworzyć bazy: %v", err)
	}
	t.Cleanup(func() { _ = baza.Zamknij() })
	return baza
}

// TestPrzejazdMigracjiOdZera jest dowodem, że instalacja u Operatora, który
// nigdy nie widział tego produktu, wstanie. Pusta baza przechodzi wszystkie
// kroki i kończy z rejestrem równym wykazowi kroków.
func TestPrzejazdMigracjiOdZera(t *testing.T) {
	baza := swiezaBaza(t)

	kroki, err := wczytajMigracje()
	if err != nil {
		t.Fatalf("nie można odczytać wykazu migracji: %v", err)
	}
	zastosowane, err := baza.zastosowaneMigracje()
	if err != nil {
		t.Fatalf("nie można odczytać rejestru migracji: %v", err)
	}
	if len(zastosowane) != len(kroki) {
		t.Fatalf("rejestr niesie %d kroków, wykaz ma %d", len(zastosowane), len(kroki))
	}
	for _, krok := range kroki {
		wpis, jest := zastosowane[krok.Wersja]
		if !jest {
			t.Errorf("krok %03d (%s) nie został odnotowany", krok.Wersja, krok.Nazwa)
			continue
		}
		if wpis.Nazwa != krok.Nazwa {
			t.Errorf("krok %03d odnotowany pod nazwą %q zamiast %q", krok.Wersja, wpis.Nazwa, krok.Nazwa)
		}
		if wpis.SumaKontrolna != krok.SumaKontrolna {
			t.Errorf("krok %03d (%s) odnotowany z inną sumą kontrolną", krok.Wersja, krok.Nazwa)
		}
	}
}

// TestPowtornyPrzejazdNiczegoNieZmienia sprawdza idempotencję. Rdzeń woła
// Migruj przy każdym starcie, więc krok wykonany po raz drugi byłby awarią
// codzienną, nie brzegową.
func TestPowtornyPrzejazdNiczegoNieZmienia(t *testing.T) {
	baza := swiezaBaza(t)

	przed, err := baza.zastosowaneMigracje()
	if err != nil {
		t.Fatalf("nie można odczytać rejestru migracji: %v", err)
	}
	if err := baza.Migruj(); err != nil {
		t.Fatalf("powtórny przejazd migracji nie powiódł się: %v", err)
	}
	po, err := baza.zastosowaneMigracje()
	if err != nil {
		t.Fatalf("nie można odczytać rejestru migracji: %v", err)
	}

	if len(przed) != len(po) {
		t.Fatalf("powtórny przejazd zmienił liczbę kroków: %d → %d", len(przed), len(po))
	}
	for wersja, wpis := range przed {
		if po[wersja].SumaKontrolna != wpis.SumaKontrolna {
			t.Errorf("powtórny przejazd zmienił sumę kontrolną kroku %03d", wersja)
		}
	}

	// Rejestr nie ma prawa dostać drugiego wiersza tej samej wersji; liczony jest wiersz, nie klucz mapy.
	var wierszy int
	if err := baza.DB.QueryRow("SELECT count(*) FROM migracja").Scan(&wierszy); err != nil {
		t.Fatalf("nie można policzyć wierszy rejestru: %v", err)
	}
	if wierszy != len(po) {
		t.Errorf("rejestr ma %d wierszy przy %d wersjach — powtórzenie kroku", wierszy, len(po))
	}
}

// TestZmienionyKrokPoZastosowaniuJestBledem pilnuje niezmienności treści.
// Poprawiony plik migracji, która już poszła na maszynę Operatora, rozjeżdża
// schematy dwóch instalacji bez jednego komunikatu.
func TestZmienionyKrokPoZastosowaniuJestBledem(t *testing.T) {
	baza := swiezaBaza(t)

	kroki, err := wczytajMigracje()
	if err != nil {
		t.Fatalf("nie można odczytać wykazu migracji: %v", err)
	}
	pierwszy := kroki[0]

	// Podmiana idzie po stronie rejestru, bo plików wkompilowanych w binarium zmienić się nie da.
	if _, err := baza.DB.Exec("UPDATE migracja SET suma_kontrolna = ? WHERE wersja = ?",
		"suma z innej treści", pierwszy.Wersja); err != nil {
		t.Fatalf("nie można podmienić sumy kontrolnej: %v", err)
	}

	if err := baza.Migruj(); err == nil {
		t.Fatal("przejazd migracji przeszedł mimo zmienionej treści kroku już zastosowanego")
	}
}

// TestRozjazdNazwyOdmawiaStartuPoUzgodnieniu odtwarza bazę wdrożenia: znacznik
// uzgodnienia stoi, sumy w rejestrze są już przepisane, a wiersze pochodzą
// z innej linii numeracji. Kontrola nazw ma sięgnąć takiej bazy przy starcie.
func TestRozjazdNazwyOdmawiaStartuPoUzgodnieniu(t *testing.T) {
	baza := swiezaBaza(t)

	kroki, err := wczytajMigracje()
	if err != nil {
		t.Fatalf("nie można odczytać wykazu migracji: %v", err)
	}
	obcy := kroki[len(kroki)/2]

	if _, err := baza.DB.Exec("UPDATE migracja SET nazwa = ? WHERE wersja = ?",
		"awatar_konta_wlasciciela", obcy.Wersja); err != nil {
		t.Fatalf("nie można podmienić nazwy kroku w rejestrze: %v", err)
	}
	if !znacznikUzgodnieniaPostawiony(t, baza) {
		t.Fatal("znacznik uzgodnienia sum nie stanął po przejeździe od zera")
	}

	blad := baza.Migruj()
	if blad == nil {
		t.Fatal("przejazd migracji przeszedł mimo rozjazdu nazwy kroku przy postawionym znaczniku")
	}
	for _, oczekiwane := range []string{"awatar_konta_wlasciciela", obcy.Nazwa, "inną linią numeracji"} {
		if !strings.Contains(blad.Error(), oczekiwane) {
			t.Errorf("odmowa nie zawiera %q: %v", oczekiwane, blad)
		}
	}
}

// TestOdmowaNazwyNieZostawiaPrzepisanychSum sprawdza kolejność: kontrola nazw
// idzie przed uzgodnieniem, więc baza z obcej linii nie wychodzi z przejazdu
// z rejestrem częściowo przepisanym.
func TestOdmowaNazwyNieZostawiaPrzepisanychSum(t *testing.T) {
	baza := swiezaBaza(t)

	kroki, err := wczytajMigracje()
	if err != nil {
		t.Fatalf("nie można odczytać wykazu migracji: %v", err)
	}
	pierwszy, obcy := kroki[0], kroki[len(kroki)/2]

	if _, err := baza.DB.Exec("UPDATE migracja SET suma_kontrolna = ? WHERE wersja = ?",
		"suma sprzed normalizacji", pierwszy.Wersja); err != nil {
		t.Fatalf("nie można podmienić sumy kontrolnej: %v", err)
	}
	if _, err := baza.DB.Exec("UPDATE migracja SET nazwa = ? WHERE wersja = ?",
		"awatar_konta_wlasciciela", obcy.Wersja); err != nil {
		t.Fatalf("nie można podmienić nazwy kroku w rejestrze: %v", err)
	}
	zdejmijZnacznikUzgodnienia(t, baza)

	if err := baza.Migruj(); err == nil {
		t.Fatal("przejazd migracji przeszedł mimo rozjazdu nazwy kroku")
	}
	var suma string
	if err := baza.DB.QueryRow("SELECT suma_kontrolna FROM migracja WHERE wersja = ?",
		pierwszy.Wersja).Scan(&suma); err != nil {
		t.Fatalf("nie można odczytać sumy kontrolnej: %v", err)
	}
	if suma != "suma sprzed normalizacji" {
		t.Errorf("odmowa zostawiła przepisaną sumę kroku %03d", pierwszy.Wersja)
	}
	if znacznikUzgodnieniaPostawiony(t, baza) {
		t.Error("znacznik uzgodnienia sum stanął mimo odmowy")
	}
}

// TestZnacznikUzgodnieniaNieSchodziZUserVersion pilnuje, by strażnika treści
// kroku nie znosiło jedno polecenie na pliku bazy. PRAGMA user_version ginie
// przy zrzucie i wczytaniu bazy narzędziem sqlite3.
func TestZnacznikUzgodnieniaNieSchodziZUserVersion(t *testing.T) {
	baza := swiezaBaza(t)

	kroki, err := wczytajMigracje()
	if err != nil {
		t.Fatalf("nie można odczytać wykazu migracji: %v", err)
	}
	pierwszy := kroki[0]

	if _, err := baza.DB.Exec("UPDATE migracja SET suma_kontrolna = ? WHERE wersja = ?",
		"suma z innej treści", pierwszy.Wersja); err != nil {
		t.Fatalf("nie można podmienić sumy kontrolnej: %v", err)
	}
	if _, err := baza.DB.Exec("PRAGMA user_version = 0"); err != nil {
		t.Fatalf("nie można cofnąć znacznika w nagłówku pliku: %v", err)
	}

	if err := baza.Migruj(); err == nil {
		t.Fatal("cofnięcie PRAGMA user_version przepuściło zmienioną treść kroku już zastosowanego")
	}
}

// TestZnacznikUzgodnieniaPrzechodziZNaglowkaPliku wykazuje, że baza uzgodniona
// przed przeniesieniem znacznika do tabeli nie uzgadnia sum po raz drugi.
func TestZnacznikUzgodnieniaPrzechodziZNaglowkaPliku(t *testing.T) {
	baza := swiezaBaza(t)

	if _, err := baza.DB.Exec("DELETE FROM uzgodnienie_sum"); err != nil {
		t.Fatalf("nie można zdjąć znacznika z tabeli: %v", err)
	}
	if _, err := baza.DB.Exec("PRAGMA user_version = 1"); err != nil {
		t.Fatalf("nie można postawić znacznika w nagłówku pliku: %v", err)
	}

	uzgodnione, err := baza.znacznikUzgodnieniaStoi()
	if err != nil {
		t.Fatalf("nie można odczytać znacznika uzgodnienia sum: %v", err)
	}
	if !uzgodnione {
		t.Fatal("znacznik z nagłówka pliku nie został uznany za postawiony")
	}
	if !znacznikUzgodnieniaPostawiony(t, baza) {
		t.Error("znacznik z nagłówka pliku nie został przepisany do tabeli")
	}
}

// TestZastanaSierotaNieObciazaKrokuMigracji odtwarza bazę po ręcznej naprawie:
// wiersz bez rodzica stoi w niej przed przejazdem, a krok migracji tej tabeli
// nie dotyka. Stan zastany zdejmowany przed przejazdem oddziela jedno od drugiego.
func TestZastanaSierotaNieObciazaKrokuMigracji(t *testing.T) {
	baza := swiezaBaza(t)
	zycie := context.Background()

	tabeleSprawdzianu := migracja{
		Wersja: 999001,
		Nazwa:  "tabele_sprawdzianu",
		Tresc: `CREATE TABLE rodzic_sprawdzianu (id INTEGER PRIMARY KEY);
		        CREATE TABLE dziecko_sprawdzianu (
		            id        INTEGER PRIMARY KEY,
		            rodzic_id INTEGER NOT NULL REFERENCES rodzic_sprawdzianu(id));`,
		SumaKontrolna: "suma tabel sprawdzianu",
	}
	if _, err := baza.zastosujMigracje(tabeleSprawdzianu, nil); err != nil {
		t.Fatalf("krok zakładający tabele sprawdzianu nie powiódł się: %v", err)
	}

	polaczenie, err := baza.DB.Conn(zycie)
	if err != nil {
		t.Fatalf("nie można zająć połączenia: %v", err)
	}
	if _, err := polaczenie.ExecContext(zycie, "PRAGMA foreign_keys = off"); err != nil {
		t.Fatalf("nie można wygasić więzów: %v", err)
	}
	if _, err := polaczenie.ExecContext(zycie,
		"INSERT INTO dziecko_sprawdzianu (id, rodzic_id) VALUES (1, 999)"); err != nil {
		t.Fatalf("nie można wstawić wiersza bez rodzica: %v", err)
	}
	if _, err := polaczenie.ExecContext(zycie, "PRAGMA foreign_keys = on"); err != nil {
		t.Fatalf("nie można przywrócić więzów: %v", err)
	}
	if err := polaczenie.Close(); err != nil {
		t.Fatalf("nie można zwolnić połączenia: %v", err)
	}

	zastane, err := baza.zastaneNaruszeniaWiezow()
	if err != nil {
		t.Fatalf("nie można zdjąć stanu więzów przed przejazdem: %v", err)
	}
	if len(zastane) == 0 {
		t.Fatal("wiersz bez rodzica nie został zdjęty jako stan zastany")
	}

	osobnaTabela := migracja{
		Wersja:        999002,
		Nazwa:         "osobna_tabela_sprawdzianu",
		Tresc:         `CREATE TABLE osobna_tabela_sprawdzianu (id INTEGER PRIMARY KEY);`,
		SumaKontrolna: "suma osobnej tabeli sprawdzianu",
	}
	if _, err := baza.zastosujMigracje(osobnaTabela, nil); err == nil {
		t.Error("bez stanu zastanego krok obcej tabeli przeszedł mimo zastanej sieroty")
	}
	if _, err := baza.zastosujMigracje(osobnaTabela, zastane); err != nil {
		t.Errorf("krok niedotykający tabeli z sierotą został odrzucony: %v", err)
	}
}

// TestPrzybylaSierotaKonczyKrokOdmowa pilnuje, by odczyt stanu zastanego nie
// przepuścił wiersza bez rodzica, który powstał dopiero w tym kroku.
func TestPrzybylaSierotaKonczyKrokOdmowa(t *testing.T) {
	baza := swiezaBaza(t)

	zrywajacyWiezy := migracja{
		Wersja: 999003,
		Nazwa:  "zrywajacy_wiezy_sprawdzianu",
		Tresc: `CREATE TABLE rodzic_zerwany (id INTEGER PRIMARY KEY);
		        CREATE TABLE dziecko_zerwane (
		            id        INTEGER PRIMARY KEY,
		            rodzic_id INTEGER NOT NULL REFERENCES rodzic_zerwany(id));
		        INSERT INTO dziecko_zerwane (id, rodzic_id) VALUES (1, 999);`,
		SumaKontrolna: "suma kroku zrywającego więzy",
	}
	_, blad := baza.zastosujMigracje(zrywajacyWiezy, map[kluczNaruszenia]int{})
	if blad == nil {
		t.Fatal("krok zostawiający wiersz bez rodzica zakończył się powodzeniem")
	}
	if !strings.Contains(blad.Error(), "dziecko_zerwane") {
		t.Errorf("odmowa nie nazywa tabeli z wierszem bez rodzica: %v", blad)
	}
}

// TestKrokMigracjiZerujeKasowaneStrony wykazuje pragmę secure_delete widzianą
// przez treść kroku: bez niej wiersz skasowany krokiem zostaje czytelny na
// zwolnionych stronach pliku bazy.
func TestKrokMigracjiZerujeKasowaneStrony(t *testing.T) {
	baza := swiezaBaza(t)

	odczytPragmy := migracja{
		Wersja:        999004,
		Nazwa:         "odczyt_pragmy_sprawdzianu",
		Tresc:         `CREATE TABLE slad_secure_delete AS SELECT * FROM pragma_secure_delete;`,
		SumaKontrolna: "suma odczytu pragmy",
	}
	if _, err := baza.zastosujMigracje(odczytPragmy, nil); err != nil {
		t.Fatalf("krok odczytujący pragmę nie powiódł się: %v", err)
	}
	var wartosc int
	if err := baza.DB.QueryRow("SELECT secure_delete FROM slad_secure_delete").Scan(&wartosc); err != nil {
		t.Fatalf("nie można odczytać wartości pragmy zdjętej w kroku: %v", err)
	}
	if wartosc != 1 {
		t.Errorf("krok migracji szedł przy secure_delete=%d", wartosc)
	}
}

// TestNieudanyKrokNieZostawiaSladu dowodzi atomowości. Krok, który wywrócił się
// w połowie, nie może zostać odnotowany jako wykonany — inaczej kolejny start
// pominąłby go i baza zostałaby ze schematem połowicznym na zawsze.
func TestNieudanyKrokNieZostawiaSladu(t *testing.T) {
	baza := swiezaBaza(t)

	// Krok wykonuje jedno polecenie poprawne, a potem jedno niepoprawne, jako miarę wycofania transakcji.
	wadliwy := migracja{
		Wersja: 999999,
		Nazwa:  "krok_sprawdzianu",
		Tresc: `CREATE TABLE slad_kroku_sprawdzianu (id INTEGER PRIMARY KEY);
		        TO NIE JEST POLECENIE SQL;`,
		SumaKontrolna: "suma kroku sprawdzianu",
	}

	if _, err := baza.zastosujMigracje(wadliwy, nil); err == nil {
		t.Fatal("krok z niepoprawnym poleceniem zakończył się powodzeniem")
	}

	var odnotowanych int
	if err := baza.DB.QueryRow("SELECT count(*) FROM migracja WHERE wersja = ?",
		wadliwy.Wersja).Scan(&odnotowanych); err != nil {
		t.Fatalf("nie można odczytać rejestru: %v", err)
	}
	if odnotowanych != 0 {
		t.Error("nieudany krok został odnotowany w rejestrze migracji")
	}

	var tabel int
	if err := baza.DB.QueryRow(
		"SELECT count(*) FROM sqlite_master WHERE type='table' AND name='slad_kroku_sprawdzianu'").
		Scan(&tabel); err != nil {
		t.Fatalf("nie można odczytać schematu: %v", err)
	}
	if tabel != 0 {
		t.Error("nieudany krok zostawił tabelę — transakcja nie została wycofana")
	}
}

// TestWykazKrokowJestRosnacyIJednoznaczny sprawdza, na czym stoi cała numeracja: kolejność rosnącą i brak dwóch kroków o tym samym numerze.
func TestWykazKrokowJestRosnacyIJednoznaczny(t *testing.T) {
	kroki, err := wczytajMigracje()
	if err != nil {
		t.Fatalf("nie można odczytać wykazu migracji: %v", err)
	}
	if len(kroki) == 0 {
		t.Fatal("wykaz migracji jest pusty")
	}
	for i := 1; i < len(kroki); i++ {
		if kroki[i].Wersja <= kroki[i-1].Wersja {
			t.Errorf("krok %03d (%s) nie stoi po kroku %03d (%s)",
				kroki[i].Wersja, kroki[i].Nazwa, kroki[i-1].Wersja, kroki[i-1].Nazwa)
		}
	}
	for _, krok := range kroki {
		if krok.Nazwa == "" {
			t.Errorf("krok %03d nie ma nazwy", krok.Wersja)
		}
		if krok.Tresc == "" {
			t.Errorf("krok %03d (%s) ma pustą treść", krok.Wersja, krok.Nazwa)
		}
		if krok.SumaKontrolna == "" {
			t.Errorf("krok %03d (%s) nie ma sumy kontrolnej", krok.Wersja, krok.Nazwa)
		}
	}
}

// TestSpojnoscPoPelnymPrzejezdzie puszcza kontrolę integralności pliku i zgodności kluczy obcych po pełnym przejeździe migracji.
func TestSpojnoscPoPelnymPrzejezdzie(t *testing.T) {
	baza := swiezaBaza(t)
	if err := baza.SprawdzSpojnosc(); err != nil {
		t.Errorf("baza po pełnym przejeździe migracji jest niespójna: %v", err)
	}
}

// TestKluczeObceSaWlaczone sprawdza pragmę, od której zależy każdy warunek
// ON DELETE CASCADE w schemacie. Pragma jest w DSN, więc obowiązuje każde
// połączenie z puli — sprawdzian bierze połączenie z puli, nie zakłada nowego.
func TestKluczeObceSaWlaczone(t *testing.T) {
	baza := swiezaBaza(t)
	wynik, err := baza.PragmaTekstowa("foreign_keys")
	if err != nil {
		t.Fatalf("nie można odczytać pragmy: %v", err)
	}
	if wynik != "1" {
		t.Errorf("klucze obce są wyłączone (foreign_keys=%q) — kaskady schematu nie działają", wynik)
	}
}

// TestPustaSciezkaJestOdmowa pilnuje, by brak wskazania pliku bazy był odmową
// nazwaną, a nie bazą założoną gdziekolwiek.
func TestPustaSciezkaJestOdmowa(t *testing.T) {
	if _, err := Otworz("   "); err == nil {
		t.Error("pusta ścieżka pliku bazy została przyjęta")
	}
}

// znacznikUzgodnieniaPostawiony mówi, czy w bazie stoi wiersz znacznika uzgodnienia sum.
func znacznikUzgodnieniaPostawiony(t *testing.T, baza *Baza) bool {
	t.Helper()
	var wierszy int
	if err := baza.DB.QueryRow("SELECT count(*) FROM uzgodnienie_sum").Scan(&wierszy); err != nil {
		t.Fatalf("nie można odczytać znacznika uzgodnienia sum: %v", err)
	}
	return wierszy > 0
}

// zdejmijZnacznikUzgodnienia cofa bazę do stanu sprzed uzgodnienia sum, w obu miejscach naraz.
func zdejmijZnacznikUzgodnienia(t *testing.T, baza *Baza) {
	t.Helper()
	if _, err := baza.DB.Exec("DELETE FROM uzgodnienie_sum"); err != nil {
		t.Fatalf("nie można zdjąć znacznika z tabeli: %v", err)
	}
	if _, err := baza.DB.Exec("PRAGMA user_version = 0"); err != nil {
		t.Fatalf("nie można zdjąć znacznika z nagłówka pliku: %v", err)
	}
}
