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

/** Zamówienie okna komunikacji: treść żądania otwarcia okna, ale bez pola sesji, którą już zna ten kanał. */
export type ZamowienieOkna = Omit<WindowCreateRequest, 'sessionId'>;

/** Kolejne etapy uzgodnienia z rdzeniem: powitanie, założenie sesji, otwarcie okna, wreszcie stan gotowy. */
export type EtapUzgodnienia = 'powitanie' | 'sesja' | 'okno' | 'gotowe';

/** Wynik pojedynczego etapu uzgodnienia — podstawa komunikatu dla operatora o postępie tego połączenia. */
export interface PostepUzgodnienia {
  etap: EtapUzgodnienia;
  udany: boolean;
  blad?: ErrorInfo;
}

/** Uzgodnienie: powitanie, założenie sesji, otwarcie okna komunikacji, w kolejności wynikającej z kontraktu. */
export interface Uzgodnienie {
  /** Rozpoczyna uzgodnienie od powitania połączenia. */
  rozpocznij(): void;
  /** Okno otwarte przez rdzeń albo `null`, dopóki nie powstało. */
  okno(): Window | null;
  /** Subskrypcja postępu uzgodnienia. */
  naPostep(sluchacz: (postep: PostepUzgodnienia) => void): Odsubskrybuj;
  /** Subskrypcja otwarcia okna komunikacji. */
  naOtwarcieOkna(sluchacz: (okno: Window) => void): Odsubskrybuj;
  // Wiąże żywe połączenie ze świeżo założoną sesją bramki, po wejściu przez bramkę logowania.
  zwiazSesjeBramki(token: string): void;
  // Tożsamość klienta przedstawiona rdzeniowi w powitaniu, ten sam identyfikator co w innych żądaniach.
  klient: TozsamoscKlienta;
}

/** Źródło tokenu sesji bramki dla powitania, dostarczane funkcją, bo sesja bramki może się między nimi zmienić. */
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
        // Token wiąże to połączenie z sesją bramki; jego brak nie wstrzymuje powitania i nie jest błędem.
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
