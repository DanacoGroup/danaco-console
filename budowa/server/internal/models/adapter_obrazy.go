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

// kanalObrazow to kanał modelu oddający bajty obrazu zamiast tekstu, oparty
// o kształt żądania API generowania obrazów: adres, prompt, parametry
// rozmiaru i liczby oraz ścieżki odczytu wyniku, wraz z obsługą materiału
// wejściowego do edycji.
type kanalObrazow struct {
	def    Definicja
	klient *http.Client
}

// NowyKanalObrazow buduje kanał obrazowy z wiersza rejestru; brak parametru
// base_url jest jedynym powodem odmowy budowy, reszta parametrów ma wartości
// domyślne.
func NowyKanalObrazow(d Definicja) (Kanal, error) {
	if strings.TrimSpace(d.Parametr("base_url")) == "" {
		return nil, fmt.Errorf("models: kanał %q (adapter %s) bez parametru base_url", d.Kod, AdapterObrazy)
	}
	return &kanalObrazow{def: d, klient: &http.Client{Timeout: limitCzasu(d)}}, nil
}

// Kod zwraca kod kanału zapisany w wierszu rejestru, którym repozytorium
// jednoznacznie identyfikuje ten kanał obrazowy.
func (k *kanalObrazow) Kod() string {
	return k.def.Kod
}

// Definicja zwraca pełny wiersz rejestru, z którego kanał powstał, wraz ze
// wszystkimi jego parametrami.
func (k *kanalObrazow) Definicja() Definicja {
	return k.def
}

// Wyslij nadaje prowenancję, metadane konta i — po udanym żądaniu — jeden
// fragment obrazu. Kolejność jest ta sama, co w kanałach tekstowych:
// odbiorca poznaje warunki wywołania, zanim zobaczy wynik.
func (k *kanalObrazow) Wyslij(ctx context.Context, z Zapytanie, u Ujscie) error {
	prowenancja := ProwenancjaZapytania(z, k.def)
	// Adres w prowenancji to adres, którego wywołanie faktycznie użyje, nie
	// zawsze base_url.
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

// sprawdzPoswiadczenie odmawia, nazywając brak, zanim poleci żądanie: punkt
// końcowy generujący obrazy bez klucza nie istnieje, więc brak poświadczenia
// rozstrzyga się tu, zamiast zwracać cudzy błąd 401.
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

// poleceniePromptu składa polecenie obrazu: punkty końcowe obrazowe
// przyjmują jedno pole `prompt`, więc nakładka systemowa okna jedzie przed
// treścią, oddzielona pustą linią.
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

// Postacie wysyłki obrazu wejściowego — wartości parametru `postac_obrazu`
// rozpoznawane przy budowie ciała żądania edycji materiału.
const (
	postacObrazuWieloczesciowa = "wieloczesciowa"
	postacObrazuBase64         = "base64"
)

// zbudujZadanie składa żądanie HTTP kanału obrazowego: materiał wejściowy
// w postaci wieloczęściowej dostaje osobne ciało, a nagłówki, adres i klucz
// są wspólne dla obu kształtów żądania.
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
	// Pole response_format jedzie tylko, gdy wiersz je wskazał; puste znaczy
	// domyślne dostawcy.
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

// cialoWieloczesciowe składa ciało jako multipart/form-data, kształt punktu
// końcowego edycji obrazu: bajty rozkodowuje z zapisu base64, bo formularz
// oczekuje pliku, nie napisu.
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

// dolozPlikObrazu wkłada jeden obraz do formularza jako plik, rozkodowując
// jego zawartość z zapisu base64 do bajtów.
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

// nazwaPlikuObrazu ustala nazwę części formularza: nazwa własna obrazu
// wygrywa, w przeciwnym razie nazwa bierze się z typu treści, a brak obu
// daje rozszerzenie png.
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

// nadajObraz wyjmuje treść wizualną z odpowiedzi i nadaje ją jednym
// fragmentem: odpowiedź czyta w całości, nie strumieniem, a brak obrazu
// w odpowiedzi traktuje jako błąd.
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
