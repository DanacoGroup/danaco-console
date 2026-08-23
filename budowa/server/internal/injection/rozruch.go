package injection

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

// Rozruch jest jedynym w drzewie miejscem, w którym powstaje i startuje proces
// modelu. Sięgają tu obie drogi kanału głównego:
//
//   - Uruchom — jedno wywołanie tury w trybie stream-json (proces.go);
//   - UruchamiaczOkien — proces okna komunikacji prowadzony przez sesję
//     (uruchamiacz_okna.go).
//
// Pakiet session nie buduje własnego `exec.Cmd`; bierze stąd gotowy uchwyt.

// Rozruch opisuje jedno uruchomienie programu zewnętrznego.
type Rozruch struct {
	// Program — ścieżka albo nazwa pliku wykonywalnego.
	Program string
	// Argumenty wiersza poleceń, bez nazwy programu.
	Argumenty []string
	// Katalog uruchomienia; pusty oznacza katalog bieżący rdzenia.
	Katalog string
	// Srodowisko w postaci KLUCZ=wartość; puste oznacza środowisko odziedziczone.
	Srodowisko []string
	// Atrybuty systemowe uruchomienia. Proces, który ma zostać objęty drzewem
	// przez warstwę sesji, dostaje tu session.AtrybutyDrzewa().
	Atrybuty *syscall.SysProcAttr
	// ZapasNaZamkniecie daje procesowi chwilę na domknięcie potoków po odwołaniu
	// kontekstu; zero oznacza brak zapasu.
	ZapasNaZamkniecie time.Duration
	// WyjscieBledowOsobno kieruje wyjście diagnostyczne do osobnego potoku
	// zamiast do bufora. Potrzebuje tego strumień okna, który czyta oba wyjścia.
	WyjscieBledowOsobno bool
}

// Start jest uchwytem uruchomionego procesu. Wypełnia interfejs
// session.UchwytProcesu, dzięki czemu warstwa sesji obejmuje ten sam proces
// drzewem potomstwa, nie startując własnego.
type Start struct {
	polecenie   *exec.Cmd
	wejscie     io.WriteCloser
	wyjscie     io.ReadCloser
	diagnostyka io.ReadCloser
	bledy       *buforBledow
}

// Wystartuj buduje i uruchamia proces według opisu rozruchu.
func Wystartuj(kontekst context.Context, r Rozruch) (*Start, error) {
	if strings.TrimSpace(r.Program) == "" {
		return nil, fmt.Errorf("injection: brak ścieżki programu kanału")
	}
	polecenie := exec.CommandContext(kontekst, r.Program, r.Argumenty...)
	polecenie.Dir = r.Katalog
	polecenie.Env = r.Srodowisko
	polecenie.SysProcAttr = r.Atrybuty
	polecenie.WaitDelay = r.ZapasNaZamkniecie

	start := &Start{polecenie: polecenie}
	wejscie, err := polecenie.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("injection: wejście procesu: %w", err)
	}
	start.wejscie = wejscie
	wyjscie, err := polecenie.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("injection: wyjście procesu: %w", err)
	}
	start.wyjscie = wyjscie
	if err := start.podepnijBledy(r.WyjscieBledowOsobno); err != nil {
		return nil, err
	}
	if err := polecenie.Start(); err != nil {
		return nil, fmt.Errorf("injection: uruchomienie %s: %w", r.Program, err)
	}
	return start, nil
}

// podepnijBledy wybiera drogę wyjścia diagnostycznego: osobny potok albo bufor.
func (s *Start) podepnijBledy(osobno bool) error {
	if !osobno {
		s.bledy = &buforBledow{}
		s.polecenie.Stderr = s.bledy
		return nil
	}
	diagnostyka, err := s.polecenie.StderrPipe()
	if err != nil {
		return fmt.Errorf("injection: wyjście diagnostyczne procesu: %w", err)
	}
	s.diagnostyka = diagnostyka
	return nil
}

// Pid zwraca identyfikator systemowy procesu.
func (s *Start) Pid() int {
	if s.polecenie.Process == nil {
		return 0
	}
	return s.polecenie.Process.Pid
}

// Wejscie zwraca strumień wejściowy procesu.
func (s *Start) Wejscie() io.WriteCloser { return s.wejscie }

// Wyjscie zwraca strumień wyjściowy procesu.
func (s *Start) Wyjscie() io.Reader { return s.wyjscie }

// Diagnostyka zwraca wyjście diagnostyczne procesu. Gdy błędy jadą do bufora,
// strumień jest pusty, a treść odczytuje Bledy.
func (s *Start) Diagnostyka() io.Reader {
	if s.diagnostyka == nil {
		return strings.NewReader("")
	}
	return s.diagnostyka
}

// Bledy zwraca zapamiętane wyjście diagnostyczne procesu.
func (s *Start) Bledy() string {
	if s.bledy == nil {
		return ""
	}
	return s.bledy.Tresc()
}

// Czekaj czeka na zakończenie procesu. Kod wyjścia różny od zera jest błędem
// bieżącego wywołania, nie awarią sesji ani konta.
func (s *Start) Czekaj() error {
	if err := s.polecenie.Wait(); err != nil {
		return fmt.Errorf("injection: proces kanału zakończył się błędem: %w", err)
	}
	return nil
}

// Ubij kończy sam proces. Drzewo potomstwa kończy warstwa sesji — tu chodzi
// wyłącznie o to, by nieudane przejęcie nie zostawiło procesu bez opieki.
func (s *Start) Ubij() error {
	if s.polecenie.Process == nil {
		return nil
	}
	if err := s.polecenie.Process.Kill(); err != nil {
		return fmt.Errorf("injection: nie można zakończyć procesu %d: %w", s.Pid(), err)
	}
	return nil
}
