// Odpowiedzialność pliku: przypisanie eksperta do projektu — okno Agent
// Manager modułu Workspace. Ekspert mieszka w module Agents, przypisanie w
// module Workspace. Tutaj zapisuje się wyłącznie fakt, że ekspert pracuje w
// projekcie.
package core

import (
	"context"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// PrzypiszAgenta zapisuje przypisanie eksperta do projektu. Wskazanie
// domyślnego wykonawcy zdejmuje je z poprzedniego — wykonawca domyślny jest
// jeden (repozytorium pilnuje tego w jednej transakcji).
func (a *adapterPrzestrzeniRoboczej) PrzypiszAgenta(ctx context.Context,
	z shared.WorkspaceAgentAssignRequest) (shared.WorkspaceAgentAssignResponse, error) {

	projekt, err := a.projektDlaZapisu(ctx, z.ProjectId)
	if err != nil {
		return shared.WorkspaceAgentAssignResponse{}, err
	}
	if z.AgentId == "" {
		return shared.WorkspaceAgentAssignResponse{}, bladProjektu("przypisanie bez wskazania eksperta")
	}
	przypisanie := dane.PrzypisanieAgenta{
		ProjektID:         projekt.ID,
		AgentKod:          z.AgentId,
		Rola:              z.Role,
		DomyslnyWykonawca: z.DefaultExecutor != nil && *z.DefaultExecutor,
	}
	zapisane, err := a.repozytorium.PrzypiszAgenta(ctx, przypisanie)
	if err != nil {
		return shared.WorkspaceAgentAssignResponse{}, err
	}
	if err := a.repozytorium.OdnotujCzynnosc(ctx, projekt.ID); err != nil {
		return shared.WorkspaceAgentAssignResponse{}, err
	}
	return shared.WorkspaceAgentAssignResponse{
		Assignment: przypisanieKontraktu(projekt.Kod, zapisane),
	}, nil
}

// przypisanieKontraktu przekłada wiersz przypisania eksperta z bazy na byt
// kontraktu widoczny w oknie.
func przypisanieKontraktu(idProjektu string, p dane.PrzypisanieAgenta) shared.WorkspaceAgentAssignment {
	domyslny := p.DomyslnyWykonawca
	return shared.WorkspaceAgentAssignment{
		ProjectId:       idProjektu,
		AgentId:         p.AgentKod,
		Role:            p.Rola,
		DefaultExecutor: &domyslny,
		AssignedAt:      chwilaBazy(p.Przypisano),
	}
}
