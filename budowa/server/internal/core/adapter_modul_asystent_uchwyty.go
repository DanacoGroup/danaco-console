// Odpowiedzialność pliku: wpięcie komend modułu Assistant przez port Asystent, którym
// rejestr komend rdzenia dociera do adaptera wypełnianego w dwóch plikach jednego typu.
package core

import (
	"context"

	"danacoconsole/shared"
)

// Asystent jest portem modułu Assistant, wymieniającym metody obsługi jego czterech komend rdzenia platformy.
type Asystent interface {
	PolecenieGlosowe(ctx context.Context, z shared.AssistantVoiceCommandRequest) (shared.AssistantVoiceCommandResponse, error)
	StanCzynnosci(ctx context.Context, z shared.AssistantActionStatusRequest) (shared.AssistantActionStatusResponse, error)
	WykazCzynnosci(ctx context.Context, z shared.AssistantActivityListRequest) (shared.AssistantActivityListResponse, error)
	OznaczWpisDziennika(ctx context.Context, z shared.AssistantActivityFlagRequest) (shared.AssistantActivityFlagResponse, error)
}

// zarejestrujAsystenta wpina komendy modułu Assistant w rejestr komend rdzenia platformy Danaco Console.
func zarejestrujAsystenta(r *Rejestr, m Asystent, e *emiter) {
	if r == nil || m == nil {
		return
	}

	r.Zarejestruj(shared.CommandAssistantVoiceCommand, obsluz(m.PolecenieGlosowe))
	r.Zarejestruj(shared.CommandAssistantActionStatus, obsluz(m.StanCzynnosci))
	r.Zarejestruj(shared.CommandAssistantActivityList, obsluz(m.WykazCzynnosci))
	r.Zarejestruj(shared.CommandAssistantActivityFlag, obsluz(m.OznaczWpisDziennika))

	// Zdarzenie zmiany zlecenia rozgłasza wykonawca zleceń własnym emiterem, nie obsługa komend.
}
