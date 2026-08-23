/**
 * Zaznaczenie wielokrotne pozycji wykazu.
 *
 * Jedna odpowiedzialność: zbiór zaznaczonych identyfikatorów wraz z jego
 * porządkowaniem po zmianie wykazu. Wydzielone, bo dotyczy dwóch okien naraz —
 * Sources Manager (wybór wielu źródeł do eksportu) i Findings Panel (wybór
 * ustaleń do raportu).
 *
 * Zaznaczenie nie warunkuje klikalności: przy pustym zbiorze akcja zostaje
 * czynna i odpowiada zdaniem opisowym, zamiast być wygaszona.
 *
 * Zbiór ogłasza swoją zmianę i dlatego konstruktor żąda wywołania zwrotnego.
 * Zbiór jest wspólny obu oknom, więc bez ogłoszenia drugie okno przerysowałoby
 * się dopiero przy najbliższej zmianie treści badania: ta sama nastawa
 * pokazywałaby w dwóch oknach dwie różne wartości, a powiązanie
 * źródło↔ustalenie brałoby to, czego w panelu ustaleń nie widać. Wywołanie
 * zwrotne jest więc obowiązkowe, nie domyślne.
 *
 * Ogłoszenie idzie wyłącznie po faktycznej zmianie. Odbiorcą jest przerysowanie
 * okien, a przerysowanie woła `ogranicz` — ogłoszenie bezwarunkowe zamknęłoby
 * pętlę bez końca.
 */
export interface ZaznaczeniePozycji {
  /** Zaznaczone identyfikatory w kolejności wykazu. */
  wybrane(): readonly string[];
  /** Czy pozycja jest zaznaczona. */
  czyWybrana(identyfikator: string): boolean;
  /** Przestawia zaznaczenie pozycji na przeciwne. */
  przelacz(identyfikator: string): void;
  /** Zdejmuje zaznaczenie z pozycji, których nie ma już w wykazie. */
  ogranicz(istniejace: readonly string[]): void;
  /** Zdejmuje całe zaznaczenie. */
  wyczysc(): void;
}

/**
 * @param oglos wywoływane po każdej faktycznej zmianie zbioru; przez nie
 *              nastawa dochodzi do wszystkich okien patrzących na ten sam zbiór
 */
export function utworzZaznaczenie(oglos: () => void): ZaznaczeniePozycji {
  const zbior = new Set<string>();

  return {
    wybrane: () => [...zbior],
    czyWybrana: (identyfikator) => zbior.has(identyfikator),

    przelacz(identyfikator) {
      if (zbior.has(identyfikator)) zbior.delete(identyfikator);
      else zbior.add(identyfikator);
      oglos();
    },

    ogranicz(istniejace) {
      const znane = new Set(istniejace);
      let zmienione = false;
      for (const identyfikator of [...zbior]) {
        if (znane.has(identyfikator)) continue;
        zbior.delete(identyfikator);
        zmienione = true;
      }
      if (zmienione) oglos();
    },

    wyczysc() {
      if (zbior.size === 0) return;
      zbior.clear();
      oglos();
    },
  };
}
