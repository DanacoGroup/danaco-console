import { pokazKomunikat } from '../../aplikacja/komunikaty';
import type { BrakDrogi } from './etykiety-designu';

/**
 * Kontrolka czynności, dla której kontrakt nie ma komendy.
 *
 * Element klikalny, który zamiast pozorować wykonanie nazywa brakującą komendę:
 * przycisk nie jest ani `disabled`, ani pusty, a naciśnięcie daje dymek mówiący,
 * czego w kontrakcie nie ma. Wzór zachowania: `aplikacja/zamiary-pulpitu.ts`.
 */
export function przyciskBrakuDrogi(brak: BrakDrogi, klasa = 'dn-btn dn-btn--zarys dn-btn--sm'): HTMLButtonElement {
  const element = document.createElement('button');
  element.type = 'button';
  element.className = klasa;
  element.textContent = brak.nazwa;
  element.dataset['brak'] = brak.kod;
  // Znak „bez drogi" czytany przez technologie wspomagające zanim Operator
  // naciśnie — powód poznaje z dymka, zapowiedź z etykiety.
  element.setAttribute('aria-description', 'czynność bez komendy w kontrakcie');
  element.addEventListener('click', () => zglosBrakDrogi(brak));
  return element;
}

/** Pokazuje dymek z powodem braku drogi; jedyna reakcja na naciśnięcie. */
export function zglosBrakDrogi(brak: BrakDrogi): void {
  pokazKomunikat({
    tytul: `${brak.nazwa} — bez drogi w kontrakcie`,
    tresc: brak.powod,
    waga: 'ostrz',
  });
}

