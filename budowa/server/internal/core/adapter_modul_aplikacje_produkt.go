// Moduł Apps — Product Builder: produkt okna, etapy budowy, kamienie milowe,
// powiązania z innymi modułami i całościowa oś czasu projektu.
//
// Obsługiwane komendy: `apps.product.get`, `apps.product.save`,
// `apps.product.link.list`, `apps.stage.list`, `apps.stage.save`,
// `apps.milestone.list`, `apps.milestone.save`, `apps.milestone.delete`,
// `apps.timeline.list`. Typ adaptera i port stoją w
// `adapter_modul_aplikacje.go` i `adapter_modul_aplikacje_uchwyty.go`.
//
// OŚ CZASU NIE MA WŁASNEJ TABELI I MIEĆ JEJ NIE MOŻE. `AppTimelineKind` niesie
// dokładnie pięć wartości — architecture, workspace, stage, deployment,
// package — i każda z nich ma już w bazie swój wiersz ze znacznikiem czasu:
// historię wersji układu, plik warsztatu, etap, przebieg wdrożenia, pakiet.
// Osobna tabela zdarzeń byłaby SZÓSTĄ kopią tych samych faktów, rozjeżdżającą
// się przy pierwszym zapisie, który zapomni ją dopisać. Oś czasu składa się
// więc z odczytów tych pięciu źródeł i sortuje po czasie.
//
// POWIĄZANIA MODUŁÓW SĄ MIERZONE, NIE DEKLAROWANE. Opracowanie mówi wprost:
// „powiązania nie są aktywne domyślnie", a kontrakt nie daje komendy, którą
// dałoby się je włączyć — `apps.product.link.list` jest w rodzinie sam. Gdyby
// rdzeń oddawał tu stałą listę z `enabled: false`, meldowałby wykaz, za którym
// nic nie stoi. Zamiast tego każde powiązanie liczy w bazie to, czym naprawdę
// żyje w tym oknie, a pole `detail` mówi, co policzono. Powiązanie czynne
// znaczy więc „w tym oknie to powiązanie ma już treść", a nie „ktoś zaznaczył
// przełącznik".
package core

import (
	"context"
	"errors"
	"sort"
	"strconv"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// Przedrostki identyfikatorów bytów strony budowy produktu.
const (
	przedrostekProduktuApp = "prod-"
	przedrostekEtapuApp    = "etap-"
	przedrostekKamieniaApp = "kmil-"
)

// granicaOsiCzasuApp jest górną granicą strony osi czasu przy braku wskazania
// w żądaniu. Oś czasu składa się z pięciu źródeł naraz, więc projekt prowadzony
// od miesięcy oddawałby przy każdym otwarciu Product Buildera wszystko, co
// kiedykolwiek w nim zaszło.
const granicaOsiCzasuApp = 200

// PobierzProdukt obsługuje `apps.product.get`. Okno bez zapisanego produktu
// oddaje wynik z pustym polem, nie odmowę: świeże okno jeszcze niczego nie
// nazwało, a to jest stan normalny, nie brak.
func (a *adapterAplikacji) PobierzProdukt(ctx context.Context,
	z shared.AppsProductGetRequest) (shared.AppsProductGetResponse, error) {

	okno, err := oknoAplikacji(z.WindowId, "apps.product.get")
	if err != nil {
		return shared.AppsProductGetResponse{}, err
	}
	wiersz, err := a.repozytorium.ProduktApp(ctx, okno)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.AppsProductGetResponse{}, nil
	}
	if err != nil {
		return shared.AppsProductGetResponse{}, bladAplikacji(err)
	}
	produkt := produktKontraktu(wiersz)
	return shared.AppsProductGetResponse{Product: &produkt}, nil
}

