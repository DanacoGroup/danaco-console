package core

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"image"
	"image/color"
	"image/png"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"

	"danacoconsole/shared"
)

// Skutek DROGI NEURONOWEJ czterech czynności warsztatu fotografii —
// `design.photo.upscale`, `.background.remove`, `.inpaint`, `.expand`.
//
// ── Co było, a co jest ──────────────────────────────────────────────────────
// Przez jedną turę te cztery czynności miały wariant neuronowy tylko na papierze:
// kanał obrazowy przyjmował samo polecenie tekstowe, więc wskazany `channelId`
// był SPRAWDZANY i pomijany, a odpowiedź oddawała `computedBy: rachunekRdzenia`.
// Po dołożeniu obrazu wejściowego do warstwy modeli droga stoi i ten plik mierzy
// ją tak, jak trzeba mierzyć wywołanie cudzego silnika: sprawdzianem, który
// SPRAWDZA ŻĄDANIE WYCHODZĄCE, a nie tylko odpowiedź.
//
// Trzy rzeczy są mierzone przy każdej czynności:
//
//  1. żądanie do kanału NIESIE zdjęcie Operatora (materiał), a przy domalowaniu
//     i rozszerzeniu kadru także maskę. Wywołanie bez materiału kazałoby silnikowi
//     wygenerować obraz NOWY i Operator dostałby cudzą treść pod swoim zdjęciem;
//  2. plik wyniku ma wymiar, który czynność OBIECAŁA — nie ten, który kanał
//     akurat oddał;
//  3. `computedBy` mówi `kanalModelu`, a łańcuch edycji zapisuje tę samą drogę.

// zadanieKanaluSprawdzianu jest jednym żądaniem przechwyconym przez serwer
// próbny — rozłożonym na pola, żeby sprawdzian pytał o materiał i maskę, a nie
// o bajty ciała.
type zadanieKanaluSprawdzianu struct {
	Sciezka   string
	Polecenie string
	Material  []byte
	Maska     []byte
}

// serwerKanaluEdycjiSprawdzianu stawia kanał obrazowy, który zapisuje żądania
// i odpowiada wskazanym obrazem.
//
// Serwer przyjmuje OBA kształty żądania warstwy modeli (wieloczęściowy
// i base64), bo sprawdzian nie ma prawa zakładać, którym wiersz rejestru
// pojedzie — a oba niosą materiał w tym samym polu.
func serwerKanaluEdycjiSprawdzianu(t *testing.T, odpowiedz []byte) (*[]zadanieKanaluSprawdzianu,
	http.HandlerFunc) {

	t.Helper()

	zadania := &[]zadanieKanaluSprawdzianu{}
	return zadania, func(w http.ResponseWriter, z *http.Request) {
		zapis := zadanieKanaluSprawdzianu{Sciezka: z.URL.Path}
		typ, parametry, err := mime.ParseMediaType(z.Header.Get("Content-Type"))
		switch {
		case err == nil && strings.HasPrefix(typ, "multipart/"):
			czesci := multipart.NewReader(z.Body, parametry["boundary"])
			for {
				czesc, err := czesci.NextPart()
				if err != nil {
					break
				}
				bajty := &bytes.Buffer{}
				_, _ = bajty.ReadFrom(czesc)
				switch czesc.FormName() {
				case "image":
					zapis.Material = bajty.Bytes()
				case "mask":
					zapis.Maska = bajty.Bytes()
				case "prompt":
					zapis.Polecenie = bajty.String()
				}
			}
		default:
			var cialo map[string]any
			if err := json.NewDecoder(z.Body).Decode(&cialo); err == nil {
				if napis, jest := cialo["image"].(string); jest {
					zapis.Material, _ = base64.StdEncoding.DecodeString(napis)
				}
				if napis, jest := cialo["mask"].(string); jest {
					zapis.Maska, _ = base64.StdEncoding.DecodeString(napis)
				}
				if napis, jest := cialo["prompt"].(string); jest {
					zapis.Polecenie = napis
				}
			}
		}
		*zadania = append(*zadania, zapis)

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": []any{map[string]any{"b64_json": wBase64(odpowiedz)}},
		})
	}
}

