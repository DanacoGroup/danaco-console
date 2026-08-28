package core

import "sync"

// pamiecSesjiKolejek trzyma tymczasowe powiązanie kolejki z sesją i oknami kontraktu, dopóki sesja nie ma własnego wiersza w bazie i klucza obcego nie da się wypełnić.
type pamiecSesjiKolejek struct {
	mu    sync.RWMutex
	sesje map[int64]string
	okna  map[int64][]string
}

// nowaPamiecSesjiKolejek zakłada pustą pamięć powiązań kolejki z sesją i oknami, gotową do zapisywania i odczytu.
func nowaPamiecSesjiKolejek() *pamiecSesjiKolejek {
	return &pamiecSesjiKolejek{sesje: map[int64]string{}, okna: map[int64][]string{}}
}

// Zapamietaj zapisuje w pamięci powiązanie kolejki z sesją oraz z listą okien wskazanych przez klienta.
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
