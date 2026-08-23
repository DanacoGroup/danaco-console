package core

import (
	"context"
	"log"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// adapterPrzenoszenia wypełnia port Przenoszenie: przekazuje komplet kontekstu
// między modułami jedną komendą.
//
// Przenoszony jest komplet, nie samo okno. Adapter składa go z żądania,
// z kompletu zapisanego przy oknie źródłowym i z realnego stanu tego okna
// (przenoszenie_komplet.go), a następnie: nakłada parametry wykonania na okno
// docelowe i zapisuje komplet przy tym oknie (przenoszenie_magazyn.go).
// Dzięki temu okno docelowe naprawdę ma polecenie, dokumenty, projekt, agentów,
// historię, źródła wiedzy i parametry — a nie tylko zmieniony moduł.
//
// Okno docelowe wskazane w żądaniu zostaje przestawione; brak wskazania zakłada
// nowe okno, z ustawieniami okna źródłowego i modułem docelowym. Ścieżka jest
// jedna dla wszystkich modułów — rdzeń się tu nie rozgałęzia.
type adapterPrzenoszenia struct {
	nadzorca *session.Nadzorca
	magazyn  *magazynKontekstu
	historia zrodloHistorii
}

// nowyAdapterPrzenoszenia wiąże port z nadzorcą sesji. Magazyn kontekstu
// powstaje od razu, w wariancie pamięciowym: przenoszenie kompletu ma działać
// także wtedy, gdy trwałości nie wpięto.
func nowyAdapterPrzenoszenia(nadzorca *session.Nadzorca) *adapterPrzenoszenia {
	return &adapterPrzenoszenia{nadzorca: nadzorca, magazyn: nowyMagazynKontekstu(nil, nil, nil)}
}

// ZTrwaloscia przenosi magazyn kontekstu na repozytorium konfiguracji, czyli na
// ósmy poziom zasięgu. Komplet przeżywa wtedy restart rdzenia.
func (a *adapterPrzenoszenia) ZTrwaloscia(zycie context.Context,
	repozytorium dane.RepozytoriumKonfiguracji, dziennik *log.Logger) *adapterPrzenoszenia {

	a.magazyn = nowyMagazynKontekstu(zycie, repozytorium, dziennik)
	return a
}

// ZHistoria wskazuje dziennik rozmowy, z którego adapter czyta historię okna
// źródłowego. Bez niego komplet niesie historię tylko wtedy, gdy wskazał ją
// Operator w żądaniu.
func (a *adapterPrzenoszenia) ZHistoria(zrodlo zrodloHistorii) *adapterPrzenoszenia {
	a.historia = zrodlo
	return a
}

// Przenies wykonuje przeniesienie kompletu kontekstu.
func (a *adapterPrzenoszenia) Przenies(_ context.Context, z shared.ContextTransferRequest) (shared.ContextTransferResponse, error) {
	zrodlo, err := a.nadzorca.Rejestr().Okno(z.SourceWindowId)
	if err != nil {
		return shared.ContextTransferResponse{}, bladSesji(err)
	}
	komplet := a.komplet(zrodlo, z)
	cel, err := a.oknoDocelowe(zrodlo, z, komplet)
	if err != nil {
		return shared.ContextTransferResponse{}, bladSesji(err)
	}
	a.magazyn.Zapisz(cel.Id, komplet)
	return shared.ContextTransferResponse{Window: oknoKontraktu(cel), Transferred: true}, nil
}

// komplet składa przenoszony komplet kontekstu okna źródłowego.
func (a *adapterPrzenoszenia) komplet(zrodlo session.Okno, z shared.ContextTransferRequest) shared.ContextBundle {
	komplet := polaczKomplety(z.Bundle, a.magazyn.Odczytaj(zrodlo.Id))
	return uzupelnijZrodlem(komplet, zrodlo, a.projekt(zrodlo.IdSesji), a.historiaOkna(zrodlo.Id))
}

// oknoDocelowe wskazuje okno, do którego trafia komplet: wskazane w żądaniu
// albo świeżo założone.
func (a *adapterPrzenoszenia) oknoDocelowe(zrodlo session.Okno, z shared.ContextTransferRequest,
	komplet shared.ContextBundle) (session.Okno, error) {

	if z.TargetWindowId != nil && *z.TargetWindowId != "" {
		return a.nadzorca.Rejestr().ZmienOkno(*z.TargetWindowId, zmianaCelu(z.TargetModuleId, komplet))
	}
	ustawienia := nalozParametry(zrodlo.Ustawienia, komplet)
	ustawienia.Modul = z.TargetModuleId
	return a.nadzorca.OtworzOkno(sesjaDocelowa(zrodlo, z.TargetSessionId), ustawienia)
}

// projekt odczytuje projekt sesji źródłowej. Sesja nieznana rejestrowi daje
// projekt pusty — przeniesienie i tak się odbywa.
func (a *adapterPrzenoszenia) projekt(idSesji string) string {
	sesja, err := a.nadzorca.Rejestr().Sesja(idSesji)
	if err != nil {
		return ""
	}
	return sesja.IdProjektu
}

// historiaOkna odczytuje identyfikatory wiadomości okna źródłowego.
func (a *adapterPrzenoszenia) historiaOkna(idOkna string) []string {
	if a.historia == nil {
		return nil
	}
	wiadomosci, _ := a.historia.Wykaz(idOkna, nil, nil)
	return identyfikatoryWiadomosci(wiadomosci)
}

// sesjaDocelowa wskazuje sesję okna docelowego: wskazaną w żądaniu albo sesję
// okna źródłowego, gdy żądanie jej nie zmienia.
func sesjaDocelowa(zrodlo session.Okno, wskazana *string) string {
	if wskazana != nil && *wskazana != "" {
		return *wskazana
	}
	return zrodlo.IdSesji
}
