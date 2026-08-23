// Odpowiedzialność pliku: silnik przeglądarki modułu Browser — uruchomienie
// Chromium bez okna i rozmowa z nim protokołem Chrome DevTools (CDP).
//
// ── Po co silnik, skoro rdzeń umie pobrać stronę HTTP-em ─────────────────────
// `przegladarka_pobieranie.go` pobiera dokument biblioteką standardową i to
// wystarcza do migawki tekstu. Nie wystarcza do niczego, co wymaga wykonania
// strony: zrzut ekranu, drzewo DOM po zbudowaniu przez skrypty, komunikaty
// konsoli, rejestr żądań sieciowych, emulacja urządzenia i przewinięcie są
// własnościami strony URUCHOMIONEJ, nie jej źródła. Opracowanie modułu wskazuje
// tu wprost Chrome DevTools Protocol (rozdz. 2.11, 2.14) i tą drogą to idzie.
//
// ── Dlaczego program, a nie biblioteka wkompilowana ───────────────────────────
// Zasada produktu mówi: żadna funkcja nie zależy od programu, którego instalka
// nie niesie. Cała aplikacja z arsenałem stoi na serwerze, u Operatora jest samo
// okno — więc Chromium jest programem SERWEROWYM, tak samo jak ffmpeg czy
// Tesseract, i jak one stoi w sondzie zależności (`zaleznosci_zewnetrzne.go`).
// Silnika przeglądarki nie da się wkompilować w binarium Go; wyjątek na
// biblioteki wkompilowane obejmuje PDF i kryptografię, a nie renderowanie stron.
//
// ── Dlaczego proces startuje tu, a nie przez `zewnetrzne.Wolaj` ──────────────
// `Wolaj` prowadzi uruchomienie DO KOŃCA i oddaje bajty po zakończeniu programu.
// Przeglądarka ma żyć, dopóki trwa rozmowa: startuje, przyjmuje polecenia
// protokołem i dopiero potem gaśnie. Sekwencja jest jednak ta sama — port
// `session.Uruchamiacz`, brama izolacji `session.SprawdzPolecenie`, objęcie
// drzewa procesów `session.PrzejmijDrzewo` — bo Chromium rozgałęzia procesy
// renderowania i sieci, a przerwana sesja bez objęcia drzewa zostawiłaby je
// na maszynie Operatora.
//
// ── Granica czasu jest zawsze ─────────────────────────────────────────────────
// Strona, która nie kończy wczytywania, jest zjawiskiem codziennym. Każde
// otwarcie ma granicę; po jej przekroczeniu sesja oddaje to, co zdążyła zebrać,
// albo odmawia — nigdy nie czeka bez końca.
package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/coder/websocket"

	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/session"
	"danacoconsole/server/internal/zewnetrzne"
	"danacoconsole/shared"
)

// narzedzieChromium opisuje silnik przeglądarki. Program bywa pod dwiema
// nazwami zależnie od dystrybucji i sposobu instalacji, więc deklaracja wybiera
// tę, która na tej maszynie naprawdę stoi — podpowiedź instalacyjna kierująca do
// pakietu, którego nie ma, byłaby podpowiedzią donikąd.
func narzedzieChromium() zewnetrzne.Narzedzie {
	for _, program := range []string{"chromium-browser", "chromium", "google-chrome"} {
		kandydat := zewnetrzne.Narzedzie{Nazwa: "Chromium", Program: program, Pakiet: "chromium-browser"}
		if zewnetrzne.Stoi(kandydat) {
			return kandydat
		}
	}
	return zewnetrzne.Narzedzie{Nazwa: "Chromium", Program: "chromium-browser", Pakiet: "chromium-browser"}
}

