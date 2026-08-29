// Zaplecze wspólne czterech narzędzi obrazu modelu — image.inspect,
// image.transform, image.adjust, image.convert — rozwiązuje źródło żądania
// i odkłada wynik jako nowy zasób Designu.
package core

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/session"
	"danacoconsole/server/internal/zewnetrzne"
	"danacoconsole/shared"
)

const (
	// granicaNarzedziObrazu ogranicza jedno uruchomienie programu na drodze
	// zapasowej: dwie minuty starczają na materiał aparatowy każdej rozsądnej
	// wielkości, a plik uszkodzony nie zawiesza procesu bez końca.
	granicaNarzedziObrazu = 120 * time.Second
)

// narzedzieImageMagick opisuje binarium przetwarzające — tryb `magick` (wersja
// siódma). Nazwa czytelna i pakiet wchodzą do treści odmowy, żeby Operator
// przeczytał, czego brakuje i czym to dociągnąć.
func narzedzieImageMagick(program string) zewnetrzne.Narzedzie {
	return zewnetrzne.Narzedzie{
		Nazwa:   "ImageMagick",
		Program: program,
		Pakiet:  "imagemagick",
	}
}

// adapterNarzedziObrazu wypełnia port `NarzedziaObrazu`: repozytorium Designu
// daje odczyt i zapis zasobu, magazyn trzyma bajty, uruchamiacz startuje
// proces. Adapter nie trzyma stanu między wywołaniami, bo obraz nie jest sesją.
type adapterNarzedziObrazu struct {
	repozytorium dane.RepozytoriumDesignu
	magazyn      *magazynTresciBiblioteki
	uruchamiacz  session.Uruchamiacz
	rozstrzygacz *konfig.Rozstrzygacz
	katalog      *KatalogRoboczy
}

// nowyAdapterNarzedziObrazu wiąże port z repozytorium modułu Design i magazynem
// jego zasobów. Uruchamiacz oraz izolacja wchodzą osobno metodą ZArsenalem, bo
// montaż zna je dopiero po złożeniu warstwy kanału.
func nowyAdapterNarzedziObrazu(repozytorium dane.RepozytoriumDesignu,
	katalogDanych string) *adapterNarzedziObrazu {

	return &adapterNarzedziObrazu{
		repozytorium: repozytorium,
		magazyn:      magazynZasobowDesignu(katalogDanych),
	}
}

// ZArsenalem wpina drogę do binariów: uruchamiacz procesów i te same dwa źródła
// izolacji, którymi jadą Terminal, Developer oraz silniki mowy. Brak zależności
// nie psuje montażu — odmawia sama droga zapasowa, nazywając brak zamiast
// udawać przetworzenie.
func (a *adapterNarzedziObrazu) ZArsenalem(uruchamiacz session.Uruchamiacz,
	rozstrzygacz *konfig.Rozstrzygacz, katalog *KatalogRoboczy) *adapterNarzedziObrazu {

	a.uruchamiacz = uruchamiacz
	a.rozstrzygacz = rozstrzygacz
	a.katalog = katalog
	return a
}

// zrodloObrazu jest rozwiązanym wskazaniem modelu: gdzie leżą bajty, ile ich
// jest i do którego okna należy wynik.
type zrodloObrazu struct {
	// sciezka wskazuje plik do odczytu — blob magazynu albo plik na dysku.
	sciezka string
	// rozmiar jest miarą źródła w bajtach; image.convert liczy z niego oszczędność.
	rozmiar int
	// okno jest oknem, w którym ma wylądować wynik.
	okno string
	// nazwa źródła — wchodzi w nazwę zasobu powstałego, żeby model widział
	// w panelu, z czego co powstało.
	nazwa *string
}

