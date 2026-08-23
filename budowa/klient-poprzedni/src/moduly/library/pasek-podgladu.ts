import { przycisk } from '../../modele/kontrolki-formularza';

/**
 * Pasek czynności okna File Preview.
 *
 * Jedna odpowiedzialność: kontrolki wiersza pod podglądem wraz z licznikiem
 * stron. Pasek nie wie, co robi każda z nich — dostaje czynności z zewnątrz,
 * więc okno może je zmienić bez dotykania układu.
 */
export interface CzynnosciPodgladu {
  wstecz(): void;
  dalej(): void;
  otworz(): void;
  zamknij(): void;
}

export interface PasekPodgladu {
  element: HTMLElement;
  /** Ustawia zdanie o stronie podglądu. */
  ustawStrone(tresc: string): void;
}

export function utworzPasekPodgladu(
  sterowanie: HTMLElement,
  czynnosci: CzynnosciPodgladu,
): PasekPodgladu {
  const strona = document.createElement('span');
  strona.className = 'ml-podglad__strona';

  const element = document.createElement('div');
  element.className = 'ml-podglad__pasek';
  element.append(
    strona,
    przyciskPaska('Poprzednia strona', 'strona-wstecz', czynnosci.wstecz),
    przyciskPaska('Następna strona', 'strona-dalej', czynnosci.dalej),
    sterowanie,
    przyciskPaska('Otwórz w module źródłowym', 'edycja', czynnosci.otworz, 'dn-btn--atrament'),
    przyciskPaska('Zamknij', 'zamknij', czynnosci.zamknij, 'dn-btn--zarys'),
  );

  return {
    element,
    ustawStrone: (tresc) => {
      strona.textContent = tresc;
    },
  };
}

/** Przycisk paska wraz ze znacznikiem czynności w `dataset`. */
function przyciskPaska(
  etykieta: string,
  kod: string,
  czynnosc: () => void,
  wariant = 'dn-btn--duch',
): HTMLElement {
  const element = przycisk(etykieta, `dn-btn dn-btn--sm ${wariant}`);
  element.dataset['czynnosc'] = kod;
  element.addEventListener('click', czynnosc);
  return element;
}
