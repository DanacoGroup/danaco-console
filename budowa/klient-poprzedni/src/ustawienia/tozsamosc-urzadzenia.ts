/**
 * Tożsamość urządzenia, na którym stoi ta przeglądarka.
 *
 * Wartość jest potrzebna w dwóch trybach. Sekcja „Uwierzytelnianie" zakłada
 * PIN, więc identyfikator musi powstać, jeśli go nie ma. Czytelnik, który
 * wyłącznie pokazuje, czym maszyna się przedstawia, niczego nadawać nie może:
 * samo obejrzenie nie jest czynnością Operatora na urządzeniu.
 *
 * Jedna funkcja z domyślnym trybem byłaby pułapką: pomyłka w argumencie
 * nadałaby tożsamość przy odczycie. Stąd dwie nazwy i jedno miejsce, w którym
 * stoi klucz pamięci.
 *
 * Sekcja „Urządzenia" po tę wartość nie sięga: wykaz `device.list` niesie pole
 * `current`, którym rdzeń sam wskazuje maszynę pytającego. Tożsamość
 * z pamięci przeglądarki byłaby przy nim drugą, słabszą prawdą.
 *
 * Wartość mieszka w pamięci przeglądarki, nie w powitaniu połączenia:
 * PIN jest właściwy urządzeniu, a tożsamość z `connection.hello` jest nadawana
 * na czas uruchomienia. PIN założony wczoraj ma zostać PIN-em tej samej
 * maszyny, więc wartość musi przeżyć zamknięcie okna.
 */

/** Klucz zapisu w pamięci trwałej przeglądarki — jedyne jego wystąpienie. */
export const KLUCZ_URZADZENIA = 'danaco.urzadzenie';

/**
 * Odczytuje tożsamość urządzenia bez nadawania jej.
 *
 * Pusty wynik znaczy „ta przeglądarka jeszcze się nie przedstawiła" i jest
 * prawdą, którą sekcja ma pokazać wprost — a nie powodem, żeby tożsamość
 * wymyślić przy okazji odczytu.
 */
export function odczytajTozsamoscUrzadzenia(): string {
  try {
    return window.localStorage.getItem(KLUCZ_URZADZENIA) ?? '';
  } catch {
    // Pamięć niedostępna (tryb prywatny, zablokowane ciasteczka) — brak
    // tożsamości nie zatrzymuje okna.
    return '';
  }
}

/**
 * Oddaje tożsamość urządzenia, nadając ją przy pierwszym wywołaniu.
 *
 * Pusty wynik znaczy, że pamięć trwała jest niedostępna — wywołanie idzie
 * wtedy bez wymyślonego identyfikatora i to rdzeń orzeka odmowę, zamiast
 * klienta zgadującego za niego.
 */
export function tozsamoscUrzadzenia(): string {
  try {
    const zapisane = odczytajTozsamoscUrzadzenia();
    if (zapisane !== '') return zapisane;
    const nowe = `urzadzenie-${window.crypto.randomUUID()}`;
    window.localStorage.setItem(KLUCZ_URZADZENIA, nowe);
    return nowe;
  } catch {
    return '';
  }
}
