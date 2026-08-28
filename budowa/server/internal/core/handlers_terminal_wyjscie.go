// Plik wpina dwie komendy wyjścia obszaru terminal.*: podgląd na żywo i odczyt wyjścia jednego
// procesu, jako rozszerzenie portu Terminal, bo wyjście ma tego samego właściciela co karty i procesy.
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

	// OdczytajWyjscie jest osobną czynnością, bo wykonanie komendy kończy się przy starcie procesu.
	OdczytajWyjscie(ctx context.Context,
		z shared.TerminalOutputReadRequest) (shared.TerminalOutputReadResponse, error)
}

// zarejestrujWyjscieTerminala wpina obie komendy wyjścia modułu Terminal w rejestrze rdzenia platformy.
func zarejestrujWyjscieTerminala(r *Rejestr, w WyjscieTerminala) {
	if r == nil || w == nil {
		return
	}
	r.Zarejestruj(shared.CommandTerminalOutputStream, obsluz(w.StrumienWyjscia))
	r.Zarejestruj(shared.CommandTerminalOutputRead, obsluz(w.OdczytajWyjscie))
}
