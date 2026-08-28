// Plik przełącza turę na kanał zapasowy: gdy kanał wskazany w zapytaniu
// odmawia, tura próbuje kanałem następnym, ale nigdy po cichu.
package core

import (
	"context"
	"fmt"
	"strings"

	"danacoconsole/server/internal/models"
)

// parametrKanaluZapasowego jest kluczem parametru wiersza rejestru kanałów
// wskazującego kanał zapasowy. Wartość danych, nie stała polityki.
const parametrKanaluZapasowego = "kanal_zapasowy"

// zapasMozliwy mówi, czy odmowę kanału wolno spróbować naprawić zapasem,
// według trzech warunków granic uczciwości opisanych osobno.
func zapasMozliwy(kontekst context.Context, przyczyna error, tekstPoszedl, zamknieta bool) bool {
	return przyczyna != nil && kontekst.Err() == nil && !tekstPoszedl && !zamknieta
}

// pojedzZapasem wykonuje jeden krok przełączenia na kanał zapasowy, gdy kanał
// wskazany w zapytaniu odmówił, zwracając wynik drugiego wywołania i znacznik,
// czy przełączenie zaszło.
func (a *adapterRozmowy) pojedzZapasem(kontekst context.Context, zapytanie models.Zapytanie,
	ujscie models.Ujscie, przyczyna error) (error, bool) {

	kanal, jest := a.kanaly.Kanal(zapytanie.Kanal)
	if !jest {
		return przyczyna, false
	}
	zapas := strings.TrimSpace(kanal.Definicja().ParametrLub(parametrKanaluZapasowego, ""))
	if zapas == "" || zapas == zapytanie.Kanal {
		return przyczyna, false
	}
	if _, czynny := a.kanaly.Kanal(zapas); !czynny {
		// Zapas wskazany, ale nieczynny: mówimy o tym w błędzie tury, zamiast milczeć.
		return fmt.Errorf("%w; kanał zapasowy %q nie stoi w rejestrze albo jest nieczynny", przyczyna, zapas), false
	}

	powod := fmt.Sprintf("kanał %s odmówił: %v", zapytanie.Kanal, przyczyna)
	// Noga pierwsza: fragment w strumieniu — przełączenie widoczne w rozmowie.
	if fragment, err := models.FragmentKonta(zapytanie, models.MetadaneKonta{
		Powod: powod + " — tura jedzie kanałem zapasowym " + zapas,
	}); err == nil {
		_ = ujscie.Fragment(kontekst, fragment)
	}
	// Noga druga: ślad trwały w tabeli `przelaczenie_kanalu`.
	if a.zdarzenia != nil {
		a.zdarzenia.ZanotujPrzelaczenie(zapytanie.Okno(), zapytanie.Wiadomosc,
			zapytanie.Kanal, zapas, powod)
	}
	// Noga trzecia: prowenancję drugiego wywołania nadaje kanał zapasowy sam.
	zapytanie.Kanal = zapas
	return a.kanaly.Wyslij(kontekst, zapytanie, ujscie), true
}
