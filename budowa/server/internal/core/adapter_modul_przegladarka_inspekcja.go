// Moduł narzędzi inspekcyjnych obsługuje pracę na uruchomionej stronie
// komendami `browser.dom.inspect`, `browser.network.har`,
// `browser.console.read`, `browser.device.emulate` i `browser.scroll`, na
// adresie wziętym z ostatniej migawki okna.
package core

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// ZbadajDrzewo obsługuje `browser.dom.inspect` i oddaje drzewo elementów
// strony po zbudowaniu przez skrypty, z wybraną głębokością i stylami.
func (a *adapterPrzegladarki) ZbadajDrzewo(ctx context.Context,
	z shared.BrowserDomInspectRequest) (shared.BrowserDomInspectResponse, error) {

	migawka, err := a.adresOstatniejStrony(ctx, z.WindowId, "dom.inspect")
	if err != nil {
		return shared.BrowserDomInspectResponse{}, err
	}
	sesja, _, err := a.otworzStrone(ctx, "browser.dom.inspect", nastawyStrony{Url: migawka.Url})
	if err != nil {
		return shared.BrowserDomInspectResponse{}, err
	}
	defer sesja.Zamknij()

	glebokosc := wartoscLiczby(z.Depth)
	if glebokosc <= 0 {
		glebokosc = 6
	}
	korzen := wartoscTekstuLubPusta(z.Selector)
	if korzen == "" {
		korzen = "html"
	}
	zeStylami := z.IncludeStyles != nil && *z.IncludeStyles

	wyrazenie := `(() => {
		const korzen = document.querySelector(` + jakoLiteral(korzen) + `);
		if (!korzen) return null;
		const glebokoscMax = ` + jakoLiczba(glebokosc) + `;
		const zeStylami = ` + jakoLogiczna(zeStylami) + `;
		const wynik = [];
		const sciezka = (el) => {
			if (el.id) return '#' + el.id;
			const nazwa = el.tagName.toLowerCase();
			if (!el.parentElement) return nazwa;
			const rodzenstwo = Array.from(el.parentElement.children).filter(x => x.tagName === el.tagName);
			return nazwa + (rodzenstwo.length > 1 ? ':nth-of-type(' + (rodzenstwo.indexOf(el) + 1) + ')' : '');
		};
		const obejdz = (el, rodzic, poziom, numer) => {
			const wlasny = 'w' + numer.n++;
			const atrybuty = {};
			for (const a of el.attributes) atrybuty[a.name] = a.value;
			const wpis = {
				nodeId: wlasny, parentNodeId: rodzic, tagName: el.tagName.toLowerCase(),
				selector: sciezka(el), depth: poziom, attributes: atrybuty,
				text: (el.childNodes.length && el.textContent ? el.textContent.trim().slice(0, 400) : '')
			};
			if (zeStylami) {
				const s = window.getComputedStyle(el);
				wpis.computedStyles = {
					display: s.display, position: s.position, color: s.color,
					backgroundColor: s.backgroundColor, fontSize: s.fontSize, visibility: s.visibility
				};
			}
			wynik.push(wpis);
			if (poziom >= glebokoscMax) return;
			for (const dziecko of el.children) obejdz(dziecko, wlasny, poziom + 1, numer);
		};
		obejdz(korzen, null, 0, { n: 1 });
		return wynik;
	})()`

	wartosc, err := sesja.ocenNaStronie(ctx, wyrazenie)
	if err != nil {
		return shared.BrowserDomInspectResponse{}, bladSilnikaPrzegladarki("browser.dom.inspect", err)
	}
	var odczyt []struct {
		NodeId     string          `json:"nodeId"`
		ParentId   *string         `json:"parentNodeId"`
		TagName    string          `json:"tagName"`
		Selector   string          `json:"selector"`
		Text       string          `json:"text"`
		Attributes json.RawMessage `json:"attributes"`
		Styles     json.RawMessage `json:"computedStyles"`
		Depth      int             `json:"depth"`
	}
	if err := json.Unmarshal(wartosc, &odczyt); err != nil || len(odczyt) == 0 {
		return shared.BrowserDomInspectResponse{}, bladWskazaniaPrzegladarki(
			"na stronie " + migawka.Url + " nie ma elementu o wskazaniu " + korzen)
	}

	wezly := make([]shared.BrowserDomNode, 0, len(odczyt))
	for _, wpis := range odczyt {
		wezel := shared.BrowserDomNode{
			NodeId:       wpis.NodeId,
			ParentNodeId: wpis.ParentId,
			TagName:      wpis.TagName,
			Depth:        wpis.Depth,
		}
		if wpis.Selector != "" {
			selektor := wpis.Selector
			wezel.Selector = &selektor
		}
		if wpis.Text != "" {
			tekst := wpis.Text
			wezel.Text = &tekst
		}
		if len(wpis.Attributes) > 0 {
			wezel.Attributes = json.RawMessage(wpis.Attributes)
		}
		if len(wpis.Styles) > 0 {
			wezel.ComputedStyles = json.RawMessage(wpis.Styles)
		}
		wezly = append(wezly, wezel)
	}
	return shared.BrowserDomInspectResponse{Nodes: wezly, CapturedAt: time.Now().UnixMilli()}, nil
}

