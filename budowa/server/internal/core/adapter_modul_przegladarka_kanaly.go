// Odpowiedzialność pliku: kanały RSS/Atom/JSON Feed i kolejka czytania —
// `browser.feed.subscribe`, `.list`, `.remove` oraz `browser.readlist.add`,
// `.list`, `.remove`. Subskrypcja naprawdę odpytuje kanał i rozpoznaje jego
// postać przed zapisem.
package core

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// kanalRss jest kształtem dokumentu RSS 2.0 w zakresie, którego używamy:
// kanał, jego wpisy i pola potrzebne do listy oraz odczytu.
type kanalRss struct {
	Kanal struct {
		Tytul string `xml:"title"`
		Wpisy []struct {
			Tytul string `xml:"title"`
			Adres string `xml:"link"`
			Opis  string `xml:"description"`
			Data  string `xml:"pubDate"`
			Klucz string `xml:"guid"`
		} `xml:"item"`
	} `xml:"channel"`
}

// kanalAtom jest kształtem dokumentu Atom w zakresie, którego używamy:
// kanał, jego wpisy i pola potrzebne do listy oraz odczytu.
type kanalAtom struct {
	Tytul string `xml:"title"`
	Wpisy []struct {
		Tytul string `xml:"title"`
		Adres []struct {
			Href string `xml:"href,attr"`
			Rel  string `xml:"rel,attr"`
		} `xml:"link"`
		Streszczenie string `xml:"summary"`
		Tresc        string `xml:"content"`
		Data         string `xml:"updated"`
		Klucz        string `xml:"id"`
	} `xml:"entry"`
}

// kanalJson jest kształtem dokumentu JSON Feed w zakresie, którego używamy:
// kanał, jego wpisy i pola potrzebne do listy oraz odczytu.
type kanalJson struct {
	Tytul string `json:"title"`
	Wpisy []struct {
		Klucz        string `json:"id"`
		Adres        string `json:"url"`
		Tytul        string `json:"title"`
		Streszczenie string `json:"summary"`
		Tresc        string `json:"content_text"`
		Data         string `json:"date_published"`
	} `json:"items"`
}

// wpisKanaluOdczytany jest jednym wpisem po rozbiorze, niezależnie od postaci
// dokumentu — trzy formaty, jeden kształt dalej.
type wpisKanaluOdczytany struct {
	Adres        string
	Tytul        string
	Streszczenie string
	Opublikowano *string
}

