import { PermissionMode } from '../../../shared/contract';
import type { PozycjaMenu } from '../komponenty/menu-drzewo';
import {
  nazwaKrotkaTrybuUprawnien,
  nazwaTrybuUprawnien,
  opisTrybuUprawnien,
} from '../sterowanie/etykiety-sterowania';
import { utworzSterNastawy, type SterPaska } from './ster-nastawy';

/**
 * Tryb zatwierdzania — ster paska zlecenia.
 *
 * Rozstrzyga, czy model pyta o zgodę przed krokami. Wartości odpowiadają
 * dosłownie przełącznikowi `--permission-mode` kanału głównego i pochodzą
 * z wyliczenia kontraktu — pasek nie prowadzi własnego katalogu trybów, tak
 * samo jak nie prowadzi go kolumna sterowania.
 *
 * Żadna pozycja nie jest ukryta ani wyszarzona, w tym pominięcie kontroli
 * uprawnień: jedyną kontrolą dostępu jest uwierzytelnianie, a interfejs nie
 * stawia blokad. Skutek wyboru mówi opis pozycji — to jest właściwe miejsce na
 * ostrzeżenie, nie odebranie kliknięcia.
 */

/** Nazwa zmiany w komunikacie — ta sama, którą wysyła lista w kolumnie. */
const NAZWA = 'Tryb uprawnień';

/** Zależności steru — wąskie i wstrzykiwane. */
export interface ZaleznosciSteruUprawnien {
  /** Tryb uprawnień okna ze stanu potwierdzonego przez rdzeń. */
  migawka(): { tryb: PermissionMode };
  /** Wysyła zmianę pól okna; odrzucenie niesie zdanie odmowy rdzenia. */
  zastosuj(nazwa: string, zmiana: { permissionMode: PermissionMode }): Promise<void>;
}

export function utworzSterUprawnien(
  zaleznosci: ZaleznosciSteruUprawnien,
): SterPaska {
  const ster = utworzSterNastawy({
    nastawa: NAZWA,
    ikona: 'tarcza',
    wykonaj: (klucz) =>
      zaleznosci.zastosuj(NAZWA, { permissionMode: tryb(klucz, zaleznosci) }),
    odswiez: () => odswiez(),
  });

  function odswiez(): void {
    const biezacy = zaleznosci.migawka().tryb;
    const drzewo: PozycjaMenu[] = Object.values(PermissionMode).map((wartosc) => ({
      rodzaj: 'wybor',
      klucz: wartosc,
      nazwa: nazwaTrybuUprawnien(wartosc),
      opis: opisTrybuUprawnien(wartosc),
      wybrany: wartosc === biezacy,
    }));
    ster.ustaw(nazwaKrotkaTrybuUprawnien(biezacy), drzewo);
  }

  odswiez();
  return { element: ster.element, odswiez };
}

/** Klucz pozycji sprowadzony do trybu kontraktu; nierozpoznany zostawia stan. */
function tryb(klucz: string, zaleznosci: ZaleznosciSteruUprawnien): PermissionMode {
  const znaleziony = Object.values(PermissionMode).find((wartosc) => wartosc === klucz);
  return znaleziony ?? zaleznosci.migawka().tryb;
}
