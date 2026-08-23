/** Polityka odstępu między kolejnymi próbami połączenia. */
export interface PolitykaPonawiania {
  /** Odstęp w milisekundach przed próbą o podanym numerze (od 1). */
  opoznienie(numerProby: number): number;
}

/**
 * Ponawianie wykładnicze z górnym pułapem odstępu i rozproszeniem losowym.
 *
 * Liczba prób nie jest ograniczona: klient ponawia, dopóki nie uzyska
 * połączenia, a brak rdzenia nie wstrzymuje pracy interfejsu.
 */
export function wykladniczePonawianie(bazaMs = 500, pulapMs = 15_000): PolitykaPonawiania {
  return {
    opoznienie(numerProby) {
      const proba = Math.max(1, numerProby);
      const rosnace = bazaMs * 2 ** (proba - 1);
      const podstawa = Math.min(rosnace, pulapMs);
      // Rozproszenie 0–25% zapobiega zbieganiu się prób wielu okien naraz.
      return Math.round(podstawa * (1 + Math.random() * 0.25));
    },
  };
}
