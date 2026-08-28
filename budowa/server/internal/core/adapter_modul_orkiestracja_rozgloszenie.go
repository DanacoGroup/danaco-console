// Odpowiedzialność pliku: rozgłoszenie zmiany podagenta — zdarzenie `subagent.changed`
// i jedyne przejście, którym stan podagenta wchodzi do bazy, dla powołania, zmiany i zatrzymania.
package core

import (
	"context"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// ustawStanPodagenta zapisuje stan podagenta i rozgłasza zmianę; jest jedynym przejściem
// do zapisu stanu, więc żaden inny zapis nie wymknie się rozgłoszeniu.
func (a *adapterPodagentow) ustawStanPodagenta(ctx context.Context, kod string,
	stan string, komunikat *string) error {

	if err := a.repozytorium.UstawStan(ctx, kod, stan, komunikat); err != nil {
		return err
	}
	a.rozglosPodagenta(ctx, kod, shared.ChangeKindUpdated)
	return nil
}

// rozglosPowolanie rozgłasza podagentów świeżo założonych, z wierszy odczytu bazy przy tym samym założeniu.
func (a *adapterPodagentow) rozglosPowolanie(wiersze []dane.Podagent) {
	if a.rozgloszenie == nil {
		return
	}
	for _, wiersz := range wiersze {
		a.rozgloszenie.podagent(shared.ChangeKindCreated, podagentKontraktu(wiersz))
	}
}

// rozglosPodagenta dobiera wiersz po kodzie i rozgłasza jego stan bieżący; nieudany dobór
// kończy wyłącznie rozgłoszenie, zostawiając ślad w dzienniku.
func (a *adapterPodagentow) rozglosPodagenta(ctx context.Context, kod string,
	zmiana shared.ChangeKind) {

	if a.rozgloszenie == nil || kod == "" {
		return
	}
	wiersze, err := a.repozytorium.PodagenciPoKodach(ctx, []string{kod})
	if err != nil {
		a.zapisz("subagent.changed: stan podagenta %s zapisany, ale nie rozgłoszony: %v", kod, err)
		return
	}
	for _, wiersz := range wiersze {
		a.rozgloszenie.podagent(zmiana, podagentKontraktu(wiersz))
	}
}

// podagent rozgłasza `subagent.changed` ze wskazaniem karty sesji, gdy okno wykonawcy
// do niej należy; podagent bez karty rozgłasza się bez niej, zamiast wcale.
func (e *emiter) podagent(zmiana shared.ChangeKind, p shared.Subagent) {
	idSesji := ""
	if p.SessionId != nil {
		idSesji = *p.SessionId
	}
	e.wyslij(shared.EventSubagentChanged, idSesji,
		shared.SubagentChangedEvent{Change: zmiana, Subagent: p})
}
