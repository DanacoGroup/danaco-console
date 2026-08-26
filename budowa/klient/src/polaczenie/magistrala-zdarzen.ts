/** Odbiorca powiadomienia o zdarzeniu. */
export type Sluchacz<T> = (dane: T) => void;

/** Odwołanie subskrypcji. Wywołanie odłącza słuchacza. */
export type Odsubskrybuj = () => void;

/**
 * Magistrala zdarzeń jednego rodzaju.
 *
 * Rozgłaszanie jest odporne na błąd pojedynczego słuchacza — wyjątek jednego
 * odbiorcy nie przerywa powiadamiania pozostałych (błąd dotyczy
 * wyłącznie bieżącego wywołania).
 */
export interface Magistrala<T> {
  /** Rejestruje słuchacza i zwraca odwołanie subskrypcji. */
  subskrybuj(sluchacz: Sluchacz<T>): Odsubskrybuj;
  /** Rozgłasza dane do wszystkich zarejestrowanych słuchaczy. */
  oglos(dane: T): void;
  /** Liczba aktywnych słuchaczy. */
  liczbaSluchaczy(): number;
}

export function utworzMagistrale<T>(): Magistrala<T> {
  const sluchacze = new Set<Sluchacz<T>>();

  return {
    subskrybuj(sluchacz) {
      sluchacze.add(sluchacz);
      return () => {
        sluchacze.delete(sluchacz);
      };
    },

    oglos(dane) {
      for (const sluchacz of [...sluchacze]) {
        try {
          sluchacz(dane);
        } catch (blad) {
          console.error('[połączenie] błąd słuchacza', blad);
        }
      }
    },

    liczbaSluchaczy() {
      return sluchacze.size;
    },
  };
}
