// Plik wpina dwanaście komend rodziny `isolation.*`, jeden port okna
// konfiguracji punktów izolacji, którego zdarzenia rozgłasza wyłącznie
// adapter.
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

// zarejestrujIzolacje wpina dwanaście komend rodziny `isolation.*`,
// obsługujących jedno okno konfiguracji punktów izolacji.
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
