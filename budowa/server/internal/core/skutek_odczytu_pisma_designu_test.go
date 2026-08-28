package core

import (
	"context"
	"os/exec"
	"strings"
	"testing"

	"github.com/tdewolff/canvas"

	"danacoconsole/shared"
)

// Sprawdzian mierzy odczyt treści napisów zrzutu w design.mockup.import z polem recognizeText.

// czyStoiCzytnikPismaSprawdzianu mówi, czy program rozpoznający pismo jest na maszynie, żeby sprawdzian mógł rozgałęzić się na mierzony odczyt albo mierzoną odmowę.
func czyStoiCzytnikPismaSprawdzianu() bool {
	_, err := exec.LookPath("tesseract")
	return err == nil
}

// zrzutZNapisemPNG składa zrzut ekranu z jednym napisem na białym tle, krojem wkompilowanym, w rozmiarze spełniającym warunek linii tekstu rdzenia.
func zrzutZNapisemPNG(t *testing.T, napis string) []byte {
	t.Helper()

	const szerokosc, wysokosc = 360.0, 120.0
	plotno := canvas.New(szerokosc, wysokosc)
	kontekst := canvas.NewContext(plotno)
	kontekst.RenderPath(canvas.Rectangle(szerokosc, wysokosc),
		stylWypelnieniaDanychDesignu("#ffffff"), canvas.Identity)
	// Rozmiar pisma 20 mieści się w granicy linii tekstu rdzenia; pismo większe się rozsypuje.
	if err := napisWyrysuDanychDesignu(kontekst, napis, 24, 56, 20, "#101010"); err != nil {
		t.Fatalf("nie można narysować napisu sprawdzianu: %v", err)
	}
	bajty, _, err := wydajPlotnoDanychDesignu(plotno, szerokosc, wysokosc, "png")
	if err != nil {
		t.Fatalf("nie można wydać zrzutu sprawdzianu: %v", err)
	}
	return bajty
}

// wciagnijZrzutSprawdzianu wnosi zrzut do magazynu i zakłada kompozycję, w której makieta ma stanąć, zwracając identyfikatory zasobu i planszy.
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

// adnotacjeLiniiTekstuSprawdzianu zbiera adnotacje warstw, które rdzeń uznał za linie tekstu, pomijając warstwy pozostałych rodzajów.
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

	// Napis wielkimi literami bez znaków spoza łacińskiego mierzy drogę odczytu, nie skuteczność czytnika.
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
	// Odczyt ma stać w adnotacji: miarą jest słowo ze zrzutu, nie sama obecność dwukropka.
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

	// Odczyt wyłączony wprost nie czyta: recognizeText=false znaczy nie czytaj, nie milcz i czytaj.
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
