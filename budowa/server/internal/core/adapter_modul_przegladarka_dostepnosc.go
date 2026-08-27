// Odpowiedzialność pliku: audyt dostępności bieżącej strony okna —
// `browser.accessibility.audit`.
//
// ── Czwarta sonda strony uruchomionej ────────────────────────────────────────
// Moduł czyta stronę żywą trzema sondami: drzewem elementów
// (`browser.dom.inspect`), rejestrem żądań (`browser.network.har`) i konsolą
// (`browser.console.read`). Audyt dostępności jest czwartą i pyta o to samo, co
// tamte trzy — o stronę PO zbudowaniu przez skrypty, nie o jej źródło. Reguła
// dostępności orzeka o etykiecie kontrolki wstawionej skryptem tak samo jak
// o etykiecie wpisanej w źródle; audyt czytający sam HTML odpowiedziałby
// o dokumencie, którego Operator nigdy nie ogląda.
//
// Adres bierze się z ostatniej migawki okna — tak samo jak w trzech sondach
// starszych. Kontrakt nie niesie w tym żądaniu pola adresu, bo pyta o bieżącą
// stronę okna, a bieżącą stroną okna jest to, dokąd okno ostatnio przeszło.
//
// ── Czym jest badane, skoro rdzeń ma własny silnik ───────────────────────────
// Reguły WCAG są cudzą wiedzą i rdzeń jej nie przepisuje: między normą a jej
// sprawdzeniem stoją setki reguł, które ktoś utrzymuje wraz z kolejnymi
// wydaniami normy. Dlatego audyt idzie programem (`pa11y`), a nie własnym
// obchodem drzewa. Program dostaje TĘ SAMĄ przeglądarkę, którą rdzeń już
// deklaruje dla sond starszych (`narzedzieChromium`) — nie własną kopię
// pobieraną z sieci przy pierwszym uruchomieniu.
//
// ── Zero naruszeń jest wynikiem dopiero po potwierdzeniu pomiaru ─────────────
// Program audytujący nie mówi, czy strona się wczytała: dokument błędu 404
// bywa poprawny wobec normy i wychodzi z audytu jako pusty wykaz naruszeń.
// Odpowiedź „zero naruszeń” dla strony, której pod tym adresem nie ma, byłaby
// brakiem pomiaru podanym jako pomiar — i to najgorszą jego postacią, bo
// wygląda dobrze i nie wzywa nikogo do sprawdzenia. Dlatego przed audytem
// rdzeń sięga po stronę własną drogą modułu (`pobierzStrone`), która orzeka
// o stanie odpowiedzi i o tym, czy zasób jest w ogóle stroną. Odmowa stąd
// nazywa, czego nie zmierzono, zamiast podawać zero.
//
// Drugie potwierdzenie jest po stronie odpowiedzi programu: pusty wykaz jest
// wynikiem tylko wtedy, gdy program oddał tablicę JSON. Wyjście, które tablicą
// nie jest, znaczy audyt nieodbyty i wraca odmową.
package core

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/zewnetrzne"
	"danacoconsole/shared"
)

// narzedziePa11y opisuje program audytu dostępności. Deklaracja stoi przy
// miejscu użycia, tak jak Pandoc przy dokumentach i ffmpeg przy nagraniach;
// wykaz zależności (`zaleznosci_zewnetrzne.go`) odwołuje się do niej, zamiast
// powtarzać nazwę po raz drugi.
var narzedziePa11y = zewnetrzne.Narzedzie{
	Nazwa: "Pa11y", Program: "pa11y", Pakiet: "npm i -g pa11y"}

const (
	// granicaAudytuDostepnosci obejmuje start przeglądarki, wczytanie strony
	// i przebieg reguł. Hojna, bo obejmuje uruchomienie procesu; skończona, bo
	// strona wisząca bez końca jest zjawiskiem codziennym.
	granicaAudytuDostepnosci = 3 * time.Minute
	// najdluzszyAudytDostepnosci jest granicą, której żądanie nie przekroczy
	// nawet wtedy, gdy poprosi o więcej. Audyt wiszący pół godziny trzyma
	// połączenie klienta i proces przeglądarki, a wynik, na który nikt już nie
	// czeka, nie jest wynikiem.
	najdluzszyAudytDostepnosci = 10 * time.Minute
	// progZgloszenPa11y zdejmuje z programu prawo kończenia się kodem
	// niezerowym z powodu ZNALEZIONYCH zgłoszeń. Program odróżnia „znalazłem
	// naruszenia" (kod 2) od „nie dałem rady zbadać" (kod 1); rdzeń potrzebuje
	// wyłącznie tego drugiego jako niepowodzenia, bo pierwsze JEST wynikiem.
	progZgloszenPa11y = "1000000"
)

