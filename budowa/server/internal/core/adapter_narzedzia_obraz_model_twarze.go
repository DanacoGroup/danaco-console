// Osobny przebieg poprawiania twarzy w `image.upscale` — pole `faces`
// kontraktu. Pomocnik pythonowy liczy na wagach PyTorcha, których wydanie
// ncnn nie publikuje, i pracuje jako drugi przebieg nad wynikiem powiększenia.
package core

import (
	"context"
	_ "embed"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"danacoconsole/server/internal/zewnetrzne"
)

// skryptPomocnikaTwarzy — treść pomocnika wkompilowana w binarium, żeby
// wdrożenie, w którym ktoś przeniósł samo binarium, nadal działało.
//
//go:embed adapter_narzedzia_obraz_pomocnik_twarzy.py
var skryptPomocnikaTwarzy string

const (
	// granicaOdtwarzaniaTwarzy to granica czasu jednego przebiegu twarzowego:
	// bez karty graficznej jedna twarz liczy się kilka sekund, a start
	// interpretera i wczytanie wag dokładają swoje.
	granicaOdtwarzaniaTwarzy = 10 * time.Minute

	// katalogWagTwarzyLinux to miejsce wag sieci twarzowej na serwerze, poza
	// `/usr/local/share/<silnik>`, bo wagi nie są składnikiem pakietu programu.
	katalogWagTwarzyLinux = "/opt/danaco-modele/twarze"

	// wagiOdtwarzaniaTwarzy to plik wag samej sieci odtwarzającej twarz,
	// wydany jako format PyTorcha, którego ncnn nie publikuje.
	wagiOdtwarzaniaTwarzy = "GFPGANv1.4.pth"

	// wagiWykrywaniaTwarzy to wagi wykrywacza twarzy (RetinaFace), bez których
	// nie ma czego odtwarzać: wskazują wycinek wyrównany do twarzy.
	wagiWykrywaniaTwarzy = "detection_Resnet50_Final.pth"

	// wagiPodzialuTwarzy to wagi sieci dzielącej twarz na obszary (ParseNet),
	// z których powstaje maska wklejenia wycinka.
	wagiPodzialuTwarzy = "parsing_parsenet.pth"

	// nazwaSkryptuTwarzy jest nazwą pliku wyłożonego w katalogu przebiegu,
	// pod którą przebieg twarzowy odnajduje własny pomocnik.
	nazwaSkryptuTwarzy = "pomocnik_twarzy.py"
	// prawaSkryptuTwarzy: skrypt czyta wyłącznie proces, który go wyłożył,
	// żaden inny użytkownik systemu nie ma do niego dostępu.
	prawaSkryptuTwarzy = 0o600
)

// narzedzieOdtwarzaniaTwarzy opisuje interpreter pomocnika twarzowego: idzie
// opakowaniem `/usr/local/bin/danaco-twarze`, a nie plikiem środowiska
// pythonowego, tak samo jak przy `rembg`.
func narzedzieOdtwarzaniaTwarzy() zewnetrzne.Narzedzie {
	return zewnetrzne.Narzedzie{
		Nazwa:   "GFPGAN (pomocnik pythonowy)",
		Program: "danaco-twarze",
		Pakiet: "torch, torchvision i facexlib w osobnym środowisku pythonowym wraz " +
			"z architekturą GFPGAN (gfpganv1_clean_arch, stylegan2_clean_arch z wydania " +
			"github.com/TencentARC/GFPGAN), wystawione opakowaniem /usr/local/bin/danaco-twarze; " +
			"wagi w " + katalogWagTwarzy(),
	}
}

// katalogWagTwarzy oddaje miejsce trzech zestawów wag przebiegu twarzowego:
// na Linuksie drzewo modeli stoi pod `/opt`, w wydaniu Windows obok rdzenia.
func katalogWagTwarzy() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(katalogWagWindows(), "twarze", "modele")
	}
	return katalogWagTwarzyLinux
}

// odpowiedzPomocnikaTwarzy to obiekt, który pomocnik wypisuje na standardowe
// wyjście. Nazwy pól odpowiadają co do znaku kluczom
// w `adapter_narzedzia_obraz_pomocnik_twarzy.py` — rozjazd zamieniłby nazwany
// brak w milczenie.
type odpowiedzPomocnikaTwarzy struct {
	Ok     bool   `json:"ok"`
	Powod  string `json:"powod"`
	Twarze int    `json:"twarze"`
}

// sprawdzWagiTwarzy upewnia się, że wszystkie trzy zestawy wag leżą na dysku,
// zanim ruszy pomocnik, żeby odmowa nie przyszła po minutach startu.
func sprawdzWagiTwarzy(katalog string) error {
	wykaz := []struct {
		plik string
		waga string
		skad string
	}{
		{wagiOdtwarzaniaTwarzy, "333 MB",
			"github.com/TencentARC/GFPGAN/releases (GFPGANv1.4.pth)"},
		{wagiWykrywaniaTwarzy, "104 MB",
			"github.com/xinntao/facexlib/releases (detection_Resnet50_Final.pth)"},
		{wagiPodzialuTwarzy, "82 MB",
			"github.com/xinntao/facexlib/releases (parsing_parsenet.pth)"},
	}
	for _, pozycja := range wykaz {
		if err := sprawdzWagi(filepath.Join(katalog, pozycja.plik),
			pozycja.plik, pozycja.waga, pozycja.skad); err != nil {
			return err
		}
	}
	return nil
}

