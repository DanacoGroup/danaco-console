// Moduł Apps — dzienniki i artefakty oddające to, co praca naprawdę zapisała:
// `apps.service.log.read`, `apps.deployment.log.read`, `apps.artifact.list`.
package core

import (
	"context"
	"strconv"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// granicaDziennikaApp jest górną granicą strony dziennika przy braku wskazania w treści żądania tego okna.
const granicaDziennikaApp = 500

// OdczytajDziennikUslugi obsługuje `apps.service.log.read` i zwraca wiersze dziennika usługi tego okna.
func (a *adapterAplikacji) OdczytajDziennikUslugi(ctx context.Context,
	z shared.AppsServiceLogReadRequest) (shared.AppsServiceLogReadResponse, error) {

	okno, err := oknoAplikacji(z.WindowId, "apps.service.log.read")
	if err != nil {
		return shared.AppsServiceLogReadResponse{}, err
	}
	komponent := strings.TrimSpace(wartoscTekstu(z.ComponentId))
	od := int64(0)
	if z.Since != nil {
		od = *z.Since
	}
	granica := granicaDziennikaApp
	if z.Limit != nil && *z.Limit > 0 {
		granica = *z.Limit
	}

	wiersze, razem, err := a.repozytorium.DziennikUslugiApp(ctx, okno, komponent, od, granica)
	if err != nil {
		return shared.AppsServiceLogReadResponse{}, bladAplikacji(err)
	}
	return shared.AppsServiceLogReadResponse{
		Lines: wierszeDziennikaApp(wiersze), Total: razem, Streaming: false,
	}, nil
}

// OdczytajDziennikWdrozenia obsługuje `apps.deployment.log.read`; wdrożenie spoza okna
// żądania jest odmową, bo dziennik cudzego przebiegu nie należy do tego panelu.
func (a *adapterAplikacji) OdczytajDziennikWdrozenia(ctx context.Context,
	z shared.AppsDeploymentLogReadRequest) (shared.AppsDeploymentLogReadResponse, error) {

	okno, err := oknoAplikacji(z.WindowId, "apps.deployment.log.read")
	if err != nil {
		return shared.AppsDeploymentLogReadResponse{}, err
	}
	kod := strings.TrimSpace(z.DeploymentId)
	if kod == "" {
		return shared.AppsDeploymentLogReadResponse{}, bladWskazaniaAplikacji(
			"apps.deployment.log.read wymaga wdrożenia")
	}
	wdrozenie, err := a.repozytorium.Wdrozenie(ctx, kod)
	if err != nil {
		return shared.AppsDeploymentLogReadResponse{}, bladNieznanegoBytuApp("wdrożenie", kod, err)
	}
	if wdrozenie.OknoKod != okno {
		return shared.AppsDeploymentLogReadResponse{}, bladWskazaniaAplikacji(
			"wdrożenie " + kod + " należy do okna " + wdrozenie.OknoKod +
				", a żądanie przyszło z okna " + okno)
	}

	od := int64(0)
	if z.Since != nil {
		od = *z.Since
	}
	granica := granicaDziennikaApp
	if z.Limit != nil && *z.Limit > 0 {
		granica = *z.Limit
	}

	wiersze, razem, err := a.repozytorium.DziennikWdrozeniaApp(ctx, kod, od, granica)
	if err != nil {
		return shared.AppsDeploymentLogReadResponse{}, bladAplikacji(err)
	}
	return shared.AppsDeploymentLogReadResponse{
		Lines: wierszeDziennikaApp(wiersze), Total: razem, Streaming: false,
	}, nil
}

// WypiszArtefakty obsługuje `apps.artifact.list`; każdy wiersz wskazuje plik w magazynie
// treści rdzenia wraz z jego rozmiarem i sumą kontrolną.
func (a *adapterAplikacji) WypiszArtefakty(ctx context.Context,
	z shared.AppsArtifactListRequest) (shared.AppsArtifactListResponse, error) {

	okno, err := oknoAplikacji(z.WindowId, "apps.artifact.list")
	if err != nil {
		return shared.AppsArtifactListResponse{}, err
	}
	wdrozenie := strings.TrimSpace(wartoscTekstu(z.DeploymentId))
	wiersze, err := a.repozytorium.ArtefaktyApp(ctx, okno, wdrozenie)
	if err != nil {
		return shared.AppsArtifactListResponse{}, bladAplikacji(err)
	}
	artefakty := make([]shared.AppArtifact, 0, len(wiersze))
	for _, wiersz := range wiersze {
		artefakty = append(artefakty, shared.AppArtifact{
			Id: wiersz.Kod, WindowId: wiersz.Okno, DeploymentId: wiersz.WdrozenieKod,
			Kind: shared.AppArtifactKind(wiersz.Rodzaj), Path: wiersz.Sciezka,
			SizeBytes: wiersz.Rozmiar, ChecksumSha256: wiersz.SumaKontrolna,
			CreatedAt: chwilaBazy(wiersz.Utworzono),
		})
	}
	return shared.AppsArtifactListResponse{Artifacts: artefakty, Total: len(artefakty)}, nil
}

// wierszeDziennikaApp składa wiersze dziennika w postać, którą niesie kontrakt: wykaz napisów
// ze znacznikiem czasu w treści, bo kontrakt nie ma osobnego pola na chwilę.
func wierszeDziennikaApp(wiersze []dane.WierszDziennikaApp) []string {
	linie := make([]string, 0, len(wiersze))
	for _, wiersz := range wiersze {
		linie = append(linie, strconv.FormatInt(wiersz.Chwila, 10)+" "+wiersz.Tresc)
	}
	return linie
}
