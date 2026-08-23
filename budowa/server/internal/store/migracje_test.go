package store

import (
	"path/filepath"
	"testing"
)

// Sprawdziany warstwy trwałości.
//
// Migracje są jedyną drogą, którą schemat dojeżdża do maszyny Operatora, i nie
// mają wersji zapasowej: krok wykonany błędnie zostaje w bazie na zawsze, bo
// kroku nie da się cofnąć. Dlatego mierzone jest tu nie „czy przechodzi", lecz
// cztery obietnice, na których stoi cała reszta: pełny przejazd od zera,
// powtarzalność, niezmienność treści kroku już zastosowanego oraz atomowość
// kroku nieudanego.
//
// Sterownik jest czystym Go, więc każdy sprawdzian zakłada własny plik bazy
// w katalogu tymczasowym i przejeżdża komplet migracji od nowa. Koszt tego
// przejazdu jest na tyle mały, że nie opłaca się dzielić bazy między
// sprawdzianami — a baza dzielona zamieniłaby je w jeden sprawdzian
// zależny od kolejności.

// swiezaBaza zakłada bazę w katalogu sprawdzianu i doprowadza ją do bieżącej
// wersji schematu.
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
		suma, jest := zastosowane[krok.Wersja]
		if !jest {
			t.Errorf("krok %03d (%s) nie został odnotowany", krok.Wersja, krok.Nazwa)
			continue
		}
		if suma != krok.SumaKontrolna {
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
	for wersja, suma := range przed {
		if po[wersja] != suma {
			t.Errorf("powtórny przejazd zmienił sumę kontrolną kroku %03d", wersja)
		}
	}

	// Rejestr nie ma prawa dostać drugiego wiersza tej samej wersji — warunek
	// UNIQUE na kolumnie `wersja` odrzuciłby wpis, ale schemat wykonałby się
	// wcześniej. Liczony jest więc wiersz, nie klucz mapy.
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
// schematy dwóch instalacji bez jednego komunikatu — chyba że start się o to
// zatrzyma.
func TestZmienionyKrokPoZastosowaniuJestBledem(t *testing.T) {
	baza := swiezaBaza(t)

	kroki, err := wczytajMigracje()
	if err != nil {
		t.Fatalf("nie można odczytać wykazu migracji: %v", err)
	}
	pierwszy := kroki[0]

	// Podmiana idzie po stronie rejestru, bo plików wkompilowanych w binarium
	// zmienić się nie da. Skutek jest ten sam: suma zapisana rozjeżdża się
	// z sumą policzoną z treści.
	if _, err := baza.DB.Exec("UPDATE migracja SET suma_kontrolna = ? WHERE wersja = ?",
		"suma z innej treści", pierwszy.Wersja); err != nil {
		t.Fatalf("nie można podmienić sumy kontrolnej: %v", err)
	}

	if err := baza.Migruj(); err == nil {
		t.Fatal("przejazd migracji przeszedł mimo zmienionej treści kroku już zastosowanego")
	}
}

// TestNieudanyKrokNieZostawiaSladu dowodzi atomowości. Krok, który wywrócił się
// w połowie, nie może zostać odnotowany jako wykonany — inaczej kolejny start
// pominąłby go i baza zostałaby ze schematem połowicznym na zawsze.
func TestNieudanyKrokNieZostawiaSladu(t *testing.T) {
	baza := swiezaBaza(t)

	// Krok wykonuje jedno polecenie poprawne, a potem jedno niepoprawne.
	// Tabela z pierwszego polecenia jest tu miarą wycofania transakcji.
	wadliwy := migracja{
		Wersja: 999999,
		Nazwa:  "krok_sprawdzianu",
		Tresc: `CREATE TABLE slad_kroku_sprawdzianu (id INTEGER PRIMARY KEY);
		        TO NIE JEST POLECENIE SQL;`,
		SumaKontrolna: "suma kroku sprawdzianu",
	}

	if err := baza.zastosujMigracje(wadliwy); err == nil {
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

// TestWykazKrokowJestRosnacyIJednoznaczny sprawdza to, na czym stoi cała
// numeracja: kolejność rosnąca i brak dwóch kroków o tym samym numerze.
// Ciągłość numeracji nie jest wymagana i nie jest tu sprawdzana — luki
// powstają przy pracy równoległej i są z zamysłu.
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

// TestSpojnoscPoPelnymPrzejezdzie puszcza kontrolę, którą rdzeń wystawia
// diagnostyce: integralność pliku i zgodność kluczy obcych. Schemat złożony
// z kroków poprawnych osobno może być niespójny razem — najprościej wtedy, gdy
// krok późniejszy odwołuje się do tabeli usuniętej wcześniej.
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
