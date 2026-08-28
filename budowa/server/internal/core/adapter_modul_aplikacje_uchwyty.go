// Plik wpina komendy portu Aplikacje modulu Apps do rejestru rdzenia; adapter laczy
// adapter_modul_aplikacje.go, adapter_modul_aplikacje_wdrozenie.go i
// adapter_modul_aplikacje_odczyt.go w jeden typ adapterAplikacji.
package core

import (
	"context"

	"danacoconsole/shared"
)

// Aplikacje jest portem modulu Apps: zbiorem metod, przez ktore rejestr komend
// rdzenia dociera do adaptera obslugujacego produkt, architekture, warsztat i wdrozenia.
type Aplikacje interface {
	ZdefiniujArchitekture(ctx context.Context, z shared.AppsArchitectureDefineRequest) (shared.AppsArchitectureDefineResponse, error)
	ZaktualizujPrzestrzen(ctx context.Context, z shared.AppsWorkspaceUpdateRequest) (shared.AppsWorkspaceUpdateResponse, error)
	UruchomWdrozenie(ctx context.Context, z shared.AppsDeploymentRunRequest) (shared.AppsDeploymentRunResponse, error)
	// Trzy komendy odczytu daja droge powrotna do wierszy zapisanych trzema komendami wyzej.
	WypiszWdrozenia(ctx context.Context, z shared.AppsDeploymentListRequest) (shared.AppsDeploymentListResponse, error)
	PobierzArchitekture(ctx context.Context, z shared.AppsArchitectureGetRequest) (shared.AppsArchitectureGetResponse, error)
	WypiszPlikiWarsztatu(ctx context.Context, z shared.AppsWorkspaceListRequest) (shared.AppsWorkspaceListResponse, error)
	// PodepnijPrzyrostWdrozenia oddaje adapterowi droge do zdarzenia zmiany wdrozenia.
	PodepnijPrzyrostWdrozenia(rozglos func(shared.ChangeKind, shared.AppDeployment))
	// PodepnijPrzyrostWarsztatu oddaje adapterowi droge do zdarzenia apps.workspace.changed.
	PodepnijPrzyrostWarsztatu(rozglos func(shared.ChangeKind, string, shared.AppWorkspaceLayer, shared.DeveloperFile))
	// PodepnijPrzyrostEtapu oddaje adapterowi druga droge do zdarzenia apps.build.changed.
	PodepnijPrzyrostEtapu(rozglos func(shared.ChangeKind, shared.AppStage))

	// --- Product Builder: produkt, etapy, kamienie milowe, oś czasu ---
	PobierzProdukt(ctx context.Context, z shared.AppsProductGetRequest) (shared.AppsProductGetResponse, error)
	ZapiszProdukt(ctx context.Context, z shared.AppsProductSaveRequest) (shared.AppsProductSaveResponse, error)
	WypiszPowiazaniaProduktu(ctx context.Context, z shared.AppsProductLinkListRequest) (shared.AppsProductLinkListResponse, error)
	WypiszEtapy(ctx context.Context, z shared.AppsStageListRequest) (shared.AppsStageListResponse, error)
	ZapiszEtap(ctx context.Context, z shared.AppsStageSaveRequest) (shared.AppsStageSaveResponse, error)
	WypiszKamienieMilowe(ctx context.Context, z shared.AppsMilestoneListRequest) (shared.AppsMilestoneListResponse, error)
	ZapiszKamienMilowy(ctx context.Context, z shared.AppsMilestoneSaveRequest) (shared.AppsMilestoneSaveResponse, error)
	UsunKamienMilowy(ctx context.Context, z shared.AppsMilestoneDeleteRequest) (shared.AppsMilestoneDeleteResponse, error)
	WypiszOsCzasu(ctx context.Context, z shared.AppsTimelineListRequest) (shared.AppsTimelineListResponse, error)

	// --- Architecture Designer: walidacja, wersje, adnotacje, eksport ---
	SprawdzArchitekture(ctx context.Context, z shared.AppsArchitectureValidateRequest) (shared.AppsArchitectureValidateResponse, error)
	WypiszWersjeArchitektury(ctx context.Context, z shared.AppsArchitectureVersionListRequest) (shared.AppsArchitectureVersionListResponse, error)
	ZapiszAdnotacje(ctx context.Context, z shared.AppsArchitectureAnnotationSaveRequest) (shared.AppsArchitectureAnnotationSaveResponse, error)
	WyeksportujArchitekture(ctx context.Context, z shared.AppsArchitectureExportRequest) (shared.AppsArchitectureExportResponse, error)

	// --- Frontend i Backend Workspace: podgląd, trasy, motyw, punkty końcowe, schemat ---
	UruchomPodglad(ctx context.Context, z shared.AppsPreviewStartRequest) (shared.AppsPreviewStartResponse, error)
	ZatrzymajPodglad(ctx context.Context, z shared.AppsPreviewStopRequest) (shared.AppsPreviewStopResponse, error)
	WypiszTrasy(ctx context.Context, z shared.AppsRouteListRequest) (shared.AppsRouteListResponse, error)
	PobierzMotyw(ctx context.Context, z shared.AppsThemeGetRequest) (shared.AppsThemeGetResponse, error)
	UstawMotyw(ctx context.Context, z shared.AppsThemeSetRequest) (shared.AppsThemeSetResponse, error)
	WypiszPunktyKoncowe(ctx context.Context, z shared.AppsEndpointListRequest) (shared.AppsEndpointListResponse, error)
	ZapytajPunktKoncowy(ctx context.Context, z shared.AppsEndpointProbeRequest) (shared.AppsEndpointProbeResponse, error)
	PobierzSchemat(ctx context.Context, z shared.AppsSchemaGetRequest) (shared.AppsSchemaGetResponse, error)

	// --- Deployment Panel: środowiska, zmienne, domena, skalowanie, kondycja ---
	WypiszSrodowiska(ctx context.Context, z shared.AppsEnvironmentListRequest) (shared.AppsEnvironmentListResponse, error)
	WypiszZmienneSrodowiska(ctx context.Context, z shared.AppsEnvironmentVariableListRequest) (shared.AppsEnvironmentVariableListResponse, error)
	UstawZmiennaSrodowiska(ctx context.Context, z shared.AppsEnvironmentVariableSetRequest) (shared.AppsEnvironmentVariableSetResponse, error)
	UstawDomene(ctx context.Context, z shared.AppsDeploymentDomainSetRequest) (shared.AppsDeploymentDomainSetResponse, error)
	UstawSkalowanie(ctx context.Context, z shared.AppsDeploymentScaleSetRequest) (shared.AppsDeploymentScaleSetResponse, error)
	PobierzKondycje(ctx context.Context, z shared.AppsDeploymentHealthGetRequest) (shared.AppsDeploymentHealthGetResponse, error)
	ZmierzWydajnosc(ctx context.Context, z shared.AppsPerformanceAuditRequest) (shared.AppsPerformanceAuditResponse, error)

	// --- Dzienniki i artefakty ---
	OdczytajDziennikUslugi(ctx context.Context, z shared.AppsServiceLogReadRequest) (shared.AppsServiceLogReadResponse, error)
	OdczytajDziennikWdrozenia(ctx context.Context, z shared.AppsDeploymentLogReadRequest) (shared.AppsDeploymentLogReadResponse, error)
	WypiszArtefakty(ctx context.Context, z shared.AppsArtifactListRequest) (shared.AppsArtifactListResponse, error)

	// --- Publisher Panel: pakowanie, manifest, walidacja, podpis, publikacja ---
	ZbudujPakiet(ctx context.Context, z shared.AppsPackageBuildRequest) (shared.AppsPackageBuildResponse, error)
	ZapiszManifestPakietu(ctx context.Context, z shared.AppsPackageManifestSaveRequest) (shared.AppsPackageManifestSaveResponse, error)
	SprawdzPakiet(ctx context.Context, z shared.AppsPackageValidateRequest) (shared.AppsPackageValidateResponse, error)
	PodpiszPakiet(ctx context.Context, z shared.AppsPackageSignRequest) (shared.AppsPackageSignResponse, error)
	OpublikujPakiet(ctx context.Context, z shared.AppsPackagePublishRequest) (shared.AppsPackagePublishResponse, error)
}

