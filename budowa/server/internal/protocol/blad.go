package protocol

import (
	"errors"

	"danacoconsole/shared"
)

// KodBledu jest maszynowym kodem błędu protokołu. Katalog kodów należy
// w całości do kontraktu — pakiet protocol nie zna ani jednej jego wartości
// i żadnej nie wybiera za wywołującego.
type KodBledu = shared.ErrorCode

// Blad jest treścią błędu odpowiedzi w kształcie kontraktu. Pole Retryable
// odróżnia usterkę przejściową od trwałej: błąd techniczny dotyczy wyłącznie
// bieżącego wywołania i nigdy nie blokuje sesji, konta ani kolejnych prób.
type Blad = shared.ErrorInfo

// NowyBlad tworzy błąd o podanym kodzie i komunikacie dla użytkownika.
// Ponawialność bierze wprost z katalogu kontraktu, a nie z decyzji
// wywołującego.
func NowyBlad(kod KodBledu, komunikat string) Blad {
	return Blad{Code: kod, Message: komunikat, Retryable: shared.KodyPonawialne[kod]}
}

// Opis składa czytelny tekst błędu z kodu kontraktu i komunikatu.
func Opis(b Blad) string {
	return string(b.Code) + ": " + b.Message
}

// JakoError przenosi błąd kontraktu jako zwykły błąd Go, dzięki czemu warstwy
// wyżej przenoszą go bez tłumaczenia na własne typy.
func JakoError(b Blad) error {
	return bladGo{blad: b}
}

// BladZeZrodla przenosi dowolny błąd Go do protokołu, zachowując kod, który
// błędowi już nadano. Zwraca zerowy Blad dla braku błędu.
func BladZeZrodla(kod KodBledu, err error) Blad {
	if err == nil {
		return Blad{}
	}
	var niesiony bladGo
	if errors.As(err, &niesiony) {
		return niesiony.blad
	}
	return NowyBlad(kod, err.Error())
}

// bladGo niesie błąd kontraktu w interfejsie error. Jest opakowaniem
// zachowania, nie kształtem komunikatu — na drut idzie wyłącznie Blad.
type bladGo struct {
	blad Blad
}

// Error czyni bladGo zwykłym błędem Go.
func (b bladGo) Error() string {
	return Opis(b.blad)
}
