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

/**
 * Wyszukiwanie po znaczeniu widziane przez klienta — dwie komendy `knowledge.*`.
 *
 * Stała `Command.*` pada wyłącznie w plikach `zrodlo-*.ts` modułu. Okno zna
 * czynność („poszukaj", „przebuduj wskaźnik"), a nie nazwę komendy, więc zmiana
 * nazwy w kontrakcie przerywa kompilację w jednym pliku, a nie w każdym widoku
 * z osobna.
 *
 * Obie komendy siedzą razem, choć robią co innego: `knowledge.search` jest
 * odczytem, a `knowledge.index` przebudową wskaźnika — ale bez wskaźnika odczyt
 * nie ma czego oddać, więc okno pokazujące wyniki potrzebuje pod ręką drogi do
 * przebudowy. Rozdzielenie ich dałoby oknu dwie zależności na jedną dziedzinę.
 *
 * Źródło oddaje `Wynik`, nie samą treść, i nie ma tu ani jednego `?? []`.
 * Wyszukiwanie, które po odmowie rdzenia oddaje pustą tablicę, mówi „nic nie
 * znalazłem" zamiast „nie udało się zapytać" — a to są dwa różne zdania i tylko
 * jedno z nich jest prawdziwe.
 *
 * Kontrakt nie niesie zdarzenia zmiany wskaźnika wiedzy — `ZDARZENIA` nie ma
 * pozycji `knowledge.*` — więc okno nie ma się na czym zawiesić i odświeża się
 * wyłącznie na czynność Operatora. Nasłuch stanąłby na zdarzeniu, które nigdy
 * nie przychodzi.
 */

/** Owoc przeszukania: fragmenty od najtrafniejszego wraz z ich liczbą. */
export interface OwocSzukania {
  /** Fragmenty wraz ze wskazaniem źródła — bez źródła trafienie jest bezwartościowe. */
  wyniki: readonly KnowledgeHit[];
  /** Liczba fragmentów spełniających warunki; bywa większa niż długość wykazu. */
  wszystkich: number;
}

/** Owoc przebudowy wskaźnika: ile weszło i ile jest po przebiegu. */
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
      // Sprawdzamy `results`, bo to jedyne pole, po którym okno iteruje —
      // rdzeń starszej wersji przysyłający treść bez niego wywróciłby pętlę
      // na pierwszym kroku, zamiast oddać uczciwe niepowodzenie.
      const sprawdzony = sprawdzKsztalt(surowy, Command.KnowledgeSearch, (tresc) =>
        czyTablica(tresc.results),
      );
      return przenies(sprawdzony, (tresc) => ({
        wyniki: tresc.results,
        // `total` bywa większe niż wykaz (rdzeń przycina do `limit`) — okno ma
        // to powiedzieć, więc licznik nie może zostać zastąpiony długością.
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
