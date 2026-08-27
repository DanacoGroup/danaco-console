import type { Kanal } from '../protokol/kanal';
import { utworzOknoMobile, type OknoMobile } from './okno-mobile';

/** Punkt zbiorczy katalogu mobilnego centrum dowodzenia. */

/**
 * Okno zbudowane przy pierwszym otwarciu; `null` przed nim. Okno jest jedno na
 * klienta i żyje między otwarciami, więc powtórne otwarcie ponawia odczyt
 * zamiast budować okno od nowa.
 */
let okno: OknoMobile | null = null;

/**
 * Kanał, na którym zbudowano okno — podstawa rozpoznania zmiany połączenia.
 * Wywołanie z innym kanałem buduje okno na nowo, żeby komendy `mobile.*` nie
 * szły przez transport, którego już nie ma.
 */
let osadzonyKanal: Kanal | null = null;

/**
 * Otwiera okno mobilnego centrum dowodzenia i zleca odczyt stanu.
 *
 * Otwarcie nie czeka na rdzeń: okno pojawia się od razu, a odczyt
 * `mobile.status.get` i `mobile.process.list` dojeżdża do niego odpowiedzią.
 */
export function otworzOknoMobile(kanal: Kanal): OknoMobile {
  if (okno !== null && osadzonyKanal !== kanal) {
    okno.rozlacz();
    okno = null;
  }

  if (okno === null) {
    okno = utworzOknoMobile(kanal);
    osadzonyKanal = kanal;
  }

  okno.otworz();
  return okno;
}

export type { OknoMobile } from './okno-mobile';
