import { czyRozwiniete, zapamietajZwiniecie } from './pamiec-zwiniecia';

/**
 * Strefa przywoływana — zapowiedź stoi zawsze, treść dopiero na życzenie.
 *
 * Buduje `<details>` z `<summary>`: postać rozwinięcia trzyma przeglądarka,
 * a dane strefy czytane są dopiero przy rozwinięciu. Tak samo działa archiwum
 * sesji (`strefa-archiwum.ts`); ten plik uogólnia zabieg na całe strefy
 * i dokłada pamięć postaci między wejściami, żeby przywołanie nie powtarzało
 * się przy każdym wejściu na stronę.
 *
 * Strefa zwinięta nie jest strefą ukrytą: zapowiedź — nazwa, zdanie
 * wyjaśnienia i dopisek z liczbą — stoi na ekranie zawsze i mówi, co jest pod
 * spodem, więc jedno naciśnięcie otwiera treść znaną z opisu.
 *
 * Rozwinięcia nie dubluje żadna klasa CSS; jedynym jego nośnikiem jest atrybut
 * `open`, bo druga kopia stanu mogłaby się z pierwszą wyłącznie rozminąć.
 */

export interface OpisStrefyZwijanej {
  /** Etykieta wersalikowa — nazwa strefy. */
  etykieta: string;
  /** Jedno zdanie: po co strefa stoi na ekranie. */
  wyjasnienie: string;
  /**
   * Klucz pamięci postaci, stały między wejściami i między wydaniami.
   * Zmiana klucza znaczy „zacznij od domyślnej", więc zmienia się go świadomie.
   */
  klucz: string;
  /** Postać przed pierwszym rozstrzygnięciem Operatora. */
  domyslnieRozwiniete: boolean;
}

export interface StrefaZwijana {
  /** Element `<details>` do osadzenia w stronie. */
  element: HTMLDetailsElement;
  /** Miejsce na treść strefy — rysowane niezależnie od postaci. */
  cialo: HTMLElement;
  /**
   * Dopisek zapowiedzi — liczba pod spodem, żeby zwinięcie nie było ślepe.
   * Napis pusty zdejmuje dopisek zamiast zostawiać puste miejsce.
   */
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

  // Dopisek stoi w wierszu etykiety, nie pod wyjaśnieniem: to on niesie liczbę,
  // dla której Operator zdecyduje, czy strefę w ogóle rozwijać.
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
