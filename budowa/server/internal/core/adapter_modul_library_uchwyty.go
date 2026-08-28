// Plik wpina dziesięć komend modułu Library na porcie `Biblioteka`, przez
// który rejestr komend rdzenia dociera do adaptera złożonego z dwóch plików
// tego samego typu `adapterBiblioteki`: pliku i wersji oraz kolekcji i
// etykiet.
package core

import (
	"context"

	"danacoconsole/shared"
)

// Biblioteka jest portem modułu Library, wymieniającym dziesięć komend
// dostarczanych przez adapter złożony z kilku plików tego samego typu.
type Biblioteka interface {
	Wgraj(ctx context.Context, z shared.LibraryFileUploadRequest) (shared.LibraryFileUploadResponse, error)
	Wykaz(ctx context.Context, z shared.LibraryFileListRequest) (shared.LibraryFileListResponse, error)
	Podglad(ctx context.Context, z shared.LibraryFilePreviewRequest) (shared.LibraryFilePreviewResponse, error)
	Szukaj(ctx context.Context, z shared.LibraryFileSearchRequest) (shared.LibraryFileSearchResponse, error)
	DolozWersje(ctx context.Context, z shared.LibraryVersionAddRequest) (shared.LibraryVersionAddResponse, error)
	Wersje(ctx context.Context, z shared.LibraryVersionListRequest) (shared.LibraryVersionListResponse, error)
	PrzywrocWersje(ctx context.Context, z shared.LibraryVersionRestoreRequest) (shared.LibraryVersionRestoreResponse, error)
	UstawEtykiety(ctx context.Context, z shared.LibraryTagSetRequest) (shared.LibraryTagSetResponse, error)
	UtworzKolekcje(ctx context.Context, z shared.LibraryCollectionCreateRequest) (shared.LibraryCollectionCreateResponse, error)
	PrzypiszDoKolekcji(ctx context.Context, z shared.LibraryCollectionAssignRequest) (shared.LibraryCollectionAssignResponse, error)

	// Opis zasobu i schemat metadanych (`adapter_modul_library_metadane.go`).
	Opis(ctx context.Context, z shared.LibraryMetadataGetRequest) (shared.LibraryMetadataGetResponse, error)
	ZapiszOpis(ctx context.Context, z shared.LibraryMetadataSetRequest) (shared.LibraryMetadataSetResponse, error)
	SchematMetadanych(ctx context.Context, z shared.LibrarySchemaGetRequest) (shared.LibrarySchemaGetResponse, error)
	UstawPoleSchematu(ctx context.Context, z shared.LibrarySchemaSetRequest) (shared.LibrarySchemaSetResponse, error)

	// Słownik etykiet, tezaurus i wykaz kolekcji (`adapter_modul_library_slownik.go`).
	SlownikEtykiet(ctx context.Context, z shared.LibraryTagListRequest) (shared.LibraryTagListResponse, error)
	ZmienEtykiete(ctx context.Context, z shared.LibraryTagUpdateRequest) (shared.LibraryTagUpdateResponse, error)
	PolaczEtykiety(ctx context.Context, z shared.LibraryTagMergeRequest) (shared.LibraryTagMergeResponse, error)
	UsunEtykiete(ctx context.Context, z shared.LibraryTagRemoveRequest) (shared.LibraryTagRemoveResponse, error)
	WykazKolekcji(ctx context.Context, z shared.LibraryCollectionListRequest) (shared.LibraryCollectionListResponse, error)
	UstawRelacjeTezaurusa(ctx context.Context, z shared.LibraryThesaurusRelateRequest) (shared.LibraryThesaurusRelateResponse, error)
	WywiezTezaurus(ctx context.Context, z shared.LibraryThesaurusExportRequest) (shared.LibraryThesaurusExportResponse, error)

	// Reguły repozytorium (`adapter_modul_library_reguly.go`).
	UstawRegule(ctx context.Context, z shared.LibraryRuleSetRequest) (shared.LibraryRuleSetResponse, error)
	WykazRegul(ctx context.Context, z shared.LibraryRuleListRequest) (shared.LibraryRuleListResponse, error)
	UsunRegule(ctx context.Context, z shared.LibraryRuleRemoveRequest) (shared.LibraryRuleRemoveResponse, error)

	// Higiena repozytorium (`adapter_modul_library_higiena.go`).
	SkanujDuplikaty(ctx context.Context, z shared.LibraryDuplicateScanRequest) (shared.LibraryDuplicateScanResponse, error)
	PolaczDuplikaty(ctx context.Context, z shared.LibraryDuplicateMergeRequest) (shared.LibraryDuplicateMergeResponse, error)
	SprawdzIntegralnosc(ctx context.Context, z shared.LibraryFixityCheckRequest) (shared.LibraryFixityCheckResponse, error)
	NormalizujNazwy(ctx context.Context, z shared.LibraryNameNormalizeRequest) (shared.LibraryNameNormalizeResponse, error)
	PulpitStanu(ctx context.Context, z shared.LibraryStatsGetRequest) (shared.LibraryStatsGetResponse, error)

	// Dziennik audytu (`adapter_modul_library_audyt.go`).
	DziennikAudytu(ctx context.Context, z shared.LibraryAuditListRequest) (shared.LibraryAuditListResponse, error)

	// Cykl życia zasobu (`adapter_modul_library_cykl.go`).
	PrzeniesZasoby(ctx context.Context, z shared.LibraryFileMoveRequest) (shared.LibraryFileMoveResponse, error)
	ZarchiwizujZasoby(ctx context.Context, z shared.LibraryFileArchiveRequest) (shared.LibraryFileArchiveResponse, error)
	PrzywrocZasoby(ctx context.Context, z shared.LibraryFileRestoreRequest) (shared.LibraryFileRestoreResponse, error)
	UsunZasoby(ctx context.Context, z shared.LibraryFileDeleteRequest) (shared.LibraryFileDeleteResponse, error)

	// Utrwalenie i przechowywanie (`adapter_modul_library_archiwum.go`).
	UtrwalArchiwalnie(ctx context.Context, z shared.LibraryPreservationRunRequest) (shared.LibraryPreservationRunResponse, error)
	UstawRetencje(ctx context.Context, z shared.LibraryRetentionSetRequest) (shared.LibraryRetentionSetResponse, error)
	WykazRetencji(ctx context.Context, z shared.LibraryRetentionListRequest) (shared.LibraryRetentionListResponse, error)
	WywiezPaczke(ctx context.Context, z shared.LibraryPackageExportRequest) (shared.LibraryPackageExportResponse, error)

	// Udostępnienia i nasłuchy (`adapter_modul_library_udostepnienia.go`).
	WystawUdostepnienie(ctx context.Context, z shared.LibraryShareCreateRequest) (shared.LibraryShareCreateResponse, error)
	WykazUdostepnien(ctx context.Context, z shared.LibraryShareListRequest) (shared.LibraryShareListResponse, error)
	OdwolajUdostepnienie(ctx context.Context, z shared.LibraryShareRevokeRequest) (shared.LibraryShareRevokeResponse, error)
	UstawWebhook(ctx context.Context, z shared.LibraryWebhookSetRequest) (shared.LibraryWebhookSetResponse, error)
	WykazWebhookow(ctx context.Context, z shared.LibraryWebhookListRequest) (shared.LibraryWebhookListResponse, error)
	UsunWebhook(ctx context.Context, z shared.LibraryWebhookRemoveRequest) (shared.LibraryWebhookRemoveResponse, error)

	// Klasyfikacja i sugestie (`adapter_modul_library_sugestie.go`).
	KlasyfikujWsadowo(ctx context.Context, z shared.LibraryClassifyRunRequest) (shared.LibraryClassifyRunResponse, error)
	WykazSugestii(ctx context.Context, z shared.LibrarySuggestionListRequest) (shared.LibrarySuggestionListResponse, error)
	RozstrzygnijSugestie(ctx context.Context, z shared.LibrarySuggestionApplyRequest) (shared.LibrarySuggestionApplyResponse, error)

	// Porównanie treści (`adapter_modul_library_porownanie.go`).
	PorownajTresci(ctx context.Context, z shared.LibraryDiffCompareRequest) (shared.LibraryDiffCompareResponse, error)
}

