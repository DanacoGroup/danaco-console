// Plik obsługuje design.asset.generate: składa polecenie z promptu, wywołuje
// kanał obrazowy, utrwala bajty obrazu w magazynie treści i dopiero potem
// zakłada wiersz zasobu. Każdy brak kończy komendę odmową, nigdy obrazem
// zastępczym.
package core

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/models"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

const (
	// limitPobraniaObrazu ogranicza bajty wciągane spod adresu dostawcy, żeby
	// odpowiedź bez końca nie wyczerpała pamięci procesu podczas pobierania obrazu.
	limitPobraniaObrazu = 64 << 20
	// limitCzasuPobraniaObrazu ogranicza czas pobrania obrazu spod adresu
	// dostawcy; czas samego generowania ustala kanał obrazowy parametrem limit_sekund.
	limitCzasuPobraniaObrazu = 120 * time.Second
)

// szczegolyOdmowyGenerowania niesie w error.details treść złożonego polecenia,
// liczbę zamówionych wariantów i wykaz braków tego wywołania jako dane
// diagnostyczne dla Operatora.
type szczegolyOdmowyGenerowania struct {
	// Polecenie jest promptem strukturalnym złożonym w jeden tekst.
	Polecenie string `json:"polecenie"`
	// Wariantow niesie liczbę zamówionych wariantów.
	Wariantow int `json:"wariantow"`
	// Brakujace wymienia braki tego wywołania, żeby okno mogło je wypisać bez
	// rozbioru zdania odmowy.
	Brakujace []string `json:"brakujace"`
}

// bladOdmowyGenerowania składa odmowę design.asset.generate wraz ze
// szczegółami ocalonymi z promptu; niepowodzenie zapisu szczegółów nie
// zmienia treści odmowy.
func bladOdmowyGenerowania(kod shared.ErrorCode, tresc string,
	p shared.DesignPrompt, wariantow int, brakujace ...string) error {

	blad := protocol.NowyBlad(kod, "moduł Design: "+tresc)
	szczegoly, err := json.Marshal(szczegolyOdmowyGenerowania{
		Polecenie: zlozPolecenieObrazu(p),
		Wariantow: wariantow,
		Brakujace: brakujace,
	})
	if err == nil {
		blad.Details = szczegoly
	}
	return protocol.JakoError(blad)
}

// GenerujZasob obsługuje design.asset.generate: sprawdza żądanie, wywołuje
// kanał obrazowy i zakłada wiersz zasobu dla każdego wariantu. Wariant
// nieudany przerywa całość; zasoby już utrwalone zostają w bazie i magazynie.
func (a *adapterDesignu) GenerujZasob(ctx context.Context,
	z shared.DesignAssetGenerateRequest) (shared.DesignAssetGenerateResponse, error) {

	if z.WindowId == "" {
		return shared.DesignAssetGenerateResponse{}, bladWskazaniaDesignu("komenda bez okna")
	}
	if z.Prompt.Subject == "" {
		return shared.DesignAssetGenerateResponse{}, bladWskazaniaDesignu("prompt bez tematu")
	}

	// Brak wskazania liczby wariantów znaczy jeden wynik, nie zero.
	wariantow := 1
	if z.Prompt.Variants != nil && *z.Prompt.Variants > 0 {
		wariantow = *z.Prompt.Variants
	}

	kanal, err := a.kanalObrazowyZadania(z.ChannelId, z.Prompt, wariantow)
	if err != nil {
		return shared.DesignAssetGenerateResponse{}, err
	}
	if a.magazyn == nil {
		return shared.DesignAssetGenerateResponse{}, bladOdmowyGenerowania(
			shared.ErrorCodeInternalError,
			"magazyn treści zasobów nie jest wpięty — nie ma gdzie odłożyć bajtów obrazu, "+
				"a zasób bez bajtów nie powstanie",
			z.Prompt, wariantow, "magazyn treści zasobów wizualnych")
	}

	// Prompt utrwala się raz na całe wywołanie: wszystkie warianty dzielą jedno polecenie.

	// Niepowodzenie zapisu promptu nie przerywa generowania; zasób wychodzi bez promptId.
	var promptID *int64
	kodKanalu := kanal.Identyfikator()
	if zapisany, err := a.repozytorium.ZapiszPrompt(ctx, promptDoZapisuDesignu(z.Prompt,
		z.WindowId, kodKanalu)); err == nil {
		promptID = &zapisany.ID
	}

	polecenie := zlozPolecenieObrazu(z.Prompt)
	rodzaj := shared.DesignAssetKindImage
	if z.Kind != nil && strings.TrimSpace(string(*z.Kind)) != "" {
		// Sprawdzenie rodzaju stoi przed wywołaniem kanału, żeby odmowa nie
		// kosztowała generowania.
		if err := sprawdzRodzajZasobu("design.asset.generate", *z.Kind); err != nil {
			return shared.DesignAssetGenerateResponse{}, err
		}
		rodzaj = string(*z.Kind)
	}

	zasoby := make([]shared.DesignAsset, 0, wariantow)
	for numer := 1; numer <= wariantow; numer++ {
		// Każdy wariant jest osobnym wywołaniem kanału; fragment niesie jeden
		// obraz na wywołanie.
		bajty, typTresci, err := a.wytworzObraz(ctx, kanal, z.WindowId, polecenie, nil)
		if err != nil {
			return shared.DesignAssetGenerateResponse{}, err
		}
		// Pierwszy wariant jest pniem, kolejne wskazują na niego przez variantOfAssetId.
		var pien *string
		if len(zasoby) > 0 {
			pien = &zasoby[0].Id
		}
		zasob, err := a.zalozZasobZBajtow(ctx, z, rodzaj, polecenie, numer, wariantow, bajty, typTresci, pien, promptID)
		if err != nil {
			return shared.DesignAssetGenerateResponse{}, err
		}
		zasoby = append(zasoby, zasob)
	}
	return shared.DesignAssetGenerateResponse{Assets: zasoby}, nil
}

