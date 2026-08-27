// Wspólny trzon rodziny narzędzi `media.*`: dwie komendy modelu dzielą
// źródło bajtów, zasięg izolacji, wołanie binarium i magazyn wyniku, a wynik
// jest zawsze nowym zasobem, źródło zostaje nietknięte.
package core

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/konfiguracja"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/session"
	"danacoconsole/server/internal/zewnetrzne"
	"danacoconsole/shared"
)

// Dwa binaria arsenału. Rozdzielone, bo rozdzielone są ich braki: maszyna bez
// `ffprobe` nie zmierzy materiału, ale bez `ffmpeg` nie przetworzy go — a
// Operator ma przeczytać w odmowie nazwę tego programu, którego naprawdę
// zabrakło.
var (
	narzedzieFfprobe = zewnetrzne.Narzedzie{
		Nazwa: "ffprobe", Program: "ffprobe", Pakiet: "ffmpeg",
	}
	narzedzieFfmpeg = zewnetrzne.Narzedzie{
		Nazwa: "ffmpeg", Program: "ffmpeg", Pakiet: "ffmpeg",
	}
)

const (
	// Granica pomiaru jest krótka, bo `ffprobe` czyta nagłówki, a nie materiał.
	// Przekroczenie takiej granicy znaczy plik uszkodzony albo nośnik, który
	// nie oddaje bajtów — jedno i drugie ma się skończyć odmową, nie czekaniem.
	granicaPomiaruMediow = 60 * time.Second

	// Granica przetworzenia jest hojniejsza niż dla obrazu, bo film przelicza
	// każdą klatkę i minuta materiału bywa minutą pracy. Granica jednak jest,
	// bo `ffmpeg` na uszkodzonym strumieniu nie kończy się nigdy.
	granicaPrzetworzeniaMediow = 30 * time.Minute
)

// adapterNarzedziMediow wypełnia port NarzedziaMedia i trzyma uruchamiacz
// procesów, repozytorium zasobów, magazyn treści oraz reguły izolacji
// rodziny narzędzi mediów.
type adapterNarzedziMediow struct {
	// uruchamiacz jest jedyną drogą startu procesu; bez niego rodzina nie
	// ruszy ani ffprobe, ani ffmpeg.
	uruchamiacz session.Uruchamiacz
	// repozytorium daje odczyt zasobu wskazanego assetId i zapis wiersza
	// zasobu wynikowego.
	repozytorium dane.RepozytoriumDesignu
	// magazyn trzyma bajty — ten sam magazyn zasobów Designu, którym jedzie
	// wynik media.transcode.
	magazyn *magazynTresciBiblioteki
	// rozstrzygacz i katalog składają zasady izolacji oraz obszar zasięgu
	// tej rodziny narzędzi.
	rozstrzygacz *konfig.Rozstrzygacz
	katalog      *KatalogRoboczy
}

// nowyAdapterNarzedziMediow wiąże rodzinę z uruchamiaczem procesów,
// repozytorium zasobów i katalogiem danych rdzenia. Pusty katalog danych
// spada na domyślny, żeby konstruktor nigdy nie oddał adaptera bez magazynu.
func nowyAdapterNarzedziMediow(uruchamiacz session.Uruchamiacz,
	repozytorium dane.RepozytoriumDesignu, katalogDanych string) *adapterNarzedziMediow {

	if strings.TrimSpace(katalogDanych) == "" {
		katalogDanych = konfiguracja.KatalogDanychDomyslny()
	}
	return &adapterNarzedziMediow{
		uruchamiacz:  uruchamiacz,
		repozytorium: repozytorium,
		magazyn:      magazynZasobowDesignu(katalogDanych),
	}
}

// ZIzolacja podpina rozstrzygacz zasięgu i ustalacz katalogu roboczego — dwa
// źródła, z których powstają zasady i obszar egzekwowane przy uruchomieniu.
func (a *adapterNarzedziMediow) ZIzolacja(rozstrzygacz *konfig.Rozstrzygacz,
	katalog *KatalogRoboczy) *adapterNarzedziMediow {

	a.rozstrzygacz, a.katalog = rozstrzygacz, katalog
	return a
}

