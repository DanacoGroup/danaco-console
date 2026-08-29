// Plik obsługuje dwie czynności na sieciach neuronowych: `image.upscale`
// (Real-ESRGAN) i `image.background.remove` (rembg), wraz z opisem obu
// silników i składaniem ich wiersza poleceń, licząc wyłącznie na procesorze.
package core

import (
	"context"
	"path/filepath"
	"strconv"
	"strings"

	"danacoconsole/server/internal/zewnetrzne"
	"danacoconsole/shared"
)

// narzedziePowiekszenia opisuje binarium superrozdzielczości: nazwa czytelna
// i pakiet wchodzą do treści odmowy, żeby Operator wiedział, czego brakuje
// i skąd to wziąć.
func narzedziePowiekszenia() zewnetrzne.Narzedzie {
	return zewnetrzne.Narzedzie{
		Nazwa:   "Real-ESRGAN (ncnn)",
		Program: "realesrgan-ncnn-vulkan",
		Pakiet: "wydanie realesrgan-ncnn-vulkan z github.com/xinntao/Real-ESRGAN/releases " +
			"rozpakowane do /usr/local/bin wraz z modelami w " + katalogModeliPowiekszenia() +
			", oraz mesa-vulkan-drivers dla liczenia na procesorze",
	}
}

// narzedzieWycinaniaTla opisuje binarium wycinania tła: wołane jest
// opakowanie `/usr/local/bin/rembg`, które samo ustawia katalog wag, więc
// silnik jest samowystarczalny niezależnie od tego, kto go woła.
func narzedzieWycinaniaTla() zewnetrzne.Narzedzie {
	return zewnetrzne.Narzedzie{
		Nazwa:   "rembg (U²-Net / ONNX Runtime)",
		Program: "rembg",
		Pakiet: "rembg[cli] w osobnym środowisku pythonowym, wystawiony opakowaniem " +
			"/usr/local/bin/rembg ustawiającym U2NET_HOME=" + katalogWagWycinania(),
	}
}

// modelPowiekszeniaZdjec to jedyna sieć powiększająca, którą ta komenda woła.
// Nazwa modelu wchodzi w argument programu i dlatego jest stałą:
// `realesrgan-x4plus` jest wyborem dla zdjęć, bo sieci uczone na rysunku
// zostawiają płaskie płaszczyzny.
const modelPowiekszeniaZdjec = "realesrgan-x4plus"

// dopuszczalneKrotnosci to krotności, które silnik ncnn przyjmuje. Krotność
// spoza tego zbioru kończy się odmową wymieniającą dopuszczalne, a nie cichym
// zaokrągleniem do najbliższej.
var dopuszczalneKrotnosci = map[int]struct{}{2: {}, 3: {}, 4: {}}

// wagiModelu opisuje jeden plik wag: gdzie leży w katalogu wag wycinania,
// ile waży i skąd się go pobiera, do treści odmowy.
type wagiModelu struct {
	// plik jest nazwą pliku wag w katalogu wag wycinania.
	plik string
	// waga wchodzi do treści odmowy, żeby czytelnik wiedział, na co się pisze.
	waga string
	// skad wskazuje źródło wag — odmowa ma mówić, gdzie ich szukać.
	skad string
}

// modeleWycinaniaTla to zamknięty zbiór sieci wycinających tło wraz z wagą
// pliku i miejscem, z którego się go bierze; wszystko troje wchodzi do treści
// odmowy przy braku pliku, a wybór modelu jest zbiorem, nie dowolnym tekstem.
var modeleWycinaniaTla = map[string]wagiModelu{
	"u2net":             {plik: "u2net.onnx", waga: "176 MB", skad: "github.com/danielgatis/rembg/releases (u2net.onnx)"},
	"u2netp":            {plik: "u2netp.onnx", waga: "4,7 MB", skad: "github.com/danielgatis/rembg/releases (u2netp.onnx)"},
	"u2net_human_seg":   {plik: "u2net_human_seg.onnx", waga: "176 MB", skad: "github.com/danielgatis/rembg/releases (u2net_human_seg.onnx)"},
	"isnet-general-use": {plik: "isnet-general-use.onnx", waga: "179 MB", skad: "github.com/danielgatis/rembg/releases (isnet-general-use.onnx)"},
	"silueta":           {plik: "silueta.onnx", waga: "44 MB", skad: "github.com/danielgatis/rembg/releases (silueta.onnx)"},
}

