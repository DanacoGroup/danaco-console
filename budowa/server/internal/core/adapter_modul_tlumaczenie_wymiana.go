// Odpowiedzialność pliku: wymiana glosariusza modułu Translate — import odmawiany wprost,
// eksport zapisujący ślad zlecenia oraz podpowiedzi z pamięci dopasowaniem podciągu.
package core

import (
	"context"
	"errors"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// limitPodpowiedziPamieci ogranicza liczbę podpowiedzi zwracanych z jednego zapytania, bo
// kontrakt nie niesie własnego limitu.
const limitPodpowiedziPamieci = 5

// ImportujSlownik obsługuje `translate.glossary.import` i odmawia wprost kodem
// `channel_unavailable`, bo rdzeń nie ma dostępu do dysku Operatora, więc nie ma czym wczytać terminów.
func (a *adapterTlumaczenia) ImportujSlownik(ctx context.Context,
	z shared.TranslateGlossaryImportRequest) (shared.TranslateGlossaryImportResponse, error) {

	if z.Path == "" {
		return shared.TranslateGlossaryImportResponse{}, bladWskazaniaTlumaczenia(
			"żądanie translate.glossary.import bez ścieżki pliku")
	}

	// Zapis śladu próby daje Operatorowi dowód, że zlecenie dotarło do rdzenia.
	if _, err := a.repozytorium.ZapiszImportSlownika(ctx, dane.SladImportuSlownika{
		Kod:                   nowyIdentyfikator(przedrostekImportuSlownika),
		Sciezka:               z.Path,
		LiczbaZaimportowanych: 0,
	}); err != nil {
		return shared.TranslateGlossaryImportResponse{}, bladTlumaczenia(err)
	}

	return shared.TranslateGlossaryImportResponse{}, protocol.JakoError(protocol.NowyBlad(
		shared.ErrorCodeChannelUnavailable,
		"moduł Translate: rdzeń nie czyta plików z dysku Operatora — translate.glossary.import "+
			"nie ma czym wczytać terminów spod ścieżki "+z.Path+
			"; kontrakt nie niesie terminów wprost w żądaniu, więc nie ma innego uczciwego źródła"))
}

// EksportujSlownik obsługuje `translate.glossary.export` i zapisuje ślad zlecenia, bo rdzeń
// nie ma magazynu blobów, więc żaden plik nie powstaje.
func (a *adapterTlumaczenia) EksportujSlownik(ctx context.Context,
	z shared.TranslateGlossaryExportRequest) (shared.TranslateGlossaryExportResponse, error) {

	if z.Path == "" {
		return shared.TranslateGlossaryExportResponse{}, bladWskazaniaTlumaczenia(
			"żądanie translate.glossary.export bez ścieżki pliku")
	}

	terminy, err := a.repozytorium.Terminy(ctx)
	if err != nil {
		return shared.TranslateGlossaryExportResponse{}, bladTlumaczenia(err)
	}
	liczba := int64(len(terminy))

	if _, err := a.repozytorium.ZapiszEksportSlownika(ctx, dane.SladEksportuSlownika{
		Kod:                    nowyIdentyfikator(przedrostekEksportuSlownika),
		Sciezka:                z.Path,
		LiczbaWyeksportowanych: liczba,
	}); err != nil {
		return shared.TranslateGlossaryExportResponse{}, bladTlumaczenia(err)
	}

	return shared.TranslateGlossaryExportResponse{ExportedCount: len(terminy)}, nil
}

// Podpowiedzi obsługuje `translate.memory.suggest`, ustalając język docelowy z panelu i
// pytając repozytorium o dopasowanie podciągu; panel nieznany wraca odmową, nie pustą listą.
func (a *adapterTlumaczenia) Podpowiedzi(ctx context.Context,
	z shared.TranslateMemorySuggestRequest) (shared.TranslateMemorySuggestResponse, error) {

	if z.PanelId == "" {
		return shared.TranslateMemorySuggestResponse{}, bladWskazaniaTlumaczenia(
			"żądanie translate.memory.suggest bez panelu")
	}
	if z.Segment == "" {
		return shared.TranslateMemorySuggestResponse{}, bladWskazaniaTlumaczenia(
			"żądanie translate.memory.suggest bez segmentu źródłowego")
	}

	panel, err := a.repozytorium.Panel(ctx, z.PanelId)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.TranslateMemorySuggestResponse{}, bladNieznanegoPanelu(z.PanelId, err)
	}
	if err != nil {
		return shared.TranslateMemorySuggestResponse{}, bladTlumaczenia(err)
	}

	wpisy, err := a.repozytorium.Podpowiedzi(ctx, panel.Jezyk, z.Segment, limitPodpowiedziPamieci)
	if err != nil {
		return shared.TranslateMemorySuggestResponse{}, bladTlumaczenia(err)
	}

	podpowiedzi := make([]string, 0, len(wpisy))
	for _, wpis := range wpisy {
		podpowiedzi = append(podpowiedzi, wpis.SegmentDocelowy)
	}
	return shared.TranslateMemorySuggestResponse{Suggestions: podpowiedzi}, nil
}
