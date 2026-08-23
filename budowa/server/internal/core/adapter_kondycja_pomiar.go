// Odpowiedzialność pliku: sam pomiar — co właściwie mierzy każdy z pięciu
// rodzajów sondy kondycji. Definicje, seria i dostępność leżą
// w `adapter_kondycja.go`.
//
// ── Każdy rodzaj mierzy coś naprawdę ────────────────────────────────────────
//   - `http`      — wysyła żądanie pod adres i patrzy na kod odpowiedzi.
//   - `tcp`       — otwiera połączenie z gniazdem i patrzy, czy się otworzyło.
//   - `internal`  — dotyka wnętrza rdzenia: bazy stanu albo jego własnej
//     pamięci. To NIE jest „zwróć w porządku": baza odpytana jest bazą, która
//     odpowiedziała, a jej czas obiegu jest zmierzoną liczbą.
//   - `command`   — uruchamia program i patrzy na jego kod wyjścia.
//   - `modelCall` — wysyła krótkie zapytanie kanałem modelu i czeka na
//     odpowiedź.
//
// ── Czego tu nie ma ─────────────────────────────────────────────────────────
// Gałęzi „nie umiem zmierzyć, więc `up`". Każdy powód, dla którego pomiar się
// nie odbył — brak uruchamiacza, brak rejestru kanałów, nieznany cel sondy
// wewnętrznej — kończy się stanem `unknown` wraz ze zdaniem mówiącym, czego
// brakuje. `unknown` znaczy „nie wiem" i tylko tak wygląda w wykazie; `up`
// znaczyłoby „sprawdziłem i jest dobrze", a nikt nie sprawdzał.
package core

import (
	"context"
	"net"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/models"
	"danacoconsole/server/internal/session"
	"danacoconsole/server/internal/zewnetrzne"
	"danacoconsole/shared"
)

const (
	// celSondyWewnetrznejBaza mierzy czas obiegu magazynu stanu rdzenia.
	celSondyWewnetrznejBaza = "database"
	// celSondyWewnetrznejWykonanie mierzy stan samego procesu rdzenia.
	celSondyWewnetrznejWykonanie = "runtime"

	// progOslabieniaSondy — odpowiedź wolniejsza niż połowa limitu czasu jest
	// odpowiedzią, ale nie taką, żeby nazwać ją pełną sprawnością. Stąd stan
	// pośredni: `degraded` liczy się w dostępności jako połowa udanego.
	progOslabieniaSondy = 0.5

	// granicaSondyProgramu domyka czas jednego uruchomienia programu sondy.
	granicaSondyProgramu = 60 * time.Second

	// trescSondySilnika jest zapytaniem sondy `modelCall`. Krótkie z zamysłu:
	// sonda ma zmierzyć, czy kanał odpowiada, a nie wygenerować treść, za którą
	// Operator zapłaci przy każdym przebiegu.
	trescSondySilnika = "ping"
)

// zmierz wykonuje jeden pomiar wskazanej sondy i oddaje jego wynik.
//
// Wynik zawsze niesie chwilę pomiaru i zawsze niesie stan. Metoda nie zwraca
// błędu, bo niepowodzenie pomiaru JEST wynikiem pomiaru — sonda, która
// zawiodła, ma zostawić po sobie wiersz, a nie odmowę.
func (a *adapterKondycji) zmierz(ctx context.Context, sonda dane.SondaKondycji) dane.WynikSondyKondycji {
	limit := limitCzasuSondy(sonda)
	pomiar, odwolaj := context.WithTimeout(ctx, limit)
	defer odwolaj()

	poczatek := time.Now()
	wynik := dane.WynikSondyKondycji{Wykonano: poczatek.UnixMilli()}

	switch shared.HealthProbeKind(sonda.Rodzaj) {
	case shared.HealthProbeKindHttp:
		a.zmierzHttp(pomiar, sonda, limit, &wynik)
	case shared.HealthProbeKindTcp:
		zmierzGniazdo(sonda, limit, &wynik)
	case shared.HealthProbeKindInternal:
		a.zmierzWnetrze(pomiar, sonda, &wynik)
	case shared.HealthProbeKindCommand:
		a.zmierzProgram(pomiar, sonda, &wynik)
	case shared.HealthProbeKindModelCall:
		a.zmierzKanal(pomiar, sonda, &wynik)
	default:
		ustawStanSondy(&wynik, shared.HealthProbeStatusUnknown,
			"rdzeń nie zna rodzaju sondy „"+sonda.Rodzaj+"”, więc nie ma czym jej wykonać")
	}

	czas := time.Since(poczatek).Milliseconds()
	wynik.CzasOdpowiedziMs = &czas
	oslabPrzyPowolnejOdpowiedzi(&wynik, czas, limit)
	return wynik
}

