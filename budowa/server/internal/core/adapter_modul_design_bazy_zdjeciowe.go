// Plik obsługuje ośmiu dostawców baz zdjęciowych modułu Design: wyszukanie
// i wciągnięcie zasobu, wraz z jego pochodzeniem i licencją.
package core

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	// limitCzasuDostawcyZdjecDesignu ogranicza jedno wywołanie dostawcy. Bez
	// granicy jeden dostawca, który nie odpowiada, trzymałby całą komendę.
	limitCzasuDostawcyZdjecDesignu = 20 * time.Second

	// limitPobraniaZdjeciaDesignu jest górną granicą bajtów wciąganych od
	// dostawcy — obrona rdzenia przed odpowiedzią bez końca.
	limitPobraniaZdjeciaDesignu = 64 << 20

	// domyslnyLimitWynikowDostawcyDesignu jest liczbą zasobów na dostawcę, gdy
	// Operator granicy nie podał.
	domyslnyLimitWynikowDostawcyDesignu = 20

	// granicaWynikowDostawcyDesignu chroni odpowiedź przed wykazem tysiącznym,
	// nawet gdy żądanie poda granicę większą.
	granicaWynikowDostawcyDesignu = 100

	// przedrostekBytuSejfuZdjecDesignu znakuje wpisy sejfu z kluczami dostawców.
	// Jeden przedrostek dla wszystkich, żeby Operator wiedział, czego szukać
	// w sejfie.
	przedrostekBytuSejfuZdjecDesignu = "design.stock."
)

// zasobDostawcyZdjecDesignu to jeden zasób odczytany od dostawcy — postać
// wspólna dla wszystkich ośmiu.
type zasobDostawcyZdjecDesignu struct {
	Identyfikator string
	Tytul         string
	Podglad       string
	Pelny         string
	Licencja      string
	Autor         string
	Odsylacz      string
}

// dostawcaZdjecDesignu opisuje jedną bazę zdjęciową: jej nazwę, wymaganie
// klucza oraz drogi wyszukania i pobrania zasobu.
type dostawcaZdjecDesignu struct {
	Nazwa        string
	WymagaKlucza bool
	// Szukaj oddaje zasoby spełniające frazę. Klucz jest pusty dla dostawców bez
	// klucza.
	Szukaj func(ctx context.Context, klient *http.Client, klucz, fraza string,
		limit int) ([]zasobDostawcyZdjecDesignu, error)
	// Pobierz oddaje jeden zasób po jego identyfikatorze u dostawcy — droga
	// wciągnięcia.
	Pobierz func(ctx context.Context, klient *http.Client, klucz,
		identyfikator string) (zasobDostawcyZdjecDesignu, error)
}

// dostawcyZdjecDesignu wylicza dostawców w kolejności zapytywania: najpierw ci
// bez klucza (nagłówek pliku).
func dostawcyZdjecDesignu() []dostawcaZdjecDesignu {
	return []dostawcaZdjecDesignu{
		{Nazwa: "openverse", Szukaj: szukajOpenverseDesignu, Pobierz: pobierzOpenverseDesignu},
		{Nazwa: "wikimedia", Szukaj: szukajWikimediaDesignu, Pobierz: pobierzWikimediaDesignu},
		{Nazwa: "met", Szukaj: szukajMetDesignu, Pobierz: pobierzMetDesignu},
		{Nazwa: "nasa", Szukaj: szukajNasaDesignu, Pobierz: pobierzNasaDesignu},
		{Nazwa: "unsplash", WymagaKlucza: true, Szukaj: szukajUnsplashDesignu,
			Pobierz: pobierzUnsplashDesignu},
		{Nazwa: "pexels", WymagaKlucza: true, Szukaj: szukajPexelsDesignu,
			Pobierz: pobierzPexelsDesignu},
		{Nazwa: "pixabay", WymagaKlucza: true, Szukaj: szukajPixabayDesignu,
			Pobierz: pobierzPixabayDesignu},
		{Nazwa: "smithsonian", WymagaKlucza: true, Szukaj: szukajSmithsonianDesignu,
			Pobierz: pobierzSmithsonianDesignu},
	}
}

// dostawcaZdjecPoNazwieDesignu odnajduje dostawcę po nazwie i mówi, czy
// dostawca o takiej nazwie jest znany.
func dostawcaZdjecPoNazwieDesignu(nazwa string) (dostawcaZdjecDesignu, bool) {
	szukana := strings.ToLower(strings.TrimSpace(nazwa))
	for _, dostawca := range dostawcyZdjecDesignu() {
		if dostawca.Nazwa == szukana {
			return dostawca, true
		}
	}
	return dostawcaZdjecDesignu{}, false
}

