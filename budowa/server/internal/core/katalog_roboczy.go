// Katalog roboczy modelu — miejsce, w którym powstają katalogi sesyjne
// i pliki robocze modelu.
//
// Katalog roboczy nie jest dostępem. Dostęp mówi,
// do jakich maszyn i katalogów model ma wgląd; katalog roboczy mówi, gdzie
// model zostawia własne pliki. To dwa niezależne ustawienia, więc ten moduł
// nie zna punktów dostępu ani nadań i nigdy o nie nie pyta.
//
// Stan wyjściowy: katalog powstaje automatycznie w miejscu instalacji aplikacji
// głównej, a katalogi sesyjne leżą wewnątrz — `<instalacja>/sesje/<identyfikator>/`.
// Operator nadpisuje to z okna konfiguracji dwoma
// ustawieniami rozstrzyganymi po ośmiu poziomach zasięgu.
//
// Ten plik trzyma wyłącznie regułę składania ścieżki — funkcje czyste, bez
// dotknięcia dysku i bez rezolwera. Ustalenie z rezolwera i degradacja leżą
// w katalog_roboczy_ustalenie.go.
package core

import (
	"os"
	"path/filepath"
	"strings"
)

// Klucze ustawień katalogu roboczego. Odpowiadają kolumnie ustawienie.klucz
// oraz liście informacyjnej shared.KnownSettingKeys (lista informacyjna, nie
// brama; katalog ustawień jest sterowany danymi).
const (
	// KluczKatalogRoboczyPodstawa wskazuje katalog, w którym powstają katalogi
	// sesyjne. Brak wartości znaczy: miejsce instalacji aplikacji głównej.
	KluczKatalogRoboczyPodstawa = "katalog.roboczy.podstawa"
	// KluczKatalogRoboczyWzorzecSesji wskazuje wzorzec nazwy katalogu sesji
	// liczony względem podstawy.
	KluczKatalogRoboczyWzorzecSesji = "katalog.roboczy.wzorzec_sesji"
)

// ZnacznikIdentyfikatoraSesji jest miejscem podstawienia identyfikatora sesji
// we wzorcu. Wzorzec bez znacznika jest traktowany jako sam katalog nadrzędny,
// a identyfikator dokłada się na jego końcu — brak znacznika nie jest błędem.
const ZnacznikIdentyfikatoraSesji = "<identyfikator>"

// WzorzecSesjiDomyslny obowiązuje przy braku ustawienia na każdym z ośmiu
// poziomów zasięgu.
const WzorzecSesjiDomyslny = "sesje/" + ZnacznikIdentyfikatoraSesji

// nazwaSesjiZastepcza jest nazwą katalogu sesji, gdy identyfikator po oczyszczeniu
// nie zostawia ani jednego znaku dopuszczalnego w nazwie pliku. Sesja bez nazwy
// dostaje katalog o nazwie zastępczej, zamiast nie dostać katalogu.
const nazwaSesjiZastepcza = "sesja"

// KatalogInstalacji zwraca miejsce instalacji aplikacji głównej: katalog pliku
// wykonywalnego procesu. Gdy ścieżki pliku wykonywalnego nie da się ustalić,
// wraca katalog bieżący, a gdy i tego nie ma — katalog bieżący w zapisie
// względnym. Żadna z tych ścieżek nie kończy się błędem: brak rozpoznania
// miejsca instalacji nie może zatrzymać startu sesji.
func KatalogInstalacji() string {
	if plik, err := os.Executable(); err == nil && strings.TrimSpace(plik) != "" {
		if rozwiazany, err := filepath.EvalSymlinks(plik); err == nil {
			plik = rozwiazany
		}
		return filepath.Dir(plik)
	}
	if katalog, err := os.Getwd(); err == nil && strings.TrimSpace(katalog) != "" {
		return katalog
	}
	return "."
}

// PodstawaLubInstalacja zwraca podstawę wskazaną przez Operatora, a przy jej
// braku — miejsce instalacji aplikacji głównej.
func PodstawaLubInstalacja(wskazana, instalacja string) string {
	if przycieta := strings.TrimSpace(wskazana); przycieta != "" {
		return filepath.Clean(przycieta)
	}
	if przycieta := strings.TrimSpace(instalacja); przycieta != "" {
		return filepath.Clean(przycieta)
	}
	return KatalogInstalacji()
}

