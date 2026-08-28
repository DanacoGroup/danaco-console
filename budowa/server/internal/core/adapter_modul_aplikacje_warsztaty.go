// Moduł Apps — narzędzia obu warsztatów: mapa routingu, motyw produktu,
// eksplorator punktów końcowych, zapytanie próbne i podgląd schematu bazy,
// odczytane bezpośrednio z pracy, nie zapisane osobno.
package core

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// czasZapytaniaProbnegoApp jest granicą czekania na odpowiedź punktu końcowego.
// Bez granicy okno Backend Workspace wisiałoby tyle, ile wisi cudza usługa.
const czasZapytaniaProbnegoApp = 10 * time.Second

// granicaTresciProbnejApp przycina odczytaną odpowiedź. Kontrakt niesie ją
// w polu tekstowym jednej koperty — odpowiedź liczona w megabajtach zapchałaby
// gniazdo, a do rozpoznania wyniku wystarcza początek.
const granicaTresciProbnejApp = 64 * 1024

// Wzorce, którymi czyta się pracę Operatora. Wyrażenia regularne, nie parser
// języka: warsztat przyjmuje dowolny stos technologiczny, więc parser jednego
// z nich byłby wyborem zrobionym za Operatora.
var (
	// wzorzecTrasyPolaApp rozpoznaje deklarację trasy w postaci pola ścieżki
	// wykazu tras, niezależnie od frameworka warstwy interfejsu.
	wzorzecTrasyPolaApp = regexp.MustCompile(`(?i)\bpath\s*:\s*['"]([^'"]+)['"]`)
	// wzorzecTrasyZnacznikaApp rozpoznaje trasę zapisaną znacznikiem
	// komponentu routingu, z atrybutem ścieżki i elementem widoku.
	wzorzecTrasyZnacznikaApp = regexp.MustCompile(`(?i)<\s*Route\b[^>]*\bpath\s*=\s*["']([^"']+)["']`)
	// wzorzecWidokuTrasyApp rozpoznaje nazwę widoku obsługującego trasę, gdy
	// stoi obok jej deklaracji w tej samej linii pliku.
	wzorzecWidokuTrasyApp = regexp.MustCompile(`(?i)\b(?:component|element|view)\s*[:=]\s*["']?([A-Za-z0-9_.]+)`)
	// wzorzecPunktuKoncowegoApp rozpoznaje wiersz punktu końcowego w kontrakcie
	// API komponentu: metodę, ścieżkę i opcjonalny opis.
	wzorzecPunktuKoncowegoApp = regexp.MustCompile(
		`(?im)^\s*(GET|POST|PUT|PATCH|DELETE|GRAPHQL)\s+(\S+)\s*(?:[-—:]\s*(.*))?$`)
	// wzorzecTabeliApp rozpoznaje polecenie tworzenia tabeli w plikach
	// warstwy backendu i wyodrębnia nazwę tabeli spod polecenia.
	wzorzecTabeliApp = regexp.MustCompile(`(?is)CREATE\s+TABLE\s+(?:IF\s+NOT\s+EXISTS\s+)?` +
		"[`\"\\[]?([A-Za-z0-9_.]+)[`\"\\]]?\\s*\\(")
)

// WypiszTrasy obsługuje komendę odczytu mapy routingu warstwy interfejsu,
// złożonej z tras odnalezionych w plikach frontendu.
func (a *adapterAplikacji) WypiszTrasy(ctx context.Context,
	z shared.AppsRouteListRequest) (shared.AppsRouteListResponse, error) {

	okno, err := oknoAplikacji(z.WindowId, "apps.route.list")
	if err != nil {
		return shared.AppsRouteListResponse{}, err
	}
	pliki, err := a.plikiWarstwyApp(ctx, okno, shared.AppWorkspaceLayerFrontend)
	if err != nil {
		return shared.AppsRouteListResponse{}, err
	}

	// Mapa po ścieżce trasy: ta sama trasa w dwóch plikach jest jedną trasą
	// produktu.
	znalezione := map[string]shared.AppRoute{}
	for _, plik := range pliki {
		for _, trafienie := range wzorzecTrasyZnacznikaApp.FindAllStringSubmatch(plik.Tresc, -1) {
			dodajTraseApp(znalezione, trafienie[1], "", plik.KomponentID)
		}
		for _, trafienie := range wzorzecTrasyPolaApp.FindAllStringSubmatchIndex(plik.Tresc, -1) {
			sciezka := plik.Tresc[trafienie[2]:trafienie[3]]
			dodajTraseApp(znalezione, sciezka, widokPrzyTrasieApp(plik.Tresc, trafienie[1]),
				plik.KomponentID)
		}
	}

	trasy := make([]shared.AppRoute, 0, len(znalezione))
	for _, trasa := range znalezione {
		trasy = append(trasy, trasa)
	}
	sort.SliceStable(trasy, func(i, j int) bool { return trasy[i].Path < trasy[j].Path })
	return shared.AppsRouteListResponse{Routes: trasy, Total: len(trasy)}, nil
}

