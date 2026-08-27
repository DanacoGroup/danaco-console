// Odpowiedzialność pliku: wyłożenie pomocnika osadzeń na dysk i wskazanie
// interpretera, który go uruchomi.
//
// Skrypt jest wkompilowany w binarium, a nie szukany obok niego. Silnik mowy
// szuka swojego pomocnika w katalogu `pomocniki/` obok binarium (patrz
// `mowa/pomocnik.go`) i płaci za to odmową na każdym wdrożeniu, w którym ktoś
// przeniósł samo binarium. Tutaj skrypt jedzie w binarium przez `go:embed` —
// tak samo jak migracje schematu (`store/zrodlo_migracji.go`) — i wykłada się
// na dysk przy pierwszym użyciu. Powód jest ten sam co tam: wdrożenie nie ma
// zależeć od obecności plików obok programu.
//
// Wykładany jest do katalogu danych, nie do katalogu tymczasowego. Ten sam
// katalog niesie bazę, sejf poświadczeń i magazyn treści biblioteki — pomocnik
// przeżywa restart rdzenia tak jak one, a Operator ma go gdzie obejrzeć, zanim
// pozwoli mu ruszyć. Zapis jest atomowy (plik tymczasowy w katalogu docelowym,
// potem przemianowanie), bo skrypt obcięty w połowie wystartowałby i wywrócił
// się komunikatem o składni, którego nikt nie powiąże z przerwanym zapisem.
//
// Ten plik niczego nie uruchamia. Start procesu należy wyłącznie do
// `zewnetrzne.Wolaj`. Jedynym wyjątkiem jest `exec.LookPath`, które
// nie uruchamia procesu, a jedynie przegląda ścieżkę wyszukiwania systemu —
// dokładnie tak, jak robi to `mowa/pomocnik.go`.
package wiedza

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// skryptPomocnika — treść pomocnika wkompilowana w binarium.
//
//go:embed pomocnik_osadzen.py
var skryptPomocnika string

// skryptPrzesiewu — treść pomocnika przesiewu wkompilowana w binarium.
//
//go:embed pomocnik_przesiewu.py
var skryptPrzesiewu string

// skryptObrazu — treść pomocnika osi obrazu wkompilowana w binarium.
//
//go:embed pomocnik_obrazu.py
var skryptObrazu string

const (
	// nazwaSkryptu jest nazwą pliku wyłożonego na dysk.
	nazwaSkryptu = "pomocnik_osadzen.py"
	// nazwaSkryptuPrzesiewu i nazwaSkryptuObrazu są nazwami plików dwóch
	// pozostałych pomocników. Każdy leży pod własną nazwą w tym samym katalogu:
	// jeden plik o zmiennej treści nie dałby się obejrzeć przed uruchomieniem,
	// a właśnie po to katalog danych jest miejscem wyłożenia.
	nazwaSkryptuPrzesiewu = "pomocnik_przesiewu.py"
	nazwaSkryptuObrazu    = "pomocnik_obrazu.py"
	// podkatalogWiedzy oddziela rzeczy wskaźnika od reszty katalogu danych.
	podkatalogWiedzy = "wiedza"
	// podkatalogModeli mieści wagi pobrane przez bibliotekę.
	podkatalogModeli = "modele"
	// podkatalogZlecen mieści pliki zleceń pomocnika. Osobny poziom, bo są to
	// byty ULOTNE — kasowane zaraz po uruchomieniu — i nie mają leżeć obok
	// skryptu ani obok wag.
	podkatalogZlecen = "zlecenia"

	// prawaKatalogu i prawaPliku: treść Operatora należy do Operatora, który
	// uruchomił rdzeń — tak samo jak magazyn biblioteki.
	prawaKatalogu = 0o700
	prawaPliku    = 0o600

	// interpreterPreferowany — trójka jawnie, bo goła nazwa `python` na wielu
	// systemach wciąż wskazuje wydanie drugie.
	interpreterPreferowany = "python3"
	// interpreterZapasowy wchodzi tam, gdzie `python3` nie istnieje jako osobne
	// polecenie (typowo Windows i część obrazów kontenerowych).
	interpreterZapasowy = "python"
)

