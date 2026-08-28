// Plik wypełnia port NarzedziaSesji. Dołożenie dodaje narzędzie do sesji na
// czas jej trwania, a wykaz łączy komendy kontraktu z katalogiem rozszerzeń
// w jedną filtrowaną listę pozycji po ukośniku.
package core

import (
	"context"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// adapterNarzedziSesji wypełnia port NarzedziaSesji, opierając się na
// repozytorium dołożeń sesji, repozytorium sesji oraz repozytorium rozszerzeń
// dostępnych do dołożenia.
type adapterNarzedziSesji struct {
	dolozenia    dane.RepozytoriumNarzedziSesji
	sesje        dane.RepozytoriumSesji
	rozszerzenia dane.RepozytoriumRozszerzen
}

func nowyAdapterNarzedziSesji(dolozenia dane.RepozytoriumNarzedziSesji,
	sesje dane.RepozytoriumSesji, rozszerzenia dane.RepozytoriumRozszerzen) *adapterNarzedziSesji {
	return &adapterNarzedziSesji{dolozenia: dolozenia, sesje: sesje, rozszerzenia: rozszerzenia}
}

// Doloz dokłada narzędzie do sesji na czas jej trwania i oddaje zapisane
// dołożenie razem z pełnym zestawem narzędzi sesji po tej czynności.
func (a *adapterNarzedziSesji) Doloz(ctx context.Context,
	z shared.SessionToolAttachRequest) (shared.SessionToolAttachResponse, error) {

	sesja, err := a.sesjaZadania(ctx, z.SessionId)
	if err != nil {
		return shared.SessionToolAttachResponse{}, err
	}
	pozycja, err := a.pozycjaPoNazwie(ctx, z.ToolName)
	if err != nil {
		return shared.SessionToolAttachResponse{}, err
	}
	if !pozycja.Attachable {
		return shared.SessionToolAttachResponse{}, bladNarzedziaSesji(niedostepnaPozycja(pozycja))
	}
	zrodlo := shared.SessionToolSourceSlashCommand
	if z.Source != nil && strings.TrimSpace(string(*z.Source)) != "" {
		zrodlo = string(*z.Source)
	}
	zapisane, juzBylo, err := a.dolozenia.Doloz(ctx, sesja, dane.NarzedzieSesji{
		NazwaPelna:    pozycja.Name,
		NazwaSkrocona: pozycja.ShortName,
		Opis:          pozycja.Description,
		Rodzaj:        string(pozycja.Kind),
		Grupa:         pozycja.Group,
		ZrodloPozycji: pozycja.Origin,
		Zrodlo:        zrodlo,
		Dolozono:      time.Now().UnixMilli(),
	})
	if err != nil {
		return shared.SessionToolAttachResponse{}, bladNosnikaNarzedziSesji(err)
	}
	zestaw, err := a.zestawSesji(ctx, sesja)
	if err != nil {
		return shared.SessionToolAttachResponse{}, err
	}
	return shared.SessionToolAttachResponse{
		Tool:            narzedzieKontraktu(zapisane),
		Tools:           zestaw,
		AlreadyAttached: juzBylo,
	}, nil
}

// Zdejmij zdejmuje dołożenie z sesji. Brak dołożenia nie jest odmową, bo
// żądany stan już zachodzi. Druga zwracana wartość niesie pozycje faktycznie
// zdjęte, nazwane pełną nazwą, do rozgłoszenia zdarzenia.
func (a *adapterNarzedziSesji) Zdejmij(ctx context.Context,
	z shared.SessionToolDetachRequest) (shared.SessionToolDetachResponse, []shared.SessionTool, error) {

	sesja, err := a.sesjaZadania(ctx, z.SessionId)
	if err != nil {
		return shared.SessionToolDetachResponse{}, nil, err
	}
	nazwa := strings.TrimSpace(z.ToolName)
	if nazwa == "" {
		return shared.SessionToolDetachResponse{}, nil, bladNarzedziaSesji("żądanie zdjęcia " +
			"dołożenia bez wskazania narzędzia; Operator poda `toolName` — nazwę pełną " +
			"ze źródłem albo skróconą")
	}
	// Zdejmowanie idzie po dołożeniach sesji, nie po wykazie: wiersz zostaje
	// mimo odinstalowania.
	zastane, err := a.dolozenia.Narzedzia(ctx, sesja)
	if err != nil {
		return shared.SessionToolDetachResponse{}, nil, bladNosnikaNarzedziSesji(err)
	}
	zdjete := make([]shared.SessionTool, 0, 1)
	for _, wpis := range zastane {
		if wpis.NazwaPelna != nazwa && wpis.NazwaSkrocona != nazwa {
			continue
		}
		usuniete, err := a.dolozenia.Zdejmij(ctx, sesja, wpis.NazwaPelna)
		if err != nil {
			return shared.SessionToolDetachResponse{}, nil, bladNosnikaNarzedziSesji(err)
		}
		// Do rozgłoszenia wchodzi wyłącznie wiersz faktycznie zdjęty, nie
		// każda próba zdjęcia.
		if usuniete {
			zdjete = append(zdjete, narzedzieKontraktu(wpis))
		}
	}
	zestaw, err := a.zestawSesji(ctx, sesja)
	if err != nil {
		return shared.SessionToolDetachResponse{}, nil, err
	}
	return shared.SessionToolDetachResponse{Detached: len(zdjete) > 0, Tools: zestaw}, zdjete, nil
}

// Wykaz oddaje dołożenia bieżącej sesji przełożone na kształt kontraktu wraz
// z ich pełną, całkowitą liczbą.
func (a *adapterNarzedziSesji) Wykaz(ctx context.Context,
	z shared.SessionToolListRequest) (shared.SessionToolListResponse, error) {

	sesja, err := a.sesjaZadania(ctx, z.SessionId)
	if err != nil {
		return shared.SessionToolListResponse{}, err
	}
	zestaw, err := a.zestawSesji(ctx, sesja)
	if err != nil {
		return shared.SessionToolListResponse{}, err
	}
	return shared.SessionToolListResponse{Tools: zestaw, Total: len(zestaw)}, nil
}

// Katalog oddaje wykaz pozycji po ukośniku po zawężeniu, oznaczeniu
// dołożonych i uporządkowaniu, gotowy do stronicowania.
func (a *adapterNarzedziSesji) Katalog(ctx context.Context,
	z shared.ToolsCatalogListRequest) (shared.ToolsCatalogListResponse, error) {

	pozycje, err := a.pozycjeWykazu(ctx)
	if err != nil {
		return shared.ToolsCatalogListResponse{}, err
	}
	if z.SessionId != nil && strings.TrimSpace(*z.SessionId) != "" {
		if err := a.oznaczDolozone(ctx, *z.SessionId, pozycje); err != nil {
			return shared.ToolsCatalogListResponse{}, err
		}
	}
	wybrane := zawezPozycje(pozycje, z)
	uporzadkujPozycje(wybrane, tekstZawezenia(z.Query))
	return shared.ToolsCatalogListResponse{
		Entries: przytnijPozycje(wybrane, z.Limit, z.Offset),
		Total:   len(wybrane),
		Groups:  grupyPozycji(wybrane),
	}, nil
}

// ── składanie wykazu ─────────────────────────────────────────────────────────

// pozycjeWykazu składa wykaz z dwóch żywych źródeł: komend kontraktu modelu
// oraz zainstalowanego katalogu rozszerzeń, bez trzeciego, osobnego katalogu
// akcji.
func (a *adapterNarzedziSesji) pozycjeWykazu(ctx context.Context) ([]*shared.ToolCatalogEntry, error) {
	deklaracje := shared.NarzedziaModelu()
	pozycje := make([]*shared.ToolCatalogEntry, 0, len(deklaracje)+32)
	for _, deklaracja := range deklaracje {
		komenda := deklaracja.Command
		obszar, _, _ := strings.Cut(string(komenda), shared.SeparatorObszaru)
		pozycje = append(pozycje, &shared.ToolCatalogEntry{
			Name:        zrodloPlatformy + ":" + string(komenda),
			ShortName:   string(komenda),
			Description: deklaracja.Description,
			Kind:        shared.SlashEntryKindAction,
			Group:       obszar,
			Origin:      zrodloPlatformy,
			Command:     &komenda,
			// Komenda akcji wykonuje czynność i zestawu narzędzi nie zmienia.
			Attachable: false,
		})
	}
	if a.rozszerzenia == nil {
		return pozycje, nil
	}
	wiersze, err := a.rozszerzenia.Rozszerzenia(ctx, dane.FiltrRozszerzen{})
	if err != nil {
		return nil, bladNosnikaNarzedziSesji(err)
	}
	for _, wiersz := range wiersze {
		pozycje = append(pozycje, pozycjaRozszerzenia(wiersz))
	}
	return pozycje, nil
}

// pozycjaRozszerzenia przekłada wiersz katalogu rozszerzeń na pozycję wykazu.
// Przedrostek źródła bierze z kodu, gdy kod go niesie, a w pozostałych
// wypadkach z pochodzenia pozycji.
func pozycjaRozszerzenia(w dane.Rozszerzenie) *shared.ToolCatalogEntry {
	zrodlo, skrocona, zDwukropkiem := strings.Cut(w.Kod, ":")
	pelna := w.Kod
	if !zDwukropkiem {
		zrodlo, skrocona = w.ZrodloPochodzenia, w.Kod
		pelna = zrodlo + ":" + skrocona
	}
	opis := w.Nazwa
	if w.Opis != nil && strings.TrimSpace(*w.Opis) != "" {
		opis = strings.TrimSpace(*w.Opis)
	}
	return &shared.ToolCatalogEntry{
		Name:        pelna,
		ShortName:   skrocona,
		Description: opis,
		Kind:        shared.SlashEntryKindTool,
		Group:       w.Rodzaj,
		Origin:      zrodlo,
		// Dołożyć można wyłącznie pozycję zainstalowaną i włączoną; pozostałe
		// zostają widoczne, niedołożalne.
		Attachable: w.Zainstalowane && w.Wlaczone,
	}
}

// pozycjaPoNazwie odszukuje pozycję wykazu po nazwie pełnej albo skróconej.
// Nazwa skrócona trafiająca w kilka pozycji jest ODMOWĄ, nie wyborem pierwszej
// z brzegu: dołożone zostałoby wtedy narzędzie, o które Operator nie prosił.
func (a *adapterNarzedziSesji) pozycjaPoNazwie(ctx context.Context,
	nazwa string) (shared.ToolCatalogEntry, error) {

	szukana := strings.TrimSpace(nazwa)
	if szukana == "" {
		return shared.ToolCatalogEntry{}, bladNarzedziaSesji("żądanie dołożenia bez " +
			"wskazania narzędzia; Operator poda `toolName` — nazwę pełną ze źródłem " +
			"albo skróconą, tę samą, którą niesie wykaz po ukośniku")
	}
	pozycje, err := a.pozycjeWykazu(ctx)
	if err != nil {
		return shared.ToolCatalogEntry{}, err
	}
	trafione := make([]*shared.ToolCatalogEntry, 0, 2)
	for _, pozycja := range pozycje {
		if pozycja.Name == szukana {
			return *pozycja, nil
		}
		if pozycja.ShortName == szukana {
			trafione = append(trafione, pozycja)
		}
	}
	switch len(trafione) {
	case 1:
		return *trafione[0], nil
	case 0:
		return shared.ToolCatalogEntry{}, bladNarzedziaSesji("pozycji " + szukana +
			" nie ma w wykazie po ukośniku; Operator sprawdzi nazwę komendą " +
			"tools.catalog.list — wykaz zna nazwę pełną ze źródłem i skróconą")
	default:
		return shared.ToolCatalogEntry{}, bladNarzedziaSesji("nazwa skrócona " + szukana +
			" trafia w " + nazwyTrafionych(trafione) + "; Operator poda nazwę pełną " +
			"ze źródłem — wybór pierwszej z brzegu dołożyłby nie to narzędzie")
	}
}

// ── stan zestawu sesji ───────────────────────────────────────────────────────

// sesjaZadania przekłada identyfikator kontraktowy sesji na klucz jej wiersza
// w repozytorium sesji rdzenia.
func (a *adapterNarzedziSesji) sesjaZadania(ctx context.Context, wskazanie string) (int64, error) {
	// Straż odbiornika zerowego, bo metoda jest wołana ze wszystkich czynności
	// portu.
	if a == nil {
		return 0, bladNosnikaNarzedziSesji(nil)
	}
	identyfikator := strings.TrimSpace(wskazanie)
	if identyfikator == "" {
		return 0, bladNarzedziaSesji("żądanie bez wskazania sesji; Operator poda " +
			"`sessionId` — dołożenie żyje w stanie sesji i bez niej nie ma gdzie zamieszkać")
	}
	if a.sesje == nil {
		return 0, bladNosnikaNarzedziSesji(nil)
	}
	wiersz, err := a.sesje.PoIdentyfikatorze(ctx, identyfikator)
	if err != nil {
		if err == dane.ErrBrakWiersza {
			return 0, bladNarzedziaSesji("sesji " + identyfikator + " nie ma w historii; " +
				"Operator wskaże sesję istniejącą — dołożenie kończy się razem z sesją, " +
				"więc sesji nieistniejącej nie ma czym poszerzyć")
		}
		return 0, bladNosnikaNarzedziSesji(err)
	}
	return wiersz.ID, nil
}

// zestawSesji oddaje dołożenia sesji przełożone na kształt kontraktu, gotowy
// do zwrócenia w odpowiedzi.
func (a *adapterNarzedziSesji) zestawSesji(ctx context.Context, sesja int64) ([]shared.SessionTool, error) {
	if a.dolozenia == nil {
		return nil, bladNosnikaNarzedziSesji(nil)
	}
	wiersze, err := a.dolozenia.Narzedzia(ctx, sesja)
	if err != nil {
		return nil, bladNosnikaNarzedziSesji(err)
	}
	zestaw := make([]shared.SessionTool, 0, len(wiersze))
	for _, wiersz := range wiersze {
		zestaw = append(zestaw, narzedzieKontraktu(wiersz))
	}
	return zestaw, nil
}

// DolozeniaNarzedzi wypełnia port DolozeniaNarzedziSesji: oddaje same nazwy
// dołożeń sesji w kolejności dokładania, na potrzeby składania zestawu
// narzędzi tury modelu. Sesja bez dołożeń oddaje listę pustą.
func (a *adapterNarzedziSesji) DolozeniaNarzedzi(ctx context.Context, idSesji string) ([]string, error) {
	// Odbiornik zerowy nie jest przypadkiem niemożliwym: przechodzi
	// porównanie u wołającego niezauważony.
	if a == nil {
		return nil, bladNosnikaNarzedziSesji(nil)
	}
	sesja, err := a.sesjaZadania(ctx, idSesji)
	if err != nil {
		return nil, err
	}
	if a.dolozenia == nil {
		return nil, bladNosnikaNarzedziSesji(nil)
	}
	wiersze, err := a.dolozenia.Narzedzia(ctx, sesja)
	if err != nil {
		return nil, bladNosnikaNarzedziSesji(err)
	}
	nazwy := make([]string, 0, len(wiersze))
	for _, wiersz := range wiersze {
		nazwy = append(nazwy, wiersz.NazwaPelna)
	}
	return nazwy, nil
}

// oznaczDolozone zaznacza w przekazanym wykazie pozycje już dołożone do
// wskazanej sesji, wypełniając pole obecności przy każdej pozycji.
func (a *adapterNarzedziSesji) oznaczDolozone(ctx context.Context, wskazanie string,
	pozycje []*shared.ToolCatalogEntry) error {

	sesja, err := a.sesjaZadania(ctx, wskazanie)
	if err != nil {
		return err
	}
	wiersze, err := a.dolozenia.Narzedzia(ctx, sesja)
	if err != nil {
		return bladNosnikaNarzedziSesji(err)
	}
	dolozone := make(map[string]bool, len(wiersze))
	for _, wiersz := range wiersze {
		dolozone[wiersz.NazwaPelna] = true
	}
	// Pole wypełnia się przy każdej pozycji, także fałszem, bo puste
	// znaczyłoby brak wiedzy.
	for _, pozycja := range pozycje {
		jest := dolozone[pozycja.Name]
		pozycja.Attached = &jest
	}
	return nil
}
