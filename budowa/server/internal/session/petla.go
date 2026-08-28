package session

import (
	"errors"
	"sync"

	"danacoconsole/server/internal/protocol"
)

// Petla koordynator-wykonawca jest jawnym, testowalnym bytem domykającym pełny łańcuch obiegu.

// ErrBrakUruchomieniaObiegu — pętla nie dostała portu rozpoczynania obiegu.
// Wybudzenie zostaje wtedy policzone i rozgłoszone, ale bieg staje z powodem
// usterki; milczące zignorowanie wybudzenia byłoby zatrzymaniem cichym.
var ErrBrakUruchomieniaObiegu = errors.New("session: pętla nie ma portu rozpoczynania obiegu")

// UstawieniaPetli zbierają nastawy biegu naprawczego. Wartości niedodatnie
// schodzą na domyślne — brak nastawy nie blokuje pętli.
type UstawieniaPetli struct {
	// ProgBrakuPostepu — ile obiegów bez zmiany stanu kończy bieg.
	ProgBrakuPostepu int
	// PojemnoscStrumienia — ile ostatnich wpisów wykonawcy zostaje jawnych.
	PojemnoscStrumienia int
}

// Struktura Petla prowadzi bieg naprawczy każdego koordynatora z osobna, utrzymując licznik obiegów, strumień wykonawcy oraz wykaz otwartych tur.
type Petla struct {
	rejestr      *Rejestr
	wybudzacz    *Wybudzacz
	strumien     *StrumienWykonawcy
	uruchomienie UruchomienieObiegu
	prog         int

	mu       sync.Mutex
	liczniki map[string]*licznikObiegow
	// tury notują okna wykonawcze prowadzące turę w torach koordynatorów; wpis i zdjęcie idą parą.
	tury         map[string]map[string]struct{}
	obserwatorzy []ObserwatorObiegu
}

// NowaPetla składa pętlę nad nadzorcą i wpina ją w wybudzacz. Od tej chwili
// koniec tury wykonawcy nie kończy się na wybudzeniu — rozpoczyna kolejny obieg.
func NowaPetla(n *Nadzorca, u UruchomienieObiegu, ust UstawieniaPetli) *Petla {
	p := &Petla{
		rejestr:      n.Rejestr(),
		wybudzacz:    n.Wybudzacz(),
		strumien:     NowyStrumienWykonawcy(ust.PojemnoscStrumienia),
		uruchomienie: u,
		prog:         ust.ProgBrakuPostepu,
		liczniki:     map[string]*licznikObiegow{},
		tury:         map[string]map[string]struct{}{},
	}
	n.Wybudzacz().Zarejestruj(p)
	return p
}

// Metoda Strumien udostępnia strumień wykonawców powiązany z pętlą, wykorzystywany do obserwacji przebiegu kolejnych obiegów koordynatora.
func (p *Petla) Strumien() *StrumienWykonawcy { return p.strumien }

