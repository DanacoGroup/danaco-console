package core

// Odpowiedzialność pliku: jeden wykaz programów spoza instalki, których rdzeń
// używa, wraz z sondą sprawdzającą ich obecność przy starcie.
//
// ── Po co sonda startowa, skoro każde wywołanie i tak sprawdza ────────────────
// `zewnetrzne.Wolaj` sprawdza obecność programu przed uruchomieniem i odmawia
// zdaniem nazywającym brak (`zewnetrzne.BrakNarzedzia`). To jest właściwe
// zachowanie w chwili czynności, ale za późne jako informacja: Operator dowiaduje
// się o braku dopiero po naciśnięciu przycisku, osobno dla każdej funkcji, i nie
// ma skąd wiedzieć, ile jeszcze takich niespodzianek przed nim.
//
// Sonda startowa odwraca kolejność — rdzeń mówi na wejściu, czym dysponuje
// i czego mu brak, jednym wykazem. To jest ta sama wiedza, ale podana zanim
// zawiedzie czynność.
//
// ── Dlaczego wykaz odwołuje się do deklaracji, a nie powtarza ich ─────────────
// Każde narzędzie jest zadeklarowane tam, gdzie jest używane — Pandoc przy
// dokumentach, ffmpeg przy nagraniach, Tesseract przy rozpoznaniu pisma. Wykaz
// bierze te deklaracje, zamiast wypisywać nazwy po raz drugi: druga lista
// rozjechałaby się z pierwszą przy pierwszej zmianie pakietu i Operator
// czytałby podpowiedź instalacyjną prowadzącą donikąd.

import (
	"sort"

	"danacoconsole/server/internal/zewnetrzne"
)

// ZaleznoscZewnetrzna opisuje jeden program wraz z jego stanem na tej maszynie
// oraz zakresem, który bez niego nie działa.
type ZaleznoscZewnetrzna struct {
	// Narzedzie niesie nazwę czytelną, program i podpowiedź instalacyjną.
	Narzedzie zewnetrzne.Narzedzie
	// Zakres mówi, co przestaje działać przy braku — zdaniem, nie kodem.
	Zakres string
	// Stoi jest wynikiem sondy: czy program jest osiągalny na ścieżce.
	Stoi bool
}

