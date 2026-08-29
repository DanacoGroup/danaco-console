/**
 * Kolejka ramek oczekujących na wysłanie w czasie rozłączenia.
 *
 * Kolejka jest nieograniczona świadomie: odrzucenie ramki po przekroczeniu
 * limitu gubiłoby żądanie Operatora bez śladu widocznego dla warstwy, która
 * je złożyła.
 */
export interface KolejkaWychodzaca<T> {
  /** Dokłada element na koniec kolejki. */
  dodaj(element: T): void;
  /** Zwraca i usuwa całą zawartość kolejki w kolejności dodania. */
  wydajWszystko(): T[];
  /** Liczba elementów oczekujących. */
  rozmiar(): number;
}

export function utworzKolejkeWychodzaca<T>(): KolejkaWychodzaca<T> {
  let oczekujace: T[] = [];

  return {
    dodaj(element) {
      oczekujace.push(element);
    },

    wydajWszystko() {
      const wydane = oczekujace;
      oczekujace = [];
      return wydane;
    },

    rozmiar() {
      return oczekujace.length;
    },
  };
}
