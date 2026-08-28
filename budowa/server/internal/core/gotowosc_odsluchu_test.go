package core

import (
	"strings"
	"testing"
)

// Sprawdzian mierzy, że pomiar odsłuchu mówi prawdę i nazywa oba brakujące silniki, gdy
// żadnego syntezatora nie ma.
func TestPomiarOdsluchuNazywaObaBrakiGdyZadnegoSyntezatoraNieMa(t *testing.T) {
	// Piper i espeak wskazane w nieistniejące ścieżki, a PATH pusty, by gołe nazwy też nic
	// nie trafiły.
	t.Setenv("DANACO_PIPER", "/nie/ma/takiego/pipera")
	t.Setenv("DANACO_PIPER_GLOSY", "/nie/ma/takiego/katalogu/glosow")
	t.Setenv("DANACO_ESPEAK", "/nie/ma/takiego/espeaka")
	t.Setenv("PATH", t.TempDir())

	gotowy, powod := gotowoscOdsluchu()

	if gotowy {
		t.Fatal("odsłuch zmierzony jako gotowy przy odciętych obu syntezatorach — " +
			"pomiar rozjechał się z wykonaniem, które odmówiłoby")
	}
	// Odmowa nazywa oba silniki: naprawa pipera i espeaka różni się, więc jeden wskazałby
	// złą naprawę.
	if !strings.Contains(powod, "iper") {
		t.Errorf("powód niedostępności odsłuchu nie wymienia pipera: %q", powod)
	}
	if !strings.Contains(powod, "speak") {
		t.Errorf("powód niedostępności odsłuchu nie wymienia espeaka: %q", powod)
	}
}
