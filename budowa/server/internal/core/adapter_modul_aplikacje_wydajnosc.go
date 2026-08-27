// Odpowiedzialność pliku: audyt wydajności strony produktu —
// `apps.performance.audit`.
//
// ── Dostępność to nie pomiar ─────────────────────────────────────────────────
// Moduł umiał dotąd powiedzieć o wdrożonym produkcie jedno: czy odpowiada
// (`apps.deployment.health.get` — dostępność, czas nieprzerwanego działania,
// wynik ostatniego sprawdzenia kondycji). To jest odpowiedź na pytanie „czy
// stoi", nie na pytanie „jak szybko się otwiera". Produkt, który odpowiada
// w cztery sekundy, jest dostępny w stu procentach i nie do użycia.
//
// ── Dlaczego programem, a nie własnym stoperem ───────────────────────────────
// Core Web Vitals nie są czasem odpowiedzi serwera. Największe wymalowanie
// treści i przesunięcia układu powstają w przeglądarce, po wykonaniu skryptów,
// a ich wartość zależy od emulacji urządzenia i dławienia sieci. Rdzeń, który
// mierzyłby to własnym `net/http`, oddałby czas pobrania dokumentu i nazwał go
// wydajnością strony — liczbę prawdziwą, odpowiadającą na inne pytanie.
//
// ── Pomiar, który się nie odbył, nie wychodzi jako wynik ─────────────────────
// Program pomiarowy mówi o tym wprost: przebieg, w którym strona się nie
// wczytała, niesie w odpowiedzi pole `runtimeError` wraz z kodem powodu i NIE
// niesie ocen. Rdzeń czyta to pole przed czymkolwiek innym — bez tego odczytu
// odpowiedź o produkcie, którego pod adresem nie ma, składałaby się z samych
// zer i wyglądałaby jak strona wolna, a nie jak strona niezmierzona.
//
// ── Granica czasu jest jawna, a jej przekroczenie nazywa się przekroczeniem ──
// Audyt trwa kilkanaście sekund przy stronie zdrowej i nie kończy się nigdy
// przy stronie, która nie przestaje się wczytywać. Granica idzie do programu
// (`--max-wait-for-load`) i osobno do arsenału, z zapasem — pierwszy mija
// program, więc przekroczenie nazywa ten, kto wie, na co czekał.
package core

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/session"
	"danacoconsole/server/internal/zewnetrzne"
	"danacoconsole/shared"
)

// narzedzieLighthouse opisuje program audytu wydajności. Deklaracja stoi przy
// miejscu użycia; wykaz zależności odwołuje się do niej, zamiast powtarzać
// nazwę po raz drugi.
var narzedzieLighthouse = zewnetrzne.Narzedzie{
	Nazwa: "Lighthouse", Program: "lighthouse", Pakiet: "npm i -g lighthouse"}

const (
	// granicaAudytuWydajnosci obejmuje start przeglądarki, wczytanie strony,
	// zebranie śladu i policzenie miar.
	granicaAudytuWydajnosci = 3 * time.Minute
	// najdluzszyAudytWydajnosci jest granicą, której żądanie nie przekroczy
	// nawet wtedy, gdy poprosi o więcej.
	najdluzszyAudytWydajnosci = 10 * time.Minute
	// granicaZapasuWydajnosci jest zapasem, o który granica arsenału przewyższa
	// granicę wczytania podaną programowi. Program potrzebuje czasu na policzenie
	// miar PO wczytaniu strony; bez zapasu arsenał ubijałby go w połowie liczenia
	// i przekroczenie wyglądałoby jak awaria.
	granicaZapasuWydajnosci = 60 * time.Second
)

// miaryWydajnosciStrony wylicza miary, o które moduł pyta, w kolejności
// ustalonej. Kolejność jest ustalona, żeby dwa kolejne audyty tej samej strony
// dawały wykaz w tym samym porządku — wynik ma się różnić wtedy, gdy zmieniła
// się strona, a nie wtedy, gdy inaczej ułożyła się mapa odpowiedzi programu.
//
// Wykaz jest zamknięty i obejmuje Core Web Vitals wraz z miarami, z których te
// się liczą. Miara dopisana tu bez pokrycia w odpowiedzi programu wyszłaby
// z audytu jako zero — dlatego brak miary w odpowiedzi pomija się, zamiast
// wypełniać wartością zastępczą.
var miaryWydajnosciStrony = []string{
	"first-contentful-paint",
	"largest-contentful-paint",
	"total-blocking-time",
	"cumulative-layout-shift",
	"speed-index",
	"interactive",
}

