// Komendy `terminal.tunnel.open`, `terminal.tunnel.list` i
// `terminal.tunnel.close` — przekierowania portów okna Session Manager.
//
// ── Czego brakowało ─────────────────────────────────────────────────────────
// Tunel dało się dotąd założyć wyłącznie poleceniem wydanym w karcie. Biegł
// wtedy jako zwykły proces polecenia: nie było wykazu tuneli, nie było ich
// stanu, a zamknięcie sprowadzało się do odszukania właściwego wiersza w Process
// Monitorze. Tunel jest bytem długożyjącym o własnym stanie i tak go tu
// prowadzimy (migracja 249).
//
// ── Dlaczego `ssh`, a nie własne przekierowanie w Go ────────────────────────
// Zasada produktu każe brać bibliotekę zamiast programu wszędzie, gdzie to
// możliwe — i tak zrobiono z git, PDF i wyszukiwaniem. Tutaj wybór jest inny,
// bo nie chodzi o czynność biblioteczną, tylko o TRANSPORT, którego kształt
// produkt już rozstrzygnął (migracja 088 „hosty zdalne”): torem jest SSH
// uruchamiany programem `ssh`, bo niesie uwierzytelnienie, szyfrowanie
// i sprawdzenie klucza hosta (`known_hosts`), których własne przekierowanie
// musiałoby dorobić od zera. Przekierowanie napisane w Go byłoby drugim,
// słabszym torem obok istniejącego — a nie usunięciem zależności.
//
// ── Skąd wiadomo, że tunel stoi ─────────────────────────────────────────────
// Z procesu, nie z zapisu. Przełącznik `-o ExitOnForwardFailure=yes` sprawia, że
// `ssh` KOŃCZY SIĘ, gdy przekierowania nie udało się założyć, zamiast biec
// z otwartym połączeniem i zamkniętym portem. Rdzeń czeka chwilę na taki koniec
// przed odpowiedzią, a potem doglądą procesu do końca jego życia: stan `active`
// znaczy „proces biegnie z założonym przekierowaniem”, a nie „polecenie
// wysłano”.
package core

import (
	"context"
	"errors"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/session"
	"danacoconsole/server/internal/zewnetrzne"
	"danacoconsole/shared"
)

// przedrostekTunelu znakuje identyfikator przekierowania portu.
const przedrostekTunelu = "ttun-"

// czasNaNieudanyTunel jest chwilą, przez którą rdzeń czeka na samoistny koniec
// procesu `ssh`. Nieudane przekierowanie kończy się w tym czasie z zapasem
// (odmowa portu jest natychmiastowa), a udane biegnie dalej i odpowiedź wraca
// bez dalszego czekania.
const czasNaNieudanyTunel = 900 * time.Millisecond

// tunelZywy jest jednym biegnącym przekierowaniem wraz z uchwytami, bez których
// nie da się go zamknąć.
type tunelZywy struct {
	kod       string
	uchwyt    session.UchwytProcesu
	drzewo    *session.DrzewoProcesu
	koniec    chan struct{}
	powod     string
	zamkniety bool
}

// rejestrTuneli trzyma tunele czynne jednego biegu rdzenia. Wiersz w bazie mówi,
// że tunel był; uchwyt tutaj jest jedynym, czym da się go zamknąć.
type rejestrTuneli struct {
	mu     sync.Mutex
	tunele map[string]*tunelZywy
}

func nowyRejestrTuneli() *rejestrTuneli {
	return &rejestrTuneli{tunele: make(map[string]*tunelZywy)}
}

func (r *rejestrTuneli) zapisz(t *tunelZywy) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tunele[t.kod] = t
}

func (r *rejestrTuneli) wez(kod string) (*tunelZywy, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	t, jest := r.tunele[kod]
	return t, jest
}

func (r *rejestrTuneli) zdejmij(kod string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.tunele, kod)
}

// zamknijWszystkie kończy tunele czynne przy zatrzymaniu rdzenia.
func (r *rejestrTuneli) zamknijWszystkie() {
	r.mu.Lock()
	wykaz := make([]*tunelZywy, 0, len(r.tunele))
	for _, tunel := range r.tunele {
		wykaz = append(wykaz, tunel)
	}
	r.mu.Unlock()
	for _, tunel := range wykaz {
		_ = tunel.drzewo.Ubij()
	}
}