// dodajTraseApp wpisuje trasę do mapy tras, nie gubiąc widoku znalezionego
// wcześniej dla tej samej ścieżki.
func dodajTraseApp(mapa map[string]shared.AppRoute, sciezka, widok string, komponent *string) {
	sciezka = strings.TrimSpace(sciezka)
	if sciezka == "" || !strings.HasPrefix(sciezka, "/") {
		return
	}
	zastana, jest := mapa[sciezka]
	if !jest {
		zastana = shared.AppRoute{Path: sciezka, ComponentId: komponent}
	}
	if zastana.ViewName == nil {
		zastana.ViewName = wskaznikNapisuApp(widok)
	}
	if zastana.ComponentId == nil {
		zastana.ComponentId = komponent
	}
	mapa[sciezka] = zastana
}

// widokPrzyTrasieApp szuka nazwy widoku w tej samej deklaracji, co ścieżka —
// najbliżej za nią, w granicach jednego wpisu wykazu tras.
func widokPrzyTrasieApp(tresc string, od int) string {
	const zasieg = 160
	do := od + zasieg
	if do > len(tresc) {
		do = len(tresc)
	}
	trafienie := wzorzecWidokuTrasyApp.FindStringSubmatch(tresc[od:do])
	if trafienie == nil {
		return ""
	}
	return trafienie[1]
}

// PobierzMotyw obsługuje `apps.theme.get`. Okno bez zapisanego motywu oddaje
// wynik z pustym polem — kontrakt ma je jako opcjonalne, bo produkt świeżo
// założony jeszcze motywu nie ma.
func (a *adapterAplikacji) PobierzMotyw(ctx context.Context,
	z shared.AppsThemeGetRequest) (shared.AppsThemeGetResponse, error) {

	okno, err := oknoAplikacji(z.WindowId, "apps.theme.get")
	if err != nil {
		return shared.AppsThemeGetResponse{}, err
	}
	wiersz, err := a.repozytorium.MotywApp(ctx, okno)
	if err != nil {
		if isBrakWierszaApp(err) {
			return shared.AppsThemeGetResponse{}, nil
		}
		return shared.AppsThemeGetResponse{}, bladAplikacji(err)
	}
	return shared.AppsThemeGetResponse{Theme: json.RawMessage(wiersz.Tresc)}, nil
}

// UstawMotyw obsługuje komendę zapisu motywu produktu. Treść jest surowym
// JSON-em kontraktu, którego rdzeń nie rozkłada — sprawdza wyłącznie, czy to
// w ogóle jest poprawny JSON.
func (a *adapterAplikacji) UstawMotyw(ctx context.Context,
	z shared.AppsThemeSetRequest) (shared.AppsThemeSetResponse, error) {

	okno, err := oknoAplikacji(z.WindowId, "apps.theme.set")
	if err != nil {
		return shared.AppsThemeSetResponse{}, err
	}
	// Pustka ma tu dwie postacie: pole pominięte w żądaniu oraz pole niosące
	// wartość pustą.
	if len(z.Theme) == 0 || strings.TrimSpace(string(z.Theme)) == "null" {
		return shared.AppsThemeSetResponse{}, bladWskazaniaAplikacji(
			"apps.theme.set wymaga motywu")
	}
	if !json.Valid(z.Theme) {
		return shared.AppsThemeSetResponse{}, bladWskazaniaAplikacji(
			"motyw produktu nie jest poprawnym JSON-em")
	}
	zapisany, err := a.repozytorium.ZapiszMotywApp(ctx, dane.MotywApp{
		Okno: okno, Tresc: string(z.Theme),
	})
	if err != nil {
		return shared.AppsThemeSetResponse{}, bladAplikacji(err)
	}
	return shared.AppsThemeSetResponse{Theme: json.RawMessage(zapisany.Tresc)}, nil
}

