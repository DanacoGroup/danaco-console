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

// czyKolizja rozpoznaje naruszenie jednoznaczności po komunikacie sterownika.
// Sterownik SQLite nie wystawia typowanego błędu, do którego dałoby się dojść
// przez errors.As bez wciągania go do tego pakietu, więc rozpoznanie idzie po
// treści — i jest zamknięte w tym miejscu, żeby nie rozlało się po warstwach
// wyżej.
func czyKolizja(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "UNIQUE constraint failed") ||
		strings.Contains(err.Error(), "PRIMARY KEY must be unique")
}
