// Moduł API Client obsługuje zakładkę Dev Tools komendami
// `developer.api.request`, `developer.api.collection.save`,
// `developer.api.collection.list` i `developer.api.openapi.import`, na
// bibliotekach wkompilowanych w rdzeń, bez procesu potomnego.
package core

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

const (
	// przedrostekKolekcjiApi znakuje identyfikator kolekcji zapytań, nadawany
	// przy zakładaniu nowej kolekcji albo przy imporcie z kontraktu OpenAPI.
	przedrostekKolekcjiApi = "apic-"
	// czasZapytaniaApi jest domyślną granicą czasu jednego zapytania HTTP,
	// stosowaną, gdy żądanie nie wskazuje własnej granicy czasu.
	czasZapytaniaApi = 30 * time.Second
	// najdluzszeZapytanieApi jest granicą, której żądanie nie przekroczy nawet
	// wtedy, gdy poprosi o więcej. Zapytanie wiszące pół godziny trzyma połączenie
	// klienta i wątek rdzenia, a odpowiedź, na którą nikt już nie czeka, nie jest
	// odpowiedzią.
	najdluzszeZapytanieApi = 5 * time.Minute
	// najwiekszaOdpowiedzApi jest granicą treści wciąganej do odpowiedzi.
	// Podgląd odpowiedzi w oknie i tak nie pokaże więcej, a plik stumegabajtowy
	// wciągnięty do pamięci rdzenia byłby ceną bez odbiorcy.
	najwiekszaOdpowiedzApi = 8 << 20
)

// WykonajZapytanieApi obsługuje `developer.api.request` i wykonuje zapytanie
// HTTP po podstawieniu wartości środowiska w adresie, nagłówkach i treści.
func (a *adapterDevelopera) WykonajZapytanieApi(ctx context.Context,
	z shared.DeveloperApiRequestRequest) (shared.DeveloperApiRequestResponse, error) {

	okno, err := a.oknoDevelopera(z.WindowId)
	if err != nil {
		return shared.DeveloperApiRequestResponse{}, err
	}
	// Zapytanie HTTP zmienia stan po drugiej stronie sieci, więc tryb
	// planistyczny je wyklucza.
	if err := sprawdzZmianeSystemu(okno.TrybUprawnien, "wykonanie zapytania HTTP"); err != nil {
		return shared.DeveloperApiRequestResponse{}, err
	}

	podstawienia, err := a.podstawieniaSrodowiska(ctx, okno.Id, z.EnvironmentId)
	if err != nil {
		return shared.DeveloperApiRequestResponse{}, err
	}

	adres := podstawWSzablonie(strings.TrimSpace(z.Url), podstawienia)
	if adres == "" {
		return shared.DeveloperApiRequestResponse{}, bladZadaniaDevelopera(
			"zapytanie wymaga adresu")
	}
	if !strings.HasPrefix(adres, "http://") && !strings.HasPrefix(adres, "https://") {
		return shared.DeveloperApiRequestResponse{}, bladZadaniaDevelopera(
			"adres zapytania ma zaczynać się od http:// albo https://; otrzymano " + adres)
	}
	metoda := strings.ToUpper(strings.TrimSpace(z.Method))
	if metoda == "" {
		metoda = http.MethodGet
	}

	var tresc io.Reader
	if z.Body != nil && *z.Body != "" {
		tresc = strings.NewReader(podstawWSzablonie(*z.Body, podstawienia))
	}

	granica := czasZapytaniaApi
	if z.TimeoutMs != nil && *z.TimeoutMs > 0 {
		granica = time.Duration(*z.TimeoutMs) * time.Millisecond
		if granica > najdluzszeZapytanieApi {
			granica = najdluzszeZapytanieApi
		}
	}
	kontekst, przerwij := context.WithTimeout(ctx, granica)
	defer przerwij()

	zadanie, err := http.NewRequestWithContext(kontekst, metoda, adres, tresc)
	if err != nil {
		return shared.DeveloperApiRequestResponse{}, bladZadaniaDevelopera(
			"nie można złożyć zapytania: " + err.Error())
	}
	naglowkiZadania(zadanie, z.Headers, z.BodyKind, podstawienia)

	poczatek := time.Now()
	// Klient jest własny, nie domyślny: domyślny nie ma granicy czasu.
	klient := &http.Client{Timeout: granica}
	odpowiedz, err := klient.Do(zadanie)
	if err != nil {
		return shared.DeveloperApiRequestResponse{}, bladWykonaniaDevelopera(
			"zapytanie do " + adres + " nie doszło do skutku: " + err.Error())
	}
	defer odpowiedz.Body.Close()

	bajty, err := io.ReadAll(io.LimitReader(odpowiedz.Body, najwiekszaOdpowiedzApi))
	if err != nil && !errors.Is(err, io.EOF) {
		return shared.DeveloperApiRequestResponse{}, bladWykonaniaDevelopera(
			"odpowiedź z " + adres + " urwała się w trakcie odczytu: " + err.Error())
	}
	czas := time.Since(poczatek).Milliseconds()

	wynik := shared.ApiResponse{
		Status:     odpowiedz.StatusCode,
		DurationMs: czas,
	}
	if odpowiedz.Status != "" {
		wynik.StatusText = wskaznikTekstu(odpowiedz.Status)
	}
	if naglowki, err := json.Marshal(splaszczNaglowki(odpowiedz.Header)); err == nil {
		wynik.Headers = naglowki
	}
	if len(bajty) > 0 {
		wynik.Body = wskaznikTekstu(string(bajty))
	}
	rozmiar := int64(len(bajty))
	wynik.SizeBytes = &rozmiar
	return shared.DeveloperApiRequestResponse{Response: wynik}, nil
}

