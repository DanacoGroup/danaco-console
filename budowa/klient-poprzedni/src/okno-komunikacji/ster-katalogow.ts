import type { PozycjaMenu } from '../komponenty/menu-drzewo';
import { utworzSterNastawy, type SterPaska } from './ster-nastawy';

// Katalog roboczy to ster paska zlecenia rozstrzygający, na jakich plikach zlecenie się wykona.

/** Nazwa zmiany w komunikacie do rdzenia — ta sama, którą wysyła lista katalogów w kolumnie sterowania. */
const NAZWA = 'Katalogi robocze';

/** Etykieta uchwytu przy pustej liście katalogów — mówi wprost o braku, niczego nie udaje operatorowi wcale. */
export const BRAK_KATALOGU = 'Bez katalogu';

/** Zależności steru katalogów — wąskie i wstrzykiwane, obejmujące migawkę stanu, wysyłkę zmiany i otwarcie kolumny. */
export interface ZaleznosciSteruKatalogow {
  /** Katalogi robocze okna ze stanu potwierdzonego przez rdzeń. */
  migawka(): { katalogi: readonly string[] };
  /** Wysyła zmianę pól okna; odrzucenie niesie zdanie odmowy rdzenia. */
  zastosuj(nazwa: string, zmiana: { workingDirs: string[] }): Promise<void>;
  /** Otwiera kolumnę sterowania okna — jedyne miejsce z polem ścieżki. */
  otworzSterowanie(): void;
}

export function utworzSterKatalogow(
  zaleznosci: ZaleznosciSteruKatalogow,
): SterPaska {
  const ster = utworzSterNastawy({
    nastawa: 'Katalog roboczy',
    ikona: 'folder',
    stopka: {
      nazwa: 'Katalogi robocze w sterowaniu okna',
      opis:
        'Dodanie katalogu wymaga wpisania ścieżki, a menu przyjmuje wybór, nie tekst. ' +
        'Pole ścieżki stoi w kolumnie sterowania — ta pozycja ją otwiera.',
      ikona: 'ustawienia',
      wykonaj: () => zaleznosci.otworzSterowanie(),
    },
    wykonaj: (klucz) => {
      const katalogi = [...zaleznosci.migawka().katalogi];
      // Zgaszony przełącznik zdejmuje katalog; dołożyć nowy można tylko wpisaniem ścieżki w stopce.
      return zaleznosci.zastosuj(NAZWA, {
        workingDirs: katalogi.filter((katalog) => katalog !== klucz),
      });
    },
    odswiez: () => odswiez(),
  });

  function odswiez(): void {
    const katalogi = zaleznosci.migawka().katalogi;
    const drzewo: PozycjaMenu[] = katalogi.map((katalog) => ({
      rodzaj: 'przelacznik',
      klucz: katalog,
      nazwa: ostatniOdcinek(katalog),
      opis: `${katalog} — zgaszenie zdejmuje ten katalog z okna.`,
      wlaczony: true,
    }));
    ster.ustaw(etykieta(katalogi), drzewo);
  }

  odswiez();
  return { element: ster.element, odswiez };
}

/**
 * Wartość na uchwyt: nazwa pierwszego katalogu, a przy kilku — z licznikiem
 * pozostałych. Bez licznika trzy katalogi wyglądałyby w pasku tak samo jak
 * jeden.
 */
function etykieta(katalogi: readonly string[]): string {
  const pierwszy = katalogi[0];
  if (pierwszy === undefined) return BRAK_KATALOGU;
  const reszta = katalogi.length - 1;
  return reszta === 0
    ? ostatniOdcinek(pierwszy)
    : `${ostatniOdcinek(pierwszy)} +${reszta}`;
}

/**
 * Ostatni odcinek ścieżki — obie kreski, bo katalog roboczy bywa windowsowy.
 * Ścieżka zakończona kreską (`/praca/`) oddaje odcinek sprzed niej, a nie pustkę.
 */
function ostatniOdcinek(sciezka: string): string {
  const odcinki = sciezka.split(/[/\\]/).filter((odcinek) => odcinek.length > 0);
  return odcinki[odcinki.length - 1] ?? sciezka;
}
