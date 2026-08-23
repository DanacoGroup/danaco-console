// Odpowiedzialność pliku: zadania projektu — hub planowania modułu Workspace
// w części „lista": założenie, zmiana, usunięcie i wykaz.
//
// Zadanie jest jednym bytem trzech widoków: pozycją listy, kartą tablicy
// i słupkiem osi czasu. Tablica i oś czasu mają własne pliki adaptera, lecz
// czytają ten sam wiersz — drugiego zapisu zadania w module nie ma.
//
// ── Zmiana jest łatą, nie podmianą ─────────────────────────────────────────
// `workspace.task.update` zmienia wyłącznie pola podane w żądaniu. Pole
// pominięte zostaje bez zmiany, bo okno wysyła jedno pole na jedną czynność
// Operatora (zmiana stanu, przypisanie wykonawcy, przesunięcie terminu), a
// podmiana całego zadania kasowałaby przy każdej z nich to, czego akurat nie
// było na ekranie.
package core

import (
	"context"
	"encoding/json"
	"errors"
	"sort"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// ZalozZadanie obsługuje `workspace.task.create`.
func (a *adapterPrzestrzeniRoboczej) ZalozZadanie(ctx context.Context,
	z shared.WorkspaceTaskCreateRequest) (shared.WorkspaceTaskCreateResponse, error) {

	projekt, err := a.projektDlaZapisu(ctx, z.ProjectId)
	if err != nil {
		return shared.WorkspaceTaskCreateResponse{}, err
	}
	if strings.TrimSpace(z.Title) == "" {
		return shared.WorkspaceTaskCreateResponse{}, bladProjektu("zadanie bez tytułu")
	}
	stan := shared.WorkspaceTaskStatus(shared.WorkspaceTaskStatusTodo)
	if z.Status != nil && *z.Status != "" {
		stan = *z.Status
	}
	waga := shared.WorkspaceTaskPriority(shared.WorkspaceTaskPriorityNormal)
	if z.Priority != nil && *z.Priority != "" {
		waga = *z.Priority
	}
	kolumna := ""
	if z.BoardColumnId != nil {
		kolumna = *z.BoardColumnId
	}
	if kolumna == "" {
		kolumna, err = a.kolumnaStanuWorkspace(ctx, projekt.ID, stan)
		if err != nil {
			return shared.WorkspaceTaskCreateResponse{}, err
		}
	}
	ostatnia, err := a.ostatniaRangaWorkspace(ctx, projekt.ID, kolumna)
	if err != nil {
		return shared.WorkspaceTaskCreateResponse{}, err
	}
	zadanie := dane.ZadanieWorkspace{
		ProjektID:            projekt.ID,
		Identyfikator:        nowyIdentyfikator("wstk-"),
		Tytul:                strings.TrimSpace(z.Title),
		Opis:                 wartoscTekstuWorkspace(z.Description),
		Stan:                 stan,
		Waga:                 waga,
		RodzajWykonawcy:      rodzajWykonawcyWorkspace(z.AssigneeKind, z.AssigneeId),
		Wykonawca:            wartoscTekstuWorkspace(z.AssigneeId),
		ZadanieNadrzedne:     wartoscTekstuWorkspace(z.ParentTaskId),
		KolumnaTablicy:       kolumna,
		KluczPorzadkowy:      rangaMiedzyWorkspace(ostatnia, ""),
		PoczatekMs:           wartoscChwiliWorkspace(z.StartAt),
		TerminMs:             wartoscChwiliWorkspace(z.DueAt),
		SzacunekMinut:        wartoscLiczbyWorkspace(z.EstimateMinutes),
		KamienMilowy:         z.Milestone != nil && *z.Milestone,
		RegulaPowtarzalnosci: wartoscTekstuWorkspace(z.RecurrenceRule),
		Etykiety:             z.Labels,
	}
	zapisane, err := a.repozytorium.ZapiszZadanieWorkspace(ctx, zadanie)
	if err != nil {
		return shared.WorkspaceTaskCreateResponse{}, err
	}
	if err := a.repozytorium.OdnotujCzynnosc(ctx, projekt.ID); err != nil {
		return shared.WorkspaceTaskCreateResponse{}, err
	}
	a.odnotujZdarzenieWorkspace(ctx, projekt.ID, shared.WorkspaceActivityKindTaskChanged,
		shared.ChangeKindCreated, shared.WorkspaceEntityKindTask, zapisane.Identyfikator,
		"założono zadanie „"+zapisane.Tytul+"”")

	return shared.WorkspaceTaskCreateResponse{Task: zadanieKontraktuWorkspace(zapisane)}, nil
}

// ZmienZadanie obsługuje `workspace.task.update`.
func (a *adapterPrzestrzeniRoboczej) ZmienZadanie(ctx context.Context,
	z shared.WorkspaceTaskUpdateRequest) (shared.WorkspaceTaskUpdateResponse, error) {

	zadanie, err := a.zadanieDoZmianyWorkspace(ctx, z.TaskId)
	if err != nil {
		return shared.WorkspaceTaskUpdateResponse{}, err
	}
	if z.Title != nil && strings.TrimSpace(*z.Title) != "" {
		zadanie.Tytul = strings.TrimSpace(*z.Title)
	}
	if z.Description != nil {
		zadanie.Opis = *z.Description
	}
	if z.Status != nil && *z.Status != "" {
		zadanie.Stan = *z.Status
		// Stan końcowy niesie czas ukończenia — pasek postępu projektu liczy
		// zamknięte zadania, a nie zadania z ustawioną plakietką.
		if *z.Status == shared.WorkspaceTaskStatusDone && zadanie.UkonczonoMs == 0 {
			zadanie.UkonczonoMs = terazWMilisekundachWorkspace()
			zadanie.PostepProcent = 100
		}
		if *z.Status != shared.WorkspaceTaskStatusDone {
			zadanie.UkonczonoMs = 0
		}
	}
	if z.Priority != nil && *z.Priority != "" {
		zadanie.Waga = *z.Priority
	}
	if z.AssigneeKind != nil || z.AssigneeId != nil {
		zadanie.RodzajWykonawcy = rodzajWykonawcyWorkspace(z.AssigneeKind, z.AssigneeId)
		if z.AssigneeId != nil {
			zadanie.Wykonawca = *z.AssigneeId
		}
	}
	if z.StartAt != nil {
		zadanie.PoczatekMs = *z.StartAt
	}
	if z.DueAt != nil {
		zadanie.TerminMs = *z.DueAt
	}
	if z.EstimateMinutes != nil {
		zadanie.SzacunekMinut = *z.EstimateMinutes
	}
	if z.SpentMinutes != nil {
		zadanie.SpedzonoMinut = *z.SpentMinutes
	}
	if z.ProgressPercent != nil {
		zadanie.PostepProcent = *z.ProgressPercent
	}
	if z.Milestone != nil {
		zadanie.KamienMilowy = *z.Milestone
	}
	if z.RecurrenceRule != nil {
		zadanie.RegulaPowtarzalnosci = *z.RecurrenceRule
	}
	if z.Labels != nil {
		zadanie.Etykiety = z.Labels
	}
	if z.Checklist != nil {
		zapis, err := json.Marshal(z.Checklist)
		if err != nil {
			return shared.WorkspaceTaskUpdateResponse{},
				bladProjektu("listy kontrolnej nie da się zapisać: " + err.Error())
		}
		zadanie.ListaKontrolnaJson = string(zapis)
	}
	zapisane, err := a.repozytorium.ZapiszZadanieWorkspace(ctx, zadanie)
	if err != nil {
		return shared.WorkspaceTaskUpdateResponse{}, err
	}
	przesuniete, err := a.przesunNastepnikiWorkspace(ctx, zapisane)
	if err != nil {
		return shared.WorkspaceTaskUpdateResponse{}, err
	}
	if err := a.repozytorium.OdnotujCzynnosc(ctx, zapisane.ProjektID); err != nil {
		return shared.WorkspaceTaskUpdateResponse{}, err
	}
	a.odnotujZdarzenieWorkspace(ctx, zapisane.ProjektID, shared.WorkspaceActivityKindTaskChanged,
		shared.ChangeKindUpdated, shared.WorkspaceEntityKindTask, zapisane.Identyfikator,
		"zmieniono zadanie „"+zapisane.Tytul+"”")

	return shared.WorkspaceTaskUpdateResponse{
		Task: zadanieKontraktuWorkspace(zapisane), RescheduledTaskIds: przesuniete,
	}, nil
}

// UsunZadanie obsługuje `workspace.task.delete`. Zadania podrzędne idą wraz
// z nadrzędnym: podzadanie bez zadania jest krokiem donikąd.
func (a *adapterPrzestrzeniRoboczej) UsunZadanie(ctx context.Context,
	z shared.WorkspaceTaskDeleteRequest) (shared.WorkspaceTaskDeleteResponse, error) {

	zadanie, err := a.repozytorium.ZadanieWorkspace(ctx, z.TaskId)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.WorkspaceTaskDeleteResponse{Deleted: false}, nil
	}
	if err != nil {
		return shared.WorkspaceTaskDeleteResponse{}, err
	}
	wszystkie, err := a.repozytorium.ZadaniaWorkspace(ctx, zadanie.ProjektID)
	if err != nil {
		return shared.WorkspaceTaskDeleteResponse{}, err
	}
	podrzedne := podrzedneZadaniaWorkspace(wszystkie, zadanie.Identyfikator)
	usuwane := append([]string{zadanie.Identyfikator}, podrzedne...)
	usuniete, err := a.repozytorium.UsunZadaniaWorkspace(ctx, usuwane)
	if err != nil {
		return shared.WorkspaceTaskDeleteResponse{}, err
	}
	if err := a.repozytorium.OdnotujCzynnosc(ctx, zadanie.ProjektID); err != nil {
		return shared.WorkspaceTaskDeleteResponse{}, err
	}
	a.odnotujZdarzenieWorkspace(ctx, zadanie.ProjektID, shared.WorkspaceActivityKindTaskChanged,
		shared.ChangeKindDeleted, shared.WorkspaceEntityKindTask, zadanie.Identyfikator,
		"usunięto zadanie „"+zadanie.Tytul+"”")

	return shared.WorkspaceTaskDeleteResponse{
		Deleted: usuniete > 0, DeletedSubtaskIds: podrzedne,
	}, nil
}

