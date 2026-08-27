package session

import (
	"fmt"
	"time"

	"danacoconsole/shared"
)

// Rejestr okien zbiera metody okna jako jednostki wykonania; sesja prowadzi wiele okien naraz.

// Metoda OtworzOkno zakłada okno komunikacji w sesji wskazanej identyfikatorem i zwraca odpis nowo otwartego okna.
func (r *Rejestr) OtworzOkno(idSesji string, u Ustawienia) (Okno, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	s, jest := r.sesje[idSesji]
	if !jest {
		return Okno{}, fmt.Errorf("%w: %s", ErrBrakSesji, idSesji)
	}
	o := noweOkno(idSesji, u)
	if err := r.sprawdzKoordynatora(*o); err != nil {
		return Okno{}, err
	}
	r.okna[o.Id] = o
	s.dopiszOkno(o.Id)
	return o.Kopia(), nil
}

// Metoda Okno zwraca odpis okna o podanym identyfikatorze, bez modyfikacji jego stanu bieżącego w rejestrze.
func (r *Rejestr) Okno(id string) (Okno, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	o, jest := r.okna[id]
	if !jest {
		return Okno{}, fmt.Errorf("%w: %s", ErrBrakOkna, id)
	}
	return o.Kopia(), nil
}

// Metoda OknaSesji zwraca odpisy wszystkich okien wskazanej sesji, ułożone w kolejności ich pierwotnego otwarcia.
func (r *Rejestr) OknaSesji(idSesji string) ([]Okno, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, jest := r.sesje[idSesji]
	if !jest {
		return nil, fmt.Errorf("%w: %s", ErrBrakSesji, idSesji)
	}
	wykaz := make([]Okno, 0, len(s.IdOkien))
	for _, idOkna := range s.IdOkien {
		if o, jest := r.okna[idOkna]; jest {
			wykaz = append(wykaz, o.Kopia())
		}
	}
	return wykaz, nil
}

// ZmienOkno nanosi wybiórczą zmianę ustawień okna — obsługa window.update.
// Zmiana jednego okna nie rusza pozostałych okien sesji.
func (r *Rejestr) ZmienOkno(id string, z Zmiana) (Okno, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	o, jest := r.okna[id]
	if !jest {
		return Okno{}, fmt.Errorf("%w: %s", ErrBrakOkna, id)
	}
	poprzednie := o.Kopia()
	o.zastosuj(z)
	if err := r.sprawdzKoordynatora(*o); err != nil {
		*o = poprzednie
		return Okno{}, err
	}
	return o.Kopia(), nil
}

// ZamknijOkno zamyka okno. Sesja trwa dalej — zamknięcie okna nie kończy ani
// sesji, ani pozostałych okien.
func (r *Rejestr) ZamknijOkno(id string) (Okno, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	o, jest := r.okna[id]
	if !jest {
		return Okno{}, fmt.Errorf("%w: %s", ErrBrakOkna, id)
	}
	if o.CzyOtwarte() {
		o.Stan = shared.WindowStatusClosed
		o.Zaktualizowano = time.Now().UTC()
	}
	return o.Kopia(), nil
}

// Metoda Wykonawcy zwraca odpisy okien wykonawczych podległych koordynatorowi wskazanemu identyfikatorem.
func (r *Rejestr) Wykonawcy(idKoordynatora string) []Okno {
	r.mu.RLock()
	defer r.mu.RUnlock()
	wykaz := make([]Okno, 0, len(r.okna))
	for _, o := range r.okna {
		if o.CzyWykonawca() && o.OknoKoordynatora == idKoordynatora {
			wykaz = append(wykaz, o.Kopia())
		}
	}
	return wykaz
}

// sprawdzKoordynatora pilnuje, by wskazanie koordynatora prowadziło do
// istniejącego okna. Wywoływane pod założoną blokadą rejestru.
func (r *Rejestr) sprawdzKoordynatora(o Okno) error {
	if o.OknoKoordynatora == "" {
		return nil
	}
	if _, jest := r.okna[o.OknoKoordynatora]; !jest {
		return fmt.Errorf("%w: %s", ErrBrakKoordynatora, o.OknoKoordynatora)
	}
	return nil
}
