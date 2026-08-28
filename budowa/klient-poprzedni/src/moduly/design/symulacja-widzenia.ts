/**
 * Symulacja wad widzenia barw na próbkach żetonów modułu Design. Plik
 * przekształca barwę macierzą jednej z trzech dichromazji, podaje jej zapis
 * w składni arkusza stylów i orzeka, czy dwie barwy zlewają się ze sobą po
 * przekształceniu.
 */
import type { Barwa } from './kontrast-wcag';

/**
 * Trzy dichromazje objęte symulacją. Wartość jest jednocześnie kluczem macierzy
 * przekształcenia i kluczem wykazu nazw pokazywanych Operatorowi, więc dodanie
 * rodzaju wymaga uzupełnienia obu zbiorów.
 */
export type RodzajWidzenia = 'protanopia' | 'deuteranopia' | 'tritanopia';

/**
 * Nazwy rodzajów widoczne dla Operatora wraz ze wskazaniem brakującego rodzaju
 * czopków. Pierwszy człon pary jest kluczem macierzy przekształcenia, drugi
 * zdaniem pokazywanym w panelu.
 */
export const RODZAJE_WIDZENIA: readonly (readonly [RodzajWidzenia, string])[] = [
  ['protanopia', 'protanopia — brak czopków czerwonych'],
  ['deuteranopia', 'deuteranopia — brak czopków zielonych'],
  ['tritanopia', 'tritanopia — brak czopków niebieskich'],
];

/**
 * Macierze przekształcenia składowych sRGB, wiersz po wierszu. Wartości
 * pochodzą z powszechnie stosowanego przybliżenia liniowego dichromazji
 * i stoją jako stała nazwana, bo liczba bez nazwy w takim rachunku jest nie do
 * sprawdzenia przy czytaniu.
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

/**
 * Barwa widziana przy danej dichromazji. Każda składowa wyniku jest ważoną sumą
 * trzech składowych wejścia według wiersza macierzy, ograniczoną do przedziału
 * od zera do dwustu pięćdziesięciu pięciu.
 */
export function symuluj(barwa: Barwa, rodzaj: RodzajWidzenia): Barwa {
  const macierz = MACIERZE[rodzaj];
  const skladowe = [barwa.r, barwa.g, barwa.b];
  const [r, g, b] = macierz.map((wiersz) =>
    ogranicz(wiersz.reduce((suma, waga, numer) => suma + waga * (skladowe[numer] ?? 0), 0)),
  ) as [number, number, number];
  return { r, g, b };
}

/**
 * Zapis barwy w składni funkcji rgb arkusza stylów, gotowy do wstawienia
 * w regułę stylu próbki. Składowe są zaokrąglane do liczb całkowitych, bo
 * rachunek prowadzony jest na wartościach ułamkowych.
 */
export function zapisBarwy(barwa: Barwa): string {
  return `rgb(${String(Math.round(barwa.r))} ${String(Math.round(barwa.g))} ${String(Math.round(barwa.b))})`;
}

/**
 * Czy dwie barwy zlewają się przy danej wadzie widzenia. Pytaniem symulacji nie
 * jest wygląd pojedynczej barwy, lecz to, czy dwie barwy rozróżniające stany
 * produktu dają się jeszcze od siebie odróżnić po przekształceniu.
 */
export function czySieZlewaja(pierwsza: Barwa, druga: Barwa, rodzaj: RodzajWidzenia): boolean {
  const a = symuluj(pierwsza, rodzaj);
  const b = symuluj(druga, rodzaj);
  const odleglosc = Math.sqrt((a.r - b.r) ** 2 + (a.g - b.g) ** 2 + (a.b - b.b) ** 2);
  return odleglosc < GRANICA_ZLANIA;
}

/**
 * Poniżej tej odległości w składowych sRGB dwie barwy uznaje się za zlane. Miara
 * jest zgrubna i dobrana tak, żeby wskazywała pary wymagające obejrzenia, a nie
 * żeby orzekała za człowieka.
 */
const GRANICA_ZLANIA = 24;

function ogranicz(wartosc: number): number {
  return Math.min(255, Math.max(0, wartosc));
}
