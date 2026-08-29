// Plik jest jedyną drogą modułu Research do sieci: pobranie zasobu spod
// adresu oraz zapytania do dostawców odkrywania źródeł Crossref, OpenAlex,
// arXiv, OpenLibrary, PubMed, DuckDuckGo i kanałów RSS/Atom.
package core

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"danacoconsole/shared"
)

const (
	// granicaSieciBadania jest limitem czasu jednego wywołania dostawcy.
	// Dostawca, który milczy minutę, jest dla Operatora dostawcą niedostępnym —
	// czekanie bez granicy zawiesiłoby okno na czas nieokreślony.
	granicaSieciBadania = 25 * time.Second
	// granicaTresciBadania ogranicza pobraną treść. Materiał badawczy bywa
	// wielki, ale rdzeń nie ma prawa wciągnąć do pamięci pliku, którego rozmiaru
	// nikt nie zapowiedział.
	granicaTresciBadania = 32 << 20
	// przedstawienieBadania wchodzi w nagłówek User-Agent. Crossref prowadzi pulę
	// grzecznościową dla wywołań, które się przedstawiają — anonimowe wywołanie
	// bywa dławione.
	przedstawienieBadania = "DanacoConsole-Research/2.0 (+https://danaco-group.pl)"
)

// klientSieciBadania jest jednym klientem HTTP modułu. Jeden na moduł, nie jeden
// na wywołanie: klient własny przy każdym żądaniu porzucałby pulę połączeń
// i otwierał nowe gniazdo do tego samego dostawcy przy każdej pozycji listy.
var klientSieciBadania = &http.Client{Timeout: granicaSieciBadania}

// pobierzBadania sprowadza treść spod adresu http albo https wraz z jej
// typem zawartości, ograniczając rozmiar odpowiedzi.
func pobierzBadania(ctx context.Context, adres string) ([]byte, string, error) {
	adres = strings.TrimSpace(adres)
	if adres == "" {
		return nil, "", bladWskazaniaBadan("żądanie bez adresu")
	}
	rozbior, err := url.Parse(adres)
	if err != nil || (rozbior.Scheme != "http" && rozbior.Scheme != "https") {
		return nil, "", bladWskazaniaBadan("adres " + adres +
			" nie jest adresem http(s) — serwer pobiera wyłącznie zasoby sieciowe")
	}

	zadanie, err := http.NewRequestWithContext(ctx, http.MethodGet, adres, nil)
	if err != nil {
		return nil, "", bladBadan(err)
	}
	zadanie.Header.Set("User-Agent", przedstawienieBadania)
	zadanie.Header.Set("Accept", "*/*")

	odpowiedz, err := klientSieciBadania.Do(zadanie)
	if err != nil {
		return nil, "", bladDostawcyBadan(adres, err.Error())
	}
	defer odpowiedz.Body.Close()

	if odpowiedz.StatusCode >= 400 {
		return nil, "", bladDostawcyBadan(adres,
			"odpowiedź o stanie "+strconv.Itoa(odpowiedz.StatusCode))
	}
	bajty, err := io.ReadAll(io.LimitReader(odpowiedz.Body, granicaTresciBadania))
	if err != nil {
		return nil, "", bladDostawcyBadan(adres, err.Error())
	}
	return bajty, odpowiedz.Header.Get("Content-Type"), nil
}

// bladDostawcyBadan nazywa niedostępność usługi zewnętrznej kodem kontraktu
// `channel_unavailable`, odróżniając ją od usterki rdzenia.
func bladDostawcyBadan(dostawca, powod string) error {
	return protokolBladBadania(shared.ErrorCodeChannelUnavailable,
		"dostawca "+dostawca+" nie odpowiedział: "+powod)
}

// odpowiedzCrossref jest wycinkiem odpowiedzi Crossref REST, z którego moduł
// składa pozycję kontraktu. Wycinek, nie pełny model: Crossref oddaje kilkadziesiąt
// pól, a `ResearchDiscoveryResult` niesie osiem.
type odpowiedzCrossref struct {
	Message struct {
		TotalResults int               `json:"total-results"`
		Items        []pozycjaCrossref `json:"items"`
	} `json:"message"`
}

