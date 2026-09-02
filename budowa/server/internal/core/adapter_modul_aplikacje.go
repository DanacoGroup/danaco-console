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

type adapterAplikacji struct {
	repozytorium      dane.RepozytoriumAplikacji
	przyrostWdrozenia func(context.Context, shared.ChangeKind, shared.AppDeployment)
	przyrostWarsztatu func(context.Context, shared.ChangeKind, string, shared.AppWorkspaceLayer, shared.DeveloperFile)
	okna              *session.Rejestr
	przyrostEtapu     func(context.Context, shared.ChangeKind, shared.AppStage)
	magazyn           *magazynTresciBiblioteki
	katalogDanych     string
	sejf              sejfKluczaWydawcy
	katalogRozszerzen dane.RepozytoriumRozszerzen
	podglady          *rejestrPodgladowApp
	uruchamiacz       session.Uruchamiacz
	rozstrzygacz      *konfig.Rozstrzygacz
	katalog           *KatalogRoboczy
}

// sejfKluczaWydawcy jest wycinkiem sejfu poświadczeń, którego moduł potrzebuje. Wycinek, nie cały sejf: `apps.package.sign` czyta i — przy pierwszym użyciu odwołania — zakłada klucz wydawcy, a kasowanie poświadczeń do modułu nie należy.
type sejfKluczaWydawcy interface {
	Odczytaj(ctx context.Context, byt string) (string, bool)
	Zapisz(ctx context.Context, byt, poswiadczenie string) (string, error)
}

func (a *adapterAplikacji) ZMagazynem(katalogDanych string) *adapterAplikacji {
	a.katalogDanych = katalogDanych
	a.magazyn = magazynWytworowApp(katalogDanych)
	a.sejf = dane.NowySejfPlikowy(katalogDanych)
	return a
}

func (a *adapterAplikacji) ZKatalogiemRozszerzen(rejestr dane.RepozytoriumRozszerzen) *adapterAplikacji {
	a.katalogRozszerzen = rejestr
	return a
}

func (a *adapterAplikacji) PodepnijPrzyrostEtapu(rozglos func(context.Context, shared.ChangeKind, shared.AppStage)) {
	a.przyrostEtapu = rozglos
}

func (a *adapterAplikacji) ZUruchamiaczem(uruchamiacz session.Uruchamiacz,
	rozstrzygacz *konfig.Rozstrzygacz, katalog *KatalogRoboczy) *adapterAplikacji {

	a.uruchamiacz, a.rozstrzygacz, a.katalog = uruchamiacz, rozstrzygacz, katalog
	return a
}

func (a *adapterAplikacji) ZOknami(rejestr *session.Rejestr) *adapterAplikacji {
	a.okna = rejestr
	return a
}

func (a *adapterAplikacji) PodepnijPrzyrostWarsztatu(
	rozglos func(context.Context, shared.ChangeKind, string, shared.AppWorkspaceLayer, shared.DeveloperFile)) {

	a.przyrostWarsztatu = rozglos
}

func (a *adapterAplikacji) rozglosWarsztat(ctx context.Context, zmiana shared.ChangeKind, okno string,
	warstwa shared.AppWorkspaceLayer, plik shared.DeveloperFile) {

	if a.przyrostWarsztatu == nil {
		return
	}
	a.przyrostWarsztatu(ctx, zmiana, okno, warstwa, plik)
}

func (a *adapterAplikacji) PodepnijPrzyrostWdrozenia(
	rozglos func(context.Context, shared.ChangeKind, shared.AppDeployment)) {

	a.przyrostWdrozenia = rozglos
}

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
	if errors.Is(err, dane.ErrKolizjaWiersza) {
		return protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeConflict, err))
	}
	return protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeInternalError, err))
}

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
