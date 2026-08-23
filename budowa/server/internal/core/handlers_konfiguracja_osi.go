// Odpowiedzialność pliku: zapis i odczyt konfiguracji pod adresem złożonym —
// poziom zasięgu razem z osią rozstrzygania (pola `axis` i `axisId` żądań
// rodziny `config.*`).
//
// Oś jest prostopadła do poziomu: poziom mówi, jak wąsko obowiązuje wartość
// (globalny → okno komunikacji), oś mówi, dla czego obowiązuje — dla platformy,
// dla modelu albo dla konta. Osie są nośnikiem konfiguracji per model
// i per konto.
//
// Adapter nie powtarza adaptera ustawień, tylko go owija: oś platformy schodzi
// do `adapterUstawien` bez zmiany, a osobna jest wyłącznie droga osi modelu
// i konta. Przekład wiersza na wpis kontraktu należy w całości do pakietu
// `konfig` (`WpisKontraktu`, `ZapisZKontraktuWOsi`) — rdzeń nie ma tu własnego
// kodowania wartości.
package core

import (
	"context"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// Adapter osi wypełnia ten sam port co adapter podstawowy. Gdyby porty się
// rozjechały, kompilacja stanie tutaj.
var _ Ustawienia = (*adapterUstawienOsi)(nil)

// adapterUstawienOsi wypełnia port Ustawienia z obsługą osi rozstrzygania.
type adapterUstawienOsi struct {
	podstawa     *adapterUstawien
	repozytorium dane.RepozytoriumKonfiguracjiOsi
	rozstrzygacz *konfig.Rozstrzygacz
}

// nowyAdapterUstawienOsi wiąże port z repozytorium świadomym osi oraz
// z rozstrzygaczem ośmiu poziomów zasięgu.
func nowyAdapterUstawienOsi(repozytorium dane.RepozytoriumKonfiguracjiOsi,
	rozstrzygacz *konfig.Rozstrzygacz) *adapterUstawienOsi {

	return &adapterUstawienOsi{
		podstawa:     nowyAdapterUstawien(repozytorium, rozstrzygacz),
		repozytorium: repozytorium,
		rozstrzygacz: rozstrzygacz,
	}
}

// Odczytaj zwraca wpisy poziomu na wskazanej osi, a bez wskazania poziomu —
// politykę efektywną liczoną z uwzględnieniem osi.
func (a *adapterUstawienOsi) Odczytaj(ctx context.Context, z shared.ConfigGetRequest) (shared.ConfigGetResponse, error) {
	os, bytOsi := osZadania(z.Axis, z.AxisId)
	if osPlatformy(os, bytOsi) {
		return a.podstawa.Odczytaj(ctx, z)
	}
	if z.Scope == nil {
		kontekst := kontekstOsi(z.ScopeId, os, bytOsi)
		return shared.ConfigGetResponse{
			Entries: a.rozstrzygacz.PolitykaEfektywna(kontekst).WpisyKontraktu(),
		}, nil
	}
	ustawienia, err := a.repozytorium.ListaOsi(ctx, *z.Scope, wartoscTekstu(z.ScopeId), os, bytOsi)
	if err != nil {
		return shared.ConfigGetResponse{}, err
	}
	return shared.ConfigGetResponse{Entries: wpisyOsi(ustawienia, z.Key)}, nil
}

// Zapisz ustawia wartość pod adresem złożonym z poziomu zasięgu i osi.
func (a *adapterUstawienOsi) Zapisz(ctx context.Context, z shared.ConfigSetRequest) (shared.ConfigSetResponse, error) {
	// Klucz spoza katalogu odpada tutaj, przed rozdziałem na osie — jedna brama
	// dla obu dróg zapisu. Bez niej zapis nieznanego klucza przeszedłby bez
	// odmowy i został w bazie martwym wierszem.
	if err := sprawdzKluczKatalogu(a.rozstrzygacz, z.Key); err != nil {
		return shared.ConfigSetResponse{}, err
	}
	os, bytOsi := osZadania(z.Axis, z.AxisId)
	if osPlatformy(os, bytOsi) {
		return a.podstawa.Zapisz(ctx, z)
	}
	adres := konfig.Adres{Poziom: z.Scope, KluczZasiegu: wartoscTekstu(z.ScopeId), Os: os, KluczOsi: bytOsi}
	wpis, poprawny := konfig.ZapisZKontraktuWOsi(z.Key, z.Value, adres)
	if !poprawny {
		return shared.ConfigSetResponse{}, bladZapisuOsi(z.Key, adres)
	}
	ustawienie := ustawienieZWpisu(wpis)
	if err := a.repozytorium.Ustaw(ctx, ustawienie); err != nil {
		return shared.ConfigSetResponse{}, err
	}
	return shared.ConfigSetResponse{Entry: wpisOsi(ustawienie)}, nil
}

// Przywroc usuwa ustawienia spod adresu złożonego. Brak zapisu znaczy wartość
// poziomu szerszego, a w ostateczności wartość domyślną.
func (a *adapterUstawienOsi) Przywroc(ctx context.Context, z shared.ConfigResetRequest) (shared.ConfigResetResponse, error) {
	os, bytOsi := osZadania(z.Axis, z.AxisId)
	if osPlatformy(os, bytOsi) {
		return a.podstawa.Przywroc(ctx, z)
	}
	zasieg := wartoscTekstu(z.ScopeId)
	usuwane, err := a.repozytorium.ListaOsi(ctx, z.Scope, zasieg, os, bytOsi)
	if err != nil {
		return shared.ConfigResetResponse{}, err
	}
	wpisy := wpisyOsi(usuwane, z.Key)
	for _, wpis := range wpisy {
		if err := a.repozytorium.UsunOsi(ctx, z.Scope, zasieg, os, bytOsi, wpis.Key); err != nil {
			return shared.ConfigResetResponse{}, err
		}
	}
	return shared.ConfigResetResponse{Entries: wpisy}, nil
}

// osZadania odczytuje oś i jej byt z pól opcjonalnych żądania. Oś pominięta
// znaczy platformę.
func osZadania(os *shared.ConfigAxis, bytOsi *string) (shared.ConfigAxis, string) {
	if os == nil {
		return konfig.OsPlatformy, ""
	}
	return konfig.OsLubPlatforma(*os), wartoscTekstu(bytOsi)
}

// osPlatformy odpowiada, czy adres opisuje oś platformy — czyli drogę sprzed
// wprowadzenia osi, obsługiwaną przez adapter podstawowy bez zmiany.
func osPlatformy(os shared.ConfigAxis, bytOsi string) bool {
	return konfig.OsLubPlatforma(os) == konfig.OsPlatformy || bytOsi == ""
}

// kontekstOsi buduje kontekst rozstrzygania dla osi modelu albo konta. Wskazany
// byt poziomu trafia na poziom okna — najwęższy z ośmiu.
func kontekstOsi(idBytu *string, os shared.ConfigAxis, bytOsi string) konfig.Kontekst {
	kontekst := kontekstZasiegu(idBytu)
	switch konfig.OsLubPlatforma(os) {
	case konfig.OsModelu:
		kontekst.Model = bytOsi
	case konfig.OsKonta:
		kontekst.Konto = bytOsi
	}
	return kontekst
}

// ustawienieZWpisu przenosi wpis pakietu konfiguracji na wiersz repozytorium.
func ustawienieZWpisu(wpis konfig.Wpis) dane.Ustawienie {
	wartosc := wpis.Wartosc
	return dane.Ustawienie{
		Poziom: wpis.Poziom, KluczZasiegu: wpis.KluczZasiegu,
		Os: wpis.Os, KluczOsi: wpis.KluczOsi,
		Klucz: wpis.Klucz, Wartosc: &wartosc, RodzajWartosci: string(wpis.Rodzaj),
	}
}

// wpisyOsi przekłada wiersze repozytorium na wpisy kontraktu, zawężając je
// kluczem, gdy klucz wskazano.
func wpisyOsi(ustawienia []dane.Ustawienie, klucz *string) []shared.ConfigEntry {
	wpisy := make([]shared.ConfigEntry, 0, len(ustawienia))
	for _, ustawienie := range ustawienia {
		if klucz != nil && *klucz != "" && ustawienie.Klucz != *klucz {
			continue
		}
		wpisy = append(wpisy, wpisOsi(ustawienie))
	}
	return wpisy
}

// wpisOsi przekłada jeden wiersz na wpis kontraktu. Kodowanie wartości i zapis
// osi w kopercie należą do pakietu konfig — tutaj nie ma drugiego przekładu.
func wpisOsi(u dane.Ustawienie) shared.ConfigEntry {
	return konfig.Wynik{
		Klucz:        u.Klucz,
		Wartosc:      wartoscTekstu(u.Wartosc),
		Rodzaj:       konfig.Rodzaj(u.RodzajWartosci),
		Poziom:       u.Poziom,
		KluczZasiegu: u.KluczZasiegu,
		Os:           u.Os,
		KluczOsi:     u.KluczOsi,
	}.WpisKontraktu()
}

// bladZapisuOsi odmawia zapisu pod adresem spoza kontraktu. To odmowa
// merytoryczna dotycząca jednego wywołania, nie awaria rdzenia.
func bladZapisuOsi(klucz string, adres konfig.Adres) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"konfiguracja: klucza "+klucz+" nie da się zapisać na poziomie "+
			string(adres.Poziom)+" osi "+string(adres.Os)))
}
