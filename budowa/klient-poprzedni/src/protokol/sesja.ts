/** Bieżąca sesja klienta — identyfikator niesiony w kopertach wychodzących, wspólny dla plików, pamięci i projektu. */
export interface Sesja {
  /** Identyfikator sesji; pusty, dopóki rdzeń jej nie założył. */
  id(): string;
  /** Zapisuje identyfikator nadany przez rdzeń. */
  ustaw(identyfikator: string): void;
  /** Czy sesja została już założona. */
  zalozona(): boolean;
}

export function utworzSesje(): Sesja {
  let identyfikator = '';

  return {
    id: () => identyfikator,

    ustaw(nowy) {
      identyfikator = nowy;
    },

    zalozona: () => identyfikator.length > 0,
  };
}
