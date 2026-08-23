package core

import "sync"

// pamiecSesjiKolejek trzyma powiązanie kolejki z sesją i oknami kontraktu.
//
// Powód istnienia jest jeden i tymczasowy: kolumna `kolejka.sesja_id` wskazuje
// wiersz sesji, a sesja żyje w pamięci pakietu sesji pod identyfikatorem
// tekstowym i wiersza nie ma. Dopóki tak jest, powiązanie mieszka tutaj — nie
// w bazie, bo klucza obcego nie da się wypełnić, i nie w kontrakcie, bo kontrakt
// jest w porządku. Trwałość sesji zdejmie ten plik w całości: powiązanie wróci
// wtedy do kolumny.
type pamiecSesjiKolejek struct {
	mu    sync.RWMutex
	sesje map[int64]string
	okna  map[int64][]string
}

// nowaPamiecSesjiKolejek zakłada pustą pamięć powiązań.
func nowaPamiecSesjiKolejek() *pamiecSesjiKolejek {
	return &pamiecSesjiKolejek{sesje: map[int64]string{}, okna: map[int64][]string{}}
}

// Zapamietaj zapisuje sesję i okna kolejki.
func (p *pamiecSesjiKolejek) Zapamietaj(idKolejki int64, idSesji string, idOkien []string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.sesje[idKolejki] = idSesji
	p.okna[idKolejki] = append([]string(nil), idOkien...)
}

// Odczytaj zwraca sesję i okna kolejki. Kolejka spoza pamięci — na przykład
// zapisana przed ponownym uruchomieniem rdzenia — wraca bez powiązania, a nie
// błędem.
func (p *pamiecSesjiKolejek) Odczytaj(idKolejki int64) (string, []string) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.sesje[idKolejki], p.okna[idKolejki]
}
