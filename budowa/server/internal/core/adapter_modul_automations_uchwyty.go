// Wpięcie sześciu komend obszaru automation.* modułu Automations wraz z pięcioma oknami operacyjnymi; jedno działanie kolejki rozgłasza dwa różne zdarzenia osobno.
package core

import (
	"context"

	"danacoconsole/shared"
)

// Automatyki jest portem modułu Automations, wystawiającym sześć komend obszaru automation.*, wraz z jego oknami.
type Automatyki interface {
	Zapisz(ctx context.Context, z shared.AutomationWorkflowSaveRequest) (shared.AutomationWorkflowSaveResponse, error)
	Wykaz(ctx context.Context, z shared.AutomationWorkflowListRequest) (shared.AutomationWorkflowListResponse, error)
	UstawHarmonogram(ctx context.Context, z shared.AutomationScheduleSetRequest) (shared.AutomationScheduleSetResponse, error)
	DzialanieKolejki(ctx context.Context, z shared.AutomationQueueActionRequest) (shared.AutomationQueueActionResponse, error)
	Zaleznosci(ctx context.Context, z shared.AutomationOrchestratorDefineRequest) (shared.AutomationOrchestratorDefineResponse, error)
	Przebiegi(ctx context.Context, z shared.AutomationExecutionSubscribeRequest) (shared.AutomationExecutionSubscribeResponse, error)
	// PrzebiegKolejki oddaje przebieg wykonywany przez kolejkę, do rozgłoszenia stanu przebiegu.
	PrzebiegKolejki(ctx context.Context, idKolejki string) (shared.AutomationExecution, bool)
}

// zarejestrujAutomatyki wpina sześć komend modułu Automations do rejestru obsługiwaczy komend rdzenia.
func zarejestrujAutomatyki(r *Rejestr, m Automatyki, e *emiter) {
	if r == nil || m == nil {
		return
	}

	r.Zarejestruj(shared.CommandAutomationWorkflowSave, obsluz(m.Zapisz))
	r.Zarejestruj(shared.CommandAutomationWorkflowList, obsluz(m.Wykaz))
	r.Zarejestruj(shared.CommandAutomationScheduleSet,
		obsluz(func(ctx context.Context, z shared.AutomationScheduleSetRequest) (shared.AutomationScheduleSetResponse, error) {
			odpowiedz, err := m.UstawHarmonogram(ctx, z)
			if err == nil {
				e.powiazanieAutomatyki(ctx, shared.ChangeKindUpdated, z.WorkflowId)
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandAutomationOrchestratorDefine,
		obsluz(func(ctx context.Context, z shared.AutomationOrchestratorDefineRequest) (shared.AutomationOrchestratorDefineResponse, error) {
			odpowiedz, err := m.Zaleznosci(ctx, z)
			// Żądanie bez dependencies sprawdza układ zastany, niczego nie zmienia ani nie rozgłasza.
			if err == nil && z.Dependencies != nil {
				e.ukladOrkiestracji(ctx, shared.ChangeKindUpdated, "")
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandAutomationExecutionSubscribe, obsluz(m.Przebiegi))

	r.Zarejestruj(shared.CommandAutomationQueueAction,
		obsluz(func(ctx context.Context, z shared.AutomationQueueActionRequest) (shared.AutomationQueueActionResponse, error) {
			odpowiedz, err := m.DzialanieKolejki(ctx, z)
			if err == nil {
				e.kolejka(ctx, shared.ChangeKindUpdated, odpowiedz.Queue)
				rozglosPrzebieg(ctx, m, e, odpowiedz.Queue.Id, z.Action)
			}
			return odpowiedz, err
		}))
}

// rozglosPrzebieg dobiera przebieg kolejki po działaniu i rozgłasza jego stan.
// Nieudany dobór kończy wyłącznie rozgłoszenie — komenda już się powiodła.
// Kolejka spoza automatyki przebiegu nie ma i nic nie rozgłasza.
func rozglosPrzebieg(ctx context.Context, m Automatyki, e *emiter,
	idKolejki string, dzialanie shared.QueueAction) {

	przebieg, jest := m.PrzebiegKolejki(ctx, idKolejki)
	if !jest {
		return
	}
	e.przebiegAutomatyki(ctx, przebieg, nazwaEtapuPrzebiegu(dzialanie))
}

// nazwaEtapuPrzebiegu opisuje Operatorowi, co właśnie zrobiono z przebiegiem.
// Powtórzenie kroku nazywa się biegiem naprawczym — tak samo jak w telemetrii
// kolejki, żeby oba okna mówiły o tym samym jednym językiem.
func nazwaEtapuPrzebiegu(dzialanie shared.QueueAction) string {
	switch dzialanie {
	case shared.QueueActionRetry:
		return etapBiegNaprawczy
	case shared.QueueActionStop:
		return etapZatrzymanie
	default:
		return etapDzialanieNaKolejce
	}
}

// przebiegAutomatyki rozgłasza stan przebiegu automatyki. Automatyka jest
// komponentem własnym, nie bytem karty sesji, więc zdarzenie idzie bez jej
// wskazania — Execution Monitor otwiera się ze strony głównej.
func (e *emiter) przebiegAutomatyki(ctx context.Context, przebieg shared.AutomationExecution, nazwaEtapu string) {
	tresc := shared.AutomationExecutionStatusEvent{Execution: przebieg}
	if nazwaEtapu != "" {
		etap := nazwaEtapu
		tresc.StepLabel = &etap
	}
	if przebieg.ErrorMessage != nil {
		wiersz := *przebieg.ErrorMessage
		tresc.LogLine = &wiersz
	}
	e.wyslijDoKonta(ctx, shared.EventAutomationExecutionStatus, "", tresc)
}

// powiazanieAutomatyki rozgłasza automation.link.changed po zmianie powiązania automatyki z bytem wyzwalającym, zawsze jako rodzaj updated.
func (e *emiter) powiazanieAutomatyki(ctx context.Context, zmiana shared.ChangeKind, idAutomatyki string) {
	if idAutomatyki == "" {
		return
	}
	e.wyslijDoKonta(ctx, shared.EventAutomationLinkChanged, "",
		shared.AutomationLinkChangedEvent{Change: zmiana, AutomationId: idAutomatyki})
}
