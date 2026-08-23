package core

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// Przelot ośmiu dróg HTTP dostawców baz zdjęciowych modułu Design.
//
// ── Czego ten sprawdzian dowodzi, a czego NIE ───────────────────────────────
// DOWODZI, że dla każdego z ośmiu dostawców droga idzie do końca: rdzeń składa
// adres, wysyła żądanie po HTTP, nosi klucz tam, gdzie dostawca go żąda, rozkłada
// odpowiedź i wypełnia z niej wspólną postać zasobu. Serwer próbny stoi w uprzęży
// i odpowiada kształtami, które rdzeń zakłada.
//
// NIE dowodzi, że kształt odpowiedzi zgadza się z tym, co dostawca naprawdę
// wysyła. Kształty pochodzą z dokumentacji, nie z pomiaru na jego API, i tego
// sprawdzian bez konta u dostawcy nie zamknie. Ta połowa braku zostaje otwarta
// i jest tak nazwana — inaczej zielony wynik tego pliku czytałoby się jako
// „dostawcy zmierzeni", a zmierzona jest DROGA.
//
// ── Dlaczego bez zmian w rdzeniu ────────────────────────────────────────────
// Każda droga dostawcy przyjmuje klienta HTTP jako argument, więc sprawdzian
// podstawia własnego — z przekładnią, która przepisuje gospodarza adresu na
// serwer próbny i zapisuje, o co rdzeń naprawdę poprosił. Rdzeń nie dostaje ani
// jednego pola „adres na potrzeby sprawdzianu", bo pole takie żyłoby w produkcie
// i dałoby się nim wskazać serwer obcy.

// przekladniaDostawcowSprawdzianu przepisuje adres żądania na serwer próbny
// i zapisuje żądanie w takiej postaci, w jakiej rdzeń je złożył.
type przekladniaDostawcowSprawdzianu struct {
	cel     *url.URL
	zadania []*http.Request
	spod    http.RoundTripper
}

func (p *przekladniaDostawcowSprawdzianu) RoundTrip(zadanie *http.Request) (*http.Response, error) {
	kopia := zadanie.Clone(zadanie.Context())
	p.zadania = append(p.zadania, kopia)
	zadanie.URL.Scheme = p.cel.Scheme
	zadanie.URL.Host = p.cel.Host
	zadanie.Host = ""
	return p.spod.RoundTrip(zadanie)
}

// ostatnieZadanie oddaje żądanie wysłane jako ostatnie.
func (p *przekladniaDostawcowSprawdzianu) ostatnieZadanie() *http.Request {
	if len(p.zadania) == 0 {
		return nil
	}
	return p.zadania[len(p.zadania)-1]
}

