// Plik wpina pięć komend rodziny extension.* — katalogu rozszerzeń — na jeden port i jedną maszynerię; adapter i rozstrzygnięcie, czym rozszerzenie jest, leżą w adapter_modul_extension.go.
package core

import (
	"context"

	"danacoconsole/shared"
)

// Rozszerzenia jest portem rodziny `extension.*`. Mówi wyłącznie
// typami kontraktu; wiersze tabel `rozszerzenie` i `punkt_dostepu` leżą po
// drugiej stronie adaptera.
type Rozszerzenia interface {
	// Wykaz obsługuje `extension.list`.
	Wykaz(ctx context.Context, z shared.ExtensionListRequest) (shared.ExtensionListResponse, error)
	// Zainstaluj obsługuje `extension.install`.
	Zainstaluj(ctx context.Context, z shared.ExtensionInstallRequest) (shared.ExtensionInstallResponse, error)
	// Skonfiguruj obsługuje `extension.configure`.
	Skonfiguruj(ctx context.Context, z shared.ExtensionConfigureRequest) (shared.ExtensionConfigureResponse, error)
	// Przestaw obsługuje `extension.toggle`.
	Przestaw(ctx context.Context, z shared.ExtensionToggleRequest) (shared.ExtensionToggleResponse, error)
	// Odinstaluj obsługuje `extension.uninstall`.
	Odinstaluj(ctx context.Context, z shared.ExtensionUninstallRequest) (shared.ExtensionUninstallResponse, error)

	// --- App Catalog i Installed Apps Manager (adapter_modul_extension_katalog.go) ---
	Szukaj(ctx context.Context, z shared.ExtensionSearchRequest) (shared.ExtensionSearchResponse, error)
	PobierzSzczegol(ctx context.Context, z shared.ExtensionDetailGetRequest) (shared.ExtensionDetailGetResponse, error)
	WypiszKolekcje(ctx context.Context, z shared.ExtensionCollectionListRequest) (shared.ExtensionCollectionListResponse, error)
	ZapiszKolekcje(ctx context.Context, z shared.ExtensionCollectionSaveRequest) (shared.ExtensionCollectionSaveResponse, error)
	ZastosujKolekcje(ctx context.Context, z shared.ExtensionCollectionApplyRequest) (shared.ExtensionCollectionApplyResponse, error)
	WypiszRejestr(ctx context.Context, z shared.ExtensionRegistryListRequest) (shared.ExtensionRegistryListResponse, error)
	SprawdzAktualizacje(ctx context.Context, z shared.ExtensionUpdateCheckRequest) (shared.ExtensionUpdateCheckResponse, error)
	PrzeslijPaczke(ctx context.Context, z shared.ExtensionPackageUploadRequest) (shared.ExtensionPackageUploadResponse, error)
	PrzypnijWersje(ctx context.Context, z shared.ExtensionVersionPinRequest) (shared.ExtensionVersionPinResponse, error)
	CofnijWersje(ctx context.Context, z shared.ExtensionVersionRollbackRequest) (shared.ExtensionVersionRollbackResponse, error)
	ZainstalujZestaw(ctx context.Context, z shared.ExtensionBundleInstallRequest) (shared.ExtensionBundleInstallResponse, error)
	WypiszHistorie(ctx context.Context, z shared.ExtensionHistoryListRequest) (shared.ExtensionHistoryListResponse, error)
	WykonajZbiorczo(ctx context.Context, z shared.ExtensionAdminBulkRequest) (shared.ExtensionAdminBulkResponse, error)

	// --- Integrations Hub i MCP & Connector Console (adapter_modul_extension_integracje.go) ---
	UstawTransport(ctx context.Context, z shared.ExtensionTransportSetRequest) (shared.ExtensionTransportSetResponse, error)
	PowiazPoswiadczenie(ctx context.Context, z shared.ExtensionCredentialBindRequest) (shared.ExtensionCredentialBindResponse, error)
	WypiszNarzedzia(ctx context.Context, z shared.ExtensionToolListRequest) (shared.ExtensionToolListResponse, error)
	WywolajNarzedzie(ctx context.Context, z shared.ExtensionToolCallRequest) (shared.ExtensionToolCallResponse, error)
	WypiszLogProtokolu(ctx context.Context, z shared.ExtensionProtocolLogListRequest) (shared.ExtensionProtocolLogListResponse, error)
	UruchomPiaskownice(ctx context.Context, z shared.ExtensionSandboxRunRequest) (shared.ExtensionSandboxRunResponse, error)
	ZaimportujDefinicje(ctx context.Context, z shared.ExtensionDefinitionImportRequest) (shared.ExtensionDefinitionImportResponse, error)
	WypiszWebhooki(ctx context.Context, z shared.ExtensionWebhookListRequest) (shared.ExtensionWebhookListResponse, error)
	ZapiszWebhook(ctx context.Context, z shared.ExtensionWebhookSaveRequest) (shared.ExtensionWebhookSaveResponse, error)
	ZapiszMapowanie(ctx context.Context, z shared.ExtensionMappingSaveRequest) (shared.ExtensionMappingSaveResponse, error)
	PobierzUzycie(ctx context.Context, z shared.ExtensionUsageGetRequest) (shared.ExtensionUsageGetResponse, error)
	SprawdzKondycje(ctx context.Context, z shared.ExtensionHealthCheckRequest) (shared.ExtensionHealthCheckResponse, error)
	WypiszAudyt(ctx context.Context, z shared.ExtensionAuditListRequest) (shared.ExtensionAuditListResponse, error)

	// --- Permissions & Trust Center (adapter_modul_extension_zaufanie.go) ---
	WypiszUprawnienia(ctx context.Context, z shared.ExtensionPermissionListRequest) (shared.ExtensionPermissionListResponse, error)
	NadajUprawnienia(ctx context.Context, z shared.ExtensionPermissionGrantRequest) (shared.ExtensionPermissionGrantResponse, error)
	ZweryfikujPodpis(ctx context.Context, z shared.ExtensionSignatureVerifyRequest) (shared.ExtensionSignatureVerifyResponse, error)
	SkanujManifest(ctx context.Context, z shared.ExtensionManifestScanRequest) (shared.ExtensionManifestScanResponse, error)
	WypiszSekrety(ctx context.Context, z shared.ExtensionSecretListRequest) (shared.ExtensionSecretListResponse, error)
	UdostepnijSekret(ctx context.Context, z shared.ExtensionSecretShareRequest) (shared.ExtensionSecretShareResponse, error)
}

