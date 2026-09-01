// Plik wpina cztery komendy obszaru roundtable.* modułu Roundtable wraz z jego czterema oknami operacyjnymi: Model Panels, Debate Panel, Moderator Panel i Consensus Panel.
package core

import (
	"context"

	"danacoconsole/shared"
)

// Debata jest portem modułu Roundtable. Obejmuje wszystkie cztery okna
// operacyjne: Model Panels (skład), Debate Panel (tury i wypowiedzi), Argument
// Map & Analysis, Voting & Evaluation Center, Moderator Panel i Consensus Panel.
type Debata interface {
	DodajModel(ctx context.Context, z shared.RoundtableModelAddRequest) (shared.RoundtableModelAddResponse, error)
	Uruchom(ctx context.Context, z shared.RoundtableDebateStartRequest) (shared.RoundtableDebateStartResponse, error)
	Moderuj(ctx context.Context, z shared.RoundtableModeratorDirectRequest) (shared.RoundtableModeratorDirectResponse, error)
	Stanowisko(ctx context.Context, z shared.RoundtableConsensusGetRequest) (shared.RoundtableConsensusGetResponse, error)

	// Skład debaty — Model Panels.
	StanDebaty(ctx context.Context, z shared.RoundtableDebateGetRequest) (shared.RoundtableDebateGetResponse, error)
	Uczestnicy(ctx context.Context, z shared.RoundtableModelListRequest) (shared.RoundtableModelListResponse, error)
	UsunModel(ctx context.Context, z shared.RoundtableModelRemoveRequest) (shared.RoundtableModelRemoveResponse, error)
	ZmienModel(ctx context.Context, z shared.RoundtableModelUpdateRequest) (shared.RoundtableModelUpdateResponse, error)
	ZapiszZespol(ctx context.Context, z shared.RoundtableTeamSaveRequest) (shared.RoundtableTeamSaveResponse, error)
	Zespoly(ctx context.Context, z shared.RoundtableTeamListRequest) (shared.RoundtableTeamListResponse, error)
	WniesZespol(ctx context.Context, z shared.RoundtableTeamApplyRequest) (shared.RoundtableTeamApplyResponse, error)
	Role(ctx context.Context, z shared.RoundtableRoleListRequest) (shared.RoundtableRoleListResponse, error)

	// Wątki poboczne i warianty — Debate Panel.
	Doprecyzuj(ctx context.Context, z shared.RoundtableDebateFollowupRequest) (shared.RoundtableDebateFollowupResponse, error)
	Powtorz(ctx context.Context, z shared.RoundtableStatementRegenerateRequest) (shared.RoundtableStatementRegenerateResponse, error)
	Rozgalez(ctx context.Context, z shared.RoundtableDebateBranchRequest) (shared.RoundtableDebateBranchResponse, error)

	// Graf argumentów i analiza — Argument Map & Analysis.
	Graf(ctx context.Context, z shared.RoundtableArgumentListRequest) (shared.RoundtableArgumentListResponse, error)
	Przypnij(ctx context.Context, z shared.RoundtableArgumentPinRequest) (shared.RoundtableArgumentPinResponse, error)
	WydajGraf(ctx context.Context, z shared.RoundtableArgumentExportRequest) (shared.RoundtableArgumentExportResponse, error)
	Analizuj(ctx context.Context, z shared.RoundtableAnalysisRunRequest) (shared.RoundtableAnalysisRunResponse, error)
	Dowody(ctx context.Context, z shared.RoundtableEvidenceListRequest) (shared.RoundtableEvidenceListResponse, error)
	KatalogBledow(ctx context.Context, z shared.RoundtableFallacyCatalogGetRequest) (shared.RoundtableFallacyCatalogGetResponse, error)
	UstawKatalogBledow(ctx context.Context, z shared.RoundtableFallacyCatalogSetRequest) (shared.RoundtableFallacyCatalogSetResponse, error)

	// Zgoda, spór i rozbieżność.
	Zgodnosc(ctx context.Context, z shared.RoundtableAgreementGetRequest) (shared.RoundtableAgreementGetResponse, error)
	Klastry(ctx context.Context, z shared.RoundtableClusterGetRequest) (shared.RoundtableClusterGetResponse, error)
	PunktSporny(ctx context.Context, z shared.RoundtableCruxGetRequest) (shared.RoundtableCruxGetResponse, error)
	Zbieznosc(ctx context.Context, z shared.RoundtableConvergenceGetRequest) (shared.RoundtableConvergenceGetResponse, error)
	Dryf(ctx context.Context, z shared.RoundtableDriftGetRequest) (shared.RoundtableDriftGetResponse, error)
	Kalibracja(ctx context.Context, z shared.RoundtableCalibrationGetRequest) (shared.RoundtableCalibrationGetResponse, error)

	// Głosowanie i ocena — Voting & Evaluation Center.
	OtworzGlosowanie(ctx context.Context, z shared.RoundtableVoteStartRequest) (shared.RoundtableVoteStartResponse, error)
	OddajGlos(ctx context.Context, z shared.RoundtableVoteCastRequest) (shared.RoundtableVoteCastResponse, error)
	Glosowanie(ctx context.Context, z shared.RoundtableVoteGetRequest) (shared.RoundtableVoteGetResponse, error)
	Ocen(ctx context.Context, z shared.RoundtableRatingSetRequest) (shared.RoundtableRatingSetResponse, error)
	ZapiszRubryke(ctx context.Context, z shared.RoundtableRubricSetRequest) (shared.RoundtableRubricSetResponse, error)
	Rubryki(ctx context.Context, z shared.RoundtableRubricListRequest) (shared.RoundtableRubricListResponse, error)
	Osadz(ctx context.Context, z shared.RoundtableJudgeRunRequest) (shared.RoundtableJudgeRunResponse, error)
	Ranking(ctx context.Context, z shared.RoundtableLeaderboardGetRequest) (shared.RoundtableLeaderboardGetResponse, error)
	ZapiszMacierz(ctx context.Context, z shared.RoundtableDecisionMatrixSetRequest) (shared.RoundtableDecisionMatrixSetResponse, error)
	Macierz(ctx context.Context, z shared.RoundtableDecisionMatrixGetRequest) (shared.RoundtableDecisionMatrixGetResponse, error)

	// Stanowisko końcowe — Consensus Panel.
	ZapiszStanowisko(ctx context.Context, z shared.RoundtableConsensusSetRequest) (shared.RoundtableConsensusSetResponse, error)
	ZapiszZdanieOdrebne(ctx context.Context, z shared.RoundtableConsensusMinoritySetRequest) (shared.RoundtableConsensusMinoritySetResponse, error)
	WersjeStanowiska(ctx context.Context, z shared.RoundtableConsensusVersionListRequest) (shared.RoundtableConsensusVersionListResponse, error)
	PrzekazStanowisko(ctx context.Context, z shared.RoundtableConsensusHandoffRequest) (shared.RoundtableConsensusHandoffResponse, error)

	// Moderacja i wydanie — Moderator Panel i panel akcji Debate Panel.
	ZapiszSzablon(ctx context.Context, z shared.RoundtableModerationTemplateSaveRequest) (shared.RoundtableModerationTemplateSaveResponse, error)
	Szablony(ctx context.Context, z shared.RoundtableModerationTemplateListRequest) (shared.RoundtableModerationTemplateListResponse, error)
	WydajTranskrypt(ctx context.Context, z shared.RoundtableTranscriptExportRequest) (shared.RoundtableTranscriptExportResponse, error)
	Odsluch(ctx context.Context, z shared.RoundtableSpeechSynthesizeRequest) (shared.RoundtableSpeechSynthesizeResponse, error)

	// PodepnijRozgloszenie oddaje drogę do zdarzenia zmiany debaty; wypowiedź powstaje poza komendą.
	PodepnijRozgloszenie(rozglos func(context.Context, shared.ChangeKind, shared.RoundtableTurn, *shared.RoundtableStatement))
}

