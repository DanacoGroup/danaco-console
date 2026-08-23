import type { PozycjaMenu } from '../komponenty/menu-drzewo';
import { utworzSterNastawy, type SterPaska } from './ster-nastawy';

/**
 * Katalog roboczy — ster paska zlecenia. Rozstrzyga, na jakich plikach
 * zlecenie się wykona.
 *
 * Katalogi są listą, nie pojedynczą wartością, więc pozycje wykazu są
 * przełącznikami, a nie wyborem jednokrotnym: okno pracuje na wszystkich
 * naraz, a zdjęcie jednego z nich nie jest wyborem innego. Każda zmiana idzie
 * komendą `window.update` z pełną listą po zmianie — kontrakt niesie
 * `workingDirs` jako całość.
 *
 * W menu nie ma dodawania, bo nowy katalog wskazuje się wpisaniem ścieżki,
 * a drzewo oddaje klucz istniejącej pozycji i pola tekstowego nie ma. Rdzeń
 * nie ma też komendy dającej wykaz katalogów do wyboru, więc gałąź „dostępne
 * katalogi” byłaby atrapą. Zamiast niej stoi stopka otwierająca kolumnę
 * sterowania, gdzie pole ścieżki stoi i działa: menu obsługuje podłączenie,
 * stopka prowadzi do rejestracji.
 *
 * Na uchwycie stoi nazwa ostatniego odcinka ścieżki, nie cała ścieżka: pasek
 * ma zostać jednym rzędem, a wielokropek ucinałby ścieżkę od końca, czyli od
 * jedynej części, która ją rozróżnia. Cała ścieżka stoi w opisie pozycji.
 */

/** Nazwa zmiany w komunikacie — ta sama, którą wysyła lista w kolumnie. */
const NAZWA = 'Katalogi robocze';

/** Etykieta uchwytu przy pustej liście — mówi o braku, niczego nie udaje. */
export const BRAK_KATALOGU = 'Bez katalogu';

/** Zależności steru — wąskie i wstrzykiwane. */
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
      // Zgaszony przełącznik zdejmuje katalog z okna. Katalogu, którego na
      // liście nie ma, nie da się tu dołożyć — ścieżkę trzeba wpisać, a od tego
      // jest stopka; wysyłka listy niezmienionej byłaby ruchem bez skutku.
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
