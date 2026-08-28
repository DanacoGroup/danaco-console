package repozytorium

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"

	"danacoconsole/shared"
)

// repozytoriumProbne zakłada repozytorium wraz z pierwszym zatwierdzeniem, tworząc materiał
// sprawdzianu niezależny od stanu repozytorium produktu.
func repozytoriumProbne(t *testing.T, pliki map[string]string) (string, *git.Repository) {
	t.Helper()

	katalog := t.TempDir()
	repo, err := git.PlainInit(katalog, false)
	if err != nil {
		t.Fatalf("nie można założyć repozytorium: %v", err)
	}
	for sciezka, tresc := range pliki {
		zapiszPlik(t, katalog, sciezka, tresc)
	}
	zatwierdz(t, repo, "zaczyn")
	return katalog, repo
}

// zapiszPlik odkłada treść pod wskazaną ścieżką wraz z katalogami pośrednimi, zakładając
// brakujące katalogi nadrzędne przed zapisem pliku.
func zapiszPlik(t *testing.T, korzen, sciezka, tresc string) {
	t.Helper()

	pelna := filepath.Join(korzen, filepath.FromSlash(sciezka))
	if err := os.MkdirAll(filepath.Dir(pelna), 0o755); err != nil {
		t.Fatalf("nie można założyć katalogu dla %s: %v", sciezka, err)
	}
	if err := os.WriteFile(pelna, []byte(tresc), 0o644); err != nil {
		t.Fatalf("nie można zapisać %s: %v", sciezka, err)
	}
}

// zatwierdz przygotowuje całość zmian w katalogu roboczym i zakłada zatwierdzenie z podanym
// opisem oraz stałym autorstwem sprawdzianu.
func zatwierdz(t *testing.T, repo *git.Repository, opis string) {
	t.Helper()

	drzewo, err := repo.Worktree()
	if err != nil {
		t.Fatalf("nie można sięgnąć po katalog roboczy: %v", err)
	}
	if err := drzewo.AddWithOptions(&git.AddOptions{All: true}); err != nil {
		t.Fatalf("nie można przygotować zmian: %v", err)
	}
	_, err = drzewo.Commit(opis, &git.CommitOptions{
		Author: &object.Signature{
			Name: "Sprawdzian", Email: "sprawdzian@danaco", When: time.Now(),
		},
	})
	if err != nil {
		t.Fatalf("nie można zatwierdzić: %v", err)
	}
}

// TestStanOdrozniaKatalogBezRepozytoriumOdCzystego pilnuje rozróżnienia, które niesie pole
// isRepository. Bez niego panel repozytorium pokazywałby brak zmian tam, gdzie repozytorium
// nigdy nie założono.
func TestStanOdrozniaKatalogBezRepozytoriumOdCzystego(t *testing.T) {
	bezRepozytorium := t.TempDir()
	stan, err := Stan(bezRepozytorium)
	if err != nil {
		t.Fatalf("odczyt katalogu bez repozytorium ma się udać, a oddał błąd: %v", err)
	}
	if stan.IsRepository {
		t.Fatal("katalog bez repozytorium melduje się jako repozytorium")
	}

	katalog, _ := repozytoriumProbne(t, map[string]string{"a.txt": "treść\n"})
	czysty, err := Stan(katalog)
	if err != nil {
		t.Fatalf("odczyt stanu repozytorium: %v", err)
	}
	if !czysty.IsRepository {
		t.Fatal("repozytorium melduje się jako katalog bez repozytorium")
	}
	if len(czysty.Entries) != 0 {
		t.Fatalf("repozytorium tuż po zatwierdzeniu melduje %d zmian", len(czysty.Entries))
	}
}

