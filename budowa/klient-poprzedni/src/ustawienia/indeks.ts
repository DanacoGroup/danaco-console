import type { Kanal } from '../protokol/kanal';
import { utworzOknoUstawien, type OknoUstawien } from './okno-ustawien';

/**
 * Okno Ustawień — punkt zbiorczy katalogu.
 *
 * Reszta aplikacji zna stąd jedną czynność:
 *
 *   import { otworzOknoUstawien } from '../ustawienia/indeks';
 *   otworzOknoUstawien(rdzen.kanal);
 *
 * Okno jest jedno na klienta i żyje między otwarciami — z tego samego powodu,
 * dla którego router nie porzuca widoku opuszczonej trasy: powtórne otwarcie
 * ma wrócić do sekcji, na której Operator skończył, a nie zaczynać od
 * początku, i nie gubić stanu (np. niezapisanego formularza) sekcji, których
 * Operator akurat nie ogląda.
 *
 * Kanał podajemy przy pierwszym otwarciu. Wywołanie z innym kanałem — po
 * ponownym połączeniu z rdzeniem — buduje okno na nowo, żeby komendy nie
 * szły przez transport, którego już nie ma.
 */

/** Okno zbudowane przy pierwszym otwarciu; `null` przed nim. */
let okno: OknoUstawien | null = null;

/** Kanał, na którym zbudowano okno — podstawa rozpoznania zmiany połączenia. */
let osadzonyKanal: Kanal | null = null;

/**
 * Otwiera Okno Ustawień i zleca odczyt sekcji czynnej.
 *
 * Otwarcie nie czeka na rdzeń: okno pojawia się od razu, a każda sekcja
 * dostaje wynik swojego odczytu odpowiedzią.
 */
export function otworzOknoUstawien(kanal: Kanal): OknoUstawien {
  if (okno !== null && osadzonyKanal !== kanal) {
    okno.rozlacz();
    okno = null;
  }

  if (okno === null) {
    okno = utworzOknoUstawien(kanal);
    osadzonyKanal = kanal;
  }

  okno.otworz();
  return okno;
}

export type { OknoUstawien } from './okno-ustawien';
export type { SekcjaUstawien, KodSekcjiUstawien, OpisSekcjiUstawien } from './sekcje';
export { REJESTR_SEKCJI_USTAWIEN } from './sekcje';
export type { StanUstawien } from './stan-ustawien';

/**
 * Most motywu wychodzi z katalogu osobno od okna.
 *
 * Nastawa motywu obowiązuje także wtedy, gdy Okno Ustawień stoi zamknięte —
 * a jest zamknięte niemal zawsze, bo okno jest przywoływane, nie rysowane
 * z urzędu. Gdyby most żył wyłącznie wewnątrz okna, zmiana motywu wykonana
 * w drugim oknie albo na drugim urządzeniu dolatywałaby dopiero po otwarciu
 * Ustawień.
 *
 * Dlatego `podepnijMostMotywu(kanal)` woła się raz przy wiązaniu gniazda
 * (`aplikacja/`), a sekcja „Wygląd i język" bierze ten sam, już podpięty most.
 */
export { podepnijMostMotywu, odepnijMostMotywu, KLUCZ_MOTYWU } from './most-motywu';
export type { MostMotywu, WyborMotywu } from './most-motywu';
