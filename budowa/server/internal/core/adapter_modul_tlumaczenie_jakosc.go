// Odpowiedzialność pliku: moduł Translate — kontrola jakości panelu
// (`translate.quality.check`) i ustawienie tonu (`translate.panel.tone.set`).
// Metody na `*adapterTlumaczenia`, którego typ, konstruktor, przedrostki
// i wspólne pomocniki błędów deklaruje `adapter_modul_tlumaczenie.go`.
//
// Kontrakt zna sześć rodzajów niezgodności (`TranslationIssueKind`): number,
// date, currency, placeholder, length, omission. Kontrole wymagające porównania
// treści panelu z tekstem źródłowym okna — `number`, `currency` i pełne
// `placeholder` — leżą w `adapter_modul_tlumaczenie_jakosc_zrodlo.go`
// (`zbadajPanel`), razem z powodem, dla którego `date` i `omission` sprawdzane
// nie są.
//
// Ten plik liczy dwa rodzaje, widoczne wewnątrz samej treści panelu:
//   - `placeholder`: znacznik podstawienia `{...}` niedomknięty albo domknięty
//     bez otwarcia — usterka składniowa, nie znaczeniowa (`sprawdzZnaczniki`).
//   - `length`: panel oznaczony jako `ready`, a jego treść pusta albo samą
//     białą spacją — długość zero tam, gdzie stan twierdzi, że tłumaczenie
//     istnieje (`sprawdzDlugosc`).
//
// Kontrola jest migawką: `ZapiszNiezgodnosci` podmienia komplet wierszy panelu
// w jednej transakcji, więc wykaz pusty znaczy panel czysty, a nie brak zapisu.
package core

