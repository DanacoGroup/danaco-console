import {
  Command,
  type KnowledgeIndexRequest,
  type KnowledgeIndexResponse,
  type KnowledgeSearchRequest,
  type KnowledgeSearchResponse,
  type MemoryDeleteResponse,
  type MemoryListRequest,
  type MemoryListResponse,
  type MemorySetRequest,
  type MemorySetResponse,
  type MemoryToggleRequest,
  type MemoryToggleResponse,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyLiczba, czyObiekt, czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../../protokol/wywolanie';

/**
 * Pamięć i wiedza asystenta widziane przez okno Memory & Context Manager —
 * cztery komendy rodziny `memory.*` i dwie rodziny `knowledge.*`.
 *
 * Plik odpowiada wyłącznie za warstwę wywołań kontraktu wraz ze sprawdzianem
 * kształtu odpowiedzi; stanu nie ma i nie buduje ani jednego elementu. Wzorem
 * jest `zrodlo-assistant.ts`: żadne wywołanie nie rzuca wyjątkiem, odmowa
 * wraca polem `blad` wyniku, a okno pokazuje ją w swoim stanie błędu.
 *
 * `memory.detach` nie jest tu wystawiony. Kontrakt oznacza znaczenie odpięcia
 * jako nierozstrzygnięte przez Właściciela („rdzeń czeka na decyzję"), więc
 * okno nie nadaje mu własnego sensu — tak samo postępuje moduł Workspace
 * (`moduly/workspace/pamiec-pozycja.ts`).
 *
 * Rodziny są dwie, bo mówią o dwóch różnych bytach. `memory.*` prowadzi
 * ustalenia — zdania, które Operator kazał zapamiętać. `knowledge.*` prowadzi
 * wskaźnik znaczenia zbudowany z treści już istniejących: biblioteki, historii
 * rozmów i plików przestrzeni roboczej. Zlanie ich w jedno źródło zatarłoby, co
 * jest ustaleniem, a co odnalezionym fragmentem.
 */
export interface ZrodloPamieci {
  /** `memory.list` — wpisy pamięci widoczne w zasięgu. */
  wpisy(zadanie: MemoryListRequest): Promise<Wynik<MemoryListResponse>>;
  /** `memory.set` — zapis albo zmiana ustalenia. */
  zapisz(zadanie: MemorySetRequest): Promise<Wynik<MemorySetResponse>>;
  /** `memory.delete` — usunięcie wpisu wskazanego identyfikatorem. */
  usun(idWpisu: string): Promise<Wynik<MemoryDeleteResponse>>;
  /** `memory.toggle` — poziomy pamięci karty sesji oraz zapis pamięci. */
  przestaw(zadanie: MemoryToggleRequest): Promise<Wynik<MemoryToggleResponse>>;
  /** `knowledge.search` — fragmenty odnalezione po znaczeniu. */
  szukaj(zadanie: KnowledgeSearchRequest): Promise<Wynik<KnowledgeSearchResponse>>;
  /** `knowledge.index` — budowa albo przebudowa wskaźnika znaczenia. */
  wskaznik(zadanie: KnowledgeIndexRequest): Promise<Wynik<KnowledgeIndexResponse>>;
}

export function utworzZrodloPamieci(kanal: Kanal): ZrodloPamieci {
  return {
    async wpisy(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.MemoryList, zadanie),
        Command.MemoryList,
        (tresc) => czyTablica(tresc.entries),
      );
    },

    async zapisz(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.MemorySet, zadanie),
        Command.MemorySet,
        (tresc) => czyObiekt(tresc.entry),
      );
    },

    async usun(idWpisu) {
      return wywolaj(kanal, Command.MemoryDelete, { entryId: idWpisu });
    },

    async przestaw(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.MemoryToggle, zadanie),
        Command.MemoryToggle,
        (tresc) => czyTablica(tresc.levels),
      );
    },

    async szukaj(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.KnowledgeSearch, zadanie),
        Command.KnowledgeSearch,
        (tresc) => czyTablica(tresc.results),
      );
    },

    async wskaznik(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.KnowledgeIndex, zadanie),
        Command.KnowledgeIndex,
        (tresc) => czyLiczba(tresc.indexed) && czyLiczba(tresc.total),
      );
    },
  };
}
