package core

import (
	"strings"
	"testing"
)

// Pomiar odsłuchu jest osobny od pomiaru dyktowania — i musi być.
//
// `speech.availability.get` obiecuje modelowi sprawdzenie „dyktowania albo
// odsłuchu", a przez długi czas mierzyła sam łańcuch transkrypcji: na maszynie
// z Pythonem i modelem, lecz bez pipera i espeaka, meldowała gotowość, a odsłuch
// odmawiał. Pole `synthesisAvailable` zamyka ten rozjazd. Sprawdzian pilnuje, że
// pomiar odsłuchu mówi prawdę, gdy żadnego syntezatora nie ma.
//
// Brak wymuszony jest zmiennymi środowiska i odcięciem PATH, nie stanem maszyny:
// na stanowisku deweloperskim piper i espeak bywają doinstalowane ręcznie, więc
// sprawdzian liczący na ich nieobecność kłamałby tam, gdzie się go uruchamia.
func TestPomiarOdsluchuNazywaObaBrakiGdyZadnegoSyntezatoraNieMa(t *testing.T) {
	// Piper wskazany w nieistniejący plik, jego głosy w nieistniejący katalog,
	// espeak w nieistniejący plik — a PATH na katalog pusty, żeby goła nazwa
	// „piper" i „espeak-ng" też nie trafiła w nic.
	t.Setenv("DANACO_PIPER", "/nie/ma/takiego/pipera")
	t.Setenv("DANACO_PIPER_GLOSY", "/nie/ma/takiego/katalogu/glosow")
	t.Setenv("DANACO_ESPEAK", "/nie/ma/takiego/espeaka")
	t.Setenv("PATH", t.TempDir())

	gotowy, powod := gotowoscOdsluchu()

	if gotowy {
		t.Fatal("odsłuch zmierzony jako gotowy przy odciętych obu syntezatorach — " +
			"pomiar rozjechał się z wykonaniem, które odmówiłoby")
	}
	// Odmowa ma nazwać OBA silniki, nie jeden: Operator naprawia piper inaczej
	// (dołożenie głosu) niż espeak (instalacja pakietu), a podanie samego drugiego
	// wskazywałoby gorszą naprawę.
	if !strings.Contains(powod, "iper") {
		t.Errorf("powód niedostępności odsłuchu nie wymienia pipera: %q", powod)
	}
	if !strings.Contains(powod, "speak") {
		t.Errorf("powód niedostępności odsłuchu nie wymienia espeaka: %q", powod)
	}
}
