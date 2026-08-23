// Odpowiedzialność pliku: droga jednego wywołania — narzędzie → komenda →
// koperta → rdzeń → wynik.
//
// Rozdzielnia nie zna ani protokołu MCP, ani gniazda rdzenia. Zna wykaz
// kontraktu, odwzorowanie `shared.KomendyNarzedzi` i port `Rdzen`, więc daje
// się sprawdzić bez procesu po drugiej stronie.
//
// Błąd rdzenia wraca treścią. Odmowa, komenda nieznana rdzeniowi,
// zerwane gniazdo — każda z tych rzeczy wraca do modelu jako czytelny opis
// błędu narzędzia. Model czyta, poprawia i próbuje dalej; połączenie MCP nie
// pada ani razu.
package narzedzia

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync/atomic"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// Rdzen jest portem do rdzenia Danaco Console: koperta w jedną stronę, koperta
// w drugą. Interfejs stoi po stronie odbiorcy, więc rozdzielnia nie wie, czy po
// drugiej stronie jest gniazdo, próba, czy cokolwiek innego.
type Rdzen interface {
	// Wykonaj oddaje żądanie rdzeniowi i czeka na odpowiedź o tym samym
	// identyfikatorze. Błąd oznacza brak odpowiedzi, nie odmowę rdzenia —
	// odmowa przychodzi kopertą z wypełnionym polem błędu.
	Wykonaj(kontekst context.Context, zadanie protocol.Koperta) (protocol.Koperta, error)
}

// Gniazdo rdzenia jest realizacją tego portu. Sprawdzenie stoi przy kompilacji,
// żeby rozjazd sygnatur zatrzymał budowanie, a nie proces modelu.
var _ Rdzen = (*Polaczenie)(nil)

// Rozdzielnia kieruje wywołania narzędzi do rdzenia w zasięgu jednego okna.
type Rozdzielnia struct {
	rdzen Rdzen
	// okno jest oknem rozmowy, do którego należy ten serwer narzędzi. Puste
	// znaczy: uzupełniania nie ma — model podaje okno sam albo rdzeń odmawia.
	okno string
	// zasieg jest rolą okna, w którego imieniu serwer pracuje (`zasieg_roli.go`).
	// Przychodzi z wpisu MCP ułożonego przez rdzeń i nie zmienia się w czasie
	// życia procesu — model nie ma czym go przestawić.
	zasieg Zasieg
	// licznik nadaje identyfikatory żądań. Odpowiedzi wiąże się z żądaniami po
	// identyfikatorze, a ten wystarczy, by był jednoznaczny w obrębie jednego
	// połączenia — innych połączeń rozdzielnia nie prowadzi.
	licznik atomic.Uint64
	// dobor prowadzi zasięg eksperta (`rozdzielnia_ekspert.go`); nil znaczy
	// zasięg bez eksperta i wtedy ta droga nie rusza ani razu.
	dobor *doborEksperta
}

// NowaRozdzielnia wiąże rozdzielnię z rdzeniem, oknem rozmowy i rolą tego okna.
//
// Rola jest parametrem wymaganym, nie doklejką z wartością domyślną: zasięg
// rozstrzyga o tym, co model może zrobić, więc każdy, kto rozdzielnię składa,
// ma powiedzieć wprost, w czyim imieniu ona pracuje. Zasięgiem zwykłym jest
// `ZasiegOkna`.
func NowaRozdzielnia(rdzen Rdzen, okno string, zasieg Zasieg) *Rozdzielnia {
	return &Rozdzielnia{rdzen: rdzen, okno: okno, zasieg: zasieg}
}

// ZEkspertem wiąże rozdzielnię z kodem eksperta nałożonego na okno i włącza
// dobór narzędzi (`rozdzielnia_ekspert.go`).
//
// Osobno od konstruktora, a nie kolejnym jego parametrem: zasięg eksperta jest
// jedynym, który potrzebuje wartości, a dokładanie jej wszystkim pozostałym
// kazałoby im podawać pustkę bez znaczenia. Kod pusty albo zasięg inny niż
// `ZasiegEksperta` nie włącza niczego — rozdzielnia zostaje rozdzielnią, jaką
// była, co do znaku.
//
// Dziennik jest tu wymagany, a nie opcjonalny, bo cena zestawu i każdy brak
// zawężenia mają dokądś dojechać; dziennik nil ucisza je, więc dobór bez
// dziennika byłby dokładnie tą cichą degradacją, przeciw której powstał —
// wywołujący ma powiedzieć wprost, gdzie meldunki idą.
func (r *Rozdzielnia) ZEkspertem(kod string, dziennik *log.Logger) *Rozdzielnia {
	if r.zasieg != ZasiegEksperta || kod == "" {
		return r
	}
	r.dobor = &doborEksperta{kod: kod, dziennik: dziennik}
	return r
}

// ZDolozeniamiSesji wnosi doraźne dołożenia sesji — narzędzia dorzucone przez
// Operatora poza definicją eksperta (`--dolozenia`).
//
// Bez zawężenia dołożenia nie mają czego dołożyć: okno bez eksperta ma pełny
// wykaz, więc każde dołożenie już w nim stoi. Wołanie tej metody poza zasięgiem
// eksperta nie jest więc błędem, tylko czynnością pustą — rozstrzygnięcie, czy
// argument w ogóle wysłać, należy do strony rdzenia, która wie o oknie więcej
// (`core/adapter_rozmowa_zestaw.go` melduje ten przypadek osobno).
func (r *Rozdzielnia) ZDolozeniamiSesji(nazwy []string) *Rozdzielnia {
	if r.dobor == nil || len(nazwy) == 0 {
		return r
	}
	r.dobor.dolozenia = nazwy
	return r
}

