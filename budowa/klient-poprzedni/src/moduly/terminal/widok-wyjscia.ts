import type { WierszWyjscia } from './bufor-wyjscia';

/**
 * Rysowanie wierszy wyjścia z grepem i zawężeniem karty pokazuje wyłącznie ogon bufora, ograniczony do stałej wysokości okna.
 */
export interface NastawyWidoku {
  wzorzec: string;
  regularne: boolean;
  wielkoscLiter: boolean;
  znaczniki: boolean;
  /** Zawężenie do jednej karty; brak znaczy wynik zbiorczy. */
  karta?: string;
  /** Ostatni wiersz brany pod uwagę — nośnik odtwarzania sesji; ujemna znaczy „wszystkie”. */
  doWiersza: number;
  /**
   * Ile wierszy naraz wchodzi do dokumentu; pominięte bierze wartość domyślną okna rysowania.
   */
  oknoRysowania?: number;
}

export interface WynikRysowania {
  element: HTMLElement;
  /** Treść błędu wzorca; pusta, gdy wzorzec jest poprawny. */
  blad: string;
}

/**
 * Ile wierszy naraz wchodzi do dokumentu w widoku wyjścia, gdy wywołanie samo nie poda innej wartości.
 */
export const OKNO_RYSOWANIA = 800;

export function rysujWiersze(
  wiersze: readonly WierszWyjscia[],
  nastawy: NastawyWidoku,
): WynikRysowania {
  const element = document.createElement('pre');
  element.className = 'dt-konsola';
  element.setAttribute('aria-label', 'Zbiorcze wyjście poleceń');

  let dopasowanie: (tresc: string) => boolean;
  if (nastawy.wzorzec === '') {
    dopasowanie = () => true;
  } else if (nastawy.regularne) {
    try {
      const wyrazenie = new RegExp(nastawy.wzorzec, nastawy.wielkoscLiter ? '' : 'i');
      dopasowanie = (tresc) => wyrazenie.test(tresc);
    } catch (blad) {
      return { element, blad: blad instanceof Error ? blad.message : String(blad) };
    }
  } else {
    const szukane = nastawy.wielkoscLiter ? nastawy.wzorzec : nastawy.wzorzec.toLowerCase();
    dopasowanie = (tresc) =>
      (nastawy.wielkoscLiter ? tresc : tresc.toLowerCase()).includes(szukane);
  }

  // Zero jest poprawnym położeniem odtwarzania; wartość ujemna niesie dopiero wszystkie wiersze.
  const granica = nastawy.doWiersza < 0 ? wiersze.length : nastawy.doWiersza;
  const wybrane = wiersze
    .slice(0, granica)
    .filter((wiersz) => nastawy.karta === undefined || wiersz.karta === nastawy.karta)
    .filter((wiersz) => dopasowanie(wiersz.tresc));

  for (const wiersz of wybrane.slice(-(nastawy.oknoRysowania ?? OKNO_RYSOWANIA))) {
    const linia = document.createElement('span');
    linia.className = 'dt-konsola__wiersz';
    linia.dataset['rodzaj'] = wiersz.rodzaj;
    linia.dataset['proces'] = wiersz.proces;
    linia.textContent = nastawy.znaczniki
      ? `${new Date(wiersz.chwila).toISOString().slice(11, 23)}  ${wiersz.tresc}`
      : wiersz.tresc;
    element.append(linia);
  }
  return { element, blad: '' };
}