// RejestrSieciowy obsługuje `browser.network.har` i zapisuje rejestr żądań
// sieciowych strony w postaci pliku HAR jako wytwór sesji w magazynie modułu.
func (a *adapterPrzegladarki) RejestrSieciowy(ctx context.Context,
	z shared.BrowserNetworkHarRequest) (shared.BrowserNetworkHarResponse, error) {

	migawka, err := a.adresOstatniejStrony(ctx, z.WindowId, "network.har")
	if err != nil {
		return shared.BrowserNetworkHarResponse{}, err
	}
	sesja, _, err := a.otworzStrone(ctx, "browser.network.har", nastawyStrony{Url: migawka.Url})
	if err != nil {
		return shared.BrowserNetworkHarResponse{}, err
	}
	defer sesja.Zamknij()

	wpisy := wpisySieciowe(sesja.Zdarzenia(), wartoscLiczby(z.Limit))
	if len(wpisy) == 0 {
		return shared.BrowserNetworkHarResponse{}, bladSilnikaPrzegladarki("browser.network.har",
			errIntoKanal("przeglądarka nie odnotowała ani jednego żądania sieciowego dla "+migawka.Url))
	}

	har, err := json.Marshal(map[string]any{
		"log": map[string]any{
			"version": "1.2",
			"creator": map[string]any{"name": "Danaco Console", "version": shared.ProtocolVersion},
			"entries": wpisyHar(wpisy),
		},
	})
	if err != nil {
		return shared.BrowserNetworkHarResponse{}, bladPrzegladarki(err)
	}
	odwolanie, err := a.zapiszTresc(har)
	if err != nil {
		return shared.BrowserNetworkHarResponse{}, err
	}
	rozmiar := int64(len(har))
	tytul := "Rejestr sieciowy — " + migawka.Url
	mime := "application/json"
	if _, err := a.repozytorium.ZapiszWytwor(ctx, dane.WytworPrzegladania{
		Kod:                 nowyIdentyfikator(przedrostekWytworu),
		Okno:                z.WindowId,
		Rodzaj:              string(shared.BrowserArtifactKindNetworkLog),
		Tytul:               &tytul,
		TrescOdwolanie:      odwolanie,
		TypMime:             &mime,
		RozmiarBajtow:       &rozmiar,
		UrlZrodla:           &migawka.Url,
		MigawkaZewnetrznaID: &migawka.Kod,
	}); err != nil {
		return shared.BrowserNetworkHarResponse{}, bladPrzegladarki(err)
	}
	return shared.BrowserNetworkHarResponse{Entries: wpisy, HarRef: &odwolanie}, nil
}

// OdczytajKonsole obsługuje `browser.console.read` i oddaje komunikaty
// konsoli oraz błędy strony, zawężone do wybranych poziomów i granicy wpisów.
func (a *adapterPrzegladarki) OdczytajKonsole(ctx context.Context,
	z shared.BrowserConsoleReadRequest) (shared.BrowserConsoleReadResponse, error) {

	migawka, err := a.adresOstatniejStrony(ctx, z.WindowId, "console.read")
	if err != nil {
		return shared.BrowserConsoleReadResponse{}, err
	}
	sesja, _, err := a.otworzStrone(ctx, "browser.console.read", nastawyStrony{Url: migawka.Url})
	if err != nil {
		return shared.BrowserConsoleReadResponse{}, err
	}
	defer sesja.Zamknij()

	wpisy := wpisyKonsoli(sesja.Zdarzenia(), z.Levels, wartoscLiczby(z.Limit))
	return shared.BrowserConsoleReadResponse{Entries: wpisy}, nil
}

