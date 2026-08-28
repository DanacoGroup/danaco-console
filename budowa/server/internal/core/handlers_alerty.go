// Plik daje port i wpięcie pięciu komend rodziny alert.* — reguł wyzwalania, rejestru
// wyzwoleń i ich potwierdzania. Port bierze nadajnik, bo rodzina ma zdarzenie alert.triggered.
package core

import (
	"context"

	"danacoconsole/shared"
)

// Alerty jest portem rodziny alert.*, obsługującym reguły wyzwalania, wykaz wyzwoleń
// i ich potwierdzanie.
type Alerty interface {
	// ZapiszRegule obsługuje `alert.rule.save`.
	ZapiszRegule(ctx context.Context, z shared.AlertRuleSaveRequest) (shared.AlertRuleSaveResponse, error)
	// WykazRegul obsługuje `alert.rule.list`.
	WykazRegul(ctx context.Context, z shared.AlertRuleListRequest) (shared.AlertRuleListResponse, error)
	// UsunRegule obsługuje `alert.rule.remove`.
	UsunRegule(ctx context.Context, z shared.AlertRuleRemoveRequest) (shared.AlertRuleRemoveResponse, error)
	// WykazWyzwolen obsługuje `alert.trigger.list`.
	WykazWyzwolen(ctx context.Context, z shared.AlertTriggerListRequest) (shared.AlertTriggerListResponse, error)
	// PotwierdzWyzwolenie obsługuje `alert.trigger.acknowledge`.
	PotwierdzWyzwolenie(ctx context.Context, z shared.AlertTriggerAcknowledgeRequest) (shared.AlertTriggerAcknowledgeResponse, error)
}

// Adapter wypełnia port Alerty w całości, udostępniając komendom wszystkie jego pięć
// metod rodziny alert.
var _ Alerty = (*adapterAlertow)(nil)

// zarejestrujAlerty wpina pięć komend rodziny alert.*. Zdarzenie alert.triggered rozgłasza
// sam adapter przy zapisie wyzwolenia, nadajnikiem wpiętym metodą ZWyjsciem, nie obsługa
// komend.
func zarejestrujAlerty(r *Rejestr, a Alerty) {
	if r == nil || a == nil {
		return
	}

	r.Zarejestruj(shared.CommandAlertRuleSave, obsluz(a.ZapiszRegule))
	r.Zarejestruj(shared.CommandAlertRuleList, obsluz(a.WykazRegul))
	r.Zarejestruj(shared.CommandAlertRuleRemove, obsluz(a.UsunRegule))
	r.Zarejestruj(shared.CommandAlertTriggerList, obsluz(a.WykazWyzwolen))
	r.Zarejestruj(shared.CommandAlertTriggerAcknowledge, obsluz(a.PotwierdzWyzwolenie))
}
