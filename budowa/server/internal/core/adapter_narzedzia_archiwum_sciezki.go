// Plik podaje formaty, ścieżki i odmowy rodziny archive: rozstrzygnięcia wspólne dla pakowania i rozpakowania archiwum w katalogu roboczym okna.
package core

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/zewnetrzne"
	"danacoconsole/shared"
)

// formatArchiwum nazywa jeden z trzech formatów archiwum, które ta rodzina czynności obsługuje w rejestrze narzędzi rdzenia.
type formatArchiwum string

const (
	formatZip   formatArchiwum = "zip"
	format7z    formatArchiwum = "7z"
	formatTarGz formatArchiwum = "tar.gz"
)

// rozpoznajFormatArchiwum czyta sygnaturę pliku i nazywa format zip, 7z albo tar.gz na podstawie zawartości, a nie nazwy pliku.
func rozpoznajFormatArchiwum(sciezka string) (formatArchiwum, error) {
	plik, err := os.Open(sciezka)
	if err != nil {
		return "", bladWskazaniaArchiwum("nie można otworzyć archiwum " + sciezka + ": " + err.Error())
	}
	defer plik.Close()

	naglowek := make([]byte, 6)
	n, err := io.ReadFull(plik, naglowek)
	if err != nil && n < 2 {
		return "", bladWskazaniaArchiwum("plik " + sciezka +
			" jest za krótki, żeby być archiwum — nie ma w nim nawet sygnatury formatu")
	}
	naglowek = naglowek[:n]

	switch {
	case len(naglowek) >= 4 && string(naglowek[:4]) == "PK\x03\x04":
		return formatZip, nil
	case len(naglowek) >= 6 && string(naglowek[:6]) == "7z\xbc\xaf\x27\x1c":
		return format7z, nil
	case len(naglowek) >= 2 && naglowek[0] == 0x1f && naglowek[1] == 0x8b:
		return formatTarGz, nil
	}
	return "", bladWskazaniaArchiwum("plik " + sciezka +
		" nie jest archiwum w żadnym z obsługiwanych formatów (zip, 7z, tar.gz) — " +
		"rozpoznanie idzie po zawartości pliku, nie po jego nazwie")
}

// rozstrzygnijFormatDocelowy przekłada pole format żądania na jeden z trzech formatów archiwum, biorąc zip przy braku wskazania.
func rozstrzygnijFormatDocelowy(zadany *string) (formatArchiwum, error) {
	if bezWartosci(zadany) {
		return formatZip, nil
	}
	nazwa := strings.ToLower(strings.TrimSpace(*zadany))
	nazwa = strings.TrimPrefix(nazwa, ".")
	switch nazwa {
	case "zip":
		return formatZip, nil
	case "7z":
		return format7z, nil
	case "tar.gz", "targz", "tgz", "gz":
		return formatTarGz, nil
	}
	return "", bladWskazaniaArchiwum("format " + nazwa +
		" nie jest obsługiwany — narzędzie pakuje do zip, 7z albo tar.gz")
}

// nazwaBezKatalogow wycina samą nazwę pliku ze ścieżki wskazanej w żądaniu, odrzucając wszystkie prowadzące katalogi.
func nazwaBezKatalogow(sciezka string) string {
	sciezka = strings.TrimRight(sciezka, "/\\")
	if i := strings.LastIndexAny(sciezka, "/\\"); i >= 0 {
		return sciezka[i+1:]
	}
	return sciezka
}

// sciezkaWzgledemKatalogu rozstrzyga ścieżkę żądania względem katalogu okna i pilnuje, żeby wynik z tego katalogu nie wyszedł.
func sciezkaWzgledemKatalogu(katalog, zadana string) (string, error) {
	zadana = strings.TrimSpace(zadana)
	if zadana == "" {
		return "", bladWskazaniaArchiwum("wskazana ścieżka jest pusta")
	}
	if filepath.IsAbs(zadana) || strings.HasPrefix(zadana, "/") || strings.HasPrefix(zadana, `\`) ||
		(len(zadana) >= 2 && zadana[1] == ':') {
		return "", bladWskazaniaArchiwum("ścieżka " + zadana +
			" jest bezwzględna — narzędzia archiwum pracują wyłącznie w katalogu roboczym okna, " +
			"a nie w dowolnym miejscu maszyny Operatora")
	}
	pelna := filepath.Clean(filepath.Join(katalog, zadana))
	wzgledna, err := filepath.Rel(katalog, pelna)
	if err != nil || wzgledna == ".." || strings.HasPrefix(wzgledna, ".."+string(filepath.Separator)) {
		return "", bladWskazaniaArchiwum("ścieżka " + zadana +
			" wychodzi poza katalog roboczy okna — narzędzia archiwum poza niego nie sięgają")
	}
	return pelna, nil
}

// bladArsenaluArchiwum przekłada odmowę pakietu narzędzi zewnętrznych na kod kontraktu, rozróżniając brak binarium od niepowodzenia programu.
func bladArsenaluArchiwum(err error) error {
	var brak *zewnetrzne.BrakNarzedzia
	if errors.As(err, &brak) {
		return bladArchiwumNiedostepnego(brak.Error())
	}
	return bladPrzetwarzaniaArchiwum(err.Error())
}

// bladArchiwumNiedostepnego znakuje zaplecze niedostępne: żądanie było
// poprawne, produkt nie jest zepsuty, brakuje czegoś w instalacji.
func bladArchiwumNiedostepnego(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
		"narzędzia archiwum: "+powod))
}

// bladWskazaniaArchiwum nazywa brak albo niepoprawność wskazania w żądaniu, w tym obie granice bezpieczeństwa katalogu i rozmiaru archiwum.
func bladWskazaniaArchiwum(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"narzędzia archiwum: "+powod))
}

// bladPrzetwarzaniaArchiwum znakuje pakowanie, które ruszyło i się nie udało podczas przetwarzania archiwum.
func bladPrzetwarzaniaArchiwum(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"narzędzia archiwum: "+powod))
}

// bladZapleczaArchiwum znakuje awarię po stronie rdzenia: magazyn niewpięty,
// nośnik pełny, wiersz nie do zapisania, blob zniknął spod odwołania.
func bladZapleczaArchiwum(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"narzędzia archiwum: "+powod))
}
