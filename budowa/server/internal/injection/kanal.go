package injection

import (
	"context"
	"encoding/json"
	"fmt"

	"danacoconsole/shared"
)

// Kanal jest adapterem kanału głównego. Nie trzyma żadnych ustawień
// rozmowy — te przychodzą z każdym zapytaniem, bo rozstrzyga je resolver
// zasięgów. Trzyma wyłącznie pulę kont, bo rotacja jest
// stanem żyjącym dłużej niż jedno wywołanie.
type Kanal struct {
	pula *PulaKont
}

// NowyKanal składa kanał nad pulą kont.
func NowyKanal(pula *PulaKont) *Kanal {
	if pula == nil {
		pula = NowaPula()
	}
	return &Kanal{pula: pula}
}

// Pula udostępnia pulę kont — rejestr kanałów pokazuje jej stan Operatorowi.
func (k *Kanal) Pula() *PulaKont {
	return k.pula
}

// Rozmowa przeprowadza jedną turę i oddaje strumień fragmentów. Kanał nigdy
// nie zwraca błędu obok strumienia: każda przeszkoda jedzie fragmentem rodzaju
// „błąd" i domyka strumień, przez co jedna droga obsługuje i powodzenie,
// i niepowodzenie.
//
// Strumień zaczyna się fragmentem prowenancji, potem idą fragmenty treści,
// a kończy fragment ze znacznikiem ostatniego i podsumowaniem tury. To ono
// wybudza koordynatora w pętli koordynator–wykonawca.
func (k *Kanal) Rozmowa(kontekst context.Context, z Zapytanie) <-chan Fragment {
	strumien := make(chan Fragment)
	go func() {
		defer close(strumien)
		k.prowadz(kontekst, z, strumien)
	}()
	return strumien
}

// prowadz prowadzi turę razem z rotacją kont.
func (k *Kanal) prowadz(kontekst context.Context, z Zapytanie, na chan<- Fragment) {
	// Konto wskazane wygrywa z rotacją i jest rozkazem tożsamości:
	// tura jedzie dokładnie nim albo mówi wprost, dlaczego nie pojedzie.
	if z.Ustawienia.Konto != "" {
		k.prowadzWskazanym(kontekst, z, na)
		return
	}
	powod := ""
	for proba := 1; proba <= len(k.pula.Konta())+1; proba++ {
		konto, jest := k.pula.Biezace()
		if !jest {
			// Pula pusta ≠ pula wyczerpana. Brak kont znaczy, że
			// Operator nie wskazał tożsamości — tura idzie z tożsamością
			// otoczenia, czyli bez CLAUDE_CONFIG_DIR (proces.go pomija zmienną
			// przy pustym katalogu). Odmowa w tym miejscu zablokowałaby pierwszą
			// turę na świeżej instalacji.
			if !k.pula.Pusta() {
				zakonczBledem(kontekst, na, z, shared.ErrorCodeChannelUnavailable,
					"wszystkie konta puli mają wyczerpany limit; kolejne próby pozostają otwarte")
				return
			}
			konto = Konto{}
		}

		wynik, err := wykonajPrzebieg(kontekst, z, konto, proba, powod, na)
		if err != nil {
			zakonczBledem(kontekst, na, z, shared.ErrorCodeChannelUnavailable, err.Error())
			return
		}

		// Pułap sprawdza się przed wyczerpaniem i zatrzymuje rotację: kolejne
		// konto puli wydałoby dokładnie tę kwotę, której Operator wydać zabronił
		// (pulap.go).
		if pulap, naPulapie := rozpoznajPulap(wynik.Obserwacja, wynik.Bledy); naPulapie {
			zakonczPulapem(kontekst, na, z, pulap)
			return
		}

		wyczerpanie, wyczerpane := rozpoznajWyczerpanie(wynik.Obserwacja, wynik.Bledy)
		if !wyczerpane {
			zakoncz(kontekst, na, z, wynik.Obserwacja.Tura, konto.Kod)
			return
		}

		nastepne, dostepne := k.pula.Wyczerpane(konto.Kod, wyczerpanie.DoChwili)
		if !dostepne || nastepne.Kod == konto.Kod || wynik.Obserwacja.TekstPoszedl {
			zakoncz(kontekst, na, z, wynik.Obserwacja.Tura, konto.Kod)
			return
		}
		powod = fmt.Sprintf("rotacja konta %s → %s (%s)", konto.Kod, nastepne.Kod, wyczerpanie.Powod)
	}
	zakonczBledem(kontekst, na, z, shared.ErrorCodeChannelUnavailable,
		"pula kont wyczerpana w tej turze; sesja pozostaje czynna")
}

