// Pakiet obsługuje rodzinę komend `extension.*`: transport, poświadczenie,
// odkrywanie i wywołanie narzędzi MCP, dziennik protokołu, piaskownicę, import
// definicji API, webhooki, odwzorowania, metryki użycia, kondycję i audyt
// rozszerzeń.
package core

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// czasPiaskownicyRozszerzenia ogranicza czas trwania pojedynczego przebiegu
// próbnego umiejętności albo wtyczki uruchamianego komendą piaskownicy
// rozszerzeń.
const czasPiaskownicyRozszerzenia = 30 * time.Second

// UstawTransport obsługuje `extension.transport.set`. Pole `probe` włącza
// próbne powitanie: bez niego zapis jest zapisem nastawy, z nim — nastawą wraz
// ze sprawdzeniem, czy serwer pod nią stoi.
func (a *adapterRozszerzen) UstawTransport(ctx context.Context,
	z shared.ExtensionTransportSetRequest) (shared.ExtensionTransportSetResponse, error) {

	wiersz, err := a.pozycjaRozszerzeniaZadania(ctx, z.ExtensionId)
	if err != nil {
		return shared.ExtensionTransportSetResponse{}, err
	}
	if err := sprawdzTransportRozszerzenia(z.Transport); err != nil {
		return shared.ExtensionTransportSetResponse{}, err
	}
	transport := string(z.Transport)
	if _, err := a.rejestr.ZapiszIntegracjeRozszerzenia(ctx, dane.IntegracjaRozszerzenia{
		RozszerzenieKod: wiersz.Identyfikator, Transport: &transport,
		Adres: z.Endpoint, Polecenie: z.Command, Zaktualizowano: a.teraz(),
	}); err != nil {
		return shared.ExtensionTransportSetResponse{}, bladRozszerzenia(err)
	}
	odswiezona, err := a.rejestr.Rozszerzenie(ctx, wiersz.Identyfikator)
	if err != nil {
		return shared.ExtensionTransportSetResponse{}, bladRozszerzenia(err)
	}
	pozycja := rozszerzenieKontraktu(odswiezona)
	a.rozglosRozszerzenie(shared.ChangeKindUpdated, pozycja)
	a.odnotujCyklZycia(ctx, odswiezona, shared.ExtensionLifecycleActionConfigured, nil, nil,
		"transport ustawiony na "+transport)

	stan := shared.ExtensionHealthStatus(shared.ExtensionHealthStatusUnknown)
	if z.Probe != nil && *z.Probe {
		kondycja := a.zmierzKondycjeRozszerzenia(ctx, odswiezona)
		stan = kondycja.Status
	}
	return shared.ExtensionTransportSetResponse{Extension: pozycja, ProbeStatus: &stan}, nil
}

// PowiazPoswiadczenie obsługuje `extension.credential.bind`. Do rdzenia wchodzi
// WYŁĄCZNIE odwołanie — klucz jawny warstwy sekretów — i to ono ląduje w bazie.
// Treści poświadczenia rodzina `extension.*` nie widzi na żadnym kroku.
func (a *adapterRozszerzen) PowiazPoswiadczenie(ctx context.Context,
	z shared.ExtensionCredentialBindRequest) (shared.ExtensionCredentialBindResponse, error) {

	wiersz, err := a.pozycjaRozszerzeniaZadania(ctx, z.ExtensionId)
	if err != nil {
		return shared.ExtensionCredentialBindResponse{}, err
	}
	if err := sprawdzSposobLogowaniaRozszerzenia(z.AuthKind); err != nil {
		return shared.ExtensionCredentialBindResponse{}, err
	}
	odwolanie := strings.TrimSpace(z.CredentialRef)
	if odwolanie == "" {
		return shared.ExtensionCredentialBindResponse{}, bladWskazaniaRozszerzenia(
			"powiązanie bez odwołania do poświadczenia")
	}
	sposob := string(z.AuthKind)
	integracja, err := a.rejestr.ZapiszIntegracjeRozszerzenia(ctx, dane.IntegracjaRozszerzenia{
		RozszerzenieKod: wiersz.Identyfikator, SposobLogowania: &sposob,
		OdwolanieSekretu: &odwolanie, Zakresy: z.Scopes, Zaktualizowano: a.teraz(),
	})
	if err != nil {
		return shared.ExtensionCredentialBindResponse{}, bladRozszerzenia(err)
	}
	// Odwołanie wchodzi też do rejestru referencji sekretów rozszerzeń.
	if _, err := a.rejestr.ZapiszSekretRozszerzenia(ctx, dane.SekretRozszerzenia{
		Odwolanie: odwolanie, SposobLogowania: &sposob, Zaktualizowano: a.teraz(),
	}, false); err != nil {
		return shared.ExtensionCredentialBindResponse{}, bladRozszerzenia(err)
	}

	odswiezona, err := a.rejestr.Rozszerzenie(ctx, wiersz.Identyfikator)
	if err != nil {
		return shared.ExtensionCredentialBindResponse{}, bladRozszerzenia(err)
	}
	pozycja := rozszerzenieKontraktu(odswiezona)
	a.rozglosRozszerzenie(shared.ChangeKindUpdated, pozycja)
	a.odnotujCyklZycia(ctx, odswiezona, shared.ExtensionLifecycleActionConfigured, nil, nil,
		"poświadczenie powiązane sposobem "+sposob)

	odpowiedz := shared.ExtensionCredentialBindResponse{Extension: pozycja}
	// Adres zgody powstaje wyłącznie dla uwierzytelnienia OAuth2 z ustawionym
	// adresem integracji.
	if z.AuthKind == shared.ExtensionAuthKindOauth2 && integracja.Adres != nil && *integracja.Adres != "" {
		adres := strings.TrimSuffix(*integracja.Adres, "/") + "/authorize"
		if len(z.Scopes) > 0 {
			adres += "?scope=" + strings.Join(z.Scopes, "+")
		}
		odpowiedz.AuthorizationUrl = &adres
	}
	return odpowiedz, nil
}

// WypiszNarzedzia obsługuje `extension.tool.list`. Bez `refresh` oddaje to, co
// odkryto wcześniej; z `refresh` — pyta serwer na nowo i wymienia wykaz.
func (a *adapterRozszerzen) WypiszNarzedzia(ctx context.Context,
	z shared.ExtensionToolListRequest) (shared.ExtensionToolListResponse, error) {

	wiersz, err := a.pozycjaRozszerzeniaZadania(ctx, z.ExtensionId)
	if err != nil {
		return shared.ExtensionToolListResponse{}, err
	}
	rodzaj := ""
	if z.Kind != nil {
		if err := sprawdzRodzajNarzedziaRozszerzenia(*z.Kind); err != nil {
			return shared.ExtensionToolListResponse{}, err
		}
		rodzaj = string(*z.Kind)
	}

	wersjaProtokolu := wersjaProtokoluMcp
	if z.Refresh != nil && *z.Refresh {
		odkryta, err := a.odkryjNarzedziaRozszerzenia(ctx, wiersz)
		if err != nil {
			return shared.ExtensionToolListResponse{}, err
		}
		wersjaProtokolu = odkryta
	}

	wpisy, err := a.rejestr.NarzedziaRozszerzenia(ctx, wiersz.Identyfikator, rodzaj)
	if err != nil {
		return shared.ExtensionToolListResponse{}, bladRozszerzenia(err)
	}
	narzedzia := make([]shared.ExtensionToolEntry, 0, len(wpisy))
	odkryto := int64(0)
	for _, wpis := range wpisy {
		narzedzia = append(narzedzia, narzedzieKontraktuRozszerzen(wpis))
		if wpis.Odkryto > odkryto {
			odkryto = wpis.Odkryto
		}
		if wpis.WersjaProtokolu != nil && *wpis.WersjaProtokolu != "" {
			wersjaProtokolu = *wpis.WersjaProtokolu
		}
	}
	return shared.ExtensionToolListResponse{
		Entries: narzedzia, DiscoveredAt: odkryto,
		ProtocolVersion: wskaznikNapisuApp(wersjaProtokolu),
	}, nil
}

