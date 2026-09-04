// Plik sprawdza, czy odmowa wpisu danaco przy braku binarium serwera narzędzi
// prowadzi do naprawy, która istnieje, przez pomiar ścieżek w drzewie, nie
// samego napisu.
package narzedzia

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// skryptyZniesione wylicza skrypty instalek natywnych usunięte z repozytorium
// wraz z modelem wdrożenia hybrydowego; odmowa nie ma prawa odesłać do żadnego.
var skryptyZniesione = []string{
	"instalka-natywna-win.sh",
	"instalka-natywna-linux.sh",
	"instalka-windows.sh",
	"instalka-pelna.sh",
	"wydanie.sh",
	"pakowanie.sh",
}

// korzenBudowy oddaje katalog budowy, korzeń modułu Go, od którego liczone są
// ścieżki wskazywane w odmowie, po potwierdzeniu, że trafił we właściwe
// miejsce.
func korzenBudowy(t *testing.T) string {
	t.Helper()
	korzen := filepath.Join("..", "..", "..")
	modul := filepath.Join(korzen, "go.mod")
	if _, err := os.Stat(modul); err != nil {
		t.Fatalf("sprawdzian nie stoi w drzewie budowy — nie ma %s (%v);"+
			" pomiaru obecności ścieżek z odmowy nie da się wykonać", modul, err)
	}
	return korzen
}

// powodBrakuBinarium wywołuje Wpis w warunkach, w których serwera narzędzi nie
// da się znaleźć, i oddaje powód odmowy.
func powodBrakuBinarium(t *testing.T) string {
	t.Helper()
	t.Setenv("PATH", t.TempDir())
	_, _, _, powod, jest := Wpis("okno-probne")
	if jest {
		t.Fatal("sprawdzian nie zmierzył odmowy: serwer narzędzi znalazł się" +
			" mimo pustej ścieżki wyszukiwania")
	}
	if powod == "" {
		t.Fatal("odmowa bez powodu — nie ma czego czytać")
	}
	return powod
}

// TestOdmowaBrakuBinariumWskazujeDrogeNaprawyKtoraIstnieje sprawdza, że każda
// ścieżka wymieniona w odmowie leży w drzewie.
func TestOdmowaBrakuBinariumWskazujeDrogeNaprawyKtoraIstnieje(t *testing.T) {
	korzen := korzenBudowy(t)
	powod := powodBrakuBinarium(t)

	for _, sciezka := range []string{zrodloBinarium, skryptPakietu} {
		if !strings.Contains(powod, sciezka) {
			t.Errorf("odmowa nie wskazuje %s; treść odmowy: %s", sciezka, powod)
			continue
		}
		if _, err := os.Stat(filepath.Join(korzen, sciezka)); err != nil {
			t.Errorf("odmowa wskazuje %s, a tego w drzewie nie ma (%v)"+
				" — droga naprawy prowadzi po nic", sciezka, err)
		}
	}
}

// TestOdmowaBrakuBinariumNieWskazujeZniesionychSkryptow pilnuje, żeby odmowa
// nie wróciła do skryptów, których w repozytorium nie ma.
func TestOdmowaBrakuBinariumNieWskazujeZniesionychSkryptow(t *testing.T) {
	powod := powodBrakuBinarium(t)

	for _, skrypt := range skryptyZniesione {
		if strings.Contains(powod, skrypt) {
			t.Errorf("odmowa odsyła do zniesionego skryptu %s; treść odmowy: %s",
				skrypt, powod)
		}
	}
}

// TestZniesioneSkryptyNieLezaWDrzewie zakotwicza sprawdzian poprzedni w stanie
// drzewa: gdyby któryś ze skryptów wrócił, zakaz odsyłania do niego przestałby
// mieć podstawę i trzeba by go rozważyć od nowa, a nie utrzymywać z rozpędu.
func TestZniesioneSkryptyNieLezaWDrzewie(t *testing.T) {
	korzen := korzenBudowy(t)

	for _, skrypt := range skryptyZniesione {
		sciezka := filepath.Join(korzen, "scripts", skrypt)
		if _, err := os.Stat(sciezka); err == nil {
			t.Errorf("%s leży w drzewie — podstawa zakazu odsyłania do niego"+
				" wymaga rozważenia od nowa", sciezka)
		}
	}
}

// TestWpisBezOknaOdmawiaZPowodem pilnuje drugiej odmowy tej samej funkcji:
// okno bez identyfikatora nie daje wpisu i mówi, dlaczego.
func TestWpisBezOknaOdmawiaZPowodem(t *testing.T) {
	polecenie, argumenty, srodowisko, powod, jest := Wpis("")
	if jest {
		t.Fatal("wpis powstał dla okna bez identyfikatora — nie miałby zasięgu")
	}
	if polecenie != "" || argumenty != nil || srodowisko != nil {
		t.Errorf("odmowa oddała polecenie %q, argumenty %v i środowisko %v — miała oddać nic",
			polecenie, argumenty, srodowisko)
	}
	if powod == "" {
		t.Error("odmowa bez powodu — dziennik rdzenia nie miałby czego zameldować")
	}
}
