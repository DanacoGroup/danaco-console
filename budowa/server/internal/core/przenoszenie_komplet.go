package core

import (
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// Plik składa komplet kontekstu przenoszonego jedną komendą z trzech warstw; wcześniejsza wygrywa.

// zrodloHistorii jest tą częścią dziennika rozmowy, której przenoszenie
// naprawdę używa: wykazem wiadomości okna. Dziennik rozmowy wypełnia
// ten kształt bez żadnej zmiany.
type zrodloHistorii interface {
	Wykaz(idOkna string, przed *string, ograniczenie *int) ([]shared.Message, bool)
}

// Zgodność dziennika rozmowy z tym kształtem sprawdzana jest przy kompilacji,
// żeby wpięcie w montażu nie rozjechało się po cichu.
var _ zrodloHistorii = (*dziennikRozmowy)(nil)

// polaczKomplety nakłada komplet zapisany przy oknie na komplet z żądania.
// Składnik wskazany w żądaniu zostaje; składnik pominięty bierze się z zapisu.
func polaczKomplety(zadanie, zapisany shared.ContextBundle) shared.ContextBundle {
	komplet := zadanie
	if pustyNapis(komplet.Prompt) {
		komplet.Prompt = zapisany.Prompt
	}
	if len(komplet.DocumentIds) == 0 {
		komplet.DocumentIds = zapisany.DocumentIds
	}
	if pustyNapis(komplet.ProjectId) {
		komplet.ProjectId = zapisany.ProjectId
	}
	if len(komplet.AgentIds) == 0 {
		komplet.AgentIds = zapisany.AgentIds
	}
	if len(komplet.HistoryMessageIds) == 0 {
		komplet.HistoryMessageIds = zapisany.HistoryMessageIds
	}
	if len(komplet.KnowledgeSourceIds) == 0 {
		komplet.KnowledgeSourceIds = zapisany.KnowledgeSourceIds
	}
	if len(komplet.ExecutionParams) == 0 {
		komplet.ExecutionParams = zapisany.ExecutionParams
	}
	return komplet
}

// uzupelnijZrodlem dokłada do kompletu to, co system odczytuje sam: projekt
// sesji źródłowej, historię rozmowy okna źródłowego i parametry wykonania tego
// okna. Składnika już wypełnionego nie rusza.
func uzupelnijZrodlem(komplet shared.ContextBundle, zrodlo session.Okno,
	idProjektu string, historia []string) shared.ContextBundle {

	if pustyNapis(komplet.ProjectId) && idProjektu != "" {
		projekt := idProjektu
		komplet.ProjectId = &projekt
	}
	if len(komplet.HistoryMessageIds) == 0 && len(historia) > 0 {
		komplet.HistoryMessageIds = historia
	}
	if len(komplet.ExecutionParams) == 0 {
		komplet.ExecutionParams = parametryZOkna(zrodlo)
	}
	return komplet
}

// identyfikatoryWiadomosci wyjmuje z historii okna same identyfikatory —
// komplet kontraktu niesie odwołania, nie treści.
func identyfikatoryWiadomosci(wiadomosci []shared.Message) []string {
	if len(wiadomosci) == 0 {
		return nil
	}
	identyfikatory := make([]string, 0, len(wiadomosci))
	for _, wiadomosc := range wiadomosci {
		identyfikatory = append(identyfikatory, wiadomosc.Id)
	}
	return identyfikatory
}

// pustyNapis mówi, czy pole opcjonalne kontraktu nie niesie wartości, czyli wskaźnik jest pusty albo wskazuje na napis pusty.
func pustyNapis(p *string) bool {
	return p == nil || *p == ""
}