// Adapter wypełnia port w całości. Gdyby port i adapter się rozjechały,
// kompilacja stanie tutaj, a nie dopiero na martwej komendzie u Operatora.
var _ Rozszerzenia = (*adapterRozszerzen)(nil)

// zarejestrujRozszerzenia wpina komplet komend rodziny `extension.*` — całą
// stronę dystrybucji i konsumpcji rozszerzeń modułu Apps.
func zarejestrujRozszerzenia(r *Rejestr, x Rozszerzenia) {
	if r == nil || x == nil {
		return
	}

	r.Zarejestruj(shared.CommandExtensionList, obsluz(x.Wykaz))
	r.Zarejestruj(shared.CommandExtensionInstall, obsluz(x.Zainstaluj))
	r.Zarejestruj(shared.CommandExtensionConfigure, obsluz(x.Skonfiguruj))
	r.Zarejestruj(shared.CommandExtensionToggle, obsluz(x.Przestaw))
	r.Zarejestruj(shared.CommandExtensionUninstall, obsluz(x.Odinstaluj))

	r.Zarejestruj(shared.CommandExtensionSearch, obsluz(x.Szukaj))
	r.Zarejestruj(shared.CommandExtensionDetailGet, obsluz(x.PobierzSzczegol))
	r.Zarejestruj(shared.CommandExtensionCollectionList, obsluz(x.WypiszKolekcje))
	r.Zarejestruj(shared.CommandExtensionCollectionSave, obsluz(x.ZapiszKolekcje))
	r.Zarejestruj(shared.CommandExtensionCollectionApply, obsluz(x.ZastosujKolekcje))
	r.Zarejestruj(shared.CommandExtensionRegistryList, obsluz(x.WypiszRejestr))
	r.Zarejestruj(shared.CommandExtensionUpdateCheck, obsluz(x.SprawdzAktualizacje))
	r.Zarejestruj(shared.CommandExtensionPackageUpload, obsluz(x.PrzeslijPaczke))
	r.Zarejestruj(shared.CommandExtensionVersionPin, obsluz(x.PrzypnijWersje))
	r.Zarejestruj(shared.CommandExtensionVersionRollback, obsluz(x.CofnijWersje))
	r.Zarejestruj(shared.CommandExtensionBundleInstall, obsluz(x.ZainstalujZestaw))
	r.Zarejestruj(shared.CommandExtensionHistoryList, obsluz(x.WypiszHistorie))
	r.Zarejestruj(shared.CommandExtensionAdminBulk, obsluz(x.WykonajZbiorczo))

	r.Zarejestruj(shared.CommandExtensionTransportSet, obsluz(x.UstawTransport))
	r.Zarejestruj(shared.CommandExtensionCredentialBind, obsluz(x.PowiazPoswiadczenie))
	r.Zarejestruj(shared.CommandExtensionToolList, obsluz(x.WypiszNarzedzia))
	r.Zarejestruj(shared.CommandExtensionToolCall, obsluz(x.WywolajNarzedzie))
	r.Zarejestruj(shared.CommandExtensionProtocolLogList, obsluz(x.WypiszLogProtokolu))
	r.Zarejestruj(shared.CommandExtensionSandboxRun, obsluz(x.UruchomPiaskownice))
	r.Zarejestruj(shared.CommandExtensionDefinitionImport, obsluz(x.ZaimportujDefinicje))
	r.Zarejestruj(shared.CommandExtensionWebhookList, obsluz(x.WypiszWebhooki))
	r.Zarejestruj(shared.CommandExtensionWebhookSave, obsluz(x.ZapiszWebhook))
	r.Zarejestruj(shared.CommandExtensionMappingSave, obsluz(x.ZapiszMapowanie))
	r.Zarejestruj(shared.CommandExtensionUsageGet, obsluz(x.PobierzUzycie))
	r.Zarejestruj(shared.CommandExtensionHealthCheck, obsluz(x.SprawdzKondycje))
	r.Zarejestruj(shared.CommandExtensionAuditList, obsluz(x.WypiszAudyt))

	r.Zarejestruj(shared.CommandExtensionPermissionList, obsluz(x.WypiszUprawnienia))
	r.Zarejestruj(shared.CommandExtensionPermissionGrant, obsluz(x.NadajUprawnienia))
	r.Zarejestruj(shared.CommandExtensionSignatureVerify, obsluz(x.ZweryfikujPodpis))
	r.Zarejestruj(shared.CommandExtensionManifestScan, obsluz(x.SkanujManifest))
	r.Zarejestruj(shared.CommandExtensionSecretList, obsluz(x.WypiszSekrety))
	r.Zarejestruj(shared.CommandExtensionSecretShare, obsluz(x.UdostepnijSekret))
}
