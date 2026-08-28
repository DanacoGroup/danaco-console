// Moduł obsługuje czynności na węzłach drzewa projektu — zakładanie, zmianę
// nazwy, usuwanie i przenoszenie — oraz wyszukiwanie i zamianę w całym
// repozytorium katalogu roboczego okna, przez bibliotekę standardową.
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

// wezelPoSciezce składa węzeł kontraktu dla istniejącej ścieżki względnej,
// czytając rodzaj i inne metadane bezpośrednio z systemu plików.
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

// ZalozWezel zakłada plik albo katalog w katalogu roboczym okna. Treść
// początkowa idzie osobnym zapisem po założeniu, nie zamiast niego, żeby
// sprawdzenie zajętości miejsca docelowego zadziałało też przy treści.
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

// ZmienNazweWezla zmienia nazwę węzła w katalogu roboczym okna, zostawiając
// jego miejsce w drzewie bez zmiany.
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

// UsunWezly usuwa wskazane węzły drzewa i oddaje ścieżki naprawdę usunięte,
// nie te, o które proszono — węzeł, którego już nie było, nie wchodzi do
// wykazu, bo tego węzła ta czynność nie usunęła.
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
	// Katalog niepusty schodzi wyłącznie przy jawnym wskazaniu rekurencji.
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

// PrzeniesWezly przenosi wskazane węzły do wskazanego katalogu docelowego
// w obrębie katalogu roboczego okna.
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

// ZamienWRepozytorium zamienia trafienia wzorca w repozytorium i oddaje
// wykaz zmian także przy samym podglądzie, bez zapisu, oraz stan zapisu.
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

// wzorzecPoSkladni orzeka, czy wzorzec ma iść drogą składni zamiast drogi
// napisu, po obecności metazmiennej `$NAZWA`; wyrażenie regularne wyłącza tę
// drogę bezwarunkowo.
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
// glob wskazane przez wołającego, każde parametrem `--globs`, z wyłączeniem
// znakowanym wykrzyknikiem.
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

// odmowaBrakuSkladni nazywa brak programu składniowego wraz z drogą naprawy
// i obejściem, zamiast cichego powrotu do drogi napisu bez trafień.
func odmowaBrakuSkladni() error {
	return bladZasobuDevelopera(
		"serwer nie ma programu " + narzedzieAstGrep.Nazwa + " (" +
			narzedzieAstGrep.Program + "), a wzorzec z metazmienną `$NAZWA` " +
			"jest wzorcem składni i po napisie się nie znajdzie; naprawa: " +
			narzedzieAstGrep.Pakiet + ". Droga, która działa bez tego programu: " +
			"wzorzec bez metazmiennej albo wyrażenie regularne (`regex: true`)")
}

// szukajPoSkladni przeprowadza wyszukanie wzorca po składni programem
// zewnętrznym i przekłada jego trafienia na wynik kontraktu.
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

// zamienPoSkladni przeprowadza zamianę wzorca po składni w dwóch przejściach
// programu zewnętrznego: pierwsze bez zapisu daje wykaz zmian, drugie zapisuje.
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

	// Zapis nie oddaje trafień, więc wykaz JSON schodzi z wywołania.
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

// celeZamianySkladni składa wskazania plików do zamiany; puste żądanie
// obejmuje całe drzewo katalogu roboczego.
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

// usunWykazJson zdejmuje z wywołania parametr wykazu JSON, niedopuszczalny
// razem z parametrem zapisu masowego.
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