// ZapiszProdukt obsługuje `apps.product.save`. Produkt jest jeden na okno, więc
// zapis nadpisuje zastany, a tożsamość (`AppProduct.id`) zostaje ta sama —
// panel, który zapamiętał identyfikator, nie traci go przy zmianie nazwy.
func (a *adapterAplikacji) ZapiszProdukt(ctx context.Context,
	z shared.AppsProductSaveRequest) (shared.AppsProductSaveResponse, error) {

	okno, err := oknoAplikacji(z.WindowId, "apps.product.save")
	if err != nil {
		return shared.AppsProductSaveResponse{}, err
	}
	if strings.TrimSpace(z.Name) == "" {
		return shared.AppsProductSaveResponse{}, bladWskazaniaAplikacji(
			"apps.product.save wymaga nazwy produktu")
	}
	platformy := make([]string, 0, len(z.Platforms))
	for _, platforma := range z.Platforms {
		if err := sprawdzPlatformeProduktuApp(platforma); err != nil {
			return shared.AppsProductSaveResponse{}, err
		}
		platformy = append(platformy, string(platforma))
	}

	zapisany, err := a.repozytorium.ZapiszProduktApp(ctx, dane.ProduktApp{
		Kod:          nowyIdentyfikator(przedrostekProduktuApp),
		Okno:         okno,
		Nazwa:        z.Name,
		Opis:         z.Description,
		Platformy:    platformy,
		Repozytorium: z.RepositoryUrl,
	})
	if err != nil {
		return shared.AppsProductSaveResponse{}, bladAplikacji(err)
	}
	return shared.AppsProductSaveResponse{Product: produktKontraktu(zapisany)}, nil
}

// WypiszPowiazaniaProduktu obsługuje `apps.product.link.list`. Każde powiązanie
// jest liczone w bazie — patrz rozstrzygnięcie na czole pliku.
func (a *adapterAplikacji) WypiszPowiazaniaProduktu(ctx context.Context,
	z shared.AppsProductLinkListRequest) (shared.AppsProductLinkListResponse, error) {

	okno, err := oknoAplikacji(z.WindowId, "apps.product.link.list")
	if err != nil {
		return shared.AppsProductLinkListResponse{}, err
	}

	pliki, err := a.repozytorium.PlikiWarsztatu(ctx, okno)
	if err != nil {
		return shared.AppsProductLinkListResponse{}, bladAplikacji(err)
	}
	etapy, err := a.repozytorium.EtapyApp(ctx, okno)
	if err != nil {
		return shared.AppsProductLinkListResponse{}, bladAplikacji(err)
	}
	wdrozenia, _, err := a.repozytorium.Wdrozenia(ctx, okno, nil, granicaWdrozenApp)
	if err != nil {
		return shared.AppsProductLinkListResponse{}, bladAplikacji(err)
	}
	_, dziennikRazem, err := a.repozytorium.DziennikUslugiApp(ctx, okno, "", 0, 1)
	if err != nil {
		return shared.AppsProductLinkListResponse{}, bladAplikacji(err)
	}

	komponentowFrontendu := 0
	if architektura, err := a.repozytorium.ArchitekturaOkna(ctx, okno); err == nil {
		komponenty, err := a.repozytorium.Komponenty(ctx, architektura.ID)
		if err != nil {
			return shared.AppsProductLinkListResponse{}, bladAplikacji(err)
		}
		for _, komponent := range komponenty {
			if komponent.Rodzaj == string(shared.AppComponentKindFrontend) {
				komponentowFrontendu++
			}
		}
	} else if !errors.Is(err, dane.ErrBrakWiersza) {
		return shared.AppsProductLinkListResponse{}, bladAplikacji(err)
	}

	wdrozenNieudanych := 0
	for _, wdrozenie := range wdrozenia {
		if wdrozenie.Stan == shared.AppDeployStatusFailed {
			wdrozenNieudanych++
		}
	}
	wykonawcy := map[string]struct{}{}
	etapowZWykonawca := 0
	for _, etap := range etapy {
		if etap.Wykonawca == nil || *etap.Wykonawca == "" {
			continue
		}
		etapowZWykonawca++
		wykonawcy[*etap.Wykonawca] = struct{}{}
	}

	powiazania := []shared.AppProductLink{
		powiazanieProduktu("developer", len(pliki), "plików warsztatu w oknie"),
		powiazanieProduktu("terminal", dziennikRazem, "wierszy dziennika przebiegów w oknie"),
		powiazanieProduktu("design", komponentowFrontendu, "komponentów warstwy interfejsu w architekturze"),
		powiazanieProduktu("agents", etapowZWykonawca, "etapów z przypisanym wykonawcą"),
		powiazanieProduktu("automations", len(wdrozenia), "przebiegów wdrożenia w dzienniku"),
		powiazanieProduktu("diagnostics", wdrozenNieudanych, "przebiegów wdrożenia zakończonych błędem"),
		powiazanieProduktu("multitaskingai", len(wykonawcy), "różnych wykonawców przypisanych do etapów"),
	}
	return shared.AppsProductLinkListResponse{Links: powiazania}, nil
}

