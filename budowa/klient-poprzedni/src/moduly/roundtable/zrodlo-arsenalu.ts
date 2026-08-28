import {
  Command,
  type RoundtableAgreementGetRequest,
  type RoundtableAgreementGetResponse,
  type RoundtableAnalysisRunRequest,
  type RoundtableAnalysisRunResponse,
  type RoundtableArgumentExportRequest,
  type RoundtableArgumentExportResponse,
  type RoundtableArgumentListRequest,
  type RoundtableArgumentListResponse,
  type RoundtableArgumentPinRequest,
  type RoundtableArgumentPinResponse,
  type RoundtableCalibrationGetRequest,
  type RoundtableCalibrationGetResponse,
  type RoundtableClusterGetRequest,
  type RoundtableClusterGetResponse,
  type RoundtableConsensusHandoffRequest,
  type RoundtableConsensusHandoffResponse,
  type RoundtableConsensusMinoritySetRequest,
  type RoundtableConsensusMinoritySetResponse,
  type RoundtableConsensusSetRequest,
  type RoundtableConsensusSetResponse,
  type RoundtableConsensusVersionListRequest,
  type RoundtableConsensusVersionListResponse,
  type RoundtableConvergenceGetRequest,
  type RoundtableConvergenceGetResponse,
  type RoundtableCruxGetRequest,
  type RoundtableCruxGetResponse,
  type RoundtableDebateBranchRequest,
  type RoundtableDebateBranchResponse,
  type RoundtableDebateFollowupRequest,
  type RoundtableDebateFollowupResponse,
  type RoundtableDebateGetRequest,
  type RoundtableDebateGetResponse,
  type RoundtableDecisionMatrixGetRequest,
  type RoundtableDecisionMatrixGetResponse,
  type RoundtableDecisionMatrixSetRequest,
  type RoundtableDecisionMatrixSetResponse,
  type RoundtableDriftGetRequest,
  type RoundtableDriftGetResponse,
  type RoundtableEvidenceListRequest,
  type RoundtableEvidenceListResponse,
  type RoundtableFallacyCatalogGetRequest,
  type RoundtableFallacyCatalogGetResponse,
  type RoundtableFallacyCatalogSetRequest,
  type RoundtableFallacyCatalogSetResponse,
  type RoundtableJudgeRunRequest,
  type RoundtableJudgeRunResponse,
  type RoundtableLeaderboardGetRequest,
  type RoundtableLeaderboardGetResponse,
  type RoundtableModelListRequest,
  type RoundtableModelListResponse,
  type RoundtableModelRemoveRequest,
  type RoundtableModelRemoveResponse,
  type RoundtableModelUpdateRequest,
  type RoundtableModelUpdateResponse,
  type RoundtableModerationTemplateListRequest,
  type RoundtableModerationTemplateListResponse,
  type RoundtableModerationTemplateSaveRequest,
  type RoundtableModerationTemplateSaveResponse,
  type RoundtableRatingSetRequest,
  type RoundtableRatingSetResponse,
  type RoundtableRoleListRequest,
  type RoundtableRoleListResponse,
  type RoundtableRubricListRequest,
  type RoundtableRubricListResponse,
  type RoundtableRubricSetRequest,
  type RoundtableRubricSetResponse,
  type RoundtableSpeechSynthesizeRequest,
  type RoundtableSpeechSynthesizeResponse,
  type RoundtableStatementRegenerateRequest,
  type RoundtableStatementRegenerateResponse,
  type RoundtableTeamApplyRequest,
  type RoundtableTeamApplyResponse,
  type RoundtableTeamListRequest,
  type RoundtableTeamListResponse,
  type RoundtableTeamSaveRequest,
  type RoundtableTeamSaveResponse,
  type RoundtableTranscriptExportRequest,
  type RoundtableTranscriptExportResponse,
  type RoundtableVoteCastRequest,
  type RoundtableVoteCastResponse,
  type RoundtableVoteGetRequest,
  type RoundtableVoteGetResponse,
  type RoundtableVoteStartRequest,
  type RoundtableVoteStartResponse,
} from "../../../../shared/contract";
import type { Kanal, Wynik } from "../../protokol/kanal";
import { czyObiekt, czyTablica, czyTekst, sprawdzKsztalt } from "../../protokol/ksztalt-odpowiedzi";
import { wywolaj } from "../../protokol/wywolanie";

