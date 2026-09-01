// Plik wpina cztery komendy obszaru `diagnostics.*`, obsługujące moduł Diagnostics wraz z jego czterema oknami operacyjnymi opartymi na jednym zdarzeniu zmiany analizy.
package core

import (
	"context"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// Diagnostyka jest portem modułu Diagnostics, obsługującym dziennik, wykaz błędów, uruchomienie analizy i wykaz rekomendacji tego modułu.
type Diagnostyka interface {
	PrzeszukajDziennik(ctx context.Context, z shared.DiagnosticsLogQueryRequest) (shared.DiagnosticsLogQueryResponse, error)
	WykazBledow(ctx context.Context, z shared.DiagnosticsErrorListRequest) (shared.DiagnosticsErrorListResponse, error)
	UruchomAnalize(ctx context.Context, z shared.DiagnosticsAnalyzeRunRequest) (shared.DiagnosticsAnalyzeRunResponse, error)
	WykazRekomendacji(ctx context.Context, z shared.DiagnosticsRecommendationListRequest) (shared.DiagnosticsRecommendationListResponse, error)
	// PodepnijRozgloszenie oddaje adapterowi drogę do zdarzenia zmiany analizy.
	PodepnijRozgloszenie(rozglos func(context.Context, shared.ChangeKind, shared.DiagnosticAnalysis))
	ObserwatorNiepowodzen
}

// ObserwatorNiepowodzen przyjmuje odmowę wykonania komendy, żeby stała się faktem widocznym w Errors Panel, portem osobnym od portu modułu Diagnostics.
type ObserwatorNiepowodzen interface {
	ZapiszNiepowodzenie(ctx context.Context, n NiepowodzenieKomendy)
}

// NiepowodzenieKomendy opisuje jedną odmowę: komendę, zasięg i błąd kontraktu wraz z jego kodem ze słownika ErrorCode.
type NiepowodzenieKomendy struct {
	Komenda shared.MessageType
	IdSesji string
	IdOkna  string
	Blad    shared.ErrorInfo
}

// zarejestrujDiagnostyke wpina cztery komendy modułu Diagnostics i podpina rozgłoszenie zmiany analizy przekazane przez adapter.
func zarejestrujDiagnostyke(r *Rejestr, d Diagnostyka, e *emiter) {
	if r == nil || d == nil {
		return
	}
	d.PodepnijRozgloszenie(e.analizaDiagnostyczna)

	// Rozgłoszenia po `diagnostics.analyze.run` nie ma tutaj celowo: nadaje je adapter.
	r.Zarejestruj(shared.CommandDiagnosticsLogQuery, obsluz(d.PrzeszukajDziennik))
	r.Zarejestruj(shared.CommandDiagnosticsErrorList, obsluz(d.WykazBledow))
	r.Zarejestruj(shared.CommandDiagnosticsAnalyzeRun, obsluz(d.UruchomAnalize))
	r.Zarejestruj(shared.CommandDiagnosticsRecommendationList, obsluz(d.WykazRekomendacji))
}

// analizaDiagnostyczna rozgłasza zmianę analizy. Sesja komunikatu jest pusta: analiza jest bytem rdzenia, nie sesji, a Diagnostics Center bywa otwarte w każdym oknie modułu.
func (e *emiter) analizaDiagnostyczna(ctx context.Context, zmiana shared.ChangeKind, a shared.DiagnosticAnalysis) {
	e.wyslijDoKonta(ctx, shared.EventDiagnosticsAnalysisChanged, "",
		shared.DiagnosticsAnalysisChangedEvent{Change: zmiana, Analysis: a})
}

// odnotujNiepowodzenie oddaje odmowę obserwatorowi. Brak obserwatora nie zmienia zachowania rdzenia, a odpowiedź udana nie jest niepowodzeniem.
func (r *Rdzen) odnotujNiepowodzenie(ctx context.Context, z protocol.Request, blad *shared.ErrorInfo) {
	if r == nil || r.niepowodzenia == nil || blad == nil {
		return
	}
	r.niepowodzenia.ZapiszNiepowodzenie(ctx, NiepowodzenieKomendy{
		Komenda: z.Komenda, IdSesji: z.Zasieg.Sesja, IdOkna: z.Zasieg.Okno, Blad: *blad,
	})
}

// odmowaNieznanej buduje odpowiedź na komendę bez uchwytu i odnotowuje ją jako odmowę, ze źródłem równym typowi żądanemu, żeby Errors Panel grupował wystąpienia po nazwie odrzuconej komendy.
func (r *Rdzen) odmowaNieznanej(ctx context.Context, z protocol.Request) protocol.Koperta {
	koperta := odpowiedzNieznanej(z)
	odmowa := z
	if z.TypZadany != "" {
		odmowa.Komenda = z.TypZadany
	}
	r.odnotujNiepowodzenie(ctx, odmowa, koperta.Error)
	return koperta
}
