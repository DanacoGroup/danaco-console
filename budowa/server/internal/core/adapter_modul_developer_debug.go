// Adapter pięciu komend okna Run & Debug; punkt przerwania jest trwały w bazie, sesja żyje wyłącznie w pamięci rdzenia.
package core

import (
	"context"
	"encoding/json"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

const (
	przedrostekSesjiDebugowania = "dbg-"
	przedrostekPunktuPrzerwania = "bpt-"
)

// czasPolaczeniaAdaptera: adapter, który się wywrócił, nie może zawiesić okna.
const czasPolaczeniaAdaptera = 20 * time.Second

type sesjaDebugowania struct {
	kod     string
	oknoKod string
	adapter string
	program string
	zaczeta time.Time
	klient  *klientDap
	uchwyt  session.UchwytProcesu
	drzewo  *session.DrzewoProcesu
	// Delve mówi wyłącznie po TCP; połączenie domyka się z sesją.
	polaczenie net.Conn
	nasluch    net.Listener

	mu               sync.Mutex
	stan             shared.DebugStatus
	watek            *string
	powodZatrzymania *string
}

func (s *sesjaDebugowania) Migawka() shared.DebugSession {
	s.mu.Lock()
	defer s.mu.Unlock()
	opis := shared.DebugSession{
		Id:            s.kod,
		WindowId:      s.oknoKod,
		Adapter:       s.adapter,
		Status:        s.stan,
		ThreadId:      s.watek,
		StoppedReason: s.powodZatrzymania,
		StartedAt:     s.zaczeta.UnixMilli(),
	}
	if s.program != "" {
		opis.Configuration = wskaznikTekstu(s.program)
	}
	return opis
}

type rejestrSesjiDebugowania struct {
	mu    sync.Mutex
	sesje map[string]*sesjaDebugowania
}

func nowyRejestrSesjiDebugowania() *rejestrSesjiDebugowania {
	return &rejestrSesjiDebugowania{sesje: map[string]*sesjaDebugowania{}}
}

func (r *rejestrSesjiDebugowania) Wpisz(sesja *sesjaDebugowania) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sesje[sesja.kod] = sesja
}

func (r *rejestrSesjiDebugowania) Sesja(kod string) (*sesjaDebugowania, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	sesja, jest := r.sesje[kod]
	return sesja, jest
}

func (r *rejestrSesjiDebugowania) SesjeOkna(oknoKod string) []*sesjaDebugowania {
	r.mu.Lock()
	defer r.mu.Unlock()
	sesje := make([]*sesjaDebugowania, 0, 2)
	for _, sesja := range r.sesje {
		if sesja.oknoKod == oknoKod {
			sesje = append(sesje, sesja)
		}
	}
	return sesje
}

func (r *rejestrSesjiDebugowania) Usun(kod string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.sesje, kod)
}

func (r *rejestrSesjiDebugowania) Zamknij() {
	if r == nil {
		return
	}
	r.mu.Lock()
	sesje := make([]*sesjaDebugowania, 0, len(r.sesje))
	for _, sesja := range r.sesje {
		sesje = append(sesje, sesja)
	}
	r.sesje = map[string]*sesjaDebugowania{}
	r.mu.Unlock()
	for _, sesja := range sesje {
		sesja.Zakoncz()
	}
}

// Debugger uruchamia proces debugowany jako potomka: ubicie samego adaptera zostawiłoby program przy życiu.
func (s *sesjaDebugowania) Zakoncz() {
	if s.klient != nil {
		_, _ = s.klient.Wolaj("disconnect", map[string]any{"terminateDebuggee": true})
		s.klient.Zamknij()
	}
	if s.polaczenie != nil {
		_ = s.polaczenie.Close()
	}
	if s.nasluch != nil {
		_ = s.nasluch.Close()
	}
	if s.drzewo != nil {
		_ = s.drzewo.Ubij()
		s.drzewo.Zwolnij()
	} else if s.uchwyt != nil {
		_ = s.uchwyt.Ubij()
	}
	s.mu.Lock()
	s.stan = shared.DebugStatusTerminated
	s.mu.Unlock()
}

