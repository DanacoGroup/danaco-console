import {
  NAZWY_TRYBOW,
  OPISY_TRYBOW,
  TrybObecnosci,
  type StanObecnosci,
} from './tryb-obecnosci';
import type { KontekstWyciszenia } from './wyciszenie-kontekst';
import { utworzMenuWyciszenia, type MenuWyciszenia } from './wyciszenie-menu';

/**
 * Sterowanie obecnością w nagłówku powierzchni interakcji: przełącznik trybu
 * `[ tryb ▼ ]` i menu kebab `⋮` z wyciszaniem
 * (`docs/funkcje-globalne/always-on-display.md`, rozdz. 2.6, 3.5, 8.3, 9.2).
 *
 * Przełącznik trybu jest kontrolką trzystanową — pełny, cichy, ukryty — i
 * zmiana obowiązuje natychmiast, bez przycisku „zapisz" i bez pytania „czy na
 * pewno".
 *
 * Menu kebab nie jest tu budowane po swojemu: to ten sam komponent
 * (`wyciszenie-menu.ts`), który stoi przy awatarze. Dwie kopie jednego menu
 * czytają jeden stan, więc wyciszenie założone przy awatarze widać w nagłówku
 * bez żadnego zszywania.
 *
 * Stan nie jest tu przechowywany: obie kontrolki czytają i przestawiają
 * `StanObecnosci` (`tryb-obecnosci.ts`), a rysują się z powrotem jego
 * powiadomieniem. Dzięki temu skrót klawiszowy i menu nigdy się nie rozjeżdżają.
 */

export interface SterowanieObecnoscia {
  /** Przełącznik trybu obecności do wstawienia w nagłówek. */
  przelacznik: HTMLElement;
  /** Menu kebab do wstawienia w nagłówek. */
  kebab: HTMLElement;
  /** Przerysowuje obie kontrolki ze stanu. */
  odswiez(teraz: number): void;
  rozlacz(): void;
}

/** Nastawy sterowania wykraczające poza stan obecności. */
export interface OpisSterowaniaObecnoscia {
  /** Kontekst nakładki dla wyciszenia kontekstowego; pominięty nazywa brak wprost. */
  kontekst?: KontekstWyciszenia;
  /** Krótkie zameldowanie czynności na pasku okna. */
  zamelduj?: (zdanie: string) => void;
  /** Wołane przy otwarciu menu — miejsce na odświeżenie kontekstu z rdzenia. */
  naOtwarcieMenu?: () => void;
}

export function utworzSterowanieObecnoscia(
  stan: StanObecnosci,
  teraz: () => number,
  opis: OpisSterowaniaObecnoscia = {},
): SterowanieObecnoscia {
  // --- przełącznik trybu obecności ---------------------------------------

  const przelacznik = document.createElement('div');
  przelacznik.className = 'ao-tryby';
  przelacznik.setAttribute('role', 'radiogroup');
  przelacznik.setAttribute('aria-label', 'Tryb obecności Always On Display');

  const pigulki = new Map<TrybObecnosci, HTMLButtonElement>();
  for (const tryb of Object.values(TrybObecnosci)) {
    const pigulka = document.createElement('button');
    pigulka.type = 'button';
    pigulka.className = 'ao-tryb';
    pigulka.textContent = NAZWY_TRYBOW[tryb];
    pigulka.title = OPISY_TRYBOW[tryb];
    pigulka.setAttribute('role', 'radio');
    pigulka.addEventListener('click', () => stan.ustawTryb(tryb));
    pigulki.set(tryb, pigulka);
    przelacznik.append(pigulka);
  }

  // --- menu kebab ---------------------------------------------------------

  const menu: MenuWyciszenia = utworzMenuWyciszenia({
    stan,
    teraz,
    ...(opis.kontekst === undefined ? {} : { kontekst: opis.kontekst }),
    ...(opis.zamelduj === undefined ? {} : { zamelduj: opis.zamelduj }),
    ...(opis.naOtwarcieMenu === undefined ? {} : { naOtwarcie: opis.naOtwarcieMenu }),
    etykietaZnaku: 'Menu wyciszania Always On Display — nagłówek powierzchni interakcji',
  });

  function odswiez(chwila: number): void {
    const biezacy = stan.tryb();
    for (const [tryb, pigulka] of pigulki) {
      const wybrany = tryb === biezacy;
      pigulka.setAttribute('aria-checked', String(wybrany));
      // Stan nie idzie samą barwą: wybrana pigułka niesie też znak wyboru.
      pigulka.dataset['wybrany'] = String(wybrany);
      pigulka.textContent = wybrany ? `● ${NAZWY_TRYBOW[tryb]}` : NAZWY_TRYBOW[tryb];
    }

    menu.odswiez(chwila);
  }

  odswiez(teraz());

  return {
    przelacznik,
    kebab: menu.element,
    odswiez,
    rozlacz() {
      menu.rozlacz();
    },
  };
}