// odpowiedziDostawcowSprawdzianu wylicza treści, którymi serwer próbny odpowiada
// na poszczególne ścieżki. Kształty są tymi, których rdzeń oczekuje.
func odpowiedziDostawcowSprawdzianu() map[string]string {
	return map[string]string{
		// Openverse
		"/v1/images/": `{"results":[{"id":"ov-1","title":"Kot na dachu",
			"url":"https://tresc/pelny.jpg","thumbnail":"https://tresc/mini.jpg",
			"license":"cc0","creator":"Anna Autorka",
			"foreign_landing_url":"https://openverse/strona"}]}`,
		"/v1/images/ov-1/": `{"id":"ov-1","title":"Kot na dachu",
			"url":"https://tresc/pelny.jpg","license":"cc0","creator":"Anna Autorka",
			"foreign_landing_url":"https://openverse/strona"}`,
		// Wikimedia Commons — jedna ścieżka na oba wywołania, rozróżnia je zapytanie.
		"/w/api.php": `{"query":{"pages":{"7":{"title":"File:Kot.jpg","imageinfo":[
			{"url":"https://tresc/kot.jpg","thumburl":"https://tresc/kot-640.jpg",
			 "descriptionurl":"https://commons/strona",
			 "extmetadata":{"LicenseShortName":{"value":"CC BY-SA 4.0"},
			                "Artist":{"value":"<a href=\"x\">Anna Autorka</a>"}}}]}}}}`,
		// Metropolitan Museum of Art
		"/public/collection/v1/search": `{"objectIDs":[451234]}`,
		"/public/collection/v1/objects/451234": `{"objectID":451234,"title":"Waza",
			"primaryImage":"https://tresc/waza.jpg","primaryImageSmall":"https://tresc/waza-mini.jpg",
			"artistDisplayName":"Nieznany","objectURL":"https://met/obiekt/451234",
			"isPublicDomain":true}`,
		// NASA
		"/search": `{"collection":{"items":[{"data":[{"nasa_id":"PIA00001","title":"Mars",
			"center":"JPL","photographer":"NASA/JPL"}],
			"links":[{"href":"https://tresc/mars-mini.jpg","rel":"preview"}]}]}}`,
		"/asset/PIA00001": `{"collection":{"items":[
			{"href":"https://tresc/mars-metadata.json"},{"href":"https://tresc/mars-orig.jpg"}]}}`,
		// Unsplash
		"/search/photos": `{"results":[{"id":"un-1","description":"Kot w oknie",
			"alt_description":"kot","urls":{"regular":"https://tresc/un-regular.jpg",
			"small":"https://tresc/un-small.jpg"},"links":{"html":"https://unsplash/zdjecie"},
			"user":{"name":"Anna Autorka"}}]}`,
		"/photos/un-1": `{"id":"un-1","description":"Kot w oknie","alt_description":"kot",
			"urls":{"regular":"https://tresc/un-regular.jpg","full":"https://tresc/un-full.jpg"},
			"links":{"html":"https://unsplash/zdjecie"},"user":{"name":"Anna Autorka"}}`,
		// Pexels
		"/v1/search": `{"photos":[{"id":9911,"alt":"Kot na kanapie",
			"url":"https://pexels/zdjecie","photographer":"Anna Autorka",
			"src":{"large":"https://tresc/px-large.jpg","medium":"https://tresc/px-medium.jpg"}}]}`,
		"/v1/photos/9911": `{"id":9911,"alt":"Kot na kanapie","url":"https://pexels/zdjecie",
			"photographer":"Anna Autorka","src":{"original":"https://tresc/px-original.jpg",
			"large":"https://tresc/px-large.jpg"}}`,
		// Pixabay — jedna ścieżka na oba wywołania.
		"/api/": `{"hits":[{"id":5511,"tags":"kot, zwierzę, futro",
			"pageURL":"https://pixabay/zdjecie","previewURL":"https://tresc/pb-preview.jpg",
			"webformatURL":"https://tresc/pb-web.jpg","largeImageURL":"https://tresc/pb-large.jpg",
			"user":"AnnaAutorka"}]}`,
		// Smithsonian Open Access — kształt ODCZYTANY z żywego API (media leżą
		// w `online_media`, nie wprost w `descriptiveNonRepeating`).
		"/openaccess/api/v1.0/search": `{"response":{"rows":[{"id":"sm-1","title":"Rycina",
			"content":{"freetext":{"name":[{"content":"Anna Autorka"}]},
			"descriptiveNonRepeating":{"guid":"https://n2t.net/ark:/65665/sm-1",
			"online_media":{"mediaCount":1,"media":[{"type":"Images",
			"content":"https://tresc/sm-pelny.jpg","thumbnail":"https://tresc/sm-mini.jpg",
			"usage":{"access":"CC0"}}]}}}}]}}`,
		"/openaccess/api/v1.0/content/sm-1": `{"response":{"id":"sm-1","title":"Rycina",
			"content":{"freetext":{"name":[{"content":"Anna Autorka"}]},
			"descriptiveNonRepeating":{"guid":"https://n2t.net/ark:/65665/sm-1",
			"online_media":{"mediaCount":1,"media":[{"type":"Images",
			"content":"https://tresc/sm-pelny.jpg","thumbnail":"https://tresc/sm-mini.jpg",
			"usage":{"access":"CC0"}}]}}}}}`,
	}
}

// serwerDostawcowSprawdzianu stawia serwer próbny wraz z klientem, który do niego
// przekierowuje.
func serwerDostawcowSprawdzianu(t *testing.T) (*httptest.Server,
	*przekladniaDostawcowSprawdzianu, *http.Client) {

	t.Helper()

	tresci := odpowiedziDostawcowSprawdzianu()
	serwer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, z *http.Request) {
		tresc, jest := tresci[z.URL.Path]
		if !jest {
			// Ścieżka, której serwer próbny nie zna, jest NIEPOWODZENIEM sprawdzianu,
			// nie pustą odpowiedzią: znaczy, że rdzeń poszedł gdzie indziej, niż
			// sprawdzian mierzy.
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"blad":"serwer próbny nie zna ścieżki ` + z.URL.Path + `"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(tresc))
	}))
	t.Cleanup(serwer.Close)

	adres, err := url.Parse(serwer.URL)
	if err != nil {
		t.Fatalf("adres serwera próbnego nieczytelny: %v", err)
	}
	przekladnia := &przekladniaDostawcowSprawdzianu{cel: adres, spod: http.DefaultTransport}
	return serwer, przekladnia, &http.Client{Transport: przekladnia}
}