// nazwyDostawcowZdjecDesignu oddaje nazwy dostawców posortowane alfabetycznie
// — wchodzą w treść odmów wskazujących dostawcę nieznanego.
func nazwyDostawcowZdjecDesignu() []string {
	nazwy := []string{}
	for _, dostawca := range dostawcyZdjecDesignu() {
		nazwy = append(nazwy, dostawca.Nazwa)
	}
	sort.Strings(nazwy)
	return nazwy
}

// kluczDostawcyZdjecDesignu czyta klucz dostawcy z sejfu poświadczeń rdzenia
// i oddaje klucz pusty, gdy sejf albo wpis w nim nie istnieje.
func (a *adapterDesignu) kluczDostawcyZdjecDesignu(ctx context.Context, dostawca string) string {
	if a.sejf == nil {
		return ""
	}
	klucz, jest := a.sejf.Odczytaj(ctx, przedrostekBytuSejfuZdjecDesignu+dostawca)
	if !jest {
		return ""
	}
	return strings.TrimSpace(klucz)
}

// klientDostawcyZdjecDesignu składa klienta HTTP z granicą czasu, aby dostawca
// bez odpowiedzi nie trzymał wywołania bez końca.
func klientDostawcyZdjecDesignu() *http.Client {
	return &http.Client{Timeout: limitCzasuDostawcyZdjecDesignu}
}

// odczytajJsonDostawcyDesignu wykonuje żądanie do dostawcy i rozkłada jego
// odpowiedź zapisaną w postaci JSON do wskazanego celu.
func odczytajJsonDostawcyDesignu(ctx context.Context, klient *http.Client, adres string,
	naglowki map[string]string, cel any) error {

	zadanie, err := http.NewRequestWithContext(ctx, http.MethodGet, adres, nil)
	if err != nil {
		return err
	}
	// Nagłówek User-Agent jest wymagany przez Wikimedia i przez api.data.gov.
	zadanie.Header.Set("User-Agent", "DanacoConsole/1.0 (moduł Design, bazy zdjęciowe)")
	zadanie.Header.Set("Accept", "application/json")
	for nazwa, wartosc := range naglowki {
		zadanie.Header.Set(nazwa, wartosc)
	}
	odpowiedz, err := klient.Do(zadanie)
	if err != nil {
		return err
	}
	defer odpowiedz.Body.Close()
	if odpowiedz.StatusCode < 200 || odpowiedz.StatusCode > 299 {
		return fmt.Errorf("odpowiedź %s", odpowiedz.Status)
	}
	bajty, err := io.ReadAll(io.LimitReader(odpowiedz.Body, limitPobraniaZdjeciaDesignu))
	if err != nil {
		return err
	}
	return json.Unmarshal(bajty, cel)
}

// pobierzBajtyZdjeciaDesignu wciąga bajty zasobu spod adresu dostawcy wraz
// z zapisanym rodzajem treści i odmawia treści pustej.
func pobierzBajtyZdjeciaDesignu(ctx context.Context, klient *http.Client,
	adres string) ([]byte, string, error) {

	zadanie, err := http.NewRequestWithContext(ctx, http.MethodGet, adres, nil)
	if err != nil {
		return nil, "", err
	}
	zadanie.Header.Set("User-Agent", "DanacoConsole/1.0 (moduł Design, bazy zdjęciowe)")
	odpowiedz, err := klient.Do(zadanie)
	if err != nil {
		return nil, "", err
	}
	defer odpowiedz.Body.Close()
	if odpowiedz.StatusCode < 200 || odpowiedz.StatusCode > 299 {
		return nil, "", fmt.Errorf("odpowiedź %s", odpowiedz.Status)
	}
	bajty, err := io.ReadAll(io.LimitReader(odpowiedz.Body, limitPobraniaZdjeciaDesignu))
	if err != nil {
		return nil, "", err
	}
	if len(bajty) == 0 {
		return nil, "", fmt.Errorf("pusta treść pod adresem zasobu")
	}
	return bajty, odpowiedz.Header.Get("Content-Type"), nil
}

// Openverse — katalog materiałów na wolnych licencjach, dostawca bez klucza.

