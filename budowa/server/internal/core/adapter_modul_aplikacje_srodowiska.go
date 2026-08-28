// Moduł Apps obsługuje środowiska wdrożeniowe Deployment Panelu: wykaz i zmienne
// środowisk, domenę wraz z wpisami DNS, skalowanie instancji oraz kondycję
// produktu mierzoną zapytaniem HTTP albo stanem ostatniego wdrożenia.
package core

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// przedrostekSrodowiskaApp znakuje identyfikatory środowisk wdrożeniowych, tak
// jak każdy inny rodzaj rekordu w rdzeniu ma własny rozpoznawalny przedrostek.
const przedrostekSrodowiskaApp = "srod-"

// czasSprawdzeniaKondycjiApp jest granicą czekania na odpowiedź produktu przy
// sprawdzaniu kondycji oraz przy rozwiązywaniu nazwy domeny w systemie DNS.
const czasSprawdzeniaKondycjiApp = 5 * time.Second

// oknoPomiaruKondycjiApp mówi, z ilu ostatnich sprawdzeń liczy się udział
// dostępności. Bez granicy udział rozwadniałby się historią sprzed miesięcy
// i przestawał opisywać stan dzisiejszy.
const oknoPomiaruKondycjiApp = 50

// nazwySrodowiskWbudowanychApp niesie nazwy widoczne Operatorowi dla trzech
// środowisk kontraktu. Kod środowiska jest wartością kontraktu, nazwa — jego
// tłumaczeniem na język okna.
var nazwySrodowiskWbudowanychApp = []struct {
	Kod   shared.AppDeployEnvironment
	Nazwa string
}{
	{shared.AppDeployEnvironmentDev, "Deweloperskie"},
	{shared.AppDeployEnvironmentStaging, "Testowe"},
	{shared.AppDeployEnvironmentProduction, "Produkcyjne"},
}

// WypiszSrodowiska obsługuje `apps.environment.list`, zakładając najpierw
// brakujące wiersze trzech wbudowanych środowisk kontraktu i oddając wykaz.
func (a *adapterAplikacji) WypiszSrodowiska(ctx context.Context,
	z shared.AppsEnvironmentListRequest) (shared.AppsEnvironmentListResponse, error) {

	okno, err := oknoAplikacji(z.WindowId, "apps.environment.list")
	if err != nil {
		return shared.AppsEnvironmentListResponse{}, err
	}
	if err := a.zalozSrodowiskaWbudowaneApp(ctx, okno); err != nil {
		return shared.AppsEnvironmentListResponse{}, err
	}
	wiersze, err := a.repozytorium.SrodowiskaApp(ctx, okno)
	if err != nil {
		return shared.AppsEnvironmentListResponse{}, bladAplikacji(err)
	}
	srodowiska := make([]shared.AppEnvironment, 0, len(wiersze))
	for _, wiersz := range wiersze {
		kolejnosc := wiersz.Kolejnosc
		srodowiska = append(srodowiska, shared.AppEnvironment{
			Id: wiersz.Kod, WindowId: wiersz.Okno, Code: wiersz.KodSrodowiska,
			Name: wiersz.Nazwa, Order: &kolejnosc, Domain: wiersz.Domena,
		})
	}
	return shared.AppsEnvironmentListResponse{Environments: srodowiska, Total: len(srodowiska)}, nil
}

// zalozSrodowiskaWbudowaneApp dopisuje brakujące wiersze trzech środowisk
// kontraktu, nie tykając tego, co Operator im nadał.
func (a *adapterAplikacji) zalozSrodowiskaWbudowaneApp(ctx context.Context, okno string) error {
	for kolejnosc, wbudowane := range nazwySrodowiskWbudowanychApp {
		_, err := a.repozytorium.ZapiszSrodowiskoApp(ctx, dane.SrodowiskoApp{
			Kod: nowyIdentyfikator(przedrostekSrodowiskaApp), Okno: okno,
			KodSrodowiska: string(wbudowane.Kod), Nazwa: wbudowane.Nazwa,
			Kolejnosc: kolejnosc + 1,
		})
		if err != nil {
			return bladAplikacji(err)
		}
	}
	return nil
}