/** Arsenał modułu Roundtable widziany przez klienta: komendy obszaru czytające, analizujące i oceniające zapis debaty, każda oddająca Wynik ze sprawdzonym kształtem odpowiedzi. */
export interface ZrodloArsenaluRoundtable {
  /** `roundtable.debate.get` — pełny stan debaty okna po jego otwarciu. */
  stanDebaty(zadanie: RoundtableDebateGetRequest): Promise<Wynik<RoundtableDebateGetResponse>>;
  /** `roundtable.model.list` — odbudowa wykazu uczestników. */
  uczestnicy(zadanie: RoundtableModelListRequest): Promise<Wynik<RoundtableModelListResponse>>;
  /** `roundtable.model.remove` — usunięcie uczestnika ze składu. */
  usunModel(zadanie: RoundtableModelRemoveRequest): Promise<Wynik<RoundtableModelRemoveResponse>>;
  /** `roundtable.model.update` — tożsamość, rola, waga i oznaczenie uczestnika. */
  zmienModel(zadanie: RoundtableModelUpdateRequest): Promise<Wynik<RoundtableModelUpdateResponse>>;
  /** `roundtable.team.save` — zapis składu jako nazwanego zespołu. */
  zapiszZespol(zadanie: RoundtableTeamSaveRequest): Promise<Wynik<RoundtableTeamSaveResponse>>;
  /** `roundtable.team.list` — zapisane zespoły. */
  zespoly(zadanie: RoundtableTeamListRequest): Promise<Wynik<RoundtableTeamListResponse>>;
  /** `roundtable.team.apply` — wniesienie zespołu do okna debaty. */
  wniesZespol(zadanie: RoundtableTeamApplyRequest): Promise<Wynik<RoundtableTeamApplyResponse>>;
  /** `roundtable.role.list` — biblioteka ról wraz z promptami systemowymi. */
  role(zadanie: RoundtableRoleListRequest): Promise<Wynik<RoundtableRoleListResponse>>;

  /** `roundtable.debate.followup` — pytanie do jednego uczestnika poza turą. */
  doprecyzuj(
    zadanie: RoundtableDebateFollowupRequest,
  ): Promise<Wynik<RoundtableDebateFollowupResponse>>;
  /** `roundtable.statement.regenerate` — powtórzenie wywołania kanału. */
  powtorz(
    zadanie: RoundtableStatementRegenerateRequest,
  ): Promise<Wynik<RoundtableStatementRegenerateResponse>>;
  /** `roundtable.debate.branch` — wariant tury do porównania. */
  rozgalez(zadanie: RoundtableDebateBranchRequest): Promise<Wynik<RoundtableDebateBranchResponse>>;

  /** `roundtable.argument.list` — graf argumentów debaty albo tury. */
  graf(zadanie: RoundtableArgumentListRequest): Promise<Wynik<RoundtableArgumentListResponse>>;
  /** `roundtable.argument.pin` — oznaczenie argumentu kluczowego. */
  przypnij(zadanie: RoundtableArgumentPinRequest): Promise<Wynik<RoundtableArgumentPinResponse>>;
  /** `roundtable.argument.export` — wydanie grafu w formacie wymiany. */
  wydajGraf(
    zadanie: RoundtableArgumentExportRequest,
  ): Promise<Wynik<RoundtableArgumentExportResponse>>;
  /** `roundtable.analysis.run` — analiza zapisu debaty wskazanego rodzaju. */
  analizuj(zadanie: RoundtableAnalysisRunRequest): Promise<Wynik<RoundtableAnalysisRunResponse>>;
  /** `roundtable.evidence.list` — rejestr dowodów i twierdzeń niepopartych. */
  dowody(zadanie: RoundtableEvidenceListRequest): Promise<Wynik<RoundtableEvidenceListResponse>>;
  /** `roundtable.fallacy.catalog.get` — katalog błędów logicznych. */
  katalogBledow(
    zadanie: RoundtableFallacyCatalogGetRequest,
  ): Promise<Wynik<RoundtableFallacyCatalogGetResponse>>;
  /** `roundtable.fallacy.catalog.set` — zakres wykrywania błędów w oknie. */
  ustawKatalogBledow(
    zadanie: RoundtableFallacyCatalogSetRequest,
  ): Promise<Wynik<RoundtableFallacyCatalogSetResponse>>;

