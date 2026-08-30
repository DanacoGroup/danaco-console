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

// NowyKanal składa nowy Kanal nad przekazaną pulą kont, gotowy od razu do
// prowadzenia rozmów z modelem.
func NowyKanal(pula *PulaKont) *Kanal {
	if pula == nil {
		pula = NowaPula()
	}
	return &Kanal{pula: pula}
}

// Pula udostępnia pulę kont, nad którą kanał pracuje, do odczytu jej stanu
// przez rejestr kanałów platformy.
func (k *Kanal) Pula() *PulaKont {
	return k.pula
}

// Rozmowa przeprowadza jedną turę i oddaje strumień fragmentów; kanał nigdy
// nie zwraca błędu obok strumienia.
func (k *Kanal) Rozmowa(kontekst context.Context, z Zapytanie) <-chan Fragment {
	strumien := make(chan Fragment)
	go func() {
		defer close(strumien)
		k.prowadz(kontekst, z, strumien)
	}()
	return strumien
}

// prowadz prowadzi jedną turę razem z rotacją kont, próbując kolejnych kont
// puli po kolei aż do sukcesu.
func (k *Kanal) prowadz(kontekst context.Context, z Zapytanie, na chan<- Fragment) {
	// Konto wskazane wygrywa z rotacją i jest rozkazem tożsamości tury.
	if z.Ustawienia.Konto != "" {
		k.prowadzWskazanym(kontekst, z, na)
		return
	}
	powod := ""
	for proba := 1; proba <= len(k.pula.Konta())+1; proba++ {
		konto, jest := k.pula.Biezace()
		if !jest {
			// Pula pusta nie znaczy pula wyczerpana: brak kont oznacza tożsamość
			// otoczenia, nie odmowę.
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

		// Pułap sprawdza się przed wyczerpaniem i zatrzymuje rotację kont puli.
		if pulap, naPulapie := rozpoznajPulap(wynik.Obserwacja, wynik.Bledy); naPulapie {
			zakonczPulapem(kontekst, na, z, pulap)
			return
		}

		wyczerpanie, wyczerpane := rozpoznajWyczerpanie(wynik.Obserwacja, wynik.Bledy)
		if !wyczerpane {
			zakoncz(kontekst, na, z, wynik.Obserwacja, konto.Kod)
			return
		}

		nastepne, dostepne := k.pula.Wyczerpane(konto.Kod, wyczerpanie.DoChwili)
		if !dostepne || nastepne.Kod == konto.Kod || wynik.Obserwacja.TekstPoszedl {
			zakoncz(kontekst, na, z, wynik.Obserwacja, konto.Kod)
			return
		}
		powod = fmt.Sprintf("rotacja konta %s → %s (%s)", konto.Kod, nastepne.Kod, wyczerpanie.Powod)
	}
	zakonczBledem(kontekst, na, z, shared.ErrorCodeChannelUnavailable,
		"pula kont wyczerpana w tej turze; sesja pozostaje czynna")
}

// prowadzWskazanym prowadzi turę na koncie wskazanym konfiguracją, bez
// rotacji na inne konto po wyczerpaniu limitu.
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
		// Ślad wyczerpania idzie do puli, ale przełączenia konta tu nie ma.
		k.pula.Wyczerpane(konto.Kod, wyczerpanie.DoChwili)
	}
	zakoncz(kontekst, na, z, wynik.Obserwacja, konto.Kod)
}

// zakoncz wysyła fragment kończący turę. Gdy tura nie przyniosła podsumowania,
// zamiast niego jedzie fragment błędu — strumień zawsze ma koniec. Strumień
// domknięty wcześniej zostaje bez zmian: drugi fragment ostatni zerwałby kontrakt.
func zakoncz(kontekst context.Context, na chan<- Fragment, z Zapytanie, obserwacja obserwacja, konto string) {
	if obserwacja.Domkniety {
		return
	}
	tura := obserwacja.Tura
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
