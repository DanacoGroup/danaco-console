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

// Strumień daje koordynatorowi obraz pracy wykonawcy: rozumowanie, wywołania narzędzi, wyniki, pliki.

// PojemnoscJawnaDomyslna określa, ile ostatnich wpisów strumienia zostaje jawnych, gdy nie podano innej wartości.
const PojemnoscJawnaDomyslna = 64

// WpisStrumienia jest jednym fragmentem wykonawcy widzianym przez koordynatora w trakcie trwania obiegu.
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

// Zwiniete opisuje starszą część strumienia sprowadzoną do liczb wpisów, znaków i przedziału czasu ich powstania.
type Zwiniete struct {
	Wpisow int
	Znakow int
	Od     time.Time
	Do     time.Time
}

// MigawkaStrumienia jest odpisem strumienia wykonawców jednego koordynatora w chwili jego faktycznego pobrania.
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

// StrumienWykonawcy gromadzi fragmenty wykonawców w osobnych torach, po jednym torze na każdego koordynatora.
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

// Metoda Dopisz zapisuje fragment wykonawcy w torze jego koordynatora, zwijając najstarsze wpisy w razie potrzeby.
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

// Metoda Migawka zwraca niezależny odpis toru koordynatora wskazanego identyfikatorem, bez zmiany jego stanu.
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

// Metoda Zapomnij usuwa z pamięci strumienia tor zamkniętego koordynatora, wskazanego jego identyfikatorem.
func (s *StrumienWykonawcy) Zapomnij(idKoordynatora string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.tory, idKoordynatora)
}

// Metoda tor zwraca tor koordynatora, zakładając go przy pierwszym wpisie; wywoływana pod założoną blokadą strumienia.
func (s *StrumienWykonawcy) tor(idKoordynatora string) *torStrumienia {
	tor, jest := s.tory[idKoordynatora]
	if !jest {
		tor = &torStrumienia{odcinek: fnv.New64a()}
		s.tory[idKoordynatora] = tor
	}
	return tor
}

// torStrumienia jest strumieniem wykonawców należącym do jednego koordynatora obsługiwanego przez rejestr.
type torStrumienia struct {
	jawne         []WpisStrumienia
	zwiniete      Zwiniete
	wpisow        int
	odcinek       hash.Hash64
	odcinekWpisow int
}

// Metoda dopisz dokłada wpis do toru i zwija najstarszy, gdy część jawna przekroczy ustaloną pojemność.
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

// Metoda zwin przenosi jeden wpis z części jawnej toru do podsumowania liczbowego jego zwiniętej historii.
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
