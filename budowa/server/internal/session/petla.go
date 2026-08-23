package session

import (
	"errors"
	"sync"

	"danacoconsole/server/internal/protocol"
)

// Pętla koordynator–wykonawca w postaci jawnego, testowalnego bytu.
//
// Petla jest jedynym odbiorcą wybudzeń, który prowadzi bieg naprawczy, i domyka
// łańcuch: koniec tury wykonawcy → wybudzenie → kolejny obieg koordynatora. Nie
// powiela silnika uruchamiania tury — samą turę koordynatora rozpoczyna warstwa
// rozmowy przez port UruchomienieObiegu.
//
// Trzy obowiązki w jednym miejscu:
//   - koniec tury → wybudzenie → kolejny obieg koordynatora;
//   - licznik obiegów, brak postępu, nazwany warunek zatrzymania;
//   - strumień wykonawcy widziany przez koordynatora.

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

// Petla prowadzi bieg naprawczy każdego koordynatora z osobna.
type Petla struct {
	rejestr      *Rejestr
	wybudzacz    *Wybudzacz
	strumien     *StrumienWykonawcy
	uruchomienie UruchomienieObiegu
	prog         int

	mu           sync.Mutex
	liczniki     map[string]*licznikObiegow
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
	}
	n.Wybudzacz().Zarejestruj(p)
	return p
}

// Strumien udostępnia strumień wykonawców.
func (p *Petla) Strumien() *StrumienWykonawcy { return p.strumien }

// Obserwuj dokłada odbiorcę śladu obiegów.
func (p *Petla) Obserwuj(o ObserwatorObiegu) {
	if o == nil {
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.obserwatorzy = append(p.obserwatorzy, o)
}

// ObserwujFragment zapisuje fragment okna w strumieniu jego koordynatora.
// Fragment okna spoza pętli nie ma dokąd trafić i nie jest to usterka.
func (p *Petla) ObserwujFragment(f protocol.Chunk) {
	okno, err := p.rejestr.Okno(f.WindowId)
	if err != nil || !okno.CzyWykonawca() || okno.OknoKoordynatora == "" {
		return
	}
	p.strumien.DopiszFragment(okno.OknoKoordynatora, f)
}

// ZakonczTure zgłasza koniec tury okna. Jest to jedyne wejście warstwy rozmowy
// do pętli: stąd idzie wybudzenie koordynatora, a z niego kolejny obieg.
// Zwraca prawdę, gdy zgłoszenie znalazło koordynatora.
func (p *Petla) ZakonczTure(idOkna, powod string) bool {
	return p.wybudzacz.KoniecTury(idOkna, powod)
}

// Wybudz wypełnia interfejs OdbiorcaWybudzenia — tu zamyka się pętla.
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

// Stan zwraca odpis licznika koordynatora.
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

// Zapomnij usuwa ślad zamkniętego koordynatora — licznik i jego strumień.
func (p *Petla) Zapomnij(idKoordynatora string) {
	p.mu.Lock()
	delete(p.liczniki, idKoordynatora)
	p.mu.Unlock()
	p.strumien.Zapomnij(idKoordynatora)
}

// rozpocznij oddaje obieg portowi warstwy rozmowy.
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
