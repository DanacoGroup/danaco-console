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
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

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
//
// Wzorzec z metazmienną idzie drogą składni (`ast-grep`), reszta — drogą napisu.
// Rozstrzygnięcie stoi przy `wzorzecPoSkladni`.
func (a *adapterDevelopera) SzukajWRepozytorium(ctx context.Context,
	z shared.DeveloperGrepSearchRequest) (shared.DeveloperGrepSearchResponse, error) {

	okno, korzen, err := a.korzenRoboczyOkna(z.WindowId)
	if err != nil {
		return shared.DeveloperGrepSearchResponse{}, err
	}
	if wzorzecPoSkladni(z.Pattern, z.Regex) {
		return a.szukajPoSkladni(ctx, okno, korzen, z)
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
//
// Wzorzec z metazmienną idzie drogą składni, tą samą, co wyszukanie.
func (a *adapterDevelopera) ZamienWRepozytorium(ctx context.Context,
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
	if wzorzecPoSkladni(z.Pattern, z.Regex) {
		return a.zamienPoSkladni(ctx, okno, korzen, z, podglad)
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

// ── Droga składni ───────────────────────────────────────────────────────────

// czasSzukaniaPoSkladni jest granicą jednego przejścia po repozytorium.
// Rozbiór składni jest droższy od dopasowania napisu, a repozytorium bywa duże.
const czasSzukaniaPoSkladni = 120 * time.Second

// metazmiennaWzorca rozpoznaje metazmienną `ast-grep`: `$` wraz z nazwą pisaną
// wersalikami, cyframi i podkreśleniem, zaczynającą się od wersalika albo
// podkreślenia.
var metazmiennaWzorca = regexp.MustCompile(`\$[A-Z_][A-Z0-9_]*`)

// wzorzecPoSkladni orzeka, czy wzorzec ma iść drogą składni zamiast drogi napisu.
//
// ── Dlaczego rozstrzyga metazmienna, a nie osobne pole ──────────────────────
// Kontrakt nie niesie pola „szukaj po składni" i niniejsza praca kontraktu nie
// zmienia. Rozstrzyga więc sam wzorzec, i to jego własną, udokumentowaną cechą:
// metazmienna `$NAZWA` jest zapisem należącym do `ast-grep` i nie znaczy nic
// w wyszukiwaniu po napisie — napis `$ARG` jako napis szukany jest zapytaniem,
// którego nikt nie zadaje.
//
// Wyrażenie regularne wyłącza tę drogę BEZWARUNKOWO: w wyrażeniu `$` jest kotwicą
// końca wiersza, więc wzorzec `foo$` byłby wzięty za składniowy wbrew temu, co
// wołający napisał wprost.
func wzorzecPoSkladni(wzorzec string, wyrazenie *bool) bool {
	if wyrazenie != nil && *wyrazenie {
		return false
	}
	return metazmiennaWzorca.MatchString(wzorzec)
}

// trafienieSkladni jest kształtem jednego trafienia `ast-grep run --json`.
//
// Wiersze i kolumny `ast-grep` liczy od zera, a kontrakt od jedynki — przeliczenie
// stoi w miejscu odczytu, żeby nie rozjechało się między wyszukaniem a zamianą.
type trafienieSkladni struct {
	Plik    string `json:"file"`
	Wiersze string `json:"lines"`
	Zakres  struct {
		Poczatek struct {
			Wiersz  int `json:"line"`
			Kolumna int `json:"column"`
		} `json:"start"`
		Koniec struct {
			Wiersz  int `json:"line"`
			Kolumna int `json:"column"`
		} `json:"end"`
	} `json:"range"`
	Zamiana string `json:"replacement"`
}

// argumentySkladni składa wspólną część wywołania: wzorzec oraz zawężenia
// wskazane przez wołającego.
//
// `include` i `exclude` kontraktu są wzorcami glob i idą jednym parametrem
// `--globs`, w którym wyłączenie znakuje się wykrzyknikiem — tak, jak robi to
// `.gitignore`.
func argumentySkladni(wzorzec string, wlacz, wylacz []string) []string {
	argumenty := []string{"run", "--pattern", wzorzec, "--json=compact"}
	for _, glob := range wlacz {
		if tresc := strings.TrimSpace(glob); tresc != "" {
			argumenty = append(argumenty, "--globs", tresc)
		}
	}
	for _, glob := range wylacz {
		if tresc := strings.TrimSpace(glob); tresc != "" {
			argumenty = append(argumenty, "--globs", "!"+tresc)
		}
	}
	return argumenty
}

// odmowaBrakuSkladni nazywa brak programu wraz z drogą naprawy i obejściem.
//
// Odmowa, a nie cichy powrót do drogi napisu: wzorzec składniowy szukany jako
// napis nie znajdzie niczego, a zero trafień wyglądałoby jak odpowiedź. Wynik
// wyglądający dobrze jest groźniejszy od odmowy, bo nie wzywa do sprawdzenia.
func odmowaBrakuSkladni() error {
	return bladZasobuDevelopera(
		"serwer nie ma programu " + narzedzieAstGrep.Nazwa + " (" +
			narzedzieAstGrep.Program + "), a wzorzec z metazmienną `$NAZWA` " +
			"jest wzorcem składni i po napisie się nie znajdzie; naprawa: " +
			narzedzieAstGrep.Pakiet + ". Droga, która działa bez tego programu: " +
			"wzorzec bez metazmiennej albo wyrażenie regularne (`regex: true`)")
}

// szukajPoSkladni przeprowadza wyszukanie wzorca po składni.
func (a *adapterDevelopera) szukajPoSkladni(ctx context.Context, okno session.Okno,
	korzen string, z shared.DeveloperGrepSearchRequest) (shared.DeveloperGrepSearchResponse, error) {

	trafienia, err := a.przejdzPoSkladni(ctx, okno, korzen,
		append(argumentySkladni(z.Pattern, z.Include, z.Exclude), "."))
	if err != nil {
		return shared.DeveloperGrepSearchResponse{}, err
	}

	wykaz := make([]shared.DeveloperGrepMatch, 0, len(trafienia))
	for _, trafienie := range trafienia {
		wykaz = append(wykaz, shared.DeveloperGrepMatch{
			Path:   sciezkaWzgledemKorzenia(trafienie.Plik, korzen),
			Line:   trafienie.Zakres.Poczatek.Wiersz + 1,
			Column: wskaznikLiczby(trafienie.Zakres.Poczatek.Kolumna + 1),
			Text:   trafienie.Wiersze,
		})
	}

	wszystkich := len(wykaz)
	przyciete := false
	if z.Limit != nil && *z.Limit > 0 && len(wykaz) > *z.Limit {
		wykaz, przyciete = wykaz[:*z.Limit], true
	}
	return shared.DeveloperGrepSearchResponse{
		Matches: wykaz, Total: &wszystkich, Truncated: przyciete,
	}, nil
}

// zamienPoSkladni przeprowadza zamianę wzorca po składni.
//
// Przejścia są dwa i jest to rozstrzygnięcie: pierwsze — bez zapisu — daje wykaz
// zmian, który wraca w odpowiedzi także przy zapisie, bo Operator ma zobaczyć,
// co się stało. Drugie zapisuje. Wyprowadzenie wykazu z samego zapisu nie da się
// zrobić: przy `--update-all` program oddaje liczbę zmian, a nie ich treść.
func (a *adapterDevelopera) zamienPoSkladni(ctx context.Context, okno session.Okno,
	korzen string, z shared.DeveloperGrepReplaceRequest,
	podglad bool) (shared.DeveloperGrepReplaceResponse, error) {

	wspolne := append(argumentySkladni(z.Pattern, nil, nil),
		"--rewrite", z.Replacement)
	cele := celeZamianySkladni(z.Paths)

	trafienia, err := a.przejdzPoSkladni(ctx, okno, korzen, append(wspolne, cele...))
	if err != nil {
		return shared.DeveloperGrepReplaceResponse{}, err
	}

	zmiany := make([]shared.DeveloperTextEdit, 0, len(trafienia))
	widziane := map[string]bool{}
	zmienione := make([]string, 0, 4)
	for _, trafienie := range trafienia {
		sciezka := sciezkaWzgledemKorzenia(trafienie.Plik, korzen)
		zmiany = append(zmiany, shared.DeveloperTextEdit{
			Path:        sciezka,
			StartLine:   trafienie.Zakres.Poczatek.Wiersz + 1,
			StartColumn: wskaznikLiczby(trafienie.Zakres.Poczatek.Kolumna + 1),
			EndLine:     trafienie.Zakres.Koniec.Wiersz + 1,
			EndColumn:   wskaznikLiczby(trafienie.Zakres.Koniec.Kolumna + 1),
			NewText:     trafienie.Zamiana,
		})
		if !widziane[sciezka] {
			widziane[sciezka] = true
			zmienione = append(zmienione, sciezka)
		}
	}

	if podglad || len(zmiany) == 0 {
		return shared.DeveloperGrepReplaceResponse{
			Edits: zmiany, ChangedPaths: zmienione, Applied: false,
		}, nil
	}

	// Zapis nie oddaje trafień, więc wykaz JSON schodzi z wywołania: program
	// odmawia postawienia obu parametrów naraz.
	zapis := append(usunWykazJson(wspolne), "--update-all")
	zapis = append(zapis, cele...)
	if _, err := a.wolajNarzedzieWarsztatu(ctx, okno, narzedzieAstGrep, zapis,
		korzen, czasSzukaniaPoSkladni); err != nil {
		if brakNarzedziaWarsztatu(err) {
			return shared.DeveloperGrepReplaceResponse{}, odmowaBrakuSkladni()
		}
		return shared.DeveloperGrepReplaceResponse{}, bladWykonaniaDevelopera(
			"zamiana po składni odmówiła: " + err.Error())
	}
	return shared.DeveloperGrepReplaceResponse{
		Edits: zmiany, ChangedPaths: zmienione, Applied: true,
	}, nil
}

// celeZamianySkladni składa wskazania plików; puste żądanie obejmuje całe drzewo.
func celeZamianySkladni(sciezki []string) []string {
	cele := make([]string, 0, len(sciezki))
	for _, sciezka := range sciezki {
		if tresc := strings.TrimSpace(sciezka); tresc != "" {
			cele = append(cele, tresc)
		}
	}
	if len(cele) == 0 {
		return []string{"."}
	}
	return cele
}

// usunWykazJson zdejmuje z wywołania parametr wykazu JSON.
func usunWykazJson(argumenty []string) []string {
	bez := make([]string, 0, len(argumenty))
	for _, argument := range argumenty {
		if strings.HasPrefix(argument, "--json") {
			continue
		}
		bez = append(bez, argument)
	}
	return bez
}

// przejdzPoSkladni uruchamia program i odczytuje wykaz trafień.
//
// Wykaz pusty jest wynikiem; wykaz nieczytelny przy niepowodzeniu programu NIE
// jest — wtedy nie zmierzono niczego i odpowiedź ma to powiedzieć odmową.
func (a *adapterDevelopera) przejdzPoSkladni(ctx context.Context, okno session.Okno,
	korzen string, argumenty []string) ([]trafienieSkladni, error) {

	wynik, err := a.wolajNarzedzieWarsztatu(ctx, okno, narzedzieAstGrep, argumenty,
		korzen, czasSzukaniaPoSkladni)
	if brakNarzedziaWarsztatu(err) {
		return nil, odmowaBrakuSkladni()
	}
	var trafienia []trafienieSkladni
	tresc := strings.TrimSpace(string(wynik.Wyjscie))
	if rozbior := json.Unmarshal([]byte(tresc), &trafienia); rozbior != nil {
		if err != nil {
			return nil, bladWykonaniaDevelopera(
				"przejście po składni odmówiło: " + skrocDiagnostyke(wynik.Diagnostyka, err))
		}
		return nil, bladWykonaniaDevelopera(
			"program " + narzedzieAstGrep.Program + " oddał odpowiedź, której nie da się " +
				"odczytać jako wykazu trafień: " + rozbior.Error())
	}
	return trafienia, nil
}
