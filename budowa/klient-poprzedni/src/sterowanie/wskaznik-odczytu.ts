/**
 * Wskaźnik odczytu katalogu z rdzenia — stan ładowania widoczny obok pola tego kompletu
 * sterowania okna.
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
