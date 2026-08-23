// Odpowiedzialność pliku: rozgłoszenie zmiany podagenta — zdarzenie
// `subagent.changed` i jedyne przejście, którym stan podagenta wchodzi do bazy.
//
// Panel Subagent Network nie musi odpytywać `subagent.list`, dopóki zdarzenie
// leci przy każdej zmianie stanu, nie tylko przy tej widzianej przez uchwyt
// komendy. Stąd `ustawStanPodagenta`: jedno przejście, przez które idą wszystkie
// zapisy stanu, więc żaden nie wymknie się rozgłoszeniu.
//
// Trzy momenty, które panel musi zobaczyć:
//   - powołanie — `created`, wprost z wierszy założonych przez `subagent.spawn`;
//   - zmiana stanu — `updated`, przy wejściu w `running`, przy przepisaniu stanu
//     pozycji na stan podagenta i przy niepowodzeniu pracy w tle;
//   - zatrzymanie — `updated`, bo `subagent.stop` nie kasuje wiersza, tylko
//     przestawia go na `stopped`; `deleted` kazałoby panelowi zdjąć podagenta
//     z wykazu, a ma on tam zostać widoczny jako zatrzymany.
//
// Rozgłoszenie odczytuje wiersz po zapisie, a nie składa go z tego, co zapisał.
// Kolumny czasu (`rozpoczeto`, `zakonczono`) i pole wyniku nadaje zapytanie
// (COALESCE w `ustawStanPodagenta` schematu), więc struktura złożona w rdzeniu
// rozjechałaby się z tym, co odda `subagent.list`.
//
// Emisja pusta znosi się sama: brak nadajnika, nieudany odczyt wiersza po
// zapisie ani wykaz pusty nie wywracają czynności — praca podagenta już się
// wykonała, a zdarzenie jest jej relacją, nie jej warunkiem.
package core

import (
	"context"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// ustawStanPodagenta zapisuje stan podagenta i rozgłasza zmianę.
//
// Jedyne przejście do `UstawStan`: repozytorium woła się wyłącznie stąd —
// wywołanie z pominięciem tej metody zostawiłoby panel bez zdarzenia dla tego
// przejścia stanu. Zapis nieudany nie rozgłasza niczego: zdarzenie ma mówić
// o zmianie, która naprawdę zaszła.
func (a *adapterPodagentow) ustawStanPodagenta(ctx context.Context, kod string,
	stan string, komunikat *string) error {

	if err := a.repozytorium.UstawStan(ctx, kod, stan, komunikat); err != nil {
		return err
	}
	a.rozglosPodagenta(ctx, kod, shared.ChangeKindUpdated)
	return nil
}

// rozglosPowolanie rozgłasza podagentów świeżo założonych. Wiersze idą wprost
// z założenia — są odczytem z bazy, więc drugiego zapytania nie potrzeba.
func (a *adapterPodagentow) rozglosPowolanie(wiersze []dane.Podagent) {
	if a.rozgloszenie == nil {
		return
	}
	for _, wiersz := range wiersze {
		a.rozgloszenie.podagent(shared.ChangeKindCreated, podagentKontraktu(wiersz))
	}
}

// rozglosPodagenta dobiera wiersz po kodzie i rozgłasza jego stan bieżący.
// Nieudany dobór kończy wyłącznie rozgłoszenie — zapis stanu już się powiódł —
// ale zostawia ślad w dzienniku, bo panel zostaje wtedy z obrazem sprzed zmiany.
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

// podagent rozgłasza `subagent.changed`. Podagent wisi na oknie wykonawcy, a to
// okno należy do karty sesji — zdarzenie idzie więc z jej wskazaniem, żeby
// dotarło tam, gdzie panel Subagent Network jest otwarty. Podagent bez karty
// (powołany w rozmowie spoza sesji) rozgłasza się bez niej, zamiast nie
// rozgłaszać się wcale.
func (e *emiter) podagent(zmiana shared.ChangeKind, p shared.Subagent) {
	idSesji := ""
	if p.SessionId != nil {
		idSesji = *p.SessionId
	}
	e.wyslij(shared.EventSubagentChanged, idSesji,
		shared.SubagentChangedEvent{Change: zmiana, Subagent: p})
}
