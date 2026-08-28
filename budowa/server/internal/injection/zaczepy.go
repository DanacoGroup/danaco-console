// Plik odczytuje zdarzenia zaczepów cyklu życia kanału przesyłane przez
// program claude i składa je w strukturę ZdarzenieZaczepu dla warstwy
// składania odpowiedzi.
package injection

import "encoding/json"

// Podtypy linii `system` niosących zdarzenia zaczepów. Wartości są dosłowne
// wartości pola `subtype` strumienia CLI (nazw nie przepisujemy).
const (
	podtypZaczepStart     = "hook_started"
	podtypZaczepOdpowiedz = "hook_response"
)

// ZdarzenieZaczepu jest jednym zdarzeniem zaczepu odczytanym ze strumienia,
// niosącym podtyp, identyfikator wywołania oraz treść odpowiedzi lub błędu.
type ZdarzenieZaczepu struct {
	// Podtyp to podtypZaczepStart albo podtypZaczepOdpowiedz.
	Podtyp string
	// IdZaczepu wiąże start z odpowiedzią tego samego wywołania (hook_id).
	IdZaczepu string
	// Nazwa jest nazwą zaczepu z konfiguracji (hook_name).
	Nazwa string
	// Zdarzenie jest punktem cyklu życia, w którym zaczep zadziałał, na
	// przykład UserPromptSubmit.
	Zdarzenie string
	// Wyjscie jest treścią oddaną przez polecenie zaczepu. Pole stdout ją
	// powiela i zostaje pominięte.
	Wyjscie string
	// BladWyjscia jest treścią stderr polecenia zaczepu.
	BladWyjscia string
	// KodWyjscia jest kodem wyjścia polecenia; nil dla zdarzenia startu,
	// które kodu jeszcze nie ma.
	KodWyjscia *int
	// Wynik jest dosłowną wartością pola outcome, na przykład success; pusty
	// dla zdarzenia startu.
	Wynik string
	// IdSesjiCLI jest identyfikatorem rozmowy po stronie programu.
	IdSesjiCLI string
	// Surowe jest całą linią strumienia, dowodem zdarzenia zapisanym
	// w dzienniku bez przekładu.
	Surowe json.RawMessage
}

// Odpowiedz mówi, czy zdarzenie jest odpowiedzią zaczepu (ma wynik i kod
// wyjścia), a nie samym startem.
func (z ZdarzenieZaczepu) Odpowiedz() bool {
	return z.Podtyp == podtypZaczepOdpowiedz
}

// Niepowodzenie mówi, czy zaczep odmówił albo zawiódł: kod wyjścia różny od
// zera albo wynik inny niż success. Zdarzenie startu nie jest niepowodzeniem.
func (z ZdarzenieZaczepu) Niepowodzenie() bool {
	if !z.Odpowiedz() {
		return false
	}
	if z.KodWyjscia != nil && *z.KodWyjscia != 0 {
		return true
	}
	return z.Wynik != "" && z.Wynik != "success"
}

// zaczepZeZdarzenia składa zdarzenie zaczepu z odczytanej linii strumienia.
// Surowa linia jedzie w całości — to ona jest dowodem, nie ten przekład.
func zaczepZeZdarzenia(zdarzenie zdarzenieCLI, linia string) ZdarzenieZaczepu {
	return ZdarzenieZaczepu{
		Podtyp:      zdarzenie.Subtype,
		IdZaczepu:   zdarzenie.HookID,
		Nazwa:       zdarzenie.HookName,
		Zdarzenie:   zdarzenie.HookEvent,
		Wyjscie:     zdarzenie.Output,
		BladWyjscia: zdarzenie.Stderr,
		KodWyjscia:  zdarzenie.ExitCode,
		Wynik:       zdarzenie.Outcome,
		IdSesjiCLI:  zdarzenie.SessionID,
		Surowe:      json.RawMessage(linia),
	}
}

// zdarzenieZaczepu mówi, czy linia systemowa strumienia niesie zdarzenie
// zaczepu, sprawdzając podtyp względem stałych zdefiniowanych w pliku.
func zdarzenieZaczepu(zdarzenie zdarzenieCLI) bool {
	return zdarzenie.Type == TypSystem &&
		(zdarzenie.Subtype == podtypZaczepStart || zdarzenie.Subtype == podtypZaczepOdpowiedz)
}
