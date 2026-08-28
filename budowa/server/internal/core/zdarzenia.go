package core

import (
	"context"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// emiter rozgłasza zdarzenia zmiany do wszystkich połączeń konta. Kontrakt nie ma osobnego protokołu synchronizacji wielourządzeniowej: nośnikiem zmiany jest zdarzenie właściwe zmienionemu obszarowi. Nadajnik niepodłączony nie jest błędem.
type emiter struct {
	nadajnik Nadajnik
}

// nowyEmiter opakowuje nadajnik warstwy transportu, tworząc gotowy do rozgłaszania zdarzeń emiter zmian.
func nowyEmiter(nadajnik Nadajnik) *emiter {
	return &emiter{nadajnik: nadajnik}
}

// Sześć zdarzeń niesie sprawcę, biorąc go z kontekstu funkcją sprawca. Brak gniazda nie jest błędem.

// sesja rozgłasza zmianę sesji wraz z rodzajem zmiany i sprawcą wywołania, wziętym z kontekstu żądania.
func (e *emiter) sesja(ctx context.Context, zmiana shared.ChangeKind, s shared.Session) {
	zdarzenie := shared.SessionChangedEvent{Change: zmiana, Session: s}
	zdarzenie.Actor, zdarzenie.ActorClientId = sprawca(ctx)
	e.wyslij(shared.EventSessionChanged, s.Id, zdarzenie)
}

// okno rozgłasza zmianę okna komunikacji wraz z rodzajem zmiany i sprawcą wywołania, biorącym z kontekstu.
func (e *emiter) okno(ctx context.Context, zmiana shared.ChangeKind, o shared.Window) {
	zdarzenie := shared.WindowChangedEvent{Change: zmiana, Window: o}
	zdarzenie.Actor, zdarzenie.ActorClientId = sprawca(ctx)
	e.wyslij(shared.EventWindowChanged, o.SessionId, zdarzenie)
}

// wiadomosc rozgłasza zmianę wiadomości okna. Pole Message.role mówi user niezależnie od tego, czy wiadomość wpisał człowiek, czy asystent jego klawiaturą, bo rola opisuje miejsce w rozmowie, nie rękę. Dopiero actor odróżnia jedno od drugiego.
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

// kolejka rozgłasza zmianę kolejki jednego silnika pętli wraz z rodzajem zmiany, bez sprawcy wywołania.
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

// postep rozgłasza telemetrię postępu procesu. Jedno zdarzenie zasila Execution Monitor okna i Process Monitor warstwy wspólnej, więc rdzeń nie ma drugiego kanału telemetrii. Sesja komunikatu bywa nieznana i wtedy zdarzenie rozgłasza się bez niej.
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
