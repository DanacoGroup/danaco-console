// Wspólne zaplecze dwóch narzędzi modelu, które nie są przekształceniem
// obrazu, tylko uruchomieniem modelu nad obrazem — `image.upscale`
// (superrozdzielczość) i `image.background.remove` (wycięcie tła). Same
// czynności leżą w `adapter_narzedzia_obraz_model_silniki.go`, port i wpięcie
// w `adapter_narzedzia_obraz_model_port.go`.
//
// To osobna rodzina od czynności `image.*`, bo obie zmyślają szczegół, którego
// w źródle nie ma (piksele między pikselami, granicę obiektu i tła), i obie
// potrzebują do tego sieci neuronowej z wagami na dysku. Wspólny z tamtą
// rodziną zostaje mechanizm: rozwiązanie pary `assetId?|sourcePath?` na plik
// do odczytu, zasięg izolacji dla wołania binarium i odłożenie bajtów wyniku
// w magazynie zasobów Designu. Trzymamy je przez `wspolne`, a nie kopiujemy,
// żeby druga kopia reguły nie rozjechała się z pierwszą.
//
// Bez silnika na maszynie ta rodzina odmawia, nazywając brak
// (`*zewnetrzne.BrakNarzedzia` → `channel_unavailable`). Podstawienie
// `magick -resize 200%` oddałoby obraz dwa razy większy i ani o szczegół
// bogatszy. Tak samo z tłem: brak silnika to odmowa, a nie obraz bez zmian
// podany jako wycięty.
//
// Wagi modelu leżą obok silnika i nie ściągają się w trakcie żądania. `rembg`
// przy pierwszym uruchomieniu pobiera model sam, a wtedy sto siedemdziesiąt
// sześć megabajtów wchodzi w czas jednego żądania, którego granica jest
// liczona na przetwarzanie, nie na łącze; przy wolnym łączu albo braku sieci
// żądanie urywa się w połowie pobierania i wygląda jak usterka silnika. Wagi
// sprawdzamy więc przed uruchomieniem, a przy ich braku odmawiamy, podając,
// gdzie mają leżeć i ile ważą.
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
	// granicaPowiekszenia to granica czasu jednego przebiegu
	// superrozdzielczości. Bez karty graficznej sieć liczy się na procesorze:
	// obraz 512×512 powiększany dwukrotnie zajmuje ponad minutę, a materiał
	// aparatowy (4000 pikseli szerokości) rośnie z tym kwadratowo. Piętnaście
	// minut znaczy „coś stanęło", a nie „to długo trwa"; niższa granica
	// ucinałaby pracę udaną w połowie i oddawała plik urwany.
	granicaPowiekszenia = 15 * time.Minute

	// granicaWycinaniaTla jest krótsza, bo przebieg jest jeden i stały:
	// obraz idzie do sieci przeskalowany do 320×320, więc czas prawie nie
	// zależy od wielkości źródła. Pięć minut mieści start interpretera,
	// wczytanie wag i przebieg z ogromnym zapasem.
	granicaWycinaniaTla = 5 * time.Minute

	// katalogModeliPowiekszeniaLinux i katalogWagWycinaniaLinux to miejsca wag
	// na serwerze — arsenał pakietu Linux stawia je pod `/usr/local/share`
	// (`scripts/arsenal-serwera.sh`). Na Windowsie tego drzewa nie ma, więc
	// ścieżkę składa się względem pliku wykonywalnego — patrz `katalogModeli...`
	// i `katalogWag...` (funkcje niżej).
	katalogModeliPowiekszeniaLinux = "/usr/local/share/realesrgan-ncnn-vulkan/models"
	katalogWagWycinaniaLinux       = "/usr/local/share/rembg-modele"
)

