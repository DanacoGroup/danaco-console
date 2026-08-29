// Adapter obsługuje `image.upscale` i `image.background.remove`: uruchamia sieć neuronową
// nad obrazem, weryfikuje obecność wag na dysku i odkłada wynik w magazynie zasobów Designu,
// dzieląc mechanizm źródła i izolacji z rodziną `image.*`.
package core

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/zewnetrzne"
	"danacoconsole/shared"
)

const (
	// granicaPowiekszenia to granica czasu jednego przebiegu superrozdzielczości bez karty
	// graficznej; piętnaście minut mieści powiększenie materiału aparatowego na procesorze,
	// nie ucinając pracy udanej w połowie.
	granicaPowiekszenia = 15 * time.Minute

	// granicaWycinaniaTla jest krótsza, bo przebieg jest jeden i stały:
	// obraz idzie do sieci przeskalowany do 320×320, więc czas prawie nie
	// zależy od wielkości źródła. Pięć minut mieści start interpretera,
	// wczytanie wag i przebieg z ogromnym zapasem.
	granicaWycinaniaTla = 5 * time.Minute

	// katalogModeliPowiekszeniaLinux i katalogWagWycinaniaLinux to miejsca wag na serwerze,
	// gdzie arsenał pakietu Linux je stawia; na Windowsie ścieżkę składają osobne funkcje
	// względem pliku wykonywalnego.
	katalogModeliPowiekszeniaLinux = "/usr/local/share/realesrgan-ncnn-vulkan/models"
	katalogWagWycinaniaLinux       = "/usr/local/share/rembg-modele"
)

// katalogModeliPowiekszenia oddaje miejsce wag sieci powiększającej: silnik ncnn szuka ich
// domyślnie względem katalogu bieżącego wołania, który jest za każdym razem inny, więc
// ścieżkę rdzeń podaje jawnie i zależnie od systemu.
func katalogModeliPowiekszenia() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(katalogWagWindows(), "realesrgan-ncnn-vulkan", "models")
	}
	return katalogModeliPowiekszeniaLinux
}

// katalogWagWycinania oddaje miejsce wag sieci wycinającej tło tą samą zasadą co wagi
// powiększania; opakowanie silnika ustawia zmienną środowiskową na ten sam katalog, który
// zwraca ta funkcja.
func katalogWagWycinania() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(katalogWagWindows(), "rembg", "modele")
	}
	return katalogWagWycinaniaLinux
}

// katalogWagWindows składa katalog pomocniczych programów obok pliku wykonywalnego rdzenia,
// tą samą drogą, którą odnajdywanie programów arsenału korzysta na Windowsie.
func katalogWagWindows() string {
	if biezace, err := os.Executable(); err == nil {
		return filepath.Join(filepath.Dir(biezace), "pomocniki")
	}
	return "pomocniki"
}

// adapterNarzedziObrazuModelu wypełnia port narzędzi modelu obrazu; źródło, zasięg izolacji
// i magazyn wyniku bierze z adaptera rodziny `image.*`, dokładając dwa silniki neuronowe,
// ich wagi i granice czasu.
type adapterNarzedziObrazuModelu struct {
	wspolne *adapterNarzedziObrazu
}

// nowyAdapterNarzedziObrazuModelu składa adapter na zapleczu przychodzącym z zewnątrz, tym
// samym, na którym stoi rodzina `image.*`, żeby uchwyt magazynu i uruchamiacz procesów nie
// miały drugiej prawdy.
func nowyAdapterNarzedziObrazuModelu(wspolne *adapterNarzedziObrazu) *adapterNarzedziObrazuModelu {
	return &adapterNarzedziObrazuModelu{wspolne: wspolne}
}

// pracowniaObrazu jest katalogiem jednego przebiegu silnika, niosącym dowiązanie do źródła
// i plik wyniku; oba silniki żądają ścieżek na dysku i nie umieją pisać obrazu na wyjście
// standardowe.
type pracowniaObrazu struct {
	katalog string
	wejscie string
	wyjscie string
}

// przygotujPracownie zakłada katalog przebiegu i wystawia w nim źródło dowiązaniem pod
// nazwą, którą silnik przyjmie, a przy braku dowiązań kopią zapasową, gdy magazyn leży na
// innym nośniku.
func przygotujPracownie(zrodlo string) (pracowniaObrazu, error) {
	katalog, err := os.MkdirTemp("", "danaco-obraz-model-")
	if err != nil {
		return pracowniaObrazu{}, bladZapleczaModeluObrazu(
			"nie można założyć katalogu pracy przebiegu: " + err.Error())
	}
	wejscie := filepath.Join(katalog, "zrodlo.png")
	if err := os.Symlink(zrodlo, wejscie); err != nil {
		bajty, blad := os.ReadFile(zrodlo)
		if blad != nil {
			_ = os.RemoveAll(katalog)
			return pracowniaObrazu{}, bladZapleczaModeluObrazu(
				"nie można wystawić źródła do pracy silnika: " + blad.Error())
		}
		if blad := os.WriteFile(wejscie, bajty, 0o600); blad != nil {
			_ = os.RemoveAll(katalog)
			return pracowniaObrazu{}, bladZapleczaModeluObrazu(
				"nie można wystawić źródła do pracy silnika: " + blad.Error())
		}
	}
	return pracowniaObrazu{
		katalog: katalog,
		wejscie: wejscie,
		wyjscie: filepath.Join(katalog, "wynik.png"),
	}, nil
}

