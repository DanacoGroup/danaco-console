import {
  Command,
  KnowledgeScope,
  type KnowledgeHit,
} from '../../../../shared/contract';
import type { Wynik } from '../../protokol/kanal';
import { czyLiczba, czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import type { StrazOdmow } from './straz-odmow';

/**
 * Dwie komendy rodziny `knowledge.*` w zakresie biblioteki — wyszukiwanie po
 * znaczeniu i budowa wskaźnika znaczenia.
 *
 * Ta sama rodzina ma drugiego klienta w module Wiedza
 * (`moduly/wiedza/zrodlo-wiedzy.ts`), z którego moduł Library świadomie nie
 * korzysta z dwóch powodów. Pierwszy jest własnościowy: moduł nie sięga do
 * plików innego modułu, bo wiązałby swoje okna z cudzym cyklem życia. Drugi
 * jest techniczny i ważniejszy — tamto źródło woła kanał wprost, a rodzina
 * `knowledge.*` nie ma własnego zdarzenia odmowy i rdzeń bez wpiętego portu
 * Wiedzy odpowiada na nią `connection.unknown`, czyli kopertą bez pola
 * `status`. Library idzie więc przez straż odmów, która taką odpowiedź wiąże
 * z żądaniem po `requestId` i zamienia w zwykły `Wynik` z błędem
 * (`straz-odmow.ts`).
 *
 * Rodzina jest poza obszarem `library`, ale czyta dokładnie ten zbiór:
 * czytelnikiem zakresu `library` jest repozytorium modułu Library, a pole
 * `sourceId` trafienia niesie identyfikator pliku biblioteki
 * (`server/internal/core/adapter_modul_wiedza_zrodla.go`). Dzięki temu trafienie
 * daje się odwzorować na wiersz wykazu bez drugiego odczytu.
 *
 * Rozdział zdolności między rodzinami jest wiążący dla okna: `library.file.search`
 * dopasowuje SŁOWA (indeks pełnotekstowy repozytorium), `knowledge.search`
 * dopasowuje ZNACZENIE (wskaźnik osadzeń). Tryb hybrydowy okna nie jest komendą
 * kontraktu — to złożenie obu odpowiedzi po stronie klienta.
 *
 * Wskaźnik nie odświeża się przy wgraniu pliku: rdzeń buduje go wyłącznie na
 * żądanie `knowledge.index`, bo osadzanie potrafi trwać minutę i sięga po wagi
 * modelu. Reindeksacja jest więc czynnością Operatora w zakładce Higiena, a nie
 * skutkiem ubocznym wgrania.
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