// TestStanWidziZmianeIPlikNiesledzony wykazuje, że stan czyta dysk, a nie
// pamięta poprzedniej odpowiedzi.
func TestStanWidziZmianeIPlikNiesledzony(t *testing.T) {
	katalog, _ := repozytoriumProbne(t, map[string]string{"a.txt": "pierwsza\n"})

	zapiszPlik(t, katalog, "a.txt", "druga\n")
	zapiszPlik(t, katalog, "nowy.txt", "świeży\n")

	stan, err := Stan(katalog)
	if err != nil {
		t.Fatalf("odczyt stanu: %v", err)
	}
	widziane := map[string]shared.GitFileState{}
	for _, wpis := range stan.Entries {
		widziane[wpis.Path] = wpis.Worktree
	}
	if widziane["a.txt"] != shared.GitFileStateModified {
		t.Fatalf("plik zmieniony ma stan %q", widziane["a.txt"])
	}
	if widziane["nowy.txt"] != shared.GitFileStateUntracked {
		t.Fatalf("plik nowy ma stan %q", widziane["nowy.txt"])
	}
}

// TestRoznicaWskazujeWierszeZmienioneZNumerami wykazuje skutek: fragment
// różnicy niesie treść i numery wierszy, a nie samą liczbę zmian.
func TestRoznicaWskazujeWierszeZmienioneZNumerami(t *testing.T) {
	katalog, _ := repozytoriumProbne(t, map[string]string{
		"kod.go": "jeden\ndwa\ntrzy\n",
	})
	zapiszPlik(t, katalog, "kod.go", "jeden\nDWA\ntrzy\n")

	fragmenty, binarne, err := Roznica(katalog, ZadanieRoznicy{})
	if err != nil {
		t.Fatalf("liczenie różnicy: %v", err)
	}
	if len(binarne) != 0 {
		t.Fatalf("plik tekstowy trafił do wykazu binarnych: %v", binarne)
	}
	if len(fragmenty) == 0 {
		t.Fatal("różnica pliku zmienionego nie ma ani jednego fragmentu")
	}

	usuniete, dodane := "", ""
	for _, fragment := range fragmenty {
		if fragment.Path != "kod.go" {
			t.Fatalf("fragment wskazuje plik %q", fragment.Path)
		}
		for _, wiersz := range fragment.Lines {
			switch wiersz.Kind {
			case shared.GitDiffLineKindRemoved:
				usuniete = wiersz.Text
				if wiersz.OldLine == nil {
					t.Fatal("wiersz usunięty nie niesie numeru w wersji poprzedniej")
				}
			case shared.GitDiffLineKindAdded:
				dodane = wiersz.Text
				if wiersz.NewLine == nil {
					t.Fatal("wiersz dodany nie niesie numeru w wersji nowej")
				}
			}
		}
	}
	if usuniete != "dwa" || dodane != "DWA" {
		t.Fatalf("różnica mówi o wierszach %q → %q, a zmiana była dwa → DWA",
			usuniete, dodane)
	}
}

// TestRoznicaNieUdajeRoznicyPlikuBinarnego pilnuje wykazu osobnego: udawanie
// różnicy wierszowej na pliku binarnym pokazałoby Operatorowi śmieci.
func TestRoznicaNieUdajeRoznicyPlikuBinarnego(t *testing.T) {
	katalog, _ := repozytoriumProbne(t, map[string]string{"obraz.bin": "\x00\x01\x02"})
	zapiszPlik(t, katalog, "obraz.bin", "\x00\x09\x08")

	fragmenty, binarne, err := Roznica(katalog, ZadanieRoznicy{})
	if err != nil {
		t.Fatalf("liczenie różnicy: %v", err)
	}
	if len(fragmenty) != 0 {
		t.Fatalf("plik binarny dostał %d fragmentów różnicy", len(fragmenty))
	}
	if len(binarne) != 1 || binarne[0] != "obraz.bin" {
		t.Fatalf("wykaz binarnych niesie %v", binarne)
	}
}