// obrazZKryciemPNG składa obraz o zadanym boku, w którym RAMKA jest przezroczysta
// a środek kryjący — materiał, na którym da się zmierzyć, czy kanał oddał obraz
// z kanałem krycia.
func obrazZKryciemPNG(t *testing.T, bok int) []byte {
	t.Helper()

	plotno := image.NewNRGBA(image.Rect(0, 0, bok, bok))
	for y := 0; y < bok; y++ {
		for x := 0; x < bok; x++ {
			if x < bok/4 || y < bok/4 || x >= 3*bok/4 || y >= 3*bok/4 {
				plotno.SetNRGBA(x, y, color.NRGBA{})
				continue
			}
			plotno.SetNRGBA(x, y, color.NRGBA{R: 0x20, G: 0x60, B: 0xc0, A: 0xff})
		}
	}
	var bufor bytes.Buffer
	if err := png.Encode(&bufor, plotno); err != nil {
		t.Fatalf("nie można złożyć obrazu sprawdzianu: %v", err)
	}
	return bufor.Bytes()
}

// wymiaryObrazuSprawdzianu odczytuje wymiary z bajtów obrazu.
func wymiaryObrazuSprawdzianu(t *testing.T, bajty []byte) (int, int) {
	t.Helper()

	obraz, _, err := image.Decode(bytes.NewReader(bajty))
	if err != nil {
		t.Fatalf("bajty nie są obrazem: %v", err)
	}
	granice := obraz.Bounds()
	return granice.Dx(), granice.Dy()
}

// TestPowiekszenieKanalemNiesieZdjecieOperatoraIWymiarZamowiony mierzy drogę
// neuronową powiększenia: żądanie do kanału ma nieść zdjęcie wniesione przez
// Operatora, plik wyniku ma mieć wymiar zamówionej krotności, a odpowiedź ma
// mówić `kanalModelu`.
func TestPowiekszenieKanalemNiesieZdjecieOperatoraIWymiarZamowiony(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)

	// Kanał oddaje obraz o wymiarze WŁASNYM (kwadrat 64×64), różnym od
	// zamówionego 80×60 — tak zachowuje się punkt końcowy z nastawą `size`.
	zadania, obsluga := serwerKanaluEdycjiSprawdzianu(t, obrazPNG(t, 64, 64))
	serwer := serwerObrazow(t, obsluga)
	kanal := wpiszKanalObrazowy(t, zmontowany, zycie, serwer.URL)

	zrodloBajty := obrazPNG(t, 20, 15)
	zrodlo := wniesObrazSprawdzianu(t, zmontowany, zycie, "okno-neuronowe", zrodloBajty)

	var wynik shared.DesignPhotoUpscaleResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignPhotoUpscale,
		shared.DesignPhotoUpscaleRequest{AssetId: zrodlo.Id, Factor: 4, ChannelId: &kanal},
		&wynik)

	if wynik.ComputedBy != shared.DesignPhotoComputeRouteKanalModelu {
		t.Errorf("odpowiedź mówi, że policzył %q, a wskazano kanał — droga ma być kanałem modelu",
			wynik.ComputedBy)
	}
	if len(*zadania) != 1 {
		t.Fatalf("kanał dostał %d żądań, a czynność jest jedna", len(*zadania))
	}
	zadanie := (*zadania)[0]
	if len(zadanie.Material) == 0 {
		t.Fatal("żądanie do kanału nie niosło ani jednego bajtu materiału — silnik dostałby " +
			"samo polecenie i oddał obraz NOWY, nie powiększone zdjęcie Operatora")
	}
	// Materiał ma być zdjęciem Operatora, nie czymkolwiek: wymiary z pliku, który
	// pojechał, muszą być wymiarami źródła.
	szerokosc, wysokosc := wymiaryObrazuSprawdzianu(t, zadanie.Material)
	if szerokosc != 20 || wysokosc != 15 {
		t.Errorf("do kanału pojechał obraz %d×%d, a zdjęcie Operatora ma 20×15",
			szerokosc, wysokosc)
	}
	if len(zadanie.Maska) != 0 {
		t.Error("powiększenie wysłało maskę — maski nie ma w tym żądaniu i punkt końcowy " +
			"mógłby ją wziąć za wskazanie obszaru pracy")
	}
	if !strings.Contains(zadanie.Polecenie, "4") {
		t.Errorf("polecenie do kanału nie mówi o krotności: %q", zadanie.Polecenie)
	}

	// PLIK ma wymiar ZAMÓWIONY (20×15 razy cztery), nie ten, który oddał kanał.
	obraz := obrazZMagazynuFotografii(t, katalog, wynik.Asset)
	granice := obraz.Bounds()
	if granice.Dx() != 80 || granice.Dy() != 60 {
		t.Errorf("plik po powiększeniu ×4 ma %d×%d, a zamówienie z 20×15 daje 80×60",
			granice.Dx(), granice.Dy())
	}
	if wynik.Width != granice.Dx() || wynik.Height != granice.Dy() {
		t.Errorf("odpowiedź podaje %d×%d, a plik ma %d×%d",
			wynik.Width, wynik.Height, granice.Dx(), granice.Dy())
	}

	// Łańcuch edycji zapisuje DROGĘ — bez tego po tygodniu nie da się powiedzieć,
	// którym rachunkiem powstał wariant.
	var historia shared.DesignPhotoHistoryGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignPhotoHistoryGet,
		shared.DesignPhotoHistoryGetRequest{AssetId: wynik.Asset.Id}, &historia)
	if len(historia.Edits) == 0 {
		t.Fatal("łańcuch edycji wariantu jest pusty")
	}
	znaleziona := false
	for _, ogniwo := range historia.Edits {
		if ogniwo.ComputedBy != nil &&
			*ogniwo.ComputedBy == shared.DesignPhotoComputeRouteKanalModelu {
			znaleziona = true
		}
	}
	if !znaleziona {
		t.Error("żadne ogniwo łańcucha nie mówi, że wariant policzył kanał modelu")
	}
}

