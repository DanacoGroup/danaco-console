// Osobny przebieg poprawiania twarzy w `image.upscale` — pole `faces`
// kontraktu. Powiększanie samo stoi
// w `adapter_narzedzia_obraz_model_silniki.go` i to ono woła ten plik; wspólne
// zaplecze (pracownia, wołanie binarium, odmowy) — w
// `adapter_narzedzia_obraz_model.go`.
//
// ── Dlaczego pomocnik pythonowy, a nie wydanie ncnn ─────────────────────────
// Reszta tej rodziny to binaria `ncnn` bez Pythona i bez Torcha, więc pomocnik
// pythonowy jest tu wyłomem i wymaga powodu. Powód jest taki: sieć twarzowa
// wydana jest jako wagi PyTorcha (`GFPGANv1.4.pth`), a wydania `ncnn` tej sieci
// nie publikuje jej autor — chodzące po sieci przeróbki niosą wagi przeliczone
// przez osoby trzecie, więc rdzeń liczyłby nie tym modelem, który leży na
// dysku, tylko czyjąś kopią o nieustalonym pochodzeniu. Pomocnik na wagach
// stojących liczy dokładnie tym, co Operator ma u siebie, i tą samą drogą, co
// wektory znaczenia (`internal/wiedza/pomocnik_osadzen.py`).
//
// ── Dlaczego przebieg jest drugi, a nie jeden wspólny ───────────────────────
// Kolejność jest zamierzona: najpierw Real-ESRGAN powiększa CAŁY obraz, potem
// pomocnik odnajduje twarze w wyniku i podmienia same wycinki. Dzięki temu
// wymiary odpowiedzi pochodzą wyłącznie z powiększenia, a `faces: false`
// i `faces: true` różnią się dokładnie tym, co obiecuje opis pola — twarzami,
// nie rozmiarem.
//
// ── Czego tu celowo nie ma ──────────────────────────────────────────────────
// Gałęzi „gdy pomocnika nie ma, oddaj samo powiększenie". Obraz bez poprawki
// twarzy podany jako poprawiony jest tą samą atrapą, co rozciągnięcie podane
// jako powiększenie — rdzeń odmawia, nazywając brak i drogę naprawy.
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

// skryptPomocnikaTwarzy — treść pomocnika wkompilowana w binarium.
//
// Skrypt jedzie w binarium, a nie leży obok niego, z tego samego powodu, co
// pomocnik osadzeń: wdrożenie, w którym ktoś przeniósł samo binarium, ma
// działać. Wykładany jest do katalogu jednego przebiegu (`pracowniaObrazu`),
// bo znika razem z nim — pomocnik jest bytem wtórnym wobec binarium i nie ma
// powodu przeżywać żądania, które go potrzebowało.
//
//go:embed adapter_narzedzia_obraz_pomocnik_twarzy.py
var skryptPomocnikaTwarzy string

const (
	// granicaOdtwarzaniaTwarzy to granica czasu jednego przebiegu twarzowego.
	// Bez karty graficznej jedna twarz liczy się na procesorze kilka sekund,
	// a zdjęcie grupowe niesie ich kilkanaście; do tego dochodzi start
	// interpretera i wczytanie trzech zestawów wag. Dziesięć minut znaczy „coś
	// stanęło", a nie „to długo trwa".
	granicaOdtwarzaniaTwarzy = 10 * time.Minute

	// katalogWagTwarzyLinux to miejsce wag sieci twarzowej na serwerze.
	// Odbiega od `/usr/local/share/<silnik>`, którym idą wagi powiększania
	// i wycinania tła, bo te wagi nie są składnikiem pakietu żadnego programu —
	// są wydaniem modelu pobieranym osobno i leżą we wspólnym drzewie modeli
	// rdzenia.
	katalogWagTwarzyLinux = "/opt/danaco-modele/twarze"

	// wagiOdtwarzaniaTwarzy to plik wag samej sieci odtwarzającej twarz.
	wagiOdtwarzaniaTwarzy = "GFPGANv1.4.pth"
	// wagiWykrywaniaTwarzy to wagi wykrywacza twarzy (RetinaFace). Bez niego
	// nie ma czego odtwarzać: sieć twarzowa pracuje na wycinku wyrównanym do
	// pięciu punktów charakterystycznych, a te punkty wskazuje właśnie ten
	// model.
	wagiWykrywaniaTwarzy = "detection_Resnet50_Final.pth"
	// wagiPodzialuTwarzy to wagi sieci dzielącej twarz na obszary
	// (ParseNet). Z niej powstaje maska wklejenia — bez maski wycinek wraca do
	// obrazu prostokątem o widocznej krawędzi.
	wagiPodzialuTwarzy = "parsing_parsenet.pth"

	// nazwaSkryptuTwarzy jest nazwą pliku wyłożonego w katalogu przebiegu.
	nazwaSkryptuTwarzy = "pomocnik_twarzy.py"
	// prawaSkryptuTwarzy: skrypt czyta wyłącznie proces, który go wyłożył.
	prawaSkryptuTwarzy = 0o600
)

