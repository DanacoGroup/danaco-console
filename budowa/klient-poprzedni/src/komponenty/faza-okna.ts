/**
 * Faza okna operacyjnego — jeden zestaw wartości dla wszystkich modułów.
 *
 * Wartość trafia do atrybutu `data-faza` w drzewie DOM i do sprawdzianów, więc
 * musi brzmieć tak samo w każdym module — inaczej wspólny selektor arkusza ani
 * wspólny sprawdzian nie mają się o co oprzeć. Rodzaj nijaki bierze się z frazy
 * „okno jest puste / gotowe".
 *
 * Zestaw nie zna wartości `spoczynek`. „Jeszcze nie pytałem rdzenia" to stan
 * źródła danych, nie stan okna, i ma własny typ `FazaOdczytu`
 * (`dostepy/stan-dostepow.ts` oraz stany pozostałych modułów); moduły
 * odwzorowują go na okienne `puste` z osobną treścią komunikatu — tak robią
 * `moduly/design/okno-assets-panel.ts` i `moduly/library/okno-library-explorer.ts`.
 *
 * Widoczność wskaźnika odczytu i chowanie treści na czas ładowania zostają
 * w modułach; różnią się one tu świadomie (Assistant i Apps zostawiają treść
 * widoczną w ładowaniu, Design i Research ją chowają).
 */

/**
 * Faza, w której stoi okno operacyjne.
 *
 * `puste`     — wywołanie się udało, wyniku nie ma; także „jeszcze nie pytałem".
 * `ladowanie` — wywołanie rdzenia w toku.
 * `blad`      — odmowa rdzenia albo awaria; komunikat zostaje w układzie.
 * `gotowe`    — komunikat zdjęty, widać samą treść.
 */
export type FazaOkna = 'puste' | 'ladowanie' | 'blad' | 'gotowe';

/**
 * Znakuje powłokę okna fazą i ustawia oprawę komunikatu: atrybut `data-faza`
 * na elemencie zewnętrznym (po nim sięgają arkusze `[data-faza='blad']`
 * i sprawdziany okien), rolę dostępności na pasie komunikatu oraz jego
 * widoczność.
 *
 * Rola `alert` należy się wyłącznie odmowie, bo przerywa czytnikowi ekranu
 * bieżącą wypowiedź; pozostałe fazy idą jako `status`. Pas znika tylko
 * w `gotowe` — stan przesłania treść, a nie kasuje jej, więc po powrocie
 * do `gotowe` Operator widzi to, co widział przed nieudanym odświeżeniem.
 *
 * @param element powłoka okna — nośnik atrybutu `data-faza`
 * @param pas     pas komunikatu stanu (`dn-pusty-stan`), przesłaniający treść
 */
export function oznaczFaze(element: HTMLElement, pas: HTMLElement, faza: FazaOkna): void {
  element.dataset['faza'] = faza;
  pas.setAttribute('role', faza === 'blad' ? 'alert' : 'status');
  pas.hidden = faza === 'gotowe';
}