const (
	// granicaOtwarciaStrony obejmuje start przeglądarki, przejście pod adres
	// i zebranie wyniku. Hojna, bo obejmuje uruchomienie procesu; skończona, bo
	// strona wisząca w nieskończoność jest zjawiskiem codziennym.
	granicaOtwarciaStrony = 90 * time.Second
	// granicaWywolaniaCdp jest granicą jednej odpowiedzi protokołu. Zrzut całej
	// strony bywa wolny, ale nie trwa minutami.
	granicaWywolaniaCdp = 45 * time.Second
	// granicaStartuSilnika jest czasem, w którym przeglądarka ma zgłosić
	// gotowość swojego punktu diagnostycznego.
	granicaStartuSilnika = 30 * time.Second
	// granicaOdczytuCdp jest górnym rozmiarem jednej odpowiedzi protokołu.
	// Zrzut całej długiej strony w PNG idzie w megabajty, a odpowiedź obcięta
	// byłaby obrazem uszkodzonym podanym jako zrzut.
	granicaOdczytuCdp = 256 << 20
)

// nastawyStrony opisują, w jakich warunkach strona ma zostać otwarta.
// Wartości zerowe znaczą „jak domyślnie", nie „zero pikseli".
type nastawyStrony struct {
	Url            string
	Szerokosc      int
	Wysokosc       int
	SkalaPikseli   float64
	Mobilne        bool
	Orientacja     string
	AgentUzytkow   string
	CzekajNaWarune string
	LimitCzasu     time.Duration
}

// wynikOtwarcia niesie to, co dała jedna wizyta na stronie.
type wynikOtwarcia struct {
	Url       string
	Tytul     string
	Tekst     string
	Html      string
	Szerokosc int
	Wysokosc  int
}

// zdarzenieCdp jest jednym powiadomieniem przeglądarki: nazwą metody
// i surowymi parametrami. Rozbiór zostawiony czytelnikowi — rejestr sieciowy
// i konsola czytają z tego strumienia inne pola.
type zdarzenieCdp struct {
	Metoda     string
	Parametry  json.RawMessage
	Odnotowano time.Time
}

// silnikPrzegladarki startuje przeglądarkę i prowadzi z nią rozmowę.
// Nie pamięta niczego między wywołaniami: każda czynność otwiera własną sesję
// i gasi ją po sobie. Sesja trwała między komendami zostawiałaby na maszynie
// proces przeglądarki żyjący dłużej niż powód jego istnienia.
type silnikPrzegladarki struct {
	uruchamiacz  session.Uruchamiacz
	rozstrzygacz *konfig.Rozstrzygacz
	katalog      *KatalogRoboczy
}

// sesjaStrony jest jedną żywą stroną: procesem przeglądarki, połączeniem
// protokołu i zebranymi po drodze zdarzeniami.
type sesjaStrony struct {
	polaczenie *websocket.Conn
	adresPunkt string
	celID      string
	zamknij    func()

	numer  int64
	zamek  sync.Mutex
	czekaj map[int64]chan json.RawMessage

	zamekZdarzen sync.Mutex
	zdarzenia    []zdarzenieCdp

	zycie             context.Context
	zakoncz           context.CancelFunc
	czytanieSkonczone chan struct{}
}

// dostepny mówi, czy silnik ma czym wystartować przeglądarkę. Odmowa nazywająca
// brak jest tu potrzebna wcześnie: komenda, która i tak nie ruszy, ma powiedzieć
// dlaczego, zanim zacznie cokolwiek zapisywać.
func (s *silnikPrzegladarki) dostepny() error {
	if s == nil || s.uruchamiacz == nil {
		return errors.New("rdzeń nie ma uruchamiacza procesów — silnik przeglądarki nie ma czym wystartować; " +
			"naprawa: podpiąć warstwę kanału (injection) przy składaniu rdzenia")
	}
	narzedzie := narzedzieChromium()
	if !zewnetrzne.Stoi(narzedzie) {
		return &zewnetrzne.BrakNarzedzia{Narzedzie: narzedzie}
	}
	return nil
}

