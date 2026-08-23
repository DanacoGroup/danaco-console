import { Command } from '../../../shared/contract';
import type { Kanal } from '../protokol/kanal';
import { pozycjeModulowBezNawigacji } from './pozycje-modulow';
import type { StrefaModulow } from './strefa-modulow';

/**
 * Wpięcie wykazu modułów bez nawigacji w kafle strony głównej.
 *
 * `module.list` bez zawężenia oddaje komplet modułów platformy wraz z polem
 * `environmentCodes`, którego pustka znaczy „moduł dostępny wyłącznie ze strony
 * głównej". Jedno pytanie wystarcza — strona główna wstaje raz na wejście.
 *
 * Gdy rdzeń nie odpowie, strefa zostaje ukryta zamiast stać z nagłówkiem nad
 * pustką: brak odpowiedzi nie jest wykazem pustym.
 *
 * Przejście idzie z samym kodem modułu, bez środowiska — moduł w żadnym nie
 * stoi, więc powłoka otwiera go z katalogu `module.list`, z pominięciem macierzy
 * widoczności.
 */

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
