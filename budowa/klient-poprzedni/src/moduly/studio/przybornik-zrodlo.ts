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
 * Adnotacje przybornika — droga własna, nie generyczna.
 *
 * ── Sprawdzenie uwagi z kontraktu ───────────────────────────────────────────
 * Kontrakt przy `studio.annotation.add` niesie uwagę: „dziś czynność ta idzie
 * drogą generyczną `window.action` i wraca odmowa `not_found`, bo katalog akcji
 * nie ma jej wiersza". Sprawdzone: okno pracy z dokumentem wołało adnotację
 * właśnie tak — `window.action` z identyfikatorem `studio.diff.adnotacja`, czyli
 * nazwą, której w katalogu akcji rdzenia nie ma. Skutkiem była odmowa
 * `not_found` przy każdym naciśnięciu, mimo że **własna komenda w kontrakcie
 * jest** i ma parę: `studio.annotation.add` i `studio.annotation.list`.
 *
 * To źródło woła te dwie komendy wprost. Czy uchwyt po stronie rdzenia stoi,
 * pokaże odpowiedź: gdy go brak, wraca odmowa nazywająca komendę, a nie odmowa
 * nazywająca brakujący wiersz katalogu akcji — i to jest różnica, po której
 * Właściciel pozna, czego naprawdę brakuje. Uwaga do sprawozdania stąd
 * pochodzi.
 *
 * Źródło nie ma stanu i nie buduje elementu.
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
  /**
   * Operacje zapisane w rdzeniu — fabryczne i własne Operatora razem.
   *
   * Wykaz jest jeden, tak jak żąda uzupełnienie o narzędziach ukrytych. Pole
   * `builtin` rozdziela pochodzenie, więc przybornik odróżnia własne od
   * fabrycznych bez drugiego wykazu.
   */
  przybornikOperacje(): Promise<Wynik<StudioOperationListResponse>>;
  /** Zakłada albo zmienia operację własną Operatora. */
  przybornikZapiszOperacje(
    idOperacji: string,
    nazwa: string,
    kategoria: string,
    polecenie: string,
  ): Promise<Wynik<StudioOperationSaveResponse>>;
  /**
   * Usuwa operację własną.
   *
   * Operacji fabrycznej nie usuwa: rdzeń odpowiada odmową nazywającą powód i tak
   * ma zostać. Klient tej odmowy nie uprzedza wyłączeniem przycisku — Operator
   * ma usłyszeć powód od rdzenia, a nie domyślać się go z martwej kontrolki.
   */
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
