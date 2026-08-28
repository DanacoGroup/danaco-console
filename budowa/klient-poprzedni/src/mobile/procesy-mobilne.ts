import {
  Command,
  MobileProcessControl,
  ProgressStatus,
  type MobileProcess,
  type MobileProcessControlRequest,
  type MobileProcessListRequest,
} from '../../../shared/contract';
import type { Kanal, Wynik } from '../protokol/kanal';
import { czyObiekt, czyTablica, sprawdzKsztalt } from '../protokol/ksztalt-odpowiedzi';
import { przenies } from '../protokol/wynik-czastkowy';
import { wywolaj } from '../protokol/wywolanie';

// Procesy platformy widziane z telefonu — `mobile.process.list`, `mobile.process.control`.

/** Czynności sterowania procesem wraz z nazwą widoczną dla Operatora oraz wagą nieodwracalności czynności. */
export interface CzynnoscSterowania {
  kod: MobileProcessControl;
  /** Nazwa czynności widziana przez Operatora — pełna, bez skrótu. */
  nazwa: string;
  /** Czy czynność zatrzymuje pracę nieodwracalnie; wstrzymanie i wznowienie przerwaniem nie są. */
  nieodwracalna: boolean;
  /** Zdanie ostrzeżenia — mówi, co Operator traci, zanim to straci. */
  ostrzezenie?: string;
}

export const CZYNNOSCI_STEROWANIA: readonly CzynnoscSterowania[] = [
  { kod: MobileProcessControl.Pause, nazwa: 'Wstrzymaj', nieodwracalna: false },
  { kod: MobileProcessControl.Resume, nazwa: 'Wznów', nieodwracalna: false },
  {
    kod: MobileProcessControl.Stop,
    nazwa: 'Zatrzymaj',
    nieodwracalna: true,
    ostrzezenie:
      'Zatrzymanie kończy pracę procesu i nie da się go wznowić — wznowić można wyłącznie ' +
      'proces wstrzymany. Tego, co proces zdążył wykonać, zatrzymanie nie odwraca, a tego, ' +
      'czego nie zdążył, nie dokończy nikt.',
  },
  {
    kod: MobileProcessControl.Restart,
    nazwa: 'Uruchom ponownie',
    nieodwracalna: true,
    ostrzezenie:
      'Ponowne uruchomienie przerywa bieg trwający i zaczyna go od początku. Postęp bieżącego ' +
      'biegu przepada.',
  },
  { kod: MobileProcessControl.Approve, nazwa: 'Zatwierdź krok czekający na decyzję', nieodwracalna: false },
  { kod: MobileProcessControl.Modify, nazwa: 'Modyfikuj zlecenie w toku', nieodwracalna: false },
];

/** Stany procesu dostępne w filtrze wraz z nazwą wyświetlaną; pusty wybór filtra znaczy wszystkie stany. */
export const STANY_PROCESU: readonly { kod: ProgressStatus; nazwa: string }[] = [
  { kod: ProgressStatus.Pending, nazwa: 'Oczekuje' },
  { kod: ProgressStatus.Running, nazwa: 'W biegu' },
  { kod: ProgressStatus.Paused, nazwa: 'Wstrzymany' },
  { kod: ProgressStatus.Stopped, nazwa: 'Zatrzymany' },
  { kod: ProgressStatus.Done, nazwa: 'Zakończony' },
  { kod: ProgressStatus.Failed, nazwa: 'Zakończony błędem' },
];

export interface ZrodloProcesowMobilnych {
  /** `mobile.process.list` — procesy platformy widoczne z urządzenia. */
  wykaz(zadanie: MobileProcessListRequest): Promise<Wynik<MobileProcess[]>>;
  /** `mobile.process.control` — sterowanie jednym procesem. */
  steruj(zadanie: MobileProcessControlRequest): Promise<Wynik<MobileProcess>>;
}

export function utworzZrodloProcesowMobilnych(kanal: Kanal): ZrodloProcesowMobilnych {
  return {
    async wykaz(zadanie) {
      const wynik = await wywolaj(kanal, Command.MobileProcessList, zadanie);
      return przenies(
        sprawdzKsztalt(wynik, Command.MobileProcessList, (tresc) => czyTablica(tresc.processes)),
        (tresc) => tresc.processes,
      );
    },

    async steruj(zadanie) {
      const wynik = await wywolaj(kanal, Command.MobileProcessControl, zadanie);
      return przenies(
        sprawdzKsztalt(wynik, Command.MobileProcessControl, (tresc) => czyObiekt(tresc.process)),
        (tresc) => tresc.process,
      );
    },
  };
}

/** Nazwa stanu procesu pokazywana wprost Operatorowi; stan spoza znanego wykazu wraca swoim kodem własnym. */
export function nazwaStanu(kod: ProgressStatus): string {
  return STANY_PROCESU.find((pozycja) => pozycja.kod === kod)?.nazwa ?? kod;
}

/** Nazwa czynności sterowania pokazywana Operatorowi; czynność spoza wykazu wraca własnym kodem źródłowym. */
export function nazwaCzynnosci(kod: MobileProcessControl): string {
  return CZYNNOSCI_STEROWANIA.find((pozycja) => pozycja.kod === kod)?.nazwa ?? kod;
}

/**
 * Opis jednego procesu w jednym zdaniu.
 *
 * Stopień ukończenia nieoddany znaczy „rdzeń go nie liczy", a nie zero procent:
 * proces oczekujący i proces bez licznika postępu to dwie różne rzeczy.
 */
export function opisProcesu(proces: MobileProcess): string {
  const czesci: string[] = [`stan: ${nazwaStanu(proces.status)}`];
  if (proces.completion !== undefined) czesci.push(`ukończone ${proces.completion}%`);
  else czesci.push('stopnia ukończenia rdzeń nie podaje');
  if (proces.sessionId !== undefined) czesci.push(`karta sesji ${proces.sessionId}`);
  if (proces.windowId !== undefined) czesci.push(`okno ${proces.windowId}`);
  if (proces.startedAt !== undefined) {
    czesci.push(`uruchomiony ${new Date(proces.startedAt).toLocaleString('pl-PL')}`);
  }
  return czesci.join(', ');
}