// zasiegSilnika składa okno, zasady izolacji i obszar — tak samo jak robią to
// pozostałe rodziny wołające programy (wzór: `adapter_narzedzia_archiwum.go`).
func (s *silnikPrzegladarki) zasiegSilnika() (session.Okno, session.Zasady, session.Obszar) {
	okno := session.Okno{Ustawienia: session.Ustawienia{
		SrodowiskoWykonania: shared.ExecutionEnvCore,
	}}
	zasady := session.Zasady{}
	obszar := session.Obszar{}
	if s.rozstrzygacz != nil {
		zasady = ZasadyIzolacji(s.rozstrzygacz, konfig.Kontekst{})
	}
	if s.katalog != nil {
		obszar = ObszarOkna(s.katalog.Ustal(konfig.Kontekst{}, ""), "")
	}
	return okno, zasady, obszar
}

// otworz startuje przeglądarkę, przechodzi pod wskazany adres i oddaje sesję
// gotową do dalszych pytań. Wołający ZAWSZE zamyka sesję — inaczej proces
// przeglądarki zostaje na maszynie.
func (s *silnikPrzegladarki) otworz(ctx context.Context, nastawy nastawyStrony) (*sesjaStrony, wynikOtwarcia, error) {
	if err := s.dostepny(); err != nil {
		return nil, wynikOtwarcia{}, err
	}
	adres := strings.TrimSpace(nastawy.Url)
	if adres == "" {
		return nil, wynikOtwarcia{}, errors.New("silnik przeglądarki: przejście bez adresu")
	}
	limit := nastawy.LimitCzasu
	if limit <= 0 {
		limit = granicaOtwarciaStrony
	}

	port, err := wolnyPort()
	if err != nil {
		return nil, wynikOtwarcia{}, err
	}
	profil, err := os.MkdirTemp("", "danaco-przegladarka-")
	if err != nil {
		return nil, wynikOtwarcia{}, fmt.Errorf("silnik przeglądarki: nie można założyć katalogu profilu: %w", err)
	}

	szerokosc, wysokosc := wymiaryOkna(nastawy)
	argumenty := []string{
		"--headless=new",
		"--no-sandbox",
		"--disable-gpu",
		"--hide-scrollbars",
		"--disable-dev-shm-usage",
		// Wykaz jest wyłączony świadomie: sesja przeglądania modułu ma być
		// każdorazowo czysta, bez pierwszego uruchomienia, bez przywracania kart
		// i bez okien powitalnych, które przesłoniłyby zrzut strony.
		"--no-first-run",
		"--no-default-browser-check",
		"--disable-background-networking",
		"--remote-allow-origins=*",
		"--user-data-dir=" + profil,
		"--window-size=" + strconv.Itoa(szerokosc) + "," + strconv.Itoa(wysokosc),
		"--remote-debugging-port=" + strconv.Itoa(port),
		"about:blank",
	}
	if agent := strings.TrimSpace(nastawy.AgentUzytkow); agent != "" {
		argumenty = append([]string{"--user-agent=" + agent}, argumenty...)
	}

	okno, zasady, obszar := s.zasiegSilnika()
	sciezka, jest := zewnetrzne.Odnajdz(narzedzieChromium())
	if !jest {
		_ = os.RemoveAll(profil)
		return nil, wynikOtwarcia{}, &zewnetrzne.BrakNarzedzia{Narzedzie: narzedzieChromium()}
	}
	dopuszczone, err := session.SprawdzPolecenie(zasady, obszar, session.Polecenie{
		Program: sciezka, Argumenty: argumenty, Katalog: obszar.KatalogRoboczy,
	})
	if err != nil {
		_ = os.RemoveAll(profil)
		return nil, wynikOtwarcia{}, err
	}
	uchwyt, err := s.uruchamiacz.UruchomProces(okno, dopuszczone)
	if err != nil {
		_ = os.RemoveAll(profil)
		return nil, wynikOtwarcia{}, fmt.Errorf("silnik przeglądarki: nie można uruchomić Chromium: %w", err)
	}
	drzewo, err := session.PrzejmijDrzewo(uchwyt.Pid())
	if err != nil {
		_ = uchwyt.Ubij()
		_ = uchwyt.Czekaj()
		_ = os.RemoveAll(profil)
		return nil, wynikOtwarcia{}, fmt.Errorf("silnik przeglądarki: nie można objąć drzewa procesów Chromium: %w", err)
	}
	// Strumienie przeglądarki są pompowane do kosza, ale pompowane muszą być:
	// Chromium pisze na diagnostykę obficie, a pełny bufor potoku zatrzymałby
	// jego pracę na zapisie i wyglądałoby to jak strona, która się nie wczytuje.
	go pochlonStrumien(uchwyt)

	sprzatanie := func() {
		_ = drzewo.Ubij()
		drzewo.Zwolnij()
		_ = uchwyt.Czekaj()
		_ = os.RemoveAll(profil)
	}

	adresPunktu := "http://127.0.0.1:" + strconv.Itoa(port)
	if err := poczekajNaSilnik(ctx, adresPunktu, granicaStartuSilnika); err != nil {
		sprzatanie()
		return nil, wynikOtwarcia{}, err
	}

	sesja, err := polaczZeStrona(ctx, adresPunktu, sprzatanie)
	if err != nil {
		sprzatanie()
		return nil, wynikOtwarcia{}, err
	}

	wynik, err := sesja.przejdz(ctx, nastawy, limit, szerokosc, wysokosc)
	if err != nil {
		sesja.Zamknij()
		return nil, wynikOtwarcia{}, err
	}
	return sesja, wynik, nil
}

