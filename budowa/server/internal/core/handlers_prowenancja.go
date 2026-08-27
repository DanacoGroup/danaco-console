// Plik wpina porty i rejestrację obsługiwaczy dwóch rodzin przekrojowych: prowenancji wywołań
// modelu i rozliczenia zużycia, opartych na wspólnym magazynie wierszy śladu.
package core

import (
	"context"

	"danacoconsole/shared"
)

// Prowenancja wypełnia rodzinę provenance.* obsługującą odczyt śladu wywołań modelu z magazynu wierszy.
type Prowenancja interface {
	WykazWywolan(ctx context.Context, z shared.ProvenanceCallListRequest) (shared.ProvenanceCallListResponse, error)
	Wywolanie(ctx context.Context, z shared.ProvenanceCallGetRequest) (shared.ProvenanceCallGetResponse, error)
	OcenWywolanie(ctx context.Context, z shared.ProvenanceCallRateRequest) (shared.ProvenanceCallRateResponse, error)
	WydajSlad(ctx context.Context, z shared.ProvenanceTraceExportRequest) (shared.ProvenanceTraceExportResponse, error)
	// PowtorzWywolanie obsługuje `provenance.call.replay`.
	PowtorzWywolanie(ctx context.Context, z shared.ProvenanceCallReplayRequest) (shared.ProvenanceCallReplayResponse, error)
}

// Zuzycie wypełnia rodzinę usage.* obsługującą rozliczenie żądań, tokenów i kosztu zużycia modeli w oknie rozliczeń.
type Zuzycie interface {
	Podsumowanie(ctx context.Context, z shared.UsageSummaryGetRequest) (shared.UsageSummaryGetResponse, error)
	Raport(ctx context.Context, z shared.UsageReportBuildRequest) (shared.UsageReportBuildResponse, error)
}

// zarejestrujProwenancje wpina pięć komend rodziny śladu, z których cztery są odczytem, a piąta
// powtarza wywołanie kanału modelu z własnym kosztem i wierszem śladu.
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

// zarejestrujZuzycie wpina dwie komendy rozliczenia żądań, tokenów i kosztu zużycia modelu w tym oknie.
func zarejestrujZuzycie(r *Rejestr, z Zuzycie) {
	if r == nil || z == nil {
		return
	}
	r.Zarejestruj(shared.CommandUsageSummaryGet, obsluz(z.Podsumowanie))
	r.Zarejestruj(shared.CommandUsageReportBuild, obsluz(z.Raport))
}
