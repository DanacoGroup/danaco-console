// Odpowiedzialność pliku: jedna droga modułu Studio do bajtów — odczyt
// materiału z magazynu zasobów rdzenia i odłożenie w nim wyniku.
package core

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// przedrostki bytów wydawanych przez moduł. Stoją tutaj, bo wszystkie trzy
// rodziny wydania sięgają po ten sam magazyn.
const (
	przedrostekGaleziStudia    = "studio-gal-"
	przedrostekOdwolaniaStudia = "studio-odw-"
	przedrostekWsaduStudia     = "studio-wsad-"
)

// bladMagazynuStudia nazywa brak magazynu wprost i mówi, co podpiąć.
//
// Rdzeń złożony bez magazynu zasobów nie jest usterką Operatora, tylko
// niepełnym montażem — odmowa ma to powiedzieć, żeby wykonawca nie szukał
// przyczyny w danych.
func bladMagazynuStudia(czynnosc string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"moduł Studio: "+czynnosc+" nie ma gdzie odłożyć wyniku — rdzeń złożony bez "+
			"magazynu zasobów; naprawa: podpiąć repozytorium zasobów i katalog danych "+
			"przy składaniu rdzenia"))
}

// bajtyZasobuStudia oddaje treść zasobu magazynu rdzenia po jego kodzie,
// czytając plik ze ścieżki znalezionej w wierszu zasobu.
func (a *adapterStudia) bajtyZasobuStudia(ctx context.Context, kod string) ([]byte, error) {
	sciezka, err := a.sciezkaZasobu(ctx, kod)
	if err != nil {
		return nil, err
	}
	bajty, err := os.ReadFile(sciezka)
	if err != nil {
		return nil, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"moduł Studio: nie można odczytać treści zasobu "+kod+": "+err.Error()))
	}
	return bajty, nil
}

// odlozTrescStudia utrwala bajty wyniku i zakłada wiersz zasobu. Okno puste
// znaczy, że wynik do żadnego okna nie należy — to nie jest usterka.
func (a *adapterStudia) odlozTrescStudia(ctx context.Context, bajty []byte,
	nazwa, format, okno string) (shared.DesignAsset, error) {

	if a.magazyn == nil {
		return shared.DesignAsset{}, bladMagazynuStudia("odłożenie zasobu " + nazwa)
	}
	skrot := sha256.Sum256(bajty)
	suma := hex.EncodeToString(skrot[:])
	odwolanie, err := a.magazyn.Zapisz(bajty, suma)
	if err != nil {
		return shared.DesignAsset{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeInternalError, "moduł Studio: "+err.Error()))
	}
	return odlozWynikArsenalu(ctx, a.zasoby, wynikArsenalu{
		odwolanie: odwolanie, okno: strings.TrimSpace(okno),
		nazwa: nazwa, format: format,
	}, func(powod string) error {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"moduł Studio: "+powod))
	})
}

// odlozObrazStudia utrwala wyrys strony wraz z wymiarami — Preview Window
// skaluje kafelek po nich, więc wymiary zmierzone przy wyrysie idą do wiersza
// zamiast być liczone drugi raz przy odczycie.
func (a *adapterStudia) odlozObrazStudia(ctx context.Context, bajty []byte,
	nazwa, okno string, szerokosc, wysokosc int) (shared.DesignAsset, error) {

	if a.magazyn == nil {
		return shared.DesignAsset{}, bladMagazynuStudia("odłożenie wyrysu " + nazwa)
	}
	skrot := sha256.Sum256(bajty)
	suma := hex.EncodeToString(skrot[:])
	odwolanie, err := a.magazyn.Zapisz(bajty, suma)
	if err != nil {
		return shared.DesignAsset{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeInternalError, "moduł Studio: "+err.Error()))
	}
	return odlozWynikArsenalu(ctx, a.zasoby, wynikArsenalu{
		odwolanie: odwolanie, okno: strings.TrimSpace(okno), nazwa: nazwa, format: "png",
		szerokosc: &szerokosc, wysokosc: &wysokosc,
	}, func(powod string) error {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"moduł Studio: "+powod))
	})
}

// trescWersjiStudia oddaje treść wskazanej wersji albo, gdy wersji nie
// podano, treść bieżącą dokumentu. Treść obszerna doczytuje się z magazynu.
func (a *adapterStudia) trescWersjiStudia(ctx context.Context, kodWersji *string,
	dokument dane.DokumentStudia) (string, error) {

	if kodWersji != nil && strings.TrimSpace(*kodWersji) != "" {
		wersja, err := a.repozytorium.Wersja(ctx, strings.TrimSpace(*kodWersji))
		if err != nil {
			return "", bladWskazaniaStronyPorownania(*kodWersji, err)
		}
		return a.trescZOdwolania(wersja.Tresc, wersja.TrescOdwolanie)
	}
	return a.trescZOdwolania(dokument.Tresc, dokument.TrescOdwolanie)
}

// trescZOdwolania rozstrzyga, skąd wziąć treść: z kolumny wprost albo z pliku
// magazynu wskazanego odwołaniem.
func (a *adapterStudia) trescZOdwolania(tresc, odwolanie *string) (string, error) {
	if tresc != nil && *tresc != "" {
		return *tresc, nil
	}
	if odwolanie != nil && strings.TrimSpace(*odwolanie) != "" {
		bajty, err := os.ReadFile(*odwolanie)
		if err != nil {
			return "", protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
				"moduł Studio: treść leży poza bazą (odwołanie "+*odwolanie+
					"), a pliku nie da się odczytać: "+err.Error()))
		}
		return string(bajty), nil
	}
	if tresc != nil {
		return *tresc, nil
	}
	return "", nil
}
