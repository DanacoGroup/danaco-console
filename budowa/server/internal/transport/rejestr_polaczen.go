package transport

import "sync"

// rejestrPolaczen zna wszystkie czynne połączenia rdzenia i przeszukuje je po koncie, nie po pojedynczym połączeniu.
type rejestrPolaczen struct {
	zamek      sync.RWMutex
	polaczenia map[string]*Polaczenie
}

// Funkcja nowyRejestrPolaczen tworzy pusty rejestr połączeń, gotowy do przyjmowania kolejnych wpisów rdzenia.
func nowyRejestrPolaczen() *rejestrPolaczen {
	return &rejestrPolaczen{polaczenia: make(map[string]*Polaczenie)}
}

// Metoda dodaj wpisuje nowo nawiązane połączenie urządzenia do tego rejestru czynnych połączeń rdzenia.
func (r *rejestrPolaczen) dodaj(p *Polaczenie) {
	r.zamek.Lock()
	r.polaczenia[p.Id()] = p
	r.zamek.Unlock()
}

// Metoda usun wykreśla z rejestru połączenie po rozłączeniu urządzenia, które z tego połączenia korzystało.
func (r *rejestrPolaczen) usun(id string) {
	r.zamek.Lock()
	delete(r.polaczenia, id)
	r.zamek.Unlock()
}

// Metoda liczba zwraca liczbę czynnych połączeń zapisanych obecnie w tym rejestrze połączeń tego rdzenia.
func (r *rejestrPolaczen) liczba() int {
	r.zamek.RLock()
	defer r.zamek.RUnlock()
	return len(r.polaczenia)
}

// wszystkie zwraca kopię zbioru połączeń. Kopia, nie odwzorowanie: wysyłka idzie
// poza blokadą, więc rozłączenie w trakcie rozgłoszenia niczego nie zakleszcza.
func (r *rejestrPolaczen) wszystkie() []*Polaczenie {
	r.zamek.RLock()
	defer r.zamek.RUnlock()
	lista := make([]*Polaczenie, 0, len(r.polaczenia))
	for _, p := range r.polaczenia {
		lista = append(lista, p)
	}
	return lista
}

// Metoda konta zwraca połączenia wskazanego konta; puste konto oznacza wszystkie połączenia tego rdzenia.
func (r *rejestrPolaczen) konta(konto string) []*Polaczenie {
	wszystkie := r.wszystkie()
	if konto == "" {
		return wszystkie
	}
	lista := make([]*Polaczenie, 0, len(wszystkie))
	for _, p := range wszystkie {
		if p.Konto() == konto {
			lista = append(lista, p)
		}
	}
	return lista
}

// Metoda zamknijWszystkie kończy wszystkie połączenia rejestru; używana przy zatrzymaniu całego rdzenia.
func (r *rejestrPolaczen) zamknijWszystkie(powod string) {
	for _, p := range r.wszystkie() {
		p.Zamknij(powod)
	}
}
