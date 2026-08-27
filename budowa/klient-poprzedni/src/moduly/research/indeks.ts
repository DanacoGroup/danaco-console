import type { OpisModulu } from '../../aplikacja/rejestr-modulow';
import { utworzModulResearch, type ModulResearch } from './modul-research';

/**
 * Stała opisuje moduł Research w rejestrze modułów, tworząc jego widok i podłączając rozłączenie modułu do zamknięcia widoku.
 */
export const MODUL: OpisModulu = {
  kod: 'research',

  utworzWidok(kanal) {
    const modul: ModulResearch = utworzModulResearch(kanal);
    return {
      element: modul.element,
      wczytaj: (idSesji) => modul.wczytaj(idSesji),
      zamknij: () => modul.rozlacz(),
    };
  },
};

export { utworzModulResearch } from './modul-research';
export type { ModulResearch } from './modul-research';
export { utworzStanBadania } from './stan-badania';
export type { StanBadania } from './stan-badania';