// zaleznosciZewnetrzne oddaje komplet deklaracji w kolejności alfabetycznej
// nazwy czytelnej. Kolejność jest ustalona, żeby dwa kolejne uruchomienia dawały
// ten sam wykaz — dziennik startu ma się różnić wtedy, gdy zmienił się stan
// maszyny, a nie wtedy, gdy inaczej ułożyła się mapa.
func zaleznosciZewnetrzne() []ZaleznoscZewnetrzna {
	wykaz := []ZaleznoscZewnetrzna{
		{Narzedzie: narzedziePandoc,
			Zakres: "zamiana formatów dokumentu w modułach Studio, Translate i Library"},
		{Narzedzie: narzedziePdfDoTekstu,
			Zakres: "odczyt warstwy tekstowej dokumentu PDF"},
		{Narzedzie: narzedziePdfDoObrazu,
			Zakres: "rasteryzacja stron PDF przed rozpoznaniem pisma"},
		{Narzedzie: narzedzieTesseract,
			Zakres: "rozpoznanie pisma ze skanu i zdjęcia kartki"},
		{Narzedzie: narzedzieLibreOffice,
			Zakres: "zamiana formatów biurowych, których nie czyta Pandoc"},
		{Narzedzie: narzedzieFfprobe,
			Zakres: "rozpoznanie zawartości nagrania dźwiękowego i filmowego"},
		{Narzedzie: narzedzieFfmpeg,
			Zakres: "zamiana formatu nagrania"},
		{Narzedzie: narzedzie7z(),
			Zakres: "pakowanie i wydobycie zawartości archiwum"},
		{Narzedzie: narzedziePowiekszenia(),
			Zakres: "powiększanie obrazu w module Design"},
		{Narzedzie: narzedzieImageMagick("magick"),
			Zakres: "zapis obrazu w AVIF oraz w WEBP stratnym, a także pomiar pliku AVIF " +
				"(image.convert, image.inspect) — pozostałe czynności rodziny image.* " +
				"liczy biblioteka wkompilowana i przy braku tego programu pracują dalej"},
		{Narzedzie: narzedzieWycinaniaTla(),
			Zakres: "wycinanie tła obrazu oraz rozkład obrazu na warstwy " +
				"(image.background.remove, image.layers.split) — obie czynności " +
				"stoją na tej samej sieci segmentującej i znikają razem z nią"},
		{Narzedzie: narzedzieSyntezyMowy,
			Zakres: "odsłuch przebiegu debaty syntezą mowy w module Roundtable"},
		{Narzedzie: narzedzieGofmt,
			Zakres: "formatowanie plików Go w module Developer"},
		{Narzedzie: narzedzieGoimports,
			Zakres: "formatowanie plików Go wraz z porządkowaniem importów"},
		{Narzedzie: narzedzieGopls,
			Zakres: "przejście do definicji, wystąpienia symbolu i refaktoryzacje " +
				"semantyczne w module Developer"},
		{Narzedzie: narzedzieGolangciLint,
			Zakres: "analiza statyczna repozytorium Go w module Developer"},
		{Narzedzie: narzedzieStaticcheck,
			Zakres: "pogłębiona analiza statyczna kodu Go"},
		{Narzedzie: narzedzieDelve,
			Zakres: "debugowanie krokowe programów Go w oknie Run & Debug"},
		{Narzedzie: narzedziePrettier,
			Zakres: "formatowanie plików TypeScript, JavaScript, CSS, Markdown i YAML"},
		{Narzedzie: narzedzieEslint,
			Zakres: "analiza statyczna kodu TypeScript i JavaScript"},
		{Narzedzie: narzedzieRuff,
			Zakres: "analiza statyczna plików Pythona w module Developer oraz analiza " +
				"i formatowanie treści skryptu karty python w module Terminal"},
		{Narzedzie: narzedzieStylelint,
			Zakres: "analiza statyczna arkuszy CSS w module Developer — wyłącznie " +
				"w repozytorium niosącym własną konfigurację Stylelinta, bo program " +
				"nie ma wbudowanego zestawu reguł"},
		{Narzedzie: narzedzieTypos,
			Zakres: "wykrywanie literówek w identyfikatorach i treści plików " +
				"repozytorium w module Developer"},
		{Narzedzie: narzedzieAstGrep,
			Zakres: "wyszukanie i zamiana po składni w module Developer — wzorzec " +
				"z metazmienną (`$NAZWA`) idzie tą drogą zamiast po napisie"},
		{Narzedzie: narzedzieSemgrep,
			Zakres: "poszerzenie skanu kodu w module Developer o reguły semantyczne — " +
				"wyłącznie w repozytorium niosącym własny zestaw reguł, bo zestaw " +
				"z rejestru wymagałby sieci"},
		{Narzedzie: narzedzieJscpd,
			Zakres: "wykrywanie powtórzonych fragmentów w plikach TypeScriptu " +
				"i JavaScriptu podczas skanu kodu w module Developer"},
		{Narzedzie: narzedzieDupl,
			Zakres: "wykrywanie powtórzonych fragmentów w plikach Go podczas skanu " +
				"kodu w module Developer"},
		{Narzedzie: narzedzieSerweraTypeScript,
			Zakres: "przejście do definicji, wystąpienia symbolu i refaktoryzacje " +
				"semantyczne plików TypeScriptu w module Developer"},
		{Narzedzie: narzedzieSilnikaKontenerow,
			Zakres: "wykaz kontenerów i obrazów, budowanie obrazu i uruchomienie stosu " +
				"w zakładce Containers"},
		{Narzedzie: narzedzieSSH,
			Zakres: "powłoka zdalna, odczyt pliku na maszynie karty i przekierowania " +
				"portów w module Terminal"},
		{Narzedzie: narzedzieShellCheck,
			Zakres: "analiza statyczna skryptów powłoki bash w module Terminal"},
		{Narzedzie: narzedzieShfmt,
			Zakres: "formatowanie skryptów powłoki bash w module Terminal"},
		{Narzedzie: narzedziePowerShell,
			Zakres: "analiza statyczna i formatowanie skryptów PowerShell w module Terminal"},
		{Narzedzie: narzedzieDocker,
			Zakres: "karta powłoki wewnątrz kontenera w module Terminal"},
		{Narzedzie: narzedzieKubectl,
			Zakres: "karta powłoki wewnątrz poda w module Terminal"},
		{Narzedzie: narzedziePicocom,
			Zakres: "karta konsoli portu szeregowego w module Terminal"},
		{Narzedzie: narzedzieTelnet,
			Zakres: "karta sesji Telnet do urządzenia sieciowego w module Terminal"},
		{Narzedzie: narzedzieChromium(),
			Zakres: "zrzuty stron, drzewo DOM, konsola, rejestr sieciowy, " +
				"emulacja urządzenia i przewijanie w module Browser"},
	}
	for i := range wykaz {
		wykaz[i].Stoi = zewnetrzne.Stoi(wykaz[i].Narzedzie)
	}
	sort.SliceStable(wykaz, func(a, b int) bool {
		return wykaz[a].Narzedzie.Nazwa < wykaz[b].Narzedzie.Nazwa
	})
	return wykaz
}

