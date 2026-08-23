package core

import (
	"encoding/json"

	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// parametryWykonania są kształtem pola executionParams kompletu kontekstu.
// Kontrakt zostawia to pole surowym JSON-em, bo nie każdy moduł
// wykonuje pracę tak samo; okno komunikacji ma jednak ustalony zestaw
// parametrów wykonania i to on tędy jedzie.
//
// Nazwy pól JSON są nazwami z okna kontraktu, więc parametry odczytane
// z przeniesienia dają się nałożyć na okno bez drugiego słownika.
type parametryWykonania struct {
	KanalModelu         string                `json:"modelChannelId,omitempty"`
	KatalogiRobocze     []string              `json:"workingDirs,omitempty"`
	SrodowiskoWykonania shared.ExecutionEnv   `json:"executionEnv,omitempty"`
	TrybUprawnien       shared.PermissionMode `json:"permissionMode,omitempty"`
	RolaOkna            shared.WindowRole     `json:"windowRole,omitempty"`
	// Agent jedzie razem z kanałem, bo przekazanie kontekstu przenosi całe
	// stanowisko, nie samą pracę. Okno docelowe ma pracować tą samą
	// tożsamością — inaczej ta sama wypowiedź trafiłaby
	// do modelu z innym promptem systemowym i innym modelem bazowym, a Operator
	// zobaczyłby rozjazd bez przyczyny.
	//
	// Nazwa pola jest nazwą z okna kontraktu, tak jak reszta — parametry
	// odczytane z przeniesienia nakładają się na okno bez drugiego słownika.
	Agent string `json:"agentId,omitempty"`
}

// parametryZOkna spisuje parametry wykonania okna źródłowego. Niepowodzenie
// kodowania daje pole puste — przeniesienie kontekstu ma się odbyć także wtedy.
func parametryZOkna(o session.Okno) json.RawMessage {
	tresc, err := json.Marshal(parametryWykonania{
		KanalModelu:         o.KanalModelu,
		KatalogiRobocze:     o.KatalogiRobocze,
		SrodowiskoWykonania: o.SrodowiskoWykonania,
		TrybUprawnien:       o.TrybUprawnien,
		RolaOkna:            o.RolaOkna,
		Agent:               o.Agent,
	})
	if err != nil {
		return nil
	}
	return tresc
}

// parametryKompletu odczytuje parametry wykonania z kompletu. Pole puste
// i pole nieczytelne dają parametry puste, a nie błąd przeniesienia.
func parametryKompletu(komplet shared.ContextBundle) parametryWykonania {
	var parametry parametryWykonania
	if len(komplet.ExecutionParams) == 0 {
		return parametry
	}
	if err := json.Unmarshal(komplet.ExecutionParams, &parametry); err != nil {
		return parametryWykonania{}
	}
	return parametry
}

// nalozParametry nanosi parametry kompletu na ustawienia okna zakładanego.
// Parametr niewskazany zostawia wartość odziedziczoną po oknie źródłowym.
func nalozParametry(u session.Ustawienia, komplet shared.ContextBundle) session.Ustawienia {
	parametry := parametryKompletu(komplet)
	if parametry.KanalModelu != "" {
		u.KanalModelu = parametry.KanalModelu
	}
	if len(parametry.KatalogiRobocze) > 0 {
		u.KatalogiRobocze = append([]string(nil), parametry.KatalogiRobocze...)
	}
	if parametry.SrodowiskoWykonania != "" {
		u.SrodowiskoWykonania = parametry.SrodowiskoWykonania
	}
	if parametry.TrybUprawnien != "" {
		u.TrybUprawnien = parametry.TrybUprawnien
	}
	if parametry.RolaOkna != "" {
		u.RolaOkna = parametry.RolaOkna
	}
	if parametry.Agent != "" {
		u.Agent = parametry.Agent
	}
	return u
}

// zmianaCelu składa wybiórczą zmianę okna docelowego wskazanego w żądaniu:
// moduł docelowy oraz te parametry wykonania, które komplet naprawdę niesie.
// Parametr nieprzeniesiony zostaje w oknie docelowym bez zmiany.
//
// Roli okna przeniesienie nie rusza: rola wiąże okno z koordynatorem pętli,
// a przekazanie kontekstu nie jest zmianą układu pętli.
func zmianaCelu(modul string, komplet shared.ContextBundle) session.Zmiana {
	parametry := parametryKompletu(komplet)
	zmiana := session.Zmiana{Modul: &modul}
	if parametry.KanalModelu != "" {
		kanal := parametry.KanalModelu
		zmiana.KanalModelu = &kanal
	}
	if len(parametry.KatalogiRobocze) > 0 {
		zmiana.KatalogiRobocze = append([]string(nil), parametry.KatalogiRobocze...)
	}
	if parametry.SrodowiskoWykonania != "" {
		srodowisko := parametry.SrodowiskoWykonania
		zmiana.SrodowiskoWykonania = &srodowisko
	}
	if parametry.TrybUprawnien != "" {
		tryb := parametry.TrybUprawnien
		zmiana.TrybUprawnien = &tryb
	}
	if parametry.Agent != "" {
		agent := parametry.Agent
		zmiana.Agent = &agent
	}
	return zmiana
}
