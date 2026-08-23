// Wspólne zaplecze czterech narzędzi obrazu modelu (`image.inspect`,
// `image.transform`, `image.adjust`, `image.convert`): rozwiązanie źródła
// (zasób magazynu albo ścieżka na dysku) i odłożenie wyniku jako nowego zasobu.
// Same czynności leżą w `adapter_narzedzia_obraz_czynnosci.go`, rachunek na
// pikselach w `adapter_narzedzia_obraz_wkompilowany.go`, port i wpięcie
// w `handlers_narzedzia_obraz.go`.
//
// Pracę wykonuje biblioteka wkompilowana w binarium. Program pakietu serwera
// zostaje drogą zapasową dla dwóch wyjść, których żaden koder czysto-Go nie
// zapisze — AVIF i WEBP stratny — oraz dla pomiaru pliku AVIF, którego nie ma
// czym zdekodować. Wykaz zależności pakietu serwera niesie go z tym właśnie
// zakresem (`zaleznosci_zewnetrzne.go`), a zapora
// `zapora_procesow_rdzenia_test.go` pilnuje, żeby ta droga nie wróciła jako
// droga podstawowa.
//
// Wołanie binarium idzie wyłącznie przez `zewnetrzne.Wolaj`: stamtąd prowadzi
// port `session.Uruchamiacz`, brama izolacji okna i objęcie drzewa procesów.
// Własnego `exec.Command` w tym pliku nie ma — odstępstwo od tej sekwencji
// kończy się wyciekiem procesu albo uchwytu.
//
// Bajty wyniku lądują w tym samym magazynie zasobów, co `design.asset.upload`
// (`magazynZasobowDesignu`, blob pod sumą sha256), a wiersz — w tej samej
// tabeli `zasob_design`. Drugi magazyn byłby drugą prawdą o tym, gdzie rdzeń
// trzyma bajty poza bazą, a Assets Panel przestałby widzieć połowę zasobów,
// które sam wytworzył.
//
// Źródło zostaje nietknięte: wynikiem każdej z trzech czynności zmieniających
// jest nowy zasób. Blob źródła leży pod swoją sumą kontrolną i nikt go tu nie
// otwiera do zapisu — rachunek czyta go, a bajty wyniku składa osobno; na drodze
// zapasowej program dostaje ścieżkę do odczytu, a wynik oddaje na standardowe
// wyjście (`format:-`).
//
// Brak binarium na drodze zapasowej jest odmową nazwaną: `zewnetrzne.Wolaj`
// oddaje wtedy `*zewnetrzne.BrakNarzedzia` niosące nazwę programu i pakiet do
// doinstalowania. Przekładamy je na kod `channel_unavailable` — ten sam, co
// przy braku silnika mowy — a nie na cichy zasób bez zmian, który wyglądałby
// jak udany retusz. Dotyczy to wyłącznie AVIF-a i WEBP-a stratnego: pozostałe
// czynności rodziny nie mają czego zabraknąć.
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
	// zapasowej. Zamiana formatu idzie w ułamku sekundy, ale program karmiony plikiem
	// uszkodzonym albo bombą dekompresyjną nie kończy się nigdy, a
	// `zewnetrzne.Wolaj` granicy niedodatniej nie przyjmuje. Dwie minuty
	// starczają na materiał aparatowy każdej rozsądnej wielkości.
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

// adapterNarzedziObrazu wypełnia port `NarzedziaObrazu`.
//
// Repozytorium Designu daje odczyt zasobu wskazanego przez model i zapis zasobu
// powstałego; magazyn — miejsce na bajty; uruchamiacz wraz z dwoma źródłami
// izolacji — jedyną drogę startu procesu. Adapter nie trzyma stanu między
// wywołaniami, bo obraz nie jest sesją.
type adapterNarzedziObrazu struct {
	repozytorium dane.RepozytoriumDesignu
	magazyn      *magazynTresciBiblioteki
	uruchamiacz  session.Uruchamiacz
	rozstrzygacz *konfig.Rozstrzygacz
	katalog      *KatalogRoboczy
}

