// Plik składa drogę odczytu nastaw poziomu aplikacja: tych, które opisują sam program, a nie
// treść w nim prowadzoną, czytanych przy każdym nawiązaniu połączenia bez pamięci podręcznej.
package core

import (
	"context"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfig"
)

// NastawyAplikacji jest portem odczytu nastaw poziomu `aplikacja`. Rozszerzenie
// nieobowiązkowe portu Ustawienia — tą samą drogą, którą rdzeń pyta ten sam
// adapter o prowenancję konfiguracji (`ProwenancjaKonfiguracji`).
type NastawyAplikacji interface {
	// WymogLogowania zwraca nastawę wraz z informacją, czy w ogóle jest wskazana — trzy stany, nie dwa.
	WymogLogowania(ctx context.Context) (bool, bool)
}

// WymogLogowania wypełnia port NastawyAplikacji na adapterze ustawień. Kontekst rozstrzygania jest
// pusty z zamysłem, bo programu nie ma czym zawęzić na poziomie aplikacji.
func (a *adapterUstawienOsi) WymogLogowania(ctx context.Context) (bool, bool) {
	if a == nil || a.rozstrzygacz == nil {
		return false, false
	}
	wynik := a.rozstrzygacz.Rozstrzygnij(konfig.Kontekst{KontoOperatora: dane.KontoOperatora(ctx)},
		konfig.KluczWymogLogowania)
	return wartoscWymoguLogowania(wynik.Wartosc)
}

// wartoscWymoguLogowania przekłada zapis kolumny na parę wartość i wskazanie. Zapis pusty i zapis
// nieczytelny oba znaczą brak wskazania, bo zgadywanie sensu byłoby rozstrzyganiem za Operatora.
func wartoscWymoguLogowania(zapis string) (bool, bool) {
	switch strings.ToLower(strings.TrimSpace(zapis)) {
	case "true", "1", "tak":
		return true, true
	case "false", "0", "nie":
		return false, true
	default:
		return false, false
	}
}
