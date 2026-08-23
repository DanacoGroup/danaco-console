import { PermissionMode } from '../../../shared/contract';
import { nazwaTrybuUprawnien } from './etykiety-sterowania';
import { utworzListeWyboru, type OpcjaWyboru } from './lista-wyboru';
import type { StanSterowania } from './stan-sterowania';
import type { ZmianaOkna } from './zmiana-okna';

const NAZWA = 'Tryb uprawnień';

/**
 * Sterowanie trybem uprawnień okna.
 *
 * Wartości odpowiadają dosłownie przełącznikowi `--permission-mode` kanału
 * głównego i pochodzą z wyliczenia kontraktu — komplet nie prowadzi
 * własnego katalogu trybów.
 *
 * Żadna pozycja nie jest wyszarzona ani ukryta, w tym pominięcie kontroli
 * uprawnień: kontrolą dostępu jest uwierzytelnianie, a nie ta lista. Tryb jest
 * parametrem okna, więc okno planistyczne i okno wykonawcze mogą pracować obok
 * siebie.
 */
export function utworzSterowanieUprawnien(
  stan: StanSterowania,
  zmiana: ZmianaOkna,
): HTMLElement {
  const lista = utworzListeWyboru(NAZWA, (wartosc) => {
    zmiana.zastosuj(NAZWA, { permissionMode: tryb(wartosc, stan) });
  });

  stan.naZmiane((migawka) => lista.pokaz(opcje(), migawka.okno.permissionMode));
  lista.pokaz(opcje(), stan.migawka().okno.permissionMode);

  return lista.element;
}

/** Katalog trybów uprawnień wprost z wyliczenia kontraktu. */
function opcje(): OpcjaWyboru[] {
  return Object.values(PermissionMode).map((wartosc) => ({
    wartosc,
    nazwa: nazwaTrybuUprawnien(wartosc),
  }));
}

/** Wartość listy sprowadzona do trybu kontraktu; nierozpoznana zostawia stan bez zmiany. */
function tryb(wartosc: string, stan: StanSterowania): PermissionMode {
  const znaleziony = Object.values(PermissionMode).find((tryb) => tryb === wartosc);
  return znaleziony ?? stan.migawka().okno.permissionMode;
}