func szukajOpenverseDesignu(ctx context.Context, klient *http.Client, _, fraza string,
	limit int) ([]zasobDostawcyZdjecDesignu, error) {

	var odpowiedz struct {
		Results []struct {
			ID      string `json:"id"`
			Title   string `json:"title"`
			URL     string `json:"url"`
			Thumb   string `json:"thumbnail"`
			License string `json:"license"`
			Creator string `json:"creator"`
			Foreign string `json:"foreign_landing_url"`
		} `json:"results"`
	}
	adres := fmt.Sprintf("https://api.openverse.org/v1/images/?q=%s&page_size=%d",
		url.QueryEscape(fraza), limit)
	if err := odczytajJsonDostawcyDesignu(ctx, klient, adres, nil, &odpowiedz); err != nil {
		return nil, err
	}
	zasoby := make([]zasobDostawcyZdjecDesignu, 0, len(odpowiedz.Results))
	for _, wynik := range odpowiedz.Results {
		zasoby = append(zasoby, zasobDostawcyZdjecDesignu{
			Identyfikator: wynik.ID, Tytul: wynik.Title, Podglad: wynik.Thumb,
			Pelny: wynik.URL, Licencja: wynik.License, Autor: wynik.Creator,
			Odsylacz: wynik.Foreign,
		})
	}
	return zasoby, nil
}

func pobierzOpenverseDesignu(ctx context.Context, klient *http.Client, _,
	identyfikator string) (zasobDostawcyZdjecDesignu, error) {

	var odpowiedz struct {
		ID      string `json:"id"`
		Title   string `json:"title"`
		URL     string `json:"url"`
		License string `json:"license"`
		Creator string `json:"creator"`
		Foreign string `json:"foreign_landing_url"`
	}
	adres := "https://api.openverse.org/v1/images/" + url.PathEscape(identyfikator) + "/"
	if err := odczytajJsonDostawcyDesignu(ctx, klient, adres, nil, &odpowiedz); err != nil {
		return zasobDostawcyZdjecDesignu{}, err
	}
	return zasobDostawcyZdjecDesignu{
		Identyfikator: odpowiedz.ID, Tytul: odpowiedz.Title, Pelny: odpowiedz.URL,
		Licencja: odpowiedz.License, Autor: odpowiedz.Creator, Odsylacz: odpowiedz.Foreign,
	}, nil
}

// Wikimedia Commons — identyfikatorem zasobu jest tytuł pliku.

func szukajWikimediaDesignu(ctx context.Context, klient *http.Client, _, fraza string,
	limit int) ([]zasobDostawcyZdjecDesignu, error) {

	var odpowiedz struct {
		Query struct {
			Pages map[string]struct {
				Title     string `json:"title"`
				ImageInfo []struct {
					URL            string `json:"url"`
					ThumbURL       string `json:"thumburl"`
					DescriptionURL string `json:"descriptionurl"`
					ExtMetadata    map[string]struct {
						Value string `json:"value"`
					} `json:"extmetadata"`
				} `json:"imageinfo"`
			} `json:"pages"`
		} `json:"query"`
	}
	adres := fmt.Sprintf("https://commons.wikimedia.org/w/api.php?action=query&format=json"+
		"&generator=search&gsrsearch=%s&gsrnamespace=6&gsrlimit=%d"+
		"&prop=imageinfo&iiprop=url%%7Cextmetadata&iiurlwidth=640",
		url.QueryEscape("filetype:bitmap "+fraza), limit)
	if err := odczytajJsonDostawcyDesignu(ctx, klient, adres, nil, &odpowiedz); err != nil {
		return nil, err
	}
	zasoby := make([]zasobDostawcyZdjecDesignu, 0, len(odpowiedz.Query.Pages))
	for _, strona := range odpowiedz.Query.Pages {
		zasob := zasobDostawcyZdjecDesignu{Identyfikator: strona.Title, Tytul: strona.Title}
		if len(strona.ImageInfo) > 0 {
			opis := strona.ImageInfo[0]
			zasob.Pelny = opis.URL
			zasob.Podglad = opis.ThumbURL
			zasob.Odsylacz = opis.DescriptionURL
			if wpis, jest := opis.ExtMetadata["LicenseShortName"]; jest {
				zasob.Licencja = oczyscZnacznikiHtmlDesignu(wpis.Value)
			}
			if wpis, jest := opis.ExtMetadata["Artist"]; jest {
				zasob.Autor = oczyscZnacznikiHtmlDesignu(wpis.Value)
			}
		}
		zasoby = append(zasoby, zasob)
	}
	// Kolejność wyników Commons jest kolejnością mapy JSON, czyli żadną;
	// porządek po tytule ją naprawia.
	sort.SliceStable(zasoby, func(i, j int) bool { return zasoby[i].Tytul < zasoby[j].Tytul })
	return zasoby, nil
}

