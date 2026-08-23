import { elementIkony } from '../ikony/ikony';
import { POZYCJA_KONFIGURACJI, type PozycjaUstawienia } from '../strona-glowna/pozycje-ustawien';
import { utworzGodlo } from './godlo-aplikacji';
import { utworzMenuAplikacji, type MenuAplikacji } from './menu-aplikacji';
import { utworzMenuOperatora, type MenuOperatora } from './menu-operatora';
import { utworzPrzelacznikMotywu } from './przelacznik-motywu';
import { utworzPrzelacznikTras, type PrzelacznikTras } from './przelacznik-tras';
import type { Trasa } from './trasy';

/** Pasek aplikacji nad widokiem trasy. */
export interface PasekAplikacji {
  /** Pasek montowany jako pierwszy wiersz widoku. */
  element: HTMLElement;
  /** Przełącznik tras — oznaczenie trasy bieżącej. */
  trasy: PrzelacznikTras;
  /** Zdejmuje nasłuchy dokumentu założone przez menu. Obowiązkowe. */
  zamknij(): void;
}

/** Zależności paska: nazwa projektu i sposób przejścia na inną trasę. */
export interface OpcjePaskaAplikacji {
  /** Nazwa projektu obok godła; pochodzi z opisu okna. */
  projekt: string;
  /** Przejście na wskazaną trasę. */
  naTrase(trasa: Trasa): void;
  /**
   * Wybór pozycji ustawień z menu aplikacji albo z menu Operatora.
   *
   * Pasek nie wie, co pozycja otwiera — wykaz skutków mieszka w
   * `akcje-ustawien.ts` i jest ten sam, którego używa listwa strony głównej.
   */
  naUstawienie(pozycja: PozycjaUstawienia): void;
}

/**
 * Pasek aplikacji — tożsamość, przełącznik widoków, motyw.
 *
 * Jedna odpowiedzialność: złożenie trzech elementów paska. Pasek stoi nad
 * widokami, które własnego paska nie mają: Centrum dowodzenia i Mission
 * Control. Powłoka środowiska ma pasek własny (`powloka/pasek-gorny.ts`)
 * i drugiego nie dostaje — przełącznik tras oraz przełącznik motywu wchodzą
 * wtedy w jej grupę akcji.
 *
 * Motyw jest przełączalny z każdego widoku, bo pasek niesie ten sam
 * przełącznik, którego używa powłoka. Wartości motywu ani żadnej barwy pasek
 * nie zna — całą pracę wykonuje warstwa `motyw/`.
 *
 * Poza godłem, przełącznikiem tras i motywem pasek niesie menu aplikacji przy
 * godle, ikonę ustawień i menu Operatora w grupie akcji — te same wejścia,
 * które ma powłoka środowiska (`powloka/akcje-paska.ts`), żeby produkt
 * zachowywał się jednakowo na każdej trasie. Żadne z nich nie zakłada nowej
 * czynności ani nie woła komendy spoza wykazu: pozycje są te same, co w listwie
 * strony głównej, a skutek jeden i wspólny.
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

  // Trasy stoją w grupie akcji, nie osobno pośrodku paska: znaki tras
  // i kontrolki platformy tworzą jeden rząd ikon przy prawej krawędzi,
  // rozdzielony kreską. Dzięki temu pasek ma dwie grupy, nie trzy.
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
 * Ikona ustawień — skrót do Okna konfiguracji.
 *
 * Skrót, a nie druga droga: otwiera dokładnie tę pozycję, którą niesie listwa
 * strony głównej pod kodem `konfiguracja`, przez ten sam wykaz skutków
 * (`akcje-ustawien.ts`). Stoi w pasku, bo ustawienia platformy wykonuje się
 * z każdego miejsca, a listwa widoczna jest tylko na stronie głównej.
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
