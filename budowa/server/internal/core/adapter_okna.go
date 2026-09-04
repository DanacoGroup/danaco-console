package core

import (
	"context"
	"strconv"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

type adapterOkien struct {
	nadzorca     *session.Nadzorca
	trwalosc     *utrwalaczStanow
	petla        *session.Petla
	przerwijTure PrzerwanieTury
	turaWBiegu   TuraWBiegu
	wiezie       czytelnikWiezi
	kanaly       czytelnikKanalow
	moduly       czytelnikModulow
	straz        StrazEksperta
	// konta rozstrzyga granicę wykazu okien: bez wskazania sesji wykaz szedłby
	// po rejestrze całej instalacji, więc konto musi być czym rozstrzygnąć.
	konta dane.RepozytoriumKontaWlasciciela
}

func (a *adapterOkien) ZKontami(k dane.RepozytoriumKontaWlasciciela) *adapterOkien {
	a.konta = k
	return a
}

type czytelnikModulow interface {
	Lista(ctx context.Context) ([]dane.Modul, error)
}

type czytelnikKanalow interface {
	Lista(ctx context.Context, tylkoAktywne bool) ([]dane.Kanal, error)
}

// Pusty napis znaczy „okno samodzielne” i jest odpowiedzią poprawną; brak wiersza wraca jako błąd.
type czytelnikWiezi interface {
	Koordynator(ctx context.Context, oknoWykonawcy string) (string, error)
}

func nowyAdapterOkien(nadzorca *session.Nadzorca) *adapterOkien {
	return &adapterOkien{nadzorca: nadzorca}
}

func (a *adapterOkien) ZTrwaloscia(u *utrwalaczStanow) *adapterOkien {
	a.trwalosc = u
	return a
}

func (a *adapterOkien) ZPetla(p *session.Petla) *adapterOkien {
	a.petla = p
	return a
}

func (a *adapterOkien) ZWieziami(c czytelnikWiezi) *adapterOkien {
	a.wiezie = c
	return a
}

func (a *adapterOkien) ZKanalami(c czytelnikKanalow) *adapterOkien {
	a.kanaly = c
	return a
}

// Kanał wskazany wprost zostaje nietknięty, także błędny.
func (a *adapterOkien) kanalDomyslny(ctx context.Context) string {
	if a.kanaly == nil {
		return ""
	}
	kanaly, err := a.kanaly.Lista(ctx, true)
	if err != nil || len(kanaly) == 0 {
		return ""
	}
	return kanaly[0].Kod
}

func (a *adapterOkien) ZModulami(c czytelnikModulow) *adapterOkien {
	a.moduly = c
	return a
}

func (a *adapterOkien) modulDomyslny(ctx context.Context) string {
	if a.moduly == nil {
		return ""
	}
	moduly, err := a.moduly.Lista(ctx)
	if err != nil || len(moduly) == 0 {
		return ""
	}
	return moduly[0].Kod
}

// `module.list` niesie identyfikator i kod; okno opisuje się jednym z nich.
func (a *adapterOkien) kodModulu(ctx context.Context, wskazanie string) string {
	if a.moduly == nil {
		return wskazanie
	}
	moduly, err := a.moduly.Lista(ctx)
	if err != nil {
		return wskazanie
	}
	for _, modul := range moduly {
		if modul.Kod == wskazanie || strconv.FormatInt(modul.ID, 10) == wskazanie {
			return modul.Kod
		}
	}
	return wskazanie
}

func (a *adapterOkien) Utworz(ctx context.Context, z shared.WindowCreateRequest) (shared.WindowCreateResponse, error) {
	if z.ModelChannelId == "" {
		z.ModelChannelId = a.kanalDomyslny(ctx)
	}
	if z.ModuleId == "" {
		z.ModuleId = a.modulDomyslny(ctx)
	}
	// Kodem opisuje okno rejestr sesji i kolumna bazy; inaczej `session.open` opisałby okno inaczej, niż zostało założone.
	z.ModuleId = a.kodModulu(ctx, z.ModuleId)
	okno, err := a.nadzorca.OtworzOkno(z.SessionId, ustawieniaOkna(z))
	if err != nil {
		return shared.WindowCreateResponse{}, bladSesji(err)
	}
	// Okno idzie do bazy od razu, nie dopiero z pierwszą wypowiedzią, tak jak sesja.
	a.utrwalZalozone(ctx, okno)
	return shared.WindowCreateResponse{Window: oknoKontraktu(okno)}, nil
}

func (a *adapterOkien) Wykaz(ctx context.Context, z shared.WindowListRequest) (shared.WindowListResponse, error) {
	okna, err := a.okna(ctx, z.SessionId)
	if err != nil {
		return shared.WindowListResponse{}, bladSesji(err)
	}
	wybrane := make([]session.Okno, 0, len(okna))
	for _, okno := range okna {
		if z.Status != nil && okno.Stan != *z.Status {
			continue
		}
		wybrane = append(wybrane, okno)
	}
	lista := oknaKontraktu(wybrane)
	a.dolozWiezi(ctx, lista)
	return shared.WindowListResponse{Windows: lista}, nil
}

// Gdy wiersz okna istnieje, baza jest źródłem prawdy, także gdy oddaje pusty napis.
func (a *adapterOkien) dolozWiezi(ctx context.Context, lista []shared.Window) {
	if a.wiezie == nil {
		return
	}
	for i := range lista {
		koordynator, err := a.wiezie.Koordynator(ctx, lista[i].Id)
		if err != nil {
			continue
		}
		if koordynator == "" {
			lista[i].CoordinatorWindowId = nil
			continue
		}
		lista[i].CoordinatorWindowId = &koordynator
	}
}

// Zapis do bazy stoi po zmianie w rejestrze nadzorcy, nie przed nią.
func (a *adapterOkien) Zmien(ctx context.Context, z shared.WindowUpdateRequest) (shared.WindowUpdateResponse, error) {
	okno, err := a.nadzorca.Rejestr().ZmienOkno(z.WindowId, zmianaOkna(z))
	if err != nil {
		return shared.WindowUpdateResponse{}, bladSesji(err)
	}
	if err := a.utrwalAgentaOkna(ctx, z.WindowId, z.AgentId); err != nil {
		return shared.WindowUpdateResponse{}, err
	}
	return shared.WindowUpdateResponse{Window: oknoKontraktu(okno)}, nil
}

func (a *adapterOkien) Zamknij(ctx context.Context, z shared.WindowCloseRequest) (shared.WindowCloseResponse, error) {
	// Krok łagodny przed twardym: tura dostaje karencję, zanim nadzorca ubije proces.
	a.domknijTureLagodnie(z.WindowId)
	okno, err := a.nadzorca.ZamknijOkno(z.WindowId)
	if err != nil {
		return shared.WindowCloseResponse{}, bladSesji(err)
	}
	a.trwalosc.StanOkna(ctx, okno.Id, okno.Stan)
	if a.petla != nil {
		a.petla.Zapomnij(okno.Id)
	}
	return shared.WindowCloseResponse{Window: oknoKontraktu(okno)}, nil
}

func (a *adapterOkien) okna(ctx context.Context, idSesji *string) ([]session.Okno, error) {
	if idSesji != nil && *idSesji != "" {
		return a.nadzorca.Rejestr().OknaSesji(*idSesji)
	}
	var wszystkie []session.Okno
	for _, sesja := range a.nadzorca.Rejestr().SesjeKonta(kontoZRepozytorium(ctx, a.konta)) {
		okna, err := a.nadzorca.Rejestr().OknaSesji(sesja.Id)
		if err != nil {
			return nil, err
		}
		wszystkie = append(wszystkie, okna...)
	}
	return wszystkie, nil
}
