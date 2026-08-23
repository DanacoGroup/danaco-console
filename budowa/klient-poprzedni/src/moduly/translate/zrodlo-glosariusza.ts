import {
  Command,
  type TranslateGlossaryApplyRequest,
  type TranslateGlossaryApplyResponse,
  type TranslateGlossaryExportResponse,
  type TranslateGlossaryImportResponse,
  type TranslateGlossaryOccurrencesResponse,
  type TranslateGlossarySetRequest,
  type TranslateGlossarySetResponse,
} from '../../../../shared/contract';
import type { Kanal } from '../../protokol/kanal';
import { czyLiczba, czyObiekt, czyTablica } from '../../protokol/ksztalt-odpowiedzi';
import { zadaj, type WynikTranslate } from './wywolanie-translate';

/**
 * Pięć komend Glossary Managera — okna zarządcy modułu Translate.
 *
 * Kontrakt nie ma komendy odczytu glosariusza — nie niesie ani
 * `translate.glossary.list`, ani `.get`, choć edycja glosariusza wymaga
 * wczytania. Źródło nie zmyśla wykazu: okno pokazuje wyłącznie terminy zapisane
 * w tej sesji i mówi wprost, czego brakuje.
 */
export interface ZrodloGlosariusza {
  /** Definiuje odpowiednik terminu; brak `termId` zakłada termin nowy. */
  zapisz(
    zadanie: TranslateGlossarySetRequest,
  ): Promise<WynikTranslate<TranslateGlossarySetResponse>>;
  /** Ujednolica terminologię panelu (albo wszystkich paneli) wedle glosariusza. */
  ujednolic(
    idPanelu: string,
  ): Promise<WynikTranslate<TranslateGlossaryApplyResponse>>;
  /** Pokazuje wystąpienia terminu w tekście źródłowym i panelach. */
  wystapienia(
    termin: string,
  ): Promise<WynikTranslate<TranslateGlossaryOccurrencesResponse>>;
  /** Wczytuje glosariusz z pliku (TBX/CSV). */
  wczytajZPliku(
    sciezka: string,
  ): Promise<WynikTranslate<TranslateGlossaryImportResponse>>;
  /** Zapisuje glosariusz do pliku (TBX/CSV). */
  zapiszDoPliku(
    sciezka: string,
  ): Promise<WynikTranslate<TranslateGlossaryExportResponse>>;
}

export function utworzZrodloGlosariusza(kanal: Kanal): ZrodloGlosariusza {
  return {
    async zapisz(zadanie) {
      return zadaj(kanal, Command.TranslateGlossarySet, zadanie, (tresc) =>
        czyObiekt(tresc.term),
      );
    },

    async ujednolic(idPanelu) {
      // Puste `panelId` obejmuje wszystkie panele — tak stanowi kontrakt,
      // więc pola pustego nie wysyłamy zamiast wysyłać pusty napis.
      const zadanie: TranslateGlossaryApplyRequest = {};
      if (idPanelu !== '') zadanie.panelId = idPanelu;
      return zadaj(kanal, Command.TranslateGlossaryApply, zadanie, (tresc) =>
        czyLiczba(tresc.changedCount),
      );
    },

    async wystapienia(termin) {
      return zadaj(kanal, Command.TranslateGlossaryOccurrences, { term: termin }, (tresc) =>
        czyTablica(tresc.occurrences),
      );
    },

    async wczytajZPliku(sciezka) {
      return zadaj(kanal, Command.TranslateGlossaryImport, { path: sciezka }, (tresc) =>
        czyLiczba(tresc.importedCount),
      );
    },

    async zapiszDoPliku(sciezka) {
      return zadaj(kanal, Command.TranslateGlossaryExport, { path: sciezka }, (tresc) =>
        czyLiczba(tresc.exportedCount),
      );
    },
  };
}
