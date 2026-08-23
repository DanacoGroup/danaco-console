/**
 * Odczyt żetonów systemu wizualnego z motywu obowiązującego w tej chwili.
 *
 * Panel żetonów **nie kopiuje ani jednej wartości**. Wykaz poniżej niesie
 * wyłącznie NAZWY — role, które system wizualny nazywa — a wartość każdej
 * czytana jest z dokumentu przy każdym odświeżeniu. Dzięki temu panel pokazuje
 * stan produktu, a nie jego opis: przełączenie motywu zmienia tu wszystko bez
 * ani jednej linii kodu, a barwa poprawiona w arkuszu motywu jest widoczna
 * natychmiast.
 *
 * Konsekwencja jest zamierzona: żeton wymieniony w wykazie, którego motyw nie
 * definiuje, wychodzi na wierzch jako „motyw nie definiuje tego żetonu"
 * zamiast zniknąć. Rozjazd nazw między modułem a motywem ma być widoczny, bo
 * to jest usterka produktu, a nie szczegół panelu.
 *
 * Wykaz obejmuje żetony SEMANTYCZNE — te, po które wolno sięgać komponentom.
 * Prymitywów skali szarości tu nie ma, bo sięganie po nie wprost jest w tym
 * produkcie zabronione i panel nie ma ich po co pokazywać jako narzędzia pracy.
 *
 * Moduł nie zapisuje żetonów i nie nadpisuje motywu. Motyw jest własnością
 * powłoki; kontrakt nie zna bytu zestawu żetonów, więc trwałego zapisu nie ma
 * gdzie odłożyć — i panel mówi to wprost, zamiast udawać edytor.
 */

/** Rodzaj wartości żetonu — rozstrzyga, co panel z nią potrafi zrobić. */
export type RodzajZetonu = 'barwa' | 'miara' | 'krój' | 'czas' | 'liczba';

/** Jeden żeton wraz z wartością odczytaną z motywu. */
export interface Zeton {
  /** Nazwa bez przedrostka — tak nazywa rolę system wizualny. */
  readonly nazwa: string;
  /** Pełna nazwa własności niestandardowej. */
  readonly zmienna: string;
  readonly rodzaj: RodzajZetonu;
  /** Wartość z motywu; pusta znaczy „motyw tego żetonu nie definiuje". */
  readonly wartosc: string;
}

/** Grupa żetonów jednej roli. */
export interface GrupaZetonow {
  readonly nazwa: string;
  readonly zetony: readonly Zeton[];
}

/** Nazwy żetonów jednej grupy wraz z rodzajem ich wartości. */
interface OpisGrupy {
  readonly nazwa: string;
  readonly rodzaj: RodzajZetonu;
  readonly nazwyZetonow: readonly string[];
}

const GRUPY: readonly OpisGrupy[] = [
  {
    nazwa: 'Powierzchnie',
    rodzaj: 'barwa',
    nazwyZetonow: ['tlo', 'powierzchnia', 'powierzchnia-2', 'panel', 'nakladka'],
  },
  {
    nazwa: 'Tekst',
    rodzaj: 'barwa',
    nazwyZetonow: ['tekst', 'tekst-2', 'tekst-3', 'tekst-inv'],
  },
  {
    nazwa: 'Obrysy i wskazanie',
    rodzaj: 'barwa',
    nazwyZetonow: ['obrys', 'obrys-mocny', 'obrys-subtelny', 'hover', 'wcisniecie', 'fokus'],
  },
  {
    nazwa: 'Sygnał — jedyna barwa akcentu',
    rodzaj: 'barwa',
    nazwyZetonow: [
      'sygnal',
      'sygnal-mocny',
      'sygnal-tlo',
      'sygnal-obrys',
      'sygnal-wypelnienie',
      'sygnal-wypelnienie-hover',
      'kropka',
    ],
  },
  {
    nazwa: 'Stany',
    rodzaj: 'barwa',
    nazwyZetonow: [
      'sukces-tlo',
      'sukces-obrys',
      'sukces-tekst',
      'ostrzezenie-tlo',
      'ostrzezenie-obrys',
      'ostrzezenie-tekst',
      'blad-tlo',
      'blad-obrys',
      'blad-tekst',
      'informacja-tlo',
      'informacja-obrys',
      'informacja-tekst',
    ],
  },
  {
    nazwa: 'Rama kokpitu i atrament',
    rodzaj: 'barwa',
    nazwyZetonow: [
      'rama',
      'rama-tekst',
      'rama-tekst-2',
      'rama-hover',
      'rama-obrys',
      'atrament',
      'atrament-hover',
      'atrament-tekst',
    ],
  },
  {
    nazwa: 'Kroje pisma',
    rodzaj: 'krój',
    nazwyZetonow: ['ff-naglowek', 'ff-bazowa', 'ff-mono'],
  },
  {
    nazwa: 'Stopnie pisma',
    rodzaj: 'miara',
    nazwyZetonow: [
      'fs-2xs',
      'fs-xs',
      'fs-sm',
      'fs-base',
      'fs-md',
      'fs-lg',
      'fs-xl',
      'fs-2xl',
      'fs-3xl',
      'fs-display',
    ],
  },
  {
    nazwa: 'Grubości i interlinia',
    rodzaj: 'liczba',
    nazwyZetonow: [
      'fw-normalna',
      'fw-srednia',
      'fw-polgruba',
      'fw-gruba',
      'lh-ciasny',
      'lh-bazowy',
      'lh-luzny',
    ],
  },
  {
    nazwa: 'Odstępy liter',
    rodzaj: 'miara',
    nazwyZetonow: ['ls-naglowek', 'ls-wersaliki', 'ls-mono-wersaliki'],
  },
  {
    nazwa: 'Przestrzeń',
    rodzaj: 'miara',
    nazwyZetonow: [
      'od-0',
      'od-1',
      'od-2',
      'od-3',
      'od-4',
      'od-5',
      'od-6',
      'od-8',
      'od-10',
      'od-12',
      'od-16',
    ],
  },
  {
    nazwa: 'Promienie',
    rodzaj: 'miara',
    nazwyZetonow: ['r-xs', 'r-sm', 'r-md', 'r-lg', 'r-xl', 'r-pill'],
  },
  {
    nazwa: 'Wymiary gęstości zwartej',
    rodzaj: 'miara',
    nazwyZetonow: [
      'wym-kontrolka',
      'wym-ikonowy',
      'wym-pasek',
      'wym-pas-kart',
      'wym-wiersz',
      'wym-boczna',
      'wym-pas-komunikacji',
      'wym-modal',
      'wym-check',
      'wym-kropka',
    ],
  },
  {
    nazwa: 'Ruch',
    rodzaj: 'czas',
    nazwyZetonow: ['czas-1', 'czas-2', 'czas-3', 'czas-tetno'],
  },
  {
    nazwa: 'Cienie',
    rodzaj: 'miara',
    nazwyZetonow: ['cien-1', 'cien-2', 'cien-3', 'cien-lg', 'cien-sygnal'],
  },
];

