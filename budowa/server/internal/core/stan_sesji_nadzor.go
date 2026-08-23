package core

import (
	"encoding/json"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// Domknięcie łańcucha telemetrii: producent → szyna → odbiorca.
//
// Odbiorcą zdarzenia `progress.changed` jest kontrolka sesji trwającej w tle,
// więc nasłuch stoi na szynie, a nie w adapterze rozmowy.
//
// Nasłuch robi dwie rzeczy naraz i obie w jednym przejściu koperty:
//   - wzbogaca telemetrię okna koordynatora o stan biegu naprawczego,
//     bo pole `loop` zdarzenia postępu nie ma innego producenta;
//   - odnotowuje punkt pracy i — gdy zmienił się stan pracy okna — rozgłasza
//     `session.changed` z żywym stanem sesji, żeby kontrolka nie musiała
//     odpytywać rdzenia.
//
// Nie powstaje tu drugi producent telemetrii. Zdarzenie postępu przechodzi
// dalej dokładnie jedno, tylko pełniejsze.

// nadajnikZObecnoscia jest szyną zdarzeń widzianą przez telemetrię postępu.
type nadajnikZObecnoscia struct {
	nadajnik Nadajnik
	obecnosc *rejestrObecnosci
}

// OwinNadajnik zakłada nasłuch obecności na nadajnik transportu i zapamiętuje
// go jako drogę rozgłaszania. Brak rejestru albo brak nadajnika zwraca nadajnik
// bez zmiany.
func (r *rejestrObecnosci) OwinNadajnik(nadajnik Nadajnik) Nadajnik {
	if r == nil || nadajnik == nil {
		return nadajnik
	}
	r.nadawca = nowyEmiter(nadajnik)
	return nadajnikZObecnoscia{nadajnik: nadajnik, obecnosc: r}
}

// Rozglos przepuszcza komunikat, wzbogacając wyłącznie telemetrię postępu.
// Komunikat niebędący postępem i postęp nieczytelny idą dalej nietknięte —
// nasłuch nie ma prawa zgubić ani opóźnić komunikatu właściwego.
func (n nadajnikZObecnoscia) Rozglos(k protocol.Koperta) {
	if k.Type != shared.EventProgressChanged {
		n.nadajnik.Rozglos(k)
		return
	}
	postep, idOkna, czytelny := postepZKoperty(k)
	if !czytelny {
		n.nadajnik.Rozglos(k)
		return
	}
	idSesji := n.obecnosc.sesjaOkna(protocol.IdSesji(k), idOkna)
	n.nadajnik.Rozglos(n.obecnosc.zBiegiem(k, postep, idOkna))
	n.obecnosc.przyjmijPostep(idSesji, idOkna, postep.Status)
	// Zmiana etapu, tury albo postępu okna jest zmianą stanu okna operacyjnego
	// wspólną każdemu oknu. Telemetria postępu jest jej jedynym
	// producentem na rdzeniu, więc `window.state.changed` wychodzi tą samą drogą,
	// obok wzbogaconej telemetrii — klient okna subskrybuje właśnie to zdarzenie.
	n.obecnosc.rozglosStanOkna(idOkna)
}

// postepZKoperty rozpakowuje telemetrię postępu wraz z oknem, którego dotyczy.
// Proces bez okna jest poprawną telemetrią (kolejka założona przed pierwszą
// wiadomością), ale nie mówi nic o oknie sesji — więc dla obecności milczy.
func postepZKoperty(k protocol.Koperta) (shared.ProgressChangedEvent, string, bool) {
	var postep shared.ProgressChangedEvent
	if err := protocol.LadunekDo(k, &postep); err != nil {
		return postep, "", false
	}
	if postep.WindowId == nil || *postep.WindowId == "" {
		return postep, "", false
	}
	return postep, *postep.WindowId, true
}

// zBiegiem dokłada do telemetrii stan biegu naprawczego okna koordynatora.
// Okno bez biegu zostawia kopertę bez zmiany, tak samo jak niepowodzenie
// spakowania — telemetria ma dojść nawet wtedy, gdy wzbogacenie się nie uda.
func (r *rejestrObecnosci) zBiegiem(k protocol.Koperta, postep shared.ProgressChangedEvent,
	idOkna string) protocol.Koperta {

	bieg, jest := r.biegi.Stan(idOkna)
	if !jest {
		return k
	}
	postep.Loop = &bieg
	wzbogacona, err := protocol.NowaKoperta(k.Type, k.Id, protocol.IdSesji(k), postep)
	if err != nil {
		return k
	}
	return wzbogacona
}

// przyjmijPostep odnotowuje punkt pracy okna i rozgłasza żywy stan sesji
// wyłącznie wtedy, gdy zmienił się stan pracy. Kolejny etap tej samej tury
// niczego nie rozgłasza.
func (r *rejestrObecnosci) przyjmijPostep(idSesji, idOkna string, stan shared.ProgressStatus) {
	if !r.czynnosc.Odnotuj(idSesji, idOkna, stan) {
		return
	}
	r.rozglos(idSesji)
}

// biegZmieniony jest odbiorcą zmian rejestru biegów. Zmiana licznika obiegów,
// braku postępu albo zatrzymania biegu zmienia to, co kontrolka ma pokazać,
// więc idzie tą samą drogą co zmiana stanu pracy okna.
func (r *rejestrObecnosci) biegZmieniony(bieg shared.LoopState) {
	r.rozglos(r.sesjaOkna("", bieg.CoordinatorWindowId))
	// Bieg naprawczy jest częścią stanu okna koordynatora, więc jego
	// zmiana rozgłasza także `window.state.changed` okna, którego dotyczy.
	r.rozglosStanOkna(bieg.CoordinatorWindowId)
}

// rozglosStanOkna rozgłasza `window.state.changed` — zmianę stanu okna
// operacyjnego wspólną każdemu oknu. Zdarzenie niesie okno wraz
// z żywym odpisem jego stanu wykonania: parametrami okna, stanem procesu, turą
// strumienia i biegiem naprawczym. Nośnikiem jest ten sam emiter rdzenia, co dla
// pozostałych zmian obszarów.
//
// Okno nieznane rejestrowi nadzorcy albo brak nadajnika kończy rozgłoszenie bez
// błędu: telemetria stanu okna jest dodatkiem do pracy rdzenia.
func (r *rejestrObecnosci) rozglosStanOkna(idOkna string) {
	if r == nil || r.nadawca == nil || r.nadzorca == nil || idOkna == "" {
		return
	}
	okno, err := r.nadzorca.Rejestr().Okno(idOkna)
	if err != nil {
		return
	}
	stan := shared.WindowStateGetResponse{
		Window:        oknoKontraktu(okno),
		ProcessStatus: r.stanProcesuOkna(idOkna),
		Streaming:     r.strumien != nil && r.strumien(idOkna),
	}
	if bieg, jest := r.biegi.Stan(idOkna); jest {
		stan.Loop = &bieg
	}
	tresc, err := json.Marshal(stan)
	if err != nil {
		return
	}
	r.nadawca.wyslij(shared.EventWindowStateChanged, okno.IdSesji, shared.WindowStateChangedEvent{
		WindowId: idOkna,
		State:    tresc,
	})
}

// stanProcesuOkna przekłada obecność procesu okna na stan telemetrii postępu
// tak samo jak `adapterNawigacji.stanProcesu`: okno bez procesu oczekuje,
// proces żywy pracuje, proces zamknięty jest zatrzymany.
func (r *rejestrObecnosci) stanProcesuOkna(idOkna string) shared.ProgressStatus {
	if r.nadzorca == nil {
		return shared.ProgressStatusPending
	}
	proces, jest := r.nadzorca.Procesy().Proces(idOkna)
	switch {
	case !jest:
		return shared.ProgressStatusPending
	case proces.Zyje():
		return shared.ProgressStatusRunning
	default:
		return shared.ProgressStatusStopped
	}
}

// rozglos wypuszcza `session.changed` z żywym stanem sesji. Zdarzenie niesie
// sesję i jej obecność razem, bo kontrakt nie ma osobnego nośnika obecności —
// nośnikiem zmiany jest zdarzenie właściwe zmienionemu obszarowi.
func (r *rejestrObecnosci) rozglos(idSesji string) {
	if r == nil || r.nadawca == nil || r.nadzorca == nil || idSesji == "" {
		return
	}
	sesja, err := r.nadzorca.Rejestr().Sesja(idSesji)
	if err != nil {
		return
	}
	odpis, jest := r.Odpis(r.kontekst, idSesji)
	if !jest {
		return
	}
	// Sprawcą zmiany jest rdzeń: obecność zmienia się od przemiatania stanu
	// sesji, a nie od komendy — nie ma tu ani gniazda, ani żądania. Znak
	// sprawcy stawiamy wprost, bo brak gniazda sam z siebie znaczy „nie
	// wiadomo", a nie „rdzeń" (`sprawca.go`).
	zdarzenie := shared.SessionChangedEvent{
		Change:   shared.ChangeKindUpdated,
		Session:  sesjaKontraktu(sesja),
		Presence: &odpis,
	}
	zdarzenie.Actor, zdarzenie.ActorClientId = sprawca(zSprawcaRdzenia(r.kontekst))
	r.nadawca.wyslij(shared.EventSessionChanged, idSesji, zdarzenie)
}

// sesjaOkna rozstrzyga sesję komunikatu: wskazaną wprost, a bez wskazania —
// wyczytaną z okna. Telemetria bywa uboższa od koperty, a kontrolka i tak musi
// wiedzieć, której sesji dotyczy punkt pracy.
func (r *rejestrObecnosci) sesjaOkna(idSesji, idOkna string) string {
	if idSesji != "" || r.nadzorca == nil || idOkna == "" {
		return idSesji
	}
	okno, err := r.nadzorca.Rejestr().Okno(idOkna)
	if err != nil {
		return ""
	}
	return okno.IdSesji
}