// prowadzWskazanym prowadzi turę na koncie wskazanym konfiguracją.
//
// Rotacji tu nie ma: wskazanie konta ustala tożsamość okna, a cicha podmiana
// na inne konto po wyczerpaniu limitu wykonałaby turę tożsamością, której
// Operator nie wybrał. Wyczerpanie zostaje odnotowane w puli (ślad i trwałość
// limitu), a tura kończy się tym, co konto oddało; konto spoza puli daje
// odmowę, nie inną tożsamość.
func (k *Kanal) prowadzWskazanym(kontekst context.Context, z Zapytanie, na chan<- Fragment) {
	konto, jest := k.pula.PoKodzie(z.Ustawienia.Konto)
	if !jest {
		zakonczBledem(kontekst, na, z, shared.ErrorCodeChannelUnavailable,
			fmt.Sprintf("wskazane konto %q nie stoi w puli kanału głównego; tura nie pojedzie inną tożsamością niż wskazana", z.Ustawienia.Konto))
		return
	}
	wynik, err := wykonajPrzebieg(kontekst, z, konto, 1, "konto wskazane konfiguracją", na)
	if err != nil {
		zakonczBledem(kontekst, na, z, shared.ErrorCodeChannelUnavailable, err.Error())
		return
	}
	if pulap, naPulapie := rozpoznajPulap(wynik.Obserwacja, wynik.Bledy); naPulapie {
		zakonczPulapem(kontekst, na, z, pulap)
		return
	}
	if wyczerpanie, wyczerpane := rozpoznajWyczerpanie(wynik.Obserwacja, wynik.Bledy); wyczerpane {
		// Ślad wyczerpania idzie do puli, ale przełączenia nie ma — patrz nagłówek.
		k.pula.Wyczerpane(konto.Kod, wyczerpanie.DoChwili)
	}
	zakoncz(kontekst, na, z, wynik.Obserwacja.Tura, konto.Kod)
}

// zakoncz wysyła fragment kończący turę. Gdy tura nie przyniosła podsumowania,
// zamiast niego jedzie fragment błędu — strumień zawsze ma koniec.
func zakoncz(kontekst context.Context, na chan<- Fragment, z Zapytanie, tura *ZakonczenieTury, konto string) {
	if tura == nil {
		zakonczBledem(kontekst, na, z, shared.ErrorCodeChannelUnavailable, "tura zakończyła się bez podsumowania")
		return
	}
	tura.Konto = konto
	fragment := Fragment{
		Chunk:   shared.StreamChunkEvent{WindowId: z.IdOkna, MessageId: z.IdWiadomosci, Kind: shared.ChunkKindText},
		Ostatni: true,
		Tura:    tura,
	}
	pusty := ""
	fragment.Chunk.Text = &pusty
	if surowe, err := json.Marshal(tura); err == nil {
		fragment.Chunk.Data = surowe
	}
	_ = wyslij(kontekst, na, fragment)
}

// zakonczBledem domyka strumień fragmentem rodzaju „błąd". Błąd dotyczy
// wyłącznie tego wywołania: konto zostaje czynne, sesja też.
func zakonczBledem(kontekst context.Context, na chan<- Fragment, z Zapytanie, kod shared.ErrorCode, tresc string) {
	fragment := shared.StreamChunkEvent{
		WindowId:  z.IdOkna,
		MessageId: z.IdWiadomosci,
		Kind:      shared.ChunkKindError,
		Text:      &tresc,
	}
	if surowe, err := json.Marshal(shared.ErrorInfo{
		Code:      kod,
		Message:   tresc,
		Retryable: shared.KodyPonawialne[kod],
	}); err == nil {
		fragment.Data = surowe
	}
	_ = wyslij(kontekst, na, Fragment{Chunk: fragment, Ostatni: true})
}
