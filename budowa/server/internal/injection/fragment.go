// Pakiet injection jest kanałem głównym rozmowy z modelem, powłoką nad
// programem `claude` uruchamianym w trybie strumienia JSON-lines.
package injection

import (
	"danacoconsole/shared"
)

// Fragment jest pozycją strumienia kanału, niosącą fragment kontraktu
// shared.StreamChunkEvent oraz dwa znaczniki toru.
type Fragment struct {
	// Chunk jest fragmentem w kształcie kontraktu.
	Chunk shared.StreamChunkEvent
	// Ostatni odpowiada polu done koperty — po nim strumień się kończy.
	Ostatni bool
	// Tura jest wypełniona wyłącznie we fragmencie kończącym turę i wybudza
	// koordynatora.
	Tura *ZakonczenieTury
}

// RodzajProwenancja oznacza fragment „co poszło do modelu". Wartość
// pochodzi z kontraktu — pakiet jej nie przepisuje. Jest to wartość
// przelotowa wyliczenia ChunkKind: strumień ją niesie, baza jej nie zapisuje.
const RodzajProwenancja = shared.ChunkKindProvenance

// ZakonczenieTury jest podsumowaniem zdarzenia `result` kanału CLI —
// ostatniej linii strumienia. Zamknięcie tury wykonawcy wybudza koordynatora,
// dlatego podsumowanie jedzie osobnym polem, a nie tekstem.
type ZakonczenieTury struct {
	// IdSesjiCLI jest identyfikatorem rozmowy po stronie programu `claude`,
	// podawanym w --resume.
	IdSesjiCLI string `json:"cliSessionId"`
	// Podtyp to `success` albo `error_*` — dosłowna wartość pola subtype.
	Podtyp string `json:"subtype"`
	// Blad powiela pole is_error zdarzenia result.
	Blad bool `json:"isError"`
	// Tekst jest treścią pola result — ostateczną odpowiedzią tury.
	Tekst string `json:"text"`
	// Koszt w dolarach, liczba tur i czas trwania — dane pola result.
	Koszt  float64 `json:"costUsd"`
	Tury   int     `json:"turns"`
	CzasMs int64   `json:"durationMs"`
	// Konto, na którym tura faktycznie się wykonała (kod, nigdy sekret).
	Konto string `json:"account"`
	// TypyZdarzen jest wykazem typów linii przechwyconych w turze, łącznie
	// z liniami nierozpoznanymi.
	TypyZdarzen []string `json:"eventTypes"`
	// LinieNierozpoznane liczy linie wyjścia, których nie dało się odczytać
	// jako JSON.
	LinieNierozpoznane int `json:"unparsedLines"`
}