// WypiszZmienneSrodowiska obsługuje `apps.environment.variable.list`. Wartości
// sekretnych nie ma tu i być nie może: kolumna niesie klucz jawny, a treść
// poświadczenia zostaje w warstwie sekretów.
func (a *adapterAplikacji) WypiszZmienneSrodowiska(ctx context.Context,
	z shared.AppsEnvironmentVariableListRequest) (shared.AppsEnvironmentVariableListResponse, error) {

	okno, err := oknoAplikacji(z.WindowId, "apps.environment.variable.list")
	if err != nil {
		return shared.AppsEnvironmentVariableListResponse{}, err
	}
	if err := sprawdzSrodowiskoWdrozenia(z.Environment); err != nil {
		return shared.AppsEnvironmentVariableListResponse{}, err
	}
	wiersze, err := a.repozytorium.ZmienneSrodowiskaApp(ctx, okno, string(z.Environment))
	if err != nil {
		return shared.AppsEnvironmentVariableListResponse{}, bladAplikacji(err)
	}
	zmienne := make([]shared.AppEnvironmentVariable, 0, len(wiersze))
	for _, wiersz := range wiersze {
		zmienne = append(zmienne, zmiennaKontraktuApp(wiersz))
	}
	return shared.AppsEnvironmentVariableListResponse{
		Variables: zmienne, Total: len(zmienne),
	}, nil
}

// UstawZmiennaSrodowiska obsługuje `apps.environment.variable.set`, zapisując
// wartość jawną albo odwołanie do sekretu, wzajemnie się wykluczające.
func (a *adapterAplikacji) UstawZmiennaSrodowiska(ctx context.Context,
	z shared.AppsEnvironmentVariableSetRequest) (shared.AppsEnvironmentVariableSetResponse, error) {

	okno, err := oknoAplikacji(z.WindowId, "apps.environment.variable.set")
	if err != nil {
		return shared.AppsEnvironmentVariableSetResponse{}, err
	}
	if err := sprawdzSrodowiskoWdrozenia(z.Environment); err != nil {
		return shared.AppsEnvironmentVariableSetResponse{}, err
	}
	if strings.TrimSpace(z.Name) == "" {
		return shared.AppsEnvironmentVariableSetResponse{}, bladWskazaniaAplikacji(
			"apps.environment.variable.set wymaga nazwy zmiennej")
	}
	if z.Value != nil && z.SecretRef != nil {
		return shared.AppsEnvironmentVariableSetResponse{}, bladWskazaniaAplikacji(
			"zmienna " + z.Name + " niesie naraz wartość jawną i odwołanie do sekretu — " +
				"kontrakt wyklucza te dwa pola wzajemnie")
	}
	if z.Value == nil && z.SecretRef == nil {
		return shared.AppsEnvironmentVariableSetResponse{}, bladWskazaniaAplikacji(
			"zmienna " + z.Name + " nie niesie ani wartości jawnej, ani odwołania do sekretu")
	}

	zapisana, err := a.repozytorium.ZapiszZmiennaSrodowiskaApp(ctx, dane.ZmiennaSrodowiskaApp{
		Okno: okno, Srodowisko: string(z.Environment), Nazwa: z.Name,
		Wartosc: z.Value, OdwolanieSekretu: z.SecretRef,
	})
	if err != nil {
		return shared.AppsEnvironmentVariableSetResponse{}, bladAplikacji(err)
	}
	return shared.AppsEnvironmentVariableSetResponse{Variable: zmiennaKontraktuApp(zapisana)}, nil
}