func (a *adapterDevelopera) UruchomDebugowanie(ctx context.Context,
	z shared.DeveloperDebugSessionStartRequest) (shared.DeveloperDebugSessionStartResponse, error) {

	okno, err := a.oknoDevelopera(z.WindowId)
	if err != nil {
		return shared.DeveloperDebugSessionStartResponse{}, err
	}
	if err := sprawdzZmianeSystemu(okno.TrybUprawnien, "uruchomienie debugowania"); err != nil {
		return shared.DeveloperDebugSessionStartResponse{}, err
	}
	korzenie := a.korzenieOkna(okno)
	if len(korzenie) == 0 {
		return shared.DeveloperDebugSessionStartResponse{}, bladZadaniaDevelopera(
			"okno " + z.WindowId + " nie ma katalogu roboczego, więc nie ma czego debugować")
	}

	// Na okno przypada jedna sesja naraz.
	for _, poprzednia := range a.sesjeDebugowania.SesjeOkna(okno.Id) {
		poprzednia.Zakoncz()
		a.sesjeDebugowania.Usun(poprzednia.kod)
	}

	program := korzenie[0]
	if z.Program != nil && strings.TrimSpace(*z.Program) != "" {
		sciezka, err := sciezkaWObszarze(korzenie, *z.Program, plikIstnieje)
		if err != nil {
			return shared.DeveloperDebugSessionStartResponse{}, err
		}
		program = sciezka
	}

	nazwaAdaptera := "dlv"
	if z.Adapter != nil && strings.TrimSpace(*z.Adapter) != "" {
		nazwaAdaptera = strings.TrimSpace(*z.Adapter)
	}
	if nazwaAdaptera != "dlv" && nazwaAdaptera != "delve" {
		return shared.DeveloperDebugSessionStartResponse{}, bladZadaniaDevelopera(
			"serwer ma dziś adapter debugowania dla Go (Delve); adapter " + nazwaAdaptera +
				" nie jest zainstalowany po stronie serwera")
	}

	// Nasłuch zakłada rdzeń, adapter dzwoni do niego; port wybiera jądro.
	nasluch, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return shared.DeveloperDebugSessionStartResponse{}, bladWykonaniaDevelopera(
			"nie można założyć nasłuchu dla adaptera debugowania: " + err.Error())
	}
	polecenie, err := a.polecenieDopuszczoneDevelopera(okno, narzedzieDelve.Program,
		[]string{"dap", "--client-addr=" + nasluch.Addr().String()})
	if err != nil {
		_ = nasluch.Close()
		return shared.DeveloperDebugSessionStartResponse{}, err
	}
	if a.uruchamiacz == nil {
		_ = nasluch.Close()
		return shared.DeveloperDebugSessionStartResponse{}, bladWykonaniaDevelopera(
			"serwer nie ma uruchamiacza procesów")
	}
	uchwyt, err := a.uruchamiacz.UruchomProces(ctx, okno, polecenie)
	if err != nil {
		_ = nasluch.Close()
		return shared.DeveloperDebugSessionStartResponse{}, bladZasobuDevelopera(
			"nie można uruchomić adaptera debugowania Delve: " + err.Error() +
				"; naprawa po stronie serwera: " + narzedzieDelve.Pakiet)
	}
	drzewo, err := session.PrzejmijDrzewo(uchwyt.Pid())
	if err != nil {
		_ = nasluch.Close()
		_ = uchwyt.Ubij()
		_ = uchwyt.Czekaj()
		return shared.DeveloperDebugSessionStartResponse{}, bladWykonaniaDevelopera(
			"nie można objąć drzewa procesu debugowania: " + err.Error())
	}

	// Adapter, który nie zadzwonił, nie wystartował; odmowa niesie jego diagnostykę.
	if err := nasluch.(*net.TCPListener).SetDeadline(
		time.Now().Add(czasPolaczeniaAdaptera)); err != nil {
		_ = nasluch.Close()
	}
	polaczenie, err := nasluch.Accept()
	if err != nil {
		diagnostyka := czytajDoKonca(uchwyt.Diagnostyka())
		_ = nasluch.Close()
		_ = drzewo.Ubij()
		drzewo.Zwolnij()
		return shared.DeveloperDebugSessionStartResponse{}, bladWykonaniaDevelopera(
			"adapter debugowania Delve nie połączył się z serwerem: " +
				skrocDiagnostyke(diagnostyka, err))
	}

	sesja := &sesjaDebugowania{
		kod:        nowyIdentyfikator(przedrostekSesjiDebugowania),
		oknoKod:    okno.Id,
		adapter:    nazwaAdaptera,
		program:    program,
		zaczeta:    time.Now().UTC(),
		klient:     nowyKlientDap(polaczenie, polaczenie),
		uchwyt:     uchwyt,
		drzewo:     drzewo,
		polaczenie: polaczenie,
		nasluch:    nasluch,
		stan:       shared.DebugStatusRunning,
	}
	a.sesjeDebugowania.Wpisz(sesja)
	go a.pilnujZdarzenDebugowania(sesja)

	if err := a.rozpocznijRozmoweZAdapterem(ctx, sesja, okno, program, z); err != nil {
		sesja.Zakoncz()
		a.sesjeDebugowania.Usun(sesja.kod)
		return shared.DeveloperDebugSessionStartResponse{}, err
	}
	return shared.DeveloperDebugSessionStartResponse{Session: sesja.Migawka()}, nil
}

