import {
  Command,
  EventType,
  type ExportFormat,
  type TranslateBacktranslationRunRequest,
  type TranslateBacktranslationRunResponse,
  type TranslateMemorySuggestResponse,
  type TranslatePanelExportResponse,
  type TranslatePanelToneSetResponse,
  type TranslateQualityCheckResponse,
  type TranslateSpeechSynthesizeResponse,
  type TranslateTargetAddRequest,
  type TranslateTargetAddResponse,
  type TranslateTranslationChangedEvent,
  type TranslateTranslationSetResponse,
} from '../../../../shared/contract';
import type { Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import type { Kanal } from '../../protokol/kanal';
import { czyObiekt, czyTablica, czyTekst } from '../../protokol/ksztalt-odpowiedzi';
import { zadaj, type WynikTranslate } from './wywolanie-translate';

/**
 * Osiem komend i jedno zdarzenie Translation Panels — okna o instancji wielokrotnej — mają jedno
 * źródło na wszystkie instancje, nie jedno na panel; odmowa jednego panelu zostaje w tym panelu.
 */
export interface ZrodloPaneli {
  dodajJezyk(
    zadanie: TranslateTargetAddRequest,
  ): Promise<WynikTranslate<TranslateTargetAddResponse>>;
  zapiszKorekte(
    idPanelu: string,
    tekst: string,
  ): Promise<WynikTranslate<TranslateTranslationSetResponse>>;
  /** Pusty kanał znaczy kanał czynny okna, więc pole kanału nie wchodzi do żądania wcale. */
  tlumaczZwrotnie(
    idPanelu: string,
    idKanalu: string,
  ): Promise<WynikTranslate<TranslateBacktranslationRunResponse>>;
  kontrolaJakosci(idPanelu: string): Promise<WynikTranslate<TranslateQualityCheckResponse>>;
  eksportuj(
    idPanelu: string,
    format: ExportFormat,
  ): Promise<WynikTranslate<TranslatePanelExportResponse>>;
  ustawTon(
    idPanelu: string,
    ton: string,
  ): Promise<WynikTranslate<TranslatePanelToneSetResponse>>;
  odsluchaj(idPanelu: string): Promise<WynikTranslate<TranslateSpeechSynthesizeResponse>>;
  podpowiedzPamieci(
    idPanelu: string,
    segment: string,
  ): Promise<WynikTranslate<TranslateMemorySuggestResponse>>;
  /** Subskrypcja `translate.translation.changed` — panele aktualizują się równolegle. */
  naZmiane(sluchacz: (tresc: TranslateTranslationChangedEvent) => void): Odsubskrybuj;
}

export function utworzZrodloPaneli(kanal: Kanal): ZrodloPaneli {
  return {
    async dodajJezyk(zadanie) {
      return zadaj(kanal, Command.TranslateTargetAdd, zadanie, (tresc) =>
        czyObiekt(tresc.panel),
      );
    },

    async zapiszKorekte(idPanelu, tekst) {
      return zadaj(
        kanal,
        Command.TranslateTranslationSet,
        { panelId: idPanelu, text: tekst },
        (tresc) => czyObiekt(tresc.panel),
      );
    },

    async tlumaczZwrotnie(idPanelu, idKanalu) {
      const zadanie: TranslateBacktranslationRunRequest = { panelId: idPanelu };
      if (idKanalu !== '') zadanie.channelId = idKanalu;
      return zadaj(kanal, Command.TranslateBacktranslationRun, zadanie, (tresc) =>
        czyTekst(tresc.text),
      );
    },

    async kontrolaJakosci(idPanelu) {
      // Wykaz pusty znaczy bez zastrzeżeń i jest wynikiem, nie pustką, więc sprawdzian pyta o tablicę.
      return zadaj(kanal, Command.TranslateQualityCheck, { panelId: idPanelu }, (tresc) =>
        czyTablica(tresc.issues),
      );
    },

    async eksportuj(idPanelu, format) {
      return zadaj(
        kanal,
        Command.TranslatePanelExport,
        { panelId: idPanelu, format },
        (tresc) => czyTekst(tresc.path),
      );
    },

    async ustawTon(idPanelu, ton) {
      return zadaj(
        kanal,
        Command.TranslatePanelToneSet,
        { panelId: idPanelu, tone: ton },
        (tresc) => czyObiekt(tresc.panel),
      );
    },

    async odsluchaj(idPanelu) {
      return zadaj(
        kanal,
        Command.TranslateSpeechSynthesize,
        { panelId: idPanelu },
        (tresc) => czyTekst(tresc.path),
      );
    },

    async podpowiedzPamieci(idPanelu, segment) {
      return zadaj(
        kanal,
        Command.TranslateMemorySuggest,
        { panelId: idPanelu, segment },
        (tresc) => czyTablica(tresc.suggestions),
      );
    },

    naZmiane(sluchacz) {
      return kanal.naZdarzenie(EventType.TranslateTranslationChanged, (tresc) => sluchacz(tresc));
    },
  };
}
