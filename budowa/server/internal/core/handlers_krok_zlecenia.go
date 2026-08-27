// Plik wpina rodzinę queue.step.* — wstrzymanie kroku, decyzję o nim
// i wznowienie procesu z zastosowaną decyzją, rejestrując uchwyt tylko dla
// nazwy znanej kontraktowi.
package core

import (
	"context"

	"danacoconsole/shared"
)

// Nazwy komend rodziny sterowania krokiem queue.step.*, wyszukiwane
// w kontrakcie przy rejestracji uchwytów obsługi.
const (
	nazwaKomendyWstrzymaniaKroku shared.MessageType = "queue.step.hold"
	nazwaKomendyDecyzjiOKroku    shared.MessageType = "queue.step.decide"
	nazwaKomendyWykazuKrokow     shared.MessageType = "queue.step.list"
)

// krokZlecenia to krok widziany z zewnątrz: stan pracy i stan sterowania
// obok siebie, jako dwa odrębne fakty o kroku.
type krokZlecenia struct {
	// Identyfikator kroku (wiersz pozycji kolejki)
	Id string `json:"id"`
	// Kolejka, do ktorej krok nalezy
	QueueId string `json:"queueId"`
	// Etykieta kroku
	Title string `json:"title"`
	// Miejsce kroku w kolejnosci wykonania
	Order int `json:"order"`
	// Stan PRACY kroku: oczekuje, przydzielona, wykonywana, do_weryfikacji,
	// ukonczona, bledna, anulowana
	WorkStatus string `json:"workStatus"`
	// Stan STEROWANIA krokiem: czeka, biegnie, wstrzymany, zatwierdzony,
	// odrzucony, zamkniety
	ControlStatus string `json:"controlStatus"`
	// Licznik obiegow naprawczych kroku; bez limitu
	Cycle int `json:"cycle"`
	// Powod wstrzymania podany przy zatrzymaniu kroku
	HoldReason *string `json:"holdReason,omitempty"`
	// Uzasadnienie decyzji Operatora — jedzie do wykonawcy przy wznowieniu
	DecisionNote *string `json:"decisionNote,omitempty"`
	// Stan pracy, z ktorego krok zostanie wznowiony (nie od poczatku)
	ResumeFrom *string `json:"resumeFrom,omitempty"`
	// Chwila wstrzymania kroku
	HeldAt *string `json:"heldAt,omitempty"`
	// Chwila decyzji Operatora
	DecidedAt *string `json:"decidedAt,omitempty"`
	// Chwila zastosowania decyzji do losu kroku
	AppliedAt *string `json:"appliedAt,omitempty"`
	// Chwila doreczenia decyzji wykonawcy; puste znaczy: jeszcze nie dojechala
	DeliveredAt *string `json:"deliveredAt,omitempty"`
	// Ostatnia zmiana kroku
	UpdatedAt int64 `json:"updatedAt"`
}

// zadanieWstrzymaniaKroku niesie ładunek komendy queue.step.hold: kolejkę,
// wstrzymywany krok i powód wstrzymania.
type zadanieWstrzymaniaKroku struct {
	// Kolejka, w ktorej stoi krok
	QueueId string `json:"queueId"`
	// Krok wstrzymywany
	StepId string `json:"stepId"`
	// Powod wstrzymania
	Reason string `json:"reason,omitempty"`
}

// odpowiedzKroku niesie odpowiedź komendy queue.step.hold: krok po zmianie
// stanu oraz kolejkę, do której należy.
type odpowiedzKroku struct {
	Step  krokZlecenia `json:"step"`
	Queue shared.Queue `json:"queue"`
}

// zadanieDecyzjiKroku niesie ładunek komendy queue.step.decide: kolejkę,
// rozstrzygany krok, decyzję i uzasadnienie.
type zadanieDecyzjiKroku struct {
	// Kolejka, w ktorej stoi krok
	QueueId string `json:"queueId"`
	// Krok rozstrzygany
	StepId string `json:"stepId"`
	// Decyzja Operatora: approve albo reject
	Decision string `json:"decision"`
	// Uzasadnienie decyzji — dojezdza do wykonawcy razem ze zleceniem
	Note string `json:"note,omitempty"`
}