// OtworzTunel obsługuje `terminal.tunnel.open`.
func (a *adapterTerminala) OtworzTunel(ctx context.Context,
	z shared.TerminalTunnelOpenRequest) (shared.TerminalTunnelOpenResponse, error) {

	dziennik, err := a.dziennikWyposazenia()
	if err != nil {
		return shared.TerminalTunnelOpenResponse{}, err
	}
	oknoKod := strings.TrimSpace(z.WindowId)
	if oknoKod == "" {
		return shared.TerminalTunnelOpenResponse{}, bladZadaniaTerminala(
			"tunel wymaga wskazania okna, bo z niego bierze się zasięg izolacji")
	}
	okno, err := a.oknoWykonania(oknoKod)
	if err != nil {
		return shared.TerminalTunnelOpenResponse{}, err
	}
	// Tunel otwiera połączenie wychodzące z maszyny rdzenia, czyli zmienia stan
	// systemu — brama trybu uprawnień okna obowiązuje go tak samo jak polecenie.
	if err := sprawdzUprawnienie(okno.TrybUprawnien, shared.ProcessInitiatorOperator); err != nil {
		return shared.TerminalTunnelOpenResponse{}, err
	}
	if a.uruchamiacz == nil {
		return shared.TerminalTunnelOpenResponse{}, bladWykonaniaTerminala(
			"rdzeń nie ma uruchamiacza procesów, więc nie założy tunelu")
	}

	cel, kodHosta, kluczSciezka, portCelu, err := a.celTunelu(ctx, z)
	if err != nil {
		return shared.TerminalTunnelOpenResponse{}, err
	}
	portLokalny, err := portLokalnyTunelu(z.LocalPort)
	if err != nil {
		return shared.TerminalTunnelOpenResponse{}, err
	}
	argumenty, err := argumentyTunelu(z, cel, kluczSciezka, portCelu, portLokalny)
	if err != nil {
		return shared.TerminalTunnelOpenResponse{}, err
	}

	wiersz := dane.TunelTerminala{
		Kod:          nowyIdentyfikator(przedrostekTunelu),
		OknoKod:      oknoKod,
		Rodzaj:       z.Kind,
		HostKod:      wskaznikTekstu(kodHosta),
		Cel:          cel,
		HostDocelowy: strings.TrimSpace(wartoscTekstu(z.RemoteHost)),
		Stan:         shared.TerminalTunnelStatusInactive,
	}
	lokalny := int64(portLokalny)
	wiersz.PortLokalny = &lokalny
	if z.RemotePort != nil && *z.RemotePort > 0 {
		zdalny := int64(*z.RemotePort)
		wiersz.PortDocelowy = &zdalny
	}
	if err := dziennik.ZapiszTunel(ctx, wiersz); err != nil {
		return shared.TerminalTunnelOpenResponse{}, err
	}

	stan, powod := a.uruchomTunel(okno, wiersz.Kod, argumenty)
	if err := dziennik.ZmienStanTunelu(ctx, wiersz.Kod, stan, powod,
		stan != shared.TerminalTunnelStatusActive); err != nil {
		return shared.TerminalTunnelOpenResponse{}, err
	}
	zapisany, err := dziennik.Tunel(ctx, wiersz.Kod)
	if err != nil {
		return shared.TerminalTunnelOpenResponse{}, err
	}
	return shared.TerminalTunnelOpenResponse{Tunnel: tunelKontraktu(zapisany)}, nil
}