import (
	"context"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// SprawdzJakosc obsługuje `translate.quality.check`. Kontrola porównuje panel
// ze źródłem: tekst źródłowy jest osiągalny z panelu przez kod zewnętrzny okna
// (`tekstZrodlowyPanelu`), więc `number`, `currency` i pełne sprawdzenie
// `placeholder` są liczone w `zbadajPanel`. Rodzajów `date` i `omission`
// kontrola nie liczy — powody stoją w nagłówku
// `adapter_modul_tlumaczenie_jakosc_zrodlo.go`.
func (a *adapterTlumaczenia) SprawdzJakosc(ctx context.Context,
	z shared.TranslateQualityCheckRequest) (shared.TranslateQualityCheckResponse, error) {

	if z.PanelId == "" {
		return shared.TranslateQualityCheckResponse{}, bladWskazaniaTlumaczenia("żądanie bez panelu do kontroli")
	}

	panel, err := a.repozytorium.Panel(ctx, z.PanelId)
	if err != nil {
		return shared.TranslateQualityCheckResponse{}, bladNieznanegoPanelu(z.PanelId, err)
	}

	zrodlo, err := a.tekstZrodlowyPanelu(ctx, panel)
	if err != nil {
		return shared.TranslateQualityCheckResponse{}, err
	}
	niezgodnosci := zbadajPanel(panel, zrodlo)

	if err := a.repozytorium.ZapiszNiezgodnosci(ctx, panel.ID, niezgodnosci); err != nil {
		return shared.TranslateQualityCheckResponse{}, bladTlumaczenia(err)
	}

	return shared.TranslateQualityCheckResponse{
		PanelId: panel.Kod,
		Issues:  zdaniaNiezgodnosci(niezgodnosci),
	}, nil
}

// sprawdzZnaczniki liczy niezgodności rodzaju `placeholder`: znaczniki
// podstawienia `{...}` w treści panelu, które się nie domykają, albo które
// domykają nawias bez wcześniejszego otwarcia. To jedyna niezgodność
// znaczników stwierdzalna bez tekstu źródłowego; porównanie znaczników panelu
// ze znacznikami źródła liczy `zbadajPanel`.
func sprawdzZnaczniki(panel dane.PanelTlumaczenia) []dane.NiezgodnoscTlumaczenia {
	if panel.Tresc == nil {
		return nil
	}
	var wynik []dane.NiezgodnoscTlumaczenia
	otwarte := 0
	for _, r := range *panel.Tresc {
		switch r {
		case '{':
			otwarte++
		case '}':
			if otwarte == 0 {
				wynik = append(wynik, niezgodnoscZnacznika("znacznik podstawienia zamknięty bez otwarcia"))
				continue
			}
			otwarte--
		}
	}
	if otwarte > 0 {
		wynik = append(wynik, niezgodnoscZnacznika("znacznik podstawienia otwarty i niedomknięty"))
	}
	return wynik
}

// niezgodnoscZnacznika składa jeden wiersz niezgodności rodzaju `placeholder`.
func niezgodnoscZnacznika(szczegol string) dane.NiezgodnoscTlumaczenia {
	s := szczegol
	return dane.NiezgodnoscTlumaczenia{
		Rodzaj:   string(shared.TranslationIssueKindPlaceholder),
		Szczegol: &s,
	}
}

// sprawdzDlugosc liczy niezgodność rodzaju `length`: panel oznaczony jako
// gotowy (`ready`), którego treść jest pusta albo samą białą spacją — długość
// zero tam, gdzie stan panelu twierdzi, że tłumaczenie istnieje. To jedyny
// przypadek długości stwierdzalny z samej treści panelu.
func sprawdzDlugosc(panel dane.PanelTlumaczenia) []dane.NiezgodnoscTlumaczenia {
	if panel.Stan != string(shared.TranslationStatusReady) {
		return nil
	}
	if panel.Tresc != nil && strings.TrimSpace(*panel.Tresc) != "" {
		return nil
	}
	szczegol := "panel oznaczony jako gotowy, ale treść tłumaczenia jest pusta"
	return []dane.NiezgodnoscTlumaczenia{{
		Rodzaj:   string(shared.TranslationIssueKindLength),
		Szczegol: &szczegol,
	}}
}

// zdaniaNiezgodnosci składa niezgodności w wykaz zdań, bo `quality.check`
// oddaje w kontrakcie `[]string`, a nie `[]TranslationIssue`. Te same fakty
// wychodzą z rdzenia w dwóch kształtach: rozłożonym na pola przy panelu
// i zwięzłym przy samej kontroli.
func zdaniaNiezgodnosci(niezgodnosci []dane.NiezgodnoscTlumaczenia) []string {
	wynik := make([]string, 0, len(niezgodnosci))
	for _, n := range niezgodnosci {
		zdanie := string(n.Rodzaj)
		if n.Segment != nil && *n.Segment != "" {
			zdanie += " · segment: " + *n.Segment
		}
		if n.Szczegol != nil && *n.Szczegol != "" {
			zdanie += " — " + *n.Szczegol
		}
		wynik = append(wynik, zdanie)
	}
	return wynik
}

// opisyNiezgodnosci przekłada wiersze repozytorium na byty kontraktu.
func opisyNiezgodnosci(niezgodnosci []dane.NiezgodnoscTlumaczenia) []shared.TranslationIssue {
	wynik := make([]shared.TranslationIssue, 0, len(niezgodnosci))
	for _, n := range niezgodnosci {
		wynik = append(wynik, shared.TranslationIssue{
			Kind:    shared.TranslationIssueKind(n.Rodzaj),
			Segment: n.Segment,
			Detail:  n.Szczegol,
		})
	}
	return wynik
}

// UstawTon obsługuje `translate.panel.tone.set`. Ton nie ma osobnej tabeli —
// zmienia kolumnę `ton` panelu przez `UstawTon` warstwy danych. Panel nieznany
// wraca odmową (`bladNieznanegoPanelu`), nie cichą zgodą.
func (a *adapterTlumaczenia) UstawTon(ctx context.Context,
	z shared.TranslatePanelToneSetRequest) (shared.TranslatePanelToneSetResponse, error) {

	if z.PanelId == "" {
		return shared.TranslatePanelToneSetResponse{}, bladWskazaniaTlumaczenia("żądanie bez panelu")
	}
	if z.Tone == "" {
		return shared.TranslatePanelToneSetResponse{}, bladWskazaniaTlumaczenia("żądanie bez tonu tłumaczenia")
	}

	panel, err := a.repozytorium.UstawTon(ctx, z.PanelId, z.Tone)
	if err != nil {
		if err == dane.ErrBrakWiersza {
			return shared.TranslatePanelToneSetResponse{}, bladNieznanegoPanelu(z.PanelId, err)
		}
		return shared.TranslatePanelToneSetResponse{}, bladTlumaczenia(err)
	}

	niezgodnosci, err := a.repozytorium.Niezgodnosci(ctx, panel.ID)
	if err != nil {
		return shared.TranslatePanelToneSetResponse{}, bladTlumaczenia(err)
	}

	zlozony := zlozPanelTlumaczenia(panel)
	zlozony.Issues = opisyNiezgodnosci(niezgodnosci)
	return shared.TranslatePanelToneSetResponse{Panel: zlozony}, nil
}