// WypiszPunktyKoncowe obsługuje `apps.endpoint.list` — eksplorator punktów
// końcowych złożony z kontraktów API komponentów (czoło pliku).
func (a *adapterAplikacji) WypiszPunktyKoncowe(ctx context.Context,
	z shared.AppsEndpointListRequest) (shared.AppsEndpointListResponse, error) {

	okno, err := oknoAplikacji(z.WindowId, "apps.endpoint.list")
	if err != nil {
		return shared.AppsEndpointListResponse{}, err
	}
	zawezenie := strings.TrimSpace(wartoscTekstu(z.ComponentId))

	komponenty, err := a.komponentyOknaApp(ctx, okno)
	if err != nil {
		return shared.AppsEndpointListResponse{}, err
	}

	punkty := []shared.AppEndpoint{}
	for _, komponent := range komponenty {
		if zawezenie != "" && komponent.KodZewnetrzny != zawezenie {
			continue
		}
		if komponent.KontraktAPI == nil {
			continue
		}
		kod := komponent.KodZewnetrzny
		for _, trafienie := range wzorzecPunktuKoncowegoApp.FindAllStringSubmatch(*komponent.KontraktAPI, -1) {
			metoda := shared.AppEndpointMethod(strings.ToLower(trafienie[1]))
			sciezka := trafienie[2]
			punkt := shared.AppEndpoint{
				Id:          kod + " " + strings.ToUpper(trafienie[1]) + " " + sciezka,
				ComponentId: &kod,
				Method:      metoda,
				Path:        sciezka,
				Description: wskaznikNapisuApp(strings.TrimSpace(trafienie[3])),
			}
			// Stan punktu wynika z pracy: ścieżka odnaleziona w plikach
			// backendu znaczy punkt zaimplementowany.
			stan := shared.AppEndpointStatus(shared.AppEndpointStatusDraft)
			punkt.Status = &stan
			punkty = append(punkty, punkt)
		}
	}

	if len(punkty) > 0 {
		pliki, err := a.plikiWarstwyApp(ctx, okno, shared.AppWorkspaceLayerBackend)
		if err != nil {
			return shared.AppsEndpointListResponse{}, err
		}
		var backend strings.Builder
		for _, plik := range pliki {
			backend.WriteString(plik.Tresc)
			backend.WriteString("\n")
		}
		tresc := backend.String()
		for indeks := range punkty {
			if !strings.Contains(tresc, punkty[indeks].Path) {
				continue
			}
			stan := shared.AppEndpointStatus(shared.AppEndpointStatusImplemented)
			punkty[indeks].Status = &stan
		}
	}

	sort.SliceStable(punkty, func(i, j int) bool { return punkty[i].Id < punkty[j].Id })
	return shared.AppsEndpointListResponse{Endpoints: punkty, Total: len(punkty)}, nil
}

