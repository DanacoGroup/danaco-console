import type { TrybPrzewijania, UkladKartek } from './widok-nastawy-operatora';

/** Typ MiejsceWRzedzie opisuje jedno miejsce w rzędzie układu kartek: numer strony obecnej w tym miejscu albo wartość pustą, gdy rozkładówka zostawia miejsce bez strony. */
export type MiejsceWRzedzie = number | null;

/** Typ RzedyKartek opisuje układ kartek powierzchni dokumentu jako listę rzędów, w której każdy rząd jest listą miejsc typu MiejsceWRzedzie. */
export type RzedyKartek = readonly (readonly MiejsceWRzedzie[])[];

/**
 * Funkcja rozlozKartkiWRzedy rozkłada strony dokumentu na rzędy kartek według wybranego układu, zwracając dla dokumentu bez stron jeden rząd pusty zamiast rzędów zero.
 */
export function rozlozKartkiWRzedy(
  liczbaStron: number,
  uklad: UkladKartek,
  kartekWRzedzie: number,
): RzedyKartek {
  const ile = Math.max(1, Math.floor(liczbaStron));
  if (uklad === 'jedna') {
    return Array.from({ length: ile }, (_, numer) => [numer + 1]);
  }
  if (uklad === 'rozkladowka') {
    // Strona pierwsza jest stroną otwarcia i stoi sama po prawej; puste miejsce po lewej to okładka.
    const rzedy: MiejsceWRzedzie[][] = [[null, 1]];
    for (let numer = 2; numer <= ile; numer += 2) {
      rzedy.push(numer + 1 <= ile ? [numer, numer + 1] : [numer, null]);
    }
    return rzedy;
  }
  const wRzedzie = Math.max(1, Math.min(8, Math.round(kartekWRzedzie)));
  const rzedy: MiejsceWRzedzie[][] = [];
  for (let numer = 1; numer <= ile; numer += wRzedzie) {
    const rzad: MiejsceWRzedzie[] = [];
    for (let przesuniecie = 0; przesuniecie < wRzedzie; przesuniecie += 1) {
      const strona = numer + przesuniecie;
      rzad.push(strona <= ile ? strona : null);
    }
    rzedy.push(rzad);
  }
  return rzedy;
}

/** Funkcja kolumnyUkladu zwraca liczbę kolumn, jaką wybrany układ kartek zajmuje na powierzchni dokumentu, do zastosowania w arkuszu stylów siatki. */
export function kolumnyUkladu(uklad: UkladKartek, kartekWRzedzie: number): number {
  if (uklad === 'jedna') return 1;
  if (uklad === 'rozkladowka') return 2;
  return Math.max(1, Math.min(8, Math.round(kartekWRzedzie)));
}

/**
 * Funkcja stronaRozkladowki zwraca stronę rozkładówki, na której stoi kartka o podanym numerze: strona nieparzysta jest prawą stroną, a parzysta lewą, zgodnie z marginesyKartki.
 */
export function stronaRozkladowki(numerStrony: number): 'lewa' | 'prawa' {
  return numerStrony % 2 === 0 ? 'lewa' : 'prawa';
}

/**
 * Numer kartki widocznej przy tym przewinięciu.
 *
 * Liczony od górnej krawędzi pola widoku, nie od jego środka: Operator, który
 * przewinął o jedną kartkę, ma widzieć jej górę. Wynik jest zawsze w zakresie
 * od 1 do liczby stron.
 */
export function kartkaPrzyPrzewinieciu(
  przewiniecie: number,
  wysokoscKartki: number,
  odstep: number,
  liczbaStron: number,
): number {
  const ile = Math.max(1, Math.floor(liczbaStron));
  const skok = wysokoscKartki + odstep;
  if (skok <= 0) return 1;
  const numer = Math.floor(Math.max(0, przewiniecie) / skok) + 1;
  return Math.min(Math.max(numer, 1), ile);
}

/** Funkcja przewiniecieDoKartki zwraca przewinięcie, przy którym górna krawędź wskazanej kartki styka się z górną krawędzią widoku powierzchni dokumentu. */
export function przewiniecieDoKartki(
  numerStrony: number,
  wysokoscKartki: number,
  odstep: number,
): number {
  const numer = Math.max(1, Math.floor(numerStrony));
  return (numer - 1) * (wysokoscKartki + odstep);
}

/** Funkcja opiszUkladKartek zwraca zdanie opisujące bieżący układ kartek, tryb przewijania i liczbę stron do wyświetlenia w pasku stanu powierzchni. */
export function opiszUkladKartek(
  uklad: UkladKartek,
  kartekWRzedzie: number,
  przewijanie: TrybPrzewijania,
  liczbaStron: number,
): string {
  const nazwaUkladu =
    uklad === 'rozkladowka'
      ? 'rozkładówka (strona pierwsza sama po prawej)'
      : uklad === 'obok'
        ? `${kolumnyUkladu(uklad, kartekWRzedzie)} kartki w rzędzie`
        : 'jedna kartka w rzędzie';
  const nazwaPrzewijania =
    przewijanie === 'ciagle' ? 'przewijanie ciągłe' : 'przewijanie strona po stronie';
  return `${nazwaUkladu} · ${nazwaPrzewijania} · stron ${Math.max(1, liczbaStron)}`;
}
