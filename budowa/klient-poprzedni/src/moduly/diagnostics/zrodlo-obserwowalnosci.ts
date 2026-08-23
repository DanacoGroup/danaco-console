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
 * Telemetria procesów widziana przez moduł Diagnostics — rodzina `monitor.*`
 * wraz ze zdarzeniem `progress.changed`.
 *
 * Rodzina jest jedynym źródłem obserwowalności, jakie kontrakt niesie poza
 * obszarem `diagnostics.*`. Opracowanie modułu wywodzi metryki wydajności
 * i kontrolę stanu samej platformy z encji `proces_sesji`, a `MonitorStatus`
 * jest dokładnie jej odwzorowaniem w kontrakcie („ta sama telemetria co
 * progress.changed”, opis struktury). Zakładki Metrics & Performance oraz
 * Health & Uptime jadą więc jednym odczytem: dwie perspektywy na ten sam
 * materiał, nie dwa niezależne pomiary, które rozejdą się po pierwszej zmianie.
 *
 * Źródło oddaje `WynikDiagnostyki`, nie samą treść, i nigdzie nie podstawia
 * pustego wykazu za nieudany odczyt — z tego samego powodu, dla którego robi
 * tak `zrodlo-diagnostics.ts`. Odmowa `monitor.status` bywa tu zjawiskiem
 * zwykłym: rdzeń bez wpiętego rejestru telemetrii odmawia głośno
 * (`core/handlers_monitor.go`, gałąź monitora niewpiętego), a wykaz pusty
 * ukryłby ten fakt pod zdaniem „nie ma procesów”.
 *
 * Podział na dwie komendy jest podziałem ról, nie wygodą: `monitor.status`
 * czyta stan bieżący i niczego nie zapisuje, `monitor.subscribe` dodatkowo
 * zapisuje okno na telemetrię. Pole `subscribed` odpowiedzi mówi, czy zapis
 * doszedł do skutku — żądanie bez `windowId` jest zwykłym odczytem
 * (`core/adapter_modul_monitor.go`, `ObserwujProcesy`), a moduł Diagnostics
 * montuje się bez okna. Okno pyta o to pole i mówi Operatorowi wprost, bo
 * „obserwacja założona” i „obserwacja niezałożona” to dwa różne zdania.
 */

/** Odpowiedź `monitor.subscribe` — stan bieżący wraz z losem samego zapisu. */
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
  /**
   * Subskrypcja `progress.changed` — zmiana stanu procesu na żywo.
   *
   * Zdarzenie dochodzi do wszystkich połączeń konta niezależnie od tego, czy
   * `monitor.subscribe` zapisał okno; zapis mówi rdzeniowi wyłącznie, które
   * okno których procesów pilnuje. Dlatego wykaz zmienia się na żywo także
   * wtedy, gdy `subscribed` niesie `false`.
   */
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