// ZapytajPunktKoncowy obsługuje `apps.endpoint.probe`: wysyła prawdziwe
// zapytanie i mierzy prawdziwą odpowiedź (czoło pliku).
func (a *adapterAplikacji) ZapytajPunktKoncowy(ctx context.Context,
	z shared.AppsEndpointProbeRequest) (shared.AppsEndpointProbeResponse, error) {

	okno, err := oknoAplikacji(z.WindowId, "apps.endpoint.probe")
	if err != nil {
		return shared.AppsEndpointProbeResponse{}, err
	}
	if err := sprawdzMetodePunktuApp(z.Method); err != nil {
		return shared.AppsEndpointProbeResponse{}, err
	}
	if strings.TrimSpace(z.Path) == "" {
		return shared.AppsEndpointProbeResponse{}, bladWskazaniaAplikacji(
			"apps.endpoint.probe wymaga ścieżki punktu końcowego")
	}
	if z.Environment != nil {
		if err := sprawdzSrodowiskoWdrozenia(*z.Environment); err != nil {
			return shared.AppsEndpointProbeResponse{}, err
		}
	}

	podstawa, err := a.adresBazowyProbyApp(ctx, okno, z.Environment)
	if err != nil {
		return shared.AppsEndpointProbeResponse{}, err
	}
	adres := strings.TrimSuffix(podstawa, "/") + "/" + strings.TrimPrefix(z.Path, "/")

	// GraphQL nie jest metodą HTTP: zapytanie zawsze jedzie żądaniem POST,
	// zgodnie z protokołem.
	metoda := strings.ToUpper(string(z.Method))
	if z.Method == shared.AppEndpointMethodGraphql {
		metoda = http.MethodPost
	}

	var tresc io.Reader
	if z.Body != nil && *z.Body != "" {
		tresc = strings.NewReader(*z.Body)
	}
	zapytanie, err := http.NewRequestWithContext(ctx, metoda, adres, tresc)
	if err != nil {
		return shared.AppsEndpointProbeResponse{}, bladWskazaniaAplikacji(
			"nie można złożyć zapytania do " + adres + ": " + err.Error())
	}
	if len(z.Headers) > 0 {
		naglowki := map[string]string{}
		if err := json.Unmarshal(z.Headers, &naglowki); err != nil {
			return shared.AppsEndpointProbeResponse{}, bladWskazaniaAplikacji(
				"nagłówki zapytania próbnego nie są obiektem napisów: " + err.Error())
		}
		for nazwa, wartosc := range naglowki {
			zapytanie.Header.Set(nazwa, wartosc)
		}
	}

	klient := &http.Client{Timeout: czasZapytaniaProbnegoApp}
	poczatek := time.Now()
	odpowiedz, err := klient.Do(zapytanie)
	czas := int(time.Since(poczatek).Milliseconds())
	if err != nil {
		// Niepowodzenie sieci nie jest odmową komendy, tylko wynikiem z powodem.
		a.dopiszDziennikApp(ctx, okno, nil, nil,
			"zapytanie próbne "+metoda+" "+adres+" nie doszło: "+err.Error())
		szczegol := err.Error()
		return shared.AppsEndpointProbeResponse{Result: shared.AppEndpointProbeResult{
			StatusCode: 0, DurationMs: czas, ErrorDetail: &szczegol,
		}}, nil
	}
	defer odpowiedz.Body.Close()

	bajty, err := io.ReadAll(io.LimitReader(odpowiedz.Body, granicaTresciProbnejApp))
	if err != nil {
		szczegol := "odpowiedź przyszła, ale jej treści nie dało się odczytać: " + err.Error()
		return shared.AppsEndpointProbeResponse{Result: shared.AppEndpointProbeResult{
			StatusCode: odpowiedz.StatusCode, DurationMs: czas, ErrorDetail: &szczegol,
		}}, nil
	}

	naglowki := map[string]string{}
	for nazwa := range odpowiedz.Header {
		naglowki[nazwa] = odpowiedz.Header.Get(nazwa)
	}
	surowe, _ := json.Marshal(naglowki)
	trescOdpowiedzi := string(bajty)

	a.dopiszDziennikApp(ctx, okno, nil, nil,
		"zapytanie próbne "+metoda+" "+adres+" → "+strconv.Itoa(odpowiedz.StatusCode)+
			" w "+strconv.Itoa(czas)+" ms")

	return shared.AppsEndpointProbeResponse{Result: shared.AppEndpointProbeResult{
		StatusCode: odpowiedz.StatusCode,
		DurationMs: czas,
		Headers:    surowe,
		Body:       &trescOdpowiedzi,
	}}, nil
}

// PobierzSchemat obsługuje `apps.schema.get` — podgląd schematu bazy produktu
// odczytany z poleceń `CREATE TABLE` w warstwie backendu (czoło pliku).
func (a *adapterAplikacji) PobierzSchemat(ctx context.Context,
	z shared.AppsSchemaGetRequest) (shared.AppsSchemaGetResponse, error) {

	okno, err := oknoAplikacji(z.WindowId, "apps.schema.get")
	if err != nil {
		return shared.AppsSchemaGetResponse{}, err
	}

	komponent := strings.TrimSpace(wartoscTekstu(z.ComponentId))
	if komponent == "" {
		komponenty, err := a.komponentyOknaApp(ctx, okno)
		if err != nil {
			return shared.AppsSchemaGetResponse{}, err
		}
		for _, kandydat := range komponenty {
			if kandydat.Rodzaj == string(shared.AppComponentKindDatabase) {
				komponent = kandydat.KodZewnetrzny
				break
			}
		}
	}

	pliki, err := a.plikiWarstwyApp(ctx, okno, shared.AppWorkspaceLayerBackend)
	if err != nil {
		return shared.AppsSchemaGetResponse{}, err
	}

	tabele := []shared.AppSchemaTable{}
	for _, plik := range pliki {
		// Plik przypisany do innego komponentu nie opisuje tej bazy; plik bez
		// przypisania liczy się zawsze.
		if komponent != "" && plik.KomponentID != nil && *plik.KomponentID != komponent {
			continue
		}
		tabele = append(tabele, tabeleZeZrodlaApp(plik.Tresc)...)
	}
	sort.SliceStable(tabele, func(i, j int) bool { return tabele[i].Name < tabele[j].Name })

	if len(tabele) == 0 {
		// Brak schematu to puste pole, nie odmowa: backend bez tabeli jest
		// stanem normalnym.
		return shared.AppsSchemaGetResponse{}, nil
	}
	schemat := shared.AppSchema{
		ComponentId: wskaznikNapisuApp(komponent),
		Tables:      tabele,
		ReadAt:      time.Now().UTC().UnixMilli(),
	}
	return shared.AppsSchemaGetResponse{Schema: &schemat}, nil
}

