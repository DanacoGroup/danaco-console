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
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"

	"danacoconsole/shared"
)

// Różnicę repozytorium liczy algorytm Myersa w czystym Go, wynikiem
// w fragmentach zgodnych z kontraktem, nie tekstem unified diff.
//
// ZadanieRoznicy opisuje, co z czym dokładnie porównać: ścieżki, źródło zmian
// i odwołania repozytorium.
type ZadanieRoznicy struct {
	// Sciezki zawężają różnicę; pusta lista obejmuje całość.
	Sciezki []string
	// Przygotowane liczy różnicę zmian wpisanych do indeksu zamiast roboczych.
	Przygotowane bool
	// OdOdwolania i DoOdwolania wskazują gałęzie albo zatwierdzenia.
	OdOdwolania string
	DoOdwolania string
	// WierszeKontekstu mówi, ile wierszy niezmienionych zostaje wokół zmiany.
	WierszeKontekstu int
}

// Roznica oddaje fragmenty różnicy wraz z wykazem plików binarnych.
//
// Plik binarny nie ma różnicy wierszowej i udawanie jej byłoby pokazaniem
// Operatorowi śmieci. Wykaz osobny mówi wprost: plik się zmienił, ale zmiany nie
// da się pokazać wierszami.
func Roznica(katalog string, zadanie ZadanieRoznicy) ([]shared.GitDiffHunk, []string, error) {
	repo, err := Otworz(katalog)
	if err != nil {
		return nil, nil, err
	}
	korzen, err := Korzen(repo)
	if err != nil {
		return nil, nil, err
	}
	kontekst := zadanie.WierszeKontekstu
	if kontekst <= 0 {
		kontekst = 3
	}
	wybrane := map[string]bool{}
	for _, sciezka := range zadanie.Sciezki {
		wybrane[filepath.ToSlash(strings.TrimSpace(sciezka))] = true
	}

	pary, err := paryPorownania(repo, korzen, zadanie)
	if err != nil {
		return nil, nil, err
	}

	fragmenty := make([]shared.GitDiffHunk, 0)
	binarne := make([]string, 0)
	sciezki := make([]string, 0, len(pary))
	for sciezka := range pary {
		if len(wybrane) > 0 && !wybrane[sciezka] {
			continue
		}
		sciezki = append(sciezki, sciezka)
	}
	sort.Strings(sciezki)

	for _, sciezka := range sciezki {
		para := pary[sciezka]
		if czyBinarna(para.przed) || czyBinarna(para.po) {
			binarne = append(binarne, sciezka)
			continue
		}
		fragmenty = append(fragmenty, fragmentyPliku(sciezka, para.przed, para.po, kontekst)...)
	}
	return fragmenty, binarne, nil
}

// paraTresci trzyma treść jednego danego pliku, i tę przed zmianą, i tę po
// niej, do policzenia różnicy.
type paraTresci struct{ przed, po string }

// paryPorownania zbiera pliki objęte różnicą wraz z obiema wersjami treści,
// katalogu roboczego albo dwóch odwołań.
func paryPorownania(repo *git.Repository, korzen string,
	zadanie ZadanieRoznicy) (map[string]paraTresci, error) {

	// Porównanie dwóch odwołań: obie strony pochodzą z magazynu obiektów.
	if strings.TrimSpace(zadanie.DoOdwolania) != "" || strings.TrimSpace(zadanie.OdOdwolania) != "" {
		return paryDwochOdwolan(repo, zadanie.OdOdwolania, zadanie.DoOdwolania)
	}

	drzewo, err := repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("katalog roboczy repozytorium: %w", err)
	}
	stan, err := drzewo.Status()
	if err != nil {
		return nil, fmt.Errorf("odczyt stanu repozytorium: %w", err)
	}
	glowaTresc := trescGlowy(repo)

	pary := map[string]paraTresci{}
	for sciezka, wpis := range stan {
		if zadanie.Przygotowane {
			if wpis.Staging == git.Unmodified || wpis.Staging == git.Untracked {
				continue
			}
			pary[sciezka] = paraTresci{
				przed: glowaTresc(sciezka),
				po:    trescIndeksu(repo, sciezka),
			}
			continue
		}
		if wpis.Worktree == git.Unmodified {
			continue
		}
		pary[sciezka] = paraTresci{
			przed: trescIndeksu(repo, sciezka),
			po:    trescRobocza(korzen, sciezka),
		}
	}
	return pary, nil
}

