// Odpowiedzialność pliku: moduł Workspace — projekt jako byt, jego zarząd rodziną project.* oraz zestawienie Project Dashboard; workspace.dashboard.get zakłada projekt, gdy jeszcze go nie ma.
package core

import (
	"context"
	"errors"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// adapterPrzestrzeniRoboczej wypełnia port PrzestrzenRobocza. Zależności są
// trzy: repozytorium modułu, rozstrzygacz ośmiu poziomów zasięgu wraz
// z repozytorium konfiguracji (instrukcje warstwowe) oraz ustalacz
// katalogu roboczego (biblioteka projektu).
type adapterPrzestrzeniRoboczej struct {
	repozytorium dane.RepozytoriumPrzestrzeniRoboczej
	konfiguracja dane.RepozytoriumKonfiguracji
	rozstrzygacz *konfig.Rozstrzygacz
	katalog      *KatalogRoboczy
	// dokumenty jest warsztatem odczytu treści plików, tym samym, którym jedzie document.text.extract.
	dokumenty WydobycieTekstuDokumentu
	// sesje odpina skasowany projekt od sesji żywych w rejestrze nadzorcy — baza
	// odpięła je w transakcji kasowania, rejestr w pamięci sam tego nie widzi.
	sesje OdpinanieProjektuZywych
}

// OdpinanieProjektuZywych zdejmuje kod projektu z sesji żywych i oddaje ich identyfikatory.
type OdpinanieProjektuZywych interface {
	OdepnijProjektZywych(kod string) []string
}

// ZSesjami wpina rejestr sesji żywych; bez niego project.delete zostawia kod projektu w pamięci nadzorcy.
func (a *adapterPrzestrzeniRoboczej) ZSesjami(sesje OdpinanieProjektuZywych) *adapterPrzestrzeniRoboczej {
	a.sesje = sesje
	return a
}

// nowyAdapterPrzestrzeniRoboczej wiąże nowo utworzony port z repozytorium modułu przestrzeni roboczej.
func nowyAdapterPrzestrzeniRoboczej(repozytorium dane.RepozytoriumPrzestrzeniRoboczej) *adapterPrzestrzeniRoboczej {
	return &adapterPrzestrzeniRoboczej{repozytorium: repozytorium}
}

// ZInstrukcjami podpina zapis instrukcji i rozstrzygacz warstw. Bez nich moduł
// pracuje dalej, a instrukcje odpowiadają brakiem zapisu.
func (a *adapterPrzestrzeniRoboczej) ZInstrukcjami(konfiguracja dane.RepozytoriumKonfiguracji,
	rozstrzygacz *konfig.Rozstrzygacz) *adapterPrzestrzeniRoboczej {

	a.konfiguracja, a.rozstrzygacz = konfiguracja, rozstrzygacz
	return a
}

// ZKatalogiem podpina ustalacz katalogu roboczego, źródło biblioteki projektu w drzewie plików rdzenia.
func (a *adapterPrzestrzeniRoboczej) ZKatalogiem(katalog *KatalogRoboczy) *adapterPrzestrzeniRoboczej {
	a.katalog = katalog
	return a
}

// Pulpit zwraca zestawienie stanu projektu: karty sesji pracujące w projekcie, przypisanych ekspertów, liczbę plików biblioteki, wpisów pamięci i czas ostatniej czynności, z liczbami wypełnianymi zawsze.
func (a *adapterPrzestrzeniRoboczej) Pulpit(ctx context.Context,
	z shared.WorkspaceDashboardGetRequest) (shared.WorkspaceDashboardGetResponse, error) {

	projekt, err := a.projektDlaZapisu(ctx, z.ProjectId)
	if err != nil {
		return shared.WorkspaceDashboardGetResponse{}, err
	}
	sesje, err := a.repozytorium.SesjeProjektu(ctx, projekt.Kod)
	if err != nil {
		return shared.WorkspaceDashboardGetResponse{}, err
	}
	przypisania, err := a.repozytorium.PrzypisaniaAgentow(ctx, projekt.ID)
	if err != nil {
		return shared.WorkspaceDashboardGetResponse{}, err
	}
	wpisy, err := a.repozytorium.WpisyPamieci(ctx, projekt.ID, 0)
	if err != nil {
		return shared.WorkspaceDashboardGetResponse{}, err
	}
	zadania, err := a.repozytorium.ZadaniaWorkspace(ctx, projekt.ID)
	if err != nil {
		return shared.WorkspaceDashboardGetResponse{}, err
	}
	notatki, err := a.repozytorium.NotatkiWorkspace(ctx, projekt.ID)
	if err != nil {
		return shared.WorkspaceDashboardGetResponse{}, err
	}
	instrukcje, err := a.repozytorium.WersjeInstrukcjiWorkspace(ctx, projekt.ID)
	if err != nil {
		return shared.WorkspaceDashboardGetResponse{}, err
	}
	pliki := a.plikiProjektu(projekt.Kod, "")

	ostatnia := chwilaBazy(projekt.Zaktualizowano)
	liczbaPlikow, liczbaWpisow := len(pliki), len(wpisy)
	liczbaZadan, liczbaNotatek := len(zadania), len(notatki)
	liczbaInstrukcji := len(instrukcje)
	liczbaUkonczonych := 0
	for _, zadanie := range zadania {
		if zadanie.Stan == shared.WorkspaceTaskStatusDone {
			liczbaUkonczonych++
		}
	}
	// Zajętość pamięci liczy się treścią wpisów: magazyn nie prowadzi dla nich
	// osobnej miary, a pasek pojemności bez liczby nic Operatorowi nie mówi.
	zajetoscPamieci := int64(0)
	for _, wpis := range wpisy {
		zajetoscPamieci += int64(len(wpis.Tresc))
	}
	return shared.WorkspaceDashboardGetResponse{
		Dashboard: shared.WorkspaceDashboard{
			Project:             projektKontraktu(projekt),
			OpenSessionIds:      sesje,
			AssignedAgentIds:    kodyEkspertow(przypisania),
			LibraryFileCount:    &liczbaPlikow,
			MemoryEntryCount:    &liczbaWpisow,
			LastActivityAt:      &ostatnia,
			TaskCount:           &liczbaZadan,
			TaskDoneCount:       &liczbaUkonczonych,
			NoteCount:           &liczbaNotatek,
			InstructionSetCount: &liczbaInstrukcji,
			MemoryUsedBytes:     &zajetoscPamieci,
		},
	}, nil
}

// Projekty oddaje wykaz projektów konta. Lewy panel ramy jest jedynym miejscem,
// w którym widać całość dorobku Operatora, więc wykaz obejmuje także projekty
// bez ani jednej sesji.
func (a *adapterPrzestrzeniRoboczej) Projekty(ctx context.Context,
	z shared.ProjectListRequest) (shared.ProjectListResponse, error) {

	zArchiwalnymi := z.IncludeArchived != nil && *z.IncludeArchived
	projekty, err := a.repozytorium.Projekty(ctx, zArchiwalnymi)
	if err != nil {
		return shared.ProjectListResponse{}, err
	}
	wykaz := make([]shared.WorkspaceProject, 0, len(projekty))
	for _, projekt := range projekty {
		wykaz = append(wykaz, projektKontraktu(projekt))
	}
	return shared.ProjectListResponse{Projects: wykaz, Total: len(wykaz)}, nil
}

// ZalozProjekt obsługuje `project.create`. Kod projektu nadaje rdzeń tym samym
// przedrostkiem, którym znakuje projekt zakładany przy przenoszeniu sesji —
// Operator podaje nazwę, bo nazwa jest tym, co widzi w lewym panelu ramy.
func (a *adapterPrzestrzeniRoboczej) ZalozProjekt(ctx context.Context,
	z shared.ProjectCreateRequest) (shared.ProjectCreateResponse, error) {

	nazwa := strings.TrimSpace(z.Name)
	if nazwa == "" {
		return shared.ProjectCreateResponse{}, bladProjektu("założenie projektu wymaga nazwy")
	}
	projekt, err := a.repozytorium.ZalozProjekt(ctx, nowyIdentyfikator(przedrostekProjektu),
		nazwa, opisProjektu(z.Description))
	if err != nil {
		return shared.ProjectCreateResponse{}, err
	}
	return shared.ProjectCreateResponse{Project: projektKontraktu(projekt)}, nil
}

// PrzemianujProjekt obsługuje `project.rename`. Zmiana dotyczy samej nazwy: kod
// projektu wiąże sesje i wiersze modułu, więc zostaje nietknięty.
func (a *adapterPrzestrzeniRoboczej) PrzemianujProjekt(ctx context.Context,
	z shared.ProjectRenameRequest) (shared.ProjectRenameResponse, error) {

	kod := strings.TrimSpace(z.ProjectId)
	nazwa := strings.TrimSpace(z.Name)
	if kod == "" {
		return shared.ProjectRenameResponse{}, bladProjektu("przemianowanie bez wskazania projektu")
	}
	if nazwa == "" {
		return shared.ProjectRenameResponse{}, bladProjektu("przemianowanie wymaga nowej nazwy")
	}
	projekt, err := a.repozytorium.PrzemianujProjekt(ctx, kod, nazwa)
	if err != nil {
		if errors.Is(err, dane.ErrBrakWiersza) {
			return shared.ProjectRenameResponse{}, bladNieznanegoProjektu(kod)
		}
		return shared.ProjectRenameResponse{}, err
	}
	return shared.ProjectRenameResponse{Project: projektKontraktu(projekt)}, nil
}

// UsunProjekt obsługuje `project.delete`. Sesje projektu zostają w historii
// i tracą przypisanie — kontrakt zapowiada to wprost, a wykaz odpiętych kart
// wraca Operatorowi, żeby wiedział, co zostało bez projektu.
func (a *adapterPrzestrzeniRoboczej) UsunProjekt(ctx context.Context,
	z shared.ProjectDeleteRequest) (shared.ProjectDeleteResponse, error) {

	kod := strings.TrimSpace(z.ProjectId)
	if kod == "" {
		return shared.ProjectDeleteResponse{}, bladProjektu("usunięcie bez wskazania projektu")
	}
	odpiete, err := a.repozytorium.UsunProjekt(ctx, kod)
	if err != nil {
		if errors.Is(err, dane.ErrBrakWiersza) {
			return shared.ProjectDeleteResponse{}, bladNieznanegoProjektu(kod)
		}
		return shared.ProjectDeleteResponse{}, err
	}
	if a.sesje != nil {
		a.sesje.OdepnijProjektZywych(kod)
	}
	return shared.ProjectDeleteResponse{ProjectId: kod, ReleasedSessionIds: odpiete}, nil
}

// Projekt oddaje projekt kontraktu — obsługiwacz potrzebuje go do rozgłoszenia
// `workspace.project.changed` po komendach, których wynik projektu nie niesie.
func (a *adapterPrzestrzeniRoboczej) Projekt(ctx context.Context,
	idProjektu string) (shared.WorkspaceProject, error) {

	projekt, err := a.repozytorium.Projekt(ctx, idProjektu)
	if err != nil {
		return shared.WorkspaceProject{}, err
	}
	return projektKontraktu(projekt), nil
}

// projektDlaZapisu zwraca projekt, do którego trafia komenda, zakładając go
// przy pierwszym wejściu. Komenda pamięci albo przypisania eksperta zastaje
// projekt gotowy również wtedy, gdy pulpit nie był jeszcze otwierany.
func (a *adapterPrzestrzeniRoboczej) projektDlaZapisu(ctx context.Context,
	idProjektu string) (dane.Projekt, error) {

	if idProjektu == "" {
		return dane.Projekt{}, bladProjektu("komenda bez wskazania projektu")
	}
	projekt, _, err := a.repozytorium.ZapewnijProjekt(ctx, idProjektu, idProjektu)
	return projekt, err
}

// kodyEkspertow wyciąga identyfikatory ekspertów z przypisań, zachowując
// kolejność ustaloną przez repozytorium.
func kodyEkspertow(przypisania []dane.PrzypisanieAgenta) []string {
	kody := make([]string, 0, len(przypisania))
	for _, przypisanie := range przypisania {
		kody = append(kody, przypisanie.AgentKod)
	}
	return kody
}

// projektKontraktu przekłada wiersz projektu na byt kontraktu wymiany z klientem panelu Project Dashboard.
func projektKontraktu(p dane.Projekt) shared.WorkspaceProject {
	stan := p.Stan
	if stan == "" {
		stan = shared.WorkspaceProjectStatusActive
	}
	return shared.WorkspaceProject{
		Id:          p.Kod,
		Name:        p.Nazwa,
		Description: p.Opis,
		Status:      stan,
		CreatedAt:   chwilaBazy(p.Utworzono),
		UpdatedAt:   chwilaBazy(p.Zaktualizowano),
	}
}

// opisProjektu zdejmuje z opisu białe znaki i zamienia opis pusty na brak
// wartości, bo kolumna opisu dopuszcza NULL zamiast pustego łańcucha znaków.
func opisProjektu(opis *string) *string {
	if opis == nil {
		return nil
	}
	przyciety := strings.TrimSpace(*opis)
	if przyciety == "" {
		return nil
	}
	return &przyciety
}

// bladNieznanegoProjektu znakuje wskazanie projektu, którego nie ma, kodem
// `not_found` — inaczej klient nie odróżniłby literówki w kodzie od usterki.
func bladNieznanegoProjektu(kod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
		"moduł Workspace: projekt "+kod+" nie istnieje"))
}

// bladProjektu znakuje błąd wskazania projektu kodem kontraktu, żeby klient
// odróżnił brak danych w żądaniu od usterki rdzenia.
func bladProjektu(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"moduł Workspace: "+powod))
}
