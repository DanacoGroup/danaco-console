package session

import (
	"fmt"
	"sync"
	"time"

	"danacoconsole/shared"
)

// Rejestr trzyma sesje i ich okna komunikacji. Jest bezpieczny dla wielu
// gorutyn, bo okna jednej sesji pracują równolegle i sięgają po
// rejestr z różnych wątków obsługi połączeń.
//
// Rejestr wydaje wyłącznie odpisy bytów. Dzięki temu wywołujący nie zmieni
// stanu rejestru przez wskaźnik trzymany po stronie warstwy wyżej.
type Rejestr struct {
	mu    sync.RWMutex
	sesje map[string]*Sesja
	okna  map[string]*Okno
}

// NowyRejestr zakłada pusty rejestr.
func NowyRejestr() *Rejestr {
	return &Rejestr{
		sesje: make(map[string]*Sesja),
		okna:  make(map[string]*Okno),
	}
}

// ZalozSesje zakłada sesję czynną i zwraca jej odpis.
func (r *Rejestr) ZalozSesje(tytul, idProjektu string) Sesja {
	s := nowaSesja(tytul, idProjektu)
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sesje[s.Id] = s
	return s.Kopia()
}

// Sesja zwraca odpis sesji o podanym identyfikatorze.
func (r *Rejestr) Sesja(id string) (Sesja, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, jest := r.sesje[id]
	if !jest {
		return Sesja{}, fmt.Errorf("%w: %s", ErrBrakSesji, id)
	}
	return s.Kopia(), nil
}

// Sesje zwraca odpisy wszystkich sesji rejestru.
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

// UsunSesje wykreśla sesję wraz z jej oknami. Zwraca odpisy usuniętych okien.
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
