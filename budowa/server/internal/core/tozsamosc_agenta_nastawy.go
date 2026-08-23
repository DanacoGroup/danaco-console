package core

import (
	"strings"

	"danacoconsole/server/internal/models"
)

// Nastawy procesu nakładane przez eksperta — wszystko poza promptem.
//
// Prompt eksperta jedzie warstwami; tutaj dochodzą pozostałe nastawy procesu.
//
// `injection/argumenty.go` produkuje dziewięć rodzajów przełącznika. Co z nich
// nakłada ekspert, a co nie — i dlaczego:
//
//	--model         Nakłada: `agent.model` i `agent.kanal_kod`.
//	--mcp-config    Nakłada: konektory eksperta rodzaju `mcp` dokładają się
//	                obok mostów okna i sesji, bo droga jest wieloelementowa
//	                z założenia.
//	--settings      Nakłada: kolumna `ustawienia_json` — zaczepy i reguły
//	                narzędzi.
//	--effort        Nie nakłada: brak nośnika w bazie.
//	--permission-mode  Nie nakłada. Tryb uprawnień per ekspert zamienia się
//	                w bramkę przy pierwszym nieuważnym użyciu, a jedyną bramką
//	                platformy jest logowanie. Gdyby wszedł, to wyłącznie jako
//	                nastawa procesu CLI, nigdy jako sprawdzenie w rdzeniu.
//	--add-dir       Nie nakłada: katalogi robocze są własnością okna i sesji,
//	                a ekspert nie jest miejscem pracy.
//	--resume        Nie nakłada: wznowienie jest tożsamością rozmowy.
//	--fallback-model   Nie nakłada: brak nośnika w bazie.

// nastawyAgenta nakłada na zapytanie kanału to, co ekspert wnosi poza promptem.
//
// Wywoływane po `uzupelnijKonfiguracje`, więc ekspert ma ostatnie słowo — tak
// samo jak w warstwach promptu. Zapytanie jest zmieniane w miejscu, bo to ten
// sam obiekt, który zaraz pojedzie do kanału; kopia byłaby drugą prawdą o turze.
func nastawyAgenta(zapytanie *models.Zapytanie, tozsamosc TozsamoscAgenta) {
	if zapytanie == nil || tozsamosc.Pusta() {
		return
	}
	nalozModel(zapytanie, tozsamosc)
	nalozUstawienia(zapytanie, tozsamosc)
	nalozMosty(zapytanie, tozsamosc)
}

// nalozModel podmienia model i kanał okna na wskazane przez eksperta.
//
// Podmiana, nie dopisanie — model jest jeden. Ekspert bez wskazanego modelu nie
// rusza tego, co ustawiło okno: pusta wartość znaczy „ekspert nie ma zdania",
// a nie „ekspert każe wrócić do domyślnego".
//
// Kanał idzie razem z modelem, bo `agent.model.set` ustawia obie wartości naraz
// i rozdzielenie ich dałoby model jednego dostawcy na drodze drugiego.
func nalozModel(zapytanie *models.Zapytanie, tozsamosc TozsamoscAgenta) {
	if strings.TrimSpace(tozsamosc.KanalKod) != "" {
		zapytanie.Kanal = tozsamosc.KanalKod
	}
	if strings.TrimSpace(tozsamosc.Model) != "" {
		zapytanie.Model = tozsamosc.Model
	}
}

// nalozUstawienia podmienia treść `--settings` na ustawienia eksperta.
//
// Podmiana, nie scalenie. Plik ustawień Claude Code ma kształt należący do CLI;
// scalanie dwóch takich plików po stronie rdzenia znaczyłoby, że rdzeń zna ich
// schemat i rozstrzyga konflikty zaczepów — czyli wersjonuje cudzy format.
// Ekspert, który wnosi własny harness, wnosi go w całości.
//
// Ekspert bez ustawień nie kasuje ustawień sesji: pusty napis znaczy „nie mam
// zdania".
func nalozUstawienia(zapytanie *models.Zapytanie, tozsamosc TozsamoscAgenta) {
	if strings.TrimSpace(tozsamosc.UstawieniaJSON) == "" {
		return
	}
	zapytanie.PlikUstawien = tozsamosc.UstawieniaJSON
}

// nalozMosty dokłada konfiguracje MCP konektorów eksperta.
//
// Dokłada, nie podmienia: droga `--mcp-config` jest wieloelementowa
// z założenia. Okno wnosi swoje mosty,
// sesja swoje, a każdy jedzie osobnym przełącznikiem. Mosty eksperta dopisują
// się na końcu, więc nadania okna nie znikają przez wybór eksperta.
func nalozMosty(zapytanie *models.Zapytanie, tozsamosc TozsamoscAgenta) {
	for _, konfiguracja := range tozsamosc.KonfiguracjeMCP {
		if strings.TrimSpace(konfiguracja) == "" {
			continue
		}
		zapytanie.DodatkoweMCP = append(zapytanie.DodatkoweMCP, konfiguracja)
	}
}
