// Rozgłoszenia rodziny `isolation.*`: dwa zdarzenia kontraktu,
// `isolation.profile.changed` i `isolation.policy.changed`, wraz z regułą,
// które wywołanie które z nich wysyła.
package core

import (
	"context"

	"danacoconsole/server/internal/konfig"
	"danacoconsole/shared"
)

// rozglosProfilIzolacji rozgłasza zmianę profilu izolacji, zdarzeniem
// `isolation.profile.changed` rodziny kontraktu.
func (a *adapterIzolacji) rozglosProfilIzolacji(ctx context.Context, zmiana shared.ChangeKind, kodProfilu string) {
	if a.rozgloszenie == nil || kodProfilu == "" {
		return
	}
	a.rozgloszenie.wyslijDoKonta(ctx, shared.EventIsolationProfileChanged, "",
		shared.IsolationProfileChangedEvent{Change: zmiana, ProfileId: kodProfilu})
}

// rozglosPolitykeAdresu rozgłasza zmianę polityki obowiązującej okno, gdy adres
// zapisu okna dotyczy. Adres spod innego poziomu nie rozgłasza niczego —
// powody w nagłówku pliku.
func (a *adapterIzolacji) rozglosPolitykeAdresu(ctx context.Context, adres konfig.Adres) {
	if adres.Poziom != shared.ConfigScopeWindow {
		return
	}
	a.rozglosPolitykeOkna(ctx, adres.KluczZasiegu)
}

// rozglosPolitykeOkna rozgłasza `isolation.policy.changed` dla jednego okna.
// Rodzaj zmiany jest zawsze `updated`: polityka okna nie powstaje ani nie znika,
// okno ma ją zawsze, choćby złożoną z samych wartości domyślnych.
func (a *adapterIzolacji) rozglosPolitykeOkna(ctx context.Context, idOkna string) {
	if a.rozgloszenie == nil || idOkna == "" {
		return
	}
	a.rozgloszenie.wyslijDoKonta(ctx, shared.EventIsolationPolicyChanged, "",
		shared.IsolationPolicyChangedEvent{Change: shared.ChangeKindUpdated, WindowId: idOkna})
}

// rozglosPolitykeWyboruWarstwy rozgłasza zmianę polityki po
// `isolation.layer.set`. Przełączenie warstwy zmienia politykę, choć nie
// zmienia wartości.
func (a *adapterIzolacji) rozglosPolitykeWyboruWarstwy(ctx context.Context, poziom shared.ConfigScope, byt string) {
	if poziom != shared.ConfigScopeWindow {
		return
	}
	a.rozglosPolitykeOkna(ctx, byt)
}

// rodzajZapisuProfilu rozstrzyga, czy `isolation.profile.save` założył profil,
// czy zmienił zastany. Rozstrzyga to samo żądanie: puste `profileId` znaczy
// profil nowy — tak mówi kontrakt i tak zachowuje się adapter.
func rodzajZapisuProfilu(idProfilu *string) shared.ChangeKind {
	if idProfilu == nil || *idProfilu == "" {
		return shared.ChangeKindCreated
	}
	return shared.ChangeKindUpdated
}
