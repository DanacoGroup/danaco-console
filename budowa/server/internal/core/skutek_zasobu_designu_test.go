package core

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"danacoconsole/shared"
)

// zmiennaKluczaObrazow nazywa zmienną środowiska, z której kanał obrazowy
// bierze klucz poświadczenia wymagany nawet wtedy, gdy odbiorcą jest serwer
// podniesiony na czas sprawdzianu.
const zmiennaKluczaObrazow = "DANACO_SPRAWDZIAN_KLUCZ_OBRAZOW"

// serwerObrazow podnosi w tym samym procesie punkt końcowy generowania
// obrazów na czas trwania sprawdzianu i sam się zamyka po jego zakończeniu.
func serwerObrazow(t *testing.T, obsluga http.HandlerFunc) *httptest.Server {
	t.Helper()

	serwer := httptest.NewServer(obsluga)
	t.Cleanup(serwer.Close)
	return serwer
}

// odpowiedzZBajtami składa odpowiedź kształtu OpenAI Images niosącą bajty
// obrazu — postać, w której obraz przychodzi gotowy i nie wymaga drugiego
// pobrania.
func odpowiedzZBajtami(t *testing.T, obraz []byte) http.HandlerFunc {
	t.Helper()

	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": []any{map[string]any{"b64_json": wBase64(obraz)}},
		})
	}
}

// wpiszKanalObrazowy zakłada wiersz rejestru kanałów wskazujący podany adres
// i oddaje identyfikator kanału.
func wpiszKanalObrazowy(t *testing.T, zmontowany *Zmontowany, zycie context.Context, adres string) string {
	t.Helper()

	t.Setenv(zmiennaKluczaObrazow, "klucz-sprawdzianu")

	var wynik shared.ChannelAddResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandChannelAdd, shared.ChannelAddRequest{
		Name:  "kanał obrazowy sprawdzianu",
		Kind:  "api",
		Model: wskaznik("model-obrazow"),
		Config: jsonSurowy(t, map[string]any{
			"adapter":       "obrazy",
			"base_url":      adres,
			"credentialRef": zmiennaKluczaObrazow,
		}),
	}, &wynik)

	if wynik.Channel.Id == "" {
		t.Fatal("channel.add nie oddał identyfikatora kanału")
	}
	return wynik.Channel.Id
}

// TestGenerowanieZasobuOddajeZasobZBajtamiObrazuWMagazynie sprawdza, czy po
// udanym `design.asset.generate` odwołanie zasobu prowadzi do pliku, którego
// bajty są dokładnie tymi, które oddał kanał obrazowy.
func TestGenerowanieZasobuOddajeZasobZBajtamiObrazuWMagazynie(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)

	obraz := obrazPNG(t, 12, 7)
	serwer := serwerObrazow(t, odpowiedzZBajtami(t, obraz))
	kanal := wpiszKanalObrazowy(t, zmontowany, zycie, serwer.URL)

	var wynik shared.DesignAssetGenerateResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignAssetGenerate,
		shared.DesignAssetGenerateRequest{
			WindowId:  "okno-designu",
			Prompt:    shared.DesignPrompt{Subject: "kadr próbny sprawdzianu"},
			ChannelId: &kanal,
		}, &wynik)

	if len(wynik.Assets) != 1 {
		t.Fatalf("generowanie oddało %d zasobów, zamówiono 1", len(wynik.Assets))
	}
	zasob := wynik.Assets[0]
	if zasob.Kind != shared.DesignAssetKindImage {
		t.Errorf("zasób ma rodzaj %q, oczekiwany %q", zasob.Kind, shared.DesignAssetKindImage)
	}
	if zasob.Uri == nil {
		t.Fatal("zasób bez odwołania do treści — kafelek, za którym nic nie leży")
	}

	bajty := bajtyPodOdwolaniem(t, katalog, *zasob.Uri)
	if !bytes.Equal(bajty, obraz) {
		t.Errorf("treść pod odwołaniem ma %d bajtów i nie jest obrazem, który oddał kanał (%d bajtów)",
			len(bajty), len(obraz))
	}
	// Nazwą bloba jest suma jego zawartości, więc rozjazd znaczy inną treść
	// niż zmierzoną przy zapisie.
	if suma := sumaSha256(obraz); !bytes.Contains([]byte(*zasob.Uri), []byte(suma)) {
		t.Errorf("odwołanie %q nie niesie sumy kontrolnej utrwalonej treści (%s)", *zasob.Uri, suma)
	}
	if zasob.Format == nil || *zasob.Format != "png" {
		t.Errorf("format zasobu %v — nagłówek utrwalonego pliku mówi png", zasob.Format)
	}
	if zasob.Width == nil || *zasob.Width != 12 || zasob.Height == nil || *zasob.Height != 7 {
		t.Errorf("wymiary zasobu %v×%v — obraz w magazynie ma 12×7", zasob.Width, zasob.Height)
	}
}

