// Odpowiedzialność pliku: formaty, ścieżki i odmowy rodziny `archive.*` —
// rozstrzygnięcia wspólne dla pakowania i rozpakowania, wyjęte z trzonu
// (`adapter_narzedzia_archiwum.go`).
//
// ── format rozpoznajemy z bajtów, a nie z nazwy ────────────────────────────
// Blob magazynu nazywa się swoją sumą kontrolną i nie ma rozszerzenia. Gdyby
// format brać z nazwy, `archive.unpack` na własnym wytworze `archive.pack`
// musiałby albo odmówić („zasób nie niesie formatu"), albo zaufać polu `format`
// wiersza — czyli etykiecie, którą ktoś kiedyś wpisał, a nie zawartości pliku.
// `7z l` rozpoznaje zip i 7z po sygnaturze bez względu na nazwę, więc etykieta
// jest tu zbędna. Czytamy sygnaturę: to wiedza o pliku, nie o jego opisie.
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

// formatArchiwum nazywa jeden z trzech formatów, które ta rodzina obsługuje.
type formatArchiwum string

const (
	formatZip   formatArchiwum = "zip"
	format7z    formatArchiwum = "7z"
	formatTarGz formatArchiwum = "tar.gz"
)

// rozpoznajFormatArchiwum czyta sygnaturę pliku i nazywa format.
//
// Sygnatury są krótkie i jednoznaczne: `PK\x03\x04` otwiera zip, `7z¼¯'\x1c`
// otwiera 7z, `\x1f\x8b` otwiera strumień gzip (a więc i tar.gz). Plik krótszy
// niż sygnatura albo o sygnaturze nieznanej jest odmową, a nie domysłem „to
// pewnie zip": rozpakowywanie czegoś, czego nie rozpoznaliśmy, kończy się
// komunikatem `7z` w obcym języku zamiast zdaniem o tym, co jest nie tak.
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

// rozstrzygnijFormatDocelowy przekłada pole `format` żądania na jeden z trzech
// formatów. Brak pola bierze zip — tak mówi kontrakt i tak jest najrozsądniej:
// zip otwiera się dwukrotnym kliknięciem w każdym systemie, więc Operator
// dostający archiwum nie potrzebuje niczego doinstalowywać.
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

// nazwaBezKatalogow wycina samą nazwę pliku ze ścieżki.
func nazwaBezKatalogow(sciezka string) string {
	sciezka = strings.TrimRight(sciezka, "/\\")
	if i := strings.LastIndexAny(sciezka, "/\\"); i >= 0 {
		return sciezka[i+1:]
	}
	return sciezka
}

// sciezkaWzgledemKatalogu rozstrzyga ścieżkę żądania względem katalogu okna
// i pilnuje, żeby wynik z tego katalogu nie wyszedł.
//
// Sprawdzenia są dwa, bo są dwa sposoby ucieczki. Ścieżka bezwzględna („/etc",
// „C:\Windows") omija katalog okna wprost. Ścieżka względna wychodzi z niego
// członem `..`, i tego nie widać po samym napisie — dlatego liczymy ścieżkę
// oczyszczoną (`filepath.Rel` po `filepath.Clean`) i pytamy, czy nadal leży
// wewnątrz. Napis „a/../../b" wygląda niewinnie i wskazuje piętro wyżej.
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

// bladArsenaluArchiwum przekłada odmowę pakietu `zewnetrzne` na kod kontraktu.
//
// Brak binarium dostaje inny kod niż niepowodzenie programu, i to jest cała
// istota typu `*BrakNarzedzia`: „nie ma czym" Operator usuwa jedną instalacją,
// a „program się wywrócił" jest usterką przetwarzania.
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

// bladWskazaniaArchiwum nazywa brak albo niepoprawność wskazania w żądaniu —
// usterka wołającego, nie rdzenia. Tym kodem idą też obie granice
// bezpieczeństwa: archiwum wychodzące poza katalog i archiwum przekraczające
// granicę rozmiaru są odmową wobec materiału, a nie awarią rdzenia.
func bladWskazaniaArchiwum(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"narzędzia archiwum: "+powod))
}

// bladPrzetwarzaniaArchiwum znakuje pakowanie, które ruszyło i się nie udało.
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
