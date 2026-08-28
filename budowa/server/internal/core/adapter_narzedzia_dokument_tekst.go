// Plik obsługuje czynność document.text.extract, czytającą PDF, obraz oraz
// dokument strukturalny lub biurowy narzędziami Pandoc, poppler, Tesseract
// i Apache Tika dobieranymi według formatu materiału.
package core

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"danacoconsole/server/internal/zewnetrzne"
	"danacoconsole/shared"
)

// Domyślnym językiem rozpoznania pisma jest polski, ponieważ kontrakt czynności
// document.text.extract nakazuje ten wybór przy braku wskazania, a nazwa trafia
// do Tesseracta w jego zapisie trójliterowym.
const jezykRozpoznaniaDomyslny = "pol"

// Apache Tika wymaga osobnego rozpoznania programu i archiwum: uruchamia ją
// środowisko Javy, a samą bibliotekę dostarcza archiwum wydania Tiki, którego
// dostępność sprawdza się inaczej niż obecność programu.
const (
	// Zmienna DANACO_TIKA wskazuje katalog wydania Apache Tiki i ma
	// pierwszeństwo przed katalogiem typowym, ponieważ pakiet ten instaluje
	// się poza produktem i jego położenie bywa różne na każdej maszynie.
	zmiennaTiki = "DANACO_TIKA"
	// katalogTikiTypowy jest katalogiem, w którym rdzeń szuka wydania Apache
	// Tiki, gdy zmienna środowiskowa DANACO_TIKA nie wskazuje innego
	// położenia tego pakietu na maszynie Operatora.
	katalogTikiTypowy = "/opt/tika"
	// klasaTiki nazywa klasę wejściową wiersza poleceń Apache Tiki; archiwum
	// wskazuje się ścieżką, bo jest wersjonowane, a klasa nazwą, bo jej nazwa
	// nie zmienia się między wydaniami pakietu.
	klasaTiki = "org.apache.tika.cli.TikaCLI"
)

// narzedzieTiki opisuje maszynę wirtualną Javy uruchamiającą Apache Tikę;
// nazwa czytelna wskazuje na Tikę, ponieważ to jej brak dotyczy odmowy
// odczytu pliku zgłaszanej Operatorowi.
var narzedzieTiki = zewnetrzne.Narzedzie{
	Nazwa:   "Apache Tika (uruchamiana środowiskiem Javy)",
	Program: "java",
	Pakiet: "środowisko uruchomieniowe Javy (default-jre) wraz z wydaniem Apache Tika " +
		"w " + katalogTikiTypowy + " albo w katalogu wskazanym zmienną " + zmiennaTiki,
}

// katalogTiki oddaje katalog, w którym rdzeń szuka wydania Apache Tiki:
// wskazanie zmienną DANACO_TIKA, gdy jest ustawione, inaczej katalog typowy
// tego pakietu na maszynie.
func katalogTiki() string {
	if wskazany := strings.TrimSpace(os.Getenv(zmiennaTiki)); wskazany != "" {
		return wskazany
	}
	return katalogTikiTypowy
}

// sciezkaKlasTiki składa ścieżkę klas maszyny wirtualnej z najnowszego
// archiwum wydania Apache Tiki i z każdego katalogu lib znalezionego
// w katalogu wydania, ponieważ wydania bywają samowystarczalne albo cienkie.
func sciezkaKlasTiki() (string, bool) {
	katalog := katalogTiki()
	archiwa, err := filepath.Glob(filepath.Join(katalog, "tika-app-*.jar"))
	if err != nil || len(archiwa) == 0 {
		return "", false
	}
	// Nazwa archiwum niesie numer wersji, więc przy dwóch wydaniach bierze się późniejsze.
	sort.Strings(archiwa)
	czlony := []string{archiwa[len(archiwa)-1]}

	katalogiBibliotek, _ := filepath.Glob(filepath.Join(katalog, "lib"))
	zagniezdzone, _ := filepath.Glob(filepath.Join(katalog, "*", "lib"))
	katalogiBibliotek = append(katalogiBibliotek, zagniezdzone...)
	sort.Strings(katalogiBibliotek)
	for _, biblioteki := range katalogiBibliotek {
		if opis, err := os.Stat(biblioteki); err == nil && opis.IsDir() {
			czlony = append(czlony, filepath.Join(biblioteki, "*"))
		}
	}
	return strings.Join(czlony, string(os.PathListSeparator)), true
}

