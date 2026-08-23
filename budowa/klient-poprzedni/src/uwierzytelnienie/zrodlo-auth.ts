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

/**
 * Źródło komend bramki — jedyne miejsce w katalogu `uwierzytelnienie/`, które
 * zna nazwy komend kontraktu i kształt ich żądań (ten sam wzorzec co
 * `aod/zrodlo-komend.ts`). Ekran logowania dostaje stąd czynności, nie
 * `Command.*`.
 *
 * Z siedmiu komend rodziny `auth.*` ekran woła cztery: `auth.login`,
 * `auth.register`, `auth.token.refresh` — tyle, ile trzeba, żeby wejść — oraz
 * `auth.password.reset`, bo zmiana hasła nie wymaga zalogowania i stoi za
 * odnośnikiem „Reset hasła". `auth.method.add` i `auth.method.remove` są
 * czynnościami Ustawień wykonywanymi po zalogowaniu. `auth.confirm` nie jest
 * wołana, bo uchwytu w rdzeniu nie ma (`core/handlers_auth.go`) — wywołanie
 * wróciłoby kopertą `auth.unknown` i niczym więcej.
 *
 * Odpowiedź niesie typ koperty, nie tylko wynik. Rdzeń odpowiada na komendę
 * bez uchwytu kopertą `auth.unknown` z kodem `not_found` — tym samym kodem,
 * którym `auth.login` odmawia przy bramce nieustawionej. Sam `Wynik` kanału
 * tych dwóch odmów nie odróżnia, a znaczą one co innego: pierwsza to rdzeń bez
 * wpiętej rodziny `auth.*`, druga to pierwsze uruchomienie. Dlatego wysyłka
 * idzie tu przez `naDowolny` i oddaje typ koperty odpowiedzi razem z treścią.
 */

/** Odpowiedź rdzenia wraz z typem koperty, którym przyszła. */
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

/**
 * Rozpoznanie stanu bramki przed pokazaniem formularza:
 *
 *   `nieustawiona` — kotwicy nie ma; ekran pokazuje ustawienie hasła (`auth.register`),
 *   `ustawiona`    — kotwica jest; ekran pokazuje wejście hasłem (`auth.login`),
 *   `bez-uchwytu`  — rdzeń nie ma wpiętej rodziny `auth.*` (koperta `auth.unknown`),
 *   `niepewny`     — odmowa spoza dwóch spodziewanych, np. rdzeń bez sejfu poświadczeń.
 */
export type StanBramki = 'nieustawiona' | 'ustawiona' | 'bez-uchwytu' | 'niepewny';

/** Wynik rozpoznania: stan wraz z powodem, gdy stan wynika z odmowy. */
export interface RozpoznanieBramki {
  stan: StanBramki;
  blad?: ErrorInfo;
}

/**
 * Metoda wejścia widziana przez Operatora. Dwie, bo tyle rdzeń otwiera dziś
 * naprawdę (`adapter_modul_auth.go`, `metodaWejscia`): hasło i PIN. Windows
 * Hello odmawia z powodu pochodzenia dokumentu, nie z braku kodu, więc nie
 * jest tu metodą — segment wyszarzony byłby bramą, a segment czynny kłamstwem.
 */
export type MetodaWejscia = 'haslo' | 'pin';

/** Żądanie wejścia złożone z formularza. */
export interface WejscieBramki {
  metoda: MetodaWejscia;
  /**
   * Login — nazwa właściciela nadana przy rejestracji. Wymagany przy metodzie
   * `haslo` (kontrakt: „wymagana dla metody password"); metody właściwe
   * urządzeniu go nie potrzebują, bo wskazuje je materiał na maszynie.
   */
  login?: string;
  /** Hasło albo PIN — zależnie od metody. */
  sekret: string;
  /** Przełącznik „nie wyloguj mnie". */
  niewylogowuj: boolean;
  /** Maszyna, do której PIN należy; przy haśle bez znaczenia. */
  urzadzenie?: string;
}

