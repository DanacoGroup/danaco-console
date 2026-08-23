import type { Kanal } from '../../protokol/kanal';
import { utworzMowaZrodlo, type MowaZrodlo } from './przybornik-mowa';
import { utworzPrzybornikZrodlo, type PrzybornikZrodlo } from './przybornik-zrodlo';
import { utworzUzycieZrodlo, type UzycieZrodlo } from './przybornik-uzycie';
import { utworzOsadzenieZrodel, type OsadzenieZrodel } from './osadzenie-zrodel';
import { utworzSchowekZrodlo, type SchowekZrodlo } from './schowek-zrodlo';

/**
 * Jedno zaplecze dla całego odcinka znakowania, asystenta, schowka i osadzenia.
 *
 * ── Po co składanie w jedno ─────────────────────────────────────────────────
 * Odcinek woła pięć rodzin komend, których okno pracy dotąd nie znało:
 * `studio.annotation.*` i `studio.operation.*`, `clipboard.*`, `speech.*`,
 * `library.*` wraz z `browser.*` i `studio.ingest.url`, oraz `config.*` na
 * nastawy przybornika. Każda ma własne źródło ze sprawdzianem kształtu
 * odpowiedzi — a moduł ma dostać je jednym wywołaniem, żeby dołożenie szóstej
 * rodziny nie było zmianą w podpisie okna.
 *
 * Zaplecze nie ma stanu i nie buduje elementu: jest sumą pięciu warstw wywołań.
 * Nazwy metod noszą przedrostki swoich rodzin (`przybornik*`, `schowek*`,
 * `mowa*`, `osadzenie*`), więc suma nie ma ani jednej kolizji nazw.
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