// SubskrybujKanal obsługuje `browser.feed.subscribe`: pobiera wskazany
// adres, rozpoznaje postać dokumentu i odkłada wpisy.
func (a *adapterPrzegladarki) SubskrybujKanal(ctx context.Context,
	z shared.BrowserFeedSubscribeRequest) (shared.BrowserFeedSubscribeResponse, error) {

	if strings.TrimSpace(z.WindowId) == "" || strings.TrimSpace(z.Url) == "" {
		return shared.BrowserFeedSubscribeResponse{}, bladWskazaniaPrzegladarki(
			"komenda feed.subscribe bez okna albo adresu")
	}
	adres := strings.TrimSpace(z.Url)
	tresc, err := pobierzStrone(ctx, adres)
	if err != nil {
		return shared.BrowserFeedSubscribeResponse{}, bladPobraniaStrony(adres, err)
	}

	postac, tytulKanalu, wpisy, err := rozbierzKanal(tresc.Html)
	if err != nil {
		return shared.BrowserFeedSubscribeResponse{}, bladWskazaniaPrzegladarki(
			"pod adresem " + adres + " nie stoi kanał RSS, Atom ani JSON Feed: " + err.Error())
	}

	kanal := dane.KanalPrzegladania{
		Kod:  nowyIdentyfikator(przedrostekKanaluPrzegladania),
		Okno: z.WindowId,
		Url:  adres,
	}
	tytul := wartoscTekstuLubPusta(z.Title)
	if tytul == "" {
		tytul = tytulKanalu
	}
	if tytul != "" {
		kanal.Tytul = &tytul
	}
	if postac != "" {
		zapis := postac
		kanal.Postac = &zapis
	}
	if interwal := int64(wartoscLiczby(z.IntervalSeconds)); interwal > 0 {
		kanal.InterwalSekund = &interwal
	}
	kanal.Pobrano = terazWBazie()

	zapisany, err := a.repozytorium.ZapiszKanal(ctx, kanal)
	if err != nil {
		return shared.BrowserFeedSubscribeResponse{}, bladPrzegladarki(err)
	}
	for _, wpis := range wpisy {
		zapis := dane.WpisKanalu{
			Kod:          nowyIdentyfikator(przedrostekWpisuKanalu),
			Kanal:        zapisany.Kod,
			Url:          wpis.Adres,
			Opublikowano: wpis.Opublikowano,
		}
		if wpis.Tytul != "" {
			tytul := wpis.Tytul
			zapis.Tytul = &tytul
		}
		if wpis.Streszczenie != "" {
			streszczenie := wpis.Streszczenie
			zapis.Streszczenie = &streszczenie
		}
		if err := a.repozytorium.ZapiszWpisKanalu(ctx, zapis); err != nil {
			return shared.BrowserFeedSubscribeResponse{}, bladPrzegladarki(err)
		}
	}

	oddany, err := a.kanalKontraktu(ctx, zapisany, true, false, 0)
	if err != nil {
		return shared.BrowserFeedSubscribeResponse{}, err
	}
	return shared.BrowserFeedSubscribeResponse{Feed: oddany}, nil
}

// WykazKanalow obsługuje `browser.feed.list`: oddaje subskrypcje Operatora,
// opcjonalnie z wpisami zawężonymi do nieprzeczytanych.
func (a *adapterPrzegladarki) WykazKanalow(ctx context.Context,
	z shared.BrowserFeedListRequest) (shared.BrowserFeedListResponse, error) {

	wiersze, err := a.repozytorium.Kanaly(ctx, wartoscTekstuLubPusta(z.WindowId),
		wartoscTekstuLubPusta(z.FeedId), wartoscLiczby(z.Limit))
	if err != nil {
		return shared.BrowserFeedListResponse{}, bladPrzegladarki(err)
	}
	zWpisami := z.IncludeEntries != nil && *z.IncludeEntries
	tylkoNieprzeczytane := z.UnreadOnly != nil && *z.UnreadOnly
	kanaly := make([]shared.BrowserFeed, 0, len(wiersze))
	for _, wiersz := range wiersze {
		kanal, err := a.kanalKontraktu(ctx, wiersz, zWpisami, tylkoNieprzeczytane, wartoscLiczby(z.Limit))
		if err != nil {
			return shared.BrowserFeedListResponse{}, err
		}
		kanaly = append(kanaly, kanal)
	}
	return shared.BrowserFeedListResponse{Feeds: kanaly}, nil
}

// ZdejmijKanal obsługuje `browser.feed.remove` — subskrypcja znika wraz
// z wpisami (kaskada schematu, migracja 173).
func (a *adapterPrzegladarki) ZdejmijKanal(ctx context.Context,
	z shared.BrowserFeedRemoveRequest) (shared.BrowserFeedRemoveResponse, error) {

	if strings.TrimSpace(z.FeedId) == "" {
		return shared.BrowserFeedRemoveResponse{}, bladWskazaniaPrzegladarki("komenda feed.remove bez kanału")
	}
	usuniety, err := a.repozytorium.UsunKanal(ctx, z.FeedId)
	if err != nil {
		return shared.BrowserFeedRemoveResponse{}, bladPrzegladarki(err)
	}
	if !usuniety {
		return shared.BrowserFeedRemoveResponse{}, bladNieznanegoBytu("kanału", z.FeedId)
	}
	return shared.BrowserFeedRemoveResponse{Removed: true}, nil
}