// Zadania obsługuje `workspace.task.list`.
func (a *adapterPrzestrzeniRoboczej) Zadania(ctx context.Context,
	z shared.WorkspaceTaskListRequest) (shared.WorkspaceTaskListResponse, error) {

	projekt, err := a.projektDlaZapisu(ctx, z.ProjectId)
	if err != nil {
		return shared.WorkspaceTaskListResponse{}, err
	}
	wszystkie, err := a.repozytorium.ZadaniaWorkspace(ctx, projekt.ID)
	if err != nil {
		return shared.WorkspaceTaskListResponse{}, err
	}
	fraza := strings.ToLower(strings.TrimSpace(wartoscTekstuWorkspace(z.Query)))
	dobrane := []shared.WorkspaceTask{}
	for _, zadanie := range wszystkie {
		if z.Status != nil && *z.Status != "" && zadanie.Stan != *z.Status {
			continue
		}
		if z.AssigneeId != nil && *z.AssigneeId != "" && zadanie.Wykonawca != *z.AssigneeId {
			continue
		}
		if z.Label != nil && *z.Label != "" && !zawieraTekstWorkspace(zadanie.Etykiety, *z.Label) {
			continue
		}
		if fraza != "" && !strings.Contains(strings.ToLower(zadanie.Tytul+" "+zadanie.Opis), fraza) {
			continue
		}
		if z.DueBefore != nil && (zadanie.TerminMs == 0 || zadanie.TerminMs >= *z.DueBefore) {
			continue
		}
		if !dolaczZamknieteWorkspace(z.IncludeDone) && stanZamknietyWorkspace(zadanie.Stan) {
			continue
		}
		dobrane = append(dobrane, zadanieKontraktuWorkspace(zadanie))
	}
	wszystkich := len(dobrane)
	return shared.WorkspaceTaskListResponse{
		Tasks: przytnijWykazWorkspace(dobrane, z.Limit, z.Offset), Total: wszystkich,
	}, nil
}

