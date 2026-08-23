package session

import (
	"errors"

	"danacoconsole/shared"
)

// Katalog błędów pakietu. Sytuacje wyjątkowe są nazwane wartościami, żeby
// warstwa wyżej rozpoznawała przyczynę przez errors.Is, a nie przez treść
// komunikatu. Odwzorowanie na kody kontraktu leży tutaj — pakiet session nie
// zakłada własnego katalogu kodów.
var (
	// ErrBrakSesji — wskazana sesja nie istnieje w rejestrze.
	ErrBrakSesji = errors.New("session: sesja nie istnieje")
	// ErrBrakOkna — wskazane okno komunikacji nie istnieje w rejestrze.
	ErrBrakOkna = errors.New("session: okno komunikacji nie istnieje")
	// ErrBrakKoordynatora — okno wykonawcy wskazuje koordynatora, którego nie ma.
	ErrBrakKoordynatora = errors.New("session: wskazane okno koordynatora nie istnieje")
	// errOknoZamkniete — czynność wymaga okna otwartego.
	errOknoZamkniete = errors.New("session: okno komunikacji jest zamknięte")
	// ErrProcesNieBiegnie — okno nie ma uruchomionego procesu.
	ErrProcesNieBiegnie = errors.New("session: okno nie ma biegnącego procesu")
)

// Kod przekłada błąd pakietu na kod błędu kontraktu. Nierozpoznany błąd jest
// usterką wewnętrzną rdzenia — kończy bieżące wywołanie, nie sesję.
func Kod(err error) shared.ErrorCode {
	switch {
	case err == nil:
		return ""
	case errors.Is(err, ErrBrakSesji),
		errors.Is(err, ErrBrakOkna),
		errors.Is(err, ErrBrakKoordynatora),
		errors.Is(err, ErrProcesNieBiegnie):
		return shared.ErrorCodeNotFound
	case errors.Is(err, errOknoZamkniete):
		return shared.ErrorCodeConflict
	default:
		return shared.ErrorCodeInternalError
	}
}
