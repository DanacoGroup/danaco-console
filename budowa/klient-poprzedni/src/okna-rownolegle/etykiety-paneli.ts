/**
 * Napisy sterowania panelami z nagłówka okna rozmowy — menu, wykaz paneli, zdanie o brakach i rząd skrótów — trzymają w jednym miejscu etykietę przycisku menu, będącą nazwą czynności, nie narzędzia.
 */
export const MENU_ETYKIETA = 'Panele obok rozmowy';

/** Podpowiedź przycisku otwierającego menu paneli, pokazywana wtedy, gdy menu jest jeszcze zupełnie zwinięte. */
export const MENU_TYTUL_ZWINIETE = 'Rozwiń wykaz paneli, które można otworzyć obok tej rozmowy';

/** Podpowiedź przycisku otwierającego menu paneli, pokazywana wtedy, gdy menu jest już całkiem rozwinięte. */
export const MENU_TYTUL_ROZWINIETE = 'Zwiń wykaz paneli';

/** Nagłówek wykazu pozycji wewnątrz menu paneli, stawiany nad wierszami otwieralnymi tej samej rozmowy. */
export const WYKAZ_NAGLOWEK = 'Panele tej rozmowy';

/**
 * Zdanie pod wykazem paneli mówi o braku, nie o zakazie: menu nie stawia wierszy nieczynnych, więc jedno zdanie nazywa, ilu pozycji spisu nie da się dziś otworzyć.
 */
export function zdanieOBrakach(nieotwieralne: number): string {
  return `Spis niesie ${pozycjiDopelniacz(nieotwieralne)} więcej, ${ktorychNieMozna(nieotwieralne)} `
    + 'dziś otworzyć — okno albo nie ma komendy w rdzeniu, albo powstaje w innej pracy. '
    + 'Pełny spis wraz z powodami stoi w pasie okien pomocniczych modułu.';
}

/** Etykieta skrótu ikonowego prowadzącego do panelu; czynność zależy od tego, czy panel już stoi obok rozmowy. */
export function etykietaSkrotu(nazwa: string, otwarty: boolean): string {
  return otwarty ? `Zamknij panel ${nazwa}` : `Otwórz panel ${nazwa}`;
}

/** Etykieta pozycji wykazu w menu paneli, czytana przez technologie wspomagające zamiast samego napisu wiersza. */
export function etykietaPozycji(nazwa: string, otwarty: boolean): string {
  return otwarty ? `${nazwa} — stoi obok rozmowy` : `${nazwa} — nie stoi obok rozmowy`;
}

/** Nazwa rzędu skrótów do paneli używana przy odczycie ekranowym; rząd jest grupą przycisków, nie paskiem narzędzi. */
export const SKROTY_ETYKIETA = 'Skróty do paneli';

/** Etykieta uchwytu rzędu skrótów, pokazywana wtedy, gdy rząd skrótów jest zwinięty i można go rozwinąć. */
export const UCHWYT_ZWINIETY = 'Rozwiń';

/** Etykieta uchwytu rzędu skrótów, pokazywana wtedy, gdy rząd skrótów jest rozwinięty i można go zwinąć. */
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

/** Liczba pozycji wypisana w bierniku — forma dopasowana do zdania mówiącego, ilu pozycji spisu nie da się dziś otworzyć. */
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

/** Zaimek zgodny liczbą z pozycją — forma którego albo których, dopasowana do liczby pozycji niedostępnych. */
function ktorychNieMozna(liczba: number): string {
  return liczba === 1 ? 'której nie można' : 'których nie można';
}
