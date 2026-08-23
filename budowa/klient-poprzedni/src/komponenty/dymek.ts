/**
 * Dymek objaśnienia [?] — biblioteczna fabryka objaśnienia kontekstowego przy
 * elemencie sterowania: mały znak zapytania, nad którym `.dn-tooltip`
 * z `komponenty/drobne.css` pokazuje treść przy najechaniu albo przy ognisku —
 * bez kliknięcia i bez zamykania.
 *
 * Znak jest przyciskiem, nie ozdobą: naciśnięcie prowadzi do niego ognisko,
 * a ognisko pokazuje objaśnienie. Treść czytają technologie wspomagające
 * z `aria-label` znaku, więc dymek nie potrzebuje identyfikatora i nie zderza
 * się między oknami.
 */

/**
 * Klasy okna wywołującego — sposób na pogodzenie jednej fabryki z wyglądem
 * właściwym dla miejsca wywołania.
 *
 * `.dn-btn-ikona` to kwadrat 32 px (na dotyku 40 px) o narożniku `--dn-r-sm`,
 * z przezroczystym obrysem i wskaźnikiem `pointer`; znaki modułów to obwódki
 * 14–16 px o narożniku `--dn-r-pill`, z widocznym obrysem i wskaźnikiem `help`.
 * Powłoka niesie z kolei zachowanie w rzędzie (`flex: none`, odstęp,
 * `vertical-align`), którego `.dn-tooltip` nie zna.
 *
 * Wywołujący, który ma własny arkusz, podaje swoje klasy; wywołujący bez
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
