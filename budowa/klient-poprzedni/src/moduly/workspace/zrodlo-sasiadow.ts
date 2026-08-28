import {
  Command,
  type AgentPermission,
  type AgentPermissionGroup,
  type ContextTransferResponse,
  type LibraryFile,
  type LibraryVersion,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyObiekt, czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../../protokol/wywolanie';
import { przenies } from '../../protokol/wynik-czastkowy';
import { KOD_MODULU } from './wynik-czastkowy';

/**
 * Plik gromadzi czynności okien Workspace sięgające poza własny obszar modułu — wgranie pliku i jego wersje z modułu Library, uprawnienie eksperta z modułu Agents oraz udostępnienie materiału przy przenoszeniu kontekstu.
 */
export interface CzynnosciSasiednie {
  wgrajPlik(idProjektu: string, nazwa: string, tresc: string): Promise<Wynik<LibraryFile>>;
  wersjePliku(idPliku: string): Promise<Wynik<LibraryVersion[]>>;
  nadajEtykiete(idPliku: string, etykiety: string[]): Promise<Wynik<LibraryFile>>;
  udostepnij(
    idOkna: string,
    modulDocelowy: string,
    idProjektu: string,
    idPlikow: string[],
  ): Promise<Wynik<ContextTransferResponse>>;
  ustawUprawnienie(
    idEksperta: string,
    grupa: AgentPermissionGroup,
    przyznane: boolean,
  ): Promise<Wynik<AgentPermission[]>>;
}

export function czynnosciSasiednie(kanal: Kanal): CzynnosciSasiednie {
  return {
    async wgrajPlik(idProjektu, nazwa, tresc) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.LibraryFileUpload, {
          name: nazwa,
          projectId: idProjektu,
          sourceModuleId: KOD_MODULU,
          contentBase64: tresc,
        }),
        Command.LibraryFileUpload,
        (odpowiedz) => czyObiekt(odpowiedz.file),
      );
      return przenies(wynik, (odpowiedz) => odpowiedz.file);
    },

    async wersjePliku(idPliku) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.LibraryVersionList, { fileId: idPliku }),
        Command.LibraryVersionList,
        (tresc) => czyTablica(tresc.versions),
      );
      return przenies(wynik, (tresc) => tresc.versions);
    },

    async nadajEtykiete(idPliku, etykiety) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.LibraryTagSet, { fileId: idPliku, tags: etykiety }),
        Command.LibraryTagSet,
        (odpowiedz) => czyObiekt(odpowiedz.file),
      );
      return przenies(wynik, (odpowiedz) => odpowiedz.file);
    },

    udostepnij(idOkna, modulDocelowy, idProjektu, idPlikow) {
      return wywolaj(kanal, Command.ContextTransfer, {
        sourceWindowId: idOkna,
        targetModuleId: modulDocelowy,
        bundle: { projectId: idProjektu, documentIds: idPlikow },
      });
    },

    async ustawUprawnienie(idEksperta, grupa, przyznane) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.AgentPermissionSet, {
          agentId: idEksperta,
          group: grupa,
          granted: przyznane,
        }),
        Command.AgentPermissionSet,
        (tresc) => czyTablica(tresc.permissions),
      );
      return przenies(wynik, (tresc) => tresc.permissions);
    },
  };
}
