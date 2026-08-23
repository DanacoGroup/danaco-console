/**
 * Dymek objaśnienia [?] przy elemencie konfiguracji sceny.
 *
 * Nośnik zasady „każdy element konfiguracji zawiera objaśnienie kontekstowe":
 * mały znak zapytania, nad którym komponent `.dn-tooltip` biblioteki
 * (`komponenty/drobne.css`) pokazuje treść przy najechaniu albo przy ognisku —
 * bez kliknięcia i bez zamykania.
 *
 * Znak jest przyciskiem, nie ozdobą: naciśnięcie prowadzi do niego ognisko,
 * a ognisko pokazuje objaśnienie — każde naciśnięcie daje odpowiedź.
 * Treść czytają technologie wspomagające z `aria-label` znaku, więc dymek nie
 * potrzebuje identyfikatora i nie zderza się między oknami sceny.
 *
 * Biblioteka niesie samą chmurkę (`.dn-tooltip` + `.dn-tooltip-tresc`), ale nie
 * niesie znaku [?] — każdy widok rysuje go u siebie. Klasy znaku są tu więc
 * widokowe.
 */
export function utworzObjasnienie(objasnienie: string): HTMLElement {
  const dymek = document.createElement('span');
  dymek.className = 'dn-tooltip dn-okna__objasnienie';

  const znak = document.createElement('button');
  znak.type = 'button';
  znak.className = 'dn-okna__objasnienie-znak';
  znak.textContent = '?';
  znak.setAttribute('aria-label', objasnienie);

  const tresc = document.createElement('span');
  tresc.className = 'dn-tooltip-tresc';
  tresc.setAttribute('aria-hidden', 'true');
  tresc.textContent = objasnienie;

  dymek.append(znak, tresc);
  return dymek;
}
