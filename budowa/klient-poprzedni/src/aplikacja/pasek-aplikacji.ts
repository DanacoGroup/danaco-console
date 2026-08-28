import { elementIkony } from '../ikony/ikony';
import { POZYCJA_KONFIGURACJI, type PozycjaUstawienia } from '../strona-glowna/pozycje-ustawien';
import { utworzGodlo } from './godlo-aplikacji';
import { utworzMenuAplikacji, type MenuAplikacji } from './menu-aplikacji';
import { utworzMenuOperatora, type MenuOperatora } from './menu-operatora';
import { utworzPrzelacznikMotywu } from './przelacznik-motywu';
import { utworzPrzelacznikTras, type PrzelacznikTras } from './przelacznik-tras';
import type { Trasa } from './trasy';

/**
 * Pasek aplikacji stoi nad widokiem trasy i niesie tożsamość projektu,
 * przełącznik tras oraz przełącznik motywu. Zamknięcie paska zdejmuje nasłuchy
 * dokumentu założone przez menu i jest obowiązkowe przy zdjęciu widoku.
 */
export interface PasekAplikacji {
  /** Pasek montowany jako pierwszy wiersz widoku. */
  element: HTMLElement;
  /** Przełącznik tras — oznaczenie trasy bieżącej. */
  trasy: PrzelacznikTras;
  /** Zdejmuje nasłuchy dokumentu założone przez menu. Obowiązkowe. */
  zamknij(): void;
}

/**
 * Zależności paska obejmują nazwę projektu wyświetlaną obok godła oraz obsługę
 * przejścia na wskazaną trasę i wyboru pozycji ustawień. Pasek nie zna skutku
 * żadnej z tych czynności; wykonuje je warstwa, która pasek buduje.
 */
export interface OpcjePaskaAplikacji {
  /** Nazwa projektu obok godła; pochodzi z opisu okna. */
  projekt: string;
  /** Przejście na wskazaną trasę. */
  naTrase(trasa: Trasa): void;
  /** Wybór pozycji ustawień z menu aplikacji albo z menu Operatora. */
  naUstawienie(pozycja: PozycjaUstawienia): void;
}

/**
 * Składa pasek z tożsamości projektu, przełącznika tras, przełącznika motywu,
 * menu aplikacji, ikony ustawień oraz menu Operatora. Pasek stoi nad widokami,
 * które własnego paska nie mają, a barw nie zna — wykonuje je warstwa `motyw/`.
 */
export function utworzPasekAplikacji(opcje: OpcjePaskaAplikacji): PasekAplikacji {
  const element = document.createElement('header');
  element.className = 'dn-pasek dn-pasek-aplikacji';

  const rozpychacz = document.createElement('span');
  rozpychacz.className = 'dn-powloka__rozpychacz';

  const trasy = utworzPrzelacznikTras(opcje.naTrase);
  const motyw = utworzPrzelacznikMotywu();

  const menu: MenuAplikacji = utworzMenuAplikacji({ naPozycje: opcje.naUstawienie });
  const operator: MenuOperatora = utworzMenuOperatora({ naPozycje: opcje.naUstawienie });
  const ustawienia = przyciskUstawien(opcje.naUstawienie);

  const tozsamosc = document.createElement('div');
  tozsamosc.className = 'dn-pasek-aplikacji__tozsamosc';
  tozsamosc.append(menu.element, utworzGodlo(opcje.projekt));

  // Trasy stoją w grupie akcji, więc pasek ma dwie grupy zamiast trzech.
  const rozdzielacz = document.createElement('span');
  rozdzielacz.className = 'dn-pasek-gorny__rozdzielacz';

  const akcje = document.createElement('div');
  akcje.className = 'dn-pasek-prawa dn-powloka__akcje';
  akcje.append(trasy.element, rozdzielacz, ustawienia, motyw.element, operator.element);

  element.append(tozsamosc, rozpychacz, akcje);

  return {
    element,
    trasy,
    zamknij() {
      menu.zamknij();
      operator.zamknij();
    },
  };
}

/**
 * Ikona ustawień jest skrótem do Okna konfiguracji: otwiera tę samą pozycję,
 * którą niesie listwa strony głównej, przez wspólny wykaz skutków. Stoi
 * w pasku, ponieważ listwa widoczna jest wyłącznie na stronie głównej.
 */
function przyciskUstawien(naPozycje: (pozycja: PozycjaUstawienia) => void): HTMLButtonElement {
  const pozycja = POZYCJA_KONFIGURACJI;
  const element = document.createElement('button');
  element.type = 'button';
  element.className = 'dn-btn-ikona dn-pasek-aplikacji__ustawienia';
  element.setAttribute('aria-label', pozycja.nazwa);
  element.title = pozycja.wyjasnienie;
  element.append(elementIkony('ustawienia', { rozmiar: 18 }));
  element.addEventListener('click', () => naPozycje(pozycja));
  return element;
}
