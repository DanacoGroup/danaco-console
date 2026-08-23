package injection

import (
	"encoding/json"

	"danacoconsole/shared"
)

// Typy linii wyjściowych programu `claude` przechwytywane przez kanał.
// Katalog jest opisowy: linia typu spoza katalogu nie przerywa strumienia,
// tylko nie zamienia się we fragment.
const (
	TypSystem     = "system"
	TypAssistant  = "assistant"
	TypUser       = "user"
	TypResult     = "result"
	TypLimitTempa = "rate_limit_event"
)

// zdarzenieCLI jest jedną linią wyjścia programu. Odczytujemy wyłącznie pola
// niosące treść dla kontraktu; reszta linii jedzie dalej nietknięta.
type zdarzenieCLI struct {
	Type      string          `json:"type"`
	Subtype   string          `json:"subtype"`
	SessionID string          `json:"session_id"`
	Message   *wiadomoscCLI   `json:"message"`
	IsError   bool            `json:"is_error"`
	Result    string          `json:"result"`
	CostUSD   float64         `json:"total_cost_usd"`
	NumTurns  int             `json:"num_turns"`
	DurationM int64           `json:"duration_ms"`
	RateLimit *limitTempaCLI  `json:"rate_limit_info"`
	Content   []blokCLI       `json:"content"`
	Usage     json.RawMessage `json:"usage"`

	// Pola zdarzeń zaczepów — linie system/hook_started i system/hook_response;
	// ich przekład na zdarzenie kontraktu leży w zaczepy.go. Pole `stdout`
	// strumienia powiela `output` i nie jest odczytywane.
	HookID    string `json:"hook_id"`
	HookName  string `json:"hook_name"`
	HookEvent string `json:"hook_event"`
	Output    string `json:"output"`
	Stderr    string `json:"stderr"`
	ExitCode  *int   `json:"exit_code"`
	Outcome   string `json:"outcome"`
}

// wiadomoscCLI jest wiadomością modelu albo użytkownika. Pole content bywa
// napisem albo tablicą bloków, dlatego trzymamy je surowe.
type wiadomoscCLI struct {
	Role    string          `json:"role"`
	Model   string          `json:"model"`
	Content json.RawMessage `json:"content"`
}

// blokCLI jest pojedynczym blokiem treści wiadomości.
type blokCLI struct {
	Type      string          `json:"type"`
	Text      string          `json:"text"`
	Thinking  string          `json:"thinking"`
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Input     json.RawMessage `json:"input"`
	ToolUseID string          `json:"tool_use_id"`
	Content   json.RawMessage `json:"content"`
	IsError   bool            `json:"is_error"`
	Source    json.RawMessage `json:"source"`
}

// limitTempaCLI niesie stan limitu konta. To on, a nie treść komunikatu
// o błędzie, jest podstawą rotacji kont.
type limitTempaCLI struct {
	Status   string `json:"status"`
	ResetsAt int64  `json:"resetsAt"`
	Rodzaj   string `json:"rateLimitType"`
}

// bloki odczytuje treść wiadomości niezależnie od tego, czy przyszła napisem,
// czy tablicą bloków.
func (w *wiadomoscCLI) bloki() []blokCLI {
	if w == nil || len(w.Content) == 0 {
		return nil
	}
	var tablica []blokCLI
	if err := json.Unmarshal(w.Content, &tablica); err == nil {
		return tablica
	}
	var tekst string
	if err := json.Unmarshal(w.Content, &tekst); err == nil && tekst != "" {
		return []blokCLI{{Type: "text", Text: tekst}}
	}
	return nil
}

// fragment zamienia blok treści na fragment kontraktu. Blok rodzaju, którego
// kontrakt nie zna, nie daje fragmentu i nie przerywa strumienia.
func (b blokCLI) fragment(okno, wiadomosc string) (shared.StreamChunkEvent, bool) {
	switch b.Type {
	case "text":
		if b.Text == "" {
			return shared.StreamChunkEvent{}, false
		}
		return tekstowy(okno, wiadomosc, shared.ChunkKindText, b.Text), true
	case "thinking", "redacted_thinking":
		if b.Thinking == "" {
			return shared.StreamChunkEvent{}, false
		}
		return tekstowy(okno, wiadomosc, shared.ChunkKindThinking, b.Thinking), true
	case "tool_use":
		return daneowy(okno, wiadomosc, shared.ChunkKindToolUse, map[string]any{
			"id": b.ID, "name": b.Name, "input": b.Input,
		}), true
	case "tool_result":
		return daneowy(okno, wiadomosc, shared.ChunkKindToolResult, map[string]any{
			"toolUseId": b.ToolUseID, "content": b.Content, "isError": b.IsError,
		}), true
	case "image":
		return daneowy(okno, wiadomosc, shared.ChunkKindImage, map[string]any{
			"source": b.Source,
		}), true
	default:
		return shared.StreamChunkEvent{}, false
	}
}

func tekstowy(okno, wiadomosc string, rodzaj shared.ChunkKind, tresc string) shared.StreamChunkEvent {
	return shared.StreamChunkEvent{
		WindowId: okno, MessageId: wiadomosc, Kind: rodzaj, Text: &tresc,
	}
}

// daneowy pakuje treść nietekstową. Błąd kodowania zostawia fragment bez
// danych — rodzaj i powiązanie z oknem są ważniejsze niż ładunek.
func daneowy(okno, wiadomosc string, rodzaj shared.ChunkKind, dane any) shared.StreamChunkEvent {
	fragment := shared.StreamChunkEvent{WindowId: okno, MessageId: wiadomosc, Kind: rodzaj}
	if surowe, err := json.Marshal(dane); err == nil {
		fragment.Data = surowe
	}
	return fragment
}
