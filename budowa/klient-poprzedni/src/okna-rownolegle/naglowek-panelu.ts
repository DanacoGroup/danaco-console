import './panele.css';

import { elementIkony, type NazwaIkony } from '../ikony/ikony';

/**
 * Gęsty nagłówek panelu to jeden wiersz z tytułem i rzędem ikon, w którym tytuł skraca arkusz, nie kod, a menu panelu dokłada gospodarz, więc rząd przyjmuje opisy czynności oraz gotowe elementy jako pozycje.
 */
export interface CzynnoscPanelu {
  /** Ikona z zestawu marki — dobrana pod czynność, nie pod wygląd narzędzia. */
  ikona: NazwaIkony;
  /** Nazwa czynności: `title` przycisku oraz jego etykieta dostępności. */
  etykieta: string;
  dzialanie(): void;
}

/**
 * Pozycja rzędu ikon jest sumą typów opisu czynności albo gotowego elementu od gospodarza, bo o kolejności w rzędzie rozstrzyga gospodarz, nie sztywne pole.
 */
export type PozycjaNaglowka = CzynnoscPanelu | HTMLElement;

/** Bok ikony w rzędzie czynności nagłówka panelu — najmniejszy rozmiar kanoniczny w całym zestawie ikon. */
const ROZMIAR_IKONY = 14;

/**
 * Buduje nagłówek panelu.
 *
 * `tytul` idzie do widoku, do `title` i do etykiety dostępności w postaci
 * pełnej; przycięcie jest wyłącznie zjawiskiem rysowania.
 */
export function utworzNaglowekPanelu(
  tytul: string,
  czynnosci: readonly PozycjaNaglowka[],
): HTMLElement {
  const element = document.createElement('header');
  element.className = 'dn-okna__naglowek-panelu';
  element.setAttribute('aria-label', `Panel ${tytul}`);

  const nazwa = document.createElement('span');
  nazwa.className = 'dn-okna__naglowek-panelu-tytul';
  nazwa.textContent = tytul;
  // Pełna nazwa pod kursorem — wielokropek zabiera ją oku, nie danym.
  nazwa.title = tytul;

  const rzad = document.createElement('div');
  rzad.className = 'dn-okna__naglowek-panelu-czynnosci';
  rzad.setAttribute('aria-label', `Czynności panelu ${tytul}`);
  for (const pozycja of czynnosci) {
    // Element gotowy idzie do rzędu bez tknięcia — etykiety i zachowanie należą do tego, kto go zbudował.
    rzad.append(pozycja instanceof HTMLElement ? pozycja : przyciskCzynnosci(pozycja));
  }

  element.append(nazwa, rzad);
  return element;
}

/**
 * Przycisk jednej czynności: sama ikona, bez napisu.
 *
 * Ikona zostaje ozdobna (`elementIkony` bez etykiety ukrywa ją przed odczytem),
 * a znaczenie niesie `aria-label` przycisku — inaczej czytnik przeczytałby
 * nazwę dwa razy.
 */
function przyciskCzynnosci(czynnosc: CzynnoscPanelu): HTMLButtonElement {
  const przycisk = document.createElement('button');
  przycisk.type = 'button';
  przycisk.className = 'dn-okna__naglowek-panelu-przycisk';
  przycisk.title = czynnosc.etykieta;
  przycisk.setAttribute('aria-label', czynnosc.etykieta);
  przycisk.append(elementIkony(czynnosc.ikona, { rozmiar: ROZMIAR_IKONY }));
  przycisk.addEventListener('click', () => czynnosc.dzialanie());
  return przycisk;
}
