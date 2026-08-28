// Plik dokłada wpis danaco — serwera narzędzi modelu — do konfiguracji MCP okna rozmowy, przez
// który model dostaje sterowanie platformą.
package core

import (
	"context"
	"encoding/json"

	"danacoconsole/server/internal/narzedzia"
)

// tekstZNarzedziami zwraca konfigurację MCP okna powiększoną o wpis serwera narzędzi, przepuszczając
// każdy krok dalej zamiast przerywać na braku nadań albo identyfikatora.
func (m *mostyOkna) tekstZNarzedziami(ctx context.Context, idOkna string) string {
	dotychczasowy := m.Tekst(ctx, idOkna)
	polecenie, argumenty, powod, jest := narzedzia.Wpis(idOkna)
	if !jest {
		m.zglosBrakNarzedzi(idOkna, powod)
		return dotychczasowy
	}
	konfiguracja := KonfiguracjaMostu{McpServers: map[string]WpisMostu{}}
	if dotychczasowy != "" {
		if err := json.Unmarshal([]byte(dotychczasowy), &konfiguracja); err != nil {
			return dotychczasowy
		}
	}
	konfiguracja.McpServers[narzedzia.KluczWpisu] = WpisMostu{
		Type:    RodzajWpisuMostu,
		Command: polecenie,
		Args:    argumenty,
	}
	tresc, err := json.MarshalIndent(konfiguracja, "", "  ")
	if err != nil {
		return dotychczasowy
	}
	return string(tresc)
}

// zglosBrakNarzedzi wpisuje do dziennika rdzenia powód, dla którego okno prowadzi turę bez
// sterowania platformą, raz na powód, nie raz na turę.
func (m *mostyOkna) zglosBrakNarzedzi(idOkna, powod string) {
	if m == nil || m.dziennik == nil || powod == "" || idOkna == "" {
		return
	}
	if _, juzBylo := m.zgloszoneBraki.LoadOrStore(powod, true); juzBylo {
		return
	}
	m.dziennik.Printf("serwer narzędzi modelu poza turą: %s — wpisu %q nie ma w konfiguracji MCP,"+
		" rozmowa toczy się bez sterowania platformą", powod, narzedzia.KluczWpisu)
}
