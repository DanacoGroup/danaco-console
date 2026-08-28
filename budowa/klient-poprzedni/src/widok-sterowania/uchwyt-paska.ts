import { elementIkony } from '../ikony/ikony';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import type { Szuflada } from './szuflada';

/** Uchwyt sterowania na pasku górnym powłoki jest drugim wejściem do tej samej szuflady: przycisk ikonowy bez stanu wyłączonego, rozwijający albo zwijający komplet kontrolek w kolumnie. */
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
