import type { Kanal } from '../protokol/kanal';
import { utworzOknoPunktowIzolacji, type OknoPunktowIzolacji } from './okno-punktow-izolacji';

/**
 * Okno Punktów Izolacji — punkt zbiorczy katalogu.
 *
 * Reszta aplikacji zna stąd jedną czynność:
 *
 *   import { otworzOknoPunktowIzolacji } from '../punkty-izolacji/indeks';
 *   otworzOknoPunktowIzolacji(rdzen.kanal);
 *
 * Okno jest jedno na klienta i żyje między otwarciami — powtórne otwarcie
 * wraca do obszaru, na którym Operator skończył, tak jak Okno Konfiguracji
 * (`konfiguracja/indeks.ts`), którego ten plik jest wierną kopią wzorca.
 *
 * Kanał podajemy przy pierwszym otwarciu. Wywołanie z innym kanałem —
 * po ponownym połączeniu z rdzeniem — buduje okno na nowo, żeby komendy
 * `isolation.*` nie szły przez transport, którego już nie ma.
 */

/** Okno zbudowane przy pierwszym otwarciu; `null` przed nim. */
let okno: OknoPunktowIzolacji | null = null;

/** Kanał, na którym zbudowano okno — podstawa rozpoznania zmiany połączenia. */
let osadzonyKanal: Kanal | null = null;

/**
 * Otwiera okno Punktów Izolacji.
 *
 * Otwarcie nie czeka na rdzeń: okno staje od razu, a każdy obszar
 * sam woła swoje komendy `isolation.*` i sam melduje odmowę, gdy rdzeń odmówi.
 */
export function otworzOknoPunktowIzolacji(kanal: Kanal): OknoPunktowIzolacji {
  if (okno !== null && osadzonyKanal !== kanal) {
    okno.rozlacz();
    okno = null;
  }

  if (okno === null) {
    okno = utworzOknoPunktowIzolacji(kanal);
    osadzonyKanal = kanal;
  }

  okno.otworz();
  return okno;
}

export type { OknoPunktowIzolacji } from './okno-punktow-izolacji';
export type { KodObszaruIzolacji, ObszarIzolacji, WpisObszaru, ZaleznosciObszaru } from './obszary';
