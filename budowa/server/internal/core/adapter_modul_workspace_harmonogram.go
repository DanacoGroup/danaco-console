// Odpowiedzialność pliku: oś czasu projektu i zależności między zadaniami —
// `workspace.schedule.get`, `workspace.task.dependency.set`
// i `workspace.task.dependency.remove`. Zależność domykająca cykl czasowy
// nie zostaje zapisana.
package core

import (
	"context"
	"errors"
	"sort"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// Harmonogram obsługuje `workspace.schedule.get` i buduje słupki osi czasu
// z zależnościami między zadaniami projektu.
func (a *adapterPrzestrzeniRoboczej) Harmonogram(ctx context.Context,
	z shared.WorkspaceScheduleGetRequest) (shared.WorkspaceScheduleGetResponse, error) {

	projekt, err := a.projektDlaZapisu(ctx, z.ProjectId)
	if err != nil {
		return shared.WorkspaceScheduleGetResponse{}, err
	}
	zadania, err := a.repozytorium.ZadaniaWorkspace(ctx, projekt.ID)
	if err != nil {
		return shared.WorkspaceScheduleGetResponse{}, err
	}
	zaleznosci, err := a.repozytorium.ZaleznosciWorkspace(ctx, projekt.ID)
	if err != nil {
		return shared.WorkspaceScheduleGetResponse{}, err
	}
	krytyczne := sciezkaKrytycznaWorkspace(zadania, zaleznosci)

	slupki := []shared.WorkspaceScheduleBar{}
	poza := []string{}
	for _, zadanie := range zadania {
		poczatek, koniec := granicaSlupkaWorkspace(zadanie)
		if poczatek == 0 || koniec == 0 {
			poza = append(poza, zadanie.Identyfikator)
			continue
		}
		if z.From != nil && koniec < *z.From {
			continue
		}
		if z.To != nil && poczatek > *z.To {
			continue
		}
		postep := zadanie.PostepProcent
		kamien := zadanie.KamienMilowy
		naSciezce := krytyczne[zadanie.Identyfikator]
		slupek := shared.WorkspaceScheduleBar{
			TaskId: zadanie.Identyfikator, Title: zadanie.Tytul,
			StartAt: poczatek, EndAt: koniec,
			ProgressPercent: &postep, Milestone: &kamien, Critical: &naSciezce,
		}
		if zadanie.Wykonawca != "" {
			wykonawca := zadanie.Wykonawca
			slupek.AssigneeId = &wykonawca
		}
		slupki = append(slupki, slupek)
	}
	sort.SliceStable(slupki, func(i, j int) bool { return slupki[i].StartAt < slupki[j].StartAt })

	wykazZaleznosci := make([]shared.WorkspaceTaskDependency, 0, len(zaleznosci))
	for _, zaleznosc := range zaleznosci {
		wykazZaleznosci = append(wykazZaleznosci, zaleznoscKontraktuWorkspace(projekt.Kod, zaleznosc))
	}
	return shared.WorkspaceScheduleGetResponse{
		Bars: slupki, Dependencies: wykazZaleznosci, UnscheduledTaskIds: poza,
	}, nil
}

// ZalozZaleznosc obsługuje `workspace.task.dependency.set` i odmawia zapisu,
// gdyby nowa krawędź domykała cykl.
func (a *adapterPrzestrzeniRoboczej) ZalozZaleznosc(ctx context.Context,
	z shared.WorkspaceTaskDependencySetRequest) (shared.WorkspaceTaskDependencySetResponse, error) {

	projekt, err := a.projektDlaZapisu(ctx, z.ProjectId)
	if err != nil {
		return shared.WorkspaceTaskDependencySetResponse{}, err
	}
	if z.PredecessorTaskId == "" || z.SuccessorTaskId == "" {
		return shared.WorkspaceTaskDependencySetResponse{},
			bladProjektu("zależność bez wskazania obu zadań")
	}
	if z.PredecessorTaskId == z.SuccessorTaskId {
		return shared.WorkspaceTaskDependencySetResponse{},
			bladProjektu("zadanie nie może zależeć od samego siebie")
	}
	for _, kod := range []string{z.PredecessorTaskId, z.SuccessorTaskId} {
		zadanie, err := a.repozytorium.ZadanieWorkspace(ctx, kod)
		if errors.Is(err, dane.ErrBrakWiersza) {
			return shared.WorkspaceTaskDependencySetResponse{},
				bladProjektu("zadania " + kod + " nie ma w projekcie")
		}
		if err != nil {
			return shared.WorkspaceTaskDependencySetResponse{}, err
		}
		if zadanie.ProjektID != projekt.ID {
			return shared.WorkspaceTaskDependencySetResponse{},
				bladProjektu("zadanie " + kod + " należy do innego projektu")
		}
	}
	zaleznosci, err := a.repozytorium.ZaleznosciWorkspace(ctx, projekt.ID)
	if err != nil {
		return shared.WorkspaceTaskDependencySetResponse{}, err
	}
	if domykaCyklWorkspace(zaleznosci, z.PredecessorTaskId, z.SuccessorTaskId) {
		return shared.WorkspaceTaskDependencySetResponse{},
			bladProjektu("zależność " + z.PredecessorTaskId + " → " + z.SuccessorTaskId +
				" domyka cykl — takiego porządku nie da się ułożyć w czasie")
	}
	rodzaj := shared.WorkspaceDependencyKind(shared.WorkspaceDependencyKindFinishToStart)
	if z.Kind != nil && *z.Kind != "" {
		rodzaj = *z.Kind
	}
	zapisana, err := a.repozytorium.ZapiszZaleznoscWorkspace(ctx, dane.ZaleznoscWorkspace{
		ProjektID: projekt.ID, Identyfikator: nowyIdentyfikator("wsdp-"),
		Poprzednik: z.PredecessorTaskId, Nastepnik: z.SuccessorTaskId,
		Rodzaj: rodzaj, OdstepMinut: wartoscLiczbyWorkspace(z.LagMinutes),
	})
	if err != nil {
		return shared.WorkspaceTaskDependencySetResponse{}, err
	}
	poprzednik, err := a.repozytorium.ZadanieWorkspace(ctx, z.PredecessorTaskId)
	if err != nil {
		return shared.WorkspaceTaskDependencySetResponse{}, err
	}
	przesuniete, err := a.przesunNastepnikiWorkspace(ctx, poprzednik)
	if err != nil {
		return shared.WorkspaceTaskDependencySetResponse{}, err
	}
	if err := a.repozytorium.OdnotujCzynnosc(ctx, projekt.ID); err != nil {
		return shared.WorkspaceTaskDependencySetResponse{}, err
	}
	a.odnotujZdarzenieWorkspace(ctx, projekt.ID, shared.WorkspaceActivityKindTaskChanged,
		shared.ChangeKindCreated, shared.WorkspaceEntityKindTask, z.SuccessorTaskId,
		"założono zależność "+z.PredecessorTaskId+" → "+z.SuccessorTaskId)

	return shared.WorkspaceTaskDependencySetResponse{
		Dependency:         zaleznoscKontraktuWorkspace(projekt.Kod, zapisana),
		RescheduledTaskIds: przesuniete,
	}, nil
}

// ZniesZaleznosc obsługuje `workspace.task.dependency.remove` i usuwa
// krawędź zależności między dwoma zadaniami.
func (a *adapterPrzestrzeniRoboczej) ZniesZaleznosc(ctx context.Context,
	z shared.WorkspaceTaskDependencyRemoveRequest) (shared.WorkspaceTaskDependencyRemoveResponse, error) {

	if strings.TrimSpace(z.DependencyId) == "" {
		return shared.WorkspaceTaskDependencyRemoveResponse{},
			bladProjektu("zniesienie bez wskazania zależności")
	}
	zaleznosc, err := a.repozytorium.ZaleznoscWorkspace(ctx, z.DependencyId)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.WorkspaceTaskDependencyRemoveResponse{Removed: false}, nil
	}
	if err != nil {
		return shared.WorkspaceTaskDependencyRemoveResponse{}, err
	}
	zniesiona, err := a.repozytorium.UsunZaleznoscWorkspace(ctx, z.DependencyId)
	if err != nil {
		return shared.WorkspaceTaskDependencyRemoveResponse{}, err
	}
	if zniesiona {
		a.odnotujZdarzenieWorkspace(ctx, zaleznosc.ProjektID, shared.WorkspaceActivityKindTaskChanged,
			shared.ChangeKindDeleted, shared.WorkspaceEntityKindTask, zaleznosc.Nastepnik,
			"zniesiono zależność "+zaleznosc.Poprzednik+" → "+zaleznosc.Nastepnik)
	}
	return shared.WorkspaceTaskDependencyRemoveResponse{Removed: zniesiona}, nil
}

