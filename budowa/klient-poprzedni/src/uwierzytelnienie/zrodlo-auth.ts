import {
  AuthMethodKind,
  Command,
  ErrorCode,
  type AuthLoginResponse,
  type AuthPasswordResetResponse,
  type AuthRecoverResponse,
  type AuthRegisterResponse,
  type AuthResetResponse,
  type AuthTokenRefreshResponse,
  type AuthVerifyResponse,
  type ConnectionHelloResponse,
  type ErrorInfo,
  type RequestOf,
  type ResponseOf,
} from '../../../shared/contract';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import type { Kanal } from '../protokol/kanal';
import { czyOdpowiedz, czyUdana, tresc } from '../protokol/koperta';

// Źródło komend bramki jest jedynym miejscem znającym nazwy komend kontraktu i kształt żądań.

/** Odpowiedź rdzenia wraz z typem koperty, którym przyszła, pozwala odróżnić brak uchwytu komendy od zwykłej odmowy. */
export interface OdpowiedzBramki<T> {
  /** Czy koperta niesie wynik (`status: "ok"`). */
  udana: boolean;
  /** Typ koperty odpowiedzi — `auth.login` albo `auth.unknown` przy braku uchwytu. */
  typ: string;
  /** Treść wyniku; wyłącznie przy odpowiedzi udanej. */
  wynik?: T;
  /** Powód odmowy; wyłącznie przy odpowiedzi błędnej. */
  blad?: ErrorInfo;
}

/** Stan bramki przed pokazaniem formularza: nieustawiona, ustawiona, bez uchwytu rdzenia albo niepewny wobec nietypowej odmowy. */
export type StanBramki = 'nieustawiona' | 'ustawiona' | 'bez-uchwytu' | 'niepewny';

/** Wynik rozpoznania stanu bramki niesie sam stan wraz z powodem odmowy, gdy stan wynika właśnie z odmowy rdzenia. */
export interface RozpoznanieBramki {
  stan: StanBramki;
  blad?: ErrorInfo;
}

/** Metoda wejścia widziana przez operatora obejmuje wyłącznie hasło i PIN, bo tyle rdzeń dziś naprawdę otwiera. */
export type MetodaWejscia = 'haslo' | 'pin';

/** Żądanie wejścia złożone z pól formularza logowania, niesione do rdzenia komendą wejścia przez bramkę. */
export interface WejscieBramki {
  metoda: MetodaWejscia;
  /** Login jest nazwą właściciela nadaną przy rejestracji, wymaganą wyłącznie przy metodzie hasła. */
  login?: string;
  /** Hasło albo PIN — zależnie od metody. */
  sekret: string;
  /** Przełącznik „nie wyloguj mnie". */
  niewylogowuj: boolean;
  /** Maszyna, do której PIN należy; przy haśle bez znaczenia. */
  urzadzenie?: string;
}

/** Żądanie rejestracji złożone z formularza pierwszego uruchomienia nie zakłada sesji, bo tę wydaje dopiero potwierdzenie adresu. */
export interface ZalozenieKonta {
  /** Nazwa właściciela używana potem przy logowaniu. */
  login: string;
  /** Adres e-mail uwierzytelniający jest weryfikowany przy rejestracji i jedyną drogą odzyskania konta. */
  adres: string;
  haslo: string;
}

/**
 * Trzy pola powitania rozstrzygające o przesłonie.
 *
 * Każde puste znaczy „rdzeń nie wie", nigdy „nie": rdzeń bez wpiętej bramki
 * albo bez znanej nastawy milczy, a klient tego milczenia nie zamienia na
 * rozstrzygnięcie.
 */
export interface Powitanie {
  /** Czy na tym nasłuchu obowiązuje wymóg logowania (`gateway.requireLogin`). */
  loginRequired?: boolean;
  /** Czy sekret bramki jest już ustawiony — rozstrzyga wejście vs pierwsze hasło. */
  gatewayConfigured?: boolean;
  /** Czy TO połączenie jest już związane z żywą sesją bramki. */
  authenticated?: boolean;
}

