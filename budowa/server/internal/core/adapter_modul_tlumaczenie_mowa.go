// Plik obsługuje moduł Translate: `speech.synthesize` i `panel.export` na `*adapterTlumaczenia`. Typ i konstruktor deklaruje `adapter_modul_tlumaczenie.go`. Mowę syntezuje `espeak-ng` portem `session.Uruchamiacz`.
package core

import (
	"context"
	"errors"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// SyntezujMowe obsługuje `speech.synthesize`. Czyta treść panelu, oddaje ją syntezatorowi i zwraca ścieżkę pliku dźwiękowego. Trzy odmowy wprost: żądanie bez panelu, panel bez treści, panel bez języka.
func (a *adapterTlumaczenia) SyntezujMowe(ctx context.Context,
	z shared.TranslateSpeechSynthesizeRequest) (shared.TranslateSpeechSynthesizeResponse, error) {

	if z.PanelId == "" {
		return shared.TranslateSpeechSynthesizeResponse{}, bladWskazaniaTlumaczenia("speech.synthesize bez panelu")
	}

	panel, err := a.repozytorium.Panel(ctx, z.PanelId)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.TranslateSpeechSynthesizeResponse{}, bladNieznanegoPanelu(z.PanelId, err)
	}
	if err != nil {
		return shared.TranslateSpeechSynthesizeResponse{}, bladTlumaczenia(err)
	}

	tresc := ""
	if panel.Tresc != nil {
		tresc = strings.TrimSpace(*panel.Tresc)
	}
	if tresc == "" {
		return shared.TranslateSpeechSynthesizeResponse{}, bladWskazaniaTlumaczenia(
			"panel " + z.PanelId + " nie ma treści — nie ma czego przeczytać na głos")
	}
	jezyk := strings.TrimSpace(panel.Jezyk)
	if jezyk == "" {
		return shared.TranslateSpeechSynthesizeResponse{}, bladWskazaniaTlumaczenia(
			"panel " + z.PanelId + " nie ma języka — syntezator nie wie, którym głosem czytać")
	}

	sciezka, err := a.zsyntezujDoPliku(ctx, panel.Kod, jezyk, tresc)
	if err != nil {
		return shared.TranslateSpeechSynthesizeResponse{}, err
	}

	// Ślad zapisuje się po syntezie, z odnośnikiem do pliku, który istnieje.
	_, _ = a.repozytorium.ZapiszSyntezeMowy(ctx, dane.SyntezaMowy{PanelID: panel.ID, NagranieOdnosnik: &sciezka})

	return shared.TranslateSpeechSynthesizeResponse{PanelId: z.PanelId, Path: sciezka}, nil
}

// EksportujPanel obsługuje `panel.export`. Zapisuje ślad eksportu panelu w formacie żądanym przez Operatora; `Path` odpowiedzi zostaje pusty, bo rdzeń nie ma magazynu blobów.
func (a *adapterTlumaczenia) EksportujPanel(ctx context.Context,
	z shared.TranslatePanelExportRequest) (shared.TranslatePanelExportResponse, error) {

	if z.PanelId == "" {
		return shared.TranslatePanelExportResponse{}, bladWskazaniaTlumaczenia("panel.export bez panelu")
	}

	panel, err := a.repozytorium.Panel(ctx, z.PanelId)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.TranslatePanelExportResponse{}, bladNieznanegoPanelu(z.PanelId, err)
	}
	if err != nil {
		return shared.TranslatePanelExportResponse{}, bladTlumaczenia(err)
	}

	eksport := dane.EksportPanelu{PanelID: panel.ID, Format: string(z.Format)}
	if _, err := a.repozytorium.ZapiszEksportPanelu(ctx, eksport); err != nil {
		return shared.TranslatePanelExportResponse{}, bladTlumaczenia(err)
	}

	return shared.TranslatePanelExportResponse{PanelId: z.PanelId}, nil
}
