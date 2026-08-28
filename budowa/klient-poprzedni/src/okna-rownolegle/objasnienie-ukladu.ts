/**
 * Moduł tworzy przycisk objaśnienia przy elemencie konfiguracji sceny, pokazujący treść pomocniczą w chmurce po najechaniu albo ustawieniu ogniska.
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
