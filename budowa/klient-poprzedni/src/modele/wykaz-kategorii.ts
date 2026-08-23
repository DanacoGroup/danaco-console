import type { IdentityCategory } from '../../../shared/contract';
import type { StanTozsamosci } from './stan-tozsamosci';
import { WARSTWY_OD_NAJWAZNIEJSZEJ, nazwaWarstwy } from './warstwy-tozsamosci';
import { znakWykazu } from './znak-wykazu';

/**
 * Wykaz kategorii zasad — lewa kolumna panelu tożsamości, ułożona warstwami.
 *
 * Warstwa jest nagłówkiem grupy: mówi, jak krytyczna jest treść i w jakiej
 * kolejności wejdzie do złożonego promptu. Kategoria spoza warstw kontraktu nie
 * znika — dostaje własną grupę na końcu, żeby było widać, co przyszło z rdzenia.
 *
 * Wiersz mówi dwie rzeczy: czy kategoria jest obowiązkowa i czy ma treść
 * zapisaną na osi czynnej. Brak treści na osi nie znaczy braku treści w ogóle —
 * obowiązuje wtedy zapis z osi szerszej i wiersz nazywa to wprost.
 */
export interface WykazKategorii {
  /** Kolumna osadzana w panelu tożsamości. */
  element: HTMLElement;
  /** Przebudowuje wykaz ze stanu i oznacza kategorię czynną. */
  odswiez(): void;
}

export function utworzWykazKategorii(stan: StanTozsamosci): WykazKategorii {
  const element = document.createElement('nav');
  element.className = 'dn-boczna dm-kategorie';
  element.setAttribute('aria-label', 'Kategorie zasad i tożsamości modelu');

  function odswiez(): void {
    const kategorie = stan.kategorie();
    if (kategorie.length === 0) {
      element.replaceChildren(komunikatPustego());
      return;
    }

    const czynna = stan.wybrana();
    const pozycje: HTMLElement[] = [];
    for (const warstwa of warstwyWykazu(kategorie)) {
      pozycje.push(naglowekWarstwy(warstwa));
      for (const kategoria of kategorie.filter((pozycja) => pozycja.layer === warstwa)) {
        pozycje.push(wiersz(kategoria, kategoria.id === czynna?.id, stan));
      }
    }
    element.replaceChildren(...pozycje);
  }

  return { element, odswiez };
}

/** Warstwy obecne w katalogu: trzy z kontraktu, potem wszystkie pozostałe. */
function warstwyWykazu(kategorie: readonly IdentityCategory[]): string[] {
  const obecne = new Set<string>(kategorie.map((kategoria) => kategoria.layer));
  const znane = WARSTWY_OD_NAJWAZNIEJSZEJ.filter((warstwa) => obecne.has(warstwa));
  const pozostale = [...obecne].filter(
    (warstwa) => !(WARSTWY_OD_NAJWAZNIEJSZEJ as readonly string[]).includes(warstwa),
  );
  return [...znane, ...pozostale];
}

function naglowekWarstwy(warstwa: string): HTMLElement {
  const element = document.createElement('h4');
  element.className = 'dn-boczna-naglowek dm-kategorie__warstwa';
  element.textContent = nazwaWarstwy(warstwa);
  return element;
}

/** Jeden wiersz wykazu: nazwa kategorii oraz jej stan na osi czynnej. */
function wiersz(
  kategoria: IdentityCategory,
  czynna: boolean,
  stan: StanTozsamosci,
): HTMLElement {
  const przycisk = document.createElement('button');
  przycisk.type = 'button';
  przycisk.className = 'dn-boczna-pozycja dm-kategorie__pozycja';
  przycisk.dataset.kategoria = kategoria.id;
  if (czynna) przycisk.setAttribute('aria-current', 'true');
  przycisk.addEventListener('click', () => stan.wybierz(kategoria.id));

  const nazwa = document.createElement('span');
  nazwa.className = 'dm-kategorie__nazwa';
  nazwa.textContent = kategoria.name;

  const znaki = document.createElement('span');
  znaki.className = 'dm-kategorie__znaki';

  if (kategoria.required) {
    znaki.append(plakietka('obowiązkowa', 'dn-plakietka dn-plakietka--ostrzezenie'));
  }
  znaki.append(
    stan.dokument(kategoria.id) === null
      ? plakietka('bez zapisu na tej osi', 'dn-plakietka')
      : plakietka('treść zapisana', 'dn-plakietka dn-plakietka--sygnal'),
  );

  przycisk.append(nazwa, znaki);
  return przycisk;
}

/** Znak wiersza wykazu kategorii — plakietka biblioteki w miejscu znaków wiersza. */
function plakietka(tresc: string, klasa: string): HTMLElement {
  return znakWykazu(tresc, klasa, 'dm-kategorie__znak');
}

function komunikatPustego(): HTMLElement {
  const element = document.createElement('p');
  element.className = 'dm-kategorie__komunikat';
  element.textContent =
    'Katalog kategorii zasad nie dotarł z rdzenia albo jest pusty. Panel pozostaje czynny; odczytaj katalog ponownie, gdy rdzeń odpowie.';
  return element;
}