// EmulujUrzadzenie obsługuje `browser.device.emulate` — strona wczytywana jest
// w wymiarach i pod nazwą urządzenia, a odpowiedź niesie ZMIERZONE metryki wraz
// z migawką strony tak wyrenderowanej.
func (a *adapterPrzegladarki) EmulujUrzadzenie(ctx context.Context,
	z shared.BrowserDeviceEmulateRequest) (shared.BrowserDeviceEmulateResponse, error) {

	migawka, err := a.adresOstatniejStrony(ctx, z.WindowId, "device.emulate")
	if err != nil {
		return shared.BrowserDeviceEmulateResponse{}, err
	}
	metryki := metrykiUrzadzenia(z)
	nastawy := nastawyStrony{
		Url:          migawka.Url,
		Szerokosc:    metryki.Width,
		Wysokosc:     metryki.Height,
		SkalaPikseli: metryki.DeviceScaleFactor,
		Mobilne:      metryki.Mobile,
		Orientacja:   string(metryki.Orientation),
	}
	if metryki.UserAgent != nil {
		nastawy.AgentUzytkow = *metryki.UserAgent
	}
	sesja, wynik, err := a.otworzStrone(ctx, "browser.device.emulate", nastawy)
	if err != nil {
		return shared.BrowserDeviceEmulateResponse{}, err
	}
	defer sesja.Zamknij()

	zapisana, err := a.zapiszMigawkeZeStrony(ctx, z.WindowId, wynik)
	if err != nil {
		return shared.BrowserDeviceEmulateResponse{}, err
	}
	oddana := migawkaKontraktu(zapisana, false, false)
	return shared.BrowserDeviceEmulateResponse{Metrics: metryki, Snapshot: &oddana}, nil
}

// Przewin obsługuje `browser.scroll`: przewija stronę po stronie rdzenia
// i zapisuje migawkę stanu strony dopiero po przewinięciu.
func (a *adapterPrzegladarki) Przewin(ctx context.Context,
	z shared.BrowserScrollRequest) (shared.BrowserScrollResponse, error) {

	migawka, err := a.adresOstatniejStrony(ctx, z.WindowId, "scroll")
	if err != nil {
		return shared.BrowserScrollResponse{}, err
	}
	sesja, _, err := a.otworzStrone(ctx, "browser.scroll", nastawyStrony{Url: migawka.Url})
	if err != nil {
		return shared.BrowserScrollResponse{}, err
	}
	defer sesja.Zamknij()

	var polecenie string
	switch {
	case z.ToEnd != nil && *z.ToEnd:
		polecenie = `window.scrollTo(0, document.body.scrollHeight);`
	case z.ToSelector != nil && strings.TrimSpace(*z.ToSelector) != "":
		polecenie = `(() => { const e = document.querySelector(` + jakoLiteral(strings.TrimSpace(*z.ToSelector)) +
			`); if (!e) throw new Error('brak elementu'); e.scrollIntoView(); })();`
	default:
		polecenie = `window.scrollBy(` + jakoLiczba(wartoscLiczby(z.DeltaX)) + `, ` +
			jakoLiczba(wartoscLiczby(z.DeltaY)) + `);`
	}

	if _, err := sesja.ocenNaStronie(ctx, `(() => { `+polecenie+` return true; })()`); err != nil {
		return shared.BrowserScrollResponse{}, bladSilnikaPrzegladarki("browser.scroll", err)
	}
	// Strona dogrywająca treść przy przewijaniu potrzebuje chwili na dociągnięcie.
	time.Sleep(400 * time.Millisecond)

	stan, err := sesja.ocenNaStronie(ctx, `({
		url: location.href,
		title: document.title || '',
		text: (document.body && document.body.innerText) || '',
		html: document.documentElement ? document.documentElement.outerHTML : ''
	})`)
	if err != nil {
		return shared.BrowserScrollResponse{}, bladSilnikaPrzegladarki("browser.scroll", err)
	}
	var odczyt wynikOtwarcia
	var pomocniczy struct {
		Url   string `json:"url"`
		Tytul string `json:"title"`
		Tekst string `json:"text"`
		Html  string `json:"html"`
	}
	if err := json.Unmarshal(stan, &pomocniczy); err != nil {
		return shared.BrowserScrollResponse{}, bladPrzegladarki(err)
	}
	odczyt.Url, odczyt.Tytul, odczyt.Tekst, odczyt.Html =
		pomocniczy.Url, pomocniczy.Tytul, pomocniczy.Tekst, pomocniczy.Html

	zapisana, err := a.zapiszMigawkeZeStrony(ctx, z.WindowId, odczyt)
	if err != nil {
		return shared.BrowserScrollResponse{}, err
	}
	return shared.BrowserScrollResponse{Snapshot: migawkaKontraktu(zapisana, false, false)}, nil
}

