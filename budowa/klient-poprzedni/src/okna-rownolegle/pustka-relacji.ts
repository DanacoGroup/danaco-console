import { elementIkony } from '../ikony/ikony';
import { PUSTKA_RELACJI_OPIS, PUSTKA_RELACJI_TYTUL } from './etykiety-ukladu';

/**
 * Stan pusty pasa relacji.
 *
 * Pas przy braku pary nie znika, bo zniknięcie niczego nie mówi: Operator nie
 * wie wtedy, czy pary nie ma, czy pas się nie wczytał, czy pomylił widok. Stan
 * pusty mówi wprost, że położenie jest oczekiwane, i wskazuje jedną czynność,
 * która je zmienia.
 *
 * Bez wariantów: biblioteka daje jedną formę pustego stanu, różnicowaną
 * wyłącznie treścią.
 */
export function utworzPustkeRelacji(): HTMLElement {
  const element = document.createElement('div');
  element.className = 'dn-pusty-stan dn-okna__pustka';

  const tytul = document.createElement('p');
  tytul.className = 'dn-pusty-stan-tytul';
  tytul.textContent = PUSTKA_RELACJI_TYTUL;

  const opis = document.createElement('p');
  opis.className = 'dn-pusty-stan-opis';
  opis.textContent = PUSTKA_RELACJI_OPIS;

  element.append(elementIkony('info', { rozmiar: 24 }), tytul, opis);
  return element;
}
