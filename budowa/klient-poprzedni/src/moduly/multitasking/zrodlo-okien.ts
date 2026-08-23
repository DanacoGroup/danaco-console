import {
  Command,
  type ConfigEntry,
  type ConfigGetRequest,
  type ConfigSetRequest,
  type ConfigSessionGetRequest,
  type RoleAssignRequest,
  type RoleAssignResponse,
  type RoleListRequest,
  type RoleRemoveRequest,
  type RoleRemoveResponse,
  type RoleUpdateRequest,
  type RoleUpdateResponse,
  type SessionConfig,
  type WindowRoleAssignment,
  type Window,
  type WindowCreateRequest,
  type WindowListRequest,
  type WindowStateGetRequest,
  type WindowStateGetResponse,
  type WindowUpdateRequest,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import {
  czyLogiczna,
  czyObiekt,
  czyTablica,
  czyTekst,
  sprawdzKsztalt,
} from '../../protokol/ksztalt-odpowiedzi';
import { przenies } from '../../protokol/wynik-czastkowy';
import { wywolaj } from '../../protokol/wywolanie';

/**
 * Okna ról widziane przez klienta — ta połowa źródła, która ustala obsadę.
 *
 * Rola okna pochodzi z kontraktu, nie z oznaczenia własnego. `WindowRole`
 * (standalone | executor | coordinator) i pole `coordinatorWindowId` niosą całą
 * przynależność okna do pętli koordynator–wykonawca. Moduł nie zakłada
 * drugiego rejestru ról: nadanie roli to komenda rdzenia, a nie własna tablica
 * w kliencie.
 *
 * Nadanie roli idzie komendą przeznaczoną do nadawania roli (`role.assign`),
 * a wcielenie (`persona`) — komendą `role.update`, bo `window.update` pola
 * wcielenia nie niesie. `window.update` zostaje przy tym, co jest jego: tytuł,
 * katalogi, kanał modelu.
 */
export interface ZrodloOkien {
  /** `window.list` — okna sesji, z których powstaje obsada ról. */
  okna(zadanie: WindowListRequest): Promise<Wynik<Window[]>>;
  /** `window.create` — założenie okna roli wraz z przypięciem do koordynatora. */
  zalozOkno(zadanie: WindowCreateRequest): Promise<Wynik<Window>>;
  /** `window.update` — tytuł, katalogi i kanał modelu okna istniejącego. */
  zmienOkno(zadanie: WindowUpdateRequest): Promise<Wynik<Window>>;
  /** `role.assign` — nadanie oknu roli w pętli koordynator–wykonawca. */
  nadajRole(zadanie: RoleAssignRequest): Promise<Wynik<RoleAssignResponse>>;
  /** `role.update` — zmiana roli albo jej wcielenia (`persona`). */
  zmienRole(zadanie: RoleUpdateRequest): Promise<Wynik<RoleUpdateResponse>>;
  /**
   * `role.list` — nadania ról widziane przez rdzeń, wraz z więzią koordynatora.
   *
   * To nie jest to samo co `window.list`. Okno niesie swoją rolę jako pole, ale
   * rejestrem nadań jest rdzeń: `role.list` oddaje `WindowRoleAssignment`, więc
   * mówi też o wcieleniu (`persona`), którego `Window` nie niesie w ogóle.
   */
  nadaniaRol(zadanie: RoleListRequest): Promise<Wynik<WindowRoleAssignment[]>>;
  /** `role.remove` — zdjęcie roli z okna wraz z więzią koordynatora. */
  zdejmijRole(zadanie: RoleRemoveRequest): Promise<Wynik<RoleRemoveResponse>>;
  /** `window.state.get` — stan okna wraz z licznikiem obiegów pętli. */
  stanOkna(zadanie: WindowStateGetRequest): Promise<Wynik<WindowStateGetResponse>>;
  /** `config.get` — ustawienia zapisane na poziomie okna albo sesji. */
  ustawienia(zadanie: ConfigGetRequest): Promise<Wynik<ConfigEntry[]>>;
  /** `config.set` — zapis ustawienia roli na jego poziomie zasięgu. */
  zapiszUstawienie(zadanie: ConfigSetRequest): Promise<Wynik<ConfigEntry>>;
  /** `config.session.get` — obszary konfiguracji sesji, w tym definicje podagentów. */
  konfiguracjaSesji(zadanie: ConfigSessionGetRequest): Promise<Wynik<SessionConfig>>;
}

export function utworzZrodloOkien(kanal: Kanal): ZrodloOkien {
  return {
    async okna(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.WindowList, zadanie),
        Command.WindowList,
        (tresc) => czyTablica(tresc.windows),
      );
      return przenies(wynik, (tresc) => tresc.windows);
    },

    async zalozOkno(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.WindowCreate, zadanie),
        Command.WindowCreate,
        (tresc) => czyObiekt(tresc.window),
      );
      return przenies(wynik, (tresc) => tresc.window);
    },

    async zmienOkno(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.WindowUpdate, zadanie),
        Command.WindowUpdate,
        (tresc) => czyObiekt(tresc.window),
      );
      return przenies(wynik, (tresc) => tresc.window);
    },

    // Kształt odpowiedzi obu komend ról sprawdzamy tak samo jak każdej innej:
    // rdzeń, który odpowie `ok` bez pola `role`, nie dowiódł, że rolę nadał.
    async nadajRole(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.RoleAssign, zadanie),
        Command.RoleAssign,
        (tresc) => czyTekst(tresc.windowId) && czyTekst(tresc.role),
      );
    },

    async zmienRole(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.RoleUpdate, zadanie),
        Command.RoleUpdate,
        (tresc) => czyTekst(tresc.windowId) && czyTekst(tresc.role),
      );
    },

    async nadaniaRol(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.RoleList, zadanie),
        Command.RoleList,
        (tresc) => czyTablica(tresc.assignments),
      );
      return przenies(wynik, (tresc) => tresc.assignments);
    },

    // `removed: false` to odpowiedź udana o skutku, którego nie było — rdzeń
    // mówi „takiej roli nie miałem czego zdjąć". Kształt sprawdzamy więc na
    // obecności pola logicznego, a jego wartość czyta widok i nazywa wprost:
    // samo `status: ok` nie dowodzi, że rola została zdjęta.
    async zdejmijRole(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.RoleRemove, zadanie),
        Command.RoleRemove,
        (tresc) => czyTekst(tresc.windowId) && czyLogiczna(tresc.removed),
      );
    },

    async stanOkna(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.WindowStateGet, zadanie),
        Command.WindowStateGet,
        (tresc) => czyObiekt(tresc.window),
      );
    },

    async ustawienia(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.ConfigGet, zadanie),
        Command.ConfigGet,
        (tresc) => czyTablica(tresc.entries),
      );
      return przenies(wynik, (tresc) => tresc.entries);
    },

    async zapiszUstawienie(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.ConfigSet, zadanie),
        Command.ConfigSet,
        (tresc) => czyObiekt(tresc.entry),
      );
      return przenies(wynik, (tresc) => tresc.entry);
    },

    async konfiguracjaSesji(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.ConfigSessionGet, zadanie),
        Command.ConfigSessionGet,
        (tresc) => czyObiekt(tresc.config),
      );
      return przenies(wynik, (tresc) => tresc.config);
    },
  };
}
