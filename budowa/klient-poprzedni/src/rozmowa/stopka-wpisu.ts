import { elementIkony } from '../ikony/ikony';
import { NAPISY } from './etykiety-rozmowy';
import { opisKonta } from './metadane-konta';
import { opisPodsumowania } from './podsumowanie-tury';
import type { WarstwyZapisu } from './widok-zapisu';
import type { BladWpisu, WpisRozmowy } from './wpis-rozmowy';

// Rozliczenie tury pod wypowiedzią: stopka i blok błędów, nie transkrypt wypowiedzi
// modelu.

/**
 * Części stopki: podsumowanie tury i konto kanału zawsze, a w wybranych trybach zdarzenia
 * tury i wykaz wytworów.
 */
export function czesciStopki(wpis: WpisRozmowy, warstwy: WarstwyZapisu): HTMLElement[] {
  const czesci: HTMLElement[] = [];
  if (warstwy.stopka && wpis.podsumowanie !== null) {
    czesci.push(notka(`${NAPISY.podsumowanie}: ${opisPodsumowania(wpis.podsumowanie)}`));
  }
  if (warstwy.stopka && wpis.konto !== null) {
    czesci.push(notka(`${NAPISY.konto}: ${opisKonta(wpis.konto)}`));
  }
  // Typy linii przechwyconych w turze; ma sens wyłącznie w trybie pełnym, poza nim jest
  // szumem.
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
 * Wykaz wytworów tury złożony z nazw narzędzi, które tura zdążyła uruchomić, każde nazwane
 * po jednym razie.
 */
export function nazwyWytworow(wpis: WpisRozmowy): string[] {
  const nazwy: string[] = [];
  for (const narzedzie of wpis.narzedzia) {
    if (narzedzie.nazwa.length > 0 && !nazwy.includes(narzedzie.nazwa)) nazwy.push(narzedzie.nazwa);
  }
  return nazwy;
}

/**
 * Jeden błąd tury: ikona, treść błędu oraz informacja o jego ponawialności, pokazywany
 * niezależnie od trybu widoku.
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

/**
 * Drobna notka stopki niosąca jeden fakt rozliczeniowy tury, zapisana czcionką o stałej
 * szerokości znaków, bez ozdobników.
 */
function notka(tekst: string): HTMLElement {
  const element = document.createElement('span');
  element.className = 'dc-wpis__notka';
  element.textContent = tekst;
  return element;
}
