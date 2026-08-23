import './panele.css';

import { elementIkony, type NazwaIkony } from '../ikony/ikony';

/**
 * Gęsty nagłówek panelu — jeden wiersz, tytuł i rząd ikon, nic więcej.
 *
 * Wszystkie czynności panelu mieszczą się w jednym rzędzie ikon o wysokości
 * jednego wiersza: bez paska narzędzi, bez drugiego rzędu, bez nagłówka sekcji
 * nad zawartością.
 *
 * Tytuł skraca arkusz (`text-overflow: ellipsis`), nie kod. Cięcie napisu
 * w kodzie odebrałoby pełną nazwę czytnikowi ekranu i podpowiedzi, a wielokropek
 * zależy od szerokości kolumny, której kod nie zna. Pełna nazwa zostaje więc
 * w `title` i w `aria-label` nagłówka.
 *
 * Menu `⋮` panelu dokłada gospodarz, nie nagłówek: pozycje takiego menu są
 * własnością panelu (Przeglądarka ma inne niż Terminal), a nagłówek ich nie zna.
 * Dlatego rząd przyjmuje obok opisu czynności także gotowy element
 * (`PozycjaNaglowka`) i daje mu wyłącznie miejsce w rzędzie.
 *
 * Nagłówek nie wie, co robią czynności, i nie zna ani jednej komendy rdzenia —
 * dostaje gotowe `dzialanie()`. Nie buduje treści panelu ani menu rozwijanego
 * i nie przyjmuje drugiego rzędu.
 *
 * Każdy przycisk jest czynny zawsze. Czynność, której nie ma czym wykonać, nie
 * wchodzi do wykazu — rząd ikon jest wtedy krótszy, a nie wyszarzony.
 *
 * Rama okna (`komponenty/rama-okna.ts`) daje trzy pasy pod tytułem (nagłówek,
 * akcje, narzędzia) i należy do okien operacyjnych; panel w stosie jej nie
 * używa.
 */

/** Jedna czynność w rzędzie ikon nagłówka panelu. */
export interface CzynnoscPanelu {
  /** Ikona z zestawu marki — dobrana pod czynność, nie pod wygląd narzędzia. */
  ikona: NazwaIkony;
  /** Nazwa czynności: `title` przycisku oraz jego etykieta dostępności. */
  etykieta: string;
  dzialanie(): void;
}

/**
 * Pozycja rzędu ikon: opis czynności albo gotowy element od gospodarza.
 *
 * Suma typów, a nie osobne pole `dodatki`, bo o kolejności w rzędzie rozstrzyga
 * gospodarz — menu `⋮` bywa raz przed pełnym ekranem, raz po nim, a osobne pole
 * narzucałoby jedno miejsce na sztywno.
 */
export type PozycjaNaglowka = CzynnoscPanelu | HTMLElement;

/** Bok ikony w rzędzie czynności — najmniejszy rozmiar kanoniczny zestawu. */
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
    // Element gotowy idzie do rzędu bez tknięcia — jego etykiety, `title`
    // i zachowanie należą do tego, kto go zbudował, i nagłówek ich nie zna.
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
