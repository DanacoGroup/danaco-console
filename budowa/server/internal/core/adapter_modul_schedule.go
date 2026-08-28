// Odpowiedzialność pliku: rodzina `schedule.*` — jedna komenda `schedule.get`,
// odczyt harmonogramów do pary z `automation.schedule.set`. Metody stoją na
// tym samym typie `adapterAutomatyk`, więc odczyt i zapis widzą ten sam stan.
package core

import (
	"context"
	"errors"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// repozytoriumHarmonogramow to rozszerzenie repozytorium modułu Automations
// o dwa odczyty, po identyfikatorze harmonogramu i po całym wykazie, których
// port modułu Automations nie zna i znać nie musi.
type repozytoriumHarmonogramow interface {
	HarmonogramPoKodzie(ctx context.Context, kod string) (dane.Harmonogram, error)
	Harmonogramy(ctx context.Context, tylkoCzynne bool) ([]dane.Harmonogram, error)
}

// HarmonogramyZadania obsługuje `schedule.get` trzema drogami żądania:
// po identyfikatorze harmonogramu, po automatyce albo bez wskazania.
func (a *adapterAutomatyk) HarmonogramyZadania(ctx context.Context,
	z shared.ScheduleGetRequest) (shared.ScheduleGetResponse, error) {

	tylkoCzynne := z.EnabledOnly != nil && *z.EnabledOnly
	kodHarmonogramu := wartoscTekstu(z.ScheduleId)
	kodAutomatyki := wartoscTekstu(z.WorkflowId)

	switch {
	case kodHarmonogramu != "":
		return a.harmonogramWskazany(ctx, kodHarmonogramu, kodAutomatyki, tylkoCzynne)
	case kodAutomatyki != "":
		return a.harmonogramyAutomatyki(ctx, kodAutomatyki, tylkoCzynne)
	default:
		return a.wszystkieHarmonogramy(ctx, tylkoCzynne)
	}
}

// harmonogramWskazany oddaje jeden harmonogram wskazany jego identyfikatorem.
// Gdy żądanie niesie oba wskazania, muszą się zgadzać, inaczej odmowa niesie
// kod `conflict` i obie nazwy.
func (a *adapterAutomatyk) harmonogramWskazany(ctx context.Context, kodHarmonogramu,
	kodAutomatyki string, tylkoCzynne bool) (shared.ScheduleGetResponse, error) {

	wykaz, err := a.wykazHarmonogramow()
	if err != nil {
		return shared.ScheduleGetResponse{}, err
	}
	wiersz, err := wykaz.HarmonogramPoKodzie(ctx, kodHarmonogramu)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.ScheduleGetResponse{}, bladNieznanegoHarmonogramu(kodHarmonogramu)
	}
	if err != nil {
		return shared.ScheduleGetResponse{}, bladAutomatyki(err)
	}
	automatyka, err := a.automatykaHarmonogramu(ctx, wiersz)
	if err != nil {
		return shared.ScheduleGetResponse{}, err
	}
	if kodAutomatyki != "" && automatyka.Kod != kodAutomatyki {
		return shared.ScheduleGetResponse{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeConflict, "moduł Automations: harmonogram "+kodHarmonogramu+
				" należy do automatyki "+automatyka.Kod+", a nie do "+kodAutomatyki))
	}
	if tylkoCzynne && !wiersz.Czynny {
		return pustyWykazHarmonogramow(), nil
	}
	harmonogram, err := a.harmonogramKontraktu(ctx, automatyka.Kod, wiersz)
	if err != nil {
		return shared.ScheduleGetResponse{}, bladAutomatyki(err)
	}
	return shared.ScheduleGetResponse{Schedules: []shared.AutomationSchedule{harmonogram}}, nil
}

