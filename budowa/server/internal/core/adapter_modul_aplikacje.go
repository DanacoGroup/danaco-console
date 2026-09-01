// Moduł Apps: typ adaptera, konstruktor i komenda `apps.architecture.define`. Warsztat i wdrożenia stoją w adapter_modul_aplikacje_wdrozenie.go i adapter_modul_aplikacje_uchwyty.go. Zastrzeżenia walidacji są przechowywane, nie wyliczane.
package core

import (
	"context"
	"errors"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// Przedrostki identyfikatorów bytów modułu. Zadeklarowane tu w całości, łącznie z przedrostkiem wdrożenia, którego ten plik nie używa — jedno miejsce nadawania przedrostków modułu zamiast kilku.
const (
	przedrostekArchitekturyApp = "arch-"
	przedrostekKomponentuApp   = "kmp-"
	przedrostekWdrozeniaApp    = "wdr-"
)

// adapterAplikacji wypełnia część portu Aplikacje. Zależność jest jedna: wspólne repozytorium modułu, rozłożone po stronie danych na trzy pliki wedle odpowiedzialności, ale niosące jeden typ.
type adapterAplikacji struct {
	repozytorium dane.RepozytoriumAplikacji
	// przyrostWdrozenia rozgłasza apps.build.changed; podpina go obsługiwacz zmian modułu.
	przyrostWdrozenia func(context.Context, shared.ChangeKind, shared.AppDeployment)
	// przyrostWarsztatu rozgłasza apps.workspace.changed drugiemu oknu tej samej przestrzeni.
	przyrostWarsztatu func(context.Context, shared.ChangeKind, string, shared.AppWorkspaceLayer, shared.DeveloperFile)
	// okna to rejestr okien sesji; wdrożenie bierze z okna przestrzeń roboczą do wysłania.
	okna *session.Rejestr
	// przyrostEtapu rozgłasza apps.build.changed przy zmianie etapu budowy Product Buildera.
	przyrostEtapu func(context.Context, shared.ChangeKind, shared.AppStage)
	// magazyn trzyma bajty wytworów modułu: eksport diagramu, artefakt budowania, archiwum pakietu.
	magazyn *magazynTresciBiblioteki
	// katalogDanych jest korzeniem, względem którego liczone są odwołania magazynu.
	katalogDanych string
	// sejf wydaje klucz wydawcy po kluczu jawnym signingKeyRef; materiał nie przechodzi przez kontrakt.
	sejf sejfKluczaWydawcy
	// katalogRozszerzen jest rejestrem pozycji katalogu dla apps.package.publish.
	katalogRozszerzen dane.RepozytoriumRozszerzen
	// podglady trzymają stojące serwery podglądu warstw; rejestr żyje wyłącznie w pamięci.
	podglady *rejestrPodgladowApp
	// uruchamiacz jest portem startu procesu; moduł sięga po niego w audycie wydajności strony.
	uruchamiacz session.Uruchamiacz
	// rozstrzygacz i katalog składają zasady izolacji okna oraz katalog startu procesu.
	rozstrzygacz *konfig.Rozstrzygacz
	katalog      *KatalogRoboczy
}

// sejfKluczaWydawcy jest wycinkiem sejfu poświadczeń, którego moduł potrzebuje. Wycinek, nie cały sejf: `apps.package.sign` czyta i — przy pierwszym użyciu odwołania — zakłada klucz wydawcy, a kasowanie poświadczeń do modułu nie należy.
type sejfKluczaWydawcy interface {
	Odczytaj(ctx context.Context, byt string) (string, bool)
	Zapisz(ctx context.Context, byt, poswiadczenie string) (string, error)
}

// ZMagazynem wpina katalog danych rdzenia: magazyn bajtów wytworów modułu oraz sejf, z którego bierze się klucz wydawcy. Bez niego eksport, pakowanie i podpis odmawiają z powodem, zamiast meldować wytwór bez bajtów.
func (a *adapterAplikacji) ZMagazynem(katalogDanych string) *adapterAplikacji {
	a.katalogDanych = katalogDanych
	a.magazyn = magazynWytworowApp(katalogDanych)
	a.sejf = dane.NowySejfPlikowy(katalogDanych)
	return a
}

// ZKatalogiemRozszerzen wpina rejestr pozycji katalogu — drogę `apps.package.publish` do prywatnego rejestru organizacji.
func (a *adapterAplikacji) ZKatalogiemRozszerzen(rejestr dane.RepozytoriumRozszerzen) *adapterAplikacji {
	a.katalogRozszerzen = rejestr
	return a
}

// PodepnijPrzyrostEtapu oddaje adapterowi drogę do `apps.build.changed` przy zmianie etapu budowy Product Buildera, oddzielną od wdrożenia.
func (a *adapterAplikacji) PodepnijPrzyrostEtapu(rozglos func(context.Context, shared.ChangeKind, shared.AppStage)) {
	a.przyrostEtapu = rozglos
}

// ZUruchamiaczem wpina port startu procesu wraz z zasadami izolacji okna. Bez niego `apps.performance.audit` odmawia zdaniem nazywającym brak, zamiast oddawać ocenę wyliczoną bez pomiaru.
func (a *adapterAplikacji) ZUruchamiaczem(uruchamiacz session.Uruchamiacz,
	rozstrzygacz *konfig.Rozstrzygacz, katalog *KatalogRoboczy) *adapterAplikacji {

	a.uruchamiacz, a.rozstrzygacz, a.katalog = uruchamiacz, rozstrzygacz, katalog
	return a
}

// ZOknami wpina rejestr okien sesji. Bez niego `apps.deployment.run` zakłada przebieg dla okna, którego nie ma, i kończy go powodem o pustym warsztacie zamiast o braku okna.
func (a *adapterAplikacji) ZOknami(rejestr *session.Rejestr) *adapterAplikacji {
	a.okna = rejestr
	return a
}

// PodepnijPrzyrostWarsztatu oddaje adapterowi drogę do `apps.workspace.changed`. Osobna droga od wdrożeniowej, bo zdarzenie niesie inny byt (plik warstwy, nie przebieg) — jedna funkcja o dwóch znaczeniach byłaby dwiema prawdami.
func (a *adapterAplikacji) PodepnijPrzyrostWarsztatu(
	rozglos func(context.Context, shared.ChangeKind, string, shared.AppWorkspaceLayer, shared.DeveloperFile)) {

	a.przyrostWarsztatu = rozglos
}

// rozglosWarsztat oddaje zmianę pliku obsługiwaczowi. Brak podpięcia nie zmienia pracy modułu — rdzeń zapisuje także wtedy, gdy nikt nie słucha zdarzeń.
func (a *adapterAplikacji) rozglosWarsztat(ctx context.Context, zmiana shared.ChangeKind, okno string,
	warstwa shared.AppWorkspaceLayer, plik shared.DeveloperFile) {

	if a.przyrostWarsztatu == nil {
		return
	}
	a.przyrostWarsztatu(ctx, zmiana, okno, warstwa, plik)
}

// PodepnijPrzyrostWdrozenia oddaje adapterowi drogę do zdarzenia zmiany wdrożenia. Przebieg zakładany jest w stanie `pending`, a silnik przesuwa go przez `running` do końcowego stanu już po odesłaniu odpowiedzi.
func (a *adapterAplikacji) PodepnijPrzyrostWdrozenia(
	rozglos func(context.Context, shared.ChangeKind, shared.AppDeployment)) {

	a.przyrostWdrozenia = rozglos
}

// nowyAdapterAplikacji wiąże adapter z repozytorium modułu. Rejestr stojących podglądów powstaje od razu — `Zamknij` musi mieć co zamknąć także wtedy, gdy żaden podgląd nie ruszył.
func nowyAdapterAplikacji(repozytorium dane.RepozytoriumAplikacji) *adapterAplikacji {
	return &adapterAplikacji{repozytorium: repozytorium, podglady: nowyRejestrPodgladowApp()}
}

// Zamknij zatrzymuje serwery podglądu podniesione przez `apps.preview.start`. Nasłuch żyje poza żądaniem, więc bez tego przeżyłby zatrzymanie rdzenia i zostawił zajęty port — tak samo jak przebieg budowania modułu Developer.
func (a *adapterAplikacji) Zamknij() {
	if a == nil {
		return
	}
	a.podglady.Zamknij()
}

// ZdefiniujArchitekture obsługuje `apps.architecture.define`. Brak `architectureId` zakłada architekturę nową; wskazanie zmienia zastaną i podnosi numer wersji. Komponenty i zależności wychodzą kontraktem w komplecie i zastępują poprzedni zapis.
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

// zastaneZastrzezenia odczytuje zastrzeżenia walidacji zapisane przy architekturze, żeby zapis definicji ich nie skasował — adapter nowych nie wylicza, więc jedyna uczciwa wartość to ta zapisana wcześniej. Nowa architektura startuje z pustą listą.
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

// zloz składa architekturę kontraktu z wiersza wraz z komponentami i ich zależnościami — Architecture Designer pokazuje układ w komplecie, nie samą nazwę.
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

// rozlozKomponenty rozkłada komponenty kontraktu na wiersze komponentów i osobne łuki zależności — tabela `zaleznosc_komponentu_apps` niesie graf niezależnie od komponentu (jedna prawda o krawędzi, nie dwie kopie w `DependsOn` obu końców).
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

// szablonArchitektury rozstrzyga brak wskazania szablonu. Repozytorium ma własną wartość domyślną („monolith"), ale adapter podaje ją jawnie, żeby odpowiedź nie pokazywała pustego szablonu przed powrotem wiersza z bazy.
func szablonArchitektury(wskazanie *shared.AppArchitectureTemplate) string {
	if wskazanie == nil {
		return "monolith"
	}
	return string(*wskazanie)
}

// bladAplikacji znakuje usterkę kodem kontraktu, żeby okno modułu pokazało powód, a nie samo „nie udało się". Błąd z kodem już nadanym przechodzi bez zmiany; dopiero usterka bez kodu staje się usterką wewnętrzną rdzenia.
func bladAplikacji(err error) error {
	if err == nil {
		return nil
	}
	return protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeInternalError, err))
}

// bladWskazaniaAplikacji nazywa brak danych w żądaniu — to błąd Operatora, nie rdzenia, więc kod odmowy jest inny niż przy usterce.
func bladWskazaniaAplikacji(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"moduł Apps: "+powod))
}

// bladNieznanejArchitektury odróżnia „architektury nie ma" od „odczyt się nie powiódł". Okno pokazuje wtedy inny komunikat i inaczej podpowiada Operatorowi.
func bladNieznanejArchitektury(kod string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Apps: architektura nie istnieje: "+kod))
	}
	return bladAplikacji(err)
}
