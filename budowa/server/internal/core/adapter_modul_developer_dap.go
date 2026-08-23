// Odpowiedzialność pliku: klient protokołu debugowania (DAP) — ramkowanie
// komunikatów, korelacja odpowiedzi z żądaniami i odbiór zdarzeń adaptera.
//
// ── Dlaczego DAP, a nie własne API debuggera ────────────────────────────────
// Opracowanie modułu (rozdz. 7.5) wskazuje DAP jako warstwę wspólną: `delve`
// dla Go, `debugpy` dla Pythona, `js-debug` dla Node, `lldb` dla C. Rozmowa
// z każdym z nich po jego własnym API oznaczałaby cztery różne implementacje
// tej samej pięciu czynności — a Run & Debug ma jedno okno i jeden zestaw
// przycisków, niezależnie od języka.
//
// ── Ramkowanie ──────────────────────────────────────────────────────────────
// DAP jedzie po strumieniu bajtów, więc komunikat ma nagłówek
// `Content-Length: N`, pustą linię i N bajtów treści JSON. To jest to samo
// ramkowanie, którego używa LSP. Bez niego dwa komunikaty wysłane po sobie
// zlałyby się w jeden nieczytelny dokument.
//
// ── Korelacja ───────────────────────────────────────────────────────────────
// Adapter odpowiada w dowolnej kolejności i wtrąca między odpowiedzi zdarzenia
// (`stopped`, `terminated`, `output`), więc czekanie na „następny komunikat”
// byłoby czekaniem na cokolwiek. Każde żądanie dostaje numer kolejny, a odbiornik
// rozdziela przychodzące komunikaty: odpowiedź trafia do kanału czekającego na
// ten numer, zdarzenie — do obserwatora sesji.
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

// komunikatDap jest wspólnym kształtem żądania, odpowiedzi i zdarzenia.
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

// klientDap prowadzi jedną rozmowę z adapterem debugowania.
type klientDap struct {
	wejscie io.WriteCloser

	mu        sync.Mutex
	numer     int
	czekajacy map[int]chan komunikatDap
	zamkniety bool

	// zdarzenia oddaje komunikaty niezwiązane z żadnym żądaniem. Kanał jest
	// buforowany, bo adapter wysyła zdarzenia szybciej, niż sesja je czyta,
	// a zablokowany odbiornik wstrzymałby także odpowiedzi.
	zdarzenia chan komunikatDap
}

// nowyKlientDap zakłada klienta nad strumieniami procesu adaptera.
func nowyKlientDap(wejscie io.WriteCloser, wyjscie io.Reader) *klientDap {
	klient := &klientDap{
		wejscie:   wejscie,
		czekajacy: map[int]chan komunikatDap{},
		zdarzenia: make(chan komunikatDap, 256),
	}
	go klient.odbieraj(wyjscie)
	return klient
}

// odbieraj czyta komunikaty adaptera i rozdziela je na odpowiedzi i zdarzenia.
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
				// Kanał pełny znaczy obserwatora, który nie nadąża. Zdarzenie
				// przepada, a rozmowa idzie dalej: zablokowanie odbiornika
				// zatrzymałoby także odpowiedzi na żądania.
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

// czytajKomunikatDap czyta jeden komunikat wraz z jego nagłówkiem długości.
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

// wyslij zapisuje komunikat wraz z nagłówkiem długości.
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

// Zdarzenia oddaje kanał zdarzeń adaptera.
func (k *klientDap) Zdarzenia() <-chan komunikatDap { return k.zdarzenia }

// Zamknij kończy rozmowę i zwalnia czekających.
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
