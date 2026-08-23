// Odpowiedzialność pliku: wytworzenie zasobu (`design.asset.generate`) — droga
// od promptu strukturalnego do zasobu, za którym leżą prawdziwe bajty obrazu.
// Metody stoją na `*adapterDesignu` (`adapter_modul_design.go`); plik osobny
// wedle odpowiedzialności, jak wniesienie w `_wgranie.go`.
//
// Ciąg jest jeden i nierozerwalny: złóż polecenie → wyślij je kanałem obrazowym
// (`models/adapter_obrazy.go`) → odbierz fragment `image` → odłóż bajty
// w magazynie (`adapter_modul_design_wgranie.go`, blob pod sumą sha256) →
// zmierz format i wymiary z nagłówka utrwalonego pliku → załóż wiersz zasobu
// z `uri`, `format`, `width`, `height`. Wiersz powstaje wyłącznie po utrwaleniu
// bajtów — ta sama kolejność co przy wniesieniu i z tego samego powodu: kafelek
// w Assets Panelu, za którym nic nie leży, jest kłamstwem koperty.
//
// Każdy brak jest odmową nazywającą brak, nigdy obrazem zastępczym. Brak kanału
// obrazowego w rejestrze, kanał nieczynny, kanał wskazany a tekstowy, brak
// poświadczenia (odmawia sam kanał, zdaniem z `models/adapter_obrazy.go`),
// odpowiedź bez obrazu, bajty nie do pobrania — wszystko to kończy komendę
// błędem. Nie ma tu ani jednej drogi, którą wracałby `status: ok` bez treści.
//
// Wykaz braków w odmowie jest prawdziwy w chwili odczytu: każda odmowa niesie
// własny wykaz tego, czego zabrakło w tym wywołaniu, obok gotowej treści
// polecenia w `details.polecenie`.
//
// Adres zamiast bajtów też ląduje w magazynie. Dostawcy zgodni z OpenAI Images
// oddają albo `b64_json`, albo `url` (i `url` bywa domyślne). Odsyłacz dostawcy
// wygasa, więc zapisanie go jako `uri` zasobu dałoby zasób, który po godzinie
// przestaje mieć treść. Bajty spod adresu wciągamy tutaj i dopiero one idą do
// magazynu; niepowodzenie pobrania jest odmową, nie zasobem bez treści.
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
	// limitPobraniaObrazu jest górną granicą bajtów wciąganych spod adresu
	// dostawcy. Nie jest to polityka jakości obrazu, tylko obrona rdzenia przed
	// odpowiedzią bez końca: bez granicy `io.ReadAll` na cudzym adresie wciąga
	// tyle, ile druga strona zechce nadać, aż do wyczerpania pamięci procesu.
	limitPobraniaObrazu = 64 << 20
	// limitCzasuPobraniaObrazu ogranicza samo pobranie spod adresu; czas
	// generowania wyznacza kanał własnym parametrem `limit_sekund`.
	limitCzasuPobraniaObrazu = 120 * time.Second
)

// szczegolyOdmowyGenerowania niesie w `error.details` to, co ocalało z pracy
// Prompt Buildera: gotową treść polecenia i liczbę wariantów, o którą Operator
// prosił. Pole `details` kontrakt opisuje jako dane diagnostyczne, a treść,
// której rdzeń nie miał komu podać, jest dokładnie tym.
type szczegolyOdmowyGenerowania struct {
	// Polecenie jest promptem strukturalnym złożonym w jeden tekst.
	Polecenie string `json:"polecenie"`
	// Wariantow niesie liczbę zamówionych wariantów.
	Wariantow int `json:"wariantow"`
	// Brakujace wymienia braki maszynowo, żeby okno mogło je wypisać bez
	// rozbierania zdania odmowy na kawałki. Wykaz składa się przy każdej
	// odmowie osobno i niesie wyłącznie braki tego wywołania.
	Brakujace []string `json:"brakujace"`
}

// bladOdmowyGenerowania składa odmowę `design.asset.generate` wraz ze
// szczegółami ocalonymi z promptu.
//
// Nieudane złożenie `details` nie zmienia odmowy — Operator ma dostać zdanie
// o braku, a nie usterkę serializacji podstawioną w jego miejsce.
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

