// Adres rdzenia dla serwera narzędzi pochodzi z tego samego źródła
// konfiguracji, co port nasłuchu procesu rdzenia, i domyślnie wskazuje pętlę
// zwrotną, na której stoi proces modelu.
package narzedzia

import (
	"net"
	"os"
	"strconv"

	"danacoconsole/server/internal/konfiguracja"
	"danacoconsole/server/internal/transport"
)

const (
	// schematGniazda jest schematem adresu gniazda WebSocket rdzenia,
	// poprzedzającym host i port w zbudowanym adresie.
	schematGniazda = "ws://"
	// hostPetliZwrotnej jest domyślnym miejscem rdzenia widzianym z procesu
	// modelu, gdy oba stoją na tej samej maszynie.
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
