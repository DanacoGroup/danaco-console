/**
 * Przycięcie wartości, która jedzie do rdzenia jako kod — etykieta zasobu, pole
 * promptu, nazwa kompozycji. Jedno rozstrzygnięcie na cały moduł o tym, co jest
 * w takiej wartości brzegową pustką.
 */

/**
 * Znaki brzegowe: `\s` JavaScriptu (niesie już U+FEFF i wszystkie Zs) wraz
 * z U+0085, którego produkcja WhiteSpace nie zna, a `unicode.IsSpace` Go zna.
 * Suma obu zbiorów — nie część wspólna, bo to strona ostrożniejsza wygrywa.
 */
const BRZEGOWE = /^[\s\u0085]+|[\s\u0085]+$/gu;

/**
 * Wartość bez brzegów białych po obu notacjach — przeglądarki i rdzenia; środek
 * wartości pozostaje nietknięty, ponieważ jest to przycięcie brzegów, a nie
 * czyszczenie treści.
 */
export function przytnijKod(tekst: string): string {
  return tekst.replace(BRZEGOWE, '');
}

/**
 * Rozbija wpis operatora na zestaw etykiet: bez wartości pustych i bez
 * powtórzeń. Ta sama funkcja służy nadaniu etykiet i filtrowi wykazu, ponieważ
 * etykieta nadana i szukana muszą być tym samym ciągiem znaków.
 */
export function rozbijEtykiety(wpis: string): readonly string[] {
  const zebrane = wpis
    .split(',')
    .map((etykieta) => przytnijKod(etykieta))
    .filter((etykieta) => etykieta !== '');
  return [...new Set(zebrane)];
}
