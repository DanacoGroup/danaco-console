import {
  Command,
  EventType,
  WindowRole,
  type ErrorInfo,
  type LoopState,
  type MonitorStatus,
  type MonitorStatusRequest,
  type ProgressChangedEvent,
  type Queue,
  type QueueChangedEvent,
  type QueueListRequest,
  type RoleListRequest,
  type Window,
  type WindowListRequest,
  type WindowRoleAssignment,
  type WindowStateChangedEvent,
  type WindowStateGetRequest,
  type WindowStateGetResponse,
} from '../../../shared/contract';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import type { Kanal, Wynik } from '../protokol/kanal';
import { czyObiekt, czyTablica, sprawdzKsztalt } from '../protokol/ksztalt-odpowiedzi';
import { przenies } from '../protokol/wynik-czastkowy';
import { wywolaj } from '../protokol/wywolanie';

/**
 * Źródło obrazu interwencji — jedyne miejsce warstwy mobilnej znające komendy
 * stanu pracy: `monitor.status`, `queue.list`, `window.list`,
 * `window.state.get`, `role.list`. Odmowa jednego odczytu nie gasi obrazu —
 * każdy wynik niesie odmowę osobno.
 */
export interface ZrodloInterwencji {
  /** `monitor.status` — telemetria procesów; ta sama, z której czyta warstwa mobilna. */
  procesy(zadanie: MonitorStatusRequest): Promise<Wynik<MonitorStatus[]>>;
  /** `queue.list` — kolejki silnika pętli koordynator–wykonawca. */
  kolejki(zadanie: QueueListRequest): Promise<Wynik<Queue[]>>;
  /** `window.list` — okna komunikacji wraz z rolą i trybem uprawnień. */
  okna(zadanie: WindowListRequest): Promise<Wynik<Window[]>>;
  /** `window.state.get` — stan okna wraz z biegiem naprawczym (`loop`). */
  stanOkna(zadanie: WindowStateGetRequest): Promise<Wynik<WindowStateGetResponse>>;
  /** `role.list` — więź koordynator–wykonawca. */
  role(zadanie: RoleListRequest): Promise<Wynik<WindowRoleAssignment[]>>;
  /** Jeden odczyt składający obraz z pięciu komend powyżej. */
  zbierzObraz(): Promise<ObrazInterwencji>;
  /** `progress.changed` — postęp procesu na żywo. */
  naPostep(sluchacz: (tresc: ProgressChangedEvent) => void): Odsubskrybuj;
  /** `queue.changed` — zmiana kolejki na żywo. */
  naKolejke(sluchacz: (tresc: QueueChangedEvent) => void): Odsubskrybuj;
  /** `window.state.changed` — zmiana stanu okna na żywo. */
  naStanOkna(sluchacz: (tresc: WindowStateChangedEvent) => void): Odsubskrybuj;
}

/** Bieg naprawczy jednego koordynatora wraz z odmową odczytu, gdy stan jego okna nie doszedł od rdzenia. */
export interface BiegKoordynatora {
  windowId: string;
  loop?: LoopState;
  odmowa?: ErrorInfo;
}

/**
 * Obraz interwencji — wszystko, na czym telefon opiera decyzję, w jednym
 * odczycie. Każde pole niesie własne rozstrzygnięcie wywołania, bo odmowa
 * jednej komendy nie unieważnia pozostałych.
 */
export interface ObrazInterwencji {
  procesy: Wynik<MonitorStatus[]>;
  kolejki: Wynik<Queue[]>;
  okna: Wynik<Window[]>;
  role: Wynik<WindowRoleAssignment[]>;
  /** Bieg naprawczy każdego okna koordynatora — źródło zdania „pętla stoi”. */
  biegi: readonly BiegKoordynatora[];
  /** Chwila złożenia obrazu w milisekundach epoki. */
  odczytaneO: number;
}

