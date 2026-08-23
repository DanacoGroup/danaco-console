/**
 * Pamięć zwinięcia stref i sekcji — jedna dla obu powierzchni wejściowych.
 *
 * Strefy rozwijają się na żądanie, a raz wykonane rozwinięcie ma się utrzymać
 * między wejściami na ekran.
 *
 * Nastawa siedzi w `localStorage`, nie w rdzeniu: dotyczy powierzchni na tym
 * urządzeniu i nie ma swojego bytu w kontrakcie. Ten sam wzorzec nosi
 * `motyw/motyw.ts` i `okna-rownolegle/kolejnosc-miejscowa.ts`.
 *
 * `localStorage` bywa niedostępny (tryb prywatny, osadzenie w ramce). Awaria
 * odczytu albo zapisu zostaje przy wartości domyślnej i nie jest zgłaszana jako
 * błąd — nastawa widoku nie jest powodem, żeby ekran nie wstał.
 */

/** Przedrostek klucza — nastawy widoku nie mieszają się z danymi sesji. */
const PRZEDROSTEK = 'dn.zwiniecie.';

/**
 * Czy strefa o podanym kluczu ma być rozwinięta.
 *
 * @param klucz Nazwa strefy, stała między wejściami.
 * @param domyslnie Postać przed pierwszym zapamiętanym zwinięciem.
 */
export function czyRozwiniete(klucz: string, domyslnie: boolean): boolean {
  try {
    const zapis = globalThis.localStorage?.getItem(PRZEDROSTEK + klucz);
    if (zapis === null || zapis === undefined) return domyslnie;
    return zapis === '1';
  } catch {
    return domyslnie;
  }
}

/** Zapamiętuje postać strefy na tym urządzeniu. */
export function zapamietajZwiniecie(klucz: string, rozwiniete: boolean): void {
  try {
    globalThis.localStorage?.setItem(PRZEDROSTEK + klucz, rozwiniete ? '1' : '0');
  } catch {
    // Nastawa widoku nie jest powodem, żeby cokolwiek przerywać.
  }
}
