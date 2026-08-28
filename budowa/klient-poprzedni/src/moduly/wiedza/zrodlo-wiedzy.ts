import {
  Command,
  type KnowledgeHit,
  type KnowledgeIndexRequest,
  type KnowledgeSearchRequest,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyLiczba, czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { przenies } from '../../protokol/wynik-czastkowy';
import { wywolaj } from '../../protokol/wywolanie';

/** Wyszukiwanie po znaczeniu widziane przez klienta obejmuje dwie komendy: szukanie i przebudowę. */

/** Owoc przeszukania niesie fragmenty od najtrafniejszego wraz z ich łączną liczbą spełniającą warunki zapytania. */
export interface OwocSzukania {
  /** Fragmenty wraz ze wskazaniem źródła — bez źródła trafienie jest bezwartościowe. */
  wyniki: readonly KnowledgeHit[];
  /** Liczba fragmentów spełniających warunki; bywa większa niż długość wykazu. */
  wszystkich: number;
}

/** Owoc przebudowy wskaźnika niesie liczbę pozycji wniesionych w tym przebiegu oraz liczbę pozycji po przebiegu. */
export interface OwocWskaznika {
  /** Liczba pozycji wprowadzonych do wskaźnika w tym przebiegu. */
  wniesione: number;
  /** Liczba pozycji we wskaźniku po przebiegu. */
  wszystkich: number;
  /** Model osadzeń, którym zbudowano wskaźnik; rdzeń może go nie podać. */
  model?: string;
}

export interface ZrodloWiedzy {
  /** `knowledge.search` — pytanie w języku Operatora, fragmenty ze źródłem. */
  szukaj(zadanie: KnowledgeSearchRequest): Promise<Wynik<OwocSzukania>>;
  /** `knowledge.index` — budowa albo przebudowa wskaźnika znaczenia. */
  przebudujWskaznik(zadanie: KnowledgeIndexRequest): Promise<Wynik<OwocWskaznika>>;
}

export function utworzZrodloWiedzy(kanal: Kanal): ZrodloWiedzy {
  return {
    async szukaj(zadanie) {
      const surowy = await wywolaj(kanal, Command.KnowledgeSearch, zadanie);
      // Sprawdzamy pole wyników, jedyne pole, po którym okno iteruje, inaczej pętla wywróciłaby się od razu.
      const sprawdzony = sprawdzKsztalt(surowy, Command.KnowledgeSearch, (tresc) =>
        czyTablica(tresc.results),
      );
      return przenies(sprawdzony, (tresc) => ({
        wyniki: tresc.results,
        // Liczba całkowita bywa większa niż wykaz, bo rdzeń przycina do limitu; licznik jej nie zastępuje.
        wszystkich: czyLiczba(tresc.total) ? tresc.total : tresc.results.length,
      }));
    },

    async przebudujWskaznik(zadanie) {
      const surowy = await wywolaj(kanal, Command.KnowledgeIndex, zadanie);
      const sprawdzony = sprawdzKsztalt(surowy, Command.KnowledgeIndex, (tresc) =>
        czyLiczba(tresc.indexed) && czyLiczba(tresc.total),
      );
      return przenies(sprawdzony, (tresc) => ({
        wniesione: tresc.indexed,
        wszystkich: tresc.total,
        ...(tresc.model === undefined ? {} : { model: tresc.model }),
      }));
    },
  };
}