// zapiszMigawkeZeStrony odkłada wynik wizyty na stronie jako kolejny wiersz
// historii migawek okna w bazie danych.
func (a *adapterPrzegladarki) zapiszMigawkeZeStrony(ctx context.Context, okno string,
	wynik wynikOtwarcia) (dane.MigawkaStrony, error) {

	zapisana, err := a.repozytorium.ZapiszMigawke(ctx, dane.MigawkaStrony{
		Kod: nowyIdentyfikator(przedrostekMigawki), Okno: okno, Url: wynik.Url,
		Tytul:           wskaznikTekstu(wynik.Tytul),
		TekstOdwolanie:  wskaznikTekstu(wynik.Tekst),
		ZrodloOdwolanie: wskaznikTekstu(wynik.Html),
	})
	if err != nil {
		return dane.MigawkaStrony{}, bladPrzegladarki(err)
	}
	return zapisana, nil
}

// metrykiUrzadzenia rozstrzyga wymiary emulacji. Nastawa nazwana daje wymiary
// typowe, `custom` bierze je z żądania, a `reset` wraca do biurka — bo zdjęcie
// emulacji też jest emulacją o znanych metrykach, nie brakiem odpowiedzi.
func metrykiUrzadzenia(z shared.BrowserDeviceEmulateRequest) shared.BrowserDeviceMetrics {
	nastawa := shared.BrowserDevicePreset(shared.BrowserDevicePresetDesktop)
	if z.Preset != nil {
		nastawa = *z.Preset
	}
	if z.Reset != nil && *z.Reset {
		nastawa = shared.BrowserDevicePreset(shared.BrowserDevicePresetDesktop)
	}
	metryki := shared.BrowserDeviceMetrics{
		Preset:      nastawa,
		Orientation: shared.BrowserDeviceOrientation(shared.BrowserDeviceOrientationPortrait),
	}
	switch nastawa {
	case shared.BrowserDevicePresetMobile:
		metryki.Width, metryki.Height, metryki.DeviceScaleFactor, metryki.Mobile = 390, 844, 3, true
	case shared.BrowserDevicePresetTablet:
		metryki.Width, metryki.Height, metryki.DeviceScaleFactor, metryki.Mobile = 820, 1180, 2, true
	default:
		metryki.Width, metryki.Height, metryki.DeviceScaleFactor, metryki.Mobile = 1280, 900, 1, false
		metryki.Orientation = shared.BrowserDeviceOrientation(shared.BrowserDeviceOrientationLandscape)
	}
	if z.Width != nil && *z.Width > 0 {
		metryki.Width = *z.Width
	}
	if z.Height != nil && *z.Height > 0 {
		metryki.Height = *z.Height
	}
	if z.DeviceScaleFactor != nil && *z.DeviceScaleFactor > 0 {
		metryki.DeviceScaleFactor = *z.DeviceScaleFactor
	}
	if z.Mobile != nil {
		metryki.Mobile = *z.Mobile
	}
	if z.Orientation != nil {
		metryki.Orientation = *z.Orientation
	}
	if z.UserAgent != nil && strings.TrimSpace(*z.UserAgent) != "" {
		metryki.UserAgent = z.UserAgent
	}
	return metryki
}

