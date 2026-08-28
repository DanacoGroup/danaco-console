import './panele.css';

import type { PanelWStosie } from './panel-w-stosie';

/**
 * Pionowy stos paneli pomocniczych jednego gniazda mieści dowolny ich podzbiór, znika przy zerze paneli i nie buduje ani nie zamyka paneli, które przez niego przechodzą.
 */
export interface KolumnaPaneli {
  /** Element osadzany w gnieździe, obok kolumny rozmowy. */
  element: HTMLElement;
  /** Ustawia zawartość stosu; pusty wykaz chowa kolumnę całkowicie. */
  ustaw(panele: readonly PanelWStosie[]): void;
  pusta(): boolean;
}

export function utworzKolumnaPaneli(): KolumnaPaneli {
  let ile = 0;

  const element = document.createElement('div');
  element.className = 'dn-okna__kolumna-paneli';
  element.setAttribute('aria-label', 'Panele pomocnicze otwarte obok rozmowy');
  element.hidden = true;

  return {
    element,
    ustaw(panele) {
      ile = panele.length;
      // Zdejmowanie poprzednich paneli z drzewa nie zamyka ich — zamknięcie należy do właściciela panelu.
      element.replaceChildren(...panele.map((panel) => panel.element));
      element.hidden = ile === 0;
    },
    pusta: () => ile === 0,
  };
}
