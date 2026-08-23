// Odpowiedzialność pliku: dwa wejścia do kolejki wczytywania, których nie ma
// w `adapter_modul_studio_cyfryzacja.go` — `studio.ingest.url` (strona sieciowa)
// i `studio.ingest.device.scan` (obraz z urządzenia).
//
// ── Strona sieciowa idzie `net/http`, nie przeglądarką ──────────────────────
// Pobranie strony to żądanie HTTP i odczyt odpowiedzi. Silnik przeglądarki
// wykonałby jeszcze skrypty strony — i byłby zależnością spoza instalki,
// a więc odmową u Operatora. Strona zbudowana wyłącznie skryptem zwróci tu
// mało treści; to jest cena znana i wybrana, a nie przeoczenie. Migawkę strony
// wykonanej skryptem oddaje `browser.snapshot.get`, co kontrakt mówi wprost
// w opisie tej komendy.
//
// ── Skaner idzie warstwą urządzeń systemu ───────────────────────────────────
// Skanowanie prowadzi warstwa urządzeń rozdzielona po systemie
// (`urzadzenia_skaner.go`): SANE na Linuksie, WIA przez PowerShell na Windowsie.
// Ta komenda nie ma własnej drogi do urządzenia i nie ma własnej odmowy —
// obie rzeczy należą do warstwy, bo inaczej rozjechałyby się z wykazem
// (`studio.ingest.device.list`), który pyta tę samą warstwę.
//
// Wersja natywna Windows jest tu rzeczą rozstrzygającą: instalka natywna nie
// niesie SANE i nigdy nie poniesie, więc odmowa „brak `scanimage`" byłaby na
// Windowsie odmową na zawsze. Dlatego rozstrzygnięcie po systemie stoi w
// warstwie, a nie w tej komendzie.
//
// Rdzeń NIE oddaje tu pustej kolejki pozycji, gdy czegoś brakuje: pusta
// kolejka jest twierdzeniem „skanowałem i nic nie przyszło", którego rdzeń bez
// odpowiedzi urządzenia nie ma prawa postawić. Każdy brak — warstwy, programu,
// urządzenia, sterownika — jest odmową nazywającą, czego brakuje.
package core

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

const (
	// granicaPobraniaStrony jest granicą czasu jednego pobrania. Strona, która
	// nie odpowiada w tym czasie, jest stroną, na którą Operator i tak by nie
	// czekał przy otwartym oknie.
	granicaPobraniaStrony = 45 * time.Second
	// granicaTresciStrony chroni rdzeń przed odpowiedzią bez końca: strumień
	// nadawany w nieskończoność wypełniłby pamięć procesu.
	granicaTresciStrony = 32 << 20
	// granicaObrazowStrony ogranicza liczbę obrazów wciąganych do magazynu.
	// Strona z galerią potrafi ich mieć setki, a kolejka wczytywania ma nieść
	// materiał do redakcji, nie kopię cudzej galerii.
	granicaObrazowStrony = 20
)

