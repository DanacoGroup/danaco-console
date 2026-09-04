package session

import (
	"fmt"
	"sync"
	"time"

	"danacoconsole/shared"
)

// Rejestr trzyma sesje i ich okna komunikacji i jest bezpieczny dla wielu gorutyn pracujących równolegle.
type Rejestr struct {
	mu    sync.RWMutex
	sesje map[string]*Sesja
	okna  map[string]*Okno
}

// Funkcja NowyRejestr zakłada pusty rejestr sesji i okien, gotowy od razu do przyjmowania kolejnych wpisów.
func NowyRejestr() *Rejestr {
	return &Rejestr{
		sesje: make(map[string]*Sesja),
		okna:  make(map[string]*Okno),
	}
}

// Metoda ZalozSesje zakłada nową sesję czynną z podanym tytułem i identyfikatorem projektu i zwraca jej odpis.
func (r *Rejestr) ZalozSesje(tytul, idProjektu string, kontoId int64) Sesja {
	return r.ZalozSesjeSrodowiska(tytul, idProjektu, "", kontoId)
}

// Metoda ZalozSesjeSrodowiska zakłada sesję opisaną środowiskiem, przez które
// Operator wszedł do pracy. Po tym opisie karta środowiska liczy swoje sesje.
func (r *Rejestr) ZalozSesjeSrodowiska(tytul, idProjektu, kodSrodowiska string, kontoId int64) Sesja {
	s := nowaSesja(tytul, idProjektu)
	s.KodSrodowiska = kodSrodowiska
	s.KontoId = kontoId
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sesje[s.Id] = s
	return s.Kopia()
}

// Metoda Sesja zwraca odpis sesji o podanym identyfikatorze, bez modyfikacji jej bieżącego stanu w rejestrze.
func (r *Rejestr) Sesja(id string) (Sesja, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, jest := r.sesje[id]
	if !jest {
		return Sesja{}, fmt.Errorf("%w: %s", ErrBrakSesji, id)
	}
	return s.Kopia(), nil
}

// SesjeKonta zwraca odpisy sesji jednego konta. Wykaz Operatora idzie tędy,
// nie przez Sesje: rejestr trzyma sesje całej instalacji.
func (r *Rejestr) SesjeKonta(kontoId int64) []Sesja {
	wykaz := make([]Sesja, 0, len(r.sesje))
	for _, s := range r.Sesje() {
		if s.KontoId == kontoId {
			wykaz = append(wykaz, s)
		}
	}
	return wykaz
}

// Metoda Sesje zwraca odpisy wszystkich sesji przechowywanych obecnie w tym rejestrze sesji rdzenia aplikacji.
func (r *Rejestr) Sesje() []Sesja {
	r.mu.RLock()
	defer r.mu.RUnlock()
	wykaz := make([]Sesja, 0, len(r.sesje))
	for _, s := range r.sesje {
		wykaz = append(wykaz, s.Kopia())
	}
	return wykaz
}

// ZamknijSesje kończy sesję i zamyka wszystkie jej okna. Zwraca odpisy okien,
// które zostały zamknięte — wywołujący ubija ich procesy.
func (r *Rejestr) ZamknijSesje(id string) ([]Okno, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	s, jest := r.sesje[id]
	if !jest {
		return nil, fmt.Errorf("%w: %s", ErrBrakSesji, id)
	}
	zamkniete := make([]Okno, 0, len(s.IdOkien))
	for _, idOkna := range s.IdOkien {
		o, jest := r.okna[idOkna]
		if !jest || !o.CzyOtwarte() {
			continue
		}
		o.Stan = shared.WindowStatusClosed
		o.Zaktualizowano = time.Now().UTC()
		zamkniete = append(zamkniete, o.Kopia())
	}
	s.Stan = shared.SessionStatusFinished
	s.Zaktualizowano = time.Now().UTC()
	return zamkniete, nil
}

