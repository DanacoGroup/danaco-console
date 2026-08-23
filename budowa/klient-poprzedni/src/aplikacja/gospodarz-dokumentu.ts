import type { OpisOkna } from '../okno-komunikacji/opis-okna';

/** Identyfikator kontenera przygotowanego przez `index.html`. */
const KONTENER = '#danaco-console';

/**
 * Miejsce, w którym mieszka cała aplikacja.
 *
 * Jedna odpowiedzialność: przygotowanie dokumentu i jednego elementu, do
 * którego router wstawia widoki tras. Gospodarz nie zna ani rdzenia, ani
 * żadnego widoku — nie wie nawet, ile tras ma aplikacja.
 *
 * Brak kontenera z `index.html` nie zatrzymuje uruchomienia — gospodarz trafia
 * wtedy do ciała dokumentu. Tytuł dokumentu składa nazwę projektu z opisu okna
 * z nazwą produktu, bo projekt pochodzi z konfiguracji budowania.
 */
export function utworzGospodarza(opis: OpisOkna): HTMLElement {
  document.title = `${opis.projekt} — Danaco Console`;

  const element = document.createElement('div');
  element.className = 'dn-aplikacja';

  miejsceMontazu().append(element);
  return element;
}

/** Kontener z dokumentu, a przy jego braku ciało dokumentu. */
function miejsceMontazu(): HTMLElement {
  return document.querySelector<HTMLElement>(KONTENER) ?? document.body;
}
