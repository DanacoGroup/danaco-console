package injection

import "encoding/json"

// Zaczepy cyklu życia kanału. Konfiguracja zaczepów jedzie do programu `claude`
// sekcją `hooks` pliku ustawień; ten plik domyka drugą połowę drogi — odbiór
// ich zdarzeń ze strumienia.
//
// Kształt kopert linii `system` z przełącznikiem --include-hook-events:
//
//	{"type":"system","subtype":"hook_started","hook_id":"…","hook_name":"…",
//	 "hook_event":"UserPromptSubmit","uuid":"…","session_id":"…"}
//	{"type":"system","subtype":"hook_response","hook_id":"…","hook_name":"…",
//	 "hook_event":"…","output":"…","stdout":"…","stderr":"","exit_code":0,
//	 "outcome":"success","uuid":"…","session_id":"…"}
//
// Zdarzenie zaczepu jest zdarzeniem wykonawczym: mówi, co się stało, nie co
// model powiedział. Kanał go nie interpretuje i nie zamienia na fragment
// kontraktu — oddaje je haczykiem Zapytania (NaZdarzenieZaczepu) warstwie
// składania, która prowadzi dziennik zdarzeń i diagnostykę. Brak haczyka nie
// zmienia przebiegu tury.

// Podtypy linii `system` niosących zdarzenia zaczepów. Wartości są dosłowne
// wartości pola `subtype` strumienia CLI (nazw nie przepisujemy).
const (
	podtypZaczepStart     = "hook_started"
	podtypZaczepOdpowiedz = "hook_response"
)

// ZdarzenieZaczepu jest jednym zdarzeniem zaczepu odczytanym ze strumienia.
type ZdarzenieZaczepu struct {
	// Podtyp to podtypZaczepStart albo podtypZaczepOdpowiedz.
	Podtyp string
	// IdZaczepu wiąże start z odpowiedzią tego samego wywołania (hook_id).
	IdZaczepu string
	// Nazwa jest nazwą zaczepu z konfiguracji (hook_name).
	Nazwa string
	// Zdarzenie jest punktem cyklu życia, w którym zaczep zadziałał
	// (hook_event), np. UserPromptSubmit, PreToolUse.
	Zdarzenie string
	// Wyjscie jest treścią oddaną przez polecenie zaczepu (output). Pole
	// stdout strumienia powiela output — nie niesiemy go drugi raz.
	Wyjscie string
	// BladWyjscia jest treścią stderr polecenia zaczepu.
	BladWyjscia string
	// KodWyjscia jest kodem wyjścia polecenia; nil dla zdarzenia startu,
	// które kodu jeszcze nie ma.
	KodWyjscia *int
	// Wynik jest dosłowną wartością pola outcome (np. success); pusty dla
	// zdarzenia startu.
	Wynik string
	// IdSesjiCLI jest identyfikatorem rozmowy po stronie programu.
	IdSesjiCLI string
	// Surowe jest całą linią strumienia — dowodem pierwotnym zdarzenia.
	// Dziennik zdarzeń zapisuje ją bez przekładu.
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

// zdarzenieZaczepu mówi, czy linia systemowa niesie zdarzenie zaczepu.
func zdarzenieZaczepu(zdarzenie zdarzenieCLI) bool {
	return zdarzenie.Type == TypSystem &&
		(zdarzenie.Subtype == podtypZaczepStart || zdarzenie.Subtype == podtypZaczepOdpowiedz)
}
