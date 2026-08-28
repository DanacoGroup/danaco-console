/**
 * Objaśnienia pozycji katalogu ustawień.
 *
 * Sam dymek [?] mieszka w `komponenty/dymek.ts` i jest stąd re-eksportowany,
 * żeby okno konfiguracji miało jedno wejście. Własne tego pliku jest tylko
 * składanie zdania objaśnienia z pól katalogu ustawień.
 */
export { utworzDymekObjasnienia } from '../komponenty/dymek';

/**
 * Objaśnienie pozycji katalogu ustawień złożone z tego, co katalog o niej
 * mówi: opisu, jednostki, poziomu domyślnego i wymogu ponownego uruchomienia.
 *
 * Zdania są dopisywane, a nie zastępowane: pozycja bez opisu nadal ma czym
 * objaśnić swój klucz.
 */
export function zlozObjasnienie(czlony: readonly (string | undefined)[]): string {
  return czlony
    .map((czlon) => (czlon ?? '').trim())
    .filter((czlon) => czlon !== '')
    .join(' ');
}
