import { utworzDymekObjasnienia } from '../konfiguracja/dymek-objasnienia';
import { objasnienieSterowania } from './adnotacje-wykonania';

/**
 * Wiersz etykiety sterowania: nazwa, znak [?] z objaśnieniem i miejsce na
 * wskaźnik odczytu.
 *
 * Jeden wiersz obsługuje cały komplet — listę wyboru, suwak nakładu, pole
 * hosta i wykaz katalogów — bo wszystkie potrzebują tego samego nagłówka.
 *
 * Etykieta nie obudowuje kontrolki: znak [?] jest przyciskiem i wewnątrz
 * obudowy `<label>` jego naciśnięcie przenosiłoby się na kontrolkę zamiast
 * pokazać objaśnienie. Wiązanie idzie więc atrybutem `for`, nie zagnieżdżeniem.
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
