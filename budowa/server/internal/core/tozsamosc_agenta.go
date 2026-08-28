package core

import (
	"context"

	"danacoconsole/server/internal/dane"
)

// Tożsamość eksperta nałożona na okno: czwarte źródło treści, nie czwarta oś kontraktu.

// TozsamoscAgenta niesie wszystko, czym ekspert nakłada się na wywołanie: warstwy promptu, model bazowy i nastawy procesu. Struktura jest płaską migawką powstającą raz na turę. Ekspert nierozpoznany daje strukturę pustą, nie błąd.
type TozsamoscAgenta struct {
	// Kod eksperta; pusty znaczy „okno pracuje na modelu surowym".
	Kod string
	// ImieWlasne pokazuje Operator; rdzeń go nie używa do niczego poza podglądem.
	ImieWlasne string
	// Warstwy promptu eksperta wraz z ich trybem.
	Warstwy []dane.WarstwaAgenta
	// InstrukcjeSystemowe to warstwa zerowa eksperta, pierwsza w warstwie najbardziej krytycznej.
	InstrukcjeSystemowe string
	// Model i KanalKod nakładają `--model`; przyjmuje je `agent.model.set`.
	Model    string
	KanalKod string
	// UstawieniaJSON nakłada `--settings`: zaczepy i reguły narzędzi.
	UstawieniaJSON string
	// KonfiguracjeMCP nakładają kolejne --mcp-config z konektorów eksperta, obok konfiguracji okna.
	KonfiguracjeMCP []string
}

// Pusta mówi, czy tożsamość niesie cokolwiek do nałożenia na wywołanie, czyli czy kod eksperta jest pusty.
func (t TozsamoscAgenta) Pusta() bool {
	return t.Kod == ""
}

// ZrodloTozsamosciAgenta podaje tożsamość eksperta po jego kodzie.
//
// Osobny port, a nie sięgnięcie do repozytoriów wprost: składacz wywołania ma
// jednego dostawcę tożsamości eksperta i nie zna tego, że siedzi ona w czterech
// tabelach.
type ZrodloTozsamosciAgenta interface {
	TozsamoscAgenta(ctx context.Context, kod string) (TozsamoscAgenta, error)
}

// tozsamoscAgentaOkna zwraca tożsamość eksperta wskazanego przez okno. Trzy przypadki dają strukturę pustą bez błędu: okno bez eksperta, brak wpiętego źródła, ekspert nierozpoznany — żaden z nich nie zatrzymuje tury.
func tozsamoscAgentaOkna(ctx context.Context, zrodlo ZrodloTozsamosciAgenta, kod string) TozsamoscAgenta {
	if zrodlo == nil || kod == "" {
		return TozsamoscAgenta{}
	}
	tozsamosc, err := zrodlo.TozsamoscAgenta(ctx, kod)
	if err != nil {
		return TozsamoscAgenta{}
	}
	return tozsamosc
}
