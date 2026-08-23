// Odpowiedzialność pliku: wpięcie WSZYSTKICH komend modułu Research — portu
// `Badania`, przez który rejestr komend rdzenia dociera do adaptera złożonego
// z dziesięciu plików tego samego typu `*adapterBadan`.
//
// Port wymienia komplet komend niezależnie od tego, który plik adaptera je
// wypełnia: rejestr rdzenia potrzebuje jednego miejsca wiążącego nazwę komendy
// z metodą portu — tak samo jak `adapter_modul_library_uchwyty.go` rejestruje
// w jednej funkcji komendy swojego modułu.
//
// ── Które komendy rozgłaszają zdarzenie ────────────────────────────────────
// Kontrakt daje modułowi cztery zdarzenia: `research.report.changed`,
// `research.source.changed`, `research.finding.changed` i
// `research.monitor.changed`. Rozgłaszają wyłącznie te działania, które któreś
// z nich naprawdę opisuje — wzór z Automatyk. Działanie bez zdarzenia
// w kontrakcie nie wymyśla sobie własnego.
package core

import (
	"context"

	"danacoconsole/shared"
)

// Badania jest portem modułu Research.
type Badania interface {
	// źródła
	DodajZrodlo(ctx context.Context, z shared.ResearchSourceAddRequest) (shared.ResearchSourceAddResponse, error)
	WypiszZrodla(ctx context.Context, z shared.ResearchSourceListRequest) (shared.ResearchSourceListResponse, error)
	ZmienZrodlo(ctx context.Context, z shared.ResearchSourceUpdateRequest) (shared.ResearchSourceUpdateResponse, error)
	UsunZrodlo(ctx context.Context, z shared.ResearchSourceRemoveRequest) (shared.ResearchSourceRemoveResponse, error)
	ScalZrodla(ctx context.Context, z shared.ResearchSourceMergeRequest) (shared.ResearchSourceMergeResponse, error)
	SzukajDuplikatowZrodel(ctx context.Context, z shared.ResearchSourceDuplicatesRequest) (shared.ResearchSourceDuplicatesResponse, error)
	OznaczZrodlo(ctx context.Context, z shared.ResearchSourceTagRequest) (shared.ResearchSourceTagResponse, error)
	DodajZalacznikZrodla(ctx context.Context, z shared.ResearchSourceAttachmentAddRequest) (shared.ResearchSourceAttachmentAddResponse, error)
	WypiszZalacznikiZrodla(ctx context.Context, z shared.ResearchSourceAttachmentListRequest) (shared.ResearchSourceAttachmentListResponse, error)
	WczytajBibliografie(ctx context.Context, z shared.ResearchSourceImportRequest) (shared.ResearchSourceImportResponse, error)
	PrzechwycStrone(ctx context.Context, z shared.ResearchSourceCaptureRequest) (shared.ResearchSourceCaptureResponse, error)
	PrzepiszNagranie(ctx context.Context, z shared.ResearchSourceTranscribeRequest) (shared.ResearchSourceTranscribeResponse, error)
	RozstrzygnijIdentyfikator(ctx context.Context, z shared.ResearchSourceResolveRequest) (shared.ResearchSourceResolveResponse, error)

	// odkrywanie, monitory, import wsadowy
	SzukajZrodel(ctx context.Context, z shared.ResearchDiscoverySearchRequest) (shared.ResearchDiscoverySearchResponse, error)
	UlozZapytania(ctx context.Context, z shared.ResearchDiscoveryAssistRequest) (shared.ResearchDiscoveryAssistResponse, error)
	RozwinCytowania(ctx context.Context, z shared.ResearchDiscoverySnowballRequest) (shared.ResearchDiscoverySnowballResponse, error)
	OdrzucWyniki(ctx context.Context, z shared.ResearchDiscoveryRejectRequest) (shared.ResearchDiscoveryRejectResponse, error)
	UstawMonitor(ctx context.Context, z shared.ResearchMonitorSetRequest) (shared.ResearchMonitorSetResponse, error)
	WypiszMonitory(ctx context.Context, z shared.ResearchMonitorListRequest) (shared.ResearchMonitorListResponse, error)
	OdswiezMonitory(ctx context.Context, z shared.ResearchMonitorRefreshRequest) (shared.ResearchMonitorRefreshResponse, error)
	WczytajPartieAdresow(ctx context.Context, z shared.ResearchBatchImportRequest) (shared.ResearchBatchImportResponse, error)

	// lektura i ekstrakcja
	OtworzDoLektury(ctx context.Context, z shared.ResearchReadingOpenRequest) (shared.ResearchReadingOpenResponse, error)
	DodajAdnotacje(ctx context.Context, z shared.ResearchAnnotationAddRequest) (shared.ResearchAnnotationAddResponse, error)
	WypiszAdnotacje(ctx context.Context, z shared.ResearchAnnotationListRequest) (shared.ResearchAnnotationListResponse, error)
	UsunAdnotacje(ctx context.Context, z shared.ResearchAnnotationRemoveRequest) (shared.ResearchAnnotationRemoveResponse, error)
	WypiszWypisy(ctx context.Context, z shared.ResearchExcerptListRequest) (shared.ResearchExcerptListResponse, error)
	StreszczZrodlo(ctx context.Context, z shared.ResearchSourceSummarizeRequest) (shared.ResearchSourceSummarizeResponse, error)
	WyodrebnijTabele(ctx context.Context, z shared.ResearchSourceExtractTableRequest) (shared.ResearchSourceExtractTableResponse, error)
	WyodrebnijTwierdzenia(ctx context.Context, z shared.ResearchSourceExtractClaimsRequest) (shared.ResearchSourceExtractClaimsResponse, error)
	RozpoznajPismoZrodla(ctx context.Context, z shared.ResearchSourceOcrRequest) (shared.ResearchSourceOcrResponse, error)
	ZapytajKorpus(ctx context.Context, z shared.ResearchCorpusAskRequest) (shared.ResearchCorpusAskResponse, error)

	// ustalenia i analiza
	DodajUstalenie(ctx context.Context, z shared.ResearchFindingAddRequest) (shared.ResearchFindingAddResponse, error)
	WypiszUstalenia(ctx context.Context, z shared.ResearchFindingListRequest) (shared.ResearchFindingListResponse, error)
	ZmienUstalenie(ctx context.Context, z shared.ResearchFindingUpdateRequest) (shared.ResearchFindingUpdateResponse, error)
	UsunUstalenie(ctx context.Context, z shared.ResearchFindingRemoveRequest) (shared.ResearchFindingRemoveResponse, error)
	ScalUstalenia(ctx context.Context, z shared.ResearchFindingMergeRequest) (shared.ResearchFindingMergeResponse, error)
	SladUstalenia(ctx context.Context, z shared.ResearchFindingProvenanceRequest) (shared.ResearchFindingProvenanceResponse, error)
	PrzypiszKodyUstalenia(ctx context.Context, z shared.ResearchFindingCodeRequest) (shared.ResearchFindingCodeResponse, error)
	PobierzKsiazkeKodow(ctx context.Context, z shared.ResearchCodebookGetRequest) (shared.ResearchCodebookGetResponse, error)
	UstawKsiazkeKodow(ctx context.Context, z shared.ResearchCodebookSetRequest) (shared.ResearchCodebookSetResponse, error)
	MacierzKodowania(ctx context.Context, z shared.ResearchFindingMatrixRequest) (shared.ResearchFindingMatrixResponse, error)
	SzukajSprzecznosci(ctx context.Context, z shared.ResearchFindingContradictionsRequest) (shared.ResearchFindingContradictionsResponse, error)
	RozstrzygnijSprzecznosc(ctx context.Context, z shared.ResearchContradictionResolveRequest) (shared.ResearchContradictionResolveResponse, error)
	ZweryfikujUstalenie(ctx context.Context, z shared.ResearchFindingFactCheckRequest) (shared.ResearchFindingFactCheckResponse, error)
	PogrupujUstalenia(ctx context.Context, z shared.ResearchFindingClusterRequest) (shared.ResearchFindingClusterResponse, error)

	// przestrzeń badania
	UstawPrzestrzenPelna(ctx context.Context, z shared.ResearchWorkspaceSetRequest) (shared.ResearchWorkspaceSetResponse, error)
	PobierzPrzestrzen(ctx context.Context, z shared.ResearchWorkspaceGetRequest) (shared.ResearchWorkspaceGetResponse, error)
	UstawPytaniaBadania(ctx context.Context, z shared.ResearchWorkspaceQuestionSetRequest) (shared.ResearchWorkspaceQuestionSetResponse, error)
	PokryciePytan(ctx context.Context, z shared.ResearchWorkspaceCoverageRequest) (shared.ResearchWorkspaceCoverageResponse, error)
	UstawNotatkeBadania(ctx context.Context, z shared.ResearchWorkspaceNoteSetRequest) (shared.ResearchWorkspaceNoteSetResponse, error)
	SwiezoscBadania(ctx context.Context, z shared.ResearchWorkspaceFreshnessRequest) (shared.ResearchWorkspaceFreshnessResponse, error)
	SzukajLukBadania(ctx context.Context, z shared.ResearchGapFindRequest) (shared.ResearchGapFindResponse, error)
	PrzesiewPrisma(ctx context.Context, z shared.ResearchPrismaGetRequest) (shared.ResearchPrismaGetResponse, error)
	GrafDowodow(ctx context.Context, z shared.ResearchEvidenceGraphRequest) (shared.ResearchEvidenceGraphResponse, error)

	// cytowania
	ZlozCytowania(ctx context.Context, z shared.ResearchCitationRenderRequest) (shared.ResearchCitationRenderResponse, error)
	WypiszStyleCytowania(ctx context.Context, z shared.ResearchCitationStylesRequest) (shared.ResearchCitationStylesResponse, error)
	SprawdzCytowania(ctx context.Context, z shared.ResearchCitationCheckRequest) (shared.ResearchCitationCheckResponse, error)
	SprawdzWycofania(ctx context.Context, z shared.ResearchRetractionCheckRequest) (shared.ResearchRetractionCheckResponse, error)

	// raport
	ZbudujRaport(ctx context.Context, z shared.ResearchReportBuildRequest) (shared.ResearchReportBuildResponse, error)
	PobierzRaport(ctx context.Context, z shared.ResearchReportGetRequest) (shared.ResearchReportGetResponse, error)
	WypiszSzablonyRaportu(ctx context.Context, z shared.ResearchReportTemplateListRequest) (shared.ResearchReportTemplateListResponse, error)
	StreszczRaport(ctx context.Context, z shared.ResearchReportSummarizeRequest) (shared.ResearchReportSummarizeResponse, error)
	ZlozBibliografie(ctx context.Context, z shared.ResearchReportBibliographyRequest) (shared.ResearchReportBibliographyResponse, error)
	UstawPrzypisyRaportu(ctx context.Context, z shared.ResearchReportFootnoteSetRequest) (shared.ResearchReportFootnoteSetResponse, error)
	WstawBlokRaportu(ctx context.Context, z shared.ResearchReportInsertRequest) (shared.ResearchReportInsertResponse, error)
	WypiszWersjeRaportu(ctx context.Context, z shared.ResearchReportVersionListRequest) (shared.ResearchReportVersionListResponse, error)
	PorownajWersjeRaportu(ctx context.Context, z shared.ResearchReportDiffRequest) (shared.ResearchReportDiffResponse, error)
	DodajKomentarzRaportu(ctx context.Context, z shared.ResearchReportCommentAddRequest) (shared.ResearchReportCommentAddResponse, error)
	WypiszKomentarzeRaportu(ctx context.Context, z shared.ResearchReportCommentListRequest) (shared.ResearchReportCommentListResponse, error)
	OperacjaKontekstowaRaportu(ctx context.Context, z shared.ResearchReportContextualOpRequest) (shared.ResearchReportContextualOpResponse, error)

	// eksport
	WyeksportujRaport(ctx context.Context, z shared.ResearchReportExportRequest) (shared.ResearchReportExportResponse, error)
	PodejrzyjEksport(ctx context.Context, z shared.ResearchExportPreviewRequest) (shared.ResearchExportPreviewResponse, error)
	WypiszEksporty(ctx context.Context, z shared.ResearchExportListRequest) (shared.ResearchExportListResponse, error)
	ZapiszSzablonEksportu(ctx context.Context, z shared.ResearchExportTemplateSetRequest) (shared.ResearchExportTemplateSetResponse, error)
	WypiszSzablonyEksportu(ctx context.Context, z shared.ResearchExportTemplateListRequest) (shared.ResearchExportTemplateListResponse, error)
	UdostepnijRaport(ctx context.Context, z shared.ResearchExportShareRequest) (shared.ResearchExportShareResponse, error)
}