// granicaSlupkaWorkspace ustala początek i koniec słupka. Kamień milowy jest
// chwilą, nie odcinkiem, więc jedna granica wystarcza mu za obie.
func granicaSlupkaWorkspace(z dane.ZadanieWorkspace) (int64, int64) {
	poczatek, koniec := z.PoczatekMs, z.TerminMs
	if z.KamienMilowy {
		chwila := pierwszaChwilaWorkspace(z.TerminMs, z.PoczatekMs)
		return chwila, chwila
	}
	return poczatek, koniec
}

// pierwszaChwilaWorkspace wybiera pierwszą chwilę różną od zera spośród
// podanych znaczników czasu początku i końca zadania.
func pierwszaChwilaWorkspace(chwile ...int64) int64 {
	for _, chwila := range chwile {
		if chwila != 0 {
			return chwila
		}
	}
	return 0
}

// domykaCyklWorkspace sprawdza, czy nowa krawędź zamknęłaby obieg: idzie od
// następnika po istniejących krawędziach i szuka poprzednika.
func domykaCyklWorkspace(zaleznosci []dane.ZaleznoscWorkspace, poprzednik, nastepnik string) bool {
	nastepni := map[string][]string{}
	for _, zaleznosc := range zaleznosci {
		nastepni[zaleznosc.Poprzednik] = append(nastepni[zaleznosc.Poprzednik], zaleznosc.Nastepnik)
	}
	odwiedzone := map[string]bool{}
	kolejka := []string{nastepnik}
	for len(kolejka) > 0 {
		biezacy := kolejka[0]
		kolejka = kolejka[1:]
		if biezacy == poprzednik {
			return true
		}
		if odwiedzone[biezacy] {
			continue
		}
		odwiedzone[biezacy] = true
		kolejka = append(kolejka, nastepni[biezacy]...)
	}
	return false
}

