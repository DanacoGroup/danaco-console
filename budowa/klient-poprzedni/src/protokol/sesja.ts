/**
 * Bieżąca sesja klienta — identyfikator niesiony w kopertach wychodzących.
 *
 * Sesja pozostaje wspólna dla plików, pamięci, projektu i agentów; okno
 * komunikacji jest wobec niej bytem podrzędnym. Identyfikator nadaje
 * rdzeń w odpowiedzi na `session.create` — do tej chwili sesja jest pusta,
 * a kontrakt taką kopertę dopuszcza (pole `sessionId` opcjonalne, puste dla
 * powitania połączenia).
 */
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