// ZmierzWydajnosc obsługuje `apps.performance.audit`.
func (a *adapterAplikacji) ZmierzWydajnosc(ctx context.Context,
	z shared.AppsPerformanceAuditRequest) (shared.AppsPerformanceAuditResponse, error) {

	okno, err := oknoAplikacji(z.WindowId, "apps.performance.audit")
	if err != nil {
		return shared.AppsPerformanceAuditResponse{}, err
	}
	if err := a.sprawdzOknoWdrozenia(okno); err != nil {
		return shared.AppsPerformanceAuditResponse{}, err
	}
	adres := strings.TrimSpace(z.Url)
	if adres == "" {
		return shared.AppsPerformanceAuditResponse{}, bladWskazaniaAplikacji(
			"apps.performance.audit wymaga adresu mierzonej strony")
	}
	if !strings.HasPrefix(adres, "http://") && !strings.HasPrefix(adres, "https://") {
		return shared.AppsPerformanceAuditResponse{}, bladWskazaniaAplikacji(
			"adres mierzonej strony ma zaczynać się od http:// albo https://; otrzymano " + adres)
	}
	postac, err := postacUrzadzeniaAudytu(z.FormFactor)
	if err != nil {
		return shared.AppsPerformanceAuditResponse{}, err
	}
	if a.uruchamiacz == nil {
		return shared.AppsPerformanceAuditResponse{}, protocol.JakoError(
			protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
				"moduł Apps: rdzeń nie ma uruchamiacza procesów — audyt wydajności nie ma czym "+
					"wystartować; naprawa: podpiąć warstwę kanału (injection) przy składaniu rdzenia"))
	}

	granica := granicaAudytuWydajnosci
	if z.TimeoutMs != nil && *z.TimeoutMs > 0 {
		granica = time.Duration(*z.TimeoutMs) * time.Millisecond
		if granica > najdluzszyAudytWydajnosci {
			granica = najdluzszyAudytWydajnosci
		}
	}

	argumenty := []string{
		adres,
		"--quiet",
		"--output=json",
		"--output-path=stdout",
		"--only-categories=performance",
		"--form-factor=" + postac,
		"--max-wait-for-load=" + strconv.FormatInt(granica.Milliseconds(), 10),
		"--chrome-flags=--headless=new --no-sandbox --disable-gpu --disable-dev-shm-usage " +
			"--no-first-run --no-default-browser-check",
	}
	// Postać biurkowa ma w programie własną nastawę zbiorczą: sam `--form-factor`
	// zmienia sposób liczenia oceny, lecz zostawia emulację i dławienie telefonu,
	// więc wynik byłby oceną biurka policzoną na warunkach telefonu.
	if postac == shared.AppPerformanceFormFactorDesktop {
		argumenty = append(argumenty, "--preset=desktop")
	}

	granicaArsenalu := granica + granicaZapasuWydajnosci
	poczatek := time.Now()
	oknoProcesu, zasady, obszar := a.zasiegAplikacji(okno)
	wynik, err := zewnetrzne.Wolaj(ctx, a.uruchamiacz, oknoProcesu, zasady, obszar,
		narzedzieLighthouse, argumenty, obszar.KatalogRoboczy, granicaArsenalu)
	trwanie := time.Since(poczatek)
	// Program kończy się kodem niezerowym także wtedy, gdy pomiar się nie odbył,
	// a powód opisał w odpowiedzi. Odpowiedź czytamy więc PRZED rozpatrzeniem
	// odmowy arsenału — inaczej „strony nie ma pod tym adresem" wyszłoby jako
	// „program zakończył się niepowodzeniem".
	raport, bladOdczytu := odczytajRaportWydajnosci(wynik.Wyjscie)
	if bladOdczytu != nil {
		if err != nil {
			return shared.AppsPerformanceAuditResponse{}, bladProgramuAplikacji(
				"apps.performance.audit", err, trwanie, granicaArsenalu)
		}
		return shared.AppsPerformanceAuditResponse{}, bladOdczytu
	}
	if powod := raport.powodNieodbytegoPomiaru(); powod != "" {
		return shared.AppsPerformanceAuditResponse{}, protocol.JakoError(
			protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
				"moduł Apps: audyt wydajności strony "+adres+" nie odbył się: "+powod+
					". Ocena zerowa byłaby tu odpowiedzią o stronie, której nikt nie zmierzył"))
	}
	if err != nil {
		return shared.AppsPerformanceAuditResponse{}, bladProgramuAplikacji(
			"apps.performance.audit", err, trwanie, granicaArsenalu)
	}

	audyt, err := raport.jakoAudyt(adres, postac, poczatek)
	if err != nil {
		return shared.AppsPerformanceAuditResponse{}, err
	}
	a.dopiszDziennikApp(ctx, okno, nil, nil,
		"audyt wydajności strony "+audyt.Url+": ocena "+strconv.Itoa(audyt.PerformanceScore)+
			" na 100, miar zmierzonych "+strconv.Itoa(len(audyt.Metrics)))
	return shared.AppsPerformanceAuditResponse{Audit: audyt}, nil
}

