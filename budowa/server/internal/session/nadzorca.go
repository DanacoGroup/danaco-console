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

// Rejestr zwraca rejestr sesji i okien, przez który nadzorca prowadzi ich
// stan, przydział i wyszukiwanie.
func (n *Nadzorca) Rejestr() *Rejestr { return n.rejestr }

// Procesy zwraca rejestr procesów okien, którymi zarządza nadzorca
// w imieniu rdzenia i całej platformy.
func (n *Nadzorca) Procesy() *RejestrProcesow { return n.procesy }

// Wybudzacz zwraca kierownicę pętli koordynator–wykonawca obsługującej
// sesję komunikacji z modelem językowym.
func (n *Nadzorca) Wybudzacz() *Wybudzacz { return n.wybudzacz }

// ZalozSesje zakłada sesję wspólną dla plików, pamięci, projektu i agentów,
// zwracając jej pełny opis startowy.
func (n *Nadzorca) ZalozSesje(tytul, idProjektu string) Sesja {
	return n.rejestr.ZalozSesje(tytul, idProjektu)
}

// ZalozSesjeSrodowiska zakłada sesję opisaną środowiskiem wejścia Operatora.
func (n *Nadzorca) ZalozSesjeSrodowiska(tytul, idProjektu, kodSrodowiska string) Sesja {
	return n.rejestr.ZalozSesjeSrodowiska(tytul, idProjektu, kodSrodowiska)
}

// OtworzOkno zakłada okno komunikacji w sesji i zwraca jego opis wraz
// z identyfikatorem oraz ustawieniami.
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

// ZamknijSesje zamyka wszystkie okna sesji i ubija ich procesy wraz z całym
// drzewem procesów potomnych.
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

// Zamknij kończy wszystkie procesy okien wszystkich sesji rdzenia;
// zamknięcie rdzenia nie zostawia sierot.
func (n *Nadzorca) Zamknij() error {
	return n.procesy.ZatrzymajWszystkie()
}