// TestWykazZasobowOddajeOdwolaniaDoBajtowKazdegoZasobu sprawdza, czy każdy
// zasób oddany przez `design.asset.list` ma pod swoim odwołaniem rzeczywiste
// bajty, niezależnie od drogi, którą wszedł do modułu.
func TestWykazZasobowOddajeOdwolaniaDoBajtowKazdegoZasobu(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)

	obraz := obrazPNG(t, 9, 9)
	serwer := serwerObrazow(t, odpowiedzZBajtami(t, obraz))
	kanal := wpiszKanalObrazowy(t, zmontowany, zycie, serwer.URL)

	var zGenerowania shared.DesignAssetGenerateResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignAssetGenerate,
		shared.DesignAssetGenerateRequest{
			WindowId:  "okno-wykazu",
			Prompt:    shared.DesignPrompt{Subject: "zasób z generowania"},
			ChannelId: &kanal,
		}, &zGenerowania)

	// Druga droga zasobu do modułu: wniesienie przez operatora, sprawdzane
	// tak samo jak generowanie.
	wniesiony := obrazPNG(t, 4, 5)
	var zWniesienia shared.DesignAssetUploadResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignAssetUpload,
		shared.DesignAssetUploadRequest{
			WindowId:      "okno-wykazu",
			Name:          wskaznik("wniesiony ręcznie"),
			Kind:          shared.DesignAssetKindImage,
			ContentBase64: wskaznik(wBase64(wniesiony)),
		}, &zWniesienia)

	var wykaz shared.DesignAssetListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignAssetList,
		shared.DesignAssetListRequest{WindowId: wskaznik("okno-wykazu")}, &wykaz)

	if len(wykaz.Assets) != 2 {
		t.Fatalf("wykaz oddał %d zasobów, założono 2", len(wykaz.Assets))
	}
	if wykaz.Total == nil || *wykaz.Total != 2 {
		t.Errorf("wykaz podaje łącznie %v zasobów, w oknie leżą 2", wykaz.Total)
	}

	// Sprawdzenie porównuje zbiory treści, nie liczbę wierszy w wykazie.
	oczekiwane := map[string]bool{sumaSha256(obraz): false, sumaSha256(wniesiony): false}
	for _, zasob := range wykaz.Assets {
		if zasob.Uri == nil {
			t.Fatalf("zasób %s w wykazie bez odwołania do treści", zasob.Id)
		}
		suma := sumaSha256(bajtyPodOdwolaniem(t, katalog, *zasob.Uri))
		widziana, znana := oczekiwane[suma]
		if !znana {
			t.Errorf("zasób %s wskazuje treść, której do modułu nie wnoszono (suma %s)", zasob.Id, suma)
			continue
		}
		if widziana {
			t.Errorf("dwa zasoby wykazu wskazują tę samą treść (suma %s)", suma)
		}
		oczekiwane[suma] = true
	}
	for suma, widziana := range oczekiwane {
		if !widziana {
			t.Errorf("treść o sumie %s weszła do modułu, ale nie ma jej pod żadnym odwołaniem wykazu", suma)
		}
	}
}

// TestGenerowanieZasobuZAdresuWciagaBajtyZamiastZapisacOdsylacz sprawdza, czy
// odpowiedź niosąca odsyłacz do obrazu zamiast jego bajtów kończy się
// wciągnięciem treści spod odsyłacza, a nie zapisaniem samego odsyłacza.
func TestGenerowanieZasobuZAdresuWciagaBajtyZamiastZapisacOdsylacz(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)

	obraz := obrazPNG(t, 6, 6)
	skladnica := serwerObrazow(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		_, _ = w.Write(obraz)
	})
	adresObrazu := skladnica.URL + "/wynik.png"

	serwer := serwerObrazow(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": []any{map[string]any{"url": adresObrazu}},
		})
	})
	kanal := wpiszKanalObrazowy(t, zmontowany, zycie, serwer.URL)

	var wynik shared.DesignAssetGenerateResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignAssetGenerate,
		shared.DesignAssetGenerateRequest{
			WindowId:  "okno-adresu",
			Prompt:    shared.DesignPrompt{Subject: "obraz oddany odsyłaczem"},
			ChannelId: &kanal,
		}, &wynik)

	if len(wynik.Assets) != 1 {
		t.Fatalf("generowanie oddało %d zasobów, zamówiono 1", len(wynik.Assets))
	}
	zasob := wynik.Assets[0]
	if zasob.Uri == nil {
		t.Fatal("zasób bez odwołania do treści")
	}
	if *zasob.Uri == adresObrazu {
		t.Fatal("zasób trzyma odsyłacz dostawcy jako swoją treść — po wygaśnięciu odsyłacza zostanie bez bajtów")
	}
	if !bytes.Equal(bajtyPodOdwolaniem(t, katalog, *zasob.Uri), obraz) {
		t.Error("treść pod odwołaniem nie jest obrazem spod adresu — bajty nie zostały wciągnięte")
	}
}

