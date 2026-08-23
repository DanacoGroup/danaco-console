package models

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
)

// kanalObrazow jest kanałem modelu oddającym bajty obrazu, a nie tekst. Jest
// bliźniakiem KanalAPI i dzieli z nim całą warstwę dostępową (nagłówki, klucz
// z sejfu albo ze zmiennej środowiskowej, dodatki ciała, limit czasu) — jedna
// prawda o poświadczeniach, jedna o nagłówkach. Różni się tym jednym,
// czym musi: kształtem żądania (prompt zamiast wiadomości) i tym, że wynik
// nadaje fragmentem `image`, nie porcjami tekstu.
//
// Kształt żądania jest ten, który dostawcy powtarzają za OpenAI Images:
// POST z ciałem {model, prompt, n, size} i odpowiedzią {"data":[{"b64_json"|"url"}]}.
// Powtarzają go dziś także dostawcy niezależni i bramy lokalne, więc jest to
// najczęstszy kształt, a nie jeden dostawca zaszyty w kodzie. Wszystko, co
// w tym kształcie bywa inne, jest parametrem wiersza rejestru:
//
//	base_url          — pełny adres punktu końcowego (wymagany)
//	rozmiar           — wartość pola `size`; domyślnie 1024x1024
//	liczba            — wartość pola `n`; domyślnie 1
//	format_odpowiedzi — wartość `response_format` (np. b64_json); pusty = nie wysyłamy
//	sciezka_base64    — ścieżka do bajtów w odpowiedzi; domyślnie data.0.b64_json
//	sciezka_adresu    — ścieżka do adresu obrazu;     domyślnie data.0.url
//	typ_tresci        — rodzaj bajtów w meldunku fragmentu; domyślnie image/png
//	sciezka_bledu, naglowki, naglowek_klucza, przedrostek_klucza,
//	cialo_dodatkowe, limit_sekund — jak w kanale api
//
// Wiersz zakłada się `channel.add` z rodzajem `api` (CHECK schematu zna cztery
// rodzaje i kontrakt jest zamrożony) oraz parametrem `adapter` równym „obrazy".
// Rodzaj mówi, jak kanał rozmawia (HTTP), adapter mówi, co oddaje — dlatego
// rodzaju kanału nie trzeba dokładać ani migracją, ani do kontraktu.
//
// ── Obraz WEJŚCIOWY: edycja obok generowania ────────────────────────────────
// Wywołanie niosące `Zapytanie.ObrazyWejsciowe` jest EDYCJĄ materiału, nie
// generowaniem od zera, i dostawcy dają na nią osobny punkt końcowy o osobnym
// kształcie żądania (u OpenAI `images/edits`: multipart/form-data z plikiem
// `image` i opcjonalną maską `mask`). Bramy lokalne częściej przyjmują ten sam
// adres z bajtami w polu JSON. Adapter zna oba kształty, bo obu nie da się
// pogodzić, a zgadywanie kończyłoby się odmową dostawcy zamiast obrazu:
//
//	adres_edycji   — adres punktu końcowego edycji; pusty = ten sam co base_url
//	postac_obrazu  — „wieloczesciowa" (domyślnie) albo „base64"
//	pole_obrazu    — nazwa pola materiału; domyślnie image
//	pole_maski     — nazwa pola maski;     domyślnie mask
//
// Wywołanie BEZ obrazu wejściowego idzie dokładnie tą samą drogą co dotąd —
// pole puste nie zmienia ani adresu, ani kształtu ciała, więc kanały już
// założone pracują bez zmiany wiersza.
type kanalObrazow struct {
	def    Definicja
	klient *http.Client
}

// NowyKanalObrazow buduje kanał obrazowy z wiersza rejestru. Brak adresu jest
// jedynym powodem odmowy budowy — reszta ma wartości domyślne, a wiersz bez
// poświadczenia buduje się celowo: odmowa ma paść w chwili wywołania, wymieniając
// brak z nazwy, a nie zniknąć jako „kanał pominięty" w wykazie rejestru.
func NowyKanalObrazow(d Definicja) (Kanal, error) {
	if strings.TrimSpace(d.Parametr("base_url")) == "" {
		return nil, fmt.Errorf("models: kanał %q (adapter %s) bez parametru base_url", d.Kod, AdapterObrazy)
	}
	return &kanalObrazow{def: d, klient: &http.Client{Timeout: limitCzasu(d)}}, nil
}