// wymiaryOkna oddaje rozmiar widoku, w którym strona ma zostać wyrenderowana.
func wymiaryOkna(nastawy nastawyStrony) (int, int) {
	szerokosc, wysokosc := nastawy.Szerokosc, nastawy.Wysokosc
	if szerokosc <= 0 {
		szerokosc = 1280
	}
	if wysokosc <= 0 {
		wysokosc = 900
	}
	if strings.EqualFold(nastawy.Orientacja, "landscape") && wysokosc > szerokosc {
		szerokosc, wysokosc = wysokosc, szerokosc
	}
	if strings.EqualFold(nastawy.Orientacja, "portrait") && szerokosc > wysokosc {
		szerokosc, wysokosc = wysokosc, szerokosc
	}
	return szerokosc, wysokosc
}

// wolnyPort rezerwuje port punktu diagnostycznego. Wybór przez system zamiast
// stałej: dwie sesje przeglądania biegnące równolegle nie mają się bić o numer.
func wolnyPort() (int, error) {
	nasluch, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, fmt.Errorf("silnik przeglądarki: nie można zarezerwować portu: %w", err)
	}
	defer nasluch.Close()
	return nasluch.Addr().(*net.TCPAddr).Port, nil
}

// pochlonStrumien opróżnia strumienie procesu przeglądarki.
func pochlonStrumien(uchwyt session.UchwytProcesu) {
	bufor := make([]byte, 4096)
	czytaj := func(strumien interface{ Read([]byte) (int, error) }) {
		if strumien == nil {
			return
		}
		for {
			if _, err := strumien.Read(bufor); err != nil {
				return
			}
		}
	}
	go czytaj(uchwyt.Wyjscie())
	czytaj(uchwyt.Diagnostyka())
}

// poczekajNaSilnik czeka, aż punkt diagnostyczny przeglądarki odpowie.
func poczekajNaSilnik(ctx context.Context, adresPunktu string, limit time.Duration) error {
	koniec := time.Now().Add(limit)
	for {
		zadanie, err := http.NewRequestWithContext(ctx, http.MethodGet, adresPunktu+"/json/version", nil)
		if err != nil {
			return fmt.Errorf("silnik przeglądarki: nie można zapytać o gotowość: %w", err)
		}
		odpowiedz, err := http.DefaultClient.Do(zadanie)
		if err == nil {
			odpowiedz.Body.Close()
			if odpowiedz.StatusCode == http.StatusOK {
				return nil
			}
		}
		if time.Now().After(koniec) {
			return errors.New("silnik przeglądarki: Chromium nie zgłosił gotowości w czasie " + limit.String() +
				"; naprawa: sprawdzić, czy program uruchamia się na tej maszynie bez okna")
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(200 * time.Millisecond):
		}
	}
}

// celStrony jest odpowiedzią punktu diagnostycznego o nowo otwartej karcie.
type celStrony struct {
	Id       string `json:"id"`
	Adres    string `json:"webSocketDebuggerUrl"`
	Rodzaj   string `json:"type"`
	AdresUrl string `json:"url"`
}

