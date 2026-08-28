// Okno Scheduler: cykliczność, wyzwalacze i chwila najbliższego
// uruchomienia automatyki. Harmonogram obowiązuje dopiero po powiązaniu
// z istniejącą automatyką.
package core

import (
	"context"
	"errors"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// UstawHarmonogram zapisuje cykliczność i wyzwalacze automatyki, odmawiając
// przy wskazaniu automatyki nieznanej.
func (a *adapterAutomatyk) UstawHarmonogram(ctx context.Context,
	z shared.AutomationScheduleSetRequest) (shared.AutomationScheduleSetResponse, error) {

	wiersz, err := a.wiersz(ctx, z.WorkflowId)
	if err != nil {
		return shared.AutomationScheduleSetResponse{}, err
	}
	zastany, kod := a.kodHarmonogramu(ctx, wiersz.ID)
	czynny := czyCzynna(z.Enabled)
	wyzwalacze := wierszeWyzwalaczy(z.Triggers, zastany)
	nastepne := chwilaNastepnegoUruchomienia(z.Cron, wyzwalacze, czynny, time.Now().UTC())

	zapisany, err := a.repozytorium.ZapiszHarmonogram(ctx, dane.Harmonogram{
		Kod: kod, AutomatykaID: wiersz.ID, Cron: z.Cron, StrefaCzasowa: z.TimeZone,
		Czynny: czynny, NastepneUruchomienie: nastepne,
	}, wyzwalacze)
	if err != nil {
		return shared.AutomationScheduleSetResponse{}, bladAutomatyki(err)
	}
	harmonogram, err := a.harmonogramKontraktu(ctx, wiersz.Kod, zapisany)
	if err != nil {
		return shared.AutomationScheduleSetResponse{}, bladAutomatyki(err)
	}
	return shared.AutomationScheduleSetResponse{Schedule: harmonogram}, nil
}

// Harmonogram oddaje harmonogram automatyki albo brak. Służy oknu Scheduler
// przy otwarciu, zanim Operator cokolwiek zmieni.
//
// Kontrakt nie zna komendy `schedule.get`, więc odczyt idzie wyłącznie tą
// drogą: wewnątrz rdzenia, do wykazu automatyk.
func (a *adapterAutomatyk) Harmonogram(ctx context.Context,
	kodAutomatyki string) (shared.AutomationSchedule, bool, error) {

	wiersz, err := a.wiersz(ctx, kodAutomatyki)
	if err != nil {
		return shared.AutomationSchedule{}, false, err
	}
	zapisany, err := a.repozytorium.Harmonogram(ctx, wiersz.ID)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.AutomationSchedule{}, false, nil
	}
	if err != nil {
		return shared.AutomationSchedule{}, false, bladAutomatyki(err)
	}
	harmonogram, err := a.harmonogramKontraktu(ctx, wiersz.Kod, zapisany)
	if err != nil {
		return shared.AutomationSchedule{}, false, bladAutomatyki(err)
	}
	return harmonogram, true, nil
}

// kodHarmonogramu zwraca wyzwalacze zastane i identyfikator harmonogramu:
// zachowany, gdy harmonogram już był, nowy przy pierwszym zapisie.
func (a *adapterAutomatyk) kodHarmonogramu(ctx context.Context,
	automatykaID int64) ([]dane.WyzwalaczAutomatyki, string) {

	zapisany, err := a.repozytorium.Harmonogram(ctx, automatykaID)
	if err != nil {
		return nil, nowyIdentyfikator(przedrostekHarmonogramu)
	}
	zastane, err := a.repozytorium.Wyzwalacze(ctx, zapisany.ID)
	if err != nil {
		zastane = nil
	}
	return zastane, zapisany.Kod
}

// wierszeWyzwalaczy przekłada wyzwalacze kontraktu na wiersze. Żądanie bez pola
// `triggers` zostawia wyzwalacze zastane — zmiana samej cykliczności nie ma
// kasować webhooka ustawionego wcześniej.
func wierszeWyzwalaczy(wyzwalacze []shared.AutomationTrigger,
	zastane []dane.WyzwalaczAutomatyki) []dane.WyzwalaczAutomatyki {

	if wyzwalacze == nil {
		return zastane
	}
	wiersze := make([]dane.WyzwalaczAutomatyki, 0, len(wyzwalacze))
	for numer, wyzwalacz := range wyzwalacze {
		if wyzwalacz.Expression == "" || wyzwalacz.Kind == "" {
			// Wyzwalacz bez treści nie ma czego obserwować i jest pomijany
			// zamiast odmowy całego zapisu.
			continue
		}
		kod := wartoscTekstu(wyzwalacz.Id)
		if kod == "" {
			kod = nowyIdentyfikator(przedrostekWyzwalacza)
		}
		wiersze = append(wiersze, dane.WyzwalaczAutomatyki{
			Kod: kod, Rodzaj: string(wyzwalacz.Kind), Wyrazenie: wyzwalacz.Expression,
			Czynny: wyzwalacz.Enabled == nil || *wyzwalacz.Enabled, Kolejnosc: numer + 1,
		})
	}
	return wiersze
}

// harmonogramKontraktu składa harmonogram kontraktu wraz z wyzwalaczami,
// odczytanymi osobnym zapytaniem.
func (a *adapterAutomatyk) harmonogramKontraktu(ctx context.Context, kodAutomatyki string,
	wiersz dane.Harmonogram) (shared.AutomationSchedule, error) {

	wyzwalacze, err := a.repozytorium.Wyzwalacze(ctx, wiersz.ID)
	if err != nil {
		return shared.AutomationSchedule{}, err
	}
	harmonogram := shared.AutomationSchedule{
		Id: wiersz.Kod, WorkflowId: kodAutomatyki, Cron: wiersz.Cron,
		TimeZone: wiersz.StrefaCzasowa, Triggers: wyzwalaczeKontraktu(wyzwalacze),
		Enabled: wiersz.Czynny, UpdatedAt: chwilaBazy(wiersz.Zaktualizowano),
	}
	if wiersz.NastepneUruchomienie != nil {
		nastepne := chwilaBazy(*wiersz.NastepneUruchomienie)
		harmonogram.NextRunAt = &nastepne
	}
	return harmonogram, nil
}

// wyzwalaczeKontraktu przekłada wiersze wyzwalaczy z bazy danych na byty
// kontraktu zwracane wołającemu.
func wyzwalaczeKontraktu(wiersze []dane.WyzwalaczAutomatyki) []shared.AutomationTrigger {
	wyzwalacze := make([]shared.AutomationTrigger, 0, len(wiersze))
	for _, wiersz := range wiersze {
		kod, czynny := wiersz.Kod, wiersz.Czynny
		wyzwalacze = append(wyzwalacze, shared.AutomationTrigger{
			Id: &kod, Kind: shared.AutomationTriggerKind(wiersz.Rodzaj),
			Expression: wiersz.Wyrazenie, Enabled: &czynny,
		})
	}
	return wyzwalacze
}
