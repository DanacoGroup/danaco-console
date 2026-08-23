// Odpowiedzialność pliku: wpięcie czterech komend rodziny `orchestration.*`.
//
// Rodzina wpina się osobno od `adapter_modul_automations_uchwyty.go`, choć jedzie
// na tej samej maszynerii układu zależności: port poniżej rozszerza port
// Automatyki, bo układ zależności ma w rdzeniu jednego właściciela.
//
// Rozdział miejsc emisji zdarzeń układu:
//
//   - `orchestration.dependency.set` → `created`/`updated` ze wskazaniem łuku;
//     `orchestration.dependency.remove` → `deleted`; oba tutaj;
//   - `automation.orchestrator.define` → `updated` bez wskazania zależności
//     (przepisuje cały układ), `DependencyId` zostaje wtedy puste;
//   - `automation.schedule.set` → `automation.link.changed` z `AutomationId`
//     równym `workflowId`.
//
// Dwa ostatnie jadą z `adapter_modul_automations_uchwyty.go`.
package core

import (
	"context"

	"danacoconsole/shared"
)

// Orkiestracja jest portem rodziny `orchestration.*` — portem modułu
// Automations rozszerzonym o cztery komendy układu zależności.
type Orkiestracja interface {
	Automatyki

	// UstawZaleznosc obsługuje `orchestration.dependency.set`.
	UstawZaleznosc(ctx context.Context, z shared.OrchestrationDependencySetRequest) (shared.OrchestrationDependencySetResponse, error)
	// WykazZaleznosci obsługuje `orchestration.dependency.list`.
	WykazZaleznosci(ctx context.Context, z shared.OrchestrationDependencyListRequest) (shared.OrchestrationDependencyListResponse, error)
	// UsunZaleznoscUkladu obsługuje `orchestration.dependency.remove`.
	UsunZaleznoscUkladu(ctx context.Context, z shared.OrchestrationDependencyRemoveRequest) (shared.OrchestrationDependencyRemoveResponse, error)
	// SprawdzUkladZaleznosci obsługuje `orchestration.validate`.
	SprawdzUkladZaleznosci(ctx context.Context, z shared.OrchestrationValidateRequest) (shared.OrchestrationValidateResponse, error)

	// Cztery dopełnienia układu, których łuk nie wyraża
	// (`adapter_modul_orkiestracja_uklad.go`).
	UstawBramke(ctx context.Context, z shared.OrchestrationGateSetRequest) (shared.OrchestrationGateSetResponse, error)
	UstawGrupe(ctx context.Context, z shared.OrchestrationGroupSetRequest) (shared.OrchestrationGroupSetResponse, error)
	UstawKompensacje(ctx context.Context, z shared.OrchestrationCompensationSetRequest) (shared.OrchestrationCompensationSetResponse, error)
	SpnijZMultitaskingiem(ctx context.Context, z shared.OrchestrationMultitaskingLinkRequest) (shared.OrchestrationMultitaskingLinkResponse, error)
}