func pobierzWikimediaDesignu(ctx context.Context, klient *http.Client, _,
	identyfikator string) (zasobDostawcyZdjecDesignu, error) {

	var odpowiedz struct {
		Query struct {
			Pages map[string]struct {
				Title     string `json:"title"`
				ImageInfo []struct {
					URL            string `json:"url"`
					DescriptionURL string `json:"descriptionurl"`
					ExtMetadata    map[string]struct {
						Value string `json:"value"`
					} `json:"extmetadata"`
				} `json:"imageinfo"`
			} `json:"pages"`
		} `json:"query"`
	}
	adres := "https://commons.wikimedia.org/w/api.php?action=query&format=json&titles=" +
		url.QueryEscape(identyfikator) + "&prop=imageinfo&iiprop=url%7Cextmetadata"
	if err := odczytajJsonDostawcyDesignu(ctx, klient, adres, nil, &odpowiedz); err != nil {
		return zasobDostawcyZdjecDesignu{}, err
	}
	for _, strona := range odpowiedz.Query.Pages {
		if len(strona.ImageInfo) == 0 {
			continue
		}
		opis := strona.ImageInfo[0]
		zasob := zasobDostawcyZdjecDesignu{
			Identyfikator: strona.Title, Tytul: strona.Title,
			Pelny: opis.URL, Odsylacz: opis.DescriptionURL,
		}
		if wpis, jest := opis.ExtMetadata["LicenseShortName"]; jest {
			zasob.Licencja = oczyscZnacznikiHtmlDesignu(wpis.Value)
		}
		if wpis, jest := opis.ExtMetadata["Artist"]; jest {
			zasob.Autor = oczyscZnacznikiHtmlDesignu(wpis.Value)
		}
		return zasob, nil
	}
	return zasobDostawcyZdjecDesignu{}, fmt.Errorf("Commons nie zna pliku %q", identyfikator)
}

// oczyscZnacznikiHtmlDesignu zdejmuje znaczniki HTML z pola metadanych. Commons
// podaje autora i licencję jako fragment HTML z odsyłaczami, a licencja zapisana
// razem ze znacznikami byłaby w oknie nieczytelna.
func oczyscZnacznikiHtmlDesignu(tresc string) string {
	var wynik strings.Builder
	wZnaczniku := false
	for _, znak := range tresc {
		switch {
		case znak == '<':
			wZnaczniku = true
		case znak == '>':
			wZnaczniku = false
		case !wZnaczniku:
			wynik.WriteRune(znak)
		}
	}
	return strings.TrimSpace(wynik.String())
}

// ── Metropolitan Museum of Art ──────────────────────────────────────────────

func szukajMetDesignu(ctx context.Context, klient *http.Client, _, fraza string,
	limit int) ([]zasobDostawcyZdjecDesignu, error) {

	var wyszukanie struct {
		ObjectIDs []int `json:"objectIDs"`
	}
	adres := "https://collectionapi.metmuseum.org/public/collection/v1/search?hasImages=true&q=" +
		url.QueryEscape(fraza)
	if err := odczytajJsonDostawcyDesignu(ctx, klient, adres, nil, &wyszukanie); err != nil {
		return nil, err
	}
	// Met oddaje same identyfikatory; opis każdego trzeba dobrać osobno.
	zasoby := []zasobDostawcyZdjecDesignu{}
	for _, numer := range wyszukanie.ObjectIDs {
		if len(zasoby) >= limit {
			break
		}
		zasob, err := pobierzMetDesignu(ctx, klient, "", strconv.Itoa(numer))
		if err != nil || zasob.Pelny == "" {
			continue
		}
		zasoby = append(zasoby, zasob)
	}
	return zasoby, nil
}

func pobierzMetDesignu(ctx context.Context, klient *http.Client, _,
	identyfikator string) (zasobDostawcyZdjecDesignu, error) {

	var obiekt struct {
		ObjectID          int    `json:"objectID"`
		Title             string `json:"title"`
		PrimaryImage      string `json:"primaryImage"`
		PrimaryImageSmall string `json:"primaryImageSmall"`
		ArtistDisplayName string `json:"artistDisplayName"`
		ObjectURL         string `json:"objectURL"`
		IsPublicDomain    bool   `json:"isPublicDomain"`
	}
	adres := "https://collectionapi.metmuseum.org/public/collection/v1/objects/" +
		url.PathEscape(identyfikator)
	if err := odczytajJsonDostawcyDesignu(ctx, klient, adres, nil, &obiekt); err != nil {
		return zasobDostawcyZdjecDesignu{}, err
	}
	// Met ma w zbiorach też dzieła pod prawem autorskim, nie domenę publiczną.
	licencja := "Met — obraz udostępniony bez zgody na domenę publiczną"
	if obiekt.IsPublicDomain {
		licencja = "CC0 1.0 (domena publiczna, Met Open Access)"
	}
	return zasobDostawcyZdjecDesignu{
		Identyfikator: strconv.Itoa(obiekt.ObjectID), Tytul: obiekt.Title,
		Podglad: obiekt.PrimaryImageSmall, Pelny: obiekt.PrimaryImage,
		Licencja: licencja, Autor: obiekt.ArtistDisplayName, Odsylacz: obiekt.ObjectURL,
	}, nil
}