// tabeleZeZrodlaApp wyciąga tabele wraz z kolumnami z treści pliku backendu,
// rozpoznając każde polecenie tworzenia tabeli.
func tabeleZeZrodlaApp(zrodlo string) []shared.AppSchemaTable {
	tabele := []shared.AppSchemaTable{}
	for _, trafienie := range wzorzecTabeliApp.FindAllStringSubmatchIndex(zrodlo, -1) {
		nazwa := zrodlo[trafienie[2]:trafienie[3]]
		cialo, koniec := cialoNawiasuApp(zrodlo, trafienie[1]-1)
		if koniec < 0 {
			continue
		}
		kolumny, odwolania := kolumnyTabeliApp(cialo)
		tabele = append(tabele, shared.AppSchemaTable{
			Name: nazwa, Columns: kolumny, References: odwolania,
		})
	}
	return tabele
}

// cialoNawiasuApp oddaje treść między nawiasem otwierającym pod wskazanym
// indeksem a jego domknięciem; niedomknięty nawias daje -1.
func cialoNawiasuApp(zrodlo string, otwarcie int) (string, int) {
	glebokosc := 0
	for indeks := otwarcie; indeks < len(zrodlo); indeks++ {
		switch zrodlo[indeks] {
		case '(':
			glebokosc++
		case ')':
			glebokosc--
			if glebokosc == 0 {
				return zrodlo[otwarcie+1 : indeks], indeks
			}
		}
	}
	return "", -1
}

// kolumnyTabeliApp rozkłada ciało polecenia `CREATE TABLE` na kolumny oraz
// tabele, do których odwołują się więzy obce.
func kolumnyTabeliApp(cialo string) ([]shared.AppSchemaColumn, []string) {
	kolumny := []shared.AppSchemaColumn{}
	odwolania := []string{}
	widziane := map[string]struct{}{}

	for _, wiersz := range rozdzielPrzecinkamiApp(cialo) {
		przyciety := strings.TrimSpace(wiersz)
		if przyciety == "" {
			continue
		}
		wielkie := strings.ToUpper(przyciety)
		if odwolanie := odwolanieWiezuApp(przyciety); odwolanie != "" {
			if _, jest := widziane[odwolanie]; !jest {
				widziane[odwolanie] = struct{}{}
				odwolania = append(odwolania, odwolanie)
			}
		}
		// Wiersz więzu tabeli nie jest kolumną — zaczyna się słowem kluczowym.
		if strings.HasPrefix(wielkie, "PRIMARY KEY") || strings.HasPrefix(wielkie, "FOREIGN KEY") ||
			strings.HasPrefix(wielkie, "UNIQUE") || strings.HasPrefix(wielkie, "CHECK") ||
			strings.HasPrefix(wielkie, "CONSTRAINT") {
			continue
		}

		czesci := strings.Fields(przyciety)
		if len(czesci) == 0 {
			continue
		}
		nazwa := strings.Trim(czesci[0], "`\"[]")
		typ := ""
		if len(czesci) > 1 {
			typ = strings.Trim(czesci[1], ",")
		}
		pusta := !strings.Contains(wielkie, "NOT NULL")
		klucz := strings.Contains(wielkie, "PRIMARY KEY")
		kolumny = append(kolumny, shared.AppSchemaColumn{
			Name: nazwa, Type: typ, Nullable: &pusta, PrimaryKey: &klucz,
		})
	}
	return kolumny, odwolania
}

// rozdzielPrzecinkamiApp dzieli ciało polecenia po przecinkach najwyższego
// poziomu — przecinek wewnątrz `CHECK(...)` albo `DECIMAL(10,2)` nie kończy
// definicji kolumny.
func rozdzielPrzecinkamiApp(cialo string) []string {
	czesci := []string{}
	glebokosc := 0
	poczatek := 0
	for indeks := 0; indeks < len(cialo); indeks++ {
		switch cialo[indeks] {
		case '(':
			glebokosc++
		case ')':
			glebokosc--
		case ',':
			if glebokosc == 0 {
				czesci = append(czesci, cialo[poczatek:indeks])
				poczatek = indeks + 1
			}
		}
	}
	return append(czesci, cialo[poczatek:])
}

