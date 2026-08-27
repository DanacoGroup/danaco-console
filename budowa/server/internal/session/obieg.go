package session

import "time"

// PowodZatrzymania nazywa przyczynę wstrzymania biegu naprawczego i rozróżnia pracę przerwaną od pracy ukończonej wynikiem, co pozwala rozstrzygnąć sposób wznowienia biegu.
type PowodZatrzymania string

const (
	// ZatrzymanieBrakPostepu oznacza, że wykonawca powtórzył ten sam wynik przez tyle obiegów z rzędu, ile wynosi próg braku postępu ustawiony dla biegu.
	ZatrzymanieBrakPostepu PowodZatrzymania = "brak postępu wykonawcy"
	// ZatrzymanieRecznie oznacza zatrzymanie biegu wydane bezpośrednio przez operatora, dostępne w dowolnej chwili trwania biegu naprawczego.
	ZatrzymanieRecznie PowodZatrzymania = "zatrzymanie przez Operatora"
	// ZatrzymanieUsterka oznacza, że rozpoczęcie obiegu zakończyło się niepowodzeniem i bieg naprawczy nie mógł zostać podjęty.
	ZatrzymanieUsterka PowodZatrzymania = "usterka rozpoczęcia obiegu"
	// ZatrzymanieUkonczenie oznacza, że koordynator zamknął turę wynikiem końcowym, gdy żadne okno wykonawcze nie prowadziło już tej tury, i jako jedyny powód oznacza pracę skończoną, a nie przerwaną.
	ZatrzymanieUkonczenie PowodZatrzymania = "ukończenie z wynikiem"
)

// ProgBrakuPostepuDomyslny podaje liczbę kolejnych obiegów bez zmiany stanu, po której bieg naprawczy zostaje zatrzymany domyślnie.
const ProgBrakuPostepuDomyslny = 3

// StanObiegu przechowuje odpis licznika obiegów prowadzonych przez jednego koordynatora w ramach biegu naprawczego.
type StanObiegu struct {
	// IdKoordynatora — okno prowadzące bieg naprawczy.
	IdKoordynatora string
	// Obiegow — łączna liczba obiegów pętli od założenia licznika.
	Obiegow int
	// ObiegowBezPostepu — ile ostatnich obiegów nie przyniosło zmiany stanu.
	ObiegowBezPostepu int
	// Prog — próg braku postępu obowiązujący ten bieg.
	Prog int
	// Zatrzymany mówi, czy bieg stanął; powód niżej rozstrzyga przerwanie od ukończenia.
	Zatrzymany bool
	// Powod zatrzymania; pusty, dopóki bieg trwa.
	Powod PowodZatrzymania
	// OstatniWykonawca — okno, którego tura wywołała ostatni obieg.
	OstatniWykonawca string
	// OstatniPowodTury — powód zakończenia tury wykonawcy.
	OstatniPowodTury string
	// Zaktualizowano — chwila ostatniej zmiany licznika.
	Zaktualizowano time.Time
}

// licznikObiegow prowadzi licznik obiegów oraz rozpoznaje warunek zatrzymania biegu naprawczego jednego koordynatora.
type licznikObiegow struct {
	stan          StanObiegu
	ostatniOdcisk string
}

// nowyLicznikObiegow zakłada licznik obiegów dla podanego koordynatora z określonym progiem braku postępu, przyjmując próg domyślny, gdy podana wartość jest niedodatnia.
func nowyLicznikObiegow(idKoordynatora string, prog int) *licznikObiegow {
	if prog <= 0 {
		prog = ProgBrakuPostepuDomyslny
	}
	return &licznikObiegow{stan: StanObiegu{
		IdKoordynatora: idKoordynatora,
		Prog:           prog,
		Zaktualizowano: time.Now().UTC(),
	}}
}

