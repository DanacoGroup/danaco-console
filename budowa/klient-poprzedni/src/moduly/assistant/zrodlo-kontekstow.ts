import {
  Command,
  type ConfigScope,
  type ContextUsageGetRequest,
  type ContextUsageGetResponse,
  type MemoryContextActivateResponse,
  type MemoryContextDeleteResponse,
  type MemoryContextListResponse,
  type MemoryContextSaveRequest,
  type MemoryContextSaveResponse,
  type MemoryRetentionGetResponse,
  type MemoryRetentionSetRequest,
  type MemoryRetentionSetResponse,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyLiczba, czyLogiczna, czyObiekt, czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../../protokol/wywolanie';

/** Nazwane konteksty pamięci, zasady retencji i pomiar zajętości okna kontekstu — siedem komend. */

/** Kontekst zapisywany z okna niesie nazwę, opis, poziomy pamięci, wpisy, prompt systemowy i stan czynności; puste id zakłada nowy. */
export interface KontekstDoZapisu {
  id: string;
  nazwa: string;
  opis: string;
  poziomy: readonly ConfigScope[];
  wpisy: readonly string[];
  promptSystemowy: string;
  czynny: boolean;
}

/** Zasada retencji zapisywana z okna niesie zasięg, dni do wygasania, domyślną wrażliwość i wzorce chronione. */
export interface ZasadaDoZapisu {
  zasieg?: ConfigScope;
  dniWygasania: number;
  wrazliweDomyslnie: boolean;
  wzorceNigdyNieZapisywane: readonly string[];
  czynna: boolean;
}

export interface ZrodloKontekstow {
  /** `memory.context.list` — nazwane zestawy wraz z kontekstem czynnym karty. */
  konteksty(idSesji: string, zWylaczonymi: boolean): Promise<Wynik<MemoryContextListResponse>>;
  /** `memory.context.save` — założenie albo zmiana zestawu. */
  zapiszKontekst(kontekst: KontekstDoZapisu): Promise<Wynik<MemoryContextSaveResponse>>;
  /** `memory.context.activate` — podmiana poziomów pamięci i warstwy promptu. */
  uaktywnijKontekst(id: string, idSesji: string): Promise<Wynik<MemoryContextActivateResponse>>;
  /** `memory.context.delete` — usunięcie wskazania; wpisów pamięci nie rusza. */
  usunKontekst(id: string): Promise<Wynik<MemoryContextDeleteResponse>>;
  /** `memory.retention.get` — zasady obowiązujące w zasięgu, od najwęższej. */
  zasadyRetencji(): Promise<Wynik<MemoryRetentionGetResponse>>;
  /** `memory.retention.set` — zapis zasady wraz z liczbą wpisów zastanych. */
  zapiszZasadeRetencji(zasada: ZasadaDoZapisu): Promise<Wynik<MemoryRetentionSetResponse>>;
  /** `context.usage.get` — zajętość okna kontekstu rozmowy. */
  zajetoscKontekstu(idOkna: string, idSesji: string): Promise<Wynik<ContextUsageGetResponse>>;
}

export function utworzZrodloKontekstow(kanal: Kanal): ZrodloKontekstow {
  return {
    async konteksty(idSesji, zWylaczonymi) {
      const zadanie: Record<string, unknown> = { includeDisabled: zWylaczonymi };
      if (idSesji !== '') zadanie['sessionId'] = idSesji;
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.MemoryContextList, zadanie),
        Command.MemoryContextList,
        (tresc) => czyTablica(tresc.contexts),
      );
    },

    async zapiszKontekst(kontekst) {
      const zadanie: MemoryContextSaveRequest = {
        name: kontekst.nazwa,
        enabled: kontekst.czynny,
      };
      if (kontekst.id !== '') zadanie.contextId = kontekst.id;
      if (kontekst.opis.trim() !== '') zadanie.description = kontekst.opis.trim();
      if (kontekst.poziomy.length > 0) zadanie.levels = [...kontekst.poziomy];
      if (kontekst.wpisy.length > 0) zadanie.entryIds = [...kontekst.wpisy];
      if (kontekst.promptSystemowy.trim() !== '') {
        zadanie.systemPrompt = kontekst.promptSystemowy.trim();
      }
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.MemoryContextSave, zadanie),
        Command.MemoryContextSave,
        (tresc) => czyObiekt(tresc.context),
      );
    },

    async uaktywnijKontekst(id, idSesji) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.MemoryContextActivate, { contextId: id, sessionId: idSesji }),
        Command.MemoryContextActivate,
        (tresc) => czyLogiczna(tresc.activated) && czyTablica(tresc.levels),
      );
    },

    async usunKontekst(id) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.MemoryContextDelete, { contextId: id }),
        Command.MemoryContextDelete,
        (tresc) => czyLogiczna(tresc.deleted),
      );
    },

    async zasadyRetencji() {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.MemoryRetentionGet, {}),
        Command.MemoryRetentionGet,
        (tresc) => czyTablica(tresc.policies),
      );
    },

    async zapiszZasadeRetencji(zasada) {
      const zadanie: MemoryRetentionSetRequest = {
        ttlDays: zasada.dniWygasania,
        sensitiveDefault: zasada.wrazliweDomyslnie,
        enabled: zasada.czynna,
      };
      if (zasada.zasieg !== undefined) zadanie.scope = zasada.zasieg;
      if (zasada.wzorceNigdyNieZapisywane.length > 0) {
        zadanie.neverStorePatterns = [...zasada.wzorceNigdyNieZapisywane];
      }
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.MemoryRetentionSet, zadanie),
        Command.MemoryRetentionSet,
        (tresc) => czyObiekt(tresc.policy) && czyLiczba(tresc.affectedEntries),
      );
    },

    async zajetoscKontekstu(idOkna, idSesji) {
      const zadanie: ContextUsageGetRequest = { windowId: idOkna };
      if (idSesji !== '') zadanie.sessionId = idSesji;
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.ContextUsageGet, zadanie),
        Command.ContextUsageGet,
        (tresc) => czyLogiczna(tresc.available),
      );
    },
  };
}