// zalozZasobZBajtow utrwala bajty jednego wariantu w magazynie treści, mierzy
// format i wymiary obrazu i zakłada jego wiersz zasobu w bazie danych.
func (a *adapterDesignu) zalozZasobZBajtow(ctx context.Context,
	z shared.DesignAssetGenerateRequest, rodzaj, polecenie string, numer, wariantow int,
	bajty []byte, typTresci string, pien *string, promptID *int64) (shared.DesignAsset, error) {

	suma := sha256.Sum256(bajty)
	odwolanie, err := a.magazyn.Zapisz(bajty, hex.EncodeToString(suma[:]))
	if err != nil {
		return shared.DesignAsset{}, bladZapisuZasobuDesignu(
			"nie można utrwalić bajtów wygenerowanego obrazu: " + err.Error())
	}

	format, szerokosc, wysokosc := rozpoznajObrazZasobu(odwolanie)
	if zTypu := formatZTypuTresci(typTresci); zTypu != nil {
		format = pierwszyTekst(format, zTypu)
	}
	nazwa := nazwaWariantu(polecenie, numer, wariantow)

	zapisany, err := a.repozytorium.ZapiszZasob(ctx, dane.ZasobDesignu{
		Kod:             nowyIdentyfikator(przedrostekZasobuDesign),
		Okno:            z.WindowId,
		Nazwa:           &nazwa,
		Rodzaj:          rodzaj,
		Format:          format,
		URI:             &odwolanie,
		PromptID:        promptID,
		WariantZasobuID: pien,
		Szerokosc:       szerokosc,
		Wysokosc:        wysokosc,
	})
	if err != nil {
		return shared.DesignAsset{}, bladDesignu(err)
	}
	// Generowanie nie nadaje etykiet; pusty wykaz jest tu prawdą, nie zaniechaniem.
	etykiety, err := a.repozytorium.EtykietyZasobu(ctx, zapisany.ID)
	if err != nil {
		return shared.DesignAsset{}, bladDesignu(err)
	}
	return zasobKontraktu(zapisany, etykiety), nil
}

// promptDoZapisuDesignu przekłada prompt kontraktu na wiersz prompt_design
// wraz z oknem i kanałem wydania; identyfikator żądania nadpisuje wiersz zastany.
func promptDoZapisuDesignu(p shared.DesignPrompt, okno, kanal string) dane.PromptDesignu {
	kod := nowyIdentyfikator(przedrostekPromptuDesign)
	if p.Id != nil && strings.TrimSpace(*p.Id) != "" {
		kod = strings.TrimSpace(*p.Id)
	}
	var kanalWiersza *string
	if strings.TrimSpace(kanal) != "" {
		wartosc := kanal
		kanalWiersza = &wartosc
	}
	return dane.PromptDesignu{
		Kod:         kod,
		Temat:       p.Subject,
		Styl:        p.Style,
		Kompozycja:  p.Composition,
		Oswietlenie: p.Lighting,
		Paleta:      p.Palette,
		Proporcje:   p.AspectRatio,
		Wykluczenia: p.Exclusions,
		Ziarno:      p.Seed,
		Warianty:    p.Variants,
		Silnik:      p.Engine,
		Kreatywnosc: p.Creativity,
		Okno:        okno,
		Kanal:       kanalWiersza,
	}
}

