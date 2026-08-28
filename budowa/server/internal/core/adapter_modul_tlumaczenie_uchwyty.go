// Plik wpina komendy modułu Translate w port `Tlumaczenie`, przez który
// rejestr komend rdzenia dociera do adaptera modułu, wspólnego z pozostałymi
// plikami metod.
package core

import (
	"context"

	"danacoconsole/shared"
)

// Tlumaczenie jest portem całego modułu Translate, zbierającym wszystkie
// komendy jego adaptera w jeden interfejs.
type Tlumaczenie interface {
	UstawZrodlo(ctx context.Context, z shared.TranslateSourceSetRequest) (shared.TranslateSourceSetResponse, error)
	PodzielNaSegmenty(ctx context.Context, z shared.TranslateSourceSegmentRequest) (shared.TranslateSourceSegmentResponse, error)
	RozpoznajJezyk(ctx context.Context, z shared.TranslateSourceDetectRequest) (shared.TranslateSourceDetectResponse, error)
	DodajPanel(ctx context.Context, z shared.TranslateTargetAddRequest) (shared.TranslateTargetAddResponse, error)
	UstawTlumaczenie(ctx context.Context, z shared.TranslateTranslationSetRequest) (shared.TranslateTranslationSetResponse, error)
	TlumaczZwrotnie(ctx context.Context, z shared.TranslateBacktranslationRunRequest) (shared.TranslateBacktranslationRunResponse, error)
	UstawTerminy(ctx context.Context, z shared.TranslateGlossarySetRequest) (shared.TranslateGlossarySetResponse, error)
	ZastosujSlownik(ctx context.Context, z shared.TranslateGlossaryApplyRequest) (shared.TranslateGlossaryApplyResponse, error)
	WystapieniaTerminow(ctx context.Context, z shared.TranslateGlossaryOccurrencesRequest) (shared.TranslateGlossaryOccurrencesResponse, error)
	ImportujSlownik(ctx context.Context, z shared.TranslateGlossaryImportRequest) (shared.TranslateGlossaryImportResponse, error)
	EksportujSlownik(ctx context.Context, z shared.TranslateGlossaryExportRequest) (shared.TranslateGlossaryExportResponse, error)
	Podpowiedzi(ctx context.Context, z shared.TranslateMemorySuggestRequest) (shared.TranslateMemorySuggestResponse, error)
	SprawdzJakosc(ctx context.Context, z shared.TranslateQualityCheckRequest) (shared.TranslateQualityCheckResponse, error)
	UstawTon(ctx context.Context, z shared.TranslatePanelToneSetRequest) (shared.TranslatePanelToneSetResponse, error)
	SyntezujMowe(ctx context.Context, z shared.TranslateSpeechSynthesizeRequest) (shared.TranslateSpeechSynthesizeResponse, error)
	EksportujPanel(ctx context.Context, z shared.TranslatePanelExportRequest) (shared.TranslatePanelExportResponse, error)

	// --- pamięć tłumaczeń jako byt Operatora (`*_pamiec_wykaz.go`) ---
	WykazPamieci(ctx context.Context, z shared.TranslateMemoryListRequest) (shared.TranslateMemoryListResponse, error)
	UstawWpisPamieci(ctx context.Context, z shared.TranslateMemorySetRequest) (shared.TranslateMemorySetResponse, error)
	UsunWpisPamieci(ctx context.Context, z shared.TranslateMemoryDeleteRequest) (shared.TranslateMemoryDeleteResponse, error)
	ImportujPamiec(ctx context.Context, z shared.TranslateMemoryImportRequest) (shared.TranslateMemoryImportResponse, error)
	EksportujPamiec(ctx context.Context, z shared.TranslateMemoryExportRequest) (shared.TranslateMemoryExportResponse, error)
	UtrzymajPamiec(ctx context.Context, z shared.TranslateMemoryMaintainRequest) (shared.TranslateMemoryMaintainResponse, error)
	TlumaczWstepnie(ctx context.Context, z shared.TranslateMemoryPretranslateRequest) (shared.TranslateMemoryPretranslateResponse, error)
	WyrownajTeksty(ctx context.Context, z shared.TranslateMemoryAlignRequest) (shared.TranslateMemoryAlignResponse, error)
	PolitykaPamieci(ctx context.Context, z shared.TranslateMemoryPolicyGetRequest) (shared.TranslateMemoryPolicyGetResponse, error)
	UstawPolitykePamieci(ctx context.Context, z shared.TranslateMemoryPolicySetRequest) (shared.TranslateMemoryPolicySetResponse, error)

	// --- segmentacja (`*_segmenty.go`) ---
	WykazRegulSegmentacji(ctx context.Context, z shared.TranslateSegmentationRulesListRequest) (shared.TranslateSegmentationRulesListResponse, error)
	UstawRegulySegmentacji(ctx context.Context, z shared.TranslateSegmentationRulesSetRequest) (shared.TranslateSegmentationRulesSetResponse, error)
	ScalSegmenty(ctx context.Context, z shared.TranslateSegmentMergeRequest) (shared.TranslateSegmentMergeResponse, error)
	PodzielSegment(ctx context.Context, z shared.TranslateSegmentSplitRequest) (shared.TranslateSegmentSplitResponse, error)

	// --- terminologia (`*_terminologia.go`) ---
	WykazTerminow(ctx context.Context, z shared.TranslateGlossaryListRequest) (shared.TranslateGlossaryListResponse, error)
	WyjmijTerminy(ctx context.Context, z shared.TranslateTermExtractRequest) (shared.TranslateTermExtractResponse, error)

	// --- korekta i spójność (`*_korekta.go`) ---
	SprawdzKorekte(ctx context.Context, z shared.TranslateProofreadRunRequest) (shared.TranslateProofreadRunResponse, error)
	ZastosujKorekte(ctx context.Context, z shared.TranslateProofreadApplyRequest) (shared.TranslateProofreadApplyResponse, error)
	SprawdzSpojnosc(ctx context.Context, z shared.TranslateConsistencyCheckRequest) (shared.TranslateConsistencyCheckResponse, error)

	// --- profile kontroli jakości i obieg (`*_profile_qa.go`) ---
	WykazProfiliQa(ctx context.Context, z shared.TranslateQaProfileListRequest) (shared.TranslateQaProfileListResponse, error)
	UstawProfilQa(ctx context.Context, z shared.TranslateQaProfileSetRequest) (shared.TranslateQaProfileSetResponse, error)
	UsunProfilQa(ctx context.Context, z shared.TranslateQaProfileDeleteRequest) (shared.TranslateQaProfileDeleteResponse, error)
	UstawZatwierdzenie(ctx context.Context, z shared.TranslateApprovalSetRequest) (shared.TranslateApprovalSetResponse, error)
	WykazZatwierdzen(ctx context.Context, z shared.TranslateApprovalListRequest) (shared.TranslateApprovalListResponse, error)

	// --- dokument (`*_dokument.go`) ---
	WczytajDokument(ctx context.Context, z shared.TranslateDocumentLoadRequest) (shared.TranslateDocumentLoadResponse, error)
	ZlozDokument(ctx context.Context, z shared.TranslateDocumentRenderRequest) (shared.TranslateDocumentRenderResponse, error)
	PorownajUklad(ctx context.Context, z shared.TranslateDocumentLayoutCompareRequest) (shared.TranslateDocumentLayoutCompareResponse, error)

	// --- lokalizacja oprogramowania (`*_lokalizacja.go`, `*_xliff.go`) ---
	WczytajZasobLokalizacji(ctx context.Context, z shared.TranslateResourceImportRequest) (shared.TranslateResourceImportResponse, error)
	WydajZasobLokalizacji(ctx context.Context, z shared.TranslateResourceExportRequest) (shared.TranslateResourceExportResponse, error)
	ZastosujFormyMnogie(ctx context.Context, z shared.TranslateResourcePluralApplyRequest) (shared.TranslateResourcePluralApplyResponse, error)
	UstawKontekstKlucza(ctx context.Context, z shared.TranslateResourceKeyContextSetRequest) (shared.TranslateResourceKeyContextSetResponse, error)
	WczytajXliff(ctx context.Context, z shared.TranslateXliffImportRequest) (shared.TranslateXliffImportResponse, error)

	// --- napisy i dubbing (`*_napisy.go`) ---
	WczytajNapisy(ctx context.Context, z shared.TranslateSubtitleImportRequest) (shared.TranslateSubtitleImportResponse, error)
	WydajNapisy(ctx context.Context, z shared.TranslateSubtitleExportRequest) (shared.TranslateSubtitleExportResponse, error)
	SprawdzTaktowanieNapisow(ctx context.Context, z shared.TranslateSubtitleTimingCheckRequest) (shared.TranslateSubtitleTimingCheckResponse, error)
	ZlozScenariuszDubbingu(ctx context.Context, z shared.TranslateDubbingScriptBuildRequest) (shared.TranslateDubbingScriptBuildResponse, error)

	// --- silniki, pivot i przebieg pakietowy (`*_silniki.go`) ---
	WykazProfiliSilnikow(ctx context.Context, z shared.TranslateEngineProfileListRequest) (shared.TranslateEngineProfileListResponse, error)
	UstawProfilSilnika(ctx context.Context, z shared.TranslateEngineProfileSetRequest) (shared.TranslateEngineProfileSetResponse, error)
	PorownajSilniki(ctx context.Context, z shared.TranslateEngineCompareRequest) (shared.TranslateEngineCompareResponse, error)
	PolitykaPivota(ctx context.Context, z shared.TranslatePivotPolicyGetRequest) (shared.TranslatePivotPolicyGetResponse, error)
	UstawPolitykePivota(ctx context.Context, z shared.TranslatePivotPolicySetRequest) (shared.TranslatePivotPolicySetResponse, error)
	UruchomPakiet(ctx context.Context, z shared.TranslateBatchRunRequest) (shared.TranslateBatchRunResponse, error)

	// --- wymiana zewnętrzna (`*_wymiana_zewnetrzna.go`) ---
	ZlozPakietPrzekazania(ctx context.Context, z shared.TranslateHandoffBuildRequest) (shared.TranslateHandoffBuildResponse, error)
	PrzyjmijPakietPrzekazania(ctx context.Context, z shared.TranslateHandoffReceiveRequest) (shared.TranslateHandoffReceiveResponse, error)
	PrzyjmijZrodloMostu(ctx context.Context, z shared.TranslateBridgeSourceReceiveRequest) (shared.TranslateBridgeSourceReceiveResponse, error)
	OdesljWynikMostu(ctx context.Context, z shared.TranslateBridgeResultSendRequest) (shared.TranslateBridgeResultSendResponse, error)
	WydajWytwor(ctx context.Context, z shared.TranslateArtifactPublishRequest) (shared.TranslateArtifactPublishResponse, error)
	WykazKrokow(ctx context.Context, z shared.TranslateStepListRequest) (shared.TranslateStepListResponse, error)
}

