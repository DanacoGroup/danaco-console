// Adapter obsługuje browser.accessibility.audit: potwierdza dostępność strony,
// uruchamia program audytu na przeglądarce silnika i zwraca naruszenia z liczbami
// wag oraz wersją programu.
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
// miejscu użycia; wykaz zależności odwołuje się do niej, zamiast powtarzać
// nazwę programu.
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
	// progZgloszenPa11y zdejmuje z programu prawo kończenia się kodem niezerowym
	// z powodu znalezionych zgłoszeń; rdzeń traktuje jako niepowodzenie wyłącznie
	// niepowodzenie badania.
	progZgloszenPa11y = "1000000"
)

// ZbadajDostepnosc obsługuje browser.accessibility.audit: potwierdza dostępność
// strony, uruchamia program audytu i zwraca wykaz zgłoszeń wraz z licznikami wag.
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

	// Potwierdzenie pomiaru: pusty wykaz strony nieosiągalnej byłby brakiem
	// pomiaru podanym za pomiar.
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

// wolajAudytDostepnosci uruchamia program audytu na przeglądarce wskazanej
// nastawami i oddaje jego surowy wynik wraz z ewentualnym błędem uruchomienia.
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
		// Próg zdejmuje kod niezerowy za znalezione zgłoszenia; niepowodzenie
		// badania kończy program błędem.
		"--threshold", progZgloszenPa11y,
		// Granica programu jest krótsza od granicy arsenału, żeby przekroczenie
		// nazwał najpierw program.
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
// podaną programowi, żeby przekroczenie nazywał zawsze program, a nie arsenał.
const granicaZapasuAudytu = 30 * time.Second

// normaAudytu rozstrzyga normę audytu i odmawia wartości spoza wyliczenia
// kontraktu, ponieważ ta droga nie przechodzi przez bramę sprawdzającą
// wyliczenia komend.
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
// wyniku, bo wykaz naruszeń zależy od wydania reguł programu audytującego.
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