// OdlozDoCzytania obsługuje `browser.readlist.add`: dokłada wskazaną stronę
// do kolejki czytania tego okna.
func (a *adapterPrzegladarki) OdlozDoCzytania(ctx context.Context,
	z shared.BrowserReadlistAddRequest) (shared.BrowserReadlistAddResponse, error) {

	if strings.TrimSpace(z.WindowId) == "" || strings.TrimSpace(z.Url) == "" {
		return shared.BrowserReadlistAddResponse{}, bladWskazaniaPrzegladarki(
			"komenda readlist.add bez okna albo adresu")
	}
	pozycja := dane.PozycjaCzytania{
		Kod:     nowyIdentyfikator(przedrostekPozycjiCzytania),
		Okno:    z.WindowId,
		Url:     strings.TrimSpace(z.Url),
		Tytul:   z.Title,
		Notatka: z.Note,
	}
	if z.RemindAt != nil {
		pozycja.Przypomnienie = znacznikChwili(*z.RemindAt)
	}
	zapisana, err := a.repozytorium.ZapiszPozycjeCzytania(ctx, pozycja)
	if err != nil {
		return shared.BrowserReadlistAddResponse{}, bladPrzegladarki(err)
	}
	return shared.BrowserReadlistAddResponse{Item: pozycjaCzytaniaKontraktu(zapisana)}, nil
}

// WykazCzytania obsługuje `browser.readlist.list`: oddaje pozycje kolejki
// czytania, opcjonalnie zawężone do nieprzeczytanych.
func (a *adapterPrzegladarki) WykazCzytania(ctx context.Context,
	z shared.BrowserReadlistListRequest) (shared.BrowserReadlistListResponse, error) {

	tylkoNieprzeczytane := z.PendingOnly != nil && *z.PendingOnly
	wiersze, err := a.repozytorium.KolejkaCzytania(ctx, wartoscTekstuLubPusta(z.WindowId),
		tylkoNieprzeczytane, wartoscLiczby(z.Limit))
	if err != nil {
		return shared.BrowserReadlistListResponse{}, bladPrzegladarki(err)
	}
	pozycje := make([]shared.BrowserReadingItem, 0, len(wiersze))
	for _, wiersz := range wiersze {
		pozycje = append(pozycje, pozycjaCzytaniaKontraktu(wiersz))
	}
	return shared.BrowserReadlistListResponse{Items: pozycje}, nil
}

// ZdejmijZCzytania obsługuje `browser.readlist.remove` — zdjęcie pozycji albo
// oznaczenie jej jako przeczytanej. Dwie różne czynności w jednej komendzie
// rozstrzyga pole `markRead`: pozycja oznaczona zostaje w kolejce jako ślad
// przeczytania.
func (a *adapterPrzegladarki) ZdejmijZCzytania(ctx context.Context,
	z shared.BrowserReadlistRemoveRequest) (shared.BrowserReadlistRemoveResponse, error) {

	if strings.TrimSpace(z.ItemId) == "" {
		return shared.BrowserReadlistRemoveResponse{}, bladWskazaniaPrzegladarki(
			"komenda readlist.remove bez pozycji")
	}
	pozycja, err := a.repozytorium.PozycjaCzytania(ctx, z.ItemId)
	if err != nil {
		return shared.BrowserReadlistRemoveResponse{}, bladWierszaPrzegladania("pozycji kolejki czytania", z.ItemId, err)
	}
	if z.MarkRead != nil && *z.MarkRead {
		pozycja.Przeczytana = true
		zapisana, err := a.repozytorium.ZapiszPozycjeCzytania(ctx, pozycja)
		if err != nil {
			return shared.BrowserReadlistRemoveResponse{}, bladPrzegladarki(err)
		}
		oddana := pozycjaCzytaniaKontraktu(zapisana)
		return shared.BrowserReadlistRemoveResponse{Item: &oddana, Removed: false}, nil
	}
	usunieta, err := a.repozytorium.UsunPozycjeCzytania(ctx, z.ItemId)
	if err != nil {
		return shared.BrowserReadlistRemoveResponse{}, bladPrzegladarki(err)
	}
	return shared.BrowserReadlistRemoveResponse{Removed: usunieta}, nil
}

