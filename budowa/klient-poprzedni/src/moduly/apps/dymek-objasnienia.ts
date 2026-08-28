/**
 * Pole modułu Apps wraz z dymkiem objaśnienia po jego prawej stronie. Sam dymek
 * buduje wspólna fabryka `komponenty/dymek.ts`, a tutaj leży wstawka układu
 * właściwa temu modułowi: wiersz `mp-pole-z-dymkiem`.
 */

import { utworzDymekObjasnienia } from '../../komponenty/dymek';

/**
 * Składa etykietę pola z dymkiem objaśnienia po jej prawej stronie i oddaje
 * gotowy wiersz układu, w którym pole rośnie, a znak dymka stoi przy jego
 * krawędzi.
 */
export function opiszPole(pole: HTMLElement, objasnienie: string): HTMLElement {
  const wiersz = document.createElement('div');
  wiersz.className = 'mp-pole-z-dymkiem';
  wiersz.append(
    pole,
    utworzDymekObjasnienia(objasnienie, { powloka: 'mp-dymek', znak: 'mp-dymek__znak' }),
  );
  return wiersz;
}
