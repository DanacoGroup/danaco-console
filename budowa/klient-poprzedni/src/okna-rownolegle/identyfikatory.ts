/**
 * Katalog gniazd układu okien równoległych ustala, ile gniazd ma układ i jak się nazywają, przy czym najmniejsza dopuszczalna liczba okien na scenie wynosi jeden, bo scena nigdy nie jest pusta.
 */
export const LICZBA_MIN = 1;

/**
 * Największa liczba okien obok siebie.
 *
 * Cztery, bo tyle liczy obsada multitaskingu: koordynator, dwóch wykonawców
 * (`LICZBA_WYKONAWCOW = 2`) i analityk. Mniejszy sufit nie mieści pełnej pętli
 * — analityk nie ma gdzie stanąć obok pary, którą ocenia.
 */
export const LICZBA_MAX = 4;

/** Identyfikator gniazda układu okien równoległych — stały przez całe życie układu, niezależny od roli okna. */
export type IdGniazda = 'okno-1' | 'okno-2' | 'okno-3' | 'okno-4';

/** Gniazda w kolejności widocznej na scenie, od lewej do prawej, w liczbie odpowiadającej granicy górnej układu. */
export const ID_GNIAZD: readonly IdGniazda[] = ['okno-1', 'okno-2', 'okno-3', 'okno-4'];

/** Numer gniazda widoczny dla operatora — liczony od jedynki, nie od zera, wprost z kolejności na scenie. */
export function numerGniazda(id: IdGniazda): number {
  return ID_GNIAZD.indexOf(id) + 1;
}

/** Rozstrzyga, czy dowolny tekst jest identyfikatorem gniazda należącym do tego układu okien równoległych. */
export function czyIdGniazda(tekst: string): tekst is IdGniazda {
  return (ID_GNIAZD as readonly string[]).includes(tekst);
}

/**
 * Liczba okien sprowadzona do zakresu układu.
 *
 * Wartość spoza zakresu nie jest błędem wstrzymującym — zostaje przycięta
 * do najbliższej dopuszczalnej.
 */
export function ograniczLiczbe(liczba: number): number {
  if (!Number.isFinite(liczba)) return LICZBA_MIN;
  return Math.min(LICZBA_MAX, Math.max(LICZBA_MIN, Math.trunc(liczba)));
}
