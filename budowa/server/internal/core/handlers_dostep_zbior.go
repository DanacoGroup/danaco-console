// Plik odczytuje komplet nadań jednego okna i nakłada na nadanie pola
// wskazane w żądaniu zmiany, z pamięcią podręczną kodu punktu.
package core

import (
	"context"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// zbiorOkna zwraca komplet nadań okna przełożony na struktury kontraktu,
// w kolejności zapisanej w bazie.
func (a *adapterNadanDostepu) zbiorOkna(ctx context.Context, okno dane.Okno,
	tylkoAktywne bool) ([]shared.AccessGrant, error) {

	wiersze, err := a.nadania.ListaOkna(ctx, okno.ID, tylkoAktywne)
	if err != nil {
		return nil, err
	}
	idOkna := identyfikatorWiersza(okno.IdentyfikatorZewnetrzny, okno.ID)
	kody := map[int64]string{}
	zbior := make([]shared.AccessGrant, 0, len(wiersze))
	for _, wiersz := range wiersze {
		kod, err := a.kodPunktu(ctx, kody, wiersz.PunktDostepuID)
		if err != nil {
			return nil, err
		}
		zbior = append(zbior, nadanieKontraktu(wiersz, kod, idOkna))
	}
	return zbior, nil
}

// kodPunktu zwraca trwały kod punktu, korzystając z pamięci podręcznej odczytu.
// Punkt skasowany w międzyczasie daje kod pusty — zbiór nadań ma się pokazać,
// a nie zniknąć z powodu jednego osieroconego wiersza.
func (a *adapterNadanDostepu) kodPunktu(ctx context.Context, kody map[int64]string,
	punktID int64) (string, error) {

	if kod, jest := kody[punktID]; jest {
		return kod, nil
	}
	punkt, err := a.punkty.Pobierz(ctx, punktID)
	if brakWiersza(err) {
		kody[punktID] = ""
		return "", nil
	}
	if err != nil {
		return "", err
	}
	kody[punktID] = punkt.Kod
	return punkt.Kod, nil
}

// zeZbioru odnajduje w zbiorze nadanie o wskazanym identyfikatorze. Brak
// nadania w zbiorze daje strukturę pustą zamiast paniki — zbiór jest odczytem
// po zapisie i może się rozminąć z zapisem tylko przy równoczesnej zmianie.
func zeZbioru(zbior []shared.AccessGrant, identyfikator string) shared.AccessGrant {
	for _, nadanie := range zbior {
		if nadanie.Id == identyfikator {
			return nadanie
		}
	}
	return shared.AccessGrant{Id: identyfikator}
}

// zastosujZmianeNadania nakłada na wiersz pola wskazane w żądaniu. Oznaczenia
// głównego tu nie ma — należy do zbioru okna, nie do wiersza, więc zapisuje je
// osobna czynność repozytorium.
func zastosujZmianeNadania(nadanie *dane.Nadanie, z shared.AccessGrantUpdateRequest) {
	if z.Mode != nil {
		nadanie.Tryb = *z.Mode
	}
	if z.Roots != nil {
		nadanie.Korzenie = z.Roots
	}
	if z.Order != nil {
		nadanie.Kolejnosc = *z.Order
	}
	if z.Enabled != nil {
		nadanie.Aktywne = *z.Enabled
	}
}

// liczbaLub odczytuje pole liczbowe opcjonalne kontraktu, oddając wartość
// domyślną, gdy pole jest puste.
func liczbaLub(wartosc *int, domyslna int) int {
	if wartosc == nil {
		return domyslna
	}
	return *wartosc
}
