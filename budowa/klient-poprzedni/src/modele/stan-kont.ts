import {
  ChangeKind,
  type Account,
  type AccountAddRequest,
  type AccountDefaultSetRequest,
  type AccountKind,
  type AccountUpdateRequest,
} from '../../../shared/contract';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import type { Kanal, Wynik } from '../protokol/kanal';
import { utworzZapisyKont } from './zapisy-kont';
import { uporzadkujKonta, utworzZrodloKont, type ZrodloKont } from './zrodlo-kont';

/**
 * Stan rejestru kont — jedno źródło prawdy dla wykazu, formularza i wyboru
 * bytu osi konta.
 *
 * Wykaz i wybór trzymamy razem, ponieważ każda zmiana rejestru dotyka obu:
 * usunięcie konta czynnego musi przestawić wybór, a nie zostawić formularz
 * wskazujący byt, którego już nie ma.
 *
 * Zdarzenie `account.changed` przychodzi także wtedy, gdy konto założono na innym
 * urządzeniu; stan przyjmuje je tą samą drogą, co własną odpowiedź.
 *
 * Żadna ścieżka nie zatrzymuje widoku: rdzeń, który nie odda wykazu, zostawia go
 * pustym, a zapis nieudany wraca jako `Wynik` z błędem, bez wyjątku.
 */
export type FazaOdczytu = 'spoczynek' | 'odczyt' | 'gotowe' | 'blad';

export interface StanKont {
  /**
   * Faza odczytu rejestru kont.
   *
   * Bez niej pusty rejestr znaczy trzy rzeczy naraz: „jeszcze nie pytałem",
   * „pytam" i „rdzeń nie zna ani jednego konta". Każdej należy się inny stan
   * widoku: nic, wskaźnik odczytu, stan pusty.
   */
  faza(): FazaOdczytu;
  /** Powód ostatniego niepowodzenia odczytu; pusty, gdy odczyt się powiódł. */
  powodNiepowodzenia(): string;
  /** Konta w kolejności wykazu. */
  konta(): readonly Account[];
  /** Konto czynne; `null`, gdy rejestr jest pusty albo nic nie wybrano. */
  wybrane(): Account | null;
  /** Ustawia konto czynne; identyfikator spoza wykazu zdejmuje wybór. */
  wybierz(identyfikator: string | null): void;
  /** Ograniczenie wykazu do rodzaju; `null` znaczy wszystkie rodzaje. */
  rodzaj(): AccountKind | null;
  /** Zmienia ograniczenie rodzaju i wczytuje wykaz ponownie. */
  ustawRodzaj(rodzaj: AccountKind | null): Promise<void>;
  /** Wczytuje wykaz kont z rdzenia. */
  odswiez(): Promise<void>;
  /** `account.add` — zakłada konto i czyni je czynnym. */
  dodaj(zadanie: AccountAddRequest): Promise<Wynik<unknown>>;
  /** `account.update` — zmienia konto wskazane w żądaniu. */
  zmien(zadanie: AccountUpdateRequest): Promise<Wynik<unknown>>;
  /** `account.remove` — usuwa konto i zdejmuje jego wybór. */
  usun(identyfikator: string): Promise<Wynik<unknown>>;
  /** `account.default.set` — wskazuje konto domyślne swojego rodzaju. */
  ustawDomyslne(zadanie: AccountDefaultSetRequest): Promise<Wynik<unknown>>;
  /** Subskrypcja przeliczenia stanu. */
  naZmiane(sluchacz: () => void): void;
  /** Odłącza subskrypcję zdarzeń kanału. */
  rozlacz(): void;
}

export function utworzStanKont(kanal: Kanal): StanKont {
  let powody: string[] = [];
  const zrodlo: ZrodloKont = utworzZrodloKont(kanal, (powod) => void powody.push(powod));

  let konta: Account[] = [];
  let wybrany: string | null = null;
  let rodzaj: AccountKind | null = null;
  let faza: FazaOdczytu = 'spoczynek';

  const sluchacze: Array<() => void> = [];
  const oglos = (): void => {
    for (const sluchacz of [...sluchacze]) sluchacz();
  };

  /** Konto potwierdzone przez rdzeń zastępuje wiersz o tym samym identyfikatorze. */
  const przyjmij = (konto: Account): void => {
    konta = uporzadkujKonta([
      ...konta.filter((istniejace) => istniejace.id !== konto.id),
      konto,
    ]);
    oglos();
  };

  const odlacz = (identyfikator: string): void => {
    konta = konta.filter((istniejace) => istniejace.id !== identyfikator);
    if (wybrany === identyfikator) wybrany = null;
    oglos();
  };

  const odsubskrybuj: Odsubskrybuj = zrodlo.naZmiane((tresc) => {
    if (tresc.change === ChangeKind.Deleted) odlacz(tresc.account.id);
    else przyjmij(tresc.account);
  });

  const zapisy = utworzZapisyKont({
    zrodlo,
    przyjmij,
    odlacz,
    wybierz: (identyfikator) => {
      wybrany = identyfikator;
      oglos();
    },
    wczytaj: () => wczytaj(),
  });

  async function wczytaj(): Promise<void> {
    powody = [];
    faza = 'odczyt';
    oglos();
    konta = await zrodlo.lista(rodzaj === null ? {} : { kind: rodzaj });
    if (wybrany !== null && !konta.some((konto) => konto.id === wybrany)) wybrany = null;
    faza = powody.length > 0 ? 'blad' : 'gotowe';
    oglos();
  }

  return {
    faza: () => faza,

    powodNiepowodzenia: () => powody.join(' '),

    konta: () => konta,

    wybrane: () => konta.find((konto) => konto.id === wybrany) ?? null,

    wybierz(identyfikator) {
      wybrany = identyfikator;
      oglos();
    },

    rodzaj: () => rodzaj,

    async ustawRodzaj(nowy) {
      rodzaj = nowy;
      await wczytaj();
    },

    odswiez: wczytaj,

    dodaj: zapisy.dodaj,
    zmien: zapisy.zmien,
    usun: zapisy.usun,
    ustawDomyslne: zapisy.ustawDomyslne,

    naZmiane: (sluchacz) => void sluchacze.push(sluchacz),

    rozlacz: odsubskrybuj,
  };
}