// Adapter wypełnia port w całości. Gdyby port i adapter się rozjechały,
// kompilacja stanie tutaj, a nie dopiero na martwej komendzie u Operatora.
var _ Aplikacje = (*adapterAplikacji)(nil)

// zarejestrujAplikacje wpina komplet komend rodziny `apps.*` — całą stronę
// budowy produktu, od Product Buildera po Publisher Panel.
func zarejestrujAplikacje(r *Rejestr, m Aplikacje, e *emiter) {
	if r == nil || m == nil {
		return
	}
	m.PodepnijPrzyrostWdrozenia(e.wdrozenieApp)
	m.PodepnijPrzyrostWarsztatu(e.warsztatApp)
	m.PodepnijPrzyrostEtapu(e.etapApp)

	r.Zarejestruj(shared.CommandAppsArchitectureDefine, obsluz(m.ZdefiniujArchitekture))
	r.Zarejestruj(shared.CommandAppsWorkspaceUpdate, obsluz(m.ZaktualizujPrzestrzen))
	r.Zarejestruj(shared.CommandAppsDeploymentRun, obsluz(m.UruchomWdrozenie))
	r.Zarejestruj(shared.CommandAppsDeploymentList, obsluz(m.WypiszWdrozenia))
	r.Zarejestruj(shared.CommandAppsArchitectureGet, obsluz(m.PobierzArchitekture))
	r.Zarejestruj(shared.CommandAppsWorkspaceList, obsluz(m.WypiszPlikiWarsztatu))

	r.Zarejestruj(shared.CommandAppsProductGet, obsluz(m.PobierzProdukt))
	r.Zarejestruj(shared.CommandAppsProductSave, obsluz(m.ZapiszProdukt))
	r.Zarejestruj(shared.CommandAppsProductLinkList, obsluz(m.WypiszPowiazaniaProduktu))
	r.Zarejestruj(shared.CommandAppsStageList, obsluz(m.WypiszEtapy))
	r.Zarejestruj(shared.CommandAppsStageSave, obsluz(m.ZapiszEtap))
	r.Zarejestruj(shared.CommandAppsMilestoneList, obsluz(m.WypiszKamienieMilowe))
	r.Zarejestruj(shared.CommandAppsMilestoneSave, obsluz(m.ZapiszKamienMilowy))
	r.Zarejestruj(shared.CommandAppsMilestoneDelete, obsluz(m.UsunKamienMilowy))
	r.Zarejestruj(shared.CommandAppsTimelineList, obsluz(m.WypiszOsCzasu))

	r.Zarejestruj(shared.CommandAppsArchitectureValidate, obsluz(m.SprawdzArchitekture))
	r.Zarejestruj(shared.CommandAppsArchitectureVersionList, obsluz(m.WypiszWersjeArchitektury))
	r.Zarejestruj(shared.CommandAppsArchitectureAnnotationSave, obsluz(m.ZapiszAdnotacje))
	r.Zarejestruj(shared.CommandAppsArchitectureExport, obsluz(m.WyeksportujArchitekture))

	r.Zarejestruj(shared.CommandAppsPreviewStart, obsluz(m.UruchomPodglad))
	r.Zarejestruj(shared.CommandAppsPreviewStop, obsluz(m.ZatrzymajPodglad))
	r.Zarejestruj(shared.CommandAppsRouteList, obsluz(m.WypiszTrasy))
	r.Zarejestruj(shared.CommandAppsThemeGet, obsluz(m.PobierzMotyw))
	r.Zarejestruj(shared.CommandAppsThemeSet, obsluz(m.UstawMotyw))
	r.Zarejestruj(shared.CommandAppsEndpointList, obsluz(m.WypiszPunktyKoncowe))
	r.Zarejestruj(shared.CommandAppsEndpointProbe, obsluz(m.ZapytajPunktKoncowy))
	r.Zarejestruj(shared.CommandAppsSchemaGet, obsluz(m.PobierzSchemat))

	r.Zarejestruj(shared.CommandAppsEnvironmentList, obsluz(m.WypiszSrodowiska))
	r.Zarejestruj(shared.CommandAppsEnvironmentVariableList, obsluz(m.WypiszZmienneSrodowiska))
	r.Zarejestruj(shared.CommandAppsEnvironmentVariableSet, obsluz(m.UstawZmiennaSrodowiska))
	r.Zarejestruj(shared.CommandAppsDeploymentDomainSet, obsluz(m.UstawDomene))
	r.Zarejestruj(shared.CommandAppsDeploymentScaleSet, obsluz(m.UstawSkalowanie))
	r.Zarejestruj(shared.CommandAppsDeploymentHealthGet, obsluz(m.PobierzKondycje))
	// Audyt wydajnosci zdarzenia nie rozglasza: mierzy strone i niczego w niej nie zmienia.
	r.Zarejestruj(shared.CommandAppsPerformanceAudit, obsluz(m.ZmierzWydajnosc))

	r.Zarejestruj(shared.CommandAppsServiceLogRead, obsluz(m.OdczytajDziennikUslugi))
	r.Zarejestruj(shared.CommandAppsDeploymentLogRead, obsluz(m.OdczytajDziennikWdrozenia))
	r.Zarejestruj(shared.CommandAppsArtifactList, obsluz(m.WypiszArtefakty))

	r.Zarejestruj(shared.CommandAppsPackageBuild, obsluz(m.ZbudujPakiet))
	r.Zarejestruj(shared.CommandAppsPackageManifestSave, obsluz(m.ZapiszManifestPakietu))
	r.Zarejestruj(shared.CommandAppsPackageValidate, obsluz(m.SprawdzPakiet))
	r.Zarejestruj(shared.CommandAppsPackageSign, obsluz(m.PodpiszPakiet))
	r.Zarejestruj(shared.CommandAppsPackagePublish, obsluz(m.OpublikujPakiet))
}

