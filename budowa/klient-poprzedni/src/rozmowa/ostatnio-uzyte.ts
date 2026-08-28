/**
 * Rejestr ostatnio użytych pozycji wykazu po ukośniku, kliencki i sesyjny, żyjący w
 * pamięci okna do odświeżenia strony.
 */

/**
 * Ile kluczy rejestr pamięta w kolejności ostatniego użycia; kolejne, starsze wpisy
 * wypadają z jego ogona.
 */
const POJEMNOSC = 12;

export interface RejestrOstatnich {
  /** Klucze od najświeższego; wprost do `uszereguj(...)`. */
  klucze(): readonly string[];
  /** Zapisuje sięgnięcie po pozycję; powtórzenie przesuwa ją na wierzch. */
  zapamietaj(klucz: string): void;
}

export function utworzRejestrOstatnich(pojemnosc: number = POJEMNOSC): RejestrOstatnich {
  /** Najświeższy na początku. */
  let klucze: string[] = [];

  return {
    klucze: () => klucze,

    zapamietaj(klucz) {
      if (klucz === '') return;
      klucze = [klucz, ...klucze.filter((k) => k !== klucz)].slice(0, Math.max(1, pojemnosc));
    },
  };
}
