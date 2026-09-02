package core

import (
	"encoding/json"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// Domknięcie łańcucha telemetrii: nasłuch wzbogaca telemetrię i rozgłasza zmianę sesji.

type nadajnikZObecnoscia struct {
	nadajnik Nadajnik
	obecnosc *rejestrObecnosci
}

func (r *rejestrObecnosci) OwinNadajnik(nadajnik Nadajnik) Nadajnik {
	if r == nil || nadajnik == nil {
		return nadajnik
	}
	r.nadawca = nowyEmiter(nadajnik)
	return nadajnikZObecnoscia{nadajnik: nadajnik, obecnosc: r}
}

// Nasłuch nie ma prawa zgubić ani opóźnić komunikatu właściwego: wzbogaca wyłącznie telemetrię postępu.
func (n nadajnikZObecnoscia) Rozglos(konto string, k protocol.Koperta) {
	if k.Type != shared.EventProgressChanged {
		n.nadajnik.Rozglos(konto, k)
		return
	}
	postep, idOkna, czytelny := postepZKoperty(k)
	if !czytelny {
		n.nadajnik.Rozglos(konto, k)
		return
	}
	idSesji := n.obecnosc.sesjaOkna(protocol.IdSesji(k), idOkna)
	n.nadajnik.Rozglos(konto, n.obecnosc.zBiegiem(k, postep, idOkna))
	n.obecnosc.przyjmijPostep(konto, idSesji, idOkna, postep.Status)
	// Zmiana etapu, tury albo postępu okna jest zmianą stanu okna operacyjnego wspólną każdemu oknu.
	n.obecnosc.rozglosStanOkna(konto, idOkna)
}

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

func (r *rejestrObecnosci) przyjmijPostep(konto, idSesji, idOkna string, stan shared.ProgressStatus) {
	if !r.czynnosc.Odnotuj(idSesji, idOkna, stan) {
		return
	}
	r.rozglos(konto, idSesji)
}

func (r *rejestrObecnosci) biegZmieniony(bieg shared.LoopState) {
	// Rejestr biegów nie zna zamawiającego, a bieg jest częścią stanu okna koordynatora.
	r.rozglos("", r.sesjaOkna("", bieg.CoordinatorWindowId))
	r.rozglosStanOkna("", bieg.CoordinatorWindowId)
}

func (r *rejestrObecnosci) rozglosStanOkna(konto, idOkna string) {
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
	r.nadawca.wyslijNaKonto(konto, shared.EventWindowStateChanged, okno.IdSesji, shared.WindowStateChangedEvent{
		WindowId: idOkna,
		State:    tresc,
	})
}

func (r *rejestrObecnosci) stanProcesuOkna(idOkna string) shared.ProgressStatus {
	if r == nil {
		return shared.ProgressStatusPending
	}
	return stanProcesuOkna(r.nadzorca, r.czynnosc, idOkna)
}

// Tura modelu nie ma procesu w rejestrze nadzorcy: rozstrzyga jej telemetria.
func stanProcesuOkna(nadzorca *session.Nadzorca, czynnosc *pamiecCzynnosci,
	idOkna string) shared.ProgressStatus {

	if wpis, jest := czynnosc.Okno(idOkna); jest && turaTrwa(wpis.Stan) {
		return wpis.Stan
	}
	if nadzorca == nil {
		return shared.ProgressStatusPending
	}
	proces, jest := nadzorca.Procesy().Proces(idOkna)
	switch {
	case !jest:
		return shared.ProgressStatusPending
	case proces.Zyje():
		return shared.ProgressStatusRunning
	default:
		return shared.ProgressStatusStopped
	}
}

func turaTrwa(stan shared.ProgressStatus) bool {
	return stan == shared.ProgressStatusRunning || stan == shared.ProgressStatusPaused
}

// Kontrakt nie ma osobnego nośnika obecności — nośnikiem zmiany jest zdarzenie właściwe obszarowi.
func (r *rejestrObecnosci) rozglos(konto, idSesji string) {
	if r == nil || r.nadawca == nil || r.nadzorca == nil || idSesji == "" {
		return
	}
	sesja, err := r.nadzorca.Rejestr().Sesja(idSesji)
	if err != nil {
		return
	}
	odpis, jest := r.Odpis(r.kontekstSesji(idSesji), idSesji)
	if !jest {
		return
	}
	// Sprawcą zmiany jest rdzeń: obecność zmienia się od przemiatania stanu sesji, nie od komendy.
	zdarzenie := shared.SessionChangedEvent{
		Change:   shared.ChangeKindUpdated,
		Session:  sesjaKontraktu(sesja),
		Presence: &odpis,
	}
	zdarzenie.Actor, zdarzenie.ActorClientId = sprawca(zSprawcaRdzenia(r.kontekst))
	r.nadawca.wyslijNaKonto(konto, shared.EventSessionChanged, idSesji, zdarzenie)
}

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
