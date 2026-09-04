package transport

import "danacoconsole/server/internal/protocol"

// Rozglos wysyła kopertę do wszystkich połączeń wskazanego konta i zwraca liczbę urządzeń, które ją przyjęły, traktując puste konto jako wszystkie połączenia rdzenia.
func (s *Serwer) Rozglos(konto string, k protocol.Koperta) int {
	return s.rozglosPoza(konto, "", k)
}

// RozglosPoBramce pomija gniazda, które bramki jeszcze nie przeszły: telemetria
// opisuje pracę Operatora, choć nie ma zamawiającego ani konta.
//
// Bez wymogu logowania znacznika przejścia nie dostaje żadne gniazdo, więc
// rozgłoszenie idzie wtedy do wszystkich.
func (s *Serwer) RozglosPoBramce(konto string, k protocol.Koperta) int {
	if !s.ustawienia.dopuszczenie().wymagana {
		return s.Rozglos(konto, k)
	}
	dane, err := protocol.Zakoduj(k)
	if err != nil {
		s.ustawienia.Dziennik.Printf("transport: rozgłoszenie %s niezakodowane: %v", k.Type, err)
		return 0
	}
	odbiorcy := 0
	for _, p := range s.polaczenia.konta(konto) {
		if !p.przeszlaPrzezBramke() {
			continue
		}
		if err := p.wyslijBajty(dane); err != nil {
			s.ustawienia.Dziennik.Printf("transport: rozgłoszenie %s do %s pominięte: %v",
				k.Type, p.Id(), err)
			continue
		}
		odbiorcy++
	}
	return odbiorcy
}

// rozglosPoza rozgłasza kopertę z pominięciem jednego połączenia, zwykle nadawcy zmiany, który wynik zna już z odpowiedzi na własną komendę.
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
