import type { KomunikatZmiany } from './komunikat-zmiany';

/**
 * Pasek komunikatów kompletu sterowania — jedyna reakcja interfejsu na niepowodzenie
 * zmiany ustawienia.
 */
export interface PasekKomunikatow {
  element: HTMLElement;
  /** Pokazuje komunikat o losie ostatniej zmiany. */
  pokaz(komunikat: KomunikatZmiany): void;
}

export function utworzPasekKomunikatow(): PasekKomunikatow {
  const element = document.createElement('p');
  element.className = 'dc-ster-komunikat';
  element.setAttribute('role', 'status');
  element.textContent = 'Zmiana ustawienia idzie do rdzenia komendą window.update.';

  return {
    element,

    pokaz(komunikat) {
      element.dataset.stan = komunikat.udany ? 'udany' : 'nieudany';
      element.textContent = komunikat.tresc;
    },
  };
}
