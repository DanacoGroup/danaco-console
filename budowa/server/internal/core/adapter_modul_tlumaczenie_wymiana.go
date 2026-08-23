// Odpowiedzialność pliku: wymiana glosariusza modułu Translate —
// import, eksport i podpowiedzi z pamięci tłumaczeń, metody `ImportujSlownik`
// (`glossary.import`), `EksportujSlownik` (`glossary.export`), `Podpowiedzi`
// (`memory.suggest`) na `*adapterTlumaczenia`.
// Typ adaptera, konstruktor i przedrostki deklaruje
// `adapter_modul_tlumaczenie.go` — ten plik ich nie powtarza.
//
// Import glosariusza odmawia wprost. Kontrakt `TranslateGlossaryImportRequest`
// niesie wyłącznie ścieżkę pliku — żadnego pola z
// terminami do zapisania. Rdzeń nie czyta plików z dysku Operatora (nie ma do
// niego dostępu), więc jedynym uczciwym zachowaniem jest odmowa wprost, wzorem
// `RozpoznajJezyk` (ten sam brak, ta sama reakcja) — nie udajemy wczytania
// pliku, którego nie umiemy otworzyć.
//
// Eksport glosariusza zapisuje ślad zlecenia, nie plik. Rdzeń nie ma magazynu
// blobów (ten sam brak co w Library i Research). `TranslateGlossaryExportResponse`
// niesie wyłącznie `ExportedCount` — inaczej niż `research.report.export`,
// kontrakt nie daje ani identyfikatora pliku, ani rozmiaru, więc nie dorabiamy
// pól, których nikt uczciwie nie wypełni.
// `ExportedCount` liczy terminy zastane w słowniku w chwili zlecenia — to
// jedyna liczba, jaką rdzeń naprawdę zna, skoro pliku nie zapisuje.
//
// Podpowiedzi z pamięci oddają dopasowanie podciągu, nie podobieństwa. Warstwa danych
// (`dane/slownik_pamiec.go`): SQLite bez rozszerzenia nie ma
// miary podobieństwa napisów, więc `Podpowiedzi` repozytorium dopasowuje przez
// `LIKE '%fraza%'` — dopasowanie podciągu segmentu źródłowego, nie dopasowanie
// znaczeniowe ani odległość edycyjną. Ten adapter oddaje dokładnie to, co
// repozytorium naprawdę znalazło; nie sortuje ani nie filtruje wyniku tak, by
// wyglądał na dopasowanie przybliżone, którego nie ma.
package core

import (
	"context"
	"errors"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// limitPodpowiedziPamieci ogranicza liczbę podpowiedzi zwracanych z jednego
// zapytania — kontrakt nie niesie własnego limitu, więc adapter narzuca stałą,
// rozsądną dla panelu edycji (ta sama wartość co domyślna w warstwie danych).
const limitPodpowiedziPamieci = 5

// Kod odmowy: `channel_unavailable`. Kontrakt zna dziewięć kodów błędu
// i nie ma wśród nich takiego, który mówiłby „brak realizacji".
// `channel_unavailable` jest tym z nich, który mówi prawdę: brakuje wykonawcy
// po stronie rdzenia, a nie żądanie jest złe. Tego samego kodu używa moduł
// Studio przy `contextual.op` bez kanału modelu.
//
// ImportujSlownik obsługuje `translate.glossary.import`. Odmawia wprost —
// patrz nagłówek pliku: kontrakt niesie tylko ścieżkę pliku, a rdzeń nie ma
// dostępu do dysku Operatora, więc nie ma czym wczytać terminów. Odmowa nosi
// kod „nieobsłużone", nie usterkę wewnętrzną — to zameldowany brak kontraktu,
// nie awaria.
func (a *adapterTlumaczenia) ImportujSlownik(ctx context.Context,
	z shared.TranslateGlossaryImportRequest) (shared.TranslateGlossaryImportResponse, error) {

	if z.Path == "" {
		return shared.TranslateGlossaryImportResponse{}, bladWskazaniaTlumaczenia(
			"żądanie translate.glossary.import bez ścieżki pliku")
	}

	// Zapisujemy ślad próby (żądanej ścieżki, zero zaimportowanych) — Operator
	// ma dowód, że zlecenie dotarło do rdzenia, nawet gdy rdzeń go odmówił.
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

// EksportujSlownik obsługuje `translate.glossary.export`. Zapisuje ślad
// zlecenia — patrz nagłówek pliku: rdzeń nie ma magazynu blobów, więc żaden
// plik nie powstaje. `ExportedCount` niesie liczbę terminów zastanych w
// słowniku w chwili zlecenia — jedyna liczba, jaką rdzeń naprawdę zna.
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

// Podpowiedzi obsługuje `translate.memory.suggest`. Ustala język docelowy z
// panelu wskazanego w żądaniu (pamięć tłumaczeń jest zawężona do języka), po
// czym pyta repozytorium o dopasowanie podciągu segmentu źródłowego — patrz
// nagłówek pliku o ograniczeniu do `LIKE`, nie dopasowania przybliżonego.
// Panel nieznany wraca odmową `not_found` (`bladNieznanegoPanelu`), nie cichą
// pustą listą.
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
