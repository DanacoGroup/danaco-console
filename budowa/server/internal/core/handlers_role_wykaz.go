// Odpowiedzialność pliku: wykaz i zdjęcie nadań ról (`role.list`,
// `role.remove`) wraz z rozgłoszeniem zdarzenia `role.changed`.
//
// Rodzina `role.*` ma jedno repozytorium i jeden port: `role.list` i
// `role.remove` stoją na tym samym adapterze `adapterRolOkien`, tym samym
// rejestrze okien i tych samych dwóch kolumnach `okno_komunikacji`, co
// `role.assign` i `role.update`. Port `RoleWykaz` osadza `RoleOkien` tak samo,
// jak tamten osadza `Okna`.
//
// Układ sekcji panelu okna prowadzi osobny plik `handlers_panel_sekcje.go`.
package core

import (
	"context"
	"errors"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// RoleWykaz jest portem rodziny `role.*` rozszerzonym o `role.list` i `role.remove`.
type RoleWykaz interface {
	RoleOkien

	WykazRol(ctx context.Context, z shared.RoleListRequest) (shared.RoleListResponse, error)
	ZdejmijRole(ctx context.Context, z shared.RoleRemoveRequest) (shared.RoleRemoveResponse, error)
}

// Rozjazd portu z adapterem zatrzymuje kompilację tutaj, nie na martwej komendzie.
var _ RoleWykaz = (*adapterRolOkien)(nil)

// ── wpięcie ──────────────────────────────────────────────────────────────────

// zarejestrujWykazRol wpina `role.list` i `role.remove`.
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
			// Rozgłasza się wyłącznie zdjęcie, które się odbyło: okno bez roli nie
			// jest odmową, ale nie jest też zmianą.
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
// koordynatora. Wykaz jest widokiem na okna, nie na drugą tabelę: nadanie roli
// to para pól okna (`windowRole`, `coordinatorWindowId`), więc wykaz składa się
// z tego samego wykazu okien, co `window.list`, wraz z więzią doczytaną z bazy.
// Własne zapytanie po rolach dałoby drugą odpowiedź na pytanie o rolę okna.
// Rolę ma każde okno, także samodzielne — stąd komplet.
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
		// Wcielenie leży tam, gdzie zapisuje je `role.update`. Odczyt nieudany nie
		// kończy wykazu: wykaz bez wcieleń bije odmowę.
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

// ZdejmijRole zdejmuje z okna rolę wraz z więzią koordynatora. Zdjęcie roli to
// powrót do roli samodzielnej, a nie skasowanie pola: katalog kontraktu nie ma
// wartości „brak", a okno poza pętlą jest stanem wyjściowym pakietu sesji
// (`session.RolaDomyslna`). Więź znika razem z rolą, w tej samej czynności:
// koordynatora niesie wyłącznie okno wykonawcy.
//
// Okno bez roli nie jest odmową — wraca `removed=false`, bo stan docelowy już
// obowiązuje.
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
	// Wcielenie odchodzi razem z rolą: `persona` opisuje rolę, a nie okno, więc
	// zostawione przy oknie samodzielnym wychodziłoby w `role.list` i w wykazie
	// nadań (`client/src/moduly/multitasking/wykaz-nadan-rol.ts`) jako wcielenie
	// roli, której już nie ma. Niepowodzenie zapisu kończy komendę: pół zdjęcia
	// roli zostawiłoby ten sam rozjazd, tyle że po odmowie.
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
// rolę koordynatora.
//
// Więź ma dwa końce. Zdjęcie roli samemu koordynatorowi zostawiłoby wykonawców
// z więzią do okna, które koordynatorem już nie jest. Zamiast odmawiać zdjęcia
// roli, drugi koniec doprowadza się do stanu zgodnego: wykonawca zostaje
// wykonawcą, ale bez koordynatora — stan dopuszczalny.
//
// Niepowodzenie nie cofa zdjęcia roli: rola została zdjęta, a więź jest
// następstwem, nie warunkiem.
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
