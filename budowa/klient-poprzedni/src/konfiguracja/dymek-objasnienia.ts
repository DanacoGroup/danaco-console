export { utworzDymekObjasnienia } from '../komponenty/dymek';

/**
 * Objaśnienia pozycji katalogu ustawień.
 *
 * Sam dymek [?] mieszka w bibliotece (`komponenty/dymek.ts`) i jest stąd
 * re-eksportowany, żeby okno konfiguracji miało jedno wejście. Własne tego
 * pliku jest tylko składanie zdania objaśnienia z pól katalogu ustawień.
 */

/**
 * Objaśnienie pozycji katalogu ustawień złożone z tego, co katalog o niej
 * mówi: opisu, jednostki, poziomu domyślnego i wymogu ponownego uruchomienia.
 *
 * Zdania są dopisywane, a nie zastępowane: pozycja bez opisu nadal ma czym
 * objaśnić swój klucz, a pozycja z opisem zyskuje to, czego opis nie mówi.
 */
export function zlozObjasnienie(czlony: readonly (string | undefined)[]): string {
  return czlony
    .map((czlon) => (czlon ?? '').trim())
    .filter((czlon) => czlon !== '')
    .join(' ');
}
