/**
 * Nadajnik zdarzeń widoku: rozsyła ładunek jednego zdarzenia do wielu odbiorców
 * i oddaje odwołanie subskrypcji. Usterka odbiorcy trafia do konsoli i nie
 * przerywa rozgłoszenia pozostałym odbiorcom.
 */

/**
 * Odwołanie subskrypcji. Wywołanie zwróconej czynności odpina odbiorcę
 * i jest jedynym sposobem wypisania pojedynczego odbiorcy z nadajnika.
 */
export type Odpiecie = () => void;

/**
 * Nadajnik jednego rodzaju zdarzenia, o ładunku ustalonym parametrem typu.
 * Niesie podpięcie odbiorcy, rozesłanie ładunku oraz odpięcie wszystkich naraz.
 */
export interface Nadajnik<T> {
  /** Podpina odbiorcę; zwraca odpięcie. */
  sluchaj(odbiorca: (ladunek: T) => void): Odpiecie;
  /** Rozsyła ładunek do wszystkich odbiorców. */
  nadaj(ladunek: T): void;
  /** Odpina wszystkich odbiorców naraz. */
  rozlacz(): void;
}

/**
 * Buduje nadajnik zdarzenia o zadanym ładunku. Nazwa zdarzenia służy wyłącznie
 * wpisowi konsoli zgłaszającemu usterkę odbiorcy i nie wpływa na rozgłoszenie.
 */
export function utworzNadajnik<T>(nazwa: string): Nadajnik<T> {
  const odbiorcy = new Set<(ladunek: T) => void>();

  return {
    sluchaj(odbiorca) {
      odbiorcy.add(odbiorca);
      return () => {
        odbiorcy.delete(odbiorca);
      };
    },
    nadaj(ladunek) {
      for (const odbiorca of [...odbiorcy]) {
        try {
          odbiorca(ladunek);
        } catch (usterka) {
          console.error(`Mission Control — odbiorca zdarzenia „${nazwa}" zgłosił usterkę`, usterka);
        }
      }
    },
    rozlacz() {
      odbiorcy.clear();
    },
  };
}