// zarejestrujDebate wpina komplet komend modułu Roundtable w rejestr komend rdzenia, przy starcie modułu.
func zarejestrujDebate(r *Rejestr, d Debata, e *emiter) {
	if r == nil || d == nil {
		return
	}
	d.PodepnijRozgloszenie(e.debata)

	// roundtable.model.add i roundtable.consensus.get niczego nie rozgłaszają: skład i odczyt tur.
	r.Zarejestruj(shared.CommandRoundtableModelAdd, obsluz(d.DodajModel))
	r.Zarejestruj(shared.CommandRoundtableDebateStart, obsluz(d.Uruchom))
	r.Zarejestruj(shared.CommandRoundtableModeratorDirect, obsluz(d.Moderuj))
	r.Zarejestruj(shared.CommandRoundtableConsensusGet, obsluz(d.Stanowisko))

	r.Zarejestruj(shared.CommandRoundtableDebateGet, obsluz(d.StanDebaty))
	r.Zarejestruj(shared.CommandRoundtableModelList, obsluz(d.Uczestnicy))
	r.Zarejestruj(shared.CommandRoundtableModelRemove, obsluz(d.UsunModel))
	r.Zarejestruj(shared.CommandRoundtableModelUpdate, obsluz(d.ZmienModel))
	r.Zarejestruj(shared.CommandRoundtableTeamSave, obsluz(d.ZapiszZespol))
	r.Zarejestruj(shared.CommandRoundtableTeamList, obsluz(d.Zespoly))
	r.Zarejestruj(shared.CommandRoundtableTeamApply, obsluz(d.WniesZespol))
	r.Zarejestruj(shared.CommandRoundtableRoleList, obsluz(d.Role))

	r.Zarejestruj(shared.CommandRoundtableDebateFollowup, obsluz(d.Doprecyzuj))
	r.Zarejestruj(shared.CommandRoundtableStatementRegenerate, obsluz(d.Powtorz))
	r.Zarejestruj(shared.CommandRoundtableDebateBranch, obsluz(d.Rozgalez))

	r.Zarejestruj(shared.CommandRoundtableArgumentList, obsluz(d.Graf))
	r.Zarejestruj(shared.CommandRoundtableArgumentPin, obsluz(d.Przypnij))
	r.Zarejestruj(shared.CommandRoundtableArgumentExport, obsluz(d.WydajGraf))
	r.Zarejestruj(shared.CommandRoundtableAnalysisRun, obsluz(d.Analizuj))
	r.Zarejestruj(shared.CommandRoundtableEvidenceList, obsluz(d.Dowody))
	r.Zarejestruj(shared.CommandRoundtableFallacyCatalogGet, obsluz(d.KatalogBledow))
	r.Zarejestruj(shared.CommandRoundtableFallacyCatalogSet, obsluz(d.UstawKatalogBledow))

	r.Zarejestruj(shared.CommandRoundtableAgreementGet, obsluz(d.Zgodnosc))
	r.Zarejestruj(shared.CommandRoundtableClusterGet, obsluz(d.Klastry))
	r.Zarejestruj(shared.CommandRoundtableCruxGet, obsluz(d.PunktSporny))
	r.Zarejestruj(shared.CommandRoundtableConvergenceGet, obsluz(d.Zbieznosc))
	r.Zarejestruj(shared.CommandRoundtableDriftGet, obsluz(d.Dryf))
	r.Zarejestruj(shared.CommandRoundtableCalibrationGet, obsluz(d.Kalibracja))

	r.Zarejestruj(shared.CommandRoundtableVoteStart, obsluz(d.OtworzGlosowanie))
	r.Zarejestruj(shared.CommandRoundtableVoteCast, obsluz(d.OddajGlos))
	r.Zarejestruj(shared.CommandRoundtableVoteGet, obsluz(d.Glosowanie))
	r.Zarejestruj(shared.CommandRoundtableRatingSet, obsluz(d.Ocen))
	r.Zarejestruj(shared.CommandRoundtableRubricSet, obsluz(d.ZapiszRubryke))
	r.Zarejestruj(shared.CommandRoundtableRubricList, obsluz(d.Rubryki))
	r.Zarejestruj(shared.CommandRoundtableJudgeRun, obsluz(d.Osadz))
	r.Zarejestruj(shared.CommandRoundtableLeaderboardGet, obsluz(d.Ranking))
	r.Zarejestruj(shared.CommandRoundtableDecisionMatrixSet, obsluz(d.ZapiszMacierz))
	r.Zarejestruj(shared.CommandRoundtableDecisionMatrixGet, obsluz(d.Macierz))

	r.Zarejestruj(shared.CommandRoundtableConsensusSet, obsluz(d.ZapiszStanowisko))
	r.Zarejestruj(shared.CommandRoundtableConsensusMinoritySet, obsluz(d.ZapiszZdanieOdrebne))
	r.Zarejestruj(shared.CommandRoundtableConsensusVersionList, obsluz(d.WersjeStanowiska))
	r.Zarejestruj(shared.CommandRoundtableConsensusHandoff, obsluz(d.PrzekazStanowisko))

	r.Zarejestruj(shared.CommandRoundtableModerationTemplateSave, obsluz(d.ZapiszSzablon))
	r.Zarejestruj(shared.CommandRoundtableModerationTemplateList, obsluz(d.Szablony))
	r.Zarejestruj(shared.CommandRoundtableTranscriptExport, obsluz(d.WydajTranskrypt))
	r.Zarejestruj(shared.CommandRoundtableSpeechSynthesize, obsluz(d.Odsluch))
}

// debata rozgłasza przyrost debaty. Sesja komunikatu jest pusta, bo debata jest
// bytem okna, nie karty sesji: kontrakt kieruje wszystkie cztery komendy przez
// `windowId` i nie zna sesji nadrzędnej debaty.
func (e *emiter) debata(ctx context.Context, zmiana shared.ChangeKind, t shared.RoundtableTurn,
	w *shared.RoundtableStatement) {

	zdarzenie := shared.RoundtableDebateChangedEvent{Change: zmiana, Turn: t}
	if w != nil {
		zdarzenie.Statement = w
	}
	e.wyslijDoKonta(ctx, shared.EventRoundtableDebateChanged, "", zdarzenie)
}