// TestOdcieciecieTlaKanalemZadaKanaluKryciaWPliku mierzy warunek, bez którego
// odcięcie tła nie jest odcięciem tła: plik oddany przez kanał MUSI mieć punkty
// przezroczyste. Kanał oddający obraz kryjący dostaje odmowę nazwaną, a nie
// odpowiedź `hasAlpha: true` nad plikiem bez przezroczystości.
func TestOdcieciecieTlaKanalemZadaKanaluKryciaWPliku(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)

	// Najpierw kanał oddający obraz Z kryciem — droga udana.
	zadania, obsluga := serwerKanaluEdycjiSprawdzianu(t, obrazZKryciemPNG(t, 40))
	serwer := serwerObrazow(t, obsluga)
	kanal := wpiszKanalObrazowy(t, zmontowany, zycie, serwer.URL)

	zrodlo := wniesObrazSprawdzianu(t, zmontowany, zycie, "okno-neuronowe",
		obrazZPrzedmiotemPNG(t, 40))

	var wynik shared.DesignPhotoBackgroundRemoveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignPhotoBackgroundRemove,
		shared.DesignPhotoBackgroundRemoveRequest{AssetId: zrodlo.Id, ChannelId: &kanal},
		&wynik)

	if wynik.ComputedBy != shared.DesignPhotoComputeRouteKanalModelu {
		t.Errorf("odpowiedź mówi %q, a wskazano kanał", wynik.ComputedBy)
	}
	if len(*zadania) != 1 || len((*zadania)[0].Material) == 0 {
		t.Fatal("do kanału nie pojechał materiał — odcięcie tła bez zdjęcia nie ma czego odciąć")
	}
	if wynik.TransparentShare <= 0 || wynik.TransparentShare >= 1 {
		t.Errorf("udział punktów przezroczystych to %v; obraz z ramką przezroczystą ma go "+
			"między zerem a jednością", wynik.TransparentShare)
	}
	// Pomiar niezależny na pliku: udział z odpowiedzi ma się zgadzać z plikiem.
	obraz := obrazZMagazynuFotografii(t, katalog, wynik.Asset)
	if zmierzony := udzialPrzezroczystosciDesignu(obraz); zmierzony-wynik.TransparentShare > 0.02 ||
		wynik.TransparentShare-zmierzony > 0.02 {

		t.Errorf("plik ma %.3f punktów przezroczystych, a odpowiedź mówi %.3f",
			zmierzony, wynik.TransparentShare)
	}

	// Teraz kanał oddający obraz BEZ krycia — odmowa nazwana.
	_, obslugaKryjaca := serwerKanaluEdycjiSprawdzianu(t, obrazPNG(t, 40, 40))
	serwerKryjacy := serwerObrazow(t, obslugaKryjaca)
	kanalKryjacy := wpiszKanalObrazowy(t, zmontowany, zycie, serwerKryjacy.URL)

	odmowa := wykonajOdmowna(t, zmontowany, zycie, shared.CommandDesignPhotoBackgroundRemove,
		shared.DesignPhotoBackgroundRemoveRequest{
			AssetId: zrodlo.Id, ChannelId: &kanalKryjacy,
		})
	if !strings.Contains(odmowa.Message, "przezroczyst") {
		t.Errorf("odmowa nie nazywa braku przezroczystości: %s", odmowa.Message)
	}
}

