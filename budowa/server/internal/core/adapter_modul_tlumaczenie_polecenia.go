// Odpowiedzialność pliku: moduł Translate — składanie poleceń dla modelu
// i wiązanie słownika Operatora z przekładem. Własnego silnika tłumaczeń tu nie
// ma; przekład idzie kanałem modelu, więc cała wiedza modułu o tym, jak ma
// wyglądać dobry przekład, mieści się w treści polecenia.
//
// adapter_modul_tlumaczenie_model.go odpowiada za drogę do modelu (wybór
// kanału, wywołanie, rozgłoszenie), ten plik — za treść polecenia. Nowy termin
// słownika albo nowa zasada jakości rusza ten plik, a nie tamten.
//
// Słownik Operatora wiąże przekład dwiema drogami:
//
//  1. Przed wywołaniem terminy wchodzą do polecenia jako wykaz obowiązkowych
//     odpowiedników i wykaz nazw nietykalnych (`wskazaniaSlownika`). Model,
//     który zna słownik, użyje właściwego słowa w odmienionej formie i we
//     właściwym miejscu zdania — czego podmiana napisu nie potrafi.
//  2. Po wywołaniu wynik przechodzi tę samą mechaniczną podmianę, co
//     w `glossary.apply` (`zastosujTerminySlownika`,
//     adapter_modul_tlumaczenie_slownik.go). Podmiana wyłapuje termin
//     zostawiony w brzmieniu źródłowym, lecz nie zastępuje punktu pierwszego —
//     trafia wyłącznie w formę podstawową.
//
// Zasady jakości także wchodzą do polecenia. Kontrakt zna sześć rodzajów
// niezgodności (`TranslationIssueKind`: number, date, currency, placeholder,
// length, omission) i moduł szuka ich potem w wyniku (`quality.check`).
// `zasadyJakosci` wyprowadza treść zasad wprost ze stałych kontraktu, żeby
// jedna definicja wady przekładu nie rozjechała się na inną w poleceniu i inną
// w kontroli.
package core