// odmowaBrakuTiki nazywa brak wydania Apache Tiki jako osobną odmowę od braku
// programu Javy, ponieważ obu tych braków nie usuwa ta sama naprawa na
// maszynie Operatora.
func odmowaBrakuTiki(format string) error {
	return odmowaDokumentu(shared.ErrorCodeChannelUnavailable,
		"rdzeń nie ma czym odczytać materiału "+opisFormatuMaterialu(format)+
			": wydania Apache Tiki nie ma w "+katalogTiki()+
			" (szukane archiwum `tika-app-*.jar`); naprawa: rozpakować wydanie Tiki "+
			"do tego katalogu albo wskazać jego położenie zmienną "+zmiennaTiki)
}

// opisFormatuMaterialu nazywa format materiału w odmowie albo jego brak,
// ponieważ zdanie kończące się słowem materiału bez dopełnienia nie niesie
// informacji dla Operatora.
func opisFormatuMaterialu(format string) string {
	if strings.TrimSpace(format) == "" {
		return "o nierozpoznanym formacie"
	}
	return "w formacie " + format
}

// WyciagnijTekst obsługuje czynność document.text.extract: ustala zakres
// stron, katalog roboczy oraz źródło materiału, po czym dobiera drogę
// odczytu do rozpoznanego formatu pliku.
func (a *adapterNarzedziDokumentu) WyciagnijTekst(ctx context.Context,
	z shared.DocumentTextExtractRequest) (shared.DocumentTextExtractResponse, error) {

	odStrony, doStrony, err := zakresStronDokumentu(z.PageFrom, z.PageTo)
	if err != nil {
		return shared.DocumentTextExtractResponse{}, err
	}

	katalogPracy, posprzataj, err := katalogPracyDokumentu()
	if err != nil {
		return shared.DocumentTextExtractResponse{}, err
	}
	defer posprzataj()

	zrodlo, err := a.ustalZrodlo(ctx, katalogPracy, z.AssetId, z.SourcePath, nil, nil,
		"document.text.extract")
	if err != nil {
		return shared.DocumentTextExtractResponse{}, err
	}
	if zrodlo.format == "" {
		// Format nierozpoznany po pliku nie wyklucza odczytu — Tika rozpoznaje rodzaj z zawartości.
		if zewnetrzne.Stoi(narzedzieTiki) {
			return a.tekstTika(ctx, zrodlo)
		}
		return shared.DocumentTextExtractResponse{}, bladZadaniaDokumentu(
			"formatu materiału nie da się rozpoznać po pliku, a rdzeń nie ma czym " +
				"rozpoznać go z zawartości — naprawa: wskazać plik z rozszerzeniem, " +
				"zasób niosący format albo dołożyć środowisko Javy wraz z wydaniem " +
				"Apache Tiki")
	}

	jezyk := jezykRozpoznaniaDokumentu(z.Language)
	wymuszone := z.ForceOcr != nil && *z.ForceOcr
	obrobkaWstepna := z.Preprocess != nil && *z.Preprocess

	switch {
	case obrazyDokumentu[zrodlo.format]:
		// Obraz ma wyłącznie piksele, więc usedOcr jest tu prawdziwe zawsze i bez wyjątku.
		tekst, err := a.rozpoznajPismo(ctx, zrodlo.sciezka, jezyk, obrobkaWstepna)
		if err != nil {
			return shared.DocumentTextExtractResponse{}, err
		}
		if strings.TrimSpace(tekst) == "" {
			return shared.DocumentTextExtractResponse{}, odmowaDokumentu(
				shared.ErrorCodeInternalError,
				"rozpoznanie pisma na obrazie nie odczytało ani jednego znaku "+
					"(język "+jezyk+") — obraz może nie zawierać tekstu albo być zbyt "+
					"niskiej rozdzielczości; naprawa: podać materiał lepszej jakości "+
					"albo wskazać właściwy język polem language")
		}
		strony := 1
		return shared.DocumentTextExtractResponse{Text: tekst, Pages: &strony, UsedOcr: true}, nil

	case zrodlo.format == "pdf":
		return a.tekstZPdf(ctx, katalogPracy, zrodlo.sciezka, jezyk, odStrony, doStrony,
			wymuszone, obrobkaWstepna)

	default:
		if wymuszone {
			return shared.DocumentTextExtractResponse{}, bladZadaniaDokumentu(
				"wymuszone rozpoznanie pisma (forceOcr) dotyczy PDF-u i obrazów — " +
					"materiał w formacie " + zrodlo.format + " nie ma pikseli do rozpoznania, " +
					"a jego treść jest tekstem wprost")
		}
		return a.tekstZDokumentu(ctx, zrodlo)
	}
}

