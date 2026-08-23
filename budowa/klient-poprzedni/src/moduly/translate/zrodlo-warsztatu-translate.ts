import {
  Command,
  type RequestOf,
  type TranslateApprovalListResponse,
  type TranslateApprovalSetResponse,
  type TranslateArtifactPublishResponse,
  type TranslateBatchRunResponse,
  type TranslateBridgeResultSendResponse,
  type TranslateBridgeSourceReceiveResponse,
  type TranslateConsistencyCheckResponse,
  type TranslateDocumentLayoutCompareResponse,
  type TranslateDocumentLoadResponse,
  type TranslateDocumentRenderResponse,
  type TranslateDubbingScriptBuildResponse,
  type TranslateEngineCompareResponse,
  type TranslateEngineProfileListResponse,
  type TranslateEngineProfileSetResponse,
  type TranslateGlossaryListResponse,
  type TranslateHandoffBuildResponse,
  type TranslateHandoffReceiveResponse,
  type TranslateMemoryAlignResponse,
  type TranslateMemoryDeleteResponse,
  type TranslateMemoryExportResponse,
  type TranslateMemoryImportResponse,
  type TranslateMemoryListResponse,
  type TranslateMemoryMaintainResponse,
  type TranslateMemoryPolicyGetResponse,
  type TranslateMemoryPolicySetResponse,
  type TranslateMemoryPretranslateResponse,
  type TranslateMemorySetResponse,
  type TranslatePivotPolicyGetResponse,
  type TranslatePivotPolicySetResponse,
  type TranslateProofreadApplyResponse,
  type TranslateProofreadRunResponse,
  type TranslateQaProfileDeleteResponse,
  type TranslateQaProfileListResponse,
  type TranslateQaProfileSetResponse,
  type TranslateResourceExportResponse,
  type TranslateResourceImportResponse,
  type TranslateResourceKeyContextSetResponse,
  type TranslateResourcePluralApplyResponse,
  type TranslateSegmentMergeResponse,
  type TranslateSegmentSplitResponse,
  type TranslateSegmentationRulesListResponse,
  type TranslateSegmentationRulesSetResponse,
  type TranslateStepListResponse,
  type TranslateSubtitleExportResponse,
  type TranslateSubtitleImportResponse,
  type TranslateSubtitleTimingCheckResponse,
  type TranslateTermExtractResponse,
  type TranslateXliffImportResponse,
} from '../../../../shared/contract';
import type { Kanal } from '../../protokol/kanal';
import { czyObiekt, czyTablica } from '../../protokol/ksztalt-odpowiedzi';
import { zadaj, type WynikTranslate } from './wywolanie-translate';

/**
 * Warsztat modułu Translate — jedno źródło dla wszystkich rodzin komend, które
 * nie mieszczą się w czterech oknach pierwotnych: pamięć jako byt Operatora,
 * segmentacja, terminologia, korekta i spójność, profile kontroli jakości,
 * obieg zatwierdzeń, dokumenty, lokalizacja oprogramowania, napisy i dubbing,
 * silniki, polityka pivota, przebieg pakietowy oraz wymiana zewnętrzna.
 *
 * Jedno źródło, nie jedno na okno. Wszystkie te komendy idą tą samą drogą
 * (`zadaj` z `wywolanie-translate.ts`), różnią się wyłącznie nazwą i ładunkiem;
 * rozbicie ich na dwanaście plików dałoby dwanaście kopii tej samej obudowy,
 * a każda z nich musiałaby osobno pamiętać o odmowie `translate.unknown`.
 *
 * Sprawdzian kształtu jest przy każdej komendzie i nie jest formalnością: to on
 * odróżnia „rdzeń odpowiedział wynikiem” od „rdzeń odpowiedział pustą kopertą”.
 * Okno, które przyjmie pustą kopertę jako wynik, pokaże Operatorowi pustkę jako
 * skutek — czyli dokładnie tę szkodę, przed którą stoi cały ten moduł.
 */