// WznowSesje otwiera zamkniętą sesję wraz z jej oknami — odwrotność
// ZamknijSesje.
//
// Rejestr żywy nie odróżnia okna zamkniętego razem z sesją od zamkniętego
// wcześniej osobno, więc wznowienie otwiera wszystkie zamknięte okna sesji.
func (r *Rejestr) WznowSesje(id string) (Sesja, []Okno, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	s, jest := r.sesje[id]
	if !jest {
		return Sesja{}, nil, fmt.Errorf("%w: %s", ErrBrakSesji, id)
	}
	otwarte := make([]Okno, 0, len(s.IdOkien))
	for _, idOkna := range s.IdOkien {
		o, jest := r.okna[idOkna]
		if !jest || o.CzyOtwarte() {
			continue
		}
		o.Stan = shared.WindowStatusOpen
		o.Zaktualizowano = time.Now().UTC()
		otwarte = append(otwarte, o.Kopia())
	}
	s.Stan = shared.SessionStatusActive
	s.Zaktualizowano = time.Now().UTC()
	return s.Kopia(), otwarte, nil
}

// Metoda UsunSesje wykreśla z rejestru sesję wraz z jej oknami i zwraca odpisy trwale usuniętych okien.
func (r *Rejestr) UsunSesje(id string) ([]Okno, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	s, jest := r.sesje[id]
	if !jest {
		return nil, fmt.Errorf("%w: %s", ErrBrakSesji, id)
	}
	usuniete := make([]Okno, 0, len(s.IdOkien))
	for _, idOkna := range s.IdOkien {
		if o, jest := r.okna[idOkna]; jest {
			usuniete = append(usuniete, o.Kopia())
			delete(r.okna, idOkna)
		}
	}
	delete(r.sesje, id)
	return usuniete, nil
}

// Odtworz wnosi do rejestru sesję i jej okna odczytane z trwałości.
// Identyfikatory pochodzą z bazy, nie z generatora — klient zna je sprzed
// restartu rdzenia i po restarcie musi trafić w te same byty.
func (r *Rejestr) Odtworz(s Sesja, okna []Okno) {
	r.mu.Lock()
	defer r.mu.Unlock()
	odpis := s.Kopia()
	odpis.IdOkien = nil
	r.sesje[odpis.Id] = &odpis
	for _, o := range okna {
		kopia := o.Kopia()
		r.okna[kopia.Id] = &kopia
		odpis.dopiszOkno(kopia.Id)
	}
}

// ZmienProjektSesji wiąże sesję z projektem albo wyjmuje ją z projektu.
// Wskazanie puste znaczy „sesja bez projektu" i jest stanem poprawnym.
func (r *Rejestr) ZmienProjektSesji(id, idProjektu string) (Sesja, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	sesja, jest := r.sesje[id]
	if !jest {
		return Sesja{}, fmt.Errorf("%w: %s", ErrBrakSesji, id)
	}
	sesja.IdProjektu = idProjektu
	sesja.Zaktualizowano = time.Now().UTC()
	return sesja.Kopia(), nil
}

// ZmienTytulSesji zmienia nazwę sesji w rejestrze żywym. Moduł, kanał
// i katalogi należą do okna, nie do sesji.
func (r *Rejestr) ZmienTytulSesji(id, tytul string) (Sesja, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	sesja, jest := r.sesje[id]
	if !jest {
		return Sesja{}, fmt.Errorf("%w: %s", ErrBrakSesji, id)
	}
	sesja.Tytul = tytul
	sesja.Zaktualizowano = time.Now().UTC()
	return sesja.Kopia(), nil
}

// Wykaz sesji czyta rejestr żywy, nie bazę, więc stan idzie i tutaj.
func (r *Rejestr) ZmienStanSesji(id string, stan shared.SessionStatus) (Sesja, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	sesja, jest := r.sesje[id]
	if !jest {
		return Sesja{}, fmt.Errorf("%w: %s", ErrBrakSesji, id)
	}
	sesja.Stan = stan
	sesja.Zaktualizowano = time.Now().UTC()
	return sesja.Kopia(), nil
}
