import type { PanelJezyka } from './panel-jezyka';

/**
 * Widok porównawczy dwóch paneli.
 *
 * Czynność czysto miejscowa: zestawienie niczego nie liczy i o nic nie pyta
 * rdzenia, pokazuje obok siebie tekst źródłowy i treści dwóch paneli, które
 * Operator sam wskazał.
 *
 * Panele są dokładnie dwa; przy innej liczbie wskazań zestawienie mówi wprost,
 * czego oczekuje, zamiast pokazać cokolwiek.
 */
export interface PorownaniePaneli {
  element: HTMLElement;
  /** Przerysowuje zestawienie ze wskazań paneli i bieżącego tekstu źródłowego. */
  odswiez(panele: readonly PanelJezyka[], tekstZrodlowy: string): void;
}

export function utworzPorownaniePaneli(): PorownaniePaneli {
  const opis = document.createElement('p');
  opis.className = 'mt-porownanie__opis';

  const kolumny = document.createElement('div');
  kolumny.className = 'mt-porownanie__kolumny';

  const element = document.createElement('section');
  element.className = 'mt-porownanie';
  element.setAttribute('aria-label', 'Widok porównawczy dwóch paneli');
  element.append(opis, kolumny);

  return {
    element,

    odswiez(panele, tekstZrodlowy) {
      const wskazane = panele.filter((panel) => panel.doPorownania());
      element.dataset['wskazane'] = String(wskazane.length);
      if (wskazane.length !== 2) {
        opis.textContent =
          `Wskaż dokładnie dwa panele znacznikiem w ich nagłówkach, aby je zestawić ` +
          `(wskazane: ${wskazane.length}).`;
        kolumny.replaceChildren();
        return;
      }
      opis.textContent = `Zestawienie: ${wskazane.map((panel) => panel.jezyk()).join(' ↔ ')}`;
      kolumny.replaceChildren(
        kolumna('źródło', tekstZrodlowy),
        ...wskazane.map((panel) => kolumna(panel.jezyk(), panel.tresc())),
      );
    },
  };
}

function kolumna(naglowek: string, tresc: string): HTMLElement {
  const tytul = document.createElement('h5');
  tytul.className = 'mt-porownanie__tytul';
  tytul.textContent = naglowek;

  const cialo = document.createElement('p');
  cialo.className = 'mt-porownanie__tresc';
  cialo.textContent = tresc === '' ? '(pusto)' : tresc;

  const element = document.createElement('div');
  element.className = 'mt-porownanie__kolumna';
  element.append(tytul, cialo);
  return element;
}
