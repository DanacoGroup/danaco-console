// Obszar odczytu modułu Apps: `apps.deployment.list`, `apps.architecture.get`
// i `apps.workspace.list`. Zapis leży w `adapter_modul_aplikacje.go`
// (architektura) i `adapter_modul_aplikacje_wdrozenie.go` (warsztat, wdrożenia)
// — ten sam typ `adapterAplikacji`, osobne pliki wedle odpowiedzialności.
//
// Te trzy komendy są drogą powrotną do wierszy w bazie: po odświeżeniu okna
// klient nie ma innego sposobu, żeby odzyskać dziennik wdrożeń, architekturę
// i warsztat.
//
// Odczyt niczego nie wylicza ani nie naprawia. Wdrożenie wraca w stanie
// zapisanym przez silnik wykonania, plik warsztatu w treści zapisanej przez
// `apps.workspace.update`, architektura w kształcie złożonym tą samą funkcją
// `zloz`, którą oddaje ją `apps.architecture.define`.
//
// Brak architektury to puste pole, nie błąd: świeże okno jeszcze niczego nie
// zdefiniowało, więc `apps.architecture.get` oddaje wynik z pustym
// `architecture`, a nie odmowę `not_found`.
package core

import (
	"context"
	"errors"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// granicaWdrozenApp jest górną granicą strony wdrożeń, gdy żądanie nie podaje
// własnej. Kontrakt mówi wprost „brak bierze granicę rdzenia" — bez niej okno
// z dziennikiem liczonym w tysiącach ciągnęłoby całą historię przy każdym
// otwarciu panelu Deployment, choć pokazuje w nim ostatnie przebiegi.
const granicaWdrozenApp = 100

// WypiszWdrozenia obsługuje `apps.deployment.list`. `total` jest liczbą wdrożeń
// spełniających warunki (okno, ewentualne środowisko), a nie długością zwróconej
// strony — inaczej panel po ograniczeniu widoku twierdziłby, że historia jest
// krótsza, niż jest naprawdę. Liczbę podaje warstwa danych osobnym COUNT-em.
func (a *adapterAplikacji) WypiszWdrozenia(ctx context.Context,
	z shared.AppsDeploymentListRequest) (shared.AppsDeploymentListResponse, error) {

	okno := strings.TrimSpace(z.WindowId)
	if okno == "" {
		return shared.AppsDeploymentListResponse{}, bladWskazaniaAplikacji(
			"apps.deployment.list wymaga okna")
	}
	if z.Environment != nil {
		// Ten sam sprawdzian, co przy zapisie — wartość spoza kontraktu nie
		// zawęża niczego, więc bez odmowy okno dostałoby pustą listę i wzięło
		// ją za „nic tu nie wdrożono", zamiast dowiedzieć się o literówce.
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

// PobierzArchitekture obsługuje `apps.architecture.get`. Okno bez ani jednej
// architektury oddaje wynik z pustym polem — patrz nagłówek pliku. Kształt
// odpowiedzi składa ta sama `zloz`, którą oddaje architekturę zapis, więc
// panel dostaje po odświeżeniu dokładnie to, co widział przy definiowaniu.
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

// WypiszPlikiWarsztatu obsługuje `apps.workspace.list`. Warstwy repozytorium
// nie rozdziela — trzyma je w jednym wykazie okna posortowanym po (warstwa,
// ścieżka) — więc zawężenie do jednej warstwy robi się tutaj, na tym samym
// wykazie. `total` liczy pliki po zawężeniu, bo to one spełniają warunki
// żądania; komenda nie ma granicy strony, więc żaden plik nie wypada poza wynik.
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
		// Ten sam przekład, którym `apps.workspace.update` oddaje plik po
		// zapisie — kontrakt każe wprost, żeby odczyt niósł ten sam kształt.
		pliki = append(pliki, plikWarsztatuKontraktu(wiersz))
	}
	return shared.AppsWorkspaceListResponse{Files: pliki, Total: len(pliki)}, nil
}
