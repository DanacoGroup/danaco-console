import { utworzDymekObjasnienia } from '../konfiguracja/dymek-objasnienia';
import { objasnienieSterowania } from './adnotacje-wykonania';

/**
 * Wiersz etykiety sterowania: nazwa, znak objaśnienia oraz miejsce na wskaźnik odczytu
 * obok samej kontrolki.
 */
export function utworzNaglowekSterowania(
  nazwa: string,
  identyfikator: string,
  ...dodatki: readonly HTMLElement[]
): HTMLElement {
  const element = document.createElement('div');
  element.className = 'dc-ster-pole__naglowek';

  const etykieta = document.createElement('label');
  etykieta.className = 'dc-ster-pole__etykieta';
  etykieta.htmlFor = identyfikator;
  etykieta.textContent = nazwa;

  element.append(etykieta, utworzDymekObjasnienia(objasnienieSterowania(nazwa)), ...dodatki);
  return element;
}