// domyslnyModelWycinania jest brany, gdy żądanie nie wskazuje modelu,
// kontrakt obiecuje brak bierze domyślny silnika, a domyślnym silnika rembg
// jest właśnie u2net.
const domyslnyModelWycinania = "u2net"

// Powieksz obsługuje `image.upscale`, powiększenie z odtworzeniem szczegółu.
// Kolejność kroków jest zamierzona: najpierw żądanie, potem źródło, potem
// wagi, a dopiero na końcu rusza silnik, bo przebieg silnika kosztuje minuty.
func (a *adapterNarzedziObrazuModelu) Powieksz(ctx context.Context,
	z shared.ImageUpscaleRequest) (shared.ImageUpscaleResponse, error) {

	krotnosc, err := rozstrzygnijKrotnosc(z.Scale)
	if err != nil {
		return shared.ImageUpscaleResponse{}, err
	}
	twarze := z.Faces != nil && *z.Faces
	if a.wspolne == nil {
		return shared.ImageUpscaleResponse{}, bladZapleczaModeluObrazu(
			"zaplecze narzędzi obrazu nie jest wpięte — nie ma czym rozwiązać źródła")
	}

	zrodlo, err := a.wspolne.rozwiazZrodlo(ctx, z.AssetId, z.SourcePath)
	if err != nil {
		return shared.ImageUpscaleResponse{}, err
	}
	if err := sprawdzWagi(filepath.Join(katalogModeliPowiekszenia(), modelPowiekszeniaZdjec+".bin"),
		modelPowiekszeniaZdjec, "64 MB",
		"wydania realesrgan-ncnn-vulkan (katalog models)"); err != nil {
		return shared.ImageUpscaleResponse{}, err
	}
	// Wagi przebiegu twarzowego są sprawdzane teraz, nie po powiększeniu liczonym minutami.
	if twarze {
		if err := sprawdzWagiTwarzy(katalogWagTwarzy()); err != nil {
			return shared.ImageUpscaleResponse{}, err
		}
	}

	pracownia, err := przygotujPracownie(zrodlo.sciezka)
	if err != nil {
		return shared.ImageUpscaleResponse{}, err
	}
	defer pracownia.sprzatnij()

	// Format wyniku png jest wymuszany jawnie, bo zapis stratny odjąłby odtworzony szczegół.
	argumenty := []string{
		"-i", pracownia.wejscie,
		"-o", pracownia.wyjscie,
		"-s", strconv.Itoa(krotnosc),
		"-n", modelPowiekszeniaZdjec,
		"-m", katalogModeliPowiekszenia(),
		"-f", "png",
	}
	if err := a.wolajSilnik(ctx, narzedziePowiekszenia(), argumenty, granicaPowiekszenia); err != nil {
		return shared.ImageUpscaleResponse{}, err
	}

	bajty, err := odczytajWynikSilnika(pracownia.wyjscie, "Real-ESRGAN")
	if err != nil {
		return shared.ImageUpscaleResponse{}, err
	}
	opis := "upscale x" + strconv.Itoa(krotnosc)
	// Odczyt wyżej rozstrzyga, czy powiększenie zostawiło obraz, zanim ruszy kosztowny przebieg twarzowy.
	if twarze {
		poprawione, ile, err := a.poprawTwarze(ctx, pracownia)
		if err != nil {
			return shared.ImageUpscaleResponse{}, err
		}
		bajty = poprawione
		opis += opisPrzebieguTwarzy(ile)
	}
	zasob, _, err := a.wspolne.odlozZasob(ctx, zrodlo, z.WindowId, bajty, "png", opis)
	if err != nil {
		return shared.ImageUpscaleResponse{}, err
	}
	szerokosc, wysokosc, err := wymiaryWyniku(zasob)
	if err != nil {
		return shared.ImageUpscaleResponse{}, err
	}
	return shared.ImageUpscaleResponse{Asset: zasob, Width: szerokosc, Height: wysokosc}, nil
}

