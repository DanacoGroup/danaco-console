// Odpowiedzialność pliku: `developer.tree.get` — drzewo projektu okna Project
// Tree. Żądanie bez wskazania ścieżki zwraca korzenie wszystkich katalogów
// roboczych naraz. Ścieżki węzłów są bezwzględne, jednoznaczne przy wielu korzeniach.
package core

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"danacoconsole/shared"
)

const (
	// glebokoscDomyslna wystarcza do nawigacji: korzeń i jego zawartość,
	// bez schodzenia głębiej bez żądania.
	glebokoscDomyslna = 2
	// glebokoscNajwieksza chroni przed zejściem w drzewo zależności, na przykład
	// katalog node_modules o niekończonej głębokości.
	glebokoscNajwieksza = 8
	// granicaWezlow jest największą liczbą pozycji jednej odpowiedzi,
	// powyżej której komenda odmawia zamiast obcinać wynik.
	granicaWezlow = 4000
)

// Drzewo obsługuje `developer.tree.get`, składając płaską listę węzłów
// w drzewo katalogów okna Project Tree.
func (a *adapterDevelopera) Drzewo(_ context.Context,
	z shared.DeveloperTreeGetRequest) (shared.DeveloperTreeGetResponse, error) {

	okno, err := a.oknoDevelopera(z.WindowId)
	if err != nil {
		return shared.DeveloperTreeGetResponse{}, err
	}
	korzenie := a.korzenieOkna(okno)
	if len(korzenie) == 0 {
		return shared.DeveloperTreeGetResponse{}, bladZasobuDevelopera(
			"okno nie ma ani jednego katalogu roboczego, więc drzewo projektu nie ma korzenia")
	}
	if wskazanie := strings.TrimSpace(wartoscTekstu(z.Path)); wskazanie != "" {
		sciezka, err := sciezkaWObszarze(korzenie, wskazanie, czyIstnieje)
		if err != nil {
			return shared.DeveloperTreeGetResponse{}, err
		}
		korzenie = []string{sciezka}
	}

	glebokosc := glebokoscZadania(z.Depth)
	ukryte := z.IncludeHidden != nil && *z.IncludeHidden

	wezly := make([]shared.DeveloperTreeNode, 0, 64)
	for _, korzen := range korzenie {
		opis, err := os.Stat(korzen)
		if err != nil {
			return shared.DeveloperTreeGetResponse{}, bladZasobuDevelopera(
				"katalog " + korzen + " nie jest dostępny: " + err.Error())
		}
		if !opis.IsDir() {
			return shared.DeveloperTreeGetResponse{}, bladZadaniaDevelopera(
				korzen + " nie jest katalogiem — drzewo zaczyna się od katalogu")
		}
		wezly = append(wezly, wezelDrzewa(korzen, "", opis))
		wezly, err = zbierzWezly(wezly, korzen, glebokosc-1, ukryte)
		if err != nil {
			return shared.DeveloperTreeGetResponse{}, err
		}
	}
	return shared.DeveloperTreeGetResponse{Root: korzenie[0], Nodes: wezly}, nil
}

// zbierzWezly dokłada zawartość katalogu do listy płaskiej. Pozycje idą
// katalogami przed plikami, alfabetycznie, niezależnie od kolejności zwróconej
// przez system plików.
func zbierzWezly(wezly []shared.DeveloperTreeNode, katalog string, glebokosc int,
	ukryte bool) ([]shared.DeveloperTreeNode, error) {

	if glebokosc <= 0 {
		return wezly, nil
	}
	pozycje, err := os.ReadDir(katalog)
	if err != nil {
		// Katalog bez prawa odczytu nie unieważnia całego drzewa: reszta
		// repozytorium ma się pokazać.
		return wezly, nil
	}
	sort.Slice(pozycje, func(i, j int) bool {
		if pozycje[i].IsDir() != pozycje[j].IsDir() {
			return pozycje[i].IsDir()
		}
		return strings.ToLower(pozycje[i].Name()) < strings.ToLower(pozycje[j].Name())
	})

	for _, pozycja := range pozycje {
		if !ukryte && strings.HasPrefix(pozycja.Name(), ".") {
			continue
		}
		if len(wezly) >= granicaWezlow {
			return nil, bladZadaniaDevelopera(
				"drzewo od " + katalog + " ma ponad " + itoa(granicaWezlow) +
					" pozycji — zawęź ścieżkę albo zmniejsz głębokość")
		}
		sciezka := filepath.Join(katalog, pozycja.Name())
		opis, err := pozycja.Info()
		if err != nil {
			continue
		}
		wezly = append(wezly, wezelDrzewa(sciezka, katalog, opis))
		if !pozycja.IsDir() {
			continue
		}
		wezly, err = zbierzWezly(wezly, sciezka, glebokosc-1, ukryte)
		if err != nil {
			return nil, err
		}
	}
	return wezly, nil
}

// wezelDrzewa składa jeden węzeł kontraktu. Rodzic pusty znaczy korzeń bez
// nadrzędnego katalogu roboczego.
func wezelDrzewa(sciezka, rodzic string, opis os.FileInfo) shared.DeveloperTreeNode {
	zmieniono := opis.ModTime().UTC().UnixMilli()
	wezel := shared.DeveloperTreeNode{
		Path:      sciezka,
		Name:      filepath.Base(sciezka),
		Kind:      shared.TreeNodeKindFile,
		UpdatedAt: &zmieniono,
	}
	if opis.IsDir() {
		wezel.Kind = shared.TreeNodeKindDirectory
	} else {
		rozmiar := opis.Size()
		wezel.SizeBytes = &rozmiar
	}
	if rodzic != "" {
		wezel.ParentPath = &rodzic
	}
	return wezel
}

// glebokoscZadania czyta głębokość z żądania. Wartość niedodatnia znaczy
// domyślną, a nadmiarowa jest przycinana do największej dopuszczalnej.
func glebokoscZadania(zadana *int) int {
	if zadana == nil || *zadana <= 0 {
		return glebokoscDomyslna
	}
	if *zadana > glebokoscNajwieksza {
		return glebokoscNajwieksza
	}
	return *zadana
}
