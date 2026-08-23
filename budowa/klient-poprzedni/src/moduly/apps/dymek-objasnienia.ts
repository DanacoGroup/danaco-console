/**
 * Pole modułu Apps wraz z dymkiem objaśnienia [?] po jego prawej.
 *
 * Sam dymek buduje wspólna fabryka `komponenty/dymek.ts`; tutaj leży wstawka
 * układu właściwa temu modułowi — wiersz `mp-pole-z-dymkiem`, w którym pole
 * rośnie, a znak stoi przy jego krawędzi.
 *
 * Klasy własne idą parametrem, bo znak modułu nie jest bibliotecznym przyciskiem
 * ikonowym: `mp-dymek__znak` to obwódka `--dn-wym-ikona-sm` o narożniku
 * `--dn-r-pill` ze wskaźnikiem `help`, a `.dn-btn-ikona` to kwadrat 32 px.
 */

import { utworzDymekObjasnienia } from '../../komponenty/dymek';

/** Etykieta pola wraz z dymkiem objaśnienia po jej prawej. */
export function opiszPole(pole: HTMLElement, objasnienie: string): HTMLElement {
  const wiersz = document.createElement('div');
  wiersz.className = 'mp-pole-z-dymkiem';
  wiersz.append(
    pole,
    utworzDymekObjasnienia(objasnienie, { powloka: 'mp-dymek', znak: 'mp-dymek__znak' }),
  );
  return wiersz;
}
