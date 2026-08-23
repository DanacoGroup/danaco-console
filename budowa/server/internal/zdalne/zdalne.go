// Pakiet zdalne obsługuje to, co leży poza procesem rdzenia. Odpowiedzialności
// są dwie:
//
//  1. Tor do hosta zdalnego — przekłada polecenie procesu okna na wywołanie
//     SSH, którym proces rusza na maszynie wskazanej przez Operatora, oraz
//     przenosi pliki tym samym torem. Pliki: tor.go, hosty.go, polecenie.go,
//     pliki.go.
//  2. Dosięgnięcie Operatora — kolejka powiadomień, rejestracja urządzeń do
//     wołania, ponowienie i wygaśnięcie. Pliki: powiadomienia.go,
//     nadajnik_powiadomien.go, budzik_powiadomien.go. Silnik doręcza po łączu,
//     które Operator już otworzył; wygaszonego telefonu nie budzi — wymagałoby
//     to usługi wypychania powiadomień spoza tego systemu.
//
// SSH jest tu transportem, nie drugim wykonawcą. Pakiet nie uruchamia procesów:
// buduje wyłącznie wiersz poleceń transportu, a startuje go jedyny spawner
// platformy (injection.Wystartuj). Po stronie zdalnej proces uruchamia sshd —
// usługa, którą host już wystawia — więc w drzewie nie przybywa żaden własny
// demon ani protokół. Strumienie SSH są strumieniami procesu zdalnego: wejście,
// wyjście, wyjście diagnostyczne i kod zakończenia przechodzą wprost, a zerwanie
// połączenia kończy proces po stronie zdalnej (SIGHUP od sshd). Rola `agent`
// produktu (tor stdio, cmd/danaco-console/uruchomienie/) komponuje się z tym
// torem bez zmian: `ssh host danaco-console --role agent` wykonuje żądania
// kontraktu na hoście zdalnym tym samym binarium.
//
// Granice:
//   - Pid uchwytu jest identyfikatorem lokalnego procesu transportu (ssh),
//     nie procesu zdalnego; drzewo potomstwa obejmowane przez warstwę sesji
//     kończy transport, a proces zdalny kończy sshd po zerwaniu połączenia.
//   - Dziedziczenie środowiska rdzenia dotyczy maszyny rdzenia: na hoście
//     zdalnym proces dziedziczy środowisko logowania SSH, a wpisy własne
//     polecenia jadą w komendzie zdalnej.
//   - Zgoda na hosta jest wierszem tabeli `host_zdalny`, wydawanym przez
//     Operatora instrukcją Danaco — domyślnie jej nie ma. Rdzeń nie zainicjuje
//     połączenia, którego Operator nie oddał.
//
// Pakiet czyta bazę rdzenia przez uchwyt podany w Zasil, a każda droga wywołana
// przed zasileniem odmawia, nazywając brakujące wpięcie. Uchwyt podaje
// kompozycja (server/cmd/danaco-console/main.go), więc odmowa braku zasilenia
// dotyczy wołaczy spoza kompozycji i sprawdzianów.
//
// Nadajnik powiadomień (`zdalne.ZasilNadajnik`) jest osobnym wpięciem. Bez niego
// przebieg kolejki odmawia w całości i nie tyka ani jednego wiersza —
// powiadomienia czekają z pełnym budżetem prób, zamiast po cichu wygasać.
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

// baza zwraca uchwyt bazy albo nil, gdy wpięcie jeszcze nie nastąpiło.
func baza() *sql.DB {
	zasilenie.RLock()
	defer zasilenie.RUnlock()
	return zasilenie.baza
}