// pozycjaCrossref jest jedną pracą w odpowiedzi Crossref, niosącą tytuł,
// autorów, rok wydania i identyfikator DOI.
type pozycjaCrossref struct {
	DOI    string   `json:"DOI"`
	Title  []string `json:"title"`
	URL    string   `json:"URL"`
	Author []struct {
		Given  string `json:"given"`
		Family string `json:"family"`
	} `json:"author"`
	Issued struct {
		DateParts [][]int `json:"date-parts"`
	} `json:"issued"`
	Abstract  string `json:"abstract"`
	Reference []struct {
		DOI string `json:"DOI"`
	} `json:"reference"`
	ContainerTitle []string `json:"container-title"`
}

// jakoWynikBadania przekłada pracę Crossref na pozycję kontraktu odkrycia,
// uzupełnioną o tytuł, autorów, rok i skrót.
func (p pozycjaCrossref) jakoWynikBadania() shared.ResearchDiscoveryResult {
	wynik := shared.ResearchDiscoveryResult{
		Key: "crossref:" + p.DOI, Provider: "crossref",
		Title: pierwszyNapisBadania(p.Title, "(bez tytułu)"),
	}
	if p.DOI != "" {
		doi := p.DOI
		wynik.Identifier = &doi
	}
	if p.URL != "" {
		adres := p.URL
		wynik.Url = &adres
	}
	for _, autor := range p.Author {
		wynik.Authors = append(wynik.Authors, strings.TrimSpace(autor.Given+" "+autor.Family))
	}
	if len(p.Issued.DateParts) > 0 && len(p.Issued.DateParts[0]) > 0 {
		rok := p.Issued.DateParts[0][0]
		wynik.Year = &rok
	}
	if streszczenie := bezZnacznikowBadania(p.Abstract); streszczenie != "" {
		wynik.Snippet = &streszczenie
	}
	return wynik
}

// szukajCrossref pyta Crossref REST o prace pasujące do zapytania, zawężone
// opcjonalnie przedziałem lat wydania.
func szukajCrossref(ctx context.Context, zapytanie string, odRoku, doRoku *int,
	limit int) ([]shared.ResearchDiscoveryResult, int, error) {

	parametry := url.Values{}
	parametry.Set("query", zapytanie)
	parametry.Set("rows", strconv.Itoa(limit))
	parametry.Set("mailto", "kontakt@danaco-group.pl")
	filtry := []string{}
	if odRoku != nil {
		filtry = append(filtry, "from-pub-date:"+strconv.Itoa(*odRoku)+"-01-01")
	}
	if doRoku != nil {
		filtry = append(filtry, "until-pub-date:"+strconv.Itoa(*doRoku)+"-12-31")
	}
	if len(filtry) > 0 {
		parametry.Set("filter", strings.Join(filtry, ","))
	}

	bajty, _, err := pobierzBadania(ctx, "https://api.crossref.org/works?"+parametry.Encode())
	if err != nil {
		return nil, 0, err
	}
	var odpowiedz odpowiedzCrossref
	if err := json.Unmarshal(bajty, &odpowiedz); err != nil {
		return nil, 0, bladDostawcyBadan("crossref", "odpowiedź nieczytelna: "+err.Error())
	}
	wyniki := make([]shared.ResearchDiscoveryResult, 0, len(odpowiedz.Message.Items))
	for _, pozycja := range odpowiedz.Message.Items {
		wyniki = append(wyniki, pozycja.jakoWynikBadania())
	}
	return wyniki, odpowiedz.Message.TotalResults, nil
}

// pracaCrossref pobiera jedną pracę po DOI — rozstrzygnięcie identyfikatora
// oraz punkt wyjścia snowballingu wstecz.
func pracaCrossref(ctx context.Context, doi string) (pozycjaCrossref, error) {
	bajty, _, err := pobierzBadania(ctx, "https://api.crossref.org/works/"+url.PathEscape(doi))
	if err != nil {
		return pozycjaCrossref{}, err
	}
	var odpowiedz struct {
		Message pozycjaCrossref `json:"message"`
	}
	if err := json.Unmarshal(bajty, &odpowiedz); err != nil {
		return pozycjaCrossref{}, bladDostawcyBadan("crossref", "odpowiedź nieczytelna: "+err.Error())
	}
	return odpowiedz.Message, nil
}

