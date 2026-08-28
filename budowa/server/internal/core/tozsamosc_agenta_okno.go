package core

import (
	"context"
	"errors"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// Utrwalenie wyboru eksperta w wierszu okna: drugie ogniwo obok pamięci nadzorcy, nie druga prawda.

// ZeStrazaEksperta wpina straż zakresu eksperta. Bez niej nałożenie idzie bez
// sprawdzenia — stanem wyjściowym platformy jest pełny dostęp.
func (a *adapterOkien) ZeStrazaEksperta(straz StrazEksperta) *adapterOkien {
	a.straz = straz
	return a
}

// utrwalAgentaOkna zapisuje wybór eksperta w istniejącym wierszu okna. Wskaźnik pusty nie robi nic, wskaźnik na pusty napis zdejmuje eksperta. Brak wiersza nie jest błędem, ale błąd zapisu jest błędem komendy, bo wybór nie przeżyje restartu.
func (a *adapterOkien) utrwalAgentaOkna(ctx context.Context, idOkna string, agent *string) error {
	if a == nil || agent == nil || idOkna == "" {
		return nil
	}
	if a.trwalosc == nil || a.trwalosc.okna == nil {
		return nil
	}
	wiersz, err := a.trwalosc.okna.PoIdentyfikatorze(ctx, idOkna)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return nil
	}
	if err != nil {
		return bladZapisuAgenta("nie da się odczytać wiersza okna "+idOkna, err)
	}
	kod := *agent
	if wartoscTekstu(wiersz.AgentKod) == kod {
		return nil
	}
	// Zakres eksperta czytany jest tutaj, bo tutaj ekspert wchodzi do okna.
	if a.straz != nil && kod != "" {
		if err := a.straz.SprawdzModulOkna(ctx, kod, wiersz.ModulID); err != nil {
			return err
		}
	}
	wiersz.AgentKod = kodAgentaKolumny(kod)
	if err := a.trwalosc.okna.Aktualizuj(ctx, wiersz); err != nil {
		return bladZapisuAgenta("nie da się zapisać eksperta okna "+idOkna, err)
	}
	return nil
}

// kodAgentaKolumny przekłada kod eksperta na wartość kolumny `agent_kod`.
//
// NULL znaczy model surowy i nie jest tym samym, co pusty napis. Pusty napis
// wraca więc jako brak wskaźnika, żeby w bazie nie
// powstał drugi sposób powiedzenia tej samej rzeczy.
func kodAgentaKolumny(kod string) *string {
	if kod == "" {
		return nil
	}
	return &kod
}

// bladZapisuAgenta składa odmowę zapisu wyboru eksperta, niosącą przekazany powód i przyczynę źródłową.
func bladZapisuAgenta(powod string, err error) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"core: wybór eksperta okna — "+powod+": "+err.Error()))
}