// WczytajZAdresu obsługuje `studio.ingest.url`.
func (a *adapterStudia) WczytajZAdresu(ctx context.Context,
	z shared.StudioIngestUrlRequest) (shared.StudioIngestUrlResponse, error) {

	if strings.TrimSpace(z.WindowId) == "" {
		return shared.StudioIngestUrlResponse{}, bladWskazaniaStudio("pobranie strony bez okna")
	}
	adres, err := url.Parse(strings.TrimSpace(z.Url))
	if err != nil || adres.Host == "" {
		return shared.StudioIngestUrlResponse{}, bladWskazaniaStudio(
			"adres " + z.Url + " nie jest adresem strony")
	}
	if adres.Scheme != "http" && adres.Scheme != "https" {
		return shared.StudioIngestUrlResponse{}, bladWskazaniaStudio(
			"adres " + z.Url + " ma schemat " + adres.Scheme +
				"; kolejka wczytywania pobiera wyłącznie http i https")
	}

	bajty, typTresci, err := pobierzStroneStudia(ctx, adres.String())
	if err != nil {
		return shared.StudioIngestUrlResponse{}, err
	}

	tresc := string(bajty)
	if strings.Contains(strings.ToLower(typTresci), "html") && wlaczone(z.Readability) {
		tresc = tekstZeStronyStudia(tresc)
	}

	if z.IncludeImages != nil && *z.IncludeImages {
		wciagniete, err := a.wciagnijObrazyStronyStudia(ctx, bajty, adres, z.WindowId)
		if err != nil {
			return shared.StudioIngestUrlResponse{}, err
		}
		if wciagniete != "" {
			tresc += "\n\n" + wciagniete
		}
	}

	zrodlo := adres.String()
	pozycja, err := a.repozytorium.ZapiszPozycjeWczytywania(ctx, dane.PozycjaWczytywania{
		Kod:             nowyIdentyfikator(przedrostekPozycjiWczytywania),
		Okno:            z.WindowId,
		SciezkaZrodlowa: &zrodlo,
		// Pozycja wchodzi od razu jako gotowa: treść jest już wydobyta, więc
		// stan „oczekuje" kazałby Operatorowi puścić rozpoznanie pisma na
		// tekście, który tekstem już jest.
		Stan:  "gotowa",
		Tekst: &tresc,
	})
	if err != nil {
		return shared.StudioIngestUrlResponse{}, bladStudio(err)
	}
	return shared.StudioIngestUrlResponse{Item: złóżPozycjeWczytywania(pozycja)}, nil
}

// SkanujUrzadzenie obsługuje `studio.ingest.device.scan`.
//
// Pobrane strony wchodzą do kolejki jako pozycje w stanie „oczekuje" — tak samo,
// jak materiał dołożony ścieżką. Skan jest obrazem, więc tekstu jeszcze nie ma;
// wpisanie stanu „gotowa" kazałoby Operatorowi przyjąć pustą treść jako wynik.
func (a *adapterStudia) SkanujUrzadzenie(ctx context.Context,
	z shared.StudioIngestDeviceScanRequest) (shared.StudioIngestDeviceScanResponse, error) {

	if strings.TrimSpace(z.WindowId) == "" {
		return shared.StudioIngestDeviceScanResponse{}, bladWskazaniaStudio("skanowanie bez okna")
	}

	zamowienie := zamowienieSkanu{
		Urzadzenie: strings.TrimSpace(wartoscTekstu(z.DeviceId)),
		TrybBarwny: strings.TrimSpace(wartoscTekstu(z.ColorMode)),
	}
	if z.ResolutionDpi != nil {
		zamowienie.Rozdzielczosc = *z.ResolutionDpi
	}
	if z.Pages != nil {
		zamowienie.Stron = *z.Pages
	}

	sciezki, err := a.skanujUrzadzenie(ctx, zamowienie)
	if err != nil {
		return shared.StudioIngestDeviceScanResponse{}, err
	}
	if len(sciezki) == 0 {
		// Warstwa oddała powodzenie bez ani jednego pliku. To nie jest pusta
		// kolejka do przekazania dalej, a usterka warstwy — i jako usterka ma
		// zostać nazwana, zamiast wyglądać na „skaner nic nie podał".
		return shared.StudioIngestDeviceScanResponse{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeInternalError,
			"moduł Studio: warstwa urządzeń zgłosiła udane skanowanie, ale nie wskazała "+
				"ani jednego pliku obrazu; naprawa: to usterka rdzenia, nie nastawa "+
				"Operatora — zgłosić ją wraz z dziennikiem startu"))
	}

	pozycje := make([]shared.StudioIngestItem, 0, len(sciezki))
	for _, sciezka := range sciezki {
		wskazanie := sciezka
		zapisana, err := a.repozytorium.ZapiszPozycjeWczytywania(ctx, dane.PozycjaWczytywania{
			Kod:             nowyIdentyfikator(przedrostekPozycjiWczytywania),
			Okno:            z.WindowId,
			SciezkaZrodlowa: &wskazanie,
			Stan:            "oczekuje",
		})
		if err != nil {
			return shared.StudioIngestDeviceScanResponse{}, bladStudio(err)
		}
		pozycje = append(pozycje, złóżPozycjeWczytywania(zapisana))
	}
	return shared.StudioIngestDeviceScanResponse{Items: pozycje}, nil
}