  /** `roundtable.agreement.get` — punkty zgody, sporu i macierz zgodności. */
  zgodnosc(zadanie: RoundtableAgreementGetRequest): Promise<Wynik<RoundtableAgreementGetResponse>>;
  /** `roundtable.cluster.get` — obozy zbliżonych stanowisk. */
  klastry(zadanie: RoundtableClusterGetRequest): Promise<Wynik<RoundtableClusterGetResponse>>;
  /** `roundtable.crux.get` — kluczowy punkt sporny. */
  punktSporny(zadanie: RoundtableCruxGetRequest): Promise<Wynik<RoundtableCruxGetResponse>>;
  /** `roundtable.convergence.get` — zbieżność stanowisk tura po turze. */
  zbieznosc(
    zadanie: RoundtableConvergenceGetRequest,
  ): Promise<Wynik<RoundtableConvergenceGetResponse>>;
  /** `roundtable.drift.get` — zmiany stanowisk między turami. */
  dryf(zadanie: RoundtableDriftGetRequest): Promise<Wynik<RoundtableDriftGetResponse>>;
  /** `roundtable.calibration.get` — pewność deklarowana wobec trafności. */
  kalibracja(
    zadanie: RoundtableCalibrationGetRequest,
  ): Promise<Wynik<RoundtableCalibrationGetResponse>>;

  /** `roundtable.vote.start` — otwarcie głosowania nad stanowiskami. */
  otworzGlosowanie(
    zadanie: RoundtableVoteStartRequest,
  ): Promise<Wynik<RoundtableVoteStartResponse>>;
  /** `roundtable.vote.cast` — oddanie głosu. */
  oddajGlos(zadanie: RoundtableVoteCastRequest): Promise<Wynik<RoundtableVoteCastResponse>>;
  /** `roundtable.vote.get` — głosowanie wraz z wynikiem agregacji. */
  glosowanie(zadanie: RoundtableVoteGetRequest): Promise<Wynik<RoundtableVoteGetResponse>>;
  /** `roundtable.rating.set` — ocena Operatora. */
  ocen(zadanie: RoundtableRatingSetRequest): Promise<Wynik<RoundtableRatingSetResponse>>;
  /** `roundtable.rubric.set` — rubryka oceny wraz z kryteriami i wagami. */
  zapiszRubryke(zadanie: RoundtableRubricSetRequest): Promise<Wynik<RoundtableRubricSetResponse>>;
  /** `roundtable.rubric.list` — rubryki dostępne dla okna. */
  rubryki(zadanie: RoundtableRubricListRequest): Promise<Wynik<RoundtableRubricListResponse>>;
  /** `roundtable.judge.run` — ocena wypowiedzi modelami-sędziami. */
  osadz(zadanie: RoundtableJudgeRunRequest): Promise<Wynik<RoundtableJudgeRunResponse>>;
  /** `roundtable.leaderboard.get` — ranking akumulowany między sesjami. */
  ranking(
    zadanie: RoundtableLeaderboardGetRequest,
  ): Promise<Wynik<RoundtableLeaderboardGetResponse>>;
  /** `roundtable.decision.matrix.set` — macierz decyzyjna wariantów. */
  zapiszMacierz(
    zadanie: RoundtableDecisionMatrixSetRequest,
  ): Promise<Wynik<RoundtableDecisionMatrixSetResponse>>;
  /** `roundtable.decision.matrix.get` — macierz wraz z wynikiem ważonym. */
  macierz(
    zadanie: RoundtableDecisionMatrixGetRequest,
  ): Promise<Wynik<RoundtableDecisionMatrixGetResponse>>;

  /** `roundtable.consensus.set` — redakcja stanowiska końcowego. */
  zapiszStanowisko(
    zadanie: RoundtableConsensusSetRequest,
  ): Promise<Wynik<RoundtableConsensusSetResponse>>;
  /** `roundtable.consensus.minority.set` — podpisane zdanie odrębne. */
  zapiszZdanieOdrebne(
    zadanie: RoundtableConsensusMinoritySetRequest,
  ): Promise<Wynik<RoundtableConsensusMinoritySetResponse>>;
  /** `roundtable.consensus.version.list` — kolejne redakcje stanowiska. */
  wersjeStanowiska(
    zadanie: RoundtableConsensusVersionListRequest,
  ): Promise<Wynik<RoundtableConsensusVersionListResponse>>;
  /** `roundtable.consensus.handoff` — przekazanie stanowiska do modułu. */
  przekazStanowisko(
    zadanie: RoundtableConsensusHandoffRequest,
  ): Promise<Wynik<RoundtableConsensusHandoffResponse>>;