// sprzatnij kasuje katalog przebiegu wraz z dowiązaniem i wynikiem. Wołany
// z `defer`, więc idzie tą samą drogą po udanym przebiegu i po odmowie —
// katalog zostawiony po nieudanym uruchomieniu rósłby w nieskończoność.
func (p pracowniaObrazu) sprzatnij() {
	if p.katalog != "" {
		_ = os.RemoveAll(p.katalog)
	}
}

// odczytajWynikSilnika bierze bajty pliku, który silnik miał wytworzyć; plik nieobecny albo
// pusty jest odmową mimo zerowego kodu wyjścia, bo oba silniki potrafią zakończyć się
// powodzeniem bez wyniku.
func odczytajWynikSilnika(sciezka, silnik string) ([]byte, error) {
	bajty, err := os.ReadFile(sciezka)
	if err != nil {
		return nil, bladPrzetwarzaniaModeluObrazu(silnik +
			" zakończył się powodzeniem, ale nie zostawił pliku wyniku: " + err.Error())
	}
	if len(bajty) == 0 {
		return nil, bladPrzetwarzaniaModeluObrazu(silnik +
			" zakończył się powodzeniem, ale plik wyniku jest pusty")
	}
	return bajty, nil
}

// wolajSilnik przeprowadza jedno uruchomienie silnika neuronowego wspólną drogą wołania
// binarium, przez port uruchamiacza, bramę izolacji okna i objęcie drzewa procesów; zasięg
// bierze z rodziny `image.*`.
func (a *adapterNarzedziObrazuModelu) wolajSilnik(ctx context.Context,
	narzedzie zewnetrzne.Narzedzie, argumenty []string, limit time.Duration) error {

	if a.wspolne == nil || a.wspolne.uruchamiacz == nil {
		return bladZapleczaNiedostepnegoModeluObrazu(
			"serwer nie ma uruchamiacza procesów — silniki obrazu nie mają czym wystartować; " +
				"naprawa: podpiąć warstwę kanału (injection) przy składaniu serwera")
	}
	okno, zasady, obszar := a.wspolne.zasiegNarzedzi()
	_, err := zewnetrzne.Wolaj(ctx, a.wspolne.uruchamiacz, okno, zasady, obszar,
		narzedzie, argumenty, strings.TrimSpace(obszar.KatalogRoboczy), limit)
	if err != nil {
		return bladArsenaluModeluObrazu(err)
	}
	return nil
}

// sprawdzWagi upewnia się, że plik wag leży na dysku, zanim ruszy silnik.
// Odmowa niesie ścieżkę i wagę pliku, żeby wiadomo było, co dociągnąć i ile to
// zajmie.
func sprawdzWagi(sciezka, nazwaModelu, waga, skad string) error {
	opis, err := os.Stat(sciezka)
	if err == nil && !opis.IsDir() && opis.Size() > 0 {
		return nil
	}
	return bladZapleczaNiedostepnegoModeluObrazu("nie ma wag modelu " + nazwaModelu +
		" — mają leżeć pod " + sciezka + " i ważą około " + waga +
		"; naprawa: pobrać je z " + skad + " i położyć pod tą ścieżką " +
		"(serwer świadomie nie ściąga ich w trakcie żądania, żeby czas łącza " +
		"nie zjadał granicy czasu przetwarzania)")
}

// wymiaryWyniku wyciąga zmierzone wymiary z odłożonego zasobu; brak wymiarów jest odmową,
// bo kontrakt obiecuje liczby, a wymiary pochodzą z pomiaru pliku, nie z pomnożenia żądanej
// krotności przez wymiar źródła.
func wymiaryWyniku(zasob shared.DesignAsset) (int, int, error) {
	if zasob.Width == nil || zasob.Height == nil {
		return 0, 0, bladPrzetwarzaniaModeluObrazu(
			"wynik legł w magazynie, ale nie dało się zmierzyć jego wymiarów — " +
				"odpowiedź nie poda liczb, których nie zmierzono")
	}
	return *zasob.Width, *zasob.Height, nil
}

// bladArsenaluModeluObrazu przekłada odmowę pakietu zewnętrznego na kod kontraktu: brak
// silnika dostaje inny kod niż niepowodzenie silnika, bo są to dwa różne zakończenia
// żądania.
func bladArsenaluModeluObrazu(err error) error {
	var brak *zewnetrzne.BrakNarzedzia
	if errors.As(err, &brak) {
		return bladZapleczaNiedostepnegoModeluObrazu(brak.Error())
	}
	return bladPrzetwarzaniaModeluObrazu(err.Error())
}

// bladZapleczaNiedostepnegoModeluObrazu znakuje zaplecze niedostępne: żądanie
// było poprawne, brakuje czegoś w instalacji i komunikat mówi czego. Ten sam
// kod, co przy braku ImageMagicka i braku silnika mowy.
func bladZapleczaNiedostepnegoModeluObrazu(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
		"silniki obrazu: "+powod))
}

// bladWskazaniaModeluObrazu nazywa brak albo niepoprawność wskazania
// w żądaniu — usterka wołającego, nie rdzenia.
func bladWskazaniaModeluObrazu(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"silniki obrazu: "+powod))
}

// bladPrzetwarzaniaModeluObrazu znakuje przebieg silnika, który ruszył, ale się nie udał —
// usterka w trakcie pracy, nie w żądaniu ani w instalacji.
func bladPrzetwarzaniaModeluObrazu(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"silniki obrazu: "+powod))
}

// bladZapleczaModeluObrazu znakuje awarię po stronie rdzenia: katalog pracy
// nie do założenia, nośnik pełny, magazyn niewpięty.
func bladZapleczaModeluObrazu(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"silniki obrazu: "+powod))
}
