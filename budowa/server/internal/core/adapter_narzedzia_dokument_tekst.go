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
package core

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"danacoconsole/shared"
)

// jezykRozpoznaniaDomyslny jest polski, bo tak stanowi kontrakt tej komendy
// („brak bierze polski") i bo taki materiał dostaje ten produkt. Wskazanie
// języka trafia do Tesseracta bez zmian po sprowadzeniu skrótów do jego
// nazewnictwa trójliterowego.
const jezykRozpoznaniaDomyslny = "pol"

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
		return shared.DocumentTextExtractResponse{}, bladZadaniaDokumentu(
			"formatu materiału nie da się rozpoznać po pliku — naprawa: wskazać plik " +
				"z rozszerzeniem albo zasób niosący format")
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