// zrodloMediow rozstrzyga, skąd wziąć bajty materiału, i oddaje ścieżkę
// pliku wraz z oknem źródła. Pierwszeństwo ma zasób wskazany `assetId`,
// drugą drogą jest `sourcePath` wciągnięty do magazynu.
func (a *adapterNarzedziMediow) zrodloMediow(ctx context.Context,
	komenda string, assetId, sourcePath *string) (string, string, error) {

	if !bezWartosci(assetId) {
		if a.repozytorium == nil {
			return "", "", odmowaNarzedziMediow(shared.ErrorCodeInternalError,
				"rejestr zasobów nie jest wpięty — nie ma gdzie odczytać zasobu "+*assetId)
		}
		zasob, err := a.repozytorium.Zasob(ctx, strings.TrimSpace(*assetId))
		if err != nil {
			return "", "", odmowaNarzedziMediow(shared.ErrorCodeNotFound,
				"nie ma zasobu o identyfikatorze "+*assetId+": "+err.Error())
		}
		if zasob.URI == nil || strings.TrimSpace(*zasob.URI) == "" {
			return "", "", odmowaNarzedziMediow(shared.ErrorCodeNotFound,
				"zasób "+*assetId+" nie ma odnośnika do treści — nie ma czego zmierzyć")
		}
		if _, err := os.Stat(*zasob.URI); err != nil {
			return "", "", odmowaNarzedziMediow(shared.ErrorCodeNotFound,
				"treść zasobu "+*assetId+" nie leży pod odnośnikiem "+*zasob.URI+
					": "+err.Error())
		}
		return *zasob.URI, zasob.Okno, nil
	}

	if !bezWartosci(sourcePath) {
		if a.magazyn == nil {
			return "", "", odmowaNarzedziMediow(shared.ErrorCodeInternalError,
				"magazyn treści nie jest wpięty — nie ma gdzie wciągnąć bajtów spod ścieżki")
		}
		odwolanie, _, _, err := a.magazyn.ZapiszZePliku(strings.TrimSpace(*sourcePath))
		if err != nil {
			return "", "", odmowaNarzedziMediow(shared.ErrorCodeValidationFailed,
				"nie można wciągnąć treści spod ścieżki "+*sourcePath+": "+err.Error())
		}
		return odwolanie, "", nil
	}

	return "", "", odmowaNarzedziMediow(shared.ErrorCodeValidationFailed,
		"komenda "+komenda+" bez wskazania materiału: brakuje pola assetId albo "+
			"sourcePath — nie ma czego wziąć do przetworzenia")
}

// odlozWynikMediow utrwala plik wytworzony przez `ffmpeg` i oddaje zasób
// kontraktu. Kolejność jest zamierzona — najpierw bajty, potem wiersz — a tę
// funkcję dopełnia wspólny `odlozWynikArsenalu` (`adapter_narzedzia_wynik.go`).
func (a *adapterNarzedziMediow) odlozWynikMediow(ctx context.Context, sciezka,
	oknoWyniku, nazwa, format string, szerokosc, wysokosc *int) (shared.DesignAsset, error) {

	if a.magazyn == nil {
		return shared.DesignAsset{}, odmowaNarzedziMediow(shared.ErrorCodeInternalError,
			"magazyn treści nie jest wpięty — nie ma gdzie odłożyć bajtów wyniku")
	}

	odwolanie, _, _, err := a.magazyn.ZapiszZePliku(sciezka)
	if err != nil {
		return shared.DesignAsset{}, odmowaNarzedziMediow(shared.ErrorCodeInternalError,
			"nie można utrwalić wyniku przetworzenia: "+err.Error())
	}

	return odlozWynikArsenalu(ctx, a.repozytorium, wynikArsenalu{
		odwolanie: odwolanie,
		okno:      oknoWyniku,
		nazwa:     nazwa,
		format:    format,
		szerokosc: szerokosc,
		wysokosc:  wysokosc,
	}, func(powod string) error {
		return odmowaNarzedziMediow(shared.ErrorCodeInternalError, powod)
	})
}

// wolajMediow prowadzi jedno uruchomienie binarium arsenału. Katalog
// uruchomienia zostaje pusty z zamysłem: wszystkie ścieżki w argumentach są
// bezwzględne, więc rozstrzygnięcie katalogu roboczego zostawia się bramie
// izolacji.
func (a *adapterNarzedziMediow) wolajMediow(ctx context.Context,
	narzedzie zewnetrzne.Narzedzie, argumenty []string,
	granica time.Duration) (zewnetrzne.Wynik, error) {

	okno, zasady, obszar := a.zasiegNarzedziMediow()
	return zewnetrzne.Wolaj(ctx, a.uruchamiacz, okno, zasady, obszar,
		narzedzie, argumenty, "", granica)
}

