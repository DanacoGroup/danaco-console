// Plik wpina cztery komendy jednolitego modelu konfiguracji sesji: odczyt zapisu poziomu, zapis
// obszarów, konfigurację obowiązującą oraz deklarację zdolności adaptera dostawcy.
package core

import (
	"context"

	"danacoconsole/shared"
)

// KonfiguracjaSesji jest portem jednolitego modelu konfiguracji sesji. Mówi
// wyłącznie typami kontraktu; wiersze tabeli ustawień i wykaz obszarów leżą po
// drugiej stronie adaptera.
type KonfiguracjaSesji interface {
	// OdczytajKonfiguracjeSesji zwraca obszary zapisane na wskazanym poziomie i osi, bez dziedziczenia.
	OdczytajKonfiguracjeSesji(ctx context.Context,
		z shared.ConfigSessionGetRequest) (shared.ConfigSessionGetResponse, error)
	// ZapiszKonfiguracjeSesji zapisuje wskazane obszary; obszar spoza wykazu
	// pozostaje nietknięty.
	ZapiszKonfiguracjeSesji(ctx context.Context,
		z shared.ConfigSessionSetRequest) (shared.ConfigSessionSetResponse, error)
	// KonfiguracjaObowiazujaca rozstrzyga obszary po wszystkich poziomach i osiach z pochodzeniem.
	KonfiguracjaObowiazujaca(ctx context.Context,
		z shared.ConfigEffectiveGetRequest) (shared.ConfigEffectiveGetResponse, error)
	// ZdolnosciAdaptera zwraca deklarację zdolności adaptera dostawcy.
	ZdolnosciAdaptera(ctx context.Context,
		z shared.ConfigCapabilitiesGetRequest) (shared.ConfigCapabilitiesGetResponse, error)
}

// zarejestrujKonfiguracjeSesji wpina cztery komendy konfiguracji sesji.
//
// Brak portu nie rejestruje niczego: komendy odpowiedzą wtedy `config.unknown`,
// a pozostałe domeny pracują bez zmian.
func zarejestrujKonfiguracjeSesji(r *Rejestr, konfiguracja KonfiguracjaSesji, e *emiter) {
	if r == nil || konfiguracja == nil {
		return
	}

	r.Zarejestruj(shared.CommandConfigSessionGet, obsluz(konfiguracja.OdczytajKonfiguracjeSesji))
	r.Zarejestruj(shared.CommandConfigEffectiveGet, obsluz(konfiguracja.KonfiguracjaObowiazujaca))
	r.Zarejestruj(shared.CommandConfigCapabilitiesGet, obsluz(konfiguracja.ZdolnosciAdaptera))

	r.Zarejestruj(shared.CommandConfigSessionSet,
		obsluz(func(ctx context.Context,
			z shared.ConfigSessionSetRequest) (shared.ConfigSessionSetResponse, error) {

			w, err := konfiguracja.ZapiszKonfiguracjeSesji(ctx, z)
			if err != nil {
				return w, err
			}
			for _, wpis := range w.Entries {
				e.ustawienie(ctx, zmianaObszaru(w.StoredAreas, wpis), wpis)
			}
			return w, nil
		}))
}

// zmianaObszaru rozpoznaje, czy wpis opisuje obszar zapisany, czy wyczyszczony.
// Obszar, którego po zapisie nie ma już na poziomie, został skasowany — jego
// wartość bierze się odtąd z poziomu szerszego.
func zmianaObszaru(zapisane []shared.SessionConfigArea, wpis shared.ConfigEntry) shared.ChangeKind {
	obszar, znany := obszarKlucza(wpis.Key)
	if !znany {
		return shared.ChangeKindUpdated
	}
	for _, obecny := range zapisane {
		if obecny == obszar {
			return shared.ChangeKindUpdated
		}
	}
	return shared.ChangeKindDeleted
}
