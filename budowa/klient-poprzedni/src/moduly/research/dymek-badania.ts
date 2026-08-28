import { utworzDymekObjasnienia } from '../../komponenty/dymek';

/**
 * Funkcja tworzy znak dymku objaśnienia przy elemencie konfiguracji okna Research na bazie wspólnego komponentu dymka.
 */
export function utworzDymekBadania(objasnienie: string): HTMLElement {
  return utworzDymekObjasnienia(objasnienie, { powloka: 'mr-dymek', znak: 'mr-dymek__znak' });
}

/** Funkcja opakowuje pole formularza razem z dymkiem objaśnienia w jeden kontener, zastępując osobne budowanie elementów w oknie. */
export function zDymkiem(pole: HTMLElement, objasnienie: string): HTMLElement {
  const element = document.createElement('div');
  element.className = 'mr-pole-z-dymkiem';
  element.append(pole, utworzDymekBadania(objasnienie));
  return element;
}
