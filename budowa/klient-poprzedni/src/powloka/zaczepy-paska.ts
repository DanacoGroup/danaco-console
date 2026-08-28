import type { PozycjaUstawienia } from '../strona-glowna/pozycje-ustawien';

// Zaczepy paska — czynności, których pasek nie zna, podane mu z zewnątrz przez inne warstwy tej powłoki.
export interface ZaczepyPaska {
  // Otwiera okno platformy odpowiadające pozycji ustawień, jedna czynność na wszystkie siedem pozycji.
  otworzUstawienie?(pozycja: PozycjaUstawienia): void;

  // Otwiera albo zamyka centrum powiadomień; osobny zaczep, a nie pozycja ustawień listwy platformy.
  otworzPowiadomienia?(): void;
}

/** Gdzie postawić zdanie komunikatu paska; bez podanego odbiorcy zdanie przepada tutaj zawsze świadomie. */
export type OdbiorcaZdania = (tytul: string, tresc: string) => void;

/** Zdanie o braku drogi — jedno wspólne zdanie na cały ten pasek, wypowiadane przez trzy osobne kontrolki. */
export function zdanieBezDrogi(nazwaOkna: string): string {
  return (
    `Ten pasek nie dostał drogi do okna „${nazwaOkna}". Okno jest zbudowane `
    + 'i osiągalne z listwy strony głównej — brakuje wyłącznie zaczepu '
    + 'w widoku środowiska, nie samego okna.'
  );
}

/**
 * Wykonuje zaczep albo mówi, dlaczego go nie ma; odczyt zaczepu jest tutaj, w chwili naciśnięcia
 * kontrolki, nie w chwili jej montażu.
 */
export function wykonajZaczep(
  zaczepy: ZaczepyPaska,
  pozycja: PozycjaUstawienia,
  naKomunikat?: OdbiorcaZdania,
): void {
  const otworz = zaczepy.otworzUstawienie;
  if (otworz === undefined) {
    naKomunikat?.(pozycja.nazwa, zdanieBezDrogi(pozycja.nazwa));
    return;
  }
  otworz(pozycja);
}