// kanalKontraktu składa `BrowserFeed`, z wpisami albo bez nich, zależnie od
// tego, czy wywołanie ich zażądało.
func (a *adapterPrzegladarki) kanalKontraktu(ctx context.Context, w dane.KanalPrzegladania,
	zWpisami, tylkoNieprzeczytane bool, limit int) (shared.BrowserFeed, error) {

	kanal := shared.BrowserFeed{
		Id:            w.Kod,
		WindowId:      w.Okno,
		Url:           w.Url,
		Title:         w.Tytul,
		LastFetchedAt: chwilaZeZnacznika(w.Pobrano),
		CreatedAt:     chwilaBazy(w.Utworzono),
	}
	if w.Postac != nil {
		postac := shared.BrowserFeedFormat(*w.Postac)
		kanal.Format = &postac
	}
	if w.InterwalSekund != nil {
		interwal := int(*w.InterwalSekund)
		kanal.IntervalSeconds = &interwal
	}
	nieprzeczytane, err := a.repozytorium.NieprzeczytaneKanalu(ctx, w.Kod)
	if err != nil {
		return shared.BrowserFeed{}, bladPrzegladarki(err)
	}
	liczba := int(nieprzeczytane)
	kanal.UnreadCount = &liczba

	if !zWpisami {
		return kanal, nil
	}
	wiersze, err := a.repozytorium.WpisyKanalu(ctx, w.Kod, tylkoNieprzeczytane, limit)
	if err != nil {
		return shared.BrowserFeed{}, bladPrzegladarki(err)
	}
	wpisy := make([]shared.BrowserFeedEntry, 0, len(wiersze))
	for _, wiersz := range wiersze {
		wpisy = append(wpisy, shared.BrowserFeedEntry{
			Id:          wiersz.Kod,
			FeedId:      wiersz.Kanal,
			Url:         wiersz.Url,
			Title:       wiersz.Tytul,
			Summary:     wiersz.Streszczenie,
			Read:        wiersz.Przeczytany,
			PublishedAt: chwilaZeZnacznika(wiersz.Opublikowano),
		})
	}
	kanal.Entries = wpisy
	return kanal, nil
}

// pozycjaCzytaniaKontraktu przekłada wiersz kolejki czytania na pozycję
// kontraktu `BrowserReadingItem`.
func pozycjaCzytaniaKontraktu(w dane.PozycjaCzytania) shared.BrowserReadingItem {
	return shared.BrowserReadingItem{
		Id:        w.Kod,
		WindowId:  w.Okno,
		Url:       w.Url,
		Title:     w.Tytul,
		Note:      w.Notatka,
		Read:      w.Przeczytana,
		RemindAt:  chwilaZeZnacznika(w.Przypomnienie),
		CreatedAt: chwilaBazy(w.Utworzono),
	}
}

