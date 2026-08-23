/**
 * Przycięcie wartości, która jedzie do rdzenia jako kod — etykieta zasobu,
 * pole promptu, nazwa kompozycji.
 *
 * Jedna odpowiedzialność: jedno rozstrzygnięcie na cały moduł o tym, co jest
 * w takiej wartości brzegową pustką.
 *
 * Samo `trim()` przeglądarki tu nie wystarcza, bo `trim()` JavaScriptu
 * i `strings.TrimSpace` języka Go nie zdejmują tego samego: produkcja
 * WhiteSpace JavaScriptu nie zna znaku NEL (U+0085), a `unicode.IsSpace` Go go
 * zna; odwrotnie `trim()` zdejmuje ZWNBSP (U+FEFF), którego Go nie rusza. Rdzeń
 * etykiet nie przycina w ogóle (`dane/design_zasoby.go` odrzuca wyłącznie
 * `etykieta == ""`), więc etykieta złożona z samego NEL zapisałaby się jako
 * niewidoczna, a etykieta z NEL na końcu nie dałaby się odnaleźć filtrem
 * wpisanym bez tego znaku.
 *
 * Brzegową pustką jest tutaj wszystko, co za biały znak uważa którakolwiek ze
 * stron. Przycinamy sumą obu zbiorów, więc do rdzenia jedzie to, co człowiek
 * uznałby za wpisane — tak samo po stronie zapisu (nadanie etykiet), jak po
 * stronie odczytu (filtr wykazu), więc jedno trafia w drugie. Środka wartości
 * nie ruszamy: to jest przycięcie brzegów, a nie czyszczenie treści.
 */

/**
 * Znaki brzegowe: `\s` JavaScriptu (niesie już U+FEFF i wszystkie Zs) wraz
 * z U+0085, którego produkcja WhiteSpace nie zna, a `unicode.IsSpace` Go zna.
 * Suma obu zbiorów — nie część wspólna, bo to strona ostrożniejsza wygrywa.
 */
const BRZEGOWE = /^[\s\u0085]+|[\s\u0085]+$/gu;

/** Wartość bez brzegów białych po obu notacjach — przeglądarki i rdzenia. */
export function przytnijKod(tekst: string): string {
  return tekst.replace(BRZEGOWE, '');
}

/**
 * Rozbija wpis Operatora na zestaw etykiet: bez wartości pustych i bez
 * powtórzeń, bo rdzeń i tak nie zapisze etykiety dwa razy (klucz główny tabeli
 * `etykieta_zasobu_design` obejmuje etykietę), a odpowiedź niosłaby wtedy co
 * innego niż żądanie.
 *
 * Jedna funkcja na nadanie etykiet i na filtr wykazu — bo etykieta nadana
 * i etykieta szukana muszą być tym samym ciągiem znaków, inaczej filtr nie
 * trafia we własny zapis.
 */
export function rozbijEtykiety(wpis: string): readonly string[] {
  const zebrane = wpis
    .split(',')
    .map((etykieta) => przytnijKod(etykieta))
    .filter((etykieta) => etykieta !== '');
  return [...new Set(zebrane)];
}
