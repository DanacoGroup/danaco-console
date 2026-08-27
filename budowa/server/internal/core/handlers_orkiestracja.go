// Plik wpina cztery komendy rodziny subagent.* obsługujące powołanie, wykaz i zbieranie wyników
// podagentów okna wykonawcy dla panelu Subagent Network.
package core

import (
	"context"

	"danacoconsole/shared"
)

// Podagenci jest portem rodziny `subagent.*`. Wypełnia go adapter nad
// repozytorium podagentów i silnikiem kolejek; rdzeń nie zna implementacji.
type Podagenci interface {
	// Powolaj uruchamia podagentów okna; górna granica piętnastu przycina, nie odmawia żądania.
	Powolaj(ctx context.Context, z shared.SubagentSpawnRequest) (shared.SubagentSpawnResponse, error)
	// Wykaz oddaje podagentów okna albo karty sesji; wykaz pusty jest poprawną odpowiedzią.
	Wykaz(ctx context.Context, z shared.SubagentListRequest) (shared.SubagentListResponse, error)
	// ZbierzWyniki oddaje podagentów z wynikami i orzeka, czy wszyscy objęci zbieraniem zakończyli.
	ZbierzWyniki(ctx context.Context, z shared.SubagentResultCollectRequest) (shared.SubagentResultCollectResponse, error)
	// Zatrzymaj zatrzymuje wskazanych podagentów; zakończony podagent nie jest błędem.
	Zatrzymaj(ctx context.Context, z shared.SubagentStopRequest) (shared.SubagentStopResponse, error)
}

// zarejestrujOrkiestracjePodagentow wpina cztery komendy rodziny `subagent.*`.
//
// Port niewypełniony nie zatrzymuje startu: komendy odpowiadają wtedy
// `subagent.unknown`, a pozostałe domeny pracują bez zmian.
func zarejestrujOrkiestracjePodagentow(r *Rejestr, m Podagenci) {
	if r == nil || m == nil {
		return
	}
	r.Zarejestruj(shared.CommandSubagentSpawn, obsluz(m.Powolaj))
	r.Zarejestruj(shared.CommandSubagentList, obsluz(m.Wykaz))
	r.Zarejestruj(shared.CommandSubagentResultCollect, obsluz(m.ZbierzWyniki))
	r.Zarejestruj(shared.CommandSubagentStop, obsluz(m.Zatrzymaj))
}
