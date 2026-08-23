package core

import (
	"context"
	"os/exec"
	"strings"
	"testing"

	"github.com/tdewolff/canvas"

	"danacoconsole/shared"
)

// Skutek odczytu pisma ze zrzutu — `design.mockup.import` z polem
// `recognizeText`.
//
// ── Co ten plik mierzy ──────────────────────────────────────────────────────
// Do tej tury rdzeń rozpoznawał, KTÓRE obszary zrzutu są liniami tekstu, ale
// treści napisów nie czytał. Sprawdzian kompilacji przechodził nad tym zielono,
// bo komenda oddawała ramkę i warstwy — brakowało jedynie tego, po co Operator
// wciąga zrzut: żeby zobaczyć, co na nim NAPISANO.
//
// Odczyt idzie programem pakietu serwera, więc sprawdzian rozgałęzia się wedle
// tego, czy program na maszynie stoi — i ŻADNA z gałęzi nie jest pominięciem:
//
//   - program STOI — mierzony jest odczyt: adnotacja warstwy ma nieść treść
//     napisu, który sprawdzian sam wpisał w zrzut;
//   - programu NIE MA — mierzone są dwie odmowy: żądanie z odczytem wskazanym
//     WPROST dostaje odmowę nazwaną wraz z naprawą, a żądanie bez wskazania
//     dostaje układ obszarów, w którym każda linia tekstu mówi w adnotacji, że
//     treści nie odczytano i dlaczego.
//
// Dzięki temu ten plik świeci zielono na serwerze z pakietem i na maszynie bez
// niego, a w obu przypadkach mierzy zachowanie, nie samą kompilację.

// czyStoiCzytnikPismaSprawdzianu mówi, czy program rozpoznający pismo jest na
// maszynie. Sprawdzian wolno o to zapytać wprost — zapora obszaru Design pilnuje
// plików rdzenia, nie sprawdzianów, a rozgałęzienie bez tego pomiaru musiałoby
// zgadywać, którą odpowiedź uznać za poprawną.
func czyStoiCzytnikPismaSprawdzianu() bool {
	_, err := exec.LookPath("tesseract")
	return err == nil
}

// zrzutZNapisemPNG składa zrzut ekranu z JEDNYM napisem na białym tle.
//
// Napis jest rysowany krojem WKOMPILOWANYM, tym samym, którym rdzeń podpisuje
// wykresy — więc sprawdzian nie zależy od krojów zainstalowanych na maszynie.
// Wysokość napisu i szerokość paska są dobrane tak, żeby obszar spełnił warunek
// linii tekstu rdzenia (wysokość do sześciu kratek, szerokość co najmniej
// trzykrotność wysokości).
func zrzutZNapisemPNG(t *testing.T, napis string) []byte {
	t.Helper()

	const szerokosc, wysokosc = 360.0, 120.0
	plotno := canvas.New(szerokosc, wysokosc)
	kontekst := canvas.NewContext(plotno)
	kontekst.RenderPath(canvas.Rectangle(szerokosc, wysokosc),
		stylWypelnieniaDanychDesignu("#ffffff"), canvas.Identity)
	// Rozmiar pisma 20 jest zmierzony, nie dobrany na oko: przy nim wyraz schodzi
	// na obszar o wysokości 32 punktów — wewnątrz granicy linii tekstu (do sześciu
	// kratek) i o szerokości grubo ponad trzykrotność wysokości. Pismo większe
	// rozsypuje się na kratce po jednej literze na obszar i żadna nie jest już
	// linią tekstu.
	if err := napisWyrysuDanychDesignu(kontekst, napis, 24, 56, 20, "#101010"); err != nil {
		t.Fatalf("nie można narysować napisu sprawdzianu: %v", err)
	}
	bajty, _, err := wydajPlotnoDanychDesignu(plotno, szerokosc, wysokosc, "png")
	if err != nil {
		t.Fatalf("nie można wydać zrzutu sprawdzianu: %v", err)
	}
	return bajty
}

// wciagnijZrzutSprawdzianu wnosi zrzut do magazynu i zakłada kompozycję, w której
// makieta ma stanąć.
func wciagnijZrzutSprawdzianu(t *testing.T, zmontowany *Zmontowany, zycie context.Context,
	okno, napis string) (string, string) {

	t.Helper()

	zasob := wniesObrazSprawdzianu(t, zmontowany, zycie, okno, zrzutZNapisemPNG(t, napis))
	var plansza shared.DesignBoardUpdateResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignBoardUpdate,
		shared.DesignBoardUpdateRequest{WindowId: okno, Name: wskaznik("makieta ze zrzutu")},
		&plansza)
	return zasob.Id, plansza.Board.Id
}

// adnotacjeLiniiTekstuSprawdzianu zbiera adnotacje warstw, które rdzeń uznał za
// linie tekstu.
func adnotacjeLiniiTekstuSprawdzianu(warstwy []shared.DesignBoardLayer) []string {
	adnotacje := []string{}
	for _, warstwa := range warstwy {
		if warstwa.Note == nil {
			continue
		}
		if strings.HasPrefix(*warstwa.Note, "linia tekstu") {
			adnotacje = append(adnotacje, *warstwa.Note)
		}
	}
	return adnotacje
}

