// Odpowiedzialność pliku: sprawdzenie klucza żądania `config.set` wobec
// katalogu ustawień, zanim wiersz trafi do bazy. Klucz spoza katalogu trafiłby
// do tabeli `ustawienie` i nigdy nie został odczytany — rozstrzygacz pyta
// wyłącznie o klucze katalogu — więc wywołujący dostałby potwierdzenie zapisu
// bez skutku.
//
// Nieznany klucz kończy się błędem jednego wywołania o kodzie
// ErrorCodeValidationFailed: rdzeń pracuje dalej, a wywołujący dostaje nazwany
// powód.
//
// Gdy rejestr definicji jest pusty — baza bez katalogu ustawień albo nieudany
// odczyt katalogu — sprawdzenia nie ma. Pusty rejestr znaczy „nie wiadomo, co
// jest znane", a nie „nic nie jest znane"; sprawdzenie oparte na takim rejestrze
// zablokowałoby całą konfigurację.
package core

import (
	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// sprawdzKluczKatalogu zwraca odmowę, gdy klucz nie występuje w rejestrze
// definicji rozstrzygacza — czyli w katalogu ustawień bazy, a przy jego braku
// w rejestrze wbudowanym rdzenia. Zgodność zwraca nil.
func sprawdzKluczKatalogu(rozstrzygacz *konfig.Rozstrzygacz, klucz string) error {
	rejestr := rozstrzygacz.Rejestr()
	if rejestr == nil || rejestr.Liczba() == 0 {
		return nil
	}
	if klucz == "" {
		return bladNieznanegoKlucza(klucz)
	}
	if _, znany := rejestr.Definicja(klucz); znany {
		return nil
	}
	return bladNieznanegoKlucza(klucz)
}

// bladNieznanegoKlucza składa odmowę: podaje odrzucony klucz oraz komendę,
// która zwraca wykaz kluczy dozwolonych.
func bladNieznanegoKlucza(klucz string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"konfiguracja: klucz \""+klucz+"\" nie należy do katalogu ustawień; "+
			"wykaz kluczy dozwolonych daje komenda settings.definition.list"))
}
