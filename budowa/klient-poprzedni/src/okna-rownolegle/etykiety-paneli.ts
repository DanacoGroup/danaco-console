/**
 * Napisy sterowania panelami z nagłówka okna rozmowy — menu `⋮`, wykaz paneli,
 * zdanie o brakach i rząd skrótów w jednym miejscu, żeby ta sama czynność nie
 * nazywała się w dwóch plikach inaczej.
 *
 * Osobno od `etykiety-ukladu.ts`, bo tamten plik mówi o scenie: liczbie okien,
 * rolach i pasie relacji. Panele należą do rozmowy, nie do sceny.
 *
 * Stoją tu wyłącznie napisy oprawy — nazwy i przeznaczenia paneli przychodzą ze
 * spisu okien pomocniczych. Konwencja jak w `etykiety-ukladu.ts`: stałe
 * wersalikowe dla napisów stałych, funkcje dla zdań składanych, polska odmiana
 * rozpisana przypadkami zamiast doklejania końcówek.
 */

/** Etykieta przycisku `⋮` — czynność, nie nazwa narzędzia. */
export const MENU_ETYKIETA = 'Panele obok rozmowy';

/** Podpowiedź przycisku `⋮` przy zwiniętym menu. */
export const MENU_TYTUL_ZWINIETE = 'Rozwiń wykaz paneli, które można otworzyć obok tej rozmowy';

/** Podpowiedź przycisku `⋮` przy rozwiniętym menu. */
export const MENU_TYTUL_ROZWINIETE = 'Zwiń wykaz paneli';

/** Nagłówek wykazu pozycji wewnątrz menu. */
export const WYKAZ_NAGLOWEK = 'Panele tej rozmowy';

/**
 * Zdanie pod wykazem — mówi o braku, nie o zakazie.
 *
 * Menu nie stawia wierszy nieczynnych: brak dostępnej pozycji daje krótszą
 * listę, nie zablokowany wiersz. Brak ma jednak zostać nazwany, więc jedno
 * zdanie podaje, ilu pozycji spisu nie da się dziś otworzyć. Przy zerze zdania
 * nie ma wcale — tym zajmuje się `menu-paneli.ts`.
 */
export function zdanieOBrakach(nieotwieralne: number): string {
  return `Spis niesie ${pozycjiDopelniacz(nieotwieralne)} więcej, ${ktorychNieMozna(nieotwieralne)} `
    + 'dziś otworzyć — okno albo nie ma komendy w rdzeniu, albo powstaje w innej pracy. '
    + 'Pełny spis wraz z powodami stoi w pasie okien pomocniczych modułu.';
}

/** Etykieta skrótu ikonowego: czynność zależy od tego, czy panel już stoi. */
export function etykietaSkrotu(nazwa: string, otwarty: boolean): string {
  return otwarty ? `Zamknij panel ${nazwa}` : `Otwórz panel ${nazwa}`;
}

/** Etykieta pozycji wykazu czytana przez technologie wspomagające. */
export function etykietaPozycji(nazwa: string, otwarty: boolean): string {
  return otwarty ? `${nazwa} — stoi obok rozmowy` : `${nazwa} — nie stoi obok rozmowy`;
}

/** Nazwa rzędu skrótów dla odczytu — rząd jest grupą, nie paskiem narzędzi. */
export const SKROTY_ETYKIETA = 'Skróty do paneli';

/** Etykieta uchwytu przy zwiniętym bycie. */
export const UCHWYT_ZWINIETY = 'Rozwiń';

/** Etykieta uchwytu przy rozwiniętym bycie. */
export const UCHWYT_ROZWINIETY = 'Zwiń';

/**
 * Komunikat, gdy ani jednej pozycji nie da się dziś otworzyć.
 *
 * Nie mówi „brak danych" — tłumaczy, czego brakuje i gdzie leży pełny spis.
 */
export const BRAK_KANALU =
  'Żadnego panelu nie da się dziś otworzyć obok tej rozmowy: ani jedna pozycja '
  + 'spisu nie ma jeszcze kanału do rdzenia. Spis wraz z powodami stoi w pasie '
  + 'okien pomocniczych modułu.';

/** „dwie pozycje" — liczba pozycji w bierniku, jak żąda tego zdanie „Spis niesie …". */
function pozycjiDopelniacz(liczba: number): string {
  switch (liczba) {
    case 1:
      return 'jedną pozycję';
    case 2:
      return 'dwie pozycje';
    case 3:
      return 'trzy pozycje';
    case 4:
      return 'cztery pozycje';
    default:
      return `${liczba} pozycji`;
  }
}

/** „której" / „których" — zaimek zgodny z liczbą pozycji. */
function ktorychNieMozna(liczba: number): string {
  return liczba === 1 ? 'której nie można' : 'których nie można';
}