// ZbadajDostepnosc obsługuje `browser.accessibility.audit`.
func (a *adapterPrzegladarki) ZbadajDostepnosc(ctx context.Context,
	z shared.BrowserAccessibilityAuditRequest) (shared.BrowserAccessibilityAuditResponse, error) {

	migawka, err := a.adresOstatniejStrony(ctx, z.WindowId, "accessibility.audit")
	if err != nil {
		return shared.BrowserAccessibilityAuditResponse{}, err
	}
	norma, err := normaAudytu(z.Standard)
	if err != nil {
		return shared.BrowserAccessibilityAuditResponse{}, err
	}
	if a.silnik == nil || a.silnik.uruchamiacz == nil {
		return shared.BrowserAccessibilityAuditResponse{}, protocol.JakoError(
			protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
				"moduł Browser: rdzeń nie ma uruchamiacza procesów — audyt dostępności nie ma "+
					"czym wystartować; naprawa: podpiąć warstwę kanału (injection) przy "+
					"składaniu rdzenia"))
	}

	// Potwierdzenie pomiaru PRZED audytem: zero naruszeń na stronie, której pod
	// tym adresem nie ma, jest brakiem pomiaru podanym jako pomiar.
	if _, err := pobierzStrone(ctx, migawka.Url); err != nil {
		kod, zdanie := kodOdmowyPobrania(err)
		return shared.BrowserAccessibilityAuditResponse{}, protocol.JakoError(protocol.NowyBlad(kod,
			"moduł Browser: audyt dostępności nie odbył się, bo strona "+migawka.Url+
				" nie dała się pobrać: "+zdanie+". Wykaz pusty byłby tu odpowiedzią o stronie, "+
				"której nikt nie zmierzył"))
	}

	wyjscie, err := a.wolajAudytDostepnosci(ctx, migawka.Url, norma, z)
	if err != nil {
		return shared.BrowserAccessibilityAuditResponse{}, err
	}

	var odczyt []struct {
		Code     string `json:"code"`
		Type     string `json:"type"`
		Message  string `json:"message"`
		Context  string `json:"context"`
		Selector string `json:"selector"`
		Runner   string `json:"runner"`
	}
	if err := json.Unmarshal(wyjscie.Wyjscie, &odczyt); err != nil {
		return shared.BrowserAccessibilityAuditResponse{}, protocol.JakoError(
			protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
				"moduł Browser: program audytu dostępności nie oddał wykazu zgłoszeń dla "+
					migawka.Url+", więc audyt się nie odbył — pusty wykaz byłby tu wynikiem "+
					"zmyślonym"+opisWyjsciaProgramu(wyjscie)))
	}

	zgloszenia := make([]shared.BrowserAccessibilityIssue, 0, len(odczyt))
	naruszen, ostrzezen, uwag := 0, 0, 0
	for _, wpis := range odczyt {
		poziom := poziomDostepnosci(wpis.Type)
		switch poziom {
		case shared.BrowserAccessibilityLevelWarning:
			ostrzezen++
		case shared.BrowserAccessibilityLevelNotice:
			uwag++
		default:
			naruszen++
		}
		zgloszenie := shared.BrowserAccessibilityIssue{
			Code:    wpis.Code,
			Level:   shared.BrowserAccessibilityLevel(poziom),
			Message: wpis.Message,
		}
		if wpis.Selector != "" {
			zgloszenie.Selector = wskaznikTekstu(wpis.Selector)
		}
		if wpis.Context != "" {
			zgloszenie.Context = wskaznikTekstu(wpis.Context)
		}
		if wpis.Runner != "" {
			zgloszenie.Runner = wskaznikTekstu(wpis.Runner)
		}
		zgloszenia = append(zgloszenia, zgloszenie)
	}

	return shared.BrowserAccessibilityAuditResponse{
		Issues:       zgloszenia,
		Standard:     shared.BrowserAccessibilityStandard(norma),
		Url:          migawka.Url,
		ErrorCount:   naruszen,
		WarningCount: ostrzezen,
		NoticeCount:  uwag,
		ToolVersion:  wersjaProgramuAudytu(ctx, a.silnik),
		AuditedAt:    time.Now().UnixMilli(),
	}, nil
}