// naglowkiZadania wpisuje nagłówki żądania wraz z podstawieniami środowiska.
//
// Rodzaj treści (`bodyKind`) uzupełnia nagłówek `Content-Type` tylko wtedy, gdy
// wołający sam go nie podał: nagłówek jawny jest zawsze mocniejszy od domysłu.
func naglowkiZadania(zadanie *http.Request, surowe json.RawMessage, rodzajTresci *string,
	podstawienia map[string]string) {

	if len(surowe) > 0 {
		wpisy := map[string]string{}
		if err := json.Unmarshal(surowe, &wpisy); err == nil {
			for nazwa, wartosc := range wpisy {
				zadanie.Header.Set(nazwa, podstawWSzablonie(wartosc, podstawienia))
			}
		}
	}
	if zadanie.Header.Get("Content-Type") != "" || rodzajTresci == nil {
		return
	}
	switch strings.ToLower(strings.TrimSpace(*rodzajTresci)) {
	case "json":
		zadanie.Header.Set("Content-Type", "application/json")
	case "form":
		zadanie.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	case "text":
		zadanie.Header.Set("Content-Type", "text/plain; charset=utf-8")
	}
}

// splaszczNaglowki sprowadza nagłówki odpowiedzi do mapy nazwa → wartość.
// Nagłówek powtórzony kilka razy (np. `Set-Cookie`) zostaje sklejony przecinkiem,
// bo kontrakt niesie mapę, a zgubienie drugiej wartości byłoby gorsze od jej
// sklejenia.
func splaszczNaglowki(naglowki http.Header) map[string]string {
	wynik := make(map[string]string, len(naglowki))
	for nazwa, wartosci := range naglowki {
		wynik[nazwa] = strings.Join(wartosci, ", ")
	}
	return wynik
}

// podstawWSzablonie zastępuje `{{nazwa}}` wartością środowiska; nazwa
// nieznana zostaje w tekście nietknięta, zamiast zniknąć w pustym miejscu.
func podstawWSzablonie(tresc string, podstawienia map[string]string) string {
	if len(podstawienia) == 0 || !strings.Contains(tresc, "{{") {
		return tresc
	}
	for nazwa, wartosc := range podstawienia {
		tresc = strings.ReplaceAll(tresc, "{{"+nazwa+"}}", wartosc)
	}
	return tresc
}