  /** `roundtable.moderation.template.save` — szablon moderacji. */
  zapiszSzablon(
    zadanie: RoundtableModerationTemplateSaveRequest,
  ): Promise<Wynik<RoundtableModerationTemplateSaveResponse>>;
  /** `roundtable.moderation.template.list` — zapisane szablony moderacji. */
  szablony(
    zadanie: RoundtableModerationTemplateListRequest,
  ): Promise<Wynik<RoundtableModerationTemplateListResponse>>;
  /** `roundtable.transcript.export` — wydanie zapisu debaty jako artefaktu. */
  wydajTranskrypt(
    zadanie: RoundtableTranscriptExportRequest,
  ): Promise<Wynik<RoundtableTranscriptExportResponse>>;
  /** `roundtable.speech.synthesize` — odsłuch przebiegu debaty. */
  odsluch(
    zadanie: RoundtableSpeechSynthesizeRequest,
  ): Promise<Wynik<RoundtableSpeechSynthesizeResponse>>;
}

/**
 * Sprawdzian artefaktu — wspólny dla trzech wydań, bo wszystkie trzy oddają
 * ten sam kształt: identyfikator artefaktu i odwołanie do jego bajtów.
 * Odpowiedź bez identyfikatora jest odpowiedzią o wydaniu, którego nie ma.
 */
function czyArtefakt(tresc: { artifactId?: unknown }): boolean {
  return czyTekst(tresc.artifactId) && tresc.artifactId !== "";
}

