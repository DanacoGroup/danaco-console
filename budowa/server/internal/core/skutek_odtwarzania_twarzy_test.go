package core

import (
	"strings"
	"testing"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// Sprawdzian odmowy przebiegu twarzowego image.upscale przy braku wag, mierzonej bez ciężkich zasobów.

// TestBrakWagTwarzyOdmawiaNazywajacPlikIDrogeNaprawy pilnuje, żeby odmowa niosła nazwę brakującego pliku, ścieżkę, pod którą ma leżeć, oraz miejsce, z którego się go bierze.
func TestBrakWagTwarzyOdmawiaNazywajacPlikIDrogeNaprawy(t *testing.T) {
	pusty := t.TempDir()

	err := sprawdzWagiTwarzy(pusty)
	if err == nil {
		t.Fatal("katalog bez wag przeszedł sprawdzenie — przebieg ruszyłby i wywrócił się " +
			"dopiero we wnętrzu pomocnika, komunikatem Pythona zamiast odmową rdzenia")
	}
	blad := protocol.BladZeZrodla(shared.ErrorCodeInternalError, err)
	if blad.Code != shared.ErrorCodeChannelUnavailable {
		t.Fatalf("brak wag dostał kod %s, a brakująca instalacja jest zapleczem "+
			"niedostępnym (%s)", blad.Code, shared.ErrorCodeChannelUnavailable)
	}
	for _, oczekiwany := range []string{wagiOdtwarzaniaTwarzy, pusty, "github.com"} {
		if !strings.Contains(blad.Message, oczekiwany) {
			t.Fatalf("odmowa nie wymienia %q; treść: %s", oczekiwany, blad.Message)
		}
	}
}

// TestPomocnikTwarzyStoiWWykazieZaleznosci pilnuje, żeby silnik przebiegu twarzowego był widoczny w sondzie startowej, zanim jego brak ujawni się dopiero po wywołaniu.
func TestPomocnikTwarzyStoiWWykazieZaleznosci(t *testing.T) {
	szukany := narzedzieOdtwarzaniaTwarzy().Program
	if szukany == "" {
		t.Fatal("silnik przebiegu twarzowego nie ma nazwy programu")
	}
	for _, pozycja := range ZaleznosciZewnetrzne() {
		if pozycja.Narzedzie.Program == szukany {
			if strings.TrimSpace(pozycja.Zakres) == "" {
				t.Fatal("pozycja wykazu nie mówi, co przestaje działać przy braku")
			}
			return
		}
	}
	t.Fatalf("programu %s nie ma w wykazie zależności zewnętrznych", szukany)
}
