package transport

import "sync"

// rejestrPolaczen zna wszystkie czynne połączenia rdzenia. Jedno konto ma wiele
// urządzeń równocześnie, więc rejestr nie jest odwzorowaniem konto→połączenie,
// lecz zbiorem połączeń przeszukiwanym po koncie.
//
// Konto połączenia zmienia się w czasie życia (PrzypiszKonto), dlatego rejestr
// nie kopiuje go do własnego indeksu — odczytuje wprost z połączenia.
type rejestrPolaczen struct {
	zamek      sync.RWMutex
	polaczenia map[string]*Polaczenie
}

// nowyRejestrPolaczen tworzy pusty rejestr.
func nowyRejestrPolaczen() *rejestrPolaczen {
	return &rejestrPolaczen{polaczenia: make(map[string]*Polaczenie)}
}

// dodaj wpisuje połączenie do rejestru.
func (r *rejestrPolaczen) dodaj(p *Polaczenie) {
	r.zamek.Lock()
	r.polaczenia[p.Id()] = p
	r.zamek.Unlock()
}

// usun wykreśla połączenie po rozłączeniu urządzenia.
func (r *rejestrPolaczen) usun(id string) {
	r.zamek.Lock()
	delete(r.polaczenia, id)
	r.zamek.Unlock()
}

// liczba zwraca liczbę czynnych połączeń.
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

// konta zwraca połączenia wskazanego konta. Puste konto oznacza wszystkie
// połączenia rdzenia.
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

// zamknijWszystkie kończy wszystkie połączenia — używane przy zatrzymaniu rdzenia.
func (r *rejestrPolaczen) zamknijWszystkie(powod string) {
	for _, p := range r.wszystkie() {
		p.Zamknij(powod)
	}
}
