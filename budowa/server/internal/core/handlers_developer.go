// Plik wpina komendy obszaru `developer.*` — modułu Developer wraz z jego
// oknami operacyjnymi (Code Editor, Project Tree, Git Panel, Build Output).
//
// Kontrakt daje modułowi jedno zdarzenie, `developer.build.changed`, więc każdy
// przyrost budowania — start, kolejny wiersz logu, domknięcie — rozgłasza się
// przebiegiem po zmianie. Build Output odświeża się z jednej subskrypcji,
// a nie z odpytywania.
//
// Podział dróg w Git Panelu jest jawny: `developer.git.action` wykonuje czynności
// ZMIENIAJĄCE repozytorium wedle zamkniętego słownika, a `developer.git.status`,
// `.diff`, `.log`, `.branch.list` i `.conflict.*` wyłącznie CZYTAJĄ. Odczyt stoi
// na bibliotece `go-git` wkompilowanej w rdzeń i nie startuje ani jednego
// procesu — Git Panel ma się otwierać także tam, gdzie programu `git` nie ma,
// bo instalka go nie niesie.
package core

import (
	"context"

	"danacoconsole/shared"
)

// Developer jest portem modułu Developer.
type Developer interface {
	OtworzPlik(ctx context.Context, z shared.DeveloperFileOpenRequest) (shared.DeveloperFileOpenResponse, error)
	ZapiszPlik(ctx context.Context, z shared.DeveloperFileSaveRequest) (shared.DeveloperFileSaveResponse, error)
	Drzewo(ctx context.Context, z shared.DeveloperTreeGetRequest) (shared.DeveloperTreeGetResponse, error)
	CzynnoscRepozytorium(ctx context.Context, z shared.DeveloperGitActionRequest) (shared.DeveloperGitActionResponse, error)
	Budowanie(ctx context.Context, z shared.DeveloperBuildRunRequest) (shared.DeveloperBuildRunResponse, error)
	// PodepnijPrzyrostBudowania oddaje adapterowi drogę do zdarzenia przyrostu.
	// Budowanie kończy się poza wykonaniem komendy — czasem minuty później —
	// a log narasta przez cały ten czas, więc rozgłoszenie nie może iść
	// wyłącznie z obsługiwacza żądania.
	PodepnijPrzyrostBudowania(rozglos func(shared.ChangeKind, shared.DeveloperBuild, string))

	// Odczyt repozytorium — Git Panel. Sześć czynności na bibliotece `go-git`
	// wkompilowanej w rdzeń, bez ani jednego procesu potomnego.
	StanRepozytorium(ctx context.Context, z shared.DeveloperGitStatusRequest) (shared.DeveloperGitStatusResponse, error)
	RoznicaRepozytorium(ctx context.Context, z shared.DeveloperGitDiffRequest) (shared.DeveloperGitDiffResponse, error)
	HistoriaRepozytorium(ctx context.Context, z shared.DeveloperGitLogRequest) (shared.DeveloperGitLogResponse, error)
	GaleziRepozytorium(ctx context.Context, z shared.DeveloperGitBranchListRequest) (shared.DeveloperGitBranchListResponse, error)
	KonfliktRepozytorium(ctx context.Context, z shared.DeveloperGitConflictGetRequest) (shared.DeveloperGitConflictGetResponse, error)
	RozstrzygnijKonfliktRepozytorium(ctx context.Context, z shared.DeveloperGitConflictResolveRequest) (shared.DeveloperGitConflictResolveResponse, error)

	// Czynności na węzłach drzewa projektu oraz wyszukiwanie w repozytorium.
	ZalozWezel(ctx context.Context, z shared.DeveloperFileCreateRequest) (shared.DeveloperFileCreateResponse, error)
	ZmienNazweWezla(ctx context.Context, z shared.DeveloperFileRenameRequest) (shared.DeveloperFileRenameResponse, error)
	UsunWezly(ctx context.Context, z shared.DeveloperFileDeleteRequest) (shared.DeveloperFileDeleteResponse, error)
	PrzeniesWezly(ctx context.Context, z shared.DeveloperFileMoveRequest) (shared.DeveloperFileMoveResponse, error)
	SzukajWRepozytorium(ctx context.Context, z shared.DeveloperGrepSearchRequest) (shared.DeveloperGrepSearchResponse, error)
	ZamienWRepozytorium(ctx context.Context, z shared.DeveloperGrepReplaceRequest) (shared.DeveloperGrepReplaceResponse, error)

	// Historia pliku edytora — wersja robocza zakładana przy zapisie. To NIE
	// jest historia repozytorium: tamta należy do Gita i jedzie `git.log`.
	WykazWersjiPliku(ctx context.Context, z shared.DeveloperFileVersionListRequest) (shared.DeveloperFileVersionListResponse, error)
	PrzywrocWersjePliku(ctx context.Context, z shared.DeveloperFileVersionRestoreRequest) (shared.DeveloperFileVersionRestoreResponse, error)

	// Warstwa językowa Code Editora. Jako jedyna w module woła programy serwera
	// (gopls, gofmt, goimports, prettier, golangci-lint), bo program JEST tą
	// wiedzą o kodzie; brak programu wraca polem `serverAvailable`/`linterAvailable`,
	// a nie odmową całej komendy.
	NawigujDoSymbolu(ctx context.Context, z shared.DeveloperSymbolNavigateRequest) (shared.DeveloperSymbolNavigateResponse, error)
	Formatuj(ctx context.Context, z shared.DeveloperFormatRunRequest) (shared.DeveloperFormatRunResponse, error)
	AnalizaStatyczna(ctx context.Context, z shared.DeveloperLintGetRequest) (shared.DeveloperLintGetResponse, error)
	Refaktoryzuj(ctx context.Context, z shared.DeveloperRefactorApplyRequest) (shared.DeveloperRefactorApplyResponse, error)

	// Odczyt okna Build Output — historia przebiegów, log, wynik testów
	// i pokrycie. Wynik testów i pokrycie powstają z rozbioru wyjścia w chwili
	// biegu, a nie z ponownego czytania przyciętego dziennika.
	WykazBudowan(ctx context.Context, z shared.DeveloperBuildListRequest) (shared.DeveloperBuildListResponse, error)
	LogBudowania(ctx context.Context, z shared.DeveloperBuildLogGetRequest) (shared.DeveloperBuildLogGetResponse, error)
	WynikTestow(ctx context.Context, z shared.DeveloperTestResultGetRequest) (shared.DeveloperTestResultGetResponse, error)
	Pokrycie(ctx context.Context, z shared.DeveloperCoverageGetRequest) (shared.DeveloperCoverageGetResponse, error)

	// Okno Run & Debug na protokole DAP. Punkt przerwania jest trwały i należy
	// do okna; sesja debugowania jest żywa i gaśnie razem z procesem adaptera.
	UruchomDebugowanie(ctx context.Context, z shared.DeveloperDebugSessionStartRequest) (shared.DeveloperDebugSessionStartResponse, error)
	SterujDebugowaniem(ctx context.Context, z shared.DeveloperDebugSessionControlRequest) (shared.DeveloperDebugSessionControlResponse, error)
	UstawPunktPrzerwania(ctx context.Context, z shared.DeveloperBreakpointSetRequest) (shared.DeveloperBreakpointSetResponse, error)
	ZakresDebugowania(ctx context.Context, z shared.DeveloperDebugScopeGetRequest) (shared.DeveloperDebugScopeGetResponse, error)
	ObliczWyrazenie(ctx context.Context, z shared.DeveloperDebugEvaluateRequest) (shared.DeveloperDebugEvaluateResponse, error)

	// Zakładka API Client — `net/http` i `getkin/kin-openapi`, bez ani jednego
	// procesu potomnego.
	WykonajZapytanieApi(ctx context.Context, z shared.DeveloperApiRequestRequest) (shared.DeveloperApiRequestResponse, error)
	WykonajPrzebiegObciazeniowy(ctx context.Context, z shared.DeveloperApiLoadRunRequest) (shared.DeveloperApiLoadRunResponse, error)
	ZapiszKolekcjeApi(ctx context.Context, z shared.DeveloperApiCollectionSaveRequest) (shared.DeveloperApiCollectionSaveResponse, error)
	WykazKolekcjiApi(ctx context.Context, z shared.DeveloperApiCollectionListRequest) (shared.DeveloperApiCollectionListResponse, error)
	ImportujOpenapi(ctx context.Context, z shared.DeveloperApiOpenapiImportRequest) (shared.DeveloperApiOpenapiImportResponse, error)

	// Zakładka Data Console — `database/sql` ze sterownikami wkompilowanymi
	// w rdzeń. Hasło nie leży w opisie połączenia, tylko w sejfie.
	UstawPolaczenieDanych(ctx context.Context, z shared.DeveloperDataConnectionSetRequest) (shared.DeveloperDataConnectionSetResponse, error)
	WykazPolaczenDanych(ctx context.Context, z shared.DeveloperDataConnectionListRequest) (shared.DeveloperDataConnectionListResponse, error)
	SchematDanych(ctx context.Context, z shared.DeveloperDataSchemaGetRequest) (shared.DeveloperDataSchemaGetResponse, error)
	WykonajZapytanieDanych(ctx context.Context, z shared.DeveloperDataQueryRunRequest) (shared.DeveloperDataQueryRunResponse, error)
	UruchomMigracjeDanych(ctx context.Context, z shared.DeveloperDataMigrationRunRequest) (shared.DeveloperDataMigrationRunResponse, error)

	// Zakładka Containers — Docker SDK Go rozmawiający z gniazdem silnika.
	// Silnika nie da się wkompilować: jego brak wraca polem `engineAvailable`.
	WykazKontenerow(ctx context.Context, z shared.DeveloperContainerListRequest) (shared.DeveloperContainerListResponse, error)
	CzynnoscKontenera(ctx context.Context, z shared.DeveloperContainerActionRequest) (shared.DeveloperContainerActionResponse, error)
	BudujObraz(ctx context.Context, z shared.DeveloperImageBuildRequest) (shared.DeveloperImageBuildResponse, error)
	KompozycjaKontenerow(ctx context.Context, z shared.DeveloperComposeUpRequest) (shared.DeveloperComposeUpResponse, error)

	// Zależności, bezpieczeństwo i jakość — własne parsery manifestów i własne
	// reguły skanowania, bez zależności od programu spoza instalki.
	WykazZaleznosci(ctx context.Context, z shared.DeveloperDependencyListRequest) (shared.DeveloperDependencyListResponse, error)
	UruchomSkan(ctx context.Context, z shared.DeveloperScanRunRequest) (shared.DeveloperScanRunResponse, error)
	WykazZnalezisk(ctx context.Context, z shared.DeveloperScanResultListRequest) (shared.DeveloperScanResultListResponse, error)

	// Operacje kontekstowe paska pływającego oraz sonda programów warsztatu.
	OperacjaKontekstowa(ctx context.Context, z shared.DeveloperContextualOpRequest) (shared.DeveloperContextualOpResponse, error)
	SprawdzWarsztat(ctx context.Context, z shared.DeveloperToolchainCheckRequest) (shared.DeveloperToolchainCheckResponse, error)
}