// polaczZeStrona zakłada nową kartę w przeglądarce i łączy się z nią protokołem.
func polaczZeStrona(ctx context.Context, adresPunktu string, sprzatanie func()) (*sesjaStrony, error) {
	zadanie, err := http.NewRequestWithContext(ctx, http.MethodPut, adresPunktu+"/json/new?about:blank", nil)
	if err != nil {
		return nil, fmt.Errorf("silnik przeglądarki: nie można poprosić o kartę: %w", err)
	}
	odpowiedz, err := http.DefaultClient.Do(zadanie)
	if err != nil {
		return nil, fmt.Errorf("silnik przeglądarki: przeglądarka nie oddała karty: %w", err)
	}
	defer odpowiedz.Body.Close()

	var cel celStrony
	if err := json.NewDecoder(odpowiedz.Body).Decode(&cel); err != nil {
		return nil, fmt.Errorf("silnik przeglądarki: nieczytelna odpowiedź o karcie: %w", err)
	}
	if cel.Adres == "" {
		return nil, errors.New("silnik przeglądarki: przeglądarka nie podała adresu połączenia z kartą")
	}

	polaczenie, _, err := websocket.Dial(ctx, cel.Adres, nil)
	if err != nil {
		return nil, fmt.Errorf("silnik przeglądarki: nie można połączyć się z kartą: %w", err)
	}
	polaczenie.SetReadLimit(granicaOdczytuCdp)

	zycie, zakoncz := context.WithCancel(context.Background())
	sesja := &sesjaStrony{
		polaczenie: polaczenie, adresPunkt: adresPunktu, celID: cel.Id, zamknij: sprzatanie,
		czekaj: map[int64]chan json.RawMessage{},
		zycie:  zycie, zakoncz: zakoncz,
		czytanieSkonczone: make(chan struct{}),
	}
	go sesja.czytaj()
	return sesja, nil
}

// kopertaCdp jest jedną wiadomością protokołu — odpowiedzią na wywołanie albo
// powiadomieniem o zdarzeniu.
type kopertaCdp struct {
	Id        int64           `json:"id"`
	Metoda    string          `json:"method"`
	Wynik     json.RawMessage `json:"result"`
	Parametry json.RawMessage `json:"params"`
	Blad      *struct {
		Kod    int    `json:"code"`
		Tresc  string `json:"message"`
		Szczeg string `json:"data"`
	} `json:"error"`
}

// czytaj pompuje wiadomości protokołu: odpowiedzi trafiają do czekających
// wywołań, powiadomienia — do wykazu zdarzeń sesji.
func (s *sesjaStrony) czytaj() {
	defer close(s.czytanieSkonczone)
	for {
		_, bajty, err := s.polaczenie.Read(s.zycie)
		if err != nil {
			return
		}
		var koperta kopertaCdp
		if err := json.Unmarshal(bajty, &koperta); err != nil {
			continue
		}
		if koperta.Metoda != "" {
			s.zamekZdarzen.Lock()
			s.zdarzenia = append(s.zdarzenia, zdarzenieCdp{
				Metoda: koperta.Metoda, Parametry: koperta.Parametry, Odnotowano: time.Now(),
			})
			s.zamekZdarzen.Unlock()
			continue
		}
		s.zamek.Lock()
		odbiorca, jest := s.czekaj[koperta.Id]
		delete(s.czekaj, koperta.Id)
		s.zamek.Unlock()
		if !jest {
			continue
		}
		if koperta.Blad != nil {
			odbiorca <- json.RawMessage(`{"__blad":` + strconv.Quote(koperta.Blad.Tresc) + `}`)
			continue
		}
		odbiorca <- koperta.Wynik
	}
}