// WypiszEtapy obsługuje `apps.stage.list`.
func (a *adapterAplikacji) WypiszEtapy(ctx context.Context,
	z shared.AppsStageListRequest) (shared.AppsStageListResponse, error) {

	okno, err := oknoAplikacji(z.WindowId, "apps.stage.list")
	if err != nil {
		return shared.AppsStageListResponse{}, err
	}
	if z.Status != nil {
		if err := sprawdzStanEtapuApp(*z.Status); err != nil {
			return shared.AppsStageListResponse{}, err
		}
	}
	wiersze, err := a.repozytorium.EtapyApp(ctx, okno)
	if err != nil {
		return shared.AppsStageListResponse{}, bladAplikacji(err)
	}
	etapy := make([]shared.AppStage, 0, len(wiersze))
	for _, wiersz := range wiersze {
		if z.Status != nil && wiersz.Stan != string(*z.Status) {
			continue
		}
		etapy = append(etapy, etapKontraktu(wiersz))
	}
	return shared.AppsStageListResponse{Stages: etapy, Total: len(etapy)}, nil
}

// ZapiszEtap obsługuje `apps.stage.save`. Brak `stageId` zakłada etap nowy;
// wskazanie zmienia zastany. Pusty łańcuch w `ownerAgentId` zdejmuje
// przypisanie wykonawcy — kontrakt mówi to wprost, więc pole rozróżnia trzy
// stany: brak wskazania (zostaw), pusty łańcuch (zdejmij), wartość (przypisz).
func (a *adapterAplikacji) ZapiszEtap(ctx context.Context,
	z shared.AppsStageSaveRequest) (shared.AppsStageSaveResponse, error) {

	okno, err := oknoAplikacji(z.WindowId, "apps.stage.save")
	if err != nil {
		return shared.AppsStageSaveResponse{}, err
	}
	kod := strings.TrimSpace(wartoscTekstu(z.StageId))
	nowy := kod == ""

	zastany := dane.EtapApp{Okno: okno, Stan: string(shared.AppStageStatusPending)}
	if nowy {
		kod = nowyIdentyfikator(przedrostekEtapuApp)
	} else {
		zastany, err = a.repozytorium.EtapApp(ctx, kod)
		if err != nil {
			return shared.AppsStageSaveResponse{}, bladNieznanegoBytuApp("etap budowy", kod, err)
		}
		if zastany.Okno != okno {
			return shared.AppsStageSaveResponse{}, bladWskazaniaAplikacji(
				"etap " + kod + " należy do okna " + zastany.Okno + ", a żądanie przyszło z okna " + okno)
		}
	}

	etap := dane.EtapApp{
		Kod: kod, Okno: okno, Nazwa: zastany.Nazwa,
		Kolejnosc: zastany.Kolejnosc, Stan: zastany.Stan, Wykonawca: zastany.Wykonawca,
	}
	if z.Name != nil {
		etap.Nazwa = *z.Name
	}
	if etap.Nazwa == "" {
		return shared.AppsStageSaveResponse{}, bladWskazaniaAplikacji(
			"apps.stage.save wymaga nazwy przy zakładaniu etapu")
	}
	if z.Order != nil {
		etap.Kolejnosc = *z.Order
	}
	if z.Status != nil {
		if err := sprawdzStanEtapuApp(*z.Status); err != nil {
			return shared.AppsStageSaveResponse{}, err
		}
		etap.Stan = string(*z.Status)
	}
	if z.OwnerAgentId != nil {
		if strings.TrimSpace(*z.OwnerAgentId) == "" {
			etap.Wykonawca = nil
		} else {
			etap.Wykonawca = z.OwnerAgentId
		}
	}

	zapisany, err := a.repozytorium.ZapiszEtapApp(ctx, etap)
	if err != nil {
		return shared.AppsStageSaveResponse{}, bladAplikacji(err)
	}
	// Etap jest drugim źródłem `apps.build.changed` obok wdrożenia — panel
	// Product Buildera rysuje oś etapów z tego samego zdarzenia, którym dostaje
	// przejścia wdrożenia.
	a.rozglosEtap(zmianaZalozenia(nowy), etapKontraktu(zapisany))
	return shared.AppsStageSaveResponse{Stage: etapKontraktu(zapisany)}, nil
}

