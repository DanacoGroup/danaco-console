import type { KartaSesji } from './karta-sesji';

/**
 * Wędrujący fokus pasa zakładek — obsługa klawiatury według wzorca ARIA.
 *
 * Jedna odpowiedzialność: przełożenie klawisza na czynność pasa. Strzałki
 * przenoszą wybór z zawijaniem, Home i End skaczą na krańce, Enter i spacja
 * wybierają kartę pod fokusem, Delete ją zamyka.
 *
 * Plik nie wie, co wybór i zamknięcie znaczą: w pasie związanym z rdzeniem
 * obie czynności są komendami kontraktu, w pasie samego widoku — zmianą
 * miejscową. Rozstrzyga to `karty-sesji.ts`, który podaje tu czynności.
 */

/** Czynności pasa wywoływane klawiszem. */
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

/** Karta wskazana klawiszem przenoszącym wybór; `null` przy innym klawiszu. */
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

/** Przesunięcie wyboru wynikające z klawisza strzałki; 0 znaczy inny klawisz. */
function przesuniecie(klawisz: string): number {
  if (klawisz === 'ArrowRight' || klawisz === 'ArrowDown') return 1;
  if (klawisz === 'ArrowLeft' || klawisz === 'ArrowUp') return -1;
  return 0;
}