// pozycjaOpenAlex jest jedną pracą w odpowiedzi OpenAlex, niosącą tytuł,
// rok, autorów i status otwartego dostępu.
type pozycjaOpenAlex struct {
	ID              string `json:"id"`
	DOI             string `json:"doi"`
	Title           string `json:"display_name"`
	PublicationYear int    `json:"publication_year"`
	OpenAccess      struct {
		IsOA bool `json:"is_oa"`
	} `json:"open_access"`
	Authorships []struct {
		Author struct {
			Name string `json:"display_name"`
		} `json:"author"`
	} `json:"authorships"`
	ReferencedWorks []string `json:"referenced_works"`
	CitedByAPIURL   string   `json:"cited_by_api_url"`
}

// jakoWynikBadania przekłada pracę OpenAlex na pozycję kontraktu odkrycia,
// uzupełnioną o identyfikator DOI i autorów.
func (p pozycjaOpenAlex) jakoWynikBadania() shared.ResearchDiscoveryResult {
	wynik := shared.ResearchDiscoveryResult{
		Key: "openalex:" + p.ID, Provider: "openalex", Title: p.Title,
	}
	if p.Title == "" {
		wynik.Title = "(bez tytułu)"
	}
	if p.DOI != "" {
		doi := strings.TrimPrefix(p.DOI, "https://doi.org/")
		wynik.Identifier = &doi
		adres := p.DOI
		wynik.Url = &adres
	}
	if p.PublicationYear > 0 {
		rok := p.PublicationYear
		wynik.Year = &rok
	}
	otwarty := p.OpenAccess.IsOA
	wynik.OpenAccess = &otwarty
	for _, autorstwo := range p.Authorships {
		if autorstwo.Author.Name != "" {
			wynik.Authors = append(wynik.Authors, autorstwo.Author.Name)
		}
	}
	return wynik
}

// szukajOpenAlex pyta OpenAlex o prace pasujące do zapytania, zawężone
// opcjonalnie latami i otwartym dostępem.
func szukajOpenAlex(ctx context.Context, zapytanie string, odRoku, doRoku *int,
	tylkoOtwarte bool, limit int) ([]shared.ResearchDiscoveryResult, int, error) {

	parametry := url.Values{}
	parametry.Set("search", zapytanie)
	parametry.Set("per-page", strconv.Itoa(limit))
	parametry.Set("mailto", "kontakt@danaco-group.pl")
	filtry := []string{}
	if odRoku != nil && doRoku != nil {
		filtry = append(filtry, "publication_year:"+strconv.Itoa(*odRoku)+"-"+strconv.Itoa(*doRoku))
	} else if odRoku != nil {
		filtry = append(filtry, "from_publication_date:"+strconv.Itoa(*odRoku)+"-01-01")
	} else if doRoku != nil {
		filtry = append(filtry, "to_publication_date:"+strconv.Itoa(*doRoku)+"-12-31")
	}
	if tylkoOtwarte {
		filtry = append(filtry, "is_oa:true")
	}
	if len(filtry) > 0 {
		parametry.Set("filter", strings.Join(filtry, ","))
	}

	bajty, _, err := pobierzBadania(ctx, "https://api.openalex.org/works?"+parametry.Encode())
	if err != nil {
		return nil, 0, err
	}
	var odpowiedz struct {
		Meta struct {
			Count int `json:"count"`
		} `json:"meta"`
		Results []pozycjaOpenAlex `json:"results"`
	}
	if err := json.Unmarshal(bajty, &odpowiedz); err != nil {
		return nil, 0, bladDostawcyBadan("openalex", "odpowiedź nieczytelna: "+err.Error())
	}
	wyniki := make([]shared.ResearchDiscoveryResult, 0, len(odpowiedz.Results))
	for _, pozycja := range odpowiedz.Results {
		wyniki = append(wyniki, pozycja.jakoWynikBadania())
	}
	return wyniki, odpowiedz.Meta.Count, nil
}

