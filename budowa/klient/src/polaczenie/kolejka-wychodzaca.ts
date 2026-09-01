/**
 * Kolejka ramek oczekujących na wysłanie w czasie rozłączenia.
 *
 * Sufit jest dobrowolny: kolejka bez sufitu rośnie bez granicy, kolejka
 * z sufitem oddaje wołającemu ramkę najstarszą, którą nowa wyparła. Wołający
 * ogłasza ją porzuconą — żądanie nie znika bez śladu widocznego dla warstwy,
 * która je złożyła.
 */
export interface KolejkaWychodzaca<T> {
  /** Dokłada element na koniec kolejki; zwraca element najstarszy, gdy sufit go wyparł. */
  dodaj(element: T): T | undefined;
  /** Zwraca i usuwa całą zawartość kolejki w kolejności dodania. */
  wydajWszystko(): T[];
  /** Liczba elementów oczekujących. */
  rozmiar(): number;
}

export function utworzKolejkeWychodzaca<T>(
  sufit: number = Number.POSITIVE_INFINITY,
): KolejkaWychodzaca<T> {
  let oczekujace: T[] = [];

  return {
    dodaj(element) {
      oczekujace.push(element);
      if (oczekujace.length <= sufit) return undefined;
      return oczekujace.shift();
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