// zadanieDoZmianyWorkspace odnajduje zadanie wskazane żądaniem i nazywa brak
// zdaniem o żądaniu, a nie zdaniem o bazie.
func (a *adapterPrzestrzeniRoboczej) zadanieDoZmianyWorkspace(ctx context.Context,
	identyfikator string) (dane.ZadanieWorkspace, error) {

	if strings.TrimSpace(identyfikator) == "" {
		return dane.ZadanieWorkspace{}, bladProjektu("wywołanie bez wskazania zadania")
	}
	zadanie, err := a.repozytorium.ZadanieWorkspace(ctx, identyfikator)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return dane.ZadanieWorkspace{}, bladProjektu("zadania " + identyfikator + " nie ma w projekcie")
	}
	return zadanie, err
}

// przesunNastepnikiWorkspace przesuwa terminy zadań następujących po zadaniu
// właśnie zmienionym. Bez tego zależność „koniec–początek" byłaby ozdobą
// wykresu: przesunięcie poprzednika zostawiałoby następnik w przeszłości.
//
// Przejście idzie wszerz z licznikiem odwiedzin, żeby zależność zapętlona
// (gdyby powstała inną drogą niż komenda) nie zawiesiła zapisu.
func (a *adapterPrzestrzeniRoboczej) przesunNastepnikiWorkspace(ctx context.Context,
	zmienione dane.ZadanieWorkspace) ([]string, error) {

	zaleznosci, err := a.repozytorium.ZaleznosciWorkspace(ctx, zmienione.ProjektID)
	if err != nil {
		return nil, err
	}
	if len(zaleznosci) == 0 {
		return nil, nil
	}
	wszystkie, err := a.repozytorium.ZadaniaWorkspace(ctx, zmienione.ProjektID)
	if err != nil {
		return nil, err
	}
	wedlugKodu := map[string]dane.ZadanieWorkspace{}
	for _, zadanie := range wszystkie {
		wedlugKodu[zadanie.Identyfikator] = zadanie
	}
	wedlugKodu[zmienione.Identyfikator] = zmienione

	przesuniete := []string{}
	odwiedzone := map[string]bool{}
	kolejka := []string{zmienione.Identyfikator}
	for len(kolejka) > 0 {
		biezace := kolejka[0]
		kolejka = kolejka[1:]
		if odwiedzone[biezace] {
			continue
		}
		odwiedzone[biezace] = true
		poprzednik := wedlugKodu[biezace]
		if poprzednik.TerminMs == 0 {
			continue
		}
		for _, zaleznosc := range zaleznosci {
			if zaleznosc.Poprzednik != biezace ||
				zaleznosc.Rodzaj != shared.WorkspaceDependencyKindFinishToStart {
				continue
			}
			nastepnik, jest := wedlugKodu[zaleznosc.Nastepnik]
			if !jest || nastepnik.PoczatekMs == 0 {
				continue
			}
			najwczesniejszy := poprzednik.TerminMs + int64(zaleznosc.OdstepMinut)*60_000
			if nastepnik.PoczatekMs >= najwczesniejszy {
				continue
			}
			roznica := najwczesniejszy - nastepnik.PoczatekMs
			nastepnik.PoczatekMs += roznica
			if nastepnik.TerminMs > 0 {
				nastepnik.TerminMs += roznica
			}
			zapisany, err := a.repozytorium.ZapiszZadanieWorkspace(ctx, nastepnik)
			if err != nil {
				return nil, err
			}
			wedlugKodu[zapisany.Identyfikator] = zapisany
			przesuniete = append(przesuniete, zapisany.Identyfikator)
			kolejka = append(kolejka, zapisany.Identyfikator)
		}
	}
	return przesuniete, nil
}

