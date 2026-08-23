import './panele.css';

import type { PanelPomocniczy } from '../okna-pomocnicze/panel-pomocniczy';
import { utworzNaglowekPanelu, type PozycjaNaglowka } from './naglowek-panelu';

/**
 * Obudowa jednej pozycji stosu paneli obok rozmowy: gęsty nagłówek nad treścią
 * panelu i dwie czynności przy nim.
 *
 * Panel pomocniczy jest bytem otwieranym i zamykanym pojedynczo — nie ma
 * pulpitu, którego panele byłyby częścią — więc otwieranie i zamykanie ma
 * miejsce w samym panelu.
 *
 * Treść stoi osobnym polem (`tresc`), bo widok pełnoekranowy przenosi ją na
 * całą scenę i oddaje z powrotem. Gdyby oddawać cały `element`, na scenę
 * pojechałby razem z nagłówkiem, a w kolumnie zostałaby dziura bez obudowy.
 * Rozdział obudowy od wnętrza sprawia, że przenosiny są przenosinami, a nie
 * kopią — kopia znaczyłaby drugi egzemplarz panelu, czyli drugą subskrypcję
 * rdzenia.
 *
 * Pełny ekran bierze ikonę `link-zewnetrzny`, bo zestaw marki nie ma ikony
 * „rozwiń" (`ikony/zrodla/*.ts`), a ikona jest dobierana pod czynność, nie pod
 * wygląd narzędzia: czynnością jest wyniesienie treści poza jej ramy, a
 * `link-zewnetrzny` to jedyna ikona zestawu, która niesie „na zewnątrz".
 * Zamknięcie bierze `zamknij`.
 *
 * Obudową nie jest `komponenty/rama-okna.ts`: rama daje trzy pasy (nagłówek,
 * akcje, narzędzia), a ta powierzchnia mieści jeden wiersz. Rama zostaje
 * wewnątrz `okno-podglad-bash.ts`, które jest oknem operacyjnym pasa modułu.
 *
 * Plik nie buduje panelu i nie zna wytwórni — dostaje gotowy
 * `PanelPomocniczy`. Nie decyduje też, co się dzieje po zamknięciu ani po
 * wyjściu na pełny ekran; woła `naZamkniecie` i `naPelnyEkran`, a skutek
 * należy do gospodarza stosu.
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
  /**
   * Gotowe elementy rzędu ikon — miejsce na własne menu `⋮` panelu.
   *
   * Obudowa pozycji takiego menu nie zna i ich nie buduje: dokłada je
   * gospodarz gotowym elementem, a samo menu rozwijane powstaje osobno.
   *
   * Pominięte znaczy, że rząd niesie dwie czynności obudowy i nic więcej —
   * jedyny zbudowany panel (`podglad-bash`) własnego menu nie ma.
   */
  dodatkiNaglowka?: readonly HTMLElement[];
}

export function utworzPanelWStosie(opcje: OpcjePanelaWStosie): PanelWStosie {
  const tresc = document.createElement('div');
  tresc.className = 'dn-okna__panel-tresc';
  tresc.append(opcje.panel.element);

  // Rząd: menu własne panelu (gdy gospodarz je dał), potem dwie czynności
  // obudowy. Trzeciej czynności obudowa nie dokłada — ikona bez czynności jest
  // atrapą, a atrapa jest gorsza niż jej brak.
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
    // Zamknięcie obudowy zamyka panel, nie tylko zdejmuje element ze sceny:
    // subskrypcja rdzenia przeżyłaby usunięcie węzła z drzewa dokumentu.
    zamknij: () => opcje.panel.zamknij(),
  };
}