// pracaOpenAlex pobiera jedną pracę po identyfikatorze DOI jako punkt
// wyjścia rozwinięcia cytowań wstecz.
func pracaOpenAlex(ctx context.Context, doi string) (pozycjaOpenAlex, error) {
	bajty, _, err := pobierzBadania(ctx, "https://api.openalex.org/works/doi:"+url.PathEscape(doi))
	if err != nil {
		return pozycjaOpenAlex{}, err
	}
	var praca pozycjaOpenAlex
	if err := json.Unmarshal(bajty, &praca); err != nil {
		return pozycjaOpenAlex{}, bladDostawcyBadan("openalex", "odpowiedź nieczytelna: "+err.Error())
	}
	return praca, nil
}

// wynikiOpenAlexZAdresuBadania pobiera listę prac spod gotowego adresu OpenAlex —
// tak przychodzi lista prac cytujących (`cited_by_api_url`), której nie da się
// złożyć z samego identyfikatora.
func wynikiOpenAlexZAdresuBadania(ctx context.Context, adres string,
	limit int) ([]shared.ResearchDiscoveryResult, int, error) {

	rozdzielnik := "?"
	if strings.Contains(adres, "?") {
		rozdzielnik = "&"
	}
	bajty, _, err := pobierzBadania(ctx, adres+rozdzielnik+"per-page="+strconv.Itoa(limit))
	if err != nil {
		return nil, 0, err
	}
	var odpowiedz struct {
		Meta struct {
			Count int `json:"count"`
		} `json:"meta"`
		Results []pozycjaOpenAlex `json:"results"`
	}
	if err := json.Unmarshal(bajty, &odpowiedz); err != nil {
		return nil, 0, bladDostawcyBadan("openalex", "odpowiedź nieczytelna: "+err.Error())
	}
	wyniki := make([]shared.ResearchDiscoveryResult, 0, len(odpowiedz.Results))
	for _, pozycja := range odpowiedz.Results {
		wyniki = append(wyniki, pozycja.jakoWynikBadania())
	}
	return wyniki, odpowiedz.Meta.Count, nil
}

// kanalArxiv jest odpowiedzią arXiv w formacie Atom, niosącą wpisy z tytułem,
// streszczeniem i datą publikacji.
type kanalArxiv struct {
	XMLName xml.Name `xml:"feed"`
	Entries []struct {
		ID        string `xml:"id"`
		Title     string `xml:"title"`
		Summary   string `xml:"summary"`
		Published string `xml:"published"`
		Authors   []struct {
			Name string `xml:"name"`
		} `xml:"author"`
	} `xml:"entry"`
}

// szukajArxiv pyta arXiv o prace pasujące do zapytania i przekłada wpisy
// kanału Atom na pozycje kontraktu.
func szukajArxiv(ctx context.Context, zapytanie string,
	limit int) ([]shared.ResearchDiscoveryResult, error) {

	parametry := url.Values{}
	parametry.Set("search_query", "all:"+zapytanie)
	parametry.Set("max_results", strconv.Itoa(limit))

	bajty, _, err := pobierzBadania(ctx, "https://export.arxiv.org/api/query?"+parametry.Encode())
	if err != nil {
		return nil, err
	}
	var kanal kanalArxiv
	if err := xml.Unmarshal(bajty, &kanal); err != nil {
		return nil, bladDostawcyBadan("arxiv", "odpowiedź nieczytelna: "+err.Error())
	}
	wyniki := make([]shared.ResearchDiscoveryResult, 0, len(kanal.Entries))
	for _, wpis := range kanal.Entries {
		adres := strings.TrimSpace(wpis.ID)
		otwarty := true
		wynik := shared.ResearchDiscoveryResult{
			Key: "arxiv:" + adres, Provider: "arxiv",
			Title: strings.TrimSpace(wpis.Title), Url: &adres, OpenAccess: &otwarty,
		}
		if fragment := strings.TrimSpace(wpis.Summary); fragment != "" {
			wynik.Snippet = &fragment
		}
		if len(wpis.Published) >= 4 {
			if rok, err := strconv.Atoi(wpis.Published[:4]); err == nil {
				wynik.Year = &rok
			}
		}
		for _, autor := range wpis.Authors {
			wynik.Authors = append(wynik.Authors, autor.Name)
		}
		wyniki = append(wyniki, wynik)
	}
	return wyniki, nil
}