// sciezkaKrytycznaWorkspace wskazuje zadania leżące na najdłuższym łańcuchu
// zależności mierzonym czasem trwania. Zadania bez granic czasu w rachunku nie
// biorą udziału — nie mają czasu trwania, którym miałyby ważyć.
func sciezkaKrytycznaWorkspace(zadania []dane.ZadanieWorkspace,
	zaleznosci []dane.ZaleznoscWorkspace) map[string]bool {

	trwanie := map[string]int64{}
	for _, zadanie := range zadania {
		poczatek, koniec := granicaSlupkaWorkspace(zadanie)
		if poczatek == 0 || koniec == 0 {
			continue
		}
		trwanie[zadanie.Identyfikator] = koniec - poczatek
	}
	poprzedni := map[string][]string{}
	for _, zaleznosc := range zaleznosci {
		poprzedni[zaleznosc.Nastepnik] = append(poprzedni[zaleznosc.Nastepnik], zaleznosc.Poprzednik)
	}

	najdluzsza := map[string]int64{}
	rodzic := map[string]string{}
	// Rachunek idzie zapamiętanym przejściem w głąb, chroniącym przed zapętleniem.
	var policz func(kod string, glebokosc int) int64
	policz = func(kod string, glebokosc int) int64 {
		if wynik, jest := najdluzsza[kod]; jest {
			return wynik
		}
		if glebokosc > len(zadania)+1 {
			return trwanie[kod]
		}
		najlepszy, skad := int64(0), ""
		for _, wczesniejszy := range poprzedni[kod] {
			if _, liczony := trwanie[wczesniejszy]; !liczony {
				continue
			}
			wynik := policz(wczesniejszy, glebokosc+1)
			if wynik > najlepszy {
				najlepszy, skad = wynik, wczesniejszy
			}
		}
		najdluzsza[kod] = najlepszy + trwanie[kod]
		if skad != "" {
			rodzic[kod] = skad
		}
		return najdluzsza[kod]
	}

	koniecSciezki, najwieksza := "", int64(0)
	for kod := range trwanie {
		if wynik := policz(kod, 0); wynik > najwieksza {
			koniecSciezki, najwieksza = kod, wynik
		}
	}
	sciezka := map[string]bool{}
	for kod := koniecSciezki; kod != ""; kod = rodzic[kod] {
		if sciezka[kod] {
			break
		}
		sciezka[kod] = true
	}
	return sciezka
}

// zaleznoscKontraktuWorkspace przekłada wiersz zależności warstwy danych
// na byt kontraktu TaskDependency.
func zaleznoscKontraktuWorkspace(idProjektu string, z dane.ZaleznoscWorkspace) shared.WorkspaceTaskDependency {
	zaleznosc := shared.WorkspaceTaskDependency{
		Id: z.Identyfikator, ProjectId: idProjektu,
		PredecessorTaskId: z.Poprzednik, SuccessorTaskId: z.Nastepnik, Kind: z.Rodzaj,
	}
	if z.OdstepMinut != 0 {
		odstep := z.OdstepMinut
		zaleznosc.LagMinutes = &odstep
	}
	return zaleznosc
}