// TestDomalowanieIRozszerzenieKanalemWysylajaMaske mierzy to, co odróżnia te dwie
// czynności od pozostałych: żądanie musi nieść MASKĘ wskazującą obszar pracy,
// a przy rozszerzeniu kadru materiał musi być płótnem o wymiarze WYNIKU.
func TestDomalowanieIRozszerzenieKanalemWysylajaMaske(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)

	zadania, obsluga := serwerKanaluEdycjiSprawdzianu(t, obrazPNG(t, 96, 96))
	serwer := serwerObrazow(t, obsluga)
	kanal := wpiszKanalObrazowy(t, zmontowany, zycie, serwer.URL)

	zrodlo := wniesObrazSprawdzianu(t, zmontowany, zycie, "okno-neuronowe", obrazPNG(t, 40, 30))

	var domalowanie shared.DesignPhotoInpaintResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignPhotoInpaint,
		shared.DesignPhotoInpaintRequest{
			AssetId:   zrodlo.Id,
			Regions:   []shared.DesignPhotoRegion{{X: 8, Y: 8, Width: 12, Height: 10}},
			Prompt:    wskaznik("dorysuj ceglaną ścianę"),
			ChannelId: &kanal,
		}, &domalowanie)

	if domalowanie.ComputedBy != shared.DesignPhotoComputeRouteKanalModelu {
		t.Errorf("domalowanie mówi %q, a wskazano kanał", domalowanie.ComputedBy)
	}
	if len(*zadania) != 1 {
		t.Fatalf("kanał dostał %d żądań, a czynność jest jedna", len(*zadania))
	}
	zadanie := (*zadania)[0]
	if len(zadanie.Maska) == 0 {
		t.Fatal("domalowanie nie wysłało maski — silnik nie wie wtedy, KTÓRY fragment " +
			"domalować, i przerysowałby całe zdjęcie")
	}
	// Maska ma wymiar zdjęcia: maska mniejsza wskazywałaby inny fragment, niż
	// Operator zaznaczył.
	szerokoscMaski, wysokoscMaski := wymiaryObrazuSprawdzianu(t, zadanie.Maska)
	if szerokoscMaski != 40 || wysokoscMaski != 30 {
		t.Errorf("maska ma %d×%d, a zdjęcie 40×30", szerokoscMaski, wysokoscMaski)
	}
	// Obszar objęty jest w masce PRZEZROCZYSTY i biały naraz — powód stoi przy
	// `obrazMaskiDesignu`. Sprawdzian mierzy oba warunki, bo punkty końcowe czytają
	// raz jedno, raz drugie.
	maska, _, err := image.Decode(bytes.NewReader(zadanie.Maska))
	if err != nil {
		t.Fatalf("maska nie jest obrazem: %v", err)
	}
	// Odczyt idzie po składowych BEZ wmnożonego krycia. Zwykłe `At().RGBA()`
	// wmnaża krycie, więc punkt biały o kryciu zerowym pokazywałby się jako
	// czarny — a w pliku PNG, który pojechał do kanału, stoi biały. Sprawdzian
	// mierzy PLIK, nie sposób, w jaki biblioteka go pokazuje.
	bezKrycia, jest := maska.(*image.NRGBA)
	if !jest {
		t.Fatalf("maska rozłożyła się jako %T, a nie jako obraz bez wmnożonego krycia — "+
			"nie da się wtedy odczytać składowych punktu przezroczystego", maska)
	}
	wObszarze := bezKrycia.NRGBAAt(12, 12)
	if wObszarze.A != 0 {
		t.Errorf("punkt (12;12) leży w zaznaczonym obszarze, a jego krycie to %d", wObszarze.A)
	}
	if wObszarze.R < 0xf0 || wObszarze.G < 0xf0 || wObszarze.B < 0xf0 {
		t.Errorf("punkt zaznaczonego obszaru nie jest biały (%d;%d;%d)",
			wObszarze.R, wObszarze.G, wObszarze.B)
	}
	if poza := bezKrycia.NRGBAAt(1, 1); poza.A == 0 {
		t.Error("punkt POZA zaznaczeniem jest w masce przezroczysty — punkt końcowy czytający " +
			"krycie przerobiłby wtedy całe zdjęcie")
	}
	if zadanie.Polecenie != "dorysuj ceglaną ścianę" {
		t.Errorf("polecenie Operatora nie pojechało do kanału: %q", zadanie.Polecenie)
	}
	// Wymiar wyniku zostaje wymiarem zdjęcia: domalowanie nie zmienia kadru.
	if granice := obrazZMagazynuFotografii(t, katalog, domalowanie.Asset).Bounds(); granice.Dx() != 40 ||
		granice.Dy() != 30 {

		t.Errorf("plik po domalowaniu ma %d×%d, a zdjęcie miało 40×30", granice.Dx(), granice.Dy())
	}

	// Rozszerzenie kadru: materiał jest płótnem o wymiarze WYNIKU, maska obejmuje
	// same marginesy.
	*zadania = nil
	var rozszerzenie shared.DesignPhotoExpandResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignPhotoExpand,
		shared.DesignPhotoExpandRequest{
			AssetId: zrodlo.Id, Left: wskaznik(10), Top: wskaznik(6), ChannelId: &kanal,
		}, &rozszerzenie)

	if rozszerzenie.ComputedBy != shared.DesignPhotoComputeRouteKanalModelu {
		t.Errorf("rozszerzenie kadru mówi %q, a wskazano kanał", rozszerzenie.ComputedBy)
	}
	if len(*zadania) != 1 {
		t.Fatalf("kanał dostał %d żądań na rozszerzenie kadru", len(*zadania))
	}
	rozszerzajace := (*zadania)[0]
	szerokoscMaterialu, wysokoscMaterialu := wymiaryObrazuSprawdzianu(t, rozszerzajace.Material)
	if szerokoscMaterialu != 50 || wysokoscMaterialu != 36 {
		t.Errorf("materiał rozszerzenia ma %d×%d, a 40×30 plus 10 z lewej i 6 od góry daje 50×36",
			szerokoscMaterialu, wysokoscMaterialu)
	}
	if len(rozszerzajace.Maska) == 0 {
		t.Fatal("rozszerzenie kadru nie wysłało maski marginesów")
	}
	maskaMarginesow, _, err := image.Decode(bytes.NewReader(rozszerzajace.Maska))
	if err != nil {
		t.Fatalf("maska marginesów nie jest obrazem: %v", err)
	}
	// Margines (2;2) jest obszarem pracy, środek dawnego zdjęcia (30;20) nie jest.
	marginesyBezKrycia, jest := maskaMarginesow.(*image.NRGBA)
	if !jest {
		t.Fatalf("maska marginesów rozłożyła się jako %T", maskaMarginesow)
	}
	if krycie := marginesyBezKrycia.NRGBAAt(2, 2).A; krycie != 0 {
		t.Errorf("punkt marginesu (2;2) ma w masce krycie %d, a margines jest obszarem pracy",
			krycie)
	}
	if krycie := marginesyBezKrycia.NRGBAAt(30, 20).A; krycie == 0 {
		t.Error("punkt dawnego zdjęcia jest w masce obszarem pracy — model przerysowałby " +
			"treść, którą Operator chciał zachować")
	}
	if rozszerzenie.Width != 50 || rozszerzenie.Height != 36 {
		t.Errorf("odpowiedź rozszerzenia podaje %d×%d, a zamówienie daje 50×36",
			rozszerzenie.Width, rozszerzenie.Height)
	}
}

