/**
 * Korzeń sceny roboczej — rama, w której stoją okna komunikacji i ich
 * sterowanie.
 *
 * Przygotowuje puste miejsca montażu sceny sesji. Korzeń nie zna ani rdzenia,
 * ani układu okien, ani kompletu sterowania — wystawia wyłącznie miejsca,
 * w które montują się kolejne warstwy. Całym dokumentem zajmuje się
 * `gospodarz-dokumentu.ts`; kształt `KorzenAplikacji` jest wejściem
 * `zamontujUkladOkien` z katalogu `okna-rownolegle/`.
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
 * Buduje ramę sceny sesji.
 *
 * Klasy `dn-powloka__*` są celowe: opisują je arkusze `aplikacja/powloka.css`
 * oraz `okna-rownolegle/uklad.css`, więc scena zachowuje rozkład ustalony
 * przez oba katalogi bez powielania ani jednej reguły.
 *
 * Miejsce akcji nie powstaje tutaj — przychodzi z paska górnego powłoki
 * środowiska, żeby uchwyt szuflady sterowania stał w pasku wspólnym z resztą
 * akcji, a nie w osobnym pasku dla jednego przycisku.
 */
export function utworzKorzen(akcjePaska: HTMLElement): KorzenAplikacji {
  const element = document.createElement('div');
  // Klasa `dn-sesja` należy do karty sesji w pasie powłoki
  // (`powloka/karta-sesji.ts`), a jej arkusz zawęża szerokość do pasa nawigacji
  // bez zawężenia selektora rodzicem. Korzeń sceny nosi własną nazwę, inaczej
  // wpadłby pod tamtą regułę i zszedł do szerokości słupka.
  element.className = 'dn-scena-sesji dn-powloka__obszar';

  const scena = document.createElement('div');
  scena.className = 'dn-powloka__scena';

  const panel = document.createElement('aside');
  panel.className = 'dn-powloka__panel';
  panel.setAttribute('aria-label', 'Sterowanie okien komunikacji');

  element.append(scena, panel);

  return { element, scena, panel, akcje: akcjePaska };
}
