import {
  ChangeKind,
  type AccessPoint,
  type Extension,
  type ExtensionKind,
} from '../../../../shared/contract';
import { opisOdmowyBledu } from '../../komponenty/odmowa';
import type { Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import type { Kanal } from '../../protokol/kanal';
import {
  utworzZrodloPunktowDostepu,
  type ZrodloPunktowDostepu,
} from './zrodlo-punktow-dostepu';
import {
  utworzZrodloRozszerzenApps,
  type ZlecenieInstalacji,
  type ZrodloRozszerzenApps,
} from './zrodlo-rozszerzen-apps';
import {
  utworzZrodloDobudowyRozszerzen,
  type ZrodloDobudowyRozszerzen,
} from './zrodlo-rozszerzen-dobudowa';

/** Który z dwóch odczytów strony dystrybucji — klucz, pod którym stoi powód odmowy tego odczytu w oknie. */
export type RodzajOdczytuKatalogu = 'katalog' | 'punkty';

/**
 * Jedno źródło prawdy strony dystrybucji i konsumpcji modułu Apps: katalog
 * rozszerzeń i katalog punktów dostępu, wspólne dla sześciu okien tej
 * strony, żeby zmiana w jednym oknie była natychmiast widoczna w pozostałych.
 */
export interface StanRozszerzen {
  zrodlo: ZrodloRozszerzenApps;
  /** Dobudowa obszaru rozszerzeń; osobne źródło, bo obsługuje to, co robi się na pozycji już stojącej. */
  dobudowa: ZrodloDobudowyRozszerzen;
  punktyDostepu: ZrodloPunktowDostepu;
  /** Katalog w kolejności oddanej przez rdzeń. */
  rozszerzenia(): readonly Extension[];
  /** Pozycje katalogu jednego rodzaju albo wielu rodzajów naraz. */
  rozszerzeniaRodzaju(rodzaje: readonly ExtensionKind[]): readonly Extension[];
  /** Punkty dostępu w kolejności oddanej przez rdzeń. */
  punkty(): readonly AccessPoint[];
  /** Punkt dostępu o wskazanym identyfikatorze; `null`, gdy katalog go nie zna. */
  punkt(idPunktu: string): AccessPoint | null;
  /** Czy katalog wrócił już z rdzenia — „pusto" to nie to samo, co „nie pytano". */
  czyKatalogCzytany(): boolean;
  /** Czy katalog punktów dostępu wrócił już z rdzenia. */
  czyPunktyCzytane(): boolean;
  /** Powód odmowy ostatniego odczytu danego rodzaju; pusty, gdy odczyt się udał. */
  powodOdczytu(co: RodzajOdczytuKatalogu): string;
  /** Pozycja wskazana do wglądu w panelach bocznych; `null`, gdy żadnej nie wskazano. */
  wybrane(): Extension | null;
  /** Wskazuje pozycję panelom bocznym; pusty identyfikator zdejmuje wskazanie. */
  wybierz(idRozszerzenia: string): void;
  /** Zleca `extension.list` bez zawężenia i wchłania katalog. */
  odczytajKatalog(): Promise<void>;
  /** Zleca `access.point.list` i wchłania katalog punktów. */
  odczytajPunkty(): Promise<void>;
  /** Zapisuje pozycję potwierdzoną przez rdzeń po instalacji, zapisie albo przełączeniu. */
  wchlonPozycje(pozycja: Extension): void;
  /** Zdejmuje pozycję potwierdzoną odinstalowaniem. */
  zdejmijPozycje(idRozszerzenia: string): void;
  /** Zapisuje punkt dostępu po sprawdzeniu, żeby stan i czas szły do wszystkich okien. */
  wchlonPunkt(punkt: AccessPoint): void;
  obserwuj(sluchacz: () => void): Odsubskrybuj;
  rozlacz(): void;
}

/** Zlecenie instalacji przekazywane oknom; kształt zlecenia bierze wprost ze źródła danych rozszerzeń modułu. */
export type { ZlecenieInstalacji };

export function utworzStanRozszerzen(kanal: Kanal): StanRozszerzen {
  const zrodlo = utworzZrodloRozszerzenApps(kanal);
  const dobudowa = utworzZrodloDobudowyRozszerzen(kanal);
  const punktyDostepu = utworzZrodloPunktowDostepu(kanal);
  const sluchacze = new Set<() => void>();
  const powody: Record<RodzajOdczytuKatalogu, string> = { katalog: '', punkty: '' };

  let rozszerzenia: Extension[] = [];
  let punkty: AccessPoint[] = [];
  let katalogCzytany = false;
  let punktyCzytane = false;
  let wybrane = '';

  function oglos(): void {
    for (const sluchacz of sluchacze) sluchacz();
  }

  /** Wstawia pozycję na miejsce pozycji o tym samym identyfikatorze albo na koniec. */
  function wstaw(pozycja: Extension): void {
    const miejsce = rozszerzenia.findIndex((inna) => inna.id === pozycja.id);
    if (miejsce === -1) rozszerzenia = [...rozszerzenia, pozycja];
    else rozszerzenia = rozszerzenia.map((inna) => (inna.id === pozycja.id ? pozycja : inna));
  }

  const odsubskrybowania: Odsubskrybuj[] = [
    // Zmiana katalogu przychodzi niezależnie od tego, gdzie zaszła; rejestr jest wspólny platformie.
    zrodlo.naZmianeKatalogu((tresc) => {
      if (tresc.change === ChangeKind.Deleted) {
        rozszerzenia = rozszerzenia.filter((inna) => inna.id !== tresc.extension.id);
        if (wybrane === tresc.extension.id) wybrane = '';
      } else {
        wstaw(tresc.extension);
      }
      oglos();
    }),
  ];

  return {
    zrodlo,
    dobudowa,
    punktyDostepu,
    rozszerzenia: () => rozszerzenia,
    rozszerzeniaRodzaju: (rodzaje) =>
      rozszerzenia.filter((pozycja) => rodzaje.includes(pozycja.kind)),
    punkty: () => punkty,
    punkt: (idPunktu) => punkty.find((inny) => inny.id === idPunktu) ?? null,
    czyKatalogCzytany: () => katalogCzytany,
    czyPunktyCzytane: () => punktyCzytane,
    powodOdczytu: (co) => powody[co],
    wybrane: () => rozszerzenia.find((pozycja) => pozycja.id === wybrane) ?? null,

    wybierz(idRozszerzenia) {
      wybrane = idRozszerzenia;
      oglos();
    },

    async odczytajKatalog() {
      powody.katalog = '';
      const wynik = await zrodlo.katalog({});
      if (!wynik.udany || wynik.wynik === undefined) {
        powody.katalog = opisOdmowyBledu('Odczyt katalogu rozszerzeń', wynik.blad);
        oglos();
        return;
      }
      // Zastępuje, nie dokłada: odpowiedź rdzenia jest pełnym stanem rejestru w chwili odczytu.
      rozszerzenia = [...wynik.wynik.extensions];
      katalogCzytany = true;
      oglos();
    },

    async odczytajPunkty() {
      powody.punkty = '';
      const wynik = await punktyDostepu.punkty();
      if (!wynik.udany || wynik.wynik === undefined) {
        powody.punkty = opisOdmowyBledu('Odczyt punktów dostępu', wynik.blad);
        oglos();
        return;
      }
      punkty = [...wynik.wynik.points];
      punktyCzytane = true;
      oglos();
    },

    wchlonPozycje(pozycja) {
      wstaw(pozycja);
      oglos();
    },

    zdejmijPozycje(idRozszerzenia) {
      rozszerzenia = rozszerzenia.filter((inna) => inna.id !== idRozszerzenia);
      if (wybrane === idRozszerzenia) wybrane = '';
      oglos();
    },

    wchlonPunkt(punkt) {
      const miejsce = punkty.findIndex((inny) => inny.id === punkt.id);
      if (miejsce === -1) punkty = [...punkty, punkt];
      else punkty = punkty.map((inny) => (inny.id === punkt.id ? punkt : inny));
      oglos();
    },

    obserwuj(sluchacz) {
      sluchacze.add(sluchacz);
      return () => sluchacze.delete(sluchacz);
    },

    rozlacz() {
      for (const odsubskrybuj of odsubskrybowania.splice(0)) odsubskrybuj();
      sluchacze.clear();
    },
  };
}
