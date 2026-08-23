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
 * Dwie drogi wyjścia treści z modułu Browser: przekazanie międzymodułowe
 * i pytanie zadane w oknie rozmowy.
 *
 * Jedna odpowiedzialność: ścieżki zapisu wychodzące poza obszar `browser.*`.
 * Trzymane osobno od odczytu (`zrodlo-browser.ts`), bo dotyczą innych obszarów
 * kontraktu i mają innego adresata.
 *
 * Przekazanie międzymodułowe („→ Wyślij do Research", „Tłumacz") idzie
 * `context.transfer` — to jedyna droga przeniesienia kompletu kontekstu między
 * modułami, jaką niesie kontrakt i jaką obsługuje rdzeń.
 *
 * Pytanie o zaznaczenie („Wyjaśnij" na pasku pływającym) idzie `message.send`.
 * Oknem, do którego pytanie trafia, jest okno modułu Browser — to samo, którego
 * identyfikator niosą komendy `browser.*`.
 *
 * Adnotacja spłaszczona do PNG idzie tą samą drogą co pytanie: `message.send`
 * kładzie treść w rozmowie wskazanego okna i niesie pole `attachments`,
 * a `context.transfer` wymaga `targetModuleId` i wysłałby rysunek poza moduł,
 * w którym powstał, nie stawiając go w rozmowie wcale.
 *
 * `MessageSendRequest.attachments` jest polem kontraktu: rdzeń przepisuje je do
 * zakładanej wiadomości (`adapter_rozmowa.go`, `Attachments: z.Attachments`)
 * i oddaje w odpowiedzi, dzięki czemu okno ocenia dołączenie po wierszu, który
 * wrócił, a nie po tym, że wysłało (`skutek-zapisu.ts`).
 */
export interface ZapisyPrzekazania {
  /** Przenosi komplet kontekstu do modułu docelowego. */
  przekaz(
    idOkna: string,
    kodModulu: string,
    komplet: ContextBundle,
  ): Promise<Wynik<ContextTransferResponse>>;
  /**
   * Zadaje pytanie w oknie rozmowy modułu; `zalaczniki` idą polem
   * `attachments` komendy i są puste, gdy pytanie niesie sam tekst.
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
      // Pole `attachments` wychodzi tylko wtedy, gdy naprawdę coś niesie:
      // pusta tablica w żądaniu twierdziłaby o załączniku zerowej długości
      // tam, gdzie załącznika nie było wcale.
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
