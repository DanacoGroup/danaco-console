import { utworzDymekBadania } from './dymek-badania';

/**
 * Panel akcji okna operacyjnego — pasek narzędzi kontekstowych.
 *
 * Jedna odpowiedzialność: zamiana wykazu akcji okna na pasek przycisków.
 * Panel nie zna kontraktu i nie woła rdzenia — oddaje kod akcji wywołującemu,
 * bo droga wykonania (`window.action`, komenda dziedzinowa albo uczciwa odmowa)
 * należy do okna, nie do paska.
 *
 * Żaden przycisk nie jest wygaszany ani warunkowany zaznaczeniem. Brak warunku
 * merytorycznego nazywa odpowiedź po kliknięciu, nie odebranie klikalności.
 * Każdy przycisk ma dymek [?] mówiący, czym akcja jest i czym się kończy.
 */
export interface AkcjaOkna {
  /** Kod akcji przekazywany do `window.action` (pole `actionId`). */
  kod: string;
  /** Etykieta widoczna dla Operatora. */
  nazwa: string;
  /** Treść dymka [?] — czym akcja jest i jaka jest jej dzisiejsza droga. */
  objasnienie: string;
}

export interface PanelAkcji {
  element: HTMLElement;
}

export function utworzPanelAkcji<A extends AkcjaOkna>(
  akcje: readonly A[],
  naAkcje: (akcja: A) => void,
): PanelAkcji {
  const element = document.createElement('div');
  element.className = 'mr-akcje';
  element.setAttribute('role', 'toolbar');
  element.setAttribute('aria-label', 'Panel akcji okna');

  for (const akcja of akcje) {
    const pozycja = document.createElement('span');
    pozycja.className = 'mr-akcje__pozycja';

    const przycisk = document.createElement('button');
    przycisk.type = 'button';
    przycisk.className = 'dn-btn dn-btn--sm dn-btn--zarys';
    przycisk.textContent = akcja.nazwa;
    przycisk.dataset['akcja'] = akcja.kod;
    przycisk.addEventListener('click', () => naAkcje(akcja));

    pozycja.append(przycisk, utworzDymekBadania(akcja.objasnienie));
    element.append(pozycja);
  }

  return { element };
}
