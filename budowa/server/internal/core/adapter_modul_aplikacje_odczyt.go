// Obszar odczytu modułu Apps: `apps.deployment.list`, `apps.architecture.get` i
// `apps.workspace.list`, drogi powrotne do wierszy bazy po odświeżeniu okna klienta.
package core

import (
	"context"
	"errors"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// granicaWdrozenApp jest górną granicą strony wdrożeń, gdy żądanie nie podaje własnej wartości granicy.
const granicaWdrozenApp = 100

// WypiszWdrozenia obsługuje `apps.deployment.list`; pole całkowitej liczby jest liczbą
// wdrożeń spełniających warunki, a nie długością zwróconej strony wyniku.
func (a *adapterAplikacji) WypiszWdrozenia(ctx context.Context,
	z shared.AppsDeploymentListRequest) (shared.AppsDeploymentListResponse, error) {

	okno := strings.TrimSpace(z.WindowId)
	if okno == "" {
		return shared.AppsDeploymentListResponse{}, bladWskazaniaAplikacji(
			"apps.deployment.list wymaga okna")
	}
	if z.Environment != nil {
		// Ten sam sprawdzian, co przy zapisie: wartość spoza kontraktu odmawia, zamiast dać pustą listę.
		if err := sprawdzSrodowiskoWdrozenia(*z.Environment); err != nil {
			return shared.AppsDeploymentListResponse{}, err
		}
	}

	granica := granicaWdrozenApp
	if z.Limit != nil && *z.Limit > 0 {
		granica = *z.Limit
	}

	wiersze, razem, err := a.repozytorium.Wdrozenia(ctx, okno, z.Environment, granica)
	if err != nil {
		return shared.AppsDeploymentListResponse{}, bladAplikacji(err)
	}

	wdrozenia := make([]shared.AppDeployment, 0, len(wiersze))
	for _, wiersz := range wiersze {
		wdrozenia = append(wdrozenia, wdrozenieKontraktu(wiersz))
	}
	return shared.AppsDeploymentListResponse{Deployments: wdrozenia, Total: razem}, nil
}

// PobierzArchitekture obsługuje `apps.architecture.get`; okno bez ani jednej architektury
// oddaje wynik z pustym polem, nie odmowę, tym samym kształtem co zapis.
func (a *adapterAplikacji) PobierzArchitekture(ctx context.Context,
	z shared.AppsArchitectureGetRequest) (shared.AppsArchitectureGetResponse, error) {

	okno := strings.TrimSpace(z.WindowId)
	if okno == "" {
		return shared.AppsArchitectureGetResponse{}, bladWskazaniaAplikacji(
			"apps.architecture.get wymaga okna")
	}

	wiersz, err := a.repozytorium.ArchitekturaOkna(ctx, okno)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.AppsArchitectureGetResponse{}, nil
	}
	if err != nil {
		return shared.AppsArchitectureGetResponse{}, bladAplikacji(err)
	}

	architektura, err := a.zloz(ctx, wiersz)
	if err != nil {
		return shared.AppsArchitectureGetResponse{}, err
	}
	return shared.AppsArchitectureGetResponse{Architecture: &architektura}, nil
}

// WypiszPlikiWarsztatu obsługuje `apps.workspace.list`; repozytorium trzyma warstwy w jednym
// wykazie okna, więc zawężenie do jednej warstwy robi się tutaj, na tym samym wykazie.
func (a *adapterAplikacji) WypiszPlikiWarsztatu(ctx context.Context,
	z shared.AppsWorkspaceListRequest) (shared.AppsWorkspaceListResponse, error) {

	okno := strings.TrimSpace(z.WindowId)
	if okno == "" {
		return shared.AppsWorkspaceListResponse{}, bladWskazaniaAplikacji(
			"apps.workspace.list wymaga okna")
	}
	if z.Layer != nil {
		if err := sprawdzWarstweWarsztatu(*z.Layer); err != nil {
			return shared.AppsWorkspaceListResponse{}, err
		}
	}

	wiersze, err := a.repozytorium.PlikiWarsztatu(ctx, okno)
	if err != nil {
		return shared.AppsWorkspaceListResponse{}, bladAplikacji(err)
	}

	pliki := make([]shared.DeveloperFile, 0, len(wiersze))
	for _, wiersz := range wiersze {
		if z.Layer != nil && wiersz.Warstwa != *z.Layer {
			continue
		}
		// Ten sam przekład, którym zapis oddaje plik po zmianie.
		pliki = append(pliki, plikWarsztatuKontraktu(wiersz))
	}
	return shared.AppsWorkspaceListResponse{Files: pliki, Total: len(pliki)}, nil
}