// raportWydajnosci jest tą częścią odpowiedzi programu, którą moduł czyta.
// Reszta odpowiedzi — ślad przeglądarki, zrzut strony, wykaz wszystkich
// sprawdzeń — do kontraktu nie wchodzi i nie ma powodu przechodzić przez pamięć
// rdzenia w postaci rozebranej.
type raportWydajnosci struct {
	Wersja  string `json:"lighthouseVersion"`
	Adres   string `json:"finalDisplayedUrl"`
	Usterka *struct {
		Kod       string `json:"code"`
		Komunikat string `json:"message"`
	} `json:"runtimeError"`
	Kategorie map[string]struct {
		Ocena *float64 `json:"score"`
	} `json:"categories"`
	Sprawdzenia map[string]struct {
		Nazwa      string   `json:"title"`
		Ocena      *float64 `json:"score"`
		Wartosc    *float64 `json:"numericValue"`
		Jednostka  string   `json:"numericUnit"`
		WartoscOpi string   `json:"displayValue"`
	} `json:"audits"`
}

// odczytajRaportWydajnosci rozbiera odpowiedź programu. Wyjście, które
// odpowiedzią nie jest, znaczy pomiar nieodbyty — nie pomiar pusty.
func odczytajRaportWydajnosci(wyjscie []byte) (raportWydajnosci, error) {
	var raport raportWydajnosci
	if err := json.Unmarshal(wyjscie, &raport); err != nil {
		return raportWydajnosci{}, protocol.JakoError(
			protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
				"moduł Apps: program audytu wydajności nie oddał raportu, więc pomiar się nie "+
					"odbył: "+err.Error()))
	}
	return raport, nil
}

// powodNieodbytegoPomiaru oddaje zdanie programu o tym, dlaczego pomiaru nie
// ma. Puste znaczy, że program mierzył.
func (r raportWydajnosci) powodNieodbytegoPomiaru() string {
	if r.Usterka != nil && strings.TrimSpace(r.Usterka.Kod) != "" {
		return r.Usterka.Kod + " — " + strings.TrimSpace(r.Usterka.Komunikat)
	}
	if _, jest := r.Kategorie["performance"]; !jest {
		return "raport nie niesie oceny wydajności"
	}
	if r.Kategorie["performance"].Ocena == nil {
		return "raport niesie ocenę wydajności bez wartości"
	}
	return ""
}

// jakoAudyt składa wynik kontraktu z raportu programu.
func (r raportWydajnosci) jakoAudyt(adres, postac string,
	poczatek time.Time) (shared.AppPerformanceAudit, error) {

	miary := make([]shared.AppPerformanceMetric, 0, len(miaryWydajnosciStrony))
	for _, klucz := range miaryWydajnosciStrony {
		sprawdzenie, jest := r.Sprawdzenia[klucz]
		if !jest || sprawdzenie.Wartosc == nil {
			// Miara, której program nie policzył, NIE wchodzi z wartością zero:
			// zero jest w tych miarach wynikiem najlepszym z możliwych.
			continue
		}
		miara := shared.AppPerformanceMetric{
			Id:    klucz,
			Title: sprawdzenie.Nazwa,
			Value: *sprawdzenie.Wartosc,
			Unit:  sprawdzenie.Jednostka,
		}
		if sprawdzenie.Ocena != nil {
			ocena := wSkaliStu(*sprawdzenie.Ocena)
			miara.Score = &ocena
		}
		if sprawdzenie.WartoscOpi != "" {
			miara.DisplayValue = wskaznikNapisuApp(sprawdzenie.WartoscOpi)
		}
		miary = append(miary, miara)
	}
	if len(miary) == 0 {
		return shared.AppPerformanceAudit{}, protocol.JakoError(
			protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
				"moduł Apps: raport audytu wydajności nie niesie ani jednej miary Core Web "+
					"Vitals dla "+adres+", więc pomiar się nie odbył"))
	}

	zmierzony := strings.TrimSpace(r.Adres)
	if zmierzony == "" {
		zmierzony = adres
	}
	return shared.AppPerformanceAudit{
		Url:              zmierzony,
		FormFactor:       shared.AppPerformanceFormFactor(postac),
		PerformanceScore: wSkaliStu(*r.Kategorie["performance"].Ocena),
		Metrics:          miary,
		ToolVersion:      strings.TrimSpace(r.Wersja),
		StartedAt:        poczatek.UnixMilli(),
		FinishedAt:       time.Now().UnixMilli(),
	}, nil
}

