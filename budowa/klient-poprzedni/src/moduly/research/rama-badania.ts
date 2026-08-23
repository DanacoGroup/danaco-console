import { utworzRameOkna, type RolaOkna } from '../../komponenty/rama-okna';
import { utworzPanelAkcji, type AkcjaOkna } from './panel-akcji';

/**
 * Okno operacyjne Research: wspólna rama biblioteki z panelem akcji modułu.
 *
 * Obudowa okna jest jedna dla całego interfejsu (`komponenty/rama-okna`), ale
 * wykaz akcji i ich dymki [?] należą do modułu — panel akcji wnosi generyk `<A>`,
 * którego rama nie potrzebuje. Tutaj okna Research składają te dwie rzeczy, więc
 * przedrostek klas modułu i sposób osadzenia panelu stoją w jednym miejscu.
 *
 * Przedrostek `mr` nie wnosi wyglądu — jest uchwytem jednej reguły własnej
 * modułu: `.mr-okno[data-ognisko='tak']`, czyli obwiedzenia okna, do którego
 * nawigacja wewnątrzmodułowa przeniosła ognisko.
 *
 * Rama nie buduje elementów treści — dostaje gotowe, tak jak przestrzeń modułu
 * dostaje gotowe okna.
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
