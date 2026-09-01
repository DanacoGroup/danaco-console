package core

import (
	"context"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// emiter rozgłasza zdarzenia zmiany do połączeń konta wołającego. Kontrakt nie ma osobnego protokołu synchronizacji wielourządzeniowej: nośnikiem zmiany jest zdarzenie właściwe zmienionemu obszarowi. Nadajnik niepodłączony nie jest błędem.
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
	e.wyslijDoKonta(ctx, shared.EventSessionChanged, s.Id, zdarzenie)
}

// okno rozgłasza zmianę okna komunikacji wraz z rodzajem zmiany i sprawcą wywołania, biorącym z kontekstu.
func (e *emiter) okno(ctx context.Context, zmiana shared.ChangeKind, o shared.Window) {
	zdarzenie := shared.WindowChangedEvent{Change: zmiana, Window: o}
	zdarzenie.Actor, zdarzenie.ActorClientId = sprawca(ctx)
	e.wyslijDoKonta(ctx, shared.EventWindowChanged, o.SessionId, zdarzenie)
}

// wiadomosc rozgłasza zmianę wiadomości okna. Pole Message.role mówi user niezależnie od tego, czy wiadomość wpisał człowiek, czy asystent jego klawiaturą, bo rola opisuje miejsce w rozmowie, nie rękę. Dopiero actor odróżnia jedno od drugiego.
func (e *emiter) wiadomosc(ctx context.Context, zmiana shared.ChangeKind, w shared.Message) {
	zdarzenie := shared.MessageChangedEvent{Change: zmiana, Message: w}
	zdarzenie.Actor, zdarzenie.ActorClientId = sprawca(ctx)
	e.wyslijDoKonta(ctx, shared.EventMessageChanged, w.SessionId, zdarzenie)
}

// ustawienie rozgłasza zmianę konfiguracji na dowolnym z ośmiu poziomów zasięgu do konta wołającego. Sesji komunikatu nie da się wyznaczyć dla poziomów szerszych niż karta sesji, więc zdarzenie idzie bez niej.
func (e *emiter) ustawienie(ctx context.Context, zmiana shared.ChangeKind, w shared.ConfigEntry) {
	idSesji := ""
	if w.Scope == shared.ConfigScopeSession && w.ScopeId != nil {
		idSesji = *w.ScopeId
	}
	e.wyslijDoKonta(ctx, shared.EventConfigChanged, idSesji, shared.ConfigChangedEvent{Change: zmiana, Entry: w})
}

// kolejka rozgłasza zmianę kolejki jednego silnika pętli wraz z rodzajem zmiany do konta wołającego, bez sprawcy wywołania.
func (e *emiter) kolejka(ctx context.Context, zmiana shared.ChangeKind, k shared.Queue) {
	e.wyslijDoKonta(ctx, shared.EventQueueChanged, k.SessionId, shared.QueueChangedEvent{Change: zmiana, Queue: k})
}

// zlecenieAsystenta rozgłasza `assistant.action.changed` — zmianę stanu zlecenia modułu Assistant. Producentem jest tor wykonawcy zleceń (adapter_modul_asystent_wykonawca.go), a nośnikiem — ten sam emiter rdzenia, co dla pozostałych zmian obszarów.
func (e *emiter) zlecenieAsystenta(ctx context.Context, zmiana shared.ChangeKind, idSesji string,
	z shared.AssistantAction) {

	zdarzenie := shared.AssistantActionChangedEvent{Change: zmiana, Action: z}
	zdarzenie.Actor, zdarzenie.ActorClientId = sprawca(ctx)
	e.wyslijDoKonta(ctx, shared.EventAssistantActionChanged, idSesji, zdarzenie)
}

// postep rozgłasza telemetrię postępu procesu. Jedno zdarzenie zasila Execution Monitor okna i Process Monitor warstwy wspólnej, więc rdzeń nie ma drugiego kanału telemetrii. Sesja komunikatu bywa nieznana i wtedy zdarzenie rozgłasza się bez niej.
func (e *emiter) postep(idSesji string, z shared.ProgressChangedEvent) {
	// Telemetria opisuje proces rdzenia, nie czynność konta, i nie ma zamawiającego; idzie do wszystkich połączeń.
	e.wyslij(shared.EventProgressChanged, idSesji, z)
}

// wyslij rozgłasza zdarzenie do wszystkich połączeń rdzenia. Droga wyłącznie dla zdarzeń o stanie całego rdzenia, bez zamawiającego; zdarzenie zmiany zamówionej przez Operatora idzie wyslijDoKonta.
func (e *emiter) wyslij(typ shared.MessageType, idSesji string, tresc any) {
	e.wyslijNaKonto("", typ, idSesji, tresc)
}

// wyslijDoKonta rozgłasza zdarzenie do połączeń konta wołającego, odczytanego z kontekstu żądania tak samo, jak transport nazywa konto gniazda.
func (e *emiter) wyslijDoKonta(ctx context.Context, typ shared.MessageType, idSesji string, tresc any) {
	e.wyslijNaKonto(kontoAdresata(ctx), typ, idSesji, tresc)
}

// wyslijNaKonto składa kopertę zdarzenia i oddaje ją transportowi. Zdarzenie niesie własny identyfikator nadany przez rdzeń, bo nie odpowiada na żadne żądanie. Niepowodzenie kodowania kończy wyłącznie to jedno rozgłoszenie.
func (e *emiter) wyslijNaKonto(konto string, typ shared.MessageType, idSesji string, tresc any) {
	if e == nil || e.nadajnik == nil {
		return
	}
	k, err := protocol.NowaKoperta(typ, identyfikatorZdarzenia(), idSesji, tresc)
	if err != nil {
		return
	}
	e.nadajnik.Rozglos(konto, k)
}

// rozlaczanie oddaje zrywanie gniazd, jeżeli nadajnik transportu je niesie; nadajnik sprawdzianu albo niepodłączony oddaje zero.
func (e *emiter) rozlaczanie() RozlaczanieSesji {
	if e == nil || e.nadajnik == nil {
		return nil
	}
	r, umie := e.nadajnik.(RozlaczanieSesji)
	if !umie {
		return nil
	}
	return r
}
