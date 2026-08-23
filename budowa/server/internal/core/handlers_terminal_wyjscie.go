// Wpięcie dwóch komend wyjścia obszaru `terminal.*` —
// `terminal.output.stream` (podgląd na żywo) i `terminal.output.read` (odczyt
// wyjścia jednego procesu).
//
// Port jest rozszerzeniem, nie drugim portem: `WyjscieTerminala` osadza
// `Terminal`, bo wyjście modułu ma tego samego właściciela co karty i procesy.
// Adapter wypełniający jeden wypełnia oba.
//
// Żadna z tych komend nie rozgłasza zdarzenia. Kontrakt daje modułowi jedno
// zdarzenie — `terminal.process.changed` — i dotyczy ono procesu, nie zapisu na
// strumień. Wiersze wyjścia jadą `stream.chunk` prosto z dziennika zbiorczego,
// a nie z obsługiwacza żądania: proces pisze długo po tym, jak odpowiedź na
// komendę już wróciła.
package core

import (
	"context"

	"danacoconsole/shared"
)

// WyjscieTerminala jest portem zbiorczego wyjścia modułu Terminal —
// portem Terminal rozszerzonym o komendę strumienia.
type WyjscieTerminala interface {
	Terminal

	// StrumienWyjscia obsługuje `terminal.output.stream`.
	StrumienWyjscia(ctx context.Context,
		z shared.TerminalOutputStreamRequest) (shared.TerminalOutputStreamResponse, error)

	// OdczytajWyjscie obsługuje `terminal.output.read`. Osobna czynność, bo
	// `terminal.command.exec` kończy się przy STARCIE procesu i wyniku nieść nie
	// może (adapter_modul_terminal_wyjscie_odczyt.go).
	OdczytajWyjscie(ctx context.Context,
		z shared.TerminalOutputReadRequest) (shared.TerminalOutputReadResponse, error)
}

// zarejestrujWyjscieTerminala wpina obie komendy wyjścia modułu.
func zarejestrujWyjscieTerminala(r *Rejestr, w WyjscieTerminala) {
	if r == nil || w == nil {
		return
	}
	r.Zarejestruj(shared.CommandTerminalOutputStream, obsluz(w.StrumienWyjscia))
	r.Zarejestruj(shared.CommandTerminalOutputRead, obsluz(w.OdczytajWyjscie))
}
