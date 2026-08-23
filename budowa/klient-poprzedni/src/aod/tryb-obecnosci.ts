import { PROGI_AOD, type WyciszenieCzasowe } from './progi-aod';
import { WagaUjawnienia, czyUjawniaSieSamoczynnie } from './rodzaje-sugestii';
import {
  magazynWyciszenDomyslny,
  utworzStanWyciszen,
  wyciszenieObejmujace,
  zdanieWyciszenia,
  type MagazynWyciszen,
  type OpisSugestiiWobecWyciszenia,
  type StanWyciszen,
  type WyciszenieCzasem,
} from './wyciszenie-aod';

/**
 * Tryby obecności i reguła ujawniania funkcji Always On Display — rozdz. 9.2
 * i 3.5 opracowania `docs/funkcje-globalne/always-on-display.md`.
 *
 * Plik nie dotyka dokumentu i nie zna kanału. Trzyma tryb obecności, sięga po
 * wykaz wyciszeń do `wyciszenie-aod.ts` i rozstrzyga z obu trzy rzeczy, o które
 * pyta warstwa widoku:
 *   • czy awatar jest w polu widzenia (tryb ukryty go zabiera),
 *   • czy plakietka liczbowa się pokazuje (wyciszenie czasowe ją chowa),
 *   • czy ta sugestia otwiera dymek TERAZ (waga, tryb, trzy rodzaje wyciszenia,
 *     limit godzinowy i odstęp między dymkami razem).
 *
 * WYJĄTEK WAGI KRYTYCZNEJ (rozdz. 3.5, wiersz ostatni) jest zaszyty w regule,
 * a nie zostawiony wołającemu, i przechodzi przez WSZYSTKIE rodzaje wyciszenia:
 * czasowe, kontekstowe, klasy zdarzeń oraz tryb cichy. Punkt decyzyjny pętli
 * wykonawczej wstrzymujący proces ujawnia się mimo wyciszenia — plakietką, bez
 * dymka — a tryb cichy zachowuje ten wyjątek bez syntezy mowy. Osłabienie tego
 * wyjątku byłoby jedyną ciszą, po której Operator nie dowiaduje się, że praca
 * stoi, więc reguła sprawdza go PIERWSZY, przed wszystkim innym.
 *
 * GDZIE TEN STAN MIESZKA. Opracowanie (rozdz. 10.1) chce tych ustawień
 * w zasięgu „Always On Display" okna konfiguracji, ze zmianą obowiązującą
 * natychmiast na wszystkich urządzeniach Operatora i rozgłoszeniem
 * `config.changed`. Kontrakt takiej kategorii ustawień nie ma i nie ma komendy
 * zapisującej tryb obecności ani wyciszenie nakładki. Stan stoi więc na
 * stanowisku Operatora, a magazyn jest PODAWANY — patrz `wyciszenie-aod.ts`
 * i `wyciszenie-braki-kontraktu.ts`.
 */

/** Tryb obecności funkcji — rozdz. 9.2. */
export const TrybObecnosci = {
  /** Awatar widoczny, sugestie ujawniane zgodnie z progami, synteza mowy czynna. */
  Pelny: 'pelny',
  /** Awatar widoczny, sugestie gromadzone bez dymka, synteza mowy wyłączona. */
  Cichy: 'cichy',
  /** Awatar poza polem widzenia, funkcja czynna w tle. */
  Ukryty: 'ukryty',
} as const;
export type TrybObecnosci = (typeof TrybObecnosci)[keyof typeof TrybObecnosci];

/** Nazwa trybu widziana przez Operatora w przełączniku nagłówka. */
export const NAZWY_TRYBOW: Readonly<Record<TrybObecnosci, string>> = {
  [TrybObecnosci.Pelny]: 'pełny',
  [TrybObecnosci.Cichy]: 'cichy',
  [TrybObecnosci.Ukryty]: 'ukryty',
};

/** Zachowanie trybu opisane zdaniem — rozdz. 9.2, kolumna „Zachowanie". */
export const OPISY_TRYBOW: Readonly<Record<TrybObecnosci, string>> = {
  [TrybObecnosci.Pelny]: 'Awatar widoczny, sugestie ujawniane zgodnie z progami, synteza mowy czynna.',
  [TrybObecnosci.Cichy]:
    'Awatar widoczny, sugestie gromadzone bez dymka, synteza mowy wyłączona, wyjątek wagi ' +
    'krytycznej zachowany.',
  [TrybObecnosci.Ukryty]:
    'Awatar poza polem widzenia, funkcja czynna w tle; dostęp skrótem klawiszowym i listwą ustawień.',
};