// odkryjNarzedziaRozszerzenia prowadzi powitanie i trzy wykazy protokołu,
// zapisuje wynik i oddaje wersję protokołu podaną przez serwer.
func (a *adapterRozszerzen) odkryjNarzedziaRozszerzenia(ctx context.Context,
	wiersz dane.Rozszerzenie) (string, error) {

	polaczenie, err := a.polaczenieRozszerzenia(ctx, wiersz)
	if err != nil {
		return "", err
	}
	rozmowa, err := otworzRozmoweMcp(ctx, polaczenie)
	if err != nil {
		return "", bladWskazaniaRozszerzenia(
			"nie można otworzyć rozmowy z serwerem pozycji " + wiersz.Kod + ": " + err.Error())
	}
	defer rozmowa.Zamknij()

	wersja, err := rozmowa.Powitaj(ctx)
	if err != nil {
		a.odlozRamkiProtokolu(ctx, wiersz.Identyfikator, rozmowa)
		return "", bladWskazaniaRozszerzenia(
			"powitanie serwera pozycji " + wiersz.Kod + " nie powiodło się: " + err.Error())
	}
	wpisy, err := rozmowa.OdkryjWykazy(ctx)
	a.odlozRamkiProtokolu(ctx, wiersz.Identyfikator, rozmowa)
	if err != nil {
		return "", bladWskazaniaRozszerzenia(
			"odkrywanie narzędzi pozycji " + wiersz.Kod + " nie powiodło się: " + err.Error())
	}

	teraz := a.teraz()
	wiersze := make([]dane.NarzedzieRozszerzenia, 0, len(wpisy))
	for _, wpis := range wpisy {
		wiersze = append(wiersze, dane.NarzedzieRozszerzenia{
			Rodzaj: string(wpis.Rodzaj), Nazwa: wpis.Nazwa,
			Opis: wskaznikNapisuApp(wpis.Opis), SchematWejscia: wskaznikNapisuApp(wpis.Schemat),
			Adres: wskaznikNapisuApp(wpis.Adres), Odkryto: teraz,
			WersjaProtokolu: wskaznikNapisuApp(wersja),
		})
	}
	if err := a.rejestr.ZapiszNarzedziaRozszerzenia(ctx, wiersz.Identyfikator, wiersze); err != nil {
		return "", bladRozszerzenia(err)
	}
	return wersja, nil
}

// WywolajNarzedzie obsługuje `extension.tool.call` — próbne wywołanie
// z inspektora. Niepowodzenie serwera NIE jest odmową komendy: inspektor ma
// pokazać, co odpowiedział serwer, wraz z czasem i powodem.
func (a *adapterRozszerzen) WywolajNarzedzie(ctx context.Context,
	z shared.ExtensionToolCallRequest) (shared.ExtensionToolCallResponse, error) {

	wiersz, err := a.pozycjaRozszerzeniaZadania(ctx, z.ExtensionId)
	if err != nil {
		return shared.ExtensionToolCallResponse{}, err
	}
	narzedzie := strings.TrimSpace(z.ToolName)
	if narzedzie == "" {
		return shared.ExtensionToolCallResponse{}, bladWskazaniaRozszerzenia(
			"wywołanie bez nazwy narzędzia")
	}
	polaczenie, err := a.polaczenieRozszerzenia(ctx, wiersz)
	if err != nil {
		return shared.ExtensionToolCallResponse{}, err
	}
	granica := czasWywolaniaMcp
	if z.TimeoutMs != nil && *z.TimeoutMs > 0 {
		granica = time.Duration(*z.TimeoutMs) * time.Millisecond
	}
	argumenty := json.RawMessage(`{}`)
	if len(z.Arguments) > 0 && json.Valid(z.Arguments) {
		argumenty = z.Arguments
	}
	parametry, err := json.Marshal(map[string]any{
		"name":      narzedzie,
		"arguments": argumenty,
	})
	if err != nil {
		return shared.ExtensionToolCallResponse{}, bladRozszerzenia(err)
	}

	poczatek := time.Now()
	rozmowa, err := otworzRozmoweMcp(ctx, polaczenie)
	if err != nil {
		return a.wynikWywolaniaRozszerzenia(ctx, wiersz, narzedzie, poczatek, nil,
			"nie można otworzyć rozmowy z serwerem: "+err.Error()), nil
	}
	defer rozmowa.Zamknij()

	if _, err := rozmowa.Powitaj(ctx); err != nil {
		a.odlozRamkiProtokolu(ctx, wiersz.Identyfikator, rozmowa)
		return a.wynikWywolaniaRozszerzenia(ctx, wiersz, narzedzie, poczatek, nil,
			"powitanie serwera nie powiodło się: "+err.Error()), nil
	}
	wynik, err := rozmowa.Wywolaj(ctx, "tools/call", parametry, granica)
	a.odlozRamkiProtokolu(ctx, wiersz.Identyfikator, rozmowa)
	if err != nil {
		return a.wynikWywolaniaRozszerzenia(ctx, wiersz, narzedzie, poczatek, nil,
			err.Error()), nil
	}
	return a.wynikWywolaniaRozszerzenia(ctx, wiersz, narzedzie, poczatek, wynik, ""), nil
}

// wynikWywolaniaRozszerzenia składa odpowiedź i odkłada wiersz użycia. Wiersz
// powstaje ZAWSZE — także przy niepowodzeniu, bo metryka błędów liczy się
// właśnie z tych wierszy.
func (a *adapterRozszerzen) wynikWywolaniaRozszerzenia(ctx context.Context, wiersz dane.Rozszerzenie,
	narzedzie string, poczatek time.Time, wynik json.RawMessage,
	powod string) shared.ExtensionToolCallResponse {

	czas := time.Since(poczatek).Milliseconds()
	_ = a.rejestr.DopiszWywolanieRozszerzenia(ctx, dane.WywolanieRozszerzenia{
		RozszerzenieKod: wiersz.Identyfikator, Narzedzie: &narzedzie,
		Udane: powod == "", CzasMs: czas, Szczegol: wskaznikNapisuApp(powod),
		Zaszlo: a.teraz(),
	})

	odpowiedz := shared.ExtensionToolCallResponse{Ok: powod == "", DurationMs: int(czas)}
	if powod != "" {
		odpowiedz.ErrorDetail = &powod
		return odpowiedz
	}
	odpowiedz.Raw = wynik
	odpowiedz.Text = wskaznikNapisuApp(tekstWynikuNarzedzia(wynik))
	return odpowiedz
}

// tekstWynikuNarzedzia wyciąga treść czytelną dla człowieka z odpowiedzi
// `tools/call`. Protokół niesie ją w wykazie `content` — pozycje rodzaju `text`
// składamy w jeden napis; brak takich pozycji daje pustkę, nie zmyśloną treść.
func tekstWynikuNarzedzia(wynik json.RawMessage) string {
	var odpowiedz struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.Unmarshal(wynik, &odpowiedz); err != nil {
		return ""
	}
	czesci := make([]string, 0, len(odpowiedz.Content))
	for _, wpis := range odpowiedz.Content {
		if wpis.Type != "text" || wpis.Text == "" {
			continue
		}
		czesci = append(czesci, wpis.Text)
	}
	return strings.Join(czesci, "\n")
}

