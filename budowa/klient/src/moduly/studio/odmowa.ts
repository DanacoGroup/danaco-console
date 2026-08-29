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

import { ErrorCode, MessageStatus, type ErrorInfo, type Message } from '../../../../shared/contract.ts';
import { tekst } from './narzedzia.ts';

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

/* Wpis zakończony błędem niesie komunikat programu zewnętrznego, po angielsku
   i jego słownictwem. Wypisany wprost stawia Operatora przed zdaniem, którego
   nie ma jak wykonać, więc stan kanału rozpoznaje się po treści. */

/**
 * Znaki rozpoznawcze stanu kanału w treści od programu zewnętrznego, każdy
 * z kluczem zdania w katalogu treści. Tylko brzmienia zmierzone na kanale;
 * brzmienie nieznane idzie gałęzią zapasową, zamiast dostać zdanie o naprawie,
 * która nic nie da.
 */
const ROZPOZNANIE: ReadonlyArray<{ znak: RegExp; klucz: string }> = [
  { znak: /not logged in|\/login\b/i, klucz: 'czat.bladKanalu.niepolaczony' },
  /* Brzmienia limitu wzięte z wykazu, po którym rozpoznaje je serwer
     (`server/internal/injection/wyczerpanie.go`) — nie z domysłu. */
  { znak: /usage limit reached|rate.?limit (exceeded|error)|quota exceeded/i, klucz: 'czat.bladKanalu.limitWyczerpany' },
];

/**
 * Zdanie dla Operatora na podstawie treści wpisu zakończonego błędem.
 * Surowa treść idzie przy okazji do konsoli przeglądarki — tam służy
 * rozpoznaniu usterki, na ekran nie wchodzi.
 */
export function opisBleduKanalu(tresc: string): string {
  zalogujTresc(tresc, 'czat.wpisBledny');
  const zmierzona = tresc.trim();
  if (zmierzona === '') return tekst('czat.bladKanalu.bezOpisu');
  for (const { znak, klucz } of ROZPOZNANIE) {
    if (znak.test(zmierzona)) return tekst(klucz);
  }
  return tekst('czat.bladKanalu.nierozpoznany');
}

/**
 * Treść wpisu do pokazania w oknie rozmowy. Punkt wpięcia dla składnika
 * historii: wpis poprawny oddaje swoją treść bez zmiany, wpis o stanie błędu
 * — zdanie po polsku.
 */
export function trescWpisu(wpis: Message): string {
  if (wpis.status !== MessageStatus.Error) return wpis.content;
  return opisBleduKanalu(wpis.content);
}

/**
 * Odkłada surową treść błędu do dziennika przeglądarki. Osobno od `zaloguj`,
 * bo tamten opisuje odmowę komendy z kodem kontraktu, a tu jest sam tekst
 * od programu zewnętrznego, bez kodu.
 */
export function zalogujTresc(tresc: string, czynnosc: string): void {
  console.warn(`[studio] błąd kanału modelu (${czynnosc})`, tresc);
}
