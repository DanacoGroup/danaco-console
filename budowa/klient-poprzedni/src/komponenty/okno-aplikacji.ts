/**
 * Znakowanie okna aplikacji kodem katalogu rdzenia.
 *
 * Służy oknom, które mają wiersz w katalogu okien operacyjnych, ale nie noszą
 * ramy z `komponenty/rama-okna.ts`: Okno Konfiguracji i Okno Ustawień są
 * modalami (`<dialog>`, wygląd niesie `dn-modal`), Strona główna jest pełnym
 * widokiem, a Okno Rozmowy dostaje ramę od gniazda układu równoległego
 * (`okna-rownolegle/gniazdo-okna.ts`). Nagłówek, plakietka roli i pas akcji
 * z ramy operacyjnej byłyby w każdym z tych przypadków elementem zbędnym albo
 * powtórzonym.
 *
 * Samo `element.dataset['okno'] = kod` rozsiane po tych plikach dawałoby kilka
 * miejsc do rozjechania się. Ta funkcja nadaje przy okazji etykietę dostępności
 * z tej samej nazwy, więc kod okna i to, co słyszy czytnik ekranu, nie mogą się
 * rozejść.
 */

export interface OpisOknaAplikacji<T extends HTMLElement> {
  /** Element okna — `<dialog>` modalu albo korzeń widoku pełnoekranowego. */
  element: T;
  /** Kod okna z katalogu rdzenia (`okno_operacyjne.kod`). */
  kod: string;
  /** Nazwa widziana przez Operatora i przez czytnik ekranu. */
  nazwa: string;
}

/**
 * Znakuje element kodem katalogu i nadaje mu etykietę dostępności.
 *
 * Zwraca ten sam element i ten sam typ, żeby dało się go użyć w miejscu
 * tworzenia bez utraty rodzaju: modal potrzebuje `HTMLDialogElement`, bo woła
 * `showModal()` i `close()`. Zwracanie `HTMLElement` odbierałoby te metody
 * i zmuszało wywołujących do rzutowania.
 */
export function oznaczOknoAplikacji<T extends HTMLElement>(opis: OpisOknaAplikacji<T>): T {
  opis.element.dataset['okno'] = opis.kod;
  opis.element.setAttribute('aria-label', opis.nazwa);
  return opis.element;
}
