package session

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

// RejestrProcesow trzyma procesy kluczowane identyfikatorem okna: sesja nie ma
// tu żadnego wpisu, ma go każde okno z osobna.
//
// Rejestr nie uruchamia procesów i nie zna drogi ich uruchomienia. Proces tury
// startuje warstwa kanału, a tutaj trafia przez Przejmij — objęty uchwytem
// drzewa, gotowy do zatrzymania.
type RejestrProcesow struct {
	mu      sync.Mutex
	procesy map[string]*Proces
}

// NowyRejestrProcesow zakłada pusty rejestr procesów okien.
func NowyRejestrProcesow() *RejestrProcesow {
	return &RejestrProcesow{procesy: make(map[string]*Proces)}
}

// Proces zwraca proces okna, o ile okno go ma.
func (r *RejestrProcesow) Proces(idOkna string) (*Proces, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	proces, jest := r.procesy[idOkna]
	return proces, jest
}

// Zatrzymaj ubija proces okna wraz z całym drzewem potomstwa i wykreśla go
// z rejestru.
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

// ZatrzymajWszystkie ubija procesy wszystkich okien rejestru.
func (r *RejestrProcesow) ZatrzymajWszystkie() error {
	return r.ZatrzymajOkna(r.Okna())
}

// Okna wylicza identyfikatory okien mających wpis w rejestrze procesów.
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

// Biegnace wylicza okna, których procesy nadal pracują.
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

// Przejmij obejmuje biegnący już proces okna uchwytem systemowym i wpisuje go
// do rejestru. Jest to jedyna droga wpisu do tego rejestru: proces tury startuje
// warstwa kanału (injection.Uruchom), a bez przejęcia zamknięcie okna nie
// zatrzymałoby tego, co model uruchomił.
//
// Przejęcie zakłada Job Object na Windows albo grupę procesów na systemach
// uniksowych, dzięki czemu ubicie okna kończy także wnuki procesu, bez `taskkill`
// i bez zależności od narzędzi systemu.
//
// Niepowodzenie przejęcia nie przerywa tury: proces biegnie i odpowiada, tylko
// jego potomstwo nie jest objęte uchwytem.
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
		// Tura poprzednia tego okna już się zamknęła; jej uchwyt zwalniamy,
		// żeby rejestr nie rósł o martwe wpisy przy każdej turze.
		_ = poprzedni.Ubij()
	}
	proces := &Proces{
		IdOkna:      idOkna,
		IdProcesu:   nowyIdentyfikator(przedrostekProcesu),
		Uruchomiono: time.Now().UTC(),
		drzewo:      drzewo,
	}
	r.procesy[idOkna] = proces
	// Dogląd startuje po wstawieniu wpisu do mapy i tylko tutaj. Przejmij jest
	// jedyną drogą wpisu do rejestru, więc obserwacja pokrywa każdy wpis i
	// żaden proces kończący się sam nie zostaje w rejestrze jako „biegnący".
	// Gorutyna nie sięga po r.mu, więc start pod zamkiem niczego nie blokuje.
	go proces.dogladaj()
	return nil
}