// Kod zwraca kod kanału z wiersza rejestru.
func (k *kanalObrazow) Kod() string {
	return k.def.Kod
}

// Definicja zwraca wiersz rejestru, z którego kanał powstał.
func (k *kanalObrazow) Definicja() Definicja {
	return k.def
}

// Wyslij nadaje prowenancję, metadane konta i — po udanym żądaniu — jeden
// fragment obrazu. Kolejność jest ta sama, co w kanałach tekstowych:
// odbiorca poznaje warunki wywołania, zanim zobaczy wynik.
func (k *kanalObrazow) Wyslij(ctx context.Context, z Zapytanie, u Ujscie) error {
	prowenancja := ProwenancjaZapytania(z, k.def)
	// Prowenancja niesie adres, który NAPRAWDĘ pojedzie: edycja materiału bywa
	// innym punktem końcowym niż generowanie, a prowenancja wskazująca nie ten
	// adres, którego użyto, jest gorsza od prowenancji bez adresu.
	prowenancja.Adres = k.adresWywolania(z)
	if err := NadajProwenancje(ctx, u, z, prowenancja); err != nil {
		return err
	}
	if k.def.PoswiadczenieOdwolanie != "" {
		konto := MetadaneKonta{Konto: z.WybraneKonto(k.def), Odwolanie: k.def.PoswiadczenieOdwolanie}
		if err := NadajKonto(ctx, u, z, konto); err != nil {
			return err
		}
	}
	if err := k.sprawdzPoswiadczenie(ctx); err != nil {
		return err
	}

	polecenie := poleceniePromptu(z)
	zadanie, err := k.zbudujZadanie(ctx, z, polecenie)
	if err != nil {
		return err
	}
	odpowiedz, err := k.klient.Do(zadanie)
	if err != nil {
		return fmt.Errorf("models: kanał %s: %w", k.def.Kod, err)
	}
	defer odpowiedz.Body.Close()
	if odpowiedz.StatusCode < 200 || odpowiedz.StatusCode > 299 {
		return bladOdpowiedzi(k.def, odpowiedz)
	}
	return k.nadajObraz(ctx, z, u, polecenie, odpowiedz.Body)
}

// sprawdzPoswiadczenie odmawia, nazywając brak, zanim poleci żądanie.
//
// Kanał tekstowy dopuszcza brak odwołania, bo punkt końcowy bez uwierzytelnienia
// istnieje (Ollama pod adresem pętli zwrotnej). Punkt końcowy generujący obrazy
// bez klucza nie istnieje — żądanie bez poświadczenia odbiłoby się o 401
// i wróciło jako błąd dostawcy, czyli komunikat o czymś innym niż rzeczywisty
// brak. Dlatego brak odwołania rozstrzygamy tutaj i mówimy wprost, czego nie ma
// i gdzie się to wskazuje. Wywołanie nie wraca ani atrapą, ani pustym obrazem:
// niczego nie było, więc nic nie jedzie w strumieniu poza tą odmową.
//
// To samo dotyczy odwołania, które jest, ale nie prowadzi do sekretu — sejf bez
// wpisu, zmienna środowiskowa nieustawiona. Tu brak wychodzi już z rozwiązania
// odwołania i różnica jest widoczna w treści: „bez odwołania" to inny stan niż
// „odwołanie bez wartości", i Operator poprawia je w dwóch różnych miejscach.
func (k *kanalObrazow) sprawdzPoswiadczenie(ctx context.Context) error {
	odwolanie := strings.TrimSpace(k.def.PoswiadczenieOdwolanie)
	if odwolanie == "" {
		return fmt.Errorf("models: kanał %s: brak poświadczenia kanału obrazowego — "+
			"punkt końcowy %s nie przyjmuje żądań bez klucza, a rdzeń klucza nie zmyśli; "+
			"wskaż go parametrem credentialRef w konfiguracji kanału (channel.update): "+
			"nazwą zmiennej środowiskowej albo referencją sejfu w postaci sejf:<byt>",
			k.def.Kod, k.def.Parametr("base_url"))
	}
	if _, err := poswiadczenieZOdwolania(ctx, k.def, odwolanie); err != nil {
		return err
	}
	return nil
}