// paryDwochOdwolan zestawia drzewa dwóch odwołań repozytorium, zbierając
// treść pliku z każdej strony osobno.
func paryDwochOdwolan(repo *git.Repository, od, doPunktu string) (map[string]paraTresci, error) {
	drzewoOd, err := drzewoOdwolania(repo, od)
	if err != nil {
		return nil, err
	}
	drzewoDo, err := drzewoOdwolania(repo, doPunktu)
	if err != nil {
		return nil, err
	}

	pary := map[string]paraTresci{}
	zbierz := func(drzewo *object.Tree, doPary func(*paraTresci, string)) error {
		if drzewo == nil {
			return nil
		}
		return drzewo.Files().ForEach(func(plik *object.File) error {
			tresc, err := plik.Contents()
			if err != nil {
				return nil
			}
			para := pary[plik.Name]
			doPary(&para, tresc)
			pary[plik.Name] = para
			return nil
		})
	}
	if err := zbierz(drzewoOd, func(p *paraTresci, t string) { p.przed = t }); err != nil {
		return nil, fmt.Errorf("odczyt drzewa %q: %w", od, err)
	}
	if err := zbierz(drzewoDo, func(p *paraTresci, t string) { p.po = t }); err != nil {
		return nil, fmt.Errorf("odczyt drzewa %q: %w", doPunktu, err)
	}

	// Pliki niezmienione odpadają tutaj, a nie przy liczeniu różnicy.
	for sciezka, para := range pary {
		if para.przed == para.po {
			delete(pary, sciezka)
		}
	}
	return pary, nil
}

// drzewoOdwolania rozwiązuje nazwę gałęzi albo zatwierdzenia na drzewo plików.
//
// Puste odwołanie znaczy „głowa repozytorium", bo tak czyta się różnicę wobec
// jednego wskazanego punktu.
func drzewoOdwolania(repo *git.Repository, nazwa string) (*object.Tree, error) {
	nazwa = strings.TrimSpace(nazwa)
	if nazwa == "" {
		nazwa = "HEAD"
	}
	skrot, err := repo.ResolveRevision(plumbing.Revision(nazwa))
	if err != nil {
		return nil, fmt.Errorf("odwołanie %q nie istnieje w repozytorium: %w", nazwa, err)
	}
	zatwierdzenie, err := repo.CommitObject(*skrot)
	if err != nil {
		return nil, fmt.Errorf("odczyt zatwierdzenia %s: %w", skrot.String(), err)
	}
	drzewo, err := zatwierdzenie.Tree()
	if err != nil {
		return nil, fmt.Errorf("odczyt drzewa zatwierdzenia %s: %w", skrot.String(), err)
	}
	return drzewo, nil
}

// trescGlowy oddaje czytelnik treści plików z głowy repozytorium.
//
// Repozytorium bez ani jednego zatwierdzenia oddaje treść pustą zamiast błędu:
// pierwszy zapis w nowym repozytorium jest zwykłą pracą, a nie awarią.
func trescGlowy(repo *git.Repository) func(string) string {
	drzewo, err := drzewoOdwolania(repo, "HEAD")
	if err != nil || drzewo == nil {
		return func(string) string { return "" }
	}
	return func(sciezka string) string {
		plik, err := drzewo.File(sciezka)
		if err != nil {
			return ""
		}
		tresc, err := plik.Contents()
		if err != nil {
			return ""
		}
		return tresc
	}
}