// rozwiazZrodlo przekłada parę assetId i sourcePath na plik do odczytu. Zasób
// ma pierwszeństwo przed ścieżką, bo jego treść leży pod sumą kontrolną.
// Żądanie bez obu wskazań jest odmową, nie zgadywaniem przedmiotu czynności.
func (a *adapterNarzedziObrazu) rozwiazZrodlo(ctx context.Context,
	kodZasobu, sciezkaZrodlowa *string) (zrodloObrazu, error) {

	if !bezWartosci(kodZasobu) {
		return a.zrodloZZasobu(ctx, strings.TrimSpace(*kodZasobu))
	}
	if !bezWartosci(sciezkaZrodlowa) {
		return a.zrodloZePliku(strings.TrimSpace(*sciezkaZrodlowa))
	}
	return zrodloObrazu{}, bladWskazaniaObrazu(
		"żądanie bez wskazania obrazu: brakuje pola assetId albo sourcePath — " +
			"narzędzie obrazu nie zgaduje, na czym ma pracować")
}

// zrodloZZasobu odczytuje wiersz zasobu i sprawdza, że jego bajty leżą tam,
// gdzie wskazuje pole uri. Wiersz bez treści na dysku jest odmową, nie pustym
// obrazem: zasób, którego blob zniknął, mówi o tym wprost.
func (a *adapterNarzedziObrazu) zrodloZZasobu(ctx context.Context, kod string) (zrodloObrazu, error) {
	if a.repozytorium == nil {
		return zrodloObrazu{}, bladZapleczaObrazu(
			"serwer nie ma magazynu zasobów — nie ma gdzie szukać wskazanego obrazu")
	}
	if kod == "" {
		return zrodloObrazu{}, bladWskazaniaObrazu("wskazanie zasobu jest puste")
	}
	zasob, err := a.repozytorium.Zasob(ctx, kod)
	if err != nil {
		if errors.Is(err, dane.ErrBrakWiersza) {
			return zrodloObrazu{}, bladWskazaniaObrazu("nie ma zasobu o identyfikatorze " + kod)
		}
		return zrodloObrazu{}, bladZapleczaObrazu("nie można odczytać zasobu " + kod + ": " + err.Error())
	}
	if zasob.URI == nil || strings.TrimSpace(*zasob.URI) == "" {
		return zrodloObrazu{}, bladWskazaniaObrazu("zasób " + kod +
			" nie ma treści (puste uri) — nie ma czego przetworzyć")
	}
	opis, err := os.Stat(*zasob.URI)
	if err != nil {
		return zrodloObrazu{}, bladZapleczaObrazu("treść zasobu " + kod +
			" nie leży pod zapisanym odwołaniem: " + err.Error())
	}
	return zrodloObrazu{
		sciezka: *zasob.URI,
		rozmiar: int(opis.Size()),
		okno:    zasob.Okno,
		nazwa:   zasob.Nazwa,
	}, nil
}

// zrodloZePliku bierze plik wskazany ścieżką. Bajtów nie wciąga do magazynu,
// bo źródło ma zostać nietknięte i niepowielone — kopiowanie zostawiłoby
// w magazynie drugi egzemplarz zdjęcia, o który nikt nie prosił.
func (a *adapterNarzedziObrazu) zrodloZePliku(sciezka string) (zrodloObrazu, error) {
	if sciezka == "" {
		return zrodloObrazu{}, bladWskazaniaObrazu("wskazana ścieżka jest pusta")
	}
	opis, err := os.Stat(sciezka)
	if err != nil {
		return zrodloObrazu{}, bladWskazaniaObrazu("nie można odczytać pliku " + sciezka + ": " + err.Error())
	}
	if opis.IsDir() {
		return zrodloObrazu{}, bladWskazaniaObrazu(sciezka + " jest katalogiem, nie obrazem")
	}
	nazwa := nazwaPlikuZrodla(sciezka)
	// Okno zostaje puste: plik na dysku nie ma okna, więc wynik zostaje bez wiersza.
	return zrodloObrazu{
		sciezka: sciezka,
		rozmiar: int(opis.Size()),
		okno:    "",
		nazwa:   &nazwa,
	}, nil
}

