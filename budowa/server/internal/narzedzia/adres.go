// Odpowiedzialność pliku: adres rdzenia dla serwera narzędzi.
//
// ADRES BIERZE SIĘ STAMTĄD, SKĄD BIERZE GO RDZEŃ. Port czyta pakiet
// `konfiguracja` — ten sam, który ustala port nasłuchu procesu rdzenia — więc
// zmiana portu przez Operatora przestawia obie strony naraz. Odczyt zmiennej
// środowiska po nazwie dosłownej należy wyłącznie do tamtego pakietu, a serwer
// narzędzi dziedziczy środowisko po procesie modelu, który
// dziedziczy je po rdzeniu.
//
// PĘTLA ZWROTNA JEST WYBOREM, NIE SKRÓTEM. Serwer narzędzi stoi zawsze na tej
// samej maszynie co proces modelu, który go uruchomił, a proces modelu stoi
// przy rdzeniu, który go zrodził. Adres inny niż pętla zwrotna byłby wtedy
// zgadywaniem — Operator wskazuje go przełącznikiem, gdy układ jest inny.
package narzedzia

import (
	"net"
	"os"
	"strconv"

	"danacoconsole/server/internal/konfiguracja"
	"danacoconsole/server/internal/transport"
)

const (
	// schematGniazda jest schematem adresu gniazda rdzenia.
	schematGniazda = "ws://"
	// hostPetliZwrotnej jest domyślnym miejscem rdzenia widzianym z procesu modelu.
	hostPetliZwrotnej = "127.0.0.1"
)

// AdresRdzenia składa adres gniazda WebSocket rdzenia.
//
// Konfiguracja nieczytelna nie zatrzymuje serwera: adres schodzi na port
// domyślny, a o tym, czy rdzeń odpowiada, rozstrzyga dopiero pierwsze wywołanie.
func AdresRdzenia() string {
	port := konfiguracja.PortDomyslny
	if kon, err := konfiguracja.Wczytaj(nil, os.Getenv); err == nil {
		port = kon.Port
	}
	miejsce := net.JoinHostPort(hostPetliZwrotnej, strconv.Itoa(port))
	return schematGniazda + miejsce + transport.SciezkaGniazdaDomyslna
}
