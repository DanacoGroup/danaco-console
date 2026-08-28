import { Command, QueueAction } from '../../../shared/contract';
import type { IdSrodowiska, RodzajUtworzenia } from './model-danych';

/**
 * Ładunki zdarzeń, które pulpit wystawia powłoce, opisują zamiar Operatora
 * nazwami kontraktu. Wejście do sesji z matrycy niesie komendę `session.open`
 * wraz z identyfikatorem sesji, środowiskiem kolumny oraz tytułem sesji.
 */
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
 * Zamiar wobec kolejki: cztery przyciski transportu odpowiadają wprost
 * wyliczeniu `QueueAction` kontraktu, piąty — przekazanie — komendzie
 * `context.transfer`, a podniesienie priorytetu odpowiednika w kontrakcie
 * nie ma.
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
      /** Brak odpowiednika w kontrakcie, więc komenda pozostaje pusta. */
      komenda: null;
      idKolejki: string;
      rola: string;
    };

/**
 * Zamiar utworzenia bytu z rzędu kafli „Utwórz" niesie rodzaj tworzonego bytu
 * oraz nazwę kafla, którą nacisnął Operator; komendę dobiera z nich powłoka,
 * ponieważ pulpit komendy tworzenia sam nie wysyła.
 */
export interface ZamiarUtworzenia {
  rodzaj: RodzajUtworzenia;
  /** Nazwa kafla, którą nacisnął operator. */
  etykieta: string;
}

/**
 * Wezwanie do rozstrzygnięcia wstrzymanych przepływów niesie liczbę przepływów
 * czekających w chwili naciśnięcia oraz przepływ czekający najdłużej, więc
 * powłoka wie, ile spraw czeka i od której zacząć.
 */
export interface ZamiarDecyzji {
  /** Ile przepływów czeka w chwili naciśnięcia. */
  przeplywy: number;
  /** Przepływ czekający najdłużej. */
  najstarszy: string;
}

/**
 * Wykaz działań transportu kolejki w kolejności widocznej na pasku; kolejność
 * stoi wyłącznie tutaj, więc pasek przycisków i zamiar wysyłany do powłoki
 * nie mogą się co do niej rozminąć.
 */
export const DZIALANIA_TRANSPORTU: readonly QueueAction[] = [
  QueueAction.Start,
  QueueAction.Pause,
  QueueAction.Stop,
  QueueAction.Retry,
];
