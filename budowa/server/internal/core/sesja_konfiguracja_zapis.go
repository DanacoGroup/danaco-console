// Plik czyta i zapisuje obszary konfiguracji sesji na jednym adresie: poziom zasięgu i oś. Obszary leżą w tabeli ustawienie pod adresem złożonym, jak klucze proste rodziny config.*. Obszar skasowany wraca do dziedziczenia poziomu szerszego.
package core

import (
	"context"
	"encoding/json"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/shared"
)

// OdczytajKonfiguracjeSesji zwraca obszary zapisane pod wskazanym adresem, bez
// rozstrzygania dziedziczenia. Brak zapisu daje konfigurację pustą i pusty
// wykaz obszarów, nie odmowę.
func (a *adapterUstawienOsi) OdczytajKonfiguracjeSesji(ctx context.Context,
	z shared.ConfigSessionGetRequest) (shared.ConfigSessionGetResponse, error) {

	obszary, err := obszaryZadania(z.Areas)
	if err != nil {
		return shared.ConfigSessionGetResponse{}, err
	}
	adres := adresKonfiguracjiSesji(z.Scope, z.ScopeId, z.Axis, z.AxisId)
	zapisane, chwila, err := a.obszaryAdresu(ctx, adres, obszary)
	if err != nil {
		return shared.ConfigSessionGetResponse{}, err
	}
	konfiguracja, err := zlozKonfiguracje(zapisane)
	if err != nil {
		return shared.ConfigSessionGetResponse{}, err
	}
	wynik := shared.ConfigSessionGetResponse{
		Config:      konfiguracja,
		StoredAreas: obszaryZapisane(zapisane),
	}
	if chwila > 0 {
		wynik.UpdatedAt = &chwila
	}
	return wynik, nil
}

// ZapiszKonfiguracjeSesji zapisuje wskazane obszary pod wskazanym adresem
// i oddaje stan poziomu po zapisie wraz z wpisami rezolwera.
func (a *adapterUstawienOsi) ZapiszKonfiguracjeSesji(ctx context.Context,
	z shared.ConfigSessionSetRequest) (shared.ConfigSessionSetResponse, error) {

	obszary, err := obszaryZapisu(z.Areas)
	if err != nil {
		return shared.ConfigSessionSetResponse{}, err
	}
	tresci, err := rozbierzKonfiguracje(z.Config)
	if err != nil {
		return shared.ConfigSessionSetResponse{}, err
	}
	poziom := z.Scope
	adres := adresKonfiguracjiSesji(&poziom, z.ScopeId, z.Axis, z.AxisId)

	wpisy := make([]shared.ConfigEntry, 0, len(obszary))
	for _, obszar := range obszary {
		ustawienie := ustawienieObszaru(adres, obszar, tresci[obszar])
		if err := a.utrwalObszar(ctx, adres, obszar, tresci[obszar]); err != nil {
			return shared.ConfigSessionSetResponse{}, err
		}
		wpisy = append(wpisy, wpisOsi(ustawienie))
	}

	zapisane, _, err := a.obszaryAdresu(ctx, adres, obszaryKontraktu)
	if err != nil {
		return shared.ConfigSessionSetResponse{}, err
	}
	konfiguracja, err := zlozKonfiguracje(zapisane)
	if err != nil {
		return shared.ConfigSessionSetResponse{}, err
	}
	return shared.ConfigSessionSetResponse{
		Config:      konfiguracja,
		StoredAreas: obszaryZapisane(zapisane),
		Entries:     wpisy,
	}, nil
}

// obszaryAdresu czyta obszary spod jednego adresu i zwraca ich treści wraz
// z czasem najświeższego zapisu.
func (a *adapterUstawienOsi) obszaryAdresu(ctx context.Context, adres konfig.Adres,
	obszary []shared.SessionConfigArea) (map[shared.SessionConfigArea]json.RawMessage, int64, error) {

	tresci := map[shared.SessionConfigArea]json.RawMessage{}
	if a == nil || a.repozytorium == nil {
		return tresci, 0, nil
	}
	wiersze, err := a.repozytorium.ListaOsi(ctx, adres.Poziom, adres.KluczZasiegu, adres.Os, adres.KluczOsi)
	if err != nil {
		return nil, 0, err
	}
	wybrane := zbiorObszarow(obszary)
	najswiezszy := int64(0)
	for _, wiersz := range wiersze {
		obszar, znany := obszarKlucza(wiersz.Klucz)
		if !znany {
			continue
		}
		if _, zadany := wybrane[obszar]; !zadany {
			continue
		}
		tresc := wartoscTekstu(wiersz.Wartosc)
		if tresc == "" || !json.Valid([]byte(tresc)) {
			continue
		}
		tresci[obszar] = json.RawMessage(tresc)
		if chwila, czytelna := chwilaZapisu(wiersz.Zaktualizowano); czytelna && chwila > najswiezszy {
			najswiezszy = chwila
		}
	}
	return tresci, najswiezszy, nil
}

