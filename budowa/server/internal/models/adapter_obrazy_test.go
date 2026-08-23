package models

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Sprawdziany obrazu WEJŚCIOWEGO kanału obrazowego.
//
// Mierzone jest to, czego nie da się zobaczyć po stronie wołającego: kształt
// żądania, które naprawdę wyszło na sieć. Kanał, który przyjmie materiał
// Operatora i wyśle samo polecenie, oddaje obraz WYGENEROWANY od zera —
// wygląda to na powodzenie, a jest podmianą materiału bez ani jednego słowa.
// Dlatego sprawdziany stawiają zaślepkę punktu końcowego i czytają ciało.

// obrazPrzykladowy niesie bajty, które da się rozpoznać po drugiej stronie.
// Treść jest dowolna — sprawdzian pyta o drogę bajtów, nie o obraz.
var obrazPrzykladowy = base64.StdEncoding.EncodeToString([]byte("bajty materiału"))

// odpowiedzZObrazem jest odpowiedzią w kształcie, który kanał umie odczytać.
const odpowiedzZObrazem = `{"data":[{"b64_json":"d3luaWs="}]}`

// wierszKanaluObrazow składa wiersz rejestru wskazujący zaślepkę.
func wierszKanaluObrazow(t *testing.T, adres string, dodatki map[string]any) Definicja {
	t.Helper()
	parametry := map[string]any{"base_url": adres, "adapter": AdapterObrazy}
	for klucz, wartosc := range dodatki {
		parametry[klucz] = wartosc
	}
	t.Setenv("KLUCZ_KANALU_OBRAZOW", "sekret-sprawdzianu")
	return Definicja{
		Kod:                    "obrazy-sprawdzian",
		Rodzaj:                 "api",
		Model:                  "model-obrazowy",
		PoswiadczenieOdwolanie: "KLUCZ_KANALU_OBRAZOW",
		Parametry:              parametry,
		Aktywny:                true,
	}
}

// przyjeteZadanie niesie to, co zaświadczyła zaślepka.
type przyjeteZadanie struct {
	sciezka   string
	typTresci string
	cialo     []byte
}

// zaslepkaObrazow stawia punkt końcowy zapisujący przyjęte żądanie.
func zaslepkaObrazow(t *testing.T, przyjete *przyjeteZadanie) *httptest.Server {
	t.Helper()
	serwer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cialo, _ := io.ReadAll(r.Body)
		przyjete.sciezka = r.URL.Path
		przyjete.typTresci = r.Header.Get("Content-Type")
		przyjete.cialo = cialo
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(odpowiedzZObrazem))
	}))
	t.Cleanup(serwer.Close)
	return serwer
}

// wyslijKanalem puszcza jedno wywołanie i zbiera fragmenty.
func wyslijKanalem(t *testing.T, d Definicja, z Zapytanie) []Fragment {
	t.Helper()
	kanal, err := NowyKanalObrazow(d)
	if err != nil {
		t.Fatalf("kanał obrazowy nie powstał: %v", err)
	}
	zebrane := []Fragment{}
	ujscie := UjscieFunkcji(func(_ context.Context, f Fragment) error {
		zebrane = append(zebrane, f)
		return nil
	})
	if err := kanal.Wyslij(context.Background(), z, ujscie); err != nil {
		t.Fatalf("wysyłka nie powiodła się: %v", err)
	}
	return zebrane
}

// TestWywolanieBezObrazuIdziePoStaremu pilnuje, żeby dołożenie pola obrazu
// wejściowego nie ruszyło kanałów już założonych: wywołanie bez materiału ma
// wyjść dokumentem JSON pod adres z `base_url`.
func TestWywolanieBezObrazuIdziePoStaremu(t *testing.T) {
	przyjete := &przyjeteZadanie{}
	serwer := zaslepkaObrazow(t, przyjete)
	wiersz := wierszKanaluObrazow(t, serwer.URL+"/generacje", nil)

	wyslijKanalem(t, wiersz, Zapytanie{
		Zasiegi: Zasiegi{Okno: "okno-1"}, Wiadomosc: "wiad-1", Tresc: "kot w kapeluszu",
	})

	if przyjete.sciezka != "/generacje" {
		t.Errorf("żądanie poszło pod %q, a miało pod /generacje", przyjete.sciezka)
	}
	if !strings.HasPrefix(przyjete.typTresci, "application/json") {
		t.Errorf("rodzaj treści %q, a miał być application/json", przyjete.typTresci)
	}
	var cialo map[string]any
	if err := json.Unmarshal(przyjete.cialo, &cialo); err != nil {
		t.Fatalf("ciało nie jest dokumentem JSON: %v", err)
	}
	if cialo["prompt"] != "kot w kapeluszu" {
		t.Errorf("ciało niesie prompt %v", cialo["prompt"])
	}
	if _, jest := cialo["image"]; jest {
		t.Error("wywołanie bez materiału wysłało pole obrazu")
	}
}

