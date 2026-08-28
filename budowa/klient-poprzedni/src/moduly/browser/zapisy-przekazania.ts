import {
  Command,
  type ContextBundle,
  type ContextTransferResponse,
  type MessageSendResponse,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyObiekt, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../../protokol/wywolanie';

/**
 * Dwie drogi wyjścia treści z modułu Browser: przekazanie kompletu kontekstu do
 * innego modułu komendą `context.transfer` oraz pytanie o zaznaczenie, które
 * komendą `message.send` staje w rozmowie okna Browser i niesie załączniki.
 */
export interface ZapisyPrzekazania {
  /** Przenosi komplet kontekstu do modułu docelowego. */
  przekaz(
    idOkna: string,
    kodModulu: string,
    komplet: ContextBundle,
  ): Promise<Wynik<ContextTransferResponse>>;
  /**
   * Zadaje pytanie w oknie rozmowy modułu; załączniki idą polem `attachments`.
   */
  zapytaj(
    idOkna: string,
    tresc: string,
    zalaczniki?: readonly string[],
  ): Promise<Wynik<MessageSendResponse>>;
}

export function utworzZapisyPrzekazania(kanal: Kanal): ZapisyPrzekazania {
  return {
    async przekaz(idOkna, kodModulu, komplet) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.ContextTransfer, {
          sourceWindowId: idOkna,
          targetModuleId: kodModulu,
          bundle: komplet,
        }),
        Command.ContextTransfer,
        (tresc) => czyObiekt(tresc.window),
      );
    },

    async zapytaj(idOkna, tresc, zalaczniki = []) {
      // Puste pole `attachments` twierdziłoby o załączniku zerowej długości
      // tam, gdzie załącznika nie było.
      const zadanie =
        zalaczniki.length === 0
          ? { windowId: idOkna, content: tresc }
          : { windowId: idOkna, content: tresc, attachments: [...zalaczniki] };
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.MessageSend, zadanie),
        Command.MessageSend,
        (odpowiedz) => czyObiekt(odpowiedz.message),
      );
    },
  };
}