// WypiszLogProtokolu obsługuje komendę `extension.protocol.log.list`,
// zwracając ramki dziennika protokołu rozmowy z serwerem pozycji wraz
// z licznikiem ogółem.
func (a *adapterRozszerzen) WypiszLogProtokolu(ctx context.Context,
	z shared.ExtensionProtocolLogListRequest) (shared.ExtensionProtocolLogListResponse, error) {

	wiersz, err := a.pozycjaRozszerzeniaZadania(ctx, z.ExtensionId)
	if err != nil {
		return shared.ExtensionProtocolLogListResponse{}, err
	}
	od := int64(0)
	if z.Since != nil {
		od = *z.Since
	}
	granica := 0
	if z.Limit != nil {
		granica = *z.Limit
	}
	wiersze, razem, err := a.rejestr.RamkiProtokolu(ctx, wiersz.Identyfikator, od, granica)
	if err != nil {
		return shared.ExtensionProtocolLogListResponse{}, bladRozszerzenia(err)
	}
	ramki := make([]shared.ProtocolFrame, 0, len(wiersze))
	for _, ramka := range wiersze {
		wpis := shared.ProtocolFrame{
			Id: ramka.Kod, Direction: shared.ProtocolFrameDirection(ramka.Kierunek),
			Method: ramka.Metoda, CorrelationId: ramka.Korelacja,
			ErrorCode: ramka.KodBledu, OccurredAt: ramka.Zaszlo,
		}
		if ramka.Tresc != nil && json.Valid([]byte(*ramka.Tresc)) {
			wpis.Payload = json.RawMessage(*ramka.Tresc)
		}
		ramki = append(ramki, wpis)
	}
	return shared.ExtensionProtocolLogListResponse{Frames: ramki, Total: razem}, nil
}

// UruchomPiaskownice obsługuje `extension.sandbox.run`: przebieg próbny
// umiejętności albo wtyczki na przykładowym wejściu, bez podłączania do
// eksperta produkcyjnego, transportem procesu lokalnego albo protokołem MCP.
func (a *adapterRozszerzen) UruchomPiaskownice(ctx context.Context,
	z shared.ExtensionSandboxRunRequest) (shared.ExtensionSandboxRunResponse, error) {

	wiersz, err := a.pozycjaRozszerzeniaZadania(ctx, z.ExtensionId)
	if err != nil {
		return shared.ExtensionSandboxRunResponse{}, err
	}
	if len(z.Input) == 0 || !json.Valid(z.Input) {
		return shared.ExtensionSandboxRunResponse{}, bladWskazaniaRozszerzenia(
			"piaskownica wymaga wejścia w postaci poprawnego JSON-a")
	}
	if a.magazyn == nil {
		return shared.ExtensionSandboxRunResponse{}, bladBrakuMagazynuApp(
			"przebieg piaskownicy")
	}
	integracja, err := a.rejestr.IntegracjaRozszerzenia(ctx, wiersz.Identyfikator)
	if err != nil && !isBrakWierszaApp(err) {
		return shared.ExtensionSandboxRunResponse{}, bladRozszerzenia(err)
	}
	granica := czasPiaskownicyRozszerzenia
	if z.TimeoutMs != nil && *z.TimeoutMs > 0 {
		granica = time.Duration(*z.TimeoutMs) * time.Millisecond
	}
	przebieg, zakoncz := context.WithTimeout(ctx, granica)
	defer zakoncz()

	poczatek := time.Now()
	var wyjscie []byte
	var powod string

	switch {
	case integracja.Polecenie != nil && strings.TrimSpace(*integracja.Polecenie) != "":
		polecenie := strings.Fields(*integracja.Polecenie)
		// #nosec G204 — polecenie pochodzi od integracji podłączonej rdzeniowi.
		proces := exec.CommandContext(przebieg, polecenie[0], polecenie[1:]...)
		proces.Stdin = strings.NewReader(string(z.Input))
		wynik, err := proces.CombinedOutput()
		wyjscie = wynik
		if err != nil {
			powod = "proces rozszerzenia zakończył się niepowodzeniem: " + err.Error()
		}

	case integracja.Adres != nil && strings.TrimSpace(*integracja.Adres) != "":
		zadanie, err := http.NewRequestWithContext(przebieg, http.MethodPost, *integracja.Adres,
			strings.NewReader(string(z.Input)))
		if err != nil {
			return shared.ExtensionSandboxRunResponse{}, bladRozszerzenia(err)
		}
		zadanie.Header.Set("Content-Type", "application/json")
		klient := &http.Client{Timeout: granica}
		odpowiedz, err := klient.Do(zadanie)
		if err != nil {
			powod = "serwer rozszerzenia nie odpowiedział: " + err.Error()
			break
		}
		wyjscie, err = io.ReadAll(io.LimitReader(odpowiedz.Body, granicaOdpowiedziMcp))
		_ = odpowiedz.Body.Close()
		if err != nil {
			powod = "nie można odczytać odpowiedzi rozszerzenia: " + err.Error()
		}
		if odpowiedz.StatusCode >= 400 {
			powod = "rozszerzenie odpowiedziało kodem " + strconv.Itoa(odpowiedz.StatusCode)
		}

	default:
		return shared.ExtensionSandboxRunResponse{}, bladWskazaniaRozszerzenia(
			"pozycja " + wiersz.Kod + " nie ma ani polecenia procesu, ani adresu — " +
				"piaskownica nie ma czego uruchomić (ustaw je komendą extension.transport.set)")
	}

	czas := time.Since(poczatek).Milliseconds()
	// Log przebiegu ląduje w magazynie treści rdzenia, do którego kontrakt
	// oddaje odwołanie.
	dziennik := "wejście:\n" + string(z.Input) + "\n\nwyjście:\n" + string(wyjscie)
	if powod != "" {
		dziennik += "\n\npowód niepowodzenia:\n" + powod
	}
	odwolanie, _, err := a.wniesLogPiaskownicy(dziennik)
	if err != nil {
		return shared.ExtensionSandboxRunResponse{}, bladRozszerzenia(err)
	}

	_ = a.rejestr.DopiszWywolanieRozszerzenia(ctx, dane.WywolanieRozszerzenia{
		RozszerzenieKod: wiersz.Identyfikator,
		Narzedzie:       wskaznikNapisuApp("piaskownica"),
		Udane:           powod == "", CzasMs: czas,
		Szczegol: wskaznikNapisuApp(powod), Zaszlo: a.teraz(),
	})

	odpowiedz := shared.ExtensionSandboxRunResponse{
		Ok: powod == "", LogRef: &odwolanie, DurationMs: int(czas),
	}
	// Wyjście wchodzi w pole `output` wyłącznie wtedy, gdy jest poprawnym
	// JSON-em.
	if json.Valid(wyjscie) {
		odpowiedz.Output = json.RawMessage(wyjscie)
	}
	return odpowiedz, nil
}

// wniesLogPiaskownicy kładzie dziennik przebiegu piaskownicy w magazynie
// treści rdzenia i oddaje odwołanie do zapisanego pliku wraz z jego rozmiarem.
func (a *adapterRozszerzen) wniesLogPiaskownicy(dziennik string) (string, int64, error) {
	bajty := []byte(dziennik)
	suma := sumaTresciApp(bajty)
	sciezka, err := a.magazyn.Zapisz(bajty, suma)
	if err != nil {
		return "", 0, err
	}
	return odwolanieWytworuApp(sciezka), int64(len(bajty)), nil
}