// wolajAudytDostepnosci uruchamia program audytu i oddaje jego surowy wynik.
//
// Program dostaje przeglądarkę WSKAZANĄ, nie szukaną: jego własna warstwa
// sterowania przeglądarką pobiera wydanie Chrome do katalogu pamięci podręcznej
// użytkownika, a rdzeń takiego pobrania nie robi i nie ma prawa go wymagać od
// wdrożenia. Wskazanie idzie plikiem nastaw, bo jedyna droga do procesu
// (`zewnetrzne.Wolaj`) nie przekazuje zmiennych środowiska — i ma nie
// przekazywać, bo binarium arsenału nie ma powodu widzieć zmiennych rdzenia.
func (a *adapterPrzegladarki) wolajAudytDostepnosci(ctx context.Context, adres, norma string,
	z shared.BrowserAccessibilityAuditRequest) (zewnetrzne.Wynik, error) {

	przegladarka := narzedzieChromium()
	sciezkaPrzegladarki, jest := zewnetrzne.Odnajdz(przegladarka)
	if !jest {
		return zewnetrzne.Wynik{}, bladSilnikaPrzegladarki("browser.accessibility.audit",
			&zewnetrzne.BrakNarzedzia{Narzedzie: przegladarka})
	}

	katalog, err := os.MkdirTemp("", "danaco-dostepnosc-")
	if err != nil {
		return zewnetrzne.Wynik{}, bladZapleczaPrzegladarki(
			"nie można założyć katalogu nastaw audytu dostępności: " + err.Error())
	}
	defer os.RemoveAll(katalog)

	nastawy, err := json.Marshal(map[string]any{
		"chromeLaunchConfig": map[string]any{
			"executablePath": sciezkaPrzegladarki,
			"args": []string{
				"--no-sandbox", "--disable-gpu", "--disable-dev-shm-usage",
				"--no-first-run", "--no-default-browser-check",
			},
		},
	})
	if err != nil {
		return zewnetrzne.Wynik{}, bladZapleczaPrzegladarki(
			"nie można złożyć nastaw audytu dostępności: " + err.Error())
	}
	plikNastaw := filepath.Join(katalog, "nastawy.json")
	if err := os.WriteFile(plikNastaw, nastawy, 0o600); err != nil {
		return zewnetrzne.Wynik{}, bladZapleczaPrzegladarki(
			"nie można zapisać nastaw audytu dostępności: " + err.Error())
	}

	granica := granicaAudytuDostepnosci
	if z.TimeoutMs != nil && *z.TimeoutMs > 0 {
		granica = time.Duration(*z.TimeoutMs) * time.Millisecond
		if granica > najdluzszyAudytDostepnosci {
			granica = najdluzszyAudytDostepnosci
		}
	}

	argumenty := []string{
		"--config", plikNastaw,
		"--reporter", "json",
		"--standard", strings.ToUpper(norma),
		// Próg zdejmuje kod niezerowy za ZNALEZIONE zgłoszenia; „nie dałem rady
		// zbadać" nadal kończy program kodem błędu i dojdzie tu jako odmowa.
		"--threshold", progZgloszenPa11y,
		// Granica programu jest krótsza od granicy arsenału, żeby przekroczenie
		// nazwał najpierw ten, kto wie, co robił — a nie zabicie drzewa procesów.
		"--timeout", strconv.FormatInt(granica.Milliseconds(), 10),
	}
	if z.IncludeWarnings != nil && *z.IncludeWarnings {
		argumenty = append(argumenty, "--include-warnings")
	}
	if z.IncludeNotices != nil && *z.IncludeNotices {
		argumenty = append(argumenty, "--include-notices")
	}
	argumenty = append(argumenty, adres)

	okno, zasady, obszar := a.silnik.zasiegSilnika()
	wynik, err := zewnetrzne.Wolaj(ctx, a.silnik.uruchamiacz, okno, zasady, obszar,
		narzedziePa11y, argumenty, obszar.KatalogRoboczy, granica+granicaZapasuAudytu)
	if err != nil {
		return zewnetrzne.Wynik{}, bladSilnikaPrzegladarki("browser.accessibility.audit", err)
	}
	return wynik, nil
}

