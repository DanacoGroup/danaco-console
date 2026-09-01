package transport

import "github.com/coder/websocket"

// RozlaczKonta zamyka połączenia konta wskazane przez wybierz i zwraca ich liczbę; puste konto znaczy wszystkie połączenia rdzenia. Kod zamknięcia jest ten sam, co przy sesji, która przestała nadawać w bramce: urządzenie ma otworzyć okno logowania, nie wznawiać połączenia.
func (s *Serwer) RozlaczKonta(konto, powod string, wybierz func(id string) bool) int {
	if wybierz == nil {
		return 0
	}
	zamkniete := 0
	for _, p := range s.polaczenia.konta(konto) {
		if !wybierz(p.Id()) {
			continue
		}
		p.ZamknijKodem(websocket.StatusPolicyViolation, powod)
		zamkniete++
	}
	return zamkniete
}
