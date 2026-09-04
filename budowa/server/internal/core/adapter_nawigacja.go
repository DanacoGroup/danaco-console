package core

import (
	"context"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// adapterNawigacji wypełnia port Nawigacja słownikami platformy z bazy oraz
// żywym rejestrem sesji i okien nadzorcy, łącząc wiersz utrwalony z bytem
// żywym tak, aby rozłączenie klienta nie kończyło sesji ani procesów.
type adapterNawigacji struct {
	zestaw     *dane.Zestaw
	nadzorca   *session.Nadzorca
	ustawienia Ustawienia
	klienci    *wieziKlientow
	strumien   zrodloStrumienia
	obecnosc   *rejestrObecnosci
	biegi      *rejestrBiegow
}

// zrodloStrumienia mówi, czy w oknie trwa tura odpowiedzi modelu. Wypełnia je
// domena rozmowy, bo tam biegnie tura. Brak źródła daje fałsz — stan okna ma
// wrócić także wtedy, gdy rozmowa nie jest wpięta.
type zrodloStrumienia func(idOkna string) bool

// nowyAdapterNawigacji wiąże port z repozytoriami, nadzorcą, konfiguracją
// warstwową i rejestrem więzi klientów.
func nowyAdapterNawigacji(zestaw *dane.Zestaw, nadzorca *session.Nadzorca, ustawienia Ustawienia,
	klienci *wieziKlientow, strumien zrodloStrumienia) *adapterNawigacji {
	return &adapterNawigacji{
		zestaw: zestaw, nadzorca: nadzorca, ustawienia: ustawienia,
		klienci: klienci, strumien: strumien,
	}
}

// ZObecnoscia dokłada żywy stan sesji trwających w tle oraz stan biegu
// naprawczego okien koordynatorów.
func (a *adapterNawigacji) ZObecnoscia(obecnosc *rejestrObecnosci, biegi *rejestrBiegow) *adapterNawigacji {
	a.obecnosc, a.biegi = obecnosc, biegi
	return a
}

// StronaGlowna zwraca strefę wyboru środowiska i sesje czynne konta. Strona jest
// jedna, platformowa — dlatego żądanie niesie wyłącznie klienta.
func (a *adapterNawigacji) StronaGlowna(ctx context.Context, z shared.HomeEnterRequest) (shared.HomeEnterResponse, error) {
	srodowiska, err := a.srodowiska(ctx, true)
	if err != nil {
		return shared.HomeEnterResponse{}, err
	}
	sesje, err := a.sesjeCzynne(ctx)
	if err != nil {
		return shared.HomeEnterResponse{}, err
	}
	wynik := shared.HomeEnterResponse{Environments: srodowiska, Sessions: sesje}
	if idSesji, _ := a.klienci.Ognisko(z.ClientId); idSesji != "" {
		wynik.FocusedSessionId = &idSesji
	}
	// Kontrolka powrotu potrzebuje wiedzieć, czy jest dokąd wracać i co tam trwa.
	wynik.Presence = a.obecnosc.Odpisy(ctx, kontoRejestru(ctx, a.zestaw))
	return wynik, nil
}

// Srodowiska zwraca środowiska platformy w kolejności wyświetlania. Kody modułów
// dokłada wyłącznie na żądanie — karta wejścia ich nie rysuje, a boczna
// nawigacja i tak przychodzi z `environment.enter`.
func (a *adapterNawigacji) Srodowiska(ctx context.Context, z shared.EnvironmentListRequest) (shared.EnvironmentListResponse, error) {
	srodowiska, err := a.srodowiska(ctx, z.IncludeModules != nil && *z.IncludeModules)
	if err != nil {
		return shared.EnvironmentListResponse{}, err
	}
	return shared.EnvironmentListResponse{Environments: srodowiska}, nil
}

// WejdzDoSrodowiska zwraca środowisko wraz z jego boczną nawigacją i kartami
// sesji. Wskazanie spoza słownika nie wywraca wejścia: klient dostaje odpowiedź
// pustą i pracuje dalej.
func (a *adapterNawigacji) WejdzDoSrodowiska(ctx context.Context, z shared.EnvironmentEnterRequest) (shared.EnvironmentEnterResponse, error) {
	wynik := shared.EnvironmentEnterResponse{Modules: []shared.Module{}, Sessions: []shared.Session{}}
	wiersz, jest, err := a.srodowiskoPoWskazaniu(ctx, z.EnvironmentId)
	if err != nil || !jest {
		return wynik, err
	}
	wiersze, err := a.zestaw.Moduly.ListaSrodowiska(ctx, wiersz.ID)
	if err != nil {
		return shared.EnvironmentEnterResponse{}, err
	}
	moduly, err := a.moduly(ctx, wiersze)
	if err != nil {
		return shared.EnvironmentEnterResponse{}, err
	}
	sesje, err := a.sesjeSrodowiska(ctx, wiersz.ID)
	if err != nil {
		return shared.EnvironmentEnterResponse{}, err
	}
	wynik.Environment = srodowiskoKontraktu(wiersz, wiersz.Kolejnosc, kodyModulow(moduly), true)
	wynik.Modules, wynik.Sessions = moduly, sesje
	wynik.FocusedSessionId = a.ogniskoWykazu(z.ClientId, z.SessionId, sesje)
	return wynik, nil
}

// srodowiska buduje strefę wyboru: wiersze słownika wraz z macierzą widoczności,
// z której wynika rodzaj bocznej nawigacji.
func (a *adapterNawigacji) srodowiska(ctx context.Context, zModulami bool) ([]shared.Environment, error) {
	wiersze, err := a.zestaw.Srodowiska.Lista(ctx)
	if err != nil {
		return nil, err
	}
	sesji := a.sesjeSrodowisk(ctx)
	wykaz := make([]shared.Environment, 0, len(wiersze))
	for pozycja, wiersz := range wiersze {
		moduly, err := a.zestaw.Moduly.ListaSrodowiska(ctx, wiersz.ID)
		if err != nil {
			return nil, err
		}
		srodowisko := srodowiskoKontraktu(wiersz, pozycja+1, kodyWierszyModulow(moduly), zModulami)
		// Karta środowiska w Centrum niesie miarę własnych sesji czynnych.
		liczba := sesji[srodowisko.Code]
		srodowisko.SessionCount = &liczba
		wykaz = append(wykaz, srodowisko)
	}
	return wykaz, nil
}

// sesjeSrodowisk liczy sesje czynne przypadające na środowisko wejścia.
// Sesja bez środowiska nie wchodzi do żadnej miary.
func (a *adapterNawigacji) sesjeSrodowisk(ctx context.Context) map[string]int {
	miary := map[string]int{}
	if a.nadzorca == nil {
		return miary
	}
	for _, sesja := range a.nadzorca.Rejestr().SesjeKonta(kontoRejestru(ctx, a.zestaw)) {
		if sesja.KodSrodowiska == "" || sesja.Stan != shared.SessionStatusActive {
			continue
		}
		miary[sesja.KodSrodowiska]++
	}
	return miary
}

// sesjeCzynne łączy sesje żywe rejestru nadzorcy z sesjami utrwalonymi w bazie;
// sesja żywa niesie procesy i okna, więc wygrywa z własnym wierszem utrwalonym.
func (a *adapterNawigacji) sesjeCzynne(ctx context.Context) ([]shared.Session, error) {
	wykaz := []shared.Session{}
	widziane := map[string]struct{}{}
	if a.nadzorca != nil {
		for _, sesja := range a.nadzorca.Rejestr().SesjeKonta(kontoRejestru(ctx, a.zestaw)) {
			if !czynna(sesja.Stan) {
				continue
			}
			wykaz = append(wykaz, sesjaKontraktu(sesja))
			widziane[sesja.Id] = struct{}{}
		}
	}
	wiersze, err := a.zestaw.Sesje.Lista(ctx, 0)
	if err != nil {
		return nil, err
	}
	for _, wiersz := range wiersze {
		sesja := sesjaWierszaKontraktu(wiersz)
		if _, jest := widziane[sesja.Id]; jest || !czynna(sesja.Status) {
			continue
		}
		wykaz = append(wykaz, sesja)
		widziane[sesja.Id] = struct{}{}
	}
	return wykaz, nil
}

// sesjeSrodowiska zwraca karty sesji jednego środowiska. Przypisanie sesji do
// środowiska niesie wyłącznie łańcuch `srodowisko → karta_sesji → sesja`, więc
// wykaz idzie z bazy — rejestr nadzorcy o środowiskach nie wie.
func (a *adapterNawigacji) sesjeSrodowiska(ctx context.Context, srodowiskoID int64) ([]shared.Session, error) {
	karty, err := a.zestaw.KartySesji.Lista(ctx, srodowiskoID, dane.KontoOperatora(ctx))
	if err != nil {
		return nil, err
	}
	wykaz := []shared.Session{}
	for _, karta := range karty {
		wiersze, err := a.zestaw.Sesje.Lista(ctx, karta.ID)
		if err != nil {
			return nil, err
		}
		for _, wiersz := range wiersze {
			if czynna(wiersz.Stan) {
				wykaz = append(wykaz, sesjaWierszaKontraktu(wiersz))
			}
		}
	}
	return wykaz, nil
}

// ogniskoWykazu rozstrzyga kartę ogniskowaną po wejściu: wskazana przez klienta,
// a bez wskazania — ostatnio ogniskowana, o ile należy do tego wykazu. Karta
// pusta zostaje pusta.
func (a *adapterNawigacji) ogniskoWykazu(idKlienta string, wskazana *string, sesje []shared.Session) *string {
	if wskazana != nil && *wskazana != "" {
		return wskazana
	}
	ogniskowana, _ := a.klienci.Ognisko(idKlienta)
	for _, sesja := range sesje {
		if sesja.Id == ogniskowana {
			return &ogniskowana
		}
	}
	return nil
}

// czynna odróżnia kartę sesji, do której można wrócić, od karty już zamkniętej i niedostępnej dla konta.
func czynna(stan shared.SessionStatus) bool {
	return stan == shared.SessionStatusActive || stan == shared.SessionStatusPaused
}
