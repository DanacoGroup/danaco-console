package core

import (
	"context"
	"errors"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// Utrwalenie wyboru eksperta w wierszu okna — drugie ogniwo, nie druga prawda.
//
// Wybór eksperta żyje w rejestrze pamięciowym nadzorcy (`session.Okno.Agent`)
// i w kolumnie `agent_kod` wiersza okna. Wiersz zakładany jest leniwie i
// wypełnia kolumnę tylko przy założeniu, więc bez tego zapisu wskazanie
// eksperta oknu mającemu już wiersz nie przeżyłoby restartu — pamięć odtwarza
// się wtedy z wierszy (`odtworzenie_stanu.go`).
//
// Zapis idzie tym samym wzorem, co utrwalenie kanału (`utrwalKanalOkna`,
// adapter_modul_model.go): odczyt wiersza po identyfikatorze, porównanie, zapis
// wyłącznie przy różnicy. Ekspert stoi obok kanału i utrwalany jest obok niego,
// nie zamiast.
//
// Nie sprawdza, czy ekspert o wskazanym kodzie istnieje: kolumna `agent_kod`
// nie ma więzu obcego z rozmysłem, a kod nierozpoznany jest faktem czytelnym —
// składacz nakładki nie znajduje eksperta i rusza z samą osią. Nie zakłada
// wiersza okna. Nie dotyka ani nakładki, ani trybu silnika — to robi
// `nakladkaZAgentem`.

// ZeStrazaEksperta wpina straż zakresu eksperta. Bez niej nałożenie idzie bez
// sprawdzenia — stanem wyjściowym platformy jest pełny dostęp.
func (a *adapterOkien) ZeStrazaEksperta(straz StrazEksperta) *adapterOkien {
	a.straz = straz
	return a
}

// utrwalAgentaOkna zapisuje wybór eksperta w istniejącym wierszu okna.
//
// Wskaźnik pusty znaczy „żądanie nie ruszało eksperta" i nie robi nic —
// dokładnie tak, jak rozumie go `session.Zmiana`. Wskaźnik na pusty napis
// zdejmuje eksperta: kolumna wraca do NULL, czyli do modelu surowego.
//
// Brak wiersza nie jest błędem. Błąd zapisu jest błędem komendy —
// wiersz istnieje, a nie przyjął wyboru, więc wybór nie przeżyje restartu.
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
	// Zakres eksperta czytany jest TUTAJ, bo tutaj ekspert wchodzi do okna.
	// Zawężenie zapisane w Permissions Center, którego nikt by w tym miejscu
	// nie sprawdził, byłoby napisem w oknie konfiguracji i niczym więcej
	// (`straz_eksperta.go`).
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

// bladZapisuAgenta składa odmowę zapisu wyboru eksperta.
func bladZapisuAgenta(powod string, err error) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"core: wybór eksperta okna — "+powod+": "+err.Error()))
}
