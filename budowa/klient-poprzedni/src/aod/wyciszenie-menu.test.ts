// @vitest-environment jsdom

import { afterEach, describe, expect, it, vi } from 'vitest';
import { utworzStanObecnosci, type StanObecnosci } from './tryb-obecnosci';
import { KlasaZdarzen, ZakresKontekstu } from './wyciszenie-aod';
import type { KontekstWyciszenia } from './wyciszenie-kontekst';
import { utworzMenuWyciszenia, type MenuWyciszenia } from './wyciszenie-menu';

/**
 * Sprawdziany menu kebab wyciszania — czynność Operatora, nie kształt pliku.
 *
 * Menu jest tą samą powierzchnią przy awatarze i w nagłówku kolumny, więc
 * pilnowane jest to, czego zlecenie żąda od obu:
 *   1. komplet rozdz. 3.5 stoi w jednym menu — trzy czasy, moduł, karta sesji,
 *      sześć klas zdarzeń, tryb cichy;
 *   2. ŻADNA pozycja nie jest wyszarzana i ŻADNA nie pyta o potwierdzenie;
 *   3. wyciszenie i jego zniesienie idą jednym kliknięciem;
 *   4. podgląd wyciszeń czynnych mówi, CO jest wyciszone i DO KIEDY;
 *   5. pozycja, której nakładka nie ma czym wykonać, mówi czego brakuje —
 *      zamiast zniknąć albo zmilczeć.
 */

const TERAZ = 1_700_000_000_000;

/** Kontekst nakładki podstawiony: moduł i karta sesji znane z nazwy. */
function kontekstZnany(): KontekstWyciszenia {
  return {
    odswiez: () => Promise.resolve(),
    biezacyModul: () => ({ id: 'mod_studio', nazwa: 'Studio' }),
    biezacaKarta: () => ({ id: 'ses_1', nazwa: 'Praca nad ofertą' }),
    modulOkna: () => 'mod_studio',
    nazwaModulu: (id) => id,
    powodBrakuModulu: () => 'nie dotyczy',
    powodBrakuKarty: () => 'nie dotyczy',
  };
}

/** Kontekst nakładki bez bieżącego bytu — rdzeń nie wskazał okna. */
function kontekstNieznany(): KontekstWyciszenia {
  return {
    odswiez: () => Promise.resolve(),
    biezacyModul: () => null,
    biezacaKarta: () => null,
    modulOkna: () => undefined,
    nazwaModulu: (id) => id,
    powodBrakuModulu: () => 'Nakładka nie wie, w którym module Operator pracuje.',
    powodBrakuKarty: () => 'Nakładka nie ma dziś karty sesji.',
  };
}

let zbudowane: MenuWyciszenia[] = [];

function zbuduj(kontekst: KontekstWyciszenia = kontekstZnany()): {
  menu: MenuWyciszenia;
  stan: StanObecnosci;
} {
  const stan = utworzStanObecnosci({ magazyn: null });
  const menu = utworzMenuWyciszenia({ stan, teraz: () => TERAZ, kontekst });
  document.body.append(menu.element);
  zbudowane.push(menu);
  menu.otworz();
  return { menu, stan };
}

/** Pozycja menu o treści zawierającej podany fragment. */
function pozycja(menu: MenuWyciszenia, fragment: string): HTMLButtonElement {
  const znaleziona = [...menu.element.querySelectorAll('button')].find((przycisk) =>
    (przycisk.textContent ?? '').includes(fragment),
  );
  expect(znaleziona, `pozycja „${fragment}" nie stoi w menu`).toBeDefined();
  return znaleziona as HTMLButtonElement;
}

afterEach(() => {
  for (const menu of zbudowane) menu.rozlacz();
  zbudowane = [];
  document.body.replaceChildren();
});

