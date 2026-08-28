/**
 * Faza okna operacyjnego — jeden zestaw wartości dla wszystkich modułów.
 * Wartość trafia do atrybutu `data-faza` w drzewie DOM oraz do sprawdzianów,
 * więc musi brzmieć tak samo w każdym module.
 */

/**
 * Faza, w której stoi okno operacyjne: `puste` — wywołanie udane bez wyniku,
 * `ladowanie` — wywołanie rdzenia w toku, `blad` — odmowa albo awaria
 * z komunikatem w układzie, `gotowe` — komunikat zdjęty i widać samą treść.
 */
export type FazaOkna = 'puste' | 'ladowanie' | 'blad' | 'gotowe';

/**
 * Znakuje powłokę okna fazą: atrybut `data-faza` na elemencie zewnętrznym, rolę
 * dostępności pasa komunikatu oraz jego widoczność.
 * @param element powłoka okna, nośnik atrybutu `data-faza`
 * @param pas pas komunikatu stanu przesłaniający treść
 */
export function oznaczFaze(element: HTMLElement, pas: HTMLElement, faza: FazaOkna): void {
  element.dataset['faza'] = faza;
  pas.setAttribute('role', faza === 'blad' ? 'alert' : 'status');
  pas.hidden = faza === 'gotowe';
}
