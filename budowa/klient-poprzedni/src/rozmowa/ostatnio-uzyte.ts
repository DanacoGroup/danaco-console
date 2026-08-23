/**
 * Rejestr ostatnio użytych pozycji wykazu po ukośniku — podstawa dla
 * szeregowania, które trzyma świeżo użyte pozycje bliżej wierzchu.
 *
 * Rejestr jest kliencki i sesyjny: żyje w pamięci okna, przeżywa otwarcie
 * i zamknięcie wykazu, nie przeżywa odświeżenia strony. Kontrakt takiego
 * pojęcia nie niesie — `tools.catalog.list` oddaje wykaz, nie historię
 * sięgania po niego — więc rejestr nie idzie do rdzenia żadną komendą i nie
 * udaje jego stanu.
 */

/** Ile kluczy rejestr pamięta; dalsze wypadają z ogona. */
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
