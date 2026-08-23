/**
 * Słownik niepowodzeń modułu — mówi, kto zawiódł, zamiast „coś poszło nie tak".
 *
 * `Wynik` z polem `udany === false` powstaje w warstwie protokołu na więcej niż
 * jeden sposób i tylko jeden z nich jest odmową rdzenia:
 *
 *  - odmowa rdzenia — koperta ze `status: "error"` i `ErrorInfo`;
 *  - odpowiedź nieczytelna — rdzeń odpowiedział powodzeniem, ale treść nie ma
 *    kształtu z kontraktu, więc `protokol/ksztalt-odpowiedzi.ts` zamienia ją
 *    na niepowodzenie z kodem `validation_failed`.
 *
 * Nazwanie odpowiedzi nieczytelnej odmową byłoby oskarżeniem rdzenia o czyn,
 * którego nie popełnił: rdzeń odpowiedział i odpowiedział powodzeniem, a nie
 * zrozumiał go klient. Okno nie ma prawa twierdzić o braku, którego rdzeń nie
 * orzekł — tak samo jak nie ma prawa zamienić odmowy w pustkę.
 *
 * Stan `odpowiedz-bez-tresci` jest strażą, nie drogą osiągalną dziś.
 * `protokol/wynik-czastkowy.ts` oddaje niepowodzenie bez pola `blad`, gdy
 * odpowiedź jest udana, lecz pusta; do `przenies` taka odpowiedź nie dochodzi,
 * bo `sprawdzKsztalt` przechwytuje ją wcześniej i okno mówi wtedy zdanie
 * o odpowiedzi nieczytelnej. Stan zostaje w słowniku na wypadek zmiany warstwy
 * protokołu.
 *
 * Zdania stoją tu, a nie w oknach, z tego samego powodu, dla którego stan treści
 * stoi w `komponenty/stan-tresci.ts`: cztery okna mówiące o tej samej ciszy
 * czterema zdaniami rozjeżdżają się przy pierwszej poprawce.
 */

/** Skąd wzięło się niepowodzenie odczytu — więcej niż jedna droga. */
export type PowodNiepowodzenia =
  /** Rdzeń orzekł odmowę: koperta błędu wraz z kodem kontraktu. */
  | 'odmowa-rdzenia'
  /** Rdzeń odpowiedział powodzeniem, lecz treść nie ma kształtu z kontraktu. */
  | 'odpowiedz-nieczytelna'
  /** Straż: odpowiedź udana bez treści. Dziś nieosiągalna — patrz opis pliku. */
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