// podstawieniaSrodowiska składa wartości wskazanego środowiska z kolekcji
// zapytań okna, przeszukując wszystkie kolekcje po nazwie środowiska.
func (a *adapterDevelopera) podstawieniaSrodowiska(ctx context.Context, oknoKod string,
	srodowisko *string) (map[string]string, error) {

	if srodowisko == nil || strings.TrimSpace(*srodowisko) == "" || a.repozytorium == nil {
		return nil, nil
	}
	szukane := strings.TrimSpace(*srodowisko)
	kolekcje, err := a.repozytorium.KolekcjeApi(ctx, oknoKod, "")
	if err != nil {
		return nil, bladWykonaniaDevelopera("nie można odczytać kolekcji zapytań: " + err.Error())
	}
	for _, kolekcja := range kolekcje {
		if kolekcja.Srodowiska == nil {
			continue
		}
		var srodowiska map[string]map[string]string
		if err := json.Unmarshal([]byte(*kolekcja.Srodowiska), &srodowiska); err != nil {
			continue
		}
		if wartosci, jest := srodowiska[szukane]; jest {
			return wartosci, nil
		}
	}
	return nil, bladZasobuDevelopera(
		"żadna kolekcja okna nie ma środowiska o nazwie " + szukane)
}

// ZapiszKolekcjeApi obsługuje `developer.api.collection.save` i zakłada
// kolekcję zapytań nową albo zmienia zastaną, wraz ze środowiskami.
func (a *adapterDevelopera) ZapiszKolekcjeApi(ctx context.Context,
	z shared.DeveloperApiCollectionSaveRequest) (shared.DeveloperApiCollectionSaveResponse, error) {

	okno, err := a.oknoDevelopera(z.WindowId)
	if err != nil {
		return shared.DeveloperApiCollectionSaveResponse{}, err
	}
	if a.repozytorium == nil {
		return shared.DeveloperApiCollectionSaveResponse{}, bladZasobuDevelopera(
			"rdzeń nie ma miejsca na kolekcje zapytań")
	}
	nazwa := strings.TrimSpace(z.Name)
	if nazwa == "" {
		return shared.DeveloperApiCollectionSaveResponse{}, bladZadaniaDevelopera(
			"kolekcja zapytań wymaga nazwy")
	}
	if !json.Valid(z.Requests) {
		return shared.DeveloperApiCollectionSaveResponse{}, bladZadaniaDevelopera(
			"wykaz zapytań kolekcji nie jest poprawnym dokumentem JSON")
	}

	kod := strings.TrimSpace(tekstWskazaniaDevelopera(z.CollectionId))
	if kod == "" {
		kod = nowyIdentyfikator(przedrostekKolekcjiApi)
	}
	wiersz := dane.KolekcjaApi{
		Kod:       kod,
		OknoKod:   okno.Id,
		Nazwa:     nazwa,
		Zapytania: string(z.Requests),
	}
	if len(z.Environments) > 0 && json.Valid(z.Environments) {
		wiersz.Srodowiska = wskaznikTekstu(string(z.Environments))
	}
	if err := a.repozytorium.ZapiszKolekcjeApi(ctx, wiersz); err != nil {
		return shared.DeveloperApiCollectionSaveResponse{}, bladWykonaniaDevelopera(
			"nie można zapisać kolekcji " + nazwa + ": " + err.Error())
	}

	zapisane, err := a.repozytorium.KolekcjeApi(ctx, okno.Id, kod)
	if err != nil || len(zapisane) == 0 {
		return shared.DeveloperApiCollectionSaveResponse{}, bladWykonaniaDevelopera(
			"kolekcja " + nazwa + " zapisała się, lecz nie daje się odczytać")
	}
	return shared.DeveloperApiCollectionSaveResponse{
		Collection: kolekcjaKontraktu(zapisane[0]),
	}, nil
}

