import type { Kanal } from '../protokol/kanal';
import { utworzOknoModeli, type OknoModeli } from './okno-modeli';

/**
 * Punkt zbiorczy katalogu modeli, kont i tożsamości: reszta aplikacji otwiera
 * stąd jedyne okno sekcji. Okno powstaje przy pierwszym otwarciu i trwa między
 * kolejnymi, a wartość `null` obowiązuje, dopóki żadne otwarcie nie nastąpiło.
 */
let okno: OknoModeli | null = null;

/**
 * Kanał, na którym zbudowano okno. Porównanie z kanałem podanym przy otwarciu
 * rozpoznaje ponowne połączenie z rdzeniem i nakazuje przebudowę okna, zanim
 * pierwsza komenda pójdzie przez transport już nieczynny.
 */
let osadzonyKanal: Kanal | null = null;

/**
 * Otwiera okno modeli i zleca odczyt rejestrów oraz katalogów. Zerwane
 * połączenie zamyka poprzednie okno, a odczyt rusza dopiero po tym, jak okno
 * stanie na kanale czynnym.
 */
export function otworzOknoModeli(kanal: Kanal): OknoModeli {
  if (okno !== null && osadzonyKanal !== kanal) {
    okno.rozlacz();
    okno = null;
  }

  if (okno === null) {
    okno = utworzOknoModeli(kanal);
    osadzonyKanal = kanal;
  }

  okno.otworz();
  return okno;
}

export { utworzSekcjeModeli } from './sekcja-modeli';
export type { OknoModeli } from './okno-modeli';
export type { SekcjaModeli } from './sekcja-modeli';
export type { StanKont } from './stan-kont';
export type { StanTozsamosci } from './stan-tozsamosci';
export type { WskazanieOsi } from './wybor-osi';
