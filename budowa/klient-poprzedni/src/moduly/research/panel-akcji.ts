import { utworzDymekBadania } from './dymek-badania';

/**
 * Panel akcji okna operacyjnego zamienia wykaz akcji na pasek przycisków, przekazując kod
 * wybranej akcji wywołującemu, ponieważ droga jej wykonania należy do okna.
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
