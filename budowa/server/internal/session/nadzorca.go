package session

// Nadzorca składa trzy części pakietu w jeden punkt wejścia dla rdzenia:
// rejestr sesji i okien, rejestr procesów kluczowany oknem oraz wybudzacz
// koordynatora. Rdzeń komponuje pakiet przez ten typ i nie sięga do wnętrza
// pozostałych.
type Nadzorca struct {
	rejestr   *Rejestr
	procesy   *RejestrProcesow
	wybudzacz *Wybudzacz
}

// NowyNadzorca składa nadzorcę. Wypełnienia z zewnątrz nadzorca nie przyjmuje,
// bo nie ma czym ich wypełnić: procesu okna sesja nie startuje ani nie czyta —
// startuje go warstwa kanału, a sesja obejmuje go uchwytem przez
// RejestrProcesow.Przejmij.
func NowyNadzorca() *Nadzorca {
	rejestr := NowyRejestr()
	wybudzacz := NowyWybudzacz(rejestr)
	return &Nadzorca{
		rejestr:   rejestr,
		procesy:   NowyRejestrProcesow(),
		wybudzacz: wybudzacz,
	}
}

// Rejestr zwraca rejestr sesji i okien.
func (n *Nadzorca) Rejestr() *Rejestr { return n.rejestr }

// Procesy zwraca rejestr procesów okien.
func (n *Nadzorca) Procesy() *RejestrProcesow { return n.procesy }

// Wybudzacz zwraca kierownicę pętli koordynator–wykonawca.
func (n *Nadzorca) Wybudzacz() *Wybudzacz { return n.wybudzacz }

// ZalozSesje zakłada sesję wspólną dla plików, pamięci, projektu i agentów.
func (n *Nadzorca) ZalozSesje(tytul, idProjektu string) Sesja {
	return n.rejestr.ZalozSesje(tytul, idProjektu)
}

// OtworzOkno zakłada okno komunikacji w sesji.
func (n *Nadzorca) OtworzOkno(idSesji string, u Ustawienia) (Okno, error) {
	return n.rejestr.OtworzOkno(idSesji, u)
}

// ZamknijOkno ubija proces okna wraz z drzewem potomstwa i zamyka okno.
// Pozostałe okna sesji pracują dalej bez zakłócenia.
func (n *Nadzorca) ZamknijOkno(idOkna string) (Okno, error) {
	if _, jest := n.procesy.Proces(idOkna); jest {
		if err := n.procesy.Zatrzymaj(idOkna); err != nil {
			return Okno{}, err
		}
	}
	return n.rejestr.ZamknijOkno(idOkna)
}

// ZamknijSesje zamyka wszystkie okna sesji i ubija ich procesy.
func (n *Nadzorca) ZamknijSesje(idSesji string) ([]Okno, error) {
	zamkniete, err := n.rejestr.ZamknijSesje(idSesji)
	if err != nil {
		return nil, err
	}
	idOkien := make([]string, 0, len(zamkniete))
	for _, okno := range zamkniete {
		idOkien = append(idOkien, okno.Id)
	}
	return zamkniete, n.procesy.ZatrzymajOkna(idOkien)
}

// Zamknij kończy wszystkie procesy okien — zamknięcie rdzenia nie zostawia
// sierot.
func (n *Nadzorca) Zamknij() error {
	return n.procesy.ZatrzymajWszystkie()
}