/** Czynności bramki widziane przez ekran logowania obejmują powitanie, sondy stanu oraz wejście, rejestrację i odzyskanie konta. */
export interface ZrodloUwierzytelnienia {
  /** Powitanie połączenia jest podsłuchane z kanału, nie wysyłane osobno, i ma zadany czas oczekiwania. */
  powitanie(limitMs?: number): Promise<Powitanie>;
  /** Sonda stanu bramki — rozstrzyga, który formularz pokazać. */
  zbadaj(): Promise<RozpoznanieBramki>;
  /** Sonda PIN-u tej maszyny sprawdza, czy PIN jest tu w ogóle założony, bez podnoszenia dławika prób. */
  zbadajPin(urzadzenie: string): Promise<boolean>;
  /** Wejście przez bramkę hasłem albo PIN-em niesie do rdzenia także przełącznik przedłużenia sesji. */
  wejdz(zadanie: WejscieBramki): Promise<OdpowiedzBramki<AuthLoginResponse>>;
  /** Rejestracja zakłada jedyne konto właściciela przy pierwszym uruchomieniu i nie zakłada sesji. */
  zaloz(zadanie: ZalozenieKonta): Promise<OdpowiedzBramki<AuthRegisterResponse>>;
  /** Potwierdzenie adresu drogą z listu zamyka rejestrację i wydaje urządzeniu token dostępu. */
  potwierdz(
    droga: string,
    niewylogowuj: boolean,
  ): Promise<OdpowiedzBramki<AuthVerifyResponse>>;
  /** Prośba o drogę odzyskania konta idzie na adres uwierzytelniający i zawsze niesie tę samą odpowiedź. */
  odzyskaj(adres: string): Promise<OdpowiedzBramki<AuthRecoverResponse>>;
  /** `auth.reset` — ustawienie nowego hasła po potwierdzeniu odzyskania. */
  ustawNoweHaslo(
    droga: string,
    nowe: string,
  ): Promise<OdpowiedzBramki<AuthResetResponse>>;
  /** `auth.token.refresh` — przedłużenie sesji zapisanej z poprzedniego uruchomienia. */
  przedluz(token: string): Promise<OdpowiedzBramki<AuthTokenRefreshResponse>>;
  /** Zmiana hasła ze znanym hasłem dotychczasowym jest wykonalna przed zalogowaniem, bez sesji. */
  zmienHaslo(
    biezace: string,
    nowe: string,
  ): Promise<OdpowiedzBramki<AuthPasswordResetResponse>>;
}

/** Wysyłka jednej komendy odczytuje całą kopertę odpowiedzi, a subskrypcja staje przed wysłaniem, więc odpowiedź nie ma jak przepaść. */
function poslij<K extends Command>(
  kanal: Kanal,
  komenda: K,
  zadanie: RequestOf<K>,
): Promise<OdpowiedzBramki<ResponseOf<K>>> {
  return new Promise((rozstrzygnij) => {
    let id = '';
    const odepnij = kanal.naDowolny((koperta) => {
      if (koperta.id !== id || !czyOdpowiedz(koperta)) return;
      odepnij();
      const udana = czyUdana(koperta);
      rozstrzygnij({
        udana,
        typ: koperta.type,
        wynik: udana ? tresc<ResponseOf<K>>(koperta) : undefined,
        blad: koperta.error,
      });
    });
    id = kanal.wyslij(komenda, zadanie);
  });
}

/** Funkcja rozstrzyga, czy koperta odpowiedzi jest odmową oznaczającą brak uchwytu komendy w rdzeniu serwera. */
export function bezUchwytu(odpowiedz: OdpowiedzBramki<unknown>): boolean {
  return !odpowiedz.udana && odpowiedz.typ.endsWith('.unknown');
}

/** Funkcja zamienia metodę wejścia widzianą przez operatora na rodzaj metody zgodny z kontraktem rdzenia. */
function rodzajMetody(metoda: MetodaWejscia): AuthMethodKind {
  return metoda === 'pin' ? AuthMethodKind.Pin : AuthMethodKind.Password;
}

