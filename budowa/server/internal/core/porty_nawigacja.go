package core

import (
	"context"

	"danacoconsole/shared"
)

// Porty ścieżki wejścia do platformy i więzi klienta z kartą sesji.
//
// Plik dokłada dwa porty do kompletu z porty.go i niczego w nim nie zmienia.
// Reguła pozostaje ta sama: rdzeń mówi wyłącznie typami kontraktu,
// a kto je wypełni — repozytoria bazy, nadzorca sesji czy atrapa testu — jest
// wiedzą warstwy składania.

// Nawigacja obsługuje ścieżkę wejścia: strona główna → środowisko → moduł →
// przestrzeń robocza → stan okna.
//
// Byty tej ścieżki — środowisko, moduł, okno operacyjne, karta sesji — są
// wierszami słowników platformy, nie stałymi kodu. Nowe środowisko albo nowy
// moduł to nowy wiersz, nie nowa gałąź obsługiwacza.
type Nawigacja interface {
	StronaGlowna(ctx context.Context, z shared.HomeEnterRequest) (shared.HomeEnterResponse, error)
	Srodowiska(ctx context.Context, z shared.EnvironmentListRequest) (shared.EnvironmentListResponse, error)
	WejdzDoSrodowiska(ctx context.Context, z shared.EnvironmentEnterRequest) (shared.EnvironmentEnterResponse, error)
	Moduly(ctx context.Context, z shared.ModuleListRequest) (shared.ModuleListResponse, error)
	WejdzDoPrzestrzeni(ctx context.Context, z shared.WorkspaceEnterRequest) (shared.WorkspaceEnterResponse, error)
	StanOkna(ctx context.Context, z shared.WindowStateGetRequest) (shared.WindowStateGetResponse, error)
}

// WiazanieSesji obsługuje więź klienta z kartą sesji: ognisko oraz powiązanie
// połączenia z sesją trwającą na rdzeniu.
//
// Obie czynności są właściwością klienta, nie konta: dwa urządzenia tego samego
// konta mają własne ognisko i własne powiązanie, a rozłączenie któregokolwiek
// nie kończy sesji ani procesów.
type WiazanieSesji interface {
	Ogniskuj(ctx context.Context, z shared.SessionFocusRequest) (shared.SessionFocusResponse, error)
	Powiaz(ctx context.Context, z shared.SessionBindRequest) (shared.SessionBindResponse, error)
}
