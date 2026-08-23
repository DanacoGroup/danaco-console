/**
 * Nadajnik zdarzeń widoku — jedno zdarzenie, wielu odbiorców.
 *
 * Jedna odpowiedzialność: rozesłanie ładunku do subskrybentów i oddanie
 * odwołania subskrypcji. Pulpit nie sięga po magistralę połączenia, bo jego
 * zdarzenia są zdarzeniami widoku, nie kontraktu — powłoka decyduje, którą
 * komendę kontraktu z nich zbuduje.
 *
 * Błąd jednego odbiorcy nie odbiera zdarzenia pozostałym: wywołanie
 * biegnie w bloku ochronnym, a usterka trafia do konsoli, nie do przerwania
 * rozgłoszenia.
 */

/** Odwołanie subskrypcji — wywołanie odpina odbiorcę. */
export type Odpiecie = () => void;

/** Nadajnik jednego rodzaju zdarzenia. */
export interface Nadajnik<T> {
  /** Podpina odbiorcę; zwraca odpięcie. */
  sluchaj(odbiorca: (ladunek: T) => void): Odpiecie;
  /** Rozsyła ładunek do wszystkich odbiorców. */
  nadaj(ladunek: T): void;
  /** Odpina wszystkich odbiorców naraz. */
  rozlacz(): void;
}

/** Buduje nadajnik zdarzenia o zadanym ładunku. */
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
