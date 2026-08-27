/**
 * Składnik — pas działań. Domyka okno przez całą jego szerokość, pod kolumną
 * tożsamości też; czynności zbierają się przy prawej krawędzi.
 */

import { ikony, type NazwaZnaku } from '../ikony.ts';
import { el, tekst, zeZnacznika, type Dziecko } from '../narzedzia.ts';

export interface Czynnosc {
  /** Klucz katalogu — napis na kontrolce. */
  klucz: string;
  /** Klasa wyglądu; puste daje kontrolkę poboczną. */
  klasa?: string | null;
  /** Nazwa odsłony, do której czynność prowadzi bez wysyłania czegokolwiek. */
  cel?: string | null;
  /** Nazwa czynności obsługiwanej przez przebieg, na przykład logowanie. */
  czynnosc?: string | null;
  /** Gałąź katalogu komunikatu pokazywanego zamiast czynności. */
  komunikat?: string | null;
  ikona?: NazwaZnaku | null;
  id?: string | null;
}

export interface WlasciwosciPasa {
  /** Nazwa odsłony, z którą pas się przełącza; puste dla pasa jedynego. */
  widok?: string | null;
  /** Czy pas jest widoczny na starcie. */
  aktywny?: boolean;
  czynnosci: Czynnosc[];
}

/** Klasa kontrolki pobocznej — czynność, która nie jest wyjściem z odsłony, tylko dodatkiem obok czynności głównej. */
const KLASA_POBOCZNA = 'dn-btn dn-btn--duch dn-btn--sm';

export function pasDzialan(w: WlasciwosciPasa): HTMLElement {
  const czynnosci = w.czynnosci.map((c) => {
    const dzieci: Dziecko[] = [];
    if (c.ikona) {
      const znak = zeZnacznika(ikony[c.ikona]);
      znak.setAttribute('aria-hidden', 'true');
      dzieci.push(znak);
    }
    dzieci.push(tekst(c.klucz));
    return el(
      'button',
      {
        klasa: c.klasa ?? KLASA_POBOCZNA,
        type: 'button',
        id: c.id ?? null,
        dane: {
          idz: c.cel ?? null,
          czynnosc: c.czynnosc ?? null,
          komunikat: c.komunikat ? tekst(`komunikaty.${c.komunikat}.tresc`) : null,
          'komunikat-tytul': c.komunikat ? tekst(`komunikaty.${c.komunikat}.tytul`) : null,
        },
      },
      dzieci,
    );
  });

  return el(
    'div',
    {
      klasa: 'we-pas',
      dane: {
        widok: w.widok ?? null,
        'grupa-widoku': w.widok ? 'stan' : null,
        'widok-aktywny': w.widok ? (w.aktywny === true ? 'tak' : 'nie') : null,
      },
    },
    czynnosci,
  );
}
