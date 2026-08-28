package session

import (
	"errors"

	"danacoconsole/shared"
)

// Katalog błędów pakietu niesie sytuacje wyjątkowe jako nazwane wartości, żeby
// warstwa wyżej rozpoznawała przyczynę przez porównanie błędów.
var (
	// ErrBrakSesji mówi, że wskazana sesja tej rozmowy w ogóle nie istnieje
	// w rejestrze sesji tej platformy.
	ErrBrakSesji = errors.New("session: sesja nie istnieje")
	// ErrBrakOkna mówi, że wskazane okno tej komunikacji w ogóle nie istnieje
	// dziś w rejestrze okien sesji.
	ErrBrakOkna = errors.New("session: okno komunikacji nie istnieje")
	// ErrBrakKoordynatora mówi, że okno wykonawcy wskazuje koordynatora,
	// którego dziś wcale nie ma w rejestrze.
	ErrBrakKoordynatora = errors.New("session: wskazane okno koordynatora nie istnieje")
	// errOknoZamkniete mówi, że dana czynność wymaga okna otwartego, a to
	// konkretne okno jest już zamknięte.
	errOknoZamkniete = errors.New("session: okno komunikacji jest zamknięte")
	// ErrProcesNieBiegnie mówi, że dane okno nie ma dziś żadnego uruchomionego
	// procesu wykonującego jego pracę.
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