// ZaimportujDefinicje obsługuje `extension.definition.import`: rozkłada opis
// OpenAPI 3 albo schemat GraphQL wskazanej pozycji na wykaz operacji
// rozszerzenia.
func (a *adapterRozszerzen) ZaimportujDefinicje(ctx context.Context,
	z shared.ExtensionDefinitionImportRequest) (shared.ExtensionDefinitionImportResponse, error) {

	if err := a.sprawdzKatalog(); err != nil {
		return shared.ExtensionDefinitionImportResponse{}, err
	}
	kod := strings.TrimSpace(z.Code)
	if kod == "" {
		return shared.ExtensionDefinitionImportResponse{}, bladWskazaniaRozszerzenia(
			"import definicji bez kodu pozycji")
	}
	if err := sprawdzFormatDefinicjiRozszerzenia(z.Format); err != nil {
		return shared.ExtensionDefinitionImportResponse{}, err
	}

	tresc, zrodloOpisu, err := a.trescDefinicjiRozszerzenia(ctx, z)
	if err != nil {
		return shared.ExtensionDefinitionImportResponse{}, err
	}

	var operacje []shared.ExtensionToolEntry
	var nazwa string
	switch z.Format {
	case shared.ExtensionDefinitionFormatOpenapi3:
		operacje, nazwa, err = operacjeZOpenapi(tresc)
	case shared.ExtensionDefinitionFormatGraphql:
		operacje, nazwa, err = operacjeZGraphql(tresc)
	}
	if err != nil {
		return shared.ExtensionDefinitionImportResponse{}, bladWskazaniaRozszerzenia(
			"nie można odczytać definicji: " + err.Error())
	}
	if nazwa == "" {
		nazwa = kod
	}

	teraz := a.teraz()
	konfiguracja, err := json.Marshal(map[string]any{
		"definitionFormat": string(z.Format),
		"definitionSource": zrodloOpisu,
		"operationCount":   len(operacje),
	})
	if err != nil {
		return shared.ExtensionDefinitionImportResponse{}, bladRozszerzenia(err)
	}

	pozycjaWiersz, err := a.zalozAlboOdswiezPozycje(ctx, kod, nazwa, string(konfiguracja), teraz)
	if err != nil {
		return shared.ExtensionDefinitionImportResponse{}, err
	}

	// Operacje lądują w tym samym wykazie narzędzi pozycji, który wypełnia
	// odkrywanie MCP.
	wiersze := make([]dane.NarzedzieRozszerzenia, 0, len(operacje))
	for _, operacja := range operacje {
		wiersze = append(wiersze, dane.NarzedzieRozszerzenia{
			Rodzaj: string(operacja.Kind), Nazwa: operacja.Name,
			Opis: operacja.Description, Odkryto: teraz,
			SchematWejscia: wskaznikNapisuApp(string(operacja.InputSchema)),
		})
	}
	if err := a.rejestr.ZapiszNarzedziaRozszerzenia(ctx, pozycjaWiersz.Identyfikator, wiersze); err != nil {
		return shared.ExtensionDefinitionImportResponse{}, bladRozszerzenia(err)
	}
	odswiezona, err := a.rejestr.Rozszerzenie(ctx, pozycjaWiersz.Identyfikator)
	if err != nil {
		return shared.ExtensionDefinitionImportResponse{}, bladRozszerzenia(err)
	}
	pozycja := rozszerzenieKontraktu(odswiezona)
	a.rozglosRozszerzenie(shared.ChangeKindUpdated, pozycja)
	a.odnotujCyklZycia(ctx, odswiezona, shared.ExtensionLifecycleActionConfigured, nil, nil,
		"import definicji "+string(z.Format)+", operacji: "+strconv.Itoa(len(operacje)))

	return shared.ExtensionDefinitionImportResponse{Extension: pozycja, Operations: operacje}, nil
}

// trescDefinicjiRozszerzenia bierze treść opisu wprost z żądania albo pobiera
// ją spod adresu wskazanego w żądaniu importu definicji rozszerzenia.
func (a *adapterRozszerzen) trescDefinicjiRozszerzenia(ctx context.Context,
	z shared.ExtensionDefinitionImportRequest) ([]byte, string, error) {

	if z.Content != nil && strings.TrimSpace(*z.Content) != "" {
		return []byte(*z.Content), "treść w żądaniu", nil
	}
	adres := strings.TrimSpace(wartoscTekstu(z.Source))
	if adres == "" {
		return nil, "", bladWskazaniaRozszerzenia(
			"import definicji wymaga treści opisu albo adresu, spod którego go wziąć")
	}
	pobranie, zakoncz := context.WithTimeout(ctx, czasPowitaniaMcp)
	defer zakoncz()
	zadanie, err := http.NewRequestWithContext(pobranie, http.MethodGet, adres, nil)
	if err != nil {
		return nil, "", bladWskazaniaRozszerzenia("adres opisu jest niepoprawny: " + err.Error())
	}
	klient := &http.Client{Timeout: czasPowitaniaMcp}
	odpowiedz, err := klient.Do(zadanie)
	if err != nil {
		return nil, "", bladWskazaniaRozszerzenia("nie można pobrać opisu spod " + adres + ": " + err.Error())
	}
	defer odpowiedz.Body.Close()
	if odpowiedz.StatusCode >= 400 {
		return nil, "", bladWskazaniaRozszerzenia("adres opisu odpowiedział kodem " +
			strconv.Itoa(odpowiedz.StatusCode))
	}
	tresc, err := io.ReadAll(io.LimitReader(odpowiedz.Body, granicaOdpowiedziMcp))
	if err != nil {
		return nil, "", bladWskazaniaRozszerzenia("nie można odczytać opisu: " + err.Error())
	}
	return tresc, adres, nil
}

// zalozAlboOdswiezPozycje zakłada nową pozycję katalogu rozszerzeń dla
// podanego kodu albo odświeża nazwę i konfigurację pozycji już zastanej.
func (a *adapterRozszerzen) zalozAlboOdswiezPozycje(ctx context.Context,
	kod, nazwa, konfiguracja string, teraz int64) (dane.Rozszerzenie, error) {

	zastane, err := a.rejestr.RozszerzeniePoKodzie(ctx, kod)
	if err == nil {
		zmieniona, err := a.rejestr.ZmienRozszerzenie(ctx, zastane.Identyfikator,
			dane.ZmianaRozszerzenia{
				Nazwa: &nazwa, Konfiguracja: &konfiguracja,
				Zainstalowane: wartoscLogicznaApp(true),
			}, teraz)
		if err != nil {
			return dane.Rozszerzenie{}, bladRozszerzenia(err)
		}
		return zmieniona, nil
	}
	if !isBrakWierszaApp(err) {
		return dane.Rozszerzenie{}, bladRozszerzenia(err)
	}
	zalozona, err := a.rejestr.ZalozRozszerzenie(ctx, dane.Rozszerzenie{
		Identyfikator: nowyIdentyfikator(przedrostekRozszerzenia),
		Kod:           kod, Rodzaj: shared.ExtensionKindApi, Nazwa: nazwa,
		Zainstalowane: true, Wlaczone: false,
		ZrodloPochodzenia: shared.ExtensionOriginPersonal,
		Konfiguracja:      konfiguracja, Zaktualizowano: teraz,
	})
	if err != nil {
		return dane.Rozszerzenie{}, bladRozszerzenia(err)
	}
	return zalozona, nil
}

