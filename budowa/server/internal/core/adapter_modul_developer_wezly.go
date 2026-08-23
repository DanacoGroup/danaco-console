// Odpowiedzialność pliku: czynności na węzłach drzewa projektu — zakładanie,
// zmiana nazwy, usuwanie i przenoszenie — oraz wyszukiwanie i zamiana w całym
// repozytorium.
//
// ── Biblioteka wkompilowana, nie program zewnętrzny ─────────────────────────
// Czynności plikowe idą przez `os` biblioteki standardowej, a wyszukiwanie przez
// `regexp` — nie przez `ripgrep`. Opracowanie modułu wskazywało `ripgrep`, ale
// wyszukiwanie w repozytorium jest czynnością, bez której Code Editor przestaje
// być edytorem kodu: nie może zależeć od programu, którego instalka nie niesie.
//
// ── Granica obszaru jest sprawdzana zawsze ──────────────────────────────────
// Każda ścieżka przechodzi przez `repozytorium`, który sprowadza ją do wnętrza
// katalogu roboczego okna i odmawia, gdy z niego wychodzi — także wtedy, gdy
// wyjściem jest dowiązanie symboliczne, a nie `..` w napisie.
package core

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"danacoconsole/server/internal/repozytorium"
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// korzenRoboczyOkna oddaje katalog, w którym pracują czynności na węzłach, wraz
// z oknem — potrzebnym do sprawdzenia trybu uprawnień.
func (a *adapterDevelopera) korzenRoboczyOkna(oknoKod string) (session.Okno, string, error) {
	okno, err := a.oknoDevelopera(oknoKod)
	if err != nil {
		return session.Okno{}, "", err
	}
	korzenie := a.korzenieOkna(okno)
	if len(korzenie) == 0 {
		return session.Okno{}, "", bladZadaniaDevelopera("okno " + oknoKod +
			" nie ma katalogu roboczego, więc czynność nie ma gdzie pracować")
	}
	korzen, err := repozytorium.KorzenRoboczy(korzenie[0])
	if err != nil {
		return session.Okno{}, "", bladZasobuDevelopera(err.Error())
	}
	return okno, korzen, nil
}

// wezelPoSciezce składa węzeł kontraktu dla istniejącej ścieżki względnej.
func wezelPoSciezce(korzen, wzgledna string) (shared.DeveloperTreeNode, error) {
	pelna := filepath.Join(korzen, filepath.FromSlash(wzgledna))
	opis, err := os.Stat(pelna)
	if err != nil {
		return shared.DeveloperTreeNode{},
			bladZasobuDevelopera("węzeł " + wzgledna + " nie istnieje po czynności")
	}
	rodzic := filepath.ToSlash(filepath.Dir(wzgledna))
	if rodzic == "." {
		rodzic = ""
	}
	return wezelDrzewa(wzgledna, rodzic, opis), nil
}

// ZalozWezel zakłada plik albo katalog w katalogu roboczym okna.
//
// Treść początkowa idzie osobnym zapisem po założeniu, a nie zamiast niego:
// zakładanie odmawia, gdy w miejscu docelowym coś już leży, i to sprawdzenie ma
// zadziałać także wtedy, gdy żądanie niesie treść.
func (a *adapterDevelopera) ZalozWezel(_ context.Context,
	z shared.DeveloperFileCreateRequest) (shared.DeveloperFileCreateResponse, error) {

	okno, korzen, err := a.korzenRoboczyOkna(z.WindowId)
	if err != nil {
		return shared.DeveloperFileCreateResponse{}, err
	}
	if err := sprawdzZmianeSystemu(okno.TrybUprawnien, "założenie węzła"); err != nil {
		return shared.DeveloperFileCreateResponse{}, err
	}
	if strings.TrimSpace(z.Path) == "" {
		return shared.DeveloperFileCreateResponse{},
			bladZadaniaDevelopera("czynność wymaga wskazania ścieżki zakładanego węzła")
	}
	katalog := z.Kind == shared.TreeNodeKindDirectory
	if err := repozytorium.Zaloz(korzen, z.Path, katalog); err != nil {
		return shared.DeveloperFileCreateResponse{}, bladRepozytorium(err)
	}
	if !katalog && z.Content != nil && *z.Content != "" {
		pelna := filepath.Join(korzen, filepath.FromSlash(z.Path))
		if err := os.WriteFile(pelna, []byte(*z.Content), 0o644); err != nil {
			return shared.DeveloperFileCreateResponse{},
				bladWykonaniaDevelopera("zapis treści początkowej: " + err.Error())
		}
	}
	wezel, err := wezelPoSciezce(korzen, filepath.ToSlash(filepath.Clean(z.Path)))
	if err != nil {
		return shared.DeveloperFileCreateResponse{}, err
	}
	return shared.DeveloperFileCreateResponse{Node: wezel}, nil
}