// Metoda Obserwuj dokłada do pętli odbiorcę śladu obiegów, który otrzymuje kolejne zdarzenia rozgłaszane przy zamknięciu każdego obiegu.
func (p *Petla) Obserwuj(o ObserwatorObiegu) {
	if o == nil {
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.obserwatorzy = append(p.obserwatorzy, o)
}

// ObserwujFragment zapisuje fragment okna w strumieniu jego koordynatora
// i notuje turę wykonawcy, z której fragment przyszedł. Fragment okna spoza
// pętli nie ma dokąd trafić i nie jest to usterka.
func (p *Petla) ObserwujFragment(f protocol.Chunk) {
	okno, err := p.rejestr.Okno(f.WindowId)
	if err != nil || !okno.CzyWykonawca() || okno.OknoKoordynatora == "" {
		return
	}
	p.otworzTure(okno.OknoKoordynatora, okno.Id)
	p.strumien.DopiszFragment(okno.OknoKoordynatora, f)
}

// Metoda ZakonczTure zgłasza koniec tury okna i jest jedynym wejściem warstwy rozmowy do pętli dla każdego zakończenia tury.
func (p *Petla) ZakonczTure(idOkna, powod string) bool {
	okno, err := p.rejestr.Okno(idOkna)
	if err != nil {
		return false
	}
	if okno.CzyKoordynator() {
		return p.ukonczBieg(okno.Id, powod)
	}
	if okno.CzyWykonawca() && okno.OknoKoordynatora != "" {
		p.zamknijTure(okno.OknoKoordynatora, okno.Id)
	}
	return p.wybudzacz.KoniecTury(idOkna, powod)
}

// Metoda Wybudz wypełnia interfejs OdbiorcaWybudzenia i w tym miejscu zamyka się pętla koordynator-wykonawca, rozpoczynając kolejny obieg.
func (p *Petla) Wybudz(w Wybudzenie) {
	odcisk := p.strumien.Odetnij(w.OknoKoordynatora)

	p.mu.Lock()
	licznik := p.licznik(w.OknoKoordynatora)
	stan, wolno := licznik.zanotuj(w.OknoWykonawcy, w.Powod, odcisk)
	p.mu.Unlock()

	obieg := Obieg{
		IdKoordynatora: w.OknoKoordynatora,
		IdWykonawcy:    w.OknoWykonawcy,
		Numer:          stan.Obiegow,
		PowodTury:      w.Powod,
		Strumien:       p.strumien.Migawka(w.OknoKoordynatora),
	}
	if !wolno {
		p.rozglos(ZdarzenieObiegu{Stan: stan, Obieg: obieg})
		return
	}
	if err := p.rozpocznij(obieg); err != nil {
		p.mu.Lock()
		stan = licznik.zatrzymaj(ZatrzymanieUsterka)
		p.mu.Unlock()
		p.rozglos(ZdarzenieObiegu{Stan: stan, Obieg: obieg, Blad: err})
		return
	}
	p.rozglos(ZdarzenieObiegu{Stan: stan, Obieg: obieg, Rozpoczety: true})
}

// Metoda Stan zwraca odpis bieżącego licznika obiegów koordynatora wskazanego identyfikatorem, bez modyfikacji jego stanu wewnętrznego.
func (p *Petla) Stan(idKoordynatora string) StanObiegu {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.licznik(idKoordynatora).stan
}

// Zatrzymaj wstrzymuje bieg naprawczy koordynatora. Przycisk zatrzymania jest
// czynny zawsze, także przed pierwszym obiegiem.
func (p *Petla) Zatrzymaj(idKoordynatora string) StanObiegu {
	p.mu.Lock()
	stan := p.licznik(idKoordynatora).zatrzymaj(ZatrzymanieRecznie)
	p.mu.Unlock()
	p.rozglos(ZdarzenieObiegu{Stan: stan, Obieg: Obieg{IdKoordynatora: idKoordynatora, Numer: stan.Obiegow}})
	return stan
}

// Wznow podejmuje bieg wstrzymany. Decyzja o wznowieniu należy do Operatora —
// pętla nie wznawia się sama.
func (p *Petla) Wznow(idKoordynatora string) StanObiegu {
	p.mu.Lock()
	stan := p.licznik(idKoordynatora).wznow()
	p.mu.Unlock()
	p.rozglos(ZdarzenieObiegu{Stan: stan, Obieg: Obieg{IdKoordynatora: idKoordynatora, Numer: stan.Obiegow}})
	return stan
}

// Metoda Zapomnij usuwa ślad zamkniętego koordynatora, czyli jego licznik, wykaz tur jego wykonawców oraz jego strumień zdarzeń.
func (p *Petla) Zapomnij(idKoordynatora string) {
	p.mu.Lock()
	delete(p.liczniki, idKoordynatora)
	delete(p.tury, idKoordynatora)
	p.mu.Unlock()
	p.strumien.Zapomnij(idKoordynatora)
}

// Metoda ukonczBieg zamyka bieg koordynatora ukończeniem z wynikiem, gdy spełniony jest podwójny warunek mierzony przez pętlę.
func (p *Petla) ukonczBieg(idKoordynatora, powod string) bool {
	if powod != PowodWynik {
		return false
	}
	p.mu.Lock()
	licznik, prowadzony := p.liczniki[idKoordynatora]
	if !prowadzony || licznik.stan.Obiegow == 0 || licznik.stan.Zatrzymany ||
		len(p.tury[idKoordynatora]) > 0 {

		p.mu.Unlock()
		return false
	}
	stan := licznik.zatrzymaj(ZatrzymanieUkonczenie)
	p.mu.Unlock()

	p.rozglos(ZdarzenieObiegu{Stan: stan,
		Obieg: Obieg{IdKoordynatora: idKoordynatora, Numer: stan.Obiegow}})
	return true
}

// Metoda otworzTure notuje w torze koordynatora turę okna wykonawczego, oznaczając je jako prowadzące bieżącą turę obiegu.
func (p *Petla) otworzTure(idKoordynatora, idWykonawcy string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	okna, jest := p.tury[idKoordynatora]
	if !jest {
		okna = map[string]struct{}{}
		p.tury[idKoordynatora] = okna
	}
	okna[idWykonawcy] = struct{}{}
}

// Metoda zamknijTure zdejmuje z toru koordynatora turę okna wykonawczego, oznaczając zakończenie jego udziału w obiegu.
func (p *Petla) zamknijTure(idKoordynatora, idWykonawcy string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	okna, jest := p.tury[idKoordynatora]
	if !jest {
		return
	}
	delete(okna, idWykonawcy)
	if len(okna) == 0 {
		delete(p.tury, idKoordynatora)
	}
}

// Metoda rozpocznij oddaje obieg portowi warstwy rozmowy, który odpowiada za faktyczne uruchomienie kolejnej tury koordynatora.
func (p *Petla) rozpocznij(o Obieg) error {
	if p.uruchomienie == nil {
		return ErrBrakUruchomieniaObiegu
	}
	return p.uruchomienie.RozpocznijObieg(o)
}

// licznik zwraca licznik koordynatora, zakładając go przy pierwszym obiegu.
// Wywoływane pod założoną blokadą.
func (p *Petla) licznik(idKoordynatora string) *licznikObiegow {
	l, jest := p.liczniki[idKoordynatora]
	if !jest {
		l = nowyLicznikObiegow(idKoordynatora, p.prog)
		p.liczniki[idKoordynatora] = l
	}
	return l
}

// rozglos przekazuje ślad obiegu obserwatorom. Obserwatorzy pracują na odpisie
// zdarzenia i nie wpływają na siebie nawzajem.
func (p *Petla) rozglos(z ZdarzenieObiegu) {
	p.mu.Lock()
	obserwatorzy := append([]ObserwatorObiegu(nil), p.obserwatorzy...)
	p.mu.Unlock()
	for _, o := range obserwatorzy {
		o.Obieg(z)
	}
}
