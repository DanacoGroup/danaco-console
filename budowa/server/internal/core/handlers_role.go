// Plik wpina dwie komendy rodziny role.* — nadania roli oknu i zmiany roli wraz z jej wcieleniem —
// jako rozszerzenie portu okien, bo rola jest polem okna komunikacji, nie osobnym bytem.
package core

import (
	"context"

	"danacoconsole/shared"
)

// RoleOkien jest portem rodziny role.* — portem okien komunikacji rozszerzonym o dwie komendy roli okna.
type RoleOkien interface {
	Okna

	// NadajRole obsługuje `role.assign`.
	NadajRole(ctx context.Context, z shared.RoleAssignRequest) (shared.RoleAssignResponse, error)
	// ZmienRole obsługuje `role.update`.
	ZmienRole(ctx context.Context, z shared.RoleUpdateRequest) (shared.RoleUpdateResponse, error)

	// OknoRoli oddaje okno po zmianie roli; kontrakt rodziny role.* oddaje samą rolę, nie całe okno.
	OknoRoli(ctx context.Context, idOkna string) (shared.Window, bool)
}

// Adapter wypełnia port w całości. Rozjazd portu z adapterem zatrzymuje
// kompilację tutaj, a nie dopiero na martwej komendzie.
var _ RoleOkien = (*adapterRolOkien)(nil)

// zarejestrujRole wpina dwie komendy rodziny role.* obsługujące nadanie i zmianę roli okna komunikacji.
func zarejestrujRole(r *Rejestr, ro RoleWykaz, e *emiter) {
	if r == nil || ro == nil {
		return
	}

	r.Zarejestruj(shared.CommandRoleAssign,
		obsluz(func(ctx context.Context, z shared.RoleAssignRequest) (shared.RoleAssignResponse, error) {
			odpowiedz, err := ro.NadajRole(ctx, z)
			if err == nil {
				rozglosOknoRoli(ctx, ro, e, odpowiedz.WindowId)
				rozglosZmianeRoli(ctx, ro, e, shared.ChangeKindUpdated, odpowiedz.WindowId)
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandRoleUpdate,
		obsluz(func(ctx context.Context, z shared.RoleUpdateRequest) (shared.RoleUpdateResponse, error) {
			odpowiedz, err := ro.ZmienRole(ctx, z)
			if err == nil {
				rozglosOknoRoli(ctx, ro, e, odpowiedz.WindowId)
				rozglosZmianeRoli(ctx, ro, e, shared.ChangeKindUpdated, odpowiedz.WindowId)
			}
			return odpowiedz, err
		}))
}

// rozglosOknoRoli rozgłasza zmianę okna po nadaniu roli. Okno, którego nie da
// się złożyć, kończy wyłącznie rozgłoszenie — komenda już się powiodła i jej
// wynik nie zależy od tego, czy ktoś słucha zdarzeń.
func rozglosOknoRoli(ctx context.Context, ro RoleWykaz, e *emiter, idOkna string) {
	okno, jest := ro.OknoRoli(ctx, idOkna)
	if !jest {
		return
	}
	e.okno(ctx, shared.ChangeKindUpdated, okno)
}
