import {
  Command,
  KnowledgeScope,
  type KnowledgeHit,
} from '../../../../shared/contract';
import type { Wynik } from '../../protokol/kanal';
import { czyLiczba, czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import type { StrazOdmow } from './straz-odmow';

/**
 * Dwie komendy rodziny `knowledge.*` w zakresie biblioteki: wyszukiwanie po
 * znaczeniu i budowa wskaźnika osadzeń. Obie idą przez straż odmów, bo rdzeń
 * bez portu Wiedzy odpowiada kopertą bez pola `status`, której
 * korelacja żądań nie rozstrzyga.
 */
export interface ZrodloZnaczenia {
  /** Fragmenty biblioteki najbliższe pytaniu wraz z identyfikatorem pliku. */
  szukajZnaczeniem(
    pytanie: string,
    granica: number,
  ): Promise<Wynik<{ results: KnowledgeHit[]; total: number }>>;
  /** Buduje wskaźnik znaczenia dla zakresu biblioteki. */
  przelicz(
    odNowa: boolean,
  ): Promise<Wynik<{ indexed: number; total: number; model?: string }>>;
}

export function utworzZrodloZnaczenia(straz: StrazOdmow): ZrodloZnaczenia {
  return {
    async szukajZnaczeniem(pytanie, granica) {
      return sprawdzKsztalt(
        await straz.wywolaj(Command.KnowledgeSearch, {
          query: pytanie,
          scope: KnowledgeScope.Library,
          limit: granica,
        }),
        Command.KnowledgeSearch,
        (tresc) => czyTablica(tresc.results) && czyLiczba(tresc.total),
      );
    },

    async przelicz(odNowa) {
      return sprawdzKsztalt(
        await straz.wywolaj(Command.KnowledgeIndex, {
          scope: KnowledgeScope.Library,
          rebuild: odNowa,
        }),
        Command.KnowledgeIndex,
        (tresc) => czyLiczba(tresc.indexed) && czyLiczba(tresc.total),
      );
    },
  };
}