// operacjeZOpenapi rozkłada opis OpenAPI 3, wczytany biblioteką
// `getkin/kin-openapi`, na uporządkowany wykaz operacji dostępnych jako
// narzędzia.
func operacjeZOpenapi(tresc []byte) ([]shared.ExtensionToolEntry, string, error) {
	czytnik := openapi3.NewLoader()
	czytnik.IsExternalRefsAllowed = false
	opis, err := czytnik.LoadFromData(tresc)
	if err != nil {
		return nil, "", err
	}
	if opis.Paths == nil {
		return nil, "", fmt.Errorf("opis OpenAPI nie wymienia ani jednej ścieżki")
	}
	nazwa := ""
	if opis.Info != nil {
		nazwa = opis.Info.Title
	}

	operacje := []shared.ExtensionToolEntry{}
	for _, sciezka := range sort.StringSlice(kluczeSciezekOpenapi(opis)) {
		pozycja := opis.Paths.Find(sciezka)
		if pozycja == nil {
			continue
		}
		for metoda, operacja := range pozycja.Operations() {
			nazwaOperacji := operacja.OperationID
			if nazwaOperacji == "" {
				nazwaOperacji = strings.ToLower(metoda) + " " + sciezka
			}
			wpis := shared.ExtensionToolEntry{
				Name: nazwaOperacji, Kind: shared.ExtensionToolKind(shared.ExtensionToolKindTool),
				Description: wskaznikNapisuApp(operacja.Summary),
				Uri:         wskaznikNapisuApp(strings.ToUpper(metoda) + " " + sciezka),
			}
			if schemat := schematWejsciaOperacji(operacja); schemat != nil {
				wpis.InputSchema = schemat
			}
			operacje = append(operacje, wpis)
		}
	}
	sort.SliceStable(operacje, func(i, j int) bool { return operacje[i].Name < operacje[j].Name })
	return operacje, nazwa, nil
}

// kluczeSciezekOpenapi oddaje ścieżki opisu w porządku ustalonym — mapa Go
// przechodzi się losowo, a wykaz operacji ma być ten sam przy każdym imporcie.
func kluczeSciezekOpenapi(opis *openapi3.T) []string {
	klucze := make([]string, 0, len(opis.Paths.Map()))
	for sciezka := range opis.Paths.Map() {
		klucze = append(klucze, sciezka)
	}
	sort.Strings(klucze)
	return klucze
}

// schematWejsciaOperacji składa schemat argumentów z parametrów operacji.
// Nie jest to przepisanie całego opisu — to wykaz nazw i rodzajów, który
// wystarcza formularzowi argumentów inspektora.
func schematWejsciaOperacji(operacja *openapi3.Operation) json.RawMessage {
	if len(operacja.Parameters) == 0 {
		return nil
	}
	wlasciwosci := map[string]any{}
	wymagane := []string{}
	for _, odwolanie := range operacja.Parameters {
		if odwolanie == nil || odwolanie.Value == nil {
			continue
		}
		parametr := odwolanie.Value
		rodzaj := "string"
		if parametr.Schema != nil && parametr.Schema.Value != nil &&
			parametr.Schema.Value.Type != nil && len(*parametr.Schema.Value.Type) > 0 {
			rodzaj = (*parametr.Schema.Value.Type)[0]
		}
		wlasciwosci[parametr.Name] = map[string]any{
			"type":        rodzaj,
			"description": parametr.Description,
		}
		if parametr.Required {
			wymagane = append(wymagane, parametr.Name)
		}
	}
	if len(wlasciwosci) == 0 {
		return nil
	}
	sort.Strings(wymagane)
	schemat, err := json.Marshal(map[string]any{
		"type": "object", "properties": wlasciwosci, "required": wymagane,
	})
	if err != nil {
		return nil
	}
	return schemat
}

// wzorzecPolaGraphql wyłuskuje ze schematu GraphQL pojedyncze pola typów
// `Query` i `Mutation` wraz z ich opisem zwracanym po dwukropku.
var wzorzecPolaGraphql = regexp.MustCompile(`(?m)^\s*([A-Za-z_][A-Za-z0-9_]*)\s*(\([^)]*\))?\s*:\s*([^\n#]+)`)

// wzorzecTypuGraphql wyłuskuje ze schematu GraphQL bloki treści typów
// operacji `Query` i `Mutation`, ograniczone parą nawiasów klamrowych.
var wzorzecTypuGraphql = regexp.MustCompile(`(?s)\btype\s+(Query|Mutation)\s*\{(.*?)\n\}`)

// operacjeZGraphql rozkłada schemat GraphQL na wykaz operacji. Operacją jest
// pole typu `Query` albo `Mutation` — reszta schematu opisuje kształty danych,
// nie czynności, więc do wykazu narzędzi nie należy.
func operacjeZGraphql(tresc []byte) ([]shared.ExtensionToolEntry, string, error) {
	bloki := wzorzecTypuGraphql.FindAllStringSubmatch(string(tresc), -1)
	if len(bloki) == 0 {
		return nil, "", fmt.Errorf("schemat GraphQL nie ma typu Query ani Mutation")
	}
	operacje := []shared.ExtensionToolEntry{}
	for _, blok := range bloki {
		rodzajOperacji := blok[1]
		for _, pole := range wzorzecPolaGraphql.FindAllStringSubmatch(blok[2], -1) {
			nazwa := pole[1]
			if nazwa == "" {
				continue
			}
			operacje = append(operacje, shared.ExtensionToolEntry{
				Name: nazwa, Kind: shared.ExtensionToolKind(shared.ExtensionToolKindTool),
				Description: wskaznikNapisuApp(rodzajOperacji + " → " + strings.TrimSpace(pole[3])),
				Uri:         wskaznikNapisuApp(strings.ToLower(rodzajOperacji) + " " + nazwa),
			})
		}
	}
	sort.SliceStable(operacje, func(i, j int) bool { return operacje[i].Name < operacje[j].Name })
	return operacje, "", nil
}

// WypiszWebhooki obsługuje komendę `extension.webhook.list`, zwracając wykaz
// webhooków rozszerzeń filtrowany kodem pozycji i kierunkiem.
func (a *adapterRozszerzen) WypiszWebhooki(ctx context.Context,
	z shared.ExtensionWebhookListRequest) (shared.ExtensionWebhookListResponse, error) {

	if err := a.sprawdzKatalog(); err != nil {
		return shared.ExtensionWebhookListResponse{}, err
	}
	kod := strings.TrimSpace(wartoscTekstu(z.ExtensionId))
	if kod != "" {
		if _, err := a.pozycjaRozszerzeniaZadania(ctx, kod); err != nil {
			return shared.ExtensionWebhookListResponse{}, err
		}
	}
	kierunek := ""
	if z.Direction != nil {
		if err := sprawdzKierunekWebhooka(*z.Direction); err != nil {
			return shared.ExtensionWebhookListResponse{}, err
		}
		kierunek = string(*z.Direction)
	}
	wiersze, err := a.rejestr.WebhookiRozszerzen(ctx, kod, kierunek)
	if err != nil {
		return shared.ExtensionWebhookListResponse{}, bladRozszerzenia(err)
	}
	webhooki := make([]shared.ExtensionWebhook, 0, len(wiersze))
	for _, wiersz := range wiersze {
		webhooki = append(webhooki, webhookKontraktu(wiersz))
	}
	return shared.ExtensionWebhookListResponse{Webhooks: webhooki, Total: len(webhooki)}, nil
}

