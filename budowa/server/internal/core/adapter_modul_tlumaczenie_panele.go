// Odpowiedzialność pliku: moduł Translate — dodanie panelu docelowego,
// zapis korekty Operatora w panelu i tłumaczenie zwrotne do kontroli
// wierności. Typ adapterTlumaczenia deklaruje adapter_modul_tlumaczenie.go.
package core

import (
	"context"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// DodajPanel obsługuje target.add. Zakłada panel języka docelowego dla
// wskazanego okna i, gdy okno ma tekst źródłowy, wypełnia go przekładem
// modelu przed założeniem panelu w bazie.
func (a *adapterTlumaczenia) DodajPanel(ctx context.Context,
	z shared.TranslateTargetAddRequest) (shared.TranslateTargetAddResponse, error) {

	if z.WindowId == "" {
		return shared.TranslateTargetAddResponse{}, bladWskazaniaTlumaczenia("żądanie bez okna tłumaczenia")
	}
	if z.Language == "" {
		return shared.TranslateTargetAddResponse{}, bladWskazaniaTlumaczenia("żądanie bez języka docelowego")
	}

	okno, err := a.repozytorium.Okno(ctx, z.WindowId)
	if err != nil {
		return shared.TranslateTargetAddResponse{}, bladNieznanegoOkna(z.WindowId, err)
	}

	// Przekład liczy się przed zapisem panelu, żeby odmowa modelu nie zostawiła pustego wiersza.
	var tresc *string
	tekstZrodlowy := ""
	if okno.TekstZrodlowy != nil {
		tekstZrodlowy = strings.TrimSpace(*okno.TekstZrodlowy)
	}
	if tekstZrodlowy != "" {
		// Kanał wskazuje pole channelId; wskazanie niedobre daje odmowę nazwaną, nie zejście na domyślny.
		przeklad, err := a.przetlumaczModelem(ctx, okno.Kod, z.Language, z.Tone, tekstZrodlowy, z.ChannelId)
		if err != nil {
			return shared.TranslateTargetAddResponse{}, err
		}
		tresc = &przeklad
	}

	kod := nowyIdentyfikator(przedrostekPaneluTlumaczenia)
	panel, err := a.repozytorium.ZapiszPanel(ctx, okno.ID, dane.PanelTlumaczenia{
		Kod:   kod,
		Jezyk: z.Language,
		Tresc: tresc,
		Ton:   z.Tone,
	})
	if err != nil {
		return shared.TranslateTargetAddResponse{}, bladTlumaczenia(err)
	}

	zlozony := zlozPanelTlumaczenia(panel)

	// Panel z treścią przekładu jest zmianą, rozgłoszoną zdarzeniem translate.translation.changed.
	if tresc != nil {
		a.rozglosZmianePanelu(ctx, shared.ChangeKindCreated, panel)

		// Migawka jakości powstaje zaraz po przekładzie, bez czekania na quality.check.
		niezgodnosci := zbadajPanel(panel, tekstZrodlowy)
		_ = a.repozytorium.ZapiszNiezgodnosci(ctx, panel.ID, niezgodnosci)
		zlozony.Issues = opisyNiezgodnosci(niezgodnosci)

		// Para segmentów wchodzi do pamiec_tlumaczen, z której czyta memory.suggest.
		a.zapamietajPary(ctx, panel, tekstZrodlowy, *tresc)
	}
	return shared.TranslateTargetAddResponse{Panel: zlozony}, nil
}

// UstawTlumaczenie obsługuje translation.set — zapisuje korektę Operatora
// w panelu docelowym, drogą treści niezależną od przekładu modelu.
func (a *adapterTlumaczenia) UstawTlumaczenie(ctx context.Context,
	z shared.TranslateTranslationSetRequest) (shared.TranslateTranslationSetResponse, error) {

	if z.PanelId == "" {
		return shared.TranslateTranslationSetResponse{}, bladWskazaniaTlumaczenia("żądanie bez panelu docelowego")
	}

	tresc := z.Text
	panel, err := a.repozytorium.UstawTlumaczenie(ctx, z.PanelId, &tresc, nil)
	if err != nil {
		return shared.TranslateTranslationSetResponse{}, przelozBladPanelu(err, z.PanelId)
	}
	a.rozglosZmianePanelu(ctx, shared.ChangeKindUpdated, panel)

	// Korekta Operatora także idzie do pamięci tłumaczeń, jeśli okno panelu ma tekst źródłowy.
	if okno, err := a.repozytorium.Okno(ctx, panel.OknoKod); err == nil && okno.TekstZrodlowy != nil {
		a.zapamietajPary(ctx, panel, *okno.TekstZrodlowy, z.Text)
	}
	return shared.TranslateTranslationSetResponse{Panel: zlozPanelTlumaczenia(panel)}, nil
}

// TlumaczZwrotnie obsługuje backtranslation.run — przekład zwrotny wykonuje
// model. Cztery odmowy: panel bez treści, okno bez języka źródłowego, brak
// czynnego kanału modelu, kanał wskazany polem channelId nieczynny.
func (a *adapterTlumaczenia) TlumaczZwrotnie(ctx context.Context,
	z shared.TranslateBacktranslationRunRequest) (shared.TranslateBacktranslationRunResponse, error) {

	if z.PanelId == "" {
		return shared.TranslateBacktranslationRunResponse{}, bladWskazaniaTlumaczenia("żądanie bez panelu docelowego")
	}

	panel, err := a.repozytorium.Panel(ctx, z.PanelId)
	if err != nil {
		return shared.TranslateBacktranslationRunResponse{}, przelozBladPanelu(err, z.PanelId)
	}
	if panel.Tresc == nil || strings.TrimSpace(*panel.Tresc) == "" {
		return shared.TranslateBacktranslationRunResponse{}, bladWskazaniaTlumaczenia(
			"panel " + z.PanelId + " nie ma jeszcze treści — nie ma czego tłumaczyć z powrotem")
	}

	okno, err := a.repozytorium.Okno(ctx, panel.OknoKod)
	if err != nil {
		return shared.TranslateBacktranslationRunResponse{}, bladNieznanegoOkna(panel.OknoKod, err)
	}
	jezykZrodlowy := ""
	if okno.JezykZrodlowy != nil {
		jezykZrodlowy = strings.TrimSpace(*okno.JezykZrodlowy)
	}
	if jezykZrodlowy == "" {
		return shared.TranslateBacktranslationRunResponse{}, bladWskazaniaTlumaczenia(
			"okno " + panel.OknoKod + " nie ma języka źródłowego — kontrola wierności nie wie, na jaki język " +
				"przełożyć panel z powrotem; ustaw język źródłowy (translate.source.set) albo rozpoznaj go " +
				"komendą translate.source.detect")
	}

	zwrotne, err := a.przetlumaczZwrotnieModelem(ctx, panel.OknoKod, jezykZrodlowy, *panel.Tresc, z.ChannelId)
	if err != nil {
		return shared.TranslateBacktranslationRunResponse{}, err
	}

	zaktualizowany, err := a.repozytorium.UstawTlumaczenieZwrotne(ctx, z.PanelId, zwrotne)
	if err != nil {
		return shared.TranslateBacktranslationRunResponse{}, przelozBladPanelu(err, z.PanelId)
	}
	tekst := ""
	if zaktualizowany.TrescZwrotna != nil {
		tekst = *zaktualizowany.TrescZwrotna
	}
	return shared.TranslateBacktranslationRunResponse{PanelId: z.PanelId, Text: tekst}, nil
}

// przelozBladPanelu tłumaczy ErrBrakWiersza na odmowę zrozumiałą dla
// Operatora przy operacjach adresowanych panelem docelowym — wspólne dla
// trzech metod tego pliku, żeby komunikat nie rozjeżdżał się między nimi.
func przelozBladPanelu(err error, kodPanelu string) error {
	return bladNieznanegoPanelu(kodPanelu, err)
}
