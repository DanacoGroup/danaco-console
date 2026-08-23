import { elementIkony, type NazwaIkony } from '../ikony/ikony';
import { NAZWY_TRAS, Trasa, TRASY } from './trasy';

/** Przełącznik tras zamontowany na pasku. */
export interface PrzelacznikTras {
  /** Grupa przycisków montowana na pasku aplikacji. */
  element: HTMLElement;
  /** Oznacza trasę pokazywaną w tej chwili. */
  ustawBiezaca(trasa: Trasa): void;
}

/** Znak trasy w przełączniku; ikony z zestawu `ikony/`. */
const ZNAKI: Readonly<Record<Trasa, NazwaIkony>> = {
  [Trasa.StronaGlowna]: 'dom',
  [Trasa.Srodowisko]: 'folder',
  [Trasa.Pulpit]: 'zegar',
};

/**
 * Przełącznik trzech widoków najwyższego rzędu.
 *
 * Buduje przyciski tras i oznacza trasę bieżącą. Sam nie przełącza — zgłasza
 * wybór warstwie, która trzyma router.
 *
 * Żaden przycisk nie jest bramą i żaden nie zostaje wyszarzony: do środowiska
 * można wejść także wprost, bo środowisko domyślne istnieje od pierwszej
 * chwili. Trasa bieżąca jest oznaczona przez `aria-current`, nie przez
 * odebranie klikalności.
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

    // Sam znak, bez napisu: trasy stoją w grupie akcji paska obok ustawień,
    // motywu i menu Operatora, które również są ikonami, a nazwa widoku stoi
    // tuż pod paskiem jako tytuł strony. Nazwę niesie `title` i etykieta
    // dostępności przycisku.
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
