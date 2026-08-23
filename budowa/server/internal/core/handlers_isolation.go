// Odpowiedzialność pliku: wpięcie dwunastu komend rodziny `isolation.*` — okna
// konfiguracji punktów izolacji.
//
// Jeden port na całe okno. Dwanaście komend obsługuje jedno okno i jedną
// maszynerię: te same jedenaście punktów izolacji, ten sam adres zapisu i ten
// sam rozstrzygacz ośmiu poziomów. Trzy porty (macierz, profile, podgląd) byłyby
// trzema prawdami o jednym module.
//
// Rodzina ma trzy zdarzenia i żadne nie leci stąd. Rozgłasza je adapter, bo
// tylko on wie, co naprawdę poszło do bazy i pod jaki adres:
//
//   - `config.changed` — zapis punktu izolacji zmienia wiersz tabeli
//     `ustawienie`, więc idzie tym samym zdarzeniem, co zapis rodziny `config.*`;
//   - `isolation.profile.changed` — założenie, zmiana i skasowanie profilu
//     (`adapter_modul_isolation_rozgloszenie.go`);
//   - `isolation.policy.changed` — okno, którego polityka obowiązująca stała się
//     inna: po zapisie punktu pod adresem okna, po przypisaniu profilu do okna
//     i po przełączeniu warstwy okna.
//
// Zapis samego profilu polityki nie zmienia (profil jest szablonem), więc
// `isolation.profile.changed` i `isolation.policy.changed` nie chodzą parami.
// Uchwyty poniżej nie biorą nadajnika: brałyby go po to, żeby go nie użyć.
//
// Port niewypełniony nie rejestruje niczego: komendy odpowiedzą wtedy
// `isolation.unknown`, a pozostałe domeny pracują bez zmian.
package core

import (
	"context"

	"danacoconsole/shared"
)

// Izolacja jest portem rodziny `isolation.*`. Mówi wyłącznie typami
// kontraktu; wiersze tabel `ustawienie`, `profil_izolacji` i `warstwa_izolacji`
// leżą po drugiej stronie adaptera.
type Izolacja interface {
	// PoziomyZasiegu obsługuje `isolation.scope.list`.
	PoziomyZasiegu(ctx context.Context, z shared.IsolationScopeListRequest) (shared.IsolationScopeListResponse, error)
	// IzolacjaKontekstu obsługuje `isolation.context.get`.
	IzolacjaKontekstu(ctx context.Context, z shared.IsolationContextGetRequest) (shared.IsolationContextGetResponse, error)
	// ZapiszIzolacjeKontekstu obsługuje `isolation.context.set`.
	ZapiszIzolacjeKontekstu(ctx context.Context, z shared.IsolationContextSetRequest) (shared.IsolationContextSetResponse, error)
	// IzolacjaTechniczna obsługuje `isolation.technical.get`.
	IzolacjaTechniczna(ctx context.Context, z shared.IsolationTechnicalGetRequest) (shared.IsolationTechnicalGetResponse, error)
	// ZapiszIzolacjeTechniczna obsługuje `isolation.technical.set`.
	ZapiszIzolacjeTechniczna(ctx context.Context, z shared.IsolationTechnicalSetRequest) (shared.IsolationTechnicalSetResponse, error)
	// ZapiszProfilIzolacji obsługuje `isolation.profile.save`.
	ZapiszProfilIzolacji(ctx context.Context, z shared.IsolationProfileSaveRequest) (shared.IsolationProfileSaveResponse, error)
	// ProfileIzolacji obsługuje `isolation.profile.list`.
	ProfileIzolacji(ctx context.Context, z shared.IsolationProfileListRequest) (shared.IsolationProfileListResponse, error)
	// WczytajProfilIzolacji obsługuje `isolation.profile.load`.
	WczytajProfilIzolacji(ctx context.Context, z shared.IsolationProfileLoadRequest) (shared.IsolationProfileLoadResponse, error)
	// PrzypiszProfilIzolacji obsługuje `isolation.profile.assign`.
	PrzypiszProfilIzolacji(ctx context.Context, z shared.IsolationProfileAssignRequest) (shared.IsolationProfileAssignResponse, error)
	// UsunProfilIzolacji obsługuje `isolation.profile.delete`.
	UsunProfilIzolacji(ctx context.Context, z shared.IsolationProfileDeleteRequest) (shared.IsolationProfileDeleteResponse, error)
	// UstawWarstweIzolacji obsługuje `isolation.layer.set`.
	UstawWarstweIzolacji(ctx context.Context, z shared.IsolationLayerSetRequest) (shared.IsolationLayerSetResponse, error)
	// PodgladPolitykiIzolacji obsługuje `isolation.policy.preview`.
	PodgladPolitykiIzolacji(ctx context.Context, z shared.IsolationPolicyPreviewRequest) (shared.IsolationPolicyPreviewResponse, error)
}

// Adapter wypełnia port w całości. Gdyby port i adapter się rozjechały,
// kompilacja stanie tutaj, a nie dopiero na martwej komendzie u Operatora.
var _ Izolacja = (*adapterIzolacji)(nil)

// zarejestrujIzolacje wpina dwanaście komend rodziny `isolation.*`.
func zarejestrujIzolacje(r *Rejestr, i Izolacja) {
	if r == nil || i == nil {
		return
	}

	r.Zarejestruj(shared.CommandIsolationScopeList, obsluz(i.PoziomyZasiegu))
	r.Zarejestruj(shared.CommandIsolationContextGet, obsluz(i.IzolacjaKontekstu))
	r.Zarejestruj(shared.CommandIsolationContextSet, obsluz(i.ZapiszIzolacjeKontekstu))
	r.Zarejestruj(shared.CommandIsolationTechnicalGet, obsluz(i.IzolacjaTechniczna))
	r.Zarejestruj(shared.CommandIsolationTechnicalSet, obsluz(i.ZapiszIzolacjeTechniczna))
	r.Zarejestruj(shared.CommandIsolationProfileSave, obsluz(i.ZapiszProfilIzolacji))
	r.Zarejestruj(shared.CommandIsolationProfileList, obsluz(i.ProfileIzolacji))
	r.Zarejestruj(shared.CommandIsolationProfileLoad, obsluz(i.WczytajProfilIzolacji))
	r.Zarejestruj(shared.CommandIsolationProfileAssign, obsluz(i.PrzypiszProfilIzolacji))
	r.Zarejestruj(shared.CommandIsolationProfileDelete, obsluz(i.UsunProfilIzolacji))
	r.Zarejestruj(shared.CommandIsolationLayerSet, obsluz(i.UstawWarstweIzolacji))
	r.Zarejestruj(shared.CommandIsolationPolicyPreview, obsluz(i.PodgladPolitykiIzolacji))
}