// tekstZDokumentu czyta treść formatu strukturalnego Pandokiem; usedOcr
// pozostaje fałszem, bo żaden piksel nie bierze udziału, a pole pages zostaje
// puste, ponieważ formaty strumieniowe stron nie mają.
func (a *adapterNarzedziDokumentu) tekstZDokumentu(ctx context.Context,
	zrodlo zrodloDokumentu) (shared.DocumentTextExtractResponse, error) {

	opis, jest := formatyDokumentu[zrodlo.format]
	if !jest || !opis.czytaPandoc {
		return shared.DocumentTextExtractResponse{}, bladZadaniaDokumentu(
			"rdzeń nie umie odczytać treści formatu " + zrodlo.format +
				"; formaty znane: " + wykazFormatowDokumentu())
	}
	// Pandoc zna nazwę formatu txt wyłącznie jako zapis plain, bez czytnika o tej nazwie.
	if zrodlo.format == "txt" {
		bajty, err := os.ReadFile(zrodlo.sciezka)
		if err != nil {
			return shared.DocumentTextExtractResponse{}, odmowaDokumentu(
				shared.ErrorCodeInternalError,
				"nie można odczytać pliku tekstowego "+zrodlo.sciezka+": "+err.Error())
		}
		return shared.DocumentTextExtractResponse{Text: string(bajty), UsedOcr: false}, nil
	}
	wyjscie, err := a.wolaj(ctx, narzedziePandoc, []string{
		"--from", opis.pandoc, "--to", "plain", "--wrap", "none",
		"--output", "-", zrodlo.sciezka,
	}, granicaOdczytuDokumentu)
	if err != nil {
		return shared.DocumentTextExtractResponse{}, err
	}
	return shared.DocumentTextExtractResponse{Text: string(wyjscie), UsedOcr: false}, nil
}

// tekstTika czyta materiał nieznany słownikowi rdzenia wierszem poleceń
// Apache Tiki; usedOcr pozostaje fałszem, bo Tika czyta znaki zapisane
// w pliku, a nie piksele, pole pages zostaje puste, a diagnostyka Tiki idzie
// osobnym strumieniem.
func (a *adapterNarzedziDokumentu) tekstTika(ctx context.Context,
	zrodlo zrodloDokumentu) (shared.DocumentTextExtractResponse, error) {

	sciezkaKlas, jest := sciezkaKlasTiki()
	if !jest {
		return shared.DocumentTextExtractResponse{}, odmowaBrakuTiki(zrodlo.format)
	}
	wyjscie, err := a.wolaj(ctx, narzedzieTiki, []string{
		"-cp", sciezkaKlas, klasaTiki, "--text", "--encoding=UTF-8", zrodlo.sciezka,
	}, granicaOdczytuDokumentu)
	if err != nil {
		return shared.DocumentTextExtractResponse{}, err
	}
	tekst := string(wyjscie)
	if strings.TrimSpace(tekst) == "" {
		// Pustka po powodzeniu nie jest zdaniem o dokumencie — Tika oddaje ją tak samo bez czytnika.
		return shared.DocumentTextExtractResponse{}, odmowaDokumentu(shared.ErrorCodeInternalError,
			"Apache Tika nie odczytała z tego materiału ani jednego znaku — plik może "+
				"być pusty albo być rodzajem, dla którego Tika nie ma czytnika; "+
				"naprawa: sprawdzić plik albo wskazać jego format polem fromFormat "+
				"komendy document.convert")
	}
	return shared.DocumentTextExtractResponse{Text: tekst, UsedOcr: false}, nil
}