// odwolanieWiezuApp oddaje nazwę tabeli wskazanej słowem kluczowym więzu
// obcego, gdy definicja kolumny go niesie.
var wzorzecOdwolaniaApp = regexp.MustCompile(
	"(?i)REFERENCES\\s+[`\"\\[]?([A-Za-z0-9_.]+)[`\"\\]]?")

func odwolanieWiezuApp(wiersz string) string {
	trafienie := wzorzecOdwolaniaApp.FindStringSubmatch(wiersz)
	if trafienie == nil {
		return ""
	}
	return trafienie[1]
}

// komponentyOknaApp zwraca komponenty bieżącej architektury okna; okno bez
// architektury oddaje pustkę, nie odmowę — punktów końcowych ani schematu po
// prostu jeszcze nie ma z czego przeczytać.
func (a *adapterAplikacji) komponentyOknaApp(ctx context.Context,
	okno string) ([]dane.KomponentArchitektury, error) {

	architektura, err := a.repozytorium.ArchitekturaOkna(ctx, okno)
	if err != nil {
		if isBrakWierszaApp(err) {
			return nil, nil
		}
		return nil, bladAplikacji(err)
	}
	komponenty, err := a.repozytorium.Komponenty(ctx, architektura.ID)
	if err != nil {
		return nil, bladAplikacji(err)
	}
	return komponenty, nil
}

// adresBazowyProbyApp rozstrzyga, dokąd idzie zapytanie próbne: pod domenę
// środowiska, a gdy jej nie ma — pod adres stojącego podglądu. Brak obu jest
// odmową z powodem: zapytanie donikąd nie ma jak zwrócić kodu odpowiedzi.
func (a *adapterAplikacji) adresBazowyProbyApp(ctx context.Context, okno string,
	srodowisko *shared.AppDeployEnvironment) (string, error) {

	kod := ""
	if srodowisko != nil {
		kod = string(*srodowisko)
	}
	if kod != "" {
		wiersz, err := a.repozytorium.SrodowiskoApp(ctx, okno, kod)
		if err == nil && wiersz.Domena != nil && *wiersz.Domena != "" {
			return zAdresemHttpApp(*wiersz.Domena), nil
		}
		if err != nil && !isBrakWierszaApp(err) {
			return "", bladAplikacji(err)
		}
	} else {
		srodowiska, err := a.repozytorium.SrodowiskaApp(ctx, okno)
		if err != nil {
			return "", bladAplikacji(err)
		}
		for _, wiersz := range srodowiska {
			if wiersz.Domena != nil && *wiersz.Domena != "" {
				return zAdresemHttpApp(*wiersz.Domena), nil
			}
		}
	}

	if adres := a.adresPodgladuApp(okno); adres != "" {
		return adres, nil
	}
	return "", bladWskazaniaAplikacji(
		"okno " + okno + " nie ma dokąd wysłać zapytania próbnego — nadaj domenę środowiska " +
			"(apps.deployment.domain.set) albo podnieś podgląd (apps.preview.start)")
}

// zAdresemHttpApp dokłada schemat, gdy domena go nie niesie. Domena zapisywana
// przez `apps.deployment.domain.set` jest nazwą hosta, nie adresem.
func zAdresemHttpApp(domena string) string {
	if strings.HasPrefix(domena, "http://") || strings.HasPrefix(domena, "https://") {
		return domena
	}
	return "http://" + domena
}

// sprawdzMetodePunktuApp dopuszcza wyłącznie metody protokołu HTTP, które
// kontrakt zna dla punktu końcowego.
func sprawdzMetodePunktuApp(metoda shared.AppEndpointMethod) error {
	switch metoda {
	case shared.AppEndpointMethodGet, shared.AppEndpointMethodPost,
		shared.AppEndpointMethodPut, shared.AppEndpointMethodPatch,
		shared.AppEndpointMethodDelete, shared.AppEndpointMethodGraphql:
		return nil
	}
	return bladWskazaniaAplikacji("nieznana metoda punktu końcowego " +
		strconv.Quote(string(metoda)) + " — dopuszczalne: get, post, put, patch, delete, graphql")
}