// limitCzasuSondy rozstrzyga granicę czasu jednego przebiegu.
func limitCzasuSondy(sonda dane.SondaKondycji) time.Duration {
	if sonda.LimitCzasuMs == nil || *sonda.LimitCzasuMs <= 0 {
		return domyslnyLimitCzasuSondy
	}
	limit := time.Duration(*sonda.LimitCzasuMs) * time.Millisecond
	if limit > gornyLimitCzasuSondy {
		return gornyLimitCzasuSondy
	}
	return limit
}

// ustawStanSondy zapisuje stan wraz z jego uzasadnieniem. Uzasadnienie idzie
// zawsze, także przy powodzeniu: wiersz serii bez zdania o tym, co zmierzono,
// jest liczbą bez świadka.
func ustawStanSondy(wynik *dane.WynikSondyKondycji, stan shared.HealthProbeStatus, szczegol string) {
	wynik.Stan = string(stan)
	if strings.TrimSpace(szczegol) != "" {
		kopia := szczegol
		wynik.Szczegol = &kopia
	}
}

// oslabPrzyPowolnejOdpowiedzi obniża stan udany do pośredniego, gdy odpowiedź
// zajęła więcej niż połowę limitu. Stan nieudany zostaje nieudany — powolna
// awaria jest awarią.
func oslabPrzyPowolnejOdpowiedzi(wynik *dane.WynikSondyKondycji, czas int64, limit time.Duration) {
	if wynik.Stan != shared.HealthProbeStatusUp {
		return
	}
	if float64(czas) <= float64(limit.Milliseconds())*progOslabieniaSondy {
		return
	}
	wynik.Stan = shared.HealthProbeStatusDegraded
	opis := "odpowiedź przyszła po " + strconv.FormatInt(czas, 10) +
		" ms, czyli po ponad połowie limitu czasu tej sondy"
	if wynik.Szczegol != nil {
		opis = *wynik.Szczegol + "; " + opis
	}
	wynik.Szczegol = &opis
}