// wytworzObraz wykonuje jedno wywołanie kanału obrazowego i oddaje bajty
// obrazu wraz z typem treści, zbierając wyłącznie fragment obrazu z odpowiedzi
// kanału.
func (a *adapterDesignu) wytworzObraz(ctx context.Context,
	kanal models.Definicja, okno, polecenie string,
	obrazy []models.ObrazWejsciowy) ([]byte, string, error) {

	var tresc models.TrescObrazu
	ujscie := models.UjscieFunkcji(func(_ context.Context, f models.Fragment) error {
		if f.Kind != shared.ChunkKindImage || len(f.Data) == 0 {
			return nil
		}
		// Pierwszy obraz wygrywa: fragment niesie jeden obraz, należący do
		// zamówionego wariantu.
		if tresc.Base64 != "" || tresc.Adres != "" {
			return nil
		}
		return json.Unmarshal(f.Data, &tresc)
	})
	zapytanie := models.Zapytanie{
		Zasiegi:         models.Zasiegi{Okno: okno},
		Wiadomosc:       okno,
		Tresc:           polecenie,
		Kanal:           kanal.Identyfikator(),
		ObrazyWejsciowe: obrazy,
	}
	// Zdanie błędu pochodzi od kanału, który zna powód swojej odmowy
	// dokładniej niż ten moduł.
	if err := a.kanaly.Wyslij(ctx, zapytanie, ujscie); err != nil {
		return nil, "", protocol.JakoError(protocol.BladZeZrodla(
			shared.ErrorCodeChannelUnavailable,
			fmt.Errorf("moduł Design: kanał obrazowy nie oddał obrazu: %w", err)))
	}

	if tresc.Base64 == "" && tresc.Adres == "" {
		return nil, "", protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
			"moduł Design: kanał "+kanal.Identyfikator()+" zakończył wywołanie bez fragmentu obrazu — "+
				"zasób nie powstanie, bo nie ma z czego; naprawa: sprawdzić parametry ścieżek "+
				"odpowiedzi kanału (sciezka_base64, sciezka_adresu)"))
	}
	if tresc.Base64 != "" {
		bajty, err := base64.StdEncoding.DecodeString(tresc.Base64)
		if err != nil {
			return nil, "", protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
				"moduł Design: kanał "+kanal.Identyfikator()+" oddał treść, która nie jest poprawnym base64: "+err.Error()))
		}
		if len(bajty) == 0 {
			return nil, "", protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
				"moduł Design: kanał "+kanal.Identyfikator()+" oddał obraz o zerowej długości — "+
					"zasób bez bajtów nie ma czego pokazać"))
		}
		return bajty, tresc.TypTresci, nil
	}
	return pobierzObrazSpodAdresu(ctx, tresc.Adres, tresc.TypTresci)
}

// pobierzObrazSpodAdresu wciąga bajty obrazu spod odsyłacza dostawcy, ponieważ
// odsyłacz wygasa i sam nie może zostać zapisany jako trwała treść zasobu.
func pobierzObrazSpodAdresu(ctx context.Context, adres, typTresci string) ([]byte, string, error) {
	odmowa := func(powod string) error {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
			"moduł Design: kanał oddał adres obrazu zamiast bajtów, a bajtów nie da się wciągnąć ("+
				powod+") — zasób nie powstanie z samego odsyłacza, bo ten wygasa"))
	}
	zadanie, err := http.NewRequestWithContext(ctx, http.MethodGet, adres, nil)
	if err != nil {
		return nil, "", odmowa(err.Error())
	}
	klient := &http.Client{Timeout: limitCzasuPobraniaObrazu}
	odpowiedz, err := klient.Do(zadanie)
	if err != nil {
		return nil, "", odmowa(err.Error())
	}
	defer odpowiedz.Body.Close()
	if odpowiedz.StatusCode < 200 || odpowiedz.StatusCode > 299 {
		return nil, "", odmowa("odpowiedź " + odpowiedz.Status)
	}
	bajty, err := io.ReadAll(io.LimitReader(odpowiedz.Body, limitPobraniaObrazu))
	if err != nil {
		return nil, "", odmowa(err.Error())
	}
	if len(bajty) == 0 {
		return nil, "", odmowa("pusta treść pod adresem")
	}
	if strings.TrimSpace(typTresci) == "" {
		typTresci = odpowiedz.Header.Get("Content-Type")
	}
	return bajty, typTresci, nil
}

// nazwaWariantu składa etykietę zasobu z treści polecenia, skracając ją do stu
// dwudziestu znaków i dopisując numer wariantu, gdy wariantów jest więcej niż jeden.
func nazwaWariantu(polecenie string, numer, wariantow int) string {
	nazwa := strings.TrimSpace(polecenie)
	// Granica liczona w znakach, nie w bajtach, żeby cięcie nie rozłupało
	// litery polskiej.
	const granica = 120
	if znaki := []rune(nazwa); len(znaki) > granica {
		nazwa = strings.TrimSpace(string(znaki[:granica])) + "…"
	}
	if wariantow > 1 {
		nazwa = fmt.Sprintf("%s (wariant %d z %d)", nazwa, numer, wariantow)
	}
	return nazwa
}

// formatZTypuTresci wyciąga nazwę formatu z typu treści odpowiedzi kanału
// obrazowego; wartość spoza rodziny image/ nie jest formatem obrazu.
func formatZTypuTresci(typTresci string) *string {
	typTresci = strings.TrimSpace(strings.ToLower(typTresci))
	if typTresci == "" {
		return nil
	}
	if kropka := strings.IndexByte(typTresci, ';'); kropka >= 0 {
		typTresci = strings.TrimSpace(typTresci[:kropka])
	}
	if !strings.HasPrefix(typTresci, "image/") {
		return nil
	}
	format := strings.TrimPrefix(typTresci, "image/")
	if format == "" {
		return nil
	}
	return &format
}