// katalogModeliPowiekszenia oddaje miejsce, w którym leżą wagi sieci
// powiększającej. Silnik ncnn szuka ich domyślnie w katalogu `models` względem
// katalogu bieżącego, a katalog bieżący wołania arsenału jest obszarem okna —
// czyli za każdym razem innym. Ścieżka podana jawnie (`-m`) jest jedyną, która
// od tego nie zależy.
//
// Zależy od systemu, bo wagi są SKŁADNIKIEM PAKIETU, a pakiet jest inny na
// serwerze i inny w wersji natywnej Windows. Na Linuksie arsenał serwera stawia
// je pod `/usr/local/share`; w wersji natywnej Windows jadą obok rdzenia
// w `pomocniki/`, więc ścieżkę bezwzględną liczy się dopiero w czasie pracy,
// względem pliku wykonywalnego. Ścieżka zaszyta po linuksowemu odmawiała na
// Windowsie ZANIM doszło do wołania silnika, choćby wagi były w paczce.
func katalogModeliPowiekszenia() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(katalogWagWindows(), "realesrgan-ncnn-vulkan", "models")
	}
	return katalogModeliPowiekszeniaLinux
}

// katalogWagWycinania oddaje miejsce wag sieci wycinającej tło — po tej samej
// zasadzie co wagi powiększania. Na Linuksie ta sama ścieżka stoi w opakowaniu
// `/usr/local/bin/rembg` jako `U2NET_HOME`, bo `zewnetrzne.Wolaj` nie dziedziczy
// środowiska; na Windowsie opakowanie `pomocniki/rembg/rembg.exe` ustawia
// `U2NET_HOME` na ten sam katalog, który zwraca ta funkcja.
func katalogWagWycinania() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(katalogWagWindows(), "rembg", "modele")
	}
	return katalogWagWycinaniaLinux
}

// katalogWagWindows składa katalog `pomocniki` obok pliku wykonywalnego rdzenia
// — tą samą drogą, którą `zewnetrzne` odnajduje programy arsenału. Gdy miejsca
// procesu nie da się ustalić, zostaje ścieżka względna: proces i tak startuje
// z katalogu rdzenia w wydaniu natywnym.
func katalogWagWindows() string {
	if biezace, err := os.Executable(); err == nil {
		return filepath.Join(filepath.Dir(biezace), "pomocniki")
	}
	return "pomocniki"
}

// adapterNarzedziObrazuModelu wypełnia port `NarzedziaObrazuModelu`. Nie ma
// własnego stanu: źródło, zasięg izolacji i magazyn wyniku bierze z adaptera
// rodziny `image.*`, a dokłada wyłącznie to, czego tamten nie zna — dwa
// silniki neuronowe, ich wagi i granice czasu liczone w minutach.
type adapterNarzedziObrazuModelu struct {
	wspolne *adapterNarzedziObrazu
}

// nowyAdapterNarzedziObrazuModelu składa adapter na tym samym zapleczu, na
// którym stoi rodzina `image.*`. Zaplecze przychodzi z zewnątrz, bo instancja
// zbudowana tutaj miałaby własny uchwyt magazynu i własny uruchamiacz, czyli
// drugą prawdę o tym, gdzie rdzeń odkłada bajty. Brak zaplecza nie psuje
// montażu — psuje dwie komendy, które wtedy odmawiają, nazywając brak.
func nowyAdapterNarzedziObrazuModelu(wspolne *adapterNarzedziObrazu) *adapterNarzedziObrazuModelu {
	return &adapterNarzedziObrazuModelu{wspolne: wspolne}
}

// pracowniaObrazu jest katalogiem jednego przebiegu: leży w nim dowiązanie do
// źródła i plik wyniku.
//
// Pliki pośrednie są tu konieczne, choć rodzina `image.*` bierze wynik
// ImageMagicka ze standardowego wyjścia: ani `realesrgan-ncnn-vulkan`, ani
// `rembg` nie umieją pisać obrazu na wyjście, oba żądają ścieżki wyniku.
// Katalog własny na przebieg nie miesza równoległych żądań i znika jednym
// `RemoveAll` niezależnie od tego, czy silnik się udał.
type pracowniaObrazu struct {
	katalog string
	wejscie string
	wyjscie string
}