// TestWciagnietyZrzutCzytaTrescNapisow mierzy odczyt treści napisów: adnotacja
// warstwy ma nieść słowo, które sprawdzian sam wpisał w zrzut. Gałąź bez programu
// mierzy dwie odmowy — nagłówek pliku mówi, dlaczego obie są pomiarem.
func TestWciagnietyZrzutCzytaTrescNapisow(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	// Napis wielkimi literami i bez znaków spoza alfabetu łacińskiego: sprawdzian
	// mierzy DROGĘ odczytu, nie skuteczność czytnika na piśmie ozdobnym.
	const napis = "DANACO KONSOLA"
	zasob, plansza := wciagnijZrzutSprawdzianu(t, zmontowany, zycie, "okno-zrzutu", napis)

	if !czyStoiCzytnikPismaSprawdzianu() {
		// Odczyt wskazany WPROST — odmowa nazwana wraz z naprawą.
		odmowa := wykonajOdmowna(t, zmontowany, zycie, shared.CommandDesignMockupImport,
			shared.DesignMockupImportRequest{
				WindowId: "okno-zrzutu", BoardId: plansza, AssetId: zasob,
				RecognizeText: wskaznik(true),
			})
		if odmowa.Code != shared.ErrorCodeChannelUnavailable {
			t.Errorf("brak programu rozpoznającego pismo dał kod %s, a jest zapleczem "+
				"niedostępnym (ponawialnym), nie wadą żądania", odmowa.Code)
		}
		if !strings.Contains(odmowa.Message, "recognizeText") {
			t.Errorf("odmowa nie mówi, czym Operator ma ją obejść: %s", odmowa.Message)
		}

		// Odczyt DOMYŚLNY — układ powstaje, a brak odczytu jest nazwany w adnotacji.
		var wynik shared.DesignMockupImportResponse
		wykonajUdana(t, zmontowany, zycie, shared.CommandDesignMockupImport,
			shared.DesignMockupImportRequest{
				WindowId: "okno-zrzutu", BoardId: plansza, AssetId: zasob,
			}, &wynik)
		adnotacje := adnotacjeLiniiTekstuSprawdzianu(wynik.Layers)
		if len(adnotacje) == 0 {
			t.Fatalf("ze zrzutu z napisem nie wyszła ani jedna linia tekstu (warstw: %d)",
				len(wynik.Layers))
		}
		for _, adnotacja := range adnotacje {
			if !strings.Contains(adnotacja, "nie odczytano") {
				t.Errorf("adnotacja %q milczy o tym, że treści nie odczytano — Operator odczyta "+
					"to jako linię tekstu bez napisu", adnotacja)
			}
		}
		return
	}

	var wynik shared.DesignMockupImportResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignMockupImport,
		shared.DesignMockupImportRequest{
			WindowId: "okno-zrzutu", BoardId: plansza, AssetId: zasob,
			RecognizeText: wskaznik(true),
		}, &wynik)

	adnotacje := adnotacjeLiniiTekstuSprawdzianu(wynik.Layers)
	if len(adnotacje) == 0 {
		t.Fatalf("ze zrzutu z napisem %q nie wyszła ani jedna linia tekstu (warstw: %d)",
			napis, len(wynik.Layers))
	}
	// Odczyt ma stać W ADNOTACJI, bo tam Operator go widzi. Miarą jest SŁOWO ze
	// zrzutu, nie sama obecność dwukropka: adnotacja „linia tekstu 1: " byłaby
	// odczytem pustym udającym odczyt.
	odczytane := false
	for _, adnotacja := range adnotacje {
		if strings.Contains(strings.ToUpper(adnotacja), "DANACO") {
			odczytane = true
		}
		if strings.Contains(adnotacja, "nie odczytano") {
			t.Errorf("program rozpoznający pismo stoi na tej maszynie, a adnotacja mówi, "+
				"że treści nie odczytano: %q", adnotacja)
		}
	}
	if !odczytane {
		t.Errorf("żadna adnotacja linii tekstu nie niesie słowa ze zrzutu (%q); adnotacje: %v",
			napis, adnotacje)
	}

	// Odczyt WYŁĄCZONY wprost nie ma prawa czytać: pole `recognizeText=false`
	// znaczy „nie czytaj", a nie „czytaj i nie mów".
	var bezOdczytu shared.DesignMockupImportResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignMockupImport,
		shared.DesignMockupImportRequest{
			WindowId: "okno-zrzutu", BoardId: plansza, AssetId: zasob,
			RecognizeText: wskaznik(false),
		}, &bezOdczytu)
	for _, warstwa := range bezOdczytu.Layers {
		if warstwa.Note == nil {
			continue
		}
		if strings.Contains(strings.ToUpper(*warstwa.Note), "DANACO") {
			t.Errorf("żądanie z recognizeText=false oddało odczytaną treść w adnotacji %q",
				*warstwa.Note)
		}
	}
}