func (a *adapterDevelopera) rozpocznijRozmoweZAdapterem(ctx context.Context,
	sesja *sesjaDebugowania, okno session.Okno, program string,
	z shared.DeveloperDebugSessionStartRequest) error {

	if _, err := sesja.klient.Wolaj("initialize", map[string]any{
		"clientID":        "danaco-console",
		"adapterID":       sesja.adapter,
		"linesStartAt1":   true,
		"columnsStartAt1": true,
		"pathFormat":      "path",
	}); err != nil {
		return bladWykonaniaDevelopera("adapter debugowania nie przyjął powitania: " + err.Error())
	}

	uruchomienie := map[string]any{
		"request": "launch",
		"mode":    "debug",
		"program": program,
		"name":    "Danaco Console",
	}
	if len(z.Arguments) > 0 {
		uruchomienie["args"] = z.Arguments
	}
	if z.StopOnEntry != nil {
		uruchomienie["stopOnEntry"] = *z.StopOnEntry
	}
	if _, err := sesja.klient.Wolaj("launch", uruchomienie); err != nil {
		return bladWykonaniaDevelopera("adapter debugowania nie uruchomił programu " +
			program + ": " + err.Error())
	}

	if err := a.podajPunktyAdapterowi(ctx, sesja, okno.Id); err != nil {
		return err
	}
	if _, err := sesja.klient.Wolaj("configurationDone", map[string]any{}); err != nil {
		return bladWykonaniaDevelopera("adapter debugowania nie domknął konfiguracji: " + err.Error())
	}
	return nil
}

// Protokół DAP zna wyłącznie cały wykaz punktów pliku naraz, więc każde ustawienie zastępuje wykaz pliku.
func (a *adapterDevelopera) podajPunktyAdapterowi(ctx context.Context,
	sesja *sesjaDebugowania, oknoKod string) error {

	if a.repozytorium == nil {
		return nil
	}
	punkty, err := a.repozytorium.PunktyPrzerwania(ctx, oknoKod, "")
	if err != nil {
		return bladWykonaniaDevelopera("nie można odczytać punktów przerwania okna: " + err.Error())
	}
	wedlugPlikow := map[string][]dane.PunktPrzerwania{}
	for _, punkt := range punkty {
		wedlugPlikow[punkt.Sciezka] = append(wedlugPlikow[punkt.Sciezka], punkt)
	}
	for sciezka, zestaw := range wedlugPlikow {
		if _, err := sesja.klient.Wolaj("setBreakpoints",
			zadanieUstawieniaPunktow(sciezka, zestaw)); err != nil {
			// Plik nieznany adapterowi nie zatrzymuje startu sesji.
			continue
		}
	}
	return nil
}

func zadanieUstawieniaPunktow(sciezka string, punkty []dane.PunktPrzerwania) map[string]any {
	wykaz := make([]map[string]any, 0, len(punkty))
	for _, punkt := range punkty {
		pozycja := map[string]any{"line": punkt.Wiersz}
		if punkt.Warunek != nil && *punkt.Warunek != "" {
			pozycja["condition"] = *punkt.Warunek
		}
		if punkt.WarunekTrafien != nil && *punkt.WarunekTrafien != "" {
			pozycja["hitCondition"] = *punkt.WarunekTrafien
		}
		if punkt.Wpis != nil && *punkt.Wpis != "" {
			pozycja["logMessage"] = *punkt.Wpis
		}
		wykaz = append(wykaz, pozycja)
	}
	return map[string]any{
		"source":      map[string]any{"path": sciezka},
		"breakpoints": wykaz,
	}
}

