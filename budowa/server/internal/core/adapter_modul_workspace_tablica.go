// Odpowiedzialność pliku: tablica kanban projektu — kolumny, kolejność kart
// i przeniesienie karty (`workspace.board.get`, `workspace.task.move`).
//
// ── Klucz porządkowy jest napisem ──────────────────────────────────────────
// Karta ma klucz porządkowy z liter, a nie numer pozycji. Przeciągnięcie karty
// między dwie sąsiednie dopisuje klucz leżący pomiędzy ich kluczami — jeden
// zapis. Numer pozycji wymagałby przepisania całej kolumny przy każdym
// przeciągnięciu, a tablica projektu bywa przeciągana kilkanaście razy pod rząd.
//
// ── Granica prac w toku ostrzega, nie odmawia ──────────────────────────────
// Przekroczenie granicy WIP wraca polem `wipExceeded`, a karta i tak staje
// w kolumnie. Platforma nie stawia twardych blokad w interfejsie: granica jest
// sygnałem dla Operatora, nie bramką.
package core

import (
	"context"
	"sort"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// kolumnyWyjscioweWorkspace to zestaw wyjściowy kolumn tablicy: po jednej na
// każdy stan zadania z kontraktu. Kolumna własna Operatora stanu nowego nie
// zakłada, więc ten zestaw pokrywa tablicę w całości.
var kolumnyWyjscioweWorkspace = []struct {
	stan  shared.WorkspaceTaskStatus
	nazwa string
}{
	{shared.WorkspaceTaskStatusTodo, "Do zrobienia"},
	{shared.WorkspaceTaskStatusInProgress, "W realizacji"},
	{shared.WorkspaceTaskStatusInReview, "W kontroli"},
	{shared.WorkspaceTaskStatusBlocked, "Zablokowane"},
	{shared.WorkspaceTaskStatusDone, "Ukończone"},
	{shared.WorkspaceTaskStatusCancelled, "Odwołane"},
}

// Tablica obsługuje `workspace.board.get`.
func (a *adapterPrzestrzeniRoboczej) Tablica(ctx context.Context,
	z shared.WorkspaceBoardGetRequest) (shared.WorkspaceBoardGetResponse, error) {

	projekt, err := a.projektDlaZapisu(ctx, z.ProjectId)
	if err != nil {
		return shared.WorkspaceBoardGetResponse{}, err
	}
	kolumny, err := a.kolumnyTablicyWorkspace(ctx, projekt.ID)
	if err != nil {
		return shared.WorkspaceBoardGetResponse{}, err
	}
	zadania, err := a.repozytorium.ZadaniaWorkspace(ctx, projekt.ID)
	if err != nil {
		return shared.WorkspaceBoardGetResponse{}, err
	}
	dolaczZamkniete := z.IncludeDone == nil || *z.IncludeDone

	wykazKolumn := []shared.WorkspaceBoardColumn{}
	widoczne := map[string]bool{}
	for _, kolumna := range kolumny {
		if !dolaczZamkniete && stanZamknietyWorkspace(kolumna.Stan) {
			continue
		}
		widoczne[kolumna.Identyfikator] = true
		wykazKolumn = append(wykazKolumn, kolumnaKontraktuWorkspace(kolumna))
	}
	karty := []dane.ZadanieWorkspace{}
	for _, zadanie := range zadania {
		if !dolaczZamkniete && stanZamknietyWorkspace(zadanie.Stan) {
			continue
		}
		if zadanie.KolumnaTablicy != "" && !widoczne[zadanie.KolumnaTablicy] && !dolaczZamkniete {
			continue
		}
		karty = append(karty, zadanie)
	}
	sort.SliceStable(karty, func(i, j int) bool {
		if karty[i].KolumnaTablicy != karty[j].KolumnaTablicy {
			return karty[i].KolumnaTablicy < karty[j].KolumnaTablicy
		}
		return karty[i].KluczPorzadkowy < karty[j].KluczPorzadkowy
	})
	wykazKart := make([]shared.WorkspaceTask, 0, len(karty))
	for _, karta := range karty {
		wykazKart = append(wykazKart, zadanieKontraktuWorkspace(karta))
	}
	return shared.WorkspaceBoardGetResponse{Board: shared.WorkspaceBoard{
		ProjectId: projekt.Kod, Columns: wykazKolumn, Tasks: wykazKart,
	}}, nil
}

// PrzeniesZadanie obsługuje `workspace.task.move`.
func (a *adapterPrzestrzeniRoboczej) PrzeniesZadanie(ctx context.Context,
	z shared.WorkspaceTaskMoveRequest) (shared.WorkspaceTaskMoveResponse, error) {

	zadanie, err := a.zadanieDoZmianyWorkspace(ctx, z.TaskId)
	if err != nil {
		return shared.WorkspaceTaskMoveResponse{}, err
	}
	kolumny, err := a.kolumnyTablicyWorkspace(ctx, zadanie.ProjektID)
	if err != nil {
		return shared.WorkspaceTaskMoveResponse{}, err
	}
	docelowa := zadanie.KolumnaTablicy
	if z.BoardColumnId != nil && *z.BoardColumnId != "" {
		docelowa = *z.BoardColumnId
	}
	stanKolumny := shared.WorkspaceTaskStatus("")
	for _, kolumna := range kolumny {
		if kolumna.Identyfikator == docelowa {
			stanKolumny = kolumna.Stan
		}
	}
	if docelowa != "" && stanKolumny == "" {
		return shared.WorkspaceTaskMoveResponse{},
			bladProjektu("kolumny " + docelowa + " nie ma na tablicy projektu")
	}
	switch {
	case z.Status != nil && *z.Status != "":
		zadanie.Stan = *z.Status
	case stanKolumny != "":
		zadanie.Stan = stanKolumny
	}
	if zadanie.Stan == shared.WorkspaceTaskStatusDone && zadanie.UkonczonoMs == 0 {
		zadanie.UkonczonoMs = terazWMilisekundachWorkspace()
	}
	if zadanie.Stan != shared.WorkspaceTaskStatusDone {
		zadanie.UkonczonoMs = 0
	}
	zadanie.KolumnaTablicy = docelowa

	// „Przed którą kartą" wyznacza granicę prawą, „za którą" — lewą. Bez obu
	// wskazań karta staje na końcu kolumny: przeniesienie samą zmianą kolumny
	// nie ma powodu wchodzić między karty już ułożone.
	rangaKartyPrzed, rangaKartyZa, err := a.sasiedziKartyWorkspace(ctx, zadanie.ProjektID,
		z.BeforeTaskId, z.AfterTaskId)
	if err != nil {
		return shared.WorkspaceTaskMoveResponse{}, err
	}
	if rangaKartyPrzed == "" && rangaKartyZa == "" {
		rangaKartyZa, err = a.ostatniaRangaWorkspace(ctx, zadanie.ProjektID, docelowa)
		if err != nil {
			return shared.WorkspaceTaskMoveResponse{}, err
		}
	}
	zadanie.KluczPorzadkowy = rangaMiedzyWorkspace(rangaKartyZa, rangaKartyPrzed)

	zapisane, err := a.repozytorium.ZapiszZadanieWorkspace(ctx, zadanie)
	if err != nil {
		return shared.WorkspaceTaskMoveResponse{}, err
	}
	if err := a.repozytorium.OdnotujCzynnosc(ctx, zapisane.ProjektID); err != nil {
		return shared.WorkspaceTaskMoveResponse{}, err
	}
	a.odnotujZdarzenieWorkspace(ctx, zapisane.ProjektID, shared.WorkspaceActivityKindTaskChanged,
		shared.ChangeKindUpdated, shared.WorkspaceEntityKindTask, zapisane.Identyfikator,
		"karta „"+zapisane.Tytul+"” przeniesiona do stanu "+string(zapisane.Stan))

	przekroczona, err := a.granicaPrzekroczonaWorkspace(ctx, zapisane.ProjektID, docelowa, kolumny)
	if err != nil {
		return shared.WorkspaceTaskMoveResponse{}, err
	}
	return shared.WorkspaceTaskMoveResponse{
		Task: zadanieKontraktuWorkspace(zapisane), WipExceeded: &przekroczona,
	}, nil
}

// kolumnyTablicyWorkspace zwraca kolumny projektu, zakładając zestaw wyjściowy
// przy pierwszym wejściu na tablicę. Kontrakt nie ma komendy zakładającej
// kolumnę, a tablica bez kolumn nie miałaby gdzie postawić kart.
func (a *adapterPrzestrzeniRoboczej) kolumnyTablicyWorkspace(ctx context.Context,
	projektID int64) ([]dane.KolumnaTablicyWorkspace, error) {

	kolumny, err := a.repozytorium.KolumnyTablicyWorkspace(ctx, projektID)
	if err != nil {
		return nil, err
	}
	if len(kolumny) > 0 {
		return kolumny, nil
	}
	for numer, wzorzec := range kolumnyWyjscioweWorkspace {
		kolumna := dane.KolumnaTablicyWorkspace{
			ProjektID:     projektID,
			Identyfikator: "kolumna-" + string(wzorzec.stan),
			Nazwa:         wzorzec.nazwa,
			Stan:          wzorzec.stan,
			Kolejnosc:     numer,
		}
		if err := a.repozytorium.ZapiszKolumneTablicyWorkspace(ctx, kolumna); err != nil {
			return nil, err
		}
	}
	return a.repozytorium.KolumnyTablicyWorkspace(ctx, projektID)
}

// kolumnaStanuWorkspace wskazuje kolumnę odwzorowującą stan zadania.
func (a *adapterPrzestrzeniRoboczej) kolumnaStanuWorkspace(ctx context.Context, projektID int64,
	stan shared.WorkspaceTaskStatus) (string, error) {

	kolumny, err := a.kolumnyTablicyWorkspace(ctx, projektID)
	if err != nil {
		return "", err
	}
	for _, kolumna := range kolumny {
		if kolumna.Stan == stan {
			return kolumna.Identyfikator, nil
		}
	}
	return "", nil
}

// ostatniaRangaWorkspace zwraca największy klucz porządkowy w kolumnie —
// miejsce, za którym staje karta dołożona.
func (a *adapterPrzestrzeniRoboczej) ostatniaRangaWorkspace(ctx context.Context, projektID int64,
	kolumna string) (string, error) {

	zadania, err := a.repozytorium.ZadaniaWorkspace(ctx, projektID)
	if err != nil {
		return "", err
	}
	ostatnia := ""
	for _, zadanie := range zadania {
		if zadanie.KolumnaTablicy != kolumna {
			continue
		}
		if zadanie.KluczPorzadkowy > ostatnia {
			ostatnia = zadanie.KluczPorzadkowy
		}
	}
	return ostatnia, nil
}

// sasiedziKartyWorkspace odczytuje klucze porządkowe kart wskazanych żądaniem.
func (a *adapterPrzestrzeniRoboczej) sasiedziKartyWorkspace(ctx context.Context, projektID int64,
	przedKarta, zaKarta *string) (string, string, error) {

	odczytaj := func(wskazanie *string) (string, error) {
		if wskazanie == nil || *wskazanie == "" {
			return "", nil
		}
		zadanie, err := a.repozytorium.ZadanieWorkspace(ctx, *wskazanie)
		if err != nil {
			return "", nil
		}
		if zadanie.ProjektID != projektID {
			return "", bladProjektu("karta " + *wskazanie + " należy do innego projektu")
		}
		return zadanie.KluczPorzadkowy, nil
	}
	przed, err := odczytaj(przedKarta)
	if err != nil {
		return "", "", err
	}
	za, err := odczytaj(zaKarta)
	if err != nil {
		return "", "", err
	}
	return przed, za, nil
}

// granicaPrzekroczonaWorkspace mówi, czy kolumna docelowa przekroczyła własną
// granicę prac w toku. Odpowiedź jest sygnałem — karta stoi w kolumnie tak czy
// inaczej.
func (a *adapterPrzestrzeniRoboczej) granicaPrzekroczonaWorkspace(ctx context.Context,
	projektID int64, kolumna string, kolumny []dane.KolumnaTablicyWorkspace) (bool, error) {

	granica := 0
	for _, pozycja := range kolumny {
		if pozycja.Identyfikator == kolumna {
			granica = pozycja.GranicaWip
		}
	}
	if granica <= 0 {
		return false, nil
	}
	zadania, err := a.repozytorium.ZadaniaWorkspace(ctx, projektID)
	if err != nil {
		return false, err
	}
	stojace := 0
	for _, zadanie := range zadania {
		if zadanie.KolumnaTablicy == kolumna {
			stojace++
		}
	}
	return stojace > granica, nil
}

// kolumnaKontraktuWorkspace przekłada wiersz kolumny na byt kontraktu.
func kolumnaKontraktuWorkspace(k dane.KolumnaTablicyWorkspace) shared.WorkspaceBoardColumn {
	kolumna := shared.WorkspaceBoardColumn{
		Id: k.Identyfikator, Name: k.Nazwa, Status: k.Stan, Order: k.Kolejnosc,
	}
	if k.GranicaWip > 0 {
		granica := k.GranicaWip
		kolumna.WipLimit = &granica
	}
	return kolumna
}

// rangaMiedzyWorkspace układa klucz porządkowy leżący ściśle między dwoma
// kluczami sąsiadów. Pusty klucz z lewej znaczy początek kolumny, pusty
// z prawej — jej koniec.
//
// Klucz składa się z liter `a`–`z`; wynik jest zawsze większy od lewego
// i mniejszy od prawego, a przy sąsiadujących literach schodzi o znak niżej
// zamiast oddawać klucz równy któremuś z sąsiadów.
func rangaMiedzyWorkspace(lewy, prawy string) string {
	const dolna, gorna = byte('a'), byte('z')
	if prawy != "" && lewy >= prawy {
		// Wskazania sprzeczne (karta „przed" leży za kartą „za") nie mogą
		// zatrzymać przeniesienia: bierzemy stronę lewą i stawiamy kartę za nią.
		prawy = ""
	}
	wynik := make([]byte, 0, 8)
	for i := 0; ; i++ {
		znakLewy := dolna - 1
		if i < len(lewy) {
			znakLewy = lewy[i]
		}
		znakPrawy := gorna + 1
		if i < len(prawy) {
			znakPrawy = prawy[i]
		}
		if znakPrawy-znakLewy > 1 {
			return string(append(wynik, znakLewy+(znakPrawy-znakLewy)/2))
		}
		if i < len(lewy) {
			wynik = append(wynik, lewy[i])
			continue
		}
		wynik = append(wynik, dolna)
	}
}

// terazWMilisekundachWorkspace oddaje bieżącą chwilę w mierze kontraktu.
func terazWMilisekundachWorkspace() int64 {
	return time.Now().UnixMilli()
}

// pierwszaNiepustaWorkspace wybiera pierwszy niepusty napis — służy nazwom,
// które kontrakt podaje opcjonalnie, a byt musi mieć zawsze.
func pierwszaNiepustaWorkspace(wartosci ...string) string {
	for _, wartosc := range wartosci {
		if strings.TrimSpace(wartosc) != "" {
			return wartosc
		}
	}
	return ""
}
