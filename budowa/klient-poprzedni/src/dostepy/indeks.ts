/**
 * Sekcja dostępów i katalogu roboczego — punkt zbiorczy katalogu. Reszta
 * aplikacji zna stąd jedną czynność `otworzOknoDostepow`, a gospodarz z własnym
 * obszarem na ekranie bierze samą sekcję `utworzSekcjeDostepow`.
 */
import type { Kanal } from '../protokol/kanal';
import { utworzOknoDostepow, type OknoDostepow } from './okno-dostepow';



/**
 * Okno zbudowane przy pierwszym otwarciu; przed nim wartość `null`. Okno jest
 * jedno na klienta i żyje między otwarciami, więc powtórne otwarcie wraca do
 * zastanego stanu.
 */
let okno: OknoDostepow | null = null;

/**
 * Kanał, na którym zbudowano okno — podstawa rozpoznania zmiany połączenia;
 * kanał inny niż osadzony znaczy nowe połączenie z rdzeniem i wymusza budowę
 * okna od nowa.
 */
let osadzonyKanal: Kanal | null = null;

/**
 * Otwiera okno dostępów i wiąże je ze wskazanym oknem rozmowy; identyfikator
 * pusty jest poprawny, ponieważ wykaz punktów i katalog roboczy od okna nie
 * zależą.
 */
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
