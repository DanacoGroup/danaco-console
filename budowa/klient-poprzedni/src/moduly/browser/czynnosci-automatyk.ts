/**
 * Rozmowa okna Automation Studio z rdzeniem: odczyt automatyk, zapis definicji,
 * harmonogram, przebiegi i przekazanie scenariusza do modułu Automations.
 */

import type {
  AutomationExecution,
  AutomationSchedule,
  AutomationStep,
  AutomationWorkflow,
} from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import { skutekHarmonogramu, skutekPrzekazania, skutekZapisuAutomatyki } from './skutek-zapisu';
import type { StanPrzegladania } from './stan-przegladania';

/** Scenariusz przeglądania w postaci, w której idzie do rdzenia — zestaw kroków wraz z nazwą, opisem i stanem aktywności. */
export interface DefinicjaAutomatyki {
  /** Automatyka zmieniana; pusty napis zakłada nową. */
  identyfikator: string;
  nazwa: string;
  opis: string;
  czynna: boolean;
  kroki: readonly AutomationStep[];
}

export interface CzynnosciAutomatyk {
  /** `automation.workflow.list`; `null` znaczy odmowę odczytu. */
  odczytaj(): Promise<readonly AutomationWorkflow[] | null>;
  /** `automation.workflow.save`; `null` znaczy odmowę zapisu. */
  zapisz(definicja: DefinicjaAutomatyki): Promise<AutomationWorkflow | null>;
  /** `automation.schedule.set` — cykliczność w notacji cron. */
  ustawHarmonogram(idAutomatyki: string, cron: string, obowiazuje: boolean): Promise<void>;
  /** `schedule.get` — harmonogramy jednej automatyki. */
  odczytajHarmonogramy(idAutomatyki: string): Promise<readonly AutomationSchedule[] | null>;
  /** `automation.execution.subscribe` — stan przebiegów wraz z obserwacją okna. */
  odczytajPrzebiegi(idAutomatyki: string): Promise<readonly AutomationExecution[] | null>;
  /** `context.transfer` definicji do modułu Automations. */
  przekaz(definicja: DefinicjaAutomatyki): Promise<void>;
}

