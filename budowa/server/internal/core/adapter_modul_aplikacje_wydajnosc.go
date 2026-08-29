// Audyt wydajności strony produktu obsługuje apps.performance.audit. Core Web
// Vitals mierzy zewnętrzny program pomiarowy, nie własny stoper rdzenia.
// Przebieg, w którym strona się nie wczytała, niesie pole runtimeError i nie
// niesie ocen.
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
	// granicę wczytania podaną programowi, na policzenie miar po wczytaniu.
	granicaZapasuWydajnosci = 60 * time.Second
)

// miaryWydajnosciStrony wylicza miary, o które moduł pyta, w kolejności
// ustalonej, żeby dwa kolejne audyty tej samej strony dawały wykaz w tym
// samym porządku. Wykaz jest zamknięty i obejmuje Core Web Vitals.
var miaryWydajnosciStrony = []string{
	"first-contentful-paint",
	"largest-contentful-paint",
	"total-blocking-time",
	"cumulative-layout-shift",
	"speed-index",
	"interactive",
}

// ZmierzWydajnosc obsługuje apps.performance.audit: uruchamia program
// pomiarowy na wskazanej stronie i przekłada jego raport na wynik kontraktu.
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
				"moduł Apps: serwer nie ma uruchamiacza procesów — audyt wydajności nie ma czym "+
					"wystartować; naprawa: podpiąć warstwę kanału (injection) przy składaniu serwera"))
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
	// Postać biurkowa wymaga w programie osobnej nastawy zbiorczej.
	if postac == shared.AppPerformanceFormFactorDesktop {
		argumenty = append(argumenty, "--preset=desktop")
	}

	granicaArsenalu := granica + granicaZapasuWydajnosci
	poczatek := time.Now()
	oknoProcesu, zasady, obszar := a.zasiegAplikacji(okno)
	wynik, err := zewnetrzne.Wolaj(ctx, a.uruchamiacz, oknoProcesu, zasady, obszar,
		narzedzieLighthouse, argumenty, obszar.KatalogRoboczy, granicaArsenalu)
	trwanie := time.Since(poczatek)
	// Odpowiedź programu jest czytana przed rozpatrzeniem jego odmowy zakończenia.
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

// jakoAudyt składa wynik kontraktu z raportu programu, pomijając miary, które
// program nie policzył, i zaokrąglając ocenę do skali setnej.
func (r raportWydajnosci) jakoAudyt(adres, postac string,
	poczatek time.Time) (shared.AppPerformanceAudit, error) {

	miary := make([]shared.AppPerformanceMetric, 0, len(miaryWydajnosciStrony))
	for _, klucz := range miaryWydajnosciStrony {
		sprawdzenie, jest := r.Sprawdzenia[klucz]
		if !jest || sprawdzenie.Wartosc == nil {
			// Miara, której program nie policzył, nie wchodzi z wartością zero.
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
// wyliczenia, ponieważ brama kontraktu sprawdza wyliczenia wyłącznie przy
// komendach wystawionych jako narzędzia modelu, a ta do nich nie należy.
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

// bladProgramuAplikacji odróżnia brak programu i przekroczenie granicy czasu
// od usterki rdzenia. Przekroczenie rozstrzyga się zmierzonym czasem, nie
// treścią komunikatu arsenału, ponieważ treść przy zmianie brzmi inaczej.
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
				" — to jest przekroczenie granicy, nie usterka serwera; naprawa: podnieść "+
				"granicę polem timeoutMs albo wskazać stronę, która wczytuje się szybciej. "+
				"Diagnostyka warstwy: "+err.Error()))
	}
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
		"moduł Apps: komenda "+komenda+" nie doszła do skutku: "+err.Error()))
}