// wywolaj wysyła jedno polecenie protokołu i czeka na jego odpowiedź.
func (s *sesjaStrony) wywolaj(ctx context.Context, metoda string, parametry map[string]any) (json.RawMessage, error) {
	s.zamek.Lock()
	s.numer++
	numer := s.numer
	odbiorca := make(chan json.RawMessage, 1)
	s.czekaj[numer] = odbiorca
	s.zamek.Unlock()

	if parametry == nil {
		parametry = map[string]any{}
	}
	bajty, err := json.Marshal(map[string]any{"id": numer, "method": metoda, "params": parametry})
	if err != nil {
		return nil, fmt.Errorf("silnik przeglądarki: nie można złożyć polecenia %s: %w", metoda, err)
	}
	if err := s.polaczenie.Write(ctx, websocket.MessageText, bajty); err != nil {
		return nil, fmt.Errorf("silnik przeglądarki: nie można wysłać polecenia %s: %w", metoda, err)
	}

	zegar := time.NewTimer(granicaWywolaniaCdp)
	defer zegar.Stop()
	select {
	case wynik := <-odbiorca:
		var odmowa struct {
			Blad *string `json:"__blad"`
		}
		if err := json.Unmarshal(wynik, &odmowa); err == nil && odmowa.Blad != nil {
			return nil, fmt.Errorf("silnik przeglądarki: %s odmówiło: %s", metoda, *odmowa.Blad)
		}
		return wynik, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-zegar.C:
		return nil, errors.New("silnik przeglądarki: " + metoda + " nie odpowiedziało w czasie " +
			granicaWywolaniaCdp.String())
	}
}

// Zdarzenia oddaje kopię zebranych powiadomień. Kopia, a nie wykaz wprost:
// czytelnik przegląda je w swoim tempie, a pompa dopisuje kolejne.
func (s *sesjaStrony) Zdarzenia() []zdarzenieCdp {
	s.zamekZdarzen.Lock()
	defer s.zamekZdarzen.Unlock()
	kopia := make([]zdarzenieCdp, len(s.zdarzenia))
	copy(kopia, s.zdarzenia)
	return kopia
}

// Zamknij gasi kartę, połączenie i cały proces przeglądarki.
func (s *sesjaStrony) Zamknij() {
	if s == nil {
		return
	}
	if s.polaczenie != nil {
		_ = s.polaczenie.Close(websocket.StatusNormalClosure, "koniec sesji przeglądania")
	}
	s.zakoncz()
	<-s.czytanieSkonczone
	if s.zamknij != nil {
		s.zamknij()
	}
}

// przejdz przechodzi pod adres, czeka na wczytanie i zbiera migawkę treści.
//
// Czekanie jest dwustopniowe: najpierw zdarzenie wczytania dokumentu, potem
// krótkie osiadanie, w którym skrypty zdążą dopisać treść. Bez drugiego kroku
// migawka strony budowanej skryptem byłaby pusta, choć strona jest pełna.
func (s *sesjaStrony) przejdz(ctx context.Context, nastawy nastawyStrony,
	limit time.Duration, szerokosc, wysokosc int) (wynikOtwarcia, error) {

	zycie, zakoncz := context.WithTimeout(ctx, limit)
	defer zakoncz()

	for _, metoda := range []string{"Page.enable", "Runtime.enable", "Log.enable", "Network.enable"} {
		if _, err := s.wywolaj(zycie, metoda, nil); err != nil {
			return wynikOtwarcia{}, err
		}
	}
	skala := nastawy.SkalaPikseli
	if skala <= 0 {
		skala = 1
	}
	if _, err := s.wywolaj(zycie, "Emulation.setDeviceMetricsOverride", map[string]any{
		"width": szerokosc, "height": wysokosc, "deviceScaleFactor": skala, "mobile": nastawy.Mobilne,
	}); err != nil {
		return wynikOtwarcia{}, err
	}
	if agent := strings.TrimSpace(nastawy.AgentUzytkow); agent != "" {
		if _, err := s.wywolaj(zycie, "Emulation.setUserAgentOverride", map[string]any{
			"userAgent": agent,
		}); err != nil {
			return wynikOtwarcia{}, err
		}
	}

	if _, err := s.wywolaj(zycie, "Page.navigate", map[string]any{"url": nastawy.Url}); err != nil {
		return wynikOtwarcia{}, err
	}
	s.poczekajNaWczytanie(zycie, nastawy.CzekajNaWarune)

	wynik := wynikOtwarcia{Szerokosc: szerokosc, Wysokosc: wysokosc}
	stan, err := s.ocenNaStronie(zycie, `({
		url: location.href,
		title: document.title || '',
		text: (document.body && document.body.innerText) || '',
		html: document.documentElement ? document.documentElement.outerHTML : '',
		width: Math.max(document.documentElement.scrollWidth, window.innerWidth),
		height: Math.max(document.documentElement.scrollHeight, window.innerHeight)
	})`)
	if err != nil {
		return wynikOtwarcia{}, err
	}
	var odczyt struct {
		Url       string `json:"url"`
		Tytul     string `json:"title"`
		Tekst     string `json:"text"`
		Html      string `json:"html"`
		Szerokosc int    `json:"width"`
		Wysokosc  int    `json:"height"`
	}
	if err := json.Unmarshal(stan, &odczyt); err != nil {
		return wynikOtwarcia{}, fmt.Errorf("silnik przeglądarki: nieczytelna migawka strony: %w", err)
	}
	wynik.Url, wynik.Tytul, wynik.Tekst, wynik.Html = odczyt.Url, odczyt.Tytul, odczyt.Tekst, odczyt.Html
	if odczyt.Szerokosc > 0 {
		wynik.Szerokosc = odczyt.Szerokosc
	}
	if odczyt.Wysokosc > 0 {
		wynik.Wysokosc = odczyt.Wysokosc
	}
	return wynik, nil
}