// UstawDomene obsługuje `apps.deployment.domain.set`, zapisuje domenę wraz
// z wpisami DNS i sprawdza rozwiązanie nazwy, zamiast przyjąć ją na słowo.
func (a *adapterAplikacji) UstawDomene(ctx context.Context,
	z shared.AppsDeploymentDomainSetRequest) (shared.AppsDeploymentDomainSetResponse, error) {

	okno, err := oknoAplikacji(z.WindowId, "apps.deployment.domain.set")
	if err != nil {
		return shared.AppsDeploymentDomainSetResponse{}, err
	}
	if err := sprawdzSrodowiskoWdrozenia(z.Environment); err != nil {
		return shared.AppsDeploymentDomainSetResponse{}, err
	}
	domena := strings.TrimSpace(z.Domain)
	if domena == "" {
		return shared.AppsDeploymentDomainSetResponse{}, bladWskazaniaAplikacji(
			"apps.deployment.domain.set wymaga domeny")
	}
	if len(z.DnsRecords) > 0 && !json.Valid(z.DnsRecords) {
		return shared.AppsDeploymentDomainSetResponse{}, bladWskazaniaAplikacji(
			"wpisy DNS nie są poprawnym JSON-em")
	}

	if err := a.zalozSrodowiskaWbudowaneApp(ctx, okno); err != nil {
		return shared.AppsDeploymentDomainSetResponse{}, err
	}
	nazwa := string(z.Environment)
	for _, wbudowane := range nazwySrodowiskWbudowanychApp {
		if wbudowane.Kod == z.Environment {
			nazwa = wbudowane.Nazwa
		}
	}
	var wpisy *string
	if len(z.DnsRecords) > 0 {
		tresc := string(z.DnsRecords)
		wpisy = &tresc
	}
	if _, err := a.repozytorium.ZapiszSrodowiskoApp(ctx, dane.SrodowiskoApp{
		Kod: nowyIdentyfikator(przedrostekSrodowiskaApp), Okno: okno,
		KodSrodowiska: string(z.Environment), Nazwa: nazwa,
		Domena: &domena, WpisyDNS: wpisy,
	}); err != nil {
		return shared.AppsDeploymentDomainSetResponse{}, bladAplikacji(err)
	}

	potwierdzona := czyDomenaRozwiazujeApp(ctx, domena)
	a.dopiszDziennikApp(ctx, okno, nil, nil,
		"domena środowiska "+string(z.Environment)+" ustawiona na "+domena+
			", rozwiązanie nazwy: "+strconv.FormatBool(potwierdzona))
	return shared.AppsDeploymentDomainSetResponse{Domain: domena, Verified: potwierdzona}, nil
}

// czyDomenaRozwiazujeApp sprawdza, czy nazwa domeny rozwiązuje się na tej
// maszynie. Nierozwiązywalna nazwa nie jest odmową, bo domena bywa nadawana,
// zanim wpisy DNS zdążą się rozejść po sieci.
func czyDomenaRozwiazujeApp(ctx context.Context, domena string) bool {
	nazwa := domena
	nazwa = strings.TrimPrefix(strings.TrimPrefix(nazwa, "https://"), "http://")
	nazwa = strings.TrimSuffix(nazwa, "/")
	if host, _, err := net.SplitHostPort(nazwa); err == nil {
		nazwa = host
	}
	if nazwa == "" {
		return false
	}
	sprawdzenie, zakoncz := context.WithTimeout(ctx, czasSprawdzeniaKondycjiApp)
	defer zakoncz()
	adresy, err := net.DefaultResolver.LookupHost(sprawdzenie, nazwa)
	return err == nil && len(adresy) > 0
}

// UstawSkalowanie obsługuje `apps.deployment.scale.set` zapisem konfiguracji,
// a nie rozkazem dla orkiestratora. `effectiveInstances` oddaje instancje
// stałe, a gdy ich nie podano, dolną granicę skalowania.
func (a *adapterAplikacji) UstawSkalowanie(ctx context.Context,
	z shared.AppsDeploymentScaleSetRequest) (shared.AppsDeploymentScaleSetResponse, error) {

	okno, err := oknoAplikacji(z.WindowId, "apps.deployment.scale.set")
	if err != nil {
		return shared.AppsDeploymentScaleSetResponse{}, err
	}
	if err := sprawdzSrodowiskoWdrozenia(z.Environment); err != nil {
		return shared.AppsDeploymentScaleSetResponse{}, err
	}
	if z.Instances == nil && z.MinInstances == nil && z.MaxInstances == nil && len(z.Rules) == 0 {
		return shared.AppsDeploymentScaleSetResponse{}, bladWskazaniaAplikacji(
			"apps.deployment.scale.set nie niesie ani jednej nastawy do zapisania")
	}
	for nazwa, wartosc := range map[string]*int{
		"instances": z.Instances, "minInstances": z.MinInstances, "maxInstances": z.MaxInstances,
	} {
		if wartosc != nil && *wartosc < 0 {
			return shared.AppsDeploymentScaleSetResponse{}, bladWskazaniaAplikacji(
				"nastawa " + nazwa + " nie może być ujemna")
		}
	}
	if z.MinInstances != nil && z.MaxInstances != nil && *z.MinInstances > *z.MaxInstances {
		return shared.AppsDeploymentScaleSetResponse{}, bladWskazaniaAplikacji(
			"dolna granica skalowania (" + strconv.Itoa(*z.MinInstances) +
				") jest wyższa od górnej (" + strconv.Itoa(*z.MaxInstances) + ")")
	}
	if len(z.Rules) > 0 && !json.Valid(z.Rules) {
		return shared.AppsDeploymentScaleSetResponse{}, bladWskazaniaAplikacji(
			"reguły skalowania nie są poprawnym JSON-em")
	}

	var reguly *string
	if len(z.Rules) > 0 {
		tresc := string(z.Rules)
		reguly = &tresc
	}
	zapisana, err := a.repozytorium.ZapiszSkalowanieApp(ctx, dane.SkalowanieApp{
		Okno: okno, Srodowisko: string(z.Environment),
		Instancje:     liczbaDuzaApp(z.Instances),
		MinInstancji:  liczbaDuzaApp(z.MinInstances),
		MaksInstancji: liczbaDuzaApp(z.MaxInstances),
		Reguly:        reguly,
	})
	if err != nil {
		return shared.AppsDeploymentScaleSetResponse{}, bladAplikacji(err)
	}

	odpowiedz := shared.AppsDeploymentScaleSetResponse{Applied: true}
	if zapisana.Instancje != nil {
		ile := int(*zapisana.Instancje)
		odpowiedz.EffectiveInstances = &ile
	} else if zapisana.MinInstancji != nil {
		ile := int(*zapisana.MinInstancji)
		odpowiedz.EffectiveInstances = &ile
	}
	a.dopiszDziennikApp(ctx, okno, nil, nil,
		"nastawa skalowania środowiska "+string(z.Environment)+" zapisana")
	return odpowiedz, nil
}

