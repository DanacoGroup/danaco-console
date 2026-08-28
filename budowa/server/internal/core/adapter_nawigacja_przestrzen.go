package core

import (
	"context"
	"errors"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// Moduly zwraca moduły wraz z katalogiem ich okien operacyjnych. Wskazanie
// środowiska zawęża wykaz do bocznej nawigacji tego środowiska; suma pozostaje
// liczbą wszystkich modułów platformy, nie długością zwróconego wykazu.
func (a *adapterNawigacji) Moduly(ctx context.Context, z shared.ModuleListRequest) (shared.ModuleListResponse, error) {
	wszystkie, err := a.zestaw.Moduly.Lista(ctx)
	if err != nil {
		return shared.ModuleListResponse{}, err
	}
	wiersze := wszystkie
	if z.EnvironmentId != nil && *z.EnvironmentId != "" {
		if wiersze, err = a.modulySrodowiska(ctx, *z.EnvironmentId); err != nil {
			return shared.ModuleListResponse{}, err
		}
	}
	moduly, err := a.moduly(ctx, wiersze)
	if err != nil {
		return shared.ModuleListResponse{}, err
	}
	return shared.ModuleListResponse{Modules: moduly, Total: len(wszystkie)}, nil
}

// WejdzDoPrzestrzeni przeładowuje przestrzeń roboczą karty sesji na wskazany
// moduł, przestawiając wskazane okno rozmowy albo zakładając okno nowe, gdy
// takiego okna brak, i zwraca kody okien operacyjnych z katalogu modułu.
func (a *adapterNawigacji) WejdzDoPrzestrzeni(ctx context.Context, z shared.WorkspaceEnterRequest) (shared.WorkspaceEnterResponse, error) {
	wynik := shared.WorkspaceEnterResponse{OperationalWindowCodes: []string{}}
	modul, jest, err := a.modulPoWskazaniu(ctx, z.ModuleId)
	if err != nil {
		return shared.WorkspaceEnterResponse{}, err
	}
	kodModulu := z.ModuleId
	if jest {
		kodModulu = modul.Kod
		moduly, err := a.moduly(ctx, []dane.Modul{modul})
		if err != nil {
			return shared.WorkspaceEnterResponse{}, err
		}
		// Kolejność pozycji dotyczy nawigacji jednego środowiska.
		moduly[0].Order = 0
		wynik.Module = moduly[0]
		wynik.OperationalWindowCodes = moduly[0].OperationalWindowCodes
	}
	sesja, err := a.sesjaPrzestrzeni(ctx, z.SessionId)
	if err != nil {
		return shared.WorkspaceEnterResponse{}, err
	}
	okno, err := a.oknoPrzestrzeni(ctx, z, kodModulu)
	if err != nil {
		return shared.WorkspaceEnterResponse{}, err
	}
	wynik.Session, wynik.Window = sesja, okno
	return wynik, nil
}

// modulySrodowiska zwraca moduły bocznej nawigacji środowiska. Środowisko
// nierozpoznane daje wykaz pusty — komenda ma odpowiedzieć, nie odmówić.
func (a *adapterNawigacji) modulySrodowiska(ctx context.Context, wskazanie string) ([]dane.Modul, error) {
	srodowisko, jest, err := a.srodowiskoPoWskazaniu(ctx, wskazanie)
	if err != nil || !jest {
		return nil, err
	}
	return a.zestaw.Moduly.ListaSrodowiska(ctx, srodowisko.ID)
}

// sesjaPrzestrzeni zwraca kartę sesji: żywą z rejestru nadzorcy, a gdy tej nie
// ma — utrwaloną wierszem. Sesja nieznana obu warstwom daje kartę pustą.
func (a *adapterNawigacji) sesjaPrzestrzeni(ctx context.Context, idSesji string) (shared.Session, error) {
	if a.nadzorca != nil {
		if sesja, err := a.nadzorca.Rejestr().Sesja(idSesji); err == nil {
			return sesjaKontraktu(sesja), nil
		}
	}
	wiersz, err := a.zestaw.Sesje.PoIdentyfikatorze(ctx, idSesji)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.Session{}, nil
	}
	if err != nil {
		return shared.Session{}, err
	}
	return sesjaWierszaKontraktu(wiersz), nil
}

// oknoPrzestrzeni zwraca okno komunikacji przestrzeni roboczej: wskazane
// przestawia na nowy moduł, brak wskazania zakłada okno nowe zakładane przez
// nadzorcę, właściciela cyklu życia okna.
func (a *adapterNawigacji) oknoPrzestrzeni(ctx context.Context, z shared.WorkspaceEnterRequest,
	kodModulu string) (shared.Window, error) {
	if a.nadzorca == nil {
		return shared.Window{}, nil
	}
	if z.WindowId != nil && *z.WindowId != "" {
		okno, err := a.nadzorca.Rejestr().ZmienOkno(*z.WindowId, session.Zmiana{Modul: &kodModulu})
		if err == nil {
			return oknoKontraktu(okno), nil
		}
		if !errors.Is(err, session.ErrBrakOkna) {
			return shared.Window{}, bladSesji(err)
		}
	}
	// Okno nowe dostaje komplet parametrów wykonania, nie sam moduł.
	okno, err := a.nadzorca.OtworzOkno(z.SessionId, a.parametryOknaModulu(ctx, z.SessionId, kodModulu))
	if errors.Is(err, session.ErrBrakSesji) {
		return shared.Window{}, nil
	}
	if err != nil {
		return shared.Window{}, bladSesji(err)
	}
	return oknoKontraktu(okno), nil
}
