// Plik obsługuje launcher.hotkey.get i launcher.hotkey.set: skrót globalny otwierający wywoływacz poleceń; rdzeń trzyma nastawę, powłoka programu okiennego rejestruje skrót u siebie.
package core

import (
	"context"
	"strings"
	"sync"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

const (
	// kluczSkrotuWywolywacza to nastawa niosąca zapis skrótu globalnego w konfiguracji tego Operatora rdzenia.
	kluczSkrotuWywolywacza = "wywolywacz_skrot_globalny"

	// zdolnoscSkrotuGlobalnego to nazwa zdolności deklarowanej przez powłokę
	// w powitaniu. Powłoka okienna ją deklaruje, przeglądarka nie.
	zdolnoscSkrotuGlobalnego = "launcher.hotkey"
)

// zdolnosciKlientow zapamiętuje, co klienci zadeklarowali w powitaniu; rejestr jest w pamięci i żyje tyle, co rdzeń, bo deklaracja dotyczy połączenia, nie Operatora.
type zdolnosciKlientow struct {
	zamek sync.RWMutex
	wpisy map[string][]string
}

// noweZdolnosciKlientow zakłada pusty rejestr deklaracji zdolności klientów przy starcie tego rdzenia.
func noweZdolnosciKlientow() *zdolnosciKlientow {
	return &zdolnosciKlientow{wpisy: map[string][]string{}}
}

// zapamietaj odkłada deklarację jednego klienta w rejestrze zdolności tego bieżącego połączenia rdzenia.
func (z *zdolnosciKlientow) zapamietaj(klient string, zdolnosci []string) {
	if z == nil || strings.TrimSpace(klient) == "" {
		return
	}
	z.zamek.Lock()
	defer z.zamek.Unlock()
	if len(zdolnosci) == 0 {
		delete(z.wpisy, klient)
		return
	}
	kopia := make([]string, len(zdolnosci))
	copy(kopia, zdolnosci)
	z.wpisy[klient] = kopia
}

// ktokolwiekDeklaruje mówi, czy którykolwiek ze znanych klientów zadeklarował wskazaną zdolność w powitaniu.
func (z *zdolnosciKlientow) ktokolwiekDeklaruje(zdolnosc string) bool {
	if z == nil {
		return false
	}
	z.zamek.RLock()
	defer z.zamek.RUnlock()
	for _, deklaracja := range z.wpisy {
		for _, pozycja := range deklaracja {
			if strings.EqualFold(strings.TrimSpace(pozycja), zdolnosc) {
				return true
			}
		}
	}
	return false
}

// adapterWywolywacza wypełnia port Wywolywacz magazynem nastaw i rejestrem deklaracji tych klientów rdzenia.
type adapterWywolywacza struct {
	konfiguracja dane.RepozytoriumKonfiguracji
	rozstrzygacz *konfig.Rozstrzygacz
	zdolnosci    *zdolnosciKlientow
}

// nowyAdapterWywolywacza wiąże port z magazynem nastaw i rejestrem deklaracji klientów przy tym montażu.
func nowyAdapterWywolywacza(konfiguracja dane.RepozytoriumKonfiguracji,
	rozstrzygacz *konfig.Rozstrzygacz, zdolnosci *zdolnosciKlientow) *adapterWywolywacza {

	return &adapterWywolywacza{
		konfiguracja: konfiguracja, rozstrzygacz: rozstrzygacz, zdolnosci: zdolnosci,
	}
}

// SkrotWywolywacza obsługuje launcher.hotkey.get, oddając nastawę skrótu wraz z jej pełną wykonalnością.
func (a *adapterWywolywacza) SkrotWywolywacza(_ context.Context,
	_ shared.LauncherHotkeyGetRequest) (shared.LauncherHotkeyGetResponse, error) {

	skrot := a.odczytajSkrot()
	wspierany := a.zdolnosci.ktokolwiekDeklaruje(zdolnoscSkrotuGlobalnego)

	odpowiedz := shared.LauncherHotkeyGetResponse{
		Hotkey:     skrot,
		Supported:  wspierany,
		Registered: wspierany && skrot != "",
	}
	switch {
	case !wspierany:
		powod := "skrót globalny rejestruje powłoka programu okiennego, a żadne " +
			"z połączonych okien nie zadeklarowało tej zdolności w powitaniu " +
			"(capabilities: " + zdolnoscSkrotuGlobalnego + "); w przeglądarce skrótu " +
			"globalnego nie ma jak przechwycić — naprawa: otworzyć Danaco Console " +
			"w programie okiennym"
		odpowiedz.Reason = &powod
	case skrot == "":
		powod := "nastawa skrótu jest pusta, więc nie ma czego rejestrować; " +
			"naprawa: zapisać skrót komendą launcher.hotkey.set"
		odpowiedz.Reason = &powod
	}
	return odpowiedz, nil
}

// ZapiszSkrotWywolywacza obsługuje launcher.hotkey.set, zapisując nastawę niezależnie od jej rejestracji.
func (a *adapterWywolywacza) ZapiszSkrotWywolywacza(ctx context.Context,
	z shared.LauncherHotkeySetRequest) (shared.LauncherHotkeySetResponse, error) {

	if a.konfiguracja == nil {
		return shared.LauncherHotkeySetResponse{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeInternalError,
			"wywoływacz: rdzeń nie ma wpiętego magazynu konfiguracji — nastawy skrótu "+
				"nie ma gdzie zapisać"))
	}
	skrot := strings.TrimSpace(z.Hotkey)
	if skrot != "" {
		if err := sprawdzZapisSkrotu(skrot); err != nil {
			return shared.LauncherHotkeySetResponse{}, err
		}
	}
	poziom := shared.ConfigScopeGlobal
	if z.Scope != nil && strings.TrimSpace(string(*z.Scope)) != "" {
		poziom = string(*z.Scope)
	}
	if err := sprawdzPoziomKontekstu(shared.ConfigScope(poziom)); err != nil {
		return shared.LauncherHotkeySetResponse{}, err
	}

	kopia := skrot
	err := a.konfiguracja.Ustaw(ctx, dane.Ustawienie{
		Poziom: shared.ConfigScope(poziom), Klucz: kluczSkrotuWywolywacza, Wartosc: &kopia,
	})
	if err != nil {
		return shared.LauncherHotkeySetResponse{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeInternalError, "wywoływacz: nie można zapisać skrótu: "+err.Error()))
	}
	if a.rozstrzygacz != nil {
		a.rozstrzygacz.Oglos(kluczSkrotuWywolywacza)
	}

	wspierany := a.zdolnosci.ktokolwiekDeklaruje(zdolnoscSkrotuGlobalnego)
	odpowiedz := shared.LauncherHotkeySetResponse{
		Hotkey:     skrot,
		Registered: wspierany && skrot != "",
	}
	if !odpowiedz.Registered {
		powod := "zapis się powiódł, ale rejestracji nie ma kto wykonać: skrót globalny " +
			"przechwytuje powłoka programu okiennego, a żadne z połączonych okien nie " +
			"zadeklarowało tej zdolności"
		if skrot == "" {
			powod = "zapis się powiódł: pusty skrót zdejmuje rejestrację"
		}
		odpowiedz.Reason = &powod
	}
	return odpowiedz, nil
}

