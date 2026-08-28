// Pakiet session prowadzi trzy byty przekroju wykonawczego: sesję wspólną dla
// plików i agentów, okno komunikacji oraz proces okna, wiązany z oknem, nie
// z sesją.
package session

import (
	"crypto/rand"
	"encoding/hex"
	"strconv"
	"sync/atomic"
)

// Przedrostki identyfikatorów. Rozpoznawalny przedrostek skraca diagnostykę
// dziennika — po samym identyfikatorze widać, jakiego bytu dotyczy wpis.
const (
	przedrostekSesji   = "ses"
	przedrostekOkna    = "okn"
	przedrostekProcesu = "prc"
)

// dlugoscLosowa jest liczbą bajtów losowych identyfikatora, dającą szesnaście
// znaków zapisu szesnastkowego.
const dlugoscLosowa = 8

// licznikZapasowy zasila identyfikator, gdy generator losowy systemu odmówi
// odpowiedzi. Identyfikator ma być unikalny, nie nieprzewidywalny.
var licznikZapasowy atomic.Uint64

// nowyIdentyfikator składa identyfikator bytu. Brak entropii nie zatrzymuje
// założenia sesji ani okna — zamiast błędu wchodzi licznik.
func nowyIdentyfikator(przedrostek string) string {
	los := make([]byte, dlugoscLosowa)
	if _, err := rand.Read(los); err != nil {
		return przedrostek + "_" + strconv.FormatUint(licznikZapasowy.Add(1), 16)
	}
	return przedrostek + "_" + hex.EncodeToString(los)
}
