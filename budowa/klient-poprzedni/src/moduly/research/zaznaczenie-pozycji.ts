/**
 * Zaznaczenie wielokrotne pozycji wykazu: zbiór zaznaczonych identyfikatorów wraz z jego
 * porządkowaniem po zmianie wykazu, wspólny dla dwóch okien.
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
 *              nastawa dochodzi do wszystkich okien korzystających z tego samego zbioru
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
