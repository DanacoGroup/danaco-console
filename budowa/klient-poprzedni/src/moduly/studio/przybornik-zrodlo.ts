import {
  Command,
  ConfigScope,
  type StudioAnnotationAddResponse,
  type StudioAnnotationListResponse,
  type StudioOperationDeleteResponse,
  type StudioOperationListResponse,
  type StudioOperationSaveResponse,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyLogiczna, czyObiekt, czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../../protokol/wywolanie';

/**
 * Interfejs PrzybornikZrodlo woła bezpośrednio komendy adnotacji i operacji przybornika w kontrakcie, zamiast drogi generycznej window.action, która dla adnotacji zwracała odmowę not_found.
 */
export interface PrzybornikZrodlo {
  /** Zakłada adnotację przy fragmencie różnicy. */
  przybornikDodajAdnotacje(
    idDokumentu: string,
    numerFragmentu: number,
    tresc: string,
    wersje: { odniesienie?: string; porownywana?: string; propozycja?: string },
  ): Promise<Wynik<StudioAnnotationAddResponse>>;
  /** Adnotacje dokumentu w kolejności fragmentów. */
  przybornikAdnotacje(
    idDokumentu: string,
    wersje: { odniesienie?: string; porownywana?: string },
  ): Promise<Wynik<StudioAnnotationListResponse>>;
  /** Wykaz operacji łączy własne i fabryczne pod jednym wywołaniem; pole builtin rozróżnia pochodzenie. */
  przybornikOperacje(): Promise<Wynik<StudioOperationListResponse>>;
  /** Zakłada albo zmienia operację własną Operatora. */
  przybornikZapiszOperacje(
    idOperacji: string,
    nazwa: string,
    kategoria: string,
    polecenie: string,
  ): Promise<Wynik<StudioOperationSaveResponse>>;
  /** Nie usuwa operacji fabrycznej: rdzeń odmawia z nazwanym powodem, a przycisk pozostaje aktywny. */
  przybornikUsunOperacje(idOperacji: string): Promise<Wynik<StudioOperationDeleteResponse>>;
}

export function utworzPrzybornikZrodlo(kanal: Kanal): PrzybornikZrodlo {
  return {
    async przybornikDodajAdnotacje(idDokumentu, numerFragmentu, tresc, wersje) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.StudioAnnotationAdd, {
          documentId: idDokumentu,
          hunkIndex: numerFragmentu,
          body: tresc,
          ...(wersje.odniesienie === undefined || wersje.odniesienie === ''
            ? {}
            : { baseVersionId: wersje.odniesienie }),
          ...(wersje.porownywana === undefined || wersje.porownywana === ''
            ? {}
            : { targetVersionId: wersje.porownywana }),
          ...(wersje.propozycja === undefined || wersje.propozycja === ''
            ? {}
            : { proposalId: wersje.propozycja }),
        }),
        Command.StudioAnnotationAdd,
        (odpowiedz) => czyObiekt(odpowiedz.annotation),
      );
    },

    async przybornikAdnotacje(idDokumentu, wersje) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.StudioAnnotationList, {
          documentId: idDokumentu,
          ...(wersje.odniesienie === undefined || wersje.odniesienie === ''
            ? {}
            : { baseVersionId: wersje.odniesienie }),
          ...(wersje.porownywana === undefined || wersje.porownywana === ''
            ? {}
            : { targetVersionId: wersje.porownywana }),
        }),
        Command.StudioAnnotationList,
        (odpowiedz) => czyTablica(odpowiedz.annotations),
      );
    },

    async przybornikOperacje() {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.StudioOperationList, { scope: ConfigScope.Global }),
        Command.StudioOperationList,
        (odpowiedz) => czyTablica(odpowiedz.operations),
      );
    },

    async przybornikZapiszOperacje(idOperacji, nazwa, kategoria, polecenie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.StudioOperationSave, {
          ...(idOperacji === '' ? {} : { operationId: idOperacji }),
          name: nazwa,
          category: kategoria,
          prompt: polecenie,
          scope: ConfigScope.Global,
        }),
        Command.StudioOperationSave,
        (odpowiedz) => czyObiekt(odpowiedz.operation),
      );
    },

    async przybornikUsunOperacje(idOperacji) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.StudioOperationDelete, { operationId: idOperacji }),
        Command.StudioOperationDelete,
        (odpowiedz) => czyLogiczna(odpowiedz.deleted),
      );
    },
  };
}