// TestCzynnoscBezWskazanegoKanaluIdzieRachunkiemRdzenia pilnuje drugiej połowy
// rozstrzygnięcia kontraktu: pominięty `channelId` znaczy rachunek wkompilowany.
// Zejście na „pierwszy czynny kanał obrazowy" wysyłałoby zdjęcie Operatora do
// dostawcy, o którego nie prosił, i kosztowało go pieniądze bez słowa.
func TestCzynnoscBezWskazanegoKanaluIdzieRachunkiemRdzenia(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	// Kanał obrazowy STOI i jest czynny — a mimo to nie ma prawa dostać żądania.
	zadania, obsluga := serwerKanaluEdycjiSprawdzianu(t, obrazPNG(t, 64, 64))
	serwer := serwerObrazow(t, obsluga)
	_ = wpiszKanalObrazowy(t, zmontowany, zycie, serwer.URL)

	zrodlo := wniesObrazSprawdzianu(t, zmontowany, zycie, "okno-neuronowe", obrazPNG(t, 20, 15))

	var wynik shared.DesignPhotoUpscaleResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignPhotoUpscale,
		shared.DesignPhotoUpscaleRequest{AssetId: zrodlo.Id, Factor: 2}, &wynik)

	if wynik.ComputedBy != shared.DesignPhotoComputeRouteRachunekRdzenia {
		t.Errorf("bez wskazanego kanału odpowiedź mówi %q, a kontrakt mówi „rachunek "+
			"wkompilowany\"", wynik.ComputedBy)
	}
	if len(*zadania) != 0 {
		t.Errorf("kanał dostał %d żądań, a Operator o niego nie prosił", len(*zadania))
	}
}

