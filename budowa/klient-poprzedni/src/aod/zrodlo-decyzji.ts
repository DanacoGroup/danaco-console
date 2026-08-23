import {
  Command,
  ConfigScope,
  QueueAction,
  WindowRole,
  type ConfigWindowOpenResponse,
  type MessageStopResponse,
  type MonitorStatusResponse,
  type QueueActionResponse,
  type QueueListResponse,
  type RequestOf,
  type ResponseOf,
  type RoleUpdateResponse,
  type SessionFocusRequest,
  type SessionFocusResponse,
  type SessionStopResponse,
} from '../../../shared/contract';
import type { Kanal, Wynik } from '../protokol/kanal';

/**
 * Źródło komend kolejki decyzji — jedyne miejsce w katalogu `aod/`, które zna
 * nazwy komend spoza rodziny `aod.*`. Stoi osobno od `zrodlo-komend.ts`, bo
 * tamto niesie wyłącznie rodzinę nakładki; sklejenie obu zatarłoby tę granicę.
 *
 * Wpisana tu jest tylko komenda, którą rdzeń obsługuje. Źródło komend
 * deklarujące komendę nieobsługiwaną produkuje przyciski donikąd.
 */

/** Opakowuje `kanal.wyslij` w Promise — tak samo jak `zrodlo-komend.ts`. */
function poslijKomende<K extends Command>(
  kanal: Kanal,
  komenda: K,
  zadanie: RequestOf<K>,
): Promise<Wynik<ResponseOf<K>>> {
  return new Promise((rozstrzygnij) => {
    kanal.wyslij(komenda, zadanie, (wynik) => rozstrzygnij(wynik));
  });
}

/** Czynności, po które kolejka decyzji sięga poza rodzinę `aod.*`. */
export interface ZrodloDecyzji {
  /** `monitor.status` bez argumentów: stan wszystkich procesów, które rdzeń zna. */
  procesy(): Promise<Wynik<MonitorStatusResponse>>;
  /** `queue.list` bez argumentów — kolejki, przez które decyzja dochodzi do kroku. */
  kolejki(): Promise<Wynik<QueueListResponse>>;
  /** `queue.action{action:'pause'}` — wstrzymanie kolejki niosącej decyzję. */
  wstrzymajKolejke(queueId: string): Promise<Wynik<QueueActionResponse>>;
  /** `message.stop` — zatrzymanie tury biegnącej w oknie decyzji. */
  zatrzymajTureOkna(windowId: string): Promise<Wynik<MessageStopResponse>>;
  /** `session.stop` — zatrzymanie tur we wszystkich oknach sesji. */
  zatrzymajTurySesji(sessionId: string): Promise<Wynik<SessionStopResponse>>;
  /** `config.window.open` na zasięgu okna — konfiguracja Koordynatora. */
  otworzKonfiguracjeOkna(windowId: string): Promise<Wynik<ConfigWindowOpenResponse>>;
  /** `role.update{role:'standalone'}` — wyjęcie okna z pętli koordynatora. */
  wyjmijZPetli(windowId: string): Promise<Wynik<RoleUpdateResponse>>;
  /** `session.focus` — przeniesienie ogniska na okno decyzji. */
  przeniesOgnisko(zadanie: SessionFocusRequest): Promise<Wynik<SessionFocusResponse>>;
}

export function utworzZrodloDecyzji(kanal: Kanal): ZrodloDecyzji {
  return {
    // Bez argumentów świadomie: nakładka pyta o wszystkie procesy, które rdzeń
    // zna, bo proces czekający na decyzję nie musi należeć do sesji tego klienta.
    procesy: () => poslijKomende(kanal, Command.MonitorStatus, {}),

    kolejki: () => poslijKomende(kanal, Command.QueueList, {}),

    wstrzymajKolejke: (queueId) =>
      poslijKomende(kanal, Command.QueueAction, { queueId, action: QueueAction.Pause }),

    zatrzymajTureOkna: (windowId) => poslijKomende(kanal, Command.MessageStop, { windowId }),

    zatrzymajTurySesji: (sessionId) => poslijKomende(kanal, Command.SessionStop, { sessionId }),

    // Zasięg `window` jest poziomem najwęższym i wygrywa z każdym szerszym —
    // konfiguracja Koordynatora dotyczy tego jednego okna.
    otworzKonfiguracjeOkna: (windowId) =>
      poslijKomende(kanal, Command.ConfigWindowOpen, { windowId, scope: ConfigScope.Window }),

    wyjmijZPetli: (windowId) =>
      poslijKomende(kanal, Command.RoleUpdate, { windowId, role: WindowRole.Standalone }),

    przeniesOgnisko: (zadanie) => poslijKomende(kanal, Command.SessionFocus, zadanie),
  };
}