// ── NASA Image and Video Library ────────────────────────────────────────────

func szukajNasaDesignu(ctx context.Context, klient *http.Client, _, fraza string,
	limit int) ([]zasobDostawcyZdjecDesignu, error) {

	var odpowiedz struct {
		Collection struct {
			Items []struct {
				Data []struct {
					NasaID       string `json:"nasa_id"`
					Title        string `json:"title"`
					Center       string `json:"center"`
					Photographer string `json:"photographer"`
				} `json:"data"`
				Links []struct {
					Href string `json:"href"`
					Rel  string `json:"rel"`
				} `json:"links"`
			} `json:"items"`
		} `json:"collection"`
	}
	adres := fmt.Sprintf("https://images-api.nasa.gov/search?q=%s&media_type=image&page_size=%d",
		url.QueryEscape(fraza), limit)
	if err := odczytajJsonDostawcyDesignu(ctx, klient, adres, nil, &odpowiedz); err != nil {
		return nil, err
	}
	zasoby := make([]zasobDostawcyZdjecDesignu, 0, len(odpowiedz.Collection.Items))
	for _, wpis := range odpowiedz.Collection.Items {
		if len(wpis.Data) == 0 {
			continue
		}
		zasob := zasobDostawcyZdjecDesignu{
			Identyfikator: wpis.Data[0].NasaID, Tytul: wpis.Data[0].Title,
			// NASA ma wyjątki od braku ochrony (logo agencji, wizerunki astronautów).
			Licencja: "NASA Media Usage Guidelines (bez ochrony prawem autorskim, " +
				"z wyjątkami dla logo i wizerunków)",
			Autor:    pierwszyNiepustyDesignu(wpis.Data[0].Photographer, wpis.Data[0].Center),
			Odsylacz: "https://images.nasa.gov/details/" + wpis.Data[0].NasaID,
		}
		for _, odsylacz := range wpis.Links {
			if odsylacz.Rel == "preview" {
				zasob.Podglad = odsylacz.Href
				break
			}
		}
		zasoby = append(zasoby, zasob)
	}
	return zasoby, nil
}

func pobierzNasaDesignu(ctx context.Context, klient *http.Client, _,
	identyfikator string) (zasobDostawcyZdjecDesignu, error) {

	var odpowiedz struct {
		Collection struct {
			Items []struct {
				Href string `json:"href"`
			} `json:"items"`
		} `json:"collection"`
	}
	adres := "https://images-api.nasa.gov/asset/" + url.PathEscape(identyfikator)
	if err := odczytajJsonDostawcyDesignu(ctx, klient, adres, nil, &odpowiedz); err != nil {
		return zasobDostawcyZdjecDesignu{}, err
	}
	// Wpisy są uszeregowane od największego; bierze się pierwszy plik obrazu.
	pelny := ""
	for _, wpis := range odpowiedz.Collection.Items {
		maly := strings.ToLower(wpis.Href)
		if strings.HasSuffix(maly, ".jpg") || strings.HasSuffix(maly, ".jpeg") ||
			strings.HasSuffix(maly, ".png") || strings.HasSuffix(maly, ".tif") {
			pelny = wpis.Href
			break
		}
	}
	if pelny == "" {
		return zasobDostawcyZdjecDesignu{}, fmt.Errorf(
			"zasób %q nie ma u NASA pliku obrazu — są tylko metadane", identyfikator)
	}
	return zasobDostawcyZdjecDesignu{
		Identyfikator: identyfikator, Tytul: identyfikator, Pelny: pelny,
		Licencja: "NASA Media Usage Guidelines (bez ochrony prawem autorskim, " +
			"z wyjątkami dla logo i wizerunków)",
		Odsylacz: "https://images.nasa.gov/details/" + identyfikator,
	}, nil
}

// ── Unsplash ────────────────────────────────────────────────────────────────

