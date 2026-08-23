package core

import (
	"context"

	"danacoconsole/server/internal/dane"
)

// Tożsamość eksperta nałożona na okno: czwarte źródło treści, nie czwarta oś.
//
// `ConfigAxis` kontraktu jest zbiorem zamkniętym trzech wartości: `platform`,
// `model`, `account`. Oś odpowiada na pytanie, dla czego wartość obowiązuje
// przy rozstrzyganiu konfiguracji, i stoi prostopadle do poziomu zasięgu.
// Ekspert nie jest ani zasięgiem, ani przedmiotem rozstrzygania — jest
// tożsamością nałożoną na kanał konkretnego okna, tak jak mówi opis pola
// w kontrakcie: „agent nie zastepuje kanalu — kanal jest droga do modelu,
// agent tozsamoscia nalozona na te droge".
//
// Ekspert wchodzi więc jako czwarte źródło treści do tych samych trzech warstw
// (konstytucja · profil roli · ekspertyza), którymi jadą osie, i przez ten sam
// składacz. Jako źródło nie wymaga zmiany kontraktu ani rozstrzygacza
// konfiguracji, a skutek jest ten sam: treść eksperta dojeżdża do procesu.
//
// Ekspert dopisuje i nigdy nie zastępuje. Prompt systemowy platformy jest
// pierwszą paczką, wbudowaną i zawsze obowiązującą; ekspert dokłada do niej
// warstwę instrukcji kontraktowej jako zakres własny Operatora. Treść eksperta
// idzie ostatnia w każdej warstwie i niczego nie kasuje, a przełącznik wiersza
// wywołania jest przy nim zawsze dopisujący.
//
// Nośniki trybu zastępowania (`agent_warstwa.tryb`, `Agent.mode`) zostają
// w bazie i w kontrakcie, ale składacz ich nie czyta.

// TozsamoscAgenta niesie wszystko, czym ekspert nakłada się na wywołanie:
// warstwy promptu, model bazowy i nastawy procesu.
//
// Struktura jest płaską migawką, nie uchwytem do bazy — powstaje raz na turę
// i jedzie do składacza wywołania. Ekspert nierozpoznany daje strukturę pustą,
// nie błąd: kod, którego nikt nie zna, jest brakiem, a nie przeszkodą
// w rozmowie.
type TozsamoscAgenta struct {
	// Kod eksperta; pusty znaczy „okno pracuje na modelu surowym".
	Kod string
	// ImieWlasne pokazuje Operator; rdzeń go nie używa do niczego poza podglądem.
	ImieWlasne string
	// Warstwy promptu eksperta wraz z ich trybem.
	Warstwy []dane.WarstwaAgenta
	// InstrukcjeSystemowe to warstwa zerowa eksperta (kolumna
	// `agent.instrukcje_systemowe`).
	//
	// Zerowa znaczy: pierwsza spośród tego, co ekspert wnosi, w warstwie
	// najbardziej krytycznej. Nie jest czwartą warstwą obok trzech, bo nie ma
	// własnego pola w `models.Nakladka` i mieć go nie musi — jest wstępem do
	// konstytucji eksperta, nie osobnym poziomem krytyczności.
	InstrukcjeSystemowe string
	// Model i KanalKod nakładają `--model`; przyjmuje je `agent.model.set`.
	Model    string
	KanalKod string
	// UstawieniaJSON nakłada `--settings`: zaczepy i reguły narzędzi.
	UstawieniaJSON string
	// KonfiguracjeMCP nakładają kolejne `--mcp-config` z konektorów eksperta.
	// Dokładają się obok konfiguracji okna i sesji, nie zamiast niej.
	KonfiguracjeMCP []string
}

// Pusta mówi, czy jest cokolwiek do nałożenia.
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

// tozsamoscAgentaOkna zwraca tożsamość eksperta wskazanego przez okno.
//
// Trzy przypadki dają to samo — strukturę pustą, bez błędu i bez zatrzymania
// tury: okno bez eksperta, brak wpiętego źródła, ekspert nierozpoznany. Turę
// zatrzymuje wyłącznie to, co ją naprawdę uniemożliwia, a brak nałożenia nią
// nie jest.
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
