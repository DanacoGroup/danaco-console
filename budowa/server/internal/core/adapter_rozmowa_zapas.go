// Odpowiedzialność pliku: przełączenie tury na kanał zapasowy — gdy kanał
// odmawia, tura próbuje kanałem następnym, ale nigdy po cichu.
//
// Zapas bierze się z parametru `kanal_zapasowy` wiersza rejestru kanałów
// (wartość danych, jak `program` i `adapter`) — Operator ustawia go komendą
// `channel.update` polem config. Kolejność zapasowa nie jest polityką kodu: kod
// zna wyłącznie jeden krok „kanał → jego zapas".
//
// Jawność ma trzy nogi: fragment metadanych konta w strumieniu tury (Operator
// widzi przełączenie w rozmowie), wiersz w tabeli `przelaczenie_kanalu` jako
// ślad trwały i prowenancja drugiego wywołania (widać, czym tura faktycznie
// pojechała). Przełączenie bez którejkolwiek nogi byłoby przełączeniem
// po cichu.
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

// zapasMozliwy mówi, czy odmowę kanału wolno spróbować naprawić zapasem.
//
// Trzy warunki są granicami uczciwości: tura przerwana przez Operatora nie
// jest odmową kanału; tura, której tekst już poszedł do odbiorcy, nie może
// pojechać drugi raz (powtórzyłaby wypowiedź — ta sama reguła co przy rotacji
// kont, injection/strumien.go TekstPoszedl); tura zamknięta zdarzeniem result
// skończyła się po stronie modelu, więc nie ma czego ponawiać.
func zapasMozliwy(kontekst context.Context, przyczyna error, tekstPoszedl, zamknieta bool) bool {
	return przyczyna != nil && kontekst.Err() == nil && !tekstPoszedl && !zamknieta
}

// pojedzZapasem wykonuje jeden krok przełączenia: kanał wskazany w zapytaniu
// odmówił, więc tura jedzie jego kanałem zapasowym. Zwraca wynik drugiego
// wywołania i znacznik, czy przełączenie w ogóle zaszło.
//
// Krok jest jeden z zamysłu. Łańcuch zapasów (A→B→C…) wykonywałby turę
// kanałem odległym od wyboru Operatora o wiele decyzji, z których każda
// zapadłaby bez niego. Jeden krok jest widoczny i odwracalny; łańcuch to
// polityka, której kod nie zna.
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
		// Zapas wskazany, ale nieczynny — mówimy o tym w błędzie tury, zamiast
		// milczeć: Operator wskazał drogę, która dziś nie istnieje.
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