func szukajUnsplashDesignu(ctx context.Context, klient *http.Client, klucz, fraza string,
	limit int) ([]zasobDostawcyZdjecDesignu, error) {

	var odpowiedz struct {
		Results []struct {
			ID             string `json:"id"`
			Description    string `json:"description"`
			AltDescription string `json:"alt_description"`
			Urls           struct {
				Regular string `json:"regular"`
				Small   string `json:"small"`
			} `json:"urls"`
			Links struct {
				HTML string `json:"html"`
			} `json:"links"`
			User struct {
				Name string `json:"name"`
			} `json:"user"`
		} `json:"results"`
	}
	adres := fmt.Sprintf("https://api.unsplash.com/search/photos?query=%s&per_page=%d",
		url.QueryEscape(fraza), limit)
	naglowki := map[string]string{"Authorization": "Client-ID " + klucz}
	if err := odczytajJsonDostawcyDesignu(ctx, klient, adres, naglowki, &odpowiedz); err != nil {
		return nil, err
	}
	zasoby := make([]zasobDostawcyZdjecDesignu, 0, len(odpowiedz.Results))
	for _, wynik := range odpowiedz.Results {
		zasoby = append(zasoby, zasobDostawcyZdjecDesignu{
			Identyfikator: wynik.ID,
			Tytul:         pierwszyNiepustyDesignu(wynik.Description, wynik.AltDescription),
			Podglad:       wynik.Urls.Small, Pelny: wynik.Urls.Regular,
			Licencja: "Unsplash License", Autor: wynik.User.Name, Odsylacz: wynik.Links.HTML,
		})
	}
	return zasoby, nil
}

func pobierzUnsplashDesignu(ctx context.Context, klient *http.Client, klucz,
	identyfikator string) (zasobDostawcyZdjecDesignu, error) {

	var zdjecie struct {
		ID             string `json:"id"`
		Description    string `json:"description"`
		AltDescription string `json:"alt_description"`
		Urls           struct {
			Regular string `json:"regular"`
			Full    string `json:"full"`
		} `json:"urls"`
		Links struct {
			HTML string `json:"html"`
		} `json:"links"`
		User struct {
			Name string `json:"name"`
		} `json:"user"`
	}
	adres := "https://api.unsplash.com/photos/" + url.PathEscape(identyfikator)
	naglowki := map[string]string{"Authorization": "Client-ID " + klucz}
	if err := odczytajJsonDostawcyDesignu(ctx, klient, adres, naglowki, &zdjecie); err != nil {
		return zasobDostawcyZdjecDesignu{}, err
	}
	return zasobDostawcyZdjecDesignu{
		Identyfikator: zdjecie.ID,
		Tytul:         pierwszyNiepustyDesignu(zdjecie.Description, zdjecie.AltDescription),
		Pelny:         pierwszyNiepustyDesignu(zdjecie.Urls.Regular, zdjecie.Urls.Full),
		Licencja:      "Unsplash License", Autor: zdjecie.User.Name,
		Odsylacz: zdjecie.Links.HTML,
	}, nil
}

// ── Pexels ──────────────────────────────────────────────────────────────────

func szukajPexelsDesignu(ctx context.Context, klient *http.Client, klucz, fraza string,
	limit int) ([]zasobDostawcyZdjecDesignu, error) {

	var odpowiedz struct {
		Photos []struct {
			ID           int    `json:"id"`
			Alt          string `json:"alt"`
			URL          string `json:"url"`
			Photographer string `json:"photographer"`
			Src          struct {
				Large  string `json:"large"`
				Medium string `json:"medium"`
			} `json:"src"`
		} `json:"photos"`
	}
	adres := fmt.Sprintf("https://api.pexels.com/v1/search?query=%s&per_page=%d",
		url.QueryEscape(fraza), limit)
	naglowki := map[string]string{"Authorization": klucz}
	if err := odczytajJsonDostawcyDesignu(ctx, klient, adres, naglowki, &odpowiedz); err != nil {
		return nil, err
	}
	zasoby := make([]zasobDostawcyZdjecDesignu, 0, len(odpowiedz.Photos))
	for _, zdjecie := range odpowiedz.Photos {
		zasoby = append(zasoby, zasobDostawcyZdjecDesignu{
			Identyfikator: strconv.Itoa(zdjecie.ID), Tytul: zdjecie.Alt,
			Podglad: zdjecie.Src.Medium, Pelny: zdjecie.Src.Large,
			Licencja: "Pexels License", Autor: zdjecie.Photographer, Odsylacz: zdjecie.URL,
		})
	}
	return zasoby, nil
}

