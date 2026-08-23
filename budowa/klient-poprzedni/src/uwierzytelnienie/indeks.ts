/**
 * Moduł uwierzytelnienia — bramka Operatora (rodzina `auth.*`).
 *
 * Powierzchnia publiczna katalogu: ekran logowania dla punktu wejścia
 * aplikacji oraz trwałość sesji bramki dla warstw, które po wejściu potrzebują
 * tokenu. Źródło komend i odmowy pozostają wewnątrz katalogu — żaden inny moduł
 * nie ma powodu wołać `auth.login` z pominięciem ekranu.
 *
 * Wpięcie ekranu należy do punktu wejścia (`main.ts`): przesłona montuje się
 * po złożeniu aplikacji, a przed pierwszym uruchomieniem tras.
 */
export { utworzEkranLogowania } from './ekran-logowania';
export type { EkranLogowania, OpisEkranuLogowania } from './ekran-logowania';
export { odczytajSesje, opisWaznosci, sesjaTrwala } from './sesja-bramki';
/**
 * Tożsamość maszyny — sam odczyt. Bramka identyfikatora nie nadaje: powstaje on
 * przy zakładaniu PIN-u w Ustawieniach, a ekran logowania tylko sprawdza, czy
 * jest czym wskazać maszynę. Ten sam klucz odczytuje
 * `ustawienia/tozsamosc-urzadzenia.ts`.
 */
export { odczytajTozsamoscUrzadzenia, KLUCZ_URZADZENIA } from './tozsamosc-urzadzenia';