// podrzedneZadaniaWorkspace zbiera całe poddrzewo zadań podrzędnych.
func podrzedneZadaniaWorkspace(wszystkie []dane.ZadanieWorkspace, korzen string) []string {
	podrzedne := []string{}
	kolejka := []string{korzen}
	for len(kolejka) > 0 {
		biezace := kolejka[0]
		kolejka = kolejka[1:]
		for _, zadanie := range wszystkie {
			if zadanie.ZadanieNadrzedne != biezace || zadanie.Identyfikator == korzen {
				continue
			}
			podrzedne = append(podrzedne, zadanie.Identyfikator)
			kolejka = append(kolejka, zadanie.Identyfikator)
		}
	}
	sort.Strings(podrzedne)
	return podrzedne
}

// zadanieKontraktuWorkspace przekłada wiersz zadania na byt kontraktu.
func zadanieKontraktuWorkspace(z dane.ZadanieWorkspace) shared.WorkspaceTask {
	waga := z.Waga
	kamien := z.KamienMilowy
	zadanie := shared.WorkspaceTask{
		Id: z.Identyfikator, ProjectId: z.ProjektKod, Title: z.Tytul,
		Status: z.Stan, Priority: &waga, Milestone: &kamien,
		Labels:    z.Etykiety,
		CreatedAt: chwilaBazy(z.Utworzono), UpdatedAt: chwilaBazy(z.Zaktualizowano),
	}
	if z.Opis != "" {
		opis := z.Opis
		zadanie.Description = &opis
	}
	if z.RodzajWykonawcy != "" {
		rodzaj := shared.WorkspaceAssigneeKind(z.RodzajWykonawcy)
		zadanie.AssigneeKind = &rodzaj
	}
	if z.Wykonawca != "" {
		wykonawca := z.Wykonawca
		zadanie.AssigneeId = &wykonawca
	}
	if z.ZadanieNadrzedne != "" {
		nadrzedne := z.ZadanieNadrzedne
		zadanie.ParentTaskId = &nadrzedne
	}
	if z.KolumnaTablicy != "" {
		kolumna := z.KolumnaTablicy
		zadanie.BoardColumnId = &kolumna
	}
	if z.KluczPorzadkowy != "" {
		ranga := z.KluczPorzadkowy
		zadanie.Rank = &ranga
	}
	zadanie.StartAt = chwilaKontraktuWorkspace(z.PoczatekMs)
	zadanie.DueAt = chwilaKontraktuWorkspace(z.TerminMs)
	zadanie.CompletedAt = chwilaKontraktuWorkspace(z.UkonczonoMs)
	if z.SzacunekMinut > 0 {
		szacunek := z.SzacunekMinut
		zadanie.EstimateMinutes = &szacunek
	}
	if z.SpedzonoMinut > 0 {
		spedzono := z.SpedzonoMinut
		zadanie.SpentMinutes = &spedzono
	}
	postep := z.PostepProcent
	zadanie.ProgressPercent = &postep
	if z.RegulaPowtarzalnosci != "" {
		regula := z.RegulaPowtarzalnosci
		zadanie.RecurrenceRule = &regula
	}
	if z.ListaKontrolnaJson != "" {
		var lista []shared.WorkspaceTaskChecklistItem
		if err := json.Unmarshal([]byte(z.ListaKontrolnaJson), &lista); err == nil {
			zadanie.Checklist = lista
		}
	}
	return zadanie
}