// nowyAdapterNarzedziObrazu wiąże port z repozytorium modułu Design i magazynem
// jego zasobów. Uruchamiacz oraz izolacja wchodzą osobno (`ZArsenalem`), bo
// montaż zna je dopiero po złożeniu warstwy kanału.
func nowyAdapterNarzedziObrazu(repozytorium dane.RepozytoriumDesignu,
	katalogDanych string) *adapterNarzedziObrazu {

	return &adapterNarzedziObrazu{
		repozytorium: repozytorium,
		magazyn:      magazynZasobowDesignu(katalogDanych),
	}
}

// ZArsenalem wpina drogę do binariów: uruchamiacz procesów i te same dwa źródła
// izolacji, którymi jadą Terminal, Developer oraz silniki mowy.
//
// Brak tej zależności nie psuje montażu i nie psuje już czterech komend: cztery
// czynności liczą się biblioteką wkompilowaną. Bez uruchamiacza odmawia sama
// droga zapasowa — zapis AVIF-a i WEBP-a stratnego oraz pomiar pliku AVIF —
// i odmawia nazywając brak, zamiast udawać, że obraz przetworzyła.
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
	// rozmiar jest miarą źródła w bajtach; `image.convert` liczy z niego
	// oszczędność, `image.inspect` oddaje go wprost.
	rozmiar int
	// okno jest oknem, w którym ma wylądować wynik.
	okno string
	// nazwa źródła — wchodzi w nazwę zasobu powstałego, żeby model widział
	// w panelu, z czego co powstało.
	nazwa *string
}

