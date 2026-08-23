/**
 * Pas zakładek okna modułu Developer — przełącznik obszarów wewnątrz jednej
 * kolumny.
 *
 * Dwa okna modułu dzielą kolumnę na obszary. Monitor łączy Build Output
 * i Run & Debug, a Dev Tools zbiera cztery integracje deweloperskie; oba
 * przełączają obszary zakładkami w nagłówku kolumny, tak jak opisuje
 * opracowanie modułu. Wygląd w całości z biblioteki
 * (`komponenty/zakladki.css`, klasy `dn-zakladki` i `dn-zakladka`); tutaj leży
 * wyłącznie zachowanie.
 *
 * Zakładka niewidoczna nie jest zakładką porzuconą: obszar zostaje w drzewie
 * i traci wyłącznie widoczność, więc treść pola zadania, zebrany log i wpisany
 * filtr przeżywają zajrzenie do sąsiedniej zakładki.
 *
 * Wędrówka strzałkami należy do wzorca zakładek: pas ma jeden przystanek
 * tabulatora (zakładka czynna), a strzałki przenoszą wybór między zakładkami.
 * Bez tego pas byłby tyloma przystankami przed treścią, ile ma pozycji.
 *
 * Bliźniaczy mechanizm stoi w modułach Diagnostics i w oknie modeli. Wspólnego
 * komponentu zakładek biblioteka `komponenty/` dziś nie ma — rozstrzygnięcie,
 * czy ma powstać, należy do właściciela projektu.
 */

/** Jedna zakładka okna: kod obszaru, nazwa własna i jego treść. */
export interface PozycjaZakladkiOkna {
  /** Kod obszaru — nośnik wyboru i wartość `data-zakladka`, nie tekst na ekran. */
  kod: string;
  /** Nazwa własna obszaru, dokładnie jak w opracowaniu modułu. */
  nazwa: string;
  /** Obszar osadzany pod pasem zakładek. */
  element: HTMLElement;
}

export interface ZakladkiOkna {
  /** Pas zakładek — osadzany w nagłówku kolumny. */
  pasek: HTMLElement;
  /** Obszary zakładek — osadzane w ciele kolumny. */
  obszary: HTMLElement;
  /** Kod obszaru czynnego. */
  czynna(): string;
  /** Przełącza obszar; kod spoza wykazu nie zmienia niczego. */
  pokaz(kod: string): void;
  /** Subskrypcja przełączenia obszaru. */
  naZmiane(sluchacz: (kod: string) => void): void;
}

export function utworzZakladkiOkna(
  etykietaPasa: string,
  pozycje: readonly PozycjaZakladkiOkna[],
): ZakladkiOkna {
  const sluchacze = new Set<(kod: string) => void>();

  const pasek = document.createElement('div');
  pasek.className = 'dn-zakladki mdev-zakladki';
  pasek.setAttribute('role', 'tablist');
  pasek.setAttribute('aria-label', etykietaPasa);

  const obszary = document.createElement('div');
  obszary.className = 'mdev-zakladki__obszary';

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

  /** Przenosi wybór o wskazany krok i zabiera za nim ognisko. */
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

    pozycja.element.classList.add('mdev-zakladki__obszar');
    pozycja.element.setAttribute('role', 'tabpanel');
    pozycja.element.setAttribute('aria-label', pozycja.nazwa);
    obszary.append(pozycja.element);
  }

  oznacz();

  return {
    pasek,
    obszary,
    czynna: () => czynna,
    pokaz,
    naZmiane: (sluchacz) => void sluchacze.add(sluchacz),
  };
}

/**
 * Ciało zakładki — jedna obudowa dla wszystkich obszarów okna.
 *
 * Obudowa stoi tu, a nie w każdej zakładce z osobna, bo zakładki składające
 * własne pudełko rozjadą się przy pierwszej zmianie odstępu, a różnią się
 * treścią, nie kształtem.
 */
export function cialoZakladki(...czesci: readonly HTMLElement[]): HTMLElement {
  const element = document.createElement('div');
  element.className = 'mdev-zakladka-tresc';
  element.append(...czesci);
  return element;
}

/** Pasek czynności zakładki — przyciski i pola nad miejscem treści. */
export function pasekZakladki(...kontrolki: readonly HTMLElement[]): HTMLElement {
  const element = document.createElement('div');
  element.className = 'mdev-zakladka-pasek';
  element.append(...kontrolki);
  return element;
}

/** Akapit objaśnienia zakładki: co widać i czego kontrakt nie niesie. */
export function objasnienieZakladki(zdanie: string): HTMLElement {
  const element = document.createElement('p');
  element.className = 'dn-pole-opis mdev-zakladka-opis';
  element.textContent = zdanie;
  return element;
}
