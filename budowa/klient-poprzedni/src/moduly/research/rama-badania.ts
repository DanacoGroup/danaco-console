import { utworzRameOkna, type RolaOkna } from '../../komponenty/rama-okna';
import { utworzPanelAkcji, type AkcjaOkna } from './panel-akcji';

/**
 * Okno operacyjne Research zestawia wspólną ramę biblioteki interfejsu z panelem akcji modułu,
 * wnoszącym generyk akcji, którego sama rama nie potrzebuje.
 */
export interface RamaBadania {
  /** Sekcja osadzana w przestrzeni modułu. */
  element: HTMLElement;
  /** Miejsce na treść okna wraz z jego stanem. */
  cialo: HTMLElement;
}

export function utworzRameBadania<A extends AkcjaOkna>(
  kodOkna: string,
  tytul: string,
  rola: RolaOkna,
  akcje: readonly A[],
  naAkcje: (akcja: A) => void,
): RamaBadania {
  const rama = utworzRameOkna({
    tytul,
    rola,
    kod: kodOkna,
    przedrostek: 'mr',
    ogniskowalne: true,
  });
  rama.akcje.append(utworzPanelAkcji(akcje, naAkcje).element);
  return { element: rama.element, cialo: rama.cialo };
}
