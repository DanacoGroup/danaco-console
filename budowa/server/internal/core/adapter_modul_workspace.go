// Odpowiedzialność pliku: moduł Workspace — projekt jako byt i zestawienie Project Dashboard, zakładany przez workspace.dashboard.get, gdy jeszcze go nie ma.
package core

import (
	"context"

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
	pliki := a.plikiProjektu(projekt.Kod, "")

	ostatnia := chwilaBazy(projekt.Zaktualizowano)
	liczbaPlikow, liczbaWpisow := len(pliki), len(wpisy)
	return shared.WorkspaceDashboardGetResponse{
		Dashboard: shared.WorkspaceDashboard{
			Project:          projektKontraktu(projekt),
			OpenSessionIds:   sesje,
			AssignedAgentIds: kodyEkspertow(przypisania),
			LibraryFileCount: &liczbaPlikow,
			MemoryEntryCount: &liczbaWpisow,
			LastActivityAt:   &ostatnia,
		},
	}, nil
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

// bladProjektu znakuje błąd wskazania projektu kodem kontraktu, żeby klient
// odróżnił brak danych w żądaniu od usterki rdzenia.
func bladProjektu(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"moduł Workspace: "+powod))
}
