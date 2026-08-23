import { nazwaModulu } from './etykiety-sterowania';
import { utworzListeWyboru, type OpcjaWyboru } from './lista-wyboru';
import type { RejestrModulow } from './rejestr-modulow';
import type { StanSterowania } from './stan-sterowania';
import type { ZmianaOkna } from './zmiana-okna';

const NAZWA = 'Moduł';

/**
 * Sterowanie modułem okna.
 *
 * Moduł jest parametrem okna, nie sesji, więc dwa okna jednej sesji mogą
 * pracować w dwóch różnych modułach.
 *
 * Wykaz pochodzi z katalogu modułów rdzenia, nie z `KnownModuleIds` — tamta
 * lista niesie środowiska (`talkin`, `workspace`, `codestudio`,
 * `multitaskingai`), a nie moduły.
 *
 * Moduł spoza katalogu pozostaje widoczny i wybieralny — katalog jest listą
 * informacyjną, nie bramą.
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

/** Katalog modułów rdzenia; pusty, dopóki rdzeń nie odpowie. */
function opcje(rejestr: RejestrModulow): OpcjaWyboru[] {
  return rejestr.moduly().map((modul) => ({
    wartosc: modul.code,
    nazwa: modul.name === '' ? nazwaModulu(modul.code) : modul.name,
  }));
}
