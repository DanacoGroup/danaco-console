/**
 * Wskaźnik odczytu katalogu z rdzenia — stan ładowania kompletu sterowania
 * (`.dn-spinner` z biblioteki).
 *
 * Wskaźnik stoi obok pola, nigdy zamiast pola: podmiana kontrolki na wskaźnik
 * byłaby blokadą, a Operator ma móc wybrać wartość także wtedy, gdy katalog
 * jeszcze jedzie z rdzenia.
 *
 * Wskaźnik mówi też technologiom wspomagającym, co się dzieje: `role="status"`
 * ogłasza zmianę bez zabierania ogniska, a `aria-label` niesie nazwę katalogu,
 * którego dotyczy odczyt.
 */
export interface WskaznikOdczytu {
  /** Element osadzany przy etykiecie sterowania. */
  element: HTMLElement;
  /** Zapala wskaźnik na czas odczytu i gasi go po rozstrzygnięciu. */
  ustaw(trwa: boolean): void;
}

export function utworzWskaznikOdczytu(opis: string): WskaznikOdczytu {
  const element = document.createElement('span');
  element.className = 'dn-spinner';
  element.setAttribute('role', 'status');
  element.setAttribute('aria-label', opis);
  element.hidden = true;

  return {
    element,
    ustaw(trwa) {
      element.hidden = !trwa;
    },
  };
}
