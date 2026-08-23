// Odpowiedzialność pliku: `developer.git.action` — czynności okna Git Panel
// przełożone na wiersz polecenia gita. Uruchomienie procesu i odczyt stanu
// repozytorium leżą w `adapter_modul_developer_git_wykonanie.go`.
//
// Słownik czynności jest mapą, nie drabiną warunków: nowa wartość GitActionKind
// w kontrakcie to nowa pozycja mapy. Czynność spoza mapy zostaje odrzucona.
//
// Każdy wynik niesie w polu `output` dokładny wiersz polecenia, który został
// uruchomiony, oraz stan repozytorium po czynności.
//
// `DeveloperGitActionRequest` nie niesie osobnego pola z odwołaniem do wersji,
// więc `revert` i `tag` biorą swoje wskazanie z pola `branch`.
package core

import (
	"context"
	"strings"
	"time"

	"danacoconsole/shared"
)

// czynnoscGita opisuje jedną pozycję słownika czynności repozytorium.
type czynnoscGita struct {
	// argumenty składa wiersz polecenia gita z treści żądania.
	argumenty func(z shared.DeveloperGitActionRequest) ([]string, error)
	// wymuszenie mówi, czy czynność ma postać wymuszoną. Żądanie `force` dla
	// czynności bez takiej postaci jest odrzucane, a nie po cichu pomijane.
	wymuszenie bool
	// siec znaczy czynność sięgającą do repozytorium zdalnego — dostaje dłuższą
	// granicę czasu niż czynność lokalna.
	siec bool
}

// granicaCzynnosciLokalnej i granicaCzynnosciSieciowej pilnują, żeby zawieszony
// git nie trzymał obsługiwacza komendy w nieskończoność.
const (
	granicaCzynnosciLokalnej  = 60 * time.Second
	granicaCzynnosciSieciowej = 5 * time.Minute
)

// czynnosciGita wiąże słownik kontraktu z wierszem polecenia.
var czynnosciGita = map[shared.GitActionKind]czynnoscGita{
	shared.GitActionKindStage: {argumenty: func(z shared.DeveloperGitActionRequest) ([]string, error) {
		if len(sciezkiZadania(z)) == 0 {
			return []string{"add", "--all"}, nil
		}
		return append([]string{"add", "--"}, sciezkiZadania(z)...), nil
	}},
	shared.GitActionKindUnstage: {argumenty: func(z shared.DeveloperGitActionRequest) ([]string, error) {
		if len(sciezkiZadania(z)) == 0 {
			return []string{"reset", "HEAD", "--"}, nil
		}
		return append([]string{"reset", "HEAD", "--"}, sciezkiZadania(z)...), nil
	}},
	shared.GitActionKindCommit: {argumenty: func(z shared.DeveloperGitActionRequest) ([]string, error) {
		opis, err := opisZatwierdzenia(z, true)
		if err != nil {
			return nil, err
		}
		polecenie := []string{"commit", "-m", opis}
		if sciezki := sciezkiZadania(z); len(sciezki) > 0 {
			polecenie = append(append(polecenie, "--"), sciezki...)
		}
		return polecenie, nil
	}},
	shared.GitActionKindAmend: {argumenty: func(z shared.DeveloperGitActionRequest) ([]string, error) {
		opis, err := opisZatwierdzenia(z, false)
		if err != nil {
			return nil, err
		}
		if opis == "" {
			// Bez nowego opisu poprawiamy wyłącznie treść zatwierdzenia;
			// `--no-edit` powstrzymuje gita przed otwarciem edytora, którego
			// w procesie rdzenia nie ma komu obsłużyć.
			return []string{"commit", "--amend", "--no-edit"}, nil
		}
		return []string{"commit", "--amend", "-m", opis}, nil
	}},
	shared.GitActionKindRevert: {argumenty: func(z shared.DeveloperGitActionRequest) ([]string, error) {
		wersja, err := wskazanieWersji(z, "odwrócenie zatwierdzenia")
		if err != nil {
			return nil, err
		}
		return []string{"revert", "--no-edit", wersja}, nil
	}},
	shared.GitActionKindCheckout: {wymuszenie: true, argumenty: func(z shared.DeveloperGitActionRequest) ([]string, error) {
		wersja, err := wskazanieWersji(z, "przełączenie gałęzi")
		if err != nil {
			return nil, err
		}
		polecenie := []string{"checkout"}
		if z.Force != nil && *z.Force {
			polecenie = append(polecenie, "--force")
		}
		return append(polecenie, wersja), nil
	}},
	shared.GitActionKindMerge: {argumenty: func(z shared.DeveloperGitActionRequest) ([]string, error) {
		wersja, err := wskazanieWersji(z, "scalenie gałęzi")
		if err != nil {
			return nil, err
		}
		return []string{"merge", "--no-edit", wersja}, nil
	}},
	shared.GitActionKindRebase: {argumenty: func(z shared.DeveloperGitActionRequest) ([]string, error) {
		wersja, err := wskazanieWersji(z, "przestawienie gałęzi")
		if err != nil {
			return nil, err
		}
		return []string{"rebase", wersja}, nil
	}},
	shared.GitActionKindTag: {argumenty: func(z shared.DeveloperGitActionRequest) ([]string, error) {
		nazwa, err := wskazanieWersji(z, "nadanie etykiety")
		if err != nil {
			return nil, err
		}
		if opis := strings.TrimSpace(wartoscTekstu(z.Message)); opis != "" {
			return []string{"tag", "-a", nazwa, "-m", opis}, nil
		}
		return []string{"tag", nazwa}, nil
	}},
	shared.GitActionKindFetch: {siec: true, argumenty: func(z shared.DeveloperGitActionRequest) ([]string, error) {
		return dolaczZdalne([]string{"fetch"}, z), nil
	}},
	shared.GitActionKindPull: {siec: true, argumenty: func(z shared.DeveloperGitActionRequest) ([]string, error) {
		return dolaczZdalne([]string{"pull"}, z), nil
	}},
	shared.GitActionKindPush: {siec: true, wymuszenie: true, argumenty: func(z shared.DeveloperGitActionRequest) ([]string, error) {
		polecenie := []string{"push"}
		if z.Force != nil && *z.Force {
			// `--force-with-lease` zamiast `--force`: wymuszenie nadpisuje stan
			// zdalny znany lokalnie, a nie zmiany dosłane w międzyczasie.
			polecenie = append(polecenie, "--force-with-lease")
		}
		return dolaczZdalne(polecenie, z), nil
	}},
	shared.GitActionKindStash: {argumenty: func(z shared.DeveloperGitActionRequest) ([]string, error) {
		polecenie := []string{"stash", "push"}
		if opis := strings.TrimSpace(wartoscTekstu(z.Message)); opis != "" {
			polecenie = append(polecenie, "-m", opis)
		}
		if sciezki := sciezkiZadania(z); len(sciezki) > 0 {
			polecenie = append(append(polecenie, "--"), sciezki...)
		}
		return polecenie, nil
	}},
	shared.GitActionKindStashPop: {argumenty: func(shared.DeveloperGitActionRequest) ([]string, error) {
		return []string{"stash", "pop"}, nil
	}},
}

