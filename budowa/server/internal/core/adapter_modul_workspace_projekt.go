// Odpowiedzialność pliku: czynności na samym projekcie i na jego
// wykonawcach — stan projektu, odłączenie eksperta oraz historia instrukcji
// systemowych. Tu leży też odkładanie zdarzeń osi czasu, wspólne modułowi.
package core

import (
	"context"
	"errors"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// UstawStanProjektu obsługuje `workspace.project.status.set` i odnotowuje
// zmianę stanu na osi czasu projektu.
func (a *adapterPrzestrzeniRoboczej) UstawStanProjektu(ctx context.Context,
	z shared.WorkspaceProjectStatusSetRequest) (shared.WorkspaceProjectStatusSetResponse, error) {

	projekt, err := a.projektDlaZapisu(ctx, z.ProjectId)
	if err != nil {
		return shared.WorkspaceProjectStatusSetResponse{}, err
	}
	stan := z.Status
	if !znanyStanProjektuWorkspace(stan) {
		return shared.WorkspaceProjectStatusSetResponse{},
			bladProjektu("stan projektu " + string(stan) + " nie należy do stanów kontraktu")
	}
	if err := a.repozytorium.UstawStanProjektuWorkspace(ctx, projekt.ID, stan); err != nil {
		return shared.WorkspaceProjectStatusSetResponse{}, err
	}
	a.odnotujZdarzenieWorkspace(ctx, projekt.ID, shared.WorkspaceActivityKindProjectChanged,
		shared.ChangeKindUpdated, shared.WorkspaceEntityKindProject, projekt.Kod,
		"stan projektu ustawiony na "+string(stan))

	poZmianie, err := a.repozytorium.Projekt(ctx, projekt.Kod)
	if err != nil {
		return shared.WorkspaceProjectStatusSetResponse{}, err
	}
	return shared.WorkspaceProjectStatusSetResponse{Project: projektKontraktu(poZmianie)}, nil
}

// OdlaczAgenta obsługuje `workspace.agent.unassign`. Brak przypisania nie jest
// odmową: czynność miała doprowadzić do stanu „ekspert nie pracuje w projekcie"
// i ten stan zastała.
func (a *adapterPrzestrzeniRoboczej) OdlaczAgenta(ctx context.Context,
	z shared.WorkspaceAgentUnassignRequest) (shared.WorkspaceAgentUnassignResponse, error) {

	projekt, err := a.projektDlaZapisu(ctx, z.ProjectId)
	if err != nil {
		return shared.WorkspaceAgentUnassignResponse{}, err
	}
	if strings.TrimSpace(z.AgentId) == "" {
		return shared.WorkspaceAgentUnassignResponse{},
			bladProjektu("odłączenie bez wskazania eksperta")
	}
	odlaczony, err := a.repozytorium.OdlaczAgentaWorkspace(ctx, projekt.ID, z.AgentId)
	if err != nil {
		return shared.WorkspaceAgentUnassignResponse{}, err
	}
	if odlaczony {
		if err := a.repozytorium.OdnotujCzynnosc(ctx, projekt.ID); err != nil {
			return shared.WorkspaceAgentUnassignResponse{}, err
		}
		a.odnotujZdarzenieWorkspace(ctx, projekt.ID,
			shared.WorkspaceActivityKindAgentAssignmentChanged, shared.ChangeKindDeleted,
			shared.WorkspaceEntityKindProject, projekt.Kod,
			"ekspert "+z.AgentId+" odłączony od projektu")
	}
	return shared.WorkspaceAgentUnassignResponse{Unassigned: odlaczony}, nil
}

// WersjeInstrukcji obsługuje `workspace.instructions.version.list` i oddaje
// historię instrukcji systemowych projektu.
func (a *adapterPrzestrzeniRoboczej) WersjeInstrukcji(ctx context.Context,
	z shared.WorkspaceInstructionsVersionListRequest) (shared.WorkspaceInstructionsVersionListResponse, error) {

	projekt, err := a.projektDlaZapisu(ctx, z.ProjectId)
	if err != nil {
		return shared.WorkspaceInstructionsVersionListResponse{}, err
	}
	wszystkie, err := a.repozytorium.WersjeInstrukcjiWorkspace(ctx, projekt.ID)
	if err != nil {
		return shared.WorkspaceInstructionsVersionListResponse{}, err
	}
	poziom := shared.ConfigScope(shared.ConfigScopeProject)
	if z.Scope != nil && *z.Scope != "" {
		poziom = *z.Scope
	}
	bytPoziomu := ""
	if z.ScopeId != nil {
		bytPoziomu = *z.ScopeId
	}
	dobrane := []shared.WorkspaceInstructionsVersion{}
	for _, wersja := range wszystkie {
		if wersja.Poziom != poziom {
			continue
		}
		if bytPoziomu != "" && wersja.KluczZasiegu != bytPoziomu {
			continue
		}
		dobrane = append(dobrane, wersjaInstrukcjiKontraktuWorkspace(wersja))
	}
	wszystkich := len(dobrane)
	return shared.WorkspaceInstructionsVersionListResponse{
		Versions: przytnijWykazWorkspace(dobrane, z.Limit, z.Offset),
		Total:    wszystkich,
	}, nil
}

// PrzywrocInstrukcje obsługuje `workspace.instructions.version.restore`.
// Przywrócenie zakłada wersję NOWĄ o treści wersji wskazanej.
func (a *adapterPrzestrzeniRoboczej) PrzywrocInstrukcje(ctx context.Context,
	z shared.WorkspaceInstructionsVersionRestoreRequest) (shared.WorkspaceInstructionsVersionRestoreResponse, error) {

	projekt, err := a.projektDlaZapisu(ctx, z.ProjectId)
	if err != nil {
		return shared.WorkspaceInstructionsVersionRestoreResponse{}, err
	}
	zrodlowa, err := a.repozytorium.WersjaInstrukcjiWorkspace(ctx, z.VersionId)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.WorkspaceInstructionsVersionRestoreResponse{},
			bladProjektu("wersji instrukcji " + z.VersionId + " nie ma w historii projektu")
	}
	if err != nil {
		return shared.WorkspaceInstructionsVersionRestoreResponse{}, err
	}
	if zrodlowa.ProjektID != projekt.ID {
		return shared.WorkspaceInstructionsVersionRestoreResponse{},
			bladProjektu("wersja instrukcji " + z.VersionId + " należy do innego projektu")
	}

	poziom := zrodlowa.Poziom
	bytPoziomu := zrodlowa.KluczZasiegu
	odpowiedz, err := a.ZapiszInstrukcje(ctx, shared.WorkspaceInstructionsSetRequest{
		ProjectId: projekt.Kod, Content: zrodlowa.Tresc,
		Scope: &poziom, ScopeId: &bytPoziomu,
	})
	if err != nil {
		return shared.WorkspaceInstructionsVersionRestoreResponse{}, err
	}
	nowa, err := a.repozytorium.ZapiszWersjeInstrukcjiWorkspace(ctx, dane.WersjaInstrukcjiWorkspace{
		ProjektID: projekt.ID, Identyfikator: nowyIdentyfikator("wsiv-"),
		Tresc: zrodlowa.Tresc, Odcisk: odciskTresci(zrodlowa.Tresc),
		Poziom: poziom, KluczZasiegu: bytPoziomu, PrzywroconoZ: zrodlowa.Identyfikator,
	})
	if err != nil {
		return shared.WorkspaceInstructionsVersionRestoreResponse{}, err
	}
	a.odnotujZdarzenieWorkspace(ctx, projekt.ID, shared.WorkspaceActivityKindInstructionsChanged,
		shared.ChangeKindUpdated, shared.WorkspaceEntityKindInstructions, nowa.Identyfikator,
		"instrukcje przywrócone z wersji "+zrodlowa.Identyfikator)

	return shared.WorkspaceInstructionsVersionRestoreResponse{
		Instructions: odpowiedz.Instructions,
		Version:      wersjaInstrukcjiKontraktuWorkspace(nowa),
	}, nil
}

