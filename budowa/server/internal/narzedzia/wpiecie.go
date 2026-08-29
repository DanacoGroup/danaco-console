// Plik niesie dane wpisu danaco w konfiguracji MCP okna rozmowy; okno wchodzi
// argumentem uruchomienia, adres rdzenia wchodzi zmienną środowiska.
package narzedzia

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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

// Wpis zwraca polecenie i argumenty wpisu danaco dla wskazanego okna; okno
// puste daje fałsz, wpisu wtedy nie dokłada się wcale.
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
