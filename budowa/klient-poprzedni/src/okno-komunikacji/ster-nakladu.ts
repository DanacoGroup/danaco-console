import type { PozycjaMenu } from '../komponenty/menu-drzewo';
import {
  nazwaKrotkaNakladu,
  nazwaNakladu,
  opisNakladu,
} from '../sterowanie/etykiety-sterowania';
import { KluczUstawieniaOkna, STOPNIE_NAKLADU } from '../sterowanie/klucze-ustawien';
import { utworzSterNastawy, type SterPaska } from './ster-nastawy';

// Wysiłek to ster paska zlecenia. W pasku stoi menu, a nie suwak: etykieta ma być bieżącą wartością.

/** Nazwa zmiany w komunikacie do rdzenia — ta sama, którą wysyła suwak wysiłku w kolumnie sterowania oknem. */
const NAZWA = 'Nakład rozumowania';

/**
 * Przedrostek klucza pozycji.
 *
 * Klucz pusty należy w tym mechanizmie do stopki (`komponenty/menu-drzewo.ts`
 * → `stopka()`), a stopień „bez wskazania" jest napisem pustym. Przedrostek
 * rozdziela jedno od drugiego.
 */
const PRZEDROSTEK = 'naklad:';

/** Zależności steru wysiłku — wąskie i wstrzykiwane, obejmujące stopień nakładu oraz zapis ustawienia okna. */
export interface ZaleznosciSteruNakladu {
  /** Stopień nakładu ze stanu potwierdzonego przez rdzeń. */
  migawka(): { naklad: string };
  /** Zapisuje ustawienie poziomu okna; odrzucenie niesie zdanie odmowy. */
  zapisz(nazwa: string, klucz: string, wartosc: unknown): Promise<void>;
}

export function utworzSterNakladu(zaleznosci: ZaleznosciSteruNakladu): SterPaska {
  const ster = utworzSterNastawy({
    nastawa: 'Wysiłek',
    ikona: 'aktywnosc',
    wykonaj: (klucz) =>
      zaleznosci.zapisz(
        NAZWA,
        KluczUstawieniaOkna.NakladRozumowania,
        klucz.slice(PRZEDROSTEK.length),
      ),
    odswiez: () => odswiez(),
  });

  function odswiez(): void {
    const stopien = zaleznosci.migawka().naklad;
    const drzewo: PozycjaMenu[] = STOPNIE_NAKLADU.map((pozycja) => ({
      rodzaj: 'wybor',
      klucz: `${PRZEDROSTEK}${pozycja}`,
      // Nazwa pełna w wykazie, krótka na uchwycie — oba napisy pochodzą z jednego słownika.
      nazwa: nazwaNakladu(pozycja),
      opis: opisNakladu(pozycja),
      wybrany: pozycja === stopien,
    }));
    ster.ustaw(nazwaKrotkaNakladu(stopien), drzewo);
  }

  odswiez();
  return { element: ster.element, odswiez };
}