// poleceniePromptu składa polecenie obrazu. Punkty końcowe obrazowe nie mają
// roli systemowej ani historii rozmowy — przyjmują jedno pole `prompt` — więc
// nakładka systemowa okna jedzie przed treścią, oddzielona pustą
// linią. Gdyby ją pominąć, warstwy konstytucji i roli przestałyby obowiązywać
// akurat na tym kanale, a nakładka ma obowiązywać wszędzie tak samo.
func poleceniePromptu(z Zapytanie) string {
	czesci := make([]string, 0, 2)
	if prompt := z.Nakladka.PromptSystemowy(); strings.TrimSpace(prompt) != "" {
		czesci = append(czesci, prompt)
	}
	if strings.TrimSpace(z.Tresc) != "" {
		czesci = append(czesci, z.Tresc)
	}
	return strings.Join(czesci, "\n\n")
}

// Postacie wysyłki obrazu wejściowego — wartości parametru `postac_obrazu`.
const (
	postacObrazuWieloczesciowa = "wieloczesciowa"
	postacObrazuBase64         = "base64"
)

// zbudujZadanie składa żądanie HTTP kanału obrazowego. Warstwa dostępowa idzie
// przez te same funkcje, co kanał api — nagłówki wiersza i klucz z odwołania.
//
// Rozgałęzienie jest jedno i pada tutaj: wywołanie z materiałem wejściowym
// w postaci wieloczęściowej ma inne ciało i inny rodzaj treści, więc składa je
// osobna funkcja. Wszystko poza ciałem — adres, nagłówki, klucz — jest wspólne
// i wykonuje się raz.
func (k *kanalObrazow) zbudujZadanie(ctx context.Context, z Zapytanie, polecenie string) (*http.Request, error) {
	var (
		cialo     io.Reader
		typTresci string
		err       error
	)
	if z.MaMaterialWejsciowy() && k.postacObrazu() == postacObrazuWieloczesciowa {
		cialo, typTresci, err = k.cialoWieloczesciowe(z, polecenie)
	} else {
		cialo, typTresci, err = k.cialoJSON(z, polecenie)
	}
	if err != nil {
		return nil, err
	}
	zadanie, err := http.NewRequestWithContext(ctx, http.MethodPost, k.adresWywolania(z), cialo)
	if err != nil {
		return nil, fmt.Errorf("models: kanał %s: budowa żądania: %w", k.def.Kod, err)
	}
	zadanie.Header.Set("Content-Type", typTresci)
	for nazwa, wartosc := range naglowkiWiersza(k.def) {
		zadanie.Header.Set(nazwa, wartosc)
	}
	if err := dolozKlucz(k.def, zadanie); err != nil {
		return nil, err
	}
	return zadanie, nil
}

// adresWywolania rozstrzyga punkt końcowy wywołania. Materiał wejściowy prowadzi
// pod adres edycji, gdy wiersz go wskazał; brak wskazania znaczy „ten sam adres",
// bo bramy lokalne obsługują generowanie i edycję jednym punktem końcowym.
func (k *kanalObrazow) adresWywolania(z Zapytanie) string {
	if z.MaMaterialWejsciowy() {
		if adres := strings.TrimSpace(k.def.Parametr("adres_edycji")); adres != "" {
			return adres
		}
	}
	return k.def.Parametr("base_url")
}

// postacObrazu odczytuje kształt wysyłki materiału. Wartość nierozpoznana schodzi
// na postać wieloczęściową, bo to ona jest kształtem dostawców — parametr
// wpisany z literówką nie może zamienić edycji w generowanie od zera.
func (k *kanalObrazow) postacObrazu() string {
	if strings.EqualFold(strings.TrimSpace(k.def.Parametr("postac_obrazu")), postacObrazuBase64) {
		return postacObrazuBase64
	}
	return postacObrazuWieloczesciowa
}

