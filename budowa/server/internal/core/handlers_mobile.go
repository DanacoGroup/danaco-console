// Plik wpina trzy komendy rodziny mobile.* mobilnego centrum dowodzenia,
// jako rozszerzenie portu nawigacji obsługujące stan warstwy wspólnej
// i sterowanie procesami.
package core

import (
	"context"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// mobilnosc jest portem rodziny `mobile.*` — portem nawigacji
// rozszerzonym o trzy komendy mobilnego centrum dowodzenia.
type mobilnosc interface {
	Nawigacja

	// StanMobilny obsługuje `mobile.status.get`.
	StanMobilny(ctx context.Context, z shared.MobileStatusGetRequest) (shared.MobileStatusGetResponse, error)
	// ProcesyMobilne obsługuje `mobile.process.list`.
	ProcesyMobilne(ctx context.Context, z shared.MobileProcessListRequest) (shared.MobileProcessListResponse, error)
	// SterujProcesemMobilnym obsługuje `mobile.process.control`.
	SterujProcesemMobilnym(ctx context.Context, z shared.MobileProcessControlRequest) (shared.MobileProcessControlResponse, error)
}

// Adapter wypełnia port w całości. Gdyby port i adapter się rozjechały,
// kompilacja stanie tutaj, a nie dopiero na martwej komendzie u Operatora.
var _ mobilnosc = (*adapterMobilny)(nil)

// zarejestrujWarstweMobilna wpina trzy komendy rodziny mobile.*, rejestrując
// je zawsze; port bez warstwy mobilnej odpowiada odmową zamiast milczeć.
func zarejestrujWarstweMobilna(r *Rejestr, n Nawigacja) {
	if r == nil {
		return
	}
	m, niesie := n.(mobilnosc)
	if !niesie {
		zarejestrujWarstweMobilnaNiewpieta(r)
		return
	}
	r.Zarejestruj(shared.CommandMobileStatusGet, obsluz(m.StanMobilny))
	r.Zarejestruj(shared.CommandMobileProcessList, obsluz(m.ProcesyMobilne))
	r.Zarejestruj(shared.CommandMobileProcessControl, obsluz(m.SterujProcesemMobilnym))
}

// zarejestrujWarstweMobilnaNiewpieta wpina odmowę na miejsce obsługiwaczy.
// Komenda zostaje znana rdzeniowi i odpowiada wprost, czego brakuje.
func zarejestrujWarstweMobilnaNiewpieta(r *Rejestr) {
	r.Zarejestruj(shared.CommandMobileStatusGet,
		obsluz(func(context.Context, shared.MobileStatusGetRequest) (shared.MobileStatusGetResponse, error) {
			return shared.MobileStatusGetResponse{}, bladWarstwyMobilnejNiewpietej()
		}))
	r.Zarejestruj(shared.CommandMobileProcessList,
		obsluz(func(context.Context, shared.MobileProcessListRequest) (shared.MobileProcessListResponse, error) {
			return shared.MobileProcessListResponse{}, bladWarstwyMobilnejNiewpietej()
		}))
	r.Zarejestruj(shared.CommandMobileProcessControl,
		obsluz(func(context.Context, shared.MobileProcessControlRequest) (shared.MobileProcessControlResponse, error) {
			return shared.MobileProcessControlResponse{}, bladWarstwyMobilnejNiewpietej()
		}))
}

// bladWarstwyMobilnejNiewpietej składa odmowę portu nawigacji, który nie
// niesie mobilnego centrum dowodzenia, zamiast zwracać pustą odpowiedź.
func bladWarstwyMobilnejNiewpietej() error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"warstwa mobilna: port nawigacji nie niesie mobilnego centrum dowodzenia — "+
			"stanu platformy ani procesów nie ma skąd odczytać"))
}
