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

// nadawcaStrumienia zamienia fragmenty kanału modelu na koperty stream.chunk i oddaje je transportowi, numerując fragmenty i domykając strumień dokładnie jednym znacznikiem końca dla całej tury.
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

// nowyNadawcaStrumienia zakłada nadawcę dla jednej tury, powiązanego z zadaniem, sesją i transportem nadajnika.
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
		// Fragment po znaczniku końca nie ma dokąd pójść: klient zamknął już wpis tury, więc milczy.
		return nil
	}
	n.wyslijZawieszony(false)
	n.zawieszony, n.jestZawiesz = &f, true
	return nil
}

// Zakoncz domyka strumień jednym zdarzeniem domykającym: fragment ostateczny zastępuje tekst złożony z fragmentów, bo ten jest przybliżeniem, a fragment błędu domyka strumień tą samą drogą, gdy tura nie dobiegła końca powodzeniem.
func (n *nadawcaStrumienia) Zakoncz(idOkna, idWiadomosci, wersjaOstateczna string, przyczyna error) {
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.domkniety {
		return
	}
	// Fragment w zawieszeniu wychodzi jako zwykły, nie ostatni: domknięcie idzie zdarzeniem niżej.
	n.wyslijZawieszony(false)
	if przyczyna != nil {
		blad := protocol.BladZeZrodla(shared.ErrorCodeChannelUnavailable, przyczyna)
		fragment := protocol.ChunkBledu(idOkna, idWiadomosci, blad)
		n.zawieszony, n.jestZawiesz = &fragment, true
	} else {
		// Tura bez ani jednego znaku domyka się pustą wersją ostateczną, by wpis nie został w stanie strumień.
		fragment := protocol.ChunkWersjiOstatecznej(idOkna, idWiadomosci, wersjaOstateczna)
		n.zawieszony, n.jestZawiesz = &fragment, true
	}
	n.wyslijZawieszony(true)
	n.domkniety = true
}

// DomknijAwaryjnie domyka strumień, jeżeli nikt nie domknął go wcześniej. Stoi w defer tury i odzywa się wyłącznie, gdy tura wyszła drogą nieprzewidzianą: panika w kanale, w rejestratorze bloków albo w koordynatorze pętli.
func (n *nadawcaStrumienia) DomknijAwaryjnie(idOkna, idWiadomosci string) {
	n.mu.Lock()
	domkniety := n.domkniety
	n.mu.Unlock()
	if domkniety {
		return
	}
	n.Zakoncz(idOkna, idWiadomosci, "", bladTuryPrzerwanej)
}

// wyslijZawieszony wysyła fragment trzymany w zawieszeniu, jako zwykły albo ostatni. Wywoływać wyłącznie pod zamkiem.
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