export function utworzZrodloArsenaluRoundtable(kanal: Kanal): ZrodloArsenaluRoundtable {
  return {
    stanDebaty: async (zadanie) =>
      sprawdzKsztalt(
        await wywolaj(kanal, Command.RoundtableDebateGet, zadanie),
        Command.RoundtableDebateGet,
        (tresc) => czyObiekt(tresc.snapshot) && czyTablica(tresc.snapshot.participants),
      ),
    uczestnicy: async (zadanie) =>
      sprawdzKsztalt(
        await wywolaj(kanal, Command.RoundtableModelList, zadanie),
        Command.RoundtableModelList,
        (tresc) => czyTablica(tresc.participants),
      ),
    usunModel: async (zadanie) =>
      sprawdzKsztalt(
        await wywolaj(kanal, Command.RoundtableModelRemove, zadanie),
        Command.RoundtableModelRemove,
        (tresc) => czyTablica(tresc.participants),
      ),
    zmienModel: async (zadanie) =>
      sprawdzKsztalt(
        await wywolaj(kanal, Command.RoundtableModelUpdate, zadanie),
        Command.RoundtableModelUpdate,
        (tresc) => czyObiekt(tresc.participant),
      ),
    zapiszZespol: async (zadanie) =>
      sprawdzKsztalt(
        await wywolaj(kanal, Command.RoundtableTeamSave, zadanie),
        Command.RoundtableTeamSave,
        (tresc) => czyObiekt(tresc.team),
      ),
    zespoly: async (zadanie) =>
      sprawdzKsztalt(
        await wywolaj(kanal, Command.RoundtableTeamList, zadanie),
        Command.RoundtableTeamList,
        (tresc) => czyTablica(tresc.teams),
      ),
    wniesZespol: async (zadanie) =>
      sprawdzKsztalt(
        await wywolaj(kanal, Command.RoundtableTeamApply, zadanie),
        Command.RoundtableTeamApply,
        (tresc) => czyTablica(tresc.participants),
      ),
    role: async (zadanie) =>
      sprawdzKsztalt(
        await wywolaj(kanal, Command.RoundtableRoleList, zadanie),
        Command.RoundtableRoleList,
        (tresc) => czyTablica(tresc.roles),
      ),

    doprecyzuj: async (zadanie) =>
      sprawdzKsztalt(
        await wywolaj(kanal, Command.RoundtableDebateFollowup, zadanie),
        Command.RoundtableDebateFollowup,
        (tresc) => czyObiekt(tresc.statement),
      ),
    powtorz: async (zadanie) =>
      sprawdzKsztalt(
        await wywolaj(kanal, Command.RoundtableStatementRegenerate, zadanie),
        Command.RoundtableStatementRegenerate,
        (tresc) => czyObiekt(tresc.statement),
      ),
    rozgalez: async (zadanie) =>
      sprawdzKsztalt(
        await wywolaj(kanal, Command.RoundtableDebateBranch, zadanie),
        Command.RoundtableDebateBranch,
        (tresc) => czyObiekt(tresc.turn),
      ),

    graf: async (zadanie) =>
      sprawdzKsztalt(
        await wywolaj(kanal, Command.RoundtableArgumentList, zadanie),
        Command.RoundtableArgumentList,
        (tresc) => czyObiekt(tresc.graph) && czyTablica(tresc.graph.nodes),
      ),
    przypnij: async (zadanie) =>
      sprawdzKsztalt(
        await wywolaj(kanal, Command.RoundtableArgumentPin, zadanie),
        Command.RoundtableArgumentPin,
        (tresc) => czyObiekt(tresc.node),
      ),
    wydajGraf: async (zadanie) =>
      sprawdzKsztalt(
        await wywolaj(kanal, Command.RoundtableArgumentExport, zadanie),
        Command.RoundtableArgumentExport,
        czyArtefakt,
      ),
    analizuj: async (zadanie) =>
      sprawdzKsztalt(
        await wywolaj(kanal, Command.RoundtableAnalysisRun, zadanie),
        Command.RoundtableAnalysisRun,
        (tresc) => czyTablica(tresc.findings),
      ),
    dowody: async (zadanie) =>
      sprawdzKsztalt(
        await wywolaj(kanal, Command.RoundtableEvidenceList, zadanie),
        Command.RoundtableEvidenceList,
        (tresc) => czyTablica(tresc.evidence),
      ),
    katalogBledow: async (zadanie) =>
      sprawdzKsztalt(
        await wywolaj(kanal, Command.RoundtableFallacyCatalogGet, zadanie),
        Command.RoundtableFallacyCatalogGet,
        (tresc) => czyTablica(tresc.definitions),
      ),
    ustawKatalogBledow: async (zadanie) =>
      sprawdzKsztalt(
        await wywolaj(kanal, Command.RoundtableFallacyCatalogSet, zadanie),
        Command.RoundtableFallacyCatalogSet,
        (tresc) => czyTablica(tresc.definitions),
      ),

    zgodnosc: async (zadanie) =>
      sprawdzKsztalt(
        await wywolaj(kanal, Command.RoundtableAgreementGet, zadanie),
        Command.RoundtableAgreementGet,
        (tresc) => czyTablica(tresc.points),
      ),
    klastry: async (zadanie) =>
      sprawdzKsztalt(
        await wywolaj(kanal, Command.RoundtableClusterGet, zadanie),
        Command.RoundtableClusterGet,
        (tresc) => czyTablica(tresc.clusters),
      ),
    // Punkt sporny jest nieobowiązkowy — sprawdzian pilnuje wyłącznie tego, że odpowiedź jest obiektem.
    punktSporny: async (zadanie) =>
      sprawdzKsztalt(
        await wywolaj(kanal, Command.RoundtableCruxGet, zadanie),
        Command.RoundtableCruxGet,
        (tresc) => czyObiekt(tresc),
      ),
    zbieznosc: async (zadanie) =>
      sprawdzKsztalt(
        await wywolaj(kanal, Command.RoundtableConvergenceGet, zadanie),
        Command.RoundtableConvergenceGet,
        (tresc) => czyTablica(tresc.points),
      ),
    dryf: async (zadanie) =>
      sprawdzKsztalt(
        await wywolaj(kanal, Command.RoundtableDriftGet, zadanie),
        Command.RoundtableDriftGet,
        (tresc) => czyTablica(tresc.entries),
      ),
    kalibracja: async (zadanie) =>
      sprawdzKsztalt(
        await wywolaj(kanal, Command.RoundtableCalibrationGet, zadanie),
        Command.RoundtableCalibrationGet,
        (tresc) => czyTablica(tresc.calibrations),
      ),

    otworzGlosowanie: async (zadanie) =>
      sprawdzKsztalt(
        await wywolaj(kanal, Command.RoundtableVoteStart, zadanie),
        Command.RoundtableVoteStart,
        (tresc) => czyObiekt(tresc.vote) && czyTablica(tresc.vote.options),
      ),
    oddajGlos: async (zadanie) =>
      sprawdzKsztalt(
        await wywolaj(kanal, Command.RoundtableVoteCast, zadanie),
        Command.RoundtableVoteCast,
        (tresc) => czyObiekt(tresc.ballot),
      ),
    glosowanie: async (zadanie) =>
      sprawdzKsztalt(
        await wywolaj(kanal, Command.RoundtableVoteGet, zadanie),
        Command.RoundtableVoteGet,
        (tresc) => czyObiekt(tresc.vote),
      ),
    ocen: async (zadanie) =>
      sprawdzKsztalt(
        await wywolaj(kanal, Command.RoundtableRatingSet, zadanie),
        Command.RoundtableRatingSet,
        (tresc) => czyObiekt(tresc.rating),
      ),
    zapiszRubryke: async (zadanie) =>
      sprawdzKsztalt(
        await wywolaj(kanal, Command.RoundtableRubricSet, zadanie),
        Command.RoundtableRubricSet,
        (tresc) => czyObiekt(tresc.rubric) && czyTablica(tresc.rubric.criteria),
      ),
    rubryki: async (zadanie) =>
      sprawdzKsztalt(
        await wywolaj(kanal, Command.RoundtableRubricList, zadanie),
        Command.RoundtableRubricList,
        (tresc) => czyTablica(tresc.rubrics),
      ),
    osadz: async (zadanie) =>
      sprawdzKsztalt(
        await wywolaj(kanal, Command.RoundtableJudgeRun, zadanie),
        Command.RoundtableJudgeRun,
        (tresc) => czyTablica(tresc.judgements),
      ),
    ranking: async (zadanie) =>
      sprawdzKsztalt(
        await wywolaj(kanal, Command.RoundtableLeaderboardGet, zadanie),
        Command.RoundtableLeaderboardGet,
        (tresc) => czyTablica(tresc.entries),
      ),
    zapiszMacierz: async (zadanie) =>
      sprawdzKsztalt(
        await wywolaj(kanal, Command.RoundtableDecisionMatrixSet, zadanie),
        Command.RoundtableDecisionMatrixSet,
        (tresc) => czyObiekt(tresc.matrix) && czyTablica(tresc.matrix.options),
      ),
    macierz: async (zadanie) =>
      sprawdzKsztalt(
        await wywolaj(kanal, Command.RoundtableDecisionMatrixGet, zadanie),
        Command.RoundtableDecisionMatrixGet,
        (tresc) => czyObiekt(tresc.matrix),
      ),

    zapiszStanowisko: async (zadanie) =>
      sprawdzKsztalt(
        await wywolaj(kanal, Command.RoundtableConsensusSet, zadanie),
        Command.RoundtableConsensusSet,
        (tresc) => czyObiekt(tresc.consensus),
      ),
    zapiszZdanieOdrebne: async (zadanie) =>
      sprawdzKsztalt(
        await wywolaj(kanal, Command.RoundtableConsensusMinoritySet, zadanie),
        Command.RoundtableConsensusMinoritySet,
        (tresc) => czyObiekt(tresc.minority),
      ),
    wersjeStanowiska: async (zadanie) =>
      sprawdzKsztalt(
        await wywolaj(kanal, Command.RoundtableConsensusVersionList, zadanie),
        Command.RoundtableConsensusVersionList,
        (tresc) => czyTablica(tresc.versions),
      ),
    przekazStanowisko: async (zadanie) =>
      sprawdzKsztalt(
        await wywolaj(kanal, Command.RoundtableConsensusHandoff, zadanie),
        Command.RoundtableConsensusHandoff,
        (tresc) => czyObiekt(tresc.handoff),
      ),

    zapiszSzablon: async (zadanie) =>
      sprawdzKsztalt(
        await wywolaj(kanal, Command.RoundtableModerationTemplateSave, zadanie),
        Command.RoundtableModerationTemplateSave,
        (tresc) => czyObiekt(tresc.template),
      ),
    szablony: async (zadanie) =>
      sprawdzKsztalt(
        await wywolaj(kanal, Command.RoundtableModerationTemplateList, zadanie),
        Command.RoundtableModerationTemplateList,
        (tresc) => czyTablica(tresc.templates),
      ),
    wydajTranskrypt: async (zadanie) =>
      sprawdzKsztalt(
        await wywolaj(kanal, Command.RoundtableTranscriptExport, zadanie),
        Command.RoundtableTranscriptExport,
        czyArtefakt,
      ),
    odsluch: async (zadanie) =>
      sprawdzKsztalt(
        await wywolaj(kanal, Command.RoundtableSpeechSynthesize, zadanie),
        Command.RoundtableSpeechSynthesize,
        czyArtefakt,
      ),
  };
}