/**
 * Sugestia opisana tym, co reguła ujawniania musi o niej wiedzieć.
 *
 * Klasa zdarzenia, moduł i karta sesji przychodzą z `OpisSugestiiWobecWyciszenia`
 * i są opcjonalne: sugestia bez modułu nie wpada w wyciszenie modułu, bo cisza
 * bez podstawy jest gorsza od ujawnienia.
 */
export interface OpisUjawnienia extends OpisSugestiiWobecWyciszenia {
  /** Waga ujawnienia — rozdz. 4.5. */
  waga: WagaUjawnienia;
  /** Punkt decyzyjny wstrzymujący proces — wyjątek wagi krytycznej rozdz. 3.5. */
  krytyczna?: boolean;
}

/** Nastawy stanu obecności — magazyn i wykaz wyciszeń są podawane, nie brane na sztywno. */
export interface OpisStanuObecnosci {
  /**
   * Magazyn zapisu trybu obecności i wyciszeń.
   *
   * Pominięty znaczy zapis miejscowy przeglądarki, `null` znaczy „bez zapisu".
   * Gdy kontrakt poniesie wyciszenie nakładki jako byt rdzenia, magazyn zmieni
   * się w jednym wywołaniu, bez zmiany ani jednego wołacza.
   */
  magazyn?: MagazynWyciszen | null;
  /** Gotowy wykaz wyciszeń; pominięty buduje się nad {@link OpisStanuObecnosci.magazyn}. */
  wyciszenia?: StanWyciszen;
}

/** Stan obecności widziany przez warstwę widoku. */
export interface StanObecnosci {
  tryb(): TrybObecnosci;
  ustawTryb(tryb: TrybObecnosci): void;
  /** Przełącza między trybem ukrytym a poprzednim — skrót `Ctrl/Cmd + Shift + H`. */
  przelaczUkrycie(): void;
  /** Przełącza tryb cichy — jedno kliknięcie w menu wyciszania, bez potwierdzania. */
  przelaczTrybCichy(): void;

  /** Wykaz wyciszeń czynnych wraz z czynnościami, które go zmieniają. */
  wyciszenia: StanWyciszen;

  /** Trwające wyciszenie czasowe albo `null`. Przeterminowane znosi się samo. */
  wyciszenie(teraz: number): WyciszenieCzasem | null;
  /** Włącza wyciszenie czasowe wskazanej wartości. */
  wycisz(czas: WyciszenieCzasowe, teraz: number): void;
  /** Znosi wyciszenie czasowe — jedno kliknięcie, bez pytania o potwierdzenie. */
  zniesWyciszenie(): void;
  /** Przełącza wyciszenie kwadransowe — skrót `Ctrl/Cmd + Shift + M`. */
  przelaczWyciszenieKwadransem(teraz: number): void;
  /** Czy funkcję obejmuje dziś jakiekolwiek wyciszenie — stan „Wyciszony" rozdz. 9.1. */
  czyJakiekolwiekWyciszenie(teraz: number): boolean;

  /** Czy awatar jest w polu widzenia. */
  czyAwatarWidoczny(): boolean;
  /** Czy plakietka liczbowa się pokazuje (wyciszenie czasowe ją chowa). */
  czyPlakietkaWidoczna(teraz: number): boolean;
  /** Czy synteza mowy jest czynna (tryb cichy ją wyłącza). */
  czySyntezaMowy(): boolean;

  /** Czy ta sugestia otwiera dymek w tej chwili. */
  czyOtworzycDymek(opis: OpisUjawnienia, teraz: number): boolean;
  /** Odnotowuje otwarcie dymka — podstawa limitu godzinowego i odstępu. */
  odnotujDymek(teraz: number): void;
  /** Zdanie mówiące, dlaczego dymek się teraz nie otworzy; pusty napis, gdy się otworzy. */
  powodMilczenia(opis: OpisUjawnienia, teraz: number): string;
  /**
   * Czy ta sugestia jest wstrzymana wyciszeniem kontekstowym albo klasy zdarzeń.
   *
   * Wyciszenie czasowe wstrzymuje wszystko, więc nie liczy się tu jako wstrzymanie
   * wybiórcze: plakietka i tak jest wtedy schowana. Sugestia krytyczna nie jest
   * wstrzymana nigdy — ujawnia się plakietką mimo każdego wyciszenia.
   */
  czySugestiaWstrzymana(opis: OpisUjawnienia, teraz: number): boolean;

