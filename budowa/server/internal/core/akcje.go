package core

import (
	"context"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// Kształt katalogu akcji ma jedno źródło prawdy — shared/contract.json. Typy poniżej są
// aliasami kształtów wygenerowanych z kontraktu, więc zmiana kontraktu przerywa kompilację
// rdzenia, zamiast rozjeżdżać go po cichu.
type (
	// AkcjaKatalogu to pozycja katalogu akcji odsyłana klientowi i kanałowi modelu, jako element
	// listy wykazu.
	AkcjaKatalogu = shared.Action
	// ZadanieKatalogAkcji jest żądaniem odczytu katalogu akcji. Pominięty poziom
	// zasięgu zwraca katalog w całości — brak zawężenia nie jest błędem.
	ZadanieKatalogAkcji = shared.ActionListRequest
	// WynikKatalogAkcji jest odpowiedzią odczytu katalogu akcji, niosącą wykaz pozycji całego
	// katalogu akcji.
	WynikKatalogAkcji = shared.ActionListResponse
)

// Akcje obsługuje katalog akcji sterowany danymi: nowa akcja to nowy
// wiersz katalogu, nie nowy przycisk w kodzie widoku ani nowa gałąź w rdzeniu.
// Ten sam katalog zasila panel akcji i narzędzia kanału modelu.
type Akcje interface {
	Wykaz(ctx context.Context, z ZadanieKatalogAkcji) (WynikKatalogAkcji, error)
}

// ZrodloAkcji podaje rejestrowi wiersze katalogu. Rdzeń nie wie, skąd pochodzą
// — w montażu jest to repozytorium warstwy danych, w teście tablica.
type ZrodloAkcji interface {
	Akcje(ctx context.Context) ([]dane.Akcja, error)
}

// zrodloAkcjiFunkcja pozwala podać źródło katalogu samą funkcją, bez tworzenia osobnego
// typu pomocniczego.
type zrodloAkcjiFunkcja func(ctx context.Context) ([]dane.Akcja, error)

// Akcje wypełnia interfejs ZrodloAkcji zwykłą funkcją, adaptując ją do kształtu wymaganego
// przez rejestr.
func (f zrodloAkcjiFunkcja) Akcje(ctx context.Context) ([]dane.Akcja, error) {
	return f(ctx)
}

// ZrodloAkcjiZRepozytorium wiąże rejestr z katalogiem w bazie. Odczytuje
// komplet wierszy, także nieczynne — o ograniczeniu do czynnych rozstrzyga
// żądanie, nie zapytanie SQL.
func ZrodloAkcjiZRepozytorium(repozytorium dane.RepozytoriumAkcji) ZrodloAkcji {
	return zrodloAkcjiFunkcja(func(ctx context.Context) ([]dane.Akcja, error) {
		if repozytorium == nil {
			return nil, nil
		}
		return repozytorium.Lista(ctx, false)
	})
}

// akcjaKatalogu przekłada wiersz katalogu na pozycję kształtu odpowiedzi.
// Pola opcjonalne kontraktu wychodzą puste, gdy wiersz ich nie wypełnia —
// pusty opis nie jest brakiem danych.
func akcjaKatalogu(a dane.Akcja) AkcjaKatalogu {
	pozycja := AkcjaKatalogu{
		Id: a.Kod, Name: a.Nazwa, Scope: a.PoziomZasiegu,
		Command: string(a.Komenda), Order: a.Kolejnosc, Enabled: a.Aktywna,
	}
	pozycja.Description = tekstOpcjonalny(a.Opis)
	pozycja.Icon = tekstOpcjonalny(a.Ikona)
	pozycja.ScopeId = tekstOpcjonalny(a.KluczZasiegu)
	pozycja.Requires = tekstOpcjonalny(a.WarunekDostepnosci)
	return pozycja
}

// tekstOpcjonalny przenosi wartość tekstową do pola opcjonalnego kontraktu.
// Tekst pusty daje brak wartości, a nie wskaźnik na pusty łańcuch.
func tekstOpcjonalny(wartosc string) *string {
	if wartosc == "" {
		return nil
	}
	return &wartosc
}
