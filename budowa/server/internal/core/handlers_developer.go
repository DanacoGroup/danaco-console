// Plik wpina komendy obszaru `developer.*` modułu Developer wraz z jego oknami operacyjnymi: Code Editor, Project Tree, Git Panel i Build Output.
package core

import (
	"context"

	"danacoconsole/shared"
)

// Developer jest portem modułu Developer, wpiętym wraz z oknami operacyjnymi jego czterech zakładek roboczych.
type Developer interface {
	OtworzPlik(ctx context.Context, z shared.DeveloperFileOpenRequest) (shared.DeveloperFileOpenResponse, error)
	ZapiszPlik(ctx context.Context, z shared.DeveloperFileSaveRequest) (shared.DeveloperFileSaveResponse, error)
	Drzewo(ctx context.Context, z shared.DeveloperTreeGetRequest) (shared.DeveloperTreeGetResponse, error)
	CzynnoscRepozytorium(ctx context.Context, z shared.DeveloperGitActionRequest) (shared.DeveloperGitActionResponse, error)
	Budowanie(ctx context.Context, z shared.DeveloperBuildRunRequest) (shared.DeveloperBuildRunResponse, error)
	// PodepnijPrzyrostBudowania oddaje drogę do zdarzenia; budowanie kończy się poza wykonaniem komendy.
	PodepnijPrzyrostBudowania(rozglos func(shared.ChangeKind, shared.DeveloperBuild, string))

	// Odczyt repozytorium Git Panelu: sześć czynności na bibliotece `go-git`, bez procesu potomnego.
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

	// Historia pliku edytora to wersja robocza zakładana przy zapisie, nie historia repozytorium.
	WykazWersjiPliku(ctx context.Context, z shared.DeveloperFileVersionListRequest) (shared.DeveloperFileVersionListResponse, error)
	PrzywrocWersjePliku(ctx context.Context, z shared.DeveloperFileVersionRestoreRequest) (shared.DeveloperFileVersionRestoreResponse, error)

	// Warstwa językowa Code Editora woła programy serwera; brak wraca polem dostępności, nie odmową.
	NawigujDoSymbolu(ctx context.Context, z shared.DeveloperSymbolNavigateRequest) (shared.DeveloperSymbolNavigateResponse, error)
	Formatuj(ctx context.Context, z shared.DeveloperFormatRunRequest) (shared.DeveloperFormatRunResponse, error)
	AnalizaStatyczna(ctx context.Context, z shared.DeveloperLintGetRequest) (shared.DeveloperLintGetResponse, error)
	Refaktoryzuj(ctx context.Context, z shared.DeveloperRefactorApplyRequest) (shared.DeveloperRefactorApplyResponse, error)

	// Odczyt okna Build Output: historia przebiegów, log, wynik testów i pokrycie z rozbioru wyjścia.
	WykazBudowan(ctx context.Context, z shared.DeveloperBuildListRequest) (shared.DeveloperBuildListResponse, error)
	LogBudowania(ctx context.Context, z shared.DeveloperBuildLogGetRequest) (shared.DeveloperBuildLogGetResponse, error)
	WynikTestow(ctx context.Context, z shared.DeveloperTestResultGetRequest) (shared.DeveloperTestResultGetResponse, error)
	Pokrycie(ctx context.Context, z shared.DeveloperCoverageGetRequest) (shared.DeveloperCoverageGetResponse, error)

	// Okno Run & Debug na DAP: punkt przerwania trwały, sesja debugowania żywa razem z adapterem.
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

	// Zakładka Data Console na `database/sql`; hasło nie leży w opisie połączenia, tylko w sejfie.
	UstawPolaczenieDanych(ctx context.Context, z shared.DeveloperDataConnectionSetRequest) (shared.DeveloperDataConnectionSetResponse, error)
	WykazPolaczenDanych(ctx context.Context, z shared.DeveloperDataConnectionListRequest) (shared.DeveloperDataConnectionListResponse, error)
	SchematDanych(ctx context.Context, z shared.DeveloperDataSchemaGetRequest) (shared.DeveloperDataSchemaGetResponse, error)
	WykonajZapytanieDanych(ctx context.Context, z shared.DeveloperDataQueryRunRequest) (shared.DeveloperDataQueryRunResponse, error)
	UruchomMigracjeDanych(ctx context.Context, z shared.DeveloperDataMigrationRunRequest) (shared.DeveloperDataMigrationRunResponse, error)

	// Zakładka Containers na Docker SDK Go; brak silnika wraca polem `engineAvailable`.
	WykazKontenerow(ctx context.Context, z shared.DeveloperContainerListRequest) (shared.DeveloperContainerListResponse, error)
	CzynnoscKontenera(ctx context.Context, z shared.DeveloperContainerActionRequest) (shared.DeveloperContainerActionResponse, error)
	BudujObraz(ctx context.Context, z shared.DeveloperImageBuildRequest) (shared.DeveloperImageBuildResponse, error)
	KompozycjaKontenerow(ctx context.Context, z shared.DeveloperComposeUpRequest) (shared.DeveloperComposeUpResponse, error)

	// Zależności, bezpieczeństwo i jakość: własne parsery manifestów, bez programu spoza instalki.
	WykazZaleznosci(ctx context.Context, z shared.DeveloperDependencyListRequest) (shared.DeveloperDependencyListResponse, error)
	UruchomSkan(ctx context.Context, z shared.DeveloperScanRunRequest) (shared.DeveloperScanRunResponse, error)
	WykazZnalezisk(ctx context.Context, z shared.DeveloperScanResultListRequest) (shared.DeveloperScanResultListResponse, error)

	// Operacje kontekstowe paska pływającego oraz sonda programów warsztatu.
	OperacjaKontekstowa(ctx context.Context, z shared.DeveloperContextualOpRequest) (shared.DeveloperContextualOpResponse, error)
	SprawdzWarsztat(ctx context.Context, z shared.DeveloperToolchainCheckRequest) (shared.DeveloperToolchainCheckResponse, error)
}

// zarejestrujDevelopera wpina komplet pięćdziesięciu komend modułu Developer w rejestr komend rdzenia.
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
	// Przebieg obciążeniowy stoi przy zapytaniu pojedynczym powtórzonym pod obciążeniem.
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
	// Sesja komunikatu zostaje pusta: przebieg należy do okna, rozgłaszany też po jego zamknięciu.
	e.wyslij(shared.EventDeveloperBuildChanged, "", tresc)
}

// PodepnijPrzyrostBudowania wypełnia port: adapter zapamiętuje drogę do zdarzenia przyrostu budowania.
func (a *adapterDevelopera) PodepnijPrzyrostBudowania(
	rozglos func(shared.ChangeKind, shared.DeveloperBuild, string)) {

	a.przyrost = rozglos
}