// zasiegNarzedziMediow składa trójkę okno–zasady–obszar dla zasięgu
// platformy. Żądania tej rodziny okna nie niosą, bo pytają o zdolność
// maszyny, nie okna; środowisko wykonania jest rdzeniowe wprost, bo materiał
// leży na hoście rdzenia.
func (a *adapterNarzedziMediow) zasiegNarzedziMediow() (session.Okno,
	session.Zasady, session.Obszar) {

	okno := session.Okno{Ustawienia: session.Ustawienia{
		SrodowiskoWykonania: shared.ExecutionEnvCore,
	}}
	zasady := session.Zasady{}
	if a.rozstrzygacz != nil {
		zasady = ZasadyIzolacji(a.rozstrzygacz, konfig.Kontekst{})
	}
	obszar := session.Obszar{}
	if a.katalog != nil {
		obszar = ObszarOkna(a.katalog.Ustal(konfig.Kontekst{}, ""), "")
	}
	return okno, zasady, obszar
}

// katalogPrzejsciowyMediow zakłada katalog na plik wynikowy i oddaje go wraz
// ze sprzątaczem. Katalog jest tymczasowy, bo wynik nie jest jeszcze zasobem,
// i ma zniknąć również wtedy, gdy binarium się wywróci.
func katalogPrzejsciowyMediow() (string, func(), error) {
	katalog, err := os.MkdirTemp("", "danaco-media-")
	if err != nil {
		return "", func() {}, odmowaNarzedziMediow(shared.ErrorCodeInternalError,
			"nie można założyć katalogu na wynik przetworzenia: "+err.Error())
	}
	return katalog, func() { _ = os.RemoveAll(katalog) }, nil
}

// rozmiarPlikuMediow mierzy plik na dysku. Nieudany pomiar oddaje zero, a nie
// zgadniętą liczbę: zero rozpoznawalnie znaczy „nie zmierzono", a liczba
// wymyślona wygląda w wyniku identycznie jak zmierzona.
func rozmiarPlikuMediow(sciezka string) int {
	stan, err := os.Stat(sciezka)
	if err != nil {
		return 0
	}
	return int(stan.Size())
}

// bladNarzedziMediow przekłada odmowy warstwy uruchomieniowej na kody
// kontraktu: brak binarium na `channel_unavailable`, naruszenie izolacji na
// `permission_denied`, a pozostałe usterki na `internal_error`.
func bladNarzedziMediow(err error) error {
	var brak *zewnetrzne.BrakNarzedzia
	if errors.As(err, &brak) {
		return odmowaNarzedziMediow(shared.ErrorCodeChannelUnavailable, brak.Error())
	}
	if errors.Is(err, session.ErrIzolacja) {
		return odmowaNarzedziMediow(shared.ErrorCodePermissionDenied, err.Error())
	}
	return odmowaNarzedziMediow(shared.ErrorCodeInternalError, err.Error())
}

// odmowaNarzedziMediow składa odmowę rodziny z kodem kontraktu. Przedrostek
// nazywa rodzinę, żeby czytający Errors Panel wiedział, kto odmówił, zanim
// przeczyta dlaczego.
func odmowaNarzedziMediow(kod shared.ErrorCode, powod string) error {
	return protocol.JakoError(protocol.NowyBlad(kod, "narzędzia mediów: "+powod))
}

// rozszerzenieMediow sprowadza wskazanie formatu do samego rozszerzenia pliku.
// Kropka wiodąca i wielkość liter są ozdobą wołającego, nie treścią wskazania;
// `ffmpeg` rozstrzyga kontener po rozszerzeniu nazwy wyjściowej, więc musi ono
// być czyste.
func rozszerzenieMediow(format string) string {
	return strings.ToLower(strings.TrimPrefix(strings.TrimSpace(format), "."))
}

// plikWynikowyMediow składa ścieżkę pliku, który wytworzy binarium,
// w katalogu przejściowym przygotowanym dla tego uruchomienia.
func plikWynikowyMediow(katalog, rozszerzenie string) string {
	return filepath.Join(katalog, "wynik."+rozszerzenie)
}
