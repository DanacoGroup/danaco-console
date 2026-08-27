// Odpowiedzialność pliku: wyłożenie pomocnika osadzeń na dysk i wskazanie
// interpretera, który go uruchomi. Skrypt jest wkompilowany w binarium przez
// `go:embed`, a nie szukany obok niego.
package wiedza

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// skryptPomocnika — treść pomocnika osadzeń wkompilowana w binarium przez
// dyrektywę embed poniżej, gotowa do wyłożenia na dysk przy uruchomieniu.
//go:embed pomocnik_osadzen.py
var skryptPomocnika string

// skryptPrzesiewu — treść pomocnika przesiewu wkompilowana w binarium
// przez dyrektywę embed poniżej, gotowa do wyłożenia na dysk.
//go:embed pomocnik_przesiewu.py
var skryptPrzesiewu string

// skryptObrazu — treść pomocnika osi obrazu wkompilowana w binarium przez
// dyrektywę embed poniżej, gotowa do wyłożenia na dysk przy uruchomieniu.
//go:embed pomocnik_obrazu.py
var skryptObrazu string

const (
	// nazwaSkryptu jest nazwą pliku pomocnika osadzeń wyłożonego na dysk,
	// pod którą go odnajduje interpreter przy uruchomieniu.
	nazwaSkryptu = "pomocnik_osadzen.py"
	// nazwaSkryptuPrzesiewu i nazwaSkryptuObrazu są nazwami plików dwóch
	// pozostałych pomocników, leżących pod własną nazwą w tym samym katalogu.
	nazwaSkryptuPrzesiewu = "pomocnik_przesiewu.py"
	nazwaSkryptuObrazu    = "pomocnik_obrazu.py"
	// podkatalogWiedzy oddziela rzeczy wskaźnika znaczenia od reszty katalogu
	// danych rdzenia, bazy i magazynu biblioteki.
	podkatalogWiedzy = "wiedza"
	// podkatalogModeli mieści wagi modelu pobrane przez bibliotekę osadzeń,
	// osobno od skryptów pomocników.
	podkatalogModeli = "modele"
	// podkatalogZlecen mieści pliki zleceń pomocnika, byty ulotne, kasowane
	// zaraz po uruchomieniu, nie leżące obok skryptu ani wag.
	podkatalogZlecen = "zlecenia"

	// prawaKatalogu i prawaPliku: treść Operatora należy do Operatora, który
	// uruchomił rdzeń, tak samo jak magazyn biblioteki.
	prawaKatalogu = 0o700
	prawaPliku    = 0o600

	// interpreterPreferowany — trójka jawnie, bo goła nazwa `python` na wielu
	// systemach wciąż wskazuje wydanie drugie języka.
	interpreterPreferowany = "python3"
	// interpreterZapasowy wchodzi tam, gdzie `python3` nie istnieje jako osobne
	// polecenie (typowo Windows i część obrazów kontenerowych).
	interpreterZapasowy = "python"
)

// wylozSkrypt zapisuje pomocnika pod katalogiem danych i oddaje jego ścieżkę.
// Zapis powtarza się przy każdym wywołaniu, nie tylko przy pierwszym.
func wylozSkrypt(katalogDanych string) (string, error) {
	return wylozPomocnika(katalogDanych, nazwaSkryptu, skryptPomocnika)
}