// polaCiala składa pola wspólne obu kształtom żądania. Wartości są napisami, bo
// postać wieloczęściowa nie zna typów — a rozjazd typów między dwiema drogami
// dawałby dwa różne żądania z jednego wiersza rejestru.
func (k *kanalObrazow) polaCiala(z Zapytanie, polecenie string) map[string]any {
	pola := map[string]any{
		"model":  z.WybranyModel(k.def),
		"prompt": polecenie,
		"n":      k.liczbaObrazow(),
		"size":   k.def.ParametrLub("rozmiar", "1024x1024"),
	}
	// `response_format` znają nie wszystkie bramy zgodne z OpenAI — pole
	// nieznane bywa odrzucane, więc jedzie wyłącznie wtedy, gdy wiersz je
	// wskazał. Milczenie parametru znaczy „zostaw dostawcy jego domyślne".
	if format := strings.TrimSpace(k.def.Parametr("format_odpowiedzi")); format != "" {
		pola["response_format"] = format
	}
	for klucz, wartosc := range dodatkiCiala(k.def) {
		pola[klucz] = wartosc
	}
	return pola
}

// cialoJSON składa ciało dokumentem JSON. Materiał i maska jadą tą drogą jako
// napisy base64 pod nazwami z wiersza — tak przyjmują je bramy lokalne, które
// nie znają wysyłki wieloczęściowej.
func (k *kanalObrazow) cialoJSON(z Zapytanie, polecenie string) (io.Reader, string, error) {
	cialo := k.polaCiala(z, polecenie)
	if material, jest := z.ObrazRoli(RolaObrazuMaterial); jest {
		cialo[k.def.ParametrLub("pole_obrazu", "image")] = material.Base64
		if maska, jest := z.ObrazRoli(RolaObrazuMaska); jest {
			cialo[k.def.ParametrLub("pole_maski", "mask")] = maska.Base64
		}
	}
	surowe, err := json.Marshal(cialo)
	if err != nil {
		return nil, "", fmt.Errorf("models: kanał %s: kodowanie żądania: %w", k.def.Kod, err)
	}
	return bytes.NewReader(surowe), "application/json", nil
}

// cialoWieloczesciowe składa ciało jako multipart/form-data — kształt punktu
// końcowego edycji obrazu. Bajty rozkodowuje się TUTAJ, bo pole niosło je
// zapisem base64, a formularz oczekuje pliku; base64 wpisane w pole pliku
// dostawca odrzuciłby jako obraz nieczytelny.
//
// Materiał nieczytelny jest odmową, nie wysyłką bez materiału: żądanie samego
// polecenia wróciłoby obrazem wygenerowanym od zera, a Operator prosił o obróbkę
// swojego zdjęcia i dostałby cudze bez ani jednego słowa o podmianie.
func (k *kanalObrazow) cialoWieloczesciowe(z Zapytanie, polecenie string) (io.Reader, string, error) {
	bufor := &bytes.Buffer{}
	formularz := multipart.NewWriter(bufor)
	for klucz, wartosc := range k.polaCiala(z, polecenie) {
		if err := formularz.WriteField(klucz, fmt.Sprint(wartosc)); err != nil {
			return nil, "", fmt.Errorf("models: kanał %s: zapis pola %s: %w", k.def.Kod, klucz, err)
		}
	}
	material, _ := z.ObrazRoli(RolaObrazuMaterial)
	if err := k.dolozPlikObrazu(formularz, k.def.ParametrLub("pole_obrazu", "image"), material); err != nil {
		return nil, "", err
	}
	if maska, jest := z.ObrazRoli(RolaObrazuMaska); jest {
		if err := k.dolozPlikObrazu(formularz, k.def.ParametrLub("pole_maski", "mask"), maska); err != nil {
			return nil, "", err
		}
	}
	if err := formularz.Close(); err != nil {
		return nil, "", fmt.Errorf("models: kanał %s: domknięcie formularza: %w", k.def.Kod, err)
	}
	return bufor, formularz.FormDataContentType(), nil
}