/**
 * Żądanie rejestracji złożone z formularza pierwszego uruchomienia.
 *
 * Przełącznika „nie wyloguj mnie" tu nie ma z rozmysłem: rejestracja nie zakłada
 * sesji, więc nie ma czego przedłużać. Wybór trwania idzie dopiero
 * z potwierdzeniem adresu, które sesję wydaje.
 */
export interface ZalozenieKonta {
  /** Nazwa właściciela używana potem przy logowaniu. */
  login: string;
  /**
   * Adres e-mail uwierzytelniający — weryfikowany przy rejestracji, później
   * jedyna droga odzyskania konta.
   */
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

/** Czynności bramki widziane przez ekran logowania. */
export interface ZrodloUwierzytelnienia {
  /**
   * Powitanie połączenia (`connection.hello`) — podsłuchane, nie wysłane.
   *
   * Rdzeń podaje w powitaniu wszystko, czego ekran potrzebuje do rozpoznania:
   * `loginRequired`, `gatewayConfigured`, `authenticated`. Powitanie i tak leci
   * przy każdym nawiązaniu połączenia (`protokol/uzgodnienie.ts`, wywoływane ze
   * stanu transportu), a warstwa uzgodnienia bierze z odpowiedzi samo
   * powodzenie etapu. Zamiast dokładać drugie powitanie z własnym
   * identyfikatorem klienta, bramka słucha koperty przechodzącej przez kanał
   * (`kanal.naDowolny`). Subskrypcja staje przy składaniu źródła — zanim
   * gniazdo się otworzy — więc odpowiedź nie ma jak przepaść.
   *
   * Milczenie ma kres: gdy powitanie nie wróci w zadanym czasie (rdzeń bez
   * uchwytu `connection.hello`, gniazdo, które nie doszło do skutku), wynikiem
   * jest powitanie puste, czyli „rdzeń nie wie", i ekran schodzi na sondę
   * zamiast wisieć.
   *
   * @param limitMs kres oczekiwania; `0` znaczy „nie czekaj, weź, co już jest".
   */
  powitanie(limitMs?: number): Promise<Powitanie>;
  /** Sonda stanu bramki — rozstrzyga, który formularz pokazać. */
  zbadaj(): Promise<RozpoznanieBramki>;
  /**
   * Sonda PIN-u tej maszyny: czy PIN w ogóle jest tu założony.
   *
   * Idzie tą samą drogą co sonda hasła i z tego samego powodu: kontrakt nie ma
   * komendy odczytu metod (metody oddaje dopiero odpowiedź na czynność),
   * a segment „PIN" nie ma prawa stanąć na ekranie, jeżeli nie ma czego
   * otworzyć. Rdzeń szuka metody po parze rodzaj + urządzenie, zanim spojrzy
   * na sekret, więc żądanie bez sekretu odpowiada wprost: `not_found` → PIN-u
   * tu nie ma, `validation_failed` → PIN jest, zabrakło tylko sekretu
   * (`adapter_modul_auth.go`).
   *
   * Sonda nie podnosi dławika prób: zwłokę zwiększa wyłącznie sekret niezgodny,
   * a tu sekretu nie ma.
   */
  zbadajPin(urzadzenie: string): Promise<boolean>;
  /**
   * `auth.login` — wejście przez bramkę hasłem albo PIN-em.
   *
   * `niewylogowuj` to przełącznik „nie wyloguj mnie" z formularza. Idzie do
   * rdzenia, nie tylko do magazynu przeglądarki: bez tego zaznaczenie
   * przedłużałoby wyłącznie życie zapisu w `localStorage`, a sesja gasłaby po
   * dobie roboczej.
   */
  wejdz(zadanie: WejscieBramki): Promise<OdpowiedzBramki<AuthLoginResponse>>;
  /**
   * `auth.register` — założenie jedynego konta właściciela przy pierwszym
   * uruchomieniu: login, adres e-mail uwierzytelniający i hasło.
   *
   * Sesji nie zakłada. Konto powstaje niepotwierdzone, a token dostępu wydaje
   * dopiero `potwierdz` po drodze przysłanej listem — dlatego odpowiedź nie
   * niesie sesji i ekran przechodzi wtedy do kroku potwierdzenia.
   */
  zaloz(zadanie: ZalozenieKonta): Promise<OdpowiedzBramki<AuthRegisterResponse>>;
  /**
   * `auth.verify` — potwierdzenie adresu drogą z listu; zamyka rejestrację
   * i wydaje urządzeniu token dostępu.
   */
  potwierdz(
    droga: string,
    niewylogowuj: boolean,
  ): Promise<OdpowiedzBramki<AuthVerifyResponse>>;
  /**
   * `auth.recover` — prośba o drogę odzyskania konta wysyłaną na adres
   * uwierzytelniający.
   *
   * Odpowiedź jest zawsze taka sama, niezależnie od tego, czy adres pasuje do
   * konta. Ekran nie ma z czego wywnioskować, czy konto istnieje — i tak ma być.
   */
  odzyskaj(adres: string): Promise<OdpowiedzBramki<AuthRecoverResponse>>;
  /** `auth.reset` — ustawienie nowego hasła po potwierdzeniu odzyskania. */
  ustawNoweHaslo(
    droga: string,
    nowe: string,
  ): Promise<OdpowiedzBramki<AuthResetResponse>>;
  /** `auth.token.refresh` — przedłużenie sesji zapisanej z poprzedniego uruchomienia. */
  przedluz(token: string): Promise<OdpowiedzBramki<AuthTokenRefreshResponse>>;
  /**
   * `auth.password.reset` — zmiana hasła ze znanym hasłem dotychczasowym.
   *
   * Należy do ekranu wejścia, choć wygląda na czynność Ustawień: handler
   * (`core/adapter_modul_auth_metody.go`) nie pyta o sesję, sprawdza wyłącznie
   * zgodność hasła dotychczasowego z zapisem bramki. Zmiana jest więc wykonalna
   * przed zalogowaniem i odnośnik „Reset hasła" ma co wywołać. Odzyskaniem
   * hasła zapomnianego to nie jest — takiej drogi kontrakt nie ma.
   */
  zmienHaslo(
    biezace: string,
    nowe: string,
  ): Promise<OdpowiedzBramki<AuthPasswordResetResponse>>;
}

/**
 * Wysyłka jednej komendy z odczytem całej koperty odpowiedzi.
 *
 * Subskrypcja staje przed wysłaniem, a dopasowanie idzie po identyfikatorze
 * żądania — dokładnie tak, jak koreluje kanał. Odpowiedź przychodzi wyłącznie
 * z ramki gniazda, więc między `naDowolny` a `wyslij` nie ma okna, w którym
 * mogłaby przepaść. Transport kolejkuje ramki do chwili otwarcia połączenia,
 * więc wywołanie przed nawiązaniem łączności czeka, zamiast przepaść.
 */
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

/** Czy koperta odpowiedzi jest odmową „rdzeń nie ma uchwytu" (`*.unknown`). */
export function bezUchwytu(odpowiedz: OdpowiedzBramki<unknown>): boolean {
  return !odpowiedz.udana && odpowiedz.typ.endsWith('.unknown');
}

/** Rodzaj metody kontraktu dla metody widzianej przez Operatora. */
function rodzajMetody(metoda: MetodaWejscia): AuthMethodKind {
  return metoda === 'pin' ? AuthMethodKind.Pin : AuthMethodKind.Password;
}

/** Buduje źródło czynności bramki nad kanałem rdzenia. */
export function utworzZrodloUwierzytelnienia(kanal: Kanal): ZrodloUwierzytelnienia {
  // ── nasłuch powitania ──────────────────────────────────────────────────────
  // Zapisane powitanie zostaje na później: rozpoznanie bramki bywa wołane po
  // raz drugi (przycisk „Spróbuj ponownie"), a koperta przechodzi raz.
  let uslyszane: Powitanie | null = null;
  const czekajacy: ((powitanie: Powitanie) => void)[] = [];
  let odepnijPowitanie: Odsubskrybuj = () => undefined;

  odepnijPowitanie = kanal.naDowolny((koperta) => {
    if (koperta.type !== Command.ConnectionHello || !czyOdpowiedz(koperta)) return;
    // Odmowa powitania to też milczenie, nie „nie": rdzeń, który powitania nie
    // umie, nie powiedział niczego o wymogu logowania ani o stanie bramki.
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

    /**
     * Sonda stanu bramki: `auth.login` metodą `password` bez sekretu.
     *
     * Kontrakt nie ma komendy `auth.status`, a ekran musi wiedzieć, czy pokazać
     * wejście, czy pierwsze ustawienie hasła. Zamiast zgadywać, pytamy rdzeń
     * i czytamy JEGO rozstrzygnięcie: adapter bramki najpierw szuka kotwicy,
     * a dopiero potem patrzy na sekret (`core/adapter_modul_auth.go`,
     * `WejdzPrzezBramke` → `metodaWejscia`), więc odmowa mówi wprost:
     *
     *   `not_found`          → „hasło bramki nie istnieje — bramki jeszcze nie
     *                          ustawiono (auth.register)" → `nieustawiona`,
     *   `validation_failed`  → „wejście metodą password bez sekretu" — kotwica
     *                          jest, zabrakło tylko sekretu → `ustawiona`.
     *
     * Sonda jest drogą zapasową, nie pierwszą: odpowiedź na to samo pytanie
     * niesie powitanie (`gatewayConfigured`), które nie wyprowadza stanu
     * z odmowy i nie przechodzi przez dławik prób. Sonda zostaje na wypadek
     * milczenia rdzenia — pole puste znaczy „nie wiem", a wtedy trzeba zapytać
     * inaczej, nie zgadnąć.
     *
     * Sonda niczego nie zmienia w rdzeniu i nie niesie żadnego sekretu.
     */
    async zbadaj() {
      const odpowiedz = await poslij(kanal, Command.AuthLogin, {
        method: AuthMethodKind.Password,
      });
      if (bezUchwytu(odpowiedz)) return { stan: 'bez-uchwytu', blad: odpowiedz.blad };
      if (odpowiedz.blad?.code === ErrorCode.NotFound) return { stan: 'nieustawiona' };
      if (odpowiedz.blad?.code === ErrorCode.ValidationFailed) return { stan: 'ustawiona' };
      // Odpowiedź udana jest tu niemożliwa (sonda nie niesie sekretu); każda
      // inna odmowa — np. rdzeń bez sejfu poświadczeń (`internal_error`) —
      // idzie do ekranu z powodem, nie jest zamiatana pod formularz.
      return { stan: 'niepewny', blad: odpowiedz.blad };
    },

    async zbadajPin(urzadzenie) {
      if (urzadzenie === '') return false;
      const odpowiedz = await poslij(kanal, Command.AuthLogin, {
        method: AuthMethodKind.Pin,
        deviceId: urzadzenie,
      });
      // Wyłącznie `validation_failed` znaczy „PIN tu jest". Każda inna odpowiedź
      // — `not_found`, brak uchwytu, rdzeń bez sejfu — zostawia segment PIN-u
      // schowany: metoda niepewna nie ma prawa stanąć na ekranie jako czynna.
      return !bezUchwytu(odpowiedz) && odpowiedz.blad?.code === ErrorCode.ValidationFailed;
    },

    // `deviceId` idzie tylko przy PIN-ie. Dla metody `password` kontrakt go nie
    // wymaga (hasło otwiera bramkę z każdej maszyny), a PIN jest właściwy
    // maszynie i bez niego rdzeń odmawia wprost. Tożsamość klienta z powitania
    // urządzeniem nie jest — nadaje się ją na czas uruchomienia
    // (`protokol/tozsamosc-klienta.ts`), więc PIN chodzi po zapisie trwałym
    // (`tozsamosc-urzadzenia.ts`). Puste pole zostaje puste.
    wejdz: (zadanie) =>
      poslij(kanal, Command.AuthLogin, {
        method: rodzajMetody(zadanie.metoda),
        secret: zadanie.sekret,
        keepSignedIn: zadanie.niewylogowuj,
        // Login idzie wyłącznie przy metodzie hasła — rdzeń rozpoznaje po nim
        // konto, a metody właściwe urządzeniu go nie czytają. Bez tego pola
        // ekran pokazywał login, którego nikt nie wysyłał.
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