// wolajImageMagick przeprowadza jedno uruchomienie i oddaje wyjście programu:
// bajty obrazu przy przetwarzaniu, tekst opisu przy badaniu. Tryb wybiera
// zamiast zakładać, bo siódma wersja stoi plikiem magick, a szósta osobnymi
// convert i identify.
func (a *adapterNarzedziObrazu) wolajImageMagick(ctx context.Context,
	podpolecenie string, argumenty []string, zapasowy string) ([]byte, error) {

	if a.uruchamiacz == nil {
		return nil, bladZapleczaNiedostepnego("serwer nie ma uruchamiacza procesów — " +
			"narzędzia obrazu nie mają czym wystartować; " +
			"naprawa: podpiąć warstwę kanału (injection) przy składaniu serwera")
	}

	program := a.wybierzProgram(zapasowy)
	if podpolecenie != "" && program == "magick" {
		argumenty = append([]string{podpolecenie}, argumenty...)
	}
	narzedzie := narzedzieImageMagick(program)
	okno, zasady, obszar := a.zasiegNarzedzi()

	wynik, err := zewnetrzne.Wolaj(ctx, a.uruchamiacz, okno, zasady, obszar,
		narzedzie, argumenty, strings.TrimSpace(obszar.KatalogRoboczy), granicaNarzedziObrazu)
	if err != nil {
		return nil, bladArsenaluObrazu(err)
	}
	return wynik.Wyjscie, nil
}

// wybierzProgram rozstrzyga, którym plikiem wykonywalnym wołać ImageMagicka.
// Parametr zapasowy jest odpowiednikiem z wersji szóstej dla danej czynności —
// convert dla przetwarzania, identify dla badania.
func (a *adapterNarzedziObrazu) wybierzProgram(zapasowy string) string {
	if zewnetrzne.Stoi(narzedzieImageMagick("magick")) {
		return "magick"
	}
	if zapasowy != "" && zewnetrzne.Stoi(narzedzieImageMagick(zapasowy)) {
		return zapasowy
	}
	// Nic nie stoi — oddaje się nazwę trybu siódmego, żeby odmowa niosła
	// właściwą wskazówkę instalacyjną.
	return "magick"
}