func pobierzPexelsDesignu(ctx context.Context, klient *http.Client, klucz,
	identyfikator string) (zasobDostawcyZdjecDesignu, error) {

	var zdjecie struct {
		ID           int    `json:"id"`
		Alt          string `json:"alt"`
		URL          string `json:"url"`
		Photographer string `json:"photographer"`
		Src          struct {
			Original string `json:"original"`
			Large    string `json:"large"`
		} `json:"src"`
	}
	adres := "https://api.pexels.com/v1/photos/" + url.PathEscape(identyfikator)
	naglowki := map[string]string{"Authorization": klucz}
	if err := odczytajJsonDostawcyDesignu(ctx, klient, adres, naglowki, &zdjecie); err != nil {
		return zasobDostawcyZdjecDesignu{}, err
	}
	return zasobDostawcyZdjecDesignu{
		Identyfikator: strconv.Itoa(zdjecie.ID), Tytul: zdjecie.Alt,
		Pelny:    pierwszyNiepustyDesignu(zdjecie.Src.Large, zdjecie.Src.Original),
		Licencja: "Pexels License", Autor: zdjecie.Photographer, Odsylacz: zdjecie.URL,
	}, nil
}

// ── Pixabay ─────────────────────────────────────────────────────────────────

func szukajPixabayDesignu(ctx context.Context, klient *http.Client, klucz, fraza string,
	limit int) ([]zasobDostawcyZdjecDesignu, error) {

	var odpowiedz struct {
		Hits []struct {
			ID            int    `json:"id"`
			Tags          string `json:"tags"`
			PageURL       string `json:"pageURL"`
			PreviewURL    string `json:"previewURL"`
			WebformatURL  string `json:"webformatURL"`
			LargeImageURL string `json:"largeImageURL"`
			User          string `json:"user"`
		} `json:"hits"`
	}
	adres := fmt.Sprintf("https://pixabay.com/api/?key=%s&q=%s&per_page=%d&image_type=photo",
		url.QueryEscape(klucz), url.QueryEscape(fraza), limit)
	if err := odczytajJsonDostawcyDesignu(ctx, klient, adres, nil, &odpowiedz); err != nil {
		return nil, err
	}
	zasoby := make([]zasobDostawcyZdjecDesignu, 0, len(odpowiedz.Hits))
	for _, wynik := range odpowiedz.Hits {
		zasoby = append(zasoby, zasobDostawcyZdjecDesignu{
			Identyfikator: strconv.Itoa(wynik.ID), Tytul: wynik.Tags,
			Podglad:  wynik.PreviewURL,
			Pelny:    pierwszyNiepustyDesignu(wynik.LargeImageURL, wynik.WebformatURL),
			Licencja: "Pixabay Content License", Autor: wynik.User, Odsylacz: wynik.PageURL,
		})
	}
	return zasoby, nil
}

func pobierzPixabayDesignu(ctx context.Context, klient *http.Client, klucz,
	identyfikator string) (zasobDostawcyZdjecDesignu, error) {

	var odpowiedz struct {
		Hits []struct {
			ID            int    `json:"id"`
			Tags          string `json:"tags"`
			PageURL       string `json:"pageURL"`
			WebformatURL  string `json:"webformatURL"`
			LargeImageURL string `json:"largeImageURL"`
			User          string `json:"user"`
		} `json:"hits"`
	}
	adres := fmt.Sprintf("https://pixabay.com/api/?key=%s&id=%s",
		url.QueryEscape(klucz), url.QueryEscape(identyfikator))
	if err := odczytajJsonDostawcyDesignu(ctx, klient, adres, nil, &odpowiedz); err != nil {
		return zasobDostawcyZdjecDesignu{}, err
	}
	if len(odpowiedz.Hits) == 0 {
		return zasobDostawcyZdjecDesignu{}, fmt.Errorf(
			"Pixabay nie zna zasobu %q", identyfikator)
	}
	wynik := odpowiedz.Hits[0]
	return zasobDostawcyZdjecDesignu{
		Identyfikator: strconv.Itoa(wynik.ID), Tytul: wynik.Tags,
		Pelny:    pierwszyNiepustyDesignu(wynik.LargeImageURL, wynik.WebformatURL),
		Licencja: "Pixabay Content License", Autor: wynik.User, Odsylacz: wynik.PageURL,
	}, nil
}

// Smithsonian Open Access — wymaga klucza api.data.gov, darmowego i wydawanego
// natychmiast.

