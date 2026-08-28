import type { OpisOkna } from '../okno-komunikacji/opis-okna';

/**
 * Wybiórka wskazująca kontener przygotowany przez `index.html`, do którego
 * gospodarz wstawia element aplikacji. Stała trzyma ją w jednym miejscu, dzięki
 * czemu wybiórka nie pada literałem w ciele funkcji szukającej kontenera.
 */
const KONTENER = '#danaco-console';

/**
 * Miejsce, w którym mieszka cała aplikacja. Jedyną odpowiedzialnością jest tu
 * przygotowanie dokumentu i jednego elementu, do którego router wstawia widoki
 * tras. Gospodarz nie zna ani rdzenia, ani żadnego widoku.
 */
export function utworzGospodarza(opis: OpisOkna): HTMLElement {
  document.title = `${opis.projekt} — Danaco Console`;

  const element = document.createElement('div');
  element.className = 'dn-aplikacja';

  miejsceMontazu().append(element);
  return element;
}

/**
 * Oddaje kontener wskazany wybiórką, a przy jego braku ciało dokumentu, dzięki
 * czemu brak kontenera w `index.html` nie zatrzymuje uruchomienia aplikacji.
 */
function miejsceMontazu(): HTMLElement {
  return document.querySelector<HTMLElement>(KONTENER) ?? document.body;
}