// wylozPomocnikaTwarzy zapisuje skrypt w katalogu przebiegu i oddaje jego
// ścieżkę, przez plik tymczasowy i przemianowanie, żeby skrypt obcięty
// w połowie nie wywrócił się niepowiązanym błędem składni.
func wylozPomocnikaTwarzy(katalog string) (string, error) {
	czesciowy, err := os.CreateTemp(katalog, nazwaSkryptuTwarzy+".*.czesciowy")
	if err != nil {
		return "", bladZapleczaModeluObrazu(
			"nie można wyłożyć pomocnika twarzy: " + err.Error())
	}
	nazwa := czesciowy.Name()
	if _, err := czesciowy.WriteString(skryptPomocnikaTwarzy); err != nil {
		czesciowy.Close()
		os.Remove(nazwa)
		return "", bladZapleczaModeluObrazu("nie można zapisać pomocnika twarzy: " + err.Error())
	}
	if err := czesciowy.Close(); err != nil {
		os.Remove(nazwa)
		return "", bladZapleczaModeluObrazu("nie można domknąć pomocnika twarzy: " + err.Error())
	}
	if err := os.Chmod(nazwa, prawaSkryptuTwarzy); err != nil {
		os.Remove(nazwa)
		return "", bladZapleczaModeluObrazu("nie można nadać praw pomocnikowi twarzy: " + err.Error())
	}
	docelowy := filepath.Join(katalog, nazwaSkryptuTwarzy)
	if err := os.Rename(nazwa, docelowy); err != nil {
		os.Remove(nazwa)
		return "", bladZapleczaModeluObrazu("nie można wyłożyć pomocnika twarzy: " + err.Error())
	}
	return docelowy, nil
}

// poprawTwarze przeprowadza drugi przebieg nad wynikiem powiększenia i oddaje
// bajty obrazu z odtworzonymi twarzami, do osobnego pliku w tej samej
// pracowni, żeby wynik powiększenia został nietknięty na wypadek odmowy.
func (a *adapterNarzedziObrazuModelu) poprawTwarze(ctx context.Context,
	pracownia pracowniaObrazu) ([]byte, int, error) {

	katalog := katalogWagTwarzy()
	if err := sprawdzWagiTwarzy(katalog); err != nil {
		return nil, 0, err
	}
	skrypt, err := wylozPomocnikaTwarzy(pracownia.katalog)
	if err != nil {
		return nil, 0, err
	}
	wyjscie := filepath.Join(pracownia.katalog, "twarze.png")

	if a.wspolne == nil || a.wspolne.uruchamiacz == nil {
		return nil, 0, bladZapleczaNiedostepnegoModeluObrazu(
			"serwer nie ma uruchamiacza procesów — pomocnik twarzy nie ma czym wystartować; " +
				"naprawa: podpiąć warstwę kanału (injection) przy składaniu serwera")
	}
	okno, zasady, obszar := a.wspolne.zasiegNarzedzi(ctx)
	wynik, err := zewnetrzne.Wolaj(ctx, a.wspolne.uruchamiacz, okno, zasady, obszar,
		narzedzieOdtwarzaniaTwarzy(),
		[]string{skrypt, pracownia.wyjscie, wyjscie,
			filepath.Join(katalog, wagiOdtwarzaniaTwarzy), katalog},
		obszar.KatalogRoboczy, granicaOdtwarzaniaTwarzy)
	if err != nil {
		return nil, 0, bladArsenaluModeluObrazu(err)
	}

	// Pomocnik kończy pracę kodem zerowym także wtedy, gdy nazywa brak.

	// Rozstrzyga więc treść odpowiedzi, a nie kod wyjścia procesu.
	var odpowiedz odpowiedzPomocnikaTwarzy
	if err := json.Unmarshal(wynik.Wyjscie, &odpowiedz); err != nil {
		return nil, 0, bladPrzetwarzaniaModeluObrazu(
			"pomocnik twarzy oddał odpowiedź, której nie da się odczytać: " + err.Error())
	}
	if !odpowiedz.Ok {
		return nil, 0, bladZapleczaNiedostepnegoModeluObrazu(
			"przebieg poprawiania twarzy nie doszedł do skutku: " + odpowiedz.Powod)
	}
	bajty, err := odczytajWynikSilnika(wyjscie, "GFPGAN")
	if err != nil {
		return nil, 0, err
	}
	return bajty, odpowiedz.Twarze, nil
}

// opisPrzebieguTwarzy składa dopisek do opisu zasobu, żeby w magazynie było
// widać, ile twarzy przebieg poprawił — liczbę zmierzoną, nie założoną.
func opisPrzebieguTwarzy(twarze int) string {
	return ", twarze poprawione: " + strconv.Itoa(twarze)
}