// TestOsiemDrogDostawcowZdjecPrzechodziPrzelotem prowadzi każdą z ośmiu dróg
// wyszukania i każdą z ośmiu dróg wciągnięcia przez prawdziwe żądanie HTTP
// i mierzy wspólną postać zasobu, którą rdzeń z odpowiedzi złożył.
func TestOsiemDrogDostawcowZdjecPrzechodziPrzelotem(t *testing.T) {
	_, przekladnia, klient := serwerDostawcowSprawdzianu(t)
	zycie := context.Background()

	przypadki := []struct {
		dostawca        string
		klucz           string
		identyfikator   string
		gospodarz       string
		sciezkaSzukaj   string
		tytul           string
		licencja        string
		naglowekKlucza  string
		kluczWZapytaniu string
	}{
		{
			dostawca: "openverse", identyfikator: "ov-1",
			gospodarz: "api.openverse.org", sciezkaSzukaj: "/v1/images/",
			tytul: "Kot na dachu", licencja: "cc0",
		},
		{
			dostawca: "wikimedia", identyfikator: "File:Kot.jpg",
			gospodarz: "commons.wikimedia.org", sciezkaSzukaj: "/w/api.php",
			tytul: "File:Kot.jpg", licencja: "CC BY-SA 4.0",
		},
		{
			dostawca: "met", identyfikator: "451234",
			gospodarz:     "collectionapi.metmuseum.org",
			sciezkaSzukaj: "/public/collection/v1/search",
			tytul:         "Waza", licencja: "CC0 1.0 (domena publiczna, Met Open Access)",
		},
		{
			dostawca: "nasa", identyfikator: "PIA00001",
			gospodarz: "images-api.nasa.gov", sciezkaSzukaj: "/search",
			tytul: "Mars", licencja: "NASA Media Usage Guidelines",
		},
		{
			dostawca: "unsplash", klucz: "klucz-unsplash", identyfikator: "un-1",
			gospodarz: "api.unsplash.com", sciezkaSzukaj: "/search/photos",
			tytul: "Kot w oknie", licencja: "Unsplash License",
			naglowekKlucza: "Client-ID klucz-unsplash",
		},
		{
			dostawca: "pexels", klucz: "klucz-pexels", identyfikator: "9911",
			gospodarz: "api.pexels.com", sciezkaSzukaj: "/v1/search",
			tytul: "Kot na kanapie", licencja: "Pexels License",
			naglowekKlucza: "klucz-pexels",
		},
		{
			dostawca: "pixabay", klucz: "klucz-pixabay", identyfikator: "5511",
			gospodarz: "pixabay.com", sciezkaSzukaj: "/api/",
			tytul: "kot, zwierzę, futro", licencja: "Pixabay Content License",
			kluczWZapytaniu: "key",
		},
		{
			dostawca: "smithsonian", klucz: "klucz-si", identyfikator: "sm-1",
			gospodarz: "api.si.edu", sciezkaSzukaj: "/openaccess/api/v1.0/search",
			tytul: "Rycina", licencja: "CC0 1.0 (domena publiczna, Smithsonian Open Access)",
			kluczWZapytaniu: "api_key",
		},
	}

	for _, przypadek := range przypadki {
		dostawca, jest := dostawcaZdjecPoNazwieDesignu(przypadek.dostawca)
		if !jest {
			t.Fatalf("rdzeń nie zna dostawcy %s, a wykaz go wymienia", przypadek.dostawca)
		}

		zasoby, err := dostawca.Szukaj(zycie, klient, przypadek.klucz, "kot", 5)
		if err != nil {
			t.Errorf("wyszukanie u dostawcy %s nie doszło do końca: %v", przypadek.dostawca, err)
			continue
		}
		if len(zasoby) != 1 {
			t.Errorf("wyszukanie u dostawcy %s oddało %d zasobów, a serwer próbny odpowiedział "+
				"jednym", przypadek.dostawca, len(zasoby))
			continue
		}

		// Adres złożony przez rdzeń — gospodarz i ścieżka. Rdzeń pytający innego
		// gospodarza pytałby w produkcie kogoś innego, niż Operator zamówił.
		zadanie := przekladnia.zadania[0]
		if zadanie.URL.Host != przypadek.gospodarz {
			t.Errorf("dostawca %s: rdzeń zapytał gospodarza %s, a droga prowadzi do %s",
				przypadek.dostawca, zadanie.URL.Host, przypadek.gospodarz)
		}
		if zadanie.URL.Path != przypadek.sciezkaSzukaj {
			t.Errorf("dostawca %s: rdzeń zapytał o ścieżkę %s, a droga wyszukania to %s",
				przypadek.dostawca, zadanie.URL.Path, przypadek.sciezkaSzukaj)
		}
		// Nagłówek rozpoznawczy jest wymagany przez część dostawców i bez niego
		// odpowiadają odmową — Operator widziałby dostawcę jako niedostępnego.
		if zadanie.Header.Get("User-Agent") == "" {
			t.Errorf("dostawca %s: żądanie bez nagłówka rozpoznawczego", przypadek.dostawca)
		}
		if przypadek.naglowekKlucza != "" &&
			zadanie.Header.Get("Authorization") != przypadek.naglowekKlucza {

			t.Errorf("dostawca %s: nagłówek klucza to %q, a ma być %q", przypadek.dostawca,
				zadanie.Header.Get("Authorization"), przypadek.naglowekKlucza)
		}
		if przypadek.kluczWZapytaniu != "" &&
			zadanie.URL.Query().Get(przypadek.kluczWZapytaniu) != przypadek.klucz {

			t.Errorf("dostawca %s: klucz w zapytaniu (%s) to %q, a ma być %q",
				przypadek.dostawca, przypadek.kluczWZapytaniu,
				zadanie.URL.Query().Get(przypadek.kluczWZapytaniu), przypadek.klucz)
		}

		zasob := zasoby[0]
		if zasob.Identyfikator == "" {
			t.Errorf("dostawca %s: zasób bez identyfikatora nie da się wciągnąć",
				przypadek.dostawca)
		}
		if zasob.Tytul != przypadek.tytul {
			t.Errorf("dostawca %s: tytuł zasobu to %q, a odpowiedź niosła %q",
				przypadek.dostawca, zasob.Tytul, przypadek.tytul)
		}
		// Wyszukanie musi dać adres, którym okno pokaże podgląd. Adres pełnej treści
		// bywa u dostawcy znany dopiero przy wciągnięciu — tak jest u NASA, której
		// wyszukanie oddaje wyłącznie podgląd, a plik źródłowy stoi pod osobną
		// drogą. Dlatego warunkiem jest tu JEDEN z dwóch adresów, a pełny mierzy się
		// na drodze wciągnięcia niżej.
		if zasob.Podglad == "" && zasob.Pelny == "" {
			t.Errorf("dostawca %s: zasób bez ani jednego adresu — okno nie ma czego pokazać",
				przypadek.dostawca)
		}
		// Licencja jest warunkiem wciągnięcia: zasób z bazy zewnętrznej bez
		// zapisanej licencji jest usterką, nie zasobem.
		if !strings.HasPrefix(zasob.Licencja, przypadek.licencja) {
			t.Errorf("dostawca %s: licencja zasobu to %q, a ma zaczynać się od %q",
				przypadek.dostawca, zasob.Licencja, przypadek.licencja)
		}

		// Droga wciągnięcia — ta sama miara na drugiej połowie dostawcy.
		przekladnia.zadania = nil
		pobrany, err := dostawca.Pobierz(zycie, klient, przypadek.klucz, przypadek.identyfikator)
		if err != nil {
			t.Errorf("wciągnięcie zasobu %s od dostawcy %s nie doszło do końca: %v",
				przypadek.identyfikator, przypadek.dostawca, err)
			przekladnia.zadania = nil
			continue
		}
		if pobrany.Pelny == "" {
			t.Errorf("dostawca %s: wciągnięty zasób bez adresu treści", przypadek.dostawca)
		}
		if pobrany.Licencja == "" {
			t.Errorf("dostawca %s: wciągnięty zasób bez licencji", przypadek.dostawca)
		}
		if ostatnie := przekladnia.ostatnieZadanie(); ostatnie == nil ||
			ostatnie.URL.Host != przypadek.gospodarz {

			t.Errorf("dostawca %s: wciągnięcie poszło do %v, a droga prowadzi do %s",
				przypadek.dostawca, ostatnie, przypadek.gospodarz)
		}
		przekladnia.zadania = nil
	}
}