// wzorzecTrafieniaWeb wyławia pozycje z odpowiedzi DuckDuckGo w postaci HTML.
// Dostawca ten nie wymaga klucza, więc wyszukiwanie webowe działa u Operatora
// bez konfiguracji — a to jest różnica między funkcją a odmową.
var wzorzecTrafieniaWeb = regexp.MustCompile(
	`(?s)<a[^>]+class="result__a"[^>]+href="([^"]+)"[^>]*>(.*?)</a>`)

// szukajWeb pyta dostawcę webowego DuckDuckGo i oddaje trafienia wyłowione
// z jego odpowiedzi w postaci HTML.
func szukajWeb(ctx context.Context, zapytanie string,
	limit int) ([]shared.ResearchDiscoveryResult, error) {

	parametry := url.Values{}
	parametry.Set("q", zapytanie)
	bajty, _, err := pobierzBadania(ctx,
		"https://html.duckduckgo.com/html/?"+parametry.Encode())
	if err != nil {
		return nil, err
	}
	trafienia := wzorzecTrafieniaWeb.FindAllStringSubmatch(string(bajty), limit)
	wyniki := make([]shared.ResearchDiscoveryResult, 0, len(trafienia))
	for _, trafienie := range trafienia {
		adres := adresTrafieniaWeb(trafienie[1])
		tytul := bezZnacznikowBadania(trafienie[2])
		if adres == "" || tytul == "" {
			continue
		}
		wyniki = append(wyniki, shared.ResearchDiscoveryResult{
			Key: "web:" + adres, Provider: "duckduckgo", Title: tytul, Url: &adres,
		})
	}
	return wyniki, nil
}

// adresTrafieniaWeb rozwija adres przekierowania dostawcy do adresu docelowego.
// Bez tego źródło zapisałoby adres pośrednika, a nie adres pracy — i przestałby
// on działać, gdy pośrednik zmieni postać odnośnika.
func adresTrafieniaWeb(surowy string) string {
	surowy = strings.TrimSpace(surowy)
	if strings.HasPrefix(surowy, "//") {
		surowy = "https:" + surowy
	}
	rozbior, err := url.Parse(surowy)
	if err != nil {
		return ""
	}
	if docelowy := rozbior.Query().Get("uddg"); docelowy != "" {
		return docelowy
	}
	if rozbior.Scheme == "" {
		return ""
	}
	return surowy
}

// rozstrzygnijISBN pobiera metadane pozycji książkowej z serwisu OpenLibrary
// po podanym numerze ISBN wprost.
func rozstrzygnijISBN(ctx context.Context, isbn string) (shared.ResearchDiscoveryResult, error) {
	bajty, _, err := pobierzBadania(ctx,
		"https://openlibrary.org/isbn/"+url.PathEscape(isbn)+".json")
	if err != nil {
		return shared.ResearchDiscoveryResult{}, err
	}
	var pozycja struct {
		Title         string `json:"title"`
		PublishDate   string `json:"publish_date"`
		NumberOfPages int    `json:"number_of_pages"`
	}
	if err := json.Unmarshal(bajty, &pozycja); err != nil {
		return shared.ResearchDiscoveryResult{}, bladDostawcyBadan("openlibrary",
			"odpowiedź nieczytelna: "+err.Error())
	}
	adres := "https://openlibrary.org/isbn/" + isbn
	kod := isbn
	wynik := shared.ResearchDiscoveryResult{
		Key: "isbn:" + isbn, Provider: "openlibrary", Title: pozycja.Title,
		Url: &adres, Identifier: &kod,
	}
	if rok := rokZTekstuBadania(pozycja.PublishDate); rok != nil {
		wynik.Year = rok
	}
	return wynik, nil
}