// WykazKolekcjiApi obsługuje `developer.api.collection.list` i oddaje
// kolekcje zapytań okna, z możliwością zawężenia do jednej kolekcji.
func (a *adapterDevelopera) WykazKolekcjiApi(ctx context.Context,
	z shared.DeveloperApiCollectionListRequest) (shared.DeveloperApiCollectionListResponse, error) {

	okno, err := a.oknoDevelopera(z.WindowId)
	if err != nil {
		return shared.DeveloperApiCollectionListResponse{}, err
	}
	if a.repozytorium == nil {
		return shared.DeveloperApiCollectionListResponse{}, bladZasobuDevelopera(
			"rdzeń nie ma miejsca na kolekcje zapytań")
	}
	kod := strings.TrimSpace(tekstWskazaniaDevelopera(z.CollectionId))
	wiersze, err := a.repozytorium.KolekcjeApi(ctx, okno.Id, kod)
	if err != nil {
		return shared.DeveloperApiCollectionListResponse{}, bladWykonaniaDevelopera(
			"nie można odczytać kolekcji zapytań okna " + okno.Id + ": " + err.Error())
	}
	kolekcje := make([]shared.ApiCollection, 0, len(wiersze))
	for _, wiersz := range wiersze {
		kolekcje = append(kolekcje, kolekcjaKontraktu(wiersz))
	}
	return shared.DeveloperApiCollectionListResponse{Collections: kolekcje}, nil
}

// kolekcjaKontraktu przekłada wiersz kolekcji z bazy danych na kształt
// odpowiedzi zgodny z kontraktem, jaki widzi klient.
func kolekcjaKontraktu(wiersz dane.KolekcjaApi) shared.ApiCollection {
	kolekcja := shared.ApiCollection{
		Id:        wiersz.Kod,
		WindowId:  wiersz.OknoKod,
		Name:      wiersz.Nazwa,
		Requests:  json.RawMessage(wiersz.Zapytania),
		UpdatedAt: chwilaBazy(wiersz.Zmieniono),
	}
	if wiersz.Srodowiska != nil {
		kolekcja.Environments = json.RawMessage(*wiersz.Srodowiska)
	}
	return kolekcja
}

// ImportujOpenapi obsługuje `developer.api.openapi.import` i wytwarza
// kolekcję zapytań z kontraktu, wziętego z pliku repozytorium albo z adresu.
func (a *adapterDevelopera) ImportujOpenapi(ctx context.Context,
	z shared.DeveloperApiOpenapiImportRequest) (shared.DeveloperApiOpenapiImportResponse, error) {

	okno, err := a.oknoDevelopera(z.WindowId)
	if err != nil {
		return shared.DeveloperApiOpenapiImportResponse{}, err
	}
	if a.repozytorium == nil {
		return shared.DeveloperApiOpenapiImportResponse{}, bladZasobuDevelopera(
			"rdzeń nie ma miejsca na kolekcje zapytań")
	}

	dokument, zrodlo, err := a.wczytajKontraktOpenapi(ctx, okno.Id, z)
	if err != nil {
		return shared.DeveloperApiOpenapiImportResponse{}, err
	}

	zapytania, adres := zapytaniaZKontraktu(dokument)
	if len(zapytania) == 0 {
		return shared.DeveloperApiOpenapiImportResponse{}, bladZadaniaDevelopera(
			"kontrakt " + zrodlo + " nie opisuje ani jednej ścieżki, więc nie ma z czego " +
				"wytworzyć zapytań")
	}

	nazwa := strings.TrimSpace(tekstWskazaniaDevelopera(z.CollectionName))
	if nazwa == "" {
		nazwa = tytulKontraktu(dokument, zrodlo)
	}
	tresc, err := json.Marshal(zapytania)
	if err != nil {
		return shared.DeveloperApiOpenapiImportResponse{}, bladWykonaniaDevelopera(
			"nie można złożyć wykazu zapytań kolekcji: " + err.Error())
	}

	wiersz := dane.KolekcjaApi{
		Kod:       nowyIdentyfikator(przedrostekKolekcjiApi),
		OknoKod:   okno.Id,
		Nazwa:     nazwa,
		Zapytania: string(tresc),
	}
	// Serwery kontraktu wchodzą jako środowisko `openapi` z `{{baseUrl}}`.
	if adres != "" {
		srodowiska, err := json.Marshal(map[string]map[string]string{
			"openapi": {"baseUrl": adres},
		})
		if err == nil {
			wiersz.Srodowiska = wskaznikTekstu(string(srodowiska))
		}
	}
	if err := a.repozytorium.ZapiszKolekcjeApi(ctx, wiersz); err != nil {
		return shared.DeveloperApiOpenapiImportResponse{}, bladWykonaniaDevelopera(
			"nie można zapisać kolekcji z importu: " + err.Error())
	}

	zapisane, err := a.repozytorium.KolekcjeApi(ctx, okno.Id, wiersz.Kod)
	if err != nil || len(zapisane) == 0 {
		return shared.DeveloperApiOpenapiImportResponse{}, bladWykonaniaDevelopera(
			"kolekcja z importu zapisała się, lecz nie daje się odczytać")
	}
	return shared.DeveloperApiOpenapiImportResponse{
		Collection:   kolekcjaKontraktu(zapisane[0]),
		RequestCount: len(zapytania),
	}, nil
}