// narzedzieOdtwarzaniaTwarzy opisuje interpreter pomocnika twarzowego.
//
// Wołamy opakowanie `/usr/local/bin/danaco-twarze`, a nie plik z wnętrza
// środowiska pythonowego — tak samo jak przy `rembg`. `zewnetrzne.Wolaj` nie
// dziedziczy środowiska rdzenia, a biblioteki pomocnika szukają katalogu pamięci
// podręcznej i katalogu domowego; opakowanie ustawia je samo, więc pomocnik jest
// samowystarczalny niezależnie od tego, kto go woła.
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

// katalogWagTwarzy oddaje miejsce, w którym leżą trzy zestawy wag przebiegu
// twarzowego. Zależy od systemu z tego samego powodu, co katalogi wag
// powiększania i wycinania tła: na Linuksie drzewo modeli stoi pod `/opt`,
// a w wydaniu natywnym Windows jedzie obok rdzenia w `pomocniki/`.
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
// zanim ruszy pomocnik. Sprawdzane są wszystkie naraz, bo przebieg potrzebuje
// każdego z nich, a odmowa po dwóch minutach startu interpretera z powodu
// trzeciego pliku byłaby czasem straconym.
//
// Katalog przychodzi argumentem, a nie jest brany ze stałej, żeby sprawdzian
// mógł zmierzyć samą odmowę na katalogu bez wag — reguła, której nie da się
// uruchomić na pustym miejscu, nie jest zmierzona.
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
// ścieżkę. Zapis idzie przez plik tymczasowy i przemianowanie, bo skrypt obcięty
// w połowie wystartowałby i wywrócił się komunikatem o składni, którego nikt nie
// powiąże z przerwanym zapisem.
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
// bajty obrazu z odtworzonymi twarzami.
//
// Wejściem jest plik wyniku Real-ESRGAN-a, a nie źródło żądania: przebieg
// twarzowy pracuje na tym, co powiększenie już wytworzyło. Wyjście idzie do
// osobnego pliku w tej samej pracowni, żeby wynik powiększenia został nietknięty
// na wypadek odmowy pomocnika — nadpisanie go w miejscu zostawiłoby przy
// przerwanym zapisie plik, który nie jest już ani jednym, ani drugim.
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
			"rdzeń nie ma uruchamiacza procesów — pomocnik twarzy nie ma czym wystartować; " +
				"naprawa: podpiąć warstwę kanału (injection) przy składaniu rdzenia")
	}
	okno, zasady, obszar := a.wspolne.zasiegNarzedzi()
	wynik, err := zewnetrzne.Wolaj(ctx, a.wspolne.uruchamiacz, okno, zasady, obszar,
		narzedzieOdtwarzaniaTwarzy(),
		[]string{skrypt, pracownia.wyjscie, wyjscie,
			filepath.Join(katalog, wagiOdtwarzaniaTwarzy), katalog},
		obszar.KatalogRoboczy, granicaOdtwarzaniaTwarzy)
	if err != nil {
		return nil, 0, bladArsenaluModeluObrazu(err)
	}

	// Pomocnik kończy pracę kodem zerowym także wtedy, gdy nazywa brak, więc
	// rozstrzyga treść odpowiedzi, a nie kod wyjścia. Odpowiedź nieczytelna
	// jest odmową: proces, który wypisał coś innego niż umówiony obiekt, nie
	// zrobił tego, po co go zawołano.
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
// widać, ile twarzy przebieg poprawił. Liczba jest ZMIERZONA przez pomocnika,
// a nie założona — zdjęcie bez rozpoznanej twarzy przechodzi przebieg
// nietknięte i opis ma to mówić wprost.
func opisPrzebieguTwarzy(twarze int) string {
	return ", twarze poprawione: " + strconv.Itoa(twarze)
}
