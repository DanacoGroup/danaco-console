// Odpowiedzialność pliku: rząd kart okna przeglądarki i jego przestrzenie
// robocze — dziewięć komend `browser.tab.*` i `browser.workspace.*`. Karta
// otwarta z adresem naprawdę pod ten adres przechodzi. Przestrzeń zapisuje
// skład kart, nie ich kopię.
package core

import (
	"context"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// OtworzKarte obsługuje `browser.tab.open`: otwiera nową kartę w oknie
// i wpisuje jej wiersz do rzędu kart przeglądarki.
func (a *adapterPrzegladarki) OtworzKarte(ctx context.Context,
	z shared.BrowserTabOpenRequest) (shared.BrowserTabOpenResponse, error) {

	if strings.TrimSpace(z.WindowId) == "" {
		return shared.BrowserTabOpenResponse{}, bladWskazaniaPrzegladarki("komenda tab.open bez okna")
	}
	kolejnosc, err := a.repozytorium.NastepnaKolejnoscKarty(ctx, z.WindowId)
	if err != nil {
		return shared.BrowserTabOpenResponse{}, bladPrzegladarki(err)
	}

	// Karta w tle nie odbiera czynności bieżącej — rozstrzyga pole
	// `background`, nie kolejność wywołań.
	wTle := z.Background != nil && *z.Background
	stan := "active"
	if wTle {
		stan = "inactive"
	}

	adres := wartoscTekstuLubPusta(z.Url)
	var migawka *shared.BrowserSnapshot
	tytul := (*string)(nil)
	if adres != "" {
		odpowiedz, err := a.Nawiguj(ctx, shared.BrowserNavigateRequest{WindowId: z.WindowId, Url: adres})
		if err != nil {
			return shared.BrowserTabOpenResponse{}, err
		}
		kopia := odpowiedz.Snapshot
		migawka = &kopia
		tytul = kopia.Title
		adres = kopia.Url
	}

	karta := dane.KartaPrzegladania{
		Kod:              nowyIdentyfikator(przedrostekKartyPrzegladania),
		Okno:             z.WindowId,
		Stan:             stan,
		Grupa:            z.GroupId,
		KartaOtwierajaca: z.OpenerTabId,
		Kolejnosc:        kolejnosc,
		OstatnioCzynna:   terazWBazie(),
	}
	if adres != "" {
		karta.Url = &adres
	}
	if tytul != nil {
		karta.Tytul = tytul
	}
	zapisana, err := a.repozytorium.ZapiszKarte(ctx, karta)
	if err != nil {
		return shared.BrowserTabOpenResponse{}, bladPrzegladarki(err)
	}
	if !wTle {
		if err := a.repozytorium.OdznaczPozostaleKarty(ctx, z.WindowId, zapisana.Kod); err != nil {
			return shared.BrowserTabOpenResponse{}, bladPrzegladarki(err)
		}
	}
	return shared.BrowserTabOpenResponse{Tab: kartaPrzegladaniaKontraktu(zapisana), Snapshot: migawka}, nil
}

// WykazKart obsługuje `browser.tab.list` — oddaje rząd kart okna wraz
// z grupami, do których karty należą.
func (a *adapterPrzegladarki) WykazKart(ctx context.Context,
	z shared.BrowserTabListRequest) (shared.BrowserTabListResponse, error) {

	if strings.TrimSpace(z.WindowId) == "" {
		return shared.BrowserTabListResponse{}, bladWskazaniaPrzegladarki("komenda tab.list bez okna")
	}
	wiersze, err := a.repozytorium.Karty(ctx, dane.FiltrKartPrzegladania{
		Okno:          z.WindowId,
		Przestrzen:    wartoscTekstuLubPusta(z.WorkspaceId),
		ZZawieszonymi: z.IncludeSuspended == nil || *z.IncludeSuspended,
	})
	if err != nil {
		return shared.BrowserTabListResponse{}, bladPrzegladarki(err)
	}
	karty := make([]shared.BrowserTab, 0, len(wiersze))
	for _, wiersz := range wiersze {
		karty = append(karty, kartaPrzegladaniaKontraktu(wiersz))
	}

	grupy, err := a.repozytorium.GrupyKart(ctx, z.WindowId, 0)
	if err != nil {
		return shared.BrowserTabListResponse{}, bladPrzegladarki(err)
	}
	wykazGrup := make([]shared.BrowserTabGroup, 0, len(grupy))
	for _, grupa := range grupy {
		wykazGrup = append(wykazGrup, grupaKartKontraktu(grupa, kodyKartGrupy(wiersze, grupa.Kod)))
	}
	return shared.BrowserTabListResponse{Tabs: karty, Groups: wykazGrup}, nil
}

// ZmienKarte obsługuje `browser.tab.update` — czynność, zawieszenie, przypięcie,
// grupę i kolejność karty. Karta wznowiona z zawieszenia przechodzi pod
// zapamiętany adres i oddaje migawkę, zamiast pokazać pustkę.
func (a *adapterPrzegladarki) ZmienKarte(ctx context.Context,
	z shared.BrowserTabUpdateRequest) (shared.BrowserTabUpdateResponse, error) {

	if strings.TrimSpace(z.TabId) == "" {
		return shared.BrowserTabUpdateResponse{}, bladWskazaniaPrzegladarki("komenda tab.update bez karty")
	}
	karta, err := a.repozytorium.Karta(ctx, z.TabId)
	if err != nil {
		return shared.BrowserTabUpdateResponse{}, bladWierszaPrzegladania("karty przeglądania", z.TabId, err)
	}

	zawieszona := z.Suspended != nil && *z.Suspended
	czynna := z.Active != nil && *z.Active
	switch {
	case zawieszona:
		karta.Stan = "suspended"
	case czynna:
		karta.Stan = "active"
	}
	if z.Pinned != nil {
		karta.Przypieta = *z.Pinned
	}
	if z.GroupId != nil {
		karta.Grupa = z.GroupId
	}
	if z.Order != nil {
		karta.Kolejnosc = int64(*z.Order)
	}

	var migawka *shared.BrowserSnapshot
	if czynna && !zawieszona && karta.Url != nil && *karta.Url != "" {
		odpowiedz, err := a.Nawiguj(ctx, shared.BrowserNavigateRequest{WindowId: karta.Okno, Url: *karta.Url})
		if err != nil {
			return shared.BrowserTabUpdateResponse{}, err
		}
		kopia := odpowiedz.Snapshot
		migawka = &kopia
		karta.Tytul = kopia.Title
		karta.OstatnioCzynna = terazWBazie()
	}

	zapisana, err := a.repozytorium.ZapiszKarte(ctx, karta)
	if err != nil {
		return shared.BrowserTabUpdateResponse{}, bladPrzegladarki(err)
	}
	if karta.Stan == "active" {
		if err := a.repozytorium.OdznaczPozostaleKarty(ctx, karta.Okno, karta.Kod); err != nil {
			return shared.BrowserTabUpdateResponse{}, bladPrzegladarki(err)
		}
	}
	return shared.BrowserTabUpdateResponse{Tab: kartaPrzegladaniaKontraktu(zapisana), Snapshot: migawka}, nil
}

// ZamknijKarte obsługuje `browser.tab.close`: usuwa wiersz karty z rzędu kart
// okna przeglądarki na trwałe.
func (a *adapterPrzegladarki) ZamknijKarte(ctx context.Context,
	z shared.BrowserTabCloseRequest) (shared.BrowserTabCloseResponse, error) {

	if strings.TrimSpace(z.TabId) == "" {
		return shared.BrowserTabCloseResponse{}, bladWskazaniaPrzegladarki("komenda tab.close bez karty")
	}
	if _, err := a.repozytorium.Karta(ctx, z.TabId); err != nil {
		return shared.BrowserTabCloseResponse{}, bladWierszaPrzegladania("karty przeglądania", z.TabId, err)
	}
	zamknieta, err := a.repozytorium.ZamknijKarte(ctx, z.TabId)
	if err != nil {
		return shared.BrowserTabCloseResponse{}, bladPrzegladarki(err)
	}
	return shared.BrowserTabCloseResponse{Closed: zamknieta}, nil
}

// UstawGrupeKart obsługuje `browser.tab.group.set` — założenie grupy, zmianę
// jej nazwy, barwy i składu albo jej zdjęcie.
func (a *adapterPrzegladarki) UstawGrupeKart(ctx context.Context,
	z shared.BrowserTabGroupSetRequest) (shared.BrowserTabGroupSetResponse, error) {

	if strings.TrimSpace(z.WindowId) == "" || strings.TrimSpace(z.Name) == "" {
		return shared.BrowserTabGroupSetResponse{}, bladWskazaniaPrzegladarki(
			"komenda tab.group.set bez okna albo nazwy grupy")
	}
	kod := wartoscTekstuLubPusta(z.GroupId)

	if z.Removed != nil && *z.Removed {
		if kod == "" {
			return shared.BrowserTabGroupSetResponse{}, bladWskazaniaPrzegladarki(
				"zdjęcie grupy kart bez jej wskazania")
		}
		grupa, err := a.repozytorium.GrupaKart(ctx, kod)
		if err != nil {
			return shared.BrowserTabGroupSetResponse{}, bladWierszaPrzegladania("grupy kart", kod, err)
		}
		if _, err := a.repozytorium.UsunGrupeKart(ctx, kod); err != nil {
			return shared.BrowserTabGroupSetResponse{}, bladPrzegladarki(err)
		}
		// Grupa zdjęta wychodzi kontraktem taka jak przy zdjęciu, ze składem
		// pustym: karty zostały otwarte.
		return shared.BrowserTabGroupSetResponse{Group: grupaKartKontraktu(grupa, nil)}, nil
	}

	if kod == "" {
		kod = nowyIdentyfikator(przedrostekGrupyKart)
	}
	zwinieta := z.Collapsed != nil && *z.Collapsed
	zapisana, err := a.repozytorium.ZapiszGrupeKart(ctx, dane.GrupaKart{
		Kod: kod, Okno: z.WindowId, Nazwa: z.Name, Barwa: z.Color, Zwinieta: zwinieta,
	})
	if err != nil {
		return shared.BrowserTabGroupSetResponse{}, bladPrzegladarki(err)
	}
	for _, karta := range z.TabIds {
		if err := a.repozytorium.PrzypiszKarteDoGrupy(ctx, z.WindowId, karta, kod); err != nil {
			return shared.BrowserTabGroupSetResponse{}, bladPrzegladarki(err)
		}
	}
	wiersze, err := a.repozytorium.Karty(ctx, dane.FiltrKartPrzegladania{Okno: z.WindowId, ZZawieszonymi: true})
	if err != nil {
		return shared.BrowserTabGroupSetResponse{}, bladPrzegladarki(err)
	}
	return shared.BrowserTabGroupSetResponse{
		Group: grupaKartKontraktu(zapisana, kodyKartGrupy(wiersze, kod)),
	}, nil
}

// ZapiszPrzestrzen obsługuje `browser.workspace.save` — przypisuje wskazane
// karty do nazwanej przestrzeni roboczej.
func (a *adapterPrzegladarki) ZapiszPrzestrzen(ctx context.Context,
	z shared.BrowserWorkspaceSaveRequest) (shared.BrowserWorkspaceSaveResponse, error) {

	if strings.TrimSpace(z.WindowId) == "" || strings.TrimSpace(z.Name) == "" {
		return shared.BrowserWorkspaceSaveResponse{}, bladWskazaniaPrzegladarki(
			"komenda workspace.save bez okna albo nazwy przestrzeni")
	}
	kod := wartoscTekstuLubPusta(z.WorkspaceId)
	if kod == "" {
		kod = nowyIdentyfikator(przedrostekPrzestrzeni)
	}
	okno := z.WindowId
	zapisana, err := a.repozytorium.ZapiszPrzestrzen(ctx, dane.PrzestrzenPrzegladania{
		Kod: kod, Okno: &okno, Nazwa: z.Name,
	})
	if err != nil {
		return shared.BrowserWorkspaceSaveResponse{}, bladPrzegladarki(err)
	}

	// Żądanie bez wykazu kart zapisuje rząd kart taki, jaki jest, a nie
	// przestrzeń pustą.
	karty := z.TabIds
	if len(karty) == 0 {
		wiersze, err := a.repozytorium.Karty(ctx, dane.FiltrKartPrzegladania{Okno: okno, ZZawieszonymi: true})
		if err != nil {
			return shared.BrowserWorkspaceSaveResponse{}, bladPrzegladarki(err)
		}
		for _, wiersz := range wiersze {
			karty = append(karty, wiersz.Kod)
		}
	}
	for _, karta := range karty {
		if err := a.repozytorium.PrzypiszKarteDoPrzestrzeni(ctx, okno, karta, kod); err != nil {
			return shared.BrowserWorkspaceSaveResponse{}, bladPrzegladarki(err)
		}
	}
	przestrzen, err := a.przestrzenKontraktu(ctx, zapisana)
	if err != nil {
		return shared.BrowserWorkspaceSaveResponse{}, err
	}
	return shared.BrowserWorkspaceSaveResponse{Workspace: przestrzen}, nil
}

// WykazPrzestrzeni obsługuje `browser.workspace.list` — oddaje wykaz
// przestrzeni roboczych zapisanych w oknie.
func (a *adapterPrzegladarki) WykazPrzestrzeni(ctx context.Context,
	z shared.BrowserWorkspaceListRequest) (shared.BrowserWorkspaceListResponse, error) {

	wiersze, err := a.repozytorium.Przestrzenie(ctx, wartoscTekstuLubPusta(z.WindowId), wartoscLiczby(z.Limit))
	if err != nil {
		return shared.BrowserWorkspaceListResponse{}, bladPrzegladarki(err)
	}
	przestrzenie := make([]shared.BrowserWorkspace, 0, len(wiersze))
	for _, wiersz := range wiersze {
		przestrzen, err := a.przestrzenKontraktu(ctx, wiersz)
		if err != nil {
			return shared.BrowserWorkspaceListResponse{}, err
		}
		przestrzenie = append(przestrzenie, przestrzen)
	}
	return shared.BrowserWorkspaceListResponse{Workspaces: przestrzenie}, nil
}

// OtworzPrzestrzen obsługuje `browser.workspace.open` — przywraca zapisany
// zestaw kart w oknie, budząc je uśpione zamiast otwierać od razu.
func (a *adapterPrzegladarki) OtworzPrzestrzen(ctx context.Context,
	z shared.BrowserWorkspaceOpenRequest) (shared.BrowserWorkspaceOpenResponse, error) {

	if strings.TrimSpace(z.WorkspaceId) == "" || strings.TrimSpace(z.WindowId) == "" {
		return shared.BrowserWorkspaceOpenResponse{}, bladWskazaniaPrzegladarki(
			"komenda workspace.open bez przestrzeni albo okna")
	}
	wiersz, err := a.repozytorium.Przestrzen(ctx, z.WorkspaceId)
	if err != nil {
		return shared.BrowserWorkspaceOpenResponse{}, bladWierszaPrzegladania("przestrzeni przeglądania", z.WorkspaceId, err)
	}

	if z.CloseCurrent != nil && *z.CloseCurrent {
		biezace, err := a.repozytorium.Karty(ctx, dane.FiltrKartPrzegladania{Okno: z.WindowId, ZZawieszonymi: true})
		if err != nil {
			return shared.BrowserWorkspaceOpenResponse{}, bladPrzegladarki(err)
		}
		for _, karta := range biezace {
			if karta.Przestrzen != nil && *karta.Przestrzen == z.WorkspaceId {
				continue
			}
			if _, err := a.repozytorium.ZamknijKarte(ctx, karta.Kod); err != nil {
				return shared.BrowserWorkspaceOpenResponse{}, bladPrzegladarki(err)
			}
		}
	}

	// Karty przestrzeni wracają otwarte, ale uśpione — budzi je
	// `browser.tab.update` z polem `active`.
	zapisane, err := a.repozytorium.Karty(ctx, dane.FiltrKartPrzegladania{
		Okno:       wartoscTekstuLubPusta(wiersz.Okno),
		Przestrzen: z.WorkspaceId, ZZawieszonymi: true, ZZamknietymi: true,
	})
	if err != nil {
		return shared.BrowserWorkspaceOpenResponse{}, bladPrzegladarki(err)
	}
	karty := make([]shared.BrowserTab, 0, len(zapisane))
	for _, karta := range zapisane {
		karta.Okno = z.WindowId
		karta.Zamknieta = false
		if karta.Stan == "active" {
			karta.Stan = "suspended"
		}
		odtworzona, err := a.repozytorium.ZapiszKarte(ctx, karta)
		if err != nil {
			return shared.BrowserWorkspaceOpenResponse{}, bladPrzegladarki(err)
		}
		karty = append(karty, kartaPrzegladaniaKontraktu(odtworzona))
	}

	przestrzen, err := a.przestrzenKontraktu(ctx, wiersz)
	if err != nil {
		return shared.BrowserWorkspaceOpenResponse{}, err
	}
	return shared.BrowserWorkspaceOpenResponse{Workspace: przestrzen, Tabs: karty}, nil
}

// UsunPrzestrzen obsługuje `browser.workspace.remove`. Karty zostają otwarte —
// kontrakt mówi o tym wprost.
func (a *adapterPrzegladarki) UsunPrzestrzen(ctx context.Context,
	z shared.BrowserWorkspaceRemoveRequest) (shared.BrowserWorkspaceRemoveResponse, error) {

	if strings.TrimSpace(z.WorkspaceId) == "" {
		return shared.BrowserWorkspaceRemoveResponse{}, bladWskazaniaPrzegladarki(
			"komenda workspace.remove bez przestrzeni")
	}
	usunieta, err := a.repozytorium.UsunPrzestrzen(ctx, z.WorkspaceId)
	if err != nil {
		return shared.BrowserWorkspaceRemoveResponse{}, bladPrzegladarki(err)
	}
	if !usunieta {
		return shared.BrowserWorkspaceRemoveResponse{}, bladNieznanegoBytu("przestrzeni przeglądania", z.WorkspaceId)
	}
	return shared.BrowserWorkspaceRemoveResponse{Removed: true}, nil
}

// przestrzenKontraktu składa `BrowserWorkspace` wraz z liczbą kart liczoną
// z bazy, nie z pamięci wołającego.
func (a *adapterPrzegladarki) przestrzenKontraktu(ctx context.Context,
	wiersz dane.PrzestrzenPrzegladania) (shared.BrowserWorkspace, error) {

	liczba, err := a.repozytorium.LiczbaKartPrzestrzeni(ctx, wiersz.Kod)
	if err != nil {
		return shared.BrowserWorkspace{}, bladPrzegladarki(err)
	}
	return shared.BrowserWorkspace{
		Id:        wiersz.Kod,
		Name:      wiersz.Nazwa,
		WindowId:  wiersz.Okno,
		TabCount:  int(liczba),
		ProfileId: wiersz.Profil,
		CreatedAt: chwilaBazy(wiersz.Utworzono),
		UpdatedAt: chwilaBazy(wiersz.Zaktualizowano),
	}, nil
}

// kartaPrzegladaniaKontraktu przekłada wiersz karty z bazy na byt kontraktu,
// który komendy oddają wołającemu.
func kartaPrzegladaniaKontraktu(w dane.KartaPrzegladania) shared.BrowserTab {
	karta := shared.BrowserTab{
		Id:           w.Kod,
		WindowId:     w.Okno,
		Url:          w.Url,
		Title:        w.Tytul,
		State:        shared.BrowserTabState(w.Stan),
		GroupId:      w.Grupa,
		WorkspaceId:  w.Przestrzen,
		OpenerTabId:  w.KartaOtwierajaca,
		Order:        int(w.Kolejnosc),
		LastActiveAt: chwilaZeZnacznika(w.OstatnioCzynna),
		CreatedAt:    chwilaBazy(w.Utworzono),
	}
	if w.Przypieta {
		przypieta := true
		karta.Pinned = &przypieta
	}
	return karta
}

// grupaKartKontraktu przekłada wiersz grupy kart wraz z jej składem na byt
// kontraktu oddawany wołającemu.
func grupaKartKontraktu(w dane.GrupaKart, karty []string) shared.BrowserTabGroup {
	grupa := shared.BrowserTabGroup{
		Id:        w.Kod,
		WindowId:  w.Okno,
		Name:      w.Nazwa,
		Color:     w.Barwa,
		Collapsed: w.Zwinieta,
		CreatedAt: chwilaBazy(w.Utworzono),
	}
	if len(karty) > 0 {
		grupa.TabIds = karty
	}
	return grupa
}

// kodyKartGrupy wybiera z rzędu kart te należące do wskazanej grupy. Skład
// liczony z kart, nie z osobnej listy — jedna prawda o przynależności.
func kodyKartGrupy(karty []dane.KartaPrzegladania, grupa string) []string {
	if grupa == "" {
		return nil
	}
	var kody []string
	for _, karta := range karty {
		if karta.Grupa != nil && *karta.Grupa == grupa {
			kody = append(kody, karta.Kod)
		}
	}
	return kody
}

// bladWierszaPrzegladania odróżnia odpowiedź „nie ma takiego bytu” od usterki
// samego odczytu wiersza z bazy.
func bladWierszaPrzegladania(co, wskazanie string, err error) error {
	if err == nil {
		return nil
	}
	if isBrakWiersza(err) {
		return bladNieznanegoBytu(co, wskazanie)
	}
	return bladPrzegladarki(err)
}
