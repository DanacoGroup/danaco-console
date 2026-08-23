package core

import (
	"context"
	"errors"
	"strconv"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// Odczyt słowników platformy na potrzeby nawigacji: rozpoznanie wskazania
// środowiska i modułu oraz złożenie modułu wraz z macierzą widoczności
// i katalogiem okien operacyjnych.
//
// Wskazanie z żądania rozpoznajemy dwojako — kodem (`talkin`, `studio`)
// i identyfikatorem wiersza — bo kontrakt niesie jedno pole `environmentId`
// / `moduleId`, a klient może mieć w ręku dowolne z nich. Wskazanie
// nierozpoznane nie jest błędem: wraca fałsz, a obsługiwacz odpowiada pusto.

// srodowiskoPoWskazaniu odnajduje środowisko po kodzie albo po identyfikatorze
// wiersza.
func (a *adapterNawigacji) srodowiskoPoWskazaniu(ctx context.Context, wskazanie string) (dane.Srodowisko, bool, error) {
	if wskazanie == "" {
		return dane.Srodowisko{}, false, nil
	}
	srodowisko, err := a.zestaw.Srodowiska.PoKodzie(ctx, wskazanie)
	if err == nil {
		return srodowisko, true, nil
	}
	if !errors.Is(err, dane.ErrBrakWiersza) {
		return dane.Srodowisko{}, false, err
	}
	id, err := strconv.ParseInt(wskazanie, 10, 64)
	if err != nil {
		return dane.Srodowisko{}, false, nil
	}
	wiersze, err := a.zestaw.Srodowiska.Lista(ctx)
	if err != nil {
		return dane.Srodowisko{}, false, err
	}
	for _, wiersz := range wiersze {
		if wiersz.ID == id {
			return wiersz, true, nil
		}
	}
	return dane.Srodowisko{}, false, nil
}

// modulPoWskazaniu odnajduje moduł po kodzie albo po identyfikatorze wiersza.
func (a *adapterNawigacji) modulPoWskazaniu(ctx context.Context, wskazanie string) (dane.Modul, bool, error) {
	if wskazanie == "" {
		return dane.Modul{}, false, nil
	}
	modul, err := a.zestaw.Moduly.PoKodzie(ctx, wskazanie)
	if err == nil {
		return modul, true, nil
	}
	if !errors.Is(err, dane.ErrBrakWiersza) {
		return dane.Modul{}, false, err
	}
	id, err := strconv.ParseInt(wskazanie, 10, 64)
	if err != nil {
		return dane.Modul{}, false, nil
	}
	wiersze, err := a.zestaw.Moduly.Lista(ctx)
	if err != nil {
		return dane.Modul{}, false, err
	}
	for _, wiersz := range wiersze {
		if wiersz.ID == id {
			return wiersz, true, nil
		}
	}
	return dane.Modul{}, false, nil
}

// moduly przekłada wiersze modułów na byty kontraktu, dokładając macierz
// widoczności i katalog okien operacyjnych. Kolejność bierze z porządku
// zapytania — dla wykazu środowiska jest to kolumna `srodowisko_modul.kolejnosc`.
func (a *adapterNawigacji) moduly(ctx context.Context, wiersze []dane.Modul) ([]shared.Module, error) {
	macierz, err := a.macierzSrodowisk(ctx)
	if err != nil {
		return nil, err
	}
	wykaz := make([]shared.Module, 0, len(wiersze))
	for pozycja, wiersz := range wiersze {
		okna, err := a.kodyOkienModulu(ctx, wiersz.ID)
		if err != nil {
			return nil, err
		}
		wykaz = append(wykaz, modulKontraktu(wiersz, pozycja+1, macierz[wiersz.ID], okna))
	}
	return wykaz, nil
}

// macierzSrodowisk buduje odwzorowanie moduł → kody środowisk, w których moduł
// jest widoczny. Moduł nieobecny w macierzy nie ma okna modułowego w żadnym
// środowisku — jest dostępny wyłącznie ze strony głównej.
func (a *adapterNawigacji) macierzSrodowisk(ctx context.Context) (map[int64][]string, error) {
	// Jedno złączenie zamiast pętli N+1: całą macierz zwraca pojedyncze
	// zapytanie, a nie wykaz środowisk plus zapytanie o moduły każdego z nich.
	// Porządek wyniku ustala `dane/macierz.go`:
	// ORDER BY s.kolejnosc, s.kod, sm.kolejnosc, m.kod.
	return nowyAdapterMacierzy(a.zestaw.Macierz).KodySrodowisk(ctx)
}

// kodyWierszyModulow wylicza kody z wierszy słownika.
func kodyWierszyModulow(moduly []dane.Modul) []string {
	kody := make([]string, 0, len(moduly))
	for _, modul := range moduly {
		kody = append(kody, modul.Kod)
	}
	return kody
}

// kodyModulow wylicza kody z modułów kontraktu.
func kodyModulow(moduly []shared.Module) []string {
	kody := make([]string, 0, len(moduly))
	for _, modul := range moduly {
		kody = append(kody, modul.Code)
	}
	return kody
}
