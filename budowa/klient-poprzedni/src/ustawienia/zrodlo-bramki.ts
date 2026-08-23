import {
  Command,
  EventType,
  type AuthChangedEvent,
  type AuthMethod,
  type AuthMethodAddRequest,
  type AuthPasswordResetResponse,
  type AuthTokenRefreshResponse,
} from '../../../shared/contract';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import type { Kanal, Wynik } from '../protokol/kanal';
import { czyTablica, sprawdzKsztalt } from '../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../protokol/wywolanie';

/**
 * Czynności bramki wykonywane po zalogowaniu — komendy `auth.*` należące do
 * Okna Ustawień, nie do ekranu logowania.
 *
 * Rozdział względem `uwierzytelnienie/zrodlo-auth.ts` przebiega po bramce, nie
 * po rodzinie komend: tamto źródło woła `auth.login`, `auth.register`
 * i `auth.token.refresh`, czyli tyle, ile trzeba, żeby wejść. Tych dwóch komend
 * tu nie ma, bo w Oknie Ustawień odmawiałyby zawsze: `auth.register` jest
 * wykonalna tylko raz, a bramka jest już założona w chwili, gdy okno da się
 * otworzyć; `auth.login` jest samą bramką i po jej przejściu nie ma czego
 * otwierać.
 *
 * Metody `hello` (Windows Hello przez WebAuthn) rdzeń nie ma zbudowanej
 * i odmawia jej założenia; przyjmuje hasło i PIN.
 *
 * Komendy odczytu wykazu metod kontrakt nie niesie. Wykaz dociera zdarzeniem
 * `auth.changed`, które rdzeń rozsyła do wszystkich gniazd — także po czynności
 * wykonanej w innym oknie.
 */
export interface ZrodloBramki {
  /**
   * `auth.method.add` — założenie metody szybkiego wejścia na urządzeniu.
   *
   * Hasła ta komenda nie zakłada: kotwica powstaje przy `auth.register`.
   * Przyjmowana metoda to PIN.
   */
  zalozMetode(zadanie: AuthMethodAddRequest): Promise<Wynik<{ methods: AuthMethod[] }>>;
  /**
   * `auth.method.remove` — zdjęcie metody szybkiego wejścia z urządzenia.
   *
   * Ani ostatniej metody, ani hasła zdjąć się nie da — hasło jest kotwicą
   * bramki. Odmowę orzeka rdzeń; okno jej nie uprzedza.
   */
  zdejmijMetode(
    idMetody: string,
    idUrzadzenia?: string,
  ): Promise<Wynik<{ methods: AuthMethod[]; removed: boolean }>>;
  /**
   * `auth.password.reset` — zmiana hasła ze znanym hasłem bieżącym.
   *
   * Drogi odzyskania listem na adres e-mail nie ma — rdzeń poczty nie wysyła.
   * Wynik niesie liczbę sesji unieważnionych zmianą, poza bieżącą.
   */
  zmienHaslo(
    biezace: string,
    nowe: string,
  ): Promise<Wynik<AuthPasswordResetResponse>>;
  /** `auth.token.refresh` — jawne przedłużenie sesji bramki. */
  przedluzSesje(token: string): Promise<Wynik<AuthTokenRefreshResponse>>;
  /**
   * Subskrypcja `auth.changed` — jedyna droga, którą wykaz metod dociera bez
   * czynności wykonanej w tym oknie.
   *
   * Zdarzenie niesie powód zmiany i komplet metod po niej. Pole `methods` jest
   * w kontrakcie opcjonalne; gdy go brak, słuchacz dostaje sam powód —
   * podstawianie w to miejsce wykazu poprzedniego byłoby zgadywaniem za rdzeń.
   */
  naZmianeBramki(sluchacz: (zmiana: AuthChangedEvent) => void): Odsubskrybuj;
}

export function utworzZrodloBramki(kanal: Kanal): ZrodloBramki {
  return {
    async zalozMetode(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AuthMethodAdd, zadanie),
        Command.AuthMethodAdd,
        (tresc) => czyTablica(tresc.methods),
      );
    },

    async zdejmijMetode(idMetody, idUrzadzenia) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AuthMethodRemove, {
          methodId: idMetody,
          ...(idUrzadzenia === undefined || idUrzadzenia === '' ? {} : { deviceId: idUrzadzenia }),
        }),
        Command.AuthMethodRemove,
        (tresc) => czyTablica(tresc.methods),
      );
    },

    async zmienHaslo(biezace, nowe) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AuthPasswordReset, {
          currentPassword: biezace,
          newPassword: nowe,
        }),
        Command.AuthPasswordReset,
        (tresc) => typeof tresc.changed === 'boolean',
      );
    },

    async przedluzSesje(token) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AuthTokenRefresh, token === '' ? {} : { token }),
        Command.AuthTokenRefresh,
        (tresc) => tresc.session !== undefined,
      );
    },

    naZmianeBramki(sluchacz) {
      return kanal.naZdarzenie(EventType.AuthChanged, (tresc) => sluchacz(tresc));
    },
  };
}

