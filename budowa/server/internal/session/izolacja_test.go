package session

import (
	"errors"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"danacoconsole/server/internal/konfig"
)

// Przynależność ścieżki do obszaru okna.
//
// SciezkaWewnatrz jest zaporą, przez którą przechodzi każde sięgnięcie rdzenia
// po plik w imieniu okna, oraz każdy korzeń podawany procesowi modelu przy
// uruchomieniu. Jej pomyłka nie zgłasza się błędem — otwiera katalog, którego
// Operator nie dał.
//
// Sprawdzian pisany jest tak, by mierzył to samo na obu systemach, bo produkt
// jedzie i na Linuksa, i na Windowsa. Jedyna różnica jest zamierzona i dotyczy
// wielkości liter: system plików Windowsa jej nie rozróżnia, więc porównanie
// też nie może — inaczej ta sama ścieżka zapisana inną wielkością omijałaby
// obszar. Ta różnica ma tu własny sprawdzian, a nie założenie milczące.

// TestSciezkaWewnatrzPrzepuszczaKorzenIJegoWnetrze mierzy drogę zwykłą.
func TestSciezkaWewnatrzPrzepuszczaKorzenIJegoWnetrze(t *testing.T) {
	korzen := t.TempDir()

	przypadki := []string{
		korzen,
		filepath.Join(korzen, "plik.txt"),
		filepath.Join(korzen, "podkatalog"),
		filepath.Join(korzen, "podkatalog", "glebiej", "plik.txt"),
		// Człon `.` znika przy oczyszczaniu i nie wyprowadza ścieżki z obszaru.
		filepath.Join(korzen, ".", "plik.txt"),
		// Zejście i powrót zostaje wewnątrz.
		filepath.Join(korzen, "podkatalog", "..", "plik.txt"),
	}
	for _, sciezka := range przypadki {
		t.Run(sciezka, func(t *testing.T) {
			if !SciezkaWewnatrz(korzen, sciezka) {
				t.Errorf("ścieżka %q uznana za leżącą poza korzeniem %q", sciezka, korzen)
			}
		})
	}
}

// TestSciezkaWewnatrzOdrzucaWyjsciePozaKorzen mierzy zaporę. Wyjście w górę
// drzewa jest tu przypadkiem właściwym, nie egzotycznym: ścieżka składana
// z członów podanych przez model dochodzi w takiej postaci.
func TestSciezkaWewnatrzOdrzucaWyjsciePozaKorzen(t *testing.T) {
	nadrzedny := t.TempDir()
	korzen := filepath.Join(nadrzedny, "obszar")

	przypadki := map[string]string{
		"katalog nadrzędny":            nadrzedny,
		"wyjście raz w górę":           filepath.Join(korzen, "..", "cudze.txt"),
		"wyjście dwa razy w górę":      filepath.Join(korzen, "..", "..", "cudze.txt"),
		"rodzeństwo o podobnej nazwie": filepath.Join(nadrzedny, "obszar-cudzy", "plik.txt"),
		"rodzeństwo o nazwie dłuższej": filepath.Join(nadrzedny, "obszarowy", "plik.txt"),
	}
	for nazwa, sciezka := range przypadki {
		t.Run(nazwa, func(t *testing.T) {
			if SciezkaWewnatrz(korzen, sciezka) {
				t.Errorf("ścieżka %q uznana za leżącą w korzeniu %q", sciezka, korzen)
			}
		})
	}
}

// TestPustyKorzenAlboPustaSciezkaNiczegoNieOtwiera pilnuje wartości domyślnej.
// Pusty korzeń nie może znaczyć „wszystko wolno" — a to jest dokładnie ten
// przypadek, który powstaje przy niewypełnionym polu formularza.
func TestPustyKorzenAlboPustaSciezkaNiczegoNieOtwiera(t *testing.T) {
	korzen := t.TempDir()

	przypadki := []struct{ korzen, sciezka string }{
		{"", filepath.Join(korzen, "plik.txt")},
		{"   ", filepath.Join(korzen, "plik.txt")},
		{korzen, ""},
		{korzen, "   "},
		{"", ""},
	}
	for _, przypadek := range przypadki {
		if SciezkaWewnatrz(przypadek.korzen, przypadek.sciezka) {
			t.Errorf("puste wskazanie przepuściło: korzeń %q, ścieżka %q",
				przypadek.korzen, przypadek.sciezka)
		}
	}
}

