// Odpowiedzialność pliku: biblioteka projektu — okno Project Library modułu
// Workspace.
//
// Źródłem wykazu jest katalog roboczy projektu, nie osobna tabela: pliki
// powstają tam, gdzie pracuje model, a komenda wgrania pliku należy do modułu
// Library (`library.file.upload`), więc tabela plików w module Workspace nie
// miałaby pisarza.
//
// Project Library jest odpowiednikiem okna Library Explorer zawężonym do
// zakresu projektu, którym jest katalog roboczy projektu.
package core

import (
	"context"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"

	"danacoconsole/server/internal/konfig"
	"danacoconsole/shared"
)

// granicaPrzegladuKatalogu ogranicza liczbę wpisów przeglądanych w katalogu
// projektu, żeby odczyt biblioteki nie rósł wraz z katalogiem rozrosłym do
// dziesiątek tysięcy plików.
const granicaPrzegladuKatalogu = 5000

// Biblioteka zwraca pliki katalogu projektu, zawężone frazą wyszukiwania.
// Suma liczy pliki spełniające warunki przed przycięciem do limitu żądania,
// dzięki czemu okno rozpoznaje wykaz urwany.
func (a *adapterPrzestrzeniRoboczej) Biblioteka(_ context.Context,
	z shared.WorkspaceLibraryListRequest) (shared.WorkspaceLibraryListResponse, error) {

	if z.ProjectId == "" {
		return shared.WorkspaceLibraryListResponse{}, bladProjektu("komenda bez wskazania projektu")
	}
	fraza := ""
	if z.Query != nil {
		fraza = *z.Query
	}
	pliki := a.plikiProjektu(z.ProjectId, fraza)
	wszystkie := len(pliki)
	if z.Limit != nil && *z.Limit > 0 && len(pliki) > *z.Limit {
		pliki = pliki[:*z.Limit]
	}
	return shared.WorkspaceLibraryListResponse{Files: pliki, Total: &wszystkie}, nil
}

// plikiProjektu przegląda katalog roboczy projektu i składa wykaz plików
// kontraktu. Katalog niedostępny daje wykaz pusty — okno pokaże stan pusty,
// a nie błąd.
func (a *adapterPrzestrzeniRoboczej) plikiProjektu(idProjektu, fraza string) []shared.LibraryFile {
	pliki := []shared.LibraryFile{}
	korzen := a.katalogProjektu(idProjektu)
	if korzen == "" {
		return pliki
	}
	szukane := strings.ToLower(strings.TrimSpace(fraza))
	przejrzane := 0
	_ = filepath.WalkDir(korzen, func(sciezka string, wpis fs.DirEntry, err error) error {
		if err != nil {
			// Wpis nieczytelny pomijamy; jeden plik bez uprawnień nie przerywa
			// wykazu pozostałych.
			return nil
		}
		if przejrzane >= granicaPrzegladuKatalogu {
			return filepath.SkipAll
		}
		przejrzane++
		if wpis.IsDir() {
			return pomijanyKatalog(wpis.Name())
		}
		wzgledna, blad := filepath.Rel(korzen, sciezka)
		if blad != nil || strings.HasPrefix(wpis.Name(), ".") {
			return nil
		}
		if szukane != "" && !strings.Contains(strings.ToLower(wzgledna), szukane) {
			return nil
		}
		pliki = append(pliki, plikKontraktu(idProjektu, filepath.ToSlash(wzgledna), wpis))
		return nil
	})
	sort.Slice(pliki, func(i, j int) bool { return pliki[i].UpdatedAt > pliki[j].UpdatedAt })
	return pliki
}

// katalogProjektu ustala katalog roboczy projektu. Brak ustalacza daje ścieżkę
// pustą, a wykaz plików — pusty.
func (a *adapterPrzestrzeniRoboczej) katalogProjektu(idProjektu string) string {
	if a.katalog == nil {
		return ""
	}
	return a.katalog.Ustal(konfig.Kontekst{Projekt: idProjektu}, "projekt-"+idProjektu).Sciezka
}

// pomijanyKatalog odsiewa katalogi robocze narzędzi. Ich zawartość nie jest
// materiałem projektu, a potrafi przeważyć wykaz liczbą wpisów.
func pomijanyKatalog(nazwa string) error {
	if strings.HasPrefix(nazwa, ".") || nazwa == "node_modules" {
		return filepath.SkipDir
	}
	return nil
}

// plikKontraktu składa pozycję wykazu z wpisu katalogu. Identyfikatorem jest
// ścieżka względna — jest trwała w obrębie projektu i nie wymaga rejestru.
func plikKontraktu(idProjektu, wzgledna string, wpis fs.DirEntry) shared.LibraryFile {
	plik := shared.LibraryFile{
		Id:        wzgledna,
		Name:      wpis.Name(),
		Path:      &wzgledna,
		ProjectId: &idProjektu,
	}
	if rodzaj := rodzajTresci(wzgledna); rodzaj != "" {
		plik.MimeType = &rodzaj
	}
	opis, err := wpis.Info()
	if err != nil {
		return plik
	}
	rozmiar := opis.Size()
	zmieniony := opis.ModTime().UnixMilli()
	plik.SizeBytes = &rozmiar
	plik.CreatedAt, plik.UpdatedAt = zmieniony, zmieniony
	return plik
}

// rodzajTresci rozpoznaje rodzaj treści pliku po rozszerzeniu. Rozszerzenie
// spoza wykazu zostawia pole puste zamiast zgadywanego rodzaju.
func rodzajTresci(sciezka string) string {
	switch strings.ToLower(filepath.Ext(sciezka)) {
	case ".pdf":
		return "application/pdf"
	case ".docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case ".md":
		return "text/markdown"
	case ".txt":
		return "text/plain"
	case ".json":
		return "application/json"
	case ".csv":
		return "text/csv"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	}
	return ""
}
