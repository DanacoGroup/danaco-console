// Plik niesie dane wpisu danaco w konfiguracji MCP okna rozmowy; okno wchodzi
// argumentem uruchomienia, a poświadczenie i adres rdzenia — zmienną środowiska.
package narzedzia

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"danacoconsole/server/internal/transport"
)

const (
	// KluczWpisu nazywa wpis tego serwera narzędzi w mapie serwerów MCP
	// w konfiguracji sesji tej rozmowy okna.
	KluczWpisu = "danaco"
	// PrzelacznikOkna jest nazwą przełącznika uruchomieniowego, który niesie
	// identyfikator tego okna rozmowy.
	PrzelacznikOkna = "okno"
	// PrzelacznikRdzenia jest nazwą przełącznika uruchomieniowego niosącego
	// adres gniazda WebSocket rdzenia.
	PrzelacznikRdzenia = "rdzen"
	// Przełącznik zostaje dla uruchomienia ręcznego; wpis MCP podaje poświadczenie zmienną.
	PrzelacznikPoswiadczenia = "poswiadczenie"
	// Wiersz poleceń procesu czyta każdy program użytkownika, środowisko — właściciel procesu.
	ZmiennaPoswiadczenia = "DANACO_POSWIADCZENIE_NARZEDZI"
	// NazwaBinarium jest nazwą binarium serwera narzędzi, bez rozszerzenia
	// właściwego systemowi operacyjnemu.
	NazwaBinarium = "danaco-narzedzia"

	// zrodloBinarium nazywa pakiet Go, z którego binarium serwera narzędzi
	// powstaje. Ścieżka liczona od katalogu `budowa/` — korzenia modułu.
	zrodloBinarium = "server/cmd/danaco-narzedzia"
	// skryptPakietu nazywa skrypt składający pakiet serwera: buduje rdzeń
	// i serwer narzędzi, po czym stawia oba w JEDNYM katalogu, bo tam ich
	// szuka `sciezkaProgramu`. Ścieżka liczona od katalogu `budowa/`.
	skryptPakietu = "scripts/pakiet-serwera.sh"
)

// Wpis zwraca polecenie, argumenty i środowisko wpisu danaco dla wskazanego
// okna. Okno puste albo poświadczenie niewydane daje fałsz, wpisu wtedy nie
// dokłada się wcale.
func Wpis(idOkna string) (polecenie string, argumenty []string,
	srodowisko map[string]string, powod string, jest bool) {

	if idOkna == "" {
		return "", nil, nil, "okno rozmowy bez identyfikatora — wpis nie miałby zasięgu", false
	}
	poswiadczenie := transport.PoswiadczenieNarzedzi()
	if poswiadczenie == "" {
		return "", nil, nil,
			"poświadczenie serwera narzędzi niewydane — brak źródła losowego procesu", false
	}
	sciezka, err := sciezkaProgramu()
	if err != nil {
		return "", nil, nil, err.Error(), false
	}
	return sciezka,
		[]string{"--" + PrzelacznikOkna, idOkna},
		map[string]string{ZmiennaPoswiadczenia: poswiadczenie},
		"", true
}

// sciezkaProgramu wskazuje binarium serwera narzędzi albo mówi, czemu go nie
// ma, szukając go najpierw obok binarium, które właśnie pracuje.
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
type brakBinarium struct {
	Nazwa string
	Powod error
}

// Error mówi wprost, czego brakuje, gdzie tego szukano i skąd to wziąć,
// wskazując skrypt, który stawia oba binaria w jednym katalogu.
func (b *brakBinarium) Error() string {
	return "brak binarium serwera narzędzi " + b.Nazwa +
		" — nie ma go obok serwera ani na ścieżce wyszukiwania systemu;" +
		" powstaje z " + zrodloBinarium + " i ma stać w jednym katalogu z serwerem," +
		" co składa " + skryptPakietu
}

// Unwrap oddaje błąd źródłowy zwrócony przez wyszukanie tego binarium na
// ścieżce tego systemu operacyjnego.
func (b *brakBinarium) Unwrap() error { return b.Powod }
