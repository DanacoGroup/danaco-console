package session

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

// RejestrProcesow trzyma procesy kluczowane identyfikatorem okna; sesja nie ma tu żadnego wpisu, ma go każde okno z osobna.
type RejestrProcesow struct {
	mu      sync.Mutex
	procesy map[string]*Proces
}

// Funkcja NowyRejestrProcesow zakłada pusty rejestr procesów okien, gotowy do przyjmowania kolejnych wpisów.
func NowyRejestrProcesow() *RejestrProcesow {
	return &RejestrProcesow{procesy: make(map[string]*Proces)}
}

// Metoda Proces zwraca proces powiązany z oknem o podanym identyfikatorze, o ile takie okno go posiada.
func (r *RejestrProcesow) Proces(idOkna string) (*Proces, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	proces, jest := r.procesy[idOkna]
	return proces, jest
}

// Metoda Zatrzymaj ubija proces okna wraz z całym drzewem potomstwa i wykreśla jego wpis z rejestru procesów.
func (r *RejestrProcesow) Zatrzymaj(idOkna string) error {
	r.mu.Lock()
	proces, jest := r.procesy[idOkna]
	delete(r.procesy, idOkna)
	r.mu.Unlock()
	if !jest {
		return fmt.Errorf("%w: %s", ErrProcesNieBiegnie, idOkna)
	}
	return proces.Ubij()
}

// ZatrzymajOkna ubija procesy wskazanych okien. Usterka jednego okna nie
// przerywa zamykania pozostałych.
func (r *RejestrProcesow) ZatrzymajOkna(idOkien []string) error {
	var pierwszaUsterka error
	for _, idOkna := range idOkien {
		if _, jest := r.Proces(idOkna); !jest {
			continue
		}
		if err := r.Zatrzymaj(idOkna); err != nil && pierwszaUsterka == nil {
			pierwszaUsterka = err
		}
	}
	return pierwszaUsterka
}

// Metoda ZatrzymajWszystkie ubija procesy wszystkich okien zapisanych obecnie w tym rejestrze procesów.
func (r *RejestrProcesow) ZatrzymajWszystkie() error {
	return r.ZatrzymajOkna(r.Okna())
}

// Metoda Okna wylicza identyfikatory wszystkich okien, które mają aktualnie wpis w tym rejestrze procesów.
func (r *RejestrProcesow) Okna() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	wykaz := make([]string, 0, len(r.procesy))
	for idOkna := range r.procesy {
		wykaz = append(wykaz, idOkna)
	}
	sort.Strings(wykaz)
	return wykaz
}

// Metoda Biegnace wylicza spośród okien tego rejestru te, których procesy nadal faktycznie pracują teraz.
func (r *RejestrProcesow) Biegnace() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	wykaz := make([]string, 0, len(r.procesy))
	for idOkna, proces := range r.procesy {
		if proces.Zyje() {
			wykaz = append(wykaz, idOkna)
		}
	}
	sort.Strings(wykaz)
	return wykaz
}

// Metoda Przejmij obejmuje biegnący już proces okna uchwytem systemowym i wpisuje go do rejestru procesów.
func (r *RejestrProcesow) Przejmij(idOkna string, pid int) error {
	if r == nil || idOkna == "" || pid <= 0 {
		return nil
	}
	drzewo, err := PrzejmijDrzewo(pid)
	if err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if poprzedni, jest := r.procesy[idOkna]; jest && poprzedni != nil {
		// Tura poprzednia tego okna już się zamknęła; zwalniamy jej uchwyt, by rejestr nie rósł.
		_ = poprzedni.Ubij()
	}
	proces := &Proces{
		IdOkna:      idOkna,
		IdProcesu:   nowyIdentyfikator(przedrostekProcesu),
		Uruchomiono: time.Now().UTC(),
		drzewo:      drzewo,
	}
	r.procesy[idOkna] = proces
	// Dogląd startuje po wstawieniu wpisu do mapy, jako jedyna droga obserwacji każdego procesu.
	go proces.dogladaj()
	return nil
}
