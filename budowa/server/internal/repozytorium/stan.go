// Odpowiedzialność pakietu: czynności modułu Developer na repozytorium
// katalogu roboczego okna — stan, różnica, historia, gałęzie i konflikty.
//
// ── Biblioteka wkompilowana, nie program `git` ──────────────────────────────
// Cały pakiet stoi na `go-git` wkompilowanym w binarium rdzenia. Nie ma tu ani
// jednego uruchomienia procesu i mieć nie będzie: funkcja zależna od programu,
// którego instalka nie niesie, jest u odbiorcy odmową, a nie funkcją. Odczyt
// stanu repozytorium ma działać zawsze, bo od niego zaczyna się każda inna
// czynność okna Git Panel.
//
// ── Pakiet nie zna kontraktu komend, zna wyłącznie kształt danych ───────────
// Stąd nie wychodzi ani jedna odmowa protokołu: pakiet oddaje wynik albo błąd
// Go, a nazwanie go Operatorowi należy do adaptera. Dzięki temu ten sam silnik
// obsługuje Git Panel, margines zmian w edytorze i przegląd różnicy przed
// scaleniem, nie ucząc się o żadnym z nich.
package repozytorium

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"

	"danacoconsole/shared"
)

// ErrBrakRepozytorium mówi, że katalog roboczy nie jest repozytorium.
//
// Błąd jest osobny, bo odpowiedź kontraktu rozróżnia katalog czysty od katalogu
// bez repozytorium polem `isRepository`: bez tego rozróżnienia Git Panel
// pokazywałby „brak zmian" tam, gdzie repozytorium nigdy nie założono.
var ErrBrakRepozytorium = errors.New("katalog roboczy nie jest repozytorium")

// ErrBrakKonfliktu mówi, że wskazany plik nie jest skonfliktowany.
var ErrBrakKonfliktu = errors.New("plik nie ma konfliktu do rozstrzygnięcia")

