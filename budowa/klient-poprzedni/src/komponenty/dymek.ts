/**
 * Dymek objaśnienia — biblioteczna fabryka objaśnienia kontekstowego przy
 * elemencie sterowania: znak zapytania, nad którym `.dn-tooltip`
 * z `komponenty/drobne.css` pokazuje treść przy najechaniu albo przy ognisku,
 * bez kliknięcia.
 */

/**
 * Klasy okna wywołującego godzą jedną fabrykę z wyglądem właściwym dla miejsca
 * wywołania: wywołujący z własnym arkuszem podaje swoje klasy, a wywołujący bez
 * arkusza nie podaje nic i dostaje wygląd biblioteczny.
 */
export interface KlasyDymka {
  /** Klasa dokładana do bibliotecznego `.dn-tooltip` na powłoce dymka. */
  powloka?: string;
  /** Klasa znaku [?], która zastępuje biblioteczne `.dn-btn-ikona`. */
  znak?: string;
}

export function utworzDymekObjasnienia(
  objasnienie: string,
  klasy: KlasyDymka = {},
): HTMLElement {
  const dymek = document.createElement('span');
  dymek.className = 'dn-tooltip';
  if (klasy.powloka !== undefined && klasy.powloka !== '') {
    dymek.classList.add(klasy.powloka);
  }

  const znak = document.createElement('button');
  znak.type = 'button';
  znak.className = klasy.znak !== undefined && klasy.znak !== '' ? klasy.znak : 'dn-btn-ikona';
  znak.textContent = '?';
  znak.setAttribute('aria-label', objasnienie);

  const tresc = document.createElement('span');
  tresc.className = 'dn-tooltip-tresc';
  tresc.setAttribute('aria-hidden', 'true');
  tresc.textContent = objasnienie;

  dymek.append(znak, tresc);
  return dymek;
}
