import { Command, type Window } from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../../protokol/wywolanie';

/**
 * Kod modułu w katalogu rdzenia, odpowiadający kolumnie `modul.kod`. Rdzeń
 * przypisuje ten kod oknu przy wejściu w obszar roboczy, więc wybór okna
 * modułu opiera się na porównaniu z tą stałą.
 */
export const KOD_MODULU = 'apps';

/**
 * Źródło identyfikatora okna, którego żądają wszystkie komendy `apps.*`.
 * Identyfikator pochodzi z rdzenia: podaje go komenda `window.list` zawężona
 * do sesji i do okien modułu Apps. Brak okna nie jest błędem, bo sesja bywa
 * jeszcze nieznana.
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
