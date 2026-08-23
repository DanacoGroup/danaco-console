// Odpowiedzialność pliku: `launcher.hotkey.get` i `launcher.hotkey.set` —
// skrót globalny otwierający wywoływacz poleceń.
//
// ── Kto naprawdę rejestruje skrót globalny ──────────────────────────────────
// Skrót globalny przechwytuje powłoka programu okiennego na maszynie Operatora
// — nie rdzeń. Rdzeń stoi na serwerze i klawiatury tamtej maszyny nie widzi.
// Podział jest więc taki:
//
//   - Rdzeń TRZYMA nastawę. Dzięki temu skrót jest ten sam na każdej maszynie
//     tego samego Operatora i przeżywa ponowne zainstalowanie okna.
//   - Powłoka REJESTRUJE skrót u siebie i to ona wie, czy się udało; skrót
//     zajęty przez inny program zajmie go dalej, cokolwiek rdzeń o tym sądzi.
//
// ── Uczciwy stan wykonalności ───────────────────────────────────────────────
// `supported` mówi, czy po drugiej stronie stoi powłoka, która w ogóle umie
// zarejestrować skrót globalny. Rdzeń wie to z jednego miejsca: z powitania.
// Klient deklaruje w nim swoje zdolności (`connection.hello`, pole
// `capabilities`), a rdzeń zapamiętuje deklarację. Przeglądarka takiej
// zdolności nie zadeklaruje i wtedy odpowiedź mówi wprost, że skrótu nie ma
// kto przechwycić — zamiast obiecywać skrót, który nikogo nie obudzi.
//
// `registered` nie jest zgadywane: jest prawdą wtedy i tylko wtedy, gdy skrót
// jest niepusty ORAZ stoi powłoka deklarująca zdolność. Rdzeń nie twierdzi, że
// rejestracja się powiodła, gdy nie ma komu jej wykonać.
//
// ── Zapis zostaje nawet wtedy, gdy rejestracja się nie uda ──────────────────
// Kontrakt mówi to wprost i tak jest tutaj: zapis nastawy idzie pierwszy,
// a odpowiedź mówi osobno o zapisie i osobno o rejestracji. Skrót zajęty przez
// inny program nie jest błędem zapisu — Operator zwolni go później i nie będzie
// musiał wpisywać nastawy od nowa.
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
	// kluczSkrotuWywolywacza to nastawa niosąca zapis skrótu globalnego.
	kluczSkrotuWywolywacza = "wywolywacz_skrot_globalny"

	// zdolnoscSkrotuGlobalnego to nazwa zdolności deklarowanej przez powłokę
	// w powitaniu. Powłoka okienna ją deklaruje, przeglądarka nie.
	zdolnoscSkrotuGlobalnego = "launcher.hotkey"
)

// zdolnosciKlientow zapamiętuje, co klienci zadeklarowali w powitaniu.
//
// Rejestr jest w pamięci i żyje tyle, co rdzeń: deklaracja dotyczy połączenia,
// a nie Operatora. Wiersz w bazie przeżyłby zamknięcie okna i twierdziłby, że
// skrót globalny ma kto przechwycić, gdy po tamtej stronie nie ma już nikogo.
type zdolnosciKlientow struct {
	zamek sync.RWMutex
	wpisy map[string][]string
}

// noweZdolnosciKlientow zakłada pusty rejestr deklaracji.
func noweZdolnosciKlientow() *zdolnosciKlientow {
	return &zdolnosciKlientow{wpisy: map[string][]string{}}
}

// zapamietaj odkłada deklarację jednego klienta.
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

// ktokolwiekDeklaruje mówi, czy którykolwiek ze znanych klientów zadeklarował
// wskazaną zdolność.
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

// adapterWywolywacza wypełnia port `Wywolywacz`.
type adapterWywolywacza struct {
	konfiguracja dane.RepozytoriumKonfiguracji
	rozstrzygacz *konfig.Rozstrzygacz
	zdolnosci    *zdolnosciKlientow
}

// nowyAdapterWywolywacza wiąże port z magazynem nastaw i rejestrem deklaracji.
func nowyAdapterWywolywacza(konfiguracja dane.RepozytoriumKonfiguracji,
	rozstrzygacz *konfig.Rozstrzygacz, zdolnosci *zdolnosciKlientow) *adapterWywolywacza {

	return &adapterWywolywacza{
		konfiguracja: konfiguracja, rozstrzygacz: rozstrzygacz, zdolnosci: zdolnosci,
	}
}

// SkrotWywolywacza obsługuje `launcher.hotkey.get`.
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

// ZapiszSkrotWywolywacza obsługuje `launcher.hotkey.set`.
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

// odczytajSkrot oddaje nastawę obowiązującą po rozstrzygnięciu poziomów.
func (a *adapterWywolywacza) odczytajSkrot() string {
	if a.rozstrzygacz == nil {
		return ""
	}
	return strings.TrimSpace(
		a.rozstrzygacz.Rozstrzygnij(konfig.Kontekst{}, kluczSkrotuWywolywacza).Wartosc)
}

// sprawdzZapisSkrotu odbija zapis, którego powłoka nie zrozumie.
//
// Sprawdzenie jest celowo zachowawcze: rdzeń nie zna wykazu klawiszy każdej
// powłoki i nie będzie go zgadywał. Odrzuca to, co na pewno jest błędem —
// zapis bez klawisza głównego albo z samymi modyfikatorami — a resztę
// przepuszcza, bo to powłoka jest tu autorytetem.
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