// tekstZPdf czyta PDF warstwą tekstową w pierwszej kolejności, a rozpoznanie
// pisma stosuje dopiero po jej braku albo na wyraźne żądanie Operatora, bo
// pikselowy odczyt bywa mniej wierny niż zapisane znaki.
func (a *adapterNarzedziDokumentu) tekstZPdf(ctx context.Context, katalogPracy, plik, jezyk string,
	odStrony, doStrony *int, wymuszone, obrobkaWstepna bool) (shared.DocumentTextExtractResponse, error) {

	warstwa, stron, bladWarstwy := a.warstwaTekstowaPdf(ctx, plik, odStrony, doStrony)
	if bladWarstwy != nil && !wymuszone {
		// Bez wymuszenia błąd warstwy jest odmową wprost, nie cichym zejściem na gorszy odczyt.
		return shared.DocumentTextExtractResponse{}, bladWarstwy
	}
	if !wymuszone && strings.TrimSpace(warstwa) != "" {
		return shared.DocumentTextExtractResponse{
			Text: warstwa, Pages: liczbaStronDokumentu(stron), UsedOcr: false,
		}, nil
	}

	tekst, przetworzone, err := a.rozpoznajPismoWPdf(ctx, katalogPracy, plik, jezyk,
		odStrony, doStrony, obrobkaWstepna)
	if err != nil {
		return shared.DocumentTextExtractResponse{}, err
	}
	if strings.TrimSpace(tekst) == "" {
		return shared.DocumentTextExtractResponse{}, odmowaDokumentu(shared.ErrorCodeInternalError,
			"dokument nie ma warstwy tekstowej, a rozpoznanie pisma (język "+jezyk+
				") nie odczytało z jego stron ani jednego znaku — rdzeń nie umie "+
				"przeczytać tego materiału; naprawa: sprawdzić, czy strony nie są puste, "+
				"albo wskazać właściwy język polem language")
	}
	return shared.DocumentTextExtractResponse{
		Text: tekst, Pages: liczbaStronDokumentu(przetworzone), UsedOcr: true,
	}, nil
}

// warstwaTekstowaPdf czyta warstwę tekstową dokumentu i liczy strony ze
// znaków wysuwu strony, którymi pdftotext rozdziela strony w wyjściu, więc
// liczba stron jest policzeniem, nie szacunkiem.
func (a *adapterNarzedziDokumentu) warstwaTekstowaPdf(ctx context.Context, plik string,
	odStrony, doStrony *int) (string, int, error) {

	// Kodowanie wymuszone na UTF-8, inaczej polskie znaki wracają przekręcone ustawieniami lokalnymi.
	argumenty := []string{"-enc", "UTF-8"}
	argumenty = append(argumenty, zakresDlaPopplera(odStrony, doStrony)...)
	argumenty = append(argumenty, plik, "-")

	wyjscie, err := a.wolaj(ctx, narzedziePdfDoTekstu, argumenty, granicaOdczytuDokumentu)
	if err != nil {
		return "", 0, err
	}
	tekst := string(wyjscie)
	strony := strings.Count(tekst, "\f")
	if strony == 0 && strings.TrimSpace(tekst) != "" {
		strony = 1
	}
	// Wysuw strony liczy strony wyżej i znika z treści dopiero teraz, po policzeniu.
	return strings.ReplaceAll(tekst, "\f", "\n"), strony, nil
}

