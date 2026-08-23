// Dwie czynności stojące na sieciach neuronowych — `image.upscale`
// (Real-ESRGAN) i `image.background.remove` (rembg / U²-Net) — wraz z opisem
// obu silników i składaniem ich wiersza poleceń. Wspólne zaplecze (źródło,
// pracownia, wołanie binarium, odmowy) stoi
// w `adapter_narzedzia_obraz_model.go`; metody stoją na tym samym adapterze.
//
// Oba silniki liczą na procesorze. Real-ESRGAN w wydaniu `ncnn-vulkan` jest
// jednym plikiem wykonywalnym bez Pythona i bez Torcha, a Vulkana dostaje od
// sterownika programowego (lavapipe z Mesy); wydanie pythonowe (`basicsr` +
// Torch) dałoby ten sam wynik za cenę kilku gigabajtów zależności. `rembg`
// jest Pythonem, ale jego runtime (ONNX Runtime) ma tryb procesorowy jako
// podstawowy, nie awaryjny.
//
// Czego tu celowo nie ma:
//  1. Gałęzi „gdy silnika nie ma, przeskaluj ImageMagickiem" — rozciągnięcie
//     oddane jako powiększenie jest atrapą, której nie widać do przybliżenia.
//  2. Poprawiania twarzy. Pole `faces` niesie kontrakt, ale przebieg twarzowy
//     robi osobna sieć (GFPGAN/CodeFormer), której wydanie `ncnn` tego silnika
//     nie zawiera. Żądanie z `faces: true` kończy się więc odmową nazywającą
//     brak, zamiast oddać obraz bez poprawki twarzy jako poprawiony.
package core

import (
	"context"
	"path/filepath"
	"strconv"
	"strings"

	"danacoconsole/server/internal/zewnetrzne"
	"danacoconsole/shared"
)

// narzedziePowiekszenia opisuje binarium superrozdzielczości. Nazwa czytelna
// i pakiet wchodzą do treści odmowy — Operator ma przeczytać, czego brakuje
// i skąd to wziąć, a ten silnik nie stoi w repozytorium dystrybucji, więc
// „pakietem" jest tu wydanie z sieci.
func narzedziePowiekszenia() zewnetrzne.Narzedzie {
	return zewnetrzne.Narzedzie{
		Nazwa:   "Real-ESRGAN (ncnn)",
		Program: "realesrgan-ncnn-vulkan",
		Pakiet: "wydanie realesrgan-ncnn-vulkan z github.com/xinntao/Real-ESRGAN/releases " +
			"rozpakowane do /usr/local/bin wraz z modelami w " + katalogModeliPowiekszenia() +
			", oraz mesa-vulkan-drivers dla liczenia na procesorze",
	}
}

// narzedzieWycinaniaTla opisuje binarium wycinania tła. Wołamy opakowanie
// `/usr/local/bin/rembg`, a nie plik z wnętrza środowiska pythonowego, bo
// `zewnetrzne.Wolaj` nie dziedziczy środowiska rdzenia, a `rembg` potrzebuje
// wskazania katalogu wag (`U2NET_HOME`). Opakowanie ustawia je samo, więc
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
//
// Nazwa modelu wchodzi w argument programu i dlatego jest stałą — ta sama
// reguła, co przy formatach `image.convert`. Sieć dobiera rdzeń, nie model
// językowy: kontrakt `image.upscale` nie ma pola na jej nazwę.
//
// `realesrgan-x4plus` jest wyborem dla zdjęć. Sieci `-anime` i `animevideov3`
// są uczone na rysunku i na fotografii zostawiają płaskie, plakatowe
// płaszczyzny.
const modelPowiekszeniaZdjec = "realesrgan-x4plus"

// dopuszczalneKrotnosci to krotności, które silnik ncnn przyjmuje: `-s 2|3|4`
// i nic więcej. Krotność spoza tego zbioru kończy się odmową wymieniającą
// dopuszczalne, a nie cichym zaokrągleniem do najbliższej — model, który prosił
// o ośmiokrotne, ma wiedzieć, że dostałby czterokrotne, zanim zobaczy wynik.
var dopuszczalneKrotnosci = map[int]struct{}{2: {}, 3: {}, 4: {}}

// wagiModelu opisuje jeden plik wag: gdzie leży, ile waży i skąd się go bierze.
type wagiModelu struct {
	// plik jest nazwą pliku wag w katalogu `katalogWagWycinania`.
	plik string
	// waga wchodzi do treści odmowy, żeby czytelnik wiedział, na co się pisze,
	// zanim ruszy pobieranie.
	waga string
	// skad wskazuje źródło wag — odmowa ma mówić, gdzie ich szukać.
	skad string
}

