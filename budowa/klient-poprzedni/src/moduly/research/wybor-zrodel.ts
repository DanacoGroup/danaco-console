import type { ResearchSource } from '../../../../shared/contract';
import type { ZaznaczeniePozycji } from './zaznaczenie-pozycji';

/**
 * Wybór źródeł powiązanych z ustaleniem: pole wielokrotnego wyboru zasilane katalogiem
 * Sources Manager, jedzie w polu sourceIds komendy zapisu ustalenia.
 */
export interface WyborZrodel {
  element: HTMLElement;
  /** Przerysowuje wybór po zmianie katalogu źródeł. */
  odswiez(zrodla: readonly ResearchSource[]): void;
}

export function utworzWyborZrodel(
  zaznaczenie: ZaznaczeniePozycji,
  naZmiane: () => void,
): WyborZrodel {
  const element = document.createElement('fieldset');
  element.className = 'mr-wybor-zrodel';

  const opis = document.createElement('legend');
  opis.className = 'dn-pole-etykieta mr-wybor-zrodel__opis';
  opis.textContent = 'Źródła powiązane z ustaleniem';

  const lista = document.createElement('div');
  lista.className = 'mr-wybor-zrodel__lista';

  element.append(opis, lista);

  return {
    element,

    odswiez(zrodla) {
      if (zrodla.length === 0) {
        const pusto = document.createElement('p');
        pusto.className = 'dn-pole-opis mr-wybor-zrodel__pusto';
        pusto.textContent =
          'Katalog źródeł jest pusty — skataloguj źródło w Sources Manager, aby móc je tu powiązać.';
        lista.replaceChildren(pusto);
        return;
      }
      lista.replaceChildren(
        ...zrodla.map((zrodlo) => pozycja(zrodlo, zaznaczenie, naZmiane)),
      );
    },
  };
}

/** Jedno źródło jako pole wyboru wraz z etykietą tytułu, zaznaczane w powiązaniu ze wskazanym ustaleniem. */
function pozycja(
  zrodlo: ResearchSource,
  zaznaczenie: ZaznaczeniePozycji,
  naZmiane: () => void,
): HTMLElement {
  const element = document.createElement('label');
  element.className = 'mr-wybor-zrodel__pozycja';

  const wybor = document.createElement('input');
  wybor.type = 'checkbox';
  wybor.className = 'dn-check';
  wybor.checked = zaznaczenie.czyWybrana(zrodlo.id);
  wybor.addEventListener('change', () => {
    zaznaczenie.przelacz(zrodlo.id);
    naZmiane();
  });

  const nazwa = document.createElement('span');
  nazwa.textContent = zrodlo.title;

  element.append(wybor, nazwa);
  return element;
}