// wylozSkrypt zapisuje pomocnika pod katalogiem danych i oddaje jego ścieżkę.
//
// Zapis powtarza się przy każdym wywołaniu, a nie tylko przy pierwszym: skrypt
// jest bytem wtórnym wobec binarium, więc po podmianie rdzenia na nowsze
// wydanie na dysku ma leżeć wersja z tego binarium. Koszt to zapis kilku
// kilobajtów raz na żądanie indeksowania.
func wylozSkrypt(katalogDanych string) (string, error) {
	return wylozPomocnika(katalogDanych, nazwaSkryptu, skryptPomocnika)
}

// wylozPomocnika wykłada jeden wkompilowany skrypt pod jego własną nazwą.
// Trzej pomocnicy pakietu — osadzenia, przesiew i oś obrazu — jadą tą samą
// drogą: druga droga wykładania byłaby drugą prawdą o tym, gdzie Operator ma
// szukać kodu, który rdzeń uruchamia na jego maszynie.
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
// podaje mu standardowego wejścia (`zewnetrzne/wolanie.go`). Oddaje ścieżkę
// i funkcję sprzątającą — wołający kasuje plik `defer`-em, także na ścieżce
// błędu, żeby katalog zleceń nie rósł o jeden plik na każde żądanie.
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

// odnajdzInterpreter rozstrzyga, który program uruchomi skrypt.
//
// Wskazanie Operatora bierzemy wprost i bez sprawdzania na dysku — świadomie,
// bo może być nazwą do rozwinięcia przez system albo dowiązaniem środowiska
// wirtualnego, a odmowa na podstawie własnego sprawdzenia unieważniałaby
// ustawienie. Gdy wskazania nie ma, szukamy `python3`, a gdy i tego nie ma —
// zostaje `python`; ta ostatnia wartość jest zgadywana i odmowa przyjdzie
// dopiero z uruchomienia, bo tylko ono zna prawdę o wykonywalności.
func odnajdzInterpreter(program string) string {
	if program != "" {
		return program
	}
	if zeSciezki, err := exec.LookPath(interpreterPreferowany); err == nil {
		return zeSciezki
	}
	return interpreterZapasowy
}

// katalogWag rozstrzyga, gdzie leżą pobrane wagi modelu.
//
// Wskazanie Operatora ma pierwszeństwo; jego brak znaczy podkatalog katalogu
// danych rdzenia, a nie katalog domyślny biblioteki. Różnica jest istotna:
// biblioteka domyślnie pisze do katalogu pamięci podręcznej użytkownika, a droga
// wołania procesu nie dziedziczy środowiska, więc pomocnik nie zna nawet `HOME`
// i wagi wylądowałyby w miejscu zależnym od tego, gdzie akurat stoi rdzeń.
func katalogWag(katalogDanych, wskazanie string) string {
	if wskazanie != "" {
		return wskazanie
	}
	return filepath.Join(katalogDanych, podkatalogWiedzy, podkatalogModeli)
}

// katalogWagOsobny rozstrzyga, gdzie leżą wagi modelu innego niż osadzenia.
//
// Osobny podkatalog na model, a nie wspólny worek: biblioteka wykłada wagi
// modelu stojącego wprost w katalogu wskazanym, więc dwa modele w jednym
// katalogu byłyby dwoma plikami `model.safetensors` w tym samym miejscu —
// czyli jednym z nich nadpisanym przez drugi.
func katalogWagOsobny(katalogDanych, wskazanie, podkatalog string) string {
	if wskazanie != "" {
		return wskazanie
	}
	return filepath.Join(katalogDanych, podkatalogWiedzy, podkatalogModeli, podkatalog)
}
