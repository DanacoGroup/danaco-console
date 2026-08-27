package core

import (
	"context"

	"danacoconsole/shared"
)

// Plik dokłada porty wejścia do platformy i więzi klienta z kartą sesji, uzupełniając porty.go.

// Nawigacja obsługuje ścieżkę wejścia: strona główna, środowisko, moduł, przestrzeń robocza i stan okna, których byty są wierszami słowników platformy.
type Nawigacja interface {
	StronaGlowna(ctx context.Context, z shared.HomeEnterRequest) (shared.HomeEnterResponse, error)
	Srodowiska(ctx context.Context, z shared.EnvironmentListRequest) (shared.EnvironmentListResponse, error)
	WejdzDoSrodowiska(ctx context.Context, z shared.EnvironmentEnterRequest) (shared.EnvironmentEnterResponse, error)
	Moduly(ctx context.Context, z shared.ModuleListRequest) (shared.ModuleListResponse, error)
	WejdzDoPrzestrzeni(ctx context.Context, z shared.WorkspaceEnterRequest) (shared.WorkspaceEnterResponse, error)
	StanOkna(ctx context.Context, z shared.WindowStateGetRequest) (shared.WindowStateGetResponse, error)
}

// WiazanieSesji obsługuje więź klienta z kartą sesji: ognisko oraz powiązanie połączenia z sesją, jako właściwość klienta, a nie konta.
type WiazanieSesji interface {
	Ogniskuj(ctx context.Context, z shared.SessionFocusRequest) (shared.SessionFocusResponse, error)
	Powiaz(ctx context.Context, z shared.SessionBindRequest) (shared.SessionBindResponse, error)
}
