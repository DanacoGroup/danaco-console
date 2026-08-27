// Plik obsługuje zakładki okna przeglądarki: `browser.bookmark.add`, `.list`, `.remove`. Zakładka jest czynnością świadomą Operatora, w odróżnieniu od źródła, które narasta samo w toku przeglądania.
package core

import (
	"context"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// DodajZakladke obsługuje `browser.bookmark.add`, zakładając wiersz zakładki wraz z folderem, etykietami i notatką.
func (a *adapterPrzegladarki) DodajZakladke(ctx context.Context,
	z shared.BrowserBookmarkAddRequest) (shared.BrowserBookmarkAddResponse, error) {

	if strings.TrimSpace(z.WindowId) == "" || strings.TrimSpace(z.Url) == "" {
		return shared.BrowserBookmarkAddResponse{}, bladWskazaniaPrzegladarki(
			"komenda bookmark.add bez okna albo adresu")
	}
	zakladka := dane.ZakladkaPrzegladania{
		Kod:          nowyIdentyfikator(przedrostekZakladki),
		Okno:         z.WindowId,
		Url:          strings.TrimSpace(z.Url),
		Tytul:        z.Title,
		Folder:       z.Folder,
		EtykietyJson: wykazJson(z.Tags),
		Notatka:      z.Note,
	}
	zapisana, err := a.repozytorium.ZapiszZakladke(ctx, zakladka)
	if err != nil {
		return shared.BrowserBookmarkAddResponse{}, bladPrzegladarki(err)
	}
	return shared.BrowserBookmarkAddResponse{Bookmark: zakladkaKontraktu(zapisana)}, nil
}

// WykazZakladek obsługuje `browser.bookmark.list`, przeszukując zakładki po adresie, tytule i notatce naraz.
func (a *adapterPrzegladarki) WykazZakladek(ctx context.Context,
	z shared.BrowserBookmarkListRequest) (shared.BrowserBookmarkListResponse, error) {

	wiersze, err := a.repozytorium.Zakladki(ctx, dane.FiltrZakladek{
		Okno:   wartoscTekstuLubPusta(z.WindowId),
		Folder: wartoscTekstuLubPusta(z.Folder),
		Szukaj: wartoscTekstuLubPusta(z.Query),
		Limit:  wartoscLiczby(z.Limit),
	})
	if err != nil {
		return shared.BrowserBookmarkListResponse{}, bladPrzegladarki(err)
	}
	zakladki := make([]shared.BrowserBookmark, 0, len(wiersze))
	for _, wiersz := range wiersze {
		zakladki = append(zakladki, zakladkaKontraktu(wiersz))
	}
	return shared.BrowserBookmarkListResponse{Bookmarks: zakladki}, nil
}

// UsunZakladke obsługuje `browser.bookmark.remove`, zdejmując wiersz zakładki wskazany identyfikatorem żądania.
func (a *adapterPrzegladarki) UsunZakladke(ctx context.Context,
	z shared.BrowserBookmarkRemoveRequest) (shared.BrowserBookmarkRemoveResponse, error) {

	if strings.TrimSpace(z.BookmarkId) == "" {
		return shared.BrowserBookmarkRemoveResponse{}, bladWskazaniaPrzegladarki(
			"komenda bookmark.remove bez zakładki")
	}
	usunieta, err := a.repozytorium.UsunZakladke(ctx, z.BookmarkId)
	if err != nil {
		return shared.BrowserBookmarkRemoveResponse{}, bladPrzegladarki(err)
	}
	if !usunieta {
		return shared.BrowserBookmarkRemoveResponse{}, bladNieznanegoBytu("zakładki", z.BookmarkId)
	}
	return shared.BrowserBookmarkRemoveResponse{Removed: true}, nil
}

// zakladkaKontraktu przekłada wiersz zakładki odczytany z repozytorium na byt kontraktu odpowiedzi okna.
func zakladkaKontraktu(w dane.ZakladkaPrzegladania) shared.BrowserBookmark {
	zakladka := shared.BrowserBookmark{
		Id:        w.Kod,
		WindowId:  w.Okno,
		Url:       w.Url,
		Title:     w.Tytul,
		Folder:    w.Folder,
		Note:      w.Notatka,
		CreatedAt: chwilaBazy(w.Utworzono),
	}
	if etykiety := wykazZJson(w.EtykietyJson); len(etykiety) > 0 {
		zakladka.Tags = etykiety
	}
	return zakladka
}
