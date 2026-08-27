/**
 * SKŁADNIK — DROGA POTWIERDZENIA Z LISTU.
 *
 * Pola po jednym znaku, odliczanie ważności i czynność poboczna. Pola
 * zaczynają PUSTE: okno pokazuje odsłonę przed wpisaniem, a nie udaje, że
 * użytkownik zdążył już coś wpisać.
 *
 * Każdy zestaw rządzi się sam — mechanika wiąże pola grupami, więc kursor nie
 * przeskakuje między odsłonami.
 *
 * Składnik zwraca WYKAZ węzłów: zestaw pól i stopka z odliczaniem są
 * rodzeństwem w kolumnie panelu.
 *
 * Liczba pól jest właściwością, nie stałą tego pliku: rozstrzyga ją długość
 * drogi potwierdzenia wydawanej przez rdzeń, a tę odczytuje się z rdzenia,
 * nie z okna.
 */

import { el, tekst } from '../narzedzia.ts';

/** Czynność w stopce zestawu. */
export type CzynnoscKodu = 'wklej' | 'ponow';

export interface WlasciwosciKodu {
  /** Ile znaków ma droga potwierdzenia. */
  znakow: number;
  /** Czas ważności w sekundach — mechanika odlicza go i sama wypisuje. */
  odliczanie: number;
  /**
   * Wklejenie jest wygodą, więc niesie przycisk; ponowienie wysyłki jest
   * wyjściem z sytuacji bez wyjścia, więc niesie odsyłacz — waga czynności
   * rozstrzyga o postaci kontrolki, nie odwrotnie.
   */
  czynnosc: CzynnoscKodu;
  /** Przedrostek identyfikatorów pól; wiąże zestaw z jedną odsłoną. */
  grupa: string;
}

export function poleKodu(w: WlasciwosciKodu): HTMLElement[] {
  const pola = Array.from({ length: w.znakow }, (_, i) =>
    el('input', {
      klasa: 'au-kod-pole',
      type: 'text',
      maxlength: '1',
      id: `${w.grupa}-znak-${i + 1}`,
      autocomplete: i === 0 ? 'one-time-code' : null,
      'aria-label': tekst('dostep.kod.znak', { numer: i + 1, ile: w.znakow }),
    }),
  );

  return [
    el(
      'div',
      {
        klasa: 'au-kod',
        role: 'group',
        'aria-label': tekst('dostep.kod.obszar'),
        dane: { 'kod-grupa': w.grupa },
      },
      pola,
    ),
    el('div', { klasa: 'au-kod-stopka' }, [
      el('span', { klasa: 'au-odliczanie' }, [
        `${tekst('dostep.kod.odliczanie')} `,
        el('b', { dane: { odliczanie: w.odliczanie } }),
      ]),
      w.czynnosc === 'ponow'
        ? el('button', {
            klasa: 'au-link',
            type: 'button',
            tekst: tekst('dostep.kod.ponow'),
            dane: { czynnosc: 'ponow-droge' },
          })
        : el('button', {
            klasa: 'dn-btn dn-btn--duch dn-btn--sm',
            type: 'button',
            tekst: tekst('dostep.kod.wklej'),
            dane: { 'wklej-kod': w.grupa },
          }),
    ]),
  ];
}

/**
 * Odczytuje drogę potwierdzenia z zestawu pól jednej grupy.
 *
 * Pierwszeństwo ma wartość odłożona przez wklejenie. Rdzeń wydaje drogę
 * dłuższą niż zestaw pól — zmierzone — więc wklejenie musi unieść ją w całości,
 * inaczej przycięłaby się do liczby pól i rdzeń odmówiłby drogi, która
 * przyszła listem poprawna. Wpisywanie znak po znaku zostaje bez zmiany.
 */
export function odczytajDroge(korzen: ParentNode, grupa: string): string {
  const zestaw = korzen.querySelector(`[data-kod-grupa="${grupa}"]`);
  if (zestaw === null) return '';
  const wklejona = (zestaw as HTMLElement).dataset['kodWklejony'];
  if (wklejona !== undefined && wklejona.length > 0) return wklejona;
  return [...zestaw.querySelectorAll('input')].map((pole) => pole.value).join('');
}

/** Odkłada wklejoną drogę na grupie i rozsypuje jej początek po polach. */
export function przyjmijWklejenie(zestaw: HTMLElement, wklejona: string): void {
  const droga = wklejona.trim();
  zestaw.dataset['kodWklejony'] = droga;
  const pola = [...zestaw.querySelectorAll('input')];
  pola.forEach((pole, i) => {
    pole.value = droga[i] ?? '';
  });
  pola[Math.min(droga.length, pola.length - 1)]?.focus();
}

/** Zdejmuje odłożone wklejenie, gdy Operator wpisuje drogę ręcznie. */
export function zapomnijWklejenie(zestaw: HTMLElement): void {
  delete zestaw.dataset['kodWklejony'];
}