// TestMaterialWejsciowyIdzieWieloczesciowoPodAdresEdycji jest sprawdzianem
// głównym: materiał i maska mają wyjść jako pliki formularza, pod adres edycji,
// z bajtami rozkodowanymi z base64.
func TestMaterialWejsciowyIdzieWieloczesciowoPodAdresEdycji(t *testing.T) {
	przyjete := &przyjeteZadanie{}
	serwer := zaslepkaObrazow(t, przyjete)
	wiersz := wierszKanaluObrazow(t, serwer.URL+"/generacje", map[string]any{
		"adres_edycji": serwer.URL + "/edycje",
	})

	wyslijKanalem(t, wiersz, Zapytanie{
		Zasiegi:   Zasiegi{Okno: "okno-1"},
		Wiadomosc: "wiad-1",
		Tresc:     "usuń tło",
		ObrazyWejsciowe: []ObrazWejsciowy{
			{Rola: RolaObrazuMaterial, TypTresci: "image/png", Base64: obrazPrzykladowy},
			{Rola: RolaObrazuMaska, TypTresci: "image/png", Base64: obrazPrzykladowy},
		},
	})

	if przyjete.sciezka != "/edycje" {
		t.Errorf("edycja poszła pod %q, a miała pod /edycje", przyjete.sciezka)
	}
	pliki, pola := rozbierzFormularz(t, przyjete)
	if got := string(pliki["image"]); got != "bajty materiału" {
		t.Errorf("pole image niesie %q — bajty nie zostały rozkodowane z base64", got)
	}
	if _, jest := pliki["mask"]; !jest {
		t.Error("maska nie wyszła osobnym polem formularza")
	}
	if pola["prompt"] != "usuń tło" {
		t.Errorf("formularz niesie prompt %q", pola["prompt"])
	}
	if pola["model"] != "model-obrazowy" {
		t.Errorf("formularz niesie model %q", pola["model"])
	}
}

// TestNazwyPolObrazuBioraSieZWiersza pilnuje, żeby dostawca nazywający pola
// inaczej dał się obsłużyć wierszem rejestru, a nie zmianą w kodzie.
func TestNazwyPolObrazuBioraSieZWiersza(t *testing.T) {
	przyjete := &przyjeteZadanie{}
	serwer := zaslepkaObrazow(t, przyjete)
	wiersz := wierszKanaluObrazow(t, serwer.URL+"/edycje", map[string]any{
		"pole_obrazu": "init_image",
		"pole_maski":  "mask_image",
	})

	wyslijKanalem(t, wiersz, Zapytanie{
		Zasiegi:   Zasiegi{Okno: "okno-1"},
		Wiadomosc: "wiad-1",
		Tresc:     "rozszerz kadr",
		ObrazyWejsciowe: []ObrazWejsciowy{
			{Base64: obrazPrzykladowy},
			{Rola: RolaObrazuMaska, Base64: obrazPrzykladowy},
		},
	})

	pliki, _ := rozbierzFormularz(t, przyjete)
	if _, jest := pliki["init_image"]; !jest {
		t.Error("materiał nie wyszedł pod nazwą pola z wiersza rejestru")
	}
	if _, jest := pliki["mask_image"]; !jest {
		t.Error("maska nie wyszła pod nazwą pola z wiersza rejestru")
	}
}