// PobierzKondycje obsługuje `apps.deployment.health.get`, mierzy dostępność
// zapytaniem HTTP albo stanem ostatniego wdrożenia i dopisuje dziennik kondycji.
func (a *adapterAplikacji) PobierzKondycje(ctx context.Context,
	z shared.AppsDeploymentHealthGetRequest) (shared.AppsDeploymentHealthGetResponse, error) {

	okno, err := oknoAplikacji(z.WindowId, "apps.deployment.health.get")
	if err != nil {
		return shared.AppsDeploymentHealthGetResponse{}, err
	}
	if z.Environment != nil {
		if err := sprawdzSrodowiskoWdrozenia(*z.Environment); err != nil {
			return shared.AppsDeploymentHealthGetResponse{}, err
		}
	}

	wdrozenia, _, err := a.repozytorium.Wdrozenia(ctx, okno, z.Environment, granicaWdrozenApp)
	if err != nil {
		return shared.AppsDeploymentHealthGetResponse{}, bladAplikacji(err)
	}
	srodowisko := shared.AppDeployEnvironment(shared.AppDeployEnvironmentDev)
	switch {
	case z.Environment != nil:
		srodowisko = *z.Environment
	case len(wdrozenia) > 0:
		// Puste środowisko bierze ostatnie wdrożenie — wykaz wraca posortowany od najnowszego.
		srodowisko = wdrozenia[0].Srodowisko
	}

	dostepna, szczegol := a.zmierzDostepnoscApp(ctx, okno, srodowisko, wdrozenia)
	teraz := time.Now().UTC().UnixMilli()
	_ = a.repozytorium.ZapiszKondycjeApp(ctx, dane.KondycjaWdrozeniaApp{
		Okno: okno, Srodowisko: string(srodowisko), Dostepna: dostepna,
		Szczegol: wskaznikNapisuApp(szczegol), Sprawdzono: teraz,
	})

	kondycja := shared.AppDeploymentHealth{
		Environment: srodowisko,
		Available:   dostepna,
		LastCheckAt: teraz,
		Detail:      wskaznikNapisuApp(szczegol),
	}

	sprawdzenia, err := a.repozytorium.KondycjeApp(ctx, okno, string(srodowisko), oknoPomiaruKondycjiApp)
	if err != nil {
		return shared.AppsDeploymentHealthGetResponse{}, bladAplikacji(err)
	}
	if len(sprawdzenia) > 0 {
		udanych := 0
		for _, sprawdzenie := range sprawdzenia {
			if sprawdzenie.Dostepna {
				udanych++
			}
		}
		udzial := strconv.FormatFloat(float64(udanych)*100/float64(len(sprawdzenia)), 'f', 2, 64)
		kondycja.AvailabilityPercent = &udzial
	}

	// Czas działania liczy się od zakończenia ostatniego udanego wdrożenia tego środowiska.
	for _, wdrozenie := range wdrozenia {
		if wdrozenie.Srodowisko != srodowisko || wdrozenie.Stan != shared.AppDeployStatusSucceeded {
			continue
		}
		if wdrozenie.Zakonczono == nil {
			continue
		}
		sekundy := (teraz - chwilaBazy(*wdrozenie.Zakonczono)) / 1000
		if sekundy < 0 {
			sekundy = 0
		}
		kondycja.UptimeSeconds = &sekundy
		break
	}

	return shared.AppsDeploymentHealthGetResponse{Health: kondycja}, nil
}

