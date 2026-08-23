import { elementIkony } from '../ikony/ikony';
import { NAPISY } from './etykiety-rozmowy';
import { opisKonta } from './metadane-konta';
import { opisPodsumowania } from './podsumowanie-tury';
import type { WarstwyZapisu } from './widok-zapisu';
import type { BladWpisu, WpisRozmowy } from './wpis-rozmowy';

/**
 * Rozliczenie tury pod wypowiedzią — stopka i blok błędów.
 *
 * Rozliczenie, nie transkrypt: podsumowanie tury, konto kanału, typy zdarzeń
 * i wykaz wytworów mówią, ile tura kosztowała i co po sobie zostawiła, a nie co
 * model powiedział. Widok wpisu składa warstwy i steruje nimi; ten plik wie,
 * jak wygląda jedna notka.
 *
 * O tym, które części pokazać, ten plik nie rozstrzyga — rozkład warstw
 * przychodzi gotowy z `widok-zapisu.ts`.
 */

/**
 * Części stopki: podsumowanie tury, konto kanału, a w wybranych trybach także
 * zdarzenia tury (`pelny`) i wykaz wytworów (`streszczenie`).
 *
 * Podsumowanie i konto stoją w każdym trybie — na tej warstwie opiera się tryb
 * „Streszczenie", a jej zdjęcie w trybie „Zwykłym" odebrałoby Operatorowi
 * rozliczenie tury dostępne bez przełącznika.
 */
export function czesciStopki(wpis: WpisRozmowy, warstwy: WarstwyZapisu): HTMLElement[] {
  const czesci: HTMLElement[] = [];
  if (warstwy.stopka && wpis.podsumowanie !== null) {
    czesci.push(notka(`${NAPISY.podsumowanie}: ${opisPodsumowania(wpis.podsumowanie)}`));
  }
  if (warstwy.stopka && wpis.konto !== null) {
    czesci.push(notka(`${NAPISY.konto}: ${opisKonta(wpis.konto)}`));
  }
  // Typy linii przechwyconych w turze — przejrzystość kanału, nie diagnostyka
  // błędu. Tryb „Pełny" jest jedynym miejscem, w którym ta lista ma sens: poza
  // nim jest szumem nad odpowiedzią.
  const zdarzenia = wpis.podsumowanie?.typyZdarzen ?? [];
  if (warstwy.zdarzenia && zdarzenia.length > 0) {
    czesci.push(notka(`${NAPISY.zdarzenia}: ${zdarzenia.join(', ')}`));
  }
  const wytwory = nazwyWytworow(wpis);
  if (warstwy.wytwory && wytwory.length > 0) {
    czesci.push(notka(`${NAPISY.wytwory}: ${wytwory.join(', ')}`));
  }
  return czesci;
}

/**
 * Wykaz wytworów tury.
 *
 * Spisu plików tu nie ma: rdzeń go nie nadaje — wpis niesie wywołania narzędzi
 * (`WywolanieNarzedzia`), a nie listę tego, co po nich zostało na dysku.
 * Wytworem tury są więc nazwy narzędzi, które tura uruchomiła, każda raz.
 */
export function nazwyWytworow(wpis: WpisRozmowy): string[] {
  const nazwy: string[] = [];
  for (const narzedzie of wpis.narzedzia) {
    if (narzedzie.nazwa.length > 0 && !nazwy.includes(narzedzie.nazwa)) nazwy.push(narzedzie.nazwa);
  }
  return nazwy;
}

/**
 * Jeden błąd tury: ikona, treść, informacja o ponawialności.
 *
 * Błąd nie podlega trybowi widoku transkryptu. Wpis, w którym tura padła, mówi
 * o tym zawsze — schowanie błędu za ustawieniem widoku byłoby ciszą w miejscu,
 * gdzie Operator musi wiedzieć, że kanał odmówił.
 */
export function wierszBledu(blad: BladWpisu): HTMLElement {
  const element = document.createElement('p');
  element.className = 'dc-blad';
  element.append(
    elementIkony('blad', { rozmiar: 16, klasa: 'dn-ikona', etykieta: NAPISY.bledy }),
  );

  const tresc = document.createElement('span');
  const ogon = blad.ponawialny ? ' — ponowienie ma sens' : '';
  tresc.textContent =
    blad.kod.length > 0 ? `${blad.tresc} (${blad.kod})${ogon}` : `${blad.tresc}${ogon}`;
  element.append(tresc);
  return element;
}

/** Drobna notka stopki — jeden fakt rozliczeniowy w kroju monospacjowym. */
function notka(tekst: string): HTMLElement {
  const element = document.createElement('span');
  element.className = 'dc-wpis__notka';
  element.textContent = tekst;
  return element;
}