// WypiszKamienieMilowe obsługuje `apps.milestone.list`.
func (a *adapterAplikacji) WypiszKamienieMilowe(ctx context.Context,
	z shared.AppsMilestoneListRequest) (shared.AppsMilestoneListResponse, error) {

	okno, err := oknoAplikacji(z.WindowId, "apps.milestone.list")
	if err != nil {
		return shared.AppsMilestoneListResponse{}, err
	}
	if z.Status != nil {
		if err := sprawdzStanKamieniaApp(*z.Status); err != nil {
			return shared.AppsMilestoneListResponse{}, err
		}
	}
	wiersze, err := a.repozytorium.KamienieMiloweApp(ctx, okno)
	if err != nil {
		return shared.AppsMilestoneListResponse{}, bladAplikacji(err)
	}
	kamienie := make([]shared.AppMilestone, 0, len(wiersze))
	for _, wiersz := range wiersze {
		if z.Status != nil && wiersz.Stan != string(*z.Status) {
			continue
		}
		kamienie = append(kamienie, kamienKontraktu(wiersz))
	}
	return shared.AppsMilestoneListResponse{Milestones: kamienie, Total: len(kamienie)}, nil
}

// ZapiszKamienMilowy obsługuje `apps.milestone.save`.
func (a *adapterAplikacji) ZapiszKamienMilowy(ctx context.Context,
	z shared.AppsMilestoneSaveRequest) (shared.AppsMilestoneSaveResponse, error) {

	okno, err := oknoAplikacji(z.WindowId, "apps.milestone.save")
	if err != nil {
		return shared.AppsMilestoneSaveResponse{}, err
	}
	if strings.TrimSpace(z.Name) == "" {
		return shared.AppsMilestoneSaveResponse{}, bladWskazaniaAplikacji(
			"apps.milestone.save wymaga nazwy kamienia milowego")
	}
	kod := strings.TrimSpace(wartoscTekstu(z.MilestoneId))
	nowy := kod == ""
	if nowy {
		kod = nowyIdentyfikator(przedrostekKamieniaApp)
	} else {
		zastany, err := a.repozytorium.KamienMilowyApp(ctx, kod)
		if err != nil {
			return shared.AppsMilestoneSaveResponse{}, bladNieznanegoBytuApp("kamień milowy", kod, err)
		}
		if zastany.Okno != okno {
			return shared.AppsMilestoneSaveResponse{}, bladWskazaniaAplikacji(
				"kamień milowy " + kod + " należy do okna " + zastany.Okno +
					", a żądanie przyszło z okna " + okno)
		}
	}

	stan := string(shared.AppMilestoneStatusPlanned)
	if z.Status != nil {
		if err := sprawdzStanKamieniaApp(*z.Status); err != nil {
			return shared.AppsMilestoneSaveResponse{}, err
		}
		stan = string(*z.Status)
	}

	zapisany, err := a.repozytorium.ZapiszKamienMilowyApp(ctx, dane.KamienMilowyApp{
		Kod: kod, Okno: okno, Nazwa: z.Name, Termin: z.DueAt, Stan: stan,
		KodyEtapow: z.StageIds,
	})
	if err != nil {
		return shared.AppsMilestoneSaveResponse{}, bladAplikacji(err)
	}
	return shared.AppsMilestoneSaveResponse{Milestone: kamienKontraktu(zapisany)}, nil
}

