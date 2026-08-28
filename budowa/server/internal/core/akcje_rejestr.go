package core

import (
	"context"
	"sync"

	"danacoconsole/shared"
)

// RejestrAkcji trzyma katalog akcji budowany w czasie działania z wierszy tabeli akcja.
// Ani jedna akcja nie jest wpisana w kod: dopisanie akcji to dopisanie wiersza i odświeżenie
// rejestru.
type RejestrAkcji struct {
	mu      sync.RWMutex
	zrodlo  ZrodloAkcji
	pozycje []AkcjaKatalogu
	czytany bool
}

// NowyRejestrAkcji zakłada pusty rejestr nad źródłem wierszy. Konstruktor
// niczego nie odpytuje — katalog powstaje dopiero przy Odswiez albo przy
// pierwszym odczycie.
func NowyRejestrAkcji(zrodlo ZrodloAkcji) *RejestrAkcji {
	return &RejestrAkcji{zrodlo: zrodlo}
}

// Odswiez czyta wiersze na nowo i podmienia katalog. Niepowodzenie odczytu nie
// kasuje katalogu poprzedniego — rejestr woli treść starszą niż żadną.
func (r *RejestrAkcji) Odswiez(ctx context.Context) error {
	if r == nil {
		return nil
	}
	r.mu.RLock()
	zrodlo := r.zrodlo
	r.mu.RUnlock()
	if zrodlo == nil {
		r.oznaczCzytany()
		return nil
	}
	wiersze, err := zrodlo.Akcje(ctx)
	if err != nil {
		return err
	}
	pozycje := make([]AkcjaKatalogu, 0, len(wiersze))
	for _, wiersz := range wiersze {
		pozycje = append(pozycje, akcjaKatalogu(wiersz))
	}
	r.mu.Lock()
	r.pozycje, r.czytany = pozycje, true
	r.mu.Unlock()
	return nil
}

// Pozycje zwraca cały katalog akcji w kolejności odczytu z bazy danych rdzenia, wraz
// z wierszami nieczynnymi.
func (r *RejestrAkcji) Pozycje() []AkcjaKatalogu {
	if r == nil {
		return nil
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]AkcjaKatalogu(nil), r.pozycje...)
}

// Zasieg zwraca akcje jednego zasięgu. Pusty poziom zwraca katalog w całości. Wskazanie
// bytu poziomu dokłada do jego akcji własnych akcje wspólne całemu poziomowi.
func (r *RejestrAkcji) Zasieg(poziom shared.ConfigScope, kluczZasiegu string,
	tylkoCzynne bool) []AkcjaKatalogu {

	pozycje := r.Pozycje()
	wybrane := make([]AkcjaKatalogu, 0, len(pozycje))
	for _, pozycja := range pozycje {
		if tylkoCzynne && !pozycja.Enabled {
			continue
		}
		if !pasujeDoZasiegu(pozycja, poziom, kluczZasiegu) {
			continue
		}
		wybrane = append(wybrane, pozycja)
	}
	return wybrane
}

// Wykaz wypełnia port Akcje. Katalog pusty jest poprawną odpowiedzią: brak wierszy nie może
// zatrzymać ani panelu akcji, ani narzędzi modelu.
func (r *RejestrAkcji) Wykaz(ctx context.Context, z ZadanieKatalogAkcji) (WynikKatalogAkcji, error) {
	if r == nil {
		return WynikKatalogAkcji{Actions: []AkcjaKatalogu{}}, nil
	}
	if !r.wczytany() {
		_ = r.Odswiez(ctx)
	}
	poziom := shared.ConfigScope("")
	if z.Scope != nil {
		poziom = *z.Scope
	}
	klucz := ""
	if z.ScopeId != nil {
		klucz = *z.ScopeId
	}
	tylkoCzynne := z.EnabledOnly == nil || *z.EnabledOnly
	return WynikKatalogAkcji{Actions: r.Zasieg(poziom, klucz, tylkoCzynne)}, nil
}

// pasujeDoZasiegu rozstrzyga przynależność pozycji katalogu do zapytanego zasięgu oraz
// jego klucza bytu.
func pasujeDoZasiegu(pozycja AkcjaKatalogu, poziom shared.ConfigScope, kluczZasiegu string) bool {
	if poziom == "" {
		return true
	}
	if pozycja.Scope != poziom {
		return false
	}
	if kluczZasiegu == "" || pozycja.ScopeId == nil || *pozycja.ScopeId == "" {
		return true
	}
	return *pozycja.ScopeId == kluczZasiegu
}

// wczytany mówi, czy rejestr próbował już odczytać wiersze katalogu z bazy przy poprzednim
// wywołaniu metody.
func (r *RejestrAkcji) wczytany() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.czytany
}

// oznaczCzytany zamyka próbę odczytu rejestru bez źródła — bez tego każdy
// odczyt katalogu ponawiałby budowę, której nie ma z czego wykonać.
func (r *RejestrAkcji) oznaczCzytany() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.czytany = true
}
