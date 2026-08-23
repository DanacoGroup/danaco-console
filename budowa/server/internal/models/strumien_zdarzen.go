package models

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

// maksymalnaLiniaStrumienia — górny rozmiar jednej linii strumienia zdarzeń.
// Pojedyncze zdarzenie potrafi nieść całą treść narzędzia, dlatego domyślny
// bufor skanera jest za mały.
const maksymalnaLiniaStrumienia = 4 << 20

// przedrostekDanych — przedrostek linii z treścią w strumieniu zdarzeń SSE.
const przedrostekDanych = "data:"

// znacznikKonca — umowny znacznik zamknięcia strumienia zdarzeń.
const znacznikKonca = "[DONE]"

// CzytajZdarzenia czyta strumień zdarzeń linia po linii i podaje obsłudze samą
// treść JSON zdarzenia. Rozpoznaje obie postacie, w których dostawcy nadają
// strumień: linie SSE z przedrostkiem "data:" oraz gołe linie JSON. Linie puste,
// komentarze SSE i znacznik końca pomija.
func CzytajZdarzenia(zrodlo io.Reader, obsluga func(zdarzenie []byte) error) error {
	skaner := bufio.NewScanner(zrodlo)
	skaner.Buffer(make([]byte, 0, 64<<10), maksymalnaLiniaStrumienia)
	for skaner.Scan() {
		zdarzenie, jest := trescLinii(skaner.Bytes())
		if !jest {
			continue
		}
		if err := obsluga(zdarzenie); err != nil {
			return err
		}
	}
	if err := skaner.Err(); err != nil {
		return fmt.Errorf("models: odczyt strumienia zdarzeń: %w", err)
	}
	return nil
}

// trescLinii zwraca treść JSON zdarzenia z pojedynczej linii strumienia.
// Drugi wynik fałszywy oznacza linię bez treści albo koniec strumienia.
func trescLinii(linia []byte) ([]byte, bool) {
	oczyszczona := bytes.TrimSpace(linia)
	if len(oczyszczona) == 0 {
		return nil, false
	}
	tekst := string(oczyszczona)
	if strings.HasPrefix(tekst, przedrostekDanych) {
		tekst = strings.TrimSpace(strings.TrimPrefix(tekst, przedrostekDanych))
	} else if strings.HasPrefix(tekst, ":") {
		return nil, false
	} else if !strings.HasPrefix(tekst, "{") && !strings.HasPrefix(tekst, "[") {
		return nil, false
	}
	if tekst == "" || tekst == znacznikKonca {
		return nil, false
	}
	return []byte(tekst), true
}
