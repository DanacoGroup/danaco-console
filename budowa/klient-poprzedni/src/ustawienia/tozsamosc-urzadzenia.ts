/** Klucz zapisu tożsamości urządzenia w pamięci trwałej przeglądarki jest jedynym jego wystąpieniem, potrzebnym w dwóch trybach: zakładania PIN i samego odczytu bez nadawania. */
export const KLUCZ_URZADZENIA = 'danaco.urzadzenie';

/** Odczytuje tożsamość urządzenia bez nadawania jej; pusty wynik znaczy, że przeglądarka jeszcze się nie przedstawiła. */
export function odczytajTozsamoscUrzadzenia(): string {
  try {
    return window.localStorage.getItem(KLUCZ_URZADZENIA) ?? '';
  } catch {
    // Pamięć niedostępna w trybie prywatnym albo przy zablokowanych ciasteczkach.
    return '';
  }
}

/** Oddaje tożsamość urządzenia, nadając ją przy pierwszym wywołaniu, gdy pamięć trwała jeszcze jej nie niesie. */
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
