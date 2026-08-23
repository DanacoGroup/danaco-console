import type { LibraryVersion } from '../../../../shared/contract';
import { poleWyboru } from '../../modele/kontrolki-formularza';

/**
 * Filtr historii wersji — zawężenie wykazu po sprawcy zmiany i po oznaczeniu.
 *
 * Filtr działa w całości po stronie okna: `library.version.list` nie przyjmuje
 * pola zawężającego, a odpowiedź niesie komplet pól, po których dokumentacja
 * każe zawężać (`author`, `label`). Wysyłanie w tym celu drugiego odczytu
 * byłoby pytaniem o to, co okno już ma.
 *
 * Sprawca jest napisem, nie wyliczeniem kontraktu: rdzeń zapisuje w polu
 * `author` to, co poda wołający (`Operator albo modul`). Pozycje „zmiany
 * Operatora" i „zmiany modelu" dopasowują więc po zawartości napisu, a wersja
 * o sprawcy nienazwanym wchodzi wyłącznie do pozycji „wszystkie" — zgadywanie
 * po stronie okna kazałoby jej trafić do jednej z dwóch grup bez podstawy.
 */
export type ZakresHistorii = 'wszystkie' | 'operator' | 'model' | 'oznaczone';

export interface FiltrWersji {
  element: HTMLElement;
  /** Zakres bieżący. */
  zakres(): ZakresHistorii;
  /** Wersje po zawężeniu. */
  zastosuj(wersje: readonly LibraryVersion[]): readonly LibraryVersion[];
}

const ZAKRESY: ReadonlyArray<{ kod: ZakresHistorii; etykieta: string }> = [
  { kod: 'wszystkie', etykieta: 'Wszystkie wersje' },
  { kod: 'operator', etykieta: 'Tylko zmiany Operatora' },
  { kod: 'model', etykieta: 'Tylko zmiany modelu' },
  { kod: 'oznaczone', etykieta: 'Tylko wersje z etykietą' },
];

/** Napisy, po których poznaje się sprawcę będącego Operatorem. */
const SPRAWCA_OPERATOR = ['operator', 'uzytkownik', 'użytkownik'];

/** Napisy, po których poznaje się sprawcę będącego modelem albo modułem. */
const SPRAWCA_MODEL = ['model', 'agent', 'modul', 'moduł'];

export function utworzFiltrWersji(naZmiane: () => void): FiltrWersji {
  const pole = poleWyboru(
    {
      etykieta: 'Filtr historii',
      opis:
        'Zawężenie liczone w oknie z odpowiedzi library.version.list — komenda nie przyjmuje ' +
        'pola filtru, a odpowiedź niesie sprawcę i etykietę każdej wersji.',
    },
    ZAKRESY.map((pozycja) => ({ wartosc: pozycja.kod, etykieta: pozycja.etykieta })),
  );
  pole.kontrolka.dataset['ster'] = 'filtr-historii';
  pole.kontrolka.addEventListener('change', () => naZmiane());

  function zakres(): ZakresHistorii {
    const pozycja = ZAKRESY.find((wpis) => wpis.kod === pole.kontrolka.value);
    return pozycja === undefined ? 'wszystkie' : pozycja.kod;
  }

  return {
    element: pole.element,
    zakres,

    zastosuj(wersje) {
      const wybrany = zakres();
      if (wybrany === 'wszystkie') return wersje;
      if (wybrany === 'oznaczone') {
        return wersje.filter((wersja) => (wersja.label ?? '').trim() !== '');
      }
      const wzorce = wybrany === 'operator' ? SPRAWCA_OPERATOR : SPRAWCA_MODEL;
      return wersje.filter((wersja) => {
        const sprawca = (wersja.author ?? '').trim().toLowerCase();
        return sprawca !== '' && wzorce.some((wzorzec) => sprawca.includes(wzorzec));
      });
    },
  };
}

/** Nazwa zakresu do zdania o pustym wyniku zawężenia. */
export function nazwaZakresu(zakres: ZakresHistorii): string {
  return ZAKRESY.find((pozycja) => pozycja.kod === zakres)?.etykieta ?? zakres;
}
