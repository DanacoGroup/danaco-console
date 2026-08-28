import { ChangeKind, type Channel, type DesignAsset, type Module } from '../../../../shared/contract';
import type { Kanal } from '../../protokol/kanal';
import { utworzCzuwanieRdzenia, type CzuwanieRdzenia } from './czuwanie-rdzenia';
import {
  odczytajZaplecze,
  przyjmijOdczytZasobow,
  pustyZapisDesignu,
  ustalOknoModulu,
  usunZasob,
  wchlonZasob,
  type FazaZasobow,
  type ZapisDesignu,
} from './zapis-designu';
import { utworzZrodloDesignu, type ZapytanieZasobow, type ZrodloDesignu } from './zrodlo-designu';
import { utworzZrodloZaplecza, type ZrodloZaplecza } from './zrodlo-zaplecza';

/**
 * Jeden zbiór zasobów i jedno okno modułu na cały moduł Design, rozdane trzem oknom wraz
 * z powiadomieniem o zmianie.
 */
export type { FazaZasobow } from './zapis-designu';

export interface StanDesignu {
  /** Źródło komend obszaru `design.*` — okna wołają je wprost. */
  zrodlo: ZrodloDesignu;
  /** Zaplecze: okno modułu, rejestry, przekazanie międzymodułowe, postęp. */
  zaplecze: ZrodloZaplecza;
  /** Czuwanie nad czynnościami okien — jedno na moduł, bo próba życia pyta o wspólną drogę, nie o okno. */
  czuwanie: CzuwanieRdzenia;
  /** Okno modułu Design w bieżącej sesji; puste, gdy rdzeń go nie wskazał. */
  idOkna(): string;
  /** Zdanie o oknie modułu — z `window.state.get` albo o jego braku. */
  opisOkna(): string;
  zasoby(): readonly DesignAsset[];
  wybrany(): DesignAsset | null;
  /** Wybiera zasób; `null` zdejmuje wybór. */
  wybierz(idZasobu: string | null): void;
  /** Rejestr kanałów modelu — kandydaci na silnik generujący. */
  silniki(): readonly Channel[];
  /** Katalog modułów rdzenia — cele przekazania, bez modułu własnego. */
  celePrzekazania(): readonly Module[];
  faza(): FazaZasobow;
  /** Powód odmowy ostatniego odczytu; pusty, gdy odczyt się udał. */
  powod(): string;
  /** Czy powód mówi o braku rozstrzygnięcia (zerwane gniazdo), a nie o odmowie. */
  czyBezRozstrzygniecia(): boolean;
  /** Odczytuje zasoby z rdzenia; warunki pominięte zostają bez zmian. */
  odswiez(zapytanie?: Partial<ZapytanieZasobow>): Promise<void>;
  /** Ustala okno modułu i jego stan; woła `window.list` i `window.state.get`. */
  ustalOkno(idSesji: string): Promise<void>;
  /** Odczytuje rejestr kanałów modelu i katalog modułów rdzenia. */
  odswiezZaplecze(): Promise<void>;
  /** Wciąga zasób po własnej zmianie, bez czekania na zdarzenie. */
  wchlon(zasob: DesignAsset): void;
  /** Zdejmuje zasób po własnym usunięciu, bez czekania na zdarzenie o niezagwarantowanej kolejności. */
  zdejmij(idZasobu: string): void;
  obserwuj(sluchacz: () => void): () => void;
  rozlacz(): void;
}