// rozstrzygnijKrotnosc przekłada nieobowiązkowe `scale` na krotność silnika.
// Brak wskazania bierze krotność dwukrotną, tak jak obiecuje kontrakt.
func rozstrzygnijKrotnosc(zadana *int) (int, error) {
	if zadana == nil {
		return 2, nil
	}
	if _, jest := dopuszczalneKrotnosci[*zadana]; !jest {
		return 0, bladWskazaniaModeluObrazu("krotność powiększenia " +
			strconv.Itoa(*zadana) + " nie jest obsługiwana — silnik przyjmuje 2, 3 albo 4")
	}
	return *zadana, nil
}

// UsunTlo obsługuje `image.background.remove`, wycięcie obiektu z tła. Wynik
// jest zawsze w PNG, bo przezroczystość ma gdzie się zapisać wyłącznie
// w formacie z kanałem alfa.
func (a *adapterNarzedziObrazuModelu) UsunTlo(ctx context.Context,
	z shared.ImageBackgroundRemoveRequest) (shared.ImageBackgroundRemoveResponse, error) {

	nazwaModelu, opis, err := rozstrzygnijModelWycinania(z.Model)
	if err != nil {
		return shared.ImageBackgroundRemoveResponse{}, err
	}
	if a.wspolne == nil {
		return shared.ImageBackgroundRemoveResponse{}, bladZapleczaModeluObrazu(
			"zaplecze narzędzi obrazu nie jest wpięte — nie ma czym rozwiązać źródła")
	}

	zrodlo, err := a.wspolne.rozwiazZrodlo(ctx, z.AssetId, z.SourcePath)
	if err != nil {
		return shared.ImageBackgroundRemoveResponse{}, err
	}
	if err := sprawdzWagi(filepath.Join(katalogWagWycinania(), opis.plik),
		nazwaModelu, opis.waga, opis.skad); err != nil {
		return shared.ImageBackgroundRemoveResponse{}, err
	}

	pracownia, err := przygotujPracownie(zrodlo.sciezka)
	if err != nil {
		return shared.ImageBackgroundRemoveResponse{}, err
	}
	defer pracownia.sprzatnij()

	// Podpolecenie i jest trybem jeden plik na jeden plik; tryby p i s są tu niepotrzebne.
	argumenty := []string{"i", "-m", nazwaModelu, pracownia.wejscie, pracownia.wyjscie}
	if err := a.wolajSilnik(ctx, narzedzieWycinaniaTla(), argumenty, granicaWycinaniaTla); err != nil {
		return shared.ImageBackgroundRemoveResponse{}, err
	}

	bajty, err := odczytajWynikSilnika(pracownia.wyjscie, "rembg")
	if err != nil {
		return shared.ImageBackgroundRemoveResponse{}, err
	}
	zasob, _, err := a.wspolne.odlozZasob(ctx, zrodlo, z.WindowId, bajty, "png", "bez tła")
	if err != nil {
		return shared.ImageBackgroundRemoveResponse{}, err
	}
	return shared.ImageBackgroundRemoveResponse{Asset: zasob}, nil
}

// rozstrzygnijModelWycinania przekłada nieobowiązkowe `model` na pozycję
// zamkniętego zbioru sieci. Nazwa spoza zbioru kończy się odmową wymieniającą
// znane.
func rozstrzygnijModelWycinania(zadany *string) (string, wagiModelu, error) {
	nazwa := domyslnyModelWycinania
	if !bezWartosci(zadany) {
		nazwa = strings.ToLower(strings.TrimSpace(*zadany))
	}
	opis, jest := modeleWycinaniaTla[nazwa]
	if !jest {
		return "", opis, bladWskazaniaModeluObrazu("nie znam modelu wycinania „" + nazwa +
			"” — serwer zna: " + strings.Join(znaneModeleWycinania(), ", "))
	}
	return nazwa, opis, nil
}

// znaneModeleWycinania wypisuje zbiór modeli w kolejności stałej, żeby treść
// odmowy nie zmieniała się między wywołaniami.
func znaneModeleWycinania() []string {
	return []string{"u2net", "u2netp", "u2net_human_seg", "isnet-general-use", "silueta"}
}