// wylozPomocnika wykłada jeden wkompilowany skrypt pod jego własną nazwą.
// Trzej pomocnicy pakietu — osadzenia, przesiew i oś obrazu — jadą tą samą
// drogą.
func wylozPomocnika(katalogDanych, nazwa, tresc string) (string, error) {
	katalog := filepath.Join(katalogDanych, podkatalogWiedzy)
	if err := os.MkdirAll(katalog, prawaKatalogu); err != nil {
		return "", fmt.Errorf("wskaźnik znaczenia: katalog %s: %w", katalog, err)
	}
	docelowy := filepath.Join(katalog, nazwa)

	tymczasowy, err := os.CreateTemp(katalog, nazwa+".*.czesciowy")
	if err != nil {
		return "", fmt.Errorf("wskaźnik znaczenia: plik tymczasowy w %s: %w", katalog, err)
	}
	nazwaTymczasowa := tymczasowy.Name()
	if _, err := tymczasowy.WriteString(tresc); err != nil {
		tymczasowy.Close()
		os.Remove(nazwaTymczasowa)
		return "", fmt.Errorf("wskaźnik znaczenia: zapis pomocnika: %w", err)
	}
	if err := tymczasowy.Close(); err != nil {
		os.Remove(nazwaTymczasowa)
		return "", fmt.Errorf("wskaźnik znaczenia: domknięcie pomocnika: %w", err)
	}
	if err := os.Chmod(nazwaTymczasowa, prawaPliku); err != nil {
		os.Remove(nazwaTymczasowa)
		return "", fmt.Errorf("wskaźnik znaczenia: prawa pomocnika: %w", err)
	}
	if err := os.Rename(nazwaTymczasowa, docelowy); err != nil {
		os.Remove(nazwaTymczasowa)
		return "", fmt.Errorf("wskaźnik znaczenia: wyłożenie pomocnika: %w", err)
	}
	return docelowy, nil
}

// zapiszZlecenie odkłada treść zlecenia w pliku, bo droga wołania procesu nie
// podaje mu standardowego wejścia. Oddaje ścieżkę i funkcję sprzątającą.
func zapiszZlecenie(katalogDanych string, tresc []byte) (string, func(), error) {
	katalog := filepath.Join(katalogDanych, podkatalogWiedzy, podkatalogZlecen)
	if err := os.MkdirAll(katalog, prawaKatalogu); err != nil {
		return "", func() {}, fmt.Errorf("wskaźnik znaczenia: katalog zleceń %s: %w", katalog, err)
	}
	plik, err := os.CreateTemp(katalog, "zlecenie-*.json")
	if err != nil {
		return "", func() {}, fmt.Errorf("wskaźnik znaczenia: plik zlecenia: %w", err)
	}
	nazwa := plik.Name()
	sprzatanie := func() { os.Remove(nazwa) }
	if _, err := plik.Write(tresc); err != nil {
		plik.Close()
		sprzatanie()
		return "", func() {}, fmt.Errorf("wskaźnik znaczenia: zapis zlecenia: %w", err)
	}
	if err := plik.Close(); err != nil {
		sprzatanie()
		return "", func() {}, fmt.Errorf("wskaźnik znaczenia: domknięcie zlecenia: %w", err)
	}
	return nazwa, sprzatanie, nil
}

// odnajdzInterpreter rozstrzyga, który program uruchomi skrypt. Wskazanie
// Operatora idzie wprost i bez sprawdzania na dysku, świadomie.
func odnajdzInterpreter(program string) string {
	if program != "" {
		return program
	}
	if zeSciezki, err := exec.LookPath(interpreterPreferowany); err == nil {
		return zeSciezki
	}
	return interpreterZapasowy
}

// katalogWag rozstrzyga, gdzie leżą pobrane wagi modelu. Wskazanie Operatora
// ma pierwszeństwo; jego brak znaczy podkatalog katalogu danych rdzenia.
func katalogWag(katalogDanych, wskazanie string) string {
	if wskazanie != "" {
		return wskazanie
	}
	return filepath.Join(katalogDanych, podkatalogWiedzy, podkatalogModeli)
}

// katalogWagOsobny rozstrzyga, gdzie leżą wagi modelu innego niż osadzenia,
// w osobnym podkatalogu, nie wspólnym worku.
func katalogWagOsobny(katalogDanych, wskazanie, podkatalog string) string {
	if wskazanie != "" {
		return wskazanie
	}
	return filepath.Join(katalogDanych, podkatalogWiedzy, podkatalogModeli, podkatalog)
}
