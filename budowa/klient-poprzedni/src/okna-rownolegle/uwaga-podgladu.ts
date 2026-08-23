import { elementIkony } from '../ikony/ikony';
import { PRZEKAZANIE_OPIS, PRZEKAZANIE_TYTUL } from './etykiety-ukladu';

/** Komunikat blokowy sceny — trwały, dopóki sytuacja nie ustanie. */
export interface UwagaPodgladu {
  /** Element montowany w pasie relacji; do czasu pokazania stoi ukryty. */
  element: HTMLElement;
  /** Wprowadza komunikat do układu. Powtórne wywołanie nic nie zmienia. */
  pokaz(): void;
}

/**
 * Uwaga: przekazanie zlecenia jest podglądem układu, nie wykonaną pracą.
 *
 * Komunikat blokowy pełnej szerokości rodzica, oznaczony wyłącznie cienką kreską
 * po lewej — bez tła i obwódki, żeby nie czytał się jako druga karta. Forma stoi
 * w arkuszu tego widoku pod nazwą własną, bo biblioteka nie niesie klasy `.dn-alert`.
 * Komunikat zostaje na scenie, dopóki nie ustanie powód — brak komendy przekazania
 * w kontrakcie — więc odpowiada na pytanie „dlaczego nic się nie stało" także później.
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