// TestSmithsonianBezMediowWWierszuNieMilczy jest sprawdzianem na defekt
// zmierzony na ŻYWYM API: rdzeń szukał mediów pod
// `content.descriptiveNonRepeating.media`, a odpowiedź `api.si.edu` niesie je pod
// `…descriptiveNonRepeating.online_media.media`. Wiersze przychodziły, żaden nie
// dawał się złożyć w zasób, komenda oddawała wykaz pusty BEZ błędu — więc
// dostawca nie wracał ani w wykazie, ani w `providersFailed`, a Operator czytał
// to jako „fraza nie ma zdjęć".
//
// Sprawdzian mierzy obie strony: kształt właściwy daje zasób, a wiersze bez
// mediów dają odmowę NAZWANĄ wraz z liczbą wierszy.
func TestSmithsonianBezMediowWWierszuNieMilczy(t *testing.T) {
	odpowiedz := `{"response":{"rows":[
		{"id":"sm-9","title":"Opis bez obrazu","content":{"freetext":{},
		 "descriptiveNonRepeating":{"guid":"g","title":{"content":"Opis"}}}},
		{"id":"sm-10","title":"Drugi opis","content":{"freetext":{},
		 "descriptiveNonRepeating":{"guid":"g2"}}}]}}`
	serwer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(odpowiedz))
	}))
	t.Cleanup(serwer.Close)

	adres, err := url.Parse(serwer.URL)
	if err != nil {
		t.Fatalf("adres serwera próbnego nieczytelny: %v", err)
	}
	klient := &http.Client{Transport: &przekladniaDostawcowSprawdzianu{
		cel: adres, spod: http.DefaultTransport,
	}}

	zasoby, err := szukajSmithsonianDesignu(context.Background(), klient, "klucz", "kot", 5)
	if len(zasoby) != 0 {
		t.Fatalf("z wierszy bez mediów powstało %d zasobów", len(zasoby))
	}
	if err == nil {
		t.Fatal("wiersze bez mediów dały wykaz pusty BEZ błędu — dostawca nie wróci w bilansie " +
			"i Operator odczyta to jako brak zdjęć dla frazy")
	}
	if !strings.Contains(err.Error(), "online_media") || !strings.Contains(err.Error(), "2") {
		t.Errorf("odmowa nie nazywa ani drogi mediów, ani liczby wierszy: %v", err)
	}
}