// przygotujPracownie zakłada katalog przebiegu i wystawia w nim źródło pod
// nazwą, którą silnik przyjmie.
//
// Źródło wchodzi dowiązaniem, nie kopią: blob zasobu leży pod swoją sumą
// kontrolną i bywa wielkim plikiem, więc kopiowanie go tylko po to, żeby
// zmienić nazwę, dawałoby drugi egzemplarz zdjęcia przy każdym żądaniu. Oba
// silniki rozpoznają format po zawartości, a nie po rozszerzeniu, więc nazwa
// `zrodlo.png` jest tylko uchwytem, nie deklaracją formatu.
//
// Kopia jest drogą zapasową na wypadek, gdy dowiązanie się nie uda (magazyn na
// innym nośniku niż katalog tymczasowy, system plików bez dowiązań).
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

// odczytajWynikSilnika bierze bajty pliku, który silnik miał wytworzyć.
//
// Plik nieobecny albo pusty jest odmową mimo zerowego kodu wyjścia. Oba
// silniki potrafią zakończyć się powodzeniem i nie napisać nic — `rembg` przy
// wagach uszkodzonych, ncnn przy nieudanej alokacji pamięci na procesorze
// programowym — a oddanie pustego zasobu byłoby sukcesem bez skutku.
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

// wolajSilnik przeprowadza jedno uruchomienie silnika neuronowego wspólną
// drogą wołania binarium (`zewnetrzne.Wolaj`): przez port
// `session.Uruchamiacz`, bramę izolacji okna i objęcie drzewa procesów.
//
// Zasięg bierzemy z rodziny `image.*`, bo jest ten sam: żądanie niesie obraz,
// a nie okno rozmowy, więc adresem jest zasięg platformy, a wykonanie idzie na
// dysku rdzenia, bo to rdzeń odkłada potem bajty do swojego magazynu.
func (a *adapterNarzedziObrazuModelu) wolajSilnik(ctx context.Context,
	narzedzie zewnetrzne.Narzedzie, argumenty []string, limit time.Duration) error {

	if a.wspolne == nil || a.wspolne.uruchamiacz == nil {
		return bladZapleczaNiedostepnegoModeluObrazu(
			"rdzeń nie ma uruchamiacza procesów — silniki obrazu nie mają czym wystartować; " +
				"naprawa: podpiąć warstwę kanału (injection) przy składaniu rdzenia")
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
		"(rdzeń świadomie nie ściąga ich w trakcie żądania, żeby czas łącza " +
		"nie zjadał granicy czasu przetwarzania)")
}

// wymiaryWyniku wyciąga zmierzone wymiary z odłożonego zasobu.
//
// Brak wymiarów jest odmową, a nie zerem w odpowiedzi: kontrakt
// `image.upscale` obiecuje `width` i `height` jako liczby, a odpowiedź
// z zerami mówiłaby, że powstał obraz bez pikseli. Wymiary pochodzą z pomiaru
// pliku, który legł w magazynie (`rozpoznajObrazZasobu`), a nie z pomnożenia
// żądanej krotności przez wymiar źródła — przy zaokrągleniach silnika te dwie
// liczby się różnią, a zgadnięta wygląda w odpowiedzi jak zmierzona.
func wymiaryWyniku(zasob shared.DesignAsset) (int, int, error) {
	if zasob.Width == nil || zasob.Height == nil {
		return 0, 0, bladPrzetwarzaniaModeluObrazu(
			"wynik legł w magazynie, ale nie dało się zmierzyć jego wymiarów — " +
				"odpowiedź nie poda liczb, których nie zmierzono")
	}
	return *zasob.Width, *zasob.Height, nil
}

// bladArsenaluModeluObrazu przekłada odmowę pakietu `zewnetrzne` na kod
// kontraktu. Brak silnika dostaje inny kod niż niepowodzenie silnika, bo „nie
// ma czym" jest jedynym zakończeniem żądania bez zainstalowanego modelu —
// alternatywą byłoby rozciągnięcie podane jako powiększenie.
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

// bladPrzetwarzaniaModeluObrazu znakuje przebieg, który ruszył i się nie udał.
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