// ZmienNazweWezla zmienia nazwę węzła bez zmiany jego miejsca.
func (a *adapterDevelopera) ZmienNazweWezla(_ context.Context,
	z shared.DeveloperFileRenameRequest) (shared.DeveloperFileRenameResponse, error) {

	okno, korzen, err := a.korzenRoboczyOkna(z.WindowId)
	if err != nil {
		return shared.DeveloperFileRenameResponse{}, err
	}
	if err := sprawdzZmianeSystemu(okno.TrybUprawnien, "zmiana nazwy węzła"); err != nil {
		return shared.DeveloperFileRenameResponse{}, err
	}
	nowa, err := repozytorium.ZmienNazwe(korzen, z.Path, z.NewName)
	if err != nil {
		return shared.DeveloperFileRenameResponse{}, bladRepozytorium(err)
	}
	wezel, err := wezelPoSciezce(korzen, nowa)
	if err != nil {
		return shared.DeveloperFileRenameResponse{}, err
	}
	return shared.DeveloperFileRenameResponse{Node: wezel}, nil
}

// UsunWezly usuwa wskazane węzły drzewa.
//
// Odpowiedź niesie ścieżki NAPRAWDĘ usunięte, a nie te, o które proszono:
// węzeł, którego już nie było, nie zatrzymuje czynności, ale nie wchodzi też do
// wykazu — bo tego węzła ta czynność nie usunęła.
//
// Znacznika cofnięcia odpowiedź nie niesie i to jest stan świadomy: cofnięcie
// wymagałoby odłożenia treści usuniętych węzłów w magazynie rdzenia, a takiego
// magazynu moduł nie ma. Znacznik wypełniony bez pokrycia obiecywałby Operatorowi
// powrót, którego nikt nie wykona.
func (a *adapterDevelopera) UsunWezly(_ context.Context,
	z shared.DeveloperFileDeleteRequest) (shared.DeveloperFileDeleteResponse, error) {

	okno, korzen, err := a.korzenRoboczyOkna(z.WindowId)
	if err != nil {
		return shared.DeveloperFileDeleteResponse{}, err
	}
	if err := sprawdzZmianeSystemu(okno.TrybUprawnien, "usunięcie węzłów"); err != nil {
		return shared.DeveloperFileDeleteResponse{}, err
	}
	if len(z.Paths) == 0 {
		return shared.DeveloperFileDeleteResponse{},
			bladZadaniaDevelopera("czynność nie wskazuje ani jednego węzła do usunięcia")
	}
	// Katalog niepusty schodzi wyłącznie przy jawnym wskazaniu: usunięcie
	// katalogu razem z zawartością na skutek pomyłki w zaznaczeniu byłoby
	// stratą, po której nie ma powrotu.
	if z.Recursive == nil || !*z.Recursive {
		for _, sciezka := range z.Paths {
			pelna := filepath.Join(korzen, filepath.FromSlash(sciezka))
			opis, err := os.Stat(pelna)
			if err != nil || !opis.IsDir() {
				continue
			}
			wpisy, err := os.ReadDir(pelna)
			if err == nil && len(wpisy) > 0 {
				return shared.DeveloperFileDeleteResponse{}, bladZadaniaDevelopera(
					"katalog " + sciezka + " nie jest pusty; usunięcie razem z zawartością " +
						"wymaga jawnego wskazania (recursive)")
			}
		}
	}

	usuniete, err := repozytorium.Usun(korzen, z.Paths)
	if err != nil {
		return shared.DeveloperFileDeleteResponse{}, bladRepozytorium(err)
	}
	return shared.DeveloperFileDeleteResponse{DeletedPaths: usuniete}, nil
}