// UsunKamienMilowy obsługuje `apps.milestone.delete`. Kamień, którego nie ma,
// kończy się odmową `not_found`, a nie polem `deleted: false` — klient
// odróżnia „usunięto" od „nie było czego usunąć" po odmowie, bo pole logiczne
// oddane jako powodzenie kazałoby mu zgadywać, czy operacja się odbyła.
func (a *adapterAplikacji) UsunKamienMilowy(ctx context.Context,
	z shared.AppsMilestoneDeleteRequest) (shared.AppsMilestoneDeleteResponse, error) {

	okno, err := oknoAplikacji(z.WindowId, "apps.milestone.delete")
	if err != nil {
		return shared.AppsMilestoneDeleteResponse{}, err
	}
	kod := strings.TrimSpace(z.MilestoneId)
	if kod == "" {
		return shared.AppsMilestoneDeleteResponse{}, bladWskazaniaAplikacji(
			"apps.milestone.delete wymaga kamienia milowego")
	}
	zastany, err := a.repozytorium.KamienMilowyApp(ctx, kod)
	if err != nil {
		return shared.AppsMilestoneDeleteResponse{}, bladNieznanegoBytuApp("kamień milowy", kod, err)
	}
	if zastany.Okno != okno {
		return shared.AppsMilestoneDeleteResponse{}, bladWskazaniaAplikacji(
			"kamień milowy " + kod + " należy do okna " + zastany.Okno +
				", a żądanie przyszło z okna " + okno)
	}
	usuniety, err := a.repozytorium.UsunKamienMilowyApp(ctx, kod)
	if err != nil {
		return shared.AppsMilestoneDeleteResponse{}, bladAplikacji(err)
	}
	return shared.AppsMilestoneDeleteResponse{Deleted: usuniety}, nil
}