// harmonogramyAutomatyki oddaje harmonogramy jednej automatyki. Automatyka
// musi istnieć, harmonogram nie musi — automatyka bez harmonogramu oddaje
// wykaz pusty, stan poprawny schematu, a nie ciszę udającą wynik.
func (a *adapterAutomatyk) harmonogramyAutomatyki(ctx context.Context, kodAutomatyki string,
	tylkoCzynne bool) (shared.ScheduleGetResponse, error) {

	automatyka, err := a.wiersz(ctx, kodAutomatyki)
	if err != nil {
		return shared.ScheduleGetResponse{}, err
	}
	wiersz, err := a.repozytorium.Harmonogram(ctx, automatyka.ID)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return pustyWykazHarmonogramow(), nil
	}
	if err != nil {
		return shared.ScheduleGetResponse{}, bladAutomatyki(err)
	}
	if tylkoCzynne && !wiersz.Czynny {
		return pustyWykazHarmonogramow(), nil
	}
	harmonogram, err := a.harmonogramKontraktu(ctx, automatyka.Kod, wiersz)
	if err != nil {
		return shared.ScheduleGetResponse{}, bladAutomatyki(err)
	}
	return shared.ScheduleGetResponse{Schedules: []shared.AutomationSchedule{harmonogram}}, nil
}

// wszystkieHarmonogramy oddaje komplet harmonogramów platformy. Żądanie bez
// żadnego wskazania o nic nie zawęża, więc i odpowiedź niczego nie ukrywa.
func (a *adapterAutomatyk) wszystkieHarmonogramy(ctx context.Context,
	tylkoCzynne bool) (shared.ScheduleGetResponse, error) {

	wykaz, err := a.wykazHarmonogramow()
	if err != nil {
		return shared.ScheduleGetResponse{}, err
	}
	wiersze, err := wykaz.Harmonogramy(ctx, tylkoCzynne)
	if err != nil {
		return shared.ScheduleGetResponse{}, bladAutomatyki(err)
	}
	odpowiedz := pustyWykazHarmonogramow()
	for _, wiersz := range wiersze {
		automatyka, err := a.automatykaHarmonogramu(ctx, wiersz)
		if err != nil {
			return shared.ScheduleGetResponse{}, err
		}
		harmonogram, err := a.harmonogramKontraktu(ctx, automatyka.Kod, wiersz)
		if err != nil {
			return shared.ScheduleGetResponse{}, bladAutomatyki(err)
		}
		odpowiedz.Schedules = append(odpowiedz.Schedules, harmonogram)
	}
	return odpowiedz, nil
}

// automatykaHarmonogramu oddaje automatykę, do której harmonogram należy.
// Brak automatyki jest usterką rdzenia, nie pomyłką Operatora, bo klucz obcy
// harmonogramu jest wymagany, więc wiersz osierocony znaczy uszkodzoną bazę.
func (a *adapterAutomatyk) automatykaHarmonogramu(ctx context.Context,
	wiersz dane.Harmonogram) (dane.Automatyka, error) {

	automatyka, err := a.repozytorium.AutomatykaPoID(ctx, wiersz.AutomatykaID)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return dane.Automatyka{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeInternalError, "moduł Automations: harmonogram "+wiersz.Kod+
				" wskazuje automatykę, której nie ma w bazie"))
	}
	if err != nil {
		return dane.Automatyka{}, bladAutomatyki(err)
	}
	return automatyka, nil
}

// wykazHarmonogramow oddaje rozszerzenie repozytorium albo odmowę — bez tych
// odczytów nie ma z czego złożyć odpowiedzi, bez udawania, że wykaz jest pusty.
func (a *adapterAutomatyk) wykazHarmonogramow() (repozytoriumHarmonogramow, error) {
	wykaz, ok := a.repozytorium.(repozytoriumHarmonogramow)
	if !ok {
		return nil, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"moduł Automations: repozytorium nie niesie odczytu harmonogramów"))
	}
	return wykaz, nil
}

// pustyWykazHarmonogramow składa odpowiedź z wykazem pustym, lecz niezerowym.
// Pole `schedules` wychodzi jako `[]`, nie jako `null`, wymagane kontraktem.
func pustyWykazHarmonogramow() shared.ScheduleGetResponse {
	return shared.ScheduleGetResponse{Schedules: []shared.AutomationSchedule{}}
}

// bladNieznanegoHarmonogramu nazywa harmonogram, którego wskazany kod nie
// odpowiada żadnemu wierszowi repozytorium.
func bladNieznanegoHarmonogramu(kod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
		"moduł Automations: harmonogram nie istnieje: "+kod))
}
