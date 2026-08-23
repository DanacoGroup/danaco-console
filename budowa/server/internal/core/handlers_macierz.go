package core

import "context"

// Port macierzy widoczności modułów; komend nie rejestruje.
//
// Kontrakt nie definiuje ani jednej komendy macierzy — nie ma `matrix.*` ani
// `module.visibility.*` — więc nie ma czego wpiąć. Macierz dociera do klienta
// wyłącznie jako pola bytów nawigacji: `Environment.moduleCodes`,
// `Module.environmentCodes` i `Environment.navigationKind`, obsługiwane przez
// `handlers_nawigacja.go`. Rejestracja nazwy spoza kontraktu byłaby ogłoszeniem
// zdolności, której kontrakt nie opisuje.
//
// Port istnieje, bo macierz ma w rdzeniu czytelnika — nawigację — a ta bierze
// ją stąd zamiast sięgać po repozytorium wprost. Odwzorowanie „moduł →
// środowiska, w których jest widoczny" ma dzięki temu jedno miejsce, a gdy
// komendy macierzy powstaną, będzie już co zarejestrować.

// Macierz udostępnia macierz widoczności modułów w środowiskach.
type Macierz interface {
	// KodySrodowisk zwraca odwzorowanie identyfikatora modułu na kody
	// środowisk, w których moduł jest widoczny, w kolejności kart środowisk.
	// Moduł nieobecny w wyniku nie ma okna modułowego w żadnym środowisku.
	KodySrodowisk(ctx context.Context) (map[int64][]string, error)
}

// Asercja rozjazdu portu z adapterem na etapie kompilacji.
var _ Macierz = (*adapterMacierzy)(nil)
