import type { Kanal } from '../../protokol/kanal';
import { utworzMowaZrodlo, type MowaZrodlo } from './przybornik-mowa';
import { utworzPrzybornikZrodlo, type PrzybornikZrodlo } from './przybornik-zrodlo';
import { utworzUzycieZrodlo, type UzycieZrodlo } from './przybornik-uzycie';
import { utworzOsadzenieZrodel, type OsadzenieZrodel } from './osadzenie-zrodel';
import { utworzSchowekZrodlo, type SchowekZrodlo } from './schowek-zrodlo';

/**
 * Interfejs ZapleczePrzybornika oraz funkcja utworzZapleczePrzybornika składają źródła komend przybornika, schowka, mowy, osadzenia i użycia w jeden zestaw metod przekazywany do odcinka studia.
 */
export interface ZapleczePrzybornika
  extends PrzybornikZrodlo,
    SchowekZrodlo,
    MowaZrodlo,
    OsadzenieZrodel,
    UzycieZrodlo {}

export function utworzZapleczePrzybornika(kanal: Kanal): ZapleczePrzybornika {
  return {
    ...utworzPrzybornikZrodlo(kanal),
    ...utworzSchowekZrodlo(kanal),
    ...utworzMowaZrodlo(kanal),
    ...utworzOsadzenieZrodel(kanal),
    ...utworzUzycieZrodlo(kanal),
  };
}
