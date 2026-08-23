// Odpowiedzialność pliku: wpięcie czterech komend rodziny `subagent.*` —
// powołania podagentów okna wykonawcy, ich wykazu dla panelu Subagent Network
// i zbierania wyników ich pracy.
//
// Rodzina ma port osobny, nie rozszerzenie portu Automatyk: podagent nie jest
// krokiem automatyki ani jej przebiegiem — wisi na oknie wykonawcy, a nie na
// zapisanej definicji, i powstaje w rozmowie, nie w Workflow Builderze.
//
// Zdarzenia `subagent.changed` nie rozgłasza ten plik, tylko adapter: podagent
// zmienia stan głównie poza żądaniem — wejście w `running`, zakończenie
// i niepowodzenie dzieją się w pracy puszczonej w tle, której uchwyt komendy
// nie widzi. Treść i reguła rodzaju zmiany:
// `adapter_modul_orkiestracja_rozgloszenie.go`. Dlatego funkcja wpinająca nie
// bierze nadawcy.
package core

import (
	"context"

	"danacoconsole/shared"
)

// Podagenci jest portem rodziny `subagent.*`. Wypełnia go adapter nad
// repozytorium podagentów i silnikiem kolejek; rdzeń nie zna implementacji.
type Podagenci interface {
	// Powolaj uruchamia podagentów okna wykonawcy. Górna granica piętnastu jest
	// przycięciem, nie odmową: żądanie o większej liczbie wykonuje się na
	// piętnastu.
	Powolaj(ctx context.Context, z shared.SubagentSpawnRequest) (shared.SubagentSpawnResponse, error)
	// Wykaz oddaje podagentów okna albo karty sesji, w razie wskazania zawężonych
	// stanem. Wykaz pusty jest poprawną odpowiedzią, nie odmową.
	Wykaz(ctx context.Context, z shared.SubagentListRequest) (shared.SubagentListResponse, error)
	// ZbierzWyniki oddaje podagentów wraz z wynikami i orzeka, czy wszyscy objęci
	// zbieraniem zakończyli pracę.
	ZbierzWyniki(ctx context.Context, z shared.SubagentResultCollectRequest) (shared.SubagentResultCollectResponse, error)
	// Zatrzymaj zatrzymuje wskazanych podagentów, nie ruszając tury okna ani
	// podagentów niewskazanych. Podagent już zakończony nie jest błędem —
	// wraca w wykazie nieczynnych.
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
