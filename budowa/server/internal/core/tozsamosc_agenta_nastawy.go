package core

import (
	"strings"

	"danacoconsole/server/internal/models"
)

// Nastawy procesu nakładane przez eksperta poza promptem: model, kanał, ustawienia i mosty MCP.

// nastawyAgenta nakłada na zapytanie kanału to, co ekspert wnosi poza promptem. Wywoływane po uzupelnijKonfiguracje, więc ekspert ma ostatnie słowo. Zapytanie jest zmieniane w miejscu, bo to ten sam obiekt, który pojedzie do kanału.
func nastawyAgenta(zapytanie *models.Zapytanie, tozsamosc TozsamoscAgenta) {
	if zapytanie == nil || tozsamosc.Pusta() {
		return
	}
	nalozModel(zapytanie, tozsamosc)
	nalozUstawienia(zapytanie, tozsamosc)
	nalozMosty(zapytanie, tozsamosc)
}

// nalozModel podmienia model i kanał okna na wskazane przez eksperta. Podmiana, nie dopisanie — model jest jeden. Ekspert bez wskazanego modelu nie rusza tego, co ustawiło okno. Kanał idzie z modelem, bo agent.model.set ustawia je razem.
func nalozModel(zapytanie *models.Zapytanie, tozsamosc TozsamoscAgenta) {
	if strings.TrimSpace(tozsamosc.KanalKod) != "" {
		zapytanie.Kanal = tozsamosc.KanalKod
	}
	if strings.TrimSpace(tozsamosc.Model) != "" {
		zapytanie.Model = tozsamosc.Model
	}
}

// nalozUstawienia podmienia treść --settings na ustawienia eksperta. Podmiana, nie scalenie: plik ustawień CLI ma własny kształt, a ekspert, który wnosi własny harness, wnosi go w całości. Ekspert bez ustawień nie kasuje ustawień sesji.
func nalozUstawienia(zapytanie *models.Zapytanie, tozsamosc TozsamoscAgenta) {
	if strings.TrimSpace(tozsamosc.UstawieniaJSON) == "" {
		return
	}
	zapytanie.PlikUstawien = tozsamosc.UstawieniaJSON
}

// nalozMosty dokłada konfiguracje MCP konektorów eksperta. Dokłada, nie podmienia: droga --mcp-config jest wieloelementowa z założenia, więc mosty eksperta dopisują się na końcu, a nadania okna nie znikają przez wybór eksperta.
func nalozMosty(zapytanie *models.Zapytanie, tozsamosc TozsamoscAgenta) {
	for _, konfiguracja := range tozsamosc.KonfiguracjeMCP {
		if strings.TrimSpace(konfiguracja) == "" {
			continue
		}
		zapytanie.DodatkoweMCP = append(zapytanie.DodatkoweMCP, konfiguracja)
	}
}
