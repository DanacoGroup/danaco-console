// Odpowiedzialność pliku: stan żywy modułu Terminal — karty powłok i procesy biegnące
// w tej chwili wraz z uchwytami do ich drzew potomstwa, jako ewidencja, nie uruchamianie.
package core

import (
	"sync"
	"sync/atomic"
	"time"

	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// kartaTerminala jest profilem powłoki jednej karty okna Terminal Tabs, wraz z celem zdalnym powłoki karty.
type kartaTerminala struct {
	kod        string
	oknoKod    string
	idSesji    string
	powloka    shared.TerminalShell
	tytul      string
	katalog    string
	srodowisko map[string]string
	stan       shared.TerminalSessionStatus
	utworzono  time.Time
	// celZdalny to adres powłoki zdalnej karty, trwały w polu, a nie w zmiennej środowiska.
	celZdalny string
	// portZdalny bierze port domyślny protokołu, gdy jest zerowy.
	portZdalny int
	// hostKod wskazuje wpis książki hostów, z którego karta wzięła adres.
	hostKod string
	// kluczSciezka to ścieżka klucza prywatnego; pusta znaczy klucz domyślny maszyny rdzenia.
	kluczSciezka string
	// Wskazanie celu powłok urządzeniowych: kontenera, poda, klastra i portu szeregowego.
	kontener        string
	pod             string
	przestrzenNazw  string
	kontekstKlastra string
	urzadzenie      string
	predkoscPortu   int
}

// procesTerminala jest jednym przebiegiem polecenia wraz z uchwytami, bez których nie da się go zakończyć.
type procesTerminala struct {
	kod          string
	kartaKod     string
	oknoKod      string
	idSesji      string
	polecenie    string
	inicjator    shared.ProcessInitiator
	pid          int
	pidNadrzedny int

	mu sync.Mutex
	// wstrzymany mówi, czy drzewo procesu stoi wstrzymane, choć stan kontraktu wciąż jest "running".
	wstrzymany  bool
	stan        shared.TerminalProcessStatus
	kodWyjscia  *int
	uruchomiono time.Time
	zakonczono  time.Time

	uchwyt session.UchwytProcesu
	drzewo *session.DrzewoProcesu
	koniec chan struct{}

	// numerFragmentu numeruje fragmenty strumienia wyjścia, bo dwie gorutyny czytają je naraz.
	numerFragmentu atomic.Int64
}

// rejestrTerminala trzyma karty i procesy czynne jednego biegu rdzenia, w pamięci procesu serwera produktu.
type rejestrTerminala struct {
	mu      sync.Mutex
	karty   map[string]*kartaTerminala
	procesy map[string]*procesTerminala
}

func nowyRejestrTerminala() *rejestrTerminala {
	return &rejestrTerminala{
		karty:   make(map[string]*kartaTerminala),
		procesy: make(map[string]*procesTerminala),
	}
}

// ZapiszKarte wstawia albo podmienia kartę w ewidencji rejestru terminala rdzenia produktu Danaco Console.
func (r *rejestrTerminala) ZapiszKarte(karta *kartaTerminala) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.karty[karta.kod] = karta
}

// Karta zwraca kartę o wskazanym kodzie, jeśli jest w ewidencji rejestru terminala rdzenia tego serwera.
func (r *rejestrTerminala) Karta(kod string) (*kartaTerminala, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	karta, jest := r.karty[kod]
	return karta, jest
}

// Karty zwraca wszystkie karty ewidencji, na czym stoi żądanie `terminal.session.list` z kontraktu terminala.
func (r *rejestrTerminala) Karty() []*kartaTerminala {
	r.mu.Lock()
	defer r.mu.Unlock()
	wykaz := make([]*kartaTerminala, 0, len(r.karty))
	for _, karta := range r.karty {
		wykaz = append(wykaz, karta)
	}
	return wykaz
}

// ZapiszProces wstawia proces do ewidencji procesów czynnych rejestru terminala rdzenia produktu Danaco.
func (r *rejestrTerminala) ZapiszProces(proces *procesTerminala) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.procesy[proces.kod] = proces
}

// Proces zwraca proces o wskazanym kodzie, jeśli jest w ewidencji rejestru terminala rdzenia tego serwera.
func (r *rejestrTerminala) Proces(kod string) (*procesTerminala, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	proces, jest := r.procesy[kod]
	return proces, jest
}

// Procesy zwraca wszystkie procesy ewidencji rejestru terminala, czynne i zakończone naraz, bez rozróżnienia.
func (r *rejestrTerminala) Procesy() []*procesTerminala {
	r.mu.Lock()
	defer r.mu.Unlock()
	wykaz := make([]*procesTerminala, 0, len(r.procesy))
	for _, proces := range r.procesy {
		wykaz = append(wykaz, proces)
	}
	return wykaz
}

// Przytnij usuwa najstarszy przebieg zakończony, gdy ewidencja przekracza pojemność, tak
// by Process Monitor mógł odpowiedzieć tak samo, jak z trwałym dziennikiem w bazie.
func (r *rejestrTerminala) Przytnij() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.procesy) <= pojemnoscHistorii {
		return
	}
	najstarszy, chwila := "", time.Time{}
	for kodWpisu, proces := range r.procesy {
		stan, _, _ := proces.Migawka()
		if stan == shared.TerminalProcessStatusRunning {
			continue
		}
		if najstarszy == "" || proces.uruchomiono.Before(chwila) {
			najstarszy, chwila = kodWpisu, proces.uruchomiono
		}
	}
	// Gdy wszystkie wpisy są czynne, nie ma czego przyciąć.
	if najstarszy != "" {
		delete(r.procesy, najstarszy)
	}
}

// pojemnoscHistorii ogranicza liczbę procesów trzymanych w pamięci. Powyżej
// progu wypada najstarszy przebieg zakończony; procesu czynnego nie wykreśla
// nic poza jego własnym końcem.
const pojemnoscHistorii = 500

// Zamknij kończy procesy czynne wraz z ich potomstwem; wywołuje się przy
// zamykaniu serwera, żeby nie zostawić procesów osieroconych. Niepowodzenie
// jednego zakończenia nie wstrzymuje pozostałych.
func (r *rejestrTerminala) Zamknij() {
	for _, proces := range r.Procesy() {
		if stan, _, _ := proces.Migawka(); stan != shared.TerminalProcessStatusRunning {
			continue
		}
		_ = proces.Zakoncz(true)
	}
}