// modeleWycinaniaTla to zamknięty zbiór sieci wycinających tło wraz z wagą
// pliku i miejscem, z którego się go bierze; wszystko troje wchodzi do treści
// odmowy przy braku pliku.
//
// Kontrakt ma pole `model`, więc wybór należy do modelu językowego, ale jest to
// wybór ze zbioru, a nie dowolny tekst wpisany w argument programu. Każda
// pozycja ma powód: `u2net` jest domyślną siecią ogólną, `isnet-general-use`
// bywa dokładniejsza na cienkim szczególe (włosy, gałęzie), `u2net_human_seg`
// jest uczona na ludziach i na portrecie bije obie, `u2netp` jest wersją lekką
// dla maszyn bez zapasu pamięci.
var modeleWycinaniaTla = map[string]wagiModelu{
	"u2net":             {plik: "u2net.onnx", waga: "176 MB", skad: "github.com/danielgatis/rembg/releases (u2net.onnx)"},
	"u2netp":            {plik: "u2netp.onnx", waga: "4,7 MB", skad: "github.com/danielgatis/rembg/releases (u2netp.onnx)"},
	"u2net_human_seg":   {plik: "u2net_human_seg.onnx", waga: "176 MB", skad: "github.com/danielgatis/rembg/releases (u2net_human_seg.onnx)"},
	"isnet-general-use": {plik: "isnet-general-use.onnx", waga: "179 MB", skad: "github.com/danielgatis/rembg/releases (isnet-general-use.onnx)"},
	"silueta":           {plik: "silueta.onnx", waga: "44 MB", skad: "github.com/danielgatis/rembg/releases (silueta.onnx)"},
}

// domyslnyModelWycinania jest brany, gdy żądanie nie wskazuje modelu —
// kontrakt obiecuje „brak bierze domyslny silnika", a domyślnym silnika
// `rembg` jest właśnie `u2net`.
const domyslnyModelWycinania = "u2net"

// Powieksz obsługuje `image.upscale` — powiększenie z odtworzeniem szczegółu.
//
// Kolejność kroków jest zamierzona: najpierw rozstrzygamy żądanie (krotność,
// twarze), potem źródło, potem obecność wag, a dopiero na końcu ruszamy silnik.
// Każde z trzech pierwszych sprawdzeń kosztuje mikrosekundy, a przebieg silnika
// kosztuje minuty; odmowa dopiero po nim byłaby czasem straconym.
func (a *adapterNarzedziObrazuModelu) Powieksz(ctx context.Context,
	z shared.ImageUpscaleRequest) (shared.ImageUpscaleResponse, error) {

	krotnosc, err := rozstrzygnijKrotnosc(z.Scale)
	if err != nil {
		return shared.ImageUpscaleResponse{}, err
	}
	if z.Faces != nil && *z.Faces {
		return shared.ImageUpscaleResponse{}, bladZapleczaNiedostepnegoModeluObrazu(
			"osobny przebieg poprawiania twarzy (faces) wymaga sieci GFPGAN, " +
				"której wydanie ncnn silnika Real-ESRGAN nie zawiera — " +
				"rdzeń odmawia zamiast oddać obraz bez poprawki twarzy jako poprawiony; " +
				"naprawa: powtórzyć żądanie bez pola faces albo doinstalować gfpgan-ncnn-vulkan")
	}
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

	pracownia, err := przygotujPracownie(zrodlo.sciezka)
	if err != nil {
		return shared.ImageUpscaleResponse{}, err
	}
	defer pracownia.sprzatnij()

	// `-f png` wymuszamy jawnie, bo format wyniku ma być bezstratny: sieć
	// właśnie odtworzyła szczegół, a zapis stratny odjąłby część tego, za co
	// zapłacono minutami liczenia.
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
	zasob, _, err := a.wspolne.odlozZasob(ctx, zrodlo, z.WindowId, bajty, "png",
		"upscale x"+strconv.Itoa(krotnosc))
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
// Brak wskazania bierze krotność dwukrotną, tak jak obiecuje kontrakt. Wartość
// spoza zbioru silnika kończy się odmową wymieniającą dopuszczalne — patrz
// `dopuszczalneKrotnosci`.
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

// UsunTlo obsługuje `image.background.remove` — wycięcie obiektu z tła.
//
// Wynik jest zawsze w PNG, bo przezroczystość ma gdzie się zapisać wyłącznie
// w formacie z kanałem alfa. Ten sam wynik w JPEG-u dałby biały prostokąt
// w miejscu przezroczystości.
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

	// Podpolecenie `i` jest trybem „jeden plik na jeden plik". Tryby `p`
	// (katalog) i `s` (serwer HTTP) są tu niepotrzebne i szersze, niż trzeba:
	// drugi otwierałby gniazdo, o które nikt nie prosił.
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
// znane — model językowy, który zgadł nazwę z pamięci, ma przeczytać listę,
// a nie milczący błąd silnika o nieznanej sesji ONNX.
func rozstrzygnijModelWycinania(zadany *string) (string, wagiModelu, error) {
	nazwa := domyslnyModelWycinania
	if !bezWartosci(zadany) {
		nazwa = strings.ToLower(strings.TrimSpace(*zadany))
	}
	opis, jest := modeleWycinaniaTla[nazwa]
	if !jest {
		return "", opis, bladWskazaniaModeluObrazu("nie znam modelu wycinania „" + nazwa +
			"” — rdzeń zna: " + strings.Join(znaneModeleWycinania(), ", "))
	}
	return nazwa, opis, nil
}

// znaneModeleWycinania wypisuje zbiór modeli w kolejności stałej, żeby treść
// odmowy nie zmieniała się między wywołaniami — mapa w Go chodzi losowo,
// a odmowa czytana dwa razy ma brzmieć tak samo.
func znaneModeleWycinania() []string {
	return []string{"u2net", "u2netp", "u2net_human_seg", "isnet-general-use", "silueta"}
}