// granicaZapasuAudytu jest zapasem, o który granica arsenału przewyższa granicę
// podaną programowi. Bez zapasu obie granice mijałyby w tej samej chwili
// i przekroczenie nazywałoby się raz przekroczeniem programu, raz zabiciem
// drzewa procesów — zależnie od tego, kto zdążył pierwszy.
const granicaZapasuAudytu = 30 * time.Second

// normaAudytu rozstrzyga normę audytu i odmawia wartości spoza wyliczenia.
//
// Odmowa jest tu potrzebna mimo bramy kontraktu: brama sprawdza wartości
// wyliczeń wyłącznie przy komendach wystawionych jako narzędzia modelu
// (`brama_kontraktu.go`), a ta do nich nie należy. Wartość nierozpoznana
// przekazana programowi wróciłaby jego własnym komunikatem o nieznanej normie.
func normaAudytu(wskazanie *shared.BrowserAccessibilityStandard) (string, error) {
	if wskazanie == nil || strings.TrimSpace(string(*wskazanie)) == "" {
		return shared.BrowserAccessibilityStandardWcag2aa, nil
	}
	szukana := strings.TrimSpace(string(*wskazanie))
	for _, znana := range shared.WartosciBrowserAccessibilityStandard() {
		if string(znana) == szukana {
			return szukana, nil
		}
	}
	return "", bladWskazaniaPrzegladarki(
		"norma audytu dostępności " + szukana + " nie należy do kontraktu")
}

// poziomDostepnosci przekłada wagę zgłoszenia programu na wyliczenie kontraktu.
// Nazwa nieznana spada na naruszenie: zgłoszenie ma się pokazać w najcięższej
// postaci, a nie zniknąć przez nieznane słowo w polu wagi.
func poziomDostepnosci(nazwa string) string {
	switch strings.ToLower(strings.TrimSpace(nazwa)) {
	case "warning":
		return shared.BrowserAccessibilityLevelWarning
	case "notice":
		return shared.BrowserAccessibilityLevelNotice
	default:
		return shared.BrowserAccessibilityLevelError
	}
}

// wersjaProgramuAudytu pyta program o jego własną wersję. Odpowiedź wchodzi do
// wyniku, bo wykaz naruszeń jest orzeczeniem zestawu reguł konkretnego wydania —
// dwa wydania programu potrafią policzyć tę samą stronę inaczej, a wynik bez
// wersji nie daje się porównać z wynikiem sprzed miesiąca.
func wersjaProgramuAudytu(ctx context.Context, silnik *silnikPrzegladarki) string {
	if silnik == nil || silnik.uruchamiacz == nil {
		return ""
	}
	okno, zasady, obszar := silnik.zasiegSilnika()
	wynik, err := zewnetrzne.Wolaj(ctx, silnik.uruchamiacz, okno, zasady, obszar,
		narzedziePa11y, []string{"--version"}, obszar.KatalogRoboczy, granicaSondyWersji)
	if err != nil {
		return ""
	}
	return pierwszyWierszWersji(string(wynik.Wyjscie))
}

// granicaSondyWersji jest granicą pytania programu o jego wersję. Krótka, bo
// program, który na to pytanie nie odpowiada w kilka sekund, nie odpowie wcale,
// a sonda wersji nie ma prawa opóźniać audytu.
const granicaSondyWersji = 20 * time.Second

// opisWyjsciaProgramu dokłada do odmowy początek tego, co program powiedział.
// Bez tego członu czytający odmowę wie tylko, że wykazu nie ma.
func opisWyjsciaProgramu(wynik zewnetrzne.Wynik) string {
	tresc := strings.TrimSpace(wynik.Diagnostyka)
	if tresc == "" {
		tresc = strings.TrimSpace(string(wynik.Wyjscie))
	}
	if tresc == "" {
		return ""
	}
	if len(tresc) > 400 {
		tresc = tresc[:400] + "…"
	}
	return "; program powiedział: " + tresc
}
