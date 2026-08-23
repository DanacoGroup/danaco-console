import { elementIkony, type NazwaIkony } from '../ikony/ikony';

/**
 * Nagłówek sekcji pulpitu — wspólna szata tytułu, wspólna dla wszystkich sekcji.
 *
 * Tytuł idzie krojem szeryfowym zastrzeżonym dla nagłówków. Dopisek obok tytułu
 * idzie krojem bazowym i nie niesie stanu samym kolorem — obok barwy stoi słowo
 * albo ikona.
 */
export interface OpisNaglowka {
  /** Tytuł sekcji — krój szeryfowy. */
  tytul: string;
  /** Zdanie objaśniające, po co sekcja stoi na ekranie. */
  dopisek?: string;
  /** Ikona wiodąca sekcji. */
  ikona?: NazwaIkony;
  /** Element dostawiony przy prawej krawędzi nagłówka (plakietka, licznik). */
  przyPrawej?: HTMLElement;
}

/** Buduje nagłówek sekcji wraz z identyfikatorem do powiązania `aria-labelledby`. */
export function utworzNaglowekSekcji(opis: OpisNaglowka, idTytulu: string): HTMLElement {
  const naglowek = document.createElement('header');
  naglowek.className = 'mc-naglowek';

  const blok = document.createElement('div');
  blok.className = 'mc-naglowek__blok';

  const tytul = document.createElement('h2');
  tytul.className = 'mc-naglowek__tytul';
  tytul.id = idTytulu;
  if (opis.ikona) {
    tytul.append(elementIkony(opis.ikona, { rozmiar: 18, klasa: 'dn-ikona mc-naglowek__ikona' }));
  }
  const napis = document.createElement('span');
  napis.textContent = opis.tytul;
  tytul.append(napis);
  blok.append(tytul);

  if (opis.dopisek) {
    const dopisek = document.createElement('p');
    dopisek.className = 'mc-naglowek__dopisek';
    dopisek.textContent = opis.dopisek;
    blok.append(dopisek);
  }

  naglowek.append(blok);
  if (opis.przyPrawej) {
    opis.przyPrawej.classList.add('mc-naglowek__prawa');
    naglowek.append(opis.przyPrawej);
  }

  return naglowek;
}

/** Buduje ramę sekcji: element `section` z nagłówkiem i pustym ciałem. */
export function utworzSekcje(
  klasa: string,
  opis: OpisNaglowka,
  idTytulu: string,
): { element: HTMLElement; cialo: HTMLElement } {
  const element = document.createElement('section');
  element.className = `mc-sekcja ${klasa}`;
  element.setAttribute('aria-labelledby', idTytulu);

  const cialo = document.createElement('div');
  cialo.className = 'mc-sekcja__cialo';

  element.append(utworzNaglowekSekcji(opis, idTytulu), cialo);
  return { element, cialo };
}