// pobierzStroneStudia wykonuje żądanie i oddaje treść wraz z jej typem.
func pobierzStroneStudia(ctx context.Context, adres string) ([]byte, string, error) {
	kontekst, zamknij := context.WithTimeout(ctx, granicaPobraniaStrony)
	defer zamknij()

	zadanie, err := http.NewRequestWithContext(kontekst, http.MethodGet, adres, nil)
	if err != nil {
		return nil, "", bladWskazaniaStudio("nie można złożyć żądania do " + adres + ": " + err.Error())
	}
	odpowiedz, err := http.DefaultClient.Do(zadanie)
	if err != nil {
		return nil, "", protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
			"moduł Studio: strona "+adres+" nie odpowiedziała: "+err.Error()))
	}
	defer odpowiedz.Body.Close()

	if odpowiedz.StatusCode >= 400 {
		return nil, "", protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Studio: strona "+adres+" odpowiedziała kodem "+
				strconv.Itoa(odpowiedz.StatusCode)))
	}
	bajty, err := io.ReadAll(io.LimitReader(odpowiedz.Body, granicaTresciStrony))
	if err != nil {
		return nil, "", protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
			"moduł Studio: przerwany odczyt strony "+adres+": "+err.Error()))
	}
	return bajty, odpowiedz.Header.Get("Content-Type"), nil
}

// tekstZeStronyStudia sprowadza stronę do postaci czytelnej: zdejmuje skrypty,
// style i znaczniki, zostawia tekst i akapity.
//
// Czyszczenie jest własne i proste, bez biblioteki czytelności: te biblioteki
// rozstrzygają, KTÓRA część strony jest treścią główną, a rozstrzygnięcie
// błędne wycina Operatorowi połowę artykułu bez ostrzeżenia. Tutaj nic nie
// znika poza nawigacją i reklamą, które i tak nie mają tekstu własnego —
// Operator dostaje więcej, niż prosił, a nie mniej.
func tekstZeStronyStudia(strona string) string {
	bez := usunElementStudia(strona, "script")
	bez = usunElementStudia(bez, "style")
	bez = usunElementStudia(bez, "nav")
	bez = usunElementStudia(bez, "noscript")
	bez = usunElementStudia(bez, "svg")

	var tekst strings.Builder
	wZnaczniku := false
	for _, znak := range bez {
		switch {
		case znak == '<':
			wZnaczniku = true
			// Znacznik jest granicą wiersza: bez tego cały akapit skleiłby się
			// z następnym w jedno długie zdanie.
			tekst.WriteRune('\n')
		case znak == '>':
			wZnaczniku = false
		case !wZnaczniku:
			tekst.WriteRune(znak)
		}
	}

	zamiennik := strings.NewReplacer("&nbsp;", " ", "&amp;", "&", "&lt;", "<",
		"&gt;", ">", "&quot;", `"`, "&#39;", "'", "&apos;", "'")
	wiersze := []string{}
	for _, wiersz := range strings.Split(zamiennik.Replace(tekst.String()), "\n") {
		if oczyszczony := strings.TrimSpace(wiersz); oczyszczony != "" {
			wiersze = append(wiersze, oczyszczony)
		}
	}
	return strings.Join(wiersze, "\n")
}

// usunElementStudia wycina element wraz z jego zawartością.
func usunElementStudia(strona, nazwa string) string {
	male := strings.ToLower(strona)
	otwarcie := "<" + nazwa
	zamkniecie := "</" + nazwa + ">"
	for {
		od := strings.Index(male, otwarcie)
		if od < 0 {
			return strona
		}
		do_ := strings.Index(male[od:], zamkniecie)
		if do_ < 0 {
			return strona[:od]
		}
		strona = strona[:od] + strona[od+do_+len(zamkniecie):]
		male = strings.ToLower(strona)
	}
}