// odczytajSkrot oddaje nastawę obowiązującą po rozstrzygnięciu poziomów zasięgu tej całej konfiguracji.
func (a *adapterWywolywacza) odczytajSkrot() string {
	if a.rozstrzygacz == nil {
		return ""
	}
	return strings.TrimSpace(
		a.rozstrzygacz.Rozstrzygnij(konfig.Kontekst{}, kluczSkrotuWywolywacza).Wartosc)
}

// sprawdzZapisSkrotu odbija zapis, którego powłoka nie zrozumie: odrzuca zapis bez klawisza głównego albo z samymi modyfikatorami, a resztę przepuszcza, bo to powłoka jest tu autorytetem.
func sprawdzZapisSkrotu(skrot string) error {
	czlony := strings.Split(skrot, "+")
	modyfikatory := map[string]struct{}{
		"ctrl": {}, "control": {}, "alt": {}, "shift": {}, "meta": {},
		"cmd": {}, "command": {}, "super": {}, "win": {},
	}
	for _, czlon := range czlony {
		oczyszczony := strings.ToLower(strings.TrimSpace(czlon))
		if oczyszczony == "" {
			return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
				"wywoływacz: zapis skrótu „"+skrot+"” ma pusty człon — człony rozdziela się plusem"))
		}
		if _, jest := modyfikatory[oczyszczony]; !jest {
			return nil
		}
	}
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"wywoływacz: zapis skrótu „"+skrot+"” niesie same modyfikatory — brakuje klawisza "+
			"głównego, np. Ctrl+Shift+Space"))
}