describe('menu kebab wyciszania', () => {
  it('niesie komplet pięciu rodzajów wyciszenia z rozdz. 3.5', () => {
    const { menu } = zbuduj();
    const tresc = menu.element.textContent ?? '';

    expect(tresc).toContain('Wycisz na 15 minut');
    expect(tresc).toContain('Wycisz na godzinę');
    expect(tresc).toContain('Wycisz do końca dnia');
    expect(tresc).toContain('Wycisz bieżący moduł „Studio"');
    expect(tresc).toContain('Wycisz bieżącą kartę sesji „Praca nad ofertą"');
    expect(tresc).toContain('Przejdź w tryb cichy');
    // Sześć klas rozdz. 3.2, każda swoją pozycją.
    const klasy = [...menu.element.querySelectorAll('button')].filter((przycisk) =>
      (przycisk.textContent ?? '').startsWith('Wycisz klasę'),
    );
    expect(klasy).toHaveLength(6);
  });

  it('żadna pozycja nie jest wyszarzana i żadna nie pyta o potwierdzenie', () => {
    const potwierdzenie = vi.spyOn(globalThis, 'confirm').mockReturnValue(true);
    const { menu } = zbuduj();

    const przyciski = [...menu.element.querySelectorAll('button')];
    expect(przyciski.length).toBeGreaterThan(10);
    for (const przycisk of przyciski) {
      expect(przycisk.disabled).toBe(false);
      expect(przycisk.getAttribute('aria-disabled')).toBeNull();
      przycisk.click();
    }

    expect(potwierdzenie).not.toHaveBeenCalled();
  });

  it('wycisza czasowo jednym kliknięciem i mówi, do kiedy', () => {
    const { menu, stan } = zbuduj();

    pozycja(menu, 'Wycisz na godzinę').click();

    expect(stan.wyciszenie(TERAZ)?.doChwili).toBe(TERAZ + 3_600_000);
    expect(menu.element.textContent).toContain(
      new Date(TERAZ + 3_600_000).toLocaleTimeString(),
    );
  });

  it('wycisza bieżący moduł i znosi to jednym kliknięciem', () => {
    const { menu, stan } = zbuduj();

    pozycja(menu, 'Wycisz bieżący moduł').click();
    expect(stan.wyciszenia.czyKontekstWyciszony(ZakresKontekstu.Modul, 'mod_studio')).toBe(true);
    expect(menu.element.textContent).toContain('Studio');

    pozycja(menu, 'Znieś wyciszenie modułu').click();
    expect(stan.wyciszenia.czyKontekstWyciszony(ZakresKontekstu.Modul, 'mod_studio')).toBe(false);
  });

  it('wycisza bieżącą kartę sesji', () => {
    const { menu, stan } = zbuduj();
    pozycja(menu, 'Wycisz bieżącą kartę sesji').click();
    expect(stan.wyciszenia.czyKontekstWyciszony(ZakresKontekstu.KartaSesji, 'ses_1')).toBe(true);
  });

  it('wycisza wskazaną klasę zdarzeń i przywraca ją', () => {
    const { menu, stan } = zbuduj();

    pozycja(menu, 'Wycisz klasę „stan kolejki zadań"').click();
    expect(stan.wyciszenia.czyKlasaWyciszona(KlasaZdarzen.StanKolejkiZadan)).toBe(true);

    pozycja(menu, 'Przywróć klasę „stan kolejki zadań"').click();
    expect(stan.wyciszenia.czyKlasaWyciszona(KlasaZdarzen.StanKolejkiZadan)).toBe(false);
  });

  it('przełącza tryb cichy z menu', () => {
    const { menu, stan } = zbuduj();
    pozycja(menu, 'Przejdź w tryb cichy').click();
    expect(stan.czySyntezaMowy()).toBe(false);
    pozycja(menu, 'Wyjdź z trybu cichego').click();
    expect(stan.czySyntezaMowy()).toBe(true);
  });

  it('podgląd wyciszeń czynnych znosi każde osobno i wszystkie naraz', () => {
    const { menu, stan } = zbuduj();

    pozycja(menu, 'Wycisz na 15 minut').click();
    pozycja(menu, 'Wycisz klasę „harmonogram"').click();
    expect(stan.wyciszenia.czynne(TERAZ)).toHaveLength(2);

    const znoszace = [...menu.element.querySelectorAll<HTMLButtonElement>(
      '.ao-wyciszenie__znies',
    )];
    expect(znoszace).toHaveLength(2);
    znoszace[0]?.click();
    expect(stan.wyciszenia.czynne(TERAZ)).toHaveLength(1);

    pozycja(menu, 'Znieś wszystkie wyciszenia').click();
    expect(stan.wyciszenia.czynne(TERAZ)).toHaveLength(0);
    expect(menu.element.textContent).toContain('Wyciszeń czynnych: brak');
  });

  it('pozycja bez bieżącego bytu zostaje klikalna i mówi, czego brakuje', () => {
    const { menu, stan } = zbuduj(kontekstNieznany());

    const przycisk = pozycja(menu, 'Wycisz bieżący moduł');
    expect(przycisk.disabled).toBe(false);
    przycisk.click();

    expect(stan.wyciszenia.czynne(TERAZ)).toHaveLength(0);
    expect(menu.element.textContent).toContain('Nakładka nie wie, w którym module');
  });

  it('niesie zdanie o wyjątku wagi krytycznej i o granicy kontraktu', () => {
    const { menu } = zbuduj();
    const tresc = menu.element.textContent ?? '';
    expect(tresc).toContain('Wyjątek wagi krytycznej');
    expect(tresc).toContain('stanem TEGO okna');
  });

  it('otwiera się i zamyka z klawiatury, znak niesie stan rozwinięcia', () => {
    const { menu } = zbuduj();
    expect(menu.znak.getAttribute('aria-expanded')).toBe('true');

    menu.element.dispatchEvent(
      new KeyboardEvent('keydown', { key: 'Escape', bubbles: true }),
    );
    expect(menu.czyOtwarte()).toBe(false);
    expect(menu.znak.getAttribute('aria-expanded')).toBe('false');

    menu.znak.click();
    expect(menu.czyOtwarte()).toBe(true);
  });

  it('dwie kopie menu nad jednym stanem mówią to samo', () => {
    const stan = utworzStanObecnosci({ magazyn: null });
    const pierwsze = utworzMenuWyciszenia({ stan, teraz: () => TERAZ, kontekst: kontekstZnany() });
    const drugie = utworzMenuWyciszenia({ stan, teraz: () => TERAZ, kontekst: kontekstZnany() });
    zbudowane.push(pierwsze, drugie);
    document.body.append(pierwsze.element, drugie.element);
    pierwsze.otworz();
    drugie.otworz();

    pozycja(pierwsze, 'Wycisz klasę „harmonogram"').click();

    expect(drugie.element.textContent).toContain('Przywróć klasę „harmonogram"');
  });
});