// TestHistoriaLiczyWszystkieAOddajeTyleIleGranica pilnuje rozróżnienia między
// „ile ich jest" a „ile pokazano". Bez niego klient nie ma jak napisać
// „pokazano 1 z 3".
func TestHistoriaLiczyWszystkieAOddajeTyleIleGranica(t *testing.T) {
	katalog, repo := repozytoriumProbne(t, map[string]string{"a.txt": "raz\n"})
	zapiszPlik(t, katalog, "a.txt", "dwa\n")
	zatwierdz(t, repo, "druga zmiana")
	zapiszPlik(t, katalog, "a.txt", "trzy\n")
	zatwierdz(t, repo, "trzecia zmiana")

	zatwierdzenia, wszystkich, err := Historia(katalog, ZadanieHistorii{Granica: 1})
	if err != nil {
		t.Fatalf("odczyt historii: %v", err)
	}
	if len(zatwierdzenia) != 1 {
		t.Fatalf("granica jeden oddała %d zatwierdzeń", len(zatwierdzenia))
	}
	if wszystkich != 3 {
		t.Fatalf("liczba wszystkich zatwierdzeń to %d, a założono trzy", wszystkich)
	}
	if zatwierdzenia[0].Message != "trzecia zmiana" {
		t.Fatalf("historia zaczyna się od %q, a najnowsze jest trzecia zmiana",
			zatwierdzenia[0].Message)
	}
	if len(zatwierdzenia[0].ShortId) != 7 {
		t.Fatalf("skrót identyfikatora ma %d znaków", len(zatwierdzenia[0].ShortId))
	}
}

// TestSzukanieOmijaPlikiIgnorowane wykazuje, że wynik jest użyteczny: bez
// pominięcia katalogów budowania wyszukiwanie oddaje tysiące trafień w kodzie,
// którego Operator nie pisał.
func TestSzukanieOmijaPlikiIgnorowane(t *testing.T) {
	katalog, _ := repozytoriumProbne(t, map[string]string{
		".gitignore":          "zaleznosci/\n",
		"kod.go":              "szukanaFraza w kodzie\n",
		"zaleznosci/cudze.go": "szukanaFraza w cudzym kodzie\n",
	})

	trafienia, wszystkich, przyciete, err := Szukaj(katalog, ZadanieSzukania{
		Wzorzec: "szukanaFraza", RozrozniajWielkosc: true,
	})
	if err != nil {
		t.Fatalf("wyszukiwanie: %v", err)
	}
	if przyciete {
		t.Fatal("wynik przycięto, choć trafień było kilka")
	}
	if wszystkich != 1 || len(trafienia) != 1 {
		t.Fatalf("wyszukiwanie oddało %d trafień, a poza katalogiem ignorowanym "+
			"jest jedno: %v", len(trafienia), trafienia)
	}
	if trafienia[0].Path != "kod.go" || trafienia[0].Line != 1 {
		t.Fatalf("trafienie wskazuje %s:%d", trafienia[0].Path, trafienia[0].Line)
	}
}

// TestWzorzecDoslownyNieJestWyrazeniem pilnuje znaczenia wzorca bez zaznaczenia
// „wyrażenie regularne": Operator szukający `config.json` nie szuka „config"
// z dowolnym znakiem w środku.
func TestWzorzecDoslownyNieJestWyrazeniem(t *testing.T) {
	katalog, _ := repozytoriumProbne(t, map[string]string{
		"a.txt": "config.json\nconfigXjson\n",
	})
	trafienia, _, _, err := Szukaj(katalog, ZadanieSzukania{
		Wzorzec: "config.json", RozrozniajWielkosc: true,
	})
	if err != nil {
		t.Fatalf("wyszukiwanie: %v", err)
	}
	if len(trafienia) != 1 || trafienia[0].Line != 1 {
		t.Fatalf("wzorzec dosłowny trafił w %d wierszy zamiast w jeden", len(trafienia))
	}
}

