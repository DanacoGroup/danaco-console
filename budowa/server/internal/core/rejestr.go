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

// Rejestr wiąże nazwę komendy z obsługiwaczem. Jest jedynym miejscem, w którym
// rdzeń rozstrzyga „co wykonać" — nie ma drugiego łańcucha warunków rozsianego
// po obsługiwaczach.
//
// Rejestr nie zawiera ani jednego literału nazwy. Nazwy wstrzykują pliki
// handlers_*.go, biorąc je ze stałych pakietu shared.
type Rejestr struct {
	wpisy map[shared.MessageType]Obsluga
}

// NowyRejestr zakłada pusty rejestr obsługiwaczy.
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

// Obsluga zwraca obsługiwacza komendy oraz informację, czy jest zarejestrowany.
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

// Liczba zwraca liczbę zarejestrowanych komend — do dziennika startu rdzenia.
func (r *Rejestr) Liczba() int {
	if r == nil {
		return 0
	}
	return len(r.wpisy)
}

// RejestrKomend buduje rozpoznanie nazw warstwy protokołu ze zbioru komend
// rzeczywiście obsługiwanych. Komenda spoza tego zbioru — także taka, która jest
// w kontrakcie, lecz nie ma jeszcze obsługiwacza — dostaje odpowiedź
// `*.unknown` zamiast błędu zrywającego.
func (r *Rejestr) RejestrKomend() *protocol.RejestrKomend {
	return protocol.NowyRejestrKomend(r.Nazwy()...)
}
