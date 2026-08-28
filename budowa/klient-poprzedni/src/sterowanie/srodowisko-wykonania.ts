import { ExecutionEnv } from '../../../shared/contract';
import { nazwaSrodowiska } from '../okno-komunikacji/etykiety-okna';
import { utworzListeWyboru, type OpcjaWyboru } from './lista-wyboru';
import type { StanSterowania } from './stan-sterowania';
import type { ZmianaOkna } from './zmiana-okna';

const NAZWA = 'Środowisko wykonania';

/** Sterowanie zasięgiem wykonania modelu: urządzenie operatora, host rdzenia albo host zdalny, jako parametr okna. */
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

/** Katalog zasięgów wykonania zbudowany wprost z wyliczenia kontraktu, bez ręcznego dopisywania nazw serwerów. */
function opcje(): OpcjaWyboru[] {
  return Object.values(ExecutionEnv).map((wartosc) => ({
    wartosc,
    nazwa: nazwaSrodowiska(wartosc),
  }));
}

/** Wartość listy sprowadzona do rodzaju wyliczenia kontraktu; wartość nierozpoznana zostawia stan okna bez zmiany. */
function rodzaj(wartosc: string, stan: StanSterowania): ExecutionEnv {
  const znaleziona = Object.values(ExecutionEnv).find((rodzaj) => rodzaj === wartosc);
  return znaleziona ?? stan.migawka().okno.executionEnv;
}