// wSkaliStu przekłada ocenę programu (ułamek od zera do jedynki) na skalę
// setną kontraktu. Zaokrąglenie w dół dałoby 99 dla strony ocenionej idealnie,
// więc idzie do najbliższej.
func wSkaliStu(ocena float64) int {
	if ocena < 0 {
		return 0
	}
	if ocena > 1 {
		return 100
	}
	return int(ocena*100 + 0.5)
}

// postacUrzadzeniaAudytu rozstrzyga postać urządzenia i odmawia wartości spoza
// wyliczenia. Odmowa jest tu potrzebna mimo bramy kontraktu: brama sprawdza
// wartości wyliczeń wyłącznie przy komendach wystawionych jako narzędzia modelu,
// a ta do nich nie należy.
func postacUrzadzeniaAudytu(wskazanie *shared.AppPerformanceFormFactor) (string, error) {
	if wskazanie == nil || strings.TrimSpace(string(*wskazanie)) == "" {
		return shared.AppPerformanceFormFactorDesktop, nil
	}
	szukana := strings.TrimSpace(string(*wskazanie))
	for _, znana := range shared.WartosciAppPerformanceFormFactor() {
		if string(znana) == szukana {
			return szukana, nil
		}
	}
	return "", bladWskazaniaAplikacji(
		"postać urządzenia " + szukana + " nie należy do kontraktu")
}

// zasiegAplikacji składa okno, zasady izolacji i obszar dla programu wołanego
// przez moduł — tak samo jak robią to pozostałe rodziny wołające programy.
func (a *adapterAplikacji) zasiegAplikacji(oknoKod string) (session.Okno, session.Zasady, session.Obszar) {
	okno := session.Okno{Id: oknoKod, Ustawienia: session.Ustawienia{
		SrodowiskoWykonania: shared.ExecutionEnvCore,
	}}
	if a.okna != nil {
		if znalezione, err := a.okna.Okno(oknoKod); err == nil {
			okno = znalezione
		}
	}
	zasady := session.Zasady{}
	obszar := session.Obszar{IdOkna: okno.Id}
	if a.rozstrzygacz != nil {
		zasady = ZasadyIzolacji(a.rozstrzygacz, konfig.Kontekst{Okno: okno.Id})
	}
	if a.katalog != nil {
		obszar = ObszarOkna(a.katalog.Ustal(konfig.Kontekst{Okno: okno.Id}, okno.IdSesji), okno.Id)
	}
	return okno, zasady, obszar
}

// bladProgramuAplikacji odróżnia brak programu i przekroczenie granicy czasu od
// usterki rdzenia.
//
// Brak programu jest brakiem, który Operator serwera usuwa jedną instalacją,
// a przekroczenie granicy jest przekroczeniem — nie awarią. Obie sytuacje bez
// tego rozróżnienia wychodziłyby jako `internal_error`, czyli zdanie „usterka
// rdzenia, zgłoś ją", mówiące czytającemu coś nieprawdziwego o tym, co się
// stało (wzór: `odmowaSkanuSane`).
//
// Przekroczenie rozstrzyga się ZMIERZONYM czasem, nie treścią komunikatu:
// zdanie, którym arsenał opisuje przerwanie, jest napisem i przy następnej
// zmianie brzmiałoby inaczej, a stoper mierzy to samo, o co tu chodzi.
func bladProgramuAplikacji(komenda string, err error, trwanie, granica time.Duration) error {
	if err == nil {
		return nil
	}
	var brak *zewnetrzne.BrakNarzedzia
	if errors.As(err, &brak) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
			"moduł Apps: komenda "+komenda+" nie ma czym zmierzyć strony: "+err.Error()))
	}
	if granica > 0 && trwanie >= granica {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
			"moduł Apps: komenda "+komenda+" przekroczyła granicę czasu "+granica.String()+
				" — to jest przekroczenie granicy, nie usterka rdzenia; naprawa: podnieść "+
				"granicę polem timeoutMs albo wskazać stronę, która wczytuje się szybciej. "+
				"Diagnostyka warstwy: "+err.Error()))
	}
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
		"moduł Apps: komenda "+komenda+" nie doszła do skutku: "+err.Error()))
}
