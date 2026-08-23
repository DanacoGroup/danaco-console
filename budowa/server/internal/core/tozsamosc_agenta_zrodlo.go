package core

import (
	"context"

	"danacoconsole/server/internal/dane"
)

// Odczyt tożsamości eksperta z warstwy danych.
//
// Tożsamość eksperta leży w kilku tabelach naraz: `agent` (model, kanał,
// nastawy procesu), `agent_warstwa` (warstwy promptu) i `agent_konektor`
// (mosty MCP). Składacz
// wywołania nie ma prawa o tym wiedzieć, więc łączy je ten adapter.
//
// Każdy brak jest brakiem, nie błędem. Ekspert bez warstw, bez modelu, bez
// nastaw i bez mostów jest poprawnym ekspertem — po prostu nakłada mniej.
// Odmowa idzie wyłącznie wtedy, gdy eksperta o wskazanym kodzie nie ma
// w katalogu, a i ona nie zatrzymuje tury: wywołujący zamienia ją na tożsamość
// pustą.
//
// Pole `agent.instrukcje_systemowe` jest warstwą zerową: idzie pierwsze spośród
// tego, co ekspert wnosi, w warstwie najbardziej krytycznej, i jedzie do
// składacza tak samo jak warstwy.

var _ ZrodloTozsamosciAgenta = (*adapterTozsamosciAgenta)(nil)

type adapterTozsamosciAgenta struct {
	agenci  dane.RepozytoriumAgentow
	warstwy dane.RepozytoriumWarstwAgenta
	// mosty podaje konfiguracje MCP konektorów eksperta. Port opcjonalny —
	// bez niego ekspert nakłada wszystko poza mostami.
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

// ZMostami dokłada źródło konfiguracji MCP konektorów eksperta.
func (a *adapterTozsamosciAgenta) ZMostami(m zrodloMostowAgenta) *adapterTozsamosciAgenta {
	a.mosty = m
	return a
}

// TozsamoscAgenta składa migawkę tożsamości eksperta o wskazanym kodzie.
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
	// Warstwy i mosty są dokładkami: ich brak zostawia tożsamość uboższą,
	// a nie unieważnia jej. Ekspert bez warstw nadal nakłada model i nastawy.
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
