import { elementIkony } from '../ikony/ikony';
import type { Kanal } from '../protokol/kanal';
import { utworzSekcjeDostepow, type SekcjaDostepow } from './sekcja-dostepow';

/**
 * Okno dostępów i katalogu roboczego: rama na natywnym elemencie `dialog`,
 * która osadza sekcję dostępów, dokłada nagłówek i stopkę z zamknięciem,
 * a wykazy zleca sekcji dopiero przy otwarciu.
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

/**
 * Nagłówek okna składa ikonę kłódki, tytuł, rozpychacz i przycisk zamknięcia;
 * przycisk niesie opis dostępnościowy oraz podpowiedź, a wskazanie go wywołuje
 * przekazane domknięcie.
 */
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