// TestPodgladZamianyNieDotykaPlikow wykazuje różnicę między podglądem
// a zapisem: zamiana masowa nie ma „Cofnij", więc podgląd musi być naprawdę
// podglądem.
func TestPodgladZamianyNieDotykaPlikow(t *testing.T) {
	katalog, _ := repozytoriumProbne(t, map[string]string{"a.txt": "stara nazwa\n"})

	zmiany, zmienione, zapisano, err := Zamien(katalog, ZadanieZamiany{
		Wzorzec: "stara", Zamiana: "nowa", Podglad: true,
	})
	if err != nil {
		t.Fatalf("zamiana: %v", err)
	}
	if zapisano {
		t.Fatal("podgląd melduje zapis")
	}
	if len(zmiany) != 1 || len(zmienione) != 1 {
		t.Fatalf("podgląd oddał %d zmian w %d plikach", len(zmiany), len(zmienione))
	}
	bajty, err := os.ReadFile(filepath.Join(katalog, "a.txt"))
	if err != nil {
		t.Fatalf("odczyt pliku: %v", err)
	}
	if string(bajty) != "stara nazwa\n" {
		t.Fatalf("podgląd zmienił plik na dysku: %q", string(bajty))
	}

	if _, _, zapisano, err = Zamien(katalog, ZadanieZamiany{
		Wzorzec: "stara", Zamiana: "nowa",
	}); err != nil {
		t.Fatalf("zamiana z zapisem: %v", err)
	}
	if !zapisano {
		t.Fatal("zamiana z zapisem melduje brak zapisu")
	}
	bajty, err = os.ReadFile(filepath.Join(katalog, "a.txt"))
	if err != nil {
		t.Fatalf("odczyt pliku: %v", err)
	}
	if string(bajty) != "nowa nazwa\n" {
		t.Fatalf("po zamianie plik niesie %q", string(bajty))
	}
}

// TestCzynnosciPlikoweOdmawiajaWyjsciaPozaKatalogRoboczy pilnuje granicy
// obszaru. Bez niej `../../` w nazwie pliku byłoby drogą do dowolnego miejsca
// na dysku serwera.
func TestCzynnosciPlikoweOdmawiajaWyjsciaPozaKatalogRoboczy(t *testing.T) {
	katalog, _ := repozytoriumProbne(t, map[string]string{"a.txt": "treść\n"})

	if err := Zaloz(katalog, "../poza.txt", false); err == nil {
		t.Fatal("założenie pliku poza katalogiem roboczym udało się")
	} else if !strings.Contains(err.Error(), "poza katalog") {
		t.Fatalf("odmowa nie mówi o wyjściu poza obszar: %v", err)
	}
	if _, err := Usun(katalog, []string{"../../etc/passwd"}); err == nil {
		t.Fatal("usunięcie poza katalogiem roboczym udało się")
	}
	if _, err := os.Stat(filepath.Join(filepath.Dir(katalog), "poza.txt")); err == nil {
		t.Fatal("plik poza katalogiem roboczym jednak powstał")
	}
}

// TestZakladanieNieNadpisujeIstniejacego pilnuje, żeby zakładanie nie kasowało
// cudzej treści: zakładanie nie jest zapisem.
func TestZakladanieNieNadpisujeIstniejacego(t *testing.T) {
	katalog, _ := repozytoriumProbne(t, map[string]string{"a.txt": "treść\n"})

	if err := Zaloz(katalog, "a.txt", false); err == nil {
		t.Fatal("założenie pliku istniejącego udało się")
	}
	bajty, err := os.ReadFile(filepath.Join(katalog, "a.txt"))
	if err != nil {
		t.Fatalf("odczyt pliku: %v", err)
	}
	if string(bajty) != "treść\n" {
		t.Fatalf("treść pliku zmieniła się na %q", string(bajty))
	}

	if err := Zaloz(katalog, "pkg/nowy/plik.go", false); err != nil {
		t.Fatalf("zakładanie z katalogami pośrednimi: %v", err)
	}
	if _, err := os.Stat(filepath.Join(katalog, "pkg", "nowy", "plik.go")); err != nil {
		t.Fatalf("plik z katalogami pośrednimi nie powstał: %v", err)
	}
}

