import { elementIkony } from '../ikony/ikony';
import type { KodSrodowiska, PozycjaSrodowiska } from './pozycje-srodowisk';

/**
 * Karta środowiska — element strefy pierwszej, środek ciężkości strony.
 *
 * Buduje jedną kartę i zgłasza jej wybór. Układ treści: godło → tytuł krojem
 * nagłówkowym 24/30 px → motto → jednozdaniowy opis trybu pracy krojem bazowym.
 *
 * Stany: w spoczynku powierzchnia neutralna, bez akcentu; przy najechaniu i w
 * stanie czynnym wstęga górna 2 px w błękicie sygnałowym, cień sygnału
 * i uniesienie o 2 px — całość z wariantu `.dn-karta--akcent` biblioteki
 * komponentów. Arkusz strony wygasza wstęgę w spoczynku, bo biblioteka pokazuje
 * ją stale. Sygnał jest jedyną barwą akcentu systemu wizualnego v2.0.
 *
 * Nośnikiem jest `<button>`, nie `<div role="button">`: karta ma być celem
 * nawigacji klawiaturą z pierwszeństwem natywnym, a pierścień fokusu wnosi
 * `.dn-karta--klikalna:focus-visible`.
 */

/** Godło środowiska w największej skali renderowania zestawu ikon. */
const ROZMIAR_GODLA = 24;

export interface KartaSrodowiska {
  element: HTMLButtonElement;
  kod: KodSrodowiska;
  /** Nadaje albo zdejmuje oznaczenie środowiska czynnego. */
  oznaczCzynna(czynna: boolean): void;
}

export function utworzKarteSrodowiska(
  pozycja: PozycjaSrodowiska,
  przyWyborze: (pozycja: PozycjaSrodowiska) => void,
): KartaSrodowiska {
  const element = document.createElement('button');
  element.type = 'button';
  element.className =
    'dn-karta dn-karta--klikalna dn-karta--akcent dn-strona__karta';
  element.dataset.srodowisko = pozycja.kod;

  const tresc = document.createElement('span');
  tresc.className = 'dn-karta-tresc dn-strona__karta-tresc';

  const godlo = document.createElement('span');
  godlo.className = 'dn-strona__godlo';
  godlo.append(elementIkony(pozycja.godlo, { rozmiar: ROZMIAR_GODLA }));

  const tytul = document.createElement('span');
  tytul.className = 'dn-karta-tytul dn-strona__karta-tytul';
  tytul.textContent = pozycja.nazwa;

  const motto = document.createElement('span');
  motto.className = 'dn-strona__karta-motto';
  motto.textContent = pozycja.motto;

  const opis = document.createElement('span');
  opis.className = 'dn-karta-opis dn-strona__karta-opis';
  opis.textContent = pozycja.opis;

  tresc.append(godlo, tytul, motto, opis);
  element.append(tresc);

  element.addEventListener('click', () => przyWyborze(pozycja));

  return {
    element,
    kod: pozycja.kod,
    oznaczCzynna(czynna) {
      element.classList.toggle('dn-karta--wybrana', czynna);
      // Stan czynny nie opiera się na samym kolorze — niesie go także
      // oznaczenie odczytywane przez czytnik ekranu.
      if (czynna) {
        element.setAttribute('aria-current', 'true');
      } else {
        element.removeAttribute('aria-current');
      }
    },
  };
}