func szukajSmithsonianDesignu(ctx context.Context, klient *http.Client, klucz, fraza string,
	limit int) ([]zasobDostawcyZdjecDesignu, error) {

	var odpowiedz struct {
		Response struct {
			Rows []smithsonianWierszDesignu `json:"rows"`
		} `json:"response"`
	}
	adres := fmt.Sprintf(
		"https://api.si.edu/openaccess/api/v1.0/search?api_key=%s&q=%s&rows=%d",
		url.QueryEscape(klucz),
		url.QueryEscape(fraza+" AND online_media_type:\"Images\""), limit)
	if err := odczytajJsonDostawcyDesignu(ctx, klient, adres, nil, &odpowiedz); err != nil {
		return nil, err
	}
	zasoby := make([]zasobDostawcyZdjecDesignu, 0, len(odpowiedz.Response.Rows))
	for _, wiersz := range odpowiedz.Response.Rows {
		if zasob, jest := zasobZeSmithsonianDesignu(wiersz); jest {
			zasoby = append(zasoby, zasob)
		}
	}
	// Wiersze przyszły, ale ani jeden nie dał się złożyć w zasób — odmowa
	// nazwana zamiast wykazu pustego.
	if len(zasoby) == 0 && len(odpowiedz.Response.Rows) > 0 {
		return nil, fmt.Errorf(
			"Smithsonian oddał %d wierszy, z których żaden nie niesie obrazu pod "+
				"content.descriptiveNonRepeating.online_media.media — kształt odpowiedzi nie "+
				"zgadza się z tym, czego serwer w niej szuka", len(odpowiedz.Response.Rows))
	}
	return zasoby, nil
}

// smithsonianWierszDesignu to jeden wiersz odpowiedzi Smithsonian; media leżą
// w polu online_media, zagnieżdżonym pod polem descriptiveNonRepeating.
type smithsonianWierszDesignu struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content struct {
		Descriptive struct {
			Names []struct {
				Content string `json:"content"`
			} `json:"name"`
		} `json:"freetext"`
		Opis struct {
			Media struct {
				Wykaz []struct {
					Type      string `json:"type"`
					Content   string `json:"content"`
					Thumbnail string `json:"thumbnail"`
					Usage     struct {
						Access string `json:"access"`
					} `json:"usage"`
				} `json:"media"`
			} `json:"online_media"`
		} `json:"descriptiveNonRepeating"`
	} `json:"content"`
}

// zasobZeSmithsonianDesignu składa zasób z wiersza. Wiersz bez obrazu nie wchodzi
// do wykazu — Smithsonian ma w zbiorach także same opisy.
func zasobZeSmithsonianDesignu(wiersz smithsonianWierszDesignu) (zasobDostawcyZdjecDesignu, bool) {
	for _, media := range wiersz.Content.Opis.Media.Wykaz {
		if !strings.EqualFold(media.Type, "Images") || media.Content == "" {
			continue
		}
		// Wartość CC0 pola usage.access znaczy domenę publiczną.
		licencja := "Smithsonian — użycie ograniczone (" + media.Usage.Access + ")"
		if strings.EqualFold(media.Usage.Access, "CC0") {
			licencja = "CC0 1.0 (domena publiczna, Smithsonian Open Access)"
		}
		autor := ""
		if len(wiersz.Content.Descriptive.Names) > 0 {
			autor = wiersz.Content.Descriptive.Names[0].Content
		}
		return zasobDostawcyZdjecDesignu{
			Identyfikator: wiersz.ID, Tytul: wiersz.Title,
			Podglad: media.Thumbnail, Pelny: media.Content,
			Licencja: licencja, Autor: autor,
			Odsylacz: "https://www.si.edu/object/" + wiersz.ID,
		}, true
	}
	return zasobDostawcyZdjecDesignu{}, false
}

func pobierzSmithsonianDesignu(ctx context.Context, klient *http.Client, klucz,
	identyfikator string) (zasobDostawcyZdjecDesignu, error) {

	var odpowiedz struct {
		Response smithsonianWierszDesignu `json:"response"`
	}
	adres := fmt.Sprintf("https://api.si.edu/openaccess/api/v1.0/content/%s?api_key=%s",
		url.PathEscape(identyfikator), url.QueryEscape(klucz))
	if err := odczytajJsonDostawcyDesignu(ctx, klient, adres, nil, &odpowiedz); err != nil {
		return zasobDostawcyZdjecDesignu{}, err
	}
	zasob, jest := zasobZeSmithsonianDesignu(odpowiedz.Response)
	if !jest {
		return zasobDostawcyZdjecDesignu{}, fmt.Errorf(
			"zasób %q nie ma u Smithsonian pliku obrazu", identyfikator)
	}
	return zasob, nil
}

// pierwszyNiepustyDesignu oddaje pierwszą niepustą wartość z podanych, po
// obcięciu białych znaków, albo pusty napis, gdy wszystkie są puste.
func pierwszyNiepustyDesignu(wartosci ...string) string {
	for _, wartosc := range wartosci {
		if strings.TrimSpace(wartosc) != "" {
			return strings.TrimSpace(wartosc)
		}
	}
	return ""
}