// odnotujZdarzenieWorkspace odkłada zdarzenie osi czasu projektu;
// niepowodzenie zapisu nie przerywa czynności wywołującej.
func (a *adapterPrzestrzeniRoboczej) odnotujZdarzenieWorkspace(ctx context.Context, projektID int64,
	rodzaj shared.WorkspaceActivityKind, zmiana shared.ChangeKind,
	rodzajBytu shared.WorkspaceEntityKind, byt, opis string) {

	_ = a.repozytorium.ZapiszZdarzenieWorkspace(ctx, dane.ZdarzenieWorkspace{
		ProjektID: projektID, Identyfikator: nowyIdentyfikator("wsev-"),
		Rodzaj: rodzaj, Zmiana: zmiana, RodzajBytu: rodzajBytu, Byt: byt, Opis: opis,
	})
}

// wersjaInstrukcjiKontraktuWorkspace przekłada wiersz historii instrukcji
// na byt kontraktu zwracany wołającemu.
func wersjaInstrukcjiKontraktuWorkspace(w dane.WersjaInstrukcjiWorkspace) shared.WorkspaceInstructionsVersion {
	wersja := shared.WorkspaceInstructionsVersion{
		Id: w.Identyfikator, ProjectId: w.ProjektKod, Content: w.Tresc,
		Scope: w.Poziom, CreatedAt: chwilaBazy(w.Utworzono),
	}
	if w.Odcisk != "" {
		odcisk := w.Odcisk
		wersja.ContentHash = &odcisk
	}
	if w.KluczZasiegu != "" {
		byt := w.KluczZasiegu
		wersja.ScopeId = &byt
	}
	rodzaj := w.RodzajAutora
	wersja.AuthorKind = &rodzaj
	if w.Autor != "" {
		autor := w.Autor
		wersja.AuthorId = &autor
	}
	if w.PrzywroconoZ != "" {
		zrodlo := w.PrzywroconoZ
		wersja.RestoredFromId = &zrodlo
	}
	return wersja
}

// znanyStanProjektuWorkspace pilnuje, żeby do kolumny stanu nie trafiła wartość
// spoza kontraktu — warunek kolumny odmówiłby wtedy zdaniem o schemacie, a nie
// o żądaniu.
func znanyStanProjektuWorkspace(stan shared.WorkspaceProjectStatus) bool {
	for _, znany := range shared.WartosciWorkspaceProjectStatus() {
		if znany == stan {
			return true
		}
	}
	return false
}

// przytnijWykazWorkspace stosuje granicę i pominięcie żądania do wykazu już
// złożonego. Jedna funkcja na cały moduł, bo wszystkie wykazy tego obszaru
// dostają tę samą parę pól kontraktu.
func przytnijWykazWorkspace[T any](wykaz []T, granica, pominiete *int) []T {
	poczatek := 0
	if pominiete != nil && *pominiete > 0 {
		poczatek = *pominiete
	}
	if poczatek >= len(wykaz) {
		return []T{}
	}
	wykaz = wykaz[poczatek:]
	if granica != nil && *granica > 0 && len(wykaz) > *granica {
		wykaz = wykaz[:*granica]
	}
	return wykaz
}
