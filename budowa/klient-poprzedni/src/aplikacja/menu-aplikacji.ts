import { elementIkony } from '../ikony/ikony';
import { utworzMenuRozwijane, type MenuRozwijane } from '../okna-rownolegle/menu-rozwijane';
import { POZYCJE_USTAWIEN, type PozycjaUstawienia } from '../strona-glowna/pozycje-ustawien';

/**
 * Menu aplikacji w pasku — wejście do ustawień platformy z każdej trasy,
 * która ten pasek nosi (Centrum dowodzenia, Mission Control).
 *
 * Pozycje bierze z `POZYCJE_USTAWIEN` — tego samego wykazu, z którego powstaje
 * listwa ustawień strony głównej; drugi wykaz rozjechałby się przy dopisaniu
 * pozycji. Skutek naciśnięcia również jest wspólny: `wykonajZamiarUstawien`
 * z `akcje-ustawien.ts`. Menu nie wie, co się po naciśnięciu stanie.
 *
 * Rozwijanie niesie `okna-rownolegle/menu-rozwijane.ts` wraz z obsługą
 * klawiatury, `aria-haspopup`, zamykaniem na `pointerdown` poza obszarem
 * i warstwą `--dn-z-przybornik`.
 *
 * Menu nie zastępuje listwy strony głównej — daje te same pozycje tam, gdzie
 * listwy nie ma: na Mission Control i po zwinięciu strony w dół.
 */
export interface MenuAplikacji {
  element: HTMLElement;
  /** Zdejmuje nasłuchy dokumentu. Obowiązkowe przy zejściu gospodarza. */
  zamknij(): void;
}

export interface OpcjeMenuAplikacji {
  /** Skutek wyboru pozycji — wykaz skutków mieszka w `akcje-ustawien`. */
  naPozycje(pozycja: PozycjaUstawienia): void;
}

export function utworzMenuAplikacji(opcje: OpcjeMenuAplikacji): MenuAplikacji {
  const menu: MenuRozwijane = utworzMenuRozwijane({
    ikona: 'menu',
    etykieta: 'Menu aplikacji',
  });

  menu.ustawTresc(POZYCJE_USTAWIEN.map((pozycja) => wpis(pozycja, opcje.naPozycje, menu)));

  return {
    element: menu.element,
    zamknij: () => menu.zamknij(),
  };
}

/**
 * Jedna pozycja menu.
 *
 * Menu zamyka się po wyborze: pozycja otwiera okno albo modal, więc menu
 * zostawione otwarte wisiałoby nad tym, co właśnie otworzyło.
 */
function wpis(
  pozycja: PozycjaUstawienia,
  naPozycje: (pozycja: PozycjaUstawienia) => void,
  menu: MenuRozwijane,
): HTMLElement {
  const element = document.createElement('button');
  element.type = 'button';
  element.className = 'dn-menu__pozycja';
  element.setAttribute('role', 'menuitem');
  element.dataset['ustawienie'] = pozycja.kod;
  // Rozwinięcie nazwy idzie do technologii wspomagających i do dymka —
  // etykieta w menu musi zostać krótka, żeby lista dała się przebiec wzrokiem.
  element.title = pozycja.wyjasnienie;
  element.setAttribute('aria-description', pozycja.wyjasnienie);

  const ikona = document.createElement('span');
  ikona.className = 'dn-menu-aplikacji__ikona';
  ikona.append(elementIkony(pozycja.ikona, { rozmiar: 16 }));

  const nazwa = document.createElement('span');
  nazwa.className = 'dn-menu-aplikacji__nazwa';
  nazwa.textContent = pozycja.nazwa;

  element.append(ikona, nazwa);
  element.addEventListener('click', () => {
    menu.ustaw(false);
    naPozycje(pozycja);
  });
  return element;
}
