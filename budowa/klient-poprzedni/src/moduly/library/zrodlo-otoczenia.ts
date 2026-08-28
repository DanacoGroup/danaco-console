import {
  Command,
  ConfigScope,
  EventType,
  type Action,
  type ContextBundle,
  type Module,
  type Window,
  type WindowChangedEvent,
} from '../../../../shared/contract';
import type { Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyObiekt, czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import type { StrazOdmow } from './straz-odmow';

/**
 * Otoczenie modułu Library udostępnia komendy, z których moduł korzysta, a
 * których nie prowadzi: okna komunikacji, kontekst okna, katalog i wykonanie
 * akcji, przeniesienie kontekstu, katalog modułów oraz odczyt i zdarzenie
 * przekazanego kompletu.
 */
export interface ZrodloOtoczenia {
  okna(idSesji: string): Promise<Wynik<{ windows: Window[] }>>;
  stanOkna(idOkna: string): Promise<Wynik<{ window: Window; messageCount: number }>>;
  akcje(): Promise<Wynik<{ actions: Action[] }>>;
  wykonajAkcje(idOkna: string, idAkcji: string): Promise<Wynik<{ actionId: string }>>;
  // Przenosi komplet kontekstu; jedynym świadkiem, dokąd trafił, jest okno z odpowiedzi, nie wpisany.
  przeniesKontekst(
    idOkna: string,
    modulDocelowy: string,
    komplet: ContextBundle,
  ): Promise<Wynik<{ transferred: boolean; window?: Window }>>;
  /** Katalog modułów platformy — pozycje stera modułu docelowego. */
  moduly(): Promise<Wynik<{ modules: Module[] }>>;
  /** Komplet kontekstu okna — treść tego, co przyniosło przekazanie. */
  komplet(idOkna: string): Promise<Wynik<{ context: ContextBundle }>>;
  /** Subskrypcja `window.changed` — jedyne zdarzenie przekazania kontekstu. */
  naZmianeOkna(sluchacz: (tresc: WindowChangedEvent) => void): Odsubskrybuj;
}

/** Kod modułu z kolumny katalogu modułów rdzenia, wyznaczający zasięg katalogu akcji tego panelu okna, wprost. */
export const KOD_MODULU = 'library';

export function utworzZrodloOtoczenia(kanal: Kanal, straz: StrazOdmow): ZrodloOtoczenia {
  return {
    async okna(idSesji) {
      return sprawdzKsztalt(
        await straz.wywolaj(Command.WindowList, { sessionId: idSesji }),
        Command.WindowList,
        (tresc) => czyTablica(tresc.windows),
      );
    },

    async stanOkna(idOkna) {
      return sprawdzKsztalt(
        await straz.wywolaj(Command.WindowStateGet, { windowId: idOkna }),
        Command.WindowStateGet,
        (tresc) => czyObiekt(tresc.window),
      );
    },

    async akcje() {
      return sprawdzKsztalt(
        await straz.wywolaj(Command.ActionList, {
          scope: ConfigScope.Module,
          scopeId: KOD_MODULU,
          enabledOnly: true,
        }),
        Command.ActionList,
        (tresc) => czyTablica(tresc.actions),
      );
    },

    async wykonajAkcje(idOkna, idAkcji) {
      return sprawdzKsztalt(
        await straz.wywolaj(Command.WindowAction, { windowId: idOkna, actionId: idAkcji }),
        Command.WindowAction,
        (tresc) => typeof tresc.actionId === 'string',
      );
    },

    async przeniesKontekst(idOkna, modulDocelowy, komplet) {
      return sprawdzKsztalt(
        await straz.wywolaj(Command.ContextTransfer, {
          sourceWindowId: idOkna,
          targetModuleId: modulDocelowy,
          bundle: komplet,
        }),
        Command.ContextTransfer,
        (tresc) => typeof tresc.transferred === 'boolean',
      );
    },

    async moduly() {
      return sprawdzKsztalt(
        await straz.wywolaj(Command.ModuleList, {}),
        Command.ModuleList,
        (tresc) => czyTablica(tresc.modules),
      );
    },

    async komplet(idOkna) {
      return sprawdzKsztalt(
        await straz.wywolaj(Command.AodContextGet, { windowId: idOkna }),
        Command.AodContextGet,
        // Sprawdzany jest sam komplet, nie jego zawartość: pusty komplet jest odpowiedzią, nie usterką.
        (tresc) => czyObiekt(tresc.context),
      );
    },

    naZmianeOkna(sluchacz) {
      return kanal.naZdarzenie(EventType.WindowChanged, (tresc) => sluchacz(tresc));
    },
  };
}
