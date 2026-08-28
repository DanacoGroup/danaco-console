/**
 * Telemetria procesów widziana przez moduł Diagnostics: rodzina komend
 * `monitor.*` wraz ze zdarzeniem `progress.changed`. Źródło oddaje wynik
 * diagnostyki, a nie samą treść, i nie podstawia pustego wykazu za nieudany
 * odczyt.
 */
import {
  Command,
  EventType,
  type MonitorStatus,
  type MonitorStatusRequest,
  type MonitorSubscribeRequest,
  type ProgressChangedEvent,
} from '../../../../shared/contract';
import type { Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import type { Kanal } from '../../protokol/kanal';
import { czyLogiczna, czyTablica } from '../../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../../protokol/wywolanie';
import { rozstrzygnij, type WynikDiagnostyki } from './zrodlo-diagnostics';

/**
 * Odpowiedź komendy `monitor.subscribe`: stan bieżący procesów wraz z losem
 * samego zapisu okna na telemetrię. Pole zapisu mówi wprost, czy obserwacja
 * została założona.
 */
export interface ObserwacjaProcesow {
  statuses: MonitorStatus[];
  /** Czy rdzeń zapisał okno na telemetrię; `false` przy żądaniu bez okna. */
  subscribed: boolean;
}

export interface ZrodloObserwowalnosci {
  /** `monitor.status` — stan bieżący procesów bez zakładania obserwacji. */
  stanProcesow(zadanie: MonitorStatusRequest): Promise<WynikDiagnostyki<MonitorStatus[]>>;
  /** `monitor.subscribe` — zapis okna na telemetrię wraz ze stanem bieżącym. */
  obserwujProcesy(
    zadanie: MonitorSubscribeRequest,
  ): Promise<WynikDiagnostyki<ObserwacjaProcesow>>;
  /** Subskrypcja `progress.changed`, czyli zmiana stanu procesu na żywo. */
  naPostep(sluchacz: (tresc: ProgressChangedEvent) => void): Odsubskrybuj;
}

export function utworzZrodloObserwowalnosci(kanal: Kanal): ZrodloObserwowalnosci {
  return {
    async stanProcesow(zadanie) {
      return rozstrzygnij(
        await wywolaj(kanal, Command.MonitorStatus, zadanie),
        Command.MonitorStatus,
        (tresc) => czyTablica(tresc.statuses),
        (tresc) => tresc.statuses,
      );
    },

    async obserwujProcesy(zadanie) {
      return rozstrzygnij(
        await wywolaj(kanal, Command.MonitorSubscribe, zadanie),
        Command.MonitorSubscribe,
        (tresc) => czyTablica(tresc.statuses) && czyLogiczna(tresc.subscribed),
        (tresc) => ({ statuses: tresc.statuses, subscribed: tresc.subscribed }),
      );
    },

    naPostep(sluchacz) {
      return kanal.naZdarzenie(EventType.ProgressChanged, (tresc) => sluchacz(tresc));
    },
  };
}