// rozstrzygnijPMID pobiera metadane pozycji z serwisu PubMed E-utilities po
// podanym identyfikatorze PMID.
func rozstrzygnijPMID(ctx context.Context, pmid string) (shared.ResearchDiscoveryResult, error) {
	bajty, _, err := pobierzBadania(ctx,
		"https://eutils.ncbi.nlm.nih.gov/entrez/eutils/esummary.fcgi?db=pubmed&retmode=json&id="+
			url.QueryEscape(pmid))
	if err != nil {
		return shared.ResearchDiscoveryResult{}, err
	}
	var odpowiedz struct {
		Result map[string]json.RawMessage `json:"result"`
	}
	if err := json.Unmarshal(bajty, &odpowiedz); err != nil {
		return shared.ResearchDiscoveryResult{}, bladDostawcyBadan("pubmed",
			"odpowiedź nieczytelna: "+err.Error())
	}
	surowa, jest := odpowiedz.Result[pmid]
	if !jest {
		return shared.ResearchDiscoveryResult{}, protokolBladBadania(shared.ErrorCodeNotFound,
			"PubMed nie zna pozycji o identyfikatorze "+pmid)
	}
	var pozycja struct {
		Title   string `json:"title"`
		PubDate string `json:"pubdate"`
		Authors []struct {
			Name string `json:"name"`
		} `json:"authors"`
	}
	if err := json.Unmarshal(surowa, &pozycja); err != nil {
		return shared.ResearchDiscoveryResult{}, bladDostawcyBadan("pubmed",
			"pozycja nieczytelna: "+err.Error())
	}
	adres := "https://pubmed.ncbi.nlm.nih.gov/" + pmid + "/"
	kod := pmid
	wynik := shared.ResearchDiscoveryResult{
		Key: "pmid:" + pmid, Provider: "pubmed", Title: pozycja.Title,
		Url: &adres, Identifier: &kod,
	}
	if rok := rokZTekstuBadania(pozycja.PubDate); rok != nil {
		wynik.Year = rok
	}
	for _, autor := range pozycja.Authors {
		wynik.Authors = append(wynik.Authors, autor.Name)
	}
	return wynik, nil
}

// kanalWiadomosciBadania obejmuje obie postacie kanału naraz: RSS 2.0
// (`channel/item`) i Atom (`entry`). Jeden typ, bo różnią się nazwami węzłów,
// a nie treścią — dwa parsery byłyby dwiema prawdami o jednej pozycji.
type kanalWiadomosciBadania struct {
	Kanal struct {
		Pozycje []struct {
			Tytul string `xml:"title"`
			Adres string `xml:"link"`
			Opis  string `xml:"description"`
			Data  string `xml:"pubDate"`
		} `xml:"item"`
	} `xml:"channel"`
	Wpisy []struct {
		Tytul string `xml:"title"`
		Adres struct {
			Href string `xml:"href,attr"`
		} `xml:"link"`
		Opis string `xml:"summary"`
		Data string `xml:"updated"`
	} `xml:"entry"`
}

// czytajKanal pobiera kanał wiadomości w formacie RSS albo Atom spod adresu
// i przekłada jego pozycje na wyniki odkrycia.
func czytajKanal(ctx context.Context, adres string) ([]shared.ResearchDiscoveryResult, error) {
	bajty, _, err := pobierzBadania(ctx, adres)
	if err != nil {
		return nil, err
	}
	var kanal kanalWiadomosciBadania
	if err := xml.Unmarshal(bajty, &kanal); err != nil {
		return nil, bladDostawcyBadan(adres, "kanał nieczytelny: "+err.Error())
	}
	wyniki := []shared.ResearchDiscoveryResult{}
	for _, pozycja := range kanal.Kanal.Pozycje {
		wyniki = append(wyniki, wynikKanaluBadania(pozycja.Tytul, pozycja.Adres,
			pozycja.Opis, pozycja.Data))
	}
	for _, wpis := range kanal.Wpisy {
		wyniki = append(wyniki, wynikKanaluBadania(wpis.Tytul, wpis.Adres.Href,
			wpis.Opis, wpis.Data))
	}
	return wyniki, nil
}

