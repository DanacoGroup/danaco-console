import { PermissionMode } from '../../../shared/contract';
import type { PozycjaMenu } from '../komponenty/menu-drzewo';
import {
  nazwaKrotkaTrybuUprawnien,
  nazwaTrybuUprawnien,
  opisTrybuUprawnien,
} from '../sterowanie/etykiety-sterowania';
import { utworzSterNastawy, type SterPaska } from './ster-nastawy';

// Tryb zatwierdzania to ster paska zlecenia, rozstrzygający, czy model pyta o zgodę przed krokami.

/** Nazwa zmiany w komunikacie do rdzenia — ta sama, którą wysyła lista trybów w kolumnie sterowania oknem. */
const NAZWA = 'Tryb uprawnień';

/** Zależności steru trybu zatwierdzania — wąskie i wstrzykiwane, obejmujące migawkę stanu oraz wysyłkę zmiany. */
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

/** Klucz pozycji sprowadzony do trybu kontraktu; nierozpoznany klucz zostawia bieżący stan uprawnień okna. */
function tryb(klucz: string, zaleznosci: ZaleznosciSteruUprawnien): PermissionMode {
  const znaleziony = Object.values(PermissionMode).find((wartosc) => wartosc === klucz);
  return znaleziony ?? zaleznosci.migawka().tryb;
}
