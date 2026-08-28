/**
 * Wskazanie źródła otwartego w Reading View: jeden identyfikator wraz z ogłoszeniem jego
 * zmiany, wspólny dla Sources Manager i Reading View.
 */
export interface WskazanieLektury {
  /** Wskazane źródło; pusty napis znaczy „nie wskazano". */
  wskazane(): string;
  /** Wskazuje źródło do lektury. */
  wskaz(identyfikator: string): void;
  /** Zdejmuje wskazanie źródła, którego nie ma już w katalogu. */
  ogranicz(istniejace: readonly string[]): void;
}

/**
 * @param oglos wywoływane po każdej faktycznej zmianie wskazania; przez nie
 *              zmiana dochodzi do obu okien korzystających z tego samego źródła
 */
export function utworzWskazanieLektury(oglos: () => void): WskazanieLektury {
  let wskazany = '';

  return {
    wskazane: () => wskazany,

    wskaz(identyfikator) {
      if (wskazany === identyfikator) return;
      wskazany = identyfikator;
      oglos();
    },

    ogranicz(istniejace) {
      if (wskazany === '' || istniejace.includes(wskazany)) return;
      wskazany = '';
      oglos();
    },
  };
}
