import { elementIkony } from '../ikony/ikony';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import type { Szuflada } from './szuflada';

/**
 * Uchwyt sterowania na pasku górnym powłoki.
 *
 * Drugie wejście do tej samej szuflady — z paska, obok przełącznika motywu.
 * Na pasku stoi wyłącznie przycisk ikonowy o wymiarze kontrolki paska; komplet
 * kontrolek mieszka w kolumnie obok sceny, bo w prawym rogu paska nie miałby
 * się gdzie zmieścić.
 *
 * Przycisk nie ma stanu wyłączonego. Naciśnięty przy zwiniętej szufladzie
 * rozwija ją, przy rozwiniętej — zwija.
 */
export interface UchwytPaska {
  /** Przycisk montowany w miejscu akcji paska górnego. */
  element: HTMLElement;
  /** Odłącza subskrypcję stanu szuflady. */
  rozlacz(): void;
}

export function utworzUchwytPaska(szuflada: Szuflada): UchwytPaska {
  // Wariant „na ramie" — pasek górny jest zawsze atramentowy.
  const element = document.createElement('button');
  element.type = 'button';
  element.className = 'dn-btn-ikona dn-btn-ikona--na-ramie dc-widok-ster__uchwyt-paska';
  element.append(elementIkony('ustawienia', { rozmiar: 18 }));
  element.addEventListener('click', szuflada.przelacz);

  function odrysuj(otwarta: boolean): void {
    element.setAttribute('aria-expanded', String(otwarta));
    const opis = otwarta
      ? 'Zwiń sterowanie okna'
      : 'Rozwiń sterowanie okna — osiem ustawień';
    element.setAttribute('aria-label', opis);
    element.title = opis;
  }

  const odsubskrybuj: Odsubskrybuj = szuflada.naZmiane(odrysuj);
  odrysuj(szuflada.otwarta());

  return { element, rozlacz: odsubskrybuj };
}