// Obserwator pracuje poza żądaniem: program zatrzymuje się na punkcie, gdy do niego dojdzie.
func (a *adapterDevelopera) pilnujZdarzenDebugowania(sesja *sesjaDebugowania) {
	for zdarzenie := range sesja.klient.Zdarzenia() {
		switch zdarzenie.Event {
		case "stopped":
			var tresc struct {
				Reason   string `json:"reason"`
				ThreadId int    `json:"threadId"`
			}
			_ = json.Unmarshal(zdarzenie.Body, &tresc)
			sesja.mu.Lock()
			sesja.stan = shared.DebugStatusStopped
			sesja.powodZatrzymania = wskaznikTekstu(tresc.Reason)
			if tresc.ThreadId != 0 {
				sesja.watek = wskaznikTekstu(strconv.Itoa(tresc.ThreadId))
			}
			sesja.mu.Unlock()
		case "continued":
			sesja.mu.Lock()
			sesja.stan = shared.DebugStatusRunning
			sesja.powodZatrzymania = nil
			sesja.mu.Unlock()
		case "terminated", "exited":
			sesja.mu.Lock()
			sesja.stan = shared.DebugStatusTerminated
			sesja.mu.Unlock()
		}
	}
	// Kanał zamknięty znaczy adapter, który przestał mówić.
	sesja.mu.Lock()
	sesja.stan = shared.DebugStatusTerminated
	sesja.mu.Unlock()
}

func (a *adapterDevelopera) SterujDebugowaniem(_ context.Context,
	z shared.DeveloperDebugSessionControlRequest) (shared.DeveloperDebugSessionControlResponse, error) {

	sesja, err := a.sesjaDebugowaniaZadania(z.SessionId)
	if err != nil {
		return shared.DeveloperDebugSessionControlResponse{}, err
	}

	if z.Step == shared.DebugStepKindStop {
		sesja.Zakoncz()
		a.sesjeDebugowania.Usun(sesja.kod)
		return shared.DeveloperDebugSessionControlResponse{Session: sesja.Migawka()}, nil
	}

	komenda, jest := komendaKrokuDebugowania(z.Step)
	if !jest {
		return shared.DeveloperDebugSessionControlResponse{}, bladZadaniaDevelopera(
			"nieznany krok debugowania: " + string(z.Step))
	}
	argumenty := map[string]any{}
	if numer, jest := numerWatku(z.ThreadId, sesja); jest {
		argumenty["threadId"] = numer
	}
	if _, err := sesja.klient.Wolaj(komenda, argumenty); err != nil {
		return shared.DeveloperDebugSessionControlResponse{}, bladWykonaniaDevelopera(
			"adapter debugowania odmówił kroku " + string(z.Step) + ": " + err.Error())
	}

	if z.Step != shared.DebugStepKindPause {
		sesja.mu.Lock()
		sesja.stan = shared.DebugStatusRunning
		sesja.powodZatrzymania = nil
		sesja.mu.Unlock()
	}
	return shared.DeveloperDebugSessionControlResponse{Session: sesja.Migawka()}, nil
}

func komendaKrokuDebugowania(krok shared.DebugStepKind) (string, bool) {
	switch krok {
	case shared.DebugStepKindContinue:
		return "continue", true
	case shared.DebugStepKindPause:
		return "pause", true
	case shared.DebugStepKindStepOver:
		return "next", true
	case shared.DebugStepKindStepInto:
		return "stepIn", true
	case shared.DebugStepKindStepOut:
		return "stepOut", true
	case shared.DebugStepKindRestart:
		return "restart", true
	default:
		return "", false
	}
}

func numerWatku(zadany *string, sesja *sesjaDebugowania) (int, bool) {
	wskazanie := ""
	if zadany != nil {
		wskazanie = strings.TrimSpace(*zadany)
	}
	if wskazanie == "" {
		sesja.mu.Lock()
		if sesja.watek != nil {
			wskazanie = *sesja.watek
		}
		sesja.mu.Unlock()
	}
	numer, err := strconv.Atoi(wskazanie)
	if err != nil {
		return 0, false
	}
	return numer, true
}

