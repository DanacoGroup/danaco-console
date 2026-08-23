import { elementIkony } from '../ikony/ikony';
import type { Kanal } from '../protokol/kanal';
import { przycisk } from './kontrolki-formularza';
import { utworzSekcjeModeli, type SekcjaModeli } from './sekcja-modeli';

/**
 * Okno sekcji modeli — rama dla czterech obszarów: kont, ustawień osi, zasad
 * i tożsamości oraz podglądu złożonego promptu.
 *
 * Okno stoi na natywnym `<dialog>`: warstwę tła, stos okien i zamknięcie
 * klawiszem Esc daje przeglądarka, a nie własna nakładka. Wygląd bierze
 * z biblioteki `komponenty/` (`dn-modal`), więc plik nie ustala barw.
 *
 * Rama wyłącznie osadza sekcję — treść i stan mieszkają w `sekcja-modeli`,
 * żeby tę samą sekcję dało się wstawić także w widok osadzony bez powielania.
 *
 * Okno otwiera się natychmiast, przed odpowiedzią rdzenia.
 */
export interface OknoModeli {
  /** Element `<dialog>` osadzany w dokumencie. */
  element: HTMLDialogElement;
  /** Otwiera okno i zleca odczyt rejestrów oraz katalogów. */
  otworz(): void;
  /** Zamyka okno; stan i subskrypcje zostają. */
  zamknij(): void;
  /** Odłącza subskrypcje kanału i usuwa okno z dokumentu. */
  rozlacz(): void;
}

export function utworzOknoModeli(kanal: Kanal): OknoModeli {
  const sekcja: SekcjaModeli = utworzSekcjeModeli(kanal);

  const element = document.createElement('dialog');
  element.className = 'dn-modal dm-okno';
  element.setAttribute('aria-label', 'Modele, konta i tożsamość');

  const cialo = document.createElement('div');
  cialo.className = 'dm-okno__cialo';
  cialo.append(sekcja.element);

  const odswiez = przycisk('Odczytaj rejestry ponownie', 'dn-btn dn-btn--zarys');
  const zamknij = przycisk('Zamknij', 'dn-btn dn-btn--atrament');

  const stopka = document.createElement('footer');
  stopka.className = 'dn-modal-stopka dm-okno__stopka';
  stopka.append(odswiez, zamknij);

  element.append(naglowek(() => element.close()), cialo, stopka);

  odswiez.addEventListener('click', () => sekcja.odswiez());
  zamknij.addEventListener('click', () => element.close());

  return {
    element,

    otworz() {
      if (!element.isConnected) document.body.append(element);
      if (!element.open) element.showModal();
      sekcja.odswiez();
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
  element.className = 'dn-modal-naglowek dm-okno__naglowek';

  const tytul = document.createElement('h2');
  tytul.className = 'dn-modal-tytul';
  tytul.textContent = 'Modele, konta i tożsamość';

  const rozpychacz = document.createElement('span');
  rozpychacz.className = 'dm-okno__rozpychacz';

  const zamkniecie = document.createElement('button');
  zamkniecie.type = 'button';
  zamkniecie.className = 'dn-btn-ikona';
  zamkniecie.setAttribute('aria-label', 'Zamknij okno modeli');
  zamkniecie.title = 'Zamknij okno modeli';
  zamkniecie.append(elementIkony('zamknij', { rozmiar: 18 }));
  zamkniecie.addEventListener('click', naZamkniecie);

  element.append(elementIkony('agent', { rozmiar: 20 }), tytul, rozpychacz, zamkniecie);
  return element;
}
