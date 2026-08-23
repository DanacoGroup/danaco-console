package core

// Odpowiedzialność pliku: wypisać wykaz zależności zewnętrznych oraz arsenał
// mowy w postaci nadającej się do maszynowego odczytu, żeby prowizjonowanie
// serwera brało nazwy pakietów z tego samego rejestru, z którego bierze je sonda
// startowa.
//
// ── Po co osobny tryb, skoro sonda już wypisuje braki do dziennika ────────────
// Dziennik startu mówi Operatorowi, czego brakuje na TEJ maszynie — jest
// diagnozą stanu zastanego. Prowizjonowanie potrzebuje czego innego: pełnego
// wykazu pakietów DO POSTAWIENIA na serwerze docelowym, niezależnie od tego, co
// stoi na maszynie budującej. To ta sama wiedza (te same deklaracje narzędzi),
// ale wyprowadzona kompletnie i w postaci, którą skrypt rozbierze na pola.
//
// ── Dlaczego to jedyne miejsce nazw pakietów ──────────────────────────────────
// `zaleznosci_zewnetrzne.go` w nagłówku ostrzega: druga lista rozjedzie się
// z pierwszą przy pierwszej zmianie pakietu. Skrypt prowizjonowania
// (`scripts/arsenal-serwera.sh`) NIE przepisuje nazw — woła binarium rdzenia
// w tym trybie i konsumuje wynik. Zmiana pola `Pakiet` w deklaracji narzędzia
// dojeżdża więc i do sondy startowej, i do prowizjonowania jednym ruchem.
//
// ── Dlaczego warstwa jest liczona tutaj, a nie w skrypcie ─────────────────────
// Rozdział na warstwę obowiązkową i decyzyjną jest rozstrzygnięciem, nie
// formatowaniem: silnik kontenerów został WSTRZYMANY przez Właściciela i nie
// może zostać postawiony milcząco. Gdyby ten rozdział robił skrypt dopasowaniem
// napisów w powłoce, byłby drugą regułą obok deklaracji — nietypowaną,
// niesprawdzalną i cichą przy pomyłce. Tutaj jest jedną funkcją z jednym
// sprawdzianem, a skrypt tylko czyta gotową kolumnę.

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"danacoconsole/server/internal/mowa"
)

// Znaczniki wiersza poleceń przełączające rdzeń w tryb wypisania wykazu zamiast
// pracy właściwej. Przyjmowane w obu postaciach zapisu flagi (`-` i `--`), bo
// obie są w użyciu w tym projekcie.
const (
	znacznikWykazuZaleznosci = "wykaz-zaleznosci"
	znacznikWykazuMowy       = "wykaz-mowy"
)

// Nazwy warstw prowizjonowania. Warstwa mówi, CZYM postawić program — a przy
// silniku kontenerów: że nie stawiać go bez wyraźnego żądania.
const (
	// WarstwaObowiazkowa — pakiety dystrybucji, bez których moduły odmawiają.
	WarstwaObowiazkowa = "obowiazkowa-apt"
	// WarstwaWarsztatGo — programy dokładane przez `go install`.
	WarstwaWarsztatGo = "warsztat-go"
	// WarstwaWarsztatNpm — programy dokładane przez `npm i -g`.
	WarstwaWarsztatNpm = "warsztat-npm"
	// WarstwaSnap — programy, których dystrybucja nie ma w apt.
	WarstwaSnap = "snap"
	// WarstwaModelRecznie — silniki i wagi z wydań spoza repozytoriów
	// dystrybucji; kroki ręczne, bo automat na nich zawodzi.
	WarstwaModelRecznie = "model-recznie"
	// WarstwaDecyzyjna — silnik kontenerów, WSTRZYMANY decyzją Właściciela.
	// Prowizjonowanie go pomija, dopóki nie zażąda się go wprost.
	WarstwaDecyzyjna = "decyzyjna"
)