// wynikKanaluBadania składa pozycję kontraktu odkrycia z jednego wpisu
// kanału wiadomości RSS albo Atom.
func wynikKanaluBadania(tytul, adres, opis, data string) shared.ResearchDiscoveryResult {
	tytul = strings.TrimSpace(tytul)
	if tytul == "" {
		tytul = "(bez tytułu)"
	}
	adres = strings.TrimSpace(adres)
	wynik := shared.ResearchDiscoveryResult{
		Key: "feed:" + adres, Provider: "feed", Title: tytul,
	}
	if adres != "" {
		wynik.Url = &adres
	}
	if fragment := bezZnacznikowBadania(opis); fragment != "" {
		wynik.Snippet = &fragment
	}
	if rok := rokZTekstuBadania(data); rok != nil {
		wynik.Year = rok
	}
	return wynik
}

// wzorzecZnacznikaBadania dopasowuje pojedynczy znacznik dokumentu HTML
// albo XML do usunięcia z tekstu.
var wzorzecZnacznikaBadania = regexp.MustCompile(`(?s)<[^>]*>`)

// wzorzecRokuBadania wyławia czterocyfrowy rok z dowolnej postaci tekstu
// daty zwróconej przez dostawcę.
var wzorzecRokuBadania = regexp.MustCompile(`(19|20)\d{2}`)

// bezZnacznikowBadania sprowadza fragment dokumentu do czystego tekstu.
// Fragment zapisany ze znacznikami wyszedłby do okna jako kod, a do modelu jako
// szum — w obu miejscach jako coś innego niż zdanie, które ktoś napisał.
func bezZnacznikowBadania(tekst string) string {
	bez := wzorzecZnacznikaBadania.ReplaceAllString(tekst, " ")
	bez = strings.ReplaceAll(bez, "&nbsp;", " ")
	bez = strings.ReplaceAll(bez, "&amp;", "&")
	bez = strings.ReplaceAll(bez, "&lt;", "<")
	bez = strings.ReplaceAll(bez, "&gt;", ">")
	bez = strings.ReplaceAll(bez, "&quot;", "\"")
	bez = strings.ReplaceAll(bez, "&#39;", "'")
	return strings.Join(strings.Fields(bez), " ")
}

// rokZTekstuBadania wyławia czterocyfrowy rok publikacji z dowolnej postaci
// daty zwróconej przez dostawcę.
func rokZTekstuBadania(tekst string) *int {
	trafienie := wzorzecRokuBadania.FindString(tekst)
	if trafienie == "" {
		return nil
	}
	rok, err := strconv.Atoi(trafienie)
	if err != nil {
		return nil
	}
	return &rok
}

// pierwszyNapisBadania oddaje pierwszy niepusty napis z listy tytułów albo
// podaną wartość zastępczą, jeśli brak.
func pierwszyNapisBadania(lista []string, zastepcza string) string {
	for _, wartosc := range lista {
		if strings.TrimSpace(wartosc) != "" {
			return strings.TrimSpace(wartosc)
		}
	}
	return zastepcza
}

// tekstZDokumentuHtmlBadania wyciąga z dokumentu HTML tytuł i czystą treść —
// tryb lektury Web Clippera bez nawigacji i reklam.
func tekstZDokumentuHtmlBadania(dokument string) (string, string) {
	tytul := ""
	if trafienie := regexp.MustCompile(`(?is)<title[^>]*>(.*?)</title>`).
		FindStringSubmatch(dokument); len(trafienie) == 2 {
		tytul = bezZnacznikowBadania(trafienie[1])
	}
	// Skrypty i style wychodzą przed zdejmowaniem znaczników, bo ich treść inaczej trafi do tekstu źródła.
	bez := regexp.MustCompile(`(?is)<(script|style|nav|header|footer|aside)[^>]*>.*?</\1>`).
		ReplaceAllString(dokument, " ")
	return tytul, bezZnacznikowBadania(bez)
}
