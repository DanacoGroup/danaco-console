// Odpowiedzialność pliku: porty i rejestracja obsługiwaczy dwóch rodzin
// przekrojowych — prowenancji wywołań modelu i rozliczenia zużycia.
//
// Rodziny są dwie, magazyn jeden. Provenance Explorer pyta o pojedyncze
// wywołanie, rozliczenie liczy sumy po wymiarze, ale obie odpowiedzi powstają
// z tych samych wierszy — druga tabela z tymi samymi liczbami rozjechałaby się
// z pierwszą przy pierwszej korekcie cennika.
//
// Port jest osobny od portu modułu Diagnostics, mimo że to jego okna po niego
// sięgają. Powód jest ten sam co przy obserwatorze niepowodzeń: rdzeń nie ma
// prawa wiedzieć, że istnieje moduł Diagnostics. Ślad wywołania czyta też
// Execution Loop Window, a alerty sięgają po niego z Always On Display.
package core

import (
	"context"

	"danacoconsole/shared"
)

// Prowenancja wypełnia rodzinę `provenance.*` — odczyt śladu wywołań modelu.
type Prowenancja interface {
	WykazWywolan(ctx context.Context, z shared.ProvenanceCallListRequest) (shared.ProvenanceCallListResponse, error)
	Wywolanie(ctx context.Context, z shared.ProvenanceCallGetRequest) (shared.ProvenanceCallGetResponse, error)
	OcenWywolanie(ctx context.Context, z shared.ProvenanceCallRateRequest) (shared.ProvenanceCallRateResponse, error)
	WydajSlad(ctx context.Context, z shared.ProvenanceTraceExportRequest) (shared.ProvenanceTraceExportResponse, error)
	// PowtorzWywolanie obsługuje `provenance.call.replay`.
	PowtorzWywolanie(ctx context.Context, z shared.ProvenanceCallReplayRequest) (shared.ProvenanceCallReplayResponse, error)
}

// Zuzycie wypełnia rodzinę `usage.*` — rozliczenie żądań, tokenów i kosztu.
type Zuzycie interface {
	Podsumowanie(ctx context.Context, z shared.UsageSummaryGetRequest) (shared.UsageSummaryGetResponse, error)
	Raport(ctx context.Context, z shared.UsageReportBuildRequest) (shared.UsageReportBuildResponse, error)
}

// zarejestrujProwenancje wpina pięć komend rodziny śladu.
//
// Cztery pierwsze są odczytem. Piąta — `provenance.call.replay` — nie jest:
// powtórzenie to nowe wywołanie kanału modelu, z własnym kosztem i własnym
// wierszem śladu. Czekała na warstwę, która kanały prowadzi, i wchodzi razem
// z nią: adapter bierze ten sam rejestr kanałów, którym jedzie okno rozmowy
// (`ZKanalami`). Rejestr niewpięty daje odmowę nazywającą brak, a nie pusty
// uchwyt udający zdolność.
func zarejestrujProwenancje(r *Rejestr, p Prowenancja) {
	if r == nil || p == nil {
		return
	}
	r.Zarejestruj(shared.CommandProvenanceCallList, obsluz(p.WykazWywolan))
	r.Zarejestruj(shared.CommandProvenanceCallGet, obsluz(p.Wywolanie))
	r.Zarejestruj(shared.CommandProvenanceCallRate, obsluz(p.OcenWywolanie))
	r.Zarejestruj(shared.CommandProvenanceTraceExport, obsluz(p.WydajSlad))
	r.Zarejestruj(shared.CommandProvenanceCallReplay, obsluz(p.PowtorzWywolanie))
}

// zarejestrujZuzycie wpina dwie komendy rozliczenia.
func zarejestrujZuzycie(r *Rejestr, z Zuzycie) {
	if r == nil || z == nil {
		return
	}
	r.Zarejestruj(shared.CommandUsageSummaryGet, obsluz(z.Podsumowanie))
	r.Zarejestruj(shared.CommandUsageReportBuild, obsluz(z.Raport))
}