// etapApp rozgłasza `apps.build.changed` przy zmianie etapu budowy. Pole
// `Deployment` zostaje puste: zmiana dotyczy osi etapów Product Buildera,
// a przebiegu wdrożenia w niej nie ma.
func (e *emiter) etapApp(zmiana shared.ChangeKind, etap shared.AppStage) {
	e.wyslij(shared.EventAppsBuildChanged, "", shared.AppsBuildChangedEvent{
		Change: zmiana,
		Stage:  etap,
	})
}

// warsztatApp rozglasza zdarzenie apps.workspace.changed po zapisie pliku warsztatu;
// sesja komunikatu zostaje pusta, bo warsztat nalezy do okna, a nie do sesji rdzenia.
func (e *emiter) warsztatApp(zmiana shared.ChangeKind, okno string,
	warstwa shared.AppWorkspaceLayer, plik shared.DeveloperFile) {

	e.wyslij(shared.EventAppsWorkspaceChanged, "", shared.AppsWorkspaceChangedEvent{
		Change:   zmiana,
		WindowId: okno,
		Layer:    warstwa,
		File:     &plik,
	})
}

// wdrozenieApp rozglasza zdarzenie apps.build.changed przy zmianie stanu przebiegu
// wdrozenia; pole Stage niesie okno zmiany, a pole Deployment sam przebieg.
func (e *emiter) wdrozenieApp(zmiana shared.ChangeKind, wdrozenie shared.AppDeployment) {
	e.wyslij(shared.EventAppsBuildChanged, "", shared.AppsBuildChangedEvent{
		Change:     zmiana,
		Stage:      shared.AppStage{WindowId: wdrozenie.WindowId},
		Deployment: &wdrozenie,
	})
}
