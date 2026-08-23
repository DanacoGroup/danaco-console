// Odpowiedzialność pliku: szablony promptu strukturalnego
// (`design.prompt.template.save`, `design.prompt.template.list`) oraz historia
// promptów wydanych w oknie (`design.prompt.history.list`). Metody stoją na
// `*adapterDesignu` (`adapter_modul_design.go`).
//
// Do tej pory historia promptów i szablony żyły w oknie do zamknięcia karty
// przeglądarki i tyle o nich było wiadomo. Prompt Builder wypracowywał
// polecenie, Operator zamykał kartę i praca znikała.
//
// Historia domyka też prowenancję. Zasób niesie pole `promptId`, którego rdzeń
// dotąd NIE wypełniał, bo nie było przekładu klucza wiersza promptu na kod
// kontraktu. Odczyt zasobu bierze teraz kod promptu podzapytaniem
// (`dane/design_zasoby.go`), a ta komenda oddaje drugą stronę tego samego
// powiązania: prompt wraz z kodami zasobów, które z niego powstały.
package core

import (
	"context"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// przedrostekSzablonuPromptuDesign znakuje identyfikatory zewnętrzne szablonów.
const przedrostekSzablonuPromptuDesign = "szablon-promptu-"

// ZapiszSzablonPromptu utrwala prompt jako szablon do wielokrotnego użycia —
// obsługuje `design.prompt.template.save`.
//
// Nadpisanie idzie po `templateId` z żądania; brak zakłada szablon nowy.
// Szablon wskazany a nieznany jest ODMOWĄ, nie cichym założeniem nowego:
// Operator, który nadpisuje szablon, oczekuje że nadpisał ten jeden, a nie że
// dostał drugi obok — a wykaz z dwoma szablonami tej samej nazwy wygląda
// dokładnie tak, jak wygląda zgubiona praca.
func (a *adapterDesignu) ZapiszSzablonPromptu(ctx context.Context,
	z shared.DesignPromptTemplateSaveRequest) (shared.DesignPromptTemplateSaveResponse, error) {

	if z.WindowId == "" {
		return shared.DesignPromptTemplateSaveResponse{}, bladWskazaniaDesignu(
			"komenda design.prompt.template.save bez wskazania okna")
	}
	if strings.TrimSpace(z.Name) == "" {
		return shared.DesignPromptTemplateSaveResponse{}, bladWskazaniaDesignu(
			"komenda design.prompt.template.save bez nazwy szablonu")
	}
	if strings.TrimSpace(z.Prompt.Subject) == "" {
		return shared.DesignPromptTemplateSaveResponse{}, bladWskazaniaDesignu(
			"komenda design.prompt.template.save z promptem bez tematu — szablon bez tematu " +
				"nie wygeneruje niczego, gdy Operator po niego sięgnie")
	}

	kod := nowyIdentyfikator(przedrostekSzablonuPromptuDesign)
	if z.TemplateId != nil && strings.TrimSpace(*z.TemplateId) != "" {
		kod = strings.TrimSpace(*z.TemplateId)
		if err := a.sprawdzSzablonPromptuDesignu(ctx, z.WindowId, kod); err != nil {
			return shared.DesignPromptTemplateSaveResponse{}, err
		}
	}

	zapisany, err := a.repozytorium.ZapiszSzablonPromptuDesignu(ctx, dane.SzablonPromptuDesignu{
		Kod:         kod,
		Okno:        z.WindowId,
		Nazwa:       strings.TrimSpace(z.Name),
		Temat:       z.Prompt.Subject,
		Styl:        z.Prompt.Style,
		Kompozycja:  z.Prompt.Composition,
		Oswietlenie: z.Prompt.Lighting,
		Paleta:      z.Prompt.Palette,
		Proporcje:   z.Prompt.AspectRatio,
		Wykluczenia: z.Prompt.Exclusions,
		Ziarno:      z.Prompt.Seed,
		Warianty:    z.Prompt.Variants,
		Silnik:      z.Prompt.Engine,
		Kreatywnosc: z.Prompt.Creativity,
	})
	if err != nil {
		return shared.DesignPromptTemplateSaveResponse{}, bladDesignu(err)
	}
	return shared.DesignPromptTemplateSaveResponse{
		Template: szablonPromptuKontraktuDesignu(zapisany),
	}, nil
}

// SzablonyPromptu zwraca szablony promptów okna — obsługuje
// `design.prompt.template.list`.
func (a *adapterDesignu) SzablonyPromptu(ctx context.Context,
	z shared.DesignPromptTemplateListRequest) (shared.DesignPromptTemplateListResponse, error) {

	if z.WindowId == "" {
		return shared.DesignPromptTemplateListResponse{}, bladWskazaniaDesignu(
			"komenda design.prompt.template.list bez wskazania okna")
	}
	wiersze, err := a.repozytorium.SzablonyPromptuDesignu(ctx, z.WindowId)
	if err != nil {
		return shared.DesignPromptTemplateListResponse{}, bladDesignu(err)
	}
	szablony := make([]shared.DesignPromptTemplate, 0, len(wiersze))
	for _, wiersz := range wiersze {
		szablony = append(szablony, szablonPromptuKontraktuDesignu(wiersz))
	}
	return shared.DesignPromptTemplateListResponse{Templates: szablony, Total: len(szablony)}, nil
}

// HistoriaPromptow zwraca prompty wydane w oknie wraz z zasobami, które z nich
// powstały — obsługuje `design.prompt.history.list`.
//
// Prompt bez ani jednego zasobu zostaje w historii. Kanał bywa odmawiał
// i wtedy prompt jest zapisem próby — czyli dokładnie tego, po co Operator do
// historii sięga, żeby poprawić polecenie i wydać je jeszcze raz.
func (a *adapterDesignu) HistoriaPromptow(ctx context.Context,
	z shared.DesignPromptHistoryListRequest) (shared.DesignPromptHistoryListResponse, error) {

	if z.WindowId == "" {
		return shared.DesignPromptHistoryListResponse{}, bladWskazaniaDesignu(
			"komenda design.prompt.history.list bez wskazania okna")
	}
	wiersze, razem, err := a.repozytorium.PromptyOknaDesignu(ctx, z.WindowId, wartoscLiczby(z.Limit))
	if err != nil {
		return shared.DesignPromptHistoryListResponse{}, bladDesignu(err)
	}
	zapisy := make([]shared.DesignPromptRecord, 0, len(wiersze))
	for _, wiersz := range wiersze {
		zasoby, err := a.repozytorium.ZasobyPromptuDesignu(ctx, wiersz.ID)
		if err != nil {
			return shared.DesignPromptHistoryListResponse{}, bladDesignu(err)
		}
		zapisy = append(zapisy, shared.DesignPromptRecord{
			Id:        wiersz.Kod,
			WindowId:  wiersz.Okno,
			Prompt:    promptKontraktuDesignu(wiersz),
			AssetIds:  zasoby,
			ChannelId: wiersz.Kanal,
			CreatedAt: chwilaBazy(wiersz.Utworzono),
		})
	}
	return shared.DesignPromptHistoryListResponse{Prompts: zapisy, Total: razem}, nil
}

// sprawdzSzablonPromptuDesignu odmawia nadpisania szablonu, którego nie ma,
// i szablonu należącego do innego okna. Drugi warunek nie jest formalnością:
// szablony są bytem okna, a nadpisanie cudzego przez podanie jego kodu byłoby
// zmianą stanu, którego wołający nawet nie widzi.
func (a *adapterDesignu) sprawdzSzablonPromptuDesignu(ctx context.Context, okno, kod string) error {
	szablony, err := a.repozytorium.SzablonyPromptuDesignu(ctx, okno)
	if err != nil {
		return bladDesignu(err)
	}
	for _, szablon := range szablony {
		if szablon.Kod == kod {
			return nil
		}
	}
	return bladNieznanegoBytuDesignu("szablonu promptu " + kod + " nie ma w oknie " + okno)
}

// szablonPromptuKontraktuDesignu składa `DesignPromptTemplate` kontraktu
// z wiersza repozytorium.
func szablonPromptuKontraktuDesignu(s dane.SzablonPromptuDesignu) shared.DesignPromptTemplate {
	return shared.DesignPromptTemplate{
		Id:       s.Kod,
		WindowId: s.Okno,
		Name:     s.Nazwa,
		Prompt: shared.DesignPrompt{
			Id:          &s.Kod,
			Subject:     s.Temat,
			Style:       s.Styl,
			Composition: s.Kompozycja,
			Lighting:    s.Oswietlenie,
			Palette:     s.Paleta,
			AspectRatio: s.Proporcje,
			Exclusions:  s.Wykluczenia,
			Seed:        s.Ziarno,
			Variants:    s.Warianty,
			Engine:      s.Silnik,
			Creativity:  s.Kreatywnosc,
		},
		UpdatedAt: chwilaBazy(s.Zaktualizowano),
	}
}

// promptKontraktuDesignu składa `DesignPrompt` kontraktu z wiersza promptu
// wydanego. Pole `Id` niesie identyfikator ZEWNĘTRZNY — ten sam, który zasób
// niesie w `promptId`, więc panel złoży jedno z drugim bez zgadywania.
func promptKontraktuDesignu(p dane.PromptDesignu) shared.DesignPrompt {
	kod := p.Kod
	return shared.DesignPrompt{
		Id:          &kod,
		Subject:     p.Temat,
		Style:       p.Styl,
		Composition: p.Kompozycja,
		Lighting:    p.Oswietlenie,
		Palette:     p.Paleta,
		AspectRatio: p.Proporcje,
		Exclusions:  p.Wykluczenia,
		Seed:        p.Ziarno,
		Variants:    p.Warianty,
		Engine:      p.Silnik,
		Creativity:  p.Kreatywnosc,
	}
}
