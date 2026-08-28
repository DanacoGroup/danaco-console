import type { BrowserNote, BrowserSource } from '../../../../shared/contract';
import { przycisk, wybor } from '../../modele/kontrolki-formularza';
import { KLASYFIKACJE } from './etykiety-browser';
import { KLASA_PRZYCISKU, przyciskCzynnosci } from './przyciski-browser';

/**
 * Jedna pozycja wykazu notatek wraz z panelem akcji: treść, powiązanie ze
 * źródłem, lista wyboru klasyfikacji, cytat oraz przyciski edycji, otwarcia
 * źródła i przypięcia. Formularz, wykaz i przekazanie do rdzenia mieszkają
 * w oknie modułu.
 */
export interface AkcjeNotatki {
  edytuj(notatka: BrowserNote): void;
  otworzZrodlo(notatka: BrowserNote): void;
  przypnij(notatka: BrowserNote): void;
  /** Nadaje notatce klasyfikację; pusty kod ją zdejmuje. */
  sklasyfikuj(notatka: BrowserNote, kod: string): void;
}

export function utworzWierszNotatki(
  notatka: BrowserNote,
  zrodlo: BrowserSource | null,
  przypieta: boolean,
  klasyfikacja: string,
  akcje: AkcjeNotatki,
): HTMLElement {
  const tresc = document.createElement('p');
  tresc.className = 'mb-notatki__tresc';
  tresc.textContent = notatka.content;

  const powiazanie = document.createElement('p');
  powiazanie.className = 'mb-notatki__powiazanie';
  powiazanie.textContent =
    zrodlo === null
      ? 'Notatka nie jest powiązana ze źródłem.'
      : `Źródło: ${(zrodlo.title ?? '').trim() === '' ? zrodlo.url : zrodlo.title ?? ''}`;

  // Napis zależy od stanu przypięcia, więc jeden przycisk niesie oba warianty.
  const przypnij = przycisk(przypieta ? 'Odepnij' : 'Przypnij', KLASA_PRZYCISKU.duch);
  przypnij.setAttribute('aria-pressed', String(przypieta));
  przypnij.addEventListener('click', () => akcje.przypnij(notatka));

  const rodzaj = wybor(
    'Klasyfikacja notatki',
    KLASYFIKACJE.map((pozycja) => [pozycja.kod, pozycja.nazwa] as const),
  );
  rodzaj.value = klasyfikacja;
  rodzaj.addEventListener('change', () => akcje.sklasyfikuj(notatka, rodzaj.value));

  const element = document.createElement('li');
  element.className = 'mb-notatki__wiersz';
  element.dataset['notatka'] = notatka.id;
  element.dataset['przypieta'] = przypieta ? 'tak' : 'nie';
  element.dataset['klasyfikacja'] = klasyfikacja === '' ? 'brak' : klasyfikacja;
  element.append(tresc, powiazanie, rodzaj);
  if ((notatka.quote ?? '').trim() !== '') element.append(cytat(notatka.quote ?? ''));
  element.append(
    przyciskCzynnosci('Edytuj', KLASA_PRZYCISKU.duch, () => akcje.edytuj(notatka)),
    przyciskCzynnosci('Otwórz źródło', KLASA_PRZYCISKU.duch, () => akcje.otworzZrodlo(notatka)),
    przypnij,
  );
  return element;
}

/**
 * Cytowany fragment strony zapisany razem z notatką. Wiersz dokłada go tylko
 * wtedy, gdy notatka niesie niepusty cytat, żeby pusty blok nie zajmował
 * miejsca w wykazie.
 */
function cytat(tresc: string): HTMLElement {
  const element = document.createElement('blockquote');
  element.className = 'mb-notatki__cytat';
  element.textContent = tresc;
  return element;
}

