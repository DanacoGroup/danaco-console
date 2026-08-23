import { Command, QueueAction } from '../../../shared/contract';
import type { IdSrodowiska, RodzajUtworzenia } from './model-danych';

/**
 * Ładunki zdarzeń, które pulpit wystawia powłoce.
 *
 * Jedna odpowiedzialność: opis zamiaru operatora wyrażony nazwami kontraktu.
 * Pulpit nie wysyła komend samodzielnie — nadaje zamiar, a powłoka składa
 * z niego kopertę. Dzięki temu nazwa komendy pojawia się w jednym miejscu i jest
 * importowana z `shared/contract`, nie przepisana.
 */

/** Wejście do sesji z matrycy — powłoka wysyła `session.open`. */
export interface WejscieDoSesji {
  /** Komenda kontraktu, którą zamiar realizuje. */
  komenda: typeof Command.SessionOpen;
  /** Identyfikator sesji — pole `sessionId` żądania. */
  sessionId: string;
  /** Środowisko kolumny, z której operator wszedł. */
  srodowisko: IdSrodowiska;
  /** Nazwa sesji, do wypisania w powłoce. */
  tytul: string;
}

/**
 * Zamiar wobec kolejki.
 *
 * Cztery przyciski transportu odpowiadają wprost `QueueAction` kontraktu,
 * piąty — przekazanie — komendzie `context.transfer`.
 *
 * Pole `rola` niesie nazwę kolejki widoczną dla operatora; kontrakt nie niesie
 * roli kolejki, więc zamiar jej nie podstawia.
 *
 * „Podnieś priorytet" nie ma odpowiednika ani w `QueueAction`, ani wśród komend,
 * dlatego zamiar niesie rodzaj `priorytet` z komendą `null` zamiast literału
 * wymyślonego po stronie klienta.
 */
export type ZamiarKolejki =
  | {
      rodzaj: 'kolejka';
      /** Komenda kontraktu, którą zamiar realizuje. */
      komenda: typeof Command.QueueAction;
      /** Działanie transportu z wyliczenia kontraktu. */
      dzialanie: QueueAction;
      idKolejki: string;
      rola: string;
    }
  | {
      rodzaj: 'przekazanie';
      komenda: typeof Command.ContextTransfer;
      idKolejki: string;
      rola: string;
    }
  | {
      rodzaj: 'priorytet';
      /** Brak odpowiednika w kontrakcie — patrz uwaga wyżej. */
      komenda: null;
      idKolejki: string;
      rola: string;
    };

/** Zamiar utworzenia bytu z rzędu „Utwórz". */
export interface ZamiarUtworzenia {
  rodzaj: RodzajUtworzenia;
  /** Nazwa kafla, którą nacisnął operator. */
  etykieta: string;
}

/** Wezwanie do rozstrzygnięcia wstrzymanych przepływów. */
export interface ZamiarDecyzji {
  /** Ile przepływów czeka w chwili naciśnięcia. */
  przeplywy: number;
  /** Przepływ czekający najdłużej. */
  najstarszy: string;
}

/** Wykaz działań transportu kolejki w kolejności widocznej na pasku. */
export const DZIALANIA_TRANSPORTU: readonly QueueAction[] = [
  QueueAction.Start,
  QueueAction.Pause,
  QueueAction.Stop,
  QueueAction.Retry,
];
