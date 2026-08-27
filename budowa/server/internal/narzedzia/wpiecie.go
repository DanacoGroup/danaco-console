// Odpowiedzialność pliku: dane wpisu `danaco` w konfiguracji MCP okna rozmowy.
//
// Okno wchodzi argumentem uruchomienia, nie zmienną środowiska. Podział
// powtarza to, co most `mcp-danaco-pulpit-console` już stosuje
// (`core/most_mcp.go`): to, co rozstrzyga zasięg jednego wpisu, idzie
// argumentem (tam tryb uprawnień), a to, co dzielą wszystkie wpisy, idzie
// środowiskiem (tam korzenie). Cztery powody:
//
//  1. Wpis powstaje osobno dla KAŻDEGO okna, a proces modelu jest jeden na
//     okno — argument należy do wpisu, więc różni się wpis po wpisie. Zmienna
//     środowiska należy do procesu i tej rozdzielczości nie ma.
//  2. Zmienną środowiska dziedziczy każdy proces potomny modelu; argument nie
//     wychodzi poza to jedno uruchomienie. Identyfikator okna jest uchwytem do
//     sterowania platformą i nie ma powodu, by wędrował dalej.
//  3. Argument widać w samej konfiguracji MCP: czytając plik `--mcp-config`
//     wiadomo, którego okna dotyczy. Wpis bez argumentu byłby dla wszystkich
//     okien identyczny i nie do odróżnienia.
//  4. Argument pominięty widać od razu w wierszu uruchomienia; zmienna pusta
//     jest nieodróżnialna od nieustawionej.
//
// Adres rdzenia idzie drogą przeciwną — środowiskiem — bo jest wspólny dla
// całej instalacji i czyta go ten sam pakiet `konfiguracja`, co w rdzeniu
// (zob. `adres.go`). Dwie różne rzeczy, dwie różne drogi.
package narzedzia

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

const (
	// KluczWpisu nazywa wpis serwera narzędzi w mapie `mcpServers`.
	KluczWpisu = "danaco"
	// PrzelacznikOkna jest nazwą przełącznika niosącego okno rozmowy.
	PrzelacznikOkna = "okno"
	// PrzelacznikRdzenia jest nazwą przełącznika niosącego adres gniazda rdzenia.
	PrzelacznikRdzenia = "rdzen"
	// NazwaBinarium jest nazwą binarium serwera narzędzi.
	NazwaBinarium = "danaco-narzedzia"

	// zrodloBinarium nazywa pakiet Go, z którego binarium serwera narzędzi
	// powstaje. Ścieżka liczona od katalogu `budowa/` — korzenia modułu.
	zrodloBinarium = "server/cmd/danaco-narzedzia"
	// skryptPakietu nazywa skrypt składający pakiet serwera: buduje rdzeń
	// i serwer narzędzi, po czym stawia oba w JEDNYM katalogu, bo tam ich
	// szuka `sciezkaProgramu`. Ścieżka liczona od katalogu `budowa/`.
	skryptPakietu = "scripts/pakiet-serwera.sh"
)

// Wpis zwraca polecenie i argumenty wpisu `danaco` dla wskazanego okna.
//
// Okno puste daje fałsz: serwer narzędzi bez okna nie miałby zasięgu, więc
// zamiast wpisu bez zasięgu lepiej wpisu nie dokładać wcale. Rozmowa toczy się
// wtedy bez sterowania platformą, a nie z narzędziami mierzącymi w nikąd.
//
// Wynik czwarty — `powod` — niesie powód odmowy. Ścieżkę binarium oddajemy
// dopiero po sprawdzeniu, że plik istnieje: gdy pakiet serwera złożono bez
// serwera narzędzi (`scripts/pakiet-serwera.sh`), wpis `danaco` wskazywałby
// plik, którego nie ma, a narzędzia sterowania platformą nie działałyby bez
// śladu. Odmowa musi być powiedziana wprost, bo brak wpisu i wpis martwy
// wyglądają dla operatora tak samo, a naprawa jest inna: złożyć pakiet serwera
// od nowa, nie szukać w rozmowie.
func Wpis(idOkna string) (polecenie string, argumenty []string, powod string, jest bool) {
	if idOkna == "" {
		return "", nil, "okno rozmowy bez identyfikatora — wpis nie miałby zasięgu", false
	}
	sciezka, err := sciezkaProgramu()
	if err != nil {
		return "", nil, err.Error(), false
	}
	return sciezka, []string{"--" + PrzelacznikOkna, idOkna}, "", true
}

// sciezkaProgramu wskazuje binarium serwera narzędzi albo mówi, czemu go nie ma.
//
// Szuka go OBOK binarium, które właśnie pracuje: rdzeń i serwer narzędzi
// wychodzą z jednego budowania i jadą w jednym pakiecie instalacyjnym, więc
// stoją w tym samym katalogu. Gdy obok go nie ma — albo gdy miejsca bieżącego
// procesu nie da się ustalić — zostaje ścieżka wyszukiwania systemu (droga
// uruchomienia z `go run`, gdzie proces stoi w katalogu tymczasowym). Dopiero
// gdy zawiodą OBIE drogi, funkcja zwraca błąd zamiast ścieżki: ścieżka zmyślona
// jest gorsza od braku wpisu, bo o braku wpisu da się powiedzieć.
func sciezkaProgramu() (string, error) {
	nazwa := NazwaBinarium
	if runtime.GOOS == "windows" {
		nazwa += ".exe"
	}
	if biezace, err := os.Executable(); err == nil {
		obok := filepath.Join(filepath.Dir(biezace), nazwa)
		if opis, err := os.Stat(obok); err == nil && !opis.IsDir() {
			return obok, nil
		}
	}
	zeSciezki, err := exec.LookPath(nazwa)
	if err != nil {
		return "", &brakBinarium{Nazwa: nazwa, Powod: err}
	}
	return zeSciezki, nil
}

// brakBinarium jest błędem: serwera narzędzi nie ma ani obok rdzenia, ani na
// ścieżce wyszukiwania systemu.
//
// Osobny typ, a nie sam napis, żeby wywołujący mógł ten jeden przypadek
// odróżnić od pozostałych odmów: to jedyna odmowa, która znaczy „produkt
// zbudowano lub spakowano niekompletnie”, a nie „to okno nie ma zasięgu”.
type brakBinarium struct {
	Nazwa string
	Powod error
}

// Error mówi wprost, czego brakuje, gdzie tego szukano i skąd to wziąć — bez
// ostatniej części meldunek w dzienniku rdzenia nazywa brak, ale do naprawy nie
// prowadzi. Droga wskazana tu jest jedyną, którą serwer narzędzi w produkcie
// powstaje: instalka Operatora nie niesie ani rdzenia, ani serwera narzędzi,
// więc oba stoją wyłącznie w pakiecie serwera.
func (b *brakBinarium) Error() string {
	return "brak binarium serwera narzędzi " + b.Nazwa +
		" — nie ma go obok rdzenia ani na ścieżce wyszukiwania systemu;" +
		" powstaje z " + zrodloBinarium + " i ma stać w jednym katalogu z rdzeniem," +
		" co składa " + skryptPakietu
}

// Unwrap oddaje błąd źródłowy z `exec.LookPath`.
func (b *brakBinarium) Unwrap() error { return b.Powod }