// Odpowiedź niesie komplet punktów okna: margines Code Editora rysuje wszystkie naraz.
func (a *adapterDevelopera) UstawPunktPrzerwania(ctx context.Context,
	z shared.DeveloperBreakpointSetRequest) (shared.DeveloperBreakpointSetResponse, error) {

	okno, sciezka, err := a.plikOkna(z.WindowId, z.Path)
	if err != nil {
		return shared.DeveloperBreakpointSetResponse{}, err
	}
	if a.repozytorium == nil {
		return shared.DeveloperBreakpointSetResponse{}, bladZasobuDevelopera(
			"serwer nie ma miejsca na punkty przerwania")
	}
	if z.Line < 1 {
		return shared.DeveloperBreakpointSetResponse{}, bladZadaniaDevelopera(
			"punkt przerwania wymaga wiersza liczonego od jedynki")
	}

	if z.Remove != nil && *z.Remove {
		if err := a.repozytorium.UsunPunktPrzerwania(ctx, okno.Id, sciezka,
			int64(z.Line)); err != nil {
			return shared.DeveloperBreakpointSetResponse{}, bladWykonaniaDevelopera(
				"nie można zdjąć punktu przerwania: " + err.Error())
		}
	} else {
		rodzaj := shared.BreakpointKind(shared.BreakpointKindLine)
		if z.Kind != nil && *z.Kind != "" {
			rodzaj = *z.Kind
		} else if z.LogMessage != nil && *z.LogMessage != "" {
			rodzaj = shared.BreakpointKindLogpoint
		} else if z.Condition != nil && *z.Condition != "" {
			rodzaj = shared.BreakpointKindConditional
		}
		if err := a.repozytorium.ZapiszPunktPrzerwania(ctx, dane.PunktPrzerwania{
			Kod:            nowyIdentyfikator(przedrostekPunktuPrzerwania),
			OknoKod:        okno.Id,
			Sciezka:        sciezka,
			Wiersz:         int64(z.Line),
			Rodzaj:         rodzaj,
			Warunek:        z.Condition,
			WarunekTrafien: z.HitCondition,
			Wpis:           z.LogMessage,
		}); err != nil {
			return shared.DeveloperBreakpointSetResponse{}, bladWykonaniaDevelopera(
				"nie można zapisać punktu przerwania: " + err.Error())
		}
	}

	punkty, err := a.repozytorium.PunktyPrzerwania(ctx, okno.Id, "")
	if err != nil {
		return shared.DeveloperBreakpointSetResponse{}, bladWykonaniaDevelopera(
			"nie można odczytać punktów przerwania okna: " + err.Error())
	}

	a.dosleZmienionePunkty(okno.Id, sciezka, punkty)

	wykaz := make([]shared.Breakpoint, 0, len(punkty))
	for _, punkt := range punkty {
		wykaz = append(wykaz, shared.Breakpoint{
			Id:           punkt.Kod,
			Path:         punkt.Sciezka,
			Line:         int(punkt.Wiersz),
			Kind:         punkt.Rodzaj,
			Condition:    punkt.Warunek,
			HitCondition: punkt.WarunekTrafien,
			LogMessage:   punkt.Wpis,
			Verified:     punkt.Zweryfikowany,
		})
	}
	return shared.DeveloperBreakpointSetResponse{Breakpoints: wykaz}, nil
}

func (a *adapterDevelopera) dosleZmienionePunkty(oknoKod, sciezka string,
	punkty []dane.PunktPrzerwania) {

	pliku := make([]dane.PunktPrzerwania, 0, 4)
	for _, punkt := range punkty {
		if punkt.Sciezka == sciezka {
			pliku = append(pliku, punkt)
		}
	}
	for _, sesja := range a.sesjeDebugowania.SesjeOkna(oknoKod) {
		_, _ = sesja.klient.Wolaj("setBreakpoints", zadanieUstawieniaPunktow(sciezka, pliku))
	}
}

