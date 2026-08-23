package session

import (
	"fmt"
	"hash"
	"hash/fnv"
	"sync"
	"time"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// Strumień daje koordynatorowi pełny obraz pracy wykonawcy: tok rozumowania,
// wywołania narzędzi, wyniki i pliki. Jest obserwatorem zwykłej drogi fragmentu,
// a nie drugim silnikiem — fragment idzie do interfejsu swoją drogą, a tutaj
// zostaje jego odpis dla koordynatora.
//
// Kontekst koordynatora jest skończony, więc starsza część strumienia zwija się
// do liczb — ile wpisów, ile znaków i w jakim przedziale czasu. Ostatnie wpisy
// w liczbie PojemnoscJawnaDomyslna zostają w całości.

// PojemnoscJawnaDomyslna — ile ostatnich wpisów strumienia zostaje jawnych.
const PojemnoscJawnaDomyslna = 64

// WpisStrumienia jest jednym fragmentem wykonawcy widzianym przez koordynatora.
type WpisStrumienia struct {
	// IdOkna — okno wykonawcy, z którego fragment pochodzi.
	IdOkna string
	// IdWiadomosci — tura, w której fragment powstał.
	IdWiadomosci string
	// Rodzaj fragmentu ze słownika kontraktu.
	Rodzaj shared.ChunkKind
	// Tresc tekstowa fragmentu; pusta dla fragmentów niosących wyłącznie dane.
	Tresc string
	// Czas zapisania wpisu.
	Czas time.Time
}

// Zwiniete opisuje starszą część strumienia sprowadzoną do liczb.
type Zwiniete struct {
	Wpisow int
	Znakow int
	Od     time.Time
	Do     time.Time
}

// MigawkaStrumienia jest odpisem strumienia wykonawców jednego koordynatora.
type MigawkaStrumienia struct {
	// IdKoordynatora — okno, dla którego strumień jest prowadzony.
	IdKoordynatora string
	// Jawne — ostatnie wpisy w całości, w kolejności napływu.
	Jawne []WpisStrumienia
	// Zwiniete — starsza część strumienia sprowadzona do liczb.
	Zwiniete Zwiniete
	// Wpisow — łączna liczba wpisów, także zwiniętych.
	Wpisow int
}

// StrumienWykonawcy gromadzi fragmenty wykonawców w torach koordynatorów.
type StrumienWykonawcy struct {
	mu        sync.Mutex
	pojemnosc int
	tory      map[string]*torStrumienia
}

// NowyStrumienWykonawcy zakłada strumień o podanej pojemności części jawnej.
// Pojemność niedodatnia schodzi na wartość domyślną.
func NowyStrumienWykonawcy(pojemnosc int) *StrumienWykonawcy {
	if pojemnosc <= 0 {
		pojemnosc = PojemnoscJawnaDomyslna
	}
	return &StrumienWykonawcy{pojemnosc: pojemnosc, tory: map[string]*torStrumienia{}}
}

// Dopisz zapisuje fragment wykonawcy w torze jego koordynatora.
func (s *StrumienWykonawcy) Dopisz(idKoordynatora string, w WpisStrumienia) {
	if idKoordynatora == "" {
		return
	}
	if w.Czas.IsZero() {
		w.Czas = time.Now().UTC()
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tor(idKoordynatora).dopisz(w, s.pojemnosc)
}

// DopiszFragment zapisuje fragment w kształcie kontraktu. Jest to
// droga używana przez warstwę rozmowy: ten sam fragment, który jedzie do
// interfejsu, zostawia tu odpis dla koordynatora.
func (s *StrumienWykonawcy) DopiszFragment(idKoordynatora string, f protocol.Chunk) {
	s.Dopisz(idKoordynatora, WpisStrumienia{
		IdOkna:       f.WindowId,
		IdWiadomosci: f.MessageId,
		Rodzaj:       f.Kind,
		Tresc:        protocol.Tresc(f),
	})
}

// Migawka zwraca niezależny odpis toru koordynatora.
func (s *StrumienWykonawcy) Migawka(idKoordynatora string) MigawkaStrumienia {
	s.mu.Lock()
	defer s.mu.Unlock()
	tor, jest := s.tory[idKoordynatora]
	if !jest {
		return MigawkaStrumienia{IdKoordynatora: idKoordynatora}
	}
	return MigawkaStrumienia{
		IdKoordynatora: idKoordynatora,
		Jawne:          append([]WpisStrumienia(nil), tor.jawne...),
		Zwiniete:       tor.zwiniete,
		Wpisow:         tor.wpisow,
	}
}

// Odetnij zamyka bieżący odcinek strumienia i zwraca jego odcisk. Odcisk jest
// miarą postępu obiegu: dwa obiegi o tym samym odcisku znaczą, że
// wykonawca powtórzył się co do znaku.
func (s *StrumienWykonawcy) Odetnij(idKoordynatora string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.tor(idKoordynatora).odetnij()
}

// Zapomnij usuwa tor zamkniętego koordynatora.
func (s *StrumienWykonawcy) Zapomnij(idKoordynatora string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.tory, idKoordynatora)
}

// tor zwraca tor koordynatora, zakładając go przy pierwszym wpisie.
// Wywoływane pod założoną blokadą.
func (s *StrumienWykonawcy) tor(idKoordynatora string) *torStrumienia {
	tor, jest := s.tory[idKoordynatora]
	if !jest {
		tor = &torStrumienia{odcinek: fnv.New64a()}
		s.tory[idKoordynatora] = tor
	}
	return tor
}

// torStrumienia jest strumieniem wykonawców jednego koordynatora.
type torStrumienia struct {
	jawne         []WpisStrumienia
	zwiniete      Zwiniete
	wpisow        int
	odcinek       hash.Hash64
	odcinekWpisow int
}

// dopisz dokłada wpis i zwija najstarszy, gdy część jawna przekroczy pojemność.
func (t *torStrumienia) dopisz(w WpisStrumienia, pojemnosc int) {
	t.jawne = append(t.jawne, w)
	t.wpisow++
	t.odcinekWpisow++
	_, _ = t.odcinek.Write([]byte(string(w.Rodzaj) + "\x00" + w.Tresc + "\x00"))
	for len(t.jawne) > pojemnosc {
		t.zwin(t.jawne[0])
		t.jawne = t.jawne[1:]
	}
}

// zwin przenosi wpis z części jawnej do podsumowania liczbowego.
func (t *torStrumienia) zwin(w WpisStrumienia) {
	if t.zwiniete.Wpisow == 0 {
		t.zwiniete.Od = w.Czas
	}
	t.zwiniete.Wpisow++
	t.zwiniete.Znakow += len(w.Tresc)
	t.zwiniete.Do = w.Czas
}

// odetnij zwraca odcisk bieżącego odcinka i zaczyna odcinek następny.
// Odcisk niesie także liczbę wpisów, więc powtórzenie tej samej treści
// większą liczbę razy nie ginie w sumie kontrolnej.
func (t *torStrumienia) odetnij() string {
	odcisk := fmt.Sprintf("%d:%016x", t.odcinekWpisow, t.odcinek.Sum64())
	t.odcinek = fnv.New64a()
	t.odcinekWpisow = 0
	return odcisk
}
