// Wpięcie dwudziestu sześciu komend obszaru `automation.*` dopełniających
// sześć komend z `adapter_modul_automations_uchwyty.go`.
//
// Osobny plik, bo osobna jest odpowiedzialność: tamten wpina rdzeń modułu
// (definicja, harmonogram, kolejka, zależności, przebiegi), ten — panele
// i szuflady okien operacyjnych, którymi Operator sięga po wersje, szablony,
// logi, ładunki, alarmy, sekrety i audyt. Port jest jeden, bo moduł jest jeden.
//
// Rozgłasza się tylko to, co zmienia byt widziany przez inne okno. Zapis
// definicji, publikacja, udostępnienie, etykiety i przywrócenie wersji zmieniają
// automatykę, więc idą zdarzeniem `automation.link.changed` — tym samym, którym
// idzie zmiana harmonogramu. Odczyty (wykaz wersji, porównanie, log, kroki,
// ładunek, audyt, wykaz szablonów, wykaz reguł, wykaz sekretów) nie rozgłaszają
// niczego: zdarzenie po odczycie byłoby szumem, na który okna reagowałyby
// odświeżeniem bez powodu.
//
// Wznowienie i odtworzenie przebiegu rozgłaszają stan przebiegu
// (`automation.execution.status`), bo obie zmieniają to, co pokazuje Execution
// Monitor — i robią to przez ten sam nośnik, którym idzie działanie na kolejce.
package core

