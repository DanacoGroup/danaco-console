package core

import (
	"strings"
	"testing"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// Sprawdzian odmowy przebiegu twarzowego `image.upscale` przy braku wag.
//
// ── Dlaczego mierzymy odmowę, a nie skutek ──────────────────────────────────
// Skutek przebiegu twarzowego mierzy się zdjęciem twarzy, siecią liczącą
// minutami i trzema zestawami wag ważącymi pół gigabajta. Sprawdzian tego
// rodzaju byłby na maszynie bez wag „pominięty", czyli świeciłby na zielono, nie
// mierząc niczego — a to jest wprost ta klasa błędu, przed którą ostrzega ustrój
// budowy. Odmowa natomiast jest zachowaniem, które MA działać wszędzie i daje
// się zmierzyć na pustym katalogu.
//
// Skutek sieci na zdjęciu wykazuje się uruchomieniem na maszynie z wagami:
// `faces: false` i `faces: true` nad tym samym źródłem dają obrazy różne,
// a różnicę podaje się liczbą.

// TestBrakWagTwarzyOdmawiaNazywajacPlikIDrogeNaprawy pilnuje, żeby odmowa niosła
// trzy rzeczy naraz: nazwę brakującego pliku, ścieżkę, pod którą ma leżeć, oraz
// miejsce, z którego się go bierze. Odmowa mówiąca samo „brak wag" zostawia
// Operatora z pytaniem, na które ten kod zna odpowiedź.
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

// TestPomocnikTwarzyStoiWWykazieZaleznosci pilnuje, żeby silnik przebiegu
// twarzowego był widoczny w sondzie startowej. Program wołany przez rdzeń, ale
// nieobecny w wykazie, jest brakiem, o którym Operator dowiaduje się dopiero po
// naciśnięciu przycisku — a wykaz istnieje właśnie po to, żeby się nie dowiadywał
// tą drogą.
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
