// Odpowiedzialność pliku: klient protokołu debugowania (DAP) — ramkowanie
// komunikatów nagłówkiem `Content-Length`, korelacja odpowiedzi z żądaniami
// numerem kolejnym i odbiór zdarzeń adaptera osobnym kanałem.
package core

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"
)

// czasOdpowiedziDap jest granicą czekania na jedną odpowiedź adaptera.
// Zatrzymanie programu bywa powolne, lecz brak granicy oznaczałby okno wiszące
// bez końca na adapterze, który przestał odpowiadać.
const czasOdpowiedziDap = 30 * time.Second

// komunikatDap jest wspólnym kształtem żądania, odpowiedzi i zdarzenia
// protokołu DAP, ramkowanym nagłówkiem długości treści.
type komunikatDap struct {
	Seq        int             `json:"seq"`
	Type       string          `json:"type"`
	Command    string          `json:"command,omitempty"`
	Event      string          `json:"event,omitempty"`
	RequestSeq int             `json:"request_seq,omitempty"`
	Success    *bool           `json:"success,omitempty"`
	Message    string          `json:"message,omitempty"`
	Arguments  json.RawMessage `json:"arguments,omitempty"`
	Body       json.RawMessage `json:"body,omitempty"`
}

// klientDap prowadzi jedną rozmowę z adapterem debugowania jednego języka,
// korelując odpowiedzi numerem kolejnym żądania.
type klientDap struct {
	wejscie io.WriteCloser

	mu        sync.Mutex
	numer     int
	czekajacy map[int]chan komunikatDap
	zamkniety bool

	// zdarzenia oddaje komunikaty niezwiązane z żadnym żądaniem, buforowane.
	zdarzenia chan komunikatDap
}

// nowyKlientDap zakłada klienta nad strumieniami procesu adaptera i startuje
// gorutynę odbiorczą jego komunikatów.
func nowyKlientDap(wejscie io.WriteCloser, wyjscie io.Reader) *klientDap {
	klient := &klientDap{
		wejscie:   wejscie,
		czekajacy: map[int]chan komunikatDap{},
		zdarzenia: make(chan komunikatDap, 256),
	}
	go klient.odbieraj(wyjscie)
	return klient
}

// odbieraj czyta komunikaty adaptera i rozdziela je na odpowiedzi i zdarzenia,
// aż do zamknięcia strumienia albo błędu odczytu.
func (k *klientDap) odbieraj(wyjscie io.Reader) {
	defer k.zamknijCzekajacych()
	if wyjscie == nil {
		return
	}
	czytnik := bufio.NewReaderSize(wyjscie, 64*1024)
	for {
		komunikat, err := czytajKomunikatDap(czytnik)
		if err != nil {
			return
		}
		if komunikat.Type == "event" {
			select {
			case k.zdarzenia <- komunikat:
			default:
				// Kanał pełny znaczy obserwatora bez nadążania: zdarzenie przepada.
			}
			continue
		}
		k.mu.Lock()
		odbiorca, jest := k.czekajacy[komunikat.RequestSeq]
		if jest {
			delete(k.czekajacy, komunikat.RequestSeq)
		}
		k.mu.Unlock()
		if jest {
			odbiorca <- komunikat
			close(odbiorca)
		}
	}
}

// czytajKomunikatDap czyta jeden komunikat wraz z jego nagłówkiem długości,
// z bufora strumienia wyjściowego adaptera.
func czytajKomunikatDap(czytnik *bufio.Reader) (komunikatDap, error) {
	dlugosc := 0
	for {
		wiersz, err := czytnik.ReadString('\n')
		if err != nil {
			return komunikatDap{}, err
		}
		tresc := strings.TrimRight(wiersz, "\r\n")
		if tresc == "" {
			break
		}
		if wartosc, jest := strings.CutPrefix(tresc, "Content-Length:"); jest {
			dlugosc, err = strconv.Atoi(strings.TrimSpace(wartosc))
			if err != nil {
				return komunikatDap{}, errors.New("adapter debugowania podał nieczytelną długość komunikatu")
			}
		}
	}
	if dlugosc <= 0 {
		return komunikatDap{}, errors.New("adapter debugowania przysłał komunikat bez długości")
	}
	bajty := make([]byte, dlugosc)
	if _, err := io.ReadFull(czytnik, bajty); err != nil {
		return komunikatDap{}, err
	}
	var komunikat komunikatDap
	if err := json.Unmarshal(bajty, &komunikat); err != nil {
		return komunikatDap{}, err
	}
	return komunikat, nil
}

