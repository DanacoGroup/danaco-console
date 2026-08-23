/**
 * Katalog gniazd układu okien równoległych.
 *
 * Jedna odpowiedzialność: ustalenie, ile gniazd ma układ i jak się nazywają.
 * Identyfikator gniazda jest stały przez całe życie układu — rola, model
 * i katalogi zmieniają się w oknie, ale gniazdo pozostaje tym samym miejscem
 * na scenie. Dzięki temu `nadajRole` i `pokazPrzekazanie` przyjmują wprost
 * identyfikator, bez pośrednictwa indeksu tablicy.
 */

/** Najmniejsza liczba okien na scenie — scena nigdy nie jest pusta. */
export const LICZBA_MIN = 1;

/**
 * Największa liczba okien obok siebie.
 *
 * Cztery, bo tyle liczy obsada multitaskingu: koordynator, dwóch wykonawców
 * (`LICZBA_WYKONAWCOW = 2`) i analityk. Mniejszy sufit nie mieści pełnej pętli
 * — analityk nie ma gdzie stanąć obok pary, którą ocenia.
 */
export const LICZBA_MAX = 4;

/** Identyfikator gniazda układu. */
export type IdGniazda = 'okno-1' | 'okno-2' | 'okno-3' | 'okno-4';

/** Gniazda w kolejności widocznej na scenie, od lewej. */
export const ID_GNIAZD: readonly IdGniazda[] = ['okno-1', 'okno-2', 'okno-3', 'okno-4'];

/** Numer gniazda widoczny dla operatora — liczony od jedynki, nie od zera. */
export function numerGniazda(id: IdGniazda): number {
  return ID_GNIAZD.indexOf(id) + 1;
}

/** Rozstrzyga, czy dowolny tekst jest identyfikatorem gniazda układu. */
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