/** Funkcja buduje źródło czynności bramki działające nad kanałem połączenia z rdzeniem serwera aplikacji. */
export function utworzZrodloUwierzytelnienia(kanal: Kanal): ZrodloUwierzytelnienia {
  // Zapisane powitanie zostaje na później, bo koperta powitania przechodzi przez kanał tylko raz.
  let uslyszane: Powitanie | null = null;
  const czekajacy: ((powitanie: Powitanie) => void)[] = [];
  let odepnijPowitanie: Odsubskrybuj = () => undefined;

  odepnijPowitanie = kanal.naDowolny((koperta) => {
    if (koperta.type !== Command.ConnectionHello || !czyOdpowiedz(koperta)) return;
    // Odmowa powitania jest milczeniem rdzenia, nie rozstrzygnięciem o wymogu logowania.
    const wynik = czyUdana(koperta) ? tresc<ConnectionHelloResponse>(koperta) : undefined;
    uslyszane = {
      loginRequired: wynik?.loginRequired,
      gatewayConfigured: wynik?.gatewayConfigured,
      authenticated: wynik?.authenticated,
    };
    odepnijPowitanie();
    for (const oddaj of czekajacy.splice(0)) oddaj(uslyszane);
  });

  return {
    powitanie(limitMs = 1500) {
      if (uslyszane !== null) return Promise.resolve(uslyszane);
      if (limitMs <= 0) return Promise.resolve({});
      return new Promise<Powitanie>((rozstrzygnij) => {
        let rozstrzygniete = false;
        const zamknij = (powitanie: Powitanie): void => {
          if (rozstrzygniete) return;
          rozstrzygniete = true;
          rozstrzygnij(powitanie);
        };
        czekajacy.push(zamknij);
        window.setTimeout(() => zamknij({}), limitMs);
      });
    },

    /** Sonda stanu bramki wysyła wejście metodą hasła bez sekretu i czyta rozstrzygnięcie rdzenia z odmowy. */
    async zbadaj() {
      const odpowiedz = await poslij(kanal, Command.AuthLogin, {
        method: AuthMethodKind.Password,
      });
      if (bezUchwytu(odpowiedz)) return { stan: 'bez-uchwytu', blad: odpowiedz.blad };
      if (odpowiedz.blad?.code === ErrorCode.NotFound) return { stan: 'nieustawiona' };
      if (odpowiedz.blad?.code === ErrorCode.ValidationFailed) return { stan: 'ustawiona' };
      // Odpowiedź udana jest tu niemożliwa; każda inna odmowa idzie do ekranu z powodem, nie jest ukrywana.
      return { stan: 'niepewny', blad: odpowiedz.blad };
    },

    async zbadajPin(urzadzenie) {
      if (urzadzenie === '') return false;
      const odpowiedz = await poslij(kanal, Command.AuthLogin, {
        method: AuthMethodKind.Pin,
        deviceId: urzadzenie,
      });
      // Wyłącznie odmowa walidacji znaczy, że PIN tu jest; każda inna zostawia segment PIN-u schowany.
      return !bezUchwytu(odpowiedz) && odpowiedz.blad?.code === ErrorCode.ValidationFailed;
    },

    // Identyfikator urządzenia idzie tylko przy PIN-ie, bo hasło otwiera bramkę z każdej maszyny.
    wejdz: (zadanie) =>
      poslij(kanal, Command.AuthLogin, {
        method: rodzajMetody(zadanie.metoda),
        secret: zadanie.sekret,
        keepSignedIn: zadanie.niewylogowuj,
        // Login idzie wyłącznie przy metodzie hasła, bo metody właściwe urządzeniu go nie czytają.
        ...(zadanie.metoda === 'haslo' && zadanie.login !== undefined && zadanie.login !== ''
          ? { login: zadanie.login }
          : {}),
        ...(zadanie.metoda === 'pin' && zadanie.urzadzenie !== undefined && zadanie.urzadzenie !== ''
          ? { deviceId: zadanie.urzadzenie }
          : {}),
      }),

    zaloz: (zadanie) =>
      poslij(kanal, Command.AuthRegister, {
        login: zadanie.login,
        email: zadanie.adres,
        password: zadanie.haslo,
      }),

    potwierdz: (droga, niewylogowuj) =>
      poslij(kanal, Command.AuthVerify, { token: droga, keepSignedIn: niewylogowuj }),

    odzyskaj: (adres) => poslij(kanal, Command.AuthRecover, { email: adres }),

    ustawNoweHaslo: (droga, nowe) =>
      poslij(kanal, Command.AuthReset, { token: droga, newPassword: nowe }),

    przedluz: (token) => poslij(kanal, Command.AuthTokenRefresh, { token }),

    zmienHaslo: (biezace, nowe) =>
      poslij(kanal, Command.AuthPasswordReset, {
        currentPassword: biezace,
        newPassword: nowe,
      }),
  };
}