// ZapiszWebhook obsługuje komendę `extension.webhook.save`, zakładając nowy
// webhook rozszerzenia albo zapisując zmianę webhooka zastanego.
func (a *adapterRozszerzen) ZapiszWebhook(ctx context.Context,
	z shared.ExtensionWebhookSaveRequest) (shared.ExtensionWebhookSaveResponse, error) {

	wiersz, err := a.pozycjaRozszerzeniaZadania(ctx, z.ExtensionId)
	if err != nil {
		return shared.ExtensionWebhookSaveResponse{}, err
	}
	if err := sprawdzKierunekWebhooka(z.Direction); err != nil {
		return shared.ExtensionWebhookSaveResponse{}, err
	}
	// Webhook wychodzący bez adresu docelowego nie ma dokąd wypchnąć zdarzenia.
	if z.Direction == shared.ExtensionWebhookDirectionOutbound &&
		strings.TrimSpace(wartoscTekstu(z.Url)) == "" {
		return shared.ExtensionWebhookSaveResponse{}, bladWskazaniaRozszerzenia(
			"webhook wychodzący wymaga adresu docelowego")
	}

	kod := strings.TrimSpace(wartoscTekstu(z.WebhookId))
	if kod == "" {
		kod = nowyIdentyfikator(przedrostekWebhookaRozszerzen)
	} else if _, err := a.rejestr.WebhookRozszerzenia(ctx, kod); err != nil {
		return shared.ExtensionWebhookSaveResponse{},
			bladNieznanegoBytuRozszerzenia("webhook", kod, err)
	}

	czynny := false
	if z.Enabled != nil {
		czynny = *z.Enabled
	}
	// Adres nasłuchu webhooka przychodzącego składa rdzeń, nie integracja
	// zewnętrzna.
	var nasluch *string
	if z.Direction == shared.ExtensionWebhookDirectionInbound {
		nasluch = wskaznikNapisuApp("/webhook/" + wiersz.Kod + "/" + kod)
	}

	zapisany, err := a.rejestr.ZapiszWebhookRozszerzenia(ctx, dane.WebhookRozszerzenia{
		Kod: kod, RozszerzenieKod: wiersz.Identyfikator, Kierunek: string(z.Direction),
		Adres: z.Url, AdresNasluchu: nasluch, Zdarzenia: z.EventTypes,
		OdwolanieSekretu: z.SecretRef, Czynny: czynny, Zaktualizowano: a.teraz(),
	})
	if err != nil {
		return shared.ExtensionWebhookSaveResponse{}, bladRozszerzenia(err)
	}
	return shared.ExtensionWebhookSaveResponse{Webhook: webhookKontraktu(zapisany)}, nil
}

// ZapiszMapowanie obsługuje `extension.mapping.save`. Zastrzeżenia walidacji
// wracają obok odwzorowania i niczego nie wstrzymują — reguła niepełna zostaje
// zapisana, bo Operator ma prawo pracować nad nią dalej.
func (a *adapterRozszerzen) ZapiszMapowanie(ctx context.Context,
	z shared.ExtensionMappingSaveRequest) (shared.ExtensionMappingSaveResponse, error) {

	wiersz, err := a.pozycjaRozszerzeniaZadania(ctx, z.ExtensionId)
	if err != nil {
		return shared.ExtensionMappingSaveResponse{}, err
	}
	if strings.TrimSpace(z.Name) == "" {
		return shared.ExtensionMappingSaveResponse{}, bladWskazaniaRozszerzenia(
			"odwzorowanie bez nazwy")
	}
	if len(z.Rules) == 0 || !json.Valid(z.Rules) {
		return shared.ExtensionMappingSaveResponse{}, bladWskazaniaRozszerzenia(
			"reguły odwzorowania nie są poprawnym JSON-em")
	}
	kod := strings.TrimSpace(wartoscTekstu(z.MappingId))
	if kod == "" {
		kod = nowyIdentyfikator(przedrostekMapowaniaRozszerzen)
	} else if _, err := a.rejestr.MapowanieRozszerzenia(ctx, kod); err != nil {
		return shared.ExtensionMappingSaveResponse{},
			bladNieznanegoBytuRozszerzenia("odwzorowanie", kod, err)
	}

	zapisane, err := a.rejestr.ZapiszMapowanieRozszerzenia(ctx, dane.MapowanieRozszerzenia{
		Kod: kod, RozszerzenieKod: wiersz.Identyfikator, Nazwa: z.Name,
		Reguly: string(z.Rules), Zaktualizowano: a.teraz(),
	})
	if err != nil {
		return shared.ExtensionMappingSaveResponse{}, bladRozszerzenia(err)
	}
	return shared.ExtensionMappingSaveResponse{
		Mapping: shared.ExtensionMapping{
			Id: zapisane.Kod, ExtensionId: zapisane.RozszerzenieKod, Name: zapisane.Nazwa,
			Rules: json.RawMessage(zapisane.Reguly), UpdatedAt: zapisane.Zaktualizowano,
		},
		ValidationIssues: zastrzezeniaOdwzorowania(z.Rules),
	}, nil
}

// zastrzezeniaOdwzorowania sprawdza kształt reguł. Reguła odwzorowania wiąże
// pole systemu zewnętrznego z polem platformy, więc para bez jednej ze stron
// nie odwzorowuje niczego — i o tym mówi zastrzeżenie.
func zastrzezeniaOdwzorowania(reguly json.RawMessage) []string {
	var pola map[string]any
	if err := json.Unmarshal(reguly, &pola); err != nil {
		return []string{"reguły nie są obiektem — odwzorowanie wiąże pola parami"}
	}
	if len(pola) == 0 {
		return []string{"reguły nie wiążą ani jednego pola"}
	}
	zastrzezenia := []string{}
	klucze := make([]string, 0, len(pola))
	for klucz := range pola {
		klucze = append(klucze, klucz)
	}
	sort.Strings(klucze)
	for _, klucz := range klucze {
		wartosc, jest := pola[klucz].(string)
		if !jest || strings.TrimSpace(wartosc) == "" {
			zastrzezenia = append(zastrzezenia,
				"pole "+klucz+" nie wskazuje odpowiednika po drugiej stronie")
		}
	}
	return zastrzezenia
}

// PobierzUzycie obsługuje komendę `extension.usage.get`: metryki wywołań
// i niepowodzeń liczone z wierszy dziennika użycia rozszerzeń.
func (a *adapterRozszerzen) PobierzUzycie(ctx context.Context,
	z shared.ExtensionUsageGetRequest) (shared.ExtensionUsageGetResponse, error) {

	if err := a.sprawdzKatalog(); err != nil {
		return shared.ExtensionUsageGetResponse{}, err
	}
	kod := strings.TrimSpace(wartoscTekstu(z.ExtensionId))
	if kod != "" {
		if _, err := a.pozycjaRozszerzeniaZadania(ctx, kod); err != nil {
			return shared.ExtensionUsageGetResponse{}, err
		}
	}
	// Okno domyślne obejmuje ostatnią dobę.
	do := a.teraz()
	if z.Until != nil && *z.Until > 0 {
		do = *z.Until
	}
	od := do - 24*60*60*1000
	if z.Since != nil {
		od = *z.Since
	}

	metryki, err := a.rejestr.MetrykiUzyciaRozszerzen(ctx, kod, od, do)
	if err != nil {
		return shared.ExtensionUsageGetResponse{}, bladRozszerzenia(err)
	}
	uzycie := make([]shared.ExtensionUsage, 0, len(metryki))
	for _, metryka := range metryki {
		wpis := shared.ExtensionUsage{
			ExtensionId: metryka.RozszerzenieKod,
			Calls:       metryka.Wywolan, Failures: metryka.Niepowodzen,
		}
		if metryka.SredniCzasMs != nil {
			sredni := int(*metryka.SredniCzasMs)
			wpis.AvgLatencyMs = &sredni
		}
		uzycie = append(uzycie, wpis)
	}
	return shared.ExtensionUsageGetResponse{
		Usage: uzycie, WindowStart: od, WindowEnd: do,
	}, nil
}

