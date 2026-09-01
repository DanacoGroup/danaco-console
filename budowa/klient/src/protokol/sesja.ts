/**
 * Bieżąca sesja klienta — identyfikator niesiony w kopertach wychodzących.
 * Nadaje go rdzeń i potwierdza zdarzeniem `session.focus.changed`. Wartość
 * jest pochodną okna roboczego bieżącego (`wiazanie/sesja-biezaca.ts`): okno
 * bez sesji zostawia ją pustą, przełączenie okna ją przestawia.
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
Jedna wartość na klienta, choć zmienna: niesie sesję okna roboczego bieżącego.
Koperta wychodząca bierze stąd identyfikator, a wiązania okien rozstrzygają po
nim, dokąd wraca praca — druga instancja rozeszłaby te dwa odczyty.
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
