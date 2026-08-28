import type { LibraryVersion } from '../../../../shared/contract';
import { poleWyboru } from '../../modele/kontrolki-formularza';

/**
 * Filtr historii wersji zawęża wykaz po sprawcy zmiany i po oznaczeniu, licząc
 * zawężenie w całości po stronie okna z pól odpowiedzi wykazu wersji, ponieważ
 * komenda odczytu nie przyjmuje pola zawężającego.
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

/**
 * Napisy rozpoznające sprawcę będącego Operatorem: pole sprawcy sprowadzone do
 * małych liter musi zawierać jeden z nich, bo rdzeń zapisuje tam napis
 * swobodny, a nie wartość wyliczenia kontraktu.
 */
const SPRAWCA_OPERATOR = ['operator', 'uzytkownik', 'użytkownik'];

/**
 * Napisy rozpoznające sprawcę będącego modelem albo modułem: dopasowanie idzie
 * tak samo po zawartości pola sprawcy, więc oznaczenie modelu wystarczy podać
 * w dowolnym otoczeniu wyrazowym.
 */
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

/**
 * Nazwa zakresu do zdania o pustym wyniku zawężenia: oddaje etykietę pozycji
 * wykazu zakresów, a dla kodu spoza wykazu oddaje sam kod, żeby zdanie miało
 * czym nazwać zawężenie.
 */
export function nazwaZakresu(zakres: ZakresHistorii): string {
  return ZAKRESY.find((pozycja) => pozycja.kod === zakres)?.etykieta ?? zakres;
}