// TestOdstepyWokolSciezkiNieZmieniajaRozstrzygniecia sprawdza przycinanie.
// Ścieżka z odstępem na początku przychodzi z pola formularza i z argumentu
// podanego przez model.
func TestOdstepyWokolSciezkiNieZmieniajaRozstrzygniecia(t *testing.T) {
	korzen := t.TempDir()
	plik := filepath.Join(korzen, "plik.txt")

	if !SciezkaWewnatrz("  "+korzen+"  ", "  "+plik+"  ") {
		t.Error("odstępy wokół ścieżki wyprowadziły ją poza obszar")
	}
}

// TestSciezkaWzgledaJestRozstrzyganaWzgledemKataloguBiezacego utrwala
// zachowanie normalizacji: ścieżka względna staje się bezwzględna, zanim
// cokolwiek zostanie porównane. Bez tego ta sama ścieżka raz mieściłaby się
// w obszarze, a raz nie — zależnie od tego, skąd wołano.
func TestSciezkaWzgledaJestRozstrzyganaWzgledemKataloguBiezacego(t *testing.T) {
	katalog := t.TempDir()
	t.Chdir(katalog)

	if !SciezkaWewnatrz(katalog, "plik.txt") {
		t.Error("ścieżka względna w katalogu bieżącym uznana za leżącą poza obszarem")
	}
	if !SciezkaWewnatrz(".", "plik.txt") {
		t.Error("korzeń wskazany kropką nie objął pliku obok")
	}
	if SciezkaWewnatrz(katalog, filepath.Join("..", "cudze.txt")) {
		t.Error("ścieżka względna wychodząca w górę uznana za leżącą w obszarze")
	}
}

// TestWielkoscLiterRozstrzygaSystemPlikow jest jedynym sprawdzianem w tym pliku,
// który daje różne wyniki na różnych systemach — i taka jest reguła produktu.
//
// Windows nie rozróżnia wielkości liter w nazwach ścieżek, więc porównanie też
// nie może: `C:\Dane` i `c:\dane` to ten sam katalog i ścieżka zapisana inną
// wielkością nie może omijać obszaru. Systemy plików Linuksa rozróżniają, więc
// `dane` i `Dane` są dwoma różnymi katalogami i sklejenie ich otwierałoby
// obszar, którego Operator nie wskazał.
//
// Gałąź windowsowa normalizacji nie wykonuje się nigdy na maszynie budującej —
// dopiero zadanie Windows w bramce sprawdzianów ją uruchamia. Ten sprawdzian
// jest po to, żeby miała co uruchomić.
func TestWielkoscLiterRozstrzygaSystemPlikow(t *testing.T) {
	nadrzedny := t.TempDir()
	korzen := filepath.Join(nadrzedny, "obszar")
	innaWielkosc := filepath.Join(nadrzedny, "OBSZAR", "plik.txt")

	wewnatrz := SciezkaWewnatrz(korzen, innaWielkosc)

	if runtime.GOOS == "windows" {
		if !wewnatrz {
			t.Error("na Windowsie ścieżka zapisana inną wielkością liter ominęła obszar")
		}
		return
	}
	if wewnatrz {
		t.Error("na systemie rozróżniającym wielkość liter dwa różne katalogi zostały sklejone")
	}
}

// TestNaruszenieWiazeSieZeWspolnymKorzeniem sprawdza rozpoznawanie przyczyny.
// Warstwa wyżej pyta przez errors.Is, nie przez treść komunikatu — treść jest
// dla Operatora, korzeń dla kodu.
func TestNaruszenieWiazeSieZeWspolnymKorzeniem(t *testing.T) {
	err := NoweNaruszenie(konfig.KluczIzolacjaPliki, "ścieżka poza obszarem okna")

	if !errors.Is(err, ErrIzolacja) {
		t.Error("naruszenie nie wiąże się ze wspólnym korzeniem")
	}
	if !strings.Contains(err.Error(), konfig.KluczIzolacjaPliki) {
		t.Errorf("komunikat nie nazywa punktu izolacji: %v", err)
	}
	if !strings.Contains(err.Error(), "ścieżka poza obszarem okna") {
		t.Errorf("komunikat nie niesie powodu: %v", err)
	}
	if errors.Is(errors.New("błąd obcy"), ErrIzolacja) {
		t.Error("błąd obcy rozpoznany jako naruszenie izolacji")
	}
}

// politykaZWartosciami składa politykę efektywną z par klucz → wartość.
func politykaZWartosciami(wartosci map[string]string) konfig.Polityka {
	pozycje := make([]konfig.Wynik, 0, len(wartosci))
	for klucz, wartosc := range wartosci {
		pozycje = append(pozycje, konfig.Wynik{
			Klucz: klucz, Wartosc: wartosc, Rodzaj: konfig.RodzajTekst,
			Pochodzenie: konfig.PochodzenieZapis,
		})
	}
	return konfig.Polityka{Pozycje: pozycje}
}

