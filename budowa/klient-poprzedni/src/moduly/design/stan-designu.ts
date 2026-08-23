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
 * Jeden zbiór zasobów i jedno okno modułu na cały moduł Design.
 *
 * Jedna odpowiedzialność: rozdanie trzem oknom tej samej prawdy i powiadomienie
 * ich o każdej jej zmianie.
 *
 * Trzy okna patrzą na ten sam zbiór. Prompt Builder oddaje wynik generowania do
 * Assets Panel, Assets Panel oddaje zasób na kanwę Design Board. Gdyby każde
 * okno prowadziło własny wykaz, zasób wygenerowany w kreatorze nie pojawiłby
 * się w panelu, a kanwa układałaby warstwy z zasobów, których panel już nie ma.
 *
 * Zdarzenie `design.asset.changed` jest drugim źródłem odświeżenia: wciąga
 * zasób powstały gdziekolwiek — także po stronie rdzenia — dokładnie tak samo
 * jak własny odczyt. Odpytywania w pętli tu nie ma.
 */
export type { FazaZasobow } from './zapis-designu';

export interface StanDesignu {
  /** Źródło komend obszaru `design.*` — okna wołają je wprost. */
  zrodlo: ZrodloDesignu;
  /** Zaplecze: okno modułu, rejestry, przekazanie międzymodułowe, postęp. */
  zaplecze: ZrodloZaplecza;
  /**
   * Czuwanie nad czynnościami okien — jedno na moduł.
   *
   * Stoi tutaj, a nie w każdym oknie osobno, bo próba życia kanału jest
   * pytaniem o jedną wspólną drogę do rdzenia, nie o okno.
   */
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
  /**
   * Zdejmuje zasób po własnym usunięciu, bez czekania na zdarzenie.
   *
   * Idzie tą samą funkcją co gałąź `deleted` zdarzenia. Rdzeń rozgłasza
   * usunięcie i zdarzenie i tak przyjdzie, ale okno, które właśnie kazało zasób
   * usunąć, nie ma prawa pokazywać go dalej ani przez chwilę — a kolejność
   * ramek nie jest niczym zagwarantowana.
   */
  zdejmij(idZasobu: string): void;
  obserwuj(sluchacz: () => void): () => void;
  rozlacz(): void;
}

export function utworzStanDesignu(kanal: Kanal): StanDesignu {
  const zrodlo = utworzZrodloDesignu(kanal);
  const zaplecze = utworzZrodloZaplecza(kanal);
  const sluchacze = new Set<() => void>();
  const zapis: ZapisDesignu = pustyZapisDesignu();

  // Próbą życia kanału jest najtańszy odczyt obszaru: jeden zasób, warunki
  // bieżącego okna. Komenda ma uchwyt w rdzeniu i odpowiada w milisekundach,
  // a odmowa też jest odpowiedzią — dowodzi, że kanał żyje. Wynik nigdzie nie
  // wsiąka: próba niczego nie zapisuje w stanie i nie rusza wykazu.
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

  // Rozdział po rodzaju zmiany. Nadawcą rodzaju `deleted` jest komenda
  // `design.asset.remove`: usunięcie zlecone w drugim oknie albo w obcym
  // połączeniu zdejmuje tu zasób z wykazu.
  const odsubskrybuj = zrodlo.naZmianeZasobu((tresc) => {
    if (tresc.change === ChangeKind.Deleted) usunZasob(zapis, tresc.asset.id);
    else wchlonZasob(zapis, tresc.asset);
    oglos();
  });

  /**
   * Odczyt zasobów pod czuwaniem.
   *
   * Odczyt zlecony przed zerwaniem gniazda nie dostaje odpowiedzi nigdy, więc
   * okno zarządcy stałoby w „Odczyt zasobów w toku…" także po powrocie rdzenia.
   * Cisza kanału trafia więc do fazy błędu wraz z powodem, który mówi prawdę:
   * odczytu nie ma, a wykaz na ekranie jest sprzed zerwania.
   */
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
        // Zdanie o powrocie należy się wyłącznie oknu, które nadal stoi na
        // ciszy. Gdy odpowiedź zdążyła przyjść, prawdą jest ona, nie zapowiedź
        // jej braku.
        if (!zapis.bezRozstrzygniecia) return;
        zapis.powod = zdanie;
        oglos();
      },
    });
    if (wynik !== null) {
      przyjmijOdczytZasobow(zapis, wynik);
      return;
    }
    // Odpowiedź spóźniona jest nadal prawdą o rdzeniu i odczyt może ją przyjąć.
    // Żądanie złożone przy martwym rdzeniu czeka w kolejce wychodzącej, a rdzeń
    // odpowiada na nie po ponownym połączeniu. Odczyt jest powtarzalny i niczego
    // nie zmienia w rdzeniu, więc przyjęcie takiej odpowiedzi nie niesie ryzyka.
    // Czynności zmieniające stan tej drogi nie mają: tam spóźniona odpowiedź
    // jest tylko zdaniem, bo skutek i tak trzeba odczytać.
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