// ZaleznosciZewnetrzne oddaje wykaz wraz z wynikiem sondy. Wystawione poza
// pakiet, bo czytelnikami są dwie warstwy: punkt wejścia procesu, który wpisuje
// stan do dziennika startu, oraz moduł diagnostyki, który pokazuje go Operatorowi.
func ZaleznosciZewnetrzne() []ZaleznoscZewnetrzna {
	return zaleznosciZewnetrzne()
}

// BrakujaceZaleznosci oddaje same pozycje nieobecne na maszynie. Pusty wynik
// znaczy, że rdzeń ma czym wykonać każdą czynność wymagającą programu spoza
// instalki.
func BrakujaceZaleznosci() []ZaleznoscZewnetrzna {
	var brakujace []ZaleznoscZewnetrzna
	for _, pozycja := range zaleznosciZewnetrzne() {
		if !pozycja.Stoi {
			brakujace = append(brakujace, pozycja)
		}
	}
	return brakujace
}

// zglosZaleznosci wpisuje wynik sondy do dziennika startu.
//
// Wiersz zbiorczy idzie zawsze — także wtedy, gdy nie brakuje niczego, bo
// „wszystkie obecne" jest informacją, nie ciszą. Każdy brak dostaje własny
// wiersz z zakresem, który przestaje działać, i podpowiedzią instalacyjną:
// Operator ma po starcie wiedzieć, czego nie zrobi, zanim spróbuje.
func (r *Rdzen) zglosZaleznosci() {
	wykaz := zaleznosciZewnetrzne()
	brakujace := 0
	for _, pozycja := range wykaz {
		if !pozycja.Stoi {
			brakujace++
		}
	}
	r.zapisz("zależności zewnętrzne: %d z %d obecnych", len(wykaz)-brakujace, len(wykaz))
	for _, pozycja := range wykaz {
		if pozycja.Stoi {
			continue
		}
		zdanie := "brak programu " + pozycja.Narzedzie.Nazwa +
			" (" + pozycja.Narzedzie.Program + ") — nie zadziała: " + pozycja.Zakres
		if pozycja.Narzedzie.Pakiet != "" {
			zdanie += "; naprawa: " + pozycja.Narzedzie.Pakiet
		}
		r.zapisz("%s", zdanie)
	}
}
