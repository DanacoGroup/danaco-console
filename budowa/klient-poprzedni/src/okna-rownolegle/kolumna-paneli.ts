import './panele.css';

import type { PanelWStosie } from './panel-w-stosie';

/**
 * Pionowy stos paneli pomocniczych jednego gniazda.
 *
 * Panele nie są stałym elementem okna: otwiera się i zamyka każdy osobno,
 * więc kolumna mieści dowolny ich podzbiór — od zera do wielu.
 *
 * Panel należy do rozmowy, nie do ekranu. Kolumna obsługuje jedno gniazdo,
 * każde okno sceny ma własny stos, a kolumna nie zna ani identyfikatora sceny,
 * ani pozostałych kolumn.
 *
 * Zero paneli to stan poprawny i wtedy kolumna znika. Pusty stos nie dostaje
 * komunikatu ani ramki, tylko `hidden` — miejsce wraca do rozmowy, z którą
 * dzieli szerokość gniazda.
 *
 * Kolumna nie buduje paneli i nie zna wytwórni — dostaje gotowe, tak jak układ
 * dostaje gotowe okna. Nie zamyka też paneli, które z niej schodzą: cykl życia
 * panelu wraz z subskrypcją rdzenia należy do tego, kto panel powołał.
 * Zamykanie ich przy każdym `ustaw` ubiłoby panel przeniesiony na pełny ekran.
 *
 * Podział szerokości między rozmowę a kolumnę należy do sceny i jej arkusza —
 * kolumna wypełnia miejsce, które dostała, i nie zna liczb kontraktu szerokości.
 */
export interface KolumnaPaneli {
  /** Element osadzany w gnieździe, obok kolumny rozmowy. */
  element: HTMLElement;
  /** Ustawia zawartość stosu; pusty wykaz chowa kolumnę całkowicie. */
  ustaw(panele: readonly PanelWStosie[]): void;
  pusta(): boolean;
}

export function utworzKolumnaPaneli(): KolumnaPaneli {
  let ile = 0;

  const element = document.createElement('div');
  element.className = 'dn-okna__kolumna-paneli';
  element.setAttribute('aria-label', 'Panele pomocnicze otwarte obok rozmowy');
  element.hidden = true;

  return {
    element,
    ustaw(panele) {
      ile = panele.length;
      // `replaceChildren` zdejmuje poprzednie panele z drzewa, ale ich nie
      // zamyka — zamknięcie należy do właściciela panelu, nie do przerysowania
      // stosu.
      element.replaceChildren(...panele.map((panel) => panel.element));
      element.hidden = ile === 0;
    },
    pusta: () => ile === 0,
  };
}
