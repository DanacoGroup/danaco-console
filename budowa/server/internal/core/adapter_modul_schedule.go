// Odpowiedzialność pliku: rodzina `schedule.*` — jedna komenda `schedule.get`,
// odczyt harmonogramów do pary z `automation.schedule.set`.
//
// To nie jest drugi moduł harmonogramu. Harmonogram ma w rdzeniu jednego
// właściciela — moduł Automations. Zapis prowadzi
// `adapter_modul_automations_harmonogram.go` i to on składa harmonogram
// kontraktu wraz z wyzwalaczami; ten plik dokłada wyłącznie odczyt trzema
// drogami, których zapis nie potrzebował. Metody stoją na tym samym typie
// `adapterAutomatyk`, więc odczyt i zapis widzą ten sam stan i to samo
// repozytorium; osobny adapter byłby drugą prawdą o harmonogramie.
//
// Trzy drogi żądania, bo tyle niesie kontrakt:
//   - `scheduleId` — jeden harmonogram po własnym identyfikatorze;
//   - `workflowId` — harmonogramy jednej automatyki (schemat dopuszcza najwyżej
//     jeden: UNIQUE na `harmonogram_automatyki.automatyka_id`
//     w `migracja_040_harmonogramy_przebiegi.sql`);
//   - bez wskazania — komplet harmonogramów platformy.
//
// `enabledOnly` jest sitem, nie warunkiem istnienia. Harmonogram wyłączony
// istnieje; żądanie z `enabledOnly` mówi „oddaj wyłącznie obowiązujące", więc
// odsianie wyłączonego oddaje wykaz pusty, a nie odmowę. Odmowa `not_found`
// należy się bytowi, którego nie ma — i tak też jest tu użyta.
package core

import (
	"context"
	"errors"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// repozytoriumHarmonogramow to rozszerzenie repozytorium modułu Automations
// o dwa odczyty, których okno Scheduler nie potrzebowało: po identyfikatorze
// samego harmonogramu i po całym wykazie.
//
// Interfejs stoi po stronie czytelnika: deklaracja mieszka tutaj, a nie
// w `dane.RepozytoriumAutomatyk`, bo wymaga jej wyłącznie rodzina `schedule.*`;
// port modułu Automations tych czynności nie zna i znać nie musi (wzorzec
// `repozytoriumWpisowPamieci` z rodziny `memory.*`).
type repozytoriumHarmonogramow interface {
	HarmonogramPoKodzie(ctx context.Context, kod string) (dane.Harmonogram, error)
	Harmonogramy(ctx context.Context, tylkoCzynne bool) ([]dane.Harmonogram, error)
}

// HarmonogramyZadania obsługuje `schedule.get`.
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
//
// Gdy żądanie niesie oba wskazania, muszą się zgadzać. Harmonogram należy do
// dokładnie jednej automatyki, więc żądanie „harmonogram H automatyki W", w
// którym H należy do innej automatyki, jest wewnętrznie sprzeczne. Sprzeczności
// nie da się spełnić ani po cichu zamienić na jedno ze wskazań — odmowa niesie
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

// harmonogramyAutomatyki oddaje harmonogramy jednej automatyki.
//
// Automatyka musi istnieć, harmonogram nie musi. Wskazanie nieistniejącej
// automatyki jest odmową `not_found` z jej nazwą — inaczej pomyłka w
// identyfikatorze wyglądałaby jak „automatyka bez harmonogramu". Automatyka
// istniejąca, której Operator harmonogramu jeszcze nie nadał, oddaje wykaz
// pusty: to stan poprawny schematu (kolumna `harmonogram_automatyki.automatyka_id`
// jest w tabeli harmonogramu, więc automatyka bez wiersza jest dopuszczona),
// a nie cisza udająca wynik.
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
// Kontrakt niesie w `AutomationSchedule.workflowId` identyfikator zewnętrzny
// automatyki, a harmonogram trzyma klucz wiersza — bez tego odczytu pole
// wyszłoby puste albo z liczbą, której klient nie zna.
//
// Brak automatyki jest usterką rdzenia, nie pomyłką Operatora: klucz obcy
// harmonogramu jest wymagany i kasuje się kaskadowo
// (`migracja_040_harmonogramy_przebiegi.sql`), więc wiersz osierocony znaczy
// uszkodzoną bazę. Stąd `internal_error`, a nie `not_found`.
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

// wykazHarmonogramow oddaje rozszerzenie repozytorium albo odmowę. Repozytorium
// bez tych odczytów nie ma z czego złożyć odpowiedzi — cisza albo pusty wykaz
// mówiłyby wtedy „nie ma harmonogramów", choć mogą być wszystkie.
func (a *adapterAutomatyk) wykazHarmonogramow() (repozytoriumHarmonogramow, error) {
	wykaz, ok := a.repozytorium.(repozytoriumHarmonogramow)
	if !ok {
		return nil, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"moduł Automations: repozytorium nie niesie odczytu harmonogramów"))
	}
	return wykaz, nil
}

// pustyWykazHarmonogramow składa odpowiedź z wykazem pustym, lecz niezerowym.
// Pole `schedules` jest w kontrakcie wymagane, więc wychodzi jako `[]`, a nie
// jako `null` — klient odróżnia „nic nie spełnia warunków" od braku pola.
func pustyWykazHarmonogramow() shared.ScheduleGetResponse {
	return shared.ScheduleGetResponse{Schedules: []shared.AutomationSchedule{}}
}

// bladNieznanegoHarmonogramu nazywa harmonogram, którego nie ma.
func bladNieznanegoHarmonogramu(kod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
		"moduł Automations: harmonogram nie istnieje: "+kod))
}