// wczytajKontraktOpenapi bierze kontrakt OpenAPI z pliku repozytorium albo
// z adresu wskazanego w żądaniu i zwraca też źródło, z którego przyszedł.
func (a *adapterDevelopera) wczytajKontraktOpenapi(ctx context.Context, oknoKod string,
	z shared.DeveloperApiOpenapiImportRequest) (*openapi3.T, string, error) {

	czytnik := openapi3.NewLoader()
	czytnik.IsExternalRefsAllowed = true

	if z.Path != nil && strings.TrimSpace(*z.Path) != "" {
		_, sciezka, err := a.plikOkna(oknoKod, *z.Path)
		if err != nil {
			return nil, "", err
		}
		bajty, err := os.ReadFile(sciezka)
		if err != nil {
			return nil, "", bladZasobuDevelopera(
				"nie można odczytać kontraktu " + *z.Path + ": " + err.Error())
		}
		dokument, err := czytnik.LoadFromData(bajty)
		if err != nil {
			return nil, "", bladZadaniaDevelopera(
				"plik " + *z.Path + " nie jest czytelnym kontraktem OpenAPI: " + err.Error())
		}
		return dokument, sciezka, nil
	}

	if z.Url == nil || strings.TrimSpace(*z.Url) == "" {
		return nil, "", bladZadaniaDevelopera(
			"import kontraktu wymaga wskazania pliku repozytorium albo adresu")
	}
	adres := strings.TrimSpace(*z.Url)
	kontekst, przerwij := context.WithTimeout(ctx, czasZapytaniaApi)
	defer przerwij()

	zadanie, err := http.NewRequestWithContext(kontekst, http.MethodGet, adres, nil)
	if err != nil {
		return nil, "", bladZadaniaDevelopera("adres kontraktu jest niepoprawny: " + err.Error())
	}
	odpowiedz, err := (&http.Client{Timeout: czasZapytaniaApi}).Do(zadanie)
	if err != nil {
		return nil, "", bladWykonaniaDevelopera(
			"nie można pobrać kontraktu z " + adres + ": " + err.Error())
	}
	defer odpowiedz.Body.Close()
	if odpowiedz.StatusCode >= 400 {
		return nil, "", bladZasobuDevelopera(
			"adres " + adres + " odpowiedział stanem " + odpowiedz.Status)
	}
	bajty, err := io.ReadAll(io.LimitReader(odpowiedz.Body, najwiekszaOdpowiedzApi))
	if err != nil {
		return nil, "", bladWykonaniaDevelopera("odczyt kontraktu urwał się: " + err.Error())
	}
	dokument, err := czytnik.LoadFromData(bajty)
	if err != nil {
		return nil, "", bladZadaniaDevelopera(
			"treść spod " + adres + " nie jest czytelnym kontraktem OpenAPI: " + err.Error())
	}
	return dokument, adres, nil
}

