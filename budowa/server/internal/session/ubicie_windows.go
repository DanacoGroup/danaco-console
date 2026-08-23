//go:build windows

package session

import (
	"fmt"
	"os"
	"sync"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Drzewo procesów okna na Windows opiera się na Job Object. Proces okna wraz
// z całym potomstwem należy do jednego zadania, więc jedno wywołanie jądra
// kończy całe drzewo. Zamiast `taskkill` idzie wywołanie jądra — bez zależności
// od narzędzi zewnętrznych systemu.
//
// Uchwyt zadania trzyma zamek: ubij (obserwator albo zamknięcie okna) czyta go
// wtedy, gdy zwolnij (obserwator zakończenia procesu) może go właśnie oddawać.
// Bez zamka byłby to wyścig o pole `zadanie`.
//
// pid to identyfikator procesu okna utrwalony w chwili przejęcia. Dogląd
// posługuje się tym polem, a nie os.Process.Pid — to drugie zeruje Release
// (wywoływany przez Zwolnij) na wartość -1, więc czytanie go z gorutyny doglądu
// byłoby wyścigiem danych z ubiciem idącym równolegle. Pole zapisuje się raz,
// w przejmij, zanim struktura wyjdzie poza jedną gorutynę, i tylko czyta później.
type drzewoProcesow struct {
	mu      sync.Mutex
	zadanie windows.Handle
	pid     uint32
}

// atrybutyProcesu tworzy proces wstrzymany. Przypisanie do zadania musi nastąpić,
// zanim proces zdąży urodzić potomstwo — inaczej wnuk mógłby powstać poza
// zadaniem i przeżyć ubicie okna.
func atrybutyProcesu() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{CreationFlags: windows.CREATE_SUSPENDED}
}

// przygotujDrzewo zakłada zadanie zamykające potomstwo wraz z uchwytem.
// JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE gwarantuje, że nawet nagłe zakończenie
// rdzenia nie zostawi sierot.
func przygotujDrzewo() (*drzewoProcesow, error) {
	zadanie, err := windows.CreateJobObject(nil, nil)
	if err != nil {
		return nil, fmt.Errorf("session: nie można utworzyć Job Object: %w", err)
	}
	var limity windows.JOBOBJECT_EXTENDED_LIMIT_INFORMATION
	limity.BasicLimitInformation.LimitFlags = windows.JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE
	_, err = windows.SetInformationJobObject(
		zadanie,
		windows.JobObjectExtendedLimitInformation,
		uintptr(unsafe.Pointer(&limity)),
		uint32(unsafe.Sizeof(limity)),
	)
	if err != nil {
		windows.CloseHandle(zadanie)
		return nil, fmt.Errorf("session: nie można ustawić granic Job Object: %w", err)
	}
	return &drzewoProcesow{zadanie: zadanie}, nil
}

// przejmij przypisuje wstrzymany proces do zadania i wznawia jego pracę.
func (d *drzewoProcesow) przejmij(p *os.Process) error {
	d.pid = uint32(p.Pid)
	uchwyt, err := windows.OpenProcess(
		windows.PROCESS_SET_QUOTA|windows.PROCESS_TERMINATE, false, uint32(p.Pid))
	if err != nil {
		return fmt.Errorf("session: nie można otworzyć procesu %d: %w", p.Pid, err)
	}
	defer windows.CloseHandle(uchwyt)

	if err := windows.AssignProcessToJobObject(d.zadanie, uchwyt); err != nil {
		return fmt.Errorf("session: nie można przypisać procesu %d do Job Object: %w", p.Pid, err)
	}
	return wznowProces(uint32(p.Pid))
}

// ubij kończy zadanie, a wraz z nim proces okna i całe jego potomstwo. Uchwyt
// zadania oddany już przez zwolnij (zadanie == 0) znaczy, że proces zakończył
// się sam — nie ma czego ubijać, więc ubicie jest wtedy ciche.
func (d *drzewoProcesow) ubij(_ *os.Process) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.zadanie == 0 {
		return nil
	}
	if err := windows.TerminateJobObject(d.zadanie, 1); err != nil {
		return fmt.Errorf("session: nie można zakończyć Job Object: %w", err)
	}
	return nil
}

// zwolnij oddaje uchwyt zadania po zakończeniu procesu.
func (d *drzewoProcesow) zwolnij() {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.zadanie != 0 {
		windows.CloseHandle(d.zadanie)
		d.zadanie = 0
	}
}

// czekaj blokuje wywołującego do faktycznego zakończenia procesu. Otwiera własny
// uchwyt synchronizujący, żeby nie ruszać uchwytu trzymanego przez os.Process
// (ten oddaje zwolnij). Nieudane otwarcie znaczy, że proces już zniknął.
//
// Posługuje się zapamiętanym d.pid, nie os.Process.Pid — patrz opis pola.
func (d *drzewoProcesow) czekaj(p *os.Process) {
	if p == nil || d.pid == 0 {
		return
	}
	uchwyt, err := windows.OpenProcess(windows.SYNCHRONIZE, false, d.pid)
	if err != nil {
		return
	}
	defer windows.CloseHandle(uchwyt)
	_, _ = windows.WaitForSingleObject(uchwyt, windows.INFINITE)
}

// wznowProces wznawia główny wątek procesu utworzonego jako wstrzymany.
// Proces wstrzymany ma dokładnie jeden wątek, więc wystarczy pierwszy wątek
// należący do jego identyfikatora.
func wznowProces(pid uint32) error {
	migawka, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPTHREAD, 0)
	if err != nil {
		return fmt.Errorf("session: nie można odczytać wykazu wątków: %w", err)
	}
	defer windows.CloseHandle(migawka)

	var wpis windows.ThreadEntry32
	wpis.Size = uint32(unsafe.Sizeof(wpis))
	for err := windows.Thread32First(migawka, &wpis); err == nil; err = windows.Thread32Next(migawka, &wpis) {
		if wpis.OwnerProcessID != pid {
			continue
		}
		return wznowWatek(wpis.ThreadID, pid)
	}
	return fmt.Errorf("session: nie znaleziono wątku głównego procesu %d", pid)
}

// wznowWatek zdejmuje wstrzymanie z jednego wątku.
func wznowWatek(idWatku, pid uint32) error {
	watek, err := windows.OpenThread(windows.THREAD_SUSPEND_RESUME, false, idWatku)
	if err != nil {
		return fmt.Errorf("session: nie można otworzyć wątku %d procesu %d: %w", idWatku, pid, err)
	}
	defer windows.CloseHandle(watek)
	if _, err := windows.ResumeThread(watek); err != nil {
		return fmt.Errorf("session: nie można wznowić wątku %d procesu %d: %w", idWatku, pid, err)
	}
	return nil
}
