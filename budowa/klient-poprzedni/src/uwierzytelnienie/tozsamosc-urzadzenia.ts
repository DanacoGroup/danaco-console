/** Klucz zapisu tożsamości maszyny w pamięci trwałej przeglądarki, nadawany przy zakładaniu metody PIN i czytany wyłącznie przez ekran logowania. */
export const KLUCZ_URZADZENIA = 'danaco.urzadzenie';

/** Zwraca identyfikator tej maszyny, jeżeli PIN był już założony, albo pusty napis, gdy pamięć przeglądarki jest niedostępna lub tożsamości brak. */
export function odczytajTozsamoscUrzadzenia(): string {
  try {
    return window.localStorage.getItem(KLUCZ_URZADZENIA) ?? '';
  } catch {
    return '';
  }
}