// wpisySieciowe wybiera z powiadomień protokołu żądania sieciowe strony wraz
// z ich odpowiedziami, w kolejności, w jakiej żądania ruszyły.
func wpisySieciowe(zdarzenia []zdarzenieCdp, limit int) []shared.BrowserNetworkEntry {
	wpisy := map[string]*shared.BrowserNetworkEntry{}
	kolejnosc := []string{}

	for _, zdarzenie := range zdarzenia {
		switch zdarzenie.Metoda {
		case "Network.requestWillBeSent":
			var tresc struct {
				RequestId string `json:"requestId"`
				Request   struct {
					Url     string          `json:"url"`
					Method  string          `json:"method"`
					Headers json.RawMessage `json:"headers"`
				} `json:"request"`
			}
			if err := json.Unmarshal(zdarzenie.Parametry, &tresc); err != nil {
				continue
			}
			wpis := &shared.BrowserNetworkEntry{
				RequestId: tresc.RequestId,
				Url:       tresc.Request.Url,
				Method:    tresc.Request.Method,
				StartedAt: zdarzenie.Odnotowano.UnixMilli(),
			}
			if len(tresc.Request.Headers) > 0 {
				wpis.RequestHeaders = json.RawMessage(tresc.Request.Headers)
			}
			if _, jest := wpisy[tresc.RequestId]; !jest {
				kolejnosc = append(kolejnosc, tresc.RequestId)
			}
			wpisy[tresc.RequestId] = wpis
		case "Network.responseReceived":
			var tresc struct {
				RequestId string `json:"requestId"`
				Response  struct {
					Status   int             `json:"status"`
					MimeType string          `json:"mimeType"`
					Headers  json.RawMessage `json:"headers"`
				} `json:"response"`
			}
			if err := json.Unmarshal(zdarzenie.Parametry, &tresc); err != nil {
				continue
			}
			wpis, jest := wpisy[tresc.RequestId]
			if !jest {
				continue
			}
			stan := tresc.Response.Status
			wpis.StatusCode = &stan
			if tresc.Response.MimeType != "" {
				mime := tresc.Response.MimeType
				wpis.MimeType = &mime
			}
			if len(tresc.Response.Headers) > 0 {
				wpis.ResponseHeaders = json.RawMessage(tresc.Response.Headers)
			}
		case "Network.loadingFinished":
			var tresc struct {
				RequestId string  `json:"requestId"`
				Rozmiar   float64 `json:"encodedDataLength"`
			}
			if err := json.Unmarshal(zdarzenie.Parametry, &tresc); err != nil {
				continue
			}
			wpis, jest := wpisy[tresc.RequestId]
			if !jest {
				continue
			}
			rozmiar := int64(tresc.Rozmiar)
			wpis.SizeBytes = &rozmiar
			trwanie := int(zdarzenie.Odnotowano.UnixMilli() - wpis.StartedAt)
			if trwanie >= 0 {
				wpis.DurationMs = &trwanie
			}
		}
	}

	wynik := make([]shared.BrowserNetworkEntry, 0, len(kolejnosc))
	for _, klucz := range kolejnosc {
		if limit > 0 && len(wynik) >= limit {
			break
		}
		wynik = append(wynik, *wpisy[klucz])
	}
	return wynik
}

// wpisyHar przekłada rejestr żądań sieciowych na kształt wpisów pliku HAR
// w wersji 1.2, gotowy do zapisu jako pole `entries` dokumentu HAR.
func wpisyHar(wpisy []shared.BrowserNetworkEntry) []map[string]any {
	zapis := make([]map[string]any, 0, len(wpisy))
	for _, wpis := range wpisy {
		stan := 0
		if wpis.StatusCode != nil {
			stan = *wpis.StatusCode
		}
		rozmiar := int64(0)
		if wpis.SizeBytes != nil {
			rozmiar = *wpis.SizeBytes
		}
		trwanie := 0
		if wpis.DurationMs != nil {
			trwanie = *wpis.DurationMs
		}
		mime := ""
		if wpis.MimeType != nil {
			mime = *wpis.MimeType
		}
		zapis = append(zapis, map[string]any{
			"startedDateTime": time.UnixMilli(wpis.StartedAt).UTC().Format(time.RFC3339Nano),
			"time":            trwanie,
			"request":         map[string]any{"method": wpis.Method, "url": wpis.Url, "httpVersion": "HTTP/1.1"},
			"response": map[string]any{
				"status": stan, "httpVersion": "HTTP/1.1",
				"content": map[string]any{"size": rozmiar, "mimeType": mime},
			},
			"cache":   map[string]any{},
			"timings": map[string]any{"send": 0, "wait": trwanie, "receive": 0},
		})
	}
	return zapis
}