// Jedna komenda oddaje ramki, zakresy i zmienne naraz: okno rysuje stos i drzewo zmiennych w tej samej chwili.
func (a *adapterDevelopera) ZakresDebugowania(_ context.Context,
	z shared.DeveloperDebugScopeGetRequest) (shared.DeveloperDebugScopeGetResponse, error) {

	sesja, err := a.sesjaDebugowaniaZadania(z.SessionId)
	if err != nil {
		return shared.DeveloperDebugScopeGetResponse{}, err
	}

	if z.VariablesRef != nil && strings.TrimSpace(*z.VariablesRef) != "" {
		zmienne, err := zmienneZOdwolania(sesja, *z.VariablesRef)
		if err != nil {
			return shared.DeveloperDebugScopeGetResponse{}, err
		}
		return shared.DeveloperDebugScopeGetResponse{
			Frames:    []shared.DebugFrame{},
			Scopes:    []shared.DebugScope{},
			Variables: zmienne,
		}, nil
	}

	watek, jest := numerWatku(nil, sesja)
	if !jest {
		return shared.DeveloperDebugScopeGetResponse{}, bladZadaniaDevelopera(
			"sesja " + sesja.kod + " nie stoi na żadnym wątku — stos wywołań istnieje " +
				"wtedy, gdy program jest zatrzymany")
	}
	tresc, err := sesja.klient.Wolaj("stackTrace", map[string]any{"threadId": watek})
	if err != nil {
		return shared.DeveloperDebugScopeGetResponse{}, bladWykonaniaDevelopera(
			"adapter debugowania nie oddał stosu wywołań: " + err.Error())
	}
	var stos struct {
		StackFrames []struct {
			Id     int    `json:"id"`
			Name   string `json:"name"`
			Line   int    `json:"line"`
			Column int    `json:"column"`
			Source *struct {
				Path string `json:"path"`
			} `json:"source"`
		} `json:"stackFrames"`
	}
	_ = json.Unmarshal(tresc, &stos)

	ramki := make([]shared.DebugFrame, 0, len(stos.StackFrames))
	for _, ramka := range stos.StackFrames {
		opis := shared.DebugFrame{
			Id:       strconv.Itoa(ramka.Id),
			Name:     ramka.Name,
			ThreadId: wskaznikTekstu(strconv.Itoa(watek)),
		}
		if ramka.Line > 0 {
			opis.Line = wskaznikLiczby(ramka.Line)
		}
		if ramka.Column > 0 {
			opis.Column = wskaznikLiczby(ramka.Column)
		}
		if ramka.Source != nil && ramka.Source.Path != "" {
			opis.Path = wskaznikTekstu(ramka.Source.Path)
		}
		ramki = append(ramki, opis)
	}

	ramkaWybrana := ""
	if z.FrameId != nil && strings.TrimSpace(*z.FrameId) != "" {
		ramkaWybrana = strings.TrimSpace(*z.FrameId)
	} else if len(ramki) > 0 {
		ramkaWybrana = ramki[0].Id
	}

	zakresy, zmienne := zakresyRamki(sesja, ramkaWybrana)
	return shared.DeveloperDebugScopeGetResponse{
		Frames:    ramki,
		Scopes:    zakresy,
		Variables: zmienne,
	}, nil
}

func zakresyRamki(sesja *sesjaDebugowania, ramka string) ([]shared.DebugScope, []shared.DebugVariable) {
	zakresy := make([]shared.DebugScope, 0, 4)
	zmienne := make([]shared.DebugVariable, 0, 16)
	if ramka == "" {
		return zakresy, zmienne
	}
	numer, err := strconv.Atoi(ramka)
	if err != nil {
		return zakresy, zmienne
	}
	tresc, err := sesja.klient.Wolaj("scopes", map[string]any{"frameId": numer})
	if err != nil {
		return zakresy, zmienne
	}
	var odpowiedz struct {
		Scopes []struct {
			Name               string `json:"name"`
			VariablesReference int    `json:"variablesReference"`
			Expensive          bool   `json:"expensive"`
		} `json:"scopes"`
	}
	_ = json.Unmarshal(tresc, &odpowiedz)

	for _, zakres := range odpowiedz.Scopes {
		odwolanie := strconv.Itoa(zakres.VariablesReference)
		zakresy = append(zakresy, shared.DebugScope{
			Name:         zakres.Name,
			VariablesRef: odwolanie,
			FrameId:      wskaznikTekstu(ramka),
		})
		// Zakres kosztowny zostaje bez rozwinięcia: wstrzymywałby odpowiedź na każdym kroku.
		if zakres.Expensive {
			continue
		}
		if wykaz, err := zmienneZOdwolania(sesja, odwolanie); err == nil {
			zmienne = append(zmienne, wykaz...)
		}
	}
	return zakresy, zmienne
}