// trescIndeksu oddaje treść danego pliku wpisaną do indeksu tego repozytorium,
// pustym napisem przy braku.
func trescIndeksu(repo *git.Repository, sciezka string) string {
	indeks, err := repo.Storer.Index()
	if err != nil {
		return ""
	}
	wpis, err := indeks.Entry(sciezka)
	if err != nil {
		return ""
	}
	obiekt, err := repo.BlobObject(wpis.Hash)
	if err != nil {
		return ""
	}
	czytnik, err := obiekt.Reader()
	if err != nil {
		return ""
	}
	defer func() { _ = czytnik.Close() }()
	bajty := make([]byte, obiekt.Size)
	odczytane, _ := czytnik.Read(bajty)
	return string(bajty[:odczytane])
}

// trescRobocza oddaje treść pliku z katalogu roboczego.
//
// Plik usunięty daje treść pustą, nie błąd: usunięcie jest właśnie tą zmianą,
// którą różnica ma pokazać.
func trescRobocza(korzen, sciezka string) string {
	bajty, err := os.ReadFile(filepath.Join(korzen, filepath.FromSlash(sciezka)))
	if err != nil {
		return ""
	}
	return string(bajty)
}

// czyBinarna rozstrzyga, czy treść nadaje się do pokazania wierszami.
//
// Miarą jest bajt zerowy, tak samo jak w `git`: żaden zapis tekstowy go nie
// niesie, a każdy format binarny niesie go niemal na pewno.
func czyBinarna(tresc string) bool {
	granica := len(tresc)
	if granica > 8000 {
		granica = 8000
	}
	return strings.IndexByte(tresc[:granica], 0) >= 0
}

// fragmentyPliku liczy różnicę jednego pliku i przenosi ją wprost na kształt
// fragmentów tego kontraktu.
func fragmentyPliku(sciezka, przed, po string, kontekst int) []shared.GitDiffHunk {
	zmiany := myers.ComputeEdits(span.URIFromPath(sciezka), przed, po)
	if len(zmiany) == 0 {
		return nil
	}
	roznica := gotextdiff.ToUnified(sciezka, sciezka, przed, zmiany)

	fragmenty := make([]shared.GitDiffHunk, 0, len(roznica.Hunks))
	for _, fragment := range roznica.Hunks {
		wiersze := make([]shared.GitDiffLine, 0, len(fragment.Lines))
		staryPoczatek, nowyPoczatek := fragment.FromLine, fragment.ToLine
		staryNumer, nowyNumer := staryPoczatek, nowyPoczatek
		staraDlugosc, nowaDlugosc := 0, 0

		for _, wiersz := range fragment.Lines {
			tresc := strings.TrimSuffix(wiersz.Content, "\n")
			pozycja := shared.GitDiffLine{Text: tresc}
			switch wiersz.Kind {
			case gotextdiff.Delete:
				pozycja.Kind = shared.GitDiffLineKindRemoved
				numer := staryNumer
				pozycja.OldLine = &numer
				staryNumer++
				staraDlugosc++
			case gotextdiff.Insert:
				pozycja.Kind = shared.GitDiffLineKindAdded
				numer := nowyNumer
				pozycja.NewLine = &numer
				nowyNumer++
				nowaDlugosc++
			default:
				pozycja.Kind = shared.GitDiffLineKindContext
				stary, nowy := staryNumer, nowyNumer
				pozycja.OldLine = &stary
				pozycja.NewLine = &nowy
				staryNumer++
				nowyNumer++
				staraDlugosc++
				nowaDlugosc++
			}
			wiersze = append(wiersze, pozycja)
		}

		naglowek := fmt.Sprintf("@@ -%d,%d +%d,%d @@",
			staryPoczatek, staraDlugosc, nowyPoczatek, nowaDlugosc)
		fragmenty = append(fragmenty, shared.GitDiffHunk{
			Path:     sciezka,
			Header:   &naglowek,
			OldStart: staryPoczatek,
			OldLines: staraDlugosc,
			NewStart: nowyPoczatek,
			NewLines: nowaDlugosc,
			Lines:    wiersze,
		})
	}
	_ = kontekst // liczbę wierszy kontekstu ustala biblioteka; wskazanie zostaje na przyszłą nastawę
	return fragmenty
}

