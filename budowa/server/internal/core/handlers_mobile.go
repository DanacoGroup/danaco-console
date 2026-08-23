// Odpowiedzialność pliku: wpięcie trzech komend rodziny `mobile.*` — mobilnego
// centrum dowodzenia.
//
// Port jest rozszerzeniem portu nawigacji, nie drugim portem. Warstwa mobilna
// pyta o stan warstwy wspólnej (sesje, okna, procesy, kolejki), a tę obsługuje
// Nawigacja; rejestr telemetrii postępu ma w rdzeniu jednego właściciela.
// Wzorem jest `monitor.*`, które tą samą drogą rozszerza ten sam port, oraz
// `memory.*`, które rozszerza port modułu Workspace.
//
// Brak portu nie jest ciszą — tak samo jak w monitorze. Obsługiwacze rejestrują
// się zawsze, a port, który nie niesie warstwy mobilnej, odpowiada
// `internal_error` z nazwą brakującego bytu. Operator ma zobaczyć „warstwa
// mobilna nie jest wpięta", a nie „nie znam takiej komendy" ani — najgorsze —
// pusty wykaz procesów i wyzerowany stan platformy.
//
// Ta rodzina nie rozgłasza zdarzeń, mimo że kontrakt je zna. Kontrakt niesie
// `mobile.process.changed`, ale proces mobilny jest procesem telemetrii
// postępu, a telemetria ma jednego producenta zdarzeń (`telemetria.go`,
// `progress.changed`). Rozgłaszanie z tego pliku dałoby zdarzenie wyłącznie po
// zmianie zleconej z telefonu, a milczałoby przy każdej zmianie zleconej
// z pulpitu — czyli drugą prawdę o procesie, i to niepełną. Dlatego
// `zarejestrujWarstweMobilna` nie bierze emitera.
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

// zarejestrujWarstweMobilna wpina trzy komendy rodziny `mobile.*`.
//
// Asercja jest miękka, a nie twarda, dokładnie z tego powodu, co w monitorze:
// odmowa ma dojść do Operatora, a rdzeń, który padłby przy montażu, powiedziałby
// ją wyłącznie temu, kto czyta dziennik startu.
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

// bladWarstwyMobilnejNiewpietej składa odmowę portu, który nie niesie warstwy
// mobilnej.
func bladWarstwyMobilnejNiewpietej() error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"warstwa mobilna: port nawigacji nie niesie mobilnego centrum dowodzenia — "+
			"stanu platformy ani procesów nie ma skąd odczytać"))
}
