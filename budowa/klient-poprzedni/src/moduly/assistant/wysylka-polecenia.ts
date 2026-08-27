import {
  AssistantActionControl,
  AssistantActionStatus,
  type AssistantVoiceCommandResponse,
} from '../../../../shared/contract';
import type { FazaOkna } from '../../komponenty/faza-okna';
import { opisOdmowy } from '../../komponenty/odmowa';
import { BRAKI, zglosBrak } from './braki-kontraktu';
import { ODCZYTY } from './etykiety-assistant';
import type { StanAssistant } from './stan-assistant';

/**
 * Droga polecenia z Voice Console do rdzenia i z powrotem odpowiada wyłącznie za wydanie
 * polecenia i przerwanie zlecenia, które z niego powstało.
 */
export interface WysylkaPolecenia {
  wyslij(polecenie: TrescPolecenia): Promise<void>;
  przerwij(): Promise<void>;
  /** Zdejmuje zlecenie z toru, gdy rdzeń zamknął je albo anulował. */
  zsynchronizuj(): void;
}

/** Wartości ustawione w pasku promptu: transkrypcja polecenia, wybrany profil i nastawa syntezy mowy odpowiedzi. */
export interface TrescPolecenia {
  transkrypcja: string;
  profil: string;
  czytaj: boolean;
}

/**
 * Sposób, w jaki okno pokazuje przebieg wysyłki, nazwaną wysyłką, nie oknem, tym samym
 * słownikiem fazy co inne okna.
 */
export interface WidokWysylki {
  /** Zdanie o tym, co się właśnie dzieje albo co odpowiedział rdzeń. */
  powiedz(tresc: string, powodzenie: boolean): void;
  /** Nanosi fazę okna wraz z jej zdaniem; kontrolki zostają klikalne. */
  faza(faza: FazaOkna, opis: string): void;
  /** Nanosi odpowiedź rdzenia; `null` ją zdejmuje. */
  odpowiedz(tresc: AssistantVoiceCommandResponse | null): void;
}

export function utworzWysylkePolecenia(
  stan: StanAssistant,
  widok: WidokWysylki,
): WysylkaPolecenia {
  let idPolecenia = '';

  return {
    async wyslij(polecenie) {
      // Pustego pola nie meldujemy jako błędu okna: to brak w polu formularza, nie widoku.
      if (polecenie.transkrypcja === '') {
        widok.powiedz('Wpisz albo podyktuj treść polecenia — rdzeń odmówi pustego.', false);
        return;
      }
      // Brak przydzielonego okna modułu jest błędem okna: polecenia nie ma dokąd wysłać.
      if (stan.idOkna() === '') {
        widok.powiedz(BRAKI.brakOkna, false);
        widok.faza('blad', BRAKI.brakOkna);
        return;
      }
      widok.faza('ladowanie', ODCZYTY.polecenie);
      widok.powiedz(ODCZYTY.polecenie, true);
      widok.odpowiedz(null);
      const wynik = await stan.zrodlo.polecenie({
        idOkna: stan.idOkna(),
        transkrypcja: polecenie.transkrypcja,
        profil: polecenie.profil,
        czytaj: polecenie.czytaj,
      });
      if (!wynik.udany || wynik.wynik === undefined) {
        const powod = opisOdmowy('Wydanie polecenia', wynik.blad?.code, wynik.blad?.message);
        widok.powiedz(powod, false);
        // Komunikat odmowy zostaje w układzie, a treść paska pozostaje widoczna pod nim.
        widok.faza('blad', powod);
        return;
      }
      idPolecenia = wynik.wynik.action.id;
      stan.wchlon([wynik.wynik.action]);
      widok.odpowiedz(wynik.wynik);
      widok.faza('gotowe', '');
      widok.powiedz('Zlecenie założone. Przebieg śledzi Actions Monitor.', true);
      await stan.odswiezDziennik();
    },

    async przerwij() {
      if (idPolecenia === '') {
        zglosBrak(
          'Przerwanie polecenia',
          'Nie ma czego przerywać: żadne polecenie nie zostało jeszcze wysłane z tego ' +
            'okna. Nagranie w toku zatrzymasz tym samym przyciskiem mikrofonu, którym ' +
            'je zacząłeś.',
        );
        return;
      }
      widok.faza('ladowanie', ODCZYTY.przerwanie);
      widok.powiedz('Anuluję zlecenie ostatniego polecenia…', true);
      const wynik = await stan.zrodlo.steruj({
        idZlecenia: idPolecenia,
        sterowanie: AssistantActionControl.Cancel,
      });
      if (!wynik.udany || wynik.wynik === undefined) {
        const powod = opisOdmowy('Przerwanie', wynik.blad?.code, wynik.blad?.message);
        widok.powiedz(powod, false);
        widok.faza('blad', powod);
        return;
      }
      stan.wchlon(wynik.wynik.actions);
      widok.faza('gotowe', '');
      widok.powiedz('Rdzeń przyjął przerwanie zlecenia.', true);
    },

    zsynchronizuj() {
      const biezace = stan.zlecenia().find((zlecenie) => zlecenie.id === idPolecenia);
      if (biezace === undefined) return;
      if (
        biezace.status === AssistantActionStatus.Running ||
        biezace.status === AssistantActionStatus.Paused ||
        biezace.status === AssistantActionStatus.Queued
      ) {
        return;
      }
      // Zlecenie zeszło z toru: nie ma już czego anulować, a przycisk powie o tym wprost.
      idPolecenia = '';
    },
  };
}
