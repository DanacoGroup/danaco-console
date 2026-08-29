// Plik niesie sam pomiar sondy kondycji: co mierzy każdy z pięciu rodzajów — http, tcp, internal, command i modelCall; definicje, seria i dostępność leżą w adapter_kondycja.go.
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
	// celSondyWewnetrznejBaza mierzy czas obiegu magazynu stanu rdzenia, wykonując prawdziwe zapytanie do bazy.
	celSondyWewnetrznejBaza = "database"
	// celSondyWewnetrznejWykonanie mierzy stan samego procesu rdzenia, czytając jego liczniki wykonania w danej chwili.
	celSondyWewnetrznejWykonanie = "runtime"

	// progOslabieniaSondy — odpowiedź wolniejsza niż połowa limitu czasu jest
	// odpowiedzią, ale nie taką, żeby nazwać ją pełną sprawnością. Stąd stan
	// pośredni: `degraded` liczy się w dostępności jako połowa udanego.
	progOslabieniaSondy = 0.5

	// granicaSondyProgramu domyka czas jednego uruchomienia programu sondy, licząc od startu procesu do jego zakończenia.
	granicaSondyProgramu = 60 * time.Second

	// trescSondySilnika jest zapytaniem sondy `modelCall`. Krótkie z zamysłu:
	// sonda ma zmierzyć, czy kanał odpowiada, a nie wygenerować treść, za którą
	// Operator zapłaci przy każdym przebiegu.
	trescSondySilnika = "ping"
)

// zmierz wykonuje jeden pomiar wskazanej sondy i oddaje jego wynik; metoda nie zwraca błędu, bo niepowodzenie pomiaru jest wynikiem pomiaru.
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
			"serwer nie zna rodzaju sondy „"+sonda.Rodzaj+"”, więc nie ma czym jej wykonać")
	}

	czas := time.Since(poczatek).Milliseconds()
	wynik.CzasOdpowiedziMs = &czas
	oslabPrzyPowolnejOdpowiedzi(&wynik, czas, limit)
	return wynik
}

// limitCzasuSondy rozstrzyga granicę czasu jednego przebiegu sondy wskazanego rodzaju, czytając ją z konfiguracji.
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

// zmierzHttp wysyła żądanie pod adres sondy; klient jest budowany na jeden przebieg i nie chodzi za przekierowaniami dalej niż pięć razy.
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

// zmierzGniazdo otwiera połączenie z gniazdem sondy i patrzy, czy się w ogóle udało je otworzyć na czas.
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

// zmierzWnetrze mierzy sam rdzeń: database odpytuje bazę stanu prawdziwym zapytaniem, runtime czyta liczniki procesu, a cel spoza tych dwóch kończy się stanem nie wiem.
func (a *adapterKondycji) zmierzWnetrze(ctx context.Context, sonda dane.SondaKondycji,
	wynik *dane.WynikSondyKondycji) {

	switch strings.ToLower(strings.TrimSpace(sonda.Cel)) {
	case celSondyWewnetrznejBaza:
		if a.repozytorium == nil {
			ustawStanSondy(wynik, shared.HealthProbeStatusUnknown,
				"serwer nie ma wpiętego magazynu — nie ma czego odpytać")
			return
		}
		if err := a.repozytorium.Puls(ctx); err != nil {
			ustawStanSondy(wynik, shared.HealthProbeStatusDown,
				"baza stanu serwera nie odpowiedziała: "+err.Error())
			return
		}
		ustawStanSondy(wynik, shared.HealthProbeStatusUp,
			"baza stanu serwera odpowiedziała na zapytanie kontrolne")
	case celSondyWewnetrznejWykonanie:
		var pamiec runtime.MemStats
		runtime.ReadMemStats(&pamiec)
		ustawStanSondy(wynik, shared.HealthProbeStatusUp,
			"proces serwera prowadzi "+strconv.Itoa(runtime.NumGoroutine())+
				" wątków i trzyma "+strconv.FormatUint(pamiec.HeapAlloc/1024/1024, 10)+
				" MB sterty")
	default:
		ustawStanSondy(wynik, shared.HealthProbeStatusUnknown,
			"nie znam wewnętrznego celu „"+sonda.Cel+"” — serwer mierzy: "+
				celSondyWewnetrznejBaza+", "+celSondyWewnetrznejWykonanie)
	}
}

// zmierzProgram uruchamia program sondy i patrzy na jego kod wyjścia, idąc tą samą drogą co każde inne wołanie arsenału: przez port uruchamiacza i bramę izolacji.
func (a *adapterKondycji) zmierzProgram(ctx context.Context, sonda dane.SondaKondycji,
	wynik *dane.WynikSondyKondycji) {

	if a.uruchamiacz == nil {
		ustawStanSondy(wynik, shared.HealthProbeStatusUnknown,
			"serwer nie ma uruchamiacza procesów — sonda programowa nie ma czym wystartować; "+
				"naprawa: podpiąć warstwę kanału (injection) przy składaniu serwera")
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

// zmierzKanal wysyła krótkie zapytanie kanałem modelu i czeka na jego odpowiedź w granicy czasu sondy.
func (a *adapterKondycji) zmierzKanal(ctx context.Context, sonda dane.SondaKondycji,
	wynik *dane.WynikSondyKondycji) {

	if a.kanaly == nil {
		ustawStanSondy(wynik, shared.HealthProbeStatusUnknown,
			"serwer nie ma wpiętego rejestru kanałów — sonda kanału nie ma czym zapytać")
		return
	}
	kod := strings.TrimSpace(sonda.Cel)
	if _, jest := a.kanaly.Kanal(kod); !jest {
		ustawStanSondy(wynik, shared.HealthProbeStatusDown,
			"kanału „"+kod+"” nie ma w rejestrze kanałów serwera")
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
