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
 * Otoczenie modułu Library — komendy, z których moduł korzysta, a których nie
 * prowadzi.
 *
 *   `window.list`      — żywe okno komunikacji sesji. Przenoszenie kontekstu
 *                        żąda okna źródłowego, a moduł zna wyłącznie kody
 *                        okien operacyjnych katalogu rdzenia; to dwa różne
 *                        byty, więc okno bierze się z rejestru.
 *   `window.state.get` — kontekst okna wymagany przez wiersz każdego z czterech
 *                        okien wykazu; moduł pyta raz i dzieli odpowiedź.
 *   `action.list`      — katalog akcji panelu. Panel akcji nie jest zaszytym
 *                        wykazem po stronie klienta: pozycje przychodzą
 *                        z rdzenia albo panel zostaje pusty i mówi to wprost.
 *   `window.action`    — wykonanie pozycji panelu; rdzeń nie ma uchwytu, więc
 *                        odpowiada odmową `window.unknown`.
 *   `context.transfer` — otwarcie zasobu w module źródłowym.
 *   `module.list`      — katalog modułów platformy. Obsadza ster modułu
 *                        docelowego: rdzeń przyjmuje każdy kod i zakłada okno
 *                        z dokładnie tym kodem, więc kodu wpisanego z ręki nie
 *                        miałby kto sprawdzić.
 *   `aod.context.get`  — treść kompletu przeniesionego do okna. Jedyna komenda
 *                        kontraktu, która czyta magazyn zapisany przez
 *                        `context.transfer` (`adapterPrzenoszenia.KompletOkna`).
 *                        Prefiks `aod.*` należy do nakładki, ale komenda jest
 *                        odczytem kompletu okna — bez niej przybycie
 *                        przekazania jest dla modułu nieme.
 *   `window.changed`   — zdarzenie przybycia. `context.transfer` nie rozgłasza
 *                        `library.file.changed` i nie zakłada pliku
 *                        w repozytorium; rozgłasza wyłącznie `window.changed`
 *                        z oknem docelowym.
 */
export interface ZrodloOtoczenia {
  okna(idSesji: string): Promise<Wynik<{ windows: Window[] }>>;
  stanOkna(idOkna: string): Promise<Wynik<{ window: Window; messageCount: number }>>;
  akcje(): Promise<Wynik<{ actions: Action[] }>>;
  wykonajAkcje(idOkna: string, idAkcji: string): Promise<Wynik<{ actionId: string }>>;
  /**
   * Przenosi komplet kontekstu i oddaje okno, które rdzeń wskazał jako docelowe.
   *
   * Rdzeń przyjmuje każdy kod modułu — także nieistniejący i pusty — i zakłada
   * okno z dokładnie tym kodem, więc jedynym świadkiem tego, dokąd komplet
   * trafił, jest okno z odpowiedzi, a nie kod wpisany w oknie.
   */
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

/** Kod modułu z kolumny `modul.kod` rdzenia — zasięg katalogu akcji. */
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
        // Sprawdzany jest sam komplet, nie jego zawartość: rdzeń oddaje go
        // także dla okna bez dokumentów. Pusty komplet jest odpowiedzią,
        // nie usterką kształtu — rozstrzyga go wołający.
        (tresc) => czyObiekt(tresc.context),
      );
    },

    naZmianeOkna(sluchacz) {
      return kanal.naZdarzenie(EventType.WindowChanged, (tresc) => sluchacz(tresc));
    },
  };
}