// WzorzecLubDomyslny zwraca wzorzec wskazany przez Operatora, a przy jego braku
// wzorzec domyślny.
func WzorzecLubDomyslny(wskazany string) string {
	if przyciety := strings.TrimSpace(wskazany); przyciety != "" {
		return przyciety
	}
	return WzorzecSesjiDomyslny
}

// SciezkaSesji składa ścieżkę katalogu jednej sesji z podstawy, wzorca
// i identyfikatora. Funkcja jest czysta — nie czyta ustawień i nie dotyka dysku.
//
// Wzorzec liczy się względem podstawy. Wzorzec bez znacznika identyfikatora
// dostaje identyfikator na końcu, więc `pliki` daje `<podstawa>/pliki/<id>`.
// Wzorzec, który wyprowadzałby poza podstawę (`..`, ścieżka bezwzględna),
// jest odrzucany na rzecz wzorca domyślnego: ustawienie Operatora steruje
// układem katalogów wewnątrz podstawy, nie omija samej podstawy.
func SciezkaSesji(podstawa, wzorzec, identyfikator string) string {
	korzen := PodstawaLubInstalacja(podstawa, "")
	wzgledna := wzglednaSciezkaSesji(WzorzecLubDomyslny(wzorzec), identyfikator)
	if wzgledna == "" {
		wzgledna = wzglednaSciezkaSesji(WzorzecSesjiDomyslny, identyfikator)
	}
	return filepath.Join(korzen, wzgledna)
}

// wzglednaSciezkaSesji podstawia identyfikator we wzorcu i sprowadza wynik do
// ścieżki względnej mieszczącej się w podstawie. Wynik pusty znaczy, że wzorzec
// wyprowadzał poza podstawę.
func wzglednaSciezkaSesji(wzorzec, identyfikator string) string {
	nazwa := nazwaKataloguSesji(identyfikator)
	ujednolicony := strings.ReplaceAll(wzorzec, "\\", "/")
	podstawiony := strings.TrimRight(ujednolicony, "/") + "/" + nazwa
	if strings.Contains(ujednolicony, ZnacznikIdentyfikatoraSesji) {
		podstawiony = strings.ReplaceAll(ujednolicony, ZnacznikIdentyfikatoraSesji, nazwa)
	}
	if strings.HasPrefix(podstawiony, "/") || filepath.IsAbs(filepath.FromSlash(podstawiony)) {
		return ""
	}
	czesci := make([]string, 0, 4)
	for _, czesc := range strings.Split(podstawiony, "/") {
		if czesc == "" || czesc == "." {
			continue
		}
		czesci = append(czesci, czesc)
	}
	if len(czesci) == 0 {
		return ""
	}
	zlozona := filepath.Join(czesci...)
	if zlozona == ".." || strings.HasPrefix(zlozona, ".."+string(filepath.Separator)) {
		return ""
	}
	return zlozona
}

// nazwaKataloguSesji zamienia identyfikator sesji na nazwę katalogu bezpieczną
// dla systemu plików: litery, cyfry, kreska, podkreślenie i kropka zostają,
// wszystko pozostałe staje się kreską. Identyfikator, z którego nie zostaje ani
// jeden znak, dostaje nazwę zastępczą.
func nazwaKataloguSesji(identyfikator string) string {
	budowana := strings.Builder{}
	for _, znak := range strings.TrimSpace(identyfikator) {
		switch {
		case znak >= 'a' && znak <= 'z', znak >= 'A' && znak <= 'Z',
			znak >= '0' && znak <= '9', znak == '-', znak == '_':
			budowana.WriteRune(znak)
		case znak == '.':
			budowana.WriteRune('.')
		default:
			budowana.WriteRune('-')
		}
	}
	nazwa := strings.Trim(strings.Trim(budowana.String(), "-"), ".")
	if nazwa == "" {
		return nazwaSesjiZastepcza
	}
	return nazwa
}
