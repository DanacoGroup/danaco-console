package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// materialTwarzySprawdzianu to fotografia domeny publicznej z jedną twarzą.
const materialTwarzySprawdzianu = "testdata/twarz_sprawdzianu.jpg"

// Sprawdzian odmowy przebiegu twarzowego image.upscale przy braku wag, mierzonej bez ciężkich zasobów.

// TestBrakWagTwarzyOdmawiaNazywajacPlikIDrogeNaprawy pilnuje, żeby odmowa niosła nazwę brakującego pliku, ścieżkę, pod którą ma leżeć, oraz miejsce, z którego się go bierze.
func TestBrakWagTwarzyOdmawiaNazywajacPlikIDrogeNaprawy(t *testing.T) {
	pusty := t.TempDir()

	err := sprawdzWagiTwarzy(pusty)
	if err == nil {
		t.Fatal("katalog bez wag przeszedł sprawdzenie — przebieg ruszyłby i wywrócił się " +
			"dopiero we wnętrzu pomocnika, komunikatem Pythona zamiast odmową rdzenia")
	}
	blad := protocol.BladZeZrodla(shared.ErrorCodeInternalError, err)
	if blad.Code != shared.ErrorCodeChannelUnavailable {
		t.Fatalf("brak wag dostał kod %s, a brakująca instalacja jest zapleczem "+
			"niedostępnym (%s)", blad.Code, shared.ErrorCodeChannelUnavailable)
	}
	for _, oczekiwany := range []string{wagiOdtwarzaniaTwarzy, pusty, "github.com"} {
		if !strings.Contains(blad.Message, oczekiwany) {
			t.Fatalf("odmowa nie wymienia %q; treść: %s", oczekiwany, blad.Message)
		}
	}
}

// TestPomocnikTwarzyStoiWWykazieZaleznosci pilnuje, żeby silnik przebiegu twarzowego był widoczny w sondzie startowej, zanim jego brak ujawni się dopiero po wywołaniu.
func TestPomocnikTwarzyStoiWWykazieZaleznosci(t *testing.T) {
	szukany := narzedzieOdtwarzaniaTwarzy().Program
	if szukany == "" {
		t.Fatal("silnik przebiegu twarzowego nie ma nazwy programu")
	}
	for _, pozycja := range ZaleznosciZewnetrzne() {
		if pozycja.Narzedzie.Program == szukany {
			if strings.TrimSpace(pozycja.Zakres) == "" {
				t.Fatal("pozycja wykazu nie mówi, co przestaje działać przy braku")
			}
			return
		}
	}
	t.Fatalf("programu %s nie ma w wykazie zależności zewnętrznych", szukany)
}

func pominSprawdzianTwarzyBezSilnika(t *testing.T) {
	t.Helper()

	if err := sprawdzWagiTwarzy(katalogWagTwarzy()); err != nil {
		t.Skipf("wagi przebiegu twarzowego nie stoją na tej maszynie: %v", err)
	}
	if _, err := os.Stat(filepath.Join(katalogModeliPowiekszenia(), modelPowiekszeniaZdjec+".bin")); err != nil {
		t.Skipf("wagi Real-ESRGAN nie stoją na tej maszynie: %v", err)
	}
}

// Sprawdzian porównuje piksele wyniku bez i z przebiegiem twarzowym.
func TestPrzebiegTwarzowyZmieniaPikseleTwarzyNaPrawdziwymZdjeciu(t *testing.T) {
	pominSprawdzianTwarzyBezSilnika(t)

	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)

	bajty, err := os.ReadFile(materialTwarzySprawdzianu)
	if err != nil {
		t.Fatalf("nie można odczytać materiału sprawdzianu %s: %v", materialTwarzySprawdzianu, err)
	}
	zrodlo := wniesObraz(t, zmontowany, zycie, bajty, "twarz.jpg")

	var bezTwarzy shared.ImageUpscaleResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandImageUpscale,
		shared.ImageUpscaleRequest{
			WindowId: wskaznik("okno-sprawdzianu"),
			AssetId:  wskaznik(zrodlo),
			Scale:    wskaznik(2),
		}, &bezTwarzy)

	var zTwarza shared.ImageUpscaleResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandImageUpscale,
		shared.ImageUpscaleRequest{
			WindowId: wskaznik("okno-sprawdzianu"),
			AssetId:  wskaznik(zrodlo),
			Scale:    wskaznik(2),
			Faces:    wskaznik(true),
		}, &zTwarza)

	if zTwarza.Asset.Name == nil || !strings.Contains(*zTwarza.Asset.Name, "twarze poprawione: 1") {
		nazwa := ""
		if zTwarza.Asset.Name != nil {
			nazwa = *zTwarza.Asset.Name
		}
		t.Fatalf("nazwa zasobu %q nie mówi o jednej poprawionej twarzy na materiale, w którym jest dokładnie jedna", nazwa)
	}

	obrazBez := obrazZMagazynu(t, zmontowany, zycie, katalog, bezTwarzy.Asset.Id)
	obrazZ := obrazZMagazynu(t, zmontowany, zycie, katalog, zTwarza.Asset.Id)
	if obrazBez.Bounds() != obrazZ.Bounds() {
		t.Fatalf("wynik z przebiegiem twarzowym ma wymiary %v, a bez niego %v — powiększenie miało dać "+
			"te same wymiary, przebieg twarzowy tylko poprawia treść wycinka twarzy",
			obrazZ.Bounds(), obrazBez.Bounds())
	}

	najwiekszaRoznica := 0
	zmienionychPikseli := 0
	granice := obrazBez.Bounds()
	for y := granice.Min.Y; y < granice.Max.Y; y++ {
		for x := granice.Min.X; x < granice.Max.X; x++ {
			r1, g1, b1, _ := obrazBez.At(x, y).RGBA()
			r2, g2, b2, _ := obrazZ.At(x, y).RGBA()
			roznica := abs(int(r1>>8) - int(r2>>8))
			if d := abs(int(g1>>8) - int(g2>>8)); d > roznica {
				roznica = d
			}
			if d := abs(int(b1>>8) - int(b2>>8)); d > roznica {
				roznica = d
			}
			if roznica > najwiekszaRoznica {
				najwiekszaRoznica = roznica
			}
			if roznica > 15 {
				zmienionychPikseli++
			}
		}
	}
	if najwiekszaRoznica < 40 {
		t.Fatalf("największa różnica składowej koloru między wynikiem bez i z przebiegiem twarzowym "+
			"to %d — GFPGAN praktycznie nie zmienił obrazu", najwiekszaRoznica)
	}
	if zmienionychPikseli < 1000 {
		t.Fatalf("tylko %d pikseli zmieniło się zauważalnie — przebieg twarzowy miał przemalować "+
			"cały wycinek twarzy, nie garstkę pikseli", zmienionychPikseli)
	}
}