// zasiegNarzedzi składa trójkę okno–zasady–obszar dla zasięgu platformy:
// żądanie image niesie sam obraz, a nie okno rozmowy, więc adresem jest
// najszerszy poziom zasięgu. Okno dostaje ExecutionEnvCore jawnie, bo
// przetwarzanie idzie na dysku rdzenia.
func (a *adapterNarzedziObrazu) zasiegNarzedzi() (session.Okno, session.Zasady, session.Obszar) {
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

// odlozZasob utrwala bajty wyniku i zakłada jego wiersz tą samą drogą, którą
// idzie design.asset.upload. Bajty zapisuje przed wierszem, bo wiersz bez
// utrwalonej treści pokazywałby w Assets Panel kafelek, za którym nie ma nic.
func (a *adapterNarzedziObrazu) odlozZasob(ctx context.Context, zrodlo zrodloObrazu,
	oknoZadane *string, bajty []byte, formatDocelowy, czynnosc string) (shared.DesignAsset, int, error) {

	if a.magazyn == nil {
		return shared.DesignAsset{}, 0, bladZapleczaObrazu(
			"magazyn treści nie jest wpięty — nie ma gdzie odłożyć bajtów wyniku")
	}
	if len(bajty) == 0 {
		// Rachunek albo program zakończył się powodzeniem i nie oddał ani bajtu —
		// to odmowa udająca sukces.
		return shared.DesignAsset{}, 0, bladPrzetwarzaniaObrazu(
			"przetworzenie zakończyło się powodzeniem, ale nie oddało ani bajtu obrazu")
	}
	suma := sha256.Sum256(bajty)
	odwolanie, err := a.magazyn.Zapisz(bajty, hex.EncodeToString(suma[:]))
	if err != nil {
		return shared.DesignAsset{}, 0, bladZapleczaObrazu(
			"nie można utrwalić treści wyniku: " + err.Error())
	}

	format, szerokosc, wysokosc := rozpoznajObrazZasobu(odwolanie)
	// Format zmierzony bije zadeklarowany, ale AVIF nie ma dekodera nagłówka —
	// zostaje format zadany.
	if format == nil && strings.TrimSpace(formatDocelowy) != "" {
		nazwaFormatu := strings.ToLower(strings.TrimSpace(formatDocelowy))
		format = &nazwaFormatu
	}

	nazwaFormatuWyniku := ""
	if format != nil {
		nazwaFormatuWyniku = *format
	}
	zasob, err := odlozWynikArsenalu(ctx, a.repozytorium, wynikArsenalu{
		odwolanie: odwolanie,
		okno:      oknoWynikuArsenalu(oknoZadane, zrodlo.okno),
		nazwa:     nazwaZasobuWyniku(zrodlo.nazwa, czynnosc),
		format:    nazwaFormatuWyniku,
		szerokosc: szerokosc,
		wysokosc:  wysokosc,
	}, bladZapleczaObrazu)
	if err != nil {
		return shared.DesignAsset{}, 0, err
	}
	return zasob, len(bajty), nil
}

// nazwaZasobuWyniku składa nazwę czytelną dla Operatora, na przykład
// portret.png i resize. Nazwa źródła nieznana zostawia samą czynność, bo pusty
// tekst byłby kafelkiem bez podpisu.
func nazwaZasobuWyniku(nazwaZrodla *string, czynnosc string) string {
	podstawa := ""
	if nazwaZrodla != nil {
		podstawa = strings.TrimSpace(*nazwaZrodla)
	}
	if podstawa == "" {
		return czynnosc
	}
	return podstawa + " · " + czynnosc
}

// nazwaPlikuZrodla wycina samą nazwę pliku ze ścieżki — bez katalogów, których
// Operator w podpisie kafelka nie potrzebuje.
func nazwaPlikuZrodla(sciezka string) string {
	sciezka = strings.TrimRight(sciezka, "/\\")
	if i := strings.LastIndexAny(sciezka, "/\\"); i >= 0 {
		return sciezka[i+1:]
	}
	return sciezka
}

// bladArsenaluObrazu przekłada odmowę pakietu zewnetrzne na kod kontraktu.
// Brak binarium dostaje inny kod niż niepowodzenie programu, bo brak usuwa się
// jedną instalacją, a wywrócony program jest usterką przetwarzania.
func bladArsenaluObrazu(err error) error {
	var brak *zewnetrzne.BrakNarzedzia
	if errors.As(err, &brak) {
		return bladZapleczaNiedostepnego(brak.Error())
	}
	return bladPrzetwarzaniaObrazu(err.Error())
}

// bladZapleczaNiedostepnego znakuje zaplecze niedostępne: żądanie było
// poprawne, brakuje czegoś w instalacji i komunikat mówi czego. Ten sam kod, co
// przy braku silnika mowy.
func bladZapleczaNiedostepnego(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
		"narzędzia obrazu: "+powod))
}

// bladWskazaniaObrazu nazywa brak albo niepoprawność wskazania obrazu
// w żądaniu modelu — usterkę wołającego, a nie rdzenia.
func bladWskazaniaObrazu(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"narzędzia obrazu: "+powod))
}

// bladPrzetwarzaniaObrazu znakuje przetwarzanie obrazu, które ruszyło
// i zakończyło się niepowodzeniem po stronie programu zewnętrznego.
func bladPrzetwarzaniaObrazu(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"narzędzia obrazu: "+powod))
}

// bladZapleczaObrazu znakuje awarię po stronie rdzenia: magazyn niewpięty,
// nośnik pełny, wiersz nie do zapisania, blob zniknął spod odwołania.
func bladZapleczaObrazu(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"narzędzia obrazu: "+powod))
}
