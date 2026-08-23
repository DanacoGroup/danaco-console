// Odpowiedzialność pliku: wpięcie trzech komend modułu Assistant — portu
// `Asystent`, przez który rejestr komend rdzenia dociera do adaptera
// wypełnianego w `adapter_modul_asystent.go` (przyjęcie polecenia głosowego)
// oraz w `adapter_modul_asystent_czynnosci.go` (stan zleceń i dziennik
// działań) — ten sam typ `adapterAsystenta`, dwa pliki.
//
// Port wymienia wszystkie trzy komendy, wzorem
// `adapter_modul_library_uchwyty.go`: rejestr rdzenia potrzebuje jednego
// miejsca wiążącego nazwę komendy z metodą portu, niezależnie od tego, który
// plik adaptera którą metodę implementuje.
package core

import (
	"context"

	"danacoconsole/shared"
)

// Asystent jest portem modułu Assistant.
type Asystent interface {
	PolecenieGlosowe(ctx context.Context, z shared.AssistantVoiceCommandRequest) (shared.AssistantVoiceCommandResponse, error)
	StanCzynnosci(ctx context.Context, z shared.AssistantActionStatusRequest) (shared.AssistantActionStatusResponse, error)
	WykazCzynnosci(ctx context.Context, z shared.AssistantActivityListRequest) (shared.AssistantActivityListResponse, error)
	OznaczWpisDziennika(ctx context.Context, z shared.AssistantActivityFlagRequest) (shared.AssistantActivityFlagResponse, error)
}

// zarejestrujAsystenta wpina cztery komendy modułu Assistant.
func zarejestrujAsystenta(r *Rejestr, m Asystent, e *emiter) {
	if r == nil || m == nil {
		return
	}

	r.Zarejestruj(shared.CommandAssistantVoiceCommand, obsluz(m.PolecenieGlosowe))
	r.Zarejestruj(shared.CommandAssistantActionStatus, obsluz(m.StanCzynnosci))
	r.Zarejestruj(shared.CommandAssistantActivityList, obsluz(m.WykazCzynnosci))
	r.Zarejestruj(shared.CommandAssistantActivityFlag, obsluz(m.OznaczWpisDziennika))

	// Zdarzenie `assistant.action.changed` rozgłasza wykonawca zleceń
	// (`adapter_modul_asystent_wykonawca.go`) przy zmianie stanu podjętego
	// zlecenia, własnym emiterem na nadajniku wpiętym metodą `ZWyjsciem`, a nie
	// obsługa komend. Trzy komendy modułu (polecenie, stan, dziennik) nic nie
	// rozgłaszają same, więc parametr `e` zostaje nieużyty, na wzór pozostałych
	// `zarejestruj<Modul>` (Automatyki, Biblioteka).
}
