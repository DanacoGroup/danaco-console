// Odpowiedzialność pliku: moduł Translate — kontrola jakości panelu widoczna wewnątrz
// jego treści (`translate.quality.check`) i ustawienie tonu (`translate.panel.tone.set`).
package core

import (
	"context"
	"errors"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

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

func niezgodnoscZnacznika(szczegol string) dane.NiezgodnoscTlumaczenia {
	s := szczegol
	return dane.NiezgodnoscTlumaczenia{
		Rodzaj:   string(shared.TranslationIssueKindPlaceholder),
		Szczegol: &s,
	}
}

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

// Panel nieznany wraca odmową, nie cichą zgodą.
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
		if errors.Is(err, dane.ErrBrakWiersza) {
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
