import type { PokrycieKomend } from '../pokrycie-komend';
import type { PozycjaBezObslugi } from './etykiety-browser';

/**
 * Pozycje modułu Browser, których okno jeszcze nie wykonuje, wraz z odczytem
 * ich pokrycia w rdzeniu. Taka pozycja staje pełnoprawnym przyciskiem, a jej
 * naciśnięcie oddaje powód wzięty z odczytu, a nie z napisu. Do rdzenia nie
 * idzie nic.
 */
export function przyciskPozycji(
  pokrycie: PokrycieKomend,
  pozycja: PozycjaBezObslugi,
  etykieta: string,
  powiedz: (tresc: string, powodzenie: boolean) => void,
): HTMLButtonElement {
  const element = pokrycie.przycisk(etykieta, pozycja.komenda, pozycja.czynnosc);
  element.classList.add('dn-btn--sm');
  // Powód czytany przy naciśnięciu, bo byt przerysowuje go po każdej zmianie wykazu.
  element.addEventListener('click', () => powiedz(element.title, false));
  return element;
}