// dolozPlikObrazu wkłada jeden obraz do formularza jako plik.
func (k *kanalObrazow) dolozPlikObrazu(formularz *multipart.Writer, pole string, obraz ObrazWejsciowy) error {
	bajty, err := base64.StdEncoding.DecodeString(strings.TrimSpace(obraz.Base64))
	if err != nil {
		return fmt.Errorf("models: kanał %s: obraz wejściowy w polu %s nie jest czytelnym zapisem base64: %w",
			k.def.Kod, pole, err)
	}
	czesc, err := formularz.CreateFormFile(pole, nazwaPlikuObrazu(obraz))
	if err != nil {
		return fmt.Errorf("models: kanał %s: zapis pliku %s: %w", k.def.Kod, pole, err)
	}
	if _, err := czesc.Write(bajty); err != nil {
		return fmt.Errorf("models: kanał %s: zapis bajtów pola %s: %w", k.def.Kod, pole, err)
	}
	return nil
}

// nazwaPlikuObrazu ustala nazwę części formularza. Nazwa własna obrazu wygrywa;
// bez niej nazwa bierze się z typu treści, bo część dostawców rozpoznaje format
// po rozszerzeniu, a nie po nagłówku części. Brak obu daje PNG — format, w którym
// warsztat fotografii oddaje wszystkie swoje wyniki.
func nazwaPlikuObrazu(obraz ObrazWejsciowy) string {
	if nazwa := strings.TrimSpace(obraz.Nazwa); nazwa != "" {
		return nazwa
	}
	rozszerzenie := "png"
	if _, podtyp, jest := strings.Cut(strings.TrimSpace(obraz.TypTresci), "/"); jest && podtyp != "" {
		rozszerzenie = podtyp
	}
	return obraz.RolaLub() + "." + rozszerzenie
}

// liczbaObrazow odczytuje zamówioną liczbę obrazów; wartość nieczytelna albo
// niedodatnia daje jeden — parametr pomocniczy nie wywraca wywołania.
func (k *kanalObrazow) liczbaObrazow() int {
	liczba, err := strconv.Atoi(strings.TrimSpace(k.def.ParametrLub("liczba", "1")))
	if err != nil || liczba <= 0 {
		return 1
	}
	return liczba
}

// nadajObraz wyjmuje treść wizualną z odpowiedzi i nadaje ją jednym fragmentem.
//
// Odpowiedź czytamy w całości, nie strumieniem: punkty końcowe obrazowe oddają
// jeden dokument JSON po zakończeniu generowania, a nie ciąg zdarzeń. Bajty mają
// pierwszeństwo przed adresem — obraz, który już przyszedł, nie wymaga drugiego
// pobrania i nie wygaśnie odbiorcy w ręku.
//
// Odpowiedź bez obrazu jest błędem, nie pustym wynikiem. Cisza w tym miejscu
// dawałaby wołającemu strumień bez treści, nieodróżnialny od udanego wywołania,
// które nic nie narysowało — dlatego mówimy, pod jakimi ścieżkami szukaliśmy.
func (k *kanalObrazow) nadajObraz(ctx context.Context, z Zapytanie, u Ujscie, polecenie string, tresc io.Reader) error {
	surowe, err := io.ReadAll(tresc)
	if err != nil {
		return fmt.Errorf("models: kanał %s: odczyt odpowiedzi: %w", k.def.Kod, err)
	}
	sciezkaBajtow := k.def.ParametrLub("sciezka_base64", "data.0.b64_json")
	sciezkaAdresu := k.def.ParametrLub("sciezka_adresu", "data.0.url")

	obraz := TrescObrazu{TypTresci: k.def.ParametrLub("typ_tresci", "image/png"), Prompt: polecenie}
	if bajty, jest := WartoscZeSciezki(surowe, sciezkaBajtow); jest {
		obraz.Base64 = bajty
	} else if adres, jest := WartoscZeSciezki(surowe, sciezkaAdresu); jest {
		obraz.Adres = adres
	} else {
		return fmt.Errorf("models: kanał %s: odpowiedź bez obrazu — ani pod ścieżką %s (bajty), "+
			"ani pod %s (adres)", k.def.Kod, sciezkaBajtow, sciezkaAdresu)
	}
	return NadajObraz(ctx, u, z, obraz)
}
