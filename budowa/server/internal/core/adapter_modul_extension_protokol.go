// Plik obsługuje rodzinę komend extension.*: klienta protokołu Model Context
// Protocol wkompilowanego w rdzeń, nad transportem stdio, SSE albo HTTP,
// z każdą ramką rozmowy zapisywaną do dziennika.
package core

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"danacoconsole/shared"
)

// Granice rozmowy z serwerem MCP. Bez nich pojedyncze wywołanie mogłoby wisieć
// tak długo, jak długo milczy cudzy proces.
const (
	czasPowitaniaMcp = 15 * time.Second
	czasWywolaniaMcp = 60 * time.Second
	// granicaOdpowiedziMcp przycina odczyt odpowiedzi. Serwer, który odeśle
	// gigabajt, nie ma prawa wywrócić rdzenia.
	granicaOdpowiedziMcp = 8 << 20
)

// wersjaProtokoluMcp jest wersją deklarowaną w powitaniu. Wartość pochodzi ze
// specyfikacji protokołu, nie z kontraktu platformy — to jest cudza umowa.
const wersjaProtokoluMcp = "2024-11-05"

// polaczenieMcp opisuje, czym rdzeń woła serwer jednej pozycji katalogu
// rozszerzeń: transportem, adresem sieciowym albo poleceniem procesu
// lokalnego.
type polaczenieMcp struct {
	Transport shared.McpTransport
	Adres     string
	Polecenie string
}

// zadanieMcp to jedna ramka żądania JSON-RPC 2.0 wysyłana do serwera
// rozszerzenia wraz z numerem korelującym ją z odpowiedzią.
type zadanieMcp struct {
	JsonRpc string          `json:"jsonrpc"`
	Id      int             `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// odpowiedzMcp to jedna ramka odpowiedzi JSON-RPC 2.0 odebrana od serwera
// rozszerzenia, niosąca wynik wywołania albo odmowę.
type odpowiedzMcp struct {
	JsonRpc string          `json:"jsonrpc"`
	Id      int             `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *bladMcp        `json:"error,omitempty"`
}

// bladMcp to odmowa serwera wyrażona w protokole JSON-RPC 2.0: kod, komunikat
// i opcjonalny ładunek danych dodatkowych.
type bladMcp struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