// uruchomTunel startuje proces `ssh` i rozstrzyga jego stan początkowy.
//
// Rozstrzygnięcie bierze się z procesu: krótkie czekanie łapie przekierowanie,
// którego nie udało się założyć (`ExitOnForwardFailure`), i wtedy stanem jest
// `failed` wraz z powodem, który program wypisał na diagnostykę. Proces, który
// przeżył to czekanie, biegnie z założonym przekierowaniem.
func (a *adapterTerminala) uruchomTunel(okno session.Okno, kod string,
	argumenty []string) (shared.TerminalTunnelStatus, string) {

	sciezka, jest := zewnetrzne.Odnajdz(narzedzieSSH)
	if !jest {
		return shared.TerminalTunnelStatusFailed, (&zewnetrzne.BrakNarzedzia{Narzedzie: narzedzieSSH}).Error()
	}
	polecenie := session.Polecenie{
		Program:             sciezka,
		Argumenty:           argumenty,
		Katalog:             a.obszarOkna(okno).KatalogRoboczy,
		DziedziczSrodowisko: true,
	}
	dopuszczone, err := session.SprawdzPolecenie(a.zasadyOkna(okno), a.obszarOkna(okno), polecenie)
	if err != nil {
		return shared.TerminalTunnelStatusFailed, err.Error()
	}
	uchwyt, err := a.uruchamiacz.UruchomProces(okno, rozwinSrodowisko(dopuszczone))
	if err != nil {
		return shared.TerminalTunnelStatusFailed, "nie można uruchomić programu ssh: " + err.Error()
	}
	drzewo, err := session.PrzejmijDrzewo(uchwyt.Pid())
	if err != nil {
		_ = uchwyt.Ubij()
		_ = uchwyt.Czekaj()
		return shared.TerminalTunnelStatusFailed, "nie można objąć drzewa procesu tunelu: " + err.Error()
	}

	tunel := &tunelZywy{kod: kod, uchwyt: uchwyt, drzewo: drzewo, koniec: make(chan struct{})}
	a.tunele.zapisz(tunel)
	// Diagnostykę czyta osobna gorutyna, bo bufor potoku ma kilkadziesiąt
	// kilobajtów: `ssh` piszący ostrzeżenia bez odbiorcy stanąłby na zapisie.
	diagnostyka := make(chan string, 1)
	go func() { diagnostyka <- czytajDiagnostykeTunelu(uchwyt) }()
	go a.dogladajTunel(tunel, diagnostyka)

	select {
	case <-tunel.koniec:
		powod := strings.TrimSpace(tunel.powod)
		if powod == "" {
			powod = "program ssh zakończył się zaraz po starcie, nie zakładając przekierowania"
		}
		return shared.TerminalTunnelStatusFailed, powod
	case <-time.After(czasNaNieudanyTunel):
		return shared.TerminalTunnelStatusActive, ""
	}
}

// dogladajTunel czeka na koniec procesu tunelu i domyka jego wiersz.
//
// Dogląd pracuje we własnej gorutynie i własnym kontekście: tunel przeżywa
// żądanie, które go założyło, i ma być domknięty także wtedy, gdy padnie łącze
// godzinę później.
func (a *adapterTerminala) dogladajTunel(tunel *tunelZywy, diagnostyka <-chan string) {
	blad := tunel.uchwyt.Czekaj()
	tunel.powod = strings.TrimSpace(<-diagnostyka)
	if tunel.powod == "" && blad != nil {
		tunel.powod = blad.Error()
	}
	close(tunel.koniec)
	tunel.drzewo.Zwolnij()
	a.tunele.zdejmij(tunel.kod)

	if a.repozytorium == nil {
		return
	}
	stan := shared.TerminalTunnelStatus(shared.TerminalTunnelStatusFailed)
	if tunel.zamkniety {
		// Tunel zamknięty poleceniem Operatora nie jest tunelem, który zawiódł.
		stan, tunel.powod = shared.TerminalTunnelStatusInactive, ""
	}
	_ = a.repozytorium.ZmienStanTunelu(context.Background(), tunel.kod, stan, tunel.powod, true)
}

// czytajDiagnostykeTunelu zbiera diagnostykę procesu do jego końca i przycina ją
// do ostatniego wiersza — powodem niepowodzenia jest to, co `ssh` powiedział
// na końcu, a nie cała jego gadatliwość.
func czytajDiagnostykeTunelu(uchwyt session.UchwytProcesu) string {
	bufor := make([]byte, 0, 4096)
	czytnik := uchwyt.Diagnostyka()
	if czytnik == nil {
		return ""
	}
	porcja := make([]byte, 1024)
	for len(bufor) < 8*1024 {
		odczytane, err := czytnik.Read(porcja)
		if odczytane > 0 {
			bufor = append(bufor, porcja[:odczytane]...)
		}
		if err != nil {
			break
		}
	}
	wiersze := strings.Split(strings.TrimSpace(string(bufor)), "\n")
	return strings.TrimSpace(wiersze[len(wiersze)-1])
}

