/**
 * Korzeń sceny roboczej — rama, w której stoją okna komunikacji i ich
 * sterowanie. Przygotowuje puste miejsca montażu sceny sesji i nie zna ani
 * rdzenia, ani układu okien, ani kompletu sterowania; wystawia wyłącznie
 * miejsca dla kolejnych warstw.
 */
export interface KorzenAplikacji {
  /** Element ramy, montowany w obszarze roboczym powłoki. */
  element: HTMLElement;
  /** Scena okien komunikacji — miejsce układu okien równoległych. */
  scena: HTMLElement;
  /** Kolumna sterowania okna, obok sceny. */
  panel: HTMLElement;
  /** Miejsce akcji paska górnego, pożyczone od powłoki środowiska. */
  akcje: HTMLElement;
}

/**
 * Buduje ramę sceny sesji. Klasy `dn-powloka__*` opisują arkusze
 * `aplikacja/powloka.css` oraz `okna-rownolegle/uklad.css`, więc scena
 * zachowuje rozkład ustalony przez oba katalogi bez powielania ani jednej
 * reguły.
 */
export function utworzKorzen(akcjePaska: HTMLElement): KorzenAplikacji {
  const element = document.createElement('div');
  // Korzeń sceny nosi własną nazwę, inaczej wpadłby pod regułę klasy `dn-sesja` z pasa powłoki.
  element.className = 'dn-scena-sesji dn-powloka__obszar';

  const scena = document.createElement('div');
  scena.className = 'dn-powloka__scena';

  const panel = document.createElement('aside');
  panel.className = 'dn-powloka__panel';
  panel.setAttribute('aria-label', 'Sterowanie okien komunikacji');

  element.append(scena, panel);

  return { element, scena, panel, akcje: akcjePaska };
}