// CzynnoscRepozytorium obsługuje `developer.git.action`.
func (a *adapterDevelopera) CzynnoscRepozytorium(ctx context.Context,
	z shared.DeveloperGitActionRequest) (shared.DeveloperGitActionResponse, error) {

	opis, jest := czynnosciGita[z.Action]
	if !jest {
		return shared.DeveloperGitActionResponse{}, bladZadaniaDevelopera(
			"czynność " + string(z.Action) + " nie należy do słownika kontraktu")
	}
	if z.Force != nil && *z.Force && !opis.wymuszenie {
		return shared.DeveloperGitActionResponse{}, bladZadaniaDevelopera(
			"czynność " + string(z.Action) + " nie ma postaci wymuszonej")
	}
	okno, err := a.oknoDevelopera(z.WindowId)
	if err != nil {
		return shared.DeveloperGitActionResponse{}, err
	}
	if err := sprawdzZmianeSystemu(okno.TrybUprawnien,
		"czynność repozytorium "+string(z.Action)); err != nil {
		return shared.DeveloperGitActionResponse{}, err
	}
	argumenty, err := opis.argumenty(z)
	if err != nil {
		return shared.DeveloperGitActionResponse{}, err
	}

	granica := granicaCzynnosciLokalnej
	if opis.siec {
		granica = granicaCzynnosciSieciowej
	}
	wynik, err := a.wykonajGit(ctx, okno, z.Action, argumenty, granica)
	if err != nil {
		return shared.DeveloperGitActionResponse{}, err
	}
	return shared.DeveloperGitActionResponse{Result: wynik}, nil
}

// sciezkiZadania odsiewa puste wskazania ścieżek.
func sciezkiZadania(z shared.DeveloperGitActionRequest) []string {
	sciezki := make([]string, 0, len(z.Paths))
	for _, sciezka := range z.Paths {
		if tresc := strings.TrimSpace(sciezka); tresc != "" {
			sciezki = append(sciezki, tresc)
		}
	}
	return sciezki
}

// opisZatwierdzenia czyta treść zatwierdzenia; `wymagany` odmawia pustej.
func opisZatwierdzenia(z shared.DeveloperGitActionRequest, wymagany bool) (string, error) {
	opis := strings.TrimSpace(wartoscTekstu(z.Message))
	if opis == "" && wymagany {
		return "", bladZadaniaDevelopera("zatwierdzenie bez opisu nie powstanie")
	}
	return opis, nil
}

// wskazanieWersji czyta gałąź, etykietę albo zatwierdzenie z pola `branch`.
func wskazanieWersji(z shared.DeveloperGitActionRequest, czynnosc string) (string, error) {
	wersja := strings.TrimSpace(wartoscTekstu(z.Branch))
	if wersja == "" {
		return "", bladZadaniaDevelopera(czynnosc + " wymaga wskazania gałęzi, etykiety albo zatwierdzenia")
	}
	if strings.HasPrefix(wersja, "-") {
		// Wskazanie zaczynające się od myślnika weszłoby do wiersza polecenia
		// jako przełącznik gita, a nie jako wersja.
		return "", bladZadaniaDevelopera("wskazanie " + wersja + " nie jest nazwą gałęzi ani zatwierdzenia")
	}
	return wersja, nil
}

// dolaczZdalne dokłada repozytorium zdalne i gałąź, gdy żądanie je podaje.
func dolaczZdalne(polecenie []string, z shared.DeveloperGitActionRequest) []string {
	zdalne := strings.TrimSpace(wartoscTekstu(z.Remote))
	galaz := strings.TrimSpace(wartoscTekstu(z.Branch))
	if zdalne == "" {
		return polecenie
	}
	polecenie = append(polecenie, zdalne)
	if galaz != "" {
		polecenie = append(polecenie, galaz)
	}
	return polecenie
}