// GenerujZasob wytwarza zasoby wizualne kanałem obrazowym — obsługuje
// `design.asset.generate`.
//
// Kolejność sprawdzeń ma znaczenie: braki żądania (okno, temat) idą pierwsze,
// bo Operator poprawia je sam; dopiero prompt kompletny wchodzi na drogę,
// na której mogą zabraknąć kanał, poświadczenie albo magazyn.
//
// Wariant nieudany przerywa całość. Gdyby drugi wariant padł, a pierwszy został,
// odpowiedź niosłaby mniej zasobów, niż Operator zamówił, bez słowa o tym, czemu
// — a `assets` nie ma pola na „ten się nie udał". Zasoby, które zdążyły powstać,
// zostają w bazie i w magazynie: ich bajty są prawdziwe, więc kasowanie ich
// byłoby niszczeniem cudzej treści z powodu, który jej nie dotyczy.
func (a *adapterDesignu) GenerujZasob(ctx context.Context,
	z shared.DesignAssetGenerateRequest) (shared.DesignAssetGenerateResponse, error) {

	if z.WindowId == "" {
		return shared.DesignAssetGenerateResponse{}, bladWskazaniaDesignu("komenda bez okna")
	}
	if z.Prompt.Subject == "" {
		return shared.DesignAssetGenerateResponse{}, bladWskazaniaDesignu("prompt bez tematu")
	}

	// Brak wskazania liczby wariantów znaczy jeden wynik, nie zero — tak czytał
	// to pole Prompt Builder od początku.
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

	// Prompt utrwala się RAZ na całe wywołanie, przed pierwszym wariantem:
	// wszystkie warianty powstają z tego samego polecenia, więc jeden wiersz
	// promptu jest o nich jedną prawdą (`migracja_048_design.sql`). Wiersz
	// niesie okno i kanał, bo `design.prompt.history.list` czyta prompty okna
	// i pyta o kanał, którym poszły.
	//
	// Niepowodzenie zapisu promptu NIE przerywa generowania. Prowenancja jest
	// wiedzą o zasobie, a nie zasobem: odmowa wytworzenia obrazu dlatego, że nie
	// dało się zapisać jego opisu, zabierałaby Operatorowi rzecz, o którą prosił,
	// z powodu, który jej nie dotyczy. Zasób wychodzi wtedy bez `promptId` —
	// czyli mówi prawdę o tym, czego rdzeń o nim nie wie.
	var promptID *int64
	kodKanalu := kanal.Identyfikator()
	if zapisany, err := a.repozytorium.ZapiszPrompt(ctx, promptDoZapisuDesignu(z.Prompt,
		z.WindowId, kodKanalu)); err == nil {
		promptID = &zapisany.ID
	}

	polecenie := zlozPolecenieObrazu(z.Prompt)
	rodzaj := shared.DesignAssetKindImage
	if z.Kind != nil && strings.TrimSpace(string(*z.Kind)) != "" {
		// Ta sama droga i ten sam powód, co przy wnoszeniu: rodzaj spoza
		// kontraktu odbiłby się od warunku schematu i wrócił jako awaria rdzenia
		// oznaczona jako ponawialna. Sprawdzenie stoi przed wytworzeniem obrazu,
		// żeby odmowa nie kosztowała wywołania kanału.
		if err := sprawdzRodzajZasobu("design.asset.generate", *z.Kind); err != nil {
			return shared.DesignAssetGenerateResponse{}, err
		}
		rodzaj = string(*z.Kind)
	}

	zasoby := make([]shared.DesignAsset, 0, wariantow)
	for numer := 1; numer <= wariantow; numer++ {
		// Warianty są oddzielnymi wywołaniami kanału. Parametr `n` kanału
		// obrazowego (`liczba`) należy do wiersza rejestru i opisuje wolę
		// Operatora co do kanału, a nie co do tego jednego promptu; poza tym
		// fragment obrazu niesie jeden obraz, więc drugiego nie ma skąd wziąć.
		// Bez obrazów wejściowych: `design.asset.generate` tworzy obraz OD ZERA
		// z samego polecenia — materiał wejściowy zrobiłby z niego edycję.
		bajty, typTresci, err := a.wytworzObraz(ctx, kanal, z.WindowId, polecenie, nil)
		if err != nil {
			return shared.DesignAssetGenerateResponse{}, err
		}
		// Pierwszy wariant jest pniem, kolejne wskazują na niego
		// `variantOfAssetId` — to jedyne powiązanie, które kontrakt tu zna,
		// i jest prawdziwe: wszystkie powstały z tego samego polecenia.
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

// zalozZasobZBajtow utrwala bajty jednego wariantu i zakłada jego wiersz.
//
// Format i wymiary mierzymy z tego, co leży w magazynie — tą samą funkcją co
// przy wniesieniu (`rozpoznajObrazZasobu`), więc zasób wygenerowany i zasób
// wniesiony opisują się jedną miarą. Typ treści z fragmentu wchodzi
// dopiero tam, gdzie nagłówek zamilkł: format spoza trójki PNG/JPEG/GIF (WEBP,
// SVG) zostawia wymiary puste, ale nazwa formatu z `mimeType` jest wtedy jedyną
// prawdziwą wiadomością o pliku, jaką mamy.
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
	// Etykiet generowanie nie nadaje — kontrakt `design.asset.generate` nie ma
	// pola `tags`, a wymyślenie etykiety byłoby dopisaniem cechy, o którą nikt
	// nie prosił. Odczyt zostaje, bo `zasobKontraktu` żąda wykazu, a pusty
	// wykaz jest tu prawdą, nie zaniechaniem.
	etykiety, err := a.repozytorium.EtykietyZasobu(ctx, zapisany.ID)
	if err != nil {
		return shared.DesignAsset{}, bladDesignu(err)
	}
	return zasobKontraktu(zapisany, etykiety), nil
}

// promptDoZapisuDesignu przekłada prompt kontraktu na wiersz `prompt_design`
// wraz z oknem i kanałem wydania.
//
// Identyfikator z żądania (`DesignPrompt.Id`) ma pierwszeństwo: Prompt Builder
// wysyła nim ten sam prompt przy regeneracji wariantów, a zapis po tym
// identyfikatorze nadpisuje wiersz zastany, zamiast zakładać drugi obok.
// Historia pokazuje wtedy jedno wydanie polecenia, a nie tyle wpisów, ile razy
// Operator kliknął „Generuj".
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

// wytworzObraz wykonuje jedno wywołanie kanału obrazowego i oddaje bajty obrazu
// wraz z typem treści.
//
// Zbieramy wyłącznie fragment `image`. Prowenancja i metadane konta jadą przed
// treścią i nie są obrazem; fragment tekstowy też obrazem nie jest — zebranie go
// dałoby „obraz", którego bajty są zdaniem po polsku. Dlatego pusty wynik pętli
// jest odmową, a nie zasobem.
// Obrazy WEJŚCIOWE są tu opcjonalne i to jedna droga dla dwóch czynności:
// generowanie od zera podaje je puste (`design.asset.generate`), a warsztat
// fotografii — materiałem i maską (`design.photo.*`). Druga funkcja wołająca ten
// sam kanał inaczej dałaby dwie prawdy o tym, jak rdzeń rozmawia z silnikiem
// obrazów; kształt żądania edycji składa warstwa modeli, nie ten moduł.
func (a *adapterDesignu) wytworzObraz(ctx context.Context,
	kanal models.Definicja, okno, polecenie string,
	obrazy []models.ObrazWejsciowy) ([]byte, string, error) {

	var tresc models.TrescObrazu
	ujscie := models.UjscieFunkcji(func(_ context.Context, f models.Fragment) error {
		if f.Kind != shared.ChunkKindImage || len(f.Data) == 0 {
			return nil
		}
		// Pierwszy obraz wygrywa: fragment niesie jeden obraz, a drugi (gdyby
		// kanał go nadał) należałby do innego wariantu niż ten zamówiony.
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
	// Błąd kanału (brak poświadczenia, odmowa dostawcy, odpowiedź bez obrazu)
	// jedzie do Operatora ze zdaniem kanału — kanał wie o swoim braku więcej
	// niż ten moduł i nazywa go dokładniej, niż zrobiłoby to zdanie ogólne.
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

// pobierzObrazSpodAdresu wciąga bajty spod odsyłacza wystawionego przez
// dostawcę — patrz nagłówek pliku: odsyłacz wygasa, treść zasobu nie ma prawa.
// Typ treści bierzemy z nagłówka odpowiedzi, gdy fragment go nie niósł.
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

// nazwaWariantu składa etykietę zasobu z treści polecenia. Nazwa jest opisem,
// nie wynikiem generowania: obraz leży w magazynie, a nazwa go tylko podpisuje.
// Numer dopisujemy wyłącznie wtedy, gdy wariantów jest więcej niż jeden, bo
// „wariant 1 z 1" nie jest żadnym rozróżnieniem.
func nazwaWariantu(polecenie string, numer, wariantow int) string {
	nazwa := strings.TrimSpace(polecenie)
	// Granica liczona w znakach, nie w bajtach: polecenia są po polsku, a cięcie
	// bajtowe rozłupałoby „ł" na pół i wpisało do bazy nazwę, która nie jest
	// poprawnym UTF-8 — panel pokazałby wtedy znak zastępczy zamiast litery.
	const granica = 120
	if znaki := []rune(nazwa); len(znaki) > granica {
		nazwa = strings.TrimSpace(string(znaki[:granica])) + "…"
	}
	if wariantow > 1 {
		nazwa = fmt.Sprintf("%s (wariant %d z %d)", nazwa, numer, wariantow)
	}
	return nazwa
}

// formatZTypuTresci wyciąga nazwę formatu z typu treści („image/png" → „png").
// Wartość spoza rodziny `image/` nie daje formatu: „application/json" nie jest
// formatem obrazu, a wpisany w kolumnę `format` wyglądałby jak zmierzony.
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
