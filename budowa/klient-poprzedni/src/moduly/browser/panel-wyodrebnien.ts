import type { BrowserSnapshot } from '../../../../shared/contract';
import type { StanPrzegladania } from './stan-przegladania';

/**
 * Panel wyodrębnionej treści Browser Window — wynik pozycji
 * „Wyodrębnij dane ▾" paska dolnego i przycisku „Wyodrębnij" paska zaznaczenia.
 *
 * Jedna odpowiedzialność: pokazanie tego, co Operator wyjął z migawki, wraz
 * z wykazem wyodrębnionych fragmentów.
 *
 * Wyodrębnienie dzieje się w kliencie z treści, którą klient ma. Kontrakt nie
 * ma komendy wyodrębniania danych ze strony; migawka jest jedyną treścią,
 * którą moduł dostał od rdzenia, i to z niej wyjmowane są trzy rodzaje danych.
 * Panel nie zmyśla żadnego pola — do przyjścia migawki stoi pusty.
 */
export interface PanelWyodrebnien {
  element: HTMLElement;
  /** Wyodrębnia wskazany rodzaj treści z bieżącej migawki. */
  wyodrebnij(rodzaj: string): string;
  /** Nanosi stan zbioru wyodrębnionych fragmentów. */
  odswiez(): void;
}

export function utworzPanelWyodrebnien(stan: StanPrzegladania): PanelWyodrebnien {
  const tytul = document.createElement('h4');
  tytul.className = 'mb-panel__tytul';
  tytul.textContent = 'Wyodrębniona treść';

  const wynik = document.createElement('pre');
  wynik.className = 'mb-wyodrebnienia__wynik';

  const lista = document.createElement('ul');
  lista.className = 'mb-wyodrebnienia__lista';

  const element = document.createElement('section');
  element.className = 'mb-panel mb-wyodrebnienia';
  element.append(tytul, wynik, lista);

  return {
    element,

    wyodrebnij(rodzaj) {
      const migawka = stan.migawka();
      if (migawka === null) {
        wynik.textContent = 'Brak migawki — najpierw przejdź do strony albo odśwież migawkę.';
        return '';
      }
      const tresc = wyjmij(migawka, rodzaj);
      wynik.textContent = tresc === '' ? 'Migawka nie niesie tej treści.' : tresc;
      return tresc;
    },

    odswiez() {
      lista.replaceChildren(
        ...stan.zebrane.wyodrebnione().map((fragment) => {
          const wiersz = document.createElement('li');
          wiersz.className = 'mb-wyodrebnienia__wiersz';
          wiersz.textContent = fragment;
          return wiersz;
        }),
      );
    },
  };
}

/** Trzy rodzaje danych, które migawka naprawdę niesie. */
function wyjmij(migawka: BrowserSnapshot, rodzaj: string): string {
  if (rodzaj === 'zrodlo') return migawka.html ?? '';
  if (rodzaj === 'adres') return `${migawka.title ?? ''}\n${migawka.url}`.trim();
  return migawka.text ?? '';
}
