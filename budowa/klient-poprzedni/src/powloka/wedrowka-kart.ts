import type { KartaSesji } from './karta-sesji';

// Wędrujący fokus pasa zakładek — obsługa klawiatury pasa kart sesji, wzorzec ARIA.

/** Czynności pasa wywoływane klawiszem: wybór karty pod ogniskiem albo jej zamknięcie klawiszem Delete. */
export interface CzynnosciPasa {
  wybierz(id: string): void;
  zamknij(id: string): void;
}

export function przeniesWedrowke(
  karty: readonly KartaSesji[],
  id: string,
  zdarzenie: KeyboardEvent,
  czynnosci: CzynnosciPasa,
): void {
  const miejsce = karty.findIndex((karta) => karta.id === id);
  if (miejsce < 0) return;

  const cel = celKlawisza(karty, miejsce, zdarzenie.key);
  if (cel !== null) {
    zdarzenie.preventDefault();
    czynnosci.wybierz(cel.id);
    cel.ustawOgnisko();
    return;
  }

  if (zdarzenie.key === 'Enter' || zdarzenie.key === ' ') {
    zdarzenie.preventDefault();
    czynnosci.wybierz(id);
    return;
  }

  if (zdarzenie.key === 'Delete') {
    zdarzenie.preventDefault();
    czynnosci.zamknij(id);
  }
}

/** Karta wskazana klawiszem przenoszącym wybór ogniska; wartość pusta oznacza naciśnięcie innego klawisza. */
function celKlawisza(
  karty: readonly KartaSesji[],
  miejsce: number,
  klawisz: string,
): KartaSesji | null {
  if (klawisz === 'Home') return karty[0] ?? null;
  if (klawisz === 'End') return karty[karty.length - 1] ?? null;

  const krok = przesuniecie(klawisz);
  if (krok === 0 || karty.length === 0) return null;
  return karty[(miejsce + krok + karty.length) % karty.length] ?? null;
}

/** Przesunięcie wyboru wynikające z klawisza strzałki; zero oznacza naciśnięcie zupełnie innego klawisza. */
function przesuniecie(klawisz: string): number {
  if (klawisz === 'ArrowRight' || klawisz === 'ArrowDown') return 1;
  if (klawisz === 'ArrowLeft' || klawisz === 'ArrowUp') return -1;
  return 0;
}