// TestPostacBase64WysylaMaterialWCieleJSON sprawdza drugą postać wysyłki — tę,
// którą przyjmują bramy lokalne nieznające formularzy wieloczęściowych.
func TestPostacBase64WysylaMaterialWCieleJSON(t *testing.T) {
	przyjete := &przyjeteZadanie{}
	serwer := zaslepkaObrazow(t, przyjete)
	wiersz := wierszKanaluObrazow(t, serwer.URL+"/edycje", map[string]any{
		"postac_obrazu": postacObrazuBase64,
	})

	wyslijKanalem(t, wiersz, Zapytanie{
		Zasiegi:   Zasiegi{Okno: "okno-1"},
		Wiadomosc: "wiad-1",
		Tresc:     "powiększ",
		ObrazyWejsciowe: []ObrazWejsciowy{
			{Rola: RolaObrazuMaterial, Base64: obrazPrzykladowy},
		},
	})

	if !strings.HasPrefix(przyjete.typTresci, "application/json") {
		t.Fatalf("postać base64 wysłała %q, a miała dokument JSON", przyjete.typTresci)
	}
	var cialo map[string]any
	if err := json.Unmarshal(przyjete.cialo, &cialo); err != nil {
		t.Fatalf("ciało nie jest dokumentem JSON: %v", err)
	}
	if cialo["image"] != obrazPrzykladowy {
		t.Errorf("pole image niesie %v, a miało zapis base64 materiału", cialo["image"])
	}
}

// TestMaterialNieczytelnyJestOdmowa pilnuje rozstrzygnięcia z nagłówka adaptera:
// materiał, którego nie da się rozkodować, ZATRZYMUJE wywołanie. Wysyłka samego
// polecenia oddałaby obraz cudzy pod postacią obróbki zdjęcia Operatora.
func TestMaterialNieczytelnyJestOdmowa(t *testing.T) {
	przyjete := &przyjeteZadanie{}
	serwer := zaslepkaObrazow(t, przyjete)
	wiersz := wierszKanaluObrazow(t, serwer.URL+"/edycje", nil)

	kanal, err := NowyKanalObrazow(wiersz)
	if err != nil {
		t.Fatalf("kanał obrazowy nie powstał: %v", err)
	}
	err = kanal.Wyslij(context.Background(), Zapytanie{
		Zasiegi:         Zasiegi{Okno: "okno-1"},
		Wiadomosc:       "wiad-1",
		Tresc:           "popraw",
		ObrazyWejsciowe: []ObrazWejsciowy{{Base64: "to nie jest base64 %%%"}},
	}, nil)
	if err == nil {
		t.Fatal("materiał nieczytelny przeszedł bez odmowy")
	}
	if !strings.Contains(err.Error(), "base64") {
		t.Errorf("odmowa nie nazywa powodu: %v", err)
	}
	if przyjete.cialo != nil {
		t.Error("żądanie mimo nieczytelnego materiału poleciało na sieć")
	}
}

// rozbierzFormularz czyta przyjęte ciało jako multipart/form-data.
func rozbierzFormularz(t *testing.T, przyjete *przyjeteZadanie) (map[string][]byte, map[string]string) {
	t.Helper()
	rodzaj, parametry, err := mime.ParseMediaType(przyjete.typTresci)
	if err != nil || rodzaj != "multipart/form-data" {
		t.Fatalf("ciało nie jest formularzem wieloczęściowym: %q (%v)", przyjete.typTresci, err)
	}
	czytnik := multipart.NewReader(strings.NewReader(string(przyjete.cialo)), parametry["boundary"])
	formularz, err := czytnik.ReadForm(1 << 20)
	if err != nil {
		t.Fatalf("nie można rozebrać formularza: %v", err)
	}
	pliki := map[string][]byte{}
	for nazwa, naglowki := range formularz.File {
		otwarty, err := naglowki[0].Open()
		if err != nil {
			t.Fatalf("nie można otworzyć części %s: %v", nazwa, err)
		}
		bajty, _ := io.ReadAll(otwarty)
		_ = otwarty.Close()
		pliki[nazwa] = bajty
	}
	pola := map[string]string{}
	for nazwa, wartosci := range formularz.Value {
		pola[nazwa] = wartosci[0]
	}
	return pliki, pola
}