// poczekajNaWczytanie czeka na zdarzenie wczytania strony, a potem daje jej
// chwilę na osiadanie. Brak zdarzenia nie jest tu odmową: strona bywa wczytana
// wcześniej, niż zdążyliśmy zacząć nasłuchiwać, a migawka strony częściowej jest
// więcej warta niż odmowa z powodu niedoczekanego powiadomienia.
func (s *sesjaStrony) poczekajNaWczytanie(ctx context.Context, warunek string) {
	oczekiwane := "Page.loadEventFired"
	if strings.EqualFold(warunek, "domContentLoaded") {
		oczekiwane = "Page.domContentEventFired"
	}
	koniec := time.Now().Add(20 * time.Second)
	for time.Now().Before(koniec) {
		for _, zdarzenie := range s.Zdarzenia() {
			if zdarzenie.Metoda == oczekiwane {
				time.Sleep(300 * time.Millisecond)
				return
			}
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(100 * time.Millisecond):
		}
	}
}

// oceń wykonuje wyrażenie w kontekście strony i oddaje jego wartość.
// Nazwa metody mówi wprost, co się dzieje: rdzeń nie „pyta strony", tylko
// URUCHAMIA na niej wyrażenie — i tylko takie, które sam napisał.
func (s *sesjaStrony) ocenNaStronie(ctx context.Context, wyrazenie string) (json.RawMessage, error) {
	wynik, err := s.wywolaj(ctx, "Runtime.evaluate", map[string]any{
		"expression": wyrazenie, "returnByValue": true, "awaitPromise": true,
	})
	if err != nil {
		return nil, err
	}
	var odpowiedz struct {
		Wynik struct {
			Wartosc json.RawMessage `json:"value"`
		} `json:"result"`
		Wyjatek *struct {
			Tekst string `json:"text"`
		} `json:"exceptionDetails"`
	}
	if err := json.Unmarshal(wynik, &odpowiedz); err != nil {
		return nil, fmt.Errorf("silnik przeglądarki: nieczytelna odpowiedź strony: %w", err)
	}
	if odpowiedz.Wyjatek != nil {
		return nil, errors.New("silnik przeglądarki: strona odrzuciła wyrażenie: " + odpowiedz.Wyjatek.Tekst)
	}
	if len(odpowiedz.Wynik.Wartosc) == 0 {
		return json.RawMessage("null"), nil
	}
	return odpowiedz.Wynik.Wartosc, nil
}

