/**
 * Wiersz opisu MultitaskingAI — jedyna kontrolka własna modułu.
 *
 * Rama okna stoi w `komponenty/rama-okna.ts`, a przycisk, pole, wybór, wykaz
 * i pozycja wykazu w `modele/kontrolki-formularza.ts` — moduł bierze je stamtąd
 * zamiast trzymać własne kopie. Wiersz klucz–wartość odpowiednika tam nie ma:
 * niosą go wyłącznie okna ról, gdzie zastępuje tabelę stanu (więź wykonawcy
 * z koordynatorem, licznik obiegów, powód zatrzymania biegu).
 *
 * Wygląd w całości z arkusza rodziny `dm-` — plik nie zna ani jednej barwy i ani
 * jednego odstępu.
 */

/**
 * Przedrostek klas rodziny modułu, z którego wspólne kontrolki składają nazwy
 * pozycji wykazu. Stoi jedną stałą, bo arkusz rodziny należy do modułu, a trzy
 * wykazy okien ról muszą wskazywać ten sam arkusz.
 */
export const PRZEDROSTEK = 'dm';

/**
 * Wysokość pól redakcyjnych okien ról w wierszach.
 *
 * Stoi tu jedną liczbą, bo trzy okna ról mają być jednakowo wysokie: kreator
 * promptu, polecenie wykonawcy i uzasadnienie oceny stoją obok siebie na scenie
 * i różnica wysokości czytałaby się jako różnica wagi tych pól.
 */
export const WIERSZE_POLA = 3;

/** Wiersz opisu: etykieta i wartość — licznik obiegów, próg, powód zatrzymania. */
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
