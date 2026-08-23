/**
 * Wskazanie źródła otwartego w Reading View.
 *
 * Jedna odpowiedzialność: jeden identyfikator wraz z ogłoszeniem jego zmiany.
 * Wydzielone, bo dotyczy dwóch okien naraz — Sources Manager naciska „Czytaj",
 * Reading View wczytuje wskazany materiał — a `ZaznaczeniePozycji` niesie zbiór
 * i tu byłby narzędziem o jeden wymiar za dużym: czytać można jedno źródło.
 *
 * Wskazanie nie warunkuje klikalności ani jednej kontrolki. Reading View bez
 * wskazanego źródła pokazuje zdanie „nie wskazano czego czytać", a nie wygaszony
 * przycisk.
 *
 * Ogłoszenie idzie wyłącznie po faktycznej zmianie — odbiorcą jest przerysowanie
 * okien, a przerysowanie woła `ogranicz`; ogłoszenie bezwarunkowe zamknęłoby
 * pętlę bez końca. Ten sam wzorzec niesie `zaznaczenie-pozycji.ts`.
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
 *              zmiana dochodzi do obu okien patrzących na to samo źródło
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