// zrzut wykonuje zrzut strony i oddaje jego bajty w zapisie base64 protokołu.
// Tryb rozstrzyga, co znaczy „strona": widok, całość przewijana, prostokąt albo
// pojedynczy element wskazany selektorem.
func (s *sesjaStrony) zrzut(ctx context.Context, tryb, format, selektor string, jakosc int,
	x, y, szerokosc, wysokosc int, wymiarStrony [2]int) (string, int, int, error) {

	parametry := map[string]any{"format": format}
	if format == "jpeg" || format == "webp" {
		if jakosc <= 0 || jakosc > 100 {
			jakosc = 85
		}
		parametry["quality"] = jakosc
	}
	szerokoscZrzutu, wysokoscZrzutu := wymiarStrony[0], wymiarStrony[1]

	switch tryb {
	case "fullPage":
		parametry["captureBeyondViewport"] = true
		parametry["clip"] = map[string]any{
			"x": 0, "y": 0, "width": wymiarStrony[0], "height": wymiarStrony[1], "scale": 1,
		}
	case "region":
		if szerokosc <= 0 || wysokosc <= 0 {
			return "", 0, 0, errors.New("silnik przeglądarki: zrzut obszaru bez szerokości albo wysokości")
		}
		parametry["captureBeyondViewport"] = true
		parametry["clip"] = map[string]any{
			"x": x, "y": y, "width": szerokosc, "height": wysokosc, "scale": 1,
		}
		szerokoscZrzutu, wysokoscZrzutu = szerokosc, wysokosc
	case "element":
		if strings.TrimSpace(selektor) == "" {
			return "", 0, 0, errors.New("silnik przeglądarki: zrzut elementu bez selektora")
		}
		ramka, err := s.ramkaElementu(ctx, selektor)
		if err != nil {
			return "", 0, 0, err
		}
		parametry["captureBeyondViewport"] = true
		parametry["clip"] = map[string]any{
			"x": ramka[0], "y": ramka[1], "width": ramka[2], "height": ramka[3], "scale": 1,
		}
		szerokoscZrzutu, wysokoscZrzutu = int(ramka[2]), int(ramka[3])
	default:
		// Widok — bez przycięcia; wymiar zrzutu jest wymiarem okna, nie strony.
		szerokoscZrzutu, wysokoscZrzutu = wymiarStrony[0], wymiarStrony[1]
	}

	wynik, err := s.wywolaj(ctx, "Page.captureScreenshot", parametry)
	if err != nil {
		return "", 0, 0, err
	}
	var odpowiedz struct {
		Dane string `json:"data"`
	}
	if err := json.Unmarshal(wynik, &odpowiedz); err != nil {
		return "", 0, 0, fmt.Errorf("silnik przeglądarki: nieczytelny zrzut: %w", err)
	}
	if odpowiedz.Dane == "" {
		return "", 0, 0, errors.New("silnik przeglądarki: przeglądarka oddała zrzut bez treści")
	}
	return odpowiedz.Dane, szerokoscZrzutu, wysokoscZrzutu, nil
}

// ramkaElementu mierzy położenie i rozmiar elementu wskazanego selektorem.
func (s *sesjaStrony) ramkaElementu(ctx context.Context, selektor string) ([4]float64, error) {
	wyrazenie := `(() => { const e = document.querySelector(` + strconv.Quote(selektor) + `);
		if (!e) return null; const r = e.getBoundingClientRect();
		return [r.x + window.scrollX, r.y + window.scrollY, r.width, r.height]; })()`
	wartosc, err := s.ocenNaStronie(ctx, wyrazenie)
	if err != nil {
		return [4]float64{}, err
	}
	var ramka []float64
	if err := json.Unmarshal(wartosc, &ramka); err != nil || len(ramka) != 4 {
		return [4]float64{}, errors.New("silnik przeglądarki: na stronie nie ma elementu o wskazaniu " + selektor)
	}
	if ramka[2] <= 0 || ramka[3] <= 0 {
		return [4]float64{}, errors.New("silnik przeglądarki: element " + selektor + " nie ma wymiarów na stronie")
	}
	return [4]float64{ramka[0], ramka[1], ramka[2], ramka[3]}, nil
}
