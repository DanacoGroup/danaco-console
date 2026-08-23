// Odpowiedzialność pliku: moduł Translate — dodanie panelu docelowego
// (`target.add`), zapis korekty Operatora w panelu (`translation.set`)
// i tłumaczenie zwrotne do kontroli wierności (`backtranslation.run`). Typ
// `adapterTlumaczenia`, konstruktor, przedrostki identyfikatorów i wspólne
// pomocniki błędów (`bladTlumaczenia`, `bladWskazaniaTlumaczenia`,
// `bladNieznanegoPanelu`) deklaruje adapter_modul_tlumaczenie.go — ten plik
// dokłada wyłącznie własne metody na tym samym typie.
//
// Przekład wykonuje model. Gdy okno ma tekst źródłowy, `target.add` przekłada
// go na język panelu, a wynik ląduje w kolumnie `tresc`. Korekta Operatora
// (`translation.set`) jest drugą drogą treści — poprawką przekładu modelu.
//
// Przekład jest związany słownikiem i zasadami jakości: do polecenia dla modelu
// wchodzą terminy Operatora, jego zakazy tłumaczenia, ton panelu i zasady
// jakości wywiedzione z rodzajów niezgodności kontraktu
// (adapter_modul_tlumaczenie_polecenia.go), a wynik przechodzi jeszcze
// mechaniczną podmianę terminów i migawkę kontroli jakości.
package core

import (
	"context"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// DodajPanel obsługuje `target.add`. Zakłada panel języka docelowego dla
// wskazanego okna i — gdy okno ma tekst źródłowy — wypełnia go przekładem
// modelu (`przetlumaczModelem`, adapter_modul_tlumaczenie_model.go). Przekład
// idzie przed założeniem panelu: nieudane wywołanie modelu (brak czynnego
// kanału, pusta odpowiedź) kończy się odmową i nie zostawia w bazie pustego
// panelu. Okno bez tekstu źródłowego daje panel bez treści; treść dołoży
// korekta Operatora przez `translation.set`.
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

	// Przekład liczy się przed zapisem panelu, żeby odmowa modelu nie zostawiła
	// pustego wiersza w bazie. Panel wychodzi od razu z treścią (albo bez niej,
	// gdy okno nie ma jeszcze tekstu źródłowego).
	var tresc *string
	tekstZrodlowy := ""
	if okno.TekstZrodlowy != nil {
		tekstZrodlowy = strings.TrimSpace(*okno.TekstZrodlowy)
	}
	if tekstZrodlowy != "" {
		// Kanał wskazuje pole `channelId`; jego brak bierze kanał domyślny
		// czynny, a wskazanie niedobre — odmowę nazwaną, nie ciche zejście na
		// domyślny (`kanalZadania`).
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

	// Panel z treścią przekładu jest zmianą, którą inne połączenia okna mają
	// zobaczyć — rozgłoszona jednym zdarzeniem `translate.translation.changed`.
	if tresc != nil {
		a.rozglosZmianePanelu(shared.ChangeKindCreated, panel)

		// Migawka jakości powstaje zaraz po przekładzie. Zasady jakości poszły
		// do polecenia (`zasadyJakosci`), lecz model gubi znaczniki i liczby
		// mimo zakazu, więc świeży panel wychodzi z wykazem zastrzeżeń już
		// wypełnionym, bez czekania na `quality.check`.
		//
		// Nieudany zapis niezgodności nie przewraca całej komendy: przekład jest
		// zapisany, a migawkę kontroli można powtórzyć komendą `quality.check`.
		// Zastrzeżenia trafiają do odpowiedzi niezależnie od wyniku zapisu.
		niezgodnosci := zbadajPanel(panel, tekstZrodlowy)
		_ = a.repozytorium.ZapiszNiezgodnosci(ctx, panel.ID, niezgodnosci)
		zlozony.Issues = opisyNiezgodnosci(niezgodnosci)

		// Para (segment źródłowy, segment przekładu) wchodzi do
		// `pamiec_tlumaczen`, z której czyta `memory.suggest`. Sparowanie i jego
		// granice: adapter_modul_tlumaczenie_pamiec.go.
		a.zapamietajPary(ctx, panel, tekstZrodlowy, *tresc)
	}
	return shared.TranslateTargetAddResponse{Panel: zlozony}, nil
}

// UstawTlumaczenie obsługuje `translation.set` — zapisuje korektę Operatora
// w panelu docelowym. Jest to droga treści niezależna od przekładu modelu
// z `target.add`. Treść trafia wprost do kolumny `tresc`; kontrakt oddaje jeden
// napis, nie odwołanie do pliku, więc ta warstwa nie rozstrzyga o pliku dla
// treści obszernej.
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
	a.rozglosZmianePanelu(shared.ChangeKindUpdated, panel)

	// Korekta Operatora także idzie do pamięci tłumaczeń. Tekst źródłowy bierze
	// się z okna panelu; okno nieosiągalne albo bez tekstu źródłowego znaczy
	// „nie ma z czym parować" i pamięć zostaje bez wiersza, co nie unieważnia
	// zapisanej już korekty.
	if okno, err := a.repozytorium.Okno(ctx, panel.OknoKod); err == nil && okno.TekstZrodlowy != nil {
		a.zapamietajPary(ctx, panel, *okno.TekstZrodlowy, z.Text)
	}
	return shared.TranslateTranslationSetResponse{Panel: zlozPanelTlumaczenia(panel)}, nil
}

// TlumaczZwrotnie obsługuje `backtranslation.run` — przekład zwrotny wykonuje
// model. Język źródłowy bierze się z okna wskazanego przez wiersz panelu
// (`PanelTlumaczenia.OknoKod`).
//
// Cztery odmowy, każda o czym innym:
//  1. panel bez treści — nie ma czego tłumaczyć z powrotem;
//  2. okno bez rozpoznanego języka źródłowego — nie wiadomo, na jaki język
//     przełożyć; języka źródłowego nie zgadujemy z panelu, od tego jest
//     `source.detect` wołany osobno;
//  3. brak czynnego kanału modelu — odmowa z `przetlumaczZwrotnieModelem`;
//  4. kanał wskazany polem `channelId`, którego nie ma albo który jest
//     nieczynny — odmowa nazwana z `kanalZadania`, nie ciche zejście na kanał
//     domyślny; kontrola wierności wykonana innym modelem sprawdzałaby co
//     innego, niż wskazano.
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