// Adapter wypełnia port w całości. Gdyby port i adapter się rozjechały,
// kompilacja stanie tutaj, a nie dopiero na martwej komendzie u Operatora.
var _ Badania = (*adapterBadan)(nil)

// zarejestrujBadania wpina komplet komend modułu Research.
func zarejestrujBadania(r *Rejestr, m Badania, e *emiter) {
	if r == nil || m == nil {
		return
	}

	// ── źródła ──
	r.Zarejestruj(shared.CommandResearchSourceList, obsluz(m.WypiszZrodla))
	r.Zarejestruj(shared.CommandResearchSourceMerge, obsluz(m.ScalZrodla))
	r.Zarejestruj(shared.CommandResearchSourceDuplicates, obsluz(m.SzukajDuplikatowZrodel))
	r.Zarejestruj(shared.CommandResearchSourceTag, obsluz(m.OznaczZrodlo))
	r.Zarejestruj(shared.CommandResearchSourceAttachmentAdd, obsluz(m.DodajZalacznikZrodla))
	r.Zarejestruj(shared.CommandResearchSourceAttachmentList, obsluz(m.WypiszZalacznikiZrodla))
	r.Zarejestruj(shared.CommandResearchSourceImport, obsluz(m.WczytajBibliografie))
	r.Zarejestruj(shared.CommandResearchSourceTranscribe, obsluz(m.PrzepiszNagranie))
	r.Zarejestruj(shared.CommandResearchSourceResolve, obsluz(m.RozstrzygnijIdentyfikator))

	// Trzy drogi wnoszące źródło rozgłaszają `research.source.changed`: okno
	// katalogu ma zobaczyć nową pozycję bez odpytywania. Zmiana i zdjęcie idą
	// tym samym zdarzeniem z innym rodzajem zmiany.
	r.Zarejestruj(shared.CommandResearchSourceAdd,
		obsluz(func(ctx context.Context, z shared.ResearchSourceAddRequest) (shared.ResearchSourceAddResponse, error) {
			odpowiedz, err := m.DodajZrodlo(ctx, z)
			if err == nil {
				e.zrodloBadania(shared.ChangeKindCreated, odpowiedz.Source)
			}
			return odpowiedz, err
		}))
	r.Zarejestruj(shared.CommandResearchSourceCapture,
		obsluz(func(ctx context.Context, z shared.ResearchSourceCaptureRequest) (shared.ResearchSourceCaptureResponse, error) {
			odpowiedz, err := m.PrzechwycStrone(ctx, z)
			if err == nil {
				e.zrodloBadania(shared.ChangeKindCreated, odpowiedz.Source)
			}
			return odpowiedz, err
		}))
	r.Zarejestruj(shared.CommandResearchSourceUpdate,
		obsluz(func(ctx context.Context, z shared.ResearchSourceUpdateRequest) (shared.ResearchSourceUpdateResponse, error) {
			odpowiedz, err := m.ZmienZrodlo(ctx, z)
			if err == nil {
				e.zrodloBadania(shared.ChangeKindUpdated, odpowiedz.Source)
			}
			return odpowiedz, err
		}))
	r.Zarejestruj(shared.CommandResearchSourceRemove,
		obsluz(func(ctx context.Context, z shared.ResearchSourceRemoveRequest) (shared.ResearchSourceRemoveResponse, error) {
			odpowiedz, err := m.UsunZrodlo(ctx, z)
			if err == nil {
				// Zdarzenie niesie samo wskazanie zdjętej pozycji: wiersza już nie
				// ma, więc pełny byt byłby odtworzeniem czegoś, co nie istnieje.
				e.zrodloBadania(shared.ChangeKindDeleted,
					shared.ResearchSource{Id: odpowiedz.SourceId})
			}
			return odpowiedz, err
		}))

	// ── odkrywanie, monitory, import wsadowy ──
	r.Zarejestruj(shared.CommandResearchDiscoverySearch, obsluz(m.SzukajZrodel))
	r.Zarejestruj(shared.CommandResearchDiscoveryAssist, obsluz(m.UlozZapytania))
	r.Zarejestruj(shared.CommandResearchDiscoverySnowball, obsluz(m.RozwinCytowania))
	r.Zarejestruj(shared.CommandResearchDiscoveryReject, obsluz(m.OdrzucWyniki))
	r.Zarejestruj(shared.CommandResearchMonitorList, obsluz(m.WypiszMonitory))
	r.Zarejestruj(shared.CommandResearchBatchImport, obsluz(m.WczytajPartieAdresow))

	r.Zarejestruj(shared.CommandResearchMonitorSet,
		obsluz(func(ctx context.Context, z shared.ResearchMonitorSetRequest) (shared.ResearchMonitorSetResponse, error) {
			odpowiedz, err := m.UstawMonitor(ctx, z)
			if err == nil {
				e.monitorBadania(odpowiedz.Monitor, 0)
			}
			return odpowiedz, err
		}))
	r.Zarejestruj(shared.CommandResearchMonitorRefresh,
		obsluz(func(ctx context.Context, z shared.ResearchMonitorRefreshRequest) (shared.ResearchMonitorRefreshResponse, error) {
			odpowiedz, err := m.OdswiezMonitory(ctx, z)
			if err == nil {
				for _, monitor := range odpowiedz.Monitors {
					nowe := 0
					if monitor.PendingCount != nil {
						nowe = *monitor.PendingCount
					}
					e.monitorBadania(monitor, nowe)
				}
			}
			return odpowiedz, err
		}))

	// ── lektura i ekstrakcja ──
	r.Zarejestruj(shared.CommandResearchReadingOpen, obsluz(m.OtworzDoLektury))
	r.Zarejestruj(shared.CommandResearchAnnotationAdd, obsluz(m.DodajAdnotacje))
	r.Zarejestruj(shared.CommandResearchAnnotationList, obsluz(m.WypiszAdnotacje))
	r.Zarejestruj(shared.CommandResearchAnnotationRemove, obsluz(m.UsunAdnotacje))
	r.Zarejestruj(shared.CommandResearchExcerptList, obsluz(m.WypiszWypisy))
	r.Zarejestruj(shared.CommandResearchSourceSummarize, obsluz(m.StreszczZrodlo))
	r.Zarejestruj(shared.CommandResearchSourceExtractTable, obsluz(m.WyodrebnijTabele))
	r.Zarejestruj(shared.CommandResearchSourceExtractClaims, obsluz(m.WyodrebnijTwierdzenia))
	r.Zarejestruj(shared.CommandResearchSourceOcr, obsluz(m.RozpoznajPismoZrodla))
	r.Zarejestruj(shared.CommandResearchCorpusAsk, obsluz(m.ZapytajKorpus))

	// ── ustalenia i analiza ──
	r.Zarejestruj(shared.CommandResearchFindingList, obsluz(m.WypiszUstalenia))
	r.Zarejestruj(shared.CommandResearchFindingProvenance, obsluz(m.SladUstalenia))
	r.Zarejestruj(shared.CommandResearchFindingCode, obsluz(m.PrzypiszKodyUstalenia))
	r.Zarejestruj(shared.CommandResearchCodebookGet, obsluz(m.PobierzKsiazkeKodow))
	r.Zarejestruj(shared.CommandResearchCodebookSet, obsluz(m.UstawKsiazkeKodow))
	r.Zarejestruj(shared.CommandResearchFindingMatrix, obsluz(m.MacierzKodowania))
	r.Zarejestruj(shared.CommandResearchFindingContradictions, obsluz(m.SzukajSprzecznosci))
	r.Zarejestruj(shared.CommandResearchContradictionResolve, obsluz(m.RozstrzygnijSprzecznosc))
	r.Zarejestruj(shared.CommandResearchFindingFactCheck, obsluz(m.ZweryfikujUstalenie))
	r.Zarejestruj(shared.CommandResearchFindingCluster, obsluz(m.PogrupujUstalenia))

	r.Zarejestruj(shared.CommandResearchFindingAdd,
		obsluz(func(ctx context.Context, z shared.ResearchFindingAddRequest) (shared.ResearchFindingAddResponse, error) {
			odpowiedz, err := m.DodajUstalenie(ctx, z)
			if err == nil {
				e.ustalenieBadania(zmianaUstaleniaBadania(z.FindingId), odpowiedz.Finding)
			}
			return odpowiedz, err
		}))
	r.Zarejestruj(shared.CommandResearchFindingUpdate,
		obsluz(func(ctx context.Context, z shared.ResearchFindingUpdateRequest) (shared.ResearchFindingUpdateResponse, error) {
			odpowiedz, err := m.ZmienUstalenie(ctx, z)
			if err == nil {
				e.ustalenieBadania(shared.ChangeKindUpdated, odpowiedz.Finding)
			}
			return odpowiedz, err
		}))
	r.Zarejestruj(shared.CommandResearchFindingMerge,
		obsluz(func(ctx context.Context, z shared.ResearchFindingMergeRequest) (shared.ResearchFindingMergeResponse, error) {
			odpowiedz, err := m.ScalUstalenia(ctx, z)
			if err == nil {
				e.ustalenieBadania(shared.ChangeKindUpdated, odpowiedz.Finding)
			}
			return odpowiedz, err
		}))
	r.Zarejestruj(shared.CommandResearchFindingRemove,
		obsluz(func(ctx context.Context, z shared.ResearchFindingRemoveRequest) (shared.ResearchFindingRemoveResponse, error) {
			odpowiedz, err := m.UsunUstalenie(ctx, z)
			if err == nil {
				e.ustalenieBadania(shared.ChangeKindDeleted,
					shared.ResearchFinding{Id: odpowiedz.FindingId})
			}
			return odpowiedz, err
		}))

	// ── przestrzeń badania ──
	r.Zarejestruj(shared.CommandResearchWorkspaceSet, obsluz(m.UstawPrzestrzenPelna))
	r.Zarejestruj(shared.CommandResearchWorkspaceGet, obsluz(m.PobierzPrzestrzen))
	r.Zarejestruj(shared.CommandResearchWorkspaceQuestionSet, obsluz(m.UstawPytaniaBadania))
	r.Zarejestruj(shared.CommandResearchWorkspaceCoverage, obsluz(m.PokryciePytan))
	r.Zarejestruj(shared.CommandResearchWorkspaceNoteSet, obsluz(m.UstawNotatkeBadania))
	r.Zarejestruj(shared.CommandResearchWorkspaceFreshness, obsluz(m.SwiezoscBadania))
	r.Zarejestruj(shared.CommandResearchGapFind, obsluz(m.SzukajLukBadania))
	r.Zarejestruj(shared.CommandResearchPrismaGet, obsluz(m.PrzesiewPrisma))
	r.Zarejestruj(shared.CommandResearchEvidenceGraph, obsluz(m.GrafDowodow))

	// ── cytowania ──
	r.Zarejestruj(shared.CommandResearchCitationRender, obsluz(m.ZlozCytowania))
	r.Zarejestruj(shared.CommandResearchCitationStyles, obsluz(m.WypiszStyleCytowania))
	r.Zarejestruj(shared.CommandResearchCitationCheck, obsluz(m.SprawdzCytowania))
	r.Zarejestruj(shared.CommandResearchRetractionCheck, obsluz(m.SprawdzWycofania))

	// ── raport ──
	r.Zarejestruj(shared.CommandResearchReportGet, obsluz(m.PobierzRaport))
	r.Zarejestruj(shared.CommandResearchReportTemplateList, obsluz(m.WypiszSzablonyRaportu))
	r.Zarejestruj(shared.CommandResearchReportBibliography, obsluz(m.ZlozBibliografie))
	r.Zarejestruj(shared.CommandResearchReportFootnoteSet, obsluz(m.UstawPrzypisyRaportu))
	r.Zarejestruj(shared.CommandResearchReportInsert, obsluz(m.WstawBlokRaportu))
	r.Zarejestruj(shared.CommandResearchReportVersionList, obsluz(m.WypiszWersjeRaportu))
	r.Zarejestruj(shared.CommandResearchReportDiff, obsluz(m.PorownajWersjeRaportu))
	r.Zarejestruj(shared.CommandResearchReportCommentAdd, obsluz(m.DodajKomentarzRaportu))
	r.Zarejestruj(shared.CommandResearchReportCommentList, obsluz(m.WypiszKomentarzeRaportu))
	r.Zarejestruj(shared.CommandResearchReportSummarize, obsluz(m.StreszczRaport))
	r.Zarejestruj(shared.CommandResearchReportContextualOp, obsluz(m.OperacjaKontekstowaRaportu))

	r.Zarejestruj(shared.CommandResearchReportBuild,
		obsluz(func(ctx context.Context, z shared.ResearchReportBuildRequest) (shared.ResearchReportBuildResponse, error) {
			odpowiedz, err := m.ZbudujRaport(ctx, z)
			if err == nil {
				e.raportBadania(zmianaRaportu(z.ReportId), odpowiedz.Report)
			}
			return odpowiedz, err
		}))

	// ── eksport ──
	r.Zarejestruj(shared.CommandResearchReportExport, obsluz(m.WyeksportujRaport))
	r.Zarejestruj(shared.CommandResearchExportPreview, obsluz(m.PodejrzyjEksport))
	r.Zarejestruj(shared.CommandResearchExportList, obsluz(m.WypiszEksporty))
	r.Zarejestruj(shared.CommandResearchExportTemplateSet, obsluz(m.ZapiszSzablonEksportu))
	r.Zarejestruj(shared.CommandResearchExportTemplateList, obsluz(m.WypiszSzablonyEksportu))
	r.Zarejestruj(shared.CommandResearchExportShare, obsluz(m.UdostepnijRaport))
}