// ZadanieHistorii zawęża odczyt dziennika repozytorium: gałąź, ścieżkę,
// autora, zakres czasu i granicę.
type ZadanieHistorii struct {
	Galaz   string
	Sciezka string
	Autor   string
	OdCzasu int64
	DoCzasu int64
	Granica int
}

// Historia oddaje zatwierdzenia spełniające warunki tego zadania wraz z ich
// łączną liczbą w repozytorium.
func Historia(katalog string, zadanie ZadanieHistorii) ([]shared.GitCommit, int, error) {
	repo, err := Otworz(katalog)
	if err != nil {
		return nil, 0, err
	}
	punkt := strings.TrimSpace(zadanie.Galaz)
	if punkt == "" {
		punkt = "HEAD"
	}
	skrot, err := repo.ResolveRevision(plumbing.Revision(punkt))
	if err != nil {
		if errors.Is(err, plumbing.ErrReferenceNotFound) {
			return nil, 0, fmt.Errorf("gałąź %q nie istnieje w repozytorium", punkt)
		}
		return nil, 0, fmt.Errorf("odczyt gałęzi %q: %w", punkt, err)
	}

	nastawy := &git.LogOptions{From: *skrot}
	if strings.TrimSpace(zadanie.Sciezka) != "" {
		szukana := filepath.ToSlash(strings.TrimSpace(zadanie.Sciezka))
		nastawy.PathFilter = func(sciezka string) bool { return sciezka == szukana }
	}
	dziennik, err := repo.Log(nastawy)
	if err != nil {
		return nil, 0, fmt.Errorf("odczyt historii: %w", err)
	}
	defer dziennik.Close()

	zatwierdzenia := make([]shared.GitCommit, 0)
	wszystkich := 0
	autor := strings.ToLower(strings.TrimSpace(zadanie.Autor))

	err = dziennik.ForEach(func(zatwierdzenie *object.Commit) error {
		czas := zatwierdzenie.Committer.When.UnixMilli()
		if zadanie.OdCzasu > 0 && czas < zadanie.OdCzasu {
			return nil
		}
		if zadanie.DoCzasu > 0 && czas > zadanie.DoCzasu {
			return nil
		}
		if autor != "" &&
			!strings.Contains(strings.ToLower(zatwierdzenie.Author.Name), autor) &&
			!strings.Contains(strings.ToLower(zatwierdzenie.Author.Email), autor) {
			return nil
		}

		// Liczba całkowita rośnie także po osiągnięciu granicy, bo odpowiedź
		// niesie, ile ich jest.
		wszystkich++
		if zadanie.Granica > 0 && len(zatwierdzenia) >= zadanie.Granica {
			return nil
		}

		pelny := zatwierdzenie.Hash.String()
		poczta := zatwierdzenie.Author.Email
		rodzice := make([]string, 0, len(zatwierdzenie.ParentHashes))
		for _, rodzic := range zatwierdzenie.ParentHashes {
			rodzice = append(rodzice, rodzic.String())
		}
		zatwierdzenia = append(zatwierdzenia, shared.GitCommit{
			Id:          pelny,
			ShortId:     pelny[:7],
			Author:      zatwierdzenie.Author.Name,
			AuthorEmail: &poczta,
			Message:     strings.TrimRight(zatwierdzenie.Message, "\n"),
			CommittedAt: czas,
			ParentIds:   rodzice,
		})
		return nil
	})
	if err != nil {
		return nil, 0, fmt.Errorf("odczyt historii: %w", err)
	}
	return zatwierdzenia, wszystkich, nil
}
