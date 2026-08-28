import type { PozycjaCzynnosciMenu } from './wiersz-czynnosci';

/**
 * Otwarcie w nowym oknie jest pierwszą pozycją sekcji czynności sesji: nie woła komendy rdzenia wprost, tylko podnosi liczbę gniazd sceny, a gniazda bierze z tej liczby układ okien.
 */
export interface PortNowegoOkna {
  /** Czy podniesienie liczby gniazd cokolwiek zmieni; pytanie zadane układowi, nie policzone tutaj. */
  wolneGniazdo(): boolean;
  /** Podnosi liczbę gniazd sceny o jedno. */
  naNoweOkno(): void;
}

/**
 * Buduje pozycję albo oddaje wartość pustą, gdy scena nie ma już wolnego gniazda; po wykonaniu przerysowuje sekcję zaraz po podniesieniu liczby.
 */
export function pozycjaNowegoOkna(
  port: PortNowegoOkna,
  poWykonaniu: () => void,
): PozycjaCzynnosciMenu | null {
  if (!port.wolneGniazdo()) return null;

  return {
    klucz: 'nowe-okno',
    nazwa: 'Otwórz w nowym oknie',
    przeznaczenie: 'Stawia obok kolejne okno tej sesji — z własną rozmową i własnymi panelami.',
    ikona: 'karta-okna',
    wykonaj: () => {
      port.naNoweOkno();
      poWykonaniu();
    },
  };
}
