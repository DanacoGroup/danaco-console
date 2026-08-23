// Wpięcie komend modułu Studio: port i funkcja rejestrująca. Metody
// portu leżą w plikach `adapter_modul_studio.go`,
// `adapter_modul_studio_wersje.go` i `adapter_modul_studio_roznice.go`; ten
// plik jest jedynym miejscem, które je razem nazywa.
//
// Zdarzenie `studio.document.changed` dotyczy stanu dokumentu, nie samej
// komendy, więc rozgłasza się po `document.save`, po `repository.restore` i po
// `contextual.op` — wszystkie trzy zmieniają treść widoczną w oknie pracy
// z dokumentem. `document.open`, `repository.list` i `diff.compare` treści nie
// zmieniają i zdarzenia nie mają.
package core

import (
	"context"
	"encoding/json"
	"strings"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// Studio jest portem modułu Studio.
type Studio interface {
	OtworzDokument(ctx context.Context, z shared.StudioDocumentOpenRequest) (shared.StudioDocumentOpenResponse, error)
	ZapiszDokument(ctx context.Context, z shared.StudioDocumentSaveRequest) (shared.StudioDocumentSaveResponse, error)
	OperacjaKontekstowa(ctx context.Context, z shared.StudioContextualOpRequest) (shared.StudioContextualOpResponse, error)
	Porownaj(ctx context.Context, z shared.StudioDiffCompareRequest) (shared.StudioDiffCompareResponse, error)
	Wersje(ctx context.Context, z shared.StudioRepositoryListRequest) (shared.StudioRepositoryListResponse, error)
	PrzywrocWersje(ctx context.Context, z shared.StudioRepositoryRestoreRequest) (shared.StudioRepositoryRestoreResponse, error)

	// Ingest/OCR Panel — kolejka wczytywania i cyfryzacja
	// (`adapter_modul_studio_cyfryzacja.go`).
	DolozDoKolejki(ctx context.Context, z shared.StudioIngestQueueAddRequest) (shared.StudioIngestQueueAddResponse, error)
	KolejkaWczytywania(ctx context.Context, z shared.StudioIngestQueueListRequest) (shared.StudioIngestQueueListResponse, error)
	Rozpoznaj(ctx context.Context, z shared.StudioIngestRecognizeRequest) (shared.StudioIngestRecognizeResponse, error)
	PoprawRozpoznanie(ctx context.Context, z shared.StudioIngestCorrectionSetRequest) (shared.StudioIngestCorrectionSetResponse, error)
	PrzyjmijPozycje(ctx context.Context, z shared.StudioIngestItemAcceptRequest) (shared.StudioIngestItemAcceptResponse, error)
	UrzadzeniaWejsciowe(ctx context.Context, z shared.StudioIngestDeviceListRequest) (shared.StudioIngestDeviceListResponse, error)

	// Komentarze redakcyjne, adnotacje różnic, śledzenie zmian i decyzja
	// o propozycji (`adapter_modul_studio_adnotacje.go`).
	DodajKomentarz(ctx context.Context, z shared.StudioCommentAddRequest) (shared.StudioCommentAddResponse, error)
	Komentarze(ctx context.Context, z shared.StudioCommentListRequest) (shared.StudioCommentListResponse, error)
	RozstrzygnijKomentarz(ctx context.Context, z shared.StudioCommentResolveRequest) (shared.StudioCommentResolveResponse, error)
	DodajAdnotacje(ctx context.Context, z shared.StudioAnnotationAddRequest) (shared.StudioAnnotationAddResponse, error)
	Adnotacje(ctx context.Context, z shared.StudioAnnotationListRequest) (shared.StudioAnnotationListResponse, error)
	UstawSledzenie(ctx context.Context, z shared.StudioTrackingSetRequest) (shared.StudioTrackingSetResponse, error)
	ZmianySledzone(ctx context.Context, z shared.StudioTrackingListRequest) (shared.StudioTrackingListResponse, error)
	RozstrzygnijZmiany(ctx context.Context, z shared.StudioTrackingDecideRequest) (shared.StudioTrackingDecideResponse, error)
	RozstrzygnijPropozycje(ctx context.Context, z shared.StudioProposalDecideRequest) (shared.StudioProposalDecideResponse, error)

	// Katalogi zasięgu i cechy dokumentu (`adapter_modul_studio_katalogi.go`).
	ZapiszOperacje(ctx context.Context, z shared.StudioOperationSaveRequest) (shared.StudioOperationSaveResponse, error)
	Operacje(ctx context.Context, z shared.StudioOperationListRequest) (shared.StudioOperationListResponse, error)
	UsunOperacje(ctx context.Context, z shared.StudioOperationDeleteRequest) (shared.StudioOperationDeleteResponse, error)
	ZapiszLancuch(ctx context.Context, z shared.StudioChainSaveRequest) (shared.StudioChainSaveResponse, error)
	Lancuchy(ctx context.Context, z shared.StudioChainListRequest) (shared.StudioChainListResponse, error)
	UruchomLancuch(ctx context.Context, z shared.StudioChainRunRequest) (shared.StudioChainRunResponse, error)
	Szablony(ctx context.Context, z shared.StudioTemplateListRequest) (shared.StudioTemplateListResponse, error)
	ZastosujSzablon(ctx context.Context, z shared.StudioTemplateApplyRequest) (shared.StudioTemplateApplyResponse, error)
	ZapiszProfilWydania(ctx context.Context, z shared.StudioExportProfileSaveRequest) (shared.StudioExportProfileSaveResponse, error)
	ProfileWydania(ctx context.Context, z shared.StudioExportProfileListRequest) (shared.StudioExportProfileListResponse, error)
	UstawFormatDokumentu(ctx context.Context, z shared.StudioDocumentFormatSetRequest) (shared.StudioDocumentFormatSetResponse, error)
	UstawEtykieteWersji(ctx context.Context, z shared.StudioVersionLabelSetRequest) (shared.StudioVersionLabelSetResponse, error)

	// Gałęzie dokumentu i odwołanie do wersji
	// (`adapter_modul_studio_galezie.go`).
	ZalozGalaz(ctx context.Context, z shared.StudioBranchCreateRequest) (shared.StudioBranchCreateResponse, error)
	Galezie(ctx context.Context, z shared.StudioBranchListRequest) (shared.StudioBranchListResponse, error)
	ScalGalezie(ctx context.Context, z shared.StudioBranchMergeRequest) (shared.StudioBranchMergeResponse, error)
	UtworzOdwolanieWersji(ctx context.Context, z shared.StudioVersionReferenceCreateRequest) (shared.StudioVersionReferenceCreateResponse, error)

	// Wydawanie pracy na zewnątrz (`adapter_modul_studio_wydanie.go`).
	WydajRepozytorium(ctx context.Context, z shared.StudioRepositoryExportRequest) (shared.StudioRepositoryExportResponse, error)
	WydajPaczke(ctx context.Context, z shared.StudioPackageExportRequest) (shared.StudioPackageExportResponse, error)
	WydajRaportRoznicy(ctx context.Context, z shared.StudioDiffReportExportRequest) (shared.StudioDiffReportExportResponse, error)

	// Preview Window i różnica wizualna (`adapter_modul_studio_podglad.go`).
	WyrenderujPodglad(ctx context.Context, z shared.StudioPreviewRenderRequest) (shared.StudioPreviewRenderResponse, error)
	PorownajWizualnie(ctx context.Context, z shared.StudioDiffVisualRequest) (shared.StudioDiffVisualResponse, error)

	// Praca na treści: wsad, osadzenie zasobu, wyszukiwanie znaczeniowe
	// i zestawienie ze źródłem (`adapter_modul_studio_wsad.go`).
	UruchomWsad(ctx context.Context, z shared.StudioBatchRunRequest) (shared.StudioBatchRunResponse, error)
	OsadzZasob(ctx context.Context, z shared.StudioAssetEmbedRequest) (shared.StudioAssetEmbedResponse, error)
	WyszukajZnaczeniowo(ctx context.Context, z shared.StudioSearchSemanticRequest) (shared.StudioSearchSemanticResponse, error)
	PorownajZeZrodlem(ctx context.Context, z shared.StudioDiffSourceRequest) (shared.StudioDiffSourceResponse, error)

	// Wejścia do kolejki wczytywania (`adapter_modul_studio_wczytanie.go`).
	WczytajZAdresu(ctx context.Context, z shared.StudioIngestUrlRequest) (shared.StudioIngestUrlResponse, error)
	SkanujUrzadzenie(ctx context.Context, z shared.StudioIngestDeviceScanRequest) (shared.StudioIngestDeviceScanResponse, error)

	// Postać dokumentu i praca na treści fragmentu
	// (`adapter_modul_studio_postac.go`).
	PostacDokumentu(ctx context.Context, z shared.StudioDocumentFormGetRequest) (shared.StudioDocumentFormGetResponse, error)
	ZapiszPostacDokumentu(ctx context.Context, z shared.StudioDocumentFormSaveRequest) (shared.StudioDocumentFormSaveResponse, error)
	TrescFragmentu(ctx context.Context, z shared.StudioTextGetRequest) (shared.StudioTextGetResponse, error)
	ZmienTresc(ctx context.Context, z shared.StudioTextEditRequest) (shared.StudioTextEditResponse, error)

	// Styl znaku i akapitu, malarz postaci, zamiana z postacią
	// (`adapter_modul_studio_format.go`).
	PostacZnaku(ctx context.Context, z shared.StudioFormatCharacterGetRequest) (shared.StudioFormatCharacterGetResponse, error)
	UstawPostacZnaku(ctx context.Context, z shared.StudioFormatCharacterSetRequest) (shared.StudioFormatCharacterSetResponse, error)
	PostacAkapitu(ctx context.Context, z shared.StudioFormatParagraphGetRequest) (shared.StudioFormatParagraphGetResponse, error)
	UstawPostacAkapitu(ctx context.Context, z shared.StudioFormatParagraphSetRequest) (shared.StudioFormatParagraphSetResponse, error)
	CzyscPostac(ctx context.Context, z shared.StudioFormatClearRequest) (shared.StudioFormatClearResponse, error)
	UstawWielkoscLiter(ctx context.Context, z shared.StudioFormatCaseSetRequest) (shared.StudioFormatCaseSetResponse, error)
	ZabierzPostac(ctx context.Context, z shared.StudioFormatPainterCopyRequest) (shared.StudioFormatPainterCopyResponse, error)
	PolozPostac(ctx context.Context, z shared.StudioFormatPainterApplyRequest) (shared.StudioFormatPainterApplyResponse, error)
	ZaznaczPodobne(ctx context.Context, z shared.StudioFormatSimilarSelectRequest) (shared.StudioFormatSimilarSelectResponse, error)
	ZamienZPostacia(ctx context.Context, z shared.StudioFormatReplaceRequest) (shared.StudioFormatReplaceResponse, error)

	// Arkusz stylów nazwanych (`adapter_modul_studio_style.go`).
	StyleDokumentu(ctx context.Context, z shared.StudioStyleListRequest) (shared.StudioStyleListResponse, error)
	ZapiszStyl(ctx context.Context, z shared.StudioStyleSaveRequest) (shared.StudioStyleSaveResponse, error)
	ZastosujStyl(ctx context.Context, z shared.StudioStyleApplyRequest) (shared.StudioStyleApplyResponse, error)
	UsunStyl(ctx context.Context, z shared.StudioStyleDeleteRequest) (shared.StudioStyleDeleteResponse, error)

	// Blokady fragmentów (`adapter_modul_studio_blokady.go`).
	ZalozBlokade(ctx context.Context, z shared.StudioLockAddRequest) (shared.StudioLockAddResponse, error)
	Blokady(ctx context.Context, z shared.StudioLockListRequest) (shared.StudioLockListResponse, error)
	ZdejmijBlokade(ctx context.Context, z shared.StudioLockRemoveRequest) (shared.StudioLockRemoveResponse, error)

	// Nastawy strony, sekcje, nagłówki i stopki, numeracja, znak wodny,
	// koperta, podział i tabulatory linijki (`adapter_modul_studio_strona.go`).
	WstawPodzial(ctx context.Context, z shared.StudioPageBreakInsertRequest) (shared.StudioPageBreakInsertResponse, error)
	UstawNadrukKoperty(ctx context.Context, z shared.StudioPageEnvelopeSetRequest) (shared.StudioPageEnvelopeSetResponse, error)
	NaglowkiIStopki(ctx context.Context, z shared.StudioPageHeaderfooterGetRequest) (shared.StudioPageHeaderfooterGetResponse, error)
	UstawNaglowekIStopke(ctx context.Context, z shared.StudioPageHeaderfooterSetRequest) (shared.StudioPageHeaderfooterSetResponse, error)
	UstawNumeracjeStron(ctx context.Context, z shared.StudioPageNumberingSetRequest) (shared.StudioPageNumberingSetResponse, error)
	NosnikiStrony(ctx context.Context, z shared.StudioPagePaperListRequest) (shared.StudioPagePaperListResponse, error)
	NastawyStrony(ctx context.Context, z shared.StudioPageSetupGetRequest) (shared.StudioPageSetupGetResponse, error)
	UstawNastawyStrony(ctx context.Context, z shared.StudioPageSetupSetRequest) (shared.StudioPageSetupSetResponse, error)
	UstawZnakWodny(ctx context.Context, z shared.StudioPageWatermarkSetRequest) (shared.StudioPageWatermarkSetResponse, error)
	UstawTabulatorLinijki(ctx context.Context, z shared.StudioRulerTabstopSetRequest) (shared.StudioRulerTabstopSetResponse, error)
	UsunSekcjeDokumentu(ctx context.Context, z shared.StudioSectionDeleteRequest) (shared.StudioSectionDeleteResponse, error)
	SekcjeDokumentu(ctx context.Context, z shared.StudioSectionListRequest) (shared.StudioSectionListResponse, error)
	ZapiszSekcjeDokumentu(ctx context.Context, z shared.StudioSectionSaveRequest) (shared.StudioSectionSaveResponse, error)

	// Listy i punktatory (`adapter_modul_studio_listy.go`).
	ZastosujListe(ctx context.Context, z shared.StudioListApplyRequest) (shared.StudioListApplyResponse, error)
	UstawPunktatorListy(ctx context.Context, z shared.StudioListBulletSetRequest) (shared.StudioListBulletSetResponse, error)
	PrzestawPoziomListy(ctx context.Context, z shared.StudioListLevelIndentRequest) (shared.StudioListLevelIndentResponse, error)
	UstawNumeracjeListy(ctx context.Context, z shared.StudioListNumberingSetRequest) (shared.StudioListNumberingSetResponse, error)
	WznowNumeracjeListy(ctx context.Context, z shared.StudioListRestartRequest) (shared.StudioListRestartResponse, error)

	// Symbole i autozamiana znaków (`adapter_modul_studio_symbole.go`).
	ZasadyAutozamiany(ctx context.Context, z shared.StudioSymbolAutoreplaceListRequest) (shared.StudioSymbolAutoreplaceListResponse, error)
	UstawZasadeAutozamiany(ctx context.Context, z shared.StudioSymbolAutoreplaceSetRequest) (shared.StudioSymbolAutoreplaceSetResponse, error)
	WstawZnak(ctx context.Context, z shared.StudioSymbolInsertRequest) (shared.StudioSymbolInsertResponse, error)
	TabliceZnakow(ctx context.Context, z shared.StudioSymbolListRequest) (shared.StudioSymbolListResponse, error)

	// Tabele dokumentu (`adapter_modul_studio_tabele.go`).
	ZamienTabeleITekst(ctx context.Context, z shared.StudioTableConvertRequest) (shared.StudioTableConvertResponse, error)
	UstawPostacTabeli(ctx context.Context, z shared.StudioTableFormatSetRequest) (shared.StudioTableFormatSetResponse, error)
	WstawTabele(ctx context.Context, z shared.StudioTableInsertRequest) (shared.StudioTableInsertResponse, error)
	WykazTabel(ctx context.Context, z shared.StudioTableListRequest) (shared.StudioTableListResponse, error)
	SortujTabele(ctx context.Context, z shared.StudioTableSortRequest) (shared.StudioTableSortResponse, error)
	ZmienBudoweTabeli(ctx context.Context, z shared.StudioTableStructureEditRequest) (shared.StudioTableStructureEditResponse, error)

	// Obiekty osadzone — obrazy, kształty, pola tekstowe
	// (`adapter_modul_studio_wstawienia.go`).
	UstawPostacObiektu(ctx context.Context, z shared.StudioObjectFormatSetRequest) (shared.StudioObjectFormatSetResponse, error)
	WstawObiekt(ctx context.Context, z shared.StudioObjectInsertRequest) (shared.StudioObjectInsertResponse, error)
	WykazObiektow(ctx context.Context, z shared.StudioObjectListRequest) (shared.StudioObjectListResponse, error)
	UsunObiekt(ctx context.Context, z shared.StudioObjectRemoveRequest) (shared.StudioObjectRemoveResponse, error)

	// Aparat dokumentu — spisy, przypisy, bibliografia, indeks
	// (`adapter_modul_studio_aparat.go`).
	WstawElementAparatu(ctx context.Context, z shared.StudioApparatusInsertRequest) (shared.StudioApparatusInsertResponse, error)
	WykazAparatu(ctx context.Context, z shared.StudioApparatusListRequest) (shared.StudioApparatusListResponse, error)
	OdswiezAparat(ctx context.Context, z shared.StudioApparatusRefreshRequest) (shared.StudioApparatusRefreshResponse, error)
	UsunElementAparatuDokumentu(ctx context.Context, z shared.StudioApparatusRemoveRequest) (shared.StudioApparatusRemoveResponse, error)

	// Pola dokumentu (`adapter_modul_studio_pola.go`).
	WstawPole(ctx context.Context, z shared.StudioFieldInsertRequest) (shared.StudioFieldInsertResponse, error)
	WykazPol(ctx context.Context, z shared.StudioFieldListRequest) (shared.StudioFieldListResponse, error)
	OdswiezPola(ctx context.Context, z shared.StudioFieldRefreshRequest) (shared.StudioFieldRefreshResponse, error)

	// Wejście do edytora, zapis pod nazwą, kopia i wydanie do formatów
	// (`adapter_modul_studio_wejscie_czynnosci.go`, `_wydanie_formatu.go`).
	SkopiujDokument(ctx context.Context, z shared.StudioDocumentCopyRequest) (shared.StudioDocumentCopyResponse, error)
	ZalozDokument(ctx context.Context, z shared.StudioDocumentCreateRequest) (shared.StudioDocumentCreateResponse, error)
	WniesObraz(ctx context.Context, z shared.StudioDocumentImageImportRequest) (shared.StudioDocumentImageImportResponse, error)
	WniesPlikDoEdytora(ctx context.Context, z shared.StudioDocumentImportFileRequest) (shared.StudioDocumentImportFileResponse, error)
	WniesPdfDoEdytora(ctx context.Context, z shared.StudioDocumentImportPdfRequest) (shared.StudioDocumentImportPdfResponse, error)
	ZapiszDokumentPodNazwa(ctx context.Context, z shared.StudioDocumentSaveAsRequest) (shared.StudioDocumentSaveAsResponse, error)
	WniesZBiblioteki(ctx context.Context, z shared.StudioInsertFromLibraryRequest) (shared.StudioInsertFromLibraryResponse, error)
	WniesZeSieci(ctx context.Context, z shared.StudioInsertFromWebRequest) (shared.StudioInsertFromWebResponse, error)
	WydajWsadowo(ctx context.Context, z shared.StudioDocumentExportBatchRequest) (shared.StudioDocumentExportBatchResponse, error)
	WydajDoFormatu(ctx context.Context, z shared.StudioDocumentExportFormatRequest) (shared.StudioDocumentExportFormatResponse, error)

	// Warsztat szablonów pism (`adapter_modul_studio_szablony_pism.go`).
	UsunSzablonPisma(ctx context.Context, z shared.StudioTemplateDeleteRequest) (shared.StudioTemplateDeleteResponse, error)
	OddajSzablonDoPliku(ctx context.Context, z shared.StudioTemplateExportRequest) (shared.StudioTemplateExportResponse, error)
	PolaSzablonu(ctx context.Context, z shared.StudioTemplateFieldListRequest) (shared.StudioTemplateFieldListResponse, error)
	UstawPoleSzablonu(ctx context.Context, z shared.StudioTemplateFieldSetRequest) (shared.StudioTemplateFieldSetResponse, error)
	WypelnijSzablon(ctx context.Context, z shared.StudioTemplateFillRequest) (shared.StudioTemplateFillResponse, error)
	WniesSzablonZPliku(ctx context.Context, z shared.StudioTemplateImportRequest) (shared.StudioTemplateImportResponse, error)
	ZapiszSzablonPisma(ctx context.Context, z shared.StudioTemplateSaveRequest) (shared.StudioTemplateSaveResponse, error)

	// Odwracalny dziennik czynności dokumentu (`adapter_modul_studio_dziennik.go`).
	DziennikCzynnosci(ctx context.Context, z shared.StudioJournalListRequest) (shared.StudioJournalListResponse, error)
	PonowCzynnosc(ctx context.Context, z shared.StudioJournalRedoRequest) (shared.StudioJournalRedoResponse, error)
	CofnijCzynnosc(ctx context.Context, z shared.StudioJournalRevertRequest) (shared.StudioJournalRevertResponse, error)

	// Przełącznik zmian modelu (`adapter_modul_studio_zmiany_modelu.go`).
	ZmianyModelu(ctx context.Context, z shared.StudioModelChangesListRequest) (shared.StudioModelChangesListResponse, error)
	PrzeskocDoZmianyModelu(ctx context.Context, z shared.StudioModelChangesNavigateRequest) (shared.StudioModelChangesNavigateResponse, error)
	CofnijZmianyModelu(ctx context.Context, z shared.StudioModelChangesRevertRequest) (shared.StudioModelChangesRevertResponse, error)

	// Różnica wersji na postaci i przeniesienie fragmentu
	// (`adapter_modul_studio_roznica_wersji.go`).
	PorownajPostac(ctx context.Context, z shared.StudioDiffFormCompareRequest) (shared.StudioDiffFormCompareResponse, error)
	PrzeniesFragmentRoznicy(ctx context.Context, z shared.StudioDiffHunkApplyRequest) (shared.StudioDiffHunkApplyResponse, error)

	// Autozapis, kopie zapasowe i szeregi wersji
	// (`adapter_modul_studio_autozapis.go`).
	NastawyAutozapisu(ctx context.Context, z shared.StudioAutosaveGetRequest) (shared.StudioAutosaveGetResponse, error)
	WykonajAutozapis(ctx context.Context, z shared.StudioAutosaveRunRequest) (shared.StudioAutosaveRunResponse, error)
	UstawAutozapis(ctx context.Context, z shared.StudioAutosaveSetRequest) (shared.StudioAutosaveSetResponse, error)
	ZalozKopieZapasowa(ctx context.Context, z shared.StudioBackupCreateRequest) (shared.StudioBackupCreateResponse, error)
	KopieDokumentu(ctx context.Context, z shared.StudioBackupListRequest) (shared.StudioBackupListResponse, error)
	PrzywrocKopie(ctx context.Context, z shared.StudioBackupRestoreRequest) (shared.StudioBackupRestoreResponse, error)
	PrzywrocWersjeZalozycielska(ctx context.Context, z shared.StudioVersionRestoreInitialRequest) (shared.StudioVersionRestoreInitialResponse, error)
	WersjeWSzeregach(ctx context.Context, z shared.StudioVersionSeriesListRequest) (shared.StudioVersionSeriesListResponse, error)

	// Znakowanie: komentarze, propozycje, wyróżnienia i rodzaje znaczników
	// (`adapter_modul_studio_znakowanie.go`).
	DodajZnakowanie(ctx context.Context, z shared.StudioMarkupAddRequest) (shared.StudioMarkupAddResponse, error)
	RozstrzygnijZnakowanie(ctx context.Context, z shared.StudioMarkupDecideRequest) (shared.StudioMarkupDecideResponse, error)
	Znakowania(ctx context.Context, z shared.StudioMarkupListRequest) (shared.StudioMarkupListResponse, error)
	ZdejmijZnakowanie(ctx context.Context, z shared.StudioMarkupRemoveRequest) (shared.StudioMarkupRemoveResponse, error)
	UsunRodzajZnacznika(ctx context.Context, z shared.StudioMarkupTypeDeleteRequest) (shared.StudioMarkupTypeDeleteResponse, error)
	RodzajeZnacznika(ctx context.Context, z shared.StudioMarkupTypeListRequest) (shared.StudioMarkupTypeListResponse, error)
	ZapiszRodzajZnacznika(ctx context.Context, z shared.StudioMarkupTypeSaveRequest) (shared.StudioMarkupTypeSaveResponse, error)

	// Zajęcia fragmentów i spięcia wykonawców (`adapter_modul_studio_agenci.go`).
	ZajmijFragment(ctx context.Context, z shared.StudioAgentsClaimRequest) (shared.StudioAgentsClaimResponse, error)
	SpieciaWykonawcowDokumentu(ctx context.Context, z shared.StudioAgentsConflictsListRequest) (shared.StudioAgentsConflictsListResponse, error)
	ZwolnijFragment(ctx context.Context, z shared.StudioAgentsReleaseRequest) (shared.StudioAgentsReleaseResponse, error)
	NastawyWykonawcow(ctx context.Context, z shared.StudioAgentsSettingsGetRequest) (shared.StudioAgentsSettingsGetResponse, error)
	UstawNastawyWykonawcow(ctx context.Context, z shared.StudioAgentsSettingsSetRequest) (shared.StudioAgentsSettingsSetResponse, error)
	ZajeciaWykonawcow(ctx context.Context, z shared.StudioAgentsSlotsListRequest) (shared.StudioAgentsSlotsListResponse, error)

	// Schowek Studia (`adapter_modul_studio_schowek.go`).
	SkopiujDoSchowka(ctx context.Context, z shared.StudioClipboardCopyRequest) (shared.StudioClipboardCopyResponse, error)
	WklejZeSchowka(ctx context.Context, z shared.StudioClipboardPasteRequest) (shared.StudioClipboardPasteResponse, error)
	PochodzenieFragmentow(ctx context.Context, z shared.StudioProvenanceListRequest) (shared.StudioProvenanceListResponse, error)

	// Nastawy widoku dokumentu i pochodzenie fragmentów.
	NastawyWidoku(ctx context.Context, z shared.StudioViewGetRequest) (shared.StudioViewGetResponse, error)
	UstawWidok(ctx context.Context, z shared.StudioViewSetRequest) (shared.StudioViewSetResponse, error)
}

// zarejestrujStudio wpina komendy modułu Studio.
func zarejestrujStudio(r *Rejestr, m Studio, e *emiter) {
	if r == nil || m == nil {
		return
	}

	r.Zarejestruj(shared.CommandStudioDocumentOpen, obsluz(m.OtworzDokument))
	r.Zarejestruj(shared.CommandStudioDiffCompare, obsluz(m.Porownaj))

	// Operacja kontekstowa zmienia od teraz TREŚĆ dokumentu — wpisuje wynik
	// modelu jako zmianę śledzoną — więc rozgłasza zmianę tak samo jak zapis.
	// Odpowiedź komendy dokumentu nie niesie (kontrakt oddaje wynik i propozycję),
	// więc dokument po zmianie czytamy tą samą drogą, którą czyta go okno:
	// `document.open`. Odczyt bez powodzenia gasi samo rozgłoszenie, a nie
	// operację — wynik jest już zapisany i odmowa tutaj byłaby nieprawdą.
	r.Zarejestruj(shared.CommandStudioContextualOp,
		obsluz(func(ctx context.Context, z shared.StudioContextualOpRequest) (shared.StudioContextualOpResponse, error) {
			odpowiedz, err := m.OperacjaKontekstowa(ctx, z)
			if err != nil {
				return odpowiedz, err
			}
			po, bladOdczytu := m.OtworzDokument(ctx, shared.StudioDocumentOpenRequest{
				WindowId: z.WindowId, DocumentId: &z.DocumentId,
			})
			if bladOdczytu == nil {
				e.dokumentStudio(shared.ChangeKindUpdated, po.Document)
			}
			return odpowiedz, nil
		}))
	r.Zarejestruj(shared.CommandStudioRepositoryList, obsluz(m.Wersje))

	r.Zarejestruj(shared.CommandStudioDocumentSave,
		obsluz(func(ctx context.Context, z shared.StudioDocumentSaveRequest) (shared.StudioDocumentSaveResponse, error) {
			odpowiedz, err := m.ZapiszDokument(ctx, z)
			if err == nil {
				e.dokumentStudio(shared.ChangeKindUpdated, odpowiedz.Document)
			}
			return odpowiedz, err
		}))

	// ── Ingest/OCR Panel ──────────────────────────────────────────────────────
	// Kolejka i jej odczyt zdarzenia nie mają: pozycja czeka, a nie zmienia
	// dokumentu. Rozpoznanie i korekta zmieniają wyłącznie pozycję kolejki,
	// więc rozgłaszają zmianę kolejki, nie zmianę dokumentu. Dopiero przyjęcie
	// pozycji zakłada dokument i to ono rozgłasza `studio.document.changed`.
	r.Zarejestruj(shared.CommandStudioIngestQueueAdd, obsluz(m.DolozDoKolejki))
	r.Zarejestruj(shared.CommandStudioIngestQueueList, obsluz(m.KolejkaWczytywania))
	r.Zarejestruj(shared.CommandStudioIngestDeviceList, obsluz(m.UrzadzeniaWejsciowe))

	r.Zarejestruj(shared.CommandStudioIngestRecognize,
		obsluz(func(ctx context.Context, z shared.StudioIngestRecognizeRequest) (shared.StudioIngestRecognizeResponse, error) {
			odpowiedz, err := m.Rozpoznaj(ctx, z)
			if err == nil {
				e.pozycjaWczytywania(odpowiedz.Item)
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandStudioIngestCorrectionSet,
		obsluz(func(ctx context.Context, z shared.StudioIngestCorrectionSetRequest) (shared.StudioIngestCorrectionSetResponse, error) {
			odpowiedz, err := m.PoprawRozpoznanie(ctx, z)
			if err == nil {
				e.pozycjaWczytywania(odpowiedz.Item)
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandStudioIngestItemAccept,
		obsluz(func(ctx context.Context, z shared.StudioIngestItemAcceptRequest) (shared.StudioIngestItemAcceptResponse, error) {
			odpowiedz, err := m.PrzyjmijPozycje(ctx, z)
			if err == nil {
				e.dokumentStudio(shared.ChangeKindCreated, odpowiedz.Document)
			}
			return odpowiedz, err
		}))

	// ── Komentarze, adnotacje, śledzenie zmian ────────────────────────────────
	// Odczyty i zapisy komentarzy nie ruszają treści dokumentu, więc zdarzenia
	// nie mają. Obie decyzje — o zmianach śledzonych i o propozycji — treść
	// zmieniają i dlatego rozgłaszają zmianę dokumentu.
	r.Zarejestruj(shared.CommandStudioCommentAdd, obsluz(m.DodajKomentarz))
	r.Zarejestruj(shared.CommandStudioCommentList, obsluz(m.Komentarze))
	r.Zarejestruj(shared.CommandStudioCommentResolve, obsluz(m.RozstrzygnijKomentarz))
	r.Zarejestruj(shared.CommandStudioAnnotationAdd, obsluz(m.DodajAdnotacje))
	r.Zarejestruj(shared.CommandStudioAnnotationList, obsluz(m.Adnotacje))
	r.Zarejestruj(shared.CommandStudioTrackingSet, obsluz(m.UstawSledzenie))
	r.Zarejestruj(shared.CommandStudioTrackingList, obsluz(m.ZmianySledzone))

	r.Zarejestruj(shared.CommandStudioTrackingDecide,
		obsluz(func(ctx context.Context, z shared.StudioTrackingDecideRequest) (shared.StudioTrackingDecideResponse, error) {
			odpowiedz, err := m.RozstrzygnijZmiany(ctx, z)
			if err == nil {
				e.dokumentStudio(shared.ChangeKindUpdated, odpowiedz.Document)
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandStudioProposalDecide,
		obsluz(func(ctx context.Context, z shared.StudioProposalDecideRequest) (shared.StudioProposalDecideResponse, error) {
			odpowiedz, err := m.RozstrzygnijPropozycje(ctx, z)
			if err == nil {
				e.dokumentStudio(shared.ChangeKindUpdated, odpowiedz.Document)
			}
			return odpowiedz, err
		}))

	// ── Katalogi zasięgu i cechy dokumentu ────────────────────────────────────
	// Zapisy katalogowe nie ruszają dokumentu i zdarzenia nie mają. Zastosowanie
	// szablonu zakłada dokument, a przestawienie formatu zmienia jego cechę —
	// obie rozgłaszają zmianę dokumentu.
	r.Zarejestruj(shared.CommandStudioOperationSave, obsluz(m.ZapiszOperacje))
	r.Zarejestruj(shared.CommandStudioOperationList, obsluz(m.Operacje))
	r.Zarejestruj(shared.CommandStudioOperationDelete, obsluz(m.UsunOperacje))
	r.Zarejestruj(shared.CommandStudioChainSave, obsluz(m.ZapiszLancuch))
	r.Zarejestruj(shared.CommandStudioChainList, obsluz(m.Lancuchy))
	r.Zarejestruj(shared.CommandStudioChainRun, obsluz(m.UruchomLancuch))
	r.Zarejestruj(shared.CommandStudioTemplateList, obsluz(m.Szablony))
	r.Zarejestruj(shared.CommandStudioExportProfileSave, obsluz(m.ZapiszProfilWydania))
	r.Zarejestruj(shared.CommandStudioExportProfileList, obsluz(m.ProfileWydania))
	r.Zarejestruj(shared.CommandStudioVersionLabelSet, obsluz(m.UstawEtykieteWersji))

	r.Zarejestruj(shared.CommandStudioTemplateApply,
		obsluz(func(ctx context.Context, z shared.StudioTemplateApplyRequest) (shared.StudioTemplateApplyResponse, error) {
			odpowiedz, err := m.ZastosujSzablon(ctx, z)
			if err == nil {
				e.dokumentStudio(shared.ChangeKindCreated, odpowiedz.Document)
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandStudioDocumentFormatSet,
		obsluz(func(ctx context.Context, z shared.StudioDocumentFormatSetRequest) (shared.StudioDocumentFormatSetResponse, error) {
			odpowiedz, err := m.UstawFormatDokumentu(ctx, z)
			if err == nil {
				e.dokumentStudio(shared.ChangeKindUpdated, odpowiedz.Document)
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandStudioRepositoryRestore,
		obsluz(func(ctx context.Context, z shared.StudioRepositoryRestoreRequest) (shared.StudioRepositoryRestoreResponse, error) {
			odpowiedz, err := m.PrzywrocWersje(ctx, z)
			if err == nil {
				e.dokumentStudio(shared.ChangeKindUpdated, odpowiedz.Document)
			}
			return odpowiedz, err
		}))

	// ── Gałęzie, wydanie, podgląd, wsad i wczytywanie ─────────────────────────
	// Zdarzenie dokumentu rozgłaszają wyłącznie te czynności, które zmieniają
	// treść widoczną w edytorze: założenie gałęzi (otwiera ją jako treść
	// bieżącą), scalenie (gdy doszło do skutku) i osadzenie zasobu. Wydania,
	// wyrysy, wyszukiwanie i wsad treści dokumentu nie ruszają — wsad zakłada
	// propozycje, a propozycja staje się treścią dopiero decyzją Operatora.
	r.Zarejestruj(shared.CommandStudioBranchList, obsluz(m.Galezie))
	r.Zarejestruj(shared.CommandStudioVersionReferenceCreate, obsluz(m.UtworzOdwolanieWersji))
	r.Zarejestruj(shared.CommandStudioRepositoryExport, obsluz(m.WydajRepozytorium))
	r.Zarejestruj(shared.CommandStudioPackageExport, obsluz(m.WydajPaczke))
	r.Zarejestruj(shared.CommandStudioDiffReportExport, obsluz(m.WydajRaportRoznicy))
	r.Zarejestruj(shared.CommandStudioPreviewRender, obsluz(m.WyrenderujPodglad))
	r.Zarejestruj(shared.CommandStudioDiffVisual, obsluz(m.PorownajWizualnie))
	r.Zarejestruj(shared.CommandStudioBatchRun, obsluz(m.UruchomWsad))
	r.Zarejestruj(shared.CommandStudioSearchSemantic, obsluz(m.WyszukajZnaczeniowo))
	r.Zarejestruj(shared.CommandStudioDiffSource, obsluz(m.PorownajZeZrodlem))
	r.Zarejestruj(shared.CommandStudioIngestDeviceScan, obsluz(m.SkanujUrzadzenie))

	r.Zarejestruj(shared.CommandStudioIngestUrl,
		obsluz(func(ctx context.Context, z shared.StudioIngestUrlRequest) (shared.StudioIngestUrlResponse, error) {
			odpowiedz, err := m.WczytajZAdresu(ctx, z)
			if err == nil {
				e.pozycjaWczytywania(odpowiedz.Item)
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandStudioBranchCreate,
		obsluz(func(ctx context.Context, z shared.StudioBranchCreateRequest) (shared.StudioBranchCreateResponse, error) {
			odpowiedz, err := m.ZalozGalaz(ctx, z)
			if err == nil {
				e.dokumentStudio(shared.ChangeKindUpdated, odpowiedz.Document)
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandStudioBranchMerge,
		obsluz(func(ctx context.Context, z shared.StudioBranchMergeRequest) (shared.StudioBranchMergeResponse, error) {
			odpowiedz, err := m.ScalGalezie(ctx, z)
			// Scalenie zatrzymane konfliktem NIE rozgłasza zmiany dokumentu:
			// dokument został taki, jaki był, a zdarzenie kazałoby oknu
			// przeładować treść, która się nie zmieniła.
			if err == nil && odpowiedz.Merged && odpowiedz.Document != nil {
				e.dokumentStudio(shared.ChangeKindUpdated, *odpowiedz.Document)
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandStudioAssetEmbed,
		obsluz(func(ctx context.Context, z shared.StudioAssetEmbedRequest) (shared.StudioAssetEmbedResponse, error) {
			odpowiedz, err := m.OsadzZasob(ctx, z)
			if err == nil {
				e.dokumentStudio(shared.ChangeKindUpdated, odpowiedz.Document)
			}
			return odpowiedz, err
		}))

	// ── Postać dokumentu, formatowanie, style i blokady ───────────────────────
	// Zdarzenia `studio.document.changed` te czynności NIE rozgłaszają, choć
	// zmieniają dokument: każda oddaje postać po zmianie wprost w odpowiedzi
	// (`Form`), a okno pracy z dokumentem tę odpowiedź już ma. Zdarzenie kazałoby
	// mu przeładować to samo drugą drogą i przy pisaniu litera po literze byłoby
	// przeładowaniem na każde naciśnięcie klawisza.
	//
	// Wyjątek jest jeden: `document.form.save` oddaje także DOKUMENT wraz
	// z wersją, więc rozgłasza zmianę tak samo jak `document.save` — po nim
	// odświeżają się wykazy dokumentów i historia wersji, a nie tylko powierzchnia.
	r.Zarejestruj(shared.CommandStudioDocumentFormGet, obsluz(m.PostacDokumentu))
	r.Zarejestruj(shared.CommandStudioTextGet, obsluz(m.TrescFragmentu))
	r.Zarejestruj(shared.CommandStudioTextEdit, obsluz(m.ZmienTresc))
	r.Zarejestruj(shared.CommandStudioFormatCharacterGet, obsluz(m.PostacZnaku))
	r.Zarejestruj(shared.CommandStudioFormatCharacterSet, obsluz(m.UstawPostacZnaku))
	r.Zarejestruj(shared.CommandStudioFormatParagraphGet, obsluz(m.PostacAkapitu))
	r.Zarejestruj(shared.CommandStudioFormatParagraphSet, obsluz(m.UstawPostacAkapitu))
	r.Zarejestruj(shared.CommandStudioFormatClear, obsluz(m.CzyscPostac))
	r.Zarejestruj(shared.CommandStudioFormatCaseSet, obsluz(m.UstawWielkoscLiter))
	r.Zarejestruj(shared.CommandStudioFormatPainterCopy, obsluz(m.ZabierzPostac))
	r.Zarejestruj(shared.CommandStudioFormatPainterApply, obsluz(m.PolozPostac))
	r.Zarejestruj(shared.CommandStudioFormatSimilarSelect, obsluz(m.ZaznaczPodobne))
	r.Zarejestruj(shared.CommandStudioFormatReplace, obsluz(m.ZamienZPostacia))
	r.Zarejestruj(shared.CommandStudioStyleList, obsluz(m.StyleDokumentu))
	r.Zarejestruj(shared.CommandStudioStyleSave, obsluz(m.ZapiszStyl))
	r.Zarejestruj(shared.CommandStudioStyleApply, obsluz(m.ZastosujStyl))
	r.Zarejestruj(shared.CommandStudioStyleDelete, obsluz(m.UsunStyl))
	r.Zarejestruj(shared.CommandStudioLockAdd, obsluz(m.ZalozBlokade))
	r.Zarejestruj(shared.CommandStudioLockList, obsluz(m.Blokady))
	r.Zarejestruj(shared.CommandStudioLockRemove, obsluz(m.ZdejmijBlokade))

	// ── Postać dokumentu: strona, listy, symbole, tabele, obiekty, aparat ─────
	// Te czynności oddają postać po zmianie wprost w odpowiedzi, więc zdarzenia
	// dokumentu nie rozgłaszają — z jednym wyjątkiem: czynność, która ZAKŁADA
	// dokument albo podmienia dokument zastany, oddaje go w odpowiedzi i wtedy
	// zdarzenie jedzie, bo odświeżają się wykazy dokumentów, nie sama powierzchnia.
	r.Zarejestruj(shared.CommandStudioPageBreakInsert, obsluz(m.WstawPodzial))
	r.Zarejestruj(shared.CommandStudioPageEnvelopeSet, obsluz(m.UstawNadrukKoperty))
	r.Zarejestruj(shared.CommandStudioPageHeaderfooterGet, obsluz(m.NaglowkiIStopki))
	r.Zarejestruj(shared.CommandStudioPageHeaderfooterSet, obsluz(m.UstawNaglowekIStopke))
	r.Zarejestruj(shared.CommandStudioPageNumberingSet, obsluz(m.UstawNumeracjeStron))
	r.Zarejestruj(shared.CommandStudioPagePaperList, obsluz(m.NosnikiStrony))
	r.Zarejestruj(shared.CommandStudioPageSetupGet, obsluz(m.NastawyStrony))
	r.Zarejestruj(shared.CommandStudioPageSetupSet, obsluz(m.UstawNastawyStrony))
	r.Zarejestruj(shared.CommandStudioPageWatermarkSet, obsluz(m.UstawZnakWodny))
	r.Zarejestruj(shared.CommandStudioRulerTabstopSet, obsluz(m.UstawTabulatorLinijki))
	r.Zarejestruj(shared.CommandStudioSectionDelete, obsluz(m.UsunSekcjeDokumentu))
	r.Zarejestruj(shared.CommandStudioSectionList, obsluz(m.SekcjeDokumentu))
	r.Zarejestruj(shared.CommandStudioSectionSave, obsluz(m.ZapiszSekcjeDokumentu))
	r.Zarejestruj(shared.CommandStudioListApply, obsluz(m.ZastosujListe))
	r.Zarejestruj(shared.CommandStudioListBulletSet, obsluz(m.UstawPunktatorListy))
	r.Zarejestruj(shared.CommandStudioListLevelIndent, obsluz(m.PrzestawPoziomListy))
	r.Zarejestruj(shared.CommandStudioListNumberingSet, obsluz(m.UstawNumeracjeListy))
	r.Zarejestruj(shared.CommandStudioListRestart, obsluz(m.WznowNumeracjeListy))
	r.Zarejestruj(shared.CommandStudioSymbolAutoreplaceList, obsluz(m.ZasadyAutozamiany))
	r.Zarejestruj(shared.CommandStudioSymbolAutoreplaceSet, obsluz(m.UstawZasadeAutozamiany))
	r.Zarejestruj(shared.CommandStudioSymbolInsert, obsluz(m.WstawZnak))
	r.Zarejestruj(shared.CommandStudioSymbolList, obsluz(m.TabliceZnakow))
	r.Zarejestruj(shared.CommandStudioTableConvert, obsluz(m.ZamienTabeleITekst))
	r.Zarejestruj(shared.CommandStudioTableFormatSet, obsluz(m.UstawPostacTabeli))
	r.Zarejestruj(shared.CommandStudioTableInsert, obsluz(m.WstawTabele))
	r.Zarejestruj(shared.CommandStudioTableList, obsluz(m.WykazTabel))
	r.Zarejestruj(shared.CommandStudioTableSort, obsluz(m.SortujTabele))
	r.Zarejestruj(shared.CommandStudioTableStructureEdit, obsluz(m.ZmienBudoweTabeli))
	r.Zarejestruj(shared.CommandStudioObjectFormatSet, obsluz(m.UstawPostacObiektu))
	r.Zarejestruj(shared.CommandStudioObjectInsert, obsluz(m.WstawObiekt))
	r.Zarejestruj(shared.CommandStudioObjectList, obsluz(m.WykazObiektow))
	r.Zarejestruj(shared.CommandStudioObjectRemove, obsluz(m.UsunObiekt))
	r.Zarejestruj(shared.CommandStudioApparatusInsert, obsluz(m.WstawElementAparatu))
	r.Zarejestruj(shared.CommandStudioApparatusList, obsluz(m.WykazAparatu))
	r.Zarejestruj(shared.CommandStudioApparatusRefresh, obsluz(m.OdswiezAparat))
	r.Zarejestruj(shared.CommandStudioApparatusRemove, obsluz(m.UsunElementAparatuDokumentu))
	r.Zarejestruj(shared.CommandStudioFieldInsert, obsluz(m.WstawPole))
	r.Zarejestruj(shared.CommandStudioFieldList, obsluz(m.WykazPol))
	r.Zarejestruj(shared.CommandStudioFieldRefresh, obsluz(m.OdswiezPola))
	r.Zarejestruj(shared.CommandStudioDocumentImageImport, obsluz(m.WniesObraz))
	r.Zarejestruj(shared.CommandStudioInsertFromLibrary, obsluz(m.WniesZBiblioteki))
	r.Zarejestruj(shared.CommandStudioInsertFromWeb, obsluz(m.WniesZeSieci))
	r.Zarejestruj(shared.CommandStudioDocumentExportBatch, obsluz(m.WydajWsadowo))
	r.Zarejestruj(shared.CommandStudioDocumentExportFormat, obsluz(m.WydajDoFormatu))
	r.Zarejestruj(shared.CommandStudioTemplateDelete, obsluz(m.UsunSzablonPisma))
	r.Zarejestruj(shared.CommandStudioTemplateExport, obsluz(m.OddajSzablonDoPliku))
	r.Zarejestruj(shared.CommandStudioTemplateFieldList, obsluz(m.PolaSzablonu))
	r.Zarejestruj(shared.CommandStudioTemplateFieldSet, obsluz(m.UstawPoleSzablonu))
	r.Zarejestruj(shared.CommandStudioTemplateImport, obsluz(m.WniesSzablonZPliku))
	r.Zarejestruj(shared.CommandStudioTemplateSave, obsluz(m.ZapiszSzablonPisma))

	r.Zarejestruj(shared.CommandStudioDocumentCopy,
		obsluz(func(ctx context.Context, z shared.StudioDocumentCopyRequest) (shared.StudioDocumentCopyResponse, error) {
			odpowiedz, err := m.SkopiujDokument(ctx, z)
			if err == nil {
				e.dokumentStudio(shared.ChangeKindCreated, odpowiedz.Document)
			}
			return odpowiedz, err
		}))
	r.Zarejestruj(shared.CommandStudioDocumentCreate,
		obsluz(func(ctx context.Context, z shared.StudioDocumentCreateRequest) (shared.StudioDocumentCreateResponse, error) {
			odpowiedz, err := m.ZalozDokument(ctx, z)
			if err == nil {
				e.dokumentStudio(shared.ChangeKindCreated, odpowiedz.Document)
			}
			return odpowiedz, err
		}))
	r.Zarejestruj(shared.CommandStudioDocumentImportFile,
		obsluz(func(ctx context.Context, z shared.StudioDocumentImportFileRequest) (shared.StudioDocumentImportFileResponse, error) {
			odpowiedz, err := m.WniesPlikDoEdytora(ctx, z)
			if err == nil {
				e.dokumentStudio(shared.ChangeKindCreated, odpowiedz.Document)
			}
			return odpowiedz, err
		}))
	r.Zarejestruj(shared.CommandStudioDocumentImportPdf,
		obsluz(func(ctx context.Context, z shared.StudioDocumentImportPdfRequest) (shared.StudioDocumentImportPdfResponse, error) {
			odpowiedz, err := m.WniesPdfDoEdytora(ctx, z)
			if err == nil {
				e.dokumentStudio(shared.ChangeKindCreated, odpowiedz.Document)
			}
			return odpowiedz, err
		}))
	r.Zarejestruj(shared.CommandStudioDocumentSaveAs,
		obsluz(func(ctx context.Context, z shared.StudioDocumentSaveAsRequest) (shared.StudioDocumentSaveAsResponse, error) {
			odpowiedz, err := m.ZapiszDokumentPodNazwa(ctx, z)
			if err == nil {
				e.dokumentStudio(shared.ChangeKindCreated, odpowiedz.Document)
			}
			return odpowiedz, err
		}))
	r.Zarejestruj(shared.CommandStudioTemplateFill,
		obsluz(func(ctx context.Context, z shared.StudioTemplateFillRequest) (shared.StudioTemplateFillResponse, error) {
			odpowiedz, err := m.WypelnijSzablon(ctx, z)
			if err == nil {
				e.dokumentStudio(shared.ChangeKindUpdated, odpowiedz.Document)
			}
			return odpowiedz, err
		}))

	// ── Kontrola pracy: dziennik, zmiany modelu, różnica, autozapis, znakowanie,
	// zajęcia wykonawców, schowek, widok ─────────────────────────────────────
	// Zdarzenie dokumentu rozgłaszają tylko te czynności, które PODMIENIAJĄ treść
	// widoczną w oknie i oddają dokument w odpowiedzi: cofnięcie i ponowienie
	// czynności, cofnięcie zmian modelu, przyjęcie propozycji, przeniesienie
	// fragmentu różnicy, przywrócenie kopii i powrót do wersji założycielskiej.
	// Odczyty, nastawy i wykazy treści nie ruszają.
	r.Zarejestruj(shared.CommandStudioJournalList, obsluz(m.DziennikCzynnosci))
	r.Zarejestruj(shared.CommandStudioModelChangesList, obsluz(m.ZmianyModelu))
	r.Zarejestruj(shared.CommandStudioModelChangesNavigate, obsluz(m.PrzeskocDoZmianyModelu))
	r.Zarejestruj(shared.CommandStudioDiffFormCompare, obsluz(m.PorownajPostac))
	r.Zarejestruj(shared.CommandStudioAutosaveGet, obsluz(m.NastawyAutozapisu))
	r.Zarejestruj(shared.CommandStudioAutosaveRun, obsluz(m.WykonajAutozapis))
	r.Zarejestruj(shared.CommandStudioAutosaveSet, obsluz(m.UstawAutozapis))
	r.Zarejestruj(shared.CommandStudioBackupCreate, obsluz(m.ZalozKopieZapasowa))
	r.Zarejestruj(shared.CommandStudioBackupList, obsluz(m.KopieDokumentu))
	r.Zarejestruj(shared.CommandStudioVersionSeriesList, obsluz(m.WersjeWSzeregach))
	r.Zarejestruj(shared.CommandStudioMarkupAdd, obsluz(m.DodajZnakowanie))
	r.Zarejestruj(shared.CommandStudioMarkupList, obsluz(m.Znakowania))
	r.Zarejestruj(shared.CommandStudioMarkupRemove, obsluz(m.ZdejmijZnakowanie))
	r.Zarejestruj(shared.CommandStudioMarkupTypeDelete, obsluz(m.UsunRodzajZnacznika))
	r.Zarejestruj(shared.CommandStudioMarkupTypeList, obsluz(m.RodzajeZnacznika))
	r.Zarejestruj(shared.CommandStudioMarkupTypeSave, obsluz(m.ZapiszRodzajZnacznika))
	r.Zarejestruj(shared.CommandStudioAgentsClaim, obsluz(m.ZajmijFragment))
	r.Zarejestruj(shared.CommandStudioAgentsConflictsList, obsluz(m.SpieciaWykonawcowDokumentu))
	r.Zarejestruj(shared.CommandStudioAgentsRelease, obsluz(m.ZwolnijFragment))
	r.Zarejestruj(shared.CommandStudioAgentsSettingsGet, obsluz(m.NastawyWykonawcow))
	r.Zarejestruj(shared.CommandStudioAgentsSettingsSet, obsluz(m.UstawNastawyWykonawcow))
	r.Zarejestruj(shared.CommandStudioAgentsSlotsList, obsluz(m.ZajeciaWykonawcow))
	r.Zarejestruj(shared.CommandStudioProvenanceList, obsluz(m.PochodzenieFragmentow))
	r.Zarejestruj(shared.CommandStudioViewGet, obsluz(m.NastawyWidoku))
	r.Zarejestruj(shared.CommandStudioViewSet, obsluz(m.UstawWidok))

	r.Zarejestruj(shared.CommandStudioBackupRestore,
		obsluz(func(ctx context.Context, z shared.StudioBackupRestoreRequest) (shared.StudioBackupRestoreResponse, error) {
			odpowiedz, err := m.PrzywrocKopie(ctx, z)
			if err == nil {
				e.dokumentStudio(shared.ChangeKindUpdated, odpowiedz.Document)
			}
			return odpowiedz, err
		}))
	r.Zarejestruj(shared.CommandStudioDiffHunkApply,
		obsluz(func(ctx context.Context, z shared.StudioDiffHunkApplyRequest) (shared.StudioDiffHunkApplyResponse, error) {
			odpowiedz, err := m.PrzeniesFragmentRoznicy(ctx, z)
			if err == nil {
				e.dokumentStudio(shared.ChangeKindUpdated, odpowiedz.Document)
			}
			return odpowiedz, err
		}))
	r.Zarejestruj(shared.CommandStudioJournalRedo,
		obsluz(func(ctx context.Context, z shared.StudioJournalRedoRequest) (shared.StudioJournalRedoResponse, error) {
			odpowiedz, err := m.PonowCzynnosc(ctx, z)
			if err == nil {
				e.dokumentStudio(shared.ChangeKindUpdated, odpowiedz.Document)
			}
			return odpowiedz, err
		}))
	r.Zarejestruj(shared.CommandStudioJournalRevert,
		obsluz(func(ctx context.Context, z shared.StudioJournalRevertRequest) (shared.StudioJournalRevertResponse, error) {
			odpowiedz, err := m.CofnijCzynnosc(ctx, z)
			if err == nil {
				e.dokumentStudio(shared.ChangeKindUpdated, odpowiedz.Document)
			}
			return odpowiedz, err
		}))
	r.Zarejestruj(shared.CommandStudioMarkupDecide,
		obsluz(func(ctx context.Context, z shared.StudioMarkupDecideRequest) (shared.StudioMarkupDecideResponse, error) {
			odpowiedz, err := m.RozstrzygnijZnakowanie(ctx, z)
			if err == nil {
				e.dokumentStudio(shared.ChangeKindUpdated, odpowiedz.Document)
			}
			return odpowiedz, err
		}))
	r.Zarejestruj(shared.CommandStudioModelChangesRevert,
		obsluz(func(ctx context.Context, z shared.StudioModelChangesRevertRequest) (shared.StudioModelChangesRevertResponse, error) {
			odpowiedz, err := m.CofnijZmianyModelu(ctx, z)
			if err == nil {
				e.dokumentStudio(shared.ChangeKindUpdated, odpowiedz.Document)
			}
			return odpowiedz, err
		}))
	r.Zarejestruj(shared.CommandStudioVersionRestoreInitial,
		obsluz(func(ctx context.Context, z shared.StudioVersionRestoreInitialRequest) (shared.StudioVersionRestoreInitialResponse, error) {
			odpowiedz, err := m.PrzywrocWersjeZalozycielska(ctx, z)
			if err == nil {
				e.dokumentStudio(shared.ChangeKindUpdated, odpowiedz.Document)
			}
			return odpowiedz, err
		}))

	// Schowek dokumentu zdarzenia NIE rozgłasza, choć wycięcie i wklejenie zmieniają
	// treść: kontrakt tych dwóch odpowiedzi dokumentu nie niesie, a doczytanie go
	// drugą drogą tylko po to, żeby rozgłosić zmianę, kosztowałoby dwa odczyty na
	// każde wklejenie. Okno pracy dostaje w odpowiedzi postać i tym się odświeża.
	r.Zarejestruj(shared.CommandStudioClipboardCopy, obsluz(m.SkopiujDoSchowka))
	r.Zarejestruj(shared.CommandStudioClipboardPaste, obsluz(m.WklejZeSchowka))

	r.Zarejestruj(shared.CommandStudioDocumentFormSave,
		obsluz(func(ctx context.Context, z shared.StudioDocumentFormSaveRequest) (shared.StudioDocumentFormSaveResponse, error) {
			odpowiedz, err := m.ZapiszPostacDokumentu(ctx, z)
			if err == nil {
				e.dokumentStudio(shared.ChangeKindUpdated, odpowiedz.Document)
			}
			return odpowiedz, err
		}))

	// Podpis wykonawcy wchodzi na drogę WSZYSTKICH wpisanych wyżej komend — na
	// końcu, bo owija to, co w rejestrze już stoi.
	podpisWykonawcyStudia(r)
}

// podpisWykonawcyStudia owija w rejestrze każdą komendę rodziny `studio.*`
// czytaniem podpisu wykonawcy z ładunku i wstawieniem go do kontekstu.
//
// ── Dlaczego w rejestrze, a nie w obsługiwaczach ─────────────────────────────
// To ten sam rachunek, którym zapora blokad stanęła w drzwiach, a nie przy
// każdym stoliku. Tożsamość wykonawcy niesie żądanie, a odkłada ją jedna droga
// wyjścia czynności postaci (`postacZakoncz`) — tej drogi żądanie nie widzi.
// Przełożenie podpisu przez trzydzieści osiem sygnatur czynności postaci
// znaczyłoby trzydzieści osiem miejsc do pominięcia przez pomyłkę, a pominięcie
// nie jest tu widoczne: zmiana zapisuje się dalej, tylko podpisana
// „nienazwanym". Owinięcie rejestru obejmuje wszystkie te komendy jednym
// warunkiem i obejmuje też te, których jeszcze nikt nie napisał.
//
// ── Dlaczego wpięcie nie odmawia i nie zmienia ładunku ──────────────────────
// Ładunek innego kształtu nie jest usterką — nie każda komenda Studia niesie
// podpis, a ta, która go nie niesie, jedzie dalej jako czynność Operatora.
// Wpięcie NIE dotyka ładunku ani odpowiedzi: dokłada wyłącznie wiedzę o tym, kto
// woła, więc nie ma jak zmienić skutku komendy.
func podpisWykonawcyStudia(r *Rejestr) {
	if r == nil || r.wpisy == nil {
		return
	}
	for nazwa, obsluga := range r.wpisy {
		if !strings.HasPrefix(string(nazwa), "studio.") {
			continue
		}
		r.wpisy[nazwa] = podpisemWykonawcy(obsluga)
	}
}

// podpisemWykonawcy składa obsługiwacza, który zna wykonawcę z żądania.
//
// ── Dlaczego wpięcie STEMPLUJE pole `author` w ładunku ──────────────────────
// Rozpoznanie wykonawcy jedzie dalej kontekstem, ale czynności postaci
// dokumentu rozstrzygają autora z POLA ŻĄDANIA (`postacAutor`), bo tak stanowi
// kontrakt: „narzędzie modelu podaje autora wprost". Wykonawca, który tego pola
// nie poda, byłby wtedy zapisany jako Operator — a wówczas jego zmiana nie
// odkłada się jako zmiana śledzona i NIE DA SIĘ JEJ PODŚWIETLIĆ przełącznikiem
// „pokaż wszystko, co zrobił model". Właściciel nazwał to wprost USTERKĄ do
// naprawy, nie ograniczeniem do zgłoszenia.
//
// Naprawa stoi tutaj, a nie w trzydziestu ośmiu czynnościach postaci: gdy fakt
// gniazda mówi „to wykonawca", wpięcie dopisuje `author: model` do ładunku,
// zanim ładunek zobaczy obsługiwacz. Dzięki temu każda droga — także te, których
// jeszcze nikt nie napisał — czyta autora prawdziwego, a nie zatajonego.
//
// Stempel idzie WYŁĄCZNIE w jedną stronę: podnosi Operatora do wykonawcy, nigdy
// odwrotnie. Żądanie podpisane `author: uzytkownik` przyszłe z gniazda serwera
// narzędzi jest twierdzeniem modelu o sobie, nie faktem — i dlatego przegrywa
// z gniazdem. Tożsamość agenta stempluje się tylko wtedy, gdy żądanie jej nie
// podało: podpis wie, KTÓRY wykonawca woła, a gniazdo tego nie wie.
func podpisemWykonawcy(obsluga Obsluga) Obsluga {
	return func(ctx context.Context, z protocol.Request) protocol.Odpowiedz {
		var podpis kontrolaPodpisZadania
		if err := json.Unmarshal(z.Ladunek, &podpis); err != nil {
			return obsluga(ctx, z)
		}
		wykonawca := kontrolaRozpoznajWykonawce(ctx, podpis)
		if wykonawca.czyWykonawca() &&
			(podpis.Author == nil || *podpis.Author != shared.StudioAuthorModel) {

			ostemplowany, err := podpisStempelAutora(z.Ladunek)
			if err != nil {
				// Ładunku nie da się ostemplować — a bez stempla zmiana wykonawcy
				// zapisałaby się jako zmiana Operatora i zniknęłaby z
				// podświetlenia. Cisza byłaby tu gorsza niż odmowa.
				return porazka(kontrolaBladZaplecza(
					"żądania wykonawcy nie da się podpisać autorem: " + err.Error()))
			}
			z.Ladunek = ostemplowany
		}
		return obsluga(kontrolaZapiszWykonawce(ctx, wykonawca), z)
	}
}

// podpisStempelAutora dopisuje `author: model` do ładunku żądania.
//
// Przez mapę, nie przez strukturę: wpięcie nie zna kształtów żądań wszystkich
// odcinków, a złożenie ładunku ze znanej mu struktury zgubiłoby każde pole,
// o którym nie wie.
func podpisStempelAutora(ladunek json.RawMessage) (json.RawMessage, error) {
	pola := map[string]json.RawMessage{}
	if err := json.Unmarshal(ladunek, &pola); err != nil {
		return nil, err
	}
	zapis, err := json.Marshal(shared.StudioAuthorModel)
	if err != nil {
		return nil, err
	}
	pola["author"] = zapis
	return json.Marshal(pola)
}

// dokumentStudio rozgłasza zmianę dokumentu Studio. Rodzaj zmiany podaje
// wołający: zapis i przywrócenie wersji zmieniają dokument zastany („updated"),
// a przyjęcie pozycji cyfryzacji zakłada dokument nowy („created") — i to jest
// jedyne miejsce, w którym Studio zakłada dokument z treścią od razu.
// pozycjaWczytywania rozgłasza zmianę stanu pozycji kolejki cyfryzacji.
// Ingest/OCR Panel odświeża się nim zamiast odpytywać rdzeń w pętli — a przy
// wsadzie wielostronicowym pytanie w pętli byłoby pytaniem o kilkadziesiąt
// pozycji naraz.
func (e *emiter) pozycjaWczytywania(pozycja shared.StudioIngestItem) {
	e.wyslij(shared.EventStudioIngestChanged, "", shared.StudioIngestChangedEvent{
		WindowId: pozycja.WindowId, Item: pozycja,
	})
}

func (e *emiter) dokumentStudio(zmiana shared.ChangeKind, dokument shared.StudioDocument) {
	e.wyslij(shared.EventStudioDocumentChanged, "", shared.StudioDocumentChangedEvent{
		Change: zmiana, Document: dokument,
	})
}