// TestGenerowanieBezKanaluObrazowegoOdmawiaZamiastZalozycKafelek sprawdza, czy
// generowanie bez kanału obrazowego kończy się odmową i czy okno po niej
// zostaje puste, potwierdzone osobnym odczytem `design.asset.list`.
func TestGenerowanieBezKanaluObrazowegoOdmawiaZamiastZalozycKafelek(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	odmowa := wykonajOdmowna(t, zmontowany, zycie, shared.CommandDesignAssetGenerate,
		shared.DesignAssetGenerateRequest{
			WindowId: "okno-bez-kanalu",
			Prompt:   shared.DesignPrompt{Subject: "obraz, którego nie ma czym wytworzyć"},
		})
	if odmowa.Code != shared.ErrorCodeChannelUnavailable {
		t.Errorf("odmowa niesie kod %q, brak kanału obrazowego to %q",
			odmowa.Code, shared.ErrorCodeChannelUnavailable)
	}
	// Praca Prompt Buildera ma przeżyć odmowę — kontrakt niesie ją w `details`.
	var szczegoly szczegolyOdmowyGenerowania
	if len(odmowa.Details) == 0 {
		t.Fatal("odmowa bez szczegółów — gotowa treść polecenia ginie razem z nią")
	}
	if err := json.Unmarshal(odmowa.Details, &szczegoly); err != nil {
		t.Fatalf("szczegóły odmowy nieczytelne: %v", err)
	}
	if szczegoly.Polecenie == "" {
		t.Error("odmowa nie oddaje złożonego polecenia — Operator traci pracę Prompt Buildera")
	}
	if len(szczegoly.Brakujace) == 0 {
		t.Error("odmowa nie wymienia braków maszynowo — okno nie ma czego wypisać")
	}

	var wykaz shared.DesignAssetListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignAssetList,
		shared.DesignAssetListRequest{WindowId: wskaznik("okno-bez-kanalu")}, &wykaz)
	if len(wykaz.Assets) != 0 {
		t.Errorf("po odmowie generowania w oknie leży %d zasobów — kafelek bez bajtów", len(wykaz.Assets))
	}
}

// TestGenerowanieBezObrazuWOdpowiedziNieZostawiaZasobu sprawdza, czy odpowiedź
// kanału bez treści obrazu kończy się odmową, a wykaz zasobów w oknie zostaje
// pusty.
func TestGenerowanieBezObrazuWOdpowiedziNieZostawiaZasobu(t *testing.T) {
	przypadki := map[string]http.HandlerFunc{
		"odpowiedź bez pola obrazu": func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{"data": []any{map[string]any{}}})
		},
		"obraz zerowej długości": func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": []any{map[string]any{"b64_json": ""}},
			})
		},
		"treść niebędąca base64": func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": []any{map[string]any{"b64_json": "to nie jest base64 !!!"}},
			})
		},
	}

	for nazwa, obsluga := range przypadki {
		t.Run(nazwa, func(t *testing.T) {
			zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)
			serwer := serwerObrazow(t, obsluga)
			kanal := wpiszKanalObrazowy(t, zmontowany, zycie, serwer.URL)

			wykonajOdmowna(t, zmontowany, zycie, shared.CommandDesignAssetGenerate,
				shared.DesignAssetGenerateRequest{
					WindowId:  "okno-bez-obrazu",
					Prompt:    shared.DesignPrompt{Subject: "kanał milczy o obrazie"},
					ChannelId: &kanal,
				})

			var wykaz shared.DesignAssetListResponse
			wykonajUdana(t, zmontowany, zycie, shared.CommandDesignAssetList,
				shared.DesignAssetListRequest{WindowId: wskaznik("okno-bez-obrazu")}, &wykaz)
			if len(wykaz.Assets) != 0 {
				t.Errorf("po odmowie w oknie leży %d zasobów — wiersz powstał bez bajtów", len(wykaz.Assets))
			}
		})
	}
}

