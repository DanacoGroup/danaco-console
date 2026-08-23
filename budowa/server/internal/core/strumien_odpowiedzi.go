package core

import (
	"context"
	"errors"
	"sync"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// bladTuryPrzerwanej wyjaśnia domknięcie awaryjne. Zdanie czyta Operator
// w oknie rozmowy, więc jest po polsku i mówi, co dalej.
var bladTuryPrzerwanej = errors.New(
	"tura przerwana usterką rdzenia — strumień odpowiedzi domknięto bez wyniku; " +
		"powtórz wysłanie wiadomości, a jeśli powtórzy się to samo, zajrzyj do " +
		"Errors Panel po zapis usterki")

// nadawcaStrumienia zamienia fragmenty kanału modelu na koperty stream.chunk
// i oddaje je transportowi.
//
// Numer fragmentu i znacznik końca żyją w kopercie, nie w ładunku, więc nadawca
// je nadaje: numeruje od jedynki i domyka strumień znacznikiem ostatniego.
// Ostatni fragment poznaje się dopiero wtedy, gdy nadejdzie następny albo
// skończy się tura — dlatego nadawca trzyma jeden fragment w zawieszeniu.
//
// Identyfikator jest niezmienny przez cały strumień: `idZadania` to
// identyfikator żądania, które turę otworzyło — ten sam, którym wróciła
// odpowiedź na `message.send`. Kontrakt wymaga tego wprost („odpowiedz i
// fragmenty strumienia powtarzaja identyfikator zadania"), a nadawca dostaje
// go raz, przy założeniu, i nie ma metody pozwalającej go podmienić.
//
// Domknięcie jest dokładnie jedno: znacznik `done` wychodzi raz na turę. Pole
// `domkniety` zamyka nadawcę na stałe, więc ani powtórzone `Zakoncz`, ani
// domknięcie awaryjne po panice nie wystawią drugiego końca tego samego
// strumienia, a fragment przyjęty po domknięciu jest odrzucany.
type nadawcaStrumienia struct {
	nadajnik  Nadajnik
	idZadania string
	idSesji   string

	mu          sync.Mutex
	numer       int
	zawieszony  *protocol.Chunk
	jestZawiesz bool
	domkniety   bool
}

// nowyNadawcaStrumienia zakłada nadawcę dla jednej tury.
func nowyNadawcaStrumienia(nadajnik Nadajnik, idZadania, idSesji string) *nadawcaStrumienia {
	return &nadawcaStrumienia{nadajnik: nadajnik, idZadania: idZadania, idSesji: idSesji}
}

// Fragment przyjmuje fragment kanału. Wypełnia interfejs models.Ujscie.
// Nie zwraca błędu: niepowodzenie wysyłki dotyczy jednego fragmentu i nie ma
// prawa przerwać tury.
func (n *nadawcaStrumienia) Fragment(_ context.Context, f protocol.Chunk) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.domkniety {
		// Fragment po znaczniku końca nie ma dokąd pójść: klient zamknął już
		// wpis tury. Milczenie jest tu poprawną odpowiedzią — jeden zbłąkany
		// fragment nie ma prawa przerwać niczego.
		return nil
	}
	n.wyslijZawieszony(false)
	n.zawieszony, n.jestZawiesz = &f, true
	return nil
}

// Zakoncz domyka strumień jednym zdarzeniem domykającym.
//
// Wersja ostateczna zastępuje fragmenty, a nie dokleja się do nich: fragment
// domykający niesie rodzaj `final` i całą treść odpowiedzi, a odbiorca podmienia
// nią tekst złożony z fragmentów (ChunkKind.final w kontrakcie,
// `zlozenie-tury.ts` po stronie klienta).
//
// Zastąpienie, a nie doklejenie, bo tekst złożony z fragmentów jest
// przybliżeniem: kanał może fragment powtórzyć po rotacji konta, może zerwać go
// w połowie znaku wielobajtowego, a tura zapasowa nadaje własne. Wersja
// ostateczna jest jedyną prawdą o tym, co model powiedział — tą samą, którą
// dziennik rozmowy zapisuje jako `content` wiadomości.
//
// Błąd domyka strumień tak samo, ale innym rodzajem: fragment rodzaju „błąd"
// jest wtedy zdarzeniem domykającym, bo strumień ma jedną drogę dla powodzenia
// i niepowodzenia. Wersji ostatecznej wtedy nie ma — tekst, który zdążył dojść,
// zostaje na wpisie taki, jaki jest.
func (n *nadawcaStrumienia) Zakoncz(idOkna, idWiadomosci, wersjaOstateczna string, przyczyna error) {
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.domkniety {
		return
	}
	// Fragment trzymany w zawieszeniu wychodzi jako zwykły, nie jako ostatni:
	// domknięcie należy do zdarzenia domykającego nadawanego niżej.
	n.wyslijZawieszony(false)
	if przyczyna != nil {
		blad := protocol.BladZeZrodla(shared.ErrorCodeChannelUnavailable, przyczyna)
		fragment := protocol.ChunkBledu(idOkna, idWiadomosci, blad)
		n.zawieszony, n.jestZawiesz = &fragment, true
	} else {
		// Tura bez ani jednego znaku też domyka się wersją ostateczną — pustą.
		// Bez tego klient czekałby na znacznik końca, który nigdy nie przyjdzie,
		// a wpis stałby w stanie `strumien`.
		fragment := protocol.ChunkWersjiOstatecznej(idOkna, idWiadomosci, wersjaOstateczna)
		n.zawieszony, n.jestZawiesz = &fragment, true
	}
	n.wyslijZawieszony(true)
	n.domkniety = true
}

// DomknijAwaryjnie domyka strumień, jeżeli nikt nie domknął go wcześniej.
//
// Stoi w `defer` tury i milczy w przypadku zwykłym — tam domknięcie należy do
// `Zakoncz`, które zna przyczynę i potrafi ją nazwać. Odzywa się wyłącznie
// wtedy, gdy tura wyszła drogą nieprzewidzianą: panika w kanale, w rejestratorze
// bloków albo w koordynatorze pętli. Bez tego wyjścia panika zostawia klienta
// z wpisem w stanie `strumien` do końca sesji, bez odpowiedzi i bez wyjaśnienia.
//
// Fragment domykający niesie rodzaj „błąd", bo tura nie dobiegła końca. Treść
// odmowy mówi, co odmówiło, dlaczego i co z tym zrobić.
func (n *nadawcaStrumienia) DomknijAwaryjnie(idOkna, idWiadomosci string) {
	n.mu.Lock()
	domkniety := n.domkniety
	n.mu.Unlock()
	if domkniety {
		return
	}
	n.Zakoncz(idOkna, idWiadomosci, "", bladTuryPrzerwanej)
}

// wyslijZawieszony wysyła fragment trzymany w zawieszeniu. Wywoływać wyłącznie
// pod zamkiem.
func (n *nadawcaStrumienia) wyslijZawieszony(ostatni bool) {
	if !n.jestZawiesz || n.nadajnik == nil {
		n.zawieszony, n.jestZawiesz = nil, false
		return
	}
	n.numer++
	k, err := protocol.KopertaFragmentu(n.idZadania, n.idSesji, n.numer, ostatni, *n.zawieszony)
	n.zawieszony, n.jestZawiesz = nil, false
	if err != nil {
		return
	}
	n.nadajnik.Rozglos(k)
}
