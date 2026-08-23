package models

import (
	"reflect"
	"strings"
)

// zbudujKanal zwraca kanał dla wiersza: zachowany, gdy wiersz nie zmienił się
// od poprzedniego odświeżenia, w przeciwnym razie zbudowany fabryką.
func zbudujKanal(d Definicja, fabryki Fabryki, poprzednie map[string]Kanal,
	poprzednieDef map[string]Definicja) (Kanal, bool, string) {
	klucz := d.Identyfikator()
	if istniejacy, jest := poprzednie[klucz]; jest && reflect.DeepEqual(poprzednieDef[klucz], d) {
		return istniejacy, true, ""
	}
	fabryka, jest := fabryki[d.KluczAdaptera()]
	if !jest {
		return nil, false, "brak fabryki adaptera " + d.KluczAdaptera()
	}
	kanal, err := fabryka(d)
	if err != nil {
		return nil, false, err.Error()
	}
	if kanal == nil {
		return nil, false, "fabryka adaptera " + d.KluczAdaptera() + " nie zwróciła kanału"
	}
	return kanal, false, ""
}

// kluczeKanalu wskazuje, pod czym kanał jest osiągalny: identyfikatorem wiersza
// oraz kodem kanału. Dwie drogi, bo okno komunikacji trzyma identyfikator,
// a konfiguracja i test posługują się kodem.
func kluczeKanalu(d Definicja) []string {
	klucze := []string{d.Identyfikator()}
	if kod := strings.TrimSpace(d.Kod); kod != "" {
		klucze = append(klucze, kod)
	}
	return klucze
}

// zamknijNieuzywane zwalnia zasoby kanałów, które wypadły z rejestru.
func zamknijNieuzywane(poprzednie map[string]Kanal, zachowane map[Kanal]bool) {
	zamkniete := map[Kanal]bool{}
	for _, kanal := range poprzednie {
		if zachowane[kanal] || zamkniete[kanal] {
			continue
		}
		zamkniete[kanal] = true
		if zamykalny, jest := kanal.(Zamykalny); jest {
			_ = zamykalny.Zamknij()
		}
	}
}
