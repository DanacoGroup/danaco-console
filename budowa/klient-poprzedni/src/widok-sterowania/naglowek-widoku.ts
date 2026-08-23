import { WindowRole } from '../../../shared/contract';
import { nazwaRoli } from '../okno-komunikacji/etykiety-okna';
import type { ObserwatorUstawien } from './obserwator-ustawien';

/**
 * Nagłówek kolumny sterowania: tytuł, rola okna i jego identyfikator.
 *
 * Tytuł jest jedynym miejscem kroju szeryfowego w tym widoku — krój ten należy
 * wyłącznie do nagłówków.
 *
 * Rola okna stoi w nagłówku, a nie tylko w podsumowaniu, ponieważ to ona
 * rozstrzyga, czym okno jest w pętli koordynator–wykonawca:
 * przy dwóch oknach obok siebie operator musi widzieć rolę bez rozwijania
 * czegokolwiek. Plakietka niesie ikonę słowną — nazwę roli — więc stan nie
 * opiera się na samej barwie.
 */
export interface NaglowekWidoku {
  /** Element montowany na szczycie kolumny sterowania. */
  element: HTMLElement;
  /** Odłącza subskrypcję stanu. */
  rozlacz(): void;
}

export function utworzNaglowekWidoku(
  obserwator: ObserwatorUstawien,
  idOkna: string,
): NaglowekWidoku {
  const element = document.createElement('header');
  element.className = 'dc-widok-ster__naglowek';

  const tytul = document.createElement('h2');
  tytul.className = 'dc-widok-ster__tytul';
  tytul.textContent = 'Sterowanie okna';

  // Plakietka roli okna (`--rola`: mono wersaliki).
  const rola = document.createElement('span');
  rola.className = 'dn-plakietka dn-plakietka--rola dc-widok-ster__rola';

  const identyfikator = document.createElement('span');
  identyfikator.className = 'dc-widok-ster__identyfikator';
  identyfikator.textContent = idOkna;
  identyfikator.title = `Identyfikator okna komunikacji: ${idOkna}`;

  const wiersz = document.createElement('div');
  wiersz.className = 'dc-widok-ster__wiersz-tytulu';
  wiersz.append(tytul, rola);

  element.append(wiersz, identyfikator);

  function odrysuj(): void {
    const biezaca = obserwator.migawka().okno.windowRole;
    rola.textContent = nazwaRoli(biezaca);
    rola.dataset.rola = biezaca;
    rola.title = opisRoli(biezaca);
  }

  const odsubskrybuj = obserwator.naZmiane(odrysuj);
  odrysuj();

  return { element, rozlacz: odsubskrybuj };
}

/** Co rola oznacza w pętli koordynator–wykonawca. */
function opisRoli(rola: WindowRole): string {
  switch (rola) {
    case WindowRole.Coordinator:
      return 'Okno planuje i rozdziela zlecenia oknom wykonawczym';
    case WindowRole.Executor:
      return 'Okno wykonuje zlecenia i raportuje koordynatorowi';
    case WindowRole.Standalone:
      return 'Okno pracuje samodzielnie, poza pętlą koordynator–wykonawca';
  }
}
