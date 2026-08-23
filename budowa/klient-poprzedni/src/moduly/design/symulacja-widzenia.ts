import type { Barwa } from './kontrast-wcag';

/**
 * Symulacja wad widzenia barw na próbkach żetonów.
 *
 * Zakres jest tu ważniejszy niż sam rachunek: symulacja obejmuje PRÓBKI ŻETONÓW,
 * a nie obrazy modułu. Obrazu nie ma czym pobrać — pole odsyłacza zasobu jest
 * ścieżką w systemie plików rdzenia — więc nałożenie symulacji na zasób
 * wizualny nie ma dziś drogi i panel tego nie udaje.
 *
 * Rachunek jest przybliżeniem i tak jest nazwany. Macierze odwzorowują trzy
 * dichromazje w przestrzeni sRGB bez przejścia przez przestrzeń długofalową,
 * co jest uproszczeniem przyjętym w narzędziach projektowych: wystarcza, żeby
 * zobaczyć, które dwie barwy systemu zlewają się w jedną, i nie wystarcza do
 * orzeczenia medycznego. Panel mówi to Operatorowi wprost, bo różnica między
 * „podglądem" a „badaniem" jest tu istotna.
 */

/** Rodzaj dichromazji objęty symulacją. */
export type RodzajWidzenia = 'protanopia' | 'deuteranopia' | 'tritanopia';

/** Nazwy widoczne dla Operatora wraz z tym, czego dotyczy każda wada. */
export const RODZAJE_WIDZENIA: readonly (readonly [RodzajWidzenia, string])[] = [
  ['protanopia', 'protanopia — brak czopków czerwonych'],
  ['deuteranopia', 'deuteranopia — brak czopków zielonych'],
  ['tritanopia', 'tritanopia — brak czopków niebieskich'],
];

/**
 * Macierze przekształcenia składowych sRGB, wiersz po wierszu.
 *
 * Wartości pochodzą z powszechnie stosowanego przybliżenia dichromazji
 * (przekształcenie liniowe w przestrzeni sRGB). Stoją jako stałe nazwane, bo
 * liczba bez nazwy w takim rachunku jest nie do sprawdzenia przy czytaniu.
 */
const MACIERZE: Readonly<Record<RodzajWidzenia, readonly (readonly number[])[]>> = {
  protanopia: [
    [0.567, 0.433, 0.0],
    [0.558, 0.442, 0.0],
    [0.0, 0.242, 0.758],
  ],
  deuteranopia: [
    [0.625, 0.375, 0.0],
    [0.7, 0.3, 0.0],
    [0.0, 0.3, 0.7],
  ],
  tritanopia: [
    [0.95, 0.05, 0.0],
    [0.0, 0.433, 0.567],
    [0.0, 0.475, 0.525],
  ],
};

/** Barwa widziana przy danej dichromazji. */
export function symuluj(barwa: Barwa, rodzaj: RodzajWidzenia): Barwa {
  const macierz = MACIERZE[rodzaj];
  const skladowe = [barwa.r, barwa.g, barwa.b];
  const [r, g, b] = macierz.map((wiersz) =>
    ogranicz(wiersz.reduce((suma, waga, numer) => suma + waga * (skladowe[numer] ?? 0), 0)),
  ) as [number, number, number];
  return { r, g, b };
}

/** Zapis barwy do wstawienia w regułę stylu próbki. */
export function zapisBarwy(barwa: Barwa): string {
  return `rgb(${String(Math.round(barwa.r))} ${String(Math.round(barwa.g))} ${String(Math.round(barwa.b))})`;
}

/**
 * Czy dwie barwy zlewają się przy danej wadzie widzenia.
 *
 * To jest właściwe pytanie symulacji: nie „jak wygląda barwa", tylko „czy dwie
 * barwy, którymi produkt rozróżnia stany, dają się jeszcze rozróżnić".
 * Granica jest odległością w składowych sRGB — miarą zgrubną, dobraną tak, żeby
 * wskazywała pary wymagające obejrzenia, a nie żeby orzekać za człowieka.
 */
export function czySieZlewaja(pierwsza: Barwa, druga: Barwa, rodzaj: RodzajWidzenia): boolean {
  const a = symuluj(pierwsza, rodzaj);
  const b = symuluj(druga, rodzaj);
  const odleglosc = Math.sqrt((a.r - b.r) ** 2 + (a.g - b.g) ** 2 + (a.b - b.b) ** 2);
  return odleglosc < GRANICA_ZLANIA;
}

/** Poniżej tej odległości w składowych sRGB dwie barwy uznajemy za zlane. */
const GRANICA_ZLANIA = 24;

function ogranicz(wartosc: number): number {
  return Math.min(255, Math.max(0, wartosc));
}
