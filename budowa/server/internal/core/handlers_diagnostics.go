// Odpowiedzialność pliku: wpięcie czterech komend obszaru `diagnostics.*` —
// modułu Diagnostics wraz z jego czterema oknami operacyjnymi (Diagnostics
// Center, Logs Viewer, Errors Panel, Recommendations Panel).
//
// Cały moduł ma jedno zdarzenie. Kontrakt daje modułowi wyłącznie
// `diagnostics.analysis.changed`, więc Recommendations Panel odświeża się
// z jednej subskrypcji po każdym uruchomieniu analizy w Diagnostics Center.
//
// Dziennik i błędy zdarzenia nie mają: kontrakt nie niesie ani zdarzenia
// dopisania wpisu, ani zgłoszenia błędu. Logs Viewer odczytuje więc
// `diagnostics.log.query` w odstępie zadanym przez Operatora, a rdzeń niczego
// nie rozgłasza na wyrost.
package core

import (
	"context"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// Diagnostyka jest portem modułu Diagnostics.
type Diagnostyka interface {
	PrzeszukajDziennik(ctx context.Context, z shared.DiagnosticsLogQueryRequest) (shared.DiagnosticsLogQueryResponse, error)
	WykazBledow(ctx context.Context, z shared.DiagnosticsErrorListRequest) (shared.DiagnosticsErrorListResponse, error)
	UruchomAnalize(ctx context.Context, z shared.DiagnosticsAnalyzeRunRequest) (shared.DiagnosticsAnalyzeRunResponse, error)
	WykazRekomendacji(ctx context.Context, z shared.DiagnosticsRecommendationListRequest) (shared.DiagnosticsRecommendationListResponse, error)
	// PodepnijRozgloszenie oddaje adapterowi drogę do zdarzenia zmiany analizy.
	PodepnijRozgloszenie(rozglos func(shared.ChangeKind, shared.DiagnosticAnalysis))
	ObserwatorNiepowodzen
}

// ObserwatorNiepowodzen przyjmuje odmowę wykonania komendy, żeby stała się
// faktem widocznym w Errors Panel.
//
// Port jest osobny od portu modułu ze względu na kierunek zależności: rdzeń nie
// ma prawa wiedzieć, że istnieje moduł Diagnostics — wie wyłącznie, że ktoś może
// chcieć usłyszeć o niepowodzeniu. Bez tego rozdzielenia dyspozytor komend
// zależałby od jednego z modułów.
type ObserwatorNiepowodzen interface {
	ZapiszNiepowodzenie(ctx context.Context, n NiepowodzenieKomendy)
}

// NiepowodzenieKomendy opisuje jedną odmowę: komendę, zasięg i błąd kontraktu
// wraz z jego kodem ze słownika ErrorCode.
type NiepowodzenieKomendy struct {
	Komenda shared.MessageType
	IdSesji string
	IdOkna  string
	Blad    shared.ErrorInfo
}

// zarejestrujDiagnostyke wpina cztery komendy modułu Diagnostics.
func zarejestrujDiagnostyke(r *Rejestr, d Diagnostyka, e *emiter) {
	if r == nil || d == nil {
		return
	}
	d.PodepnijRozgloszenie(e.analizaDiagnostyczna)

	// Rozgłoszenia po `diagnostics.analyze.run` nie ma tutaj z zamysłem:
	// zdarzenie nadaje adapter, bo tylko on wie, czy migawka rzeczywiście
	// powstała. Rozgłoszenie z obsługiwacza powiadamiałoby także o analizie,
	// której zapis się nie powiódł.
	r.Zarejestruj(shared.CommandDiagnosticsLogQuery, obsluz(d.PrzeszukajDziennik))
	r.Zarejestruj(shared.CommandDiagnosticsErrorList, obsluz(d.WykazBledow))
	r.Zarejestruj(shared.CommandDiagnosticsAnalyzeRun, obsluz(d.UruchomAnalize))
	r.Zarejestruj(shared.CommandDiagnosticsRecommendationList, obsluz(d.WykazRekomendacji))
}

// analizaDiagnostyczna rozgłasza zmianę analizy. Sesja komunikatu jest pusta:
// analiza jest bytem rdzenia, nie sesji, a Diagnostics Center bywa otwarte
// w każdym oknie modułu.
func (e *emiter) analizaDiagnostyczna(zmiana shared.ChangeKind, a shared.DiagnosticAnalysis) {
	e.wyslij(shared.EventDiagnosticsAnalysisChanged, "",
		shared.DiagnosticsAnalysisChangedEvent{Change: zmiana, Analysis: a})
}

// odnotujNiepowodzenie oddaje odmowę obserwatorowi. Brak obserwatora nie zmienia
// zachowania rdzenia, a odpowiedź udana nie jest niepowodzeniem.
func (r *Rdzen) odnotujNiepowodzenie(ctx context.Context, z protocol.Request, blad *shared.ErrorInfo) {
	if r == nil || r.niepowodzenia == nil || blad == nil {
		return
	}
	r.niepowodzenia.ZapiszNiepowodzenie(ctx, NiepowodzenieKomendy{
		Komenda: z.Komenda, IdSesji: z.Zasieg.Sesja, IdOkna: z.Zasieg.Okno, Blad: *blad,
	})
}

// odmowaNieznanej buduje odpowiedź na komendę bez uchwytu i odnotowuje ją jako
// odmowę — tak samo jak odmowę merytoryczną obsługiwacza. Źródłem odmowy jest
// typ żądany, nie `<obszar>.unknown`: Errors Panel stawia `source` w tytule
// pozycji jako nazwę odrzuconej komendy i po niej grupuje wystąpienia. Nazwa
// zdarzenia obszaru zlepiłaby wszystkie nieobsłużone komendy obszaru w jeden
// nierozróżnialny wiersz.
func (r *Rdzen) odmowaNieznanej(ctx context.Context, z protocol.Request) protocol.Koperta {
	koperta := odpowiedzNieznanej(z)
	odmowa := z
	if z.TypZadany != "" {
		odmowa.Komenda = z.TypZadany
	}
	r.odnotujNiepowodzenie(ctx, odmowa, koperta.Error)
	return koperta
}
