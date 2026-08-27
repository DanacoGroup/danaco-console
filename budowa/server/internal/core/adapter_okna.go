package core

import (
	"context"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// adapterOkien wypełnia port Okna rejestrem okien pakietu sesji.
//
// Okno jest jednostką wykonania: własny moduł, własny kanał modelu, własna lista
// katalogów. Wiele okien jednej sesji żyje równolegle, dlatego każda
// czynność wskazuje okno wprost.
type adapterOkien struct {
	nadzorca *session.Nadzorca
	trwalosc *utrwalaczStanow
	petla    *session.Petla
	// przerwijTure i turaWBiegu składają krok łagodnego zamykania; puste znaczy brak karencji.
	przerwijTure PrzerwanieTury
	turaWBiegu   TuraWBiegu
	// wiezie czyta z bazy więź koordynator-wykonawca, składając odpowiedź z pamięci i z bazy.
	wiezie czytelnikWiezi
	// kanaly rozstrzyga kanał modelu okna zakładanego bez wskazania.
	kanaly czytelnikKanalow
	// moduly rozstrzyga moduł okna zakładanego bez wskazania.
	moduly czytelnikModulow
	// straz czyta zakres eksperta przy nakładaniu na okno; pusta znaczy pełny dostęp.
	straz StrazEksperta
}

// czytelnikModulow oddaje moduły znane katalogowi rdzenia, gotowe do wyboru
// modułu domyślnego okna zakładanego bez wskazania.
type czytelnikModulow interface {
	Lista(ctx context.Context) ([]dane.Modul, error)
}

// czytelnikKanalow oddaje kanały modelu znane rejestrowi, gotowe do wyboru
// kanału domyślnego okna zakładanego bez wskazania.
type czytelnikKanalow interface {
	Lista(ctx context.Context, tylkoAktywne bool) ([]dane.Kanal, error)
}

// czytelnikWiezi oddaje koordynatora wskazanego okna. Pusty napis znaczy „okno
// samodzielne” i jest odpowiedzią poprawną; brak wiersza wraca jako błąd.
type czytelnikWiezi interface {
	Koordynator(ctx context.Context, oknoWykonawcy string) (string, error)
}

// nowyAdapterOkien wiąże port Okna z nadzorcą, czytelnikiem modułów
// i czytelnikiem kanałów, gotowy do dalszego dołożenia trwałości.
func nowyAdapterOkien(nadzorca *session.Nadzorca) *adapterOkien {
	return &adapterOkien{nadzorca: nadzorca}
}

// ZTrwaloscia dokłada zapis stanu okna do bazy. Bez niego stan zamkniętego okna
// żyje wyłącznie w pamięci procesu — zapis jest dodatkiem, nie warunkiem pracy
// adaptera.
func (a *adapterOkien) ZTrwaloscia(u *utrwalaczStanow) *adapterOkien {
	a.trwalosc = u
	return a
}

// ZPetla wskazuje pętlę koordynator–wykonawca, z której trzeba wykreślić
// zamknięte okno. Bez tego licznik obiegów i strumień okna zostają w pamięci
// rdzenia po jego zamknięciu.
func (a *adapterOkien) ZPetla(p *session.Petla) *adapterOkien {
	a.petla = p
	return a
}

// ZWieziami wskazuje repozytorium, z którego wykaz okien bierze więź
// koordynator–wykonawca. Bez niego wykaz pracuje na samej pamięci i nie widzi
// przekazań utrwalonych przez `window.handoff`.
func (a *adapterOkien) ZWieziami(c czytelnikWiezi) *adapterOkien {
	a.wiezie = c
	return a
}

// ZKanalami wskazuje rejestr, z którego okno bierze kanał domyślny, gdy
// zakładane jest bez wskazania kanału wprost.
func (a *adapterOkien) ZKanalami(c czytelnikKanalow) *adapterOkien {
	a.kanaly = c
	return a
}

// kanalDomyslny oddaje kod pierwszego czynnego kanału rejestru, dla okna
// zakładanego bez wskazania kanału. Rdzeń zna swój rejestr i to on rozstrzyga
// wartość domyślną — kanał wskazany wprost zostaje nietknięty, także błędny.
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

// ZModulami wskazuje katalog, z którego okno bierze moduł domyślny, gdy
// zakładane jest bez wskazania modułu wprost.
func (a *adapterOkien) ZModulami(c czytelnikModulow) *adapterOkien {
	a.moduly = c
	return a
}

// modulDomyslny oddaje kod pierwszego modułu katalogu, dla okna zakładanego
// bez wskazania modułu. Katalog rdzenia niesie moduły studio, library, agents
// i pozostałe.
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

// Utworz zakłada okno komunikacji w sesji, zapisując je do bazy od razu, nie
// dopiero z pierwszą wypowiedzią.
func (a *adapterOkien) Utworz(ctx context.Context, z shared.WindowCreateRequest) (shared.WindowCreateResponse, error) {
	if z.ModelChannelId == "" {
		z.ModelChannelId = a.kanalDomyslny(ctx)
	}
	if z.ModuleId == "" {
		z.ModuleId = a.modulDomyslny(ctx)
	}
	okno, err := a.nadzorca.OtworzOkno(z.SessionId, ustawieniaOkna(z))
	if err != nil {
		return shared.WindowCreateResponse{}, bladSesji(err)
	}
	// Okno idzie do bazy od razu, nie dopiero z pierwszą wypowiedzią, tak jak sesja.
	a.utrwalZalozone(ctx, okno)
	return shared.WindowCreateResponse{Window: oknoKontraktu(okno)}, nil
}

// Wykaz zwraca okna sesji albo okna wszystkich sesji, gdy sesji nie wskazano.
// Odpowiedź składa się z dwóch źródeł: fakt bieżący z rejestru pamięciowego
// i fakt postanowiony z bazy.
func (a *adapterOkien) Wykaz(ctx context.Context, z shared.WindowListRequest) (shared.WindowListResponse, error) {
	okna, err := a.okna(z.SessionId)
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

// dolozWiezi nanosi na wykaz więź koordynator-wykonawca odczytaną z bazy. Gdy
// wiersz okna istnieje, baza jest źródłem prawdy, także gdy oddaje pusty
// napis; gdy wiersza nie ma, zostaje wartość z pamięci.
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

// Zmien zmienia ustawienia okna wybiórczo — pole niewskazane zostaje bez
// zmiany. Zapis do bazy stoi po zmianie w rejestrze nadzorcy, nie przed nią.
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

// Zamknij zamyka okno i ubija jego proces wraz z drzewem potomstwa. Pozostałe
// okna sesji pracują dalej.
func (a *adapterOkien) Zamknij(_ context.Context, z shared.WindowCloseRequest) (shared.WindowCloseResponse, error) {
	// Krok łagodny przed twardym: tura dostaje karencję, zanim nadzorca ubije proces.
	a.domknijTureLagodnie(z.WindowId)
	okno, err := a.nadzorca.ZamknijOkno(z.WindowId)
	if err != nil {
		return shared.WindowCloseResponse{}, bladSesji(err)
	}
	a.trwalosc.StanOkna(okno.Id, okno.Stan)
	if a.petla != nil {
		a.petla.Zapomnij(okno.Id)
	}
	return shared.WindowCloseResponse{Window: oknoKontraktu(okno)}, nil
}

// okna zwraca okna jednej sesji albo okna wszystkich sesji, gdy sesji nie
// wskazano, jako wykaz gotowy do naniesienia więzi.
func (a *adapterOkien) okna(idSesji *string) ([]session.Okno, error) {
	if idSesji != nil && *idSesji != "" {
		return a.nadzorca.Rejestr().OknaSesji(*idSesji)
	}
	var wszystkie []session.Okno
	for _, sesja := range a.nadzorca.Rejestr().Sesje() {
		okna, err := a.nadzorca.Rejestr().OknaSesji(sesja.Id)
		if err != nil {
			return nil, err
		}
		wszystkie = append(wszystkie, okna...)
	}
	return wszystkie, nil
}
