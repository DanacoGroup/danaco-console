// Pakiet session prowadzi trzy byty przekroju wykonawczego Danaco Console:
// sesję — wspólną dla plików, pamięci, projektu i agentów, okno
// komunikacji — nośnik parametrów wykonania, oraz proces okna. Proces wiąże
// się z oknem, nie z sesją: jedna sesja prowadzi wiele okien i wiele procesów
// biegnących równolegle, każde z własnym kanałem modelu i własnym katalogiem
// roboczym.
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

// dlugoscLosowa — liczba bajtów losowych identyfikatora (16 znaków szesnastkowych).
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