// zapytanieZKontraktu jest jednym zapytaniem kolekcji wytworzonym ze ścieżki
// kontraktu OpenAPI, gotowym do zapisu w polu `requests` kolekcji.
type zapytanieZKontraktu struct {
	Name    string            `json:"name"`
	Method  string            `json:"method"`
	Url     string            `json:"url"`
	Headers map[string]string `json:"headers,omitempty"`
	Body    string            `json:"body,omitempty"`
	Summary string            `json:"summary,omitempty"`
}

// zapytaniaZKontraktu przekłada ścieżki kontraktu na zapytania kolekcji,
// w ustalonej kolejności ścieżki, a potem metody.
func zapytaniaZKontraktu(dokument *openapi3.T) ([]zapytanieZKontraktu, string) {
	adres := ""
	if dokument.Servers != nil && len(dokument.Servers) > 0 {
		adres = strings.TrimSuffix(dokument.Servers[0].URL, "/")
	}
	if dokument.Paths == nil {
		return nil, adres
	}

	sciezki := make([]string, 0, dokument.Paths.Len())
	for sciezka := range dokument.Paths.Map() {
		sciezki = append(sciezki, sciezka)
	}
	sort.Strings(sciezki)

	zapytania := make([]zapytanieZKontraktu, 0, len(sciezki))
	for _, sciezka := range sciezki {
		pozycja := dokument.Paths.Value(sciezka)
		if pozycja == nil {
			continue
		}
		metody := pozycja.Operations()
		nazwyMetod := make([]string, 0, len(metody))
		for metoda := range metody {
			nazwyMetod = append(nazwyMetod, metoda)
		}
		sort.Strings(nazwyMetod)

		for _, metoda := range nazwyMetod {
			czynnosc := metody[metoda]
			zapytanie := zapytanieZKontraktu{
				Name:   nazwaZapytaniaKontraktu(czynnosc, metoda, sciezka),
				Method: metoda,
				Url:    "{{baseUrl}}" + sciezka,
			}
			if czynnosc.Summary != "" {
				zapytanie.Summary = czynnosc.Summary
			}
			if czynnosc.RequestBody != nil && czynnosc.RequestBody.Value != nil {
				if _, jest := czynnosc.RequestBody.Value.Content["application/json"]; jest {
					zapytanie.Headers = map[string]string{"Content-Type": "application/json"}
					zapytanie.Body = "{}"
				}
			}
			zapytania = append(zapytania, zapytanie)
		}
	}
	return zapytania, adres
}

// nazwaZapytaniaKontraktu dobiera czytelną nazwę pozycji kolekcji z operacji
// kontraktu: identyfikator operacji, potem streszczenie, potem metodę i ścieżkę.
func nazwaZapytaniaKontraktu(czynnosc *openapi3.Operation, metoda, sciezka string) string {
	if czynnosc.OperationID != "" {
		return czynnosc.OperationID
	}
	if czynnosc.Summary != "" {
		return czynnosc.Summary
	}
	return metoda + " " + sciezka
}

// tytulKontraktu bierze nazwę kolekcji z tytułu kontraktu, a gdy go nie ma —
// ze źródła, z którego kontrakt przyszedł.
func tytulKontraktu(dokument *openapi3.T, zrodlo string) string {
	if dokument.Info != nil && strings.TrimSpace(dokument.Info.Title) != "" {
		return strings.TrimSpace(dokument.Info.Title)
	}
	return "Kontrakt " + zrodlo
}

// tekstWskazaniaDevelopera oddaje treść wskaźnika tekstowego albo tekst
// pusty, gdy wskaźnik jest nieustawiony.
func tekstWskazaniaDevelopera(wskaznik *string) string {
	if wskaznik == nil {
		return ""
	}
	return *wskaznik
}
