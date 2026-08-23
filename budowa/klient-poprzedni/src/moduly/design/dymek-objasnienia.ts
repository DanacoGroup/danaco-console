/**
 * Dymek objaśnienia [?] wpięty w układ pola modułu Design.
 *
 * Odpowiada wyłącznie za miejsce znaku objaśnienia w wierszu pola — sam znak
 * buduje wspólna fabryka `komponenty/dymek.ts`.
 */
import { utworzDymekObjasnienia, type KlasyDymka } from '../../komponenty/dymek';

/**
 * Wygląd znaku [?] w tym module — pierścień 14 px zamiast bibliotecznego
 * przycisku ikonowego, powłoka wyrównana do linii pisma etykiety.
 *
 * Nazwy klas stoją tu w jednym miejscu, bo korzystają z nich dwaj odbiorcy:
 * wiersz pola i wiersz suwaka.
 */
export const KLASY_DYMKA: KlasyDymka = { powloka: 'md-dymek', znak: 'md-dymek__znak' };

/**
 * Dopina dymek do wiersza pola zbudowanego przez kontrolki formularza.
 *
 * Etykieta pola i znak objaśnienia stoją w jednym rzędzie — dymek wchodzi
 * obok etykiety, a nie pod kontrolką, żeby objaśnienie było widoczne przed
 * wpisaniem wartości, nie po nim.
 */
export function dopnijDymek(pole: HTMLElement, objasnienie: string): HTMLElement {
  const etykieta = pole.querySelector('label');
  if (etykieta === null) {
    pole.prepend(utworzDymekObjasnienia(objasnienie, KLASY_DYMKA));
    return pole;
  }
  const rzad = document.createElement('span');
  rzad.className = 'md-pole__rzad';
  etykieta.replaceWith(rzad);
  rzad.append(etykieta, utworzDymekObjasnienia(objasnienie, KLASY_DYMKA));
  return pole;
}