// ZadanoWykazZaleznosci mówi, czy w argumentach procesu stoi żądanie wypisania
// wykazu zależności. Sprawdzane przed odczytem konfiguracji, bo wykaz nie
// potrzebuje ani bazy, ani katalogu danych — a zestaw flag konfiguracji
// odrzuciłby ten argument jako nierozpoznany.
func ZadanoWykazZaleznosci(argumenty []string) bool {
	return zadanoZnacznik(argumenty, znacznikWykazuZaleznosci)
}

// ZadanoWykazMowy mówi, czy zażądano wypisania arsenału mowy.
func ZadanoWykazMowy(argumenty []string) bool {
	return zadanoZnacznik(argumenty, znacznikWykazuMowy)
}

// zadanoZnacznik odpowiada, czy argumenty niosą wskazany znacznik.
func zadanoZnacznik(argumenty []string, znacznik string) bool {
	for _, argument := range argumenty {
		if argument == "-"+znacznik || argument == "--"+znacznik {
			return true
		}
	}
	return false
}

// WarstwaZaleznosci rozstrzyga, do której warstwy prowizjonowania należy
// pozycja wykazu.
//
// Silnik kontenerów rozpoznajemy po programie, nie po podpowiedzi
// instalacyjnej: w wykazie stoją dwie jego deklaracje („Docker" warsztatu
// Developera i „Docker (klient wiersza poleceń)" modułu Terminal), niosą różne
// podpowiedzi, a obie mają trafić do warstwy decyzyjnej. Pozostałe warstwy
// bierzemy z treści podpowiedzi, bo ona już dziś mówi, czym program dociągnąć.
// Kolejność pytań jest istotna: podpowiedź `go install github.com/...` niesie
// adres GitHuba, a nie jest krokiem ręcznym.
func WarstwaZaleznosci(pozycja ZaleznoscZewnetrzna) string {
	program := strings.TrimSpace(pozycja.Narzedzie.Program)
	pakiet := strings.TrimSpace(pozycja.Narzedzie.Pakiet)

	if program == "docker" || program == "podman" {
		return WarstwaDecyzyjna
	}
	switch {
	case strings.Contains(pakiet, "go install"):
		return WarstwaWarsztatGo
	case strings.HasPrefix(pakiet, "npm "):
		return WarstwaWarsztatNpm
	case strings.Contains(pakiet, "(snap)"):
		return WarstwaSnap
	case strings.Contains(pakiet, "github.com"),
		strings.Contains(pakiet, "środowisku pythonowym"):
		return WarstwaModelRecznie
	default:
		return WarstwaObowiazkowa
	}
}

// WypiszWykazZaleznosci wypisuje komplet zależności zewnętrznych, po jednym
// wierszu na pozycję wykazu, w polach rozdzielonych znakiem tabulacji:
//
//	warstwa <TAB> program <TAB> pakiet <TAB> stoi <TAB> nazwa <TAB> zakres
//
// Wiersze komentarza zaczynają się od `#` — skrypt konsumujący je pomija.
// Kolejność jest ta sama, którą ustala `zaleznosciZewnetrzne` (alfabetyczna po
// nazwie czytelnej), więc dwa kolejne wywołania dają ten sam wykaz.
func WypiszWykazZaleznosci(wyjscie io.Writer) error {
	if _, err := fmt.Fprintln(wyjscie,
		"# wykaz zależności zewnętrznych rdzenia — pola: "+
			"warstwa\tprogram\tpakiet\tstoi\tnazwa\tzakres"); err != nil {
		return err
	}
	for _, pozycja := range ZaleznosciZewnetrzne() {
		stoi := "nie"
		if pozycja.Stoi {
			stoi = "tak"
		}
		if _, err := fmt.Fprintf(wyjscie, "%s\t%s\t%s\t%s\t%s\t%s\n",
			WarstwaZaleznosci(pozycja),
			pozycja.Narzedzie.Program,
			pozycja.Narzedzie.Pakiet,
			stoi,
			pozycja.Narzedzie.Nazwa,
			pozycja.Zakres,
		); err != nil {
			return err
		}
	}
	return nil
}