export function utworzCzynnosciAutomatyk(
  stan: StanPrzegladania,
  powiedz: (tresc: string, powodzenie: boolean) => void,
): CzynnosciAutomatyk {
  return {
    async odczytaj() {
      powiedz('Odczyt automatyk z rdzenia…', true);
      const wynik = await stan.automatyki.wykaz({});
      if (!wynik.udany || wynik.wynik === undefined) {
        powiedz(opisOdmowy('Odczyt automatyk', wynik.blad?.code, wynik.blad?.message), false);
        return null;
      }
      const automatyki = wynik.wynik.workflows;
      powiedz(`Wykaz zaciągnięty: ${automatyki.length} automatyk.`, true);
      return automatyki;
    },

    async zapisz(definicja) {
      const nazwa = definicja.nazwa.trim();
      if (nazwa === '') {
        powiedz('Nazwa scenariusza jest wymagana — rdzeń odmówi zapisu bez niej.', false);
        return null;
      }
      const opis = definicja.opis.trim();
      powiedz(`Zapis automatyki „${nazwa}"…`, true);
      const wynik = await stan.automatyki.zapisz({
        name: nazwa,
        ...(definicja.identyfikator === '' ? {} : { workflowId: definicja.identyfikator }),
        ...(opis === '' ? {} : { description: opis }),
        steps: [...definicja.kroki],
        enabled: definicja.czynna,
      });
      if (!wynik.udany || wynik.wynik === undefined) {
        powiedz(opisOdmowy('Zapis automatyki', wynik.blad?.code, wynik.blad?.message), false);
        return null;
      }
      const skutek = skutekZapisuAutomatyki(wynik.wynik.workflow, {
        nazwa,
        krokow: definicja.kroki.length,
      });
      powiedz(skutek.zdanie, skutek.udany);
      // Definicja z rdzenia wraca mimo rozjazdu — wiersz pokazuje jej postać obok zdania o rozbieżności.
      return wynik.wynik.workflow;
    },

    async ustawHarmonogram(idAutomatyki, cron, obowiazuje) {
      if (idAutomatyki === '') {
        powiedz(
          'Harmonogram wiąże się z automatyką zapisaną — najpierw zapisz scenariusz w rdzeniu.',
          false,
        );
        return;
      }
      const cyklicznosc = cron.trim();
      if (cyklicznosc === '') {
        powiedz('Podaj cykliczność w notacji cron — bez niej harmonogram nie ma treści.', false);
        return;
      }
      powiedz('Zapis harmonogramu automatyki…', true);
      const wynik = await stan.automatyki.ustawHarmonogram({
        workflowId: idAutomatyki,
        cron: cyklicznosc,
        enabled: obowiazuje,
      });
      if (!wynik.udany || wynik.wynik === undefined) {
        powiedz(opisOdmowy('Zapis harmonogramu', wynik.blad?.code, wynik.blad?.message), false);
        return;
      }
      const skutek = skutekHarmonogramu(wynik.wynik.schedule, cyklicznosc);
      powiedz(skutek.zdanie, skutek.udany);
    },

    async odczytajHarmonogramy(idAutomatyki) {
      if (idAutomatyki === '') {
        powiedz('Wskaż automatykę — harmonogram czyta się dla jednej definicji.', false);
        return null;
      }
      powiedz('Odczyt harmonogramów automatyki…', true);
      const wynik = await stan.automatyki.harmonogramy({ workflowId: idAutomatyki });
      if (!wynik.udany || wynik.wynik === undefined) {
        powiedz(opisOdmowy('Odczyt harmonogramów', wynik.blad?.code, wynik.blad?.message), false);
        return null;
      }
      const harmonogramy = wynik.wynik.schedules;
      powiedz(
        harmonogramy.length === 0
          ? 'Automatyka nie ma harmonogramu — uruchamia się wyłącznie na żądanie.'
          : `Harmonogramów automatyki: ${harmonogramy.length}.`,
        true,
      );
      return harmonogramy;
    },

    async odczytajPrzebiegi(idAutomatyki) {
      if (idAutomatyki === '') {
        powiedz('Wskaż automatykę — przebiegi czyta się dla jednej definicji.', false);
        return null;
      }
      const idOkna = stan.idOkna();
      powiedz('Odczyt przebiegów automatyki…', true);
      // Okno zgłasza się odbiorcą telemetrii, gdy rdzeń je wskazał — inaczej trafi na nieistniejące okno.
      const wynik = await stan.automatyki.przebiegi({
        workflowId: idAutomatyki,
        ...(idOkna === '' ? {} : { windowId: idOkna }),
      });
      if (!wynik.udany || wynik.wynik === undefined) {
        powiedz(opisOdmowy('Odczyt przebiegów', wynik.blad?.code, wynik.blad?.message), false);
        return null;
      }
      const przebiegi = wynik.wynik.executions;
      powiedz(
        przebiegi.length === 0
          ? 'Automatyka nie ma jeszcze ani jednego przebiegu.'
          : `Przebiegów: ${przebiegi.length}; obserwacja okna ${wynik.wynik.subscribed ? 'założona' : 'NIE założona'}.`,
        true,
      );
      return przebiegi;
    },

    async przekaz(definicja) {
      const idOkna = stan.idOkna();
      if (idOkna === '') {
        powiedz(stan.powod(), false);
        return;
      }
      const nazwa = definicja.nazwa.trim();
      if (nazwa === '' || definicja.kroki.length === 0) {
        powiedz(
          'Przekazanie idzie z gotowym scenariuszem — nadaj mu nazwę i dołóż co najmniej jeden krok.',
          false,
        );
        return;
      }
      powiedz('Przekazanie scenariusza do modułu Automations…', true);
      const wynik = await stan.zapisy.przekaz(idOkna, 'automations', {
        prompt: `Scenariusz przeglądania „${nazwa}" z modułu Browser`,
        executionParams: {
          name: nazwa,
          description: definicja.opis.trim(),
          enabled: definicja.czynna,
          steps: [...definicja.kroki],
        },
      });
      if (!wynik.udany || wynik.wynik === undefined) {
        powiedz(
          opisOdmowy('Przekazanie do Automations', wynik.blad?.code, wynik.blad?.message),
          false,
        );
        return;
      }
      const skutek = skutekPrzekazania(
        wynik.wynik,
        'automations',
        `Wysłano scenariusz „${nazwa}" o ${definicja.kroki.length} krokach`,
      );
      powiedz(skutek.zdanie, skutek.udany);
    },
  };
}
