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
 * Trzy komendy Source Panel — okna wiodącego modułu Translate.
 *
 * Źródło nie ma własnego stanu i nie buduje ani jednego elementu: jest warstwą
 * wywołania i sprawdzianu kształtu odpowiedzi. Tekst źródłowy mieszka
 * w `stan-translate.ts`, żeby Translation Panels patrzyły na ten sam tekst,
 * a nie na własną kopię.
 *
 * Wszystkie trzy komendy rdzeń rejestruje i obsługuje: `source.set` oddaje
 * `{sourceLanguage, segmentCount, panels}`, `source.segment` → `{segments:[…]}`,
 * `source.detect` → `{language}`. Bez zalogowanego modelu rdzeń oddaje stan
 * zdegradowany (`language` niesie wtedy komunikat „Not logged in") — to jest
 * odpowiedź rdzenia, nie brak uchwytu. Ścieżka odmowy
 * w `wywolanie-translate.ts` zostaje na wypadek starszego rdzenia albo
 * pośrednika.
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
      // Kontrakt dopuszcza żądanie bez tekstu — wtedy rdzeń bierze tekst
      // bieżący. Pustego pola nie wysyłamy, żeby nie kasować nim treści.
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
