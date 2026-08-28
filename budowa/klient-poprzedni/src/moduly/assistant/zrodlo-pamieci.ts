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
 * Pamięć i wiedza asystenta widziane przez okno Memory & Context Manager:
 * cztery komendy rodziny `memory.*` oraz dwie rodziny `knowledge.*`. Plik
 * prowadzi samą warstwę wywołań kontraktu wraz ze sprawdzianem kształtu
 * odpowiedzi i stanu nie buduje.
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