// utrwalObszar zapisuje treść obszaru albo — gdy żądanie jej nie niesie —
// kasuje obszar z tego poziomu, oddając go dziedziczeniu.
func (a *adapterUstawienOsi) utrwalObszar(ctx context.Context, adres konfig.Adres,
	obszar shared.SessionConfigArea, tresc json.RawMessage) error {

	if a == nil || a.repozytorium == nil {
		return bladBrakuKatalogu("konfiguracji sesji")
	}
	if len(tresc) == 0 {
		return a.repozytorium.UsunOsi(ctx, adres.Poziom, adres.KluczZasiegu,
			adres.Os, adres.KluczOsi, kluczObszaru(obszar))
	}
	return a.repozytorium.Ustaw(ctx, ustawienieObszaru(adres, obszar, tresc))
}

// ustawienieObszaru buduje wiersz tabeli ustawień niosący treść obszaru pod adresem złożonym z poziomu, osi i klucza obszaru.
func ustawienieObszaru(adres konfig.Adres, obszar shared.SessionConfigArea,
	tresc json.RawMessage) dane.Ustawienie {

	wartosc := string(tresc)
	return dane.Ustawienie{
		Poziom: adres.Poziom, KluczZasiegu: adres.KluczZasiegu,
		Os: adres.Os, KluczOsi: adres.KluczOsi,
		Klucz: kluczObszaru(obszar), Wartosc: &wartosc,
		RodzajWartosci: string(konfig.RodzajJSON),
	}
}

// adresKonfiguracjiSesji składa adres zapisu z pól żądania. Poziom pominięty
// znaczy globalny, oś pominięta — platformę.
func adresKonfiguracjiSesji(poziom *shared.ConfigScope, bytPoziomu *string,
	os *shared.ConfigAxis, bytOsi *string) konfig.Adres {

	wybrany := shared.ConfigScope(shared.ConfigScopeGlobal)
	if poziom != nil && *poziom != "" {
		wybrany = *poziom
	}
	osZasiegu, kluczOsi := osZadania(os, bytOsi)
	return konfig.Adres{
		Poziom: wybrany, KluczZasiegu: wartoscTekstu(bytPoziomu),
		Os: osZasiegu, KluczOsi: kluczOsi,
	}
}

// obszaryZapisu zawęża wykaz obszarów objętych zapisem. Wykaz pusty znaczy
// tutaj „żadnego obszaru”, nie „komplet”: domyślne objęcie wszystkiego
// skasowałoby poziom.
func obszaryZapisu(zadane []shared.SessionConfigArea) ([]shared.SessionConfigArea, error) {
	if len(zadane) == 0 {
		return nil, nil
	}
	return obszaryZadania(zadane)
}

// obszaryZapisane wylicza obszary obecne w mapie treści zapisu, zachowując kolejność pól kontraktu, a nie kolejność wstawienia do mapy.
func obszaryZapisane(tresci map[shared.SessionConfigArea]json.RawMessage) []shared.SessionConfigArea {
	obecne := make([]shared.SessionConfigArea, 0, len(tresci))
	for _, obszar := range obszaryKontraktu {
		if _, jest := tresci[obszar]; jest {
			obecne = append(obecne, obszar)
		}
	}
	return obecne
}

// zbiorObszarow składa zbiór obszarów do szybkiego sprawdzenia przynależności obszaru do wykazu obszarów zapisu.
func zbiorObszarow(obszary []shared.SessionConfigArea) map[shared.SessionConfigArea]struct{} {
	zbior := make(map[shared.SessionConfigArea]struct{}, len(obszary))
	for _, obszar := range obszary {
		zbior[obszar] = struct{}{}
	}
	return zbior
}