// wpisNarzedziaMcp to jeden wpis wykazu tools/list, resources/list albo
// prompts/list. Trzy wykazy mają w protokole trzy nieco różne kształty; pola
// niżej są ich sumą, a nieobecne zostają puste.
type wpisNarzedziaMcp struct {
	Name        string          `json:"name"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	InputSchema json.RawMessage `json:"inputSchema"`
	Uri         string          `json:"uri"`
}

// wykazNarzedziMcp to odpowiedź jednego z trzech wykazów protokołu: narzędzi,
// zasobów albo podpowiedzi serwera rozszerzenia.
type wykazNarzedziMcp struct {
	Tools     []wpisNarzedziaMcp `json:"tools"`
	Resources []wpisNarzedziaMcp `json:"resources"`
	Prompts   []wpisNarzedziaMcp `json:"prompts"`
}

// rozmowaMcp jest jedną sesją z serwerem: powitanie i tyle wywołań, ile trzeba.
// Sesja żyje w granicach jednej komendy, bo serwer HTTP jest bezstanowy między
// żądaniami, a proces stdio wiszący między komendami byłby zasobem, którego
// nikt nie zamyka.
type rozmowaMcp struct {
	polaczenie polaczenieMcp
	klient     *http.Client
	// proces i strumienie żyją wyłącznie przy transporcie stdio.
	proces  *exec.Cmd
	wejscie io.WriteCloser
	wyjscie *bufio.Reader
	numer   int
	// ramki zbiera przebieg rozmowy do dziennika, odkładany w bazie po jej
	// zakończeniu.
	ramki []ramkaRozmowyMcp
}

// ramkaRozmowyMcp to jedna ramka zapisana do dziennika protokołu wraz
// z kierunkiem, metodą, korelacją i chwilą zajścia.
type ramkaRozmowyMcp struct {
	Kierunek  shared.ProtocolFrameDirection
	Metoda    string
	Korelacja string
	Tresc     string
	KodBledu  string
	Zaszlo    int64
}

// otworzRozmoweMcp podnosi rozmowę wedle transportu pozycji: uruchamia proces
// lokalny przy stdio albo zakłada klienta HTTP przy pozostałych transportach.
func otworzRozmoweMcp(ctx context.Context, polaczenie polaczenieMcp) (*rozmowaMcp, error) {
	rozmowa := &rozmowaMcp{polaczenie: polaczenie}

	switch polaczenie.Transport {
	case shared.McpTransportStdio:
		polecenie := strings.Fields(polaczenie.Polecenie)
		if len(polecenie) == 0 {
			return nil, fmt.Errorf("transport stdio wymaga polecenia procesu serwera")
		}
		// #nosec G204 — polecenie pochodzi od operatora, który świadomie podłączył
		// ten serwer.
		proces := exec.CommandContext(ctx, polecenie[0], polecenie[1:]...)
		wejscie, err := proces.StdinPipe()
		if err != nil {
			return nil, fmt.Errorf("nie można otworzyć wejścia procesu serwera: %w", err)
		}
		wyjscie, err := proces.StdoutPipe()
		if err != nil {
			return nil, fmt.Errorf("nie można otworzyć wyjścia procesu serwera: %w", err)
		}
		if err := proces.Start(); err != nil {
			return nil, fmt.Errorf("nie można uruchomić procesu serwera %q: %w", polecenie[0], err)
		}
		rozmowa.proces = proces
		rozmowa.wejscie = wejscie
		rozmowa.wyjscie = bufio.NewReaderSize(wyjscie, 64<<10)

	case shared.McpTransportSse, shared.McpTransportHttp, shared.McpTransportStreamableHttp:
		if strings.TrimSpace(polaczenie.Adres) == "" {
			return nil, fmt.Errorf("transport %s wymaga adresu serwera", polaczenie.Transport)
		}
		rozmowa.klient = &http.Client{Timeout: czasWywolaniaMcp}

	default:
		return nil, fmt.Errorf("nieznany transport %q", polaczenie.Transport)
	}
	return rozmowa, nil
}

// Zamknij kończy rozmowę. Proces serwera zamyka się przez zamknięcie jego
// wejścia; gdy nie odejdzie sam, kontekst rozmowy go ubija.
func (r *rozmowaMcp) Zamknij() {
	if r == nil {
		return
	}
	if r.wejscie != nil {
		_ = r.wejscie.Close()
	}
	if r.proces != nil {
		_ = r.proces.Wait()
	}
}

// Powitaj wykonuje metodę initialize i oddaje wersję protokołu podaną przez
// serwer, a przy odpowiedzi nieczytelnej wersję zadeklarowaną przez rdzeń.
func (r *rozmowaMcp) Powitaj(ctx context.Context) (string, error) {
	parametry, err := json.Marshal(map[string]any{
		"protocolVersion": wersjaProtokoluMcp,
		"capabilities":    map[string]any{},
		"clientInfo": map[string]any{
			"name":    "danaco-console",
			"version": shared.ProtocolVersion,
		},
	})
	if err != nil {
		return "", err
	}
	wynik, err := r.Wywolaj(ctx, "initialize", parametry, czasPowitaniaMcp)
	if err != nil {
		return "", err
	}
	var powitanie struct {
		ProtocolVersion string `json:"protocolVersion"`
	}
	if err := json.Unmarshal(wynik, &powitanie); err != nil {
		// Serwer odpowiedział, ale nie tym kształtem — wersja protokołu jest podawana
		// wtedy własna.
		return wersjaProtokoluMcp, nil
	}
	if powitanie.ProtocolVersion == "" {
		return wersjaProtokoluMcp, nil
	}
	return powitanie.ProtocolVersion, nil
}

// Wywolaj wysyła jedną metodę i czeka na jej odpowiedź, odkładając obie ramki
// do dziennika rozmowy niezależnie od wyniku wywołania.
func (r *rozmowaMcp) Wywolaj(ctx context.Context, metoda string, parametry json.RawMessage,
	granica time.Duration) (json.RawMessage, error) {

	r.numer++
	korelacja := strconv.Itoa(r.numer)
	zadanie := zadanieMcp{JsonRpc: "2.0", Id: r.numer, Method: metoda, Params: parametry}
	bajty, err := json.Marshal(zadanie)
	if err != nil {
		return nil, err
	}
	r.ramki = append(r.ramki, ramkaRozmowyMcp{
		Kierunek: shared.ProtocolFrameDirectionOutgoing, Metoda: metoda,
		Korelacja: korelacja, Tresc: string(bajty), Zaszlo: time.Now().UTC().UnixMilli(),
	})

	wykonanie, zakoncz := context.WithTimeout(ctx, granica)
	defer zakoncz()

	var surowa []byte
	if r.proces != nil {
		surowa, err = r.wymienStdio(bajty)
	} else {
		surowa, err = r.wymienHttp(wykonanie, bajty)
	}
	if err != nil {
		r.ramki = append(r.ramki, ramkaRozmowyMcp{
			Kierunek: shared.ProtocolFrameDirectionIncoming, Korelacja: korelacja,
			KodBledu: "transport", Tresc: err.Error(), Zaszlo: time.Now().UTC().UnixMilli(),
		})
		return nil, err
	}

	var odpowiedz odpowiedzMcp
	if err := json.Unmarshal(surowa, &odpowiedz); err != nil {
		r.ramki = append(r.ramki, ramkaRozmowyMcp{
			Kierunek: shared.ProtocolFrameDirectionIncoming, Korelacja: korelacja,
			KodBledu: "parse", Tresc: string(surowa), Zaszlo: time.Now().UTC().UnixMilli(),
		})
		return nil, fmt.Errorf("odpowiedź serwera nie jest ramką JSON-RPC: %w", err)
	}

	ramka := ramkaRozmowyMcp{
		Kierunek: shared.ProtocolFrameDirectionIncoming, Korelacja: korelacja,
		Tresc: string(surowa), Zaszlo: time.Now().UTC().UnixMilli(),
	}
	if odpowiedz.Error != nil {
		ramka.KodBledu = strconv.Itoa(odpowiedz.Error.Code)
		r.ramki = append(r.ramki, ramka)
		return nil, fmt.Errorf("serwer odmówił (%d): %s", odpowiedz.Error.Code, odpowiedz.Error.Message)
	}
	r.ramki = append(r.ramki, ramka)
	return odpowiedz.Result, nil
}

// wymienStdio wysyła ramkę do procesu i czyta jedną ramkę odpowiedzi.
// Kadrowanie jest liniowe (JSON Lines) — tak nadaje większość serwerów stdio,
// a ramka nagłówkowa `Content-Length` jest rozpoznawana i pomijana niżej.
func (r *rozmowaMcp) wymienStdio(bajty []byte) ([]byte, error) {
	if _, err := r.wejscie.Write(append(bajty, '\n')); err != nil {
		return nil, fmt.Errorf("nie można wysłać ramki do procesu serwera: %w", err)
	}
	for {
		linia, err := r.wyjscie.ReadBytes('\n')
		if err != nil {
			return nil, fmt.Errorf("proces serwera nie odesłał ramki: %w", err)
		}
		przycieta := bytes.TrimSpace(linia)
		if len(przycieta) == 0 {
			continue
		}
		// Nagłówek kadrowania długością nie jest ramką — jest pomijany, a odczyt trwa
		// dalej.
		if bytes.HasPrefix(przycieta, []byte("Content-Length:")) {
			continue
		}
		if przycieta[0] != '{' {
			continue
		}
		return przycieta, nil
	}
}

// wymienHttp wysyła ramkę żądaniem POST i czyta odpowiedź. Transport SSE
// odpowiada strumieniem zdarzeń — ramka jedzie wtedy w polu `data:`, więc
// odczyt rozpoznaje obie postacie zamiast zakładać jedną.
func (r *rozmowaMcp) wymienHttp(ctx context.Context, bajty []byte) ([]byte, error) {
	zadanie, err := http.NewRequestWithContext(ctx, http.MethodPost, r.polaczenie.Adres,
		bytes.NewReader(bajty))
	if err != nil {
		return nil, err
	}
	zadanie.Header.Set("Content-Type", "application/json")
	zadanie.Header.Set("Accept", "application/json, text/event-stream")

	odpowiedz, err := r.klient.Do(zadanie)
	if err != nil {
		return nil, fmt.Errorf("serwer pod %s nie odpowiada: %w", r.polaczenie.Adres, err)
	}
	defer odpowiedz.Body.Close()

	tresc, err := io.ReadAll(io.LimitReader(odpowiedz.Body, granicaOdpowiedziMcp))
	if err != nil {
		return nil, fmt.Errorf("nie można odczytać odpowiedzi serwera: %w", err)
	}
	if odpowiedz.StatusCode >= 400 {
		return nil, fmt.Errorf("serwer odpowiedział kodem %d: %s",
			odpowiedz.StatusCode, poczatekTekstuApp(string(tresc)))
	}
	if strings.Contains(odpowiedz.Header.Get("Content-Type"), "text/event-stream") {
		return ramkaZeStrumieniaMcp(tresc)
	}
	return bytes.TrimSpace(tresc), nil
}

// ramkaZeStrumieniaMcp wyciąga pierwszą ramkę JSON-RPC z odpowiedzi SSE,
// pomijając wiersze niebędące polem danych zdarzenia.
func ramkaZeStrumieniaMcp(tresc []byte) ([]byte, error) {
	for _, linia := range strings.Split(string(tresc), "\n") {
		przycieta := strings.TrimSpace(linia)
		if !strings.HasPrefix(przycieta, "data:") {
			continue
		}
		ladunek := strings.TrimSpace(strings.TrimPrefix(przycieta, "data:"))
		if ladunek == "" || ladunek[0] != '{' {
			continue
		}
		return []byte(ladunek), nil
	}
	return nil, fmt.Errorf("strumień zdarzeń nie niósł ani jednej ramki JSON-RPC")
}

// OdkryjWykazy pobiera trzy wykazy protokołu i składa je w jeden zbiór wpisów.
// Wykaz, którego serwer nie udostępnia, jest pustką, nie awarią: `resources`
// i `prompts` są w protokole opcjonalne.
func (r *rozmowaMcp) OdkryjWykazy(ctx context.Context) ([]wpisWykazuMcp, error) {
	wpisy := []wpisWykazuMcp{}
	for _, wykaz := range []struct {
		metoda string
		rodzaj shared.ExtensionToolKind
	}{
		{"tools/list", shared.ExtensionToolKindTool},
		{"resources/list", shared.ExtensionToolKindResource},
		{"prompts/list", shared.ExtensionToolKindPrompt},
	} {
		wynik, err := r.Wywolaj(ctx, wykaz.metoda, json.RawMessage(`{}`), czasPowitaniaMcp)
		if err != nil {
			// Brak wykazu nie przewraca odkrycia — pierwszy wykaz jest jednak
			// obowiązkowy.
			if wykaz.rodzaj == shared.ExtensionToolKindTool {
				return nil, err
			}
			continue
		}
		var wykazOdpowiedzi wykazNarzedziMcp
		if err := json.Unmarshal(wynik, &wykazOdpowiedzi); err != nil {
			continue
		}
		zrodlo := wykazOdpowiedzi.Tools
		switch wykaz.rodzaj {
		case shared.ExtensionToolKindResource:
			zrodlo = wykazOdpowiedzi.Resources
		case shared.ExtensionToolKindPrompt:
			zrodlo = wykazOdpowiedzi.Prompts
		}
		for _, wpis := range zrodlo {
			nazwa := wpis.Name
			if nazwa == "" {
				nazwa = wpis.Uri
			}
			if nazwa == "" {
				continue
			}
			wpisy = append(wpisy, wpisWykazuMcp{
				Rodzaj: wykaz.rodzaj, Nazwa: nazwa, Opis: wpis.Description,
				Schemat: string(wpis.InputSchema), Adres: wpis.Uri,
			})
		}
	}
	return wpisy, nil
}

// wpisWykazuMcp to jeden wpis odkryty u serwera rozszerzenia, w kształcie
// własnym rdzenia, wspólnym dla narzędzi, zasobów i podpowiedzi.
type wpisWykazuMcp struct {
	Rodzaj  shared.ExtensionToolKind
	Nazwa   string
	Opis    string
	Schemat string
	Adres   string
}

// poczatekTekstuApp przycina cudzą odpowiedź do wielkości czytelnej
// w odmowie, żeby długi ładunek serwera nie zalał komunikatu błędu.
func poczatekTekstuApp(tekst string) string {
	const granica = 400
	if len(tekst) <= granica {
		return tekst
	}
	return tekst[:granica] + "…"
}
