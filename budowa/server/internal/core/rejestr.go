package core

import (
	"context"
	"sort"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// Obsluga wykonuje jedną komendę kontraktu i zwraca odpowiedź protokołu.
// Nie zwraca błędu Go: błąd wykonania jest treścią odpowiedzi, bo dotyczy
// wyłącznie bieżącego wywołania i nie ma prawa zatrzymać sesji ani połączenia.
type Obsluga func(ctx context.Context, z protocol.Request) protocol.Odpowiedz

// Rejestr wiąże nazwę komendy z obsługiwaczem i jest jedynym miejscem, w którym rdzeń rozstrzyga, co wykonać. Nie zawiera ani jednego literału nazwy — nazwy wstrzykują pliki handlers_*.go ze stałych pakietu shared.
type Rejestr struct {
	wpisy map[shared.MessageType]Obsluga
}

// NowyRejestr zakłada pusty rejestr obsługiwaczy, gotowy do wypełnienia wpisami komend przez Zarejestruj.
func NowyRejestr() *Rejestr {
	return &Rejestr{wpisy: make(map[shared.MessageType]Obsluga)}
}

// Zarejestruj wiąże nazwę komendy z obsługiwaczem. Pusta nazwa i pusty
// obsługiwacz są pomijane — rejestr nie ma jak wywołać niczego, więc komenda
// trafi na ścieżkę `*.unknown` zamiast wywrócić proces.
func (r *Rejestr) Zarejestruj(nazwa shared.MessageType, obsluga Obsluga) {
	if r == nil || r.wpisy == nil || nazwa == "" || obsluga == nil {
		return
	}
	r.wpisy[nazwa] = obsluga
}

// Obsluga zwraca obsługiwacza komendy oraz informację, czy jest zarejestrowany, żeby komenda nieznana trafiła na ścieżkę „*.unknown".
func (r *Rejestr) Obsluga(nazwa shared.MessageType) (Obsluga, bool) {
	if r == nil || r.wpisy == nil {
		return nil, false
	}
	obsluga, jest := r.wpisy[nazwa]
	return obsluga, jest
}

// Nazwy zwraca komendy obsługiwane przez rdzeń, uporządkowane rosnąco.
// Zbiór ten jest jednocześnie odpowiedzią na pytanie klienta o zdolności rdzenia
// w powitaniu oraz podstawą rozpoznania komendy nieznanej.
func (r *Rejestr) Nazwy() []shared.MessageType {
	if r == nil || r.wpisy == nil {
		return nil
	}
	nazwy := make([]shared.MessageType, 0, len(r.wpisy))
	for nazwa := range r.wpisy {
		nazwy = append(nazwy, nazwa)
	}
	sort.Strings(nazwy)
	return nazwy
}

// Liczba zwraca liczbę zarejestrowanych komend, wykorzystywaną wyłącznie do zapisu w dzienniku startu rdzenia.
func (r *Rejestr) Liczba() int {
	if r == nil {
		return 0
	}
	return len(r.wpisy)
}

// RejestrKomend buduje rozpoznanie nazw warstwy protokołu ze zbioru komend rzeczywiście obsługiwanych. Komenda spoza tego zbioru dostaje odpowiedź „*.unknown" zamiast błędu zrywającego.
func (r *Rejestr) RejestrKomend() *protocol.RejestrKomend {
	return protocol.NowyRejestrKomend(r.Nazwy()...)
}