// Otworz otwiera repozytorium katalogu roboczego.
//
// Szukanie idzie w górę drzewa, bo katalogiem roboczym okna bywa podkatalog
// repozytorium — a Operator, który otworzył `server/internal`, oczekuje stanu
// swojego repozytorium, nie komunikatu o jego braku.
func Otworz(katalog string) (*git.Repository, error) {
	if strings.TrimSpace(katalog) == "" {
		return nil, ErrBrakRepozytorium
	}
	repo, err := git.PlainOpenWithOptions(katalog, &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if errors.Is(err, git.ErrRepositoryNotExists) {
		return nil, ErrBrakRepozytorium
	}
	if err != nil {
		return nil, fmt.Errorf("otwarcie repozytorium: %w", err)
	}
	return repo, nil
}

// Korzen oddaje katalog główny repozytorium.
func Korzen(repo *git.Repository) (string, error) {
	drzewo, err := repo.Worktree()
	if err != nil {
		return "", fmt.Errorf("katalog roboczy repozytorium: %w", err)
	}
	return drzewo.Filesystem.Root(), nil
}

// stanPliku odwzorowuje kod stanu `go-git` na wartość kontraktu.
//
// Odwzorowanie jest jawne i pełne: wartość nieznana idzie na „bez zmian", bo
// stan zmyślony byłby gorszy od stanu nieokreślonego — Operator zobaczyłby
// plakietkę zmiany na pliku, którego nikt nie ruszał.
func stanPliku(kod git.StatusCode) shared.GitFileState {
	switch kod {
	case git.Unmodified:
		return shared.GitFileStateUnmodified
	case git.Modified:
		return shared.GitFileStateModified
	case git.Added:
		return shared.GitFileStateAdded
	case git.Deleted:
		return shared.GitFileStateDeleted
	case git.Renamed:
		return shared.GitFileStateRenamed
	case git.Copied:
		return shared.GitFileStateCopied
	case git.Untracked:
		return shared.GitFileStateUntracked
	case git.UpdatedButUnmerged:
		return shared.GitFileStateConflicted
	default:
		return shared.GitFileStateUnmodified
	}
}

// Stan oddaje stan repozytorium katalogu roboczego.
//
// Rozbieżność wobec gałęzi zdalnej liczy się z lokalnych referencji, bez
// sięgania do sieci: odczyt stanu ma być natychmiastowy i ma działać bez
// łączności. Liczba mówi więc o rozbieżności wobec tego, co rdzeń ostatnio
// pobrał — i tak samo działa `git status`.
func Stan(katalog string) (shared.DeveloperGitStatus, error) {
	repo, err := Otworz(katalog)
	if err != nil {
		if errors.Is(err, ErrBrakRepozytorium) {
			return shared.DeveloperGitStatus{
				Entries: []shared.GitStatusEntry{}, IsRepository: false,
			}, nil
		}
		return shared.DeveloperGitStatus{}, err
	}

	drzewo, err := repo.Worktree()
	if err != nil {
		return shared.DeveloperGitStatus{}, fmt.Errorf("katalog roboczy repozytorium: %w", err)
	}
	stanDrzewa, err := drzewo.Status()
	if err != nil {
		return shared.DeveloperGitStatus{}, fmt.Errorf("odczyt stanu repozytorium: %w", err)
	}

	wpisy := make([]shared.GitStatusEntry, 0, len(stanDrzewa))
	konflikty := false
	for sciezka, wpis := range stanDrzewa {
		if wpis.Staging == git.Unmodified && wpis.Worktree == git.Unmodified {
			continue
		}
		pozycja := shared.GitStatusEntry{
			Path:     sciezka,
			Index:    stanPliku(wpis.Staging),
			Worktree: stanPliku(wpis.Worktree),
		}
		if wpis.Extra != "" {
			poprzednia := wpis.Extra
			pozycja.RenamedFrom = &poprzednia
		}
		if pozycja.Index == shared.GitFileStateConflicted ||
			pozycja.Worktree == shared.GitFileStateConflicted {
			konflikty = true
		}
		wpisy = append(wpisy, pozycja)
	}
	// Kolejność mapy w Go jest losowa, a wykaz zmian pokazywany Operatorowi ma
	// stać w miejscu między jednym odczytem a drugim.
	sort.Slice(wpisy, func(i, j int) bool { return wpisy[i].Path < wpisy[j].Path })

	stan := shared.DeveloperGitStatus{
		Entries: wpisy, HasConflicts: konflikty, IsRepository: true,
	}

	glowa, err := repo.Head()
	if err == nil && glowa.Name().IsBranch() {
		nazwa := glowa.Name().Short()
		stan.Branch = &nazwa
		if sledzona, przed, za, znaleziona := rozbieznosc(repo, nazwa); znaleziona {
			stan.Upstream = &sledzona
			stan.Ahead = &przed
			stan.Behind = &za
		}
	}
	return stan, nil
}

// rozbieznosc liczy, o ile zatwierdzeń gałąź wyprzedza gałąź zdalną i o ile za
// nią zostaje.
//
// Brak gałęzi śledzonej nie jest błędem: gałąź lokalna bez odpowiednika zdalnego
// jest zwykłym stanem pracy, więc czynność oddaje „nie znaleziono" zamiast
// zatrzymywać odczyt stanu.
func rozbieznosc(repo *git.Repository, galaz string) (string, int, int, bool) {
	nastawy, err := repo.Config()
	if err != nil {
		return "", 0, 0, false
	}
	opis, jest := nastawy.Branches[galaz]
	if !jest || opis.Remote == "" {
		return "", 0, 0, false
	}
	nazwaZdalnej := opis.Remote + "/" + galaz
	zdalna, err := repo.Reference(
		plumbing.NewRemoteReferenceName(opis.Remote, galaz), true)
	if err != nil {
		return nazwaZdalnej, 0, 0, false
	}
	lokalna, err := repo.Reference(plumbing.NewBranchReferenceName(galaz), true)
	if err != nil {
		return nazwaZdalnej, 0, 0, false
	}

	przed, err := policzDoPrzodkow(repo, lokalna.Hash(), zdalna.Hash())
	if err != nil {
		return nazwaZdalnej, 0, 0, false
	}
	za, err := policzDoPrzodkow(repo, zdalna.Hash(), lokalna.Hash())
	if err != nil {
		return nazwaZdalnej, przed, 0, true
	}
	return nazwaZdalnej, przed, za, true
}

// policzDoPrzodkow liczy zatwierdzenia osiągalne z `od`, a nieosiągalne z `do`.
func policzDoPrzodkow(repo *git.Repository, od, doPunktu plumbing.Hash) (int, error) {
	osiagalne := map[plumbing.Hash]bool{}
	if err := przejdzPrzodkow(repo, doPunktu, osiagalne); err != nil {
		return 0, err
	}
	wlasne := map[plumbing.Hash]bool{}
	if err := przejdzPrzodkow(repo, od, wlasne); err != nil {
		return 0, err
	}
	liczba := 0
	for skrot := range wlasne {
		if !osiagalne[skrot] {
			liczba++
		}
	}
	return liczba, nil
}

// przejdzPrzodkow zbiera zatwierdzenia osiągalne ze wskazanego punktu.
func przejdzPrzodkow(repo *git.Repository, punkt plumbing.Hash,
	zebrane map[plumbing.Hash]bool) error {

	if punkt.IsZero() {
		return nil
	}
	doOdwiedzenia := []plumbing.Hash{punkt}
	for len(doOdwiedzenia) > 0 {
		biezacy := doOdwiedzenia[len(doOdwiedzenia)-1]
		doOdwiedzenia = doOdwiedzenia[:len(doOdwiedzenia)-1]
		if zebrane[biezacy] {
			continue
		}
		zebrane[biezacy] = true
		zatwierdzenie, err := repo.CommitObject(biezacy)
		if err != nil {
			// Zatwierdzenie nieobecne w magazynie obiektów kończy tę gałąź
			// przejścia zamiast całego liczenia: repozytorium pobrane płytko
			// (`shallow`) nie ma przodków i to jest jego zwykły stan.
			continue
		}
		doOdwiedzenia = append(doOdwiedzenia, zatwierdzenie.ParentHashes...)
	}
	return nil
}

// Galezie oddaje gałęzie repozytorium wraz z rozbieżnością wobec gałęzi zdalnej.
func Galezie(katalog string, zeZdalnymi bool) ([]shared.GitBranch, *string, error) {
	repo, err := Otworz(katalog)
	if err != nil {
		return nil, nil, err
	}
	var biezaca *string
	if glowa, err := repo.Head(); err == nil && glowa.Name().IsBranch() {
		nazwa := glowa.Name().Short()
		biezaca = &nazwa
	}

	galezie := make([]shared.GitBranch, 0)
	wykaz, err := repo.Branches()
	if err != nil {
		return nil, nil, fmt.Errorf("odczyt gałęzi: %w", err)
	}
	err = wykaz.ForEach(func(odwolanie *plumbing.Reference) error {
		nazwa := odwolanie.Name().Short()
		skrot := odwolanie.Hash().String()
		galaz := shared.GitBranch{
			Name:         nazwa,
			Current:      biezaca != nil && *biezaca == nazwa,
			LastCommitId: &skrot,
		}
		if sledzona, przed, za, znaleziona := rozbieznosc(repo, nazwa); znaleziona {
			galaz.Upstream = &sledzona
			galaz.Ahead = &przed
			galaz.Behind = &za
		}
		galezie = append(galezie, galaz)
		return nil
	})
	if err != nil {
		return nil, nil, fmt.Errorf("odczyt gałęzi: %w", err)
	}

	if zeZdalnymi {
		odwolania, err := repo.References()
		if err != nil {
			return nil, nil, fmt.Errorf("odczyt gałęzi zdalnych: %w", err)
		}
		err = odwolania.ForEach(func(odwolanie *plumbing.Reference) error {
			if !odwolanie.Name().IsRemote() {
				return nil
			}
			pelna := odwolanie.Name().Short()
			skrot := odwolanie.Hash().String()
			zdalne := strings.SplitN(pelna, "/", 2)
			galaz := shared.GitBranch{Name: pelna, LastCommitId: &skrot}
			if len(zdalne) == 2 {
				nazwaZdalnego := zdalne[0]
				galaz.Remote = &nazwaZdalnego
			}
			galezie = append(galezie, galaz)
			return nil
		})
		if err != nil {
			return nil, nil, fmt.Errorf("odczyt gałęzi zdalnych: %w", err)
		}
	}

	sort.Slice(galezie, func(i, j int) bool { return galezie[i].Name < galezie[j].Name })
	return galezie, biezaca, nil
}

// znacznikiKonfliktu wymieniają wiersze, którymi scalanie znaczy fragment
// sporny w pliku katalogu roboczego.
const (
	znacznikBiezacej      = "<<<<<<<"
	znacznikRozdzielajacy = "======="
	znacznikPrzodka       = "|||||||"
	znacznikPrzychodzacej = ">>>>>>>"
)

// Konflikt rozkłada plik skonfliktowany na trzy wersje.
//
// Rozbiór idzie po znacznikach w pliku katalogu roboczego, a nie po indeksie
// repozytorium: to plik na dysku widzi Operator i to jego treść ma pokazać widok
// trójstronny. Wersja bazowa bywa nieobecna, bo znacznik przodka wypisuje
// wyłącznie scalanie w stylu `diff3`.
func Konflikt(katalog, sciezka string) (shared.GitConflict, error) {
	repo, err := Otworz(katalog)
	if err != nil {
		return shared.GitConflict{}, err
	}
	korzen, err := Korzen(repo)
	if err != nil {
		return shared.GitConflict{}, err
	}
	bajty, err := os.ReadFile(filepath.Join(korzen, filepath.FromSlash(sciezka)))
	if err != nil {
		return shared.GitConflict{}, fmt.Errorf("odczyt pliku %s: %w", sciezka, err)
	}

	konflikt := shared.GitConflict{Path: sciezka}
	var biezaca, przychodzaca, przodek []string
	var przodekObecny bool
	gdzie := 0 // 0 — poza konfliktem, 1 — bieżąca, 2 — przodek, 3 — przychodząca
	znaleziony := false

	for _, wiersz := range strings.Split(string(bajty), "\n") {
		switch {
		case strings.HasPrefix(wiersz, znacznikBiezacej):
			gdzie, znaleziony = 1, true
		case strings.HasPrefix(wiersz, znacznikPrzodka) && gdzie == 1:
			gdzie, przodekObecny = 2, true
		case strings.HasPrefix(wiersz, znacznikRozdzielajacy) && gdzie != 0:
			gdzie = 3
		case strings.HasPrefix(wiersz, znacznikPrzychodzacej) && gdzie == 3:
			gdzie = 0
		default:
			switch gdzie {
			case 0:
				// Treść wspólna wchodzi do obu stron, bo widok trójstronny
				// pokazuje całe pliki, nie same fragmenty sporne.
				biezaca = append(biezaca, wiersz)
				przychodzaca = append(przychodzaca, wiersz)
				przodek = append(przodek, wiersz)
			case 1:
				biezaca = append(biezaca, wiersz)
			case 2:
				przodek = append(przodek, wiersz)
			case 3:
				przychodzaca = append(przychodzaca, wiersz)
			}
		}
	}
	if !znaleziony {
		return shared.GitConflict{}, ErrBrakKonfliktu
	}

	konflikt.Current = strings.Join(biezaca, "\n")
	konflikt.Incoming = strings.Join(przychodzaca, "\n")
	if przodekObecny {
		tresc := strings.Join(przodek, "\n")
		konflikt.Base = &tresc
	}
	return konflikt, nil
}

// RozstrzygnijKonflikt zapisuje wybraną treść i przygotowuje plik do
// zatwierdzenia.
//
// Zapis idzie razem z przygotowaniem zmiany (`add`), bo scalanie uznaje konflikt
// za rozstrzygnięty dopiero po wpisaniu pliku do indeksu. Sam zapis na dysk
// zostawiłby Operatora z plikiem poprawionym i repozytorium nadal
// skonfliktowanym.
func RozstrzygnijKonflikt(katalog, sciezka string,
	sposob shared.ConflictResolutionKind, tresc string) ([]string, error) {

	repo, err := Otworz(katalog)
	if err != nil {
		return nil, err
	}
	korzen, err := Korzen(repo)
	if err != nil {
		return nil, err
	}

	var wynik string
	switch sposob {
	case shared.ConflictResolutionKindManual:
		wynik = tresc
	case shared.ConflictResolutionKindTakeCurrent,
		shared.ConflictResolutionKindTakeIncoming,
		shared.ConflictResolutionKindTakeBoth:
		konflikt, err := Konflikt(katalog, sciezka)
		if err != nil {
			return nil, err
		}
		switch sposob {
		case shared.ConflictResolutionKindTakeCurrent:
			wynik = konflikt.Current
		case shared.ConflictResolutionKindTakeIncoming:
			wynik = konflikt.Incoming
		default:
			// Obie wersje idą jedna po drugiej w kolejności scalania: bieżąca
			// jest tą, na której Operator stoi, więc stoi wyżej.
			wynik = konflikt.Current + "\n" + konflikt.Incoming
		}
	default:
		return nil, fmt.Errorf("sposób rozstrzygnięcia %q nie jest znany", sposob)
	}

	pelna := filepath.Join(korzen, filepath.FromSlash(sciezka))
	poprzednie, err := os.Stat(pelna)
	prawa := os.FileMode(0o644)
	if err == nil {
		prawa = poprzednie.Mode().Perm()
	}
	if err := os.WriteFile(pelna, []byte(wynik), prawa); err != nil {
		return nil, fmt.Errorf("zapis pliku %s: %w", sciezka, err)
	}
	drzewo, err := repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("katalog roboczy repozytorium: %w", err)
	}
	if _, err := drzewo.Add(sciezka); err != nil {
		return nil, fmt.Errorf("przygotowanie pliku %s: %w", sciezka, err)
	}

	stan, err := Stan(katalog)
	if err != nil {
		return nil, err
	}
	pozostale := make([]string, 0)
	for _, wpis := range stan.Entries {
		if wpis.Index == shared.GitFileStateConflicted ||
			wpis.Worktree == shared.GitFileStateConflicted {
			pozostale = append(pozostale, wpis.Path)
		}
	}
	return pozostale, nil
}
