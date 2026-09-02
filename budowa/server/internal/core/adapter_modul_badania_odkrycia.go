// Moduł Discovery Panel obsługuje wyszukiwanie u dostawców, pomoc w układaniu
// zapytań, snowballing cytowań, odrzucanie pozycji, monitory tematów
// i kanałów oraz import wsadowy adresów.
package core

import (
	"context"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// limitOdkryciaBadania jest domyślną liczbą pozycji pobieranych z jednego
// dostawcy, gdy żądanie nie wskazuje własnej granicy.
const limitOdkryciaBadania = 20

// SzukajZrodel obsługuje `research.discovery.search` i pyta wielu dostawców
// naraz, zapamiętując każdą zwróconą pozycję w bazie pod jej kluczem.
func (a *adapterBadan) SzukajZrodel(ctx context.Context,
	z shared.ResearchDiscoverySearchRequest) (shared.ResearchDiscoverySearchResponse, error) {

	if z.WindowId == "" {
		return shared.ResearchDiscoverySearchResponse{}, bladWskazaniaBadan("discovery.search bez okna badania")
	}
	if strings.TrimSpace(z.Query) == "" {
		return shared.ResearchDiscoverySearchResponse{}, bladWskazaniaBadan("discovery.search bez zapytania")
	}
	limit := limitOdkryciaBadania
	if z.Limit != nil && *z.Limit > 0 {
		limit = *z.Limit
	}
	tylkoOtwarte := z.OpenAccessOnly != nil && *z.OpenAccessOnly

	wyniki := []shared.ResearchDiscoveryResult{}
	uzyci := []string{}
	zawiedli := []string{}
	razem := 0

	dodaj := func(nazwa string, pozycje []shared.ResearchDiscoveryResult, suma int, err error) {
		if err != nil {
			zawiedli = append(zawiedli, nazwa)
			return
		}
		uzyci = append(uzyci, nazwa)
		wyniki = append(wyniki, pozycje...)
		razem += suma
	}

	for _, dostawca := range dostawcyOdkryciaBadania(z.Mode, z.Providers) {
		switch dostawca {
		case "crossref":
			pozycje, suma, err := szukajCrossref(ctx, z.Query, z.YearFrom, z.YearTo, limit)
			dodaj(dostawca, pozycje, suma, err)
		case "openalex":
			pozycje, suma, err := szukajOpenAlex(ctx, z.Query, z.YearFrom, z.YearTo, tylkoOtwarte, limit)
			dodaj(dostawca, pozycje, suma, err)
		case "arxiv":
			pozycje, err := szukajArxiv(ctx, z.Query, limit)
			dodaj(dostawca, pozycje, len(pozycje), err)
		case "duckduckgo":
			pozycje, err := szukajWeb(ctx, z.Query, limit)
			dodaj(dostawca, pozycje, len(pozycje), err)
		}
	}
	if len(uzyci) == 0 {
		return shared.ResearchDiscoverySearchResponse{}, protokolBladBadania(
			shared.ErrorCodeChannelUnavailable,
			"żaden z dostawców wyszukiwania nie odpowiedział ("+strings.Join(zawiedli, ", ")+
				") — naprawa: sprawdzić łącze albo wskazać innych dostawców polem providers")
	}

	wzbogacone, err := a.zapamietajWynikiBadania(ctx, z.WindowId, wyniki)
	if err != nil {
		return shared.ResearchDiscoverySearchResponse{}, err
	}
	if razem < len(wzbogacone) {
		razem = len(wzbogacone)
	}
	return shared.ResearchDiscoverySearchResponse{
		Results: wzbogacone, Total: razem, ProvidersUsed: uzyci, ProvidersFailed: zawiedli,
	}, nil
}

// dostawcyOdkryciaBadania wybiera dostawców właściwych trybowi wyszukiwania.
// Wskazanie wprost ma pierwszeństwo — Operator, który wskazał dostawcę, dostaje
// jego wyniki, a nie zestaw domyślny obok nich.
func dostawcyOdkryciaBadania(tryb shared.ResearchDiscoveryMode, wskazani []string) []string {
	if len(wskazani) > 0 {
		return wskazani
	}
	if tryb == shared.ResearchDiscoveryModeWeb {
		return []string{"duckduckgo"}
	}
	return []string{"crossref", "openalex", "arxiv"}
}

// zapamietajWynikiBadania utrwala pozycje i dokłada im wiedzę, której dostawca
// nie ma: czy pozycja jest już źródłem badania i czy została odrzucona.
func (a *adapterBadan) zapamietajWynikiBadania(ctx context.Context, okno string,
	wyniki []shared.ResearchDiscoveryResult) ([]shared.ResearchDiscoveryResult, error) {

	zrodla, err := a.repozytorium.Zrodla(ctx, okno)
	if err != nil {
		return nil, bladBadan(err)
	}
	poAdresie := map[string]string{}
	for _, zrodlo := range zrodla {
		if zrodlo.Adres != nil && strings.TrimSpace(*zrodlo.Adres) != "" {
			poAdresie[strings.ToLower(strings.TrimSpace(*zrodlo.Adres))] = zrodlo.Kod
		}
	}
	zastane, err := a.repozytorium.WynikiOdkrycia(ctx, okno)
	if err != nil {
		return nil, bladBadan(err)
	}
	znane := map[string]dane.WynikOdkryciaBadania{}
	for _, wynik := range zastane {
		znane[wynik.Klucz] = wynik
	}

	doZapisu := make([]dane.WynikOdkryciaBadania, 0, len(wyniki))
	wzbogacone := make([]shared.ResearchDiscoveryResult, 0, len(wyniki))
	widziane := map[string]bool{}
	for _, wynik := range wyniki {
		if widziane[wynik.Key] {
			continue
		}
		widziane[wynik.Key] = true

		wiersz := dane.WynikOdkryciaBadania{
			Klucz: wynik.Key, Okno: okno, Tytul: wynik.Title, Adres: wynik.Url,
			Autorzy: wynik.Authors, Dostawca: wynik.Provider, Identyfikator: wynik.Identifier,
			Fragment: wynik.Snippet,
		}
		if wynik.Year != nil {
			rok := int64(*wynik.Year)
			wiersz.Rok = &rok
		}
		if wynik.OpenAccess != nil {
			otwarty := int64(0)
			if *wynik.OpenAccess {
				otwarty = 1
			}
			wiersz.OtwartyDostep = &otwarty
		}
		if wynik.Url != nil {
			if kod, jest := poAdresie[strings.ToLower(strings.TrimSpace(*wynik.Url))]; jest {
				wiersz.ZrodloKod = &kod
				wiersz.Duplikat = true
				wynik.AlreadySourceId = &kod
			}
		}
		if zapamietany, jest := znane[wynik.Key]; jest && zapamietany.ZrodloKod != nil {
			wynik.AlreadySourceId = zapamietany.ZrodloKod
		}
		doZapisu = append(doZapisu, wiersz)
		wzbogacone = append(wzbogacone, wynik)
	}
	if err := a.repozytorium.ZapiszWynikiOdkrycia(ctx, doZapisu); err != nil {
		return nil, bladBadan(err)
	}
	return wzbogacone, nil
}

// UlozZapytania obsługuje `research.discovery.assist` i przekłada pytanie
// badawcze na zestaw zapytań wyszukiwawczych wraz z uzasadnieniem doboru.
func (a *adapterBadan) UlozZapytania(ctx context.Context,
	z shared.ResearchDiscoveryAssistRequest) (shared.ResearchDiscoveryAssistResponse, error) {

	if strings.TrimSpace(z.Question) == "" {
		return shared.ResearchDiscoveryAssistResponse{}, bladWskazaniaBadan("discovery.assist bez pytania")
	}
	kanal, err := a.domyslnyKanalBadania(ctx)
	if err != nil {
		return shared.ResearchDiscoveryAssistResponse{}, err
	}
	tryb := "naukowego"
	if z.Mode != nil && *z.Mode == shared.ResearchDiscoveryModeWeb {
		tryb = "webowego"
	}
	polecenie := "Przekształć pytanie badawcze w zestaw zapytań wyszukiwawczych dla wyszukiwania " +
		tryb + ". Użyj operatorów logicznych i synonimów. Odpowiedz WYŁĄCZNIE obiektem JSON " +
		"o polach \"zapytania\" (lista napisów) i \"uzasadnienie\" (jedno zdanie).\n\nPytanie: " +
		z.Question

	odpowiedz, err := a.zapytajModel(ctx, "research.discovery.assist", kanal, polecenie)
	if err != nil {
		return shared.ResearchDiscoveryAssistResponse{}, err
	}
	var rozlozone struct {
		Zapytania    []string `json:"zapytania"`
		Uzasadnienie string   `json:"uzasadnienie"`
	}
	if err := jsonModeluBadania(odpowiedz, &rozlozone); err != nil || len(rozlozone.Zapytania) == 0 {
		// Odpowiedź nierozłożona nie przepada: każdy niepusty wiersz jest kandydatem.
		rozlozone.Zapytania = nil
		for _, wiersz := range strings.Split(odpowiedz, "\n") {
			wiersz = strings.TrimSpace(strings.TrimLeft(wiersz, "-*0123456789. "))
			if wiersz != "" {
				rozlozone.Zapytania = append(rozlozone.Zapytania, wiersz)
			}
		}
	}
	odpowiedzKontraktu := shared.ResearchDiscoveryAssistResponse{Queries: rozlozone.Zapytania}
	if strings.TrimSpace(rozlozone.Uzasadnienie) != "" {
		uzasadnienie := rozlozone.Uzasadnienie
		odpowiedzKontraktu.Rationale = &uzasadnienie
	}
	return odpowiedzKontraktu, nil
}

// RozwinCytowania obsługuje `research.discovery.snowball` i rozwija sieć
// cytowań pozycji wstecz, wprzód albo w obie strony naraz.
func (a *adapterBadan) RozwinCytowania(ctx context.Context,
	z shared.ResearchDiscoverySnowballRequest) (shared.ResearchDiscoverySnowballResponse, error) {

	identyfikator := wartoscTekstu(z.Identifier)
	if identyfikator == "" && z.SourceId != nil && *z.SourceId != "" {
		lektura, err := a.repozytorium.LekturaZrodlaBadania(ctx, *z.SourceId)
		if err != nil {
			return shared.ResearchDiscoverySnowballResponse{}, bladBadan(err)
		}
		if lektura.Identyfikat != nil {
			identyfikator = *lektura.Identyfikat
		}
	}
	if strings.TrimSpace(identyfikator) == "" {
		return shared.ResearchDiscoverySnowballResponse{}, bladWskazaniaBadan(
			"discovery.snowball bez identyfikatora pozycji — snowballing idzie po DOI, " +
				"więc źródło bez identyfikatora trzeba najpierw rozstrzygnąć przez research.source.resolve")
	}
	limit := limitOdkryciaBadania
	if z.Limit != nil && *z.Limit > 0 {
		limit = *z.Limit
	}

	odpowiedz := shared.ResearchDiscoverySnowballResponse{}
	wstecz := z.Direction == shared.ResearchSnowballDirectionBackward ||
		z.Direction == shared.ResearchSnowballDirectionBoth
	wprzod := z.Direction == shared.ResearchSnowballDirectionForward ||
		z.Direction == shared.ResearchSnowballDirectionBoth

	if wstecz {
		praca, err := pracaCrossref(ctx, identyfikator)
		if err != nil {
			return shared.ResearchDiscoverySnowballResponse{}, err
		}
		for _, odwolanie := range praca.Reference {
			if strings.TrimSpace(odwolanie.DOI) == "" || len(odpowiedz.Referenced) >= limit {
				continue
			}
			cytowana, err := pracaCrossref(ctx, odwolanie.DOI)
			if err != nil {
				continue
			}
			odpowiedz.Referenced = append(odpowiedz.Referenced, cytowana.jakoWynikBadania())
		}
	}
	if wprzod {
		praca, err := pracaOpenAlex(ctx, identyfikator)
		if err != nil {
			return shared.ResearchDiscoverySnowballResponse{}, err
		}
		if strings.TrimSpace(praca.CitedByAPIURL) != "" {
			cytujace, _, err := wynikiOpenAlexZAdresuBadania(ctx, praca.CitedByAPIURL, limit)
			if err == nil {
				odpowiedz.Citing = cytujace
			}
		}
	}
	if !wstecz && !wprzod {
		return shared.ResearchDiscoverySnowballResponse{}, bladWskazaniaBadan(
			"discovery.snowball bez kierunku — wskaż backward, forward albo both")
	}
	return odpowiedz, nil
}

// OdrzucWyniki obsługuje `research.discovery.reject` i odrzuca wskazane
// pozycje z uzasadnieniem, aktualizując liczniki diagramu PRISMA.
func (a *adapterBadan) OdrzucWyniki(ctx context.Context,
	z shared.ResearchDiscoveryRejectRequest) (shared.ResearchDiscoveryRejectResponse, error) {

	if z.WindowId == "" {
		return shared.ResearchDiscoveryRejectResponse{}, bladWskazaniaBadan("discovery.reject bez okna badania")
	}
	if len(z.ResultKeys) == 0 {
		return shared.ResearchDiscoveryRejectResponse{},
			bladWskazaniaBadan("discovery.reject bez pozycji — nie ma czego odrzucić")
	}
	if strings.TrimSpace(z.Reason) == "" {
		return shared.ResearchDiscoveryRejectResponse{}, bladWskazaniaBadan(
			"discovery.reject bez uzasadnienia — przesiew bez powodu nie da się udokumentować " +
				"w diagramie PRISMA")
	}
	odrzucone, err := a.repozytorium.OdrzucWynikiOdkrycia(ctx, z.WindowId, z.ResultKeys, z.Reason)
	if err != nil {
		return shared.ResearchDiscoveryRejectResponse{}, bladBadan(err)
	}
	liczniki, err := a.licznikiPrismaBadania(ctx, z.WindowId)
	if err != nil {
		return shared.ResearchDiscoveryRejectResponse{}, err
	}
	return shared.ResearchDiscoveryRejectResponse{Rejected: odrzucone, Counts: liczniki}, nil
}

// ── Monitory tematów i kanałów ─────────────────────────────────────────────

// UstawMonitor obsługuje `research.monitor.set` i zakłada monitor tematu
// albo kanału nowy, albo zmienia zastany po wskazanym identyfikatorze.
func (a *adapterBadan) UstawMonitor(ctx context.Context,
	z shared.ResearchMonitorSetRequest) (shared.ResearchMonitorSetResponse, error) {

	if z.WindowId == "" {
		return shared.ResearchMonitorSetResponse{}, bladWskazaniaBadan("monitor.set bez okna badania")
	}
	if z.Kind == "" {
		return shared.ResearchMonitorSetResponse{}, bladWskazaniaBadan("monitor.set bez rodzaju monitora")
	}
	if z.Kind == shared.ResearchMonitorKindTopic && (z.Query == nil || strings.TrimSpace(*z.Query) == "") {
		return shared.ResearchMonitorSetResponse{},
			bladWskazaniaBadan("monitor tematu bez zapytania — nie ma czego pilnować")
	}
	if z.Kind == shared.ResearchMonitorKindFeed && (z.Url == nil || strings.TrimSpace(*z.Url) == "") {
		return shared.ResearchMonitorSetResponse{},
			bladWskazaniaBadan("monitor kanału bez adresu — nie ma czego czytać")
	}

	kod := nowyIdentyfikator(przedrostekMonitoraBadania)
	if z.MonitorId != nil && *z.MonitorId != "" {
		kod = *z.MonitorId
	}
	wlaczony := true
	if z.Enabled != nil {
		wlaczony = *z.Enabled
	}
	monitor := dane.MonitorBadania{
		Kod: kod, Okno: z.WindowId, Rodzaj: string(z.Kind),
		Zapytanie: z.Query, Adres: z.Url, Wlaczony: wlaczony,
	}
	if z.IntervalMinutes != nil {
		interwal := int64(*z.IntervalMinutes)
		monitor.InterwalMinut = &interwal
	}
	zapisany, err := a.repozytorium.ZapiszMonitor(ctx, monitor)
	if err != nil {
		return shared.ResearchMonitorSetResponse{}, bladBadan(err)
	}
	return shared.ResearchMonitorSetResponse{Monitor: zlozMonitorBadania(zapisany)}, nil
}

// WypiszMonitory obsługuje `research.monitor.list` i oddaje monitory okna,
// z możliwością zawężenia do monitorów włączonych.
func (a *adapterBadan) WypiszMonitory(ctx context.Context,
	z shared.ResearchMonitorListRequest) (shared.ResearchMonitorListResponse, error) {

	if z.WindowId == "" {
		return shared.ResearchMonitorListResponse{}, bladWskazaniaBadan("monitor.list bez okna badania")
	}
	tylkoWlaczone := z.EnabledOnly != nil && *z.EnabledOnly
	monitory, err := a.repozytorium.Monitory(ctx, z.WindowId, tylkoWlaczone)
	if err != nil {
		return shared.ResearchMonitorListResponse{}, bladBadan(err)
	}
	przelozone := make([]shared.ResearchMonitor, 0, len(monitory))
	oczekujace := 0
	for _, monitor := range monitory {
		przelozone = append(przelozone, zlozMonitorBadania(monitor))
		oczekujace += monitor.Oczekujace
	}
	return shared.ResearchMonitorListResponse{Monitors: przelozone, PendingTotal: oczekujace}, nil
}

// OdswiezMonitory obsługuje `research.monitor.refresh`. Pozycje nowe wchodzą do
// bazy odkryć, więc skrzynka „nowe źródła" przeżywa zamknięcie okna.
func (a *adapterBadan) OdswiezMonitory(ctx context.Context,
	z shared.ResearchMonitorRefreshRequest) (shared.ResearchMonitorRefreshResponse, error) {

	if z.WindowId == "" {
		return shared.ResearchMonitorRefreshResponse{}, bladWskazaniaBadan("monitor.refresh bez okna badania")
	}
	monitory, err := a.repozytorium.Monitory(ctx, z.WindowId, false)
	if err != nil {
		return shared.ResearchMonitorRefreshResponse{}, bladBadan(err)
	}

	wyniki := []shared.ResearchDiscoveryResult{}
	zawiodly := []string{}
	odswiezone := []shared.ResearchMonitor{}
	teraz := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")

	for _, monitor := range monitory {
		if z.MonitorId != nil && *z.MonitorId != "" && monitor.Kod != *z.MonitorId {
			continue
		}
		if !monitor.Wlaczony {
			odswiezone = append(odswiezone, zlozMonitorBadania(monitor))
			continue
		}
		pozycje, err := a.pozycjeMonitoraBadania(ctx, monitor)
		if err != nil {
			zawiodly = append(zawiodly, monitor.Kod)
			odswiezone = append(odswiezone, zlozMonitorBadania(monitor))
			continue
		}
		zapamietane, err := a.zapamietajWynikiBadania(ctx, z.WindowId, pozycje)
		if err != nil {
			return shared.ResearchMonitorRefreshResponse{}, err
		}
		wyniki = append(wyniki, zapamietane...)

		monitor.Oczekujace = len(zapamietane)
		monitor.OdswiezonoO = &teraz
		zapisany, err := a.repozytorium.ZapiszMonitorZeStanem(ctx, monitor)
		if err != nil {
			return shared.ResearchMonitorRefreshResponse{}, bladBadan(err)
		}
		odswiezone = append(odswiezone, zlozMonitorBadania(zapisany))
	}
	return shared.ResearchMonitorRefreshResponse{
		Results: wyniki, Monitors: odswiezone, FailedMonitorIds: zawiodly,
	}, nil
}

// pozycjeMonitoraBadania pobiera nowe pozycje monitora wedle jego rodzaju,
// czytając kanał albo powtarzając zapytanie tematu.
func (a *adapterBadan) pozycjeMonitoraBadania(ctx context.Context,
	monitor dane.MonitorBadania) ([]shared.ResearchDiscoveryResult, error) {

	if monitor.Rodzaj == string(shared.ResearchMonitorKindFeed) {
		if monitor.Adres == nil {
			return nil, bladWskazaniaBadan("monitor kanału " + monitor.Kod + " nie ma adresu")
		}
		return czytajKanal(ctx, *monitor.Adres)
	}
	if monitor.Zapytanie == nil {
		return nil, bladWskazaniaBadan("monitor tematu " + monitor.Kod + " nie ma zapytania")
	}
	pozycje, _, err := szukajCrossref(ctx, *monitor.Zapytanie, nil, nil, limitOdkryciaBadania)
	return pozycje, err
}

// zlozMonitorBadania przekłada wiersz monitora z bazy danych na kształt
// odpowiedzi zgodny z kontraktem, jaki widzi klient.
func zlozMonitorBadania(m dane.MonitorBadania) shared.ResearchMonitor {
	monitor := shared.ResearchMonitor{
		Id: m.Kod, WindowId: m.Okno, Kind: shared.ResearchMonitorKind(m.Rodzaj),
		Query: m.Zapytanie, Url: m.Adres, Enabled: m.Wlaczony,
	}
	if m.InterwalMinut != nil {
		interwal := int(*m.InterwalMinut)
		monitor.IntervalMinutes = &interwal
	}
	oczekujace := m.Oczekujace
	monitor.PendingCount = &oczekujace
	if m.OdswiezonoO != nil {
		chwila := chwilaBazy(*m.OdswiezonoO)
		monitor.LastRefreshedAt = &chwila
	}
	return monitor
}

// ── Import wsadowy ─────────────────────────────────────────────────────────

// WczytajPartieAdresow obsługuje `research.batch.import` i pozyskuje wiele
// adresów naraz, znakując każde pozyskane źródło identyfikatorem partii.
func (a *adapterBadan) WczytajPartieAdresow(ctx context.Context,
	z shared.ResearchBatchImportRequest) (shared.ResearchBatchImportResponse, error) {

	if z.WindowId == "" {
		return shared.ResearchBatchImportResponse{}, bladWskazaniaBadan("batch.import bez okna badania")
	}
	if len(z.Urls) == 0 {
		return shared.ResearchBatchImportResponse{},
			bladWskazaniaBadan("batch.import bez adresów — nie ma czego pozyskać")
	}
	tryb := shared.ResearchCaptureMode(shared.ResearchCaptureModeReadable)
	if z.Mode != nil && *z.Mode != "" {
		tryb = *z.Mode
	}
	partia := nowyIdentyfikator(przedrostekPartiiBadania)

	przyjete := 0
	odrzucone := 0
	powody := []string{}
	for _, adres := range z.Urls {
		adres = strings.TrimSpace(adres)
		if adres == "" {
			odrzucone++
			powody = append(powody, "pusty wiersz na liście adresów")
			continue
		}
		odpowiedz, err := a.PrzechwycStrone(ctx, shared.ResearchSourceCaptureRequest{
			WindowId: z.WindowId, Url: adres, Mode: tryb,
		})
		if err != nil {
			odrzucone++
			powody = append(powody, adres+": "+err.Error())
			continue
		}
		zrodlo, err := a.repozytorium.Zrodlo(ctx, odpowiedz.Source.Id)
		if err == nil {
			pochodzenie := partia
			zrodlo.Pochodzenie = &pochodzenie
			if _, err := a.repozytorium.ZapiszZrodlo(ctx, zrodlo); err != nil {
				return shared.ResearchBatchImportResponse{}, bladBadan(err)
			}
		}
		przyjete++
	}
	return shared.ResearchBatchImportResponse{
		QueueItemId: partia, Accepted: przyjete, Rejected: odrzucone, RejectReasons: powody,
	}, nil
}
