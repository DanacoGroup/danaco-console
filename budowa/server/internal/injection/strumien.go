package injection

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"io"
	"strings"

	"danacoconsole/shared"
)

// rozmiarBuforaLinii jest wstępnym rozmiarem bufora czytnika. Linia zdarzenia
// bywa duża (wynik narzędzia), więc czytnik rośnie w razie potrzeby.
const rozmiarBuforaLinii = 256 << 10

// obserwacja jest tym, co kanał wyniósł z jednego przebiegu strumienia
// wyjścia procesu programu zewnętrznego.
type obserwacja struct {
	// Tura jest podsumowaniem zdarzenia result; nil znaczy, że tura nie
	// domknęła się poprawnie.
	Tura *ZakonczenieTury
	// Limit jest ostatnim stanem limitu tempa przekazanym przez program.
	Limit *limitTempaCLI
	// TekstPoszedl mówi, czy odbiorca zobaczył już jakikolwiek fragment tekstu
	// tury.
	TekstPoszedl bool
}

// czytajStrumien czyta wyjście procesu linia po linii i wysyła fragmenty do
// odbiorcy. Linia nieczytelna jako JSON nie przerywa strumienia — zostaje
// policzona i widoczna w podsumowaniu tury.
func czytajStrumien(kontekst context.Context, zrodlo io.Reader, z Zapytanie, na chan<- Fragment) (obserwacja, error) {
	czytnik := bufio.NewReaderSize(zrodlo, rozmiarBuforaLinii)
	wynik := obserwacja{}
	typy := nowyWykazTypow()
	nierozpoznane := 0

	for {
		linia, err := czytajLinie(czytnik)
		if linia != "" {
			zdarzenie, blad := odczytaj(linia)
			switch {
			case blad != nil:
				nierozpoznane++
			default:
				typy.dodaj(zdarzenie.Type)
				if blad := przetworz(kontekst, zdarzenie, linia, z, na, &wynik); blad != nil {
					return wynik, blad
				}
			}
		}
		if err != nil {
			if !errors.Is(err, io.EOF) {
				return wynik, err
			}
			break
		}
		if wynik.Tura != nil {
			break
		}
	}
	if wynik.Tura != nil {
		wynik.Tura.TypyZdarzen = typy.wykaz()
		wynik.Tura.LinieNierozpoznane = nierozpoznane
	}
	return wynik, nil
}

// czytajLinie pobiera jedną linię dowolnej długości; zwraca też błąd końca
// strumienia razem z ostatnią, niedomkniętą linią.
func czytajLinie(czytnik *bufio.Reader) (string, error) {
	linia, err := czytnik.ReadString('\n')
	return strings.TrimSpace(linia), err
}

func odczytaj(linia string) (zdarzenieCLI, error) {
	var zdarzenie zdarzenieCLI
	if err := json.Unmarshal([]byte(linia), &zdarzenie); err != nil {
		return zdarzenieCLI{}, err
	}
	return zdarzenie, nil
}

// przetworz zamienia jedno zdarzenie na fragmenty i uzupełnia obserwację.
// Surowa linia jedzie obok zdarzenia, bo zdarzenie zaczepu przekazuje ją dalej
// w postaci pierwotnej, nie w przekładzie na strukturę.
func przetworz(kontekst context.Context, zdarzenie zdarzenieCLI, linia string, z Zapytanie, na chan<- Fragment, wynik *obserwacja) error {
	if zdarzenie.RateLimit != nil {
		wynik.Limit = zdarzenie.RateLimit
	}
	// Zdarzenie zaczepu nie jest fragmentem kontraktu, idzie haczykiem do
	// warstwy składania.
	if zdarzenieZaczepu(zdarzenie) {
		if z.NaZdarzenieZaczepu != nil {
			z.NaZdarzenieZaczepu(zaczepZeZdarzenia(zdarzenie, linia))
		}
		return nil
	}
	switch zdarzenie.Type {
	case TypAssistant, TypUser:
		for _, blok := range zdarzenie.Message.bloki() {
			fragment, jest := blok.fragment(z.IdOkna, z.IdWiadomosci)
			if !jest {
				continue
			}
			if fragment.Kind == shared.ChunkKindText {
				wynik.TekstPoszedl = true
			}
			if err := wyslij(kontekst, na, Fragment{Chunk: fragment}); err != nil {
				return err
			}
		}
	case TypResult:
		wynik.Tura = &ZakonczenieTury{
			IdSesjiCLI: zdarzenie.SessionID,
			Podtyp:     zdarzenie.Subtype,
			Blad:       zdarzenie.IsError,
			Tekst:      zdarzenie.Result,
			Koszt:      zdarzenie.CostUSD,
			Tury:       zdarzenie.NumTurns,
			CzasMs:     zdarzenie.DurationM,
		}
	}
	return nil
}

// wyslij oddaje fragment odbiorcy strumienia przez kanał, honorując przy
// tym odwołanie kontekstu wywołania.
func wyslij(kontekst context.Context, na chan<- Fragment, fragment Fragment) error {
	select {
	case na <- fragment:
		return nil
	case <-kontekst.Done():
		return kontekst.Err()
	}
}

// wykazTypow zbiera typy zdarzeń strumienia w kolejności ich pierwszego
// wystąpienia w danej turze rozmowy.
type wykazTypow struct {
	widziane map[string]bool
	kolejno  []string
}

func nowyWykazTypow() *wykazTypow {
	return &wykazTypow{widziane: map[string]bool{}}
}

func (w *wykazTypow) dodaj(typ string) {
	if typ == "" || w.widziane[typ] {
		return
	}
	w.widziane[typ] = true
	w.kolejno = append(w.kolejno, typ)
}

func (w *wykazTypow) wykaz() []string {
	return w.kolejno
}
