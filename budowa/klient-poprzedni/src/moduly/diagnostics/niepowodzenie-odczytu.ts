/**
 * Słownik niepowodzeń odczytu modułu Diagnostics, który nazywa stronę
 * odpowiedzialną, zamiast podawać ogólnikowe zdanie o awarii. Rozróżnia odmowę
 * rdzenia od odpowiedzi udanej, której klient nie potrafi odczytać.
 */

/**
 * Powód niepowodzenia odczytu. Niepowodzenie powstaje w warstwie protokołu
 * więcej niż jedną drogą, a tylko jedna z nich jest odmową rdzenia, dlatego
 * każda droga nosi tutaj osobną nazwę.
 */
export type PowodNiepowodzenia =
  /** Rdzeń orzekł odmowę: koperta błędu wraz z kodem kontraktu. */
  | 'odmowa-rdzenia'
  /** Rdzeń odpowiedział powodzeniem, lecz treść nie ma kształtu z kontraktu. */
  | 'odpowiedz-nieczytelna'
  /** Straż: odpowiedź udana bez treści, dziś nieosiągalna w warstwie protokołu. */
  | 'odpowiedz-bez-tresci';

/**
 * Zdanie o nieudanej czynności. Powód nieznany nie staje się oskarżeniem.
 *
 * @param czynnosc dopełniacz czynności — „odczytu błędów", „uruchomienia analizy".
 */
export function zdanieNiepowodzenia(czynnosc: string, powod?: PowodNiepowodzenia): string {
  switch (powod) {
    case 'odmowa-rdzenia':
      return `Rdzeń odmówił ${czynnosc}.`;
    case 'odpowiedz-nieczytelna':
      return `Rdzeń NIE odmówił ${czynnosc} — odpowiedział, ale jego odpowiedzi nie da się odczytać: kształt niezgodny z kontraktem.`;
    case 'odpowiedz-bez-tresci':
      return `Rdzeń przyjął żądanie ${czynnosc} i nie oddał ani treści, ani powodu — okno nie ma czego pokazać i nie zgaduje dlaczego.`;
    default:
      // Wywołujący nie podał drogi. Okno mówi wtedy o sobie, nie o rdzeniu.
      return `Czynność ${czynnosc} nie doszła do skutku, a okno nie ustaliło dlaczego — o rdzeniu nic tu nie orzeka.`;
  }
}