export interface ZrodloWarsztatuTranslate {
  // --- pamięć tłumaczeń jako byt Operatora ---
  wykazPamieci(
    zadanie: RequestOf<typeof Command.TranslateMemoryList>,
  ): Promise<WynikTranslate<TranslateMemoryListResponse>>;
  zapiszPare(
    zadanie: RequestOf<typeof Command.TranslateMemorySet>,
  ): Promise<WynikTranslate<TranslateMemorySetResponse>>;
  usunPare(
    zadanie: RequestOf<typeof Command.TranslateMemoryDelete>,
  ): Promise<WynikTranslate<TranslateMemoryDeleteResponse>>;
  importujPamiec(
    zadanie: RequestOf<typeof Command.TranslateMemoryImport>,
  ): Promise<WynikTranslate<TranslateMemoryImportResponse>>;
  eksportujPamiec(
    zadanie: RequestOf<typeof Command.TranslateMemoryExport>,
  ): Promise<WynikTranslate<TranslateMemoryExportResponse>>;
  utrzymajPamiec(
    zadanie: RequestOf<typeof Command.TranslateMemoryMaintain>,
  ): Promise<WynikTranslate<TranslateMemoryMaintainResponse>>;
  tlumaczWstepnie(
    zadanie: RequestOf<typeof Command.TranslateMemoryPretranslate>,
  ): Promise<WynikTranslate<TranslateMemoryPretranslateResponse>>;
  wyrownajTeksty(
    zadanie: RequestOf<typeof Command.TranslateMemoryAlign>,
  ): Promise<WynikTranslate<TranslateMemoryAlignResponse>>;
  politykaPamieci(
    zadanie: RequestOf<typeof Command.TranslateMemoryPolicyGet>,
  ): Promise<WynikTranslate<TranslateMemoryPolicyGetResponse>>;
  ustawPolitykePamieci(
    zadanie: RequestOf<typeof Command.TranslateMemoryPolicySet>,
  ): Promise<WynikTranslate<TranslateMemoryPolicySetResponse>>;

  // --- segmentacja ---
  wykazRegulSegmentacji(
    zadanie: RequestOf<typeof Command.TranslateSegmentationRulesList>,
  ): Promise<WynikTranslate<TranslateSegmentationRulesListResponse>>;
  ustawRegulySegmentacji(
    zadanie: RequestOf<typeof Command.TranslateSegmentationRulesSet>,
  ): Promise<WynikTranslate<TranslateSegmentationRulesSetResponse>>;
  scalSegmenty(
    zadanie: RequestOf<typeof Command.TranslateSegmentMerge>,
  ): Promise<WynikTranslate<TranslateSegmentMergeResponse>>;
  podzielSegment(
    zadanie: RequestOf<typeof Command.TranslateSegmentSplit>,
  ): Promise<WynikTranslate<TranslateSegmentSplitResponse>>;

  // --- terminologia ---
  wykazTerminow(
    zadanie: RequestOf<typeof Command.TranslateGlossaryList>,
  ): Promise<WynikTranslate<TranslateGlossaryListResponse>>;
  wyjmijTerminy(
    zadanie: RequestOf<typeof Command.TranslateTermExtract>,
  ): Promise<WynikTranslate<TranslateTermExtractResponse>>;

  // --- korekta i spójność ---
  sprawdzKorekte(
    zadanie: RequestOf<typeof Command.TranslateProofreadRun>,
  ): Promise<WynikTranslate<TranslateProofreadRunResponse>>;
  zastosujKorekte(
    zadanie: RequestOf<typeof Command.TranslateProofreadApply>,
  ): Promise<WynikTranslate<TranslateProofreadApplyResponse>>;
  sprawdzSpojnosc(
    zadanie: RequestOf<typeof Command.TranslateConsistencyCheck>,
  ): Promise<WynikTranslate<TranslateConsistencyCheckResponse>>;

