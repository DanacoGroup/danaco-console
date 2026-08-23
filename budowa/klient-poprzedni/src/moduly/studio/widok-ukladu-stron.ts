import type { TrybPrzewijania, UkladKartek } from './widok-nastawy-operatora';

/**
 * Układ kartek na powierzchni — jedna, wiele obok siebie, rozkładówka.
 *
 * ── Trzy rzeczy, których nie wolno pomylić ──────────────────────────────────
 * 1. „Widok dwóch stron" to dwie strony TEGO SAMEGO dokumentu obok siebie i stoi
 *    tutaj. 2. „Dwa dokumenty obok siebie" to podział powierzchni i stoi
 *    w `widok-podzialu-powierzchni.ts`. 3. Rozkładówka to nie „dwie strony obok
 *    siebie": to książka, w której strona pierwsza stoi SAMA po prawej, a dalej
 *    idą pary parzysta–nieparzysta. Bez tego rozkładówka pokazywałaby parę 1–2,
 *    której w oprawionym pismie nigdy nie widać naraz.
 *
 * Plik nie zna DOM: bierze liczbę stron i nastawy, oddaje rzędy numerów.
 */

/** Miejsce w rzędzie: numer strony albo puste miejsce rozkładówki. */
export type MiejsceWRzedzie = number | null;

/** Rzędy kartek — po jednym wierszu na rząd. */
export type RzedyKartek = readonly (readonly MiejsceWRzedzie[])[];

/**
 * Rozkłada strony na rzędy.
 *
 * @param liczbaStron Liczba stron po podziale; zero oddaje jeden rząd pusty,
 *   bo dokument bez treści ma jedną kartkę gotową do pisania, a nie zero kartek.
 * @param uklad Wybór Operatora.
 * @param kartekWRzedzie Liczba kartek w rzędzie przy układzie „obok siebie".
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
    // Strona pierwsza jest stroną otwarcia: stoi sama, po prawej. Puste miejsce
    // po lewej jest miejscem okładki, nie stroną — dlatego `null`, nie numer.
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

/** Liczba kolumn, jaką układ zajmuje — do arkusza siatki. */
export function kolumnyUkladu(uklad: UkladKartek, kartekWRzedzie: number): number {
  if (uklad === 'jedna') return 1;
  if (uklad === 'rozkladowka') return 2;
  return Math.max(1, Math.min(8, Math.round(kartekWRzedzie)));
}

/**
 * Strona rozkładówki, na której stoi kartka — do marginesów odbicia.
 *
 * Strona nieparzysta jest prawą stroną rozkładówki. Zgadza się to
 * z `marginesyKartki`, gdzie margines wewnętrzny strony nieparzystej stoi po
 * lewej — inaczej oprawa rysowałaby się po przeciwnej stronie rowka.
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

/** Przewinięcie, przy którym górna krawędź wskazanej kartki stoi u góry widoku. */
export function przewiniecieDoKartki(
  numerStrony: number,
  wysokoscKartki: number,
  odstep: number,
): number {
  const numer = Math.max(1, Math.floor(numerStrony));
  return (numer - 1) * (wysokoscKartki + odstep);
}

/** Zdanie o układzie kartek — do paska stanu powierzchni. */
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