import (
	"context"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// wiazaniaSlownika to słownik Operatora zawężony do jednego przekładu:
// odpowiedniki obowiązujące w języku docelowym i nazwy nietykalne.
//
// Rozdział na dwa wykazy jest istotny, nie kosmetyczny. Odpowiednik dotyczy
// jednego języka docelowego — „invoice → faktura" nie obowiązuje w przekładzie
// na niemiecki. Nietykalność dotyczy terminu, nie języka: nazwa własna produktu
// ma zostać nietknięta w każdym przekładzie, więc wykaz nietykalnych zbiera
// terminy wszystkich języków słownika. Zawężenie ich do języka panelu byłoby
// cichym rozluźnieniem zakazu Operatora.
type wiazaniaSlownika struct {
	odpowiedniki []dane.TerminSlownika
	nietykalne   []dane.TerminSlownika
}

// puste mówi, czy słownik nie wniósł do tego przekładu niczego — wtedy
// polecenie nie niesie bloku słownika w ogóle. Pusty nagłówek „SŁOWNIK
// OPERATORA" bez ani jednej pozycji jest szumem, który tylko rozprasza model.
func (w wiazaniaSlownika) puste() bool {
	return len(w.odpowiedniki) == 0 && len(w.nietykalne) == 0
}

// wiazaniaSlownikaDlaJezyka czyta słownik Operatora i zawęża go do przekładu
// na wskazany język wedle zasady opisanej przy `wiazaniaSlownika`.
//
// Błąd odczytu nie jest tu odmową całego przekładu — jest nią w wywołującym
// (`przetlumaczModelem`). Przekład bez słownika, o który Operator prosił, jest
// przekładem cudzym, a nie „prawie dobrym", więc lepiej odmówić wprost niż po
// cichu oddać wynik ignorujący terminologię. Dlatego ta funkcja oddaje błąd,
// zamiast go połykać.
func (a *adapterTlumaczenia) wiazaniaSlownikaDlaJezyka(ctx context.Context,
	jezykDocelowy string) (wiazaniaSlownika, error) {

	terminy, err := a.repozytorium.Terminy(ctx)
	if err != nil {
		return wiazaniaSlownika{}, bladTlumaczenia(err)
	}

	var wynik wiazaniaSlownika
	for _, termin := range terminy {
		if strings.TrimSpace(termin.Zrodlo) == "" {
			continue
		}
		if termin.NieTlumaczyc {
			// Zakaz tłumaczenia obowiązuje bez względu na język panelu.
			wynik.nietykalne = append(wynik.nietykalne, termin)
			continue
		}
		if termin.Jezyk != jezykDocelowy {
			continue
		}
		if termin.Cel == nil || strings.TrimSpace(*termin.Cel) == "" {
			// Termin bez odpowiednika nie niesie żadnego polecenia dla modelu:
			// Operator zaznaczył słowo, ale nie powiedział, czym je zastąpić.
			continue
		}
		wynik.odpowiedniki = append(wynik.odpowiedniki, termin)
	}
	return wynik, nil
}

// wskazaniaSlownika przekłada zawężony słownik na blok polecenia. Uwaga
// Operatora (`Uwaga`, kolumna `uwaga`) idzie do modelu razem z odpowiednikiem —
// bywa, że to ona niesie właściwy powód wyboru słowa („w umowach, nie w
// marketingu"), a bez niej model dostałby samą parę napisów.
func wskazaniaSlownika(w wiazaniaSlownika) string {
	if w.puste() {
		return ""
	}
	var b strings.Builder
	b.WriteString("SŁOWNIK OPERATORA — OBOWIĄZUJE BEZWZGLĘDNIE.\n")
	if len(w.odpowiedniki) > 0 {
		b.WriteString("Poniższe terminy tłumacz WYŁĄCZNIE podanym odpowiednikiem " +
			"(odmieniaj go zgodnie z gramatyką języka docelowego, ale nie zastępuj synonimem):\n")
		for _, termin := range w.odpowiedniki {
			b.WriteString("  • „" + termin.Zrodlo + "” → „" + *termin.Cel + "”")
			if termin.Uwaga != nil && strings.TrimSpace(*termin.Uwaga) != "" {
				b.WriteString(" (uwaga Operatora: " + strings.TrimSpace(*termin.Uwaga) + ")")
			}
			b.WriteString("\n")
		}
	}
	if len(w.nietykalne) > 0 {
		b.WriteString("Poniższych nazw NIE TŁUMACZ — przepisz je w brzmieniu źródłowym, znak w znak:\n")
		for _, termin := range w.nietykalne {
			b.WriteString("  • „" + termin.Zrodlo + "”")
			if termin.Uwaga != nil && strings.TrimSpace(*termin.Uwaga) != "" {
				b.WriteString(" (uwaga Operatora: " + strings.TrimSpace(*termin.Uwaga) + ")")
			}
			b.WriteString("\n")
		}
	}
	return b.String()
}

// zasadyJakosci składa blok zasad z rodzajów niezgodności kontraktu. Każda
// zasada jest wprost sparowana ze stałą `TranslationIssueKind`, żeby widać
// było, że polecenie i kontrola mówią o tym samym (nagłówek pliku). Rodzaj
// `omission` jest tu obecny, choć `quality.check` go nie sprawdza (wymagałby
// rozumienia treści) — polecenie może żądać rzeczy, której kontrola potem nie
// zweryfikuje; odwrotnie byłoby nieuczciwie.
func zasadyJakosci() string {
	zasady := []struct {
		rodzaj shared.TranslationIssueKind
		tresc  string
	}{
		{shared.TranslationIssueKindNumber,
			"przepisz wszystkie liczby bez zmiany wartości — nie zaokrąglaj, nie przeliczaj jednostek"},
		{shared.TranslationIssueKindDate,
			"zachowaj wszystkie daty; zapis możesz dostosować do zwyczaju języka docelowego, ale dzień, miesiąc i rok muszą zostać te same"},
		{shared.TranslationIssueKindCurrency,
			"nie przeliczaj kwot ani nie podmieniaj walut — symbol i kod waluty zostają takie, jak w źródle"},
		{shared.TranslationIssueKindPlaceholder,
			"znaczniki podstawienia w klamrach {…} przepisz znak w znak, razem z nazwą w środku — nie tłumacz ich i nie zmieniaj kolejności klamer"},
		{shared.TranslationIssueKindLength,
			"nie streszczaj i nie rozbudowuj — długość przekładu ma odpowiadać długości źródła"},
		{shared.TranslationIssueKindOmission,
			"przetłumacz KAŻDE zdanie źródła; nie pomijaj żadnego fragmentu, nawet gdy wydaje się powtórzeniem"},
	}
	var b strings.Builder
	b.WriteString("ZASADY JAKOŚCI (moduł sprawdza je po przekładzie):\n")
	for _, z := range zasady {
		b.WriteString("  • [" + string(z.rodzaj) + "] " + z.tresc + "\n")
	}
	return b.String()
}

// poleceniePrzekladu składa treść wywołania modelu dla przekładu jednego
// panelu: żądany język, ton panelu (gdy Operator go ustawił), słownik
// Operatora, zasady jakości i dopiero na końcu tekst źródłowy.
//
// Kolejność nie jest przypadkowa — polecenia idą przed tekstem, żeby model
// czytał je jako instrukcję, a nie jako część materiału do przełożenia; sam
// tekst źródłowy zamyka polecenie pod wyraźnym nagłówkiem z tego samego
// powodu. Żądanie „oddaj wyłącznie przekład" zostaje, bo treść panelu ma być
// tłumaczeniem, nie rozmową o tłumaczeniu.
func poleceniePrzekladu(jezykDocelowy string, ton *string, w wiazaniaSlownika, tekstZrodlowy string) string {
	var b strings.Builder
	b.WriteString("Przetłumacz poniższy tekst na język: " + jezykDocelowy + ".\n")
	b.WriteString("Oddaj wyłącznie sam przekład — bez komentarza, bez cudzysłowów, bez powtarzania oryginału.\n")
	if ton != nil && strings.TrimSpace(*ton) != "" {
		b.WriteString("Ton przekładu: " + strings.TrimSpace(*ton) + ".\n")
	}
	if wskazania := wskazaniaSlownika(w); wskazania != "" {
		b.WriteString("\n" + wskazania)
	}
	b.WriteString("\n" + zasadyJakosci())
	b.WriteString("\nTekst źródłowy:\n" + tekstZrodlowy)
	return b.String()
}

// polecenieTlumaczeniaZwrotnego składa treść wywołania dla `backtranslation.run`.
// Tłumaczenie zwrotne służy kontroli wierności, więc polecenie jest odwrotnością
// polecenia przekładu w jednym istotnym punkcie: słownik do niego nie wchodzi. Gdyby
// wszedł, przekład zwrotny „naprawiałby” terminologię z powrotem na brzmienie
// źródłowe i Operator zobaczyłby zgodność tam, gdzie jej nie ma — kontrola
// przestałaby cokolwiek kontrolować. Z tego samego powodu żądamy przekładu
// dosłownego, nie gładkiego.
func polecenieTlumaczeniaZwrotnego(jezykZrodlowy, tekstPanelu string) string {
	return "Przetłumacz poniższy tekst z powrotem na język: " + jezykZrodlowy + ".\n" +
		"To jest KONTROLA WIERNOŚCI przekładu: tłumacz możliwie dosłownie, oddając to, " +
		"co tekst NAPRAWDĘ mówi, a nie to, co prawdopodobnie miał mówić oryginał. " +
		"Nie poprawiaj błędów, nie wygładzaj stylu, nie uzupełniaj pominięć.\n" +
		"Oddaj wyłącznie sam przekład — bez komentarza.\n\n" +
		"Tekst do przełożenia z powrotem:\n" + tekstPanelu
}