// PrzeniesWezly przenosi węzły do wskazanego katalogu.
func (a *adapterDevelopera) PrzeniesWezly(_ context.Context,
	z shared.DeveloperFileMoveRequest) (shared.DeveloperFileMoveResponse, error) {

	okno, korzen, err := a.korzenRoboczyOkna(z.WindowId)
	if err != nil {
		return shared.DeveloperFileMoveResponse{}, err
	}
	if err := sprawdzZmianeSystemu(okno.TrybUprawnien, "przeniesienie węzłów"); err != nil {
		return shared.DeveloperFileMoveResponse{}, err
	}
	if len(z.Paths) == 0 {
		return shared.DeveloperFileMoveResponse{},
			bladZadaniaDevelopera("czynność nie wskazuje ani jednego węzła do przeniesienia")
	}
	przeniesione, err := repozytorium.Przenies(korzen, z.Paths, z.TargetPath)
	if err != nil {
		return shared.DeveloperFileMoveResponse{}, bladRepozytorium(err)
	}
	wezly := make([]shared.DeveloperTreeNode, 0, len(przeniesione))
	for _, sciezka := range przeniesione {
		wezel, err := wezelPoSciezce(korzen, sciezka)
		if err != nil {
			return shared.DeveloperFileMoveResponse{}, err
		}
		wezly = append(wezly, wezel)
	}
	return shared.DeveloperFileMoveResponse{Nodes: wezly}, nil
}

// SzukajWRepozytorium przeszukuje repozytorium katalogu roboczego okna.
func (a *adapterDevelopera) SzukajWRepozytorium(_ context.Context,
	z shared.DeveloperGrepSearchRequest) (shared.DeveloperGrepSearchResponse, error) {

	_, korzen, err := a.korzenRoboczyOkna(z.WindowId)
	if err != nil {
		return shared.DeveloperGrepSearchResponse{}, err
	}
	zadanie := repozytorium.ZadanieSzukania{
		Wzorzec: z.Pattern, Wlacz: z.Include, Wylacz: z.Exclude,
	}
	if z.Regex != nil {
		zadanie.Wyrazenie = *z.Regex
	}
	if z.CaseSensitive != nil {
		zadanie.RozrozniajWielkosc = *z.CaseSensitive
	}
	if z.Limit != nil {
		zadanie.Granica = *z.Limit
	}

	trafienia, wszystkich, przyciete, err := repozytorium.Szukaj(korzen, zadanie)
	if err != nil {
		return shared.DeveloperGrepSearchResponse{}, bladZadaniaDevelopera(err.Error())
	}
	return shared.DeveloperGrepSearchResponse{
		Matches: trafienia, Total: &wszystkich, Truncated: przyciete,
	}, nil
}

// ZamienWRepozytorium zamienia trafienia wzorca w repozytorium.
//
// Zamiana masowa dotyka wielu plików naraz i nie ma po niej „Cofnij", więc
// odpowiedź mówi wprost, czy zmiany zapisano, a wykaz zmian wraca także przy
// podglądzie — Operator ma zobaczyć, co się stanie, zanim to się stanie.
func (a *adapterDevelopera) ZamienWRepozytorium(_ context.Context,
	z shared.DeveloperGrepReplaceRequest) (shared.DeveloperGrepReplaceResponse, error) {

	okno, korzen, err := a.korzenRoboczyOkna(z.WindowId)
	if err != nil {
		return shared.DeveloperGrepReplaceResponse{}, err
	}
	podglad := z.Preview != nil && *z.Preview
	if !podglad {
		if err := sprawdzZmianeSystemu(okno.TrybUprawnien, "zamiana w repozytorium"); err != nil {
			return shared.DeveloperGrepReplaceResponse{}, err
		}
	}

	zadanie := repozytorium.ZadanieZamiany{
		Wzorzec: z.Pattern, Zamiana: z.Replacement, Sciezki: z.Paths, Podglad: podglad,
	}
	if z.Regex != nil {
		zadanie.Wyrazenie = *z.Regex
	}
	zmiany, zmienione, zapisano, err := repozytorium.Zamien(korzen, zadanie)
	if err != nil {
		return shared.DeveloperGrepReplaceResponse{}, bladZadaniaDevelopera(err.Error())
	}
	return shared.DeveloperGrepReplaceResponse{
		Edits: zmiany, ChangedPaths: zmienione, Applied: zapisano,
	}, nil
}