// zarejestrujOrkiestracje wpina cztery komendy rodziny `orchestration.*`.
// Dwie z nich zmieniają układ i po udanej zmianie rozgłaszają
// `orchestration.changed`; dwie pozostałe są odczytem i milczą.
func zarejestrujOrkiestracje(r *Rejestr, m Orkiestracja, e *emiter) {
	if r == nil || m == nil {
		return
	}

	r.Zarejestruj(shared.CommandOrchestrationDependencySet,
		obsluz(func(ctx context.Context, z shared.OrchestrationDependencySetRequest) (shared.OrchestrationDependencySetResponse, error) {
			zmiana := rodzajZmianyLuku(ctx, m, z)
			odpowiedz, err := m.UstawZaleznosc(ctx, z)
			if err == nil {
				e.ukladOrkiestracji(zmiana, wskazanieLuku(z.Dependency))
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandOrchestrationDependencyList, obsluz(m.WykazZaleznosci))

	r.Zarejestruj(shared.CommandOrchestrationDependencyRemove,
		obsluz(func(ctx context.Context, z shared.OrchestrationDependencyRemoveRequest) (shared.OrchestrationDependencyRemoveResponse, error) {
			odpowiedz, err := m.UsunZaleznoscUkladu(ctx, z)
			if err == nil {
				e.ukladOrkiestracji(shared.ChangeKindDeleted,
					wskazanieLuku(shared.AutomationDependency{FromStepId: z.FromStepId, ToStepId: z.ToStepId}))
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandOrchestrationValidate, obsluz(m.SprawdzUkladZaleznosci))

	// Cztery dopełnienia układu. Każde zmienia układ, więc każde rozgłasza
	// `orchestration.changed` — bez wskazania łuku, bo bramka, grupa,
	// kompensacja i spięcie nie dotyczą żadnej jednej pary kroków. Okno czyta
	// takie zdarzenie jako „przelicz układ od nowa" (patrz `ukladOrkiestracji`).
	r.Zarejestruj(shared.CommandOrchestrationGateSet,
		obsluz(func(ctx context.Context, z shared.OrchestrationGateSetRequest) (shared.OrchestrationGateSetResponse, error) {
			odpowiedz, err := m.UstawBramke(ctx, z)
			if err == nil {
				e.ukladOrkiestracji(shared.ChangeKindUpdated, "")
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandOrchestrationGroupSet,
		obsluz(func(ctx context.Context, z shared.OrchestrationGroupSetRequest) (shared.OrchestrationGroupSetResponse, error) {
			odpowiedz, err := m.UstawGrupe(ctx, z)
			if err == nil {
				e.ukladOrkiestracji(shared.ChangeKindUpdated, "")
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandOrchestrationCompensationSet,
		obsluz(func(ctx context.Context, z shared.OrchestrationCompensationSetRequest) (shared.OrchestrationCompensationSetResponse, error) {
			odpowiedz, err := m.UstawKompensacje(ctx, z)
			if err == nil {
				e.ukladOrkiestracji(shared.ChangeKindUpdated, "")
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandOrchestrationMultitaskingLink,
		obsluz(func(ctx context.Context, z shared.OrchestrationMultitaskingLinkRequest) (shared.OrchestrationMultitaskingLinkResponse, error) {
			odpowiedz, err := m.SpnijZMultitaskingiem(ctx, z)
			if err == nil {
				e.ukladOrkiestracji(shared.ChangeKindUpdated, "")
			}
			return odpowiedz, err
		}))
}

// rodzajZmianyLuku orzeka, czy `dependency.set` łuk dokłada, czy przepisuje.
// Odpowiedź komendy tego nie mówi (oddaje cały układ po zapisie), więc pyta się
// o układ PRZED zapisem. Nieudany odczyt nie wstrzymuje niczego — komenda ma
// się wykonać tak samo, a rodzaj zmiany schodzi wtedy do `updated`, bo
// „zmieniło się" jest zdaniem prawdziwym w obu przypadkach.
func rodzajZmianyLuku(ctx context.Context, m Orkiestracja,
	z shared.OrchestrationDependencySetRequest) shared.ChangeKind {

	wykaz, err := m.WykazZaleznosci(ctx, shared.OrchestrationDependencyListRequest{WorkflowId: z.WorkflowId})
	if err != nil {
		return shared.ChangeKindUpdated
	}
	for _, luk := range wykaz.Dependencies {
		if luk.FromStepId == z.Dependency.FromStepId && luk.ToStepId == z.Dependency.ToStepId {
			return shared.ChangeKindUpdated
		}
	}
	return shared.ChangeKindCreated
}

// wskazanieLuku składa identyfikator zależności z obu jej końców.
//
// Łuk nie ma własnego klucza: tabela `zaleznosc_kroku_automatyki` rozpoznaje go
// parą (krok_z, krok_do) i tą samą parą posługują się obie komendy kontraktu.
// Pole `dependencyId` jest złożeniem tej pary; łuk bez obu końców wskazania
// nie ma.
func wskazanieLuku(luk shared.AutomationDependency) string {
	if luk.FromStepId == "" || luk.ToStepId == "" {
		return ""
	}
	return luk.FromStepId + "→" + luk.ToStepId
}

// ukladOrkiestracji rozgłasza `orchestration.changed`. Układ zależności jest
// komponentem własnym automatyki, nie bytem karty sesji, więc zdarzenie idzie
// bez jej wskazania — tak samo jak przebieg automatyki.
//
// Puste wskazanie jest stanem zamierzonym: `dependencyId` jest polem
// nieobowiązkowym, bo przepisanie całego układu
// (`automation.orchestrator.define`) nie dotyczy żadnego jednego łuku. Okno
// czyta wtedy zdarzenie jako „przelicz układ od nowa”.
func (e *emiter) ukladOrkiestracji(zmiana shared.ChangeKind, idZaleznosci string) {
	tresc := shared.OrchestrationChangedEvent{Change: zmiana}
	if idZaleznosci != "" {
		wskazanie := idZaleznosci
		tresc.DependencyId = &wskazanie
	}
	e.wyslij(shared.EventOrchestrationChanged, "", tresc)
}