/** Przedrostek własności niestandardowych systemu wizualnego. */
const PRZEDROSTEK = '--dn-';

/** Pełna nazwa własności niestandardowej danego żetonu. */
export function zmiennaZetonu(nazwa: string): string {
  return `${PRZEDROSTEK}${nazwa}`;
}

/**
 * Wartość żetonu z motywu obowiązującego.
 *
 * Czytana z korzenia dokumentu, bo tam motyw kładzie swoje definicje. Pusty
 * wynik znaczy „motyw tego żetonu nie definiuje" — stan inny niż wartość pusta
 * i panel go rozróżnia.
 */
export function wartoscZetonu(nazwa: string): string {
  const korzen = document.documentElement;
  return getComputedStyle(korzen).getPropertyValue(zmiennaZetonu(nazwa)).trim();
}

/** Wszystkie grupy wraz z wartościami odczytanymi w chwili wywołania. */
export function odczytajZetony(): readonly GrupaZetonow[] {
  return GRUPY.map((grupa) => ({
    nazwa: grupa.nazwa,
    zetony: grupa.nazwyZetonow.map((nazwa) => ({
      nazwa,
      zmienna: zmiennaZetonu(nazwa),
      rodzaj: grupa.rodzaj,
      wartosc: wartoscZetonu(nazwa),
    })),
  }));
}

/** Liczba żetonów, których motyw nie definiuje — miara rozjazdu, nie napis. */
export function liczbaBezDefinicji(grupy: readonly GrupaZetonow[]): number {
  return grupy.reduce(
    (suma, grupa) => suma + grupa.zetony.filter((zeton) => zeton.wartosc === '').length,
    0,
  );
}

/**
 * Selektory arkuszy produktu, w których żeton występuje — powiązanie żetonu
 * z komponentami.
 *
 * Przegląd idzie po arkuszach wczytanych do dokumentu. Arkusz z obcego źródła
 * odmawia dostępu do swoich reguł i rzuca wyjątkiem — taki arkusz pomijamy,
 * zamiast przerywać cały przegląd, bo w tym produkcie arkusze są własne
 * i pominięcie dotyczy co najwyżej zasobu wstrzykniętego przez środowisko.
 *
 * @param granica górna liczba selektorów w wyniku; przegląd kończy się po niej,
 *   bo wykaz na kilkaset pozycji nie jest odpowiedzią, tylko zrzutem.
 */
export function selektoryZetonu(nazwa: string, granica = 40): readonly string[] {
  const szukane = `var(${zmiennaZetonu(nazwa)}`;
  const znalezione: string[] = [];

  for (const arkusz of Array.from(document.styleSheets)) {
    let reguly: CSSRuleList;
    try {
      reguly = arkusz.cssRules;
    } catch {
      // Arkusz spoza źródła dokumentu nie oddaje reguł. Pomijamy go w ciszy:
      // to nie jest usterka panelu ani produktu.
      continue;
    }
    if (zbierzZRegul(reguly, szukane, znalezione, granica)) return znalezione;
  }
  return znalezione;
}

/**
 * Zbiera selektory z listy reguł, wchodząc w reguły zagnieżdżone.
 *
 * Zagnieżdżenie ma znaczenie: definicje motywu ciemnego stoją w regule
 * warunkowej, więc przegląd płaski przeoczyłby połowę produktu.
 *
 * @returns czy granica została osiągnięta.
 */
function zbierzZRegul(
  reguly: CSSRuleList,
  szukane: string,
  znalezione: string[],
  granica: number,
): boolean {
  for (const regula of Array.from(reguly)) {
    if (znalezione.length >= granica) return true;
    if (regula instanceof CSSStyleRule) {
      if (regula.style.cssText.includes(szukane)) znalezione.push(regula.selectorText);
      continue;
    }
    // Reguła grupująca (warunek motywu, punkt łamania, warstwa) niesie własne
    // reguły — wchodzimy w nie, zamiast je pomijać.
    const zagniezdzone = (regula as CSSGroupingRule).cssRules as CSSRuleList | undefined;
    if (zagniezdzone !== undefined && zbierzZRegul(zagniezdzone, szukane, znalezione, granica)) {
      return true;
    }
  }
  return false;
}
