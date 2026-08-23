// Odpowiedzialność pliku: wpięcie czterech komend jednolitego modelu
// konfiguracji sesji — odczytu zapisu poziomu, zapisu obszarów, konfiguracji
// obowiązującej oraz deklaracji zdolności adaptera dostawcy.
//
// To jest ta sama konfiguracja co rodzina `config.*`, nie drugi rejestr obok
// niej: konfiguracja sesji jest innym kształtem wartości w tym samym rejestrze
// ośmiu poziomów zasięgu i trzech osi. Dlatego port konfiguracji sesji wchodzi
// w skład portu Ustawienia, a nie obok niego — rdzeń ma jedną bramę do
// konfiguracji i jedno miejsce rozgłaszania `config.changed`.
//
// Zapis rozgłasza zmianę obszaru. Każdy dotknięty obszar wraca wpisem rezolwera
// i idzie zdarzeniem `config.changed` — obszar zapisany jako zmieniony, obszar
// wyczyszczony jako usunięty, bo skasowanie zapisu przywraca dziedziczenie.
// Klient nie musi odpytywać poziomu, żeby dowiedzieć się o zmianie.
package core

import (
	"context"

	"danacoconsole/shared"
)

// KonfiguracjaSesji jest portem jednolitego modelu konfiguracji sesji. Mówi
// wyłącznie typami kontraktu; wiersze tabeli ustawień i wykaz obszarów leżą po
// drugiej stronie adaptera.
type KonfiguracjaSesji interface {
	// OdczytajKonfiguracjeSesji zwraca obszary zapisane na wskazanym poziomie
	// zasięgu i wskazanej osi, bez rozstrzygania dziedziczenia.
	OdczytajKonfiguracjeSesji(ctx context.Context,
		z shared.ConfigSessionGetRequest) (shared.ConfigSessionGetResponse, error)
	// ZapiszKonfiguracjeSesji zapisuje wskazane obszary; obszar spoza wykazu
	// pozostaje nietknięty.
	ZapiszKonfiguracjeSesji(ctx context.Context,
		z shared.ConfigSessionSetRequest) (shared.ConfigSessionSetResponse, error)
	// KonfiguracjaObowiazujaca rozstrzyga obszary po ośmiu poziomach zasięgu
	// i trzech osiach wraz z pochodzeniem każdego obszaru.
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
				e.ustawienie(zmianaObszaru(w.StoredAreas, wpis), wpis)
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
