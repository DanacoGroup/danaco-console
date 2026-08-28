import {
  Command,
  type TranslateSourceDetectRequest,
  type TranslateSourceDetectResponse,
  type TranslateSourceSegmentRequest,
  type TranslateSourceSegmentResponse,
  type TranslateSourceSetRequest,
  type TranslateSourceSetResponse,
} from '../../../../shared/contract';
import type { Kanal } from '../../protokol/kanal';
import { czyTablica, czyTekst } from '../../protokol/ksztalt-odpowiedzi';
import { zadaj, type WynikTranslate } from './wywolanie-translate';

/**
 * Trzy komendy panelu źródła okna wiodącego modułu, stanowiące wyłącznie warstwę
 * wywołania i sprawdzianu kształtu odpowiedzi rdzenia, bez własnego stanu.
 */
export interface ZrodloZrodla {
  /** Ustawia tekst źródłowy; wyzwala jednoczesną aktualizację wszystkich paneli. */
  ustaw(
    zadanie: TranslateSourceSetRequest,
  ): Promise<WynikTranslate<TranslateSourceSetResponse>>;
  /** Dzieli tekst źródłowy na segmenty ponownie. */
  segmentuj(tekst: string): Promise<WynikTranslate<TranslateSourceSegmentResponse>>;
  /** Rozpoznaje język tekstu źródłowego. */
  rozpoznajJezyk(tekst: string): Promise<WynikTranslate<TranslateSourceDetectResponse>>;
}

export function utworzZrodloZrodla(kanal: Kanal): ZrodloZrodla {
  return {
    async ustaw(zadanie) {
      return zadaj(
        kanal,
        Command.TranslateSourceSet,
        zadanie,
        (tresc) => czyTekst(tresc.sourceLanguage) && czyTablica(tresc.panels),
      );
    },

    async segmentuj(tekst) {
      // Puste pole nie jest wysyłane, żeby nie skasować nim tekstu bieżącego w rdzeniu.
      const zadanie: TranslateSourceSegmentRequest = {};
      if (tekst.trim() !== '') zadanie.text = tekst;
      return zadaj(kanal, Command.TranslateSourceSegment, zadanie, (tresc) =>
        czyTablica(tresc.segments),
      );
    },

    async rozpoznajJezyk(tekst) {
      const zadanie: TranslateSourceDetectRequest = {};
      if (tekst.trim() !== '') zadanie.text = tekst;
      return zadaj(kanal, Command.TranslateSourceDetect, zadanie, (tresc) =>
        czyTekst(tresc.language),
      );
    },
  };
}