export function utworzZrodloInterwencji(kanal: Kanal): ZrodloInterwencji {
  const zrodlo: ZrodloInterwencji = {
    async procesy(zadanie) {
      const wynik = await wywolaj(kanal, Command.MonitorStatus, zadanie);
      return przenies(
        sprawdzKsztalt(wynik, Command.MonitorStatus, (t) => czyTablica(t.statuses)),
        (t) => t.statuses,
      );
    },

    async kolejki(zadanie) {
      const wynik = await wywolaj(kanal, Command.QueueList, zadanie);
      return przenies(
        sprawdzKsztalt(wynik, Command.QueueList, (t) => czyTablica(t.queues)),
        (t) => t.queues,
      );
    },

    async okna(zadanie) {
      const wynik = await wywolaj(kanal, Command.WindowList, zadanie);
      return przenies(
        sprawdzKsztalt(wynik, Command.WindowList, (t) => czyTablica(t.windows)),
        (t) => t.windows,
      );
    },

    async stanOkna(zadanie) {
      const wynik = await wywolaj(kanal, Command.WindowStateGet, zadanie);
      return sprawdzKsztalt(wynik, Command.WindowStateGet, (t) => czyObiekt(t.window));
    },

    async role(zadanie) {
      const wynik = await wywolaj(kanal, Command.RoleList, zadanie);
      return przenies(
        sprawdzKsztalt(wynik, Command.RoleList, (t) => czyTablica(t.assignments)),
        (t) => t.assignments,
      );
    },

    async zbierzObraz() {
      // Cztery odczyty naraz — telefon czeka raz, nie cztery razy pod rząd.
      const [procesy, kolejki, okna, role] = await Promise.all([
        zrodlo.procesy({}),
        zrodlo.kolejki({}),
        zrodlo.okna({}),
        zrodlo.role({}),
      ]);

      // O bieg naprawczy pytamy wyłącznie okna koordynatorów — `loop` jest puste dla pozostałych ról.
      const koordynatorzy = wskazKoordynatorow(okna, role);
      const biegi = await Promise.all(
        koordynatorzy.map(async (windowId): Promise<BiegKoordynatora> => {
          const wynik = await zrodlo.stanOkna({ windowId });
          if (!wynik.udany || wynik.wynik === undefined) {
            return { windowId, odmowa: wynik.blad };
          }
          return { windowId, loop: wynik.wynik.loop };
        }),
      );

      return { procesy, kolejki, okna, role, biegi, odczytaneO: Date.now() };
    },

    naPostep: (sluchacz) => kanal.naZdarzenie(EventType.ProgressChanged, sluchacz),
    naKolejke: (sluchacz) => kanal.naZdarzenie(EventType.QueueChanged, sluchacz),
    naStanOkna: (sluchacz) => kanal.naZdarzenie(EventType.WindowStateChanged, sluchacz),
  };

  return zrodlo;
}

/**
 * Okna koordynatorów wskazane z dwóch źródeł naraz: `window.windowRole`
 * i nadania z `role.list`. Dwa źródła, bo odmowa jednego z odczytów nie może
 * zabrać telefonowi całej wiedzy o tym, gdzie stoi pętla.
 */
export function wskazKoordynatorow(
  okna: Wynik<Window[]>,
  role: Wynik<WindowRoleAssignment[]>,
): string[] {
  const zbior = new Set<string>();
  for (const okno of okna.wynik ?? []) {
    if (okno.windowRole === WindowRole.Coordinator) zbior.add(okno.id);
    if (okno.coordinatorWindowId !== undefined) zbior.add(okno.coordinatorWindowId);
  }
  for (const nadanie of role.wynik ?? []) {
    if (nadanie.role === WindowRole.Coordinator) zbior.add(nadanie.windowId);
    if (nadanie.coordinatorWindowId !== undefined) zbior.add(nadanie.coordinatorWindowId);
  }
  return [...zbior];
}
