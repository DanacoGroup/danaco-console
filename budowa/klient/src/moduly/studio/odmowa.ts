/**
 * Przekład odmowy serwera na zdanie dla Operatora.
 *
 * Kontrakt opisuje pole `message` jako „opis dla użytkownika", ale serwer
 * wypełnia je językiem swojego wnętrza — nazwami warstw, komend i tabel.
 * Wypisanie go wprost stawia Operatora przed zdaniem, którego nie ma jak
 * zrozumieć ani na które nie ma jak odpowiedzieć. Wiążący jest więc kod
 * odmowy: on jest ustalony kontraktem i niesie tyle, ile potrzeba, żeby
 * powiedzieć, co zrobić dalej.
 *
 * Opis od serwera nie ginie — idzie do dziennika przeglądarki, gdzie służy
 * rozpoznaniu usterki.
 *
 * Ten sam przekład stoi w drodze wejścia (`wejscie/tresci.ts`, gałąź
 * `usterki.odmowaSerwera.wedlugKodu`). Osobno, bo katalogi treści drogi
 * wejścia i modułu są rozdzielone z założenia; brzmienie zdań jest wspólne.
 */

import { ErrorCode, type ErrorInfo } from '../../../../shared/contract.ts';

/** Zdanie na każdy kod z katalogu kontraktu. Brak pozycji zostawiłby Operatora bez wyjaśnienia. */
const ZDANIE: Readonly<Record<ErrorCode, string>> = {
  [ErrorCode.ValidationFailed]: 'Podane dane są nieprawidłowe. Popraw je i spróbuj ponownie.',
  [ErrorCode.NotFound]: 'Nie znaleziono wskazanej pozycji. Odśwież widok i spróbuj ponownie.',
  [ErrorCode.NotAuthenticated]: 'Sesja nie jest zalogowana. Zaloguj się ponownie.',
  [ErrorCode.PermissionDenied]: 'Brak uprawnień do tej czynności. Skontaktuj się z administratorem.',
  [ErrorCode.Conflict]: 'Nie można wykonać tej czynności w obecnym stanie. Odśwież widok i sprawdź, czy jest nadal potrzebna.',
  [ErrorCode.ChannelUnavailable]: 'Model jest chwilowo niedostępny. Spróbuj ponownie za chwilę.',
  [ErrorCode.RateLimited]: 'Za dużo żądań w krótkim czasie. Odczekaj chwilę i spróbuj ponownie.',
  [ErrorCode.InternalError]: 'Wystąpił błąd po stronie serwera. Spróbuj ponownie. Jeżeli problem się powtarza, skontaktuj się z administratorem.',
};

/** Zdanie ostatniej szansy: odmowa bez kodu albo z kodem spoza katalogu kontraktu. */
const ZDANIE_ZAPASOWE =
  'Nie udało się wykonać czynności. Spróbuj ponownie. Jeżeli problem się powtarza, skontaktuj się z administratorem.';

/**
 * Zdanie dla Operatora na podstawie odmowy. `czynnosc` nazywa miejsce wywołania
 * i trafia wyłącznie do dziennika — pomaga odnaleźć usterkę, gdy Operator
 * zgłasza, że „coś nie działa".
 */
export function opisOdmowy(blad: ErrorInfo | undefined, czynnosc: string): string {
  zaloguj(blad, czynnosc);
  if (blad === undefined) return ZDANIE_ZAPASOWE;
  return ZDANIE[blad.code] ?? ZDANIE_ZAPASOWE;
}

/**
 * Odkłada odmowę do dziennika, nie zwracając nic. Dla miejsc, w których panel
 * ma własne zdanie na tę jedną czynność — trafniejsze niż zdanie ogólne po
 * kodzie — a opis od serwera i tak nie może pójść do widoku.
 */
export function zaloguj(blad: ErrorInfo | undefined, czynnosc: string): void {
  if (blad === undefined) return;
  console.warn(`[studio] odmowa serwera (${czynnosc})`, blad.code, blad.message);
}
