import { Command, type Window } from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../../protokol/wywolanie';

/** Kod modułu w katalogu rdzenia — kolumna `modul.kod`. */
export const KOD_MODULU = 'apps';

/**
 * Skąd moduł bierze `windowId` wymagany przez wszystkie trzy komendy `apps.*`.
 *
 * Każde żądanie obszaru niesie pole `windowId` opisane w kontrakcie jako
 * „Okno modułu Apps”, więc identyfikator musi pochodzić z rdzenia. Bierze go
 * komenda `window.list` zawężona do sesji i do okien modułu Apps:
 * `Window.moduleId` jest tym, co rdzeń przypisał oknu przy `workspace.enter`.
 *
 * Brak okna nie jest błędem: sesja bywa jeszcze nieznana w chwili wejścia
 * w moduł; wtedy wykaz wraca pusty, a okna modułu pokazują stan pusty
 * nazywający brakujący warunek zamiast wysyłać żądanie bez identyfikatora.
 */
export interface ZrodloOknaModulu {
  /** Okna komunikacji sesji należące do modułu Apps. */
  okna(idSesji: string): Promise<Wynik<{ windows: Window[] }>>;
}

export function utworzZrodloOknaModulu(kanal: Kanal): ZrodloOknaModulu {
  return {
    async okna(idSesji) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.WindowList, { sessionId: idSesji }),
        Command.WindowList,
        (tresc) => czyTablica(tresc.windows),
      );
    },
  };
}

/**
 * Okno modułu Apps spośród okien sesji.
 *
 * Wybieramy pierwsze okno o module zgodnym z kodem modułu; gdy rdzeń nie
 * przypisał żadnego, oddajemy pusty łańcuch, a nie okno przypadkowe.
 */
export function wybierzOknoModulu(okna: readonly Window[]): string {
  const nasze = okna.find((okno) => okno.moduleId === KOD_MODULU);
  return nasze?.id ?? '';
}
