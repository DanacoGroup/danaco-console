// Plik wpina obszar dostępów, punktów dostępu i nadań dostępu okna rozmowy,
// rozgłaszając zmianę obu bytów zdarzeniem właściwym obszarowi.
package core

import (
	"context"

	"danacoconsole/shared"
)

// PunktyDostepu jest portem katalogu punktów dostępu, obsługującym dodanie,
// wykaz, zmianę, usunięcie i sprawdzenie punktu.
type PunktyDostepu interface {
	Dodaj(ctx context.Context, z shared.AccessPointAddRequest) (shared.AccessPointAddResponse, error)
	Wykaz(ctx context.Context, z shared.AccessPointListRequest) (shared.AccessPointListResponse, error)
	Zmien(ctx context.Context, z shared.AccessPointUpdateRequest) (shared.AccessPointUpdateResponse, error)
	Usun(ctx context.Context, z shared.AccessPointRemoveRequest) (shared.AccessPointRemoveResponse, error)
	Sprawdz(ctx context.Context, z shared.AccessPointCheckRequest) (shared.AccessPointCheckResponse, error)
}

// NadaniaDostepu jest portem zbioru nadań okna rozmowy, obsługującym dodanie,
// wykaz, zmianę i usunięcie nadania.
type NadaniaDostepu interface {
	Dodaj(ctx context.Context, z shared.AccessGrantAddRequest) (shared.AccessGrantAddResponse, error)
	Wykaz(ctx context.Context, z shared.AccessGrantListRequest) (shared.AccessGrantListResponse, error)
	Zmien(ctx context.Context, z shared.AccessGrantUpdateRequest) (shared.AccessGrantUpdateResponse, error)
	Usun(ctx context.Context, z shared.AccessGrantRemoveRequest) (shared.AccessGrantRemoveResponse, error)
}

// zarejestrujPunktyDostepu wpina pięć komend katalogu punktów i rozgłasza
// zmianę punktu po dodaniu, zmianie i usunięciu.
func zarejestrujPunktyDostepu(r *Rejestr, punkty PunktyDostepu, e *emiter) {
	if r == nil || punkty == nil {
		return
	}
	r.Zarejestruj(shared.CommandAccessPointList, obsluz(punkty.Wykaz))
	r.Zarejestruj(shared.CommandAccessPointCheck, obsluz(punkty.Sprawdz))

	r.Zarejestruj(shared.CommandAccessPointAdd,
		obsluz(func(ctx context.Context, z shared.AccessPointAddRequest) (shared.AccessPointAddResponse, error) {
			w, err := punkty.Dodaj(ctx, z)
			if err == nil {
				e.punktDostepu(shared.ChangeKindCreated, w.Point)
			}
			return w, err
		}))

	r.Zarejestruj(shared.CommandAccessPointUpdate,
		obsluz(func(ctx context.Context, z shared.AccessPointUpdateRequest) (shared.AccessPointUpdateResponse, error) {
			w, err := punkty.Zmien(ctx, z)
			if err == nil {
				e.punktDostepu(shared.ChangeKindUpdated, w.Point)
			}
			return w, err
		}))

	r.Zarejestruj(shared.CommandAccessPointRemove,
		obsluz(func(ctx context.Context, z shared.AccessPointRemoveRequest) (shared.AccessPointRemoveResponse, error) {
			w, err := punkty.Usun(ctx, z)
			if err == nil && w.Removed {
				e.punktDostepu(shared.ChangeKindDeleted, shared.AccessPoint{Id: z.AccessPointId})
			}
			return w, err
		}))
}

// zarejestrujNadaniaDostepu wpina cztery komendy zbioru nadań okna i rozgłasza
// zmianę nadania po dodaniu, zmianie i usunięciu.
func zarejestrujNadaniaDostepu(r *Rejestr, nadania NadaniaDostepu, e *emiter) {
	if r == nil || nadania == nil {
		return
	}
	r.Zarejestruj(shared.CommandAccessGrantList, obsluz(nadania.Wykaz))

	r.Zarejestruj(shared.CommandAccessGrantAdd,
		obsluz(func(ctx context.Context, z shared.AccessGrantAddRequest) (shared.AccessGrantAddResponse, error) {
			w, err := nadania.Dodaj(ctx, z)
			if err == nil {
				e.nadanieDostepu(shared.ChangeKindCreated, w.Grant)
			}
			return w, err
		}))

	r.Zarejestruj(shared.CommandAccessGrantUpdate,
		obsluz(func(ctx context.Context, z shared.AccessGrantUpdateRequest) (shared.AccessGrantUpdateResponse, error) {
			w, err := nadania.Zmien(ctx, z)
			if err == nil {
				e.nadanieDostepu(shared.ChangeKindUpdated, w.Grant)
			}
			return w, err
		}))

	r.Zarejestruj(shared.CommandAccessGrantRemove,
		obsluz(func(ctx context.Context, z shared.AccessGrantRemoveRequest) (shared.AccessGrantRemoveResponse, error) {
			w, err := nadania.Usun(ctx, z)
			if err == nil && w.Removed {
				e.nadanieDostepu(shared.ChangeKindDeleted,
					shared.AccessGrant{Id: z.GrantId, WindowId: oknoZbioru(w.Grants)})
			}
			return w, err
		}))
}

// punktDostepu rozgłasza zmianę punktu dostępu. Punkt nie należy do żadnej
// sesji, więc zdarzenie idzie bez jej wskazania — dociera do wszystkich
// połączeń konta.
func (e *emiter) punktDostepu(zmiana shared.ChangeKind, p shared.AccessPoint) {
	e.wyslij(shared.EventAccessPointChanged, "",
		shared.AccessPointChangedEvent{Change: zmiana, Point: p})
}

// nadanieDostepu rozgłasza zmianę nadania wraz z oknem, którego zbioru
// dotyczy, zdarzeniem dostępnym wszystkim połączeniom konta.
func (e *emiter) nadanieDostepu(zmiana shared.ChangeKind, n shared.AccessGrant) {
	e.wyslij(shared.EventAccessGrantChanged, "",
		shared.AccessGrantChangedEvent{Change: zmiana, WindowId: n.WindowId, Grant: n})
}

// oknoZbioru odczytuje okno ze zbioru nadań pozostałych po odebraniu jednego
// z nich. Zbiór opróżniony do zera nie pozwala wskazać okna; zdarzenie idzie
// wtedy bez niego, zamiast nie iść wcale.
func oknoZbioru(nadania []shared.AccessGrant) string {
	if len(nadania) == 0 {
		return ""
	}
	return nadania[0].WindowId
}
