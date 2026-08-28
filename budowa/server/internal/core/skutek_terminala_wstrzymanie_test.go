//go:build !windows

package core

import (
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"danacoconsole/shared"
)

// Sprawdzian skutku wstrzymania procesu czyta stan procesu z /proc/<pid>/stat na systemach uniksowych.

// TestWstrzymanieProcesuZatrzymujeGoWSystemie uruchamia proces długi, wstrzymuje
// go i sprawdza u systemu, że naprawdę stanął; potem wznawia i sprawdza, że
// naprawdę wrócił do biegu.
func TestWstrzymanieProcesuZatrzymujeGoWSystemie(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)
	roboczy := t.TempDir()
	oknoKod := oknoTerminalaSprawdzianu(t, zmontowany, zycie, roboczy)
	kartaKod := kartaSprawdzianu(t, zmontowany, zycie, oknoKod, roboczy)

	var uruchomiony shared.TerminalCommandExecResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTerminalCommandExec,
		shared.TerminalCommandExecRequest{
			SessionId: kartaKod,
			// Pętla zajmuje procesor, więc proces biegnący jest w stanie R albo S, nie granicznym.
			Command: "while true; do :; done",
		}, &uruchomiony)

	if uruchomiony.Process.Pid == nil {
		t.Fatal("uruchomiony proces nie ma identyfikatora systemowego, więc nie ma czego zmierzyć")
	}
	pid := *uruchomiony.Process.Pid
	t.Cleanup(func() {
		prawda := true
		_ = wykonajKomende(t, zmontowany, zycie, shared.CommandTerminalProcessKill,
			shared.TerminalProcessKillRequest{ProcessId: uruchomiony.Process.Id, Force: &prawda})
	})
	if !doczekajStanuProcesu(pid, "RS", 5*time.Second) {
		t.Fatalf("proces %d nie ruszył — sprawdzian nie ma czego wstrzymywać", pid)
	}

	var wstrzymany shared.TerminalProcessSuspendResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTerminalProcessSuspend,
		shared.TerminalProcessSuspendRequest{ProcessId: uruchomiony.Process.Id}, &wstrzymany)
	if !wstrzymany.Supported {
		t.Fatal("rdzeń zgłosił brak wsparcia wstrzymania na systemie uniksowym")
	}
	if !doczekajStanuProcesu(pid, "T", 5*time.Second) {
		t.Fatalf("proces %d nie stoi w stanie zatrzymanym, choć rdzeń zameldował wstrzymanie "+
			"(stan systemu: %q)", pid, stanProcesu(pid))
	}

	prawda := true
	var wznowiony shared.TerminalProcessSuspendResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTerminalProcessSuspend,
		shared.TerminalProcessSuspendRequest{ProcessId: uruchomiony.Process.Id, Resume: &prawda},
		&wznowiony)
	if !doczekajStanuProcesu(pid, "RS", 5*time.Second) {
		t.Errorf("proces %d nie wrócił do biegu po wznowieniu (stan systemu: %q)",
			pid, stanProcesu(pid))
	}

	// Proces zakończony nie ma czego wstrzymywać — odmowa ma to nazwać, nie meldować powodzenie.
	wykonajUdana(t, zmontowany, zycie, shared.CommandTerminalProcessKill,
		shared.TerminalProcessKillRequest{ProcessId: uruchomiony.Process.Id, Force: &prawda},
		&shared.TerminalProcessKillResponse{})
	odmowa := wykonajOdmowna(t, zmontowany, zycie, shared.CommandTerminalProcessSuspend,
		shared.TerminalProcessSuspendRequest{ProcessId: uruchomiony.Process.Id})
	if odmowa.Code != shared.ErrorCodeConflict {
		t.Errorf("wstrzymanie procesu zakończonego odmówiło kodem %s zamiast conflict", odmowa.Code)
	}
}

// stanProcesu czyta jednoliterowy stan procesu z `/proc/<pid>/stat`.
//
// Nazwa programu w tym pliku stoi w nawiasach i może zawierać spacje, więc pole
// stanu wyławia się PO ostatnim nawiasie zamykającym, a nie przez podział całego
// wiersza po spacjach.
func stanProcesu(pid int) string {
	tresc, err := os.ReadFile("/proc/" + strconv.Itoa(pid) + "/stat")
	if err != nil {
		return ""
	}
	nawias := strings.LastIndex(string(tresc), ")")
	if nawias < 0 || nawias+2 >= len(tresc) {
		return ""
	}
	pola := strings.Fields(string(tresc)[nawias+1:])
	if len(pola) == 0 {
		return ""
	}
	return pola[0]
}

// doczekajStanuProcesu czeka, aż proces wejdzie w jeden ze stanów wskazanych zbiorem liter, albo upłynie limit czasu oczekiwania.
func doczekajStanuProcesu(pid int, stany string, najdluzej time.Duration) bool {
	koniec := time.Now().Add(najdluzej)
	for time.Now().Before(koniec) {
		if stan := stanProcesu(pid); stan != "" && strings.Contains(stany, stan) {
			return true
		}
		time.Sleep(50 * time.Millisecond)
	}
	return false
}