import (
	"context"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// AutomatykiDobudowa jest portem dwudziestu sześciu komend dopełniających
// moduł Automations. Rozszerza port `Automatyki`, nie stoi obok niego:
// moduł ma w rdzeniu jednego właściciela i jedno repozytorium.
type AutomatykiDobudowa interface {
	Automatyki

	// Workflow Builder: przebieg próbny, historia wersji, etykiety, publikacja,
	// udostępnienie.
	Symulacja(ctx context.Context, z shared.AutomationWorkflowSimulateRequest) (shared.AutomationWorkflowSimulateResponse, error)
	WykazWersji(ctx context.Context, z shared.AutomationWorkflowVersionListRequest) (shared.AutomationWorkflowVersionListResponse, error)
	PrzywrocWersje(ctx context.Context, z shared.AutomationWorkflowVersionRestoreRequest) (shared.AutomationWorkflowVersionRestoreResponse, error)
	PorownajWersje(ctx context.Context, z shared.AutomationWorkflowVersionDiffRequest) (shared.AutomationWorkflowVersionDiffResponse, error)
	UstawEtykiety(ctx context.Context, z shared.AutomationWorkflowTagSetRequest) (shared.AutomationWorkflowTagSetResponse, error)
	Opublikuj(ctx context.Context, z shared.AutomationWorkflowPublishRequest) (shared.AutomationWorkflowPublishResponse, error)
	Udostepnij(ctx context.Context, z shared.AutomationWorkflowShareRequest) (shared.AutomationWorkflowShareResponse, error)

	// Panel zmiennych i kanwa.
	UstawZmienne(ctx context.Context, z shared.AutomationWorkflowVariablesSetRequest) (shared.AutomationWorkflowVariablesSetResponse, error)
	UstawNotatkeKroku(ctx context.Context, z shared.AutomationStepNoteSetRequest) (shared.AutomationStepNoteSetResponse, error)
	UstawUkladKanwy(ctx context.Context, z shared.AutomationStepLayoutSetRequest) (shared.AutomationStepLayoutSetResponse, error)

	// Biblioteka szablonów przepływów.
	ZapiszSzablon(ctx context.Context, z shared.AutomationTemplateSaveRequest) (shared.AutomationTemplateSaveResponse, error)
	WykazSzablonow(ctx context.Context, z shared.AutomationTemplateListRequest) (shared.AutomationTemplateListResponse, error)
	ZastosujSzablon(ctx context.Context, z shared.AutomationTemplateApplyRequest) (shared.AutomationTemplateApplyResponse, error)

	// Execution Monitor: log, kroki, punkty wznowienia, ładunki, odtworzenie.
	DziennikPrzebiegu(ctx context.Context, z shared.AutomationExecutionLogRequest) (shared.AutomationExecutionLogResponse, error)
	KrokiPrzebiegu(ctx context.Context, z shared.AutomationExecutionStepsRequest) (shared.AutomationExecutionStepsResponse, error)
	WykazPunktowWznowienia(ctx context.Context, z shared.AutomationExecutionCheckpointListRequest) (shared.AutomationExecutionCheckpointListResponse, error)
	WznowPrzebieg(ctx context.Context, z shared.AutomationExecutionResumeRequest) (shared.AutomationExecutionResumeResponse, error)
	PodgladLadunku(ctx context.Context, z shared.AutomationExecutionPayloadGetRequest) (shared.AutomationExecutionPayloadGetResponse, error)
	OdtworzPrzebieg(ctx context.Context, z shared.AutomationExecutionReplayRequest) (shared.AutomationExecutionReplayResponse, error)

	// Alarmowanie, budżety, skarbiec i audyt.
	UstawReguleAlarmowania(ctx context.Context, z shared.AutomationAlertRuleSetRequest) (shared.AutomationAlertRuleSetResponse, error)
	WykazRegulAlarmowania(ctx context.Context, z shared.AutomationAlertRuleListRequest) (shared.AutomationAlertRuleListResponse, error)
	UstawBudzetyPrzebiegu(ctx context.Context, z shared.AutomationExecutionBudgetSetRequest) (shared.AutomationExecutionBudgetSetResponse, error)
	ZapiszPoswiadczenie(ctx context.Context, z shared.AutomationSecretSetRequest) (shared.AutomationSecretSetResponse, error)
	WykazPoswiadczen(ctx context.Context, z shared.AutomationSecretListRequest) (shared.AutomationSecretListResponse, error)
	UsunPoswiadczenie(ctx context.Context, z shared.AutomationSecretRemoveRequest) (shared.AutomationSecretRemoveResponse, error)
	OdczytajAudyt(ctx context.Context, z shared.AutomationAuditListRequest) (shared.AutomationAuditListResponse, error)
}

// zarejestrujDobudoweAutomatyk wpina dwadzieścia sześć komend dopełniających.
//
// Port modułu bez tych czynności to co innego niż brak portu: moduł jest,
// a rdzeń nie umie wykonać części jego pracy. Wtedy komendy i tak zostają
// wpięte i odmawiają wprost — odpowiedź „nieznana komenda" wskazywałaby na brak
// modułu, a nie na usterkę montażu.
func zarejestrujDobudoweAutomatyk(r *Rejestr, m Automatyki, e *emiter) {
	if r == nil || m == nil {
		return
	}
	dobudowa, ok := m.(AutomatykiDobudowa)
	if !ok {
		zarejestrujOdmoweDobudowyAutomatyk(r)
		return
	}

	r.Zarejestruj(shared.CommandAutomationWorkflowSimulate, obsluz(dobudowa.Symulacja))
	r.Zarejestruj(shared.CommandAutomationWorkflowVersionList, obsluz(dobudowa.WykazWersji))
	r.Zarejestruj(shared.CommandAutomationWorkflowVersionDiff, obsluz(dobudowa.PorownajWersje))
	r.Zarejestruj(shared.CommandAutomationTemplateList, obsluz(dobudowa.WykazSzablonow))
	r.Zarejestruj(shared.CommandAutomationExecutionLog, obsluz(dobudowa.DziennikPrzebiegu))
	r.Zarejestruj(shared.CommandAutomationExecutionSteps, obsluz(dobudowa.KrokiPrzebiegu))
	r.Zarejestruj(shared.CommandAutomationExecutionCheckpointList,
		obsluz(dobudowa.WykazPunktowWznowienia))
	r.Zarejestruj(shared.CommandAutomationExecutionPayloadGet, obsluz(dobudowa.PodgladLadunku))
	r.Zarejestruj(shared.CommandAutomationAlertRuleList, obsluz(dobudowa.WykazRegulAlarmowania))
	r.Zarejestruj(shared.CommandAutomationSecretList, obsluz(dobudowa.WykazPoswiadczen))
	r.Zarejestruj(shared.CommandAutomationAuditList, obsluz(dobudowa.OdczytajAudyt))
	r.Zarejestruj(shared.CommandAutomationStepLayoutSet, obsluz(dobudowa.UstawUkladKanwy))
	r.Zarejestruj(shared.CommandAutomationStepNoteSet, obsluz(dobudowa.UstawNotatkeKroku))
	r.Zarejestruj(shared.CommandAutomationWorkflowVariablesSet, obsluz(dobudowa.UstawZmienne))
	r.Zarejestruj(shared.CommandAutomationTemplateSave, obsluz(dobudowa.ZapiszSzablon))
	r.Zarejestruj(shared.CommandAutomationAlertRuleSet, obsluz(dobudowa.UstawReguleAlarmowania))
	r.Zarejestruj(shared.CommandAutomationSecretSet, obsluz(dobudowa.ZapiszPoswiadczenie))
	r.Zarejestruj(shared.CommandAutomationSecretRemove, obsluz(dobudowa.UsunPoswiadczenie))

	zarejestrujZmianyAutomatyki(r, dobudowa, e)
	zarejestrujZmianyPrzebiegu(r, dobudowa, e)
}

// zarejestrujZmianyAutomatyki wpina pięć komend zmieniających automatykę.
// Każda rozgłasza `automation.link.changed`, bo każda czyni wykaz automatyk
// okna nieaktualnym — i robi to tym samym nośnikiem, co zmiana harmonogramu.
func zarejestrujZmianyAutomatyki(r *Rejestr, m AutomatykiDobudowa, e *emiter) {
	r.Zarejestruj(shared.CommandAutomationWorkflowVersionRestore,
		obsluz(func(ctx context.Context, z shared.AutomationWorkflowVersionRestoreRequest) (shared.AutomationWorkflowVersionRestoreResponse, error) {
			odpowiedz, err := m.PrzywrocWersje(ctx, z)
			if err == nil {
				e.powiazanieAutomatyki(shared.ChangeKindUpdated, z.WorkflowId)
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandAutomationWorkflowTagSet,
		obsluz(func(ctx context.Context, z shared.AutomationWorkflowTagSetRequest) (shared.AutomationWorkflowTagSetResponse, error) {
			odpowiedz, err := m.UstawEtykiety(ctx, z)
			if err == nil {
				e.powiazanieAutomatyki(shared.ChangeKindUpdated, z.WorkflowId)
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandAutomationWorkflowPublish,
		obsluz(func(ctx context.Context, z shared.AutomationWorkflowPublishRequest) (shared.AutomationWorkflowPublishResponse, error) {
			odpowiedz, err := m.Opublikuj(ctx, z)
			if err == nil {
				e.powiazanieAutomatyki(shared.ChangeKindUpdated, z.WorkflowId)
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandAutomationWorkflowShare,
		obsluz(func(ctx context.Context, z shared.AutomationWorkflowShareRequest) (shared.AutomationWorkflowShareResponse, error) {
			odpowiedz, err := m.Udostepnij(ctx, z)
			if err == nil {
				e.powiazanieAutomatyki(shared.ChangeKindUpdated, z.WorkflowId)
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandAutomationExecutionBudgetSet,
		obsluz(func(ctx context.Context, z shared.AutomationExecutionBudgetSetRequest) (shared.AutomationExecutionBudgetSetResponse, error) {
			odpowiedz, err := m.UstawBudzetyPrzebiegu(ctx, z)
			if err == nil {
				e.powiazanieAutomatyki(shared.ChangeKindUpdated, z.WorkflowId)
			}
			return odpowiedz, err
		}))

	// Zastosowanie szablonu ZAKŁADA automatykę, więc rozgłasza jej powstanie,
	// a identyfikator bierze z odpowiedzi — żądanie go nie zna, bo automatyki
	// jeszcze nie było.
	r.Zarejestruj(shared.CommandAutomationTemplateApply,
		obsluz(func(ctx context.Context, z shared.AutomationTemplateApplyRequest) (shared.AutomationTemplateApplyResponse, error) {
			odpowiedz, err := m.ZastosujSzablon(ctx, z)
			if err == nil {
				e.powiazanieAutomatyki(shared.ChangeKindCreated, odpowiedz.Workflow.Id)
			}
			return odpowiedz, err
		}))
}

// zarejestrujZmianyPrzebiegu wpina dwie komendy zmieniające stan przebiegu.
func zarejestrujZmianyPrzebiegu(r *Rejestr, m AutomatykiDobudowa, e *emiter) {
	r.Zarejestruj(shared.CommandAutomationExecutionResume,
		obsluz(func(ctx context.Context, z shared.AutomationExecutionResumeRequest) (shared.AutomationExecutionResumeResponse, error) {
			odpowiedz, err := m.WznowPrzebieg(ctx, z)
			if err == nil {
				e.przebiegAutomatyki(odpowiedz.Execution, etapWznowieniePrzebiegu)
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandAutomationExecutionReplay,
		obsluz(func(ctx context.Context, z shared.AutomationExecutionReplayRequest) (shared.AutomationExecutionReplayResponse, error) {
			odpowiedz, err := m.OdtworzPrzebieg(ctx, z)
			if err == nil {
				e.przebiegAutomatyki(odpowiedz.Execution, etapOdtworzeniePrzebiegu)
			}
			return odpowiedz, err
		}))
}

// Nazwy etapów przebiegu wprowadzone dobudową. Stoją obok etapów kolejki
// (`etapBiegNaprawczy`, `etapZatrzymanie`), żeby Execution Monitor mówił
// o wszystkich czynnościach na przebiegu jednym językiem.
const (
	etapWznowieniePrzebiegu  = "wznowienie od punktu"
	etapOdtworzeniePrzebiegu = "odtworzenie z ładunku"
)

// zarejestrujOdmoweDobudowyAutomatyk wpina wszystkie dwadzieścia sześć komend
// jako odmowę montażu — patrz uzasadnienie przy `zarejestrujDobudoweAutomatyk`.
func zarejestrujOdmoweDobudowyAutomatyk(r *Rejestr) {
	const powod = "moduł Automations: port modułu nie niesie czynności dobudowy"

	odmowa := func(typ shared.MessageType) {
		r.Zarejestruj(typ, func(context.Context, protocol.Request) protocol.Odpowiedz {
			return porazka(bladNiedostepnegoSilnika(powod))
		})
	}
	for _, typ := range []shared.MessageType{
		shared.CommandAutomationWorkflowSimulate,
		shared.CommandAutomationWorkflowVersionList,
		shared.CommandAutomationWorkflowVersionRestore,
		shared.CommandAutomationWorkflowVersionDiff,
		shared.CommandAutomationWorkflowTagSet,
		shared.CommandAutomationWorkflowVariablesSet,
		shared.CommandAutomationStepNoteSet,
		shared.CommandAutomationStepLayoutSet,
		shared.CommandAutomationTemplateSave,
		shared.CommandAutomationTemplateList,
		shared.CommandAutomationTemplateApply,
		shared.CommandAutomationWorkflowPublish,
		shared.CommandAutomationWorkflowShare,
		shared.CommandAutomationExecutionLog,
		shared.CommandAutomationExecutionSteps,
		shared.CommandAutomationExecutionCheckpointList,
		shared.CommandAutomationExecutionResume,
		shared.CommandAutomationExecutionPayloadGet,
		shared.CommandAutomationExecutionReplay,
		shared.CommandAutomationAlertRuleSet,
		shared.CommandAutomationAlertRuleList,
		shared.CommandAutomationExecutionBudgetSet,
		shared.CommandAutomationSecretSet,
		shared.CommandAutomationSecretList,
		shared.CommandAutomationSecretRemove,
		shared.CommandAutomationAuditList,
	} {
		odmowa(typ)
	}
}