// wpisyKonsoli wybiera z powiadomień protokołu komunikaty konsoli i błędy
// strony, zawężone do poziomów wskazanych w żądaniu.
func wpisyKonsoli(zdarzenia []zdarzenieCdp, poziomy []string, limit int) []shared.BrowserConsoleEntry {
	dopuszczony := map[string]bool{}
	for _, poziom := range poziomy {
		dopuszczony[strings.ToLower(strings.TrimSpace(poziom))] = true
	}

	wpisy := []shared.BrowserConsoleEntry{}
	for numer, zdarzenie := range zdarzenia {
		var poziom, tekst, zrodlo, adres string
		var wiersz *int

		switch zdarzenie.Metoda {
		case "Runtime.consoleAPICalled":
			var tresc struct {
				Rodzaj    string `json:"type"`
				Argumenty []struct {
					Wartosc     json.RawMessage `json:"value"`
					Opis        string          `json:"description"`
					OpisWartosc string          `json:"unserializableValue"`
				} `json:"args"`
			}
			if err := json.Unmarshal(zdarzenie.Parametry, &tresc); err != nil {
				continue
			}
			poziom = poziomKonsoli(tresc.Rodzaj)
			czesci := make([]string, 0, len(tresc.Argumenty))
			for _, argument := range tresc.Argumenty {
				if len(argument.Wartosc) > 0 {
					czesci = append(czesci, strings.Trim(string(argument.Wartosc), `"`))
					continue
				}
				if argument.Opis != "" {
					czesci = append(czesci, argument.Opis)
				}
			}
			tekst = strings.Join(czesci, " ")
			zrodlo = "console-api"
		case "Log.entryAdded":
			var tresc struct {
				Wpis struct {
					Poziom string `json:"level"`
					Tekst  string `json:"text"`
					Zrodlo string `json:"source"`
					Url    string `json:"url"`
					Wiersz int    `json:"lineNumber"`
				} `json:"entry"`
			}
			if err := json.Unmarshal(zdarzenie.Parametry, &tresc); err != nil {
				continue
			}
			poziom = poziomKonsoli(tresc.Wpis.Poziom)
			tekst, zrodlo, adres = tresc.Wpis.Tekst, tresc.Wpis.Zrodlo, tresc.Wpis.Url
			if tresc.Wpis.Wiersz > 0 {
				numer := tresc.Wpis.Wiersz
				wiersz = &numer
			}
		case "Runtime.exceptionThrown":
			var tresc struct {
				Szczegoly struct {
					Tekst  string `json:"text"`
					Url    string `json:"url"`
					Wiersz int    `json:"lineNumber"`
				} `json:"exceptionDetails"`
			}
			if err := json.Unmarshal(zdarzenie.Parametry, &tresc); err != nil {
				continue
			}
			poziom, tekst, zrodlo = string(shared.BrowserConsoleLevelError), tresc.Szczegoly.Tekst, "javascript"
			adres = tresc.Szczegoly.Url
			if tresc.Szczegoly.Wiersz > 0 {
				numer := tresc.Szczegoly.Wiersz
				wiersz = &numer
			}
		default:
			continue
		}

		if len(dopuszczony) > 0 && !dopuszczony[poziom] {
			continue
		}
		if tekst == "" {
			continue
		}
		wpis := shared.BrowserConsoleEntry{
			Id:         "cons-" + jakoLiczba(numer+1),
			Level:      shared.BrowserConsoleLevel(poziom),
			Text:       tekst,
			LineNumber: wiersz,
			LoggedAt:   zdarzenie.Odnotowano.UnixMilli(),
		}
		if zrodlo != "" {
			kopia := zrodlo
			wpis.Source = &kopia
		}
		if adres != "" {
			kopia := adres
			wpis.Url = &kopia
		}
		wpisy = append(wpisy, wpis)
		if limit > 0 && len(wpisy) >= limit {
			break
		}
	}
	return wpisy
}

// poziomKonsoli przekłada nazwy poziomów przeglądarki na wyliczenie kontraktu.
// Nazwa nieznana spada na `log` — komunikat ma się pokazać, a nie zniknąć przez
// nieznane słowo w polu poziomu.
func poziomKonsoli(nazwa string) string {
	switch strings.ToLower(strings.TrimSpace(nazwa)) {
	case "debug", "verbose":
		return string(shared.BrowserConsoleLevelDebug)
	case "info":
		return string(shared.BrowserConsoleLevelInfo)
	case "warning", "warn":
		return string(shared.BrowserConsoleLevelWarning)
	case "error", "assert":
		return string(shared.BrowserConsoleLevelError)
	default:
		return string(shared.BrowserConsoleLevelLog)
	}
}