// zmierzDostepnoscApp rozstrzyga dostępność produktu przez zapytanie do domeny
// albo do podglądu w oknie i nazywa źródło rozstrzygnięcia w szczególe.
func (a *adapterAplikacji) zmierzDostepnoscApp(ctx context.Context, okno string,
	srodowisko shared.AppDeployEnvironment, wdrozenia []dane.WdrozenieApp) (bool, string) {

	adres := ""
	if wiersz, err := a.repozytorium.SrodowiskoApp(ctx, okno, string(srodowisko)); err == nil &&
		wiersz.Domena != nil && *wiersz.Domena != "" {
		adres = zAdresemHttpApp(*wiersz.Domena)
	}
	if adres == "" {
		adres = a.adresPodgladuApp(okno)
	}

	if adres != "" {
		sprawdzenie, zakoncz := context.WithTimeout(ctx, czasSprawdzeniaKondycjiApp)
		defer zakoncz()
		zapytanie, err := http.NewRequestWithContext(sprawdzenie, http.MethodGet, adres, nil)
		if err != nil {
			return false, "nie można złożyć zapytania kondycji do " + adres + ": " + err.Error()
		}
		klient := &http.Client{Timeout: czasSprawdzeniaKondycjiApp}
		poczatek := time.Now()
		odpowiedz, err := klient.Do(zapytanie)
		if err != nil {
			return false, "produkt pod " + adres + " nie odpowiada: " + err.Error()
		}
		defer odpowiedz.Body.Close()
		czas := time.Since(poczatek).Milliseconds()
		dostepna := odpowiedz.StatusCode < 500
		return dostepna, "sprawdzenie " + adres + " → " + strconv.Itoa(odpowiedz.StatusCode) +
			" w " + strconv.FormatInt(czas, 10) + " ms"
	}

	// Bez adresu do zapytania dostępność wynika ze stanu ostatniego wdrożenia środowiska.
	for _, wdrozenie := range wdrozenia {
		if wdrozenie.Srodowisko != srodowisko {
			continue
		}
		return wdrozenie.Stan == shared.AppDeployStatusSucceeded,
			"produkt nie ma nadanej domeny ani stojącego podglądu — dostępność wzięta ze stanu " +
				"ostatniego wdrożenia " + wdrozenie.Kod + " (" + string(wdrozenie.Stan) + ")"
	}
	return false, "środowisko " + string(srodowisko) + " nie ma ani jednego wdrożenia, " +
		"nadanej domeny ani stojącego podglądu — nie ma czego sprawdzić"
}

// zmiennaKontraktuApp przekłada wiersz zmiennej środowiskowej z bazy danych
// na kształt zmiennej środowiskowej z kontraktu, wraz z chwilą aktualizacji.
func zmiennaKontraktuApp(wiersz dane.ZmiennaSrodowiskaApp) shared.AppEnvironmentVariable {
	return shared.AppEnvironmentVariable{
		Name:        wiersz.Nazwa,
		Environment: shared.AppDeployEnvironment(wiersz.Srodowisko),
		Value:       wiersz.Wartosc,
		SecretRef:   wiersz.OdwolanieSekretu,
		UpdatedAt:   chwilaBazy(wiersz.Zaktualizowano),
	}
}

// liczbaDuzaApp przekłada wskaźnik do liczby całkowitej z kontraktu na
// wskaźnik do liczby całkowitej o szerszym zakresie, jakiego wymaga baza.
func liczbaDuzaApp(wartosc *int) *int64 {
	if wartosc == nil {
		return nil
	}
	duza := int64(*wartosc)
	return &duza
}
