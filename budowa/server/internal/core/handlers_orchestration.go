// Plik wpina cztery komendy rodziny orchestration.* obsługujące układ zależności, rozszerzając
// port modułu Automations o operacje na łuku układu.
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
				e.ukladOrkiestracji(ctx, zmiana, wskazanieLuku(z.Dependency))
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandOrchestrationDependencyList, obsluz(m.WykazZaleznosci))

	r.Zarejestruj(shared.CommandOrchestrationDependencyRemove,
		obsluz(func(ctx context.Context, z shared.OrchestrationDependencyRemoveRequest) (shared.OrchestrationDependencyRemoveResponse, error) {
			odpowiedz, err := m.UsunZaleznoscUkladu(ctx, z)
			if err == nil {
				e.ukladOrkiestracji(ctx, shared.ChangeKindDeleted,
					wskazanieLuku(shared.AutomationDependency{FromStepId: z.FromStepId, ToStepId: z.ToStepId}))
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandOrchestrationValidate, obsluz(m.SprawdzUkladZaleznosci))

	// Cztery dopełnienia układu rozgłaszają zmianę układu bez wskazania łuku, bo dotyczą całości.
	r.Zarejestruj(shared.CommandOrchestrationGateSet,
		obsluz(func(ctx context.Context, z shared.OrchestrationGateSetRequest) (shared.OrchestrationGateSetResponse, error) {
			odpowiedz, err := m.UstawBramke(ctx, z)
			if err == nil {
				e.ukladOrkiestracji(ctx, shared.ChangeKindUpdated, "")
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandOrchestrationGroupSet,
		obsluz(func(ctx context.Context, z shared.OrchestrationGroupSetRequest) (shared.OrchestrationGroupSetResponse, error) {
			odpowiedz, err := m.UstawGrupe(ctx, z)
			if err == nil {
				e.ukladOrkiestracji(ctx, shared.ChangeKindUpdated, "")
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandOrchestrationCompensationSet,
		obsluz(func(ctx context.Context, z shared.OrchestrationCompensationSetRequest) (shared.OrchestrationCompensationSetResponse, error) {
			odpowiedz, err := m.UstawKompensacje(ctx, z)
			if err == nil {
				e.ukladOrkiestracji(ctx, shared.ChangeKindUpdated, "")
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandOrchestrationMultitaskingLink,
		obsluz(func(ctx context.Context, z shared.OrchestrationMultitaskingLinkRequest) (shared.OrchestrationMultitaskingLinkResponse, error) {
			odpowiedz, err := m.SpnijZMultitaskingiem(ctx, z)
			if err == nil {
				e.ukladOrkiestracji(ctx, shared.ChangeKindUpdated, "")
			}
			return odpowiedz, err
		}))
}

// rodzajZmianyLuku orzeka, czy zapis zależności łuk dokłada, czy przepisuje, pytając o układ
// przed zapisem, bo odpowiedź komendy tego nie mówi.
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

// wskazanieLuku składa identyfikator zależności z obu jej końców, bo łuk nie ma własnego klucza
// w tabeli zależności kroku automatyki.
func wskazanieLuku(luk shared.AutomationDependency) string {
	if luk.FromStepId == "" || luk.ToStepId == "" {
		return ""
	}
	return luk.FromStepId + "→" + luk.ToStepId
}

// ukladOrkiestracji rozgłasza orchestration.changed dla komponentu własnego automatyki, bez
// wskazania karty sesji, tak samo jak przebieg automatyki.
func (e *emiter) ukladOrkiestracji(ctx context.Context, zmiana shared.ChangeKind, idZaleznosci string) {
	tresc := shared.OrchestrationChangedEvent{Change: zmiana}
	if idZaleznosci != "" {
		wskazanie := idZaleznosci
		tresc.DependencyId = &wskazanie
	}
	e.wyslijDoKonta(ctx, shared.EventOrchestrationChanged, "", tresc)
}
