/**
 * Pas zakładek kontenera Observability Tools steruje przełączaniem obszarów wraz
 * z wędrówką strzałkami: wygląd pochodzi z biblioteki komponentów, tutaj leży
 * wyłącznie zachowanie.
 */

/** Jedna zakładka kontenera: kod obszaru sterujący wyborem, nazwa widoczna na przycisku i treść osadzana pod pasem. */
export interface PozycjaZakladkiNarzedzi {
  /** Kod obszaru — nośnik wyboru i wartość `data-zakladka`, nie tekst na ekran. */
  kod: string;
  /** Nazwa własna narzędzia widoczna na przycisku zakładki. */
  nazwa: string;
  /** Obszar osadzany pod pasem zakładek. */
  element: HTMLElement;
}

export interface ZakladkiNarzedzi {
  /** Pas zakładek wraz z obszarami pod nim. */
  element: HTMLElement;
  /** Kod obszaru czynnego. */
  czynna(): string;
  /** Przełącza obszar; kod spoza wykazu nie zmienia niczego. */
  pokaz(kod: string): void;
  /** Subskrypcja przełączenia obszaru — po niej okno dociąga materiał. */
  naZmiane(sluchacz: (kod: string) => void): void;
}

export function utworzZakladkiNarzedzi(
  pozycje: readonly PozycjaZakladkiNarzedzi[],
): ZakladkiNarzedzi {
  const sluchacze = new Set<(kod: string) => void>();

  const pasek = document.createElement('div');
  pasek.className = 'dn-zakladki dg-narzedzia__pasek';
  pasek.setAttribute('role', 'tablist');
  pasek.setAttribute('aria-label', 'Narzędzia obserwowalności');

  const obszary = document.createElement('div');
  obszary.className = 'dg-narzedzia__obszary';

  const element = document.createElement('div');
  element.className = 'dg-narzedzia';
  element.append(pasek, obszary);

  const przyciski = new Map<string, HTMLButtonElement>();
  const kody = pozycje.map((pozycja) => pozycja.kod);
  let czynna = kody[0] ?? '';

  function oznacz(): void {
    for (const [kod, przycisk] of przyciski) {
      const wybrana = kod === czynna;
      przycisk.setAttribute('aria-selected', String(wybrana));
      // Jeden przystanek tabulatora na cały pas — reszta idzie strzałkami.
      przycisk.tabIndex = wybrana ? 0 : -1;
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

  /** Przenosi wybór o wskazany krok i zabiera za nim ognisko klawiatury. */
  function przesun(krok: number): void {
    const teraz = kody.indexOf(czynna);
    if (teraz < 0 || kody.length === 0) return;
    const nastepny = kody[(teraz + krok + kody.length) % kody.length] ?? czynna;
    pokaz(nastepny);
    przyciski.get(nastepny)?.focus();
  }

  for (const pozycja of pozycje) {
    const przycisk = document.createElement('button');
    przycisk.type = 'button';
    przycisk.className = 'dn-zakladka';
    przycisk.setAttribute('role', 'tab');
    przycisk.dataset['zakladka'] = pozycja.kod;
    przycisk.textContent = pozycja.nazwa;
    przycisk.addEventListener('click', () => pokaz(pozycja.kod));
    przycisk.addEventListener('keydown', (zdarzenie) => {
      if (zdarzenie.key === 'ArrowRight') przesun(1);
      else if (zdarzenie.key === 'ArrowLeft') przesun(-1);
      else return;
      zdarzenie.preventDefault();
    });
    przyciski.set(pozycja.kod, przycisk);
    pasek.append(przycisk);

    pozycja.element.classList.add('dg-narzedzia__obszar');
    pozycja.element.setAttribute('role', 'tabpanel');
    pozycja.element.setAttribute('aria-label', pozycja.nazwa);
    obszary.append(pozycja.element);
  }

  oznacz();

  return {
    element,
    czynna: () => czynna,
    pokaz,
    naZmiane: (sluchacz) => void sluchacze.add(sluchacz),
  };
}

/**
 * Ciało zakładki — jedna obudowa dla wszystkich pięciu.
 *
 * Obudowa stoi tu, a nie w każdej zakładce z osobna, bo pięć zakładek
 * składających własne pudełko rozjedzie się przy pierwszej zmianie odstępu,
 * a różnią się treścią, nie kształtem.
 */
export function cialoNarzedzia(...czesci: readonly HTMLElement[]): HTMLElement {
  const element = document.createElement('div');
  element.className = 'dg-narzedzie';
  element.append(...czesci);
  return element;
}

/** Pasek czynności zakładki: przyciski akcji i pola formularza nad miejscem treści, ułożone w jednym rzędzie poziomym. */
export function pasekNarzedzia(...kontrolki: readonly HTMLElement[]): HTMLElement {
  const element = document.createElement('div');
  element.className = 'dg-narzedzie__pasek';
  element.append(...kontrolki);
  return element;
}

/**
 * Akapit objaśnienia zakładki: co widać i czego kontrakt nie niesie.
 *
 * Zakładki warstwy eksperckiej mówią o własnej granicy zdaniem, nie milczeniem
 * — dlatego akapit jest częścią wyposażenia kontenera, a nie ozdobą jednej
 * zakładki.
 */
export function objasnienieNarzedzia(zdanie: string): HTMLElement {
  const element = document.createElement('p');
  element.className = 'dn-pole-opis dg-narzedzie__opis';
  element.textContent = zdanie;
  return element;
}