// TestZmianaNazwyOdmawiaSciezki pilnuje granicy między zmianą nazwy
// a przeniesieniem: obie czynności mają w kontrakcie własne komendy i własne
// odpowiedzi.
func TestZmianaNazwyOdmawiaSciezki(t *testing.T) {
	katalog, _ := repozytoriumProbne(t, map[string]string{"a.txt": "treść\n"})

	if _, err := ZmienNazwe(katalog, "a.txt", "pkg/b.txt"); err == nil {
		t.Fatal("zmiana nazwy przyjęła ścieżkę zamiast nazwy")
	}
	nowa, err := ZmienNazwe(katalog, "a.txt", "b.txt")
	if err != nil {
		t.Fatalf("zmiana nazwy: %v", err)
	}
	if nowa != "b.txt" {
		t.Fatalf("czynność melduje nową ścieżkę %q", nowa)
	}
	if _, err := os.Stat(filepath.Join(katalog, "b.txt")); err != nil {
		t.Fatalf("plik pod nową nazwą nie istnieje: %v", err)
	}
}

// TestKonfliktRozkladaPlikNaTrzyWersje wykazuje skutek widoku trójstronnego:
// z pliku ze znacznikami scalania wychodzą trzy czytelne wersje.
func TestKonfliktRozkladaPlikNaTrzyWersje(t *testing.T) {
	katalog, _ := repozytoriumProbne(t, map[string]string{"a.txt": "start\n"})
	zapiszPlik(t, katalog, "a.txt", strings.Join([]string{
		"wspólny nagłówek",
		"<<<<<<< HEAD",
		"moja wersja",
		"=======",
		"cudza wersja",
		">>>>>>> galaz",
		"wspólna stopka",
	}, "\n"))

	konflikt, err := Konflikt(katalog, "a.txt")
	if err != nil {
		t.Fatalf("rozbiór konfliktu: %v", err)
	}
	if !strings.Contains(konflikt.Current, "moja wersja") ||
		strings.Contains(konflikt.Current, "cudza wersja") {
		t.Fatalf("wersja bieżąca niesie %q", konflikt.Current)
	}
	if !strings.Contains(konflikt.Incoming, "cudza wersja") ||
		strings.Contains(konflikt.Incoming, "moja wersja") {
		t.Fatalf("wersja przychodząca niesie %q", konflikt.Incoming)
	}
	// Treść wspólna wchodzi do obu stron: widok trójstronny pokazuje całe pliki.
	if !strings.Contains(konflikt.Current, "wspólny nagłówek") ||
		!strings.Contains(konflikt.Incoming, "wspólna stopka") {
		t.Fatal("treść wspólna nie weszła do obu wersji")
	}

	pozostale, err := RozstrzygnijKonflikt(katalog, "a.txt",
		shared.ConflictResolutionKindTakeIncoming, "")
	if err != nil {
		t.Fatalf("rozstrzygnięcie konfliktu: %v", err)
	}
	if len(pozostale) != 0 {
		t.Fatalf("po rozstrzygnięciu zostało %v", pozostale)
	}
	bajty, err := os.ReadFile(filepath.Join(katalog, "a.txt"))
	if err != nil {
		t.Fatalf("odczyt pliku: %v", err)
	}
	if strings.Contains(string(bajty), "<<<<<<<") {
		t.Fatal("znaczniki konfliktu zostały w pliku po rozstrzygnięciu")
	}
	if !strings.Contains(string(bajty), "cudza wersja") {
		t.Fatalf("plik nie niesie wersji przychodzącej: %q", string(bajty))
	}
}