func zmienneZOdwolania(sesja *sesjaDebugowania, odwolanie string) ([]shared.DebugVariable, error) {
	numer, err := strconv.Atoi(strings.TrimSpace(odwolanie))
	if err != nil {
		return nil, bladZadaniaDevelopera("odwołanie do zmiennych ma być liczbą adaptera")
	}
	tresc, err := sesja.klient.Wolaj("variables", map[string]any{"variablesReference": numer})
	if err != nil {
		return nil, bladWykonaniaDevelopera("adapter debugowania nie oddał zmiennych: " + err.Error())
	}
	var odpowiedz struct {
		Variables []struct {
			Name               string `json:"name"`
			Value              string `json:"value"`
			Type               string `json:"type"`
			VariablesReference int    `json:"variablesReference"`
		} `json:"variables"`
	}
	_ = json.Unmarshal(tresc, &odpowiedz)

	zmienne := make([]shared.DebugVariable, 0, len(odpowiedz.Variables))
	for _, zmienna := range odpowiedz.Variables {
		opis := shared.DebugVariable{Name: zmienna.Name, Value: zmienna.Value}
		if zmienna.Type != "" {
			opis.Type = wskaznikTekstu(zmienna.Type)
		}
		if zmienna.VariablesReference > 0 {
			opis.VariablesRef = wskaznikTekstu(strconv.Itoa(zmienna.VariablesReference))
			// Zmienna z odwołaniem ma zawartość do rozwinięcia; edycja w miejscu jej nie dotyczy.
		} else {
			opis.Editable = wskaznikPrawdy(true)
		}
		zmienne = append(zmienne, opis)
	}
	return zmienne, nil
}

// assignTo zamienia pytanie w przypisanie tą samą drogą protokołu.
func (a *adapterDevelopera) ObliczWyrazenie(_ context.Context,
	z shared.DeveloperDebugEvaluateRequest) (shared.DeveloperDebugEvaluateResponse, error) {

	sesja, err := a.sesjaDebugowaniaZadania(z.SessionId)
	if err != nil {
		return shared.DeveloperDebugEvaluateResponse{}, err
	}
	wyrazenie := strings.TrimSpace(z.Expression)
	if wyrazenie == "" {
		return shared.DeveloperDebugEvaluateResponse{}, bladZadaniaDevelopera(
			"obliczenie wymaga wyrażenia")
	}
	ramka, err := strconv.Atoi(strings.TrimSpace(z.FrameId))
	if err != nil {
		return shared.DeveloperDebugEvaluateResponse{}, bladZadaniaDevelopera(
			"obliczenie wymaga ramki stosu — wyrażenie liczy się w jej kontekście")
	}

	if z.AssignTo != nil && strings.TrimSpace(*z.AssignTo) != "" {
		wyrazenie = strings.TrimSpace(*z.AssignTo) + " = " + wyrazenie
	}
	tresc, err := sesja.klient.Wolaj("evaluate", map[string]any{
		"expression": wyrazenie,
		"frameId":    ramka,
		"context":    "repl",
	})
	if err != nil {
		return shared.DeveloperDebugEvaluateResponse{}, bladWykonaniaDevelopera(
			"adapter debugowania nie obliczył wyrażenia: " + err.Error())
	}
	var odpowiedz struct {
		Result             string `json:"result"`
		Type               string `json:"type"`
		VariablesReference int    `json:"variablesReference"`
	}
	_ = json.Unmarshal(tresc, &odpowiedz)

	wynik := shared.DeveloperDebugEvaluateResponse{Value: odpowiedz.Result}
	if odpowiedz.Type != "" {
		wynik.Type = wskaznikTekstu(odpowiedz.Type)
	}
	if odpowiedz.VariablesReference > 0 {
		wynik.VariablesRef = wskaznikTekstu(strconv.Itoa(odpowiedz.VariablesReference))
	}
	return wynik, nil
}

func (a *adapterDevelopera) sesjaDebugowaniaZadania(kod string) (*sesjaDebugowania, error) {
	wskazanie := strings.TrimSpace(kod)
	if wskazanie == "" {
		return nil, bladZadaniaDevelopera("czynność debugowania wymaga wskazania sesji")
	}
	sesja, jest := a.sesjeDebugowania.Sesja(wskazanie)
	if !jest {
		return nil, bladZasobuDevelopera("nie ma czynnej sesji debugowania " + wskazanie)
	}
	return sesja, nil
}