// Narzedzia zwraca wykaz narzędzi podawany modelowi: wykaz kontraktu wraz z tym,
// co dokłada rola okna, a w zasięgu eksperta — podzbiór wskazany przez jego
// definicję. Rozdzielnia niczego nie wymyśla: o rozszerzeniu rozstrzyga `zasieg`
// ustalony przy uruchomieniu, o zawężeniu — rdzeń zapytany o eksperta.
//
// Kontekst wchodzi parametrem, bo w zasięgu eksperta ta droga pyta rdzeń, a
// pytanie bez kontekstu nie dałoby się przerwać razem z resztą procesu. Zasięgi
// pozostałe kontekstu nie tykają.
func (r *Rozdzielnia) Narzedzia(kontekst context.Context) []Narzedzie {
	wykaz := WykazZasiegu(r.zasieg)
	if r.dobor == nil {
		return wykaz
	}
	return r.dobor.wykaz(kontekst, r.rdzen, wykaz)
}

// Wywolaj wykonuje jedno wywołanie narzędzia i zwraca treść wyniku.
//
// Zwrócony błąd jest treścią dla modelu, nie usterką procesu: warstwa protokołu
// odsyła go jako wynik narzędzia oznaczony jako błędny.
func (r *Rozdzielnia) Wywolaj(kontekst context.Context, nazwa string,
	argumenty map[string]any) (string, error) {

	komenda, parametry, dostepne := r.rozpoznaj(nazwa)
	if !dostepne {
		return "", odmowaNarzedzia(nazwa)
	}
	// Zawężenie ma obowiązywać także przy wywołaniu — inaczej byłoby wyłącznie
	// oszczędnością żetonów, a nie doborem narzędzi eksperta.
	if r.dobor != nil {
		if wolno, _ := r.dobor.wolno(nazwa); !wolno {
			return "", r.dobor.odmowaPozaDoborem(nazwa)
		}
	}
	zadanie, err := r.koperta(komenda, zZasiegiemOkna(argumenty, parametry, r.okno))
	if err != nil {
		return "", err
	}
	odpowiedz, err := r.rdzen.Wykonaj(kontekst, zadanie)
	if err != nil {
		return "", fmt.Errorf("narzędzie %s: rdzeń nie odpowiedział: %w", nazwa, err)
	}
	return wynikKoperty(nazwa, odpowiedz)
}

// rozpoznaj rozstrzyga, czy nazwa jest w zasięgu tego okna, i oddaje komendę
// wraz z nazwami pól, które wolno uzupełnić oknem serwera.
//
// Narzędzie dołożone przez rolę nie dostaje uzupełnienia okna — i to jest jego
// istota, nie przeoczenie. Rozszerzenie okna asystenta służy nastawianiu okna
// docelowego; podstawienie okna serwera w brakujące `windowId` kazałoby
// asystentowi przestawić kanał modelu samemu sobie, czyli zrobić dokładnie to,
// czego zakaz trzyma te komendy poza wykazem kontraktu. Okno docelowe model
// wskazuje jawnie albo rdzeń odmawia — i tak ma być.
func (r *Rozdzielnia) rozpoznaj(nazwa string) (shared.MessageType, []string, bool) {
	pozycja, wWykazie := deklaracja(nazwa)
	if komenda, objeta := shared.KomendyNarzedzi[nazwa]; wWykazie && objeta {
		return komenda, nazwyParametrow(pozycja), true
	}
	if komenda, zRoli := komendaRoli(nazwa, r.zasieg); zRoli {
		return komenda, nil, true
	}
	return "", nil, false
}

// koperta pakuje argumenty w kopertę kontraktu {type, id, payload}.
func (r *Rozdzielnia) koperta(komenda shared.MessageType, argumenty map[string]any) (protocol.Koperta, error) {
	identyfikator := "narzedzia-" + strconv.FormatUint(r.licznik.Add(1), 10)
	zadanie, err := protocol.NowaKoperta(komenda, identyfikator, "", argumenty)
	if err != nil {
		return protocol.Koperta{}, fmt.Errorf("komenda %s: argumenty nie dają się zakodować: %w", komenda, err)
	}
	return zadanie, nil
}

// wynikKoperty wyjmuje z odpowiedzi treść dla modelu albo opis odmowy rdzenia.
func wynikKoperty(nazwa string, odpowiedz protocol.Koperta) (string, error) {
	if odpowiedz.Error != nil {
		return "", fmt.Errorf("narzędzie %s: %s", nazwa, protocol.Opis(*odpowiedz.Error))
	}
	if odpowiedz.Status != nil && *odpowiedz.Status != shared.EnvelopeStatusOk {
		return "", fmt.Errorf("narzędzie %s: rdzeń odpowiedział stanem %s bez opisu błędu", nazwa, *odpowiedz.Status)
	}
	if len(odpowiedz.Payload) == 0 {
		// Komenda potwierdzona bez treści właściwej. Napis pusty byłby dla modelu
		// nieodróżnialny od usterki, więc wraca poprawny, pusty obiekt JSON.
		return "{}", nil
	}
	czytelny, err := json.MarshalIndent(json.RawMessage(odpowiedz.Payload), "", "  ")
	if err != nil {
		return string(odpowiedz.Payload), nil
	}
	return string(czytelny), nil
}

// nazwyParametrow zwraca nazwy pól treści żądania deklaracji.
func nazwyParametrow(pozycja shared.ToolDeclaration) []string {
	nazwy := make([]string, 0, len(pozycja.Parameters))
	for _, parametr := range pozycja.Parameters {
		nazwy = append(nazwy, parametr.Name)
	}
	return nazwy
}
