package session

import (
	"sync"
	"time"
)

// Zakończenie tury wykonawcy wybudza koordynatora; to jedyne miejsce, w którym reguła ta żyje.

// PowodWynik — tura zamknęła się wynikiem kanału. Powód zakończenia tury nazywa
// warstwa rozmowy, bo to ona prowadzi turę; sesja podaje wspólną nazwę, żeby
// koordynator i dziennik czytały tę samą wartość.
const PowodWynik = "wynik tury"

// Wybudzenie opisuje zdarzenie zakończenia tury wykonawcy, niosące powód, oba okna oraz czas zgłoszenia.
type Wybudzenie struct {
	// OknoWykonawcy — okno, którego tura się zakończyła.
	OknoWykonawcy string
	// OknoKoordynatora — okno, które ma zostać wybudzone.
	OknoKoordynatora string
	// Powod zakończenia tury: znacznik kanału albo zakończenie procesu.
	Powod string
	// Czas zgłoszenia.
	Czas time.Time
}

// OdbiorcaWybudzenia przyjmuje wybudzenie koordynatora. Implementuje go
// warstwa kolejek rdzenia; pakiet session zna wyłącznie ten interfejs.
type OdbiorcaWybudzenia interface {
	Wybudz(w Wybudzenie)
}

// OdbiorcaFunkcja pozwala podać odbiorcę wybudzeń zwykłą funkcją zamiast osobnej struktury implementującej interfejs.
type OdbiorcaFunkcja func(w Wybudzenie)

// Metoda Wybudz wypełnia interfejs OdbiorcaWybudzenia, przekazując wybudzenie owiniętej funkcji odbiorcy.
func (f OdbiorcaFunkcja) Wybudz(w Wybudzenie) { f(w) }

// Wybudzacz kieruje zgłoszenia końca tury do właściwych koordynatorów, korzystając z rejestru okien tej sesji.
type Wybudzacz struct {
	rejestr  *Rejestr
	mu       sync.RWMutex
	odbiorcy []OdbiorcaWybudzenia
}

// Funkcja NowyWybudzacz składa nowy wybudzacz nad podanym rejestrem okien, gotowy do przyjmowania zgłoszeń.
func NowyWybudzacz(rejestr *Rejestr) *Wybudzacz {
	return &Wybudzacz{rejestr: rejestr}
}

// Metoda Zarejestruj dokłada do wybudzacza odbiorcę wybudzeń, który otrzyma kolejne zdarzenia końca tury.
func (w *Wybudzacz) Zarejestruj(o OdbiorcaWybudzenia) {
	if o == nil {
		return
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	w.odbiorcy = append(w.odbiorcy, o)
}

// KoniecTury przyjmuje zgłoszenie od procesu okna. Zwraca prawdę, gdy
// wybudzenie zostało skierowane do koordynatora.
//
// Okno samodzielne, okno koordynatora oraz wykonawca bez wskazanego
// koordynatora nie wybudzają nikogo — i nie jest to usterka.
func (w *Wybudzacz) KoniecTury(idOkna, powod string) bool {
	okno, err := w.rejestr.Okno(idOkna)
	if err != nil || !okno.CzyWykonawca() || okno.OknoKoordynatora == "" {
		return false
	}
	if _, err := w.rejestr.Okno(okno.OknoKoordynatora); err != nil {
		return false
	}
	w.rozeslij(Wybudzenie{
		OknoWykonawcy:    okno.Id,
		OknoKoordynatora: okno.OknoKoordynatora,
		Powod:            powod,
		Czas:             time.Now().UTC(),
	})
	return true
}

// rozeslij przekazuje wybudzenie wszystkim odbiorcom. Odbiorcy pracują na
// odpisie zdarzenia i nie wpływają na siebie nawzajem.
func (w *Wybudzacz) rozeslij(zdarzenie Wybudzenie) {
	w.mu.RLock()
	odbiorcy := append([]OdbiorcaWybudzenia(nil), w.odbiorcy...)
	w.mu.RUnlock()
	for _, odbiorca := range odbiorcy {
		odbiorca.Wybudz(zdarzenie)
	}
}
