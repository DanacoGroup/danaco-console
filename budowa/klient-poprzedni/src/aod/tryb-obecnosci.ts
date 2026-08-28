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

// Plik trzyma tryb obecności funkcji Always On Display i regułę samoczynnego ujawniania dymka.

/** Stała wylicza trzy tryby obecności funkcji Always On Display: pełny, cichy i ukryty, każdy z odmiennym zachowaniem awatara. */
export const TrybObecnosci = {
  /** Awatar widoczny, sugestie ujawniane zgodnie z progami, synteza mowy czynna. */
  Pelny: 'pelny',
  /** Awatar widoczny, sugestie gromadzone bez dymka, synteza mowy wyłączona. */
  Cichy: 'cichy',
  /** Awatar poza polem widzenia, funkcja czynna w tle. */
  Ukryty: 'ukryty',
} as const;
export type TrybObecnosci = (typeof TrybObecnosci)[keyof typeof TrybObecnosci];

/** Stała podaje nazwę każdego trybu obecności, widoczną dla Operatora w przełączniku trybu w nagłówku okna. */
export const NAZWY_TRYBOW: Readonly<Record<TrybObecnosci, string>> = {
  [TrybObecnosci.Pelny]: 'pełny',
  [TrybObecnosci.Cichy]: 'cichy',
  [TrybObecnosci.Ukryty]: 'ukryty',
};

/** Stała opisuje zachowanie każdego trybu obecności jednym zdaniem, widocznym w kolumnie objaśnień przełącznika trybu. */
export const OPISY_TRYBOW: Readonly<Record<TrybObecnosci, string>> = {
  [TrybObecnosci.Pelny]: 'Awatar widoczny, sugestie ujawniane zgodnie z progami, synteza mowy czynna.',
  [TrybObecnosci.Cichy]:
    'Awatar widoczny, sugestie gromadzone bez dymka, synteza mowy wyłączona, wyjątek wagi ' +
    'krytycznej zachowany.',
  [TrybObecnosci.Ukryty]:
    'Awatar poza polem widzenia, funkcja czynna w tle; dostęp skrótem klawiszowym i listwą ustawień.',
};

/** Interfejs opisuje sugestię tym, co reguła ujawniania musi o niej wiedzieć: klasę zdarzenia, moduł i kartę sesji, opcjonalne, gdy sugestia ich nie niesie. */
export interface OpisUjawnienia extends OpisSugestiiWobecWyciszenia {
  /** Waga ujawnienia. */
  waga: WagaUjawnienia;
  /** Punkt decyzyjny wstrzymujący proces — wyjątek wagi krytycznej. */
  krytyczna?: boolean;
}

/** Interfejs nazywa nastawy stanu obecności: magazyn zapisu i wykaz wyciszeń są podawane parametrem, nie brane na sztywno. */
export interface OpisStanuObecnosci {
  // Magazyn zapisu trybu i wyciszeń jest podawany; pominięty znaczy zapis miejscowy przeglądarki.
  magazyn?: MagazynWyciszen | null;
  /** Gotowy wykaz wyciszeń; pominięty buduje się nad {@link OpisStanuObecnosci.magazyn}. */
  wyciszenia?: StanWyciszen;
}

/** Interfejs opisuje stan obecności widziany przez warstwę widoku: tryb, wyciszenia oraz regułę otwierania dymka. */
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
  /** Czy funkcję obejmuje dziś jakiekolwiek wyciszenie — stan „Wyciszony". */
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
  // Wyciszenie czasowe wstrzymuje wszystko, nie liczy się jako wybiórcze; wyjątkiem jest waga krytyczna.
  czySugestiaWstrzymana(opis: OpisUjawnienia, teraz: number): boolean;

  /** Zgłasza obserwatora zmiany stanu; zwraca odsubskrybowanie. */
  obserwuj(sluchacz: () => void): () => void;
}

/** Stała podaje klucz zapisu trybu obecności w magazynie miejscowym; jeden klucz na stanowisko, bo funkcja jest globalna. */
const KLUCZ_ZAPISU = 'danaco.aod.obecnosc';

/** Interfejs opisuje kształt zapisu miejscowego trybu obecności, czytany defensywnie, bo zapis bywa cudzy albo starszy. */
interface ZapisObecnosci {
  tryb?: string;
}

export function utworzStanObecnosci(opis: OpisStanuObecnosci = {}): StanObecnosci {
  const magazyn = opis.magazyn === undefined ? magazynWyciszenDomyslny() : opis.magazyn;
  const wyciszenia = opis.wyciszenia ?? utworzStanWyciszen(magazyn);

  let tryb: TrybObecnosci = wczytajTryb(magazyn);

  /** Tryb sprzed ukrycia — żeby `Ctrl/Cmd + Shift + H` wracał tam, skąd wyszedł. */
  let trybPrzedUkryciem: TrybObecnosci = tryb === TrybObecnosci.Ukryty ? TrybObecnosci.Pelny : tryb;

  /** Chwile otwarcia dymków — podstawa limitu godzinowego i odstępu. */
  const dymki: number[] = [];

  const sluchacze = new Set<() => void>();

  function rozglos(): void {
    for (const sluchacz of sluchacze) sluchacz();
  }

  // Zmiana wykazu wyciszeń jest zmianą obecności: awatar wraca ze stanu „Wyciszony” bez odświeżenia.
  wyciszenia.obserwuj(rozglos);

  function zapiszTryb(): void {
    if (magazyn === null) return;
    try {
      magazyn.setItem(KLUCZ_ZAPISU, JSON.stringify({ tryb } satisfies ZapisObecnosci));
    } catch {
      // Zapis miejscowy bywa wyłączony ustawieniem przeglądarki; funkcja traci pamięć między wejściami.
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

    // Wyjątek wagi krytycznej sprawdza się pierwszy, obejmuje każde wyciszenie, ujawnia się plakietką.
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
      // Plakietkę chowa wyłącznie wyciszenie czasowe; kontekstowe wstrzymuje część sugestii, nie plakietkę.
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

/** Funkcja odczytuje tryb obecności z magazynu. Zapis bywa nieobecny, cudzy albo starszy od tego pliku — każda niezgodność daje tryb pełny, nie wyjątek. */
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