// SprawdzKondycje obsługuje komendę `extension.health.check`: powitanie
// protokołu MCP wskazanej pozycji albo wszystkich pozycji włączonych.
func (a *adapterRozszerzen) SprawdzKondycje(ctx context.Context,
	z shared.ExtensionHealthCheckRequest) (shared.ExtensionHealthCheckResponse, error) {

	if err := a.sprawdzKatalog(); err != nil {
		return shared.ExtensionHealthCheckResponse{}, err
	}
	wiersze := []dane.Rozszerzenie{}
	kod := strings.TrimSpace(wartoscTekstu(z.ExtensionId))
	if kod != "" {
		wiersz, err := a.pozycjaRozszerzeniaZadania(ctx, kod)
		if err != nil {
			return shared.ExtensionHealthCheckResponse{}, err
		}
		wiersze = append(wiersze, wiersz)
	} else {
		// Bez wskazania kodu sprawdzeniu podlegają wyłącznie pozycje włączone.
		wszystkie, err := a.rejestr.Rozszerzenia(ctx, dane.FiltrRozszerzen{TylkoZainstalowane: true})
		if err != nil {
			return shared.ExtensionHealthCheckResponse{}, bladRozszerzenia(err)
		}
		for _, wiersz := range wszystkie {
			if wiersz.Wlaczone {
				wiersze = append(wiersze, wiersz)
			}
		}
	}

	wyniki := make([]shared.ExtensionHealth, 0, len(wiersze))
	for _, wiersz := range wiersze {
		wyniki = append(wyniki, a.zmierzKondycjeRozszerzenia(ctx, wiersz))
	}
	return shared.ExtensionHealthCheckResponse{Results: wyniki, CheckedAt: a.teraz()}, nil
}

// zmierzKondycjeRozszerzenia prowadzi powitanie i odkłada wynik w dzienniku
// kondycji. Pozycja bez transportu wraca stanem `unknown` — to jest prawda
// o niej, a nie awaria.
func (a *adapterRozszerzen) zmierzKondycjeRozszerzenia(ctx context.Context,
	wiersz dane.Rozszerzenie) shared.ExtensionHealth {

	teraz := a.teraz()
	wynik := shared.ExtensionHealth{
		ExtensionId: wiersz.Identyfikator,
		Status:      shared.ExtensionHealthStatus(shared.ExtensionHealthStatusUnknown),
		CheckedAt:   teraz,
	}

	polaczenie, err := a.polaczenieRozszerzenia(ctx, wiersz)
	if err != nil {
		wynik.HandshakeError = wskaznikNapisuApp("pozycja nie ma ustawionego transportu")
		a.odlozKondycjeRozszerzenia(ctx, wiersz.Identyfikator, wynik)
		return wynik
	}

	poczatek := time.Now()
	rozmowa, err := otworzRozmoweMcp(ctx, polaczenie)
	if err != nil {
		wynik.Status = shared.ExtensionHealthStatus(shared.ExtensionHealthStatusUnavailable)
		wynik.HandshakeError = wskaznikNapisuApp(err.Error())
		a.odlozKondycjeRozszerzenia(ctx, wiersz.Identyfikator, wynik)
		return wynik
	}
	defer rozmowa.Zamknij()

	if _, err := rozmowa.Powitaj(ctx); err != nil {
		a.odlozRamkiProtokolu(ctx, wiersz.Identyfikator, rozmowa)
		wynik.Status = shared.ExtensionHealthStatus(shared.ExtensionHealthStatusUnavailable)
		wynik.HandshakeError = wskaznikNapisuApp(err.Error())
		a.odlozKondycjeRozszerzenia(ctx, wiersz.Identyfikator, wynik)
		return wynik
	}
	wpisy, err := rozmowa.OdkryjWykazy(ctx)
	a.odlozRamkiProtokolu(ctx, wiersz.Identyfikator, rozmowa)
	czas := int(time.Since(poczatek).Milliseconds())
	wynik.LatencyMs = &czas

	if err != nil {
		// Serwer odpowiada, lecz wykazu nie oddał — kontrakt ma dla tego stanu
		// osobną wartość.
		wynik.Status = shared.ExtensionHealthStatus(shared.ExtensionHealthStatusDegraded)
		wynik.HandshakeError = wskaznikNapisuApp(err.Error())
		a.odlozKondycjeRozszerzenia(ctx, wiersz.Identyfikator, wynik)
		return wynik
	}
	liczba := len(wpisy)
	wynik.ToolCount = &liczba
	wynik.Status = shared.ExtensionHealthStatus(shared.ExtensionHealthStatusHealthy)
	a.odlozKondycjeRozszerzenia(ctx, wiersz.Identyfikator, wynik)
	return wynik
}

// odlozKondycjeRozszerzenia zapisuje wynik sprawdzenia. Nieudany zapis nie
// przewraca sprawdzenia — wynik i tak jedzie w odpowiedzi.
func (a *adapterRozszerzen) odlozKondycjeRozszerzenia(ctx context.Context, kod string,
	wynik shared.ExtensionHealth) {

	wiersz := dane.KondycjaRozszerzenia{
		RozszerzenieKod: kod, Stan: string(wynik.Status),
		BladPowitania: wynik.HandshakeError, Sprawdzono: wynik.CheckedAt,
	}
	if wynik.LatencyMs != nil {
		czas := int64(*wynik.LatencyMs)
		wiersz.CzasMs = &czas
	}
	if wynik.ToolCount != nil {
		liczba := int64(*wynik.ToolCount)
		wiersz.LiczbaNarzedzi = &liczba
	}
	_ = a.rejestr.DopiszKondycjeRozszerzenia(ctx, wiersz)
}

// odlozRamkiProtokolu zapisuje ramki odbytej rozmowy protokołu MCP
// w dzienniku protokołu rozszerzenia i czyści bufor ramek rozmowy.
func (a *adapterRozszerzen) odlozRamkiProtokolu(ctx context.Context, kod string, rozmowa *rozmowaMcp) {
	if rozmowa == nil {
		return
	}
	for _, ramka := range rozmowa.ramki {
		_ = a.rejestr.DopiszRamkeProtokolu(ctx, dane.RamkaProtokolu{
			Kod: nowyIdentyfikator(przedrostekRamkiProtokolu), RozszerzenieKod: kod,
			Kierunek: string(ramka.Kierunek), Metoda: wskaznikNapisuApp(ramka.Metoda),
			Korelacja: wskaznikNapisuApp(ramka.Korelacja), Tresc: wskaznikNapisuApp(ramka.Tresc),
			KodBledu: wskaznikNapisuApp(ramka.KodBledu), Zaszlo: ramka.Zaszlo,
		})
	}
	rozmowa.ramki = nil
}

