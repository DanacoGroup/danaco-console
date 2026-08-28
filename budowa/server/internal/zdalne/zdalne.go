// Pakiet zdalne obsługuje to, co leży poza procesem rdzenia: tor do hosta zdalnego oraz dosięgnięcie operatora powiadomieniami.
package zdalne

import (
	"database/sql"
	"sync"
)

// zasilenie trzyma uchwyt bazy rdzenia. Pakiet nie otwiera bazy samodzielnie —
// plik SQLite ma jedną pulę połączeń w całym procesie, a jej
// właścicielem jest kompozycja (main.go), nie warstwa kanału.
var zasilenie struct {
	sync.RWMutex
	baza *sql.DB
}

// Zasil podaje pakietowi uchwyt otwartej bazy rdzenia. Wywołanie należy do
// kompozycji (server/cmd/danaco-console/main.go), po kontroli spójności bazy.
func Zasil(baza *sql.DB) {
	zasilenie.Lock()
	defer zasilenie.Unlock()
	zasilenie.baza = baza
}

// Odetnij zdejmuje uchwyt bazy — dla porządku sprawdzianów, które zasilają
// pakiet własną bazą tymczasową i nie mogą zostawić jej następcom.
func Odetnij() {
	Zasil(nil)
}

// Funkcja baza zwraca uchwyt otwartej bazy rdzenia albo pustą wartość, gdy wpięcie jeszcze nie nastąpiło.
func baza() *sql.DB {
	zasilenie.RLock()
	defer zasilenie.RUnlock()
	return zasilenie.baza
}
