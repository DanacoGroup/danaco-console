import { WindowRole } from '../../../shared/contract';
import { nazwaRoli } from '../okno-komunikacji/etykiety-okna';
import { utworzListeWyboru, type OpcjaWyboru } from './lista-wyboru';
import type { StanSterowania } from './stan-sterowania';
import type { ZmianaOkna } from './zmiana-okna';

const NAZWA = 'Rola okna';

/**
 * Sterowanie rolą okna w pętli koordynator-wykonawca; rola przypisana jest do okna, nie do
 * całej sesji.
 */
export function utworzSterowanieRoli(
  stan: StanSterowania,
  zmiana: ZmianaOkna,
): HTMLElement {
  const lista = utworzListeWyboru(NAZWA, (wartosc) => {
    zmiana.zastosuj(NAZWA, { windowRole: rola(wartosc, stan) });
  });

  stan.naZmiane((migawka) => lista.pokaz(opcje(), migawka.okno.windowRole));
  lista.pokaz(opcje(), stan.migawka().okno.windowRole);

  return lista.element;
}

/**
 * Katalog ról okna wprost z wyliczenia typów zdefiniowanych w kontrakcie tej całej
 * aplikacji klienckiej.
 */
function opcje(): OpcjaWyboru[] {
  return Object.values(WindowRole).map((wartosc) => ({
    wartosc,
    nazwa: nazwaRoli(wartosc),
  }));
}

/**
 * Wartość listy wyboru sprowadzona do roli z kontraktu aplikacji; nierozpoznana zostawia
 * stan bez zmiany.
 */
function rola(wartosc: string, stan: StanSterowania): WindowRole {
  const znaleziona = Object.values(WindowRole).find((rola) => rola === wartosc);
  return znaleziona ?? stan.migawka().okno.windowRole;
}