// odpowiedzDecyzjiKroku niesie odpowiedź komendy queue.step.decide wraz
// z faktem doręczenia decyzji wykonawcy kroku.
type odpowiedzDecyzjiKroku struct {
	Step  krokZlecenia `json:"step"`
	Queue shared.Queue `json:"queue"`
	// Czy decyzja dojechala do wykonawcy kroku
	Delivered bool `json:"delivered"`
	// Dlaczego nie dojechala, gdy nie dojechala
	DeliveryNote string `json:"deliveryNote,omitempty"`
}

// zadanieWykazuKrokow niesie ładunek komendy queue.step.list: identyfikator
// kolejki, której kroki mają być wypisane.
type zadanieWykazuKrokow struct {
	// Kolejka, ktorej kroki maja byc wypisane
	QueueId string `json:"queueId"`
}

// odpowiedzWykazuKrokow niesie odpowiedź komendy queue.step.list: wykaz
// kroków kolejki wskazanej w żądaniu.
type odpowiedzWykazuKrokow struct {
	Steps []krokZlecenia `json:"steps"`
}

// sterowanieKrokiem jest rozszerzeniem portu kolejek o sterowanie pojedynczym
// krokiem, a nie drugim portem: silnik zleceń jest jeden.
type sterowanieKrokiem interface {
	Kolejki

	// WstrzymajKrok obsługuje `queue.step.hold`.
	WstrzymajKrok(ctx context.Context, z zadanieWstrzymaniaKroku) (odpowiedzKroku, error)
	// ZdecydujOKroku obsługuje `queue.step.decide`.
	ZdecydujOKroku(ctx context.Context, z zadanieDecyzjiKroku) (odpowiedzDecyzjiKroku, error)
	// WykazKrokow obsługuje `queue.step.list`.
	WykazKrokow(ctx context.Context, z zadanieWykazuKrokow) (odpowiedzWykazuKrokow, error)
}

// zarejestrujSterowanieKrokiem wpina rodzinę queue.step.*; komendy
// wstrzymania i decyzji rozgłaszają zdarzenie zmiany kolejki, a wykaz kroków
// nie rozgłasza niczego.
func zarejestrujSterowanieKrokiem(r *Rejestr, kolejki Kolejki, e *emiter) {
	if r == nil || kolejki == nil {
		return
	}
	sterowanie, ok := kolejki.(sterowanieKrokiem)
	if !ok {
		// Port kolejek bez sterowania krokiem zostawia rodzinę queue.step.*
		// nieznaną.
		return
	}

	zarejestrujKomendeKontraktu(r, nazwaKomendyWstrzymaniaKroku,
		obsluz(func(ctx context.Context, z zadanieWstrzymaniaKroku) (odpowiedzKroku, error) {
			w, err := sterowanie.WstrzymajKrok(ctx, z)
			if err == nil {
				e.kolejka(shared.ChangeKindUpdated, w.Queue)
			}
			return w, err
		}))

	zarejestrujKomendeKontraktu(r, nazwaKomendyDecyzjiOKroku,
		obsluz(func(ctx context.Context, z zadanieDecyzjiKroku) (odpowiedzDecyzjiKroku, error) {
			w, err := sterowanie.ZdecydujOKroku(ctx, z)
			if err == nil {
				e.kolejka(shared.ChangeKindUpdated, w.Queue)
			}
			return w, err
		}))

	zarejestrujKomendeKontraktu(r, nazwaKomendyWykazuKrokow, obsluz(sterowanie.WykazKrokow))
}

// zarejestrujKomendeKontraktu wpina uchwyt wyłącznie pod nazwą, którą zna
// kontrakt. Nazwa spoza kontraktu nie wpina niczego — rdzeń nie ma prawa
// ogłaszać zdolności, której kontrakt nie opisuje.
func zarejestrujKomendeKontraktu(r *Rejestr, nazwa shared.MessageType, obsluga Obsluga) {
	if !shared.CzyKomenda(nazwa) {
		return
	}
	r.Zarejestruj(nazwa, obsluga)
}
