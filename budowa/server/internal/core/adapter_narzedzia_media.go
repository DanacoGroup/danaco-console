// Odpowiedzialność pliku: wspólny trzon rodziny narzędzi `media.*` — dwie
// komendy modelu (`media.inspect`, `media.transcode`) stoją na tym samym
// źródle bajtów, tym samym zasięgu izolacji, tym samym wołaniu binarium
// i tym samym magazynie wyniku. Pomiar (`ffprobe`) leży w
// `adapter_narzedzia_media_pomiar.go`, przetworzenie (`ffmpeg`) w
// `adapter_narzedzia_media_przetworzenie.go`, wpięcie komend w
// `handlers_narzedzia_media.go`.
//
// Film i dźwięk są materiałem, którego model nie zmierzy ani nie przetworzy bez
// cudzego programu. Kontrakt wpisuje obie komendy jako narzędzia modelu
// (`danaco_media_inspect`, `danaco_media_transcode`), więc ich wołaczem jest
// model w turze, a nie panel okna.
//
// Binarium wołane jest jedną drogą — `zewnetrzne.Wolaj`. Własnego
// `exec.Command` w tym pliku nie ma i być nie może: tamta droga idzie przez
// port `session.Uruchamiacz`, bramę izolacji okna i objęcie drzewa procesów.
// Ostatnie jest tu ważniejsze niż gdziekolwiek: `ffmpeg` rozgałęzia wątki
// dekodera i filtrów, a przerwane transkodowanie bez objęcia drzewa zostawia
// na maszynie Operatora procesy mielące film w nieskończoność.
//
// Wynik jest zawsze nowym zasobem, źródło zostaje nietknięte. `ffmpeg` pisze
// do pliku w katalogu tymczasowym, a stamtąd bajty wciąga magazyn zasobów pod
// sumę sha256 — ten sam magazyn, którym jedzie `design.asset.upload`
// (`adapter_modul_design_wgranie.go`). Przetworzenie „w miejscu" byłoby
// zniszczeniem materiału Operatora przy pierwszej pomyłce w parametrach.
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
	// każdą klatkę i minuta materiału bywa minutą pracy. Granica jednak jest:
	// uruchomienie bez niej odrzuca sam `zewnetrzne.Wolaj`, bo `ffmpeg` na
	// uszkodzonym strumieniu nie kończy się nigdy.
	granicaPrzetworzeniaMediow = 30 * time.Minute
)

// adapterNarzedziMediow wypełnia port NarzedziaMedia.
type adapterNarzedziMediow struct {
	// uruchamiacz jest portem warstwy kanału — jedyną drogą startu procesu
	// w drzewie. Bez niego rodzina nie ruszy ani `ffprobe`, ani
	// `ffmpeg`, i mówi to wprost zamiast milczeć.
	uruchamiacz session.Uruchamiacz
	// repozytorium daje dwie rzeczy: odczyt zasobu wskazanego `assetId`
	// (źródło bajtów) i zapis wiersza zasobu wynikowego.
	repozytorium dane.RepozytoriumDesignu
	// magazyn trzyma bajty — ten sam magazyn zasobów Designu, bo wynikiem
	// `media.transcode` jest `DesignAsset` i drugi skład bajtów zasobu byłby
	// drugą prawdą o tym, gdzie leży treść.
	magazyn *magazynTresciBiblioteki
	// rozstrzygacz i katalog składają zasady izolacji oraz obszar zasięgu —
	// te same dwa źródła, którymi jadą Terminal, Developer i silnik mowy.
	rozstrzygacz *konfig.Rozstrzygacz
	katalog      *KatalogRoboczy
}

// nowyAdapterNarzedziMediow wiąże rodzinę z uruchamiaczem procesów,
// repozytorium zasobów i katalogiem danych rdzenia.
//
// Pusty katalog danych spada na domyślny, żeby konstruktor nigdy nie oddał
// adaptera bez magazynu — wzorem `nowyAdapterDesignu`. Adapter bez magazynu
// odmawiałby z powodu własnego montażu, a nie z powodu czegokolwiek, co zrobił
// Operator.
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

