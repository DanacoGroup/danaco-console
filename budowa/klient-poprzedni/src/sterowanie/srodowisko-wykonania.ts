import { ExecutionEnv } from '../../../shared/contract';
import { nazwaSrodowiska } from '../okno-komunikacji/etykiety-okna';
import { utworzListeWyboru, type OpcjaWyboru } from './lista-wyboru';
import type { StanSterowania } from './stan-sterowania';
import type { ZmianaOkna } from './zmiana-okna';

const NAZWA = 'Środowisko wykonania';

/**
 * Sterowanie zasięgiem wykonania modelu.
 *
 * Wyliczenie ma trzy rodzaje: urządzenie operatora, host rdzenia, host zdalny.
 * Nazwa konkretnego serwera nie jest wartością wyliczenia — wskazuje ją
 * sterowanie hostem, jako ustawienie poziomu okna. Dzięki temu dopisanie
 * kolejnego serwera nie wymaga zmiany kontraktu ani schematu bazy.
 *
 * Zasięg wykonania jest parametrem okna, nie właściwością wdrożenia: dwa okna
 * jednej sesji mogą pracować w dwóch różnych zasięgach równocześnie.
 */
export function utworzSterowanieSrodowiska(
  stan: StanSterowania,
  zmiana: ZmianaOkna,
): HTMLElement {
  const lista = utworzListeWyboru(NAZWA, (wartosc) => {
    zmiana.zastosuj(NAZWA, { executionEnv: rodzaj(wartosc, stan) });
  });

  stan.naZmiane((migawka) => lista.pokaz(opcje(), migawka.okno.executionEnv));
  lista.pokaz(opcje(), stan.migawka().okno.executionEnv);

  return lista.element;
}

/** Katalog zasięgów wykonania wprost z wyliczenia kontraktu. */
function opcje(): OpcjaWyboru[] {
  return Object.values(ExecutionEnv).map((wartosc) => ({
    wartosc,
    nazwa: nazwaSrodowiska(wartosc),
  }));
}

/** Wartość listy sprowadzona do rodzaju kontraktu; nierozpoznana zostawia stan bez zmiany. */
function rodzaj(wartosc: string, stan: StanSterowania): ExecutionEnv {
  const znaleziona = Object.values(ExecutionEnv).find((rodzaj) => rodzaj === wartosc);
  return znaleziona ?? stan.migawka().okno.executionEnv;
}