// Wolaj wysyła żądanie i czeka na odpowiedź adaptera.
//
// Odpowiedź nieudana wraca błędem niosącym treść adaptera. To jest odpowiedź,
// a nie awaria rdzenia: „nie da się obliczyć tego wyrażenia” jest zdaniem, które
// Operator ma przeczytać.
func (k *klientDap) Wolaj(komenda string, argumenty any) (json.RawMessage, error) {
	tresc, err := json.Marshal(argumenty)
	if err != nil {
		return nil, err
	}
	if argumenty == nil {
		tresc = nil
	}

	k.mu.Lock()
	if k.zamkniety {
		k.mu.Unlock()
		return nil, errors.New("adapter debugowania zakończył pracę")
	}
	k.numer++
	numer := k.numer
	odbiorca := make(chan komunikatDap, 1)
	k.czekajacy[numer] = odbiorca
	k.mu.Unlock()

	zadanie := komunikatDap{Seq: numer, Type: "request", Command: komenda, Arguments: tresc}
	if err := k.wyslij(zadanie); err != nil {
		k.mu.Lock()
		delete(k.czekajacy, numer)
		k.mu.Unlock()
		return nil, err
	}

	select {
	case odpowiedz, otwarty := <-odbiorca:
		if !otwarty {
			return nil, errors.New("adapter debugowania rozłączył się przed odpowiedzią na " + komenda)
		}
		if odpowiedz.Success != nil && !*odpowiedz.Success {
			tresc := odpowiedz.Message
			if tresc == "" {
				tresc = "adapter odmówił czynności " + komenda
			}
			return nil, errors.New(tresc)
		}
		return odpowiedz.Body, nil
	case <-time.After(czasOdpowiedziDap):
		k.mu.Lock()
		delete(k.czekajacy, numer)
		k.mu.Unlock()
		return nil, errors.New("adapter debugowania nie odpowiedział na " + komenda +
			" w granicy czasu")
	}
}

// wyslij zapisuje komunikat wraz z nagłówkiem długości do strumienia
// wejściowego procesu adaptera, pod zamkiem wykluczającym równoległy zapis.
func (k *klientDap) wyslij(komunikat komunikatDap) error {
	bajty, err := json.Marshal(komunikat)
	if err != nil {
		return err
	}
	if k.wejscie == nil {
		return errors.New("adapter debugowania nie ma strumienia wejścia")
	}
	naglowek := "Content-Length: " + strconv.Itoa(len(bajty)) + "\r\n\r\n"
	if _, err := k.wejscie.Write(append([]byte(naglowek), bajty...)); err != nil {
		return errors.New("nie można wysłać komunikatu do adaptera debugowania: " + err.Error())
	}
	return nil
}

// Zdarzenia oddaje kanał zdarzeń adaptera, niezwiązanych z żadnym żądaniem,
// bezpośrednio dla obserwatora sesji.
func (k *klientDap) Zdarzenia() <-chan komunikatDap { return k.zdarzenia }

// Zamknij kończy rozmowę i zwalnia wszystkich czekających na odpowiedź,
// zamykając strumień wejściowy procesu.
func (k *klientDap) Zamknij() {
	k.mu.Lock()
	if k.zamkniety {
		k.mu.Unlock()
		return
	}
	k.zamkniety = true
	k.mu.Unlock()
	if k.wejscie != nil {
		_ = k.wejscie.Close()
	}
}

// zamknijCzekajacych zwalnia wszystkich, którzy czekają na odpowiedź, gdy
// adapter przestał mówić. Bez tego kroku każdy z nich czekałby do granicy czasu
// osobno, a okno stałoby przez pół minuty na komunikacie, którego już nie będzie.
func (k *klientDap) zamknijCzekajacych() {
	k.mu.Lock()
	k.zamkniety = true
	czekajacy := k.czekajacy
	k.czekajacy = map[int]chan komunikatDap{}
	k.mu.Unlock()
	for _, odbiorca := range czekajacy {
		close(odbiorca)
	}
}
