import { czyRozwiniete, zapamietajZwiniecie } from './pamiec-zwiniecia';

/** Strefa przywoływana trzyma zapowiedź na ekranie zawsze i czyta dane treści dopiero przy rozwinięciu, pamiętając jego postać między wejściami. */
export interface OpisStrefyZwijanej {
  /** Etykieta wersalikowa — nazwa strefy. */
  etykieta: string;
  /** Jedno zdanie: po co strefa stoi na ekranie. */
  wyjasnienie: string;
  /** Klucz pamięci postaci, stały między wejściami; jego zmiana znaczy „zacznij od domyślnej". */
  klucz: string;
  /** Postać przed pierwszym rozstrzygnięciem Operatora. */
  domyslnieRozwiniete: boolean;
}

export interface StrefaZwijana {
  /** Element `<details>` do osadzenia w stronie. */
  element: HTMLDetailsElement;
  /** Miejsce na treść strefy — rysowane niezależnie od postaci. */
  cialo: HTMLElement;
  /** Dopisek zapowiedzi z liczbą pod spodem; napis pusty zdejmuje dopisek zamiast zostawić puste miejsce. */
  ustawDopisek(tekst: string): void;
  /** Podpina odbiorcę zmiany postaci; wywoływany także przy przywołaniu. */
  naZmianePostaci(sluchacz: (rozwiniete: boolean) => void): void;
  /** Czy strefa jest w tej chwili rozwinięta. */
  rozwiniete(): boolean;
}

export function utworzStrefeZwijana(opis: OpisStrefyZwijanej): StrefaZwijana {
  const sluchacze: ((rozwiniete: boolean) => void)[] = [];

  const element = document.createElement('details');
  element.className = 'dn-strona__strefa dn-strona__strefa--zwijana';
  element.open = czyRozwiniete(opis.klucz, opis.domyslnieRozwiniete);

  const zapowiedz = document.createElement('summary');
  zapowiedz.className = 'dn-strona__naglowek-strefy dn-strona__zapowiedz';

  const napis = document.createElement('p');
  napis.className = 'dn-etykieta-wersalikowa dn-strona__etykieta-strefy';
  napis.textContent = opis.etykieta;

  // Dopisek stoi w wierszu etykiety: on niesie liczbę rozstrzygającą o rozwinięciu.
  const dopisek = document.createElement('span');
  dopisek.className = 'dn-strona__dopisek-strefy';
  dopisek.hidden = true;

  const wiersz = document.createElement('span');
  wiersz.className = 'dn-strona__wiersz-zapowiedzi';
  wiersz.append(napis, dopisek);

  const wyjasnienie = document.createElement('p');
  wyjasnienie.className = 'dn-strona__wyjasnienie-strefy';
  wyjasnienie.textContent = opis.wyjasnienie;

  zapowiedz.append(wiersz, wyjasnienie);

  const cialo = document.createElement('div');
  cialo.className = 'dn-strona__cialo-strefy';

  element.append(zapowiedz, cialo);

  element.addEventListener('toggle', () => {
    zapamietajZwiniecie(opis.klucz, element.open);
    for (const sluchacz of sluchacze) sluchacz(element.open);
  });

  return {
    element,
    cialo,

    ustawDopisek(tekst) {
      dopisek.textContent = tekst;
      dopisek.hidden = tekst === '';
    },

    naZmianePostaci(sluchacz) {
      sluchacze.push(sluchacz);
    },

    rozwiniete: () => element.open,
  };
}