// zarejestrujTlumaczenie wpina wszystkie komendy modułu Translate do
// rejestru rdzenia i wiąże adapter z emiterem zdarzeń.
func zarejestrujTlumaczenie(r *Rejestr, m Tlumaczenie, e *emiter) {
	if r == nil || m == nil {
		return
	}

	// Szyna zdarzeń trafia do adaptera tutaj, nie w montażu portów — emiter
	// powstaje w kompozycji.
	if adapter, ok := m.(*adapterTlumaczenia); ok {
		adapter.ZWyjsciem(e)
	}

	// przekład
	r.Zarejestruj(shared.CommandTranslateSourceSet, obsluz(m.UstawZrodlo))
	r.Zarejestruj(shared.CommandTranslateSourceSegment, obsluz(m.PodzielNaSegmenty))
	r.Zarejestruj(shared.CommandTranslateSourceDetect, obsluz(m.RozpoznajJezyk))
	r.Zarejestruj(shared.CommandTranslateTargetAdd, obsluz(m.DodajPanel))
	r.Zarejestruj(shared.CommandTranslateTranslationSet, obsluz(m.UstawTlumaczenie))
	r.Zarejestruj(shared.CommandTranslateBacktranslationRun, obsluz(m.TlumaczZwrotnie))

	// słownik
	r.Zarejestruj(shared.CommandTranslateGlossarySet, obsluz(m.UstawTerminy))
	r.Zarejestruj(shared.CommandTranslateGlossaryApply, obsluz(m.ZastosujSlownik))
	r.Zarejestruj(shared.CommandTranslateGlossaryOccurrences, obsluz(m.WystapieniaTerminow))
	r.Zarejestruj(shared.CommandTranslateGlossaryImport, obsluz(m.ImportujSlownik))
	r.Zarejestruj(shared.CommandTranslateGlossaryExport, obsluz(m.EksportujSlownik))
	r.Zarejestruj(shared.CommandTranslateMemorySuggest, obsluz(m.Podpowiedzi))

	// jakość i mowa
	r.Zarejestruj(shared.CommandTranslateQualityCheck, obsluz(m.SprawdzJakosc))
	r.Zarejestruj(shared.CommandTranslatePanelToneSet, obsluz(m.UstawTon))
	r.Zarejestruj(shared.CommandTranslateSpeechSynthesize, obsluz(m.SyntezujMowe))
	r.Zarejestruj(shared.CommandTranslatePanelExport, obsluz(m.EksportujPanel))

	// pamięć tłumaczeń jako byt Operatora
	r.Zarejestruj(shared.CommandTranslateMemoryList, obsluz(m.WykazPamieci))
	r.Zarejestruj(shared.CommandTranslateMemorySet, obsluz(m.UstawWpisPamieci))
	r.Zarejestruj(shared.CommandTranslateMemoryDelete, obsluz(m.UsunWpisPamieci))
	r.Zarejestruj(shared.CommandTranslateMemoryImport, obsluz(m.ImportujPamiec))
	r.Zarejestruj(shared.CommandTranslateMemoryExport, obsluz(m.EksportujPamiec))
	r.Zarejestruj(shared.CommandTranslateMemoryMaintain, obsluz(m.UtrzymajPamiec))
	r.Zarejestruj(shared.CommandTranslateMemoryPretranslate, obsluz(m.TlumaczWstepnie))
	r.Zarejestruj(shared.CommandTranslateMemoryAlign, obsluz(m.WyrownajTeksty))
	r.Zarejestruj(shared.CommandTranslateMemoryPolicyGet, obsluz(m.PolitykaPamieci))
	r.Zarejestruj(shared.CommandTranslateMemoryPolicySet, obsluz(m.UstawPolitykePamieci))

	// segmentacja
	r.Zarejestruj(shared.CommandTranslateSegmentationRulesList, obsluz(m.WykazRegulSegmentacji))
	r.Zarejestruj(shared.CommandTranslateSegmentationRulesSet, obsluz(m.UstawRegulySegmentacji))
	r.Zarejestruj(shared.CommandTranslateSegmentMerge, obsluz(m.ScalSegmenty))
	r.Zarejestruj(shared.CommandTranslateSegmentSplit, obsluz(m.PodzielSegment))

	// terminologia
	r.Zarejestruj(shared.CommandTranslateGlossaryList, obsluz(m.WykazTerminow))
	r.Zarejestruj(shared.CommandTranslateTermExtract, obsluz(m.WyjmijTerminy))

	// korekta i spójność
	r.Zarejestruj(shared.CommandTranslateProofreadRun, obsluz(m.SprawdzKorekte))
	r.Zarejestruj(shared.CommandTranslateProofreadApply, obsluz(m.ZastosujKorekte))
	r.Zarejestruj(shared.CommandTranslateConsistencyCheck, obsluz(m.SprawdzSpojnosc))

	// profile kontroli jakości i obieg zatwierdzeń
	r.Zarejestruj(shared.CommandTranslateQaProfileList, obsluz(m.WykazProfiliQa))
	r.Zarejestruj(shared.CommandTranslateQaProfileSet, obsluz(m.UstawProfilQa))
	r.Zarejestruj(shared.CommandTranslateQaProfileDelete, obsluz(m.UsunProfilQa))
	r.Zarejestruj(shared.CommandTranslateApprovalSet, obsluz(m.UstawZatwierdzenie))
	r.Zarejestruj(shared.CommandTranslateApprovalList, obsluz(m.WykazZatwierdzen))

	// dokument
	r.Zarejestruj(shared.CommandTranslateDocumentLoad, obsluz(m.WczytajDokument))
	r.Zarejestruj(shared.CommandTranslateDocumentRender, obsluz(m.ZlozDokument))
	r.Zarejestruj(shared.CommandTranslateDocumentLayoutCompare, obsluz(m.PorownajUklad))

	// lokalizacja oprogramowania
	r.Zarejestruj(shared.CommandTranslateResourceImport, obsluz(m.WczytajZasobLokalizacji))
	r.Zarejestruj(shared.CommandTranslateResourceExport, obsluz(m.WydajZasobLokalizacji))
	r.Zarejestruj(shared.CommandTranslateResourcePluralApply, obsluz(m.ZastosujFormyMnogie))
	r.Zarejestruj(shared.CommandTranslateResourceKeyContextSet, obsluz(m.UstawKontekstKlucza))
	r.Zarejestruj(shared.CommandTranslateXliffImport, obsluz(m.WczytajXliff))

	// napisy i dubbing
	r.Zarejestruj(shared.CommandTranslateSubtitleImport, obsluz(m.WczytajNapisy))
	r.Zarejestruj(shared.CommandTranslateSubtitleExport, obsluz(m.WydajNapisy))
	r.Zarejestruj(shared.CommandTranslateSubtitleTimingCheck, obsluz(m.SprawdzTaktowanieNapisow))
	r.Zarejestruj(shared.CommandTranslateDubbingScriptBuild, obsluz(m.ZlozScenariuszDubbingu))

	// silniki, pivot i przebieg pakietowy
	r.Zarejestruj(shared.CommandTranslateEngineProfileList, obsluz(m.WykazProfiliSilnikow))
	r.Zarejestruj(shared.CommandTranslateEngineProfileSet, obsluz(m.UstawProfilSilnika))
	r.Zarejestruj(shared.CommandTranslateEngineCompare, obsluz(m.PorownajSilniki))
	r.Zarejestruj(shared.CommandTranslatePivotPolicyGet, obsluz(m.PolitykaPivota))
	r.Zarejestruj(shared.CommandTranslatePivotPolicySet, obsluz(m.UstawPolitykePivota))
	r.Zarejestruj(shared.CommandTranslateBatchRun, obsluz(m.UruchomPakiet))

	// wymiana zewnętrzna
	r.Zarejestruj(shared.CommandTranslateHandoffBuild, obsluz(m.ZlozPakietPrzekazania))
	r.Zarejestruj(shared.CommandTranslateHandoffReceive, obsluz(m.PrzyjmijPakietPrzekazania))
	r.Zarejestruj(shared.CommandTranslateBridgeSourceReceive, obsluz(m.PrzyjmijZrodloMostu))
	r.Zarejestruj(shared.CommandTranslateBridgeResultSend, obsluz(m.OdesljWynikMostu))
	r.Zarejestruj(shared.CommandTranslateArtifactPublish, obsluz(m.WydajWytwor))
	r.Zarejestruj(shared.CommandTranslateStepList, obsluz(m.WykazKrokow))

	// Zdarzenie translate.translation.changed rozgłasza się po przekładzie
	// modelu i po korekcie operatora.
}
