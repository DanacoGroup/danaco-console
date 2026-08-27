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
 * Źródło podglądu w tle stoi na jednej komendzie kontraktu i jednej subskrypcji: komenda zapisuje okno na zbiorcze wyjście i oddaje ogon historii, a nowe wiersze dochodzą tą samą subskrypcją, nie osobnym zdarzeniem.
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
      // Sprawdzian obejmuje oba pola odpowiedzi: historię oraz prawo do zdania, że podgląd jest na żywo.
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
