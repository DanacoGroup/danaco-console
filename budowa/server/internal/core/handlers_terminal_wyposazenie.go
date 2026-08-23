// Wpięcie dwudziestu jeden komend obszaru `terminal.*` spoza rdzenia
// wykonawczego modułu: cyklu życia karty, odczytu pliku, wstrzymania procesu,
// książki hostów, biblioteki skryptów, wykazu kluczy, tuneli i obserwacji.
//
// Port jest ROZSZERZENIEM, nie drugim portem — tak samo jak wyjście
// (`handlers_terminal_wyjscie.go`): osadza `Terminal`, bo wyposażenie modułu ma
// tego samego właściciela co karty i procesy, a adapter wypełniający jeden
// wypełnia wszystkie.
//
// Żadna z tych komend nie rozgłasza własnego zdarzenia i nie jest to
// przeoczenie: kontrakt daje modułowi jedno zdarzenie — `terminal.process.changed`
// — i dotyczy ono procesu. Zmiany, które procesu dotyczą (wstrzymanie,
// zakończenie procesów zamykanej karty, wyzwolenie obserwacji), idą tym
// zdarzeniem z adaptera. Zmiany wpisu książki hostów czy pozycji biblioteki nie
// mają w kontrakcie zdarzenia, więc klient odświeża wykaz po własnym zapisie —
// zamiast dostawać rozgłoszenie nazwą, której kontrakt nie zna.
package core

import (
	"context"

	"danacoconsole/shared"
)

// WyposazenieTerminala jest portem wyposażenia modułu Terminal.
type WyposazenieTerminala interface {
	Terminal

	// ── Cykl życia karty i jej zawartość ────────────────────────────────────
	ZamknijKarte(ctx context.Context, z shared.TerminalSessionCloseRequest) (shared.TerminalSessionCloseResponse, error)
	WykazKart(ctx context.Context, z shared.TerminalSessionListRequest) (shared.TerminalSessionListResponse, error)
	OdczytajPlik(ctx context.Context, z shared.TerminalFileReadRequest) (shared.TerminalFileReadResponse, error)
	WstrzymajProces(ctx context.Context, z shared.TerminalProcessSuspendRequest) (shared.TerminalProcessSuspendResponse, error)

	// ── Książka hostów ──────────────────────────────────────────────────────
	ZapiszHosta(ctx context.Context, z shared.TerminalHostSaveRequest) (shared.TerminalHostSaveResponse, error)
	WykazHostow(ctx context.Context, z shared.TerminalHostListRequest) (shared.TerminalHostListResponse, error)
	UsunHosta(ctx context.Context, z shared.TerminalHostRemoveRequest) (shared.TerminalHostRemoveResponse, error)

	// ── Biblioteka skryptów ─────────────────────────────────────────────────
	ZapiszSkrypt(ctx context.Context, z shared.TerminalScriptSaveRequest) (shared.TerminalScriptSaveResponse, error)
	WykazSkryptow(ctx context.Context, z shared.TerminalScriptListRequest) (shared.TerminalScriptListResponse, error)
	UsunSkrypt(ctx context.Context, z shared.TerminalScriptRemoveRequest) (shared.TerminalScriptRemoveResponse, error)
	SprawdzSkrypt(ctx context.Context, z shared.TerminalScriptLintRequest) (shared.TerminalScriptLintResponse, error)

	// ── Przekierowania portów ───────────────────────────────────────────────
	OtworzTunel(ctx context.Context, z shared.TerminalTunnelOpenRequest) (shared.TerminalTunnelOpenResponse, error)
	WykazTuneli(ctx context.Context, z shared.TerminalTunnelListRequest) (shared.TerminalTunnelListResponse, error)
	ZamknijTunel(ctx context.Context, z shared.TerminalTunnelCloseRequest) (shared.TerminalTunnelCloseResponse, error)

	// ── Klucze SSH ──────────────────────────────────────────────────────────
	WytworzKlucz(ctx context.Context, z shared.TerminalKeyGenerateRequest) (shared.TerminalKeyGenerateResponse, error)
	WciagnijKlucz(ctx context.Context, z shared.TerminalKeyImportRequest) (shared.TerminalKeyImportResponse, error)
	WykazKluczy(ctx context.Context, z shared.TerminalKeyListRequest) (shared.TerminalKeyListResponse, error)
	UsunKlucz(ctx context.Context, z shared.TerminalKeyRemoveRequest) (shared.TerminalKeyRemoveResponse, error)

	// ── Obserwacje plików ───────────────────────────────────────────────────
	ZalozObserwacje(ctx context.Context, z shared.TerminalWatchStartRequest) (shared.TerminalWatchStartResponse, error)
	ZatrzymajObserwacje(ctx context.Context, z shared.TerminalWatchStopRequest) (shared.TerminalWatchStopResponse, error)
	WykazObserwacji(ctx context.Context, z shared.TerminalWatchListRequest) (shared.TerminalWatchListResponse, error)
}

// Zgodność adaptera z portem wyposażenia sprawdzana jest przy kompilacji.
var _ WyposazenieTerminala = (*adapterTerminala)(nil)

// zarejestrujWyposazenieTerminala wpina wszystkie komendy wyposażenia modułu.
func zarejestrujWyposazenieTerminala(r *Rejestr, w WyposazenieTerminala) {
	if r == nil || w == nil {
		return
	}
	r.Zarejestruj(shared.CommandTerminalSessionClose, obsluz(w.ZamknijKarte))
	r.Zarejestruj(shared.CommandTerminalSessionList, obsluz(w.WykazKart))
	r.Zarejestruj(shared.CommandTerminalFileRead, obsluz(w.OdczytajPlik))
	r.Zarejestruj(shared.CommandTerminalProcessSuspend, obsluz(w.WstrzymajProces))

	r.Zarejestruj(shared.CommandTerminalHostSave, obsluz(w.ZapiszHosta))
	r.Zarejestruj(shared.CommandTerminalHostList, obsluz(w.WykazHostow))
	r.Zarejestruj(shared.CommandTerminalHostRemove, obsluz(w.UsunHosta))

	r.Zarejestruj(shared.CommandTerminalScriptSave, obsluz(w.ZapiszSkrypt))
	r.Zarejestruj(shared.CommandTerminalScriptList, obsluz(w.WykazSkryptow))
	r.Zarejestruj(shared.CommandTerminalScriptRemove, obsluz(w.UsunSkrypt))
	r.Zarejestruj(shared.CommandTerminalScriptLint, obsluz(w.SprawdzSkrypt))

	r.Zarejestruj(shared.CommandTerminalTunnelOpen, obsluz(w.OtworzTunel))
	r.Zarejestruj(shared.CommandTerminalTunnelList, obsluz(w.WykazTuneli))
	r.Zarejestruj(shared.CommandTerminalTunnelClose, obsluz(w.ZamknijTunel))

	r.Zarejestruj(shared.CommandTerminalKeyGenerate, obsluz(w.WytworzKlucz))
	r.Zarejestruj(shared.CommandTerminalKeyImport, obsluz(w.WciagnijKlucz))
	r.Zarejestruj(shared.CommandTerminalKeyList, obsluz(w.WykazKluczy))
	r.Zarejestruj(shared.CommandTerminalKeyRemove, obsluz(w.UsunKlucz))

	r.Zarejestruj(shared.CommandTerminalWatchStart, obsluz(w.ZalozObserwacje))
	r.Zarejestruj(shared.CommandTerminalWatchStop, obsluz(w.ZatrzymajObserwacje))
	r.Zarejestruj(shared.CommandTerminalWatchList, obsluz(w.WykazObserwacji))
}