// rodzajWykonawcyWorkspace ustala rodzaj wykonawcy: wskazanie żądania, a przy
// jego braku — rodzaj wynikający z obecności eksperta.
func rodzajWykonawcyWorkspace(rodzaj *shared.WorkspaceAssigneeKind, ekspert *string) string {
	if rodzaj != nil && *rodzaj != "" {
		return string(*rodzaj)
	}
	if ekspert != nil && *ekspert != "" {
		return shared.WorkspaceAssigneeKindAgent
	}
	return ""
}

// stanZamknietyWorkspace odróżnia stany końcowe od stanów pracy.
func stanZamknietyWorkspace(stan shared.WorkspaceTaskStatus) bool {
	return stan == shared.WorkspaceTaskStatusDone || stan == shared.WorkspaceTaskStatusCancelled
}

// dolaczZamknieteWorkspace czyta pole `includeDone`; brak znaczy „nie”.
func dolaczZamknieteWorkspace(pole *bool) bool {
	return pole != nil && *pole
}

// zawieraTekstWorkspace szuka wartości w wykazie bez oglądania się na wielkość
// liter — etykieta wpisana ręcznie rzadko trafia w wielkość liter poprzedniej.
func zawieraTekstWorkspace(wykaz []string, szukana string) bool {
	for _, pozycja := range wykaz {
		if strings.EqualFold(pozycja, szukana) {
			return true
		}
	}
	return false
}

// wartoscTekstuWorkspace rozwija pole opcjonalne kontraktu do napisu.
func wartoscTekstuWorkspace(pole *string) string {
	if pole == nil {
		return ""
	}
	return *pole
}

// wartoscLiczbyWorkspace rozwija pole opcjonalne kontraktu do liczby.
func wartoscLiczbyWorkspace(pole *int) int {
	if pole == nil {
		return 0
	}
	return *pole
}

// wartoscChwiliWorkspace rozwija chwilę opcjonalną; zero znaczy brak granicy.
func wartoscChwiliWorkspace(pole *int64) int64 {
	if pole == nil {
		return 0
	}
	return *pole
}

// chwilaKontraktuWorkspace zwija zero kolumny do braku pola kontraktu. Zero
// wypuszczone jako chwila znaczyłoby rok 1970, a nie „bez granicy".
func chwilaKontraktuWorkspace(wartosc int64) *int64 {
	if wartosc == 0 {
		return nil
	}
	chwila := wartosc
	return &chwila
}
