import type { PozycjaCzynnosciMenu } from './wiersz-czynnosci';

/**
 * `Otwórz w nowym oknie` — pierwsza pozycja sekcji czynności sesji.
 *
 * Pozycja nie woła komendy rdzenia wprost, tylko podnosi liczbę gniazd sceny.
 * Gniazda bierze z tej liczby `uklad-okien.ts`, a `aplikacja/scena-sesji.ts`
 * zamawia okno rdzenia dopiero wtedy, gdy gniazdo wejdzie na scenę; okna
 * zakładanego z boku nic w układzie by nie zauważyło.
 *
 * Sufit gniazd bywa węższy od maksimum sceny (figura modułu —
 * `figura-modulu.ts`), więc pytamy o niego przy każdym rysowaniu wiersza. Przy
 * suficie osiągniętym pozycja nie powstaje: lista jest o jedną pozycję krótsza
 * zamiast pokazywać wiersz wygaszony. Skrótu klawiszowego pozycja nie ma.
 */

/** Dojście do liczby gniazd sceny; wypełnia je układ okien. */
export interface PortNowegoOkna {
  /**
   * Czy podniesienie liczby gniazd cokolwiek zmieni.
   *
   * Pytanie zadane układowi, nie policzone tutaj: sufit zależy od figury modułu
   * sceny, a ta przestawia się z rdzenia. Druga rachuba tej samej granicy
   * rozjechałaby się z pierwszą przy zmianie modułu.
   */
  wolneGniazdo(): boolean;
  /** Podnosi liczbę gniazd sceny o jedno. */
  naNoweOkno(): void;
}

/**
 * Buduje pozycję albo oddaje `null`, gdy scena nie ma już wolnego gniazda.
 *
 * `poWykonaniu` przerysowuje sekcję zaraz po podniesieniu liczby: menu zostaje
 * rozwinięte po naciśnięciu, a wiersz, który właśnie zajął ostatnie wolne
 * gniazdo, ma z tego menu zniknąć.
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