  // --- profile kontroli jakości i obieg zatwierdzeń ---
  wykazProfiliQa(
    zadanie: RequestOf<typeof Command.TranslateQaProfileList>,
  ): Promise<WynikTranslate<TranslateQaProfileListResponse>>;
  ustawProfilQa(
    zadanie: RequestOf<typeof Command.TranslateQaProfileSet>,
  ): Promise<WynikTranslate<TranslateQaProfileSetResponse>>;
  usunProfilQa(
    zadanie: RequestOf<typeof Command.TranslateQaProfileDelete>,
  ): Promise<WynikTranslate<TranslateQaProfileDeleteResponse>>;
  ustawZatwierdzenie(
    zadanie: RequestOf<typeof Command.TranslateApprovalSet>,
  ): Promise<WynikTranslate<TranslateApprovalSetResponse>>;
  wykazZatwierdzen(
    zadanie: RequestOf<typeof Command.TranslateApprovalList>,
  ): Promise<WynikTranslate<TranslateApprovalListResponse>>;

  // --- dokument ---
  wczytajDokument(
    zadanie: RequestOf<typeof Command.TranslateDocumentLoad>,
  ): Promise<WynikTranslate<TranslateDocumentLoadResponse>>;
  zlozDokument(
    zadanie: RequestOf<typeof Command.TranslateDocumentRender>,
  ): Promise<WynikTranslate<TranslateDocumentRenderResponse>>;
  porownajUklad(
    zadanie: RequestOf<typeof Command.TranslateDocumentLayoutCompare>,
  ): Promise<WynikTranslate<TranslateDocumentLayoutCompareResponse>>;

  // --- lokalizacja oprogramowania ---
  wczytajZasob(
    zadanie: RequestOf<typeof Command.TranslateResourceImport>,
  ): Promise<WynikTranslate<TranslateResourceImportResponse>>;
  wydajZasob(
    zadanie: RequestOf<typeof Command.TranslateResourceExport>,
  ): Promise<WynikTranslate<TranslateResourceExportResponse>>;
  zastosujFormyMnogie(
    zadanie: RequestOf<typeof Command.TranslateResourcePluralApply>,
  ): Promise<WynikTranslate<TranslateResourcePluralApplyResponse>>;
  ustawKontekstKlucza(
    zadanie: RequestOf<typeof Command.TranslateResourceKeyContextSet>,
  ): Promise<WynikTranslate<TranslateResourceKeyContextSetResponse>>;
  wczytajXliff(
    zadanie: RequestOf<typeof Command.TranslateXliffImport>,
  ): Promise<WynikTranslate<TranslateXliffImportResponse>>;

  // --- napisy i dubbing ---
  wczytajNapisy(
    zadanie: RequestOf<typeof Command.TranslateSubtitleImport>,
  ): Promise<WynikTranslate<TranslateSubtitleImportResponse>>;
  wydajNapisy(
    zadanie: RequestOf<typeof Command.TranslateSubtitleExport>,
  ): Promise<WynikTranslate<TranslateSubtitleExportResponse>>;
  sprawdzTaktowanie(
    zadanie: RequestOf<typeof Command.TranslateSubtitleTimingCheck>,
  ): Promise<WynikTranslate<TranslateSubtitleTimingCheckResponse>>;
  zlozScenariuszDubbingu(
    zadanie: RequestOf<typeof Command.TranslateDubbingScriptBuild>,
  ): Promise<WynikTranslate<TranslateDubbingScriptBuildResponse>>;

