// Odpowiedzialność pliku: graf wiedzy projektu (`workspace.knowledge.graph.get`)
// oraz tablica wizualna (`workspace.canvas.get`, `workspace.canvas.save`).
// Rodzina `workspace.knowledge` nie miesza się z rodziną `knowledge.*`.
package core

import (
	"context"
	"encoding/json"
	"sort"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// granicaWezlowGrafuWorkspace jest domyślną granicą wielkości grafu wiedzy
// zwracanego bez jawnie podanego limitu w żądaniu.
const granicaWezlowGrafuWorkspace = 500

// GrafWiedzy obsługuje `workspace.knowledge.graph.get` i rysuje sieć wprost
// napisanych odnośników między bytami projektu.
func (a *adapterPrzestrzeniRoboczej) GrafWiedzy(ctx context.Context,
	z shared.WorkspaceKnowledgeGraphGetRequest) (shared.WorkspaceKnowledgeGraphGetResponse, error) {

	projekt, err := a.projektDlaZapisu(ctx, z.ProjectId)
	if err != nil {
		return shared.WorkspaceKnowledgeGraphGetResponse{}, err
	}
	wezly, krawedzie, err := a.siecProjektuWorkspace(ctx, projekt.ID, projekt.Kod)
	if err != nil {
		return shared.WorkspaceKnowledgeGraphGetResponse{}, err
	}
	if len(z.Kinds) > 0 {
		wezly, krawedzie = zawezRodzajeGrafuWorkspace(wezly, krawedzie, z.Kinds)
	}
	if z.RootId != nil && *z.RootId != "" {
		glebokosc := 1
		if z.Depth != nil && *z.Depth > 0 {
			glebokosc = *z.Depth
		}
		wezly, krawedzie = rozwinOdKorzeniaWorkspace(wezly, krawedzie, *z.RootId, glebokosc)
	}
	granica := granicaWezlowGrafuWorkspace
	if z.MaxNodes != nil && *z.MaxNodes > 0 {
		granica = *z.MaxNodes
	}
	przyciety := false
	if len(wezly) > granica {
		wezly, przyciety = wezly[:granica], true
		krawedzie = krawedzieWezlowWorkspace(wezly, krawedzie)
	}
	stopnie := map[string]int{}
	for _, krawedz := range krawedzie {
		stopnie[krawedz.SourceId]++
		stopnie[krawedz.TargetId]++
	}
	for i := range wezly {
		stopien := stopnie[wezly[i].Id]
		wezly[i].Degree = &stopien
	}
	return shared.WorkspaceKnowledgeGraphGetResponse{Graph: shared.WorkspaceKnowledgeGraph{
		ProjectId: projekt.Kod, Nodes: wezly, Edges: krawedzie, Truncated: &przyciety,
	}}, nil
}

// Kanwa obsługuje `workspace.canvas.get` i oddaje treść tablicy wizualnej
// wskazanego projektu wraz z wykazem jej identyfikatorów.
func (a *adapterPrzestrzeniRoboczej) Kanwa(ctx context.Context,
	z shared.WorkspaceCanvasGetRequest) (shared.WorkspaceCanvasGetResponse, error) {

	projekt, err := a.projektDlaZapisu(ctx, z.ProjectId)
	if err != nil {
		return shared.WorkspaceCanvasGetResponse{}, err
	}
	tablice, err := a.repozytorium.TabliceWorkspace(ctx, projekt.ID)
	if err != nil {
		return shared.WorkspaceCanvasGetResponse{}, err
	}
	kody := make([]string, 0, len(tablice))
	for _, tablica := range tablice {
		kody = append(kody, tablica.Identyfikator)
	}
	odpowiedz := shared.WorkspaceCanvasGetResponse{CanvasIds: kody}
	if len(tablice) == 0 {
		// Projekt bez tablicy oddaje brak pola, nie odmowę: „jeszcze nie ma”
		// jest stanem początkowym płótna.
		return odpowiedz, nil
	}
	wybrana := tablice[0]
	if z.CanvasId != nil && *z.CanvasId != "" {
		znaleziona := false
		for _, tablica := range tablice {
			if tablica.Identyfikator == *z.CanvasId {
				wybrana, znaleziona = tablica, true
			}
		}
		if !znaleziona {
			return shared.WorkspaceCanvasGetResponse{},
				bladProjektu("tablicy " + *z.CanvasId + " nie ma w projekcie")
		}
	}
	kanwa := kanwaKontraktuWorkspace(wybrana)
	odpowiedz.Canvas = &kanwa
	return odpowiedz, nil
}

// ZapiszKanwe obsługuje `workspace.canvas.save` i zapisuje treść tablicy
// wizualnej pod jej identyfikatorem w projekcie.
func (a *adapterPrzestrzeniRoboczej) ZapiszKanwe(ctx context.Context,
	z shared.WorkspaceCanvasSaveRequest) (shared.WorkspaceCanvasSaveResponse, error) {

	projekt, err := a.projektDlaZapisu(ctx, z.ProjectId)
	if err != nil {
		return shared.WorkspaceCanvasSaveResponse{}, err
	}
	scena := strings.TrimSpace(string(z.Scene))
	if scena == "" {
		scena = "{}"
	}
	if !json.Valid([]byte(scena)) {
		return shared.WorkspaceCanvasSaveResponse{},
			bladProjektu("scena tablicy nie jest poprawnym zapisem JSON")
	}
	identyfikator := wartoscTekstuWorkspace(z.CanvasId)
	nazwa := pierwszaNiepustaWorkspace(wartoscTekstuWorkspace(z.Name), "Tablica projektu")
	if identyfikator == "" {
		identyfikator = nowyIdentyfikator("wscv-")
	} else if z.Name == nil {
		// Zapis samej sceny nie ma prawa przemianować tablicy.
		istniejaca, err := a.repozytorium.TablicaWorkspace(ctx, identyfikator)
		if err == nil {
			nazwa = istniejaca.Nazwa
		}
	}
	zapisana, err := a.repozytorium.ZapiszTabliceWorkspace(ctx, dane.TablicaWorkspace{
		ProjektID: projekt.ID, Identyfikator: identyfikator, Nazwa: nazwa, Scena: scena,
	})
	if err != nil {
		return shared.WorkspaceCanvasSaveResponse{}, err
	}
	if err := a.repozytorium.OdnotujCzynnosc(ctx, projekt.ID); err != nil {
		return shared.WorkspaceCanvasSaveResponse{}, err
	}
	a.odnotujZdarzenieWorkspace(ctx, projekt.ID, shared.WorkspaceActivityKindNoteChanged,
		shared.ChangeKindUpdated, shared.WorkspaceEntityKindNote, zapisana.Identyfikator,
		"zapisano tablicę wizualną „"+zapisana.Nazwa+"”")

	return shared.WorkspaceCanvasSaveResponse{Canvas: kanwaKontraktuWorkspace(zapisana)}, nil
}

// siecProjektuWorkspace składa węzły i krawędzie z bytów projektu, tworząc
// graf odnośników zwracany wołającemu.
func (a *adapterPrzestrzeniRoboczej) siecProjektuWorkspace(ctx context.Context, projektID int64,
	kodProjektu string) ([]shared.WorkspaceKnowledgeNode, []shared.WorkspaceKnowledgeEdge, error) {

	notatki, err := a.repozytorium.NotatkiWorkspace(ctx, projektID)
	if err != nil {
		return nil, nil, err
	}
	zadania, err := a.repozytorium.ZadaniaWorkspace(ctx, projektID)
	if err != nil {
		return nil, nil, err
	}
	odnosniki, err := a.repozytorium.OdnosnikiWorkspace(ctx, projektID)
	if err != nil {
		return nil, nil, err
	}
	komentarze, err := a.repozytorium.KomentarzeWorkspace(ctx, projektID)
	if err != nil {
		return nil, nil, err
	}
	wpisy, err := a.repozytorium.WpisyPamieci(ctx, projektID, 0)
	if err != nil {
		return nil, nil, err
	}

	wezly := []shared.WorkspaceKnowledgeNode{}
	rodzaje := map[string]shared.WorkspaceEntityKind{}
	dolozWezel := func(kod string, rodzaj shared.WorkspaceEntityKind, nazwa string) {
		if kod == "" {
			return
		}
		if _, jest := rodzaje[kod]; jest {
			return
		}
		rodzaje[kod] = rodzaj
		wezly = append(wezly, shared.WorkspaceKnowledgeNode{Id: kod, Kind: rodzaj, Label: nazwa})
	}
	for _, notatka := range notatki {
		dolozWezel(notatka.Identyfikator, shared.WorkspaceEntityKindNote, notatka.Tytul)
	}
	for _, zadanie := range zadania {
		dolozWezel(zadanie.Identyfikator, shared.WorkspaceEntityKindTask, zadanie.Tytul)
	}
	for _, wpis := range wpisy {
		dolozWezel(wpis.Identyfikator, shared.WorkspaceEntityKindMemoryEntry,
			skrocOpisWorkspace(wpis.Tresc, 60))
	}
	for _, plik := range a.plikiProjektu(kodProjektu, "") {
		dolozWezel(plik.Id, shared.WorkspaceEntityKindFile, plik.Name)
	}
	for _, komentarz := range komentarze {
		dolozWezel(komentarz.Identyfikator, shared.WorkspaceEntityKindComment,
			skrocOpisWorkspace(komentarz.Tresc, 60))
	}

	krawedzie := []shared.WorkspaceKnowledgeEdge{}
	dolozKrawedz := func(zrodlo, cel string, rodzaj shared.WorkspaceGraphEdgeKind) {
		rodzajZrodla, mamZrodlo := rodzaje[zrodlo]
		rodzajCelu, mamCel := rodzaje[cel]
		if !mamZrodlo || !mamCel {
			return
		}
		krawedzie = append(krawedzie, shared.WorkspaceKnowledgeEdge{
			SourceId: zrodlo, SourceKind: rodzajZrodla,
			TargetId: cel, TargetKind: rodzajCelu, Kind: rodzaj,
		})
	}
	for _, odnosnik := range odnosniki {
		dolozKrawedz(odnosnik.NotatkaZrodlowa, odnosnik.NotatkaDocelowa, odnosnik.Rodzaj)
	}
	for _, zadanie := range zadania {
		dolozKrawedz(zadanie.Identyfikator, zadanie.ZadanieNadrzedne,
			shared.WorkspaceGraphEdgeKindReference)
	}
	for _, komentarz := range komentarze {
		dolozKrawedz(komentarz.Identyfikator, komentarz.Byt, shared.WorkspaceGraphEdgeKindReference)
	}
	sort.SliceStable(wezly, func(i, j int) bool { return wezly[i].Label < wezly[j].Label })
	return wezly, krawedzie, nil
}

// zawezRodzajeGrafuWorkspace zostawia węzły wskazanych rodzajów wraz
// z krawędziami, których oba końce zostały.
func zawezRodzajeGrafuWorkspace(wezly []shared.WorkspaceKnowledgeNode,
	krawedzie []shared.WorkspaceKnowledgeEdge,
	rodzaje []shared.WorkspaceEntityKind) ([]shared.WorkspaceKnowledgeNode, []shared.WorkspaceKnowledgeEdge) {

	dozwolone := map[shared.WorkspaceEntityKind]bool{}
	for _, rodzaj := range rodzaje {
		dozwolone[rodzaj] = true
	}
	zostawione := []shared.WorkspaceKnowledgeNode{}
	for _, wezel := range wezly {
		if dozwolone[wezel.Kind] {
			zostawione = append(zostawione, wezel)
		}
	}
	return zostawione, krawedzieWezlowWorkspace(zostawione, krawedzie)
}

// rozwinOdKorzeniaWorkspace zostawia węzły osiągalne z bytu początkowego
// w zadanej liczbie kroków — sąsiedztwo jest nieskierowane, bo panel ogląda
// otoczenie bytu, a nie jego następstwa.
func rozwinOdKorzeniaWorkspace(wezly []shared.WorkspaceKnowledgeNode,
	krawedzie []shared.WorkspaceKnowledgeEdge, korzen string,
	glebokosc int) ([]shared.WorkspaceKnowledgeNode, []shared.WorkspaceKnowledgeEdge) {

	sasiedzi := map[string][]string{}
	for _, krawedz := range krawedzie {
		sasiedzi[krawedz.SourceId] = append(sasiedzi[krawedz.SourceId], krawedz.TargetId)
		sasiedzi[krawedz.TargetId] = append(sasiedzi[krawedz.TargetId], krawedz.SourceId)
	}
	osiagalne := map[string]bool{korzen: true}
	warstwa := []string{korzen}
	for krok := 0; krok < glebokosc; krok++ {
		nastepna := []string{}
		for _, kod := range warstwa {
			for _, sasiad := range sasiedzi[kod] {
				if osiagalne[sasiad] {
					continue
				}
				osiagalne[sasiad] = true
				nastepna = append(nastepna, sasiad)
			}
		}
		warstwa = nastepna
	}
	zostawione := []shared.WorkspaceKnowledgeNode{}
	for _, wezel := range wezly {
		if osiagalne[wezel.Id] {
			zostawione = append(zostawione, wezel)
		}
	}
	return zostawione, krawedzieWezlowWorkspace(zostawione, krawedzie)
}

// krawedzieWezlowWorkspace zostawia krawędzie, których oba końce są w wykazie
// węzłów — krawędź w próżnię nie jest powiązaniem, tylko śladem po przycięciu.
func krawedzieWezlowWorkspace(wezly []shared.WorkspaceKnowledgeNode,
	krawedzie []shared.WorkspaceKnowledgeEdge) []shared.WorkspaceKnowledgeEdge {

	obecne := map[string]bool{}
	for _, wezel := range wezly {
		obecne[wezel.Id] = true
	}
	zostawione := []shared.WorkspaceKnowledgeEdge{}
	for _, krawedz := range krawedzie {
		if obecne[krawedz.SourceId] && obecne[krawedz.TargetId] {
			zostawione = append(zostawione, krawedz)
		}
	}
	return zostawione
}

// kanwaKontraktuWorkspace przekłada wiersz tablicy wizualnej z bazy danych
// na byt kontraktu zwracany wołającemu.
func kanwaKontraktuWorkspace(t dane.TablicaWorkspace) shared.WorkspaceCanvas {
	return shared.WorkspaceCanvas{
		Id: t.Identyfikator, ProjectId: t.ProjektKod, Name: t.Nazwa,
		Scene:     json.RawMessage(t.Scena),
		CreatedAt: chwilaBazy(t.Utworzono), UpdatedAt: chwilaBazy(t.Zaktualizowano),
	}
}

// skrocOpisWorkspace przycina treść do nazwy węzła. Graf pokazuje nazwy, a wpis
// pamięci bywa akapitem — nieprzycięty zamalowałby rysunek.
func skrocOpisWorkspace(tresc string, granica int) string {
	tresc = strings.TrimSpace(strings.ReplaceAll(tresc, "\n", " "))
	znaki := []rune(tresc)
	if len(znaki) <= granica {
		return tresc
	}
	return string(znaki[:granica]) + "…"
}
