import { elementIkony, type NazwaIkony } from '../ikony/ikony';
import { NAZWY_TRAS, Trasa, TRASY } from './trasy';

/**
 * Przełącznik tras zamontowany na pasku aplikacji: grupa przycisków trzech
 * widoków najwyższego rzędu wraz z oznaczeniem trasy pokazywanej w tej chwili.
 */
export interface PrzelacznikTras {
  /** Grupa przycisków montowana na pasku aplikacji. */
  element: HTMLElement;
  /** Oznacza trasę pokazywaną w tej chwili. */
  ustawBiezaca(trasa: Trasa): void;
}

/**
 * Znak trasy w przełączniku — nazwa ikony z zestawu `ikony/` przypisana każdej
 * z trzech tras, ponieważ przyciski paska stoją bez napisu.
 */
const ZNAKI: Readonly<Record<Trasa, NazwaIkony>> = {
  [Trasa.StronaGlowna]: 'dom',
  [Trasa.Srodowisko]: 'folder',
  [Trasa.Pulpit]: 'zegar',
};

/**
 * Przełącznik trzech widoków najwyższego rzędu. Buduje przyciski tras i oznacza
 * trasę bieżącą atrybutem `aria-current`; sam nie przełącza — zgłasza wybór
 * warstwie, która trzyma router.
 */
export function utworzPrzelacznikTras(naWybor: (trasa: Trasa) => void): PrzelacznikTras {
  const element = document.createElement('nav');
  element.className = 'dn-trasy dn-trasy--znaki';
  element.setAttribute('aria-label', 'Widoki aplikacji');

  const przyciski = new Map<Trasa, HTMLButtonElement>();

  for (const trasa of TRASY) {
    const przycisk = document.createElement('button');
    przycisk.type = 'button';
    przycisk.className = 'dn-trasy__przycisk';
    przycisk.dataset.trasa = trasa;
    przycisk.title = NAZWY_TRAS[trasa];

    // Sam znak, bez napisu: nazwę widoku niosą `title` oraz etykieta dostępności przycisku.
    przycisk.setAttribute('aria-label', NAZWY_TRAS[trasa]);
    przycisk.append(elementIkony(ZNAKI[trasa], { rozmiar: 18 }));
    przycisk.addEventListener('click', () => naWybor(trasa));

    przyciski.set(trasa, przycisk);
    element.append(przycisk);
  }

  return {
    element,
    ustawBiezaca(trasa) {
      for (const [nazwa, przycisk] of przyciski) {
        if (nazwa === trasa) przycisk.setAttribute('aria-current', 'page');
        else przycisk.removeAttribute('aria-current');
      }
    },
  };
}
