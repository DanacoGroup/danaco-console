package core

import (
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// Składanie kompletu kontekstu przenoszonego jedną komendą.
//
// Komplet ma siedem składników: polecenie wyjściowe, dokumenty, projekt,
// agenci, historia rozmowy, źródła wiedzy i parametry wykonania. Żaden z nich
// nie może zginąć po drodze, więc komplet powstaje z trzech warstw, w tej
// kolejności: to, co przyniosło żądanie, potem komplet zapisany przy oknie
// źródłowym, a na końcu realny stan okna i sesji źródłowej.
//
// Warstwa wcześniejsza wygrywa: Operator, który wskazał składnik wprost,
// nie zostaje nadpisany tym, co system odczytał sam.

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

// pustyNapis mówi, czy pole opcjonalne kontraktu nie niesie wartości.
func pustyNapis(p *string) bool {
	return p == nil || *p == ""
}