// rozbierzKanal rozpoznaje postać dokumentu i wyjmuje z niego wpisy.
// Kolejność prób idzie od najczęstszej postaci; dokument, którego nie rozpoznaje
// żadna z trzech, jest odmową nazywającą powód, a nie pustą subskrypcją.
func rozbierzKanal(dokument string) (string, string, []wpisKanaluOdczytany, error) {
	przyciety := strings.TrimSpace(dokument)
	if przyciety == "" {
		return "", "", nil, errIntoKanal("dokument jest pusty")
	}

	if strings.HasPrefix(przyciety, "{") {
		var kanal kanalJson
		if err := json.Unmarshal([]byte(przyciety), &kanal); err == nil && len(kanal.Wpisy) > 0 {
			wpisy := make([]wpisKanaluOdczytany, 0, len(kanal.Wpisy))
			for _, wpis := range kanal.Wpisy {
				adres := wpis.Adres
				if adres == "" {
					adres = wpis.Klucz
				}
				streszczenie := wpis.Streszczenie
				if streszczenie == "" {
					streszczenie = wpis.Tresc
				}
				wpisy = append(wpisy, wpisKanaluOdczytany{
					Adres: adres, Tytul: wpis.Tytul, Streszczenie: streszczenie,
					Opublikowano: znacznikDaty(wpis.Data),
				})
			}
			return string(shared.BrowserFeedFormatJsonFeed), kanal.Tytul, wpisy, nil
		}
	}

	var rss kanalRss
	if err := xml.Unmarshal([]byte(przyciety), &rss); err == nil && len(rss.Kanal.Wpisy) > 0 {
		wpisy := make([]wpisKanaluOdczytany, 0, len(rss.Kanal.Wpisy))
		for _, wpis := range rss.Kanal.Wpisy {
			adres := strings.TrimSpace(wpis.Adres)
			if adres == "" {
				adres = strings.TrimSpace(wpis.Klucz)
			}
			wpisy = append(wpisy, wpisKanaluOdczytany{
				Adres: adres, Tytul: strings.TrimSpace(wpis.Tytul),
				Streszczenie: strings.TrimSpace(wpis.Opis),
				Opublikowano: znacznikDaty(wpis.Data),
			})
		}
		return string(shared.BrowserFeedFormatRss), strings.TrimSpace(rss.Kanal.Tytul), wpisy, nil
	}

	var atom kanalAtom
	if err := xml.Unmarshal([]byte(przyciety), &atom); err == nil && len(atom.Wpisy) > 0 {
		wpisy := make([]wpisKanaluOdczytany, 0, len(atom.Wpisy))
		for _, wpis := range atom.Wpisy {
			adres := ""
			for _, odnosnik := range wpis.Adres {
				if odnosnik.Rel == "" || odnosnik.Rel == "alternate" {
					adres = odnosnik.Href
					break
				}
			}
			if adres == "" {
				adres = strings.TrimSpace(wpis.Klucz)
			}
			streszczenie := strings.TrimSpace(wpis.Streszczenie)
			if streszczenie == "" {
				streszczenie = strings.TrimSpace(wpis.Tresc)
			}
			wpisy = append(wpisy, wpisKanaluOdczytany{
				Adres: adres, Tytul: strings.TrimSpace(wpis.Tytul),
				Streszczenie: streszczenie, Opublikowano: znacznikDaty(wpis.Data),
			})
		}
		return string(shared.BrowserFeedFormatAtom), strings.TrimSpace(atom.Tytul), wpisy, nil
	}

	return "", "", nil, errIntoKanal("dokument nie jest kanałem RSS, Atom ani JSON Feed albo nie ma ani jednego wpisu")
}

// znacznikDaty przekłada datę publikacji wpisu na znacznik kolumny czasu.
// Formatów jest kilka, bo kanały ich używają kilku; data nierozpoznana daje
// brak, a nie datę zmyśloną.
func znacznikDaty(zapis string) *string {
	przyciety := strings.TrimSpace(zapis)
	if przyciety == "" {
		return nil
	}
	formaty := []string{
		time.RFC1123Z, time.RFC1123, time.RFC3339, time.RFC822Z, time.RFC822,
		"2006-01-02T15:04:05Z07:00", "2006-01-02 15:04:05", "2006-01-02",
	}
	for _, format := range formaty {
		if chwila, err := time.Parse(format, przyciety); err == nil {
			znak := chwila.UTC().Format(formatZnacznikaBazy)
			return &znak
		}
	}
	return nil
}

// errIntoKanal opakowuje podany powód niepowodzenia rozbioru kanału w błąd
// zgodny z interfejsem error.
func errIntoKanal(powod string) error {
	return &bladRozbioruKanalu{powod: powod}
}

// bladRozbioruKanalu nazywa powód, dla którego dokument nie jest kanałem
// żadnej ze znanych postaci — RSS, Atom ani JSON Feed.
type bladRozbioruKanalu struct{ powod string }

func (b *bladRozbioruKanalu) Error() string { return b.powod }
