// Odpowiedzialność pliku: egzekucja trzech punktów izolacji, które rozstrzygają
// się w chwili uruchomienia procesu okna — katalogu roboczego sesji, środowiska
// procesu i katalogu danych modelu.
//
// Miejsce egzekucji jest jedno: polecenie zbudowane przez warstwę kanału, zanim
// pójdzie do uruchamiacza (uruchamiacz.go). Polecenie niosące ścieżkę
// albo środowisko spoza obszaru okna zostaje ODRZUCONE — proces z takim
// poleceniem nie startuje.
package session

import (
	"strings"

	"danacoconsole/server/internal/konfig"
)

// zmienneKataloguDanychDomyslne wylicza zmienne środowiska, którymi kanał
// modelu dostaje położenie własnego katalogu danych. Wykaz jest danymi, nie
// kodem: obszar okna może go zastąpić własnym.
var zmienneKataloguDanychDomyslne = []string{"CLAUDE_CONFIG_DIR"}

// Obszar to wydzielone dla jednego okna komunikacji miejsce na dysku. Wypełnia
// go punkt kompozycji rdzenia; egzekutor obszaru nie wylicza i nie tworzy.
type Obszar struct {
	// IdOkna — okno, którego obszar opisuje struktura.
	IdOkna string
	// KatalogRoboczy — własny katalog roboczy okna.
	KatalogRoboczy string
	// KatalogDanych — własny katalog danych kanału modelu okna.
	KatalogDanych string
	// ZmienneKataloguDanych zastępuje wykaz domyślny. Puste znaczy wykaz
	// domyślny, nie brak sprawdzenia.
	ZmienneKataloguDanych []string
}

// SprawdzPolecenie egzekwuje izolację na poleceniu uruchomienia procesu okna.
// Zwraca polecenie dopuszczone do wykonania albo naruszenie.
//
// Wartość pusta, którą izolacja jednoznacznie wyznacza, zostaje uzupełniona —
// brak wskazania nie jest naruszeniem, bo brak danych ma dawać poprawny wynik,
// nie awarię. Wskazanie wychodzące poza obszar okna jest naruszeniem,
// bo jest decyzją sprzeczną z ustawieniem Operatora.
func SprawdzPolecenie(zasady Zasady, obszar Obszar, polecenie Polecenie) (Polecenie, error) {
	polecenie, err := sprawdzKatalogRoboczy(zasady, obszar, polecenie)
	if err != nil {
		return polecenie, err
	}
	if err := sprawdzSrodowisko(zasady, obszar, polecenie); err != nil {
		return polecenie, err
	}
	return polecenie, sprawdzZmienneDanych(zasady, obszar, polecenie)
}

// sprawdzKatalogRoboczy pilnuje, żeby proces okna startował we własnym katalogu.
func sprawdzKatalogRoboczy(zasady Zasady, obszar Obszar, polecenie Polecenie) (Polecenie, error) {
	if !zasady.KatalogRoboczy {
		return polecenie, nil
	}
	wlasny := strings.TrimSpace(obszar.KatalogRoboczy)
	if wlasny == "" {
		return polecenie, NoweNaruszenie(konfig.KluczIzolacjaKatalogRoboczy,
			"okno "+obszar.IdOkna+" nie ma wydzielonego katalogu roboczego, "+
				"więc proces poszedłby do katalogu wspólnego zasięgu")
	}
	if strings.TrimSpace(polecenie.Katalog) == "" {
		polecenie.Katalog = wlasny
		return polecenie, nil
	}
	if !SciezkaWewnatrz(wlasny, polecenie.Katalog) {
		return polecenie, NoweNaruszenie(konfig.KluczIzolacjaKatalogRoboczy,
			"katalog uruchomienia "+polecenie.Katalog+" leży poza własnym katalogiem okna "+wlasny)
	}
	return polecenie, nil
}

// sprawdzSrodowisko pilnuje, żeby proces okna nie odziedziczył środowiska
// rdzenia. Dziedziczenie jest dokładnie tym, co punkt izolacji wyłącza:
// przy własnym zestawie zmiennych wspólne środowisko serwera nie wchodzi.
func sprawdzSrodowisko(zasady Zasady, obszar Obszar, polecenie Polecenie) error {
	if !zasady.SrodowiskoProcesu || !polecenie.DziedziczSrodowisko {
		return nil
	}
	return NoweNaruszenie(konfig.KluczIzolacjaSrodowiskoProcesu,
		"polecenie okna "+obszar.IdOkna+" dziedziczy środowisko rdzenia, "+
			"a okno ma pracować na własnym zestawie zmiennych")
}

// sprawdzZmienneDanych pilnuje, żeby zmienne niosące katalog danych kanału
// wskazywały wyłącznie własny katalog danych okna.
func sprawdzZmienneDanych(zasady Zasady, obszar Obszar, polecenie Polecenie) error {
	if !zasady.KatalogDanychModelu {
		return nil
	}
	nazwy := zmienneKataloguDanych(obszar)
	for _, wpis := range polecenie.Srodowisko {
		nazwa, wartosc, jest := strings.Cut(wpis, "=")
		if !jest || !zawiera(nazwy, nazwa) {
			continue
		}
		if err := SprawdzKatalogDanych(zasady, obszar, wartosc); err != nil {
			return err
		}
	}
	return nil
}

// SprawdzKatalogDanych egzekwuje izolację katalogu danych modelu na jednej
// ścieżce. Wywołuje to zarówno sprawdzenie polecenia, jak i punkt, w którym
// kanał wybiera katalog konfiguracji konta — reguła jest jedna.
//
// Wskazanie puste jest naruszeniem: kanał bez wskazania sięga po katalog
// wspólny, a właśnie tego izolacja zabrania.
func SprawdzKatalogDanych(zasady Zasady, obszar Obszar, sciezka string) error {
	if !zasady.KatalogDanychModelu {
		return nil
	}
	wlasny := strings.TrimSpace(obszar.KatalogDanych)
	if wlasny == "" {
		return NoweNaruszenie(konfig.KluczIzolacjaKatalogDanych,
			"okno "+obszar.IdOkna+" nie ma wydzielonego katalogu danych modelu")
	}
	if strings.TrimSpace(sciezka) == "" {
		return NoweNaruszenie(konfig.KluczIzolacjaKatalogDanych,
			"brak wskazania katalogu danych modelu dla okna "+obszar.IdOkna+
				" znaczy katalog wspólny kanału")
	}
	if !SciezkaWewnatrz(wlasny, sciezka) {
		return NoweNaruszenie(konfig.KluczIzolacjaKatalogDanych,
			"katalog danych modelu "+sciezka+" leży poza własnym katalogiem danych okna "+wlasny)
	}
	return nil
}

// zmienneKataloguDanych zwraca wykaz obszaru albo wykaz domyślny.
func zmienneKataloguDanych(obszar Obszar) []string {
	if len(obszar.ZmienneKataloguDanych) > 0 {
		return obszar.ZmienneKataloguDanych
	}
	return zmienneKataloguDanychDomyslne
}

// zawiera odpowiada, czy wykaz niesie nazwę.
func zawiera(wykaz []string, nazwa string) bool {
	for _, pozycja := range wykaz {
		if pozycja == nazwa {
			return true
		}
	}
	return false
}