// WypiszWykazMowy wypisuje arsenał mowy — tę jego część, której nie widać
// w wykazie zależności zewnętrznych.
//
// Wykaz zależności niesie z mowy tylko syntezator zapasowy (eSpeak NG), bo tylko
// on jest zwykłym programem na ścieżce. Reszta arsenału mowy to piper wraz
// z plikami głosów, biblioteka pythonowa rozpoznawania (faster-whisper) i wagi
// jej modelu — rzeczy stawiane inaczej niż pakietem dystrybucji, a bez nich
// mikrofon i odsłuch odmawiają Operatorowi tak samo. Wartości pochodzą ze
// stałych rdzenia (nazwy programów, miejsca arsenału, zmienne wskazania) oraz
// z ustawień silnika mowy (`mowa.ModelDomyslny`) — nie są tu wpisane po raz
// drugi. Postać wiersza: klucz <TAB> wartość.
func WypiszWykazMowy(wyjscie io.Writer) error {
	wiersze := [][2]string{
		{"# arsenał mowy — pola: klucz", "wartość"},
		{"synteza.piper-program", silnikPiper},
		{"synteza.piper-arsenal", piperArsenalu},
		{"synteza.piper-glosy", glosyPiperaArsenalu},
		{"synteza.piper-zmienna", zmiennaPipera},
		{"synteza.piper-glosy-zmienna", zmiennaGlosowPiper},
		{"synteza.espeak-program", silnikEspeak},
		{"synteza.espeak-zmienna", zmiennaEspeaka},
		{"rozpoznanie.model-domyslny", mowa.ModelDomyslny},
		{"rozpoznanie.ustawienie-modelu", mowa.KluczModel},
		{"rozpoznanie.ustawienie-katalogu-modeli", mowa.KluczKatalogModeli},
		{"rozpoznanie.ustawienie-interpretera", mowa.KluczProgram},
	}
	wiersze = append(wiersze, wierszePomocnikaMowy()...)
	for _, wiersz := range wiersze {
		if _, err := fmt.Fprintf(wyjscie, "%s\t%s\n", wiersz[0], wiersz[1]); err != nil {
			return err
		}
	}
	return nil
}

// wierszePomocnikaMowy oddaje położenie pomocnika transkrypcji wraz z plikiem
// jego zależności pythonowych.
//
// Plik zależności wskazujemy ścieżką wyliczoną z położenia skryptu, a nie
// wypisaną wprost: `pomocniki/transkrypcja/wymagania.txt` jest jedynym miejscem,
// w którym stoi nazwa i wersja biblioteki rozpoznawania, więc prowizjonowanie ma
// go zainstalować przez `-r`, zamiast powtarzać nazwę pakietu u siebie.
//
// Gdy skryptu nie widać (rdzeń uruchomiony poza pakietem produktu), wypisujemy
// przeszukane miejsca — to jedyna wskazówka naprawy, którą warto podać
// prowizjonowaniu.
func wierszePomocnikaMowy() [][2]string {
	pomocnik, err := mowa.OdnajdzPomocnika("")
	if err != nil {
		var brak *mowa.BrakPomocnika
		if errors.As(err, &brak) {
			return [][2]string{
				{"rozpoznanie.pomocnik-skrypt", ""},
				{"rozpoznanie.pomocnik-szukano", strings.Join(brak.Szukano, " ")},
			}
		}
		return [][2]string{{"rozpoznanie.pomocnik-skrypt", ""}}
	}
	return [][2]string{
		{"rozpoznanie.pomocnik-skrypt", pomocnik.Skrypt},
		{"rozpoznanie.pomocnik-interpreter", pomocnik.Program},
		{"rozpoznanie.wymagania", filepath.Join(filepath.Dir(pomocnik.Skrypt), "wymagania.txt")},
	}
}