// TestPolitykaPustaDajeStanWyjsciowyPlatformy pilnuje wartości domyślnej
// jedenastu punktów izolacji. Klucz nierozstrzygnięty ma dawać kontekst
// odrębny i żaden zakres techniczny niewłączony — czyli pełną swobodę
// operacyjną przy odrębnych wymiarach kontekstu.
func TestPolitykaPustaDajeStanWyjsciowyPlatformy(t *testing.T) {
	zasady := ZasadyZPolityki(konfig.Polityka{})

	if !zasady.HistoriaOdrebna || !zasady.PamiecOdrebna || !zasady.KontekstOdrebny {
		t.Errorf("polityka pusta zeszła ze stanu odrębnego: %+v", zasady)
	}
	for nazwa, wlaczony := range map[string]bool{
		"katalog roboczy":    zasady.KatalogRoboczy,
		"środowisko procesu": zasady.SrodowiskoProcesu,
		"katalog danych":     zasady.KatalogDanychModelu,
		"dostęp sieciowy":    zasady.DostepSieciowy,
		"pliki":              zasady.Pliki,
		"konto i token":      zasady.KontoIToken,
		"model procesu":      zasady.ModelProcesu,
		"serwer wykonania":   zasady.SerwerWykonania,
	} {
		if wlaczony {
			t.Errorf("polityka pusta włączyła zakres %q", nazwa)
		}
	}
}

// TestWartoscSpozaSlownikaNieWlaczaZakresu sprawdza regułę fail-open zakresów
// technicznych: włączenie jest decyzją zapisaną wprost, a wartość nierozpoznana
// zostawia zakres wyłączony.
func TestWartoscSpozaSlownikaNieWlaczaZakresu(t *testing.T) {
	for _, wartosc := range []string{"", "   ", "tak", "wlaczona", "1", "true"} {
		zasady := ZasadyZPolityki(politykaZWartosciami(map[string]string{
			konfig.KluczIzolacjaPliki: wartosc,
		}))
		if zasady.Pliki {
			t.Errorf("wartość %q spoza słownika włączyła zakres plikowy", wartosc)
		}
	}

	zasady := ZasadyZPolityki(politykaZWartosciami(map[string]string{
		konfig.KluczIzolacjaPliki: konfig.IzolacjaWlaczona,
	}))
	if !zasady.Pliki {
		t.Errorf("wartość %q ze słownika nie włączyła zakresu plikowego", konfig.IzolacjaWlaczona)
	}
}

// TestWspoldzielenieWymiaruJestDecyzjaZapisanaWprost sprawdza regułę odwrotną
// dla trzech wymiarów kontekstu: odrębność jest stanem wyjściowym, a zejście
// z niej wymaga wartości wskazanej wprost. Każda inna wartość zostawia wymiar
// odrębny — pomyłka w zapisie nie może rozszczelnić kontekstu.
func TestWspoldzielenieWymiaruJestDecyzjaZapisanaWprost(t *testing.T) {
	wspoldzielona := ZasadyZPolityki(politykaZWartosciami(map[string]string{
		konfig.KluczIzolacjaHistoria: konfig.IzolacjaWspoldzielona,
	}))
	if wspoldzielona.HistoriaOdrebna {
		t.Error("wartość wskazująca współdzielenie nie zeszła z odrębności")
	}

	for _, wartosc := range []string{"", "   ", "wspolna", "nie", "false"} {
		zasady := ZasadyZPolityki(politykaZWartosciami(map[string]string{
			konfig.KluczIzolacjaHistoria: wartosc,
		}))
		if !zasady.HistoriaOdrebna {
			t.Errorf("wartość %q spoza słownika rozszczelniła wymiar historii", wartosc)
		}
	}
}

// TestOdstepyWokolWartosciNieZmieniajaRozstrzygniecia pilnuje przycinania po
// stronie polityki — wartość wpisana z odstępem to ta sama wartość.
func TestOdstepyWokolWartosciNieZmieniajaRozstrzygniecia(t *testing.T) {
	zasady := ZasadyZPolityki(politykaZWartosciami(map[string]string{
		konfig.KluczIzolacjaPliki:    "  " + konfig.IzolacjaWlaczona + "  ",
		konfig.KluczIzolacjaHistoria: "  " + konfig.IzolacjaWspoldzielona + "  ",
	}))

	if !zasady.Pliki {
		t.Error("odstępy wokół wartości nie włączyły zakresu plikowego")
	}
	if zasady.HistoriaOdrebna {
		t.Error("odstępy wokół wartości nie zeszły z odrębności historii")
	}
}
