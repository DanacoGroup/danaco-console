// Plik obsługuje odczyt i zapis konfiguracji pod adresem złożonym z poziomu
// zasięgu oraz osi rozstrzygania, przez port zgodny z adapterem podstawowym
// ustawień, bez powielania kodowania wartości.
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

type adapterUstawienOsi struct {
	podstawa     *adapterUstawien
	repozytorium dane.RepozytoriumKonfiguracjiOsi
	rozstrzygacz *konfig.Rozstrzygacz
}

func nowyAdapterUstawienOsi(repozytorium dane.RepozytoriumKonfiguracjiOsi,
	rozstrzygacz *konfig.Rozstrzygacz) *adapterUstawienOsi {

	return &adapterUstawienOsi{
		podstawa:     nowyAdapterUstawien(repozytorium, rozstrzygacz),
		repozytorium: repozytorium,
		rozstrzygacz: rozstrzygacz,
	}
}

func (a *adapterUstawienOsi) Odczytaj(ctx context.Context, z shared.ConfigGetRequest) (shared.ConfigGetResponse, error) {
	os, bytOsi := osZadania(z.Axis, z.AxisId)
	if osPlatformy(os, bytOsi) {
		return a.podstawa.Odczytaj(ctx, z)
	}
	if z.Scope == nil {
		kontekst := kontekstOsi(ctx, z.ScopeId, os, bytOsi)
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

func (a *adapterUstawienOsi) Zapisz(ctx context.Context, z shared.ConfigSetRequest) (shared.ConfigSetResponse, error) {
	// Klucz spoza katalogu odpada tutaj, przed rozdziałem zapisu na osie.
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

func osZadania(os *shared.ConfigAxis, bytOsi *string) (shared.ConfigAxis, string) {
	if os == nil {
		return konfig.OsPlatformy, ""
	}
	return konfig.OsLubPlatforma(*os), wartoscTekstu(bytOsi)
}

func osPlatformy(os shared.ConfigAxis, bytOsi string) bool {
	return konfig.OsLubPlatforma(os) == konfig.OsPlatformy || bytOsi == ""
}

func kontekstOsi(ctx context.Context, idBytu *string, os shared.ConfigAxis, bytOsi string) konfig.Kontekst {
	kontekst := kontekstZasiegu(ctx, idBytu)
	switch konfig.OsLubPlatforma(os) {
	case konfig.OsModelu:
		kontekst.Model = bytOsi
	case konfig.OsKonta:
		kontekst.Konto = bytOsi
	}
	return kontekst
}

func ustawienieZWpisu(wpis konfig.Wpis) dane.Ustawienie {
	wartosc := wpis.Wartosc
	return dane.Ustawienie{
		Poziom: wpis.Poziom, KluczZasiegu: wpis.KluczZasiegu,
		Os: wpis.Os, KluczOsi: wpis.KluczOsi,
		Klucz: wpis.Klucz, Wartosc: &wartosc, RodzajWartosci: string(wpis.Rodzaj),
	}
}

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

func bladZapisuOsi(klucz string, adres konfig.Adres) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"konfiguracja: klucza "+klucz+" nie da się zapisać na poziomie "+
			string(adres.Poziom)+" osi "+string(adres.Os)))
}
