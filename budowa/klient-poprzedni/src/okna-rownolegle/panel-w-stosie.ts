import './panele.css';

import type { PanelPomocniczy } from '../okna-pomocnicze/panel-pomocniczy';
import { utworzNaglowekPanelu, type PozycjaNaglowka } from './naglowek-panelu';

/**
 * Obudowa jednej pozycji stosu paneli obok rozmowy niesie gęsty nagłówek nad treścią panelu i dwie czynności przy nim, oddzielając wnętrze od obudowy, żeby przeniesienie na pełny ekran nie było kopią panelu.
 */
export interface PanelWStosie {
  /** Kod pozycji — ten sam, którym woła się wytwórnię. */
  kod: string;
  /** Obudowa osadzana w kolumnie paneli. */
  element: HTMLElement;
  /** Treść panelu — do przeniesienia do widoku pełnoekranowego i z powrotem. */
  tresc: HTMLElement;
  odswiez(): void;
  zamknij(): void;
}

export interface OpcjePanelaWStosie {
  kod: string;
  /** Nazwa panelu w nagłówku; skracana wielokropkiem, nigdy cięta w kodzie. */
  tytul: string;
  panel: PanelPomocniczy;
  /** Operator nacisnął zamknięcie — gospodarz zdejmuje panel ze stosu. */
  naZamkniecie(): void;
  /** Operator nacisnął pełny ekran — gospodarz przenosi `tresc` na scenę. */
  naPelnyEkran(): void;
  /** Gotowe elementy rzędu ikon — miejsce na własne menu panelu, dokładane przez gospodarza. */
  dodatkiNaglowka?: readonly HTMLElement[];
}

export function utworzPanelWStosie(opcje: OpcjePanelaWStosie): PanelWStosie {
  const tresc = document.createElement('div');
  tresc.className = 'dn-okna__panel-tresc';
  tresc.append(opcje.panel.element);

  // Rząd mieści menu własne panelu, potem dwie czynności obudowy, bez trzeciej, by uniknąć atrapy.
  const czynnosci: readonly PozycjaNaglowka[] = [
    ...(opcje.dodatkiNaglowka ?? []),
    {
      ikona: 'link-zewnetrzny',
      etykieta: `Pokaż panel ${opcje.tytul} na całej scenie`,
      dzialanie: opcje.naPelnyEkran,
    },
    {
      ikona: 'zamknij',
      etykieta: `Zamknij panel ${opcje.tytul}`,
      dzialanie: opcje.naZamkniecie,
    },
  ];

  const element = document.createElement('section');
  element.className = 'dn-okna__panel';
  element.dataset['panel'] = opcje.kod;
  element.setAttribute('aria-label', `Panel ${opcje.tytul}`);
  element.append(utworzNaglowekPanelu(opcje.tytul, czynnosci), tresc);

  return {
    kod: opcje.kod,
    element,
    tresc,
    odswiez: () => opcje.panel.odswiez(),
    // Zamknięcie obudowy zamyka panel, nie tylko zdejmuje element ze sceny, bo subskrypcja by przeżyła.
    zamknij: () => opcje.panel.zamknij(),
  };
}
