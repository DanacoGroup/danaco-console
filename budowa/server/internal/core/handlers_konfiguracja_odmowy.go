// Plik ustala wspólne słownictwo odmów okna konfiguracji, jednakowe dla
// czterech rodzin komend, które to okno składają.
package core

import (
	"errors"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// bladWskazania przekłada brak wiersza na odmowę wskazania. Każdy inny błąd
// przechodzi nietknięty — jest awarią odczytu, nie pomyłką Operatora.
func bladWskazania(err error, byt, identyfikator string) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			byt+" "+identyfikator+" nie istnieje"))
	}
	return err
}

// bladBrakuKatalogu odmawia zapisu, gdy domena nie ma wpiętego repozytorium.
// Odczyt w tej samej sytuacji zwraca wykaz pusty — czytać nie ma czego, ale
// zapis, który nigdzie nie trafia, byłby potwierdzeniem czynności niewykonanej.
func bladBrakuKatalogu(byt string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"rdzeń: katalog "+byt+" nie jest wpięty"))
}

// brakWiersza odpowiada, czy błąd oznacza wyłącznie nieobecność wiersza.
// Czynności odwracalne — usunięcie, odczyt jednego bytu — kończą się wtedy
// wynikiem pustym, a nie odmową.
func brakWiersza(err error) bool {
	return errors.Is(err, dane.ErrBrakWiersza)
}