export function utworzStanDesignu(kanal: Kanal): StanDesignu {
  const zrodlo = utworzZrodloDesignu(kanal);
  const zaplecze = utworzZrodloZaplecza(kanal);
  const sluchacze = new Set<() => void>();
  const zapis: ZapisDesignu = pustyZapisDesignu();

  // Próbą życia kanału jest najtańszy odczyt; odmowa też jest odpowiedzią, dowodzi, że kanał żyje.
  const czuwanie = utworzCzuwanieRdzenia(() =>
    zrodlo.zasoby({ ...zapis.warunki, idOkna: zapis.oknoModulu, granica: 1 }),
  );

  function oglos(): void {
    for (const sluchacz of [...sluchacze]) sluchacz();
  }

  /** Wykonuje odczyt i rozgłasza jego wynik — także wtedy, gdy jest odmową. */
  async function przezOdczyt(czynnosc: () => Promise<void>): Promise<void> {
    await czynnosc();
    oglos();
  }

  // Rozdział po rodzaju zmiany: usunięcie zlecone gdzie indziej zdejmuje tu zasób z wykazu.
  const odsubskrybuj = zrodlo.naZmianeZasobu((tresc) => {
    if (tresc.change === ChangeKind.Deleted) usunZasob(zapis, tresc.asset.id);
    else wchlonZasob(zapis, tresc.asset);
    oglos();
  });

  /** Odczyt zasobów pod czuwaniem; cisza kanału trafia do fazy błędu z powodem mówiącym prawdę o wykazie. */
  async function odczytajPodCzuwaniem(): Promise<void> {
    const wywolanie = zrodlo.zasoby(zapis.warunki);
    const wynik = await czuwanie.prowadz('odczyt zasobów', wywolanie, {
      cisza: (zdanie) => {
        zapis.faza = 'blad';
        zapis.powod = zdanie;
        zapis.bezRozstrzygniecia = true;
        oglos();
      },
      wToku: (zdanie) => {
        zapis.powod = zdanie;
        oglos();
      },
      powrot: (zdanie) => {
        // Zdanie o powrocie należy się wyłącznie oknu, które stoi na ciszy, nie temu, co już ma odpowiedź.
        if (!zapis.bezRozstrzygniecia) return;
        zapis.powod = zdanie;
        oglos();
      },
    });
    if (wynik !== null) {
      przyjmijOdczytZasobow(zapis, wynik);
      return;
    }
    // Odpowiedź spóźniona jest wciąż prawdą; powtarzalny odczyt nie niesie ryzyka jej przyjęcia.
    void wywolanie.then((spozniony) => {
      przyjmijOdczytZasobow(zapis, spozniony);
      oglos();
    });
  }

  return {
    zrodlo,
    zaplecze,
    czuwanie,
    idOkna: () => zapis.oknoModulu,
    opisOkna: () => zapis.zdanieOOknie,
    zasoby: () => zapis.zbior,
    wybrany: () => zapis.zbior.find((wpis) => wpis.id === zapis.wybor) ?? null,
    silniki: () => zapis.kanaly,
    celePrzekazania: () => zapis.cele,
    faza: () => zapis.faza,
    powod: () => zapis.powod,
    czyBezRozstrzygniecia: () => zapis.bezRozstrzygniecia,
    ustalOkno: (idSesji) => przezOdczyt(() => ustalOknoModulu(zapis, zaplecze, idSesji)),
    odswiezZaplecze: () => przezOdczyt(() => odczytajZaplecze(zapis, zaplecze)),

    wybierz(idZasobu) {
      if (zapis.wybor === idZasobu) return;
      zapis.wybor = idZasobu;
      oglos();
    },

    odswiez(zapytanie) {
      zapis.warunki = { ...zapis.warunki, ...zapytanie, idOkna: zapis.oknoModulu };
      zapis.faza = 'odczyt';
      zapis.powod = '';
      zapis.bezRozstrzygniecia = false;
      oglos();
      return przezOdczyt(odczytajPodCzuwaniem);
    },

    wchlon(zasob) {
      wchlonZasob(zapis, zasob);
      oglos();
    },

    zdejmij(idZasobu) {
      usunZasob(zapis, idZasobu);
      oglos();
    },

    obserwuj(sluchacz) {
      sluchacze.add(sluchacz);
      return () => sluchacze.delete(sluchacz);
    },

    rozlacz() {
      odsubskrybuj();
      sluchacze.clear();
    },
  };
}
