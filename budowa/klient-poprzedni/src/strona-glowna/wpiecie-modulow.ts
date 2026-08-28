import { Command } from '../../../shared/contract';
import type { Kanal } from '../protokol/kanal';
import { pozycjeModulowBezNawigacji } from './pozycje-modulow';
import type { StrefaModulow } from './strefa-modulow';

/** Wpięcie wykazu modułów bez nawigacji w kafle strony głównej pobiera komplet modułów jednym zapytaniem na wejście. */
export interface ZaleznosciWpieciaModulow {
  /** Kanał kontraktu — źródło wykazu modułów. */
  kanal: Kanal;
  /** Strefa kafli modułów poza nawigacją. */
  strefa: StrefaModulow;
  /** Przejście do pracy w module o wskazanym kodzie. */
  naPrzejscie(kodModulu: string): void;
}

export function wepnijModuly(zaleznosci: ZaleznosciWpieciaModulow): void {
  zaleznosci.kanal.wyslij(Command.ModuleList, {}, (wynik) => {
    const moduly = wynik.wynik?.modules;
    if (!wynik.udany || moduly === undefined) return;
    zaleznosci.strefa.ustaw(pozycjeModulowBezNawigacji(moduly));
  });

  zaleznosci.strefa.naWybor((pozycja) => zaleznosci.naPrzejscie(pozycja.kod));
}
