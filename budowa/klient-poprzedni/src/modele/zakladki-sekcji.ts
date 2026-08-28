/**
 * Zakładki sekcji modeli, czyli przełącznik czterech obszarów jednej sekcji:
 * kont, ustawień osi, tożsamości oraz podglądu promptu. Przełączenie ukrywa
 * obszar opuszczany, pozostawiając go w drzewie dokumentu wraz z treścią
 * wpisaną przez Operatora.
 */
export interface PozycjaZakladki {
  /** Kod obszaru — nośnik wyboru, nie tekst do wydruku. */
  kod: string;
  /** Nazwa obszaru pokazywana Operatorowi. */
  nazwa: string;
  /** Obszar osadzany pod paskiem zakładek. */
  element: HTMLElement;
}

export interface ZakladkiSekcji {
  /** Pasek zakładek wraz z obszarami pod nim. */
  element: HTMLElement;
  /** Kod obszaru czynnego. */
  czynna(): string;
  /** Przełącza obszar; kod spoza wykazu nie zmienia niczego. */
  pokaz(kod: string): void;
  /** Subskrypcja przełączenia obszaru. */
  naZmiane(sluchacz: (kod: string) => void): void;
}

export function utworzZakladkiSekcji(pozycje: readonly PozycjaZakladki[]): ZakladkiSekcji {
  const sluchacze: Array<(kod: string) => void> = [];

  const pasek = document.createElement('div');
  pasek.className = 'dn-zakladki dm-zakladki';
  pasek.setAttribute('role', 'tablist');

  const obszary = document.createElement('div');
  obszary.className = 'dm-zakladki__obszary';

  const element = document.createElement('div');
  element.className = 'dm-zakladki__rama';
  element.append(pasek, obszary);

  let czynna = pozycje[0]?.kod ?? '';

  const przyciski = new Map<string, HTMLButtonElement>();

  function oznacz(): void {
    for (const [kod, przycisk] of przyciski) {
      przycisk.setAttribute('aria-selected', String(kod === czynna));
    }
    for (const pozycja of pozycje) {
      pozycja.element.hidden = pozycja.kod !== czynna;
    }
  }

  function pokaz(kod: string): void {
    if (!przyciski.has(kod) || kod === czynna) return;
    czynna = kod;
    oznacz();
    for (const sluchacz of [...sluchacze]) sluchacz(kod);
  }

  for (const pozycja of pozycje) {
    const przycisk = document.createElement('button');
    przycisk.type = 'button';
    przycisk.className = 'dn-zakladka';
    przycisk.setAttribute('role', 'tab');
    przycisk.textContent = pozycja.nazwa;
    przycisk.addEventListener('click', () => pokaz(pozycja.kod));
    przyciski.set(pozycja.kod, przycisk);
    pasek.append(przycisk);

    pozycja.element.classList.add('dm-zakladki__obszar');
    obszary.append(pozycja.element);
  }

  oznacz();

  return {
    element,
    czynna: () => czynna,
    pokaz,
    naZmiane: (sluchacz) => void sluchacze.push(sluchacz),
  };
}