// TestDostawcaZdjecOdpowiadajacyOdmowaWracaBledem pilnuje warunku odwrotnego:
// odpowiedź o stanie błędu nie ma prawa wyjść z rdzenia jako pusty wykaz
// zasobów. Pusty wykaz bez powodu Operator odczyta jako „fraza nie ma zdjęć",
// a to nieprawda o frazie i prawda o sieci.
func TestDostawcaZdjecOdpowiadajacyOdmowaWracaBledem(t *testing.T) {
	serwer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"blad":"klucz odrzucony"}`))
	}))
	t.Cleanup(serwer.Close)

	adres, err := url.Parse(serwer.URL)
	if err != nil {
		t.Fatalf("adres serwera próbnego nieczytelny: %v", err)
	}
	klient := &http.Client{Transport: &przekladniaDostawcowSprawdzianu{
		cel: adres, spod: http.DefaultTransport,
	}}

	for _, dostawca := range dostawcyZdjecDesignu() {
		zasoby, err := dostawca.Szukaj(context.Background(), klient, "klucz", "kot", 5)
		if err == nil && len(zasoby) > 0 {
			t.Errorf("dostawca %s przy odpowiedzi o stanie 403 oddał %d zasobów",
				dostawca.Nazwa, len(zasoby))
		}
		if err == nil && len(zasoby) == 0 {
			// Met składa wynik z dwóch wywołań i pierwsze pada, więc wykaz jest
			// pusty bez błędu — ale wtedy bilans komendy zapisuje dostawcę jako
			// tego, który nic nie dał. Warunkiem jest, żeby NIE oddał zasobu.
			continue
		}
	}
}