// zanotuj zapisuje kolejny obieg pętli, aktualizuje licznik braku postępu na podstawie odcisku stanu i zwraca stan po zapisie wraz z zgodą na rozpoczęcie następnego obiegu.
func (l *licznikObiegow) zanotuj(idWykonawcy, powodTury, odcisk string) (StanObiegu, bool) {
	if l.stan.Zatrzymany && l.stan.Powod == ZatrzymanieUkonczenie {
		l.stan.Zatrzymany = false
		l.stan.Powod = ""
	}
	if l.stan.Zatrzymany {
		return l.stan, false
	}
	l.stan.Obiegow++
	l.stan.OstatniWykonawca = idWykonawcy
	l.stan.OstatniPowodTury = powodTury
	l.stan.Zaktualizowano = time.Now().UTC()

	if odcisk != "" && odcisk == l.ostatniOdcisk {
		l.stan.ObiegowBezPostepu++
	} else {
		l.stan.ObiegowBezPostepu = 0
	}
	l.ostatniOdcisk = odcisk

	if l.stan.ObiegowBezPostepu >= l.stan.Prog {
		l.zatrzymaj(ZatrzymanieBrakPostepu)
		return l.stan, false
	}
	return l.stan, true
}

// zatrzymaj wstrzymuje bieg naprawczy, zapisując rozpoznany powód zatrzymania oraz chwilę tej zmiany w stanie licznika.
func (l *licznikObiegow) zatrzymaj(powod PowodZatrzymania) StanObiegu {
	l.stan.Zatrzymany = true
	l.stan.Powod = powod
	l.stan.Zaktualizowano = time.Now().UTC()
	return l.stan
}

// wznow podejmuje bieg naprawczy wstrzymany wcześniej, zerując licznik braku postępu i ostatni zapamiętany odcisk stanu.
func (l *licznikObiegow) wznow() StanObiegu {
	l.stan.Zatrzymany = false
	l.stan.Powod = ""
	l.stan.ObiegowBezPostepu = 0
	l.stan.Zaktualizowano = time.Now().UTC()
	l.ostatniOdcisk = ""
	return l.stan
}

// Obieg opisuje jeden obieg pętli naprawczej przekazywany oknu, które ma podjąć kolejną turę pracy koordynatora biegu.
type Obieg struct {
	// IdKoordynatora — okno, które ma podjąć pracę.
	IdKoordynatora string
	// IdWykonawcy — okno, którego tura zamknęła poprzedni obieg.
	IdWykonawcy string
	// Numer obiegu w bieżącym biegu naprawczym, licząc od jedynki.
	Numer int
	// PowodTury — powód zakończenia tury wykonawcy.
	PowodTury string
	// Strumien — odpis strumienia wykonawców dla koordynatora.
	Strumien MigawkaStrumienia
}

// UruchomienieObiegu jest portem rozpoczynającym turę koordynatora, oddzielającym decyzję o starcie obiegu od sposobu jej wykonania.
type UruchomienieObiegu interface {
	RozpocznijObieg(o Obieg) error
}

// UruchomienieFunkcja pozwala podać implementację portu UruchomienieObiegu zwykłą funkcją zamiast osobnego typu.
type UruchomienieFunkcja func(o Obieg) error

// RozpocznijObieg wypełnia interfejs UruchomienieObiegu, wywołując funkcję przekazaną jako implementacja portu.
func (f UruchomienieFunkcja) RozpocznijObieg(o Obieg) error { return f(o) }

// ZdarzenieObiegu niesie jawny ślad jednego obiegu pętli naprawczej, w tym również obiegu odmówionego z powodu zatrzymania.
type ZdarzenieObiegu struct {
	Stan       StanObiegu
	Obieg      Obieg
	Rozpoczety bool
	Blad       error
}

// ObserwatorObiegu jest portem przyjmującym ślad obiegu, wykorzystywanym przez dziennik zdarzeń, nadajnik powiadomień oraz telemetrię.
type ObserwatorObiegu interface {
	Obieg(z ZdarzenieObiegu)
}

// ObserwatorFunkcja pozwala podać implementację portu ObserwatorObiegu zwykłą funkcją zamiast osobnego typu.
type ObserwatorFunkcja func(z ZdarzenieObiegu)

// Obieg wypełnia interfejs ObserwatorObiegu, przekazując zdarzenie obiegu do funkcji pełniącej rolę obserwatora.
func (f ObserwatorFunkcja) Obieg(z ZdarzenieObiegu) { f(z) }