  // --- silniki, pivot i przebieg pakietowy ---
  wykazProfiliSilnikow(
    zadanie: RequestOf<typeof Command.TranslateEngineProfileList>,
  ): Promise<WynikTranslate<TranslateEngineProfileListResponse>>;
  ustawProfilSilnika(
    zadanie: RequestOf<typeof Command.TranslateEngineProfileSet>,
  ): Promise<WynikTranslate<TranslateEngineProfileSetResponse>>;
  porownajSilniki(
    zadanie: RequestOf<typeof Command.TranslateEngineCompare>,
  ): Promise<WynikTranslate<TranslateEngineCompareResponse>>;
  politykaPivota(
    zadanie: RequestOf<typeof Command.TranslatePivotPolicyGet>,
  ): Promise<WynikTranslate<TranslatePivotPolicyGetResponse>>;
  ustawPolitykePivota(
    zadanie: RequestOf<typeof Command.TranslatePivotPolicySet>,
  ): Promise<WynikTranslate<TranslatePivotPolicySetResponse>>;
  uruchomPakiet(
    zadanie: RequestOf<typeof Command.TranslateBatchRun>,
  ): Promise<WynikTranslate<TranslateBatchRunResponse>>;

  // --- wymiana zewnętrzna ---
  zlozPakietPrzekazania(
    zadanie: RequestOf<typeof Command.TranslateHandoffBuild>,
  ): Promise<WynikTranslate<TranslateHandoffBuildResponse>>;
  przyjmijZwrot(
    zadanie: RequestOf<typeof Command.TranslateHandoffReceive>,
  ): Promise<WynikTranslate<TranslateHandoffReceiveResponse>>;
  przyjmijZrodloMostu(
    zadanie: RequestOf<typeof Command.TranslateBridgeSourceReceive>,
  ): Promise<WynikTranslate<TranslateBridgeSourceReceiveResponse>>;
  odesljWynikMostu(
    zadanie: RequestOf<typeof Command.TranslateBridgeResultSend>,
  ): Promise<WynikTranslate<TranslateBridgeResultSendResponse>>;
  wydajWytwor(
    zadanie: RequestOf<typeof Command.TranslateArtifactPublish>,
  ): Promise<WynikTranslate<TranslateArtifactPublishResponse>>;
  wykazKrokow(
    zadanie: RequestOf<typeof Command.TranslateStepList>,
  ): Promise<WynikTranslate<TranslateStepListResponse>>;
}

