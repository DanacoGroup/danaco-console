import { PermissionMode } from '../../../shared/contract';
import { nazwaTrybuUprawnien } from './etykiety-sterowania';
import { utworzListeWyboru, type OpcjaWyboru } from './lista-wyboru';
import type { StanSterowania } from './stan-sterowania';
import type { ZmianaOkna } from './zmiana-okna';

const NAZWA = 'Tryb uprawnień';

/** Sterowanie trybem uprawnień okna, z wyliczenia kontraktu, bez wyszarzania ani ukrywania żadnej pozycji. */
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

/** Katalog trybów uprawnień zbudowany wprost z wyliczenia kontraktu, bez własnego katalogu po stronie klienta. */
function opcje(): OpcjaWyboru[] {
  return Object.values(PermissionMode).map((wartosc) => ({
    wartosc,
    nazwa: nazwaTrybuUprawnien(wartosc),
  }));
}

/** Wartość listy sprowadzona do trybu wyliczenia kontraktu; wartość nierozpoznana zostawia stan okna bez zmiany. */
function tryb(wartosc: string, stan: StanSterowania): PermissionMode {
  const znaleziony = Object.values(PermissionMode).find((tryb) => tryb === wartosc);
  return znaleziony ?? stan.migawka().okno.permissionMode;
}
