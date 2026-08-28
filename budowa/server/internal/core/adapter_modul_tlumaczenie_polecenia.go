// Odpowiedzialność pliku: moduł Translate — składanie poleceń dla modelu i wiązanie
// słownika Operatora z przekładem, bo cała wiedza o dobrym przekładzie mieści się w poleceniu.
package core

import (
	"context"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// wiazaniaSlownika to słownik Operatora zawężony do jednego przekładu: odpowiedniki
// obowiązujące w języku docelowym i nazwy nietykalne obowiązujące we wszystkich językach.
type wiazaniaSlownika struct {
	odpowiedniki []dane.TerminSlownika
	nietykalne   []dane.TerminSlownika
}

// puste mówi, czy słownik nie wniósł do tego przekładu niczego — wtedy polecenie nie
// niesie bloku słownika w ogóle, żeby nie rozpraszać modelu pustym nagłówkiem.
func (w wiazaniaSlownika) puste() bool {
	return len(w.odpowiedniki) == 0 && len(w.nietykalne) == 0
}

// wiazaniaSlownikaDlaJezyka czyta słownik Operatora i zawęża go do przekładu na wskazany
// język; błąd odczytu oddaje wprost, zamiast po cichu ignorować terminologię.
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
			// Termin bez odpowiednika nie niesie żadnego polecenia dla modelu.
			continue
		}
		wynik.odpowiedniki = append(wynik.odpowiedniki, termin)
	}
	return wynik, nil
}

// wskazaniaSlownika przekłada zawężony słownik na blok polecenia; uwaga Operatora idzie
// do modelu razem z odpowiednikiem, bo bywa, że niesie powód wyboru słowa.
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

// zasadyJakosci składa blok zasad z rodzajów niezgodności kontraktu; każda zasada jest
// wprost sparowana ze stałą kontraktu, żeby polecenie i kontrola mówiły o tym samym.
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

// poleceniePrzekladu składa treść wywołania modelu dla przekładu jednego panelu: żądany
// język, ton, słownik Operatora, zasady jakości i dopiero na końcu tekst źródłowy.
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

// polecenieTlumaczeniaZwrotnego składa treść wywołania dla `backtranslation.run`; słownik
// do niego nie wchodzi, bo tłumaczenie zwrotne służy kontroli wierności przekładu.
func polecenieTlumaczeniaZwrotnego(jezykZrodlowy, tekstPanelu string) string {
	return "Przetłumacz poniższy tekst z powrotem na język: " + jezykZrodlowy + ".\n" +
		"To jest KONTROLA WIERNOŚCI przekładu: tłumacz możliwie dosłownie, oddając to, " +
		"co tekst NAPRAWDĘ mówi, a nie to, co prawdopodobnie miał mówić oryginał. " +
		"Nie poprawiaj błędów, nie wygładzaj stylu, nie uzupełniaj pominięć.\n" +
		"Oddaj wyłącznie sam przekład — bez komentarza.\n\n" +
		"Tekst do przełożenia z powrotem:\n" + tekstPanelu
}
