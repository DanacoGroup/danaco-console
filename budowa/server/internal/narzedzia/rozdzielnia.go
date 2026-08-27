// Plik prowadzi drogę jednego wywołania narzędzia: narzędzie, komenda,
// koperta, rdzeń, wynik; błąd rdzenia wraca modelowi treścią, nie zrywa
// połączenia MCP.
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
	// Wykonaj oddaje żądanie rdzeniowi i czeka na odpowiedź; błąd oznacza brak
	// odpowiedzi, nie odmowę.
	Wykonaj(kontekst context.Context, zadanie protocol.Koperta) (protocol.Koperta, error)
}

// Gniazdo rdzenia jest realizacją tego portu. Sprawdzenie stoi przy kompilacji,
// żeby rozjazd sygnatur zatrzymał budowanie, a nie proces modelu.
var _ Rdzen = (*Polaczenie)(nil)

// Rozdzielnia kieruje wywołania narzędzi modelu do rdzenia, w zasięgu jednego
// okna rozmowy tej platformy.
type Rozdzielnia struct {
	rdzen Rdzen
	// okno jest oknem rozmowy, do którego należy ten serwer; puste znaczy brak
	// uzupełniania okna.
	okno string
	// zasieg jest rolą okna, w którego imieniu serwer pracuje; nie zmienia się
	// w czasie życia procesu.
	zasieg Zasieg
	// licznik nadaje identyfikatory żądań, jednoznaczne w obrębie jednego
	// prowadzonego połączenia.
	licznik atomic.Uint64
	// dobor prowadzi zasięg eksperta; nil znaczy zasięg bez eksperta, droga
	// doboru wtedy nie rusza.
	dobor *doborEksperta
}

// NowaRozdzielnia wiąże rozdzielnię z rdzeniem, oknem rozmowy i rolą tego
// okna; rola jest parametrem wymaganym, nie wartością domyślną.
func NowaRozdzielnia(rdzen Rdzen, okno string, zasieg Zasieg) *Rozdzielnia {
	return &Rozdzielnia{rdzen: rdzen, okno: okno, zasieg: zasieg}
}

// ZEkspertem wiąże rozdzielnię z kodem eksperta nałożonego na okno i włącza
// dobór narzędzi; kod pusty albo zasięg inny niż eksperta nie włącza niczego.
func (r *Rozdzielnia) ZEkspertem(kod string, dziennik *log.Logger) *Rozdzielnia {
	if r.zasieg != ZasiegEksperta || kod == "" {
		return r
	}
	r.dobor = &doborEksperta{kod: kod, dziennik: dziennik}
	return r
}

// ZDolozeniamiSesji wnosi doraźne dołożenia sesji — narzędzia dorzucone przez
// Operatora poza definicją eksperta; poza zasięgiem eksperta jest czynnością
// pustą.
func (r *Rozdzielnia) ZDolozeniamiSesji(nazwy []string) *Rozdzielnia {
	if r.dobor == nil || len(nazwy) == 0 {
		return r
	}
	r.dobor.dolozenia = nazwy
	return r
}

// Narzedzia zwraca wykaz narzędzi podawany modelowi: wykaz kontraktu wraz
// z rozszerzeniem roli okna, a w zasięgu eksperta zawężony do jego definicji.
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
	// Zawężenie ma obowiązywać także przy wywołaniu, nie tylko przy doborze
	// wykazu.
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
// wraz z nazwami pól uzupełnianych oknem serwera.
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

// koperta pakuje argumenty wywołania w kopertę kontraktu, gotową do wysłania
// do rdzenia jako żądanie tury.
func (r *Rozdzielnia) koperta(komenda shared.MessageType, argumenty map[string]any) (protocol.Koperta, error) {
	identyfikator := "narzedzia-" + strconv.FormatUint(r.licznik.Add(1), 10)
	zadanie, err := protocol.NowaKoperta(komenda, identyfikator, "", argumenty)
	if err != nil {
		return protocol.Koperta{}, fmt.Errorf("komenda %s: argumenty nie dają się zakodować: %w", komenda, err)
	}
	return zadanie, nil
}

// wynikKoperty wyjmuje z odpowiedzi rdzenia treść wyniku gotową dla modelu
// albo opis odmowy tego rdzenia.
func wynikKoperty(nazwa string, odpowiedz protocol.Koperta) (string, error) {
	if odpowiedz.Error != nil {
		return "", fmt.Errorf("narzędzie %s: %s", nazwa, protocol.Opis(*odpowiedz.Error))
	}
	if odpowiedz.Status != nil && *odpowiedz.Status != shared.EnvelopeStatusOk {
		return "", fmt.Errorf("narzędzie %s: rdzeń odpowiedział stanem %s bez opisu błędu", nazwa, *odpowiedz.Status)
	}
	if len(odpowiedz.Payload) == 0 {
		// Komenda potwierdzona bez treści właściwej; napis pusty byłby
		// nieodróżnialny od usterki.
		return "{}", nil
	}
	czytelny, err := json.MarshalIndent(json.RawMessage(odpowiedz.Payload), "", "  ")
	if err != nil {
		return string(odpowiedz.Payload), nil
	}
	return string(czytelny), nil
}

// nazwyParametrow zwraca nazwy pól treści żądania zadeklarowanych w deklaracji
// narzędzia tego kontraktu.
func nazwyParametrow(pozycja shared.ToolDeclaration) []string {
	nazwy := make([]string, 0, len(pozycja.Parameters))
	for _, parametr := range pozycja.Parameters {
		nazwy = append(nazwy, parametr.Name)
	}
	return nazwy
}