// WypiszOsCzasu obsługuje `apps.timeline.list` — chronologię projektu złożoną
// z pięciu źródeł (patrz czoło pliku), od najnowszego zdarzenia.
func (a *adapterAplikacji) WypiszOsCzasu(ctx context.Context,
	z shared.AppsTimelineListRequest) (shared.AppsTimelineListResponse, error) {

	okno, err := oknoAplikacji(z.WindowId, "apps.timeline.list")
	if err != nil {
		return shared.AppsTimelineListResponse{}, err
	}
	od := int64(0)
	if z.Since != nil {
		od = *z.Since
	}

	zdarzenia := []shared.AppTimelineEntry{}

	architektura, err := a.repozytorium.ArchitekturaOkna(ctx, okno)
	switch {
	case err == nil:
		wersje, err := a.repozytorium.WersjeArchitekturyApp(ctx, architektura.ID)
		if err != nil {
			return shared.AppsTimelineListResponse{}, bladAplikacji(err)
		}
		for _, wersja := range wersje {
			zdarzenia = append(zdarzenia, shared.AppTimelineEntry{
				Id:   architektura.Kod + "#" + strconv.Itoa(wersja.Wersja),
				Kind: shared.AppTimelineKindArchitecture,
				Summary: "architektura " + architektura.Kod + ": wersja " +
					strconv.Itoa(wersja.Wersja) + ", komponentów " +
					strconv.Itoa(wersja.LiczbaKomponentow),
				RefId:      wskaznikNapisuApp(architektura.Kod),
				OccurredAt: chwilaBazy(wersja.Utworzono),
			})
		}
	case errors.Is(err, dane.ErrBrakWiersza):
	default:
		return shared.AppsTimelineListResponse{}, bladAplikacji(err)
	}

	pliki, err := a.repozytorium.PlikiWarsztatu(ctx, okno)
	if err != nil {
		return shared.AppsTimelineListResponse{}, bladAplikacji(err)
	}
	for _, plik := range pliki {
		zdarzenia = append(zdarzenia, shared.AppTimelineEntry{
			Id:   "plik:" + string(plik.Warstwa) + ":" + plik.Sciezka,
			Kind: shared.AppTimelineKindWorkspace,
			Summary: "warsztat " + string(plik.Warstwa) + ": zapis pliku " + plik.Sciezka +
				" (" + strconv.FormatInt(plik.Rozmiar, 10) + " bajtów)",
			RefId:      wskaznikNapisuApp(plik.Sciezka),
			OccurredAt: chwilaBazy(plik.Zaktualizowano),
		})
	}

	etapy, err := a.repozytorium.EtapyApp(ctx, okno)
	if err != nil {
		return shared.AppsTimelineListResponse{}, bladAplikacji(err)
	}
	for _, etap := range etapy {
		zdarzenia = append(zdarzenia, shared.AppTimelineEntry{
			Id:         "etap:" + etap.Kod,
			Kind:       shared.AppTimelineKindStage,
			Summary:    "etap " + strconv.Quote(etap.Nazwa) + " w stanie " + etap.Stan,
			RefId:      wskaznikNapisuApp(etap.Kod),
			OccurredAt: chwilaBazy(etap.Zaktualizowano),
		})
	}

	wdrozenia, _, err := a.repozytorium.Wdrozenia(ctx, okno, nil, granicaWdrozenApp)
	if err != nil {
		return shared.AppsTimelineListResponse{}, bladAplikacji(err)
	}
	for _, wdrozenie := range wdrozenia {
		chwila := chwilaBazy(wdrozenie.Rozpoczeto)
		if wdrozenie.Zakonczono != nil {
			chwila = chwilaBazy(*wdrozenie.Zakonczono)
		}
		zdarzenia = append(zdarzenia, shared.AppTimelineEntry{
			Id:   "wdrozenie:" + wdrozenie.Kod,
			Kind: shared.AppTimelineKindDeployment,
			Summary: "wdrożenie na " + string(wdrozenie.Srodowisko) + " w stanie " +
				string(wdrozenie.Stan),
			RefId:      wskaznikNapisuApp(wdrozenie.Kod),
			OccurredAt: chwila,
		})
	}

	pakiety, err := a.repozytorium.PakietyApp(ctx, okno)
	if err != nil {
		return shared.AppsTimelineListResponse{}, bladAplikacji(err)
	}
	for _, pakiet := range pakiety {
		opis := "pakiet " + pakiet.Kod + " w formacie " + pakiet.Format
		if pakiet.RozszerzenieKod != nil {
			opis += ", opublikowany jako " + *pakiet.RozszerzenieKod
		}
		zdarzenia = append(zdarzenia, shared.AppTimelineEntry{
			Id:         "pakiet:" + pakiet.Kod,
			Kind:       shared.AppTimelineKindPackage,
			Summary:    opis,
			RefId:      wskaznikNapisuApp(pakiet.Kod),
			OccurredAt: chwilaBazy(pakiet.Zaktualizowano),
		})
	}

	wybrane := make([]shared.AppTimelineEntry, 0, len(zdarzenia))
	for _, zdarzenie := range zdarzenia {
		if zdarzenie.OccurredAt < od {
			continue
		}
		wybrane = append(wybrane, zdarzenie)
	}
	// Porządek malejący po czasie, a przy równym czasie po identyfikatorze:
	// pięć źródeł zapisuje znaczniki z rozdzielczością milisekundy, więc dwa
	// zdarzenia tej samej milisekundy bez drugiego klucza zamieniałyby się
	// miejscami między wywołaniami.
	sort.SliceStable(wybrane, func(i, j int) bool {
		if wybrane[i].OccurredAt != wybrane[j].OccurredAt {
			return wybrane[i].OccurredAt > wybrane[j].OccurredAt
		}
		return wybrane[i].Id < wybrane[j].Id
	})

	razem := len(wybrane)
	granica := granicaOsiCzasuApp
	if z.Limit != nil && *z.Limit > 0 {
		granica = *z.Limit
	}
	if len(wybrane) > granica {
		wybrane = wybrane[:granica]
	}
	return shared.AppsTimelineListResponse{Entries: wybrane, Total: razem}, nil
}