// WykazTuneli obsługuje `terminal.tunnel.list`.
func (a *adapterTerminala) WykazTuneli(ctx context.Context,
	z shared.TerminalTunnelListRequest) (shared.TerminalTunnelListResponse, error) {

	dziennik, err := a.dziennikWyposazenia()
	if err != nil {
		return shared.TerminalTunnelListResponse{}, err
	}
	filtr := dane.FiltrTuneliTerminala{OknoKod: strings.TrimSpace(wartoscTekstu(z.WindowId))}
	if z.Status != nil {
		filtr.Stan = *z.Status
	}
	wiersze, err := dziennik.Tunele(ctx, filtr)
	if err != nil {
		return shared.TerminalTunnelListResponse{}, err
	}
	wykaz := make([]shared.TerminalTunnel, 0, len(wiersze))
	for _, wiersz := range wiersze {
		wykaz = append(wykaz, tunelKontraktu(wiersz))
	}
	return shared.TerminalTunnelListResponse{Tunnels: wykaz, Total: len(wykaz)}, nil
}

// ZamknijTunel obsługuje `terminal.tunnel.close`.
func (a *adapterTerminala) ZamknijTunel(ctx context.Context,
	z shared.TerminalTunnelCloseRequest) (shared.TerminalTunnelCloseResponse, error) {

	dziennik, err := a.dziennikWyposazenia()
	if err != nil {
		return shared.TerminalTunnelCloseResponse{}, err
	}
	kod := strings.TrimSpace(z.TunnelId)
	if kod == "" {
		return shared.TerminalTunnelCloseResponse{}, bladZadaniaTerminala(
			"zamknięcie tunelu wymaga jego wskazania")
	}
	if _, err := dziennik.Tunel(ctx, kod); err != nil {
		if errors.Is(err, dane.ErrBrakWiersza) {
			return shared.TerminalTunnelCloseResponse{}, bladBrakuZasobuTerminala("tunel " + kod)
		}
		return shared.TerminalTunnelCloseResponse{}, err
	}
	if tunel, jest := a.tunele.wez(kod); jest {
		tunel.zamkniety = true
		if err := tunel.drzewo.Ubij(); err != nil {
			return shared.TerminalTunnelCloseResponse{}, bladWykonaniaTerminala(
				"nie można zakończyć procesu tunelu " + kod + ": " + err.Error())
		}
		// Krótkie czekanie na dogląd sprawia, że odpowiedź nie mówi „active”
		// o tunelu właśnie zamkniętym.
		select {
		case <-tunel.koniec:
		case <-time.After(czasNaDomkniecie):
		}
	} else {
		// Tunelu nie ma w rejestrze żywym: albo zakończył się sam, albo pochodzi
		// z poprzedniego biegu rdzenia. Wiersz i tak ma zejść na stan zamknięty.
		if err := dziennik.ZmienStanTunelu(ctx, kod,
			shared.TerminalTunnelStatusInactive, "", true); err != nil {
			return shared.TerminalTunnelCloseResponse{}, err
		}
	}
	zamkniety, err := dziennik.Tunel(ctx, kod)
	if err != nil {
		return shared.TerminalTunnelCloseResponse{}, err
	}
	return shared.TerminalTunnelCloseResponse{Tunnel: tunelKontraktu(zamkniety)}, nil
}

// celTunelu ustala adres celu, wpis książki, ścieżkę klucza i port celu.
func (a *adapterTerminala) celTunelu(ctx context.Context,
	z shared.TerminalTunnelOpenRequest) (string, string, string, int, error) {

	kodHosta := strings.TrimSpace(wartoscTekstu(z.HostId))
	cel := strings.TrimSpace(wartoscTekstu(z.RemoteTarget))
	if kodHosta == "" && cel == "" {
		return "", "", "", 0, bladZadaniaTerminala(
			"tunel wymaga celu — podaj wpis książki hostów polem hostId albo adres polem remoteTarget")
	}
	if kodHosta == "" {
		return cel, "", "", 0, nil
	}
	pomocnicza := &kartaTerminala{}
	if err := a.celZWpisuKsiazki(ctx, pomocnicza, kodHosta); err != nil {
		return "", "", "", 0, err
	}
	return pomocnicza.celZdalny, pomocnicza.hostKod, pomocnicza.kluczSciezka, pomocnicza.portZdalny, nil
}

