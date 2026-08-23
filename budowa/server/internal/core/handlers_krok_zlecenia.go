// Odpowiedzialność pliku: wpięcie rodziny `queue.step.*` — wstrzymania kroku,
// decyzji o nim i wznowienia procesu z zastosowaną decyzją.
//
// Kontrakt nie niesie dziś ani jednej komendy dotyczącej pozycji kolejki, ani
// struktury pozycji, ani wyliczenia jej stanów. Rejestracja niżej pyta
// `shared.CzyKomenda` o każdą nazwę i wpina uchwyt wyłącznie wtedy, gdy nazwa
// do kontraktu należy. Skutki są dwa i oba są zamierzone:
//
//   - dopóki kontrakt tych nazw nie zna, nie wpina się nic, a klient wołający
//     `queue.step.hold` dostaje uczciwe `queue.unknown`: rdzeń nie ogłasza
//     zdolności, której kontrakt nie opisuje;
//   - po wniesieniu komend do kontraktu uchwyty wpinają się bez dotykania tego
//     pliku, bo warunek jest tym samym pytaniem, które zadaje sprawdzian
//     zgodności rejestru z kontraktem.
//
// Literały nazw stoją tu wyjątkowo, bo są kluczem wyszukania w kontrakcie,
// a nie deklaracją komendy; po wniesieniu komend ustępują stałym
// `shared.CommandQueueStep*`. Struktury niżej mają wtedy stać się aliasami
// typów kontraktu, bez ruszania ciał uchwytów i adaptera.
package core

import (
	"context"

	"danacoconsole/shared"
)

// Nazwy komend rodziny sterowania krokiem, wyszukiwane w kontrakcie.
const (
	nazwaKomendyWstrzymaniaKroku shared.MessageType = "queue.step.hold"
	nazwaKomendyDecyzjiOKroku    shared.MessageType = "queue.step.decide"
	nazwaKomendyWykazuKrokow     shared.MessageType = "queue.step.list"
)

// krokZlecenia to krok widziany z zewnątrz: stan pracy i stan sterowania obok
// siebie. Dwa stany, nie jeden, bo są to dwa różne fakty — „gdzie krok stoi
// w pracy" i „czego oczekuje od człowieka". Sklejenie ich w jedno pole
// zmusiłoby do wyboru, który z nich zataić.
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

// zadanieWstrzymaniaKroku — ładunek `queue.step.hold`.
type zadanieWstrzymaniaKroku struct {
	// Kolejka, w ktorej stoi krok
	QueueId string `json:"queueId"`
	// Krok wstrzymywany
	StepId string `json:"stepId"`
	// Powod wstrzymania
	Reason string `json:"reason,omitempty"`
}

// odpowiedzKroku — odpowiedź `queue.step.hold`.
type odpowiedzKroku struct {
	Step  krokZlecenia `json:"step"`
	Queue shared.Queue `json:"queue"`
}

// zadanieDecyzjiKroku — ładunek `queue.step.decide`.
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

// odpowiedzDecyzjiKroku — odpowiedź `queue.step.decide`.
//
// Pola `delivered` i `deliveryNote` niosą fakt, którego wynik inaczej by nie
// pokazał: decyzja zastosowana do stanu kroku, ale niedoręczona wykonawcy.
// Bez nich brak efektu trzeba by odgadywać.
type odpowiedzDecyzjiKroku struct {
	Step  krokZlecenia `json:"step"`
	Queue shared.Queue `json:"queue"`
	// Czy decyzja dojechala do wykonawcy kroku
	Delivered bool `json:"delivered"`
	// Dlaczego nie dojechala, gdy nie dojechala
	DeliveryNote string `json:"deliveryNote,omitempty"`
}

// zadanieWykazuKrokow — ładunek `queue.step.list`.
type zadanieWykazuKrokow struct {
	// Kolejka, ktorej kroki maja byc wypisane
	QueueId string `json:"queueId"`
}

// odpowiedzWykazuKrokow — odpowiedź `queue.step.list`.
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

// zarejestrujSterowanieKrokiem wpina rodzinę `queue.step.*`.
//
// Zdarzenie `queue.changed` rozgłaszają obie komendy zmieniające, bo obie
// zmieniają kolejkę — wstrzymanie kroku bieżącego przestawia ją w `paused`,
// a decyzja podejmuje ją z powrotem. Wykaz kroków niczego nie zmienia i niczego
// nie rozgłasza: wykaz ogłoszony jako zmiana byłby zdarzeniem bez faktu.
func zarejestrujSterowanieKrokiem(r *Rejestr, kolejki Kolejki, e *emiter) {
	if r == nil || kolejki == nil {
		return
	}
	sterowanie, ok := kolejki.(sterowanieKrokiem)
	if !ok {
		// Port kolejek bez sterowania krokiem zostawia rodzinę nieznaną, tak
		// samo jak każdą domenę bez portu. Montaż nie odmawia: dopóki kontrakt
		// nazw nie niesie, odmowa dotyczyłaby komendy, której i tak nikt nie
		// może zawołać.
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
