package core

import (
	"context"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/shared"
)

// adapterUstawien wypełnia port Ustawienia dwiema warstwami: repozytorium konfiguracji
// zapisuje i czyta wiersze, a rozstrzygacz dziewięciu poziomów zasięgu odpowiada na pytanie
// o wartość obowiązującą.
type adapterUstawien struct {
	repozytorium dane.RepozytoriumKonfiguracji
	rozstrzygacz *konfig.Rozstrzygacz
}

// nowyAdapterUstawien wiąże port z repozytorium konfiguracji i rozstrzygaczem poziomów
// zasięgu, tworząc adapter gotowy do obsługi komend ustawień.
func nowyAdapterUstawien(repozytorium dane.RepozytoriumKonfiguracji, rozstrzygacz *konfig.Rozstrzygacz) *adapterUstawien {
	return &adapterUstawien{repozytorium: repozytorium, rozstrzygacz: rozstrzygacz}
}

// Odczytaj zwraca wpisy wskazanego poziomu zasięgu, a bez wskazania poziomu —
// politykę efektywną, czyli wartość obowiązującą wraz z poziomem, z którego
// pochodzi.
func (a *adapterUstawien) Odczytaj(ctx context.Context, z shared.ConfigGetRequest) (shared.ConfigGetResponse, error) {
	if z.Scope == nil {
		return shared.ConfigGetResponse{
			Entries: a.rozstrzygacz.PolitykaEfektywna(kontekstZasiegu(z.ScopeId)).WpisyKontraktu(),
		}, nil
	}
	wpisy, err := a.repozytorium.ListaPoziomu(ctx, *z.Scope, wartoscTekstu(z.ScopeId))
	if err != nil {
		return shared.ConfigGetResponse{}, err
	}
	return shared.ConfigGetResponse{Entries: wpisyKontraktu(wpisy, z.Key)}, nil
}

// Zapisz ustawia wartość na wskazanym poziomie zasięgu, zapisując nowy wiersz konfiguracji
// w repozytorium.
func (a *adapterUstawien) Zapisz(ctx context.Context, z shared.ConfigSetRequest) (shared.ConfigSetResponse, error) {
	wartosc, rodzaj := konfig.DekodujJSON(z.Value)
	tekst := wartosc
	ustawienie := dane.Ustawienie{
		Poziom: z.Scope, KluczZasiegu: wartoscTekstu(z.ScopeId), Klucz: z.Key,
		Wartosc: &tekst, RodzajWartosci: string(rodzaj),
	}
	if err := a.repozytorium.Ustaw(ctx, ustawienie); err != nil {
		return shared.ConfigSetResponse{}, err
	}
	return shared.ConfigSetResponse{Entry: wpisKontraktu(ustawienie)}, nil
}

// Przywroc usuwa ustawienie z poziomu. Brak ustawienia znaczy wartość domyślną,
// więc usunięcie wpisu jest przywróceniem wartości domyślnej.
func (a *adapterUstawien) Przywroc(ctx context.Context, z shared.ConfigResetRequest) (shared.ConfigResetResponse, error) {
	usuwane, err := a.repozytorium.ListaPoziomu(ctx, z.Scope, wartoscTekstu(z.ScopeId))
	if err != nil {
		return shared.ConfigResetResponse{}, err
	}
	wpisy := wpisyKontraktu(usuwane, z.Key)
	for _, wpis := range wpisy {
		if err := a.repozytorium.Usun(ctx, z.Scope, wartoscTekstu(z.ScopeId), wpis.Key); err != nil {
			return shared.ConfigResetResponse{}, err
		}
	}
	return shared.ConfigResetResponse{Entries: wpisy}, nil
}

// kontekstZasiegu buduje kontekst rozstrzygania. Wskazany byt trafia na poziom
// okna — najwęższy z dziewięciu; jego brak daje kontekst globalny.
func kontekstZasiegu(idBytu *string) konfig.Kontekst {
	return konfig.Kontekst{Okno: wartoscTekstu(idBytu)}
}

// wpisyKontraktu przekłada wiersze repozytorium na wpisy kontraktu, zawężając
// je kluczem, gdy klucz wskazano.
func wpisyKontraktu(ustawienia []dane.Ustawienie, klucz *string) []shared.ConfigEntry {
	wpisy := make([]shared.ConfigEntry, 0, len(ustawienia))
	for _, ustawienie := range ustawienia {
		if klucz != nil && *klucz != "" && ustawienie.Klucz != *klucz {
			continue
		}
		wpisy = append(wpisy, wpisKontraktu(ustawienie))
	}
	return wpisy
}

// wpisKontraktu przekłada jeden wiersz repozytorium na wpis kontraktu, dołączając kod
// zasięgu tylko dla wpisu, który go ma.
func wpisKontraktu(u dane.Ustawienie) shared.ConfigEntry {
	wpis := shared.ConfigEntry{
		Key:   u.Klucz,
		Value: konfig.KodujJSON(wartoscTekstu(u.Wartosc), konfig.Rodzaj(u.RodzajWartosci)),
		Scope: u.Poziom,
	}
	if u.KluczZasiegu != "" {
		klucz := u.KluczZasiegu
		wpis.ScopeId = &klucz
	}
	return wpis
}

// wartoscTekstu odczytuje pole opcjonalne kontraktu, zwracając pusty napis, gdy pole nie
// niesie wartości.
func wartoscTekstu(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}
