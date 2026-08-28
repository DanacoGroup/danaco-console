// Odpowiedzialność pliku: odczyt repozytorium modułu Developer — stan, różnica,
// historia, gałęzie i konflikty Git Panelu, biblioteką `go-git` wkompilowaną
// w binarium rdzenia, bez uruchamiania programu `git`.
package core

import (
	"context"
	"errors"
	"strings"

	"danacoconsole/server/internal/repozytorium"
	"danacoconsole/shared"
)

// katalogRepozytorium oddaje pierwszy katalog roboczy okna — repozytorium
// Git Panelu jest jedno, ten sam katalog, w którym startują procesy okna.
func (a *adapterDevelopera) katalogRepozytorium(oknoKod string) (string, error) {
	okno, err := a.oknoDevelopera(oknoKod)
	if err != nil {
		return "", err
	}
	korzenie := a.korzenieOkna(okno)
	if len(korzenie) == 0 {
		return "", bladZadaniaDevelopera("okno " + oknoKod + " nie ma katalogu roboczego, " +
			"więc nie wiadomo, którego repozytorium czynność dotyczy")
	}
	return korzenie[0], nil
}

// bladRepozytorium przekłada błąd silnika na odmowę nazwaną kodem kontraktu.
//
// Brak repozytorium i wyjście poza obszar okna to dwie różne sprawy i dwie różne
// naprawy, więc dostają różne kody: pierwsze jest brakiem zasobu, drugie
// odmową uprawnienia.
func bladRepozytorium(err error) error {
	switch {
	case errors.Is(err, repozytorium.ErrBrakRepozytorium):
		return bladZasobuDevelopera("katalog roboczy okna nie jest repozytorium")
	case errors.Is(err, repozytorium.ErrBrakKonfliktu):
		return bladZasobuDevelopera("wskazany plik nie ma konfliktu do rozstrzygnięcia")
	case errors.Is(err, repozytorium.ErrPozaObszarem):
		return bladDostepuDevelopera("ścieżka wychodzi poza katalog roboczy okna")
	default:
		return bladWykonaniaDevelopera(err.Error())
	}
}

// StanRepozytorium oddaje stan repozytorium katalogu roboczego okna. Katalog
// bez repozytorium nie jest odmową: odpowiedź niesie `isRepository: false`.
func (a *adapterDevelopera) StanRepozytorium(_ context.Context,
	z shared.DeveloperGitStatusRequest) (shared.DeveloperGitStatusResponse, error) {

	katalog, err := a.katalogRepozytorium(z.WindowId)
	if err != nil {
		return shared.DeveloperGitStatusResponse{}, err
	}
	stan, err := repozytorium.Stan(katalog)
	if err != nil {
		return shared.DeveloperGitStatusResponse{}, bladRepozytorium(err)
	}
	return shared.DeveloperGitStatusResponse{Status: stan}, nil
}

// RoznicaRepozytorium oddaje fragmenty różnicy wraz z wykazem plików
// binarnych, których treści różnicowej nie da się pokazać tekstem.
func (a *adapterDevelopera) RoznicaRepozytorium(_ context.Context,
	z shared.DeveloperGitDiffRequest) (shared.DeveloperGitDiffResponse, error) {

	katalog, err := a.katalogRepozytorium(z.WindowId)
	if err != nil {
		return shared.DeveloperGitDiffResponse{}, err
	}
	zadanie := repozytorium.ZadanieRoznicy{Sciezki: z.Paths}
	if z.Staged != nil {
		zadanie.Przygotowane = *z.Staged
	}
	if z.FromRef != nil {
		zadanie.OdOdwolania = *z.FromRef
	}
	if z.ToRef != nil {
		zadanie.DoOdwolania = *z.ToRef
	}
	if z.ContextLines != nil {
		zadanie.WierszeKontekstu = *z.ContextLines
	}

	fragmenty, binarne, err := repozytorium.Roznica(katalog, zadanie)
	if err != nil {
		return shared.DeveloperGitDiffResponse{}, bladRepozytorium(err)
	}
	return shared.DeveloperGitDiffResponse{Hunks: fragmenty, BinaryPaths: binarne}, nil
}

