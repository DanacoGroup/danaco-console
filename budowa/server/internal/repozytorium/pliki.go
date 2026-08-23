package repozytorium

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Czynności plikowe drzewa projektu.
//
// ── Ścieżka nie wychodzi poza katalog roboczy ───────────────────────────────
// Każda czynność sprowadza wskazanie Operatora do ścieżki wewnątrz katalogu
// roboczego okna i odmawia, gdy wskazanie z niego wychodzi. Bez tego `../../`
// w nazwie pliku byłoby drogą do dowolnego miejsca na dysku serwera — a okno
// modułu Developer stoi po to, żeby pracować w JEDNYM repozytorium.
//
// Sprawdzenie idzie po ścieżce rozwiniętej z dowiązań, nie po samym napisie:
// dowiązanie symboliczne wskazujące poza katalog jest tym samym wyjściem, tylko
// zapisanym inaczej.

// ErrPozaObszarem mówi, że wskazanie wychodzi poza katalog roboczy okna.
var ErrPozaObszarem = errors.New("ścieżka wychodzi poza katalog roboczy okna")

// ErrJuzIstnieje mówi, że w miejscu docelowym coś już leży.
var ErrJuzIstnieje = errors.New("w miejscu docelowym już coś leży")

// wObszarze sprowadza wskazanie względne do ścieżki bezwzględnej w katalogu
// roboczym i pilnuje, żeby z niego nie wyszła.
func wObszarze(korzen, wskazanie string) (string, error) {
	korzenPelny, err := filepath.Abs(korzen)
	if err != nil {
		return "", fmt.Errorf("katalog roboczy %s: %w", korzen, err)
	}
	if rozwiniety, err := filepath.EvalSymlinks(korzenPelny); err == nil {
		korzenPelny = rozwiniety
	}

	oczyszczone := filepath.Clean(filepath.FromSlash(strings.TrimSpace(wskazanie)))
	if oczyszczone == "." || oczyszczone == string(filepath.Separator) {
		return "", ErrPozaObszarem
	}
	if filepath.IsAbs(oczyszczone) {
		return "", ErrPozaObszarem
	}
	pelna := filepath.Join(korzenPelny, oczyszczone)

	// Ścieżka rozwija się z dowiązań tylko wtedy, gdy istnieje. Przy zakładaniu
	// pliku sprawdza się katalog nadrzędny — on istnieje, a to on rozstrzyga,
	// gdzie plik naprawdę powstanie.
	doSprawdzenia := pelna
	if _, err := os.Lstat(pelna); err != nil {
		doSprawdzenia = filepath.Dir(pelna)
	}
	if rozwiniety, err := filepath.EvalSymlinks(doSprawdzenia); err == nil {
		wzgledna, err := filepath.Rel(korzenPelny, rozwiniety)
		if err != nil || wzgledna == ".." || strings.HasPrefix(wzgledna, ".."+string(filepath.Separator)) {
			return "", ErrPozaObszarem
		}
	}
	wzgledna, err := filepath.Rel(korzenPelny, pelna)
	if err != nil || wzgledna == ".." || strings.HasPrefix(wzgledna, ".."+string(filepath.Separator)) {
		return "", ErrPozaObszarem
	}
	return pelna, nil
}

// Zaloz zakłada plik albo katalog w katalogu roboczym okna.
//
// Katalogi pośrednie powstają same: Operator, który podaje `pkg/nowy/plik.go`,
// prosi o plik pod tą ścieżką, a nie o komunikat, że `pkg/nowy` nie istnieje.
// Plik istniejący NIE jest nadpisywany — zakładanie nie jest zapisem i nie ma
// prawa skasować cudzej treści.
func Zaloz(korzen, sciezka string, katalog bool) error {
	pelna, err := wObszarze(korzen, sciezka)
	if err != nil {
		return err
	}
	if _, err := os.Lstat(pelna); err == nil {
		return ErrJuzIstnieje
	}
	if katalog {
		if err := os.MkdirAll(pelna, 0o755); err != nil {
			return fmt.Errorf("założenie katalogu %s: %w", sciezka, err)
		}
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(pelna), 0o755); err != nil {
		return fmt.Errorf("założenie katalogu nadrzędnego dla %s: %w", sciezka, err)
	}
	plik, err := os.OpenFile(pelna, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("założenie pliku %s: %w", sciezka, err)
	}
	return plik.Close()
}

