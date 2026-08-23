import {
  Command,
  EventType,
  type Channel,
  type ContextBundle,
  type ContextTransferResponse,
  type Module,
  type ProgressChangedEvent,
  type Window,
  type WindowStateGetResponse,
} from '../../../../shared/contract';
import type { Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyLiczba, czyObiekt, czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../../protokol/wywolanie';

/**
 * Zaplecze modułu Design — wywołania czterech cudzych obszarów kontraktu,
 * z których moduł korzysta, a których nie prowadzi.
 *
 *   `window.list` / `window.state.get`
 *        — okno modułu i jego stan; komenda wspólna każdemu oknu operacyjnemu.
 *   `channel.list`
 *        — rejestr kanałów modelu. Najbliższy istniejący odpowiednik „wyboru
 *          silnika generującego" z panelu akcji Prompt Buildera; pole `engine`
 *          promptu przyjmuje identyfikator kanału.
 *   `module.list` / `context.transfer`
 *        — katalog modułów rdzenia i jedyna zbudowana droga przekazania
 *          międzymodułowego. Nią idzie „Wyślij do modułu docelowego"
 *          z Assets Panel. Wykaz modułów docelowych bierze się z rdzenia,
 *          nigdy z kopii katalogu po stronie klienta.
 *   `progress.changed`
 *        — telemetria postępu. Generowanie zasobu nim nie jedzie:
 *          `design.asset.generate` wraca w jednej turze, bez pola `processId`,
 *          i nie wysyła po drodze żadnego `progress.changed`.
 */
export interface ZrodloZaplecza {
  okna(idSesji: string): Promise<Wynik<{ windows: Window[] }>>;
  stanOkna(idOkna: string): Promise<Wynik<WindowStateGetResponse>>;
  silniki(): Promise<Wynik<{ channels: Channel[] }>>;
  modulyDocelowe(): Promise<Wynik<{ modules: Module[] }>>;
  przekaz(zlecenie: ZleceniePrzekazania): Promise<Wynik<ContextTransferResponse>>;
  /** Subskrypcja `progress.changed` — postęp procesu wskazanego okna. */
  naPostep(sluchacz: (tresc: ProgressChangedEvent) => void): Odsubskrybuj;
}

/** Zlecenie przekazania kompletu kontekstu do modułu docelowego. */
export interface ZleceniePrzekazania {
  idOknaZrodlowego: string;
  kodModuluDocelowego: string;
  komplet: ContextBundle;
}

export function utworzZrodloZaplecza(kanal: Kanal): ZrodloZaplecza {
  return {
    async okna(idSesji) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.WindowList, { sessionId: idSesji }),
        Command.WindowList,
        (tresc) => czyTablica(tresc.windows),
      );
    },

    async stanOkna(idOkna) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.WindowStateGet, { windowId: idOkna }),
        Command.WindowStateGet,
        (tresc) => czyObiekt(tresc.window) && czyLiczba(tresc.messageCount),
      );
    },

    async silniki() {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.ChannelList, { enabledOnly: true }),
        Command.ChannelList,
        (tresc) => czyTablica(tresc.channels),
      );
    },

    async modulyDocelowe() {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.ModuleList, {}),
        Command.ModuleList,
        (tresc) => czyTablica(tresc.modules),
      );
    },

    async przekaz(zlecenie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.ContextTransfer, {
          sourceWindowId: zlecenie.idOknaZrodlowego,
          targetModuleId: zlecenie.kodModuluDocelowego,
          bundle: zlecenie.komplet,
        }),
        Command.ContextTransfer,
        (tresc) => czyObiekt(tresc.window) && typeof tresc.transferred === 'boolean',
      );
    },

    naPostep(sluchacz) {
      return kanal.naZdarzenie(EventType.ProgressChanged, (tresc) => sluchacz(tresc));
    },
  };
}