// argumentyTunelu składa wiersz `ssh` właściwy rodzajowi przekierowania.
func argumentyTunelu(z shared.TerminalTunnelOpenRequest, cel, kluczSciezka string,
	portCelu, portLokalny int) ([]string, error) {

	argumenty := []string{"-N", "-T",
		// Bez pytań interaktywnych: proces rdzenia nie ma komu ich zadać, a `ssh`
		// czekający na hasło wyglądałby jak tunel, który zaraz stanie.
		"-o", "BatchMode=yes",
		// Nieudane przekierowanie ma KOŃCZYĆ proces, a nie zostawiać połączenie
		// z zamkniętym portem. Na tym stoi rozpoznanie stanu tunelu.
		"-o", "ExitOnForwardFailure=yes",
	}
	if portCelu > 0 {
		argumenty = append(argumenty, "-p", strconv.Itoa(portCelu))
	}
	if kluczSciezka != "" {
		argumenty = append(argumenty, "-i", kluczSciezka)
	}

	switch z.Kind {
	case shared.TerminalTunnelKindLocal, shared.TerminalTunnelKindRemote:
		hostDocelowy := strings.TrimSpace(wartoscTekstu(z.RemoteHost))
		if hostDocelowy == "" {
			hostDocelowy = "localhost"
		}
		if z.RemotePort == nil || *z.RemotePort <= 0 {
			return nil, bladZadaniaTerminala("przekierowanie " + string(z.Kind) +
				" wymaga portu docelowego po drugiej stronie tunelu")
		}
		przelacznik := "-L"
		if z.Kind == shared.TerminalTunnelKindRemote {
			przelacznik = "-R"
		}
		argumenty = append(argumenty, przelacznik,
			strconv.Itoa(portLokalny)+":"+hostDocelowy+":"+strconv.Itoa(*z.RemotePort))
	case shared.TerminalTunnelKindDynamic:
		// Przekierowanie dynamiczne (SOCKS) nie ma drugiej strony: cel wybiera
		// każde połączenie z osobna, więc pola maszyny i portu docelowego są
		// w kontrakcie opisane jako nieużywane i tu ich nie ma.
		argumenty = append(argumenty, "-D", strconv.Itoa(portLokalny))
	default:
		return nil, bladZadaniaTerminala("rodzaj przekierowania " + string(z.Kind) +
			" nie należy do słownika kontraktu")
	}
	return append(argumenty, cel), nil
}

// portLokalnyTunelu bierze port z żądania albo wskazuje wolny.
//
// Wolny port wybiera SYSTEM, nie licznik rdzenia: nasłuch na porcie zerowym
// oddaje port, o którym jądro wie, że jest wolny w tej chwili. Licznik po stronie
// rdzenia trafiałby prędzej czy później na port zajęty przez cudzy proces.
func portLokalnyTunelu(wskazany *int) (int, error) {
	if wskazany != nil && *wskazany > 0 {
		if *wskazany > 65535 {
			return 0, bladZadaniaTerminala("port lokalny leży poza zakresem portów (1–65535)")
		}
		return *wskazany, nil
	}
	nasluch, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, bladWykonaniaTerminala("nie można wskazać wolnego portu: " + err.Error())
	}
	defer nasluch.Close()
	return nasluch.Addr().(*net.TCPAddr).Port, nil
}

// tunelKontraktu przekłada wiersz tunelu na byt kontraktu.
//
// Liczników `bytesIn` i `bytesOut` nie wypełniamy. Kontrakt ma na nie pola
// opcjonalne, a rdzeń nie stoi w torze bajtów — przenosi je `ssh` we własnym
// procesie. Zero znaczyłoby „nic nie przeszło" zamiast „nie wiadomo".
func tunelKontraktu(w dane.TunelTerminala) shared.TerminalTunnel {
	tunel := shared.TerminalTunnel{
		Id:         w.Kod,
		WindowId:   w.OknoKod,
		Kind:       w.Rodzaj,
		HostId:     w.HostKod,
		LocalPort:  wskaznikMalej(w.PortLokalny),
		RemoteHost: wskaznikTekstu(w.HostDocelowy),
		RemotePort: wskaznikMalej(w.PortDocelowy),
		Status:     w.Stan,
		OpenedAt:   chwilaBazy(w.Zalozono),
	}
	if w.Powod != "" {
		powod := w.Powod
		tunel.ErrorMessage = &powod
	}
	if w.Zamknieto != nil {
		chwila := chwilaBazy(*w.Zamknieto)
		tunel.ClosedAt = &chwila
	}
	return tunel
}
