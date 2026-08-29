/**
 * Składnik — droga potwierdzenia z listu. Pola po jednym znaku, odliczanie
 * ważności i czynność poboczna; pola zaczynają puste, bez udawania, że coś
 * już wpisano.
 */

import { el, tekst } from '../narzedzia.ts';

/** Czynność w stopce zestawu, rozstrzygająca, czy stoi tam przycisk wklejenia, czy odsyłacz ponowienia. */
export type CzynnoscKodu = 'wklej' | 'ponow';

export interface WlasciwosciKodu {
  /** Ile znaków ma droga potwierdzenia. */
  znakow: number;
  /** Czas ważności w sekundach — mechanika odlicza go i sama wypisuje. */
  odliczanie: number;
  /** Wklejenie jest wygodą, więc niesie przycisk; ponowienie wysyłki niesie odsyłacz, nie przycisk. */
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
 * Odczytuje drogę potwierdzenia z zestawu pól jednej grupy. Pierwszeństwo
 * ma wartość odłożona przez wklejenie, bo rdzeń wydaje drogę dłuższą niż
 * zestaw pól. Wpisywanie znak po znaku zostaje bez zmiany.
 */
export function odczytajDroge(korzen: ParentNode, grupa: string): string {
  const zestaw = korzen.querySelector(`[data-kod-grupa="${grupa}"]`);
  if (zestaw === null) return '';
  const wklejona = (zestaw as HTMLElement).dataset['kodWklejony'];
  if (wklejona !== undefined && wklejona.length > 0) return wklejona;
  return [...zestaw.querySelectorAll('input')].map((pole) => pole.value).join('');
}

/**
 * Odkłada wklejoną drogę na grupie i rozsypuje jej początek po polach, ustawiając skupienie na ostatnim wypełnionym polu.
 *
 * Odstępy znikają w całości, nie tylko brzegowe: list pokazuje kod rozdzielony
 * spacją co trzy znaki, więc skopiowany stamtąd niesie ją w środku. Zostawiona
 * zajęłaby jedno z pól i przesunęła resztę kodu o znak.
 */
export function przyjmijWklejenie(zestaw: HTMLElement, wklejona: string): void {
  const droga = wklejona.replace(/\s+/gu, '');
  zestaw.dataset['kodWklejony'] = droga;
  const pola = [...zestaw.querySelectorAll('input')];
  pola.forEach((pole, i) => {
    pole.value = droga[i] ?? '';
  });
  pola[Math.min(droga.length, pola.length - 1)]?.focus();
}

/** Zdejmuje odłożone wklejenie, gdy Operator wpisuje drogę ręcznie, znak po znaku, bez udziału wklejenia. */
export function zapomnijWklejenie(zestaw: HTMLElement): void {
  delete zestaw.dataset['kodWklejony'];
}
