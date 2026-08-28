package core

import (
	"context"

	"danacoconsole/server/internal/dane"
)

// Odczyt tożsamości eksperta z warstwy danych: agent, warstwy promptu i mosty MCP w jednym adapterze.

var _ ZrodloTozsamosciAgenta = (*adapterTozsamosciAgenta)(nil)

type adapterTozsamosciAgenta struct {
	agenci  dane.RepozytoriumAgentow
	warstwy dane.RepozytoriumWarstwAgenta
	// mosty podaje konfiguracje MCP konektorów eksperta. Port opcjonalny — bez niego brak mostów.
	mosty zrodloMostowAgenta
}

// zrodloMostowAgenta podaje konfiguracje MCP konektorów eksperta.
//
// Osobny port, bo konektory eksperta wiążą się z punktami dostępu, a te
// prowadzi inne repozytorium. Wpięcie opcjonalne: brak mostów nie unieważnia
// reszty tożsamości.
type zrodloMostowAgenta interface {
	MostyAgenta(ctx context.Context, kodAgenta string) ([]string, error)
}

func nowyAdapterTozsamosciAgenta(agenci dane.RepozytoriumAgentow,
	warstwy dane.RepozytoriumWarstwAgenta) *adapterTozsamosciAgenta {

	return &adapterTozsamosciAgenta{agenci: agenci, warstwy: warstwy}
}

// ZMostami dokłada źródło konfiguracji MCP konektorów eksperta jako port opcjonalny wpinany do adaptera.
func (a *adapterTozsamosciAgenta) ZMostami(m zrodloMostowAgenta) *adapterTozsamosciAgenta {
	a.mosty = m
	return a
}

// TozsamoscAgenta składa migawkę tożsamości eksperta o wskazanym kodzie, łącząc dane z trzech tabel bazy.
func (a *adapterTozsamosciAgenta) TozsamoscAgenta(ctx context.Context,
	kod string) (TozsamoscAgenta, error) {

	if a == nil || a.agenci == nil || kod == "" {
		return TozsamoscAgenta{}, nil
	}
	agent, err := a.agenci.PoKodzie(ctx, kod)
	if err != nil {
		return TozsamoscAgenta{}, err
	}
	tozsamosc := TozsamoscAgenta{
		Kod:                 agent.Kod,
		ImieWlasne:          agent.ImieWlasne,
		Model:               wartoscTekstu(agent.Model),
		KanalKod:            wartoscTekstu(agent.KanalKod),
		UstawieniaJSON:      agent.UstawieniaJSON,
		InstrukcjeSystemowe: agent.InstrukcjeSystemowe,
	}
	// Warstwy i mosty są dokładkami: ich brak zostawia tożsamość uboższą, nie unieważnia jej.
	if a.warstwy != nil {
		if warstwy, err := a.warstwy.Warstwy(ctx, kod); err == nil {
			tozsamosc.Warstwy = warstwy
		}
	}
	if a.mosty != nil {
		if mosty, err := a.mosty.MostyAgenta(ctx, kod); err == nil {
			tozsamosc.KonfiguracjeMCP = mosty
		}
	}
	return tozsamosc, nil
}
