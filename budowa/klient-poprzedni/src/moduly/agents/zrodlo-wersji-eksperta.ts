import {
  Command,
  type Agent,
  type AgentArchiveListRequest,
  type AgentArchiveRequest,
  type AgentRestoreRequest,
  type AgentVersion,
  type AgentVersionListRequest,
  type AgentVersionRestoreRequest,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyObiekt, czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../../protokol/wywolanie';

/**
 * Pięć komend historii tożsamości i archiwum eksperta — warstwa wywołań. Żadne
 * wywołanie nie rzuca wyjątkiem ani nie odrzuca obietnicy: niepowodzenie wraca
 * polem `blad` wyniku, więc okno rozstrzyga o nim tak samo jak o powodzeniu.
 */
export interface ZrodloWersjiEksperta {
  /** Historia tożsamości od wersji najnowszej. */
  wersje(
    idEksperta: string,
    ile?: number,
  ): Promise<Wynik<{ agentId: string; versions: AgentVersion[]; total: number }>>;
  /** Przywraca wskazaną wersję jako kolejną — historii nie skraca. */
  przywrocWersje(
    idEksperta: string,
    idWersji: string,
  ): Promise<Wynik<{ agent: Agent; version: AgentVersion }>>;
  /** Odkłada eksperta poza wykaz czynnych bez utraty definicji i historii. */
  zarchiwizuj(idEksperta: string): Promise<Wynik<{ agentId: string; archived: boolean }>>;
  /** Oddaje eksperta wykazowi czynnych wraz z zapamiętanym stanem czynności. */
  przywrocZArchiwum(idEksperta: string): Promise<Wynik<{ agentId: string; restored: boolean }>>;
  /** Wykaz zarchiwizowanych — osobny od biblioteki czynnej. */
  archiwum(ile?: number): Promise<Wynik<{ agents: Agent[]; total: number }>>;
}

export function utworzZrodloWersjiEksperta(kanal: Kanal): ZrodloWersjiEksperta {
  return {
    async wersje(idEksperta, ile) {
      const zadanie: AgentVersionListRequest = { agentId: idEksperta };
      if (ile !== undefined) zadanie.limit = ile;
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AgentVersionList, zadanie),
        Command.AgentVersionList,
        (tresc) => czyTablica(tresc.versions),
      );
    },

    async przywrocWersje(idEksperta, idWersji) {
      // Wersję wskazuje identyfikator, nie numer; numer jest porządkiem historii.
      const zadanie: AgentVersionRestoreRequest = { agentId: idEksperta, versionId: idWersji };
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AgentVersionRestore, zadanie),
        Command.AgentVersionRestore,
        (tresc) => czyObiekt(tresc.agent) && czyObiekt(tresc.version),
      );
    },

    async zarchiwizuj(idEksperta) {
      const zadanie: AgentArchiveRequest = { agentId: idEksperta };
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AgentArchive, zadanie),
        Command.AgentArchive,
        (tresc) => typeof tresc.archived === 'boolean',
      );
    },

    async przywrocZArchiwum(idEksperta) {
      const zadanie: AgentRestoreRequest = { agentId: idEksperta };
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AgentRestore, zadanie),
        Command.AgentRestore,
        (tresc) => typeof tresc.restored === 'boolean',
      );
    },

    async archiwum(ile) {
      const zadanie: AgentArchiveListRequest = {};
      if (ile !== undefined) zadanie.limit = ile;
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AgentArchiveList, zadanie),
        Command.AgentArchiveList,
        (tresc) => czyTablica(tresc.agents),
      );
    },
  };
}

/**
 * Zdanie o wersji dla wykazu: etykieta, sprawca i chwila.
 *
 * Wersja nie niesie migawki tożsamości — rodzina oddaje sumę kontrolną, a treść
 * pobiera się przy przywróceniu. Wykaz pokazuje więc wyłącznie to, co odróżnia
 * wersje od siebie.
 */
export function zdanieOWersji(wersja: AgentVersion): string {
  const czesci = [wersja.label ?? wersja.id];
  if (wersja.summary !== undefined && wersja.summary !== '') czesci.push(wersja.summary);
  if (wersja.author !== undefined && wersja.author !== '') czesci.push(`sprawca: ${wersja.author}`);
  czesci.push(new Date(wersja.createdAt).toLocaleString('pl-PL'));
  return czesci.join(' · ');
}
