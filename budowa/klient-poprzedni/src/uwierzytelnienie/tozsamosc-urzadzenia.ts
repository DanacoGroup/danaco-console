/**
 * Tożsamość maszyny — odczyt na potrzeby wejścia PIN-em.
 *
 * PIN jest właściwy maszynie, nie Operatorowi: rdzeń odnajduje metodę po parze
 * rodzaj + urządzenie (`adapter_modul_auth.go`, `metodaWejscia`), więc PIN
 * założony pod jednym identyfikatorem nie otworzy bramki podanym pod innym.
 * Ekran logowania musi znać tę wartość, żeby wiedzieć, czy segment „PIN" ma
 * w ogóle stanąć, i żeby móc PIN wysłać.
 *
 * Bramka wyłącznie czyta identyfikator; nadaje go zakładanie PIN-u
 * w Ustawieniach (`auth.method.add`), czyli czynność wykonana na urządzeniu.
 * Ekran logowania staje przy każdym uruchomieniu, także pierwszym na świeżej
 * przeglądarce — gdyby nadawał tożsamość, zapisywałby maszynę, która żadnego
 * PIN-u nie ma. Pusty wynik jest prawdą: ta przeglądarka jeszcze się nie
 * przedstawiła.
 *
 * Ten sam klucz niesie `ustawienia/tozsamosc-urzadzenia.ts` i tak samo
 * rozdziela odczyt od nadania. Bramka nie importuje pliku Ustawień, bo `auth.*`
 * leży niżej niż okno Ustawień i zależność szłaby pod prąd warstw; scalenie obu
 * plików wymaga warstwy wspólnej poniżej nich.
 */

/** Klucz zapisu tożsamości maszyny w pamięci trwałej przeglądarki. */
export const KLUCZ_URZADZENIA = 'danaco.urzadzenie';

/**
 * Identyfikator tej maszyny, jeżeli już się przedstawiła.
 *
 * @returns identyfikator albo pusty napis — gdy PIN-u tu nie zakładano albo
 *   pamięci przeglądarki nie ma. Wołający ma wtedy o co mniej pytać rdzenia.
 */
export function odczytajTozsamoscUrzadzenia(): string {
  try {
    return window.localStorage.getItem(KLUCZ_URZADZENIA) ?? '';
  } catch {
    return '';
  }
}