// rozwiazZrodlo przekłada parę `assetId?|sourcePath?` na plik do odczytu.
//
// Zasób ma pierwszeństwo przed ścieżką: jego treść leży pod sumą kontrolną
// i nie zmieni się między wskazaniem a odczytem, a plik na dysku jest treścią
// żywą. Gdy model podał oba wskazania, bierzemy pewniejsze.
//
// Żądanie bez jednego i drugiego jest odmową — podstawienie ostatniego zasobu
// byłoby zgadywaniem przedmiotu czynności.
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
// gdzie wskazuje `uri`. Wiersz bez treści na dysku jest odmową, nie pustym
// obrazem: zasób, którego blob zniknął, mówi o tym wprost.
func (a *adapterNarzedziObrazu) zrodloZZasobu(ctx context.Context, kod string) (zrodloObrazu, error) {
	if a.repozytorium == nil {
		return zrodloObrazu{}, bladZapleczaObrazu(
			"rdzeń nie ma magazynu zasobów — nie ma gdzie szukać wskazanego obrazu")
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

// zrodloZePliku bierze plik wskazany ścieżką. Bajtów nie wciągamy do magazynu:
// źródło ma zostać nietknięte i niepowielone, a wynik i tak wyląduje w magazynie
// własnym blobem. Kopiowanie źródła zostawiłoby w magazynie drugi egzemplarz
// zdjęcia, o który nikt nie prosił.
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
	// Okno zostaje puste: plik na dysku nie należy do żadnego okna, więc nie ma
	// stąd czego odziedziczyć. Puste okno znaczy wynik bez wiersza, chyba że
	// żądanie poda `windowId` (`oknoWynikuArsenalu`).
	return zrodloObrazu{
		sciezka: sciezka,
		rozmiar: int(opis.Size()),
		okno:    "",
		nazwa:   &nazwa,
	}, nil
}

// wolajImageMagick przeprowadza jedno uruchomienie i oddaje to, co program
// wypisał na wyjście — bajty obrazu przy przetwarzaniu, tekst opisu przy
// badaniu.
//
// Tryb wybieramy, nie zakładamy. ImageMagick siódmy stoi jednym plikiem
// `magick`, szósty — osobnymi `convert` i `identify`. Pytamy `zewnetrzne.Stoi`
// po kolei i pierwszy obecny wygrywa; gdy nie ma żadnego, odmowa nazywa tryb
// siódmy, bo to on jest wskazówką do instalacji. `podpolecenie` jest nazwą
// trybu wersji siódmej (`identify`) i wchodzi na początek argumentów tylko
// wtedy, gdy wołamy pliku `magick`; wersja szósta ma na to osobne binarium
// i podpolecenia nie rozumie.
func (a *adapterNarzedziObrazu) wolajImageMagick(ctx context.Context,
	podpolecenie string, argumenty []string, zapasowy string) ([]byte, error) {

	if a.uruchamiacz == nil {
		return nil, bladZapleczaNiedostepnego("rdzeń nie ma uruchamiacza procesów — " +
			"narzędzia obrazu nie mają czym wystartować; " +
			"naprawa: podpiąć warstwę kanału (injection) przy składaniu rdzenia")
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

// wybierzProgram rozstrzyga, którym plikiem wykonywalnym wołamy ImageMagicka.
// `zapasowy` jest odpowiednikiem z wersji szóstej dla danej czynności
// (`convert` dla przetwarzania, `identify` dla badania).
func (a *adapterNarzedziObrazu) wybierzProgram(zapasowy string) string {
	if zewnetrzne.Stoi(narzedzieImageMagick("magick")) {
		return "magick"
	}
	if zapasowy != "" && zewnetrzne.Stoi(narzedzieImageMagick(zapasowy)) {
		return zapasowy
	}
	// Nic nie stoi — oddajemy nazwę trybu siódmego, żeby odmowa niosła
	// właściwą wskazówkę instalacyjną.
	return "magick"
}

// zasiegNarzedzi składa trójkę okno–zasady–obszar dla zasięgu platformy, tą
// samą drogą i z tego samego powodu, co silnik syntezy mowy: żądanie `image.*`
// niesie sam obraz, a nie okno rozmowy, więc adresem jest najszerszy poziom
// zasięgu, a nie puste struktury znaczące wyłączoną izolację.
//
// Okno dostaje `ExecutionEnvCore` jawnie: przetwarzanie idzie na dysku rdzenia,
// bo to rdzeń odkłada potem bajty wyniku do własnego magazynu.
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

// odlozZasob utrwala bajty wyniku i zakłada jego wiersz — tą samą drogą, którą
// idzie `design.asset.upload`.
//
// Kolejność jest zamierzona: najpierw bajty, potem wiersz. Wiersz powstaje
// dopiero po utrwaleniu treści, inaczej Assets Panel pokazywałby kafelek, za
// którym nie ma nic.
//
// Wymiary mierzymy z tego, co legło w magazynie — `rozpoznajObrazZasobu` czyta
// sam nagłówek pliku. Wyliczenie ich z argumentów żądania dałoby liczbę
// zgadniętą, nieodróżnialną w panelu od zmierzonej: przy zachowanych
// proporcjach wysokość wychodzi inna, niż podał model.
//
// Okno i rodzaj rozstrzyga wspólny kod arsenału (`adapter_narzedzia_wynik.go`),
// a nie ten plik. `oknoZadane` jest polem `windowId` żądania — nieobowiązkowym,
// więc jego brak nie wstrzymuje czynności: bajty idą do magazynu tak czy owak,
// tylko zasób nie pojawia się wtedy w wykazie okna.
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
	// Format zmierzony bije zadeklarowany, ale AVIF nie ma dekodera nagłówka
	// w drzewie — wtedy zostaje format, który sami kazaliśmy wytworzyć.
	//
	// Wymiary zostają wtedy puste, bo pusto znaczy „nie wiem", a zgadnięta
	// liczba wygląda w panelu identycznie jak zmierzona; tak samo rozstrzyga to
	// `design.asset.upload`. Model, który tych wymiarów potrzebuje, ma na to
	// `image.inspect` — czynność pytająca o nie program, a nie bibliotekę
	// standardową.
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

// nazwaZasobuWyniku składa nazwę czytelną dla Operatora: „portret.png ·
// resize". Nazwa źródła nieznana zostawia samą czynność — pusty tekst byłby
// kafelkiem bez podpisu.
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

// bladArsenaluObrazu przekłada odmowę pakietu `zewnetrzne` na kod kontraktu.
//
// Brak binarium dostaje inny kod niż niepowodzenie programu, bo brak usuwa się
// jedną instalacją, a wywrócony program jest usterką przetwarzania. Zlanie obu
// w jeden kod kazałoby szukać usterki tam, gdzie jej nie ma.
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

// bladWskazaniaObrazu nazywa brak albo niepoprawność wskazania w żądaniu —
// usterka wołającego, nie rdzenia.
func bladWskazaniaObrazu(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"narzędzia obrazu: "+powod))
}

// bladPrzetwarzaniaObrazu znakuje przetwarzanie, które ruszyło i się nie udało.
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