// zarejestrujDevelopera wpina komplet pięćdziesięciu komend modułu Developer.
func zarejestrujDevelopera(r *Rejestr, d Developer, e *emiter) {
	if r == nil || d == nil {
		return
	}
	d.PodepnijPrzyrostBudowania(e.przyrostBudowania)

	r.Zarejestruj(shared.CommandDeveloperFileOpen, obsluz(d.OtworzPlik))
	r.Zarejestruj(shared.CommandDeveloperFileSave, obsluz(d.ZapiszPlik))
	r.Zarejestruj(shared.CommandDeveloperTreeGet, obsluz(d.Drzewo))
	r.Zarejestruj(shared.CommandDeveloperGitAction, obsluz(d.CzynnoscRepozytorium))
	r.Zarejestruj(shared.CommandDeveloperBuildRun, obsluz(d.Budowanie))

	r.Zarejestruj(shared.CommandDeveloperGitStatus, obsluz(d.StanRepozytorium))
	r.Zarejestruj(shared.CommandDeveloperGitDiff, obsluz(d.RoznicaRepozytorium))
	r.Zarejestruj(shared.CommandDeveloperGitLog, obsluz(d.HistoriaRepozytorium))
	r.Zarejestruj(shared.CommandDeveloperGitBranchList, obsluz(d.GaleziRepozytorium))
	r.Zarejestruj(shared.CommandDeveloperGitConflictGet, obsluz(d.KonfliktRepozytorium))
	r.Zarejestruj(shared.CommandDeveloperGitConflictResolve, obsluz(d.RozstrzygnijKonfliktRepozytorium))

	r.Zarejestruj(shared.CommandDeveloperFileCreate, obsluz(d.ZalozWezel))
	r.Zarejestruj(shared.CommandDeveloperFileRename, obsluz(d.ZmienNazweWezla))
	r.Zarejestruj(shared.CommandDeveloperFileDelete, obsluz(d.UsunWezly))
	r.Zarejestruj(shared.CommandDeveloperFileMove, obsluz(d.PrzeniesWezly))
	r.Zarejestruj(shared.CommandDeveloperGrepSearch, obsluz(d.SzukajWRepozytorium))
	r.Zarejestruj(shared.CommandDeveloperGrepReplace, obsluz(d.ZamienWRepozytorium))

	r.Zarejestruj(shared.CommandDeveloperFileVersionList, obsluz(d.WykazWersjiPliku))
	r.Zarejestruj(shared.CommandDeveloperFileVersionRestore, obsluz(d.PrzywrocWersjePliku))

	r.Zarejestruj(shared.CommandDeveloperSymbolNavigate, obsluz(d.NawigujDoSymbolu))
	r.Zarejestruj(shared.CommandDeveloperFormatRun, obsluz(d.Formatuj))
	r.Zarejestruj(shared.CommandDeveloperLintGet, obsluz(d.AnalizaStatyczna))
	r.Zarejestruj(shared.CommandDeveloperRefactorApply, obsluz(d.Refaktoryzuj))

	r.Zarejestruj(shared.CommandDeveloperBuildList, obsluz(d.WykazBudowan))
	r.Zarejestruj(shared.CommandDeveloperBuildLogGet, obsluz(d.LogBudowania))
	r.Zarejestruj(shared.CommandDeveloperTestResultGet, obsluz(d.WynikTestow))
	r.Zarejestruj(shared.CommandDeveloperCoverageGet, obsluz(d.Pokrycie))

	r.Zarejestruj(shared.CommandDeveloperDebugSessionStart, obsluz(d.UruchomDebugowanie))
	r.Zarejestruj(shared.CommandDeveloperDebugSessionControl, obsluz(d.SterujDebugowaniem))
	r.Zarejestruj(shared.CommandDeveloperBreakpointSet, obsluz(d.UstawPunktPrzerwania))
	r.Zarejestruj(shared.CommandDeveloperDebugScopeGet, obsluz(d.ZakresDebugowania))
	r.Zarejestruj(shared.CommandDeveloperDebugEvaluate, obsluz(d.ObliczWyrazenie))

	r.Zarejestruj(shared.CommandDeveloperApiRequest, obsluz(d.WykonajZapytanieApi))
	// Przebieg obciążeniowy stoi przy zapytaniu pojedynczym, bo jest tym samym
	// zapytaniem powtórzonym pod obciążeniem — z tym samym podstawianiem zmiennych
	// środowiska kolekcji.
	r.Zarejestruj(shared.CommandDeveloperApiLoadRun, obsluz(d.WykonajPrzebiegObciazeniowy))
	r.Zarejestruj(shared.CommandDeveloperApiCollectionSave, obsluz(d.ZapiszKolekcjeApi))
	r.Zarejestruj(shared.CommandDeveloperApiCollectionList, obsluz(d.WykazKolekcjiApi))
	r.Zarejestruj(shared.CommandDeveloperApiOpenapiImport, obsluz(d.ImportujOpenapi))

	r.Zarejestruj(shared.CommandDeveloperDataConnectionSet, obsluz(d.UstawPolaczenieDanych))
	r.Zarejestruj(shared.CommandDeveloperDataConnectionList, obsluz(d.WykazPolaczenDanych))
	r.Zarejestruj(shared.CommandDeveloperDataSchemaGet, obsluz(d.SchematDanych))
	r.Zarejestruj(shared.CommandDeveloperDataQueryRun, obsluz(d.WykonajZapytanieDanych))
	r.Zarejestruj(shared.CommandDeveloperDataMigrationRun, obsluz(d.UruchomMigracjeDanych))

	r.Zarejestruj(shared.CommandDeveloperContainerList, obsluz(d.WykazKontenerow))
	r.Zarejestruj(shared.CommandDeveloperContainerAction, obsluz(d.CzynnoscKontenera))
	r.Zarejestruj(shared.CommandDeveloperImageBuild, obsluz(d.BudujObraz))
	r.Zarejestruj(shared.CommandDeveloperComposeUp, obsluz(d.KompozycjaKontenerow))

	r.Zarejestruj(shared.CommandDeveloperDependencyList, obsluz(d.WykazZaleznosci))
	r.Zarejestruj(shared.CommandDeveloperScanRun, obsluz(d.UruchomSkan))
	r.Zarejestruj(shared.CommandDeveloperScanResultList, obsluz(d.WykazZnalezisk))

	r.Zarejestruj(shared.CommandDeveloperContextualOp, obsluz(d.OperacjaKontekstowa))
	r.Zarejestruj(shared.CommandDeveloperToolchainCheck, obsluz(d.SprawdzWarsztat))
}

// przyrostBudowania rozgłasza przyrost przebiegu. Wiersz logu jedzie osobnym
// polem kontraktu, a nie doklejony do przebiegu: Build Output dopisuje go do
// ogona bez przerysowywania całego widoku.
func (e *emiter) przyrostBudowania(zmiana shared.ChangeKind, budowanie shared.DeveloperBuild,
	wiersz string) {

	tresc := shared.DeveloperBuildChangedEvent{Change: zmiana, Build: budowanie}
	if wiersz != "" {
		tresc.LogLine = &wiersz
	}
	// Sesja komunikatu zostaje pusta: przebieg należy do okna, a rdzeń rozgłasza
	// go także wtedy, gdy okna nie ma już w rejestrze — budowanie przeżywa
	// zamknięcie okna, a klient ma prawo zobaczyć jego koniec.
	e.wyslij(shared.EventDeveloperBuildChanged, "", tresc)
}

// PodepnijPrzyrostBudowania wypełnia port: adapter zapamiętuje drogę do zdarzenia.
func (a *adapterDevelopera) PodepnijPrzyrostBudowania(
	rozglos func(shared.ChangeKind, shared.DeveloperBuild, string)) {

	a.przyrost = rozglos
}
