package core

import (
	"sort"
	"strings"
	"testing"

	"danacoconsole/shared"
)

// Pokrycie kontraktu mierzone na ŻYWYM rejestrze zmontowanego rdzenia.
//
// Powód istnienia tego pliku jest konkretny. Pokrycie liczono dotąd czytaniem
// źródeł wyrażeniem `Zarejestruj\(shared\.(Command\w+)` — i ta miara kłamie
// w obie strony. Komendy wpinane przez parametr, a nie literałem
// (`zarejestrujAkcje(rejestr, p.Akcje, shared.CommandActionList)`
// w `kompozycja.go`), wyrażenie omija, więc `action.list` wychodził z pomiaru
// jako komenda bez obsługi, choć rdzeń odpowiada na nią od dawna. W drugą
// stronę: wywołanie `Zarejestruj` w gałęzi, do której montaż nigdy nie dochodzi,
// pomiar liczy jako pokrycie.
//
// Rejestr zna prawdę, bo to on rozstrzyga, czy komenda dostanie uchwyt, czy
// odpowiedź `*.unknown`. Sprawdzian pyta jego, a nie źródeł.
func TestKazdaKomendaKontraktuMaUchwytWRejestrze(t *testing.T) {
	zmontowany, _, _ := zmontujDoPomiaruSkutku(t)

	bezUchwytu := []string{}
	for _, komenda := range shared.WszystkieKomendy() {
		if _, jest := zmontowany.Rdzen.rejestr.Obsluga(komenda); !jest {
			bezUchwytu = append(bezUchwytu, string(komenda))
		}
	}
	sort.Strings(bezUchwytu)

	if len(bezUchwytu) > 0 {
		t.Errorf("komend kontraktu bez uchwytu: %d z %d\n%s",
			len(bezUchwytu), len(shared.WszystkieKomendy()), strings.Join(bezUchwytu, "\n"))
	}
}

// Rejestr nie ma prawa znać nazwy spoza kontraktu.
//
// Uchwyt pod nazwą, której kontrakt nie zna, jest komendą niewidoczną dla
// klienta i dla narzędzi modelu — martwym kodem, który wygląda na czynny.
func TestRejestrNieZnaNazwSpozaKontraktu(t *testing.T) {
	zmontowany, _, _ := zmontujDoPomiaruSkutku(t)

	kontrakt := make(map[shared.MessageType]struct{}, len(shared.WszystkieKomendy()))
	for _, komenda := range shared.WszystkieKomendy() {
		kontrakt[komenda] = struct{}{}
	}

	obce := []string{}
	for _, nazwa := range zmontowany.Rdzen.rejestr.Nazwy() {
		if _, jest := kontrakt[nazwa]; !jest {
			obce = append(obce, string(nazwa))
		}
	}
	sort.Strings(obce)

	if len(obce) > 0 {
		t.Errorf("rejestr zna nazwy spoza kontraktu: %s", strings.Join(obce, "\n"))
	}
}
