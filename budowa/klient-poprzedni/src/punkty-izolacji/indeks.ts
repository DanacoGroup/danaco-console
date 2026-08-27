import type { Kanal } from '../protokol/kanal';
import { utworzOknoPunktowIzolacji, type OknoPunktowIzolacji } from './okno-punktow-izolacji';

// Okno Punktów Izolacji — punkt zbiorczy katalogu, jedna czynność otwarcia dla reszty aplikacji.

/** Okno zbudowane przy pierwszym otwarciu tego katalogu; wartość pusta oznacza stan sprzed jego budowy. */
let okno: OknoPunktowIzolacji | null = null;

/** Kanał, na którym zbudowano to okno — podstawa rozpoznania zmiany połączenia po ponownym jego wejściu. */
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
