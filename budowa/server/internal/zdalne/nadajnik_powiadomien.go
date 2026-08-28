// Pakiet zdalne opisuje własny port nadajnika wołania urządzenia, niezależny od pakietu transportu, oraz prowadzi kolejkę i rozstrzyga kolejność doręczania powiadomień.
package zdalne

import (
	"context"

	"danacoconsole/server/internal/dane"
)

// Cel jest jednym adresem wołania odpowiadającym wierszowi rejestracji urządzenia powiadomień, złożonym z kanału doręczenia i klucza właściwego temu kanałowi.
type Cel struct {
	// Rejestracja jest kluczem wiersza rejestracji, którym ślad doręczenia wskazuje drogę powiadomienia.
	Rejestracja int64
	// Kanal jest dziś zawsze wartością połączenia — jedyną drogą, którą rdzeń umie.
	Kanal string
	// KluczKanalu dla kanału połączenia jest identyfikatorem klienta niesionym przez transport.
	KluczKanalu string
	// Etykieta jest napisem operatora o tym urządzeniu; może być pusta.
	Etykieta string
}

// Nadajnik doręcza powiadomienie pod jeden adres i zwraca prawdę wyłącznie wtedy, gdy urządzenie kopertę przyjęło, ponieważ odpowiedź niepewna zamieniłaby się w bazie w twierdzenie pewne.
type Nadajnik interface {
	Zawolaj(ctx context.Context, cel Cel, p dane.Powiadomienie) bool
}