export function utworzZrodloWarsztatuTranslate(kanal: Kanal): ZrodloWarsztatuTranslate {
  /** Komenda oddająca wykaz — wynikiem jest tablica pod wskazanym polem. */
  function zWykazem<K extends keyof typeof Command & string>(
    komenda: (typeof Command)[K],
    pole: string,
  ) {
    return (zadanie: RequestOf<(typeof Command)[K]>) =>
      zadaj(kanal, komenda, zadanie, (tresc) =>
        czyTablica((tresc as unknown as Record<string, unknown>)[pole]),
      );
  }

  /** Komenda oddająca byt — wynikiem jest obiekt pod wskazanym polem. */
  function zBytem<K extends keyof typeof Command & string>(
    komenda: (typeof Command)[K],
    pole: string,
  ) {
    return (zadanie: RequestOf<(typeof Command)[K]>) =>
      zadaj(kanal, komenda, zadanie, (tresc) =>
        czyObiekt((tresc as unknown as Record<string, unknown>)[pole]),
      );
  }

  /**
   * Komenda oddająca samą liczbę albo ścieżkę. Sprawdzianem jest obecność
   * odpowiedzi jako obiektu: pola liczbowe bywają zerem i zero jest wynikiem
   * poprawnym („nic nie zmieniono”), więc sprawdzanie ich wartości zamieniałoby
   * prawdziwy wynik w odmowę.
   */
  function zLiczba<K extends keyof typeof Command & string>(komenda: (typeof Command)[K]) {
    return (zadanie: RequestOf<(typeof Command)[K]>) =>
      zadaj(kanal, komenda, zadanie, (tresc) => czyObiekt(tresc as unknown));
  }

  return {
    wykazPamieci: zWykazem(Command.TranslateMemoryList, 'entries') as ZrodloWarsztatuTranslate['wykazPamieci'],
    zapiszPare: zBytem(Command.TranslateMemorySet, 'entry') as ZrodloWarsztatuTranslate['zapiszPare'],
    usunPare: zLiczba(Command.TranslateMemoryDelete) as ZrodloWarsztatuTranslate['usunPare'],
    importujPamiec: zLiczba(Command.TranslateMemoryImport) as ZrodloWarsztatuTranslate['importujPamiec'],
    eksportujPamiec: zLiczba(Command.TranslateMemoryExport) as ZrodloWarsztatuTranslate['eksportujPamiec'],
    utrzymajPamiec: zLiczba(Command.TranslateMemoryMaintain) as ZrodloWarsztatuTranslate['utrzymajPamiec'],
    tlumaczWstepnie: zWykazem(
      Command.TranslateMemoryPretranslate,
      'panels',
    ) as ZrodloWarsztatuTranslate['tlumaczWstepnie'],
    wyrownajTeksty: zWykazem(Command.TranslateMemoryAlign, 'pairs') as ZrodloWarsztatuTranslate['wyrownajTeksty'],
    politykaPamieci: zBytem(
      Command.TranslateMemoryPolicyGet,
      'policy',
    ) as ZrodloWarsztatuTranslate['politykaPamieci'],
    ustawPolitykePamieci: zBytem(
      Command.TranslateMemoryPolicySet,
      'policy',
    ) as ZrodloWarsztatuTranslate['ustawPolitykePamieci'],

    wykazRegulSegmentacji: zWykazem(
      Command.TranslateSegmentationRulesList,
      'rulesets',
    ) as ZrodloWarsztatuTranslate['wykazRegulSegmentacji'],
    ustawRegulySegmentacji: zBytem(
      Command.TranslateSegmentationRulesSet,
      'ruleset',
    ) as ZrodloWarsztatuTranslate['ustawRegulySegmentacji'],
    scalSegmenty: zWykazem(
      Command.TranslateSegmentMerge,
      'segments',
    ) as ZrodloWarsztatuTranslate['scalSegmenty'],
    podzielSegment: zWykazem(
      Command.TranslateSegmentSplit,
      'segments',
    ) as ZrodloWarsztatuTranslate['podzielSegment'],

    wykazTerminow: zWykazem(Command.TranslateGlossaryList, 'terms') as ZrodloWarsztatuTranslate['wykazTerminow'],
    wyjmijTerminy: zWykazem(
      Command.TranslateTermExtract,
      'candidates',
    ) as ZrodloWarsztatuTranslate['wyjmijTerminy'],

    sprawdzKorekte: zWykazem(
      Command.TranslateProofreadRun,
      'findings',
    ) as ZrodloWarsztatuTranslate['sprawdzKorekte'],
    zastosujKorekte: zBytem(
      Command.TranslateProofreadApply,
      'panel',
    ) as ZrodloWarsztatuTranslate['zastosujKorekte'],
    sprawdzSpojnosc: zWykazem(
      Command.TranslateConsistencyCheck,
      'findings',
    ) as ZrodloWarsztatuTranslate['sprawdzSpojnosc'],

    wykazProfiliQa: zWykazem(
      Command.TranslateQaProfileList,
      'profiles',
    ) as ZrodloWarsztatuTranslate['wykazProfiliQa'],
    ustawProfilQa: zBytem(Command.TranslateQaProfileSet, 'profile') as ZrodloWarsztatuTranslate['ustawProfilQa'],
    usunProfilQa: zLiczba(Command.TranslateQaProfileDelete) as ZrodloWarsztatuTranslate['usunProfilQa'],
    ustawZatwierdzenie: zBytem(
      Command.TranslateApprovalSet,
      'record',
    ) as ZrodloWarsztatuTranslate['ustawZatwierdzenie'],
    wykazZatwierdzen: zWykazem(
      Command.TranslateApprovalList,
      'records',
    ) as ZrodloWarsztatuTranslate['wykazZatwierdzen'],

    wczytajDokument: zBytem(
      Command.TranslateDocumentLoad,
      'document',
    ) as ZrodloWarsztatuTranslate['wczytajDokument'],
    zlozDokument: zLiczba(Command.TranslateDocumentRender) as ZrodloWarsztatuTranslate['zlozDokument'],
    porownajUklad: zWykazem(
      Command.TranslateDocumentLayoutCompare,
      'differences',
    ) as ZrodloWarsztatuTranslate['porownajUklad'],

    wczytajZasob: zBytem(
      Command.TranslateResourceImport,
      'resource',
    ) as ZrodloWarsztatuTranslate['wczytajZasob'],
    wydajZasob: zLiczba(Command.TranslateResourceExport) as ZrodloWarsztatuTranslate['wydajZasob'],
    zastosujFormyMnogie: zWykazem(
      Command.TranslateResourcePluralApply,
      'keys',
    ) as ZrodloWarsztatuTranslate['zastosujFormyMnogie'],
    ustawKontekstKlucza: zBytem(
      Command.TranslateResourceKeyContextSet,
      'key',
    ) as ZrodloWarsztatuTranslate['ustawKontekstKlucza'],
    wczytajXliff: zWykazem(
      Command.TranslateXliffImport,
      'panels',
    ) as ZrodloWarsztatuTranslate['wczytajXliff'],

    wczytajNapisy: zWykazem(
      Command.TranslateSubtitleImport,
      'cues',
    ) as ZrodloWarsztatuTranslate['wczytajNapisy'],
    wydajNapisy: zLiczba(Command.TranslateSubtitleExport) as ZrodloWarsztatuTranslate['wydajNapisy'],
    sprawdzTaktowanie: zWykazem(
      Command.TranslateSubtitleTimingCheck,
      'issues',
    ) as ZrodloWarsztatuTranslate['sprawdzTaktowanie'],
    zlozScenariuszDubbingu: zWykazem(
      Command.TranslateDubbingScriptBuild,
      'lines',
    ) as ZrodloWarsztatuTranslate['zlozScenariuszDubbingu'],

    wykazProfiliSilnikow: zWykazem(
      Command.TranslateEngineProfileList,
      'profiles',
    ) as ZrodloWarsztatuTranslate['wykazProfiliSilnikow'],
    ustawProfilSilnika: zBytem(
      Command.TranslateEngineProfileSet,
      'profile',
    ) as ZrodloWarsztatuTranslate['ustawProfilSilnika'],
    porownajSilniki: zWykazem(
      Command.TranslateEngineCompare,
      'variants',
    ) as ZrodloWarsztatuTranslate['porownajSilniki'],
    politykaPivota: zBytem(
      Command.TranslatePivotPolicyGet,
      'policy',
    ) as ZrodloWarsztatuTranslate['politykaPivota'],
    ustawPolitykePivota: zBytem(
      Command.TranslatePivotPolicySet,
      'policy',
    ) as ZrodloWarsztatuTranslate['ustawPolitykePivota'],
    uruchomPakiet: zLiczba(Command.TranslateBatchRun) as ZrodloWarsztatuTranslate['uruchomPakiet'],

    zlozPakietPrzekazania: zBytem(
      Command.TranslateHandoffBuild,
      'package',
    ) as ZrodloWarsztatuTranslate['zlozPakietPrzekazania'],
    przyjmijZwrot: zWykazem(
      Command.TranslateHandoffReceive,
      'panels',
    ) as ZrodloWarsztatuTranslate['przyjmijZwrot'],
    przyjmijZrodloMostu: zLiczba(
      Command.TranslateBridgeSourceReceive,
    ) as ZrodloWarsztatuTranslate['przyjmijZrodloMostu'],
    odesljWynikMostu: zLiczba(
      Command.TranslateBridgeResultSend,
    ) as ZrodloWarsztatuTranslate['odesljWynikMostu'],
    wydajWytwor: zLiczba(Command.TranslateArtifactPublish) as ZrodloWarsztatuTranslate['wydajWytwor'],
    wykazKrokow: zWykazem(Command.TranslateStepList, 'steps') as ZrodloWarsztatuTranslate['wykazKrokow'],
  };
}
