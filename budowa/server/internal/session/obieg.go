package session

import "time"

// Bieg naprawczy nie ma limitu obiegów. Zamiast bramy licznikowej
// wchodzi przejrzystość: licznik obiegów, wykrywanie braku postępu i jawny,
// nazwany warunek zatrzymania. Zatrzymanie nigdy nie jest ciche — zawsze niesie
// rozpoznany powód, który idzie do obserwatorów pętli.

// PowodZatrzymania nazywa przyczynę wstrzymania biegu naprawczego.
type PowodZatrzymania string

const (
	// ZatrzymanieBrakPostepu — wykonawca powtórzył się co do znaku przez tyle
	// obiegów z rzędu, ile wynosi próg braku postępu.
	ZatrzymanieBrakPostepu PowodZatrzymania = "brak postępu wykonawcy"
	// ZatrzymanieRecznie — zatrzymanie przez Operatora. Przycisk zatrzymania
	// jest czynny w każdej chwili biegu.
	ZatrzymanieRecznie PowodZatrzymania = "zatrzymanie przez Operatora"
	// ZatrzymanieUsterka — obiegu nie udało się rozpocząć.
	ZatrzymanieUsterka PowodZatrzymania = "usterka rozpoczęcia obiegu"
)

// ProgBrakuPostepuDomyslny — ile obiegów bez zmiany stanu z rzędu kończy bieg.
const ProgBrakuPostepuDomyslny = 3

// StanObiegu jest odpisem licznika jednego koordynatora.
type StanObiegu struct {
	// IdKoordynatora — okno prowadzące bieg naprawczy.
	IdKoordynatora string
	// Obiegow — łączna liczba obiegów pętli od założenia licznika.
	Obiegow int
	// ObiegowBezPostepu — ile ostatnich obiegów nie przyniosło zmiany stanu.
	ObiegowBezPostepu int
	// Prog — próg braku postępu obowiązujący ten bieg.
	Prog int
	// Zatrzymany mówi, czy bieg został wstrzymany.
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

// licznikObiegow prowadzi bieg naprawczy jednego koordynatora.
type licznikObiegow struct {
	stan          StanObiegu
	ostatniOdcisk string
}

// nowyLicznikObiegow zakłada licznik z podanym progiem braku postępu.
// Próg niedodatni schodzi na wartość domyślną.
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

// zanotuj zapisuje kolejny obieg i rozstrzyga, czy bieg trwa dalej.
// Zwraca stan po zapisie oraz zgodę na rozpoczęcie obiegu.
func (l *licznikObiegow) zanotuj(idWykonawcy, powodTury, odcisk string) (StanObiegu, bool) {
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

// zatrzymaj wstrzymuje bieg z rozpoznanym powodem.
func (l *licznikObiegow) zatrzymaj(powod PowodZatrzymania) StanObiegu {
	l.stan.Zatrzymany = true
	l.stan.Powod = powod
	l.stan.Zaktualizowano = time.Now().UTC()
	return l.stan
}

// wznow podejmuje bieg wstrzymany. Historia obiegów zostaje — kasuje się
// wyłącznie licznik braku postępu i ostatni odcisk, bo to one rozstrzygały
// o zatrzymaniu.
func (l *licznikObiegow) wznow() StanObiegu {
	l.stan.Zatrzymany = false
	l.stan.Powod = ""
	l.stan.ObiegowBezPostepu = 0
	l.stan.Zaktualizowano = time.Now().UTC()
	l.ostatniOdcisk = ""
	return l.stan
}

// Obieg opisuje jeden obieg pętli przekazywany temu, kto go wykonuje.
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

// UruchomienieObiegu rozpoczyna turę koordynatora. Port trzyma pętlę po jednej
// stronie granicy: session wie kiedy zacząć obieg, warstwa rozmowy wie jak.
type UruchomienieObiegu interface {
	RozpocznijObieg(o Obieg) error
}

// UruchomienieFunkcja pozwala podać port zwykłą funkcją.
type UruchomienieFunkcja func(o Obieg) error

// RozpocznijObieg wypełnia interfejs UruchomienieObiegu.
func (f UruchomienieFunkcja) RozpocznijObieg(o Obieg) error { return f(o) }

// ZdarzenieObiegu jest jawnym śladem jednego obiegu — także obiegu odmówionego.
type ZdarzenieObiegu struct {
	Stan       StanObiegu
	Obieg      Obieg
	Rozpoczety bool
	Blad       error
}

// ObserwatorObiegu przyjmuje ślad obiegu: dziennik, nadajnik zdarzeń, telemetria.
type ObserwatorObiegu interface {
	Obieg(z ZdarzenieObiegu)
}

// ObserwatorFunkcja pozwala podać obserwatora zwykłą funkcją.
type ObserwatorFunkcja func(z ZdarzenieObiegu)

// Obieg wypełnia interfejs ObserwatorObiegu.
func (f ObserwatorFunkcja) Obieg(z ZdarzenieObiegu) { f(z) }
