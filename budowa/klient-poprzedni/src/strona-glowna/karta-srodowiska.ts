import { elementIkony } from '../ikony/ikony';
import type { KodSrodowiska, PozycjaSrodowiska } from './pozycje-srodowisk';

/** Karta środowiska jest elementem strefy pierwszej, środkiem ciężkości strony głównej. */

/** Godło środowiska renderowane jest w największej skali dostępnej w zestawie ikon całego systemu wizualnego. */
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
      // Stan czynny niesie też oznaczenie odczytywane przez czytnik ekranu, nie tylko kolor.
      if (czynna) {
        element.setAttribute('aria-current', 'true');
      } else {
        element.removeAttribute('aria-current');
      }
    },
  };
}
