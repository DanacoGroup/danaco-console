import {
  Command,
  PROTOCOL_VERSION,
  type ErrorInfo,
  type Window,
  type WindowCreateRequest,
} from '../../../shared/contract';
import { utworzMagistrale, type Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import type { Kanal, Wynik } from './kanal';
import type { TozsamoscKlienta } from './tozsamosc-klienta';

/** Zamówienie okna komunikacji: treść `window.create` bez sesji, którą zna kanał. */
export type ZamowienieOkna = Omit<WindowCreateRequest, 'sessionId'>;

/** Kolejne etapy uzgodnienia z rdzeniem. */
export type EtapUzgodnienia = 'powitanie' | 'sesja' | 'okno' | 'gotowe';

/** Wynik pojedynczego etapu — podstawa komunikatu dla operatora. */
export interface PostepUzgodnienia {
  etap: EtapUzgodnienia;
  udany: boolean;
  blad?: ErrorInfo;
}

/**
 * Uzgodnienie: powitanie → założenie sesji → otwarcie okna komunikacji.
 *
 * Kolejność wynika z kontraktu (`connection.hello`, `session.create`,
 * `window.create`): okno jest bytem pośrednim między sesją a wiadomością, więc
 * wiadomość można wysłać dopiero po jego otwarciu. Niepowodzenie etapu nie
 * blokuje połączenia ani kolejnych prób — dotyczy wyłącznie bieżącego
 * wywołania.
 */
export interface Uzgodnienie {
  /** Rozpoczyna uzgodnienie od powitania połączenia. */
  rozpocznij(): void;
  /** Okno otwarte przez rdzeń albo `null`, dopóki nie powstało. */
  okno(): Window | null;
  /** Subskrypcja postępu uzgodnienia. */
  naPostep(sluchacz: (postep: PostepUzgodnienia) => void): Odsubskrybuj;
  /** Subskrypcja otwarcia okna komunikacji. */
  naOtwarcieOkna(sluchacz: (okno: Window) => void): Odsubskrybuj;
  /**
   * Wiąże żywe połączenie ze świeżo założoną sesją bramki.
   *
   * Uzgodnienie wita rdzeń w chwili nawiązania połączenia, czyli zanim podano
   * hasło, więc niesie wtedy token sesji poprzedniej albo żaden. Po wejściu
   * przez bramkę rdzeń musi dowiedzieć się, które gniazdo należy teraz do
   * której sesji; inaczej `auth.password.reset` rozłączyłby tego, kto właśnie
   * zmienił hasło.
   *
   * Powitanie jest odczytem — nie zakłada ani sesji pracy, ani okna — więc
   * powtórzenie niczego nie dubluje. Ciągu dalszego uzgodnienia to wywołanie
   * nie uruchamia; tamten idzie przy nawiązaniu połączenia.
   */
  zwiazSesjeBramki(token: string): void;
  /**
   * Tożsamość klienta przedstawiona rdzeniowi w powitaniu.
   *
   * Pole `id` jest tym samym `clientId`, którego żądają `home.enter`,
   * `environment.enter`, `session.focus` i `session.bind`. Widok bierze je
   * stąd, zamiast wołać `tozsamoscKlienta()` po raz drugi: każde wywołanie
   * nadaje identyfikator nowy, a ognisko jest właściwością klienta — drugi
   * identyfikator rozdzieliłby ognisko od połączenia, które je zgłosiło.
   */
  klient: TozsamoscKlienta;
}

/**
 * Źródło tokenu sesji bramki dla powitania.
 *
 * Token dostarcza funkcja, nie wartość, ponieważ powitanie idzie przy każdym
 * nawiązaniu połączenia — także po zerwaniu i ponownym połączeniu — a sesja
 * bramki może się między nimi zmienić (wejście, wylogowanie, wygaśnięcie).
 * Wartość zamrożona przy składaniu warstwy niosłaby stan sprzed uruchomienia
 * i wiązałaby połączenie z sesją, której już nie ma.
 *
 * Warstwa protokołu nie zna pochodzenia tokenu: magazyn sesji należy do
 * `uwierzytelnienie/`, stąd wstrzyknięcie od składającego.
 */
export type ZrodloTokenuBramki = () => string | undefined;

export function utworzUzgodnienie(
  kanal: Kanal,
  zamowienie: ZamowienieOkna,
  klient: TozsamoscKlienta,
  tokenBramki?: ZrodloTokenuBramki,
): Uzgodnienie {
  const postepy = utworzMagistrale<PostepUzgodnienia>();
  const otwarcia = utworzMagistrale<Window>();
  let otwarte: Window | null = null;

  /** Ogłasza wynik etapu i rozstrzyga, czy uzgodnienie idzie dalej. */
  function zglos(etap: EtapUzgodnienia, wynik: Wynik<unknown>): boolean {
    postepy.oglos({ etap, udany: wynik.udany, blad: wynik.blad });
    return wynik.udany;
  }

  function powitaj(): void {
    kanal.wyslij(
      Command.ConnectionHello,
      {
        clientId: klient.id,
        clientVersion: klient.wersja,
        protocolVersion: PROTOCOL_VERSION,
        // Token wiąże to połączenie z sesją bramki. Jego brak nie wstrzymuje
        // powitania i nie jest błędem — rdzeń odpowiada `authenticated: false`,
        // a ekran logowania i tak stoi nad aplikacją.
        token: tokenBramki?.(),
      },
      (wynik) => {
        if (zglos('powitanie', wynik)) zalozSesje();
      },
    );
  }

  function zalozSesje(): void {
    kanal.wyslij(Command.SessionCreate, { title: zamowienie.title }, (wynik) => {
      const sesja = wynik.wynik?.session;
      if (!zglos('sesja', wynik) || sesja === undefined) return;
      kanal.sesja().ustaw(sesja.id);
      otworzOkno();
    });
  }

  function otworzOkno(): void {
    kanal.wyslij(
      Command.WindowCreate,
      { ...zamowienie, sessionId: kanal.sesja().id() },
      (wynik) => {
        const okno = wynik.wynik?.window;
        if (!zglos('okno', wynik) || okno === undefined) return;
        otwarte = okno;
        otwarcia.oglos(okno);
        postepy.oglos({ etap: 'gotowe', udany: true });
      },
    );
  }

  return {
    rozpocznij: powitaj,
    zwiazSesjeBramki: (token) => {
      kanal.wyslij(Command.ConnectionHello, {
        clientId: klient.id,
        clientVersion: klient.wersja,
        protocolVersion: PROTOCOL_VERSION,
        token,
      });
    },
    okno: () => otwarte,
    naPostep: (sluchacz) => postepy.subskrybuj(sluchacz),
    naOtwarcieOkna: (sluchacz) => otwarcia.subskrybuj(sluchacz),
    klient,
  };
}
