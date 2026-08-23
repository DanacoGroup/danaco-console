import { elementIkony } from '../ikony/ikony';
import type { Kanal } from '../protokol/kanal';
import { utworzSekcjeDostepow, type SekcjaDostepow } from './sekcja-dostepow';

/**
 * Okno, w którym mieszka sekcja dostępów — rama, nie treść.
 *
 * Sekcja sama nie zakłada, gdzie zostanie osadzona: oddaje element. Rama
 * potrzebna jest jednak od razu, bo pierwszym gospodarzem sekcji jest listwa
 * ustawień Centrum dowodzenia, a ta otwiera widoki jako okna modalne — tak jak
 * okno konfiguracji.
 *
 * Okno stoi na natywnym `<dialog>`, więc warstwę tła, stos okien i zamknięcie
 * klawiszem Esc daje przeglądarka, a nie własna nakładka. Wygląd bierze
 * z biblioteki `komponenty/` (`dn-modal`) — plik nie zna ani jednej barwy.
 *
 * Okno otwiera się natychmiast, przed odpowiedzią rdzenia; wykazy dojeżdżają
 * do niego odpowiedzią.
 */
export interface OknoDostepow {
  /** Element `<dialog>` osadzany w dokumencie. */
  element: HTMLDialogElement;
  /** Sekcja osadzona w oknie — przez nią idzie wiązanie z oknem rozmowy. */
  sekcja: SekcjaDostepow;
  /** Otwiera okno i wczytuje wykazy. */
  otworz(): void;
  /** Zamyka okno; stan i subskrypcje zostają. */
  zamknij(): void;
  /** Odłącza subskrypcje kanału i usuwa okno z dokumentu. */
  rozlacz(): void;
}

export function utworzOknoDostepow(kanal: Kanal, oknoID = ''): OknoDostepow {
  const sekcja = utworzSekcjeDostepow({ kanal, oknoID });

  const element = document.createElement('dialog');
  element.className = 'dn-modal dd-okno';
  element.setAttribute('aria-label', 'Dostępy i katalog roboczy');

  const cialo = document.createElement('div');
  cialo.className = 'dn-modal-cialo dd-okno__cialo';
  cialo.append(sekcja.element);

  const stopka = document.createElement('footer');
  stopka.className = 'dn-modal-stopka';

  const zamknij = document.createElement('button');
  zamknij.type = 'button';
  zamknij.className = 'dn-btn dn-btn--atrament';
  zamknij.textContent = 'Zamknij';
  zamknij.addEventListener('click', () => element.close());
  stopka.append(zamknij);

  element.append(naglowek(() => element.close()), cialo, stopka);

  return {
    element,
    sekcja,

    otworz() {
      if (!element.isConnected) document.body.append(element);
      if (!element.open) element.showModal();
      sekcja.wczytaj();
    },

    zamknij: () => element.close(),

    rozlacz() {
      sekcja.rozlacz();
      element.remove();
    },
  };
}

/** Nagłówek okna: ikona, tytuł, przycisk zamknięcia. */
function naglowek(naZamkniecie: () => void): HTMLElement {
  const element = document.createElement('header');
  element.className = 'dn-modal-naglowek dd-okno__naglowek';

  const tytul = document.createElement('h2');
  tytul.className = 'dn-modal-tytul';
  tytul.textContent = 'Dostępy i katalog roboczy';

  const rozpychacz = document.createElement('span');
  rozpychacz.className = 'dd-okno__rozpychacz';

  const zamkniecie = document.createElement('button');
  zamkniecie.type = 'button';
  zamkniecie.className = 'dn-btn-ikona';
  zamkniecie.setAttribute('aria-label', 'Zamknij okno dostępów');
  zamkniecie.title = 'Zamknij okno dostępów';
  zamkniecie.append(elementIkony('zamknij', { rozmiar: 18 }));
  zamkniecie.addEventListener('click', naZamkniecie);

  element.append(elementIkony('klodka', { rozmiar: 20 }), tytul, rozpychacz, zamkniecie);
  return element;
}