// zmianaRaportu odróżnia raport nowy od rozbudowywanego na podstawie sygnału
// dostępnego w tej warstwie: żądanie bez `reportId` zakłada raport nowy
// (created), żądanie z `reportId` — aktualizację zastanego (updated).
func zmianaRaportu(reportId *string) shared.ChangeKind {
	if reportId == nil || *reportId == "" {
		return shared.ChangeKindCreated
	}
	return shared.ChangeKindUpdated
}

// zmianaUstaleniaBadania odróżnia ustalenie nowe od zmienianego tym samym
// sygnałem co przy raporcie.
func zmianaUstaleniaBadania(findingId *string) shared.ChangeKind {
	if findingId == nil || *findingId == "" {
		return shared.ChangeKindCreated
	}
	return shared.ChangeKindUpdated
}

// raportBadania rozgłasza `research.report.changed`. Raport jest bytem okna
// badania, nie karty sesji, więc zdarzenie idzie bez wskazania sesji — okno
// Report Builder odbiera je stroną własną (wzór `przebiegAutomatyki`).
func (e *emiter) raportBadania(zmiana shared.ChangeKind, raport shared.ResearchReport) {
	e.wyslij(shared.EventResearchReportChanged, "",
		shared.ResearchReportChangedEvent{Change: zmiana, Report: raport})
}

// zrodloBadania rozgłasza `research.source.changed` — Sources Manager odbiera
// zmianę katalogu bez odpytywania.
func (e *emiter) zrodloBadania(zmiana shared.ChangeKind, zrodlo shared.ResearchSource) {
	e.wyslij(shared.EventResearchSourceChanged, "",
		shared.ResearchSourceChangedEvent{Change: zmiana, Source: zrodlo})
}

// ustalenieBadania rozgłasza `research.finding.changed` — Findings Panel
// odbiera zmianę rejestru ustaleń bez odpytywania.
func (e *emiter) ustalenieBadania(zmiana shared.ChangeKind, ustalenie shared.ResearchFinding) {
	e.wyslij(shared.EventResearchFindingChanged, "",
		shared.ResearchFindingChangedEvent{Change: zmiana, Finding: ustalenie})
}

// monitorBadania rozgłasza `research.monitor.changed` — skrzynka „nowe źródła"
// Discovery Panelu dowiaduje się o odświeżeniu bez odpytywania.
func (e *emiter) monitorBadania(monitor shared.ResearchMonitor, nowych int) {
	e.wyslij(shared.EventResearchMonitorChanged, "",
		shared.ResearchMonitorChangedEvent{Monitor: monitor, NewCount: nowych})
}