// zmierzHttp wysyła żądanie pod adres sondy.
//
// Klient jest budowany na jeden przebieg i nie chodzi za przekierowaniami dalej
// niż pięć razy: sonda ma zmierzyć adres, który podał Operator, a nie zwiedzić
// łańcuch przekierowań do cudzej strony błędu.
func (a *adapterKondycji) zmierzHttp(ctx context.Context, sonda dane.SondaKondycji,
	limit time.Duration, wynik *dane.WynikSondyKondycji) {

	adres := strings.TrimSpace(sonda.Cel)
	if !strings.HasPrefix(adres, "http://") && !strings.HasPrefix(adres, "https://") {
		ustawStanSondy(wynik, shared.HealthProbeStatusUnknown,
			"cel sondy http musi zaczynać się od http:// albo https://, a zaczyna się od czegoś innego")
		return
	}

	metoda := http.MethodGet
	var tresc strings.Reader
	if sonda.TrescWysylana != nil && *sonda.TrescWysylana != "" {
		metoda = http.MethodPost
		tresc = *strings.NewReader(*sonda.TrescWysylana)
	}
	zadanie, err := http.NewRequestWithContext(ctx, metoda, adres, &tresc)
	if err != nil {
		ustawStanSondy(wynik, shared.HealthProbeStatusUnknown,
			"nie można złożyć żądania pod adres sondy: "+err.Error())
		return
	}

	klient := &http.Client{
		Timeout: limit,
		CheckRedirect: func(_ *http.Request, poprzednie []*http.Request) error {
			if len(poprzednie) >= 5 {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}
	odpowiedz, err := klient.Do(zadanie)
	if err != nil {
		ustawStanSondy(wynik, shared.HealthProbeStatusDown,
			"adres nie odpowiedział: "+err.Error())
		return
	}
	defer odpowiedz.Body.Close()

	status := int64(odpowiedz.StatusCode)
	wynik.StatusHttp = &status

	oczekiwany := int64(200)
	if sonda.OczekiwanyKod != nil && *sonda.OczekiwanyKod > 0 {
		oczekiwany = *sonda.OczekiwanyKod
	}
	if status == oczekiwany {
		ustawStanSondy(wynik, shared.HealthProbeStatusUp,
			"adres odpowiedział kodem "+strconv.FormatInt(status, 10))
		return
	}
	ustawStanSondy(wynik, shared.HealthProbeStatusDown,
		"adres odpowiedział kodem "+strconv.FormatInt(status, 10)+
			", a sonda oczekiwała "+strconv.FormatInt(oczekiwany, 10))
}

// zmierzGniazdo otwiera połączenie z gniazdem sondy.
func zmierzGniazdo(sonda dane.SondaKondycji, limit time.Duration, wynik *dane.WynikSondyKondycji) {
	adres := strings.TrimSpace(sonda.Cel)
	if !strings.Contains(adres, ":") {
		ustawStanSondy(wynik, shared.HealthProbeStatusUnknown,
			"cel sondy tcp ma mieć postać host:port, a nie ma dwukropka")
		return
	}
	polaczenie, err := net.DialTimeout("tcp", adres, limit)
	if err != nil {
		ustawStanSondy(wynik, shared.HealthProbeStatusDown,
			"gniazdo nie przyjęło połączenia: "+err.Error())
		return
	}
	_ = polaczenie.Close()
	ustawStanSondy(wynik, shared.HealthProbeStatusUp, "gniazdo przyjęło połączenie")
}

// zmierzWnetrze mierzy sam rdzeń.
//
// Dwa cele, oba mierzalne: `database` odpytuje bazę stanu (i to jest prawdziwe
// zapytanie, nie sprawdzenie wskaźnika), `runtime` czyta liczniki procesu.
// Cel spoza tych dwóch kończy się stanem „nie wiem" wraz z wykazem znanych —
// zgadywanie, o co Operatorowi chodziło, dałoby pomiar czegoś innego niż prosił.
func (a *adapterKondycji) zmierzWnetrze(ctx context.Context, sonda dane.SondaKondycji,
	wynik *dane.WynikSondyKondycji) {

	switch strings.ToLower(strings.TrimSpace(sonda.Cel)) {
	case celSondyWewnetrznejBaza:
		if a.repozytorium == nil {
			ustawStanSondy(wynik, shared.HealthProbeStatusUnknown,
				"rdzeń nie ma wpiętego magazynu — nie ma czego odpytać")
			return
		}
		if err := a.repozytorium.Puls(ctx); err != nil {
			ustawStanSondy(wynik, shared.HealthProbeStatusDown,
				"baza stanu rdzenia nie odpowiedziała: "+err.Error())
			return
		}
		ustawStanSondy(wynik, shared.HealthProbeStatusUp,
			"baza stanu rdzenia odpowiedziała na zapytanie kontrolne")
	case celSondyWewnetrznejWykonanie:
		var pamiec runtime.MemStats
		runtime.ReadMemStats(&pamiec)
		ustawStanSondy(wynik, shared.HealthProbeStatusUp,
			"proces rdzenia prowadzi "+strconv.Itoa(runtime.NumGoroutine())+
				" wątków i trzyma "+strconv.FormatUint(pamiec.HeapAlloc/1024/1024, 10)+
				" MB sterty")
	default:
		ustawStanSondy(wynik, shared.HealthProbeStatusUnknown,
			"nie znam wewnętrznego celu „"+sonda.Cel+"” — rdzeń mierzy: "+
				celSondyWewnetrznejBaza+", "+celSondyWewnetrznejWykonanie)
	}
}

// zmierzProgram uruchamia program sondy i patrzy na jego kod wyjścia.
//
// Program idzie tą samą drogą co każde inne wołanie arsenału
// (`zewnetrzne.Wolaj`): przez port uruchamiacza, bramę izolacji i objęcie
// drzewa procesów. Własnego `exec.Command` tu nie ma — proces uruchomiony obok
// tej drogi wypada spod nadzoru i zostaje po nim uchwyt.
func (a *adapterKondycji) zmierzProgram(ctx context.Context, sonda dane.SondaKondycji,
	wynik *dane.WynikSondyKondycji) {

	if a.uruchamiacz == nil {
		ustawStanSondy(wynik, shared.HealthProbeStatusUnknown,
			"rdzeń nie ma uruchamiacza procesów — sonda programowa nie ma czym wystartować; "+
				"naprawa: podpiąć warstwę kanału (injection) przy składaniu rdzenia")
		return
	}
	czesci := strings.Fields(strings.TrimSpace(sonda.Cel))
	if len(czesci) == 0 {
		ustawStanSondy(wynik, shared.HealthProbeStatusUnknown,
			"cel sondy programowej jest pusty — nie ma czego uruchomić")
		return
	}

	narzedzie := zewnetrzne.Narzedzie{
		Nazwa: czesci[0], Program: czesci[0], Pakiet: czesci[0],
	}
	okno, zasady, obszar := a.zasiegKondycji()
	_, err := zewnetrzne.Wolaj(ctx, a.uruchamiacz, okno, zasady, obszar,
		narzedzie, czesci[1:], strings.TrimSpace(obszar.KatalogRoboczy), granicaSondyProgramu)
	if err != nil {
		ustawStanSondy(wynik, shared.HealthProbeStatusDown,
			"program sondy zakończył się niepowodzeniem: "+err.Error())
		return
	}
	ustawStanSondy(wynik, shared.HealthProbeStatusUp,
		"program „"+czesci[0]+"” zakończył się powodzeniem")
}

// zmierzKanal wysyła krótkie zapytanie kanałem modelu.
func (a *adapterKondycji) zmierzKanal(ctx context.Context, sonda dane.SondaKondycji,
	wynik *dane.WynikSondyKondycji) {

	if a.kanaly == nil {
		ustawStanSondy(wynik, shared.HealthProbeStatusUnknown,
			"rdzeń nie ma wpiętego rejestru kanałów — sonda kanału nie ma czym zapytać")
		return
	}
	kod := strings.TrimSpace(sonda.Cel)
	if _, jest := a.kanaly.Kanal(kod); !jest {
		ustawStanSondy(wynik, shared.HealthProbeStatusDown,
			"kanału „"+kod+"” nie ma w rejestrze kanałów rdzenia")
		return
	}
	tresc := trescSondySilnika
	if sonda.TrescWysylana != nil && strings.TrimSpace(*sonda.TrescWysylana) != "" {
		tresc = *sonda.TrescWysylana
	}
	odpowiedzial := false
	ujscie := models.UjscieFunkcji(func(_ context.Context, f models.Fragment) error {
		if f.Kind == shared.ChunkKindText {
			odpowiedzial = true
		}
		return nil
	})
	zapytanie := models.Zapytanie{
		Wiadomosc: nowyIdentyfikator(przedrostekWynikuKondycji),
		Tresc:     tresc,
		Kanal:     kod,
	}
	if err := a.kanaly.Wyslij(ctx, zapytanie, ujscie); err != nil {
		ustawStanSondy(wynik, shared.HealthProbeStatusDown,
			"kanał „"+kod+"” nie odpowiedział: "+err.Error())
		return
	}
	if !odpowiedzial {
		ustawStanSondy(wynik, shared.HealthProbeStatusDegraded,
			"kanał „"+kod+"” zakończył wywołanie bez błędu, ale nie oddał ani jednego fragmentu treści")
		return
	}
	ustawStanSondy(wynik, shared.HealthProbeStatusUp, "kanał „"+kod+"” odpowiedział treścią")
}

// zasiegKondycji składa trójkę okno–zasady–obszar dla wołania programu sondy.
// Ta sama droga i ten sam powód, co przy narzędziach obrazu: żądanie niesie
// sondę, a nie okno rozmowy, więc adresem jest najszerszy poziom zasięgu.
func (a *adapterKondycji) zasiegKondycji() (session.Okno, session.Zasady, session.Obszar) {
	okno := session.Okno{Ustawienia: session.Ustawienia{
		SrodowiskoWykonania: shared.ExecutionEnvCore,
	}}
	zasady := session.Zasady{}
	if a.rozstrzygacz != nil {
		zasady = ZasadyIzolacji(a.rozstrzygacz, konfig.Kontekst{})
	}
	obszar := session.Obszar{}
	if a.katalog != nil {
		obszar = ObszarOkna(a.katalog.Ustal(konfig.Kontekst{}, ""), "")
	}
	return okno, zasady, obszar
}
