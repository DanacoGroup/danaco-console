import type { KomunikatZmiany } from './komunikat-zmiany';

/**
 * Pasek komunikatów kompletu sterowania.
 *
 * Jedyna reakcja interfejsu na niepowodzenie zmiany: informacja. Pasek nie
 * wyłącza sterowań, nie zamyka okna i nie wymusza potwierdzenia — kolejna próba
 * idzie zwyczajnie, bo błąd dotyczy wyłącznie bieżącego wywołania.
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