// rozpoznajPismoWPdf rozkłada strony PDF na obrazy o rozdzielczości 300 DPI
// i puszcza każdy przez rozpoznanie pisma; niższa rozdzielczość gubi znaki
// diakrytyczne polskiego materiału.
func (a *adapterNarzedziDokumentu) rozpoznajPismoWPdf(ctx context.Context, katalogPracy, plik,
	jezyk string, odStrony, doStrony *int, obrobkaWstepna bool) (string, int, error) {

	przedrostek := filepath.Join(katalogPracy, "strona")
	argumenty := append([]string{"-r", "300", "-png"}, zakresDlaPopplera(odStrony, doStrony)...)
	argumenty = append(argumenty, plik, przedrostek)
	if _, err := a.wolaj(ctx, narzedziePdfDoObrazu, argumenty, granicaRozpoznaniaDokumentu); err != nil {
		return "", 0, err
	}

	obrazy, err := filepath.Glob(przedrostek + "*.png")
	if err != nil || len(obrazy) == 0 {
		return "", 0, odmowaDokumentu(shared.ErrorCodeInternalError,
			"rasteryzacja nie dała ani jednej strony do rozpoznania — "+
				"dokument może być pusty albo uszkodzony")
	}
	// Nazwy niosą numer strony z wiodącymi zerami, więc porządek leksykalny jest porządkiem stron.
	sort.Strings(obrazy)

	czesci := make([]string, 0, len(obrazy))
	for _, obraz := range obrazy {
		tekst, err := a.rozpoznajPismo(ctx, obraz, jezyk, obrobkaWstepna)
		if err != nil {
			return "", 0, err
		}
		czesci = append(czesci, tekst)
	}
	return strings.Join(czesci, "\n"), len(obrazy), nil
}

// rozpoznajPismo puszcza jeden obraz przez Tesseracta, poprzedzone obróbką
// wstępną unpaperem przy zamówieniu jej polem preprocess, i oddaje odczytany
// tekst wprost ze standardowego wyjścia programu.
func (a *adapterNarzedziDokumentu) rozpoznajPismo(ctx context.Context, obraz, jezyk string,
	obrobkaWstepna bool) (string, error) {

	material, posprzataj, err := a.obrazPoObrobceWstepnej(ctx, obraz, obrobkaWstepna)
	if err != nil {
		return "", err
	}
	defer posprzataj()

	wyjscie, err := a.wolaj(ctx, narzedzieTesseract,
		[]string{material, "stdout", "-l", jezyk}, granicaRozpoznaniaDokumentu)
	if err != nil {
		return "", err
	}
	return string(wyjscie), nil
}

// obrazPoObrobceWstepnej oddaje ścieżkę obrazu bez zmiany, gdy żądanie nie
// zamówiło obróbki wstępnej; przy zamówieniu przepuszcza obraz przez unpapera.
func (a *adapterNarzedziDokumentu) obrazPoObrobceWstepnej(ctx context.Context, obraz string,
	obrobkaWstepna bool) (string, func(), error) {

	pusto := func() {}
	if !obrobkaWstepna {
		return obraz, pusto, nil
	}
	if !zewnetrzne.Stoi(narzedzieCzyszczeniaSkanu) {
		return "", pusto, odmowaDokumentu(shared.ErrorCodeChannelUnavailable,
			"żądanie zamówiło obróbkę wstępną obrazu (pole preprocess), a programu "+
				narzedzieCzyszczeniaSkanu.Nazwa+" ("+narzedzieCzyszczeniaSkanu.Program+
				") nie ma na tej maszynie; naprawa: zainstalować pakiet "+
				narzedzieCzyszczeniaSkanu.Pakiet+
				". Droga, która działa bez niego: wysłać żądanie bez pola preprocess — "+
				"rozpoznanie pobiegnie na materiale bez obróbki")
	}

	katalog, err := os.MkdirTemp("", "danaco-dokument-obrobka-")
	if err != nil {
		return "", pusto, odmowaDokumentu(shared.ErrorCodeInternalError,
			"nie można założyć katalogu roboczego obróbki wstępnej: "+err.Error())
	}
	posprzataj := func() { _ = os.RemoveAll(katalog) }

	wejscie, err := materialWPnm(obraz, katalog)
	if err != nil {
		posprzataj()
		return "", pusto, err
	}
	wyjscie := filepath.Join(katalog, "oczyszczony.ppm")
	if _, err := a.wolaj(ctx, narzedzieCzyszczeniaSkanu,
		argumentyObrobkiWstepnejDokumentu(wejscie, wyjscie), granicaRozpoznaniaDokumentu); err != nil {
		posprzataj()
		return "", pusto, err
	}
	// unpaper potrafi skończyć się powodzeniem i nie zostawić pliku wyniku.
	if opis, err := os.Stat(wyjscie); err != nil || opis.Size() == 0 {
		posprzataj()
		return "", pusto, odmowaDokumentu(shared.ErrorCodeInternalError,
			"unpaper zakończył pracę, ale obrazu po obróbce nie ma pod "+wyjscie+
				" — materiał do rozpoznania nie powstał")
	}
	return wyjscie, posprzataj, nil
}