// wciagnijObrazyStronyStudia pobiera obrazy strony do magazynu i oddaje ich
// wykaz zapisem osadzenia — tym samym, którym posługuje się `studio.asset.embed`.
// Dzięki temu obraz wciągnięty ze strony wstawia się do dokumentu tak samo jak
// grafika z Design, a nie drugim sposobem zapisu.
func (a *adapterStudia) wciagnijObrazyStronyStudia(ctx context.Context, strona []byte,
	adres *url.URL, okno string) (string, error) {

	odnosniki := odnosnikiObrazowStudia(string(strona))
	wykaz := []string{}
	for _, odnosnik := range odnosniki {
		if len(wykaz) >= granicaObrazowStrony {
			break
		}
		pelny, err := adres.Parse(odnosnik)
		if err != nil || (pelny.Scheme != "http" && pelny.Scheme != "https") {
			continue
		}
		bajty, typTresci, err := pobierzStroneStudia(ctx, pelny.String())
		if err != nil || len(bajty) == 0 {
			// Obraz niepobrany nie unieważnia strony: Operator prosił o tekst,
			// a obrazy są dodatkiem. Odmowa całości z powodu jednego martwego
			// odnośnika byłaby odmową wczytania artykułu.
			continue
		}
		zasob, err := a.odlozTrescStudia(ctx, bajty, nazwaObrazuStudia(pelny),
			formatObrazuStudia(typTresci, pelny), okno)
		if err != nil {
			return "", err
		}
		wykaz = append(wykaz, "![obraz "+strconv.Itoa(len(wykaz)+1)+
			"](danaco://zasob/"+zasob.Id+")")
	}
	if len(wykaz) == 0 {
		return "", nil
	}
	return "Obrazy strony:\n\n" + strings.Join(wykaz, "\n"), nil
}

// odnosnikiObrazowStudia wyjmuje wartości `src` ze znaczników obrazu.
func odnosnikiObrazowStudia(strona string) []string {
	odnosniki := []string{}
	male := strings.ToLower(strona)
	pozycja := 0
	for {
		od := strings.Index(male[pozycja:], "<img")
		if od < 0 {
			return odnosniki
		}
		od += pozycja
		koniec := strings.Index(male[od:], ">")
		if koniec < 0 {
			return odnosniki
		}
		znacznik := strona[od : od+koniec]
		if wartosc := wartoscCechyStudia(znacznik, "src"); wartosc != "" {
			odnosniki = append(odnosniki, wartosc)
		}
		pozycja = od + koniec
	}
}

// wartoscCechyStudia wyjmuje wartość cechy znacznika w cudzysłowie prostym albo
// podwójnym.
func wartoscCechyStudia(znacznik, cecha string) string {
	male := strings.ToLower(znacznik)
	od := strings.Index(male, cecha+"=")
	if od < 0 {
		return ""
	}
	reszta := znacznik[od+len(cecha)+1:]
	if reszta == "" {
		return ""
	}
	granica := reszta[0]
	if granica != '"' && granica != '\'' {
		if koniec := strings.IndexAny(reszta, " \t\r\n"); koniec >= 0 {
			return reszta[:koniec]
		}
		return reszta
	}
	koniec := strings.IndexByte(reszta[1:], granica)
	if koniec < 0 {
		return ""
	}
	return reszta[1 : 1+koniec]
}

// nazwaObrazuStudia nadaje obrazowi nazwę czytelną w Assets Panelu.
func nazwaObrazuStudia(adres *url.URL) string {
	czesci := strings.Split(strings.Trim(adres.Path, "/"), "/")
	nazwa := czesci[len(czesci)-1]
	if strings.TrimSpace(nazwa) == "" {
		return "obraz-" + adres.Host
	}
	return nazwa
}

// formatObrazuStudia rozstrzyga format po typie treści, a gdy serwer go nie
// podał — po rozszerzeniu w adresie.
func formatObrazuStudia(typTresci string, adres *url.URL) string {
	typ := strings.ToLower(strings.TrimSpace(strings.Split(typTresci, ";")[0]))
	if czesci := strings.SplitN(typ, "/", 2); len(czesci) == 2 && czesci[0] == "image" {
		if czesci[1] == "jpeg" {
			return "jpg"
		}
		if czesci[1] == "svg+xml" {
			return "svg"
		}
		return czesci[1]
	}
	if kropka := strings.LastIndex(adres.Path, "."); kropka >= 0 {
		return strings.ToLower(strings.TrimPrefix(adres.Path[kropka:], "."))
	}
	return "bin"
}