// powiazanieProduktu składa jedno powiązanie z policzonej wartości. Zdanie
// `detail` mówi, co policzono — bez niego „czynne" byłoby stwierdzeniem bez
// pokrycia (czoło pliku).
func powiazanieProduktu(kod string, ile int, co string) shared.AppProductLink {
	opis := strconv.Itoa(ile) + " " + co
	return shared.AppProductLink{ModuleCode: kod, Enabled: ile > 0, Detail: &opis}
}

// produktKontraktu przekłada wiersz produktu na kształt kontraktu.
func produktKontraktu(wiersz dane.ProduktApp) shared.AppProduct {
	platformy := make([]shared.AppProductPlatform, 0, len(wiersz.Platformy))
	for _, platforma := range wiersz.Platformy {
		platformy = append(platformy, shared.AppProductPlatform(platforma))
	}
	return shared.AppProduct{
		Id: wiersz.Kod, WindowId: wiersz.Okno, Name: wiersz.Nazwa,
		Description: wiersz.Opis, Platforms: platformy, RepositoryUrl: wiersz.Repozytorium,
		UpdatedAt: chwilaBazy(wiersz.Zaktualizowano),
	}
}

// etapKontraktu przekłada wiersz etapu na kształt kontraktu.
func etapKontraktu(wiersz dane.EtapApp) shared.AppStage {
	kolejnosc := wiersz.Kolejnosc
	return shared.AppStage{
		Id: wiersz.Kod, WindowId: wiersz.Okno, Name: wiersz.Nazwa,
		Order: &kolejnosc, Status: shared.AppStageStatus(wiersz.Stan),
		OwnerAgentId: wiersz.Wykonawca, UpdatedAt: chwilaBazy(wiersz.Zaktualizowano),
	}
}

// kamienKontraktu przekłada wiersz kamienia milowego na kształt kontraktu.
func kamienKontraktu(wiersz dane.KamienMilowyApp) shared.AppMilestone {
	return shared.AppMilestone{
		Id: wiersz.Kod, WindowId: wiersz.Okno, Name: wiersz.Nazwa,
		DueAt: wiersz.Termin, Status: shared.AppMilestoneStatus(wiersz.Stan),
		StageIds: wiersz.KodyEtapow, UpdatedAt: chwilaBazy(wiersz.Zaktualizowano),
	}
}

// sprawdzStanEtapuApp dopuszcza wyłącznie stany kontraktu.
func sprawdzStanEtapuApp(stan shared.AppStageStatus) error {
	switch stan {
	case shared.AppStageStatusPending, shared.AppStageStatusActive,
		shared.AppStageStatusDone, shared.AppStageStatusBlocked:
		return nil
	}
	return bladWskazaniaAplikacji("nieznany stan etapu " + strconv.Quote(string(stan)) +
		" — dopuszczalne: pending, active, done, blocked")
}

// sprawdzStanKamieniaApp dopuszcza wyłącznie stany kontraktu.
func sprawdzStanKamieniaApp(stan shared.AppMilestoneStatus) error {
	switch stan {
	case shared.AppMilestoneStatusPlanned, shared.AppMilestoneStatusActive,
		shared.AppMilestoneStatusReached, shared.AppMilestoneStatusMissed:
		return nil
	}
	return bladWskazaniaAplikacji("nieznany stan kamienia milowego " + strconv.Quote(string(stan)) +
		" — dopuszczalne: planned, active, reached, missed")
}

// sprawdzPlatformeProduktuApp dopuszcza wyłącznie platformy kontraktu.
func sprawdzPlatformeProduktuApp(platforma shared.AppProductPlatform) error {
	switch platforma {
	case shared.AppProductPlatformWeb, shared.AppProductPlatformMobile,
		shared.AppProductPlatformDesktop:
		return nil
	}
	return bladWskazaniaAplikacji("nieznana platforma produktu " + strconv.Quote(string(platforma)) +
		" — dopuszczalne: web, mobile, desktop")
}