// argumentyObrobkiWstepnejDokumentu składa wiersz wywołania unpapera dla pola
// preprocess: prostowanie, odszumianie i przycinanie marginesów włączone,
// filtr czerni wyłączony jawnie — tak samo jak w drodze Studia.
func argumentyObrobkiWstepnejDokumentu(wejscie, wyjscie string) []string {
	return []string{"--layout", "single", "--no-blackfilter", wejscie, wyjscie}
}

// jezykRozpoznaniaDokumentu sprowadza wskazanie wołającego do nazwy języka
// znanej Tesseractowi; nazwa nierozpoznana zostaje bez zmian, ponieważ
// podmiana na język domyślny dałaby odczyt zmyślony.
func jezykRozpoznaniaDokumentu(wskazanie *string) string {
	nazwa := strings.ToLower(strings.TrimSpace(wartoscTekstu(wskazanie)))
	switch nazwa {
	case "":
		return jezykRozpoznaniaDomyslny
	case "pl", "pol", "polski", "polish", "pl-pl", "pl_pl":
		return jezykRozpoznaniaDomyslny
	case "en", "eng", "angielski", "english", "en-us", "en_us", "en-gb":
		return "eng"
	}
	return nazwa
}

// zakresStronDokumentu sprawdza wskazanie stron: strona zerowa albo ujemna
// jest błędem, a zakres odwrócony zwraca odmowę wprost, zamiast oddać pusty
// tekst wołającemu bez wyjaśnienia przyczyny.
func zakresStronDokumentu(od, do *int) (*int, *int, error) {
	if od != nil && *od < 1 {
		return nil, nil, bladZadaniaDokumentu("pierwsza strona zakresu jest mniejsza od jedynki")
	}
	if do != nil && *do < 1 {
		return nil, nil, bladZadaniaDokumentu("ostatnia strona zakresu jest mniejsza od jedynki")
	}
	if od != nil && do != nil && *do < *od {
		return nil, nil, bladZadaniaDokumentu("zakres stron jest odwrócony — " +
			"ostatnia strona leży przed pierwszą")
	}
	return od, do, nil
}

// zakresDlaPopplera przekłada zakres kontraktu na przełączniki -f i -l,
// wspólne obu narzędziom pakietu poppler-utils; brak wskazania nie dokłada
// przełącznika i narzędzie bierze cały dokument.
func zakresDlaPopplera(od, do *int) []string {
	argumenty := make([]string, 0, 4)
	if od != nil {
		argumenty = append(argumenty, "-f", strconv.Itoa(*od))
	}
	if do != nil {
		argumenty = append(argumenty, "-l", strconv.Itoa(*do))
	}
	return argumenty
}

// liczbaStronDokumentu oddaje wskaźnik na liczbę stron albo nic, gdy liczenie
// nic nie dało, ponieważ zero stron w odpowiedzi wyglądałoby jak dokument bez
// żadnej strony, a nie brak policzenia.
func liczbaStronDokumentu(stron int) *int {
	if stron <= 0 {
		return nil
	}
	return &stron
}