// TestWniesienieZasobuBezTresciNieZakladaWiersza domyka drugą drogę zasobu do
// modułu. Wniesienie bez treści i bez ścieżki nie ma skąd wziąć bajtów, więc ma
// odmówić — a nie założyć zasób „do uzupełnienia".
func TestWniesienieZasobuBezTresciNieZakladaWiersza(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	przypadki := map[string]shared.DesignAssetUploadRequest{
		"bez treści i bez ścieżki": {
			WindowId: "okno-wniesienia", Kind: shared.DesignAssetKindImage,
		},
		"treść pusta": {
			WindowId: "okno-wniesienia", Kind: shared.DesignAssetKindImage,
			ContentBase64: wskaznik(""),
		},
		"ścieżka wskazująca nieistniejący plik": {
			WindowId: "okno-wniesienia", Kind: shared.DesignAssetKindImage,
			SourcePath: wskaznik(t.TempDir() + "/nie-ma-takiego-pliku.png"),
		},
	}
	for nazwa, zadanie := range przypadki {
		t.Run(nazwa, func(t *testing.T) {
			wykonajOdmowna(t, zmontowany, zycie, shared.CommandDesignAssetUpload, zadanie)
		})
	}

	var wykaz shared.DesignAssetListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignAssetList,
		shared.DesignAssetListRequest{WindowId: wskaznik("okno-wniesienia")}, &wykaz)
	if len(wykaz.Assets) != 0 {
		t.Errorf("po odmowach wniesienia w oknie leży %d zasobów", len(wykaz.Assets))
	}
}

// TestWniesienieZasobuSciezkaZamrazaTrescWMagazynie sprawdza, czy zasób
// wniesiony ścieżką zachowuje treść z chwili wniesienia także po nadpisaniu
// pliku źródłowego.
func TestWniesienieZasobuSciezkaZamrazaTrescWMagazynie(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)

	pierwotny := obrazPNG(t, 8, 3)
	sciezka := t.TempDir() + "/zrodlo.png"
	zapiszPlikSprawdzianu(t, sciezka, pierwotny)

	var wniesiony shared.DesignAssetUploadResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignAssetUpload,
		shared.DesignAssetUploadRequest{
			WindowId:   "okno-sciezki",
			Kind:       shared.DesignAssetKindImage,
			SourcePath: &sciezka,
		}, &wniesiony)

	if wniesiony.Asset.Uri == nil {
		t.Fatal("zasób wniesiony ścieżką nie ma odwołania do treści")
	}
	// Nadpisanie pliku źródłowego ujawniłoby wskaźnik zamiast zamrożonej
	// treści zmianą zasobu.
	zapiszPlikSprawdzianu(t, sciezka, obrazPNG(t, 40, 40))

	if !bytes.Equal(bajtyPodOdwolaniem(t, katalog, *wniesiony.Asset.Uri), pierwotny) {
		t.Error("treść zasobu poszła za nadpisanym plikiem źródłowym — zasób jest wskaźnikiem, nie treścią zamrożoną")
	}
}

// TestRodzajSpozaKontraktuWracaJakoPomylkaWolajacego sprawdza, czy rodzaj
// zasobu spoza wartości kontraktu wraca kodem pomyłki wołającego, a nie
// ponawialną awarią rdzenia, i czy po odmowie okno zostaje bez zasobu.
func TestRodzajSpozaKontraktuWracaJakoPomylkaWolajacego(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	tresc := obrazPNG(t, 6, 4)
	blad := wykonajOdmowna(t, zmontowany, zycie, shared.CommandDesignAssetUpload,
		shared.DesignAssetUploadRequest{
			WindowId:      "okno-rodzaju",
			Kind:          shared.DesignAssetKind("rzezba"),
			ContentBase64: wskaznik(wBase64(tresc)),
		})

	if blad.Code != shared.ErrorCodeValidationFailed {
		t.Errorf("odmowa niesie kod %q, a pomyłka wołającego ma wracać jako %q — "+
			"kod awarii rdzenia jest ponawialny, więc klient z pętlą powtarza żądanie bez końca",
			blad.Code, shared.ErrorCodeValidationFailed)
	}
	if strings.Contains(blad.Message, "CHECK constraint") || strings.Contains(blad.Message, "rodzaj IN") {
		t.Errorf("odmowa cytuje warunek schematu zamiast nazwać brak: %s", blad.Message)
	}
	for _, dopuszczalny := range shared.WartosciDesignAssetKind() {
		if !strings.Contains(blad.Message, string(dopuszczalny)) {
			t.Errorf("odmowa nie wymienia rodzaju dopuszczalnego %q — "+
				"Operator nie dowie się, czym zastąpić swój: %s", dopuszczalny, blad.Message)
		}
	}

	var wykaz shared.DesignAssetListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignAssetList,
		shared.DesignAssetListRequest{WindowId: wskaznik("okno-rodzaju")}, &wykaz)
	if len(wykaz.Assets) != 0 {
		t.Errorf("po odmowie w oknie leży %d zasobów — sprawdzenie rodzaju stoi po zapisie, "+
			"a ma stać przed nim", len(wykaz.Assets))
	}
}
