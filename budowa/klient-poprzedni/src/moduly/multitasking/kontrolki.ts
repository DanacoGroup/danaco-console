/**
 * Wiersz opisu MultitaskingAI jest jedyną kontrolką własną modułu: ramę okna
 * niesie `komponenty/rama-okna.ts`, a kontrolki formularza
 * `modele/kontrolki-formularza.ts`. Wygląd pochodzi w całości z arkusza
 * rodziny `dm-`.
 */

/**
 * Przedrostek klas rodziny modułu, z którego wspólne kontrolki składają nazwy
 * pozycji wykazu. Stoi jedną stałą, bo arkusz rodziny należy do modułu, a trzy
 * wykazy okien ról muszą wskazywać ten sam arkusz.
 */
export const PRZEDROSTEK = 'dm';

/**
 * Wysokość pól redakcyjnych okien ról w wierszach. Stoi tu jedną liczbą, bo
 * trzy okna ról mają być jednakowo wysokie, a różnica wysokości czytałaby się
 * jako różnica wagi tych pól.
 */
export const WIERSZE_POLA = 3;

/**
 * Wiersz opisu złożony z etykiety i wartości, na przykład licznika obiegów,
 * progu albo powodu zatrzymania. Oddaje element klasy `dm-wiersz`, w którym
 * etykieta trafia zarówno do treści, jak i do atrybutu danych.
 */
export function wierszOpisu(etykieta: string, wartosc: string): HTMLElement {
  const element = document.createElement('div');
  element.className = 'dm-wiersz';
  element.dataset['klucz'] = etykieta;

  const nazwa = document.createElement('span');
  nazwa.className = 'dm-wiersz__etykieta';
  nazwa.textContent = etykieta;

  const tresc = document.createElement('span');
  tresc.className = 'dm-wiersz__wartosc';
  tresc.textContent = wartosc;

  element.append(nazwa, tresc);
  return element;
}
