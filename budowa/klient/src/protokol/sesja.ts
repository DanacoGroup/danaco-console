/**
 * Bieżąca sesja klienta — identyfikator niesiony w kopertach wychodzących.
 * Identyfikator nadaje rdzeń, a klient poznaje go ze zdarzenia `session.focus.changed`
 * potwierdzającego ognisko tego klienta; do tej chwili sesja jest pusta.
 */
export interface Sesja {
  /** Identyfikator sesji; pusty, dopóki rdzeń nie potwierdził ogniska. */
  id(): string;
  /** Zapisuje identyfikator nadany przez rdzeń. */
  ustaw(identyfikator: string): void;
  /** Czy sesja została już założona. */
  zalozona(): boolean;
}

/*
Sesja jest jedna na uruchomienie klienta. Koperta wychodząca niesie jej
identyfikator, a wiązania okien rozstrzygają po nim, dokąd wraca praca —
druga instancja rozeszłaby te dwa odczyty i koperta mówiłaby co innego niż okno.
*/
let identyfikator = '';

const SESJA: Sesja = {
  id: () => identyfikator,

  ustaw(nowy) {
    identyfikator = nowy;
  },

  zalozona: () => identyfikator.length > 0,
};

/** Sesja bieżąca klienta — jedyne źródło identyfikatora sesji w całym kliencie. */
export function sesjaKlienta(): Sesja {
  return SESJA;
}