  /** Zgłasza obserwatora zmiany stanu; zwraca odsubskrybowanie. */
  obserwuj(sluchacz: () => void): () => void;
}

/** Klucz zapisu trybu obecności. Jeden na stanowisko — funkcja jest globalna. */
const KLUCZ_ZAPISU = 'danaco.aod.obecnosc';

/** Kształt zapisu miejscowego — czytany defensywnie, bo zapis bywa cudzy albo stary. */
interface ZapisObecnosci {
  tryb?: string;
}

export function utworzStanObecnosci(opis: OpisStanuObecnosci = {}): StanObecnosci {
  const magazyn = opis.magazyn === undefined ? magazynWyciszenDomyslny() : opis.magazyn;
  const wyciszenia = opis.wyciszenia ?? utworzStanWyciszen(magazyn);

  let tryb: TrybObecnosci = wczytajTryb(magazyn);

  /** Tryb sprzed ukrycia — żeby `Ctrl/Cmd + Shift + H` wracał tam, skąd wyszedł. */
  let trybPrzedUkryciem: TrybObecnosci = tryb === TrybObecnosci.Ukryty ? TrybObecnosci.Pelny : tryb;

  /** Chwile otwarcia dymków — podstawa limitu godzinowego i odstępu (rozdz. 3.4). */
  const dymki: number[] = [];

  const sluchacze = new Set<() => void>();

  function rozglos(): void {
    for (const sluchacz of sluchacze) sluchacz();
  }

  // Zmiana wykazu wyciszeń jest zmianą stanu obecności: awatar przechodzi w stan
  // „Wyciszony" i wraca z niego bez osobnego odświeżenia po stronie widoku.
  wyciszenia.obserwuj(rozglos);

  function zapiszTryb(): void {
    if (magazyn === null) return;
    try {
      magazyn.setItem(KLUCZ_ZAPISU, JSON.stringify({ tryb } satisfies ZapisObecnosci));
    } catch {
      // Zapis miejscowy bywa wyłączony ustawieniem przeglądarki. Funkcja działa
      // dalej, tracąc wyłącznie pamięć między przeładowaniami strony.
    }
  }

  /** Ile dymków otwarto w ostatniej godzinie; przy okazji czyści starsze wpisy. */
  function dymkiWOstatniejGodzinie(teraz: number): number {
    const granica = teraz - 60 * 60_000;
    while (dymki.length > 0 && (dymki[0] ?? 0) < granica) dymki.shift();
    return dymki.length;
  }

  function ostatniDymek(): number | undefined {
    return dymki[dymki.length - 1];
  }

  function powodMilczenia(opis: OpisUjawnienia, teraz: number): string {
    const obejmujace = wyciszenieObejmujace(wyciszenia.czynne(teraz), opis);

    // WYJĄTEK WAGI KRYTYCZNEJ — sprawdzany pierwszy i obejmujący wszystkie
    // rodzaje wyciszenia (rozdz. 3.5, wiersz „Wyjątek wagi krytycznej").
    // Sugestia ujawnia się mimo wyciszenia, ale PLAKIETKĄ, nie dymkiem.
    if (opis.krytyczna === true && (tryb === TrybObecnosci.Cichy || obejmujace !== null)) {
      const czego =
        obejmujace === null
          ? 'trybu cichego'
          : `wyciszenia — ${zdanieWyciszenia(obejmujace)}`;
      return (
        `Punkt decyzyjny wstrzymujący proces ujawnia się mimo ${czego} — plakietką, bez dymka ` +
        '(rozdz. 3.5, wyjątek wagi krytycznej).'
      );
    }

    if (tryb === TrybObecnosci.Ukryty) {
      return 'Tryb ukryty: awatar jest poza polem widzenia, więc dymek nie ma się przy czym otworzyć.';
    }

    if (tryb === TrybObecnosci.Cichy) {
      return 'Tryb cichy: sugestie gromadzą się w liście oczekujących, dymek się nie otwiera.';
    }

    if (obejmujace !== null) return zdanieWyciszenia(obejmujace);

    if (!czyUjawniaSieSamoczynnie(opis.waga)) {
      return 'Samoczynnie ujawnia się wyłącznie waga wysoka; ta sugestia czeka pod plakietką.';
    }

    if (dymkiWOstatniejGodzinie(teraz) >= PROGI_AOD.liczbaDymkowNaGodzine) {
      return `Limit ${PROGI_AOD.liczbaDymkowNaGodzine} dymków na godzinę wyczerpany — sugestia trafia do listy oczekujących.`;
    }

    const ostatni = ostatniDymek();
    if (ostatni !== undefined && teraz - ostatni < PROGI_AOD.odstepMiedzyDymkamiMs) {
      return 'Odstęp między dymkami krótszy niż 5 minut — sugestia łączy się z poprzednią w pozycję zbiorczą.';
    }

    return '';
  }

  return {
    tryb: () => tryb,

    ustawTryb(nowy) {
      if (nowy === tryb) return;
      if (tryb !== TrybObecnosci.Ukryty) trybPrzedUkryciem = tryb;
      tryb = nowy;
      zapiszTryb();
      rozglos();
    },

    przelaczUkrycie() {
      if (tryb === TrybObecnosci.Ukryty) {
        tryb = trybPrzedUkryciem;
      } else {
        trybPrzedUkryciem = tryb;
        tryb = TrybObecnosci.Ukryty;
      }
      zapiszTryb();
      rozglos();
    },

    przelaczTrybCichy() {
      if (tryb === TrybObecnosci.Cichy) {
        tryb = TrybObecnosci.Pelny;
      } else {
        if (tryb !== TrybObecnosci.Ukryty) trybPrzedUkryciem = tryb;
        tryb = TrybObecnosci.Cichy;
      }
      zapiszTryb();
      rozglos();
    },

    wyciszenia,

    wyciszenie: (teraz) => wyciszenia.czasowe(teraz),

    wycisz(czas, teraz) {
      wyciszenia.wyciszCzasem(czas, teraz);
    },

    zniesWyciszenie() {
      wyciszenia.znies('czasowe');
    },

    przelaczWyciszenieKwadransem(teraz) {
      wyciszenia.przelaczCzas('kwadrans', teraz);
    },

    czyJakiekolwiekWyciszenie: (teraz) => wyciszenia.czynne(teraz).length > 0,

    czyAwatarWidoczny: () => tryb !== TrybObecnosci.Ukryty,

    czyPlakietkaWidoczna(teraz) {
      if (tryb === TrybObecnosci.Ukryty) return false;
      // Plakietkę chowa wyłącznie wyciszenie czasowe (rozdz. 3.5, wiersz
      // pierwszy: „plakietka pozostaje ukryta"). Wyciszenie kontekstowe i klasy
      // zdarzeń wstrzymują część sugestii, więc plakietka pozostaje prawdziwa
      // dla pozostałych — chowanie jej zabrałoby Operatorowi wiedzę o nich.
      return wyciszenia.czasowe(teraz) === null;
    },

    czySyntezaMowy: () => tryb === TrybObecnosci.Pelny,

    czyOtworzycDymek: (opis, teraz) => powodMilczenia(opis, teraz) === '',

    odnotujDymek(teraz) {
      dymki.push(teraz);
      dymkiWOstatniejGodzinie(teraz);
    },

    powodMilczenia,

    czySugestiaWstrzymana(opis, teraz) {
      if (opis.krytyczna === true) return false;
      const obejmujace = wyciszenieObejmujace(wyciszenia.czynne(teraz), opis);
      return obejmujace !== null && obejmujace.rodzaj !== 'czasowe';
    },

    obserwuj(sluchacz) {
      sluchacze.add(sluchacz);
      return () => sluchacze.delete(sluchacz);
    },
  };
}

/**
 * Odczyt trybu obecności z magazynu.
 *
 * Zapis bywa nieobecny, cudzy albo starszy od tego pliku — każda niezgodność
 * daje stan domyślny opracowania (tryb pełny), a nie wyjątek.
 */
function wczytajTryb(magazyn: MagazynWyciszen | null): TrybObecnosci {
  if (magazyn === null) return TrybObecnosci.Pelny;

  let zapis: ZapisObecnosci;
  try {
    const surowy = magazyn.getItem(KLUCZ_ZAPISU);
    if (surowy === null || surowy === '') return TrybObecnosci.Pelny;
    zapis = JSON.parse(surowy) as ZapisObecnosci;
  } catch {
    return TrybObecnosci.Pelny;
  }

  return czyTryb(zapis.tryb) ? zapis.tryb : TrybObecnosci.Pelny;
}

function czyTryb(wartosc: unknown): wartosc is TrybObecnosci {
  return (
    typeof wartosc === 'string' && (Object.values(TrybObecnosci) as string[]).includes(wartosc)
  );
}
