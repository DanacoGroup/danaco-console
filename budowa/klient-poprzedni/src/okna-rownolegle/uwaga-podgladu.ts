import { elementIkony } from '../ikony/ikony';
import { PRZEKAZANIE_OPIS, PRZEKAZANIE_TYTUL } from './etykiety-ukladu';

/**
 * Komunikat blokowy sceny, trwały na ekranie dopóki sytuacja go wywołująca nie ustanie,
 * montowany w pasie relacji i ukryty do chwili pokazania.
 */
export interface UwagaPodgladu {
  /** Element montowany w pasie relacji; do czasu pokazania stoi ukryty. */
  element: HTMLElement;
  /** Wprowadza komunikat do układu. Powtórne wywołanie nic nie zmienia. */
  pokaz(): void;
}

/**
 * Tworzy komunikat blokowy informujący, że przekazanie zlecenia jest wyłącznie podglądem
 * układu, a nie wykonaną pracą, i pozostaje widoczny, dopóki nie ustanie jego przyczyna.
 */
export function utworzUwagePodgladu(): UwagaPodgladu {
  const element = document.createElement('p');
  element.className = 'dn-okna__uwaga';
  element.setAttribute('role', 'status');
  element.hidden = true;

  const tytul = document.createElement('strong');
  tytul.className = 'dn-okna__uwaga-tytul';
  tytul.textContent = PRZEKAZANIE_TYTUL;

  const opis = document.createElement('span');
  opis.textContent = PRZEKAZANIE_OPIS;

  element.append(elementIkony('ostrzezenie', { rozmiar: 16 }), tytul, opis);

  return {
    element,
    pokaz() {
      element.hidden = false;
    },
  };
}