// ZmienNazwe zmienia nazwę węzła bez zmiany jego miejsca.
//
// Nowa nazwa jest nazwą, nie ścieżką: zmiana nazwy, która przy okazji przenosi
// plik do innego katalogu, jest przeniesieniem — a przeniesienie ma własną
// czynność i własny wykaz w odpowiedzi.
func ZmienNazwe(korzen, sciezka, nowaNazwa string) (string, error) {
	nowaNazwa = strings.TrimSpace(nowaNazwa)
	if nowaNazwa == "" || strings.ContainsAny(nowaNazwa, `/\`) {
		return "", fmt.Errorf("nowa nazwa %q jest ścieżką, a nie nazwą; "+
			"do przeniesienia służy developer.file.move", nowaNazwa)
	}
	pelna, err := wObszarze(korzen, sciezka)
	if err != nil {
		return "", err
	}
	if _, err := os.Lstat(pelna); err != nil {
		return "", fmt.Errorf("węzeł %s nie istnieje", sciezka)
	}
	docelowa := filepath.Join(filepath.Dir(pelna), nowaNazwa)
	if _, err := os.Lstat(docelowa); err == nil {
		return "", ErrJuzIstnieje
	}
	if err := os.Rename(pelna, docelowa); err != nil {
		return "", fmt.Errorf("zmiana nazwy %s: %w", sciezka, err)
	}
	wzgledna, err := filepath.Rel(korzen, docelowa)
	if err != nil {
		return nowaNazwa, nil
	}
	return filepath.ToSlash(wzgledna), nil
}

// Usun usuwa wskazane węzły drzewa i oddaje te, które naprawdę usunięto.
//
// Wykaz usuniętych jest wykazem skutku, nie zamiaru: węzeł, którego nie było,
// nie zatrzymuje czynności — Operator zaznaczający kilkanaście plików nie ma
// dostawać odmowy przez jeden, który zniknął wcześniej. Ale nie ma go też
// w wykazie, bo tego pliku ta czynność nie usunęła.
func Usun(korzen string, sciezki []string) ([]string, error) {
	usuniete := make([]string, 0, len(sciezki))
	for _, sciezka := range sciezki {
		pelna, err := wObszarze(korzen, sciezka)
		if err != nil {
			return nil, err
		}
		if _, err := os.Lstat(pelna); err != nil {
			continue
		}
		if err := os.RemoveAll(pelna); err != nil {
			return nil, fmt.Errorf("usunięcie %s: %w", sciezka, err)
		}
		usuniete = append(usuniete, filepath.ToSlash(filepath.Clean(sciezka)))
	}
	sort.Strings(usuniete)
	return usuniete, nil
}

// Przenies przenosi węzły do wskazanego katalogu.
//
// Katalog docelowy musi istnieć: przeniesienie do katalogu, którego nie ma, jest
// zwykle literówką w nazwie, a założenie go po cichu zostawiłoby Operatora
// z plikami w miejscu, którego nie zna.
func Przenies(korzen string, sciezki []string, katalogDocelowy string) ([]string, error) {
	celPelny, err := wObszarze(korzen, katalogDocelowy)
	if err != nil {
		return nil, err
	}
	opis, err := os.Stat(celPelny)
	if err != nil || !opis.IsDir() {
		return nil, fmt.Errorf("katalog docelowy %s nie istnieje", katalogDocelowy)
	}

	przeniesione := make([]string, 0, len(sciezki))
	for _, sciezka := range sciezki {
		pelna, err := wObszarze(korzen, sciezka)
		if err != nil {
			return nil, err
		}
		if _, err := os.Lstat(pelna); err != nil {
			return nil, fmt.Errorf("węzeł %s nie istnieje", sciezka)
		}
		docelowa := filepath.Join(celPelny, filepath.Base(pelna))
		if docelowa == pelna {
			continue
		}
		if _, err := os.Lstat(docelowa); err == nil {
			return nil, fmt.Errorf("w katalogu %s już leży %s",
				katalogDocelowy, filepath.Base(pelna))
		}
		if err := os.Rename(pelna, docelowa); err != nil {
			return nil, fmt.Errorf("przeniesienie %s: %w", sciezka, err)
		}
		wzgledna, err := filepath.Rel(korzen, docelowa)
		if err != nil {
			continue
		}
		przeniesione = append(przeniesione, filepath.ToSlash(wzgledna))
	}
	sort.Strings(przeniesione)
	return przeniesione, nil
}

// KorzenRoboczy oddaje katalog, w którym pracują czynności plikowe okna.
//
// Repozytorium wskazuje swój katalog główny; katalog bez repozytorium jest
// swoim własnym korzeniem. Dzięki temu czynności plikowe działają także tam,
// gdzie Operator jeszcze nie założył repozytorium.
func KorzenRoboczy(katalog string) (string, error) {
	if strings.TrimSpace(katalog) == "" {
		return "", fmt.Errorf("katalog roboczy okna nie jest wskazany")
	}
	if repo, err := Otworz(katalog); err == nil {
		if korzen, err := Korzen(repo); err == nil {
			return korzen, nil
		}
	}
	opis, err := os.Stat(katalog)
	if err != nil || !opis.IsDir() {
		return "", fmt.Errorf("katalog roboczy %s nie istnieje", katalog)
	}
	return katalog, nil
}
