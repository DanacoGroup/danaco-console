import {
  Command,
  EventType,
  type Envelope,
  type StreamChunkEvent,
  type TerminalOutputStreamRequest,
  type TerminalOutputStreamResponse,
} from '../../../shared/contract';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import type { Kanal, Wynik } from '../protokol/kanal';
import { czyLogiczna, czyTablica, sprawdzKsztalt } from '../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../protokol/wywolanie';

/**
 * Źródło podglądu w tle — jedna komenda kontraktu (`terminal.output.stream`)
 * i jedna subskrypcja (`stream.chunk`).
 *
 * `shared/contract.ts` niesie `Command.TerminalOutputStream`, a rdzeń rejestruje
 * dla niej uchwyt (`server/internal/core/handlers_terminal_wyjscie.go`), więc
 * podgląd w tle stoi na tej komendzie, a nie na obejściu.
 *
 * Komenda robi dwie czynności naraz, a odpowiedź rozdziela je na dwa pola:
 * zapisuje wskazane okno na zbiorcze wyjście wszystkich otwartych kart
 * terminala (`subscribed`) oraz oddaje ogon historii (`lines`). Żądanie bez
 * `windowId` jest kontraktem dopuszczone i znaczy sam odczyt ogona — wraca
 * wtedy `subscribed: false`. To nie jest awaria i okno ma to powiedzieć wprost,
 * zamiast milczeć albo udawać podgląd na żywo, którego nie ma.
 *
 * Nowe wiersze jadą `stream.chunk`, nie osobnym zdarzeniem: kontrakt nie ma
 * zdarzenia zbiorczego wyjścia, a rdzeń rozsyła wiersze obserwatorowi wspólnym
 * strumieniem fragmentów z `windowId` okna obserwującego i `messageId` równym
 * identyfikatorowi procesu. Dlatego druga czynność tego źródła jest
 * subskrypcją, a nie drugą komendą.
 *
 * To nie jest drugie źródło terminala: `moduly/terminal/zrodlo-terminala.ts`
 * niesie komendy kart i procesów, a ten plik wyłącznie tę jedną i żadnej
 * z tamtych nie powiela.
 */
export interface ZrodloPodgladuBash {
  /**
   * `terminal.output.stream` — zapis okna na zbiorcze wyjście kart wraz
   * z ogonem historii.
   */
  zapiszNaWyjscie(
    zadanie: TerminalOutputStreamRequest,
  ): Promise<Wynik<TerminalOutputStreamResponse>>;
  /** Subskrypcja `stream.chunk` — wiersze dochodzące po zapisie. */
  naFragmentWyjscia(
    sluchacz: (tresc: StreamChunkEvent, koperta: Envelope) => void,
  ): Odsubskrybuj;
}

export function utworzZrodloPodgladuBash(kanal: Kanal): ZrodloPodgladuBash {
  return {
    async zapiszNaWyjscie(zadanie) {
      // Sprawdzian kształtu obejmuje oba pola odpowiedzi, bo okno rozstrzyga
      // z nich dwie różne rzeczy: `lines` daje historię, `subscribed` daje
      // prawo do zdania „podgląd jest na żywo". Brak któregokolwiek zamienia
      // odpowiedź w zwykłe niepowodzenie wywołania, a nie w pustą historię.
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.TerminalOutputStream, zadanie),
        Command.TerminalOutputStream,
        (tresc) => czyTablica(tresc.lines) && czyLogiczna(tresc.subscribed),
      );
    },

    naFragmentWyjscia(sluchacz) {
      return kanal.naZdarzenie(EventType.StreamChunk, (tresc, koperta) => sluchacz(tresc, koperta));
    },
  };
}
