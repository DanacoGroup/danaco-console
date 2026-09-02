/* Kształty znaków i czas stania powiadomienia. Wartości przepisane z warstwy
   prototypu `design/zasoby/prototyp.js`, której produkt nie wczytuje. */

/** Wagi komunikatu rozpoznawane przez `komponenty.css`. */
export type RodzajOgloszenia = 'informacja' | 'sukces' | 'ostrzezenie' | 'blad';

/** Kształt znaku: nazwa elementu SVG wraz z cechami. */
export interface KsztaltZnaku {
  element: string;
  cechy: Record<string, string>;
}

/** Znak wagi komunikatu; ten sam zestaw, który niesie warstwa projektowa. */
export const KSZTALTY_ZNAKU: Record<RodzajOgloszenia, KsztaltZnaku[]> = {
  informacja: [
    { element: 'circle', cechy: { cx: '12', cy: '12', r: '10' } },
    { element: 'path', cechy: { d: 'M12 16v-4' } },
    { element: 'path', cechy: { d: 'M12 8h.01' } },
  ],
  sukces: [
    { element: 'circle', cechy: { cx: '12', cy: '12', r: '10' } },
    { element: 'path', cechy: { d: 'm9 12 2 2 4-4' } },
  ],
  ostrzezenie: [
    { element: 'path', cechy: { d: 'm21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3' } },
    { element: 'path', cechy: { d: 'M12 9v4' } },
    { element: 'path', cechy: { d: 'M12 17h.01' } },
  ],
  blad: [
    { element: 'circle', cechy: { cx: '12', cy: '12', r: '10' } },
    { element: 'path', cechy: { d: 'm15 9-6 6' } },
    { element: 'path', cechy: { d: 'm9 9 6 6' } },
  ],
};

/** Czas stania komunikatu na ekranie w milisekundach. */
const CZAS_DOMYSLNY = 3200;

/** Czas stania podany przez wołającego; brak i wartość niedodatnia biorą domyślny. */
export function godzinaZycia(ms: number): number {
  return Number.isFinite(ms) && ms > 0 ? ms : CZAS_DOMYSLNY;
}
