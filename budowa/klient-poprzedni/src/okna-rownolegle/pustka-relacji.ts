import { elementIkony } from '../ikony/ikony';
import { PUSTKA_RELACJI_OPIS, PUSTKA_RELACJI_TYTUL } from './etykiety-ukladu';

/**
 * Moduł tworzy widok stanu pustego pasa relacji, informujący, że położenie pary jest oczekiwane i nie zostało jeszcze wczytane.
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