// zrodloMediow rozstrzyga, skąd wziąć bajty materiału, i oddaje ścieżkę pliku
// gotowego do podania binarium.
//
// Drogi są dwie, a pierwszeństwo ma zasób. `assetId` wskazuje treść już
// utrwaloną pod sumą kontrolną — nikt jej pod narzędziem nie podmieni.
// `sourcePath` wskazuje plik żywy na dysku Operatora, więc bajty są wciągane do
// magazynu (`ZapiszZePliku`), zamiast podawać binarium cudzą ścieżkę: tak samo
// robi wniesienie zasobu i wgranie do biblioteki, i z tego samego powodu —
// pomiar ma opisywać treść, która po pomiarze nadal jest tą samą treścią.
//
// Żądanie bez obu wskazań jest odmową nazywającą brak; domyślanie się materiału
// („weź ostatni zasób") byłoby zgadywaniem.
//
// Okno źródła wraca razem ze ścieżką i to jedyny powód, dla którego ta funkcja
// oddaje parę zamiast jednej wartości. Gdy materiałem jest zasób, jego wiersz
// zna okno, w którym Operator ten materiał widzi — i wynik ma prawo zostać tam
// samo, gdy żądanie nie powie inaczej. Plik wskazany ścieżką nie należy do
// żadnego okna i oddaje pusty tekst; rozstrzyga wtedy `windowId` żądania albo
// nie powstaje wiersz (`oknoWynikuArsenalu`).
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
// kontraktu.
//
// Kolejność jest zamierzona — najpierw bajty, potem wiersz — dokładnie jak przy
// wniesieniu zasobu: wiersz bez treści byłby kafelkiem, za którym nie ma nic.
// Ta funkcja odpowiada wyłącznie za pierwszy krok; drugi robi wspólny
// `odlozWynikArsenalu` (`adapter_narzedzia_wynik.go`).
//
// Rodzaj zasobu nie jest przybliżeniem. Kontrakt zna `video`, `audio`,
// `document` i `archive`, warunek CHECK kolumny je dopuszcza
// (`migracja_113_rodzaje_zasobow_arsenalu.sql`), a rodzaj wyprowadza z formatu
// wyniku wspólna tablica arsenału. Ma to znaczenie właśnie w tej rodzinie:
// `media.transcode` z czynnością `frame` daje obraz, a `extractAudio` — dźwięk,
// choć komenda jest ta sama.
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

// wolajMediow prowadzi jedno uruchomienie binarium arsenału.
//
// Katalog uruchomienia zostaje pusty z zamysłem. Wszystkie ścieżki
// w argumentach są bezwzględne (blob magazynu na wejściu, plik katalogu
// tymczasowego na wyjściu), więc katalog bieżący procesu nie ma na nic wpływu.
// Pusty oddaje rozstrzygnięcie bramie izolacji: przy włączonym punkcie
// „katalog roboczy" brama sama wstawi katalog własny zasięgu, zamiast odrzucić
// katalog tymczasowy leżący poza nim.
func (a *adapterNarzedziMediow) wolajMediow(ctx context.Context,
	narzedzie zewnetrzne.Narzedzie, argumenty []string,
	granica time.Duration) (zewnetrzne.Wynik, error) {

	okno, zasady, obszar := a.zasiegNarzedziMediow()
	return zewnetrzne.Wolaj(ctx, a.uruchamiacz, okno, zasady, obszar,
		narzedzie, argumenty, "", granica)
}

// zasiegNarzedziMediow składa trójkę okno–zasady–obszar dla zasięgu platformy.
//
// Powód jest ten sam co przy silniku mowy (`adapter_modul_mowa.go`): żądania
// tej rodziny okna nie niosą, bo pytają o zdolność maszyny, nie okna. Pusty
// `konfig.Kontekst{}` jest poprawnym adresem najszerszego z poziomów zasięgu,
// a nie podstawieniem pustych struktur po cichu — gdy Operator włączy punkt
// izolacji globalnie, brama zadziała tu tak samo jak dla Terminala.
//
// Środowisko wykonania jest rdzeniowe wprost, bo materiał leży na hoście
// rdzenia — w magazynie zasobów albo w jego katalogu tymczasowym — a binarium
// musi mieć te bajty pod ręką.
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
// ze sprzątaczem.
//
// Katalog jest tymczasowy, bo wynik nie jest jeszcze zasobem. Zasobem staje
// się dopiero po wciągnięciu do magazynu pod sumę kontrolną; plik pośredni ma
// zniknąć również wtedy, gdy binarium się wywróci — inaczej katalog tymczasowy
// Operatora zapełniałby się urwanymi filmami.
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

// bladNarzedziMediow przekłada odmowy warstwy uruchomieniowej na kody kontraktu.
//
// Trzy przypadki, trzy kody:
//
//   - `*zewnetrzne.BrakNarzedzia` → `channel_unavailable`. „Nie ma czym" jest
//     brakiem po stronie instalacji, który Operator usuwa jedną komendą
//     pakietu — nie wadą żądania (`validation_failed` obwiniłby model) i nie
//     usterką rdzenia (`internal_error` kazałby zgłaszać produkt, który działa
//     i właśnie powiedział, czego dołożyć). Kod jest ponawialny i tak ma być:
//     po instalacji `ffmpeg` to samo żądanie przechodzi bez zmiany.
//
//   - `session.ErrIzolacja` → `permission_denied`. Punkt izolacji Operatora
//     zatrzymał uruchomienie; tak samo znakuje je Terminal i silnik mowy,
//     a dwie reguły dla jednej bramy byłyby rozjazdem.
//
//   - reszta → `internal_error`. Granica czasu, wywrócenie binarium, brak
//     uruchamiacza. Treść niesie już to, co program powiedział o sobie sam
//     (`zewnetrzne` dokleja początek diagnostyki), więc nic nie zostaje
//     przemilczane.
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

// plikWynikowyMediow składa ścieżkę pliku, który wytworzy binarium.
func plikWynikowyMediow(katalog, rozszerzenie string) string {
	return filepath.Join(katalog, "wynik."+rozszerzenie)
}
