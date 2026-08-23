// Odpowiedzialność pliku: port i wpięcie pięciu komend rodziny `alert.*` —
// reguł wyzwalania, rejestru wyzwoleń i ich potwierdzania.
//
// Rodzina jest przekrojowa, jak kondycja: reguły czyta Alerts Panel modułu
// Diagnostics, ale wyzwolenie dociera też do Always On Display i do poczty,
// a rdzeń nie ma prawa wiedzieć, że istnieje moduł Diagnostics.
//
// Port bierze nadajnik, bo rodzina MA zdarzenie: `alert.triggered` jest drogą
// alertu do okien. Bez niego Operator dowiadywałby się o wyzwoleniu dopiero
// przy następnym otwarciu wykazu, czyli wtedy, kiedy i tak już patrzy.
//
// Port niewypełniony nie rejestruje niczego: pięć komend odpowie wtedy
// `alert.unknown`, a pozostałe domeny pracują bez zmian.
package core

import (
	"context"

	"danacoconsole/shared"
)

// Alerty jest portem rodziny `alert.*`.
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

// Adapter wypełnia port w całości.
var _ Alerty = (*adapterAlertow)(nil)

// zarejestrujAlerty wpina pięć komend rodziny `alert.*`.
//
// Zdarzenie `alert.triggered` rozgłasza sam adapter przy zapisie wyzwolenia,
// nadajnikiem wpiętym metodą `ZWyjsciem` — nie obsługa komend. Powód jest
// prosty: wyzwolenie powstaje w trakcie ewaluacji, a nie w odpowiedzi na
// komendę, więc nadanie zdarzenia z obsługi komendy pominęłoby wyzwolenia,
// których nikt akurat nie odczytywał.
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
