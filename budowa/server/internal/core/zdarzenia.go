package core

import (
	"context"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// emiter rozgłasza zdarzenia zmiany do wszystkich połączeń konta.
//
// Kontrakt nie ma osobnego protokołu synchronizacji wielourządzeniowej:
// nośnikiem zmiany jest zdarzenie właściwe zmienionemu obszarowi. Dlatego
// obsługiwacz po udanej zmianie stanu zgłasza ją tutaj, a transport roznosi
// zdarzenie dalej.
//
// Nadajnik niepodłączony nie jest błędem: rdzeń wykonuje komendy także wtedy,
// gdy nikt nie słucha zdarzeń.
type emiter struct {
	nadajnik Nadajnik
}

// nowyEmiter opakowuje nadajnik warstwy transportu.
func nowyEmiter(nadajnik Nadajnik) *emiter {
	return &emiter{nadajnik: nadajnik}
}

// SZEŚĆ ZDARZEŃ NIESIE SPRAWCĘ i wszystkie sześć biorą go z JEDNEGO miejsca —
// `sprawca(ctx)` (`sprawca.go`). Stąd kontekst w podpisie: bez niego rdzeń nie
// ma jak odpowiedzieć na pytanie „czyja ręka", bo odpowiedź jest własnością
// WYWOŁANIA, nie zmienionego bytu. Kontekst niesie już tożsamość żądania
// i tożsamość połączenia, więc sprawca jest tu bytem tej samej klasy i nie
// dokłada ani jednego parametru domenowego.
//
// KONTEKST BEZ GNIAZDA NIE JEST BŁĘDEM: pola zostają puste, a zdarzenie idzie
// tak samo. Kontrakt mówi wprost, że brak znaczy „rdzeń nie potrafił tego
// rozstrzygnąć" — i to jest wtedy prawda.

// sesja rozgłasza zmianę sesji.
func (e *emiter) sesja(ctx context.Context, zmiana shared.ChangeKind, s shared.Session) {
	zdarzenie := shared.SessionChangedEvent{Change: zmiana, Session: s}
	zdarzenie.Actor, zdarzenie.ActorClientId = sprawca(ctx)
	e.wyslij(shared.EventSessionChanged, s.Id, zdarzenie)
}

// okno rozgłasza zmianę okna komunikacji.
func (e *emiter) okno(ctx context.Context, zmiana shared.ChangeKind, o shared.Window) {
	zdarzenie := shared.WindowChangedEvent{Change: zmiana, Window: o}
	zdarzenie.Actor, zdarzenie.ActorClientId = sprawca(ctx)
	e.wyslij(shared.EventWindowChanged, o.SessionId, zdarzenie)
}

// wiadomosc rozgłasza zmianę wiadomości okna.
//
// TU SPRAWCA WAŻY NAJWIĘCEJ. Pole `Message.role` mówi `user` niezależnie od
// tego, czy wiadomość wpisał człowiek, czy asystent jego klawiaturą — bo rola
// opisuje MIEJSCE W ROZMOWIE, nie rękę. Dopiero `actor` odróżnia jedno od
// drugiego i dopiero z nim interfejs może napisać „Asystent" zamiast „spoza
// tego połączenia".
func (e *emiter) wiadomosc(ctx context.Context, zmiana shared.ChangeKind, w shared.Message) {
	zdarzenie := shared.MessageChangedEvent{Change: zmiana, Message: w}
	zdarzenie.Actor, zdarzenie.ActorClientId = sprawca(ctx)
	e.wyslij(shared.EventMessageChanged, w.SessionId, zdarzenie)
}

// ustawienie rozgłasza zmianę konfiguracji na dowolnym z ośmiu poziomów
// zasięgu. Sesji komunikatu nie da się wyznaczyć dla poziomów
// szerszych niż karta sesji, więc zdarzenie idzie bez niej.
func (e *emiter) ustawienie(zmiana shared.ChangeKind, w shared.ConfigEntry) {
	idSesji := ""
	if w.Scope == shared.ConfigScopeSession && w.ScopeId != nil {
		idSesji = *w.ScopeId
	}
	e.wyslij(shared.EventConfigChanged, idSesji, shared.ConfigChangedEvent{Change: zmiana, Entry: w})
}

// kolejka rozgłasza zmianę kolejki jednego silnika pętli.
func (e *emiter) kolejka(zmiana shared.ChangeKind, k shared.Queue) {
	e.wyslij(shared.EventQueueChanged, k.SessionId, shared.QueueChangedEvent{Change: zmiana, Queue: k})
}

// zlecenieAsystenta rozgłasza `assistant.action.changed` — zmianę stanu zlecenia
// modułu Assistant. Producentem jest tor wykonawcy zleceń
// (adapter_modul_asystent_wykonawca.go), a nośnikiem — ten sam emiter rdzenia,
// co dla pozostałych zmian obszarów.
func (e *emiter) zlecenieAsystenta(ctx context.Context, zmiana shared.ChangeKind, idSesji string,
	z shared.AssistantAction) {

	zdarzenie := shared.AssistantActionChangedEvent{Change: zmiana, Action: z}
	zdarzenie.Actor, zdarzenie.ActorClientId = sprawca(ctx)
	e.wyslij(shared.EventAssistantActionChanged, idSesji, zdarzenie)
}

// postep rozgłasza telemetrię postępu procesu. Jedno zdarzenie zasila
// Execution Monitor okna i Process Monitor warstwy wspólnej, więc rdzeń nie ma
// drugiego kanału telemetrii — producentem jest wyłącznie telemetria.go.
//
// Sesja komunikatu bywa nieznana: proces bez okna i bez sesji (na przykład
// kolejka założona przed pierwszą wiadomością) rozgłasza się bez niej, zamiast
// nie rozgłaszać się wcale.
func (e *emiter) postep(idSesji string, z shared.ProgressChangedEvent) {
	e.wyslij(shared.EventProgressChanged, idSesji, z)
}

// wyslij składa kopertę zdarzenia i oddaje ją transportowi. Zdarzenie niesie
// własny identyfikator nadany przez rdzeń, bo nie odpowiada na żadne żądanie.
// Niepowodzenie kodowania kończy wyłącznie to jedno rozgłoszenie.
func (e *emiter) wyslij(typ shared.MessageType, idSesji string, tresc any) {
	if e == nil || e.nadajnik == nil {
		return
	}
	k, err := protocol.NowaKoperta(typ, identyfikatorZdarzenia(), idSesji, tresc)
	if err != nil {
		return
	}
	e.nadajnik.Rozglos(k)
}
