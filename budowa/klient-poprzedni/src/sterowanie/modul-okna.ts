import { nazwaModulu } from './etykiety-sterowania';
import { utworzListeWyboru, type OpcjaWyboru } from './lista-wyboru';
import type { RejestrModulow } from './rejestr-modulow';
import type { StanSterowania } from './stan-sterowania';
import type { ZmianaOkna } from './zmiana-okna';

const NAZWA = 'Moduł';

/**
 * Sterowanie modułem okna; moduł jest parametrem tego okna, nie całej sesji, więc dwa okna
 * mogą pracować w różnych modułach.
 */
export function utworzSterowanieModulu(
  stan: StanSterowania,
  zmiana: ZmianaOkna,
  rejestr: RejestrModulow,
): HTMLElement {
  const lista = utworzListeWyboru(NAZWA, (wartosc) => {
    zmiana.zastosuj(NAZWA, { moduleId: wartosc });
  });

  const odswiezWidok = (): void => lista.pokaz(opcje(rejestr), stan.migawka().okno.moduleId);

  stan.naZmiane((migawka) => lista.pokaz(opcje(rejestr), migawka.okno.moduleId));
  rejestr.naZmiane(odswiezWidok);
  odswiezWidok();
  rejestr.odswiez();

  return lista.element;
}

/**
 * Katalog modułów rdzenia, pusty dopóki rdzeń nie odpowie na zapytanie tego klienta o
 * dostępne moduły pracy.
 */
function opcje(rejestr: RejestrModulow): OpcjaWyboru[] {
  return rejestr.moduly().map((modul) => ({
    wartosc: modul.code,
    nazwa: modul.name === '' ? nazwaModulu(modul.code) : modul.name,
  }));
}
