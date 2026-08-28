// Odpowiedzialność pliku: wykaz i zdjęcie nadań ról (`role.list`,
// `role.remove`) wraz z rozgłoszeniem zdarzenia `role.changed`. Port
// `RoleWykaz` osadza `RoleOkien` tak samo, jak `role.assign` osadza `Okna`.
package core

import (
	"context"
	"errors"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// RoleWykaz jest portem rodziny `role.*` rozszerzonym o `role.list`
// i `role.remove`, obok `role.assign` i `role.update` osadzonych w `Okna`.
type RoleWykaz interface {
	RoleOkien

	WykazRol(ctx context.Context, z shared.RoleListRequest) (shared.RoleListResponse, error)
	ZdejmijRole(ctx context.Context, z shared.RoleRemoveRequest) (shared.RoleRemoveResponse, error)
}

// Rozjazd portu z adapterem zatrzymuje kompilację tutaj, nie na martwej
// komendzie odkrytej dopiero w czasie działania rdzenia.
var _ RoleWykaz = (*adapterRolOkien)(nil)

// zarejestrujWykazRol wpina `role.list` i `role.remove`, rozgłaszając
// `role.changed` po każdym udanym zdjęciu roli oknu.
func zarejestrujWykazRol(r *Rejestr, rw RoleWykaz, e *emiter) {
	if r == nil || rw == nil {
		return
	}

	r.Zarejestruj(shared.CommandRoleList,
		obsluz(func(ctx context.Context, z shared.RoleListRequest) (shared.RoleListResponse, error) {
			return rw.WykazRol(ctx, z)
		}))

	r.Zarejestruj(shared.CommandRoleRemove,
		obsluz(func(ctx context.Context, z shared.RoleRemoveRequest) (shared.RoleRemoveResponse, error) {
			odpowiedz, err := rw.ZdejmijRole(ctx, z)
			// Rozgłasza się wyłącznie zdjęcie, które się odbyło.

			// Okno bez roli nie jest odmową, ale nie jest też zmianą.
			if err == nil && odpowiedz.Removed {
				rozglosZmianeRoli(ctx, rw, e, shared.ChangeKindDeleted, odpowiedz.WindowId)
			}
			return odpowiedz, err
		}))
}

// rozglosZmianeRoli rozgłasza `role.changed` dla okna po zmianie jego roli.
// Okno, którego nie da się złożyć, kończy wyłącznie rozgłoszenie — komenda już
// się powiodła.
func rozglosZmianeRoli(ctx context.Context, rw RoleWykaz, e *emiter,
	zmiana shared.ChangeKind, idOkna string) {

	okno, jest := rw.OknoRoli(ctx, idOkna)
	if !jest {
		return
	}
	e.wyslij(shared.EventRoleChanged, okno.SessionId, shared.RoleChangedEvent{
		Change: zmiana,
		Assignment: shared.WindowRoleAssignment{
			WindowId:            okno.Id,
			Role:                okno.WindowRole,
			CoordinatorWindowId: okno.CoordinatorWindowId,
		},
	})
}

// WykazRol oddaje nadania ról oknom, z zawężeniem do wykonawców jednego
// koordynatora, jako widok na wykaz okien, nie na drugą tabelę.
func (a *adapterRolOkien) WykazRol(ctx context.Context,
	z shared.RoleListRequest) (shared.RoleListResponse, error) {

	okna, err := a.Wykaz(ctx, shared.WindowListRequest{})
	if err != nil {
		return shared.RoleListResponse{}, err
	}
	zawezenie := strings.TrimSpace(wartoscTekstu(z.CoordinatorWindowId))

	nadania := make([]shared.WindowRoleAssignment, 0, len(okna.Windows))
	for _, okno := range okna.Windows {
		koordynator := strings.TrimSpace(wartoscTekstu(okno.CoordinatorWindowId))
		if zawezenie != "" && koordynator != zawezenie {
			continue
		}
		nadanie := shared.WindowRoleAssignment{
			WindowId:            okno.Id,
			Role:                okno.WindowRole,
			CoordinatorWindowId: wskaznikPolaRoli(koordynator),
		}
		// Wcielenie leży tam, gdzie zapisuje je `role.update`.

		// Odczyt nieudany nie kończy wykazu: wykaz bez wcieleń bije odmowę.
		if wcielenie, err := a.wcielenieWykazu(ctx, okno.Id); err == nil {
			nadanie.Persona = wskaznikPolaRoli(wcielenie)
		}
		nadania = append(nadania, nadanie)
	}
	return shared.RoleListResponse{Assignments: nadania, Total: len(nadania)}, nil
}

// wcielenieWykazu czyta wcielenie okna. Rdzeń bez konfiguracji oddaje wcielenie
// puste, a nie usterkę: praca na samej pamięci rejestru jest dopuszczalna.
func (a *adapterRolOkien) wcielenieWykazu(ctx context.Context, idOkna string) (string, error) {
	if a.zestaw == nil || a.zestaw.Konfiguracja == nil {
		return "", nil
	}
	return a.wcielenieObowiazujace(ctx, idOkna)
}

// ZdejmijRole zdejmuje z okna rolę wraz z więzią koordynatora, powrotem do
// roli samodzielnej. Okno bez roli nie jest odmową — wraca `removed=false`.
func (a *adapterRolOkien) ZdejmijRole(ctx context.Context,
	z shared.RoleRemoveRequest) (shared.RoleRemoveResponse, error) {

	idOkna := strings.TrimSpace(z.WindowId)
	if idOkna == "" {
		return shared.RoleRemoveResponse{}, bladRoli(
			"żądanie zdjęcia roli bez wskazania okna; Operator poda `windowId` okna, " +
				"z którego rola ma zostać zdjęta — identyfikator oddaje `role.list`")
	}
	zastana, err := a.rolaZastana(ctx, idOkna)
	if err != nil {
		return shared.RoleRemoveResponse{}, err
	}
	if zastana.Rola == shared.WindowRoleStandalone && zastana.Koordynator == "" {
		return shared.RoleRemoveResponse{WindowId: idOkna, Removed: false}, nil
	}

	samodzielna := shared.WindowRole(shared.WindowRoleStandalone)
	bezKoordynatora := ""
	if _, err := a.nadajRole(ctx, idOkna, &samodzielna, &bezKoordynatora); err != nil {
		return shared.RoleRemoveResponse{}, err
	}
	// Wcielenie odchodzi razem z rolą: `persona` opisuje rolę, a nie okno.

	// Niepowodzenie zapisu kończy komendę: pół zdjęcia roli byłoby rozjazdem.
	if _, err := a.zdejmijWcielenieRoli(ctx, idOkna); err != nil {
		return shared.RoleRemoveResponse{}, err
	}
	a.rozwiazWiezWykonawcow(ctx, idOkna)
	return shared.RoleRemoveResponse{WindowId: idOkna, Removed: true}, nil
}

// zdejmijWcielenieRoli kasuje wcielenie okna. Rdzeń bez konfiguracji nie ma
// gdzie trzymać wcielenia, więc nie ma też czego zdejmować — i to nie jest
// usterka.
func (a *adapterRolOkien) zdejmijWcielenieRoli(ctx context.Context, idOkna string) (string, error) {
	if a.zestaw == nil || a.zestaw.Konfiguracja == nil {
		return "", nil
	}
	puste := ""
	return a.zapiszWcielenie(ctx, idOkna, &puste)
}

// rozwiazWiezWykonawcow zdejmuje więź z okien, które podlegały oknu tracącemu
// rolę koordynatora, doprowadzając drugi koniec więzi do stanu zgodnego.
func (a *adapterRolOkien) rozwiazWiezWykonawcow(ctx context.Context, idKoordynatora string) {
	okna, err := a.Wykaz(ctx, shared.WindowListRequest{})
	if err != nil {
		return
	}
	bezKoordynatora := ""
	for _, okno := range okna.Windows {
		if strings.TrimSpace(wartoscTekstu(okno.CoordinatorWindowId)) != idKoordynatora {
			continue
		}
		_, _ = a.nadajRole(ctx, okno.Id, nil, &bezKoordynatora)
	}
}

// rolaZastana odczytuje rolę tą samą dwudrożnością, co `nadajRole`: okno otwarte
// stoi w rejestrze, okno sprzed restartu — w wierszu. Nieznane obu jest odmową.
func (a *adapterRolOkien) rolaZastana(ctx context.Context, idOkna string) (dane.RolaOkna, error) {
	if a.nadzorca != nil {
		if okno, err := a.nadzorca.Rejestr().Okno(idOkna); err == nil {
			return dane.RolaOkna{Rola: okno.RolaOkna, Koordynator: okno.OknoKoordynatora}, nil
		}
	}
	if a.role == nil {
		return dane.RolaOkna{}, bladBrakuOknaRoli(idOkna)
	}
	stan, err := a.role.RolaOkna(ctx, idOkna)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return dane.RolaOkna{}, bladBrakuOknaRoli(idOkna)
	}
	if err != nil {
		return dane.RolaOkna{}, err
	}
	return stan, nil
}
