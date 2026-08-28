// Odpowiedzialność pliku: sygnały błędów warstwy danych rozpoznawane przez
// warstwy wyższe. Rozpoznanie „wiersza nie ma” musi być odróżnialne od „odczyt
// się nie powiódł” — pierwsze bywa stanem normalnym (zapis dopiero powstanie),
// drugie jest awarią.
package dane

import (
	"errors"
	"strings"
)

// ErrBrakWiersza oznacza, że zapytanie nie znalazło wiersza. Warstwa wyższa
// odróżnia ten przypadek przez errors.Is, nie przez treść komunikatu.
var ErrBrakWiersza = errors.New("dane: brak wiersza")

// ErrKolizjaWiersza oznacza, że zapis rozbił się o jednoznaczność pilnowaną
// przez bazę — indeks UNIQUE albo klucz główny. Warstwa wyżej składa z tego
// odmowę kontraktu (`conflict`), zamiast przepuszczać komunikat silnika
// z nazwami tabel i kolumn.
var ErrKolizjaWiersza = errors.New("dane: kolizja z istniejącym wierszem")

// czyKolizja rozpoznaje naruszenie jednoznaczności po treści komunikatu
// sterownika SQLite, ponieważ sterownik nie udostępnia typowanego błędu
// dostępnego przez errors.As bez zależności od niego w tym pakiecie.
func czyKolizja(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "UNIQUE constraint failed") ||
		strings.Contains(err.Error(), "PRIMARY KEY must be unique")
}
