// Odpowiedzialność pliku: doprowadzenie do tury zlecenia asystenta
// konfiguracji MCP okna wraz z wpisem serwera narzędzi modelu, bez czego
// asystent nie jest klawiaturą Operatora, tylko drugim rozmówcą.
package core

import (
	"context"
	"encoding/json"

	"danacoconsole/server/internal/models"
	"danacoconsole/server/internal/narzedzia"
	"danacoconsole/server/internal/session"
)

// ZMostami wpina składacz konfiguracji mostów MCP okna — ten sam, którym
// jedzie rozmowa. Sprawdzenie na nil jest tu z tego samego powodu co przy
// `ZMowa`.
func (a *adapterAsystenta) ZMostami(m *mostyOkna) *adapterAsystenta {
	if m == nil {
		return a
	}
	a.mosty = m
	return a
}

// uzupelnijNarzedzia dokłada do zapytania zlecenia konfigurację MCP okna wraz
// z wpisem serwera narzędzi modelu, osobno od `zapytanieZlecenia`, bo tamta
// funkcja jest czysta.
func (a *adapterAsystenta) uzupelnijNarzedzia(ctx context.Context, okno session.Okno,
	zapytanie *models.Zapytanie) {

	if a == nil || a.mosty == nil || zapytanie == nil || okno.Id == "" {
		return
	}
	tekst := a.mosty.tekstZNarzedziami(ctx, okno.Id)
	if tekst == "" {
		return
	}
	zapytanie.KonfiguracjaMCP = zZasiegiemKlawiatury(tekst, okno)
}

// modulAsystenta jest identyfikatorem modułu, którego okno jest klawiaturą
// Operatora. Wartość jest wartością katalogu modułów, tą samą, którą oddaje
// `window.list` w polu `moduleId`.
const modulAsystenta = "assistant"

// zZasiegiemKlawiatury dopisuje do wpisu `danaco` rolę okna asystenta.
// Rozszerzenie wchodzi wyłącznie wtedy, gdy okno tury jest oknem modułu
// Assistant, i wchodzi do wpisu MCP, zanim proces modelu wystartuje.
func zZasiegiemKlawiatury(tekst string, okno session.Okno) string {
	if okno.Modul != modulAsystenta {
		return tekst
	}
	argumenty := narzedzia.ArgumentyZasiegu(narzedzia.ZasiegKlawiatury)
	if len(argumenty) == 0 {
		return tekst
	}
	var konfiguracja KonfiguracjaMostu
	if err := json.Unmarshal([]byte(tekst), &konfiguracja); err != nil {
		return tekst
	}
	wpis, jest := konfiguracja.McpServers[narzedzia.KluczWpisu]
	if !jest {
		return tekst
	}
	wpis.Args = append(append([]string{}, wpis.Args...), argumenty...)
	konfiguracja.McpServers[narzedzia.KluczWpisu] = wpis
	tresc, err := json.MarshalIndent(konfiguracja, "", "  ")
	if err != nil {
		return tekst
	}
	return string(tresc)
}
