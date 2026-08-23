import type { Kanal } from '../protokol/kanal';
import { utworzOknoDostepow, type OknoDostepow } from './okno-dostepow';

/**
 * Sekcja dostępów i katalogu roboczego — punkt zbiorczy katalogu.
 *
 * Reszta aplikacji zna stąd jedną czynność:
 *
 *   import { otworzOknoDostepow } from '../dostepy/indeks';
 *   otworzOknoDostepow(rdzen.kanal, idOknaRozmowy);
 *
 * Gospodarz, który ma własny obszar na ekranie i nie chce okna modalnego,
 * bierze zamiast tego samą sekcję (`utworzSekcjeDostepow`) i osadza jej
 * element u siebie — czynność `ustawOkno` wiąże ją wtedy z oknem rozmowy.
 *
 * Okno jest jedno na klienta i żyje między otwarciami — ta sama zasada, co
 * w oknie konfiguracji: powtórne otwarcie wraca do stanu, na którym się
 * skończyło, i nie gubi subskrypcji zdarzeń `access.*.changed`. Wywołanie
 * z innym kanałem, po ponownym połączeniu z rdzeniem, buduje okno od nowa,
 * żeby komendy nie szły przez transport, którego już nie ma.
 *
 * Nadanie dostępu żyje per okno, więc identyfikator okna podaje się przy
 * otwarciu. Otwarcie bez niego jest poprawne: wykaz punktów i katalog roboczy
 * nie zależą od okna, a w miejscu nadań sekcja mówi, na co czeka.
 */

/** Okno zbudowane przy pierwszym otwarciu; `null` przed nim. */
let okno: OknoDostepow | null = null;

/** Kanał, na którym zbudowano okno — podstawa rozpoznania zmiany połączenia. */
let osadzonyKanal: Kanal | null = null;

/** Otwiera okno dostępów i wiąże je ze wskazanym oknem rozmowy. */
export function otworzOknoDostepow(kanal: Kanal, oknoID = ''): OknoDostepow {
  if (okno !== null && osadzonyKanal !== kanal) {
    okno.rozlacz();
    okno = null;
  }

  if (okno === null) {
    okno = utworzOknoDostepow(kanal, oknoID);
    osadzonyKanal = kanal;
  } else if (oknoID !== '') {
    okno.sekcja.ustawOkno(oknoID);
  }

  okno.otworz();
  return okno;
}

export { utworzSekcjeDostepow } from './sekcja-dostepow';
export type { SekcjaDostepow, ZaleznosciSekcji } from './sekcja-dostepow';
export type { OknoDostepow } from './okno-dostepow';
export type { StanDostepow } from './stan-dostepow';
export { MASZYNY_CHRONIONE, ostrzezenieZapisu } from './ostrzezenie-zapisu';
export { KLUCZ_PODSTAWA, KLUCZ_WZORZEC, WZORZEC_DOMYSLNY } from './klucze-katalogu';
export { czyPowlokaNatywna, wskazKatalog } from './dialog-katalogu';
