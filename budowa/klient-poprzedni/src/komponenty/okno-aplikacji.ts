/**
 * Znakowanie okna aplikacji kodem katalogu rdzenia. Służy oknom, które mają
 * wiersz w katalogu okien operacyjnych, lecz nie noszą ramy operacyjnej
 * z `komponenty/rama-okna.ts`, ponieważ nagłówek, plakietka roli i pas akcji
 * byłyby w nich powtórzeniem.
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
 * Znakuje element kodem katalogu i nadaje mu etykietę dostępności o tej samej
 * nazwie. Zwraca ten sam element i ten sam typ, żeby dało się go użyć w miejscu
 * tworzenia bez utraty rodzaju.
 */
export function oznaczOknoAplikacji<T extends HTMLElement>(opis: OpisOknaAplikacji<T>): T {
  opis.element.dataset['okno'] = opis.kod;
  opis.element.setAttribute('aria-label', opis.nazwa);
  return opis.element;
}