// TestKanalNieoddajacyObrazuJestOdmowaNieRachunkiem pilnuje, żeby niepowodzenie
// wskazanego kanału NIE schodziło po cichu na rachunek wkompilowany: Operator
// prosił o drogę neuronową, a obraz policzony inaczej byłby odpowiedzią na inne
// żądanie.
func TestKanalNieoddajacyObrazuJestOdmowaNieRachunkiem(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	serwer := serwerObrazow(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte(`{"error":{"message":"silnik chwilowo nieczynny"}}`))
	})
	kanal := wpiszKanalObrazowy(t, zmontowany, zycie, serwer.URL)

	zrodlo := wniesObrazSprawdzianu(t, zmontowany, zycie, "okno-neuronowe", obrazPNG(t, 20, 15))

	odmowa := wykonajOdmowna(t, zmontowany, zycie, shared.CommandDesignPhotoUpscale,
		shared.DesignPhotoUpscaleRequest{AssetId: zrodlo.Id, Factor: 4, ChannelId: &kanal})
	if odmowa.Code != shared.ErrorCodeChannelUnavailable {
		t.Errorf("niepowodzenie kanału dało kod %s, a jest zapleczem niedostępnym", odmowa.Code)
	}

	// Wariant nie ma prawa powstać: zasób oddany po cichu rachunkiem wkompilowanym
	// leżałby w magazynie jako wynik drogi, którą nie poszedł.
	var wykaz shared.DesignAssetListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignAssetList,
		shared.DesignAssetListRequest{WindowId: wskaznik("okno-neuronowe")}, &wykaz)
	for _, zasob := range wykaz.Assets {
		if zasob.VariantOfAssetId != nil && *zasob.VariantOfAssetId == zrodlo.Id {
			t.Errorf("po odmowie kanału w magazynie leży wariant %s — komenda odmówiła "+
				"i zapisała", zasob.Id)
		}
	}
}
