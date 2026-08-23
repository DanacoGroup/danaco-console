// Moduł Apps: typ adaptera, konstruktor i komenda `apps.architecture.define`.
// Warsztat frontendu i backendu (`apps.workspace.update`) oraz przebiegi
// wdrożenia (`apps.deployment.run`) stoją w `adapter_modul_aplikacje_wdrozenie.go`
// i `adapter_modul_aplikacje_uchwyty.go`; tam też deklarowany jest port
// `Aplikacje` i `zarejestrujAplikacje`.
//
// Zastrzeżenia walidacji są przechowywane, nie wyliczane. Kontrakt niesie
// `AppArchitecture.ValidationIssues` jako wynik, ale `AppsArchitectureDefineRequest`
// nie niesie żadnej reguły ani listy naruszeń do sprawdzenia. Adapter zapisuje
// więc to pole tak, jak zastał je w bazie przy poprzednim zapisie
// (`dane/aplikacje.go` utrzymuje je w UPSERT-cie), i nie liczy niczego nowego.
package core

import (
	"context"
	"errors"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// Przedrostki identyfikatorów bytów modułu. Zadeklarowane tu w całości, łącznie
// z przedrostkiem wdrożenia, którego ten plik nie używa — jedno miejsce nadawania
// przedrostków modułu zamiast kilku.
const (
	przedrostekArchitekturyApp = "arch-"
	przedrostekKomponentuApp   = "kmp-"
	przedrostekWdrozeniaApp    = "wdr-"
)

// adapterAplikacji wypełnia część portu Aplikacje. Zależność jest jedna:
// wspólne repozytorium modułu, rozłożone po stronie danych na trzy pliki wedle
// odpowiedzialności, ale niosące jeden typ.
type adapterAplikacji struct {
	repozytorium dane.RepozytoriumAplikacji
	// przyrostWdrozenia rozgłasza `apps.build.changed`. Podpina go obsługiwacz
	// (`adapter_modul_aplikacje_uchwyty.go`). Silnik wykonania wdrożenia
	// przesuwa stan przebiegu poza wykonaniem komendy
	// (`adapter_modul_aplikacje_wdrozenie_bieg.go`), więc rozgłoszenie nie może
	// iść wyłącznie z obsługiwacza żądania — tak jak `przyrost` w module
	// Developer.
	przyrostWdrozenia func(shared.ChangeKind, shared.AppDeployment)
	// przyrostWarsztatu rozgłasza `apps.workspace.changed`. Bez niego drugie okno
	// tej samej przestrzeni nie dowiaduje się o zmianie pliku warsztatu.
	przyrostWarsztatu func(shared.ChangeKind, string, shared.AppWorkspaceLayer, shared.DeveloperFile)
	// okna są rejestrem okien sesji. Wdrożenie dzieje się w oknie — z jego
	// przestrzeni roboczej bierze się to, co ma pojechać — więc okno musi
	// istnieć, zanim powstanie wiersz przebiegu. Zależność opcjonalna: bez niej
	// moduł pracuje, ale montaż ją wpina.
	okna *session.Rejestr
	// przyrostEtapu rozgłasza `apps.build.changed` przy zmianie etapu budowy.
	// Osobna droga od wdrożeniowej, bo zdarzenie niesie inny byt: etap Product
	// Buildera, nie przebieg wdrożenia.
	przyrostEtapu func(shared.ChangeKind, shared.AppStage)
	// magazyn trzyma bajty wytworów modułu — eksport diagramu, artefakt
	// budowania, archiwum pakietu. To ten sam magazyn treści, którym jadą
	// zasoby modułu Design i pliki biblioteki: wiersz w bazie wskazuje plik na
	// dysku, a nie udaje, że go ma.
	magazyn *magazynTresciBiblioteki
	// katalogDanych jest korzeniem, względem którego liczone są odwołania
	// magazynu — bez niego odwołanie wypuszczone z rdzenia wynosiłoby układ
	// katalogów maszyny.
	katalogDanych string
	// sejf wydaje klucz wydawcy po jego kluczu jawnym. `apps.package.sign`
	// dostaje w żądaniu WYŁĄCZNIE odwołanie (`signingKeyRef`), nigdy treść —
	// materiał klucza nie przechodzi przez kontrakt ani przez bazę modułu.
	sejf sejfKluczaWydawcy
	// katalogRozszerzen jest rejestrem pozycji katalogu. `apps.package.publish`
	// zakłada w nim pozycję z manifestu pakietu — prywatny rejestr organizacji
	// nie jest drugim rejestrem obok `extension.*`, tylko tym samym.
	katalogRozszerzen dane.RepozytoriumRozszerzen
	// podglady trzymają stojące serwery podglądu warstw. Rejestr żyje wyłącznie
	// w pamięci: serwer nie przeżywa restartu rdzenia, więc wiersz w bazie
	// mówiłby po restarcie o nasłuchu, którego nie ma.
	podglady *rejestrPodgladowApp
}

// sejfKluczaWydawcy jest wycinkiem sejfu poświadczeń, którego moduł potrzebuje.
// Wycinek, nie cały sejf: `apps.package.sign` czyta i — przy pierwszym użyciu
// odwołania — zakłada klucz wydawcy, a kasowanie poświadczeń do modułu nie
// należy.
type sejfKluczaWydawcy interface {
	Odczytaj(ctx context.Context, byt string) (string, bool)
	Zapisz(ctx context.Context, byt, poswiadczenie string) (string, error)
}

// ZMagazynem wpina katalog danych rdzenia: magazyn bajtów wytworów modułu oraz
// sejf, z którego bierze się klucz wydawcy. Bez niego eksport, pakowanie
// i podpis odmawiają z powodem, zamiast meldować wytwór bez bajtów.
func (a *adapterAplikacji) ZMagazynem(katalogDanych string) *adapterAplikacji {
	a.katalogDanych = katalogDanych
	a.magazyn = magazynWytworowApp(katalogDanych)
	a.sejf = dane.NowySejfPlikowy(katalogDanych)
	return a
}

// ZKatalogiemRozszerzen wpina rejestr pozycji katalogu — drogę
// `apps.package.publish` do prywatnego rejestru organizacji.
func (a *adapterAplikacji) ZKatalogiemRozszerzen(rejestr dane.RepozytoriumRozszerzen) *adapterAplikacji {
	a.katalogRozszerzen = rejestr
	return a
}

// PodepnijPrzyrostEtapu oddaje adapterowi drogę do `apps.build.changed` przy
// zmianie etapu budowy.
func (a *adapterAplikacji) PodepnijPrzyrostEtapu(rozglos func(shared.ChangeKind, shared.AppStage)) {
	a.przyrostEtapu = rozglos
}

// ZOknami wpina rejestr okien sesji. Bez niego `apps.deployment.run` zakłada
// przebieg dla okna, którego nie ma, i kończy go powodem o pustym warsztacie
// zamiast o braku okna.
func (a *adapterAplikacji) ZOknami(rejestr *session.Rejestr) *adapterAplikacji {
	a.okna = rejestr
	return a
}

// PodepnijPrzyrostWarsztatu oddaje adapterowi drogę do `apps.workspace.changed`.
// Osobna droga od wdrożeniowej, bo zdarzenie niesie inny byt (plik warstwy, nie
// przebieg) — jedna funkcja o dwóch znaczeniach byłaby dwiema prawdami.
func (a *adapterAplikacji) PodepnijPrzyrostWarsztatu(
	rozglos func(shared.ChangeKind, string, shared.AppWorkspaceLayer, shared.DeveloperFile)) {

	a.przyrostWarsztatu = rozglos
}

// rozglosWarsztat oddaje zmianę pliku obsługiwaczowi. Brak podpięcia nie
// zmienia pracy modułu — rdzeń zapisuje także wtedy, gdy nikt nie słucha
// zdarzeń.
func (a *adapterAplikacji) rozglosWarsztat(zmiana shared.ChangeKind, okno string,
	warstwa shared.AppWorkspaceLayer, plik shared.DeveloperFile) {

	if a.przyrostWarsztatu == nil {
		return
	}
	a.przyrostWarsztatu(zmiana, okno, warstwa, plik)
}

// PodepnijPrzyrostWdrozenia oddaje adapterowi drogę do zdarzenia zmiany
// wdrożenia. Wdrożenie kończy się poza wykonaniem komendy: przebieg zakładany
// jest w stanie `pending`, a silnik przesuwa go przez `running` do
// `succeeded`/`failed` już po odesłaniu odpowiedzi, więc rozgłoszenie tych
// przejść nie może wychodzić z obsługiwacza żądania (wzorzec
// `PodepnijPrzyrostBudowania` modułu Developer).
func (a *adapterAplikacji) PodepnijPrzyrostWdrozenia(
	rozglos func(shared.ChangeKind, shared.AppDeployment)) {

	a.przyrostWdrozenia = rozglos
}

// nowyAdapterAplikacji wiąże adapter z repozytorium modułu. Rejestr stojących
// podglądów powstaje od razu — `Zamknij` musi mieć co zamknąć także wtedy, gdy
// żaden podgląd nie ruszył.
func nowyAdapterAplikacji(repozytorium dane.RepozytoriumAplikacji) *adapterAplikacji {
	return &adapterAplikacji{repozytorium: repozytorium, podglady: nowyRejestrPodgladowApp()}
}

// Zamknij zatrzymuje serwery podglądu podniesione przez `apps.preview.start`.
// Nasłuch żyje poza żądaniem, więc bez tego przeżyłby zatrzymanie rdzenia
// i zostawił zajęty port — tak samo jak przebieg budowania modułu Developer.
func (a *adapterAplikacji) Zamknij() {
	if a == nil {
		return
	}
	a.podglady.Zamknij()
}

// ZdefiniujArchitekture obsługuje `apps.architecture.define`. Brak
// `architectureId` zakłada architekturę nową; wskazanie zmienia zastaną i
// podnosi numer wersji — tak jak `automatyka.wersja` w module Automations.
// Komponenty i ich zależności wychodzą kontraktem w komplecie przy każdym
// zapisie (Architecture Designer nadsyła cały układ na nowo), więc warstwa
// danych wymienia je w jednej transakcji („usuń, wstaw od nowa") — jedna
// droga zapisu, nie dwie.
func (a *adapterAplikacji) ZdefiniujArchitekture(ctx context.Context,
	z shared.AppsArchitectureDefineRequest) (shared.AppsArchitectureDefineResponse, error) {

	if z.WindowId == "" {
		return shared.AppsArchitectureDefineResponse{}, bladWskazaniaAplikacji("żądanie bez okna modułu")
	}

	kod := wartoscTekstu(z.ArchitectureId)
	nowaArchitektura := kod == ""
	if nowaArchitektura {
		kod = nowyIdentyfikator(przedrostekArchitekturyApp)
	}

	zastane, err := a.zastaneZastrzezenia(ctx, kod, nowaArchitektura)
	if err != nil {
		return shared.AppsArchitectureDefineResponse{}, err
	}

	szablon := szablonArchitektury(z.Template)
	komponenty, zaleznosci := rozlozKomponenty(z.Components)

	zapisana, err := a.repozytorium.ZapiszArchitekture(ctx, dane.ArchitekturaApp{
		Kod: kod, Okno: z.WindowId, Nazwa: z.Name, Szablon: szablon,
		ZastrzezeniaWalidacji: zastane,
	}, komponenty, zaleznosci)
	if err != nil {
		return shared.AppsArchitectureDefineResponse{}, bladAplikacji(err)
	}

	architektura, err := a.zloz(ctx, zapisana)
	if err != nil {
		return shared.AppsArchitectureDefineResponse{}, err
	}
	return shared.AppsArchitectureDefineResponse{Architecture: architektura}, nil
}

// zastaneZastrzezenia odczytuje zastrzeżenia walidacji już zapisane przy
// architekturze, żeby zapis definicji ich nie skasował — adapter nie ma z
// czego wyliczyć nowych (patrz nagłówek pliku), więc jedyna uczciwa wartość
// to ta, która tam już była. Nowa architektura startuje z pustą listą.
func (a *adapterAplikacji) zastaneZastrzezenia(ctx context.Context, kod string, nowaArchitektura bool) ([]string, error) {
	if nowaArchitektura {
		return nil, nil
	}
	zastana, err := a.repozytorium.Architektura(ctx, kod)
	if err != nil {
		return nil, bladNieznanejArchitektury(kod, err)
	}
	return zastana.ZastrzezeniaWalidacji, nil
}

// zloz składa architekturę kontraktu z wiersza wraz z komponentami i ich
// zależnościami — Architecture Designer pokazuje układ w komplecie, nie samą nazwę.
func (a *adapterAplikacji) zloz(ctx context.Context, wiersz dane.ArchitekturaApp) (shared.AppArchitecture, error) {
	komponenty, err := a.repozytorium.Komponenty(ctx, wiersz.ID)
	if err != nil {
		return shared.AppArchitecture{}, bladAplikacji(err)
	}
	zaleznosci, err := a.repozytorium.ZaleznosciKomponentow(ctx, wiersz.ID)
	if err != nil {
		return shared.AppArchitecture{}, bladAplikacji(err)
	}
	zaleznosciZ := map[string][]string{}
	for _, zaleznosc := range zaleznosci {
		zaleznosciZ[zaleznosc.KomponentDo] = append(zaleznosciZ[zaleznosc.KomponentDo], zaleznosc.KomponentZ)
	}

	skladowe := make([]shared.AppComponent, 0, len(komponenty))
	for _, komponent := range komponenty {
		skladowe = append(skladowe, shared.AppComponent{
			Id: komponent.KodZewnetrzny, Name: komponent.Nazwa,
			Kind: shared.AppComponentKind(komponent.Rodzaj), Stack: komponent.Stos,
			Description: komponent.Opis, DependsOn: zaleznosciZ[komponent.KodZewnetrzny],
			ApiContract: komponent.KontraktAPI,
		})
	}

	wersja := wiersz.Wersja
	return shared.AppArchitecture{
		Id: wiersz.Kod, WindowId: wiersz.Okno, Name: wiersz.Nazwa,
		Template: shared.AppArchitectureTemplate(wiersz.Szablon), Components: skladowe,
		Version: &wersja, ValidationIssues: wiersz.ZastrzezeniaWalidacji,
		UpdatedAt: chwilaBazy(wiersz.Zaktualizowano),
	}, nil
}

// rozlozKomponenty rozkłada komponenty kontraktu na wiersze komponentów i
// osobne łuki zależności — tabela `zaleznosc_komponentu_apps` niesie graf
// niezależnie od komponentu (jedna prawda o krawędzi, nie dwie kopie
// w `DependsOn` obu końców).
func rozlozKomponenty(komponenty []shared.AppComponent) ([]dane.KomponentArchitektury, []dane.ZaleznoscKomponentu) {
	wiersze := make([]dane.KomponentArchitektury, 0, len(komponenty))
	zaleznosci := []dane.ZaleznoscKomponentu{}
	for _, komponent := range komponenty {
		wiersze = append(wiersze, dane.KomponentArchitektury{
			KodZewnetrzny: komponent.Id, Nazwa: komponent.Name,
			Rodzaj: string(komponent.Kind), Stos: komponent.Stack,
			Opis: komponent.Description, KontraktAPI: komponent.ApiContract,
		})
		for _, zrodlo := range komponent.DependsOn {
			zaleznosci = append(zaleznosci, dane.ZaleznoscKomponentu{
				KomponentZ: zrodlo, KomponentDo: komponent.Id,
			})
		}
	}
	return wiersze, zaleznosci
}

// szablonArchitektury rozstrzyga brak wskazania szablonu. Repozytorium ma
// własną wartość domyślną w UPSERT-cie („monolith"), ale adapter podaje ją
// jawnie, żeby złożona odpowiedź nie pokazywała pustego szablonu przy
// pierwszym zapisie zanim wiersz wróci z bazy.
func szablonArchitektury(wskazanie *shared.AppArchitectureTemplate) string {
	if wskazanie == nil {
		return "monolith"
	}
	return string(*wskazanie)
}

// bladAplikacji znakuje usterkę kodem kontraktu, żeby okno modułu pokazało
// powód, a nie samo „nie udało się". Błąd, któremu kod już nadano — odmowa
// wskazania, brak bytu — przechodzi tędy bez zmiany kodu; dopiero usterka bez
// kodu staje się usterką wewnętrzną rdzenia.
func bladAplikacji(err error) error {
	if err == nil {
		return nil
	}
	return protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeInternalError, err))
}

// bladWskazaniaAplikacji nazywa brak danych w żądaniu — to błąd Operatora,
// nie rdzenia, więc kod odmowy jest inny niż przy usterce.
func bladWskazaniaAplikacji(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"moduł Apps: "+powod))
}

// bladNieznanejArchitektury odróżnia „architektury nie ma" od „odczyt się nie
// powiódł". Okno pokazuje wtedy inny komunikat i inaczej podpowiada Operatorowi.
func bladNieznanejArchitektury(kod string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Apps: architektura nie istnieje: "+kod))
	}
	return bladAplikacji(err)
}