// zarejestrujBiblioteke wpina dziesięć komend modułu Library. Cztery z nich
// rozgłaszają zdarzenie `library.file.changed` po udanym wykonaniu: wgranie,
// dołożenie i przywrócenie wersji oraz ustawienie etykiet.
func zarejestrujBiblioteke(r *Rejestr, m Biblioteka, e *emiter) {
	if r == nil || m == nil {
		return
	}

	r.Zarejestruj(shared.CommandLibraryFileUpload,
		obsluz(func(ctx context.Context, z shared.LibraryFileUploadRequest) (shared.LibraryFileUploadResponse, error) {
			odpowiedz, err := m.Wgraj(ctx, z)
			if err == nil {
				e.plikBiblioteki(shared.ChangeKindCreated, odpowiedz.File)
			}
			return odpowiedz, err
		}))
	r.Zarejestruj(shared.CommandLibraryVersionAdd,
		obsluz(func(ctx context.Context, z shared.LibraryVersionAddRequest) (shared.LibraryVersionAddResponse, error) {
			odpowiedz, err := m.DolozWersje(ctx, z)
			if err == nil {
				e.plikBiblioteki(shared.ChangeKindUpdated, odpowiedz.File)
			}
			return odpowiedz, err
		}))
	r.Zarejestruj(shared.CommandLibraryFileList, obsluz(m.Wykaz))
	r.Zarejestruj(shared.CommandLibraryFilePreview, obsluz(m.Podglad))
	r.Zarejestruj(shared.CommandLibraryFileSearch, obsluz(m.Szukaj))
	r.Zarejestruj(shared.CommandLibraryVersionList, obsluz(m.Wersje))
	r.Zarejestruj(shared.CommandLibraryVersionRestore,
		obsluz(func(ctx context.Context, z shared.LibraryVersionRestoreRequest) (shared.LibraryVersionRestoreResponse, error) {
			odpowiedz, err := m.PrzywrocWersje(ctx, z)
			if err == nil {
				e.plikBiblioteki(shared.ChangeKindUpdated, odpowiedz.File)
			}
			return odpowiedz, err
		}))
	r.Zarejestruj(shared.CommandLibraryTagSet,
		obsluz(func(ctx context.Context, z shared.LibraryTagSetRequest) (shared.LibraryTagSetResponse, error) {
			odpowiedz, err := m.UstawEtykiety(ctx, z)
			if err == nil {
				e.plikBiblioteki(shared.ChangeKindUpdated, odpowiedz.File)
			}
			return odpowiedz, err
		}))
	r.Zarejestruj(shared.CommandLibraryCollectionCreate, obsluz(m.UtworzKolekcje))
	r.Zarejestruj(shared.CommandLibraryCollectionAssign, obsluz(m.PrzypiszDoKolekcji))

	// ── Opis zasobu i schemat metadanych ────────────────────────────────────
	r.Zarejestruj(shared.CommandLibraryMetadataGet, obsluz(m.Opis))
	r.Zarejestruj(shared.CommandLibraryMetadataSet,
		obsluz(func(ctx context.Context, z shared.LibraryMetadataSetRequest) (shared.LibraryMetadataSetResponse, error) {
			odpowiedz, err := m.ZapiszOpis(ctx, z)
			if err == nil {
				e.plikBiblioteki(shared.ChangeKindUpdated, odpowiedz.File)
			}
			return odpowiedz, err
		}))
	r.Zarejestruj(shared.CommandLibrarySchemaGet, obsluz(m.SchematMetadanych))
	r.Zarejestruj(shared.CommandLibrarySchemaSet, obsluz(m.UstawPoleSchematu))

	// ── Słownik etykiet, tezaurus, kolekcje ─────────────────────────────────
	r.Zarejestruj(shared.CommandLibraryTagList, obsluz(m.SlownikEtykiet))
	r.Zarejestruj(shared.CommandLibraryTagUpdate, obsluz(m.ZmienEtykiete))
	r.Zarejestruj(shared.CommandLibraryTagMerge, obsluz(m.PolaczEtykiety))
	r.Zarejestruj(shared.CommandLibraryTagRemove, obsluz(m.UsunEtykiete))
	r.Zarejestruj(shared.CommandLibraryCollectionList, obsluz(m.WykazKolekcji))
	r.Zarejestruj(shared.CommandLibraryThesaurusRelate, obsluz(m.UstawRelacjeTezaurusa))
	r.Zarejestruj(shared.CommandLibraryThesaurusExport, obsluz(m.WywiezTezaurus))

	// ── Reguły repozytorium ─────────────────────────────────────────────────
	r.Zarejestruj(shared.CommandLibraryRuleSet, obsluz(m.UstawRegule))
	r.Zarejestruj(shared.CommandLibraryRuleList, obsluz(m.WykazRegul))
	r.Zarejestruj(shared.CommandLibraryRuleRemove, obsluz(m.UsunRegule))

	// ── Higiena repozytorium ────────────────────────────────────────────────
	r.Zarejestruj(shared.CommandLibraryDuplicateScan, obsluz(m.SkanujDuplikaty))
	r.Zarejestruj(shared.CommandLibraryDuplicateMerge,
		obsluz(func(ctx context.Context, z shared.LibraryDuplicateMergeRequest) (shared.LibraryDuplicateMergeResponse, error) {
			odpowiedz, err := m.PolaczDuplikaty(ctx, z)
			if err == nil {
				e.plikBiblioteki(shared.ChangeKindUpdated, odpowiedz.File)
			}
			return odpowiedz, err
		}))
	r.Zarejestruj(shared.CommandLibraryFixityCheck, obsluz(m.SprawdzIntegralnosc))
	r.Zarejestruj(shared.CommandLibraryNameNormalize, obsluz(m.NormalizujNazwy))
	r.Zarejestruj(shared.CommandLibraryStatsGet, obsluz(m.PulpitStanu))
	r.Zarejestruj(shared.CommandLibraryAuditList, obsluz(m.DziennikAudytu))

	// ── Cykl życia zasobu ───────────────────────────────────────────────────

	// Rozgłaszają `library.file.changed` osobno dla każdego zasobu.
	r.Zarejestruj(shared.CommandLibraryFileMove,
		obsluz(func(ctx context.Context, z shared.LibraryFileMoveRequest) (shared.LibraryFileMoveResponse, error) {
			odpowiedz, err := m.PrzeniesZasoby(ctx, z)
			if err == nil {
				rozglosPlikiBiblioteki(e, shared.ChangeKindUpdated, odpowiedz.Files)
			}
			return odpowiedz, err
		}))
	r.Zarejestruj(shared.CommandLibraryFileArchive,
		obsluz(func(ctx context.Context, z shared.LibraryFileArchiveRequest) (shared.LibraryFileArchiveResponse, error) {
			odpowiedz, err := m.ZarchiwizujZasoby(ctx, z)
			if err == nil {
				rozglosPlikiBiblioteki(e, shared.ChangeKindUpdated, odpowiedz.Files)
			}
			return odpowiedz, err
		}))
	r.Zarejestruj(shared.CommandLibraryFileRestore,
		obsluz(func(ctx context.Context, z shared.LibraryFileRestoreRequest) (shared.LibraryFileRestoreResponse, error) {
			odpowiedz, err := m.PrzywrocZasoby(ctx, z)
			if err == nil {
				rozglosPlikiBiblioteki(e, shared.ChangeKindUpdated, odpowiedz.Files)
			}
			return odpowiedz, err
		}))
	r.Zarejestruj(shared.CommandLibraryFileDelete, obsluz(m.UsunZasoby))

	// ── Utrwalenie, retencja, paczka ────────────────────────────────────────
	r.Zarejestruj(shared.CommandLibraryPreservationRun, obsluz(m.UtrwalArchiwalnie))
	r.Zarejestruj(shared.CommandLibraryRetentionSet, obsluz(m.UstawRetencje))
	r.Zarejestruj(shared.CommandLibraryRetentionList, obsluz(m.WykazRetencji))
	r.Zarejestruj(shared.CommandLibraryPackageExport, obsluz(m.WywiezPaczke))

	// ── Udostępnienia i nasłuchy ────────────────────────────────────────────
	r.Zarejestruj(shared.CommandLibraryShareCreate, obsluz(m.WystawUdostepnienie))
	r.Zarejestruj(shared.CommandLibraryShareList, obsluz(m.WykazUdostepnien))
	r.Zarejestruj(shared.CommandLibraryShareRevoke, obsluz(m.OdwolajUdostepnienie))
	r.Zarejestruj(shared.CommandLibraryWebhookSet, obsluz(m.UstawWebhook))
	r.Zarejestruj(shared.CommandLibraryWebhookList, obsluz(m.WykazWebhookow))
	r.Zarejestruj(shared.CommandLibraryWebhookRemove, obsluz(m.UsunWebhook))

	// ── Klasyfikacja, sugestie, porównanie ──────────────────────────────────
	r.Zarejestruj(shared.CommandLibraryClassifyRun, obsluz(m.KlasyfikujWsadowo))
	r.Zarejestruj(shared.CommandLibrarySuggestionList, obsluz(m.WykazSugestii))
	r.Zarejestruj(shared.CommandLibrarySuggestionApply,
		obsluz(func(ctx context.Context, z shared.LibrarySuggestionApplyRequest) (shared.LibrarySuggestionApplyResponse, error) {
			odpowiedz, err := m.RozstrzygnijSugestie(ctx, z)
			if err == nil {
				rozglosPlikiBiblioteki(e, shared.ChangeKindUpdated, odpowiedz.Files)
			}
			return odpowiedz, err
		}))
	r.Zarejestruj(shared.CommandLibraryDiffCompare, obsluz(m.PorownajTresci))
}

// rozglosPlikiBiblioteki rozgłasza zmianę wielu zasobów naraz — komendy
// zbiorowe zmieniają wykaz w kilku miejscach jednocześnie, a okno ma odświeżyć
// każde z nich.
func rozglosPlikiBiblioteki(e *emiter, zmiana shared.ChangeKind, pliki []shared.LibraryFile) {
	for _, plik := range pliki {
		e.plikBiblioteki(zmiana, plik)
	}
}

// plikBiblioteki rozgłasza `library.file.changed`, plik repozytorium po
// zmianie. Plik nie jest bytem karty sesji, więc zdarzenie idzie bez jej
// wskazania, a Library Explorer odświeża się ze strony głównej.
func (e *emiter) plikBiblioteki(zmiana shared.ChangeKind, plik shared.LibraryFile) {
	e.wyslij(shared.EventLibraryFileChanged, "",
		shared.LibraryFileChangedEvent{Change: zmiana, File: plik})
}
