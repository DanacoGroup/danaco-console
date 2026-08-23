import type { Kanal } from '../protokol/kanal';
import { utworzOknoModeli, type OknoModeli } from './okno-modeli';

/**
 * Sekcja modeli, kont i tożsamości — punkt zbiorczy katalogu.
 *
 * Reszta aplikacji zna stąd jedną czynność:
 *
 *   import { otworzOknoModeli } from '../modele/indeks';
 *   otworzOknoModeli(rdzen.kanal);
 *
 * Okno jest jedno na klienta i żyje między otwarciami — z tego samego powodu,
 * dla którego router nie porzuca widoku opuszczonej trasy: powtórne otwarcie ma
 * wrócić do zakładki i kategorii, na której Operator skończył. Przy okazji nie
 * gubi się subskrypcja zdarzeń `account.changed` i `identity.changed`, więc
 * zmiana dokonana na innym urządzeniu konta dociera także wtedy, gdy okno jest
 * zamknięte.
 *
 * Wywołanie z innym kanałem — po ponownym połączeniu z rdzeniem — buduje okno
 * od nowa, żeby komendy nie szły przez transport, którego już nie ma.
 *
 * Sekcję można też osadzić bez ramy okna: `utworzSekcjeModeli` zwraca ten sam
 * byt bez `<dialog>` wokół niego — jedna implementacja, dwie oprawy.
 */

/** Okno zbudowane przy pierwszym otwarciu; `null` przed nim. */
let okno: OknoModeli | null = null;

/** Kanał, na którym zbudowano okno — podstawa rozpoznania zmiany połączenia. */
let osadzonyKanal: Kanal | null = null;

/** Otwiera okno modeli i zleca odczyt rejestrów oraz katalogów. */
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
