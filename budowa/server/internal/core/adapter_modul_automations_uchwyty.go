// Wpięcie sześciu komend obszaru `automation.*` — modułu Automations wraz z
// jego pięcioma oknami operacyjnymi (Workflow Builder, Scheduler, Queue
// Manager, Orchestrator, Execution Monitor).
//
// Jedno działanie na kolejce rozgłasza dwa zdarzenia i nie jest to powtórzenie.
// `queue.changed` niesie kolejkę i dotyczy Queue Managera oraz Mission Control;
// `automation.execution.status` niesie przebieg automatyki i zasila Execution
// Monitor. Są to dwa różne byty tej samej czynności, więc rozgłaszają się
// osobno. Trzeci nośnik — telemetria postępu `progress.changed` — wychodzi
// z adaptera kolejek i tu się go nie powtarza.
//
// Osobnych komend odczytu harmonogramu, wykazu kolejek i układu zależności
// `shared/contract.json` nie zna. Ich pracę wykonują komendy istniejące: układ
// zależności prowadzi `automation.orchestrator.define`, a stan przebiegów
// `automation.execution.subscribe`.
package core

import (
	"context"

	"danacoconsole/shared"
)

// Automatyki jest portem modułu Automations.
type Automatyki interface {
	Zapisz(ctx context.Context, z shared.AutomationWorkflowSaveRequest) (shared.AutomationWorkflowSaveResponse, error)
	Wykaz(ctx context.Context, z shared.AutomationWorkflowListRequest) (shared.AutomationWorkflowListResponse, error)
	UstawHarmonogram(ctx context.Context, z shared.AutomationScheduleSetRequest) (shared.AutomationScheduleSetResponse, error)
	DzialanieKolejki(ctx context.Context, z shared.AutomationQueueActionRequest) (shared.AutomationQueueActionResponse, error)
	Zaleznosci(ctx context.Context, z shared.AutomationOrchestratorDefineRequest) (shared.AutomationOrchestratorDefineResponse, error)
	Przebiegi(ctx context.Context, z shared.AutomationExecutionSubscribeRequest) (shared.AutomationExecutionSubscribeResponse, error)
	// PrzebiegKolejki oddaje przebieg wykonywany przez kolejkę. Służy
	// rozgłoszeniu `automation.execution.status` po działaniu, którego wynik
	// niesie kolejkę, a nie przebieg.
	PrzebiegKolejki(ctx context.Context, idKolejki string) (shared.AutomationExecution, bool)
}

// zarejestrujAutomatyki wpina sześć komend modułu Automations.
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
				e.powiazanieAutomatyki(shared.ChangeKindUpdated, z.WorkflowId)
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandAutomationOrchestratorDefine,
		obsluz(func(ctx context.Context, z shared.AutomationOrchestratorDefineRequest) (shared.AutomationOrchestratorDefineResponse, error) {
			odpowiedz, err := m.Zaleznosci(ctx, z)
			// Żądanie bez `dependencies` jest samym sprawdzeniem układu
			// zastanego („Waliduj graf”) — niczego nie zmienia, więc niczego
			// nie rozgłasza. Zdarzenie po odczycie byłoby szumem.
			if err == nil && z.Dependencies != nil {
				e.ukladOrkiestracji(shared.ChangeKindUpdated, "")
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandAutomationExecutionSubscribe, obsluz(m.Przebiegi))

	r.Zarejestruj(shared.CommandAutomationQueueAction,
		obsluz(func(ctx context.Context, z shared.AutomationQueueActionRequest) (shared.AutomationQueueActionResponse, error) {
			odpowiedz, err := m.DzialanieKolejki(ctx, z)
			if err == nil {
				e.kolejka(shared.ChangeKindUpdated, odpowiedz.Queue)
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
	e.przebiegAutomatyki(przebieg, nazwaEtapuPrzebiegu(dzialanie))
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
func (e *emiter) przebiegAutomatyki(przebieg shared.AutomationExecution, nazwaEtapu string) {
	tresc := shared.AutomationExecutionStatusEvent{Execution: przebieg}
	if nazwaEtapu != "" {
		etap := nazwaEtapu
		tresc.StepLabel = &etap
	}
	if przebieg.ErrorMessage != nil {
		wiersz := *przebieg.ErrorMessage
		tresc.LogLine = &wiersz
	}
	e.wyslij(shared.EventAutomationExecutionStatus, "", tresc)
}

// powiazanieAutomatyki rozgłasza `automation.link.changed` — zmianę powiązania
// automatyki z bytem wyzwalającym.
//
// Bytem wyzwalającym jest harmonogram: `automation.schedule.set` zapisuje
// cykliczność i wyzwalacze, czyli jedyne w rdzeniu wskazanie, co ma automatykę
// uruchomić. Ładunek kontraktu niesie samo `automationId`, więc mówi
// „powiązanie tej automatyki jest inne niż było", a nie jakie. Bliski krewny,
// `queue.link`, wiąże kolejkę z automatyką i rozgłasza `queue.changed`, bo
// bytem zmienianym jest tam kolejka, nie automatyka.
//
// Rodzaj zmiany jest zawsze `updated`. Automatyka powiązanie ma zawsze,
// choćby puste (harmonogram nieczynny, zero wyzwalaczy) — `created`
// i `deleted` opisywałyby byt o własnym cyklu życia, którego tu nie ma.
func (e *emiter) powiazanieAutomatyki(zmiana shared.ChangeKind, idAutomatyki string) {
	if idAutomatyki == "" {
		return
	}
	e.wyslij(shared.EventAutomationLinkChanged, "",
		shared.AutomationLinkChangedEvent{Change: zmiana, AutomationId: idAutomatyki})
}