// WypiszAudyt obsługuje `extension.audit.list` — użycie widziane od strony
// eksperta, wraz z uprawnieniami obowiązującymi w chwili wywołania.
func (a *adapterRozszerzen) WypiszAudyt(ctx context.Context,
	z shared.ExtensionAuditListRequest) (shared.ExtensionAuditListResponse, error) {

	if err := a.sprawdzKatalog(); err != nil {
		return shared.ExtensionAuditListResponse{}, err
	}
	kod := strings.TrimSpace(wartoscTekstu(z.ExtensionId))
	if kod != "" {
		if _, err := a.pozycjaRozszerzeniaZadania(ctx, kod); err != nil {
			return shared.ExtensionAuditListResponse{}, err
		}
	}
	if err := a.sprawdzEksperta(ctx, wartoscTekstu(z.AgentId)); err != nil {
		return shared.ExtensionAuditListResponse{}, err
	}
	od := int64(0)
	if z.Since != nil {
		od = *z.Since
	}
	granica := 0
	if z.Limit != nil {
		granica = *z.Limit
	}

	wiersze, razem, err := a.rejestr.AudytRozszerzen(ctx, kod,
		strings.TrimSpace(wartoscTekstu(z.AgentId)), od, granica)
	if err != nil {
		return shared.ExtensionAuditListResponse{}, bladRozszerzenia(err)
	}

	// Uprawnienia czytamy raz na pozycję rozszerzenia, nie raz na wpis audytu.
	uprawnieniaPozycji := map[string][]shared.ExtensionPermission{}
	wpisy := make([]shared.ExtensionAuditEntry, 0, len(wiersze))
	for _, wiersz := range wiersze {
		nadane, jest := uprawnieniaPozycji[wiersz.RozszerzenieKod]
		if !jest {
			wszystkie, err := a.rejestr.UprawnieniaRozszerzenia(ctx, wiersz.RozszerzenieKod)
			if err != nil {
				return shared.ExtensionAuditListResponse{}, bladRozszerzenia(err)
			}
			for _, uprawnienie := range wszystkie {
				if !uprawnienie.Nadane {
					continue
				}
				nadane = append(nadane, uprawnienieKontraktuRozszerzen(uprawnienie))
			}
			uprawnieniaPozycji[wiersz.RozszerzenieKod] = nadane
		}
		wpisy = append(wpisy, shared.ExtensionAuditEntry{
			ExtensionId: wiersz.RozszerzenieKod, AgentId: wiersz.AgentKod,
			ToolName: wiersz.Narzedzie, Permissions: nadane, UsedAt: wiersz.Zaszlo,
		})
	}
	return shared.ExtensionAuditListResponse{Entries: wpisy, Total: razem}, nil
}

// polaczenieRozszerzenia składa opis połączenia z serwerem pozycji na
// podstawie zapisanej integracji: transportu, adresu i polecenia.
func (a *adapterRozszerzen) polaczenieRozszerzenia(ctx context.Context,
	wiersz dane.Rozszerzenie) (polaczenieMcp, error) {

	integracja, err := a.rejestr.IntegracjaRozszerzenia(ctx, wiersz.Identyfikator)
	if err != nil {
		if isBrakWierszaApp(err) {
			return polaczenieMcp{}, bladWskazaniaRozszerzenia(
				"pozycja " + wiersz.Kod + " nie ma ustawionego transportu — " +
					"ustaw go komendą extension.transport.set")
		}
		return polaczenieMcp{}, bladRozszerzenia(err)
	}
	if integracja.Transport == nil || *integracja.Transport == "" {
		return polaczenieMcp{}, bladWskazaniaRozszerzenia(
			"pozycja " + wiersz.Kod + " nie ma ustawionego transportu")
	}
	return polaczenieMcp{
		Transport: shared.McpTransport(*integracja.Transport),
		Adres:     wartoscTekstu(integracja.Adres),
		Polecenie: wartoscTekstu(integracja.Polecenie),
	}, nil
}

// webhookKontraktu przekłada wiersz webhooka rozszerzenia z bazy danych na
// kształt webhooka zwracany kontraktem komunikacji.
func webhookKontraktu(wiersz dane.WebhookRozszerzenia) shared.ExtensionWebhook {
	return shared.ExtensionWebhook{
		Id: wiersz.Kod, ExtensionId: wiersz.RozszerzenieKod,
		Direction: shared.ExtensionWebhookDirection(wiersz.Kierunek),
		Url:       wiersz.Adres, ReceiveUrl: wiersz.AdresNasluchu,
		EventTypes: wiersz.Zdarzenia, SecretRef: wiersz.OdwolanieSekretu,
		Enabled: wiersz.Czynny, UpdatedAt: wiersz.Zaktualizowano,
	}
}

// sprawdzTransportRozszerzenia dopuszcza wyłącznie wartości transportu
// wymienione w kontrakcie: stdio, sse, http, streamableHttp.
func sprawdzTransportRozszerzenia(transport shared.McpTransport) error {
	switch transport {
	case shared.McpTransportStdio, shared.McpTransportSse,
		shared.McpTransportHttp, shared.McpTransportStreamableHttp:
		return nil
	}
	return bladWskazaniaRozszerzenia("nieznany transport " + strconv.Quote(string(transport)) +
		" — dopuszczalne: stdio, sse, http, streamableHttp")
}

// sprawdzSposobLogowaniaRozszerzenia dopuszcza wyłącznie sposoby logowania
// integracji wymienione w kontrakcie komunikacji.
func sprawdzSposobLogowaniaRozszerzenia(sposob shared.ExtensionAuthKind) error {
	switch sposob {
	case shared.ExtensionAuthKindOauth2, shared.ExtensionAuthKindApiKey,
		shared.ExtensionAuthKindToken, shared.ExtensionAuthKindBasic,
		shared.ExtensionAuthKindNone:
		return nil
	}
	return bladWskazaniaRozszerzenia("nieznany sposób uwierzytelnienia " +
		strconv.Quote(string(sposob)) + " — dopuszczalne: oauth2, apiKey, token, basic, none")
}

// sprawdzRodzajNarzedziaRozszerzenia dopuszcza wyłącznie rodzaje wpisów
// wykazu narzędzi wymienione w kontrakcie: narzędzie, zasób, podpowiedź.
func sprawdzRodzajNarzedziaRozszerzenia(rodzaj shared.ExtensionToolKind) error {
	switch rodzaj {
	case shared.ExtensionToolKindTool, shared.ExtensionToolKindResource,
		shared.ExtensionToolKindPrompt:
		return nil
	}
	return bladWskazaniaRozszerzenia("nieznany rodzaj wpisu " + strconv.Quote(string(rodzaj)) +
		" — dopuszczalne: tool, resource, prompt")
}

// sprawdzKierunekWebhooka dopuszcza wyłącznie kierunki webhooka wymienione
// w kontrakcie: przychodzący albo wychodzący.
func sprawdzKierunekWebhooka(kierunek shared.ExtensionWebhookDirection) error {
	switch kierunek {
	case shared.ExtensionWebhookDirectionInbound, shared.ExtensionWebhookDirectionOutbound:
		return nil
	}
	return bladWskazaniaRozszerzenia("nieznany kierunek webhooka " +
		strconv.Quote(string(kierunek)) + " — dopuszczalne: inbound, outbound")
}

// sprawdzFormatDefinicjiRozszerzenia dopuszcza wyłącznie formaty definicji
// wymienione w kontrakcie: OpenAPI 3 albo GraphQL.
func sprawdzFormatDefinicjiRozszerzenia(format shared.ExtensionDefinitionFormat) error {
	switch format {
	case shared.ExtensionDefinitionFormatOpenapi3, shared.ExtensionDefinitionFormatGraphql:
		return nil
	}
	return bladWskazaniaRozszerzenia("nieznany format definicji " +
		strconv.Quote(string(format)) + " — dopuszczalne: openapi3, graphql")
}
