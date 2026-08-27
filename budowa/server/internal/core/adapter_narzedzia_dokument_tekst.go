// Odpowiedzialność pliku: czynność `document.text.extract` — narzędzie, którym
// model czyta plik dostany od Operatora: PDF, skan, zdjęcie kartki. Bez niej
// model nie ma żadnej drogi do treści pliku leżącego na dysku Operatora:
// zamiana formatów jest wygodą, a odczyt warunkiem rozmowy o dokumencie.
//
// ── `usedOcr` ma mówić prawdę ──────────────────────────────────────────────
// Pole `usedOcr` nie jest ciekawostką techniczną — mówi modelowi, czy wolno mu
// zacytować tę treść jako fakt:
//
//   - `usedOcr: false` — tekst pochodzi z warstwy tekstowej dokumentu. To są te
//     same znaki, które wpisał autor; „0" nie zamieni się w „O", a kwota nie
//     zgubi przecinka. Model może cytować dosłownie.
//   - `usedOcr: true` — tekst odczytano z pikseli. Rozpoznanie pisma myli znaki
//     podobne, gubi kolumny i wymyśla spacje. Model ma to traktować jak relację
//     świadka, nie jak dokument: cytując, powinien zaznaczyć, skąd treść
//     pochodzi.
//
// Dlatego kolejność jest jedna i nieodwracalna: najpierw warstwa tekstowa,
// dopiero po jej braku (albo na wyraźne `forceOcr`) rasteryzacja i rozpoznanie.
// Odwrotna kolejność byłaby szybsza do napisania i kłamliwa w skutkach:
// dokument z doskonałą warstwą tekstową wracałby jako odczyt z pikseli, a model
// bez potrzeby przestałby ufać własnemu materiałowi.
//
// Wartość `usedOcr` nie wynika z długości tekstu i nie wolno jej z niej
// wyprowadzać — wynika wyłącznie z tego, która droga dała treść.
//
// ── Brak obu dróg naraz jest odmową ────────────────────────────────────────
// Skan bez warstwy tekstowej, na maszynie bez Tesseracta, to dokument, którego
// rdzeń nie umie przeczytać. Wraca odmowa nazywająca brak — nie pusty tekst
// z `usedOcr: false`, bo pusty tekst znaczy „dokument jest pusty", a to jest
// zdanie o dokumencie, nie o rdzeniu.
//
// ── Format spoza słownika rdzenia ──────────────────────────────────────────
// Słownik `formatyDokumentu` zna dziewięć formatów i jest wykazem tego, co
// rdzeń umie ZAMIENIAĆ. Odczyt jest czymś innym niż zamiana: model dostaje od
// Operatora arkusz, prezentację, wiadomość poczty albo plik biurowy spoza tej
// dziewiątki i pytanie brzmi „co tam jest napisane", a nie „na co to zamienić".
// Materiał, którego słownik nie zna, idzie więc do Apache Tiki — biblioteki,
// której cała robota polega na rozpoznaniu rodzaju pliku i wydobyciu z niego
// tekstu. Tika nie wypiera żadnej z istniejących dróg: format, który słownik
// zna, jedzie jak jechał, bo Pandoc i poppler znają jego strukturę lepiej.
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

// jezykRozpoznaniaDomyslny jest polski, bo tak stanowi kontrakt tej komendy
// („brak bierze polski") i bo taki materiał dostaje ten produkt. Wskazanie
// języka trafia do Tesseracta bez zmian po sprowadzeniu skrótów do jego
// nazewnictwa trójliterowego.
const jezykRozpoznaniaDomyslny = "pol"

// ── Apache Tika ─────────────────────────────────────────────────────────────
//
// Tika nie jest plikiem wykonywalnym: jest zbiorem archiwów Javy, które ktoś
// musi uruchomić maszyną wirtualną. Rdzeń rozdziela więc dwie rzeczy, tak samo
// jak przy silniku mowy, gdzie osobno stoi binarium `piper`, a osobno plik
// głosu:
//
//   - PROGRAM to `java` — i to jego dotyczy deklaracja narzędzia, sonda
//     obecności oraz wykaz zależności. Ścieżka wyszukiwania systemu odpowiada
//     na pytanie o niego wprost;
//   - ARCHIWUM to `tika-app-*.jar` wraz z bibliotekami wydania. Nie jest
//     programem, więc `zewnetrzne.Stoi` nie ma o co go zapytać — jego brak jest
//     osobną odmową, z osobną naprawą („dołożyć wydanie Tiki"), bo naprawa
//     „zainstalować Javę" niczego by tu nie załatwiła.
const (
	// zmiennaTiki jest wskazaniem Operatora, gdzie leży wydanie Tiki —
	// pierwszeństwo przed miejscem typowym, tą samą zasadą co `DANACO_PIPER`
	// przy syntezie mowy. Arsenał instaluje się poza produktem i bywa na każdej
	// maszynie gdzie indziej.
	zmiennaTiki = "DANACO_TIKA"
	// katalogTikiTypowy jest miejscem sprawdzanym, gdy zmiennej nie ma.
	katalogTikiTypowy = "/opt/tika"
	// klasaTiki to punkt wejścia wiersza poleceń Tiki. Nazwa klasy, nie nazwa
	// pliku — archiwum wskazujemy ścieżką, a klasę nazwą, bo archiwum jest
	// wersjonowane, a klasa nie.
	klasaTiki = "org.apache.tika.cli.TikaCLI"
)