// HistoriaRepozytorium oddaje zatwierdzenia spełniające warunki żądania,
// w porządku od najświeższego zatwierdzenia repozytorium.
func (a *adapterDevelopera) HistoriaRepozytorium(_ context.Context,
	z shared.DeveloperGitLogRequest) (shared.DeveloperGitLogResponse, error) {

	katalog, err := a.katalogRepozytorium(z.WindowId)
	if err != nil {
		return shared.DeveloperGitLogResponse{}, err
	}
	zadanie := repozytorium.ZadanieHistorii{}
	if z.Branch != nil {
		zadanie.Galaz = *z.Branch
	}
	if z.Path != nil {
		zadanie.Sciezka = *z.Path
	}
	if z.Author != nil {
		zadanie.Autor = *z.Author
	}
	if z.FromTime != nil {
		zadanie.OdCzasu = *z.FromTime
	}
	if z.ToTime != nil {
		zadanie.DoCzasu = *z.ToTime
	}
	if z.Limit != nil {
		zadanie.Granica = *z.Limit
	}

	zatwierdzenia, wszystkich, err := repozytorium.Historia(katalog, zadanie)
	if err != nil {
		return shared.DeveloperGitLogResponse{}, bladRepozytorium(err)
	}
	return shared.DeveloperGitLogResponse{Commits: zatwierdzenia, Total: &wszystkich}, nil
}

// GaleziRepozytorium oddaje gałęzie wraz z rozbieżnością wobec gałęzi zdalnej,
// liczoną z odwołań już pobranych, bez sięgania do sieci.
func (a *adapterDevelopera) GaleziRepozytorium(_ context.Context,
	z shared.DeveloperGitBranchListRequest) (shared.DeveloperGitBranchListResponse, error) {

	katalog, err := a.katalogRepozytorium(z.WindowId)
	if err != nil {
		return shared.DeveloperGitBranchListResponse{}, err
	}
	zeZdalnymi := z.IncludeRemote != nil && *z.IncludeRemote
	galezie, biezaca, err := repozytorium.Galezie(katalog, zeZdalnymi)
	if err != nil {
		return shared.DeveloperGitBranchListResponse{}, bladRepozytorium(err)
	}
	return shared.DeveloperGitBranchListResponse{Branches: galezie, Current: biezaca}, nil
}

// KonfliktRepozytorium rozkłada plik skonfliktowany na trzy wersje: bazową,
// lokalną i zdalną, do wyboru rozstrzygnięcia w Git Panelu.
func (a *adapterDevelopera) KonfliktRepozytorium(_ context.Context,
	z shared.DeveloperGitConflictGetRequest) (shared.DeveloperGitConflictGetResponse, error) {

	katalog, err := a.katalogRepozytorium(z.WindowId)
	if err != nil {
		return shared.DeveloperGitConflictGetResponse{}, err
	}
	if strings.TrimSpace(z.Path) == "" {
		return shared.DeveloperGitConflictGetResponse{},
			bladZadaniaDevelopera("czynność wymaga wskazania pliku z konfliktem")
	}
	konflikt, err := repozytorium.Konflikt(katalog, z.Path)
	if err != nil {
		return shared.DeveloperGitConflictGetResponse{}, bladRepozytorium(err)
	}
	return shared.DeveloperGitConflictGetResponse{Conflict: konflikt}, nil
}

// RozstrzygnijKonfliktRepozytorium zapisuje wybraną treść i przygotowuje plik.
//
// Czynność zmienia plik na dysku, więc przechodzi przez bramę trybu uprawnień:
// okno w trybie planistycznym pracuje bez zmian w systemie.
func (a *adapterDevelopera) RozstrzygnijKonfliktRepozytorium(_ context.Context,
	z shared.DeveloperGitConflictResolveRequest) (shared.DeveloperGitConflictResolveResponse, error) {

	okno, err := a.oknoDevelopera(z.WindowId)
	if err != nil {
		return shared.DeveloperGitConflictResolveResponse{}, err
	}
	if err := sprawdzZmianeSystemu(okno.TrybUprawnien, "rozstrzygnięcie konfliktu"); err != nil {
		return shared.DeveloperGitConflictResolveResponse{}, err
	}
	korzenie := a.korzenieOkna(okno)
	if len(korzenie) == 0 {
		return shared.DeveloperGitConflictResolveResponse{},
			bladZadaniaDevelopera("okno " + z.WindowId + " nie ma katalogu roboczego")
	}
	if strings.TrimSpace(z.Path) == "" {
		return shared.DeveloperGitConflictResolveResponse{},
			bladZadaniaDevelopera("czynność wymaga wskazania pliku z konfliktem")
	}
	// Rozstrzygnięcie ręczne bez treści skasowałoby plik: pustą treść trzeba
	// podać wprost.
	if z.Resolution == shared.ConflictResolutionKindManual && z.Content == nil {
		return shared.DeveloperGitConflictResolveResponse{},
			bladZadaniaDevelopera("rozstrzygnięcie ręczne wymaga treści wynikowej")
	}

	tresc := ""
	if z.Content != nil {
		tresc = *z.Content
	}
	pozostale, err := repozytorium.RozstrzygnijKonflikt(korzenie[0], z.Path, z.Resolution, tresc)
	if err != nil {
		return shared.DeveloperGitConflictResolveResponse{}, bladRepozytorium(err)
	}
	return shared.DeveloperGitConflictResolveResponse{
		Resolved: true, RemainingPaths: pozostale,
	}, nil
}
