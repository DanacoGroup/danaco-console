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

/** Czynności bramki wykonywane po zalogowaniu obejmują komendy uwierzytelniania należące do okna ustawień, nie do ekranu logowania, i różnią się od źródła logowania właśnie tym momentem wywołania. */
export interface ZrodloBramki {
  /** Zakłada metodę szybkiego wejścia; hasła nie zakłada, bo kotwica powstaje przy rejestracji. */
  zalozMetode(zadanie: AuthMethodAddRequest): Promise<Wynik<{ methods: AuthMethod[] }>>;
  /** Zdejmuje metodę szybkiego wejścia z urządzenia; ani ostatniej metody, ani hasła zdjąć się nie da. */
  zdejmijMetode(
    idMetody: string,
    idUrzadzenia?: string,
  ): Promise<Wynik<{ methods: AuthMethod[]; removed: boolean }>>;
  /** Zmienia hasło ze znanym hasłem bieżącym; wynik niesie liczbę sesji unieważnionych zmianą. */
  zmienHaslo(
    biezace: string,
    nowe: string,
  ): Promise<Wynik<AuthPasswordResetResponse>>;
  /** Jawne przedłużenie sesji bramki. */
  przedluzSesje(token: string): Promise<Wynik<AuthTokenRefreshResponse>>;
  /** Subskrypcja zmiany bramki: jedyna droga, którą wykaz metod dociera bez czynności w tym oknie. */
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

