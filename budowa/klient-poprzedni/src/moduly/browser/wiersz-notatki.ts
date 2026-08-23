import type { BrowserNote, BrowserSource } from '../../../../shared/contract';
import { przycisk, wybor } from '../../modele/kontrolki-formularza';
import { KLASYFIKACJE } from './etykiety-browser';
import { KLASA_PRZYCISKU, przyciskCzynnosci } from './przyciski-browser';

/**
 * Jedna pozycja wykazu Notes Panel wraz z jej panelem akcji.
 *
 * Jedna odpowiedzialność: jeden wiersz notatki. Formularz, wykaz i przekazanie
 * mieszkają w oknie.
 *
 * Powiązanie ze źródłem jest widoczne, nie domyślne: notatka bez `sourceId`
 * mówi o tym wprost, bo powiązanie zakłada się z wykazu okna i jego brak jest
 * informacją, nie pustką do przemilczenia.
 *
 * Klasyfikacja stoi przy wierszu jako lista wyboru, a nie jako ikona: rodzaj
 * notatki ma być czytelny bez najeżdżania na znak i bez rozróżniania barw,
 * a zmiana ma być jednym gestem.
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

  // Napis zależy od stanu przypięcia, więc jeden przycisk niesie oba warianty
  // zamiast dwóch wywołań z literałami dublujących tę samą czynność.
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

/** Cytowany fragment strony zapisany razem z notatką. */
function cytat(tresc: string): HTMLElement {
  const element = document.createElement('blockquote');
  element.className = 'mb-notatki__cytat';
  element.textContent = tresc;
  return element;
}