// narzedzieTiki opisuje maszynę wirtualną, którą Tika się uruchamia.
//
// Nazwa czytelna mówi o Tice, nie o Javie, bo Operator, któremu odmówiono
// odczytu pliku, ma przeczytać, czego brakuje do ODCZYTU. Pakiet wymienia obie
// rzeczy, bo obie są warunkiem.
var narzedzieTiki = zewnetrzne.Narzedzie{
	Nazwa:   "Apache Tika (uruchamiana środowiskiem Javy)",
	Program: "java",
	Pakiet: "środowisko uruchomieniowe Javy (default-jre) wraz z wydaniem Apache Tika " +
		"w " + katalogTikiTypowy + " albo w katalogu wskazanym zmienną " + zmiennaTiki,
}

// katalogTiki oddaje katalog, w którym rdzeń szuka wydania Tiki.
func katalogTiki() string {
	if wskazany := strings.TrimSpace(os.Getenv(zmiennaTiki)); wskazany != "" {
		return wskazany
	}
	return katalogTikiTypowy
}

// sciezkaKlasTiki składa ścieżkę klas dla maszyny wirtualnej: archiwum wiersza
// poleceń oraz każdy katalog bibliotek wydania.
//
// Wersja archiwum NIE jest wpisana — nazwa pliku niesie numer wydania, a numer
// wpisany w kod rdzenia rozjechałby się z pierwszą aktualizacją Tiki i objawił
// odmową u Operatora. Wzorzec `tika-app-*.jar` jest nazwą, którą to wydanie
// nosi od lat.
//
// Katalogi bibliotek dokłada się dlatego, że wydania Tiki bywają dwojakie:
// archiwum samowystarczalne (niesie zależności w sobie) albo archiwum cienkie
// obok katalogu `lib`. Rdzeń nie zgaduje, które ma przed sobą — dokłada każdy
// `lib`, jaki w katalogu wydania stoi, a gdy nie stoi żaden, ścieżka klas
// zostaje samym archiwum i wydanie samowystarczalne rusza tak samo. Gwiazdka na
// końcu katalogu jest wieloznacznikiem MASZYNY WIRTUALNEJ, nie powłoki —
// rozwija ją Java i znaczy „wszystkie archiwa w tym katalogu".
func sciezkaKlasTiki() (string, bool) {
	katalog := katalogTiki()
	archiwa, err := filepath.Glob(filepath.Join(katalog, "tika-app-*.jar"))
	if err != nil || len(archiwa) == 0 {
		return "", false
	}
	// Kolejność wydań jest kolejnością nazw, a nazwa niesie numer wersji —
	// przy dwóch wydaniach obok siebie bierzemy późniejsze, zamiast pozwalać
	// systemowi plików rozstrzygnąć to za rdzeń.
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

// odmowaBrakuTiki nazywa brak wydania Tiki. Osobna od odmowy braku programu,
// bo naprawa jest inna: Java może stać, a archiwum i tak nie ma.
func odmowaBrakuTiki(format string) error {
	return odmowaDokumentu(shared.ErrorCodeChannelUnavailable,
		"rdzeń nie ma czym odczytać materiału "+opisFormatuMaterialu(format)+
			": wydania Apache Tiki nie ma w "+katalogTiki()+
			" (szukane archiwum `tika-app-*.jar`); naprawa: rozpakować wydanie Tiki "+
			"do tego katalogu albo wskazać jego położenie zmienną "+zmiennaTiki)
}

// opisFormatuMaterialu nazywa format materiału albo jego brak. Odmowa mówiąca
// „materiału ” nie mówi nic.
func opisFormatuMaterialu(format string) string {
	if strings.TrimSpace(format) == "" {
		return "o nierozpoznanym formacie"
	}
	return "w formacie " + format
}

// WyciagnijTekst obsługuje `document.text.extract`.
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
		// Format nierozpoznany nie znaczy jeszcze „nie do odczytania": słownik
		// rdzenia zna dziewięć formatów, a Tika rozpoznaje rodzaj pliku sama,
		// z jego zawartości. Odmowa zostaje na wypadek, gdy Tiki nie ma.
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

	switch {
	case obrazyDokumentu[zrodlo.format]:
		// Obraz ma wyłącznie piksele. `forceOcr` niczego tu nie zmienia — nie ma
		// warstwy tekstowej, którą dałoby się pominąć, więc `usedOcr` jest
		// prawdziwe zawsze i bez wyjątku.
		tekst, err := a.rozpoznajPismo(ctx, zrodlo.sciezka, jezyk)
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
		return a.tekstZPdf(ctx, katalogPracy, zrodlo.sciezka, jezyk, odStrony, doStrony, wymuszone)

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

// tekstZDokumentu czyta treść formatu strukturalnego (docx, odt, html,
// markdown, rtf, epub, csv, txt) Pandokiem.
//
// `usedOcr` jest tu fałszem zawsze i zasłużenie: żaden piksel nie brał udziału,
// znaki pochodzą wprost z pliku. Pole `pages` zostaje puste, bo formaty
// strumieniowe stron nie mają — wpisana jedynka byłaby liczbą zmyśloną.
func (a *adapterNarzedziDokumentu) tekstZDokumentu(ctx context.Context,
	zrodlo zrodloDokumentu) (shared.DocumentTextExtractResponse, error) {

	opis, jest := formatyDokumentu[zrodlo.format]
	if !jest || !opis.czytaPandoc {
		return shared.DocumentTextExtractResponse{}, bladZadaniaDokumentu(
			"rdzeń nie umie odczytać treści formatu " + zrodlo.format +
				"; formaty znane: " + wykazFormatowDokumentu())
	}
	// Pliku tekstowego nie ma z czego wydobywać — jego treść JEST tekstem, więc
	// czytamy go wprost. Droga przez Pandoc kończyła się tu odmową: nazwa
	// pandokowa formatu `txt` brzmi `plain`, a Pandoc zna `plain` wyłącznie jako
	// format zapisu i nie ma czytnika o tej nazwie. Odczyt własny jest przy tym
	// jedyną drogą, która nie potrzebuje programu spoza instalki.
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

// tekstTika czyta materiał, którego słownik rdzenia nie zna, wierszem poleceń
// Apache Tiki.
//
// `usedOcr` jest tu fałszem i zasłużenie — Tika czyta ZNAKI zapisane w pliku,
// nie piksele. Pole `pages` zostaje puste: wyjście `--text` jest strumieniem
// treści bez znaków podziału stron, a jedynka wpisana z góry byłaby liczbą
// zmyśloną.
//
// Diagnostyka Tiki idzie osobnym strumieniem (wiersze `INFO` o włączonych
// rozszerzeniach) i nie miesza się z treścią — `zewnetrzne.Wolaj` trzyma oba
// strumienie osobno właśnie po to.
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
		// Pustka po programie, który skończył się powodzeniem, nie jest zdaniem
		// o dokumencie: Tika oddaje ją tak samo wtedy, gdy plik jest pusty, jak
		// wtedy, gdy nie ma czytnika dla jego rodzaju. Rdzeń nie ma czym tych
		// dwóch rzeczy rozróżnić, więc nie orzeka o żadnej.
		return shared.DocumentTextExtractResponse{}, odmowaDokumentu(shared.ErrorCodeInternalError,
			"Apache Tika nie odczytała z tego materiału ani jednego znaku — plik może "+
				"być pusty albo być rodzajem, dla którego Tika nie ma czytnika; "+
				"naprawa: sprawdzić plik albo wskazać jego format polem fromFormat "+
				"komendy document.convert")
	}
	return shared.DocumentTextExtractResponse{Text: tekst, UsedOcr: false}, nil
}

// tekstZPdf prowadzi rozstrzygnięcie opisane w nagłówku pliku: warstwa
// tekstowa, a gdy jej nie ma albo Operator wymusił — rasteryzacja i
// rozpoznanie pisma.
func (a *adapterNarzedziDokumentu) tekstZPdf(ctx context.Context, katalogPracy, plik, jezyk string,
	odStrony, doStrony *int, wymuszone bool) (shared.DocumentTextExtractResponse, error) {

	warstwa, stron, bladWarstwy := a.warstwaTekstowaPdf(ctx, plik, odStrony, doStrony)
	if bladWarstwy != nil && !wymuszone {
		// Bez wymuszenia niepowodzenie odczytu warstwy jest odmową wprost:
		// zejście po cichu na rozpoznanie pisma oddałoby tekst gorszej jakości
		// bez powiedzenia, dlaczego.
		return shared.DocumentTextExtractResponse{}, bladWarstwy
	}
	if !wymuszone && strings.TrimSpace(warstwa) != "" {
		return shared.DocumentTextExtractResponse{
			Text: warstwa, Pages: liczbaStronDokumentu(stron), UsedOcr: false,
		}, nil
	}

	tekst, przetworzone, err := a.rozpoznajPismoWPdf(ctx, katalogPracy, plik, jezyk, odStrony, doStrony)
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

// warstwaTekstowaPdf czyta warstwę tekstową dokumentu i przy okazji liczy
// strony. Liczba bierze się ze znaków wysuwu strony, którymi `pdftotext`
// rozdziela strony w wyjściu — to policzenie, nie szacunek, i nie wymaga
// kolejnego binarium na maszynie.
func (a *adapterNarzedziDokumentu) warstwaTekstowaPdf(ctx context.Context, plik string,
	odStrony, doStrony *int) (string, int, error) {

	// Kodowanie wymuszone na UTF-8: bez tego `pdftotext` bierze zestaw znaków
	// z ustawień lokalnych maszyny i polskie znaki wracają do modelu
	// przekręcone. Wysuwu strony nie wyłączamy — jest miarą liczby stron niżej.
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
	// Wysuwy strony znikają z treści dopiero po policzeniu — model
	// dostaje tekst, a nie znaki sterujące terminala.
	return strings.ReplaceAll(tekst, "\f", "\n"), strony, nil
}

// rozpoznajPismoWPdf rozkłada strony na obrazy i puszcza każdą przez
// rozpoznanie pisma.
//
// 300 DPI jest wyborem, nie przypadkiem: to rozdzielczość, przy której
// Tesseract czyta pismo drukowane pewnie, a strona A4 mieści się w kilku
// megabajtach. Niżej gubi znaki diakrytyczne — a w polskim materiale różnica
// między „gęślą" a „geslą" jest różnicą między odczytem a zmyśleniem.
func (a *adapterNarzedziDokumentu) rozpoznajPismoWPdf(ctx context.Context, katalogPracy, plik,
	jezyk string, odStrony, doStrony *int) (string, int, error) {

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
	// Kolejność stron jest treścią. `pdftoppm` numeruje pliki z wiodącymi
	// zerami, więc porządek leksykalny pokrywa się z porządkiem stron;
	// bez sortowania kolejność zależałaby od systemu plików i akapity
	// wracałyby przestawione.
	sort.Strings(obrazy)

	czesci := make([]string, 0, len(obrazy))
	for _, obraz := range obrazy {
		tekst, err := a.rozpoznajPismo(ctx, obraz, jezyk)
		if err != nil {
			return "", 0, err
		}
		czesci = append(czesci, tekst)
	}
	return strings.Join(czesci, "\n"), len(obrazy), nil
}

// rozpoznajPismo puszcza jeden obraz przez Tesseracta i oddaje odczytany tekst.
// Wynik idzie na wyjście standardowe (`stdout`), więc nic nie ląduje na dysku
// poza materiałem, który i tak zniknie z katalogiem roboczym czynności.
func (a *adapterNarzedziDokumentu) rozpoznajPismo(ctx context.Context, obraz, jezyk string) (string, error) {
	wyjscie, err := a.wolaj(ctx, narzedzieTesseract,
		[]string{obraz, "stdout", "-l", jezyk}, granicaRozpoznaniaDokumentu)
	if err != nil {
		return "", err
	}
	return string(wyjscie), nil
}

// jezykRozpoznaniaDokumentu sprowadza wskazanie wołającego do nazwy, którą zna
// Tesseract. Nazwa nierozpoznana nie jest podmieniana na domyślną: Tesseract ma
// setkę języków, rdzeń nie ma prawa udawać, że zna ich wykaz, a ciche zejście
// na polski przy wskazaniu „deu" dałoby odczyt niemieckiego skanu polskim
// słownikiem — wynik wygląda jak tekst i jest zmyśleniem.
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

// zakresStronDokumentu sprawdza wskazanie stron. Strona zerowa i ujemna nie
// istnieje, a zakres odwrócony jest pomyłką wołającego, nie zakresem pustym:
// odmowa mówi mu o niej wprost, zamiast oddać pusty tekst.
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

// zakresDlaPopplera przekłada zakres kontraktu na przełączniki `-f`/`-l`, które
// rozumieją oba narzędzia poppler-utils tak samo. Brak wskazania nie dokłada
// przełącznika — narzędzie bierze wtedy cały dokument.
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
// nic nie dało. Zero stron nie jest liczbą stron — jest jej brakiem, a wpisane
// w odpowiedź wyglądałoby jak dokument bez stron.
func liczbaStronDokumentu(stron int) *int {
	if stron <= 0 {
		return nil
	}
	return &stron
}
