// Odpowiedzialność pliku: okno Orchestrator — ustalenie zależności między
// krokami, sprawdzenie układu i wskazanie ścieżki krytycznej.
//
// Sprawdzenie układu nie odmawia zapisu: układ zapisuje się także wtedy, gdy ma
// cykl, a odpowiedź niesie `valid=false` wraz z zastrzeżeniami. Odmowa zapisu
// kasowałaby pracę wykonaną do chwili wykrycia usterki, zamiast pokazać, co
// wymaga poprawki.
package core

import (
	"context"
	"sort"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// Zaleznosci zapisuje układ zależności automatyki i oddaje jego ocenę.
// Żądanie bez pola `dependencies` jest samym sprawdzeniem układu zastanego —
// tak działa przycisk „Waliduj graf” panelu akcji.
func (a *adapterAutomatyk) Zaleznosci(ctx context.Context,
	z shared.AutomationOrchestratorDefineRequest) (shared.AutomationOrchestratorDefineResponse, error) {

	wiersz, err := a.wiersz(ctx, z.WorkflowId)
	if err != nil {
		return shared.AutomationOrchestratorDefineResponse{}, err
	}
	if z.Dependencies != nil {
		if err := a.repozytorium.ZapiszZaleznosci(ctx, wiersz.ID, wierszeZaleznosci(z.Dependencies)); err != nil {
			return shared.AutomationOrchestratorDefineResponse{}, bladAutomatyki(err)
		}
	}
	kroki, err := a.repozytorium.Kroki(ctx, wiersz.ID)
	if err != nil {
		return shared.AutomationOrchestratorDefineResponse{}, bladAutomatyki(err)
	}
	zapisane, err := a.repozytorium.Zaleznosci(ctx, wiersz.ID)
	if err != nil {
		return shared.AutomationOrchestratorDefineResponse{}, bladAutomatyki(err)
	}
	zastrzezenia := sprawdzUklad(kroki, zapisane)
	return shared.AutomationOrchestratorDefineResponse{
		Dependencies:        zaleznosciKontraktu(zapisane),
		Valid:               len(zastrzezenia) == 0,
		Issues:              zastrzezenia,
		CriticalPathStepIds: sciezkaKrytyczna(kroki, zapisane),
	}, nil
}

// wierszeZaleznosci przekłada zależności kontraktu na wiersze. Łuk bez rodzaju
// jest sekwencyjny — to znaczenie domyślne „krok po kroku”.
func wierszeZaleznosci(zaleznosci []shared.AutomationDependency) []dane.ZaleznoscKroku {
	wiersze := make([]dane.ZaleznoscKroku, 0, len(zaleznosci))
	for _, zaleznosc := range zaleznosci {
		if zaleznosc.FromStepId == "" || zaleznosc.ToStepId == "" ||
			zaleznosc.FromStepId == zaleznosc.ToStepId {
			// Łuk pusty i pętla własna nie przechodzą przez więzy schematu;
			// zastrzeżenie zgłosi walidacja układu, a zapis nie ma się wywrócić.
			continue
		}
		rodzaj := string(zaleznosc.Kind)
		if rodzaj == "" {
			rodzaj = shared.AutomationDependencyKindSequential
		}
		wiersze = append(wiersze, dane.ZaleznoscKroku{
			KrokZ: zaleznosc.FromStepId, KrokDo: zaleznosc.ToStepId,
			Rodzaj: rodzaj, Warunek: zaleznosc.Condition,
		})
	}
	return wiersze
}

// zaleznosciKontraktu przekłada wiersze układu na zależności kontraktu.
func zaleznosciKontraktu(wiersze []dane.ZaleznoscKroku) []shared.AutomationDependency {
	zaleznosci := make([]shared.AutomationDependency, 0, len(wiersze))
	for _, wiersz := range wiersze {
		zaleznosci = append(zaleznosci, shared.AutomationDependency{
			FromStepId: wiersz.KrokZ, ToStepId: wiersz.KrokDo,
			Kind: shared.AutomationDependencyKind(wiersz.Rodzaj), Condition: wiersz.Warunek,
		})
	}
	return zaleznosci
}

// sprawdzUklad wylicza zastrzeżenia do układu: łuk do kroku nieistniejącego,
// łuk warunkowy bez warunku i cykl. Wykaz pusty znaczy układ poprawny.
func sprawdzUklad(kroki []dane.KrokAutomatyki, zaleznosci []dane.ZaleznoscKroku) []string {
	istniejace := make(map[string]bool, len(kroki))
	for _, krok := range kroki {
		istniejace[krok.Kod] = true
	}
	zastrzezenia := []string{}
	for _, zaleznosc := range zaleznosci {
		if !istniejace[zaleznosc.KrokZ] {
			zastrzezenia = append(zastrzezenia,
				"zależność wychodzi z kroku, którego nie ma w automatyce: "+zaleznosc.KrokZ)
		}
		if !istniejace[zaleznosc.KrokDo] {
			zastrzezenia = append(zastrzezenia,
				"zależność prowadzi do kroku, którego nie ma w automatyce: "+zaleznosc.KrokDo)
		}
		if zaleznosc.Rodzaj == shared.AutomationDependencyKindConditional &&
			(zaleznosc.Warunek == nil || *zaleznosc.Warunek == "") {
			zastrzezenia = append(zastrzezenia,
				"zależność warunkowa bez warunku: "+zaleznosc.KrokZ+" → "+zaleznosc.KrokDo)
		}
	}
	if cykl := znajdzCykl(kroki, zaleznosci); cykl != "" {
		zastrzezenia = append(zastrzezenia, "układ zawiera cykl: "+cykl)
	}
	return zastrzezenia
}

// znajdzCykl orzeka, czy układ da się ułożyć w kolejność wykonania. Zwraca
// nazwę pierwszego kroku, który do kolejności nie wszedł — poprawka ma od czego
// zacząć, zamiast samego „graf jest zły”.
func znajdzCykl(kroki []dane.KrokAutomatyki, zaleznosci []dane.ZaleznoscKroku) string {
	_, pozostale := porzadekTopologiczny(kroki, zaleznosci)
	if len(pozostale) == 0 {
		return ""
	}
	sort.Strings(pozostale)
	return pozostale[0]
}

// porzadekTopologiczny układa kroki w kolejność wykonania metodą Kahna. Drugi
// wynik to kroki, które do kolejności nie weszły — czyli te uwikłane w cykl.
func porzadekTopologiczny(kroki []dane.KrokAutomatyki,
	zaleznosci []dane.ZaleznoscKroku) ([]string, []string) {

	stopien := map[string]int{}
	nastepnicy := map[string][]string{}
	for _, krok := range kroki {
		stopien[krok.Kod] = 0
	}
	for _, zaleznosc := range zaleznosci {
		if _, jest := stopien[zaleznosc.KrokDo]; !jest {
			continue
		}
		if _, jest := stopien[zaleznosc.KrokZ]; !jest {
			continue
		}
		stopien[zaleznosc.KrokDo]++
		nastepnicy[zaleznosc.KrokZ] = append(nastepnicy[zaleznosc.KrokZ], zaleznosc.KrokDo)
	}
	gotowe := []string{}
	for _, krok := range kroki {
		if stopien[krok.Kod] == 0 {
			gotowe = append(gotowe, krok.Kod)
		}
	}
	kolejnosc := []string{}
	for len(gotowe) > 0 {
		biezacy := gotowe[0]
		gotowe = gotowe[1:]
		kolejnosc = append(kolejnosc, biezacy)
		for _, nastepny := range nastepnicy[biezacy] {
			stopien[nastepny]--
			if stopien[nastepny] == 0 {
				gotowe = append(gotowe, nastepny)
			}
		}
	}
	wykonane := make(map[string]bool, len(kolejnosc))
	for _, kod := range kolejnosc {
		wykonane[kod] = true
	}
	pozostale := []string{}
	for _, krok := range kroki {
		if !wykonane[krok.Kod] {
			pozostale = append(pozostale, krok.Kod)
		}
	}
	return kolejnosc, pozostale
}

// sciezkaKrytyczna wskazuje najdłuższy łańcuch kroków układu — pozycję
// „podświetlenie ścieżki krytycznej” panelu akcji Orchestratora. Układ z cyklem
// nie ma najdłuższej ścieżki, więc wskazanie jest wtedy puste; okno pokazuje
// wówczas zastrzeżenia, a nie zmyśloną kolejność.
func sciezkaKrytyczna(kroki []dane.KrokAutomatyki, zaleznosci []dane.ZaleznoscKroku) []string {
	kolejnosc, pozostale := porzadekTopologiczny(kroki, zaleznosci)
	if len(pozostale) > 0 || len(kolejnosc) == 0 {
		return nil
	}
	poprzednicy := map[string][]string{}
	for _, zaleznosc := range zaleznosci {
		poprzednicy[zaleznosc.KrokDo] = append(poprzednicy[zaleznosc.KrokDo], zaleznosc.KrokZ)
	}
	dlugosc := map[string]int{}
	wiodacy := map[string]string{}
	najdluzszy := kolejnosc[0]
	for _, kod := range kolejnosc {
		dlugosc[kod] = 1
		for _, poprzednik := range poprzednicy[kod] {
			if dlugosc[poprzednik]+1 > dlugosc[kod] {
				dlugosc[kod] = dlugosc[poprzednik] + 1
				wiodacy[kod] = poprzednik
			}
		}
		if dlugosc[kod] > dlugosc[najdluzszy] {
			najdluzszy = kod
		}
	}
	sciezka := []string{}
	for kod := najdluzszy; kod != ""; kod = wiodacy[kod] {
		sciezka = append([]string{kod}, sciezka...)
	}
	return sciezka
}
