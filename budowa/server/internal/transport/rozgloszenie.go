package transport

import "danacoconsole/server/internal/protocol"

// Rozglos wysyła kopertę do wszystkich połączeń konta i zwraca liczbę urządzeń,
// które ją przyjęły. Puste konto oznacza wszystkie połączenia rdzenia.
//
// Synchronizacja wielourządzeniowa nie ma własnego protokołu: nośnikiem jest
// zdarzenie właściwe zmienionemu obszarowi, rozgłoszone tą drogą.
func (s *Serwer) Rozglos(konto string, k protocol.Koperta) int {
	return s.rozglosPoza(konto, "", k)
}

// rozglosPoza rozgłasza z pominięciem jednego połączenia — zwykle nadawcy
// zmiany, który wynik zna już z odpowiedzi na własną komendę.
func (s *Serwer) rozglosPoza(konto, pomijaneId string, k protocol.Koperta) int {
	dane, err := protocol.Zakoduj(k)
	if err != nil {
		s.ustawienia.Dziennik.Printf("transport: rozgłoszenie %s niezakodowane: %v", k.Type, err)
		return 0
	}
	odbiorcy := 0
	for _, p := range s.polaczenia.konta(konto) {
		if p.Id() == pomijaneId {
			continue
		}
		if err := p.wyslijBajty(dane); err != nil {
			s.ustawienia.Dziennik.Printf("transport: rozgłoszenie %s do %s pominięte: %v", k.Type, p.Id(), err)
			continue
		}
		odbiorcy++
	}
	return odbiorcy
}
